"""Verify the retained Phase 3B.7 package without opening SQLite or rerunning it."""
import argparse
import gzip
import hashlib
import json
from pathlib import Path
import subprocess
import tarfile


def sha(data):
    return hashlib.sha256(data).hexdigest()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("root", type=Path)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    root = args.root.resolve(strict=True)
    evidence = root / "docs/foundational/phase3b7/evidence"
    qualification = evidence / "qualification"
    subprocess.run([
        "python3", str(qualification / "freeze-inputs.py"), str(root),
        "--check", str(qualification / "freeze-corrected.json"),
    ], check=True)
    result = {"head": subprocess.check_output(
        ["git", "rev-parse", "HEAD"], cwd=root, text=True).strip(),
        "verifier_sha256": sha(Path(__file__).read_bytes()),
        "frozen_inputs_unchanged": True, "archives": {}, "commands": {}}
    summary = json.loads((qualification / "corrected-checks.json").read_text())
    initial = json.loads((qualification / "initial-checks.json").read_text())
    for mode in ("normal", "race"):
        directory = evidence / "capture" / (mode + "-v1")
        record = json.loads((directory / "archive.json").read_text())
        packed = (directory / "raw.tar.gz").read_bytes()
        assert len(packed) == record["archive_bytes"]
        assert sha(packed) == record["archive_sha256"]
        members, inventory = {}, {}
        with tarfile.open(directory / "raw.tar.gz", "r:gz") as archive:
            for member in archive:
                assert member.isfile() and member.name not in members
                data = archive.extractfile(member).read()
                members[member.name] = data
                inventory[member.name] = {
                    "bytes": len(data), "sha256": sha(data), "mode": member.mode,
                }
        assert inventory == record["files"] and len(members) == 122
        for cut in ("first-owner/post-open", "stopped", "reopened", "final"):
            assert members[f"zero/{cut}/inspect.json"] == members[f"applied/{cut}/inspect.json"]
        assert members["zero/run.json"] == members["applied/run.json"]
        final = json.loads(members["zero/final/inspect.json"])
        assert len(final["Tables"]) == 41
        assert sum(map(len, final["Tables"].values())) == 94
        assert summary[mode]["Supported"] and not initial[mode]["Supported"]
        for stage, name, expected in (("initial", "v1", False), ("corrected", "corrected", True)):
            outputs = []
            for run in (1, 2):
                path = qualification / (stage + "-commands") / f"check-{mode}-{name}-{run}"
                data = gzip.decompress((path / "output.log.gz").read_bytes())
                report = json.loads(data)
                assert report["Supported"] is expected
                assert all(inventory[key]["sha256"] == value["sha256"]
                           and inventory[key]["bytes"] == value["bytes"]
                           for key, value in report["Inputs"].items())
                outputs.append(data)
            assert outputs[0] == outputs[1]
            expected_hash = (summary if expected else initial)[mode]["report_sha256"]
            assert sha(outputs[0]) == expected_hash
        result["archives"][mode] = {"sha256": sha(packed), "members": len(members),
            "four_exports_and_final_report_byte_equal": True, "final_tables": 41,
            "final_rows": 94, "initial_rejections_and_corrected_reports_preserved": True}
    for group in ("initial-commands", "corrected-commands", "verification-commands"):
        directory = qualification / group
        index = json.loads((directory / "index.json").read_text())
        for name, expected in index.items():
            command = directory / name
            for filename, digest in expected["files"].items():
                assert sha((command / filename).read_bytes()) == digest
            record = json.loads((command / "result.json").read_text())
            output = gzip.decompress((command / "output.log.gz").read_bytes())
            assert len(output) == record["output_bytes"]
            assert sha(output) == record["output_sha256"]
            assert sha((command / "output.log.gz").read_bytes()) == record["compressed_output_sha256"]
            assert sha((command / "job.json").read_bytes()) == record["job_sha256"]
            assert record["runner_sha256"] == sha((qualification / "run-command.py").read_bytes())
            assert record["exit"] == expected["exit"]
        result["commands"][group] = {"records": len(index), "hashes_and_outputs_valid": True}
    result["limitations"] = [
        "This verifies retained package consistency, not artifact authenticity or a new experiment.",
        "SQLite members are opaque bytes; no database was opened.",
        "Full protocol/world/model judgments remain in the separately reviewed checker reports.",
    ]
    with args.output.open("x") as output:
        json.dump(result, output, indent=2, sort_keys=True)
        output.write("\n")
    print(json.dumps({"verified": True, "output": str(args.output)}))


if __name__ == "__main__":
    main()
