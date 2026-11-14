package contractcomposition

import (
	"fmt"
	"github.com/JakeFAU/alan-research/research/twoagent"
)

type Availability struct {
	Down       bool
	Crashed    bool
	Activation uint8
}
type Transport struct {
	Binding        Binding
	Committed      bool
	Admitted       bool
	Applied        bool
	Consumed       bool
	Volatile       twoagent.Value
	Recorded       twoagent.Value
	Retained       twoagent.Value
	RequestChannel bool
	ResultChannel  twoagent.Value
	ACKChannel     bool
	WorldResult    twoagent.Value
}
type ConcreteState struct {
	Agents    [3]Availability
	Originals [3]Transport
}
type Event struct {
	Kind      twoagent.Kind
	Actor     string
	Original  string
	Binding   Binding
	Result    twoagent.Value
	Duplicate bool
}
type MicroStep struct {
	Event Event
	State ConcreteState
}
type Witness struct {
	Name        string
	Config      Config
	Initial     ConcreteState
	Steps       []MicroStep
	Projection  []Action
	Complete    bool
	CanComplete bool
}

func NewConcrete(c Config) (ConcreteState, error) {
	if err := c.Validate(); err != nil {
		return ConcreteState{}, err
	}
	s := ConcreteState{}
	for _, b := range c.Originals {
		s.Originals[slot(b.Original)].Binding = b
	}
	return s, nil
}

// Validate checks retained premises and exact records, not oracle authenticity
// or global reachability. It never reads WorldResult.
// Reachable sources are still required by the universal lifting proposition.
func (s ConcreteState) Validate(c Config) error {
	if err := c.Validate(); err != nil {
		return err
	}
	for _, a := range s.Agents {
		initial := !a.Down && !a.Crashed && a.Activation == 0
		down := a.Down && a.Crashed && a.Activation == 0
		recovered := !a.Down && a.Crashed && a.Activation == 1
		if !initial && !down && !recovered {
			return fmt.Errorf("invalid physical availability")
		}
	}
	for i, e := range s.Originals {
		if i >= len(c.Originals) {
			if e != (Transport{}) {
				return fmt.Errorf("undeclared original")
			}
			continue
		}
		if e.Binding != c.Originals[i] {
			return fmt.Errorf("forged concrete binding")
		}
		if (e.Admitted && !e.Committed) || (e.Applied && !e.Admitted) || (e.Consumed && e.Recorded == "") {
			return fmt.Errorf("missing causal predecessor")
		}
		// Validate retained records without consulting the inaccessible oracle.
		var known twoagent.Value
		for _, value := range []twoagent.Value{e.Volatile, e.Recorded, e.Retained, e.ResultChannel} {
			if value == "" {
				continue
			}
			if !e.Applied || (value != twoagent.R1 && value != twoagent.R2) || (known != "" && value != known) {
				return fmt.Errorf("inconsistent retained receipt")
			}
			known = value
		}
		if e.Applied && e.Retained == "" && e.Recorded == "" {
			return fmt.Errorf("lost durable custody")
		}
		if (e.ACKChannel && e.Recorded == "") || (e.RequestChannel && !e.Committed) {
			return fmt.Errorf("channel lacks durable causal source")
		}
		if s.Agents[agentSlot(e.Binding.Caller)].Down && e.Volatile != "" {
			return fmt.Errorf("down caller retains volatile receipt")
		}
	}
	return nil
}

func Project(c Config, s ConcreteState) State {
	q := Initial(c)
	for _, b := range c.Originals {
		i := slot(b.Original)
		e := s.Originals[i]
		if e.Committed {
			q.Stages[i] = W
		}
		if e.Applied {
			q.Stages[i] = P
		}
		if e.Consumed {
			q.Stages[i] = D
		}
	}
	return q
}

