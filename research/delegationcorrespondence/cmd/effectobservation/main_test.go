package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestEffectObservationUsage(t *testing.T) {
	for _, args := range [][]string{nil, {"-unknown"}, {"-case", "dir", "extra"}, {"-case", t.TempDir()}} {
		var diagnostics bytes.Buffer
		if status := run(args, io.Discard, &diagnostics); status != 2 || diagnostics.Len() == 0 {
			t.Errorf("run(%v) = %d, %q", args, status, &diagnostics)
		}
	}
}

type failedWriter struct{}

func (failedWriter) Write([]byte) (int, error) { return 0, errors.New("failed write") }

func TestEffectObservationCommand(t *testing.T) {
	root := os.Getenv("ALAN_EFFECT_OBSERVATION_CORPUS")
	if root == "" {
		t.Skip("set ALAN_EFFECT_OBSERVATION_CORPUS to the verified extraction root")
	}
	data, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]struct{ Cases map[string]string }
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	cases := manifest["raw-normal.tar.gz"].Cases
	for _, tc := range []struct {
		name string
		want int
	}{
		{"N06", 1},
		{"N07", 0},
		{"applied-unknown-revoked-before-fresh-gate", 0},
	} {
		var output, diagnostics bytes.Buffer
		args := []string{"-case", filepath.Join(root, "raw-normal", cases[tc.name])}
		if status := run(args, &output, &diagnostics); status != tc.want || diagnostics.Len() != 0 {
			t.Fatalf("%s status=%d stderr=%q", tc.name, status, &diagnostics)
		}
		var report struct {
			Supported, RawValid bool
			Captures            []struct {
				Path      string
				Supported bool
				Store     json.RawMessage
			}
		}
		if err := json.Unmarshal(output.Bytes(), &report); err != nil || !report.RawValid || report.Supported != (tc.want == 0) {
			t.Fatalf("%s status flags or serialized report: %v", tc.name, err)
		}

		if tc.want == 0 {
			count := 3
			if tc.name != "N07" {
				count = 5
			}
			if len(report.Captures) != count {
				t.Fatalf("%s lost serialized captures", tc.name)
			}
			for i, cut := range report.Captures {
				supported := tc.name == "N07" || i >= 2
				if cut.Supported != supported || (string(cut.Store) != "null") != supported {
					t.Fatalf("%s cut %d serialized support/observation mismatch", tc.name, i)
				}
			}
		}

		if status := run(args, failedWriter{}, &diagnostics); status != 2 {
			t.Errorf("writer error gave status %d", status)
		}
	}
}
