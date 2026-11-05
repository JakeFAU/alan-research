package delegationcorrespondence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBindingExportCorpus(t *testing.T) {
	root := os.Getenv("ALAN_PHASE3B4_CHECK_ROOT")
	if root == "" {
		t.Skip("ALAN_PHASE3B4_CHECK_ROOT not set: genuine binding export campaign")
	}
	paths, err := filepath.Glob(filepath.Join(root, "*", "case.json"))
	if err != nil || len(paths) == 0 {
		t.Fatalf("case discovery: %v (%d)", err, len(paths))
	}
	for _, p := range paths {
		t.Run(filepath.Base(filepath.Dir(p)), func(t *testing.T) {
			b, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			var c BindingCase
			if err = json.Unmarshal(b, &c); err != nil {
				t.Fatal(err)
			}
			report, err := CheckBindingDirectory(filepath.Dir(p))
			if err != nil {
				t.Fatal(err)
			}
			for _, phase := range report.Phases {
				if !phase.Verdict.Valid {
					t.Errorf("%s: %+v", phase.Path, phase.Verdict.Failure)
				}
			}
			if len(c.Captures) > 0 && !report.External.Valid {
				t.Errorf("external: %+v", report.External.Failure)
			}
			if report.ContinuityFailure != nil {
				t.Errorf("continuity: %+v", report.ContinuityFailure)
			}
			switch c.Name {
			case "P01-v2", "P02-v2", "P04", "P05", "P06", "P07", "P08", "P09", "P11":
				if !report.Projection.Supported || report.Projection.Error != "" {
					t.Errorf("required supported projection: %+v", report.Projection)
				}
			}

			if out := os.Getenv("ALAN_PHASE3B4_CHECK_OUTPUT"); out != "" {
				if err = os.MkdirAll(out, 0755); err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(filepath.Join(out, c.Name+".json"), bytesJSON(report), 0644); err != nil {
					t.Fatal(err)
				}
			}

		})
	}
}

func TestBindingOrderedRenewals(t *testing.T) {
	root := os.Getenv("ALAN_PHASE3B4_CHECK_ROOT")
	if root == "" {
		t.Skip("ALAN_PHASE3B4_CHECK_ROOT not set: genuine two-loss C07")
	}
	caseBytes := bindingInput(t, "C07", "case.json")
	var c BindingCase
	if err := json.Unmarshal(caseBytes, &c); err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{"case.json": caseBytes}
	paths, err := filepath.Glob(filepath.Join(root, "C07-*"))
	if err != nil || len(paths) != 1 {
		t.Fatalf("C07 paths: %v", paths)
	}
	err = filepath.WalkDir(paths[0], func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(p, ".json") || strings.HasSuffix(p, ".jsonl") {
			rel, e := filepath.Rel(paths[0], p)
			if e != nil {
				return e
			}
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			files[rel] = b
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// A symlinked preliminary selection is followed explicitly by ReadDir.
	if len(files) == 1 {
		target, err := filepath.EvalSymlinks(paths[0])
		if err != nil {
			t.Fatal(err)
		}
		err = filepath.WalkDir(target, func(p string, d os.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if d.IsDir() {
				return nil
			}
			if strings.HasSuffix(p, ".json") || strings.HasSuffix(p, ".jsonl") {
				rel, e := filepath.Rel(target, p)
				if e != nil {
					return e
				}
				b, e := os.ReadFile(p)
				if e != nil {
					return e
				}
				files[rel] = b
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	raw := files["final/inspect.json"]
	v := CheckExternal(raw, files)
	if !v.Valid || v.AcknowledgedLosses != 2 || len(v.Owners) != 3 {
		t.Fatalf("C07 renewal: valid=%v losses=%d owners=%d failure=%+v", v.Valid, v.AcknowledgedLosses, len(v.Owners), v.Failure)
	}
	for name, o := range v.Owners {
		if !o.RenewalChecked {
			t.Fatalf("unchecked genuine renewal %s", name)
		}
	}
	delete(files, "initial/inspect.json")
	v = CheckExternal(raw, files)
	if !v.Valid || v.Owners["owner-01"].RenewalChecked {
		t.Fatalf("missing pre-Open was not explicit unchecked: %+v", v.Failure)
	}
	loss := files["loss-02.json"]
	files["loss-02.json"] = []byte(`{"PID":1,"Wait":"signal: killed"}`)
	v = CheckExternal(raw, files)
	if v.Valid || v.Failure == nil || v.Failure.Edge != "missing_acknowledged_loss" {
		t.Fatalf("unacknowledged second loss accepted: %+v", v.Failure)
	}
	files["loss-02.json"] = loss
}

func TestBindingPrefixMutationDiagnostics(t *testing.T) {
	if os.Getenv("ALAN_PHASE3B4_CHECK_ROOT") == "" {
		t.Skip("ALAN_PHASE3B4_CHECK_ROOT not set: genuine incomplete captures")
	}
	for _, name := range []string{"C01", "C06", "C07", "N03", "M03"} {
		t.Run(name, func(t *testing.T) {
			var c BindingCase
			if err := json.Unmarshal(bindingInput(t, name, "case.json"), &c); err != nil {
				t.Fatal(err)
			}
			for _, capture := range c.Captures {
				raw := bindingInput(t, name, capture.Directory+"/inspect.json")
				ms := RawMutations(raw)
				if len(ms) != 15 {
					t.Fatalf("%s lost diagnostics: %d", capture.Directory, len(ms))
				}
				for _, m := range ms {
					if m.Applicable {
						if m.OriginalSHA256 == m.MutatedSHA256 || len(m.Input) == 0 {
							t.Fatalf("%s/%s no mutation", capture.Directory, m.Name)
						}
					} else if m.Diagnostic == "" || m.Verdict.Valid || m.Verdict.Failure != nil {
						t.Fatalf("%s/%s inapplicability misreported", capture.Directory, m.Name)
					}
				}
			}
		})
	}
}
