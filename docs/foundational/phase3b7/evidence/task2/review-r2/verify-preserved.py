"""Verify copied evidence and frozen inputs without writing to any input."""
import gzip
import hashlib
import json
from pathlib import Path

out = Path('docs/foundational/phase3b7/evidence/task2/review-r2')
manifest = json.loads(gzip.decompress(Path('docs/foundational/phase3b7/evidence/task1/prepared-source.json.gz').read_bytes()))
trace = json.loads((out / 'root-cause.json').read_bytes())
for row in trace['source_trace']:
    expected = manifest['overlays'].get(row['path'], manifest['runtime_files'].get(row['path']))
    assert all(expected[k] == row[k] for k in ('bytes', 'sha256'))
print('All three source trace members match the frozen manifest lengths/hashes.')
assert (out / 'initial-checks.json').read_bytes() == Path('docs/foundational/phase3b7/evidence/qualification/initial-checks.json').read_bytes()
for path, expected in json.loads((out / 'native-inputs.json').read_bytes())['inputs'].items():
    data = Path(path).read_bytes()
    assert len(data) == expected['bytes'] and hashlib.sha256(data).hexdigest() == expected['sha256']
for path, expected in trace['initial_reports'].items():
    data = Path(path).read_bytes()
    assert len(data) == expected['bytes'] and hashlib.sha256(data).hexdigest() == expected['sha256']
before = json.loads(Path('docs/foundational/phase3b7/evidence/qualification/checker-before.json').read_bytes())
assert hashlib.sha256(Path('/tmp/alan-phase3b7/matched-check').read_bytes()).hexdigest() == before['checker_sha256']
print('Initial summary, 12 excerpt inputs, four initial reports, and old frozen checker are unchanged.')
for name in ('green-normal.log', 'green-race.log'):
    text = (out / name).read_text()
    assert 'FAIL' not in text and 'WARNING: DATA RACE' not in text
    print(name, [line for line in text.splitlines() if line.startswith('ok')])
assert not (out / 'staticcheck.log').read_bytes()
