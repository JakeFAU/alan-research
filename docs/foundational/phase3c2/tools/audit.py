#!/usr/bin/env python3
"""Independent declarative audit of alan/phase3c2/v1 (standard library only).

No producer/runtime code is loaded or executed. Source bytes are hashed only.
Operational probes are structurally checked reports, not independently rerun.
"""
import copy
import gzip
import hashlib
import json
import pathlib
import sys
from collections import deque

ROOT = pathlib.Path(__file__).resolve().parents[4]
NAMES = ('AB', 'BC')
LAYOUTS = ('shared-1', 'shared-2', 'separate-1')
POLICIES = ('early', 'ready-now', 'stable-ready')
WITNESSES = ('early-AB-holder', 'early-BC-holder', 'early-BC-completed',
             'ready-now-AB-holder', 'ready-now-BC-holder', 'ready-now-BC-completed',
             'old-AB-before-BC-commit', 'early-AB-holder-crashed',
             'stable-recovery-custody-r1', 'stable-recovery-custody-r2')


class InvalidEvidence(ValueError):
    """Malformed or semantically false evidence, with a diagnostic location."""


def require(condition, message):
    if not condition:
        raise InvalidEvidence(message)


def exact(actual, expected, path):
    """Structural equality with exact types (not Python's True == 1)."""
    require(type(actual) is type(expected), path + ': type mismatch')
    if isinstance(expected, dict):
        require(actual.keys() == expected.keys(), path + ': keys mismatch')
        for key, value in expected.items():
            exact(actual[key], value, path + '.' + key)
    elif isinstance(expected, list):
        require(len(actual) == len(expected), path + ': length mismatch')
        for i, value in enumerate(expected):
            exact(actual[i], value, path + '[' + str(i) + ']')
    else:
        require(actual == expected, path + ': value mismatch')


def keys(record, fields, path):
    require(type(record) is dict and set(record) == set(fields.split()), path + ': keys mismatch')


def integer(value, path):
    require(type(value) is int and value >= 0, path + ': nonnegative integer required')


def decode(raw):
    def pairs(items):
        obj = {}
        for key, value in items:
            require(key not in obj, 'duplicate key: ' + key)
            obj[key] = value
        return obj

    def constant(value):
        raise InvalidEvidence('nonfinite JSON number: ' + value)

    try:
        return json.loads(raw.decode('utf-8'), object_pairs_hook=pairs, parse_constant=constant)
    except (ValueError, UnicodeError) as exc:
        raise InvalidEvidence(str(exc)) from exc


def load(path):
    raw = pathlib.Path(path).read_bytes()
    if raw.startswith(b'\x1f\x8b'):
        raw = gzip.decompress(raw)
    return decode(raw)


def binding(i):
    caller, callee = 'ABC'[i:i+2]
    return dict(Original=NAMES[i], Caller=caller, CallerGeneration=1, Operation=1,
                Callee=callee, CalleeGeneration=1, Request='increment',
                Grant=f'grant:{caller}/g1->{callee}/g1:counter-{callee}', Resource='counter-' + callee)


def zero_binding():
    return {k: (0 if type(v) is int else '') for k, v in binding(0).items()}


def config(layout, policy, mask):
    roots = ([dict(ID='Q-AB', Limit=1), dict(ID='Q-BC', Limit=1)]
             if layout == 'separate-1' else [dict(ID='Q', Limit=int(layout[-1]))])
    memberships = [dict(Original=n, Root=('Q-' + n if layout == 'separate-1' else 'Q')) for n in NAMES]
    return dict(Base=dict(Topology='chain', ServiceMask=mask,
                          Agents=[dict(ID=a, Generation=1, Servicing=bool(mask & (1 << i)))
                                  for i, a in enumerate('ABC')],
                          Originals=[binding(0), binding(1)]),
                Layout=layout, Policy=policy, Roots=roots, Memberships=memberships)


def config_id(c):
    return f"{c['Layout']}/{c['Policy']}/{c['Base']['ServiceMask']}"


