package delegationcorrespondence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	m "github.com/JakeFAU/alan-research/research/authorityservice"
	b "github.com/JakeFAU/alan-research/research/mutualobligations"
)

func readRevocationTestFile(t *testing.T, dir, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func editRevocationDirectory(t *testing.T, parent string, edits map[string][]byte) string {
	t.Helper()
	var hashes map[string]string
	if err := strictJSON(readRevocationTestFile(t, parent, "digests.json"), &hashes); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for name := range hashes {
		data, edited := edits[name]
		if !edited {
			data = readRevocationTestFile(t, parent, name)
		}
		if edited && data == nil {
			delete(hashes, name)
			continue
		}
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatal(err)
		}
		hashes[name] = sha(data)
	}
	if err := os.WriteFile(filepath.Join(dir, "digests.json"), bytesJSON(hashes), 0644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestRevocationManifestRootAttacks(t *testing.T) {
	parent := revocationFixture(t, "normal", "N07")
	for _, attack := range []string{"swapped", "aliased"} {
		t.Run(attack, func(t *testing.T) {
			var manifest map[string]any
			if err := json.Unmarshal(readRevocationTestFile(t, parent, "case.json"), &manifest); err != nil {
				t.Fatal(err)
			}
			rec := manifest["Manifest"].(map[string]any)["Reciprocal"].(map[string]any)
			want := "exact_root_mapping"
			if attack == "swapped" {
				rec["ABServiceRoot"], rec["BAServiceRoot"] = rec["BAServiceRoot"], rec["ABServiceRoot"]
			} else {
				rec["ABServiceRoot"] = rec["BAServiceRoot"]
				want = "aliased_exact_roots"
			}
			dir := editRevocationDirectory(t, parent, map[string][]byte{"case.json": bytesJSON(manifest)})
			report, err := CheckRevocationDirectory(dir)
			if err != nil {
				t.Fatal(err)
			}
			if !report.RawValid || !report.Physical.Valid || report.Legacy.ContinuityFailure != nil ||
				report.Supported || !slices.Contains(report.Reasons, want) {
				t.Fatalf("semantic assertion: %s complete roots must reject at %s, got support=%v raw=%v physical=%v reasons=%v",
					attack, want, report.Supported, report.RawValid, report.Physical.Valid, report.Reasons)
			}
			t.Logf("raw and physical valid; projector rejected at %s", want)
		})
	}
}

func TestRevocationPhysicalArtifactsRequired(t *testing.T) {
	for _, missing := range []string{"endpoint-a.jsonl", "endpoint-b.jsonl", "endpoints.json"} {
		t.Run(missing, func(t *testing.T) {
			parent := revocationFixture(t, "normal", "N06")
			dir := editRevocationDirectory(t, parent, map[string][]byte{missing: nil})
			report, err := CheckRevocationDirectory(dir)
			if err != nil {
				t.Fatal(err)
			}
			f := report.Physical.Failure
			if !report.RawValid || !report.Legacy.External.Valid || report.Supported || report.Physical.Valid ||
				f == nil || f.Edge != "required_endpoint_artifact_missing" || f.ID != missing {
				t.Fatalf("missing %s: support=%v raw=%v physical=%+v", missing, report.Supported, report.RawValid, f)
			}
			t.Logf("raw and legacy external valid; physical evidence rejected at %s", f.Edge)
		})
	}
}

func TestRevocationRawAttacks(t *testing.T) {
	parent := revocationFixture(t, "normal", "N06")
	for _, tc := range []struct{ attack, edge string }{
		{"grant-before-source", "lifecycle_before_source"},
		{"revoke-before-first-allow", "first_authority_policy"},
		{"fresh-allow-after-revoke", "fresh_gate_policy"},
		{"stale-gate-epoch", "activation_bounds"},
		{"result-from-refused-gate", "result_without_allowed_gate"},
		{"future-selection-witness", "first_authority_policy"},
		{"wrong-original-witness", "first_authority_policy"},
	} {
		t.Run(tc.attack, func(t *testing.T) {
			var raw RawSnapshot
			if err := strictJSON(readRevocationTestFile(t, parent, "final/inspect.json"), &raw); err != nil {
				t.Fatal(err)
			}
			mutateRevocationRaw(t, &raw, tc.attack)
			dir := editRevocationDirectory(t, parent, map[string][]byte{"final/inspect.json": bytesJSON(raw)})
			report, err := CheckRevocationDirectory(dir)
			if err != nil {
				t.Fatal(err)
			}
			f := report.Legacy.Final.Failure
			if report.RawValid || report.Supported || f == nil || f.Edge != tc.edge {
				t.Fatalf("semantic assertion: %s must reject at %s, got support=%v raw=%v failure=%+v",
					tc.attack, tc.edge, report.Supported, report.RawValid, f)
			}
			for _, cut := range report.Captures {
				if len(cut.Events) != 0 {
					t.Fatalf("raw-invalid dossier produced projection events at %s", cut.Path)
				}
			}
			t.Logf("raw rejected at %s order %d", f.Edge, f.Order)
		})
	}
}

func mutateRevocationRaw(t *testing.T, raw *RawSnapshot, attack string) {
	t.Helper()
	attemptRow := raw.Tables["external_attempts"][0]
	attempt := decodeCell[ExternalAttempt](attemptRow, "descriptor")
	id := attempt.Intent.ID
	if attack == "stale-gate-epoch" {
		advanceRevocationActivation(raw, attempt.Intent.Run)
		attempt.Request.Epoch++
		put(attemptRow, "descriptor", attempt)
	}
	if attack == "future-selection-witness" || attack == "wrong-original-witness" {
		if attack == "future-selection-witness" {
			attempt.Authority.RequestConsumedSequence = attempt.Intent.Sequence + 1
		} else {
			q := identity("a", 1, 1, 0)
			attempt.Authority.OriginalRequestID = q
			attempt.Authority.RequestObservationID = identity("observation", "message", q)
		}
		put(attemptRow, "descriptor", attempt)
	}
	for _, row := range raw.Tables["composition_reservations"] {
		r := decodeCell[Reservation](row, "descriptor")
		r.AttemptDigest = digest("alan/attempt-row/v1", rawCell(attemptRow, "descriptor"))
		put(row, "descriptor", r)
	}
	for _, row := range raw.Tables["composition_facts"] {
		f := decodeCell[Fact](row, "descriptor")
		switch attack {
		case "grant-before-source":
			if f.Order == 7 || f.Order == 8 {
				f.Order = 15 - f.Order
			}
		case "revoke-before-first-allow":
			if f.Order == 16 {
				f.Order = 20
			} else if f.Order > 16 && f.Order <= 20 {
				f.Order--
			}
		case "fresh-allow-after-revoke":
			if f.AuthorityAdmission != nil {
				f.AuthorityAdmission.Allowed, f.AuthorityAdmission.Reason = true, "authorized"
			}
		}
		if f.AuthorityAdmission != nil {
			f.AuthorityAdmission.AttemptDigest = digest("alan/attempt-row/v1", rawCell(attemptRow, "descriptor"))
		}
		row["decision_order"] = bytesJSON(f.Order)
		put(row, "descriptor", f)
	}
	if attack == "revoke-before-first-allow" {
		for _, row := range raw.Tables["composition_charges"] {
			c := decodeCell[Charge](row, "descriptor")
			c.DecisionOrder = 20
			row["decision_order"] = bytesJSON(c.DecisionOrder)
			put(row, "descriptor", c)
		}
	}
	if attack == "result-from-refused-gate" {
		receipt := Receipt{Key: attempt.Request.Key, Profile: attempt.Request.Profile, Delta: attempt.Request.Delta, Mutation: 1, Total: 1}
		raw.Tables["external_receipts"] = []Row{{"intent_id": bytesJSON(id), "descriptor": bytesJSON(string(bytesJSON(receipt)))}}
		f := Fact{
			Version: "alan/composition-fact/v1", ID: digest("alan/fact/external_result/v1", id), Kind: "external_result",
			Order: 23, RootIDs: []string{attempt.Authority.ServingRootID}, Run: attempt.Intent.Run, Sequence: attempt.Intent.Sequence, OperationID: id,
		}
		raw.Tables["composition_facts"] = append(raw.Tables["composition_facts"],
			Row{"id": bytesJSON(f.ID), "kind": bytesJSON(f.Kind), "decision_order": bytesJSON(f.Order), "descriptor": bytesJSON(string(bytesJSON(f)))})
	}
	rehashFacts(raw)
}

// Move the serving history and attempt to activation 2 while keeping the
// gate's old activation 1 and recomputing its dependent attempt digest.
func advanceRevocationActivation(raw *RawSnapshot, run RunID) {
	for _, table := range []string{"runs", "history"} {
		for _, row := range raw.Tables[table] {
			if (RunID{textCell(row, "subject"), intCell(row, "generation")}) == run {
				row["epoch"] = bytesJSON(2)
			}
		}
	}
	for _, row := range raw.Tables["attempts"] {
		a := decodeCell[Authorization](row, "result")
		if a.Source.Run == run {
			a.Source.Epoch = 2
			put(row, "result", a)
		}
	}
	for _, row := range raw.Tables["lifecycle_results"] {
		l := decodeCell[LifecycleResult](row, "result")
		for _, target := range []*TargetState{l.Before, l.After} {
			if target != nil && target.Run == run {
				target.Epoch = 2
			}
		}
		put(row, "result", l)
	}
}

func TestRevocationWrongTarget(t *testing.T) {
	parent := revocationFixture(t, "normal", "N07")
	legacy, files, err := checkBindingDirectory(parent)
	if err != nil {
		t.Fatal(err)
	}
	var entry revocationManifest
	if err := json.Unmarshal(files["case.json"], &entry); err != nil {
		t.Fatal(err)
	}
	phase := legacy.Phases[1]
	var raw RawSnapshot
	if err := strictJSON(files[phase.Path], &raw); err != nil {
		t.Fatal(err)
	}
	changeRevocationTarget(t, &raw, entry.Manifest.Reciprocal.ABServiceRoot)
	data := bytesJSON(raw)
	rec := entry.Manifest.Reciprocal
	phase.Verdict = CheckBindingRaw(data, BindingSpec{Principals: []BindingChain{
		{rec.ABRequest, rec.ABService, rec.ABReply}, {rec.BARequest, rec.BAService, rec.BAReply},
	}})
	phase.SHA256 = sha(data)
	if !phase.Verdict.Valid {
		t.Fatalf("coherent c-root revoke must remain raw valid: %+v", phase.Verdict.Failure)
	}
	cut := projectRevocationCapture(phase, false, data, entry, legacy.External.Exposures, m.GrantRevocationConfig(), files)
	if cut.Supported || cut.Failure == nil || cut.Failure.Edge != "exact_authority_lifecycle_root" {
		t.Fatalf("semantic assertion: coherent c-root revoke accepted into b-only model: %+v", cut.Failure)
	}
	if cut.Principals[b.BA].Reservation == nil || len(cut.Principals[b.BA].Gates) != 1 {
		t.Fatal("projector rejected before mapping genuine BA authorization and gate")
	}
	t.Logf("raw accepted c-root revoke; projector rejected at %s order %d", cut.Failure.Edge, cut.Failure.Order)
}

func changeRevocationTarget(t *testing.T, raw *RawSnapshot, target AuthorityRoot) {
	t.Helper()
	var operation string
	rewrite := func(i *Intent) {
		if i.Payload["operation"] == stringValue("RevokeAuthority") {
			i.Payload["claim"] = stringValue(string(bytesJSON(Revoke{bindingVersion, target})))
			operation = i.ID
		}
	}
	rewritten := false
	replaceProgram(raw, "b", func(p *Program) {
		for _, instruction := range p.Code {
			if instruction.Effect == nil {
				continue
			}
			p := instruction.Effect.Payload
			if expr := p["operation"]; expr.Literal != nil && *expr.Literal == stringValue("RevokeAuthority") {
				expr := p["claim"]
				*expr.Literal = stringValue(string(bytesJSON(Revoke{bindingVersion, target})))
				rewritten = true
			}
		}
	})
	for _, row := range raw.Tables["intents"] {
		i := decodeCell[Intent](row, "descriptor")
		rewrite(&i)
		put(row, "descriptor", i)
	}
	if !rewritten || operation == "" {
		t.Fatal("missing revocation program or operation")
	}
	for _, row := range raw.Tables["history"] {
		is := decodeCell[[]Intent](row, "intents")
		for j := range is {
			rewrite(&is[j])
		}
		put(row, "intents", is)
	}
	for _, row := range raw.Tables["attempts"] {
		q := decodeCell[AuthorizationRequest](row, "descriptor")
		rewrite(&q.Intent)
		put(row, "descriptor", q)
	}
	for _, row := range raw.Tables["lifecycle_results"] {
		if textCell(row, "id") != operation {
			continue
		}
		l := decodeCell[LifecycleResult](row, "result")
		rewrite(&l.Request.Intent)
		l.Authority.Root = target
		put(row, "descriptor", l.Request)
		put(row, "result", l)
	}
	for _, row := range raw.Tables["settlements"] {
		if textCell(row, "intent_id") != operation {
			continue
		}
		settlement := decodeCell[struct {
			Request     LifecycleRequest
			Disposition string
			Accepted    bool
			Reason      string
			Successor   int64
		}](row, "descriptor")
		rewrite(&settlement.Request.Intent)
		put(row, "descriptor", settlement)
	}

	rehashFacts(raw)
}

func TestRevocationEndpointSnapshotCoverage(t *testing.T) {
	parent := revocationFixture(t, "normal", "N07")
	var snapshot map[string]map[string]json.RawMessage
	if err := strictJSON(readRevocationTestFile(t, parent, "endpoints.json"), &snapshot); err != nil {
		t.Fatal(err)
	}
	snapshot["b"]["Applications"] = json.RawMessage("null")
	dir := editRevocationDirectory(t, parent, map[string][]byte{"endpoints.json": bytesJSON(snapshot)})
	report, err := CheckRevocationDirectory(dir)
	if err != nil {
		t.Fatal(err)
	}
	f := report.Physical.Failure
	if !report.RawValid || !report.Legacy.External.Valid || report.Supported || report.Physical.Valid ||
		f == nil || f.Edge != "endpoint_snapshot_coverage" || f.ID != "b" {
		t.Fatalf("incomplete snapshot: support=%v raw=%v physical=%+v", report.Supported, report.RawValid, f)
	}
}

func TestRevocationCaptureOrder(t *testing.T) {
	parent := revocationFixture(t, "normal", "N06")
	var manifest map[string]any
	if err := json.Unmarshal(readRevocationTestFile(t, parent, "case.json"), &manifest); err != nil {
		t.Fatal(err)
	}
	captures := manifest["Captures"].([]any)
	captures[0], captures[1] = captures[1], captures[0]
	dir := editRevocationDirectory(t, parent, map[string][]byte{"case.json": bytesJSON(manifest)})
	report, err := CheckRevocationDirectory(dir)
	if err != nil {
		t.Fatal(err)
	}
	f := report.Legacy.ContinuityFailure
	if !report.RawValid || !report.Physical.Valid || report.Supported || f == nil || f.Edge != "capture_order" {
		t.Fatalf("reordered captures: support=%v raw=%v continuity=%+v", report.Supported, report.RawValid, f)
	}
}
