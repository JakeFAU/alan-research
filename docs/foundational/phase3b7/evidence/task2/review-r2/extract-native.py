"""Read only JSON/JSONL; preserve copied bytes and descriptor strings in a test excerpt."""
import base64
import gzip
import hashlib
import json
from pathlib import Path

out = Path('docs/foundational/phase3b7/evidence/task2/review-r2')
entries, inputs = {}, {}
service = '["a",1,4,0]'
for mode in ('normal', 'race'):
    for branch in ('zero', 'applied'):
        root = Path('/tmp/alan-phase3b7/raw-' + mode + '-v1')
        files = {}
        for suffix in ('first-owner/boundaries.jsonl', 'first-owner/pre-call.json', 'stopped/inspect.json'):
            path = root / branch / suffix
            data = path.read_bytes()
            inputs[str(path)] = {'bytes': len(data), 'sha256': hashlib.sha256(data).hexdigest()}
            files[suffix] = data
        raw = json.loads(files['stopped/inspect.json'])
        attempts = [r['descriptor'] for r in raw['Tables']['external_attempts'] if r['intent_id'] == service]
        gates = [r['descriptor'] for r in raw['Tables']['composition_facts'] if json.loads(r['descriptor']).get('OperationID') == service and json.loads(r['descriptor']).get('AuthorityAdmission')]
        assert len(attempts) == len(gates) == 1
        entries[mode + '/' + branch] = {
            'Boundaries': base64.b64encode(files['first-owner/boundaries.jsonl']).decode(),
            'Worker': base64.b64encode(files['first-owner/pre-call.json']).decode(),
            'Attempt': attempts[0], 'Gate': gates[0],
        }
encoded = json.dumps(entries, sort_keys=True, separators=(',', ':')).encode()
Path('research/matchedcapture/testdata/native-authorization.json.gz').write_bytes(gzip.compress(encoded, mtime=0))
(out / 'native-inputs.json').write_text(json.dumps({'inputs': inputs, 'excerpt_sha256': hashlib.sha256(encoded).hexdigest()}, indent=2) + '\n')
