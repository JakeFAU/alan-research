package effectobservation

import "fmt"

// Ledger separately represents event-derived custody and lifecycle facts.
// Historical admission is fixed at the declared initial cut.
type Ledger struct {
	PacketSent  bool   `json:"packet_sent"`
	Contact     bool   `json:"contact"`
	World       Result `json:"world"`
	Inbox       Result `json:"inbox"`
	Stored      Result `json:"stored"`
	Used        Result `json:"used"`
	Opportunity bool   `json:"opportunity"`
	Dead        bool   `json:"dead"`
	Restarted   bool   `json:"restarted"`
	Revoked     bool   `json:"revoked"`
	CValid      bool   `json:"c_valid"`
	Refused     bool   `json:"refused"`
}

// AuditTrace interprets complete event descriptors without operational helpers.
// It assumes the fixed post-admission cut, not arbitrary supplied admission facts.
func AuditTrace(events []Event) (Ledger, error) {
	l := Ledger{Opportunity: true, CValid: true}
	for i, e := range events {
		next, err := l.accept(e)
		if err != nil {
			return Ledger{}, fmt.Errorf("event %d (%s): %w", i, e.Kind, err)
		}
		l = next
	}
	return l, nil
}

func ledgerBinding() Binding {
	return Binding{Original: "BA/1", SourceGeneration: 1, SourceRevision: 1,
		Grant: "grant:BA:b", Root: RootB, OriginalActivation: 1, Wire: "wire:BA:1"}
}

func (l Ledger) activation() uint8 {
	if l.Restarted {
		return 2
	}
	return 1
}

func (l Ledger) accept(e Event) (Ledger, error) {
	fail := func(why string) (Ledger, error) { return Ledger{}, fmt.Errorf("%s", why) }
	if e.Binding != ledgerBinding() {
		return fail("wrong original binding")
	}
	actor, incarnation, root := Actor(""), l.activation(), Root("")
	hasValue := false
	switch e.Kind {
	case Launch:
		actor = ActorA
	case Ingress:
		actor, incarnation = ActorEndpoint, 0
	case Apply:
		actor, incarnation, hasValue = ActorEndpoint, 0, true
	case Response, Record:
		actor, hasValue = ActorA, true
	case Consume:
		actor, hasValue = ActorB, true
	case RevokeB:
		actor, root = ActorB, RootB
	case LoseOwner:
		actor = ActorOwner
	case Reopen:
		actor, incarnation = ActorOwner, 2
	case FreshRefuse:
		actor, root = ActorA, RootB
	default:
		return fail("unknown event kind")
	}
	if e.Actor != actor || e.Incarnation != incarnation || e.Root != root {
		return fail("wrong actor, incarnation or targeted root")
	}
	if (hasValue && e.Value != Value1 && e.Value != Value2) || (!hasValue && e.Value != NoResult) {
		return fail("missing, invalid or irrelevant value")
	}
	original := !l.Dead && !l.Restarted
	switch e.Kind {
	case Launch:
		if !original || !l.Opportunity || l.PacketSent {
			return fail("no original launch opportunity")
		}
		l.Opportunity, l.PacketSent = false, true
	case Ingress:
		if !l.PacketSent || l.Contact {
			return fail("ingress needs one outstanding packet")
		}
		l.Contact = true
	case Apply:
		if !l.Contact || l.World != NoResult {
			return fail("application needs unapplied ingress")
		}
		l.World = e.Value
	case Response:
		if !original || l.World == NoResult || l.Inbox != NoResult || l.Stored != NoResult || e.Value != l.World {
			return fail("response needs authentic value and original live recipient")
		}
		l.Inbox = e.Value
	case Record:
		if !original || l.Inbox == NoResult || l.Stored != NoResult || e.Value != l.Inbox {
			return fail("record needs authentic original-owner custody")
		}
		l.Stored, l.Inbox = e.Value, NoResult
	case Consume:
		if l.Dead || l.Stored == NoResult || l.Used != NoResult || e.Value != l.Stored {
			return fail("consumption needs unconsumed retained value and live owner")
		}
		l.Used = e.Value
	case RevokeB:
		if l.Dead || l.Revoked {
			return fail("revoke needs live owner and current b")
		}
		l.Revoked = true
	case LoseOwner:
		if !original {
			return fail("loss needs original live owner")
		}
		l.Dead, l.Opportunity, l.Inbox = true, false, NoResult
	case Reopen:
		if !l.Dead || l.Restarted {
			return fail("reopen needs first stopped owner")
		}
		l.Dead, l.Restarted = false, true
	case FreshRefuse:
		if !l.Restarted || l.Dead || !l.Revoked || l.Stored != NoResult || l.Refused {
			return fail("fresh refusal needs reopened owner, revoked b and no retained result")
		}
		l.Refused = true
	}
	return l, nil
}

// StoreView derives local evidence. Owner stage is supplied by the death/Open
// protocol; admission constants are hypotheses of the initial cut.
func (l Ledger) StoreView() StoreView {
	owner := OwnerOriginal
	if l.Dead {
		owner = OwnerStopped
	} else if l.Restarted {
		owner = OwnerReopened
	}
	control := Root("")
	if l.CValid {
		control = RootC
	}
	return StoreView{Binding: ledgerBinding(), AcceptedGrant: true, FirstAllow: true,
		Debit: 1, FirstGateAllowed: true, Owner: owner, Activation: l.activation(),
		BRevoked: l.Revoked, ControlRoot: control, FreshRefused: l.Refused,
		Retained: l.Stored, Consumed: l.Used}
}

// DossierView requires complete identity-bound endpoint evidence at the cut.
// Missing/partial evidence cannot be encoded as a negative observation here.
func (l Ledger) DossierView() DossierView {
	return DossierView{StoreView: l.StoreView(), Ingress: l.Contact, Applied: l.World}
}
