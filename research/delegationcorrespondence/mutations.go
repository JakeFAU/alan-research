package delegationcorrespondence

import "encoding/json"

type RawMutation struct {
	Applicable                          bool
	Diagnostic                          string
	Name, OriginalSHA256, MutatedSHA256 string
	SemanticField                       string
	OriginalIDs, MutatedIDs             []string
	Changes                             []RowRef
	CellChanges                         []RawCellChange
	Input                               json.RawMessage
	Verdict                             RawVerdict
}

type RawCellChange struct {
	Table, Key, Column string
	Before, After      json.RawMessage
}

// RawMutations preserves complete mutated bytes and unchanged result rows.
// Changed typed witnesses are canonically encoded and their attempt/gate hash
// mirrors are recomputed where relevant. Encoding integrity is not the oracle.
func RawMutations(original []byte) []RawMutation {
	out := []RawMutation{}
	names := []string{"missing-grant", "forged-grant-origin", "altered-holder", "altered-recipient", "altered-root", "altered-operation", "altered-resource", "altered-request", "altered-generation", "reused-grant", "result-as-authority", "equal-payload-result-substitution", "gate-before-reservation", "stale-activation", "forged-gate-attempt"}
	for _, name := range names {
		var raw RawSnapshot
		if strictJSON(original, &raw) != nil {
			return nil
		}
		m := RawMutation{Name: name, OriginalSHA256: sha(original)}
		m.Diagnostic = mutationPrerequisite(name, raw)
		if m.Diagnostic != "" {
			m.Input = append(json.RawMessage(nil), original...)
			m.MutatedSHA256 = m.OriginalSHA256
			out = append(out, m)
			continue
		}
		m.Applicable = true
		update := func(t string, r Row, col string, v any) {
			before := rawCell(r, col)
			after := bytesJSON(v)
			r[col] = bytesJSON(string(after))
			key := rowKey(t, r)
			m.Changes = append(m.Changes, RowRef{t, key, col, sha(before)})
			m.CellChanges = append(m.CellChanges, RawCellChange{t, key, col, before, after})
		}
		switch name {
		case "missing-grant":
			for _, r := range raw.Tables["lifecycle_results"] {
				m.OriginalIDs = append(m.OriginalIDs, textCell(r, "id"))
				m.Changes = append(m.Changes, RowRef{"lifecycle_results", textCell(r, "id"), "result", sha(rawCell(r, "result"))})
			}
			raw.Tables["lifecycle_results"] = []Row{}
		case "forged-grant-origin":
			r := raw.Tables["lifecycle_results"][0]
			q := decodeCell[LifecycleRequest](r, "descriptor")
			m.OriginalIDs = append(m.OriginalIDs, q.Intent.ID)
			q.Intent.Run.Subject = "a"
			update("lifecycle_results", r, "descriptor", q)
		case "altered-recipient":
			r := raw.Tables["lifecycle_results"][0]
			q := decodeCell[LifecycleRequest](r, "descriptor")
			var claim Claim
			if strictJSON([]byte(q.Intent.Payload["claim"].Text), &claim) != nil {
				continue
			}
			m.SemanticField = "lifecycle_results.descriptor.Intent.Payload.claim.Text.Recipient.Subject"
			m.OriginalIDs = []string{q.Intent.ID, identity(claim.Recipient)}
			claim.Recipient.Subject = "b"
			m.MutatedIDs = []string{q.Intent.ID, identity(claim.Recipient)}
			q.Intent.Payload["claim"] = stringValue(string(bytesJSON(claim)))
			update("lifecycle_results", r, "descriptor", q)
			// The claimed decision also names the changed recipient, but its
			// committed source intent/history remain the retained original.
			result := decodeCell[LifecycleResult](r, "result")
			result.Request = q
			update("lifecycle_results", r, "result", result)
		case "equal-payload-result-substitution":
			rs := raw.Tables["external_receipts"]
			if len(rs) < 2 {
				continue
			}
			a := decodeCell[Receipt](rs[0], "descriptor")
			b := decodeCell[Receipt](rs[1], "descriptor")
			m.OriginalIDs = []string{a.Key, b.Key}
			a.Key = b.Key
			m.MutatedIDs = []string{a.Key}
			update("external_receipts", rs[0], "descriptor", a)
		case "gate-before-reservation":
			rs := raw.Tables["composition_facts"]
			var gate, res Row
			for _, r := range rs {
				f := decodeCell[Fact](r, "descriptor")
				if f.Kind == "reservation" && res == nil {
					res = r
				}
				if f.Kind == "authority_admission" && gate == nil {
					gate = r
				}
			}
			if gate == nil || res == nil {
				continue
			}
			g := decodeCell[Fact](gate, "descriptor")
			r := decodeCell[Fact](res, "descriptor")
			g.Order, r.Order = r.Order, g.Order
			gate["decision_order"] = bytesJSON(g.Order)
			res["decision_order"] = bytesJSON(r.Order)
			update("composition_facts", gate, "descriptor", g)
			update("composition_facts", res, "descriptor", r)
		case "stale-activation", "forged-gate-attempt":
			for _, r := range raw.Tables["composition_facts"] {
				f := decodeCell[Fact](r, "descriptor")
				if f.Kind != "authority_admission" {
					continue
				}
				m.OriginalIDs = append(m.OriginalIDs, f.OperationID)
				if name == "stale-activation" {
					f.AuthorityAdmission.Source.Epoch = 0
				} else {
					f.AuthorityAdmission.AttemptDigest = sha([]byte("other attempt"))
				}
				f.RecordDigest = digest("alan/authority-admission/v1", f.AuthorityAdmission)
				update("composition_facts", r, "descriptor", f)
				break
			}
		default:
			for _, r := range raw.Tables["external_attempts"] {
				a := decodeCell[ExternalAttempt](r, "descriptor")
				if a.Authority == nil || a.Authority.Use == nil || a.Authority.Use.Path != "delegated" {
					continue
				}
				m.OriginalIDs = []string{a.Intent.ID, a.Authority.Use.GrantID, a.Authority.OriginalRequestID}
				switch name {
				case "altered-holder":
					m.SemanticField = "external_attempts.descriptor.Authority.Use.Root.Run.Subject"
					a.Authority.Use.Root.Run.Subject = "a"
				case "altered-root":
					a.Authority.RootID = sha([]byte("absent root"))
				case "altered-operation":
					a.Authority.ServiceOperationID = identity("a", 1, 5, 0)
				case "altered-resource":
					a.Authority.Use.Root.Permission.Target = "https://other-resource.invalid"
				case "altered-request":
					a.Authority.OriginalRequestID = identity("a", 1, 1, 0)
				case "altered-generation":
					a.Authority.Use.Root.Run.Generation = 2
				case "reused-grant":
					a.Authority.ServiceOperationID = identity("b", 1, 9, 0)
				case "result-as-authority":
					a.Authority.Use.GrantID = identity("b", 1, 9, 0)
				}
				m.MutatedIDs = []string{a.Authority.Use.GrantID, a.Authority.OriginalRequestID, a.Authority.ServiceOperationID}
				update("external_attempts", r, "descriptor", a)
				newDigest := digest("alan/attempt-row/v1", a)
				for _, sr := range raw.Tables["settlements"] {
					if textCell(sr, "intent_id") != a.Intent.ID {
						continue
					}
					var ss struct {
						Attempt     ExternalAttempt
						Disposition string
						Receipt     *Receipt
					}
					if strictJSON(rawCell(sr, "descriptor"), &ss) == nil {
						ss.Attempt = a
						update("settlements", sr, "descriptor", ss)
					}
				}
				for _, rr := range raw.Tables["composition_reservations"] {
					if textCell(rr, "operation_id") != a.Intent.ID {
						continue
					}
					v := decodeCell[Reservation](rr, "descriptor")
					v.AttemptDigest = newDigest
					update("composition_reservations", rr, "descriptor", v)
					for _, fr := range raw.Tables["composition_facts"] {
						f := decodeCell[Fact](fr, "descriptor")
						if f.OperationID != a.Intent.ID {
							continue
						}
						if f.Kind == "reservation" {
							f.RecordDigest = digest("alan/reservation-row/v1", v)
							update("composition_facts", fr, "descriptor", f)
						} else if f.Kind == "authority_admission" {
							f.AuthorityAdmission.AttemptDigest = newDigest
							f.RecordDigest = digest("alan/authority-admission/v1", f.AuthorityAdmission)
							update("composition_facts", fr, "descriptor", f)
						}
					}
				}
				break
			}
		}
		m.Input = bytesJSON(raw)
		m.MutatedSHA256 = sha(m.Input)
		m.Verdict = CheckRaw(m.Input)
		out = append(out, m)
	}
	return out
}

