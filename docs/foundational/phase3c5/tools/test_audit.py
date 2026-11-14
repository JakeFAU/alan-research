"""Independent semantic controls and exact CLI attack generators."""
# This direct-script bootstrap reads the helper once and compiles those bytes.
# Accepted entrypoints then execute their own recorded source buffer as well.
if '__source_bundle__' not in globals():
    import pathlib as _pathlib
    import types as _types
    _root = _pathlib.Path(__file__).resolve().parents[4]
    _helper_path = _root / 'docs/foundational/phase3c5/tools/source_loading.py'
    _helper_bytes = _helper_path.read_bytes()
    _helper = _types.ModuleType('phase3c5_source_loading')
    _helper.__file__ = str(_helper_path)
    exec(compile(_helper_bytes, str(_helper_path), 'exec'), _helper.__dict__)
    __source_bundle__ = _helper.SourceBundle(_root, _helper_bytes)
    if __name__ == '__main__':
        __source_bundle__.load('__main__', str(_pathlib.Path(__file__).resolve().relative_to(_root)))
        raise SystemExit(0)

import copy
import itertools
import unittest
audit = __source_bundle__.load('audit', 'docs/foundational/phase3c5/tools/audit.py')

FIXTURE = audit.ROOT / 'docs/foundational/phase3c5/evidence/task1/capture-3/export.json.gz'


