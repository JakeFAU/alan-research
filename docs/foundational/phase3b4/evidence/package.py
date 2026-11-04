import pathlib,shutil,json,hashlib,tarfile,subprocess,gzip
src=pathlib.Path('/tmp/alan-phase3b4-task1'); dst=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3b4')
for f in ['01-charter.md','02-contract.md']:shutil.copyfile(src/f,dst/f)
base=dst/'evidence/baseline'; cons=dst/'evidence/construction';base.mkdir(parents=True,exist_ok=True);cons.mkdir(parents=True,exist_ok=True)
shutil.copyfile('/tmp/alan-phase3b4-evidence-recon.md',dst/'evidence/preflight/evidence-recon.md')
shutil.copyfile('/tmp/alan-phase3b4-preserved-manifest.json',base/'preserved-model-manifest.json')
for n in [0,1]: shutil.copyfile(f'/tmp/phase3b4Start{n}-metadata.json',base/f'pr-metadata-{n}.json')
shutil.copytree(src/'provenance',base/'provenance',dirs_exist_ok=True,ignore=shutil.ignore_patterns('*.diff'))
rawdiff=(src/'provenance/runtime-historical-production.diff').read_text()
(base/'provenance/runtime-historical-production.json').write_text(json.dumps({'command':'git diff 91963d5a81e4ce16f424aae36cc4fb382a56e2cd 1fdd4c9e1720bf60b344914bcbc3d0938ba8d678 -- internal/executor/composition_store.go internal/executor/composition_types.go','exit':0,'diff':rawdiff},indent=2)+'\n')
(base/'baseline-output.json').write_text(json.dumps({'exit':1,'output':(src/'baseline.log').read_text()},indent=2)+'\n')
for n in ['baseline-command.json','baseline-allowed.log','baseline-allowed-command.json','baseline-assertions.json','assert-baseline.py']:shutil.copyfile(src/n,base/n)
for n in ['construction-build.log','construction.log','construction-command.json','construction-final.log','construction-race.log','staticcheck.log','format.log','diff-check.log','verification.json']:shutil.copyfile(src/n,cons/n)
fixture=pathlib.Path('/private/tmp/alan-phase3b3-runtime/internal/durabletest/phase3b4_construction_test.go');shutil.copyfile(fixture,cons/'phase3b4_construction_test.go.txt')
for target,names in [(base,['baseline','baseline-allowed']),(cons,['construction','construction-final','construction-race'])]:
 manifest={}
 with tarfile.open(target/'evidence.tar.gz','w:gz') as tf:
  for name in names:
   for f in sorted((src/name).rglob('*')):
    if f.is_file():
     rel=f.relative_to(src).as_posix();tf.add(f,arcname=rel);raw=f.read_bytes();manifest[rel]={'sha256':hashlib.sha256(raw).hexdigest(),'bytes':len(raw)}
 raw=(target/'evidence.tar.gz').read_bytes();(target/'archive-manifest.json').write_text(json.dumps({'archive_sha256':hashlib.sha256(raw).hexdigest(),'archive_bytes':len(raw),'files':manifest},indent=2)+'\n')
 # Hash archive members as bytes, without SQLite or extraction.
 with tarfile.open(target/'evidence.tar.gz','r:gz') as tf:
  members=[m for m in tf.getmembers() if m.isfile()];assert len(members)==len(manifest)
  for m in members: assert hashlib.sha256(tf.extractfile(m).read()).hexdigest()==manifest[m.name]['sha256']
 (target/'archive-check.json').write_text(json.dumps({'members_checked':len(manifest),'all_sha256_match':True,'sqlite_opened':False},indent=2)+'\n')
# Preserve source identity of production and fixture separately, without reading environment files.
runtime=pathlib.Path('/private/tmp/alan-phase3b3-runtime'); files={}
for name in subprocess.check_output(['git','ls-files','*.go','go.mod','go.sum'],cwd=runtime,text=True).splitlines():
 f=runtime/name;files[name]=hashlib.sha256(f.read_bytes()).hexdigest()
(cons/'source.json').write_text(json.dumps({'runtime_base':'1fdd4c9e1720bf60b344914bcbc3d0938ba8d678','fixture_commit':subprocess.check_output(['git','rev-parse','HEAD'],cwd=runtime,text=True).strip(),'test_time_identity':'base plus separately hashed uncommitted construction fixture; final commit preserves tested bytes','files_sha256':files},indent=2)+'\n')
# Recheck parent-supplied preservation hashes now.
preserved=json.loads((base/'preserved-model-manifest.json').read_text());research=dst.parents[2]
for item in preserved['preserved']:
 raw=(research/item['path']).read_bytes();assert hashlib.sha256(raw).hexdigest()==item['sha256']
 if 'decoded_sha256' in item:assert hashlib.sha256(gzip.decompress(raw)).hexdigest()==item['decoded_sha256']
(base/'preservation-recheck.json').write_text(json.dumps({'count':len(preserved['preserved']),'all_sha256_match':True,'decoded_sha256_match':True},indent=2)+'\n')
print('packaged and verified baseline/construction archives')
