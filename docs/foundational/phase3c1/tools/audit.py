#!/usr/bin/env python3
"""Independent stdlib audit of the frozen Phase 3C.1 evidence interface.

No producer imports, execution, or reducer calls. Graphs are reconstructed from
named progress registers; concrete histories use a causal custody ledger and a
declarative protocol table. See ../03-independent-audit.md for scope.
"""
import copy
import itertools
import json
from pathlib import Path
import sys


class AuditError(ValueError):
    """A malformed record or a falsified evidence assertion."""


def require(condition, reason):
    if not condition:
        raise AuditError(reason)


# The schema grammar is separate from semantic validation. Exact primitive types
# matter: Python's bool is an int subclass, but the frozen JSON interface differs.
RECORDS = {
    'Evidence': 'Schema:s Graphs:[Graph] Groupings:[Grouping] Witnesses:[Witness] Collisions:[Collision] Retry:Retry Claims:[Claim]',
    'Config': 'Topology:s ServiceMask:i Agents:[Agent] Originals:[Binding]',
    'Agent': 'ID:s Generation:i Servicing:b',
    'Binding': 'Original:s Caller:s CallerGeneration:i Operation:i Callee:s CalleeGeneration:i Request:s Grant:s Resource:s',
    'State': 'Stages:[s]',
    'Action': 'Kind:s Original:s Binding:Binding',
    'Arc': 'Action:Action To:i',
    'Node': 'State:State Arcs:[Arc] Complete:b CanComplete:b Parent:i Via:?Action',
    'Dependency': 'Original:s Requires:s',
    'Certificate': 'Dependencies:[Dependency] Acyclic:b Order:[s]',
    'Graph': 'Config:Config Nodes:[Node] States:i Arcs:i CanComplete:i CannotComplete:i Certificate:Certificate',
    'Component': 'Owners:[Agent] Ports:[Agent] Originals:[Binding]',
    'Grouping': 'ConfigID:s A:Component B:Component C:Component AB:Component BC:Component Left:Component Right:Component Flat:Config LeftClosed:Config RightClosed:Config Equal:b',
    'Availability': 'Down:b Crashed:b Activation:i',
    'Transport': 'Binding:Binding Committed:b Admitted:b Applied:b Consumed:b Volatile:s Recorded:s Retained:s RequestChannel:b ResultChannel:s ACKChannel:b WorldResult:s',
    'ConcreteState': 'Agents:[Availability] Originals:[Transport]',
    'Event': 'Kind:s Actor:s Original:s Binding:Binding Result:s Duplicate:b',
    'MicroStep': 'Event:Event State:ConcreteState',
    'Witness': 'Name:s Config:Config Initial:ConcreteState Steps:[MicroStep] Projection:[Action] Complete:b CanComplete:b',
    'Collision': 'Rule:s Left:s Right:s HiddenOriginal:s ExposedOriginals:[s] Observation:s',
    'Retry': 'Witness:s FromState:i ToState:i',
    'Claim': 'ID:s Status:s Statement:s Evidence:[s] Limits:[s]',
}


def shape(value, kind, path='$'):
    if kind.startswith('?'):
        if value is not None:
            shape(value, kind[1:], path)
    elif kind.startswith('['):
        require(type(value) is list, 'schema: {} must be array'.format(path))
        for i, entry in enumerate(value):
            shape(entry, kind[1:-1], '{}[{}]'.format(path, i))
    elif kind in ('s', 'i', 'b'):
        require(type(value) is {'s': str, 'i': int, 'b': bool}[kind],
                'schema: {} has wrong primitive type'.format(path))
    else:
        fields = dict(part.split(':') for part in RECORDS[kind].split())
        require(type(value) is dict and set(value) == set(fields),
                'schema: {} has missing/unknown fields'.format(path))
        for key, child in fields.items():
            shape(value[key], child, path + '.' + key)


NAMES = ('AB', 'BC', 'CA')
AGENTS = ('A', 'B', 'C')


def binding(index):
    caller, callee = AGENTS[index], AGENTS[(index + 1) % 3]
    resource = 'counter-' + callee
    return dict(Original=NAMES[index], Caller=caller, CallerGeneration=1,
                Operation=1, Callee=callee, CalleeGeneration=1,
                Request='increment', Grant='grant:{}/g1->{}/g1:{}'.format(caller, callee, resource),
                Resource=resource)


