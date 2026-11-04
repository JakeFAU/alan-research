// Package authorityservice extends the unchanged reciprocal custody model with
// a finite, explicit authority contract. It is research machinery, not runtime.
package authorityservice

import base "github.com/JakeFAU/alan-research/research/mutualobligations"

// Config selects one of the sixteen frozen experiments.
type Config struct {
	Policy   string
	Scenario string
}

func Configurations() []Config {
	out := []Config{}
	for _, p := range []string{"S1", "S2", "S3"} {
		for _, s := range []string{"Local", "Explicit", "Missing", "ScopeMix"} {
			out = append(out, Config{Policy: p, Scenario: s})
		}
	}
	for _, s := range []string{"Shared", "LateGrant", "Revocable"} {
		out = append(out, Config{Policy: "S2", Scenario: s})
	}
	return append(out, Config{Policy: "Cartesian", Scenario: "ScopeMix"})
}
func (c Config) Name() string { return c.Policy + "/" + c.Scenario }
func valid(c Config) bool {
	for _, v := range Configurations() {
		if v == c {
			return true
		}
	}
	return false
}

type Permission struct {
	Operation string
	Resource  string
	Scope     string
}
type Tuple uint8

const (
	NoTuple Tuple = iota
	APermission
	BPermission
	SharedPermission
	InventedPermission
)

func (t Tuple) Permission() Permission {
	switch t {
	case APermission:
		return Permission{Operation: "increment", Resource: "counter-A", Scope: "red"}
	case BPermission:
		return Permission{Operation: "increment", Resource: "counter-B", Scope: "blue"}
	case SharedPermission:
		return Permission{Operation: "read", Resource: "status", Scope: "shared"}
	case InventedPermission:
		return Permission{Operation: "increment", Resource: "counter-A", Scope: "blue"}
	default:
		return Permission{}
	}
}

type PermissionSet uint8

func Set(ts ...Tuple) PermissionSet {
	var s PermissionSet
	for _, t := range ts {
		s |= 1 << t
	}
	return s
}
func (s PermissionSet) Contains(t Tuple) bool { return s&Set(t) != 0 }

// ServiceBinding explicitly maps original base custody identities to a complete
// authority tuple; the inherited Binding.Grant is never authority evidence.
type ServiceBinding struct {
	Original         string
	Caller           base.Agent
	CallerGeneration uint8
	Callee           base.Agent
	CalleeGeneration uint8
	Operation        uint8
	Permission       Permission
}

func Binding(c Config, d base.Obligation) ServiceBinding {
	t := BPermission
	if d == base.BA {
		switch c.Scenario {
		case "Local":
			t = APermission
		case "ScopeMix":
			t = InventedPermission
		}
	}
	if c.Scenario == "Shared" {
		t = SharedPermission
	}
	name := "AB"
	if d == base.BA {
		name = "BA"
	}
	return ServiceBinding{
		Original: name, Caller: base.Agent(d), CallerGeneration: 1,
		Callee: base.Agent(1 - d), CalleeGeneration: 1,
		Operation: 1, Permission: t.Permission(),
	}
}
func tuple(p Permission) Tuple {
	for _, t := range []Tuple{APermission, BPermission, SharedPermission, InventedPermission} {
		if p == t.Permission() {
			return t
		}
	}
	return NoTuple
}

// GrantRecord is a retained claim, including rejected issuance. Accepted is a
// fact of the issuance event; current validity is checked separately.
type GrantRecord struct {
	ID                string
	Parent            string
	Root              string
	Issuer            base.Agent
	IssuerGeneration  uint8
	Grantee           base.Agent
	GranteeGeneration uint8
	Permissions       PermissionSet
	Binding           ServiceBinding
	Revision          uint8
	NoRedelegation    bool
	Accepted          bool
}
type Path uint8

const (
	NoPath Path = iota
	OwnRoot
	Delegated
	Ambient
	Cartesian
)

type AuthorityWitness struct {
	Path       Path
	Root       string
	Grant      string
	Principal  base.Agent
	Generation uint8
	Permission Permission
	Revision   uint8
}
type Decision uint8

const (
	Undecided Decision = iota
	Allowed
	Refused
)

type GrantStatus uint8

const (
	NoGrant GrantStatus = iota
	RetainedGrant
	RejectedGrant
)

