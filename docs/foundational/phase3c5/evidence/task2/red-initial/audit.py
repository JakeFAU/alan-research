#!/usr/bin/env python3
"""Independent Phase 3C.5 reconstruction; no Go predicates or old verdicts.

Unchanged, hash-pinned Python utilities provide descriptors, epsilon-NFA
operations, synchronized products, shortest paths and strict JSON equality.
All inclusion, equivalence reduction, coverage and finite-hit logic is new.
"""
import copy
import gzip
import hashlib
import importlib.util
import itertools
import json
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parents[4]
DEPENDENCIES = {
    'docs/foundational/phase3c4/tools/audit.py': '163fb14ea65af59a507e52b88c188c0b4fe2df65b40e8d2d5c7aef7f076fc88e',
    'docs/foundational/phase3c3/tools/audit.py': '287be61a4e632b37510d401a658be7fa3fea465697d5586b43f07a58b2b7ad75',
    'docs/foundational/phase3c2/tools/audit.py': '1dc485fb2e3f0594b6122daa3a63ec05a12c8f5a6b9988afe5d2a775a768f193',
}


def digest(data):
    return hashlib.sha256(data).hexdigest()


# Check every executable dependency before importing (including transitive ones).
for dependency, sha in DEPENDENCIES.items():
    if digest((ROOT / dependency).read_bytes()) != sha:
        raise RuntimeError('checker dependency hash mismatch: ' + dependency)
sys.dont_write_bytecode = True
_spec = importlib.util.spec_from_file_location('phase3c4_utilities', ROOT / next(iter(DEPENDENCIES)))
old = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(old)
prior = old.prior
InvalidEvidence = old.InvalidEvidence
require, exact, load, decode = old.require, old.exact, old.load, old.decode
Automaton, subsets, product, shortest = old.Automaton, old.subsets, old.product, old.shortest
atom, edge, descriptor = prior.atom, prior.edge, prior.descriptor
ALPHABET = old.ALPHABET
NAMES = old.NAMES + ('minimal-base', 'minimal-augmented')
PREDECESSOR = 'docs/foundational/phase3c4/evidence/task1/export.json.gz'


def monitor(name, signature=prior.CORE):
    if name in old.NAMES:
        return old.monitor(name, signature)
    require(name in NAMES[-2:], 'unknown monitor')
    r, a, c = 'reserve-capacity:AB', 'apply:AB', 'commit:BC'
    count = 5 if name == 'minimal-augmented' else 4
    edges = [edge(s, label, s) for s in range(count) for label in ALPHABET if label not in (r, a, c)]
    edges.append(edge(0, 'tau', 1))
    for label in (r, a, c):
        edges.extend([edge(0, label, 3), edge(3, label, 3)])
    edges.extend([edge(1, r, 2), edge(2, a, 3)])
    if count == 5:
        edges.extend([edge(0, 'tau', 4), edge(4, r, 2), edge(4, c, 3)])
    return descriptor(name, 'M', signature, '0', [atom(s, True) for s in range(count)], edges)


def language(left, right, ls, rs):
    pairs, words, arcs = subsets(left, right, ls, rs)
    nodes = [dict(Left=list(l), Right=list(r), LeftAccept=left.accepts(l),
                  RightAccept=right.accepts(r), Arcs=row) for (l, r), row in zip(pairs, arcs)]
    lo = next((i for i, n in enumerate(nodes) if n['LeftAccept'] and not n['RightAccept']), None)
    ro = next((i for i, n in enumerate(nodes) if n['RightAccept'] and not n['LeftAccept']), None)
    bad = [i for i in (lo, ro) if i is not None]
    first = min(bad) if bad else None
    return dict(Left=left.descriptor['ID'], Right=right.descriptor['ID'], Equal=not bad,
                Word=[] if first is None else words[first],
                LeftAccept=False if first is None else nodes[first]['LeftAccept'],
                RightAccept=False if first is None else nodes[first]['RightAccept'], Pairs=nodes,
                LeftSubset=lo is None, RightSubset=ro is None,
                LeftOutside=None if lo is None else words[lo], RightOutside=None if ro is None else words[ro])


