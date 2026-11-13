package effectobservation

import "testing"

// These descriptors deliberately do not use Initial, Successors or Replay.
func auditEvent(kind EventKind, incarnation uint8, value Result) Event {
	e := Event{Kind: kind, Binding: Binding{Original: "BA/1", SourceGeneration: 1, SourceRevision: 1,
		Grant: "grant:BA:b", Root: RootB, OriginalActivation: 1, Wire: "wire:BA:1"}, Incarnation: incarnation, Value: value}
	switch kind {
	case Launch, Response, Record, FreshRefuse:
		e.Actor = ActorA
	case Ingress, Apply:
		e.Actor = ActorEndpoint
	case Consume, RevokeB:
		e.Actor = ActorB
	case LoseOwner, Reopen:
		e.Actor = ActorOwner
	}
	if kind == RevokeB || kind == FreshRefuse {
		e.Root = RootB
	}
	return e
}

func auditFixture(kinds ...EventKind) []Event {
	var events []Event
	inc := uint8(1)
	for _, kind := range kinds {
		if kind == Reopen {
			inc = 2
		}
		actorInc := inc
		if kind == Ingress || kind == Apply {
			actorInc = 0
		}
		value := NoResult
		if kind == Apply || kind == Response || kind == Record || kind == Consume {
			value = Value1
		}
		events = append(events, auditEvent(kind, actorInc, value))
	}
	return events
}

func TestLedgerLawfulControls(t *testing.T) {
	tests := []struct {
		name  string
		trace []Event
		check func(Ledger) bool
	}{
		{"gate-is-not-contact", nil, func(l Ledger) bool { return l.Opportunity && !l.Contact && l.World == NoResult && l.CValid }},
		{"contact-is-not-application", auditFixture(Launch, Ingress), func(l Ledger) bool { return l.Contact && l.World == NoResult && l.Stored == NoResult }},
		{"application-is-not-retention", auditFixture(Launch, Ingress, Apply), func(l Ledger) bool { return l.World == Value1 && l.Stored == NoResult && l.Inbox == NoResult }},
		{"old-permit-survives-revoke", auditFixture(RevokeB, Launch), func(l Ledger) bool { return l.PacketSent && l.Revoked && l.CValid }},
		{"packet-survives-loss", auditFixture(Launch, RevokeB, LoseOwner, Ingress, Apply), func(l Ledger) bool {
			return l.Dead && l.PacketSent && l.Contact && l.World == Value1 && l.Stored == NoResult
		}},
		{"world-survives-loss", auditFixture(Launch, Ingress, Apply, Response, LoseOwner, Reopen), func(l Ledger) bool {
			return l.World == Value1 && l.Inbox == NoResult && l.Stored == NoResult && !l.Opportunity
		}},
		{"retained-recovery", auditFixture(Launch, Ingress, Apply, Response, Record, RevokeB, LoseOwner, Reopen, Consume), func(l Ledger) bool { return l.Stored == Value1 && l.Used == Value1 && l.Restarted && l.CValid }},
		{"refusal-with-world-result", auditFixture(Launch, Ingress, Apply, LoseOwner, Reopen, RevokeB, FreshRefuse), func(l Ledger) bool { return l.Refused && l.World == Value1 && l.Stored == NoResult }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := AuditTrace(tt.trace)
			if err != nil {
				t.Fatal(err)
			}
			if !tt.check(l) {
				t.Fatalf("lawful control failed: %+v", l)
			}
		})
	}
}

func TestLedgerRejectsInvalidCustody(t *testing.T) {
	tests := []struct {
		name  string
		trace []Event
	}{
		{"gate-as-contact", auditFixture(Apply)},
		{"gate-as-launch", auditFixture(Ingress)},
		{"contact-as-application", auditFixture(Launch, Ingress, Response)},
		{"application-as-retention", auditFixture(Launch, Ingress, Apply, Consume)},
		{"forensic-result-import", auditFixture(Launch, Ingress, Apply, Record)},
		{"revived-old-response", auditFixture(Launch, Ingress, Apply, LoseOwner, Reopen, Response)},
		{"revived-old-permit", auditFixture(LoseOwner, Reopen, Launch)},
		{"fresh-allowed-retry", auditFixture(Launch, LoseOwner, Reopen, RevokeB, FreshRefuse, Launch)},
		{"stopped-consumption", auditFixture(Launch, Ingress, Apply, Response, Record, LoseOwner, Consume)},
		{"stopped-revoke", auditFixture(LoseOwner, RevokeB)},
		{"second-loss", auditFixture(LoseOwner, Reopen, LoseOwner)},
		{"premature-refusal", auditFixture(RevokeB, FreshRefuse)},
		{"unrevoked-refusal", auditFixture(LoseOwner, Reopen, FreshRefuse)},
		{"retained-refusal", auditFixture(Launch, Ingress, Apply, Response, Record, LoseOwner, Reopen, RevokeB, FreshRefuse)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := AuditTrace(tt.trace); err == nil {
				t.Fatal("accepted invalid custody transition")
			}
		})
	}
}

