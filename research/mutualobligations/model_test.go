package mutualobligations

import (
	"slices"
	"testing"

	prior "github.com/JakeFAU/alan-research/research/twoagent"
)

func TestSerialTrapAndAsymmetricCompletion(t *testing.T) {
	priorGraph := prior.Explore(prior.DurableACK)
	projection := newProjection(priorGraph)
	for _, p := range Profiles() {
		r, g := explore(p)
		if r.States < 2 || r.Edges < r.States || r.CyclicStates == 0 {
			t.Fatalf("%s: missing finite graph/cycles: %+v", p, r)
		}
		if r.LocalSafetyViolations+r.LocalTransitionViolations+r.CustodyViolations != 0 {
			t.Fatalf("%s: safety/custody violation", p)
		}
		if r.AppliedStates == 0 || r.ReleasedStates == 0 || r.CompleteStates == 0 {
			t.Fatalf("%s: custody exercise vacuous", p)
		}
		if p == SerialSerial {
			missingWait := r.ShortestTrap == nil || r.AdmittedWait == nil || r.RecoveredWait == nil
			if r.CannotComplete == 0 || r.TrapMismatch != 0 || missingWait {
				t.Fatalf("missing exact serial obstruction: %+v", r)
			}
			if len(r.ShortestTrap.Trace) != 2 || len(r.AdmittedWait.Trace) != 6 || len(r.RecoveredWait.Trace) != 10 {
				t.Fatal("nonminimal wait witnesses")
			}
		} else if r.CannotComplete != 0 {
			t.Fatalf("%s: asymmetric/symmetric service cannot complete", p)
		}
		checkJointGraph(t, p, g)
		projection.check(t, p, g)
		checkOracleNoninterference(t, p, g)
		checkWitnesses(t, p, r, g)
		t.Logf(
			"%s states=%d edges=%d cyclic=%d components=%d applied=%d released=%d complete=%d can=%d cannot=%d trap=%d",
			p, r.States, r.Edges, r.CyclicStates, r.CyclicComponents,
			r.AppliedStates, r.ReleasedStates, r.CompleteStates,
			r.CanComplete, r.CannotComplete, r.TrapStates,
		)
	}
}

// Each BFS predecessor supplies an independently checked ledger, then EVERY
// outgoing arc is checked from it and compared to its target. No quadratic
// rebuilding of every full prefix and no call to Successors inside the checker.
func checkJointGraph(t *testing.T, p Profile, g *graph) {
	t.Helper()
	ledgers := make([]ledger, len(g.nodes))
	ids := make(map[State]uint32, len(g.nodes))
	for i, n := range g.nodes {
		ids[n.state] = uint32(i)
		if n.parent >= 0 {
			l, err := ledgers[n.parent].step(p, n.via.event())
			if err != nil {
				t.Fatalf("%s parent of %d: %v", p, i, err)
			}
			ledgers[i] = l
		}
		if ledgerState(ledgers[i]) != n.state {
			t.Fatalf("%s independent predecessor state mismatch at %d", p, i)
		}
	}
	for i, n := range g.nodes {
		generated := Successors(p, n.state)
		if len(generated) != len(g.edges(i)) {
			t.Fatal("discarded reachable arcs")
		}
		for j, e := range g.edges(i) {
			l, err := ledgers[i].step(p, e.event.event())
			if err != nil {
				t.Fatalf("%s independent edge %d -> %d: %v", p, i, e.to, err)
			}
			if ledgerState(l) != g.nodes[e.to].state {
				t.Fatalf("%s independent target mismatch %d -> %d", p, i, e.to)
			}
			got := generated[j]
			to, ok := ids[got.State]
			if !ok || to != e.to || got.Event != e.event.event() {
				t.Fatal("graph not closed under every successor")
			}
			if traceDepth(g, int(e.to)) > traceDepth(g, i)+1 {
				t.Fatal("predecessor path is not shortest")
			}
		}
	}
}
func traceDepth(g *graph, i int) int {
	depth := 0
	for g.nodes[i].parent >= 0 {
		depth++
		i = int(g.nodes[i].parent)
	}
	return depth
}
func ledgerState(l ledger) State {
	var s State
	for i, p := range l.phase {
		s.Agents[i] = Availability{Down: p == 1, Crashed: p != 0}
		if p == 2 {
			s.Agents[i].Activation = 1
		}
	}
	for i, x := range l.entries {
		s.Obligations[i] = EdgeState{Committed: x.committed, Admitted: x.admitted, Applied: x.applied, Consumed: x.consumed,
			RequestChannel: x.request, ACKChannel: x.ack, WorldResult: x.actual, Retained: x.receipt,
			ResultChannel: x.channel, Volatile: x.received, Recorded: x.durable}
	}
	return s
}

