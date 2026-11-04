package authorityservice

import (
	"fmt"
	"maps"

	base "github.com/JakeFAU/alan-research/research/mutualobligations"
)

// Audit separates acceptance by the candidate protocol from the frozen
// explicit-consent safety policy. Paid means actual authenticated consumption.
// AuthoritySafe is certified only for a protocol-valid prefix without policy
// violations; malformed full histories receive no safety certification.
type Audit struct {
	ProtocolValid bool
	AuthoritySafe bool
	Conditions    []string
	Paid          [2]bool
	Unpaid        [2]bool
	EventsChecked int
}

type ledgerOriginal struct {
	committed, admitted, applied, consumed           bool
	request, ack                                     bool
	result, retained, published, delivered, recorded base.Value
	decision                                         Decision
	witness                                          AuthorityWitness
	permit, admissionRefused                         bool
	permitEpoch                                      uint8
}
type ledgerAgent struct {
	down, interrupted bool
	epoch             uint8
}
type authorityLedger struct {
	originals  [2]ledgerOriginal
	agents     [2]ledgerAgent
	grants     map[string]GrantRecord // immutable after insertion; branches copy on write
	revoked    bool
	violations Violation
}

// AuditTrace interprets retained events independently. It does not call any
// operational successor, authorization, binding, or transition predicate.
// Protocol errors stop at the first invalid event; unsafe candidate events are
// retained and do not stop custody checking.
func AuditTrace(c Config, events []Event) (Audit, error) {
	l := authorityLedger{grants: map[string]GrantRecord{}}
	if !ledgerConfig(c) {
		a := l.audit(0)
		a.ProtocolValid = false
		a.AuthoritySafe = false
		a.Conditions = append(a.Conditions, "unknown-frozen-configuration")
		return a, fmt.Errorf("unknown frozen configuration")
	}
	for i, e := range events {
		if err := l.accept(c, e); err != nil {
			a := l.audit(i)
			a.ProtocolValid = false
			a.AuthoritySafe = false
			a.Conditions = append(a.Conditions, err.Error())
			return a, fmt.Errorf("event %d (%s): %w", i, e.Kind, err)
		}
	}
	return l.audit(len(events)), nil
}
func (l authorityLedger) audit(n int) Audit {
	a := Audit{ProtocolValid: true, AuthoritySafe: l.violations == 0, Conditions: []string{}, EventsChecked: n}
	if l.violations&ConsentViolation != 0 {
		a.Conditions = append(a.Conditions, "ambient-consent")
	}
	if l.violations&TupleInvention != 0 {
		a.Conditions = append(a.Conditions, "tuple-invention")
	}
	for d := range a.Paid {
		a.Paid[d] = l.originals[d].consumed
		a.Unpaid[d] = !a.Paid[d]
	}
	return a
}
func ledgerConfig(c Config) bool {
	switch c.Policy {
	case "S1", "S2", "S3":
		switch c.Scenario {
		case "Local", "Explicit", "Missing", "ScopeMix":
			return true
		}
		return c.Policy == "S2" && (c.Scenario == "Shared" || c.Scenario == "LateGrant" || c.Scenario == "Revocable")
	case "Cartesian":
		return c.Scenario == "ScopeMix"
	}
	return false
}
func ledgerBinding(c Config, d base.Obligation) ServiceBinding {
	p := Permission{Operation: "increment", Resource: "counter-B", Scope: "blue"}
	if d == base.BA {
		if c.Scenario == "Local" {
			p = Permission{Operation: "increment", Resource: "counter-A", Scope: "red"}
		}
		if c.Scenario == "ScopeMix" {
			p = Permission{Operation: "increment", Resource: "counter-A", Scope: "blue"}
		}
	}
	if c.Scenario == "Shared" {
		p = Permission{Operation: "read", Resource: "status", Scope: "shared"}
	}
	name := "AB"
	if d == base.BA {
		name = "BA"
	}
	return ServiceBinding{Original: name, Caller: base.Agent(d), CallerGeneration: 1, Callee: base.Agent(1 - d), CalleeGeneration: 1, Operation: 1, Permission: p}
}

type ledgerRoot struct {
	holder     base.Agent
	permission Permission
	bits       PermissionSet
}

