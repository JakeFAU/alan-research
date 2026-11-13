package main

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/JakeFAU/alan-research/research/effectobservation"
)

func TestExport(t *testing.T) {
	var out, diagnostics bytes.Buffer
	if code := run(nil, &out, &diagnostics); code != 0 || diagnostics.Len() != 0 {
		t.Fatalf("export exited %d: %s", code, &diagnostics)
	}
	var graph effectobservation.Graph
	if err := json.Unmarshal(out.Bytes(), &graph); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(graph, effectobservation.Explore()) {
		t.Fatal("JSON does not contain the complete graph")
	}
	var second bytes.Buffer
	if code := run(nil, &second, &diagnostics); code != 0 || !bytes.Equal(out.Bytes(), second.Bytes()) {
		t.Fatalf("second export differed or exited %d: %s", code, &diagnostics)
	}
}

func TestExportErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		out  io.Writer
		code int
		text string
	}{
		{name: "positional", args: []string{"graph"}, out: io.Discard, code: 2, text: "usage: export\n"},
		{name: "flags", args: []string{"--"}, out: io.Discard, code: 2, text: "usage: export\n"},
		{name: "write", out: closedWriter{}, code: 1, text: io.ErrClosedPipe.Error() + "\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var diagnostics bytes.Buffer
			if code := run(test.args, test.out, &diagnostics); code != test.code || diagnostics.String() != test.text {
				t.Fatalf("got exit %d and %q, want %d and %q", code, &diagnostics, test.code, test.text)
			}
		})
	}
}

type closedWriter struct{}

func (closedWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}