BINDINGS = tuple(binding(i) for i in range(3))
EMPTY_BINDING = {key: (0 if type(value) is int else '') for key, value in BINDINGS[0].items()}


def config(topology, mask):
    return dict(Topology=topology, ServiceMask=mask,
                Agents=[dict(ID=name, Generation=1, Servicing=bool(mask & (1 << i)))
                        for i, name in enumerate(AGENTS)],
                Originals=list(BINDINGS[:2 if topology == 'chain' else 3]))


CONFIGS = tuple(config(t, m) for t in ('chain', 'cycle') for m in range(8))


def action(index, kind):
    return dict(Kind=kind, Original=NAMES[index], Binding=BINDINGS[index])


def successors(conf, state):
    """Finite register semantics, independent of producer traversal/flags."""
    count = len(conf['Originals'])
    out = []
    for i in range(count):
        current = state[i]
        if current == 'D':
            continue
        callee = (i + 1) % 3
        waiting = callee < count and state[callee] in ('W', 'P')
        if current == 'W' and not conf['Agents'][callee]['Servicing'] and waiting:
            continue
        target = list(state)
        target[i] = {'U': 'W', 'W': 'P', 'P': 'D'}[current]
        out.append((action(i, {'U': 'commit', 'W': 'apply', 'P': 'consume'}[current]), tuple(target)))
    return out


def model(conf):
    """Construct the whole small lattice, then select reachable states and arcs."""
    count = len(conf['Originals'])
    universe = [s + (('',) if count == 2 else ()) for s in itertools.product('UWPD', repeat=count)]
    relation = {s: successors(conf, s) for s in universe}
    initial = tuple('U' for _ in range(count)) + (('',) if count == 2 else ())
    distance = {initial: 0}
    frontier = [initial]
    for source in frontier:
        for _, target in relation[source]:
            if target not in distance:
                distance[target] = distance[source] + 1
                frontier.append(target)
    graph = {s: relation[s] for s in distance}
    complete = tuple('D' for _ in range(count)) + (('',) if count == 2 else ())
    can = {complete}
    # Fixed-point predecessor closure, intentionally independent of exported flags.
    while True:
        more = {s for s in graph if any(t in can for _, t in graph[s])}
        if more <= can:
            break
        can |= more
    return graph, can, distance, complete


def certificate(conf):
    edges = []
    names = [b['Original'] for b in conf['Originals']]
    for i, name in enumerate(names):
        callee = (i + 1) % 3
        if callee < len(names) and not conf['Agents'][callee]['Servicing']:
            edges.append(dict(Original=name, Requires=names[callee]))
    order = []
    while True:
        before = len(order)
        for name in names:
            prerequisites = [e['Requires'] for e in edges if e['Original'] == name]
            if name not in order and all(p in order for p in prerequisites):
                order.append(name)
        if len(order) == before:
            break
    return dict(Dependencies=edges, Acyclic=len(order) == len(names), Order=order)


def audit_graph(record, conf):
    require(record['Config'] == conf, 'config: exact identity/policy roster mismatch')
    graph, can, distances, complete = model(conf)
    states = [tuple(n['State']['Stages']) for n in record['Nodes']]
    require(len(states) == len(set(states)) and set(states) == set(graph),
            'reachable state set: missing, duplicated or unreachable state')
    require(states[0] == min(distances, key=distances.get), 'reachable state set: wrong initial')
    indices = {s: i for i, s in enumerate(states)}
    for i, node in enumerate(record['Nodes']):
        source = states[i]
        for arc in node['Arcs']:
            a = arc['Action']
            require(a['Original'] in NAMES and a['Binding'] == BINDINGS[NAMES.index(a['Original'])],
                    'binding: graph action has forged identity')
        expected = [dict(Action=a, To=indices[t]) for a, t in graph[source]]
        require(node['Arcs'] == expected, 'transition relation: missing/extra/unguarded arc at {}'.format(source))
        require(node['Complete'] == (source == complete), 'completion: wrong Complete predicate')
        require(node['CanComplete'] == (source in can), 'progress verdict: wrong reverse reachability')
        if i == 0:
            require(node['Parent'] == -1 and node['Via'] is None, 'BFS predecessor: wrong root')
        else:
            p = node['Parent']
            require(0 <= p < i and distances[states[p]] + 1 == distances[source]
                    and (node['Via'], source) in graph[states[p]],
                    'BFS predecessor: invalid predecessor/action/depth')
        if i:
            require(distances[states[i - 1]] <= distances[source], 'BFS predecessor: non-BFS node order')
    counts = dict(States=len(graph), Arcs=sum(len(v) for v in graph.values()),
                  CanComplete=len(can), CannotComplete=len(graph) - len(can))
    require(all(record[k] == v for k, v in counts.items()), 'counts: graph totals differ')
    cert = certificate(conf)
    require(record['Certificate'] == cert, 'dependency: edges, order or acyclicity differs')
    require(not cert['Acyclic'] or len(can) == len(graph), 'dependency: acyclic sufficiency falsified')
    return graph, can, complete


