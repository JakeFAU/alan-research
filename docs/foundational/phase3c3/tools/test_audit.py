"""Behavioral controls: each mutation names the false claim it challenges."""
import copy
import importlib.util
import pathlib
import unittest

HERE = pathlib.Path(__file__).resolve().parent
spec = importlib.util.spec_from_file_location('interface_audit', HERE / 'audit.py')
audit = importlib.util.module_from_spec(spec)
spec.loader.exec_module(audit)
EXPORT = HERE.parent / 'evidence/task1/export.json.gz'


def assign(path, value):
    def change(data):
        current = data
        for component in path[:-1]:
            current = current[component]
        current[path[-1]] = value
    return change


def remove(path, index):
    def change(data):
        current = data
        for component in path:
            current = current[component]
        current.pop(index)
    return change


def non_tree_arc(data):
    """Remove a non-discovery arc; preserve every BFS parent and superficial count."""
    graph = data['Products'][0]
    lookup = {n['ID']: n for n in graph['Nodes']}
    queue, seen, parents = [graph['Initial']], {graph['Initial']}, set()
    for source in queue:
        for arc in lookup[source]['Arcs']:
            if arc['To'] not in seen:
                parents.add((source, arc['Label'], arc['To']))
                seen.add(arc['To'])
                queue.append(arc['To'])
    for node in graph['Nodes']:
        for i, arc in enumerate(node['Arcs']):
            triple = node['ID'], arc['Label'], arc['To']
            if triple not in parents:
                node['Arcs'].pop(i)
                graph['Arcs'] -= 1
                data['Counts']['ProductArcs'] -= 1
                # Parent evidence remains computed from the unmodified relation.
                return dict(Removed=dict(From=triple[0], Label=triple[1], To=triple[2]),
                            PreservedParents=[dict(From=s, Label=l, To=t) for s, l, t in sorted(parents)],
                            States=graph['States'], RemainingArcs=graph['Arcs'])
    raise AssertionError('no non-tree edge in complete graph')


def swap_recovery(data):
    steps = data['Histories'][5]['Steps']
    steps[5], steps[6] = steps[6], steps[5]


def external_veto(data):
    m = data['Monitors'][2]
    m['Edges'] = [e for e in m['Edges'] if not (e['From'] == '0' and e['Label'] == 'commit:AB')]


def omit_tau(data):
    m = data['Monitors'][3]
    m['Edges'] = [e for e in m['Edges'] if e['Label'] != 'tau']


def extra_field(data):
    data['Histories'][0]['Steps'][0]['Event']['SecretHoldRelease'] = True


