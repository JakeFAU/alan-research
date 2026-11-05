import os,pathlib,subprocess,time,json,hashlib
root=pathlib.Path('/tmp/alan-task5/runtime-filter-probe');env=dict(os.environ);env['ALAN_FILTER_PARENT']='/tmp/alan-task5/final-normal-attacks/A02-filter-parent.json';env['ALAN_FILTER_ATTACK']='/tmp/alan-task5/final-normal-attacks/A02-filter.json'
argv=['go','test','-v','-count=1','-timeout=3m','./internal/executor','-run','^TestPhase3B4IndependentEmptyOperationFilter$'];start=time.monotonic()
log=root/'result.log'
with log.open('wb') as f:r=subprocess.run(argv,cwd=root,env=env,stdout=f,stderr=subprocess.STDOUT,timeout=240)
record={'cwd':str(root),'argv':argv,'exit':r.returncode,'elapsed_seconds':time.monotonic()-start,'environment':{k:v for k,v in env.items() if k.startswith('ALAN_FILTER_')},'inputs':{k:hashlib.sha256(pathlib.Path(v).read_bytes()).hexdigest() for k,v in env.items() if k.startswith('ALAN_FILTER_')},'log_sha256':hashlib.sha256(log.read_bytes()).hexdigest()}
(root/'result.json').write_text(json.dumps(record,indent=2)+'\n');print(log.read_text());raise SystemExit(r.returncode)
