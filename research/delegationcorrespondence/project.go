package delegationcorrespondence

import (
	"fmt"
	m "github.com/JakeFAU/alan-research/research/authorityservice"
	b "github.com/JakeFAU/alan-research/research/mutualobligations"
	"sort"
)

type ProjectedEvent struct {
	Event                              m.Event
	RawFactID, RawOperationID, Mapping string
	Order                              int64
}
type Projection struct {
	ProjectionReport
	Config        m.Config
	Events        []ProjectedEvent
	Audit         *m.Audit
	Error         string
	Linearization string
	Results       map[string]b.Value
}

// ProjectRaw never contributes authority facts to CheckRaw. It only consumes a
// successful independent raw verdict and retained evidence, then asks the
// unchanged forward graph whether the derived event is enabled.
func ProjectRaw(data []byte) Projection { return projectChecked(data, CheckRaw(data)) }
func projectChecked(data []byte, v RawVerdict) (out Projection) {
	out.ProjectionReport = projectSummary(v)
	out.Linearization = "Composition facts retain their total order. Message send/admit is placed at its earlier source-authorization fact when an actual accepted admission exists; receive-result/record expands that durable admission. This is a causal witness, not a claim about unrecorded wall-clock transport order. All expansions reference the authentic operation and source fact."
	if !v.Valid {
		return out
	}
	var raw RawSnapshot
	if e := strictJSON(data, &raw); e != nil {
		out.Error = e.Error()
		out.Supported = false
		return out
	}
	intents := map[string]Intent{}
	hist := map[string]history{}
	admitted := map[string]bool{}
	replyTo := map[string]b.Obligation{}
	requestTo := map[string]b.Obligation{}
	effectTo := map[string]b.Obligation{}
	replyObservation := map[string]b.Obligation{}
	receipts := map[string]Receipt{}
	life := map[string]LifecycleResult{}
	for _, r := range raw.Tables["intents"] {
		i := decodeCell[Intent](r, "descriptor")
		intents[i.ID] = i
	}
	for _, r := range raw.Tables["history"] {
		h := history{Run: RunID{textCell(r, "subject"), intCell(r, "generation")}, Sequence: intCell(r, "sequence"), Selected: nullableText(r, "selected_observation"), Intents: decodeCell[[]Intent](r, "intents")}
		hist[historyKey(h.Run, h.Sequence)] = h
	}
	for _, o := range v.Custody.Originals {
		d := b.AB
		if o.Caller.Subject == "b" {
			d = b.BA
		}
		requestTo[o.ID] = d
		if o.EffectID != "" {
			effectTo[o.EffectID] = d
		}
		if o.ReplyID != "" {
			replyTo[o.ReplyID] = d
			replyObservation[identity("observation", "message", o.ReplyID)] = d
		}
	}
	for _, r := range raw.Tables["admissions"] {
		a := decodeCell[Admission](r, "descriptor")
		res := decodeCell[AdmissionResult](r, "result")
		if a.Kind == "message" && res.Accepted {
			admitted[a.IntentID] = true
		}
	}
	for _, r := range raw.Tables["external_receipts"] {
		receipts[textCell(r, "intent_id")] = decodeCell[Receipt](r, "descriptor")
	}
	for _, r := range raw.Tables["lifecycle_results"] {
		life[textCell(r, "id")] = decodeCell[LifecycleResult](r, "result")
	}
	out.Config = m.Config{Policy: "S2", Scenario: "LateGrant"}
	if len(v.Grants) == 0 {
		allOwn := len(v.Effects) > 0
		for _, e := range v.Effects {
			allOwn = allOwn && e.Witness.Use != nil && e.Witness.Use.Path == "own"
		}
		if allOwn {
			out.Config = m.Config{Policy: "S1", Scenario: "Local"}
		}
	}
	// Compare result content while preserving identities. Equal mutation/total
	// pairs map to the same atom, but the original-operation keys never merge.
	out.Results = map[string]b.Value{}
	values := map[string]b.Value{}
	ids := []string{}
	for id := range receipts {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		r := receipts[id]
		key := identity(r.Mutation, r.Total)
		atom, ok := values[key]
		if !ok {
			atom = b.R1
			if len(values) > 0 {
				atom = b.R2
			}
			values[key] = atom
		}
		out.Results[id] = atom
	}
	if !out.Supported && out.Mismatch == "" {
		return out
	}
	s := m.Initial(out.Config)
	events := []m.Event{}
	stopped := false
	emit := func(kind string, d b.Obligation, value b.Value, f Fact, rawID, mapping string) {
		if stopped {
			return
		}
		ts := []m.Transition{}
		for _, t := range m.Successors(out.Config, s) {
			e := t.Event
			if e.Kind == kind && e.Obligation == d && !e.Duplicate && (value == b.NoValue || e.Result == value) {
				ts = append(ts, t)
			}
		}
		if len(ts) != 1 {
			out.Error = fmt.Sprintf("raw event %s/%d at order%d has %d forward edges", kind, d, f.Order, len(ts))
			out.Supported = false
			stopped = true
			return
		}
		t := ts[0]
		s = t.State
		events = append(events, t.Event)
		out.Events = append(out.Events, ProjectedEvent{t.Event, f.ID, rawID, mapping, f.Order})
	}
	for _, f := range v.OrderedFacts {
		if stopped {
			break
		}
		switch f.Kind {
		case "execution":
			h := hist[historyKey(f.Run, f.Sequence)]
			for _, i := range h.Intents {
				if d, ok := requestTo[i.ID]; ok {
					emit("commit", d, b.NoValue, f, i.ID, "original request commitment")
				}
				if d, ok := replyTo[i.ID]; ok {
					emit("send-result", d, b.NoValue, f, i.ID, "authentic result-bearing reply emission")
				}
			}
			if d, ok := replyObservation[h.Selected]; ok {
				emit("consume", d, b.NoValue, f, h.Selected, "actual caller history selection")
			}
		case "authorization":
			id := f.OperationID
			if !admitted[id] {
				continue
			}
			if d, ok := requestTo[id]; ok {
				emit("send-request", d, b.NoValue, f, id, "accepted transport proves send; earliest causal placement")
				emit("admit", d, b.NoValue, f, id, "authenticated accepted original admission")
			}
			if d, ok := replyTo[id]; ok {
				emit("receive-result", d, b.NoValue, f, id, "accepted reply admission expansion 1/2")
				emit("record", d, b.NoValue, f, id, "durable reply observation expansion 2/2")
			}
		case "lifecycle":
			l := life[f.OperationID]
			if l.Accepted && l.Request.Operation == "DelegateAuthority" {
				emit("retain-grant", b.BA, b.NoValue, f, f.OperationID, "accepted retained lifecycle decision")
			}
		case "reservation":
			d, ok := effectTo[f.OperationID]
			if !ok {
				out.Supported = false
				out.Error = "effect lacks original mapping"
				stopped = true
				break
			}
			emit("authorize", d, b.NoValue, f, f.OperationID, "immutable nested permission decision; resource checked separately")
			if !stopped {
				actual := false
				for _, e := range v.Effects {
					if e.ID == f.OperationID {
						actual = e.Witness.Allowed
					}
				}
				if (out.Events[len(out.Events)-1].Event.Decision == m.Allowed) != actual {
					out.Mismatch = "raw_model_authorization_disagreement"
					if len(v.Discrepancies) > 0 {
						out.Mismatch = v.Discrepancies[0]
					}
					out.Supported = false
					stopped = true
				}
			}
		case "authority_admission":
			d, ok := effectTo[f.OperationID]
			if ok {
				kind := "dispatch-admit"
				if !f.AuthorityAdmission.Allowed {
					kind = "dispatch-refuse"
				}
				emit(kind, d, b.NoValue, f, f.OperationID, "actual original physical gate")
			}
		case "external_result":
			d, ok := effectTo[f.OperationID]
			if ok {
				emit("apply", d, out.Results[f.OperationID], f, f.OperationID, "atomic result abstraction only for known result with prior allowed gate")
			}
		}
	}
	audit, e := m.AuditTrace(out.Config, events)
	out.Audit = &audit
	if e != nil {
		out.Error = e.Error()
		out.Supported = false
	}
	return out
}
