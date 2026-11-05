"""Run isolated, compiled checker mutants. Build failures/panics never kill.
Usage: python3 .../qualify_mutants.py SOURCE OUTPUT [FINAL_CORPUS]
Each OUTPUT must be new. Source copies include local modules and test fixtures.
"""
import hashlib,json,os,pathlib,shutil,subprocess,sys,time
source=pathlib.Path(sys.argv[1]).resolve();output=pathlib.Path(sys.argv[2]).resolve();output.mkdir(parents=True,exist_ok=False)
base=pathlib.Path('research/delegationcorrespondence')
mutants={
 'M01':('binding.go','!exists || o.Consumed == 0 || o.Consumed >= i.Sequence','!exists || o.Consumed == 0',['TestBindingDenseAttacks/M01'],'strict selection-before-service'),
 'M02':('binding.go','q := c.intents[d.OriginalRequestID]','''q := c.intents[d.OriginalRequestID]
 latest := int64(0)
 for _,candidate:=range c.observations {
  if candidate.Run==i.Run && candidate.Provenance=="ordinary-message" && candidate.Source==d.Root.Run.Subject && candidate.SourceGeneration==d.Root.Run.Generation && candidate.Consumed>latest && candidate.Consumed<i.Sequence {
   for _,admission:=range c.admissions {if identity("observation","message",admission.IntentID)==candidate.ID {q=c.intents[admission.IntentID];latest=candidate.Consumed}}
  }
 }''',['TestBindingDenseAttacks/M02','TestBindingExplicitPrincipals/P10'],'replace exact original/selection/admission with latest same-peer receive'),
 'M03':('custody.go','(a.Generation == 0 || v.Bound.Generation == a.Generation) &&','',['TestBindingDenseAttacks/M03'],'omit requested target generation versus bound generation'),
 'M04':('check.go','w := c.deriveWitness(a.Intent)\n\tallowed, reason := w.Allowed, w.Reason','savedRevocations:=c.revoked; c.revoked=map[AuthorityRoot]bool{}\n w := c.deriveWitness(a.Intent)\n c.revoked=savedRevocations\n allowed, reason := w.Allowed, w.Reason',['TestBindingDenseAttacks/M04'],'omit revocation only at fresh gate'),
 'A02-filter-pre-fix':('custody.go','f.Kind == o.Kind && f.Provenance == o.Provenance && source == o.Source && gen == o.SourceGeneration && obligation == o.Obligation && operation == o.Operation','(f.Kind=="" || f.Kind==o.Kind) && (f.Provenance=="" || f.Provenance==o.Provenance) && (source=="" || source==o.Source) && (gen==0 || gen==o.SourceGeneration) && (obligation=="" || obligation==o.Obligation) && (operation=="" || operation==o.Operation)',['TestBindingDenseAttacks/A02-filter'],'restore historical wildcard interpreter defect'),
}
env=dict(os.environ);env.update(GOTOOLCHAIN='local',GOWORK='off',GOPROXY='off',GOSUMDB='off',GOMODCACHE='/tmp/alan-phase3b4-qualification-modcache',GOCACHE='/tmp/alan-phase3a-go-build')
if len(sys.argv)>3:env['ALAN_PHASE3B4_CHECK_ROOT']=sys.argv[3]
for name in ('ALAN_PHASE3B4_ATTACK_OUTPUT','ALAN_PHASE3B4_CHECK_OUTPUT'):env.pop(name,None)
paths=[p for p in source.rglob('*') if p.is_file() and not any(x in p.relative_to(source).parts for x in ('.git','.superpowers','docs','paper')) and (p.suffix=='.go' or p.name in ('go.mod','go.sum') or 'testdata' in p.relative_to(source).parts)]
records=[]
for name,(filename,old,new,tests,predicate) in mutants.items():
 root=output/name;root.mkdir();tree=root/'source';tree.mkdir()
 for p in paths:
  target=tree/p.relative_to(source);target.parent.mkdir(parents=True,exist_ok=True);shutil.copyfile(p,target)
 target=tree/base/filename;text=target.read_text();assert text.count(old)==1,(name,text.count(old));target.write_text(text.replace(old,new))
 record={'name':name,'predicate':predicate,'path':str(base/filename),'before':old,'after':new,'source_sha256':{str(p.relative_to(tree)):hashlib.sha256(p.read_bytes()).hexdigest() for p in tree.rglob('*') if p.is_file()},'commands':[]}
 def run(argv,label):
  start=time.monotonic();log=root/(label+'.log')
  cwd=tree if label=='build' else tree/base
  with log.open('wb') as stream:r=subprocess.run(argv,cwd=cwd,env=env,stdout=stream,stderr=subprocess.STDOUT,timeout=300)
  record['commands'].append({'argv':argv,'cwd':str(cwd),'exit':r.returncode,'elapsed':time.monotonic()-start,'log':str(log),'log_sha256':hashlib.sha256(log.read_bytes()).hexdigest()})
  return r.returncode,log.read_text()
 code,_=run(['go','test','-c','-o',str(root/'checker.test'),'./research/delegationcorrespondence'],'build')
 record['built']=code==0;record['semantic_kills']=[]
 if code==0:
  for n,test in enumerate(tests):
   code,log=run([str(root/'checker.test'),'-test.run','^'+test.replace('/','$/^')+'$','-test.count=1'],f'assertion-{n}')
   semantic=code==1 and 'panic:' not in log and ('semantic assertion:' in log and 'valid=true' in log or name=='M02' and 'genuine distinct auxiliary rejected' in log and 'first_authority_policy' in log)
   record['semantic_kills'].append({'test':test,'killed':semantic})
 record['qualified_kill']=record['built'] and len(record['semantic_kills'])==len(tests) and all(x['killed'] for x in record['semantic_kills'])
 (root/'record.json').write_text(json.dumps(record,indent=2)+'\n');records.append(record);print(name,record['qualified_kill'],flush=True)
(output/'results.json').write_text(json.dumps(records,indent=2)+'\n')
sys.exit(0 if all(r['qualified_kill'] for r in records) else 1)
