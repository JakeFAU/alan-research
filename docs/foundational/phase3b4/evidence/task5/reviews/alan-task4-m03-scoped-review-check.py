from pathlib import Path
import base64,hashlib,io,json,posixpath,sqlite3,subprocess,tarfile,tempfile
R=Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research'); RT=Path('/private/tmp/alan-phase3b3-runtime'); P=R/'docs/foundational/phase3b4/evidence/task4-completion-m03'; OLD=R/'docs/foundational/phase3b4/evidence/task4'
RB='5d7473735d550ac1550228e9991863c3601afde7';RH='c122de0d5f51210b60d6b8af878a7ccd8a3f8c01';TB='2b8f6e19a013d2b918146a7193451843004d2f86';TH='fc60391f48b46a8c8635afc46882db8338b6c93b';NEW='internal/durabletest/phase3b4_generation_parent_test.go'
sha=lambda b:hashlib.sha256(b).hexdigest()
load=lambda p:json.loads(Path(p).read_text())
def git(w,*args):return subprocess.check_output(['git',*args],cwd=w)
def one(xs):
 assert len(xs)==1,len(xs)
 return xs[0]
def tarbytes(path):
 with tarfile.open(path) as t:return {m.name:t.extractfile(m).read() for m in t.getmembers() if m.isfile()}
out={}
assert git(RT,'diff','--name-only',TB,TH).decode().splitlines()==[NEW]
changes=git(R,'diff','--name-only',RB,RH).decode().splitlines();assert len(changes)==31
assert all(x.startswith('docs/foundational/phase3b4/evidence/task4-completion-m03/') or x in ['docs/foundational/phase3b4/evidence/task4/'+n for n in ['README.md','report.md','manifest.json']] for x in changes)
out['changed_files']={'runtime':1,'research':len(changes)}
# Verify both manifest trees, retained source tar, current bytes and exact committed source.
for root in [P,OLD]:
 m=load(root/'manifest.json');assert all(sha((root/n).read_bytes())==h for n,h in m['files'].items())
out['manifest_files']=len(load(P/'manifest.json')['files'])
src=load(P/'runtime-source.json'); assert src['head']==TH and src['fix_base']==TB
st=tarbytes(P/'runtime-source.tar.gz'); assert set(st)==set(src['files']);assert sha((P/'runtime-source.tar.gz').read_bytes())==src['archive_sha256']
with tarfile.open(fileobj=io.BytesIO(git(RT,'archive',TH))) as t:
 commit={m.name:t.extractfile(m).read() for m in t.getmembers() if m.isfile() and (m.name.endswith('.go') or m.name in ['go.mod','go.sum'])}
assert set(commit)==set(src['files'])
for n,h in src['files'].items():assert sha(st[n])==h==sha(commit[n])==sha((RT/n).read_bytes()),n
oldsrc=load(OLD/'runtime-source.json')['files'];assert all(src['files'][n]==h for n,h in oldsrc.items());assert set(src['files'])-set(oldsrc)=={NEW}
out['source_files']=len(src['files']);out['unchanged_original_source_files']=len(oldsrc)
pres=load(P/'approved-package-before-completion.json');oldpath=str(OLD.relative_to(R))
with tarfile.open(fileobj=io.BytesIO(git(R,'archive',RB,oldpath))) as t: oldbytes={m.name[len(oldpath)+1:]:t.extractfile(m).read() for m in t.getmembers() if m.isfile()}
assert set(pres['files'])==set(oldbytes)
for n,h in pres['files'].items():assert sha(oldbytes[n])==h
oldchanged=[n for n,b in oldbytes.items() if (OLD/n).read_bytes()!=b];assert sorted(oldchanged)==['README.md','manifest.json','report.md']
out['original_evidence_files']=len(oldbytes);out['original_evidence_unchanged']=len(oldbytes)-len(oldchanged);out['intentional_narrative_manifest_changes']=oldchanged
for n in ['alan-phase3b4-campaign-review-r1.md','alan-task4-review-check.py','alan-task4-review-check.json','alan-task4-review-extra.json','alan-task4-review-m03.json']:
 assert (P/n).read_bytes()==(Path('/tmp')/n).read_bytes(),n
