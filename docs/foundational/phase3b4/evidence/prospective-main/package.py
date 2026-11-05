import hashlib,json,pathlib,shutil,tarfile
root=pathlib.Path('/tmp/alan-phase3b4-prospective-qualification')
out=root/'package'
if out.exists():raise SystemExit('Refusing to replace existing package')
out.mkdir()
sha=lambda b:hashlib.sha256(b).hexdigest()
archives={}
def pack(label,base):
 files=[p for p in sorted(base.rglob('*')) if p.is_file()]
 if not files:raise ValueError('Empty evidence root '+str(base))
 members={str(p.relative_to(base)):{'bytes':p.stat().st_size,'sha256':sha(p.read_bytes())} for p in files}
 dest=out/(label+'.tar.gz')
 with tarfile.open(dest,'w:gz') as tar:
  for p in files:tar.add(p,arcname=str(p.relative_to(base)),recursive=False)
 with tarfile.open(dest) as tar:
  actual={m.name:{'bytes':m.size,'sha256':sha(tar.extractfile(m).read())} for m in tar.getmembers() if m.isfile()}
 assert actual==members,label
 archives[dest.name]={'source_directory':str(base),'bytes':dest.stat().st_size,'sha256':sha(dest.read_bytes()),'members':members}
for mode in ('normal','race'):
 for label,stem in [('raw',''),('v1','v1-'),('promoted','promoted-'),('independent','verdicts-'),('attacks','attacks-')]:
  pack(label+'-'+mode,pathlib.Path('/tmp/alan-phase3b4-prospective-'+stem+mode))
pack('commands',root/'gates')
pack('sandbox-attempt',root/'sandbox-attempt')
for name in ['identity.json','plan.json','progress.md','summary.json','driver-results.json','sandbox-attempt-driver-results.json','gate-recorder.py','gates-driver.py','summarize.py','package.py','main-delta-review.diff','merge-tree.stdout','merge-tree.stderr','worktree.log','runtime-frozen-inputs.json','research-frozen-inputs.json','runtime-frozen-inputs.tar.gz','review-provisional.md','tools-and-readers.json','fresh-main-and-branch-check.json']:
 shutil.copy2(root/name,out/name)
(out/'archives.json').write_text(json.dumps(archives,indent=2)+'\n')
print(json.dumps({'archives':len(archives),'files':sum(len(a['members']) for a in archives.values()),'compressed_bytes':sum(a['bytes'] for a in archives.values())}))
