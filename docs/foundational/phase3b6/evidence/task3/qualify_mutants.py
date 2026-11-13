#!/usr/bin/env python3
"""Compile disposable causal/renewal checker mutants; require semantic failures."""
import argparse, hashlib, json, os, pathlib, shutil, subprocess, tempfile, time

def main():
 p=argparse.ArgumentParser();p.add_argument('repository',type=pathlib.Path);p.add_argument('output',type=pathlib.Path);p.add_argument('corpus',type=pathlib.Path);a=p.parse_args();a.output.mkdir(parents=True,exist_ok=True)
 env=os.environ.copy();env.update(GOTOOLCHAIN='local',GOWORK='off',GOPROXY='off',GOSUMDB='off',GOMODCACHE='/tmp/alan-phase3b5/modcache',GOCACHE='/tmp/alan-phase3b6/go-build',ALAN_EFFECT_OBSERVATION_CORPUS=str(a.corpus))
 results=[]
 with tempfile.TemporaryDirectory(prefix='effect-mutants-',dir='/tmp/alan-phase3b6') as tmp:
  root=pathlib.Path(tmp)
  shutil.copy2(a.repository/'go.mod',root/'go.mod');shutil.copy2(a.repository/'go.sum',root/'go.sum')
  for pkg in ['delegationcorrespondence','authorityservice','mutualobligations','effectobservation']:
   shutil.copytree(a.repository/'research'/pkg,root/'research'/pkg)
  target=root/'research/delegationcorrespondence/effect_observation.go';original=target.read_text()
  renewal=root/'research/delegationcorrespondence/effect_protocol.go';renewal_original=renewal.read_text()
  trials=[('baseline',None,'^TestEffectObservation(Positive|CausalAckRequired|RenewalCoverage)$',0),
   ('omit-causal-acks','checkEffectCausalAcks(files, entry, legacy, out.Producer)','^TestEffectObservation(Positive|CausalAckRequired)$',1),
   ('omit-exact-renewal','helper','^TestEffectObservation(Positive|RenewalCoverage)$',1)]
  for name,needle,test,want in trials:
   source=original
   renewal.write_text(renewal_original)
   if needle=="helper":
    start=renewal_original.index("func checkEffectExactRenewal(");end=renewal_original.index("\nfunc checkEffectCallback",start)
    renewal.write_text(renewal_original[:start]+"func checkEffectExactRenewal(before, after []byte) {}\n"+renewal_original[end:])
   elif needle:
    assert source.count(needle)==1
    source=source.replace(needle,'// mutant: omitted '+needle)
   target.write_text(source)
   binary=root/(name+'.test');build=['go','test','-c','-o',str(binary),'./research/delegationcorrespondence']
   compiled=subprocess.run(build,cwd=root,env=env,capture_output=True);assert compiled.returncode==0,compiled.stderr.decode()
   command=[str(binary),'-test.v','-test.count=1','-test.run='+test];start=time.monotonic();run=subprocess.run(command,cwd=root,env=env,capture_output=True);elapsed=time.monotonic()-start
   log=run.stdout+run.stderr;(a.output/(name+'.log')).write_bytes(log)
   assert run.returncode==want,(name,run.returncode,log.decode())
   assert b'--- PASS: TestEffectObservationPositive' in log
   if needle:assert b'semantic assertion:' in log and b'panic:' not in log
   if name=='omit-exact-renewal':assert b'semantic assertion: exact renewal helper accepted' in log
   results.append(dict(name=name,compiled=True,build_command=build,command=command,expected_exit=want,exit=run.returncode,seconds=elapsed,source_sha256=hashlib.sha256(source.encode()).hexdigest(),renewal_source_sha256=hashlib.sha256(renewal.read_bytes()).hexdigest(),scope="helper-level exact renewal" if name=="omit-exact-renewal" else "end-to-end",log_sha256=hashlib.sha256(log).hexdigest()))
 (a.output/'results.json').write_text(json.dumps(results,indent=2)+'\n')
 print(json.dumps(results,indent=2))
if __name__=='__main__':main()
