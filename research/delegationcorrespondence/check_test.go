package delegationcorrespondence

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func positive(t *testing.T) []byte {
	t.Helper()
	b, e := os.ReadFile("testdata/positive.json")
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func TestRawPositive(t *testing.T) {
	r := CheckRaw(positive(t))
	if !r.Valid {
		t.Fatalf("%+v", r.Failure)
	}
	if len(r.Effects) != 2 || !r.Custody.Complete {
		t.Fatalf("missing authority/custody: %+v", r)
	}
}
func TestStrictDuplicateAndExactIntegers(t *testing.T) {
	for _, s := range []string{`{"a":1,"a":2}`, `{"x":{"a":1,"a":1}}`, `{"a":1} {"b":2}`} {
		if strictJSON([]byte(s), new(any)) == nil {
			t.Fatalf("accepted %s", s)
		}
	}
	var n struct{ N int64 }
	if e := strictJSON([]byte(`{"N":9007199254740993}`), &n); e != nil || n.N != 9007199254740993 {
		t.Fatalf("%d %v", n.N, e)
	}
}
func TestRawHostileInputsPreserveResults(t *testing.T) {
	for _, name := range []string{"missing-grant", "forged-origin", "request-generation", "selected-request", "gate-attempt", "result-link", "permission", "ordinal", "extra-field"} {
		t.Run(name, func(t *testing.T) {
			var x map[string]any
			json.Unmarshal(positive(t), &x)
			tab := x["Tables"].(map[string]any)
			rows := func(s string) []any { return tab[s].([]any) }
			edit := func(table string, f func(map[string]any)) {
				r := rows(table)[0].(map[string]any)
				var d map[string]any
				json.Unmarshal([]byte(r["descriptor"].(string)), &d)
				f(d)
				b, _ := json.Marshal(d)
				r["descriptor"] = string(b)
			}
			switch name {
			case "missing-grant":
				tab["lifecycle_results"] = []any{}
			case "forged-origin":
				edit("lifecycle_results", func(d map[string]any) { d["Intent"].(map[string]any)["Run"].(map[string]any)["Subject"] = "a" })
			case "request-generation":
				edit("intents", func(d map[string]any) { d["Run"].(map[string]any)["Generation"] = float64(2) })
			case "selected-request":
				for _, r := range rows("history") {
					z := r.(map[string]any)
					if z["subject"] == "a" && z["sequence"] == float64(3) {
						z["selected_observation"] = nil
					}
				}
			case "gate-attempt", "ordinal":
				for _, r := range rows("composition_facts") {
					z := r.(map[string]any)
					if z["kind"] != "authority_admission" {
						continue
					}
					var d map[string]any
					json.Unmarshal([]byte(z["descriptor"].(string)), &d)
					g := d["AuthorityAdmission"].(map[string]any)
					if name == "ordinal" {
						g["Ordinal"] = float64(2)
					} else {
						g["AttemptDigest"] = "forged"
					}
					b, _ := json.Marshal(d)
					z["descriptor"] = string(b)
					break
				}
			case "result-link":
				edit("external_receipts", func(d map[string]any) { d["Key"] = "other-original" })
			case "permission":
				tab["permissions"] = rows("permissions")[1:]
			case "extra-field":
				edit("external_attempts", func(d map[string]any) { d["InventedAuthority"] = true })
			}
			b, _ := json.Marshal(x)
			r := CheckRaw(b)
			if r.Valid || r.Failure == nil || r.Failure.Edge == "" {
				t.Fatalf("mutant accepted: %+v", r)
			}
			if r.InputSHA256 == "" {
				t.Fatal("missing input digest")
			}
		})
	}
}
func TestGeneratedCorpus(t *testing.T) {
	c, e := GenerateCorpus()
	if e != nil {
		t.Fatal(e)
	}
	classes := map[string]bool{}
	for _, x := range c.Cases {
		classes[x.Class] = true
		if !x.Audit.AuthoritySafe && x.FirstFailedEdge < 0 {
			t.Fatalf("missing first policy edge: %s", x.ID)
		}
	}
	if len(classes) != 8 {
		t.Fatalf("classes: %v", classes)
	}
	if len(c.Prior.Histories) != 33 || len(c.Prior.Chains) != 4 {
		t.Fatal("old controls lost")
	}
	b, _ := json.Marshal(c)
	d, e := GenerateCorpus()
	if e != nil {
		t.Fatal(e)
	}
	bb, _ := json.Marshal(d)
	if string(b) != string(bb) {
		t.Fatal("nondeterministic")
	}
}
func TestFrozenCorpus(t *testing.T) {
	root := os.Getenv("ALAN_TASK3_RAW")
	if root == "" {
		t.Skip("set ALAN_TASK3_RAW to frozen extracted root")
	}
	r, e := CheckCorpus(root, root)
	if e != nil {
		t.Fatal(e)
	}
	if len(r.Cases) != 48 {
		t.Fatal(len(r.Cases))
	}
	for _, c := range r.Cases {
		invalid := strings.HasPrefix(c.Name, "mutation-") && c.Name != "mutation-base"
		if c.ContinuityFailure != nil {
			t.Errorf("continuity %s: %+v", c.Name, c.ContinuityFailure)
		}
		if !invalid && !c.External.Valid {
			t.Errorf("external %s: %+v", c.Name, c.External.Failure)
		}
		for _, p := range c.Phases {
			if !invalid && !p.Verdict.Valid {
				t.Errorf("phase %s %s: %+v", c.Name, p.Path, p.Verdict.Failure)
			}
		}
		if invalid == c.Final.Valid {
			t.Errorf("%s validity=%v failure=%+v", c.Name, c.Final.Valid, c.Final.Failure)
		}
	}
}

func TestPrimaryGatePolicyBeyondHashes(t *testing.T) {
	for _, which := range []string{"stale-epoch", "wrong-ordinal", "changed-attempt"} {
		t.Run(which, func(t *testing.T) {
			var raw RawSnapshot
			if e := strictJSON(positive(t), &raw); e != nil {
				t.Fatal(e)
			}
			for _, row := range raw.Tables["composition_facts"] {
				if textCell(row, "kind") != "authority_admission" {
					continue
				}
				f := decodeCell[Fact](row, "descriptor")
				switch which {
				case "stale-epoch":
					f.AuthorityAdmission.Source.Epoch = 0
				case "wrong-ordinal":
					f.AuthorityAdmission.Ordinal = 9
					f.ID = digest("alan/fact/authority-admission/v1", []any{raw.Namespace, f.OperationID, int64(9)})
					row["id"] = bytesJSON(f.ID)
				case "changed-attempt":
					f.AuthorityAdmission.AttemptDigest = strings.Repeat("0", 64)
				}
				f.RecordDigest = digest("alan/authority-admission/v1", f.AuthorityAdmission)
				row["descriptor"] = bytesJSON(string(bytesJSON(f)))
				break
			}
			r := CheckRaw(bytesJSON(raw))
			if r.Valid || r.Failure == nil {
				t.Fatal("self-consistent primary forgery accepted")
			}
			switch r.Failure.Edge {
			case "activation_bounds", "gate_identity", "gate_attempt_digest":
			default:
				t.Fatalf("did not reach policy edge: %+v", r.Failure)
			}
		})
	}
}
func TestFrozenFreshRevocationAndSelection(t *testing.T) {
	root := os.Getenv("ALAN_TASK3_RAW")
	if root == "" {
		t.Skip("frozen corpus")
	}
	data, e := os.ReadFile(root + "/case-index.json")
	if e != nil {
		t.Fatal(e)
	}
	var idx []IndexEntry
	if e = strictJSON(data, &idx); e != nil {
		t.Fatal(e)
	}
	for _, c := range idx {
		if c.Name != "control-wrong-prior-selection" && c.Name != "applied-unknown-revoked-before-fresh-gate" {
			continue
		}
		data, e := os.ReadFile(c.Directory + "/final/inspect.json")
		if e != nil {
			t.Fatal(e)
		}
		if c.Name == "control-wrong-prior-selection" {
			v := CheckRaw(data)
			p := ProjectRaw(data)
			if !v.Valid || p.Mismatch != "immediate_selection_narrows_model_admission" || p.Audit == nil || !p.Audit.ProtocolValid {
				t.Fatalf("lost narrowing: %+v %+v", v.Failure, p)
			}
			continue
		}
		var raw RawSnapshot
		strictJSON(data, &raw)
		changed := false
		for _, row := range raw.Tables["composition_facts"] {
			if textCell(row, "kind") != "authority_admission" {
				continue
			}
			f := decodeCell[Fact](row, "descriptor")
			if f.AuthorityAdmission.Allowed {
				continue
			}
			f.AuthorityAdmission.Allowed = true
			f.AuthorityAdmission.Reason = "authorized"
			f.RecordDigest = digest("alan/authority-admission/v1", f.AuthorityAdmission)
			row["descriptor"] = bytesJSON(string(bytesJSON(f)))
			changed = true
		}
		if !changed {
			t.Fatal("no revoked gate")
		}
		v := CheckRaw(bytesJSON(raw))
		if v.Valid || v.Failure.Edge != "fresh_gate_policy" {
			t.Fatalf("revocation mutant: %+v", v.Failure)
		}
	}
}

func TestClosedRawSchemaAndTerminalOrigins(t *testing.T) {
	for _, name := range []string{"unknown-table", "unknown-column", "duplicate-origin", "missing-terminal", "forged-terminal", "missing-terminal-route", "missing-subject"} {
		t.Run(name, func(t *testing.T) {
			var raw RawSnapshot
			strictJSON(positive(t), &raw)
			switch name {
			case "missing-subject":
				raw.Tables["subjects"] = raw.Tables["subjects"][:1]
			case "unknown-table":
				raw.Tables["invented_grants"] = []Row{}
			case "unknown-column":
				raw.Tables["metadata"][0]["invented"] = bytesJSON(true)
			case "duplicate-origin":
				raw.Tables["run_origins"][1] = raw.Tables["run_origins"][0]
			case "missing-terminal":
				raw.Tables["terminal_obligations"] = []Row{}
			case "missing-terminal-route":
				raw.Tables["terminal_obligations"][0]["destinations"] = bytesJSON("[]")
			case "forged-terminal":
				raw.Tables["terminal_obligations"][0]["source_revision"] = bytesJSON(int64(9))
			}
			r := CheckRaw(bytesJSON(raw))
			if r.Valid {
				t.Fatalf("accepted %s", name)
			}
		})
	}
}

func TestFrozenHostileExportsReachEvidenceChecks(t *testing.T) {
	root := os.Getenv("ALAN_TASK3_RAW")
	if root == "" {
		t.Skip("frozen corpus")
	}
	b, e := os.ReadFile(root + "/case-index.json")
	if e != nil {
		t.Fatal(e)
	}
	var entries []IndexEntry
	strictJSON(b, &entries)
	for _, entry := range entries {
		data, e := os.ReadFile(entry.Directory + "/raw-tables.json")
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			t.Fatal(e)
		}
		v := CheckRaw(data)
		if v.Valid || v.Failure == nil {
			t.Fatalf("%s accepted", entry.Name)
		}
		switch v.Failure.Edge {
		case "raw_table_set", "missing_table", "outside_static_domain", "schema":
			t.Fatalf("%s bypassed substantive hostile check: %+v", entry.Name, v.Failure)
		}
	}
}

