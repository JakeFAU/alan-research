package componentinterfaces

import (
	pc "github.com/JakeFAU/alan-research/research/capacitycomposition"
	"testing"
)

func coreConfig(t *testing.T) pc.Config {
	t.Helper()
	c, e := pc.NewConfig("shared-1", "early", 0)
	if e != nil {
		t.Fatal(e)
	}
	return c
}

// Rejecting incomplete role loops prevents a monitor from vetoing another owner.
func TestValidation(t *testing.T) {
	c := coreConfig(t)
	for _, name := range []string{"good", "forged-binding", "wrong-root", "missing-external-loop", "hidden-effect", "bad-state", "duplicate-edge"} {
		t.Run(name, func(t *testing.T) {
			a := monitors(c)[0]
			switch name {
			case "forged-binding":
				a.Signature.Base.Originals[0].Grant = "forged"
			case "wrong-root":
				a.Signature.Roots[0].Limit = 3
			case "missing-external-loop":
				a.Edges = a.Edges[1:]
			case "hidden-effect":
				a.Edges = append(a.Edges, edge{"0", "hidden-apply", "0"})
			case "bad-state":
				a.Edges = append(a.Edges, edge{"0", "tau", "missing"})
			case "duplicate-edge":
				a.Edges = append(a.Edges, a.Edges[0])
			}
			err := validate(a)
			if (err == nil) != (name == "good") {
				t.Fatalf("validation=%v", err)
			}
		})
	}
}

func TestSignatureAgreement(t *testing.T) {
	c := coreConfig(t)
	other, err := pc.NewConfig("shared-2", "early", 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := product(plant(c), monitors(other)[0]); err == nil {
		t.Fatal("accepted incompatible pool signature")
	}
	for _, name := range []string{"extra-root", "missing-original", "extra-original", "bad-initial", "duplicate-state", "unknown-role"} {
		t.Run(name, func(t *testing.T) {
			a := monitors(c)[0]
			switch name {
			case "extra-root":
				a.Signature.Roots = append(a.Signature.Roots, pc.Root{ID: "extra", Limit: 1})
			case "missing-original":
				a.Signature.Base.Originals = a.Signature.Base.Originals[:1]
			case "extra-original":
				a.Signature.Base.Originals = append(a.Signature.Base.Originals, a.Signature.Base.Originals[0])
			case "bad-initial":
				a.Initial = "missing"
			case "duplicate-state":
				a.States = append(a.States, a.States[0])
			case "unknown-role":
				a.Role = "rogue"
			}
			if err := validate(a); err == nil {
				t.Fatal("accepted malformed descriptor")
			}
		})
	}
}
