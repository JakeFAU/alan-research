#!/usr/bin/env python3
"""Capture immutable checker inputs and commands. Refuse to overwrite captures."""
import copy
import gzip
import hashlib
import importlib.util
import io
import json
import os
import pathlib
import subprocess
import sys
import time

HERE = pathlib.Path(__file__).resolve().parent
ROOT = HERE.parents[4]
TOOLS = ROOT / 'docs/foundational/phase3c3/tools'
sys.path.insert(0, str(TOOLS))
import test_audit as tests

audit = tests.audit
SCRATCH = pathlib.Path('/tmp/alan-phase3c3')
ENV = dict(os.environ, PYTHONDONTWRITEBYTECODE='1', GOCACHE=str(SCRATCH / 'gocache'))


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def write(path, raw):
    with path.open('xb') as out:
        out.write(raw)


def record(path, data):
    write(path, (json.dumps(data, indent=2, sort_keys=True) + '\n').encode())


def compressed(raw):
    output = io.BytesIO()
    with gzip.GzipFile(filename='', mode='wb', fileobj=output, mtime=0) as out:
        out.write(raw)
    return output.getvalue()


def run(argv, capture_name, stdout_archive=None):
    start = time.monotonic()
    p = subprocess.run(argv, cwd=ROOT, env=ENV, capture_output=True)
    result = dict(argv=argv, cwd=str(ROOT), environment={k: ENV[k] for k in ('PYTHONDONTWRITEBYTECODE', 'GOCACHE')},
                  exit=p.returncode, stderr=p.stderr.decode(), duration_seconds=round(time.monotonic()-start, 3))
    if stdout_archive is not None:
        result.update(stdout_archive=stdout_archive, stdout_sha256=sha(p.stdout), stdout_bytes=len(p.stdout))
    else:
        result['stdout'] = p.stdout.decode()
    record(HERE / capture_name, result)
    return p


def main():
    fresh = run(['go', 'run', './research/componentinterfaces/cmd/export'],
                'fresh-export.command.json', 'export.json.gz (decompressed)')
    if fresh.returncode:
        raise RuntimeError('fresh capture failed; command retained')
    previous = gzip.decompress((HERE.parent / 'task1/export.json.gz').read_bytes())
    write(HERE / 'export.json.gz', compressed(fresh.stdout))
    record(HERE / 'fresh-export-comparison.json',
           dict(identical=fresh.stdout == previous, rawSHA256=sha(fresh.stdout), rawBytes=len(fresh.stdout),
                task1RawSHA256=sha(previous)))
    if fresh.stdout != previous:
        raise RuntimeError('fresh export differs; both exact captures retained')
    plain = SCRATCH / 'task2-fresh-export.json'
    plain.write_bytes(fresh.stdout)
    valid = []
    for name, path in [('plain', plain), ('gzip', HERE / 'export.json.gz')]:
        p = run(['python3', str(TOOLS / 'audit.py'), str(path)], 'audit-' + name + '.command.json')
        valid.append(dict(name=name, path=str(path), rawSHA256=sha(fresh.stdout), exit=p.returncode))
        if p.returncode:
            raise RuntimeError('valid input rejected')
    evidence = audit.decode(fresh.stdout)
    (HERE / 'hostile').mkdir(exist_ok=True)
    index = []
    for name, mutate in tests.controls().items():
        forged = copy.deepcopy(evidence)
        extra = mutate(forged)
        raw = (json.dumps(forged, separators=(',', ':')) + '\n').encode()
        path = HERE / 'hostile' / (name + '.json.gz')
        packed = compressed(raw)
        write(path, packed)
        p = run(['python3', str(TOOLS / 'audit.py'), str(path)], 'hostile/' + name + '.command.json')
        index.append(dict(name=name, path=str(path.relative_to(ROOT)), rawSHA256=sha(raw),
                          gzipSHA256=sha(packed), expected='reject', exit=p.returncode, reason=p.stderr.decode().strip()))
        if extra:
            record(HERE / 'non-tree-control-check.json', extra)
        print(name + ': ' + str(p.returncode), flush=True)
    for name, raw in [('duplicate-json-key', b'{"Schema":"alan/phase3c3/v1","Schema":"forged"}\n'),
                      ('nonfinite-json-number', b'{"count":NaN}\n')]:
        path = HERE / 'hostile' / (name + '.json.gz')
        packed = compressed(raw)
        write(path, packed)
        p = run(['python3', str(TOOLS / 'audit.py'), str(path)], 'hostile/' + name + '.command.json')
        index.append(dict(name=name, path=str(path.relative_to(ROOT)), rawSHA256=sha(raw),
                          gzipSHA256=sha(packed), expected='reject', exit=p.returncode, reason=p.stderr.decode().strip()))
    record(HERE / 'hostile-index.json', index)
    record(HERE / 'accepted-index.json', valid)
    unit = run(['python3', '-m', 'unittest', 'discover', '-s', 'docs/foundational/phase3c3/tools',
                '-p', 'test_*.py', '-v'], 'unittest.command.json')
    record(HERE / 'independent-counts.json', audit.audit(evidence))
    record(HERE / 'versions.json', dict(python=sys.version))
    paths = [TOOLS / 'audit.py', TOOLS / 'test_audit.py', audit.OLD_PATH, pathlib.Path(__file__)]
    record(HERE / 'checker-sources.json', [dict(path=str(p.relative_to(ROOT)), sha256=sha(p.read_bytes())) for p in paths])
    record(HERE / 'manifest.json', [dict(path=str(p.relative_to(ROOT)), sha256=sha(p.read_bytes()), bytes=p.stat().st_size)
                                   for p in sorted(HERE.rglob('*')) if p.is_file() and '__pycache__' not in str(p)
                                   and p.name != 'manifest.json'])
    if unit.returncode or any(row['exit'] == 0 for row in index):
        raise RuntimeError('qualification failed; accepted attacks and commands retained')
    print('ALL CAPTURES COMPLETE', flush=True)


if __name__ == '__main__':
    main()
