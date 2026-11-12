"""Compare default progress bytes with Task2 base source (no historical export exists).
Usage: python3 compare_old_progress.py SOURCE NEW_OUTPUT CURRENT_BINARY EVIDENCE_DIRECTORY
"""
import hashlib
import io
import json
import os
from pathlib import Path
import subprocess
import sys
import tarfile

source, output, binary, evidence = map(lambda p: Path(p).resolve(), sys.argv[1:])
output.mkdir(parents=True, exist_ok=False)
env = dict(os.environ, GOTOOLCHAIN="local", GOWORK="off", GOPROXY="off", GOSUMDB="off",
           GOMODCACHE="/tmp/alan-phase3b5/modcache", GOCACHE="/tmp/alan-phase3b5/go-build")
commands = []
def run(argv, cwd):
    r = subprocess.run(argv, cwd=cwd, env=env, capture_output=True, timeout=120)
    commands.append(dict(argv=argv, cwd=str(cwd), exit=r.returncode,
                         stdout_sha256=hashlib.sha256(r.stdout).hexdigest(), stderr=r.stderr.decode()))
    assert r.returncode == 0, (argv, r.stderr)
    return r.stdout
archive = run(["git", "archive", "e6e53deb", "go.mod", "go.sum", "research/authorityservice",
               "research/mutualobligations", "research/delegationcorrespondence"], source)
with tarfile.open(fileobj=io.BytesIO(archive)) as tar:
    for member in tar:
        path = output / member.name
        assert output in path.resolve().parents
        if member.isdir():
            path.mkdir(parents=True, exist_ok=True)
        else:
            assert member.isfile(), member.name
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_bytes(tar.extractfile(member).read())
run(["go", "build", "-buildvcs=false", "-o", str(output / "progress"),
     "./research/delegationcorrespondence/cmd/progress"], output)
old = run([str(output / "progress")], output)
new = run([str(binary)], source)
assert old == new
record = dict(base="e6e53deb", baseline_source_archive_sha256=hashlib.sha256(archive).hexdigest(),
              commands=commands, exact_bytes=True, sha256=hashlib.sha256(old).hexdigest(), bytes=len(old))
(evidence / "old-progress.json").write_text(json.dumps(record, indent=2) + "\n")
print("old progress exact bytes:", record["sha256"])
