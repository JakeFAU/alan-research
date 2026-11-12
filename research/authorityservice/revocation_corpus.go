package authorityservice

import (
	"fmt"
	"strings"

	base "github.com/JakeFAU/alan-research/research/mutualobligations"
)

// GrantRevocationOracleCorpus generates concrete controls with the operational
// relation, then independently interprets their complete descriptors. Forged
// histories are mutations of lawful controls, not operational successors.
func GrantRevocationOracleCorpus() (OracleCorpusReport, error) {
	c := GrantRevocationConfig()
	out := OracleCorpusReport{Histories: []TraceControl{}, Chains: []ChainControl{}}
	add := func(name string, events []Event, valid bool, reason string) {
		a, err := AuditTrace(c, events)
		message := ""
		if err != nil {
			message = err.Error()
		}
		out.Histories = append(out.Histories, TraceControl{Name: name, Config: c, Events: events,
			WantProtocolValid: valid, WantAuthoritySafe: valid, WantCondition: reason, Audit: a, Error: message})
	}
	start := []corpusAction{}
	for _, kind := range []string{"commit", "retain-grant", "send-request", "admit", "authorize"} {
		start = append(start, corpusAction{kind: kind, d: base.BA})
	}
	revoke := corpusAction{kind: "revoke", d: base.NoObligation, actor: base.B}
	permit := corpusAction{kind: "dispatch-admit", d: base.BA}
	apply := corpusAction{kind: "apply", d: base.BA, value: base.R2}
	crash := corpusAction{kind: "crash", d: base.NoObligation, actor: base.A}
	recoverA := corpusAction{kind: "recover", d: base.NoObligation, actor: base.A}
	refusal := corpusAction{kind: "dispatch-refuse", d: base.BA}
	appendActions := func(extra ...corpusAction) []corpusAction {
		return append(append([]corpusAction{}, start...), extra...)
	}
	schedules := []struct {
		name    string
		actions []corpusAction
	}{
		{"historical-allow-fresh-refusal", appendActions(revoke, refusal)},
		{"old-permit-survives-revocation", appendActions(permit, revoke, apply)},
		{"crash-destroys-revoked-permit", appendActions(permit, revoke, crash, recoverA, refusal)},
		{"current-path-recovery-readmission", appendActions(permit, crash, recoverA, permit, apply)},
		{"rejected-issuance", []corpusAction{revoke, {kind: "commit", d: base.BA}, {kind: "retain-grant", d: base.BA},
			{kind: "send-request", d: base.BA}, {kind: "admit", d: base.BA}, {kind: "authorize", d: base.BA}}},
		{"sticky-first-refusal", []corpusAction{{kind: "commit", d: base.BA}, {kind: "send-request", d: base.BA},
			{kind: "admit", d: base.BA}, {kind: "authorize", d: base.BA}, {kind: "retain-grant", d: base.BA}, revoke, crash, recoverA}},
	}
	complete := appendActions(permit, revoke, apply)
	for _, kind := range []string{"send-result", "receive-result", "record", "consume"} {
		complete = append(complete, corpusAction{kind: kind, d: base.BA})
	}
	for _, kind := range []string{"commit", "send-request", "admit", "authorize", "dispatch-admit", "apply", "send-result", "receive-result", "record", "consume"} {
		complete = append(complete, corpusAction{kind: kind, d: base.AB, value: func() base.Value {
			if kind == "apply" {
				return base.R1
			}
			return base.NoValue
		}()})
	}
	schedules = append(schedules, struct {
		name    string
		actions []corpusAction
	}{"post-revocation-both-authentic-consumptions", complete})
	controls := map[string][]Event{}
	for _, s := range schedules {
		events, err := corpusSequence(c, s.actions)
		if err != nil {
			return out, err
		}
		controls[s.name] = events
		add(s.name, events, true, "")
	}
	original := controls["post-revocation-both-authentic-consumptions"]
	for _, attack := range []struct {
		name, reason string
		change       func([]Event)
	}{
		{"grant-id", "grant-identity", func(es []Event) { es[1].Grant.ID = "grant:BA:other" }},
		{"root-c-substitution", "mutated-root", func(es []Event) { es[1].Grant.Root = "B/g1:c" }},
		{"parent-c-substitution", "mutated-root", func(es []Event) { es[1].Grant.Parent = "B/g1:c" }},
		{"missing-ancestor", "missing-origin", func(es []Event) { es[1].Grant.Parent = "unknown" }},
		{"rootless-cycle", "grant-cycle", func(es []Event) { es[1].Grant.Parent = es[1].Grant.ID }},
		{"issuer-generation", "identity-adjacency", func(es []Event) { es[1].Grant.IssuerGeneration = 2 }},
		{"grantee-generation", "identity-adjacency", func(es []Event) { es[1].Grant.GranteeGeneration = 2 }},
		{"complete-tuple-attenuation", "tuple-attenuation", func(es []Event) { es[1].Grant.Permissions |= 1 << ABOwnPermission }},
		{"grant-original", "service-binding", func(es []Event) { es[1].Grant.Binding.Original = "AB" }},
		{"grant-operation", "service-binding", func(es []Event) { es[1].Grant.Binding.Operation = 2 }},
		{"grant-binding-generation", "service-binding", func(es []Event) { es[1].Grant.Binding.CallerGeneration = 2 }},
		{"grant-scope", "service-binding", func(es []Event) { es[1].Grant.Binding.Permission.Scope = "red" }},
		{"grant-redelegation", "redelegation-disposition", func(es []Event) { es[1].Grant.NoRedelegation = false }},
		{"grant-revision", "grant-revision", func(es []Event) { es[1].Grant.Revision = 1 }},
		{"false-rejected-current-issuance", "unjustified-issuance-rejection", func(es []Event) { es[1].Grant.Accepted = false }},
		{"first-allow-rewritten", "immutable-decision", func(es []Event) { es[5].Decision = Refused }},
		{"wrong-admission-epoch", "activation-binding", func(es []Event) { es[5].Activation = 1 }},
		{"wrong-revoke-actor", "revocation-order", func(es []Event) { es[6].Actor = base.A }},
		{"wrong-revoke-era", "authority-revision-order", func(es []Event) { es[6].Revision = 1 }},
		{"revoke-names-original", "revocation-order", func(es []Event) { es[6].Obligation = base.BA }},
		{"revoked-application-changed-path", "application-without-opportunity", func(es []Event) { es[7].Authority.Root = "B/g1:c" }},
		{"ab-root-b-substitution", "authority-path-or-binding", func(es []Event) { es[15].Authority.Root = "B/g1:b" }},
		{"ab-complete-tuple-substitution", "original-binding", func(es []Event) { es[15].Binding.Permission.Resource = "counter-B" }},
	} {
		events := append([]Event{}, original...)
		attack.change(events)
		add(attack.name, events, false, attack.reason)
	}
	bad := append([]Event{}, controls["historical-allow-fresh-refusal"]...)
	bad[len(bad)-1].Kind = "dispatch-admit"
	bad[len(bad)-1].Refusal = ""
	add("forged-fresh-admission-after-revoke", bad, false, "fresh-admission-revoked")
	bad = append([]Event{}, controls["rejected-issuance"]...)
	bad[2].Grant.Accepted = true
	add("forged-accepted-post-revocation-issuance", bad, false, "revoked-origin")
	bad = append([]Event{}, controls["sticky-first-refusal"]...)
	fake := bad[3]
	fake.Activation = 1
	fake.Revision = 1
	bad = append(bad, fake)
	add("cannot-repair-first-refusal", bad, false, "immutable-first-decision")
	bad = append([]Event{}, controls["crash-destroys-revoked-permit"]...)
	bad = bad[:len(bad)-1]
	fake = controls["old-permit-survives-revocation"][7]
	fake.Activation = 1
	bad = append(bad, fake)
	add("cannot-apply-crashed-permit", bad, false, "application-without-opportunity")
	for _, control := range out.Histories {
		if control.Audit.ProtocolValid != control.WantProtocolValid || control.Audit.AuthoritySafe != control.WantAuthoritySafe ||
			!strings.Contains(strings.Join(control.Audit.Conditions, " ")+" "+control.Error, control.WantCondition) {
			return out, fmt.Errorf("control %s: unexpected audit %+v (%s)", control.Name, control.Audit, control.Error)
		}
	}
	return out, nil
}
