package effectobservation_test

import (
	"reflect"
	"testing"

	m "github.com/JakeFAU/alan-research/research/effectobservation"
)

func TestCompleteGraph(t *testing.T) {
	graph := m.Explore()
	if !reflect.DeepEqual(graph, m.Explore()) {
		t.Fatal("enumeration changed between calls")
	}
	count := len(graph.States)
	if count == 0 || len(graph.Arcs) != count || len(graph.Prefixes) != count ||
		graph.States[0] != m.Initial() || len(graph.Prefixes[0]) != 0 {
		t.Fatal("graph does not start at the admitted cut or its indices are incomplete")
	}

	ids := make(map[m.State]int, count)
	for id, state := range graph.States {
		if previous, exists := ids[state]; exists {
			t.Fatalf("state %d duplicates state %d", id, previous)
		}
		ids[state] = id
		checkState(t, id, state)
		prefix := graph.Prefixes[id]
		if id > 0 && len(prefix) < len(graph.Prefixes[id-1]) {
			t.Fatalf("state %d appears before a shorter prefix", id-1)
		}
		if actual, ok := m.Replay(prefix); !ok || actual != state {
			t.Fatalf("prefix %d does not reconstruct its state: %+v", id, actual)
		}
	}

	nextID := 1
	for from, arcs := range graph.Arcs {
		state := graph.States[from]
		edges := m.Successors(state)
		if len(arcs) != len(edges) {
			t.Fatalf("state %d has %d arcs for %d successors", from, len(arcs), len(edges))
		}
		events := map[m.Event]bool{}
		for i, arc := range arcs {
			edge := edges[i]
			to, exists := ids[edge.State]
			if !exists || arc.To != to || arc.Event != edge.Event || events[arc.Event] {
				t.Fatalf("invalid or repeated arc at state %d: %+v", from, arc)
			}
			events[arc.Event] = true
			checkEdge(t, state, edge)
			fromPrefix, toPrefix := graph.Prefixes[from], graph.Prefixes[to]
			if len(toPrefix) > len(fromPrefix)+1 {
				t.Fatalf("arc %d to %d gives a shorter prefix", from, to)
			}
			if to >= nextID {
				if to != nextID || len(toPrefix) != len(fromPrefix)+1 ||
					!reflect.DeepEqual(toPrefix[:len(fromPrefix)], fromPrefix) || toPrefix[len(fromPrefix)] != arc.Event {
					t.Fatalf("first discovery of state %d does not preserve BFS order and prefix", to)
				}
				nextID++
			}
		}
	}
	if nextID != count {
		t.Fatalf("BFS first discovered %d of %d states", nextID, count)
	}
}

func checkState(t *testing.T, id int, state m.State) {
	t.Helper()
	initial := m.Initial()
	if state.Binding != initial.Binding || state.AcceptedGrant != initial.AcceptedGrant ||
		state.FirstAllow != initial.FirstAllow || state.Debit != initial.Debit ||
		state.FirstGateAllowed != initial.FirstGateAllowed || state.ControlRoot != initial.ControlRoot {
		t.Fatalf("state %d altered immutable admission or c authority: %+v", id, state)
	}
	switch state.Owner {
	case m.OwnerOriginal:
		if state.Activation != 1 || state.Permit == state.Launched {
			t.Fatalf("state %d has inconsistent original owner custody: %+v", id, state)
		}
	case m.OwnerStopped, m.OwnerReopened:
		if state.Permit || state.Response != m.NoResult {
			t.Fatalf("state %d revived volatile custody: %+v", id, state)
		}
		expectedActivation := uint8(1)
		if state.Owner == m.OwnerReopened {
			expectedActivation = 2
		}
		if state.Activation != expectedActivation {
			t.Fatalf("state %d has the wrong shared activation: %+v", id, state)
		}
	default:
		t.Fatalf("state %d has an unknown owner stage: %+v", id, state)
	}
	if state.Ingress && !state.Launched || state.Applied != m.NoResult && !state.Ingress {
		t.Fatalf("state %d skipped transport or ingress: %+v", id, state)
	}
	switch state.Applied {
	case m.NoResult, m.Value1, m.Value2:
	default:
		t.Fatalf("state %d invented endpoint value %q", id, state.Applied)
	}
	if state.Response != m.NoResult && (state.Response != state.Applied || state.Retained != m.NoResult) ||
		state.Retained != m.NoResult && state.Retained != state.Applied ||
		state.Consumed != m.NoResult && state.Consumed != state.Retained {
		t.Fatalf("state %d has unauthentic or overlapping result custody: %+v", id, state)
	}
	if state.FreshRefused && (state.Owner != m.OwnerReopened || !state.BRevoked || state.Retained != m.NoResult) {
		t.Fatalf("state %d has a premature fresh refusal: %+v", id, state)
	}
}

