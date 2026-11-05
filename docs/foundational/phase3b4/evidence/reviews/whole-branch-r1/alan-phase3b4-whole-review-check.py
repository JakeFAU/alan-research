import hashlib,json,pathlib,tarfile,subprocess
root=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research'); ev=root/'docs/foundational/phase3b4/evidence/task5'
sha=lambda b:hashlib.sha256(b).hexdigest()
art=json.loads((ev/'artifacts.json').read_text()); assert all(sha((ev/p).read_bytes())==h for p,h in art.items()); print('Task5 artifact hashes verified:',len(art))
for mode in ['normal','race']:
 with tarfile.open(ev/f'independent-{mode}.tar.gz') as t:
  reports=[json.load(t.extractfile(m)) for m in t.getmembers() if m.isfile() and not m.name.endswith('-promotion.json')]
  counts={'cases':len(reports),'stores':sum(bool(r['Phases']) for r in reports),'captures':sum(len(r['Phases'] or []) for r in reports),'complete':sum(r['Final']['Custody']['Complete'] for r in reports),'C':sum(r['Projection']['Supported'] for r in reports),'contacts':sum(r['External']['Contacts'] for r in reports),'applications':sum(r['External']['Applications'] for r in reports),'losses':sum(r['External']['AcknowledgedLosses'] for r in reports)}
  assert counts==dict(cases=81,stores=80,captures=163,complete=40,C=33,contacts=89,applications=89,losses=11),counts
  assert all(p['Verdict']['Valid'] for r in reports for p in (r['Phases'] or [])); assert all(r['ContinuityFailure'] is None and (not r['Phases'] or r['External']['Valid']) for r in reports)
  print(mode,counts)
with tarfile.open(ev/'final-normal.tar.gz') as t:
 names=t.getnames()
 def read(p):return json.load(t.extractfile(p))
 def case(prefix):
  p=next(n for n in names if n.startswith(prefix+'-') and n.endswith('/case.json')); return p.rsplit('/',1)[0],read(p)
 for prefix in ['P02-v1','P02-v2','P10','N02','N06','C07','M03']:
  base,c=case(prefix); s=read(base+'/final/inspect.json'); rows=s['Tables']; facts=[json.loads(r['descriptor']) for r in rows['composition_facts']]; attempts=[json.loads(r['descriptor']) for r in rows['external_attempts']]; rec=c['Manifest']['Reciprocal']; q=rec['BARequest']; service=rec['BAService']; obs=next((r for r in rows['observations'] if json.loads(r['descriptor'])['ID']==json.dumps(['observation','message',q],separators=(',',':'))),None)
  a=next((a for a in attempts if a['Intent']['ID']==service),None)
  gates=[f for f in facts if f['Kind']=='authority_admission' and f['OperationID']==service]
  print(prefix,'q',q,'c',None if obs is None else obs['consumed_sequence'],'e',service,'first',None if a is None else (a['Allowed'],a['Reason']),'gates',[(f['Order'],f['AuthorityAdmission']['Allowed'],f['AuthorityAdmission']['Reason']) for f in gates])
  if prefix=='P02-v2':
   assert obs['consumed_sequence']==2 and a['Intent']['Sequence']==4 and a['Allowed']; w=a['Authority']; assert w['OriginalRequestID']==q and w['RequestObservationID']==obs['id'] and w['RequestConsumedSequence']==2
   h=next(h for h in rows['history'] if h['subject']=='a' and h['sequence']==2); assert h['selected_observation']==obs['id']
   admit=[(json.loads(r['descriptor']),json.loads(r['result'])) for r in rows['admissions'] if r['observation_id']==obs['id']]; assert len(admit)==1 and admit[0][0]['IntentID']==q and admit[0][1]['Accepted']
   grant=json.loads(next(r['result'] for r in rows['lifecycle_results'] if r['id']==rec['Grant'])); claim=json.loads(grant['Request']['Intent']['Payload']['claim']['Text']); assert grant['Accepted'] and claim['OriginalRequestID']==q and claim['ServiceOperationID']==service
   def order(kind,op=None,sub=None,seq=None):return next(f['Order'] for f in facts if f['Kind']==kind and (op is None or f['OperationID']==op) and (sub is None or f['Run']['Subject']==sub and f['Sequence']==seq))
   o=[order('execution',sub='b',seq=1),order('authorization',q),order('execution',sub='a',seq=2),order('execution',sub='a',seq=4),order('reservation',service),gates[0]['Order'],order('external_result',service)]; assert o==sorted(set(o)); assert order('lifecycle',rec['Grant'])<order('reservation',service)
   print('P02 exact q/admission/o/c/g/e/decision/gate/result causal orders',o)
   for endpoint in ['a','b']:
    lines=[json.loads(l) for l in t.extractfile(base+'/endpoint-'+endpoint+'.jsonl').read().splitlines()];print('endpoint',endpoint,lines)
  if prefix=='N02':assert obs['consumed_sequence']==3 and a['Intent']['Sequence']==2 and not a['Allowed'] and a['Authority']['RequestConsumedSequence']==0
  if prefix=='N06':assert a['Allowed'] and any(not f['AuthorityAdmission']['Allowed'] and f['AuthorityAdmission']['Reason']=='authority_revoked' for f in gates)
  if prefix=='P10':assert obs['consumed_sequence']==2 and a['Authority']['RequestConsumedSequence']==2 and a['Intent']['Sequence']==5
  if prefix=='C07':
   assert sum(ca['AcknowledgedSIGKILL'] for ca in c['Captures'])==2
   snaps={ca['Directory']:read(base+'/'+ca['Directory']+'/inspect.json') for ca in c['Captures']}
   for ca in c['Captures']:
    if ca['Directory'].endswith('post-open'):
     before=snaps[ca['Predecessor']]['Tables'];after=snaps[ca['Directory']]['Tables'];assert {k:v for k,v in before.items() if k!='runs'}=={k:v for k,v in after.items() if k!='runs'}
     out=[]
     for old,new in zip(before['runs'],after['runs']):
      op=json.loads(old['projection']);want=old['epoch']+(op['Status']=='open');assert new['epoch']==want
      x=dict(new);x['epoch']=old['epoch'];assert x==old;out.append((old['subject'],old['epoch'],new['epoch']))
     print('C07 actual Open',ca['Directory'],out)
  if prefix=='M03':
   found=[]
   for row in rows['admissions']:
    d=json.loads(row['descriptor']);r=json.loads(row['result'])
    if d['Kind']=='message' and d['Generation']==2:assert not r['Accepted'] and r['Bound']['Generation']==1 and r['ObservationID']=='';found.append(d['IntentID'])
   assert found;print('M03 target2 refused genuine originals',found)
with tarfile.open(ev/'checker-mutants.tar.gz') as t:
 fmap=json.load(t.extractfile('files.json'));print('mutant mapping shape',type(fmap).__name__)
 if isinstance(fmap,dict):
  for path,digest in fmap.items():
   if path.endswith('/record.json'):
    rec=json.load(t.extractfile('sha256/'+digest)); assert rec['built'] and rec['qualified_kill'];print('compiled mutant',rec['name'],[(c['exit'],c['log_sha256']) for c in rec['commands']])
print('Installed PDF',sha((root/'paper/alan-dissertation.pdf').read_bytes()))
