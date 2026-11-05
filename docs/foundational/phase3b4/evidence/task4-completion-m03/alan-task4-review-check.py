import pathlib,json,tarfile,hashlib,collections
P=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3b4/evidence/task4')
sha=lambda b:hashlib.sha256(b).hexdigest()
index=json.loads((P/'case-index.json').read_text());archives=json.loads((P/'archives.json').read_text());data={};checks=[]
for label,rec in archives.items():
 assert sha((P/rec['archive']).read_bytes())==rec['sha256']
 with tarfile.open(P/rec['archive']) as t:
  files={m.name:t.extractfile(m).read() for m in t if m.isfile()}
 for name,digest in rec['files'].items(): assert sha(files[label+'/'+name])==digest,(label,name)
 data.update(files)
checks.append({'archive_files':len(data),'archive_hashes':'verified all'})
source=json.loads((P/'runtime-source.json').read_text())
with tarfile.open(P/'runtime-source.tar.gz') as t:
 for name,h in source['files'].items():
  assert sha(t.extractfile(name).read())==h
  assert sha((pathlib.Path('/private/tmp/alan-phase3b3-runtime')/name).read_bytes())==h
checks.append({'runtime_source_files':len(source['files']),'head':source['head']})
for n in ['final-campaign-normal','final-campaign-race','evidence-fixed-normal','evidence-fixed-race','final-promoted-normal','final-promoted-race','focused-runtime-race','prior-integration','prior-race','evidence-fixed-staticcheck','final-vet']:
 r=json.loads((P/'commands'/f'{n}.json').read_text());assert r['exit']==0
 checks.append({'command':n,'argv':r['argv'],'source_differences':[f for f,h in source['files'].items() if r['source_sha256'].get(f)!=h],'log':(P/'commands'/f'{n}.log').read_text().strip()})
def read(base,file):return json.loads(data[base+'/'+file])
def raw(base,cap='final'):return read(base,cap+'/inspect.json')['Tables']
def typed(base,cap='final'):return read(base,cap+'/typed-inspection.json')
def dec(row):return json.loads(row['descriptor'])
def oid(s,n):return json.dumps([s,1,n,0],separators=(',',':'))
def obs(t,id):return next(r for r in t['observations'] if r['id']==id)
def original(t,id):return next(r for r in t['observations'] if json.loads(r['id'])==['observation','message',id])
def hstep(t,subject,seq):return next(r for r in t['history'] if r['subject']==subject and r['sequence']==seq)
def effect(t,id):return next((e for e in (t['Effects'] or []) if e['OperationID']==id),None)
def selection(t,o,seq):
 assert o['consumed_sequence']==seq
 assert hstep(t,o['subject'],seq)['selected_observation']==o['id']
 # Genuine receive opcode from exact artifact and preceding projection PC.
 run=next(r for r in t['runs'] if r['subject']==o['subject'])
 prog=next(r for r in t['programs'] if r['id']==run['program_id'])
 assert sha(prog['bytes'].encode())==prog['id']
 before=json.loads(hstep(t,o['subject'],seq-1)['successor'])
 instr=json.loads(prog['bytes'])['Code'][before['PC']];assert instr['Op']=='receive'
 d=dec(o);m=instr['Match'];assert m['Kind']==d['Kind'] and m['Provenance']==d['Provenance']
 for k,field in [('Source','Source'),('Generation','SourceGeneration'),('Obligation','Obligation'),('Operation','Operation')]:
  v=m[k];v=before['Data'][v['Ref']] if v.get('Ref') else v['Literal'];value=v['Number'] if k=='Generation' else v['Text'];assert value==d[field],(k,value,d[field])
