package componentinterfaces

import "testing"

// These cases catch epsilon erasure, deterministic-only comparison and depth cutoffs.
func TestLanguageComparison(t *testing.T) {
	c := coreConfig(t)
	ms := monitors(c)
	for i := 0; i < len(ms); i++ {
		for j := i + 1; j < len(ms); j++ {
			got := compareLanguages(ms[i], ms[j])
			want := i >= 2
			if got.Equal != want {
				t.Fatalf("%s/%s equal=%v", ms[i].ID, ms[j].ID, got.Equal)
			}
			if !got.Equal && (accepts(ms[i], got.Word) == accepts(ms[j], got.Word)) {
				t.Fatal("invalid separator")
			}
		}
	}
	base := automaton{ID: "test", Role: "test", Signature: c, Initial: "0", States: []state{{ID: "0", Marked: true}}, Edges: []edge{}}
	t.Run("epsilon acceptance", func(t *testing.T) {
		b := base
		b.States = []state{{ID: "0"}}
		r := compareLanguages(base, b)
		if r.Equal || len(r.Word) != 0 {
			t.Fatal(r)
		}
	})
	t.Run("tau cycle and nondeterminism", func(t *testing.T) {
		b := base
		b.States = []state{{ID: "0"}, {ID: "1", Marked: true}, {ID: "2", Marked: true}}
		b.Edges = []edge{{"0", "tau", "1"}, {"1", "tau", "0"}, {"0", "tau", "2"}}
		if !compareLanguages(base, b).Equal {
			t.Fatal("epsilon branches differ")
		}
	})
	t.Run("empty subset", func(t *testing.T) {
		b := base
		b.Edges = []edge{{"0", "apply:AB", "0"}}
		r := compareLanguages(base, b)
		if r.Equal || len(r.Word) != 1 {
			t.Fatal(r)
		}
	})
	t.Run("beyond short samples", func(t *testing.T) {
		a := base
		a.States = []state{{ID: "0"}, {ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5", Marked: true}}
		a.Edges = []edge{{"0", "commit:AB", "1"}, {"1", "commit:AB", "2"}, {"2", "commit:AB", "3"}, {"3", "commit:AB", "4"}, {"4", "commit:AB", "5"}}
		b := a
		b.States = append([]state(nil), a.States...)
		b.States[5].Marked = false
		r := compareLanguages(a, b)
		if r.Equal || len(r.Word) != 5 {
			t.Fatal(r)
		}
	})
}

func TestNondeterministicVisibleBranches(t *testing.T) {
	c := coreConfig(t)
	a := automaton{ID: "a", Signature: c, Initial: "0", States: []state{{ID: "0"}, {ID: "1", Marked: true}}, Edges: []edge{{"0", "apply:AB", "1"}}}
	b := automaton{ID: "b", Signature: c, Initial: "0", States: []state{{ID: "0"}, {ID: "1", Marked: true}, {ID: "2", Marked: true}}, Edges: []edge{{"0", "apply:AB", "1"}, {"0", "apply:AB", "2"}}}
	if !compareLanguages(a, b).Equal {
		t.Fatal("equivalent nondeterministic branches rejected")
	}
}
