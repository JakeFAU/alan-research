package residualinterfaces

import (
	"fmt"
	"slices"
	"testing"
)

func TestConditionalProfileClosureIndependentSmallAutomata(t *testing.T) {
	// Equal monitors with different state counts, marks and independent private
	// paths. Contexts enumerate marks and absent/deterministic/nondeterministic
	// visible edges; no ALAN plant or historical outcomes supply these fixtures.
	a := tiny("A", []bool{false, true}, edge{"0", labels[0], "1"}, edge{"1", labels[0], "1"})
	b := tiny("B", []bool{false, false, true, true}, edge{"0", "tau", "1"}, edge{"1", "tau", "0"}, edge{"0", labels[0], "2"}, edge{"1", labels[0], "3"}, edge{"2", labels[0], "2"}, edge{"3", labels[0], "3"})
	if !compareProfiles(a, b).Equal {
		t.Fatal("control monitors must have equal profiles")
	}
	for marks := 0; marks < 4; marks++ {
		for bits := 0; bits < 32; bits++ {
			x := tiny(fmt.Sprintf("X-%d-%d", marks, bits), []bool{marks&1 != 0, marks&2 != 0})
			choices := []edge{{"0", labels[0], "0"}, {"0", labels[0], "1"}, {"1", labels[0], "0"}, {"1", labels[0], "1"}, {"0", "tau", "1"}}
			for i, e := range choices {
				if bits&(1<<i) != 0 {
					x.Edges = append(x.Edges, e)
				}
			}
			l, err := product(x, a)
			if err != nil {
				t.Fatal(err)
			}
			r, err := product(x, b)
			if err != nil {
				t.Fatal(err)
			}
			if l.Nonblocking != r.Nonblocking || !compareProfiles(graphAutomaton(l, x), graphAutomaton(r, x)).Equal {
				t.Fatalf("closure failed %s", x.ID)
			}
			// A direct reachable-pair check includes epsilon, repeated visible steps and
			// unreachable letters; compare Cartesian endpoints without language verdicts.
			for _, word := range [][]string{{}, {labels[0]}, {labels[0], labels[0]}, {labels[1]}} {
				reach := func(z automaton) []string {
					s := closure(z, []string{z.Initial})
					for _, v := range word {
						s = advance(z, s, v)
					}
					return s
				}
				want := []string{}
				for _, u := range reach(x) {
					for _, v := range reach(a) {
						want = append(want, key([]string{u, v}))
					}
				}
				slices.Sort(want)
				if !slices.Equal(reach(graphAutomaton(l, x)), want) {
					t.Fatalf("Cartesian reachability %s %v", x.ID, word)
				}
			}
		}
	}
}