def complete(base):
 t=raw(base);s=typed(base);m=read(base,'case.json')['Manifest']['Reciprocal'];end=read(base,'endpoints.json')
 chains=[]
 for req,svc,reply,server,caller,endpoint in [(m['BARequest'],m['BAService'],m['BAReply'],'a','b','b'),(m['ABRequest'],m['ABService'],m['ABReply'],'b','a','a')]:
  e=effect(s,svc);assert e['ExternalAttempt']['Allowed'];r=obs(t,e['Settlement']['ObservationID']);selection(t,r,r['consumed_sequence']);q=original(t,req);selection(t,q,q['consumed_sequence']);rep=original(t,reply);selection(t,rep,rep['consumed_sequence']);pay=dec(rep)['Payload'];assert pay['request_id']['Text']==req and pay['service_operation']['Text']==svc
  for field,F in [('mutation','Mutation'),('total','Total')]:assert pay[field]==dec(r)['Payload'][field] and pay[field]['Number']==e['ExternalReceipt'][F]
  assert q['consumed_sequence']<json.loads(svc)[2]<r['consumed_sequence']<json.loads(reply)[2]
  assert e['ExternalReceipt'] in end[endpoint]['Applications']
  assert all(g['Allowed'] for g in e['AuthorityAdmissions'])
  chains.append([req,q['consumed_sequence'],svc,r['consumed_sequence'],reply,rep['consumed_sequence']])
 assert all(json.loads(r['projection'])['Status']=='succeeded' for r in t['runs'])
 assert not s['Pending']
 messages=[dec(r)['ID'] for r in t['intents'] if dec(r)['Kind']=='message'];assert set(messages)=={r['intent_id'] for r in t['message_acknowledgments']}
 for term in t['terminal_obligations']:assert all(json.loads(term['delivery_status']).values())
 return chains
for mode,cases in index.items():
 for name,e in cases.items():assert sha(data[e['directory']+'/case.json'])==e['case_sha256'];assert read(e['directory'],'case.json')==e['case']
 for name in ['P01-v1','P01-v2','P02-v2','P04','P05','P06','P07','P08','P09','P10','P11','auxiliary-serviced-equal-total']:
  b=cases[name]['directory'];chains=complete(b);checks.append({'mode':mode,'case':name,'directory':b,'case_hash':cases[name]['case_sha256'],'chains':[list(x) for x in chains]})
 b=cases['P02-v1']['directory'];e=effect(typed(b),oid('a',4));assert not e['ExternalAttempt']['Allowed'] and e['ExternalAttempt']['Authority']['Reason']=='authority_request_not_admitted';assert all(x['Contacts']==0 for x in read(b,'endpoints.json').values())
 for name in ['P10','N01','N02']:
  b=cases[name]['directory'];t=raw(b);q=original(t,oid('b',1));e=effect(typed(b),oid('a',5 if name=='P10' else 2 if name=='N02' else 4));w=e['ExternalAttempt']['Authority'];record={'mode':mode,'case':name,'q_consumed':q['consumed_sequence'],'witness':{k:w[k] for k in ['Allowed','Reason','RequestObservationID','RequestConsumedSequence']}}
  if name!='N02':
   q2=original(t,oid('b',6));selection(t,q2,4 if name=='P10' else 3);assert dec(q)['Payload']==dec(q2)['Payload'];record['q2_consumed']=q2['consumed_sequence']
  if name=='N01':assert q['consumed_sequence'] is None and not w['Allowed']
  elif name=='N02':selection(t,q,3);assert w['RequestConsumedSequence']==0 and hstep(t,'a',2)['selected_observation'] is None
  else:selection(t,q,2);assert w['RequestObservationID']==q['id']
  checks.append(record)
 for name in ['N03','N04']:
  b=cases[name]['directory'];t=raw(b,'first-refusal');a=t['external_attempts'];assert a==raw(b,'after-later-facts')['external_attempts']==raw(b,'post-open')['external_attempts'];assert not t['composition_charges'];checks.append({'mode':mode,'case':name,'immutable_attempt':dec(a[0])['Reason'],'decision_outputs':[read(b,'decision-run-01.json'),read(b,'decision-run-02.json')]})
 for name in ['N06','C08']:
  b=cases[name]['directory'];s=typed(b);e=effect(s,oid('a',4));assert e['ExternalAttempt']['Allowed'] and not e.get('AuthorityAdmissions')[0]['Allowed'];assert len(raw(b)['composition_charges'])==1;assert all(v['Contacts']==0 for v in read(b,'endpoints.json').values());checks.append({'mode':mode,'case':name,'authority':e['ExternalAttempt']['Authority'],'gate':e['AuthorityAdmissions'][0]})
 for name in ['C01','C02','C03','C04','C05','C06','C07','C08']:
  b=cases[name]['directory'];c=cases[name]['case'];prev=None;pairs=[];kills=[]
  for cap in c['Captures']:
   t=raw(b,cap['Directory']);assert cap['Predecessor']==(prev['Directory'] if prev else '')
   if cap['Label']=='post-open':
    old=raw(b,prev['Directory']);assert old.keys()==t.keys()
    for table,rows in old.items():
     if table!='runs':assert rows==t[table],(name,table)
     else:
      for a,z in zip(rows,t[table]):
       expected=dict(a);expected['epoch']+=int(json.loads(a['projection'])['Status']=='open');assert expected==z
    pairs.append([prev['Directory'],cap['Directory'],[(r['subject'],r['epoch']) for r in t['runs']]])
   if cap['AcknowledgedSIGKILL']:
    loss=read(b,f'loss-{len(kills)+1:02d}.json');assert loss['PID']==cap['OwnerPID']==loss['Acknowledgment']['PID'];assert 'killed' in loss['Wait'];kills.append({'pid':loss['PID'],'ack':loss['Acknowledgment'],'sequences':[(r['subject'],r['sequence']) for r in t['runs']],'grant_exists':any(dec(r)['ID']==oid('b',2) for r in t['intents']),'service_exists':any(dec(r)['ID']==cases[name]['case']['Manifest']['Reciprocal']['BAService'] for r in t['intents'])})
   prev=cap
  assert len(kills)==(2 if name=='C07' else 1) and len({k['pid'] for k in kills})==len(kills)
  if name!='C08':complete(b)
  checks.append({'mode':mode,'case':name,'directory':b,'pairs':pairs,'kills':kills})