func checkEdge(t *testing.T, state m.State, edge m.Edge) {
	t.Helper()
	want := state
	event := m.Event{Kind: edge.Event.Kind, Binding: state.Binding, Incarnation: state.Activation}
	switch event.Kind {
	case m.Launch:
		want.Permit, want.Launched = false, true
		event.Actor = m.ActorA
	case m.Ingress:
		want.Ingress = true
		event.Actor, event.Incarnation = m.ActorEndpoint, 0
	case m.Apply:
		want.Applied = edge.Event.Value
		event.Actor, event.Incarnation, event.Value = m.ActorEndpoint, 0, want.Applied
		if state.Applied != m.NoResult || want.Applied == m.NoResult {
			t.Fatalf("application replaced or omitted the world result: %+v", edge)
		}
	case m.Response:
		want.Response = state.Applied
		event.Actor, event.Value = m.ActorA, state.Applied
	case m.Record:
		want.Response, want.Retained = m.NoResult, state.Response
		event.Actor, event.Value = m.ActorA, state.Response
	case m.Consume:
		want.Consumed = state.Retained
		event.Actor, event.Value = m.ActorB, state.Retained
	case m.RevokeB:
		want.BRevoked = true
		event.Actor, event.Root = m.ActorB, m.RootB
	case m.LoseOwner:
		want.Owner, want.Permit, want.Response = m.OwnerStopped, false, m.NoResult
		event.Actor = m.ActorOwner
	case m.Reopen:
		want.Owner, want.Activation = m.OwnerReopened, 2
		event.Actor, event.Incarnation = m.ActorOwner, 2
	case m.FreshRefuse:
		want.FreshRefused = true
		event.Actor, event.Root = m.ActorA, m.RootB
	default:
		t.Fatalf("unexpected event kind %q", event.Kind)
	}
	if state == want || want != edge.State || event != edge.Event {
		t.Fatalf("event should change only its own stage and carry its exact descriptor: got %+v want %+v",
			edge, m.Edge{Event: event, State: want})
	}
}

func TestReplayRejectsMalformedDescriptors(t *testing.T) {
	mutations := []struct {
		name   string
		mutate func(*m.Event)
	}{
		{name: "kind", mutate: func(event *m.Event) { event.Kind = "lookup" }},
		{name: "original", mutate: func(event *m.Event) { event.Binding.Original = "AB/1" }},
		{name: "generation", mutate: func(event *m.Event) { event.Binding.SourceGeneration++ }},
		{name: "revision", mutate: func(event *m.Event) { event.Binding.SourceRevision++ }},
		{name: "grant", mutate: func(event *m.Event) { event.Binding.Grant = "other-grant" }},
		{name: "bound-root", mutate: func(event *m.Event) { event.Binding.Root = m.RootC }},
		{name: "bound-activation", mutate: func(event *m.Event) { event.Binding.OriginalActivation++ }},
		{name: "wire", mutate: func(event *m.Event) { event.Binding.Wire = "other-wire" }},
		{name: "actor", mutate: func(event *m.Event) { event.Actor = "" }},
		{name: "incarnation", mutate: func(event *m.Event) { event.Incarnation++ }},
		{name: "target-root", mutate: func(event *m.Event) { event.Root = m.RootC }},
		{name: "value", mutate: func(event *m.Event) {
			if event.Value == m.NoResult {
				event.Value = m.Value1
				return
			}
			event.Value = "invented"
		}},
	}
	graph := m.Explore()
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			for id, arcs := range graph.Arcs {
				for _, arc := range arcs {
					event := arc.Event
					mutation.mutate(&event)
					path := append(append([]m.Event{}, graph.Prefixes[id]...), event)
					if _, ok := m.Replay(path); ok {
						t.Fatalf("accepted malformed %s event from state %d: %+v", mutation.name, id, event)
					}
				}
			}
		})
	}
}
