// Package mutualobligations models two reciprocal durable-ACK obligations with
// shared physical agents. It is a finite research contract, not a runtime.
package mutualobligations

// Profile fixes which agents may apply incoming work while waiting outgoing.
type Profile string

const (
	SerialSerial       Profile = "serial/serial"
	ServicingSerial    Profile = "servicing/serial"
	SerialServicing    Profile = "serial/servicing"
	ServicingServicing Profile = "servicing/servicing"
)

func Profiles() []Profile {
	return []Profile{SerialSerial, ServicingSerial, SerialServicing, ServicingServicing}
}

func validProfile(p Profile) bool {
	switch p {
	case SerialSerial, ServicingSerial, SerialServicing, ServicingServicing:
		return true
	}
	return false
}
func servicing(p Profile, a Agent) bool {
	return p == ServicingServicing || (a == A && p == ServicingSerial) || (a == B && p == SerialServicing)
}

type Agent uint8

const (
	A Agent = iota
	B
)

type Obligation uint8

const (
	AB Obligation = iota
	BA
	NoObligation
)

type Value uint8

const (
	NoValue Value = iota
	R1
	R2
)

type Kind string

const (
	Commit        Kind = "commit"
	SendRequest   Kind = "send-request"
	Admit         Kind = "admit"
	Deduplicate   Kind = "deduplicate"
	Apply         Kind = "apply"
	SendResult    Kind = "send-result"
	ReceiveResult Kind = "receive-result"
	Record        Kind = "record"
	Consume       Kind = "consume"
	SendACK       Kind = "send-ack"
	Release       Kind = "release"
	Crash         Kind = "crash"
	Recover       Kind = "recover"
)

// Binding is carried on every edge event, including identity-only ACKs. Both
// originals intentionally have operation=1 and equal result alphabets.
type Binding struct {
	Caller           Agent
	CallerGeneration uint8
	Operation        uint8
	Callee           Agent
	CalleeGeneration uint8
	Request          string
	Grant            string
	Resource         string
}

func FixedBinding(d Obligation) Binding {
	switch d {
	case AB:
		return Binding{Caller: A, CallerGeneration: 1, Operation: 1, Callee: B, CalleeGeneration: 1,
			Request: "increment", Grant: "grant:A/g1->B/g1:counter-B", Resource: "counter-B"}
	case BA:
		return Binding{Caller: B, CallerGeneration: 1, Operation: 1, Callee: A, CalleeGeneration: 1,
			Request: "increment", Grant: "grant:B/g1->A/g1:counter-A", Resource: "counter-A"}
	default:
		panic("invalid obligation")
	}
}

type Availability struct {
	Down       bool
	Crashed    bool
	Activation uint8
}

// EdgeState retains one request and three one-slot channels. WorldResult is an
// inaccessible oracle fact, not an agent record or a source of receipt bytes.
type EdgeState struct {
	Committed      bool
	Admitted       bool
	Applied        bool
	Consumed       bool
	Volatile       Value
	Recorded       Value
	Retained       Value
	RequestChannel bool
	ResultChannel  Value
	ACKChannel     bool
	WorldResult    Value
}
type State struct {
	Agents      [2]Availability
	Obligations [2]EdgeState
}
type Event struct {
	Kind       Kind
	Actor      Agent
	Obligation Obligation
	Binding    Binding
	Result     Value
	Duplicate  bool
}
type Edge struct {
	Event Event
	State State
}

func Initial() State        { return State{} }
func Complete(s State) bool { return s.Obligations[AB].Consumed && s.Obligations[BA].Consumed }
func Trap(s State) bool {
	a, b := s.Obligations[AB], s.Obligations[BA]
	return a.Committed && b.Committed && !a.Applied && !b.Applied
}
func Recoverable(s State) State {
	for i := range s.Obligations {
		s.Obligations[i].WorldResult = NoValue
	}
	return s
}

