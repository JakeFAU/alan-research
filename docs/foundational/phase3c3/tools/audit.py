#!/usr/bin/env python3
"""Independent reconstruction of the frozen Phase 3C.3 finite evidence.

New automata/plant/language/product logic is specified here, not imported from
Go. Inherited concrete custody replay uses the UNCHANGED Phase 3C.2 Python
checker (source-bound in the audit manifest). Go bytes are only hashed.
"""
import copy
import hashlib
import importlib.util
import itertools
import json
import pathlib
import sys
from collections import defaultdict, deque
from functools import lru_cache

ROOT = pathlib.Path(__file__).resolve().parents[4]
OLD_PATH = ROOT / 'docs/foundational/phase3c2/tools/audit.py'
_spec = importlib.util.spec_from_file_location('phase3c2_audit', OLD_PATH)
old = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(old)
InvalidEvidence = old.InvalidEvidence
require, exact, keys, load, decode = old.require, old.exact, old.keys, old.load, old.decode
KINDS = ('commit', 'reserve-capacity', 'apply', 'consume')
ALPHABET = tuple(k + ':' + n for n in ('AB', 'BC') for k in KINDS)
MONITORS = ('free', 'commit-first', 'release-first', 'release-first-tau', 'release-first-doomed')
CORE = old.config('shared-1', 'early', 0)


def key(q):
    return ','.join(str(v).lower() if type(v) is bool else v for v in q)


def tuple_id(q):
    return json.dumps(q, separators=(',', ':'))


def plant_id(c):
    return f"{c['Policy']}/{c['Layout']}/mask{c['Base']['ServiceMask']}"


def plant_moves(c, q):
    """Declarative finite service/capacity relation, independently from contract."""
    for i, name in enumerate(('AB', 'BC')):
        stage, held = q[i], q[2+i]
        service = i == 1 or bool(c['Base']['ServiceMask'] & 2) or q[1] in ('U', 'D')
        stable = i == 1 or bool(c['Base']['ServiceMask'] & 2) or q[1] == 'D'
        root = c['Memberships'][i]['Root']
        occupancy = sum(q[2+j] for j in (0, 1) if c['Memberships'][j]['Root'] == root)
        limit = next(r['Limit'] for r in c['Roots'] if r['ID'] == root)
        reservation_gate = {'early': True, 'ready-now': service, 'stable-ready': stable}[c['Policy']]
        cases = [('commit', stage == 'U', 'W', held),
                 ('reserve-capacity', stage == 'W' and not held and occupancy < limit and reservation_gate, stage, True),
                 ('apply', stage == 'W' and held and service, 'P', False),
                 ('consume', stage == 'P', 'D', held)]
        for label, enabled, target_stage, target_hold in cases:
            if enabled:
                target = list(q)
                target[i], target[2+i] = target_stage, target_hold
                yield label + ':' + name, tuple(target)


def descriptor(name, role, c, initial, states, edges):
    return dict(ID=name, Role=role, Signature=c, Initial=initial, States=states, Edges=edges)


def atom(s, marked):
    return dict(ID=str(s), Tuple=None, Marked=marked)


def edge(source, label, target):
    return dict(From=str(source), Label=label, To=str(target))


def plant(c):
    queue = [('U', 'U', False, False)]
    seen = set(queue)
    edges = []
    for q in queue:
        for label, target in plant_moves(c, q):
            edges.append(edge(key(q), label, key(target)))
            if target not in seen:
                seen.add(target)
                queue.append(target)
    return descriptor(plant_id(c), 'K', c, key(queue[0]),
                      [atom(key(q), q[:2] == ('D', 'D')) for q in queue], edges)


def monitor(name, c=CORE):
    r, a, commit = 'reserve-capacity:AB', 'apply:AB', 'commit:BC'
    count = 1 if name == 'free' else 4 if name.endswith(('tau', 'doomed')) else 3
    edges = [edge(s, label, s) for s in range(count) for label in ALPHABET
             if label not in (r, a, commit)]
    terminal = 0 if name == 'free' else 2
    edges += [edge(terminal, label, terminal) for label in (r, a, commit)]
    if name != 'free':
        edges.append(edge(0, r, 1))
        if name == 'commit-first':
            edges.append(edge(1, commit, 2))
        elif name == 'release-first-tau':
            edges.extend([edge(1, 'tau', 3), edge(3, a, 2)])
        else:
            edges.append(edge(1, a, 2))
            if name == 'release-first-doomed':
                edges.append(edge(1, 'tau', 3))
    return descriptor(name, 'M', c, '0', [atom(s, True) for s in range(count)], edges)


