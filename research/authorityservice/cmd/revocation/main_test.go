package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestRejectsArguments(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"unexpected"}, &out, &errOut); code != 2 || out.Len() != 0 || errOut.Len() == 0 {
		t.Fatalf("code=%d out=%s err=%s", code, &out, &errOut)
	}
}
func TestExportsIndependentEvidence(t *testing.T) {
	if testing.Short() {
		t.Skip("complete finite graph export")
	}
	var out, errOut bytes.Buffer
	if code := run(nil, &out, &errOut); code != 0 {
		t.Fatalf("code=%d: %s", code, &errOut)
	}
	var result report
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Graphs) != 1 || len(result.Controls.Histories) != 34 || len(result.Controls.Chains) != 0 {
		t.Fatalf("unexpected coverage: graphs=%d histories=%d chains=%d", len(result.Graphs), len(result.Controls.Histories), len(result.Controls.Chains))
	}
	if result.Coverage == "" {
		t.Fatal("missing finite-coverage qualification")
	}
}

type brokenWriter struct{}

func (brokenWriter) Write([]byte) (int, error) { return 0, fmt.Errorf("injected output failure") }
func TestOutputFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("complete finite graph export")
	}
	var errOut bytes.Buffer
	if code := run(nil, brokenWriter{}, &errOut); code != 2 || !strings.Contains(errOut.String(), "injected output failure") {
		t.Fatalf("code=%d error=%s", code, &errOut)
	}
}
