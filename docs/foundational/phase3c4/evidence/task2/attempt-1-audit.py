#!/usr/bin/env python3
"""Independent Phase 3C.4 reconstruction from the written finite contract.

Only Phase 3C.3 descriptor constructors and Phase 3C.2 JSON checking utilities
are reused. No predecessor product/verdict or Go operational code is executed.
"""
import copy
import gzip
import hashlib
import importlib.util
import itertools
import json
import pathlib
import sys
from collections import defaultdict

ROOT = pathlib.Path(__file__).resolve().parents[4]
PREDECESSOR = 'docs/foundational/phase3c3/evidence/task1/export.json.gz'
PYTHON_DEPENDENCIES = ('docs/foundational/phase3c3/tools/audit.py',
                       'docs/foundational/phase3c2/tools/audit.py')
sys.dont_write_bytecode = True
_spec = importlib.util.spec_from_file_location('phase3c3_constructors', ROOT / PYTHON_DEPENDENCIES[0])
prior = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(prior)
InvalidEvidence = prior.InvalidEvidence
require, exact, load, decode = prior.require, prior.exact, prior.load, prior.decode
ALPHABET = prior.ALPHABET
NAMES = prior.MONITORS + ('release-first-duplicate', 'release-first-loop')


def digest(raw):
    return hashlib.sha256(raw).hexdigest()


def encoded_tuple(states):
    return json.dumps(states, separators=(',', ':'))


def monitor(name, signature=prior.CORE):
    if name in prior.MONITORS:
        return prior.monitor(name, signature)
    reserve, apply, commit = 'reserve-capacity:AB', 'apply:AB', 'commit:BC'
    edges = [prior.edge(s, label, s) for s in range(4) for label in ALPHABET
             if label not in (reserve, apply, commit)]
    edges += [prior.edge(0, reserve, 1), prior.edge(1, apply, 2)]
    if name == 'release-first-duplicate':
        edges += [prior.edge(1, apply, 3)]
        terminal = (2, 3)
    else:
        require(name == 'release-first-loop', 'unknown monitor')
        edges += [prior.edge(1, 'tau', 3), prior.edge(3, 'tau', 1), prior.edge(3, apply, 2)]
        terminal = (2,)
    edges += [prior.edge(s, label, s) for s in terminal for label in (reserve, apply, commit)]
    return prior.descriptor(name, 'M', signature, '0', [prior.atom(s, True) for s in range(4)], edges)


class Automaton:
    def __init__(self, descriptor):
        self.descriptor = descriptor
        self.marks = {s['ID'] for s in descriptor['States'] if s['Marked']}
        self.adjacency = defaultdict(set)
        for edge in descriptor['Edges']:
            self.adjacency[edge['From'], edge['Label']].add(edge['To'])

    def closure(self, states):
        reached = set(states)
        pending = list(reached)
        for state in pending:
            for target in self.adjacency[state, 'tau'] - reached:
                reached.add(target)
                pending.append(target)
        return tuple(sorted(reached))

    def advance(self, states, label):
        return self.closure({t for s in states for t in self.adjacency[s, label]})

    def accepts(self, states):
        return bool(self.marks.intersection(states))


def subsets(left, right, left_start, right_start):
    """Complete BFS relation, including empty subsets and post-mismatch arcs."""
    initial = (left.closure([left_start]), right.closure([right_start]))
    pending, indices, words, arcs = [initial], {initial: 0}, [[]], []
    for index, (ls, rs) in enumerate(pending):
        row = []
        for label in ALPHABET:
            target = (left.advance(ls, label), right.advance(rs, label))
            if target not in indices:
                indices[target] = len(pending)
                pending.append(target)
                words.append(words[index] + [label])
            row.append(dict(Label=label, To=indices[target]))
        arcs.append(row)
    return pending, words, arcs


def language(left, right, left_start, right_start):
    pairs, words, arcs = subsets(left, right, left_start, right_start)
    nodes = [dict(Left=list(ls), Right=list(rs), LeftAccept=left.accepts(ls),
                  RightAccept=right.accepts(rs), Arcs=row)
             for (ls, rs), row in zip(pairs, arcs)]
    first = next((i for i, n in enumerate(nodes) if n['LeftAccept'] != n['RightAccept']), None)
    return dict(Left=left.descriptor['ID'], Right=right.descriptor['ID'], Equal=first is None,
                Word=[] if first is None else words[first],
                LeftAccept=False if first is None else nodes[first]['LeftAccept'],
                RightAccept=False if first is None else nodes[first]['RightAccept'], Pairs=nodes)


