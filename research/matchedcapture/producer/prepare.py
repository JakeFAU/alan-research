#!/usr/bin/env python3
"""Materialize a pinned ALAN runtime and verify every file before overlaying tests."""

import argparse
import hashlib
import json
from pathlib import Path
import stat
import subprocess


COMMIT = "5b4662fd682a220ca3065b15e86672b7e56567e7"
TREE = "44af03133ac7ba75fb6f6f8dc36629ed4cf52d98"
VERSION = "alan/phase3b7-matched/v1"
OVERLAYS = (
    "phase3b7_controls_test.go",
    "phase3b7_matched_test.go",
    "phase3b7_owner_test.go",
    "phase3b7_process_test.go",
)


def git(repo, *args):
    return subprocess.check_output(
        ["git", "-C", str(repo), "-c", "core.hooksPath=/dev/null", *args]
    )


def digest(raw):
    return {"bytes": len(raw), "sha256": hashlib.sha256(raw).hexdigest()}


def prepare(source, destination):
    source = source.resolve(strict=True)
    destination = destination.resolve()
    checkout = Path(git(source, "rev-parse", "--absolute-git-dir").decode().strip())
    common = Path(
        git(source, "rev-parse", "--path-format=absolute", "--git-common-dir")
        .decode().strip()
    )
    if any(destination == p or destination.is_relative_to(p) for p in (source, checkout, common)):
        raise ValueError("destination must be outside the source repository")
    if destination.exists():
        raise ValueError("destination must be fresh")
    if git(source, "rev-parse", f"{COMMIT}^{{tree}}").decode().strip() != TREE:
        raise ValueError("frozen runtime tree mismatch")
    # --no-hardlinks ensures checkout object custody is independent.
    git(source, "clone", "--local", "--no-hardlinks", "--no-checkout", "--",
        str(source), str(destination))
    git(destination, "-c", "core.autocrlf=false", "checkout", "--detach", COMMIT)
    if git(destination, "rev-parse", "HEAD", "HEAD^{tree}").splitlines() != [
        COMMIT.encode(), TREE.encode()
    ]:
        raise ValueError("prepared HEAD or tree mismatch")
    if git(destination, "status", "--porcelain", "--untracked-files=all"):
        raise ValueError("unexpected pre-overlay checkout changes")

    files = {}
    for record in git(destination, "ls-tree", "-rz", COMMIT).split(b"\x00"):
        if not record:
            continue
        header, path = record.split(b"\t", 1)
        mode, kind, oid = header.split()
        name = path.decode()
        local = destination / name
        actual_mode = local.lstat().st_mode
        if kind != b"blob" or not stat.S_ISREG(actual_mode):
            raise ValueError(f"unsupported source member {name}")
        if mode != (b"100755" if actual_mode & stat.S_IXUSR else b"100644"):
            raise ValueError(f"frozen source mode mismatch in {name}")
        raw = local.read_bytes()
        blob = f"blob {len(raw)}\0".encode() + raw
        if hashlib.sha1(blob).hexdigest().encode() != oid:
            raise ValueError(f"frozen source bytes mismatch in {name}")
        files[name] = {"mode": mode.decode(), "git_blob": oid.decode(), **digest(raw)}

    producer = Path(__file__).resolve().parent
    overlays = {}
    for name in OVERLAYS:
        overlay = producer / (name + ".txt")
        raw = overlay.read_bytes()
        local = f"internal/durabletest/{name}"
        with (destination / local).open("xb") as output:
            output.write(raw)
        overlays[local] = {"source": overlay.name, **digest(raw)}

    manifest = destination / "phase3b7-producer-source.json"
    with manifest.open("x", encoding="utf-8") as output:
        json.dump(
            {
                "version": VERSION,
                "runtime_commit": COMMIT,
                "runtime_tree": TREE,
                "runtime_files": files,
                "overlays": overlays,
                "preparation_script": digest(Path(__file__).read_bytes()),
            },
            output,
            indent=2,
            sort_keys=True,
        )
        output.write("\n")
    print(manifest)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("runtime_repository", type=Path)
    parser.add_argument("fresh_destination", type=Path)
    args = parser.parse_args()
    prepare(args.runtime_repository, args.fresh_destination)


if __name__ == "__main__":
    main()
