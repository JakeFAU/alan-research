package capacitycomposition

import (
	"reflect"
	"testing"
)

func TestJoinRetainsOneRootAndAssociates(t *testing.T) {
	for _, layout := range []string{"shared-1", "shared-2", "separate-1"} {
		for _, policy := range []string{"early", "ready-now", "stable-ready"} {
			for mask := uint8(0); mask < 8; mask++ {
				c := config(t, layout, policy, mask)
				a, e := Leaf(c, "A")
				if e != nil {
					t.Fatal(e)
				}
				b, e := Leaf(c, "B")
				if e != nil {
					t.Fatal(e)
				}
				d, e := Leaf(c, "C")
				if e != nil {
					t.Fatal(e)
				}
				ab, e := Join(a, b)
				if e != nil {
					t.Fatal(e)
				}
				bc, e := Join(b, d)
				if e != nil {
					t.Fatal(e)
				}
				left, e := Join(ab, d)
				if e != nil {
					t.Fatal(e)
				}
				right, e := Join(a, bc)
				if e != nil {
					t.Fatal(e)
				}
				l, e := Close(left)
				if e != nil {
					t.Fatal(e)
				}
				r, e := Close(right)
				if e != nil {
					t.Fatal(e)
				}
				if !reflect.DeepEqual(left, right) || !reflect.DeepEqual(c, l) || !reflect.DeepEqual(c, r) || !reflect.DeepEqual(Explore(l), Explore(r)) {
					t.Fatal("regrouping changed roots or relation")
				}
			}
		}
	}
}
func TestInvalidCompositionRejected(t *testing.T) {
	for _, mutation := range []string{"duplicate-owner", "missing-member", "extra-member", "changed-root", "copied-pool", "policy", "limit", "duplicate-original", "missing-original"} {
		c := config(t, "shared-1", "early", 0)
		a, _ := Leaf(c, "A")
		b, _ := Leaf(c, "B")
		if len(b.Roots) != 1 || len(b.Memberships) != 1 {
			t.Fatal("leaf omitted root membership")
		}
		switch mutation {
		case "duplicate-owner":
			b = a
		case "missing-member":
			b.Memberships = nil
		case "extra-member":
			b.Memberships = append(b.Memberships, Membership{"AB", "Q"})
		case "changed-root":
			b.Memberships[0].Root = "invented"
		case "copied-pool":
			b.Roots = append(b.Roots, b.Roots[0])
		case "policy":
			b.Policy = "ready-now"
		case "limit":
			b.Roots[0].Limit = 2
		case "duplicate-original":
			b.Base.Originals = append(b.Base.Originals, b.Base.Originals[0])
		case "missing-original":
			b.Base.Originals = nil
			b.Base.Ports = nil
			b.Memberships = nil
			b.Roots = nil
		}
		if _, err := Join(a, b); err == nil {
			t.Errorf("accepted %s", mutation)
		}
	}
}
