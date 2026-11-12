package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestRevocationUsage(t *testing.T) {
	for _, args := range [][]string{nil, {"-unknown"}, {"-case", "dir", "extra"}, {"-case", t.TempDir()}} {
		var diagnostics bytes.Buffer
		if status := run(args, io.Discard, &diagnostics); status != 2 || diagnostics.Len() == 0 {
			t.Errorf("run(%v) = %d, %q", args, status, &diagnostics)
		}
	}
}

type failedWriter struct{}

func (failedWriter) Write([]byte) (int, error) { return 0, errors.New("failed write") }

func TestRevocationCommand(t *testing.T) {
	root := os.Getenv("ALAN_REVOCATION_CORPUS")
	if root == "" {
		t.Skip("set ALAN_REVOCATION_CORPUS to the verified extraction root")
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
		{"N06", 0},
		{"N07", 0},
		{"applied-unknown-revoked-before-fresh-gate", 1},
	} {
		var output, diagnostics bytes.Buffer
		args := []string{"-case", filepath.Join(root, "raw-normal", cases[tc.name])}
		if status := run(args, &output, &diagnostics); status != tc.want || diagnostics.Len() != 0 {
			t.Fatalf("%s status=%d stderr=%q", tc.name, status, &diagnostics)
		}
		var report struct {
			Supported, RawValid bool
			Captures            []struct {
				Path       string
				Unresolved []string
			}
		}
		if err := json.Unmarshal(output.Bytes(), &report); err != nil || !report.RawValid || report.Supported != (tc.want == 0) {
			t.Fatalf("%s status flags or serialized report: %v", tc.name, err)
		}
		if tc.name == "N07" {
			want := []struct {
				path       string
				unresolved []string
			}{
				{"before-revoke/inspect.json", []string{"BA"}},
				{"after-revoke/inspect.json", []string{"BA"}},
				{"final/inspect.json", []string{}},
			}
			if len(report.Captures) != len(want) {
				t.Fatalf("N07 capture count = %d, want %d", len(report.Captures), len(want))
			}
			for i, cut := range report.Captures {
				if cut.Path != want[i].path || cut.Unresolved == nil || !slices.Equal(cut.Unresolved, want[i].unresolved) {
					t.Fatalf("N07 capture %d = %+v, want %+v", i, cut, want[i])
				}
			}
		}
		if status := run(args, failedWriter{}, &diagnostics); status != 2 {
			t.Errorf("writer error gave status %d", status)
		}
	}
}
