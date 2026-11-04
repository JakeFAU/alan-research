// Package authority supplies the bounded same-Store Phase 3B.2 authority probes.
package authority

import (
	"fmt"
	"github.com/JakeFAU/alan-research/author"
	"github.com/JakeFAU/alan-research/runtime"
	"github.com/JakeFAU/alan-research/schema"
)

type Scenario string

const (
	Local         Scenario = "local"
	Missing       Scenario = "missing"
	Claim         Scenario = "claim"
	ResultLaunder Scenario = "result_launder"
	MutualWait    Scenario = "mutual_wait"
	CounterA               = "phase3b2.counter-A.red"
	CounterB               = "phase3b2.counter-B.blue"
)

// New provisions static principal-local permissions. S2 operation-bound rooted
// delegation has no public runtime representation and is intentionally absent.
func New(s Scenario) (runtime.Provision, error) {
	switch s {
	case Local, Missing, Claim, ResultLaunder, MutualWait:
	default:
		return runtime.Provision{}, fmt.Errorf("unknown authority scenario %q", s)
	}
	aCapability := CounterA
	if s == Missing || s == Claim || s == ResultLaunder {
		aCapability = CounterB
	}
	a, err := root("a", "b", aCapability, s != MutualWait && s != ResultLaunder, s)
	if err != nil {
		return runtime.Provision{}, err
	}
	b, err := root("b", "a", CounterB, s == ResultLaunder, s)
	if err != nil {
		return runtime.Provision{}, err
	}
	claim := "ordinary-data"
	if s == Claim {
		claim = "grant:issuer=b;grantee=a;parent=b-root;permission=increment/counter-B/blue"
	}
	p, err := (author.Application{
		"a": {Artifact: a, Inputs: map[string]schema.Value{"own_id": schema.Text("a-request"), "claim": schema.Text(claim)}},
		"b": {Artifact: b, Inputs: map[string]schema.Value{"own_id": schema.Text("b-request"), "claim": schema.Text(claim)}},
	}).Provision()
	if err != nil {
		return p, err
	}
	// The author automatically grants declared effects. Remove the unwanted grant
	// through the public Provision before Create; never edit the positive database.
	filtered := p.Permissions[:0]
	for _, v := range p.Permissions {
		if v.Kind == "adapter" && v.Source == "a" && v.Target != CounterA {
			continue
		}
		filtered = append(filtered, v)
	}
	if aCapability != CounterA {
		filtered = append(filtered, runtime.Permission{Kind: "adapter", Source: "a", Target: CounterA, Operation: "increment", Revision: 1, Obligation: "increment"})
	}
	p.Permissions = filtered
	return p, nil
}

func root(subject, peer, capability string, serviceFirst bool, s Scenario) (author.Artifact, error) {
	request := schema.Record{"request_id": {Kind: "string", MaxBytes: 32}, "claim": {Kind: "string", MaxBytes: 256}}
	response := schema.Record{"response_id": {Kind: "string", MaxBytes: 32}, "mutation": {Kind: "int"}, "total": {Kind: "int"}, "service_operation": {Kind: "string", MaxBytes: 256}}
	hold := func() author.Step { return author.Receive("hold", map[string]string{"release": "release"}) }
	claim := author.Var("peer_claim")
	if s == ResultLaunder && subject == "a" {
		claim = author.Var("result_claim")
	}
	serve := []author.Step{
		author.Receive("incoming", map[string]string{"request_id": "peer_id", "claim": "peer_claim"}),
		author.Request("operation", "increment", map[string]author.Expr{"delta": author.Integer(1), "claim": claim, "receipt_mutation": author.Var("receipt_mutation")}, nil),
		author.Await("increment", "operation", map[string]string{"valid": "effect_valid", "disposition": "effect_disposition", "mutation": "mutation", "total": "total"}),
		author.IfEqual(author.Var("effect_valid"), author.Boolean(true), []author.Step{author.Send("incoming", map[string]author.Expr{"response_id": author.Var("peer_id"), "mutation": author.Var("mutation"), "total": author.Var("total"), "service_operation": author.Var("operation")})}, []author.Step{author.Set("denied", author.Boolean(true)), hold()}),
	}
	consume := []author.Step{
		author.Receive("outgoing", map[string]string{"response_id": "reply_id", "mutation": "reply_mutation", "total": "reply_total", "service_operation": "reply_operation"}),
		author.IfEqual(author.Var("reply_id"), author.Var("own_id"), []author.Step{author.Set("paid", author.Boolean(true))}, []author.Step{author.Failure()}),
		// This value is derived only after consuming the authentic B result. It
		// remains ordinary authored data and cannot confer B's adapter permission.
		author.Set("result_claim", author.String("result-from-b-is-not-authority")),
		author.Set("receipt_mutation", author.Var("reply_mutation")),
	}
	steps := []author.Step{author.Send("outgoing", map[string]author.Expr{"request_id": author.Var("own_id"), "claim": author.Var("claim")})}
	if serviceFirst {
		steps = append(steps, serve...)
		steps = append(steps, consume...)
	} else {
		steps = append(steps, consume...)
		steps = append(steps, serve...)
	}
	steps = append(steps, hold(), author.Success())
	return author.Compile(author.Definition{
		Name:    "phase3b2-" + string(s) + "-" + subject,
		Inputs:  schema.Record{"own_id": {Kind: "string", MaxBytes: 32}, "claim": {Kind: "string", MaxBytes: 256}},
		Values:  map[string]schema.Value{"paid": schema.Bool(false), "denied": schema.Bool(false), "receipt_mutation": schema.Int(0)},
		Effects: map[string]author.Effect{"increment": {Capability: capability, Input: schema.Record{"delta": {Kind: "int"}, "claim": {Kind: "string", MaxBytes: 256}, "receipt_mutation": {Kind: "int"}}, Output: schema.Record{"mutation": {Kind: "int"}, "total": {Kind: "int"}}}},
		Ports: map[string]author.Port{
			"outgoing": {Peer: peer, Obligation: "phase3b2_" + subject + "_to_" + peer, Send: request, Receive: response},
			"incoming": {Peer: peer, Obligation: "phase3b2_" + peer + "_to_" + subject, Send: response, Receive: request},
			"hold":     {Peer: "local", Local: true, Obligation: "phase3b2_hold", Receive: schema.Record{"release": {Kind: "bool"}}},
		}, Steps: steps,
	})
}
