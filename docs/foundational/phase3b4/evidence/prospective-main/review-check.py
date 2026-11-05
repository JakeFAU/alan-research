import json,pathlib,hashlib
p=pathlib.Path('/tmp/alan-phase3b4-prospective-qualification'); sha=lambda b:hashlib.sha256(b).hexdigest(); summary=json.loads((p/'summary.json').read_text())
for s in summary['records']:
 f=p/'gates'/(s['name']+'.json');r=json.loads(f.read_text());log=p/'gates'/r['log'];assert sha(f.read_bytes())==s['record_sha256'] and sha(log.read_bytes())==r['log_sha256']==s['log_sha256'];assert r['exit']==0 and not r['watchdog_expired'] and r['source_unchanged'] and r['source_sha256_before']==r['source_sha256_after'];assert len(r['source_sha256_after'])==(351 if s['name'].startswith('runtime') else 562)
 for file,h in r['source_sha256_after'].items():assert sha((pathlib.Path(r['cwd'])/file).read_bytes())==h
 events=[]
 for line in log.read_text().splitlines():
  if line.startswith('{'):
   try:events.append(json.loads(line))
   except ValueError:pass
 assert not any(e.get('Action')=='fail' for e in events)
 if s['name'] in ['runtime-integration','runtime-integration-race']:
  passed=[e['Test'] for e in events if e.get('Action')=='pass' and e.get('Test')]
  for test in ['TestCompetingReadCannotSettleAdopterDebt/original','TestCompetingReadCannotSettleAdopterDebt/adopter_query','TestRootCompletionsMembershipQueriesOncePerRun','TestRootCompletionsFreshLateReceiptSettlementAndConflict']:
   assert test in passed,test
 print('gate',s['name'],'verified')
