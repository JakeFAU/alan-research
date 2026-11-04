package delegationcorrespondence

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExternalEvidenceAndAppendOnly(t *testing.T) {
	root := os.Getenv("ALAN_TASK3_RAW")
	if root == "" {
		t.Skip("frozen corpus")
	}
	b, e := os.ReadFile(filepath.Join(root, "case-index.json"))
	if e != nil {
		t.Fatal(e)
	}
	var idx []IndexEntry
	strictJSON(b, &idx)
	for _, c := range idx {
		if c.Name != "ordinary-s2-handshake" {
			continue
		}
		raw, e := os.ReadFile(filepath.Join(c.Directory, "final/inspect.json"))
		if e != nil {
			t.Fatal(e)
		}
		files := map[string][]byte{}
		for _, n := range []string{"endpoint-a.jsonl", "endpoint-b.jsonl", "boundaries.json", "endpoints.json"} {
			files[n], e = os.ReadFile(filepath.Join(c.Directory, n))
			if e != nil {
				t.Fatal(e)
			}
		}
		v := CheckExternal(raw, files)
		if !v.Valid || v.Contacts != 2 || v.Applications != 2 {
			t.Fatalf("%+v", v)
		}
		files["endpoint-b.jsonl"] = []byte(strings.Replace(string(files["endpoint-b.jsonl"]), `\"Delta\":1`, `\"Delta\":9`, 1))
		v = CheckExternal(raw, files)
		if v.Valid {
			t.Fatal("altered endpoint wire accepted")
		}
	}
}
