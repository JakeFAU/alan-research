package residualantichains

import "fmt"

func extraMonitors(signature []byte) []automaton {
	out := []automaton{}
	for _, name := range []string{"minimal-base", "minimal-augmented"} {
		a := automaton{ID: name, Role: "M", Signature: signature, Initial: "0", States: []state{}, Edges: []edge{}}
		n := 4
		if name == "minimal-augmented" {
			n = 5
		}
		for i := 0; i < n; i++ {
			id := fmt.Sprint(i)
			a.States = append(a.States, state{ID: id, Marked: true})
			for _, label := range labels {
				if label != "reserve-capacity:AB" && label != "apply:AB" && label != "commit:BC" {
					a.Edges = append(a.Edges, edge{From: id, Label: label, To: id})
				}
			}
		}
		a.Edges = append(a.Edges, edge{From: "0", Label: "tau", To: "1"})
		for _, label := range []string{"reserve-capacity:AB", "apply:AB", "commit:BC"} {
			a.Edges = append(a.Edges, edge{From: "0", Label: label, To: "3"}, edge{From: "3", Label: label, To: "3"})
		}
		a.Edges = append(a.Edges, edge{From: "1", Label: "reserve-capacity:AB", To: "2"}, edge{From: "2", Label: "apply:AB", To: "3"})
		if n == 5 {
			a.Edges = append(a.Edges, edge{From: "0", Label: "tau", To: "4"}, edge{From: "4", Label: "reserve-capacity:AB", To: "2"}, edge{From: "4", Label: "commit:BC", To: "3"})
		}
		out = append(out, a)
	}
	return out
}
func genericMINLoss() (automaton, automaton) {
	x := automaton{ID: "generic-MIN-X", Role: "generic", Signature: []byte(`{}`), Initial: "s", States: []state{{ID: "s", Marked: true}, {ID: "d"}}, Edges: []edge{{From: "s", Label: "tau", To: "d"}}}
	y := automaton{ID: "generic-MIN-Y", Role: "generic", Signature: []byte(`{}`), Initial: "d", States: []state{{ID: "d"}}, Edges: []edge{}}
	for _, label := range labels {
		x.Edges = append(x.Edges, edge{From: "s", Label: label, To: "s"}, edge{From: "d", Label: label, To: "d"})
		y.Edges = append(y.Edges, edge{From: "d", Label: label, To: "d"})
	}
	return x, y
}
