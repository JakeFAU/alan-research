#!/usr/bin/env python3
"""Compile disposable ledger mutants and require targeted semantic failures.

Run from repository root: PYTHONPYCACHEPREFIX=/tmp/alan-phase3b6/python-cache
python3 docs/foundational/phase3b6/evidence/task2/run-mutants.py
No repository executable is edited; copies and binaries stay under /tmp.
"""
import difflib
import gzip
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile

root = Path.cwd()
evidence = root / "docs/foundational/phase3b6/evidence/task2"
source = root / "research/effectobservation/oracle.go"
original = source.read_text()
env = dict(os.environ, GOTOOLCHAIN="local", GOWORK="off", GOPROXY="off", GOSUMDB="off",
           GOMODCACHE="/tmp/alan-phase3b5/modcache", GOCACHE="/tmp/alan-phase3b6/go-build",
           STATICCHECK_CACHE="/tmp/alan-phase3b6/staticcheck",
           PYTHONPYCACHEPREFIX="/tmp/alan-phase3b6/python-cache")
mutants = [
    ("packet-deletion-on-loss", "l.Dead, l.Opportunity, l.Inbox = true, false, NoResult",
     "l.Dead, l.Opportunity, l.Inbox = true, false, NoResult\n\t\tl.PacketSent = false",
     "^TestLedgerLawfulControls$/^packet-survives-loss$"),
    ("application-as-retention", "l.World = e.Value", "l.World = e.Value\n\t\tl.Stored = e.Value",
     "^TestLedgerLawfulControls$/^application-is-not-retention$"),
    ("old-response-after-reopen", "if !original || l.World == NoResult", "if l.Dead || l.World == NoResult",
     "^TestLedgerRejectsInvalidCustody$/^revived-old-response$"),
    ("c-revocation", "l.Revoked = true", "l.Revoked = true\n\t\tl.CValid = false",
     "^TestLedgerLawfulControls$/^old-permit-survives-revoke$"),
]
manifest = {"source_sha256": hashlib.sha256(source.read_bytes()).hexdigest(), "mutants": []}
with tempfile.TemporaryDirectory(prefix="task2-mutants-", dir="/tmp/alan-phase3b6") as tmp:
    workspace = Path(tmp)
    shutil.copy2(root / "go.mod", workspace / "go.mod")
    shutil.copy2(root / "go.sum", workspace / "go.sum")
    shutil.copytree(root / "research/effectobservation", workspace / "research/effectobservation")
    target = workspace / "research/effectobservation/oracle.go"
    binary = workspace / "ledger.test"
    compile_cmd = ["go", "test", "-c", "-o", str(binary), "./research/effectobservation"]
    baseline = subprocess.run(compile_cmd, cwd=workspace, env=env, capture_output=True, text=True)
    assert baseline.returncode == 0, baseline.stdout + baseline.stderr
    for name, old, new, test in mutants:
        lawful = subprocess.run([str(binary), "-test.v", "-test.run", test], cwd=workspace,
                                env=env, capture_output=True, text=True)
        assert lawful.returncode == 0, lawful.stdout + lawful.stderr
    for name, old, new, test in mutants:
        assert original.count(old) == 1, (name, original.count(old))
        mutated = original.replace(old, new, 1)
        target.write_text(mutated)
        built = subprocess.run(compile_cmd, cwd=workspace, env=env, capture_output=True, text=True)
        assert built.returncode == 0, name + built.stdout + built.stderr
        result = subprocess.run([str(binary), "-test.v", "-test.run", test], cwd=workspace,
                                env=env, capture_output=True, text=True)
        output = result.stdout + result.stderr
        assert result.returncode == 1 and "--- FAIL:" in output and "panic:" not in output, name + output
        log = (f"name: {name}\ncompile: go test -c -o <disposable>/ledger.test ./research/effectobservation\n"
               f"compile_exit: {built.returncode}\n{built.stdout}{built.stderr}"
               f"run: <disposable>/ledger.test -test.v -test.run '{test}'\n"
               f"run_exit: {result.returncode}\n{output}")
        (evidence / f"mutant-{name}.log.gz").write_bytes(gzip.compress(log.encode(), mtime=0))
        patch = "".join(difflib.unified_diff(original.splitlines(True), mutated.splitlines(True),
                       fromfile="a/research/effectobservation/oracle.go", tofile="b/research/effectobservation/oracle.go", n=0))
        (evidence / f"mutant-{name}.patch").write_text(patch)
        manifest["mutants"].append({"name": name, "compile_exit": built.returncode,
                                   "targeted_test": test, "baseline_test_exit": 0,
                                   "mutant_test_exit": result.returncode,
                                   "mutant_source_sha256": hashlib.sha256(mutated.encode()).hexdigest()})
        print(log)
(evidence / "mutants.json").write_text(json.dumps(manifest, indent=2) + "\n")
