package capacitycomposition

import (
	"fmt"
	cc "github.com/JakeFAU/alan-research/research/contractcomposition"
	"github.com/JakeFAU/alan-research/research/twoagent"
	"reflect"
	"slices"
)

type Witness struct {
	Name        string
	Config      Config
	Initial     ConcreteState
	Steps       []MicroStep
	Projection  []Action
	Complete    bool
	CanComplete bool
}
type Instruction struct {
	ID        string
	Kind      twoagent.Kind
	Result    twoagent.Value
	Duplicate bool
}

func appendInstruction(script []Instruction, id string, kind twoagent.Kind) []Instruction {
	return append(script, Instruction{ID: id, Kind: kind})
}
func admittedScript() []Instruction {
	out := []Instruction{}
	for _, id := range []string{"AB", "BC"} {
		for _, kind := range []twoagent.Kind{twoagent.Commit, twoagent.SendRequest, twoagent.Admit} {
			out = appendInstruction(out, id, kind)
		}
	}
	return out
}
func completeScript(prefix []Instruction) []Instruction {
	out := slices.Clone(prefix)
	for _, id := range []string{"BC", "AB"} {
		if id == "AB" {
			out = appendInstruction(out, id, "reserve-capacity")
		}
		for _, kind := range []twoagent.Kind{twoagent.Apply, twoagent.SendResult, twoagent.ReceiveResult, twoagent.Record, twoagent.Consume} {
			out = appendInstruction(out, id, kind)
		}
	}
	return out
}
func buildWitness(name string, c Config, script []Instruction) (Witness, error) {
	s, err := NewConcrete(c)
	if err != nil {
		return Witness{}, err
	}
	w := Witness{Name: name, Config: c, Initial: s, Steps: []MicroStep{}, Projection: []Action{}}
	for _, instruction := range script {
		found := false
		for _, e := range ConcreteSuccessors(c, s) {
			if e.Event.Kind != instruction.Kind || e.Event.Duplicate != instruction.Duplicate || (e.Event.Original != instruction.ID && e.Event.Actor != instruction.ID) || (instruction.Result != "" && e.Event.Result != instruction.Result) {
				continue
			}
			if Project(c, s) != Project(c, e.State) {
				w.Projection = append(w.Projection, Action{Kind: string(e.Event.Kind), Original: e.Event.Original, Binding: e.Event.Binding})
			}
			w.Steps = append(w.Steps, e)
			s = e.State
			found = true
			break
		}
		if !found {
			return Witness{}, fmt.Errorf("%s: no legal instruction %+v", name, instruction)
		}
	}
	w.Complete = Complete(c, Project(c, s))
	for _, n := range Explore(c).Nodes {
		if n.State == Project(c, s) {
			w.CanComplete = n.CanComplete
		}
	}
	if err := CheckWitness(w); err != nil {
		return Witness{}, err
	}
	return w, nil
}

