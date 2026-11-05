package delegationcorrespondence

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func bindingInput(t *testing.T, name, file string) []byte {
	t.Helper()
	if root := os.Getenv("ALAN_PHASE3B4_CHECK_ROOT"); root != "" {
		wanted := name
		if name == "M03" {
			wanted = "M03-auxiliary-target-generation-two"
		}
		paths, err := filepath.Glob(filepath.Join(root, "*", "case.json"))
		if err != nil {
			t.Fatal(err)
		}
		found := ""
		for _, p := range paths {
			b, e := os.ReadFile(p)
			if e != nil {
				t.Fatal(e)
			}
			var c BindingCase
			if e = json.Unmarshal(b, &c); e != nil {
				t.Fatal(e)
			}
			if c.Name == wanted {
				if found != "" {
					t.Fatalf("ambiguous final parent %s", wanted)
				}
				found = filepath.Dir(p)
			}
		}
		if found == "" {
			t.Fatalf("missing genuine final parent %s in %s", wanted, root)
		}
		b, err := os.ReadFile(filepath.Join(found, file))
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	filename := name + "-" + strings.ReplaceAll(file, "/", "-") + ".gz"
	f, err := os.Open("testdata/binding/" + filename)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	z, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer z.Close()
	b, err := io.ReadAll(z)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func bindingFixture(t *testing.T, name string) []byte {
	t.Helper()
	return bindingInput(t, name, "final/inspect.json")
}

func TestBindingNonadjacent(t *testing.T) {
	v := CheckBindingRaw(bindingFixture(t, "P02-v2"), bindingSpec(t, "P02-v2"))
	if !v.Valid {
		t.Fatalf("genuine nonadjacent export rejected: %+v", v.Failure)
	}
	for _, e := range v.Effects {
		if e.Source.Subject == "a" && (e.Witness.RequestConsumedSequence != 2 || !e.Allowed) {
			t.Fatalf("incorrect original selection: %+v", e)
		}
	}
}

func TestBindingExplicitPrincipals(t *testing.T) {
	for _, name := range []string{"P10", "auxiliary-serviced-equal-total"} {
		t.Run(name, func(t *testing.T) {
			v := CheckBindingRaw(bindingFixture(t, name), bindingSpec(t, name))
			if !v.Valid {
				t.Fatalf("genuine distinct auxiliary rejected: %+v", v.Failure)
			}
			if !v.Custody.Complete || len(v.Custody.Originals) != 2 || len(v.Custody.Auxiliary) != 1 || v.Custody.Auxiliary[0].ReplySelected {
				t.Fatalf("incorrect principal/auxiliary accounting: %+v", v.Custody)
			}
		})
	}
}

func bindingSpec(t *testing.T, name string) BindingSpec {
	t.Helper()
	var c BindingCase
	if err := json.Unmarshal(bindingInput(t, name, "case.json"), &c); err != nil {
		t.Fatal(err)
	}
	return c.Spec()
}
