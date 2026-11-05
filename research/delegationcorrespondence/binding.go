package delegationcorrespondence

// bindOriginal derives the selection from the immutable grant's original ID.
// Admission has no global timestamp: its authentic observation and source
// authorization supply a causal linearization, not measured transport time.
func (c *checker) bindOriginal(i Intent, d Claim, w *Witness) bool {
	q := c.intents[d.OriginalRequestID]
	need(q.ID != "" && q.Kind == "message" && q.Run == d.Root.Run && q.Target == i.Run.Subject && q.Generation == i.Run.Generation && c.revisions[q.ID] == d.Root.Permission.Revision, "binding_original", i.ID, q.ID)
	oid := identity("observation", "message", q.ID)
	o, exists := c.observations[oid]
	if !exists || o.Consumed == 0 || o.Consumed >= i.Sequence {
		return false
	}
	need(o.Run == i.Run && o.Provenance == "ordinary-message" && o.Source == q.Run.Subject && o.SourceGeneration == q.Run.Generation, "binding_observation", i.ID, oid)
	h := c.history[historyKey(i.Run, o.Consumed)]
	need(h.Selected == oid && h.Revision == c.revisions[i.ID], "binding_selection", i.ID, oid)
	need(c.seq[i.Run] >= o.Consumed && c.has("authorization", q.ID) && c.auth[q.ID].Allowed, "binding_prefix", i.ID, oid)
	admitted := false
	for id, a := range c.admissions {
		r := c.admitResults[id]
		admitted = admitted || a.Kind == "message" && a.IntentID == q.ID && r.Accepted && r.Bound != nil && *r.Bound == i.Run && r.ObservationID == oid
	}
	need(admitted, "binding_admission", i.ID, oid)
	w.RequestObservationID, w.RequestConsumedSequence = oid, o.Consumed
	return true
}

// BindingChain names an expected causal association, not an authority fact.
// Missing future operations remain missing in a valid incomplete prefix.
type BindingChain struct{ OriginalID, ServiceID, ReplyID string }
type BindingSpec struct {
	Principals []BindingChain
	Auxiliary  []BindingChain
}