type projectionArc struct {
	from, to uint32
	event    prior.Event
}
type projection struct {
	ids  map[prior.State]uint32
	arcs map[projectionArc]bool
}

func newProjection(r prior.Report) projection {
	p := projection{ids: map[prior.State]uint32{}, arcs: map[projectionArc]bool{}}
	for i, n := range r.Nodes {
		p.ids[n.State] = uint32(i)
	}
	for i, n := range r.Nodes {
		for _, e := range n.Edges {
			p.arcs[projectionArc{from: uint32(i), to: uint32(e.To), event: e.Event}] = true
		}
	}
	return p
}
func projected(s State, d Obligation) prior.State {
	e, a := s.Obligations[d], s.Agents[d]
	return prior.State{Delegation: prior.FixedDelegation(), Grant: prior.FixedGrant, Committed: e.Committed,
		ADown: a.Down, Crashed: a.Crashed, Activation: a.Activation,
		AVolatile: oldValue(e.Volatile), ARecorded: oldValue(e.Recorded),
		AConsumed: e.Consumed, BAdmitted: e.Admitted, Applied: e.Applied, BResult: oldValue(e.Retained),
		RequestChannel: e.RequestChannel, ResultChannel: oldValue(e.ResultChannel),
		ACKChannel: e.ACKChannel, WorldResult: oldValue(e.WorldResult),
	}
}
func oldValue(v Value) prior.Value {
	switch v {
	case R1:
		return prior.R1
	case R2:
		return prior.R2
	}
	return ""
}
func (projection projection) check(t *testing.T, p Profile, g *graph) {
	t.Helper()
	ids := make([][2]uint32, len(g.nodes))
	for i, n := range g.nodes {
		for d := AB; d <= BA; d++ {
			id, ok := projection.ids[projected(n.state, d)]
			if !ok {
				t.Fatalf("%s state %d edge %d not in frozen Phase3A", p, i, d)
			}
			ids[i][d] = id
		}
	}
	steps, stutters := 0, 0
	for i := range g.nodes {
		for _, e := range g.edges(i) {
			for d := AB; d <= BA; d++ {
				from, to := ids[i][d], ids[e.to][d]
				if from == to {
					stutters++
					continue
				}
				event := e.event.event()
				physical := event.Kind == Crash || event.Kind == Recover
				if (physical && event.Actor != Agent(d)) || (!physical && event.Obligation != d) {
					t.Fatal("nonlocal work changed projection")
				}
				old := prior.Event{Kind: prior.Kind(event.Kind), Delegation: prior.FixedDelegation(), Grant: prior.FixedGrant,
					Result: oldValue(event.Result), Duplicate: event.Duplicate}
				if !projection.arcs[projectionArc{from: from, to: to, event: old}] {
					t.Fatalf("%s transition %d -> %d edge %d not Phase3A or stutter", p, i, e.to, d)
				}
				steps++
			}
		}
	}
	t.Logf("%s frozen Phase3A projection: %d nonstutter + %d stutter checks", p, steps, stutters)
}
func checkOracleNoninterference(t *testing.T, p Profile, g *graph) {
	t.Helper()
	for _, n := range g.nodes {
		for d := AB; d <= BA; d++ {
			absent := n.state
			absent.Obligations[d].WorldResult = NoValue
			expected := Successors(p, absent)
			// An absent oracle can be malformed for an applied state. This
			// perturbation probes presence guards, not a third authentic result.
			for _, value := range []Value{R1, R2} {
				replaced := n.state
				replaced.Obligations[d].WorldResult = value
				got := Successors(p, replaced)
				if len(expected) != len(got) {
					t.Fatal("oracle presence/value changed enablement")
				}
				for i := range expected {
					if expected[i].Event != got[i].Event || Recoverable(expected[i].State) != Recoverable(got[i].State) {
						t.Fatal("oracle presence/value supplied agent-visible work/data")
					}
				}
			}
		}
	}
}
func checkWitnesses(t *testing.T, p Profile, r Report, g *graph) {
	t.Helper()
	witnesses := []*Witness{
		r.ShortestTrap, r.AdmittedWait, r.RecoveredWait, r.ShortestComplete,
		r.PostApplicationCustody, r.ReleasedCustody,
	}
	for _, w := range witnesses {
		if w == nil {
			continue
		}
		if err := CheckTrace(p, w.Trace); err != nil {
			t.Fatalf("invalid exported witness: %v", err)
		}
		if got := follow(t, p, w.Trace); got != w.State {
			t.Fatal("witness state mismatch")
		}
		if len(w.Trace) != traceDepth(g, w.StateID) {
			t.Fatal("witness not shortest")
		}
	}
	if r.UnfairRetry == nil || !r.UnfairRetry.CanComplete {
		t.Fatal("missing unfair but completable retry cycle")
	}
	w := r.UnfairRetry
	trace := append(slices.Clone(w.Prefix), w.Cycle...)
	if err := CheckTrace(p, trace); err != nil {
		t.Fatal(err)
	}
	if follow(t, p, trace) != w.State || Complete(w.State) {
		t.Fatal("retry cycle not repeatable before completion")
	}
	if r.AdmittedWait != nil && (!admittedDrainedUp(r.AdmittedWait.State) || !Trap(r.AdmittedWait.State)) {
		t.Fatal("wait witness contract")
	}
	if r.RecoveredWait != nil {
		a := r.RecoveredWait.State.Agents
		if a[A].Activation != 1 || a[B].Activation != 1 {
			t.Fatal("wait did not recover both agents")
		}
	}
}
func event(d Obligation, k Kind, v Value) Event {
	actor := Agent(d)
	switch k {
	case Admit, Deduplicate, Apply, SendResult, Release:
		actor = Agent(1 - d)
	}
	return Event{Kind: k, Actor: actor, Obligation: d, Binding: FixedBinding(d), Result: v}
}
func physical(a Agent, k Kind) Event { return Event{Kind: k, Actor: a, Obligation: NoObligation} }
func follow(t *testing.T, p Profile, trace []Event) State {
	t.Helper()
	s := Initial()
	for _, e := range trace {
		found := false
		for _, next := range Successors(p, s) {
			if next.Event == e {
				s = next.State
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("model omits %s by agent %d on edge %d", e.Kind, e.Actor, e.Obligation)
		}
	}
	return s
}

func TestIndependentAdversariesAndSharedAvailability(t *testing.T) {
	p := ServicingServicing
	for d := AB; d <= BA; d++ {
		ev := func(k Kind, v Value) Event { return event(d, k, v) }
		prefix := []Event{
			ev(Commit, NoValue), ev(SendRequest, NoValue), ev(Admit, NoValue),
			ev(Apply, R1), ev(SendResult, R1), ev(ReceiveResult, R1), ev(Record, R1),
		}
		ready := [][]Event{
			prefix[:1], prefix[:2], prefix[:3], prefix[:4], prefix[:5], prefix[:6], prefix[:7],
			append(slices.Clone(prefix), ev(Consume, R1)), append(slices.Clone(prefix), ev(SendACK, NoValue)),
			append(slices.Clone(prefix), ev(SendACK, NoValue), ev(Release, NoValue)),
			append(slices.Clone(prefix[:3]), ev(SendRequest, NoValue), ev(Deduplicate, NoValue)),
		}
		for _, trace := range ready {
			last := trace[len(trace)-1]
			t.Run(string(rune('A'+d))+"/down/"+string(last.Kind), func(t *testing.T) {
				if err := CheckTrace(p, trace); err != nil {
					t.Fatalf("invalid positive control: %v", err)
				}
				attack := append(slices.Clone(trace[:len(trace)-1]), physical(last.Actor, Crash), last)
				if err := CheckTrace(p, attack); err == nil {
					t.Fatal("down agent performed role work")
				}
				s := follow(t, p, attack[:len(attack)-1])
				for _, e := range Successors(p, s) {
					if e.Event.Actor == last.Actor && e.Event.Kind != Recover {
						t.Fatal("down agent has enabled role action")
					}
				}
			})
		}
		good := append(
			slices.Clone(prefix), ev(SendACK, NoValue), ev(Release, NoValue),
			physical(Agent(d), Crash), physical(Agent(1-d), Crash),
			physical(Agent(d), Recover), physical(Agent(1-d), Recover), ev(Consume, R1),
		)
		if err := CheckTrace(p, good); err != nil {
			t.Fatal(err)
		}
		s := follow(t, p, good)
		if s.Obligations[d].Retained != NoValue || s.Obligations[d].Recorded != R1 || !s.Obligations[d].Consumed {
			t.Fatal("released custody lost across both crashes")
		}
		attacks := map[string][]Event{
			"forgery":                append(slices.Clone(prefix[:5]), ev(ReceiveResult, R2)),
			"early-ack":              append(slices.Clone(prefix[:6]), ev(SendACK, NoValue)),
			"double-apply":           append(slices.Clone(prefix[:4]), ev(Apply, R1)),
			"early-consume":          append(slices.Clone(prefix[:6]), ev(Consume, R1)),
			"unreceived-record":      append(slices.Clone(prefix[:4]), ev(Record, R1)),
			"release-without-ack":    append(slices.Clone(prefix), ev(Release, NoValue)),
			"send-without-retention": {ev(SendResult, R1)},
			"resend-deleted-receipt": append(slices.Clone(good[:9]), ev(SendResult, R1)),
			"wrong-record-value":     append(slices.Clone(prefix[:6]), ev(Record, R2)),
			"ack-payload":            append(slices.Clone(prefix), ev(SendACK, R1)),
			"double-consume":         append(slices.Clone(good), ev(Consume, R1)),
		}
		for name, trace := range attacks {
			t.Run(string(rune('A'+d))+"/"+name, func(t *testing.T) {
				if CheckTrace(p, trace) == nil {
					t.Fatal("adversarial trace accepted")
				}
			})
		}
		changes := map[string]func(*Event){
			"actor":             func(e *Event) { e.Actor = Agent(1 - e.Actor) },
			"edge":              func(e *Event) { e.Obligation = 1 - e.Obligation },
			"generation":        func(e *Event) { e.Binding.CallerGeneration++ },
			"callee-generation": func(e *Event) { e.Binding.CalleeGeneration++ },
			"operation":         func(e *Event) { e.Binding.Operation++ },
			"endpoint":          func(e *Event) { e.Binding.Callee = e.Binding.Caller },
			"request":           func(e *Event) { e.Binding.Request = "forged" },
			"grant":             func(e *Event) { e.Binding.Grant = "forged" },
			"resource":          func(e *Event) { e.Binding.Resource = "peer-counter" },
			"fully-rebound-cross-edge": func(e *Event) {
				e.Obligation = 1 - e.Obligation
				e.Binding = FixedBinding(e.Obligation)
				e.Actor = Agent(e.Obligation)
			},
		}
		for name, change := range changes {
			t.Run(string(rune('A'+d))+"/binding/"+name, func(t *testing.T) {
				trace := slices.Clone(prefix)
				change(&trace[len(trace)-1])
				if CheckTrace(p, trace) == nil {
					t.Fatal("spoofed correlation accepted")
				}
			})
		}
	}
	if CheckTrace("bad", nil) == nil {
		t.Fatal("unknown profile accepted")
	}
	sameGrant := FixedBinding(AB).Grant == FixedBinding(BA).Grant
	sameResource := FixedBinding(AB).Resource == FixedBinding(BA).Resource
	if FixedBinding(AB) == FixedBinding(BA) || sameGrant || sameResource {
		t.Fatal("reciprocal authority collapsed")
	}
}

func TestSerialGateOnlyBlocksApplication(t *testing.T) {
	trace := []Event{
		event(AB, Commit, NoValue), event(BA, Commit, NoValue),
		event(AB, SendRequest, NoValue), event(AB, Admit, NoValue),
		event(BA, SendRequest, NoValue), event(BA, Admit, NoValue),
	}
	if err := CheckTrace(SerialSerial, trace); err != nil {
		t.Fatal(err)
	}
	for d := AB; d <= BA; d++ {
		attempt := append(slices.Clone(trace), event(d, Apply, R1))
		if CheckTrace(SerialSerial, attempt) == nil {
			t.Fatal("serial gate allowed application")
		}
		p := ServicingSerial
		if d == AB {
			p = SerialServicing
		}
		if err := CheckTrace(p, attempt); err != nil {
			t.Fatalf("own policy placement wrong: %v", err)
		}
	}
	// The incoming gate uses committed AND not consumed, not mere admission or
	// the peer's wait; request delivery is legal throughout both waits.
	for _, p := range Profiles() {
		follow(t, p, trace)
	}
}

func TestPreRecordCrashRetryAndCoincidentResults(t *testing.T) {
	for _, value := range []Value{R1, R2} {
		for d := AB; d <= BA; d++ {
			e := func(k Kind, v Value) Event { return event(d, k, v) }
			trace := []Event{
				e(Commit, NoValue), e(SendRequest, NoValue), e(Admit, NoValue),
				e(Apply, value), e(SendResult, value), e(ReceiveResult, value),
				physical(Agent(d), Crash), e(SendResult, value), physical(Agent(d), Recover),
				e(SendRequest, NoValue), e(Deduplicate, NoValue),
				e(ReceiveResult, value), e(Record, value), e(Consume, value),
			}
			if err := CheckTrace(SerialSerial, trace); err != nil {
				t.Fatal(err)
			}
			if !follow(t, SerialSerial, trace).Obligations[d].Consumed {
				t.Fatal("retry lost original result")
			}
		}
		// Both edges may choose identical opaque values and local sequence numbers;
		// independent receipt provenance still requires both actual causal chains.
		trace := []Event{}
		for d := AB; d <= BA; d++ {
			trace = append(trace, []Event{
				event(d, Commit, NoValue), event(d, SendRequest, NoValue),
				event(d, Admit, NoValue), event(d, Apply, value),
				event(d, SendResult, value), event(d, ReceiveResult, value),
				event(d, Record, value), event(d, Consume, value),
			}...)
		}
		if err := CheckTrace(SerialSerial, trace); err != nil {
			t.Fatal(err)
		}
		if !Complete(follow(t, SerialSerial, trace)) {
			t.Fatal("equal values merged original consumptions")
		}
	}
}

func TestConcurrentReadersHaveNoSharedScratch(t *testing.T) {
	trace := []Event{event(AB, Commit, NoValue), event(AB, SendRequest, NoValue), event(AB, Admit, NoValue)}
	shared := follow(t, SerialSerial, trace)
	for _, p := range Profiles() {
		t.Run(string(p), func(t *testing.T) {
			t.Parallel()
			expected := Successors(p, shared)
			for range 100 {
				got := Successors(p, shared)
				if !slices.Equal(got, expected) {
					t.Fatal("concurrent call changed shared state or successors")
				}
				if err := CheckTrace(p, trace); err != nil {
					t.Fatal(err)
				}
				// Caller ownership of returned slices must not mutate another call's output.
				if len(got) > 0 {
					got[0].State.Obligations[AB].Recorded = R2
				}
			}
		})
	}
}

func TestSafetyCheckersRejectLostFacts(t *testing.T) {
	s := State{Agents: [2]Availability{{Crashed: true, Activation: 1}, {Crashed: true, Activation: 1}}}
	s.Obligations[AB] = EdgeState{
		Committed: true, Admitted: true, Applied: true, Consumed: true,
		Recorded: R1, WorldResult: R1,
	}
	tests := map[string]func(*State){
		"commit":      func(n *State) { n.Obligations[AB].Committed = false },
		"admission":   func(n *State) { n.Obligations[AB].Admitted = false },
		"application": func(n *State) { n.Obligations[AB].Applied = false },
		"consumption": func(n *State) { n.Obligations[AB].Consumed = false },
		"record":      func(n *State) { n.Obligations[AB].Recorded = NoValue },
		"actual":      func(n *State) { n.Obligations[AB].WorldResult = R2 },
		"crash":       func(n *State) { n.Agents[A].Crashed = false },
		"activation":  func(n *State) { n.Agents[A].Activation = 0 },
	}
	for name, change := range tests {
		t.Run(name, func(t *testing.T) {
			n := s
			change(&n)
			if transitionSafe(s, Edge{Event: Event{Kind: Release, Obligation: AB, Actor: B}, State: n}) {
				t.Fatal("lost durable fact accepted")
			}
		})
	}
}
