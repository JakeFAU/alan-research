package residualantichains

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestCompleteEvidence(t *testing.T) {
	e, err := build("../..")
	if err != nil {
		t.Fatal(err)
	}
	if len(e.Products) != 3564 || len(e.Comparisons) != 36 || len(e.Input.Sources) != 40 {
		t.Fatalf("incomplete roster: %+v", e.Counts)
	}
	if len(e.Substitutions) != e.Counts.EqualJ*396 {
		t.Fatal("missing contexts")
	}
	for _, s := range e.Substitutions {
		if s.InitialLeft != s.InitialRight || s.NonblockingLeft != s.NonblockingRight {
			t.Fatal("J substitution failed", s)
		}
	}
	for _, c := range e.Comparisons {
		for _, n := range c.Pairs {
			for _, cs := range [][]coverage{n.LeftToRight, n.RightToLeft} {
				for _, cov := range cs {
					if cov.Failure != nil {
						if !cov.Failure.TargetMiss {
							t.Fatal("target hit by separator")
						}
						for _, hit := range cov.Failure.Hits {
							if !hit.Hit {
								t.Fatal("source missed by separator")
							}
						}
					}
				}
			}
		}
	}
	byID := map[string]graph{}
	for _, g := range e.Products {
		byID[g.ID] = g
		ns := nodesByID(g)
		arcs, can := 0, 0
		for _, n := range g.Nodes {
			if n.CanComplete {
				can++
			}
			arcs += len(n.Arcs)
			for _, a := range n.Arcs {
				if _, ok := ns[a.To]; !ok {
					t.Fatal("dangling graph arc")
				}
			}
		}
		if len(g.Nodes) != g.States || arcs != g.Arcs || can != g.CanComplete || g.Nonblocking != (g.States == can) {
			t.Fatal("invalid graph counts")
		}
	}
	neg := e.Negative
	c := e.Comparisons[neg.Comparison]
	l, r := byID[neg.LeftProduct], byID[neg.RightProduct]
	if !c.MAX || !c.InitialLanguage.Equal || c.MIN || !l.Nonblocking || r.Nonblocking || !l.Nodes[0].CanComplete || !r.Nodes[0].CanComplete {
		t.Fatal("MAX counterexample lost")
	}
	checkPath := func(g graph, p path) {
		t.Helper()
		ns := nodesByID(g)
		at := p.Start
		for _, s := range p.Steps {
			if s.From != at || !slices.Contains(ns[at].Arcs, arc{Label: s.Label, To: s.To}) {
				t.Fatal("invented path")
			}
			at = s.To
		}
		if at != p.End {
			t.Fatal("path endpoint")
		}
	}
	checkPath(l, neg.LeftPath)
	checkPath(r, neg.RightPath)
	if len(neg.LeftPath.Steps) != 2 || len(neg.RightPath.Steps) != 3 || !nodesByID(l)[neg.LeftPath.End].CanComplete || nodesByID(r)[neg.RightPath.End].CanComplete {
		t.Fatal("MAX legal core paths")
	}
	for _, c := range e.LegalControl.Inclusions {
		if !c.LeftSubset || c.RightSubset || c.RightOutside == nil {
			t.Fatal("proper residual difference disabled by legal plant")
		}
	}
	for _, name := range []string{"minimal-base", "minimal-augmented"} {
		good, bad, kd := 0, 0, 0
		for _, g := range e.Products {
			if g.Factors[1] != name {
				continue
			}
			if len(g.Factors) == 3 {
				kd++
				if g.Nonblocking {
					t.Fatal("new fixture K/D unexpectedly nonblocking")
				}
				continue
			}
			if g.Nonblocking {
				good++
			} else {
				bad++
			}
		}
		if good != 52 || bad != 20 || kd != 324 {
			t.Fatalf("fixture outcomes %s: %d %d %d", name, good, bad, kd)
		}
	}
	if !e.Generics.MINLoss.Comparison.MIN || e.Generics.MINLoss.Comparison.InitialLanguage.Equal || !e.Generics.MINLoss.LeftProduct.Nodes[0].CanComplete || e.Generics.MINLoss.RightProduct.Nodes[0].CanComplete {
		t.Fatal("MIN loss")
	}
	if e.Generics.JNonnecessity.Comparison.J || !e.Generics.JNonnecessity.Comparison.InitialLanguage.Equal {
		t.Fatal("generic nonnecessity")
	}
	for i, g := range e.Generics.Multiword.Products {
		if g.Nodes[0].CanComplete != (i < 2) {
			t.Fatal("multiword trie initial hit")
		}
	}
	if !e.Generics.LocalHitNotNB.Product.Nodes[0].CanComplete || e.Generics.LocalHitNotNB.Product.Nonblocking {
		t.Fatal("local versus global")
	}
	t.Logf("measured counts: %+v", e.Counts)
}
func TestInputAuthenticationAndCompatibility(t *testing.T) {
	p, _, err := readInput("../..")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	path := filepath.Join(root, predecessorPath)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("wrong fixture"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readInput(root); err == nil {
		t.Fatal("unauthenticated descriptor accepted")
	}
	k, m := p.Plants[0], p.Monitors[2]
	m.Signature = k.Signature
	for _, mutation := range []string{"signature", "root", "role", "veto", "external-target"} {
		t.Run(mutation, func(t *testing.T) {
			bad := m
			bad.Edges = slices.Clone(m.Edges)
			switch mutation {
			case "signature":
				bad.Signature = append([]byte(" "), m.Signature...)
			case "root":
				var sig map[string]any
				if err := json.Unmarshal(m.Signature, &sig); err != nil {
					t.Fatal(err)
				}
				sig["Roots"].([]any)[0].(map[string]any)["Limit"] = 99
				var err error
				bad.Signature, err = json.Marshal(sig)
				if err != nil {
					t.Fatal(err)
				}
			case "role":
				bad.Role = "unknown"
			case "veto":
				for i, e := range bad.Edges {
					if e.Label == "commit:AB" {
						bad.Edges = append(bad.Edges[:i], bad.Edges[i+1:]...)
						break
					}
				}
			case "external-target":
				bad.Edges = append(bad.Edges, edge{From: "0", Label: "commit:AB", To: "1"})
			}
			if _, err := product(k, bad); err == nil {
				t.Fatal("invalid compatibility accepted")
			}
		})
	}
}

// These are deliberately false research hypotheses, enabled only for separately
// captured falsification runs. Passing counterexample assertions above are normal tests.
func TestFalseHypotheses(t *testing.T) {
	hypothesis := os.Getenv("ALAN_FALSE_HYPOTHESIS")
	if hypothesis == "" {
		t.Skip("false hypotheses require explicit diagnostic selection")
	}
	switch hypothesis {
	case "MIN-preserves-initial":
		x, y := genericMINLoss()
		c := compareProfiles(x, y)
		if c.MIN && !c.InitialLanguage.Equal {
			t.Fatal("falsified: MIN alone preserves initial language")
		}
	case "MAX-preserves-NB":
		p, _, err := readInput("../..")
		if err != nil {
			t.Fatal(err)
		}
		k := p.Plants[0]
		for _, a := range p.Plants {
			if a.ID == "early/shared-1/mask0" {
				k = a
			}
		}
		l, r := p.Monitors[2], p.Monitors[4]
		l.Signature = k.Signature
		r.Signature = k.Signature
		gl, err := product(k, l)
		if err != nil {
			t.Fatal(err)
		}
		gr, err := product(k, r)
		if err != nil {
			t.Fatal(err)
		}
		c := compareProfiles(l, r)
		if c.MAX && c.InitialLanguage.Equal && gl.Nonblocking != gr.Nonblocking {
			t.Fatal("falsified: MAX plus initial language preserves NB")
		}
	case "local-hit-implies-NB":
		g, err := makeGenerics()
		if err != nil {
			t.Fatal(err)
		}
		p := g.LocalHitNotNB.Product
		if p.Nodes[0].CanComplete && !p.Nonblocking {
			t.Fatal("falsified: a local suffix hit implies whole-product NB")
		}
	default:
		t.Fatalf("unknown false hypothesis %q", hypothesis)
	}
}
