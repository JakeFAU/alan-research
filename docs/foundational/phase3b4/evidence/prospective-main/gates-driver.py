import json,os,pathlib,subprocess,sys
out=pathlib.Path('/tmp/alan-phase3b4-prospective-qualification')
plan=json.loads((out/'plan.json').read_text())
results=[]
for job in plan['jobs']:
 env={k:v for k,v in os.environ.items() if not k.startswith('ALAN_')}
 env.update(plan['common_environment'])
 if job['kind']=='runtime': env.update(plan['runtime_environment'])
 env.update(job.get('environment',{}))
 argv=['python3',str(out/'gate-recorder.py'),'--repo',job['repo'],'--out',str(out/'gates'),'--name',job['name'],'--expected-source',str(out/(job['kind']+'-frozen-inputs.json')),'--timeout','3900','--',*job['argv']]
 code=subprocess.run(argv,env=env).returncode
 results.append({'name':job['name'],'exit':code})
 (out/'driver-results.json').write_text(json.dumps(results,indent=2)+'\n')
 if code: raise SystemExit(code)
print('PROSPECTIVE GATES CLOSED',len(results),flush=True)
