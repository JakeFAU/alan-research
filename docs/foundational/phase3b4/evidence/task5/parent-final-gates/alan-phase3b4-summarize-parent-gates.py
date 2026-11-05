import hashlib,json,pathlib,sys
root=pathlib.Path('/tmp/alan-phase3b4-parent-final-gates')
records=[]
for p in sorted(root.glob('*.json')):
 r=json.loads(p.read_text())
 if not isinstance(r,dict) or not all(k in r for k in ('argv','source_unchanged','log','exit')): continue
 log=root/r['log'];data=log.read_bytes()
 assert hashlib.sha256(data).hexdigest()==r['log_sha256'],p
 name=p.stem;tests=[];packages=[];skips=[];no_tests=[];fail=[]
 for line in data.decode(errors='replace').splitlines():
  try:o=json.loads(line)
  except json.JSONDecodeError:continue
  if not isinstance(o,dict):continue
  if o.get('Action')=='pass':
   (tests if 'Test' in o else packages).append(o.get('Test',o.get('Package')))
  if o.get('Action')=='skip':
   (skips if 'Test' in o else no_tests).append({k:o[k] for k in ('Package','Test','Elapsed') if k in o})
  if o.get('Action')=='fail':fail.append({k:o[k] for k in ('Package','Test','Elapsed') if k in o})
 inherited=name in ('research-staticcheck','research-staticcheck-integration')
 expected=1 if inherited else 0
 records.append({'name':name,'argv':r['argv'],'exit':r['exit'],'expected_exit':expected,'disposition':'inherited_non_green_compare_exact_baseline' if inherited else ('passed' if r['exit']==0 else 'failed'),'source_unchanged':r['source_unchanged'],'elapsed_seconds':r['elapsed_seconds'],'test_pass_events':len(tests),'package_passes':len(packages),'named_skips':skips,'no_test_packages':no_tests,'failure_events':fail,'record_sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'log_sha256':r['log_sha256']})
required=['runtime-ordinary','runtime-integration','runtime-integration-race','runtime-promoted-normal','runtime-promoted-race','research-ordinary','research-model-checker-race','research-prior-reciprocal-normal','research-prior-reciprocal-race','research-prior-authority-normal','research-prior-authority-race']
for kind in ('runtime','research'):
 required += [kind+'-'+label for label in ('mod-verify','mod-tidy','build','vet','vet-integration','staticcheck','staticcheck-integration','format')]
required += ['research-staticcheck-changed']
missing=sorted(set(required)-{r['name'] for r in records})
comparison=json.loads((root/'lint-and-clean-output-comparison.json').read_text())
summary={'records':records,'required':required,'missing':missing,'lint_comparison':comparison,'complete':not missing,'expected_statuses_observed':all(r['exit']==r['expected_exit'] and r['source_unchanged'] and not r['failure_events'] for r in records),'scope':'Counts are test events, not independent experiments. Inherited lint exit1 is not a passing gate. Named skips are separate from packages with no tests.'}
(root/'gate-summary.json').write_text(json.dumps(summary,indent=2)+'\n')
print(json.dumps({'closed_commands':len(records),'missing':missing,'expected_statuses_observed':summary['expected_statuses_observed']}))
if missing or not summary['expected_statuses_observed']:sys.exit(1)
