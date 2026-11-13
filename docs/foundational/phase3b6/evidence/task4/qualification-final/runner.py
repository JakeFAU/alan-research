"""Record task-local qualification commands against immutable executable inputs."""
import concurrent.futures
import gzip
import hashlib
import json
import os
from pathlib import Path
import signal
import subprocess
import sys
import time


def digest(data):
    return hashlib.sha256(data).hexdigest()


def source_map(root, exclusions=()):
    names = subprocess.check_output(["git", "ls-files", "-z"], cwd=root).decode().split("\0")
    result = {}
    for name in names:
        if not name or name.startswith("paper/") or any(name.startswith(prefix) for prefix in exclusions):
            continue
        if name.startswith("docs/") and not (name.startswith(("docs/foundational/phase3b5/evidence/", "docs/foundational/phase3b6/evidence/")) and name.endswith(".py")):
            continue
        if name.endswith((".go", ".mod", ".sum", ".py", ".sh")) or "/testdata/" in name:
            result[name] = digest((root / name).read_bytes())
    return result


def main():
    plan_path = Path(sys.argv[1]).resolve()
    plan = json.loads(plan_path.read_text())
    root = Path(plan["root"])
    output = Path(plan["output"])
    output.mkdir(parents=True, exist_ok=False)
    frozen = source_map(root, plan.get("source_exclusions", ()))
    expected = json.loads(Path(plan["source_manifest"]).read_text())
    if frozen != expected:
        raise RuntimeError("source differs from frozen manifest")
    (output / "plan.json").write_bytes(plan_path.read_bytes())
    (output / "source.json").write_text(json.dumps(frozen, indent=2, sort_keys=True) + "\n")
    (output / "runner.py").write_bytes(Path(__file__).read_bytes())
    identity = {
        "head": subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=root).decode().strip(),
        "tree": subprocess.check_output(["git", "rev-parse", "HEAD^{tree}"], cwd=root).decode().strip(),
        "runner_sha256": digest(Path(__file__).read_bytes()),
        "plan_sha256": digest(plan_path.read_bytes()),
        "source_sha256": digest((output / "source.json").read_bytes()),
    }
    (output / "identity.json").write_text(json.dumps(identity, indent=2) + "\n")

    def run(job):
        name = job["name"]
        environment = dict(plan["environment"])
        environment.update(job.get("environment", {}))
        env = dict(os.environ)
        env.update(environment)
        before = source_map(root, plan.get("source_exclusions", ()))
        if before != frozen:
            raise RuntimeError(name + ": source changed before command")
        log_path = output / (name + ".log")
        record = {"name": name, "argv": job["argv"], "cwd": str(root),
                  "environment": environment, "timeout_seconds": job["timeout"],
                  "expected_exit": job.get("expected_exit", 0)}
        started = time.monotonic()
        print("START " + name, flush=True)
        with log_path.open("wb") as stream:
            process = subprocess.Popen(job["argv"], cwd=root, env=env,
                                       stdout=stream, stderr=subprocess.STDOUT,
                                       start_new_session=True)
            try:
                record["exit"] = process.wait(timeout=job["timeout"])
            except subprocess.TimeoutExpired:
                os.killpg(process.pid, signal.SIGKILL)
                record["exit"] = process.wait()
                record["watchdog_expired"] = True
        record["elapsed_seconds"] = time.monotonic() - started
        raw = log_path.read_bytes()
        record["log_sha256"] = digest(raw)
        record["log_bytes"] = len(raw)
        archive = log_path.with_suffix(".log.gz")
        archive.write_bytes(gzip.compress(raw, mtime=0))
        record["archive_sha256"] = digest(archive.read_bytes())
        log_path.unlink()
        after = source_map(root, plan.get("source_exclusions", ()))
        record["source_unchanged"] = before == after == frozen
        if job.get("go_test_json"):
            events = []
            for line in raw.decode(errors="replace").splitlines():
                if line.startswith("{"):
                    try:
                        events.append(json.loads(line))
                    except json.JSONDecodeError:
                        pass
            record["test_summary"] = {
                "pass_events": sum(e.get("Action") == "pass" and "Test" in e for e in events),
                "package_results": [e for e in events if e.get("Action") in ("pass", "fail", "skip") and "Test" not in e],
                "named_skips": [e for e in events if e.get("Action") == "skip" and "Test" in e],
                "failures": [e for e in events if e.get("Action") == "fail"],
            }
        if job.get("expected_log"):
            expected_log = Path(job["expected_log"]).read_bytes()
            record["baseline_log_sha256"] = digest(expected_log)
            record["baseline_log_byte_equal"] = raw == expected_log
        record["matched"] = (record.get("baseline_log_byte_equal", True) and record["exit"] == record["expected_exit"] and
                             record["source_unchanged"] and not record.get("watchdog_expired", False))
        (output / (name + ".json")).write_text(json.dumps(record, indent=2) + "\n")
        print("END " + name + " exit=" + str(record["exit"]) + " matched=" + str(record["matched"]), flush=True)
        return record

    with concurrent.futures.ThreadPoolExecutor(max_workers=plan.get("parallel", 3)) as pool:
        records = list(pool.map(run, plan["jobs"]))
    (output / "summary.json").write_text(json.dumps(records, indent=2) + "\n")
    return 0 if all(r["matched"] for r in records) else 1


if __name__ == "__main__":
    sys.exit(main())
