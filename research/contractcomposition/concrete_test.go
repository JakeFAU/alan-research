package contractcomposition

import (
	"github.com/JakeFAU/alan-research/research/twoagent"
	"reflect"
	"testing"
)

func take(t *testing.T, c Config, s ConcreteState, id string, k twoagent.Kind) MicroStep {
	t.Helper()
	for _, e := range ConcreteSuccessors(c, s) {
		if e.Event.Kind == k && (e.Event.Original == id || e.Event.Actor == id) && !e.Event.Duplicate {
			return e
		}
	}
	t.Fatalf("no legal %s/%s from %+v", id, k, Project(c, s))
	return MicroStep{}
}
func runPrefix(t *testing.T, c Config, steps [][2]string) ConcreteState {
	t.Helper()
	s, e := NewConcrete(c)
	if e != nil {
		t.Fatal(e)
	}
	for _, p := range steps {
		s = take(t, c, s, p[0], twoagent.Kind(p[1])).State
	}
	return s
}
func canApply(c Config, s ConcreteState, id string) bool {
	for _, e := range ConcreteSuccessors(c, s) {
		if e.Event.Kind == twoagent.Apply && e.Event.Original == id {
			return true
		}
	}
	return false
}

func TestConcreteHiddenDebtAndHiddenWait(t *testing.T) {
	ws, err := BuildWitnesses()
	if err != nil {
		t.Fatal(err)
	}
	byName := map[string]Witness{}
	for _, w := range ws {
		if err := CheckWitness(w); err != nil {
			t.Fatalf("%s: %v", w.Name, err)
		}
		byName[w.Name] = w
	}
	state := func(name string) ConcreteState {
		w, ok := byName[name]
		if !ok {
			t.Fatalf("missing %s", name)
		}
		return w.Steps[len(w.Steps)-1].State
	}
	left, right := state("R0-unpaid"), state("R0-paid")
	if left.Agents != right.Agents || left.Originals[1] != right.Originals[1] {
		t.Fatal("R0 exposed full transport or physical availability differs")
	}
	if byName["R0-unpaid"].Complete || !byName["R0-paid"].Complete || !byName["R0-unpaid"].CanComplete {
		t.Fatal("R0 completion collision missing")
	}
	u, w := state("R1-uncommitted"), state("R1-waiting")
	c := fixture("cycle", 6)
	if u.Agents != w.Agents || u.Originals[1] != w.Originals[1] || u.Originals[2] != w.Originals[2] {
		t.Fatal("R1 exposed state differs")
	}
	if !canApply(c, u, "CA") || canApply(c, w, "CA") {
		t.Fatal("hidden waiting must change concrete CA readiness")
	}
}

func TestPreRecordLossCustodyAndNoReapplication(t *testing.T) {
	c := fixture("chain", 7)
	s := runPrefix(t, c, [][2]string{{"AB", "commit"}, {"AB", "send-request"}, {"AB", "admit"}, {"AB", "apply"}, {"AB", "send-result"}, {"AB", "receive-result"}, {"A", "crash"}})
	e := s.Originals[0]
	if e.Volatile != "" || e.Recorded != "" || e.Retained == "" {
		t.Fatal("pre-record crash must lose only volatile receipt")
	}
	a := Action{"consume", "AB", c.Originals[0]}
	path, err := Lift(c, s, a)
	if err != nil {
		t.Fatal(err)
	}
	for _, step := range path {
		if step.Event.Kind == twoagent.Apply {
			t.Fatal("lifting reapplied effect")
		}
	}
	s = path[len(path)-1].State
	s = take(t, c, s, "AB", twoagent.SendACK).State
	s = take(t, c, s, "AB", twoagent.Release).State
	if !s.Originals[0].Consumed || s.Originals[0].Retained != "" || s.Originals[0].Recorded == "" {
		t.Fatal("release destroyed durable custody")
	}
	if canApply(c, s, "AB") {
		t.Fatal("tombstone lost after release")
	}
}

func TestAbstractEligibilityIsNotImmediateReadiness(t *testing.T) {
	c := fixture("chain", 7)
	s := runPrefix(t, c, [][2]string{{"AB", "commit"}})
	if !ServiceGuard(c, Project(c, s), "B") || canApply(c, s, "AB") {
		t.Fatal("eligibility was confused with admission")
	}
	path, err := Lift(c, s, Action{"apply", "AB", c.Originals[0]})
	if err != nil {
		t.Fatal(err)
	}
	if len(path) != 3 || path[0].Event.Kind != twoagent.SendRequest || path[1].Event.Kind != twoagent.Admit || path[2].Event.Kind != twoagent.Apply {
		t.Fatal("lift must use actual request plumbing")
	}
}

