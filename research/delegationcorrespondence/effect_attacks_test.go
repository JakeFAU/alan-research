package delegationcorrespondence

import (
	"encoding/json"
	"strings"
	"testing"
)

func effectAttackJSON(t *testing.T, parent, name string, edit func(map[string]any)) []byte {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal(readRevocationTestFile(t, parent, name), &value); err != nil {
		t.Fatal(err)
	}
	edit(value)
	return bytesJSON(value)
}

func assertEffectSemanticRejection(t *testing.T, dir, reason string) {
	t.Helper()
	r, err := CheckEffectObservationDirectory(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !r.RawValid || !r.ContinuityValid || !r.Physical.Valid || r.Supported || r.ProtocolValid || !strings.Contains(strings.Join(r.Reasons, ";"), reason) {
		t.Fatalf("semantic assertion: expected %s after raw/continuity/physical success, got raw=%v continuity=%v physical=%v protocol=%v supported=%v reasons=%v", reason, r.RawValid, r.ContinuityValid, r.Physical.Valid, r.ProtocolValid, r.Supported, r.Reasons)
	}
	t.Logf("raw/continuity/physical valid; rejected at %v", r.Reasons)
}

func TestEffectObservationAckAttacks(t *testing.T) {
	parent := effectObservationFixture(t, "normal", effectLossName)
	for _, tc := range []struct {
		name, role, field string
		value             any
		reason            string
	}{
		{"endpoint-phase", "Endpoint", "Phase", "pre-commit", "effect_endpoint_causal_ack"},
		{"endpoint-boundary", "Endpoint", "Boundary", "resource.before", "effect_endpoint_causal_ack"},
		{"endpoint-key", "Endpoint", "Intent", "wrong", "effect_endpoint_ack_binding"},
		{"endpoint-auth", "Endpoint", "Attempt", "wrong", "effect_endpoint_ack_binding"},
		{"endpoint-subject", "Endpoint", "Subject", "b", "effect_endpoint_ack_binding"},
		{"endpoint-generation", "Endpoint", "Generation", 2, "effect_endpoint_ack_binding"},
		{"endpoint-epoch", "Endpoint", "Epoch", 3, "effect_endpoint_ack_binding"},
		{"endpoint-sequence", "Endpoint", "Sequence", 5, "effect_endpoint_ack_binding"},
		{"endpoint-pid", "Endpoint", "PID", 1, "effect_endpoint_causal_ack"},
		{"endpoint-token", "Endpoint", "Token", "wrong", "effect_endpoint_causal_ack"},
		{"endpoint-store", "Endpoint", "Store", "wrong", "effect_endpoint_causal_ack"},
		{"endpoint-scenario", "Endpoint", "Scenario", "wrong", "effect_endpoint_causal_ack"},
		{"runtime-phase", "Runtime", "Phase", "pre-call", "effect_runtime_inflight_ack"},
		{"runtime-scenario", "Runtime", "Scenario", "wrong", "effect_runtime_inflight_ack"},
		{"runtime-inflight-missing", "Runtime", "InFlight", nil, "effect_runtime_inflight_ack"},
		{"runtime-inflight-other", "Runtime", "InFlight", []string{"wrong"}, "effect_runtime_inflight_ack"},
		{"runtime-token", "Runtime", "Token", "wrong", "effect_sealed_process_anchor"},
		{"runtime-store", "Runtime", "Store", "wrong", "effect_sealed_process_anchor"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := effectAttackJSON(t, parent, "acknowledgments.json", func(v map[string]any) { v[tc.role].(map[string]any)[tc.field] = tc.value })
			assertEffectSemanticRejection(t, editRevocationDirectory(t, parent, map[string][]byte{"acknowledgments.json": data}), tc.reason)
		})
	}
}

