"""Run one declared research command, preserving its exact exit and output.

This controller utility does not interpret or rewrite captured evidence.
Its job file fixes cwd, argv, a small environment override and timeout.
"""
import gzip
import hashlib
import json
import os
from pathlib import Path
import signal
import subprocess
import sys
import time


def sha(data):
    return hashlib.sha256(data).hexdigest()


def main():
    job_path = Path(sys.argv[1]).resolve()
    job = json.loads(job_path.read_text())
    out = Path(sys.argv[2]).resolve()
    out.mkdir(parents=True, exist_ok=False)
    (out / "job.json").write_bytes(job_path.read_bytes())
    (out / "runner.py").write_bytes(Path(__file__).read_bytes())
    env = dict(os.environ)
    env.update(job.get("environment", {}))
    record = dict(job)
    record["job_sha256"] = sha(job_path.read_bytes())
    record["runner_sha256"] = sha(Path(__file__).read_bytes())
    record["started_utc"] = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
    log = out / "output.log"
    start = time.monotonic()
    with log.open("xb") as stream:
        process = subprocess.Popen(job["argv"], cwd=job["cwd"], env=env,
                                   stdout=stream, stderr=subprocess.STDOUT,
                                   start_new_session=True)
        record["pid"] = process.pid
        try:
            record["exit"] = process.wait(timeout=job["timeout_seconds"])
        except subprocess.TimeoutExpired:
            record["watchdog_expired"] = True
            os.killpg(process.pid, signal.SIGKILL)
            record["exit"] = process.wait()
    record["elapsed_seconds"] = time.monotonic() - start
    raw = log.read_bytes()
    record["output_bytes"] = len(raw)
    record["output_sha256"] = sha(raw)
    packed = gzip.compress(raw, mtime=0)
    (out / "output.log.gz").write_bytes(packed)
    record["compressed_output_sha256"] = sha(packed)
    if job.get("go_test_json"):
        events = []
        other = []
        for line in raw.decode(errors="replace").splitlines():
            try:
                event = json.loads(line)
            except json.JSONDecodeError:
                other.append(line)
                continue
            if isinstance(event, dict) and "Action" in event:
                events.append(event)
            else:
                other.append(line)
        record["test_summary"] = {
            "pass_events": sum(e.get("Action") == "pass" and "Test" in e for e in events),
            "package_results": [e for e in events if e.get("Action") in ("pass", "fail", "skip") and "Test" not in e],
            "named_skips": [e for e in events if e.get("Action") == "skip" and "Test" in e],
            "failures": [e for e in events if e.get("Action") == "fail"],
            "non_event_lines": other,
        }
    (out / "result.json").write_text(json.dumps(record, indent=2) + "\n")
    print(json.dumps({"exit": record["exit"], "elapsed_seconds": record["elapsed_seconds"],
                      "output": str(log), "record": str(out / "result.json")}), flush=True)
    return 0 if record["exit"] == job.get("expected_exit", 0) and not record.get("watchdog_expired") else 1


if __name__ == "__main__":
    sys.exit(main())
