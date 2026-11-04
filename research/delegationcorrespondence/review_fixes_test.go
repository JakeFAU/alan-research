package delegationcorrespondence

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setPermissions(raw *RawSnapshot, ps []Permission) {
	raw.Tables["composition_permissions"][0]["descriptor"] = bytesJSON(string(bytesJSON(ps)))
	raw.Tables["permissions"] = nil
	for _, p := range ps {
		raw.Tables["permissions"] = append(raw.Tables["permissions"], Row{"kind": bytesJSON(p.Kind), "source": bytesJSON(p.Source), "target": bytesJSON(p.Target), "operation": bytesJSON(p.Operation), "revision": bytesJSON(p.Revision), "obligation": bytesJSON(p.Obligation)})
	}
}
func TestAcceptedAdmissionRequiresOriginalRoute(t *testing.T) {
	for _, obligation := range []string{"ba-request", "ab-request", "grant-ready", "ba-reply", "ab-reply"} {
		t.Run(obligation, func(t *testing.T) {
			var raw RawSnapshot
			if err := strictJSON(positive(t), &raw); err != nil {
				t.Fatal(err)
			}
			ps := decodeCell[[]Permission](raw.Tables["composition_permissions"][0], "descriptor")
			kept := []Permission{}
			for _, p := range ps {
				if p.Kind != "message" || p.Obligation != obligation {
					kept = append(kept, p)
				}
			}
			if len(kept) != len(ps)-1 {
				t.Fatal("expected one original route")
			}
			setPermissions(&raw, kept)
			v := CheckRaw(bytesJSON(raw))
			if v.Valid || v.Failure.Edge != "admission_permission" {
				t.Fatalf("missing accepted %s route: %+v", obligation, v.Failure)
			}
		})
	}
}
func TestAcceptedLocalInputRequiresRevisionZeroRoute(t *testing.T) {
	var raw RawSnapshot
	if err := strictJSON(positive(t), &raw); err != nil {
		t.Fatal(err)
	}
	a := Admission{TransportID: "local-input-1", Kind: "input", LogicalID: "local-input-1", Source: RunID{Subject: "outside"}, Target: "a", Generation: 1, Obligation: "local-test", Payload: map[string]Value{"value": intValue(7)}}
	o := Observation{ID: identity("observation", "input", a.LogicalID), Run: RunID{"a", 1}, Kind: "input", Provenance: "local-input", Source: "outside", Obligation: a.Obligation, Payload: a.Payload}
	v := AdmissionResult{TransportID: a.TransportID, Accepted: true, Reason: "accepted", Bound: &o.Run, ObservationID: o.ID}
	raw.Tables["admissions"] = append(raw.Tables["admissions"], Row{"id": bytesJSON(identity("admission", "input", a.TransportID)), "domain": bytesJSON("input"), "transport_id": bytesJSON(a.TransportID), "descriptor": bytesJSON(string(bytesJSON(a))), "result": bytesJSON(string(bytesJSON(v))), "observation_id": bytesJSON(o.ID)})
	raw.Tables["observations"] = append(raw.Tables["observations"], Row{"id": bytesJSON(o.ID), "subject": bytesJSON("a"), "generation": bytesJSON(1), "descriptor": bytesJSON(string(bytesJSON(o))), "consumed_sequence": json.RawMessage("null")})
	ps := decodeCell[[]Permission](raw.Tables["composition_permissions"][0], "descriptor")
	route := Permission{"input", "outside", "a", "input", 0, a.Obligation}
	setPermissions(&raw, append(ps, route))
	if result := CheckRaw(bytesJSON(raw)); !result.Valid {
		t.Fatalf("permitted input: %+v", result.Failure)
	}
	route.Revision = 1
	setPermissions(&raw, append(ps, route))
	if result := CheckRaw(bytesJSON(raw)); result.Valid || result.Failure.Edge != "admission_permission" {
		t.Fatalf("input accepted with nonzero route revision: %+v", result.Failure)
	}
}
func TestImmutableSourceAndAuthorityContinuity(t *testing.T) {
	for _, change := range []string{"permission-added", "program-added", "deployment-owner"} {
		t.Run(change, func(t *testing.T) {
			original := positive(t)
			var raw RawSnapshot
			if err := strictJSON(original, &raw); err != nil {
				t.Fatal(err)
			}
			switch change {
			case "permission-added":
				ps := decodeCell[[]Permission](raw.Tables["composition_permissions"][0], "descriptor")
				setPermissions(&raw, append(ps, Permission{"external", "b", "https://extra.invalid", "apply", 1, "new-authority"}))
			case "program-added":
				program := []byte(`{"Semantics":"alan/1d/v1","Code":[{"Op":"end","Outcome":"failed"}]}`)
				raw.Tables["programs"] = append(raw.Tables["programs"], Row{"id": bytesJSON(sha(program)), "semantics": bytesJSON("alan/1d/v1"), "bytes": bytesJSON(string(program))})
			case "deployment-owner":
				d := decodeCell[Deployment](raw.Tables["composition_deployment"][0], "descriptor")
				for i := range d.Subjects {
					if d.Subjects[i].ID == "b" {
						d.Subjects[i].Owner = "changed-owner"
					}
				}
				raw.Tables["composition_deployment"][0]["descriptor"] = bytesJSON(string(bytesJSON(d)))
				for _, r := range raw.Tables["subjects"] {
					if textCell(r, "id") == "b" {
						r["owner"] = bytesJSON("changed-owner")
					}
				}
				for _, r := range raw.Tables["run_origins"] {
					if textCell(r, "subject") == "b" {
						o := decodeCell[Origin](r, "descriptor")
						o.Subject.Owner = "changed-owner"
						r["descriptor"] = bytesJSON(string(bytesJSON(o)))
					}
				}
			}
			changed := bytesJSON(raw)
			v := CheckRaw(changed)
			if !v.Valid {
				t.Fatalf("control must remain individually valid: %+v", v.Failure)
			}
			failure := CheckContinuity([]PhaseReport{{Path: "initial/inspect.json", Verdict: CheckRaw(original)}, {Path: "final/inspect.json", Verdict: v}}, map[string][]byte{"initial/inspect.json": original, "final/inspect.json": changed})
			if failure == nil || failure.Edge != "immutable_source_changed" {
				t.Fatalf("accepted %s drift: %+v", change, failure)
			}
		})
	}
}
func frozenCaseFiles(t *testing.T, name string) map[string][]byte {
	t.Helper()
	root := os.Getenv("ALAN_TASK3_RAW")
	if root == "" {
		t.Skip("frozen corpus")
	}
	index, err := os.ReadFile(filepath.Join(root, "case-index.json"))
	if err != nil {
		t.Fatal(err)
	}
	var entries []IndexEntry
	if err = strictJSON(index, &entries); err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name != name {
			continue
		}
		data, err := os.ReadFile(filepath.Join(e.Directory, "digests.json"))
		if err != nil {
			t.Fatal(err)
		}
		var digests map[string]string
		if err = strictJSON(data, &digests); err != nil {
			t.Fatal(err)
		}
		files := map[string][]byte{}
		for p, h := range digests {
			if !strings.HasSuffix(p, ".json") && !strings.HasSuffix(p, ".jsonl") {
				continue
			}
			b, err := os.ReadFile(filepath.Join(e.Directory, p))
			if err != nil {
				t.Fatal(err)
			}
			if sha(b) != h {
				t.Fatal("fixture digest mismatch", p)
			}
			files[p] = b
		}
		return files
	}
	t.Fatal("missing frozen case", name)
	return nil
}
func TestFreshOwnerGateEpochRequiresOpenContext(t *testing.T) {
	files := frozenCaseFiles(t, "cut-06b-applied-no-receipt")
	rawBytes := files["final/inspect.json"]
	if v := CheckExternal(rawBytes, files); !v.Valid {
		t.Fatalf("actual owner: %+v", v.Failure)
	}
	var raw RawSnapshot
	if err := strictJSON(rawBytes, &raw); err != nil {
		t.Fatal(err)
	}
	for _, r := range raw.Tables["composition_facts"] {
		f := decodeCell[Fact](r, "descriptor")
		if f.Kind == "authority_admission" && f.AuthorityAdmission.Ordinal == 2 {
			f.AuthorityAdmission.Source.Epoch = 2
			f.RecordDigest = digest("alan/authority-admission/v1", f.AuthorityAdmission)
			r["descriptor"] = bytesJSON(string(bytesJSON(f)))
		}
	}
	var changed bytes.Buffer
	for _, line := range bytes.Split(bytes.TrimSpace(files["fresh-owner/boundaries.jsonl"]), []byte("\n")) {
		var e boundaryEvent
		if err := strictJSON(line, &e); err != nil {
			t.Fatal(err)
		}
		if e.AuthorityAdmissionOrdinal == 2 {
			e.Epoch = 2
		}
		changed.Write(bytesJSON(e))
		changed.WriteByte('\n')
	}
	files["fresh-owner/boundaries.jsonl"] = changed.Bytes()
	rawBytes = bytesJSON(raw)
	if v := CheckRaw(rawBytes); !v.Valid {
		t.Fatalf("historical bounds control: %+v", v.Failure)
	}
	if v := CheckExternal(rawBytes, files); v.Valid || v.Failure.Edge != "boundary_owner_epoch" {
		t.Fatalf("stale pair accepted against fresh owner: %+v", v.Failure)
	}
}
func TestOwnerRenewalPreservesTerminalEpoch(t *testing.T) {
	for _, name := range []string{"cut-06b-applied-no-receipt", "cut-10-one-complete"} {
		t.Run(name, func(t *testing.T) {
			files := frozenCaseFiles(t, name)
			raw := files["final/inspect.json"]
			if v := CheckExternal(raw, files); !v.Valid {
				t.Fatalf("actual renewal: %+v", v.Failure)
			}
			var opened map[string]json.RawMessage
			if err := strictJSON(files["fresh-owner/open.json"], &opened); err != nil {
				t.Fatal(err)
			}
			var runs map[string]map[string]json.RawMessage
			if err := strictJSON(opened["Runs"], &runs); err != nil {
				t.Fatal(err)
			}
			var post RawSnapshot
			if err := strictJSON(files["fresh-owner/post-open/inspect.json"], &post); err != nil {
				t.Fatal(err)
			}
			epoch := int64(2)
			if name == "cut-10-one-complete" {
				epoch = 3
			}
			runs["a"]["Epoch"] = bytesJSON(epoch)
			opened["Runs"] = bytesJSON(runs)
			files["fresh-owner/open.json"] = bytesJSON(opened)
			for _, r := range post.Tables["runs"] {
				if textCell(r, "subject") == "a" {
					r["epoch"] = bytesJSON(epoch)
				}
			}
			files["fresh-owner/post-open/inspect.json"] = bytesJSON(post)
			if v := CheckExternal(raw, files); v.Valid || v.Failure.Edge != "owner_renewal" {
				t.Fatalf("invalid renewal accepted: %+v", v.Failure)
			}
		})
	}
}
func TestUncapturedInitialRenewalRemainsExplicit(t *testing.T) {
	files := frozenCaseFiles(t, "generic-unknown")
	v := CheckExternal(files["final/inspect.json"], files)
	if !v.Valid {
		t.Fatalf("actual incomplete initial capture: %+v", v.Failure)
	}
	first, fresh := v.Owners["first-owner"], v.Owners["fresh-owner"]
	if first.RenewalChecked || first.PriorPath != "" || first.PriorSHA256 != "" || !fresh.RenewalChecked || fresh.PriorPath != "stopped/inspect.json" || len(v.OwnerBindings) == 0 {
		t.Fatalf("renewal evidence misstated: first=%+v fresh=%+v", first, fresh)
	}
	delete(files, "fresh-owner/post-open/inspect.json")
	v = CheckExternal(files["final/inspect.json"], files)
	if v.Valid || v.Failure.Edge != "owner_snapshot_missing" {
		t.Fatalf("missing current owner snapshot accepted: %+v", v.Failure)
	}
}

