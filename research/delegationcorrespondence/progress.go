package delegationcorrespondence

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	m "github.com/JakeFAU/alan-research/research/authorityservice"
	b "github.com/JakeFAU/alan-research/research/mutualobligations"
)

type ProgressWitness struct {
	State          m.State
	Prefix, Suffix []m.Event
	Audit          m.Audit
}
type ProgressReport struct {
	Config                    m.Config
	Predicate                 string
	States                    int
	Arcs                      int
	SafeComplete              int
	SafeCompletable           int
	PredicateTrue             int
	PredicateMismatches       int
	NonCompletableNoBARefusal int
	RefusedNoGrant            int
	RefusedRetainedGrant      int
	MaxSuffix                 int
	IncompleteDeadEnds        int
	Counterexample            *ProgressWitness
	LongestSuffix             *ProgressWitness
	StateSetSHA256            string
	ArcSetSHA256              string
	HashEncoding              string
	Limits                    []string
}

type progressNode struct {
	state  m.State
	parent int
	event  m.Event
}
type progressArc struct {
	From, To, Successor int
}

func QueryProgress(config m.Config) (ProgressReport, error) {
	out := ProgressReport{Config: config, Limits: []string{"Existential safe suffix in unchanged finite graph; no inevitable scheduling claim", "No all-runtime-state completion claim; handshake, budgets and unknown external outcomes remain additional obligations"}}
	if config.Policy != "S2" || (config.Scenario != "LateGrant" && config.Scenario != "Revocable") {
		return out, fmt.Errorf("unsupported progress query %s", config.Name())
	}
	out.Predicate = "Unsafe==0 and BA first Decision!=Refused"
	if config.Scenario == "Revocable" {
		out.Predicate = "Unsafe==0 and for each AB,BA: Decision!=Refused and (!Revoked or Applied or Permit)"
	}
	nodes := []progressNode{{state: m.Initial(config), parent: -1}}
	ids := map[m.State]int{nodes[0].state: 0}
	arcs := []progressArc{}
	for from := 0; from < len(nodes); from++ {
		next := m.Successors(config, nodes[from].state)
		if len(next) == 0 && !m.Complete(nodes[from].state) {
			out.IncompleteDeadEnds++
		}
		for successor, edge := range next {
			to, ok := ids[edge.State]
			if !ok {
				to = len(nodes)
				ids[edge.State] = to
				nodes = append(nodes, progressNode{edge.State, from, edge.Event})
			}
			arcs = append(arcs, progressArc{from, to, successor})
		}
	}
	incoming := make([][]int, len(nodes))
	for j, a := range arcs {
		incoming[a.To] = append(incoming[a.To], j)
	}
	reach := make([]bool, len(nodes))
	nextArc := make([]int, len(nodes))
	distance := make([]int, len(nodes))
	queue := []int{}
	for j, n := range nodes {
		nextArc[j] = -1
		if n.state.Unsafe == 0 && m.Complete(n.state) {
			reach[j] = true
			queue = append(queue, j)
			out.SafeComplete++
		}
	}
	for q := 0; q < len(queue); q++ {
		to := queue[q]
		for _, ai := range incoming[to] {
			a := arcs[ai]
			if !reach[a.From] && nodes[a.From].state.Unsafe == 0 {
				reach[a.From] = true
				nextArc[a.From] = ai
				distance[a.From] = distance[to] + 1
				queue = append(queue, a.From)
			}
		}
	}
	witness := func(j int) ProgressWitness {
		w := ProgressWitness{State: nodes[j].state}
		for k := j; nodes[k].parent >= 0; k = nodes[k].parent {
			w.Prefix = append(w.Prefix, nodes[k].event)
		}
		for l, r := 0, len(w.Prefix)-1; l < r; l, r = l+1, r-1 {
			w.Prefix[l], w.Prefix[r] = w.Prefix[r], w.Prefix[l]
		}
		for k := j; nextArc[k] >= 0; {
			a := arcs[nextArc[k]]
			w.Suffix = append(w.Suffix, m.Successors(config, nodes[a.From].state)[a.Successor].Event)
			k = a.To
		}
		w.Audit, _ = m.AuditTrace(config, append(append([]m.Event{}, w.Prefix...), w.Suffix...))
		return w
	}
	maxID := 0
	states := make([]m.State, len(nodes))
	for j, n := range nodes {
		states[j] = n.state
		s := n.state
		predicate := s.Unsafe == 0 && s.Authority[b.BA].Decision != m.Refused
		if config.Scenario == "Revocable" {
			predicate = s.Unsafe == 0
			for d := b.AB; d <= b.BA; d++ {
				predicate = predicate && s.Authority[d].Decision != m.Refused && (!s.Revoked || s.Base.Obligations[d].Applied || s.Authority[d].Permit)
			}
		}
		if predicate {
			out.PredicateTrue++
		}
		if predicate != reach[j] {
			out.PredicateMismatches++
		}
		if !reach[j] {
			if out.Counterexample == nil {
				w := witness(j)
				out.Counterexample = &w
			}
			if s.Authority[b.BA].Decision != m.Refused {
				out.NonCompletableNoBARefusal++
			} else if s.Grant == m.NoGrant {
				out.RefusedNoGrant++
			} else if s.Grant == m.RetainedGrant {
				out.RefusedRetainedGrant++
			}
		}
		if distance[j] > out.MaxSuffix {
			out.MaxSuffix = distance[j]
			maxID = j
		}
	}
	w := witness(maxID)
	out.LongestSuffix = &w
	if !w.Audit.ProtocolValid || !w.Audit.AuthoritySafe || !w.Audit.Paid[0] || !w.Audit.Paid[1] {
		return out, fmt.Errorf("independent suffix audit failed")
	}
	out.States = len(nodes)
	out.Arcs = len(arcs)
	out.SafeCompletable = len(queue)
	stateHash, arcHash := sha256.New(), sha256.New()
	stateEncoder, arcEncoder := json.NewEncoder(stateHash), json.NewEncoder(arcHash)
	for _, state := range states {
		if e := stateEncoder.Encode(state); e != nil {
			return out, e
		}
	}
	for _, arc := range arcs {
		if e := arcEncoder.Encode(arc); e != nil {
			return out, e
		}
	}
	out.StateSetSHA256 = hex.EncodeToString(stateHash.Sum(nil))
	out.ArcSetSHA256 = hex.EncodeToString(arcHash.Sum(nil))
	out.HashEncoding = "JSONL in deterministic BFS discovery order; arc records are From/To/Successor-index in unchanged Successors"
	return out, nil
}
