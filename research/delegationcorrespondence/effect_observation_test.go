package delegationcorrespondence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func effectObservationFixture(t *testing.T, mode, name string) string {
	t.Helper()
	root := os.Getenv("ALAN_EFFECT_OBSERVATION_CORPUS")
	if root == "" {
		t.Skip("ALAN_EFFECT_OBSERVATION_CORPUS is not set")
	}
	data, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]struct{ Cases map[string]string }
	if err = json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	member := manifest["raw-"+mode+".tar.gz"].Cases[name]
	if member == "" {
		t.Fatal("missing corpus case", name)
	}
	return filepath.Join(root, "raw-"+mode, member)
}

func TestEffectObservationPositive(t *testing.T) {
	for _, mode := range []string{"normal", "race"} {
		for _, name := range []string{"N07", "applied-unknown-revoked-before-fresh-gate"} {
			t.Run(mode+"/"+name, func(t *testing.T) {
				r, e := CheckEffectObservationDirectory(effectObservationFixture(t, mode, name))
				if e != nil {
					t.Fatal(e)
				}
				if !r.RawValid || !r.ContinuityValid || !r.Physical.Valid || !r.ProtocolValid || !r.Supported {
					t.Fatalf("semantic assertion: genuine sealed cut must be supported: raw=%v continuity=%v physical=%v protocol=%v supported=%v reasons=%v", r.RawValid, r.ContinuityValid, r.Physical.Valid, r.ProtocolValid, r.Supported, r.Reasons)
				}
			})
		}
	}
}

func TestEffectObservationCausalAckRequired(t *testing.T) {
	parent := effectObservationFixture(t, "normal", "applied-unknown-revoked-before-fresh-gate")
	var ack map[string]json.RawMessage
	if err := json.Unmarshal(readRevocationTestFile(t, parent, "acknowledgments.json"), &ack); err != nil {
		t.Fatal(err)
	}
	delete(ack, "Endpoint")
	dir := editRevocationDirectory(t, parent, map[string][]byte{"acknowledgments.json": bytesJSON(ack)})
	r, e := CheckEffectObservationDirectory(dir)
	if e != nil {
		t.Fatal(e)
	}
	if r.Supported {
		t.Fatal("semantic assertion: final application without endpoint causal acknowledgment must not certify application-before-loss")
	}
}

func TestEffectObservationCutSeparation(t *testing.T) {
	for _, name := range []string{"N07", "applied-unknown-revoked-before-fresh-gate"} {
		r, err := CheckEffectObservationDirectory(effectObservationFixture(t, "normal", name))
		if err != nil || !r.Supported {
			t.Fatalf("baseline %s: %v %v", name, err, r.Reasons)
		}
		if !r.Producer.SealedBytesMatch {
			t.Fatal("genuine dossier does not match sealed member map")
		}
		var prior []EffectWitnessStep
		for i, c := range r.Captures {
			if name != "N07" && i < 2 {
				if c.Supported || c.Reason != "outside_domain_pre_admission" || c.Store != nil || c.Dossier != nil || len(c.Witness) != 0 {
					t.Fatal("pre-admission cut gained model facts")
				}
				continue
			}
			if c.Audit == nil || c.Store == nil || c.Dossier == nil || len(prior) > len(c.Witness) || !equal(prior, c.Witness[:len(prior)]) && len(prior) > 0 {
				t.Fatal("witness does not coherently extend")
			}
			prior = c.Witness
			if name == "N07" && i < 2 {
				if c.Dossier.Ingress || c.Dossier.Applied != "" || c.Store.Retained != "" || c.Store.Consumed != "" || c.Certificate.FinalCorroboration == nil {
					t.Fatal("semantic assertion: final endpoint or custody evidence leaked backward")
				}
			}
			if name != "N07" {
				if !c.Dossier.Ingress || c.Dossier.Applied != "v1" || c.Store.Retained != "" || c.Store.Consumed != "" || c.Binding.Attempt.Request.Epoch != 2 {
					t.Fatal("loss cut erased world truth or imported a receipt")
				}
			}
		}
		if name != "N07" && r.Legacy.Projection.Supported {
			t.Fatal("legacy loss exclusion changed")
		}
	}
}

func TestEffectObservationReencodedPositive(t *testing.T) {
	parent := effectObservationFixture(t, "normal", effectLossName)
	// A coherent editable dossier may preserve semantics while changing export
	// formatting. Its conditional relation can pass while archive membership fails.
	data := effectAttackJSON(t, parent, "reopened/inspect.json", func(map[string]any) {})
	r, err := CheckEffectObservationDirectory(editRevocationDirectory(t, parent, map[string][]byte{"reopened/inspect.json": data}))
	if err != nil || !r.Supported || r.Producer.SealedBytesMatch {
		t.Fatalf("reencoded lawful control: err=%v supported=%v sealed=%v reasons=%v", err, r.Supported, r.Producer.SealedBytesMatch, r.Reasons)
	}
}

func TestEffectObservationRenewalCoverage(t *testing.T) {
	parent := effectObservationFixture(t, "normal", effectLossName)
	before := readRevocationTestFile(t, parent, "stopped/inspect.json")
	checkEffectExactRenewal(before, readRevocationTestFile(t, parent, "reopened/inspect.json"))
	for _, attack := range []string{"run-missing", "run-extra", "one-epoch", "run-field", "table-missing", "table-extra", "non-run-row"} {
		t.Run(attack, func(t *testing.T) {
			var after RawSnapshot
			if err := json.Unmarshal(readRevocationTestFile(t, parent, "reopened/inspect.json"), &after); err != nil {
				t.Fatal(err)
			}
			switch attack {
			case "run-missing":
				after.Tables["runs"] = after.Tables["runs"][:1]
			case "run-extra":
				after.Tables["runs"] = append(after.Tables["runs"], after.Tables["runs"][0])
			case "one-epoch":
				after.Tables["runs"][0]["epoch"] = bytesJSON(2)
			case "run-field":
				after.Tables["runs"][0]["sequence"] = bytesJSON(99)
			case "table-missing":
				delete(after.Tables, "adapter_attempts")
			case "table-extra":
				after.Tables["other"] = []Row{}
			case "non-run-row":
				after.Tables["subjects"][0]["revision"] = bytesJSON(2)
			}
			rejected := false
			func() {
				defer func() {
					if p := recover(); p != nil {
						if f, ok := p.(Failure); ok && f.Edge == "effect_exact_renewal" {
							rejected = true
						} else {
							panic(p)
						}
					}
				}()
				checkEffectExactRenewal(before, bytesJSON(after))
			}()
			if !rejected {
				t.Fatal("semantic assertion: exact renewal helper accepted", attack)
			}
		})
	}
}
