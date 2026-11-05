import hashlib,json,pathlib,tarfile
repo=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research');out=repo/'docs/foundational/phase3b4/evidence/task5';work=pathlib.Path('/tmp/alan-task5')
def sha(x):return hashlib.sha256(x).hexdigest()
roots={}
for mode in ['normal','race']:
 for kind in ['','promoted-','v1-']:roots['final-'+kind+mode]=pathlib.Path('/tmp/alan-task5-final-'+kind+mode)
 roots['independent-'+mode]=work/('final-'+mode+'-verdicts');roots['attacks-'+mode]=work/('final-'+mode+'-attacks')
 for kind in ['authority','reciprocal']:roots['prior-'+kind+'-'+mode]=pathlib.Path('/tmp/alan-task5-prior-'+kind+'-'+mode)
records={}
for name,root in roots.items():
 expected={str(p.relative_to(root)):sha(p.read_bytes()) for p in root.rglob('*') if p.is_file()};actual={}
 with tarfile.open(out/(name+'.tar.gz')) as t:
  for m in t:
   if m.isfile():
    assert m.name not in actual;actual[m.name]=sha(t.extractfile(m).read())
 assert actual==expected,(name,set(actual)^set(expected));records[name]={'Files':len(actual),'AllBytesEqualSource':True}
for name in ['checker-mutants','preliminary-apparatus']:
 with tarfile.open(out/(name+'.tar.gz')) as t:
  mapping=json.load(t.extractfile('files.json'));blobs={m.name.removeprefix('sha256/'):sha(t.extractfile(m).read()) for m in t if m.name.startswith('sha256/')}
  assert all(k==v for k,v in blobs.items());assert set(mapping.values())==set(blobs)
  records[name]={'Files':len(mapping),'Blobs':len(blobs),'AllMappedBlobsHashExact':True}
(out/'archive-verification.json').write_text(json.dumps(records,indent=2)+'\n');print(json.dumps(records,indent=2))
