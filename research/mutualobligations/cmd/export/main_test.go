package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestExitCodesSeparateCounterexampleFromApparatus(t *testing.T) {
	for _, args := range [][]string{{"--profile", "unknown"}, {"--unknown"}, {"unexpected"}} {
		var out, diagnostic bytes.Buffer
		if got := run(args, &out, &diagnostic); got != 2 {
			t.Fatalf("bad apparatus arguments exit %d", got)
		}
	}
	var out, diagnostic bytes.Buffer
	args := []string{"--profile", "serial/serial", "--assert-naive-completion"}
	if got := run(args, &out, &diagnostic); got != 1 {
		t.Fatalf("naive completion exit %d", got)
	}
	if !json.Valid(out.Bytes()) || diagnostic.Len() == 0 {
		t.Fatal("counterexample lacks machine-readable evidence/diagnostic")
	}
}
