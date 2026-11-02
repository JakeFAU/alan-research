package substitution

import (
	"slices"
	"testing"
)

type peerBudget struct{ invocations, runs int } // E = 0 throughout this effect-free fragment
type peerLabel struct{ kind, operation, value string }
type peerStep struct {
	label peerLabel
	to    int
}
type peerNode struct {
	remaining       peerBudget
	admission, wait string
	steps           []peerStep
}
type peerGraph []peerNode

// Independent admission/selection model. The root stays open, depth/bytes are
// ample and authority is fixed. Operation names stand for distinct original
// handles. An accepted pair is serviced to its common successful empty reply.
// A refusal is selected from the original decision, never from today's budget.
func (b peerBudget) decide(members int) (peerBudget, string) {
	switch {
	case b.invocations == 0:
		return b, "invocation_limit"
	case b.runs < members:
		return b, "created_run_limit"
	default:
		return peerBudget{b.invocations - 1, b.runs - members}, "accepted"
	}
}

const (
	decideFirst = iota
	selectFirst
	decideSecond
	selectSecond
	branchOnReasons
	emitMarker
	candidateDone
)

type peerCut struct {
	remaining peerBudget
	pc        int
	original  [2]string
	selected  [2]string
	peer      string
	marker    string
}

type peerClient struct {
	refusal    func(peerBudget, string) string
	markerPair [2]string
}

var literalClient = peerClient{
	refusal:    func(_ peerBudget, reason string) string { return reason },
	markerPair: [2]string{"created_run_limit", "invocation_limit"},
}

// Enumerate X and Y separately. The product offers one static peer's original
// singleton at every interleaving of the displayed cuts. The open model instead
// offers an at-most-once successful (1,1) environment debit at every funded cut.
// Admission and matched outcome selection are diagnostic output labels. Only
// the marker port participates in the mailbox/continuation readiness projection.
func (client peerClient) graph(b peerBudget, markerBranch bool, context string) peerGraph {
	handles := [2]string{"c1", "c2"}
	start := peerCut{remaining: b}
	cuts := []peerCut{start}
	ids := map[peerCut]int{start: 0}
	var g peerGraph
	for i := 0; i < len(cuts); i++ {
		s := cuts[i]
		// The only readiness port is the root's fixed, granted marker from the
		// candidate. The observer remains open through the post-selection cut.
		n := peerNode{remaining: s.remaining, admission: "marker", wait: "marker"}
		if s.marker == "selected" {
			n.wait = ""
		}
		add := func(label peerLabel, next peerCut) {
			id, found := ids[next]
			if !found {
				id = len(cuts)
				ids[next] = id
				cuts = append(cuts, next)
			}
			n.steps = append(n.steps, peerStep{label, id})
		}
		next := s
		next.pc++
		switch s.pc {
		case decideFirst, decideSecond:
			call := s.pc / 2
			next.remaining, next.original[call] = s.remaining.decide(2)
			decision := next.original[call]
			if decision != "accepted" {
				// Fix the delivered reason at admission, even if the peer spends before selection.
				next.original[call] = client.refusal(s.remaining, decision)
			}
			add(peerLabel{"decision", handles[call], decision}, next)
		case selectFirst, selectSecond:
			call := s.pc / 2
			result := s.original[call]
			if result == "accepted" {
				result = "member_succeeded"
			}
			next.selected[call] = result
			add(peerLabel{"selected", handles[call], result}, next)
		case branchOnReasons:
			if !markerBranch || s.selected != client.markerPair {
				next.pc = candidateDone
			}
			add(peerLabel{}, next)
		case emitMarker:
			next.marker = "pending"
			add(peerLabel{"out", "marker", ""}, next)
		}
		// Permission and delivery run with the granted source and target open.
		// Only the candidate can emit the marker the root selects after admission.
		next = s
		switch s.marker {
		case "pending":
			next.marker = "admitted"
			add(peerLabel{"in", "marker", ""}, next)
		case "admitted":
			next.marker = "selected"
			add(peerLabel{}, next)
		}
		next = s
		switch context {
		case "product":
			switch s.peer {
			case "":
				next.peer = "committed"
				add(peerLabel{}, next)
			case "committed":
				next.remaining, next.peer = s.remaining.decide(1)
				add(peerLabel{"decision", "p", next.peer}, next)
			}
		case "open":
			if s.peer == "" && s.remaining.invocations > 0 && s.remaining.runs > 0 {
				next.remaining = peerBudget{s.remaining.invocations - 1, s.remaining.runs - 1}
				next.peer = "spent"
				add(peerLabel{"in", "rely", "spend"}, next)
			}
		case "isolated":
		default:
			panic(context)
		}
		g = append(g, n)
	}
	return g
}

