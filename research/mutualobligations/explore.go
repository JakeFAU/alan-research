package mutualobligations

import "slices"

// Internal graph arcs contain compact event codes, never unbounded histories.
type eventCode uint16

var kinds = [...]Kind{
	"", Commit, SendRequest, Admit, Deduplicate, Apply, SendResult,
	ReceiveResult, Record, Consume, SendACK, Release, Crash, Recover,
}

func encode(e Event) eventCode {
	var k uint16
	for i, v := range kinds {
		if v == e.Kind {
			k = uint16(i)
			break
		}
	}
	code := k | uint16(e.Actor)<<4 | uint16(e.Obligation)<<5 | uint16(e.Result)<<7
	if e.Duplicate {
		code |= 1 << 9
	}
	return eventCode(code)
}
func (c eventCode) event() Event {
	e := Event{
		Kind: kinds[uint16(c)&15], Actor: Agent((c >> 4) & 1),
		Obligation: Obligation((c >> 5) & 3), Result: Value((c >> 7) & 3),
		Duplicate: c&(1<<9) != 0,
	}
	if e.Obligation != NoObligation {
		e.Binding = FixedBinding(e.Obligation)
	}
	return e
}

type arc struct {
	to    uint32
	event eventCode
}
type node struct {
	state       State
	start       uint32
	count       uint8
	parent      int32
	via         eventCode
	canComplete bool
}
type graph struct {
	nodes []node
	arcs  []arc
}

func (g *graph) edges(i int) []arc {
	n := g.nodes[i]
	return g.arcs[int(n.start) : int(n.start)+int(n.count)]
}
func (g *graph) trace(i int) []Event {
	trace := []Event{}
	for g.nodes[i].parent >= 0 {
		trace = append(trace, g.nodes[i].via.event())
		i = int(g.nodes[i].parent)
	}
	slices.Reverse(trace)
	return trace
}

type Witness struct {
	StateID     int
	State       State
	Trace       []Event
	CanComplete bool
}
type CycleWitness struct {
	Prefix      []Event
	Cycle       []Event
	State       State
	CanComplete bool
}

// Report exports counts and finite shortest witnesses, not the full graph.
// CanComplete is existential reachability of both consumptions, not fairness.
type Report struct {
	Profile                   Profile
	States                    int
	Edges                     int
	CyclicStates              int
	CyclicComponents          int
	LocalSafetyViolations     int
	LocalTransitionViolations int
	CustodyViolations         int
	AppliedStates             int
	ReleasedStates            int
	CompleteStates            int
	CanComplete               int
	CannotComplete            int
	TrapStates                int
	TrapMismatch              int
	ShortestTrap              *Witness
	AdmittedWait              *Witness
	RecoveredWait             *Witness
	ShortestComplete          *Witness
	PostApplicationCustody    *Witness
	ReleasedCustody           *Witness
	UnfairRetry               *CycleWitness
}