checks.append({'selected_cases':{k:len(v) for k,v in index.items()}})
pathlib.Path('/tmp/alan-task4-review-check.json').write_text(json.dumps(checks,indent=2)+'\n')
print('PASS',len(checks),'records; hashes, selected source, command provenance, real receive/artifact filters, complete chains, exact q, sticky refusal, all C Open tables and kills checked')
extra=[]
for mode,cases in index.items():
 for name in ['C01','C02','C03','C04','C05','C06','C07','C08']:
  b=cases[name]['directory'];m=cases[name]['case']['Manifest']['Reciprocal'];t=raw(b,'stopped-01');s=typed(b,'stopped-01');q=original(t,m['BARequest']);e=effect(s,m['BAService']);a=next(r for r in t['runs'] if r['subject']=='a')
  if name=='C01':assert q['consumed_sequence'] is None and original(raw(b,'final-owner/post-open'),m['BARequest'])['consumed_sequence'] is None
  else:selection(t,q,2)
  if name=='C02':assert a['sequence']==2
  if name=='C03':assert a['sequence']==6 and e is None;assert any(dec(r)['ID']==m['Grant'] for r in t['intents'])
  if name=='C04':assert e is not None and e['ExternalAttempt'] is None
  if name in ['C05','C08']:
   assert e['ExternalAttempt']['Allowed'] and not e.get('AuthorityAdmissions') and e['ExternalReceipt'] is None
   assert len(t['composition_charges'])==1
   new=effect(typed(b),m['BAService']);assert new['ExternalAttempt']==e['ExternalAttempt'];assert new['AuthorityAdmissions'][0]['Source']['Epoch']==a['epoch']+1
  if name=='C06':
   result=obs(t,e['Settlement']['ObservationID']);assert result['consumed_sequence'] is None
   new=effect(typed(b),m['BAService']);assert new['ExternalAttempt']==e['ExternalAttempt'] and new['AuthorityAdmissions']==e['AuthorityAdmissions'];assert read(b,'endpoints.json')['b']['Contacts']==1
  if name=='C07':
   assert a['sequence']==2 and e is None and not t['lifecycle_results'];t2=raw(b,'stopped-02');assert hstep(t2,'a',6);assert effect(typed(b,'stopped-02'),m['BAService']) is None;assert any(dec(r)['ID']==m['Grant'] for r in t2['intents'])
  extra.append({'mode':mode,'case':name,'cut_preconditions':'checked','inspect_sha256':sha(data[b+'/stopped-01/inspect.json'])})
 b=cases['auxiliary-serviced-equal-total']['directory'];t=raw(b);s=typed(b);e1=effect(s,oid('a',5));e2=effect(s,oid('a',13));w=e2['ExternalAttempt']['Authority'];assert w['OriginalRequestID']==oid('b',6) and w['RequestConsumedSequence']==4 and w['Use']['GrantID']==oid('b',7);selection(t,obs(t,e2['Settlement']['ObservationID']),14);assert e1['ExternalReceipt']['Total']==e2['ExternalReceipt']['Total']==1;assert e1['ExternalReceipt']['Mutation']==1 and e2['ExternalReceipt']['Mutation']==2;assert original(t,oid('a',17))['consumed_sequence'] is None;assert any(r['intent_id']==oid('a',17) for r in t['message_acknowledgments']);assert len(read(b,'endpoints.json')['b']['Applications'])==2;extra.append({'mode':mode,'case':'auxiliary-serviced-equal-total','second_authority':w,'result_selected':14,'principal_reply_selected':original(t,oid('a',16))['consumed_sequence'],'auxiliary_reply_selected':None})
 for name in ['generic-unknown','applied-unknown-revoked-before-fresh-gate']:
  b=cases[name]['directory'];s=typed(b);m=cases[name]['case']['Manifest']['Reciprocal'];e=effect(s,m['BAService']);end=read(b,'endpoints.json');assert e['Settlement'] is None
  if name=='generic-unknown':
   assert e['AdapterResult'] is None and len(e['AuthorityAdmissions'])==1;report=read(b,'fresh-owner/run.json')['Report'];assert report['Dispatches']==0 and report['UnknownAdapters']==[m['BAService']];assert end['a']['Contacts']==1 and len(end['a']['Applications'])==1
  else:
   assert e['ExternalReceipt'] is None and len(e['AuthorityAdmissions'])==2 and not e.get('AuthorityAdmissions')[1]['Allowed'];assert end['b']['Contacts']==1 and len(end['b']['Applications'])==1
  extra.append({'mode':mode,'case':name,'gates':[(g['Allowed'],g['Source']['Epoch']) for g in e['AuthorityAdmissions']],'contacts':{k:v['Contacts'] for k,v in end.items()}})
 for name in ['N10e-holder','N10e-recipient']:
  b=cases[name]['directory'];obstruction=read(b,'construction-obstruction.json');assert obstruction['ActualRetry']['Reason']=='invocation_retry_unsupported' and not obstruction['ActualRetry']['Accepted'];assert all(r['generation']==1 for r in raw(b)['subjects']);extra.append({'mode':mode,'case':name,'actual_retry':obstruction['ActualRetry']})
 for name,want in [('K01',13),('K02',15)]:
  b=cases[name]['directory'];s=typed(b);t=raw(b);assert s['StorageVersion']==want;pr=next(r for r in t['programs'] if r['id']==cases[name]['case']['Manifest']['ProgramIDs']['a']);assert pr['semantics']=='alan/1e/v1';e=effect(s,oid('a',3));a=e['AdapterAttempt'];candidate=e['AdapterResult']['Candidate'];wire=json.loads(a['Prepared']['Wire']);receipt=json.loads(candidate['Response']);end=read(b,'endpoints.json')['a'];assert wire==end['Requests'][0] and receipt==end['Applications'][0] and len(end['Requests'])==1;assert wire['Delta']==1 and a['Prepared']['Request']['Payload']['delta']['Number']==999;selection(t,obs(t,e['Settlement']['ObservationID']),4)
  parents=[k[:-len('/parent.json')] for k in data if k.startswith('promoted-'+mode+'/'+name+'-') and k.endswith('/parent.json')];assert len(parents)==1;pb=parents[0];parent=read(pb,'parent.json');assert parent['SHA256']['']==sha(data[b+'/final/runtime.sqlite']);before=read(pb,'opened-inspect.json')['Tables'];after=read(pb,'promoted-inspect.json')['Tables'];assert all(rows==after[tab] for tab,rows in before.items() if tab not in ['metadata','receipt_conflicts']);assert len(after['receipt_conflicts'])==1;assert read(pb,'typed-inspection.json')['StorageVersion']==want+1
  extra.append({'mode':mode,'case':name,'public_dir':b,'parent_sha256':parent['SHA256'],'public_format':want,'promoted_dir':pb,'promoted_format':want+1,'wire_tuple':a['Prepared']['Request']})
