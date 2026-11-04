package delegationcorrespondence

import (
	"sort"
)

func stringValue(s string) Value { return Value{Kind: "string", Text: s} }
func intValue(n int64) Value     { return Value{Kind: "int", Number: n} }
func boolValue(b bool) Value     { return Value{Kind: "bool", Flag: b} }
func expr(e Expr, data map[string]Value) Value {
	need(e.Offset == nil, "instruction_domain", e.Ref, "offset expressions outside frozen fixture")
	if e.Literal != nil {
		need(e.Ref == "", "expression_encoding", e.Ref, "")
		return *e.Literal
	}
	v, ok := data[e.Ref]
	need(ok, "expression_reference", e.Ref, "")
	return v
}

// This deliberately small independent transition interpreter covers precisely
// the five instruction kinds present in exported fixtures, not general ALAN.
func (c *checker) checkInstruction(h history) {
	p := c.programs[c.runs[h.Run].ProgramID]
	prev := c.states[h.Run]
	need(prev.Status == "open" && prev.PC >= 0 && prev.PC < len(p.Code), "instruction_prefix", h.Run.Subject, "")
	in := p.Code[prev.PC]
	next := ProjectionState{PC: prev.PC + 1, Status: "open", Data: map[string]Value{}}
	for k, v := range prev.Data {
		next.Data[k] = v
	}
	emitted := []Intent{}
	switch in.Op {
	case "emit":
		need(in.Effect != nil && h.Selected == "", "emit_shape", h.Run.Subject, "")
		e := in.Effect
		i := Intent{Authority: e.Authority, ID: identity(h.Run.Subject, h.Run.Generation, h.Sequence, 0), Run: h.Run, Sequence: h.Sequence, Ordinal: 0, Kind: e.Kind, Target: expr(e.Target, prev.Data).Text, Generation: expr(e.Generation, prev.Data).Number, Obligation: expr(e.Obligation, prev.Data).Text, Payload: map[string]Value{}}
		for k, v := range e.Payload {
			i.Payload[k] = expr(v, prev.Data)
		}
		emitted = append(emitted, i)
		if in.Dst != "" {
			next.Data[in.Dst] = stringValue(i.ID)
		}
	case "receive":
		need(in.Match != nil && h.Selected != "", "receive_shape", h.Run.Subject, "")
		o := c.observations[h.Selected]
		f := in.Match
		need(o.Run == h.Run && o.Consumed == h.Sequence, "receive_original", o.ID, "")
		c.requireObservationPrefix(o)
		source := expr(f.Source, prev.Data).Text
		gen := expr(f.Generation, prev.Data).Number
		obligation := expr(f.Obligation, prev.Data).Text
		operation := expr(f.Operation, prev.Data).Text
		need((f.Kind == "" || f.Kind == o.Kind) && (f.Provenance == "" || f.Provenance == o.Provenance) && (source == "" || source == o.Source) && (gen == 0 || gen == o.SourceGeneration) && (obligation == "" || obligation == o.Obligation) && (operation == "" || operation == o.Operation), "receive_filter", o.ID, "")
		for k, dst := range in.Bind {
			v, ok := o.Payload[k]
			need(ok, "receive_payload", o.ID, k)
			next.Data[dst] = v
		}
	case "branch":
		need(h.Selected == "", "branch_selection", h.Run.Subject, "")
		next.PC = in.Otherwise
		if equal(expr(in.A, prev.Data), expr(in.B, prev.Data)) {
			next.PC = in.Target
		}
	case "set":
		need(h.Selected == "", "set_selection", h.Run.Subject, "")
		next.Data[in.Dst] = expr(in.A, prev.Data)
	case "end":
		need(h.Selected == "", "end_selection", h.Run.Subject, "")
		next.PC = prev.PC
		next.Status = in.Outcome
	default:
		fail("instruction_domain", h.Run.Subject, in.Op)
	}
	for k, e := range in.Updates {
		next.Data[k] = expr(e, prev.Data)
	}
	need(equal(emitted, h.Intents), "program_emitted_intent", historyKey(h.Run, h.Sequence), "")
	need(equal(next, h.Next), "program_successor", historyKey(h.Run, h.Sequence), in.Op)
}

