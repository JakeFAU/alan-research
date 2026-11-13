package delegationcorrespondence

import (
	"slices"
)

func checkEffectBinding(phase PhaseReport, files map[string][]byte, entry revocationManifest, loss bool, index int) EffectConcreteBinding {
	raw := effectRaw(files, phase.Path)
	rows := indexRevocationRows(raw)
	roots := effectCompositionRoots(raw)
	mapped := revocationRoots(raw, entry, files, roots[1])
	rec := entry.Manifest.Reciprocal
	id := rec.BAService
	originalEpoch, currentEpoch := int64(1), int64(1)
	if loss {
		originalEpoch = 2
		currentEpoch = 2
		if index >= 3 {
			currentEpoch = 3
		}
	}
	for _, row := range raw.Tables["runs"] {
		need(intCell(row, "epoch") == currentEpoch, "effect_cut_epoch", phase.Path, "")
	}
	originals := phase.Verdict.Custody.Originals
	need(len(originals) == 2, "effect_original_pair", phase.Path, "")
	o := originals[1]
	need(o.ID == rec.BARequest && o.Caller == revocationRuns[1] && o.Callee == revocationRuns[0] && o.Committed && o.Admitted && o.Selected && o.EffectID == id,
		"effect_selected_original", phase.Path, "")
	out := EffectConcreteBinding{Original: o, ServiceID: id, Roots: mapped, Authorization: rows.auth[id], Attempt: rows.attempts[id],
		Reservation: rows.reservations[id], Debit: rows.charges[id], Gates: []Fact{}}
	grant, ok := rows.life[rec.Grant]
	need(ok, "effect_retained_grant", phase.Path, "")
	out.Grant = grant
	checkEffectControl(grant, rows.auth[rec.Grant], entry, mapped[1], roots[0], raw.Namespace, originalEpoch)
	need(grant.Request.Operation == "DelegateAuthority", "effect_retained_grant", phase.Path, "")
	effect := EffectReport{}
	for _, e := range phase.Verdict.Effects {
		if e.ID == id {
			effect = e
		}
	}
	auth := out.Authorization
	attempt, reservation, debit := out.Attempt, out.Reservation, out.Debit
	source := Source{Run: revocationRuns[0], Epoch: originalEpoch, Sequence: 4, Revision: 1}
	use := Use{Version: bindingVersion, Path: "delegated", Root: mapped[1].Concrete, GrantID: rec.Grant}
	need(auth.IntentID == id && auth.Allowed && auth.Source == source && attempt.AuthorizationID == auth.ID && attempt.Allowed &&
		attempt.Request.Subject == source.Run.Subject && attempt.Request.Generation == source.Run.Generation && attempt.Request.Epoch == source.Epoch &&
		attempt.Request.Sequence == source.Sequence && attempt.Request.Authorization == auth.ID && attempt.Intent.Authority != nil && *attempt.Intent.Authority == use,
		"effect_immutable_attempt", phase.Path, "")
	w := effect.Witness
	need(effect.Allowed && effect.ServingCharge == 1 && w != nil && w.Allowed && w.Use != nil && *w.Use == use && w.RootID == roots[1] && w.ServingRootID == roots[0] &&
		w.OriginalRequestID == o.ID && w.ServiceOperationID == id && w.RequestObservationID == identity("observation", "message", o.ID) &&
		w.RequestConsumedSequence == o.RequestSelection && w.GrantDecisionSequence == grant.DecisionSequence,
		"effect_selected_grant_path", phase.Path, "")
	need(reservation.OperationID == id && reservation.RootID == roots[0] && reservation.Allowed && reservation.AttemptDigest != "" &&
		debit == (Charge{RootID: roots[0], OperationID: id, Kind: "effect", Effects: 1, DecisionOrder: effect.Order}),
		"effect_immutable_debit", phase.Path, "")
	selected := rows.history[historyKey(o.Callee, o.RequestSelection)]
	need(selected.Selected == identity("observation", "message", o.ID) && selected.Sequence < source.Sequence && rows.admitted[o.ID] != "",
		"effect_request_custody", phase.Path, "")
	var revokeOrder, resultOrder, consumeOrder int64
	controls := 0
	for _, fact := range phase.Verdict.OrderedFacts {
		if fact.Kind == "lifecycle" {
			control := rows.life[fact.OperationID]
			epoch := originalEpoch
			if control.Request.Operation == "RevokeAuthority" {
				epoch = currentEpoch
			}
			checkEffectControl(control, rows.auth[fact.OperationID], entry, mapped[1], roots[0], raw.Namespace, epoch)
			controls++
			if control.Request.Operation == "RevokeAuthority" {
				need(out.Revoke == nil, "effect_one_revoke", phase.Path, "")
				out.Revoke = &control
				revokeOrder = fact.Order
			}
		}
		if fact.OperationID == id && fact.Kind == "authority_admission" {
			out.Gates = append(out.Gates, fact)
		}
		if fact.OperationID == id && fact.Kind == "external_result" {
			resultOrder = fact.Order
		}
		if fact.Kind == "execution" && fact.Run == o.Caller && o.ReplySelected && fact.Sequence == o.ReplySelection {
			out.Consumption = &fact
			consumeOrder = fact.Order
		}
	}
	shouldRevoke := (!loss && index >= 1) || (loss && index == 4)
	need((out.Revoke != nil) == shouldRevoke && controls == 1+boolInt(shouldRevoke), "effect_revoke_cut", phase.Path, "")
	wantGates := 1
	if loss && index == 4 {
		wantGates = 2
	}
	need(len(out.Gates) == wantGates && len(effect.Gates) == wantGates, "effect_gate_count", phase.Path, "")
	gate := out.Gates[0].AuthorityAdmission
	need(gate != nil && gate.Allowed && gate.Reason == "authorized" && gate.Ordinal == 1 && gate.Source == source &&
		gate.AttemptDigest == reservation.AttemptDigest && out.Gates[0].Order > effect.Order, "effect_first_gate", phase.Path, "")
	if shouldRevoke {
		need(revokeOrder > out.Gates[0].Order, "effect_revoke_after_gate", phase.Path, "")
	}
	if wantGates == 2 {
		fresh := out.Gates[1].AuthorityAdmission
		wantSource := source
		wantSource.Epoch = currentEpoch
		need(fresh != nil && !fresh.Allowed && fresh.Reason == "authority_revoked" && fresh.Ordinal == 2 && fresh.Source == wantSource &&
			fresh.AttemptDigest == gate.AttemptDigest && out.Gates[1].Order > revokeOrder, "effect_fresh_gate", phase.Path, "")
	}
	if shouldRevoke {
		side := effectJSON[struct {
			Decision      LifecycleResult
			Authorization Authorization
			Admission     AdmissionResult
		}](files, "controlled-revocation.json")
		need(equal(side.Decision, *out.Revoke) && equal(side.Authorization, rows.auth[out.Revoke.Request.Intent.ID]) &&
			side.Admission.Accepted && side.Admission.Bound != nil && *side.Admission.Bound == revocationRuns[1] &&
			side.Admission.ObservationID == identity("observation", "lifecycle", out.Revoke.Request.Intent.ID),
			"effect_controlled_revoke_join", phase.Path, "")
	}
	receipt, retained := rows.receipts[id]
	wantResult := !loss && index == 2
	need(retained == wantResult && effect.Result == wantResult && (resultOrder > 0) == wantResult && o.ResultSelected == wantResult && o.ReplySelected == wantResult,
		"effect_local_result_cut", phase.Path, "")
	if retained {
		endpoint := checkEffectEndpointPair(files, attempt.Request)
		need(receipt == endpoint && resultOrder > revokeOrder && consumeOrder > resultOrder && out.Consumption != nil &&
			o.ReplyID == rec.BAReply && rows.history[historyKey(o.Caller, o.ReplySelection)].Selected == identity("observation", "message", rec.BAReply),
			"effect_authentic_consumption", phase.Path, "")
		out.Receipt = &receipt
	}
	settlements := 0
	for _, row := range raw.Tables["settlements"] {
		if textCell(row, "intent_id") == id {
			settlements++
		}
	}
	need(settlements == boolInt(wantResult), "effect_external_settlement", phase.Path, "")
	// Continuity already protects append-only raw rows. Explicitly keep the
	// selected attempt, authorization, reservation and debit equal at every cut.
	final := indexRevocationRows(effectRaw(files, "final/inspect.json"))
	need(equal(attempt, final.attempts[id]) && equal(auth, final.auth[id]) && equal(reservation, final.reservations[id]) && equal(debit, final.charges[id]),
		"effect_historical_binding_changed", phase.Path, "")
	return out
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Like the old epoch-1 relation's control checks, but verifies actual source
// incarnation instead of rewriting historical source or wire coordinates.
func checkEffectControl(l LifecycleResult, a Authorization, entry revocationManifest, root RevocationRoot, recipientRootID, namespace string, epoch int64) {
	i := l.Request.Intent
	need(l.Accepted && l.Authority != nil && l.Authority.Version == bindingVersion && l.Authority.Namespace == namespace &&
		i.Run == revocationRuns[1] && a.Allowed && a.Source == (Source{Run: i.Run, Epoch: epoch, Sequence: i.Sequence, Revision: 1}) &&
		l.Before != nil && l.Before.Epoch == epoch && equal(l.Before, l.After) && l.Request.SourceRevision == 1 && l.Request.ExpectedRevision == 1 && l.Successor == 0,
		"effect_exact_control_source", i.ID, "")
	need(l.Authority.Root == root.Concrete && l.Authority.RootID == root.Composition, "effect_exact_control_root", i.ID, "")
	rec := entry.Manifest.Reciprocal
	switch l.Request.Operation {
	case "DelegateAuthority":
		var claim Claim
		want := Claim{Version: bindingVersion, Root: root.Concrete, Recipient: revocationRuns[0], RecipientRoot: revocationRuns[0], RecipientRevision: 1,
			OriginalRequestID: rec.BARequest, ServiceOperationID: rec.BAService}
		need(strictJSON([]byte(i.Payload["claim"].Text), &claim) == nil && claim == want && rec.Claim == want && i.ID == rec.Grant &&
			l.Reason == "delegated" && i.Target == "a" && i.Generation == 1 && l.Authority.RecipientRootID == recipientRootID && l.DecisionSequence == 1,
			"effect_exact_grant", i.ID, "")
	case "RevokeAuthority":
		var revoke Revoke
		need(strictJSON([]byte(i.Payload["claim"].Text), &revoke) == nil && revoke == (Revoke{Version: bindingVersion, Root: root.Concrete}) &&
			l.Reason == "revoked" && i.Target == "b" && i.Generation == 1 && l.Authority.RecipientRootID == "" && l.DecisionSequence == 2,
			"effect_exact_revoke", i.ID, "")
	default:
		fail("effect_control_outside_domain", i.ID, l.Request.Operation)
	}
}

func effectCertificate(files map[string][]byte, b EffectConcreteBinding, loss bool, index int) EffectCausalCertificate {
	receipt := checkEffectEndpointPair(files, b.Attempt.Request)
	cert := EffectCausalCertificate{Kind: "synchronous_pretransport_callback", EndpointAtCut: "complete zero ingress/application before worker launch",
		Artifacts:   effectArtifactHashes(files, "case.json", "source.json", "binding-source.json", "endpoint-b-ready.json", "endpoint-b.jsonl", "endpoints.json"),
		Constraints: []string{}, InferredPrerequisites: []string{}, Unobserved: []string{}, FinalCorroboration: &receipt,
		OriginalEpoch: b.Attempt.Request.Epoch, CurrentEpoch: b.Attempt.Request.Epoch, EpochMap: map[int64]uint8{b.Attempt.Request.Epoch: 1}}
	if loss {
		cert.Kind = "post_commit_acknowledged_owner_loss"
		cert.EndpointAtCut = "complete one ingress/one application certified before death; no local authentic result"
		cert.Constraints = []string{"endpoint ingress < application commit < endpoint acknowledgment < acknowledged SIGKILL",
			"runtime in-flight acknowledgment < acknowledged SIGKILL",
			"acknowledged SIGKILL < stopped capture < controller endpoint release < Open < reopened capture < exact b revoke < fresh refused gate",
			"successful authentic receipt emission requires successful post-commit barrier release; controller releases only after death"}
		cert.Unobserved = []string{"relative emission order of endpoint and runtime acknowledgments",
			"continuous barrier occupancy and transport liveness until death; cancellation/timeout may return an error earlier",
			"a complete physical chronology; witness steps are an existential abstraction"}
		cert.InferredPrerequisites = []string{"launch precedes matched ingress/application; this prerequisite is not a separately retained launch record"}
		cert.EpochMap[3] = 2
		if index >= 3 {
			cert.CurrentEpoch = 3
		}
		for k, v := range effectArtifactHashes(files, "acknowledgments.json", "kill.json", "first-owner/open.json", "stopped/inspect.json", "reopened/inspect.json") {
			cert.Artifacts[k] = v
		}
	} else {
		cert.Constraints = []string{"first allowed-gate callback synchronously captures before-revoke, performs exact revoke, captures after-revoke, then returns",
			"worker launch follows callback return; neither intermediate cut has contact/application or local result"}
		cert.Unobserved = []string{"final journals alone cannot date application at either intermediate cut"}
		cert.Artifacts["temporal-boundary.json"] = sha(files["temporal-boundary.json"])
		if index == 2 {
			cert.Kind = "final_retained_authentic_result"
			cert.EndpointAtCut = "complete one ingress/one application; separately retained result and actual principal consumption"
			cert.InferredPrerequisites = []string{"launch before ingress and authentic response before retention; reply transport hidden by model"}
		}
	}
	cert.Artifacts["controlled-revocation.json"] = sha(files["controlled-revocation.json"])
	// Certificates distinguish final corroboration from at-cut knowledge above.
	cert.Constraints = slices.Clone(cert.Constraints)
	return cert
}

func effectCompositionRoots(raw RawSnapshot) (ids [2]string) {
	need(len(raw.Tables["runs"]) == 2 && len(raw.Tables["subjects"]) == 2 && len(raw.Tables["composition_roots"]) == 2 && len(raw.Tables["composition_static_members"]) == 2, "effect_composition_pair", raw.Namespace, "")
	for _, r := range raw.Tables["subjects"] {
		id := RunID{Subject: textCell(r, "id"), Generation: intCell(r, "generation")}
		need(slices.Contains(revocationRuns[:], id) && intCell(r, "revision") == 1, "effect_subject_revision", id.Subject, "")
	}
	for _, row := range raw.Tables["composition_roots"] {
		root := decodeCell[Root](row, "descriptor")
		i := slices.Index(revocationRuns[:], root.Run)
		need(i >= 0 && slices.Equal(root.Members, []RunID{root.Run}) && root.ID != "", "effect_composition_pair", root.ID, "")
		ids[i] = root.ID
	}
	need(ids[0] != "" && ids[1] != "" && ids[0] != ids[1], "effect_composition_pair", raw.Namespace, "")
	return ids
}