def mutations():
    def change(path, value):
        def edit(data):
            at = data
            for key in path[:-1]:
                at = at[key]
            at[path[-1]] = value
        return edit

    def first_failure(data):
        return next(c for p in data['Comparisons'] for n in p['Pairs'] for side in ('LeftToRight', 'RightToLeft')
                    for c in n[side] if c['Failure'] and c['Failure']['H'])['Failure']

    def missing_arc(data):
        graph = data['Products'][0]
        seen = {graph['Initial']}
        for node in graph['Nodes']:
            for i, arc in enumerate(node['Arcs']):
                if arc['To'] in seen:
                    del node['Arcs'][i]
                    graph['Arcs'] -= 1
                    data['Counts']['ProductArcs'] -= 1
                    return
                seen.add(arc['To'])
        raise AssertionError('no non-tree edge')

    def reversed_inclusion(data):
        language = next(r['Language'] for c in data['Comparisons'] for r in c['Residuals']
                        if r['Language']['LeftSubset'] != r['Language']['RightSubset'])
        for l, r in [('LeftSubset', 'RightSubset'), ('LeftOutside', 'RightOutside')]:
            language[l], language[r] = language[r], language[l]

    def dropped_equivalent(data):
        cls = next(g for c in data['Comparisons'] for n in c['Pairs'] for f in ('LeftFamily', 'RightFamily')
                   for g in n[f]['Classes'] if len(g['Members']) > 1 and g['Representative'] in n[f]['MIN'])
        cls['Members'].pop()

    def max_for_min(data):
        fam = next(n[f] for c in data['Comparisons'] for n in c['Pairs'] for f in ('LeftFamily', 'RightFamily')
                   if n[f]['MIN'] != n[f]['MAX'])
        fam['MIN'] = fam['MAX'][:]

    def empty_confusion(data):
        fam = next(n['LeftFamily'] for n in data['Generics']['JNonnecessity']['Comparison']['Pairs'] if not n['Left'])
        fam.update(Classes=[dict(Representative='0', Members=['0'], ProperSubsets=[], ProperSupersets=[])], MIN=['0'], MAX=['0'])

    def empty_h(data):
        failure = next(c['Failure'] for p in data['Comparisons'] for n in p['Pairs'] for side in ('LeftToRight', 'RightToLeft')
                       for c in n[side] if c['Failure'] and not c['Failure']['H'])
        failure['H'] = [[]]

    def bad_prefix(data):
        c = next(c for c in data['Comparisons'] if c['MINWord'] is not None)
        c['MINWord'] = ['consume:BC'] + c['MINWord']

    def bad_word(data):
        r = next(r['Language'] for c in data['Comparisons'] for r in c['Residuals'] if r['Language']['LeftOutside'] is not None)
        r['LeftOutside'] = ['consume:BC'] + r['LeftOutside']

    def bad_tau(data):
        e = next(e for e in data['Monitors'][-1]['Edges'] if e['Label'] == 'tau')
        e['To'] = '3'

    def wrong_dominance(data):
        cls = next(g for c in data['Comparisons'] for n in c['Pairs'] for f in ('LeftFamily', 'RightFamily')
                   for g in n[f]['Classes'] if g['ProperSubsets'])
        cls['ProperSubsets'] = []

    return [
        ('missing-non-tree-arc-repaired-counts', missing_arc, 'Products'),
        ('missing-language-non-tree-arc', lambda d: d['Comparisons'][0]['InitialLanguage']['Pairs'][-1]['Arcs'].pop(), 'Arcs'),
        ('missing-profile-non-tree-arc', lambda d: d['Comparisons'][0]['Pairs'][-1]['Arcs'].pop(), 'Arcs'),
        ('reversed-inclusion', reversed_inclusion, 'Subset'),
        ('dropped-equivalent-minimum-member', dropped_equivalent, 'Members'),
        ('max-substituted-for-min', max_for_min, 'MIN'),
        ('lost-proper-dominator', wrong_dominance, 'ProperSubsets'),
        ('empty-family-as-empty-language', empty_confusion, 'LeftFamily'),
        ('empty-H-as-epsilon-H', empty_h, 'H'),
        ('missing-H-member', lambda d: first_failure(d)['H'].pop(), 'H'),
        ('forged-H-member', lambda d: first_failure(d)['H'].append(['consume:BC']), 'H'),
        ('missing-outside-source', lambda d: first_failure(d)['Outside'].pop(), 'Outside'),
        ('false-source-hit', lambda d: first_failure(d)['Hits'][0].update(Hit=False), 'Hit'),
        ('false-target-miss', lambda d: first_failure(d).update(TargetMiss=False), 'TargetMiss'),
        ('bad-shortest-prefix', bad_prefix, 'MINWord'),
        ('bad-shortest-directional-word', bad_word, 'Outside'),
        ('false-initial-acceptance', change(['Comparisons',0,'InitialLanguage','Pairs',0,'LeftAccept'], False), 'LeftAccept'),
        ('wrong-generic-MIN-initial', change(['Generics','MINLoss','LeftProduct','Nodes',0,'CanComplete'],False), 'CanComplete'),
        ('wrong-generic-nonnecessity', change(['Generics','JNonnecessity','Comparison','J'],True), 'J'),
        ('wrong-local-hit-NB', change(['Generics','LocalHitNotNB','Product','Nonblocking'],True), 'Nonblocking'),
        ('wrong-multiword-target-hit', change(['Generics','Multiword','Products',2,'Nodes',0,'CanComplete'],True), 'CanComplete'),
        ('forged-trie-mark', change(['Generics','Multiword','Trie','States',0,'Marked'],True), 'Marked'),
        ('forged-monitor-mark', change(['Monitors',8,'States',4,'Marked'],False), 'Marked'),
        ('forged-tau', bad_tau, 'Edges'),
        ('forged-product-mark', change(['Products',0,'Nodes',0,'Marked'],True), 'Marked'),
        ('forged-product-completion', change(['Products',0,'Nodes',0,'CanComplete'],False), 'CanComplete'),
        ('forged-product-node', change(['Products',0,'Nodes',0,'ID'],'fiction'), 'ID'),
        ('bad-negative-path', change(['Negative','RightPath','Steps',1,'From'],'fiction'), 'From'),
        ('bad-legal-trace-inclusion', change(['LegalControl','Inclusions',0,'RightSubset'],True), 'RightSubset'),
        ('forged-legal-plant-mark', change(['LegalControl','Plant','States',0,'Marked'],False), 'Marked'),
        ('missing-blocking-substitution', lambda d: d['Substitutions'].pop(next(i for i,s in enumerate(d['Substitutions']) if not s['NonblockingLeft'])), 'Substitutions'),
        ('wrong-J-substitution-initial', change(['Substitutions',0,'InitialLeft'],False), 'InitialLeft'),
        ('signature-generation', change(['Monitors',0,'Signature','Base','Originals',0,'CallerGeneration'],2), 'CallerGeneration'),
        ('signature-root-limit', change(['Plants',0,'Signature','Roots',0,'Limit'],2), 'Limit'),
        ('alphabet-root', change(['Alphabet',0,'Root'],'other'), 'Root'),
        ('source-hash', change(['Sources',0,'SHA256'],'0'*64), 'SHA256'),
        ('source-use', change(['Sources',0,'Use'],'compiled producer'), 'Use'),
        ('source-omission', lambda d: d['Sources'].pop(), 'Sources'),
        ('inherited-source-hash', change(['Input','Sources',0,'SHA256'],'0'*64), 'SHA256'),
        ('pinned-raw-hash', change(['Input','RawSHA256'],'0'*64), 'RawSHA256'),
        ('unknown-key', lambda d: d.update(Unknown=True), 'keys'),
        ('unknown-nested-key', lambda d: d['Products'][0]['Nodes'][0].update(Unknown=True), 'keys'),
        ('bool-as-int', change(['Counts','Plants'],True), 'type'),
        ('int-as-bool', change(['Comparisons',0,'J'],0), 'type'),
        ('float-as-int', change(['Counts','Plants'],72.0), 'type'),
        ('null-array', change(['Generics','Multiword','Coverage'],None), 'type'),
        ('duplicate-product', lambda d: d['Products'].append(copy.deepcopy(d['Products'][0])), 'Products'),
    ]