// mutationPrerequisite reports absent evidence without treating an unattempted
// mutation as a policy rejection. It does not manufacture a completed prefix.
func mutationPrerequisite(name string, raw RawSnapshot) string {
	switch name {
	case "missing-grant", "forged-grant-origin", "altered-recipient":
		if len(raw.Tables["lifecycle_results"]) == 0 {
			return "missing lifecycle result"
		}
	case "equal-payload-result-substitution":
		if len(raw.Tables["external_receipts"]) < 2 {
			return "requires two external receipts"
		}
	case "gate-before-reservation", "stale-activation", "forged-gate-attempt":
		gate, reservation := false, false
		for _, r := range raw.Tables["composition_facts"] {
			f := decodeCell[Fact](r, "descriptor")
			gate = gate || f.Kind == "authority_admission"
			reservation = reservation || f.Kind == "reservation"
		}
		if !gate {
			return "missing authority admission"
		}
		if name == "gate-before-reservation" && !reservation {
			return "missing reservation"
		}
	default:
		for _, r := range raw.Tables["external_attempts"] {
			a := decodeCell[ExternalAttempt](r, "descriptor")
			if a.Authority != nil && a.Authority.Use != nil && a.Authority.Use.Path == "delegated" {
				return ""
			}
		}
		return "missing delegated external attempt"
	}
	return ""
}
