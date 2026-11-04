import base64,gzip,hashlib,json,pathlib,sys
base=pathlib.Path(sys.argv[1]) if len(sys.argv)>1 else pathlib.Path(__file__).resolve().parent
output=pathlib.Path(sys.argv[2]) if len(sys.argv)>2 else pathlib.Path('/tmp/alan-phase3b-runtime-audit.json')
cases=[]
for mode in ('normal','race'):
 root=base/('runtime-'+mode)
 files=sorted(root.glob('*/evidence.json.gz'))
 assert len(files)==16,(mode,len(files))
 for path in files:
  ev=json.loads(gzip.decompress(path.read_bytes()))
  for a in ev['Provision']['Programs']:
   raw=base64.b64decode(a['Bytes'])
   assert hashlib.sha256(raw).hexdigest()==a['ID']
   assert (path.parent/('program-'+a['ID']+'.json')).read_bytes()==raw
  raw=ev['Inspections']['after_run']['Tables']
  counts={'RequestsAdmitted':0,'RequestsConsumed':0,'ResultsAdmitted':0,'ResultsConsumed':0,'Unpaid':[]}
  paid=set()
  for row in raw['observations']:
   obs=json.loads(row['descriptor'])
   admission=next(a for a in raw['admissions'] if a['observation_id']==obs['ID'])
   iid=json.loads(admission['descriptor'])['IntentID']
   origin=next(json.loads(i['descriptor']) for i in raw['intents'] if i['id']==iid)
   assert origin['Run']['Generation']==obs['SourceGeneration']
   assert origin['Generation']==obs['Run']['Generation']
   is_req='request_id' in obs['Payload']
   key='Requests' if is_req else 'Results'
   counts[key+'Admitted']+=1
   consumed=row['consumed_sequence'] is not None
   selected=[h for h in raw['history'] if h['selected_observation']==obs['ID']]
   assert len(selected)==int(consumed)
   if consumed:
    assert selected[0]['sequence']==row['consumed_sequence']
    assert selected[0]['subject']==obs['Run']['Subject']
    assert selected[0]['generation']==obs['Run']['Generation']
    counts[key+'Consumed']+=1
    if not is_req: paid.add(obs['Payload']['response_id']['Text'])
  counts['Unpaid']=[v for v in ['a-request','b-request'] if v not in paid]
  assert counts==ev['Counts'],(path,counts,ev['Counts'])
  assert len(ev['Prefix'])==4
  assert not raw['terminal_obligations']
  states=[json.loads(r['projection']) for r in raw['runs']]
  assert all(s['Status']=='open' for s in states)
  serial=ev['Profile']=='serial_serial'
  assert counts['ResultsConsumed']==(0 if serial else 2)
  assert counts['RequestsConsumed']==(0 if serial else 2)
  first=ev['Reports']['first']
  assert first['Quiescent'] and not first['UnresolvedIntents'] and not first['PendingDeliveries']
  if serial:
   assert len(ev['Refusals'])==2 and all(x['Error']=='ineligible observation' for x in ev['Refusals'])
   reopen=ev['Inspections']['after_reopen']['Tables']
   assert [r['epoch'] for r in reopen['runs']]==[r['epoch']+1 for r in raw['runs']]
   for tab in ('observations','history','intents','admissions','programs','subjects','message_acknowledgments'):
    assert reopen[tab]==raw[tab],(path,tab)
   assert all(r['Steps']==0 and r['Dispatches']==0 for r in ev['Reports'].values())
  else:
   assert all(s['Data']['validated']['Flag'] for s in states)
  cases.append({'Mode':mode,'Case':path.parent.name,'RawCounts':counts,'PrefixCommits':len(ev['Prefix']),
                'HistoryRows':len(raw['history']),'MessageAcknowledgments':len(raw['message_acknowledgments']),
                'RunReports':{k:{f:v[f] for f in ('Steps','Dispatches','Quiescent')} for k,v in ev['Reports'].items()},
                'EvidenceSHA256':hashlib.sha256(path.read_bytes()).hexdigest()})
result={'CasesChecked':len(cases),'SerialCases':sum(c['Case'].startswith('serial_serial-') for c in cases),
        'ServicingCases':sum(not c['Case'].startswith('serial_serial-') for c in cases),'Cases':cases}
output.write_text(json.dumps(result,indent=2)+'\n')
print(json.dumps({k:v for k,v in result.items() if k!='Cases'}))
