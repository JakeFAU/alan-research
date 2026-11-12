package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	c "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
)

func TestProgressSelection(t *testing.T) {
	if testing.Short() {
		t.Skip("complete graph query")
	}
	var out, errOut bytes.Buffer
	if code := run([]string{"-config", "S2/SeparateGrantRevocation"}, &out, &errOut); code != 0 {
		t.Fatalf("exit %d: %s", code, &errOut)
	}
	var reports []c.ProgressReport
	if err := json.Unmarshal(out.Bytes(), &reports); err != nil {
		t.Fatal(err)
	}
	if len(reports) != 1 || reports[0].Config.Scenario != "SeparateGrantRevocation" || reports[0].PredicateMismatches != 0 || reports[0].Revocation == nil {
		t.Fatal("incorrect opt-in selection")
	}
}
func TestProgressRejectsArguments(t *testing.T) {
	for _, args := range [][]string{{"unexpected"}, {"-config", "S3/SeparateGrantRevocation"}, {"-unknown"}} {
		t.Run(args[0], func(t *testing.T) {
			var out, errOut bytes.Buffer
			if code := run(args, &out, &errOut); code != 2 || out.Len() != 0 || errOut.Len() == 0 {
				t.Fatalf("exit=%d out=%s error=%s", code, &out, &errOut)
			}
		})
	}
}

type failedWriter struct{}

func (failedWriter) Write([]byte) (int, error) { return 0, errors.New("injected output failure") }
func TestProgressOutputFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("complete graph query")
	}
	var errOut bytes.Buffer
	if code := run([]string{"-config", "S2/SeparateGrantRevocation"}, failedWriter{}, &errOut); code != 2 {
		t.Fatalf("exit=%d error=%s", code, &errOut)
	}
}
