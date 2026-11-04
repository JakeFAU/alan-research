#!/usr/bin/env python3
"""Post-run, stdlib-only Phase3B.2 archive inspection; no Go or Replay execution.

Usage: python3 runtime-audit.py [evidence-directory] [output-json]
Writes deterministic JSON to stdout. Reads archive bytes; decompresses SQLite
backups only into temporary scratch files opened read-only and immutable.
This bounded archive checker is not part of the frozen Go qualification binary.
"""
import copy
import gzip
import hashlib
import io
import json
from pathlib import Path
import re
import sqlite3
import sys
import tarfile
import tempfile


def check(condition, message):
    if not condition:
        raise ValueError(message)


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def compact(value):
    # Go encoding/json's escaping, retaining struct field order from descriptors.
    return json.dumps(value, ensure_ascii=False, separators=(',', ':')).replace('&', r'\u0026').replace('<', r'\u003c').replace('>', r'\u003e').replace('\u2028', r'\u2028').replace('\u2029', r'\u2029').encode()


def desc(rows):
    return [json.loads(row['descriptor']) for row in rows]


def text(payload, key):
    return payload.get(key, {}).get('Text')


def selection(tables, row, before=None):
    o = json.loads(row['descriptor'])
    return sum(h['selected_observation'] == o['ID'] and h['subject'] == o['Run']['Subject'] and h['generation'] == o['Run']['Generation'] and h['sequence'] == row['consumed_sequence'] and (before is None or h['sequence'] < before) for h in tables['history']) == 1


def paid_edges(inspection):
    t = inspection['Tables']
    intents = {v['ID']: v for v in desc(t['intents'])}
    paid = 0
    for row in t['observations']:
        o = json.loads(row['descriptor'])
        if 'response_id' not in o['Payload']:
            continue
        check(o['Kind'] == 'message' and o['Provenance'] == 'ordinary-message' and o['SourceGeneration'] == o['Run']['Generation'] == 1 and text(o['Payload'], 'response_id') == o['Run']['Subject']+'-request', 'inauthentic final result')
        linked = False
        for admission in t['admissions']:
            if admission['observation_id'] != o['ID']:
                continue
            d = json.loads(admission['descriptor'])
            v = intents.get(d['IntentID'], {})
            linked |= (v.get('Run') == {'Subject': o['Source'], 'Generation': o['SourceGeneration']} and v.get('Target') == o['Run']['Subject'] and v.get('Generation') == o['Run']['Generation'] and v.get('Obligation') == o['Obligation'] and v.get('Payload') == o['Payload'])
        check(linked, 'final result lacks admitted original message')
        if row['consumed_sequence'] is not None:
            check(selection(t, row), 'final consumption lacks exact RunID/history selection')
            paid += 1
    return paid


