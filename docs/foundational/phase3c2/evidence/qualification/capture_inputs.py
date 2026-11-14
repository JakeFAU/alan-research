from pathlib import Path
import json,hashlib,subprocess,gzip
r=Path.cwd();q=r/'docs/foundational/phase3c2/evidence/qualification';base='69534188f796efab2e866750f6a76e2c7a3a1706';h=lambda b:hashlib.sha256(b).hexdigest()
files=json.loads((r/'docs/foundational/phase3c1/evidence/qualification/go-inputs.json').read_text())['files']
assert all(h((r/p).read_bytes())==s==h(subprocess.check_output(['git','show',base+':'+p])) for p,s in files.items())
tracked=subprocess.check_output(['git','ls-files'],text=True).splitlines();inputs=[p for p in tracked if p.endswith('.go') or Path(p).name in ('go.mod','go.sum','go.work','go.work.sum')];extra={};new={}
for p in sorted(set(inputs)-set(files)):
 if p.startswith('research/capacitycomposition/'):new[p]=h((r/p).read_bytes())
 else:
  s=h((r/p).read_bytes());assert s==h(subprocess.check_output(['git','show',base+':'+p]));extra[p]=s
src=json.loads(gzip.decompress((r/'docs/foundational/phase3c2/evidence/task1/export.json.gz').read_bytes()))['Sources'];assert all(h((r/e['Path']).read_bytes())==e['SHA256'] for e in src)
a=json.loads((r/'docs/foundational/phase3c2/evidence/task2/manifest.json').read_text());assert all(h((r/e['path']).read_bytes())==e['sha256'] for e in a['artifacts'])
summary={'base':base,'historical_manifest':'docs/foundational/phase3c1/evidence/qualification/go-inputs.json','historical_files':len(files),'unchanged_current':len(files),'unchanged_merged_base':len(files),'preexisting_tracked_inputs_outside_manifest_verified_against_base':extra,'new_inputs':new,'new_source_hashes_verified':len(src),'task2_artifacts_verified':len(a['artifacts']),'historical_full_suite_reused':True,'fresh_full_suite_run':False}
historical_base='1456a0000a0e14bb3dd105026189717c71ab905b'
comparison={}
for f,expected in extra.items():
 raw=subprocess.check_output(['git','show',historical_base+':'+f]);actual=h(raw)
 comparison[f]={'historical_baseline_present':True,'historical_baseline_sha256':actual,'matches_current':actual==expected}
assert all(e['matches_current'] for e in comparison.values())
summary['supplemental_historical_baseline_commit']=historical_base
summary['supplemental_historical_baseline_comparison']=comparison
(q/'input-identity.json').write_text(json.dumps(summary,indent=2)+'\n')

print(json.dumps({k:v for k,v in summary.items() if not isinstance(v,dict)},indent=2));print('additional unchanged inputs',len(extra),'new inputs',len(new))
