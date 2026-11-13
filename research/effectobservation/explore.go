package effectobservation

// Arc is an outgoing transition whose destination indexes Graph.States.
type Arc struct {
	Event Event `json:"event"`
	To    int   `json:"to"`
}

// Graph indexes the same BFS order across states, outgoing arcs and one
// shortest event prefix per state. States[0] is Initial.
type Graph struct {
	States   []State   `json:"states"`
	Arcs     [][]Arc   `json:"arcs"`
	Prefixes [][]Event `json:"prefixes"`
}

// Explore enumerates the reachable graph without a depth limit.
func Explore() Graph {
	initial := Initial()
	graph := Graph{States: []State{initial}, Prefixes: [][]Event{{}}}
	ids := map[State]int{initial: 0}

	for from := 0; from < len(graph.States); from++ {
		arcs := []Arc{}
		for _, edge := range Successors(graph.States[from]) {
			to, seen := ids[edge.State]
			if !seen {
				to = len(graph.States)
				ids[edge.State] = to
				graph.States = append(graph.States, edge.State)

				parent := graph.Prefixes[from]
				prefix := make([]Event, len(parent)+1)
				copy(prefix, parent)
				prefix[len(parent)] = edge.Event
				graph.Prefixes = append(graph.Prefixes, prefix)
			}
			arcs = append(arcs, Arc{Event: edge.Event, To: to})
		}
		graph.Arcs = append(graph.Arcs, arcs)
	}
	return graph
}

// Replay checks complete descriptors against the operational transitions,
// starting at Initial. Independent checking must use a separate interpreter.
func Replay(events []Event) (State, bool) {
	state := Initial()
	for _, event := range events {
		matched := false
		for _, edge := range Successors(state) {
			if edge.Event == event {
				state = edge.State
				matched = true
				break
			}
		}
		if !matched {
			return State{}, false
		}
	}
	return state, true
}