def residuals(left, right):
    rows, lookup = [], {}
    for ls in left.descriptor['States']:
        for rs in right.descriptor['States']:
            l, r = ls['ID'], rs['ID']
            result = language(left, right, l, r)
            rows.append(dict(LeftState=l, RightState=r, Language=result))
            lookup[l, r] = result
    return rows, lookup


def family(states, within):
    groups = []
    for state in sorted(states):
        equal = next((g for g in groups if within[state, g['Representative']]['Equal']), None)
        if equal is None:
            groups.append(dict(Representative=state, Members=[state], ProperSubsets=[], ProperSupersets=[]))
        else:
            equal['Members'].append(state)
    for group in groups:
        rep = group['Representative']
        for other in groups:
            o = other['Representative']
            if o != rep:
                if within[o, rep]['LeftSubset']:
                    group['ProperSubsets'].append(o)
                if within[rep, o]['LeftSubset']:
                    group['ProperSupersets'].append(o)
    return dict(Classes=groups, MIN=[g['Representative'] for g in groups if not g['ProperSubsets']],
                MAX=[g['Representative'] for g in groups if not g['ProperSupersets']])


def accepts_word(automaton, root, word):
    states = automaton.closure([root])
    for label in word:
        states = automaton.advance(states, label)
    return automaton.accepts(states)


def coverage(source, target, sf, tf, lookup, reverse=False):
    out = []
    for target_class in tf['Classes']:
        t = target_class['Representative']
        reps = [g['Representative'] for g in sf['Classes']]
        decisions = {s: lookup[t, s] if reverse else lookup[s, t] for s in reps}
        included = [s for s in reps if decisions[s]['RightSubset' if reverse else 'LeftSubset']]
        failure = None
        if not included:
            outside = [dict(Source=s, Word=decisions[s]['RightOutside' if reverse else 'LeftOutside']) for s in reps]
            h = []
            for item in outside:
                require(item['Word'] is not None, 'internal uncovered source has no outside word')
                if item['Word'] not in h:
                    h.append(item['Word'])
            hits = [dict(Source=s, Hit=any(accepts_word(source, s, w) for w in h)) for s in reps]
            miss = not any(accepts_word(target, t, w) for w in h)
            require(all(row['Hit'] for row in hits) and miss, 'finite hit witness failed')
            failure = dict(Outside=outside, H=h, Hits=hits, TargetMiss=miss)
        out.append(dict(Target=t, Sources=included, Failure=failure))
    return out


def profile(ld, rd):
    left, right = Automaton(ld), Automaton(rd)
    lw, lt = residuals(left, left)
    rw, rt = residuals(right, right)
    cross, table = residuals(left, right)
    pairs, words, arcs = subsets(left, right, ld['Initial'], rd['Initial'])
    nodes = []
    first = {kind: None for kind in ('FULL', 'MIN', 'MAX')}
    for index, ((ls, rs), row) in enumerate(zip(pairs, arcs)):
        lf, rf = family(ls, lt), family(rs, rt)
        flags = {}
        for kind in first:
            lset = [g['Representative'] for g in lf['Classes']] if kind == 'FULL' else lf[kind]
            rset = [g['Representative'] for g in rf['Classes']] if kind == 'FULL' else rf[kind]
            flags[kind] = (all(any(table[l, r]['Equal'] for r in rset) for l in lset)
                           and all(any(table[l, r]['Equal'] for l in lset) for r in rset))
            if not flags[kind] and first[kind] is None:
                first[kind] = words[index]
        nodes.append(dict(Left=list(ls), Right=list(rs), LeftFamily=lf, RightFamily=rf,
                          LeftToRight=coverage(left, right, lf, rf, table),
                          RightToLeft=coverage(right, left, rf, lf, table, True), Arcs=row, **flags))
    initial = language(left, right, ld['Initial'], rd['Initial'])
    return dict(Left=ld['ID'], Right=rd['ID'], InitialLanguage=initial, LeftWithin=lw, RightWithin=rw,
                Residuals=cross, Pairs=nodes, **{k: v is None for k, v in first.items()},
                J=initial['Equal'] and first['MIN'] is None, **{k + 'Word': v for k, v in first.items()})