// Successors preserves retry and duplicate-delivery cycles. Only Apply reads
// the executing agent's outgoing waiting facts; all role actions require up.
func Successors(p Profile, s State) []Edge {
	if !validProfile(p) {
		panic("invalid reciprocal profile")
	}
	out := []Edge{}
	add := func(d Obligation, actor Agent, k Kind, v Value, dup bool, n State) {
		out = append(out, Edge{Event: Event{Kind: k, Actor: actor, Obligation: d,
			Binding: FixedBinding(d), Result: v, Duplicate: dup}, State: n})
	}
	for d := AB; d <= BA; d++ {
		e := s.Obligations[d]
		caller, callee := Agent(d), Agent(1-d)
		if !s.Agents[caller].Down {
			if !e.Committed {
				n := s
				n.Obligations[d].Committed = true
				add(d, caller, Commit, NoValue, false, n)
			}
			if e.Committed && !e.RequestChannel {
				n := s
				n.Obligations[d].RequestChannel = true
				add(d, caller, SendRequest, NoValue, false, n)
			}
			if e.ResultChannel != NoValue {
				for _, dup := range []bool{false, true} {
					n := s
					n.Obligations[d].Volatile = e.ResultChannel
					if !dup {
						n.Obligations[d].ResultChannel = NoValue
					}
					add(d, caller, ReceiveResult, e.ResultChannel, dup, n)
				}
			}
			if e.Volatile != NoValue && e.Recorded == NoValue {
				n := s
				n.Obligations[d].Recorded = e.Volatile
				add(d, caller, Record, e.Volatile, false, n)
			}
			if e.Recorded != NoValue && !e.Consumed {
				n := s
				n.Obligations[d].Consumed = true
				add(d, caller, Consume, e.Recorded, false, n)
			}
			if e.Recorded != NoValue && !e.ACKChannel {
				n := s
				n.Obligations[d].ACKChannel = true
				add(d, caller, SendACK, NoValue, false, n)
			}
		}
		if s.Agents[callee].Down {
			continue
		}
		if e.RequestChannel {
			for _, dup := range []bool{false, true} {
				n := s
				n.Obligations[d].Admitted = true
				if !dup {
					n.Obligations[d].RequestChannel = false
				}
				k := Admit
				if e.Admitted {
					k = Deduplicate
				}
				add(d, callee, k, NoValue, dup, n)
			}
		}
		outgoing := s.Obligations[Obligation(callee)]
		waiting := outgoing.Committed && !outgoing.Consumed
		if e.Admitted && !e.Applied && (!waiting || servicing(p, callee)) {
			for _, v := range []Value{R1, R2} {
				n := s
				n.Obligations[d].Applied = true
				n.Obligations[d].Retained = v
				n.Obligations[d].WorldResult = v
				add(d, callee, Apply, v, false, n)
			}
		}
		if e.Retained != NoValue && e.ResultChannel == NoValue {
			n := s
			n.Obligations[d].ResultChannel = e.Retained
			add(d, callee, SendResult, e.Retained, false, n)
		}
		if e.ACKChannel {
			for _, dup := range []bool{false, true} {
				n := s
				n.Obligations[d].Retained = NoValue
				if !dup {
					n.Obligations[d].ACKChannel = false
				}
				add(d, callee, Release, NoValue, dup, n)
			}
		}
	}
	for a := A; a <= B; a++ {
		n := s
		k := Crash
		if s.Agents[a].Down {
			n.Agents[a].Down = false
			n.Agents[a].Activation = 1
			k = Recover
		} else if !s.Agents[a].Crashed {
			n.Agents[a].Down = true
			n.Agents[a].Crashed = true
			n.Obligations[a].Volatile = NoValue
		} else {
			continue
		}
		out = append(out, Edge{Event: Event{Kind: k, Actor: a, Obligation: NoObligation}, State: n})
	}
	return out
}
