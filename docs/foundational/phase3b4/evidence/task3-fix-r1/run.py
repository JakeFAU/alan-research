import subprocess,os,json,pathlib,hashlib,sys,time
out=pathlib.Path(__file__).parent
name,cwd,*cmd=sys.argv[1:]
env=os.environ.copy()
for k in list(env):
 if k.startswith('ALAN_'): del env[k]
env.update(GOTOOLCHAIN='local',GOPROXY='off',GOSUMDB='off',GOCACHE='/tmp/alan-phase3a-go-build',GOWORK='off',STATICCHECK_CACHE='/tmp/alan-phase3b4-staticcheck')
root=pathlib.Path(cwd)
files={str(p.relative_to(root)):hashlib.sha256(p.read_bytes()).hexdigest() for p in sorted(root.rglob('*.go')) if '.git' not in p.parts}
for n in ['go.mod','go.sum']: files[n]=hashlib.sha256((root/n).read_bytes()).hexdigest()
start=time.time(); r=subprocess.run(cmd,cwd=cwd,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
(out/(name+'.log')).write_bytes(r.stdout)
(out/(name+'.json')).write_text(json.dumps({'command':cmd,'cwd':cwd,'environment':{k:env[k] for k in ['GOTOOLCHAIN','GOPROXY','GOSUMDB','GOCACHE','GOWORK','STATICCHECK_CACHE']},'exit_code':r.returncode,'elapsed_seconds':time.time()-start,'source_sha256':files,'output_sha256':hashlib.sha256(r.stdout).hexdigest()},indent=2)+'\n')
print(name,'exit',r.returncode);print(r.stdout.decode())
