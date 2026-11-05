#!/usr/bin/env python3
"""Record one qualification command without hiding its exit or source identity."""
import argparse, hashlib, json, os, pathlib, platform, signal, subprocess, sys, time
p=argparse.ArgumentParser()
p.add_argument('--repo',required=True)
p.add_argument('--out',required=True)
p.add_argument('--name',required=True)
p.add_argument('--expected-source')
p.add_argument('--timeout',type=int,default=3900)
p.add_argument('command',nargs=argparse.REMAINDER)
a=p.parse_args()
root=pathlib.Path(a.repo).resolve(); out=pathlib.Path(a.out); out.mkdir(parents=True,exist_ok=True)
argv=a.command[1:] if a.command[:1]==['--'] else a.command
if not argv: p.error('command required')
recordpath=out/(a.name+'.json'); logpath=out/(a.name+'.log')
if recordpath.exists() or logpath.exists(): raise SystemExit('Refusing to overwrite existing command evidence')
def git(*args):
 return subprocess.check_output(['git',*args],cwd=root).decode().strip()
def source():
 names=subprocess.check_output(['git','ls-files','-z','--','*.go','go.mod','go.sum','**/go.mod','**/go.sum','**/testdata/**','*.py','*.sh'],cwd=root).split(b'\0')
 return {n.decode():hashlib.sha256((root/n.decode()).read_bytes()).hexdigest() for n in names if n and n.decode().split('/')[0] not in ('docs','paper')}
before=source(); headbefore=git('rev-parse','HEAD')
if a.expected_source:
 expected=json.loads(pathlib.Path(a.expected_source).read_text())
 if before!=expected: raise SystemExit('Source differs from supplied frozen map; command not run')
env=dict(os.environ)
allowed=['GOTOOLCHAIN','GOWORK','GOPROXY','GOSUMDB','GOMODCACHE','GOCACHE','STATICCHECK_CACHE','CGO_ENABLED']
allowed += ['ALAN_TASK3_RAW','ALAN_TASK3_PROGRESS','ALAN_PHASE3B_RUNTIME_EXPORT','ALAN_PHASE2_EVIDENCE','ALAN_PHASE3B3_EVIDENCE','ALAN_PHASE3B3_OLD_READER','ALAN_PHASE3B4_EVIDENCE','ALAN_PHASE3B4_OLD_READER','ALAN_PHASE3B4_PUBLIC_PARENTS','ALAN_PHASE3B4_PROMOTED_OUTPUT','ALAN_PHASE3B4_CHECK_ROOT','ALAN_PHASE3B4_PROMOTED_CHECK_ROOT','ALAN_PHASE3B4_CHECK_OUTPUT','ALAN_PHASE3B4_ATTACK_OUTPUT']
selected={k:v for k,v in env.items() if k in allowed}
start=time.time(); timedout=False
with logpath.open('wb') as log:
 try:
  child=subprocess.Popen(argv,cwd=root,env=env,stdout=log,stderr=subprocess.STDOUT,start_new_session=True)
  code=child.wait(timeout=a.timeout)
 except subprocess.TimeoutExpired:
  code=124; timedout=True
  os.killpg(child.pid,signal.SIGKILL)
  child.wait()
  log.write(b'\nRecorder watchdog expired. This is not a semantic test verdict.\n')
after=source()
record={'cwd':str(root),'argv':argv,'environment':selected,'platform':platform.platform(),'machine':platform.machine(),'head_before':headbefore,'head_after':git('rev-parse','HEAD'),'source_scope':'Git-tracked Go/module, Python/shell qualification scripts and testdata files outside docs/paper; evidence and manuscript copies have separate manifests','source_sha256_before':before,'source_sha256_after':after,'source_unchanged':before==after,'start_unix':start,'elapsed_seconds':time.time()-start,'exit':code,'watchdog_seconds':a.timeout,'watchdog_expired':timedout,'log':logpath.name,'log_sha256':hashlib.sha256(logpath.read_bytes()).hexdigest()}
recordpath.write_text(json.dumps(record,indent=2)+'\n')
print(json.dumps({'name':a.name,'exit':code,'source_unchanged':before==after,'elapsed_seconds':record['elapsed_seconds'],'record':str(recordpath)}),flush=True)
sys.exit(code if before==after else 125)
