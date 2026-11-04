package twoagent

import "slices"

// GraphEdge preserves all transitions, including self-loops and previously
// discovered targets. A revisit is not treated as a completed execution.
type GraphEdge struct {
	Event Event
	To    int
}

type Node struct {
	State       State
	Edges       []GraphEdge
	Parent      int
	Via         Event
	CanComplete bool
}

type Witness struct {
	StateID     int
	State       State
	Trace       []Event
	CanComplete bool
}

type Pair struct {
	Left             Witness
	Right            Witness
	RecoverableEqual bool
}

// Report contains exhaustive finite-graph counts. Nodes exposes the full graph
// to callers; JSON exports the counts and shortest falsifying witnesses only.
// CanComplete is existential reachability of A's durable authentic consumption,
// not joint termination, B's release, channel drain or global quiescence. It is
// never a universal liveness claim.
type Report struct {
	Profile                   Profile
	States                    int
	Edges                     int
	CyclicStates              int
	CyclicComponents          int
	LocalSafetyViolations     int
	LocalTransitionViolations int
	CustodyViolations         int
	CanComplete               int
	CannotComplete            int
	ShortestCustodyViolation  *Witness
	LossWitnesses             []Witness
	Indistinguishable         *Pair
	Nodes                     []Node `json:"-"`
}

// Explore enumerates the entire reachable graph by BFS without a depth limit,
// then computes completion by reverse reachability and cycles by SCCs.
func Explore(p Profile) Report {
	r := Report{Profile: p, Nodes: []Node{{State: Initial(), Parent: -1}}, LossWitnesses: []Witness{}}
	ids := map[State]int{Initial(): 0}
	var incoming [][]int
	for i := 0; i < len(r.Nodes); i++ {
		state := r.Nodes[i].State
		for _, edge := range Successors(p, state) {
			if !localTransitionSafe(state, edge) {
				r.LocalTransitionViolations++
			}
			to, ok := ids[edge.State]
			if !ok {
				to = len(r.Nodes)
				ids[edge.State] = to
				r.Nodes = append(r.Nodes, Node{State: edge.State, Parent: i, Via: edge.Event})
			}
			r.Nodes[i].Edges = append(r.Nodes[i].Edges, GraphEdge{Event: edge.Event, To: to})
			r.Edges++
		}
	}
	r.States = len(r.Nodes)
	incoming = make([][]int, r.States)
	var queue []int
	for i, node := range r.Nodes {
		for _, edge := range node.Edges {
			incoming[edge.To] = append(incoming[edge.To], i)
		}
		if node.State.AConsumed {
			r.Nodes[i].CanComplete = true
			queue = append(queue, i)
		}
	}
	for i := 0; i < len(queue); i++ {
		for _, from := range incoming[queue[i]] {
			if !r.Nodes[from].CanComplete {
				r.Nodes[from].CanComplete = true
				queue = append(queue, from)
			}
		}
	}
	r.CanComplete = len(queue)
	r.CannotComplete = r.States - r.CanComplete
	seenLoss := map[Value]bool{}
	for i, node := range r.Nodes {
		s := node.State
		if !locallySafe(s) {
			r.LocalSafetyViolations++
		}
		if s.Applied && s.BResult == "" && s.ARecorded == "" {
			r.CustodyViolations++
			if r.ShortestCustodyViolation == nil {
				w := r.witness(i)
				r.ShortestCustodyViolation = &w
			}
		}
		recovered := s.Crashed && !s.ADown
		noReceipt := s.BResult == "" && s.ARecorded == "" && s.AVolatile == "" && s.ResultChannel == ""
		if s.Applied && recovered && noReceipt && !seenLoss[s.WorldResult] {
			seenLoss[s.WorldResult] = true
			r.LossWitnesses = append(r.LossWitnesses, r.witness(i))
		}
	}
	if len(r.LossWitnesses) == 2 {
		a, b := r.LossWitnesses[0], r.LossWitnesses[1]
		r.Indistinguishable = &Pair{Left: a, Right: b, RecoverableEqual: Recoverable(a.State) == Recoverable(b.State)}
	}
	r.CyclicStates, r.CyclicComponents = cycles(r.Nodes)
	return r
}

