package delegationcorrespondence

// Only ordinary closure/retirement and the composition Retry obstruction are
// in this fixed-generation domain. No generation/revision writer is inferred.
func (c *checker) checkOrdinaryLifecycle(a LifecycleResult) {
	q := a.Request
	i := q.Intent
	target := RunID{i.Target, i.Generation}
	before := a.Before
	need(before != nil && before.Run == target && before.ProgramID == c.runs[target].ProgramID && before.Owner == c.subjects[target].Owner && before.Revision == c.subjects[target].Revision && before.Retired == c.subjects[target].Retired && before.RetryCancelled == c.subjects[target].RetryCancelled && before.Sequence == c.seq[target] && equal(before.Projection, c.states[target]), "lifecycle_target_prefix", i.ID, "")
	c.epoch(Source{target, before.Epoch, before.Sequence, before.Revision})
	after := *before
	permitted := false
	for _, p := range c.permissions {
		permitted = permitted || p.Kind == "lifecycle" && p.Source == i.Run.Subject && p.Target == i.Target && p.Operation == q.Operation && p.Revision == q.ExpectedRevision
	}
	accepted, reason := false, ""
	switch {
	case q.Operation == "Retry":
		reason = "invocation_retry_unsupported"
	case !permitted:
		reason = "grant_denied"
	case before.Revision != q.ExpectedRevision:
		reason = "revision_changed"
	case before.Retired:
		reason = "retired"
	default:
		accepted = true
		reason = "already_terminal"
		if q.Operation == "RetireSubject" {
			after.Retired = true
			reason = "retired"
		}
		if before.Projection.Status == "open" {
			after.Projection.Status = "cancelled"
			if q.Operation == "AbandonRun" {
				after.Projection.Status = "abandoned"
			}
			if q.Operation != "RetireSubject" {
				reason = "closed"
			}
		}
	}
	need(a.Authority == nil && a.Successor == 0 && a.Accepted == accepted && a.Reason == reason && equal(a.After, &after), "ordinary_lifecycle_policy", i.ID, reason)
	if accepted {
		c.states[target] = after.Projection
		s := c.subjects[target]
		s.Retired = after.Retired
		c.subjects[target] = s
	}
	c.out.Discrepancies = append(c.out.Discrepancies, "ordinary_lifecycle_outside_restricted_correspondence")
}

func (c *checker) checkSourcePolicy(a Authorization) {
	reason := "authorized"
	closure := ""
	if c.subjects[a.Source.Run].Retired {
		reason = "source_retired"
	} else if c.states[a.Source.Run].Status != "open" {
		reason = "source_closed"
	}
	if reason != "authorized" {
		var latest int64
		for id, l := range c.life {
			if c.has("lifecycle", id) && l.Accepted && l.After != nil && l.After.Run == a.Source.Run && l.After.Projection.Status != "open" && l.DecisionSequence > latest {
				closure = id
				latest = l.DecisionSequence
			}
		}
	}
	need(a.Allowed == (reason == "authorized") && a.Reason == reason && a.ClosureOperation == closure, "source_policy", a.IntentID, reason)
}

func (c *checker) checkCurrentSubjects() {
	for _, r := range c.raw.Tables["subjects"] {
		id := RunID{textCell(r, "id"), intCell(r, "generation")}
		s := c.subjects[id]
		need((intCell(r, "retired") == 0 || intCell(r, "retired") == 1) && (intCell(r, "retired") == 1) == s.Retired, "current_retirement", id.Subject, "")
	}
	closures := map[string]bool{}
	for _, r := range c.raw.Tables["composition_closures"] {
		var d struct {
			RootID, LifecycleOperationID, Reason string
			Freeze                               Freeze
		}
		need(canonical(rawCell(r, "descriptor"), &d) == nil, "closure_encoding", "", "")
		l := c.life[d.LifecycleOperationID]
		need(l.Accepted && l.Request.Operation == "RetireSubject" && l.After != nil && l.After.Retired && d.RootID == c.members[l.After.Run] && d.RootID == textCell(r, "root_id") && textCell(r, "scope_id") == "root:"+d.RootID && d.Reason == "subject_retired", "closure_origin", d.LifecycleOperationID, "")
		closures[d.LifecycleOperationID] = true
		f := c.expected[digest("alan/fact/lifecycle/v1", d.LifecycleOperationID)]
		need(d.Freeze.InvocationID == "root:"+d.RootID && d.Freeze.Trigger == f.ID && d.Freeze.Run == l.After.Run && d.Freeze.Sequence == l.After.Sequence && d.Freeze.Status == l.After.Projection.Status && d.Freeze.DecisionOrder == intCell(r, "decision_order"), "closure_freeze", d.LifecycleOperationID, "")
		for _, fact := range c.out.OrderedFacts {
			if fact.ID == f.ID {
				need(fact.Order == d.Freeze.DecisionOrder, "closure_order", d.LifecycleOperationID, "")
			}
		}
	}
	for id, l := range c.life {
		if l.Accepted && l.Request.Operation == "RetireSubject" && l.After != nil && c.members[l.After.Run] != "" {
			need(closures[id], "missing_retirement_closure", id, "")
		}
	}
}
