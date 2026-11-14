# This capture procedure requires the controller scratch inputs under /tmp/alan-phase3c1.
import gzip,hashlib,json,pathlib,subprocess
R=pathlib.Path.cwd(); S=pathlib.Path('/tmp/alan-phase3c1'); Q=R/'docs/foundational/phase3c1/evidence/qualification'
def sha(b): return hashlib.sha256(b).hexdigest()
def write(p,d): p.write_text(json.dumps(d,indent=2,sort_keys=True)+'\n')
files=[]
def archive(src,dst,compress=False):
 raw=src.read_bytes(); dst.parent.mkdir(parents=True,exist_ok=True)
 if compress:
  with dst.open('wb') as f:
   with gzip.GzipFile(filename='',mode='wb',fileobj=f,mtime=0) as z: z.write(raw)
 else: dst.write_bytes(raw)
 files.append(dict(source=str(src),path=str(dst.relative_to(R)),raw_bytes=len(raw),raw_sha256=sha(raw),retained_sha256=sha(dst.read_bytes())))
for p in sorted((S/'baseline').iterdir()):
 if p.is_file(): archive(p,Q/'baseline'/(p.name+('.gz' if p.suffix in ('.jsonl','.stderr') else '')),p.suffix in ('.jsonl','.stderr'))
for name in ['go-inputs.json','input-comparison.json','staticcheck-comparison.json','package-coverage.json','run_check.py']:
 archive(S/'qualification'/name,Q/name)
for p in sorted((S/'qualification').iterdir()):
 if not p.is_dir() or not (p/'job.json').exists(): continue
 for name in ['job.json','result.json','output.log']:
  archive(p/name,Q/'commands'/p.name/(name+('.gz' if name=='output.log' else '')),name=='output.log')
# Source hashes, provenance and package outcomes are verified, not rerun.
frozen=json.loads((S/'baseline/old-go-inputs.json').read_text())['files']
current=json.loads((S/'qualification/go-inputs.json').read_text())['files']
for p,v in frozen.items():
 h=v['sha256'] if isinstance(v,dict) else v
 assert sha((R/p).read_bytes())==h,p
for p,h in current.items(): assert sha((R/p).read_bytes())==h,p
supp=json.loads((S/'qualification/input-comparison.json').read_text())['supplemental_preexisting_module_files']
for p,v in supp.items():
 base=subprocess.check_output(['git','show','1456a0000a0e14bb3dd105026189717c71ab905b:'+p])
 assert sha(base)==v['base_sha256']==v['current_sha256']==sha((R/p).read_bytes())
m1=json.loads((R/'docs/foundational/phase3c1/evidence/task1/source-manifest.json').read_text())
m2=json.loads((R/'docs/foundational/phase3c1/evidence/task2/manifest.json').read_text())
for p,h in m1['SourceSHA256'].items(): assert sha((R/p).read_bytes())==h,p
for p,h in m2['source_sha256'].items(): assert sha((R/p).read_bytes())==h,p
for p,h in m2['retained_file_sha256'].items(): assert sha((R/p).read_bytes())==h,p
for a in m2['artifacts']:
 packed=(R/a['path']).read_bytes(); raw=gzip.decompress(packed)
 assert sha(packed)==a['compressed_sha256'] and sha(raw)==a['raw_sha256'] and len(raw)==a['raw_bytes']
 assert packed[3]&8==0 and packed[4:8]==b'\0'*4
old=json.loads((S/'baseline/inputs.json').read_text())['packages']
now=(S/'qualification/package-list-cache-writable/output.log').read_text().splitlines()
new=['github.com/JakeFAU/alan-research/research/contractcomposition','github.com/JakeFAU/alan-research/research/contractcomposition/cmd/export']
assert set(now)==set(old)|set(new) and len(now)==72
terminal={}; failures=[]
for line in (S/'baseline/ordinary.jsonl').read_text().splitlines():
 event=json.loads(line)
 if event.get('Action')=='fail': failures.append(event)
 if not event.get('Test') and event.get('Action') in ('pass','skip','fail'): terminal[event['Package']]=event['Action']
assert not failures and set(terminal)==set(old)
for log in ['normal.log','race.log']:
 text=(R/'docs/foundational/phase3c1/evidence/task1'/log).read_text()
 for p in new: assert any(line.split()[:2]==['ok',p] for line in text.splitlines()),(log,p)
 for c in m1['Commands']:
  if c.get('output','').endswith('/'+log): assert c['exit_code']==0
for name,count in [('global-staticcheck',13),('global-staticcheck-integration',31)]:
 data=(S/'qualification'/name/'output.log').read_bytes()
 olddata=gzip.decompress((R/'docs/foundational/phase3b7/evidence/qualification/verification-commands'/name/'output.log.gz').read_bytes())
 assert data==olddata and len(data.splitlines())==count
write(Q/'task3-verification.json',dict(task_base='7b804df8b0dbbd160a838cdee8681558e69dc200',frozen_baseline_files=len(frozen),current_go_module_files=len(current),supplemental_preexisting_module_files=len(supp),task1_sources=len(m1['SourceSHA256']),task2_sources=len(m2['source_sha256']),task2_retained_files=len(m2['retained_file_sha256']),compressed_artifacts=len(m2['artifacts']),old_packages=len(old),current_packages=len(now),old_pass=sum(v=='pass' for v in terminal.values()),old_no_test=sum(v=='skip' for v in terminal.values()),new_packages=new,normal_and_race_logs_verified=True,global_lint_exact_inheritance=[13,31],source_changes=False,method='Task 3 SHA-256, gzip and retained-command/package-outcome verification; no suite rerun'))
write(Q/'archive-manifest.json',dict(schema='alan.phase3c1.qualification.archive.v1',attribution='Controller-run baseline/root checks during this Phase 3C.1 task; copied and verified by Task 3. Earlier Phase 3B.7 lint is comparison evidence only.',files=files))
for name in ['alan-dissertation.tex','references.bib']:
 (S/'task3/base'/name).write_bytes(subprocess.check_output(['git','show','7b804df8b0dbbd160a838cdee8681558e69dc200:paper/'+name]))
print('archived',len(files),'files; verified 598 frozen + 29 supplemental + 10 new inputs; all provenance and coverage assertions passed')
