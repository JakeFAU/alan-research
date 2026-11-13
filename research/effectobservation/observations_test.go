package effectobservation

import (
	"reflect"
	"slices"
	"testing"
)

func TestIndependentObservationPartition(t *testing.T) {
	graph := Explore()
	report, err := AuditGraph(graph)
	if err != nil {
		t.Fatal(err)
	}
	if report.CheckedStates != len(graph.States) || report.CheckedPrefixes != len(graph.Prefixes) || report.CheckedArcs != 104 || report.MergeArcs != 42 {
		t.Fatalf("incomplete audit: %+v", report)
	}
	stores, dossiers := map[StoreView]bool{}, map[DossierView]bool{}
	for i, s := range graph.States {
		l, err := AuditTrace(graph.Prefixes[i])
		if err != nil {
			t.Fatal(err)
		}
		stores[l.StoreView()], dossiers[l.DossierView()] = true, true
		for j, other := range graph.States {
			r, err := AuditTrace(graph.Prefixes[j])
			if err != nil {
				t.Fatal(err)
			}
			if (l.StoreView() == r.StoreView()) != (ObserveStore(s) == ObserveStore(other)) {
				t.Fatalf("Store partition disagreement %d,%d", i, j)
			}
			if (l.DossierView() == r.DossierView()) != (ObserveDossier(s) == ObserveDossier(other)) {
				t.Fatalf("dossier partition disagreement %d,%d", i, j)
			}
		}
	}
	if report.StoreClasses != len(stores) || report.DossierClasses != len(dossiers) {
		t.Fatal("incorrect class count")
	}
	// Graph input order and choice among equal-length prefixes do not control ties.
	reordered := cloneAuditGraph(graph)
	for i := range reordered.Arcs {
		slices.Reverse(reordered.Arcs[i])
	}
	again, err := AuditGraph(reordered)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(report, again) {
		t.Fatal("audit depends on outgoing arc order")
	}
	t.Logf("states=%d arcs=%d merged targets=%d merge arcs=%d Store classes=%d dossier classes=%d", report.CheckedStates, report.CheckedArcs, report.MergedTargets, report.MergeArcs, report.StoreClasses, report.DossierClasses)
}

func TestObservationWitnesses(t *testing.T) {
	report, err := AuditGraph(Explore())
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range []WitnessPair{report.ZeroOne, report.DistinctValues} {
		if pair.Left.Store != pair.Right.Store || pair.Left.Dossier == pair.Right.Dossier {
			t.Fatal("invalid observation pair")
		}
		for _, w := range []Witness{pair.Left, pair.Right} {
			l, err := AuditTrace(w.Events)
			if err != nil {
				t.Fatal(err)
			}
			if !l.Restarted || !l.Revoked || !l.Refused || l.Stored != NoResult || l.Used != NoResult {
				t.Fatal("pair not qualified")
			}
		}
	}
	if len(report.ZeroOne.Left.Events) != 4 || len(report.ZeroOne.Right.Events) != 7 || len(report.DistinctValues.Left.Events) != 7 || len(report.DistinctValues.Right.Events) != 7 {
		t.Fatal("pair shortestness changed")
	}
	if report.ZeroOne.Left.Dossier.Applied != NoResult || report.ZeroOne.Right.Dossier.Applied != Value1 || report.DistinctValues.Right.Dossier.Applied != Value2 {
		t.Fatal("wrong application values")
	}
	cut, later := report.DelayedApplication.Cut, report.DelayedApplication.Later
	if len(cut.Events) != 3 || len(later.Events) != 5 || !slices.Equal(cut.Events, later.Events[:len(cut.Events)]) {
		t.Fatal("negative-cut witness not coherent/shortest")
	}
	if cut.Dossier.Ingress || cut.Dossier.Applied != NoResult || later.Dossier.Applied != Value1 || later.Store.Owner != OwnerStopped {
		t.Fatal("absence at cut conflated with future nonapplication")
	}
	recovery := report.RetainedRecovery
	if len(recovery.Events) != 9 || recovery.Events[8].Kind != Consume || recovery.Events[8].Incarnation != 2 {
		t.Fatal("recovery needs consumption after reopen")
	}
	prior, err := AuditTrace(recovery.Events[:8])
	if err != nil {
		t.Fatal(err)
	}
	if !prior.Restarted || !prior.Revoked || prior.Stored != Value1 || prior.Used != NoResult {
		t.Fatal("invalid recovery source")
	}
}

func cloneAuditGraph(g Graph) Graph {
	g.States = slices.Clone(g.States)
	g.Arcs = slices.Clone(g.Arcs)
	g.Prefixes = slices.Clone(g.Prefixes)
	for i := range g.States {
		g.Arcs[i] = slices.Clone(g.Arcs[i])
		g.Prefixes[i] = slices.Clone(g.Prefixes[i])
	}
	return g
}

func TestAuditRejectsMalformedGraph(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Graph)
	}{
		{"empty", func(g *Graph) { *g = Graph{} }},
		{"array-length", func(g *Graph) { g.Arcs = g.Arcs[:1] }},
		{"missing-edge", func(g *Graph) { g.Arcs[0] = g.Arcs[0][1:] }},
		{"duplicate-edge", func(g *Graph) { g.Arcs[0] = append(g.Arcs[0], g.Arcs[0][0]) }},
		{"invalid-target", func(g *Graph) { g.Arcs[0][0].To = -1 }},
		{"merged-target", func(g *Graph) { g.Arcs[0][0].To = 0 }},
		{"wrong-event", func(g *Graph) { g.Arcs[0][0].Event.Value = Value1 }},
		{"rewritten-first-allow", func(g *Graph) { g.States[1].FirstAllow = false }},
		{"rewritten-gate", func(g *Graph) { g.States[1].FirstGateAllowed = false }},
		{"debit-amplification", func(g *Graph) { g.States[1].Debit = 2 }},
		{"c-revocation", func(g *Graph) { g.States[1].ControlRoot = RootB }},
		{"duplicate-state", func(g *Graph) { g.States[1] = g.States[0] }},
		{"wrong-prefix", func(g *Graph) { g.Prefixes[1] = nil }},
		{"extra-unreachable", func(g *Graph) {
			g.States = append(g.States, g.States[0])
			g.Arcs = append(g.Arcs, nil)
			g.Prefixes = append(g.Prefixes, nil)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := cloneAuditGraph(Explore())
			tt.change(&g)
			if _, err := AuditGraph(g); err == nil {
				t.Fatal("accepted malformed graph")
			}
		})
	}
}