func TestObservationSelectionRequiresHistoricalPredecessor(t *testing.T) {
	for _, kind := range []string{"external_result", "lifecycle"} {
		t.Run(kind, func(t *testing.T) {
			var raw RawSnapshot
			strictJSON(positive(t), &raw)
			var sourceRow, selectionRow Row
			var source Fact
			for _, r := range raw.Tables["composition_facts"] {
				f := decodeCell[Fact](r, "descriptor")
				if f.Kind == kind {
					sourceRow = r
					source = f
					break
				}
			}
			if sourceRow == nil {
				t.Fatal("source absent")
			}
			wantObservation := identity("observation", "receipt", source.OperationID)
			if kind == "lifecycle" {
				wantObservation = identity("observation", "lifecycle", source.OperationID)
			}
			var selection history
			for _, r := range raw.Tables["history"] {
				if nullableText(r, "selected_observation") == wantObservation {
					selection = history{Run: RunID{textCell(r, "subject"), intCell(r, "generation")}, Sequence: intCell(r, "sequence")}
					break
				}
			}
			for _, r := range raw.Tables["composition_facts"] {
				f := decodeCell[Fact](r, "descriptor")
				if f.Kind == "execution" && f.Run == selection.Run && f.Sequence == selection.Sequence {
					selectionRow = r
					break
				}
			}
			if selectionRow == nil {
				t.Fatal("selection absent")
			}
			f := decodeCell[Fact](selectionRow, "descriptor")
			source.Order, f.Order = f.Order, source.Order
			sourceRow["decision_order"] = bytesJSON(source.Order)
			selectionRow["decision_order"] = bytesJSON(f.Order)
			sourceRow["descriptor"] = bytesJSON(string(bytesJSON(source)))
			selectionRow["descriptor"] = bytesJSON(string(bytesJSON(f)))
			v := CheckRaw(bytesJSON(raw))
			if v.Valid {
				t.Fatalf("selected %s observation before source fact", kind)
			}
			if v.Failure.Edge != "observation_before_effect_result" && v.Failure.Edge != "observation_before_lifecycle_decision" {
				t.Fatalf("did not reach causal check: %+v", v.Failure)
			}
		})
	}
}

