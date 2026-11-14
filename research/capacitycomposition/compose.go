package capacitycomposition

import (
	"fmt"
	cc "github.com/JakeFAU/alan-research/research/contractcomposition"
	"reflect"
	"slices"
)

type Component struct {
	Base        cc.Component
	Layout      string
	Policy      string
	Roots       []Root
	Memberships []Membership
}

func Leaf(c Config, owner string) (Component, error) {
	if err := c.Validate(); err != nil {
		return Component{}, err
	}
	base, err := cc.Leaf(c.Base, owner)
	if err != nil {
		return Component{}, err
	}
	out := Component{Base: base, Layout: c.Layout, Policy: c.Policy, Roots: []Root{}, Memberships: []Membership{}}
	for i, b := range c.Base.Originals {
		if b.Caller == owner {
			out.Memberships = append(out.Memberships, c.Memberships[i])
			for _, r := range c.Roots {
				if r.ID == c.Memberships[i].Root {
					out.Roots = append(out.Roots, r)
				}
			}
		}
	}
	return out, nil
}
func (c Component) Validate() error {
	if err := c.Base.Validate(); err != nil {
		return err
	}
	full, err := NewConfig(c.Layout, c.Policy, 0)
	if err != nil {
		return err
	}
	// Each owned A/B must carry its outgoing original; C has none. The old
	// open-component validator alone permits omission, so enforce the fixed chain.
	expectedOriginals := []cc.Binding{}
	members := []Membership{}
	roots := []Root{}
	for i, b := range full.Base.Originals {
		for _, owner := range c.Base.Owners {
			if b.Caller == owner.ID {
				expectedOriginals = append(expectedOriginals, b)
				members = append(members, full.Memberships[i])
			}
		}
	}
	for _, r := range full.Roots {
		for _, m := range members {
			if m.Root == r.ID {
				roots = append(roots, r)
				break
			}
		}
	}
	if !slices.Equal(c.Base.Originals, expectedOriginals) || !slices.Equal(c.Memberships, members) || !slices.Equal(c.Roots, roots) {
		return fmt.Errorf("missing, extra, conflicting or reordered original/root membership")
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
	if a.Layout != b.Layout || a.Policy != b.Policy {
		return Component{}, fmt.Errorf("incompatible layout or policy")
	}
	base, err := cc.Join(a.Base, b.Base)
	if err != nil {
		return Component{}, err
	}
	out := Component{Base: base, Layout: a.Layout, Policy: a.Policy, Roots: slices.Clone(a.Roots), Memberships: append(slices.Clone(a.Memberships), b.Memberships...)}
	for _, r := range b.Roots {
		found := false
		for _, old := range out.Roots {
			if old.ID == r.ID {
				if old != r {
					return Component{}, fmt.Errorf("conflicting shared root")
				}
				found = true
			}
		}
		if !found {
			out.Roots = append(out.Roots, r)
		}
	}
	slices.SortFunc(out.Roots, func(a, b Root) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	})
	slices.SortFunc(out.Memberships, func(a, b Membership) int { return slot(a.Original) - slot(b.Original) })
	if err := out.Validate(); err != nil {
		return Component{}, err
	}
	return out, nil
}
func Close(component Component) (Config, error) {
	if err := component.Validate(); err != nil {
		return Config{}, err
	}
	base, err := cc.Close(component.Base)
	if err != nil {
		return Config{}, err
	}
	c := Config{Base: base, Layout: component.Layout, Policy: component.Policy, Roots: slices.Clone(component.Roots), Memberships: slices.Clone(component.Memberships)}
	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}

type Grouping struct {
	ConfigID    string
	A           Component
	B           Component
	C           Component
	AB          Component
	BC          Component
	Left        Component
	Right       Component
	Flat        Config
	LeftClosed  Config
	RightClosed Config
	Equal       bool
}

func regroup(c Config) (Grouping, error) {
	a, err := Leaf(c, "A")
	if err != nil {
		return Grouping{}, err
	}
	b, err := Leaf(c, "B")
	if err != nil {
		return Grouping{}, err
	}
	d, err := Leaf(c, "C")
	if err != nil {
		return Grouping{}, err
	}
	ab, err := Join(a, b)
	if err != nil {
		return Grouping{}, err
	}
	bc, err := Join(b, d)
	if err != nil {
		return Grouping{}, err
	}
	left, err := Join(ab, d)
	if err != nil {
		return Grouping{}, err
	}
	right, err := Join(a, bc)
	if err != nil {
		return Grouping{}, err
	}
	lc, err := Close(left)
	if err != nil {
		return Grouping{}, err
	}
	rc, err := Close(right)
	if err != nil {
		return Grouping{}, err
	}
	equal := reflect.DeepEqual(left, right) && reflect.DeepEqual(c, lc) && reflect.DeepEqual(c, rc) && reflect.DeepEqual(Explore(lc), Explore(rc))
	return Grouping{configID(c), a, b, d, ab, bc, left, right, c, lc, rc, equal}, nil
}
func configID(c Config) string {
	return fmt.Sprintf("%s/%s/%d", c.Layout, c.Policy, c.Base.ServiceMask)
}
