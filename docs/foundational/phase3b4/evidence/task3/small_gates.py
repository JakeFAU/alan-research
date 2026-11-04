from pathlib import Path
import subprocess,json,os,datetime
root=Path('/private/tmp/alan-phase3b3-runtime');out=Path('/tmp/alan-phase3b4-task3');source=json.loads((out/'frozen-source.json').read_text())
env=os.environ.copy();overrides={'GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off','GOCACHE':'/tmp/alan-phase3a-go-build','GOWORK':'off','STATICCHECK_CACHE':'/tmp/alan-phase3b4-staticcheck'};env.update(overrides)
commands=[('build',['go','build','./...']),('staticcheck',['staticcheck','./...']),('staticcheck-integration',['staticcheck','-tags','integration','./internal/executor','./internal/durabletest','./examples/delegation','./runtime']),('vet',['go','vet','./...']),('vet-integration',['go','vet','-tags','integration','./...']),('mod-verify',['go','mod','verify']),('mod-tidy',['go','mod','tidy','-diff']),('gofmt',['gofmt','-l','adapters','author','cmd','compose','effects','examples','internal','runtime','schema']),('diff-check',['git','diff','--check'])]
records=[]
for name,argv in commands:
 start=datetime.datetime.now(datetime.timezone.utc).isoformat();r=subprocess.run(argv,cwd=root,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT);log='frozen-'+name+'.log';(out/log).write_bytes(r.stdout)
 records.append(dict(argv=argv,cwd=str(root),environment=overrides,source_commit=source['head'],source_map_sha256=source['source_map_sha256'],start=start,exit=r.returncode,output=log));print(name,r.returncode)
 assert r.returncode==0
 if name in ['gofmt','mod-tidy']:assert not r.stdout
(out/'small-gates.json').write_text(json.dumps(records,indent=2)+'\n')
