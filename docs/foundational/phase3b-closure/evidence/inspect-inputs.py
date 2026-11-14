"""Read-only document audit: JSON/tar text inspection, never SQLite or model execution."""
import gzip, hashlib, json, pathlib, subprocess, tarfile
ROOT = pathlib.Path(__file__).resolve().parents[4]
OUT = pathlib.Path(__file__).resolve().parent
inputs = {}
members = {}
def read(p):
    b = (subprocess.check_output(['git','show','7ae2c5ddd5e18d2907dd91fcefc896456a25562d:'+p], cwd=ROOT)
         if p == 'docs/foundational/phase3/README.md' else (ROOT / p).read_bytes())
    inputs[p] = {'bytes': len(b), 'sha256': hashlib.sha256(b).hexdigest()}
    return b
def data(p):
    b = read(p)
    return json.loads(gzip.decompress(b) if p.endswith('.gz') else b)
def member(t, archive, name):
    assert name.endswith(('.json', '.jsonl', '.go', '.txt'))
    b = t.extractfile(name).read()
    members[archive + '::' + name] = {'bytes': len(b), 'sha256': hashlib.sha256(b).hexdigest()}
    return b
facts = {}
p = 'docs/foundational/phase3b/evidence/runtime-normal/serial_serial-send_a-deliver_a/summary.json'
facts['serial_wait'] = data(p)
p = 'docs/foundational/phase3b2/evidence/model-summary.json'
d = data(p)
facts['authority'] = {'states': d['States'], 'arcs': d['Arcs'], 'selected': [x for x in d['Configurations'] if (x['Config']['Policy'],x['Config']['Scenario']) in [('S2','Shared'),('S2','Revocable'),('S3','Missing'),('Cartesian','ScopeMix')]]}
p = 'docs/foundational/phase3b5/evidence/task2/revocation-progress.json.gz'
d = data(p)[0]
facts['revocation_progress'] = {k:v for k,v in d.items() if not isinstance(v,(dict,list))}
facts['revocation_progress_keys'] = list(d)
facts['revocation_audit'] = data('docs/foundational/phase3b5/evidence/task2/revocation-audit.json.gz')['Graphs']
d = data('docs/foundational/phase3b6/evidence/task2/audit.json.gz')['report']
facts['split_audit'] = {k:v for k,v in d.items() if not isinstance(v,(dict,list))}
facts['cancel_after_admission'] = {}
for mode in ['normal','race']:
    for route in ['idempotent','adapter']:
        p = f'docs/foundational/phase2b/evidence/process-{mode}/TestPhase2BLiveOrders_{route}_A-source-live.json'
        d = data(p)
        facts['cancel_after_admission'][mode+'/'+route] = {
            'events': [({k:e[k] for k in ['Boundary','Phase','Subject'] if k in e} if isinstance(e,dict) else e) for e in d['events']],
            'world_held_contacts': d['world_held']['Contacts'], 'world_final_contacts': d['world_final']['Contacts'],
            'applications': len(d['world_final']['Applications']),
            'revoked_lifecycle': [json.loads(r['result']) for r in d['revoked']['Tables']['lifecycle_results']]}
facts['retry_obstruction'] = {}
p = 'docs/foundational/phase3b4/evidence/task5/final-normal.tar.gz'
read(p)
with tarfile.open(ROOT/p) as t:
    for m in t.getmembers():
        if 'N10e' in m.name and m.name.endswith('construction-obstruction.json'):
            d = json.loads(member(t,p,m.name))['ActualRetry']
            facts['retry_obstruction'][m.name] = {k:d[k] for k in ['Accepted','Reason','Successor']}
            facts['retry_obstruction'][m.name]['runs'] = [d[k]['Run'] for k in ['Before','After']]
facts['matched'] = {}
for mode in ['normal','race']:
    p = f'docs/foundational/phase3b7/evidence/capture/{mode}-v1/raw.tar.gz'
    read(p)
    with tarfile.open(ROOT/p) as t:
        eq = {}
        for suffix in ['first-owner/post-open/inspect.json','stopped/inspect.json','reopened/inspect.json','final/inspect.json','run.json']:
            a,b = [member(t,p,branch+'/'+suffix) for branch in ['zero','applied']]
            eq[suffix] = a == b
        report = json.loads(member(t,p,'zero/run.json'))
        facts['matched'][mode] = {'byte_equal':eq, 'report':report,'regular_members':sum(m.isfile() for m in t.getmembers())}
facts['corrected'] = data('docs/foundational/phase3b7/evidence/qualification/corrected-checks.json')
facts['initial'] = data('docs/foundational/phase3b7/evidence/qualification/initial-checks.json')
p = 'docs/foundational/phase3b4/evidence/task4/runtime-source.tar.gz'
read(p)
with tarfile.open(ROOT/p) as t:
    excerpts = {}
    for name,lo,hi in [('internal/executor/dispatch.go',90,148),('internal/executor/external.go',185,210),('internal/executor/lifecycle.go',310,348),('internal/executor/request_binding.go',1,180)]:
        b = member(t,p,name).decode().splitlines()
        excerpts[name] = '\n'.join(f'{i+1}: {b[i]}' for i in range(lo-1,min(hi,len(b))))
    facts['archived_source_excerpts'] = excerpts
for p in [
 'docs/foundational/phase3/README.md',
 'docs/foundational/phase3b/01-contract.md', 'docs/foundational/phase3b/03-results.md',
 'docs/foundational/phase3b2/01-contract.md','docs/foundational/phase3b2/03-results.md',
 'docs/foundational/phase3b3/02-contract.md','docs/foundational/phase3b3/04-results.md',
 'docs/foundational/phase3b4/06-qualification.md','docs/foundational/phase3b4/08-results.md',
 'docs/foundational/phase3b4/evidence/task4/alan-task4-lifecycle-writer-obstruction.txt',
 'docs/foundational/phase3b5/01-contract.md','docs/foundational/phase3b5/04-correspondence.md','docs/foundational/phase3b5/07-results.md','docs/foundational/phase3b5/evidence/task3/verdicts.json',
 'docs/foundational/phase3b6/01-contract.md','docs/foundational/phase3b6/04-correspondence.md','docs/foundational/phase3b6/07-results.md','docs/foundational/phase3b6/evidence/task3/verdicts.json',
 'docs/foundational/phase3b7/01-contract.md','docs/foundational/phase3b7/03-comparison.md','docs/foundational/phase3b7/05-results.md','docs/foundational/phase3b7/evidence/task2/review-r2/root-cause.json',
 'docs/foundational/phase2b/01-contract.md', 'internal/durabletest/phase2b_live_test.go',
 'research/mutualobligations/model.go','research/authorityservice/model.go','research/delegationcorrespondence/binding.go','research/delegationcorrespondence/binding_attacks_test.go',
 'research/effectobservation/model.go','research/matchedcapture/check.go','research/matchedcapture/compare.go','research/matchedcapture/protocol.go']:
    read(p)
(OUT/'inspected-inputs.json').write_text(json.dumps({'base':'7ae2c5ddd5e18d2907dd91fcefc896456a25562d','roadmap_from_base':True,'files':inputs,'archive_text_members':members},indent=2)+'\n')
(OUT/'inspected-facts.json').write_text(json.dumps(facts,indent=2)+'\n')
print(json.dumps({'input_files':len(inputs),'text_members':len(members),'matched':{m:facts['matched'][m]['byte_equal'] for m in facts['matched']}}))
