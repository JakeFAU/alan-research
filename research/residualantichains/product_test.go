package residualantichains

import (
	"fmt"
	"testing"
)

func TestConditionalProductClosure(t *testing.T) {
	ms := extraMonitors([]byte(`{}`))
	a, b := ms[0], ms[1]
	for marks := 0; marks < 4; marks++ {
		for bits := 0; bits < 32; bits++ {
			x := tiny(fmt.Sprintf("context-%d-%d", marks, bits), []bool{marks&1 != 0, marks&2 != 0})
			choices := []edge{{From: "0", Label: labels[1], To: "0"}, {From: "0", Label: labels[1], To: "1"}, {From: "1", Label: labels[2], To: "0"}, {From: "1", Label: labels[4], To: "1"}, {From: "0", Label: "tau", To: "1"}}
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
			c := compareProfiles(graphAutomaton(l, x), graphAutomaton(r, x))
			if !c.J || l.Nonblocking != r.Nonblocking || l.Nodes[0].CanComplete != r.Nodes[0].CanComplete {
				t.Fatalf("conditional closure: %s", x.ID)
			}
		}
	}
}
func TestHitIsNotNonblocking(t *testing.T) {
	x := tiny("branch", []bool{false, true, false}, edge{From: "0", Label: labels[2], To: "1"}, edge{From: "0", Label: labels[2], To: "2"})
	h := tiny("H", []bool{false, true}, edge{From: "0", Label: labels[2], To: "1"})
	g, err := product(x, h)
	if err != nil {
		t.Fatal(err)
	}
	if !g.Nodes[0].CanComplete || g.Nonblocking || g.States != 3 || g.CanComplete != 2 {
		t.Fatalf("local hit was promoted to NB: %+v", g)
	}
}
