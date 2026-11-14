package residualinterfaces

import (
	"encoding/json"
	"fmt"
	"io"
)

type substitution struct {
	Left             string
	Right            string
	Context          string
	ProductLeft      string
	ProductRight     string
	NonblockingLeft  bool
	NonblockingRight bool
}
type counts struct {
	Plants                   int
	Monitors                 int
	Controllers              int
	Products                 int
	ProductStates            int
	ProductArcs              int
	ProductCanComplete       int
	InitialCanComplete       int
	NonblockingProducts      int
	Comparisons              int
	EqualLanguages           int
	EqualProfiles            int
	ResidualComparisons      int
	LanguageSubsetPairs      int
	ProfileSubsetPairs       int
	Substitutions            int
	NonblockingSubstitutions int
	BlockingSubstitutions    int
}
type evidence struct {
	Schema        string
	Input         inputBinding
	Alphabet      []symbol
	Monitors      []automaton
	Plants        []automaton
	Controllers   []automaton
	Comparisons   []profileComparison
	Products      []graph
	Substitutions []substitution
	CoreOutcomes  []coreOutcome
	Negative      negativeWitness
	Loop          loopWitness
	Strictness    strictnessWitness
	Counts        counts
	Sources       []sourceBinding
}

func build(root string) (evidence, error) {
	p, binding, err := readInput(root)
	if err != nil {
		return evidence{}, fmt.Errorf("input: %w", err)
	}
	e := evidence{Schema: "alan/phase3c4/v1", Input: binding, Alphabet: p.Alphabet, Plants: p.Plants, Controllers: p.Controllers,
		Monitors: append(p.Monitors, extraMonitors(p.Monitors[0].Signature)...), Comparisons: []profileComparison{}, Products: []graph{}, Substitutions: []substitution{}, CoreOutcomes: []coreOutcome{}}
	for i, a := range e.Monitors {
		for _, b := range e.Monitors[i+1:] {
			e.Comparisons = append(e.Comparisons, compareProfiles(a, b))
		}
	}
	core := map[string]graph{}
	contextProducts := map[string]map[string]graph{}
	contexts := []string{}
	var corePlant automaton
	for _, k := range p.Plants {
		contexts = append(contexts, k.ID)
		contextProducts[k.ID] = map[string]graph{}
		for _, m := range e.Monitors {
			m.Signature = k.Signature
			g, err := product(k, m)
			if err != nil {
				return e, err
			}
			e.Products = append(e.Products, g)
			contextProducts[k.ID][m.ID] = g
			if k.ID == "early/shared-1/mask0" {
				corePlant = k
				core[m.ID] = g
				e.CoreOutcomes = append(e.CoreOutcomes, coreOutcome{Monitor: m.ID, Product: g.ID, InitialCanComplete: g.Nodes[0].CanComplete, Nonblocking: g.Nonblocking})
			}
		}
	}
	if corePlant.ID == "" {
		return e, fmt.Errorf("missing core plant")
	}
	for _, d := range p.Controllers {
		context := corePlant.ID + " || " + d.ID
		contexts = append(contexts, context)
		contextProducts[context] = map[string]graph{}
		for _, m := range e.Monitors {
			m.Signature = corePlant.Signature
			g, err := product(corePlant, m, d)
			if err != nil {
				return e, err
			}
			e.Products = append(e.Products, g)
			contextProducts[context][m.ID] = g
		}
	}
	for _, c := range e.Comparisons {
		if !c.Equal {
			continue
		}
		for _, context := range contexts {
			l, r := contextProducts[context][c.Left], contextProducts[context][c.Right]
			e.Substitutions = append(e.Substitutions, substitution{Left: c.Left, Right: c.Right, Context: context, ProductLeft: l.ID, ProductRight: r.ID, NonblockingLeft: l.Nonblocking, NonblockingRight: r.Nonblocking})
		}
	}
	e.Negative, e.Loop, err = makeWitnesses(core, e.Comparisons)
	if err != nil {
		return e, err
	}
	x, y := genericStrictness()
	e.Strictness = strictnessWitness{Left: x, Right: y, Comparison: compareProfiles(x, y)}
	e.Counts = countEvidence(e)
	e.Sources, err = sourceManifest(root, p.Sources)
	return e, err
}
func countEvidence(e evidence) counts {
	c := counts{Plants: len(e.Plants), Monitors: len(e.Monitors), Controllers: len(e.Controllers), Products: len(e.Products), Comparisons: len(e.Comparisons), Substitutions: len(e.Substitutions)}
	for _, g := range e.Products {
		c.ProductStates += g.States
		c.ProductArcs += g.Arcs
		c.ProductCanComplete += g.CanComplete
		if g.Nodes[0].CanComplete {
			c.InitialCanComplete++
		}
		if g.Nonblocking {
			c.NonblockingProducts++
		}
	}
	for _, p := range e.Comparisons {
		if p.InitialLanguage.Equal {
			c.EqualLanguages++
		}
		if p.Equal {
			c.EqualProfiles++
		}
		c.ResidualComparisons += len(p.Residuals)
		c.ProfileSubsetPairs += len(p.Pairs)
		c.LanguageSubsetPairs += len(p.InitialLanguage.Pairs)
		for _, r := range p.Residuals {
			c.LanguageSubsetPairs += len(r.Language.Pairs)
		}
	}
	for _, s := range e.Substitutions {
		if s.NonblockingLeft && s.NonblockingRight {
			c.NonblockingSubstitutions++
		}
		if !s.NonblockingLeft && !s.NonblockingRight {
			c.BlockingSubstitutions++
		}
	}
	return c
}

// Export writes deterministic complete evidence. Call from the repository root.
func Export(w io.Writer) error {
	e, err := build(".")
	if err != nil {
		return err
	}
	if err := json.NewEncoder(w).Encode(e); err != nil {
		return fmt.Errorf("encode evidence: %w", err)
	}
	return nil
}
