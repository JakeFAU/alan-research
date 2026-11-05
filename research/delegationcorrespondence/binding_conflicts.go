package delegationcorrespondence

// ConflictPromotionReport qualifies appended losing evidence against an exact
// public parent. It does not establish a public conflict-production schedule,
// new service custody or another endpoint application.
type ConflictPromotionReport struct {
	Valid                                bool
	Failure                              *Failure
	Parent, OpenedSHA256, PromotedSHA256 string
	ParentVerdict                        RawVerdict
	Conflicts                            int
	Scope                                string
}

func CheckConflictPromotion(parent, opened, promoted []byte) (out ConflictPromotionReport) {
	out.Parent, out.OpenedSHA256, out.PromotedSHA256 = sha(parent), sha(opened), sha(promoted)
	out.Scope = "ancillary internal-writer compatibility; no public conflict-production schedule or additional endpoint application established"
	defer func() {
		if p := recover(); p != nil {
			if f, ok := p.(Failure); ok {
				out.Failure = &f
				out.Valid = false
			} else {
				panic(p)
			}
		}
	}()
	out.ParentVerdict = CheckRaw(parent)
	need(out.ParentVerdict.Valid, "conflict_parent_invalid", "", "")
	var original, before RawSnapshot
	var export struct {
		RawSnapshot
		ReceiptEvidence []struct{ IntentID, State string }
	}
	need(strictJSON(parent, &original) == nil && strictJSON(opened, &before) == nil && strictJSON(promoted, &export) == nil, "conflict_snapshot_encoding", "", "")
	checkOwnerRenewal(original, before, "promotion-open")
	after := export.RawSnapshot
	need(len(after.Tables["metadata"]) == 1, "conflict_metadata_shape", "", "")
	schema := intCell(before.Tables["metadata"][0], "schema_version")
	need((schema == 13 || schema == 15) && intCell(after.Tables["metadata"][0], "schema_version") == schema+1, "conflict_promotion_version", "", "")
	need(before.Namespace == after.Namespace && before.KeyVersion == after.KeyVersion && equal(before.Roots, after.Roots), "conflict_identity", "", "")
	need(len(after.Tables) == len(before.Tables)+1 && len(after.Tables["receipt_conflicts"]) > 0, "conflict_table_set", "", "")
	for table, rows := range before.Tables {
		if table == "metadata" {
			copyRow := Row{}
			for k, v := range after.Tables[table][0] {
				copyRow[k] = v
			}
			copyRow["schema_version"] = bytesJSON(schema)
			need(equal(rows, []Row{copyRow}), "conflict_metadata_changed", "", "")
			continue
		}
		need(equal(rows, after.Tables[table]), "conflict_parent_rows_changed", table, "")
	}
	winners := map[string]Row{}
	attempts := map[string]AdapterAttempt{}
	for _, r := range before.Tables["adapter_results"] {
		winners[textCell(r, "intent_id")] = r
	}
	for _, r := range before.Tables["adapter_attempts"] {
		attempts[textCell(r, "intent_id")] = decodeCell[AdapterAttempt](r, "descriptor")
	}
	seen := map[string]bool{}
	for _, r := range after.Tables["receipt_conflicts"] {
		need(len(r) == 4, "conflict_row_shape", "", "")
		id := textCell(r, "intent_id")
		winner, ok := winners[id]
		need(ok, "conflict_original_result", id, "")
		bytes := rawCell(r, "descriptor")
		offer := decodeCell[AdapterResult](r, "descriptor")
		copyOffer := offer
		copyOffer.Digest = ""
		key := id + textCell(r, "offer_digest")
		need(!seen[key] && textCell(r, "winner_digest") == sha(rawCell(winner, "descriptor")) && textCell(r, "offer_digest") == sha(bytes) && sha(bytes) != sha(rawCell(winner, "descriptor")), "conflict_digests", id, "")
		need(offer.Digest == sha(bytesJSON(copyOffer)) && offer.Candidate.ID == id && offer.Candidate.Key == attempts[id].Prepared.Request.Key, "conflict_offered_origin", id, "")
		seen[key] = true
		out.Conflicts++
	}
	need(len(export.ReceiptEvidence) == len(winners), "conflict_evidence_count", "", "")
	for _, e := range export.ReceiptEvidence {
		need(winners[e.IntentID] != nil && e.State == "conflict", "conflict_evidence_summary", e.IntentID, "")
	}
	out.Valid = true
	return
}
