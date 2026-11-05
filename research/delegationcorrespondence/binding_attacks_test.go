package delegationcorrespondence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func bindingRaw(t *testing.T, name string) RawSnapshot {
	t.Helper()
	var r RawSnapshot
	if err := strictJSON(bindingFixture(t, name), &r); err != nil {
		t.Fatal(err)
	}
	return r
}
func put(r Row, k string, v any) { r[k] = bytesJSON(string(bytesJSON(v))) }
func rewriteProgram(r *RawSnapshot, subject string, edit func(*Program)) {
	replaceProgram(r, subject, edit)
	rehashFacts(r)
}
func replaceProgram(r *RawSnapshot, subject string, edit func(*Program)) {
	old := ""
	for _, row := range r.Tables["runs"] {
		if textCell(row, "subject") == subject {
			old = textCell(row, "program_id")
		}
	}
	newID := ""
	for _, row := range r.Tables["programs"] {
		if textCell(row, "id") == old {
			var p Program
			need(strictJSON(rawCell(row, "bytes"), &p) == nil, "attack_program_encoding", old, "")
			edit(&p)
			newID = sha(bytesJSON(p))
			row["id"] = bytesJSON(newID)
			put(row, "bytes", p)
		}
	}
	for _, row := range r.Tables["runs"] {
		if textCell(row, "program_id") == old {
			row["program_id"] = bytesJSON(newID)
		}
	}
	for _, row := range r.Tables["run_origins"] {
		o := decodeCell[Origin](row, "descriptor")
		if o.ProgramID == old {
			o.ProgramID = newID
			put(row, "descriptor", o)
		}
	}
	for _, row := range r.Tables["composition_deployment"] {
		d := decodeCell[Deployment](row, "descriptor")
		for j := range d.Runs {
			if d.Runs[j].ProgramID == old {
				d.Runs[j].ProgramID = newID
			}
		}
		put(row, "descriptor", d)
	}
	for _, row := range r.Tables["lifecycle_results"] {
		l := decodeCell[LifecycleResult](row, "result")
		for _, state := range []*TargetState{l.Before, l.After} {
			if state != nil && state.ProgramID == old {
				state.ProgramID = newID
			}
		}
		put(row, "result", l)
	}
}
func rehashFacts(raw *RawSnapshot) {
	c := &checker{out: &RawVerdict{}}
	c.raw = *raw
	c.initialize()
	c.index()
	for _, r := range raw.Tables["composition_facts"] {
		f := decodeCell[Fact](r, "descriptor")
		if f.Kind == "authority_admission" {
			f.RecordDigest = digest("alan/authority-admission/v1", f.AuthorityAdmission)
		} else {
			want, ok := c.expected[f.ID]
			if ok {
				f.RecordDigest = want.RecordDigest
			}
		}
		put(r, "descriptor", f)
	}
}