def controls():
    # False evidence lives only in the returned copies and retained attack captures.
    return {
        'removed-non-tree-arc': non_tree_arc,
        'changed-plant-gate': remove(['Plants', 0, 'Edges'], 0),
        'forged-equivalence': assign(['Comparisons', 0, 'Equal'], True),
        'wrong-distinguishing-word': assign(['Comparisons', 0, 'Word'], []),
        'nonshortest-distinguishing-word': assign(['Comparisons', 0, 'Word'], ['commit:AB', 'apply:AB']),
        'wrong-subset-acceptance': assign(['Comparisons', 0, 'Pairs', 0, 'LeftAccept'], False),
        'missing-subset-arc': remove(['Comparisons', 0, 'Pairs', 0, 'Arcs'], 0),
        'omitted-tau': omit_tau,
        'forged-monitor-mark': assign(['Monitors', 2, 'States', 0, 'Marked'], False),
        'external-action-veto': external_veto,
        'hidden-effect-label': assign(['Alphabet', 2, 'ID'], 'tau'),
        'forged-effect-binding': assign(['Alphabet', 2, 'Action', 'Binding', 'Grant'], 'forged'),
        'forged-actor': assign(['Alphabet', 2, 'Actor'], 'A'),
        'wrong-root': assign(['Plants', 0, 'Signature', 'Memberships', 0, 'Root'], 'Q-copy'),
        'conflicting-root-limit': assign(['Monitors', 0, 'Signature', 'Roots', 0, 'Limit'], 2),
        'signature-missing-binding': remove(['Plants', 0, 'Signature', 'Base', 'Originals'], 0),
        'boolean-generation': assign(['Monitors', 0, 'Signature', 'Base', 'Agents', 0, 'Generation'], True),
        'fabricated-product-completion': assign(['Products', 1, 'Nodes', 0, 'CanComplete'], True),
        'forged-product-mark': assign(['Products', 0, 'Nodes', 0, 'Marked'], True),
        'wrong-grouping-normalization': assign(['Groupings', 0, 'Left', 'Nodes', 0, 'Tuple', 1], '1'),
        'grouped-internal-effect-hidden': assign(['Groupings', 0, 'Right', 'Nodes', 0, 'Arcs', 0, 'Label'], 'tau'),
        'missing-controller': remove(['Controllers'], 0),
        'duplicate-controller': assign(['Controllers', 1, 'ID'], 'D-00-0'),
        'duplicate-edge': lambda d: d['Monitors'][0]['Edges'].append(copy.deepcopy(d['Monitors'][0]['Edges'][0])),
        'duplicate-state': lambda d: d['Monitors'][0]['States'].append(copy.deepcopy(d['Monitors'][0]['States'][0])),
        'unknown-monitor-label': assign(['Monitors', 0, 'Edges', 0, 'Label'], 'apply:CA'),
        'missing-substitution': remove(['Substitutions'], 0),
        'false-substitution-completion': assign(['Substitutions', 0, 'CanCompleteLeft'], False),
        'missing-history': remove(['Histories'], 0),
        'invented-initial-receipt': assign(['Histories', 0, 'Initial', 'Old', 'Base', 'Originals', 0, 'Retained'], 'r1'),
        'wrong-common-cut-time': assign(['Histories', 0, 'Cuts', 0, 'AfterSteps'], 3),
        'wrong-common-cut-holder': assign(['Histories', 0, 'Steps', 3, 'State', 'Old', 'Held'], [False, True]),
        'false-cut-cancomplete': assign(['Histories', 0, 'Cuts', 0, 'CanComplete'], True),
        'false-final-completion': assign(['Histories', 0, 'Complete'], True),
        'forged-concrete-binding': assign(['Histories', 0, 'Steps', 0, 'Event', 'Binding', 'Operation'], 2),
        'forged-result-provenance': assign(['Histories', 1, 'Steps', 6, 'Event', 'Result'], 'r2'),
        'invented-application-poststate': assign(['Histories', 1, 'Steps', 4, 'State', 'Old', 'Base', 'Originals', 0, 'WorldResult'], 'r2'),
        'missing-real-application': remove(['Histories', 1, 'Steps'], 4),
        'private-step-while-down': swap_recovery,
        'private-step-forged-role': assign(['Histories', 5, 'Steps', 6, 'Event', 'Actor'], 'A'),
        'private-step-hidden-q-change': assign(['Histories', 5, 'Steps', 6, 'State', 'Old', 'Held'], [False, False]),
        'recovery-resets-pc': assign(['Histories', 7, 'Steps', 6, 'State', 'Monitor'], '0'),
        'recovery-releases-q': assign(['Histories', 7, 'Steps', 6, 'State', 'Old', 'Held'], [False, False]),
        'wrong-activation': assign(['Histories', 7, 'Steps', 6, 'State', 'Old', 'Base', 'Agents', 1, 'Activation'], 0),
        'physical-event-wrong-actor': assign(['Histories', 5, 'Steps', 4, 'Event', 'Actor'], 'A'),
        'false-private-flag': assign(['Histories', 5, 'Steps', 6, 'Private'], False),
        'boolean-count': assign(['Counts', 'Histories'], True),
        'float-count': assign(['Counts', 'Histories'], 8.0),
        'null-array': assign(['Monitors'], None),
        'unknown-top-field': lambda d: d.update(Unclaimed=True),
        'unknown-event-field': extra_field,
        'source-forged': assign(['Sources', 0, 'SHA256'], '0' * 64),
        'source-omitted': remove(['Sources'], 0),
    }


class AuditTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.evidence = audit.load(EXPORT)

    def test_complete_export(self):
        result = audit.audit(self.evidence)
        self.assertEqual(result['Products'], 1980)
        self.assertEqual(result['Histories'], 8)

    def test_hostile_semantic_controls(self):
        for name, mutate in controls().items():
            with self.subTest(name=name):
                evidence = copy.deepcopy(self.evidence)
                mutate(evidence)
                with self.assertRaises(audit.InvalidEvidence):
                    audit.audit(evidence)

    def test_non_tree_attack_retains_all_discovery_parents_and_counts(self):
        evidence = copy.deepcopy(self.evidence)
        report = non_tree_arc(evidence)
        graph = evidence['Products'][0]
        arcs = {(n['ID'], a['Label'], a['To']) for n in graph['Nodes'] for a in n['Arcs']}
        self.assertEqual(graph['Arcs'], len(arcs))
        self.assertEqual(evidence['Counts']['ProductArcs'], sum(g['Arcs'] for g in evidence['Products']))
        self.assertEqual(len(report['PreservedParents']), graph['States'] - 1)
        self.assertTrue(all((p['From'], p['Label'], p['To']) in arcs for p in report['PreservedParents']))
        self.assertEqual(len(audit.normalized_graph(graph)['Nodes']), graph['States'])

    def test_duplicate_json_key_and_nonfinite_number(self):
        for raw in (b'{"Schema":1,"Schema":2}', b'{"count":NaN}', b'{"count":Infinity}'):
            with self.assertRaises(audit.InvalidEvidence):
                audit.decode(raw)

    def test_equal_language_does_not_preserve_nonblocking(self):
        plants, monitors, _, products, _, comparisons, _ = audit.reconstruct()
        self.assertTrue(audit.comparison(monitors[2], monitors[3])['Equal'])
        self.assertTrue(audit.comparison(monitors[2], monitors[4])['Equal'])
        self.assertTrue(products[2]['Nodes'][0]['CanComplete'])
        self.assertTrue(products[4]['Nodes'][0]['CanComplete'])
        self.assertTrue(products[2]['Nonblocking'])
        self.assertFalse(products[4]['Nonblocking'])
        self.assertTrue(any(not n['CanComplete'] for n in products[4]['Nodes']))

    def test_epsilon_cycles_nondeterminism_empty_subsets_and_empty_words(self):
        # An epsilon cycle with a nondeterministic visible edge: existential
        # acceptance survives an unmarked sink; the empty subset stays empty.
        a = audit.descriptor('a', 'M', audit.CORE, '0',
                             [audit.atom('0', False), audit.atom('1', True), audit.atom('2', False)],
                             [audit.edge('0', 'tau', '1'), audit.edge('1', 'tau', '0'),
                              audit.edge('1', 'apply:AB', '1'), audit.edge('1', 'apply:AB', '2')])
        b = audit.descriptor('b', 'M', audit.CORE, '0', [audit.atom('0', True)],
                             [audit.edge('0', 'apply:AB', '0')])
        self.assertTrue(audit.comparison(a, b)['Equal'])
        self.assertTrue(audit.accepts(a, []))
        self.assertTrue(audit.accepts(a, ['apply:AB'] * 4))
        self.assertFalse(audit.accepts(a, ['commit:AB']))
        self.assertEqual(audit.advance(a, (), 'apply:AB'), ())
        b['States'][0]['Marked'] = False
        comp = audit.comparison(a, b)
        self.assertFalse(comp['Equal'])
        self.assertEqual(comp['Word'], [])
        self.assertTrue(comp['LeftAccept'])
        self.assertFalse(comp['RightAccept'])

    def test_shortest_separators_have_independent_acceptance(self):
        _, monitors, _, _, _, comparisons, _ = audit.reconstruct()
        by_name = {m['ID']: m for m in monitors}
        for comp in comparisons:
            if comp['Equal']:
                continue
            left, right = by_name[comp['Left']], by_name[comp['Right']]
            self.assertNotEqual(audit.accepts(left, comp['Word']), audit.accepts(right, comp['Word']))
            for length in range(len(comp['Word'])):
                for word in audit.itertools.product(audit.ALPHABET, repeat=length):
                    self.assertEqual(audit.accepts(left, word), audit.accepts(right, word))


if __name__ == '__main__':
    unittest.main()
