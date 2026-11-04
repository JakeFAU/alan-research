import pathlib,json,hashlib,tarfile,subprocess,glob,statistics,shutil
runtime=pathlib.Path('/private/tmp/alan-phase3b3-runtime')
research=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research')
out=research/'docs/foundational/phase3b4/evidence/task4';out.mkdir(parents=True,exist_ok=True)
sha=lambda p:hashlib.sha256(pathlib.Path(p).read_bytes()).hexdigest()
def write(name,v):(out/name).write_text(json.dumps(v,indent=2,sort_keys=True)+'\n')
roots={'normal':'/tmp/alan-task4-final-normal','race':'/tmp/alan-task4-final-race','supplemental-normal':'/tmp/alan-task4-evidence-fixed-normal','supplemental-race':'/tmp/alan-task4-evidence-fixed-race','promoted-normal':'/tmp/alan-task4-final-promoted-normal','promoted-race':'/tmp/alan-task4-final-promoted-race','prior-v1-normal':'/tmp/alan-task4-prior-regressions','prior-v1-race':'/tmp/alan-task4-prior-regressions-race','old-readers':'/tmp/alan-task4-final-old-readers'}
archives={}
for label,root in roots.items():
 root=pathlib.Path(root);assert root.is_dir(),root
 archive=out/(label+'.tar.gz')
 with tarfile.open(archive,'w:gz') as tar:tar.add(root,arcname=label)
 archives[label]={'archive':archive.name,'sha256':sha(archive),'original_root':str(root),'files':{str(p.relative_to(root)):sha(p) for p in sorted(root.rglob('*')) if p.is_file()}}
write('archives.json',archives)
selected={}
for mode in ['normal','race']:
 cases={}
 for label in [mode,'supplemental-'+mode]:
  root=pathlib.Path(roots[label])
  for f in sorted(root.glob('*/case.json')):
   c=json.loads(f.read_text());cases[c['Name']]={'case':c,'directory':label+'/'+f.parent.name,'original_directory':str(f.parent),'case_sha256':sha(f)}
 selected[mode]=cases
write('case-index.json',selected)
# Check capture lists mechanically, without interpreting authority/custody.
continuity=[]
for mode,cases in selected.items():
 for name,entry in cases.items():
  if not name.startswith('C0'):continue
  root=pathlib.Path(entry['original_directory']);c=entry['case'];last=''
  for capture in c['Captures']:
   assert capture['Predecessor']==last,(name,capture,last)
   assert capture['OwnerPID']>0
   path=root/capture['Directory'];assert (path/'runtime.sqlite').is_file()
   if capture['Label']=='post-open':
    before=json.loads((root/capture['Predecessor']/'inspect.json').read_text());after=json.loads((path/'inspect.json').read_text())
    assert before['Namespace']==after['Namespace'];changes={}
    assert before['Tables'].keys()==after['Tables'].keys()
    for table,left in before['Tables'].items():
     right=after['Tables'][table]
     if left==right:continue
     assert table=='runs',(name,table)
     assert len(left)==len(right)
     for a,b in zip(left,right):
      assert set(a)==set(b)
      keys=[k for k in a if a[k]!=b[k]];assert keys==['epoch'],(name,keys)
      assert b['epoch']==a['epoch']+1
     changes[table]={'only':'epoch','rows':len(left)}
    continuity.append({'mode':mode,'case':name,'before':entry['directory']+'/'+capture['Predecessor'],'after':entry['directory']+'/'+capture['Directory'],'changes':changes,'checked_all_tables':True})
   last=capture['Directory']
write('capture-continuity.json',{'scope':'mechanical raw-table epoch-only comparison; not an independent principal/authority oracle','pairs':continuity})
# Explicit role maps preserve physical namespaces and actual retained artifacts.
def role_record(entry):
 root=pathlib.Path(entry['original_directory']);c=entry['case'];m=c['Manifest']['Reciprocal'];s=json.loads((root/'final/typed-inspection.json').read_text());ops={o['ID']:o for o in s['Operations']}
 roles={k:m[k] for k in ['ABRequest','BARequest','Grant','Ready','BAService','BAReply','ABService','ABReply']}
 assert all(v in ops for v in roles.values())
 artifacts={r['ID']['Subject']:r['ProgramID'] for r in s['Runs']}
 assert artifacts==c['Manifest']['ProgramIDs']
 for ident in artifacts.values():assert sha(root/'programs'/f'{ident}.json')==ident
 return {'case':entry['directory'],'namespace':s['Namespace'],'roles':roles,'programs':artifacts,'roots':s['Composition']['Roots'],'permissions':s['Permissions'],'endpoint_a':json.loads((root/'endpoint-a-ready.json').read_text()),'endpoint_b':json.loads((root/'endpoint-b-ready.json').read_text()),'cost':json.loads((root/'final/cost-sidecar.json').read_text())}