for mode in ['normal','race']:
 root=pathlib.Path('/tmp/alan-phase3b4-prospective-'+mode);vr=pathlib.Path('/tmp/alan-phase3b4-prospective-verdicts-'+mode)
 reports=[json.loads(f.read_text()) for f in vr.glob('*.json') if not f.name.endswith('-promotion.json')]
 counts=dict(cases=len(reports),stores=sum(bool(r['Phases']) for r in reports),captures=sum(len(r['Phases'] or []) for r in reports),principal_complete=sum(r['Final']['Custody']['Complete'] for r in reports),supported_C=sum(r['Projection']['Supported'] for r in reports),contacts=sum(r['External']['Contacts'] for r in reports),applications=sum(r['External']['Applications'] for r in reports),acknowledged_losses=sum(r['External']['AcknowledgedLosses'] for r in reports));assert counts==summary['corpora'][mode];assert all(r['ContinuityFailure'] is None and (not r['Phases'] or r['External']['Valid']) and all(ph['Verdict']['Valid'] for ph in (r['Phases'] or [])) for r in reports); print(mode,counts)
 for prefix in ['P02-v1','P02-v2','P10','N02','N06','C07']:
  cp=next(root.glob(prefix+'-*/case.json'));c=json.loads(cp.read_text());base=cp.parent;rows=json.loads((base/'final/inspect.json').read_text())['Tables'];rec=c['Manifest']['Reciprocal'];q=rec['BARequest'];service=rec['BAService'];obs=next((r for r in rows['observations'] if json.loads(r['descriptor'])['ID']==json.dumps(['observation','message',q],separators=(',',':'))),None);a=next(json.loads(r['descriptor']) for r in rows['external_attempts'] if json.loads(r['descriptor'])['Intent']['ID']==service);facts=[json.loads(r['descriptor']) for r in rows['composition_facts']];gates=[f for f in facts if f['Kind']=='authority_admission' and f['OperationID']==service]
  for f,h in json.loads((base/'digests.json').read_text()).items():assert sha((base/f).read_bytes())==h
  print(mode,prefix,'c',obs['consumed_sequence'],'e',service,'first',a['Allowed'],a['Reason'],'gate',[(f['AuthorityAdmission']['Allowed'],f['AuthorityAdmission']['Reason']) for f in gates])
  if prefix=='P02-v1':assert not a['Allowed'] and a['Reason']=='authority_request_not_admitted' and not gates
  if prefix=='P02-v2':
   assert obs['consumed_sequence']==2 and a['Intent']['Sequence']==4 and a['Allowed'];w=a['Authority'];assert w['OriginalRequestID']==q and w['RequestObservationID']==obs['id'] and w['RequestConsumedSequence']==2
   h=next(r for r in rows['history'] if r['subject']=='a' and r['sequence']==2);assert h['selected_observation']==obs['id'];ad=next(r for r in rows['admissions'] if r['observation_id']==obs['id']);assert json.loads(ad['descriptor'])['IntentID']==q and json.loads(ad['result'])['Accepted'];g=json.loads(next(r['result'] for r in rows['lifecycle_results'] if r['id']==rec['Grant']));cl=json.loads(g['Request']['Intent']['Payload']['claim']['Text']);assert g['Accepted'] and cl['OriginalRequestID']==q and cl['ServiceOperationID']==service
   def order(kind,op=None,sub=None,seq=None):return next(f['Order'] for f in facts if f['Kind']==kind and (op is None or f['OperationID']==op) and (sub is None or f['Run']['Subject']==sub and f['Sequence']==seq))
   orders=[order('execution',sub='b',seq=1),order('authorization',q),order('execution',sub='a',seq=2),order('execution',sub='a',seq=4),order('reservation',service),gates[0]['Order'],order('external_result',service)];assert orders==sorted(set(orders));assert order('lifecycle',rec['Grant'])<order('reservation',service);print(mode,'P02 exact causal orders',orders)
   for endpoint in ['a','b']:
    ev=[json.loads(l) for l in (base/('endpoint-'+endpoint+'.jsonl')).read_text().splitlines()];assert len(ev)==2 and ev[0]['Event']=='ingress' and ev[1]['Event']=='committed';req=json.loads(ev[0]['Body']);assert req['Key']==ev[1]['Receipt']['Key'];print(mode,'endpoint',endpoint,req['Sequence'],ev[1]['Receipt']['Mutation'])
  if prefix=='P10':assert obs['consumed_sequence']==2 and a['Intent']['Sequence']==5 and a['Authority']['RequestConsumedSequence']==2
  if prefix=='N02':assert obs['consumed_sequence']==3 and a['Intent']['Sequence']==2 and not a['Allowed'] and a['Authority']['RequestConsumedSequence']==0
  if prefix=='N06':assert a['Allowed'] and any(not f['AuthorityAdmission']['Allowed'] and f['AuthorityAdmission']['Reason']=='authority_revoked' for f in gates)
  if prefix=='C07':
   assert sum(ca['AcknowledgedSIGKILL'] for ca in c['Captures'])==2;print(mode,'C07 acknowledged losses',[(ca['OwnerPID'],ca['Directory']) for ca in c['Captures'] if ca['AcknowledgedSIGKILL']]);snaps={ca['Directory']:json.loads((base/ca['Directory']/'inspect.json').read_text()) for ca in c['Captures']}
   for ca in c['Captures']:
    if ca['Directory'].endswith('post-open'):
     b=snaps[ca['Predecessor']]['Tables'];a0=snaps[ca['Directory']]['Tables'];assert {k:v for k,v in b.items() if k!='runs'}=={k:v for k,v in a0.items() if k!='runs'}
     for old,new in zip(b['runs'],a0['runs']):
      assert new['epoch']==old['epoch']+(json.loads(old['projection'])['Status']=='open');x=dict(new);x['epoch']=old['epoch'];assert x==old
     print(mode,'C07 renewal',ca['Directory'],[(o['epoch'],n['epoch']) for o,n in zip(b['runs'],a0['runs'])])
 print(mode,'promotions',[(f.name,json.loads(f.read_text())['Valid']) for f in vr.glob('*-promotion.json')]);assert len(list(vr.glob('*-promotion.json')))==2 and all(json.loads(f.read_text())['Valid'] for f in vr.glob('*-promotion.json'))
print('All scoped evidence assertions passed; no suite rerun or SQLite open.')
