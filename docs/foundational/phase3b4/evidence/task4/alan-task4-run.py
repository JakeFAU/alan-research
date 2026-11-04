import os,sys,json,subprocess,time,hashlib,pathlib
root=pathlib.Path('/private/tmp/alan-phase3b3-runtime')
output=pathlib.Path('/tmp/alan-task4-qualification');output.mkdir(exist_ok=True)
name=sys.argv[1];argv=sys.argv[2:]
env=dict(os.environ);env.update(GOTOOLCHAIN='local',GOWORK='off',GOPROXY='off',GOSUMDB='off',GOMODCACHE='/tmp/alan-phase3b4-qualification-modcache',GOCACHE='/tmp/alan-phase3a-go-build',STATICCHECK_CACHE='/tmp/alan-phase3b4-staticcheck')
source={str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for pattern in ('**/*.go','go.mod','go.sum') for p in root.glob(pattern) if '.git' not in p.parts}
start=time.time()
with (output/(name+'.log')).open('wb') as log: result=subprocess.run(argv,cwd=root,env=env,stdout=log,stderr=subprocess.STDOUT)
record={'cwd':str(root),'argv':argv,'environment':{k:v for k,v in env.items() if k.startswith('ALAN_PHASE3B') or k in ('GOTOOLCHAIN','GOWORK','GOPROXY','GOSUMDB','GOMODCACHE','GOCACHE','STATICCHECK_CACHE')},'exit':result.returncode,'elapsed_seconds':time.time()-start,'source_sha256':source,'log':name+'.log'}
(output/(name+'.json')).write_text(json.dumps(record,indent=2)+'\n');print(name,result.returncode,record['elapsed_seconds']);sys.exit(result.returncode)
