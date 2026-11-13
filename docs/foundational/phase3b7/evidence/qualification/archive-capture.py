"""Seal an already stopped raw capture without changing its members.

SQLite files are opaque bytes. Tar metadata is made deterministic; the raw
member bytes and their relative names remain exactly as captured.
"""
import gzip
import hashlib
import json
from pathlib import Path
import stat
import sys
import tarfile


def inventory(root):
    result = {}
    for path in sorted(root.rglob("*")):
        info = path.lstat()
        if stat.S_ISDIR(info.st_mode):
            continue
        if not stat.S_ISREG(info.st_mode):
            raise ValueError("non-regular capture member: " + str(path))
        raw = path.read_bytes()
        result[path.relative_to(root).as_posix()] = {
            "bytes": len(raw), "sha256": hashlib.sha256(raw).hexdigest(),
            "mode": stat.S_IMODE(info.st_mode),
        }
    return result


def main():
    raw_root = Path(sys.argv[1]).resolve(strict=True)
    output = Path(sys.argv[2]).resolve()
    if not raw_root.is_dir() or output == raw_root or output.is_relative_to(raw_root):
        raise ValueError("output must be outside the raw capture directory")
    output.mkdir(parents=True, exist_ok=False)
    before = inventory(raw_root)
    archive = output / "raw.tar.gz"
    with archive.open("xb") as file:
        with gzip.GzipFile(filename="", mode="wb", fileobj=file, mtime=0) as packed:
            with tarfile.open(fileobj=packed, mode="w|") as tar:
                for name, record in before.items():
                    info = tarfile.TarInfo(name)
                    info.size = record["bytes"]
                    info.mode = record["mode"]
                    info.mtime = 0
                    with (raw_root / name).open("rb") as source:
                        tar.addfile(info, source)
    if inventory(raw_root) != before:
        raise RuntimeError("raw capture changed while archiving; preserve this failed archive")
    with tarfile.open(archive, "r:gz") as tar:
        actual = {}
        for member in tar:
            if not member.isfile() or member.name in actual:
                raise ValueError("unexpected archive member")
            stream = tar.extractfile(member)
            if stream is None:
                raise ValueError("archive member missing stream")
            raw = stream.read()
            actual[member.name] = {
                "bytes": len(raw), "sha256": hashlib.sha256(raw).hexdigest(),
                "mode": member.mode,
            }
        if actual != before:
            raise RuntimeError("archive member mismatch")
    raw = archive.read_bytes()
    record = {
        "raw_directory": str(raw_root), "archive": archive.name,
        "archive_bytes": len(raw), "archive_sha256": hashlib.sha256(raw).hexdigest(),
        "files": before, "raw_unchanged": True, "archive_members_equal": True,
        "producer_member_seal_present": "members.json" in before,
        "archiver_sha256": hashlib.sha256(Path(__file__).read_bytes()).hexdigest(),
    }
    (output / "archive.json").write_text(json.dumps(record, indent=2, sort_keys=True) + "\n")
    print(json.dumps({k: v for k, v in record.items() if k != "files"}))


if __name__ == "__main__":
    main()