// ProjectEdge renames a concrete original to the unchanged twoagent names.
// Physical callee interruptions are stutters in that edge projection.
func ProjectEdge(c Config, s ConcreteState, id string) twoagent.State {
	i := slot(id)
	e := s.Originals[i]
	a := s.Agents[agentSlot(e.Binding.Caller)]
	return twoagent.State{Delegation: twoagent.FixedDelegation(), Grant: twoagent.FixedGrant,
		Committed: e.Committed, ADown: a.Down, Crashed: a.Crashed, Activation: a.Activation,
		AVolatile: e.Volatile, ARecorded: e.Recorded, AConsumed: e.Consumed,
		BAdmitted: e.Admitted, Applied: e.Applied, BResult: e.Retained,
		RequestChannel: e.RequestChannel, ResultChannel: e.ResultChannel, ACKChannel: e.ACKChannel, WorldResult: e.WorldResult}
}

func transport(b Binding, s twoagent.State) Transport {
	return Transport{Binding: b, Committed: s.Committed, Admitted: s.BAdmitted, Applied: s.Applied, Consumed: s.AConsumed,
		Volatile: s.AVolatile, Recorded: s.ARecorded, Retained: s.BResult, RequestChannel: s.RequestChannel,
		ResultChannel: s.ResultChannel, ACKChannel: s.ACKChannel, WorldResult: s.WorldResult}
}
func calleeStep(k twoagent.Kind) bool {
	switch k {
	case twoagent.Admit, twoagent.Deduplicate, twoagent.Apply, twoagent.SendResult, twoagent.Release:
		return true
	}
	return false
}

// ConcreteSuccessors reuses the original retry-enabled edge relation. Only
// actor availability and the named serial apply guard restrict edge work.
// No guard reads WorldResult or the remote caller's receipt to suppress retries.
func ConcreteSuccessors(c Config, s ConcreteState) []MicroStep {
	out := []MicroStep{}
	q := Project(c, s)
	for _, b := range c.Originals {
		i := slot(b.Original)
		for _, step := range twoagent.Successors(twoagent.DurableACK, ProjectEdge(c, s, b.Original)) {
			e := step.Event
			if e.Kind == twoagent.Crash || e.Kind == twoagent.Recover {
				continue
			}
			actor := b.Caller
			if calleeStep(e.Kind) {
				actor = b.Callee
			}
			if s.Agents[agentSlot(actor)].Down {
				continue
			}
			if e.Kind == twoagent.Apply && !ServiceGuard(c, q, b.Callee) {
				continue
			}
			next := s
			next.Originals[i] = transport(b, step.State)
			out = append(out, MicroStep{Event{e.Kind, actor, b.Original, b, e.Result, e.Duplicate}, next})
		}
	}
	for i, a := range s.Agents {
		next := s
		kind := twoagent.Crash
		if a.Down {
			next.Agents[i].Down = false
			next.Agents[i].Activation = 1
			kind = twoagent.Recover
		} else if !a.Crashed {
			next.Agents[i].Down = true
			next.Agents[i].Crashed = true
			for j, e := range next.Originals {
				if e.Binding.Caller == c.Agents[i].ID {
					next.Originals[j].Volatile = ""
				}
			}
		} else {
			continue
		}
		out = append(out, MicroStep{Event: Event{Kind: kind, Actor: c.Agents[i].ID}, State: next})
	}
	return out
}

