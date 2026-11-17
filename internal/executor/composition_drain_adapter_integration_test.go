//go:build integration

package executor

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/JakeFAU/alan-research/effects"
	recordschema "github.com/JakeFAU/alan-research/schema"
)

func TestAcceptedChildAdapterDrainsAfterRootSuccess(t *testing.T) {
	for _, mode := range []string{"drain", "retirement", "legacy"} {
		t.Run(mode, func(t *testing.T) { testAcceptedChildAdapterAfterRootSuccess(t, mode) })
	}
}

func testAcceptedChildAdapterAfterRootSuccess(t *testing.T, mode string) {
	p, c := invocationFixture(t, "succeeded", 1)
	if mode == "retirement" {
		addInvocationController(t, &p, "RetireSubject")
	}
	template := c.Templates[0]
	var child Program
	for _, a := range p.Programs {
		if a.ID == template.Members["entry"].ProgramID {
			if err := json.Unmarshal(a.Bytes, &child); err != nil {
				t.Fatal(err)
			}
		}
	}
	input, _ := (recordschema.Record{"input": {Kind: "string", MaxBytes: 10}}).Retain()
	output, _ := (recordschema.Record{"answer": {Kind: "string", MaxBytes: 20}}).Retain()
	child.Author.Effects = map[string]EffectDeclaration{"call": {Capability: "model", Input: input, Output: output}}
	effect := Instruction{Op: "emit", Dst: "model_pending", Effect: &Effect{Kind: "adapter", Target: Lit(Text("model")), Generation: Lit(Int(1)), Obligation: Lit(Text("call")), Payload: map[string]Operand{"input": Lit(Text("question"))}}}
	child.Code = append(child.Code[:1], append([]Instruction{effect}, child.Code[1:]...)...)
	child.Author.Positions = make([]string, len(child.Code))
	raw, id, err := EncodeProgram(child)
	if err != nil {
		t.Fatal(err)
	}
	p.Programs = append(p.Programs, Artifact{ID: id, Semantics: child.Semantics, Bytes: raw})
	member := template.Members["entry"]
	member.ProgramID = id
	template.Members["entry"] = member
	template.Capabilities = []string{"model"}
	c.Roots[0].Capabilities = []string{"model"}
	replaceInvocationTemplate(t, &p, &c, template)
	path := t.TempDir() + "/drain.sqlite"
	if mode == "legacy" {
		c.Version = compositionProtocolVersion
	} else {
		c.Version = compositionDrainProtocolVersion
	}
	s, err := createStore(t.Context(), path, p, DefaultLimits(), &c)
	if err != nil {
		t.Fatal(err)
	}
	a := invokeOnce(t, s)
	entry := compositionMemberID(a.Invocation.ID, "entry")
	if _, err = s.Commit(t.Context(), mustSnapshot(t, s, entry).Authority(), compositionIdentity("alan/invocation-request/v1", a.Invocation.ID)); err != nil {
		t.Fatal(err)
	}
	emitted, err := s.Commit(t.Context(), mustSnapshot(t, s, entry).Authority(), "")
	if err != nil {
		t.Fatal(err)
	}
	op := emitted.Projection.Data["model_pending"].Text
	for range 2 {
		if _, err = s.Commit(t.Context(), current(t, s, "root").Authority(), ""); err != nil {
			t.Fatal(err)
		}
	}
	auth, err := s.AuthorizeIntent(t.Context(), op)
	if err != nil || !auth.Allowed {
		t.Fatal(auth, err)
	}
	if mode == "retirement" {
		decision := authorizedDecision(t, s, committedRequest(t, s, "controller"))
		if !decision.Accepted {
			t.Fatal(decision)
		}
	}
	w, err := s.prepareAdapterWork(t.Context(), op, map[string]effects.Adapter{"model": callbackAdapter{}}, time.Second)
	if mode != "drain" {
		if err != nil || w != nil {
			t.Fatal(w, err)
		}
		attempt, e := loadAdapterAttempt(t.Context(), s.conn, op)
		if e != nil || attempt.Allowed || attempt.Reason != "scope_closed" {
			t.Fatal(attempt, e)
		}
		if err = s.Close(); err != nil {
			t.Fatal(err)
		}
		if _, err = Replay(t.Context(), path, DefaultLimits()); err != nil {
			t.Fatal(err)
		}
		reopened, e := Open(t.Context(), path, DefaultLimits())
		if e != nil {
			t.Fatal(e)
		}
		defer reopened.Close()
		retained, e := loadAdapterAttempt(t.Context(), reopened.conn, op)
		if e != nil || retained.Allowed || retained.Reason != "scope_closed" {
			t.Fatal(retained, e)
		}
		return
	}
	if err != nil || w == nil {
		t.Fatalf("child adapter blocked: work=%v err=%v", w, err)
	}
	candidate := transportAdapter(t.Context(), *w)
	result, err := s.recordAdapterResult(t.Context(), w.attempt, candidate.candidate)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.settleAdapter(t.Context(), w.attempt, result); err != nil {
		t.Fatal(err)
	}
	report, err := s.Run(t.Context(), RunOptions{MaxSteps: 100, MaxDispatch: 100, MaxInFlight: 1, Adapters: map[string]effects.Adapter{"model": callbackAdapter{}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Roots) != 1 || report.Roots[0].Outcome == nil || report.Roots[0].Outcome.Status != "succeeded" {
		t.Fatalf("root did not settle after descendant result: %+v", report.Roots)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err = Replay(t.Context(), path, DefaultLimits()); err != nil {
		t.Fatal(err)
	}
	s, err = Open(t.Context(), path, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
}
