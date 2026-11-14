package contractcomposition

import (
	"reflect"
	"testing"
)

func fixture(topology string, mask uint8) Config {
	c := Config{Topology: topology, ServiceMask: mask, Agents: []Agent{{"A", 1, mask&1 != 0}, {"B", 1, mask&2 != 0}, {"C", 1, mask&4 != 0}}}
	for _, pair := range [][2]string{{"A", "B"}, {"B", "C"}, {"C", "A"}} {
		if topology == "chain" && pair[0] == "C" {
			continue
		}
		id := pair[0] + pair[1]
		c.Originals = append(c.Originals, Binding{id, pair[0], 1, 1, pair[1], 1, "increment", "grant:" + pair[0] + "/g1->" + pair[1] + "/g1:counter-" + pair[1], "counter-" + pair[1]})
	}
	return c
}

// Deleting latent AB from Complete would merge these states incorrectly.
func TestHiddenDebtCompletionCollision(t *testing.T) {
	c := fixture("chain", 7)
	s := State{[3]Stage{W, D, ""}}
	paid := State{[3]Stage{D, D, ""}}
	if Complete(c, s) || !Complete(c, paid) {
		t.Fatal("completion must retain hidden AB debt")
	}
	if len(Successors(c, s)) != 1 || Successors(c, s)[0].Action.Original != "AB" {
		t.Fatal("hidden application must stay visible")
	}
}

// A single unpaid-debt bit cannot distinguish an uncommitted from a waiting A.
func TestHiddenWaitReadinessCollision(t *testing.T) {
	c := fixture("cycle", 6)
	u := State{[3]Stage{U, U, W}}
	w := State{[3]Stage{W, U, W}}
	if !ServiceGuard(c, u, "A") || ServiceGuard(c, w, "A") {
		t.Fatal("serial A must read its named outgoing stage")
	}
	for _, stage := range []Stage{W, P} {
		w.Stages[0] = stage
		if !Waiting(c, w, "A") {
			t.Fatalf("%s is waiting", stage)
		}
	}
	if ServiceGuard(c, u, "missing") {
		t.Fatal("absent external endpoint is not a non-waiting owner")
	}
}

func TestConfigRejectsChangedBindingsAndResources(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Config)
	}{
		{"resource alias", func(c *Config) { c.Originals[1].Resource = c.Originals[0].Resource }},
		{"grant", func(c *Config) { c.Originals[0].Grant = "forged" }},
		{"generation", func(c *Config) { c.Originals[0].CallerGeneration = 2 }},
		{"operation", func(c *Config) { c.Originals[0].Operation = 2 }},
		{"duplicate identity", func(c *Config) { c.Originals[1].Original = "AB" }},
		{"owner", func(c *Config) { c.Agents[1] = c.Agents[0] }},
		{"mask", func(c *Config) { c.ServiceMask = 8 }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := fixture("chain", 7)
			tc.mutate(&c)
			if c.Validate() == nil {
				t.Fatal("accepted incompatible contract")
			}
		})
	}
	for _, topology := range []string{"chain", "cycle"} {
		for mask := uint8(0); mask < 8; mask++ {
			c, err := NewConfig(topology, mask)
			if err != nil || !reflect.DeepEqual(c, fixture(topology, mask)) {
				t.Fatalf("fixture %s/%d: %v %#v", topology, mask, err, c)
			}
		}
	}
	if _, err := NewConfig("fanout", 0); err == nil {
		t.Fatal("accepted unsupported grammar")
	}
}

func TestStateAndActionConstructorsRejectUndeclaredRecords(t *testing.T) {
	c := fixture("chain", 7)
	for _, stages := range [][3]Stage{{U, U, U}, {U, "invalid", ""}} {
		if _, err := NewState(c, stages); err == nil {
			t.Errorf("accepted undeclared or invalid stage %v", stages)
		}
	}
	for _, input := range [][2]string{{"apply", "CA"}, {"release", "AB"}, {"consume", "unknown"}} {
		if _, err := NewAction(c, input[0], input[1]); err == nil {
			t.Errorf("accepted undeclared action %v", input)
		}
	}
	a, err := NewAction(c, "apply", "AB")
	if err != nil || a.Binding != c.Originals[0] {
		t.Error("action lacks exact declared binding")
	}
}