def controller(code, marks):
    edges = [edge(s, label, s) for s in range(2) for label in ALPHABET
             if label not in ('commit:AB', 'apply:BC')]
    digits = code
    for source, label in itertools.product(range(2), ('commit:AB', 'apply:BC')):
        digit = digits % 3
        digits //= 3
        if digit:
            edges.append(edge(source, label, digit-1))
    return descriptor(f'D-{code:02}-{marks}', 'D', CORE, '0',
                      [atom(s, bool(marks & (1 << s))) for s in range(2)], edges)


def normalized_descriptor(d):
    keys(d, 'ID Role Signature Initial States Edges', 'descriptor')
    out = dict(d)
    require(type(d['States']) is list and type(d['Edges']) is list, 'descriptor arrays')
    out['States'] = sorted(d['States'], key=lambda s: s['ID'])
    out['Edges'] = sorted(d['Edges'], key=lambda e: (e['From'], e['Label'], e['To']))
    return out


def table(d):
    adjacency = defaultdict(list)
    for e in d['Edges']:
        adjacency[(e['From'], e['Label'])].append(e['To'])
    return {s['ID']: s for s in d['States']}, adjacency


def product(*factors):
    """All synchronous visible moves, each independent private move, all arcs."""
    for f in factors[1:]:
        exact(f['Signature'], factors[0]['Signature'], 'product signature')
    tables = [table(f) for f in factors]
    initial = tuple(f['Initial'] for f in factors)
    pending, seen, nodes = [initial], {initial}, []

    def flatten(q):
        return tuple(v for i, state in enumerate(q)
                     for v in (tables[i][0][state]['Tuple'] or [state]))

    for q in pending:
        transitions = set()
        for label in ALPHABET:
            for target in itertools.product(*(t[1].get((s, label), []) for t, s in zip(tables, q))):
                transitions.add((label, target))
        for i, (t, s) in enumerate(zip(tables, q)):
            for dest in t[1].get((s, 'tau'), []):
                target = q[:i] + (dest,) + q[i+1:]
                transitions.add(('tau', target))
        # Discovery order is immaterial to relation equality; shortest path layers
        # remain BFS. Serialized arcs use the schema's label/target order.
        ordered = sorted(transitions, key=lambda x: (x[0], tuple_id(flatten(x[1]))))
        for _, target in ordered:
            if target not in seen:
                seen.add(target)
                pending.append(target)
        flat = flatten(q)
        nodes.append(dict(ID=tuple_id(flat), Tuple=list(flat),
                          Marked=all(t[0][s]['Marked'] for t, s in zip(tables, q)),
                          CanComplete=False,
                          Arcs=[dict(Label=l, To=tuple_id(flatten(t))) for l, t in ordered]))
    reverse = defaultdict(set)
    for node in nodes:
        for arc in node['Arcs']:
            reverse[arc['To']].add(node['ID'])
    can = {n['ID'] for n in nodes if n['Marked']}
    work = list(can)
    for state in work:
        for predecessor in reverse[state] - can:
            can.add(predecessor)
            work.append(predecessor)
    for n in nodes:
        n['CanComplete'] = n['ID'] in can
    return dict(ID=' || '.join(f['ID'] for f in factors), Factors=[f['ID'] for f in factors],
                Initial=tuple_id(flatten(initial)), Nodes=nodes, States=len(nodes),
                Arcs=sum(len(n['Arcs']) for n in nodes), CanComplete=len(can),
                Nonblocking=len(can) == len(nodes))


def as_factor(g, signature):
    return descriptor(g['ID'], 'product', signature, g['Initial'],
                      [dict(ID=n['ID'], Tuple=n['Tuple'], Marked=n['Marked']) for n in g['Nodes']],
                      [edge(n['ID'], a['Label'], a['To']) for n in g['Nodes'] for a in n['Arcs']])