correspondences=[]
for mode,cases in selected.items():
 ref=role_record(cases['P02-v2'])
 for name in ['P04','P05','P06','P07','P11']+[k for k in sorted(cases) if k.startswith('cost-')]:
  candidate=role_record(cases[name]);assert ref['namespace']!=candidate['namespace']
  correspondences.append({'mode':mode,'reference':ref,'candidate':candidate,'role_pairs':{k:[ref['roles'][k],candidate['roles'][k]] for k in ref['roles']},'scope':'declared designated roles across genuinely different provisions; no byte-equality, endpoint-identity, or whole-history equivalence claim'})
write('identity-correspondence.json',correspondences)
measurements=[]
for mode,cases in selected.items():
 for name,entry in cases.items():
  if name.startswith('cost-'):
   m=json.loads((pathlib.Path(entry['original_directory'])/'measurements.json').read_text());measurements.append({'mode':mode,'case':entry['directory'],**m})
write('cost-observations.json',measurements)
summary={}
for mode in ['normal','race']:
 summary[mode]={}
 for n in [0,1,8,32]:
  ms=[m for m in measurements if m['mode']==mode and m['N']==n];assert len(ms)==5
  summary[mode][n]={key:{'min_ns':min(x['Nanoseconds'][key] for x in ms),'median_ns':statistics.median(x['Nanoseconds'][key] for x in ms),'max_ns':max(x['Nanoseconds'][key] for x in ms)} for key in ms[0]['Nanoseconds']}
write('cost-summary.json',{'scope':'observed windows with explicit overhead; no allocation, threshold, speedup or total O(1) claim','measurements':summary})
# Full source bytes are retained independently of the mutable worktree.
sourcefiles={str(p.relative_to(runtime)):sha(p) for pat in ['**/*.go','go.mod','go.sum'] for p in runtime.glob(pat) if '.git' not in p.parts}
with tarfile.open(out/'runtime-source.tar.gz','w:gz') as tar:
 for name in sorted(sourcefiles):tar.add(runtime/name,arcname=name)
write('runtime-source.json',{'head':subprocess.check_output(['git','rev-parse','HEAD'],cwd=runtime,text=True).strip(),'base':'9e9c8078a0bd336a33bd1a03c85e3233249af625','files':sourcefiles,'archive_sha256':sha(out/'runtime-source.tar.gz')})
shutil.copytree('/tmp/alan-task4-qualification',out/'commands',dirs_exist_ok=True)
for p in ['/tmp/alan-phase3b4-runtime-review-r2.md','/tmp/alan-phase3b4-parent-c07-inspection.json','/tmp/alan-task4-lifecycle-writer-obstruction.txt','/tmp/alan-task4-run.py','/tmp/alan-task4-old-readers.py','/tmp/alan-task4-pre-evidence-fix-source.tar.gz','/tmp/alan-task4-final-old-reader-command.py','/tmp/alan-task4-package.py']:
 if pathlib.Path(p).is_file():shutil.copyfile(p,out/pathlib.Path(p).name)
# Keep superseded and failed apparatus roots in a separate archive, byte-identical.
excluded=set(pathlib.Path(v).resolve() for v in roots.values())|{pathlib.Path('/tmp/alan-task4-qualification')}
development=[p for p in pathlib.Path('/tmp').glob('alan-task4-*') if p.is_dir() and p.resolve() not in excluded]
with tarfile.open(out/'development.tar.gz','w:gz') as tar:
 for p in sorted(development):tar.add(p,arcname=p.name)
 for p in sorted(pathlib.Path('/tmp').glob('alan-task4-*.log')):tar.add(p,arcname=p.name)
write('development.json',{'scope':'superseded successful runs and apparatus failures; not final qualified corpus','roots':[str(p) for p in development],'sha256':sha(out/'development.tar.gz')})
write('manifest.json',{'version':'alan/phase3b4-task4-evidence/v1','runtime_base':'9e9c8078a0bd336a33bd1a03c85e3233249af625','research_base':'d8d2ec9026c6c343925c8d93c2be8f256976b0f9','status':'preliminary runtime qualification; principal oracle and frozen-pair final rerun are Task5','files':{str(p.relative_to(out)):sha(p) for p in sorted(out.rglob('*')) if p.is_file() and p.name!='manifest.json'}})
print('packaged',out)
