"""Qualify immutable research source; preserve expected RED and old lint separately."""
from concurrent.futures import ThreadPoolExecutor
from datetime import datetime, timezone
from pathlib import Path
import gzip
import hashlib
import json
import os
import shutil
import subprocess
import sys
import time

repo = Path('/Users/jacob/dev/personal/alan-research')
out = Path('/tmp/alan-phase3b3-research-frozen')
expected = sys.argv[1]
head = subprocess.check_output(['git', 'rev-parse', 'HEAD'], cwd=repo, text=True).strip()
assert head == expected, (head, expected)
assert not out.exists(), 'Never overwrite qualification originals'
assert not subprocess.check_output(['git', 'diff', 'HEAD', '--', '*.go', '*go.mod', '*go.sum'], cwd=repo)
untracked = subprocess.check_output(['git', 'ls-files', '--others', '--exclude-standard'], cwd=repo, text=True).splitlines()
assert not [p for p in untracked if p.endswith('.go') or Path(p).name in ('go.mod', 'go.sum')]
out.mkdir()
env = os.environ.copy()
env.update(GOTOOLCHAIN='local', GOCACHE='/tmp/alan-phase3a-go-build', GOPROXY='off',
           GOSUMDB='off', STATICCHECK_CACHE='/tmp/alan-phase3b3-staticcheck')
inputs = subprocess.check_output(['git', 'ls-files', '--', '*.go', '*go.mod', '*go.sum'], cwd=repo, text=True).splitlines()
source = {p: hashlib.sha256((repo/p).read_bytes()).hexdigest() for p in inputs}
old_inputs = subprocess.check_output(['git', 'ls-tree', '-r', '--name-only', '8197683c23a3ca774750fc3c2c389c5a64dd38b8'], cwd=repo, text=True).splitlines()
old_inputs = [p for p in old_inputs if p.endswith('.go') or Path(p).name in ('go.mod', 'go.sum')]
for p in old_inputs:
    original = subprocess.check_output(['git', 'show', '8197683c23a3ca774750fc3c2c389c5a64dd38b8:'+p], cwd=repo)
    assert hashlib.sha256(original).hexdigest() == source[p], p
metadata = {'SourceCommit': head, 'BaseCommit': '8197683c23a3ca774750fc3c2c389c5a64dd38b8',
            'GoModuleFiles': source, 'UnchangedBaseInputs': len(old_inputs),
            'Go': subprocess.check_output(['go', 'version'], cwd=repo, env=env, text=True).strip(),
            'Staticcheck': subprocess.check_output(['staticcheck','-version'], cwd=repo, env=env, text=True).strip(),
            'RunnerSHA256': hashlib.sha256(Path(__file__).read_bytes()).hexdigest()}
fixture_paths = subprocess.check_output(['git','ls-files','--','research/delegationcorrespondence/testdata'],cwd=repo,text=True).splitlines()
fixtures = {p:hashlib.sha256((repo/p).read_bytes()).hexdigest() for p in fixture_paths}
metadata['FixtureFiles'] = fixtures
expectations = repo/'docs/foundational/phase3b3/evidence/task3-development/expectations.json'
assert expectations.exists()
(out/'expectations.json').write_bytes(expectations.read_bytes())
metadata['ExpectationInput'] = {'Path':str(expectations.relative_to(repo)),'SHA256':hashlib.sha256(expectations.read_bytes()).hexdigest(), 'Role':'explicit expected outcomes; not the raw policy oracle'}
(out/'source.json').write_text(json.dumps(metadata, indent=2)+'\n')
records = []

def run(name, command, additions=None, expected_exit=0, stdout_name=None, workdir=None):
    local = env.copy();local.update(additions or {})
    start = datetime.now(timezone.utc).isoformat();tick = time.monotonic()
    stdout = out/(stdout_name or name+'.log')
    stderr = out/(name+'.stderr') if stdout_name else None
    with stdout.open('wb') as stream:
        if stderr:
            with stderr.open('wb') as errors:
                done = subprocess.run(command,cwd=workdir or repo,env=local,stdout=stream,stderr=errors)
        else:
            done = subprocess.run(command,cwd=workdir or repo,env=local,stdout=stream,stderr=subprocess.STDOUT)
    r = {'Name':name,'Command':command,'Cwd':str(workdir or repo),'SourceCommit':head,
         'Environment':{k:local[k] for k in ['GOTOOLCHAIN','GOCACHE','GOPROXY','GOSUMDB','STATICCHECK_CACHE']},
         'AdditionalEnvironment':additions or {},'StartedUTC':start,'Seconds':time.monotonic()-tick,
         'Exit':done.returncode,'ExpectedExit':expected_exit,'Log':stdout.name,
         'SHA256':hashlib.sha256(stdout.read_bytes()).hexdigest()}
    if stderr:r['Stderr']={'Path':stderr.name,'SHA256':hashlib.sha256(stderr.read_bytes()).hexdigest()}
    if command[:2] == ['go','test']:
        events=[]
        for line in stdout.read_text(errors='replace').splitlines():
            try:events.append(json.loads(line))
            except json.JSONDecodeError:pass
        r['Tests']={'PassEvents':sum(e.get('Action')=='pass' and 'Test' in e for e in events),
                    'PackagePasses':[e['Package'] for e in events if e.get('Action')=='pass' and 'Test' not in e],
                    'Failures':[e for e in events if e.get('Action')=='fail'],
                    'Skips':[e for e in events if e.get('Action')=='skip']}
    (out/(name+'-command.json')).write_text(json.dumps(r,indent=2)+'\n')
    print(json.dumps({'Name':name,'Exit':done.returncode,'Expected':expected_exit,'Seconds':r['Seconds']}),flush=True)
    return r

