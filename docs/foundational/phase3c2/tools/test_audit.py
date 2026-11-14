"""Behavioral controls: deleting relation, provenance or schema checks must fail."""
import pathlib
import subprocess
import sys
import unittest

ROOT = pathlib.Path(__file__).resolve().parents[4]
TOOLS = pathlib.Path(__file__).resolve().parent
EXPORT = ROOT / 'docs/foundational/phase3c2/evidence/task1/export.json.gz'


class AuditAcceptance(unittest.TestCase):
    def test_genuine_export_is_independently_accepted(self):
        result = subprocess.run([sys.executable, str(TOOLS / 'audit.py'), str(EXPORT)],
                                capture_output=True, text=True, cwd=ROOT)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn('independent_abstract_graphs', result.stdout)

# These are forged evidence controls, never asserted historical executions.
def hostile_cases(document):
    import copy
    import json

    def changed(name, path, value, reason):
        d = copy.deepcopy(document)
        target = d
        for part in path[:-1]:
            target = target[part]
        target[path[-1]] = value
        return name, json.dumps(d).encode(), reason

    for name, path, value, reason in [
        ('wrong-schema', ['Schema'], 'unknown', 'Schema'),
        ('unknown-top-field', ['Invented'], 0, 'keys'),
        ('boolean-count', ['Graphs', 0, 'States'], True, 'States'),
        ('float-count', ['Graphs', 0, 'Arcs'], 1.0, 'Arcs'),
        ('missing-config', ['Graphs'], document['Graphs'][:-1], 'Graphs'),
        ('wrong-completion', ['Graphs', 0, 'Nodes', 0, 'CanComplete'], False, 'CanComplete'),
        ('wrong-parent', ['Graphs', 0, 'Nodes', 1, 'Parent'], 1, 'Parent'),
        ('aggregate-holder', ['Graphs', 0, 'Nodes', 0, 'State', 'Held'], 0, 'Held'),
        ('changed-holder', ['Witnesses', 0, 'Steps', 6, 'State', 'Held'], [False, True], 'Held'),
        ('invented-membership', ['Groupings', 0, 'A', 'Memberships', 0, 'Root'], 'Q-new', 'Memberships'),
        ('conflicting-limit', ['Groupings', 0, 'AB', 'Roots', 0, 'Limit'], 2, 'Roots'),
        ('copied-pool', ['Groupings', 0, 'Left', 'Roots'], [{'ID':'Q','Limit':1}, {'ID':'Q','Limit':1}], 'Roots'),
        ('capacity-overflow', ['Witnesses', 0, 'Steps', 6, 'State', 'Held'], [True, True], 'Held'),
        ('early-release-poststate', ['Witnesses', 7, 'Steps', 7, 'State', 'Held'], [False, False], 'Held'),
        ('wrong-actor', ['Witnesses', 0, 'Steps', 0, 'Event', 'Actor'], 'C', 'actor'),
        ('wrong-activation', ['Witnesses', 7, 'Steps', 7, 'State', 'Base', 'Agents', 1, 'Activation'], 2, 'Activation'),
        ('forged-binding', ['Witnesses', 0, 'Steps', 0, 'Event', 'Binding', 'Grant'], 'forged', 'Binding'),
        ('forged-receipt', ['Witnesses', 8, 'Steps', 13, 'Event', 'Result'], 'r2', 'result provenance'),
        ('poststate-invention', ['Witnesses', 0, 'Steps', 0, 'State', 'Base', 'Originals', 1, 'Committed'], True, 'Committed'),
        ('invented-initial-receipt', ['Witnesses', 0, 'Initial', 'Base', 'Originals', 0, 'Recorded'], 'r1', 'Initial'),
        ('unknown-event-field', ['Witnesses', 0, 'Steps', 0, 'Event', 'Activation'], 0, 'keys'),
        ('boolean-generation', ['Graphs', 0, 'Config', 'Base', 'Agents', 0, 'Generation'], True, 'Generation'),
        ('unknown-nested-field', ['Graphs', 0, 'Nodes', 0, 'State', 'Occupancy'], 0, 'keys'),
        ('null-array', ['Groupings', 0, 'C', 'Roots'], None, 'Roots'),
        ('probe-row-omitted', ['Probes', 'Rows'], document['Probes']['Rows'][:-1], 'probe roster'),
        ('probe-total-forged', ['Probes', 'Lifts'], 0, 'probe total'),
        ('probe-boolean', ['Probes', 'OldEdgeStates'], True, 'OldEdgeStates'),
        ('source-forged', ['Sources', 0, 'SHA256'], '0'*64, 'Sources'),
        ('missing-projection', ['Witnesses', 0, 'Projection'], [], 'Projection'),
        ('wrong-final-label', ['Witnesses', 0, 'CanComplete'], True, 'CanComplete'),
    ]:
        yield changed(name, path, value, reason)

    d = copy.deepcopy(document)
    graph = d['Graphs'][0]
    for index, node in enumerate(graph['Nodes']):
        for arc_index, arc in enumerate(node['Arcs']):
            target = graph['Nodes'][arc['To']]
            if target['Parent'] != index or target['Via'] != arc['Action']:
                # Remove a non-tree arc AND adjust its count: counts agree,
                # all BFS parents survive, only complete relation checking fails.
                del node['Arcs'][arc_index]
                graph['Arcs'] -= 1
                yield 'removed-non-tree-arc', json.dumps(d).encode(), 'Arcs'
                break
        else:
            continue
        break
    else:
        raise AssertionError('fixture must contain a non-tree arc')

    for name, witness, index, step, reason in [
        ('duplicate-reservation', 0, 7, document['Witnesses'][0]['Steps'][6], 'reservation'),
        ('duplicate-application', 8, 13, document['Witnesses'][8]['Steps'][12], 'apply'),
        ('premature-ack', 8, 20, document['Witnesses'][8]['Steps'][21], 'durable ACK'),
        ('premature-receipt-release', 8, 20, document['Witnesses'][8]['Steps'][22], 'release ACK'),
        ('consume-before-record', 8, 20, document['Witnesses'][8]['Steps'][24], 'consume'),
        ('actor-down', 8, 8, document['Witnesses'][8]['Steps'][12], 'actor down'),
    ]:
        d = copy.deepcopy(document)
        d['Witnesses'][witness]['Steps'].insert(index, step)
        yield name, json.dumps(d).encode(), reason

    d = copy.deepcopy(document)
    del d['Witnesses'][8]['Steps'][12]
    yield 'missing-genuine-apply', json.dumps(d).encode(), 'send-result custody'
    d = copy.deepcopy(document)
    d['Witnesses'][0]['Steps'][6]['Event']['Kind'] = 'release-capacity'
    yield 'standalone-capacity-release', json.dumps(d).encode(), 'event kind'
    # Legal but non-corresponding r2 history: a transport-only retry changes
    # the requested paired witness, without changing completion or projection.
    import audit
    d = copy.deepcopy(document)
    w = d['Witnesses'][9]
    current = w['Steps'][-1]['State']
    for kind in ('send-request', 'deduplicate'):
        event = dict(Kind=kind, Actor='A' if kind == 'send-request' else 'B',
                     Original='AB', Binding=audit.binding(0), Result='', Duplicate=False)
        current = audit.step(w['Config'], current, event, 'paired-control')
        w['Steps'].append(dict(Event=event, State=current))
    yield 'noncorresponding-r2-history', json.dumps(d).encode(), 'corresponding recovery histories'

    # Regenerate both valid histories with consumption before ACK/release.
    # Durable custody is sound, but this no longer witnesses release-before-consume.
    d = copy.deepcopy(document)
    for index in (8, 9):
        w = d['Witnesses'][index]
        events = [s['Event'] for s in w['Steps']]
        for consume, ack in ((41, 38), (24, 21)):
            events.insert(ack, events.pop(consume))
        current = audit.initial()
        w['Steps'] = []
        for event in events:
            current = audit.step(w['Config'], current, event, 'order-control')
            w['Steps'].append(dict(Event=event, State=current))
    yield 'wrong-recovery-custody-order', json.dumps(d).encode(), 'release before consumption'

    raw = json.dumps(document).encode()
    yield 'duplicate-json-key', raw.replace(b'"Schema":', b'"Schema":"duplicate", "Schema":', 1), 'duplicate key'
    yield 'nonfinite-json-number', raw.replace(b'"States":', b'"States":NaN, "Unused":', 1), 'nonfinite'


