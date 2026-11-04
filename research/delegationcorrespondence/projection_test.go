package delegationcorrespondence

import (
	m "github.com/JakeFAU/alan-research/research/authorityservice"
	"os"
	"testing"
)

func TestRawProjection(t *testing.T) {
	p := ProjectRaw(positive(t))
	if !p.Supported || p.Audit == nil || !p.Audit.ProtocolValid || !p.Audit.AuthoritySafe || len(p.Events) == 0 {
		t.Fatalf("%+v", p)
	}
	if !p.Audit.Paid[0] || !p.Audit.Paid[1] {
		t.Fatal("not complete")
	}
}
func TestProgressQueries(t *testing.T) {
	if os.Getenv("ALAN_TASK3_PROGRESS") == "" {
		t.Skip("set ALAN_TASK3_PROGRESS=1 for full independent graph queries")
	}
	for _, v := range []struct {
		s       string
		n, a, r int
	}{{"LateGrant", 111065, 840616, 107213}, {"Revocable", 220564, 1740832, 193681}} {
		p, e := QueryProgress(m.Config{Policy: "S2", Scenario: v.s})
		if e != nil {
			t.Fatal(e)
		}
		if p.Counterexample == nil || p.SafeCompletable == p.States {
			t.Fatal("unconditional completion adverse control lost")
		}
		if v.s == "LateGrant" && (p.NonCompletableNoBARefusal != 0 || p.RefusedNoGrant != 1926 || p.RefusedRetainedGrant != 1926) {
			t.Fatal("late-grant exclusions changed")
		}
		if v.s == "Revocable" && p.NonCompletableNoBARefusal != 22887 {
			t.Fatal("weak revocable predicate no longer records adverse states")
		}
		if p.States != v.n || p.Arcs != v.a || p.SafeCompletable != v.r || p.PredicateMismatches != 0 {
			t.Fatalf("%+v", p)
		}
	}
}

func TestFrozenProjectionDomains(t *testing.T) {
	root := os.Getenv("ALAN_TASK3_RAW")
	if root == "" {
		t.Skip("frozen corpus")
	}
	r, e := CheckCorpus(root, root)
	if e != nil {
		t.Fatal(e)
	}
	for _, c := range r.Cases {
		if c.Name != "control-wrong-prior-selection" && c.Projection.Mismatch != "" {
			t.Errorf("misnamed mismatch %s: %s", c.Name, c.Projection.Mismatch)
		}
		switch c.Name {
		case "control-removed-original-tuple", "actual-prior-service-result-as-grant", "cut-07-retained-result":
			if c.Projection.Supported || c.Projection.Error != "" {
				t.Errorf("%s must preserve unsupported-domain reason without pretending agreement: %+v", c.Name, c.Projection.ProjectionReport)
			}
		}
	}
}
