package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestExitDisposition(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want int
	}{
		{name: "safe-local", args: []string{"-config", "S2/Local", "-assert-authority-safe", "-assert-universal-completion"}, want: 0},
		{name: "ambient-counterexample", args: []string{"-config", "S3/Missing", "-assert-authority-safe"}, want: 1},
		{name: "invention-counterexample", args: []string{"-config", "Cartesian/ScopeMix", "-assert-authority-safe"}, want: 1},
		{name: "missing-completion", args: []string{"-config", "S2/Missing", "-assert-universal-completion"}, want: 1},
		{name: "unknown-config", args: []string{"-config", "unknown"}, want: 2},
		{name: "bad-flag", args: []string{"-unknown"}, want: 2},
		{name: "positional", args: []string{"unexpected"}, want: 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out, err bytes.Buffer
			if got := run(tc.args, &out, &err); got != tc.want {
				t.Fatalf("exit %d want %d: %s", got, tc.want, err.String())
			}
			if tc.want == 1 && out.Len() == 0 {
				t.Fatal("counterexample output lost")
			}
		})
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("injected output failure") }
func TestOutputFailureIsApparatusError(t *testing.T) {
	var err bytes.Buffer
	if got := run([]string{"-config", "S2/Missing"}, failingWriter{}, &err); got != 2 {
		t.Fatalf("exit %d", got)
	}
}
func TestDeterministicExport(t *testing.T) {
	var first, second, err bytes.Buffer
	args := []string{"-config", "S2/Missing"}
	if run(args, &first, &err) != 0 || run(args, &second, &err) != 0 {
		t.Fatal(err.String())
	}
	if !bytes.Equal(first.Bytes(), second.Bytes()) {
		t.Fatal("nondeterministic export")
	}
}
