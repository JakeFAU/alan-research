package residualinterfaces

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

func fixture(t *testing.T) predecessor {
	t.Helper()
	p, _, err := readInput("../..")
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func TestPinnedInputAndMonitorProfiles(t *testing.T) {
	p := fixture(t)
	ms := append(p.Monitors, extraMonitors(p.Monitors[0].Signature)...)
	if len(ms) != 7 {
		t.Fatal("roster")
	}
	for _, a := range ms {
		if err := validate(a); err != nil {
			t.Fatal(err)
		}
	}
	for _, i := range []int{3, 5, 6} {
		c := compareProfiles(ms[2], ms[i])
		if !c.Equal || !c.InitialLanguage.Equal {
			t.Fatalf("RF versus %s", ms[i].ID)
		}
	}
	c := compareProfiles(ms[2], ms[4])
	if c.Equal || !c.InitialLanguage.Equal || !slices.Equal(c.Word, []string{"reserve-capacity:AB"}) || c.Mismatch.Side != "right" || c.Mismatch.State != "3" || len(c.Mismatch.Separators) != 1 {
		t.Fatalf("doomed mismatch: %+v", c)
	}
	if !reflect.DeepEqual(p.Monitors, ms[:5]) {
		t.Fatal("changed predecessor monitors")
	}
	// Exact fixture authentication is required before any descriptor is accepted.
	root := t.TempDir()
	path := filepath.Join(root, predecessorPath)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("wrong fixture"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readInput(root); err == nil {
		t.Fatal("accepted changed fixture")
	}
}
func TestCompleteEvidenceAndWitnesses(t *testing.T) {
	e, err := build("../..")
	if err != nil {
		t.Fatal(err)
	}
	if len(e.Products) != 2772 || len(e.Comparisons) != 21 {
		t.Fatalf("roster: %+v", e.Counts)
	}
	if e.Counts.EqualProfiles != 6 || len(e.Substitutions) != 6*396 {
		t.Fatalf("profile equality counts: %+v", e.Counts)
	}
	if e.Counts.BlockingSubstitutions == 0 || e.Counts.NonblockingSubstitutions == 0 {
		t.Fatal("discarded blocked or successful contexts")
	}
	for _, s := range e.Substitutions {
		if s.NonblockingLeft != s.NonblockingRight {
			t.Fatal("substitution failed", s)
		}
	}
	graphs := map[string]graph{}
	for _, g := range e.Products {
		graphs[g.ID] = g
		if len(g.Nodes) != g.States {
			t.Fatal("missing nodes")
		}
		nodes := nodesByID(g)
		arcs, can := 0, 0
		for _, n := range g.Nodes {
			arcs += len(n.Arcs)
			if n.CanComplete {
				can++
			}
			for _, a := range n.Arcs {
				if _, ok := nodes[a.To]; !ok {
					t.Fatal("dangling arc")
				}
			}
		}
		if arcs != g.Arcs || can != g.CanComplete || g.Nonblocking != (can == g.States) {
			t.Fatal("graph labels/counts")
		}
	}
	checkPath := func(g graph, p path) {
		t.Helper()
		at := p.Start
		nodes := nodesByID(g)
		if _, ok := nodes[at]; !ok {
			t.Fatal("missing start")
		}
		for _, s := range p.Steps {
			if s.From != at || !slices.Contains(nodes[at].Arcs, arc{Label: s.Label, To: s.To}) {
				t.Fatal("invented path step", s)
			}
			at = s.To
		}
		if at != p.End {
			t.Fatal("wrong endpoint")
		}
	}
	neg := e.Negative
	l, r := graphs[neg.LeftProduct], graphs[neg.RightProduct]
	checkPath(l, neg.LeftPath)
	checkPath(r, neg.RightPath)
	if !l.Nonblocking || r.Nonblocking || !l.Nodes[0].CanComplete || !r.Nodes[0].CanComplete {
		t.Fatal("negative core classification")
	}
	if !nodesByID(l)[neg.LeftPath.End].CanComplete || nodesByID(r)[neg.RightPath.End].CanComplete {
		t.Fatal("negative residual classification")
	}
	if len(neg.LeftPath.Steps) != 2 || len(neg.RightPath.Steps) != 3 || neg.LeftPath.Steps[0].Label != "commit:AB" {
		t.Fatal("monitor prefix confused with legal K path")
	}
	loop := e.Loop
	g := graphs[loop.Product]
	nodes := nodesByID(g)
	checkPath(g, loop.Prefix)
	checkPath(g, loop.Cycle)
	if !g.Nonblocking || loop.Prefix.End != loop.Cycle.Start || loop.Cycle.Start != loop.Cycle.End || len(loop.Cycle.Steps) != 2 {
		t.Fatal("invalid lasso")
	}
	for i, s := range loop.Cycle.Steps {
		n := nodes[s.From]
		if s.Label != "tau" || !unpaid(n) || n.Marked || !n.CanComplete {
			t.Fatal("cycle lost unpaid completable debt")
		}
		escape := loop.Escapes[i]
		checkPath(g, escape)
		if escape.Start != s.From || !nodes[escape.End].Marked {
			t.Fatal("invalid escape")
		}
	}
	if e.Strictness.Comparison.Equal || !e.Strictness.Comparison.InitialLanguage.Equal ||
		!slices.Equal(e.Strictness.Comparison.Word, []string{"apply:AB"}) {
		t.Fatal("strictness")
	}
	t.Logf("measured counts: %+v", e.Counts)
}
func TestExactSignatureAndExternalActions(t *testing.T) {
	p := fixture(t)
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
				var signature map[string]any
				if err := json.Unmarshal(m.Signature, &signature); err != nil {
					t.Fatal(err)
				}
				signature["Roots"].([]any)[0].(map[string]any)["Limit"] = 99
				var err error
				bad.Signature, err = json.Marshal(signature)
				if err != nil {
					t.Fatal(err)
				}
			case "role":
				bad.Role = "unrecognized"
			case "veto":
				for i, e := range bad.Edges {
					if e.Label == "commit:AB" {
						bad.Edges = append(bad.Edges[:i], bad.Edges[i+1:]...)
						break
					}
				}
			case "external-target":
				bad.Edges = append(bad.Edges, edge{"0", "commit:AB", "1"})
			}
			if _, err := product(k, bad); err == nil {
				t.Fatal("accepted invalid composition", mutation)
			}
		})
	}
}