// Final admission and settlement rows establish identity, not availability.
// Require each retained causal predecessor at the selection's actual prefix;
// local input has no source composition fact and remains an input assumption.
func (c *checker) requireObservationPrefix(o observed) {
	switch o.Provenance {
	case "ordinary-message", "runtime-terminal":
		for id, a := range c.admissions {
			v := c.admitResults[id]
			if !v.Accepted || v.ObservationID != o.ID {
				continue
			}
			if o.Provenance == "ordinary-message" {
				need(c.has("authorization", a.IntentID), "observation_before_message_authorization", o.ID, a.IntentID)
			} else {
				need(c.seq[a.Source] >= a.TerminalSequence, "observation_before_terminal_source", o.ID, historyKey(a.Source, a.TerminalSequence))
			}
			return
		}
		fail("observation_prefix_origin", o.ID, o.Provenance)
	case "runtime-lifecycle":
		need(c.has("lifecycle", o.Operation), "observation_before_lifecycle_decision", o.ID, o.Operation)
	case "resource-receipt", "runtime-adapter":
		need(c.has("reservation", o.Operation), "observation_before_effect_reservation", o.ID, o.Operation)
		if c.effects[o.Operation].Allowed {
			need(c.has("external_result", o.Operation), "observation_before_effect_result", o.ID, o.Operation)
		}
	case "local-input":
	default:
		fail("observation_prefix_domain", o.ID, o.Provenance)
	}
}

