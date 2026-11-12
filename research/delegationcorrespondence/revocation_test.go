package delegationcorrespondence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"

	m "github.com/JakeFAU/alan-research/research/authorityservice"
	b "github.com/JakeFAU/alan-research/research/mutualobligations"
)

func revocationFixture(t *testing.T, mode, name string) string {
	t.Helper()
	root := os.Getenv("ALAN_REVOCATION_CORPUS")
	if root == "" {
		t.Skip("set ALAN_REVOCATION_CORPUS to the verified extraction root")
	}
	data, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var entries map[string]struct{ Cases map[string]string }
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatal(err)
	}
	id := entries["raw-"+mode+".tar.gz"].Cases[name]
	if id == "" {
		t.Fatalf("missing corpus locator %s/%s", mode, name)
	}
	return filepath.Join(root, "raw-"+mode, id)
}

func TestRevocationSealed(t *testing.T) {
	for _, mode := range []string{"normal", "race"} {
		for _, tc := range []struct {
			name     string
			orders   []int64
			contacts int
			captures int
			legacy   bool
		}{
			{"N05", []int64{7, 17, 21}, 0, 3, false},
			{"N06", []int64{8, 16, 20, 22}, 0, 3, false},
			{"N07", []int64{8, 16, 17, 21, 23, 33, 34, 35}, 2, 3, false},
			{"P02-v2", []int64{8, 16, 17, 18, 28, 29, 30}, 2, 1, true},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				report, err := CheckRevocationDirectory(revocationFixture(t, mode, tc.name))
				if err != nil {
					t.Fatal(err)
				}
				if !report.Supported || !report.RawValid || !report.Physical.Valid || report.Legacy.ContinuityFailure != nil {
					t.Fatalf("case outside expected support: reasons=%v physical=%+v", report.Reasons, report.Physical)
				}
				if report.Legacy.Projection.Supported != tc.legacy || report.Legacy.External.Contacts != tc.contacts ||
					report.Legacy.External.Applications != tc.contacts || len(report.Captures) != tc.captures {
					t.Fatalf("legacy/evidence mismatch for %s", tc.name)
				}
				for _, cut := range report.Captures {
					if !cut.Supported || cut.Audit == nil || !cut.Audit.ProtocolValid || !cut.Audit.AuthoritySafe {
						t.Fatalf("unsupported cut %s: %+v", cut.Path, cut.Failure)
					}
				}
				final := report.Captures[len(report.Captures)-1]
				var orders []int64
				for _, e := range final.Events {
					switch e.Event.Kind {
					case "retain-grant", "authorize", "revoke", "dispatch-admit", "dispatch-refuse", "apply":
						orders = append(orders, e.Order)
					}
				}
				if !slices.Equal(orders, tc.orders) {
					t.Fatalf("mapped event order %v, want %v", orders, tc.orders)
				}
				if tc.name == "N07" {
					for _, cut := range report.Captures[:2] {
						if !slices.Equal(cut.Unresolved, []b.Obligation{b.BA}) || cut.State.Base.Obligations[b.BA].Applied {
							t.Fatalf("cut %s lost unresolved opportunity or invented apply", cut.Path)
						}
					}
				}
				ba, ab := final.Principals[b.BA], final.Principals[b.AB]
				if final.Roots[b.BA].Model != "B/g1:b" || final.Roots[b.AB].Model != "B/g1:c" ||
					final.Roots[b.BA].Concrete == final.Roots[b.AB].Concrete ||
					final.Roots[b.BA].Composition != final.Roots[b.AB].Composition ||
					ba.AuthorizationSource.Epoch != 1 || ba.AuthorizationSource.Revision != 1 {
					t.Fatalf("lost complete root map or runtime source: %+v %+v", final.Roots, ba.AuthorizationSource)
				}
				switch tc.name {
				case "N05":
					if ba.Reservation.Allowed || ba.Debit != nil || len(ba.Gates) != 0 || ba.Receipt != nil ||
						final.State.Authority[b.BA].Decision != m.Refused {
						t.Fatal("early refusal gained a debit, gate, or result")
					}
				case "N06":
					if !ba.Reservation.Allowed || ba.Debit == nil || ba.Debit.Effects != 1 || ba.Debit.DecisionOrder != 16 ||
						ba.Debit.RootID == final.Roots[b.BA].Composition || ba.Authority.RequestConsumedSequence != 2 ||
						ba.Authority.GrantDecisionSequence != 1 || len(ba.Gates) != 1 || ba.Gates[0].AuthorityAdmission.Allowed ||
						ba.Gates[0].AuthorityAdmission.Reason != "authority_revoked" || ba.Receipt != nil ||
						ba.Exposure == nil || ba.Exposure.Contacts != 0 || ba.Exposure.Applications != 0 ||
						final.State.Authority[b.BA].Decision != m.Allowed || !final.State.Authority[b.BA].AdmissionRefused {
						t.Fatal("historical allow/debit and fresh zero-contact refusal diverged")
					}
				case "N07":
					if ba.Receipt == nil || ab.Receipt == nil || !final.State.Revoked ||
						final.State.Authority[b.AB].Path != m.OwnRoot || ab.SelectedPath.Root != final.Roots[b.AB].Concrete ||
						ba.ReplySelection.Order != 28 || ab.ReplySelection.Order != 40 {
						t.Fatal("old BA opportunity or unaffected AB authority failed")
					}
				}
				if (tc.contacts == 2) != final.State.Base.Obligations[b.AB].Consumed ||
					(tc.contacts == 2) != final.State.Base.Obligations[b.BA].Consumed {
					t.Fatalf("unexpected principal consumption in %s", tc.name)
				}
			})
		}
		t.Run(mode+"/outside-domain", func(t *testing.T) {
			report, err := CheckRevocationDirectory(revocationFixture(t, mode, "applied-unknown-revoked-before-fresh-gate"))
			if err != nil {
				t.Fatal(err)
			}
			if !report.RawValid || !report.Physical.Valid || report.Supported ||
				report.Physical.Contacts != 1 || report.Physical.Applications != 1 || len(report.Captures) != 5 ||
				!slices.Contains(report.Reasons, "whole_store_loss") || !slices.Contains(report.Reasons, "applied_without_retained_result") {
				t.Fatalf("outside-domain verdict: support=%v raw=%v physical=%+v reasons=%v", report.Supported, report.RawValid, report.Physical, report.Reasons)
			}
			for _, cut := range report.Captures {
				if cut.State.Base.Obligations[b.BA].Applied || len(cut.Events) != 0 {
					t.Fatalf("invented BA application at %s", cut.Path)
				}
			}
		})
	}
}
