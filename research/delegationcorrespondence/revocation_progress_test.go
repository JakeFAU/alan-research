package delegationcorrespondence

import (
	"testing"

	m "github.com/JakeFAU/alan-research/research/authorityservice"
)

func TestGrantRevocationProgress(t *testing.T) {
	if testing.Short() {
		t.Skip("complete independent reachability")
	}
	r, err := QueryProgress(m.GrantRevocationConfig())
	if err != nil {
		t.Fatal(err)
	}
	if r.States != 229834 || r.Arcs != 1807246 || r.SafeComplete != 32768 || r.SafeCompletable != 203833 ||
		r.PredicateTrue != r.SafeCompletable || r.PredicateMismatches != 0 {
		t.Fatalf("unexpected finite predicate result: %+v", r)
	}
	if r.Revocation == nil || r.Revocation.ABRefusalStates != 0 || r.Revocation.StickyBARefusalEdges == 0 ||
		r.Revocation.CyclicStates != 188614 || r.Revocation.CyclicComponents != 75070 {
		t.Fatalf("unexpected revocation boundary: %+v", r.Revocation)
	}
	for _, cut := range []struct {
		name     string
		complete bool
	}{
		{"early-ba-refusal", false}, {"sticky-ba-refusal-with-grant", false}, {"revoked-after-allow-before-permit", false},
		{"revoked-permit-survives", true}, {"revoked-permit-crashed", false}, {"revoked-permit-applied", true},
	} {
		t.Run(cut.name, func(t *testing.T) {
			w, ok := r.Revocation.Cuts[cut.name]
			if !ok || w.CanComplete != cut.complete || !w.Witness.Audit.ProtocolValid || !w.Witness.Audit.AuthoritySafe {
				t.Fatalf("cut disagrees: %+v", w)
			}
			if cut.complete && (!w.Witness.Audit.Paid[0] || !w.Witness.Audit.Paid[1] || len(w.Witness.Suffix) == 0) {
				t.Fatal("successful suffix lacks both authentic consumptions")
			}
		})
	}
	if _, ok := r.Revocation.Cuts["predicate-mismatch"]; ok {
		t.Fatal("retained predicate counterexample")
	}
	if r.Counterexample == nil || len(r.Counterexample.Prefix) != 1 || r.Counterexample.Prefix[0].Kind != "revoke" {
		t.Fatal("lost shortest completion counterexample")
	}
}

func TestProgressRejectsUnfrozenRevocationPolicy(t *testing.T) {
	if _, err := QueryProgress(m.Config{Policy: "S3", Scenario: "SeparateGrantRevocation"}); err == nil {
		t.Fatal("unfrozen policy accepted")
	}
}

func TestProgressCycles(t *testing.T) {
	arcs := []progressArc{{From: 0, To: 1}, {From: 1, To: 0}, {From: 1, To: 2}, {From: 2, To: 2}, {From: 3, To: 2}}
	incoming := [][]int{{1}, {0}, {2, 3, 4}, {}}
	states, components := progressCycles(4, arcs, incoming)
	if states != 3 || components != 2 {
		t.Fatalf("SCC boundaries: %d states, %d components", states, components)
	}
}