def action(kind, i):
    return dict(Kind=kind, Original=NAMES[i], Binding=binding(i))


def state(q):
    return dict(Base=dict(Stages=[q[0], q[1], '']), Held=[q[2], q[3]])


def ready(c, q, i, stable=False):
    return i == 1 or c['Base']['Agents'][1]['Servicing'] or (
        q[1] == 'D' if stable else q[1] not in ('W', 'P'))


def transitions(c, q):
    """Full named abstract relation, including arcs to previously visited states."""
    for i in range(2):
        target = list(q)
        if q[i] == 'U':
            kind, target[i] = 'commit', 'W'
        elif q[i] == 'W' and q[i+2] and ready(c, q, i):
            kind, target[i], target[i+2] = 'apply', 'P', False
        elif q[i] == 'P':
            kind, target[i] = 'consume', 'D'
        else:
            continue
        yield action(kind, i), tuple(target)
    for i in range(2):
        root = c['Memberships'][i]['Root']
        limit = next(r['Limit'] for r in c['Roots'] if r['ID'] == root)
        used = sum(q[j+2] for j in range(2) if c['Memberships'][j]['Root'] == root)
        policy = c['Policy']
        allowed = policy == 'early' or ready(c, q, i, stable=policy == 'stable-ready')
        if q[i] == 'W' and not q[i+2] and used < limit and allowed:
            target = list(q)
            target[i+2] = True
            yield action('reserve-capacity', i), tuple(target)


def graph(c):
    initial = ('U', 'U', False, False)
    queue = [initial]
    indexes = {initial: 0}
    parents = [(-1, None)]
    adjacency = []
    for index, q in enumerate(queue):
        arcs = []
        for a, target in transitions(c, q):
            if target not in indexes:
                indexes[target] = len(queue)
                queue.append(target)
                parents.append((index, a))
            arcs.append(dict(Action=a, To=indexes[target]))
        adjacency.append(arcs)
    reverse = [[] for _ in queue]
    for index, arcs in enumerate(adjacency):
        for arc in arcs:
            reverse[arc['To']].append(index)
    complete = {i for i, q in enumerate(queue) if q[:2] == ('D', 'D')}
    can = set(complete)
    pending = deque(complete)
    while pending:
        for predecessor in reverse[pending.popleft()]:
            if predecessor not in can:
                can.add(predecessor)
                pending.append(predecessor)
    nodes = [dict(State=state(q), Arcs=adjacency[i], Complete=i in complete,
                  CanComplete=i in can, Parent=parents[i][0], Via=parents[i][1])
             for i, q in enumerate(queue)]
    return dict(Config=c, Nodes=nodes, States=len(queue), Arcs=sum(map(len, adjacency)),
                CanComplete=len(can), CannotComplete=len(queue)-len(can)), indexes


def component(c, owners):
    agents = c['Base']['Agents']
    originals = [i for i in range(2) if i in owners]
    ports = sorted({i+1 for i in originals} - set(owners))
    memberships = [c['Memberships'][i] for i in originals]
    roots = sorted({m['Root'] for m in memberships})
    return dict(Base=dict(Owners=[agents[i] for i in owners], Ports=[agents[i] for i in ports],
                          Originals=[binding(i) for i in originals]),
                Layout=c['Layout'], Policy=c['Policy'],
                Roots=[r for name in roots for r in c['Roots'] if r['ID'] == name],
                Memberships=memberships)


def grouping(c):
    out = dict(ConfigID=config_id(c), Flat=c, LeftClosed=c, RightClosed=c, Equal=True)
    for name, owners in [('A', [0]), ('B', [1]), ('C', [2]), ('AB', [0, 1]),
                         ('BC', [1, 2]), ('Left', [0, 1, 2]), ('Right', [0, 1, 2])]:
        out[name] = component(c, owners)
    return out


def empty_transport(b):
    return dict(Binding=b, Committed=False, Admitted=False, Applied=False, Consumed=False,
                Volatile='', Recorded='', Retained='', RequestChannel=False,
                ResultChannel='', ACKChannel=False, WorldResult='')