def json_attacks():
    return [('duplicate-top-key', b'{"Schema":1,"Schema":2}', 'duplicate key'),
            ('duplicate-nested-key', b'{"x":{"n":1,"n":2}}', 'duplicate key'),
            ('nan', b'NaN', 'nonfinite'), ('infinity', b'Infinity', 'nonfinite'),
            ('negative-infinity', b'-Infinity', 'nonfinite'),
            ('overflow-number', b'{"Counts":1e999}', 'keys'),
            ('trailing-json', b'{} {}', 'Extra data'), ('invalid-utf8', b'\xff', 'codec'),
            ('malformed-gzip', b'\x1f\x8bBAD', 'Compressed file ended')]


class IndependentAuditTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.fixture = audit.load(FIXTURE)
        cls.expected = audit.reconstruct()

    def test_complete_positive(self):
        self.assertEqual(audit.audit(self.fixture, self.expected)['Products'],3564)

    def test_hostile_semantics(self):
        for name, mutate, diagnostic in mutations():
            with self.subTest(name=name):
                forged = copy.deepcopy(self.fixture)
                mutate(forged)
                with self.assertRaisesRegex(audit.InvalidEvidence, diagnostic):
                    audit.audit(forged, self.expected)

    def test_strict_decoder(self):
        for raw in (b'{"a":1,"a":2}', b'{"a":{"x":1,"x":2}}', b'NaN', b'Infinity', b'-Infinity', b'{} trailing', b'\xff'):
            with self.subTest(raw=raw), self.assertRaises(audit.InvalidEvidence):
                audit.decode(raw)

    def test_min_loss_and_nonnecessity(self):
        generic = self.expected['Generics']
        self.assertTrue(generic['MINLoss']['Comparison']['MIN'])
        self.assertFalse(generic['MINLoss']['Comparison']['J'])
        self.assertTrue(generic['MINLoss']['LeftProduct']['Nodes'][0]['CanComplete'])
        self.assertFalse(generic['MINLoss']['RightProduct']['Nodes'][0]['CanComplete'])
        self.assertEqual(generic['JNonnecessity']['Comparison']['MINWord'], ['apply:AB'])
        for side in ('Left', 'Right'):
            self.assertFalse(any(s['Marked'] for s in generic['JNonnecessity'][side]['States']))

    def test_local_hit_is_not_nonblocking(self):
        g = self.expected['Generics']['LocalHitNotNB']['Product']
        self.assertTrue(g['Nodes'][0]['CanComplete'])
        self.assertFalse(g['Nonblocking'])
        multi = self.expected['Generics']['Multiword']
        self.assertEqual(multi['Coverage'][0]['Failure']['H'], [['apply:AB'], ['reserve-capacity:AB']])
        self.assertEqual([p['Nodes'][0]['CanComplete'] for p in multi['Products']], [True,True,False])

    def test_new_pair_and_blocked_contexts(self):
        comp = self.expected['Comparisons'][-1]
        self.assertEqual((comp['FULL'],comp['MIN'],comp['J']),(False,True,True))
        self.assertEqual(comp['FULLWord'],[])
        for name in audit.NAMES[-2:]:
            plants = [p for p in self.expected['Products'] if p['Factors'][1] == name and len(p['Factors']) == 2]
            controllers = [p for p in self.expected['Products'] if p['Factors'][1] == name and len(p['Factors']) == 3]
            self.assertEqual(sum(p['Nonblocking'] for p in plants),52)
            self.assertEqual(sum(not p['Nonblocking'] for p in plants),20)
            self.assertEqual(sum(not p['Nonblocking'] for p in controllers),324)

    def test_max_negative_and_real_paths(self):
        n = self.expected['Negative']
        c = self.expected['Comparisons'][n['Comparison']]
        self.assertTrue(c['MAX'] and c['InitialLanguage']['Equal'])
        self.assertFalse(c['MIN'])
        products = {p['ID']:p for p in self.expected['Products']}
        for side, expected in [('Left',True),('Right',False)]:
            p = products[n[side+'Product']]
            self.assertEqual(p['Nonblocking'],expected)
            nodes = {x['ID']:x for x in p['Nodes']}
            for s in n[side+'Path']['Steps']:
                self.assertIn(dict(Label=s['Label'],To=s['To']),nodes[s['From']]['Arcs'])
            self.assertEqual(nodes[n[side+'Path']['End']]['CanComplete'], expected)

    def test_legal_trace_strictness(self):
        for c in self.expected['LegalControl']['Inclusions']:
            self.assertTrue(c['LeftSubset'])
            self.assertFalse(c['RightSubset'])
            self.assertIsNotNone(c['RightOutside'])

    def test_empty_language_empty_family_epsilon_and_tau_cycle(self):
        x = audit.descriptor('x','generic',{},'0',[audit.atom(0,False),audit.atom(1,True)],
                             [audit.edge(0,'tau',1),audit.edge(1,'tau',0)])
        y = audit.descriptor('y','generic',{},'0',[audit.atom(0,False)],[])
        la = audit.language(audit.Automaton(x),audit.Automaton(y),'0','0')
        self.assertFalse(la['LeftSubset'])
        self.assertTrue(la['RightSubset'])
        self.assertEqual(la['LeftOutside'],[])
        self.assertIsNone(la['RightOutside'])
        c = audit.profile(x,x)
        self.assertEqual(c['Pairs'][0]['LeftFamily']['MIN'],['0'])
        self.assertEqual(c['Pairs'][0]['LeftFamily']['Classes'][0]['Members'],['0','1'])
        self.assertEqual(c['Pairs'][0]['LeftFamily']['Classes'][0]['ProperSubsets'],[])
        self.assertEqual(c['Pairs'][1]['LeftFamily'],dict(Classes=[],MIN=[],MAX=[]))
        for words, marked in [([],False),([[]],True)]:
            t = audit.trie('t',words)
            self.assertEqual(t['States'][0]['Marked'],marked)
            self.assertEqual(t['Edges'],[])

    def test_finite_algebra_exhaustive_three_word_universe(self):
        # Set-theoretic oracle: all 256 families of subsets of a 3-word universe.
        # Independent of automata, inclusion certificates and producer code.
        universe = frozenset(range(3))
        languages = [frozenset(i for i in universe if mask & (1<<i)) for mask in range(8)]
        families = [frozenset(languages[i] for i in range(8) if mask & (1<<i)) for mask in range(256)]
        def minima(family):
            return frozenset(l for l in family if not any(k < l for k in family))
        mins = {f:minima(f) for f in families}
        hits = {f:tuple(all(l & h for l in f) for h in languages) for f in families}
        for f,g in itertools.product(families,repeat=2):
            self.assertEqual(minima(frozenset(l & k for l in f for k in g)),
                             minima(frozenset(l & k for l in mins[f] for k in mins[g])))
            coverage = all(any(l <= k for l in f) for k in g)
            self.assertEqual(all(not x or y for x,y in zip(hits[f],hits[g])),coverage)
            mutual = coverage and all(any(k <= l for k in g) for l in f)
            self.assertEqual(mins[f] == mins[g],mutual)
        for f in families:
            self.assertEqual(hits[f],hits[mins[f]])
            self.assertTrue(all(any(k <= l for k in mins[f]) for l in f))


if __name__ == '__main__':
    unittest.main()
