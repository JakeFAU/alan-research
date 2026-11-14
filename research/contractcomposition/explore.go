package contractcomposition

import (
	"fmt"
	"reflect"
	"slices"
)

type Arc struct {
	Action Action
	To     int
}
type Node struct {
	State       State
	Arcs        []Arc
	Complete    bool
	CanComplete bool
	Parent      int
	Via         *Action
}
type Dependency struct {
	Original string
	Requires string
}
type Certificate struct {
	Dependencies []Dependency
	Acyclic      bool
	Order        []string
}
type Graph struct {
	Config         Config
	Nodes          []Node
	States         int
	Arcs           int
	CanComplete    int
	CannotComplete int
	Certificate    Certificate
}

// Explore keeps every arc, including multiple incoming arcs. Reverse search
// decides existence of a completion suffix, never inevitable completion.
func Explore(c Config) Graph {
	if err := c.Validate(); err != nil {
		panic(err)
	}
	initial := Initial(c)
	g := Graph{Config: c, Nodes: []Node{{State: initial, Parent: -1, Arcs: []Arc{}}}, Certificate: Dependencies(c)}
	ids := map[State]int{initial: 0}
	for i := 0; i < len(g.Nodes); i++ {
		for _, step := range Successors(c, g.Nodes[i].State) {
			to, ok := ids[step.State]
			if !ok {
				to = len(g.Nodes)
				ids[step.State] = to
				a := step.Action
				g.Nodes = append(g.Nodes, Node{State: step.State, Parent: i, Via: &a, Arcs: []Arc{}})
			}
			g.Nodes[i].Arcs = append(g.Nodes[i].Arcs, Arc{step.Action, to})
			g.Arcs++
		}
	}
	incoming := make([][]int, len(g.Nodes))
	queue := []int{}
	for i, n := range g.Nodes {
		for _, a := range n.Arcs {
			incoming[a.To] = append(incoming[a.To], i)
		}
		if Complete(c, n.State) {
			g.Nodes[i].Complete = true
			g.Nodes[i].CanComplete = true
			queue = append(queue, i)
		}
	}
	for i := 0; i < len(queue); i++ {
		for _, from := range incoming[queue[i]] {
			if !g.Nodes[from].CanComplete {
				g.Nodes[from].CanComplete = true
				queue = append(queue, from)
			}
		}
	}
	g.States = len(g.Nodes)
	g.CanComplete = len(queue)
	g.CannotComplete = g.States - g.CanComplete
	return g
}

// Dependencies conservatively orders a serial callee's outgoing original first.
// A cyclic certificate makes no necessity claim for a particular state.
func Dependencies(c Config) Certificate {
	cert := Certificate{Dependencies: []Dependency{}, Order: []string{}}
	for _, e := range c.Originals {
		if c.Agents[agentSlot(e.Callee)].Servicing {
			continue
		}
		for _, f := range c.Originals {
			if f.Caller == e.Callee {
				cert.Dependencies = append(cert.Dependencies, Dependency{e.Original, f.Original})
			}
		}
	}
	done := map[string]bool{}
	for len(done) < len(c.Originals) {
		advanced := false
		for _, e := range c.Originals {
			if done[e.Original] {
				continue
			}
			ready := true
			for _, d := range cert.Dependencies {
				if d.Original == e.Original && !done[d.Requires] {
					ready = false
				}
			}
			if ready {
				done[e.Original] = true
				cert.Order = append(cert.Order, e.Original)
				advanced = true
			}
		}
		if !advanced {
			return cert
		}
	}
	cert.Acyclic = true
	return cert
}