def leaf(conf, owner):
    originals = [b for b in conf['Originals'] if b['Caller'] == owner]
    ports = {b['Callee'] for b in originals} - {owner}
    return dict(Owners=[a for a in conf['Agents'] if a['ID'] == owner],
                Ports=[a for a in conf['Agents'] if a['ID'] in ports], Originals=originals)


def join(left, right):
    """Join compatible open components without erasing endpoint incidence/debt."""
    owners = left['Owners'] + right['Owners']
    bindings = left['Originals'] + right['Originals']
    owner_ids = [a['ID'] for a in owners]
    require(len(set(owner_ids)) == len(owner_ids), 'grouping: duplicated ownership')
    require(len({b['Original'] for b in bindings}) == len(bindings), 'grouping: duplicate original')
    require(len({b['Resource'] for b in bindings}) == len(bindings), 'grouping: resource alias')
    require(len({b['Caller'] for b in bindings}) == len(bindings), 'grouping: duplicate outgoing original')
    endpoints = {}
    for a in owners + left['Ports'] + right['Ports']:
        require(a['ID'] not in endpoints or endpoints[a['ID']] == a, 'grouping: incompatible endpoint policy/generation')
        endpoints[a['ID']] = a
    for b in bindings:
        for role in ('Caller', 'Callee'):
            require(b[role] in endpoints and endpoints[b[role]]['Generation'] == b[role + 'Generation'],
                    'grouping: unresolved endpoint binding')
    return dict(Owners=sorted(owners, key=lambda a: a['ID']),
                Ports=[endpoints[k] for k in sorted(endpoints) if k not in owner_ids],
                Originals=sorted(bindings, key=lambda b: b['Original']))


def audit_grouping(record, conf, flat_graph):
    require(record['ConfigID'] == '{}/{}'.format(conf['Topology'], conf['ServiceMask']), 'grouping: config ID')
    expected = {a: leaf(conf, a) for a in AGENTS}
    expected['AB'] = join(expected['A'], expected['B'])
    expected['BC'] = join(expected['B'], expected['C'])
    expected['Left'] = join(expected['AB'], expected['C'])
    expected['Right'] = join(expected['A'], expected['BC'])
    require(all(record[k] == v for k, v in expected.items()), 'grouping: erased debt/ports or incompatible component')
    require(record['Flat'] == conf and record['Equal'], 'grouping: flat config/equality assertion')
    for name, closed in (('Left', 'LeftClosed'), ('Right', 'RightClosed')):
        c = expected[name]
        require(not c['Ports'], 'grouping: unresolved closed ports')
        reconstructed = dict(Topology=conf['Topology'], ServiceMask=conf['ServiceMask'],
                             Agents=c['Owners'], Originals=c['Originals'])
        require(record[closed] == reconstructed, 'grouping: closed config dropped debt or policy')
        require(model(reconstructed)[0] == flat_graph, 'grouping: full state/arc relation differs')


def initial_state(conf):
    blank = dict(Committed=False, Admitted=False, Applied=False, Consumed=False,
                 Volatile='', Recorded='', Retained='', RequestChannel=False,
                 ResultChannel='', ACKChannel=False, WorldResult='')
    return dict(Agents=[dict(Down=False, Crashed=False, Activation=0) for _ in AGENTS],
                Originals=[dict(blank, Binding=BINDINGS[i] if i < len(conf['Originals']) else EMPTY_BINDING)
                           for i in range(3)])


def progress(state):
    return tuple('' if not e['Binding']['Original'] else
                 'D' if e['Consumed'] else 'P' if e['Applied'] else
                 'W' if e['Committed'] else 'U' for e in state['Originals'])


