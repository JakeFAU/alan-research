package authority

import (
	"context"
	"fmt"
	"github.com/JakeFAU/alan-research/effects"
	"github.com/JakeFAU/alan-research/examples/effectboundary"
)

// Registry fixes distinct resource/scope configurations in trusted deployment
// code. The request payload cannot select an endpoint, resource or scope.
func Registry(endpointA, endpointB string) map[string]effects.Adapter {
	return map[string]effects.Adapter{
		CounterA: counter{endpointA, CounterA, "counter-A", "red"},
		CounterB: counter{endpointB, CounterB, "counter-B", "blue"},
	}
}

type counter struct{ endpoint, capability, resource, scope string }

func (a counter) Prepare(r effects.Request) (effects.Prepared, error) {
	if r.Capability != a.capability || len(r.Payload) != 3 || r.Payload["claim"].Kind != "string" {
		return effects.Prepared{}, fmt.Errorf("wrong fixed capability or input")
	}
	inner := effects.Copy(r)
	delete(inner.Payload, "claim")
	delete(inner.Payload, "receipt_mutation")
	p, err := (effectboundary.Counter{Endpoint: a.endpoint}).Prepare(inner)
	p.Request = r
	p.Protocol = "phase3b2/counter/v1"
	p.Model = a.resource + "/" + a.scope
	return p, err
}
func (a counter) Execute(ctx context.Context, p effects.Prepared) effects.Candidate {
	return (effectboundary.Counter{Endpoint: a.endpoint}).Execute(ctx, p)
}