def initial():
    return dict(Base=dict(Agents=[dict(Down=False, Crashed=False, Activation=0) for _ in range(3)],
                          Originals=[empty_transport(binding(0)), empty_transport(binding(1)),
                                     empty_transport(zero_binding())]), Held=[False, False])


def project(s):
    stages = tuple('D' if t['Consumed'] else 'P' if t['Applied'] else 'W' if t['Committed'] else 'U'
                   for t in s['Base']['Originals'][:2])
    return stages + tuple(s['Held'])


def step(c, before, e, path):
    """Generate a post-state from a causal event, never from exported post-state facts."""
    keys(e, 'Kind Actor Original Binding Result Duplicate', path)
    for field in ('Kind', 'Actor', 'Original', 'Result'):
        require(type(e[field]) is str, path + '.' + field + ': string required')
    require(type(e['Duplicate']) is bool, path + '.Duplicate: boolean required')
    require(e['Actor'] in 'ABC' and len(e['Actor']) == 1, path + ': actor')
    after = copy.deepcopy(before)
    kind = e['Kind']
    owner = 'ABC'.index(e['Actor'])
    availability = after['Base']['Agents'][owner]
    if kind in ('crash', 'recover'):
        exact(e, dict(Kind=kind, Actor=e['Actor'], Original='', Binding=zero_binding(),
                      Result='', Duplicate=False), path)
        if kind == 'crash':
            require(not availability['Crashed'] and not availability['Down'], path + ': crash phase')
            availability.update(Down=True, Crashed=True)
            if owner < 2:
                after['Base']['Originals'][owner]['Volatile'] = ''
        else:
            require(availability['Down'], path + ': recover phase')
            availability.update(Down=False, Activation=1)
        return after
    require(kind in ('commit', 'send-request', 'admit', 'deduplicate', 'reserve-capacity',
                     'apply', 'send-result', 'receive-result', 'record', 'consume',
                     'send-ack', 'release'), path + ': event kind')
    require(e['Original'] in NAMES, path + ': original')
    i = NAMES.index(e['Original'])
    exact(e['Binding'], binding(i), path + '.Binding')
    caller = kind in ('commit', 'send-request', 'receive-result', 'record', 'consume', 'send-ack')
    require(owner == i + (0 if caller else 1), path + ': actor')
    require(not availability['Down'], path + ': actor down')
    duplicate = e['Duplicate']
    require(not duplicate or kind in ('admit', 'deduplicate', 'receive-result', 'release'),
            path + ': invalid duplicate')
    result = e['Result']
    if kind in ('apply', 'send-result', 'receive-result', 'record', 'consume'):
        require(result in ('r1', 'r2'), path + ': result value')
    else:
        require(result == '', path + ': unexpected result')
    t = after['Base']['Originals'][i]
    if kind in ('commit', 'reserve-capacity', 'apply', 'consume'):
        choices = {a['Kind']: q for a, q in transitions(c, project(before)) if a['Original'] == NAMES[i]}
        diagnostic = 'reservation' if kind == 'reserve-capacity' else kind
        require(kind in choices, path + ': illegal ' + diagnostic)
        if kind in ('reserve-capacity', 'apply'):
            require(t['Admitted'], path + ': admission required')
        if kind == 'commit':
            t['Committed'] = True
        elif kind == 'reserve-capacity':
            after['Held'][i] = True
        elif kind == 'apply':
            t.update(Applied=True, Retained=result, WorldResult=result)
            after['Held'][i] = False
        else:
            require(t['Recorded'] != '', path + ': consume before durable record')
            require(result == t['Recorded'], path + ': result provenance')
            t['Consumed'] = True
    elif kind == 'send-request':
        require(t['Committed'] and not t['RequestChannel'], path + ': request custody/channel')
        t['RequestChannel'] = True
    elif kind in ('admit', 'deduplicate'):
        require(t['RequestChannel'], path + ': admission channel')
        require(t['Admitted'] == (kind == 'deduplicate'), path + ': admission/deduplication phase')
        t['Admitted'] = True
        t['RequestChannel'] = duplicate
    elif kind == 'send-result':
        require(t['Retained'] != '' and t['ResultChannel'] == '', path + ': send-result custody/channel')
        require(result == t['Retained'], path + ': result provenance')
        t['ResultChannel'] = result
    elif kind == 'receive-result':
        require(t['ResultChannel'] != '', path + ': receive-result channel')
        require(result == t['ResultChannel'], path + ': result provenance')
        t['Volatile'] = result
        if not duplicate:
            t['ResultChannel'] = ''
    elif kind == 'record':
        require(t['Volatile'] != '' and t['Recorded'] == '', path + ': record custody')
        require(result == t['Volatile'], path + ': result provenance')
        t['Recorded'] = result
    elif kind == 'send-ack':
        require(t['Recorded'] != '' and not t['ACKChannel'], path + ': durable ACK')
        t['ACKChannel'] = True
    elif kind == 'release':
        require(t['ACKChannel'], path + ': release ACK')
        t['Retained'] = ''
        t['ACKChannel'] = duplicate
    return after


