package authorityservice

import (
	"fmt"
	"strings"

	base "github.com/JakeFAU/alan-research/research/mutualobligations"
)

// TraceControl retains the entire input, independent verdict, and expected
// distinction between candidate validity and explicit-consent safety.
type TraceControl struct {
	Name              string
	Config            Config
	Events            []Event
	WantProtocolValid bool
	WantAuthoritySafe bool
	WantCondition     string
	Audit             Audit
	Error             string
}
type ChainControl struct {
	Name          string
	Binding       ServiceBinding
	Claims        []GrantRecord
	Selected      string
	WantCondition string
	Error         string
}
type OracleCorpusReport struct {
	Histories []TraceControl
	Chains    []ChainControl
}

// OracleCorpus preserves controls reproducibly. Successors is used only to
// generate baseline event lists; AuditTrace and grantPath interpret independently.
func OracleCorpus() (OracleCorpusReport, error) {
	out := OracleCorpusReport{Histories: []TraceControl{}, Chains: []ChainControl{}}
	c := Config{Policy: "S2", Scenario: "Explicit"}
	original, err := corpusComplete(c, base.R1, base.R1, false)
	if err != nil {
		return out, err
	}
	add := func(name string, config Config, events []Event, valid, safe bool, reason string) {
		a, e := AuditTrace(config, events)
		message := ""
		if e != nil {
			message = e.Error()
		}
		out.Histories = append(out.Histories, TraceControl{Name: name, Config: config, Events: events,
			WantProtocolValid: valid, WantAuthoritySafe: safe, WantCondition: reason, Audit: a, Error: message})
	}
	for _, a := range []base.Value{base.R1, base.R2} {
		for _, b := range []base.Value{base.R1, base.R2} {
			events, e := corpusComplete(c, a, b, false)
			if e != nil {
				return out, e
			}
			add(fmt.Sprintf("complete-values-%d-%d", a, b), c, events, true, true, "")
		}
	}
	interrupted, e := corpusComplete(c, base.R2, base.R2, true)
	if e != nil {
		return out, e
	}
	add("released-receipt-both-crashes-then-paid", c, interrupted, true, true, "")
	for _, mutation := range []struct {
		name, reason string
		change       func([]Event)
	}{
		{"root-identity-collision", "root-identity-collision", func(es []Event) { es[0].Grant.ID = "B/g1:b" }},
		{"missing-parent", "missing-origin", func(es []Event) { es[0].Grant.Parent = "result:BA:R1" }},
		{"self-parent-cycle", "grant-cycle", func(es []Event) { es[0].Grant.Parent = es[0].Grant.ID }},
		{"mutated-root", "mutated-root", func(es []Event) { es[0].Grant.Root = "A/g1:a" }},
		{"stale-grantee", "identity-adjacency", func(es []Event) { es[0].Grant.GranteeGeneration = 2 }},
		{"stale-issuer", "identity-adjacency", func(es []Event) { es[0].Grant.IssuerGeneration = 2 }},
		{"enlarged-permission-set", "tuple-attenuation", func(es []Event) { es[0].Grant.Permissions |= 1 << APermission }},
		{"invented-tuple", "tuple-attenuation", func(es []Event) { es[0].Grant.Permissions = 1 << InventedPermission }},
		{"rebound-original", "service-binding", func(es []Event) { es[0].Grant.Binding.Original = "AB" }},
		{"rebound-operation", "service-binding", func(es []Event) { es[0].Grant.Binding.Operation = 2 }},
		{"stale-binding-generation", "service-binding", func(es []Event) { es[0].Grant.Binding.CallerGeneration = 2 }},
		{"wrong-authority-revision", "grant-revision", func(es []Event) { es[0].Grant.Revision = 1 }},
		{"enabled-redelegation", "redelegation-disposition", func(es []Event) { es[0].Grant.NoRedelegation = false }},
		{"admission-epoch", "activation-binding", func(es []Event) { es[4].Activation = 1 }},
		{"result-as-authority", "result-is-not-authority", func(es []Event) { es[6].Authority = es[3].Authority }},
		{"swapped-equal-value-provenance", "original-binding", func(es []Event) { es[9].Obligation = base.AB }},
		{"mutated-authorizer", "authority-path-or-binding", func(es []Event) { es[3].Authority.Principal = base.B }},
		{"forged-retained-decision", "immutable-decision", func(es []Event) { es[7].Decision = Refused }},
	} {
		changed := append([]Event{}, original...)
		mutation.change(changed)
		add(mutation.name, c, changed, false, false, mutation.reason)
	}
	for _, candidate := range []Config{{Policy: "S3", Scenario: "Missing"}, {Policy: "Cartesian", Scenario: "ScopeMix"}} {
		events, e := corpusComplete(candidate, base.R1, base.R1, false)
		if e != nil {
			return out, e
		}
		reason := "ambient-consent"
		if candidate.Policy == "Cartesian" {
			reason = "tuple-invention"
		}
		add("candidate-valid-"+candidate.Policy, candidate, events, true, false, reason)
	}
	// A permanent first refusal is not discharged by later grant retention.
	late := Config{Policy: "S2", Scenario: "LateGrant"}
	lateEvents, e := corpusSequence(late, []corpusAction{
		{kind: "commit", d: base.BA}, {kind: "send-request", d: base.BA}, {kind: "admit", d: base.BA},
		{kind: "authorize", d: base.BA}, {kind: "retain-grant", d: base.BA},
		{kind: "crash", d: base.NoObligation, actor: base.A}, {kind: "recover", d: base.NoObligation, actor: base.A},
	})
	if e != nil {
		return out, e
	}
	add("late-grant-sticky-refusal", late, lateEvents, true, true, "")
	retry := append([]Event{}, lateEvents...)
	fake := lateEvents[3]
	fake.Activation = 1
	retry = append(retry, fake)
	add("late-grant-cannot-rewrite-first-refusal", late, retry, false, false, "immutable-first-decision")
	rev := Config{Policy: "S2", Scenario: "Revocable"}
	starts := []corpusAction{{kind: "commit", d: base.BA}, {kind: "send-request", d: base.BA}, {kind: "admit", d: base.BA}, {kind: "authorize", d: base.BA}}
	schedules := []struct {
		name    string
		actions []corpusAction
	}{
		{"revoked-before-issuance", []corpusAction{{kind: "revoke", d: base.NoObligation, actor: base.B}, {kind: "commit", d: base.BA}, {kind: "send-request", d: base.BA}, {kind: "admit", d: base.BA}, {kind: "authorize", d: base.BA}}},
		{"revoked-before-fresh-admission", append(append([]corpusAction{}, starts...), corpusAction{kind: "revoke", d: base.NoObligation, actor: base.B}, corpusAction{kind: "dispatch-refuse", d: base.BA})},
		{"pre-revocation-admission-applies", append(append([]corpusAction{}, starts...), corpusAction{kind: "dispatch-admit", d: base.BA}, corpusAction{kind: "revoke", d: base.NoObligation, actor: base.B}, corpusAction{kind: "apply", d: base.BA, value: base.R2})},
		{"crashed-opportunity-cannot-readmit-revoked", append(append([]corpusAction{}, starts...), corpusAction{kind: "dispatch-admit", d: base.BA}, corpusAction{kind: "revoke", d: base.NoObligation, actor: base.B}, corpusAction{kind: "crash", d: base.NoObligation, actor: base.A}, corpusAction{kind: "recover", d: base.NoObligation, actor: base.A}, corpusAction{kind: "dispatch-refuse", d: base.BA})},
		{"fresh-activation-readmission", append(append([]corpusAction{}, starts...), corpusAction{kind: "dispatch-admit", d: base.BA}, corpusAction{kind: "crash", d: base.NoObligation, actor: base.A}, corpusAction{kind: "recover", d: base.NoObligation, actor: base.A}, corpusAction{kind: "dispatch-admit", d: base.BA}, corpusAction{kind: "apply", d: base.BA, value: base.R1})},
	}
	for _, schedule := range schedules {
		events, e := corpusSequence(rev, schedule.actions)
		if e != nil {
			return out, e
		}
		add(schedule.name, rev, events, true, true, "")
		if schedule.name == "revoked-before-fresh-admission" {
			bad := append([]Event{}, events...)
			bad[len(bad)-1].Kind = "dispatch-admit"
			bad[len(bad)-1].Refusal = ""
			add("revoked-path-forged-fresh-admission", rev, bad, false, false, "fresh-admission-revoked")
		}
	}
	// Longer chains are malformed ledger controls, not generated operational
	// histories: no valid arbitrary redelegation protocol is claimed.
	g := original[0].Grant
	makeChain := func(name string, claims []GrantRecord, selected, reason string) {
		l := authorityLedger{grants: map[string]GrantRecord{}}
		for _, claim := range claims {
			l.grants[claim.ID] = claim
		}
		_, e := l.grantPath(selected, g.Binding, true)
		message := ""
		if e != nil {
			message = e.Error()
		}
		out.Chains = append(out.Chains, ChainControl{Name: name, Binding: g.Binding, Claims: claims, Selected: selected, WantCondition: reason, Error: message})
	}
	makeChain("one-hop-rooted-accepted", []GrantRecord{g}, g.ID, "")
	rejected := g
	rejected.Accepted = false
	makeChain("rejected-issuance-confers-nothing", []GrantRecord{rejected}, g.ID, "rejected-issuance")
	first, second := g, g
	first.ID = "cycle-1"
	first.Parent = "cycle-2"
	second.ID = "cycle-2"
	second.Parent = "cycle-1"
	makeChain("two-node-rootless-cycle", []GrantRecord{first, second}, first.ID, "grant-cycle")
	returned := g
	returned.ID = "returned"
	returned.Parent = g.ID
	returned.Issuer = base.A
	returned.Grantee = base.B
	makeChain("unauthorized-return-delegation", []GrantRecord{g, returned}, returned.ID, "no-redelegation")
	for _, control := range out.Histories {
		conditions := strings.Join(control.Audit.Conditions, " ") + " " + control.Error
		matches := control.Audit.ProtocolValid == control.WantProtocolValid && control.Audit.AuthoritySafe == control.WantAuthoritySafe
		if !matches || !strings.Contains(conditions, control.WantCondition) {
			return out, fmt.Errorf("control %s: unexpected audit %+v (%s)", control.Name, control.Audit, control.Error)
		}
	}
	for _, control := range out.Chains {
		if control.Error != control.WantCondition {
			return out, fmt.Errorf("chain %s: %s, want %s", control.Name, control.Error, control.WantCondition)
		}
	}
	return out, nil
}

