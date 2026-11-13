"""Check closed qualification records, input hashes and ordinary package coverage."""
import gzip
import hashlib
import json
import os
from pathlib import Path
import subprocess


ROOT = Path.cwd()
OUT = ROOT / "docs/foundational/phase3b6/evidence/task4"


def digest(data):
    return hashlib.sha256(data).hexdigest()


def read(path):
    return json.loads(path.read_text())


def check_records(folder):
    records = read(folder / "summary.json")
    for record in records:
        name = record["name"]
        assert record == read(folder / (name + ".json")), name
        compressed = (folder / (name + ".log.gz")).read_bytes()
        raw = gzip.decompress(compressed)
        assert digest(compressed) == record["archive_sha256"], name
        assert digest(raw) == record["log_sha256"], name
        assert len(raw) == record["log_bytes"], name
        assert record["matched"] and record["source_unchanged"], name
        assert record["exit"] == record["expected_exit"], name
        if "test_summary" in record:
            events = [json.loads(line) for line in raw.splitlines() if line.startswith(b"{")]
            summary = record["test_summary"]
            assert not summary["failures"], name
            assert summary["pass_events"] == sum(
                e.get("Action") == "pass" and "Test" in e for e in events
            ), name
            assert summary["named_skips"] == [
                e for e in events if e.get("Action") == "skip" and "Test" in e
            ], name
            assert summary["package_results"] == [
                e for e in events
                if e.get("Action") in ("pass", "fail", "skip") and "Test" not in e
            ], name
    return {r["name"]: r for r in records}


old = check_records(OUT / "qualification-unaffected")
current = check_records(OUT / "qualification-final")
freeze = read(OUT / "freeze.json")
sources = read(OUT / "qualification-final/source.json")
assert len(sources) == freeze["source_files"]
assert digest((OUT / "qualification-final/source.json").read_bytes()) == freeze["source_manifest_sha256"]
for name, value in sources.items():
    assert digest((ROOT / name).read_bytes()) == value, name
old_sources = read(OUT / "qualification-unaffected/source.json")
assert all(sources[name] == value for name, value in old_sources.items())

env = dict(os.environ, **read(OUT / "qualification-final/plan.json")["environment"])
packages = subprocess.check_output(["go", "list", "./..."], env=env, cwd=ROOT).decode().splitlines()
groups = [old["ordinary-unaffected"], current["ordinary-affected"]]
observed = [[p["Package"] for p in r["test_summary"]["package_results"]] for r in groups]
assert not set(observed[0]) & set(observed[1])
assert set(observed[0]) | set(observed[1]) == set(packages) == set(freeze["ordinary_packages"])
assert all(len(group) == len(set(group)) for group in observed)
assert set(observed[0]) == set(freeze["ordinary_unaffected_packages"])
assert set(observed[1]) == set(freeze["ordinary_affected_packages"])

changed_go = subprocess.check_output([
    "git", "diff", "--name-only", "df840103aa3a85e33b4b95ddcf2fe3df071cb024", "--", "*.go"
], cwd=ROOT).decode().splitlines()
format_command = ["gofmt", "-l", *changed_go]
formatted = subprocess.run(format_command, cwd=ROOT, capture_output=True, check=True)
assert not formatted.stdout and not formatted.stderr

comparisons = []
for job, archive in [
    ("new-model", "docs/foundational/phase3b6/evidence/task1/graph.json.gz"),
    ("new-audit", "docs/foundational/phase3b6/evidence/task2/audit.json.gz"),
    ("prior-revocation-model", "docs/foundational/phase3b5/evidence/task1/model.json.gz"),
]:
    expected = digest(gzip.decompress((ROOT / archive).read_bytes()))
    assert current[job]["log_sha256"] == expected, job
    comparisons.append({"job": job, "archive": archive, "sha256": expected})
prior_progress = read(ROOT / "docs/foundational/phase3b5/evidence/task2/old-progress.json")
assert current["prior-default-progress"]["log_sha256"] == prior_progress["sha256"]
comparisons.append({"job": "prior-default-progress", "sha256": prior_progress["sha256"]})
prior_exports = {r["name"]: r for r in read(ROOT / "docs/foundational/phase3b5/evidence/task2/exports.json")}
for report in read(OUT / "exports/exports.json"):
    assert report["twice_byte_equal"] and report["sha256"] == prior_exports[report["name"]]["sha256"]
    comparisons.append({"job": report["name"], "sha256": report["sha256"]})

package_results = [p for r in groups for p in r["test_summary"]["package_results"]]
result = {
    "head_at_verification": subprocess.check_output(["git", "rev-parse", "HEAD"]).decode().strip(),
    "executable_freeze": freeze["head"],
    "verified_source_inputs": len(sources),
    "verified_unchanged_group_inputs": len(old_sources),
    "closed_command_records": len(old) + len(current),
    "ordinary_package_groups": [len(group) for group in observed],
    "ordinary_packages": len(packages),
    "ordinary_tested_packages": sum(p["Action"] == "pass" for p in package_results),
    "ordinary_no_test_packages": sum(p["Action"] == "skip" for p in package_results),
    "ordinary_pass_events": sum(r["test_summary"]["pass_events"] for r in groups),
    "ordinary_named_skips": [s for r in groups for s in r["test_summary"]["named_skips"]],
    "race": current["affected-race"]["test_summary"],
    "format_command": format_command,
    "format_exit": formatted.returncode,
    "format_output_sha256": digest(formatted.stdout),
    "matched_exports": comparisons,
    "verifier_sha256": digest(Path(__file__).read_bytes()),
    "scope": "Frozen executable inputs and closed grouped qualification; manuscript/PDF checks are separate. Archived recorder/verifier copies are evidence apparatus, not additional Go package inputs.",
}
(OUT / "verification.json").write_text(json.dumps(result, indent=2) + "\n")
print(json.dumps({k: result[k] for k in [
    "verified_source_inputs", "ordinary_packages", "ordinary_tested_packages",
    "ordinary_no_test_packages", "ordinary_pass_events", "closed_command_records"
]}))
