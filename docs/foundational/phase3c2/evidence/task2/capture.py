#!/usr/bin/env python3
"""Capture tooling only: the independent checker never invokes this exporter."""
import argparse
import gzip
import hashlib
import json
import os
import pathlib
import re
import subprocess
import sys
import time

HERE = pathlib.Path(__file__).resolve().parent
ROOT = HERE.parents[4]
TOOLS = ROOT / 'docs/foundational/phase3c2/tools'
SCRATCH = pathlib.Path('/tmp/alan-phase3c2')
sys.path.insert(0, str(TOOLS))
import audit
import test_audit

EXPECTED_RAW = 'a7cf316219b30cdc83d455b95a16549c022d55aa422b642b8f5976315e0c0fd2'


def write_json(path, data):
    path.write_text(json.dumps(data, sort_keys=True, indent=2) + '\n')


def packed(path, raw):
    with path.open('wb') as stream:
        with gzip.GzipFile(filename='', mode='wb', fileobj=stream, mtime=0) as archive:
            archive.write(raw)


def run(argv, name, output=None):
    start = time.monotonic()
    environment = os.environ.copy()
    environment.update(GOCACHE=str(SCRATCH / 'gocache'), PYTHONPYCACHEPREFIX=str(SCRATCH / 'pycache'))
    stdout = output or HERE / (name + '.stdout.log')
    stderr = HERE / (name + '.stderr.log')
    with stdout.open('wb') as out, stderr.open('wb') as err:
        result = subprocess.run(argv, cwd=ROOT, env=environment, stdout=out, stderr=err)
    record = dict(argv=argv, cwd=str(ROOT), exit_code=result.returncode,
                  elapsed_seconds=round(time.monotonic()-start, 6),
                  environment_overrides={k: environment[k] for k in ('GOCACHE', 'PYTHONPYCACHEPREFIX')},
                  stdout=str(stdout.relative_to(ROOT)) if stdout.is_relative_to(ROOT) else str(stdout),
                  stderr=str(stderr.relative_to(ROOT)))
    write_json(HERE / (name + '.command.json'), record)
    return result.returncode


def make_manifest():
    paths = [p for p in HERE.rglob('*') if p.is_file() and p.name != 'manifest.json']
    paths += [TOOLS / 'audit.py', TOOLS / 'test_audit.py',
              ROOT / 'docs/foundational/phase3c2/03-independent-audit.md']
    paths += [ROOT / name for name in ('docs/foundational/phase3a/01-contract.md',
                                      'docs/foundational/phase3c1/00-contract.md',
                                      'docs/foundational/phase3c1/01-argument.md',
                                      'docs/foundational/phase3c1/02-evidence-schema.md')]
    artifacts = []
    for path in sorted(paths):
        raw = path.read_bytes()
        record = dict(path=path.relative_to(ROOT).as_posix(), bytes=len(raw),
                      sha256=hashlib.sha256(raw).hexdigest())
        if path.suffix == '.gz':
            plain = gzip.decompress(raw)
            record.update(raw_bytes=len(plain), raw_sha256=hashlib.sha256(plain).hexdigest())
        artifacts.append(record)
    write_json(HERE / 'manifest.json', dict(schema='alan/phase3c2/task2-archive/v1',
                                           base_commit='60d46a43dd32caa07f0b814e6741213f506a9cb6',
                                           artifacts=artifacts))


def verify():
    m = json.loads((HERE / 'manifest.json').read_text())
    for r in m['artifacts']:
        raw = (ROOT / r['path']).read_bytes()
        assert len(raw) == r['bytes'], r['path']
        assert hashlib.sha256(raw).hexdigest() == r['sha256'], r['path']
        if 'raw_sha256' in r:
            plain = gzip.decompress(raw)
            assert len(plain) == r['raw_bytes'], r['path']
            assert hashlib.sha256(plain).hexdigest() == r['raw_sha256'], r['path']
    # Verify Task 1's manifest as well, without writing old evidence.
    old = json.loads((HERE.parent / 'task1/artifact-manifest.json').read_text())
    archive = (HERE.parent / 'task1' / old['Archive']).read_bytes()
    plain = gzip.decompress(archive)
    assert len(archive) == old['PackedBytes'] and len(plain) == old['RawBytes']
    assert hashlib.sha256(archive).hexdigest() == old['PackedSHA256']
    assert hashlib.sha256(plain).hexdigest() == old['RawSHA256'] == EXPECTED_RAW
    exported = audit.decode(plain)
    exact_sources = ''.join(r['SHA256'] + '  ' + r['Path'] + '\n' for r in exported['Sources'])
    assert (HERE.parent / 'task1/source-hashes.txt').read_text() == exact_sources
    assert (HERE / 'source-hashes.txt').read_text() == exact_sources
    audit.exact(exported['Sources'], audit.sources(ROOT), 'manifest source identity')
    print(json.dumps(dict(verified_task2_artifacts=len(m['artifacts']),
                          verified_task1_raw_sha256=old['RawSHA256'],
                          verified_source_hashes=len(exported['Sources'])), sort_keys=True))


