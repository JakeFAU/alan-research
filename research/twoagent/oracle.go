package twoagent

import "fmt"

// CheckTrace is a separately implemented causal ledger for finite prefixes. It
// never calls Successors, Explore, Initial or the model's state-safety checker.
// Authenticity is derived from the one application and retained/message copies;
// a historical actual-result fact does not permit re-sending a deleted receipt.
func CheckTrace(p Profile, events []Event) error {
	if p != DeliveryACK && p != DurableACK {
		return fmt.Errorf("unknown profile %q", p)
	}
	var committed, down, crashed, admitted, applied, consumed bool
	var requestQueued, ackQueued bool
	var actual, retained, queued, received, durable Value
	for i, e := range events {
		fail := func(reason string) error { return fmt.Errorf("event %d (%s): %s", i, e.Kind, reason) }
		if e.Delegation != FixedDelegation() || e.Grant != FixedGrant {
			return fail("changed delegation, request or fixed grant")
		}
		if e.Duplicate && e.Kind != Admit && e.Kind != Deduplicate && e.Kind != ReceiveResult && e.Kind != Release {
			return fail("duplication flag outside transport delivery")
		}
		usesValue := e.Kind == Apply || e.Kind == SendResult || e.Kind == ReceiveResult || e.Kind == Record || e.Kind == Consume
		if usesValue {
			if e.Result != R1 && e.Result != R2 {
				return fail("invalid opaque result")
			}
		} else if e.Result != "" {
			return fail("unexpected payload; ACK retains identity only")
		}
		aStep := false
		switch e.Kind {
		case Commit, SendRequest, ReceiveResult, Record, Consume, SendACK, Crash:
			aStep = true
		}
		if aStep && down {
			return fail("A step while down")
		}
		switch e.Kind {
		case Commit:
			if committed {
				return fail("duplicate commit")
			}
			committed = true
		case SendRequest:
			if !committed || requestQueued {
				return fail("request has no committed source or channel occupied")
			}
			requestQueued = true
		case Admit, Deduplicate:
			if !requestQueued || (e.Kind == Admit && admitted) || (e.Kind == Deduplicate && !admitted) {
				return fail("admission/deduplication lacks matching original request")
			}
			admitted = true
			if !e.Duplicate {
				requestQueued = false
			}
		case Apply:
			if !admitted || applied {
				return fail("application lacks admission or repeats mutation")
			}
			applied = true
			actual = e.Result
			retained = e.Result
		case SendResult:
			if retained == "" || e.Result != retained || queued != "" {
				return fail("send lacks retained authentic receipt or channel occupied")
			}
			queued = e.Result
		case ReceiveResult:
			if queued == "" || e.Result != queued {
				return fail("receipt lacks matching authentic channel bytes")
			}
			received = e.Result
			if !e.Duplicate {
				queued = ""
			}
		case Record:
			if received == "" || e.Result != received || durable != "" {
				return fail("record lacks matching volatile receipt or already recorded")
			}
			durable = e.Result
		case Consume:
			if durable == "" || e.Result != durable || consumed {
				return fail("consume lacks matching durable receipt or repeats consumption")
			}
			consumed = true
		case SendACK:
			ready := durable != ""
			if p == DeliveryACK {
				ready = ready || received != ""
			}
			if !ready || ackQueued {
				return fail("ACK lacks profile-specific causal evidence or channel occupied")
			}
			ackQueued = true
		case Release:
			if !ackQueued {
				return fail("release lacks authenticated ACK")
			}
			retained = ""
			if !e.Duplicate {
				ackQueued = false
			}
		case Crash:
			if crashed {
				return fail("second crash")
			}
			crashed = true
			down = true
			received = ""
		case Recover:
			if !down {
				return fail("recovery while A is available")
			}
			down = false
		default:
			return fail("unknown event")
		}
		for _, v := range []Value{retained, queued, received, durable} {
			if v != "" && (!applied || v != actual) {
				return fail("fabricated or changed result")
			}
		}
	}
	return nil
}