func ledgerRoots() map[string]ledgerRoot {
	return map[string]ledgerRoot{
		"A/g1:a": {holder: base.A, permission: Permission{Operation: "increment", Resource: "counter-A", Scope: "red"}, bits: 1 << APermission},
		"B/g1:b": {holder: base.B, permission: Permission{Operation: "increment", Resource: "counter-B", Scope: "blue"}, bits: 1 << BPermission},
		"A/g1:r": {holder: base.A, permission: Permission{Operation: "read", Resource: "status", Scope: "shared"}, bits: 1 << SharedPermission},
		"B/g1:r": {holder: base.B, permission: Permission{Operation: "read", Resource: "status", Scope: "shared"}, bits: 1 << SharedPermission},
	}
}
func ledgerTuple(p Permission) PermissionSet {
	switch p {
	case (Permission{Operation: "increment", Resource: "counter-A", Scope: "red"}):
		return 1 << APermission
	case (Permission{Operation: "increment", Resource: "counter-B", Scope: "blue"}):
		return 1 << BPermission
	case (Permission{Operation: "read", Resource: "status", Scope: "shared"}):
		return 1 << SharedPermission
	case (Permission{Operation: "increment", Resource: "counter-A", Scope: "blue"}):
		return 1 << InventedPermission
	}
	return 0
}
func (l authorityLedger) rootLive(id string) bool {
	_, ok := ledgerRoots()[id]
	return ok && !(id == "B/g1:b" && l.revoked)
}

