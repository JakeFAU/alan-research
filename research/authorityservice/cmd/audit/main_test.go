package main

import (
	"bytes"
	"encoding/json"
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
	if len(result.Graphs) != 16 || len(result.Controls.Histories) != 33 || len(result.Controls.Chains) != 4 {
		t.Fatalf("unexpected coverage: graphs=%d histories=%d chains=%d", len(result.Graphs), len(result.Controls.Histories), len(result.Controls.Chains))
	}
	if result.Coverage == "" {
		t.Fatal("missing finite-coverage qualification")
	}
}
