#!/usr/bin/env python3
"""Capture a frozen Phase 3C.5 producer into a fresh directory; never overwrite."""
import argparse
import gzip
import hashlib
import json
import os
from pathlib import Path
import platform
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[4]
OVERRIDES = {"GOCACHE": "/tmp/alan-phase3c5/go-cache",
             "GOMODCACHE": "/tmp/alan-phase3c5/mod-cache",
             "STATICCHECK_CACHE": "/tmp/alan-phase3c5/staticcheck-cache"}


def digest(data):
    return hashlib.sha256(data).hexdigest()


def write_json(path, value):
    path.write_text(json.dumps(value, sort_keys=True, indent=2) + "\n")


def run(out, name, argv):
    result = subprocess.run(argv, cwd=ROOT, env=dict(os.environ, **OVERRIDES), capture_output=True)
    packed = gzip.compress(result.stdout, mtime=0)
    (out / (name + ".stdout.gz")).write_bytes(packed)
    write_json(out / (name + ".command.json"), {
        "argv": argv, "cwd": str(ROOT), "environment": OVERRIDES,
        "exit": result.returncode, "stdout_archive": name + ".stdout.gz",
        "stdout_bytes": len(result.stdout), "stdout_sha256": digest(result.stdout),
        "stderr": result.stderr.decode("utf-8", "replace")})
    if result.returncode:
        raise RuntimeError(f"{name} exited {result.returncode}; exact result retained")
    return result.stdout


def objects(raw):
    decoder = json.JSONDecoder()
    text, pos = raw.decode(), 0
    while pos < len(text):
        while pos < len(text) and text[pos].isspace():
            pos += 1
        if pos == len(text):
            break
        value, pos = decoder.raw_decode(text, pos)
        yield value


def capture(out):
    write_json(out / "environment.json", {
        "Root": str(ROOT), "Python": sys.version, "Platform": platform.platform(),
        "Overrides": OVERRIDES, "CaptureScript": str(Path(__file__).resolve().relative_to(ROOT)),
        "CaptureScriptSHA256": digest(Path(__file__).read_bytes())})
    deps = run(out, "go-list-deps", ["go", "list", "-deps", "-json", "./research/residualantichains/cmd/export"])
    raw = run(out, "export", ["go", "run", "./research/residualantichains/cmd/export"])
    data = json.loads(raw)
    # The exact stdout is also available under the conventional fixture name.
    (out / "export.json.gz").write_bytes(gzip.compress(raw, mtime=0))
    compiled, modules = [], []
    for package in objects(deps):
        module = package.get("Module", {})
        if not module.get("Main"):
            continue
        if module.get("Path") != "github.com/JakeFAU/alan-research":
            raise RuntimeError("unexpected main module")
        module_dir = Path(module["Dir"])
        package_dir = Path(package["Dir"])
        relative = package_dir.relative_to(module_dir)
        module_record = {"Path": module["Path"], "Dir": module["Dir"], "GoMod": module["GoMod"]}
        if module_record not in modules:
            modules.append(module_record)
        for file in package.get("GoFiles", []) + package.get("CgoFiles", []):
            compiled.append((relative / file).as_posix())
        # No hidden local embedded/non-Go inputs are expected in this producer.
        for kind in ("EmbedFiles", "SFiles", "CFiles", "CXXFiles", "HFiles", "FFiles", "SysoFiles"):
            if package.get(kind):
                raise RuntimeError(f"unbound local compiler input {kind}")
    compiled.sort()
    manifest = [s["Path"] for s in data["Sources"] if s["Use"] == "compiled producer"]
    write_json(out / "compiled-roster.json", {
        "Modules": modules, "Normalization": "package.Dir relative to package.Module.Dir",
        "Compiled": compiled, "Manifest": manifest, "Equal": compiled == manifest})
    if compiled != manifest:
        raise RuntimeError("compiled source roster mismatch")
    reads = []
    for source in data["Sources"]:
        content = (ROOT / source["Path"]).read_bytes()
        if digest(content) != source["SHA256"]:
            raise RuntimeError("source changed during capture: " + source["Path"])
        reads.append(dict(source, Bytes=len(content)))
    if len(data["Input"]["Sources"]) != 40:
        raise RuntimeError("inherited source roster mismatch")
    for source in data["Input"]["Sources"]:
        if digest((ROOT / source["Path"]).read_bytes()) != source["SHA256"]:
            raise RuntimeError("inherited provenance changed: " + source["Path"])
    write_json(out / "read-inputs.json", reads)
    packed = (out / "export.json.gz").read_bytes()
    write_json(out / "result.json", {
        "Schema": data["Schema"], "RawBytes": len(raw), "RawSHA256": digest(raw),
        "CompressedBytes": len(packed), "CompressedSHA256": digest(packed),
        "CompiledRosterEqual": True, "InheritedSources": 40,
        "ReadInputs": len(reads), "Counts": data["Counts"]})
    write_json(out / "artifact-manifest.json", [
        {"Path": p.name, "Bytes": p.stat().st_size, "SHA256": digest(p.read_bytes())}
        for p in sorted(out.iterdir()) if p.is_file()])
    print(json.dumps({"RawBytes": len(raw), "RawSHA256": digest(raw), "Counts": data["Counts"]}, sort_keys=True))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("output_dir", type=Path)
    args = parser.parse_args()
    # mkdir is the first write and is atomic with respect to an existing path.
    # No command runs, source is read or output is touched on overwrite refusal.
    try:
        args.output_dir.mkdir()
    except FileExistsError:
        print("refusing existing capture directory", file=sys.stderr)
        return 1
    try:
        capture(args.output_dir)
    except Exception as exc:
        write_json(args.output_dir / "failure.json", {"Error": str(exc)})
        print(str(exc), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
