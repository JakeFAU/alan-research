package effectobservation

import (
	"fmt"
	"slices"
)

// Witness is an independently audited path from the post-admission initial cut.
type Witness struct {
	StateID int         `json:"state_id"`
	Events  []Event     `json:"events"`
	Store   StoreView   `json:"store"`
	Dossier DossierView `json:"dossier"`
}

type WitnessPair struct {
	Left  Witness `json:"left"`
	Right Witness `json:"right"`
}

type DelayedApplicationWitness struct {
	Cut   Witness `json:"cut"`
	Later Witness `json:"later"`
}

type AuditReport struct {
	CheckedStates      int                       `json:"checked_states"`
	CheckedArcs        int                       `json:"checked_arcs"`
	CheckedPrefixes    int                       `json:"checked_prefixes"`
	MergedTargets      int                       `json:"merged_targets"`
	MergeArcs          int                       `json:"merge_arcs"`
	StoreClasses       int                       `json:"store_classes"`
	DossierClasses     int                       `json:"dossier_classes"`
	SearchCriterion    string                    `json:"search_criterion"`
	ZeroOne            WitnessPair               `json:"zero_one"`
	DistinctValues     WitnessPair               `json:"distinct_values"`
	DelayedApplication DelayedApplicationWitness `json:"delayed_application"`
	RetainedRecovery   Witness                   `json:"retained_recovery"`
}

