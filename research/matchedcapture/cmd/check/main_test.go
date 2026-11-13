package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"
)

func TestUsageAndRejectedDirectory(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want int
	}{{"missing", nil, 2}, {"flag", []string{"-bad"}, 2}, {"extra", []string{"-pair", "x", "extra"}, 2}, {"empty", []string{"-pair", t.TempDir()}, 1}} {
		t.Run(tc.name, func(t *testing.T) {
			var out, diag bytes.Buffer
			if got := run(tc.args, &out, &diag); got != tc.want {
				t.Fatalf("status %d want %d: %s", got, tc.want, diag.String())
			}
			if tc.want == 1 {
				var report struct{ Supported bool }
				if e := json.Unmarshal(out.Bytes(), &report); e != nil || report.Supported {
					t.Fatalf("failure report: %v", e)
				}
			}
		})
	}
}

type badWriter struct{}

func (badWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }
func TestWriteError(t *testing.T) {
	if got := run([]string{"-pair", t.TempDir()}, badWriter{}, io.Discard); got != 2 {
		t.Fatal(got)
	}
}