func TestMessageAndTerminalSelectionRequireHistoricalSource(t *testing.T) {
	for _, provenance := range []string{"ordinary-message", "runtime-terminal"} {
		t.Run(provenance, func(t *testing.T) {
			var raw RawSnapshot
			strictJSON(positive(t), &raw)
			var observationID string
			for _, r := range raw.Tables["observations"] {
				o := decodeCell[Observation](r, "descriptor")
				if provenance == "ordinary-message" && o.Obligation != "grant-ready" {
					continue
				}
				if o.Provenance == provenance && string(r["consumed_sequence"]) != "null" && intCell(r, "consumed_sequence") > 0 {
					observationID = o.ID
					break
				}
			}
			var admission Admission
			for _, r := range raw.Tables["admissions"] {
				v := decodeCell[AdmissionResult](r, "result")
				if v.ObservationID == observationID {
					admission = decodeCell[Admission](r, "descriptor")
				}
			}
			var selection history
			for _, r := range raw.Tables["history"] {
				if nullableText(r, "selected_observation") == observationID {
					selection = history{Run: RunID{textCell(r, "subject"), intCell(r, "generation")}, Sequence: intCell(r, "sequence")}
				}
			}
			var sourceRow, selectionRow Row
			for _, r := range raw.Tables["composition_facts"] {
				f := decodeCell[Fact](r, "descriptor")
				if provenance == "ordinary-message" && f.Kind == "authorization" && f.OperationID == admission.IntentID || provenance == "runtime-terminal" && f.Kind == "execution" && f.Run == admission.Source && f.Sequence == admission.TerminalSequence {
					sourceRow = r
				}
				if f.Kind == "execution" && f.Run == selection.Run && f.Sequence == selection.Sequence {
					selectionRow = r
				}
			}
			if observationID == "" || sourceRow == nil || selectionRow == nil {
				t.Fatal("missing retained selection or source")
			}
			source := decodeCell[Fact](sourceRow, "descriptor")
			selected := decodeCell[Fact](selectionRow, "descriptor")
			// Move selection before the source while preserving the source's
			// order relative to its later freeze/outcome facts.
			for _, r := range raw.Tables["composition_facts"] {
				f := decodeCell[Fact](r, "descriptor")
				if f.Order == selected.Order {
					f.Order = source.Order
				} else if f.Order >= source.Order && f.Order < selected.Order {
					f.Order++
				} else {
					continue
				}
				r["decision_order"] = bytesJSON(f.Order)
				r["descriptor"] = bytesJSON(string(bytesJSON(f)))
			}
			v := CheckRaw(bytesJSON(raw))
			want := "observation_before_message_authorization"
			if provenance == "runtime-terminal" {
				want = "observation_before_terminal_source"
			}
			if v.Valid || v.Failure.Edge != want {
				t.Fatalf("expected %s: %+v", want, v.Failure)
			}
		})
	}
}