def replay(w, c, g, indexes, path):
    keys(w, 'Name Config Initial Steps Projection Complete CanComplete', path)
    exact(w['Config'], c, path + '.Config')
    current = initial()
    exact(w['Initial'], current, path + '.Initial')
    require(type(w['Steps']) is list, path + '.Steps: array required')
    projection = []
    for i, entry in enumerate(w['Steps']):
        p = path + '.Steps[' + str(i) + ']'
        keys(entry, 'Event State', p)
        after = step(c, current, entry['Event'], p + '.Event')
        exact(entry['State'], after, p + '.State')
        if project(after) != project(current):
            a = action(entry['Event']['Kind'], NAMES.index(entry['Event']['Original']))
            require((a, project(after)) in list(transitions(c, project(current))), p + ': projection relation')
            projection.append(a)
        current = after
    exact(w['Projection'], projection, path + '.Projection')
    q = project(current)
    require(q in indexes, path + ': unreachable projection')
    node = g['Nodes'][indexes[q]]
    exact(w['Complete'], node['Complete'], path + '.Complete')
    exact(w['CanComplete'], node['CanComplete'], path + '.CanComplete')
    return current


def witness_relations(ws, endpoints):
    for a, b in ((0, 1), (3, 4)):
        exact(endpoints[a]['Base'], endpoints[b]['Base'], 'holder-pair Base')
        exact(endpoints[a]['Held'], [True, False], 'holder-pair AB')
        exact(endpoints[b]['Held'], [False, True], 'holder-pair BC')
        require(sum(endpoints[a]['Held']) == sum(endpoints[b]['Held']) == 1, 'holder-pair totals')
        exact([ws[a]['CanComplete'], ws[b]['CanComplete']], [False, True], 'holder-pair completion')
        exact(list(project(endpoints[a])[:2]), ['W', 'W'], 'holder-pair waiting')
        require(all(t['Admitted'] for t in endpoints[a]['Base']['Originals'][:2]), 'holder-pair admission')
    for prefix, full in ((1, 2), (4, 5), (0, 7)):
        exact(ws[full]['Steps'][:len(ws[prefix]['Steps'])], ws[prefix]['Steps'], 'witness prefix')
    for i in (2, 5, 8, 9):
        exact(ws[i]['Complete'], True, 'completed witness')
    exact(len(ws[7]['Steps']), len(ws[0]['Steps']) + 1, 'crash extension length')
    exact(ws[7]['Steps'][-1]['Event']['Kind'], 'crash', 'crash extension')
    exact(ws[7]['Steps'][-1]['Event']['Actor'], 'B', 'crash owner')
    exact(endpoints[7]['Held'], endpoints[0]['Held'], 'crash retains hold')
    exact(ws[7]['CanComplete'], False, 'crash cannot complete')
    ready_w = ws[3]
    exact(ready_w['Steps'][3]['Event']['Kind'], 'reserve-capacity', 'ready-now reserve index')
    exact(ready_w['Steps'][4]['Event']['Kind'], 'commit', 'ready-now invalidation index')
    require(ready(ready_w['Config'], project(ready_w['Steps'][3]['State']), 0), 'ready before BC commit')
    require(not ready(ready_w['Config'], project(ready_w['Steps'][4]['State']), 0), 'ready invalidation')
    old = ws[6]
    require(any(s['Event']['Kind'] == 'apply' and s['Event']['Original'] == 'AB'
                and project(s['State'])[1] == 'U' for s in old['Steps']), 'old AB-before-BC witness')
    # Trace-language membership with epsilon closure over reserve only. This
    # checks the stated erased observation, not just a local missing guard.
    c = config('shared-1', 'stable-ready', 0)
    frontier = {('U', 'U', False, False)}
    for observed in (a for a in old['Projection'] if a['Kind'] != 'reserve-capacity'):
        closure = set(frontier)
        pending = list(frontier)
        for q in pending:
            for a, t in transitions(c, q):
                if a['Kind'] == 'reserve-capacity' and t not in closure:
                    closure.add(t)
                    pending.append(t)
        frontier = {t for q in closure for a, t in transitions(c, q) if a == observed}
    require(not frontier, 'stable-ready wrongly admits erased old trace')
    # Opaque-value symmetry is checked over complete records, not only labels.
    corresponding = json.loads(json.dumps(ws[8]).replace('r1', 'r2'))
    exact(ws[9], corresponding, 'corresponding recovery histories')
    for number, value in ((8, 'r1'), (9, 'r2')):
        w = ws[number]
        events = [s['Event'] for s in w['Steps']]
        require({e['Original'] for e in events if e['Kind'] == 'apply' and e['Result'] == value}
                == set(NAMES), 'recovery actual applications')
        require(all(e['Result'] in ('', value) for e in events), 'recovery result symmetry')
        for kind in ('deduplicate', 'receive-result', 'release'):
            require({e['Original'] for e in events if e['Kind'] == kind and e['Duplicate']}
                    == set(NAMES), 'recovery duplicate coverage')
        for original in NAMES:
            releases = [j for j, e in enumerate(events) if e['Kind'] == 'release'
                        and e['Original'] == original]
            consumes = [j for j, e in enumerate(events) if e['Kind'] == 'consume'
                        and e['Original'] == original]
            require(releases and consumes and max(releases) < min(consumes),
                    'recovery release before consumption')
        lost = set()
        for j, e in enumerate(events):
            before = w['Initial'] if j == 0 else w['Steps'][j-1]['State']
            if e['Kind'] == 'crash' and e['Actor'] in ('A', 'B'):
                t = before['Base']['Originals']['AB'.index(e['Actor'])]
                if t['Volatile'] == value and t['Recorded'] == '':
                    lost.add(e['Actor'])
        require(lost == {'A', 'B'}, 'recovery pre-record losses')
        require(any(e['Kind'] == 'crash' and e['Actor'] == 'C'
                    and w['Steps'][j]['State']['Held'] == [False, True]
                    for j, e in enumerate(events)), 'recovery C holder survives')
        require(all(t['Recorded'] == value and t['Retained'] == ''
                    for t in endpoints[number]['Base']['Originals'][:2]), 'recovery durable release custody')


