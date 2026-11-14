package residualantichains

import "fmt"

type pathStep struct {
	From  string
	Label string
	To    string
}
type path struct {
	Start string
	Steps []pathStep
	End   string
}
type coreOutcome struct {
	Monitor            string
	Product            string
	InitialCanComplete bool
	Nonblocking        bool
}
type negativeWitness struct {
	Comparison   int
	LeftProduct  string
	RightProduct string
	LeftPath     path
	RightPath    path
}

func nodesByID(g graph) map[string]node {
	out := map[string]node{}
	for _, n := range g.Nodes {
		out[n.ID] = n
	}
	return out
}

// Product path ties use the stored lexicographically sorted (label,target) arcs.
func shortestPath(g graph, start string, goal func(node) bool) (path, error) {
	nodes := nodesByID(g)
	if _, ok := nodes[start]; !ok {
		return path{}, fmt.Errorf("unknown path start")
	}
	queue := []string{start}
	seen := map[string]bool{start: true}
	parents := map[string]pathStep{}
	for i := 0; i < len(queue); i++ {
		id := queue[i]
		if goal(nodes[id]) {
			p := path{Start: start, Steps: []pathStep{}, End: id}
			for at := id; at != start; {
				e := parents[at]
				p.Steps = append(p.Steps, e)
				at = e.From
			}
			for l, r := 0, len(p.Steps)-1; l < r; l, r = l+1, r-1 {
				p.Steps[l], p.Steps[r] = p.Steps[r], p.Steps[l]
			}
			return p, nil
		}
		for _, e := range nodes[id].Arcs {
			if !seen[e.To] {
				seen[e.To] = true
				parents[e.To] = pathStep{From: id, Label: e.Label, To: e.To}
				queue = append(queue, e.To)
			}
		}
	}
	return path{}, fmt.Errorf("path target unreachable in %s", g.ID)
}
func unpaid(n node) bool { return len(n.Tuple) >= 2 && n.Tuple[0] == "W,U,true,false" }

func makeNegative(core map[string]graph, comparisons []profileComparison) (negativeWitness, error) {
	neg := negativeWitness{Comparison: -1}
	for i, c := range comparisons {
		if c.Left == "release-first" && c.Right == "release-first-doomed" {
			neg.Comparison = i
		}
	}
	if neg.Comparison < 0 {
		return neg, fmt.Errorf("missing MAX negative comparison")
	}
	l, r := core["release-first"], core["release-first-doomed"]
	neg.LeftProduct = l.ID
	neg.RightProduct = r.ID
	var err error
	neg.LeftPath, err = shortestPath(l, l.Initial, func(n node) bool { return unpaid(n) && n.Tuple[1] == "1" })
	if err != nil {
		return neg, err
	}
	neg.RightPath, err = shortestPath(r, r.Initial, func(n node) bool { return unpaid(n) && n.Tuple[1] == "3" })
	return neg, err
}

type legalControl struct {
	Plant      automaton
	Monitor    string
	Roots      []string
	Products   []graph
	Inclusions []comparison
}

func makeLegalControl(k, m automaton) (legalControl, error) {
	k.ID += "/all-marked-control"
	k.States = append([]state{}, k.States...)
	for i := range k.States {
		k.States[i].Marked = true
	}
	out := legalControl{Plant: k, Monitor: m.ID, Roots: []string{"1", "4", "0"}, Products: []graph{}, Inclusions: []comparison{}}
	as := []automaton{}
	for _, root := range out.Roots {
		a := m
		a.Initial = root
		a.ID += "@" + root
		a.Signature = k.Signature
		g, err := product(k, a)
		if err != nil {
			return out, err
		}
		out.Products = append(out.Products, g)
		as = append(as, graphAutomaton(g, k))
	}
	out.Inclusions = append(out.Inclusions, compareLanguages(as[0], as[1]), compareLanguages(as[1], as[2]))
	return out, nil
}

type genericComparison struct {
	Left         automaton
	Right        automaton
	Comparison   profileComparison
	Trie         automaton
	LeftProduct  graph
	RightProduct graph
}

