package matchedcapture

import (
	"encoding/json"
	"testing"
)

// Each mutation is made in a scratch dossier and its member seal is recomputed.
func TestReviewR1EndpointCoverage(t *testing.T) {
	for _, endpoint := range []string{"endpoint-a", "endpoint-b"} {
		for _, cut := range []string{"stopped", "final"} {
			for _, mutation := range []string{"delete-coverage", "null-total", "null-contacts", "null-applications"} {
				t.Run(endpoint+"/"+cut+"/"+mutation, func(t *testing.T) {
					f := fixture(t)
					p := "zero/" + endpoint + "/" + cut + ".json"
					editJSON(f, p, func(d map[string]any) {
						s := d["snapshot"].(map[string]any)
						switch mutation {
						case "delete-coverage":
							for _, k := range []string{"Requests", "Applications", "Total", "Contacts"} {
								delete(s, k)
							}
						case "null-total":
							s["Total"] = nil
						case "null-contacts":
							s["Contacts"] = nil
						case "null-applications":
							s["Applications"] = nil
						}
					})
					r, e := CheckDirectory(sealFixture(t, f))
					if e != nil || !r.Seal.Valid || r.World.Valid || r.Supported {
						t.Fatalf("incomplete endpoint coverage: %v seal=%+v world=%+v", e, r.Seal, r.World)
					}
				})
			}
		}
	}
}
func TestReviewR1CompleteReports(t *testing.T) {
	for _, field := range []string{"Steps", "Dispatches", "Quiescent", "OpenRuns", "UnresolvedIntents", "PendingDeliveries", "UnknownExternal", "UnknownAdapters", "Roots", "RejectedAdmissions"} {
		for _, mutation := range []string{"delete", "null", "wrong-type"} {
			t.Run(field+"/"+mutation, func(t *testing.T) {
				f := fixture(t)
				for _, branch := range []string{"zero", "applied"} {
					editJSON(f, branch+"/run.json", func(d map[string]any) {
						r := d["Report"].(map[string]any)
						switch mutation {
						case "delete":
							delete(r, field)
						case "null":
							r[field] = nil
						case "wrong-type":
							r[field] = map[string]any{}
						}
					})
				}
				r, e := CheckDirectory(sealFixture(t, f))
				if e != nil || !r.Seal.Valid || r.Raw.Valid || !r.ReportEquality.StructuralEqual || r.Supported {
					t.Fatalf("equal incomplete Run report: %v raw=%+v equality=%+v", e, r.Raw, r.ReportEquality)
				}
			})
		}
	}
}
func TestReviewR1NestedReportSchema(t *testing.T) {
	for _, which := range []string{"root-charge", "intent-ordinal", "null-root", "null-open-run"} {
		t.Run(which, func(t *testing.T) {
			f := fixture(t)
			for _, branch := range []string{"zero", "applied"} {
				editJSON(f, branch+"/run.json", func(d map[string]any) {
					r := d["Report"].(map[string]any)
					switch which {
					case "root-charge":
						delete(r["Roots"].([]any)[0].(map[string]any)["Charges"].(map[string]any), "Effects")
					case "intent-ordinal":
						delete(r["UnresolvedIntents"].([]any)[0].(map[string]any), "Ordinal")
					case "null-root":
						r["Roots"].([]any)[0] = nil
					case "null-open-run":
						r["OpenRuns"].([]any)[0] = nil
					}
				})
			}
			r, e := CheckDirectory(sealFixture(t, f))
			if e != nil || !r.Seal.Valid || r.Raw.Valid || !r.ReportEquality.StructuralEqual || r.Supported {
				t.Fatalf("incomplete nested report: %v raw=%+v", e, r.Raw)
			}
		})
	}
}
func TestReviewR1AuthorizedEvent(t *testing.T) {
	for _, tc := range []struct {
		field string
		value any
	}{{"Subject", "wrong"}, {"Generation", 2}, {"Epoch", 999}, {"Sequence", 999}, {"Phase", "wrong"}, {"Root", "wrong"}, {"AuthorityAdmissionOrdinal", 1}, {"InFlight", []string{"wrong"}}, {"InFlight", []string{}}, {"Generation", nil}, {"Intent", "wrong"}, {"Attempt", "wrong"}} {
		t.Run(tc.field+"/"+string(jsonBytes(tc.value)), func(t *testing.T) {
			f := fixture(t)
			p := "zero/first-owner/boundaries.jsonl"
			rows := lines[map[string]any](f, p)
			f[p] = nil
			for _, row := range rows {
				if row["Boundary"] == "external.authorized" {
					row[tc.field] = tc.value
				}
				f[p] = append(f[p], append(jsonBytes(row), '\n')...)
			}
			r, e := CheckDirectory(sealFixture(t, f))
			if e != nil || !r.Seal.Valid || r.Protocol.Valid || r.Supported {
				t.Fatalf("contradictory native authorization: %v protocol=%+v", e, r.Protocol)
			}
		})
	}
}
func TestReviewR1EffectiveSelection(t *testing.T) {
	selected := "^TestPhase3B7MatchedCapture$"
	cases := []struct {
		name  string
		args  []string
		valid bool
	}{
		{"single-equals", []string{"-test.run=" + selected}, true},
		{"double-equals", []string{"--test.run=" + selected}, true},
		{"single-separated", []string{"-test.run", selected}, true},
		{"double-separated", []string{"--test.run", selected}, true},
		{"normal-driver", []string{"-test.paniconexit0", "-test.v=test2json", "-test.timeout=10m0s", "-test.count=1", "-test.run=" + selected}, true},
		{"terminal-after-selector", []string{"-test.run=" + selected, "--"}, true},
		{"override", []string{"-test.run=" + selected, "-test.run=^OtherTest$"}, false},
		{"override-double-separated", []string{"-test.run=" + selected, "--test.run", "^OtherTest$"}, false},
		{"duplicate-identical", []string{"-test.run=" + selected, "--test.run=" + selected}, false},
		{"terminated-before", []string{"--", "-test.run=" + selected}, false},
		{"positional-before", []string{"positional", "-test.run=" + selected}, false},
		{"hidden-override-after-termination", []string{"-test.run=" + selected, "--", "-test.run=^OtherTest$"}, false},
		{"selector-as-another-value", []string{"-test.testlogfile", "-test.run=" + selected}, false},
		{"missing-value", []string{"-test.run"}, false},
		{"repeat-count", []string{"-test.run=" + selected, "-test.count=2"}, false},
		{"zero-count", []string{"-test.run=" + selected, "-test.count=0"}, false},
		{"skip-selected", []string{"-test.run=" + selected, "-test.skip=" + selected}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := fixture(t)
			editJSON(f, "source.json", func(d map[string]any) { d["test_argv"] = append([]string{"/synthetic/owner.test"}, tc.args...) })
			r, e := CheckDirectory(sealFixture(t, f))
			if e != nil || !r.Seal.Valid || r.Source.Valid != tc.valid || r.Supported != tc.valid {
				t.Fatalf("effective selection valid=%v: %v source=%+v supported=%v", tc.valid, e, r.Source, r.Supported)
			}
		})
	}
}

