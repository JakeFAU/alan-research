package componentinterfaces

import (
	"reflect"
	"testing"
)

func TestConcreteHistories(t *testing.T) {
	hs, e := histories()
	if e != nil {
		t.Fatal(e)
	}
	if len(hs) != 8 {
		t.Fatal(len(hs))
	}
	if !reflect.DeepEqual(hs[0].Config, hs[1].Config) || hs[0].Steps[hs[0].Cuts[0].AfterSteps-1].State.Old != hs[1].Steps[hs[1].Cuts[0].AfterSteps-1].State.Old {
		t.Fatal("candidate cuts differ")
	}
	if hs[0].Cuts[0].CanComplete || !hs[1].Cuts[0].CanComplete {
		t.Fatal("cut classification")
	}
	for i, h := range hs {
		if err := checkHistory(h); err != nil {
			t.Fatal(err)
		}
		want := i != 0 && i != 7
		if h.Complete != want || h.CanComplete != want {
			t.Fatalf("history %d completion", i)
		}
	}
	// Private control cannot advance while physical B is down; crash/recover retains PC.
	h := hs[5]
	m := monitors(h.Config)[3]
	for i, s := range h.Steps {
		if s.Event.Kind == "crash" {
			if s.State.Monitor != h.Steps[i-1].State.Monitor {
				t.Fatal("crash reset")
			}
			for _, e := range concreteSuccessors(h.Config, m, s.State) {
				if e.Private {
					t.Fatal("private step while B down")
				}
			}
		}
		if s.Event.Kind == "recover" && s.State.Monitor != h.Steps[i-1].State.Monitor {
			t.Fatal("recovery reset")
		}
	}
}

func TestConcreteCaptureTampering(t *testing.T) {
	hs, err := histories()
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"reset-PC", "invented-result", "fake-private", "changed-cut", "changed-hold"} {
		t.Run(name, func(t *testing.T) {
			h := hs[5]
			h.Steps = append([]microStep(nil), h.Steps...)
			h.Cuts = append([]cut(nil), h.Cuts...)
			switch name {
			case "reset-PC":
				for i, e := range h.Steps {
					if e.Event.Kind == "recover" {
						h.Steps[i].State.Monitor = "0"
						break
					}
				}
			case "invented-result":
				h.Steps[0].State.Old.Base.Originals[0].WorldResult = "r1"
			case "fake-private":
				h.Steps[0].Private = true
			case "changed-cut":
				h.Cuts[0].CanComplete = !h.Cuts[0].CanComplete
			case "changed-hold":
				h.Steps[0].State.Old.Held[0] = true
			}
			if err := checkHistory(h); err == nil {
				t.Fatal("accepted altered capture")
			}
		})
	}
}