func TestWitnessHostileBindingsReleaseAndDuplicateApplication(t *testing.T) {
	ws, err := BuildWitnesses()
	if err != nil {
		t.Fatal(err)
	}
	for _, mutation := range []string{"binding", "release", "duplicate"} {
		w := ws[0]
		w.Steps = append([]MicroStep{}, w.Steps...)
		switch mutation {
		case "binding":
			w.Steps[0].Event.Binding.Grant = "forged"
		case "release":
			w.Steps[0].Event.Kind = twoagent.Release
		case "duplicate":
			for i, e := range w.Steps {
				if e.Event.Kind == twoagent.Apply {
					w.Steps = append(w.Steps[:i+1], append([]MicroStep{e}, w.Steps[i+1:]...)...)
					break
				}
			}
		}
		if CheckWitness(w) == nil {
			t.Fatalf("accepted %s", mutation)
		}
	}
}

// This is factored edge coverage, not full three-original microstate reachability.
// Each named cycle edge crosses every old state with its callee's four outgoing
// stages, two policies and three physical phases. Terminal chain BC instead
// crosses both C policies and all three C phases with no outgoing original.
func TestEveryOldEdgeRepresentativeProjectionAndLifting(t *testing.T) {
	old := twoagent.Explore(twoagent.DurableACK)
	cases, lifts, arcs := 0, 0, 0
	for _, target := range []struct{ topology, id string }{{"cycle", "AB"}, {"cycle", "BC"}, {"cycle", "CA"}, {"chain", "BC"}} {
		stages := []Stage{U, W, P, D}
		if target.topology == "chain" {
			stages = []Stage{U}
		}
		targetIndex := slot(target.id)
		binding := fixture(target.topology, 0).Originals[targetIndex]
		caller, callee := agentSlot(binding.Caller), agentSlot(binding.Callee)
		for _, node := range old.Nodes {
			for _, servicing := range []bool{false, true} {
				for _, stage := range stages {
					for _, phase := range []Availability{{}, {Down: true, Crashed: true}, {Crashed: true, Activation: 1}} {
						mask := uint8(0)
						if servicing {
							mask = 1 << callee
						}
						c := fixture(target.topology, mask)
						s, _ := NewConcrete(c)
						v := node.State
						s.Agents[caller] = Availability{v.ADown, v.Crashed, v.Activation}
						s.Agents[callee] = phase
						s.Originals[targetIndex] = Transport{Binding: binding, Committed: v.Committed, Admitted: v.BAdmitted, Applied: v.Applied, Consumed: v.AConsumed, Volatile: v.AVolatile, Recorded: v.ARecorded, Retained: v.BResult, RequestChannel: v.RequestChannel, ResultChannel: v.ResultChannel, ACKChannel: v.ACKChannel, WorldResult: v.WorldResult}
						// Canonical retained representative of the callee's own outgoing phase.
						if target.topology == "cycle" {
							e := &s.Originals[callee]
							e.Committed = stage != U
							e.Admitted = stage == P || stage == D
							e.Applied = e.Admitted
							if e.Applied {
								e.Retained = twoagent.R1
								e.WorldResult = twoagent.R1
							}
							if stage == D {
								e.Recorded = twoagent.R1
								e.Consumed = true
							}
						}
						if got := ProjectEdge(c, s, target.id); got != v {
							t.Fatal("old edge renaming lost fields")
						}
						q := Project(c, s)
						expected := U
						if v.Committed {
							expected = W
						}
						if v.Applied {
							expected = P
						}
						if v.AConsumed {
							expected = D
						}
						if q.Stages[targetIndex] != expected {
							t.Fatal("incorrect progress projection")
						}
						for _, step := range ConcreteSuccessors(c, s) {
							arcs++
							r := Project(c, step.State)
							if r != q {
								found := false
								for _, a := range Successors(c, q) {
									if a.State == r && a.Action.Original == step.Event.Original && a.Action.Kind == string(step.Event.Kind) {
										found = true
									}
								}
								if !found {
									t.Fatal("microstep lacks named macro projection")
								}
							}
						}
						for _, a := range Successors(c, q) {
							if a.Action.Original != target.id {
								continue
							}
							path, err := Lift(c, s, a.Action)
							if err != nil {
								t.Fatalf("old edge %v, stage %s phase %+v: %v", v, stage, phase, err)
							}
							if len(path) == 0 {
								t.Fatal("empty lifting path")
							}
							for i, step := range path {
								got := Project(c, step.State)
								want := q
								if i == len(path)-1 {
									want = a.State
								}
								if got != want {
									t.Fatal("lift advances other original or changes progress early")
								}
							}
							lifts++
						}
						// Delete only the inaccessible oracle field. Recoverable successors must agree.
						stripped := s
						for i := range stripped.Originals {
							stripped.Originals[i].WorldResult = ""
						}
						original, next := ConcreteSuccessors(c, s), ConcreteSuccessors(c, stripped)
						for i := range original {
							for j := range original[i].State.Originals {
								original[i].State.Originals[j].WorldResult = ""
							}
						}
						for i := range next {
							for j := range next[i].State.Originals {
								next[i].State.Originals[j].WorldResult = ""
							}
						}
						if !reflect.DeepEqual(original, next) {
							t.Fatal("world oracle interfered with transition enabling or receipt bytes")
						}
						cases++
					}
				}
			}
		}
	}
	t.Logf("old durable-ACK states=%d; factored representatives=%d; micro arcs checked=%d; enabled named-edge lifts=%d", len(old.Nodes), cases, arcs, lifts)
}