def rooted(d, root):
    result = copy.deepcopy(d)
    result.update(ID=d['ID'] + '@' + root, Initial=root)
    return result


def trie(name, words):
    states, edges, prefixes = [atom(0, False)], [], {(): '0'}
    for word in words:
        prefix = ()
        for label in word:
            target = prefix + (label,)
            if target not in prefixes:
                prefixes[target] = str(len(states))
                states.append(atom(len(states), False))
                edges.append(edge(prefixes[prefix], label, prefixes[target]))
            prefix = target
        states[int(prefixes[prefix])]['Marked'] = True
    return descriptor(name, 'generic', {}, '0', states, edges)


def generics():
    a, r, c = 'apply:AB', 'reserve-capacity:AB', 'commit:BC'
    x = descriptor('generic-MIN-X', 'generic', {}, 's', [atom('s', True), atom('d', False)],
                   [edge('s', 'tau', 'd')] + [edge(s, l, s) for l in ALPHABET for s in ('s', 'd')])
    y = descriptor('generic-MIN-Y', 'generic', {}, 'd', [atom('d', False)], [edge('d', l, 'd') for l in ALPHABET])
    epsilon = trie('epsilon-trie', [[]])
    def pair(l, r):
        return dict(Left=l, Right=r, Comparison=profile(l, r), Trie=epsilon,
                    LeftProduct=product(l, epsilon), RightProduct=product(r, epsilon))
    loss = pair(x, y)
    x = descriptor('generic-unmarked-X', 'generic', {}, '0', [atom(0, False)], [])
    y = descriptor('generic-unmarked-Y', 'generic', {}, '0', [atom(0, False), atom(1, False)], [edge(0, a, 1)])
    unnecessary = pair(x, y)
    d = descriptor('generic-three-singletons', 'generic', {}, '0', [atom(i, i == 3) for i in range(4)],
                   [edge(0, a, 3), edge(1, r, 3), edge(2, c, 3)])
    automaton = Automaton(d)
    _, table = residuals(automaton, automaton)
    cov = coverage(automaton, automaton, family(['0', '1'], table), family(['2'], table), table)
    h = cov[0]['Failure']['H']
    t = trie('multiword-H', h)
    multi = dict(Automaton=d, SourceStates=['0', '1'], TargetState='2', Coverage=cov, Trie=t,
                 Products=[product(rooted(d, s), t) for s in ('0', '1', '2')])
    x = descriptor('generic-local-hit', 'generic', {}, 'q', [atom('q', False), atom('f', True), atom('d', False)],
                   [edge('q', a, 'f'), edge('q', a, 'd')])
    t = trie('local-hit-H', [[a]])
    local = dict(X=x, H=[[a]], Trie=t, Product=product(x, t))
    return dict(MINLoss=loss, JNonnecessity=unnecessary, Multiword=multi, LocalHitNotNB=local)


def pinned_input():
    compressed = (ROOT / PREDECESSOR).read_bytes()
    exact(len(compressed), 157341, 'pinned compressed length')
    exact(digest(compressed), '35115c3104cddd8a789c80a8c0e73eb7c7f173632af512776ba8502b61b9afc9', 'pinned compressed hash')
    raw = gzip.decompress(compressed)
    exact(len(raw), 8271854, 'pinned raw length')
    exact(digest(raw), 'b40f1d5ac0aea0f43999502ea318cb128f545eb8fb1fd129451915a5aa84bf1f', 'pinned raw hash')
    data = decode(raw)
    exact(data['Schema'], 'alan/phase3c4/v1', 'pinned schema')
    exact(len(data['Sources']), 40, 'inherited source count')
    # Roster itself is bound by the pinned archive; current bytes must still agree.
    for source in data['Sources']:
        exact(digest((ROOT / source['Path']).read_bytes()), source['SHA256'], 'inherited source ' + source['Path'])
    return dict(Path=PREDECESSOR, CompressedBytes=len(compressed), CompressedSHA256=digest(compressed),
                RawBytes=len(raw), RawSHA256=digest(raw), Schema=data['Schema'], Sources=data['Sources']), data


