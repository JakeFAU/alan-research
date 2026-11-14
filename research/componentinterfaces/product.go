package componentinterfaces

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type arc struct {
	Label string
	To    string
}
type node struct {
	ID          string
	Tuple       []string
	Marked      bool
	CanComplete bool
	Arcs        []arc
}
type graph struct {
	ID          string
	Factors     []string
	Initial     string
	Nodes       []node
	States      int
	Arcs        int
	CanComplete int
	Nonblocking bool
}
type grouping struct {
	Product string
	Left    graph
	Right   graph
	Equal   bool
}

func key(tuple []string) string { b, _ := json.Marshal(tuple); return string(b) }

// Visible labels synchronize in every factor; each tau advances one factor.
func product(as ...automaton) (graph, error) {
	if err := compatible(as); err != nil {
		return graph{}, err
	}
	g := graph{Factors: []string{}, Nodes: []node{}}
	start := []string{}
	for _, a := range as {
		g.Factors = append(g.Factors, a.ID)
		start = append(start, a.Initial)
	}
	g.ID = strings.Join(g.Factors, " || ")
	flatten := func(local []string) []string {
		t := []string{}
		for i, id := range local {
			t = append(t, tuple(as[i], id)...)
		}
		return t
	}
	queue := [][]string{start}
	g.Initial = key(flatten(start))
	seen := map[string]bool{key(start): true}
	for i := 0; i < len(queue); i++ {
		local := queue[i]
		flat := flatten(local)
		n := node{ID: key(flat), Tuple: flat, Marked: true, Arcs: []arc{}}
		for j, a := range as {
			n.Marked = n.Marked && marked(a, local[j])
		}
		add := func(label string, next []string) {
			k := key(next)
			if !seen[k] {
				seen[k] = true
				queue = append(queue, slices.Clone(next))
			}
			n.Arcs = append(n.Arcs, arc{label, key(flatten(next))})
		}
		for j, a := range as {
			for _, to := range targets(a, local[j], "tau") {
				next := slices.Clone(local)
				next[j] = to
				add("tau", next)
			}
		}
		for _, label := range alphabet(as[0].Signature) {
			var cross func(int, []string)
			cross = func(j int, next []string) {
				if j == len(as) {
					add(label.ID, next)
					return
				}
				for _, to := range targets(as[j], local[j], label.ID) {
					cross(j+1, append(next, to))
				}
			}
			cross(0, nil)
		}
		slices.SortFunc(n.Arcs, func(a, b arc) int {
			if x := strings.Compare(a.Label, b.Label); x != 0 {
				return x
			}
			return strings.Compare(a.To, b.To)
		})
		n.Arcs = slices.Compact(n.Arcs)
		g.Arcs += len(n.Arcs)
		g.Nodes = append(g.Nodes, n)
	}
	ids := map[string]int{}
	incoming := map[string][]int{}
	back := []int{}
	for i, n := range g.Nodes {
		ids[n.ID] = i
		for _, e := range n.Arcs {
			incoming[e.To] = append(incoming[e.To], i)
		}
		if n.Marked {
			g.Nodes[i].CanComplete = true
			back = append(back, i)
		}
	}
	for i := 0; i < len(back); i++ {
		for _, from := range incoming[g.Nodes[back[i]].ID] {
			if !g.Nodes[from].CanComplete {
				g.Nodes[from].CanComplete = true
				back = append(back, from)
			}
		}
	}
	g.States = len(g.Nodes)
	g.CanComplete = len(back)
	g.Nonblocking = g.States == g.CanComplete
	if _, ok := ids[g.Initial]; !ok {
		return graph{}, fmt.Errorf("missing initial")
	}
	return g, nil
}
func graphAutomaton(g graph, signature automaton) automaton {
	a := automaton{ID: g.ID, Role: "product", Signature: signature.Signature, Initial: g.Initial, States: []state{}, Edges: []edge{}}
	for _, n := range g.Nodes {
		a.States = append(a.States, state{ID: n.ID, Tuple: n.Tuple, Marked: n.Marked})
		for _, e := range n.Arcs {
			a.Edges = append(a.Edges, edge{n.ID, e.Label, e.To})
		}
	}
	return a
}
func sameRelation(a, b graph) bool {
	if a.Initial != b.Initial || a.States != b.States || a.Arcs != b.Arcs || a.CanComplete != b.CanComplete || a.Nonblocking != b.Nonblocking {
		return false
	}
	left, right := slices.Clone(a.Nodes), slices.Clone(b.Nodes)
	cmp := func(a, b node) int { return strings.Compare(a.ID, b.ID) }
	slices.SortFunc(left, cmp)
	slices.SortFunc(right, cmp)
	return reflect.DeepEqual(left, right)
}
func checkGrouping(k, m, d automaton, flat graph) (grouping, error) {
	km, e := product(k, m)
	if e != nil {
		return grouping{}, e
	}
	md, e := product(m, d)
	if e != nil {
		return grouping{}, e
	}
	left, e := product(graphAutomaton(km, k), d)
	if e != nil {
		return grouping{}, e
	}
	right, e := product(k, graphAutomaton(md, m))
	if e != nil {
		return grouping{}, e
	}
	return grouping{flat.ID, left, right, sameRelation(flat, left) && sameRelation(flat, right)}, nil
}
