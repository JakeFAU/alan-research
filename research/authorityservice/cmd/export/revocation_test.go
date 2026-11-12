package main

import (
	"bytes"
	"encoding/json"
	"testing"

	model "github.com/JakeFAU/alan-research/research/authorityservice"
)

func TestExportSeparateGrantRevocation(t *testing.T) {
	for _, tc := range []struct {
		flag string
		exit int
	}{
		{flag: "-assert-authority-safe", exit: 0},
		{flag: "-assert-universal-completion", exit: 1},
	} {
		t.Run(tc.flag, func(t *testing.T) {
			const name = "S2/SeparateGrantRevocation"
			var out, errOut bytes.Buffer
			if got := run([]string{"-config", name, tc.flag}, &out, &errOut); got != tc.exit {
				t.Fatalf("exit %d want %d: %s", got, tc.exit, errOut.String())
			}
			var reports []model.Report
			if err := json.Unmarshal(out.Bytes(), &reports); err != nil {
				t.Fatal(err)
			}
			if len(reports) != 1 || reports[0].Config.Name() != name || reports[0].ShortestCannotComplete == nil {
				t.Fatal("expected the opt-in report with a retained completion counterexample")
			}
		})
	}
}
