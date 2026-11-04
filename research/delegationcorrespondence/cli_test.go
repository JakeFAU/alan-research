package delegationcorrespondence

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCLIExits(t *testing.T) {
	var out, err bytes.Buffer
	if n := RunCheck([]string{"-bad"}, &out, &err); n != 2 {
		t.Fatal(n)
	}
	p := filepath.Join(t.TempDir(), "raw.json")
	os.WriteFile(p, positive(t), 0600)
	if n := RunCheck([]string{"-input", p}, &out, &err); n != 0 {
		t.Fatalf("%d %s", n, err.String())
	}
	os.WriteFile(p, []byte(`{"Tables":{}}`), 0600)
	if n := RunCheck([]string{"-input", p}, &out, &err); n != 1 {
		t.Fatal(n)
	}
	os.WriteFile(p, []byte(`{`), 0600)
	if n := RunCheck([]string{"-input", p}, &out, &err); n != 2 {
		t.Fatal(n)
	}
}
func TestPortableCorpus(t *testing.T) {
	root := os.Getenv("ALAN_TASK3_RAW")
	if root == "" {
		t.Skip("frozen corpus")
	}
	dest := filepath.Join(t.TempDir(), "relocated")
	if e := os.CopyFS(dest, os.DirFS(root)); e != nil {
		t.Fatal(e)
	}
	r, e := CheckCorpus(dest, root)
	if e != nil || len(r.Cases) != 48 {
		t.Fatalf("%v %d", e, len(r.Cases))
	}
	for _, c := range r.Cases {
		if c.CaseSHA256 == "" {
			t.Fatal("missing digest")
		}
	}
	if _, e := resolve(dest, root, filepath.Join(root, "..", "other")); e == nil {
		t.Fatal("traversal accepted")
	}
}
