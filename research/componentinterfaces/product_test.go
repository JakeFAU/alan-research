package componentinterfaces

import "testing"

func TestProductsAndGrouping(t *testing.T) {
	c := coreConfig(t)
	k := plant(c)
	ms := monitors(c)
	want := []bool{true, false, true, true, true}
	for i, m := range ms {
		g, e := product(k, m)
		if e != nil {
			t.Fatal(e)
		}
		if g.Nodes[0].CanComplete != want[i] {
			t.Fatalf("%s initial completion=%v", m.ID, g.Nodes[0].CanComplete)
		}
		if m.ID == "release-first-doomed" && g.Nonblocking {
			t.Fatal("doomed branch was lost")
		}
	}
	ds := controllers(c)
	if len(ds) != 324 {
		t.Fatal(len(ds))
	}
	for _, d := range ds {
		for _, m := range ms {
			flat, e := product(k, m, d)
			if e != nil {
				t.Fatal(e)
			}
			group, e := checkGrouping(k, m, d, flat)
			if e != nil || !group.Equal {
				t.Fatalf("grouping %v", e)
			}
		}
	}
}