// falseAllow alters the complete decision mirrors at an actual refused prefix.
// It adds no result, endpoint application or selecting execution. The original
// unconsumed refusal settlement is removed so the attack remains an unknown
// allowed attempt rather than contradicting an unrelated receipt shape.
func falseAllow(raw *RawSnapshot, original string, consumed int64) {
	var a ExternalAttempt
	for _, row := range raw.Tables["external_attempts"] {
		a = decodeCell[ExternalAttempt](row, "descriptor")
		a.Allowed, a.Reason = true, "authorized"
		a.Authority.Allowed, a.Authority.Reason = true, "authorized"
		a.Authority.RequestObservationID = identity("observation", "message", original)
		a.Authority.RequestConsumedSequence = consumed
		for _, lr := range raw.Tables["lifecycle_results"] {
			l := decodeCell[LifecycleResult](lr, "result")
			if l.Request.Intent.ID == a.Authority.Use.GrantID {
				var claim Claim
				need(canonical([]byte(l.Request.Intent.Payload["claim"].Text), &claim) == nil, "attack_grant", l.Request.Intent.ID, "")
				a.Authority.OriginalRequestID = claim.OriginalRequestID
				a.Authority.ServiceOperationID = claim.ServiceOperationID
				a.Authority.GrantDecisionSequence = l.DecisionSequence
			}
		}
		put(row, "descriptor", a)
	}
	id := a.Intent.ID
	settlements := []Row{}
	for _, r := range raw.Tables["settlements"] {
		if textCell(r, "intent_id") != id {
			settlements = append(settlements, r)
		}
	}
	raw.Tables["settlements"] = settlements
	observations := []Row{}
	for _, r := range raw.Tables["observations"] {
		if decodeCell[Observation](r, "descriptor").Operation != id {
			observations = append(observations, r)
		}
	}
	raw.Tables["observations"] = observations
	admissions := []Row{}
	for _, r := range raw.Tables["admissions"] {
		a := decodeCell[Admission](r, "descriptor")
		if a.Kind != "receipt" || a.IntentID != id {
			admissions = append(admissions, r)
		}
	}
	raw.Tables["admissions"] = admissions
	order := int64(0)
	maxOrder := int64(0)
	for _, r := range raw.Tables["composition_facts"] {
		f := decodeCell[Fact](r, "descriptor")
		if f.Kind == "reservation" && f.OperationID == id {
			order = f.Order
		}
		if f.Order > maxOrder {
			maxOrder = f.Order
		}
	}
	for _, r := range raw.Tables["composition_reservations"] {
		if textCell(r, "operation_id") == id {
			v := decodeCell[Reservation](r, "descriptor")
			v.Allowed, v.Reason, v.AttemptDigest = true, "authorized", digest("alan/attempt-row/v1", a)
			put(r, "descriptor", v)
		}
	}
	root := a.Authority.ServingRootID
	ch := Charge{RootID: root, OperationID: id, Kind: "effect", Effects: 1, DecisionOrder: order}
	raw.Tables["composition_charges"] = append(raw.Tables["composition_charges"], Row{"decision_order": bytesJSON(order), "descriptor": bytesJSON(string(bytesJSON(ch))), "kind": bytesJSON("effect"), "operation_id": bytesJSON(id), "root_id": bytesJSON(root)})
	source := Source{Run: a.Intent.Run, Revision: 1}
	for _, r := range raw.Tables["runs"] {
		if textCell(r, "subject") == source.Run.Subject {
			source.Epoch = intCell(r, "epoch")
			source.Sequence = intCell(r, "sequence")
		}
	}
	g := &Gate{Version: bindingVersion, Namespace: raw.Namespace, Ordinal: 1, Source: source, AttemptDigest: digest("alan/attempt-row/v1", a), Allowed: true, Reason: "authorized"}
	roots := []string{root, a.Authority.RootID}
	if roots[1] < roots[0] {
		roots[0], roots[1] = roots[1], roots[0]
	}
	f := Fact{AuthorityAdmission: g, Version: "alan/composition-fact/v1", ID: digest("alan/fact/authority-admission/v1", []any{raw.Namespace, id, int64(1)}), Kind: "authority_admission", Order: maxOrder + 1, RootIDs: roots, Run: source.Run, Sequence: source.Sequence, OperationID: id, RecordDigest: digest("alan/authority-admission/v1", g)}
	raw.Tables["composition_facts"] = append(raw.Tables["composition_facts"], Row{"decision_order": bytesJSON(f.Order), "descriptor": bytesJSON(string(bytesJSON(f))), "id": bytesJSON(f.ID), "kind": bytesJSON(f.Kind)})
	rehashFacts(raw)
}
func denseAttack(t *testing.T, name string) (RawSnapshot, string, string) {
	t.Helper()
	switch name {
	case "A03-future-grant":
		r := bindingRaw(t, "N03")
		falseAllow(&r, identity("b", 1, 1, 0), 3)
		return r, "N03", "first_authority_policy"
	case "M01":
		r := bindingRaw(t, "N02")
		falseAllow(&r, identity("b", 1, 1, 0), 3)
		return r, "N02", "selection_before_service"
	case "M02":
		r := bindingRaw(t, "N01")
		falseAllow(&r, identity("b", 1, 6, 0), 3)
		return r, "N01", "first_authority_policy"
	case "M03":
		r := bindingRaw(t, "M03")
		for _, row := range r.Tables["admissions"] {
			a := decodeCell[Admission](row, "descriptor")
			if a.IntentID != identity("b", 1, 6, 0) {
				continue
			}
			v := decodeCell[AdmissionResult](row, "result")
			v.Accepted = true
			v.Reason = "accepted"
			v.ObservationID = identity("observation", "message", a.IntentID)
			put(row, "result", v)
			row["observation_id"] = bytesJSON(v.ObservationID)
			o := Observation{ID: v.ObservationID, Run: *v.Bound, Kind: "message", Provenance: "ordinary-message", Source: a.Source.Subject, SourceGeneration: a.Source.Generation, Obligation: a.Obligation, Operation: a.Operation, Payload: a.Payload}
			r.Tables["observations"] = append(r.Tables["observations"], Row{"consumed_sequence": json.RawMessage("null"), "descriptor": bytesJSON(string(bytesJSON(o))), "generation": bytesJSON(o.Run.Generation), "id": bytesJSON(o.ID), "subject": bytesJSON(o.Run.Subject)})
		}
		return r, "M03", "admission_binding"
	case "M04":
		r := bindingRaw(t, "N06")
		for _, row := range r.Tables["composition_facts"] {
			f := decodeCell[Fact](row, "descriptor")
			if f.Kind == "authority_admission" {
				f.AuthorityAdmission.Allowed = true
				f.AuthorityAdmission.Reason = "authorized"
				f.RecordDigest = digest("alan/authority-admission/v1", f.AuthorityAdmission)
				put(row, "descriptor", f)
			}
		}
		return r, "N06", "fresh_gate_policy"
	case "A02-filter":
		r := bindingRaw(t, "P02-v2")
		rewriteProgram(&r, "a", func(p *Program) {
			for j := range p.Code {
				in := &p.Code[j]
				if in.Match != nil && in.Match.Kind == "receipt" {
					v := stringValue("")
					in.Match.Operation = Expr{Literal: &v}
					break
				}
			}
		})
		return r, "P02-v2", "receive_filter"
	}
	t.Fatalf("unknown attack %s", name)
	return RawSnapshot{}, "", ""
}
func TestBindingDenseAttacks(t *testing.T) {
	for _, name := range []string{"M01", "M02", "M03", "M04", "A02-filter", "A03-future-grant"} {
		t.Run(name, func(t *testing.T) {
			raw, parent, edge := denseAttack(t, name)
			before := CheckBindingRaw(bindingFixture(t, parent), bindingSpec(t, parent))
			b := bytesJSON(raw)
			v := CheckBindingRaw(b, bindingSpec(t, parent))
			if output := os.Getenv("ALAN_PHASE3B4_ATTACK_OUTPUT"); output != "" {
				if err := os.MkdirAll(output, 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(output, name+".json"), b, 0644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(output, name+"-parent.json"), bindingFixture(t, parent), 0644); err != nil {
					t.Fatal(err)
				}
				summary := map[string]any{"Parent": parent, "ParentSHA256": before.InputSHA256, "AttackSHA256": sha(b), "ExpectedEdge": edge, "Verdict": v}
				if err := os.WriteFile(filepath.Join(output, name+"-verdict.json"), bytesJSON(summary), 0644); err != nil {
					t.Fatal(err)
				}
			}
			projection := projectChecked(b, v)
			if v.Valid || v.Failure == nil || v.Failure.Edge != edge || projection.Supported {
				t.Fatalf("semantic assertion: %s must reject at %s, got valid=%v failure=%+v", name, edge, v.Valid, v.Failure)
			}
		})
	}
}

