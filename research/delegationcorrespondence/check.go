package delegationcorrespondence

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Failure struct {
	Edge, ID, Detail string
	Order            int64
}
type RowRef struct{ Table, Key, Column, SHA256 string }
type EffectReport struct {
	ID, Kind      string
	Source        RunID
	Witness       *Witness
	Allowed       bool
	Reason        string
	Order         int64
	Gates         []Gate
	Result        bool
	ResultOrder   int64
	ServingCharge int64
}
type Original struct {
	ID                                   string
	Caller, Callee                       RunID
	Committed, Admitted, Selected        bool
	RequestSelection                     int64
	EffectID, ReplyID, ResultObservation string
	ResultSelected, ReplySelected        bool
	ReplySelection                       int64
	State                                string
}
type CustodyReport struct {
	Originals                     []Original
	Complete                      bool
	BothOutstandingAtFirstService bool
	BSerialGuard                  bool
}
type RawVerdict struct {
	InputSHA256, Namespace      string
	Valid                       bool
	Failure                     *Failure
	Effects                     []EffectReport
	Custody                     CustodyReport
	FactsChecked, RowsChecked   int
	References                  []RowRef
	Limits                      []string
	Discrepancies               []string
	Revocations                 int
	RejectedAuthorityOperations []string
	Grants                      []string
	OrderedFacts                []Fact
	Programs                    map[string]string
}
type history struct {
	Run                       RunID
	Sequence, Epoch, Revision int64
	Selected                  string
	Next                      ProjectionState
	Intents                   []Intent
	Row                       Row
}
type observed struct {
	Observation
	Consumed int64
}
type effectAttempt struct {
	Intent                  Intent
	Witness                 *Witness
	Allowed                 bool
	Reason, AuthorizationID string
	Request                 WireRequest
	Prepared                *Prepared
	Raw                     json.RawMessage
}
type runRecord struct {
	ProgramID       string
	Initial         ProjectionState
	Epoch, Sequence int64
	Projection      ProjectionState
}
type checker struct {
	raw            RawSnapshot
	out            *RawVerdict
	order          int64
	members        map[RunID]string
	roots          map[string]Root
	subjects       map[RunID]Subject
	deployment     Deployment
	permissions    []Permission
	runs           map[RunID]runRecord
	programs       map[string]Program
	states         map[RunID]ProjectionState
	seq            map[RunID]int64
	intents        map[string]Intent
	revisions      map[string]int64
	history        map[string]history
	observations   map[string]observed
	admissions     map[string]Admission
	admitResults   map[string]AdmissionResult
	auth           map[string]Authorization
	life           map[string]LifecycleResult
	effects        map[string]effectAttempt
	reservations   map[string]Reservation
	charges        map[string]Charge
	receipts       map[string]Receipt
	adapterResults map[string]AdapterResult
	rowSources     map[string]json.RawMessage
	expected       map[string]Fact
	seen           map[string]bool
	revoked        map[AuthorityRoot]bool
	effectReports  map[string]int
	gateCount      map[string]int64
}