def normalized_graph(g):
    keys(g, 'ID Factors Initial Nodes States Arcs CanComplete Nonblocking', 'graph')
    require(type(g['Nodes']) is list and g['Nodes'], 'graph.Nodes: nonempty array required')
    # Check the declared BFS ordering by distance; sibling order is unspecified.
    require(g['Nodes'][0]['ID'] == g['Initial'], 'graph: initial must be first')
    lookup = {n['ID']: n for n in g['Nodes']}
    require(len(lookup) == len(g['Nodes']), 'graph: duplicate node')
    distances, pending = {g['Initial']: 0}, [g['Initial']]
    for s in pending:
        for a in lookup[s]['Arcs']:
            require(a['To'] in lookup, 'graph: undeclared arc target')
            if a['To'] not in distances:
                distances[a['To']] = distances[s] + 1
                pending.append(a['To'])
    require(len(distances) == len(lookup), 'graph: unreachable node')
    layer = [distances[n['ID']] for n in g['Nodes']]
    require(layer == sorted(layer), 'graph: not BFS ordered')
    return dict(g, Nodes=sorted(g['Nodes'], key=lambda n: n['ID']))


def closure(d, states):
    _, adj = table(d)
    reached = set(states)
    work = list(states)
    for s in work:
        for target in adj.get((s, 'tau'), []):
            if target not in reached:
                reached.add(target)
                work.append(target)
    return tuple(sorted(reached))


def advance(d, states, label):
    _, adj = table(d)
    return closure(d, {t for s in states for t in adj.get((s, label), [])})


def accepting(d, states):
    marks = {s['ID'] for s in d['States'] if s['Marked']}
    return bool(marks.intersection(states))


def accepts(d, word):
    states = closure(d, [d['Initial']])
    for label in word:
        states = advance(d, states, label)
    return accepting(d, states)


def comparison(left, right):
    start = (closure(left, [left['Initial']]), closure(right, [right['Initial']]))
    pending, indexes, words, pairs = [start], {start: 0}, [[]], []
    separator = None
    for index, (ls, rs) in enumerate(pending):
        la, ra = accepting(left, ls), accepting(right, rs)
        if la != ra and separator is None:
            separator = (words[index], la, ra)
        arcs = []
        for label in ALPHABET:
            target = (advance(left, ls, label), advance(right, rs, label))
            if target not in indexes:
                indexes[target] = len(pending)
                pending.append(target)
                words.append(words[index] + [label])
            arcs.append(dict(Label=label, To=indexes[target]))
        pairs.append(dict(Left=list(ls), Right=list(rs), LeftAccept=la, RightAccept=ra, Arcs=arcs))
    word, la, ra = separator if separator is not None else ([], False, False)
    return dict(Left=left['ID'], Right=right['ID'], Equal=separator is None, Word=word,
                LeftAccept=la, RightAccept=ra, Pairs=pairs)


@lru_cache(maxsize=1)
def reconstruct():
    plants = [plant(old.config(layout, policy, mask)) for layout in old.LAYOUTS
              for policy in old.POLICIES for mask in range(8)]
    monitors = [monitor(n) for n in MONITORS]
    controllers = [controller(code, marks) for code in range(81) for marks in range(4)]
    products = [product(p, monitor(n, p['Signature'])) for p in plants for n in MONITORS]
    groupings = []
    for d in controllers:
        for m in monitors:
            g = product(plants[0], m, d)
            products.append(g)
            left = product(as_factor(product(plants[0], m), CORE), d)
            right = product(plants[0], as_factor(product(m, d), CORE))
            for grouped in (left, right):
                canonical = dict(grouped, Factors=g['Factors'])
                exact(normalized_graph(canonical), normalized_graph(g), 'independent associativity')
            groupings.append(dict(Product=g['ID'], Left=left, Right=right, Equal=True))
    comparisons = [comparison(l, r) for l, r in itertools.combinations(monitors, 2)]
    lookup = {g['ID']: g for g in products}
    substitutions = []
    contexts = [(p['ID'], p['ID'], '') for p in plants]
    contexts += [(plants[0]['ID'] + ' || ' + d['ID'], plants[0]['ID'], ' || ' + d['ID']) for d in controllers]
    for comp in comparisons:
        if not comp['Equal']:
            continue
        for context, p, suffix in contexts:
            lp, rp = p + ' || ' + comp['Left'] + suffix, p + ' || ' + comp['Right'] + suffix
            lc, rc = lookup[lp]['Nodes'][0]['CanComplete'], lookup[rp]['Nodes'][0]['CanComplete']
            require(lc == rc, 'independently reconstructed substitution failed')
            substitutions.append(dict(Left=comp['Left'], Right=comp['Right'], Context=context,
                                      ProductLeft=lp, ProductRight=rp, CanCompleteLeft=lc, CanCompleteRight=rc))
    return plants, monitors, controllers, products, groupings, comparisons, substitutions