def ready(conf, state, index):
    callee = (index + 1) % 3
    outgoing = state['Originals'][callee]
    waiting = outgoing['Committed'] and not outgoing['Consumed']
    edge = state['Originals'][index]
    return (not state['Agents'][callee]['Down'] and edge['Admitted'] and not edge['Applied']
            and (conf['Agents'][callee]['Servicing'] or not waiting))


# Protocol obligations are read as a table, not generated from producer code.
# +field means present/true; -field means absent/false. Receipt source is distinct
# from the historical world ledger, which never supplies actor bytes or readiness.
PROTOCOL = {
    'commit': ('caller', ('-Committed',), None, {'Committed': True}, None),
    'send-request': ('caller', ('+Committed', '-RequestChannel'), None, {'RequestChannel': True}, None),
    'admit': ('callee', ('+RequestChannel', '-Admitted'), None, {'Admitted': True}, 'RequestChannel'),
    'deduplicate': ('callee', ('+RequestChannel', '+Admitted'), None, {}, 'RequestChannel'),
    'apply': ('callee', ('+Admitted', '-Applied'), 'new', {'Applied': True}, None),
    'send-result': ('callee', ('+Retained', '-ResultChannel'), 'Retained', {}, None),
    'receive-result': ('caller', ('+ResultChannel',), 'ResultChannel', {}, 'ResultChannel'),
    'record': ('caller', ('+Volatile', '-Recorded'), 'Volatile', {}, None),
    'consume': ('caller', ('+Recorded', '-Consumed'), 'Recorded', {'Consumed': True}, None),
    'send-ack': ('caller', ('+Recorded', '-ACKChannel'), None, {'ACKChannel': True}, None),
    'release': ('callee', ('+ACKChannel',), None, {'Retained': ''}, 'ACKChannel'),
}
RESULT_TARGETS = {'apply': ('Retained', 'WorldResult'), 'send-result': ('ResultChannel',),
                  'receive-result': ('Volatile',), 'record': ('Recorded',)}


class Ledger:
    def __init__(self, conf):
        self.conf = conf
        self.state = initial_state(conf)
        self.applied = {}  # Original -> actual applied bytes, never an actor source.
        self.recorded = {}  # Original -> durable caller receipt acquired from channel.
        self.acked = set()  # Exact original identities with actual post-record ACK.

    def step(self, event):
        kind, actor = event['Kind'], event['Actor']
        require(actor in AGENTS, 'event: unknown physical actor')
        owner = AGENTS.index(actor)
        physical = self.state['Agents'][owner]
        if kind in ('crash', 'recover'):
            require(event['Original'] == '' and event['Binding'] == EMPTY_BINDING
                    and event['Result'] == '' and not event['Duplicate'], 'binding: malformed physical event')
            if kind == 'crash':
                require(not physical['Down'] and not physical['Crashed'], 'availability: repeated/invalid crash')
                physical.update(Down=True, Crashed=True)
                self.state['Originals'][owner]['Volatile'] = ''
            else:
                require(physical['Down'], 'availability: recovery without crash')
                physical.update(Down=False, Activation=1)
            return
        require(kind in PROTOCOL, 'event: unknown protocol event')
        name = event['Original']
        require(name in [b['Original'] for b in self.conf['Originals']], 'binding: undeclared original')
        index = NAMES.index(name)
        require(event['Binding'] == BINDINGS[index], 'binding: forged event identity')
        role, guards, source, patch, channel = PROTOCOL[kind]
        expected_actor = BINDINGS[index]['Caller' if role == 'caller' else 'Callee']
        require(actor == expected_actor, 'binding: actor does not own this protocol role')
        require(not physical['Down'], 'actor unavailable: physical owner down')
        require(not event['Duplicate'] or channel is not None, 'event: invalid duplicate delivery flag')
        edge = self.state['Originals'][index]
        if kind == 'apply':
            require(name not in self.applied, 'repeated application: original already applied')
            require(ready(self.conf, self.state, index), 'serial guard/readiness: apply disabled')
        if kind == 'release':
            require(name in self.acked and name in self.recorded and edge['ACKChannel'],
                    'release requires recorded ACK: no causal caller custody')
        require(all(bool(edge[g[1:]]) == (g[0] == '+') for g in guards),
                'protocol guard: {} prerequisites not satisfied'.format(kind))
        result = event['Result']
        if source == 'new':
            require(result in ('r1', 'r2'), 'receipt custody: invalid application result')
            self.applied[name] = result
        elif source:
            require(result == edge[source] and result == self.applied.get(name),
                    'receipt custody: {} did not use actual {}'.format(kind, source))
        else:
            require(result == '', 'receipt custody: payload on identity-only event')
        if kind == 'record':
            self.recorded[name] = result
        if kind == 'send-ack':
            require(name in self.recorded, 'receipt custody: ACK precedes recording')
            self.acked.add(name)
        edge.update(patch)
        for target in RESULT_TARGETS.get(kind, ()):
            edge[target] = result
        if channel and not event['Duplicate']:
            edge[channel] = '' if channel == 'ResultChannel' else False


