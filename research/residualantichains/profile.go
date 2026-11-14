package residualantichains

type profileNode struct {
	Left        []string
	Right       []string
	LeftFamily  family
	RightFamily family
	LeftToRight []coverage
	RightToLeft []coverage
	FULL        bool
	MIN         bool
	MAX         bool
	Arcs        []subsetArc
}
type profileComparison struct {
	Left            string
	Right           string
	InitialLanguage comparison
	LeftWithin      []residualComparison
	RightWithin     []residualComparison
	Residuals       []residualComparison
	Pairs           []profileNode
	FULL            bool
	MIN             bool
	MAX             bool
	J               bool
	FULLWord        *[]string
	MINWord         *[]string
	MAXWord         *[]string
}

func compareProfiles(a, b automaton) profileComparison {
	c := profileComparison{Left: a.ID, Right: b.ID, InitialLanguage: compareLanguages(a, b),
		LeftWithin: residuals(a, a), RightWithin: residuals(b, b), Residuals: residuals(a, b), Pairs: []profileNode{}, FULL: true, MIN: true, MAX: true}
	left, right, cross := residualIndex(c.LeftWithin), residualIndex(c.RightWithin), residualIndex(c.Residuals)
	// The initial language graph is also the complete weak-reach graph. Its
	// acceptance bits are irrelevant to profile comparison; its closures are not.
	words := map[int][]string{0: {}}
	for i, p := range c.InitialLanguage.Pairs {
		lf, rf := reduce(p.Left, left), reduce(p.Right, right)
		n := profileNode{Left: p.Left, Right: p.Right, LeftFamily: lf, RightFamily: rf,
			LeftToRight: cover(a, b, lf, rf, cross, false), RightToLeft: cover(b, a, rf, lf, cross, true),
			FULL: equalFamilies(representatives(lf), representatives(rf), cross),
			MIN:  equalFamilies(lf.MIN, rf.MIN, cross), MAX: equalFamilies(lf.MAX, rf.MAX, cross), Arcs: p.Arcs}
		if c.FULL && !n.FULL {
			w := words[i]
			c.FULLWord = &w
		}
		c.FULL = c.FULL && n.FULL
		if c.MIN && !n.MIN {
			w := words[i]
			c.MINWord = &w
		}
		c.MIN = c.MIN && n.MIN
		if c.MAX && !n.MAX {
			w := words[i]
			c.MAXWord = &w
		}
		c.MAX = c.MAX && n.MAX
		for _, arc := range p.Arcs {
			if _, ok := words[arc.To]; !ok {
				words[arc.To] = append(append([]string{}, words[i]...), arc.Label)
			}
		}
		c.Pairs = append(c.Pairs, n)
	}
	c.J = c.MIN && c.InitialLanguage.Equal
	return c
}
