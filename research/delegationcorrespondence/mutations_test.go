package delegationcorrespondence

import "testing"

func TestRawMutationCorpus(t *testing.T) {
	rs := RawMutations(positive(t))
	if len(rs) < 12 {
		t.Fatal(len(rs))
	}
	for _, r := range rs {
		if r.Verdict.Valid || r.Verdict.Failure == nil || r.OriginalSHA256 == r.MutatedSHA256 || len(r.Input) == 0 {
			t.Fatalf("%s: %+v", r.Name, r.Verdict)
		}
	}
}

func TestRawMutationsIncompletePrefix(t *testing.T) {
	raw := RawSnapshot{Tables: map[string][]Row{}}
	mutations := RawMutations(bytesJSON(raw))
	if len(mutations) != 15 {
		t.Fatalf("got %d diagnostics, want 15", len(mutations))
	}
	for _, m := range mutations {
		if m.Applicable || m.Diagnostic == "" || m.Verdict.Failure != nil || m.Verdict.Valid {
			t.Fatalf("%s did not report inapplicability: %+v", m.Name, m)
		}
		if m.OriginalSHA256 != m.MutatedSHA256 {
			t.Fatalf("%s invented evidence", m.Name)
		}
	}
}

func TestRawMutationNonGrantLifecyclePrefix(t *testing.T) {
	q := LifecycleRequest{Operation: "CancelRun", Intent: Intent{Payload: map[string]Value{"operation": stringValue("CancelRun")}}}
	raw := RawSnapshot{Tables: map[string][]Row{"lifecycle_results": {{"id": bytesJSON("ordinary-control"), "descriptor": bytesJSON(string(bytesJSON(q))), "result": bytesJSON(string(bytesJSON(LifecycleResult{Request: q})))}}}}
	mutations := RawMutations(bytesJSON(raw))
	if len(mutations) != 15 {
		t.Fatalf("lost applicability diagnostic: got %d, want15", len(mutations))
	}
	for _, m := range mutations {
		if m.Name == "altered-recipient" && (m.Applicable || m.Diagnostic == "") {
			t.Fatalf("missing grant treated as applicable: %+v", m)
		}
	}
}
