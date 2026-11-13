package matchedcapture

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func editJSON(f map[string][]byte, p string, fn func(map[string]any)) {
	var v map[string]any
	if e := json.Unmarshal(f[p], &v); e != nil {
		panic(e)
	}
	fn(v)
	f[p] = append(jsonBytes(v), '\n')
}
func editProtocol(f map[string][]byte, branch string, fn func([]record) []record) {
	p := branch + "/protocol.jsonl"
	r := lines[record](f, p)
	r = fn(r)
	f[p] = nil
	for i := range r {
		r[i].Ordinal = i + 1
		f[p] = append(f[p], append(jsonBytes(r[i]), '\n')...)
	}
}
func changeRecord(f map[string][]byte, branch, action string, fn func(map[string]any)) {
	editProtocol(f, branch, func(r []record) []record {
		for i := range r {
			if r[i].Action == action {
				var d map[string]any
				json.Unmarshal(r[i].Data, &d)
				fn(d)
				r[i].Data = jsonBytes(d)
				break
			}
		}
		return r
	})
}
func TestSyntheticPair(t *testing.T) {
	r, e := CheckDirectory(sealFixture(t, fixture(t)))
	if e != nil || !r.Supported {
		t.Fatalf("synthetic lawful control: %v %+v", e, r)
	}
	for _, v := range r.Equality {
		if !v.StructuralEqual {
			t.Fatal(v)
		}
	}
	if r.Worlds["zero"].Applications != 0 || r.Worlds["applied"].Applications != 1 {
		t.Fatal("opposite world facts lost")
	}
}
func TestResealedAttacks(t *testing.T) {
	tests := []struct {
		name    string
		edit    func(map[string][]byte)
		verdict string
	}{
		{"missing-precall", func(f map[string][]byte) { delete(f, "zero/first-owner/pre-call.json") }, "protocol"},
		{"wrong-precall", func(f map[string][]byte) {
			changeRecord(f, "zero", "pre-call-ack", func(d map[string]any) { d["event"].(map[string]any)["Attempt"] = "wrong" })
		}, "protocol"},
		{"missing-scheduler", func(f map[string][]byte) {
			editProtocol(f, "zero", func(r []record) []record {
				for i, v := range r {
					if v.Action == "scheduler-ack" {
						return append(r[:i], r[i+1:]...)
					}
				}
				return r
			})
		}, "protocol"},
		{"wrong-scheduler", func(f map[string][]byte) {
			changeRecord(f, "zero", "scheduler-ack", func(d map[string]any) { d["event"].(map[string]any)["InFlight"] = []string{"wrong"} })
		}, "protocol"},
		{"premature-worker", func(f map[string][]byte) {
			editProtocol(f, "applied", func(r []record) []record {
				for i, v := range r {
					if v.Action == "scheduler-ack" {
						r[i], r[i+1] = r[i+1], r[i]
						break
					}
				}
				return r
			})
		}, "protocol"},
		{"zero-worker-release", func(f map[string][]byte) {
			f["zero/first-owner/worker-release.json"] = f["zero/first-owner/pre-call.json"]
		}, "protocol"},
		{"wrong-kill", func(f map[string][]byte) {
			changeRecord(f, "zero", "kill-request", func(d map[string]any) { d["signal"] = "terminated" })
		}, "protocol"},
		{"wrong-wait", func(f map[string][]byte) {
			changeRecord(f, "zero", "wait", func(d map[string]any) { d["exit_code"] = 0; d["signal"] = ""; d["error"] = "<nil>" })
		}, "protocol"},
		{"owner-abort", func(f map[string][]byte) { f["zero/first-owner/abort.json"] = []byte(`{"stage":"scheduler-return"}`) }, "protocol"},
		{"missing-journal", func(f map[string][]byte) { delete(f, "zero/endpoint-b/journal.jsonl") }, "world"},
		{"missing-snapshot", func(f map[string][]byte) { delete(f, "zero/endpoint-b/final.json") }, "world"},
		{"second-contact", func(f map[string][]byte) {
			p := "applied/endpoint-b/journal.jsonl"
			s := strings.Split(string(f[p]), "\n")[0]
			f[p] = append(f[p], []byte(s+"\n")...)
			for _, cut := range []string{"stopped", "final"} {
				editJSON(f, "applied/endpoint-b/"+cut+".json", func(d map[string]any) { d["snapshot"].(map[string]any)["Contacts"] = 2 })
			}
		}, "world"},
		{"same-applied", func(f map[string][]byte) {
			for p, b := range f {
				if strings.HasPrefix(p, "applied/endpoint-b/") {
					f["zero/"+strings.TrimPrefix(p, "applied/")] = b
				}
			}
		}, "world"},
		{"hidden-state", func(f map[string][]byte) {
			editJSON(f, "applied/final/inspect.json", func(d map[string]any) {
				d["Tables"].(map[string]any)["runs"].([]any)[0].(map[string]any)["projection"] = `{"PC":999}`
			})
		}, "equality"},
		{"nonexact-renewal", func(f map[string][]byte) {
			editJSON(f, "zero/reopened/inspect.json", func(d map[string]any) { d["Tables"].(map[string]any)["runs"].([]any)[1].(map[string]any)["epoch"] = 2 })
		}, "raw"},
		{"outer-renewal", func(f map[string][]byte) {
			editJSON(f, "zero/reopened/inspect.json", func(d map[string]any) { d["Settings"].(map[string]any)["QueryOnly"] = 0 })
		}, "raw"},
		{"stale-refusal", func(f map[string][]byte) {
			p := "zero/final/inspect.json"
			f[p] = []byte(strings.ReplaceAll(string(f[p]), `\"Epoch\":3`, `\"Epoch\":2`))
		}, "raw"},
		{"wrong-root", func(f map[string][]byte) {
			editJSON(f, "seed/manifest.json", func(d map[string]any) { d["BAServiceRoot"] = d["ABServiceRoot"] })
		}, "raw"},
		{"invented-receipt", func(f map[string][]byte) {
			editJSON(f, "zero/final/inspect.json", func(d map[string]any) {
				d["Tables"].(map[string]any)["external_receipts"] = []any{map[string]any{"intent_id": service, "descriptor": "{}"}}
			})
		}, "raw"},
		{"invented-outer-receipt", func(f map[string][]byte) {
			editJSON(f, "zero/final/inspect.json", func(d map[string]any) { d["ReceiptEvidence"] = []any{map[string]any{"Receipt": "invented"}} })
		}, "raw"},
		{"self-declared-source", func(f map[string][]byte) {
			editJSON(f, "source.json", func(d map[string]any) {
				r := d["prepared"].(map[string]any)["runtime_files"].(map[string]any)
				for k, v := range r {
					v.(map[string]any)["sha256"] = digest([]byte("altered source")).SHA256
					r[k] = v
					break
				}
			})
		}, "source"},
		{"cross-store", func(f map[string][]byte) {
			changeRecord(f, "applied", "pre-call-ack", func(d map[string]any) { d["event"].(map[string]any)["Store"] = "/synthetic/zero/runtime.sqlite" })
		}, "protocol"},
		{"wrong-application", func(f map[string][]byte) {
			changeRecord(f, "applied", "application-ack", func(d map[string]any) { d["event"].(map[string]any)["Intent"] = "wrong" })
		}, "protocol"},
		{"response-before-death", func(f map[string][]byte) {
			editProtocol(f, "applied", func(r []record) []record {
				var at, death int
				for i, v := range r {
					if v.Action == "release-request" && v.Role == "endpoint-b" {
						at = i
					}
					if v.Action == "kill-request" && v.Role == "owner" {
						death = i
					}
				}
				v := append([]record{}, r[at:at+2]...)
				r = append(r[:at], r[at+2:]...)
				return append(append(append([]record{}, r[:death]...), v...), r[death:]...)
			})
		}, "protocol"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := fixture(t)
			tc.edit(f)
			r, e := CheckDirectory(sealFixture(t, f))
			if e != nil {
				t.Fatal(e)
			}
			if !r.Seal.Valid {
				t.Fatalf("attack did not pass reseal: %+v", r.Seal)
			}
			v := map[string]Verdict{"raw": r.Raw, "source": r.Source, "protocol": r.Protocol, "world": r.World}[tc.verdict]
			if tc.verdict == "equality" {
				if r.Equality["final"].StructuralEqual {
					t.Fatal("hidden raw difference accepted")
				}
			} else if v.Valid {
				t.Fatalf("%s attack accepted: %+v", tc.verdict, r)
			}
			if r.Supported {
				t.Fatal("attack certified")
			}
		})
	}
}
func TestInvalidSeal(t *testing.T) {
	d := sealFixture(t, fixture(t))
	if e := os.WriteFile(filepath.Join(d, "unsealed.json"), []byte(`{}`), 0600); e != nil {
		t.Fatal(e)
	}
	r, e := CheckDirectory(d)
	if e != nil || r.Seal.Valid || r.Supported {
		t.Fatalf("unsealed member: %v %+v", e, r)
	}
}
func TestStatusNotSemanticAuthority(t *testing.T) {
	f := fixture(t)
	editJSON(f, "pair.json", func(d map[string]any) { d["Status"] = "failed" })
	r, e := CheckDirectory(sealFixture(t, f))
	if e != nil {
		t.Fatal(e)
	}
	if !r.World.Valid || !r.Raw.Valid || !r.Equality["final"].StructuralEqual {
		t.Fatal("status erased independent evidence")
	}
	if r.Supported {
		t.Fatal("failed attempt certified")
	}
}

