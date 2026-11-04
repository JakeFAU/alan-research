package authorityservice

import (
	base "github.com/JakeFAU/alan-research/research/mutualobligations"
	"slices"
)

// Compact graph arcs retain every transition, including self-loops. Event
// descriptors are reconstructed from their source state and fixed configuration.
type eventCode uint16

var kinds = [...]string{
	"", "commit", "send-request", "admit", "deduplicate", "apply", "send-result",
	"receive-result", "record", "consume", "send-ack", "release", "crash", "recover",
	"authorize", "dispatch-admit", "dispatch-refuse", "retain-grant", "revoke",
}

func encode(e Event) eventCode {
	var k uint16
	for i, v := range kinds {
		if v == e.Kind {
			k = uint16(i)
			break
		}
	}
	code := k | uint16(e.Actor)<<5 | uint16(e.Obligation)<<6 | uint16(e.Result)<<8
	if e.Duplicate {
		code |= 1 << 10
	}
	return eventCode(code)
}
func (code eventCode) event(c Config, s State) Event {
	e := Event{
		Kind: kinds[uint16(code)&31], Actor: base.Agent((code >> 5) & 1),
		Obligation: base.Obligation((code >> 6) & 3), Result: base.Value((code >> 8) & 3),
		Duplicate: code&(1<<10) != 0,
	}
	return describe(c, s, e)
}

type arc struct {
	to    uint32
	event eventCode
}
type node struct {
	state  State
	start  uint32
	count  uint8
	parent int32
	via    eventCode
	raw    bool
	safe   bool
}
type graph struct {
	config Config
	nodes  []node
	arcs   []arc
}

func (g *graph) edges(i int) []arc {
	n := g.nodes[i]
	return g.arcs[int(n.start) : int(n.start)+int(n.count)]
}
func (g *graph) trace(i int) []Event {
	trace := []Event{}
	for g.nodes[i].parent >= 0 {
		n := g.nodes[i]
		p := int(n.parent)
		trace = append(trace, n.via.event(g.config, g.nodes[p].state))
		i = p
	}
	slices.Reverse(trace)
	return trace
}

// Witness traces are shortest BFS prefixes, or a shortest prefix plus the named
// event. CanComplete and CanCompleteSafely are existential suffix properties.
type Witness struct {
	StateID           int
	State             State
	Trace             []Event
	CanComplete       bool
	CanCompleteSafely bool
}
type CycleWitness struct {
	Prefix []Event
	Cycle  []Event
	State  State
}
type Report struct {
	Config                         Config
	States                         int
	Arcs                           int
	CyclicStates                   int
	CyclicComponents               int
	OperationalInvariantViolations int
	CustodyViolations              int
	SafeStates                     int
	SafeTransitions                int
	UnsafeStates                   int
	ViolatingTransitions           int
	ConsentTransitions             int
	TupleInventionTransitions      int
	AppliedStates                  int
	ReleasedStates                 int
	CompleteStates                 int
	SafeCompleteStates             int
	CanComplete                    int
	CannotComplete                 int
	CanCompleteSafely              int
	CannotCompleteSafely           int
	ShortestComplete               *Witness
	ShortestSafeComplete           *Witness
	ShortestUnsafe                 *Witness
	ShortestCannotComplete         *Witness
	Cuts                           map[string]*Witness
	UnfairRetry                    *CycleWitness
}

