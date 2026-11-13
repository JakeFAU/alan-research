package delegationcorrespondence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMatchedAdapter(t *testing.T) {
	for _, mode := range []string{"normal", "race"} {
		t.Run(mode, func(t *testing.T) {
			dir := effectObservationFixture(t, mode, effectLossName)
			files := map[string][]byte{}
			for _, name := range []string{"initial/inspect.json", "first-owner/post-open/inspect.json", "stopped/inspect.json", "reopened/inspect.json", "final/inspect.json", "controlled-revocation.json", "endpoint-a-ready.json", "endpoint-b-ready.json"} {
				b, e := os.ReadFile(filepath.Join(dir, name))
				if e != nil {
					t.Fatal(e)
				}
				files[name] = b
			}
			data, e := os.ReadFile(filepath.Join(dir, "manifest.json"))
			if e != nil {
				t.Fatal(e)
			}
			var manifest struct{ Reciprocal json.RawMessage }
			if e = json.Unmarshal(data, &manifest); e != nil {
				t.Fatal(e)
			}
			result := CheckMatchedBinding(files, manifest.Reciprocal)
			if result.Failure != nil || len(result.Bindings) != 3 {
				t.Fatalf("adapter rejected historical raw bindings: %+v", result.Failure)
			}
			var native struct{ ABRequest, ABService, ABReply, BARequest, BAService, BAReply string }
			json.Unmarshal(manifest.Reciprocal, &native)
			spec := BindingSpec{Principals: []BindingChain{{native.ABRequest, native.ABService, native.ABReply}, {native.BARequest, native.BAService, native.BAReply}}}
			for _, phase := range result.Phases {
				direct := CheckBindingRaw(files[phase.Path], spec)
				if !equal(direct, phase.Verdict) {
					t.Fatal("adapter changed raw verdict", phase.Path)
				}
			}
			var bad map[string]json.RawMessage
			json.Unmarshal(manifest.Reciprocal, &bad)
			bad["BARequest"] = json.RawMessage(`"wrong"`)
			if attack := CheckMatchedBinding(files, bytesJSON(bad)); attack.Failure == nil {
				t.Fatal("adapter lost exact original identity")
			}
		})
	}
}
