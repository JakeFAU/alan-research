package residualantichains

import (
	"fmt"
	"testing"
)

func TestFiniteSetAlgebraExhaustive(t *testing.T) {
	// Universe {epsilon,r,a}; all 8 languages and all 256 finite families.
	// Expected semantics are direct bit-set intersection, independent of the
	// epsilon-NFA subset exploration and reduction implementation under test.
	universe := [][]string{{}, {labels[1]}, {labels[2]}}
	all := automaton{ID: "finite-universe", Role: "generic", Signature: []byte(`{}`), Initial: "0", States: []state{}, Edges: []edge{}}
	for lang := 0; lang < 8; lang++ {
		root := fmt.Sprint(lang)
		all.States = append(all.States, state{ID: root, Marked: lang&1 != 0})
		for bit, label := range []string{labels[1], labels[2]} {
			if lang&(1<<(bit+1)) != 0 {
				all.Edges = append(all.Edges, edge{From: root, Label: label, To: "end"})
			}
		}
	}
	all.States = append(all.States, state{ID: "end", Marked: true})
	idx := residualIndex(residuals(all, all))
	families := make([]family, 256)
	mins := make([]int, 256)
	for f := 0; f < 256; f++ {
		roots := []string{}
		for lang := 0; lang < 8; lang++ {
			if f&(1<<lang) != 0 {
				roots = append(roots, fmt.Sprint(lang))
			}
		}
		families[f] = reduce(roots, idx)
		for _, root := range families[f].MIN {
			var lang int
			if _, err := fmt.Sscan(root, &lang); err != nil {
				t.Fatal(err)
			}
			mins[f] |= 1 << lang
		}
	}
	hit := func(f, h int) bool {
		for l := 0; l < 8; l++ {
			if f&(1<<l) != 0 && l&h == 0 {
				return false
			}
		}
		return true
	}
	tensor := func(f, g int) int {
		out := 0
		for l := 0; l < 8; l++ {
			for r := 0; r < 8; r++ {
				if f&(1<<l) != 0 && g&(1<<r) != 0 {
					out |= 1 << (l & r)
				}
			}
		}
		return out
	}
	for f := 0; f < 256; f++ {
		for h := 0; h < 8; h++ {
			if hit(f, h) != hit(mins[f], h) {
				t.Fatalf("hit preservation F=%d H=%d", f, h)
			}
		}
		for g := 0; g < 256; g++ {
			if mins[tensor(f, g)] != mins[tensor(mins[f], mins[g])] {
				t.Fatalf("product reduction F=%d G=%d", f, g)
			}
			implication := true
			for h := 0; h < 8; h++ {
				if hit(f, h) && !hit(g, h) {
					implication = false
				}
			}
			cov := cover(all, all, families[f], families[g], idx, false)
			covered := true
			for _, c := range cov {
				if len(c.Sources) == 0 {
					covered = false
					if c.Failure == nil || !c.Failure.TargetMiss {
						t.Fatal("missing failed-target witness")
					}
					for _, v := range c.Failure.Hits {
						if !v.Hit {
							t.Fatal("invalid source hit")
						}
					}
				}
			}
			if implication != covered {
				t.Fatalf("local implication F=%d G=%d", f, g)
			}
		}
	}
	// All 64 finite-language intersections are realized by actual synchronized
	// trie products, including empty language and epsilon.
	for l := 0; l < 8; l++ {
		for r := 0; r < 8; r++ {
			words := func(mask int) [][]string {
				out := [][]string{}
				for i, w := range universe {
					if mask&(1<<i) != 0 {
						out = append(out, w)
					}
				}
				return out
			}
			a, b := trie("left", words(l)), trie("right", words(r))
			g, err := product(a, b)
			if err != nil {
				t.Fatal(err)
			}
			pa := graphAutomaton(g, a)
			for i, w := range universe {
				if accepts(pa, w) != (l&r&(1<<i) != 0) {
					t.Fatal("product language differs from set intersection")
				}
			}
		}
	}
	t.Log("exhaustive: 256 families, 2048 hit tests, 65536 family pairs with 8 suffix tests each, 64 trie products")
}
