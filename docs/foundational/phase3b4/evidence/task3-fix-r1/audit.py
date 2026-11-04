from pathlib import Path
import hashlib,json,re,tarfile,subprocess
r=Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research');root=r/'docs/foundational/phase3b4/evidence/task3-fix-r1'; runtime=Path('/private/tmp/alan-phase3b3-runtime')
# Verify reconstructed comparator against Git source, not against its own manifest.
base=json.loads((root/'base-source.json').read_text())
for rec in base['files']:
 data=subprocess.check_output(['git','show',base['commit']+':'+rec['path']],cwd=runtime)
 assert hashlib.sha256(data).hexdigest()==rec['sha256']
for cmd in json.loads((root/'commands.json').read_text()):
 rec=json.loads((root/cmd['record']).read_text())
 assert hashlib.sha256((root/cmd['output']).read_bytes()).hexdigest()==rec['output_sha256']
 assert cmd['exit_code']==rec['exit_code']
 assert rec['command']==cmd['command']
 assert len(rec['source_sha256']) in (307,316)
 for p,h in rec['source_sha256'].items():
  assert not p.startswith('/') and (p.endswith('.go') or p in ('go.mod','go.sum'))
  assert re.fullmatch('[a-f0-9]{64}',h)
patterns=rb'-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----|ghp_[A-Za-z0-9]{30}|github_pat_[A-Za-z0-9_]{30}|AKIA[A-Z0-9]{16}|sk-[A-Za-z0-9]{30}'
scanned=[]
for p in sorted(root.iterdir()):
 if p.name in ('manifest.json','retention-checks.json'):continue
 if p.name.endswith('.tar.gz'):
  with tarfile.open(p) as t:
   for m in t.getmembers():
    assert m.isfile() and Path(m.name).name==m.name
    assert not re.search(patterns,t.extractfile(m).read()),(p.name,m.name)
 else:
  raw=p.read_bytes();assert not re.search(patterns,raw),p.name
  if p.suffix in ['.json']:json.loads(raw)
 scanned.append(p.name)
for label,link in re.findall(r'\[([^]]+)\]\(([^)]+)\)',(root/'report.md').read_text()):assert (root/link).exists(),link
review=Path('/tmp/alan-phase3b4-runtime-review-r1.md').read_bytes();assert (root/'review.md').read_bytes()==review
checks={'review_sha256':hashlib.sha256(review).hexdigest(),'base_git_source_files_verified':len(base['files']),'complete_evidence_files_scanned':scanned,'secret_pattern_candidates':[],'scope':'Only round1 review/probes, base-produced fixture, commands/source hashes/output and correction report; archive members checked as bytes, never opened as databases. Every source hash map has only Go/module paths and SHA256 values.','report_links':'all resolve','runtime_status':subprocess.check_output(['git','status','--short'],cwd=runtime,text=True)}
(root/'retention-checks.json').write_text(json.dumps(checks,indent=2)+'\n')
manifest={'runtime_base':'5b68180bc18a5c7ee3f5be40d77b251e054631a9','runtime_head':'9e9c8078a0bd336a33bd1a03c85e3233249af625','research_base':'a9125eaee96a73fc464b7879804aad1cc326fbd7','files':{p.name:hashlib.sha256(p.read_bytes()).hexdigest() for p in sorted(root.iterdir()) if p.name!='manifest.json'}}
(root/'manifest.json').write_text(json.dumps(manifest,indent=2)+'\n')
print('Verified',len(base['files']),'base source files;',len(manifest['files']),'retained evidence files; no credential-pattern candidates; runtime clean:',checks['runtime_status']=='')