// AuditGraph checks prefixes, all outgoing edges, edge completeness, merged
// targets and projections against the separate ledger. It does not assume that
// graph state IDs or supplied prefixes encode a particular BFS tie-break.
func AuditGraph(graph Graph) (AuditReport, error) {
	var report AuditReport
	n := len(graph.States)
	if n == 0 || len(graph.Arcs) != n || len(graph.Prefixes) != n {
		return report, fmt.Errorf("graph arrays must have equal nonzero lengths")
	}
	ledgers := make([]Ledger, n)
	seenStates := map[State]bool{}
	stores, dossiers := map[StoreView]bool{}, map[DossierView]bool{}
	for i, state := range graph.States {
		l, err := AuditTrace(graph.Prefixes[i])
		if err != nil {
			return report, fmt.Errorf("prefix %d: %w", i, err)
		}
		if seenStates[state] {
			return report, fmt.Errorf("duplicate state %d", i)
		}
		seenStates[state] = true
		if !ledgerMatchesState(l, state) {
			return report, fmt.Errorf("prefix %d ledger/state mismatch", i)
		}
		if l.StoreView() != ObserveStore(state) || l.DossierView() != ObserveDossier(state) {
			return report, fmt.Errorf("state %d observation mismatch", i)
		}
		stores[l.StoreView()], dossiers[l.DossierView()] = true, true
		ledgers[i] = l
	}
	if len(graph.Prefixes[0]) != 0 {
		return report, fmt.Errorf("state zero must be the initial cut")
	}
	indegree := make([]int, n)
	for from, arcs := range graph.Arcs {
		seen := map[Event]bool{}
		for _, arc := range arcs {
			if arc.To < 0 || arc.To >= n {
				return report, fmt.Errorf("state %d target out of bounds", from)
			}
			if seen[arc.Event] {
				return report, fmt.Errorf("state %d duplicate outgoing event", from)
			}
			seen[arc.Event] = true
			target, err := ledgers[from].accept(arc.Event)
			if err != nil {
				return report, fmt.Errorf("edge %d -> %d: %w", from, arc.To, err)
			}
			if target != ledgers[arc.To] {
				return report, fmt.Errorf("edge %d -> %d merged ledger mismatch", from, arc.To)
			}
			indegree[arc.To]++
			report.CheckedArcs++
		}
		// Exhaust the independently specified descriptor domain, not Successors.
		for _, candidate := range ledgerCandidates() {
			if _, err := ledgers[from].accept(candidate); err == nil && !seen[candidate] {
				return report, fmt.Errorf("state %d missing enabled %s event", from, candidate.Kind)
			}
		}
	}
	paths, order := shortestPaths(graph)
	if len(order) != n {
		return report, fmt.Errorf("graph contains unreachable states")
	}
	for id, prefix := range graph.Prefixes {
		if len(prefix) != len(paths[id]) {
			return report, fmt.Errorf("prefix %d is not shortest", id)
		}
		current := 0
		for _, event := range prefix {
			found := false
			for _, arc := range graph.Arcs[current] {
				if arc.Event == event {
					current, found = arc.To, true
					break
				}
			}
			if !found {
				return report, fmt.Errorf("prefix %d does not follow graph arcs", id)
			}
		}
		if current != id {
			return report, fmt.Errorf("prefix %d ends at wrong state", id)
		}
		if indegree[id] > 1 {
			report.MergedTargets++
			report.MergeArcs += indegree[id] - 1
		}
	}
	report.CheckedStates, report.CheckedPrefixes = n, n
	report.StoreClasses, report.DossierClasses = len(stores), len(dossiers)
	report.SearchCriterion = "From the post-admission initial cut: minimum total event count for pairs, then left length, then left/right event-order lexicographic paths; individual witnesses minimize length then event-order lexicographic path. Event order: launch, ingress, apply(v1), apply(v2), response, record, consume, revoke-b, lose-owner, reopen, fresh-refuse. Negative cut must be stopped/revoked with packet pending ingress; its extension ingresses/applies while stopped. Recovery ends in consume from a reopened/revoked source. Supplied shortest target prefixes alone do not determine history-qualified witnesses."
	makeWitness := func(id int, events []Event) Witness {
		return Witness{StateID: id, Events: events, Store: ledgers[id].StoreView(), Dossier: ledgers[id].DossierView()}
	}
	qualified := func(l Ledger) bool {
		return l.Restarted && l.Revoked && l.Refused && l.Stored == NoResult && l.Used == NoResult
	}
	findPair := func(leftValue, rightValue Result) (WitnessPair, bool) {
		var pair WitnessPair
		found := false
		for _, left := range order {
			if !qualified(ledgers[left]) || ledgers[left].World != leftValue {
				continue
			}
			for _, right := range order {
				if !qualified(ledgers[right]) || ledgers[right].World != rightValue || ledgers[left].StoreView() != ledgers[right].StoreView() {
					continue
				}
				candidate := WitnessPair{makeWitness(left, paths[left]), makeWitness(right, paths[right])}
				if !found || pairLess(candidate, pair) {
					pair, found = candidate, true
				}
			}
		}
		return pair, found
	}
	var ok bool
	report.ZeroOne, ok = findPair(NoResult, Value1)
	if !ok {
		return AuditReport{}, fmt.Errorf("missing equal-Store zero/one pair")
	}
	report.DistinctValues, ok = findPair(Value1, Value2)
	if !ok {
		return AuditReport{}, fmt.Errorf("missing distinct unretained-value pair")
	}
	foundDelayed, foundRecovery := false, false
	for _, from := range order {
		l := ledgers[from]
		if l.Dead && l.Revoked && l.PacketSent && !l.Contact && l.World == NoResult {
			for _, ingress := range graph.Arcs[from] {
				if ingress.Event.Kind != Ingress {
					continue
				}
				for _, apply := range graph.Arcs[ingress.To] {
					if apply.Event.Kind != Apply {
						continue
					}
					path := appendEvent(appendEvent(paths[from], ingress.Event), apply.Event)
					if !foundDelayed || pathLess(path, report.DelayedApplication.Later.Events) {
						report.DelayedApplication = DelayedApplicationWitness{makeWitness(from, paths[from]), makeWitness(apply.To, path)}
						foundDelayed = true
					}
				}
			}
		}
		if l.Restarted && l.Revoked && l.Stored != NoResult && l.Used == NoResult {
			for _, arc := range graph.Arcs[from] {
				if arc.Event.Kind != Consume {
					continue
				}
				path := appendEvent(paths[from], arc.Event)
				if !foundRecovery || pathLess(path, report.RetainedRecovery.Events) {
					report.RetainedRecovery, foundRecovery = makeWitness(arc.To, path), true
				}
			}
		}
	}
	if !foundDelayed || !foundRecovery {
		return AuditReport{}, fmt.Errorf("missing history-qualified witness")
	}
	return report, nil
}