// BuildWitnesses constructs every history from empty initial custody using
// actual transitions. Full post-states are captured, never repaired or merged.
func BuildWitnesses() ([]Witness, error) {
	out := []Witness{}
	add := func(name, policy string, mask uint8, script []Instruction) error {
		c, err := NewConfig("shared-1", policy, mask)
		if err != nil {
			return err
		}
		w, err := buildWitness(name, c, script)
		if err != nil {
			return err
		}
		out = append(out, w)
		return nil
	}
	earlyAB := appendInstruction(admittedScript(), "AB", "reserve-capacity")
	earlyBC := appendInstruction(admittedScript(), "BC", "reserve-capacity")
	readyAB := []Instruction{}
	for _, kind := range []twoagent.Kind{twoagent.Commit, twoagent.SendRequest, twoagent.Admit, "reserve-capacity"} {
		readyAB = appendInstruction(readyAB, "AB", kind)
	}
	for _, kind := range []twoagent.Kind{twoagent.Commit, twoagent.SendRequest, twoagent.Admit} {
		readyAB = appendInstruction(readyAB, "BC", kind)
	}
	for _, item := range []struct {
		name, policy string
		script       []Instruction
	}{
		{"early-AB-holder", "early", earlyAB}, {"early-BC-holder", "early", earlyBC}, {"early-BC-completed", "early", completeScript(earlyBC)},
		{"ready-now-AB-holder", "ready-now", readyAB}, {"ready-now-BC-holder", "ready-now", earlyBC}, {"ready-now-BC-completed", "ready-now", completeScript(earlyBC)},
		{"old-AB-before-BC-commit", "early", appendInstruction(slices.Clone(readyAB[:4]), "AB", twoagent.Apply)},
		{"early-AB-holder-crashed", "early", appendInstruction(slices.Clone(earlyAB), "B", twoagent.Crash)},
	} {
		if err := add(item.name, item.policy, 0, item.script); err != nil {
			return nil, err
		}
	}
	for _, v := range []twoagent.Value{twoagent.R1, twoagent.R2} {
		script := slices.Clone(earlyBC)
		script = appendInstruction(script, "C", twoagent.Crash)
		script = appendInstruction(script, "C", twoagent.Recover)
		for _, id := range []string{"BC", "AB"} {
			if id == "AB" {
				script = appendInstruction(script, id, "reserve-capacity")
			}
			script = appendInstruction(script, id, twoagent.SendRequest)
			script = append(script, Instruction{ID: id, Kind: twoagent.Deduplicate, Duplicate: true}, Instruction{ID: id, Kind: twoagent.Deduplicate})
			script = append(script, Instruction{ID: id, Kind: twoagent.Apply, Result: v})
			script = appendInstruction(script, id, twoagent.SendResult)
			script = append(script, Instruction{ID: id, Kind: twoagent.ReceiveResult, Duplicate: true}, Instruction{ID: id, Kind: twoagent.ReceiveResult})
			caller := "A"
			if id == "BC" {
				caller = "B"
			}
			script = appendInstruction(script, caller, twoagent.Crash)
			script = appendInstruction(script, caller, twoagent.Recover)
			script = appendInstruction(script, id, twoagent.SendResult)
			script = appendInstruction(script, id, twoagent.ReceiveResult)
			script = appendInstruction(script, id, twoagent.Record)
			script = appendInstruction(script, id, twoagent.SendACK)
			script = append(script, Instruction{ID: id, Kind: twoagent.Release, Duplicate: true}, Instruction{ID: id, Kind: twoagent.Release})
			script = appendInstruction(script, id, twoagent.Consume)
		}
		if err := add("stable-recovery-custody-"+string(v), "stable-ready", 0, script); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// CheckWitness checks exact adjacency and delegates the erased capacity history
// to the unchanged custody ledger too. Both remain implementation-side checks;
// the independent checker must replay the exported history declaratively.
func CheckWitness(w Witness) error {
	initial, err := NewConcrete(w.Config)
	if err != nil {
		return err
	}
	if w.Initial != initial {
		return fmt.Errorf("noninitial witness custody")
	}
	s := initial
	projection := []Action{}
	old := cc.Witness{Name: w.Name, Config: w.Config.Base, Initial: initial.Base, Steps: []cc.MicroStep{}, Projection: []cc.Action{}}
	for i, e := range w.Steps {
		found := false
		for _, candidate := range ConcreteSuccessors(w.Config, s) {
			if candidate == e {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("illegal witness step %d", i)
		}
		if err := e.State.Validate(w.Config); err != nil {
			return err
		}
		if Project(w.Config, s) != Project(w.Config, e.State) {
			projection = append(projection, Action{Kind: string(e.Event.Kind), Original: e.Event.Original, Binding: e.Event.Binding})
		}
		if e.Event.Kind != "reserve-capacity" {
			old.Steps = append(old.Steps, cc.MicroStep{Event: e.Event, State: e.State.Base})
		}
		if cc.Project(w.Config.Base, s.Base) != cc.Project(w.Config.Base, e.State.Base) {
			old.Projection = append(old.Projection, cc.Action{Kind: string(e.Event.Kind), Original: e.Event.Original, Binding: e.Event.Binding})
		}
		s = e.State
	}
	if !reflect.DeepEqual(projection, w.Projection) || w.Complete != Complete(w.Config, Project(w.Config, s)) {
		return fmt.Errorf("wrong witness projection/completion")
	}
	found := false
	for _, n := range Explore(w.Config).Nodes {
		if n.State == Project(w.Config, s) {
			found = true
			if n.CanComplete != w.CanComplete {
				return fmt.Errorf("wrong witness CanComplete")
			}
		}
	}
	if !found {
		return fmt.Errorf("witness projection unreachable")
	}
	old.Complete = w.Complete
	for _, n := range cc.Explore(w.Config.Base).Nodes {
		if n.State == cc.Project(w.Config.Base, s.Base) {
			old.CanComplete = n.CanComplete
		}
	}
	return cc.CheckWitness(old)
}

type SourceHash struct {
	Path   string
	SHA256 string
}
type Evidence struct {
	Schema    string
	Sources   []SourceHash
	Graphs    []Graph
	Groupings []Grouping
	Witnesses []Witness
	Probes    ProbeSummary
}

// BuildFiniteEvidence builds graphs, regroupings and genuine witnesses; probes
// are attached separately by the exporter, avoiding repeated expensive checks.
func BuildFiniteEvidence() (Evidence, error) {
	e := Evidence{Schema: "alan/phase3c2/v1", Sources: []SourceHash{}, Graphs: []Graph{}, Groupings: []Grouping{}}
	for _, layout := range []string{"shared-1", "shared-2", "separate-1"} {
		for _, policy := range []string{"early", "ready-now", "stable-ready"} {
			for mask := uint8(0); mask < 8; mask++ {
				c, err := NewConfig(layout, policy, mask)
				if err != nil {
					return Evidence{}, err
				}
				g, err := regroup(c)
				if err != nil {
					return Evidence{}, err
				}
				if !g.Equal {
					return Evidence{}, fmt.Errorf("regrouping changed relation")
				}
				e.Graphs = append(e.Graphs, Explore(c))
				e.Groupings = append(e.Groupings, g)
			}
		}
	}
	var err error
	e.Witnesses, err = BuildWitnesses()
	if err != nil {
		return Evidence{}, err
	}
	return e, nil
}