def audit_witness(witness, conf):
    require(witness['Config'] == conf, 'config: wrong witness configuration')
    ledger = Ledger(conf)
    require(witness['Initial'] == ledger.state, 'initial state: invented durable/volatile/physical fact')
    graph, can, _, complete = model(conf)
    projection = []
    states = [copy.deepcopy(ledger.state)]
    for i, step in enumerate(witness['Steps']):
        before = progress(ledger.state)
        try:
            ledger.step(step['Event'])
        except AuditError as error:
            raise AuditError('{} step {}: {}'.format(witness['Name'], i, error)) from error
        require(step['State'] == ledger.state, 'post-state: {} step {} disagrees with causal ledger'.format(witness['Name'], i))
        after = progress(ledger.state)
        kind = step['Event']['Kind']
        if kind in ('commit', 'apply', 'consume'):
            a = action(NAMES.index(step['Event']['Original']), kind)
            require((a, after) in graph[before], 'projection: concrete event is not a named macro transition')
            projection.append(a)
        else:
            require(before == after, 'projection: transport/physical event changed progress')
        states.append(copy.deepcopy(ledger.state))
    require(witness['Projection'] == projection, 'projection: hidden or invented visible progress')
    final = progress(ledger.state)
    require(witness['Complete'] == (final == complete) and witness['CanComplete'] == (final in can),
            'witness progress: completion/reverse-reachability verdict differs')
    return states


WITNESS_CONFIGS = (
    ('R0-unpaid', 'chain', 7), ('R0-paid', 'chain', 7),
    ('R1-uncommitted', 'cycle', 6), ('R1-waiting', 'cycle', 6),
    ('serial-cycle-control', 'cycle', 0), ('retry-divergence', 'chain', 7),
    ('pre-record-loss-and-release', 'chain', 7),
)


def audit_collisions(data, histories):
    witnesses = {w['Name']: w for w in data['Witnesses']}
    expected = [('R0', 'R0-unpaid', 'R0-paid', ['BC']),
                ('R1', 'R1-uncommitted', 'R1-waiting', ['BC', 'CA'])]
    require(len(data['Collisions']) == 2, 'collision: wrong roster')
    for c, (rule, left, right, exposed) in zip(data['Collisions'], expected):
        require((c['Rule'], c['Left'], c['Right'], c['HiddenOriginal'], c['ExposedOriginals'])
                == (rule, left, right, 'AB', exposed) and bool(c['Observation']), 'collision: malformed identity/scope')
        l, r = histories[left][-1], histories[right][-1]
        require(l['Agents'] == r['Agents'], 'collision: physical availability differs')
        require(all(l['Originals'][NAMES.index(e)] == r['Originals'][NAMES.index(e)] for e in exposed),
                'collision: full exposed transport differs')
        require(witnesses[right]['Steps'][:len(witnesses[left]['Steps'])] == witnesses[left]['Steps'],
                'collision: right history is not a continuation')
        if rule == 'R0':
            require(progress(l) == ('W', 'D', '') and progress(r) == ('D', 'D', ''), 'collision: R0 debt/completion distinction missing')
            require(witnesses[left]['CanComplete'] and witnesses[right]['CanComplete'], 'collision: R0 must remain completable')
            added = witnesses[right]['Projection'][len(witnesses[left]['Projection']):]
            require(added == [action(0, 'apply'), action(0, 'consume')], 'collision: R0 lost internal visible application/consumption')
        else:
            require(progress(l)[0] == 'U' and progress(r)[0] == 'W'
                    and not l['Originals'][0]['Consumed'] and not r['Originals'][0]['Consumed'], 'collision: R1 hidden unpaid debt')
            conf = witnesses[left]['Config']
            require(ready(conf, l, 2) and not ready(conf, r, 2), 'collision: R1 concrete readiness distinction missing')
    control = histories['serial-cycle-control'][-1]
    require(progress(control) == ('W', 'W', 'W') and not witnesses['serial-cycle-control']['CanComplete'],
            'assumption classification: serial cycle must lack completion suffix')
    require(not model(config('cycle', 0))[0][progress(control)], 'assumption classification: serial cycle has progress')
    retry = data['Retry']
    require(retry == dict(Witness='retry-divergence', FromState=3, ToState=5), 'retry: wrong certificate endpoints')
    states = histories['retry-divergence']
    require(len(states) == 6 and states[3] == states[5], 'retry: concrete cycle does not close')
    events = witnesses['retry-divergence']['Steps'][3:5]
    require([s['Event']['Kind'] for s in events] == ['send-request', 'deduplicate'], 'retry: wrong divergence path')


    loss = witnesses['pre-record-loss-and-release']
    require([s['Event']['Kind'] for s in loss['Steps']] == [
        'commit', 'send-request', 'admit', 'apply', 'send-result',
        'receive-result', 'crash', 'recover', 'send-result', 'receive-result',
        'record', 'send-ack', 'release', 'consume'],
        'loss/release witness: missing loss, recovery, custody retransmission or settlement')
    trace = histories[loss['Name']]
    before_loss, after_loss = trace[6]['Originals'][0], trace[7]['Originals'][0]
    require(before_loss['Volatile'] in ('r1', 'r2') and before_loss['Recorded'] == ''
            and after_loss['Volatile'] == '' and after_loss['Retained'] == before_loss['Volatile']
            and progress(trace[-1]) == ('D', 'U', ''),
            'loss/release witness: caller loss/callee custody/final progress not witnessed')


