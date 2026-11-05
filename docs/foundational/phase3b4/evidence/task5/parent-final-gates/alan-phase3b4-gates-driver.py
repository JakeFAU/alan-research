#!/usr/bin/env python3
"""Execute a declared group only after parent writes frozen source/input maps."""
import json, os, pathlib, subprocess, sys
plan=json.loads(pathlib.Path('/tmp/alan-phase3b4-parent-gates-plan.json').read_text())
group=sys.argv[1];out=pathlib.Path('/tmp/alan-phase3b4-parent-final-gates');out.mkdir(exist_ok=True)
runtime='/private/tmp/alan-phase3b3-runtime';research='/Users/jacob/.codex/worktrees/phase3b4-research/alan-research'
repo=runtime if group.startswith('runtime') else research
kind='runtime' if repo==runtime else 'research'
expected=out/(kind+'-frozen-inputs.json')
if not expected.exists(): raise SystemExit('Frozen source/input map missing; nothing executed')
jobs=[]
if group in ('runtime-primary','runtime-ordinary','runtime-normal','runtime-race'):
 jobs=plan['runtime_suite_jobs']
 if group!='runtime-primary':jobs=[jobs[{'runtime-ordinary':0,'runtime-normal':1,'runtime-race':2}[group]]]
elif group in ('runtime-promoted','runtime-promoted-normal','runtime-promoted-race'):
 for mode in (('normal','race') if group=='runtime-promoted' else (group.removeprefix('runtime-promoted-'),)):
  jobs.append({'name':'runtime-promoted-'+mode,'argv':['go','test','-json','-count=1','-timeout=10m','-tags=integration']+(['-race'] if mode=='race' else [])+['./internal/executor','-run','^TestPhase3B4PromotedPublicParents$'],'environment':{'ALAN_PHASE3B4_PUBLIC_PARENTS':'/tmp/alan-task5-final-'+mode,'ALAN_PHASE3B4_PROMOTED_OUTPUT':'/tmp/alan-task5-final-promoted-'+mode}})
elif group in ('research-primary','research-ordinary','research-race','research-prior'):
 jobs=plan['research_suite_jobs']
 if group=='research-ordinary':jobs=jobs[:1]
 elif group=='research-race':jobs=jobs[1:2]
 elif group=='research-prior':jobs=jobs[2:]
elif group in ('runtime-small','research-small'):
 labels=['mod-verify','mod-tidy','build','vet','vet-integration','staticcheck','staticcheck-integration']
 jobs=[{'name':kind+'-'+label,'argv':argv,'expected_exit':1 if kind=='research' and label.startswith('staticcheck') else 0} for label,argv in zip(labels,plan['small_gate_argv'])]
 jobs.append({'name':kind+'-format','argv':['gofmt','-l']+(['adapters','author','cmd','compose','effects','examples','internal','runtime','schema'] if kind=='runtime' else ['research/delegationcorrespondence'])})
 if kind=='research': jobs.append({'name':'research-staticcheck-changed','argv':['staticcheck','-tags=integration','./research/delegationcorrespondence/...']})
else: raise SystemExit('Unknown group')
results=[]
for job in jobs:
 env={k:v for k,v in os.environ.items() if not k.startswith('ALAN_')}
 env.update(plan['common_environment'])
 if kind=='runtime':env.update(plan['runtime_environment'])
 env.update(job.get('environment',{}))
 argv=['python3','/tmp/alan-phase3b4-gate-recorder.py','--repo',repo,'--out',str(out),'--name',job['name'],'--expected-source',str(expected),'--timeout','3900','--',*job['argv']]
 code=subprocess.run(argv,env=env).returncode
 record={'name':job['name'],'recorder_exit':code,'expected_command_exit':job.get('expected_exit',0)}
 results.append(record)
 (out/(group+'-driver.json')).write_text(json.dumps({'group':group,'scope':'Driver expected-status accounting is not a claim that inherited non-green gates passed','jobs':results},indent=2)+'\n')
 if code!=job.get('expected_exit',0):
  print('UNEXPECTED STATUS; remaining group jobs not run',job['name'],code,flush=True)
  sys.exit(code or 125)
print('GROUP RECORDED',group,'including any explicitly expected non-green statuses',flush=True)
