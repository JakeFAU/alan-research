"""Check frozen identities and recorded qualification; do not rerun model/suites."""
from pathlib import Path
import gzip
import hashlib
import json
import re
import subprocess

ROOT = Path(__file__).resolve().parents[4]
P = ROOT / 'docs/foundational/phase3c5'
Q = P / 'evidence/qualification'
M = P / 'evidence/manuscript'
BASE = '347d93ba25a5f7a10f75144b1a44107be09a2f39'
ACCEPTED = '208e1411dd9074f4d97eb9fbba4fb9b8755b3236'
OLD = 'b079c67db3298688f32f2c5c8218a01f81d45601'


def require(value, message):
    if not value:
        raise ValueError(message)


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def read(path):
    return json.loads(path.read_text())


def blob(commit, path):
    return subprocess.check_output(['git', 'show', f'{commit}:{path}'], cwd=ROOT)


def check(rows, folder=ROOT, commit=None):
    for r in rows:
        name = r.get('Path', r.get('path'))
        raw = (folder / name).read_bytes()
        require(sha(raw) == r.get('SHA256', r.get('sha256')), 'hash: ' + name)
        require(len(raw) == r.get('Bytes', r.get('bytes', len(raw))), 'bytes: '+name)
        if commit:
            require(raw == blob(commit, name), 'Git identity: '+name)


def normal(s):
    s = s.replace(r'Theorem~\ref{thm:congruence}', 'Theorem 4.8')
    s = s.replace('$$', '').replace('$', '').replace(r'\{', '{').replace(r'\}', '}')
    s = s.replace('(Ware and Malik, 2011, [A State-Based Characterisation of the Conflict Preorder](https://arxiv.org/abs/1108.0230))', 'WARE')
    s = s.replace(r'~\cite{ware2011conflict}', 'WARE')
    s = s.replace('[the Phase 3C.5 results](../docs/foundational/phase3c5/04-results.md)', 'RESULTS')
    s = s.replace(r'\path{docs/foundational/phase3c5/04-results.md}', 'RESULTS')
    return ' '.join(s.split())


