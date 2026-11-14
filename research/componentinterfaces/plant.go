package componentinterfaces

import (
	"fmt"
	pc "github.com/JakeFAU/alan-research/research/capacitycomposition"
)

func plantID(c pc.Config) string {
	return fmt.Sprintf("%s/%s/mask%d", c.Policy, c.Layout, c.Base.ServiceMask)
}
func plantKey(s pc.State) string {
	return fmt.Sprintf("%s,%s,%t,%t", s.Base.Stages[0], s.Base.Stages[1], s.Held[0], s.Held[1])
}
func plant(c pc.Config) automaton {
	a := automaton{ID: plantID(c), Role: "K", Signature: cloneConfig(c), Initial: plantKey(pc.Initial(c)), States: []state{}, Edges: []edge{}}
	queue := []pc.State{pc.Initial(c)}
	seen := map[pc.State]bool{queue[0]: true}
	for i := 0; i < len(queue); i++ {
		s := queue[i]
		id := plantKey(s)
		a.States = append(a.States, state{ID: id, Marked: pc.Complete(c, s)})
		for _, e := range pc.Successors(c, s) {
			a.Edges = append(a.Edges, edge{id, e.Action.Kind + ":" + e.Action.Original, plantKey(e.State)})
			if !seen[e.State] {
				seen[e.State] = true
				queue = append(queue, e.State)
			}
		}
	}
	return a
}
