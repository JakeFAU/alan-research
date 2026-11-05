"""Package already completed immutable qualification outputs; never opens SQLite."""
import collections,gzip,hashlib,io,json,pathlib,statistics,tarfile,shutil,subprocess,sys
repo=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research');out=repo/'docs/foundational/phase3b4/evidence/task5';work=pathlib.Path('/tmp/alan-task5')
def sha(b):return hashlib.sha256(b).hexdigest()
def write(name,data):(out/name).write_text(json.dumps(data,indent=2)+'\n')
def archive(name,root,filter=None):
 with tarfile.open(out/name,'w:gz') as t:
  for p in sorted(root.iterdir()):t.add(p,arcname=p.name,filter=filter)
def dedup(name,paths):
 blobs={};manifest={}
 for prefix,path in paths:
  files=[path] if path.is_file() else sorted(path.rglob('*'))
  for p in files:
   if not p.is_file() or p.is_symlink() or p.name=='checker.test':continue
   rel=prefix if p==path else prefix+'/'+str(p.relative_to(path));b=p.read_bytes();h=sha(b);manifest[rel]=h;blobs.setdefault(h,b)
 with tarfile.open(out/name,'w:gz') as t:
  for path,b in [('files.json',json.dumps(manifest,indent=2).encode()+b'\n')]+[('sha256/'+h,b) for h,b in sorted(blobs.items())]:
   info=tarfile.TarInfo(path);info.size=len(b);info.mode=0o644;t.addfile(info,io.BytesIO(b))
 return {'files':len(manifest),'unique_blobs':len(blobs),'reconstruct':'files.json maps each original relative path to sha256/<digest>; copy bytes verbatim'}
for kind in ['reciprocal','authority']:
 for mode in ['normal','race']:
  archive('prior-'+kind+'-'+mode+'.tar.gz',pathlib.Path('/tmp/alan-task5-prior-'+kind+'-'+mode))
summary={};indices={};costs=[]
for mode in ['normal','race']:
 root=pathlib.Path('/tmp/alan-task5-final-'+mode)
 entries=[]
 for p in sorted(root.glob('*/case.json')):
  case=json.loads(p.read_text());d=p.parent
  entries.append({'Name':case['Name'],'Directory':d.name,'OriginalDirectory':str(d),'CaseSHA256':sha(p.read_bytes()),'Profile':case['Profile'],'PrincipalRequests':case['PrincipalRequests'],'AuxiliaryRequests':case['AuxiliaryRequests'],'Captures':case['Captures'],'Files':{str(f.relative_to(d)):sha(f.read_bytes()) for f in sorted(d.rglob('*')) if f.is_file()}})
  if (d/'measurements.json').exists():costs.append({'Mode':mode,'Case':case['Name'],'Measurements':json.loads((d/'measurements.json').read_text()),'Sidecar':json.loads((d/'final/cost-sidecar.json').read_text()),'QueryPlan':json.loads((d/'query-plan.json').read_text())})
 indices[mode]=entries;archive('final-'+mode+'.tar.gz',root)
 archive('final-promoted-'+mode+'.tar.gz',pathlib.Path('/tmp/alan-task5-final-promoted-'+mode))
 archive('final-v1-'+mode+'.tar.gz',pathlib.Path('/tmp/alan-task5-final-v1-'+mode))
 archive('independent-'+mode+'.tar.gz',work/('final-'+mode+'-verdicts'))
 archive('attacks-'+mode+'.tar.gz',work/('final-'+mode+'-attacks'))
 rs=[json.loads(p.read_text()) for p in (work/('final-'+mode+'-verdicts')).glob('*.json') if '-promotion' not in p.name]
 summary[mode]={'Cases':len(rs),'RawEvaluated':sum(bool(r['Phases']) for r in rs),'RawValidFinal':sum(r['Final']['Valid'] for r in rs),'Captures':sum(len(r['Phases'] or []) for r in rs),'AllCapturesRawValid':all(p['Verdict']['Valid'] for r in rs for p in r['Phases'] or []),'PrincipalComplete':sum(r['Final']['Custody']['Complete'] for r in rs),'SupportedModelProjection':sum(r['Projection']['Supported'] for r in rs),'Contacts':sum(r['External']['Contacts'] for r in rs),'Applications':sum(r['External']['Applications'] for r in rs),'AcknowledgedLosses':sum(r['External']['AcknowledgedLosses'] for r in rs),'AllEvaluatedExternalValid':all(r['External']['Valid'] for r in rs if r['Phases']),'AllContinuityValid':all(r['ContinuityFailure'] is None for r in rs),'UnsupportedReasons':dict(collections.Counter(reason for r in rs for reason in r['Projection']['Reasons'] or []))}
