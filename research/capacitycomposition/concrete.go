package capacitycomposition

import (
	"fmt"
	cc "github.com/JakeFAU/alan-research/research/contractcomposition"
	"github.com/JakeFAU/alan-research/research/twoagent"
)

type ConcreteState struct {
	Base cc.ConcreteState
	Held [2]bool
}
type MicroStep struct {
	Event cc.Event
	State ConcreteState
}

func NewConcrete(c Config) (ConcreteState, error) {
	if err := c.Validate(); err != nil {
		return ConcreteState{}, err
	}
	b, e := cc.NewConcrete(c.Base)
	return ConcreteState{Base: b}, e
}
func Project(c Config, s ConcreteState) State {
	return State{Base: cc.Project(c.Base, s.Base), Held: s.Held}
}

// Validate checks retained premises, not reachability or oracle authenticity.
// Holds require actual admission; no operational check reads WorldResult.
func (s ConcreteState) Validate(c Config) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if err := s.Base.Validate(c.Base); err != nil {
		return err
	}
	if err := Project(c, s).Validate(c); err != nil {
		return err
	}
	for i, held := range s.Held {
		if held && !s.Base.Originals[i].Admitted {
			return fmt.Errorf("holder lacks admitted request")
		}
	}
	return nil
}

// ConcreteSuccessors preserves every old step except application without a
// hold. Actual application and result retention atomically clear that hold.
func ConcreteSuccessors(c Config, s ConcreteState) []MicroStep {
	must(s.Validate(c))
	out := []MicroStep{}
	for _, e := range cc.ConcreteSuccessors(c.Base, s.Base) {
		next := ConcreteState{Base: e.State, Held: s.Held}
		if e.Event.Kind == twoagent.Apply {
			i := slot(e.Event.Original)
			if !s.Held[i] {
				continue
			}
			next.Held[i] = false
		}
		out = append(out, MicroStep{e.Event, next})
	}
	for _, e := range Successors(c, Project(c, s)) {
		if e.Action.Kind != "reserve-capacity" {
			continue
		}
		i := slot(e.Action.Original)
		b := e.Action.Binding
		if !s.Base.Originals[i].Admitted || s.Base.Agents[agentSlot(b.Callee)].Down {
			continue
		}
		next := s
		next.Held = e.State.Held
		out = append(out, MicroStep{cc.Event{Kind: "reserve-capacity", Actor: b.Callee, Original: b.Original, Binding: b}, next})
	}
	return out
}

// Lift targets the new resource-aware relation. Preparation neither advances
// another original nor acquires/releases any other hold. It uses real custody.
func Lift(c Config, s ConcreteState, a Action) ([]MicroStep, error) {
	if err := s.Validate(c); err != nil {
		return nil, err
	}
	enabled := false
	for _, e := range Successors(c, Project(c, s)) {
		if e.Action == a {
			enabled = true
		}
	}
	if !enabled {
		return nil, fmt.Errorf("resource-aware action disabled or forged")
	}
	path := []MicroStep{}
	take := func(id string, kind twoagent.Kind) error {
		for _, e := range ConcreteSuccessors(c, s) {
			if e.Event.Kind == kind && !e.Event.Duplicate && (e.Event.Original == id || e.Event.Actor == id) {
				path = append(path, e)
				s = e.State
				return nil
			}
		}
		return fmt.Errorf("no legal %s for %s", kind, id)
	}
	recoverOwner := func(id string) error {
		if s.Base.Agents[agentSlot(id)].Down {
			return take(id, twoagent.Recover)
		}
		return nil
	}
	b := a.Binding
	i := slot(a.Original)
	switch a.Kind {
	case "commit":
		if err := recoverOwner(b.Caller); err != nil {
			return nil, err
		}
		if err := take(b.Original, twoagent.Commit); err != nil {
			return nil, err
		}
	case "reserve-capacity":
		if err := recoverOwner(b.Callee); err != nil {
			return nil, err
		}
		if !s.Base.Originals[i].Admitted {
			if !s.Base.Originals[i].RequestChannel {
				if err := recoverOwner(b.Caller); err != nil {
					return nil, err
				}
				if err := take(b.Original, twoagent.SendRequest); err != nil {
					return nil, err
				}
			}
			if err := take(b.Original, twoagent.Admit); err != nil {
				return nil, err
			}
		}
		if err := take(b.Original, "reserve-capacity"); err != nil {
			return nil, err
		}
	case "apply":
		if err := recoverOwner(b.Callee); err != nil {
			return nil, err
		}
		if err := take(b.Original, twoagent.Apply); err != nil {
			return nil, err
		}
	case "consume":
		if err := recoverOwner(b.Caller); err != nil {
			return nil, err
		}
		if s.Base.Originals[i].Recorded == "" {
			if s.Base.Originals[i].Volatile == "" {
				if s.Base.Originals[i].ResultChannel == "" {
					if err := recoverOwner(b.Callee); err != nil {
						return nil, err
					}
					if err := take(b.Original, twoagent.SendResult); err != nil {
						return nil, err
					}
				}
				if err := take(b.Original, twoagent.ReceiveResult); err != nil {
					return nil, err
				}
			}
			if err := take(b.Original, twoagent.Record); err != nil {
				return nil, err
			}
		}
		if err := take(b.Original, twoagent.Consume); err != nil {
			return nil, err
		}
	}
	return path, nil
}