func TestLiftRejectsForgedOrLostCustodyRepresentative(t *testing.T) {
	c := fixture("chain", 7)
	original := runPrefix(t, c, [][2]string{{"AB", "commit"}, {"AB", "send-request"}, {"AB", "admit"}, {"AB", "apply"}})
	for _, mutation := range []string{"binding", "lost-custody", "activation"} {
		s := original
		switch mutation {
		case "binding":
			s.Originals[0].Binding.Resource = "counter-C"
		case "lost-custody":
			s.Originals[0].Retained = ""
		case "activation":
			s.Agents[0].Activation = 2
		}
		if _, err := Lift(c, s, Action{"consume", "AB", c.Originals[0]}); err == nil {
			t.Errorf("accepted violated premise %s", mutation)
		}
	}
}

func TestCrashDisablesBothRolesAndPreservesDurableCustody(t *testing.T) {
	c := fixture("chain", 7)
	s := runPrefix(t, c, [][2]string{{"AB", "commit"}, {"AB", "send-request"}, {"AB", "admit"}, {"AB", "apply"}, {"BC", "commit"}, {"BC", "send-request"}, {"BC", "admit"}, {"BC", "apply"}, {"BC", "send-result"}, {"BC", "receive-result"}})
	retained := s.Originals[0].Retained
	s = take(t, c, s, "B", twoagent.Crash).State
	if s.Originals[0].Retained != retained || s.Originals[1].Volatile != "" {
		t.Fatal("one physical crash must clear only B's volatile caller receipt")
	}
	for _, step := range ConcreteSuccessors(c, s) {
		if step.Event.Actor == "B" && step.Event.Kind != twoagent.Recover {
			t.Fatal("B progressed an incoming or outgoing role while down")
		}
	}
}

func TestOldLedgerRejectsReleaseAndReapplication(t *testing.T) {
	ws, err := BuildWitnesses()
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range ws {
		if w.Name != "pre-record-loss-and-release" {
			continue
		}
		events := []twoagent.Event{}
		for _, step := range w.Steps {
			e := step.Event
			events = append(events, twoagent.Event{Kind: e.Kind, Delegation: twoagent.FixedDelegation(), Grant: twoagent.FixedGrant, Result: e.Result, Duplicate: e.Duplicate})
		}
		if err := twoagent.CheckTrace(twoagent.DurableACK, events); err != nil {
			t.Fatal(err)
		}
		for _, kind := range []twoagent.Kind{twoagent.Release, twoagent.Apply} {
			bad := append([]twoagent.Event{}, events[:4]...)
			e := bad[3]
			e.Kind = kind
			if kind == twoagent.Release {
				e.Result = ""
			}
			bad = append(bad, e)
			if twoagent.CheckTrace(twoagent.DurableACK, bad) == nil {
				t.Fatalf("independent old ledger accepted %s", kind)
			}
		}
	}
}

func TestLiftingDoesNotConsultWorldOracle(t *testing.T) {
	c := fixture("chain", 7)
	s := runPrefix(t, c, [][2]string{{"AB", "commit"}, {"AB", "send-request"}, {"AB", "admit"}, {"AB", "apply"}})
	action := Action{"consume", "AB", c.Originals[0]}
	real, err := Lift(c, s, action)
	if err != nil {
		t.Fatal(err)
	}
	for i := range s.Originals {
		s.Originals[i].WorldResult = ""
	}
	erased, err := Lift(c, s, action)
	if err != nil {
		t.Fatalf("inaccessible oracle affected lifting: %v", err)
	}
	for i := range real {
		for j := range real[i].State.Originals {
			real[i].State.Originals[j].WorldResult = ""
		}
	}
	if !reflect.DeepEqual(real, erased) {
		t.Fatal("lifting read oracle for receipt or step selection")
	}
}
