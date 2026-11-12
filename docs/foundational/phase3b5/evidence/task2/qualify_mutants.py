"""Compile four independent-ledger mutants in disposable copies.
Usage: python3 qualify_mutants.py SOURCE NEW_OUTPUT EVIDENCE_DIRECTORY
A kill requires build exit 0 and targeted assertion exit 1, without a panic.
"""
import gzip
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys

source, output, evidence = map(lambda p: Path(p).resolve(), sys.argv[1:])
output.mkdir(parents=True, exist_ok=False)
evidence.mkdir(parents=True, exist_ok=True)
env = dict(os.environ, GOTOOLCHAIN="local", GOWORK="off", GOPROXY="off", GOSUMDB="off",
           GOMODCACHE="/tmp/alan-phase3b5/modcache", GOCACHE="/tmp/alan-phase3b5/go-build",
           STATICCHECK_CACHE="/tmp/alan-phase3b5/staticcheck")
mutants = [
    ("omit-fresh-root-check", 'if !live {\n\t\t\t\treturn fmt.Errorf("fresh-admission-revoked")',
     'if false && !live {\n\t\t\t\treturn fmt.Errorf("fresh-admission-revoked")',
     "TestRevocationOracleFreshAdmission", "fresh revoked admission accepted"),
    ("revoke-c-with-b", 'id == "B/g1:b" && l.revoked',
     '(id == "B/g1:b" || id == "B/g1:c") && l.revoked',
     "TestRevocationOracleDistinctRoot", "distinct c root must survive b revocation"),
    ("accept-post-revocation-issuance", '_, err := probe.grantPath(g.ID, e.Binding, true)',
     '_, err := probe.grantPath(g.ID, e.Binding, false)',
     "TestRevocationOracleRejectedIssuance", "accepted post-revocation issuance was not rejected"),
    ("retain-crash-permit", 'l.originals[1-e.Actor].permit = false\n\t\tl.originals[1-e.Actor].permitEpoch = 0',
     '// mutant: retain the activation-bound permit across crash',
     "TestRevocationOracleCrashPermit", "crashed permit must clear before fresh refusal"),
]
paths = [source / "go.mod", source / "go.sum"]
for package in ("authorityservice", "mutualobligations"):
    paths.extend(sorted((source / "research" / package).glob("*.go")))
records = []
for name, before, after, test, assertion in mutants:
    tree = output / name / "source"
    tree.mkdir(parents=True)
    for path in paths:
        target = tree / path.relative_to(source)
        target.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(path, target)
    target = tree / "research/authorityservice/oracle.go"
    text = target.read_text()
    assert text.count(before) == 1, (name, text.count(before))
    target.write_text(text.replace(before, after))
    record = dict(name=name, before=before, after=after, test=test, assertion=assertion,
                  source_sha256={str(p.relative_to(tree)): hashlib.sha256(p.read_bytes()).hexdigest()
                                 for p in sorted(tree.rglob("*")) if p.is_file()}, commands=[])
    binary = output / name / "checker.test"
    def run(argv, label):
        result = subprocess.run(argv, cwd=tree, env=env, capture_output=True, timeout=180)
        data = result.stdout + result.stderr
        log = name + "-" + label + ".log.gz"
        (evidence / log).write_bytes(gzip.compress(data, mtime=0))
        record["commands"].append(dict(argv=argv, cwd=str(tree), exit=result.returncode,
                                        log=log, sha256=hashlib.sha256(data).hexdigest()))
        return result.returncode, data.decode()
    build, _ = run(["go", "test", "-c", "-o", str(binary), "./research/authorityservice"], "build")
    code, log = run([str(binary), "-test.run", "^" + test + "$", "-test.count=1", "-test.v"], "assertion") if build == 0 else (-1, "")
    record["qualified_kill"] = build == 0 and code == 1 and assertion in log and "panic:" not in log
    records.append(record)
    print(name, "build", build, "assertion", code, "qualified", record["qualified_kill"], flush=True)
(evidence / "mutants.json").write_text(json.dumps(records, indent=2) + "\n")
sys.exit(0 if all(r["qualified_kill"] for r in records) else 1)