class HostileEvidence(unittest.TestCase):
    def test_hostile_records_rejected_for_specific_reasons(self):
        import audit
        document = audit.load(EXPORT)
        for name, raw, reason in hostile_cases(document):
            with self.subTest(control=name):
                with self.assertRaisesRegex(audit.InvalidEvidence, reason):
                    audit.check(audit.decode(raw), ROOT)

    def test_plain_json_and_gzip_have_identical_meaning(self):
        import audit
        import gzip
        self.assertEqual(audit.load(EXPORT), audit.decode(gzip.decompress(EXPORT.read_bytes())))


class IndependentlyGeneratedHistories(unittest.TestCase):
    def test_completion_under_every_configuration_and_both_results(self):
        # Disabling reservation, atomic release, crash retention, actual result
        # custody or durable ACK must invalidate these separately authored plans.
        import audit
        histories = generated_histories()
        self.assertEqual(len(histories), 144)
        for w in histories:
            with self.subTest(config=w['Name']):
                g, indexes = audit.graph(w['Config'])
                final = audit.replay(w, w['Config'], g, indexes, w['Name'])
                self.assertEqual(final['Held'], [False, False])
                for original in final['Base']['Originals'][:2]:
                    self.assertTrue(original['Consumed'])
                    self.assertEqual(original['Recorded'], w['Name'].split('/')[-1])
                    self.assertEqual(original['Retained'], '')
                self.assertEqual([a['Activation'] for a in final['Base']['Agents']], [1, 1, 1])

    def test_local_transport_guards_reject_impossible_events(self):
        import audit
        for kind in ('send-request', 'admit', 'apply', 'receive-result', 'record',
                     'consume', 'send-ack', 'release', 'recover'):
            c = audit.config('shared-1', 'early', 0)
            caller = kind in ('send-request', 'receive-result', 'record', 'consume', 'send-ack')
            physical = kind == 'recover'
            event = dict(Kind=kind, Actor='A' if caller or physical else 'B',
                         Original='' if physical else 'AB',
                         Binding=audit.zero_binding() if physical else audit.binding(0),
                         Result='r1' if kind in ('apply', 'receive-result', 'record', 'consume') else '',
                         Duplicate=False)
            with self.subTest(kind=kind):
                with self.assertRaises(audit.InvalidEvidence):
                    audit.step(c, audit.initial(), event, 'empty-custody')


