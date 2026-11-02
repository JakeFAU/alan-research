//go:build integration

package substitutiontest

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/JakeFAU/alan/author"
	"github.com/JakeFAU/alan/runtime"
	"github.com/JakeFAU/alan/schema"
)

// The two artifacts differ only in one literal under the mixed-reason guard.
// In particular the marker port, its grant, and the receiver are present in X too.
func peerSpendingProvision(t *testing.T, enableMarker bool) runtime.Provision {
	t.Helper()
	compile := func(d author.Definition) author.Artifact {
		t.Helper()
		a, err := author.Compile(d)
		must(t, err)
		return a
	}
	entry := compile(author.Definition{
		Ports: map[string]author.Port{"work": {Peer: "caller", Obligation: "work"}},
		Steps: []author.Step{author.Receive("work", nil), author.Send("work", nil), author.Success()},
	})
	members := author.Application{"entry": {Artifact: entry, Inputs: map[string]schema.Value{}}}
	one, err := author.NewTemplate(members, "entry", "work")
	must(t, err)
	members["helper"] = author.Instance{
		Artifact: compile(author.Definition{Steps: []author.Step{author.Success()}}),
		Inputs:   map[string]schema.Value{},
	}
	two, err := author.NewTemplate(members, "entry", "work")
	must(t, err)

	tag := schema.Record{"tag": {Kind: "string", MaxBytes: 16}}
	receiver := compile(author.Definition{
		Ports: map[string]author.Port{"marker": {Peer: "candidate", Obligation: "marker", Receive: tag}},
		Steps: []author.Step{author.Receive("marker", map[string]string{"tag": "seen"}), author.Success()},
	})
	peer := compile(author.Definition{
		Invocations: map[string]author.Template{"singleton": one},
		Steps:       []author.Step{author.Invoke("ticket", "singleton", nil, nil), author.Success()},
	})
	await := func(handle string) author.Step {
		return author.AwaitInvocation("pair", handle, author.InvocationOutcomeBindings{
			Status: handle + "_status", Reason: handle + "_reason", ResponsePresent: handle + "_present",
		})
	}
	mark := author.Send("marker", map[string]author.Expr{"tag": author.String("peer_spent")})
	enabled := author.IfEqual(author.Boolean(enableMarker), author.Boolean(true), []author.Step{mark}, nil)
	afterSecond := author.IfEqual(author.Var("second_reason"), author.String("invocation_limit"), []author.Step{enabled}, nil)
	candidate := compile(author.Definition{
		Ports:       map[string]author.Port{"marker": {Peer: "root", Obligation: "marker", Send: tag}},
		Invocations: map[string]author.Template{"pair": two},
		Steps: []author.Step{
			author.Invoke("first", "pair", nil, nil), await("first"),
			author.Invoke("second", "pair", nil, nil), await("second"),
			author.IfEqual(author.Var("first_reason"), author.String("created_run_limit"), []author.Step{afterSecond}, nil),
			author.Success(),
		},
	})
	p, err := (author.Application{
		"root": {Artifact: receiver}, "candidate": {Artifact: candidate}, "peer": {Artifact: peer},
	}).Provision()
	must(t, err)
	route := runtime.Permission{
		Kind: "message", Source: "candidate", Target: "root", Operation: "message", Revision: 1, Obligation: "marker",
	}
	if len(p.Permissions) != 1 || p.Permissions[0] != route {
		t.Fatal("expected the same public marker grant in both fixtures", p.Permissions)
	}
	p.Composition.Roots = []runtime.ResourceRoot{{
		Run: runtime.RunID{Subject: "root", Generation: 1},
		Members: []runtime.RunID{{Subject: "root", Generation: 1}, {Subject: "candidate", Generation: 1},
			{Subject: "peer", Generation: 1}},
		Templates: []string{one.ID(), two.ID()}, Limits: limits(1, 1, 2),
	}}
	return p
}

func TestPeerSpendingRefusalSequence(t *testing.T) {
	for _, variant := range []string{"X", "Y"} {
		provision := peerSpendingProvision(t, variant == "Y")
		for _, tc := range []struct{ order, first, second string }{
			{"before", "invocation_limit", "invocation_limit"},
			{"between", "created_run_limit", "invocation_limit"},
			{"after", "created_run_limit", "created_run_limit"},
		} {
			t.Run(variant+"/peer_"+tc.order, func(t *testing.T) {
				path := filepath.Join(t.TempDir(), "peer.db")
				store, err := runtime.Create(t.Context(), path, provision, runtime.DefaultLimits())
				must(t, err)
				t.Cleanup(func() { must(t, store.Close()) })
				w := peerSpendingRun{t: t, store: store, path: path}

				w.checkLedger()
				if tc.order == "before" {
					w.spend()
				}
				first := w.selectRefusal("first", tc.first)
				if tc.order == "between" {
					w.spend()
				}
				second := w.selectRefusal("second", tc.second)
				if tc.order == "after" {
					w.spend()
				}

				if first.Invocation.TemplateID != second.Invocation.TemplateID || first.Invocation.OperationID == second.Invocation.OperationID {
					t.Fatal("expected distinct original calls to one template", first, second)
				}
				saved, err := store.AdmitInvocation(t.Context(), first.Invocation.OperationID)
				must(t, err)
				if !reflect.DeepEqual(saved, first) {
					t.Fatal("first refusal changed after peer spending", saved, first)
				}
				w.checkLedger()
				marks := w.finishCandidate()
				want := 0
				if variant == "Y" && tc.order == "between" {
					want = 1
				}
				t.Logf("%s=%s %s=%s %s debit=(1,1,0) selected markers=%d",
					first.Invocation.OperationID, first.Reason, second.Invocation.OperationID, second.Reason, w.peerOp, marks)
				if marks != want {
					t.Fatalf("public marker count = %d, want %d", marks, want)
				}
				w.checkLedger()
			})
		}
	}
}

