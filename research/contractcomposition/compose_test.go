package contractcomposition

import (
	"reflect"
	"testing"
)

func TestClosedPortsOwnersResourcesRejected(t *testing.T) {
	c := fixture("chain", 7)
	a, _ := Leaf(c, "A")
	b, _ := Leaf(c, "B")
	if _, err := Close(a); err == nil {
		t.Error("closed unresolved B port")
	}
	if _, err := Join(a, a); err == nil {
		t.Error("accepted duplicated owner")
	}
	b.Originals = []Binding{c.Originals[1]}
	b.Originals[0].Resource = c.Originals[0].Resource
	if _, err := Join(a, b); err == nil {
		t.Error("accepted shared resource")
	}
}
func TestGroupedEqualityAllMasks(t *testing.T) {
	for _, topology := range []string{"chain", "cycle"} {
		for mask := uint8(0); mask < 8; mask++ {
			c := fixture(topology, mask)
			a, e := Leaf(c, "A")
			if e != nil {
				t.Fatal(e)
			}
			b, _ := Leaf(c, "B")
			d, _ := Leaf(c, "C")
			ab, e := Join(a, b)
			if e != nil {
				t.Fatal(e)
			}
			bc, e := Join(b, d)
			if e != nil {
				t.Fatal(e)
			}
			l, e := Join(ab, d)
			if e != nil {
				t.Fatal(e)
			}
			r, e := Join(a, bc)
			if e != nil {
				t.Fatal(e)
			}
			left, e := Close(l)
			if e != nil {
				t.Fatal(e)
			}
			right, e := Close(r)
			if e != nil {
				t.Fatal(e)
			}
			if !reflect.DeepEqual(c, left) || !reflect.DeepEqual(left, right) {
				t.Fatalf("lost named contract at %s/%d", topology, mask)
			}
			if len(ab.Originals) != 2 || len(ab.Ports) != 1 || ab.Ports[0].ID != "C" {
				t.Fatal("open AB must retain internal AB and external BC")
			}
		}
	}
}