out['original_review_artifacts_preserved']=5
ar=load(P/'archives.json');idx=load(P/'case-index.json');out['archive_members']={};out['cases']={}
for label,record in ar.items():
 data=tarbytes(P/record['archive']);assert sha((P/record['archive']).read_bytes())==record['sha256']
 files={n[len(label)+1:]:b for n,b in data.items()};assert set(files)==set(record['files']);assert all(sha(files[n])==h for n,h in record['files'].items());out['archive_members'][label]=len(files)
 if label=='apparatus-red':continue
 entry=idx[label];base=entry['directory']+'/'
 def get(n):return json.loads(data[base+n])
 assert sha(data[base+'case.json'])==entry['case_sha256'];case=get('case.json');assert case==entry['case']; assert case['AuxiliaryRequests']==['["b",1,6,0]']
 assert [x['Directory'] for x in case['Captures']]==['emitted-before-authorization','final'];assert case['Classification']=='source-authorized-target-generation-refusal-outside-C'
 manifest=case['Manifest'];principals=['["a",1,1,0]','["b",1,1,0]'];assert case['PrincipalRequests']==manifest['PrincipalRequests']==principals
 before=get('emitted-before-authorization/typed-inspection.json');s=get('final/typed-inspection.json');raw=get('final/inspect.json')['Tables'];br=get('emitted-before-authorization/inspect.json')['Tables'];q2='["b",1,6,0]';q='["b",1,1,0]';A={'Subject':'a','Generation':1};B={'Subject':'b','Generation':1}
 op=one([x for x in s['Operations'] if x['ID']==q2]);orig=one([x for x in s['Operations'] if x['ID']==q]);assert op['Run']==B and op['Sequence']==6 and op['Ordinal']==0 and op['Kind']=='message' and op['Target']=='a' and op['Generation']==2;assert op['Payload']==orig['Payload'] and op['Obligation']==orig['Obligation'];assert one([x for x in before['Operations'] if x['ID']==q2])==op
 rawop=one([x for x in raw['intents'] if x['id']==q2]);assert rawop['subject']=='b' and rawop['generation']==1 and rawop['sequence']==6 and rawop['ordinal']==0 and rawop['source_revision']==1;assert json.loads(rawop['descriptor'])==op
 hist=one([x for x in raw['history'] if x['subject']=='b' and x['generation']==1 and x['sequence']==6]);assert json.loads(hist['intents'])==[op] and hist['selected_observation'] is None and hist['source_revision']==1
 prev=one([x for x in raw['history'] if x['subject']=='b' and x['generation']==1 and x['sequence']==5]);assert json.loads(prev['successor'])['PC']==5 and json.loads(hist['successor'])['PC']==6
 hstep=one([x for x in get('final/history-000.json')['Steps'] if x['Run']==B and x['Sequence']==6]);assert hstep['Instruction']=='emit' and hstep['BeforePC']==5 and hstep['AfterPC']==6 and hstep['OperationIDs']==[q2]
 progs={x['id']:x for x in raw['programs']}
 for sub,pid in manifest['ProgramIDs'].items():
  b=data[base+'programs/'+pid+'.json'];assert sha(b)==pid and b==progs[pid]['bytes'].encode();assert next(x for x in get('provision.json')['Programs'] if x['ID']==pid)['Bytes']==base64.b64encode(b).decode()
 code=json.loads(progs[manifest['ProgramIDs']['b']]['bytes'])['Code'];inst=code[5];assert inst['Op']=='emit' and inst['Dst']=='q2';eff=inst['Effect'];assert eff['Generation']['Literal']['Number']==2 and eff['Target']['Literal']['Text']=='a';assert {k:v['Literal'] for k,v in eff['Payload'].items()}==op['Payload'];assert eff['Obligation']['Literal']['Text']==op['Obligation']
 # Check coherent prospective principal service/reply locators from actual programs, before Create.
 acode=json.loads(progs[manifest['ProgramIDs']['a']]['bytes'])['Code'];assert acode[4]['Effect']['Kind']=='external' and code[9]['Effect']['Kind']=='external';assert acode[7]['Effect']['Obligation']['Literal']['Text']=='ba-reply' and code[12]['Effect']['Obligation']['Literal']['Text']=='ab-reply';claim=json.loads(op['Payload']['claim']['Text']);assert claim['OriginalRequestID']==q and claim['ServiceOperationID']==manifest['Reciprocal']['BAService']=='["a",1,5,0]';assert manifest['Reciprocal']['ABService']=='["b",1,10,0]' and manifest['Reciprocal']['ABReply']=='["b",1,13,0]' and manifest['Reciprocal']['BAReply']=='["a",1,8,0]'
 delivery=get('auxiliary-delivery.json');auth=delivery['SourceAuthorization'];attempt=one([x for x in raw['attempts'] if x['intent_id']==q2]);assert json.loads(attempt['descriptor'])=={'Intent':op,'Revision':1};assert json.loads(attempt['result'])==auth;assert auth['Allowed'] and auth['Reason']=='authorized' and auth['IntentID']==q2 and auth['Source']=={'Run':B,'Epoch':1,'Sequence':6,'Revision':1}
 adm=one([x for x in s['Admissions'] if x['SourceOperationID']==q2]);ra=one([x for x in raw['admissions'] if json.loads(x['descriptor']).get('IntentID')==q2]);desc=json.loads(ra['descriptor']);res=json.loads(ra['result']);assert res==delivery['AdmissionResult'];assert desc['LogicalID']==desc['IntentID']==q2 and desc['Source']==B and desc['Target']=='a' and desc['Generation']==2 and desc['Payload']==op['Payload'] and desc['Obligation']==op['Obligation'];assert adm['Domain']=='message' and adm['LogicalID']==q2 and adm['Source']==B and adm['Target']=='a' and adm['RequestedGeneration']==2;assert not res['Accepted'] and res['Reason']=='target_generation_changed' and res['Bound']==A and not res['ObservationID'] and ra['observation_id'] is None;assert not any(x['SourceOperationID']==q2 for x in s['Observations'])
 assert not any(x['intent_id']==q2 for x in br['attempts']);assert not any(json.loads(x['descriptor']).get('IntentID')==q2 for x in br['admissions'])
 assert s['Composition']['Roots']==before['Composition']['Roots'] and len(s['Composition']['Roots'])==2 and s['Deployment']['Subjects']==get('provision.json')['Subjects'];assert len(s['Runs'])==2 and [x['ID'] for x in s['Runs']]==[A,B];assert all(x['Projection']['Status']=='open' and not x['Projection']['Data']['paid']['Flag'] for x in s['Runs']);assert not any(x['ID'] in [manifest['Reciprocal']['ABService'],manifest['Reciprocal']['BAService']] for x in s['Operations']);assert all(x['Contacts']==0 and not x['Applications'] for x in get('endpoints.json').values())
 changed=[k for k in raw if raw[k]!=br[k]];assert sorted(changed)==['admissions','attempts','composition_facts'];assert all(x in raw['composition_facts'] for x in br['composition_facts']);fact=one([json.loads(x['descriptor']) for x in raw['composition_facts'] if x not in br['composition_facts']]);assert fact['Kind']=='authorization' and fact['OperationID']==q2 and fact['Run']==B and fact['Sequence']==6;assert len(raw['admissions'])==len(br['admissions'])+1 and len(raw['attempts'])==len(br['attempts'])+1
 for cap in ['emitted-before-authorization','final']:
  assert set(get(cap+'/outcomes.json').values())=={'<nil>'};assert get(cap+'/replay.json')['Error']=='<nil>'
  # Only copied DB/WAL/SHM are opened, query-only. Compare all retained raw tables.
  with tempfile.TemporaryDirectory(prefix='alan-m03-review-') as td:
   db=Path(td)/'runtime.sqlite'
   for suffix in ['','-wal','-shm']:
    member=base+cap+'/runtime.sqlite'+suffix
    if member in data:(Path(str(db)+suffix)).write_bytes(data[member])
   conn=sqlite3.connect('file:'+str(db)+'?mode=ro&immutable=1',uri=True);conn.row_factory=sqlite3.Row
   conn.execute('PRAGMA query_only=ON')
   for name,rows in get(cap+'/inspect.json')['Tables'].items():
    actual=[{k:(v.decode() if isinstance(v,bytes) else v) for k,v in dict(r).items()} for r in conn.execute('SELECT * FROM "'+name+'"')]
    # Exports omit internal row-order; compare row dictionaries independent of order.
    assert sorted(map(lambda x:json.dumps(x,sort_keys=True),actual))==sorted(map(lambda x:json.dumps(x,sort_keys=True),rows)),(label,cap,name)
   conn.close()
 for n,h in get('binding-source.json').items():assert src['files'][posixpath.normpath('internal/durabletest/'+n)]==h
 for n,h in get('source.json')['FilesSHA256'].items():assert src['files'][n]==h
 out['cases'][label]={'namespace':s['Namespace'],'case_sha256':entry['case_sha256'],'A_program':manifest['ProgramIDs']['a'],'B_program':manifest['ProgramIDs']['b'],'q2':q2,'source_authorization':auth,'admission_descriptor':desc,'admission_result':res,'raw_tables_checked':len(raw),'raw_table_changes':changed,'history_steps':len(raw['history']),'runs':len(s['Runs']),'roots':len(s['Composition']['Roots'])}
for n in ['m03-final-normal','m03-final-race','m03-staticcheck']:
 d=load(P/'commands'/f'{n}.json');assert d['exit']==0 and d['source_sha256']==src['files'] and d['cwd']==str(RT);log=(P/'commands'/d['log']).read_text();assert (log.split()[:2]==['ok','github.com/JakeFAU/alan/internal/durabletest'] if n!='m03-staticcheck' else log=='')
red=load(P/'commands/m03-normal.json');assert red['exit']!=0;reds=red['source_sha256'];assert set(reds)==set(src['files']);assert all(h==src['files'][n] for n,h in reds.items() if n!=NEW);assert reds[NEW]==sha((P/'alan-task4-m03-root-order-red.go').read_bytes());assert 'principal/root set changed' in (P/'commands/m03-normal.log').read_text()
out['commands']={'passing_final_source_maps':3,'failed_apparatus_source_maps':1};out['sqlite_captures_verified']=4
Path('/tmp/alan-task4-m03-scoped-review-check.json').write_text(json.dumps(out,indent=2,sort_keys=True)+'\n')
print(json.dumps({k:v for k,v in out.items() if k!='cases'},indent=2));print('PASS: both causal joins and copied SQLite raw-table comparisons')