def main():
    baseline, embedded, new = [read(Q/n) for n in ['baseline-inputs.json','embedded-inputs.json','new-go-inputs.json']]
    require(baseline['base'] == embedded['base'] == new['base'] == BASE, 'base')
    require([len(x['files']) for x in [baseline,embedded,new]] == [695,5,15], 'input counts')
    check(baseline['files'], commit=BASE)
    check(embedded['files'], commit=BASE)
    check(new['files'], commit=new['head'])
    roster = subprocess.check_output(['git','ls-tree','-r','--name-only',BASE],cwd=ROOT).decode().splitlines()
    require(sorted(r['path'] for r in baseline['files']) == sorted(p for p in roster if p.endswith('.go') or Path(p).name in ['go.mod','go.sum','go.work','go.work.sum']), 'complete baseline roster')
    frozen = subprocess.check_output(['git','ls-tree','-r','--name-only',ACCEPTED,'--','docs/foundational/phase3c5','research/residualantichains'],cwd=ROOT).decode().splitlines()
    for name in frozen:
        require((ROOT/name).read_bytes() == blob(ACCEPTED,name), 'accepted source/archive drift: '+name)
    packages = read(Q/'package-comparison.json')
    old = blob(BASE,packages['path']); current=(Q/'package-list/stdout.log').read_bytes()
    require(old == (Q/'baseline-package-list.txt').read_bytes() and sha(old)==packages['raw_sha256'], 'actual historical package bytes')
    require(sha(current)==packages['new_raw_sha256'], 'current packages')
    a,b=old.decode().splitlines(),current.decode().splitlines()
    require(len(a)==78 and len(b)==80, 'package counts')
    require(sorted(set(b)-set(a))==packages['added'] and sorted(set(a)-set(b))==packages['removed']==[], 'package delta')
    for name in ['build','module-verify','module-tidy-diff','package-list']:
        job,result=read(Q/name/'job.json'),read(Q/name/'result.json')
        require(job['head']==new['head'] and result['exit_code']==0, 'root result '+name)
        for stream in ['stdout','stderr']:
            raw=(Q/name/(stream+'.log')).read_bytes()
            require(len(raw)==result[stream+'_bytes'] and sha(raw)==result[stream+'_sha256'], 'root stream '+name)
    task1=P/'evidence/task1'
    raw=gzip.decompress((task1/'capture-3/export.json.gz').read_bytes())
    require(sha(raw)=='ab3cbdfda577a3ae804eb76a94a8c7ac665cf00c69f992ad6804b6e0d0e27b8a' and len(raw)==14097260,'export identity')
    require(raw==gzip.decompress((task1/'capture-4/export.json.gz').read_bytes())==gzip.decompress((P/'evidence/task2/qualification-1/fresh-export/export.json.gz').read_bytes()), 'final export equality')
    e=json.loads(raw); require(e['Schema']=='alan/phase3c5/v1','schema')
    check(e['Sources']);check(e['Input']['Sources'])
    require(len(e['Sources'])==54 and len(e['Input']['Sources'])==40,'source counts')
    pred=(ROOT/e['Input']['Path']).read_bytes()
    for key,data in [('Compressed',pred),('Raw',gzip.decompress(pred))]:
        require(len(data)==e['Input'][key+'Bytes'] and sha(data)==e['Input'][key+'SHA256'],'predecessor '+key)
    capture=task1/'capture-3'
    # Interpret saved compiler paths relative to the compiler's recorded module root.
    text=gzip.decompress((capture/'go-list-deps.stdout.gz').read_bytes()).decode(); decoder=json.JSONDecoder();compiled=[];depcount=0
    recorded_root=Path(read(capture/'environment.json')['Root'])
    while text.strip():
        d,end=decoder.raw_decode(text.lstrip());text=text.lstrip()[end:];depcount+=1
        require(not d.get('Error') and not d.get('DepsErrors'),'compiler errors')
        module=d.get('Module',{})
        if module.get('Main'):
            module_root=Path(module['Dir']);require(module_root==recorded_root,'recorded module root')
            compiled.extend(str((Path(d['Dir'])/f).relative_to(module_root)) for f in d.get('GoFiles',[]))
    roster=read(capture/'compiled-roster.json')
    require(sorted(compiled)==roster['Compiled']==roster['Manifest'] and len(compiled)==10 and roster['Equal'],'compiled source roster')
    for manifest in [task1/'artifact-manifest.json',P/'evidence/task2/artifact-manifest.json',P/'evidence/task2/fix-1-artifact-manifest.json']:
        check(read(manifest),manifest.parent)
    accepted=P/'evidence/task2/qualification-2';sources=read(accepted/'checker-sources.json');require(len(sources)==9,'checker roster');check(sources)
    old_sources=read(P/'evidence/task2/qualification-1/checker-sources.json');snapshots=0
    for r in old_sources:
        name=r['Path']; data=blob(OLD,name);require(sha(data)==r['SHA256'],'historical source '+name)
        if name.startswith('docs/foundational/phase3c5/') and Path(name).name in ['audit.py','test_audit.py','capture_audit.py']:
            require(data==(P/'evidence/task2/fix-1-red'/Path(name).name).read_bytes(),'pre-fix snapshot');snapshots+=1
        else:require(data==(ROOT/name).read_bytes(),'unchanged historical source')
    check(read(accepted/'artifact-manifest.json'),accepted)
    outcomes=read(accepted/'results.json');require(not outcomes['Failures'] and not outcomes['AcceptedHostiles'],'accepted failures')
    cli=outcomes['Outcomes'][1:];require(len(cli)==58 and sum(x['ShouldAccept'] for x in cli)==2,'CLI counts')
    for r in cli:
        data=(accepted/'cases'/r['InputArchive']).read_bytes();require(sha(data)==r['InputArchiveSHA256'],'input archive')
        data=gzip.decompress(data);require(len(data)==r['InputBytes'] and sha(data)==r['InputSHA256'],'actual CLI input')
        require(r['Passed'] and r['Exit']==(0 if r['ShouldAccept'] else 1),'CLI result')
    commands=[accepted/'unittest.command.json']+sorted((accepted/'cases').glob('*.command.json'))
    require(len(commands)==59,'child count')
    for path in commands:
        c=read(path)
        for stream in ['stdout','stderr']:
            data=gzip.decompress((path.parent/c[stream+'_archive']).read_bytes())
            require(sha(data)==c[stream+'_sha256'] and len(data)==c[stream+'_bytes'],'child output')
            if stream=='stdout':
                lines=[x for x in data.decode().splitlines() if x.startswith('SOURCE_BUFFERS: ')]
                require(len(lines)==1 and json.loads(lines[0].split(': ',1)[1])==sources,'child source buffers')
    stderr=gzip.decompress((accepted/'unittest.stderr.gz').read_bytes()).decode()
    require('Ran 13 tests' in stderr and '\nOK\n' in stderr,'accepted tests')
    require(read(accepted/'unittest.command.json')['argv']==['python3','docs/foundational/phase3c5/tools/run_audit_tests.py'],'source-only entrypoint')
    require(all(x['NonblockingLeft']==x['NonblockingRight'] and x['InitialLeft']==x['InitialRight'] for x in e['Substitutions']),'substitutions')
    c=e['Counts']
    measured={'Plants':len(e['Plants']),'Monitors':len(e['Monitors']),'Controllers':len(e['Controllers']),'Products':len(e['Products']),'ProductStates':sum(len(p['Nodes']) for p in e['Products']),'ProductArcs':sum(len(n['Arcs']) for p in e['Products'] for n in p['Nodes']),'ProductCanComplete':sum(n['CanComplete'] for p in e['Products'] for n in p['Nodes']),'InitialCanComplete':sum(next(n['CanComplete'] for n in p['Nodes'] if n['ID']==p['Initial']) for p in e['Products']),'NonblockingProducts':sum(p['Nonblocking'] for p in e['Products']),'Substitutions':len(e['Substitutions']),'NonblockingSubstitutions':sum(s['NonblockingLeft'] for s in e['Substitutions'])}
    for k,v in measured.items():require(c[k]==v,'measured '+k)
    require(c['BlockingSubstitutions']==1928 and c['Comparisons']==36,'preserved counts')
    md=(ROOT/'paper/alan-dissertation.md').read_text();tex=(ROOT/'paper/alan-dissertation.tex').read_text()
    ms=md.split('### 6.20 Completion Languages and Residual Antichains (Phase 3C.5)\n\n')[1].split('\n---\n')[0]
    ts=tex.split('\\subsection{Completion Languages and Residual Antichains (Phase 3C.5)}\n\n')[1].split('\\section{Conclusion}')[0]
    require(normal(ms)==normal(ts),'new section alignment')
    require(normal(md.split('**Abstract.** ')[1].split('\n---')[0])==normal(tex.split('\\begin{abstract}')[1].split('\\end{abstract}')[0]),'Abstract alignment')
    require(normal(md.split("## 8. Conclusion: ALAN's Contributions and Open Mathematical Problems")[1].split('---')[0])==normal(tex.split('\\section{Conclusion}')[1].split('\\bibliographystyle')[0]),'Conclusion alignment')
    require('one of the eleven closed instruction forms' in md and 'Theorem 4.8) remains unproved' in md and r'Theorem~\ref{thm:congruence} remains unproved' in tex,'preserved editorial contracts')
    grammar=md.split('#### 3.1.2 Instructions and Retained Programs')[1].split('where:')[0]
    forms=re.findall(r'\\mathsf\{(Set|Add|Branch|Jump|Receive|Emit|End|ForEach|Next|Loop|Append)\}\(',grammar)
    require(len(forms)==len(set(forms))==11,'eleven forms')
    links=0
    for path in [P/'04-results.md',P/'05-qualification.md',ROOT/'docs/foundational/phase3/README.md',ROOT/'paper/alan-dissertation.md']:
        prior=blob(BASE,str(path.relative_to(ROOT))).decode() if path.name in ['README.md','alan-dissertation.md'] else ''
        for target in re.findall(r'\[[^\]]+\]\(([^)]+)\)',path.read_text()):
            if '://' not in target and target not in prior:
                require((path.parent/target.split('#')[0]).exists(),'local link '+target);links+=1
    manifest=read(M/'manifest.json');check(manifest['files']);check(manifest['build_inputs']['current'])
    for row in manifest['build_inputs']['base']:require(sha(blob(BASE,row['path']))==row['sha256'],'PDF baseline input')
    comparison=read(M/'comparison.json');require(comparison['normalized_warnings_equal'],'new PDF warning')
    require(manifest['visual_review']['inspected_pages']==comparison['rendered_pages'],'visual coverage')
    for name in ['base','current']:
        pages=gzip.decompress((M/(name+'-text.txt.gz')).read_bytes()).decode().split('\f')
        if not pages[-1].strip():pages.pop()
        require(len(pages)==comparison['page_counts'][name],'PDF page count')
    require(all(r['exit_code']==0 for r in read(M/'build-commands.json')+read(M/'render-commands.json')),'final build/render exits')
    marker=read(M/'operation-marker.json');require(marker['exit_code']==0 and marker['argv'][-6:]==['--operation-kind','edit','--expected-output-count','1','--output-format','pdf'],'marker')
    check(read(Q/'artifact-manifest.json'))
    return dict(status='PASS',baseline_inputs=695,embedded_inputs=5,new_go_inputs=15,accepted_frozen_files=len(frozen),producer_source_entries=54,compiled_local_files=10,dependency_packages=depcount,historical_snapshots=snapshots,accepted_checker_sources=9,accepted_child_reports=59,CLI={'valid':2,'hostile':56},checker_tests=13,counts=c,packages={'old':78,'current':80},local_links=links,PDF_pages=comparison['page_counts'],visual_pages=len(comparison['rendered_pages']),scope='Identity and archived-result checks only; no suite, model, runtime or SQLite execution')


if __name__ == '__main__':
    print(json.dumps(main(),indent=2,sort_keys=True))
