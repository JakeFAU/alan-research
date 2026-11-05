import pathlib,json,tarfile,hashlib
P=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3b4/evidence/task5')
sha=lambda x:hashlib.sha256(x).hexdigest()
for f in ['research','runtime']:
 m=json.loads((P/'freeze'/f'{f}-source.json').read_text())
 print('MAPSHAPE',f,type(m),list(m)[:3])
 with tarfile.open(P/'freeze'/f'{f}-source.tar.gz') as t:
  for n,h in m.items(): assert sha(t.extractfile(n).read())==h,(f,n)
 print('FROZEN_ARCHIVE',f,len(m))
for f in sorted((P/'commands').glob('*.json')):
 d=json.loads(f.read_text());assert d['source_sha256_before']==d['source_sha256_after'];assert sha((f.parent/d['log']).read_bytes())==d['log_sha256'];assert d['exit']==0
 print('COMMAND',f.name,'guard/log exact',len(d['source_sha256_before']))
for mode in ['normal','race']:
 with tarfile.open(P/f'final-{mode}.tar.gz') as t:
  total=0
  for n in t.getnames():
   if n.endswith('/digests.json'):
    root=n.rsplit('/',1)[0];ds=json.load(t.extractfile(n))
    for f,h in ds.items():assert sha(t.extractfile(root+'/'+f).read())==h,(mode,root,f)
    total+=len(ds)
  print('CORPUS_DIGESTS',mode,total)
with tarfile.open(P/'checker-mutants.tar.gz') as t:
 fs=json.load(t.extractfile('files.json'))
 for h in set(fs.values()):assert sha(t.extractfile('sha256/'+h).read())==h
 for n,h in fs.items():
  if n.endswith('/record.json'):
   d=json.load(t.extractfile('sha256/'+h));assert d['built'] and d['qualified_kill'];base=n.rsplit('/',1)[0]
   for f,h in d['source_sha256'].items():assert fs[base+'/source/'+f]==h
   for c in d['commands']:
    log=base+'/'+pathlib.Path(c['log']).name;assert fs[log]==c['log_sha256']
   print('MUTANT',d['name'],[c['exit'] for c in d['commands']],d['semantic_kills'])
 print('MUTANT_BLOBS',len(fs),len(set(fs.values())))
def val(row,k):return json.loads(row[k])
def j(t,n):return json.load(t.extractfile(n))
def summarize(raw):
 T=raw['Tables'];out={}
 intents=[val(r,'descriptor') for r in T['intents']]
 obs=[(val(r,'descriptor'),r['consumed_sequence']) for r in T['observations']]
 out['messages']=[(i['ID'],i['Generation'],i['Obligation']) for i in intents if i['Kind']=='message' and i['ID'] in ['["b",1,1,0]','["b",1,6,0]']]
 out['selections']=[(o['ID'],c) for o,c in obs if o['Provenance']=='ordinary-message' and o['Source']=='b']
 out['services']=[(val(r,'descriptor')['Intent']['ID'],val(r,'descriptor')['Allowed'],val(r,'descriptor')['Reason'],val(r,'descriptor')['Authority']) for r in T['external_attempts']]
 out['order']=[(val(r,'descriptor')['Order'],val(r,'descriptor')['Kind'],val(r,'descriptor')['OperationID'],val(r,'descriptor').get('AuthorityAdmission')) for r in T['composition_facts'] if val(r,'descriptor')['Kind'] in ['lifecycle','reservation','authority_admission']]
 return out
with tarfile.open(P/'final-normal.tar.gz') as t:
 for n in t.getnames():
  if n.endswith('/case.json'):
   c=j(t,n);base=n.rsplit('/',1)[0]
   if c['Name'] in ['P10','N01','N02','N06','M03-auxiliary-target-generation-two','auxiliary-serviced-equal-total']:
    raw=j(t,base+'/final/inspect.json');print('RAW',c['Name'],json.dumps(summarize(raw),separators=(',',':')))
   if c['Name']=='C07':
    for cap in c['Captures']:
     raw=j(t,base+'/'+cap['Directory']+'/inspect.json');print('C07',cap['Directory'],cap['OwnerPID'],[(r['subject'],r['epoch'],r['sequence']) for r in raw['Tables']['runs']])
    for loss in ['loss-01.json','loss-02.json']:print('C07LOSS',loss,j(t,base+'/'+loss))
   if c['Name'] in ['K01','K02']:
    raw=j(t,base+'/final/inspect.json');T=raw['Tables'];print('COMPAT',c['Name'],'format',T['metadata'][0]['schema_version'],'programs',[(r['id'],r.get('semantics')) for r in T['programs']],'attempts',T['adapter_attempts'],'results',T['adapter_results']);print('COMPAT_FILES',[x[len(base)+1:] for x in t.getnames() if x.startswith(base+'/') and x.endswith('.json')])