func TestSeedProvisionJoin(t *testing.T) {
	f := fixture(t)
	editJSON(f, "seed/provision.json", func(d map[string]any) { d["Programs"].([]any)[0].(map[string]any)["ID"] = "wrong" })
	r, e := CheckDirectory(sealFixture(t, f))
	if e != nil || r.Raw.Valid || r.Supported {
		t.Fatalf("altered seed provision accepted: %v %+v", e, r)
	}
}
func TestUnexpectedProcessArtifact(t *testing.T) {
	f := fixture(t)
	f["zero/endpoint-c/ready.json"] = f["zero/endpoint-a/ready.json"]
	r, e := CheckDirectory(sealFixture(t, f))
	if e != nil || r.Protocol.Valid || r.Supported {
		t.Fatalf("undeclared endpoint accepted: %v %+v", e, r)
	}
}
func TestEmptySourceArgvFailsClosed(t *testing.T) {
	f := fixture(t)
	editJSON(f, "source.json", func(d map[string]any) { d["test_argv"] = []string{} })
	r, e := CheckDirectory(sealFixture(t, f))
	if e != nil || r.Source.Valid || r.Protocol.Valid || r.Supported {
		t.Fatalf("empty argv accepted: %v %+v", e, r)
	}
}

func TestDeterministicReport(t *testing.T) {
	dir := sealFixture(t, fixture(t))
	a, e := CheckDirectory(dir)
	if e != nil {
		t.Fatal(e)
	}
	b, e := CheckDirectory(dir)
	if e != nil || string(jsonBytes(a)) != string(jsonBytes(b)) {
		t.Fatalf("nondeterministic report: %v", e)
	}
}
func TestResealedDuplicateKey(t *testing.T) {
	f := fixture(t)
	f["applied/final/inspect.json"] = []byte(strings.Replace(string(f["applied/final/inspect.json"]), `"KeyVersion":2`, `"KeyVersion":2,"KeyVersion":2`, 1))
	r, e := CheckDirectory(sealFixture(t, f))
	if e != nil || r.Seal.Valid || r.Raw.Valid || r.Equality["final"].StructuralEqual || r.Supported {
		t.Fatalf("duplicate key accepted: %v", e)
	}
}