def sources(input_record):
    uses = {s['Path']: 'hash-checked predecessor provenance' for s in input_record['Sources']}
    uses.update({p: 'module' for p in ('go.mod', 'go.sum')})
    uses[PREDECESSOR] = 'decoded pinned automata'
    for name in ('00-contract.md', '01-argument.md', '02-evidence-schema.md'):
        uses['docs/foundational/phase3c5/' + name] = 'hash-checked specification'
    for path in (ROOT / 'research/residualantichains').rglob('*.go'):
        if not path.name.endswith('_test.go'):
            uses[str(path.relative_to(ROOT))] = 'compiled producer'
    return [dict(Path=p, SHA256=digest((ROOT / p).read_bytes()), Use=use) for p, use in sorted(uses.items())]


def reconstruct():
    input_record, pinned = pinned_input()
    plants = [prior.plant(prior.old.config(layout, policy, mask)) for layout in prior.old.LAYOUTS
              for policy in prior.old.POLICIES for mask in range(8)]
    monitors = [monitor(n) for n in NAMES]
    controllers = [prior.controller(code, marks) for code in range(81) for marks in range(4)]
    alphabet = [dict(ID=k + ':' + n, Action=prior.old.action(k, i),
                     Actor='ABC'[i + (k in ('reserve-capacity', 'apply'))], Root='Q')
                for i, n in enumerate(('AB', 'BC')) for k in prior.KINDS]
    for field, expected in [('Plants', plants), ('Monitors', monitors[:7]), ('Controllers', controllers)]:
        exact([prior.normalized_descriptor(d) for d in pinned[field]],
              [prior.normalized_descriptor(d) for d in expected], 'pinned.' + field)
    exact(pinned['Alphabet'], alphabet, 'pinned.Alphabet')
    plants, controllers = pinned['Plants'], pinned['Controllers']
    monitors = pinned['Monitors'] + monitors[7:]
    products = [product(p, monitor(n, p['Signature'])) for p in plants for n in NAMES]
    core = next(p for p in plants if p['ID'] == prior.plant_id(prior.CORE))
    products += [product(core, m, d) for d in controllers for m in monitors]
    comparisons = [profile(l, r) for l, r in itertools.combinations(monitors, 2)]
    lookup = {p['ID']: p for p in products}
    contexts = [(p['ID'], p['ID'], '') for p in plants]
    contexts += [(core['ID'] + ' || ' + d['ID'], core['ID'], ' || ' + d['ID']) for d in controllers]
    substitutions = []
    for comp in comparisons:
        if comp['J']:
            for context, plant_id, suffix in contexts:
                lp, rp = (lookup[plant_id + ' || ' + comp[side] + suffix] for side in ('Left', 'Right'))
                lb, rb = lp['Nonblocking'], rp['Nonblocking']
                li, ri = lp['Nodes'][0]['CanComplete'], rp['Nodes'][0]['CanComplete']
                require(lb == rb and li == ri, 'independent J substitution falsified')
                substitutions.append(dict(Left=comp['Left'], Right=comp['Right'], Context=context,
                                          ProductLeft=lp['ID'], ProductRight=rp['ID'], NonblockingLeft=lb,
                                          NonblockingRight=rb, InitialLeft=li, InitialRight=ri))
    outcomes = [dict(Monitor=m['ID'], Product=(p := lookup[core['ID'] + ' || ' + m['ID']])['ID'],
                     InitialCanComplete=p['Nodes'][0]['CanComplete'], Nonblocking=p['Nonblocking']) for m in monitors]
    left, right = (lookup[core['ID'] + ' || ' + name] for name in ('release-first', 'release-first-doomed'))
    negative = dict(Comparison=next(i for i, c in enumerate(comparisons) if c['Left'] == 'release-first'
                                   and c['Right'] == 'release-first-doomed'),
                    LeftProduct=left['ID'], RightProduct=right['ID'],
                    LeftPath=shortest(left, left['Initial'], lambda n: n['Tuple'] == ['W,U,true,false', '1']),
                    RightPath=shortest(right, right['Initial'], lambda n: n['Tuple'] == ['W,U,true,false', '3']))
    control = copy.deepcopy(core)
    control['ID'] += '/all-marked-control'
    for state in control['States']:
        state['Marked'] = True
    roots = ['1', '4', '0']
    legal_products = [product(control, rooted(monitors[-1], root)) for root in roots]
    def as_automaton(graph):
        return Automaton(descriptor(graph['ID'], 'product', control['Signature'], graph['Initial'],
                                    [atom(n['ID'], n['Marked']) for n in graph['Nodes']],
                                    [edge(n['ID'], a['Label'], a['To']) for n in graph['Nodes'] for a in n['Arcs']]))
    automata = [as_automaton(p) for p in legal_products]
    legal = dict(Plant=control, Monitor='minimal-augmented', Roots=roots, Products=legal_products,
                 Inclusions=[language(l, r, l.descriptor['Initial'], r.descriptor['Initial'])
                             for l, r in zip(automata, automata[1:])])
    counts = dict(Plants=len(plants), Monitors=len(monitors), Controllers=len(controllers), Products=len(products),
                  ProductStates=sum(p['States'] for p in products), ProductArcs=sum(p['Arcs'] for p in products),
                  ProductCanComplete=sum(p['CanComplete'] for p in products),
                  InitialCanComplete=sum(p['Nodes'][0]['CanComplete'] for p in products),
                  NonblockingProducts=sum(p['Nonblocking'] for p in products), Comparisons=len(comparisons),
                  EqualLanguages=sum(c['InitialLanguage']['Equal'] for c in comparisons),
                  **{'Equal' + kind: sum(c[kind] for c in comparisons) for kind in ('FULL', 'MIN', 'MAX', 'J')},
                  ResidualComparisons=sum(len(c[key]) for c in comparisons for key in ('LeftWithin', 'RightWithin', 'Residuals')),
                  LanguageSubsetPairs=sum(len(c['InitialLanguage']['Pairs']) + sum(len(r['Language']['Pairs'])
                                          for key in ('LeftWithin', 'RightWithin', 'Residuals') for r in c[key]) for c in comparisons),
                  ProfileSubsetPairs=sum(len(c['Pairs']) for c in comparisons), Substitutions=len(substitutions),
                  NonblockingSubstitutions=sum(s['NonblockingLeft'] for s in substitutions),
                  BlockingSubstitutions=sum(not s['NonblockingLeft'] for s in substitutions))
    return dict(Schema='alan/phase3c5/v1', Input=input_record, Alphabet=alphabet, Monitors=monitors, Plants=plants,
                Controllers=controllers, Comparisons=comparisons, Products=products, Substitutions=substitutions,
                CoreOutcomes=outcomes, Negative=negative, LegalControl=legal, Generics=generics(), Counts=counts,
                Sources=sources(input_record))


def audit(evidence, expected=None):
    expected = reconstruct() if expected is None else expected
    exact(evidence, expected, 'export')
    return expected['Counts']


def main():
    if len(sys.argv) != 2:
        print('usage: audit.py PATH', file=sys.stderr)
        return 2
    try:
        counts = audit(load(sys.argv[1]))
    except (InvalidEvidence, KeyError, TypeError, ValueError, OSError, IndexError, EOFError) as exc:
        print('REJECT: ' + str(exc), file=sys.stderr)
        return 1
    print('ACCEPT: ' + json.dumps(counts, sort_keys=True))
    return 0


if __name__ == '__main__':
    sys.exit(main())