write('case-index.json',indices);write('verdict-summary.json',summary);write('cost-observations.json',costs)
cs={}
for mode in ['normal','race']:
 cs[mode]={}
 for n in [0,1,8,32]:
  group=[c for c in costs if c['Mode']==mode and c['Measurements']['N']==n];assert len(group)==5
  ns={key:[c['Measurements']['Nanoseconds'][key] for c in group] for key in group[0]['Measurements']['Nanoseconds']}
  cs[mode][str(n)]={'Observations':5,'Nanoseconds':{key:{'min':min(v),'median':statistics.median(v),'max':max(v)} for key,v in ns.items()},'CommittedReductions':sorted(set(c['Sidecar']['CommittedReductions'] for c in group)),'MaxProgramBytes':max(c['Sidecar']['MaxProgramBytes'] for c in group),'AllExactIndexLookup':all(c['QueryPlan']==['SEARCH admissions USING INDEX admissions_observation_id (observation_id=?)'] for c in group),'Allocations':'not measured','TotalConstantTimeClaim':False}
write('cost-summary.json',cs)
write('mutant-archive-format.json',dedup('checker-mutants.tar.gz',[('final-mutants',work/'final-mutants')]))
# Preserve preliminary source variants and failures by content address without
# duplicating unchanged repository files in every isolated mutant tree.
paths=[]
for p in sorted(work.iterdir()):
 if p.name.startswith('preliminary-') and p.name!='preliminary-selected' and not p.name.endswith('.tar.gz'):paths.append((p.name,p))
 elif p.name=='tooling' or p.suffix in ('.log','.py','.diff','.json'):paths.append((p.name,p))
write('preliminary-archive-format.json',dedup('preliminary-apparatus.tar.gz',paths))
for p in (work/'final-commands').iterdir():shutil.copy(p,out/'commands')
if '--skip-parent' not in sys.argv:
 shutil.copytree('/tmp/alan-phase3b4-parent-final-gates',out/'parent-final-gates',dirs_exist_ok=True)
probe=work/'runtime-filter-probe';(out/'runtime-filter-comparison').mkdir(exist_ok=True)
for p in [probe/'probe-provenance.json',probe/'result.json',probe/'result.log',probe/'internal/executor/phase3b4_independent_filter_test.go',work/'run-runtime-filter.py']:shutil.copy(p,out/'runtime-filter-comparison')
after={}
for label,path in [('research',repo),('runtime',pathlib.Path('/private/tmp/alan-phase3b3-runtime'))]:
 frozen=json.loads((work/'freeze'/(label+'-source.json')).read_text());current={p:sha((path/p).read_bytes()) for p in frozen};assert current==frozen
 after[label]={'AllFrozenFilesUnchanged':True,'Files':len(current),'CurrentHead':subprocess.check_output(['git','rev-parse','HEAD'],cwd=path).decode().strip()}
old=json.loads((work/'freeze/old-exports.json').read_text())
for p,h in old.items():
 b=(repo/p).read_bytes();assert sha(b)==h['sha256'];assert sha(gzip.decompress(b) if p.endswith('.gz') else b)==h['decoded_sha256']
write('source-after.json',after)
write('artifacts.json',{str(p.relative_to(out)):sha(p.read_bytes()) for p in sorted(out.rglob('*')) if p.is_file() and p.name!='artifacts.json'})
print(json.dumps(summary,indent=2))