func TestEffectObservationProtocolAttacks(t *testing.T) {
	parent := effectObservationFixture(t, "normal", effectLossName)
	for _, tc := range []struct {
		name, file, reason string
		edit               func(map[string]any)
	}{
		{"changed-producer", "source.json", "effect_producer_identity", func(v map[string]any) { v["HEAD"] = "other" }},
		{"changed-source", "source.json", "effect_producer_source", func(v map[string]any) {
			v["FilesSHA256"].(map[string]any)["internal/durabletest/phase3b3_runtime_test.go"] = "other"
		}},
		{"changed-binding-source", "binding-source.json", "effect_producer_source", func(v map[string]any) { v["phase3b4_ancillary_test.go"] = "other" }},
		{"changed-schedule", "case.json", "effect_protocol_declaration", func(v map[string]any) { v["Schedule"] = "different" }},
		{"multiple-loss", "case.json", "effect_capture_order", func(v map[string]any) { v["Captures"].([]any)[1].(map[string]any)["AcknowledgedSIGKILL"] = true }},
		{"moved-revoke", "controlled-revocation.json", "effect_controlled_revoke_join", func(v map[string]any) { v["Authorization"].(map[string]any)["Source"].(map[string]any)["Epoch"] = 2 }},
		{"aliased-revoke", "controlled-revocation.json", "effect_controlled_revoke_join", func(v map[string]any) {
			v["Decision"].(map[string]any)["Authority"].(map[string]any)["Root"].(map[string]any)["Permission"].(map[string]any)["Obligation"] = "ab-service"
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := effectAttackJSON(t, parent, tc.file, tc.edit)
			assertEffectSemanticRejection(t, editRevocationDirectory(t, parent, map[string][]byte{tc.file: data}), tc.reason)
		})
	}
	n07 := effectObservationFixture(t, "normal", "N07")
	data := effectAttackJSON(t, n07, "temporal-boundary.json", func(v map[string]any) { v["Phase"] = "returned" })
	assertEffectSemanticRejection(t, editRevocationDirectory(t, n07, map[string][]byte{"temporal-boundary.json": data}), "effect_synchronous_callback")
}

func TestEffectObservationPhysicalAndRawAttacks(t *testing.T) {
	parent := effectObservationFixture(t, "normal", effectLossName)
	for _, attack := range []string{"one-principal-renewal", "wrong-wait", "wrong-kill", "runtime-pid", "missing-ready-a", "missing-journal-a", "missing-snapshot", "duplicate-ingress", "duplicate-commit", "missing-coverage", "stale-fresh-gate", "lifecycle-as-receipt", "forged-custody", "final-only-at-cut"} {
		t.Run(attack, func(t *testing.T) {
			parent := parent
			edits := map[string][]byte{}
			switch attack {
			case "one-principal-renewal":
				edits["reopened/inspect.json"] = effectAttackJSON(t, parent, "reopened/inspect.json", func(v map[string]any) { v["Tables"].(map[string]any)["runs"].([]any)[0].(map[string]any)["epoch"] = 2 })
			case "wrong-wait", "wrong-kill":
				edits["kill.json"] = effectAttackJSON(t, parent, "kill.json", func(v map[string]any) {
					if attack == "wrong-wait" {
						v["Wait"] = "exit status 0"
					} else {
						v["PID"] = 1
					}
				})
			case "runtime-pid":
				edits["acknowledgments.json"] = effectAttackJSON(t, parent, "acknowledgments.json", func(v map[string]any) { v["Runtime"].(map[string]any)["PID"] = 1 })
			case "missing-ready-a":
				edits["endpoint-a-ready.json"] = nil
			case "missing-journal-a":
				edits["endpoint-a.jsonl"] = nil
			case "missing-snapshot":
				edits["endpoints.json"] = nil
			case "duplicate-ingress", "duplicate-commit":
				lines := strings.Split(strings.TrimSpace(string(readRevocationTestFile(t, parent, "endpoint-b.jsonl"))), "\n")
				if attack == "duplicate-ingress" {
					lines = append([]string{lines[0]}, lines...)
				} else {
					lines = append(lines, lines[1])
				}
				edits["endpoint-b.jsonl"] = []byte(strings.Join(lines, "\n") + "\n")
			case "missing-coverage":
				edits["endpoints.json"] = effectAttackJSON(t, parent, "endpoints.json", func(v map[string]any) { v["b"].(map[string]any)["Applications"] = []any{} })
			case "stale-fresh-gate":
				var raw RawSnapshot
				if err := json.Unmarshal(readRevocationTestFile(t, parent, "final/inspect.json"), &raw); err != nil {
					t.Fatal(err)
				}
				for _, row := range raw.Tables["composition_facts"] {
					f := decodeCell[Fact](row, "descriptor")
					if f.AuthorityAdmission != nil && f.AuthorityAdmission.Ordinal == 2 {
						f.AuthorityAdmission.Source.Epoch = 2
						put(row, "descriptor", f)
					}
				}
				rehashFacts(&raw)
				edits["final/inspect.json"] = bytesJSON(raw)
			case "lifecycle-as-receipt":
				var raw RawSnapshot
				if err := json.Unmarshal(readRevocationTestFile(t, parent, "final/inspect.json"), &raw); err != nil {
					t.Fatal(err)
				}
				raw.Tables["external_receipts"] = []Row{{"intent_id": bytesJSON(identity("a", 1, 4, 0)), "descriptor": raw.Tables["settlements"][0]["descriptor"]}}
				edits["final/inspect.json"] = bytesJSON(raw)
			case "forged-custody":
				n07 := effectObservationFixture(t, "normal", "N07")
				parent = n07
				var raw RawSnapshot
				if err := json.Unmarshal(readRevocationTestFile(t, parent, "final/inspect.json"), &raw); err != nil {
					t.Fatal(err)
				}
				for _, row := range raw.Tables["history"] {
					if textCell(row, "subject") == "b" && intCell(row, "sequence") == 8 {
						row["selected_observation"] = bytesJSON(identity("observation", "message", "forged"))
					}
				}
				edits["final/inspect.json"] = bytesJSON(raw)
			case "final-only-at-cut":
				parent = effectObservationFixture(t, "normal", "N07")
				edits["before-revoke/inspect.json"] = readRevocationTestFile(t, parent, "final/inspect.json")
			}
			dir := editRevocationDirectory(t, parent, edits)
			r, err := CheckEffectObservationDirectory(dir)
			if err != nil {
				t.Fatal("rehashed attack should reach a semantic checker", err)
			}
			if r.Supported {
				t.Fatalf("semantic assertion: %s falsely accepted", attack)
			}
			t.Logf("raw=%v continuity=%v physical=%v physicalFailure=%+v reasons=%v", r.RawValid, r.ContinuityValid, r.Physical.Valid, r.Physical.Failure, r.Reasons)
		})
	}
}