def chain(inspection):
    """Focused structural custody checks; not a third full semantics proof."""
    t = inspection['Tables']
    intents = {v['ID']: v for v in desc(t['intents'])}
    observations = {json.loads(r['descriptor'])['ID']: r for r in t['observations']}
    attempts = {v['Intent']['ID']: v for v in desc(t['adapter_attempts'])}
    results = {r['intent_id']: json.loads(r['descriptor']) for r in t['adapter_results']}
    for row in t['observations']:
        reply = json.loads(row['descriptor'])
        if 'response_id' not in reply['Payload']:
            continue
        op = text(reply['Payload'], 'service_operation')
        v, a, r = intents.get(op), attempts.get(op), results.get(op)
        check(v and a and r, 'missing service operation')
        check(v['Kind'] == 'adapter' and v['Run'] == {'Subject': reply['Source'], 'Generation': reply['SourceGeneration']} and a['Intent'] == v and a['Allowed'] and a['Reason'] == 'authorized', 'service operation/authority mismatch')
        check(a['Prepared']['Request']['ID'] == r['Candidate']['ID'] == op and a['Prepared']['Request']['Key'] == r['Candidate']['Key'] and r['Disposition'] == r['Candidate']['Disposition'] == 'validated' and r['Values']['valid']['Flag'], 'service candidate mismatch')
        rd = copy.deepcopy(r); rd['Digest'] = ''
        check(r['Digest'] == sha(compact(rd)), 'service result digest mismatch')
        for field in ['mutation', 'total']:
            check(reply['Payload'][field] == r['Values'][field] == r['Candidate']['Payload'][field], 'unrelated result value')
        check(any(p['kind'] == 'adapter' and p['source'] == reply['Source'] and p['target'] == v['Target'] and p['operation'] == p['obligation'] == v['Obligation'] and p['revision'] == 1 for p in t['permissions']), 'missing principal permission tuple')
        auth = [json.loads(x['result']) for x in t['attempts'] if x['id'] == a['AuthorizationID']]
        check(any(x['Allowed'] and x['Reason'] == 'authorized' and x['IntentID'] == op and x['Source']['Run'] == v['Run'] and x['Source']['Revision'] == 1 for x in auth), 'missing source authorization')
        request = False
        for reqrow in t['observations']:
            req = json.loads(reqrow['descriptor'])
            if req['Run'] != v['Run'] or req['Source'] != reply['Run']['Subject'] or req['Payload'].get('request_id') != reply['Payload']['response_id']:
                continue
            if req['Kind'] != 'message' or req['Provenance'] != 'ordinary-message' or req['SourceGeneration'] != 1:
                continue
            for admission in t['admissions']:
                if admission['observation_id'] != req['ID']:
                    continue
                d = json.loads(admission['descriptor']); orig = intents.get(d['IntentID'], {})
                request |= (orig.get('Kind') == 'message' and orig.get('Run') == {'Subject': req['Source'], 'Generation': req['SourceGeneration']} and orig.get('Target') == req['Run']['Subject'] and orig.get('Generation') == req['Run']['Generation'] and orig.get('Obligation') == req['Obligation'] and orig.get('Payload') == req['Payload'] and selection(t, reqrow, v['Sequence']))
        check(request, 'missing original request selection')
        replies = [x for x in intents.values() if x['Run'] == {'Subject': reply['Source'], 'Generation': reply['SourceGeneration']} and x['Kind'] == 'message' and x['Payload'] == reply['Payload']]
        check(len(replies) == 1, 'ambiguous reply origin')
        settled = False
        for st in t['settlements']:
            obsrow = observations.get(st['observation_id'])
            if st['intent_id'] != op or not obsrow:
                continue
            obs = json.loads(obsrow['descriptor'])
            settled |= (json.loads(st['descriptor']) == r and obs['Operation'] == op and obs['Run'] == v['Run'] and obs['Source'] == v['Target'] and obs['Provenance'] == 'runtime-adapter' and obs['Payload'] == r['Values'] and selection(t, obsrow, replies[0]['Sequence']))
        check(settled, 'missing consumed original settlement')


def rejected(fn, value):
    try:
        fn(value)
    except ValueError:
        return True
    return False


