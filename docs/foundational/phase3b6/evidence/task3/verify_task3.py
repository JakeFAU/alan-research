#!/usr/bin/env python3
"""Focused Task3 qualification, deterministic reports and frozen old-byte comparison."""
import gzip,hashlib,io,json,os,pathlib,subprocess,tarfile,time
root=pathlib.Path.cwd();evidence=root/'docs/foundational/phase3b6/evidence/task3';tmp=pathlib.Path('/tmp/alan-phase3b6/task3');tmp.mkdir(exist_ok=True)
env=os.environ.copy();env.update(GOTOOLCHAIN='local',GOWORK='off',GOPROXY='off',GOSUMDB='off',GOMODCACHE='/tmp/alan-phase3b5/modcache',GOCACHE='/tmp/alan-phase3b6/go-build',XDG_CACHE_HOME='/tmp/alan-phase3b6/cache',STATICCHECK_CACHE='/tmp/alan-phase3b6/staticcheck',ALAN_EFFECT_OBSERVATION_CORPUS='/tmp/alan-phase3b5/corpus',ALAN_REVOCATION_CORPUS='/tmp/alan-phase3b5/corpus')
commands=[];sha=lambda b:hashlib.sha256(b).hexdigest()
def run(args,label=None,cwd=root,want=0):
 start=time.monotonic();r=subprocess.run([str(a) for a in args],cwd=cwd,env=env,capture_output=True);record={'command':[str(a) for a in args],'cwd':str(cwd),'exit':r.returncode,'expected_exit':want,'seconds':time.monotonic()-start};commands.append(record)
 if label:(evidence/(label+'.log')).write_bytes(r.stdout+r.stderr)
 assert r.returncode==want,(args,r.stdout.decode(),r.stderr.decode())
 return r.stdout
run(['go','version'],'toolchain');run(['staticcheck','-version'],'staticcheck-version')
run(['python3',evidence/'verify_source.py'],'source-verification')
run(['go','test','-count=1','-v','./research/delegationcorrespondence','./research/delegationcorrespondence/cmd/effectobservation','-run','^TestEffectObservation'],'focused-tests')
run(['staticcheck','./research/delegationcorrespondence/...'],'staticcheck')
for command in ['effectobservation','revocation','check']:run(['go','build','-o',tmp/command,'./research/delegationcorrespondence/cmd/'+command],'build-'+command)
# Reconstruct old output from frozen Phase3B5 code, verify archived byte hashes,
# then compare bytes. The old individual report byte files were not archived.
base=tmp/'phase3b5-source';base.mkdir(exist_ok=True)
archive=run(['git','archive','df840103aa3a85e33b4b95ddcf2fe3df071cb024','go.mod','go.sum','research'])
with tarfile.open(fileobj=io.BytesIO(archive)) as t:
 for member in t.getmembers():assert not pathlib.PurePosixPath(member.name).is_absolute() and '..' not in pathlib.PurePosixPath(member.name).parts
 t.extractall(base)
for command in ['revocation','check']:run(['go','build','-o',tmp/('old-'+command),'./research/delegationcorrespondence/cmd/'+command],'build-old-'+command,cwd=base)
corpus=pathlib.Path(env['ALAN_EFFECT_OBSERVATION_CORPUS']);manifest=json.loads((corpus/'manifest.json').read_text())
old=json.loads((root/'docs/foundational/phase3b5/evidence/task3/verdicts.json').read_text());legacy=[]
for archived in old:
 mode,name=archived['mode'],archived['case'];path=corpus/('raw-'+mode)/manifest['raw-'+mode+'.tar.gz']['cases'][name];status=0 if archived['supported'] else 1
 a=run([tmp/'old-revocation','-case',path],want=status);b=run([tmp/'revocation','-case',path],want=status)
 c=run([tmp/'old-check','-binding-case',path]);d=run([tmp/'check','-binding-case',path])
 assert sha(a)==archived['report_sha256'] and sha(c)==archived['legacy_sha256'] and a==b and c==d
 legacy.append({'mode':mode,'case':name,'report_sha256':sha(b),'raw_legacy_sha256':sha(d),'frozen_baseline_and_archived_hash_match':True})
(evidence/'legacy-comparison.json').write_text(json.dumps(legacy,indent=2)+'\n')
exports=evidence/'reports';exports.mkdir(exist_ok=True);summary=[]
for mode in ['normal','race']:
 for name in ['N07','applied-unknown-revoked-before-fresh-gate']:
  path=corpus/('raw-'+mode)/manifest['raw-'+mode+'.tar.gz']['cases'][name]
  a=run([tmp/'effectobservation','-case',path]);b=run([tmp/'effectobservation','-case',path]);assert a==b
  r=json.loads(a);assert r['Supported'] and r['Producer']['SealedBytesMatch']
  output=mode+('-loss' if name!='N07' else '-N07')+'.json.gz';(exports/output).write_bytes(gzip.compress(a,mtime=0))
  summary.append({'mode':mode,'case':name,'report_sha256':sha(a),'report_bytes':len(a),'file':output,'captures':[{'path':c['Path'],'sha256':c['SHA256'],'supported':c['Supported'],'reason':c['Reason'],'event_count':len(c['Witness'] or []),'store':c['Store'],'dossier':c['Dossier']} for c in r['Captures']]})
(evidence/'verdicts.json').write_text(json.dumps(summary,indent=2)+'\n')
files=sorted(root.glob('research/delegationcorrespondence/effect*.go'))+sorted((root/'research/delegationcorrespondence/cmd/effectobservation').glob('*.go'))+[root/'research/delegationcorrespondence/testdata/effect-observation/anchors.json']
fmt=run(['gofmt','-l',*[p for p in files if p.suffix=='.go']],'gofmt');assert not fmt
(evidence/'source-sha256.json').write_text(json.dumps({str(p.relative_to(root)):sha(p.read_bytes()) for p in files},indent=2)+'\n')
(evidence/'commands.json').write_text(json.dumps({'environment':{k:env[k] for k in ['GOTOOLCHAIN','GOWORK','GOPROXY','GOSUMDB','GOMODCACHE','GOCACHE','XDG_CACHE_HOME','STATICCHECK_CACHE','ALAN_EFFECT_OBSERVATION_CORPUS','ALAN_REVOCATION_CORPUS']},'commands':commands},indent=2)+'\n')
print('PASS focused tests, scoped Staticcheck,4 deterministic reports (16 captures),10 old report and10 raw report byte comparisons against frozen code and archived hashes.')
