package residualinterfaces

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
type loopWitness struct {
	Product string
	Prefix  path
	Cycle   path
	Escapes []path
}
type strictnessWitness struct {
	Left       automaton
	Right      automaton
	Comparison profileComparison
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
func makeWitnesses(core map[string]graph, comparisons []profileComparison) (negativeWitness, loopWitness, error) {
	neg := negativeWitness{Comparison: -1}
	for i, c := range comparisons {
		if c.Left == "release-first" && c.Right == "release-first-doomed" {
			neg.Comparison = i
		}
	}
	if neg.Comparison < 0 {
		return neg, loopWitness{}, fmt.Errorf("missing negative comparison")
	}
	left, right := core["release-first"], core["release-first-doomed"]
	neg.LeftProduct = left.ID
	neg.RightProduct = right.ID
	var err error
	neg.LeftPath, err = shortestPath(left, left.Initial, func(n node) bool { return unpaid(n) && n.Tuple[1] == "1" })
	if err != nil {
		return neg, loopWitness{}, err
	}
	neg.RightPath, err = shortestPath(right, right.Initial, func(n node) bool { return unpaid(n) && n.Tuple[1] == "3" })
	if err != nil {
		return neg, loopWitness{}, err
	}
	g := core["release-first-loop"]
	loop := loopWitness{Product: g.ID, Escapes: []path{}}
	loop.Prefix, err = shortestPath(g, g.Initial, func(n node) bool { return unpaid(n) && n.Tuple[1] == "1" })
	if err != nil {
		return neg, loop, err
	}
	first := loop.Prefix.End
	other := key([]string{"W,U,true,false", "3"})
	loop.Cycle = path{Start: first, Steps: []pathStep{{From: first, Label: "tau", To: other}, {From: other, Label: "tau", To: first}}, End: first}
	for _, id := range []string{first, other} {
		escape, e := shortestPath(g, id, func(n node) bool { return n.Marked })
		if e != nil {
			return neg, loop, e
		}
		loop.Escapes = append(loop.Escapes, escape)
	}
	return neg, loop, nil
}