// The ledger difference telescopes along each path. There are no refunds or
// resets. Only tau edges may surround the single visible label in a weak match.
type peerTransfer struct {
	label peerLabel
	debit peerBudget
}

type peerView struct {
	strong    map[peerTransfer][]int
	weak      map[peerTransfer][]int
	divergent bool
}

func peerViews(g peerGraph) []peerView {
	var tau peerLabel
	hidden := make([][]int, len(g))
	steps := make([][]peerStep, len(g))
	for i, n := range g {
		hidden[i] = []int{i}
		for j := 0; j < len(hidden[i]); j++ {
			for _, e := range g[hidden[i][j]].steps {
				if e.label == tau && !slices.Contains(hidden[i], e.to) {
					hidden[i] = append(hidden[i], e.to)
				}
			}
		}
		steps[i] = slices.Clone(n.steps)
		if !slices.ContainsFunc(n.steps, func(e peerStep) bool { return e.label.kind != "in" }) {
			// At a quiescent node this synthetic challenge can match only a
			// zero-cost hidden path ending at an equally ready quiescent node.
			ready := peerLabel{"ready", n.admission, n.wait}
			steps[i] = append(steps[i], peerStep{ready, i})
		}
	}
	views := make([]peerView, len(g))
	for i, n := range g {
		add := func(m map[peerTransfer][]int, label peerLabel, to int) {
			end := g[to].remaining
			debit := peerBudget{n.remaining.invocations - end.invocations, n.remaining.runs - end.runs}
			k := peerTransfer{label, debit}
			m[k] = append(m[k], to)
		}
		v := peerView{strong: map[peerTransfer][]int{}, weak: map[peerTransfer][]int{}}
		for _, e := range steps[i] {
			add(v.strong, e.label, e.to)
		}
		for _, p := range hidden[i] {
			add(v.weak, tau, p)
			for _, e := range steps[p] {
				if e.label == tau {
					// Infinite hidden paths in a finite graph reach a hidden cycle.
					if slices.Contains(hidden[e.to], p) {
						v.divergent = true
					}
					continue
				}
				for _, q := range hidden[e.to] {
					add(v.weak, e.label, q)
				}
			}
		}
		views[i] = v
	}
	return views
}

// Greatest fixed point on the two finite graphs. Each single-edge challenge
// must have some weak match ending in the current relation, in both directions.
func peerBisimilar(x, y peerGraph) bool {
	a, b := peerViews(x), peerViews(y)
	pairs := map[[2]int]bool{}
	for i := range a {
		for j := range b {
			if a[i].divergent == b[j].divergent {
				pairs[[2]int{i, j}] = true
			}
		}
	}
	transfers := func(left, right peerView, reverse bool) bool {
		for action, successors := range left.strong {
			for _, next := range successors {
				ok := slices.ContainsFunc(right.weak[action], func(other int) bool {
					if reverse {
						return pairs[[2]int{other, next}]
					}
					return pairs[[2]int{next, other}]
				})
				if !ok {
					return false
				}
			}
		}
		return true
	}
	for {
		before := len(pairs)
		for p := range pairs {
			left, right := a[p[0]], b[p[1]]
			if !transfers(left, right, false) || !transfers(right, left, true) {
				delete(pairs, p)
			}
		}
		if len(pairs) == before {
			return pairs[[2]int{0, 0}]
		}
	}
}