type Violation uint8

const (
	ConsentViolation Violation = 1 << iota
	TupleInvention
)

type AuthorityState struct {
	Decision         Decision
	Path             Path
	Refusal          uint8 // 1: no covering authority; 2: selected path revoked
	Permit           bool
	PermitActivation uint8
	AdmissionRefused bool
}
type State struct {
	Base      base.State
	Authority [2]AuthorityState
	Grant     GrantStatus
	Revoked   bool
	Unsafe    Violation // monotone; completion cannot erase an unsafe prefix
}

// Event is comparable and contains all descriptors needed by a separate ledger.
// Grant carries issuance on commit (inline) or retain-grant (late). Authorization
// and dispatch/application events retain the selected original authority path.
type Event struct {
	Kind       string
	Actor      base.Agent
	Obligation base.Obligation
	Binding    ServiceBinding
	Result     base.Value
	Duplicate  bool
	Activation uint8
	Revision   uint8
	HasGrant   bool
	Grant      GrantRecord
	Authority  AuthorityWitness
	Decision   Decision
	Refusal    string
	Violations Violation
}
type Transition struct {
	Event Event
	State State
}

func Initial(c Config) State {
	if !valid(c) {
		panic("invalid authority configuration")
	}
	return State{Base: base.Initial()}
}
func Complete(s State) bool { return base.Complete(s.Base) }
func root(a base.Agent, t Tuple) string {
	if a == base.A {
		if t == APermission {
			return "A/g1:a"
		}
		if t == SharedPermission {
			return "A/g1:r"
		}
	}
	if a == base.B {
		if t == BPermission {
			return "B/g1:b"
		}
		if t == SharedPermission {
			return "B/g1:r"
		}
	}
	return ""
}
func current(s State, r string) bool { return r != "" && !(s.Revoked && r == "B/g1:b") }
func grant(c Config, s State) GrantRecord {
	return GrantRecord{ID: "grant:BA:b", Parent: "B/g1:b", Root: "B/g1:b", Issuer: base.B, IssuerGeneration: 1,
		Grantee: base.A, GranteeGeneration: 1, Permissions: Set(BPermission), Binding: Binding(c, base.BA),
		Revision: 0, NoRedelegation: true, Accepted: !s.Revoked}
}
func witness(c Config, s State, d base.Obligation, p Path) AuthorityWitness {
	b := Binding(c, d)
	t := tuple(b.Permission)
	w := AuthorityWitness{Path: p, Principal: b.Callee, Generation: b.CalleeGeneration, Permission: b.Permission}
	switch p {
	case OwnRoot:
		w.Root = root(b.Callee, t)
	case Delegated:
		w.Root = "B/g1:b"
		w.Grant = "grant:BA:b"
	case Ambient:
		w.Root = root(b.Caller, t)
	case Cartesian:
		w.Root = "A/g1:a + B/g1:b"
	}
	return w
}
func pathCurrent(c Config, s State, d base.Obligation, p Path) bool {
	if p == Cartesian {
		return !s.Revoked
	}
	return current(s, witness(c, s, d, p).Root)
}
func choose(c Config, s State, d base.Obligation) Path {
	b := Binding(c, d)
	t := tuple(b.Permission)
	if current(s, root(b.Callee, t)) {
		return OwnRoot
	}
	if c.Policy == "S1" {
		return NoPath
	}
	if d == base.BA && s.Grant == RetainedGrant && t == BPermission && !s.Revoked {
		return Delegated
	}
	if c.Policy == "S3" && current(s, root(b.Caller, t)) {
		return Ambient
	}
	if c.Policy == "Cartesian" && d == base.BA && t == InventedPermission && s.Grant == RetainedGrant && !s.Revoked {
		return Cartesian
	}
	return NoPath
}
func violation(p Path) Violation {
	switch p {
	case Ambient:
		return ConsentViolation
	case Cartesian:
		return TupleInvention
	}
	return 0
}
func describe(c Config, s State, e Event) Event {
	e.Activation = s.Base.Agents[e.Actor].Activation
	if s.Revoked {
		e.Revision = 1
	}
	if e.Obligation != base.NoObligation {
		e.Binding = Binding(c, e.Obligation)
		a := s.Authority[e.Obligation]
		e.Decision = a.Decision
		if e.Kind == "authorize" {
			p := choose(c, s, e.Obligation)
			e.Decision = Allowed
			if p == NoPath {
				e.Decision = Refused
				e.Refusal = "no-covering-authority"
			}
			e.Authority = witness(c, s, e.Obligation, p)
		} else if e.Kind == "dispatch-admit" || e.Kind == "dispatch-refuse" || e.Kind == "apply" {
			e.Authority = witness(c, s, e.Obligation, a.Path)
		}
		if e.Kind == "dispatch-refuse" {
			e.Refusal = "selected-path-revoked"
		}
		if e.Kind == "authorize" || e.Kind == "dispatch-admit" || e.Kind == "apply" {
			e.Violations = violation(e.Authority.Path)
		}
		if (e.Kind == "commit" && e.Obligation == base.BA && inline(c)) || e.Kind == "retain-grant" {
			e.HasGrant = true
			e.Grant = grant(c, s)
		}
	}
	return e
}
func inline(c Config) bool {
	return c.Scenario == "Explicit" || c.Scenario == "ScopeMix" || c.Scenario == "Revocable"
}