def profile(ld, rd):
    left, right = Automaton(ld), Automaton(rd)
    residuals, lookup = [], {}
    for ls, rs in itertools.product(ld['States'], rd['States']):
        lookup[ls['ID'], rs['ID']] = len(residuals)
        residuals.append(dict(LeftState=ls['ID'], RightState=rs['ID'],
                              Language=language(left, right, ls['ID'], rs['ID'])))
    pairs, words, arcs = subsets(left, right, ld['Initial'], rd['Initial'])
    nodes = []
    mismatch = None
    first_word = []
    for index, ((ls, rs), row) in enumerate(zip(pairs, arcs)):
        lc = [dict(State=s, Matches=[t for t in rs if residuals[lookup[s, t]]['Language']['Equal']]) for s in ls]
        rc = [dict(State=t, Matches=[s for s in ls if residuals[lookup[s, t]]['Language']['Equal']]) for t in rs]
        equal = all(c['Matches'] for c in lc + rc)
        nodes.append(dict(Left=list(ls), Right=list(rs), LeftCoverage=lc, RightCoverage=rc, Equal=equal, Arcs=row))
        if not equal and mismatch is None:
            unmatched = next(((side, c['State']) for side, coverage in [('left', lc), ('right', rc)]
                              for c in coverage if not c['Matches']))
            side, state = unmatched
            separators = []
            for opposite in rs if side == 'left' else ls:
                residual = lookup[state, opposite] if side == 'left' else lookup[opposite, state]
                decision = residuals[residual]['Language']
                require(not decision['Equal'], 'internal unmatched state has equal partner')
                separators.append(dict(OppositeState=opposite, Residual=residual, Word=decision['Word'],
                                       LeftAccept=decision['LeftAccept'], RightAccept=decision['RightAccept']))
            mismatch = dict(Side=side, State=state, Separators=separators)
            first_word = words[index]
    return dict(Left=ld['ID'], Right=rd['ID'],
                InitialLanguage=language(left, right, ld['Initial'], rd['Initial']),
                Residuals=residuals, Pairs=nodes, Equal=mismatch is None,
                Word=first_word, Mismatch=mismatch)


def product(*descriptors):
    """Fresh synchronized relation and reverse marked reachability."""
    for d in descriptors[1:]:
        exact(d['Signature'], descriptors[0]['Signature'], 'compatible signature')
    factors = [Automaton(d) for d in descriptors]
    initial = tuple(d['Initial'] for d in descriptors)
    pending, seen, nodes = [initial], {initial}, []
    for states in pending:
        transitions = []
        for i, (factor, state) in enumerate(zip(factors, states)):
            transitions += [('tau', states[:i] + (target,) + states[i+1:])
                            for target in sorted(factor.adjacency[state, 'tau'])]
        for label in ALPHABET:
            transitions += [(label, target) for target in itertools.product(
                *(sorted(f.adjacency[s, label]) for f, s in zip(factors, states)))]
        for _, target in transitions:
            if target not in seen:
                seen.add(target)
                pending.append(target)
        relation = sorted({(label, encoded_tuple(target)) for label, target in transitions})
        nodes.append(dict(ID=encoded_tuple(states), Tuple=list(states),
                          Marked=all(s in f.marks for f, s in zip(factors, states)),
                          CanComplete=False, Arcs=[dict(Label=l, To=t) for l, t in relation]))
    reverse = defaultdict(set)
    for node in nodes:
        for arc in node['Arcs']:
            reverse[arc['To']].add(node['ID'])
    can = {n['ID'] for n in nodes if n['Marked']}
    work = list(can)
    for state in work:
        for source in reverse[state] - can:
            can.add(source)
            work.append(source)
    for node in nodes:
        node['CanComplete'] = node['ID'] in can
    return dict(ID=' || '.join(d['ID'] for d in descriptors), Factors=[d['ID'] for d in descriptors],
                Initial=encoded_tuple(initial), Nodes=nodes, States=len(nodes),
                Arcs=sum(len(n['Arcs']) for n in nodes), CanComplete=len(can), Nonblocking=len(can) == len(nodes))


