"""Read-only artifact, source identity and changed-prose alignment verification.

This does not rerun the producer, independent model audit or Go qualification.
"""
import gzip
import hashlib
import json
import re
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[5]
PACKET = ROOT / 'docs/foundational/phase3c3'


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def read(name):
    return json.loads((PACKET / name).read_text())


def checked(item):
    raw = (ROOT / item['path']).read_bytes()
    assert sha(raw) == item['sha256'], item['path']
    if 'bytes' in item:
        assert len(raw) == item['bytes'], item['path']
    return raw


manifest = read('evidence/manuscript/manifest.json')
for name, expected in manifest['current_paper'].items():
    assert sha((ROOT / 'paper' / name).read_bytes()) == expected, name
for item in manifest['artifacts']:
    raw = gzip.decompress(checked(item))
    assert sha(raw) == item['raw_sha256'] and len(raw) == item['raw_bytes'], item['path']
for item in manifest['supporting_files']:
    checked(item)
for item in read('evidence/qualification/accepted-task1-records.json'):
    checked(item)
archive = read('evidence/qualification/archive-manifest.json')
for item in archive['files']:
    checked(item)
inputs = read('evidence/qualification/baseline-inputs.json')
assert len(inputs['files']) == 669
for item in inputs['files']:
    raw = checked(item)
    old = subprocess.check_output(['git', 'show', inputs['base'] + ':' + item['path']], cwd=ROOT)
    assert old == raw, item['path']
new = read('evidence/qualification/new-go-inputs.json')
for item in new['files']:
    assert checked(item) == subprocess.check_output(['git', 'show', new['accepted_head'] + ':' + item['path']], cwd=ROOT)

for name in ['build-cache-writable', 'module-verify', 'module-tidy-diff', 'package-list-cache-writable']:
    result = read('evidence/qualification/' + name + '/result.json')
    raw = gzip.decompress((PACKET / 'evidence/qualification' / name / 'output.log.gz').read_bytes())
    assert result['exit_code'] == 0
    assert sha(raw) == result['output_sha256'] and len(raw) == result['output_bytes']
    job = read('evidence/qualification/' + name + '/job.json')
    assert job['head'] == new['accepted_head']
    if name in ['build-cache-writable', 'module-tidy-diff']:
        assert not raw
comparison = read('evidence/qualification/package-comparison.json')
prior_raw = (ROOT / comparison['prior_source']).read_bytes()
current_raw = gzip.decompress((PACKET / 'evidence/qualification/package-list-cache-writable/output.log.gz').read_bytes())
assert sha(prior_raw) == comparison['prior_source_sha256']
assert sha(current_raw) == comparison['current_output_sha256']
prior, current = set(prior_raw.decode().splitlines()), set(current_raw.decode().splitlines())
assert len(prior) == comparison['prior'] == 74 and len(current) == comparison['current'] == 76
assert sorted(current - prior) == comparison['added'] and sorted(prior - current) == comparison['removed'] == []
assert set(comparison['added']) == {'github.com/JakeFAU/alan-research/research/componentinterfaces', 'github.com/JakeFAU/alan-research/research/componentinterfaces/cmd/export'}

first = read('evidence/task1/artifact-manifest.json')
packed = (PACKET / 'evidence/task1/export.json.gz').read_bytes()
raw = gzip.decompress(packed)
assert len(raw) == first['raw']['bytes'] and sha(raw) == first['raw']['sha256']
assert len(packed) == first['gzip']['bytes'] and sha(packed) == first['gzip']['sha256']
assert packed[4:8] == bytes(4) and not packed[3] & 8
assert raw == gzip.decompress((PACKET / 'evidence/task2/export.json.gz').read_bytes())
evidence = json.loads(raw)
for item in evidence['Sources']:
    assert sha((ROOT / item['Path']).read_bytes()) == item['SHA256'], item['Path']
assert len(evidence['Sources']) == 27
second = read('evidence/task2/manifest.json')
assert len(second) == 125
for item in second:
    checked(item)
checker = read('evidence/task2/checker-sources.json')
assert len(checker) == 4
for item in checker:
    checked(item)
assert evidence['Counts'] == read('evidence/task2/independent-counts.json')
assert all(s['Marked'] for monitor in evidence['Monitors'] for s in monitor['States'])
for name in ['green-test', 'green-race', 'green-staticcheck']:
    record = read('evidence/task1/' + name + '.json')
    assert record['exit'] == 0 and not record['stderr']
    if name != 'green-staticcheck':
        assert len(re.findall(r'^--- PASS: ', record['stdout'], re.M)) == 9
    else:
        assert not record['stdout']
hostile = read('evidence/task2/hostile-index.json')
assert len(hostile) == 55
for item in hostile:
    packed = (ROOT / item['path']).read_bytes()
    assert sha(packed) == item['gzipSHA256'] and sha(gzip.decompress(packed)) == item['rawSHA256']
    assert item['exit'] == 1 and item['expected'] == 'reject' and item['reason'].startswith('REJECT:')
for item in read('evidence/task2/accepted-index.json'):
    assert item['exit'] == 0 and item['rawSHA256'] == first['raw']['sha256']

