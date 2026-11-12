"""Reproduce new exports twice and compare all four historical archived bytes.
Usage: python3 qualify_exports.py SOURCE NEW_OUTPUT EVIDENCE_DIRECTORY
"""
import gzip
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys

source, output, evidence = map(lambda p: Path(p).resolve(), sys.argv[1:])
output.mkdir(parents=True, exist_ok=False)
env = dict(os.environ, GOTOOLCHAIN="local", GOWORK="off", GOPROXY="off", GOSUMDB="off",
           GOMODCACHE="/tmp/alan-phase3b5/modcache", GOCACHE="/tmp/alan-phase3b5/go-build",
           STATICCHECK_CACHE="/tmp/alan-phase3b5/staticcheck")
configs = [
    ("twoagent", "research/twoagent/cmd/export", "docs/foundational/phase3a/evidence/model.json", []),
    ("mutual", "research/mutualobligations/cmd/export", "docs/foundational/phase3b/evidence/model.json", []),
    ("authority", "research/authorityservice/cmd/export", "docs/foundational/phase3b2/evidence/model.json.gz", []),
    ("audit", "research/authorityservice/cmd/audit", "docs/foundational/phase3b2/evidence/audit.json.gz", []),
    ("revocation-audit", "research/authorityservice/cmd/revocation", None, []),
    ("revocation-progress", "research/delegationcorrespondence/cmd/progress", None, ["-config", "S2/SeparateGrantRevocation"]),
]
records = []
for name, package, archive, args in configs:
    binary = output / name
    record = dict(name=name, commands=[])
    def run(argv, label):
        result = subprocess.run(argv, cwd=source, env=env, capture_output=True, timeout=300)
        record["commands"].append(dict(argv=argv, cwd=str(source), exit=result.returncode,
                                      stdout_sha256=hashlib.sha256(result.stdout).hexdigest(),
                                      stderr=result.stderr.decode()))
        (output / (label + ".stdout")).write_bytes(result.stdout)
        assert result.returncode == 0, (label, result.stderr.decode())
        return result.stdout
    run(["go", "build", "-buildvcs=false", "-o", str(binary), "./" + package], name + "-build")
    first = run([str(binary), *args], name + "-first")
    second = run([str(binary), *args], name + "-second")
    assert first == second, name
    record["twice_byte_equal"] = True
    record["bytes"] = len(first)
    record["sha256"] = hashlib.sha256(first).hexdigest()
    if archive:
        frozen = (source / archive).read_bytes()
        decoded = gzip.decompress(frozen) if archive.endswith(".gz") else frozen
        assert first == decoded, (name, "archive bytes differ")
        record.update(archive=archive, archived_sha256=hashlib.sha256(frozen).hexdigest(),
                      decoded_sha256=hashlib.sha256(decoded).hexdigest(), exact_archived_bytes=True)
    else:
        (evidence / (name + ".json.gz")).write_bytes(gzip.compress(first, mtime=0))
    records.append(record)
    print(name, "identical twice", record.get("exact_archived_bytes", "new"), record["sha256"], flush=True)
(evidence / "exports.json").write_text(json.dumps(records, indent=2) + "\n")