func TestRecipientMutationChangesRecipient(t *testing.T) {
	var original RawSnapshot
	if err := strictJSON(positive(t), &original); err != nil {
		t.Fatal(err)
	}
	foundHolder, foundRecipient := false, false
	for _, m := range RawMutations(positive(t)) {
		if m.Name == "altered-holder" {
			foundHolder = true
		}
		if m.Name != "altered-recipient" {
			continue
		}
		foundRecipient = true
		var raw RawSnapshot
		if err := strictJSON(m.Input, &raw); err != nil {
			t.Fatal(err)
		}
		before := decodeCell[LifecycleRequest](original.Tables["lifecycle_results"][0], "descriptor")
		after := decodeCell[LifecycleRequest](raw.Tables["lifecycle_results"][0], "descriptor")
		var a, b Claim
		if err := strictJSON([]byte(before.Intent.Payload["claim"].Text), &a); err != nil {
			t.Fatal(err)
		}
		if err := strictJSON([]byte(after.Intent.Payload["claim"].Text), &b); err != nil {
			t.Fatal(err)
		}
		if a.Recipient == b.Recipient || a.Root != b.Root || !equal(original.Tables["external_receipts"], raw.Tables["external_receipts"]) {
			t.Fatal("recipient mutation changes holder or loses original result evidence")
		}
		if m.Verdict.Valid {
			t.Fatal("mutated recipient accepted")
		}
		if m.SemanticField != "lifecycle_results.descriptor.Intent.Payload.claim.Text.Recipient.Subject" || len(m.CellChanges) != 2 || m.CellChanges[0].Key != rowKey("lifecycle_results", raw.Tables["lifecycle_results"][0]) || m.Verdict.Failure.Edge != "lifecycle_origin" {
			t.Fatal("missing exact recipient coordinates or original-source rejection")
		}
	}
	if !foundHolder || !foundRecipient {
		t.Fatal("holder and recipient controls must be distinct")
	}
}
