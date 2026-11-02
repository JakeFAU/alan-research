//go:build integration

package executor

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/JakeFAU/alan/effects"
)

// Challenge the recorder directly. The public parent has no corresponding
// receipt injection method. Keep both candidates even though only one survives.
func TestPhase2ARecordConflict(t *testing.T) {
	s, path, id := committedAdapter(t)
	work, err := s.prepareAdapterWork(t.Context(), id, map[string]effects.Adapter{"model": callbackAdapter{}}, time.Second)
	if err != nil || work == nil {
		t.Fatal(work, err)
	}
	first := transportAdapter(t.Context(), *work).candidate
	retained, err := s.recordAdapterResult(t.Context(), work.attempt, first)
	if err != nil {
		t.Fatal(err)
	}
	for _, control := range []string{"foreign-key", "contradictory-result"} {
		t.Run(control, func(t *testing.T) {
			changed := effects.Copy(first)
			if control == "foreign-key" {
				changed.Key += "-other"
			} else {
				changed.Payload["answer"] = Text("incompatible")
			}
			_, err := s.recordAdapterResult(t.Context(), work.attempt, changed)
			if !errors.Is(err, ErrConflict) {
				t.Fatalf("unexpected recorder disposition %v", err)
			}
			inspection, inspectErr := Inspect(t.Context(), path, DefaultLimits())
			if inspectErr != nil {
				t.Fatal(inspectErr)
			}
			rows := inspection.Tables["adapter_results"]
			if len(rows) != 1 {
				t.Fatal("lost or duplicated original result")
			}
			var actual AdapterResult
			if e := json.Unmarshal([]byte(rows[0]["descriptor"].(string)), &actual); e != nil || !reflect.DeepEqual(actual, retained) {
				t.Fatal("changed original", e)
			}
			raw, e := json.Marshal(map[string]any{"original": first, "offered": changed, "error": err.Error(), "retained": actual, "attempt": work.attempt})
			if e != nil {
				t.Fatal(e)
			}
			t.Logf("RECORDER %s", raw)
		})
	}
}

func TestBaselineRecorderOfferOrder(t *testing.T) {
	for _, order := range [][2]int64{{5, 10}, {10, 5}} {
		name := fmt.Sprintf("%d_then_%d", order[0], order[1])
		t.Run("raw/"+name, func(t *testing.T) {
			s, path := messageStore(t, externalProvision(t, "http://127.0.0.1:1", "idempotent"))
			run, err := s.Commit(t.Context(), current(t, s, "client").Authority(), "")
			if err != nil {
				t.Fatal(err)
			}
			a, _, err := s.prepareExternal(t.Context(), run.Projection.Data["pending"].Text)
			if err != nil {
				t.Fatal(err)
			}
			checkConflictingOffers(t, path, order, func(total int64) error {
				receipt := receiptFor(a.Request)
				receipt.Total = total
				return s.recordExternalReceipt(t.Context(), a, receipt)
			})
			got, err := externalOutcome(t.Context(), s.conn, a.Intent.ID)
			if err != nil || got.Receipt == nil || got.Receipt.Total != order[0] {
				t.Fatalf("retained raw result = %+v, %v", got, err)
			}
		})
		t.Run("adapter/"+name, func(t *testing.T) {
			s, path, id := committedAdapter(t)
			work, err := s.prepareAdapterWork(t.Context(), id, map[string]effects.Adapter{"model": callbackAdapter{}}, time.Second)
			if err != nil || work == nil {
				t.Fatal(work, err)
			}
			base := transportAdapter(t.Context(), *work).candidate
			checkConflictingOffers(t, path, order, func(value int64) error {
				c := effects.Copy(base)
				c.Payload["answer"] = Text(fmt.Sprint(value))
				_, err := s.recordAdapterResult(t.Context(), work.attempt, c)
				return err
			})
			got, err := loadAdapterResult(t.Context(), s.conn, id)
			if err != nil || got.Candidate.Payload["answer"] != Text(fmt.Sprint(order[0])) {
				t.Fatalf("retained adapter result = %+v, %v", got, err)
			}
		})
	}
}

func checkConflictingOffers(t *testing.T, path string, order [2]int64, offer func(int64) error) {
	t.Helper()
	if err := offer(order[0]); err != nil {
		t.Fatal(err)
	}
	before, err := Inspect(t.Context(), path, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	if err := offer(order[1]); !errors.Is(err, ErrConflict) {
		t.Fatalf("second offer returned %v, want ErrConflict", err)
	}
	after, err := Inspect(t.Context(), path, DefaultLimits())
	if err != nil || !reflect.DeepEqual(before.Tables, after.Tables) {
		t.Fatalf("rejected offer changed the retained tables or failed inspection: %v", err)
	}
	if _, err := Replay(t.Context(), path, DefaultLimits()); err != nil {
		t.Fatal(err)
	}
}
