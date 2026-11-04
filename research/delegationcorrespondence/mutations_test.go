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
