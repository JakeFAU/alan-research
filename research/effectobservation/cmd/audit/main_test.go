package main

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/JakeFAU/alan-research/research/effectobservation"
)

func TestRun(t *testing.T) {
	var a, b, errOut bytes.Buffer
	if run(nil, &a, &errOut) != 0 || errOut.Len() != 0 {
		t.Fatalf("audit failed: %s", errOut.String())
	}
	if run(nil, &b, &errOut) != 0 || !bytes.Equal(a.Bytes(), b.Bytes()) {
		t.Fatal("nondeterministic audit")
	}
	var result struct {
		Report   effectobservation.AuditReport     `json:"report"`
		Controls []effectobservation.ControlResult `json:"controls"`
	}
	if err := json.Unmarshal(a.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Report.CheckedStates != 63 || result.Report.CheckedArcs != 104 || len(result.Controls) != 20 {
		t.Fatal("missing audit coverage/controls")
	}
	for _, c := range result.Controls {
		if !c.Passed {
			t.Fatalf("failed control: %s", c.Name)
		}
	}
}

type failedWriter struct{}

func (failedWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestRunFailures(t *testing.T) {
	t.Run("usage", func(t *testing.T) {
		var errOut bytes.Buffer
		if run([]string{"unexpected"}, io.Discard, &errOut) != 2 || !strings.Contains(errOut.String(), "usage: audit") {
			t.Fatal("missing usage error")
		}
	})
	t.Run("write", func(t *testing.T) {
		var errOut bytes.Buffer
		if run(nil, failedWriter{}, &errOut) != 1 || !strings.Contains(errOut.String(), io.ErrClosedPipe.Error()) {
			t.Fatal("missing write error")
		}
	})
}
