package author

import (
	"encoding/json"
	"github.com/JakeFAU/alan-research/internal/executor"
	"github.com/JakeFAU/alan-research/schema"
	"strings"
	"testing"
)

func TestInvocationHandleRoundTripThroughList(t *testing.T) {
	leaf := invocationTemplate(t)
	handles := schema.Field{Kind: "list", MaxItems: 2, MaxBytes: 8192, Element: schema.Record{"operation": {Kind: "string", MaxBytes: 4096}}}
	_, err := Compile(Definition{Invocations: map[string]Template{"child": leaf}, Collections: map[string]schema.Field{"handles": handles}, Steps: []Step{
		Invoke("pending", "child", nil, nil), Append("handles", map[string]Expr{"operation": Var("pending")}),
		ForEach(Var("handles"), map[string]string{"operation": "saved"}, AwaitInvocation("child", "saved", InvocationOutcomeBindings{Status: "status", Reason: "reason", ResponsePresent: "present"})), Success(),
	}})
	if err != nil {
		t.Fatal(err)
	}
}
func TestInvocationAwaitRejectsKnownWrongDeclarationAndNonString(t *testing.T) {
	leaf := invocationTemplate(t)
	for _, steps := range [][]Step{
		{Invoke("pending", "other", nil, nil), AwaitInvocation("child", "pending", InvocationOutcomeBindings{Status: "status", Reason: "reason", ResponsePresent: "present"}), Success()},
		{Set("pending", Integer(7)), AwaitInvocation("child", "pending", InvocationOutcomeBindings{Status: "status", Reason: "reason", ResponsePresent: "present"}), Success()},
	} {
		if _, err := Compile(Definition{Invocations: map[string]Template{"child": leaf, "other": leaf}, Steps: steps}); err == nil {
			t.Fatal("invalid handle accepted")
		}
	}
}
func TestAdapterAwaitStillRejectsPlainString(t *testing.T) {
	_, err := Compile(Definition{Effects: map[string]Effect{"effect": {Capability: "test", Input: schema.Record{}, Output: schema.Record{}}}, Values: map[string]schema.Value{"saved": schema.Text("opaque")}, Steps: []Step{Await("effect", "saved", nil), Success()}})
	if err == nil {
		t.Fatal("adapter Await lost request declaration proof")
	}
}

func TestInvocationHandleProvenanceRejectsUntrustedStrings(t *testing.T) {
	leaf := invocationTemplate(t)
	handles := schema.Field{Kind: "list", MaxItems: 3, MaxBytes: 8192, Element: schema.Record{"operation": {Kind: "string", MaxBytes: 4096}, "label": {Kind: "string", MaxBytes: 64}}}
	await := AwaitInvocation("child", "saved", InvocationOutcomeBindings{Status: "status", Reason: "reason", ResponsePresent: "present"})
	cases := map[string][]Step{
		"input":             {Set("saved", Var("input")), await, Success()},
		"literal row":       {Append("handles", map[string]Expr{"operation": String("forged"), "label": String("ok")}), ForEach(Var("handles"), map[string]string{"operation": "saved"}, await), Success()},
		"mixed rows":        {Invoke("pending", "child", nil, nil), Append("handles", map[string]Expr{"operation": Var("pending"), "label": String("ok")}), Append("handles", map[string]Expr{"operation": String("forged"), "label": String("bad")}), ForEach(Var("handles"), map[string]string{"operation": "saved"}, await), Success()},
		"wrong declaration": {Invoke("pending", "other", nil, nil), Append("handles", map[string]Expr{"operation": Var("pending"), "label": String("ok")}), ForEach(Var("handles"), map[string]string{"operation": "saved"}, await), Success()},
	}
	for name, steps := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Compile(Definition{Inputs: schema.Record{"input": {Kind: "string", MaxBytes: 4096}}, Invocations: map[string]Template{"child": leaf, "other": leaf}, Collections: map[string]schema.Field{"handles": handles}, Steps: steps})
			if err == nil {
				t.Fatal("untrusted invocation handle accepted")
			}
		})
	}
}