def sources(root):
    paths = [root / 'go.mod', root / 'go.sum']
    paths += [root / 'docs/foundational/phase3c2' / (n + '.md')
              for n in ('00-contract', '01-argument', '02-evidence-schema')]
    for package in ('capacitycomposition', 'contractcomposition', 'twoagent'):
        paths += list((root / 'research' / package).rglob('*.go'))
    return [dict(Path=p.relative_to(root).as_posix(), SHA256=hashlib.sha256(p.read_bytes()).hexdigest())
            for p in sorted(paths)]


def check(d, root=ROOT):
    keys(d, 'Schema Sources Graphs Groupings Witnesses Probes', 'Evidence')
    exact(d['Schema'], 'alan/phase3c2/v1', 'Schema')
    exact(d['Sources'], sources(root), 'Sources')
    configs = [config(l, p, m) for l in LAYOUTS for p in POLICIES for m in range(8)]
    require(type(d['Graphs']) is list and len(d['Graphs']) == 72, 'Graphs: roster')
    require(type(d['Groupings']) is list and len(d['Groupings']) == 72, 'Groupings: roster')
    reconstructed = {}
    totals = dict(states=0, arcs=0, can_complete=0, cannot_complete=0)
    for i, c in enumerate(configs):
        g, indexes = graph(c)
        exact(d['Graphs'][i], g, f'Graphs[{i}]')
        exact(d['Groupings'][i], grouping(c), f'Groupings[{i}]')
        if c['Policy'] == 'stable-ready' or c['Layout'] != 'shared-1' or c['Base']['ServiceMask'] & 2:
            require(g['CannotComplete'] == 0, config_id(c) + ': positive completion control')
        if c['Policy'] == 'stable-ready':
            require(all(not q[j+2] or ready(c, q, j) for q in indexes for j in range(2)),
                    config_id(c) + ': stable holder readiness')
        reconstructed[config_id(c)] = (g, indexes)
        for key, field in [('states', 'States'), ('arcs', 'Arcs'), ('can_complete', 'CanComplete'),
                           ('cannot_complete', 'CannotComplete')]:
            totals[key] += g[field]
    ws = d['Witnesses']
    require(type(ws) is list and len(ws) == 10, 'Witnesses: roster')
    endpoints = []
    for i, name in enumerate(WITNESSES):
        require(type(ws[i]) is dict, f'Witnesses[{i}]: object')
        exact(ws[i].get('Name'), name, f'Witnesses[{i}].Name')
        policy = 'ready-now' if i in (3, 4, 5) else 'stable-ready' if i >= 8 else 'early'
        c = config('shared-1', policy, 0)
        g, indexes = reconstructed[config_id(c)]
        endpoints.append(replay(ws[i], c, g, indexes, f'Witnesses[{i}]'))
    witness_relations(ws, endpoints)
    p = d['Probes']
    numeric = 'OldEdgeStates Configurations Representatives MicroArcs Lifts LiftMicroSteps'
    keys(p, 'Domain Rows ' + numeric, 'Probes')
    exact(p['Domain'], 'factored-old-edge-v1', 'Probes.Domain')
    for field in numeric.split():
        integer(p[field], 'Probes.' + field)
    exact(p['Configurations'], 72, 'Probes.Configurations')
    require(type(p['Rows']) is list and len(p['Rows']) == 72, 'probe roster')
    fields = ('Representatives', 'MicroArcs', 'Lifts', 'LiftMicroSteps')
    for c, row in zip(configs, p['Rows']):
        keys(row, 'ConfigID ' + ' '.join(fields), 'probe row')
        exact(row['ConfigID'], config_id(c), 'probe row ConfigID')
        for field in fields:
            integer(row[field], 'probe row ' + field)
    for field in fields:
        require(p[field] == sum(r[field] for r in p['Rows']), 'probe total ' + field)
    return dict(schema=d['Schema'], independent_abstract_graphs=72, **totals,
                static_groupings=72, genuine_histories=10,
                replayed_microsteps=sum(len(w['Steps']) for w in ws),
                reported_operational_probes={k: p[k] for k in ('Domain',) + tuple(numeric.split())},
                operational_probes_independently_reenumerated=False,
                source_hashes_verified=len(d['Sources']))


def main():
    if len(sys.argv) != 2:
        print('usage: audit.py EXPORT[.gz]', file=sys.stderr)
        return 2
    try:
        summary = check(load(sys.argv[1]))
    except (InvalidEvidence, OSError, EOFError) as exc:
        print('REJECT: ' + str(exc), file=sys.stderr)
        return 1
    print(json.dumps(summary, sort_keys=True, indent=2))
    return 0


if __name__ == '__main__':
    sys.exit(main())
