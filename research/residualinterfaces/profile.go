package residualinterfaces

import "slices"

type residualComparison struct {
	LeftState  string
	RightState string
	Language   comparison
}
type coverage struct {
	State   string
	Matches []string
}
type profileNode struct {
	Left          []string
	Right         []string
	LeftCoverage  []coverage
	RightCoverage []coverage
	Equal         bool
	Arcs          []subsetArc
}
type separator struct {
	OppositeState string
	Residual      int
	Word          []string
	LeftAccept    bool
	RightAccept   bool
}
type mismatch struct {
	Side       string
	State      string
	Separators []separator
}
type profileComparison struct {
	Left            string
	Right           string
	InitialLanguage comparison
	Residuals       []residualComparison
	Pairs           []profileNode
	Equal           bool
	Word            []string
	Mismatch        *mismatch
}

// Every cross-state language comparison is retained. Coverage is bidirectional
// and permits multiple partners; no minimization or state bijection is used.
func compareProfiles(a, b automaton) profileComparison {
	r := profileComparison{Left: a.ID, Right: b.ID, InitialLanguage: compareLanguages(a, b),
		Residuals: []residualComparison{}, Pairs: []profileNode{}, Equal: true, Word: []string{}}
	lookup := map[string]int{}
	for _, x := range a.States {
		for _, y := range b.States {
			aa, bb := a, b
			aa.Initial = x.ID
			bb.Initial = y.ID
			lookup[key([]string{x.ID, y.ID})] = len(r.Residuals)
			r.Residuals = append(r.Residuals, residualComparison{LeftState: x.ID, RightState: y.ID, Language: compareLanguages(aa, bb)})
		}
	}
	add := func(l, rr []string) {
		n := profileNode{Left: l, Right: rr, LeftCoverage: []coverage{}, RightCoverage: []coverage{}, Equal: true, Arcs: []subsetArc{}}
		for _, x := range l {
			c := coverage{State: x, Matches: []string{}}
			for _, y := range rr {
				if r.Residuals[lookup[key([]string{x, y})]].Language.Equal {
					c.Matches = append(c.Matches, y)
				}
			}
			n.Equal = n.Equal && len(c.Matches) > 0
			n.LeftCoverage = append(n.LeftCoverage, c)
		}
		for _, y := range rr {
			c := coverage{State: y, Matches: []string{}}
			for _, x := range l {
				if r.Residuals[lookup[key([]string{x, y})]].Language.Equal {
					c.Matches = append(c.Matches, x)
				}
			}
			n.Equal = n.Equal && len(c.Matches) > 0
			n.RightCoverage = append(n.RightCoverage, c)
		}
		r.Pairs = append(r.Pairs, n)
	}
	l, rr := closure(a, []string{a.Initial}), closure(b, []string{b.Initial})
	add(l, rr)
	ids := map[string]int{subsetKey(l, rr): 0}
	parents := []int{-1}
	via := []string{""}
	for i := 0; i < len(r.Pairs); i++ {
		n := r.Pairs[i]
		if r.Equal && !n.Equal {
			r.Equal = false
			for j := i; parents[j] >= 0; j = parents[j] {
				r.Word = append(r.Word, via[j])
			}
			slices.Reverse(r.Word)
			m := mismatch{Separators: []separator{}}
			for _, c := range n.LeftCoverage {
				if len(c.Matches) == 0 {
					m.Side = "left"
					m.State = c.State
					break
				}
			}
			if m.Side == "" {
				for _, c := range n.RightCoverage {
					if len(c.Matches) == 0 {
						m.Side = "right"
						m.State = c.State
						break
					}
				}
			}
			opposite := n.Right
			if m.Side == "right" {
				opposite = n.Left
			}
			for _, s := range opposite {
				pair := []string{m.State, s}
				if m.Side == "right" {
					pair = []string{s, m.State}
				}
				index := lookup[key(pair)]
				c := r.Residuals[index].Language
				m.Separators = append(m.Separators, separator{OppositeState: s, Residual: index, Word: c.Word, LeftAccept: c.LeftAccept, RightAccept: c.RightAccept})
			}
			r.Mismatch = &m
		}
		for _, label := range labels {
			nl, nr := advance(a, n.Left, label), advance(b, n.Right, label)
			k := subsetKey(nl, nr)
			to, ok := ids[k]
			if !ok {
				to = len(r.Pairs)
				ids[k] = to
				add(nl, nr)
				parents = append(parents, i)
				via = append(via, label)
			}
			r.Pairs[i].Arcs = append(r.Pairs[i].Arcs, subsetArc{Label: label, To: to})
		}
	}
	return r
}
