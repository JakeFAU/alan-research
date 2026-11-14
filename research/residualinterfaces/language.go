package residualinterfaces

import (
	"encoding/json"
	"slices"
)

type subsetArc struct {
	Label string
	To    int
}
type subsetNode struct {
	Left        []string
	Right       []string
	LeftAccept  bool
	RightAccept bool
	Arcs        []subsetArc
}
type comparison struct {
	Left        string
	Right       string
	Equal       bool
	Word        []string
	LeftAccept  bool
	RightAccept bool
	Pairs       []subsetNode
}

func closure(a automaton, seed []string) []string {
	seen := map[string]bool{}
	q := []string{}
	for _, s := range seed {
		if !seen[s] {
			seen[s] = true
			q = append(q, s)
		}
	}
	for i := 0; i < len(q); i++ {
		for _, to := range targets(a, q[i], "tau") {
			if !seen[to] {
				seen[to] = true
				q = append(q, to)
			}
		}
	}
	slices.Sort(q)
	return q
}
func advance(a automaton, set []string, label string) []string {
	out := []string{}
	for _, s := range set {
		out = append(out, targets(a, s, label)...)
	}
	return closure(a, out)
}
func acceptsSet(a automaton, set []string) bool {
	for _, s := range set {
		if marked(a, s) {
			return true
		}
	}
	return false
}
func accepts(a automaton, word []string) bool {
	s := closure(a, []string{a.Initial})
	for _, label := range word {
		s = advance(a, s, label)
	}
	return acceptsSet(a, s)
}
func subsetKey(l, r []string) string { b, _ := json.Marshal([][]string{l, r}); return string(b) }

// BFS over epsilon-closed subset pairs is exact for finite epsilon-NFAs.
// All pairs and all alphabet arcs are retained even after a separator is found.
func compareLanguages(a, b automaton) comparison {
	r := comparison{Left: a.ID, Right: b.ID, Equal: true, Word: []string{}, Pairs: []subsetNode{}}
	add := func(l, rr []string) {
		r.Pairs = append(r.Pairs, subsetNode{Left: l, Right: rr, LeftAccept: acceptsSet(a, l), RightAccept: acceptsSet(b, rr), Arcs: []subsetArc{}})
	}
	l, rr := closure(a, []string{a.Initial}), closure(b, []string{b.Initial})
	add(l, rr)
	ids := map[string]int{subsetKey(l, rr): 0}
	parents := []int{-1}
	via := []string{""}
	for i := 0; i < len(r.Pairs); i++ {
		n := r.Pairs[i]
		if r.Equal && n.LeftAccept != n.RightAccept {
			r.Equal = false
			r.LeftAccept = n.LeftAccept
			r.RightAccept = n.RightAccept
			for j := i; parents[j] >= 0; j = parents[j] {
				r.Word = append(r.Word, via[j])
			}
			slices.Reverse(r.Word)
		}
		for _, label := range labels {
			nl, nr := advance(a, n.Left, label), advance(b, n.Right, label)
			key := subsetKey(nl, nr)
			to, ok := ids[key]
			if !ok {
				to = len(r.Pairs)
				ids[key] = to
				add(nl, nr)
				parents = append(parents, i)
				via = append(via, label)
			}
			r.Pairs[i].Arcs = append(r.Pairs[i].Arcs, subsetArc{label, to})
		}
	}
	return r
}