func TestCoherentIdentityAttacks(t *testing.T) {
	for _, tc := range []struct {
		name, branch string
		edit         func(map[string][]byte)
	}{
		{"precall-binding-all-copies", "zero", func(f map[string][]byte) {
			for p, b := range f {
				if strings.HasPrefix(p, "zero/first-owner/") || p == "zero/protocol.jsonl" {
					f[p] = []byte(strings.ReplaceAll(string(b), `\"authorization\",\"[\\\"a\\\",1,4,0]\"`, `\"authorization\",\"wrong\"`))
				}
			}
		}},
		{"shared-store-all-copies", "applied", func(f map[string][]byte) {
			for p, b := range f {
				if strings.HasPrefix(p, "applied/") {
					f[p] = []byte(strings.ReplaceAll(string(b), "/synthetic/applied/runtime.sqlite", "/synthetic/zero/runtime.sqlite"))
				}
			}
		}},
		{"endpoint-reused-all-copies", "applied", func(f map[string][]byte) {
			for p, b := range f {
				if strings.HasPrefix(p, "applied/") {
					f[p] = []byte(strings.ReplaceAll(string(b), "applied-endpoint-b", "zero-endpoint-b"))
				}
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := fixture(t)
			before := string(f[tc.branch+"/protocol.jsonl"])
			tc.edit(f)
			if before == string(f[tc.branch+"/protocol.jsonl"]) {
				t.Fatal("attack did not change transcript")
			}
			r, e := CheckDirectory(sealFixture(t, f))
			if e != nil || !r.Seal.Valid || r.Protocol.Valid || r.Supported {
				t.Fatalf("coherent identity attack accepted: %v %+v", e, r.Protocol)
			}
		})
	}
}
func TestSourceNullDependencyFailsClosed(t *testing.T) {
	f := fixture(t)
	editJSON(f, "source.json", func(d map[string]any) { d["build_info"].(map[string]any)["Deps"] = []any{nil} })
	r, e := CheckDirectory(sealFixture(t, f))
	if e != nil || r.Source.Valid || r.Supported {
		t.Fatalf("null build dependency accepted: %v", e)
	}
}

func TestSyntheticCommand(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "check")
	build := exec.Command("go", "build", "-o", bin, "./cmd/check")
	if b, e := build.CombinedOutput(); e != nil {
		t.Fatalf("build: %v %s", e, b)
	}
	dir := sealFixture(t, fixture(t))
	var previous []byte
	for i := 0; i < 2; i++ {
		command := exec.Command(bin, "-pair", dir)
		var diagnostics bytes.Buffer
		command.Stderr = &diagnostics
		b, e := command.Output()
		if e != nil || diagnostics.Len() != 0 {
			t.Fatalf("command: %v %s", e, diagnostics.Bytes())
		}
		var report Report
		if e = json.Unmarshal(b, &report); e != nil || !report.Supported {
			t.Fatalf("command report: %v", e)
		}
		if i > 0 && !bytes.Equal(previous, b) {
			t.Fatal("command output changed")
		}
		previous = b
	}
}

