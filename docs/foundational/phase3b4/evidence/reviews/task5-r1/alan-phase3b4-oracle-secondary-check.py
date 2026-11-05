import tarfile,pathlib,json,hashlib,gzip
P=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3b4/evidence/task5');R=P.parents[4]
sha=lambda b:hashlib.sha256(b).hexdigest()
freeze=json.loads((P/'freeze/research-source.json').read_text())
for f in (P/'commands').glob('*.json'):
 d=json.loads(f.read_text());assert d['source_sha256_before']==freeze,f.name
print('COMMAND_MAPS_EQUAL_FREEZE',6)
with tarfile.open(P/'checker-mutants.tar.gz') as t:
 fs=json.load(t.extractfile('files.json'))
 for n,h in fs.items():
  if n.endswith('/record.json'):
   d=json.load(t.extractfile('sha256/'+h));changed=[f for f,h in d['source_sha256'].items() if freeze[f]!=h];assert changed==[d['path']],changed;print('MUTANT_ONLY_CHANGED',d['name'],changed)
with tarfile.open(P/'unchanged-model-exports.tar.gz') as t:
 for n in t.getnames():
  if n.endswith('comparison.json'):print('OLD_EXPORT',n,json.load(t.extractfile(n)))
 print('EXPORT_COMMANDS',t.extractfile('results.json').read().decode()[:4000])
with tarfile.open(P/'final-normal.tar.gz') as t,tarfile.open(P/'attacks-normal.tar.gz') as attacks,tarfile.open(P/'final-promoted-normal.tar.gz') as promoted:
 def load(t,n):return json.load(t.extractfile(n))
 for n in t.getnames():
  if not n.endswith('/case.json'):continue
  c=load(t,n);base=n.rsplit('/',1)[0];name=c['Name']
  if name.startswith('N10c') or name.startswith('N10d'):
   raw=load(t,base+'/final/inspect.json');print('LIFECYCLE',name)
   for row in raw['Tables']['composition_facts']:
    f=json.loads(row['descriptor'])
    if f['Kind'] in ['authorization','lifecycle','reservation','authority_admission']:print('FACT',f['Order'],f['Kind'],f['OperationID'],f.get('AuthorityAdmission',{}).get('Reason'))
   for row in raw['Tables']['lifecycle_results']:
    l=json.loads(row['result']);print('LIFE',l['Request']['Operation'],l['Accepted'],l['Reason'])
  if name=='M03-auxiliary-target-generation-two':
   for typ,raw in [('parent',load(t,base+'/final/inspect.json')),('attack',load(attacks,'M03.json'))]:
    for r in raw['Tables']['admissions']:
     a=json.loads(r['descriptor'])
     if a['IntentID']=='["b",1,6,0]':print('M03_ADMISSION',typ,a['Generation'],json.loads(r['result']))
  if name=='auxiliary-serviced-equal-total':
   for typ,raw in [('parent',load(t,base+'/final/inspect.json')),('attack',load(attacks,'equal-total-custody.json'))]:
    print('EQUAL_TOTAL',typ,'receipts',[(r['intent_id'],json.loads(r['descriptor'])['Mutation'],json.loads(r['descriptor'])['Total']) for r in raw['Tables']['external_receipts']])
    for r in raw['Tables']['intents']:
     if r['id']=='["a",1,16,0]':print('REPLY',typ,json.loads(r['descriptor'])['Payload'])
  if name in ['K01','K02']:
   root=next(x.rsplit('/',1)[0] for x in promoted.getnames() if x.startswith(name+'-') and x.endswith('parent-inspect.json'))
   assert t.extractfile(base+'/final/inspect.json').read()==promoted.extractfile(root+'/parent-inspect.json').read()
   before=load(promoted,root+'/opened-inspect.json');after=load(promoted,root+'/promoted-inspect.json')
   for k,v in before['Tables'].items():
    if k=='metadata':continue
    assert v==after['Tables'][k],(name,k)
   print('PROMOTION',name,'only metadata/additional receipt_conflicts changed',len(after['Tables']['receipt_conflicts']))
