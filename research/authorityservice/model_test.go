package authorityservice

import (
	base "github.com/JakeFAU/alan-research/research/mutualobligations"
	"testing"
)

func TestApplyRequiresFreshAuthorityAdmission(t *testing.T) {
	c := Config{Policy: "S2", Scenario: "Missing"}
	s := Initial(c)
	for _, kind := range []string{"commit", "send-request", "admit"} {
		found := false
		for _, e := range Successors(c, s) {
			if e.Event.Kind == kind && e.Event.Obligation == base.BA {
				s = e.State
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing %s", kind)
		}
	}
	for _, e := range Successors(c, s) {
		if e.Event.Kind == "apply" && e.Event.Obligation == base.BA {
			t.Fatal("unpermitted BA service applies before any authority decision or physical admission")
		}
	}
}

func step(t *testing.T, c Config, s State, kind string, d base.Obligation) (State, Event) {
	t.Helper()
	for _, e := range Successors(c, s) {
		if e.Event.Kind == kind && e.Event.Obligation == d {
			return e.State, e.Event
		}
	}
	t.Fatalf("missing %s/%d in %+v", kind, d, s)
	return s, Event{}
}
func TestDecisiveCandidate(t *testing.T) {
	for _, tc := range []struct {
		policy, scenario string
		decision         Decision
		path             Path
		unsafe           Violation
	}{
		{"S1", "Explicit", Refused, NoPath, 0}, {"S2", "Explicit", Allowed, Delegated, 0},
		{"S2", "Missing", Refused, NoPath, 0}, {"S3", "Missing", Allowed, Ambient, ConsentViolation},
		{"S3", "ScopeMix", Refused, NoPath, 0}, {"Cartesian", "ScopeMix", Allowed, Cartesian, TupleInvention},
	} {
		t.Run(tc.policy+"/"+tc.scenario, func(t *testing.T) {
			c := Config{Policy: tc.policy, Scenario: tc.scenario}
			s := Initial(c)
			trace := []Event{}
			for _, k := range []string{"commit", "send-request", "admit", "authorize"} {
				var e Event
				s, e = step(t, c, s, k, base.BA)
				trace = append(trace, e)
			}
			a := s.Authority[base.BA]
			if a.Decision != tc.decision || a.Path != tc.path || s.Unsafe != tc.unsafe {
				t.Fatalf("got %+v unsafe=%d", a, s.Unsafe)
			}
			t.Logf("four-step decisive BA prefix: %+v", trace)
		})
	}
}

func TestCompleteGraphs(t *testing.T) {
	if testing.Short() {
		t.Skip("complete graph qualification")
	}
	for _, c := range Configurations() {
		t.Run(c.Name(), func(t *testing.T) {
			r := Explore(c)
			t.Logf("states=%d arcs=%d unsafe=%d raw=%d safe=%d cannot=%d cycles=%d", r.States, r.Arcs, r.UnsafeStates, r.CanComplete, r.CanCompleteSafely, r.CannotComplete, r.CyclicStates)
			if r.OperationalInvariantViolations+r.CustodyViolations != 0 {
				t.Fatalf("apparatus/custody violations: %+v", r)
			}
			wantComplete := c.Scenario == "Local" || c.Scenario == "Shared" || c.Scenario == "LateGrant" || c.Scenario == "Revocable" || c.Policy != "S1" && c.Scenario == "Explicit" || c.Policy == "S3" && c.Scenario == "Missing" || c.Policy == "Cartesian"
			if (r.CompleteStates > 0) != wantComplete {
				t.Fatalf("completion=%d want existence %v", r.CompleteStates, wantComplete)
			}
			wantUnsafe := c.Policy == "S3" && c.Scenario == "Missing" || c.Policy == "Cartesian"
			if (r.UnsafeStates > 0) != wantUnsafe {
				t.Fatalf("unsafe=%d want %v", r.UnsafeStates, wantUnsafe)
			}
			if wantUnsafe && (r.SafeCompleteStates != 0 || r.CanCompleteSafely != 0) {
				t.Fatal("unsafe completion was laundered")
			}
			if c.Scenario == "Local" || c.Scenario == "Shared" || c.Policy != "S1" && c.Scenario == "Explicit" {
				if r.CannotCompleteSafely != 0 {
					t.Fatal("positive profile lacks safe suffix")
				}
			}
			if r.AppliedStates == 0 || r.ReleasedStates == 0 || r.CyclicStates == 0 {
				t.Fatal("vacuous custody/cycle control")
			}
		})
	}
}

func admittedBA(t *testing.T, c Config) State {
	t.Helper()
	s := Initial(c)
	for _, k := range []string{"commit", "send-request", "admit", "authorize", "dispatch-admit"} {
		s, _ = step(t, c, s, k, base.BA)
	}
	return s
}
func noStep(t *testing.T, c Config, s State, kind string, d base.Obligation) {
	t.Helper()
	for _, e := range Successors(c, s) {
		if e.Event.Kind == kind && e.Event.Obligation == d {
			t.Fatalf("unexpected %s/%d", kind, d)
		}
	}
}
func actorStep(t *testing.T, c Config, s State, kind string, actor base.Agent) State {
	t.Helper()
	for _, e := range Successors(c, s) {
		if e.Event.Kind == kind && e.Event.Actor == actor {
			return e.State
		}
	}
	t.Fatalf("missing %s by %d", kind, actor)
	return s
}
func TestRevocationTemporalBoundaries(t *testing.T) {
	c := Config{Policy: "S2", Scenario: "Revocable"}
	t.Run("rejected-inline-issuance-still-commits", func(t *testing.T) {
		s := actorStep(t, c, Initial(c), "revoke", base.B)
		s, e := step(t, c, s, "commit", base.BA)
		if !s.Base.Obligations[base.BA].Committed || s.Grant != RejectedGrant || !e.HasGrant || e.Grant.Accepted {
			t.Fatal("revocation suppressed commit or invented grant")
		}
	})
	t.Run("authorization-does-not-reserve-admission", func(t *testing.T) {
		s := Initial(c)
		for _, k := range []string{"commit", "send-request", "admit", "authorize"} {
			s, _ = step(t, c, s, k, base.BA)
		}
		s = actorStep(t, c, s, "revoke", base.B)
		s, _ = step(t, c, s, "dispatch-refuse", base.BA)
		if s.Authority[base.BA].Decision != Allowed {
			t.Fatal("historical allow rewritten")
		}
		noStep(t, c, s, "dispatch-admit", base.BA)
		noStep(t, c, s, "apply", base.BA)
	})
	t.Run("admitted-opportunity-survives-revocation", func(t *testing.T) {
		s := actorStep(t, c, admittedBA(t, c), "revoke", base.B)
		s, _ = step(t, c, s, "apply", base.BA)
		if !s.Base.Obligations[base.BA].Applied || s.Authority[base.BA].Permit {
			t.Fatal("pre-revocation opportunity lost or not consumed")
		}
	})
	t.Run("server-crash-clears-opportunity", func(t *testing.T) {
		s := actorStep(t, c, admittedBA(t, c), "revoke", base.B)
		s = actorStep(t, c, s, "crash", base.A)
		if s.Authority[base.BA].Permit {
			t.Fatal("crash retained physical opportunity")
		}
		s = actorStep(t, c, s, "recover", base.A)
		noStep(t, c, s, "dispatch-admit", base.BA)
		noStep(t, c, s, "apply", base.BA)
		s, _ = step(t, c, s, "dispatch-refuse", base.BA)
	})
	t.Run("unrevoked-path-readmits-new-activation", func(t *testing.T) {
		s := actorStep(t, c, admittedBA(t, c), "crash", base.A)
		s = actorStep(t, c, s, "recover", base.A)
		s, e := step(t, c, s, "dispatch-admit", base.BA)
		if e.Activation != 1 || s.Authority[base.BA].PermitActivation != 1 {
			t.Fatal("fresh admission has stale activation")
		}
	})
	t.Run("shared-root-is-independent", func(t *testing.T) {
		shared := Config{Policy: "S2", Scenario: "Shared"}
		s := Initial(shared)
		s.Revoked = true
		for d := base.AB; d <= base.BA; d++ {
			if choose(shared, s, d) != OwnRoot || !pathCurrent(shared, s, d, OwnRoot) {
				t.Fatal("b revocation removed shared r")
			}
		}
	})
}
func TestLateGrantCannotRewriteRefusal(t *testing.T) {
	c := Config{Policy: "S2", Scenario: "LateGrant"}
	s := Initial(c)
	for _, k := range []string{"commit", "send-request", "admit", "authorize", "retain-grant"} {
		s, _ = step(t, c, s, k, base.BA)
	}
	s = actorStep(t, c, s, "crash", base.A)
	s = actorStep(t, c, s, "recover", base.A)
	if s.Grant != RetainedGrant || s.Authority[base.BA].Decision != Refused {
		t.Fatal("late grant lost retained refusal")
	}
	noStep(t, c, s, "authorize", base.BA)
	noStep(t, c, s, "dispatch-admit", base.BA)
	noStep(t, c, s, "apply", base.BA)
}
func TestReleasedCustodySurvivesBothCrashes(t *testing.T) {
	c := Config{Policy: "S2", Scenario: "Explicit"}
	s := admittedBA(t, c)
	for _, k := range []string{"apply", "send-result", "receive-result", "record", "send-ack", "release"} {
		s, _ = step(t, c, s, k, base.BA)
	}
	s = actorStep(t, c, s, "crash", base.A)
	s = actorStep(t, c, s, "crash", base.B)
	if s.Base.Obligations[base.BA].Retained != base.NoValue || s.Base.Obligations[base.BA].Recorded == base.NoValue {
		t.Fatal("released result lost durable custody")
	}
	s = actorStep(t, c, s, "recover", base.A)
	s = actorStep(t, c, s, "recover", base.B)
	s, _ = step(t, c, s, "consume", base.BA)
	for _, k := range []string{"commit", "send-request", "admit", "authorize", "dispatch-admit", "apply", "send-result", "receive-result", "record", "consume"} {
		s, _ = step(t, c, s, k, base.AB)
	}
	if !Complete(s) || s.Unsafe != 0 {
		t.Fatal("post-interruption originals not safely consumed")
	}
}