func TestPeerSpendingRelation(t *testing.T) {
	for _, context := range []string{"isolated", "product", "open"} {
		t.Run(context, func(t *testing.T) {
			maxI, maxR := 2, 4 // total demand when both original calls are accepted
			if context != "isolated" {
				maxI, maxR = 3, 5
			}
			var profiles, distinguished, states, edges int
			for i := 0; i <= maxI; i++ {
				for r := 0; r <= maxR; r++ {
					b := peerBudget{i, r}
					x, y := literalClient.graph(b, false, context), literalClient.graph(b, true, context)
					wantDifference := context != "isolated" && b == (peerBudget{1, 1})
					if peerBisimilar(x, y) == wantDifference {
						t.Fatalf("remaining %+v: want distinguishable=%t", b, wantDifference)
					}
					marker := peerLabel{"in", "marker", ""}
					for side, g := range [2]peerGraph{x, y} {
						states += len(g)
						admitted := false
						for _, n := range g {
							edges += len(n.steps)
							if slices.ContainsFunc(n.steps, func(e peerStep) bool { return e.label == marker }) {
								admitted = true
							}
						}
						if admitted != (side == 1 && wantDifference) {
							t.Fatalf("side %d, remaining %+v: wrong public marker language", side, b)
						}
					}
					if wantDifference {
						distinguished++
						if !peerBisimilar(x, literalClient.graph(b, false, context)) {
							t.Fatal("identical fallback distinguished")
						}
					}
					profiles++
				}
			}
			t.Logf("%d allowance profiles, %d distinguished (%d states and %d edges across X and Y)", profiles, distinguished, states, edges)
		})
	}
}

func TestPeerSpendingOriginalTrace(t *testing.T) {
	predicates := func(b peerBudget, _ string) string {
		switch {
		case b.invocations == 0 && b.runs < 2:
			return "{I,R}"
		case b.invocations == 0:
			return "{I}"
		default:
			return "{R}" // Called only on refusal of a two-member invocation.
		}
	}
	coarse := func(peerBudget, string) string { return "insufficient_allowance" }
	sets := [2]string{"{R}", "{I,R}"}
	same := [2]string{"insufficient_allowance", "insufficient_allowance"}
	for _, tc := range []struct {
		name     string
		client   peerClient
		selected [2]string
		marks    bool
	}{
		{"priority", literalClient, literalClient.markerPair, true},
		{"sets_literal_guard", peerClient{predicates, literalClient.markerPair}, sets, false},
		{"sets_adapted_guard", peerClient{predicates, sets}, sets, true},
		{"coarse_literal_guard", peerClient{coarse, literalClient.markerPair}, same, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b := peerBudget{1, 1}
			if !peerBisimilar(tc.client.graph(b, false, "isolated"), tc.client.graph(b, true, "isolated")) {
				t.Fatal("distinguished without the peer")
			}
			x, y := tc.client.graph(b, false, "product"), tc.client.graph(b, true, "product")
			if peerBisimilar(x, y) == tc.marks {
				t.Fatalf("with peer: want distinguished=%t", tc.marks)
			}
			// The peer spends before c1 is selected. Both candidates must still
			// read the two original results; only the enabled guard may emit.
			sequence := []peerTransfer{
				{label: peerLabel{"decision", "c1", "created_run_limit"}},
				{peerLabel{"decision", "p", "accepted"}, peerBudget{1, 1}},
				{label: peerLabel{"selected", "c1", tc.selected[0]}},
				{label: peerLabel{"decision", "c2", "invocation_limit"}},
				{label: peerLabel{"selected", "c2", tc.selected[1]}},
				{label: peerLabel{"out", "marker", ""}},
				{label: peerLabel{"in", "marker", ""}},
			}
			for side, g := range [2]peerGraph{x, y} {
				views, states := peerViews(g), []int{0}
				for index, action := range sequence {
					var next []int
					for _, state := range states {
						next = append(next, views[state].weak[action]...)
					}
					if (len(next) > 0) != (index < 5 || side == 1 && tc.marks) {
						t.Fatalf("side %d, step %d, %v: reached %v", side, index, action, next)
					}
					states = next
				}
				// Erase diagnostics only after generating the candidate's behavior.
				for _, n := range g {
					for j, e := range n.steps {
						if e.label.kind == "decision" || e.label.kind == "selected" {
							n.steps[j].label.value = ""
						}
					}
				}
			}
			if peerBisimilar(x, y) == tc.marks {
				t.Fatal("erasing recorded reasons changed the marker distinction")
			}
		})
	}
}