// BuildWitnesses retains concrete histories, including actual transport retry.
// These are deterministic constructive witnesses, with no global minimality claim.
func BuildWitnesses() ([]Witness, error) {
	out := []Witness{}
	build := func(name, topology string, mask uint8, script [][2]string) error {
		c, err := NewConfig(topology, mask)
		if err != nil {
			return err
		}
		s, err := NewConcrete(c)
		if err != nil {
			return err
		}
		w := Witness{Name: name, Config: c, Initial: s, Steps: []MicroStep{}, Projection: []Action{}}
		for _, instruction := range script {
			found := false
			for _, step := range ConcreteSuccessors(c, s) {
				e := step.Event
				if string(e.Kind) != instruction[1] || e.Duplicate || (e.Original != instruction[0] && e.Actor != instruction[0]) {
					continue
				}
				if Project(c, s) != Project(c, step.State) {
					w.Projection = append(w.Projection, Action{string(e.Kind), e.Original, e.Binding})
				}
				w.Steps = append(w.Steps, step)
				s = step.State
				found = true
				break
			}
			if !found {
				return fmt.Errorf("%s: no concrete step %v", name, instruction)
			}
		}
		q := Project(c, s)
		w.Complete = Complete(c, q)
		for _, n := range Explore(c).Nodes {
			if n.State == q {
				w.CanComplete = n.CanComplete
			}
		}
		if err := CheckWitness(w); err != nil {
			return err
		}
		out = append(out, w)
		return nil
	}
	r0 := [][2]string{{"AB", "commit"}, {"BC", "commit"}, {"BC", "send-request"}, {"BC", "admit"}, {"BC", "apply"}, {"BC", "send-result"}, {"BC", "receive-result"}, {"BC", "record"}, {"BC", "consume"}}
	if err := build("R0-unpaid", "chain", 7, r0); err != nil {
		return nil, err
	}
	paid := append(slices.Clone(r0), [][2]string{{"AB", "send-request"}, {"AB", "admit"}, {"AB", "apply"}, {"AB", "send-result"}, {"AB", "receive-result"}, {"AB", "record"}, {"AB", "consume"}}...)
	if err := build("R0-paid", "chain", 7, paid); err != nil {
		return nil, err
	}
	r1 := [][2]string{{"CA", "commit"}, {"CA", "send-request"}, {"CA", "admit"}}
	if err := build("R1-uncommitted", "cycle", 6, r1); err != nil {
		return nil, err
	}
	if err := build("R1-waiting", "cycle", 6, append(slices.Clone(r1), [2]string{"AB", "commit"})); err != nil {
		return nil, err
	}
	if err := build("serial-cycle-control", "cycle", 0, [][2]string{{"AB", "commit"}, {"BC", "commit"}, {"CA", "commit"}}); err != nil {
		return nil, err
	}
	retry := [][2]string{{"AB", "commit"}, {"AB", "send-request"}, {"AB", "admit"}, {"AB", "send-request"}, {"AB", "deduplicate"}}
	if err := build("retry-divergence", "chain", 7, retry); err != nil {
		return nil, err
	}
	loss := [][2]string{{"AB", "commit"}, {"AB", "send-request"}, {"AB", "admit"}, {"AB", "apply"}, {"AB", "send-result"}, {"AB", "receive-result"}, {"A", "crash"}, {"A", "recover"}, {"AB", "send-result"}, {"AB", "receive-result"}, {"AB", "record"}, {"AB", "send-ack"}, {"AB", "release"}, {"AB", "consume"}}
	if err := build("pre-record-loss-and-release", "chain", 7, loss); err != nil {
		return nil, err
	}
	return out, nil
}

type Grouping struct {
	ConfigID    string
	A           Component
	B           Component
	C           Component
	AB          Component
	BC          Component
	Left        Component
	Right       Component
	Flat        Config
	LeftClosed  Config
	RightClosed Config
	Equal       bool
}
type Collision struct {
	Rule             string
	Left             string
	Right            string
	HiddenOriginal   string
	ExposedOriginals []string
	Observation      string
}
type Claim struct {
	ID        string
	Status    string
	Statement string
	Evidence  []string
	Limits    []string
}
type RetryCertificate struct {
	Witness   string
	FromState int
	ToState   int
}
type Evidence struct {
	Schema     string
	Graphs     []Graph
	Groupings  []Grouping
	Witnesses  []Witness
	Collisions []Collision
	Retry      RetryCertificate
	Claims     []Claim
}