def generated_histories():
    """Authored BC-first completion plans; no exported events/states are used.

    This fixture generator shares the Python replay transitions, not a second
    independent checker. It is independent of the Go producer and its histories.
    """
    import audit
    result = []
    for layout in ('shared-1', 'shared-2', 'separate-1'):
        for policy in ('early', 'ready-now', 'stable-ready'):
            for mask in range(8):
                for value in ('r1', 'r2'):
                    c = audit.config(layout, policy, mask)
                    current = audit.initial()
                    steps, projection = [], []

                    def emit(kind, i, duplicate=False, physical=None):
                        nonlocal current
                        caller = kind in ('commit', 'send-request', 'receive-result',
                                          'record', 'consume', 'send-ack')
                        event = dict(Kind=kind, Actor=physical or 'ABC'[i + (0 if caller else 1)],
                                     Original='' if physical else ('AB', 'BC')[i],
                                     Binding=audit.zero_binding() if physical else audit.binding(i),
                                     Result=value if kind in ('apply', 'send-result', 'receive-result',
                                                              'record', 'consume') else '',
                                     Duplicate=duplicate)
                        current = audit.step(c, current, event, 'generated')
                        steps.append(dict(Event=event, State=current))
                        if kind in ('commit', 'reserve-capacity', 'apply', 'consume'):
                            projection.append(audit.action(kind, i))

                    # Complete BC first; hold interruption and both pre-record
                    # receipt losses are chosen explicitly by this test plan.
                    for i in (1, 0):
                        for kind in ('commit', 'send-request', 'admit', 'reserve-capacity'):
                            emit(kind, i)
                        if i == 1:
                            emit('crash', i, physical='C')
                            emit('recover', i, physical='C')
                        emit('send-request', i)
                        emit('deduplicate', i, True)
                        emit('deduplicate', i)
                        emit('apply', i)
                        emit('send-result', i)
                        emit('receive-result', i, True)
                        emit('receive-result', i)
                        emit('crash', i, physical='AB'[i])
                        emit('recover', i, physical='AB'[i])
                        for kind in ('send-result', 'receive-result', 'record', 'send-ack'):
                            emit(kind, i)
                        emit('release', i, True)
                        emit('release', i)
                        emit('consume', i)
                        # Exercise retries after consumption and a fresh ACK
                        # after receipt release, without recreating any result.
                        for kind in ('send-request', 'deduplicate', 'send-ack', 'release'):
                            emit(kind, i)
                    result.append(dict(Name=audit.config_id(c) + '/' + value, Config=c,
                                       Initial=audit.initial(), Steps=steps, Projection=projection,
                                       Complete=True, CanComplete=True))
    return result


if __name__ == '__main__':
    unittest.main()