func TestLedgerMalformedDescriptors(t *testing.T) {
	graph := Explore()
	mutations := []struct {
		name   string
		change func(*Event)
	}{
		{"kind", func(e *Event) { e.Kind = "fresh-allow" }},
		{"original", func(e *Event) { e.Binding.Original = "BA/2" }},
		{"generation", func(e *Event) { e.Binding.SourceGeneration++ }},
		{"revision", func(e *Event) { e.Binding.SourceRevision++ }},
		{"grant", func(e *Event) { e.Binding.Grant = "grant:BA:c" }},
		{"binding-root", func(e *Event) { e.Binding.Root = RootC }},
		{"original-activation", func(e *Event) { e.Binding.OriginalActivation++ }},
		{"wire", func(e *Event) { e.Binding.Wire = "wire:BA:2" }},
		{"actor", func(e *Event) { e.Actor = "forger" }},
		{"stale-incarnation", func(e *Event) { e.Incarnation = 9 }},
		{"wrong-or-irrelevant-root", func(e *Event) { e.Root = RootC }},
		{"wrong-or-irrelevant-value", func(e *Event) { e.Value = "forged" }},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			for from, arcs := range graph.Arcs {
				for _, arc := range arcs {
					e := arc.Event
					mutation.change(&e)
					if _, err := AuditTrace(appendEvent(graph.Prefixes[from], e)); err == nil {
						t.Fatalf("accepted malformed %s from state %d", e.Kind, from)
					}
				}
			}
		})
	}
	for from, arcs := range graph.Arcs {
		for _, arc := range arcs {
			if _, err := AuditTrace(appendEvent(graph.Prefixes[from], arc.Event)); err != nil {
				t.Fatalf("lawful %s: %v", arc.Event.Kind, err)
			}
			for inc := uint8(0); inc <= 2; inc++ {
				if inc == arc.Event.Incarnation {
					continue
				}
				e := arc.Event
				e.Incarnation = inc
				if _, err := AuditTrace(appendEvent(graph.Prefixes[from], e)); err == nil {
					t.Fatalf("accepted wrong known incarnation %d for %s", inc, e.Kind)
				}
			}
			e := arc.Event
			if e.Value == NoResult {
				e.Value = Value1
			} else {
				e.Value = NoResult
			}
			if _, err := AuditTrace(appendEvent(graph.Prefixes[from], e)); err == nil {
				t.Fatalf("accepted missing or irrelevant value for %s", e.Kind)
			}
			e = arc.Event
			if e.Root == "" {
				e.Root = RootB
			} else {
				e.Root = ""
			}
			if _, err := AuditTrace(appendEvent(graph.Prefixes[from], e)); err == nil {
				t.Fatalf("accepted missing or irrelevant root for %s", e.Kind)
			}
			repeated := appendEvent(appendEvent(graph.Prefixes[from], arc.Event), arc.Event)
			if _, err := AuditTrace(repeated); err == nil {
				t.Fatalf("accepted duplicate %s", arc.Event.Kind)
			}
			if arc.Event.Kind == Response || arc.Event.Kind == Record || arc.Event.Kind == Consume {
				e := arc.Event
				if e.Value == Value1 {
					e.Value = Value2
				} else {
					e.Value = Value1
				}
				if _, err := AuditTrace(appendEvent(graph.Prefixes[from], e)); err == nil {
					t.Fatalf("accepted wrong authentic %s value", e.Kind)
				}
			}
		}
	}
}

func TestLedgerHistoricalConservation(t *testing.T) {
	l, err := AuditTrace(auditFixture(RevokeB, Launch, Ingress, Apply, LoseOwner, Reopen, FreshRefuse))
	if err != nil {
		t.Fatal(err)
	}
	v := l.StoreView()
	if !v.AcceptedGrant || !v.FirstAllow || !v.FirstGateAllowed || v.Debit != 1 || v.ControlRoot != RootC || v.Binding != auditEvent(Launch, 1, NoResult).Binding {
		t.Fatalf("rewritten historical allowance or aliased c: %+v", v)
	}
	if v.Retained != NoResult || v.Consumed != NoResult || l.DossierView().Applied != Value1 {
		t.Fatal("forensic application became local custody")
	}
}
