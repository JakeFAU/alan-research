package twoagent

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestExploreCustodyAndRecovery(t *testing.T) {
	for _, profile := range []Profile{DeliveryACK, DurableACK} {
		t.Run(string(profile), func(t *testing.T) {
			r := Explore(profile)
			if r.States < 2 || r.Edges < r.States || r.CyclicStates == 0 {
				t.Fatalf("expected complete graph with retry cycles: %+v", r)
			}
			if r.LocalSafetyViolations != 0 || r.LocalTransitionViolations != 0 {
				t.Fatalf("local safety violations: %d", r.LocalSafetyViolations)
			}
			if profile == DeliveryACK {
				if r.CustodyViolations == 0 || r.CannotComplete == 0 || len(r.LossWitnesses) != 2 {
					t.Fatalf("weak profile lacks loss evidence: %+v", r)
				}
				if r.Indistinguishable == nil || !r.Indistinguishable.RecoverableEqual {
					t.Fatal("missing paired indistinguishable worlds")
				}
				pair := r.Indistinguishable
				if pair.Left.State.WorldResult == pair.Right.State.WorldResult || Recoverable(pair.Left.State) != Recoverable(pair.Right.State) {
					t.Fatal("pair differs in recoverable state")
				}
				for _, w := range r.LossWitnesses {
					if w.CanComplete || !w.State.Crashed || w.State.ADown {
						t.Fatal("loss witness must be recovered and unable to complete")
					}
					if err := CheckTrace(profile, w.Trace); err != nil {
						t.Fatalf("invalid witness: %v", err)
					}
				}
			} else if r.CustodyViolations != 0 || r.CannotComplete != 0 || len(r.LossWitnesses) != 0 {
				t.Fatalf("durable ACK contract violated: %+v", r)
			}
		})
	}
}

func TestTraceAdversaries(t *testing.T) {
	id := FixedDelegation()
	event := func(kind Kind, result Value) Event {
		return Event{Kind: kind, Delegation: id, Grant: FixedGrant, Result: result}
	}
	prefix := []Event{event(Commit, ""), event(SendRequest, ""), event(Admit, ""), event(Apply, R1), event(SendResult, R1), event(ReceiveResult, R1)}
	good := append(slices.Clone(prefix), event(Record, R1), event(SendACK, ""), event(Release, ""), event(Crash, ""), event(Recover, ""), event(Consume, R1))
	if err := CheckTrace(DurableACK, good); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		profile Profile
		trace   []Event
	}{
		{"fabricated_result", DeliveryACK, append(slices.Clone(prefix[:5]), event(ReceiveResult, R2))},
		{"early_durable_ack", DurableACK, append(slices.Clone(prefix), event(SendACK, ""))},
		{"duplicate_application", DeliveryACK, append(slices.Clone(prefix[:4]), event(Apply, R1))},
		{"consume_before_record", DeliveryACK, append(slices.Clone(prefix), event(Consume, R1))},
		{"record_without_receipt", DeliveryACK, append(slices.Clone(prefix[:4]), event(Record, R1))},
		{"release_without_ack", DeliveryACK, append(slices.Clone(prefix[:4]), event(Release, ""))},
		{"send_without_retention", DeliveryACK, []Event{event(SendResult, R1)}},
		{"resend_after_release", DeliveryACK, append(slices.Clone(good[:9]), event(SendResult, R1))},
		{"wrong_record_binding", DeliveryACK, append(slices.Clone(prefix), event(Record, R2))},
		{"wrong_consumption_binding", DurableACK, append(slices.Clone(good[:7]), event(Consume, R2))},
		{"ack_payload_leak", DeliveryACK, append(slices.Clone(prefix), event(SendACK, R1))},
		{"consume_twice", DurableACK, append(slices.Clone(good), event(Consume, R1))},
	}
	changed := slices.Clone(good)
	changed[6].Delegation.Operation++
	tests = append(tests, struct {
		name    string
		profile Profile
		trace   []Event
	}{"changed_correlation", DurableACK, changed})
	changedGrant := slices.Clone(good)
	changedGrant[3].Grant = "forged"
	tests = append(tests, struct {
		name    string
		profile Profile
		trace   []Event
	}{"changed_grant", DurableACK, changedGrant})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckTrace(tt.profile, tt.trace); err == nil {
				t.Fatal("adversarial trace accepted")
			}
		})
	}
}

