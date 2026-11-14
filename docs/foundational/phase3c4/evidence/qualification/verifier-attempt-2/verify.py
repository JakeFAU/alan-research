"""Deterministic artifact/source/alignment check; no model execution or SQLite I/O."""
from pathlib import Path
import gzip
import hashlib
import json
import re
import subprocess

ROOT = Path(__file__).resolve().parents[4]
PACKET = ROOT / 'docs/foundational/phase3c4'
QUAL = PACKET / 'evidence/qualification'
PAPER = PACKET / 'evidence/manuscript'
TASK1 = PACKET / 'evidence/task1'
TASK2 = PACKET / 'evidence/task2/qualification-1'
BASE = 'c369e25c8db268d418ab7c11573d46dbc1d9425f'


def require(condition, context):
    if not condition:
        raise ValueError(context)


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def read(path):
    return json.loads(path.read_text())


def blob(commit, path):
    return subprocess.check_output(['git', 'show', f'{commit}:{path}'], cwd=ROOT)


def check_files(records, folder=ROOT, commit=None):
    for row in records:
        path = row.get('path', row.get('Path'))
        digest = row.get('sha256', row.get('SHA256'))
        raw = (folder / path).read_bytes()
        require(sha(raw) == digest, f'hash mismatch: {path}')
        size = row.get('bytes', row.get('Bytes'))
        require(size is None or len(raw) == size, f'length mismatch: {path}')
        if commit:
            require(raw == blob(commit, path), f'Git identity mismatch: {path}')


def normalize(text):
    text = text.replace('$$', '').replace(r'\[', '').replace(r'\]', '')
    text = text.replace('$', '')
    text = re.sub(r'\[the Phase 3C.4 results\]\([^)]*\)', 'RESULTS', text)
    text = text.replace(r'\path{docs/foundational/phase3c4/04-results.md}', 'RESULTS')
    return ' '.join(text.split())


