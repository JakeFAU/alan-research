package delegationcorrespondence

import (
	"fmt"
	m "github.com/JakeFAU/alan-research/research/authorityservice"
	b "github.com/JakeFAU/alan-research/research/mutualobligations"
)

type ModelCase struct {
	ID, Class, MutationOf, OriginalSHA256, SHA256, Mutation string
	Config                                                  m.Config
	Events, CompletionSuffix                                []m.Event
	Audit                                                   m.Audit
	Error                                                   string
	FirstFailedEdge                                         int
}
type CorpusReport struct {
	Version string
	Classes []string
	Cases   []ModelCase
	Prior   m.OracleCorpusReport
	Limits  []string
}
type action struct {
	kind  string
	d     b.Obligation
	value b.Value
	actor b.Agent
}

func follow(config m.Config, actions []action) ([]m.Event, error) {
	s := m.Initial(config)
	es := []m.Event{}
	for _, a := range actions {
		found := []m.Transition{}
		for _, t := range m.Successors(config, s) {
			e := t.Event
			if e.Kind != a.kind || e.Obligation != a.d || e.Duplicate || a.value != b.NoValue && e.Result != a.value || a.d == b.NoObligation && e.Actor != a.actor {
				continue
			}
			found = append(found, t)
		}
		if len(found) != 1 {
			return nil, fmt.Errorf("%s %s/%d candidates=%d", config.Name(), a.kind, a.d, len(found))
		}
		s = found[0].State
		es = append(es, found[0].Event)
	}
	return es, nil
}
func completeModel(ab, ba b.Value) ([]m.Event, error) {
	c := m.Config{Policy: "S2", Scenario: "LateGrant"}
	as := []action{}
	for _, d := range []b.Obligation{b.AB, b.BA} {
		for _, k := range []string{"commit", "send-request", "admit"} {
			as = append(as, action{kind: k, d: d})
		}
	}
	as = append(as, action{kind: "retain-grant", d: b.BA})
	for _, d := range []b.Obligation{b.BA, b.AB} {
		v := ab
		if d == b.BA {
			v = ba
		}
		for _, k := range []string{"authorize", "dispatch-admit", "apply", "send-result", "receive-result", "record", "consume"} {
			a := action{kind: k, d: d}
			if k == "apply" {
				a.value = v
			}
			as = append(as, a)
		}
	}
	return follow(c, as)
}
func GenerateCorpus() (CorpusReport, error) {
	out := CorpusReport{Version: "alan/delegation-correspondence/v1", Classes: []string{"valid_complete", "valid_incomplete", "valid_authority_blocked", "invalid_missing_authority", "invalid_scope_amplification", "invalid_request_or_generation", "invalid_revocation_timing", "invalid_causal_order"}, Limits: []string{"Finite histories supplement, never replace sixteen unchanged graph audits", "Separate grant and revocation have no direct unchanged configuration"}}
	prior, e := m.OracleCorpus()
	if e != nil {
		return out, e
	}
	out.Prior = prior
	add := func(id, class string, config m.Config, es []m.Event, parent, mutation string) {
		a, err := m.AuditTrace(config, es)
		s := ""
		if err != nil {
			s = err.Error()
		}
		failed := -1
		if err != nil {
			failed = a.EventsChecked
		} else if !a.AuthoritySafe {
			for n := 1; n <= len(es); n++ {
				prefix, _ := m.AuditTrace(config, es[:n])
				if !prefix.AuthoritySafe {
					failed = n - 1
					break
				}
			}
		}
		out.Cases = append(out.Cases, ModelCase{ID: id, Class: class, Config: config, Events: es, Audit: a, Error: s, SHA256: sha(bytesJSON(es)), MutationOf: parent, Mutation: mutation, FirstFailedEdge: failed})
	}
	config := m.Config{Policy: "S2", Scenario: "LateGrant"}
	for _, ab := range []b.Value{b.R1, b.R2} {
		for _, ba := range []b.Value{b.R1, b.R2} {
			es, e := completeModel(ab, ba)
			if e != nil {
				return out, e
			}
			need(es[7].Obligation == b.BA && es[7].Authority.Path == m.Delegated, "model_generator", "", "positive is not delegated BA")
			add(fmt.Sprintf("reciprocal-%d-%d", ab, ba), "valid_complete", config, es, "", "")
		}
	}
	positive := out.Cases[0]
	add("reciprocal-prefix", "valid_incomplete", config, positive.Events[:8], "", "")
	out.Cases[len(out.Cases)-1].CompletionSuffix = positive.Events[8:]
	for _, t := range prior.Histories {
		class := ""
		switch t.Name {
		case "late-grant-sticky-refusal":
			class = "valid_authority_blocked"
		case "candidate-valid-S3", "missing-parent", "result-as-authority":
			class = "invalid_missing_authority"
		case "enlarged-permission-set", "invented-tuple", "candidate-valid-Cartesian":
			class = "invalid_scope_amplification"
		case "rebound-original", "rebound-operation", "stale-binding-generation", "stale-grantee", "stale-issuer", "swapped-equal-value-provenance":
			class = "invalid_request_or_generation"
		case "revoked-path-forged-fresh-admission":
			class = "invalid_revocation_timing"
		case "late-grant-cannot-rewrite-first-refusal":
			class = "invalid_causal_order"
		}
		if class != "" {
			add(t.Name, class, t.Config, t.Events, "prior-oracle-corpus", t.WantCondition)
		}
	}
	missing := append([]m.Event{}, positive.Events[:6]...)
	missing = append(missing, positive.Events[7:]...)
	add("missing-retained-grant", "invalid_missing_authority", config, missing, positive.ID, "delete actual retain-grant; preserve results")
	wrong := append([]m.Event{}, positive.Events...)
	wrong[7], wrong[8] = wrong[8], wrong[7]
	add("gate-before-first-decision", "invalid_causal_order", config, wrong, positive.ID, "swap authorization and gate; preserve all results")
	parents := map[string]string{}
	for _, t := range prior.Histories {
		parents[t.Name] = sha(bytesJSON(t.Events))
	}
	for j := range out.Cases {
		c := &out.Cases[j]
		if c.MutationOf != "prior-oracle-corpus" {
			continue
		}
		switch c.ID {
		case "late-grant-sticky-refusal", "candidate-valid-S3", "candidate-valid-Cartesian":
			c.MutationOf = ""
		case "late-grant-cannot-rewrite-first-refusal":
			c.MutationOf = "prior:late-grant-sticky-refusal"
			c.OriginalSHA256 = parents["late-grant-sticky-refusal"]
		case "revoked-path-forged-fresh-admission":
			c.MutationOf = "prior:revoked-before-fresh-admission"
			c.OriginalSHA256 = parents["revoked-before-fresh-admission"]
		default:
			c.MutationOf = "prior:" + prior.Histories[0].Name
			c.OriginalSHA256 = sha(bytesJSON(prior.Histories[0].Events))
		}
	}
	for j := range out.Cases {
		if out.Cases[j].MutationOf == positive.ID {
			out.Cases[j].OriginalSHA256 = positive.SHA256
		}
	}
	return out, nil
}
