"""Read-only Phase 3C.2 source, artifact, qualification and prose verification."""
import gzip
import hashlib
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[5]
PACKET = ROOT / 'docs/foundational/phase3c2'
sha = lambda data: hashlib.sha256(data).hexdigest()
manifest = json.loads((PACKET / 'evidence/manuscript/manifest.json').read_text())
for name, expected in manifest['current_paper'].items():
    assert sha((ROOT / 'paper' / name).read_bytes()) == expected, name
for item in manifest['artifacts']:
    packed = (ROOT / item['path']).read_bytes()
    raw = gzip.decompress(packed)
    assert sha(packed) == item['compressed_sha256'], item['path']
    assert sha(raw) == item['raw_sha256'] and len(raw) == item['raw_bytes']
for item in manifest.get('supporting_files', []):
    assert sha((ROOT / item['path']).read_bytes()) == item['sha256'], item['path']
archive = json.loads((PACKET / 'evidence/qualification/archive-manifest.json').read_text())
for item in archive['files']:
    raw = (ROOT / item['path']).read_bytes()
    assert sha(raw) == item['sha256'] and len(raw) == item['bytes'], item['path']
inputs = json.loads((PACKET / 'evidence/qualification/input-identity.json').read_text())
old = json.loads((ROOT / inputs['historical_manifest']).read_text())['files']
for files in [old, inputs['preexisting_tracked_inputs_outside_manifest_verified_against_base'], inputs['new_inputs']]:
    for name, expected in files.items():
        assert sha((ROOT / name).read_bytes()) == expected, name
assert all(e['matches_current'] for e in inputs['supplemental_historical_baseline_comparison'].values())
evidence = json.loads(gzip.decompress((PACKET / 'evidence/task1/export.json.gz').read_bytes()))
for item in evidence['Sources']:
    assert sha((ROOT / item['Path']).read_bytes()) == item['SHA256'], item['Path']
second = json.loads((PACKET / 'evidence/task2/manifest.json').read_text())
for item in second['artifacts']:
    assert sha((ROOT / item['path']).read_bytes()) == item['sha256'], item['path']
md = (ROOT / 'paper/alan-dissertation.md').read_text()
tex = (ROOT / 'paper/alan-dissertation.tex').read_text()
m = md.split('### 6.17 Shared Capacity and Stable Reservation Contracts (Phase 3C.2)\n', 1)[1].split('\n---\n', 1)[0].strip()
t = tex.split(r'\subsection{Shared Capacity and Stable Reservation Contracts (Phase 3C.2)}', 1)[1].split(r'\section{Conclusion}', 1)[0].strip()
t = t.replace(r"Coffman, Elphick and Shoshani's \emph{System Deadlocks}~\cite{coffman1971deadlocks}", "Coffman, Elphick and Shoshani's *System Deadlocks* (1971)")
t = t.replace(r'\path{docs/foundational/phase3c2/04-results.md}', '[Phase 3C.2 results](../docs/foundational/phase3c2/04-results.md)')
assert m == t, 'new mathematical content differs'
for start in ['If a crash leaves a non-idempotent request', 'The Phase 3B investigation is closed within those explicit fragments']:
    assert md[md.index(start):].split('\n', 1)[0] == tex[tex.index(start):].split('\n', 1)[0], start
links = 0
for path in [PACKET / '04-results.md', PACKET / '05-qualification.md', ROOT / 'docs/foundational/phase3/README.md', ROOT / 'research/capacitycomposition/README.md']:
    for target in re.findall(r'\[[^\]]+\]\(([^)]+)\)', path.read_text()):
        if '://' not in target and not target.startswith('#'):
            assert (path.parent / target.split('#', 1)[0]).exists(), (path, target)
            links += 1
assert manifest['warning_comparison']['exact_warning_lines_equal']
assert set(manifest['changed_current_pages']) <= set(manifest['rendered_and_visually_inspected_pages'])
readme = (ROOT / 'research/capacitycomposition/README.md').read_text()
race = (PACKET / 'evidence/task1/race.log').read_text()
assert '-race -timeout 30m' in readme and '-race -timeout 30m' in race
print(json.dumps({'result': 'PASS', 'historical_manifest_inputs': len(old), 'supplemental_inputs': len(inputs['supplemental_historical_baseline_comparison']), 'new_go_inputs': len(inputs['new_inputs']), 'producer_inputs': len(evidence['Sources']), 'task2_artifacts': len(second['artifacts']), 'qualification_files': len(archive['files']), 'compressed_manuscript_artifacts': len(manifest['artifacts']), 'local_links': links, 'new_section_aligned': True, 'changed_abstract_conclusion_equal': True, 'paper_hashes_match': True, 'warnings_unchanged': True, 'changed_pages_inspected': True, 'race_readme_matches_timeout': True}, indent=2))
