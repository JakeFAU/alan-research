#!/usr/bin/env python3
"""Verify and optionally extract only the declared inputs for Task 3."""
import argparse
import hashlib
import json
from pathlib import Path, PurePosixPath
import tarfile

CASES = ("N05", "N06", "N07", "P02-v2", "applied-unknown-revoked-before-fresh-gate")


def sha(data):
    return hashlib.sha256(data).hexdigest()


def verify(root, extraction, prepare):
    sealed = root / "docs/foundational/phase3b4/evidence/prospective-main"
    archives = json.loads((sealed / "archives.json").read_bytes())
    selected = json.loads(
        (root / "docs/foundational/phase3b5/evidence/preflight/corpus-extraction.json").read_bytes()
    )
    summary = {}
    for mode in ("normal", "race"):
        name = f"raw-{mode}.tar.gz"
        path = sealed / name
        manifest = archives[name]
        archive_sha = sha(path.read_bytes())
        if archive_sha != manifest["sha256"] or path.stat().st_size != manifest["bytes"]:
            raise ValueError(f"archive mismatch {name}")
        cases = {case: selected[name]["cases"][case] for case in CASES}
        output = {"archive_sha256": archive_sha, "cases": cases, "directory": str(extraction / f"raw-{mode}"), "verified": {}}
        with tarfile.open(path, "r:gz") as archive:
            members = {}
            for member in archive:
                if member.name in members:
                    raise ValueError(f"duplicate member {member.name}")
                members[member.name] = member
            for case, prefix in cases.items():
                checked = {}
                for key, meta in sorted(manifest["members"].items()):
                    if not key.startswith(prefix + "/"):
                        continue
                    rel = PurePosixPath(key)
                    if rel.is_absolute() or ".." in rel.parts or str(rel) != key:
                        raise ValueError(f"unsafe path {key}")
                    member = members[key]
                    if not member.isfile():
                        raise ValueError(f"non-regular member {key}")
                    with archive.extractfile(member) as stream:
                        data = stream.read()
                    if sha(data) != meta["sha256"] or len(data) != meta["bytes"]:
                        raise ValueError(f"member mismatch {key}")
                    checked[str(rel.relative_to(prefix))] = data
                digests = json.loads(checked["digests.json"])
                case_json = json.loads(checked["case.json"])
                if case_json["Name"] != case:
                    raise ValueError(f"case identity mismatch {case}")
                required = {"case.json", "endpoint-a.jsonl", "endpoint-b.jsonl", "endpoints.json"}
                required.update(capture["Directory"] + "/inspect.json" for capture in case_json["Captures"])
                if not required <= digests.keys():
                    raise ValueError(f"missing declared inputs {case}: {required - digests.keys()}")
                for key, digest in digests.items():
                    if sha(checked[key]) != digest:
                        raise ValueError(f"case digest mismatch {case}/{key}")
                for key in sorted({"digests.json", *digests}):
                    dest = extraction / f"raw-{mode}" / prefix / key
                    if prepare:
                        dest.parent.mkdir(parents=True, exist_ok=True)
                        dest.write_bytes(checked[key])
                    if dest.read_bytes() != checked[key]:
                        raise ValueError(f"extracted member mismatch {dest}")
                output["verified"][case] = {
                    "archive_members": len(checked),
                    "declared_members": len(digests),
                    "bytes": sum(len(data) for data in checked.values()),
                    "case_sha256": digests["case.json"],
                    "digests_sha256": sha(checked["digests.json"]),
                }
        summary[name] = output
    if prepare:
        (extraction / "manifest.json").write_text(json.dumps(summary, indent=2) + "\n")
    return summary


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("extraction", type=Path)
    parser.add_argument("--prepare", action="store_true", help="populate a new disposable extraction root")
    args = parser.parse_args()
    if args.prepare and args.extraction.exists():
        parser.error("--prepare requires a new extraction root")
    print(json.dumps(verify(Path.cwd(), args.extraction.resolve(), args.prepare), indent=2))


if __name__ == "__main__":
    main()
