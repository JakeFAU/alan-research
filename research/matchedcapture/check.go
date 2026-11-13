package matchedcapture

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	dc "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
)

type Report struct {
	Version                                   string
	Supported                                 bool
	Seal, Source, Raw, Protocol, World, Model Verdict
	Inputs                                    map[string]Digest
	ProducerStatus                            string
	Equality                                  map[string]Comparison
	ReportEquality                            Comparison
	Worlds                                    map[string]WorldFacts
	RawBranches                               map[string]dc.MatchedBindingReport
	Models                                    map[string][]ModelCut
	Limitations                               []string
}

// CheckDirectory reads JSON/JSONL and opaque byte hashes only. Each verdict is
// independent; Supported requires all evidence and the complete observation.
func CheckDirectory(dir string) (out Report, err error) {
	out = Report{Version: Version, Equality: map[string]Comparison{}, Worlds: map[string]WorldFacts{}, RawBranches: map[string]dc.MatchedBindingReport{}, Models: map[string][]ModelCut{}, Limitations: []string{
		"Source-conditioned, unsigned producer evidence; seals establish internal byte consistency, not authenticity or proof that archived bytes came from these sources.",
		"Database files are hashed as opaque bytes. The frozen producer supplies the database/export connection; this checker does not independently inspect SQLite contents.",
		"Authority validation reuses the unchanged checker independent of the runtime producer; this is not a second implementation of that checker.",
		"One shared seed export with two pre-execution byte copies; four subsequent exports per branch. No separately inspected branch seed exports.",
		"Model projection is limited to stopped, reopened and final post-admission cuts with relative epochs 2->1 and 3->2. It is not full-export equality.",
		"Launch before matched ingress is an inferred prerequisite; endpoint callback occupancy is not established. Absence is bounded by held pre-call/death and complete endpoint process coverage.",
	}}
	f, d, e := loadDirectory(dir)
	if e != nil {
		return out, e
	}
	out.Inputs = d
	out.Seal = stage(func() { checkSeal(f, d); checkMembers(f, d) })
	var p pair
	declaration := stage(func() {
		p = read[pair](f, "pair.json")
		require(p.Version == Version, "pair.json", "version")
		require(p.Status == "complete" || p.Status == "failed", "pair.json", "status")
		require(p.Mode == "normal" || p.Mode == "race", "pair.json", "mode")
		out.ProducerStatus = p.Status
	})
	out.Source = stage(func() { require(declaration.Valid, "source", "pair declaration invalid"); checkSource(f, p) })
	rawVerdicts := []Verdict{stage(func() { checkSeedProvision(f) })}
	for _, branch := range []string{"zero", "applied"} {
		rawVerdicts = append(rawVerdicts, stage(func() { out.RawBranches[branch] = checkStore(f, branch) }))
	}
	out.Raw = combine(rawVerdicts...)
	protocols := map[string]branchProtocol{}
	pv := []Verdict{stage(func() { checkArchiveShape(d) })}
	for _, branch := range []string{"zero", "applied"} {
		pv = append(pv, stage(func() {
			require(declaration.Valid, "protocol", "pair declaration invalid")
			protocols[branch] = checkProtocol(f, p, branch)
		}))
	}
	pv = append(pv, stage(func() {
		require(len(protocols) == 2, "protocol", "branch chronology prerequisite")
		checkIsolation(protocols["zero"], protocols["applied"])
		require(same(p.Completed, []string{"zero", "applied"}) && p.Status == "complete", "pair.json", "incomplete attempt")
	}))
	out.Protocol = combine(pv...)
	wv := []Verdict{}
	for _, branch := range []string{"zero", "applied"} {
		wv = append(wv, stage(func() { out.Worlds[branch] = checkWorld(f, d, branch) }))
	}
	out.World = combine(wv...)
	// Compare raw exports even when a semantic or causal verdict fails. The seed
	// comparison is explicitly the sealed seed and its declared opaque copies.
	out.Equality["seed-copies"] = seedEquality(f, d)
	for _, cut := range []string{"first-owner/post-open", "stopped", "reopened", "final"} {
		out.Equality[cut] = comparisonFiles(f, "zero/"+cut+"/inspect.json", "applied/"+cut+"/inspect.json")
	}
	out.ReportEquality = comparisonFiles(f, "zero/run.json", "applied/run.json")
	out.Model = stage(func() {
		require(out.Raw.Valid && out.Protocol.Valid && out.World.Valid, "model", "raw/protocol/world prerequisites not certified")
		for _, branch := range []string{"zero", "applied"} {
			out.Models[branch] = mapModels(branch, out.RawBranches[branch])
		}
		for i := 0; i < 3; i++ {
			require(out.Models["zero"][i].Store == out.Models["applied"][i].Store, "model", "projected mismatch")
		}
	})
	equal := out.ReportEquality.StructuralEqual
	for _, c := range out.Equality {
		equal = equal && c.StructuralEqual
	}
	out.Supported = out.Seal.Valid && out.Source.Valid && out.Raw.Valid && out.Protocol.Valid && out.World.Valid && out.Model.Valid && equal
	return out, nil
}
func combine(vs ...Verdict) Verdict {
	out := Verdict{Valid: true, Reasons: []string{}}
	for _, v := range vs {
		out.Valid = out.Valid && v.Valid
		out.Reasons = append(out.Reasons, v.Reasons...)
	}
	return out
}
func comparisonFiles(f map[string][]byte, a, b string) Comparison {
	_, existsA := f[a]
	_, existsB := f[b]
	if !existsA || !existsB {
		return Comparison{Paths: []string{a + " <-> " + b + ": missing export"}}
	}
	r, e := CompareJSON(f[a], f[b])
	if e != nil {
		r.Paths = []string{a + " <-> " + b + ": " + e.Error()}
		return r
	}
	return r
}
func seedEquality(f map[string][]byte, d map[string]Digest) Comparison {
	out := Comparison{Paths: []string{}}
	v := stage(func() {
		sets := read[map[string]map[string]Digest](f, "seed/copies.json")
		require(len(sets) == 3 && len(sets["seed"]) > 0, "seed/copies.json", "three complete file sets")
		actual := map[string]Digest{}
		for p, h := range d {
			if strings.HasPrefix(p, "seed/initial/runtime.sqlite") {
				actual[strings.TrimPrefix(p, "seed/initial/")] = h
			}
		}
		require(actual["runtime.sqlite"].Bytes > 0, "seed/initial", "missing SQLite bytes")
		for n := range actual {
			require(n == "runtime.sqlite" || n == "runtime.sqlite-wal" || n == "runtime.sqlite-shm", "seed/initial", "unsupported DB companion")
		}
		for _, label := range []string{"seed", "zero", "applied"} {
			require(same(actual, sets[label]), "seed/copies.json."+label, "does not equal complete sealed seed file set")
		}
	})
	out.StructuralEqual = v.Valid
	out.ByteEqual = v.Valid
	out.Paths = append(out.Paths, v.Reasons...)
	return out
}
func checkMembers(f map[string][]byte, d map[string]Digest) {
	dirs := []string{"seed/initial"}
	for _, branch := range []string{"zero", "applied"} {
		for _, cut := range []string{"first-owner/post-open", "stopped", "reopened", "final"} {
			dirs = append(dirs, branch+"/"+cut)
		}
	}
	for _, dir := range dirs {
		for _, name := range []string{"inspect.json", "typed-inspection.json", "outcomes.json", "replay.json", "cost-sidecar.json", "history-000.json", "runtime.sqlite"} {
			require(d[dir+"/"+name].Bytes > 0, dir+"/"+name, "required capture member")
		}
		outcomes := read[map[string]string](f, dir+"/outcomes.json")
		require(len(outcomes) == 3 && outcomes["Inspect"] == "<nil>" && outcomes["InspectExecution"] == "<nil>" && outcomes["Replay"] == "<nil>", dir, "capture export errors")
		replay := read[map[string]json.RawMessage](f, dir+"/replay.json")
		require(decode[string](replay["Error"], dir) == "<nil>", dir, "replay export error")
	}
	// Unknown extra JSON remains sealed and parsed, never silently treated as an
	// authoritative receipt or a new control action. Failed working files/aborts
	// fail protocol; rollback journals are outside the closed-copy contract.
	names := []string{}
	for n := range d {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		require(!strings.HasSuffix(n, "-journal"), n, "unsupported rollback journal")
	}
	require(d["seed/provision.json"].Bytes > 0 && d["seed/manifest.json"].Bytes > 0, "seed", "missing source provision")
	_, e := parseJSON(f["seed/provision.json"])
	require(e == nil, "seed/provision.json", fmt.Sprint(e))
}
