import hashlib,json,pathlib,sys
root=pathlib.Path('/tmp/alan-phase3b4-prospective-qualification')
plan=json.loads((root/'plan.json').read_text());records=[]
for job in plan['jobs']:
 p=root/'gates'/(job['name']+'.json')
 if not p.exists():raise SystemExit('Gate not closed: '+job['name'])
 r=json.loads(p.read_text());log=p.with_suffix('.log');data=log.read_bytes()
 assert hashlib.sha256(data).hexdigest()==r['log_sha256'],p
 assert r['exit']==0 and r['source_unchanged'] and not r['watchdog_expired'],p
 expected=json.loads((root/(job['kind']+'-frozen-inputs.json')).read_text())
 assert r['source_sha256_before']==r['source_sha256_after']==expected,p
 events=[]
 for line in data.decode(errors='replace').splitlines():
  try:o=json.loads(line)
  except json.JSONDecodeError:continue
  if isinstance(o,dict):events.append(o)
 failures=[e for e in events if e.get('Action')=='fail'];assert not failures,p
 skips=[{k:e[k] for k in ('Package','Test','Elapsed') if k in e} for e in events if e.get('Action')=='skip' and 'Test' in e]
 if job['name'] in ('runtime-build','runtime-mod-tidy','runtime-format','runtime-vet','runtime-vet-integration','runtime-staticcheck','runtime-staticcheck-integration'): assert data==b'',p
 records.append({'name':job['name'],'argv':r['argv'],'environment':r['environment'],'exit':r['exit'],'source_unchanged':r['source_unchanged'],'elapsed_seconds':r['elapsed_seconds'],'package_passes':sum(e.get('Action')=='pass' and 'Test' not in e for e in events),'test_pass_events':sum(e.get('Action')=='pass' and 'Test' in e for e in events),'named_skips':skips,'record_sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'log_sha256':r['log_sha256']})
corpora={}
for mode in ('normal','race'):
 d=pathlib.Path('/tmp/alan-phase3b4-prospective-verdicts-'+mode)
 reports=[json.loads(p.read_text()) for p in sorted(d.rglob('*.json')) if not p.name.endswith('-promotion.json')]
 assert reports and all('Phases' in r for r in reports),d
 counts={'cases':len(reports),'stores':sum(bool(r['Phases']) for r in reports),'captures':sum(len(r['Phases'] or []) for r in reports),'principal_complete':sum(r['Final']['Custody']['Complete'] for r in reports),'supported_C':sum(r['Projection']['Supported'] for r in reports),'contacts':sum(r['External']['Contacts'] for r in reports),'applications':sum(r['External']['Applications'] for r in reports),'acknowledged_losses':sum(r['External']['AcknowledgedLosses'] for r in reports)}
 assert counts==dict(cases=81,stores=80,captures=163,principal_complete=40,supported_C=33,contacts=89,applications=89,acknowledged_losses=11),counts
 assert all(p['Verdict']['Valid'] for r in reports for p in (r['Phases'] or []))
 assert all(r['ContinuityFailure'] is None and (not r['Phases'] or r['External']['Valid']) for r in reports)
 corpora[mode]=counts
summary={'identity':json.loads((root/'identity.json').read_text()),'research_executable':'0c303c3d43d6582fed0a9caac24f50e71d844a19','records':records,'corpora':corpora,'all_15_gates_passed':len(records)==15,'limitations':'Separate prospective integration qualification; original frozen scientific pair remains unchanged. No new performance/theorem claim. No repeated full research suite/model exports. Initial sandbox listener failure is retained separately, not counted as a semantic RED.'}
(root/'summary.json').write_text(json.dumps(summary,indent=2)+'\n')
print(json.dumps({'gates':len(records),'counts':corpora,'skips':{r['name']:r['named_skips'] for r in records if r['named_skips']}},indent=2))