func (c *checker) bindingCustody(spec BindingSpec) {
	need(len(spec.Principals) == 2, "principal_declaration", "", "exactly two principal originals required")
	c.out.Custody.Declared = true
	seen := map[string]bool{}
	for group, chains := range [][]BindingChain{spec.Principals, spec.Auxiliary} {
		for _, chain := range chains {
			need(chain.OriginalID != "" && !seen[chain.OriginalID], "principal_declaration", chain.OriginalID, "duplicate/empty original")
			seen[chain.OriginalID] = true
			q, exists := c.intents[chain.OriginalID]
			o := Original{ID: chain.OriginalID, State: "not-committed"}
			if exists {
				need(q.Kind == "message", "custody_request_kind", q.ID, "")
				o.Caller, o.Callee, o.Committed, o.State = q.Run, RunID{q.Target, q.Generation}, true, "committed-unpaid"
				oid := identity("observation", "message", q.ID)
				obs := c.observations[oid]
				o.Admitted, o.Selected, o.RequestSelection = obs.ID != "", obs.Consumed > 0, obs.Consumed
				eid := chain.ServiceID
				if eid == "" { // Auxiliary services derive only from exact accepted grant originals.
					for id, e := range c.effects {
						if e.Witness.OriginalRequestID == q.ID {
							need(eid == "", "ambiguous_auxiliary_service", q.ID, "")
							eid = id
						}
					}
				}
				if e, ok := c.effects[eid]; ok {
					need(e.Intent.Run == o.Callee, "custody_service_run", eid, q.ID)
					if e.Intent.Authority != nil && e.Intent.Authority.Path == "delegated" && e.Allowed {
						need(e.Witness.OriginalRequestID == q.ID, "custody_service_original", eid, q.ID)
					}
					if e.Allowed && e.Intent.Authority != nil && e.Intent.Authority.Path == "own" {
						need(o.Selected && o.RequestSelection < e.Intent.Sequence, "custody_request_selection", eid, q.ID)
					}
					o.EffectID = eid
					o.State = "first-refused"
					if e.Allowed {
						o.State = "unknown-or-not-contacted"
					}
					for id, v := range c.observations {
						if v.Operation == eid && (v.Provenance == "resource-receipt" || v.Provenance == "runtime-adapter") {
							o.ResultObservation, o.ResultSelected = id, v.Consumed > 0
							if e.Allowed {
								o.State = "retained-result"
							}
						}
					}
				}
				for rid, r := range c.intents {
					if r.Kind != "message" || r.Payload["request_id"] != stringValue(q.ID) || r.Payload["service_operation"].Kind != "string" {
						continue
					}
					need(chain.ReplyID == "" || chain.ReplyID == rid, "custody_reply_locator", rid, chain.ReplyID)
					need(o.ReplyID == "" && r.Run == o.Callee && r.Target == o.Caller.Subject && r.Generation == o.Caller.Generation, "reply_original", rid, q.ID)
					e, ok := c.effects[r.Payload["service_operation"].Text]
					need(ok && e.Intent.ID == o.EffectID && e.Allowed && o.ResultSelected && c.observations[o.ResultObservation].Consumed < r.Sequence, "reply_result_custody", rid, q.ID)
					receipt := c.receipts[e.Intent.ID]
					need(receipt.Key != "" && r.Payload["accepted"] == boolValue(true) && r.Payload["mutation"] == intValue(receipt.Mutation) && r.Payload["total"] == intValue(receipt.Total), "reply_result_values", rid, q.ID)
					o.ReplyID = rid
					reply := c.observations[identity("observation", "message", rid)]
					o.ReplySelected, o.ReplySelection = reply.Consumed > 0, reply.Consumed
					if o.ReplySelected {
						need(reply.Run == o.Caller, "caller_generation", rid, "")
						o.State = "consumed"
					}
				}
			}
			if group == 0 {
				c.out.Custody.Originals = append(c.out.Custody.Originals, o)
			} else {
				c.out.Custody.Auxiliary = append(c.out.Custody.Auxiliary, o)
			}
		}
	}
	c.out.Custody.Complete = true
	for _, o := range c.out.Custody.Originals {
		c.out.Custody.Complete = c.out.Custody.Complete && o.ReplySelected
	}
	if len(c.out.Effects) > 0 {
		first := c.out.Effects[0].Order
		both := true
		for _, o := range c.out.Custody.Originals {
			q := c.intents[o.ID]
			found := false
			for _, f := range c.out.OrderedFacts {
				if f.Kind == "execution" && f.Run == q.Run && f.Sequence == q.Sequence {
					found = f.Order < first
				}
			}
			both = both && found
			if o.Caller.Subject == "b" && o.ReplySelected {
				for _, other := range c.out.Custody.Originals {
					if other.Callee == o.Caller && other.EffectID != "" {
						c.out.Custody.BSerialGuard = o.ReplySelection < c.effects[other.EffectID].Intent.Sequence
					}
				}
			}
		}
		c.out.Custody.BothOutstandingAtFirstService = both
	}
}

// BindingCase is the explicit producer handoff. Extra manifest detail is
// retained by the caller; only declared locators contribute custody labels.
type BindingCase struct {
	Version, Name, Schedule, Profile, Classification string
	PrincipalRequests                                [2]string
	AuxiliaryRequests                                []string
	Manifest                                         struct {
		Reciprocal struct{ ABRequest, BARequest, ABService, BAService, ABReply, BAReply string }
	}
	Captures []BindingCapture
}
type BindingCapture struct {
	Ordinal                       int
	Predecessor, Label, Directory string
	OwnerPID                      int
	AcknowledgedSIGKILL           bool
}

func (c BindingCase) Spec() BindingSpec {
	r := c.Manifest.Reciprocal
	s := BindingSpec{Principals: []BindingChain{{c.PrincipalRequests[0], r.ABService, r.ABReply}, {c.PrincipalRequests[1], r.BAService, r.BAReply}}}
	for _, id := range c.AuxiliaryRequests {
		s.Auxiliary = append(s.Auxiliary, BindingChain{OriginalID: id})
	}
	return s
}
