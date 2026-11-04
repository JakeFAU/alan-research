from pathlib import Path
import hashlib,json,shutil,subprocess
out=Path('/tmp/alan-phase3b4-fix-r1'); runtime=Path('/private/tmp/alan-phase3b3-runtime'); research=Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research')
head=subprocess.check_output(['git','rev-parse','HEAD'],cwd=runtime,text=True).strip()
files={str(p.relative_to(runtime)):hashlib.sha256(p.read_bytes()).hexdigest() for p in sorted(runtime.rglob('*.go')) if '.git' not in p.parts}
for name in ['go.mod','go.sum']:files[name]=hashlib.sha256((runtime/name).read_bytes()).hexdigest()
final_gates=['p2-green','scoped-green-escalated','scoped-race-escalated','scoped-integration','scoped-integration-race','final-staticcheck','final-staticcheck-integration','final-format','final-diff-check']
for gate in final_gates:
 rec=json.loads((out/(gate+'.json')).read_text()); assert rec['source_sha256']==files,gate; assert rec['exit_code']==0,gate
contract=hashlib.sha256((research/'docs/foundational/phase3b4/02-contract.md').read_bytes()).hexdigest();assert contract=='3b3662639f22adf9bfe3c2e5297afc7b5b6f4a389c6e8be4cb6a2f7deca9f3f3'
(out/'final-source.json').write_text(json.dumps({'runtime_base':'5b68180bc18a5c7ee3f5be40d77b251e054631a9','runtime_head':head,'commits':['c03cd09c2cba1d0bee46a6c7bf125817cb001680',head],'research_base':'a9125eaee96a73fc464b7879804aad1cc326fbd7','contract_sha256_unchanged':contract,'source_sha256':files,'successful_gates_verified_exact_final_source':final_gates},indent=2)+'\n')
commands=[]
for p in sorted(out.glob('*.json')):
 rec=json.loads(p.read_text())
 if 'command' not in rec or 'source_sha256' not in rec:continue
 commands.append({'record':p.name,'output':p.stem+'.log','exit_code':rec['exit_code'],'command':rec['command'],'execution':'require_escalated' if p.stem in ['scoped-green-escalated','scoped-race-escalated','scoped-integration','scoped-integration-race'] else 'sandbox','classification':'apparatus loopback bind failure' if p.stem in ['scoped-green','scoped-race'] else 'causal RED' if 'red' in p.stem else 'successful check'})
(out/'commands.json').write_text(json.dumps(commands,indent=2)+'\n')
for p in (runtime/'internal/executor/testdata').glob('v1-dormant*'):shutil.copy2(p,out/p.name)
dest=research/'docs/foundational/phase3b4/evidence/task3-fix-r1';dest.mkdir(exist_ok=True)
for p in out.iterdir():
 if p.suffix in ['.json','.log','.txt','.py','.gz','.md']:shutil.copy2(p,dest/p.name)
print('Retained',len(list(dest.iterdir())),'evidence files; final runtime',head,'source files',len(files))
