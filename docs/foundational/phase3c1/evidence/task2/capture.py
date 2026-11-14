"""Reproduce Task 2 export/audit artifacts from the repository root (stdlib)."""
import gzip
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[5]
OUT = Path(__file__).resolve().parent
SCRATCH = Path('/tmp/alan-phase3c1/task2')
TOOLS = ROOT / 'docs/foundational/phase3c1/tools'
sys.path.insert(0, str(TOOLS))
from test_audit import controls  # noqa: E402 -- test-only mutation generator


def digest(data):
    return hashlib.sha256(data).hexdigest()


def relative(path):
    try:
        return str(path.relative_to(ROOT))
    except ValueError:
        return str(path)


def main():
    SCRATCH.mkdir(parents=True, exist_ok=True)
    (OUT / 'hostile').mkdir(exist_ok=True)
    commands = []
    artifacts = []

    def run(argv, name, env=None, expected=0, stdout=None):
        overrides = env or {}
        proc = subprocess.run(argv, cwd=str(ROOT), env=dict(os.environ, **overrides), capture_output=True)
        log = OUT / (name + '.log')
        if stdout:
            stdout.write_bytes(proc.stdout)
            log.write_bytes(proc.stderr)
        else:
            log.write_bytes(proc.stdout + proc.stderr)
        commands.append(dict(argv=argv, environment=overrides, exit_code=proc.returncode,
                             log=relative(log), stdout=relative(stdout) if stdout else None))
        if proc.returncode != expected:
            raise RuntimeError('{} exited {}, expected {}; inspect {}'.format(name, proc.returncode, expected, log))
        return proc.stdout

    def compress(raw, path):
        with path.open('wb') as stream:
            with gzip.GzipFile(filename='', mode='wb', fileobj=stream, mtime=0) as archive:
                archive.write(raw)
        artifacts.append(dict(path=relative(path), raw_bytes=len(raw), raw_sha256=digest(raw),
                              compressed_sha256=digest(path.read_bytes()), gzip_mtime=0))

    python_version = run(['python3', '--version'], 'python-version').decode().strip()
    go_version = run(['go', 'version'], 'go-version').decode().strip()
    head = run(['git', 'rev-parse', 'HEAD'], 'source-head').decode().strip()
    first, second = SCRATCH / 'fresh-1.json', SCRATCH / 'fresh-2.json'
    env = {'GOCACHE': '/tmp/alan-phase3c1/task1/go-cache'}
    run(['go', 'run', './research/contractcomposition/cmd/export'], 'export-1', env, stdout=first)
    run(['go', 'run', './research/contractcomposition/cmd/export'], 'export-2', env, stdout=second)
    raw = first.read_bytes()
    sample = ROOT / 'docs/foundational/phase3c1/evidence/task1/export.json'
    assert raw == second.read_bytes() == sample.read_bytes(), 'fresh exports differ from each other or frozen sample'
    compress(raw, OUT / 'export.json.gz')
    run(['python3', str(TOOLS / 'audit.py'), str(first)], 'audit')
    run(['python3', '-m', 'unittest', 'discover', '-s', str(TOOLS), '-p', 'test_*.py'], 'unittest')
    genuine = json.loads(raw)
    rejections = []
    for name, reason, hostile in controls(genuine):
        modified = (json.dumps(hostile, sort_keys=True, separators=(',', ':')) + '\n').encode()
        temp = SCRATCH / (name + '.json')
        temp.write_bytes(modified)
        compress(modified, OUT / 'hostile' / (name + '.json.gz'))
        response = run(['python3', str(TOOLS / 'audit.py'), str(temp)], 'hostile/' + name, expected=1)
        verdict = json.loads(response)
        assert verdict['status'] == 'rejected' and reason in verdict['reason'], (name, verdict)
        rejections.append(dict(name=name, required_reason=reason, actual_reason=verdict['reason'],
                               exit_code=1, input_raw_sha256=digest(modified)))
    (OUT / 'hostile-rejections.json').write_text(json.dumps(rejections, indent=2, sort_keys=True) + '\n')
    # Hash inputs, including producer files, without treating producer code as
    # executable oracle logic. Claims prose remains a separate source-review duty.
    task1 = json.loads((ROOT / 'docs/foundational/phase3c1/evidence/task1/source-manifest.json').read_text())
    sources = list(task1['SourceSHA256']) + [relative(p) for p in sorted(TOOLS.glob('*.py'))]
    sources += [relative(OUT / 'protocol-interface.md'), relative(Path(__file__).resolve()),
                'docs/foundational/phase3c1/03-independent-audit.md']
    source_hashes = {p: digest((ROOT / p).read_bytes()) for p in sorted(sources)}
    assert all(source_hashes[p] == h for p, h in task1['SourceSHA256'].items()), 'producer/spec input drift'
    manifest = dict(schema='alan.phase3c1.task2.provenance.v1', task_base=head,
                    tools=dict(python=python_version, go=go_version), commands=commands,
                    source_sha256=source_hashes, artifacts=artifacts,
                    fresh_exports_equal=True, frozen_sample_sha256=digest(sample.read_bytes()),
                    hostile_rejections=len(rejections),
                    historical_logs=[
                        dict(path='red-initial.log', outcome='36 assertion failures before implementation; outer shell exit was not captured separately'),
                        dict(path='green-attempt-1.log', outcome='one wrong expected-reason failure in down-actor fixture; corrected post-states to reach availability guard'),
                        dict(path='red-loss-witness.log', exit_code=1, outcome='two assertions expose missing loss-witness coverage; subsequently implemented')])
    manifest['retained_file_sha256'] = {relative(p): digest(p.read_bytes()) for p in sorted(OUT.rglob('*'))
                                      if p.is_file() and p.name != 'manifest.json' and '__pycache__' not in p.parts}
    (OUT / 'manifest.json').write_text(json.dumps(manifest, indent=2, sort_keys=True) + '\n')
    print(json.dumps(dict(status='captured', commands=len(commands), rejections=len(rejections),
                          raw_sha256=digest(raw), bytes=len(raw)), sort_keys=True))


if __name__ == '__main__':
    main()