// grantPath walks the retained graph, including hostile rootless cycles, rather
// than trusting a root label. Every hop keeps the entire original binding.
func (l authorityLedger) grantPath(id string, binding ServiceBinding, current bool) (string, error) {
	visited := map[string]bool{}
	chain := []GrantRecord{}
	for {
		if visited[id] {
			return "", fmt.Errorf("grant-cycle")
		}
		visited[id] = true
		if r, ok := ledgerRoots()[id]; ok {
			if _, collision := l.grants[id]; collision {
				return "", fmt.Errorf("root-identity-collision")
			}
			if current && !l.rootLive(id) {
				return "", fmt.Errorf("revoked-origin")
			}
			holder, bits := r.holder, r.bits
			for i := len(chain) - 1; i >= 0; i-- {
				g := chain[i]
				if g.Root != id {
					return "", fmt.Errorf("mutated-root")
				}
				if g.Issuer != holder || g.IssuerGeneration != 1 || g.GranteeGeneration != 1 {
					return "", fmt.Errorf("identity-adjacency")
				}
				if g.Permissions == 0 || g.Permissions & ^bits != 0 {
					return "", fmt.Errorf("tuple-attenuation")
				}
				if g.Binding != binding || g.Issuer != binding.Caller || g.Grantee != binding.Callee {
					return "", fmt.Errorf("service-binding")
				}
				if g.Revision != 0 {
					return "", fmt.Errorf("grant-revision")
				}
				if !g.NoRedelegation {
					return "", fmt.Errorf("redelegation-disposition")
				}
				if i > 0 {
					return "", fmt.Errorf("no-redelegation")
				}
				holder, bits = g.Grantee, g.Permissions
			}
			return id, nil
		}
		g, ok := l.grants[id]
		if !ok {
			return "", fmt.Errorf("missing-origin")
		}
		if !g.Accepted {
			return "", fmt.Errorf("rejected-issuance")
		}
		if g.ID != id {
			return "", fmt.Errorf("mutated-grant-identity")
		}
		chain = append(chain, g)
		id = g.Parent
	}
}
func (l *authorityLedger) issue(e Event) error {
	g := e.Grant
	if g.ID == "" {
		return fmt.Errorf("empty-grant-identity")
	}
	if _, ok := l.grants[g.ID]; ok {
		return fmt.Errorf("immutable-grant")
	}
	next := maps.Clone(l.grants)
	next[g.ID] = g
	// Evaluate the claim as if accepted; the recorded acceptance must match the
	// independently established ancestry and issuance-time validity.
	probe := *l
	probe.grants = maps.Clone(next)
	claim := g
	claim.Accepted = true
	probe.grants[g.ID] = claim
	_, err := probe.grantPath(g.ID, e.Binding, true)
	if g.Accepted && err != nil {
		return err
	}
	if !g.Accepted && err == nil {
		return fmt.Errorf("unjustified-issuance-rejection")
	}
	l.grants = next
	return nil
}
func (l authorityLedger) selectWitness(c Config, b ServiceBinding) AuthorityWitness {
	w := AuthorityWitness{Principal: b.Callee, Generation: b.CalleeGeneration, Permission: b.Permission}
	for id, r := range ledgerRoots() {
		if r.holder == b.Callee && r.permission == b.Permission && l.rootLive(id) {
			w.Path = OwnRoot
			w.Root = id
			return w
		}
	}
	if c.Policy == "S1" {
		return w
	}
	for id, g := range l.grants {
		origin, err := l.grantPath(id, b, true)
		if err == nil && g.Permissions&ledgerTuple(b.Permission) != 0 {
			w.Path = Delegated
			w.Root = origin
			w.Grant = id
			return w
		}
	}
	if c.Policy == "S3" {
		for id, r := range ledgerRoots() {
			if r.holder == b.Caller && r.permission == b.Permission && l.rootLive(id) {
				w.Path = Ambient
				w.Root = id
				return w
			}
		}
	}
	// This named mutant is intentionally candidate-valid but authority-unsafe.
	if c.Policy == "Cartesian" && b.Original == "BA" && ledgerTuple(b.Permission) == 1<<InventedPermission {
		g, ok := l.grants["grant:BA:b"]
		_, err := l.grantPath("grant:BA:b", b, true)
		if ok && err == nil && g.Permissions&(1<<BPermission) != 0 && l.rootLive("A/g1:a") {
			w.Path = Cartesian
			w.Root = "A/g1:a + B/g1:b"
		}
	}
	return w
}
func (l authorityLedger) witnessLive(w AuthorityWitness, b ServiceBinding) bool {
	if w.Path == Cartesian {
		return l.rootLive("A/g1:a") && l.rootLive("B/g1:b")
	}
	if w.Path == Delegated {
		r, err := l.grantPath(w.Grant, b, true)
		return err == nil && r == w.Root
	}
	return l.rootLive(w.Root)
}
func (l *authorityLedger) accept(c Config, e Event) error {
	if e.Actor > base.B {
		return fmt.Errorf("actor-identity")
	}
	a := &l.agents[e.Actor]
	if e.Activation != a.epoch {
		return fmt.Errorf("activation-binding")
	}
	revision := uint8(0)
	if l.revoked {
		revision = 1
	}
	if e.Revision != revision {
		return fmt.Errorf("authority-revision-order")
	}
	if e.Kind == "recover" {
		if e.Obligation != base.NoObligation || !a.down {
			return fmt.Errorf("recovery-order")
		}
		if err := ledgerEmpty(e); err != nil {
			return err
		}
		a.down = false
		a.epoch++
		return nil
	}
	if a.down {
		return fmt.Errorf("actor-down")
	}
	if e.Kind == "crash" {
		if e.Obligation != base.NoObligation || a.interrupted {
			return fmt.Errorf("crash-order")
		}
		if err := ledgerEmpty(e); err != nil {
			return err
		}
		a.down = true
		a.interrupted = true
		l.originals[e.Actor].delivered = base.NoValue
		l.originals[1-e.Actor].permit = false
		l.originals[1-e.Actor].permitEpoch = 0
		return nil
	}
	if e.Kind == "revoke" {
		if c.Scenario != "Revocable" || e.Actor != base.B || l.revoked || e.Obligation != base.NoObligation {
			return fmt.Errorf("revocation-order")
		}
		if err := ledgerEmpty(e); err != nil {
			return err
		}
		l.revoked = true
		return nil
	}
	if e.Obligation > base.BA {
		return fmt.Errorf("original-identity")
	}
	b := ledgerBinding(c, e.Obligation)
	if e.Binding != b {
		return fmt.Errorf("original-binding")
	}
	o := &l.originals[e.Obligation]
	callerAction := e.Kind == "commit" || e.Kind == "send-request" || e.Kind == "receive-result" || e.Kind == "record" || e.Kind == "consume" || e.Kind == "send-ack" || e.Kind == "retain-grant"
	actor := b.Callee
	if callerAction {
		actor = b.Caller
	}
	if e.Actor != actor {
		return fmt.Errorf("role-binding")
	}
	if e.Kind != "authorize" && e.Decision != o.decision {
		return fmt.Errorf("immutable-decision")
	}
	grantExpected := e.Kind == "retain-grant" || e.Kind == "commit" && e.Obligation == base.BA && (c.Scenario == "Explicit" || c.Scenario == "ScopeMix" || c.Scenario == "Revocable")
	if e.HasGrant != grantExpected || !e.HasGrant && e.Grant != (GrantRecord{}) {
		return fmt.Errorf("grant-retention-boundary")
	}
	authorityEvent := e.Kind == "authorize" || e.Kind == "dispatch-admit" || e.Kind == "dispatch-refuse" || e.Kind == "apply"
	if !authorityEvent && e.Authority != (AuthorityWitness{}) {
		return fmt.Errorf("result-is-not-authority")
	}
	expectedViolation := Violation(0)
	if authorityEvent && e.Kind != "dispatch-refuse" {
		if e.Authority.Path == Ambient {
			expectedViolation = ConsentViolation
		}
		if e.Authority.Path == Cartesian {
			expectedViolation = TupleInvention
		}
	}
	if e.Violations != expectedViolation {
		return fmt.Errorf("violation-evidence")
	}
	if e.Duplicate && e.Kind != "admit" && e.Kind != "deduplicate" && e.Kind != "receive-result" && e.Kind != "release" {
		return fmt.Errorf("duplicate-boundary")
	}
	hasResult := e.Kind == "apply" || e.Kind == "send-result" || e.Kind == "receive-result" || e.Kind == "record" || e.Kind == "consume"
	if hasResult != (e.Result == base.R1 || e.Result == base.R2) || !hasResult && e.Result != base.NoValue {
		return fmt.Errorf("result-vocabulary")
	}
	if e.Kind != "authorize" && e.Kind != "dispatch-refuse" && e.Refusal != "" {
		return fmt.Errorf("refusal-boundary")
	}
	switch e.Kind {
	case "commit":
		if o.committed {
			return fmt.Errorf("duplicate-commit")
		}
		if e.HasGrant {
			if err := l.issue(e); err != nil {
				return err
			}
		}
		o.committed = true
	case "retain-grant":
		if c.Scenario != "LateGrant" || e.Obligation != base.BA || !o.committed || len(l.grants) != 0 {
			return fmt.Errorf("issuance-order")
		}
		if err := l.issue(e); err != nil {
			return err
		}
	case "send-request":
		if !o.committed || o.request {
			return fmt.Errorf("request-custody")
		}
		o.request = true
	case "admit", "deduplicate":
		if !o.request || (e.Kind == "deduplicate") != o.admitted {
			return fmt.Errorf("request-admission")
		}
		o.admitted = true
		if !e.Duplicate {
			o.request = false
		}
	case "authorize":
		if !o.admitted || o.decision != Undecided {
			return fmt.Errorf("immutable-first-decision")
		}
		want := l.selectWitness(c, b)
		if e.Authority != want {
			return fmt.Errorf("authority-path-or-binding")
		}
		decision := Allowed
		reason := ""
		if want.Path == NoPath {
			decision = Refused
			reason = "no-covering-authority"
		}
		if e.Decision != decision || e.Refusal != reason {
			return fmt.Errorf("authorization-decision")
		}
		o.decision = decision
		o.witness = want
	case "dispatch-admit", "dispatch-refuse":
		if o.decision != Allowed || o.permit || o.applied || e.Authority != o.witness {
			return fmt.Errorf("admission-order-or-path")
		}
		live := l.witnessLive(o.witness, b)
		if e.Kind == "dispatch-admit" {
			if !live {
				return fmt.Errorf("fresh-admission-revoked")
			}
			o.permit = true
			o.permitEpoch = a.epoch
		} else {
			if live || o.admissionRefused || e.Refusal != "selected-path-revoked" {
				return fmt.Errorf("admission-refusal")
			}
			o.admissionRefused = true
		}
	case "apply":
		if !o.permit || o.permitEpoch != a.epoch || !o.admitted || o.applied || e.Authority != o.witness {
			return fmt.Errorf("application-without-opportunity")
		}
		outgoing := l.originals[e.Actor]
		if e.Actor == base.B && outgoing.committed && !outgoing.consumed {
			return fmt.Errorf("serial-service-guard")
		}
		o.applied = true
		o.result = e.Result
		o.retained = e.Result
		o.permit = false
		o.permitEpoch = 0
	case "send-result":
		if o.retained == base.NoValue || o.published != base.NoValue || e.Result != o.retained {
			return fmt.Errorf("publication-custody")
		}
		o.published = e.Result
	case "receive-result":
		if o.published == base.NoValue || e.Result != o.published {
			return fmt.Errorf("delivery-custody")
		}
		o.delivered = e.Result
		if !e.Duplicate {
			o.published = base.NoValue
		}
	case "record":
		if o.delivered == base.NoValue || o.recorded != base.NoValue || e.Result != o.delivered {
			return fmt.Errorf("recording-custody")
		}
		o.recorded = e.Result
	case "consume":
		if o.recorded == base.NoValue || o.consumed || e.Result != o.recorded || e.Result != o.result {
			return fmt.Errorf("consumption-custody")
		}
		o.consumed = true
	case "send-ack":
		if o.recorded == base.NoValue || o.ack {
			return fmt.Errorf("ack-before-durable-custody")
		}
		o.ack = true
	case "release":
		if !o.ack || o.recorded == base.NoValue {
			return fmt.Errorf("release-before-durable-custody")
		}
		o.retained = base.NoValue
		if !e.Duplicate {
			o.ack = false
		}
	default:
		return fmt.Errorf("unknown-event")
	}
	l.violations |= expectedViolation
	return nil
}
func ledgerEmpty(e Event) error {
	if e.Binding != (ServiceBinding{}) || e.Authority != (AuthorityWitness{}) || e.HasGrant || e.Grant != (GrantRecord{}) {
		return fmt.Errorf("lifecycle-is-not-authority")
	}
	if e.Result != base.NoValue || e.Duplicate || e.Decision != Undecided || e.Refusal != "" || e.Violations != 0 {
		return fmt.Errorf("lifecycle-descriptor")
	}
	return nil
}

