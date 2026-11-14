package residualantichains

import (
	"fmt"
	"slices"
	"testing"
)

func tiny(id string, marks []bool, edges ...edge) automaton {
	a := automaton{ID: id, Role: "generic", Signature: []byte(`{}`), Initial: "0", States: []state{}, Edges: append([]edge{}, edges...)}
	for i, m := range marks {
		a.States = append(a.States, state{ID: fmt.Sprint(i), Marked: m})
	}
	return a
}

// Catches an XOR/equality test substituted for directed inclusion.
func TestDirectedInclusion(t *testing.T) {
	a := tiny("epsilon", []bool{true})
	b := tiny("epsilon-or-a", []bool{true, true}, edge{From: "0", Label: labels[2], To: "1"})
	c := compareLanguages(a, b)
	if !c.LeftSubset || c.RightSubset || c.LeftOutside != nil || c.RightOutside == nil || !slices.Equal(*c.RightOutside, []string{labels[2]}) {
		t.Fatalf("proper inclusion incorrectly decided: %+v", c)
	}
}

// Catches missing epsilon closure, leading/trailing tau, nondeterminism, graph truncation and bounded sampling.
func TestLanguageBoundaries(t *testing.T) {
	a := tiny("tau", []bool{false, false, true}, edge{From: "0", Label: "tau", To: "1"}, edge{From: "1", Label: "tau", To: "0"}, edge{From: "1", Label: labels[0], To: "0"}, edge{From: "1", Label: labels[0], To: "2"})
	if accepts(a, []string{}) || !accepts(a, []string{labels[0]}) {
		t.Fatal("tau/nondeterministic acceptance")
	}
	eps := tiny("epsilon", []bool{false, true}, edge{From: "0", Label: "tau", To: "1"})
	c := compareLanguages(eps, tiny("empty", []bool{false}))
	if c.LeftOutside == nil || len(*c.LeftOutside) != 0 || !c.RightSubset {
		t.Fatal("epsilon witness absent or wrong direction")
	}
	chain := tiny("chain", make([]bool, 14))
	chain.States[13].Marked = true
	for i := 0; i < 13; i++ {
		chain.Edges = append(chain.Edges, edge{From: fmt.Sprint(i), Label: labels[0], To: fmt.Sprint(i + 1)})
	}
	c = compareLanguages(chain, tiny("empty", []bool{false}))
	if c.LeftOutside == nil || len(*c.LeftOutside) != 13 || len(c.Pairs) != 15 {
		t.Fatalf("deep or incomplete certificate: %+v", c)
	}
	for _, n := range c.Pairs {
		if len(n.Arcs) != len(labels) {
			t.Fatal("incomplete graph")
		}
	}
}

func TestBothDirectionalWitnesses(t *testing.T) {
	a := tiny("r", []bool{false, false, true}, edge{From: "0", Label: labels[1], To: "1"}, edge{From: "1", Label: "tau", To: "2"})
	b := tiny("a", []bool{false, true}, edge{From: "0", Label: labels[2], To: "1"})
	c := compareLanguages(a, b)
	if c.LeftSubset || c.RightSubset || c.LeftOutside == nil || c.RightOutside == nil {
		t.Fatal("lost one direction of incomparable languages")
	}
	if !slices.Equal(*c.LeftOutside, []string{labels[1]}) || !slices.Equal(*c.RightOutside, []string{labels[2]}) || !slices.Equal(c.Word, []string{labels[1]}) {
		t.Fatal("shortest witnesses or alphabet tie order")
	}
	if !accepts(a, *c.LeftOutside) || accepts(b, *c.LeftOutside) || accepts(a, *c.RightOutside) || !accepts(b, *c.RightOutside) {
		t.Fatal("wrong directional acceptance")
	}
}
