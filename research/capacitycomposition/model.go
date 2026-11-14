// Package capacitycomposition models the fixed Phase 3C.2 shared-capacity contract.
// It is a research relation, not a runtime allocator.
package capacitycomposition

import (
	"fmt"
	cc "github.com/JakeFAU/alan-research/research/contractcomposition"
)

type Root struct {
	ID    string
	Limit int
}
type Membership struct {
	Original string
	Root     string
}
type Config struct {
	Base        cc.Config
	Layout      string
	Policy      string
	Roots       []Root
	Memberships []Membership
}
type State struct {
	Base cc.State
	Held [2]bool
}
type Action = cc.Action
type Step struct {
	Action Action
	State  State
}

func NewConfig(layout, policy string, mask uint8) (Config, error) {
	base, err := cc.NewConfig("chain", mask)
	if err != nil {
		return Config{}, err
	}
	if policy != "early" && policy != "ready-now" && policy != "stable-ready" {
		return Config{}, fmt.Errorf("unknown reservation policy")
	}
	c := Config{Base: base, Layout: layout, Policy: policy}
	switch layout {
	case "shared-1", "shared-2":
		limit := 1
		if layout == "shared-2" {
			limit = 2
		}
		c.Roots = []Root{{"Q", limit}}
		c.Memberships = []Membership{{"AB", "Q"}, {"BC", "Q"}}
	case "separate-1":
		c.Roots = []Root{{"Q-AB", 1}, {"Q-BC", 1}}
		c.Memberships = []Membership{{"AB", "Q-AB"}, {"BC", "Q-BC"}}
	default:
		return Config{}, fmt.Errorf("unknown capacity layout")
	}
	return c, nil
}

// The template is private and never returned or mutated. Validation compares
// exact records without allocating a new configuration at every probe step.
var chainTemplate = func() cc.Config { c, err := cc.NewConfig("chain", 0); must(err); return c }()

func (c Config) Validate() error {
	if c.Policy != "early" && c.Policy != "ready-now" && c.Policy != "stable-ready" {
		return fmt.Errorf("unknown reservation policy")
	}
	if c.Base.Topology != "chain" || c.Base.ServiceMask > 7 || len(c.Base.Agents) != 3 || len(c.Base.Originals) != 2 {
		return fmt.Errorf("wrong fixed-chain shape")
	}
	for i, a := range c.Base.Agents {
		expected := chainTemplate.Agents[i]
		expected.Servicing = c.Base.ServiceMask&(1<<i) != 0
		if a != expected {
			return fmt.Errorf("incompatible owner or policy")
		}
	}
	for i, b := range c.Base.Originals {
		if b != chainTemplate.Originals[i] {
			return fmt.Errorf("incompatible exact original binding")
		}
	}
	if len(c.Memberships) != 2 {
		return fmt.Errorf("wrong membership count")
	}
	switch c.Layout {
	case "shared-1", "shared-2":
		limit := 1
		if c.Layout == "shared-2" {
			limit = 2
		}
		if len(c.Roots) != 1 || c.Roots[0] != (Root{"Q", limit}) || c.Memberships[0] != (Membership{"AB", "Q"}) || c.Memberships[1] != (Membership{"BC", "Q"}) {
			return fmt.Errorf("wrong shared root or memberships")
		}
	case "separate-1":
		if len(c.Roots) != 2 || c.Roots[0] != (Root{"Q-AB", 1}) || c.Roots[1] != (Root{"Q-BC", 1}) || c.Memberships[0] != (Membership{"AB", "Q-AB"}) || c.Memberships[1] != (Membership{"BC", "Q-BC"}) {
			return fmt.Errorf("wrong separate roots or memberships")
		}
	default:
		return fmt.Errorf("unknown capacity layout")
	}
	return nil
}
func Initial(c Config) State { must(c.Validate()); return State{Base: cc.Initial(c.Base)} }
func (s State) Validate(c Config) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if s.Base.Stages[2] != "" {
		return fmt.Errorf("undeclared progress register")
	}
	for _, stage := range s.Base.Stages[:2] {
		if stage != cc.U && stage != cc.W && stage != cc.P && stage != cc.D {
			return fmt.Errorf("invalid progress stage")
		}
	}
	for i, held := range s.Held {
		if held && s.Base.Stages[i] != cc.W {
			return fmt.Errorf("only W may hold capacity")
		}
	}
	for i, n := range Occupancy(c, s) {
		if n > c.Roots[i].Limit {
			return fmt.Errorf("capacity overflow at %s", c.Roots[i].ID)
		}
	}
	return nil
}

// Occupancy derives totals from exact holders and membership. It is not a state replacement.
func Occupancy(c Config, s State) []int {
	out := make([]int, len(c.Roots))
	for i, m := range c.Memberships {
		if s.Held[i] {
			for j, r := range c.Roots {
				if r.ID == m.Root {
					out[j]++
				}
			}
		}
	}
	return out
}
func reservationReady(c Config, s State, i int) bool {
	callee := c.Base.Originals[i].Callee
	switch c.Policy {
	case "early":
		return true
	case "ready-now":
		return cc.ServiceGuard(c.Base, s.Base, callee)
	case "stable-ready":
		for _, a := range c.Base.Agents {
			if a.ID == callee && a.Servicing {
				return true
			}
		}
		for j, b := range c.Base.Originals {
			if b.Caller == callee && s.Base.Stages[j] != cc.D {
				return false
			}
		}
		return true
	}
	return false
}
func Successors(c Config, s State) []Step {
	must(s.Validate(c))
	out := []Step{}
	for _, e := range cc.Successors(c.Base, s.Base) {
		i := slot(e.Action.Original)
		if e.Action.Kind == "apply" && !s.Held[i] {
			continue
		}
		next := s
		next.Base = e.State
		if e.Action.Kind == "apply" {
			next.Held[i] = false
		}
		out = append(out, Step{e.Action, next})
	}
	totals := Occupancy(c, s)
	for i, b := range c.Base.Originals {
		if s.Base.Stages[i] != cc.W || s.Held[i] || !reservationReady(c, s, i) {
			continue
		}
		for j, r := range c.Roots {
			if r.ID == c.Memberships[i].Root && totals[j] < r.Limit {
				next := s
				next.Held[i] = true
				out = append(out, Step{Action{Kind: "reserve-capacity", Original: b.Original, Binding: b}, next})
			}
		}
	}
	return out
}
func Complete(c Config, s State) bool { must(s.Validate(c)); return cc.Complete(c.Base, s.Base) }
func slot(id string) int {
	switch id {
	case "AB":
		return 0
	case "BC":
		return 1
	}
	panic("undeclared original")
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
	panic("undeclared owner")
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
