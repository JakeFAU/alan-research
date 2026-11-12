package authorityservice

import (
	"testing"

	base "github.com/JakeFAU/alan-research/research/mutualobligations"
)

func TestSeparateGrantRevocationCommitWithoutGrant(t *testing.T) {
	c := GrantRevocationConfig()
	if c != (Config{Policy: "S2", Scenario: "SeparateGrantRevocation"}) {
		t.Fatalf("unexpected opt-in configuration: %s", c.Name())
	}
	for _, old := range Configurations() {
		if old == c {
			t.Fatal("opt-in configuration leaked into defaults")
		}
	}
	s := Initial(c)
	noStep(t, c, s, "retain-grant", base.BA)
	s, e := step(t, c, s, "commit", base.BA)
	if !s.Base.Obligations[base.BA].Committed || e.HasGrant || s.Grant != NoGrant {
		t.Fatal("BA commitment failed to remain separate from grant issuance")
	}
}

func TestSeparateGrantRevocationIssuance(t *testing.T) {
	c := GrantRevocationConfig()
	t.Run("accepted-while-root-current", func(t *testing.T) {
		s, _ := step(t, c, Initial(c), "commit", base.BA)
		s = actorStep(t, c, s, "crash", base.B)
		noStep(t, c, s, "retain-grant", base.BA)
		noStep(t, c, s, "revoke", base.NoObligation)
		s = actorStep(t, c, s, "recover", base.B)
		s, e := step(t, c, s, "retain-grant", base.BA)
		want := GrantRecord{
			ID: "grant:BA:b", Parent: "B/g1:b", Root: "B/g1:b",
			Issuer: base.B, IssuerGeneration: 1, Grantee: base.A, GranteeGeneration: 1,
			Permissions: Set(BPermission), NoRedelegation: true, Accepted: true,
			Binding: ServiceBinding{
				Original: "BA", Caller: base.B, CallerGeneration: 1, Callee: base.A,
				CalleeGeneration: 1, Operation: 1,
				Permission: Permission{Operation: "increment", Resource: "counter-B", Scope: "blue"},
			},
		}
		if !e.HasGrant || e.Grant != want || s.Grant != RetainedGrant ||
			e.Actor != base.B || e.Activation != 1 || e.Revision != 0 {
			t.Fatalf("unexpected accepted separate issuance: %+v (status %d)", e, s.Grant)
		}
		noStep(t, c, s, "retain-grant", base.BA)
		revoked := actorStep(t, c, s, "revoke", base.B)
		if revoked.Grant != RetainedGrant {
			t.Fatal("revocation rewrote accepted issuance")
		}
	})
	t.Run("rejected-after-revocation", func(t *testing.T) {
		s := actorStep(t, c, Initial(c), "revoke", base.B)
		noStep(t, c, s, "revoke", base.NoObligation)
		s, e := step(t, c, s, "commit", base.BA)
		if !s.Base.Obligations[base.BA].Committed || e.HasGrant || s.Grant != NoGrant {
			t.Fatal("revocation altered ordinary BA commitment")
		}
		s, e = step(t, c, s, "retain-grant", base.BA)
		if !e.HasGrant || e.Grant.Accepted || e.Grant.Root != "B/g1:b" ||
			e.Revision != 1 || e.Grant.Revision != 0 || s.Grant != RejectedGrant {
			t.Fatalf("missing rejected separate issuance: %+v (status %d)", e, s.Grant)
		}
		noStep(t, c, s, "retain-grant", base.BA)
		s = revocationSteps(t, s, base.BA, "send-request", "admit", "authorize")
		if s.Authority[base.BA].Decision != Refused || s.Grant != RejectedGrant {
			t.Fatal("rejected issuance supplied authority or was rewritten")
		}
	})
}

func TestSeparateGrantRevocationPreservesFirstAllow(t *testing.T) {
	c := GrantRevocationConfig()
	s, decision := separateFirstAllow(t)
	s = actorStep(t, c, s, "revoke", base.B)
	s, refused := step(t, c, s, "dispatch-refuse", base.BA)
	if decision.Decision != Allowed || decision.Authority.Path != Delegated ||
		refused.Authority != decision.Authority || refused.Decision != Allowed ||
		refused.Refusal != "selected-path-revoked" || refused.Revision != 1 ||
		s.Authority[base.BA].Decision != Allowed || !s.Authority[base.BA].AdmissionRefused ||
		s.Grant != RetainedGrant {
		t.Fatalf("fresh refusal did not preserve historical allow: %+v, %+v", decision, refused)
	}
	for _, kind := range []string{"authorize", "dispatch-admit", "dispatch-refuse", "apply"} {
		noStep(t, c, s, kind, base.BA)
	}
}

