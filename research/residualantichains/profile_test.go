package residualantichains

import "testing"

// Catches FULL equality used in place of strict antichain reduction.
func TestStrictReduction(t *testing.T) {
	ms := extraMonitors([]byte(`{}`))
	c := compareProfiles(ms[0], ms[1])
	if c.FULL || !c.MIN || !c.MAX || !c.InitialLanguage.Equal || !c.J {
		t.Fatalf("strict reduction: FULL=%v MIN=%v MAX=%v initial=%v J=%v", c.FULL, c.MIN, c.MAX, c.InitialLanguage.Equal, c.J)
	}
}

// Catches using MIN alone for initial possibility in the generic marked domain.
func TestGenericMINInformationLoss(t *testing.T) {
	x, y := genericMINLoss()
	c := compareProfiles(x, y)
	if !c.MIN || c.InitialLanguage.Equal || c.J {
		t.Fatalf("MIN loses initial information: MIN=%v initial=%v J=%v", c.MIN, c.InitialLanguage.Equal, c.J)
	}
}

// Equivalent duplicates must collapse before strict dominance, otherwise both minima disappear.
func TestDuplicateAndEmptyFamilies(t *testing.T) {
	a := tiny("duplicates", []bool{true, true, true}, edge{From: "0", Label: labels[0], To: "1"}, edge{From: "0", Label: labels[0], To: "2"})
	f := reduce([]string{"1", "2"}, residualIndex(residuals(a, a)))
	if len(f.Classes) != 1 || len(f.MIN) != 1 || len(f.MAX) != 1 || len(f.Classes[0].Members) != 2 {
		t.Fatalf("lost equivalent minima: %+v", f)
	}
	x := tiny("no-edge", []bool{false})
	y := tiny("edge", []bool{false, false}, edge{From: "0", Label: labels[2], To: "1"})
	c := compareProfiles(x, y)
	if c.MIN || c.J || !c.InitialLanguage.Equal {
		t.Fatal("empty family confused with empty language")
	}
	found := false
	for _, p := range c.Pairs {
		if len(p.Left) == 0 && len(p.Right) > 0 {
			found = true
			f := p.LeftToRight[0].Failure
			if f == nil || len(f.H) != 0 || len(f.Outside) != 0 || !f.TargetMiss {
				t.Fatalf("empty source requires H empty: %+v", f)
			}
		}
	}
	if !found {
		t.Fatal("missing empty-family mismatch")
	}
	// A nonempty source containing epsilon needs H={epsilon}, never H empty.
	e := tiny("epsilon", []bool{true})
	d := tiny("dead", []bool{false})
	c = compareProfiles(e, d)
	f2 := c.Pairs[0].LeftToRight[0].Failure
	if f2 == nil || len(f2.H) != 1 || len(f2.H[0]) != 0 || !f2.Hits[0].Hit {
		t.Fatalf("epsilon separator: %+v", f2)
	}
}
func TestMultiwordCoverageDirection(t *testing.T) {
	a := tiny("three-singletons", []bool{false, false, false, true}, edge{From: "0", Label: labels[2], To: "3"}, edge{From: "1", Label: labels[1], To: "3"}, edge{From: "2", Label: labels[4], To: "3"})
	cross := residualIndex(residuals(a, a))
	src := reduce([]string{"0", "1"}, cross)
	dst := reduce([]string{"2"}, cross)
	cs := cover(a, a, src, dst, cross, false)
	if len(cs) != 1 || cs[0].Failure == nil || len(cs[0].Failure.H) != 2 || !cs[0].Failure.TargetMiss {
		t.Fatal("multiword separator missing")
	}
	for _, h := range cs[0].Failure.Hits {
		if !h.Hit {
			t.Fatal("separator does not hit source")
		}
	}
	// Empty is a subset of epsilon, so dead -> epsilon coverage holds, not reverse.
	empty := tiny("dead", []bool{false})
	epsilon := tiny("epsilon", []bool{true})
	c := compareProfiles(empty, epsilon)
	if len(c.Pairs[0].LeftToRight[0].Sources) != 1 || c.Pairs[0].RightToLeft[0].Failure == nil {
		t.Fatal("coverage direction reversed")
	}
}
