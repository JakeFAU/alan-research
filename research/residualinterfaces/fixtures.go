package residualinterfaces

import "fmt"

func extraMonitors(signature []byte) []automaton {
	out := []automaton{}
	for _, name := range []string{"release-first-duplicate", "release-first-loop"} {
		a := automaton{ID: name, Role: "M", Signature: signature, Initial: "0", States: []state{}, Edges: []edge{}}
		for i := 0; i < 4; i++ {
			id := fmt.Sprint(i)
			a.States = append(a.States, state{ID: id, Marked: true})
			for _, label := range labels {
				if label != "reserve-capacity:AB" && label != "apply:AB" && label != "commit:BC" {
					a.Edges = append(a.Edges, edge{id, label, id})
				}
			}
		}
		a.Edges = append(a.Edges, edge{"0", "reserve-capacity:AB", "1"}, edge{"1", "apply:AB", "2"})
		loops := []string{"2"}
		if name == "release-first-duplicate" {
			a.Edges = append(a.Edges, edge{"1", "apply:AB", "3"})
			loops = append(loops, "3")
		}
		if name == "release-first-loop" {
			a.Edges = append(a.Edges, edge{"1", "tau", "3"}, edge{"3", "tau", "1"}, edge{"3", "apply:AB", "2"})
		}
		for _, id := range loops {
			for _, label := range []string{"reserve-capacity:AB", "apply:AB", "commit:BC"} {
				a.Edges = append(a.Edges, edge{id, label, id})
			}
		}
		out = append(out, a)
	}
	return out
}
func genericStrictness() (automaton, automaton) {
	x := automaton{ID: "generic-unmarked-X", Role: "generic", Signature: []byte(`{}`), Initial: "0", States: []state{{ID: "0"}}, Edges: []edge{}}
	y := automaton{ID: "generic-unmarked-Y", Role: "generic", Signature: []byte(`{}`), Initial: "0", States: []state{{ID: "0"}, {ID: "1"}}, Edges: []edge{{"0", "apply:AB", "1"}}}
	return x, y
}
