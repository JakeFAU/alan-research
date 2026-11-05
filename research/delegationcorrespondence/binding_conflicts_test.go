package delegationcorrespondence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBindingConflictPromotions(t *testing.T) {
	root := os.Getenv("ALAN_PHASE3B4_PROMOTED_CHECK_ROOT")
	if root == "" {
		t.Skip("ALAN_PHASE3B4_PROMOTED_CHECK_ROOT not set: ancillary writer copies")
	}
	paths, err := filepath.Glob(filepath.Join(root, "*", "promoted-inspect.json"))
	if err != nil || len(paths) != 2 {
		t.Fatalf("expected two promoted parents: %d %v", len(paths), err)
	}
	for _, p := range paths {
		t.Run(filepath.Base(filepath.Dir(p)), func(t *testing.T) {
			before, err := os.ReadFile(filepath.Join(filepath.Dir(p), "parent-inspect.json"))
			if err != nil {
				t.Fatal(err)
			}
			after, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			opened, err := os.ReadFile(filepath.Join(filepath.Dir(p), "opened-inspect.json"))
			if err != nil {
				t.Fatal(err)
			}
			v := CheckConflictPromotion(before, opened, after)
			if !v.Valid || v.Conflicts != 1 {
				t.Fatalf("promotion: %+v", v.Failure)
			}
			if out := os.Getenv("ALAN_PHASE3B4_CHECK_OUTPUT"); out != "" {
				if err = os.MkdirAll(out, 0755); err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(filepath.Join(out, filepath.Base(filepath.Dir(p))+"-promotion.json"), bytesJSON(v), 0644); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
