package capacitycomposition

import (
	"fmt"
	cc "github.com/JakeFAU/alan-research/research/contractcomposition"
	"github.com/JakeFAU/alan-research/research/twoagent"
)

type ProbeRow struct {
	ConfigID        string
	Representatives int
	MicroArcs       int
	Lifts           int
	LiftMicroSteps  int
}
type ProbeSummary struct {
	Domain          string
	OldEdgeStates   int
	Configurations  int
	Representatives int
	MicroArcs       int
	Lifts           int
	LiftMicroSteps  int
	Rows            []ProbeRow
}

func edgeStage(s twoagent.State) cc.Stage {
	if s.AConsumed {
		return cc.D
	}
	if s.Applied {
		return cc.P
	}
	if s.Committed {
		return cc.W
	}
	return cc.U
}
func edgeTransport(b cc.Binding, s twoagent.State) cc.Transport {
	return cc.Transport{Binding: b, Committed: s.Committed, Admitted: s.BAdmitted, Applied: s.Applied, Consumed: s.AConsumed, Volatile: s.AVolatile, Recorded: s.ARecorded, Retained: s.BResult, RequestChannel: s.RequestChannel, ResultChannel: s.ResultChannel, ACKChannel: s.ACKChannel, WorldResult: s.WorldResult}
}
func eraseOracle(s ConcreteState) ConcreteState {
	for i := range s.Base.Originals {
		s.Base.Originals[i].WorldResult = ""
	}
	return s
}
func equalWithoutOracle(a, b []MicroStep) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Event != b[i].Event || eraseOracle(a[i].State) != eraseOracle(b[i].State) {
			return false
		}
	}
	return true
}

// ProbeRepresentatives synthesizes a factored domain, not genuine histories or
// the complete concrete product. For every roster config and reachable abstract
// state, each target edge ranges over all 311 old durable-ACK reachable states
// with the matching stage (and admitted if held). Its caller phase is retained;
// the other two owners independently range over initial/down/recovered phases.
// The other edge is canonical: U empty; W admitted; P retained r1; D retained
// and recorded r1. No volatile receipt is invented for the canonical edge.
func ProbeRepresentatives() (ProbeSummary, error) {
	old := twoagent.Explore(twoagent.DurableACK)
	summary := ProbeSummary{Domain: "factored-old-edge-v1", OldEdgeStates: len(old.Nodes), Rows: []ProbeRow{}}
	phases := []cc.Availability{{}, {Down: true, Crashed: true}, {Crashed: true, Activation: 1}}
	for _, layout := range []string{"shared-1", "shared-2", "separate-1"} {
		for _, policy := range []string{"early", "ready-now", "stable-ready"} {
			for mask := uint8(0); mask < 8; mask++ {
				c, err := NewConfig(layout, policy, mask)
				if err != nil {
					return summary, err
				}
				row := ProbeRow{ConfigID: configID(c)}
				for _, node := range Explore(c).Nodes {
					q := node.State
					for target := 0; target < 2; target++ {
						caller := agentSlot(c.Base.Originals[target].Caller)
						for _, oldNode := range old.Nodes {
							v := oldNode.State
							if edgeStage(v) != q.Base.Stages[target] || (q.Held[target] && !v.BAdmitted) {
								continue
							}
							for _, p0 := range phases {
								for _, p1 := range phases {
									s, err := NewConcrete(c)
									if err != nil {
										return summary, err
									}
									s.Held = q.Held
									s.Base.Agents[caller] = cc.Availability{Down: v.ADown, Crashed: v.Crashed, Activation: v.Activation}
									otherOwners := []int{}
									for i := 0; i < 3; i++ {
										if i != caller {
											otherOwners = append(otherOwners, i)
										}
									}
									s.Base.Agents[otherOwners[0]] = p0
									s.Base.Agents[otherOwners[1]] = p1
									s.Base.Originals[target] = edgeTransport(c.Base.Originals[target], v)
									other := 1 - target
									e := &s.Base.Originals[other]
									stage := q.Base.Stages[other]
									e.Committed = stage != cc.U
									e.Admitted = stage != cc.U
									e.Applied = stage == cc.P || stage == cc.D
									e.Consumed = stage == cc.D
									if e.Applied {
										e.Retained = twoagent.R1
										e.WorldResult = twoagent.R1
									}
									if e.Consumed {
										e.Recorded = twoagent.R1
									}
									if err := checkProbe(c, s, q, &row); err != nil {
										return summary, fmt.Errorf("%s state %+v target %d custody %+v availability %+v: %w", row.ConfigID, q, target, v, s.Base.Agents, err)
									}
								}
							}
						}
					}
				}
				summary.Rows = append(summary.Rows, row)
				summary.Configurations++
				summary.Representatives += row.Representatives
				summary.MicroArcs += row.MicroArcs
				summary.Lifts += row.Lifts
				summary.LiftMicroSteps += row.LiftMicroSteps
			}
		}
	}
	return summary, nil
}
func checkProbe(c Config, s ConcreteState, q State, row *ProbeRow) error {
	if err := s.Validate(c); err != nil {
		return err
	}
	if Project(c, s) != q {
		return fmt.Errorf("full projection differs from prescribed state")
	}
	actions := Successors(c, q)
	steps := ConcreteSuccessors(c, s)
	for _, e := range steps {
		if err := e.State.Validate(c); err != nil {
			return err
		}
		target := Project(c, e.State)
		if target != q {
			found := false
			for _, a := range actions {
				if a.State == target && a.Action.Kind == string(e.Event.Kind) && a.Action.Original == e.Event.Original && a.Action.Binding == e.Event.Binding {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("unprojectable concrete step")
			}
		}
	}
	if !equalWithoutOracle(steps, ConcreteSuccessors(c, eraseOracle(s))) {
		return fmt.Errorf("oracle interfered with successor selection")
	}
	row.MicroArcs += len(steps)
	for _, a := range actions {
		path, err := Lift(c, s, a.Action)
		if err != nil {
			return err
		}
		if len(path) == 0 {
			return fmt.Errorf("empty lift")
		}
		cursor := s
		for i, e := range path {
			found := false
			for _, candidate := range ConcreteSuccessors(c, cursor) {
				if candidate == e {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("illegal lift adjacency")
			}
			expected := q
			if i == len(path)-1 {
				expected = a.State
			}
			if Project(c, e.State) != expected {
				return fmt.Errorf("lift changed unrelated progress or holds")
			}
			if err := e.State.Validate(c); err != nil {
				return err
			}
			cursor = e.State
		}
		stripped, err := Lift(c, eraseOracle(s), a.Action)
		if err != nil {
			return err
		}
		if !equalWithoutOracle(path, stripped) {
			return fmt.Errorf("oracle interfered with lifting")
		}
		row.Lifts++
		row.LiftMicroSteps += len(path)
	}
	row.Representatives++
	return nil
}
