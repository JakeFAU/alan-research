// Package mutual authors the bounded Phase 3B reciprocal request-order witness.
package mutual

import (
	"fmt"

	"github.com/JakeFAU/alan-research/author"
	"github.com/JakeFAU/alan-research/runtime"
	"github.com/JakeFAU/alan-research/schema"
)

// Profile fixes which roots serve incoming work before awaiting their own reply.
type Profile string

const (
	SerialSerial Profile = "serial_serial"
	ServiceA     Profile = "service_a"
	ServiceB     Profile = "service_b"
	ServiceBoth  Profile = "service_both"
)

// New builds two ordinary retained roots; runtime service is supplied separately.
func New(profile Profile) (runtime.Provision, error) {
	switch profile {
	case SerialSerial, ServiceA, ServiceB, ServiceBoth:
	default:
		return runtime.Provision{}, fmt.Errorf("unknown reciprocal service profile %q", profile)
	}
	a, err := root("a", "b", profile == ServiceA || profile == ServiceBoth)
	if err != nil {
		return runtime.Provision{}, err
	}
	b, err := root("b", "a", profile == ServiceB || profile == ServiceBoth)
	if err != nil {
		return runtime.Provision{}, err
	}
	return (author.Application{
		"a": {Artifact: a, Inputs: inputs("a", 11)},
		"b": {Artifact: b, Inputs: inputs("b", 22)},
	}).Provision()
}

func inputs(subject string, value int64) map[string]schema.Value {
	return map[string]schema.Value{"own_id": schema.Text(subject + "-request"), "own_value": schema.Int(value)}
}

func root(subject, peer string, serviceFirst bool) (author.Artifact, error) {
	request := schema.Record{"request_id": {Kind: "string", MaxBytes: 32}, "request_value": {Kind: "int"}}
	response := schema.Record{"response_id": {Kind: "string", MaxBytes: 32}, "response_value": {Kind: "int"}}
	serve := []author.Step{
		author.Receive("incoming", map[string]string{"request_id": "peer_id", "request_value": "peer_value"}),
		author.Send("incoming", map[string]author.Expr{
			"response_id": author.Var("peer_id"), "response_value": author.Var("peer_value"),
		}),
	}
	consume := []author.Step{
		author.Receive("outgoing", map[string]string{"response_id": "reply_id", "response_value": "reply_value"}),
		author.IfEqual(author.Var("reply_id"), author.Var("own_id"),
			[]author.Step{author.IfEqual(author.Var("reply_value"), author.Var("own_value"),
				[]author.Step{author.Set("validated", author.Boolean(true))}, []author.Step{author.Failure()})},
			[]author.Step{author.Failure()}),
	}
	steps := []author.Step{author.Send("outgoing", map[string]author.Expr{
		"request_id": author.Var("own_id"), "request_value": author.Var("own_value"),
	})}
	if serviceFirst {
		steps = append(steps, serve...)
		steps = append(steps, consume...)
	} else {
		steps = append(steps, consume...)
		steps = append(steps, serve...)
	}
	// This unsupplied input keeps both original senders open through delivery.
	// Root termination would introduce a different source-authorization boundary.
	steps = append(steps, author.Receive("hold", map[string]string{"release": "release"}), author.Success())
	return author.Compile(author.Definition{
		Name:   "phase3b-mutual-" + subject,
		Inputs: schema.Record{"own_id": {Kind: "string", MaxBytes: 32}, "own_value": {Kind: "int"}},
		Values: map[string]schema.Value{"validated": schema.Bool(false)},
		Ports: map[string]author.Port{
			"outgoing": {
				Peer: peer, Obligation: "phase3b_" + subject + "_to_" + peer, Send: request, Receive: response,
			},
			"incoming": {
				Peer: peer, Obligation: "phase3b_" + peer + "_to_" + subject, Send: response, Receive: request,
			},
			"hold": {
				Peer: "local", Local: true, Obligation: "phase3b_hold", Receive: schema.Record{"release": {Kind: "bool"}},
			},
		},
		Steps: steps,
	})
}