md = (ROOT / 'paper/alan-dissertation.md').read_text()
tex = (ROOT / 'paper/alan-dissertation.tex').read_text()
m = md.split('### 6.18 Completion-Preserving Continuation Interfaces (Phase 3C.3)\n', 1)[1].split('\n---\n', 1)[0].strip()
t = tex.split(r'\subsection{Completion-Preserving Continuation Interfaces (Phase 3C.3)}', 1)[1].split(r'\section{Conclusion}', 1)[0].strip()
for title in ['Snapshot counterexample.', 'Conditional source argument.', 'Residual boundary.']:
    t = t.replace('\\paragraph{' + title + '}', '**' + title + '**')
t = t.replace(r'\[', '$$').replace(r'\]', '$$')
t = t.replace(r"Tang and Moor's work on compositional non-blockingness~\cite{tang2024compositional}", "Tang and Moor's work on compositional non-blockingness (2024)")
t = t.replace(r'\path{docs/foundational/phase3c3/04-results.md}', '[Phase 3C.3 results](../docs/foundational/phase3c3/04-results.md)')
assert m == t, 'changed mathematical content differs'
assert 'The snapshot therefore omits a completion-relevant distinction in future ordering.' in m
assert 'The forgotten future transition relation is therefore semantically necessary' not in md + tex
for start in ['If a crash leaves a non-idempotent request', 'The Phase 3B investigation is closed within those explicit fragments']:
    assert md[md.index(start):].split('\n', 1)[0] == tex[tex.index(start):].split('\n', 1)[0], start
bib = (ROOT / 'paper/references.bib').read_text()
entry = bib.split('@article{tang2024compositional,', 1)[1].split('\n}', 1)[0]
for expected in ['Tang, Yiheng and Moor, Thomas', 'Compositional non-blockingness verification of finite automata with prioritised events', 'Discrete Event Dynamic Systems', '{34}', '{125--161}', '{2024}', '10.1007/s10626-024-00394-2']:
    assert expected in entry, expected
assert '28. **Tang, Y., & Moor, T.** (2024).' in md
assert 'DOI: 10.1007/s10626-024-00394-2.' in md
links = 0
for path, content in [(PACKET / '04-results.md', (PACKET / '04-results.md').read_text()), (PACKET / '05-qualification.md', (PACKET / '05-qualification.md').read_text()), (ROOT / 'docs/foundational/phase3/README.md', (ROOT / 'docs/foundational/phase3/README.md').read_text()), (ROOT / 'paper/alan-dissertation.md', m)]:
    for target in re.findall(r'\[[^\]]+\]\(([^)]+)\)', content):
        if '://' not in target and not target.startswith('#'):
            assert (path.parent / target.split('#', 1)[0]).exists(), (path, target)
            links += 1
comparison = manifest['warning_comparison']
assert comparison['exact_warning_lines_equal'] and comparison['baseline'] == comparison['current']
texts = {}
for stage, warning_key in [('base', 'baseline'), ('current', 'current')]:
    stage_prefix = manifest.get('artifact_stage_prefixes', {}).get(stage, stage)
    log = gzip.decompress((PACKET / 'evidence/manuscript' / (stage_prefix + '-alan-dissertation.log.gz')).read_bytes()).decode()
    warnings = [line for line in log.splitlines() if re.search(r'(Underfull \\hbox|Overfull \\hbox|LaTeX.*Warning)', line)]
    assert warnings == comparison[warning_key]
    texts[stage] = gzip.decompress((PACKET / 'evidence/manuscript' / (stage_prefix + '-text.txt.gz')).read_bytes()).decode().split('\f')
    if not texts[stage][-1].strip():
        texts[stage].pop()
    assert len(texts[stage]) == manifest['page_counts'][stage]
assert manifest['changed_current_pages'] == [i + 1 for i, text in enumerate(texts['current']) if i >= len(texts['base']) or text != texts['base'][i]]
for command in read('evidence/manuscript/build-commands.json') + read('evidence/manuscript/render-commands.json'):
    assert command['exit_code'] == 0
for path in manifest.get('additional_command_records', []):
    for command in json.loads((ROOT / path).read_text()):
        assert command['exit_code'] == 0
for name in ['alan-dissertation.tex', 'references.bib']:
    assert manifest['build_sources']['current'][name] == manifest['current_paper'][name]
    assert manifest['build_sources']['base'][name] == sha(subprocess.check_output(['git', 'show', manifest['baseline_commit'] + ':paper/' + name], cwd=ROOT))
assert set(manifest['changed_current_pages']) <= set(manifest['rendered_and_visually_inspected_pages'])
marker = read('evidence/manuscript/operation-marker.json')
assert marker['exit_code'] == 0 and marker['argv'][-6:] == ['--operation-kind', 'edit', '--expected-output-count', '1', '--output-format', 'pdf']
print(json.dumps({'result': 'PASS', 'baseline_inputs_against_current_and_git': 669,
                  'accepted_new_go_inputs': len(new['files']), 'producer_inputs': 27,
                  'task2_artifacts': len(second), 'checker_sources': len(checker),
                  'hostile_input_identities': len(hostile), 'root_packages': 76,
                  'go_top_level_tests_in_accepted_ordinary_and_race': 9,
                  'qualification_archive_files': len(archive['files']),
                  'manuscript_archive_files': len(manifest['artifacts']) + len(manifest['supporting_files']),
                  'local_links': links, 'changed_section_abstract_conclusion_aligned': True,
                  'new_reference_metadata_aligned': True, 'warnings_unchanged': True,
                  'changed_pages_inspected': True, 'paper_hashes_match': True}, indent=2))