func TestSeparateGrantRevocationAdmittedApplication(t *testing.T) {
	c := GrantRevocationConfig()
	s, decision := separateFirstAllow(t)
	s = revocationSteps(t, s, base.BA, "dispatch-admit")
	s = actorStep(t, c, s, "revoke", base.B)
	s, e := step(t, c, s, "apply", base.BA)
	if e.Revision != 1 || e.Authority != decision.Authority || !s.Base.Obligations[base.BA].Applied ||
		e.Result == base.NoValue || s.Base.Obligations[base.BA].Retained != e.Result ||
		s.Authority[base.BA].Permit {
		t.Fatal("pre-revocation BA opportunity failed to apply once with its authentic result")
	}
	noStep(t, c, s, "apply", base.BA)
}

func TestSeparateGrantRevocationCrashClearsPermit(t *testing.T) {
	c := GrantRevocationConfig()
	s, decision := separateFirstAllow(t)
	s = revocationSteps(t, s, base.BA, "dispatch-admit")
	s = actorStep(t, c, s, "revoke", base.B)
	s = actorStep(t, c, s, "crash", base.A)
	if s.Authority[base.BA].Permit {
		t.Fatal("crash retained BA permit")
	}
	s = actorStep(t, c, s, "recover", base.A)
	noStep(t, c, s, "apply", base.BA)
	noStep(t, c, s, "dispatch-admit", base.BA)
	s, e := step(t, c, s, "dispatch-refuse", base.BA)
	if e.Activation != 1 || e.Authority != decision.Authority ||
		s.Authority[base.BA].Permit || s.Authority[base.BA].PermitActivation != 0 ||
		!s.Authority[base.BA].AdmissionRefused {
		t.Fatal("recovery failed to refuse a fresh BA opportunity in the new activation")
	}
}

func TestSeparateGrantRevocationLateGrantRefusalIsSticky(t *testing.T) {
	c := GrantRevocationConfig()
	s := revocationSteps(t, Initial(c), base.BA, "commit", "send-request", "admit", "authorize")
	a := s.Authority[base.BA]
	s, e := step(t, c, s, "retain-grant", base.BA)
	s = actorStep(t, c, s, "revoke", base.B)
	s = actorStep(t, c, s, "crash", base.A)
	s = actorStep(t, c, s, "recover", base.A)
	if a.Decision != Refused || a.Path != NoPath || !e.Grant.Accepted ||
		s.Grant != RetainedGrant || s.Authority[base.BA] != a {
		t.Fatal("later accepted grant or recovery changed the early refusal")
	}
	for _, kind := range []string{"authorize", "dispatch-admit", "apply"} {
		noStep(t, c, s, kind, base.BA)
	}
}

func TestSeparateGrantRevocationDistinctABRootCompletes(t *testing.T) {
	c := GrantRevocationConfig()
	want := Permission{Operation: "increment", Resource: "counter-AB", Scope: "blue"}
	if Binding(c, base.AB).Permission != want || Binding(c, base.BA).Permission == want {
		t.Fatal("AB and BA must have distinct complete service permissions")
	}
	for _, old := range Configurations() {
		if Binding(old, base.AB).Permission == want || Binding(old, base.BA).Permission == want {
			t.Fatalf("new AB permission leaked into %s", old.Name())
		}
	}

	s, _ := separateFirstAllow(t)
	s = revocationSteps(t, s, base.BA, "dispatch-admit")
	s = actorStep(t, c, s, "revoke", base.B)
	s = revocationSteps(t, s, base.AB, "commit", "send-request", "admit")
	s, e := step(t, c, s, "authorize", base.AB)
	if e.Decision != Allowed || e.Authority.Path != OwnRoot || e.Authority.Root != "B/g1:c" ||
		e.Authority.Grant != "" || e.Authority.Permission != want {
		t.Fatalf("AB lacks its distinct own-root path after b revocation: %+v", e)
	}
	s, e = step(t, c, s, "dispatch-admit", base.AB)
	if e.Authority.Root != "B/g1:c" {
		t.Fatal("AB dispatch changed the historical root")
	}
	noStep(t, c, s, "apply", base.AB) // B still waits for its outgoing BA result.
	s, ba := step(t, c, s, "apply", base.BA)
	s = revocationSteps(t, s, base.BA, "send-result", "receive-result", "record")
	noStep(t, c, s, "apply", base.AB) // Recording alone does not finish BA.
	s, consumedBA := step(t, c, s, "consume", base.BA)
	var ab Event
	for _, tr := range Successors(c, s) {
		if tr.Event.Kind == "apply" && tr.Event.Obligation == base.AB && tr.Event.Result == base.R2 {
			s, ab = tr.State, tr.Event
			break
		}
	}
	s = revocationSteps(t, s, base.AB, "send-result", "receive-result", "record")
	s, consumedAB := step(t, c, s, "consume", base.AB)
	if !Complete(s) || s.Unsafe != 0 || ba.Result != base.R1 || consumedBA.Result != ba.Result ||
		ab.Result != base.R2 || consumedAB.Result != ab.Result || ab.Authority.Root != "B/g1:c" {
		t.Fatal("reciprocal completion failed to consume both distinct authentic results")
	}
}

