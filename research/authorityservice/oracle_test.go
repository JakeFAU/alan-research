package authorityservice

import (
	base "github.com/JakeFAU/alan-research/research/mutualobligations"
	"strings"
	"testing"
)

func oracleTrace(t *testing.T, c Config, kinds []string, d base.Obligation) []Event {
	t.Helper()
	s := Initial(c)
	out := []Event{}
	for _, k := range kinds {
		var e Event
		s, e = step(t, c, s, k, d)
		out = append(out, e)
	}
	return out
}
func TestOracleCandidateSafety(t *testing.T) {
	for _, tc := range []struct {
		c    Config
		safe bool
	}{
		{Config{"S2", "Explicit"}, true}, {Config{"S2", "Missing"}, true},
		{Config{"S3", "Missing"}, false}, {Config{"Cartesian", "ScopeMix"}, false},
	} {
		trace := oracleTrace(t, tc.c, []string{"commit", "send-request", "admit", "authorize"}, base.BA)
		a, err := AuditTrace(tc.c, trace)
		if err != nil || !a.ProtocolValid || a.AuthoritySafe != tc.safe || a.Paid[base.BA] {
			t.Fatalf("%s: %+v %v", tc.c.Name(), a, err)
		}
	}
}
func TestOracleRejectsMutatedGrantAndEvents(t *testing.T) {
	c := Config{Policy: "S2", Scenario: "Explicit"}
	original := oracleTrace(t, c, []string{"commit", "send-request", "admit", "authorize", "dispatch-admit", "apply", "send-result", "receive-result", "record", "consume"}, base.BA)
	for _, tc := range []struct {
		name, reason string
		mutate       func([]Event)
	}{
		{"missing-parent", "missing-origin", func(es []Event) { es[0].Grant.Parent = "absent" }},
		{"self-cycle", "grant-cycle", func(es []Event) { es[0].Grant.Parent = es[0].Grant.ID }},
		{"mutated-root", "mutated-root", func(es []Event) { es[0].Grant.Root = "A/g1:a" }},
		{"generation", "identity-adjacency", func(es []Event) { es[0].Grant.GranteeGeneration = 2 }},
		{"enlarged-tuples", "tuple-attenuation", func(es []Event) { es[0].Grant.Permissions |= 1 << APermission }},
		{"rebound-original", "service-binding", func(es []Event) { es[0].Grant.Binding.Original = "AB" }},
		{"operation", "service-binding", func(es []Event) { es[0].Grant.Binding.Operation = 2 }},
		{"admission-epoch", "activation-binding", func(es []Event) { es[4].Activation = 1 }},
		{"result-as-authority", "result-is-not-authority", func(es []Event) { es[6].Authority = es[3].Authority }},
		{"foreign-equal-value", "original-binding", func(es []Event) { es[9].Obligation = base.AB }},
		{"unpaid-result", "consumption-custody", func(es []Event) { es[9].Result = base.R2 }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			es := append([]Event{}, original...)
			tc.mutate(es)
			a, err := AuditTrace(c, es)
			if err == nil || a.ProtocolValid || !strings.Contains(err.Error(), tc.reason) {
				t.Fatalf("audit=%+v err=%v, want %s", a, err, tc.reason)
			}
		})
	}
}
func TestOracleAllGraphEdges(t *testing.T) {
	if testing.Short() {
		t.Skip("full finite graph/prefix qualification")
	}
	for _, c := range Configurations() {
		t.Run(c.Name(), func(t *testing.T) {
			r, err := AuditGraph(c)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("prefixes=%d edges=%d unsafe-prefixes=%d unsafe-edges=%d paid=%v", r.Prefixes, r.Edges, r.UnsafePrefixes, r.UnsafeEdges, r.PaidOriginals)
		})
	}
}

func TestOracleRetainedCorpus(t *testing.T) {
	corpus, err := OracleCorpus()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("retained histories=%d retained chains=%d", len(corpus.Histories), len(corpus.Chains))
	for _, control := range corpus.Histories {
		t.Run(control.Name, func(t *testing.T) {
			if strings.HasPrefix(control.Name, "complete-values-") || control.Name == "released-receipt-both-crashes-then-paid" {
				if !control.Audit.Paid[0] || !control.Audit.Paid[1] {
					t.Fatal("completion without both consumed originals")
				}
			}
			if control.Name == "late-grant-sticky-refusal" && control.Audit.Paid[base.BA] {
				t.Fatal("refusal paid an original")
			}
			if control.Name == "pre-revocation-admission-applies" && control.Audit.Paid[base.BA] {
				t.Fatal("application paid an unconsumed original")
			}
		})
	}
}
func TestOracleBadInputsDoNotPanic(t *testing.T) {
	for _, event := range []Event{{Kind: "apply", Actor: 9}, {Kind: "authorize", Obligation: 9}, {Kind: "result-grants-authority", Obligation: base.BA}} {
		a, err := AuditTrace(Config{Policy: "S2", Scenario: "Explicit"}, []Event{event})
		if err == nil || a.ProtocolValid || a.AuthoritySafe {
			t.Fatalf("accepted malformed event: %+v", event)
		}
	}
	a, err := AuditTrace(Config{}, nil)
	if err == nil || a.ProtocolValid {
		t.Fatal("accepted unknown configuration")
	}
}

func TestOracleRevocationDoesNotRemoveSharedRoot(t *testing.T) {
	l := authorityLedger{grants: map[string]GrantRecord{}, revoked: true}
	c := Config{Policy: "S2", Scenario: "Shared"}
	for _, d := range []base.Obligation{base.AB, base.BA} {
		b := ledgerBinding(c, d)
		w := l.selectWitness(c, b)
		if w.Path != OwnRoot || !l.witnessLive(w, b) {
			t.Fatalf("shared root lost: %+v", w)
		}
	}
}
func TestOracleMutationsKeepApplicationOutcomes(t *testing.T) {
	corpus, err := OracleCorpus()
	if err != nil {
		t.Fatal(err)
	}
	compared := 0
	for _, control := range corpus.Histories {
		if control.WantProtocolValid || len(control.Events) != 20 {
			continue
		}
		outcomes := [2]base.Value{}
		for _, event := range control.Events {
			if event.Kind == "apply" {
				outcomes[event.Obligation] = event.Result
			}
		}
		if outcomes != [2]base.Value{base.R1, base.R1} {
			t.Fatalf("%s changed superficial application outcomes: %v", control.Name, outcomes)
		}
		compared++
	}
	if compared != 18 {
		t.Fatalf("checked %d preserved-outcome mutations, want 18", compared)
	}
}
