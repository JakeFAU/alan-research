"""Semantic corruption controls: each mutation breaks a stated evidence contract."""
import copy
import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

import audit

SAMPLE = Path(__file__).resolve().parents[1] / 'evidence/task1/export.json'


def controls(source):
    """Separate hostile inputs; never mutate the genuine export."""
    cases = []

    def add(name, reason, mutate):
        data = copy.deepcopy(source)
        mutate(data)
        cases.append((name, reason, data))

    add('schema-version', 'schema', lambda d: d.update(Schema='wrong'))
    add('missing-field', 'schema', lambda d: d['Graphs'][0].pop('States'))
    add('bool-as-count', 'schema', lambda d: d['Graphs'][0].update(States=True))
    add('null-array', 'schema', lambda d: d['Graphs'][0]['Nodes'][0].update(Arcs=None))
    add('missing-config', 'config roster', lambda d: d['Graphs'].pop())
    add('forged-binding', 'binding', lambda d: d['Graphs'][0]['Nodes'][0]['Arcs'][0]['Action']['Binding'].update(Grant='forged'))
    add('forged-config', 'config', lambda d: d['Graphs'][0]['Config']['Agents'][1].update(Generation=2))
    add('missing-arc', 'transition relation', lambda d: d['Graphs'][0]['Nodes'][0]['Arcs'].pop())
    add('extra-arc', 'transition relation', lambda d: d['Graphs'][0]['Nodes'][0]['Arcs'].append(copy.deepcopy(d['Graphs'][0]['Nodes'][0]['Arcs'][0])))

    def incoming(d):
        g = d['Graphs'][0]
        # Remove a non-tree incoming edge: all nodes still have their BFS witness.
        for i, n in enumerate(g['Nodes']):
            for j, arc in enumerate(n['Arcs']):
                if g['Nodes'][arc['To']]['Parent'] != i:
                    n['Arcs'].pop(j)
                    g['Arcs'] -= 1
                    return
        raise AssertionError('fixture has no non-tree incoming edge')
    add('missing-incoming-arc', 'transition relation', incoming)

    def guard(d):
        g = d['Graphs'][8]
        dead = next(n for n in g['Nodes'] if n['State']['Stages'] == ['W', 'W', 'W'])
        target = next(i for i, n in enumerate(g['Nodes']) if n['State']['Stages'] == ['P', 'W', 'W'])
        dead['Arcs'].append({'Action': {'Kind': 'apply', 'Original': 'AB', 'Binding': copy.deepcopy(g['Config']['Originals'][0])}, 'To': target})
    add('wrong-serial-guard', 'transition relation', guard)
    add('wrong-complete', 'completion', lambda d: d['Graphs'][0]['Nodes'][0].update(Complete=True))
    add('wrong-can-complete', 'progress verdict', lambda d: d['Graphs'][8]['Nodes'][0].update(CanComplete=False))
    add('wrong-count', 'counts', lambda d: d['Graphs'][0].update(Arcs=0))
    add('wrong-parent', 'BFS predecessor', lambda d: d['Graphs'][0]['Nodes'][1].update(Parent=1))
    add('unreachable-state', 'reachable state set', lambda d: d['Graphs'][0]['Nodes'][0]['State'].update(Stages=['D', 'D', '']))
    add('dependency-erasure', 'dependency', lambda d: d['Graphs'][8]['Certificate'].update(Dependencies=[]))
    add('dependency-order', 'dependency', lambda d: d['Graphs'][0]['Certificate'].update(Order=['AB', 'BC']))
    add('dropped-hidden-debt', 'grouping', lambda d: d['Groupings'][0]['AB']['Originals'].pop(0))
    add('collapsed-availability', 'grouping', lambda d: d['Groupings'][0]['AB']['Owners'].pop())
    add('erased-port', 'grouping', lambda d: d['Groupings'][0]['A'].update(Ports=[]))
    add('wrong-grouped-config', 'grouping', lambda d: d['Groupings'][0]['LeftClosed']['Originals'].pop())
    add('forged-event-binding', 'binding', lambda d: d['Witnesses'][0]['Steps'][0]['Event']['Binding'].update(Operation=2))
    add('invented-initial-receipt', 'initial state', lambda d: d['Witnesses'][0]['Initial']['Originals'][0].update(Recorded='r1'))
    add('post-state-invention', 'post-state', lambda d: d['Witnesses'][0]['Steps'][0]['State']['Originals'][0].update(Recorded='r1'))
    add('hidden-application-erasure', 'projection', lambda d: d['Witnesses'][1]['Projection'].pop())
    add('wrong-witness-progress', 'witness progress', lambda d: d['Witnesses'][4].update(CanComplete=True))

    def premature(d):
        step = d['Witnesses'][-1]['Steps'][4]
        step['Event'].update(Kind='release', Result='')
    add('premature-release', 'release requires recorded ACK', premature)

    def repeat(d):
        w = d['Witnesses'][-1]
        w['Steps'].insert(4, copy.deepcopy(w['Steps'][3]))
    add('repeated-application', 'repeated application', repeat)

    def unavailable(d):
        w = d['Witnesses'][-1]
        # Actor A remains down, then receives a result: bypassing physical recovery.
        w['Steps'].pop(7)
        for step in w['Steps'][7:]:
            step['State']['Agents'][0].update(Down=True, Activation=0)
    add('down-actor', 'actor unavailable', unavailable)
    add('invented-receipt', 'receipt custody', lambda d: d['Witnesses'][0]['Steps'][5]['Event'].update(Result='r2'))
    add('wrong-collision', 'collision', lambda d: d['Collisions'][0].update(HiddenOriginal='BC'))
    add('wrong-retry', 'retry', lambda d: d['Retry'].update(ToState=4))
    add('missing-loss-witness', 'loss/release witness', lambda d: d['Witnesses'][-1].update(Steps=[], Projection=[], Complete=False, CanComplete=True))
    add('overstated-claim', 'claim classification', lambda d: d['Claims'][0].update(Status='bounded-exhaustive'))
    return cases


class AuditTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.source = json.loads(SAMPLE.read_text())

    def test_genuine_full_evidence(self):
        result = audit.audit(self.source)
        self.assertEqual(result['status'], 'accepted')
        self.assertEqual(result['graphs'], 16)
        self.assertEqual(result['states'], 639)
        self.assertEqual(result['arcs'], 1237)
        self.assertEqual(result['cannot_complete'], 1)
        self.assertEqual(result['witness_steps'], 54)

    def test_hostile_controls(self):
        for name, reason, data in controls(self.source):
            with self.subTest(name=name):
                with self.assertRaisesRegex(audit.AuditError, reason):
                    audit.audit(data)

    def test_loss_witness_really_loses_unrecorded_receipt(self):
        data = copy.deepcopy(self.source)
        w = data['Witnesses'][-1]
        w.update(Steps=[], Projection=[], Complete=False, CanComplete=True)
        with self.assertRaisesRegex(audit.AuditError, 'loss/release witness'):
            audit.audit(data)

    def test_duplicate_channel_and_post_consumption_retries(self):
        w = self.source['Witnesses'][-1]
        ledger = audit.Ledger(w['Config'])
        # Run through first receipt delivery, retaining that delivery's channel.
        for step in w['Steps'][:5]:
            ledger.step(step['Event'])
        event = copy.deepcopy(w['Steps'][5]['Event'])
        event['Duplicate'] = True
        ledger.step(event)
        self.assertEqual(ledger.state['Originals'][0]['ResultChannel'], 'r1')
        for i in (10, 11, 12, 13):
            ledger.step(w['Steps'][i]['Event'])
        # Authentic old channel receipt remains serviceable after consumption.
        event['Duplicate'] = False
        ledger.step(event)
        self.assertEqual(ledger.state['Originals'][0]['ResultChannel'], '')
        self.assertEqual(ledger.state['Originals'][0]['Recorded'], 'r1')
        # Send/receive ACK again after release; tombstone prevents reapplication.
        ledger.step(w['Steps'][11]['Event'])
        ledger.step(w['Steps'][12]['Event'])
        self.assertTrue(ledger.state['Originals'][0]['Applied'])
        with self.assertRaisesRegex(audit.AuditError, 'repeated application'):
            ledger.step(w['Steps'][3]['Event'])

    def test_protocol_guards_and_shared_owner_roles(self):
        w = self.source['Witnesses'][-1]
        ledger = audit.Ledger(w['Config'])
        with self.assertRaisesRegex(audit.AuditError, 'protocol guard'):
            ledger.step(w['Steps'][1]['Event'])  # no durable request yet
        for step in w['Steps'][:2]:
            ledger.step(step['Event'])
        with self.assertRaisesRegex(audit.AuditError, 'protocol guard'):
            ledger.step(w['Steps'][1]['Event'])  # occupied request channel
        admission = copy.deepcopy(w['Steps'][2]['Event'])
        admission['Duplicate'] = True
        ledger.step(admission)
        self.assertTrue(ledger.state['Originals'][0]['RequestChannel'])
        admission['Kind'] = 'deduplicate'
        ledger.step(admission)
        self.assertTrue(ledger.state['Originals'][0]['RequestChannel'])
        admission['Duplicate'] = False
        ledger.step(admission)
        self.assertFalse(ledger.state['Originals'][0]['RequestChannel'])
        ledger.step(w['Steps'][3]['Event'])
        ledger.step(w['Steps'][4]['Event'])
        with self.assertRaisesRegex(audit.AuditError, 'protocol guard'):
            ledger.step(w['Steps'][4]['Event'])  # occupied result channel
        ledger.step(w['Steps'][5]['Event'])
        ledger.step(w['Steps'][10]['Event'])
        with self.assertRaisesRegex(audit.AuditError, 'protocol guard'):
            ledger.step(w['Steps'][10]['Event'])  # cannot replace durable receipt
        ledger.step(w['Steps'][11]['Event'])
        release = copy.deepcopy(w['Steps'][12]['Event'])
        release['Duplicate'] = True
        ledger.step(release)
        self.assertTrue(ledger.state['Originals'][0]['ACKChannel'])
        ledger.step(release)  # repeated ACK release does not erase tombstone
        self.assertTrue(ledger.state['Originals'][0]['Applied'])
        crash = copy.deepcopy(w['Steps'][6]['Event'])
        crash['Actor'] = 'B'
        ledger.step(crash)
        with self.assertRaisesRegex(audit.AuditError, 'actor unavailable'):
            ledger.step(release)  # B as callee on AB
        bc_commit = self.source['Witnesses'][0]['Steps'][1]['Event']
        with self.assertRaisesRegex(audit.AuditError, 'actor unavailable'):
            ledger.step(bc_commit)  # same physical B as caller on BC

    def test_compatible_join_rejects_resource_and_policy_aliases(self):
        group = self.source['Groupings'][0]
        left, right = copy.deepcopy(group['A']), copy.deepcopy(group['B'])
        right['Originals'][0]['Resource'] = 'counter-B'
        with self.assertRaisesRegex(audit.AuditError, 'resource alias'):
            audit.join(left, right)
        right = copy.deepcopy(group['B'])
        right['Owners'][0]['Servicing'] = True
        with self.assertRaisesRegex(audit.AuditError, 'endpoint policy'):
            audit.join(left, right)
        with self.assertRaisesRegex(audit.AuditError, 'duplicated ownership'):
            audit.join(left, left)

    def test_input_remains_unchanged(self):
        before = copy.deepcopy(self.source)
        audit.audit(self.source)
        self.assertEqual(before, self.source)

    def test_cli_rejects_invalid_json(self):
        Path('/tmp/alan-phase3c1').mkdir(parents=True, exist_ok=True)
        with tempfile.TemporaryDirectory(dir='/tmp/alan-phase3c1') as directory:
            path = Path(directory) / 'invalid.json'
            nonfinite = copy.deepcopy(self.source)
            nonfinite['Graphs'][0]['States'] = float('nan')
            cases = [('{', 'Expecting property name'),
                     ('{"Schema":"first","Schema":"second"}', 'duplicate JSON key'),
                     (json.dumps(nonfinite), 'wrong primitive type')]
            for raw, reason in cases:
                with self.subTest(reason=reason):
                    path.write_text(raw)
                    proc = subprocess.run([sys.executable, str(Path(audit.__file__)), str(path)], capture_output=True, text=True)
                    self.assertEqual(proc.returncode, 1)
                    result = json.loads(proc.stdout)
                    self.assertEqual(result['status'], 'rejected')
                    self.assertIn(reason, result['reason'])


if __name__ == '__main__':
    unittest.main()