func ledgerMatchesState(l Ledger, s State) bool {
	v := l.StoreView()
	return s == (State{Binding: v.Binding, AcceptedGrant: v.AcceptedGrant, FirstAllow: v.FirstAllow,
		Debit: v.Debit, FirstGateAllowed: v.FirstGateAllowed, Owner: v.Owner, Activation: v.Activation,
		BRevoked: v.BRevoked, ControlRoot: v.ControlRoot, FreshRefused: v.FreshRefused,
		Permit: l.Opportunity, Launched: l.PacketSent, Ingress: l.Contact, Applied: l.World,
		Response: l.Inbox, Retained: l.Stored, Consumed: l.Used})
}

// The only valid descriptor shapes in the frozen domain. Incarnations and
// values are exhausted rather than derived from operational state.
func ledgerCandidates() []Event {
	var events []Event
	for _, kind := range []EventKind{Launch, Ingress, Apply, Response, Record, Consume, RevokeB, LoseOwner, Reopen, FreshRefuse} {
		actor, root := ActorA, Root("")
		switch kind {
		case Ingress, Apply:
			actor = ActorEndpoint
		case Consume, RevokeB:
			actor = ActorB
		case LoseOwner, Reopen:
			actor = ActorOwner
		}
		if kind == RevokeB || kind == FreshRefuse {
			root = RootB
		}
		for inc := uint8(0); inc <= 2; inc++ {
			for _, value := range []Result{NoResult, Value1, Value2} {
				events = append(events, Event{Kind: kind, Binding: ledgerBinding(), Actor: actor, Incarnation: inc, Root: root, Value: value})
			}
		}
	}
	return events
}

func appendEvent(prefix []Event, event Event) []Event {
	result := make([]Event, len(prefix)+1)
	copy(result, prefix)
	result[len(prefix)] = event
	return result
}

func eventRank(e Event) int {
	for i, kind := range []EventKind{Launch, Ingress, Apply, Response, Record, Consume, RevokeB, LoseOwner, Reopen, FreshRefuse} {
		if e.Kind == kind {
			if e.Value == Value2 {
				return i*2 + 1
			}
			return i * 2
		}
	}
	return 100
}

func pathLexLess(a, b []Event) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if eventRank(a[i]) != eventRank(b[i]) {
			return eventRank(a[i]) < eventRank(b[i])
		}
	}
	return len(a) < len(b)
}

func pathLess(a, b []Event) bool {
	if len(a) != len(b) {
		return len(a) < len(b)
	}
	return pathLexLess(a, b)
}

func pairLess(a, b WitnessPair) bool {
	aLen, bLen := len(a.Left.Events)+len(a.Right.Events), len(b.Left.Events)+len(b.Right.Events)
	if aLen != bLen {
		return aLen < bLen
	}
	if len(a.Left.Events) != len(b.Left.Events) {
		return len(a.Left.Events) < len(b.Left.Events)
	}
	if !slices.Equal(a.Left.Events, b.Left.Events) {
		return pathLexLess(a.Left.Events, b.Left.Events)
	}
	return pathLexLess(a.Right.Events, b.Right.Events)
}

func shortestPaths(graph Graph) ([][]Event, []int) {
	paths := make([][]Event, len(graph.States))
	paths[0] = []Event{}
	order := []int{0}
	for pos := 0; pos < len(order); pos++ {
		from := order[pos]
		arcs := slices.Clone(graph.Arcs[from])
		slices.SortFunc(arcs, func(a, b Arc) int { return eventRank(a.Event) - eventRank(b.Event) })
		for _, arc := range arcs {
			if paths[arc.To] == nil {
				paths[arc.To] = appendEvent(paths[from], arc.Event)
				order = append(order, arc.To)
			}
		}
	}
	return paths, order
}
