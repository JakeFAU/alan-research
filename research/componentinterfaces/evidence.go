package componentinterfaces

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	pc "github.com/JakeFAU/alan-research/research/capacitycomposition"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type substitution struct {
	Left             string
	Right            string
	Context          string
	ProductLeft      string
	ProductRight     string
	CanCompleteLeft  bool
	CanCompleteRight bool
}
type counts struct {
	Plants                  int
	Monitors                int
	Controllers             int
	Products                int
	ProductStates           int
	ProductArcs             int
	ProductCanComplete      int
	InitialCanComplete      int
	NonblockingProducts     int
	Comparisons             int
	EqualComparisons        int
	SubsetPairs             int
	Groupings               int
	Substitutions           int
	SuccessfulSubstitutions int
	BlockedSubstitutions    int
	Histories               int
	ConcreteSteps           int
}
type source struct {
	Path   string
	SHA256 string
}
type evidence struct {
	Schema        string
	Alphabet      []symbol
	Monitors      []automaton
	Plants        []automaton
	Controllers   []automaton
	Products      []graph
	Comparisons   []comparison
	Groupings     []grouping
	Substitutions []substitution
	Histories     []history
	Counts        counts
	Sources       []source
}

func buildEvidence() (evidence, error) {
	c, err := pc.NewConfig("shared-1", "early", 0)
	if err != nil {
		return evidence{}, err
	}
	e := evidence{Schema: "alan/phase3c3/v1", Alphabet: alphabet(c), Monitors: monitors(c), Controllers: controllers(c), Plants: []automaton{}, Products: []graph{}, Comparisons: []comparison{}, Groupings: []grouping{}, Substitutions: []substitution{}, Sources: []source{}}
	for i, a := range e.Monitors {
		for _, b := range e.Monitors[i+1:] {
			e.Comparisons = append(e.Comparisons, compareLanguages(a, b))
		}
	}
	recordSubstitutions := func(context string, gs []graph) error {
		for _, cmp := range e.Comparisons {
			if !cmp.Equal {
				continue
			}
			var left, right graph
			for i, m := range e.Monitors {
				if m.ID == cmp.Left {
					left = gs[i]
				}
				if m.ID == cmp.Right {
					right = gs[i]
				}
			}
			l, r := left.Nodes[0].CanComplete, right.Nodes[0].CanComplete
			if l != r {
				return fmt.Errorf("initial language substitution failed")
			}
			e.Substitutions = append(e.Substitutions, substitution{cmp.Left, cmp.Right, context, left.ID, right.ID, l, r})
		}
		return nil
	}
	for _, layout := range []string{"shared-1", "shared-2", "separate-1"} {
		for _, policy := range []string{"early", "ready-now", "stable-ready"} {
			for mask := uint8(0); mask < 8; mask++ {
				c, err := pc.NewConfig(layout, policy, mask)
				if err != nil {
					return evidence{}, err
				}
				k := plant(c)
				e.Plants = append(e.Plants, k)
				gs := []graph{}
				for _, m := range monitors(c) {
					g, err := product(k, m)
					if err != nil {
						return evidence{}, err
					}
					e.Products = append(e.Products, g)
					gs = append(gs, g)
				}
				if err := recordSubstitutions(k.ID, gs); err != nil {
					return evidence{}, err
				}
			}
		}
	}
	k := plant(c)
	for _, d := range e.Controllers {
		gs := []graph{}
		for _, m := range e.Monitors {
			g, err := product(k, m, d)
			if err != nil {
				return evidence{}, err
			}
			group, err := checkGrouping(k, m, d, g)
			if err != nil {
				return evidence{}, err
			}
			if !group.Equal {
				return evidence{}, fmt.Errorf("grouping relation differs")
			}
			e.Products = append(e.Products, g)
			e.Groupings = append(e.Groupings, group)
			gs = append(gs, g)
		}
		if err := recordSubstitutions(k.ID+" || "+d.ID, gs); err != nil {
			return evidence{}, err
		}
	}
	e.Histories, err = histories()
	if err != nil {
		return evidence{}, err
	}
	e.Counts = counts{Plants: len(e.Plants), Monitors: len(e.Monitors), Controllers: len(e.Controllers), Products: len(e.Products), Comparisons: len(e.Comparisons), Groupings: len(e.Groupings), Substitutions: len(e.Substitutions), Histories: len(e.Histories)}
	for _, g := range e.Products {
		e.Counts.ProductStates += g.States
		e.Counts.ProductArcs += g.Arcs
		e.Counts.ProductCanComplete += g.CanComplete
		if g.Nodes[0].CanComplete {
			e.Counts.InitialCanComplete++
		}
		if g.Nonblocking {
			e.Counts.NonblockingProducts++
		}
	}
	for _, c := range e.Comparisons {
		e.Counts.SubsetPairs += len(c.Pairs)
		if c.Equal {
			e.Counts.EqualComparisons++
		}
	}
	for _, s := range e.Substitutions {
		if s.CanCompleteLeft {
			e.Counts.SuccessfulSubstitutions++
		} else {
			e.Counts.BlockedSubstitutions++
		}
	}
	for _, h := range e.Histories {
		e.Counts.ConcreteSteps += len(h.Steps)
	}
	return e, nil
}

// sourceManifest binds all non-test compiled local package sources. Package
// directories are explicit: introducing another local import requires updating
// this list and the schema, and review against go list -deps.
func sourceManifest() ([]source, error) {
	paths := []string{"go.mod", "go.sum", "docs/foundational/phase3c3/00-contract.md", "docs/foundational/phase3c3/01-argument.md", "docs/foundational/phase3c3/02-evidence-schema.md"}
	for _, dir := range []string{"research/componentinterfaces", "research/componentinterfaces/cmd/export", "research/capacitycomposition", "research/contractcomposition", "research/twoagent"} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
				paths = append(paths, filepath.ToSlash(filepath.Join(dir, entry.Name())))
			}
		}
	}
	slices.Sort(paths)
	out := []source{}
	for _, path := range paths {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		hash := sha256.Sum256(b)
		out = append(out, source{path, hex.EncodeToString(hash[:])})
	}
	return out, nil
}

// Export writes complete deterministic evidence from the repository root. It
// checks all declared grouping and initial substitution obligations first.
func Export(w io.Writer) error {
	e, err := buildEvidence()
	if err != nil {
		return err
	}
	e.Sources, err = sourceManifest()
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(e)
}
