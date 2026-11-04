"""Targeted independent joins for one supplied JSON snapshot; not a general validator."""
import hashlib
import json
from pathlib import Path

case = Path('/tmp/alan-phase3b3-runtime-development/controls/control-wrong-prior-selection-1914204028')
path = case / 'final/inspect.json'
raw = path.read_bytes()
d = json.loads(raw)
t = d['Tables']
dec = json.loads
sha = lambda b: hashlib.sha256(b).hexdigest()
canon = lambda x: json.dumps(x, separators=(',', ':'), ensure_ascii=False)
ids = {'AB': '["a",1,1,0]', 'BA': '["b",1,1,0]', 'Grant': '["b",1,2,0]', 'Ready': '["b",1,5,0]', 'Service': '["a",1,4,0]'}
intents = {r['id']: dec(r['descriptor']) for r in t['intents']}
history = {(r['subject'], r['generation'], r['sequence']): r for r in t['history']}
obs = {r['id']: r for r in t['observations']}
facts = [dec(r['descriptor']) for r in t['composition_facts']]
assert [f['Order'] for f in facts] == list(range(1, len(facts)+1))
for r, f in zip(t['composition_facts'], facts):
    assert (r['id'], r['kind'], r['decision_order']) == (f['ID'], f['Kind'], f['Order'])
for r in t['intents']:
    i = intents[r['id']]
    assert i['ID'] == canon([r['subject'], r['generation'], r['sequence'], r['ordinal']])
    assert i['Run'] == {'Subject': r['subject'], 'Generation': r['generation']}
    assert (i['Sequence'], i['Ordinal']) == (r['sequence'], r['ordinal'])
    assert i in dec(history[(r['subject'], r['generation'], r['sequence'])]['intents'])
lr, = t['lifecycle_results']
grant = dec(lr['result'])
req = dec(lr['descriptor'])
assert lr['id'] == ids['Grant'] and req == grant['Request']
assert req['Intent'] == intents[ids['Grant']]
assert req['Operation'] == 'DelegateAuthority' and grant['Accepted'] and grant['Reason'] == 'delegated'
assert grant['Before'] == grant['After'] and grant['Before']['Projection']['Status'] == 'open'
claim = dec(req['Intent']['Payload']['claim']['Text'])
assert claim['OriginalRequestID'] == ids['BA'] and claim['ServiceOperationID'] == ids['Service']
service = intents[ids['Service']]
assert claim['Recipient'] == service['Run'] and claim['RecipientRevision'] == 1
assert service['Authority']['GrantID'] == ids['Grant'] and service['Authority']['Path'] == 'delegated'
assert service['Authority']['Root'] == claim['Root'] == grant['Authority']['Root']
assert claim['Root']['Permission'] in dec(t['composition_permissions'][0]['descriptor'])
p = claim['Root']['Permission']
assert (p['Kind'], p['Target'], p['Operation'], p['Obligation']) == (service['Kind'], service['Target'], 'apply', service['Obligation'])
roots = {r['id']: dec(r['descriptor']) for r in t['composition_roots']}
for run, expected in [(claim['Root']['RootRun'], grant['Authority']['RootID']), (claim['RecipientRoot'], grant['Authority']['RecipientRootID'])]:
    assert sha(canon(['alan/resource-profile/v1', [d['Namespace'], run]]).encode()) == expected
    assert roots[expected]['Run'] == run and run in roots[expected]['Members']
assert grant['Authority']['Namespace'] == d['Namespace']
assert len(t['lifecycle_results']) == 1 and all(i.get('Payload', {}).get('operation', {}).get('Text') != 'RevokeAuthority' for i in intents.values())
assert all(r['generation'] == 1 and r['revision'] == 1 and not r['retired'] for r in t['subjects'])
assert not t['composition_charges'] and not t['external_receipts'] and not t['adapter_results']
assert not any(f['Kind'] == 'authority_admission' for f in facts)
assert not any(i['Obligation'] in ['ba-reply', 'ab-reply'] for i in intents.values())

admitted = {}
for label in ['AB', 'BA', 'Ready']:
    i = intents[ids[label]]
    matches = [r for r in t['admissions'] if dec(r['descriptor']).get('IntentID') == ids[label]]
    ar, = matches
    desc, res = dec(ar['descriptor']), dec(ar['result'])
    assert res['Accepted'] and desc['Source'] == i['Run'] and desc['Target'] == i['Target']
    assert res['Bound'] == {'Subject': i['Target'], 'Generation': 1}
    assert ar['observation_id'] == res['ObservationID'] and desc['Payload'] == i['Payload']
    o = obs[res['ObservationID']]
    od = dec(o['descriptor'])
    assert od['Run'] == res['Bound'] and od['Source'] == i['Run']['Subject'] and od['SourceGeneration'] == 1
    assert od['Payload'] == i['Payload'] and od['Obligation'] == i['Obligation']
    admitted[label] = {'IntentID': ids[label], 'ObservationID': o['id'], 'ConsumedSequence': o['consumed_sequence']}
