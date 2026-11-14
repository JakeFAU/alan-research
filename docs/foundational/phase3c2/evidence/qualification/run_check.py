import datetime, hashlib, json, os, pathlib, subprocess, sys, time
root=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research')
scratch=pathlib.Path('/tmp/alan-phase3c2/qualification')
commands={'build':['go','build','./...'],'package-list':['go','list','./...'],'module-verify':['go','mod','verify'],'module-tidy-diff':['go','mod','tidy','-diff']}
commands['build-cache-writable']=commands['build']
commands['package-list-cache-writable']=commands['package-list']
name=sys.argv[1];argv=commands[name]
out=scratch/name;out.mkdir(parents=True,exist_ok=True)
env=os.environ.copy();extra={'GOCACHE':'/tmp/alan-phase3c2/gocache','STATICCHECK_CACHE':'/tmp/alan-phase3c2/staticcheck-cache'};env.update(extra)
head=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip()
job={'argv':argv,'cwd':str(root),'head':head,'environment':extra,'started_utc':datetime.datetime.now(datetime.timezone.utc).isoformat()}
(out/'job.json').write_text(json.dumps(job,indent=2)+'\n')
start=time.monotonic()
with (out/'output.log').open('wb') as f:
 result=subprocess.run(argv,cwd=root,env=env,stdout=f,stderr=subprocess.STDOUT)
record={'exit_code':result.returncode,'elapsed_seconds':round(time.monotonic()-start,3),'output_bytes':(out/'output.log').stat().st_size,'output_sha256':hashlib.sha256((out/'output.log').read_bytes()).hexdigest()}
(out/'result.json').write_text(json.dumps(record,indent=2)+'\n')
print(json.dumps({'check':name,**record}),flush=True)
sys.exit(result.returncode)
