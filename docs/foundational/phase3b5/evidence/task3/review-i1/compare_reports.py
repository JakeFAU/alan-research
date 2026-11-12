#!/usr/bin/env python3
"""Compare the Task 3 CLI twice and compare legacy CLI bytes with the task base."""
import argparse
import hashlib
import json
from pathlib import Path
import subprocess


def sha(data):
    return hashlib.sha256(data).hexdigest()


def command(bin_path, flag, path, expected):
    run = subprocess.run([str(bin_path), flag, str(path)], capture_output=True, check=False)
    if run.returncode != expected or run.stderr:
        raise ValueError(f"{bin_path.name} {path.name}: expected exit {expected}, got {run.returncode}: {run.stderr!r}")
    return run.stdout


def compare(root, baseline, legacy, checker):
    manifest = json.loads((root / "manifest.json").read_bytes())
    results = []
    for mode in ("normal", "race"):
        cases = manifest[f"raw-{mode}.tar.gz"]["cases"]
        for case in ("N05", "N06", "N07", "P02-v2", "applied-unknown-revoked-before-fresh-gate"):
            path = root / f"raw-{mode}" / cases[case]
            supported = not case.startswith("applied-unknown")
            first = command(checker, "-case", path, 0 if supported else 1)
            second = command(checker, "-case", path, 0 if supported else 1)
            before = command(baseline, "-binding-case", path, 0)
            after = command(legacy, "-binding-case", path, 0)
            report = json.loads(first)
            if first != second or before != after or json.loads(before) != report["Legacy"]:
                raise ValueError(f"determinism or legacy report changed for {mode}/{case}")
            physical = report["Physical"]
            producer = report["Producer"]
            if (
                report["Supported"] != supported or not report["RawValid"] or not physical["Valid"]
                or report["Legacy"]["ContinuityFailure"] is not None
                or producer["Race"] != (mode == "race")
                or producer["Commit"] != "5b4662fd682a220ca3065b15e86672b7e56567e7"
                or producer["GoVersion"] != "go1.27.3"
            ):
                raise ValueError(f"verdict or producer mismatch for {mode}/{case}")
            cuts = []
            for cut in report["Captures"]:
                if cut["Supported"] != supported:
                    raise ValueError(f"unexpected cut verdict for {mode}/{case}/{cut['Path']}")
                cuts.append({
                    "path": cut["Path"],
                    "sha256": cut["SHA256"],
                    "supported": cut["Supported"],
                    "event_count": len(cut["Events"] or []),
                    "unresolved": ["AB" if d == 0 else "BA" for d in (cut["Unresolved"] or [])],
                })
            results.append({
                "mode": mode, "case": case, "directory": cases[case],
                "legacy_sha256": sha(before), "report_sha256": sha(first), "report_bytes": len(first),
                "legacy_supported": report["Legacy"]["Projection"]["Supported"], "supported": supported,
                "reasons": report["Reasons"], "raw_valid": report["RawValid"], "physical_valid": physical["Valid"],
                "contacts": physical["Contacts"], "applications": physical["Applications"],
                "acknowledged_losses": physical["AcknowledgedLosses"], "producer": producer, "captures": cuts,
            })
    return results


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    for name in ("extraction", "legacy_base", "legacy_current", "revocation"):
        parser.add_argument(name, type=Path)
    args = parser.parse_args()
    print(json.dumps(compare(args.extraction, args.legacy_base, args.legacy_current, args.revocation), indent=2))


if __name__ == "__main__":
    main()