func TestInvocationHandleProvenanceTracksFieldsAndBranches(t *testing.T) {
	leaf := invocationTemplate(t)
	handles := schema.Field{Kind: "list", MaxItems: 2, MaxBytes: 8192, Element: schema.Record{"operation": {Kind: "string", MaxBytes: 4096}, "label": {Kind: "string", MaxBytes: 64}}}
	await := AwaitInvocation("child", "saved", InvocationOutcomeBindings{Status: "status", Reason: "reason", ResponsePresent: "present"})
	valid := []Step{Invoke("pending", "child", nil, nil), Set("copy", Var("pending")), Append("handles", map[string]Expr{"operation": Var("copy"), "label": String("untrusted")}), ForEach(Var("handles"), map[string]string{"operation": "saved", "label": "label"}, await), Success()}
	_, err := Compile(Definition{Invocations: map[string]Template{"child": leaf}, Collections: map[string]schema.Field{"handles": handles}, Steps: valid})
	if err != nil {
		t.Fatal(err)
	}
	invalid := []Step{Invoke("pending", "child", nil, nil), IfEqual(Var("choice"), String("yes"), []Step{Append("handles", map[string]Expr{"operation": Var("pending"), "label": String("ok")})}, []Step{Append("handles", map[string]Expr{"operation": String("bad"), "label": String("bad")})}), ForEach(Var("handles"), map[string]string{"operation": "saved"}, await), Success()}
	_, err = Compile(Definition{Inputs: schema.Record{"choice": {Kind: "string", MaxBytes: 8}}, Invocations: map[string]Template{"child": leaf}, Collections: map[string]schema.Field{"handles": handles}, Steps: invalid})
	if err == nil {
		t.Fatal("branch admitted untrusted row")
	}
}

func TestInvocationHandleProvenanceNestedLists(t *testing.T) {
	leaf := invocationTemplate(t)
	inner := schema.Field{Kind: "list", MaxItems: 2, MaxBytes: 8192, Element: schema.Record{"operation": {Kind: "string", MaxBytes: 4096}}}
	outer := schema.Field{Kind: "list", MaxItems: 1, MaxBytes: 16384, Element: schema.Record{"children": inner}}
	steps := []Step{Invoke("pending", "child", nil, nil), Append("inner", map[string]Expr{"operation": Var("pending")}), Append("outer", map[string]Expr{"children": Var("inner")}), ForEach(Var("outer"), map[string]string{"children": "selected"}, ForEach(Var("selected"), map[string]string{"operation": "saved"}, AwaitInvocation("child", "saved", InvocationOutcomeBindings{Status: "status", Reason: "reason", ResponsePresent: "present"}))), Success()}
	_, err := Compile(Definition{Invocations: map[string]Template{"child": leaf}, Collections: map[string]schema.Field{"inner": inner, "outer": outer}, Steps: steps})
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvocationHandleProvenanceSeparatesSnapshotNames(t *testing.T) {
	leaf := invocationTemplate(t)
	handles := schema.Field{Kind: "list", MaxItems: 1, MaxBytes: 8192, Element: schema.Record{"operation": {Kind: "string", MaxBytes: 4096}}}
	for _, tc := range []struct {
		name       string
		collection string
		genuine    bool
	}{
		{"ordinary forged", "untrusted", false},
		{"adjacent NUL forged", "\x00loop4", false},
		{"colliding forged", "\x00loop5", false},
		{"colliding genuine", "\x00loop5", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			operation := String("forged")
			if tc.genuine {
				operation = Var("pending")
			}
			// Two collection initializers and three steps put the outer loop
			// at PC 5. Its snapshot must not alias this legal user identifier.
			artifact, err := Compile(Definition{
				Invocations: map[string]Template{"child": leaf},
				Collections: map[string]schema.Field{"genuine": handles, tc.collection: handles},
				Steps: []Step{
					Invoke("pending", "child", nil, nil),
					Append("genuine", map[string]Expr{"operation": Var("pending")}),
					Append(tc.collection, map[string]Expr{"operation": operation}),
					ForEach(Var("genuine"), nil,
						ForEach(Var(tc.collection), map[string]string{"operation": "saved"},
							AwaitInvocation("child", "saved", InvocationOutcomeBindings{Status: "status", Reason: "reason", ResponsePresent: "present"}))),
					Success(),
				},
			})
			if tc.genuine {
				if err != nil {
					t.Fatalf("genuine handle with legal identifier %q: %v", tc.collection, err)
				}
				var program executor.Program
				if err := json.Unmarshal(artifact.Bytes(), &program); err != nil {
					t.Fatal(err)
				}
				if len(program.Code) <= 5 {
					t.Fatal("collision witness no longer has an instruction at PC 5")
				}
				if loop := program.Code[5]; loop.Op != "foreach" || loop.A.Ref != "genuine" {
					t.Fatalf("collision witness moved: PC 5 is %s over %q, want foreach over genuine", loop.Op, loop.A.Ref)
				}
			} else if err == nil || !strings.Contains(err.Error(), `await requires proven request for "invocation:child"`) {
				t.Fatalf("expected forged-handle refusal for collection %q, got %v", tc.collection, err)
			}
		})
	}
}
