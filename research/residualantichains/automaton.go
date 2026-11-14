// Package residualantichains decides finite residual-language profiles over
// pinned predecessor automata. It does not implement runtime refinement.
package residualantichains

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
)

var labels = []string{"commit:AB", "reserve-capacity:AB", "apply:AB", "consume:AB",
	"commit:BC", "reserve-capacity:BC", "apply:BC", "consume:BC"}

type symbol struct {
	ID     string
	Action json.RawMessage
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
	Signature json.RawMessage
	Initial   string
	States    []state
	Edges     []edge
}

// Signatures are opaque exact predecessor descriptors. Their trusted bytes are
// authenticated before decoding; no old constructor or outcome oracle is called.
func validate(a automaton) error {
	if a.ID == "" {
		return fmt.Errorf("empty automaton id")
	}
	if !slices.Contains([]string{"K", "M", "D", "product", "generic"}, a.Role) {
		return fmt.Errorf("unknown role %q", a.Role)
	}
	states := map[string]bool{}
	for _, s := range a.States {
		if s.ID == "" || states[s.ID] {
			return fmt.Errorf("invalid state %q", s.ID)
		}
		states[s.ID] = true
	}
	if !states[a.Initial] {
		return fmt.Errorf("undeclared initial")
	}
	seen := map[edge]bool{}
	for _, e := range a.Edges {
		if !states[e.From] || !states[e.To] || (e.Label != "tau" && !slices.Contains(labels, e.Label)) || seen[e] {
			return fmt.Errorf("invalid edge %+v", e)
		}
		seen[e] = true
	}
	if a.Role != "M" && a.Role != "D" {
		return nil
	}
	for _, s := range a.States {
		for _, label := range labels {
			restricted := label == "reserve-capacity:AB" || label == "apply:AB" || label == "commit:BC"
			if a.Role == "D" {
				restricted = label == "commit:AB" || label == "apply:BC"
			}
			if restricted {
				continue
			}
			ts := targets(a, s.ID, label)
			if len(ts) != 1 || ts[0] != s.ID {
				return fmt.Errorf("external action veto or state change: %s", label)
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
		if !bytes.Equal(as[0].Signature, a.Signature) {
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