type corpusAction struct {
	kind  string
	d     base.Obligation
	actor base.Agent
	value base.Value
}

func corpusSequence(c Config, actions []corpusAction) ([]Event, error) {
	s := Initial(c)
	events := []Event{}
	for _, action := range actions {
		found := false
		for _, edge := range Successors(c, s) {
			e := edge.Event
			if e.Kind != action.kind || e.Obligation != action.d || action.value != base.NoValue && e.Result != action.value {
				continue
			}
			if action.d == base.NoObligation && e.Actor != action.actor {
				continue
			}
			s = edge.State
			events = append(events, e)
			found = true
			break
		}
		if !found {
			return events, fmt.Errorf("corpus generation %s missing %s/%d", c.Name(), action.kind, action.d)
		}
	}
	return events, nil
}
func corpusComplete(c Config, ab, ba base.Value, interrupt bool) ([]Event, error) {
	actions := []corpusAction{}
	// B's serial outgoing wait is paid before it applies AB; A may service BA.
	for _, d := range []base.Obligation{base.BA, base.AB} {
		value := ab
		if d == base.BA {
			value = ba
		}
		for _, kind := range []string{"commit", "send-request", "admit", "authorize", "dispatch-admit", "apply", "send-result", "receive-result", "record"} {
			action := corpusAction{kind: kind, d: d}
			if kind == "apply" {
				action.value = value
			}
			actions = append(actions, action)
		}
		if interrupt && d == base.BA {
			actions = append(actions, corpusAction{kind: "send-ack", d: d}, corpusAction{kind: "release", d: d},
				corpusAction{kind: "crash", d: base.NoObligation, actor: base.A}, corpusAction{kind: "crash", d: base.NoObligation, actor: base.B},
				corpusAction{kind: "recover", d: base.NoObligation, actor: base.A}, corpusAction{kind: "recover", d: base.NoObligation, actor: base.B})
		}
		actions = append(actions, corpusAction{kind: "consume", d: d})
	}
	return corpusSequence(c, actions)
}