// Explore enumerates the complete finite graph; it does not impose a depth or
// state cutoff. Counts concern states/arcs, never the infinite history language.
func Explore(c Config) Report { r, _ := explore(c); return r }
func explore(c Config) (Report, *graph) {
	initial := Initial(c)
	g := &graph{config: c, nodes: []node{{state: initial, parent: -1}}, arcs: []arc{}}
	ids := map[State]uint32{initial: 0}
	type cut struct {
		from int
		to   int
		via  eventCode
	}
	cuts := map[string]cut{}
	r := Report{Config: c, Cuts: map[string]*Witness{}}
	for i := 0; i < len(g.nodes); i++ {
		s := g.nodes[i].state
		next := Successors(c, s)
		g.nodes[i].start = uint32(len(g.arcs))
		g.nodes[i].count = uint8(len(next))
		for _, tr := range next {
			if !transitionSafe(s, tr) {
				r.OperationalInvariantViolations++
			}
			to, ok := ids[tr.State]
			if !ok {
				to = uint32(len(g.nodes))
				ids[tr.State] = to
				g.nodes = append(g.nodes, node{state: tr.State, parent: int32(i), via: encode(tr.Event)})
			}
			code := encode(tr.Event)
			if code.event(c, s) != tr.Event {
				r.OperationalInvariantViolations++
			}
			g.arcs = append(g.arcs, arc{to: to, event: code})
			e := tr.Event
			if s.Unsafe == 0 && tr.State.Unsafe == 0 {
				r.SafeTransitions++
			}
			if e.Violations != 0 {
				r.ViolatingTransitions++
			}
			if e.Violations&ConsentViolation != 0 {
				r.ConsentTransitions++
			}
			if e.Violations&TupleInvention != 0 {
				r.TupleInventionTransitions++
			}
			remember := func(name string) {
				if _, ok := cuts[name]; !ok {
					cuts[name] = cut{from: i, to: int(to), via: code}
				}
			}
			if e.HasGrant {
				if e.Grant.Accepted {
					remember("grant-retained")
				} else {
					remember("grant-issuance-rejected")
				}
			}
			if e.Kind == "authorize" {
				if e.Decision == Refused {
					remember("authorization-refused")
				} else {
					remember("authorization-allowed")
				}
			}
			if e.Kind == "dispatch-admit" {
				remember("dispatch-admitted")
			}
			if e.Kind == "dispatch-refuse" {
				remember("fresh-admission-refused")
			}
			if e.Violations != 0 {
				remember("unsafe-" + e.Kind)
			}
			if e.Kind == "retain-grant" && s.Authority[base.BA].Decision == Refused {
				remember("late-grant-sticky-refusal")
			}
			if e.Kind == "revoke" {
				a := s.Authority[base.BA]
				b := s.Base.Obligations[base.BA]
				if a.Decision == Undecided {
					remember("revocation-before-authorization")
				}
				if a.Decision == Allowed && !a.Permit && !b.Applied {
					remember("revocation-between-authorization-admission")
				}
				if a.Permit {
					remember("revocation-after-admission")
				}
				if b.Applied {
					remember("revocation-after-application")
				}
			}
			if e.Kind == "apply" {
				remember("application-retained")
				if s.Revoked {
					remember("application-after-revocation")
				}
			}
			if e.Kind == "release" {
				remember("receipt-released")
			}
			if e.Kind == "crash" {
				for _, b := range s.Base.Obligations {
					if b.Applied {
						remember("applied-then-crashed")
					}
					if b.Applied && b.Retained == base.NoValue {
						remember("released-then-crashed")
					}
				}
				if s.Authority[base.Obligation(1-e.Actor)].Permit {
					remember("admitted-then-server-crashed")
				}
			}
			if e.Kind == "dispatch-admit" && e.Activation == 1 {
				remember("recovered-fresh-admission")
			}
		}
	}
	r.States = len(g.nodes)
	r.Arcs = len(g.arcs)
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
	reverse := func(safe bool) int {
		queue := []uint32{}
		for i := range g.nodes {
			n := &g.nodes[i]
			if Complete(n.state) && (!safe || n.state.Unsafe == 0) {
				if safe {
					n.safe = true
				} else {
					n.raw = true
				}
				queue = append(queue, uint32(i))
			}
		}
		for j := 0; j < len(queue); j++ {
			to := queue[j]
			for _, from := range incoming[offsets[to]:offsets[to+1]] {
				n := &g.nodes[from]
				if safe {
					if n.safe || n.state.Unsafe != 0 {
						continue
					}
					n.safe = true
				} else {
					if n.raw {
						continue
					}
					n.raw = true
				}
				queue = append(queue, from)
			}
		}
		return len(queue)
	}
	r.CanComplete = reverse(false)
	r.CanCompleteSafely = reverse(true)
	r.CannotComplete = r.States - r.CanComplete
	r.CannotCompleteSafely = r.States - r.CanCompleteSafely
	witness := func(i int, trace []Event) *Witness {
		n := g.nodes[i]
		return &Witness{StateID: i, State: n.state, Trace: trace, CanComplete: n.raw, CanCompleteSafely: n.safe}
	}
	for name, cut := range cuts {
		trace := append(g.trace(cut.from), cut.via.event(c, g.nodes[cut.from].state))
		r.Cuts[name] = witness(cut.to, trace)
	}
	for i, n := range g.nodes {
		s := n.state
		if !locallySafe(s) {
			r.OperationalInvariantViolations++
		}
		applied, released := false, false
		for _, b := range s.Base.Obligations {
			applied = applied || b.Applied
			released = released || (b.Applied && b.Retained == base.NoValue)
			if b.Applied && b.Retained == base.NoValue && b.Recorded == base.NoValue {
				r.CustodyViolations++
			}
		}
		if applied {
			r.AppliedStates++
		}
		if released {
			r.ReleasedStates++
		}
		if s.Unsafe != 0 {
			r.UnsafeStates++
			if r.ShortestUnsafe == nil {
				r.ShortestUnsafe = witness(i, g.trace(i))
			}
		}
		if Complete(s) {
			r.CompleteStates++
			if r.ShortestComplete == nil {
				r.ShortestComplete = witness(i, g.trace(i))
			}
			if s.Unsafe == 0 {
				r.SafeCompleteStates++
				if r.ShortestSafeComplete == nil {
					r.ShortestSafeComplete = witness(i, g.trace(i))
				}
			}
		}
		if !n.raw && r.ShortestCannotComplete == nil {
			r.ShortestCannotComplete = witness(i, g.trace(i))
		}
		if r.UnfairRetry == nil && n.raw && !Complete(s) {
			for _, a := range g.edges(i) {
				e := a.event.event(c, s)
				if e.Kind != "send-request" {
					continue
				}
				for _, b := range g.edges(int(a.to)) {
					f := b.event.event(c, g.nodes[a.to].state)
					if b.to == uint32(i) && f.Kind == "deduplicate" && e.Obligation == f.Obligation {
						r.UnfairRetry = &CycleWitness{Prefix: g.trace(i), Cycle: []Event{e, f}, State: s}
						break
					}
				}
				if r.UnfairRetry != nil {
					break
				}
			}
		}
	}
	r.SafeStates = r.States - r.UnsafeStates
	r.CyclicStates, r.CyclicComponents = cycles(g)
	return r, g
}
func locallySafe(s State) bool {
	for a, v := range s.Base.Agents {
		if v.Down && (!v.Crashed || s.Base.Obligations[a].Volatile != base.NoValue) {
			return false
		}
		if v.Activation > 1 || (v.Activation == 1 && (!v.Crashed || v.Down)) {
			return false
		}
		if v.Crashed && !v.Down && v.Activation != 1 {
			return false
		}
	}
	for d, b := range s.Base.Obligations {
		a := s.Authority[d]
		server := s.Base.Agents[1-d]
		if a.Permit && (server.Down || a.PermitActivation != server.Activation || a.Decision != Allowed || b.Applied) {
			return false
		}
		if b.Applied && a.Decision != Allowed {
			return false
		}
		if b.RequestChannel && !b.Committed || b.ACKChannel && b.Recorded == base.NoValue {
			return false
		}
		if b.Applied != (b.WorldResult != base.NoValue) || b.Applied && !b.Admitted ||
			b.Admitted && !b.Committed || b.Consumed && b.Recorded == base.NoValue {
			return false
		}
		for _, v := range []base.Value{b.Retained, b.Volatile, b.Recorded, b.ResultChannel} {
			if v != base.NoValue && (!b.Applied || v != b.WorldResult) {
				return false
			}
		}
	}
	return true
}
func transitionSafe(s State, tr Transition) bool {
	n := tr.State
	e := tr.Event
	if s.Unsafe&n.Unsafe != s.Unsafe || s.Revoked && !n.Revoked || s.Grant != NoGrant && s.Grant != n.Grant {
		return false
	}
	for d, a := range s.Authority {
		b := n.Authority[d]
		if a.Decision != Undecided && (a.Decision != b.Decision || a.Path != b.Path || a.Refusal != b.Refusal) {
			return false
		}
	}
	if e.Kind == "apply" && (!s.Authority[e.Obligation].Permit || n.Authority[e.Obligation].Permit) {
		return false
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
