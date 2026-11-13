#!/usr/bin/env python3
"""Reproduce only the reviewed inventory wording delta and command coverage."""
import gzip,hashlib,json,os,pathlib,subprocess,time
root=pathlib.Path.cwd();ev=root/'docs/foundational/phase3b6/evidence/task3';fix=ev/'review-i1';tmp=pathlib.Path('/tmp/alan-phase3b6/task3');tmp.mkdir(exist_ok=True)
correction=json.loads((fix/'correction.json').read_text());prior=json.loads((fix/'prior-verdicts.json').read_text());recorded=json.loads((ev/'commands.json').read_text());env=os.environ.copy();env.update(recorded['environment']);commands=[];sha=lambda b:hashlib.sha256(b).hexdigest()
def run(args,label=None):
 start=time.monotonic();r=subprocess.run([str(a) for a in args],cwd=root,env=env,capture_output=True)
 commands.append({'command':[str(a) for a in args],'cwd':str(root),'exit':r.returncode,'seconds':time.monotonic()-start})
 if label:(fix/(label+'.log')).write_bytes(r.stdout+r.stderr)
 assert r.returncode==0,(args,r.stdout.decode(),r.stderr.decode())
 return r.stdout
run(['go','test','-count=1','-v','./research/delegationcorrespondence/cmd/effectobservation'],'command-tests')
run(['go','build','-o',tmp/'effectobservation-review-i1','./research/delegationcorrespondence/cmd/effectobservation'],'build')
corpus=pathlib.Path(env['ALAN_EFFECT_OBSERVATION_CORPUS']);manifest=json.loads((corpus/'manifest.json').read_text());summaries=[];comparisons=[]
for old_summary in prior:
 mode,name=old_summary['mode'],old_summary['case'];filename=old_summary['file'];path=corpus/('raw-'+mode)/manifest['raw-'+mode+'.tar.gz']['cases'][name]
 # Frozen Git bytes preserve the earlier artifact, independently of regenerated files.
 old_compressed=run(['git','show',correction['prior_commit']+':docs/foundational/phase3b6/evidence/task3/reports/'+filename]);old=gzip.decompress(old_compressed)
 assert sha(old)==old_summary['report_sha256'] and len(old)==old_summary['report_bytes']
 first=run([tmp/'effectobservation-review-i1','-case',path]);second=run([tmp/'effectobservation-review-i1','-case',path]);assert first==second
 assert first==old.replace(correction['old_phrase'].encode(),correction['new_phrase'].encode()),'report changed beyond the reviewed wording'
 (ev/'reports'/filename).write_bytes(gzip.compress(first,mtime=0))
 new_summary=dict(old_summary);new_summary.update(report_sha256=sha(first),report_bytes=len(first));summaries.append(new_summary)
 comparisons.append({'mode':mode,'case':name,'old_sha256':sha(old),'new_sha256':sha(first),'old_bytes':len(old),'new_bytes':len(first),'only_reviewed_phrase_changed':True,'deterministic':True})
(ev/'verdicts.json').write_text(json.dumps(summaries,indent=2)+'\n')
hashes=json.loads((fix/'prior-source-sha256.json').read_text());changed=[]
for p,old_hash in hashes.items():
 new_hash=sha((root/p).read_bytes())
 if new_hash!=old_hash:changed.append(p)
 hashes[p]=new_hash
assert changed==[correction['source_path']],changed
(ev/'source-sha256.json').write_text(json.dumps(hashes,indent=2)+'\n')
assert not run(['gofmt','-l',correction['source_path']],'format')
(fix/'report-comparison.json').write_text(json.dumps(comparisons,indent=2)+'\n')
(fix/'commands.json').write_text(json.dumps({'environment':recorded['environment'],'commands':commands},indent=2)+'\n')
print('PASS existing command/serialization tests;4 reports exported twice; every output differs only by the reviewed phrase; only effect_protocol.go source hash changed.')