assert admitted['AB']['ConsumedSequence'] is None
assert admitted['BA']['ConsumedSequence'] == 2 and admitted['Ready']['ConsumedSequence'] == 3
assert history[('a', 1, 2)]['selected_observation'] == admitted['BA']['ObservationID']
assert history[('a', 1, 3)]['selected_observation'] == admitted['Ready']['ObservationID']
assert history[('a', 1, 4)]['selected_observation'] is None
ea, = t['external_attempts']
attempt = dec(ea['descriptor'])
assert ea['intent_id'] == ids['Service'] and attempt['Intent'] == service
assert not attempt['Allowed'] and attempt['Reason'] == 'authority_request_not_admitted'
assert not attempt['Authority']['Allowed'] and attempt['Authority']['Reason'] == attempt['Reason']
assert attempt['Authority']['OriginalRequestID'] == ids['BA'] and attempt['Authority']['ServiceOperationID'] == ids['Service']
auth, = [dec(r['result']) for r in t['attempts'] if r['intent_id'] == ids['Service']]
assert auth['Allowed'] and auth['Source']['Sequence'] == 4
order = lambda kind, subject=None, seq=None, op=None: next(f['Order'] for f in facts if f['Kind'] == kind and (subject is None or f['Run']['Subject'] == subject) and (seq is None or f['Sequence'] == seq) and (op is None or f['OperationID'] == op))
orders = {'ABCommit': order('execution', 'a', 1), 'BACommit': order('execution', 'b', 1), 'BASelection': order('execution', 'a', 2), 'GrantAccepted': order('lifecycle', op=ids['Grant']), 'ReadySelection': order('execution', 'a', 3), 'ServiceCommit': order('execution', 'a', 4), 'RefusalReservation': order('reservation', op=ids['Service'])}
assert list(orders.values()) == sorted(orders.values())
assert all(dec(history[(sub, 1, seq)]['successor'])['Status'] == 'open' for sub, seq in [('a', 4), ('b', 5)])
programs = {r['id']: r for r in t['programs']}
arun, = [r for r in t['runs'] if r['subject'] == 'a']
ap = programs[arun['program_id']]
assert sha(ap['bytes'].encode()) == ap['id']
code = dec(ap['bytes'])['Code']
run = dec((case / 'run.json').read_text())
endpoints = dec((case / 'endpoints.json').read_text())
assert all(v['Contacts'] == 0 and v['Applications'] == [] for v in endpoints.values())
refs = []
for table, key, column in [('history', ('a', 1, 2), 'selected_observation'), ('history', ('a', 1, 3), 'selected_observation')]:
    r = history[key]
    refs.append({'Table': table, 'Key': key, 'Column': column, 'Value': r[column], 'RawCellSHA256': sha(r[column].encode())})
for table, r, col in [('lifecycle_results', lr, 'result'), ('external_attempts', ea, 'descriptor')]:
    refs.append({'Table': table, 'Key': r.get('id', r.get('intent_id')), 'Column': col, 'RawCellSHA256': sha(r[col].encode())})
print(json.dumps({'Input': str(path), 'InputSHA256': sha(raw), 'Namespace': d['Namespace'], 'IDs': ids, 'GrantAccepted': grant['Accepted'], 'ExactClaim': claim, 'AdmissionJoins': admitted, 'CompositionOrders': orders, 'RuntimeAuthorityAllowed': attempt['Authority']['Allowed'], 'RuntimeReason': attempt['Reason'], 'NoRevocationIntentOrDecision': True, 'BothOriginalsCommittedAndUnpaid': True, 'ServiceSourceAuthorized': True, 'NoPhysicalGateOrApplicationReceipt': True, 'OriginalAProgramSHA256': ap['id'], 'ActualAReceiveInstructions': code[1:3], 'RunQuiescent': run['Report']['Quiescent'], 'RunOpenRuns': run['Report']['OpenRuns'], 'EndpointContacts': {k:v['Contacts'] for k,v in endpoints.items()}, 'RawCellReferences': refs}, indent=2))
