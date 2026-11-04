//go:build integration

package mutual_test

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/JakeFAU/alan-research/examples/communication/mutual"
	"github.com/JakeFAU/alan-research/runtime"
	"github.com/JakeFAU/alan-research/schema"
)

type intent struct {
	ID                       string
	Run                      runtime.RunID
	Sequence                 int64
	Ordinal                  int
	Kind, Target, Obligation string
	Generation               int64
	Payload                  map[string]schema.Value
}

type observation struct {
	ID                       string
	Run                      runtime.RunID
	Kind, Provenance, Source string
	SourceGeneration         int64
	Obligation, Operation    string
	Payload                  map[string]schema.Value
}

type commitRecord struct {
	Before, After runtime.RunState
	Selected      string
}

type refusedSelection struct {
	Authority runtime.Authority
	Selected  string
	Error     string
}

type counts struct {
	RequestsAdmitted, RequestsConsumed, ResultsAdmitted, ResultsConsumed int
	Unpaid                                                               []string
}

type evidence struct {
	Profile                  mutual.Profile
	SendOrder, DeliveryOrder []string
	Provision                runtime.Provision
	Prefix                   []commitRecord
	Refusals                 []refusedSelection
	Admissions               []runtime.AdmissionResult
	Inspections              map[string]runtime.Inspection
	Reports                  map[string]runtime.RunReport
	Replay                   runtime.ReplayReport
	Counts                   counts
}

type witness struct {
	t        *testing.T
	path     string
	store    *runtime.Store
	evidence evidence
}

func TestReciprocalRuntime(t *testing.T) {
	profiles := []mutual.Profile{mutual.SerialSerial, mutual.ServiceA, mutual.ServiceB, mutual.ServiceBoth}
	orders := [][]string{{"a", "b"}, {"b", "a"}}
	for _, profile := range profiles {
		for _, sendOrder := range orders {
			for _, deliveryOrder := range orders {
				name := fmt.Sprintf("%s/send_%s/deliver_%s", profile, sendOrder[0], deliveryOrder[0])
				t.Run(name, func(t *testing.T) {
					p, err := mutual.New(profile)
					must(t, err)
					if len(p.Runs) != 2 {
						t.Fatalf("two ordinary roots required, got %d", len(p.Runs))
					}
					d := newWitness(t, p)
					d.evidence.Profile = profile
					d.evidence.SendOrder = sendOrder
					d.evidence.DeliveryOrder = deliveryOrder
					d.prefix(sendOrder, deliveryOrder)
					before := d.inspect("admitted")
					initial := inspectCounts(t, before)
					if initial.RequestsAdmitted != 2 || initial.RequestsConsumed != 0 || len(initial.Unpaid) != 2 {
						t.Fatalf("not two admitted outstanding requests: %+v", initial)
					}
					report := d.run("first")
					after := d.inspect("after_run")
					got := inspectCounts(t, after)
					d.evidence.Counts = got
					if profile == mutual.SerialSerial {
						if !reflect.DeepEqual(got, counts{RequestsAdmitted: 2, Unpaid: got.Unpaid}) || len(got.Unpaid) != 2 {
							t.Fatalf("serial wait consumed work: %+v", got)
						}
						if report.Steps != 0 || report.Dispatches != 0 {
							t.Fatalf("serial cut enabled useful work: %+v", report)
						}
						d.checkWaitRecovery(after)
					} else {
						complete := got.RequestsConsumed == 2 && got.ResultsConsumed == 2 && len(got.Unpaid) == 0
						if !complete || got.ResultsAdmitted != 2 {
							t.Fatalf("servicing failed to consume both results: %+v", got)
						}
						for _, run := range d.states() {
							if !run.Projection.Data["validated"].Equal(schema.Bool(true)) {
								t.Fatal("reply identity/value were not validated", run)
							}
							d.requireHold(run)
						}
					}
					d.evidence.Replay, err = runtime.Replay(t.Context(), d.path, runtime.DefaultLimits())
					must(t, err)
					if !reflect.DeepEqual(d.evidence.Replay.Runs, d.states()) {
						t.Fatal("replay disagreed with retained original roots")
					}
					d.export()
				})
			}
		}
	}
}