func TestMalformedEvidenceFailsClosed(t *testing.T) {
	for _, path := range []string{"pair.json", "source.json", "seed/manifest.json", "seed/provision.json", "zero/stopped/inspect.json", "zero/first-owner/open.json", "zero/endpoint-b/ready.json", "zero/endpoint-b/final.json", "zero/recovery-open.json"} {
		for _, data := range []string{"null", "{}", "[]"} {
			t.Run(path+"/"+data, func(t *testing.T) {
				f := fixture(t)
				f[path] = []byte(data)
				r, e := CheckDirectory(sealFixture(t, f))
				if e != nil || r.Supported {
					t.Fatalf("malformed evidence accepted: %v", e)
				}
			})
		}
	}
}

func TestProjectedMatchDoesNotHideFullDifference(t *testing.T) {
	f := fixture(t)
	editJSON(f, "applied/final/inspect.json", func(d map[string]any) {
		d["Settings"].(map[string]any)["SQLiteVersion"] = "synthetic-different-version"
	})
	r, e := CheckDirectory(sealFixture(t, f))
	if e != nil || !r.Raw.Valid || !r.Model.Valid || r.Equality["final"].StructuralEqual || r.Supported {
		t.Fatalf("projection/full comparison separation: %v raw=%+v model=%+v equality=%+v", e, r.Raw, r.Model, r.Equality["final"])
	}
	if len(r.Equality["final"].Paths) != 1 || r.Equality["final"].Paths[0] != "$.Settings.SQLiteVersion" {
		t.Fatal(r.Equality["final"].Paths)
	}
}
func TestMalformedOuterSchema(t *testing.T) {
	for _, field := range []string{"Settings", "Roots"} {
		t.Run(field, func(t *testing.T) {
			f := fixture(t)
			for _, branch := range []string{"zero", "applied"} {
				editJSON(f, branch+"/final/inspect.json", func(d map[string]any) { d[field] = nil })
			}
			r, e := CheckDirectory(sealFixture(t, f))
			if e != nil || r.Raw.Valid || r.Supported {
				t.Fatalf("equal malformed outer schema accepted: %v raw=%+v", e, r.Raw)
			}
		})
	}
}