// Explore enumerates the complete reachable finite graph with no depth cutoff.
// BFS finds shortest witnesses; reverse reachability keeps every cyclic arc.
func Explore(p Profile) Report { r, _ := explore(p); return r }
func explore(p Profile) (Report, *graph) {
	r := Report{Profile: p}
	g := &graph{nodes: []node{{state: Initial(), parent: -1}}, arcs: []arc{}}
	ids := map[State]uint32{Initial(): 0}
	for i := 0; i < len(g.nodes); i++ {
		state := g.nodes[i].state
		successors := Successors(p, state)
		g.nodes[i].start = uint32(len(g.arcs))
		g.nodes[i].count = uint8(len(successors))
		for _, e := range successors {
			if !transitionSafe(state, e) {
				r.LocalTransitionViolations++
			}
			to, ok := ids[e.State]
			if !ok {
				to = uint32(len(g.nodes))
				ids[e.State] = to
				g.nodes = append(g.nodes, node{state: e.State, parent: int32(i), via: encode(e.Event)})
			}
			g.arcs = append(g.arcs, arc{to: to, event: encode(e.Event)})
		}
	}
	r.States = len(g.nodes)
	r.Edges = len(g.arcs)
	// Reverse CSR: each original arc is retained, including self-loops.
	offsets := make([]uint32, len(g.nodes)+1)
	for _, e := range g.arcs {
		offsets[int(e.to)+1]++
	}
	for i := 1; i < len(offsets); i++ {
		offsets[i] += offsets[i-1]
	}
	cursor := slices.Clone(offsets)
	incoming := make([]uint32, len(g.arcs))
	for i := range g.nodes {
		for _, e := range g.edges(i) {
			incoming[cursor[e.to]] = uint32(i)
			cursor[e.to]++
		}
	}
	queue := []uint32{}
	for i := range g.nodes {
		if Complete(g.nodes[i].state) {
			g.nodes[i].canComplete = true
			queue = append(queue, uint32(i))
			r.CompleteStates++
		}
	}
	for i := 0; i < len(queue); i++ {
		to := queue[i]
		for _, from := range incoming[offsets[to]:offsets[to+1]] {
			if !g.nodes[from].canComplete {
				g.nodes[from].canComplete = true
				queue = append(queue, from)
			}
		}
	}
	r.CanComplete = len(queue)
	r.CannotComplete = r.States - r.CanComplete
	witness := func(i int) *Witness {
		n := g.nodes[i]
		return &Witness{StateID: i, State: n.state, Trace: g.trace(i), CanComplete: n.canComplete}
	}
	for i, n := range g.nodes {
		s := n.state
		if !locallySafe(s) {
			r.LocalSafetyViolations++
		}
		applied, released := false, false
		for _, e := range s.Obligations {
			applied = applied || e.Applied
			released = released || (e.Applied && e.Retained == NoValue)
			if e.Applied && e.Retained == NoValue && e.Recorded == NoValue {
				r.CustodyViolations++
			}
		}
		if applied {
			r.AppliedStates++
			if r.PostApplicationCustody == nil {
				r.PostApplicationCustody = witness(i)
			}
		}
		if released {
			r.ReleasedStates++
			if r.ReleasedCustody == nil {
				r.ReleasedCustody = witness(i)
			}
		}
		if Complete(s) && r.ShortestComplete == nil {
			r.ShortestComplete = witness(i)
		}
		if Trap(s) {
			r.TrapStates++
			if r.ShortestTrap == nil {
				r.ShortestTrap = witness(i)
			}
		}
		if p == SerialSerial && Trap(s) == n.canComplete {
			r.TrapMismatch++
		}
		if Trap(s) && admittedDrainedUp(s) {
			if r.AdmittedWait == nil {
				r.AdmittedWait = witness(i)
			}
			if s.Agents[A].Crashed && s.Agents[B].Crashed && r.RecoveredWait == nil {
				r.RecoveredWait = witness(i)
			}
		}
		if r.UnfairRetry == nil && n.canComplete && !Complete(s) {
			for _, first := range g.edges(i) {
				e := first.event.event()
				if e.Kind != SendRequest {
					continue
				}
				for _, second := range g.edges(int(first.to)) {
					f := second.event.event()
					if second.to == uint32(i) && f.Kind == Deduplicate && f.Obligation == e.Obligation {
						r.UnfairRetry = &CycleWitness{Prefix: g.trace(i), Cycle: []Event{e, f}, State: s, CanComplete: true}
						break
					}
				}
				if r.UnfairRetry != nil {
					break
				}
			}
		}
	}
	r.CyclicStates, r.CyclicComponents = cycles(g)
	return r, g
}
func admittedDrainedUp(s State) bool {
	if s.Agents[A].Down || s.Agents[B].Down {
		return false
	}
	for _, e := range s.Obligations {
		if !e.Admitted || e.RequestChannel || e.ResultChannel != NoValue || e.ACKChannel {
			return false
		}
	}
	return true
}
func locallySafe(s State) bool {
	for a, v := range s.Agents {
		if v.Down && (!v.Crashed || s.Obligations[a].Volatile != NoValue) {
			return false
		}
		if v.Activation > 1 || (v.Activation == 1 && (!v.Crashed || v.Down)) {
			return false
		}
		if v.Crashed && !v.Down && v.Activation != 1 {
			return false
		}
	}
	for _, e := range s.Obligations {
		if (e.RequestChannel && !e.Committed) || (e.ACKChannel && e.Recorded == NoValue) {
			return false
		}
		if e.Applied != (e.WorldResult != NoValue) || (e.Applied && !e.Admitted) || (e.Admitted && !e.Committed) {
			return false
		}
		if e.WorldResult > R2 || (e.Consumed && e.Recorded == NoValue) {
			return false
		}
		for _, v := range []Value{e.Retained, e.Volatile, e.Recorded, e.ResultChannel} {
			if v != NoValue && (v != e.WorldResult || !e.Applied) {
				return false
			}
		}
	}
	return true
}
func transitionSafe(s State, e Edge) bool {
	n := e.State
	for i, a := range s.Agents {
		b := n.Agents[i]
		if a.Crashed && !b.Crashed {
			return false
		}
		if e.Event.Kind == Recover && e.Event.Actor == Agent(i) {
			if !a.Down || b.Down || b.Activation != 1 {
				return false
			}
		} else if a.Activation != b.Activation {
			return false
		}
	}
	for i, a := range s.Obligations {
		b := n.Obligations[i]
		lostCallerFact := (a.Committed && !b.Committed) || (a.Consumed && !b.Consumed)
		lostCalleeFact := (a.Admitted && !b.Admitted) || (a.Applied && !b.Applied)
		if lostCallerFact || lostCalleeFact {
			return false
		}
		if a.Recorded != NoValue && a.Recorded != b.Recorded {
			return false
		}
		if a.WorldResult != NoValue && a.WorldResult != b.WorldResult {
			return false
		}
		if e.Event.Kind == Apply && e.Event.Obligation == Obligation(i) {
			if a.Applied || !b.Applied || b.Retained != e.Event.Result {
				return false
			}
		} else if a.Applied != b.Applied || a.WorldResult != b.WorldResult {
			return false
		}
		if e.Event.Kind == Consume && e.Event.Obligation == Obligation(i) {
			if a.Consumed || !b.Consumed || b.Recorded != e.Event.Result {
				return false
			}
		} else if a.Consumed != b.Consumed {
			return false
		}
	}
	return true
}

// Tarjan SCC counts actual cycles, including duplicate-delivery self-loops.
func cycles(g *graph) (int, int) {
	index := make([]int, len(g.nodes))
	low := make([]int, len(g.nodes))
	on := make([]bool, len(g.nodes))
	for i := range index {
		index[i] = -1
	}
	stack := []int{}
	next, states, components := 0, 0, 0
	var visit func(int)
	visit = func(v int) {
		index[v] = next
		low[v] = next
		next++
		stack = append(stack, v)
		on[v] = true
		for _, e := range g.edges(v) {
			w := int(e.to)
			if index[w] < 0 {
				visit(w)
				low[v] = min(low[v], low[w])
			} else if on[w] {
				low[v] = min(low[v], index[w])
			}
		}
		if low[v] != index[v] {
			return
		}
		size := 0
		for {
			w := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			on[w] = false
			size++
			if w == v {
				break
			}
		}
		cyclic := size > 1
		if !cyclic {
			for _, e := range g.edges(v) {
				if int(e.to) == v {
					cyclic = true
					break
				}
			}
		}
		if cyclic {
			states += size
			components++
		}
	}
	for i := range g.nodes {
		if index[i] < 0 {
			visit(i)
		}
	}
	return states, components
}