// Lift realizes one enabled abstract transition from a reachable representative.
// It chooses no new crash and consults actual custody, never the world oracle.
// Reachability/custody is a premise; failure is returned for a malformed source.
func Lift(c Config, s ConcreteState, a Action) ([]MicroStep, error) {
	if err := s.Validate(c); err != nil {
		return nil, err
	}
	enabled := false
	for _, step := range Successors(c, Project(c, s)) {
		if step.Action == a {
			enabled = true
		}
	}
	if !enabled {
		return nil, fmt.Errorf("abstract action is not enabled or has a forged binding")
	}
	path := []MicroStep{}
	take := func(id string, k twoagent.Kind) error {
		for _, step := range ConcreteSuccessors(c, s) {
			if step.Event.Kind == k && !step.Event.Duplicate && (step.Event.Original == id || step.Event.Actor == id) {
				path = append(path, step)
				s = step.State
				return nil
			}
		}
		return fmt.Errorf("no legal %s for %s: missing concrete premise", k, id)
	}
	recover := func(id string) error {
		if s.Agents[agentSlot(id)].Down {
			return take(id, twoagent.Recover)
		}
		return nil
	}
	b := a.Binding
	i := slot(a.Original)
	if err := recover(b.Caller); err != nil {
		return nil, err
	}
	switch a.Kind {
	case "commit":
		if err := take(a.Original, twoagent.Commit); err != nil {
			return nil, err
		}
	case "apply":
		if err := recover(b.Callee); err != nil {
			return nil, err
		}
		if !s.Originals[i].Admitted {
			if !s.Originals[i].RequestChannel {
				if err := take(a.Original, twoagent.SendRequest); err != nil {
					return nil, err
				}
			}
			if err := take(a.Original, twoagent.Admit); err != nil {
				return nil, err
			}
		}
		if err := take(a.Original, twoagent.Apply); err != nil {
			return nil, err
		}
	case "consume":
		if s.Originals[i].Recorded == "" {
			if s.Originals[i].Volatile == "" {
				if s.Originals[i].ResultChannel == "" {
					if err := recover(b.Callee); err != nil {
						return nil, err
					}
					if err := take(a.Original, twoagent.SendResult); err != nil {
						return nil, err
					}
				}
				if err := take(a.Original, twoagent.ReceiveResult); err != nil {
					return nil, err
				}
			}
			if err := take(a.Original, twoagent.Record); err != nil {
				return nil, err
			}
		}
		if err := take(a.Original, twoagent.Consume); err != nil {
			return nil, err
		}
	}
	return path, nil
}

// CheckWitness first checks exact concrete adjacency, then invokes the old
// independently implemented causal ledger on each renamed original history.
// Adjacency shares Successors; neither check is a runtime refinement proof.
func CheckWitness(w Witness) error {
	initial, err := NewConcrete(w.Config)
	if err != nil {
		return err
	}
	if w.Initial != initial {
		return fmt.Errorf("witness must start at the declared initial state")
	}
	s := w.Initial
	projection := []Action{}
	ledgers := map[string][]twoagent.Event{}
	for i, step := range w.Steps {
		legal := false
		for _, candidate := range ConcreteSuccessors(w.Config, s) {
			if candidate == step {
				legal = true
				break
			}
		}
		if !legal {
			return fmt.Errorf("step %d: illegal concrete event, binding or target", i)
		}
		if Project(w.Config, s) != Project(w.Config, step.State) {
			projection = append(projection, Action{string(step.Event.Kind), step.Event.Original, step.Event.Binding})
		}
		for _, b := range w.Config.Originals {
			e := step.Event
			physical := (e.Kind == twoagent.Crash || e.Kind == twoagent.Recover) && e.Actor == b.Caller
			if e.Original == b.Original || physical {
				ledgers[b.Original] = append(ledgers[b.Original], twoagent.Event{Kind: e.Kind, Delegation: twoagent.FixedDelegation(), Grant: twoagent.FixedGrant, Result: e.Result, Duplicate: e.Duplicate})
			}
		}
		s = step.State
	}
	for id, events := range ledgers {
		if err := twoagent.CheckTrace(twoagent.DurableACK, events); err != nil {
			return fmt.Errorf("%s causal ledger: %w", id, err)
		}
	}
	if len(w.Projection) != len(projection) {
		return fmt.Errorf("wrong progress projection length")
	}
	for i, a := range projection {
		if w.Projection[i] != a {
			return fmt.Errorf("wrong progress observation at %d", i)
		}
	}
	if w.Complete != Complete(w.Config, Project(w.Config, s)) {
		return fmt.Errorf("wrong completion classification")
	}
	found := false
	for _, n := range Explore(w.Config).Nodes {
		if n.State == Project(w.Config, s) {
			found = true
			if n.CanComplete != w.CanComplete {
				return fmt.Errorf("wrong CanComplete classification")
			}
		}
	}
	if !found {
		return fmt.Errorf("projection absent from macro graph")
	}
	return nil
}
