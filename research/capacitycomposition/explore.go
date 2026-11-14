package capacitycomposition

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
type Graph struct {
	Config         Config
	Nodes          []Node
	States         int
	Arcs           int
	CanComplete    int
	CannotComplete int
}

// Explore uses complete BFS and reverse reachability, retaining non-tree arcs.
func Explore(c Config) Graph {
	initial := Initial(c)
	g := Graph{Config: c, Nodes: []Node{{State: initial, Arcs: []Arc{}, Parent: -1}}}
	ids := map[State]int{initial: 0}
	for i := 0; i < len(g.Nodes); i++ {
		for _, e := range Successors(c, g.Nodes[i].State) {
			to, ok := ids[e.State]
			if !ok {
				to = len(g.Nodes)
				ids[e.State] = to
				a := e.Action
				g.Nodes = append(g.Nodes, Node{State: e.State, Arcs: []Arc{}, Parent: i, Via: &a})
			}
			g.Nodes[i].Arcs = append(g.Nodes[i].Arcs, Arc{e.Action, to})
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
