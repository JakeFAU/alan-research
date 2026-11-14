package componentinterfaces

import (
	"fmt"
	pc "github.com/JakeFAU/alan-research/research/capacitycomposition"
)

func histories() ([]history, error) {
	out := []history{}
	for index, roster := range []struct {
		monitor, layout string
		mask            uint8
	}{
		{"commit-first", "shared-1", 0}, {"release-first", "shared-1", 0}, {"commit-first", "shared-2", 0}, {"commit-first", "separate-1", 0}, {"commit-first", "shared-1", 2}, {"release-first-tau", "shared-1", 0}, {"release-first-doomed", "shared-1", 0}, {"release-first-doomed", "shared-1", 0},
	} {
		c, err := pc.NewConfig(roster.layout, "early", roster.mask)
		if err != nil {
			return nil, err
		}
		var m automaton
		for _, a := range monitors(c) {
			if a.ID == roster.monitor {
				m = a
			}
		}
		old, err := pc.NewConcrete(c)
		if err != nil {
			return nil, err
		}
		s := concreteState{old, m.Initial}
		h := history{ID: fmt.Sprintf("history-%d", index+1), Config: c, Monitor: m.ID, Initial: s, Steps: []microStep{}, Cuts: []cut{}}
		g, err := product(plant(c), m)
		if err != nil {
			return nil, err
		}
		recordCut := func(name string) error {
			n, e := graphNode(g, concreteKey(c, s))
			if e != nil {
				return e
			}
			h.Cuts = append(h.Cuts, cut{name, len(h.Steps), n.ID, n.Marked, n.CanComplete})
			return nil
		}
		progress := func(label, to string) error {
			path, e := liftMonitor(c, m, s, label, to)
			if e != nil {
				return e
			}
			for _, step := range path {
				h.Steps = append(h.Steps, step)
				s = step.State
			}
			return nil
		}
		physical := func(kind string) error {
			for _, e := range concreteSuccessors(c, m, s) {
				if string(e.Event.Kind) == kind && e.Event.Actor == "B" {
					h.Steps = append(h.Steps, e)
					s = e.State
					return nil
				}
			}
			return fmt.Errorf("missing physical step %s", kind)
		}
		if err := progress("commit:AB", "0"); err != nil {
			return nil, err
		}
		if err := progress("reserve-capacity:AB", "1"); err != nil {
			return nil, err
		}
		if err := recordCut("after-reserve-AB"); err != nil {
			return nil, err
		}
		if index == 0 {
			if err := progress("commit:BC", "2"); err != nil {
				return nil, err
			}
		}
		if index == 5 {
			if err := physical("crash"); err != nil {
				return nil, err
			}
			if err := recordCut("B-down-before-private"); err != nil {
				return nil, err
			}
			if err := progress("tau", "3"); err != nil {
				return nil, err
			}
			if err := recordCut("private-after-recovery"); err != nil {
				return nil, err
			}
		}
		if index == 7 {
			if err := progress("tau", "3"); err != nil {
				return nil, err
			}
			if err := recordCut("doomed-before-crash"); err != nil {
				return nil, err
			}
			if err := physical("crash"); err != nil {
				return nil, err
			}
			if err := physical("recover"); err != nil {
				return nil, err
			}
			if err := recordCut("doomed-after-recovery"); err != nil {
				return nil, err
			}
		}
		if index != 0 && index != 7 {
			path, e := completionPath(g, concreteKey(c, s))
			if e != nil {
				return nil, e
			}
			for _, arc := range path {
				n, e := graphNode(g, arc.To)
				if e != nil {
					return nil, e
				}
				if e := progress(arc.Label, n.Tuple[1]); e != nil {
					return nil, e
				}
			}
		}
		if err := recordCut("final"); err != nil {
			return nil, err
		}
		last := h.Cuts[len(h.Cuts)-1]
		h.Complete = last.Complete
		h.CanComplete = last.CanComplete
		if err := checkHistory(h); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, nil
}
func completionPath(g graph, start string) ([]arc, error) {
	q := []string{start}
	seen := map[string]bool{start: true}
	parent := map[string]string{}
	via := map[string]arc{}
	for i := 0; i < len(q); i++ {
		n, e := graphNode(g, q[i])
		if e != nil {
			return nil, e
		}
		if n.Marked {
			path := []arc{}
			for id := n.ID; id != start; id = parent[id] {
				path = append(path, via[id])
			}
			for l, r := 0, len(path)-1; l < r; l, r = l+1, r-1 {
				path[l], path[r] = path[r], path[l]
			}
			return path, nil
		}
		for _, arc := range n.Arcs {
			if !seen[arc.To] {
				seen[arc.To] = true
				parent[arc.To] = n.ID
				via[arc.To] = arc
				q = append(q, arc.To)
			}
		}
	}
	return nil, fmt.Errorf("no completion path")
}
