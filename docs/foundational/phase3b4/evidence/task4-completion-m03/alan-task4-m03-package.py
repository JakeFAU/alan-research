from pathlib import Path
import json,hashlib,tarfile,shutil,subprocess
runtime=Path('/private/tmp/alan-phase3b3-runtime');research=Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research')
out=research/'docs/foundational/phase3b4/evidence/task4-completion-m03';out.mkdir(parents=True,exist_ok=True)
sha=lambda p:hashlib.sha256(Path(p).read_bytes()).hexdigest()
def write(name,obj):(out/name).write_text(json.dumps(obj,indent=2,sort_keys=True)+'\n')
roots={'normal':Path('/tmp/alan-task4-m03-final-normal'),'race':Path('/tmp/alan-task4-m03-final-race'),'apparatus-red':Path('/tmp/alan-task4-m03-normal')};archives={};cases={}
for label,root in roots.items():
 with tarfile.open(out/(label+'.tar.gz'),'w:gz') as t:t.add(root,arcname=label)
 archives[label]={'archive':label+'.tar.gz','sha256':sha(out/(label+'.tar.gz')),'original_root':str(root),'files':{str(p.relative_to(root)):sha(p) for p in sorted(root.rglob('*')) if p.is_file()}}
 if label=='apparatus-red':continue
 f,=root.glob('*/case.json');c=json.loads(f.read_text());cases[label]={'directory':label+'/'+f.parent.name,'original_directory':str(f.parent),'case_sha256':sha(f),'case':c}
 s=json.loads((f.parent/'final/typed-inspection.json').read_text());q2='["b",1,6,0]';intent,=[x for x in s['Operations'] if x['ID']==q2];a,=[x for x in s['Admissions'] if x['SourceOperationID']==q2]
 assert intent['Generation']==a['RequestedGeneration']==2 and intent['Target']==a['Target']=='a'
 assert not a['Accepted'] and a['Reason']=='target_generation_changed' and a['Bound']=={'Subject':'a','Generation':1}
 assert len(s['Runs'])==2 and all(x['ID']['Generation']==1 for x in s['Runs'])
 for subject,ident in c['Manifest']['ProgramIDs'].items(): assert sha(f.parent/'programs'/f'{ident}.json')==ident
write('archives.json',archives);write('case-index.json',cases)
source={str(p.relative_to(runtime)):sha(p) for pat in ['**/*.go','go.mod','go.sum'] for p in runtime.glob(pat) if '.git' not in p.parts}
head=subprocess.check_output(['git','rev-parse','HEAD'],cwd=runtime,text=True).strip()
with tarfile.open(out/'runtime-source.tar.gz','w:gz') as t:
 for n in sorted(source):t.add(runtime/n,arcname=n)
write('runtime-source.json',{'head':head,'fix_base':'2b8f6e19a013d2b918146a7193451843004d2f86','files':source,'archive_sha256':sha(out/'runtime-source.tar.gz')})
(out/'commands').mkdir(exist_ok=True)
for name in ['m03-normal','m03-final-normal','m03-final-race','m03-staticcheck']:
 for suffix in ['.json','.log']:shutil.copyfile('/tmp/alan-task4-qualification/'+name+suffix,out/'commands'/(name+suffix))
 d=json.loads((out/'commands'/(name+'.json')).read_text())
 if name!='m03-normal':assert d['exit']==0 and d['source_sha256']==source
for n in ['alan-phase3b4-campaign-review-r1.md','alan-task4-review-check.py','alan-task4-review-check.json','alan-task4-review-extra.json','alan-task4-review-m03.json','alan-task4-m03-root-order-red.go','alan-task4-run.py','alan-task4-m03-package.py']:
 shutil.copyfile('/tmp/'+n,out/n)
old=research/'docs/foundational/phase3b4/evidence/task4';old_source=json.loads((old/'runtime-source.json').read_text())['files']
assert all(source[k]==v for k,v in old_source.items());assert set(source)-set(old_source)=={'internal/durabletest/phase3b4_generation_parent_test.go'}
preserved={}
for p in sorted(old.rglob('*')):
 if not p.is_file():continue
 rel=str(p.relative_to(research));blob=subprocess.check_output(['git','show','5d7473735d550ac1550228e9991863c3601afde7:'+rel],cwd=research)
 assert hashlib.sha256(blob).hexdigest()==sha(p),rel
 preserved[str(p.relative_to(old))]=sha(p)
write('approved-package-before-completion.json',{'research_base':'5d7473735d550ac1550228e9991863c3601afde7','files':preserved,'unchanged_runtime_base_files':len(old_source),'added_runtime_file':'internal/durabletest/phase3b4_generation_parent_test.go'})
write('manifest.json',{'version':'alan/phase3b4-task4-completion-m03/v1','runtime_fix_base':'2b8f6e19a013d2b918146a7193451843004d2f86','research_fix_base':'5d7473735d550ac1550228e9991863c3601afde7','status':'genuine M03 producer parent; scoped review and Task5 hostile mutation/oracle pending','files':{str(p.relative_to(out)):sha(p) for p in sorted(out.rglob('*')) if p.is_file() and p.name!='manifest.json'}})
print(head,[(k,v['directory']) for k,v in cases.items()])
