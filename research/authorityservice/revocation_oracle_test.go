package authorityservice

import (
	"testing"

	base "github.com/JakeFAU/alan-research/research/mutualobligations"
)

func TestRevocationOracleDistinctRoot(t *testing.T) {
	actions := []corpusAction{{kind: "revoke", d: base.NoObligation, actor: base.B}}
	for _, kind := range []string{"commit", "send-request", "admit", "authorize", "dispatch-admit", "apply"} {
		actions = append(actions, corpusAction{kind: kind, d: base.AB})
	}
	events, err := corpusSequence(GrantRevocationConfig(), actions)
	if err != nil {
		t.Fatal(err)
	}
	a, err := AuditTrace(GrantRevocationConfig(), events)
	if err != nil || !a.ProtocolValid || !a.AuthoritySafe {
		t.Fatalf("distinct c root must survive b revocation: %+v, %v", a, err)
	}
}

func TestRevocationOracleRejectedIssuance(t *testing.T) {
	events, err := corpusSequence(GrantRevocationConfig(), []corpusAction{
		{kind: "revoke", d: base.NoObligation, actor: base.B}, {kind: "commit", d: base.BA}, {kind: "retain-grant", d: base.BA},
	})
	if err != nil {
		t.Fatal(err)
	}

	events[len(events)-1].Grant.Accepted = true
	a, err := AuditTrace(GrantRevocationConfig(), events)
	if err == nil || a.ProtocolValid {
		t.Fatal("accepted post-revocation issuance was not rejected")
	}
	events[len(events)-1].Grant.Accepted = false
	a, err = AuditTrace(GrantRevocationConfig(), events)
	if err != nil || !a.ProtocolValid || !a.AuthoritySafe {
		t.Fatalf("historical rejected issuance must be retained: %+v, %v", a, err)
	}

}

func TestRevocationOracleCorpus(t *testing.T) {
	report, err := GrantRevocationOracleCorpus()
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Histories) < 30 {
		t.Fatal("hostile controls missing")
	}
	for _, c := range report.Histories {
		t.Run(c.Name, func(t *testing.T) {
			if c.Audit.ProtocolValid != c.WantProtocolValid {
				t.Fatalf("unexpected verdict: %+v", c)
			}
			if c.Name == "post-revocation-both-authentic-consumptions" && (!c.Audit.Paid[0] || !c.Audit.Paid[1]) {
				t.Fatal("completion invented or lost consumption")
			}
		})
	}
}

func TestRevocationOracleFreshAdmission(t *testing.T) {
	events, err := corpusSequence(GrantRevocationConfig(), []corpusAction{
		{kind: "commit", d: base.BA}, {kind: "retain-grant", d: base.BA}, {kind: "send-request", d: base.BA},
		{kind: "admit", d: base.BA}, {kind: "authorize", d: base.BA}, {kind: "revoke", d: base.NoObligation, actor: base.B},
		{kind: "dispatch-refuse", d: base.BA},
	})
	if err != nil {
		t.Fatal(err)
	}
	if a, err := AuditTrace(GrantRevocationConfig(), events); err != nil || !a.ProtocolValid {
		t.Fatalf("fresh refusal lost historical allow: %+v %v", a, err)
	}
	events[len(events)-1].Kind = "dispatch-admit"
	events[len(events)-1].Refusal = ""
	if a, err := AuditTrace(GrantRevocationConfig(), events); err == nil || a.ProtocolValid {
		t.Fatal("fresh revoked admission accepted")
	}
}

func TestRevocationOracleCrashPermit(t *testing.T) {
	events, err := corpusSequence(GrantRevocationConfig(), []corpusAction{
		{kind: "commit", d: base.BA}, {kind: "retain-grant", d: base.BA}, {kind: "send-request", d: base.BA},
		{kind: "admit", d: base.BA}, {kind: "authorize", d: base.BA}, {kind: "dispatch-admit", d: base.BA},
		{kind: "revoke", d: base.NoObligation, actor: base.B}, {kind: "crash", d: base.NoObligation, actor: base.A},
		{kind: "recover", d: base.NoObligation, actor: base.A}, {kind: "dispatch-refuse", d: base.BA},
	})
	if err != nil {
		t.Fatal(err)
	}
	if a, err := AuditTrace(GrantRevocationConfig(), events); err != nil || !a.ProtocolValid {
		t.Fatalf("crashed permit must clear before fresh refusal: %+v %v", a, err)
	}
}

func TestRevocationOracleGraph(t *testing.T) {
	if testing.Short() {
		t.Skip("complete graph audit")
	}
	r, err := AuditGraph(GrantRevocationConfig())
	if err != nil {
		t.Fatal(err)
	}
	if r.Prefixes != 229834 || r.Edges != 1807246 || r.UnsafePrefixes != 0 || r.UnsafeEdges != 0 {
		t.Fatalf("unexpected complete audit: %+v", r)
	}
	t.Logf("independent audit: %+v", r)
}