func TestBindingAttackParents(t *testing.T) {
	for _, name := range []string{"P02-v2", "P10", "N01", "N02", "N06", "M03"} {
		t.Run(name, func(t *testing.T) {
			v := CheckRaw(bindingFixture(t, name))
			if !v.Valid {
				t.Fatalf("genuine parent: %+v", v.Failure)
			}
		})
	}
}

func exactCustodyAttack(t *testing.T) RawSnapshot {
	raw := bindingRaw(t, "auxiliary-serviced-equal-total")
	replyID := identity("a", 1, 16, 0)
	otherService := identity("a", 1, 13, 0)
	var reply Intent
	for _, r := range raw.Tables["intents"] {
		if textCell(r, "id") == replyID {
			reply = decodeCell[Intent](r, "descriptor")
			reply.Payload["service_operation"] = stringValue(otherService)
			reply.Payload["mutation"] = intValue(2)
			put(r, "descriptor", reply)
		}
	}
	oid := identity("observation", "message", replyID)
	selected := int64(0)
	for _, r := range raw.Tables["observations"] {
		if textCell(r, "id") == oid {
			o := decodeCell[Observation](r, "descriptor")
			o.Payload = reply.Payload
			put(r, "descriptor", o)
			selected = intCell(r, "consumed_sequence")
		}
	}
	for _, r := range raw.Tables["admissions"] {
		a := decodeCell[Admission](r, "descriptor")
		if a.IntentID == replyID {
			a.Payload = reply.Payload
			put(r, "descriptor", a)
		}
	}
	for _, r := range raw.Tables["attempts"] {
		if textCell(r, "intent_id") == replyID {
			a := decodeCell[AuthorizationRequest](r, "descriptor")
			a.Intent = reply
			put(r, "descriptor", a)
		}
	}
	update := func(p *ProjectionState) {
		p.Data["reply_mutation"] = intValue(2)
		p.Data["reply_service"] = stringValue(otherService)
	}
	for _, r := range raw.Tables["history"] {
		intents := decodeCell[[]Intent](r, "intents")
		for j := range intents {
			if intents[j].ID == replyID {
				intents[j] = reply
			}
		}
		put(r, "intents", intents)
		if textCell(r, "subject") == "b" && intCell(r, "sequence") >= selected {
			p := decodeCell[ProjectionState](r, "successor")
			update(&p)
			put(r, "successor", p)
		}
	}
	for _, r := range raw.Tables["runs"] {
		if textCell(r, "subject") == "b" {
			p := decodeCell[ProjectionState](r, "projection")
			update(&p)
			put(r, "projection", p)
		}
	}
	for _, r := range raw.Tables["terminal_obligations"] {
		if textCell(r, "subject") == "b" {
			p := decodeCell[ProjectionState](r, "outcome")
			update(&p)
			put(r, "outcome", p)
		}
	}
	rewriteProgram(&raw, "a", func(p *Program) {
		for _, key := range []string{"accepted", "mutation", "service_operation", "total"} {
			reference := key
			if key == "accepted" {
				reference = "service_accepted"
			}
			p.Code[15].Effect.Payload[key] = Expr{Ref: reference}
		}
	})
	return raw
}
func TestBindingEqualTotalCustodyAttack(t *testing.T) {
	raw := exactCustodyAttack(t)
	data := bytesJSON(raw)
	// Raw authority/reduction still accepts this coherent program and real result.
	// Only the explicitly declared principal custody association is violated.
	plain := CheckRaw(data)
	if !plain.Valid {
		t.Fatalf("unrelated raw rejection: %+v", plain.Failure)
	}
	v := CheckBindingRaw(data, bindingSpec(t, "auxiliary-serviced-equal-total"))
	if v.Valid || v.Failure == nil || v.Failure.Edge != "reply_result_custody" {
		t.Fatalf("exact result lineage assertion: valid=%v %+v", v.Valid, v.Failure)
	}
	if output := os.Getenv("ALAN_PHASE3B4_ATTACK_OUTPUT"); output != "" {
		if err := os.MkdirAll(output, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "equal-total-custody.json"), data, 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "equal-total-custody-verdict.json"), bytesJSON(v), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func rewriteAttempt(raw *RawSnapshot, edit func(*ExternalAttempt)) {
	for _, row := range raw.Tables["external_attempts"] {
		a := decodeCell[ExternalAttempt](row, "descriptor")
		if a.Intent.Run.Subject != "a" {
			continue
		}
		edit(&a)
		put(row, "descriptor", a)
		for _, r := range raw.Tables["settlements"] {
			if textCell(r, "intent_id") == a.Intent.ID {
				var s struct {
					Attempt     ExternalAttempt
					Disposition string
					Receipt     *Receipt
				}
				need(canonical(rawCell(r, "descriptor"), &s) == nil, "attack_settlement", a.Intent.ID, "")
				s.Attempt = a
				put(r, "descriptor", s)
			}
		}
		for _, r := range raw.Tables["composition_reservations"] {
			if textCell(r, "operation_id") == a.Intent.ID {
				v := decodeCell[Reservation](r, "descriptor")
				v.AttemptDigest = digest("alan/attempt-row/v1", a)
				put(r, "descriptor", v)
			}
		}
		for _, r := range raw.Tables["composition_facts"] {
			f := decodeCell[Fact](r, "descriptor")
			if f.OperationID == a.Intent.ID && f.AuthorityAdmission != nil {
				f.AuthorityAdmission.AttemptDigest = digest("alan/attempt-row/v1", a)
				put(r, "descriptor", f)
			}
		}
	}
	rehashFacts(raw)
}
func TestBindingWitnessAxes(t *testing.T) {
	attacks := []struct {
		name string
		edit func(*ExternalAttempt)
	}{
		{"A01-namespace", func(a *ExternalAttempt) { a.Authority.Namespace = "foreign" }},
		{"A01-holder-generation", func(a *ExternalAttempt) { a.Authority.Use.Root.Run.Generation = 2 }},
		{"A01-holder-revision", func(a *ExternalAttempt) { a.Authority.Use.Root.Permission.Revision = 2 }},
		{"A01-root", func(a *ExternalAttempt) { a.Authority.RootID = sha([]byte("foreign root")) }},
		{"A01-tuple-target", func(a *ExternalAttempt) { a.Authority.Use.Root.Permission.Target = "foreign" }},
		{"A01-tuple-operation", func(a *ExternalAttempt) { a.Authority.Use.Root.Permission.Operation = "foreign" }},
		{"A01-tuple-obligation", func(a *ExternalAttempt) { a.Authority.Use.Root.Permission.Obligation = "foreign" }},
		{"A02-original", func(a *ExternalAttempt) { a.Authority.OriginalRequestID = identity("b", 1, 6, 0) }},
		{"A02-observation", func(a *ExternalAttempt) {
			a.Authority.RequestObservationID = identity("observation", "message", identity("b", 1, 6, 0))
		}},
		{"A02-consumption", func(a *ExternalAttempt) { a.Authority.RequestConsumedSequence = 4 }},
		{"A02-service", func(a *ExternalAttempt) { a.Authority.ServiceOperationID = identity("a", 1, 6, 0) }},
		{"A04-forged-attempt-agreement", func(a *ExternalAttempt) { a.Authority.Use.GrantID = identity("b", 1, 5, 0) }},
		{"A05-witness-version", func(a *ExternalAttempt) { a.Authority.Version = authorityVersion }},
	}
	for _, attack := range attacks {
		t.Run(attack.name, func(t *testing.T) {
			raw := bindingRaw(t, "P10")
			rewriteAttempt(&raw, attack.edit)
			b := bytesJSON(raw)
			v := CheckBindingRaw(b, bindingSpec(t, "P10"))
			if v.Valid || v.Failure == nil || v.Failure.Edge != "first_authority_policy" {
				t.Fatalf("semantic witness assertion: %+v", v.Failure)
			}
			archiveAttack(t, attack.name, b, v)
		})
	}
}
func archiveAttack(t *testing.T, name string, b []byte, v RawVerdict) {
	t.Helper()
	if out := os.Getenv("ALAN_PHASE3B4_ATTACK_OUTPUT"); out != "" {
		if err := os.MkdirAll(out, 0755); err != nil {
			t.Fatal(err)
		}
		for suffix, data := range map[string][]byte{".json": b, "-verdict.json": bytesJSON(v)} {
			if err := os.WriteFile(filepath.Join(out, name+suffix), data, 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
}
func TestBindingAdmissionAndVersionAxes(t *testing.T) {
	for _, name := range []string{"A02-missing-admission", "A02-missing-source-authorization", "A05-profile", "A05-metadata", "A05-decision", "A05-gate", "A05-coherent-downgrade"} {
		t.Run(name, func(t *testing.T) {
			raw := bindingRaw(t, "P02-v2")
			edge := ""
			q := identity("b", 1, 1, 0)
			switch name {
			case "A02-missing-admission":
				rows := []Row{}
				for _, r := range raw.Tables["admissions"] {
					if decodeCell[Admission](r, "descriptor").IntentID != q {
						rows = append(rows, r)
					}
				}
				raw.Tables["admissions"] = rows
				acks := []Row{}
				for _, r := range raw.Tables["message_acknowledgments"] {
					if textCell(r, "intent_id") != q {
						acks = append(acks, r)
					}
				}
				raw.Tables["message_acknowledgments"] = acks
				edge = "unbacked_observation"
			case "A02-missing-source-authorization":
				rows := []Row{}
				for _, r := range raw.Tables["attempts"] {
					if textCell(r, "intent_id") != q {
						rows = append(rows, r)
					}
				}
				raw.Tables["attempts"] = rows
				edge = "message_origin"
			case "A05-profile":
				p := decodeCell[Profile](raw.Tables["composition_profile"][0], "descriptor")
				p.RequestAuthority = authorityVersion
				put(raw.Tables["composition_profile"][0], "descriptor", p)
				edge = "profile"
			case "A05-metadata":
				raw.Tables["metadata"][0]["schema_version"] = bytesJSON(13)
				edge = "profile"
			case "A05-decision":
				for _, r := range raw.Tables["lifecycle_results"] {
					l := decodeCell[LifecycleResult](r, "result")
					if l.Authority != nil {
						l.Authority.Version = authorityVersion
						put(r, "result", l)
					}
				}
				rehashFacts(&raw)
				edge = "authority_decision_identity"
			case "A05-gate":
				for _, r := range raw.Tables["composition_facts"] {
					f := decodeCell[Fact](r, "descriptor")
					if f.AuthorityAdmission != nil {
						f.AuthorityAdmission.Version = authorityVersion
						f.RecordDigest = digest("alan/authority-admission/v1", f.AuthorityAdmission)
						put(r, "descriptor", f)
						break
					}
				}
				edge = "gate_identity"
			case "A05-coherent-downgrade":
				for _, rows := range raw.Tables {
					for _, r := range rows {
						for col, b := range r {
							r[col] = json.RawMessage(strings.ReplaceAll(string(b), bindingVersion, authorityVersion))
						}
					}
				}
				raw.Tables["metadata"][0]["schema_version"] = bytesJSON(13)
				replaceProgram(&raw, "a", func(*Program) {})
				replaceProgram(&raw, "b", func(*Program) {})
				// Every descriptor changed by the downgrade retains its typed canonical
				// field order; all execution/lifecycle/attempt/reservation mirrors rehash.
				for _, row := range raw.Tables["external_attempts"] {
					a := decodeCell[ExternalAttempt](row, "descriptor")
					for _, r := range raw.Tables["composition_reservations"] {
						if textCell(r, "operation_id") == a.Intent.ID {
							v := decodeCell[Reservation](r, "descriptor")
							v.AttemptDigest = digest("alan/attempt-row/v1", a)
							put(r, "descriptor", v)
						}
					}
					for _, r := range raw.Tables["composition_facts"] {
						f := decodeCell[Fact](r, "descriptor")
						if f.OperationID == a.Intent.ID && f.AuthorityAdmission != nil {
							f.AuthorityAdmission.AttemptDigest = digest("alan/attempt-row/v1", a)
							put(r, "descriptor", f)
						}
					}
				}
				rehashFacts(&raw)
				edge = "first_authority_policy"
			}
			b := bytesJSON(raw)
			v := CheckRaw(b)
			if v.Valid || v.Failure == nil || v.Failure.Edge != edge {
				t.Fatalf("semantic axis %s expected %s: %+v", name, edge, v.Failure)
			}
			archiveAttack(t, name, b, v)
		})
	}
}

func TestBindingNonReceiveSelections(t *testing.T) {
	for _, op := range []string{"emit", "set", "branch", "end"} {
		t.Run(op, func(t *testing.T) {
			raw := bindingRaw(t, "P02-v2")
			rewriteProgram(&raw, "a", func(p *Program) {
				for j := range p.Code {
					if p.Code[j].Match != nil && p.Code[j].Match.Kind == "receipt" {
						in := Instruction{Op: op, Outcome: "succeeded", Target: j + 1, Otherwise: j + 1, Dst: "scratch", A: Expr{Literal: ptrValue(intValue(0))}, B: Expr{Literal: ptrValue(intValue(0))}}
						if op == "emit" {
							in.Effect = p.Code[0].Effect
						}
						p.Code[j] = in
						break
					}
				}
			})
			v := CheckRaw(bytesJSON(raw))
			edge := op + "_selection"
			if op == "emit" {
				edge = "emit_shape"
			}
			if v.Valid || v.Failure == nil || v.Failure.Edge != edge {
				t.Fatalf("nonreceive selection %s: %+v", op, v.Failure)
			}
			archiveAttack(t, "A02-selected-on-"+op, bytesJSON(raw), v)
		})
	}
}
func ptrValue(v Value) *Value { return &v }

func TestBindingFutureSourceAuthorization(t *testing.T) {
	raw := bindingRaw(t, "P02-v2")
	q := identity("b", 1, 1, 0)
	facts := []Fact{}
	var moved Fact
	for _, r := range raw.Tables["composition_facts"] {
		f := decodeCell[Fact](r, "descriptor")
		if f.Kind == "authorization" && f.OperationID == q {
			moved = f
		} else {
			facts = append(facts, f)
		}
	}
	sort.Slice(facts, func(i, j int) bool { return facts[i].Order < facts[j].Order })
	at := -1
	for j, f := range facts {
		if f.Kind == "execution" && f.Run == (RunID{"a", 1}) && f.Sequence == 2 {
			at = j + 1
		}
	}
	need(at > 0, "attack_parent_selection", q, "")
	facts = append(facts, Fact{})
	copy(facts[at+1:], facts[at:])
	facts[at] = moved
	orders := map[string]int64{}
	for j := range facts {
		facts[j].Order = int64(j + 1)
		orders[facts[j].ID] = facts[j].Order
	}
	for _, r := range raw.Tables["composition_facts"] {
		f := decodeCell[Fact](r, "descriptor")
		f.Order = orders[f.ID]
		r["decision_order"] = bytesJSON(f.Order)
		put(r, "descriptor", f)
	}
	for _, r := range raw.Tables["composition_charges"] {
		c := decodeCell[Charge](r, "descriptor")
		c.DecisionOrder = orders[digest("alan/fact/reservation/v1", c.OperationID)]
		r["decision_order"] = bytesJSON(c.DecisionOrder)
		put(r, "descriptor", c)
	}
	for _, r := range raw.Tables["composition_freezes"] {
		f := decodeCell[Freeze](r, "descriptor")
		f.DecisionOrder = orders[digest("alan/fact/freeze/v1", f.InvocationID)]
		put(r, "descriptor", f)
	}
	for _, r := range raw.Tables["composition_root_publications"] {
		o := decodeCell[Outcome](r, "descriptor")
		o.DecisionOrder = orders[digest("alan/fact/outcome/v1", o.InvocationID)]
		r["decision_order"] = bytesJSON(o.DecisionOrder)
		put(r, "descriptor", o)
	}
	rehashFacts(&raw)
	v := CheckRaw(bytesJSON(raw))
	if v.Valid || v.Failure == nil || v.Failure.Edge != "observation_before_message_authorization" {
		t.Fatalf("future source fact: %+v", v.Failure)
	}
	archiveAttack(t, "A03-future-source-authorization", bytesJSON(raw), v)
}

func TestBindingRecipientClaimAxes(t *testing.T) {
	for _, name := range []string{"generation", "revision", "root", "claim-version"} {
		t.Run(name, func(t *testing.T) {
			raw := bindingRaw(t, "P02-v2")
			id := identity("b", 1, 2, 0)
			var intent Intent
			var claim Claim
			for _, r := range raw.Tables["intents"] {
				if textCell(r, "id") == id {
					intent = decodeCell[Intent](r, "descriptor")
					need(canonical([]byte(intent.Payload["claim"].Text), &claim) == nil, "attack_claim", id, "")
				}
			}
			edge := "lifecycle_policy"
			switch name {
			case "generation":
				claim.Recipient.Generation = 2
			case "revision":
				claim.RecipientRevision = 2
			case "root":
				claim.RecipientRoot = RunID{"b", 1}
			case "claim-version":
				claim.Version = authorityVersion
				edge = "claim_encoding"
			}
			intent.Payload["claim"] = stringValue(string(bytesJSON(claim)))
			for _, r := range raw.Tables["intents"] {
				if textCell(r, "id") == id {
					put(r, "descriptor", intent)
				}
			}
			for _, r := range raw.Tables["history"] {
				is := decodeCell[[]Intent](r, "intents")
				for j := range is {
					if is[j].ID == id {
						is[j] = intent
					}
				}
				put(r, "intents", is)
			}
			for _, r := range raw.Tables["attempts"] {
				if textCell(r, "intent_id") == id {
					q := decodeCell[AuthorizationRequest](r, "descriptor")
					q.Intent = intent
					put(r, "descriptor", q)
				}
			}
			for _, r := range raw.Tables["lifecycle_results"] {
				if textCell(r, "id") == id {
					q := decodeCell[LifecycleRequest](r, "descriptor")
					q.Intent = intent
					put(r, "descriptor", q)
					l := decodeCell[LifecycleResult](r, "result")
					l.Request = q
					l.Authority.RecipientRootID = rootID(raw.Namespace, claim.RecipientRoot)
					put(r, "result", l)
				}
			}
			for _, r := range raw.Tables["settlements"] {
				if textCell(r, "intent_id") == id {
					var s struct {
						Request     LifecycleRequest
						Disposition string
						Accepted    bool
						Reason      string
						Successor   int64
					}
					need(canonical(rawCell(r, "descriptor"), &s) == nil, "attack_lifecycle_settlement", id, "")
					s.Request.Intent = intent
					put(r, "descriptor", s)
				}
			}
			rewriteProgram(&raw, "b", func(p *Program) { v := intent.Payload["claim"]; p.Code[1].Effect.Payload["claim"] = Expr{Literal: &v} })
			b := bytesJSON(raw)
			v := CheckRaw(b)
			if v.Valid || v.Failure == nil || v.Failure.Edge != edge {
				t.Fatalf("recipient axis %s expected %s: %+v", name, edge, v.Failure)
			}
			archiveAttack(t, "A01-recipient-"+name, b, v)
		})
	}
}