func newWitness(t *testing.T, p runtime.Provision) *witness {
	t.Helper()
	path := filepath.Join(t.TempDir(), "mutual.db")
	s, err := runtime.Create(t.Context(), path, p, runtime.DefaultLimits())
	must(t, err)
	d := &witness{t: t, path: path, store: s, evidence: evidence{
		Provision: p, Prefix: []commitRecord{}, Admissions: []runtime.AdmissionResult{},
		Inspections: map[string]runtime.Inspection{}, Reports: map[string]runtime.RunReport{},
	}}
	t.Cleanup(func() { must(t, d.store.Close()) })
	return d
}

// Values and other compiler-generated instructions may precede the Send. Execute
// actual ordinary commits until its retained intent exists, without rewriting code.
func (d *witness) prefix(sendOrder, deliveryOrder []string) {
	t := d.t
	requests := map[string]intent{}
	for _, subject := range sendOrder {
		for range 32 {
			before, err := d.store.Snapshot(t.Context(), runID(subject))
			must(t, err)
			after, err := d.store.Commit(t.Context(), before.Authority(), "")
			must(t, err)
			d.evidence.Prefix = append(d.evidence.Prefix, commitRecord{Before: before, After: after})
			snapshot := d.inspect("committed_requests")
			for _, row := range snapshot.Tables["intents"] {
				in := decode[intent](t, row["descriptor"])
				if in.Run.Subject == subject {
					requests[subject] = in
				}
			}
			if requests[subject].ID != "" {
				break
			}
		}
		if requests[subject].ID == "" {
			t.Fatal("no initial request within the finite prefix bound", subject)
		}
	}
	if len(d.evidence.Inspections["committed_requests"].Tables["intents"]) != 2 {
		t.Fatal("prefix did not commit exactly the original two requests")
	}
	for _, subject := range deliveryOrder {
		in := requests[subject]
		transport := "phase3b-request-" + subject
		admission, err := d.store.DeliverMessage(t.Context(), in.ID, transport)
		must(t, err)
		if !admission.Accepted || admission.Bound == nil || *admission.Bound != runID(in.Target) {
			t.Fatalf("original request not admitted: %+v", admission)
		}
		d.evidence.Admissions = append(d.evidence.Admissions, admission)
		must(t, d.store.AcknowledgeMessage(t.Context(), in.ID, transport))
	}
}

func (d *witness) run(name string) runtime.RunReport {
	t := d.t
	report, err := d.store.Run(t.Context(), runtime.RunOptions{MaxSteps: 128, MaxDispatch: 32})
	must(t, err)
	if !report.Quiescent || len(report.OpenRuns) != 2 || len(report.PendingDeliveries) != 0 {
		t.Fatalf("unexpected runner disposition: %+v", report)
	}
	if len(report.UnresolvedIntents) != 0 {
		t.Fatalf("transport acknowledgments not drained: %+v", report)
	}
	d.evidence.Reports[name] = report
	return report
}

func (d *witness) inspect(name string) runtime.Inspection {
	snapshot, err := runtime.Inspect(d.t.Context(), d.path, runtime.DefaultLimits())
	must(d.t, err)
	d.evidence.Inspections[name] = snapshot
	return snapshot
}

func (d *witness) states() []runtime.RunState {
	out := []runtime.RunState{}
	for _, subject := range []string{"a", "b"} {
		run, err := d.store.Snapshot(d.t.Context(), runID(subject))
		must(d.t, err)
		out = append(out, run)
	}
	return out
}

func (d *witness) requireHold(run runtime.RunState) {
	t := d.t
	var program struct {
		Code []struct {
			Op    string
			Match struct{ Kind, Provenance string }
		}
	}
	for _, artifact := range d.evidence.Provision.Programs {
		if artifact.ID == run.ProgramID {
			must(t, json.Unmarshal(artifact.Bytes, &program))
		}
	}
	if run.Projection.Status != "open" || run.Projection.PC >= len(program.Code) {
		t.Fatal("root closed or lost its retained instruction", run)
	}
	in := program.Code[run.Projection.PC]
	if in.Op != "receive" || in.Match.Kind != "input" || in.Match.Provenance != "local-input" {
		t.Fatal("completed root did not stop at unsupplied local hold", run)
	}
}

