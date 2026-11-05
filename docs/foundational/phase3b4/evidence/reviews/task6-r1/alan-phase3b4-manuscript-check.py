import pathlib,json,hashlib,tarfile,re
R=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research');U=pathlib.Path('/private/tmp/alan-phase3b3-runtime');P=R/'docs/foundational/phase3b4';E=P/'evidence/task5'
sha=lambda b:hashlib.sha256(b).hexdigest()
for repo,root in [('research',R),('runtime',U)]:
 m=json.loads((E/'freeze'/f'{repo}-source.json').read_text())
 for name,h in m.items(): assert sha((root/name).read_bytes())==h,(repo,name)
 print('UNCHANGED_EXECUTABLE',repo,len(m))
rv=P/'evidence/reviews/task5-r1';manifest=json.loads((rv/'manifest.json').read_text())['files']
for name,d in manifest.items():assert sha((rv/name).read_bytes())==d['sha256']==sha(pathlib.Path(d['source']).read_bytes()),name
print('EXACT_REVIEW_PROVENANCE',len(manifest))
m=(R/'paper/alan-dissertation.md').read_text();t=(R/'paper/alan-dissertation.tex').read_text()
ms=m.split('### 6.12 Explicit Request Binding Across Intervening Computation\n\n')[1].split('\n## 7.')[0].strip();ts=t.split('\\subsection{Explicit Request Binding Across Intervening Computation}\n\n')[1].split('\\section{Conclusion}')[0].strip()
converted=re.sub(r'`([^`]+)`',lambda x:r'\path{'+x[1]+'}',ms).replace('[Phase 3B.4 results](../docs/foundational/phase3b4/08-results.md)',r'\path{docs/foundational/phase3b4/08-results.md}')
assert converted==ts
assert m.split('**Abstract.** ')[1].split('\n---')[0].strip()==t.split('\\begin{abstract}\n')[1].split('\\end{abstract}')[0].strip()
assert m.split("## 8. Conclusion: ALAN's Contributions and Open Mathematical Problems\n")[1].split('\n---')[0].strip()==t.split('\\section{Conclusion}\n')[1].split('\\bibliographystyle')[0].strip()
print('EXACT_EDITED_TEXT_PARITY','Abstract, 6.12, Conclusion')
for mode in ['normal','race']:
 with tarfile.open(E/f'independent-{mode}.tar.gz') as a:
  alljson=[json.load(a.extractfile(n)) for n in a.getnames() if n.endswith('.json')]
  reports=[d for d in alljson if 'Name' in d]
  vals=[v for d in reports for v in d.get('Phases') or []]
  assert all(v['Verdict']['Valid'] for v in vals)
  counts=[len(reports),sum(bool(d.get('Final',{}).get('Valid')) for d in reports),len(vals),sum(bool((d.get('Final',{}).get('Custody') or {}).get('Complete')) for d in reports),sum(bool(d.get('Projection',{}).get('Supported')) for d in reports),sum(d.get('External',{}).get('Contacts',0) for d in reports),sum(d.get('External',{}).get('Applications',0) for d in reports),sum(d.get('External',{}).get('AcknowledgedLosses',0) for d in reports)]
  assert counts==[81,80,163,40,33,89,89,11],counts
  print('INDEPENDENT_REPORT_COUNTS',mode,counts)
  for name in ['P02-v1','P02-v2','P10','auxiliary-serviced-equal-total']:
   d=next(d for d in reports if d['Name']==name); print('CUSTODY',mode,name,d['Final']['Custody']['Complete'],[(x['ID'],x['State']) for x in d['Final']['Custody'].get('Auxiliary') or []],d['External']['Contacts'],d['External']['Applications'])
 with tarfile.open(E/f'final-{mode}.tar.gz') as a:
  for name in ['P02-v1','P02-v2']:
   n=next(n for n in a.getnames() if n.startswith(name+'-') and n.endswith('/case.json'));base=n.rsplit('/',1)[0];c=json.load(a.extractfile(n));assert c['Schedule']=='automatic-run'
   raw=json.load(a.extractfile(base+'/final/inspect.json'));T=raw['Tables'];h=[r for r in T['history'] if r['subject']=='a' and r['sequence'] in [2,3,4]]
   ats=[json.loads(r['descriptor']) for r in T['external_attempts']];at=next(d for d in ats if d['Intent']['ID']=='["a",1,4,0]')
   assert h[0]['selected_observation']=='["observation","message","[\\"b\\",1,1,0]"]'
   assert h[1]['selected_observation']=='["observation","message","[\\"b\\",1,5,0]"]'
   assert at['Allowed']==(name=='P02-v2')
   if name=='P02-v1':assert at['Reason']=='authority_request_not_admitted'
   else:assert at['Authority']['RequestConsumedSequence']==2 and at['Authority']['OriginalRequestID']=='["b",1,1,0]'
   print('RAW_P02',mode,name,[(r['sequence'],r['selected_observation']) for r in h],at['Reason'])
print('MANUSCRIPT_EVIDENCE_CHECK_OK')
