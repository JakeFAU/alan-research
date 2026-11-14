package contractcomposition

import (
	"fmt"
	"slices"
)

// Component retains separate physical policies and original registers. Ports
// carry the exact expected endpoint descriptor until their owner is joined.
type Component struct {
	Owners    []Agent
	Ports     []Agent
	Originals []Binding
}

func Leaf(c Config, owner string) (Component, error) {
	if err := c.Validate(); err != nil {
		return Component{}, err
	}
	i := agentSlot(owner)
	if i < 0 {
		return Component{}, fmt.Errorf("unknown owner")
	}
	out := Component{Owners: []Agent{c.Agents[i]}, Ports: []Agent{}, Originals: []Binding{}}
	for _, b := range c.Originals {
		if b.Caller == owner {
			out.Originals = append(out.Originals, b)
			out.Ports = append(out.Ports, c.Agents[agentSlot(b.Callee)])
		}
	}
	return out, nil
}

func (c Component) Validate() error {
	owners := map[string]bool{}
	nodes := map[string]Agent{}
	for _, a := range c.Owners {
		if agentSlot(a.ID) < 0 || a.Generation != 1 || owners[a.ID] {
			return fmt.Errorf("invalid or duplicated owner")
		}
		owners[a.ID] = true
		nodes[a.ID] = a
	}
	if len(owners) == 0 {
		return fmt.Errorf("empty owner set")
	}
	for _, a := range c.Ports {
		if _, ok := nodes[a.ID]; ok {
			return fmt.Errorf("duplicate or owned port")
		}
		if agentSlot(a.ID) < 0 || a.Generation != 1 {
			return fmt.Errorf("invalid port")
		}
		nodes[a.ID] = a
	}
	originals := map[string]bool{}
	resources := map[string]bool{}
	outgoing := map[string]bool{}
	needed := map[string]bool{}
	for _, b := range c.Originals {
		if b != fixedBinding(b.Original) || b.Original == "" {
			return fmt.Errorf("incompatible original binding")
		}
		if originals[b.Original] || resources[b.Resource] || outgoing[b.Caller] {
			return fmt.Errorf("duplicate original, resource or outgoing obligation")
		}
		if !owners[b.Caller] {
			return fmt.Errorf("original caller is not owned")
		}
		if _, ok := nodes[b.Callee]; !ok {
			return fmt.Errorf("unrepresented endpoint")
		}
		originals[b.Original] = true
		resources[b.Resource] = true
		outgoing[b.Caller] = true
		if !owners[b.Callee] {
			needed[b.Callee] = true
		}
	}
	if len(needed) != len(c.Ports) {
		return fmt.Errorf("unreferenced port")
	}
	return nil
}

func Join(a, b Component) (Component, error) {
	if err := a.Validate(); err != nil {
		return Component{}, err
	}
	if err := b.Validate(); err != nil {
		return Component{}, err
	}
	out := Component{Owners: []Agent{}, Ports: []Agent{}, Originals: []Binding{}}
	nodes := map[string]Agent{}
	owned := map[string]bool{}
	for _, c := range []Component{a, b} {
		for _, n := range c.Owners {
			if owned[n.ID] {
				return Component{}, fmt.Errorf("duplicated owner %s", n.ID)
			}
			owned[n.ID] = true
			out.Owners = append(out.Owners, n)
		}
		for _, n := range append(slices.Clone(c.Owners), c.Ports...) {
			if old, ok := nodes[n.ID]; ok && old != n {
				return Component{}, fmt.Errorf("incompatible endpoint policy %s", n.ID)
			}
			nodes[n.ID] = n
		}
		out.Originals = append(out.Originals, c.Originals...)
	}
	for id, n := range nodes {
		if !owned[id] {
			out.Ports = append(out.Ports, n)
		}
	}
	slices.SortFunc(out.Owners, func(a, b Agent) int { return agentSlot(a.ID) - agentSlot(b.ID) })
	slices.SortFunc(out.Ports, func(a, b Agent) int { return agentSlot(a.ID) - agentSlot(b.ID) })
	slices.SortFunc(out.Originals, func(a, b Binding) int { return slot(a.Original) - slot(b.Original) })
	if err := out.Validate(); err != nil {
		return Component{}, err
	}
	return out, nil
}

func Close(component Component) (Config, error) {
	if err := component.Validate(); err != nil {
		return Config{}, err
	}
	if len(component.Ports) > 0 {
		return Config{}, fmt.Errorf("unresolved ports")
	}
	topology := "chain"
	if len(component.Originals) == 3 {
		topology = "cycle"
	}
	c := Config{Topology: topology, Agents: slices.Clone(component.Owners), Originals: slices.Clone(component.Originals)}
	for _, a := range c.Agents {
		if a.Servicing {
			c.ServiceMask |= 1 << agentSlot(a.ID)
		}
	}
	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}