func (d *witness) checkWaitRecovery(before runtime.Inspection) {
	t := d.t
	old := d.states()
	for _, run := range old {
		if run.Projection.Status != "open" || !run.Projection.Data["validated"].Equal(schema.Bool(false)) {
			t.Fatal("serial root did not retain an unresolved reply", run)
		}
	}
	for _, row := range before.Tables["observations"] {
		o := decode[observation](t, row["descriptor"])
		run, err := d.store.Snapshot(t.Context(), o.Run)
		must(t, err)
		_, err = d.store.Commit(t.Context(), run.Authority(), o.ID)
		if err == nil {
			t.Fatal("admitted peer request bypassed the current outbound-result receive")
		}
		// The public runtime facade does not alias ErrIneligible. Pin its current
		// public error text so unrelated storage or authority errors cannot pass.
		if err.Error() != "ineligible observation" {
			t.Fatalf("request selection failed for another reason: %v", err)
		}
		d.evidence.Refusals = append(d.evidence.Refusals, refusedSelection{
			Authority: run.Authority(), Selected: o.ID, Error: err.Error(),
		})
	}
	if !reflect.DeepEqual(before.Tables, d.inspect("after_refused_selections").Tables) {
		t.Fatal("ineligible request selection changed retained state")
	}
	report := d.run("repeated")
	if report.Steps != 0 || report.Dispatches != 0 {
		t.Fatal("repeated run advanced serial wait", report)
	}
	repeated := d.inspect("after_repeated_run")
	if !reflect.DeepEqual(before.Tables, repeated.Tables) {
		t.Fatal("repeated run changed retained wait")
	}
	must(t, d.store.Close())
	var err error
	d.store, err = runtime.Open(t.Context(), d.path, runtime.DefaultLimits())
	must(t, err)
	for i, run := range d.states() {
		if run.Epoch != old[i].Epoch+1 {
			t.Fatal("whole-Store reopen did not renew both activation epochs", old[i], run)
		}
		run.Epoch = old[i].Epoch
		if !reflect.DeepEqual(run, old[i]) {
			t.Fatal("reopen changed logical identity, program, sequence or retained control", old[i], run)
		}
	}
	reopened := d.inspect("after_reopen")
	for _, table := range []string{"programs", "subjects", "history", "intents", "admissions", "observations", "message_acknowledgments"} {
		if !reflect.DeepEqual(before.Tables[table], reopened.Tables[table]) {
			t.Fatal("reopen changed stable evidence", table)
		}
	}
	report = d.run("reopened")
	if report.Steps != 0 || report.Dispatches != 0 {
		t.Fatal("reopened run advanced serial wait", report)
	}
	final := d.inspect("after_reopened_run")
	if !reflect.DeepEqual(inspectCounts(t, before), inspectCounts(t, final)) {
		t.Fatal("recovery changed unpaid obligations")
	}
}

