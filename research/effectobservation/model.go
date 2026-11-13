// Package effectobservation models one admitted effect across a shared Store owner loss.
package effectobservation

type OwnerStage string

const (
	OwnerOriginal OwnerStage = "original"
	OwnerStopped  OwnerStage = "stopped"
	OwnerReopened OwnerStage = "reopened"
)

type EventKind string

const (
	Launch      EventKind = "launch"
	Ingress     EventKind = "ingress"
	Apply       EventKind = "apply"
	Response    EventKind = "response"
	Record      EventKind = "record"
	Consume     EventKind = "consume"
	RevokeB     EventKind = "revoke-b"
	LoseOwner   EventKind = "lose-owner"
	Reopen      EventKind = "reopen"
	FreshRefuse EventKind = "fresh-refuse"
)

type Result string

const (
	NoResult Result = ""
	Value1   Result = "v1"
	Value2   Result = "v2"
)

type Root string

const (
	RootB Root = "B/g1:b"
	RootC Root = "B/g1:c"
)

type Actor string

const (
	ActorA        Actor = "A"
	ActorB        Actor = "B"
	ActorEndpoint Actor = "endpoint"
	ActorOwner    Actor = "owner"
)

// Binding names the one original request, its admission path and its wire.
type Binding struct {
	Original           string `json:"original"`
	SourceGeneration   uint8  `json:"source_generation"`
	SourceRevision     uint8  `json:"source_revision"`
	Grant              string `json:"grant"`
	Root               Root   `json:"root"`
	OriginalActivation uint8  `json:"original_activation"`
	Wire               string `json:"wire"`
}

// State distinguishes historical admission, current authority, volatile owner
// custody, external activity, and locally retained results.
type State struct {
	Binding          Binding    `json:"binding"`
	AcceptedGrant    bool       `json:"accepted_grant"`
	FirstAllow       bool       `json:"first_allow"`
	Debit            uint8      `json:"debit"`
	FirstGateAllowed bool       `json:"first_gate_allowed"`
	Owner            OwnerStage `json:"owner"`
	Activation       uint8      `json:"activation"`
	BRevoked         bool       `json:"b_revoked"`
	ControlRoot      Root       `json:"control_root"`
	FreshRefused     bool       `json:"fresh_refused"`
	Permit           bool       `json:"permit"`
	Launched         bool       `json:"launched"`
	Ingress          bool       `json:"ingress"`
	Applied          Result     `json:"applied"`
	Response         Result     `json:"response"`
	Retained         Result     `json:"retained"`
	Consumed         Result     `json:"consumed"`
}

// Event binds each transition to the same original, with a separate target
// root for authority actions and a result value for result-bearing actions.
type Event struct {
	Kind        EventKind `json:"kind"`
	Binding     Binding   `json:"binding"`
	Actor       Actor     `json:"actor"`
	Incarnation uint8     `json:"incarnation"`
	Root        Root      `json:"root"`
	Value       Result    `json:"value"`
}

type Edge struct {
	Event Event `json:"event"`
	State State `json:"state"`
}

type StoreView struct {
	Binding          Binding    `json:"binding"`
	AcceptedGrant    bool       `json:"accepted_grant"`
	FirstAllow       bool       `json:"first_allow"`
	Debit            uint8      `json:"debit"`
	FirstGateAllowed bool       `json:"first_gate_allowed"`
	Owner            OwnerStage `json:"owner"`
	Activation       uint8      `json:"activation"`
	BRevoked         bool       `json:"b_revoked"`
	ControlRoot      Root       `json:"control_root"`
	FreshRefused     bool       `json:"fresh_refused"`
	Retained         Result     `json:"retained"`
	Consumed         Result     `json:"consumed"`
}

type DossierView struct {
	StoreView
	Ingress bool   `json:"ingress"`
	Applied Result `json:"applied"`
}

// Initial is the cut after the first successful gate and before launch.
func Initial() State {
	return State{
		Binding: Binding{
			Original: "BA/1", SourceGeneration: 1, SourceRevision: 1,
			Grant: "grant:BA:b", Root: RootB, OriginalActivation: 1, Wire: "wire:BA:1",
		},
		AcceptedGrant: true, FirstAllow: true, Debit: 1, FirstGateAllowed: true,
		Owner: OwnerOriginal, Activation: 1, ControlRoot: RootC, Permit: true,
	}
}

// Successors returns every enabled transition in a fixed order. Its domain is
// the post-admission states reachable from Initial.
func Successors(s State) []Edge {
	edges := []Edge{}
	add := func(kind EventKind, actor Actor, incarnation uint8, root Root, value Result, target State) {
		edges = append(edges, Edge{
			Event: Event{Kind: kind, Binding: s.Binding, Actor: actor, Incarnation: incarnation, Root: root, Value: value},
			State: target,
		})
	}

	original := s.Owner == OwnerOriginal
	live := s.Owner != OwnerStopped

	if original && s.Permit {
		next := s
		next.Permit = false
		next.Launched = true
		add(Launch, ActorA, s.Activation, "", NoResult, next)
	}
	if s.Launched && !s.Ingress {
		next := s
		next.Ingress = true
		add(Ingress, ActorEndpoint, 0, "", NoResult, next)
	}
	if s.Ingress && s.Applied == NoResult {
		for _, value := range []Result{Value1, Value2} {
			next := s
			next.Applied = value
			add(Apply, ActorEndpoint, 0, "", value, next)
		}
	}
	if original && s.Applied != NoResult && s.Response == NoResult && s.Retained == NoResult {
		next := s
		next.Response = s.Applied
		add(Response, ActorA, s.Activation, "", s.Applied, next)
	}
	if original && s.Response != NoResult {
		next := s
		next.Response = NoResult
		next.Retained = s.Response
		add(Record, ActorA, s.Activation, "", s.Response, next)
	}
	if live && s.Retained != NoResult && s.Consumed == NoResult {
		next := s
		next.Consumed = s.Retained
		add(Consume, ActorB, s.Activation, "", s.Retained, next)
	}
	if live && !s.BRevoked {
		next := s
		next.BRevoked = true
		add(RevokeB, ActorB, s.Activation, RootB, NoResult, next)
	}
	if original {
		next := s
		next.Owner = OwnerStopped
		next.Permit = false
		next.Response = NoResult
		add(LoseOwner, ActorOwner, s.Activation, "", NoResult, next)
	}
	if s.Owner == OwnerStopped {
		next := s
		next.Owner = OwnerReopened
		next.Activation++
		add(Reopen, ActorOwner, next.Activation, "", NoResult, next)
	}
	if s.Owner == OwnerReopened && s.BRevoked && s.Retained == NoResult && !s.FreshRefused {
		next := s
		next.FreshRefused = true
		add(FreshRefuse, ActorA, s.Activation, RootB, NoResult, next)
	}
	return edges
}

// ObserveStore projects retained local facts and the protocol-supplied owner stage.
func ObserveStore(s State) StoreView {
	return StoreView{
		Binding: s.Binding, AcceptedGrant: s.AcceptedGrant,
		FirstAllow: s.FirstAllow, Debit: s.Debit, FirstGateAllowed: s.FirstGateAllowed,
		Owner: s.Owner, Activation: s.Activation, BRevoked: s.BRevoked, ControlRoot: s.ControlRoot,
		FreshRefused: s.FreshRefused, Retained: s.Retained, Consumed: s.Consumed,
	}
}

// ObserveDossier adds the modeled truth at a cut with complete endpoint evidence.
func ObserveDossier(s State) DossierView {
	return DossierView{StoreView: ObserveStore(s), Ingress: s.Ingress, Applied: s.Applied}
}