type peerSpendingRun struct {
	t            *testing.T
	store        *runtime.Store
	path, peerOp string
}

func (w *peerSpendingRun) commit(subject, selected string) runtime.RunState {
	w.t.Helper()
	state, err := w.store.Snapshot(w.t.Context(), runtime.RunID{Subject: subject, Generation: 1})
	must(w.t, err)
	state, err = w.store.Commit(w.t.Context(), state.Authority(), selected)
	must(w.t, err)
	return state
}

func (w *peerSpendingRun) checkLedger() {
	w.t.Helper()
	s := inspect(w.t, w.path)
	if len(s.Roots) != 1 {
		w.t.Fatal("expected one shared root", s.Roots)
	}
	root := s.Roots[0]
	rows := s.Tables["composition_charges"]
	spent := int64(0)
	if w.peerOp != "" {
		spent = 1
		if len(rows) != 1 || rows[0]["operation_id"] != w.peerOp || rows[0]["kind"] != "invocation" {
			w.t.Fatal("expected only the peer's original debit", rows)
		}
	}
	if root.LocalStatus != "open" || root.Limits != limits(1, 1, 2) || int64(len(rows)) != spent ||
		root.Charges.Invocations != spent || root.Charges.CreatedRuns != spent || root.Charges.Effects != 0 ||
		len(s.Tables["runs"]) != 3+int(spent) {
		w.t.Fatal("unexpected open-root ledger", root, len(s.Tables["runs"]))
	}
}

func (w *peerSpendingRun) spend() {
	w.t.Helper()
	w.peerOp = w.commit("peer", "").Projection.Data["ticket"].Text
	a, err := w.store.AdmitInvocation(w.t.Context(), w.peerOp)
	must(w.t, err)
	if !a.Accepted || a.Invocation.Caller.Subject != "peer" || len(a.Members) != 1 {
		w.t.Fatal("singleton admission", a)
	}
	w.checkLedger()
}

func (w *peerSpendingRun) selectRefusal(handle, reason string) runtime.InvocationAdmission {
	w.t.Helper()
	op := w.commit("candidate", "").Projection.Data[handle].Text
	a, err := w.store.AdmitInvocation(w.t.Context(), op)
	must(w.t, err)
	if a.Accepted || a.Reason != reason || len(a.Members) != 0 || a.Invocation.OperationID != op {
		w.t.Fatal("pair admission", a)
	}
	w.checkLedger()
	d, err := w.store.SettleInvocation(w.t.Context(), op)
	must(w.t, err)
	if !d.Accepted || d.Reason != "recorded_outcome" || d.ObservationID == "" {
		w.t.Fatal("refusal settlement", d)
	}
	post := w.commit("candidate", d.ObservationID)
	data := post.Projection.Data
	if data[handle].Text != op || data[handle+"_reason"] != schema.Text(reason) ||
		data[handle+"_status"] != schema.Text("rejected") || data[handle+"_present"] != schema.Bool(false) {
		w.t.Fatal("candidate did not bind this original refusal", data)
	}
	s := inspect(w.t, w.path)
	row := communicationRow(s, "observations", "id", d.ObservationID)
	o, err := readEvidence[observation](row, "descriptor")
	must(w.t, err)
	if o.Operation != op || o.Kind != "invocation" || o.Provenance != "runtime-invocation" ||
		o.Run != post.ID || o.Source != a.Invocation.TemplateID || o.SourceGeneration != 1 ||
		o.Obligation != "pair" || row["consumed_sequence"] != post.Sequence {
		w.t.Fatal("wrong selected refusal", o, row["consumed_sequence"], post.Sequence)
	}
	return a
}

// Service any actual message before the candidate's next local step, while
// both ends are open. A silent candidate follows the same driver through End.
func (w *peerSpendingRun) finishCandidate() int {
	w.t.Helper()
	marks := 0
	for range 16 { // at most three guards, their joins, a send and End
		post := w.commit("candidate", "")
		s := inspect(w.t, w.path)
		for _, row := range s.Tables["intents"] {
			i := decode[intent](row, "descriptor")
			if i.Run != post.ID || i.Sequence != post.Sequence {
				continue
			}
			if i.Kind != "message" || i.Target != "root" || i.Generation != 1 || i.Obligation != "marker" ||
				i.Payload["tag"] != schema.Text("peer_spent") || post.Projection.Status != "open" {
				w.t.Fatal("unexpected candidate suffix intent", i)
			}
			auth, err := w.store.AuthorizeIntent(w.t.Context(), i.ID)
			must(w.t, err)
			d, err := w.store.DeliverMessage(w.t.Context(), i.ID, "marker-delivery")
			must(w.t, err)
			if !auth.Allowed || !d.Accepted || d.ObservationID == "" {
				w.t.Fatal("marker route refused", auth, d)
			}
			must(w.t, w.store.AcknowledgeMessage(w.t.Context(), i.ID, "marker-delivery"))
			if w.commit("root", d.ObservationID).Projection.Data["seen"] != schema.Text("peer_spent") {
				w.t.Fatal("root failed to select the admitted marker")
			}
			marks++
		}
		if post.Projection.Status == "succeeded" {
			return marks
		}
	}
	w.t.Fatal("candidate did not reach its common successful End")
	return 0
}
