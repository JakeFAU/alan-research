#!/usr/bin/env python3
"""Read-only checks for the distinct resolved dependency record."""
import gzip
import hashlib
import json
from pathlib import Path

root = Path.cwd()
dest = root / "docs/foundational/phase3c4/evidence/task1"
record = json.loads((dest / "go-list-deps-resolved-execution.json").read_text())
blob = (dest / record["stdout_file"]).read_bytes()
raw = gzip.decompress(blob)
assert record["exit"] == 0 and record["stderr"] == ""
assert hashlib.sha256(blob).hexdigest() == record["gzip_sha256"]
assert hashlib.sha256(raw).hexdigest() == record["stdout_sha256"]
assert len(raw) == record["stdout_bytes"]
text = raw.decode()
decoder = json.JSONDecoder()
packages, compiled = [], []
while text.strip():
    package, end = decoder.raw_decode(text.lstrip())
    text = text.lstrip()[end:]
    assert not package.get("Error") and not package.get("DepsErrors"), package
    packages.append(package["ImportPath"])
    directory = Path(package.get("Dir", "/"))
    if directory.is_relative_to(root):
        compiled.extend(str((directory / name).relative_to(root)) for name in package.get("GoFiles", []))
sources = json.loads((dest / "sources.json").read_text())
expected = sorted(row["Path"] for row in sources if row["Use"] == "compiled producer")
assert len(compiled) == 9 and sorted(compiled) == expected
for row in sources:
    assert hashlib.sha256((root / row["Path"]).read_bytes()).hexdigest() == row["SHA256"], row["Path"]
originals = {}
for prefix in ("", "pre-strictness-correction/"):
    old_record_path = dest / (prefix + "go-list-deps-execution.json")
    old_blob_path = dest / (prefix + "go-list-deps.json.gz")
    old = json.loads(old_record_path.read_text())
    assert old["exit"] == 0 and "writing stat cache" in old["stderr"] and "operation not permitted" in old["stderr"]
    assert gzip.decompress(old_blob_path.read_bytes()) == raw
    for path in (old_record_path, old_blob_path):
        originals[str(path.relative_to(root))] = hashlib.sha256(path.read_bytes()).hexdigest()
manifest = json.loads((dest / "artifact-manifest.json").read_text())
for name in ("export.json.gz", "export-second.json.gz"):
    capture = (dest / name).read_bytes()
    assert hashlib.sha256(capture).hexdigest() == manifest["GzipSHA256"]
    assert hashlib.sha256(gzip.decompress(capture)).hexdigest() == manifest["RawSHA256"]
print(json.dumps({"PackageCount": len(packages), "PackageErrors": [],
    "Compiled": sorted(compiled), "ManifestCompiled": expected, "CompiledEqual": True,
    "CheckedSourceHashes": len(sources), "SourceHashesUnchanged": True,
    "CaptureHashesUnchanged": True, "DependencyStdoutIdenticalToOriginals": True,
    "PreservedOriginalRecords": originals}, indent=2))