def source_records():
    paths = ['go.mod', 'go.sum'] + ['docs/foundational/phase3c3/' + n for n in
                                       ('00-contract.md', '01-argument.md', '02-evidence-schema.md')]
    for directory in ('research/componentinterfaces', 'research/componentinterfaces/cmd/export',
                      'research/capacitycomposition', 'research/contractcomposition', 'research/twoagent'):
        paths += [str(p.relative_to(ROOT)) for p in (ROOT / directory).glob('*.go')
                  if not p.name.endswith('_test.go')]
    return [dict(Path=p, SHA256=hashlib.sha256((ROOT / p).read_bytes()).hexdigest()) for p in sorted(paths)]


def projection(s):
    return tuple_id((key(old.project(s['Old'])), s['Monitor']))


def replay(h, monitor_name, c, graph, path):
    keys(h, 'ID Config Monitor Initial Steps Cuts Complete CanComplete', path)
    exact(h['Config'], c, path + '.Config')
    exact(h['Monitor'], monitor_name, path + '.Monitor')
    current = dict(Old=old.initial(), Monitor='0')
    exact(h['Initial'], current, path + '.Initial')
    automaton = monitor(monitor_name, c)
    _, adj = table(automaton)
    nodes = {n['ID']: n for n in graph['Nodes']}
    require(type(h['Steps']) is list, path + '.Steps: array required')
    states = [current]
    for i, step in enumerate(h['Steps']):
        p = f'{path}.Steps[{i}]'
        keys(step, 'Event Private State', p)
        e = step['Event']
        keys(e, 'Kind Actor Original Binding Result Duplicate', p + '.Event')
        private = e['Kind'] == 'tau'
        exact(step['Private'], private, p + '.Private')
        pc = current['Monitor']
        if private:
            exact(e, dict(Kind='tau', Actor='B', Original='', Binding=old.zero_binding(),
                          Result='', Duplicate=False), p + '.Event')
            require(not current['Old']['Base']['Agents'][1]['Down'], p + ': private actor B down')
            targets = adj.get((pc, 'tau'), [])
            require(len(targets) == 1, p + ': disabled private step')
            after = dict(Old=copy.deepcopy(current['Old']), Monitor=targets[0])
            label = 'tau'
        else:
            after = dict(Old=old.step(c, current['Old'], e, p + '.Event'), Monitor=pc)
            label = e['Kind'] + ':' + e['Original']
            if label in ALPHABET:
                targets = adj.get((pc, label), [])
                require(len(targets) == 1, p + ': monitor forbids visible step')
                after['Monitor'] = targets[0]
                if e['Kind'] == 'apply':
                    exact(e['Result'], 'r1', p + ': deterministic first application result')
        exact(step['State'], after, p + '.State')
        source, target = projection(current), projection(after)
        require(target in nodes, p + ': projection outside product')
        if private or label in ALPHABET:
            require(dict(Label=label, To=target) in nodes[source]['Arcs'], p + ': product projection edge')
        else:
            exact(target, source, p + ': transport must stutter')
        states.append(after)
        current = after
    cuts = []
    require(type(h['Cuts']) is list, path + '.Cuts: array required')
    expected_names = ['after-reserve-AB']
    if monitor_name == 'release-first-tau':
        expected_names += ['B-down-before-private', 'private-after-recovery']
    if h['ID'] == 'history-8':
        expected_names += ['doomed-before-crash', 'doomed-after-recovery']
    expected_names.append('final')
    require(len(h['Cuts']) == len(expected_names), path + ': cut count')
    for cut, name in zip(h['Cuts'], expected_names):
        keys(cut, 'Name AfterSteps Node Complete CanComplete', path + '.Cut')
        old.integer(cut['AfterSteps'], path + '.Cut.AfterSteps')
        index = cut['AfterSteps']
        require(index < len(states), path + ': cut index outside history')
        s = states[index]
        n = nodes[projection(s)]
        exact(cut, dict(Name=name, AfterSteps=index, Node=n['ID'], Complete=n['Marked'],
                        CanComplete=n['CanComplete']), path + '.Cut.' + name)
        if name == 'after-reserve-AB':
            exact(old.project(s['Old']), ('W', 'U', True, False), path + ': common-cut projection')
            require(index > 0 and h['Steps'][index-1]['Event']['Kind'] == 'reserve-capacity'
                    and h['Steps'][index-1]['Event']['Original'] == 'AB', path + ': common-cut boundary')
        elif name == 'final':
            exact(index, len(states)-1, path + ': final cut index')
        elif name == 'B-down-before-private':
            exact(s['Old']['Base']['Agents'][1]['Down'], True, path + ': B down cut')
            exact(s['Monitor'], '1', path + ': retained pre-private PC')
        elif name == 'private-after-recovery':
            exact(s['Monitor'], '3', path + ': private PC')
            exact(s['Old']['Base']['Agents'][1], dict(Down=False, Crashed=True, Activation=1), path + ': recovery')
        elif name.startswith('doomed-'):
            exact(s['Monitor'], '3', path + ': doomed PC')
            exact(n['CanComplete'], False, path + ': doomed graph label')
            if name.endswith('recovery'):
                exact(s['Old']['Base']['Agents'][1], dict(Down=False, Crashed=True, Activation=1), path + ': doomed recovery')
            else:
                exact(s['Old']['Base']['Agents'][1], dict(Down=False, Crashed=False, Activation=0), path + ': before crash')
        cuts.append(s)
    n = nodes[projection(current)]
    exact(h['Complete'], n['Marked'], path + '.Complete')
    exact(h['CanComplete'], n['CanComplete'], path + '.CanComplete')
    exact(h['Complete'], h['ID'] not in ('history-1', 'history-8'), path + ': roster completion')
    return states, cuts


