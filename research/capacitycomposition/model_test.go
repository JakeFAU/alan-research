package capacitycomposition

import (
	cc "github.com/JakeFAU/alan-research/research/contractcomposition"
	"reflect"
	"testing"
)

func config(t *testing.T, layout, policy string, mask uint8) Config {
	t.Helper()
	c, err := NewConfig(layout, policy, mask)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
func advance(t *testing.T, c Config, s State, id, kind string) State {
	t.Helper()
	for _, e := range Successors(c, s) {
		if e.Action.Original == id && e.Action.Kind == kind {
			return e.State
		}
	}
	t.Fatalf("missing %s %s from %+v", id, kind, s)
	return State{}
}
func canComplete(t *testing.T, c Config, s State) bool {
	t.Helper()
	for _, n := range Explore(c).Nodes {
		if n.State == s {
			return n.CanComplete
		}
	}
	t.Fatalf("state not reachable: %+v", s)
	return false
}

// These tests catch copied root limits, identity erasure, stale readiness and
// accidental old-q lifting, rather than merely checking configuration spelling.
func TestExactRootAccounting(t *testing.T) {
	c := config(t, "shared-1", "early", 0)
	if !reflect.DeepEqual(c.Roots, []Root{{"Q", 1}}) || !reflect.DeepEqual(c.Memberships, []Membership{{"AB", "Q"}, {"BC", "Q"}}) {
		t.Fatal("shared root is not one exact retained record")
	}
	s := advance(t, c, Initial(c), "AB", "commit")
	s = advance(t, c, s, "AB", "reserve-capacity")
	s = advance(t, c, s, "BC", "commit")
	for _, e := range Successors(c, s) {
		if e.Action.Kind == "reserve-capacity" {
			t.Fatal("over-allocated shared root")
		}
	}
	s.Held[1] = true
	if s.Validate(c) == nil {
		t.Fatal("accepted occupancy overflow")
	}
	c.Roots[0].Limit = 0
	if c.Validate() == nil {
		t.Fatal("coerced zero capacity")
	}
}
func TestAggregateEqualDifferentCompletion(t *testing.T) {
	c := config(t, "shared-1", "early", 0)
	s := advance(t, c, Initial(c), "AB", "commit")
	s = advance(t, c, s, "BC", "commit")
	a := advance(t, c, s, "AB", "reserve-capacity")
	b := advance(t, c, s, "BC", "reserve-capacity")
	if a.Base != b.Base || !reflect.DeepEqual(Occupancy(c, a), Occupancy(c, b)) || a.Held == b.Held {
		t.Fatal("not an aggregate collision")
	}
	if canComplete(t, c, a) || !canComplete(t, c, b) {
		t.Fatal("holder identity did not distinguish completion")
	}
}
func TestReadyNowInvalidatedByCommit(t *testing.T) {
	c := config(t, "shared-1", "ready-now", 0)
	s := advance(t, c, Initial(c), "AB", "commit")
	s = advance(t, c, s, "AB", "reserve-capacity")
	if !cc.ServiceGuard(c.Base, s.Base, "B") {
		t.Fatal("guard false at reservation")
	}
	s = advance(t, c, s, "BC", "commit")
	if cc.ServiceGuard(c.Base, s.Base, "B") || canComplete(t, c, s) {
		t.Fatal("commit must invalidate readiness")
	}
}
func TestStableReadyAndControlsEveryReachableCompletion(t *testing.T) {
	for _, layout := range []string{"shared-1", "shared-2", "separate-1"} {
		for _, policy := range []string{"early", "ready-now", "stable-ready"} {
			for mask := uint8(0); mask < 8; mask++ {
				c := config(t, layout, policy, mask)
				g := Explore(c)
				if len(g.Nodes) == 0 {
					t.Fatal("empty exploration")
				}
				for _, n := range g.Nodes {
					if err := n.State.Validate(c); err != nil {
						t.Fatal(err)
					}
				}
				if policy == "stable-ready" || layout != "shared-1" || mask&2 != 0 {
					if g.CannotComplete != 0 {
						t.Fatalf("control %s/%s/%d has %d failures", layout, policy, mask, g.CannotComplete)
					}
				}
			}
		}
	}
}
func TestOldABApplyTraceExcluded(t *testing.T) {
	early := config(t, "shared-1", "early", 0)
	stable := config(t, "shared-1", "stable-ready", 0)
	s := advance(t, early, Initial(early), "AB", "commit")
	held := advance(t, early, s, "AB", "reserve-capacity")
	applied := advance(t, early, held, "AB", "apply")
	if applied.Base.Stages[1] != cc.U || applied.Held[0] {
		t.Fatal("apply must release atomically while BC remains U")
	}
	for _, e := range Successors(stable, s) {
		if e.Action.Original == "AB" {
			t.Fatal("stable-ready reproduced excluded old trace")
		}
	}
}

func TestMalformedConfigurationsAndStatesRejected(t *testing.T) {
	for _, bad := range []struct {
		layout, policy string
		mask           uint8
	}{{"shared-0", "early", 0}, {"shared-1", "eager", 0}, {"shared-1", "early", 8}} {
		if _, err := NewConfig(bad.layout, bad.policy, bad.mask); err == nil {
			t.Fatal("accepted unknown domain", bad)
		}
	}
	for _, mutate := range []func(*Config){
		func(c *Config) { c.Base.Topology = "cycle" }, func(c *Config) { c.Base.Agents[1].Generation = 2 }, func(c *Config) { c.Base.Agents[1].Servicing = true },
		func(c *Config) { c.Base.Originals[0].Grant = "forged" }, func(c *Config) { c.Memberships[0].Original = "BC" },
		func(c *Config) { c.Roots = append(c.Roots, c.Roots[0]) }, func(c *Config) { c.Roots[0].Limit = -1 }, func(c *Config) { c.Roots[0].Limit = 2 },
	} {
		c := config(t, "shared-1", "early", 0)
		mutate(&c)
		if c.Validate() == nil {
			t.Fatal("accepted malformed config", c)
		}
	}
	c := config(t, "shared-1", "early", 0)
	for _, s := range []State{
		{Base: cc.State{Stages: [3]cc.Stage{cc.U, cc.U, cc.U}}},
		{Base: cc.State{Stages: [3]cc.Stage{"bad", cc.U, ""}}},
		{Base: cc.State{Stages: [3]cc.Stage{cc.P, cc.U, ""}}, Held: [2]bool{true, false}},
	} {
		if s.Validate(c) == nil {
			t.Fatal("accepted malformed state", s)
		}
	}
}
