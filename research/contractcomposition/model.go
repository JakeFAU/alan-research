// Package contractcomposition models the fixed Phase 3C.1 progress contract.
// It is a research model, not a runtime or an arbitrary-topology theorem.
package contractcomposition

import "fmt"

type Stage string

const (
	U Stage = "U"
	W Stage = "W"
	P Stage = "P"
	D Stage = "D"
)

type Agent struct {
	ID         string
	Generation uint8
	Servicing  bool
}

type Binding struct {
	Original         string
	Caller           string
	CallerGeneration uint8
	Operation        uint8
	Callee           string
	CalleeGeneration uint8
	Request          string
	Grant            string
	Resource         string
}

type Config struct {
	Topology    string
	ServiceMask uint8
	Agents      []Agent
	Originals   []Binding
}

// Stages has fixed named slots AB, BC, CA; absent CA in a chain is "".
// No progress register disappears when its endpoints become internal.
type State struct{ Stages [3]Stage }
type Action struct {
	Kind     string
	Original string
	Binding  Binding
}
type Step struct {
	Action Action
	State  State
}

func fixedBinding(id string) Binding {
	i := slot(id)
	if i < 0 {
		return Binding{}
	}
	caller, callee := []string{"A", "B", "C"}[i], []string{"B", "C", "A"}[i]
	return Binding{Original: id, Caller: caller, CallerGeneration: 1, Operation: 1, Callee: callee, CalleeGeneration: 1,
		Request: "increment", Grant: "grant:" + caller + "/g1->" + callee + "/g1:counter-" + callee, Resource: "counter-" + callee}
}

func NewConfig(topology string, mask uint8) (Config, error) {
	if (topology != "chain" && topology != "cycle") || mask > 7 {
		return Config{}, fmt.Errorf("unsupported topology or service mask")
	}
	c := Config{Topology: topology, ServiceMask: mask, Agents: []Agent{}, Originals: []Binding{}}
	for i, id := range []string{"A", "B", "C"} {
		c.Agents = append(c.Agents, Agent{id, 1, mask&(1<<i) != 0})
	}
	for _, id := range []string{"AB", "BC", "CA"} {
		if id != "CA" || topology == "cycle" {
			c.Originals = append(c.Originals, fixedBinding(id))
		}
	}
	return c, nil
}

// Validate rejects deviations from the declared fixed identities and grammar.
func (c Config) Validate() error {
	expected, err := NewConfig(c.Topology, c.ServiceMask)
	if err != nil {
		return err
	}
	if len(c.Agents) != 3 || len(c.Originals) != len(expected.Originals) {
		return fmt.Errorf("wrong owners or original count")
	}
	for i, a := range c.Agents {
		if a != expected.Agents[i] {
			return fmt.Errorf("incompatible owner or policy at %d", i)
		}
	}
	for i, b := range c.Originals {
		if b != expected.Originals[i] {
			return fmt.Errorf("incompatible original binding at %d", i)
		}
	}
	return nil
}

func Initial(c Config) State {
	s := State{}
	for _, b := range c.Originals {
		s.Stages[slot(b.Original)] = U
	}
	return s
}
func (s State) Validate(c Config) error {
	if err := c.Validate(); err != nil {
		return err
	}
	for i, v := range s.Stages {
		if i >= len(c.Originals) {
			if v != "" {
				return fmt.Errorf("undeclared progress register")
			}
			continue
		}
		if v != U && v != W && v != P && v != D {
			return fmt.Errorf("invalid stage %q", v)
		}
	}
	return nil
}
func Complete(c Config, s State) bool {
	for _, b := range c.Originals {
		if s.Stages[slot(b.Original)] != D {
			return false
		}
	}
	return true
}
func Waiting(c Config, s State, agent string) bool {
	for _, b := range c.Originals {
		if b.Caller == agent {
			v := s.Stages[slot(b.Original)]
			return v == W || v == P
		}
	}
	return false
}

// ServiceGuard is policy eligibility, not admission or physical readiness.
// An unresolved external port is never interpreted as an idle agent.
func ServiceGuard(c Config, s State, agent string) bool {
	for _, a := range c.Agents {
		if a.ID == agent {
			return a.Servicing || !Waiting(c, s, agent)
		}
	}
	return false
}
func Successors(c Config, s State) []Step {
	out := []Step{}
	for _, b := range c.Originals {
		i := slot(b.Original)
		next := s
		kind := ""
		switch s.Stages[i] {
		case U:
			next.Stages[i] = W
			kind = "commit"
		case W:
			if ServiceGuard(c, s, b.Callee) {
				next.Stages[i] = P
				kind = "apply"
			}
		case P:
			next.Stages[i] = D
			kind = "consume"
		}
		if kind != "" {
			out = append(out, Step{Action{kind, b.Original, b}, next})
		}
	}
	return out
}
func slot(id string) int {
	switch id {
	case "AB":
		return 0
	case "BC":
		return 1
	case "CA":
		return 2
	}
	return -1
}
func agentSlot(id string) int {
	switch id {
	case "A":
		return 0
	case "B":
		return 1
	case "C":
		return 2
	}
	return -1
}

// NewState validates the exact declared register set, not reachability.
func NewState(c Config, stages [3]Stage) (State, error) {
	s := State{Stages: stages}
	if err := s.Validate(c); err != nil {
		return State{}, err
	}
	return s, nil
}

// NewAction binds a progress kind to one declared original. Enablement still
// depends on the source state and is checked separately by Successors or Lift.
func NewAction(c Config, kind, original string) (Action, error) {
	if err := c.Validate(); err != nil {
		return Action{}, err
	}
	if kind != "commit" && kind != "apply" && kind != "consume" {
		return Action{}, fmt.Errorf("unknown progress kind")
	}
	for _, b := range c.Originals {
		if b.Original == original {
			return Action{kind, original, b}, nil
		}
	}
	return Action{}, fmt.Errorf("undeclared original")
}
