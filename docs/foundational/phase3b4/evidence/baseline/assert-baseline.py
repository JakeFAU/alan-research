# Targeted assertion on this baseline pair, not the Phase 3B.4 policy oracle.
import json, pathlib
root=pathlib.Path('/tmp/alan-phase3b4-task1/baseline-allowed')
rows=[]
for name,selected,allowed,reason in [('ordinary-s2-handshake',3,True,'authorized'),('control-wrong-prior-selection',2,False,'authority_request_not_admitted')]:
 d=next(root.glob(name+'-*')); s=json.load(open(d/'final/typed-inspection.json'));m=json.load(open(d/'manifest.json')); ep=json.load(open(d/'endpoints.json'))
 e=next(x for x in s['Effects'] if x['OperationID']==m['BAService']);w=e['ExternalAttempt']['Authority'];g=next(x for x in s['LifecycleDecisions'] if x['Request']['Intent']['ID']==m['Grant']);claim=json.loads(g['Request']['Intent']['Payload']['claim']['Text']);o=next(x for x in s['Observations'] if x['SourceOperationID']==m['BARequest'] and x['Run']==m['A']);ready=next(x for x in s['Observations'] if x['SourceOperationID']==m['Ready'] and x['Run']==m['A'])
 assert m['BAService']=='["a",1,4,0]' and e['Sequence']==4
 assert g['Accepted'] and len(s['LifecycleDecisions'])==1 and claim==m['Claim'] and claim['ServiceOperationID']==e['OperationID'] and claim['OriginalRequestID']==o['SourceOperationID']
 assert e['Authorization']['Allowed'] and e['ExternalAttempt']['Allowed']==allowed and w['Allowed']==allowed and w['Reason']==reason and w['Use']['Root']==claim['Root']
 assert o['SelectedAt']['Sequence']==selected and ready['SelectedAt']['Sequence']==5-selected and o['Admissions']
 h=json.load(open(d/'final/history-000.json'));actual=next(x for x in h['Steps'] if x['Run']==m['A'] and x['Sequence']==selected);assert actual['SelectedObservationID']==o['ID']
 unpaid=[not x['Projection']['Data']['paid']['Flag'] for x in s['Runs']];assert unpaid==[not allowed,not allowed]
 contacts={k:v['Contacts'] for k,v in ep.items()};assert contacts=={'a':int(allowed),'b':int(allowed)}
 if allowed: assert all(x['Projection']['Status']=='succeeded' for x in s['Runs']) and e['ExternalReceipt']
 else: assert not e.get('AuthorityAdmissions') and e['ExternalReceipt'] is None and all(x['Charges']['Effects']==0 for x in s['Roots'])
 rows.append({'case':name,'directory':d.name,'namespace':s['Namespace'],'service':e['OperationID'],'service_sequence':e['Sequence'],'request':o['SourceOperationID'],'observation':o['ID'],'selection':o['SelectedAt'],'ready_selection':ready['SelectedAt'],'grant':m['Grant'],'claim':claim,'authority':w,'unpaid':unpaid,'contacts':contacts,'pending':s['Pending'],'run':json.load(open(d/'run.json'))})
print(json.dumps({'targeted_assertions':'passed','cases':rows},indent=2))