measure=json.loads((P/'cost-observations.json').read_text());assert collections.Counter((x['mode'],x['N']) for x in measure)==collections.Counter({(m,n):5 for m in index for n in [0,1,8,32]})
for x in measure:
 assert all(v>0 for v in x['Nanoseconds'].values());assert x=={'mode':x['mode'],'case':x['case'],**read(x['case'],'measurements.json')};assert read(x['case'],'query-plan.json')==['SEARCH admissions USING INDEX admissions_observation_id (observation_id=?)']
ident=json.loads((P/'identity-correspondence.json').read_text());assert len(ident)==50
for pair in ident:
 for side in ['reference','candidate']:
  v=pair[side];s=typed(v['case']);assert v['namespace']==s['Namespace'] and v['programs']=={r['ID']['Subject']:r['ProgramID'] for r in s['Runs']};assert set(v['roles'].values())<={r['ID'] for r in s['Operations']};assert v['cost']==read(v['case'],'final/cost-sidecar.json')
 assert pair['reference']['namespace']!=pair['candidate']['namespace']
extra.append({'cost_observations':len(measure),'identity_maps':len(ident),'status':'all checked'})
old=read('old-readers','results.json');assert old['binary_sha256']=='f08fc613f6411adb005adb9bf67409bb456eff69cd39599fa9652957e1a1db36'
for r in old['records']:
 if r['format'] in ['15','16']:assert r['open']['Error']=='unsupported version' and r['copy_before_sha256']==r['copy_after_sha256']
extra.append({'old_reader':old})
pathlib.Path('/tmp/alan-task4-review-extra.json').write_text(json.dumps(extra,indent=2)+'\n');print('PASS extra',len(extra),'cut, obstruction, K, cost, unknown and auxiliary checks')
target2=[];intent_counts=collections.Counter()
for mode,cases in index.items():
 for name,entry in cases.items():
  for cap in entry['case']['Captures'] or []:
   for row in raw(entry['directory'],cap['Directory'])['intents']:
    d=dec(row)
    if d['Kind']=='message':
     intent_counts[(mode,d['Generation'])]+=1
     if d['Target']=='a' and d['Generation']==2:target2.append((mode,name,cap['Directory'],d))
print('M03 actual message target2 records:',target2,'all message generation counts:',dict(intent_counts))
pathlib.Path('/tmp/alan-task4-review-m03.json').write_text(json.dumps({'target_a_generation2':target2,'message_generation_counts':{str(k):v for k,v in intent_counts.items()}},indent=2)+'\n')
