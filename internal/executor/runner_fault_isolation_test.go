package executor

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	recordschema "github.com/JakeFAU/alan-research/schema"
)

func TestRunnerIsolatesReductionFault(t *testing.T) {
	input, err := (recordschema.Record{"short": {Kind: "string", MaxBytes: 8}}).Retain()
	if err != nil {
		t.Fatal(err)
	}
	output, err := (recordschema.Record{"short": {Kind: "string", MaxBytes: 1}}).Retain()
	if err != nil {
		t.Fatal(err)
	}
	// Retained code can reach an invalid output even if the current compiler rejects it.
	program := Program{
		Semantics: AuthorSemanticVersion,
		Code:      []Instruction{{Op: "end", Outcome: "succeeded"}},
		Author: &ProgramMetadata{
			Compiler: AuthorCompilerVersion, Positions: []string{"end"},
			Inputs: input, Outputs: output,
		},
	}
	encoded, id, err := EncodeProgram(program)
	if err != nil {
		t.Fatal(err)
	}
	fault := RunID{Subject: "a-fault", Generation: 1}
	provision := Provision{
		Programs: []Artifact{{ID: id, Semantics: program.Semantics, Bytes: encoded}},
		Subjects: []Subject{{ID: fault.Subject, Generation: 1, Revision: 1, Owner: "owner"}},
		Runs: []InitialRun{{
			ID: fault, ProgramID: id,
			Projection: Projection{Status: "open", Data: map[string]Value{"short": Text("oversize")}},
		}},
	}
	addProgram(t, &provision, "z-healthy", []Instruction{
		{Op: "set", Dst: "ran", A: Lit(Int(1))},
		{Op: "end", Outcome: "succeeded"},
	}, map[string]Value{})
	limits := DefaultLimits()
	path := filepath.Join(t.TempDir(), "fault.sqlite")
	s, err := Create(t.Context(), path, provision, limits)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	for _, wantSteps := range []int{2, 0} {
		report, err := s.Run(t.Context(), RunOptions{MaxSteps: 8, MaxDispatch: 8})
		healthy := current(t, s, "z-healthy")
		blocked := current(t, s, fault.Subject)
		if err == nil || !strings.Contains(err.Error(), "a-fault") || !strings.Contains(err.Error(), "output schema") ||
			report.Quiescent || report.Steps != wantSteps || !slices.Equal(report.OpenRuns, []RunID{fault}) ||
			healthy.Projection.Status != "succeeded" || healthy.Sequence != 2 ||
			blocked.Projection.Status != "open" || blocked.Sequence != 0 {
			t.Fatalf("Run: %+v, %v, healthy %+v, blocked %+v", report, err, healthy, blocked)
		}
	}
	replay, err := Replay(t.Context(), path, limits)
	if err != nil || replay.HistoryCount != 2 {
		t.Fatalf("Replay: %+v, %v", replay, err)
	}
}

func TestRunnerSkipsObservationWhoseBindingExceedsProjectionLimit(t *testing.T) {
	provision := Provision{Permissions: []Permission{
		{Kind: "input", Source: "local", Target: "receiver", Operation: "input", Obligation: "work"},
	}}
	addProgram(t, &provision, "receiver", inputProgram(), map[string]Value{})
	limits := DefaultLimits()
	limits.MaxDataBytes = 256
	path := filepath.Join(t.TempDir(), "observations.sqlite")
	s, err := Create(t.Context(), path, provision, limits)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	admit := func(id, value string) string {
		t.Helper()
		a, err := s.AdmitInput(t.Context(), Input{TransportID: id, LogicalID: id, Target: "receiver",
			Obligation: "work", Payload: map[string]Value{"number": Text(value)}})
		if err != nil || !a.Accepted {
			t.Fatalf("admit %s: %+v, %v", id, a, err)
		}
		return a.ObservationID
	}
	larger := admit("1-oversize", strings.Repeat("a", 512))
	if _, err := s.Commit(t.Context(), current(t, s, "receiver").Authority(), larger); !errors.Is(err, ErrLimit) {
		t.Fatalf("oversize binding: %v", err)
	}
	for range 2 {
		report, err := s.Run(t.Context(), RunOptions{MaxSteps: 8, MaxDispatch: 8})
		receiver := current(t, s, "receiver")
		if !errors.Is(err, ErrLimit) || !strings.Contains(err.Error(), "1-oversize") ||
			report.Quiescent || report.Steps != 0 || receiver.Sequence != 0 {
			t.Fatalf("oversize-only Run: %+v, %v, receiver %+v", report, err, receiver)
		}
	}
	smaller := admit("2-fits", "fits")

	report, err := s.Run(t.Context(), RunOptions{MaxSteps: 8, MaxDispatch: 8})
	receiver := current(t, s, "receiver")
	if err != nil || !report.Quiescent || report.Steps != 2 ||
		receiver.Projection.Status != "succeeded" || receiver.Projection.Data["number"] != Text("fits") {
		t.Fatalf("Run: %+v, %v, receiver %+v", report, err, receiver)
	}
	inbox, err := readObservations(t.Context(), s.conn, limits)
	if err != nil || len(inbox) != 2 || inbox[larger].consumed.Valid || !inbox[smaller].consumed.Valid ||
		inbox[smaller].consumed.Int64 != 1 {
		t.Fatalf("observations: %+v, %v", inbox, err)
	}
	replay, err := Replay(t.Context(), path, limits)
	if err != nil || replay.HistoryCount != 2 {
		t.Fatalf("Replay: %+v, %v", replay, err)
	}
}
