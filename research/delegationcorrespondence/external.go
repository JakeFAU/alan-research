package delegationcorrespondence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type boundaryEvent struct {
	AuthorityAdmissionOrdinal                        int64
	Scenario, Token, Store, Boundary, Phase, Subject string
	PID                                              int
	Generation, Epoch, Sequence                      int64
	Intent, Attempt                                  string
	InFlight                                         []string
}
type endpointEvent struct {
	Body          string
	Contact       int
	Event, Method string
	PID           int
	Path, Process string
	Receipt       *Receipt
}

// CheckExternal checks original-level contact and application evidence against
// independently validated raw rows. It does not assert a wire carries a gate
// ordinal; returned ordinal claims come only from local boundary observations.
func CheckExternal(data []byte, files map[string][]byte) (out ExternalReport) {
	out.Valid = true
	out.OwnerEpochs = map[string]map[string]int64{}
	out.Limitations = []string{"Unsigned Store and endpoint exports are contextual evidence, not authenticity proofs", "Endpoint wire has original epoch and no gate ordinal; unique contact-to-gate attribution is not generally provable", "Passive ordinary boundary events may have blank process/case identifiers; process acknowledgments and Open snapshots supply separate context"}
	defer func() {
		if p := recover(); p != nil {
			if f, ok := p.(Failure); ok {
				out.Valid = false
				out.Failure = &f
			} else {
				panic(p)
			}
		}
	}()
	v := CheckRaw(data)
	if !v.Valid {
		out.Valid = false
		out.Failure = &Failure{Edge: "raw_invalid_external_not_evaluated"}
		return
	}
	var raw RawSnapshot
	need(strictJSON(data, &raw) == nil, "external_raw", "", "")
	attempts := map[string]WireRequest{}
	endpoint := map[string]string{}
	byKey := map[string]string{}
	expectedResult := map[string]Receipt{}
	exposures := map[string]*Exposure{}
	add := func(id, url string, w WireRequest) {
		attempts[id] = w
		endpoint[id] = url
		byKey[w.Key] = id
		exposures[id] = &Exposure{Original: id, Endpoint: url, WireSHA256: sha(bytesJSON(w))}
	}
	for _, r := range raw.Tables["external_attempts"] {
		a := decodeCell[ExternalAttempt](r, "descriptor")
		add(a.Intent.ID, a.Intent.Target, a.Request)
	}
	var mappings []struct{ Endpoint, Capability string }
	if d := files["deployment-adapters.json"]; len(d) > 0 {
		need(strictJSON(d, &mappings) == nil, "adapter_deployment", "", "")
	}
	for _, r := range raw.Tables["adapter_attempts"] {
		a := decodeCell[AdapterAttempt](r, "descriptor")
		var wire WireRequest
		need(strictJSON([]byte(a.Prepared.Wire), &wire) == nil, "adapter_wire_encoding", a.Intent.ID, "")
		delta := int64(0)
		switch a.Prepared.Request.Effect {
		case "add-one":
			delta = 1
		case "add-ten":
			delta = 10
		default:
			fail("trusted_adapter_domain", a.Intent.ID, a.Prepared.Request.Effect)
		}
		need(wire.Key == a.Prepared.Request.Key && wire.Profile == "idempotent" && wire.Delta == delta && a.Prepared.Protocol == "phase3b3/fixed-counter/v1", "trusted_adapter_action", a.Intent.ID, "")
		matched := false
		for _, m := range mappings {
			matched = matched || m.Capability == a.Intent.Target && m.Endpoint == a.Prepared.Endpoint
		}
		need(matched, "adapter_deployment", a.Intent.ID, "")
		add(a.Intent.ID, a.Prepared.Endpoint, wire)
	}
	for _, r := range raw.Tables["external_receipts"] {
		expectedResult[textCell(r, "intent_id")] = decodeCell[Receipt](r, "descriptor")
	}
	for _, r := range raw.Tables["adapter_results"] {
		a := decodeCell[AdapterResult](r, "descriptor")
		var receipt Receipt
		need(strictJSON([]byte(a.Candidate.Response), &receipt) == nil, "adapter_response", a.Candidate.ID, "")
		need(a.Candidate.Disposition == "validated" && a.Disposition == "validated" && a.Candidate.Payload["mutation"] == intValue(receipt.Mutation) && a.Candidate.Payload["total"] == intValue(receipt.Total) && a.Values["valid"] == boolValue(true) && a.Values["disposition"] == stringValue("validated") && a.Values["mutation"] == intValue(receipt.Mutation) && a.Values["total"] == intValue(receipt.Total), "adapter_candidate_values", a.Candidate.ID, "")
		expectedResult[a.Candidate.ID] = receipt
	}
	for _, e := range v.Effects {
		for _, g := range e.Gates {
			if g.Allowed {
				exposures[e.ID].AllowedGates = append(exposures[e.ID].AllowedGates, g.Ordinal)
			} else {
				exposures[e.ID].RefusedGates = append(exposures[e.ID].RefusedGates, g.Ordinal)
			}
		}
	}
	committed := map[string]Receipt{}
	endpointContacts := map[string]int{}
	keys := []string{}
	for key := range files {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, name := range []string{"endpoint-a.jsonl", "endpoint-b.jsonl"} {
		d, ok := files[name]
		if !ok {
			continue
		}
		out.Present = true
		readyName := strings.TrimSuffix(name, ".jsonl") + "-ready.json"
		var ready struct {
			URL, Token, Store string
			PID               int
		}
		if rd := files[readyName]; len(rd) > 0 {
			need(strictJSON(rd, &ready) == nil, "endpoint_ready", readyName, "")
		}
		contacts := map[int]WireRequest{}
		for lineNo, line := range bytes.Split(bytes.TrimSpace(d), []byte("\n")) {
			if len(line) == 0 {
				continue
			}
			var event endpointEvent
			need(strictJSON(line, &event) == nil, "endpoint_event", name, fmt.Sprint(lineNo+1))
			out.References = append(out.References, RowRef{name, fmt.Sprint(lineNo + 1), "line", sha(line)})
			if ready.PID > 0 {
				need(event.PID == ready.PID && event.Process == ready.Token, "endpoint_process", name, "")
			}
			switch event.Event {
			case "ingress":
				need(event.Contact == len(contacts)+1 && event.Method == "POST" && event.Path == "/apply", "endpoint_contact_order", name, "")
				var wire WireRequest
				need(strictJSON([]byte(event.Body), &wire) == nil, "endpoint_wire", name, "")
				id := byKey[wire.Key]
				want, exists := attempts[id]
				need(exists && equal(wire, want), "endpoint_original_wire", id, "")
				if ready.URL != "" {
					need(ready.URL == endpoint[id], "endpoint_resource", id, "")
				}
				need(len(exposures[id].AllowedGates) > 0, "contact_without_gate", id, "")
				contacts[event.Contact] = wire
				exposures[id].Contacts++
				out.Contacts++
				endpointContacts[name]++
			case "committed":
				wire, exists := contacts[event.Contact]
				need(exists && event.Receipt != nil, "endpoint_commit_without_contact", name, "")
				r := *event.Receipt
				id := byKey[wire.Key]
				need(r.Key == wire.Key && r.Profile == wire.Profile && r.Delta == wire.Delta && r.Mutation > 0, "endpoint_receipt", id, "")
				if old, ok := committed[id]; ok {
					need(equal(old, r), "idempotent_endpoint_result_changed", id, "")
				} else {
					committed[id] = r
					exposures[id].Applications++
					out.Applications++
				}
			default:
				fail("endpoint_event_domain", name, event.Event)
			}
		}
	}
	need(out.Present || len(v.Effects) == 0, "missing_endpoint_evidence", "", "")
	snapshotName := "endpoints.json"
	if _, ok := files["endpoints-final.json"]; ok {
		snapshotName = "endpoints-final.json"
	}
	if d, ok := files[snapshotName]; ok {
		var endpoints map[string]map[string]json.RawMessage
		need(strictJSON(d, &endpoints) == nil, "endpoint_snapshot", snapshotName, "")
		for label, fields := range endpoints {
			var contacts int
			need(strictJSON(fields["Contacts"], &contacts) == nil && contacts == endpointContacts["endpoint-"+label+".jsonl"], "endpoint_snapshot_contacts", label, "")
			var requests []WireRequest
			var applications []Receipt
			need(strictJSON(fields["Requests"], &requests) == nil && strictJSON(fields["Applications"], &applications) == nil, "endpoint_snapshot_rows", label, "")
			for _, q := range requests {
				need(equal(attempts[byKey[q.Key]], q), "endpoint_snapshot_request", q.Key, "")
			}
			for _, r := range applications {
				need(equal(committed[byKey[r.Key]], r), "endpoint_snapshot_application", r.Key, "")
			}
		}
	}

	if out.Present {
		for id, r := range expectedResult {
			need(equal(committed[id], r), "retained_result_endpoint", id, "")
		}
		for id, x := range exposures {
			need(x.Contacts <= len(x.AllowedGates), "more_contacts_than_opportunities", id, "")
		}
	}
	checkOwnerEvidence(raw, files, &out)
	for _, key := range keys {
		if key == "boundaries.json" || strings.HasSuffix(key, "/boundaries.jsonl") {
			events := []boundaryEvent{}
			if key == "boundaries.json" {
				need(strictJSON(files[key], &events) == nil, "boundary_encoding", key, "")
			} else {
				for _, line := range bytes.Split(bytes.TrimSpace(files[key]), []byte("\n")) {
					if len(line) == 0 {
						continue
					}
					var e boundaryEvent
					need(strictJSON(line, &e) == nil, "boundary_encoding", key, "")
					events = append(events, e)
				}
			}
			for eventIndex, e := range events {
				owner := strings.TrimSuffix(key, "/boundaries.jsonl")
				if key != "boundaries.json" && e.Subject != "" {
					context, ok := out.Owners[owner]
					need(ok, "boundary_owner_missing", key, "")
					r, ok := context.Runs[e.Subject]
					need(ok && r.ID == (RunID{e.Subject, e.Generation}) && r.Epoch == e.Epoch, "boundary_owner_epoch", key, e.Intent)
					need(e.PID == 0 || e.PID == context.PID, "boundary_owner_process", key, e.Intent)
					if e.AuthorityAdmissionOrdinal > 0 {
						out.OwnerBindings = append(out.OwnerBindings, OwnerBoundaryBinding{key, owner, e.Intent, eventIndex + 1, context.PID, r.ID, r.Epoch, e.AuthorityAdmissionOrdinal})
					}
				}
				if e.AuthorityAdmissionOrdinal == 0 {
					continue
				}
				x := exposures[e.Intent]
				need(x != nil, "boundary_original", e.Intent, "")
				found := false
				for _, effect := range v.Effects {
					if effect.ID != e.Intent {
						continue
					}
					for _, g := range effect.Gates {
						if g.Ordinal == e.AuthorityAdmissionOrdinal {
							found = g.Source.Run == (RunID{e.Subject, e.Generation}) && g.Source.Epoch == e.Epoch && g.Source.Sequence == e.Sequence
							if e.Boundary == "external.result-recorded" || e.Boundary == "adapter.result-recorded" {
								need(g.Allowed && effect.Result, "result_boundary_gate", e.Intent, "")
								x.ReturnedOrdinals = append(x.ReturnedOrdinals, g.Ordinal)
							}
						}
					}
				}
				need(found, "boundary_gate_binding", e.Intent, "")
			}
		}
	}
	ids := []string{}
	for id := range exposures {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		out.Exposures = append(out.Exposures, *exposures[id])
	}
	return out
}

// CheckContinuity checks immutable evidence between ordered captured prefixes;
// activations and observation consumption are allowed to advance independently.
func CheckContinuity(phases []PhaseReport, files map[string][]byte) (failure *Failure) {
	defer func() {
		if p := recover(); p != nil {
			if f, ok := p.(Failure); ok {
				failure = &f
			} else {
				panic(p)
			}
		}
	}()
	rank := func(p string) int {
		switch p {
		case "initial/inspect.json", "before-grant/inspect.json":
			return 0
		case "first-owner/post-open/inspect.json":
			return 1
		case "stopped/inspect.json":
			return 2
		case "fresh-owner/post-open/inspect.json", "reopened/inspect.json":
			return 3
		case "final/inspect.json", "after-grant/inspect.json":
			if _, stopped := files["stopped/inspect.json"]; !stopped {
				// Ordinary controls capture final execution before inspecting
				// the reopened Store; process cuts capture final after recovery.
				return 2
			}
			return 4
		}
		return 5
	}
	ps := append([]PhaseReport{}, phases...)
	sort.Slice(ps, func(i, j int) bool { return rank(ps[i].Path) < rank(ps[j].Path) })
	var before *RawSnapshot
	for _, p := range ps {
		if !p.Verdict.Valid {
			continue
		}
		var now RawSnapshot
		need(strictJSON(files[p.Path], &now) == nil, "phase_encoding", p.Path, "")
		if before != nil {
			need(now.Namespace == before.Namespace, "phase_namespace", p.Path, "")
			checkImmutableSource(*before, now, p.Path)
			for _, table := range []string{"intents", "attempts", "lifecycle_results", "external_attempts", "adapter_attempts", "composition_reservations", "composition_charges", "composition_facts", "external_receipts", "adapter_results", "history"} {
				set := map[string]bool{}
				for _, r := range now.Tables[table] {
					set[string(bytesJSON(r))] = true
				}
				for _, r := range before.Tables[table] {
					need(set[string(bytesJSON(r))], "immutable_prefix_changed", p.Path, table)
				}
			}
		}
		before = &now
	}
	return nil
}

// This fixture has fixed programs, subjects, permissions and static roots.
// Dynamic run epochs/projections and observation consumption are not compared
// as immutable rows; their narrower transitions are checked separately.
func checkImmutableSource(before, now RawSnapshot, path string) {
	for _, table := range []string{"metadata", "programs", "subjects", "permissions", "run_origins", "composition_profile", "composition_deployment", "composition_permissions", "composition_roots", "composition_static_members"} {
		set := map[string]bool{}
		for _, r := range before.Tables[table] {
			set[string(bytesJSON(r))] = true
		}
		need(len(now.Tables[table]) == len(before.Tables[table]), "immutable_source_changed", path, table)
		for _, r := range now.Tables[table] {
			need(set[string(bytesJSON(r))], "immutable_source_changed", path, table)
		}
	}
	oldRuns := map[string]Row{}
	for _, r := range before.Tables["runs"] {
		oldRuns[rowKey("runs", r)] = r
	}
	need(len(oldRuns) == len(now.Tables["runs"]), "immutable_source_changed", path, "runs")
	for _, r := range now.Tables["runs"] {
		old, ok := oldRuns[rowKey("runs", r)]
		need(ok, "immutable_source_changed", path, "runs")
		for _, column := range []string{"program_id", "initial"} {
			need(equal(old[column], r[column]), "immutable_source_changed", path, "runs."+column)
		}
	}
}

func ownerSnapshot(files map[string][]byte, path string) RawSnapshot {
	data, ok := files[path]
	need(ok, "owner_snapshot_missing", path, "")
	v := CheckRaw(data)
	need(v.Valid, "owner_snapshot_invalid", path, fmt.Sprint(v.Failure))
	var snapshot RawSnapshot
	need(strictJSON(data, &snapshot) == nil, "owner_snapshot", path, "")
	return snapshot
}

// Open renews exactly the nonterminal retained runs. It does not reduce a
// program, change source authority, or renew an already-terminal activation.
func checkOwnerRenewal(before, after RawSnapshot, path string) {
	need(before.Namespace == after.Namespace, "owner_namespace", path, "")
	checkImmutableSource(before, after, path)
	old := map[string]Row{}
	for _, r := range before.Tables["runs"] {
		old[rowKey("runs", r)] = r
	}
	for _, r := range after.Tables["runs"] {
		previous, ok := old[rowKey("runs", r)]
		need(ok, "owner_renewal", path, "run set")
		want := intCell(previous, "epoch")
		if decodeCell[ProjectionState](previous, "projection").Status == "open" {
			want++
		}
		need(intCell(r, "epoch") == want && intCell(r, "sequence") == intCell(previous, "sequence") && equal(r["projection"], previous["projection"]), "owner_renewal", path, textCell(r, "subject"))
	}
}

func checkOwnerEvidence(raw RawSnapshot, files map[string][]byte, out *ExternalReport) {
	out.Owners = map[string]OwnerEvidence{}
	for _, prefix := range []string{"first-owner", "fresh-owner"} {
		key := prefix + "/open.json"
		data, exists := files[key]
		if !exists {
			continue
		}
		var opened struct {
			Error, Namespace string
			PID              int
			Runs             map[string]OwnerRun
		}
		need(strictJSON(data, &opened) == nil, "owner_encoding", key, "")
		need(opened.Namespace == raw.Namespace, "owner_namespace", key, "")
		need(opened.Error == "<nil>" && opened.PID > 0, "owner_open", key, opened.Error)
		priorPath := "initial/inspect.json"
		if prefix == "fresh-owner" {
			priorPath = "stopped/inspect.json"
		}
		postPath := prefix + "/post-open/inspect.json"
		post := ownerSnapshot(files, postPath)
		_, hasPrior := files[priorPath]
		if hasPrior || prefix == "fresh-owner" {
			checkOwnerRenewal(ownerSnapshot(files, priorPath), post, key)
		} else {
			priorPath = ""
			out.Limitations = append(out.Limitations, prefix+": no captured pre-Open snapshot; initial renewal not checked, owner/post-open activation still bound")
		}
		// The final snapshot may advance execution, but it must retain the same
		// original source/authority as the owner snapshots supplied for this Store.
		checkImmutableSource(post, raw, key)
		need(len(opened.Runs) == len(post.Tables["runs"]), "owner_run_set", key, "")
		revisions := map[string]int64{}
		for _, r := range post.Tables["subjects"] {
			revisions[textCell(r, "id")] = intCell(r, "revision")
		}
		out.OwnerEpochs[prefix] = map[string]int64{}
		for _, row := range post.Tables["runs"] {
			name := textCell(row, "subject")
			r, ok := opened.Runs[name]
			need(ok && r.ID == (RunID{name, intCell(row, "generation")}) && r.Epoch == intCell(row, "epoch") && r.Sequence == intCell(row, "sequence") && r.ProgramID == textCell(row, "program_id") && r.Revision == revisions[name] && equal(r.Projection, decodeCell[ProjectionState](row, "projection")), "empty_reopen_snapshot", key, name)
			out.OwnerEpochs[prefix][name] = r.Epoch
		}
		priorHash := ""
		if hasPrior {
			priorHash = sha(files[priorPath])
		}
		out.Owners[prefix] = OwnerEvidence{key, sha(data), postPath, sha(files[postPath]), priorPath, priorHash, opened.PID, opened.Runs, hasPrior}
	}
	if _, ok := files["reopened/inspect.json"]; ok {
		before := "stopped/inspect.json"
		if _, stopped := files[before]; !stopped {
			before = "final/inspect.json"
		}
		checkOwnerRenewal(ownerSnapshot(files, before), ownerSnapshot(files, "reopened/inspect.json"), "reopened/inspect.json")
	}
}