// Successors keeps every inherited transition except ungated Apply. Authority
// steps stutter under projection. A physical crash clears that server's permits.
func Successors(c Config, s State) []Transition {
	if !valid(c) {
		panic("invalid authority configuration")
	}
	out := []Transition{}
	add := func(e Event, n State) {
		e = describe(c, s, e)
		n.Unsafe |= e.Violations
		out = append(out, Transition{Event: e, State: n})
	}
	for _, edge := range base.Successors(base.ServicingSerial, s.Base) {
		e := edge.Event
		n := s
		n.Base = edge.State
		if e.Kind == base.Apply {
			if !s.Authority[e.Obligation].Permit {
				continue
			}
			n.Authority[e.Obligation].Permit = false
			n.Authority[e.Obligation].PermitActivation = 0
		}
		if e.Kind == base.Crash {
			d := base.Obligation(1 - e.Actor)
			n.Authority[d].Permit = false
			n.Authority[d].PermitActivation = 0
		}
		if e.Kind == base.Commit && e.Obligation == base.BA && inline(c) {
			n.Grant = RetainedGrant
			if s.Revoked {
				n.Grant = RejectedGrant
			}
		}
		add(Event{
			Kind: string(e.Kind), Actor: e.Actor, Obligation: e.Obligation,
			Result: e.Result, Duplicate: e.Duplicate,
		}, n)
	}
	for d := base.AB; d <= base.BA; d++ {
		server := base.Agent(1 - d)
		a := s.Authority[d]
		if s.Base.Agents[server].Down || !s.Base.Obligations[d].Admitted {
			continue
		}
		if a.Decision == Undecided {
			n := s
			p := choose(c, s, d)
			n.Authority[d].Path = p
			n.Authority[d].Decision = Allowed
			if p == NoPath {
				n.Authority[d].Decision = Refused
				n.Authority[d].Refusal = 1
			}
			add(Event{Kind: "authorize", Actor: server, Obligation: d}, n)
		}
		if a.Decision == Allowed && !a.Permit && !s.Base.Obligations[d].Applied {
			if pathCurrent(c, s, d, a.Path) {
				n := s
				n.Authority[d].Permit = true
				n.Authority[d].PermitActivation = s.Base.Agents[server].Activation
				add(Event{Kind: "dispatch-admit", Actor: server, Obligation: d}, n)
			} else if !a.AdmissionRefused {
				n := s
				n.Authority[d].AdmissionRefused = true
				add(Event{Kind: "dispatch-refuse", Actor: server, Obligation: d}, n)
			}
		}
	}
	if c.Scenario == "LateGrant" && s.Base.Obligations[base.BA].Committed &&
		s.Grant == NoGrant && !s.Base.Agents[base.B].Down {
		n := s
		n.Grant = RetainedGrant
		add(Event{Kind: "retain-grant", Actor: base.B, Obligation: base.BA}, n)
	}
	if c.Scenario == "Revocable" && !s.Revoked && !s.Base.Agents[base.B].Down {
		n := s
		n.Revoked = true
		add(Event{Kind: "revoke", Actor: base.B, Obligation: base.NoObligation}, n)
	}
	return out
}