// GraphAudit counts finite checks, not histories. Each node uses its checked
// shortest BFS prefix; every outgoing arc is interpreted from that ledger.
// Equal target ledgers establish consistency for that finite quotient only.
type GraphAudit struct {
	Config         Config
	Prefixes       int
	Edges          int
	UnsafePrefixes int
	UnsafeEdges    int
	PaidOriginals  [2]int
}

func AuditGraph(c Config) (GraphAudit, error) {
	if !ledgerConfig(c) {
		return GraphAudit{Config: c}, fmt.Errorf("unknown frozen configuration")
	}
	_, g := explore(c)
	return auditGraph(g)
}
func auditGraph(g *graph) (GraphAudit, error) {
	report := GraphAudit{Config: g.config}
	ledgers := make([]authorityLedger, len(g.nodes))
	ledgers[0].grants = map[string]GrantRecord{}
	for i, n := range g.nodes {
		if n.parent >= 0 {
			l := ledgers[n.parent]
			event := n.via.event(g.config, g.nodes[n.parent].state)
			if err := l.accept(g.config, event); err != nil {
				return report, fmt.Errorf("prefix %d: %w", i, err)
			}
			ledgers[i] = l
		}
		if ledgers[i].violations != n.state.Unsafe {
			return report, fmt.Errorf("prefix %d: unsafe marker differs", i)
		}
		for d, o := range ledgers[i].originals {
			if o.consumed != n.state.Base.Obligations[d].Consumed {
				return report, fmt.Errorf("prefix %d: paid original %d differs", i, d)
			}
		}
		report.Prefixes++
		if ledgers[i].violations != 0 {
			report.UnsafePrefixes++
		}
		for d, o := range ledgers[i].originals {
			if o.consumed {
				report.PaidOriginals[d]++
			}
		}
	}
	for i, n := range g.nodes {
		for _, arc := range g.edges(i) {
			l := ledgers[i]
			if err := l.accept(g.config, arc.event.event(g.config, n.state)); err != nil {
				return report, fmt.Errorf("edge %d -> %d: %w", i, arc.to, err)
			}
			report.Edges++
			if l.violations != 0 {
				report.UnsafeEdges++
			}
			target := ledgers[arc.to]
			equal := l.originals == target.originals && l.agents == target.agents && l.revoked == target.revoked
			if !equal || l.violations != target.violations || !maps.Equal(l.grants, target.grants) {
				return report, fmt.Errorf("edge %d -> %d: distinct independent ledger at merged target", i, arc.to)
			}
		}
	}
	return report, nil
}
