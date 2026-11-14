package residualantichains

import "slices"

type residualComparison struct {
	LeftState  string
	RightState string
	Language   comparison
}
type languageClass struct {
	Representative  string
	Members         []string
	ProperSubsets   []string
	ProperSupersets []string
}
type family struct {
	Classes []languageClass
	MIN     []string
	MAX     []string
}

func residuals(a, b automaton) []residualComparison {
	out := []residualComparison{}
	for _, x := range a.States {
		for _, y := range b.States {
			aa, bb := a, b
			aa.Initial = x.ID
			bb.Initial = y.ID
			out = append(out, residualComparison{LeftState: x.ID, RightState: y.ID, Language: compareLanguages(aa, bb)})
		}
	}
	return out
}
func residualIndex(rs []residualComparison) map[string]comparison {
	out := map[string]comparison{}
	for _, r := range rs {
		out[key([]string{r.LeftState, r.RightState})] = r.Language
	}
	return out
}
func reduce(reached []string, within map[string]comparison) family {
	f := family{Classes: []languageClass{}, MIN: []string{}, MAX: []string{}}
	for _, s := range reached {
		found := false
		for i, c := range f.Classes {
			if within[key([]string{s, c.Representative})].Equal {
				f.Classes[i].Members = append(f.Classes[i].Members, s)
				found = true
				break
			}
		}
		if !found {
			f.Classes = append(f.Classes, languageClass{Representative: s, Members: []string{s}, ProperSubsets: []string{}, ProperSupersets: []string{}})
		}
	}
	for i, c := range f.Classes {
		for _, d := range f.Classes {
			r := within[key([]string{c.Representative, d.Representative})]
			if r.LeftSubset && !r.RightSubset {
				f.Classes[i].ProperSupersets = append(f.Classes[i].ProperSupersets, d.Representative)
			}
			if r.RightSubset && !r.LeftSubset {
				f.Classes[i].ProperSubsets = append(f.Classes[i].ProperSubsets, d.Representative)
			}
		}
		if len(f.Classes[i].ProperSubsets) == 0 {
			f.MIN = append(f.MIN, c.Representative)
		}
		if len(f.Classes[i].ProperSupersets) == 0 {
			f.MAX = append(f.MAX, c.Representative)
		}
	}
	return f
}
func representatives(f family) []string {
	out := []string{}
	for _, c := range f.Classes {
		out = append(out, c.Representative)
	}
	return out
}
func equalFamilies(l, r []string, cross map[string]comparison) bool {
	for _, s := range l {
		found := false
		for _, t := range r {
			found = found || cross[key([]string{s, t})].Equal
		}
		if !found {
			return false
		}
	}
	for _, t := range r {
		found := false
		for _, s := range l {
			found = found || cross[key([]string{s, t})].Equal
		}
		if !found {
			return false
		}
	}
	return true
}

type outsideWord struct {
	Source string
	Word   []string
}
type sourceHit struct {
	Source string
	Hit    bool
}
type localFailure struct {
	Outside    []outsideWord
	H          [][]string
	Hits       []sourceHit
	TargetMiss bool
}
type coverage struct {
	Target  string
	Sources []string
	Failure *localFailure
}

// Coverage direction is source <= target. A failed target gets one shortest
// outside word per source class, not a separator from the source union.
func cover(source, target automaton, sf, tf family, cross map[string]comparison, reverse bool) []coverage {
	out := []coverage{}
	for _, t := range tf.Classes {
		c := coverage{Target: t.Representative, Sources: []string{}}
		fail := localFailure{Outside: []outsideWord{}, H: [][]string{}, Hits: []sourceHit{}, TargetMiss: true}
		seen := map[string]bool{}
		for _, s := range sf.Classes {
			pair := []string{s.Representative, t.Representative}
			if reverse {
				slices.Reverse(pair)
			}
			r := cross[key(pair)]
			sub, w := r.LeftSubset, r.LeftOutside
			if reverse {
				sub, w = r.RightSubset, r.RightOutside
			}
			if sub {
				c.Sources = append(c.Sources, s.Representative)
				continue
			}
			word := append([]string{}, (*w)...)
			fail.Outside = append(fail.Outside, outsideWord{Source: s.Representative, Word: word})
			if !seen[key(word)] {
				seen[key(word)] = true
				fail.H = append(fail.H, word)
			}
		}
		if len(c.Sources) == 0 {
			for _, s := range sf.Classes {
				a := source
				a.Initial = s.Representative
				hit := false
				for _, w := range fail.H {
					hit = hit || accepts(a, w)
				}
				fail.Hits = append(fail.Hits, sourceHit{Source: s.Representative, Hit: hit})
			}
			a := target
			a.Initial = t.Representative
			for _, w := range fail.H {
				fail.TargetMiss = fail.TargetMiss && !accepts(a, w)
			}
			c.Failure = &fail
		}
		out = append(out, c)
	}
	return out
}
