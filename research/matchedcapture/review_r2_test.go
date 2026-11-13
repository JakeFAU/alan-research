package matchedcapture

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
	"testing"

	dc "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
)

// Exact boundary/pre-call bytes and descriptor strings copied from the four
// sealed branches. Extraction paths and hashes live in evidence/task2/review-r2.
// This narrow event test does not qualify an entire capture or open SQLite.
//
//go:embed testdata/native-authorization.json.gz
var nativeAuthorizationGzip []byte

func TestNativeAuthorizationOrdinal(t *testing.T) {
	z, err := gzip.NewReader(bytes.NewReader(nativeAuthorizationGzip))
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(z)
	if err != nil {
		t.Fatal(err)
	}
	if err = z.Close(); err != nil {
		t.Fatal(err)
	}
	captures := decode[map[string]struct {
		Boundaries, Worker []byte
		Attempt, Gate      string
	}](b, "native authorization excerpts")
	if len(captures) != 4 {
		t.Fatal("expected both branches in both modes")
	}
	for name, c := range captures {
		t.Run(name, func(t *testing.T) {
			worker := decodeComplete[event](c.Worker, name)
			attempt := decode[dc.ExternalAttempt]([]byte(c.Attempt), name)
			gate := decode[dc.Fact]([]byte(c.Gate), name)
			if gate.AuthorityAdmission == nil || gate.AuthorityAdmission.Ordinal != 1 || worker.AuthorityAdmissionOrdinal != 1 {
				t.Fatal("native gate/worker must carry first admission ordinal")
			}
			path := "native/first-owner/boundaries.jsonl"
			v := stage(func() {
				checkBoundaries(map[string][]byte{path: c.Boundaries}, "native", process{Config: worker, PID: worker.PID}, worker, attempt.Request)
			})
			if !v.Valid {
				t.Fatalf("native emitted boundary: %+v", v)
			}
		})
	}
}

func TestResealedAuthorizationOrdinal(t *testing.T) {
	for _, branch := range []string{"zero", "applied"} {
		for _, ordinal := range []int64{0, 2, -1} {
			t.Run(branch+"/"+string(jsonBytes(ordinal)), func(t *testing.T) {
				f := fixture(t)
				p := branch + "/first-owner/boundaries.jsonl"
				rows := lines[map[string]any](f, p)
				f[p] = nil
				for _, row := range rows {
					if row["Boundary"] == "external.authorized" {
						row["AuthorityAdmissionOrdinal"] = ordinal
						if ordinal == 0 {
							delete(row, "AuthorityAdmissionOrdinal")
						}
					}
					f[p] = append(f[p], append(jsonBytes(row), '\n')...)
				}
				worker := read[event](f, branch+"/first-owner/pre-call.json")
				local := stage(func() {
					checkBoundaries(f, branch, process{Config: worker, PID: worker.PID}, worker, selectedRequest(f, branch))
				})
				if local.Valid {
					t.Fatal("wrong ordinal accepted in mutated branch")
				}
				r, e := CheckDirectory(sealFixture(t, f))
				if e != nil || !r.Seal.Valid || !r.Raw.Valid || r.Protocol.Valid || r.Supported {
					t.Fatalf("wrong ordinal must fail despite valid reseal/raw: %v seal=%+v raw=%+v protocol=%+v", e, r.Seal, r.Raw, r.Protocol)
				}
			})
		}
	}
}
