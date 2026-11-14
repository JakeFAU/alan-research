package componentinterfaces

import (
	"fmt"
	pc "github.com/JakeFAU/alan-research/research/capacitycomposition"
	cc "github.com/JakeFAU/alan-research/research/contractcomposition"
)

func monitors(c pc.Config) []automaton {
	out := []automaton{}
	for i, name := range []string{"free", "commit-first", "release-first", "release-first-tau", "release-first-doomed"} {
		a := automaton{ID: name, Role: "M", Signature: cloneConfig(c), Initial: "0", States: []state{}, Edges: []edge{}}
		n := 3
		if i == 0 {
			n = 1
		}
		if i >= 3 {
			n = 4
		}
		for j := 0; j < n; j++ {
			id := fmt.Sprint(j)
			a.States = append(a.States, state{ID: id, Marked: true})
			for _, s := range alphabet(c) {
				if s.ID != "reserve-capacity:AB" && s.ID != "apply:AB" && s.ID != "commit:BC" {
					a.Edges = append(a.Edges, edge{id, s.ID, id})
				}
			}
		}
		loop := "2"
		if i == 0 {
			loop = "0"
		}
		for _, label := range []string{"reserve-capacity:AB", "apply:AB", "commit:BC"} {
			a.Edges = append(a.Edges, edge{loop, label, loop})
		}
		if i > 0 {
			a.Edges = append(a.Edges, edge{"0", "reserve-capacity:AB", "1"})
		}
		switch i {
		case 1:
			a.Edges = append(a.Edges, edge{"1", "commit:BC", "2"})
		case 2, 4:
			a.Edges = append(a.Edges, edge{"1", "apply:AB", "2"})
		case 3:
			a.Edges = append(a.Edges, edge{"1", "tau", "3"}, edge{"3", "apply:AB", "2"})
		}
		if i == 4 {
			a.Edges = append(a.Edges, edge{"1", "tau", "3"})
		}
		out = append(out, a)
	}
	return out
}
func cloneConfig(c pc.Config) pc.Config {
	c.Base.Agents = append([]cc.Agent(nil), c.Base.Agents...)
	c.Base.Originals = append([]cc.Binding(nil), c.Base.Originals...)
	c.Roots = append([]pc.Root(nil), c.Roots...)
	c.Memberships = append([]pc.Membership(nil), c.Memberships...)
	return c
}
func controllers(c pc.Config) []automaton {
	out := []automaton{}
	for code := 0; code < 81; code++ {
		for mask := 0; mask < 4; mask++ {
			a := automaton{ID: fmt.Sprintf("D-%02d-%d", code, mask), Role: "D", Signature: cloneConfig(c), Initial: "0", States: []state{{ID: "0", Marked: mask&1 != 0}, {ID: "1", Marked: mask&2 != 0}}, Edges: []edge{}}
			digits := code
			for _, from := range []string{"0", "1"} {
				for _, s := range alphabet(c) {
					if s.ID != "commit:AB" && s.ID != "apply:BC" {
						a.Edges = append(a.Edges, edge{from, s.ID, from})
					}
				}
				for _, label := range []string{"commit:AB", "apply:BC"} {
					v := digits % 3
					digits /= 3
					if v > 0 {
						a.Edges = append(a.Edges, edge{from, label, fmt.Sprint(v - 1)})
					}
				}
			}
			out = append(out, a)
		}
	}
	return out
}
