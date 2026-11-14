#!/usr/bin/env python3
"""Capture Task 2 tests and exact CLI attacks in a fresh directory only."""
import argparse
import copy
import gzip
import json
import os
from pathlib import Path
import platform
import subprocess
import sys
import tempfile

sys.dont_write_bytecode = True
import audit
import test_audit

ROOT = audit.ROOT
OVERRIDES = {'GOCACHE':'/tmp/alan-phase3c5/go-cache',
             'GOMODCACHE':'/tmp/alan-phase3c5/mod-cache',
             'STATICCHECK_CACHE':'/tmp/alan-phase3c5/staticcheck-cache',
             'PYTHONDONTWRITEBYTECODE':'1'}


def write_json(path, value):
    path.write_text(json.dumps(value, sort_keys=True, indent=2) + '\n')


def run(out, name, argv):
    result = subprocess.run(argv, cwd=ROOT, env=dict(os.environ, **OVERRIDES), capture_output=True)
    (out / (name + '.stdout.gz')).write_bytes(gzip.compress(result.stdout, mtime=0))
    (out / (name + '.stderr.gz')).write_bytes(gzip.compress(result.stderr, mtime=0))
    record = dict(argv=argv,cwd=str(ROOT),environment=OVERRIDES,exit=result.returncode,
                  stdout_bytes=len(result.stdout),stdout_sha256=audit.digest(result.stdout),
                  stderr_bytes=len(result.stderr),stderr_sha256=audit.digest(result.stderr),
                  stdout_archive=name+'.stdout.gz',stderr_archive=name+'.stderr.gz')
    write_json(out / (name + '.command.json'),record)
    return result


def capture(out):
    sources = list(audit.DEPENDENCIES) + [
        'docs/foundational/phase3c5/tools/' + name for name in ('audit.py','test_audit.py','capture_audit.py','capture.py')]
    write_json(out/'checker-sources.json',[dict(Path=p,SHA256=audit.digest((ROOT/p).read_bytes())) for p in sources])
    write_json(out/'environment.json',dict(Python=sys.version,Platform=platform.platform(),Root=str(ROOT),Overrides=OVERRIDES))
    export_dir = out/'fresh-export'
    result = run(out,'fresh-export-capture',['python3','docs/foundational/phase3c5/tools/capture.py',str(export_dir)])
    if result.returncode:
        raise RuntimeError('fresh producer capture failed; preserved')
    fresh = gzip.decompress((export_dir/'export.json.gz').read_bytes())
    previous = gzip.decompress(test_audit.FIXTURE.read_bytes())
    audit.exact(fresh,previous,'fresh semantic bytes')
    fixture = audit.decode(fresh)
    # Both captures bind current source files independently through capture.py.
    write_json(out/'semantic-identity.json',dict(Equal=True,RawBytes=len(fresh),RawSHA256=audit.digest(fresh),
               Previous=str(test_audit.FIXTURE.relative_to(ROOT)),Fresh='fresh-export/export.json.gz',
               Sources=fixture['Sources']))
    outcomes = []
    unit = run(out,'unittest',['python3','-m','unittest','discover','-s','docs/foundational/phase3c5/tools','-p','test_*.py','-v'])
    outcomes.append(dict(Name='unittest',Exit=unit.returncode,Passed=unit.returncode==0))
    cases = out/'cases'
    cases.mkdir()
    # Plain bytes are archived losslessly compressed, then materialized exactly
    # for the CLI run. Record both wire input and archived identities.
    def cli(name, raw, gzip_input, should_accept, diagnostic=''):
        wire = gzip.compress(raw,mtime=0) if gzip_input else raw
        archive = cases/(name+'.input.gz')
        archive.write_bytes(gzip.compress(wire,mtime=0))
        with tempfile.TemporaryDirectory(prefix='alan-phase3c5-audit-') as temp:
            path = Path(temp)/('input.json.gz' if gzip_input else 'input.json')
            path.write_bytes(wire)
            result = run(cases,name,['python3','docs/foundational/phase3c5/tools/audit.py',str(path)])
        reason = result.stderr.decode('utf-8','replace')
        passed = (result.returncode == 0 if should_accept else result.returncode != 0 and diagnostic in reason)
        record = dict(Name=name,ShouldAccept=should_accept,Passed=passed,Exit=result.returncode,
                      Diagnostic=diagnostic,InputArchive=archive.name,InputArchiveSHA256=audit.digest(archive.read_bytes()),
                      InputBytes=len(wire),InputSHA256=audit.digest(wire),Format='gzip' if gzip_input else 'plain',
                      Recovery='gzip-decompress InputArchive once to recover exact CLI input bytes')
        write_json(cases/(name+'.result.json'),record)
        outcomes.append(record)
        print(name + ': ' + ('PASS' if passed else 'FAIL (preserved)'),flush=True)
    cli('valid-plain',fresh,False,True)
    cli('valid-gzip',fresh,True,True)
    for index,(name,mutate,diagnostic) in enumerate(test_audit.mutations()):
        forged = copy.deepcopy(fixture)
        mutate(forged)
        raw = (json.dumps(forged,separators=(',',':'),ensure_ascii=False)+'\n').encode()
        cli(name,raw,index % 2 == 0,False,diagnostic)
    for name,raw,diagnostic in test_audit.json_attacks():
        cli(name,raw,False,False,diagnostic)
    failures = [r for r in outcomes if not r['Passed']]
    write_json(out/'results.json',dict(Outcomes=outcomes,Failures=failures,
               AcceptedHostiles=[r['Name'] for r in outcomes if r.get('ShouldAccept') is False and r['Exit']==0]))
    write_json(out/'artifact-manifest.json',[dict(Path=str(p.relative_to(out)),Bytes=p.stat().st_size,SHA256=audit.digest(p.read_bytes()))
               for p in sorted(out.rglob('*')) if p.is_file()])
    if failures:
        raise RuntimeError(str(len(failures))+' qualifications failed; exact bytes/results preserved')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('output_dir',type=Path)
    args = parser.parse_args()
    # Refuse before any capture command or write; never replace evidence.
    try:
        args.output_dir.mkdir()
    except FileExistsError:
        print('refusing existing capture directory',file=sys.stderr)
        return 1
    try:
        capture(args.output_dir)
    except Exception as exc:
        write_json(args.output_dir/'failure.json',dict(Error=str(exc)))
        print(str(exc),file=sys.stderr)
        return 1
    return 0


if __name__ == '__main__':
    sys.exit(main())