func TestEffectObservationCoherentRepeatedIngress(t *testing.T) {
	parent := effectObservationFixture(t, "normal", effectLossName)
	lines := strings.Split(strings.TrimSpace(string(readRevocationTestFile(t, parent, "endpoint-b.jsonl"))), "\n")
	var ingress, commit endpointEvent
	if err := json.Unmarshal([]byte(lines[0]), &ingress); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(lines[1]), &commit); err != nil {
		t.Fatal(err)
	}
	ingress.Contact = 2
	commit.Contact = 2
	// A second real ingress with the same idempotency key may return the same
	// application receipt. Keep interval counts, request list and commit ordinals
	// coherent; the one-contact relation must still reject this repeated exposure.
	journal := strings.Join(append(lines, string(bytesJSON(ingress)), string(bytesJSON(commit))), "\n") + "\n"
	snapshot := effectAttackJSON(t, parent, "endpoints.json", func(v map[string]any) {
		b := v["b"].(map[string]any)
		requests := b["Requests"].([]any)
		b["Contacts"] = 2
		b["Requests"] = append(requests, requests[0])
	})
	r, err := CheckEffectObservationDirectory(editRevocationDirectory(t, parent, map[string][]byte{"endpoint-b.jsonl": []byte(journal), "endpoints.json": snapshot}))
	if err != nil {
		t.Fatal(err)
	}
	if !r.RawValid || !r.ContinuityValid || r.Physical.Valid || r.Supported || r.Physical.Failure == nil || r.Physical.Failure.Edge != "more_contacts_than_opportunities" {
		t.Fatalf("semantic assertion: coherent repeated ingress accepted or failed for unrelated reason: raw=%v continuity=%v physical=%+v supported=%v", r.RawValid, r.ContinuityValid, r.Physical.Failure, r.Supported)
	}
}

func TestEffectObservationManifestProgramJoin(t *testing.T) {
	parent := effectObservationFixture(t, "normal", "N07")
	edits := map[string][]byte{}
	// Change only non-operational author metadata, consistently replacing the
	// content-addressed program across every raw capture. Raw execution and
	// continuity remain valid, but this is no longer the pinned fixture program.
	for _, cut := range []string{"before-revoke", "after-revoke", "final"} {
		path := cut + "/inspect.json"
		var raw RawSnapshot
		if err := json.Unmarshal(readRevocationTestFile(t, parent, path), &raw); err != nil {
			t.Fatal(err)
		}
		rewriteProgram(&raw, "a", func(p *Program) { p.Author = json.RawMessage(`{"changed":true}`) })
		edits[path] = bytesJSON(raw)
	}
	assertEffectSemanticRejection(t, editRevocationDirectory(t, parent, edits), "effect_manifest_programs")
}
