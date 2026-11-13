package effectobservation

import "testing"

func TestAuditControls(t *testing.T) {
	controls, err := AuditControls()
	if err != nil {
		t.Fatal(err)
	}
	accepted, rejected := 0, 0
	for _, c := range controls {
		t.Run(c.Name, func(t *testing.T) {
			if !c.Passed {
				t.Fatal("control failed")
			}
			if c.Rejection != "" {
				rejected++
			} else {
				accepted++
			}
		})
	}
	if accepted != 9 || rejected != 11 {
		t.Fatalf("expected 9 lawful and 11 hostile controls, got %d/%d", accepted, rejected)
	}
}
