import pathlib,json,hashlib,subprocess,collections
R=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research');P=R/'docs/foundational/phase3b4/evidence/task5';sha=lambda b:hashlib.sha256(b).hexdigest()
head=subprocess.check_output(['git','rev-parse','HEAD'],cwd=R,text=True).strip();assert head=='b940f0d2a30f3d0362cc24ea2d7c0bbfdf4d74cc';assert not subprocess.check_output(['git','status','--porcelain'],cwd=R,text=True)
files=subprocess.check_output(['git','diff','--name-only','0c303c3d43d6582fed0a9caac24f50e71d844a19',head],cwd=R,text=True).splitlines();assert len(files)==150 and all(f.startswith('docs/foundational/phase3b4/') for f in files)
print('DELTA_EVIDENCE_ONLY',head,len(files),'clean')
art=json.loads((P/'artifacts.json').read_text())
for f,h in art.items():assert sha((P/f).read_bytes())==h,f
actual={str(f.relative_to(P)) for f in P.rglob('*') if f.is_file() and f.name!='artifacts.json'};assert set(art)==actual,(set(art)-actual,actual-set(art));print('ARTIFACT_MANIFEST',len(art),'complete exact')
for repo,root in [('research',R),('runtime',pathlib.Path('/private/tmp/alan-phase3b3-runtime'))]:
 m=json.loads((P/'freeze'/f'{repo}-source.json').read_text())
 for f,h in m.items():assert sha((root/f).read_bytes())==h,(repo,f)
 print('CURRENT_EXECUTABLE_MAP',repo,len(m),'exact')
parent=P/'parent-final-gates';summary=json.loads((parent/'gate-summary.json').read_text());print('SUMMARY_KEYS',list(summary))
for entry in summary['records']:
 name=entry['name'];record=parent/(name+'.json');d=json.loads(record.read_text());assert sha(record.read_bytes())==entry['record_sha256'];assert sha((parent/d['log']).read_bytes())==d['log_sha256']==entry['log_sha256'];assert d['source_sha256_before']==d['source_sha256_after'];repo='research' if name.startswith('research') else 'runtime';assert d['source_sha256_before']==json.loads((P/'freeze'/f'{repo}-source.json').read_text());assert d['exit']==entry['expected_exit'];assert not d['watchdog_expired']
 events=[]
 for line in (parent/d['log']).read_text().splitlines():
  try:e=json.loads(line)
  except:continue
  if isinstance(e,dict):events.append(e)
 assert not [e for e in events if e.get('Action')=='fail']
 if name=='research-ordinary':
  pp=sum(e.get('Action')=='pass' and not e.get('Test') for e in events);tp=sum(e.get('Action')=='pass' and bool(e.get('Test')) for e in events);skips=[e['Test'] for e in events if e.get('Action')=='skip' and e.get('Test')];assert pp==48 and tp==2528;print('ORDINARY',d['exit'],pp,tp,d['elapsed_seconds'],skips)
print('PARENT_RECORDS',len(summary['records']),'all hashes, exit statuses, frozen maps and no fail events exact')
new=json.loads((P/'mutant-source-after.json').read_text());print('MUTANT_POST_MAP',new)
print('FINAL_DELTA_OK')