func genericControl(x, y automaton) (genericComparison, error) {
	h := trie("epsilon-trie", [][]string{{}})
	l, err := product(x, h)
	if err != nil {
		return genericComparison{}, err
	}
	r, err := product(y, h)
	if err != nil {
		return genericComparison{}, err
	}
	return genericComparison{Left: x, Right: y, Comparison: compareProfiles(x, y), Trie: h, LeftProduct: l, RightProduct: r}, nil
}

type hitControl struct {
	X       automaton
	H       [][]string
	Trie    automaton
	Product graph
}
type multiwordControl struct {
	Automaton    automaton
	SourceStates []string
	TargetState  string
	Coverage     []coverage
	Trie         automaton
	Products     []graph
}
type genericWitnesses struct {
	MINLoss       genericComparison
	JNonnecessity genericComparison
	Multiword     multiwordControl
	LocalHitNotNB hitControl
}

// A trie accepts exactly the finite suffix set, including the distinction
// between no words (unmarked root) and the singleton empty word (marked root).
func trie(id string, words [][]string) automaton {
	a := automaton{ID: id, Role: "generic", Signature: []byte(`{}`), Initial: "0", States: []state{{ID: "0"}}, Edges: []edge{}}
	next := map[string]int{}
	for _, word := range words {
		at := 0
		for _, label := range word {
			k := key([]string{fmt.Sprint(at), label})
			to, ok := next[k]
			if !ok {
				to = len(a.States)
				next[k] = to
				a.States = append(a.States, state{ID: fmt.Sprint(to)})
				a.Edges = append(a.Edges, edge{From: fmt.Sprint(at), Label: label, To: fmt.Sprint(to)})
			}
			at = to
		}
		a.States[at].Marked = true
	}
	return a
}
func makeGenerics() (genericWitnesses, error) {
	out := genericWitnesses{}
	x, y := genericMINLoss()
	var err error
	out.MINLoss, err = genericControl(x, y)
	if err != nil {
		return out, err
	}
	x = automaton{ID: "generic-unmarked-X", Role: "generic", Signature: []byte(`{}`), Initial: "0", States: []state{{ID: "0"}}, Edges: []edge{}}
	y = automaton{ID: "generic-unmarked-Y", Role: "generic", Signature: []byte(`{}`), Initial: "0", States: []state{{ID: "0"}, {ID: "1"}}, Edges: []edge{{From: "0", Label: labels[2], To: "1"}}}
	out.JNonnecessity, err = genericControl(x, y)
	if err != nil {
		return out, err
	}
	a := automaton{ID: "generic-three-singletons", Role: "generic", Signature: []byte(`{}`), Initial: "0",
		States: []state{{ID: "0"}, {ID: "1"}, {ID: "2"}, {ID: "3", Marked: true}},
		Edges:  []edge{{From: "0", Label: labels[2], To: "3"}, {From: "1", Label: labels[1], To: "3"}, {From: "2", Label: labels[4], To: "3"}}}
	idx := residualIndex(residuals(a, a))
	f, g := reduce([]string{"0", "1"}, idx), reduce([]string{"2"}, idx)
	cov := cover(a, a, f, g, idx, false)
	h := trie("multiword-H", cov[0].Failure.H)
	multi := multiwordControl{Automaton: a, SourceStates: []string{"0", "1"}, TargetState: "2", Coverage: cov, Trie: h, Products: []graph{}}
	for _, root := range []string{"0", "1", "2"} {
		b := a
		b.Initial = root
		b.ID += "@" + root
		g, err := product(b, h)
		if err != nil {
			return out, err
		}
		multi.Products = append(multi.Products, g)
	}
	out.Multiword = multi
	x = automaton{ID: "generic-local-hit", Role: "generic", Signature: []byte(`{}`), Initial: "q",
		States: []state{{ID: "q"}, {ID: "f", Marked: true}, {ID: "d"}}, Edges: []edge{{From: "q", Label: labels[2], To: "f"}, {From: "q", Label: labels[2], To: "d"}}}
	words := [][]string{{labels[2]}}
	h = trie("local-hit-H", words)
	prod, err := product(x, h)
	if err != nil {
		return out, err
	}
	out.LocalHitNotNB = hitControl{X: x, H: words, Trie: h, Product: prod}
	return out, nil
}
