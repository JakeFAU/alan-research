import concurrent.futures,gzip,hashlib,json,os,pathlib,subprocess,time
root=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research');out=pathlib.Path('/tmp/alan-task5/final-exports');out.mkdir(exist_ok=False)
configs=[('twoagent','research/twoagent/cmd/export','docs/foundational/phase3a/evidence/model.json'),('mutual','research/mutualobligations/cmd/export','docs/foundational/phase3b/evidence/model.json'),('authority','research/authorityservice/cmd/export','docs/foundational/phase3b2/evidence/model.json.gz'),('audit','research/authorityservice/cmd/audit','docs/foundational/phase3b2/evidence/audit.json.gz')]
def work(config):
 name,pkg,old=config;binary=out/name;record={'name':name,'commands':[]};archived=(root/old).read_bytes();decoded=gzip.decompress(archived) if old.endswith('.gz') else archived
 record['archived_path']=old;record['archived_sha256']=hashlib.sha256(archived).hexdigest();record['decoded_sha256']=hashlib.sha256(decoded).hexdigest()
 def run(argv,label):
  log=out/(label+'.log');start=time.monotonic()
  with log.open('wb') as stream:r=subprocess.run(argv,cwd=root,stdout=stream,stderr=subprocess.STDOUT,timeout=3500)
  record['commands'].append({'argv':argv,'cwd':str(root),'exit':r.returncode,'elapsed_seconds':time.monotonic()-start,'output':str(log),'sha256':hashlib.sha256(log.read_bytes()).hexdigest()});print(name,label,r.returncode,flush=True)
  if r.returncode:raise RuntimeError(label)
  return log
 run(['go','build','-buildvcs=false','-o',str(binary),'./'+pkg],name+'-build')
 for number in (1,2):
  log=run([str(binary)],name+'-'+str(number));assert log.read_bytes()==decoded,(name,number,'decoded mismatch')
 record['twice_byte_equal_to_archive']=True;assert (root/old).read_bytes()==archived
 (out/(name+'-comparison.json')).write_text(json.dumps(record,indent=2)+'\n');return record
with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:records=list(pool.map(work,configs))
(out/'results.json').write_text(json.dumps(records,indent=2)+'\n')