func (c *checker) validateAdmissions() {
	backed := map[string]bool{}
	for _, r := range c.raw.Tables["admissions"] {
		a := decodeCell[Admission](r, "descriptor")
		v := decodeCell[AdmissionResult](r, "result")
		id := textCell(r, "id")
		need(id == identity("admission", a.Kind, a.TransportID) && textCell(r, "domain") == a.Kind && textCell(r, "transport_id") == a.TransportID && v.TransportID == a.TransportID && nullableText(r, "observation_id") == v.ObservationID && c.admissions[id].TransportID == "", "admission_index", id, "")
		c.admissions[id] = a
		c.admitResults[id] = v
		if a.Kind == "message" {
			i, ok := c.intents[a.IntentID]
			need(ok && i.Kind == "message" && i.Run == a.Source && i.Target == a.Target && i.Generation == a.Generation && i.Obligation == a.Obligation && equal(i.Payload, a.Payload) && a.LogicalID == i.ID && a.TransportID == identity("message-delivery", i.ID) && a.Operation == "" && a.TerminalSequence == 0 && c.auth[i.ID].Allowed, "message_origin", a.IntentID, "")
		}
		if v.Accepted {
			if a.Kind == "message" || a.Kind == "input" {
				revision := int64(0)
				if a.Kind == "message" {
					revision = c.revisions[a.IntentID]
				}
				want := Permission{a.Kind, a.Source.Subject, a.Target, a.Kind, revision, a.Obligation}
				permitted := false
				for _, p := range c.permissions {
					permitted = permitted || p == want
				}
				need(permitted, "admission_permission", id, string(bytesJSON(want)))
			}
			need(v.Bound != nil && v.Bound.Subject == a.Target && (a.Generation == 0 || v.Bound.Generation == a.Generation) && v.ObservationID == identity("observation", a.Kind, a.LogicalID), "admission_binding", id, "")
			o := c.observations[v.ObservationID]
			provenance := map[string]string{"message": "ordinary-message", "input": "local-input", "terminal": "runtime-terminal", "lifecycle": "runtime-lifecycle", "receipt": "resource-receipt", "adapter": "runtime-adapter"}[a.Kind]
			need(provenance != "", "admission_domain", id, a.Kind)
			want := Observation{ID: v.ObservationID, Run: *v.Bound, Kind: a.Kind, Provenance: provenance, Source: a.Source.Subject, SourceGeneration: a.Source.Generation, Obligation: a.Obligation, Operation: a.Operation, Payload: a.Payload}
			need(equal(o.Observation, want), "admission_observation", id, "")
			backed[o.ID] = true
			if a.Kind == "terminal" {
				h := c.history[historyKey(a.Source, a.TerminalSequence)]
				need(h.Next.Status != "" && h.Next.Status != "open" && a.Payload["outcome"] == stringValue(h.Next.Status), "terminal_origin", id, "")
			}
		} else {
			need(v.Bound == nil || c.members[*v.Bound] != "", "refused_admission_binding", id, "")
		}
	}
	for _, r := range c.raw.Tables["message_acknowledgments"] {
		id := textCell(r, "intent_id")
		a := textCell(r, "admission_id")
		need(c.admissions[a].IntentID == id && c.admitResults[a].Accepted, "message_ack", id, "")
	}
	for _, r := range c.raw.Tables["settlements"] {
		id := textCell(r, "intent_id")
		i, ok := c.intents[id]
		need(ok, "settlement_original", id, "")
		obsID := textCell(r, "observation_id")
		o := c.observations[obsID]
		need(o.ID != "" && o.Run == i.Run && o.Operation == id && o.Source == i.Target && o.SourceGeneration == i.Generation && o.Obligation == i.Obligation, "settlement_observation", id, "")
		raw := rawCell(r, "descriptor")
		switch i.Kind {
		case "lifecycle":
			var s struct {
				Request     LifecycleRequest
				Disposition string
				Accepted    bool
				Reason      string
				Successor   int64
			}
			need(canonical(raw, &s) == nil, "settlement_encoding", id, "")
			a := c.life[id]
			need(equal(s.Request, a.Request) && s.Disposition == "target_decision" && s.Accepted == a.Accepted && s.Reason == a.Reason && s.Successor == a.Successor && textCell(r, "id") == identity("lifecycle-outcome", id), "lifecycle_settlement", id, "")
			want := map[string]Value{"accepted": boolValue(a.Accepted), "disposition": stringValue(s.Disposition), "reason": stringValue(a.Reason), "successor": intValue(a.Successor)}
			need(o.Kind == "lifecycle" && o.Provenance == "runtime-lifecycle" && equal(o.Payload, want), "lifecycle_observation", id, "")
		case "external":
			var s struct {
				Attempt     ExternalAttempt
				Disposition string
				Receipt     *Receipt
			}
			need(canonical(raw, &s) == nil, "settlement_encoding", id, "")
			need(string(bytesJSON(s.Attempt)) == string(c.effects[id].Raw) && textCell(r, "id") == identity("external-outcome", id), "external_settlement_attempt", id, "")
			a := c.effects[id]
			want := map[string]Value{"accepted": boolValue(a.Allowed), "disposition": stringValue(s.Disposition), "mutation": intValue(0), "total": intValue(0)}
			if a.Allowed {
				v, ok := c.receipts[id]
				need(ok && s.Receipt != nil && equal(*s.Receipt, v) && s.Disposition == "applied", "external_settlement_receipt", id, "")
				want["mutation"] = intValue(v.Mutation)
				want["total"] = intValue(v.Total)
			} else {
				need(s.Receipt == nil && s.Disposition == "grant_denied", "external_refusal_settlement", id, "")
			}
			need(o.Kind == "receipt" && o.Provenance == "resource-receipt" && equal(o.Payload, want), "receipt_observation", id, "")
		case "adapter":
			var s AdapterResult
			need(canonical(raw, &s) == nil, "settlement_encoding", id, "")
			v, ok := c.adapterResults[id]
			need(ok && equal(s, v) && equal(o.Payload, v.Values) && o.Kind == "adapter" && o.Provenance == "runtime-adapter", "adapter_settlement_result", id, "")
		default:
			fail("settlement_domain", id, i.Kind)
		}
		backed[o.ID] = true
	}
	for id := range c.observations {
		need(backed[id], "unbacked_observation", id, "")
	}
}
func (c *checker) custody() {
	ids := []string{}
	for id, i := range c.intents {
		if i.Kind == "message" && (i.Obligation == "ab-request" || i.Obligation == "ba-request") {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	for _, id := range ids {
		i := c.intents[id]
		o := Original{ID: id, Caller: i.Run, Callee: RunID{i.Target, i.Generation}, Committed: true, State: "committed-unpaid"}
		for key, a := range c.admissions {
			r := c.admitResults[key]
			if a.Kind == "message" && a.IntentID == id && r.Accepted {
				v := c.observations[r.ObservationID]
				o.Admitted = true
				o.Selected = v.Consumed > 0
				o.RequestSelection = v.Consumed
			}
		}
		for eid, e := range c.effects {
			if e.Intent.Run != o.Callee {
				continue
			}
			if e.Witness.OriginalRequestID == id || e.Intent.Obligation == serviceDirection(i.Obligation)+"-service" {
				o.EffectID = eid
				if !e.Allowed {
					o.State = "first-refused"
				} else if c.receipts[eid].Key == "" && c.adapterResults[eid].Digest == "" {
					o.State = "unknown-or-not-contacted"
				} else {
					o.State = "retained-result"
				}
				for oid, v := range c.observations {
					if v.Operation == eid && (v.Kind == "receipt" || v.Kind == "adapter") {
						o.ResultObservation = oid
						o.ResultSelected = v.Consumed > 0
					}
				}
			}
		}
		for rid, r := range c.intents {
			if r.Kind != "message" || r.Obligation != serviceDirection(i.Obligation)+"-reply" {
				continue
			}
			need(r.Run == o.Callee && r.Target == o.Caller.Subject && r.Generation == o.Caller.Generation && r.Payload["request_id"] == stringValue(id), "reply_original", rid, "")
			eid := r.Payload["service_operation"].Text
			e, ok := c.effects[eid]
			need(ok && eid == o.EffectID && e.Allowed && o.ResultSelected && c.observations[o.ResultObservation].Consumed < r.Sequence, "reply_result_custody", rid, "")
			v := c.receipts[eid]
			need(v.Key != "" && r.Payload["accepted"] == boolValue(true) && r.Payload["mutation"] == intValue(v.Mutation) && r.Payload["total"] == intValue(v.Total), "reply_result_values", rid, "")
			o.ReplyID = rid
			for aid, a := range c.admissions {
				res := c.admitResults[aid]
				if a.IntentID == rid && res.Accepted {
					ob := c.observations[res.ObservationID]
					need(ob.Run == o.Caller, "caller_generation", rid, "")
					o.ReplySelected = ob.Consumed > 0
					o.ReplySelection = ob.Consumed
					if o.ReplySelected {
						o.State = "consumed"
					}
				}
			}
		}
		c.out.Custody.Originals = append(c.out.Custody.Originals, o)
	}
	c.out.Custody.Complete = len(ids) == 2
	for _, o := range c.out.Custody.Originals {
		c.out.Custody.Complete = c.out.Custody.Complete && o.ReplySelected
	}
	if len(ids) == 2 && len(c.out.Effects) > 0 {
		first := c.out.Effects[0].Order
		all := true
		for _, id := range ids {
			i := c.intents[id]
			f := c.expected[digest("alan/fact/execution/v1", []any{i.Run, i.Sequence})]
			order := int64(0)
			for _, v := range c.out.OrderedFacts {
				if v.ID == f.ID {
					order = v.Order
				}
			}
			all = all && order > 0 && order < first
		}
		c.out.Custody.BothOutstandingAtFirstService = all
		for _, o := range c.out.Custody.Originals {
			if o.Caller.Subject == "b" && o.ReplySelected {
				for _, e := range c.out.Effects {
					if e.Source == o.Caller {
						c.out.Custody.BSerialGuard = o.ReplySelection < c.effects[e.ID].Intent.Sequence
					}
				}
			}
		}
	}
}
