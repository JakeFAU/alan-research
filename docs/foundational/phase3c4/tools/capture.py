#!/usr/bin/env python3
"""Archive exact CLI qualification. Creates a new attempt, never overwrites one."""
import copy
import gzip
import hashlib
import json
import os
import pathlib
import platform
import subprocess
import sys

sys.dont_write_bytecode = True
import audit
import test_audit

ROOT = audit.ROOT
BASE = ROOT / 'docs/foundational/phase3c4/evidence/task2'
ENVIRONMENT = dict(PYTHONDONTWRITEBYTECODE='1', GOCACHE='/tmp/alan-phase3c4/go-cache',
                   STATICCHECK_CACHE='/tmp/alan-phase3c4/staticcheck-cache')
ENV = dict(os.environ, **ENVIRONMENT)


def write_json(path, data):
    path.write_text(json.dumps(data, indent=2, sort_keys=True) + '\n')


def run(argv):
    result = subprocess.run(argv, cwd=ROOT, env=ENV, capture_output=True)
    return result, dict(argv=argv, cwd=str(ROOT), environment=ENVIRONMENT, exit=result.returncode,
                        stdout=result.stdout.decode('utf-8', 'replace'), stderr=result.stderr.decode('utf-8', 'replace'))


def main():
    attempt = BASE / sys.argv[1]
    attempt.mkdir()  # fail on reuse
    dependencies = ['docs/foundational/phase3c4/tools/' + n for n in ('audit.py', 'test_audit.py', 'capture.py')]
    dependencies += list(audit.PYTHON_DEPENDENCIES)
    write_json(attempt / 'checker-sources.json', [dict(Path=p, SHA256=audit.digest((ROOT / p).read_bytes())) for p in dependencies])
    write_json(attempt / 'environment.json', dict(Python=sys.version, Platform=platform.platform(), Overrides=ENVIRONMENT))
    result, record = run(['python3', '-m', 'unittest', 'discover', '-s', 'docs/foundational/phase3c4/tools', '-p', 'test_*.py', '-v'])
    write_json(attempt / 'tests.json', record)
    failures = [] if result.returncode == 0 else ['unittest']
    result, record = run(['go', 'run', './research/residualinterfaces/cmd/export'])
    raw = result.stdout
    record.pop('stdout')
    record.update(stdout_archive='export-third.json.gz', RawBytes=len(raw), RawSHA256=audit.digest(raw))
    (attempt / 'export-third.json.gz').write_bytes(gzip.compress(raw, mtime=0))
    write_json(attempt / 'export-third-command.json', record)
    prior_raw = gzip.decompress(test_audit.FIXTURE.read_bytes())
    same = raw == prior_raw
    write_json(attempt / 'fresh-equality.json', dict(Equal=same, Task1Path=str(test_audit.FIXTURE.relative_to(ROOT)),
               Task1RawSHA256=audit.digest(prior_raw), FreshRawSHA256=audit.digest(raw),
               SourceManifestEqual=audit.decode(raw)['Sources'] == audit.decode(prior_raw)['Sources']))
    if result.returncode or not same:
        failures.append('fresh producer equality')
    data = audit.decode(raw)
    # Cross-check the saved clean Task1 compiler enumeration, without repeating
    # Go tests or importing Go semantics. Every local compiled file must match.
    enumeration_path = ROOT / 'docs/foundational/phase3c4/evidence/task1/go-list-deps-resolved.json.gz'
    enumeration = gzip.decompress(enumeration_path.read_bytes()).decode()
    decoder, position, compiled = json.JSONDecoder(), 0, []
    while position < len(enumeration):
        while position < len(enumeration) and enumeration[position].isspace():
            position += 1
        if position == len(enumeration):
            break
        package, position = decoder.raw_decode(enumeration, position)
        directory = pathlib.Path(package.get('Dir', '/'))
        if directory.is_relative_to(ROOT):
            compiled += [str((directory / f).relative_to(ROOT)) for f in package.get('GoFiles', []) + package.get('CgoFiles', [])]
    manifest = [s['Path'] for s in data['Sources'] if s['Use'] == 'compiled producer']
    write_json(attempt / 'compiled-roster-check.json', dict(Enumeration=str(enumeration_path.relative_to(ROOT)),
               EnumerationSHA256=audit.digest(enumeration_path.read_bytes()), Compiled=sorted(compiled),
               Manifest=manifest, Equal=sorted(compiled) == manifest))
    if sorted(compiled) != manifest:
        failures.append('compiled roster')
    scratch = pathlib.Path('/tmp/alan-phase3c4/task2') / sys.argv[1]
    scratch.mkdir()
    plain = scratch / 'valid.json'
    plain.write_bytes(raw)
    cli = []

    def check(name, input_path, wanted_exit, input_raw, archive, diagnostic=None):
        result, record = run(['python3', 'docs/foundational/phase3c4/tools/audit.py', str(input_path)])
        record.update(Name=name, ExpectedExit=wanted_exit, InputRawSHA256=audit.digest(input_raw),
                      InputRawBytes=len(input_raw), InputArchive=archive,
                      InputArchiveSHA256=audit.digest((attempt / archive).read_bytes()))
        if diagnostic:
            record['ExpectedDiagnostic'] = diagnostic
        cli.append(record)
        write_json(attempt / 'cli-commands.json', cli)
        if result.returncode != wanted_exit or (diagnostic and diagnostic not in record['stderr']):
            failures.append(name)
        print(name, result.returncode, flush=True)

    check('valid-plain', plain, 0, raw, 'export-third.json.gz')
    check('valid-gzip', attempt / 'export-third.json.gz', 0, raw, 'export-third.json.gz')
    hostile = attempt / 'hostile'
    hostile.mkdir()
    for name, mutate, diagnostic in test_audit.mutations():
        forged = copy.deepcopy(data)
        mutate(forged)
        hostile_raw = json.dumps(forged, separators=(',', ':'), allow_nan=False).encode() + b'\n'
        archive = 'hostile/' + name + '.json.gz'
        path = attempt / archive
        path.write_bytes(gzip.compress(hostile_raw, mtime=0))
        check(name, path, 1, hostile_raw, archive, diagnostic)
    controls = [('duplicate-key', b'{"Schema":1,"Schema":2}', 'duplicate key'),
                ('nested-duplicate-key', b'{"x":{"a":1,"a":2}}', 'duplicate key'),
                ('nan', b'NaN', 'nonfinite'), ('infinity', b'Infinity', 'nonfinite'),
                ('negative-infinity', b'-Infinity', 'nonfinite'),
                ('overflow-number', b'{"x":1e999}', 'keys'),
                ('invalid-utf8', b'\xff', 'utf-8'), ('trailing-json', b'{} {}', 'Extra data'),
                ('truncated-json', b'{"Schema":', 'Expecting value')]
    for name, hostile_raw, diagnostic in controls:
        archive = 'hostile/' + name + '.json'
        path = attempt / archive
        path.write_bytes(hostile_raw)
        check(name, path, 1, hostile_raw, archive, diagnostic)
    write_json(attempt / 'result.json', dict(Failures=failures, CLICommands=len(cli),
               Accepted=sum(r['exit'] == 0 for r in cli), Rejected=sum(r['exit'] != 0 for r in cli), Counts=data['Counts']))
    archives = [dict(Path=str(p.relative_to(attempt)), Bytes=p.stat().st_size, SHA256=audit.digest(p.read_bytes()))
                for p in sorted(attempt.rglob('*')) if p.is_file()]
    write_json(attempt / 'archive-manifest.json', archives)
    return bool(failures)


if __name__ == '__main__':
    sys.exit(main())
