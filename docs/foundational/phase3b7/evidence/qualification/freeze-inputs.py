"""Record or verify the exact executable and acceptance inputs for Phase 3B.7."""
import argparse
import hashlib
import json
from pathlib import Path
import subprocess


ACCEPTANCE = {
    "docs/foundational/phase3b7/01-contract.md",
    "docs/foundational/phase3b7/02-protocol.md",
    "docs/foundational/phase3b7/03-comparison.md",
}
CONTROLLER = "docs/foundational/phase3b7/evidence/qualification/"
EMBEDDED = {"examples/evidence/corpus.txt"}


def inputs(root):
    names = subprocess.check_output(["git", "ls-files", "-z"], cwd=root).decode().split("\0")
    result = {}
    for name in names:
        if not name:
            continue
        path = Path(name)
        controller = name.startswith(CONTROLLER) and "/" not in name[len(CONTROLLER):] and name.endswith(".py")
        executable = not name.startswith(("docs/", "paper/")) and (
            name.endswith((".go", ".go.txt", ".mod", ".sum", ".py", ".sh")) or "/testdata/" in name
        )
        if name not in ACCEPTANCE and name not in EMBEDDED and not executable and not controller:
            continue
        file = root / path
        if file.is_symlink() or not file.is_file():
            raise ValueError("unsupported input: " + name)
        raw = file.read_bytes()
        result[name] = {"bytes": len(raw), "sha256": hashlib.sha256(raw).hexdigest()}
    if not (ACCEPTANCE | EMBEDDED).issubset(result):
        raise ValueError("missing acceptance input")
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("root", type=Path)
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument("--write", type=Path)
    group.add_argument("--check", type=Path)
    args = parser.parse_args()
    root = args.root.resolve(strict=True)
    current = inputs(root)
    if args.check:
        expected = json.loads(args.check.read_text())["files"]
        if expected != current:
            changed = sorted(k for k in set(current) | set(expected) if current.get(k) != expected.get(k))
            raise ValueError("frozen input mismatch: " + repr(changed))
        print(json.dumps({"input_count": len(current), "unchanged": True}))
        return
    manifest = {
        "head": subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=root).decode().strip(),
        "tree": subprocess.check_output(["git", "rev-parse", "HEAD^{tree}"], cwd=root).decode().strip(),
        "files": current,
    }
    with args.write.open("x", encoding="utf-8") as output:
        json.dump(manifest, output, indent=2, sort_keys=True)
        output.write("\n")
    print(json.dumps({"input_count": len(current), "manifest": str(args.write)}))


if __name__ == "__main__":
    main()
