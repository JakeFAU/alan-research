#!/usr/bin/env python3
"""Compile focused checker mutants through isolated Go overlays."""
import argparse
import gzip
import hashlib
import json
from pathlib import Path
import subprocess


def sha(data):
    return hashlib.sha256(data).hexdigest()


MUTANTS = {
    "collapse-revoke-root": {
        "file": "revocation.go",
        "changes": [
            ("l.Authority.Root == root.Concrete && l.Authority.RootID", "l.Authority.RootID"),
            ("revoke == (Revoke{bindingVersion, root.Concrete})", "revoke == (Revoke{bindingVersion, l.Authority.Root})"),
        ],
        "run": "^TestRevocationWrongTarget$",
        "failure": "semantic assertion: coherent c-root revoke accepted into b-only model: <nil>",
    },
    "omit-current-gate-revocation": {
        "file": "check.go",
        "changes": [
            (
                "w := c.deriveWitness(a.Intent)\n\tallowed, reason := w.Allowed, w.Reason",
                "savedRevocations := c.revoked\n\tc.revoked = map[AuthorityRoot]bool{}\n"
                "\tw := c.deriveWitness(a.Intent)\n\tc.revoked = savedRevocations\n"
                "\tallowed, reason := w.Allowed, w.Reason",
            ),
        ],
        "run": "^TestRevocationRawAttacks/fresh-allow-after-revoke$",
        "failure": "fresh-allow-after-revoke must reject at fresh_gate_policy, got support=false raw=true failure=<nil>",
    },
}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("scratch", type=Path)
    parser.add_argument("evidence", type=Path)
    args = parser.parse_args()
    args.scratch.mkdir(parents=True, exist_ok=True)
    args.evidence.mkdir(parents=True, exist_ok=True)
    result = {}
    for name, mutant in MUTANTS.items():
        d = args.scratch.resolve() / name
        d.mkdir(exist_ok=True)
        owner = Path("research/delegationcorrespondence") / mutant["file"]
        base = owner.read_text()
        content = base
        for old, replacement in mutant["changes"]:
            if content.count(old) != 1:
                raise ValueError(f"mutation anchor must occur exactly once for {name}")
            content = content.replace(old, replacement, 1)
        changed = d / owner.name
        changed.write_text(content)
        overlay = d / "overlay.json"
        overlay.write_text(json.dumps({"Replace": {str(owner.resolve()): str(changed)}}))
        exe = d / "checker.test"
        build = ["go", "test", "-overlay", str(overlay), "-c", "-o", str(exe), "./research/delegationcorrespondence"]
        assertion = [str(exe), "-test.v", "-test.run", mutant["run"], "-test.count=1"]
        recorded = {"source": str(owner), "base_sha256": sha(base.encode()), "mutant_sha256": sha(content.encode())}
        for stage, command in (("build", build), ("assertion", assertion)):
            process = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, check=False)
            log = process.stdout
            log_name = name + "-" + stage + ".log.gz"
            (args.evidence / log_name).write_bytes(gzip.compress(log, mtime=0))
            recorded[stage] = {"command": command, "exit": process.returncode, "log": log_name, "log_sha256": sha(log)}
            if stage == "build" and process.returncode != 0:
                raise ValueError(f"{name} failed to compile: {log.decode()}")
            if stage == "assertion" and (
                process.returncode != 1 or mutant["failure"].encode() not in log or b"panic:" in log
            ):
                raise ValueError(f"{name} missed intended assertion: {log.decode()}")
        result[name] = recorded
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
