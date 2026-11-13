package effectobservation_test

import (
	"testing"

	m "github.com/JakeFAU/alan-research/research/effectobservation"
)

func step(t *testing.T, state m.State, kind m.EventKind, value m.Result) m.State {
	t.Helper()
	for _, edge := range m.Successors(state) {
		if edge.Event.Kind == kind && edge.Event.Value == value {
			return edge.State
		}
	}
	t.Fatalf("missing %s %q from %+v", kind, value, state)
	return m.State{}
}

func absent(t *testing.T, state m.State, kinds ...m.EventKind) {
	t.Helper()
	for _, edge := range m.Successors(state) {
		for _, kind := range kinds {
			if edge.Event.Kind == kind {
				t.Fatalf("unexpected %+v from %+v", edge.Event, state)
			}
		}
	}
}

func apply(t *testing.T, state m.State, value m.Result) m.State {
	t.Helper()
	state = step(t, state, m.Launch, m.NoResult)
	state = step(t, state, m.Ingress, m.NoResult)
	return step(t, state, m.Apply, value)
}

func TestSeparateStages(t *testing.T) {
	initial := m.Initial()
	absent(t, initial, m.Ingress, m.Apply, m.Response, m.Record, m.Consume, m.FreshRefuse)
	launched := step(t, initial, m.Launch, m.NoResult)
	if launched.Permit || !launched.Launched || m.ObserveDossier(launched) != m.ObserveDossier(initial) {
		t.Fatalf("launch should only use the permit and transfer the packet: %+v", launched)
	}
	absent(t, launched, m.Launch, m.Apply, m.Response, m.Record, m.Consume)
	ingressed := step(t, launched, m.Ingress, m.NoResult)
	if !m.ObserveDossier(ingressed).Ingress || m.ObserveDossier(ingressed).Applied != m.NoResult {
		t.Fatalf("ingress should be visible without application: %+v", ingressed)
	}
	absent(t, ingressed, m.Ingress, m.Response, m.Record, m.Consume)
	for _, value := range []m.Result{m.Value1, m.Value2} {
		t.Run(string(value), func(t *testing.T) {
			applied := step(t, ingressed, m.Apply, value)
			if applied.Applied != value || applied.Response != m.NoResult ||
				applied.Retained != m.NoResult || applied.Consumed != m.NoResult {
				t.Fatalf("application should introduce only world truth: %+v", applied)
			}
			absent(t, applied, m.Apply, m.Record, m.Consume)
			response := step(t, applied, m.Response, value)
			if response.Response != value || m.ObserveStore(response) != m.ObserveStore(initial) {
				t.Fatalf("successful arrival should add volatile custody only: %+v", response)
			}
			absent(t, response, m.Response, m.Consume)
			retained := step(t, response, m.Record, value)
			if retained.Response != m.NoResult || m.ObserveStore(retained).Retained != value || retained.Consumed != m.NoResult {
				t.Fatalf("record should move the response into durable custody: %+v", retained)
			}
			absent(t, retained, m.Response, m.Record)
			consumed := step(t, retained, m.Consume, value)
			if m.ObserveStore(consumed).Consumed != value || consumed.Retained != value {
				t.Fatalf("consume should preserve the retained result and add its marker: %+v", consumed)
			}
			absent(t, consumed, m.Consume)
		})
	}
}

func TestLossPreservesExternalPacket(t *testing.T) {
	launched := step(t, m.Initial(), m.Launch, m.NoResult)
	revoked := step(t, launched, m.RevokeB, m.NoResult)
	stopped := step(t, revoked, m.LoseOwner, m.NoResult)
	if stopped.Owner != m.OwnerStopped || !stopped.Launched || stopped.Permit || stopped.Ingress {
		t.Fatalf("loss should preserve the launched packet outside the owner: %+v", stopped)
	}
	if cut := m.ObserveDossier(stopped); cut.Ingress || cut.Applied != m.NoResult {
		t.Fatalf("early cut should have no ingress or application: %+v", cut)
	}
	ingressed := step(t, stopped, m.Ingress, m.NoResult)
	applied := step(t, ingressed, m.Apply, m.Value2)
	if !applied.BRevoked || applied.Applied != m.Value2 || applied.Owner != m.OwnerStopped {
		t.Fatalf("endpoint should apply while owner remains stopped: %+v", applied)
	}
	absent(t, applied, m.Launch, m.Response, m.Record, m.Consume, m.RevokeB)
	unlaunched := step(t, m.Initial(), m.LoseOwner, m.NoResult)
	if unlaunched.Permit || m.ObserveDossier(unlaunched).Ingress {
		t.Fatalf("loss should clear an unlaunched opportunity: %+v", unlaunched)
	}
	absent(t, unlaunched, m.Launch, m.Ingress, m.Apply)
	absent(t, step(t, unlaunched, m.Reopen, m.NoResult), m.Launch, m.Ingress, m.Apply)
}