def shortest(graph, start, goal):
    nodes = {n['ID']: n for n in graph['Nodes']}
    pending, paths = [start], {start: []}
    for state in pending:
        if goal(nodes[state]):
            return dict(Start=start, Steps=paths[state], End=state)
        for arc in nodes[state]['Arcs']:
            target = arc['To']
            if target not in paths:
                paths[target] = paths[state] + [dict(From=state, Label=arc['Label'], To=target)]
                pending.append(target)
    raise InvalidEvidence('no independent witness path')


def witnesses(products, comparisons):
    by_id = {p['ID']: p for p in products}
    prefix = prior.plant_id(prior.CORE) + ' || '
    left, right, loop = (by_id[prefix + m] for m in ('release-first', 'release-first-doomed', 'release-first-loop'))
    debt1, debt3 = (encoded_tuple(['W,U,true,false', s]) for s in ('1', '3'))
    negative = dict(Comparison=next(i for i, c in enumerate(comparisons)
                                   if c['Left'] == 'release-first' and c['Right'] == 'release-first-doomed'),
                    LeftProduct=left['ID'], RightProduct=right['ID'],
                    LeftPath=shortest(left, left['Initial'], lambda n: n['ID'] == debt1),
                    RightPath=shortest(right, right['Initial'], lambda n: n['ID'] == debt3))
    require(left['Nonblocking'] and not right['Nonblocking'], 'union falsification missing')
    require(left['Nodes'][0]['CanComplete'] and right['Nodes'][0]['CanComplete'], 'initial completion differs')
    nodes = {n['ID']: n for n in loop['Nodes']}
    steps = []
    for source, target in ((debt1, debt3), (debt3, debt1)):
        require(dict(Label='tau', To=target) in nodes[source]['Arcs'], 'cycle arc missing')
        require(not nodes[source]['Marked'] and nodes[source]['CanComplete'], 'cycle marking/completion')
        steps.append(dict(From=source, Label='tau', To=target))
    lasso = dict(Product=loop['ID'], Prefix=shortest(loop, loop['Initial'], lambda n: n['ID'] == debt1),
                 Cycle=dict(Start=debt1, Steps=steps, End=debt1),
                 Escapes=[shortest(loop, s, lambda n: n['Marked']) for s in (debt1, debt3)])
    x = prior.descriptor('generic-unmarked-X', 'generic', {}, '0', [prior.atom(0, False)], [])
    y = prior.descriptor('generic-unmarked-Y', 'generic', {}, '0', [prior.atom(s, False) for s in range(2)],
                         [prior.edge(0, 'apply:AB', 1)])
    strictness = dict(Left=x, Right=y, Comparison=profile(x, y))
    require(strictness['Comparison']['InitialLanguage']['Equal'] and not strictness['Comparison']['Equal'],
            'generic strictness control failed')
    return negative, lasso, strictness


def pinned_input():
    compressed = (ROOT / PREDECESSOR).read_bytes()
    exact(digest(compressed), '86a9d22f9808a67dd0066b6b77dd2b3c01a2125d10863c6fbe1369097f6624c8', 'pinned gzip')
    raw = gzip.decompress(compressed)
    exact(len(raw), 25384964, 'pinned raw length')
    exact(digest(raw), 'a073a22fe8ac62e05f2d15b0f5146b19db9606e38ce4619261e7e4108e27269d', 'pinned raw')
    data = decode(raw)
    exact(data['Schema'], 'alan/phase3c3/v1', 'pinned schema')
    exact(data['Sources'], prior.source_records(), 'predecessor source roster')
    exact(len(data['Sources']), 27, 'predecessor source count')
    record = dict(Path=PREDECESSOR, CompressedBytes=len(compressed), CompressedSHA256=digest(compressed),
                  RawBytes=len(raw), RawSHA256=digest(raw), Schema=data['Schema'], Sources=data['Sources'])
    return record, data


