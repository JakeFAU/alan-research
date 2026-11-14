package componentinterfaces

import (
	"fmt"
	pc "github.com/JakeFAU/alan-research/research/capacitycomposition"
	cc "github.com/JakeFAU/alan-research/research/contractcomposition"
)

type concreteState struct {
	Old     pc.ConcreteState
	Monitor string
}
type microStep struct {
	Event   cc.Event
	Private bool
	State   concreteState
}
type cut struct {
	Name        string
	AfterSteps  int
	Node        string
	Complete    bool
	CanComplete bool
}
type history struct {
	ID          string
	Config      pc.Config
	Monitor     string
	Initial     concreteState
	Steps       []microStep
	Cuts        []cut
	Complete    bool
	CanComplete bool
}

func concreteSuccessors(c pc.Config, m automaton, s concreteState) []microStep {
	out := []microStep{}
	for _, e := range pc.ConcreteSuccessors(c, s.Old) {
		if pc.Project(c, s.Old) == pc.Project(c, e.State) {
			out = append(out, microStep{Event: e.Event, State: concreteState{e.State, s.Monitor}})
			continue
		}
		label := string(e.Event.Kind) + ":" + e.Event.Original
		for _, to := range targets(m, s.Monitor, label) {
			out = append(out, microStep{Event: e.Event, State: concreteState{e.State, to}})
		}
	}
	if !s.Old.Base.Agents[1].Down {
		for _, to := range targets(m, s.Monitor, "tau") {
			out = append(out, microStep{Event: cc.Event{Kind: "tau", Actor: "B"}, Private: true, State: concreteState{s.Old, to}})
		}
	}
	return out
}
func concreteKey(c pc.Config, s concreteState) string {
	return key([]string{plantKey(pc.Project(c, s.Old)), s.Monitor})
}
func graphNode(g graph, id string) (node, error) {
	for _, n := range g.Nodes {
		if n.ID == id {
			return n, nil
		}
	}
	return node{}, fmt.Errorf("missing graph state %s", id)
}

// liftMonitor uses inherited custody preparation, checking each resulting step
// against the wrapped relation. Preparation cannot silently advance M.
func liftMonitor(c pc.Config, m automaton, s concreteState, label, to string) ([]microStep, error) {
	out := []microStep{}
	take := func(e microStep) error {
		for _, candidate := range concreteSuccessors(c, m, s) {
			if candidate == e {
				out = append(out, e)
				s = e.State
				return nil
			}
		}
		return fmt.Errorf("illegal wrapped lift step")
	}
	if label == "tau" {
		if s.Old.Base.Agents[1].Down {
			found := false
			for _, e := range concreteSuccessors(c, m, s) {
				if e.Event.Kind == "recover" && e.Event.Actor == "B" {
					if err := take(e); err != nil {
						return nil, err
					}
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("cannot recover B")
			}
		}
		for _, e := range concreteSuccessors(c, m, s) {
			if e.Private && e.State.Monitor == to {
				if err := take(e); err != nil {
					return nil, err
				}
				return out, nil
			}
		}
		return nil, fmt.Errorf("disabled private step")
	}
	var action cc.Action
	found := false
	for _, sym := range alphabet(c) {
		if sym.ID == label {
			action = sym.Action
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("unknown visible label")
	}
	allowed := false
	for _, target := range targets(m, s.Monitor, label) {
		allowed = allowed || target == to
	}
	if !allowed {
		return nil, fmt.Errorf("monitor disables visible action")
	}
	path, err := pc.Lift(c, s.Old, action)
	if err != nil {
		return nil, err
	}
	for _, e := range path {
		pcNext := s.Monitor
		if pc.Project(c, s.Old) != pc.Project(c, e.State) {
			pcNext = to
		}
		if err := take(microStep{Event: e.Event, State: concreteState{e.State, pcNext}}); err != nil {
			return nil, err
		}
	}
	return out, nil
}
func checkHistory(h history) error {
	var m automaton
	found := false
	for _, a := range monitors(h.Config) {
		if a.ID == h.Monitor {
			m = a
			found = true
		}
	}
	if !found {
		return fmt.Errorf("unknown monitor")
	}
	old, err := pc.NewConcrete(h.Config)
	if err != nil {
		return err
	}
	s := concreteState{old, m.Initial}
	if h.Initial != s {
		return fmt.Errorf("noninitial history")
	}
	g, err := product(plant(h.Config), m)
	if err != nil {
		return err
	}
	states := []concreteState{s}
	for i, e := range h.Steps {
		legal := false
		for _, candidate := range concreteSuccessors(h.Config, m, s) {
			if candidate == e {
				legal = true
				break
			}
		}
		if !legal {
			return fmt.Errorf("history %s step %d illegal", h.ID, i)
		}
		s = e.State
		states = append(states, s)
	}
	for _, cut := range h.Cuts {
		if cut.AfterSteps < 0 || cut.AfterSteps >= len(states) {
			return fmt.Errorf("cut outside history")
		}
		id := concreteKey(h.Config, states[cut.AfterSteps])
		n, e := graphNode(g, id)
		if e != nil {
			return e
		}
		if cut.Node != id || cut.Complete != n.Marked || cut.CanComplete != n.CanComplete {
			return fmt.Errorf("wrong cut classification")
		}
	}
	n, err := graphNode(g, concreteKey(h.Config, s))
	if err != nil {
		return err
	}
	if h.Complete != n.Marked || h.CanComplete != n.CanComplete {
		return fmt.Errorf("wrong final classification")
	}
	return nil
}