// Count application obligations from actual request/result and consumption
// evidence. Delivery ACKs and runner quiescence are deliberately not inputs.
func inspectCounts(t *testing.T, snapshot runtime.Inspection) counts {
	t.Helper()
	out := counts{Unpaid: []string{}}
	requests, results := map[string]int{}, map[string]int{}
	intents := map[string]intent{}
	for _, row := range snapshot.Tables["intents"] {
		in := decode[intent](t, row["descriptor"])
		intents[in.ID] = in
	}
	for _, row := range snapshot.Tables["observations"] {
		o := decode[observation](t, row["descriptor"])
		if o.Kind != "message" || o.Provenance != "ordinary-message" || o.SourceGeneration != 1 {
			t.Fatal("unexpected unauthenticated observation", o)
		}
		var origin intent
		for _, admission := range snapshot.Tables["admissions"] {
			if admission["observation_id"] == o.ID {
				descriptor := decode[struct{ IntentID string }](t, admission["descriptor"])
				origin = intents[descriptor.IntentID]
			}
		}
		sourceMatches := origin.Run.Subject == o.Source && origin.Run.Generation == o.SourceGeneration
		targetMatches := origin.Target == o.Run.Subject && origin.Generation == o.Run.Generation
		linked := origin.ID != "" && sourceMatches && targetMatches
		if !linked || origin.Obligation != o.Obligation || !reflect.DeepEqual(origin.Payload, o.Payload) {
			t.Fatal("observation not linked to exact original message", o, origin)
		}
		if o.ID != row["id"] || o.Run.Generation != 1 || o.Operation != "" {
			t.Fatal("wrong retained observation identity", row)
		}
		consumed := row["consumed_sequence"] != nil
		selected := 0
		for _, h := range snapshot.Tables["history"] {
			if h["selected_observation"] != o.ID {
				continue
			}
			sameRun := h["subject"] == o.Run.Subject && h["generation"] == o.Run.Generation
			if !consumed || h["sequence"] != row["consumed_sequence"] || !sameRun {
				t.Fatal("selected history does not certify this consumption", row, h)
			}
			selected++
		}
		if consumed && selected != 1 || !consumed && selected != 0 {
			t.Fatal("observation consumption/selected-history mismatch", row, selected)
		}
		if request, ok := o.Payload["request_id"]; ok {
			want := o.Source + "-request"
			if !request.Equal(schema.Text(want)) || o.Obligation != "phase3b_"+o.Source+"_to_"+o.Run.Subject {
				t.Fatal("wrong original request/edge", o)
			}
			out.RequestsAdmitted++
			requests[want]++
			if consumed {
				out.RequestsConsumed++
			}
			continue
		}
		want := o.Run.Subject + "-request"
		value := int64(11)
		if o.Run.Subject == "b" {
			value = 22
		}
		validID := o.Payload["response_id"].Equal(schema.Text(want))
		validValue := o.Payload["response_value"].Equal(schema.Int(value))
		if !validID || !validValue || o.Obligation != "phase3b_"+o.Run.Subject+"_to_"+o.Source {
			t.Fatal("wrong original response identity/value/edge", o)
		}
		out.ResultsAdmitted++
		if consumed {
			out.ResultsConsumed++
			results[want]++
		}
	}
	for _, subject := range []string{"a", "b"} {
		id := subject + "-request"
		if requests[id] != 1 || results[id] > 1 {
			t.Fatal("lost/duplicated original obligation", id, requests, results)
		}
		if results[id] == 0 {
			out.Unpaid = append(out.Unpaid, id)
		}
	}
	return out
}

func (d *witness) export() {
	t := d.t
	root := os.Getenv("ALAN_PHASE3B_RUNTIME_EXPORT")
	if root == "" {
		return
	}
	name := fmt.Sprintf("%s-send_%s-deliver_%s", d.evidence.Profile,
		d.evidence.SendOrder[0], d.evidence.DeliveryOrder[0])
	dir := filepath.Join(root, name)
	must(t, os.MkdirAll(dir, 0755))
	for _, program := range d.evidence.Provision.Programs {
		must(t, os.WriteFile(filepath.Join(dir, "program-"+program.ID+".json"), program.Bytes, 0644))
	}
	file, err := os.Create(filepath.Join(dir, "evidence.json.gz"))
	must(t, err)
	writer := gzip.NewWriter(file)
	must(t, json.NewEncoder(writer).Encode(d.evidence))
	must(t, writer.Close())
	must(t, file.Close())
	compact := struct {
		Profile                  mutual.Profile
		SendOrder, DeliveryOrder []string
		Counts                   counts
		Reports                  map[string]runtime.RunReport
	}{
		Profile: d.evidence.Profile, SendOrder: d.evidence.SendOrder, DeliveryOrder: d.evidence.DeliveryOrder,
		Counts: d.evidence.Counts, Reports: d.evidence.Reports,
	}
	raw, err := json.MarshalIndent(compact, "", "  ")
	must(t, err)
	must(t, os.WriteFile(filepath.Join(dir, "summary.json"), append(raw, '\n'), 0644))
}

func runID(subject string) runtime.RunID { return runtime.RunID{Subject: subject, Generation: 1} }

func decode[T any](t *testing.T, v any) T {
	t.Helper()
	var out T
	must(t, json.Unmarshal([]byte(fmt.Sprint(v)), &out))
	return out
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
