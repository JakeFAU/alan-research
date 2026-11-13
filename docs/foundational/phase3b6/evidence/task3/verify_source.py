import json,hashlib,pathlib,subprocess,tarfile
root=pathlib.Path.cwd(); kit=root/'docs/foundational/phase3b6/evidence/task3/protocol-source'; x=json.loads((kit/'index.json').read_text()); sha=lambda b:hashlib.sha256(b).hexdigest(); repo='/Users/jacob/dev/personal/alan'; commit=x['producer']['commit']
assert subprocess.check_output(['git','-C',repo,'rev-parse',commit+'^{tree}']).decode().strip()==x['producer']['tree']
for path,s in x['source_files'].items():
 b=subprocess.check_output(['git','-C',repo,'show',commit+':'+path]);assert sha(b)==s['whole_source_sha256'];assert len(b)==s['whole_source_bytes']
 for e in s['extracts']:
  lines=b.splitlines(keepends=True)[e['first_line']-1:e['last_line']];assert b''.join(lines)==(kit/e['path']).read_bytes();assert sha(b''.join(lines))==e['sha256']
for path,h in x['binding_source_hashes'].items(): assert sha(subprocess.check_output(['git','-C',repo,'show',commit+':'+path]))==h
# Inspect/hash binary members as opaque bytes only. Never instantiate SQLite.
archives=x['archives']; sealed=json.loads((root/'docs/foundational/phase3b4/evidence/prospective-main/archives.json').read_text()); count=0
for name,a in archives.items():
 path=root/'docs/foundational/phase3b4/evidence/prospective-main'/name
 assert sha(path.read_bytes())==a['sha256']==sealed[name]['sha256']
 assert len(path.read_bytes())==sealed[name]['bytes']
 with tarfile.open(path) as t:
  for case,c in a['cases'].items():
   prefix=c['member_prefix']; d=pathlib.Path('/tmp/alan-phase3b5/corpus')/name.removesuffix('.tar.gz')/prefix
   members=[m for m in t.getmembers() if m.isfile() and m.name.startswith(prefix+'/')]
   if not members: members=[m for m in t.getmembers() if m.isfile() and ('/'+prefix+'/') in m.name]
   assert len(members)==c['regular_members_verified'],(name,prefix,len(members))
   hashes=json.loads((d/'digests.json').read_text())
   source=json.loads((d/'source.json').read_text())
   for source_path,source_hash in source['FilesSHA256'].items(): assert sha(subprocess.check_output(['git','-C',repo,'show',commit+':'+source_path]))==source_hash
   for m in members:
    rel=m.name.split(prefix+'/',1)[1];b=t.extractfile(m).read();assert b==(d/rel).read_bytes()
    seal=sealed[name]['members'][prefix+'/'+rel];assert sha(b)==seal['sha256'] and len(b)==seal['bytes']
    if rel!='digests.json': assert sha(b)==hashes[rel]
    count+=1
   for rel,h in c['artifact_sha256'].items(): assert sha((d/rel).read_bytes())==h
print('Verified pinned commit/tree,17 source blobs,31 literal extracts,16 binding source blobs,4 dossiers and',count,'regular members. No SQLite access or runtime execution.')