// BuildEvidence is deterministic: identities, policy masks, BFS successors and
// witness scripts all have a fixed order; every count comes from exploration.
func BuildEvidence() (Evidence, error) {
	e := Evidence{Schema: "alan.phase3c1.progress.v1", Graphs: []Graph{}, Groupings: []Grouping{}}
	for _, topology := range []string{"chain", "cycle"} {
		for mask := uint8(0); mask < 8; mask++ {
			c, err := NewConfig(topology, mask)
			if err != nil {
				return Evidence{}, err
			}
			a, err := Leaf(c, "A")
			if err != nil {
				return Evidence{}, err
			}
			b, err := Leaf(c, "B")
			if err != nil {
				return Evidence{}, err
			}
			d, err := Leaf(c, "C")
			if err != nil {
				return Evidence{}, err
			}
			ab, err := Join(a, b)
			if err != nil {
				return Evidence{}, err
			}
			bc, err := Join(b, d)
			if err != nil {
				return Evidence{}, err
			}
			left, err := Join(ab, d)
			if err != nil {
				return Evidence{}, err
			}
			right, err := Join(a, bc)
			if err != nil {
				return Evidence{}, err
			}
			lc, err := Close(left)
			if err != nil {
				return Evidence{}, err
			}
			rc, err := Close(right)
			if err != nil {
				return Evidence{}, err
			}
			graph := Explore(c)
			equal := reflect.DeepEqual(left, right) && reflect.DeepEqual(graph, Explore(lc)) && reflect.DeepEqual(graph, Explore(rc))
			e.Graphs = append(e.Graphs, graph)
			e.Groupings = append(e.Groupings, Grouping{fmt.Sprintf("%s/%d", topology, mask), a, b, d, ab, bc, left, right, c, lc, rc, equal})
		}
	}
	ws, err := BuildWitnesses()
	if err != nil {
		return Evidence{}, err
	}
	e.Witnesses = ws
	e.Collisions = []Collision{
		{"R0", "R0-unpaid", "R0-paid", "AB", []string{"BC"}, "all-original Complete and named internal apply/consume events"},
		{"R1", "R1-uncommitted", "R1-waiting", "AB", []string{"BC", "CA"}, "immediate concrete CA application readiness with A up and CA admitted"},
	}
	e.Retry = RetryCertificate{"retry-divergence", 3, 5}
	e.Claims = []Claim{
		{"finite-path-equivalence", "source-argument", "Every reachable concrete representative supports each enabled macro transition by finite legal stutters then that transition; completion suffix existence agrees.", []string{"docs/foundational/phase3c1/01-argument.md", "research/contractcomposition/concrete.go"}, []string{"Fixed chain/cycle and declared atomic custody/recovery premises only.", "Source argument, not machine-checked proof; probes are factored, not full concrete product enumeration."}},
		{"macro-graphs", "bounded-exhaustive", "Complete reachable progress graphs and reverse CanComplete for both topologies and all eight masks.", []string{"Graphs"}, []string{"Value-abstracted progress; no fairness, divergence equivalence or concrete immediate readiness claim."}},
		{"compatible-grouping", "bounded-check-and-source-argument", "Compatible grouping retains original bindings, debt and physical-node policy incidence.", []string{"Groupings", "docs/foundational/phase3c1/01-argument.md"}, []string{"Algebraic equality alone does not prove concrete abstraction.", "No arbitrary contexts or shared-authority fan-out."}},
		{"candidate-failures", "concrete-witnesses", "R0 loses completion and visible progress; R1 loses the specified readiness distinction.", []string{"Collisions", "Witnesses"}, []string{"No global shortest-history or smallest-topology claim.", "R1 alone does not refute a relation omitting readiness."}},
	}
	return e, nil
}