def main():
    baseline = read(QUAL / 'baseline-inputs.json')
    new_go = read(QUAL / 'new-go-inputs.json')
    embedded = read(QUAL / 'embedded-inputs.json')
    require(baseline['base'] == embedded['base'] == BASE, 'baseline source')
    require(len(baseline['files']) == 683, 'baseline input count')
    require(len(new_go['files']) == 12, 'new Go input count')
    require(len(embedded['files']) == 5, 'embedded input count')
    check_files(baseline['files'], commit=BASE)
    check_files(new_go['files'], commit=new_go['accepted_go_head'])
    check_files(embedded['files'], commit=BASE)
    tracked = subprocess.check_output(['git', 'ls-tree', '-r', '--name-only', BASE], cwd=ROOT).decode().splitlines()
    expected = sorted(p for p in tracked if p.endswith('.go') or Path(p).name in ('go.mod', 'go.sum', 'go.work', 'go.work.sum'))
    require(expected == sorted(row['path'] for row in baseline['files']), 'complete baseline Go/module/workspace roster')

    package = read(QUAL / 'package-comparison.json')
    previous_archive = blob(BASE, package['old_path'])
    previous_raw = gzip.decompress(previous_archive)
    current_raw = (QUAL / 'package-list/output.log').read_bytes()
    require(sha(previous_archive) == package['old_archive_sha256'], 'old package archive')
    require(sha(previous_raw) == package['old_raw_sha256'], 'old package raw')
    require(sha(current_raw) == package['new_raw_sha256'], 'new package raw')
    previous, current = previous_raw.decode().splitlines(), current_raw.decode().splitlines()
    require(len(previous) == package['old_count'] == 76, 'old package count')
    require(len(current) == package['new_count'] == 78, 'new package count')
    require(sorted(set(current)-set(previous)) == package['added'], 'added packages')
    require(sorted(set(previous)-set(current)) == package['removed'] == [], 'removed packages')
    for name in ('build', 'module-verify', 'module-tidy-diff', 'package-list'):
        job, result = read(QUAL / name / 'job.json'), read(QUAL / name / 'result.json')
        raw = (QUAL / name / 'output.log').read_bytes()
        require(job['head'] == new_go['accepted_go_head'], name + ' source head')
        require(result['exit_code'] == 0, name + ' exit')
        require(len(raw) == result['output_bytes'] and sha(raw) == result['output_sha256'], name + ' output')

    raw = gzip.decompress((TASK1 / 'export.json.gz').read_bytes())
    export = json.loads(raw)
    manifest = read(TASK1 / 'artifact-manifest.json')
    require(sha(raw) == manifest['RawSHA256'] == 'b40f1d5ac0aea0f43999502ea318cb128f545eb8fb1fd129451915a5aa84bf1f', 'final export hash')
    require(len(raw) == manifest['RawBytes'] == 8271854, 'final export size')
    require(raw == gzip.decompress((TASK1 / 'export-second.json.gz').read_bytes()) == gzip.decompress((TASK2 / 'export-third.json.gz').read_bytes()), 'three exact exports')
    require(export['Sources'] == read(TASK1 / 'sources.json'), 'source roster binding')
    check_files(export['Sources'])
    check_files(export['Input']['Sources'])
    require(len(export['Sources']) == 40 and len(export['Input']['Sources']) == 27, 'source counts')
    predecessor = (ROOT / export['Input']['Path']).read_bytes()
    decoded = gzip.decompress(predecessor)
    for key, data in [('Compressed', predecessor), ('Raw', decoded)]:
        require(len(data) == export['Input'][key+'Bytes'], 'predecessor '+key+' bytes')
        require(sha(data) == export['Input'][key+'SHA256'], 'predecessor '+key+' hash')
    require(json.loads(decoded)['Schema'] == export['Input']['Schema'] == 'alan/phase3c3/v1', 'predecessor schema')
    require(export['Schema'] == 'alan/phase3c4/v1', 'export schema')
    check_files(read(TASK2 / 'checker-sources.json'))
    archive = read(TASK2 / 'archive-manifest.json')
    require(len(archive) == 57, 'Task 2 archive count')
    check_files(archive, TASK2)
    commands = read(TASK2 / 'cli-commands.json')
    require(len(commands) == 50, 'CLI count')
    require(sum(c['exit'] == 0 for c in commands) == 2, 'valid CLI count')
    require(sum(c['exit'] != 0 for c in commands) == 48, 'hostile CLI count')
    for command in commands:
        require(command['exit'] == command['ExpectedExit'], 'CLI exit '+command['Name'])
        data = (TASK2 / command['InputArchive']).read_bytes()
        require(sha(data) == command['InputArchiveSHA256'], 'CLI archived input')
        if command['InputArchive'].endswith('.gz'):
            data = gzip.decompress(data)
        require(len(data) == command['InputRawBytes'] and sha(data) == command['InputRawSHA256'], 'CLI raw input')
    tests = read(TASK2 / 'tests.json')
    require(tests['exit'] == 0 and 'Ran 9 tests' in tests['stderr'] and '\nOK\n' in tests['stderr'], 'checker tests')
    c = export['Comparisons']
    require(all(s['NonblockingLeft'] == s['NonblockingRight'] for s in export['Substitutions']), 'substitution equality')
    counts = {
        'Plants': len(export['Plants']), 'Monitors': len(export['Monitors']),
        'Controllers': len(export['Controllers']), 'Products': len(export['Products']),
        'ProductStates': sum(len(p['Nodes']) for p in export['Products']),
        'ProductArcs': sum(len(n['Arcs']) for p in export['Products'] for n in p['Nodes']),
        'ProductCanComplete': sum(n['CanComplete'] for p in export['Products'] for n in p['Nodes']),
        'InitialCanComplete': sum(next(n['CanComplete'] for n in p['Nodes'] if n['ID'] == p['Initial']) for p in export['Products']),
        'NonblockingProducts': sum(p['Nonblocking'] for p in export['Products']),
        'Comparisons': len(c), 'EqualLanguages': sum(p['InitialLanguage']['Equal'] for p in c),
        'EqualProfiles': sum(p['Equal'] for p in c),
        'ResidualComparisons': sum(len(p['Residuals']) for p in c),
        'LanguageSubsetPairs': sum(len(p['InitialLanguage']['Pairs']) for p in c) + sum(len(r['Language']['Pairs']) for p in c for r in p['Residuals']),
        'ProfileSubsetPairs': sum(len(p['Pairs']) for p in c),
        'Substitutions': len(export['Substitutions']),
        'NonblockingSubstitutions': sum(s['NonblockingLeft'] for s in export['Substitutions']),
        'BlockingSubstitutions': sum(not s['NonblockingLeft'] for s in export['Substitutions']),
    }
    require(counts == export['Counts'] == manifest['Counts'] == read(TASK2 / 'result.json')['Counts'], 'measured counts')

    md, tex = (ROOT / 'paper/alan-dissertation.md').read_text(), (ROOT / 'paper/alan-dissertation.tex').read_text()
    ms = md.split('### 6.19 Residual Interfaces and Nonblocking Substitution (Phase 3C.4)\n\n')[1].split('\n---\n')[0]
    ts = tex.split('\\subsection{Residual Interfaces and Nonblocking Substitution (Phase 3C.4)}\n\n')[1].split('\\section{Conclusion}')[0]
    require(normalize(ms) == normalize(ts), 'new section alignment')
    ma = md.split('**Abstract.** ')[1].split('\n---')[0]
    ta = tex.split('\\begin{abstract}')[1].split('\\end{abstract}')[0]
    require(normalize(ma) == normalize(ta), 'Abstract alignment')
    mc = md.split('## 8. Conclusion: ALAN\'s Contributions and Open Mathematical Problems')[1].split('---')[0]
    tc = tex.split('\\section{Conclusion}')[1].split('\\bibliographystyle')[0]
    require(normalize(mc) == normalize(tc), 'Conclusion alignment')
    require('one of the eleven closed instruction forms' in md and 'one of the ten closed instruction forms' not in md, 'eleven forms')
    grammar = md.split('#### 3.1.2 Instructions and Retained Programs')[1].split('where:')[0]
    forms = re.findall(r'\\mathsf\{(Set|Add|Branch|Jump|Receive|Emit|End|ForEach|Next|Loop|Append)\}\(', grammar)
    require(len(forms) == len(set(forms)) == 11, 'displayed eleven-form grammar')
    require('bound semantic version and its validator' in grammar and 'bound semantic version and its validator' in tex, 'version scope')
    correction = md.split('### 3.4')[1].split('## 4.')[0]
    require('in Section 4.5' in correction and 'Theorem 4.8) remains unproved' in correction and 'We resolve these three failures' not in md, 'congruence correction')
    require(r'Theorem~\ref{thm:congruence} remains unproved' in tex and r'\label{thm:congruence}' in tex, 'TeX open theorem reference')
    require((ROOT / 'paper/references.bib').read_bytes() == blob(BASE, 'paper/references.bib'), 'unchanged bibliography')
    for n in (2772,29156,35672,14410,886,656,223,850,89,2376,792,1584):
        require(f'{n:,}' in ms, 'manuscript count '+str(n))
    links = 0
    files = [PACKET/'04-results.md', PACKET/'05-qualification.md']
    for path in files:
        for target in re.findall(r'\[[^\]]+\]\(([^)]+)\)', path.read_text()):
            if '://' not in target:
                require((path.parent / target.split('#')[0]).exists(), 'local link '+target)
                links += 1
    for source in ['docs/foundational/phase3/README.md', 'paper/alan-dissertation.md']:
        old = blob(BASE, source).decode()
        for target in re.findall(r'\[[^\]]+\]\(([^)]+)\)', (ROOT/source).read_text()):
            if target not in old and '://' not in target:
                require(((ROOT/source).parent / target.split('#')[0]).exists(), 'new local link '+target)
                links += 1
    paper_manifest = read(PAPER / 'manifest.json')
    check_files(paper_manifest['files'])
    require(paper_manifest['visual_review']['inspected_pages'] == [1,2,3]+list(range(9,52)), 'visual coverage')
    require(paper_manifest['normalized_warnings_equal'], 'warning comparison')
    require(all(row['exit_code'] == 0 for row in read(PAPER/'build-commands.json')), 'build exits')
    require(all(row['exit_code'] == 0 for row in read(PAPER/'render-commands.json')), 'render exits')
    check_files(read(QUAL/'artifact-manifest.json'))
    return dict(status='PASS', baseline_inputs=683, new_go_inputs=12, embedded_inputs=5,
                source_records=40, predecessor_sources=27, checker_sources=5, task2_archive_entries=57,
                packages={'old':76,'current':78}, counts=counts, CLI={'valid':2,'hostile':48},
                checker_test_methods=9, new_local_links=links,
                alignment=['new section','Abstract','Conclusion','editorial corrections'],
                PDF_pages={'baseline':50,'current':51}, visual_pages=46,
                scope='Identity and archived-result checks; no suite, model, runtime or SQLite execution')


if __name__ == '__main__':
    print(json.dumps(main(), indent=2, sort_keys=True))