func TestAsynchronousBAndLateReceipt(t *testing.T) {
	id := FixedDelegation()
	e := func(k Kind, v Value) Event { return Event{Kind: k, Delegation: id, Grant: FixedGrant, Result: v} }
	trace := []Event{e(Commit, ""), e(SendRequest, ""), e(Crash, ""), e(Admit, ""), e(Apply, R2), e(SendResult, R2), e(Recover, ""), e(ReceiveResult, R2), e(Record, R2), e(Consume, R2), e(SendACK, ""), e(Release, "")}
	if err := CheckTrace(DurableACK, trace); err != nil {
		t.Fatal(err)
	}
	// A genuine retry/resend survives a pre-record crash and does not reapply.
	retry := []Event{e(Commit, ""), e(SendRequest, ""), e(Admit, ""), e(Apply, R1), e(SendResult, R1), e(ReceiveResult, R1), e(Crash, ""), e(SendResult, R1), e(Recover, ""), e(SendRequest, ""), e(Deduplicate, ""), e(ReceiveResult, R1), e(Record, R1), e(Consume, R1)}
	if err := CheckTrace(DurableACK, retry); err != nil {
		t.Fatal(err)
	}
	s := Initial()
	for _, event := range trace {
		var found bool
		for _, edge := range Successors(DurableACK, s) {
			if edge.Event == event {
				s = edge.State
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("model omits asynchronous event %s from %+v", event.Kind, s)
		}
	}
}

func TestEveryReachableEdgeHasIndependentTrace(t *testing.T) {
	for _, p := range []Profile{DeliveryACK, DurableACK} {
		t.Run(string(p), func(t *testing.T) {
			r := Explore(p)
			for i, node := range r.Nodes {
				if err := CheckTrace(p, r.Trace(i)); err != nil {
					t.Fatalf("state %d: %v", i, err)
				}
				for _, edge := range node.Edges {
					trace := append(r.Trace(i), edge.Event)
					if err := CheckTrace(p, trace); err != nil {
						t.Fatalf("edge %d -> %d: %v", i, edge.To, err)
					}
				}
			}
		})
	}
}

func TestDeterministicEvidence(t *testing.T) {
	a, _ := json.Marshal(Explore(DeliveryACK))
	b, _ := json.Marshal(Explore(DeliveryACK))
	if string(a) != string(b) {
		t.Fatal("export is nondeterministic")
	}
}

func TestOracleFactCannotEnableOrSupplySteps(t *testing.T) {
	for _, p := range []Profile{DeliveryACK, DurableACK} {
		t.Run(string(p), func(t *testing.T) {
			for _, node := range Explore(p).Nodes {
				a, b := node.State, node.State
				a.WorldResult = R1
				b.WorldResult = R2
				ae, be := Successors(p, a), Successors(p, b)
				if len(ae) != len(be) {
					t.Fatal("oracle changed enabled steps")
				}
				for i := range ae {
					if ae[i].Event != be[i].Event || Recoverable(ae[i].State) != Recoverable(be[i].State) {
						t.Fatal("oracle supplied agent-visible data")
					}
				}
			}
		})
	}
}

func TestBDoesNotInspectAPrivateState(t *testing.T) {
	for _, p := range []Profile{DeliveryACK, DurableACK} {
		t.Run(string(p), func(t *testing.T) {
			for _, node := range Explore(p).Nodes {
				altered := node.State
				altered.ADown = !altered.ADown
				altered.ARecorded = R2
				altered.AVolatile = R1
				altered.AConsumed = !altered.AConsumed
				bEvents := func(s State) []Event {
					var out []Event
					for _, edge := range Successors(p, s) {
						switch edge.Event.Kind {
						case Admit, Deduplicate, Apply, SendResult, Release:
							out = append(out, edge.Event)
						}
					}
					return out
				}
				if !slices.Equal(bEvents(node.State), bEvents(altered)) {
					t.Fatal("B read A's private state")
				}
			}
		})
	}
}

func TestCommittedFactsSurviveEveryEdge(t *testing.T) {
	for _, p := range []Profile{DeliveryACK, DurableACK} {
		t.Run(string(p), func(t *testing.T) {
			r := Explore(p)
			for _, node := range r.Nodes {
				for _, edge := range node.Edges {
					a, b := node.State, r.Nodes[edge.To].State
					if a.Delegation != b.Delegation || a.Grant != b.Grant || (a.Committed && !b.Committed) || (a.BAdmitted && !b.BAdmitted) || (a.Applied && !b.Applied) || (a.AConsumed && !b.AConsumed) || (a.Crashed && !b.Crashed) || (a.ARecorded != "" && a.ARecorded != b.ARecorded) || (a.WorldResult != "" && a.WorldResult != b.WorldResult) {
						t.Fatalf("committed fact lost on %s", edge.Event.Kind)
					}
				}
			}
		})
	}
}

func TestShortestWitnessAndGraphClosure(t *testing.T) {
	for _, p := range []Profile{DeliveryACK, DurableACK} {
		t.Run(string(p), func(t *testing.T) {
			r := Explore(p)
			ids := make(map[State]int, len(r.Nodes))
			for i, n := range r.Nodes {
				ids[n.State] = i
			}
			for i, n := range r.Nodes {
				edges := Successors(p, n.State)
				if len(edges) != len(n.Edges) {
					t.Fatal("explorer discarded transitions")
				}
				for j, e := range edges {
					to, ok := ids[e.State]
					if !ok || n.Edges[j].To != to || n.Edges[j].Event != e.Event {
						t.Fatal("graph is not closed under successors")
					}
					if len(r.Trace(to)) > len(r.Trace(i))+1 {
						t.Fatal("BFS witness is not shortest")
					}
				}
			}
			for _, w := range r.LossWitnesses {
				if len(w.Trace) != 10 {
					t.Fatalf("expected ten-event minimal recovered-loss schedule, got %d", len(w.Trace))
				}
			}
		})
	}
}

func TestCanonicalLossAndDurableCrashBoundaries(t *testing.T) {
	for _, result := range []Value{R1, R2} {
		t.Run(string(result), func(t *testing.T) {
			e := func(k Kind, v Value) Event {
				return Event{Kind: k, Delegation: FixedDelegation(), Grant: FixedGrant, Result: v}
			}
			receipt := []Event{e(Commit, ""), e(SendRequest, ""), e(Admit, ""), e(Apply, result), e(SendResult, result), e(ReceiveResult, result)}
			loss := append(slices.Clone(receipt), e(SendACK, ""), e(Release, ""), e(Crash, ""), e(Recover, ""))
			if err := CheckTrace(DeliveryACK, loss); err != nil {
				t.Fatal(err)
			}
			s := follow(t, DeliveryACK, loss)
			if !s.Applied || s.ADown || s.AVolatile != "" || s.ARecorded != "" || s.BResult != "" || s.ResultChannel != "" {
				t.Fatalf("canonical trace has not lost result: %+v", s)
			}
			for _, consumeBeforeCrash := range []bool{false, true} {
				trace := append(slices.Clone(receipt), e(Record, result), e(SendACK, ""), e(Release, ""))
				if consumeBeforeCrash {
					trace = append(trace, e(Consume, result))
				}
				trace = append(trace, e(Crash, ""), e(Recover, ""))
				if !consumeBeforeCrash {
					trace = append(trace, e(Consume, result))
				}
				if err := CheckTrace(DurableACK, trace); err != nil {
					t.Fatal(err)
				}
				if end := follow(t, DurableACK, trace); end.ARecorded != result || !end.AConsumed {
					t.Fatal("durable facts lost across crash")
				}
			}
		})
	}
}

func follow(t *testing.T, p Profile, trace []Event) State {
	t.Helper()
	s := Initial()
	for _, e := range trace {
		found := false
		for _, edge := range Successors(p, s) {
			if edge.Event == e {
				s = edge.State
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("event %s is not enabled", e.Kind)
		}
	}
	return s
}

func TestLocalTransitionCheckerRejectsLostFacts(t *testing.T) {
	s := Initial()
	s.Committed = true
	s.BAdmitted = true
	s.Applied = true
	s.WorldResult = R1
	s.BResult = R1
	s.ARecorded = R1
	s.AConsumed = true
	s.Crashed = true
	s.Activation = 1
	tests := []struct {
		name   string
		change func(*State)
	}{
		{"commit", func(n *State) { n.Committed = false }},
		{"admission", func(n *State) { n.BAdmitted = false }},
		{"tombstone", func(n *State) { n.Applied = false }},
		{"record", func(n *State) { n.ARecorded = "" }},
		{"consumption", func(n *State) { n.AConsumed = false }},
		{"world", func(n *State) { n.WorldResult = R2 }},
		{"identity", func(n *State) { n.Delegation.Operation++ }},
		{"grant", func(n *State) { n.Grant = "changed" }},
		{"activation", func(n *State) { n.Activation = 0 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := s
			tt.change(&next)
			if localTransitionSafe(s, Edge{Event: Event{Kind: Release}, State: next}) {
				t.Fatal("lost committed fact accepted")
			}
		})
	}
}

func TestExistentialCompletionAllowsUnfairRetryCycle(t *testing.T) {
	e := func(k Kind) Event { return Event{Kind: k, Delegation: FixedDelegation(), Grant: FixedGrant} }
	prefix := []Event{e(Commit), e(SendRequest), e(Admit)}
	before := follow(t, DurableACK, prefix)
	after := follow(t, DurableACK, append(slices.Clone(prefix), e(SendRequest), e(Deduplicate)))
	if before != after || before.AConsumed {
		t.Fatal("expected unconsumed request retry cycle")
	}
	if err := CheckTrace(DurableACK, append(slices.Clone(prefix), e(SendRequest), e(Deduplicate), e(SendRequest), e(Deduplicate))); err != nil {
		t.Fatal(err)
	}
}
