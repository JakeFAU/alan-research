package effectobservation

import "fmt"

// ControlResult retains full descriptors and the actual semantic verdict.
type ControlResult struct {
	Name        string  `json:"name"`
	Expectation string  `json:"expectation"`
	Passed      bool    `json:"passed"`
	Events      []Event `json:"events"`
	Rejection   string  `json:"rejection,omitempty"`
	Ledger      *Ledger `json:"ledger,omitempty"`
}

// AuditControls runs hostile traces alongside lawful controls. Accepting no
// events is not sufficient: positive controls require specific preserved facts.
func AuditControls() ([]ControlResult, error) {
	type control struct {
		name   string
		events []Event
		// A nil check requires rejection; otherwise acceptance and this predicate.
		check func(Ledger) bool
	}
	controls := []control{
		{"gate-is-not-contact", controlTrace(), func(l Ledger) bool { return l.Opportunity && !l.PacketSent && !l.Contact && l.World == NoResult }},
		{"contact-is-not-application", controlTrace(Launch, Ingress), func(l Ledger) bool { return l.Contact && l.World == NoResult && l.Stored == NoResult }},
		{"application-is-not-retention", controlTrace(Launch, Ingress, Apply), func(l Ledger) bool { return l.World == Value1 && l.Inbox == NoResult && l.Stored == NoResult }},
		{"loss-preserves-packet", controlTrace(Launch, RevokeB, LoseOwner, Ingress, Apply), func(l Ledger) bool { return l.Dead && l.PacketSent && l.World == Value1 && l.Stored == NoResult }},
		{"loss-preserves-world", controlTrace(Launch, Ingress, Apply, Response, LoseOwner, Reopen), func(l Ledger) bool { return l.World == Value1 && l.Inbox == NoResult && l.Stored == NoResult }},
		{"old-permit-survives-revocation", controlTrace(RevokeB, Launch), func(l Ledger) bool { return l.PacketSent && l.Revoked && l.CValid }},
		{"historical-allowance-survives", controlTrace(RevokeB, LoseOwner, Reopen, FreshRefuse), func(l Ledger) bool {
			v := l.StoreView()
			return v.FirstAllow && v.FirstGateAllowed && v.AcceptedGrant && v.Debit == 1 && v.ControlRoot == RootC && v.FreshRefused
		}},
		{"retained-result-survives", controlTrace(Launch, Ingress, Apply, Response, Record, RevokeB, LoseOwner, Reopen, Consume), func(l Ledger) bool { return l.Stored == Value1 && l.Used == Value1 && l.CValid }},
		{"refusal-does-not-erase-application", controlTrace(Launch, Ingress, Apply, LoseOwner, Reopen, RevokeB, FreshRefuse), func(l Ledger) bool { return l.Refused && l.World == Value1 && l.Stored == NoResult }},
		{"gate-as-contact", controlTrace(Ingress), nil},
		{"contact-as-application", controlTrace(Launch, Ingress, Response), nil},
		{"application-as-consumable-result", controlTrace(Launch, Ingress, Apply, Consume), nil},
		{"forensic-result-import", controlTrace(Launch, Ingress, Apply, LoseOwner, Reopen, Record), nil},
		{"old-response-revived", controlTrace(Launch, Ingress, Apply, LoseOwner, Reopen, Response), nil},
		{"old-permit-revived", controlTrace(LoseOwner, Reopen, Launch), nil},
		{"fresh-allowed-retry", controlTrace(Launch, LoseOwner, Reopen, RevokeB, FreshRefuse, Launch), nil},
	}
	wrongRoot := controlTrace(RevokeB)
	wrongRoot[0].Root = RootC
	controls = append(controls, control{"b-c-aliasing", wrongRoot, nil})
	wrongValue := controlTrace(Launch, Ingress, Apply, Response)
	wrongValue[3].Value = Value2
	controls = append(controls, control{"wrong-authentic-response", wrongValue, nil})
	stale := controlTrace(LoseOwner, Reopen, RevokeB, FreshRefuse)
	stale[3].Incarnation = 1
	controls = append(controls, control{"stale-fresh-refusal", stale, nil})
	irrelevant := controlTrace(Launch)
	irrelevant[0].Value = Value1
	controls = append(controls, control{"irrelevant-value", irrelevant, nil})
	results := make([]ControlResult, 0, len(controls))
	failed := false
	for _, c := range controls {
		l, err := AuditTrace(c.events)
		r := ControlResult{Name: c.name, Events: c.events, Expectation: "reject"}
		if err != nil {
			r.Rejection = err.Error()
		} else {
			r.Ledger = &l
		}
		if c.check == nil {
			r.Passed = err != nil
		} else {
			r.Expectation = "accept and preserve stated facts"
			r.Passed = err == nil && c.check(l)
		}
		failed = failed || !r.Passed
		results = append(results, r)
	}
	if failed {
		return results, fmt.Errorf("one or more semantic controls failed")
	}
	return results, nil
}

// This corpus constructor uses only the declared descriptor schema.
func controlTrace(kinds ...EventKind) []Event {
	events := make([]Event, 0, len(kinds))
	inc := uint8(1)
	for _, kind := range kinds {
		if kind == Reopen {
			inc = 2
		}
		e := Event{Kind: kind, Binding: ledgerBinding(), Actor: ActorA, Incarnation: inc}
		switch kind {
		case Ingress, Apply:
			e.Actor, e.Incarnation = ActorEndpoint, 0
		case Consume, RevokeB:
			e.Actor = ActorB
		case LoseOwner, Reopen:
			e.Actor = ActorOwner
		}
		if kind == RevokeB || kind == FreshRefuse {
			e.Root = RootB
		}
		if kind == Apply || kind == Response || kind == Record || kind == Consume {
			e.Value = Value1
		}
		events = append(events, e)
	}
	return events
}
