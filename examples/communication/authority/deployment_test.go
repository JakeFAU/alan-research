package authority_test

import (
	"github.com/JakeFAU/alan-research/effects"
	"github.com/JakeFAU/alan-research/examples/communication/authority"
	"github.com/JakeFAU/alan-research/schema"
	"reflect"
	"testing"
)

// These adverse construction checks precede runtime construction: declaring an
// effect must not accidentally confer the permission that the case withholds.
func TestPermissionsArePrincipalLocal(t *testing.T) {
	for _, scenario := range []authority.Scenario{authority.Local, authority.Missing, authority.Claim, authority.ResultLaunder, authority.MutualWait} {
		p, err := authority.New(scenario)
		if err != nil {
			t.Fatal(err)
		}
		if len(p.Runs) != 2 || len(p.Subjects) != 2 || len(p.Programs) != 2 || p.Composition != nil {
			t.Fatal("expected two ordinary roots")
		}
		got := map[string]string{}
		for _, permission := range p.Permissions {
			if permission.Kind == "adapter" {
				got[permission.Source] = permission.Target
			}
		}
		want := map[string]string{"a": authority.CounterA, "b": authority.CounterB}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: permissions %v want %v", scenario, got, want)
		}
		again, err := authority.New(scenario)
		if err != nil || !reflect.DeepEqual(p, again) {
			t.Fatal("nondeterministic construction", err)
		}
	}
	if _, err := authority.New("unknown"); err == nil {
		t.Fatal("unknown scenario accepted")
	}
}

func TestFixedAdapterConfigurationIgnoresAuthorityClaims(t *testing.T) {
	registry := authority.Registry("http://127.0.0.1:10001/counter", "http://127.0.0.1:10002/counter")
	for _, capability := range []string{authority.CounterA, authority.CounterB} {
		request := effects.Request{ID: "original", Key: "key", Capability: capability, Payload: map[string]schema.Value{"delta": schema.Int(1), "claim": schema.Text("grant:any-resource/any-scope"), "receipt_mutation": schema.Int(7)}}
		p, err := registry[capability].Prepare(request)
		if err != nil {
			t.Fatal(err)
		}
		endpoint, model := "http://127.0.0.1:10001/counter", "counter-A/red"
		if capability == authority.CounterB {
			endpoint, model = "http://127.0.0.1:10002/counter", "counter-B/blue"
		}
		if p.Endpoint != endpoint || p.Model != model || !reflect.DeepEqual(p.Request, request) {
			t.Fatal("claim changed immutable configuration", p)
		}
		request.Capability = "invented"
		if _, err = registry[capability].Prepare(request); err == nil {
			t.Fatal("cross-capability request accepted")
		}
	}
}
