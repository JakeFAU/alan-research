#!/usr/bin/env python3
"""Reproduce Task 1 qualification from repository root; retain exact diagnostics."""
import gzip
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys

DEST = Path(__file__).resolve().parent
ROOT = DEST.parents[4]
ENV = {"GOCACHE": "/tmp/alan-phase3c4/go-cache", "STATICCHECK_CACHE": "/tmp/alan-phase3c4/staticcheck-cache"}

def digest(data):
    return hashlib.sha256(data).hexdigest()

def write_json(name, value):
    (DEST / name).write_text(json.dumps(value, indent=2) + "\n")

def run(name, argv, large=False, expected_exit=0):
    result = subprocess.run(argv, cwd=ROOT, env=os.environ | ENV, capture_output=True)
    record = {"argv": argv, "cwd": "repository root", "env": ENV,
              "exit": result.returncode, "stderr": result.stderr.decode()}
    if large:
        blob = gzip.compress(result.stdout, mtime=0)
        (DEST / (name + ".json.gz")).write_bytes(blob)
        record.update(stdout_file=name + ".json.gz", stdout_bytes=len(result.stdout),
                      stdout_sha256=digest(result.stdout), gzip_sha256=digest(blob))
    else:
        record["stdout"] = result.stdout.decode()
    write_json(name + "-execution.json", record)
    print(name, "exit", result.returncode, flush=True)
    if result.returncode != expected_exit:
        raise SystemExit(f"unexpected exit {result.returncode}, expected {expected_exit}")
    return result.stdout

if __name__ == "__main__":
    selection = ["-run", "TestCompleteEvidenceAndWitnesses"] if "--affected" in sys.argv else []
    run("ordinary", ["go", "test", "-v", "-timeout", "30m"] + selection + ["./research/residualinterfaces/..."])
    run("race", ["go", "test", "-race", "-v", "-timeout", "30m"] + selection + ["./research/residualinterfaces/..."])
    run("staticcheck", ["staticcheck", "./research/residualinterfaces/..."])
    deps_raw = run("go-list-deps", ["go", "list", "-deps", "-json", "./research/residualinterfaces/cmd/export"], large=True)
    first = run("export", ["go", "run", "./research/residualinterfaces/cmd/export"], large=True)
    second = run("export-second", ["go", "run", "./research/residualinterfaces/cmd/export"], large=True)
    if first != second:
        raise SystemExit("two captures differ")
    payload = json.loads(first)
    decoder = json.JSONDecoder()
    text = deps_raw.decode()
    compiled = []
    while text.strip():
        obj, end = decoder.raw_decode(text.lstrip())
        text = text.lstrip()[end:]
        directory = Path(obj.get("Dir", "/"))
        if directory.is_relative_to(ROOT):
            compiled.extend(str((directory / name).relative_to(ROOT)) for name in obj.get("GoFiles", []))
    expected = sorted(row["Path"] for row in payload["Sources"] if row["Use"] == "compiled producer")
    if sorted(compiled) != expected:
        raise SystemExit("compiled manifest roster differs from go list")
    write_json("compiled-inputs.json", {"argv": ["go", "list", "-deps", "-json", "./research/residualinterfaces/cmd/export"],
        "Compiled": sorted(compiled), "ManifestCompiled": expected, "Equal": True})
    write_json("artifact-manifest.json", {"Schema": payload["Schema"], "RawBytes": len(first),
        "RawSHA256": digest(first), "GzipBytes": (DEST / "export.json.gz").stat().st_size,
        "GzipSHA256": digest((DEST / "export.json.gz").read_bytes()),
        "SecondRawSHA256": digest(second), "Identical": first == second,
        "SourceManifestSHA256": digest(json.dumps(payload["Sources"], separators=(",", ":")).encode()),
        "Counts": payload["Counts"]})
    write_json("sources.json", payload["Sources"])
    run("falsification-union", ["python3", "docs/foundational/phase3c4/evidence/task1/falsify.py", "union-sufficiency"], expected_exit=1)
    run("falsification-inevitable", ["python3", "docs/foundational/phase3c4/evidence/task1/falsify.py", "inevitable-completion"], expected_exit=1)
