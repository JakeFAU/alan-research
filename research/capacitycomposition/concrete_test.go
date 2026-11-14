package capacitycomposition

import (
	cc "github.com/JakeFAU/alan-research/research/contractcomposition"
	"github.com/JakeFAU/alan-research/research/twoagent"
	"reflect"
	"testing"
)

func take(t *testing.T, c Config, s ConcreteState, id, kind string) MicroStep {
	t.Helper()
	for _, e := range ConcreteSuccessors(c, s) {
		if string(e.Event.Kind) == kind && !e.Event.Duplicate && (e.Event.Original == id || e.Event.Actor == id) {
			return e
		}
	}
	t.Fatalf("missing concrete %s/%s", id, kind)
	return MicroStep{}
}
func prefix(t *testing.T, c Config, script [][2]string) ConcreteState {
	t.Helper()
	s, e := NewConcrete(c)
	if e != nil {
		t.Fatal(e)
	}
	for _, p := range script {
		s = take(t, c, s, p[0], p[1]).State
	}
	return s
}

var bothAdmitted = [][2]string{{"AB", "commit"}, {"AB", "send-request"}, {"AB", "admit"}, {"BC", "commit"}, {"BC", "send-request"}, {"BC", "admit"}}

func TestGenuineConcreteHolderCollision(t *testing.T) {
	c := config(t, "shared-1", "early", 0)
	s := prefix(t, c, bothAdmitted)
	a := take(t, c, s, "AB", "reserve-capacity").State
	b := take(t, c, s, "BC", "reserve-capacity").State
	if a.Base != b.Base || a.Held == b.Held || !reflect.DeepEqual(Occupancy(c, Project(c, a)), Occupancy(c, Project(c, b))) {
		t.Fatal("full base/aggregate collision missing")
	}
	if canComplete(t, c, Project(c, a)) || !canComplete(t, c, Project(c, b)) {
		t.Fatal("wrong collision classification")
	}
	crashed := take(t, c, a, "B", "crash").State
	if crashed.Held != a.Held {
		t.Fatal("crash released hold")
	}
	recovered := take(t, c, crashed, "B", "recover").State
	if recovered.Held != a.Held {
		t.Fatal("recovery changed original ownership")
	}
}
func TestApplyNeedsAdmittedHoldAndReleasesAtomically(t *testing.T) {
	c := config(t, "shared-1", "early", 7)
	s := prefix(t, c, bothAdmitted)
	for _, e := range ConcreteSuccessors(c, s) {
		if e.Event.Kind == twoagent.Apply {
			t.Fatal("apply without hold")
		}
	}
	s = take(t, c, s, "AB", "reserve-capacity").State
	for _, v := range []twoagent.Value{twoagent.R1, twoagent.R2} {
		found := false
		for _, e := range ConcreteSuccessors(c, s) {
			if e.Event.Kind == twoagent.Apply && e.Event.Result == v {
				found = true
				if e.State.Held[0] || e.State.Base.Originals[0].Retained != v || e.State.Base.Originals[0].WorldResult != v {
					t.Fatal("release not atomic with real retention")
				}
			}
		}
		if !found {
			t.Fatal("missing actual result", v)
		}
	}
	s.Base.Originals[0].Admitted = false
	if s.Validate(c) == nil {
		t.Fatal("accepted invented hold without admission")
	}
}
func TestLiftUsesNewResourceAwareAction(t *testing.T) {
	c := config(t, "shared-1", "stable-ready", 0)
	s := prefix(t, c, [][2]string{{"AB", "commit"}})
	if _, err := Lift(c, s, Action{Kind: "apply", Original: "AB", Binding: c.Base.Originals[0]}); err == nil {
		t.Fatal("lifted excluded old apply")
	}
	c = config(t, "shared-1", "early", 0)
	s = prefix(t, c, [][2]string{{"AB", "commit"}, {"A", "crash"}, {"B", "crash"}})
	a := Action{Kind: "reserve-capacity", Original: "AB", Binding: c.Base.Originals[0]}
	path, err := Lift(c, s, a)
	if err != nil {
		t.Fatal(err)
	}
	checkLift(t, c, s, a, path)
	if !path[len(path)-1].State.Base.Originals[0].Admitted {
		t.Fatal("reservation fabricated admission")
	}
	forged := a
	forged.Binding.Grant = "forged"
	if _, err := Lift(c, s, forged); err == nil {
		t.Fatal("accepted forged permission")
	}
}
func checkLift(t *testing.T, c Config, s ConcreteState, a Action, path []MicroStep) {
	t.Helper()
	q := Project(c, s)
	var target State
	found := false
	for _, e := range Successors(c, q) {
		if e.Action == a {
			target = e.State
			found = true
		}
	}
	if !found || len(path) == 0 {
		t.Fatal("not an enabled nonempty lift")
	}
	for i, e := range path {
		legal := false
		for _, candidate := range ConcreteSuccessors(c, s) {
			if candidate == e {
				legal = true
			}
		}
		if !legal {
			t.Fatal("illegal adjacency")
		}
		want := q
		if i == len(path)-1 {
			want = target
		}
		if Project(c, e.State) != want {
			t.Fatal("unrelated progress or early hold change")
		}
		if err := e.State.Validate(c); err != nil {
			t.Fatal(err)
		}
		s = e.State
	}
}
func TestCrashCustodyAndOracleIndependence(t *testing.T) {
	c := config(t, "shared-1", "early", 7)
	s := prefix(t, c, [][2]string{{"AB", "commit"}, {"AB", "send-request"}, {"AB", "admit"}, {"AB", "reserve-capacity"}, {"AB", "apply"}, {"AB", "send-result"}, {"AB", "receive-result"}, {"A", "crash"}})
	if s.Base.Originals[0].Volatile != "" || s.Base.Originals[0].Retained == "" {
		t.Fatal("wrong crash custody")
	}
	a := Action{Kind: "consume", Original: "AB", Binding: c.Base.Originals[0]}
	path, err := Lift(c, s, a)
	if err != nil {
		t.Fatal(err)
	}
	checkLift(t, c, s, a, path)
	stripped := s
	for i := range stripped.Base.Originals {
		stripped.Base.Originals[i].WorldResult = ""
	}
	other, err := Lift(c, stripped, a)
	if err != nil {
		t.Fatal(err)
	}
	for i := range path {
		for j := range path[i].State.Base.Originals {
			path[i].State.Base.Originals[j].WorldResult = ""
		}
	}
	for i := range other {
		for j := range other[i].State.Base.Originals {
			other[i].State.Base.Originals[j].WorldResult = ""
		}
	}
	if !reflect.DeepEqual(path, other) {
		t.Fatal("lift consulted WorldResult")
	}
	end := path[len(path)-1].State
	end = take(t, c, end, "AB", "send-ack").State
	end = take(t, c, end, "AB", "release").State
	if end.Base.Originals[0].Recorded == "" || end.Base.Originals[0].Retained != "" || end.Held != ([2]bool{}) {
		t.Fatal("ACK release confused custody with capacity")
	}
	if Project(c, end).Base.Stages[0] != cc.D {
		t.Fatal("consumption lost")
	}
}

func TestOnePhysicalCrashDisablesBothRolesAndPreservesHold(t *testing.T) {
	c := config(t, "shared-2", "early", 7)
	s := prefix(t, c, bothAdmitted)
	s = take(t, c, s, "AB", "reserve-capacity").State
	s = take(t, c, s, "BC", "reserve-capacity").State
	s = take(t, c, s, "BC", "apply").State
	s = take(t, c, s, "BC", "send-result").State
	s = take(t, c, s, "BC", "receive-result").State
	crashed := take(t, c, s, "B", "crash").State
	if crashed.Held != s.Held || !crashed.Held[0] || crashed.Base.Originals[1].Volatile != "" || crashed.Base.Originals[1].Retained == "" {
		t.Fatal("B crash violated incoming hold or outgoing custody")
	}
	for _, e := range ConcreteSuccessors(c, crashed) {
		if e.Event.Actor == "B" && e.Event.Kind != twoagent.Recover {
			t.Fatal("down owner progressed either role")
		}
	}
}
