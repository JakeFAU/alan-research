from pathlib import Path
import subprocess,tarfile,io,hashlib,json,os
root=Path('/private/tmp/alan-phase3b3-runtime');out=Path('/tmp/alan-phase3b4-task3');copy=out/'red-reproduction';copy.mkdir(exist_ok=True)
archive=subprocess.check_output(['git','archive','HEAD'],cwd=root)
with tarfile.open(fileobj=io.BytesIO(archive)) as tar:tar.extractall(copy,filter='data')
changed=subprocess.check_output(['git','diff','--name-only'],cwd=root,text=True).splitlines()+subprocess.check_output(['git','ls-files','--others','--exclude-standard'],cwd=root,text=True).splitlines()
for name in changed:
 p=root/name
 if p.is_file():dest=copy/name;dest.parent.mkdir(parents=True,exist_ok=True);dest.write_bytes(p.read_bytes())
env=os.environ.copy();env.update(GOTOOLCHAIN='local',GOPROXY='off',GOSUMDB='off',GOCACHE='/tmp/alan-phase3a-go-build',GOWORK='off')
base=(copy/'internal/executor/request_authority.go').read_bytes();binding=(copy/'internal/executor/request_binding.go').read_bytes()
cases=[('profile', 'TestRequestBindingProfile', 'internal/executor/request_authority.go', 'return v == RequestAuthorityVersion || v == RequestAuthorityBindingVersion','return v == RequestAuthorityVersion'),('selection','TestRequestBindingNonadjacent','internal/executor/request_authority.go','if x.profile.RequestAuthority == RequestAuthorityBindingVersion {','if false {'),('commit-join','TestRequestBindingOriginalCommitJoin','internal/executor/request_binding.go','if err := validateBindingCommit(ctx, q, intent); err != nil {','if err := validateBindingCommit(ctx, q, intent); false && err != nil {')]
records=[]
for name,test,file,old,new in cases:
 (copy/'internal/executor/request_authority.go').write_bytes(base);(copy/'internal/executor/request_binding.go').write_bytes(binding)
 p=copy/file;s=p.read_text();assert s.count(old)==1;(out/(name+'-reproduction-original.go.txt')).write_text(s);p.write_text(s.replace(old,new));(out/(name+'-reproduction-mutant.go.txt')).write_bytes(p.read_bytes())
 files={str(p.relative_to(copy)):hashlib.sha256(p.read_bytes()).hexdigest() for p in sorted(copy.rglob('*.go'))}
 argv=['go','test','./internal/executor','-run','^'+test+'$','-count=1'];r=subprocess.run(argv,cwd=copy,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
 log=name+'-red-reproduced.log';(out/log).write_bytes(r.stdout)
 records.append(dict(name=name,cwd=str(copy),argv=argv,exit=r.returncode,mutation=dict(file=file,old=old,new=new),go_source_sha256=files,output=log))
 print(name,r.returncode,r.stdout.decode())
 assert r.returncode==1
(out/'red-reproduction.json').write_text(json.dumps(records,indent=2)+'\n')
