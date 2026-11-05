import json,pathlib
out=pathlib.Path('/tmp/alan-task5/gate-summary.json');records=[]
for label,root in [('independent',pathlib.Path('/tmp/alan-task5/final-commands')),('controller',pathlib.Path('/tmp/alan-phase3b4-parent-final-gates'))]:
 for p in sorted(root.glob('*.json')):
  d=json.loads(p.read_text())
  if not isinstance(d,dict) or 'argv' not in d or 'exit' not in d:continue
  rec={k:d[k] for k in ['cwd','argv','environment','exit','source_unchanged','elapsed_seconds','watchdog_expired'] if k in d};rec.update(owner=label,name=p.stem,record=str(p),test_skips=[],package_no_tests=[])
  log=p.with_suffix('.log');latest={}
  for line in log.read_text().splitlines() if log.exists() else []:
   try:j=json.loads(line)
   except (ValueError,TypeError):continue
   if not isinstance(j,dict):continue
   key=(j.get('Package'),j.get('Test'))
   if j.get('Action')=='output':latest[key]=j.get('Output')
   if j.get('Action')=='skip':
    target='test_skips' if j.get('Test') else 'package_no_tests';rec[target].append({'package':key[0],'test':key[1],'last_output':latest.get(key)})
  records.append(rec)
out.write_text(json.dumps(records,indent=2)+'\n');print([(r['name'],r['exit'],len(r['test_skips'])) for r in records])