func TestReviewR1CoherentTruncation(t *testing.T) {
	for _, kind := range []string{"endpoint-a", "endpoint-b", "run"} {
		t.Run(kind, func(t *testing.T) {
			f := fixture(t)
			if kind == "run" {
				for _, branch := range []string{"zero", "applied"} {
					f[branch+"/run.json"] = jsonBytes(map[string]any{"Error": "<nil>", "Report": map[string]any{"UnknownExternal": []string{service}}})
				}
			} else {
				for _, cut := range []string{"stopped", "final"} {
					editJSON(f, "zero/"+kind+"/"+cut+".json", func(d map[string]any) {
						s := d["snapshot"].(map[string]any)
						for _, k := range []string{"Requests", "Applications", "Total", "Contacts"} {
							delete(s, k)
						}
					})
				}
			}
			r, e := CheckDirectory(sealFixture(t, f))
			if e != nil || !r.Seal.Valid || r.Supported || (kind == "run" && (r.Raw.Valid || !r.ReportEquality.StructuralEqual)) || (kind != "run" && r.World.Valid) {
				t.Fatalf("coherent truncation accepted: %v raw=%+v world=%+v equality=%+v", e, r.Raw, r.World, r.ReportEquality)
			}
		})
	}
}

func TestReviewR1RequiredSnapshotMembers(t *testing.T) {
	f := fixture(t)
	for _, branch := range []string{"zero", "applied"} {
		path := branch + "/endpoint-b/final.json"
		outer := read[map[string]json.RawMessage](f, path)
		snapshot := decode[map[string]json.RawMessage](outer["snapshot"], path)
		for field := range snapshot {
			for _, mutation := range []string{"delete", "wrong-type", "null"} {
				// The native zero Requests slice is nil, and its explicit null is valid.
				if field == "Requests" && mutation == "null" {
					continue
				}
				t.Run(branch+"/"+field+"/"+mutation, func(t *testing.T) {
					s := decode[map[string]json.RawMessage](outer["snapshot"], path)
					switch mutation {
					case "delete":
						delete(s, field)
					case "wrong-type":
						s[field] = json.RawMessage(`{}`)
					case "null":
						s[field] = json.RawMessage(`null`)
					}
					v := stage(func() {
						got := decodeComplete[endpointSnapshot](jsonBytes(s), path)
						require(got.Applications != nil, path, "native Applications array")
					})
					if v.Valid {
						t.Fatal("missing or mistyped native member accepted")
					}
				})
			}
		}
	}
}

func TestReviewR1RequiredAuthorizationMembers(t *testing.T) {
	for _, field := range []string{"Subject", "Generation", "Epoch", "Sequence", "Phase", "Intent", "Attempt", "InFlight"} {
		t.Run(field, func(t *testing.T) {
			f := fixture(t)
			path := "applied/first-owner/boundaries.jsonl"
			rows := lines[map[string]any](f, path)
			f[path] = nil
			for _, row := range rows {
				if row["Boundary"] == "external.authorized" {
					delete(row, field)
				}
				f[path] = append(f[path], append(jsonBytes(row), '\n')...)
			}
			r, e := CheckDirectory(sealFixture(t, f))
			if e != nil || !r.Seal.Valid || r.Protocol.Valid || r.Supported {
				t.Fatalf("missing authorization member accepted: %v %+v", e, r.Protocol)
			}
		})
	}
}
