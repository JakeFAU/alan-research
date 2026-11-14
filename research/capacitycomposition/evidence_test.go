package capacitycomposition

import (
	"reflect"
	"slices"
	"testing"
)

func TestExportWitnessesAreActualHistories(t *testing.T) {
	ws, err := BuildWitnesses()
	if err != nil {
		t.Fatal(err)
	}
	if len(ws) < 10 {
		t.Fatal("missing required witness histories")
	}
	byName := map[string]Witness{}
	for _, w := range ws {
		if err := CheckWitness(w); err != nil {
			t.Fatalf("%s: %v", w.Name, err)
		}
		byName[w.Name] = w
	}
	for _, policy := range []string{"early", "ready-now"} {
		a := byName[policy+"-AB-holder"]
		b := byName[policy+"-BC-holder"]
		if len(a.Steps) == 0 || len(b.Steps) == 0 {
			t.Fatal("missing holder pair")
		}
		x := a.Steps[len(a.Steps)-1].State
		y := b.Steps[len(b.Steps)-1].State
		if x.Base != y.Base || !reflect.DeepEqual(Occupancy(a.Config, Project(a.Config, x)), Occupancy(b.Config, Project(b.Config, y))) || a.CanComplete || !b.CanComplete {
			t.Fatal("genuine holder collision not established")
		}
	}
	for _, mutation := range []string{"initial-receipt", "holder", "early-release", "actor", "activation", "grant", "receipt", "post-state", "duplicate-reserve"} {
		w := byName["early-AB-holder"]
		w.Steps = slices.Clone(w.Steps)
		i := len(w.Steps) - 1
		switch mutation {
		case "initial-receipt":
			w.Initial.Base.Originals[0].Retained = "r1"
		case "holder":
			w.Steps[i].State.Held = [2]bool{false, true}
		case "early-release":
			w.Steps[i].State.Held = [2]bool{}
		case "actor":
			w.Steps[i].Event.Actor = "A"
		case "activation":
			w.Steps[i].State.Base.Agents[1].Activation = 1
		case "grant":
			w.Steps[i].Event.Binding.Grant = "forged"
		case "receipt":
			w.Steps[i].State.Base.Originals[0].Retained = "r2"
		case "post-state":
			w.Steps[i].State.Base.Originals[1].Consumed = true
		case "duplicate-reserve":
			w.Steps = append(w.Steps, w.Steps[i])
		}
		if CheckWitness(w) == nil {
			t.Errorf("accepted hostile %s", mutation)
		}
	}
}
func TestFactoredRepresentativeProjectionAndEveryEnabledLift(t *testing.T) {
	summary, err := ProbeRepresentatives()
	if err != nil {
		t.Fatal(err)
	}
	if summary.Representatives == 0 || summary.Lifts == 0 || summary.MicroArcs == 0 {
		t.Fatal("empty probe domain")
	}
	t.Logf("probe summary: %+v", summary)
}