# The explicit expected-outcome file does not remove the default semantic RED.
# This script is not execution evidence until its command records exist.
tasks = [
 ('format-new',['gofmt','-l','research/delegationcorrespondence'],{},0),
 ('ordinary-module',['go','test','-json','-count=1','-timeout=60m','./...'],{'ALAN_TASK3_RAW':'/tmp/alan-phase3b3-runtime-frozen/normal','ALAN_TASK3_PROGRESS':'1'},0),
 ('model-race',['go','test','-json','-race','-count=1','./research/twoagent/...','./research/mutualobligations/...','./research/authorityservice/...','./research/delegationcorrespondence/...'],{'ALAN_TASK3_RAW':'/tmp/alan-phase3b3-runtime-frozen/race','ALAN_TASK3_PROGRESS':'1'},0),
 ('prior-reciprocal-normal',['go','test','-json','-tags=integration','-count=1','./examples/communication/mutual'],{'ALAN_PHASE3B_RUNTIME_EXPORT':str(out/'prior-reciprocal-normal')},0),
 ('prior-reciprocal-race',['go','test','-json','-race','-tags=integration','-count=1','./examples/communication/mutual'],{'ALAN_PHASE3B_RUNTIME_EXPORT':str(out/'prior-reciprocal-race')},0),
 ('prior-authority-normal',['go','test','-json','-tags=integration','-count=1','./internal/durabletest','-run','^TestPhase3B2'],{'ALAN_PHASE2_EVIDENCE':str(out/'prior-authority-normal')},0),
 ('prior-authority-race',['go','test','-json','-race','-tags=integration','-count=1','./internal/durabletest','-run','^TestPhase3B2'],{'ALAN_PHASE2_EVIDENCE':str(out/'prior-authority-race')},0),
 ('staticcheck-new-model-fixtures',['staticcheck','-tags=integration','./research/twoagent/...','./research/mutualobligations/...','./research/authorityservice/...','./research/delegationcorrespondence/...','./examples/communication/authority','./examples/communication/mutual'],{},0),
 ('staticcheck-durable-normal',['staticcheck','./internal/durabletest'],{},0),
 ('staticcheck-durable-integration',['staticcheck','-tags=integration','./internal/durabletest'],{},1),
]

