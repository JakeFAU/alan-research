package executor

import (
	"encoding/json"
	"testing"
)

func TestAcceptedDescendantScopeRemainsOpenAfterRootSuccess(t *testing.T) {
	s := invocationStore(t, "succeeded", 1, nil)
	admitted := invokeOnce(t, s)
	if !admitted.Accepted {
		t.Fatal(admitted)
	}
	for range 2 {
		if _, err := s.Commit(t.Context(), current(t, s, "root").Authority(), ""); err != nil {
			t.Fatal(err)
		}
	}
	entry := compositionMemberID(admitted.Invocation.ID, "entry")
	closed, err := compositionScopeClosed(t.Context(), s.conn, entry)
	if err != nil || closed {
		t.Fatalf("accepted descendant fenced after root success: closed=%v error=%v", closed, err)
	}
}

func TestLegacyCompositionRetainsRootSuccessFenceAcrossReopen(t *testing.T) {
	p, c := invocationFixture(t, "succeeded", 1)
	c.Version = compositionProtocolVersion
	path := t.TempDir() + "/legacy.sqlite"
	s, err := createStore(t.Context(), path, p, DefaultLimits(), &c)
	if err != nil {
		t.Fatal(err)
	}
	a := invokeOnce(t, s)
	for range 2 {
		if _, err = s.Commit(t.Context(), current(t, s, "root").Authority(), ""); err != nil {
			t.Fatal(err)
		}
	}
	entry := compositionMemberID(a.Invocation.ID, "entry")
	closed, err := compositionScopeClosed(t.Context(), s.conn, entry)
	if err != nil || !closed {
		t.Fatal(closed, err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err = Replay(t.Context(), path, DefaultLimits()); err != nil {
		t.Fatal(err)
	}
	s, err = Open(t.Context(), path, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	closed, err = compositionScopeClosed(t.Context(), s.conn, entry)
	if err != nil || !closed {
		t.Fatal(closed, err)
	}
}

func TestDrainProfileCannotBeIndividuallyDowngraded(t *testing.T) {
	for _, mutation := range []string{
		`UPDATE composition_profile SET version='alan/composition/v1', descriptor=json_set(descriptor,'$.Version','alan/composition/v1')`,
		`UPDATE metadata SET schema_version=4; PRAGMA user_version=4`,
	} {
		t.Run(mutation, func(t *testing.T) {
			p, c := invocationFixture(t, "succeeded", 1)
			path := t.TempDir() + "/profile.sqlite"
			s, err := createInvocationStore(t.Context(), path, p, c, DefaultLimits())
			if err != nil {
				t.Fatal(err)
			}
			if _, err = s.conn.ExecContext(t.Context(), mutation); err != nil {
				t.Fatal(err)
			}
			if err = s.Close(); err != nil {
				t.Fatal(err)
			}
			if _, err = Replay(t.Context(), path, DefaultLimits()); err == nil {
				t.Fatal("replay accepted partial downgrade")
			}
			opened, err := Open(t.Context(), path, DefaultLimits())
			if err == nil {
				opened.Close()
				t.Fatal("open accepted partial downgrade")
			}
		})
	}
}

func TestFailedCancelledAndRetiredRootFenceAcceptedDescendant(t *testing.T) {
	for _, outcome := range []string{"failed", "cancelled", "retired"} {
		t.Run(outcome, func(t *testing.T) {
			p, c := invocationFixture(t, "succeeded", 1)
			if outcome == "failed" {
				var root Program
				for _, a := range p.Programs {
					if a.ID == p.Runs[0].ProgramID {
						if err := json.Unmarshal(a.Bytes, &root); err != nil {
							t.Fatal(err)
						}
					}
				}
				root.Code[2].Outcome = "failed"
				raw, id, err := EncodeProgram(root)
				if err != nil {
					t.Fatal(err)
				}
				p.Programs = append(p.Programs, Artifact{ID: id, Semantics: root.Semantics, Bytes: raw})
				p.Runs[0].ProgramID = id
			} else if outcome == "cancelled" {
				addInvocationController(t, &p, "CancelRun")
			} else {
				addInvocationController(t, &p, "RetireSubject")
			}
			path := t.TempDir() + "/closed.sqlite"
			s, err := createInvocationStore(t.Context(), path, p, c, DefaultLimits())
			if err != nil {
				t.Fatal(err)
			}
			a := invokeOnce(t, s)
			if outcome == "failed" || outcome == "retired" {
				for range 2 {
					if _, err = s.Commit(t.Context(), current(t, s, "root").Authority(), ""); err != nil {
						t.Fatal(err)
					}
				}
			}
			if outcome != "failed" {
				decision := authorizedDecision(t, s, committedRequest(t, s, "controller"))
				if !decision.Accepted {
					t.Fatal(decision)
				}
			}
			entry := compositionMemberID(a.Invocation.ID, "entry")
			closed, err := compositionScopeClosed(t.Context(), s.conn, entry)
			if err != nil || !closed {
				t.Fatal(closed, err)
			}
			if err = s.Close(); err != nil {
				t.Fatal(err)
			}
			if _, err = Replay(t.Context(), path, DefaultLimits()); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestForeignRuntimeFormatsAreNotRecognized(t *testing.T) {
	for version := 9; version <= 16; version++ {
		if isCompositionSchema(protocolSchemaNumber(version)) || hasReceiptConflictSchema(version) {
			t.Fatalf("foreign runtime format %d was recognized", version)
		}
	}
}
