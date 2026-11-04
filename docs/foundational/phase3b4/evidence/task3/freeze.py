from pathlib import Path
import subprocess,hashlib,json,os
root=Path('/private/tmp/alan-phase3b3-runtime');out=Path('/tmp/alan-phase3b4-task3')
files=subprocess.check_output(['git','ls-files'],cwd=root,text=True).splitlines()
source={name:hashlib.sha256((root/name).read_bytes()).hexdigest() for name in files if name.endswith('.go') or name in ['go.mod','go.sum']}
production={name:digest for name,digest in source.items() if not name.endswith('_test.go')}
aggregate=lambda x:hashlib.sha256(json.dumps(x,sort_keys=True,separators=(',',':')).encode()).hexdigest()
manifest=dict(base='3f8307664b6543331a41c943e91fb5bac8b376ef',head=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),commits=subprocess.check_output(['git','log','--format=%H %s','3f8307664b6543331a41c943e91fb5bac8b376ef..HEAD'],cwd=root,text=True).splitlines(),cwd=str(root),source_sha256=source,source_map_sha256=aggregate(source),production_source_map_sha256=aggregate(production),toolchain=subprocess.check_output(['go','version'],text=True).strip(),staticcheck=subprocess.check_output(['staticcheck','-version'],text=True).strip(),environment={'GOTOOLCHAIN':'local','GOPROXY':'off','GOSUMDB':'off','GOCACHE':'/tmp/alan-phase3a-go-build','GOWORK':'off','STATICCHECK_CACHE':'/tmp/alan-phase3b4-staticcheck'},old_reader={'path':'/tmp/alan-phase3b4-task1/compatibility/old-reader.test','sha256':hashlib.sha256(Path('/tmp/alan-phase3b4-task1/compatibility/old-reader.test').read_bytes()).hexdigest(),'source':'1fdd4c9e1720bf60b344914bcbc3d0938ba8d678'})
(out/'frozen-source.json').write_text(json.dumps(manifest,indent=2)+'\n')
print(json.dumps({k:v for k,v in manifest.items() if k!='source_sha256'},indent=2))