def audit_run(run, files, kills, frozen):
    raw = {n[4:]: v for n, v in files.items() if n.startswith('raw/')}
    data = {n: json.loads(v) for n, v in raw.items() if n.endswith('.json')}
    def loadpath(path):
        return raw[Path(path).name]
    db_checks = 0
    def reconcile(path, inspection):
        nonlocal db_checks
        with tempfile.TemporaryDirectory(prefix='phase3b2-readonly-') as scratch:
            p = Path(scratch)/'backup.sqlite'; p.write_bytes(gzip.decompress(loadpath(path)))
            conn = sqlite3.connect(p.as_uri()+'?mode=ro&immutable=1', uri=True)
            conn.execute('PRAGMA query_only=ON')
            check(conn.execute('PRAGMA integrity_check').fetchone()[0] == 'ok', 'SQLite backup integrity')
            for table, expected in inspection['Tables'].items():
                check(re.fullmatch('[a-z_]+', table), 'unsafe table identifier')
                cur = conn.execute('SELECT * FROM '+table)
                columns = [x[0] for x in cur.description]
                rows = [dict(zip(columns, row)) for row in cur]
                for row in rows:
                    for key, value in row.items():
                        if isinstance(value, bytes): row[key] = value.decode()
                canonical = lambda values: sorted(json.dumps(v, sort_keys=True) for v in values)
                check(canonical(rows) == canonical(expected), 'SQLite/inspection table mismatch: '+table)
            conn.close()
        db_checks += 1
    manifests = [v for n,v in data.items() if '-sources-' in n]
    source_digests = set()
    for m in manifests:
        expected_sources = {path: v['sha256'] for path,v in frozen.items() if not set(Path(path).parts) & {'.git', 'evidence', 'consumer'}}
        check(m == expected_sources, 'source file differs from exact frozen helper subset')
        digest = sha(''.join(m[k]+'  '+k+'\n' for k in sorted(m)).encode())
        source_digests.add(digest)
        check(any('-sources-'+digest+'.json' in n for n in data), 'source manifest digest/name mismatch')
    check(len(source_digests) == 1, 'source manifest identity differs')
    source_digest = next(iter(source_digests))
    program_files = {n.rsplit('-program-',1)[1][:-5]: v for n,v in raw.items() if '-program-' in n and n.endswith('.json')}
    for pid, value in program_files.items():
        check(sha(value) == pid, 'raw retained program identity')
    starts = [v for n,v in data.items() if '-process-start-' in n]
    binaries = set()
    for p in starts:
        check(p['source_manifest_sha256'] == source_digest and len(p['executable_sha256']) == 64, 'child source/binary provenance')
        settings = {x['Key']: x.get('Value') for x in p['child_build_info']['Settings']}
        check((settings.get('-race') == 'true') == (run == 'race'), 'race build provenance')
        check(settings.get('-tags') == 'integration', 'integration build provenance')
        binaries.add(p['executable_sha256'])
    check(len(binaries) == 1, 'runtime binary differs within run')
    extracted = [json.loads(l.split('DURABLE ',1)[1]) for l in files['run.log'].decode().splitlines() if 'DURABLE ' in l]
    check(extracted == kills and len(kills) == 9, 'actual DURABLE kill records mismatch')
    case_wants = {'local':(2,1,1),'missing':(0,0,0),'claim':(0,0,0),'result_launder':(1,0,1),'mutual_wait':(0,0,0)}
    cases = [v for n,v in data.items() if n.endswith('-phase3b2-case.json')]
    check({c['scenario'] for c in cases} == set(case_wants) and len(cases) == 5, 'runtime case matrix')
    summary_cases = {}
    for c in cases:
        paid,ca,cb = case_wants[c['scenario']]
        for p in c['after']['Tables']['programs']:
            check(p['id'] in program_files and p['bytes'].encode() == program_files[p['id']], 'retained program bytes differ from raw artifact')
        chain(c['after']);check(paid_edges(c['after']) == c['paid_originals'] == paid and c['unpaid_originals'] == 2-paid, 'runtime case payment')
        for label, contacts in [('A',ca),('B',cb)]:
            w = c['endpoint_'+label]
            check(w['Contacts'] == len(w['Applications']) == contacts and w['Total'] == contacts, 'case endpoint contacts/application mismatch')
            check(sum(json.loads(line).get('Event') == 'ingress' for line in loadpath(c['ingress_'+label]).splitlines()) == contacts, 'case ingress count')
        reconcile(c['store'],c['after'])
        summary_cases[c['scenario']] = {'paid':paid,'unpaid':2-paid,'contacts':[ca,cb]}
    crash_wants = {'before-authorization':('phase3b2.manual','before-authorization',2),'after-source-authorization':('intent.authorized','authorized',2),'reserved-before-physical-admission':('dispatch.consider','pre-admission',0),'physical-admission-before-call':('adapter.authorized','pre-transport',0),'result-before-settlement':('adapter.result-recorded','pre-settlement',2),'settled-before-consumption':('adapter.settled','pre-consumption',2),'after-refusal':('phase3b2.manual','after-refusal',0),'mutual-wait':('phase3b2.manual','mutual-wait',0),'one-paid-one-pending':('intent.authorized','authorized',2)}
    crashes = [v for n,v in data.items() if n.endswith('-phase3b2-crash.json')]
    check({c['cut'] for c in crashes} == set(crash_wants) and len(crashes) == 9, 'crash matrix')
    summary_crashes = {}
    for c in crashes:
        boundary,phase,paid = crash_wants[c['cut']];b=c['applicable_boundary']
        check((b['Boundary'],b['Phase']) == (boundary,phase) and b['PID'] == c['crash_PID'] != c['recovery_PID'] and c['signal'] == 'SIGKILL', 'barrier/owner identity')
        k = [k for k in kills if k['pid'] == c['crash_PID']]
        check(len(k) == 1 and k[0]['exit_signal'] == 'killed' and k[0]['race_children'] == (run == 'race') and k[0]['binary_sha256'] in binaries, 'checked kill disposition')
        for key in ['Scenario','Token','Store','Boundary','Phase','Subject']:
            check(k[0]['barrier_config'][key] == b[key], 'kill/barrier configuration mismatch')
        check(any(p['PID'] == c['crash_PID'] for p in starts) and any(p['PID'] == c['recovery_PID'] for p in starts), 'missing distinct process launch provenance')
        recs = [v for n,v in data.items() if n.endswith('-process-'+str(c['recovery_PID'])+'.json')]
        check(len(recs)==1 and recs[0]['exit_code']==0 and recs[0]['race_children']==(run=='race') and recs[0]['binary_sha256'] in binaries, 'recovery process result')
        report = next(json.loads(l.split('PHASE3B2_RESULT ',1)[1]) for l in recs[0]['output'].splitlines() if l.startswith('PHASE3B2_RESULT '))
        check(report['PID']==c['recovery_PID'], 'recovery report PID')
        old,new=c['stopped']['Tables'],c['after']['Tables']
        for table in ['programs','subjects','permissions']:
            check(old[table] == new[table] == c['initial']['Tables'][table], 'immutable program/authority changed')
        for table in ['attempts','adapter_attempts','adapter_results','settlements']:
            check(all(r in new[table] for r in old[table]), 'retained decision/result prefix changed')
        epochs=[]
        for subject in ['a','b']:
            a=next(r for r in old['runs'] if r['subject']==subject and r['generation']==1)
            z=next(r for r in new['runs'] if r['subject']==subject and r['generation']==1)
            check(z['epoch']==a['epoch']+1, 'physical activation not renewed');epochs.append([a['epoch'],z['epoch']])
        chain(c['stopped']);chain(c['after']);check(paid_edges(c['after'])==c['paid_originals']==paid and c['unpaid_originals']==2-paid, 'crash payment')
        unknown=c['cut'] in ['reserved-before-physical-admission','physical-admission-before-call']
        check(c['allowed_without_result_remains_unknown']==unknown, 'unknown label')
        if unknown:
            check(report['Report']['UnknownAdapters']==[b['Intent']], 'unknown exact original identity')
            check(not old['adapter_results'] and not new['adapter_results'] and not old['settlements'] and not new['settlements'], 'unknown invented outcome')
            check(any(x['Intent']['ID']==b['Intent'] and x['Allowed'] for x in desc(old['adapter_attempts'])), 'unknown lacked retained reservation')
        else:
            check(not report['Report']['UnknownAdapters'], 'unexpected unknown original')
        contacts=[]
        for label in ['A','B']:
            before,after=c['world_'+label+'_cut'],c['world_'+label+'_after']
            contacts.append(after['Contacts'])
            check(after['Contacts']==len(after['Applications'])==after['Total'], 'crash endpoint counts')
            if unknown or c['cut'] in ['after-refusal','mutual-wait']:
                check(before==after and after['Contacts']==0, 'unpaid recovery contacted endpoint')
            else:
                check(after['Contacts']==1, 'known recovery endpoint count')
                if before['Contacts']:
                    check(before==after, 'known operation reissued')
        if c['cut']=='one-paid-one-pending':check(paid_edges(c['stopped'])==1, 'one-paid cut missing actual consumption')
        reconcile(c['cut_store'],c['stopped']);reconcile(c['after_store'],c['after'])
        summary_crashes[c['cut']]={'paid':paid,'unpaid':2-paid,'contacts':contacts,'epochs':epochs,'unknown_original':b['Intent'] if unknown else None}
    # Tie every validated candidate to exact configured resource/scope world key/value.
    endpoint_links=0
    for c in cases+crashes:
        t=c['after']['Tables'];attempts={x['Intent']['ID']:x for x in desc(t['adapter_attempts'])}
        for row in t['adapter_results']:
            r=json.loads(row['descriptor'])
            if r['Disposition']!='validated':continue
            a=attempts[row['intent_id']];principal=a['Intent']['Run']['Subject'];label='A' if principal=='a' else 'B'
            check(a['Intent']['Target']=='phase3b2.counter-'+label+('.red' if label=='A' else '.blue') and a['Prepared']['Model']=='counter-'+label+('/red' if label=='A' else '/blue'), 'configured resource/scope')
            configs = [v for n,v in data.items() if '-phase3b2-config-' in n]
            check(any(x['endpoint_'+label] == a['Prepared']['Endpoint'] for x in configs), 'prepared endpoint lacks fixed config provenance')
            w=c.get('endpoint_'+label,c.get('world_'+label+'_after'))
            apps=[x for x in w['Applications'] if x['Key']==r['Candidate']['Key']]
            check(len(apps)==1 and apps[0]['Mutation']==r['Values']['mutation']['Number'] and apps[0]['Total']==r['Values']['total']['Number'], 'endpoint original key/value provenance')
            endpoint_links+=1
    controls={v['mutation']:v for n,v in data.items() if '-oracle-' in n}
    names={'unrelated-result','missing-original-request','permission-tuple','all-selected-foreign-generation','final-consumption-foreign-generation'}
    check(set(controls)==names,'inspection control names')
    for name,c in controls.items():
        chain(c['before']);check(paid_edges(c['before'])==2,'control positive input')
        check(c['before']['Tables']['observations']==c['after']['Tables']['observations'] and len(c['after']['Tables']['adapter_results'])==c['adapter_results']==c['superficial_paid_observations']==2,'control surface changed')
        chain_bad=rejected(chain,c['after']);paid_bad=rejected(paid_edges,c['after'])
        check(chain_bad==(name!='final-consumption-foreign-generation'),'control intermediate outcome')
        check(paid_bad==(name in {'all-selected-foreign-generation','final-consumption-foreign-generation'}),'control final generation outcome')
        check((c['chain_error']!='<nil>')==chain_bad and (c['paid_error']!='<nil>')==paid_bad and c['validated_paid']==(0 if paid_bad else 2),'retained control diagnostics')
    rejections={v['mutation']:v for n,v in data.items() if n.endswith('-rejection.json')}
    check(set(rejections)=={'remove-principal-permission','launder-result-as-allow'},'Replay control names')
    for c in rejections.values():
        with tempfile.TemporaryDirectory(prefix='phase3b2-corrupt-readonly-') as scratch:
            p = Path(scratch)/'backup.sqlite'; p.write_bytes(gzip.decompress(loadpath(c['mutated_copy'])))
            conn = sqlite3.connect(p.as_uri()+'?mode=ro&immutable=1', uri=True)
            conn.execute('PRAGMA query_only=ON')
            if c['mutation'] == 'remove-principal-permission':
                check(conn.execute("SELECT count(*) FROM permissions WHERE kind='adapter' AND source='a'").fetchone()[0] == 0, 'disposable permission deletion absent')
            else:
                attempts = [json.loads(row[0]) for row in conn.execute('SELECT descriptor FROM adapter_attempts')]
                check(any(x['Intent']['Run']['Subject']=='a' and x['Allowed'] and x['Reason']=='grant_denied' for x in attempts), 'disposable contradictory attempt absent')
            conn.close()
        check(c['positive_source_unchanged'] and 'corrupt' in c['replay_error'] and Path(c['mutated_copy']).name in raw,'Replay rejection evidence')
    replays=[v for n,v in data.items() if '-replay-' in n and n.endswith('.json')]
    for r in replays:
        check(r['before']==r['after'] and r['endpoint_before']==r['endpoint_after'] and r['ingress_before']==r['ingress_after'],'Replay changed inspection/world/ingress')
        reconcile(r['before_local_store'],r['before']);reconcile(r['after_local_store'],r['after'])
    # Meaningful in-memory self-controls; archived bytes are never rewritten.
    positive=copy.deepcopy(controls['final-consumption-foreign-generation']['before'])
    final_ids={json.loads(r['descriptor'])['ID'] for r in positive['Tables']['observations'] if 'response_id' in json.loads(r['descriptor'])['Payload']}
    moved=0
    for h in positive['Tables']['history']:
        if h['selected_observation'] in final_ids:h['generation']+=1;moved+=1
    check(moved==2 and rejected(paid_edges,positive) and not rejected(chain,positive),'generation-only self-control')
    return {'raw_members':len(raw),'source_manifest_files':len(manifests[0]),'source_manifest_sha256':source_digest,'binary_sha256':next(iter(binaries)),'checked_kills':len(kills),'cases':summary_cases,'crashes':summary_crashes,'inspection_controls':sorted(controls),'disposable_Replay_controls':sorted(rejections),'readonly_Replay_records':len(replays),'sqlite_backups_reconciled':db_checks,'endpoint_original_links':endpoint_links,'negative_self_controls':['final-only-generation rejected; intermediate chain retained','artifact digest tamper rejected']}