func TestSeparateGrantRevocationConsumedResultIsImmutable(t *testing.T) {
	c := GrantRevocationConfig()
	s, _ := separateFirstAllow(t)
	s = revocationSteps(t, s, base.BA, "dispatch-admit", "apply", "send-result", "receive-result", "record")
	before := s.Base.Obligations[base.BA]
	s = actorStep(t, c, s, "revoke", base.B)
	if s.Base.Obligations[base.BA] != before {
		t.Fatal("revocation altered the recorded original result")
	}
	s, e := step(t, c, s, "consume", base.BA)
	if e.Result != before.WorldResult || e.Result == base.NoValue {
		t.Fatal("consumption after revocation lost the authentic result")
	}
	s = revocationSteps(t, s, base.BA, "send-ack", "release")
	for _, actor := range []base.Agent{base.A, base.B} {
		s = actorStep(t, c, s, "crash", actor)
		s = actorStep(t, c, s, "recover", actor)
	}
	b := s.Base.Obligations[base.BA]
	if !b.Consumed || b.Recorded != e.Result || b.WorldResult != e.Result || b.Retained != base.NoValue {
		t.Fatal("release or recovery changed durable consumed result")
	}
	noStep(t, c, s, "consume", base.BA)
	noStep(t, c, s, "apply", base.BA)
}

func TestSeparateGrantRevocationGraph(t *testing.T) {
	if testing.Short() {
		t.Skip("complete graph qualification")
	}
	c := GrantRevocationConfig()
	r, g := explore(c)
	t.Logf("states=%d arcs=%d safe-complete=%d safe-suffix=%d no-safe-suffix=%d cyclic-states=%d cyclic-components=%d",
		r.States, r.Arcs, r.SafeCompleteStates, r.CanCompleteSafely, r.CannotCompleteSafely, r.CyclicStates, r.CyclicComponents)
	if r.OperationalInvariantViolations != 0 || r.CustodyViolations != 0 || r.UnsafeStates != 0 {
		t.Fatal("graph contains descriptor, operational, custody or authority violations")
	}
	if r.SafeCompleteStates == 0 || r.CannotCompleteSafely == 0 || r.ShortestCannotComplete == nil ||
		r.CyclicComponents == 0 || r.UnfairRetry == nil {
		t.Fatal("graph lacks genuine completion, a losing safe state or a retry cycle")
	}
	if w := r.ShortestCannotComplete; len(w.Trace) != 1 || w.Trace[0].Kind != "revoke" || w.State.Unsafe != 0 {
		t.Fatal("expected immediate b revocation as the shortest universal completion counterexample")
	}
	for i, n := range g.nodes {
		ba := n.state.Authority[base.BA]
		want := n.state.Unsafe == 0 && ba.Decision != Refused &&
			(!n.state.Revoked || n.state.Base.Obligations[base.BA].Applied || ba.Permit)
		if n.safe != want || n.state.Authority[base.AB].Decision == Refused {
			t.Fatalf("completion predicate or AB own-root counterexample at state %d: %+v", i, n.state)
		}
	}
	for _, tc := range []struct {
		name        string
		kind        string
		d           base.Obligation
		canComplete bool
	}{
		{name: "ba-application-after-revocation", kind: "apply", d: base.BA, canComplete: true},
		{name: "revoked-ba-opportunity-crashed", kind: "crash", d: base.NoObligation, canComplete: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := r.Cuts[tc.name]
			if w == nil || len(w.Trace) == 0 {
				t.Fatal("missing temporal witness cut")
			}
			before := Initial(c)
			for i, e := range w.Trace {
				if i == len(w.Trace)-1 {
					if e.Kind != tc.kind || e.Obligation != tc.d || e.Actor != base.A ||
						!before.Revoked || !before.Authority[base.BA].Permit || e.Revision != 1 {
						t.Fatalf("incorrect target for revoked BA opportunity: %+v", e)
					}
				}
				found := false
				for _, tr := range Successors(c, before) {
					if tr.Event == e {
						before = tr.State
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("full descriptor fails to reconstruct event %d: %+v", i, e)
				}
			}
			if before != w.State || w.State.Authority[base.BA].Permit || w.CanCompleteSafely != tc.canComplete {
				t.Fatal("witness target did not consume or clear the BA permit as expected")
			}
		})
	}
}

func separateFirstAllow(t *testing.T) (State, Event) {
	t.Helper()
	c := GrantRevocationConfig()
	s := revocationSteps(t, Initial(c), base.BA, "commit", "retain-grant", "send-request", "admit")
	return step(t, c, s, "authorize", base.BA)
}

func revocationSteps(t *testing.T, s State, d base.Obligation, kinds ...string) State {
	t.Helper()
	for _, kind := range kinds {
		s, _ = step(t, GrantRevocationConfig(), s, kind, d)
	}
	return s
}