func CheckRaw(data []byte) (out RawVerdict) {
	out.InputSHA256 = sha(data)
	out.Limits = []string{"Independent finite instruction slice; no general reducer/compiler proof", "Unsigned retained rows do not prove contact completeness or exact empty-reopen timing", "Opaque adapter preparation is a trusted deployment boundary; external evidence checked separately"}
	c := &checker{out: &out}
	defer func() {
		if p := recover(); p != nil {
			if f, ok := p.(Failure); ok {
				f.Order = c.order
				out.Failure = &f
				out.Valid = false
			} else {
				panic(p)
			}
		}
	}()
	if e := strictJSON(data, &c.raw); e != nil {
		fail("raw_encoding", "", e.Error())
	}
	checkSchema(c.raw)
	c.initialize()
	c.index()
	c.validateAdmissions()
	c.checkTerminals()
	c.walk()
	c.custody()
	out.Valid = true
	return
}
func (c *checker) singleton(t string) Row {
	r := c.raw.Tables[t]
	need(len(r) == 1, "singleton", t, fmt.Sprint(len(r)))
	need(intCell(r[0], "singleton") == 1, "singleton", t, "")
	return r[0]
}
func (c *checker) initialize() {
	m := c.singleton("metadata")
	need(textCell(m, "semantics") == "alan/1d/v1", "store_semantics", "metadata", "")
	n := textCell(m, "namespace")
	need(n != "" && (c.raw.Namespace == "" || c.raw.Namespace == n), "namespace", "metadata", "")
	need(intCell(m, "schema_version") == 13, "schema", "metadata", "this corpus requires profile13; conflict14 is separately qualified")
	need(intCell(m, "key_version") == 2 && (c.raw.KeyVersion == 0 || c.raw.KeyVersion == 2), "key_version", "metadata", "")
	c.raw.Namespace = n
	c.out.Namespace = n
	c.members = map[RunID]string{}
	c.roots = map[string]Root{}
	c.subjects = map[RunID]Subject{}
	c.runs = map[RunID]runRecord{}
	c.programs = map[string]Program{}
	c.states = map[RunID]ProjectionState{}
	c.seq = map[RunID]int64{}
	c.intents = map[string]Intent{}
	c.revisions = map[string]int64{}
	c.history = map[string]history{}
	c.observations = map[string]observed{}
	c.admissions = map[string]Admission{}
	c.admitResults = map[string]AdmissionResult{}
	c.auth = map[string]Authorization{}
	c.life = map[string]LifecycleResult{}
	c.effects = map[string]effectAttempt{}
	c.reservations = map[string]Reservation{}
	c.charges = map[string]Charge{}
	c.receipts = map[string]Receipt{}
	c.adapterResults = map[string]AdapterResult{}
	c.rowSources = map[string]json.RawMessage{}
	c.expected = map[string]Fact{}
	c.seen = map[string]bool{}
	c.revoked = map[AuthorityRoot]bool{}
	c.effectReports = map[string]int{}
	c.gateCount = map[string]int64{}
	c.out.Programs = map[string]string{}
	tableNames := []string{}
	for t := range c.raw.Tables {
		tableNames = append(tableNames, t)
	}
	sort.Strings(tableNames)
	for _, t := range tableNames {
		rs := c.raw.Tables[t]
		for _, r := range rs {
			c.out.RowsChecked++
			for _, k := range []string{"descriptor", "result", "intents", "successor", "bytes"} {
				if _, ok := r[k]; ok {
					b := rawCell(r, k)
					var v any
					need(strictJSON(b, &v) == nil, "nested_json", t, "")
					c.out.References = append(c.out.References, RowRef{t, rowKey(t, r), k, sha(b)})
				}
			}
		}
	}
	sort.Slice(c.out.References, func(i, j int) bool {
		a, b := c.out.References[i], c.out.References[j]
		return a.Table+"/"+a.Key+"/"+a.Column < b.Table+"/"+b.Key+"/"+b.Column
	})
	p := decodeCell[Profile](c.singleton("composition_profile"), "descriptor")
	need(p.Version == "alan/composition/v1" && p.RequestAuthority == authorityVersion && len(p.Roots) == 2 && len(p.Templates) == 0, "profile", "composition_profile", "")
	for _, r := range p.Roots {
		need(r.ID == rootID(n, r.Run) && r.Version == "alan/resource-profile/v1" && r.Run.Generation == 1 && len(r.Members) == 1 && r.Members[0] == r.Run && len(r.Templates) == 0, "root_membership", r.ID, "")
		need(c.roots[r.ID].ID == "", "duplicate_root", r.ID, "")
		c.roots[r.ID] = r
		c.members[r.Run] = r.ID
	}
	need(len(c.raw.Tables["composition_roots"]) == 2 && len(c.raw.Tables["composition_static_members"]) == 2, "root_count", "", "")
	for _, r := range c.raw.Tables["composition_roots"] {
		d := decodeCell[Root](r, "descriptor")
		need(equal(c.roots[textCell(r, "id")], d), "root_row", d.ID, "")
	}
	for _, r := range c.raw.Tables["composition_static_members"] {
		d := decodeCell[StaticMember](r, "descriptor")
		need(d.Run == (RunID{textCell(r, "subject"), intCell(r, "generation")}) && d.RootID == textCell(r, "root_id") && c.members[d.Run] == d.RootID, "member_row", d.RootID, "")
	}
	for _, t := range []string{"composition_templates", "composition_members", "composition_invocations", "composition_admission_decisions", "composition_bindings", "composition_closures", "composition_requests", "composition_results", "composition_replies", "composition_outcome_deliveries", "composition_publications", "corpus_assets"} {
		need(len(c.raw.Tables[t]) == 0, "outside_static_domain", t, "")
	}
	c.deployment = decodeCell[Deployment](c.singleton("composition_deployment"), "descriptor")
	need(len(c.deployment.Subjects) == 2 && len(c.deployment.Runs) == 2, "deployment", "", "")
	for _, s := range c.deployment.Subjects {
		r := RunID{s.ID, s.Generation}
		need(s.Generation == 1 && s.Revision == 1 && !s.Retired && c.members[r] != "", "initial_subject", s.ID, "")
		c.subjects[r] = s
	}
	c.permissions = decodeCell[[]Permission](c.singleton("composition_permissions"), "descriptor")
	actual := map[Permission]bool{}
	for _, r := range c.raw.Tables["permissions"] {
		p := Permission{textCell(r, "kind"), textCell(r, "source"), textCell(r, "target"), textCell(r, "operation"), intCell(r, "revision"), textCell(r, "obligation")}
		need(!actual[p], "duplicate_permission", p.Source, "")
		actual[p] = true
	}
	need(len(actual) == len(c.permissions), "original_permissions", "", "")
	for _, p := range c.permissions {
		need(actual[p], "original_permissions", p.Source, fmt.Sprint(p))
	}
	for _, r := range c.raw.Tables["programs"] {
		id := textCell(r, "id")
		raw := rawCell(r, "bytes")
		need(sha(raw) == id, "program_hash", id, "")
		var p Program
		err := strictJSON(raw, &p)
		need(err == nil && (p.Semantics == "alan/1d/v1" || p.Semantics == "alan/1e/v1" && len(p.Author) > 0), "program_encoding", id, fmt.Sprint(err))
		need(textCell(r, "semantics") == p.Semantics, "program_semantics", id, "")
		c.programs[id] = p
		c.out.Programs[id] = sha(raw)
	}
	for _, r := range c.raw.Tables["runs"] {
		id := RunID{textCell(r, "subject"), intCell(r, "generation")}
		need(c.runs[id].ProgramID == "", "duplicate_run", id.Subject, "")
		c.runs[id] = runRecord{textCell(r, "program_id"), decodeCell[ProjectionState](r, "initial"), intCell(r, "epoch"), intCell(r, "sequence"), decodeCell[ProjectionState](r, "projection")}
	}
	for _, r := range c.deployment.Runs {
		cur, ok := c.runs[r.ID]
		need(ok && cur.ProgramID == r.ProgramID && equal(cur.Initial, r.Projection), "deployment_run", r.ID.Subject, "")
		_, ok = c.programs[r.ProgramID]
		need(ok, "program_missing", r.ProgramID, "")
		c.states[r.ID] = r.Projection
	}
	need(len(c.raw.Tables["run_origins"]) == 2, "origins", "", "")
	for _, r := range c.raw.Tables["run_origins"] {
		d := decodeCell[Origin](r, "descriptor")
		id := RunID{textCell(r, "subject"), intCell(r, "generation")}
		need(textCell(r, "kind") == "root" && nullableText(r, "operation_id") == "" && string(r["predecessor_generation"]) == "null" && equal(c.subjects[id], d.Subject) && c.runs[id].ProgramID == d.ProgramID && equal(c.runs[id].Initial, d.Initial), "retained_origin", id.Subject, "")
	}
	need(len(c.raw.Tables["subjects"]) == len(c.subjects), "current_subject_set", "subjects", "current rows must exactly cover original principals")
	for _, r := range c.raw.Tables["subjects"] {
		id := RunID{textCell(r, "id"), intCell(r, "generation")}
		s := c.subjects[id]
		need(s.ID != "" && s.Owner == textCell(r, "owner") && s.Revision == intCell(r, "revision") && intCell(r, "retired") == 0 && intCell(r, "retry_cancelled") == 0, "subject_mutation", id.Subject, "outside fixed generation-one authority operations")
	}
}
func historyKey(r RunID, q int64) string { return identity(r, q) }
func (c *checker) expectedFact(kind, id string, run RunID, seq int64, roots []string, raw any) Fact {
	return Fact{Version: "alan/composition-fact/v1", ID: digest("alan/fact/"+kind+"/v1", id), Kind: kind, RootIDs: roots, Run: run, Sequence: seq, OperationID: id, RecordDigest: digest("alan/"+kind+"-row/v1", raw)}
}
func (c *checker) addExpected(f Fact) {
	need(c.expected[f.ID].ID == "", "duplicate_fact_source", f.ID, "")
	c.expected[f.ID] = f
}
func (c *checker) index() {
	for _, r := range c.raw.Tables["intents"] {
		i := decodeCell[Intent](r, "descriptor")
		run, q, o := locator(i.ID)
		need(i.ID == textCell(r, "id") && i.Run == run && i.Sequence == q && i.Ordinal == o && run == (RunID{textCell(r, "subject"), intCell(r, "generation")}) && q == intCell(r, "sequence") && int64(o) == intCell(r, "ordinal"), "intent_index", i.ID, "")
		need(c.intents[i.ID].ID == "", "duplicate_intent", i.ID, "")
		validateValues(i.Payload)
		validateUse(i.Authority)
		c.intents[i.ID] = i
		c.revisions[i.ID] = intCell(r, "source_revision")
	}
	seenIntents := map[string]bool{}
	for _, r := range c.raw.Tables["history"] {
		h := history{Run: RunID{textCell(r, "subject"), intCell(r, "generation")}, Sequence: intCell(r, "sequence"), Epoch: intCell(r, "epoch"), Revision: intCell(r, "source_revision"), Selected: nullableText(r, "selected_observation"), Next: decodeCell[ProjectionState](r, "successor"), Intents: decodeCell[[]Intent](r, "intents"), Row: r}
		key := historyKey(h.Run, h.Sequence)
		need(c.history[key].Sequence == 0 && h.Sequence > 0 && h.Revision == c.subjects[h.Run].Revision, "history_index", key, "")
		c.history[key] = h
		for j, i := range h.Intents {
			need(!seenIntents[i.ID] && equal(c.intents[i.ID], i) && i.Run == h.Run && i.Sequence == h.Sequence && i.Ordinal == j && c.revisions[i.ID] == h.Revision, "history_intent", i.ID, "")
			seenIntents[i.ID] = true
		}
		evidence := struct {
			Run                       RunID
			Sequence, Epoch, Revision int64
			Selected                  string
			Successor, Intents        json.RawMessage
		}{h.Run, h.Sequence, h.Epoch, h.Revision, h.Selected, rawCell(r, "successor"), rawCell(r, "intents")}
		f := Fact{Version: "alan/composition-fact/v1", ID: digest("alan/fact/execution/v1", []any{h.Run, h.Sequence}), Kind: "execution", RootIDs: []string{c.members[h.Run]}, Run: h.Run, Sequence: h.Sequence, RecordDigest: digest("alan/history-row/v1", evidence)}
		c.addExpected(f)
	}
	need(len(seenIntents) == len(c.intents), "orphan_intent", "", "")
	for _, r := range c.raw.Tables["observations"] {
		o := decodeCell[Observation](r, "descriptor")
		n := int64(0)
		if string(r["consumed_sequence"]) != "null" {
			n = intCell(r, "consumed_sequence")
		}
		need(o.ID == textCell(r, "id") && o.Run == (RunID{textCell(r, "subject"), intCell(r, "generation")}) && c.observations[o.ID].ID == "", "observation_index", o.ID, "")
		validateValues(o.Payload)
		if n > 0 {
			need(c.history[historyKey(o.Run, n)].Selected == o.ID, "observation_consumption", o.ID, "")
		}
		c.observations[o.ID] = observed{o, n}
	}
	for _, h := range c.history {
		if h.Selected != "" {
			o, ok := c.observations[h.Selected]
			need(ok && o.Run == h.Run && o.Consumed == h.Sequence, "selected_observation", h.Selected, "")
		}
	}
	for _, r := range c.raw.Tables["attempts"] {
		q := decodeCell[AuthorizationRequest](r, "descriptor")
		a := decodeCell[Authorization](r, "result")
		id := textCell(r, "intent_id")
		need(equal(c.intents[id], q.Intent) && q.Revision == c.revisions[id] && a.ID == textCell(r, "id") && a.ID == identity("authorization", id) && a.IntentID == id && a.Source.Run == q.Intent.Run && a.Source.Revision == q.Revision && c.auth[id].ID == "", "source_authorization", id, "")
		c.auth[id] = a
		c.addExpected(c.expectedFact("authorization", id, q.Intent.Run, q.Intent.Sequence, []string{c.members[q.Intent.Run]}, rawCell(r, "result")))
	}
	for _, r := range c.raw.Tables["lifecycle_results"] {
		q := decodeCell[LifecycleRequest](r, "descriptor")
		a := decodeCell[LifecycleResult](r, "result")
		id := textCell(r, "id")
		need(equal(c.intents[id], q.Intent) && q.SourceRevision == c.revisions[id] && equal(q, a.Request) && a.DecisionSequence == intCell(r, "decision_sequence") && c.life[id].Request.Intent.ID == "", "lifecycle_origin", id, "")
		need(q.Operation == q.Intent.Payload["operation"].Text && q.ExpectedRevision == q.Intent.Payload["revision"].Number && len(q.Intent.Payload) == 3, "lifecycle_payload", id, "")
		c.life[id] = a
		roots := []string{c.members[q.Intent.Run]}
		run := RunID{q.Intent.Target, q.Intent.Generation}
		seq := int64(0)
		if a.Before != nil {
			run = a.Before.Run
			seq = a.Before.Sequence
			if root := c.members[run]; root != "" && root != roots[0] {
				roots = append(roots, root)
			}
		}
		sort.Strings(roots)
		c.addExpected(c.expectedFact("lifecycle", id, run, seq, roots, []json.RawMessage{rawCell(r, "descriptor"), rawCell(r, "result")}))
	}
	for _, t := range []string{"external_attempts", "adapter_attempts"} {
		for _, r := range c.raw.Tables[t] {
			id := textCell(r, "intent_id")
			var a effectAttempt
			if t == "external_attempts" {
				e := decodeCell[ExternalAttempt](r, "descriptor")
				a = effectAttempt{Intent: e.Intent, Witness: e.Authority, Allowed: e.Allowed, Reason: e.Reason, AuthorizationID: e.AuthorizationID, Request: e.Request}
			} else {
				e := decodeCell[AdapterAttempt](r, "descriptor")
				a = effectAttempt{Intent: e.Intent, Witness: e.Authority, Allowed: e.Allowed, Reason: e.Reason, AuthorizationID: e.AuthorizationID, Prepared: &e.Prepared}
			}
			a.Raw = rawCell(r, "descriptor")
			need(equal(a.Intent, c.intents[id]) && a.AuthorizationID == identity("authorization", id) && a.Witness != nil && c.effects[id].Intent.ID == "", "effect_original", id, "")
			c.effects[id] = a
			c.validateAttemptWire(id, a)
		}
	}
	for _, r := range c.raw.Tables["composition_reservations"] {
		a := decodeCell[Reservation](r, "descriptor")
		id := textCell(r, "operation_id")
		e, ok := c.effects[id]
		need(ok && a.OperationID == id && a.RootID == textCell(r, "root_id") && a.RootID == c.members[e.Intent.Run] && a.Run == e.Intent.Run && a.Sequence == e.Intent.Sequence && a.Kind == e.Intent.Kind && a.Allowed == e.Allowed && a.Reason == e.Reason && a.AttemptDigest == digest("alan/attempt-row/v1", e.Raw), "reservation", id, "")
		c.reservations[id] = a
		c.addExpected(c.expectedFact("reservation", id, a.Run, a.Sequence, []string{a.RootID}, rawCell(r, "descriptor")))
	}
	need(len(c.reservations) == len(c.effects), "missing_reservation", "", "")
	for _, r := range c.raw.Tables["composition_charges"] {
		a := decodeCell[Charge](r, "descriptor")
		id := textCell(r, "operation_id")
		need(a.OperationID == id && a.RootID == textCell(r, "root_id") && a.Kind == textCell(r, "kind") && a.Kind == "effect" && a.Effects == 1 && a.Invocations == 0 && a.CreatedRuns == 0 && a.DecisionOrder == intCell(r, "decision_order") && c.charges[id].OperationID == "", "charge", id, "")
		c.charges[id] = a
	}
	for _, t := range []string{"external_receipts", "adapter_results"} {
		for _, r := range c.raw.Tables[t] {
			id := textCell(r, "intent_id")
			e, ok := c.effects[id]
			need(ok && e.Allowed, "result_attempt", id, "")
			if t == "external_receipts" {
				v := decodeCell[Receipt](r, "descriptor")
				need(v.Key == e.Request.Key && v.Profile == e.Request.Profile && v.Delta == e.Request.Delta && v.Mutation > 0, "result_original", id, "")
				c.receipts[id] = v
			} else {
				v := decodeCell[AdapterResult](r, "descriptor")
				copyResult := v
				copyResult.Digest = ""
				need(v.Digest == sha(bytesJSON(copyResult)), "adapter_result_digest", id, "")
				need(e.Prepared != nil && v.Candidate.ID == id && v.Candidate.Key == e.Prepared.Request.Key, "adapter_result_original", id, "")
				c.adapterResults[id] = v
			}
			c.rowSources["result:"+id] = rawCell(r, "descriptor")
			c.addExpected(c.expectedFact("external_result", id, e.Intent.Run, e.Intent.Sequence, []string{c.members[e.Intent.Run]}, rawCell(r, "descriptor")))
		}
	}
	for _, t := range []string{"composition_freezes", "composition_root_publications"} {
		for _, r := range c.raw.Tables[t] {
			root := textCell(r, "root_id")
			scope := textCell(r, "scope_id")
			need(scope == "root:"+root && c.roots[root].ID != "", "root_scope", scope, "")
			if t == "composition_freezes" {
				f := decodeCell[Freeze](r, "descriptor")
				need(f.InvocationID == scope && f.Run == c.roots[root].Run, "freeze_origin", scope, "")
				c.addExpected(c.expectedFact("freeze", scope, f.Run, f.Sequence, []string{root}, rawCell(r, "descriptor")))
			} else {
				o := decodeCell[Outcome](r, "descriptor")
				need(o.ID == digest("alan/root-outcome/v1", root) && o.InvocationID == scope && o.OperationID == scope && o.Caller == c.roots[root].Run && o.DecisionOrder == intCell(r, "decision_order"), "outcome_origin", scope, "")
				c.addExpected(c.expectedFact("outcome", scope, o.Caller, 0, []string{root}, rawCell(r, "descriptor")))
			}
		}
	}
}
func validateValues(m map[string]Value) {
	for k, v := range m {
		switch v.Kind {
		case "int":
			need(v.Text == "" && !v.Flag, "value_encoding", k, "")
		case "string":
			need(v.Number == 0 && !v.Flag, "value_encoding", k, "")
		case "bool":
			need(v.Number == 0 && v.Text == "", "value_encoding", k, "")
		default:
			fail("value_domain", k, v.Kind)
		}
	}
}
func validateRoot(r AuthorityRoot) {
	p := r.Permission
	need(r.Run.Subject != "" && r.Run.Generation > 0 && r.RootRun.Subject != "" && r.RootRun.Generation > 0 && p.Source != "" && p.Target != "" && p.Operation != "" && p.Obligation != "" && p.Revision > 0 && (p.Kind == "external" || p.Kind == "adapter"), "root_encoding", r.Run.Subject, "")
}
func validateUse(u *Use) {
	if u == nil {
		return
	}
	validateRoot(u.Root)
	need(u.Version == authorityVersion && (u.Path == "own" && u.GrantID == "" || u.Path == "delegated" && u.GrantID != ""), "selector_encoding", u.GrantID, "")
}
func (c *checker) validateAttemptWire(id string, a effectAttempt) {
	key := identity(2, c.raw.Namespace, id)
	i := a.Intent
	if a.Prepared == nil {
		q := a.Request
		need(q.Key == key && q.Subject == i.Run.Subject && q.Authorization == a.AuthorizationID && q.Generation == i.Run.Generation && q.Sequence == i.Sequence && q.Epoch == c.auth[id].Source.Epoch && q.Profile == i.Payload["profile"].Text && q.Delta == i.Payload["delta"].Number, "external_wire", id, "")
	} else {
		p := a.Prepared
		q := p.Request
		copyPrepared := *p
		copyPrepared.Digest = ""
		need(p.Digest == sha(bytesJSON(copyPrepared)), "prepared_digest", id, "")
		need(q.ID == id && q.Key == key && q.Capability == i.Target && q.Effect == i.Obligation && equal(q.Payload, i.Payload), "adapter_prepared_origin", id, "")
	}
}
func (c *checker) committed(i Intent) bool  { return i.ID != "" && c.seq[i.Run] >= i.Sequence }
func (c *checker) has(kind, id string) bool { return c.seen[digest("alan/fact/"+kind+"/v1", id)] }
func (c *checker) walk() {
	fs := []Fact{}
	for _, r := range c.raw.Tables["composition_facts"] {
		f := decodeCell[Fact](r, "descriptor")
		need(f.ID == textCell(r, "id") && f.Kind == textCell(r, "kind") && f.Order == intCell(r, "decision_order"), "fact_index", f.ID, "")
		fs = append(fs, f)
	}
	sort.Slice(fs, func(i, j int) bool { return fs[i].Order < fs[j].Order })
	c.out.OrderedFacts = fs
	lifecycleSequence := int64(0)
	for j, f := range fs {
		c.order = f.Order
		need(f.Order == int64(j+1) && !c.seen[f.ID], "fact_order", f.ID, "")
		if f.Kind == "authority_admission" {
			c.checkGate(f)
		} else {
			e, ok := c.expected[f.ID]
			e.Order = f.Order
			need(ok && equal(e, f), "mirrored_fact", f.ID, fmt.Sprintf("kind=%s operation=%s", f.Kind, f.OperationID))
			switch f.Kind {
			case "execution":
				h := c.history[historyKey(f.Run, f.Sequence)]
				need(h.Sequence == c.seq[h.Run]+1 && h.Epoch >= 1 && h.Epoch <= c.runs[h.Run].Epoch, "history_sequence", f.ID, "")
				if h.Sequence > 1 {
					need(h.Epoch >= c.history[historyKey(h.Run, h.Sequence-1)].Epoch, "history_epoch", f.ID, "")
				}
				c.checkInstruction(h)
				c.states[h.Run] = h.Next
				c.seq[h.Run] = h.Sequence
			case "authorization":
				a := c.auth[f.OperationID]
				i := c.intents[f.OperationID]
				need(c.committed(i) && a.Source.Run == i.Run && a.Source.Sequence == c.seq[i.Run] && a.Source.Revision == c.subjects[i.Run].Revision, "source_prefix", i.ID, "")
				c.epoch(a.Source)
				need(a.Allowed && a.Reason == "authorized" && a.ClosureOperation == "", "source_policy", i.ID, "fixture source authorizations are open generation-one originals")
			case "lifecycle":
				lifecycleSequence++
				a := c.life[f.OperationID]
				need(a.DecisionSequence == lifecycleSequence, "lifecycle_sequence", f.OperationID, "")
				c.checkLifecycle(a)
			case "reservation":
				c.checkReservation(f)
			case "external_result":
				id := f.OperationID
				k, ok := c.effectReports[id]
				need(ok && len(c.out.Effects[k].Gates) > 0, "result_without_gate", id, "")
				allowed := false
				for _, g := range c.out.Effects[k].Gates {
					allowed = allowed || g.Allowed
				}
				need(allowed, "result_without_allowed_gate", id, "")
				c.out.Effects[k].Result = true
				c.out.Effects[k].ResultOrder = f.Order
			case "freeze":
				need(c.states[f.Run].Status != "open", "premature_freeze", f.ID, "")
			case "outcome":
				need(c.has("freeze", f.OperationID), "outcome_without_freeze", f.ID, "")
			default:
				fail("fact_domain", f.ID, f.Kind)
			}
		}
		c.seen[f.ID] = true
		c.out.FactsChecked++
	}
	for id := range c.expected {
		need(c.seen[id], "missing_mirrored_fact", id, "")
	}
	for r, q := range c.runs {
		need(c.seq[r] == q.Sequence && equal(c.states[r], q.Projection), "final_projection", r.Subject, "")
	}
	for id, ch := range c.charges {
		need(c.effects[id].Allowed && ch.RootID == c.members[c.effects[id].Intent.Run], "orphan_charge", id, "")
	}
}
func (c *checker) epoch(a Source) {
	run := c.runs[a.Run]
	low := int64(1)
	if a.Sequence > 0 {
		low = c.history[historyKey(a.Run, a.Sequence)].Epoch
	}
	hi := run.Epoch
	if h, ok := c.history[historyKey(a.Run, a.Sequence+1)]; ok {
		hi = h.Epoch
	}
	need(a.Epoch >= low && a.Epoch <= hi, "activation_bounds", a.Run.Subject, fmt.Sprintf("%d not [%d,%d]", a.Epoch, low, hi))
}
func (c *checker) rootReason(r AuthorityRoot, allowRevoked bool) string {
	if r.Permission.Source != r.Run.Subject {
		return "authority_holder_mismatch"
	}
	if c.members[r.Run] == "" || c.members[r.Run] != rootID(c.raw.Namespace, r.RootRun) {
		return "authority_root_mismatch"
	}
	s, ok := c.subjects[r.Run]
	if !ok || r.Run.Generation != 1 {
		return "authority_generation_changed"
	}
	if s.Revision != r.Permission.Revision {
		return "authority_revision_changed"
	}
	if s.Retired {
		return "authority_retired"
	}
	found := false
	for _, p := range c.permissions {
		found = found || p == r.Permission
	}
	if !found {
		return "authority_permission_missing"
	}
	if c.revoked[r] && !allowRevoked {
		return "authority_revoked"
	}
	return ""
}
func (c *checker) checkLifecycle(a LifecycleResult) {
	q := a.Request
	i := q.Intent
	need(c.committed(i) && c.has("authorization", i.ID) && c.auth[i.ID].Allowed, "lifecycle_before_source", i.ID, "")
	need(equal(a.Before, a.After) && a.Successor == 0, "authority_lifecycle_mutated_state", i.ID, "")
	var claim Claim
	var rev Revoke
	root := AuthorityRoot{}
	delegates := q.Operation == "DelegateAuthority"
	if delegates {
		need(canonical([]byte(i.Payload["claim"].Text), &claim) == nil && claim.Version == authorityVersion, "claim_encoding", i.ID, "")
		root = claim.Root
		_, _, o := locator(claim.ServiceOperationID)
		need(o == 0, "service_locator", i.ID, "")
	} else {
		need(q.Operation == "RevokeAuthority" && canonical([]byte(i.Payload["claim"].Text), &rev) == nil && rev.Version == authorityVersion, "revoke_encoding", i.ID, "")
		root = rev.Root
	}
	validateRoot(root)
	d := &Decision{Version: authorityVersion, Namespace: c.raw.Namespace, Root: root, RootID: rootID(c.raw.Namespace, root.RootRun)}
	if delegates {
		d.RecipientRootID = rootID(c.raw.Namespace, claim.RecipientRoot)
	}
	need(equal(a.Authority, d), "authority_decision_identity", i.ID, "")
	target := RunID{i.Target, i.Generation}
	_, targetExists := c.subjects[target]
	reason := ""
	if !targetExists {
		need(a.Before == nil, "missing_target_state", i.ID, "")
		reason = "target_missing"
	} else {
		need(a.Before != nil && a.Before.Run == target && a.Before.ProgramID == c.runs[target].ProgramID && a.Before.Owner == c.subjects[target].Owner && a.Before.Revision == c.subjects[target].Revision && a.Before.Sequence == c.seq[target] && !a.Before.Retired && !a.Before.RetryCancelled && equal(a.Before.Projection, c.states[target]), "lifecycle_target_prefix", i.ID, "")
		c.epoch(Source{target, a.Before.Epoch, a.Before.Sequence, a.Before.Revision})
	}
	if reason == "" && (i.Run != root.Run || q.SourceRevision != root.Permission.Revision) {
		reason = "authority_origin_mismatch"
	}
	if reason == "" && c.states[i.Run].Status != "open" {
		reason = "authority_source_closed"
	}
	if reason == "" {
		reason = c.rootReason(root, !delegates)
	}
	if reason == "" {
		if delegates {
			switch {
			case target != claim.Recipient || q.ExpectedRevision != claim.RecipientRevision:
				reason = "authority_recipient_mismatch"
			case c.subjects[claim.Recipient].ID == "" || c.subjects[claim.Recipient].Revision != claim.RecipientRevision || c.states[claim.Recipient].Status != "open":
				reason = "authority_recipient_closed"
			case c.members[claim.Recipient] != rootID(c.raw.Namespace, claim.RecipientRoot):
				reason = "authority_recipient_root_mismatch"
			default:
				request, ok := c.intents[claim.OriginalRequestID]
				if !ok {
					reason = "authority_request_missing"
				} else if !c.committed(request) || request.Kind != "message" || request.Run != root.Run || c.revisions[request.ID] != root.Permission.Revision || request.Target != claim.Recipient.Subject || request.Generation != claim.Recipient.Generation {
					reason = "authority_request_mismatch"
				} else {
					run, _, _ := locator(claim.ServiceOperationID)
					if run != claim.Recipient {
						reason = "authority_service_mismatch"
					}
				}
			}
		} else if target != root.Run || q.ExpectedRevision != root.Permission.Revision {
			reason = "authority_target_mismatch"
		}
	}
	accepted := reason == ""
	if accepted {
		if delegates {
			reason = "delegated"
		} else if c.revoked[root] {
			reason = "already_revoked"
		} else {
			reason = "revoked"
		}
	}
	need(a.Accepted == accepted && a.Reason == reason, "lifecycle_policy", i.ID, fmt.Sprintf("derived %v/%s retained %v/%s", accepted, reason, a.Accepted, a.Reason))
	if !accepted {
		c.out.RejectedAuthorityOperations = append(c.out.RejectedAuthorityOperations, i.ID)
	}
	if accepted {
		if delegates {
			c.out.Grants = append(c.out.Grants, i.ID)
		} else {
			c.revoked[root] = true
			c.out.Revocations++
		}
	}
}
func (c *checker) deriveWitness(i Intent) *Witness {
	w := &Witness{Version: authorityVersion, Namespace: c.raw.Namespace, Use: i.Authority, ServingRootID: c.members[i.Run]}
	deny := func(s string) *Witness { w.Reason = s; return w }
	u := i.Authority
	if u == nil {
		return deny("authority_missing")
	}
	w.RootID = rootID(c.raw.Namespace, u.Root.RootRun)
	if s := c.rootReason(u.Root, false); s != "" {
		return deny(s)
	}
	p := u.Root.Permission
	op := "apply"
	if i.Kind == "adapter" {
		op = i.Obligation
	}
	if p.Kind != i.Kind || p.Target != i.Target || p.Operation != op || p.Obligation != i.Obligation {
		return deny("authority_action_mismatch")
	}
	if u.Path == "own" {
		if u.Root.Run != i.Run || p.Revision != c.revisions[i.ID] {
			return deny("authority_owner_mismatch")
		}
	} else {
		if !c.has("lifecycle", u.GrantID) {
			return deny("authority_grant_missing")
		}
		grant := c.life[u.GrantID]
		if grant.Request.Operation != "DelegateAuthority" || !grant.Accepted || grant.Authority == nil {
			return deny("authority_grant_refused")
		}
		var d Claim
		need(canonical([]byte(grant.Request.Intent.Payload["claim"].Text), &d) == nil, "claim_encoding", u.GrantID, "")
		if d.Root != u.Root || d.Recipient != i.Run || d.RecipientRevision != c.revisions[i.ID] || c.members[i.Run] != rootID(c.raw.Namespace, d.RecipientRoot) || d.ServiceOperationID != i.ID || i.Ordinal != 0 {
			return deny("authority_service_mismatch")
		}
		w.OriginalRequestID = d.OriginalRequestID
		w.ServiceOperationID = d.ServiceOperationID
		w.GrantDecisionSequence = grant.DecisionSequence
		h := c.history[historyKey(i.Run, i.Sequence-1)]
		if i.Sequence <= 1 || h.Selected == "" {
			return deny("authority_request_not_selected")
		}
		o := c.observations[h.Selected]
		if o.Run != i.Run || o.Provenance != "ordinary-message" || o.Source != d.Root.Run.Subject || o.SourceGeneration != d.Root.Run.Generation || o.Consumed != i.Sequence-1 {
			return deny("authority_request_not_selected")
		}
		admitted := false
		for id, a := range c.admissions {
			r := c.admitResults[id]
			if a.Kind == "message" && a.IntentID == d.OriginalRequestID && r.Accepted && r.Bound != nil && *r.Bound == i.Run && r.ObservationID == h.Selected {
				admitted = true
			}
		}
		if !admitted || !c.has("authorization", d.OriginalRequestID) || !c.auth[d.OriginalRequestID].Allowed {
			return deny("authority_request_not_admitted")
		}
		w.RequestObservationID = h.Selected
		w.RequestConsumedSequence = o.Consumed
	}
	w.Allowed = true
	w.Reason = "authorized"
	return w
}
func (c *checker) checkReservation(f Fact) {
	id := f.OperationID
	a := c.effects[id]
	need(c.committed(a.Intent) && c.has("authorization", id), "effect_before_source", id, "")
	w := c.deriveWitness(a.Intent)
	need(equal(w, a.Witness), "first_authority_policy", id, fmt.Sprintf("derived %s retained %s", w.Reason, a.Witness.Reason))
	allowed := w.Allowed && c.auth[id].Allowed
	reason := w.Reason
	// The fixture has one serving-root effect limit. Count historical charges,
	// preserving permission allowance separately from reservation refusal.
	used := int64(0)
	for _, e := range c.out.Effects {
		if e.Allowed && c.members[e.Source] == c.members[a.Intent.Run] {
			used++
		}
	}
	if allowed && used >= c.roots[c.members[a.Intent.Run]].Limits.MaxEffects {
		allowed = false
		reason = "effect_limit"
	}
	need(a.Allowed == allowed && a.Reason == reason, "reservation_policy", id, fmt.Sprintf("derived %v/%s retained %v/%s", allowed, reason, a.Allowed, a.Reason))
	charge := c.charges[id]
	if allowed {
		need(charge.OperationID == id && charge.RootID == c.members[a.Intent.Run] && charge.DecisionOrder == f.Order, "serving_root_debit", id, "")
	} else {
		need(charge.OperationID == "", "refusal_debit", id, "")
	}
	c.effectReports[id] = len(c.out.Effects)
	c.out.Effects = append(c.out.Effects, EffectReport{ID: id, Kind: a.Intent.Kind, Source: a.Intent.Run, Witness: w, Allowed: allowed, Reason: reason, Order: f.Order, ServingCharge: charge.Effects})
	if reason == "authority_request_not_admitted" && w.OriginalRequestID != "" {
		for _, o := range c.observations {
			if o.ID == identity("observation", "message", w.OriginalRequestID) && o.Consumed > 0 && o.Consumed < a.Intent.Sequence-1 {
				c.out.Discrepancies = append(c.out.Discrepancies, "immediate_selection_narrows_model_admission")
			}
		}
	}
}
func (c *checker) checkGate(f Fact) {
	g := f.AuthorityAdmission
	id := f.OperationID
	a, ok := c.effects[id]
	need(ok && g != nil && a.Allowed && c.has("reservation", id), "gate_before_reservation", id, "")
	need(g.Version == authorityVersion && g.Namespace == c.raw.Namespace && g.Ordinal == c.gateCount[id]+1 && g.Source.Run == a.Intent.Run && g.Source.Sequence == c.seq[a.Intent.Run] && g.Source.Revision == c.subjects[a.Intent.Run].Revision, "gate_identity", id, "")
	c.epoch(g.Source)
	need(g.AttemptDigest == digest("alan/attempt-row/v1", a.Raw), "gate_attempt_digest", id, "")
	w := c.deriveWitness(a.Intent)
	allowed, reason := w.Allowed, w.Reason
	if c.states[a.Intent.Run].Status != "open" {
		allowed = false
		reason = "source_closed"
	}
	need(g.Allowed == allowed && g.Reason == reason, "fresh_gate_policy", id, fmt.Sprintf("derived %v/%s retained %v/%s", allowed, reason, g.Allowed, g.Reason))
	roots := []string{c.members[a.Intent.Run]}
	if a.Intent.Authority != nil && a.Intent.Authority.Path == "delegated" {
		r := c.members[a.Intent.Authority.Root.Run]
		if r != "" && r != roots[0] {
			roots = append(roots, r)
		}
	}
	sort.Strings(roots)
	expected := Fact{AuthorityAdmission: g, Version: "alan/composition-fact/v1", ID: digest("alan/fact/authority-admission/v1", []any{c.raw.Namespace, id, g.Ordinal}), Kind: "authority_admission", Order: f.Order, RootIDs: roots, Run: g.Source.Run, Sequence: g.Source.Sequence, OperationID: id, RecordDigest: digest("alan/authority-admission/v1", g)}
	need(equal(expected, f), "primary_gate_encoding", id, "")
	c.gateCount[id]++
	k := c.effectReports[id]
	c.out.Effects[k].Gates = append(c.out.Effects[k].Gates, *g)
}

// Keep names stable for reports; strings are diagnostic categories, not policy.
func serviceDirection(obligation string) string { return strings.TrimSuffix(obligation, "-request") }