func TestLateResponseCannotReachReopenedOwner(t *testing.T) {
	applied := apply(t, m.Initial(), m.Value2)
	arrived := step(t, applied, m.Response, m.Value2)
	for _, source := range []m.State{applied, arrived} {
		stopped := step(t, source, m.LoseOwner, m.NoResult)
		reopened := step(t, stopped, m.Reopen, m.NoResult)
		if reopened.Response != m.NoResult || reopened.Applied != m.Value2 ||
			reopened.Owner != m.OwnerReopened || reopened.Activation != 2 {
			t.Fatalf("loss should clear old volatile custody without changing endpoint truth: %+v", reopened)
		}
		absent(t, reopened, m.Launch, m.Response, m.Record, m.Consume, m.LoseOwner, m.Reopen)
	}
}

func TestExactRootAndHistoricalAuthority(t *testing.T) {
	initial := m.Initial()
	if !initial.AcceptedGrant || !initial.FirstAllow || initial.Debit != 1 ||
		!initial.FirstGateAllowed || initial.Binding.Root != m.RootB || initial.ControlRoot != m.RootC {
		t.Fatalf("initial cut should include fixed admission and distinct roots: %+v", initial)
	}
	revoked := step(t, initial, m.RevokeB, m.NoResult)
	want := initial
	want.BRevoked = true
	if revoked != want {
		t.Fatalf("revocation should change only exact b validity: got %+v want %+v", revoked, want)
	}
	absent(t, revoked, m.RevokeB, m.FreshRefuse)
	applied := apply(t, revoked, m.Value1)
	if !applied.BRevoked || !applied.FirstAllow || applied.Applied != m.Value1 {
		t.Fatalf("surviving opportunity should still apply after b revocation: %+v", applied)
	}
}

func TestRefusalAfterUnretainedApplication(t *testing.T) {
	refuse := func(source m.State) m.State {
		t.Helper()
		stopped := step(t, source, m.LoseOwner, m.NoResult)
		absent(t, stopped, m.FreshRefuse, m.RevokeB)
		reopened := step(t, stopped, m.Reopen, m.NoResult)
		absent(t, reopened, m.FreshRefuse)
		revoked := step(t, reopened, m.RevokeB, m.NoResult)
		refused := step(t, revoked, m.FreshRefuse, m.NoResult)
		want := revoked
		want.FreshRefused = true
		if refused != want {
			t.Fatalf("fresh refusal should only record its disposition: got %+v want %+v", refused, want)
		}
		absent(t, refused, m.FreshRefuse, m.Launch, m.Response, m.Record, m.Consume)
		return refused
	}
	unapplied := refuse(m.Initial())
	first := refuse(apply(t, m.Initial(), m.Value1))
	second := refuse(apply(t, m.Initial(), m.Value2))
	for _, state := range []m.State{first, second} {
		if m.ObserveStore(state) != m.ObserveStore(unapplied) {
			t.Fatalf("external application leaked into the local observation: %+v", state)
		}
	}
	if m.ObserveDossier(unapplied).Applied != m.NoResult || m.ObserveDossier(first).Applied != m.Value1 ||
		m.ObserveDossier(second).Applied != m.Value2 {
		t.Fatal("dossier should distinguish no application and both authentic values")
	}
}

func TestRetainedResultSurvivesLoss(t *testing.T) {
	applied := apply(t, m.Initial(), m.Value2)
	arrived := step(t, applied, m.Response, m.Value2)
	retained := step(t, arrived, m.Record, m.Value2)
	stopped := step(t, retained, m.LoseOwner, m.NoResult)
	absent(t, stopped, m.Consume)
	reopened := step(t, stopped, m.Reopen, m.NoResult)
	revoked := step(t, reopened, m.RevokeB, m.NoResult)
	absent(t, revoked, m.FreshRefuse)
	consumed := step(t, revoked, m.Consume, m.Value2)
	if consumed.Retained != m.Value2 || consumed.Consumed != m.Value2 || consumed.Activation != 2 {
		t.Fatalf("new owner should consume the surviving result: %+v", consumed)
	}
}
