// Package componentinterfaces studies completion-preserving finite continuation
// monitors over the unchanged Phase 3C.2 plant. It is not runtime refinement.
package componentinterfaces

import (
	"fmt"
	pc "github.com/JakeFAU/alan-research/research/capacitycomposition"
	cc "github.com/JakeFAU/alan-research/research/contractcomposition"
	"reflect"
	"slices"
)

type symbol struct {
	ID     string
	Action cc.Action
	Actor  string
	Root   string
}
type state struct {
	ID     string
	Tuple  []string
	Marked bool
}
type edge struct {
	From  string
	Label string
	To    string
}
type automaton struct {
	ID        string
	Role      string
	Signature pc.Config
	Initial   string
	States    []state
	Edges     []edge
}

func alphabet(c pc.Config) []symbol {
	out := []symbol{}
	for i, b := range c.Base.Originals {
		for _, kind := range []string{"commit", "reserve-capacity", "apply", "consume"} {
			actor := b.Caller
			if kind == "reserve-capacity" || kind == "apply" {
				actor = b.Callee
			}
			out = append(out, symbol{kind + ":" + b.Original, cc.Action{Kind: kind, Original: b.Original, Binding: b}, actor, c.Memberships[i].Root})
		}
	}
	return out
}
func validate(a automaton) error {
	if err := a.Signature.Validate(); err != nil {
		return err
	}
	if a.ID == "" {
		return fmt.Errorf("empty automaton ID")
	}
	if a.Role != "K" && a.Role != "M" && a.Role != "D" && a.Role != "product" {
		return fmt.Errorf("unknown role")
	}
	states := map[string]bool{}
	for _, s := range a.States {
		if s.ID == "" || states[s.ID] {
			return fmt.Errorf("empty or duplicate state")
		}
		states[s.ID] = true
	}
	if !states[a.Initial] {
		return fmt.Errorf("undeclared initial state")
	}
	labels := map[string]bool{"tau": true}
	for _, s := range alphabet(a.Signature) {
		labels[s.ID] = true
	}
	seen := map[edge]bool{}
	for _, e := range a.Edges {
		if !states[e.From] || !states[e.To] || !labels[e.Label] || seen[e] {
			return fmt.Errorf("invalid or duplicate edge %+v", e)
		}
		seen[e] = true
	}
	if a.Role == "M" || a.Role == "D" {
		for _, s := range a.States {
			for _, label := range alphabet(a.Signature) {
				restricted := label.ID == "reserve-capacity:AB" || label.ID == "apply:AB" || label.ID == "commit:BC"
				if a.Role == "D" {
					restricted = label.ID == "commit:AB" || label.ID == "apply:BC"
				}
				if !restricted {
					if !seen[edge{s.ID, label.ID, s.ID}] {
						return fmt.Errorf("external action veto: %s", label.ID)
					}
					for _, e := range a.Edges {
						if e.From == s.ID && e.Label == label.ID && e.To != s.ID {
							return fmt.Errorf("external action changes private state")
						}
					}
				}
			}
		}
	}
	return nil
}
func compatible(as []automaton) error {
	if len(as) < 2 {
		return fmt.Errorf("product requires factors")
	}
	for _, a := range as {
		if err := validate(a); err != nil {
			return fmt.Errorf("%s: %w", a.ID, err)
		}
		if !reflect.DeepEqual(as[0].Signature, a.Signature) {
			return fmt.Errorf("signature conflict")
		}
	}
	return nil
}
func targets(a automaton, from, label string) []string {
	out := []string{}
	for _, e := range a.Edges {
		if e.From == from && e.Label == label {
			out = append(out, e.To)
		}
	}
	slices.Sort(out)
	return out
}
func marked(a automaton, id string) bool {
	for _, s := range a.States {
		if s.ID == id {
			return s.Marked
		}
	}
	return false
}
func tuple(a automaton, id string) []string {
	for _, s := range a.States {
		if s.ID == id && len(s.Tuple) > 0 {
			return s.Tuple
		}
	}
	return []string{id}
}
