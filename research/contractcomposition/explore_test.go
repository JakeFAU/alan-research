package contractcomposition

import "testing"

func TestSerialCycleObstructionAndAcyclicSufficiency(t *testing.T) {
	for _, topology := range []string{"chain", "cycle"} {
		for mask := uint8(0); mask < 8; mask++ {
			c := fixture(topology, mask)
			g := Explore(c)
			if g.States == 0 || g.Arcs == 0 {
				t.Fatal("missing complete graph")
			}
			for _, n := range g.Nodes {
				if topology == "cycle" && mask == 0 && n.State.Stages == [3]Stage{W, W, W} && n.CanComplete {
					t.Fatal("all-serial circular wait incorrectly completes")
				}
				if g.Certificate.Acyclic && !n.CanComplete {
					t.Fatal("acyclic certificate failed")
				}
			}
			if (topology == "chain" || mask != 0) && g.CannotComplete != 0 {
				t.Fatal("service assumption should discharge progress")
			}
			if topology == "cycle" && mask == 0 && g.CannotComplete == 0 {
				t.Fatal("missing trapped state")
			}
		}
	}
}
func TestCompleteGraphRetainsAllIncomingArcs(t *testing.T) {
	g := Explore(fixture("chain", 7))
	if g.States != 16 || g.Arcs != 24 || g.CanComplete != 16 {
		t.Fatalf("two independent four-stage paths: got %d states %d arcs %d completing", g.States, g.Arcs, g.CanComplete)
	}
	incoming := map[int]int{}
	for _, n := range g.Nodes {
		for _, a := range n.Arcs {
			incoming[a.To]++
		}
	}
	found := false
	for _, n := range incoming {
		found = found || n > 1
	}
	if !found {
		t.Fatal("BFS discarded repeated-target arcs")
	}
}

func TestEvidenceIncludesAllConfigsOpenInterfacesAndRetry(t *testing.T) {
	e, err := BuildEvidence()
	if err != nil {
		t.Fatal(err)
	}
	if e.Schema != "alan.phase3c1.progress.v1" || len(e.Graphs) != 16 || len(e.Groupings) != 16 || len(e.Collisions) != 2 {
		t.Fatal("missing evidence domain or candidate failures")
	}
	for _, g := range e.Groupings {
		if !g.Equal || len(g.AB.Originals) != 2 || len(g.AB.Ports) != 1 {
			t.Fatal("group export erased an open intermediate or latent original")
		}
	}
	for _, w := range e.Witnesses {
		if err := CheckWitness(w); err != nil {
			t.Fatal(err)
		}
		if w.Name == e.Retry.Witness {
			states := []ConcreteState{w.Initial}
			for _, step := range w.Steps {
				states = append(states, step.State)
			}
			if e.Retry.FromState >= e.Retry.ToState || states[e.Retry.FromState] != states[e.Retry.ToState] {
				t.Fatal("retry certificate is not a real concrete loop")
			}
		}
	}
}
