package executor

import (
	"reflect"
	"testing"
)

func TestRootPendingInvocationsDistinguishesRejectedAndUnadmitted(t *testing.T) {
	p, c := invocationFixture(t, "succeeded", 1)
	c.Roots[0].Limits.MaxCreatedRuns = 1
	path := t.TempDir() + "/pending.sqlite"
	s, err := createInvocationStore(t.Context(), path, p, c, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	first := invokeOnce(t, s)
	if first.Accepted {
		t.Fatal("expected rejected child", first)
	}
	next, err := s.Commit(t.Context(), current(t, s, "root").Authority(), "")
	if err != nil {
		t.Fatal(err)
	}
	second := next.Projection.Data["pending"].Text
	roots, err := readRootCompletions(t.Context(), s.conn)
	if err != nil || !reflect.DeepEqual(roots[0].PendingInvocations, []string{second}) {
		t.Fatalf("rejected child should not be pending; unadmitted child should: %+v, %v", roots, err)
	}
	a, err := s.AdmitInvocation(t.Context(), second)
	if err != nil || a.Accepted {
		t.Fatal(a, err)
	}
	if _, err = s.Commit(t.Context(), current(t, s, "root").Authority(), ""); err != nil {
		t.Fatal(err)
	}
	root, err := compositionRootFor(t.Context(), s.conn, RunID{"root", 1})
	if err != nil {
		t.Fatal(err)
	}
	result, err := s.SettleRoot(t.Context(), root.ID)
	if err != nil || !result.Ready || result.Outcome.Status != "succeeded" {
		t.Fatal(result, err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	replay, err := Replay(t.Context(), path, DefaultLimits())
	if err != nil || len(replay.Roots) != 1 || len(replay.Roots[0].PendingInvocations) != 0 {
		t.Fatal(replay.Roots, err)
	}
	s, err = Open(t.Context(), path, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	roots, err = readRootCompletions(t.Context(), s.conn)
	if err != nil || len(roots[0].PendingInvocations) != 0 {
		t.Fatal(roots, err)
	}
}
