package main

import (
	"bytes"
	"encoding/json"
	capacity "github.com/JakeFAU/alan-research/research/capacitycomposition"
	"testing"
)

// Exercise actual finite evidence and the output encoder; full CLI and expensive
// probe execution are checked by the two recorded deterministic export runs.
func TestExportUsesFrozenSchemaAndCompleteRoster(t *testing.T) {
	evidence, err := capacity.BuildFiniteEvidence()
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := writeEvidence(&out, evidence); err != nil {
		t.Fatal(err)
	}
	var e struct {
		Schema    string
		Graphs    []json.RawMessage
		Groupings []json.RawMessage
		Witnesses []json.RawMessage
	}
	if err := json.Unmarshal(out.Bytes(), &e); err != nil {
		t.Fatal(err)
	}
	if e.Schema != "alan/phase3c2/v1" || len(e.Graphs) != 72 || len(e.Groupings) != 72 || len(e.Witnesses) != 10 {
		t.Fatalf("incomplete frozen export: schema %s graphs %d groupings %d witnesses %d", e.Schema, len(e.Graphs), len(e.Groupings), len(e.Witnesses))
	}
}
func TestExportRejectsArguments(t *testing.T) {
	var out bytes.Buffer
	if run(&out, []string{"unexpected"}) == nil {
		t.Fatal("ignored argument")
	}
}
