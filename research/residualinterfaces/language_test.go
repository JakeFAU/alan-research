package residualinterfaces

import (
	"fmt"
	"slices"
	"testing"
)

func tiny(id string, marks []bool, edges ...edge) automaton {
	a := automaton{ID: id, Role: "generic", Signature: []byte(`{}`), Initial: "0", States: []state{}, Edges: edges}
	for i, m := range marks {
		a.States = append(a.States, state{ID: fmt.Sprint(i), Marked: m})
	}
	return a
}
func TestLanguageBoundaries(t *testing.T) {
	a := tiny("a", []bool{false, true}, edge{"0", "tau", "1"}, edge{"1", "tau", "0"})
	if !accepts(a, []string{}) {
		t.Fatal("empty word must include tau closure")
	}
	if accepts(a, []string{labels[0]}) {
		t.Fatal("empty subset accepted")
	}
	b := tiny("b", []bool{true})
	if !compareLanguages(a, b).Equal {
		t.Fatal("tau cycle changes language")
	}
	chain := tiny("chain", make([]bool, 13))
	chain.States[12].Marked = true
	for i := 0; i < 12; i++ {
		chain.Edges = append(chain.Edges, edge{fmt.Sprint(i), labels[0], fmt.Sprint(i + 1)})
	}
	c := compareLanguages(chain, tiny("empty", []bool{false}))
	if c.Equal || len(c.Word) != 12 || !c.LeftAccept || c.RightAccept {
		t.Fatalf("deep separator: %+v", c)
	}
}
func TestProfileMultiplicityAndEmptyLanguage(t *testing.T) {
	a := tiny("a", []bool{true, true}, edge{"0", labels[0], "1"})
	b := tiny("b", []bool{true, true, true}, edge{"0", labels[0], "1"}, edge{"0", labels[0], "2"})
	if !compareProfiles(a, b).Equal {
		t.Fatal("one-to-many equal languages must cover")
	}
	x := tiny("x", []bool{false})
	y := tiny("y", []bool{false, false}, edge{"0", labels[0], "1"})
	c := compareProfiles(x, y)
	if !c.InitialLanguage.Equal || c.Equal || !slices.Equal(c.Word, []string{labels[0]}) || c.Mismatch.Side != "right" || len(c.Mismatch.Separators) != 0 {
		t.Fatalf("empty versus dead: %+v", c)
	}
}