def audit_claims(claims):
    expected = [('finite-path-equivalence', 'source-argument'), ('macro-graphs', 'bounded-exhaustive'),
                ('compatible-grouping', 'bounded-check-and-source-argument'), ('candidate-failures', 'concrete-witnesses')]
    require([(c['ID'], c['Status']) for c in claims] == expected, 'claim classification: unjustified evidence category/roster')
    for claim in claims:
        require(bool(claim['Statement']) and bool(claim['Evidence']) and bool(claim['Limits'])
                and all(claim['Evidence']) and all(claim['Limits']), 'claim classification: missing statement/evidence/limits')


def audit(data):
    shape(data, 'Evidence')
    require(data['Schema'] == 'alan.phase3c1.progress.v1', 'schema: unsupported version')
    require(len(data['Graphs']) == 16 and len(data['Groupings']) == 16, 'config roster: expected all 16 ordered profiles')
    states = arcs = cannot = 0
    for g, grouping, conf in zip(data['Graphs'], data['Groupings'], CONFIGS):
        graph, can, _ = audit_graph(g, conf)
        audit_grouping(grouping, conf, graph)
        states += len(graph)
        arcs += sum(len(v) for v in graph.values())
        cannot += len(graph) - len(can)
    require([w['Name'] for w in data['Witnesses']] == [w[0] for w in WITNESS_CONFIGS], 'witness roster: wrong names/order')
    histories = {}
    for w, (_, topology, mask) in zip(data['Witnesses'], WITNESS_CONFIGS):
        histories[w['Name']] = audit_witness(w, config(topology, mask))
    audit_collisions(data, histories)
    audit_claims(data['Claims'])
    return dict(status='accepted', graphs=16, grouped_graphs=32, states=states, arcs=arcs,
                cannot_complete=cannot, witnesses=len(histories),
                witness_steps=sum(len(v) - 1 for v in histories.values()),
                collisions=2, retry_cycles=1,
                serial_cycle='undischarged service assumptions; custody and identity preserved',
                scope='complete abstract graphs and selected concrete replay; no full concrete product enumeration',
                prose_claims='classification checked; statements and limits require source review')


def unique_object(pairs):
    result = {}
    for key, value in pairs:
        require(key not in result, 'schema: duplicate JSON key ' + key)
        result[key] = value
    return result


def main():
    try:
        require(len(sys.argv) == 2, 'usage: audit.py INPUT.json')
        data = json.loads(Path(sys.argv[1]).read_text(), object_pairs_hook=unique_object)
        result = audit(data)
        code = 0
    except (AuditError, ValueError, OSError) as error:
        result = dict(status='rejected', reason=str(error))
        code = 1
    print(json.dumps(result, sort_keys=True, separators=(',', ':')))
    return code


if __name__ == '__main__':
    sys.exit(main())