def capture(skip_export):
    SCRATCH.mkdir(parents=True, exist_ok=True)
    plain = SCRATCH / 'task2-fresh-export.json'
    if not skip_export:
        status = run(['go', 'run', './research/capacitycomposition/cmd/export'], 'fresh-export', plain)
        assert status == 0, 'fresh exporter failed'
    raw = plain.read_bytes()
    old = gzip.decompress((HERE.parent / 'task1/export.json.gz').read_bytes())
    assert hashlib.sha256(old).hexdigest() == EXPECTED_RAW, 'reviewed capture hash mismatch'
    assert raw == old, 'fresh export differs from reviewed immutable export'
    packed(HERE / 'export.json.gz', raw)
    write_json(HERE / 'fresh-export-comparison.json', dict(equal_bytes=True, bytes=len(raw),
                                                          raw_sha256=hashlib.sha256(raw).hexdigest(),
                                                          reviewed_raw_sha256=EXPECTED_RAW))
    assert run([sys.executable, str(TOOLS / 'audit.py'), str(plain)], 'audit-plain') == 0
    assert run([sys.executable, str(TOOLS / 'audit.py'), str(HERE / 'export.json.gz')], 'audit-gzip') == 0
    assert (HERE / 'audit-plain.stdout.log').read_bytes() == (HERE / 'audit-gzip.stdout.log').read_bytes()
    document = audit.decode(raw)
    (HERE / 'source-hashes.txt').write_text(''.join(r['SHA256'] + '  ' + r['Path'] + '\n'
                                                 for r in document['Sources']))
    controls = HERE / 'hostile'
    controls.mkdir(exist_ok=True)
    records = []
    for name, mutant, reason in test_audit.hostile_cases(document):
        path = controls / (name + '.json.gz')
        packed(path, mutant)
        status = run([sys.executable, str(TOOLS / 'audit.py'), str(path)], 'hostile/' + name)
        stderr = (controls / (name + '.stderr.log')).read_text()
        assert status == 1 and re.search(reason, stderr), (name, status, stderr)
        records.append(dict(name=name, expected_reason=reason, exit_code=status,
                            input=str(path.relative_to(ROOT)),
                            raw_sha256=hashlib.sha256(mutant).hexdigest(),
                            stderr=stderr.strip(), classification='forged evidence control'))
    write_json(HERE / 'hostile-index.json', records)
    genuine = document['Graphs'][0]
    mutant = audit.load(controls / 'removed-non-tree-arc.json.gz')['Graphs'][0]
    assert len(genuine['Nodes']) == len(mutant['Nodes']) == mutant['States']
    assert mutant['Arcs'] == sum(len(n['Arcs']) for n in mutant['Nodes'])
    removed = [dict(**{'from': i}, arc=arc)
               for i, (original, changed) in enumerate(zip(genuine['Nodes'], mutant['Nodes']))
               for arc in original['Arcs'] if arc not in changed['Arcs']]
    for i, node in enumerate(mutant['Nodes'][1:], 1):
        assert dict(Action=node['Via'], To=i) in mutant['Nodes'][node['Parent']]['Arcs']
    assert len(removed) == 1
    write_json(HERE / 'non-tree-control-check.json',
               dict(removed=removed, states=mutant['States'], original_arcs=genuine['Arcs'],
                    mutated_arcs=mutant['Arcs'], all_bfs_parent_arcs_preserved=True,
                    reported_arc_total_matches_actual=True,
                    rejection='Graphs[0].Nodes[2].Arcs: length mismatch'))
    histories = test_audit.generated_histories()
    for w in histories:
        g, indexes = audit.graph(w['Config'])
        audit.replay(w, w['Config'], g, indexes, w['Name'])
    packed(HERE / 'python-generated-histories.json.gz',
           (json.dumps(histories, sort_keys=True, indent=2) + '\n').encode())
    write_json(HERE / 'python-generated-summary.json',
               dict(classification='checker test fixtures, not independent semantic evidence',
                    histories=len(histories), microsteps=sum(len(w['Steps']) for w in histories),
                    configurations=72, result_values=['r1', 'r2'],
                    independence='independent of Go; generator and replay share Python transition rules'))
    assert run([sys.executable, '-m', 'unittest', 'discover', '-s',
                'docs/foundational/phase3c2/tools', '-p', 'test_*.py', '-v'], 'unittest') == 0
    write_json(HERE / 'versions.json', dict(python=sys.version, executable=sys.executable))
    print(json.dumps(dict(hostile_controls=len(records), generated_histories=len(histories),
                          fresh_export_sha256=hashlib.sha256(raw).hexdigest()), sort_keys=True))


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--skip-export', action='store_true', help='reuse the already captured scratch stdout')
    parser.add_argument('--manifest', action='store_true')
    parser.add_argument('--verify', action='store_true')
    args = parser.parse_args()
    if args.manifest:
        make_manifest()
    elif args.verify:
        verify()
    else:
        capture(args.skip_export)