func TestSeedSourceMultiplicity(t *testing.T) {
	for _, field := range []string{"Programs", "Runs", "Subjects", "Permissions"} {
		t.Run(field, func(t *testing.T) {
			f := fixture(t)
			editJSON(f, "seed/provision.json", func(d map[string]any) { a := d[field].([]any); a[1] = a[0] })
			r, e := CheckDirectory(sealFixture(t, f))
			if e != nil || r.Raw.Valid || r.Supported {
				t.Fatalf("duplicate seed %s accepted: %v", field, e)
			}
		})
	}
}
func TestWrongFixtureManifest(t *testing.T) {
	f := fixture(t)
	editJSON(f, "seed/manifest.json", func(d map[string]any) { d["Options"].(map[string]any)["Scenario"] = "other" })
	r, e := CheckDirectory(sealFixture(t, f))
	if e != nil || r.Raw.Valid || r.Supported {
		t.Fatalf("wrong fixture declaration accepted: %v", e)
	}
}

func TestSyntheticRacePair(t *testing.T) {
	f := fixture(t)
	editJSON(f, "pair.json", func(d map[string]any) { d["Mode"] = "race" })
	editJSON(f, "source.json", func(d map[string]any) {
		d["race"] = true
		info := d["build_info"].(map[string]any)
		info["Settings"] = append(info["Settings"].([]any), map[string]any{"Key": "-race", "Value": "true"})
	})
	r, e := CheckDirectory(sealFixture(t, f))
	if e != nil || !r.Supported {
		t.Fatalf("synthetic race metadata: %v source=%+v", e, r.Source)
	}
}
func TestWrongTestSelection(t *testing.T) {
	f := fixture(t)
	editJSON(f, "source.json", func(d map[string]any) {
		d["test_argv"] = []string{"/synthetic/owner.test", "-test.run=OtherTestPhase3B7MatchedCapture"}
	})
	r, e := CheckDirectory(sealFixture(t, f))
	if e != nil || r.Source.Valid || r.Supported {
		t.Fatalf("unfrozen test selector accepted: %v", e)
	}
}

func TestManifestCannotReplaceMemberWithMissingZeroDigest(t *testing.T) {
	dir := sealFixture(t, fixture(t))
	p := filepath.Join(dir, "members.json")
	b, e := os.ReadFile(p)
	if e != nil {
		t.Fatal(e)
	}
	var m map[string]any
	if e = json.Unmarshal(b, &m); e != nil {
		t.Fatal(e)
	}
	files := m["files"].(map[string]any)
	delete(files, "zero/first-owner/process.log")
	files["absent.json"] = map[string]any{"bytes": 0, "sha256": ""}
	if e = os.WriteFile(p, jsonBytes(m), 0600); e != nil {
		t.Fatal(e)
	}
	r, e := CheckDirectory(dir)
	if e != nil || r.Seal.Valid || r.Supported {
		t.Fatalf("manifest silently lost actual member: %v seal=%+v", e, r.Seal)
	}
}

func TestMissingExportsDoNotClaimByteEquality(t *testing.T) {
	r, e := CheckDirectory(t.TempDir())
	if e != nil || r.Equality["final"].ByteEqual || r.ReportEquality.ByteEqual {
		t.Fatalf("missing exports treated as equal byte artifacts: %v", e)
	}
}
