package main

import (
	"bytes"
	"encoding/json"
	"github.com/JakeFAU/alan-research/research/contractcomposition"
	"testing"
)

func TestExportDeterministicCompleteAndDecodable(t *testing.T) {
	var a, b bytes.Buffer
	if err := writeExport(&a); err != nil {
		t.Fatal(err)
	}
	if err := writeExport(&b); err != nil {
		t.Fatal(err)
	}
	var e contractcomposition.Evidence
	if err := json.Unmarshal(a.Bytes(), &e); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a.Bytes(), b.Bytes()) || len(e.Graphs) != 16 {
		t.Fatal("nondeterministic or incomplete export")
	}
	for _, g := range e.Graphs {
		arcs := 0
		for _, n := range g.Nodes {
			arcs += len(n.Arcs)
		}
		if arcs != g.Arcs {
			t.Fatal("lost graph arcs in JSON")
		}
	}
}
