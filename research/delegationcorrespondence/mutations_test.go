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
