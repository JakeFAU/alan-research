"""Positive semantic controls and exact hostile export mutations.

mutations() is also used by capture.py to retain and execute each CLI input.
"""
import copy
import json
import pathlib
import unittest
import audit

FIXTURE = audit.ROOT / 'docs/foundational/phase3c4/evidence/task1/export.json.gz'


def mutations():
    def change(path, value):
        def edit(data):
            at = data
            for key in path[:-1]:
                at = at[key]
            at[path[-1]] = value
        return edit

    def missing_arc(data):
        # Remove a genuine non-discovery edge, preserving all BFS parents and
        # repairing superficial arc counts. No reachability-only check suffices.
        graph = data['Products'][0]
        seen = {graph['Initial']}
        for node in graph['Nodes']:
            for index, arc in enumerate(node['Arcs']):
                if arc['To'] in seen:
                    del node['Arcs'][index]
                    graph['Arcs'] -= 1
                    data['Counts']['ProductArcs'] -= 1
                    return
                seen.add(arc['To'])
        raise AssertionError('fixture has no non-tree edge')

    def missing_separator(data):
        next(c for c in data['Comparisons'] if c['Mismatch'] and c['Mismatch']['Separators'])['Mismatch']['Separators'].pop()

    def duplicate_edge(data):
        data['Monitors'][0]['Edges'].append(copy.deepcopy(data['Monitors'][0]['Edges'][0]))

    def external_veto(data):
        data['Monitors'][2]['Edges'] = [e for e in data['Monitors'][2]['Edges']
                                      if not (e['From'] == '1' and e['Label'] == 'consume:BC')]

    def missing_coverage(data):
        next(p for p in data['Comparisons'][0]['Pairs'] if p['LeftCoverage'])['LeftCoverage'].pop()

    def missing_match(data):
        comparison = next(c for c in data['Comparisons'] if c['Left'] == 'release-first' and c['Right'] == 'release-first-duplicate')
        next(x for p in comparison['Pairs'] for x in p['LeftCoverage'] if len(x['Matches']) > 1)['Matches'].pop()

    def union(data):
        c = data['Comparisons'][12]
        c.update(Equal=True, Word=[], Mismatch=None)
        for node in c['Pairs']:
            node['Equal'] = True
            for coverage, opposite in [('LeftCoverage', 'Right'), ('RightCoverage', 'Left')]:
                for item in node[coverage]:
                    item['Matches'] = list(node[opposite])
        data['Counts']['EqualProfiles'] += 1

    def drop_language_arc(data):
        data['Comparisons'][0]['Residuals'][0]['Language']['Pairs'][-1]['Arcs'].pop()

    def dead_unreachable(data):
        c = data['Strictness']['Comparison']
        c.update(Equal=True, Word=[], Mismatch=None)
        for p in c['Pairs']:
            if not p['Left'] and p['Right']:
                p.update(Right=[], RightCoverage=[], Equal=True)

    return [
        ('missing-non-tree-arc-repaired-counts', missing_arc, 'Products'),
        ('union-for-profile', union, 'Comparisons'),
        ('missing-coverage', missing_coverage, 'Coverage'),
        ('missing-one-to-many-match', missing_match, 'Matches'),
        ('missing-separator', missing_separator, 'Separators'),
        ('bad-shortest-prefix', change(['Comparisons', 12, 'Word'], ['commit:AB', 'reserve-capacity:AB']), 'Word'),
        ('wrong-separator-side', change(['Comparisons', 12, 'Mismatch', 'Side'], 'left'), 'Side'),
        ('bad-residual-shortest-word', change(['Comparisons', 0, 'Residuals', 0, 'Language', 'Word'], ['apply:BC']), 'Word'),
        ('missing-language-non-tree-arc', drop_language_arc, 'Arcs'),
        ('forged-acceptance', change(['Comparisons', 0, 'InitialLanguage', 'Pairs', 0, 'LeftAccept'], False), 'LeftAccept'),
        ('dead-versus-unreachable', dead_unreachable, 'Strictness'),
        ('forged-mark', change(['Monitors', 6, 'States', 3, 'Marked'], False), 'Marked'),
        ('forged-tau', change(['Monitors', 6, 'Edges', 22, 'Label'], 'apply:AB'), 'Edges'),
        ('external-veto', external_veto, 'Monitors'),
        ('identity-generation', change(['Monitors', 0, 'Signature', 'Base', 'AB', 'CallerGeneration'], 2), 'CallerGeneration'),
        ('root-mismatch', change(['Alphabet', 0, 'Root'], 'other'), 'Root'),
        ('signature-root-limit', change(['Plants', 0, 'Signature', 'Roots', 0, 'Limit'], 2), 'Limit'),
        ('invented-nonblocking', change(['Products', 2, 'Nonblocking'], False), 'Nonblocking'),
        ('invented-node-completion', change(['Products', 2, 'Nodes', 0, 'CanComplete'], False), 'CanComplete'),
        ('forged-product-mark', change(['Products', 2, 'Nodes', 0, 'Marked'], True), 'Marked'),
        ('invalid-lasso-closure', change(['Loop', 'Cycle', 'End'], 'not-the-start'), 'End'),
        ('invalid-lasso-edge', change(['Loop', 'Cycle', 'Steps', 1, 'Label'], 'apply:AB'), 'Label'),
        ('invalid-lasso-debt', change(['Loop', 'Prefix', 'End'], '["W,U,false,false","1"]'), 'End'),
        ('incomplete-escape', change(['Loop', 'Escapes', 0, 'Steps'], []), 'Steps'),
        ('discontinuous-negative-path', change(['Negative', 'RightPath', 'Steps', 1, 'From'], 'elsewhere'), 'From'),
        ('wrong-substitution-outcome', change(['Substitutions', 0, 'NonblockingLeft'], False), 'NonblockingLeft'),
        ('missing-false-context', lambda d: d['Substitutions'].pop(next(i for i,s in enumerate(d['Substitutions']) if not s['NonblockingLeft'])), 'Substitutions'),
        ('source-hash', change(['Sources', 0, 'SHA256'], '0' * 64), 'SHA256'),
        ('source-omission', lambda d: d['Sources'].pop(), 'Sources'),
        ('pinned-raw-hash', change(['Input', 'RawSHA256'], '0' * 64), 'RawSHA256'),
        ('unknown-top-key', lambda d: d.update(Unknown=True), 'keys'),
        ('unknown-nested-key', lambda d: d['Products'][0]['Nodes'][0].update(Unknown=True), 'keys'),
        ('bool-as-int', change(['Counts', 'Plants'], True), 'type'),
        ('int-as-bool', change(['Comparisons', 0, 'Equal'], 0), 'type'),
        ('float-as-int', change(['Counts', 'Plants'], 72.0), 'type'),
        ('null-array', change(['Loop', 'Escapes'], None), 'type'),
        ('duplicate-edge', duplicate_edge, 'Edges'),
        ('duplicate-product', lambda d: d['Products'].append(copy.deepcopy(d['Products'][0])), 'Products'),
        ('duplicate-coverage', lambda d: d['Comparisons'][0]['Pairs'][0]['LeftCoverage'].append(copy.deepcopy(d['Comparisons'][0]['Pairs'][0]['LeftCoverage'][0])), 'Coverage'),
    ]


class IndependentAuditTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.fixture = audit.load(FIXTURE)
        cls.expected = audit.reconstruct()

    def test_complete_positive(self):
        self.assertEqual(audit.audit(self.fixture, self.expected)['Products'], 2772)

    def test_hostile_semantics(self):
        for name, mutate, diagnostic in mutations():
            with self.subTest(name=name):
                forged = copy.deepcopy(self.fixture)
                mutate(forged)
                with self.assertRaisesRegex(audit.InvalidEvidence, diagnostic):
                    audit.audit(forged, self.expected)

    def test_json_controls(self):
        for raw in (b'{"a":1,"a":2}', b'{"a":{"x":1,"x":2}}', b'NaN', b'Infinity', b'-Infinity', b'{} trailing', b'\xff'):
            with self.subTest(raw=raw), self.assertRaises(audit.InvalidEvidence):
                audit.decode(raw)

    def test_union_does_not_preserve_nonblocking(self):
        c = self.expected['Comparisons'][12]
        self.assertTrue(c['InitialLanguage']['Equal'])
        self.assertFalse(c['Equal'])
        self.assertEqual(c['Word'], ['reserve-capacity:AB'])
        self.assertEqual(self.expected['Negative']['LeftPath']['Steps'][0]['Label'], 'commit:AB')
        self.assertTrue(self.expected['CoreOutcomes'][2]['Nonblocking'])
        self.assertFalse(self.expected['CoreOutcomes'][4]['Nonblocking'])

    def test_one_to_many_coverage(self):
        c = next(c for c in self.expected['Comparisons'] if c['Left'] == 'release-first' and c['Right'] == 'release-first-duplicate')
        self.assertTrue(c['Equal'])
        self.assertTrue(any(len(x['Matches']) == 2 for p in c['Pairs'] for x in p['LeftCoverage']))

    def test_empty_profile_is_not_empty_language(self):
        c = self.expected['Strictness']['Comparison']
        self.assertTrue(c['InitialLanguage']['Equal'])
        self.assertFalse(c['Equal'])
        self.assertEqual(c['Word'], ['apply:AB'])
        self.assertEqual(c['Mismatch'], dict(Side='right', State='1', Separators=[]))
        # No marked state exists in either factor, so every pointed product is
        # blocking; profile mismatch therefore need not separate nonblocking.
        for side in ('Left', 'Right'):
            self.assertFalse(any(s['Marked'] for s in self.expected['Strictness'][side]['States']))

    def test_private_cycle_and_completion_escapes(self):
        loop = self.expected['Loop']
        self.assertEqual(loop['Prefix']['End'], loop['Cycle']['Start'])
        self.assertEqual(loop['Cycle']['Start'], loop['Cycle']['End'])
        self.assertEqual([s['Label'] for s in loop['Cycle']['Steps']], ['tau', 'tau'])
        for escape in loop['Escapes']:
            self.assertIn('D,D,false,false', escape['End'])
            self.assertEqual(escape['Steps'][0]['Label'], 'apply:AB')

    def test_each_opposite_residual_needs_own_separator(self):
        # After c, left accepts {a,b}, right has residuals {a} and {b}.
        # No word separates left from their union; two distinct witnesses do.
        a, b, c = 'apply:AB', 'apply:BC', 'commit:BC'
        atom, edge, descriptor = audit.prior.atom, audit.prior.edge, audit.prior.descriptor
        left = descriptor('L', 'generic', {}, '0', [atom(i, i == 2) for i in range(3)],
                          [edge(0,c,1), edge(1,a,2), edge(1,b,2)])
        right = descriptor('R', 'generic', {}, '0', [atom(i, i == 3) for i in range(4)],
                           [edge(0,c,1), edge(0,c,2), edge(1,a,3), edge(2,b,3)])
        decision = audit.profile(left, right)
        self.assertTrue(decision['InitialLanguage']['Equal'])
        self.assertEqual(decision['Word'], [c])
        self.assertEqual(decision['Mismatch']['State'], '1')
        self.assertEqual([s['Word'] for s in decision['Mismatch']['Separators']], [[b], [a]])

    def test_leading_trailing_tau_and_empty_acceptance(self):
        atom, edge, descriptor = audit.prior.atom, audit.prior.edge, audit.prior.descriptor
        left = descriptor('L', 'generic', {}, '0', [atom(i, i == 2) for i in range(3)],
                          [edge(0,'tau',1), edge(1,'tau',0), edge(1,'tau',2)])
        right = descriptor('R', 'generic', {}, '0', [atom(0, True)], [])
        decision = audit.profile(left, right)
        self.assertTrue(decision['Equal'])
        self.assertEqual(decision['InitialLanguage']['Pairs'][0]['Left'], ['0','1','2'])
        self.assertTrue(decision['InitialLanguage']['Pairs'][0]['LeftAccept'])
        self.assertFalse(decision['InitialLanguage']['LeftAccept'])  # equality sentinel


if __name__ == '__main__':
    unittest.main()
