package componentinterfaces

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
)

func TestCompleteEvidence(t *testing.T) {
	e, err := buildEvidence()
	if err != nil {
		t.Fatal(err)
	}
	if len(e.Plants) != 72 || len(e.Monitors) != 5 || len(e.Controllers) != 324 || len(e.Products) != 1980 || len(e.Comparisons) != 10 || len(e.Groupings) != 1620 || len(e.Substitutions) != 1188 || len(e.Histories) != 8 {
		t.Fatal("incomplete declared roster")
	}
	if e.Counts.SuccessfulSubstitutions == 0 || e.Counts.BlockedSubstitutions == 0 {
		t.Fatal("lost one substitution class")
	}
	for _, s := range e.Substitutions {
		if s.CanCompleteLeft != s.CanCompleteRight {
			t.Fatal("equal-language substitution failed")
		}
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write failure") }
func TestExport(t *testing.T) {
	// Export resolves repository-relative source paths; tests run in package cwd.
	t.Chdir("../..")
	var a, b bytes.Buffer
	if err := Export(&a); err != nil {
		t.Fatal(err)
	}
	if err := Export(&b); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a.Bytes(), b.Bytes()) {
		t.Fatal("nondeterministic export")
	}
	var e evidence
	if err := json.Unmarshal(a.Bytes(), &e); err != nil {
		t.Fatal(err)
	}
	if len(e.Sources) == 0 || e.Schema != "alan/phase3c3/v1" {
		t.Fatal("missing source/schema")
	}
	if err := Export(failingWriter{}); err == nil {
		t.Fatal("swallowed writer failure")
	}
}