// Trace returns the BFS shortest path to a node. Edges closing cycles remain
// available on Nodes; the shortest predecessor tree does not erase them.
func (r Report) Trace(i int) []Event {
	trace := []Event{}
	for r.Nodes[i].Parent >= 0 {
		trace = append(trace, r.Nodes[i].Via)
		i = r.Nodes[i].Parent
	}
	slices.Reverse(trace)
	return trace
}

func (r Report) witness(i int) Witness {
	return Witness{StateID: i, State: r.Nodes[i].State, Trace: r.Trace(i), CanComplete: r.Nodes[i].CanComplete}
}

func locallySafe(s State) bool {
	if s.Delegation != FixedDelegation() || s.Grant != FixedGrant {
		return false
	}
	if s.Applied != (s.WorldResult != "") || (s.Applied && !s.BAdmitted) || (s.BAdmitted && !s.Committed) {
		return false
	}
	if s.WorldResult != "" && s.WorldResult != R1 && s.WorldResult != R2 {
		return false
	}
	for _, v := range []Value{s.AVolatile, s.ARecorded, s.BResult, s.ResultChannel} {
		if v != "" && (s.WorldResult == "" || v != s.WorldResult) {
			return false
		}
	}
	consumedHasReceipt := !s.AConsumed || s.ARecorded != ""
	downLostVolatile := !s.ADown || (s.Crashed && s.AVolatile == "")
	activationValid := s.Activation == 0 || (s.Activation == 1 && s.Crashed && !s.ADown)
	return consumedHasReceipt && downLostVolatile && activationValid
}

// Tarjan's algorithm counts states in nontrivial SCCs or singleton SCCs with a
// self-loop. Reachability merges alone are not evidence of cycles.
func cycles(nodes []Node) (int, int) {
	index := make([]int, len(nodes))
	low := make([]int, len(nodes))
	onStack := make([]bool, len(nodes))
	for i := range index {
		index[i] = -1
	}
	var stack []int
	next, states, components := 0, 0, 0
	var visit func(int)
	visit = func(v int) {
		index[v] = next
		low[v] = next
		next++
		stack = append(stack, v)
		onStack[v] = true
		for _, e := range nodes[v].Edges {
			w := e.To
			if index[w] < 0 {
				visit(w)
				low[v] = min(low[v], low[w])
			} else if onStack[w] {
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
			onStack[w] = false
			size++
			if w == v {
				break
			}
		}
		cyclic := size > 1
		if !cyclic {
			for _, e := range nodes[v].Edges {
				if e.To == v {
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
	for i := range nodes {
		if index[i] < 0 {
			visit(i)
		}
	}
	return states, components
}

// localTransitionSafe separately checks preservation of committed facts and
// the one-application/one-consumption transition bounds on every graph edge.
func localTransitionSafe(s State, e Edge) bool {
	n := e.State
	if s.Delegation != n.Delegation || s.Grant != n.Grant {
		return false
	}
	lostACommit := (s.Committed && !n.Committed) || (s.AConsumed && !n.AConsumed)
	lostBCommit := (s.BAdmitted && !n.BAdmitted) || (s.Applied && !n.Applied)
	if lostACommit || lostBCommit || (s.Crashed && !n.Crashed) {
		return false
	}
	if s.ARecorded != "" && n.ARecorded != s.ARecorded {
		return false
	}
	if s.WorldResult != "" && n.WorldResult != s.WorldResult {
		return false
	}
	if e.Event.Kind == Apply {
		if s.Applied || !n.Applied || n.BResult != e.Event.Result {
			return false
		}
	} else if s.Applied != n.Applied || s.WorldResult != n.WorldResult {
		return false
	}
	if e.Event.Kind == Consume {
		if s.AConsumed || !n.AConsumed || n.ARecorded != e.Event.Result {
			return false
		}
	} else if s.AConsumed != n.AConsumed {
		return false
	}
	if e.Event.Kind == Recover {
		return s.ADown && !n.ADown && n.Activation == 1
	}
	return s.Activation == n.Activation
}