def audit(evidence):
    keys(evidence, 'Schema Alphabet Monitors Plants Controllers Products Comparisons Groupings '
                   'Substitutions Histories Counts Sources', 'export')
    exact(evidence['Schema'], 'alan/phase3c3/v1', 'Schema')
    alphabet = [dict(ID=k + ':' + n, Action=old.action(k, i),
                     Actor='ABC'[i + (k in ('reserve-capacity', 'apply'))], Root='Q')
                for i, n in enumerate(('AB', 'BC')) for k in KINDS]
    exact(evidence['Alphabet'], alphabet, 'Alphabet')
    exact(evidence['Sources'], source_records(), 'Sources')
    plants, monitors, controllers, products, groupings, comparisons, substitutions = reconstruct()
    for name, expected in [('Plants', plants), ('Monitors', monitors), ('Controllers', controllers)]:
        require(type(evidence[name]) is list and len(evidence[name]) == len(expected), name + ': roster length')
        for i, (actual, wanted) in enumerate(zip(evidence[name], expected)):
            exact(normalized_descriptor(actual), normalized_descriptor(wanted), f'{name}[{i}]')
    require(type(evidence['Products']) is list and len(evidence['Products']) == len(products), 'Products: roster length')
    for i, (actual, wanted) in enumerate(zip(evidence['Products'], products)):
        exact(normalized_graph(actual), normalized_graph(wanted), f'Products[{i}]')
    require(type(evidence['Groupings']) is list and len(evidence['Groupings']) == len(groupings), 'Groupings: roster length')
    for i, (actual, wanted) in enumerate(zip(evidence['Groupings'], groupings)):
        keys(actual, 'Product Left Right Equal', f'Groupings[{i}]')
        exact(actual['Product'], wanted['Product'], f'Groupings[{i}].Product')
        exact(actual['Equal'], True, f'Groupings[{i}].Equal')
        for side in ('Left', 'Right'):
            exact(normalized_graph(actual[side]), normalized_graph(wanted[side]), f'Groupings[{i}].{side}')
    exact(evidence['Comparisons'], comparisons, 'Comparisons')
    # The substitution schema defines a complete set, without prescribing order.
    require(type(evidence['Substitutions']) is list, 'Substitutions: array required')
    order = lambda r: (r['Left'], r['Right'], r['Context'])
    exact(sorted(evidence['Substitutions'], key=order), sorted(substitutions, key=order), 'Substitutions')
    histories = evidence['Histories']
    require(type(histories) is list and len(histories) == 8, 'Histories: roster length')
    roster = [('commit-first', 'shared-1', 0), ('release-first', 'shared-1', 0),
              ('commit-first', 'shared-2', 0), ('commit-first', 'separate-1', 0),
              ('commit-first', 'shared-1', 2), ('release-first-tau', 'shared-1', 0),
              ('release-first-doomed', 'shared-1', 0), ('release-first-doomed', 'shared-1', 0)]
    lookup = {g['ID']: g for g in products}
    replays = []
    for i, (h, (name, layout, mask)) in enumerate(zip(histories, roster)):
        exact(h['ID'], f'history-{i+1}', f'Histories[{i}].ID')
        c = old.config(layout, 'early', mask)
        replays.append(replay(h, name, c, lookup[plant_id(c) + ' || ' + name], f'Histories[{i}]'))
    exact(replays[0][1][0]['Old'], replays[1][1][0]['Old'], 'common-cut entire Old state')
    exact([h['Cuts'][0]['CanComplete'] for h in histories[:2]], [False, True], 'common-cut completion separation')
    exact(old.project(replays[0][0][-1]['Old']), ('W', 'W', True, False), 'blocked BC final')
    exact(histories[0]['Steps'][-1]['Event']['Kind'], 'commit', 'blocked BC final commitment')
    # Recovery witnesses must actually contain the specified physical schedule.
    for i, sequence in [(5, ['crash', 'recover', 'tau']), (7, ['tau', 'crash', 'recover'])]:
        physical = [s['Event'] for s in histories[i]['Steps'] if s['Event']['Kind'] in ('crash', 'recover', 'tau')]
        exact([e['Kind'] for e in physical], sequence, f'Histories[{i}]: physical order')
        exact([e['Actor'] for e in physical], ['B'] * 3, f'Histories[{i}]: physical actor')
    counts = dict(Plants=len(plants), Monitors=len(monitors), Controllers=len(controllers),
                  Products=len(products), ProductStates=sum(g['States'] for g in products),
                  ProductArcs=sum(g['Arcs'] for g in products),
                  ProductCanComplete=sum(g['CanComplete'] for g in products),
                  InitialCanComplete=sum(g['Nodes'][0]['CanComplete'] for g in products),
                  NonblockingProducts=sum(g['Nonblocking'] for g in products),
                  Comparisons=len(comparisons), EqualComparisons=sum(c['Equal'] for c in comparisons),
                  SubsetPairs=sum(len(c['Pairs']) for c in comparisons), Groupings=len(groupings),
                  Substitutions=len(substitutions),
                  SuccessfulSubstitutions=sum(s['CanCompleteLeft'] for s in substitutions),
                  BlockedSubstitutions=sum(not s['CanCompleteLeft'] for s in substitutions),
                  Histories=len(histories), ConcreteSteps=sum(len(h['Steps']) for h in histories))
    exact(evidence['Counts'], counts, 'Counts')
    return counts


def main():
    if len(sys.argv) != 2:
        print('usage: audit.py PATH', file=sys.stderr)
        return 2
    try:
        counts = audit(load(sys.argv[1]))
    except (InvalidEvidence, KeyError, TypeError, ValueError, OSError, IndexError) as exc:
        print('REJECT: ' + str(exc), file=sys.stderr)
        return 1
    print('ACCEPT: ' + json.dumps(counts, sort_keys=True))
    return 0


if __name__ == '__main__':
    sys.exit(main())