func TestPeerSpendingRelationControls(t *testing.T) {
	sink := peerNode{admission: "marker", wait: "marker"}
	exit := peerLabel{"out", "finish", ""}
	before := peerNode{steps: []peerStep{{exit, 1}}}
	base := peerGraph{before, sink}
	check := func(name string, left, right peerGraph, want bool) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			if got := peerBisimilar(left, right); got != want {
				t.Fatalf("bisimilar = %t, want %t", got, want)
			}
		})
	}
	check("private_stutter", base, peerGraph{{steps: []peerStep{{peerLabel{}, 1}}}, {steps: []peerStep{{exit, 2}}}, sink}, true)
	// Both initial states can exit and are nonquiescent. The optional hidden
	// cycle is distinguishable even though every finite visible trace matches.
	check("hidden_divergence", base, peerGraph{{steps: []peerStep{{peerLabel{}, 0}, {exit, 1}}}, sink}, false)
	check("reverse_transfer", base, peerGraph{{steps: []peerStep{{exit, 1}, {peerLabel{"out", "extra", ""}, 1}}}, sink}, false)
	check("admission_readiness", peerGraph{sink}, peerGraph{{wait: "marker"}}, false)
	check("continuation_readiness", peerGraph{sink}, peerGraph{{admission: "marker"}}, false)
	check("quiescence_after_private_step", peerGraph{sink}, peerGraph{{steps: []peerStep{{peerLabel{}, 1}}}, sink}, true)
	inputOnly := sink
	inputOnly.steps = []peerStep{{peerLabel{"in", "marker", ""}, 0}}
	otherWait := inputOnly
	otherWait.wait = ""
	check("input_enabled_quiescence", peerGraph{inputOnly}, peerGraph{otherWait}, false)
	// The terminal totals agree. The two labeled transfers charge at different cuts.
	early := peerGraph{{remaining: peerBudget{1, 1}, steps: []peerStep{{exit, 1}}}, {steps: []peerStep{{exit, 2}}}, sink}
	late := slices.Clone(early)
	late[1].remaining = peerBudget{1, 1}
	check("exact_debit_per_transfer", early, late, false)

	t.Run("charged_weak_path", func(t *testing.T) {
		path := peerGraph{
			{remaining: peerBudget{1, 2}, steps: []peerStep{{peerLabel{}, 1}}},
			{remaining: peerBudget{0, 2}, steps: []peerStep{{exit, 2}}},
			{remaining: peerBudget{0, 1}, steps: []peerStep{{peerLabel{}, 3}}},
			sink,
		}
		weak := peerViews(path)[0].weak
		ends := weak[peerTransfer{exit, peerBudget{1, 2}}]
		if !slices.Equal(ends, []int{3}) {
			t.Fatalf("charged hidden prefix, visible action, hidden suffix: got endpoints %v", ends)
		}
	})
}
