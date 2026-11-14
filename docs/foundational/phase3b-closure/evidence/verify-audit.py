"""Document/input checks only; does not execute historical models or open databases."""
import hashlib, json, pathlib, re, subprocess, tarfile
R=pathlib.Path(__file__).resolve().parents[4]
O=pathlib.Path(__file__).resolve().parent
BASE='7ae2c5ddd5e18d2907dd91fcefc896456a25562d'
def old(p): return subprocess.check_output(['git','show',BASE+':'+p],cwd=R)
def sha(b): return hashlib.sha256(b).hexdigest()
m=json.loads((O/'inspected-inputs.json').read_text())
for p,h in m['files'].items():
    b=old(p) if p=='docs/foundational/phase3/README.md' else (R/p).read_bytes()
    assert len(b)==h['bytes'] and sha(b)==h['sha256'],p
road=(R/'docs/foundational/phase3/README.md').read_text()
original=old('docs/foundational/phase3/README.md').decode()
for start,end in [('**Question.**','## First slice'),('Before calling Phase 3B established,',None)]:
    text=original[original.index(start):original.index(end) if end else len(original)].strip()
    assert text in road, start
assessment=(O.parent/'01-assessment.md').read_text()
for key in ['E1','E2','E3']+['G'+str(i) for i in range(1,10)]: assert len(re.findall(r'^\| '+key+r':',assessment,re.M))==1,key
links=0
for p in [O.parent/'01-assessment.md',O.parent/'02-verification.md',O.parent/'03-reviews.md',O/'manuscript/visual-review.md',R/'docs/foundational/phase3/README.md']:
    for target in re.findall(r'\]\(([^)]+)\)',p.read_text()):
        if '://' in target or target.startswith('#'): continue
        assert (p.parent/target.split('#')[0]).exists(),(str(p),target)
        links+=1
md=(R/'paper/alan-dissertation.md').read_text();tex=(R/'paper/alan-dissertation.tex').read_text()
start='The Phase 3B closure audit answers';end=' See '
md_body=md[md.index(start):md.index(end,md.index(start))]
tex_body=tex[tex.index(start):tex.index(end,tex.index(start))]
assert md_body==tex_body
c='The Phase 3B investigation is closed within those explicit fragments and negative boundaries. This audit adds no runtime qualification or unified refinement theorem; Phase 3C\'s compositional closure remains open.'
assert c in md and c in tex
oldmd=old('paper/alan-dissertation.md').decode()
def sec7(s): return s[s.index('## 7. Formal Falsification:'):s.index('## 8. Conclusion:')]
assert sec7(md)==sec7(oldmd)
assert (R/'paper/references.bib').read_bytes()==old('paper/references.bib')
changed=subprocess.check_output(['git','diff',BASE,'--name-only'],cwd=R,text=True).splitlines()
allowed=['docs/foundational/phase3b-closure/','docs/foundational/phase3/README.md','paper/alan-dissertation.']
assert all(any(p.startswith(a) for a in allowed) for p in changed),changed
facts=json.loads((O/'inspected-facts.json').read_text())
for label,d in facts['cancel_after_admission'].items():
    events=d['events'];admit=next(i for i,e in enumerate(events) if isinstance(e,dict) and e.get('Boundary')=='dispatch.admitted')
    decision=next(i for i,e in enumerate(events) if isinstance(e,dict) and e.get('Boundary')=='lifecycle.decided' and e.get('Phase')=='accepted')
    assert admit<decision<events.index('release held opportunity')
    assert (d['world_held_contacts'],d['world_final_contacts'],d['applications'])==(0,1,1)
    v=d['revoked_lifecycle'][0]
    assert v['Request']['Operation']=='CancelRun' and v['Accepted'] and v['Before']['Projection']['Status']=='open' and v['After']['Projection']['Status']=='cancelled'
for label,d in facts['matched'].items():
    assert all(d['byte_equal'].values())
    q=d['report']['Report'];assert q['Quiescent'] and q['Steps']==0 and q['Dispatches']==1 and len(q['OpenRuns'])==2
    assert q['UnknownExternal']==['["a",1,4,0]'] and q['PendingDeliveries']
    assert all(r['Pending']=='members' and r['Outcome'] is None for r in q['Roots'])
build_manifest=json.loads((O/'manuscript/manifest.json').read_text())
archive=O/'manuscript/build-and-renders.tar.gz'
assert sha(archive.read_bytes())==build_manifest['archive']['sha256']
with tarfile.open(archive) as t:
    for name,h in build_manifest['members'].items():
        b=t.extractfile(name).read()
        assert len(b)==h['bytes'] and sha(b)==h['sha256'],name
    for name in ['alan-dissertation.md','alan-dissertation.tex','alan-dissertation.pdf','references.bib']:
        assert t.extractfile(name).read()==(R/'paper'/name).read_bytes(),name
result={'base':BASE,'hashed_inputs_unchanged':len(m['files']),'matrix_rows':12,'local_links_checked':links,'original_exit_and_gates_preserved':True,'affected_md_tex_equal':True,'inherited_markdown_section7_unchanged':True,'bibliography_unchanged':True,'tracked_diff_within_document_scope':True,'build_archive_and_installed_sources_match':True,'cancel_trace_checks':4,'matched_raw_pair_checks':2,'pdf_sha256':sha((R/'paper/alan-dissertation.pdf').read_bytes())}
(O/'audit-checks.json').write_text(json.dumps(result,indent=2)+'\n')
print(json.dumps(result))
