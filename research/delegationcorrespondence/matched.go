package delegationcorrespondence

import "encoding/json"

// MatchedBindingReport exposes existing raw checks without accepting any
// historical case protocol, endpoint application claim or producer verdict.
type MatchedBindingReport struct {
	Phases   []PhaseReport
	Bindings []EffectConcreteBinding
	Failure  *Failure
}

// CheckMatchedBinding reuses the unchanged independent authority semantics.
// files names initial, first-owner/post-open, stopped, reopened and final raw
// exports, controlled-revocation.json and endpoint-{a,b}-ready.json. manifest
// is the native reciprocal fixture manifest, not a historical case wrapper.
// Caller separately validates complete outer JSON and exact Open transitions.
func CheckMatchedBinding(files map[string][]byte, manifest []byte) (out MatchedBindingReport) {
	defer func() {
		if p := recover(); p != nil {
			if f, ok := p.(Failure); ok {
				out.Failure = &f
			} else {
				panic(p)
			}
		}
	}()
	var entry revocationManifest
	need(json.Unmarshal(manifest, &entry.Manifest.Reciprocal) == nil, "matched_manifest", "", "")
	r := entry.Manifest.Reciprocal
	need(r.A == revocationRuns[0] && r.B == revocationRuns[1] && r.BAService == `["a",1,4,0]`, "matched_principals", "", "")
	spec := BindingSpec{Principals: []BindingChain{{r.ABRequest, r.ABService, r.ABReply}, {r.BARequest, r.BAService, r.BAReply}}}
	for _, dir := range []string{"initial", "first-owner/post-open", "stopped", "reopened", "final"} {
		path := dir + "/inspect.json"
		v := CheckBindingRaw(files[path], spec)
		out.Phases = append(out.Phases, PhaseReport{Path: path, SHA256: sha(files[path]), Verdict: v})
		need(v.Valid, "matched_raw", path, fmtFailure(v.Failure))
	}
	if f := CheckContinuity(out.Phases, files); f != nil {
		panic(*f)
	}
	for i := 2; i < len(out.Phases); i++ {
		out.Bindings = append(out.Bindings, checkEffectBinding(out.Phases[i], files, entry, true, i))
	}
	return
}
func fmtFailure(f *Failure) string {
	if f == nil {
		return ""
	}
	return f.Edge + ": " + f.ID + ": " + f.Detail
}