def verify_digest(raw, expected):
    check(len(raw)==expected['bytes'] and sha(raw)==expected['sha256'],'artifact digest/size mismatch')


def main():
    directory=Path(sys.argv[1]) if len(sys.argv)>1 else Path(__file__).resolve().parent
    index_raw=(directory/'runtime-index.json').read_bytes();index=json.loads(index_raw)
    check(index['source_commit']=='f75e26984c84fe649f26d5e4eaf4ed70e7506e46','frozen commit identity')
    for name,expected in index['artifacts'].items():verify_digest((directory/name).read_bytes(),expected)
    frozen_doc=json.loads((directory/'frozen-source.json').read_bytes());frozen={x['path']:x for x in frozen_doc['Files']}
    check(frozen_doc['ExecutableSourceCommit']==index['source_commit'],'source identity conflict')
    result={'format':1,'scope':'post-run archive inspection, not a third full semantics proof or whole-history enumeration','executable_source_commit':index['source_commit'],'audit_script_sha256':sha(Path(__file__).read_bytes()),'index_sha256':sha(index_raw),'runs':{}}
    for run,a in sorted(index['archives'].items()):
        archive=(directory/a['file']).read_bytes();verify_digest(archive,a)
        tampered = bytearray(archive); tampered[len(tampered)//2] ^= 1
        try:verify_digest(tampered,a)
        except ValueError:pass
        else:raise ValueError('digest tamper self-control accepted')
        with tarfile.open(fileobj=io.BytesIO(archive),mode='r:gz') as tar:
            members=tar.getmembers();check([m.name for m in members]==sorted(m.name for m in members) and len({m.name for m in members})==len(members),'archive member ordering/duplicates')
            files={}
            for m in members:
                check(m.isfile() and m.uid==m.gid==m.mtime==0 and m.mode==0o644 and not m.uname and not m.gname and not Path(m.name).is_absolute() and '..' not in Path(m.name).parts,'archive metadata/path')
                files[m.name]=tar.extractfile(m).read()
        check(set(files)=={m['path'] for m in a['members']},'archive index coverage')
        for m in a['members']:verify_digest(files[m['path']],m)
        kills=json.loads((directory/('kill-records-'+run+'.json')).read_bytes())
        result['runs'][run]=audit_run(run,files,kills,frozen)
    rendered = json.dumps(result,sort_keys=True,indent=2)+'\n'
    if len(sys.argv)>2:
        Path(sys.argv[2]).write_text(rendered)
    else:
        sys.stdout.write(rendered)


if __name__=='__main__':
    main()