def sources(input_record):
    uses = {s['Path']: 'hash-checked predecessor provenance' for s in input_record['Sources']}
    uses.update({p: 'module' for p in ('go.mod', 'go.sum')})
    uses[PREDECESSOR] = 'decoded pinned automata'
    for name in ('00-contract.md', '01-argument.md', '02-evidence-schema.md'):
        uses['docs/foundational/phase3c4/' + name] = 'hash-checked specification'
    for path in (ROOT / 'research/residualinterfaces').rglob('*.go'):
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
    # Source-bound constructors reconstruct the pinned descriptor relations.
    # Predecessor descriptor edge order is unspecified, unlike this new export.
    for field, expected in [('Plants', plants), ('Monitors', monitors[:5]), ('Controllers', controllers)]:
        exact([prior.normalized_descriptor(d) for d in pinned[field]],
              [prior.normalized_descriptor(d) for d in expected], 'pinned.' + field)
    exact(pinned['Alphabet'], alphabet, 'pinned.Alphabet')
    products = [product(p, monitor(n, p['Signature'])) for p in plants for n in NAMES]
    core = next(p for p in plants if p['ID'] == prior.plant_id(prior.CORE))
    products += [product(core, m, d) for d in controllers for m in monitors]
    comparisons = [profile(l, r) for l, r in itertools.combinations(monitors, 2)]
    lookup = {g['ID']: g for g in products}
    contexts = [(p['ID'], p['ID'], '') for p in plants]
    contexts += [(core['ID'] + ' || ' + d['ID'], core['ID'], ' || ' + d['ID']) for d in controllers]
    substitutions = []
    for comp in comparisons:
        if comp['Equal']:
            for context, plant_id, suffix in contexts:
                lp, rp = (plant_id + ' || ' + comp[side] + suffix for side in ('Left', 'Right'))
                lb, rb = lookup[lp]['Nonblocking'], lookup[rp]['Nonblocking']
                require(lb == rb, 'independent profile substitution falsified')
                substitutions.append(dict(Left=comp['Left'], Right=comp['Right'], Context=context,
                                          ProductLeft=lp, ProductRight=rp, NonblockingLeft=lb, NonblockingRight=rb))
    outcomes = [dict(Monitor=m['ID'], Product=(g := lookup[core['ID'] + ' || ' + m['ID']])['ID'],
                     InitialCanComplete=g['Nodes'][0]['CanComplete'], Nonblocking=g['Nonblocking']) for m in monitors]
    negative, loop, strictness = witnesses(products, comparisons)
    counts = dict(Plants=len(plants), Monitors=len(monitors), Controllers=len(controllers), Products=len(products),
                  ProductStates=sum(g['States'] for g in products), ProductArcs=sum(g['Arcs'] for g in products),
                  ProductCanComplete=sum(g['CanComplete'] for g in products),
                  InitialCanComplete=sum(g['Nodes'][0]['CanComplete'] for g in products),
                  NonblockingProducts=sum(g['Nonblocking'] for g in products), Comparisons=len(comparisons),
                  EqualLanguages=sum(c['InitialLanguage']['Equal'] for c in comparisons),
                  EqualProfiles=sum(c['Equal'] for c in comparisons),
                  ResidualComparisons=sum(len(c['Residuals']) for c in comparisons),
                  LanguageSubsetPairs=sum(len(c['InitialLanguage']['Pairs']) + sum(len(r['Language']['Pairs'])
                                          for r in c['Residuals']) for c in comparisons),
                  ProfileSubsetPairs=sum(len(c['Pairs']) for c in comparisons), Substitutions=len(substitutions),
                  NonblockingSubstitutions=sum(s['NonblockingLeft'] for s in substitutions),
                  BlockingSubstitutions=sum(not s['NonblockingLeft'] for s in substitutions))
    return dict(Schema='alan/phase3c4/v1', Input=input_record, Alphabet=alphabet, Plants=plants,
                Monitors=monitors, Controllers=controllers, Products=products, Comparisons=comparisons,
                Substitutions=substitutions, CoreOutcomes=outcomes, Negative=negative, Loop=loop,
                Strictness=strictness, Counts=counts, Sources=sources(input_record))


def audit(evidence, expected=None):
    # Full typed structural equality checks all finite keys, arrays, certificates,
    # arcs and witnesses. No supplied counts, labels or decisions become oracles.
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