def export_all():
    local=[]
    exports=[('phase3a','./research/twoagent/cmd/export','docs/foundational/phase3a/evidence/model.json'),
             ('phase3b1','./research/mutualobligations/cmd/export','docs/foundational/phase3b/evidence/model.json'),
             ('phase3b2','./research/authorityservice/cmd/export','docs/foundational/phase3b2/evidence/model.json.gz'),
             ('phase3b2-audit','./research/authorityservice/cmd/audit','docs/foundational/phase3b2/evidence/audit.json.gz')]
    comparisons=[]
    for name,package,archive in exports:
        binary=str(out/(name+'-export'))
        local.append(run('build-'+name,['go','build','-buildvcs=false','-o',binary,package]))
        assert local[-1]['Exit']==0
        local.append(run(name+'-export',[binary],stdout_name=name+'.json'))
        old=repo/archive;expected_bytes=gzip.decompress(old.read_bytes()) if old.suffix=='.gz' else old.read_bytes()
        actual=(out/(name+'.json')).read_bytes();assert actual==expected_bytes,name
        comparisons.append({'Name':name,'PriorArchive':archive,'Equal':True,'SHA256':hashlib.sha256(actual).hexdigest()})
    (out/'prior-model-comparison.json').write_text(json.dumps(comparisons,indent=2)+'\n')
    controls=[('phase3a-loss','phase3a',['-assert-naive-recovery'],1),
              ('phase3b1-serial','phase3b1',['--profile','serial/serial','--assert-naive-completion'],1),
              ('phase3b1-positive','phase3b1',['--profile','servicing/serial','--assert-naive-completion'],0)]
    for config,assertion in [('S3/Missing','-assert-authority-safe'),('Cartesian/ScopeMix','-assert-authority-safe'),('S2/Missing','-assert-universal-completion'),('S2/LateGrant','-assert-universal-completion'),('S2/Revocable','-assert-universal-completion')]:
        controls.append(('phase3b2-'+config.replace('/','-'),'phase3b2',['-config',config,assertion],1))
    controls.append(('phase3b2-positive','phase3b2',['-config','S2/Explicit','-assert-authority-safe','-assert-universal-completion'],0))
    for name,model,flags,exitcode in controls:
        local.append(run(name,[str(out/(model+'-export'))]+flags,expected_exit=exitcode,stdout_name=name+'.json'))
    for name in ['corpus','check','progress','mutations']:
        binary=str(out/('correspondence-'+name))
        local.append(run('build-correspondence-'+name,['go','build','-buildvcs=false','-o',binary,'./research/delegationcorrespondence/cmd/'+name]))
        assert local[-1]['Exit']==0
    local.append(run('correspondence-corpus',[str(out/'correspondence-corpus')],stdout_name='model-corpus.json'))
    local.append(run('correspondence-progress',[str(out/'correspondence-progress')],stdout_name='progress.json'))
    local.append(run('correspondence-mutations',[str(out/'correspondence-mutations'),str(repo/'research/delegationcorrespondence/testdata/positive.json')],stdout_name='raw-mutations.json'))
    for profile in ['normal','race']:
        root='/tmp/alan-phase3b3-runtime-frozen/'+profile
        args=[str(out/'correspondence-check'),'-root',root]
        local.append(run('correspondence-'+profile+'-unqualified',args,expected_exit=1,stdout_name=profile+'-unqualified.json'))
        local.append(run('correspondence-'+profile+'-expected',args+['-expect',str(out/'expectations.json')],stdout_name=profile+'-expected.json'))
    mutants = [
        ('duplicate-keys','wire.go','if !ok || seen[s] {','if !ok {','TestStrictDuplicateAndExactIntegers'),
        ('attempt-binding','check.go','need(g.AttemptDigest == digest("alan/attempt-row/v1", a.Raw),','need(true || g.AttemptDigest == digest("alan/attempt-row/v1", a.Raw),','TestPrimaryGatePolicyBeyondHashes'),
        ('fresh-policy','check.go','need(g.Allowed == allowed && g.Reason == reason,','need(true || g.Allowed == allowed && g.Reason == reason,','TestFrozenFreshRevocationAndSelection'),
        ('causal-prefix','custody.go','c.requireObservationPrefix(o)','// omitted by mutation control','TestObservationSelectionRequiresHistoricalPredecessor|TestMessageAndTerminalSelectionRequireHistoricalSource')]
    mutant_records=[]
    for name,file,old,new,test in mutants:
        directory=out/'checker-mutants'/name
        shutil.copytree(repo/'research/delegationcorrespondence',directory)
        path=directory/file;before=path.read_bytes();text=before.decode();assert text.count(old)==1
        path.write_text(text.replace(old,new))
        (directory/'go.mod').write_text('module phase3b3-checker-mutant\n\ngo 1.27.3\n\nrequire github.com/JakeFAU/alan-research v0.0.0\nreplace github.com/JakeFAU/alan-research => '+str(repo)+'\n')
        result=run('checker-mutant-'+name,['go','test','-json','-count=1','-run','^('+test+')$','.'],{'ALAN_TASK3_RAW':'/tmp/alan-phase3b3-runtime-frozen/normal'},expected_exit=1,workdir=directory)
        local.append(result)
        assert result['Exit']==1 and any('Test' in event for event in result['Tests']['Failures']), name+' was not killed by a test assertion'
        mutant_records.append({'Name':name,'File':file,'OriginalSHA256':hashlib.sha256(before).hexdigest(),'MutatedSHA256':hashlib.sha256(path.read_bytes()).hexdigest(),'Old':old,'New':new,'CommandRecord':'checker-mutant-'+name+'-command.json'})
    (out/'checker-mutants.json').write_text(json.dumps(mutant_records,indent=2)+'\n')
    return local

with ThreadPoolExecutor(max_workers=3) as pool:
    futures=[pool.submit(run,*t) for t in tasks]
    exports=pool.submit(export_all)
    records=[f.result() for f in futures]+exports.result()
assert source=={p:hashlib.sha256((repo/p).read_bytes()).hexdigest() for p in inputs},'source changed'
assert not (out/'format-new.log').read_bytes(),'new source needs gofmt'
assert fixtures=={p:hashlib.sha256((repo/p).read_bytes()).hexdigest() for p in fixture_paths},'fixture changed'
oldlint=(repo/'docs/foundational/phase3b3/evidence/research-baseline/staticcheck-integration.log').read_bytes()
assert (out/'staticcheck-durable-integration.log').read_bytes()==oldlint,'old lint changed'
(out/'commands.json').write_text(json.dumps(records,indent=2)+'\n')
unexpected=[r['Name'] for r in records if r['Exit']!=r['ExpectedExit']]
(out/'completion.json').write_text(json.dumps({'SourceStable':True,'UnchangedBaseInputs':len(old_inputs),'OldLintByteIdentical':True,'UnexpectedStatuses':unexpected},indent=2)+'\n')
assert not unexpected,unexpected
print('research qualification complete',flush=True)
