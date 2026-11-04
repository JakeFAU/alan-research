package mutualobligations

import "fmt"

// A causal ledger independently implements event preconditions. It shares no
// transition generation, state initialization or safety predicate with the model.
type ledgerEntry struct {
	committed, admitted, applied, consumed      bool
	request, ack                                bool
	actual, receipt, channel, received, durable Value
}
type ledger struct {
	phase   [2]uint8
	entries [2]ledgerEntry
}

// CheckTrace validates a finite joint trace, including exact original bindings,
// both roles of each physical agent, strong ACK and the serial service guard.
// It never calls Successors, Explore, Initial or the local safety checker.
func CheckTrace(p Profile, events []Event) error {
	switch p {
	case SerialSerial, ServicingSerial, SerialServicing, ServicingServicing:
	default:
		return fmt.Errorf("unknown profile %q", p)
	}
	var l ledger
	for i, e := range events {
		next, err := l.step(p, e)
		if err != nil {
			return fmt.Errorf("event %d (%s): %w", i, e.Kind, err)
		}
		l = next
	}
	return nil
}

func (l ledger) step(p Profile, e Event) (ledger, error) {
	fail := func(reason string) (ledger, error) { return l, fmt.Errorf("%s", reason) }
	if e.Actor > B {
		return fail("unknown physical actor")
	}
	if e.Kind == Crash || e.Kind == Recover {
		if e.Obligation != NoObligation || e.Binding != (Binding{}) || e.Result != NoValue || e.Duplicate {
			return fail("physical interruption must not carry edge data")
		}
		if e.Kind == Crash {
			if l.phase[e.Actor] != 0 {
				return fail("second crash or crash while down")
			}
			l.phase[e.Actor] = 1
			l.entries[e.Actor].received = NoValue
		} else {
			if l.phase[e.Actor] != 1 {
				return fail("recovery while available")
			}
			l.phase[e.Actor] = 2
		}
		return l, nil
	}
	if e.Obligation > BA {
		return fail("unknown obligation")
	}
	// Deliberately repeat the frozen authority table rather than asking model code
	// for its binding: agreement cannot hide a changed endpoint/grant/resource.
	expected := Binding{Caller: A, CallerGeneration: 1, Operation: 1, Callee: B, CalleeGeneration: 1,
		Request: "increment", Grant: "grant:A/g1->B/g1:counter-B", Resource: "counter-B"}
	if e.Obligation == BA {
		expected = Binding{Caller: B, CallerGeneration: 1, Operation: 1, Callee: A, CalleeGeneration: 1,
			Request: "increment", Grant: "grant:B/g1->A/g1:counter-A", Resource: "counter-A"}
	}
	if e.Binding != expected {
		return fail("wrong original correlation, grant or resource")
	}
	responsible := expected.Caller
	switch e.Kind {
	case Admit, Deduplicate, Apply, SendResult, Release:
		responsible = expected.Callee
	case Commit, SendRequest, ReceiveResult, Record, Consume, SendACK:
	default:
		return fail("unknown event")
	}
	if e.Actor != responsible {
		return fail("event attributed to wrong role/agent")
	}
	if l.phase[e.Actor] == 1 {
		return fail("agent work while down")
	}
	if e.Duplicate && e.Kind != Admit && e.Kind != Deduplicate && e.Kind != ReceiveResult && e.Kind != Release {
		return fail("duplicate flag outside delivery")
	}
	usesResult := false
	switch e.Kind {
	case Apply, SendResult, ReceiveResult, Record, Consume:
		usesResult = true
	}
	if usesResult {
		if e.Result != R1 && e.Result != R2 {
			return fail("invalid opaque result")
		}
	} else if e.Result != NoValue {
		return fail("payload on identity-only operation")
	}
	x := &l.entries[e.Obligation]
	switch e.Kind {
	case Commit:
		if x.committed {
			return fail("duplicate commit")
		}
		x.committed = true
	case SendRequest:
		if !x.committed || x.request {
			return fail("request lacks committed source or slot")
		}
		x.request = true
	case Admit, Deduplicate:
		if !x.request || (e.Kind == Admit && x.admitted) || (e.Kind == Deduplicate && !x.admitted) {
			return fail("admission/dedup lacks original request")
		}
		x.admitted = true
		if !e.Duplicate {
			x.request = false
		}
	case Apply:
		if !x.admitted || x.applied {
			return fail("application without admission or repeated")
		}
		own := l.entries[e.Actor]
		mayServe := p == ServicingServicing
		switch e.Actor {
		case A:
			mayServe = mayServe || p == ServicingSerial
		case B:
			mayServe = mayServe || p == SerialServicing
		}
		if own.committed && !own.consumed && !mayServe {
			return fail("serial agent applies incoming while waiting")
		}
		x.applied = true
		x.actual = e.Result
		x.receipt = e.Result
	case SendResult:
		if x.receipt == NoValue || e.Result != x.receipt || x.channel != NoValue {
			return fail("send lacks retained receipt or channel slot")
		}
		x.channel = e.Result
	case ReceiveResult:
		if x.channel == NoValue || e.Result != x.channel {
			return fail("receipt forgery or cross-edge bytes")
		}
		x.received = e.Result
		if !e.Duplicate {
			x.channel = NoValue
		}
	case Record:
		if x.received == NoValue || x.received != e.Result || x.durable != NoValue {
			return fail("record lacks received authentic bytes or repeats")
		}
		x.durable = e.Result
	case Consume:
		if x.durable == NoValue || x.durable != e.Result || x.consumed {
			return fail("consume lacks durable result or repeats")
		}
		x.consumed = true
	case SendACK:
		if x.durable == NoValue || x.ack {
			return fail("ACK lacks durable record or slot")
		}
		x.ack = true
	case Release:
		if !x.ack {
			return fail("release lacks authenticated ACK")
		}
		x.receipt = NoValue
		if !e.Duplicate {
			x.ack = false
		}
	}
	for _, v := range []Value{x.receipt, x.channel, x.received, x.durable} {
		if v != NoValue && (!x.applied || v != x.actual) {
			return fail("inauthentic result")
		}
	}
	return l, nil
}
