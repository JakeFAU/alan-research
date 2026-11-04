// Package twoagent models a fixed, abstract two-agent result-custody contract.
package twoagent

type Profile string

const (
	DeliveryACK Profile = "delivery-ack"
	DurableACK  Profile = "durable-ack"
)

type Value string

const (
	R1 Value = "r1"
	R2 Value = "r2"
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
const FixedGrant = "grant:A/g1->B/g1:increment-counter"

type Delegation struct {
	Caller           string
	CallerGeneration uint8
	Operation        uint8
	Callee           string
	CalleeGeneration uint8
	Request          string
}

func FixedDelegation() Delegation {
	return Delegation{
		Caller: "A", CallerGeneration: 1, Operation: 1,
		Callee: "B", CalleeGeneration: 1, Request: "increment-counter",
	}
}

type State struct {
	Delegation     Delegation
	Grant          string
	Committed      bool
	ADown          bool
	Crashed        bool
	Activation     uint8
	AVolatile      Value
	ARecorded      Value
	AConsumed      bool
	BAdmitted      bool
	Applied        bool
	BResult        Value
	RequestChannel bool
	ResultChannel  Value
	ACKChannel     bool
	WorldResult    Value
}

// Event carries a fixed delegation/grant binding. Duplicate retains a transport
// copy on delivery; it is legal only for admission, receipt and ACK delivery.
// ACKs carry identity only, never result bytes or a result digest.
type Event struct {
	Kind       Kind
	Delegation Delegation
	Grant      string
	Result     Value
	Duplicate  bool
}
type Edge struct {
	Event Event
	State State
}

func Initial() State            { return State{Delegation: FixedDelegation(), Grant: FixedGrant} }
func Recoverable(s State) State { s.WorldResult = ""; return s }

// Successors is the executable contract. WorldResult is assigned only at Apply
// and never read by a guard or used as a source of agent-visible result bytes.
func Successors(p Profile, s State) []Edge {
	if p != DeliveryACK && p != DurableACK {
		panic("unknown custody profile")
	}
	var out []Edge
	add := func(kind Kind, result Value, duplicate bool, next State) {
		out = append(out, Edge{
			Event: Event{Kind: kind, Delegation: s.Delegation, Grant: s.Grant, Result: result, Duplicate: duplicate},
			State: next,
		})
	}
	if !s.ADown {
		if !s.Committed {
			n := s
			n.Committed = true
			add(Commit, "", false, n)
		}
		if s.Committed && !s.RequestChannel {
			n := s
			n.RequestChannel = true
			add(SendRequest, "", false, n)
		}
		if s.ResultChannel != "" {
			for _, dup := range []bool{false, true} {
				n := s
				n.AVolatile = s.ResultChannel
				if !dup {
					n.ResultChannel = ""
				}
				add(ReceiveResult, s.ResultChannel, dup, n)
			}
		}
		if s.AVolatile != "" && s.ARecorded == "" {
			n := s
			n.ARecorded = s.AVolatile
			add(Record, s.AVolatile, false, n)
		}
		if s.ARecorded != "" && !s.AConsumed {
			n := s
			n.AConsumed = true
			add(Consume, s.ARecorded, false, n)
		}
		ready := s.ARecorded != ""
		if p == DeliveryACK {
			ready = ready || s.AVolatile != ""
		}
		if ready && !s.ACKChannel {
			n := s
			n.ACKChannel = true
			add(SendACK, "", false, n)
		}
	} else {
		n := s
		n.ADown = false
		n.Activation = 1
		add(Recover, "", false, n)
	}
	// B and the channels can progress independently of A's availability and
	// without inspecting A's private volatile or durable records.
	if s.RequestChannel {
		for _, dup := range []bool{false, true} {
			n := s
			n.BAdmitted = true
			if !dup {
				n.RequestChannel = false
			}
			kind := Admit
			if s.BAdmitted {
				kind = Deduplicate
			}
			add(kind, "", dup, n)
		}
	}
	if s.BAdmitted && !s.Applied {
		for _, v := range []Value{R1, R2} {
			n := s
			n.Applied = true
			n.BResult = v
			n.WorldResult = v
			add(Apply, v, false, n)
		}
	}
	if s.BResult != "" && s.ResultChannel == "" {
		n := s
		n.ResultChannel = s.BResult
		add(SendResult, s.BResult, false, n)
	}
	if s.ACKChannel {
		for _, dup := range []bool{false, true} {
			n := s
			n.BResult = ""
			if !dup {
				n.ACKChannel = false
			}
			add(Release, "", dup, n)
		}
	}
	// Prefer ordinary edge work before the interruption when shortest witnesses
	// have equal length; every crash position remains in the reachable graph.
	if !s.ADown && !s.Crashed {
		n := s
		n.ADown = true
		n.Crashed = true
		n.AVolatile = ""
		add(Crash, "", false, n)
	}
	return out
}
