package delegationcorrespondence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	m "github.com/JakeFAU/alan-research/research/authorityservice"
	b "github.com/JakeFAU/alan-research/research/mutualobligations"
)

// RevocationCaseReport preserves the prior checker verdicts alongside the
// separately qualified physical evidence and bounded model correspondence.
type RevocationCaseReport struct {
	Legacy              CaseReport
	RawValid, Supported bool
	Physical            ExternalReport
	Config              m.Config
	Producer            RevocationProducer
	Reasons             []string
	Captures            []RevocationCapture
}

type RevocationProducer struct {
	Commit, GoVersion, Test string
	Race                    bool
	SourceSHA256            string
}

type RevocationCapture struct {
	Path, SHA256 string
	Final        bool
	Supported    bool
	Failure      *Failure
	Events       []ProjectedEvent
	Audit        *m.Audit
	State        m.State
	Unresolved   []b.Obligation
	Roots        [2]RevocationRoot
	Principals   [2]RevocationPrincipal
	Controls     []RevocationControl
}

type RevocationRoot struct {
	Concrete    AuthorityRoot
	Composition string
	Model       string
}

type RevocationPrincipal struct {
	Original                                     Original
	RequestAdmissionID, ReplyAdmissionID         string
	RequestObservationID, ReplyObservationID     string
	RequestSelection, ReplySelection             *Fact
	RequestSelectionSource, ReplySelectionSource Source
	SelectedPath                                 *Use
	Authority                                    *Witness
	AuthorizationSource                          Source
	ReservationFact                              *Fact
	Reservation                                  *Reservation
	Debit                                        *Charge
	Gates                                        []Fact
	ResultFact                                   *Fact
	Receipt                                      *Receipt
	Exposure                                     *Exposure
}

type RevocationControl struct {
	Fact      Fact
	Lifecycle LifecycleResult
	Source    Source
	ModelRoot string
	ModelEra  uint8
}

type revocationManifest struct {
	Profile           string
	PrincipalRequests [2]string
	AuxiliaryRequests []string
	Captures          []BindingCapture
	Manifest          struct {
		PrincipalRequests [2]string
		Reciprocal        struct {
			A, B                         RunID
			ABRequest, BARequest         string
			ABService, BAService         string
			ABReply, BAReply, Grant      string
			ABServiceRoot, BAServiceRoot AuthorityRoot
			Claim                        Claim
		}
	}
}

var revocationRuns = [2]RunID{{Subject: "a", Generation: 1}, {Subject: "b", Generation: 1}}
var revocationRootNames = [2]string{"B/g1:c", "B/g1:b"}

const revocationGrantName = "grant:BA:b"

// CheckRevocationDirectory verifies declared artifact bytes without opening
// archived databases, then projects independently checked capture prefixes.
func CheckRevocationDirectory(dir string) (RevocationCaseReport, error) {
	legacy, files, err := checkBindingDirectory(dir)
	out := RevocationCaseReport{Legacy: legacy, Physical: legacy.External, Config: m.GrantRevocationConfig()}
	if err != nil {
		return out, err
	}
	var entry revocationManifest
	var manifestKeys map[string]json.RawMessage
	if err := strictJSON(files["case.json"], &manifestKeys); err != nil {
		return out, fmt.Errorf("case encoding: %w", err)
	}
	if err := json.Unmarshal(files["case.json"], &entry); err != nil {
		return out, err
	}
	var source struct {
		HEAD, GoVersion, Test string
		RaceChildren          bool
	}
	if data, ok := files["source.json"]; ok {
		if err := json.Unmarshal(data, &source); err != nil {
			return out, fmt.Errorf("source encoding: %w", err)
		}
		out.Producer = RevocationProducer{strings.TrimSpace(source.HEAD), source.GoVersion, source.Test, source.RaceChildren, sha(data)}
	}
	out.RawValid = len(legacy.Phases) != 0
	for _, phase := range legacy.Phases {
		out.RawValid = out.RawValid && phase.Verdict.Valid && phase.Verdict.Custody.Declared
	}
	checkRevocationPhysical(&out.Physical, files)
	addReason := func(reason string) { out.Reasons = append(out.Reasons, reason) }
	if !out.RawValid {
		addReason("raw_binding_not_valid")
	}
	if legacy.ContinuityFailure != nil {
		addReason("capture_continuity_failed")
	}
	if !out.Physical.Valid {
		addReason("required_physical_evidence_not_valid")
	}
	if len(legacy.Phases) > 0 {
		previous := ""
		seen := map[string]bool{}
		for i, capture := range entry.Captures {
			if capture.Directory != filepath.Clean(capture.Directory) || !filepath.IsLocal(capture.Directory) ||
				capture.Ordinal != i+1 || capture.Predecessor != previous || seen[capture.Directory] ||
				capture.OwnerPID <= 0 || capture.Label == "" || (i == len(entry.Captures)-1 && capture.Directory != "final") {
				addReason("capture_declaration_outside_domain")
				break
			}
			previous, seen[capture.Directory] = capture.Directory, true
		}
		for _, capture := range entry.Captures {
			if capture.AcknowledgedSIGKILL || capture.OwnerPID != entry.Captures[0].OwnerPID {
				addReason("whole_store_loss")
				break
			}
		}
		if legacy.External.Valid {
			out.Reasons = append(out.Reasons, finalRevocationExposureReasons(legacy.Final, legacy.External.Exposures)...)
		}
	}
	if entry.Profile != bindingVersion || len(entry.AuxiliaryRequests) != 0 || out.Producer.Commit == "" ||
		legacy.Files["case.json"] != legacy.CaseSHA256 {
		addReason("case_declaration_outside_domain")
	}
	if len(out.Reasons) != 0 {
		for i, phase := range legacy.Phases {
			out.Captures = append(out.Captures, RevocationCapture{Path: phase.Path, SHA256: phase.SHA256, Final: i == len(legacy.Phases)-1})
		}
		return out, nil
	}

	var prior []ProjectedEvent
	for i, phase := range legacy.Phases {
		cut := projectRevocationCapture(phase, i == len(legacy.Phases)-1, files[phase.Path], entry, out.Physical.Exposures, out.Config, files)
		if cut.Supported && (len(prior) > len(cut.Events) || !slices.Equal(prior, cut.Events[:len(prior)])) {
			cut.Supported = false
			cut.Failure = &Failure{Edge: "projected_capture_prefix", ID: phase.Path}
		}
		out.Captures = append(out.Captures, cut)
		if !cut.Supported {
			addReason(cut.Failure.Edge)
		}
		prior = cut.Events
	}
	final := out.Captures[len(out.Captures)-1]
	if final.Supported && final.State.Grant != m.RetainedGrant {
		addReason("accepted_separate_grant_missing")
	}
	out.Supported = len(out.Reasons) == 0
	return out, nil
}

func finalRevocationExposureReasons(v RawVerdict, exposures []Exposure) []string {
	var reasons []string
	for _, effect := range v.Effects {
		var x Exposure
		for _, e := range exposures {
			if e.Original == effect.ID {
				x = e
				break
			}
		}
		switch {
		case x.Applications > 0 && !effect.Result:
			reasons = append(reasons, "applied_without_retained_result")
		case x.Contacts > 1 || x.Applications > 1 || len(effect.Gates) > 1:
			reasons = append(reasons, "repeated_physical_opportunity_or_exposure")
		case effect.Result:
			if x.Contacts != 1 || x.Applications != 1 || len(x.AllowedGates) != 1 || len(effect.Gates) != 1 {
				reasons = append(reasons, "result_without_unique_physical_application")
			}
		case effect.Allowed:
			if len(effect.Gates) != 1 || effect.Gates[0].Allowed || x.Contacts != 0 || x.Applications != 0 {
				reasons = append(reasons, "unresolved_final_effect")
			}
		}
	}
	return reasons
}

// CheckExternal accepts optional endpoint artifacts in its legacy scope. This
// relation requires the two ready-bound journals and a complete snapshot.
func checkRevocationPhysical(out *ExternalReport, files map[string][]byte) {
	for _, name := range []string{"endpoint-a.jsonl", "endpoint-b.jsonl", "endpoint-a-ready.json", "endpoint-b-ready.json", "endpoints.json"} {
		if _, ok := files[name]; !ok {
			out.Valid = false
			out.Failure = &Failure{Edge: "required_endpoint_artifact_missing", ID: name}
			return
		}
	}
	if !out.Valid {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			f, ok := p.(Failure)
			if !ok {
				panic(p)
			}
			out.Valid, out.Failure = false, &f
		}
	}()
	var snapshot map[string]map[string]json.RawMessage
	need(strictJSON(files["endpoints.json"], &snapshot) == nil && len(snapshot) == 2,
		"endpoint_snapshot_pair", "endpoints.json", "")
	for _, label := range []string{"a", "b"} {
		name := "endpoint-" + label + ".jsonl"
		ready := readRevocationReady(label, files)
		need(ready.URL != "" && ready.PID > 0 && ready.Token != "", "endpoint_ready_binding", name, "")
		fields := snapshot[label]
		var contacts int
		var requests []WireRequest
		var applications []Receipt
		need(strictJSON(fields["Contacts"], &contacts) == nil && strictJSON(fields["Requests"], &requests) == nil &&
			strictJSON(fields["Applications"], &applications) == nil, "endpoint_snapshot_rows", label, "")
		var observedRequests []WireRequest
		var observedApplications []Receipt
		seen := map[string]bool{}
		committedContacts := map[int]bool{}
		for _, line := range bytes.Split(bytes.TrimSpace(files[name]), []byte("\n")) {
			if len(line) == 0 {
				continue
			}
			var event endpointEvent
			need(strictJSON(line, &event) == nil, "endpoint_event", name, "")
			switch event.Event {
			case "ingress":
				var request WireRequest
				need(strictJSON([]byte(event.Body), &request) == nil, "endpoint_wire", name, "")
				observedRequests = append(observedRequests, request)
			case "committed":
				need(event.Receipt != nil && !committedContacts[event.Contact], "endpoint_commit_coverage", name, "")
				committedContacts[event.Contact] = true
				if !seen[event.Receipt.Key] {
					seen[event.Receipt.Key] = true
					observedApplications = append(observedApplications, *event.Receipt)
				}
			}
		}
		need(contacts == len(observedRequests) && slices.Equal(requests, observedRequests) &&
			slices.Equal(applications, observedApplications), "endpoint_snapshot_coverage", label, "")
	}
}

type revocationReady struct {
	URL, Token, Store string
	PID               int
}

func readRevocationReady(label string, files map[string][]byte) (out revocationReady) {
	name := "endpoint-" + label + "-ready.json"
	need(strictJSON(files[name], &out) == nil, "endpoint_ready_encoding", name, "")
	return out
}

func projectRevocationCapture(phase PhaseReport, final bool, data []byte, entry revocationManifest, exposures []Exposure, config m.Config, files map[string][]byte) (out RevocationCapture) {
	out.Path, out.SHA256, out.Final = phase.Path, phase.SHA256, final
	var order int64
	defer func() {
		if p := recover(); p != nil {
			f, ok := p.(Failure)
			if !ok {
				panic(p)
			}
			f.Order = order
			out.Failure = &f
		}
	}()
	var raw RawSnapshot
	need(strictJSON(data, &raw) == nil, "capture_encoding", phase.Path, "")
	v := phase.Verdict
	rec := entry.Manifest.Reciprocal
	need(rec.A == revocationRuns[b.A] && rec.B == revocationRuns[b.B] &&
		entry.PrincipalRequests == entry.Manifest.PrincipalRequests &&
		entry.PrincipalRequests == [2]string{rec.ABRequest, rec.BARequest},
		"exact_principal_declaration", phase.Path, "")
	rootIDs := revocationCompositionRoots(raw)
	out.Roots = revocationRoots(raw, entry, files, rootIDs[b.B])
	rows := indexRevocationRows(raw)
	effectIDs := [2]string{rec.ABService, rec.BAService}
	replyIDs := [2]string{rec.ABReply, rec.BAReply}
	requestTo, effectTo, replyTo, replyObservation := map[string]b.Obligation{}, map[string]b.Obligation{}, map[string]b.Obligation{}, map[string]b.Obligation{}
	for d := b.AB; d <= b.BA; d++ {
		o := v.Custody.Originals[d]
		p := &out.Principals[d]
		p.Original = o
		need(o.ID == entry.PrincipalRequests[d] && effectIDs[d] != "" && replyIDs[d] != "",
			"exact_principal_declaration", o.ID, "")
		need(o.Committed && o.Caller == revocationRuns[d] && o.Callee == revocationRuns[1-d] &&
			rows.intents[o.ID].Obligation == []string{"ab-request", "ba-request"}[d],
			"exact_principal_original", o.ID, "")
		p.RequestObservationID = identity("observation", "message", o.ID)
		p.RequestAdmissionID = rows.admitted[o.ID]
		need(o.Admitted == (p.RequestAdmissionID != ""), "exact_principal_admission", o.ID, "")
		requestTo[o.ID] = d
		if o.EffectID != "" {
			need(o.EffectID == effectIDs[d], "exact_principal_service", o.EffectID, "")
			effectTo[o.EffectID] = d
		}
		if o.ReplyID != "" {
			need(o.ReplyID == replyIDs[d], "exact_principal_reply", o.ReplyID, "")
			replyTo[o.ReplyID] = d
			p.ReplyObservationID = identity("observation", "message", o.ReplyID)
			p.ReplyAdmissionID = rows.admitted[o.ReplyID]
			replyObservation[p.ReplyObservationID] = d
		}
	}
	for d := b.AB; d <= b.BA; d++ {
		p := &out.Principals[d]
		for _, f := range v.OrderedFacts {
			if f.Kind != "execution" {
				continue
			}
			h := rows.history[historyKey(f.Run, f.Sequence)]
			if h.Selected == p.RequestObservationID && f.Run == p.Original.Callee {
				p.RequestSelection = &f
				p.RequestSelectionSource = Source{h.Run, h.Epoch, h.Sequence, h.Revision}
			}
			if p.ReplyObservationID != "" && h.Selected == p.ReplyObservationID && f.Run == p.Original.Caller {
				p.ReplySelection = &f
				p.ReplySelectionSource = Source{h.Run, h.Epoch, h.Sequence, h.Revision}
			}
		}
		need(p.Original.Selected == (p.RequestSelection != nil) && p.Original.ReplySelected == (p.ReplySelection != nil),
			"exact_principal_selection", p.Original.ID, "")
	}
	effects := map[string]EffectReport{}
	for _, e := range v.Effects {
		_, ok := effectTo[e.ID]
		need(ok && e.Kind == "external", "unmapped_external_attempt", e.ID, "")
		effects[e.ID] = e
	}
	for _, x := range exposures {
		if d, ok := effectTo[x.Original]; ok && (final || effects[x.Original].Result) {
			out.Principals[d].Exposure = &x
		}
	}

	state := m.Initial(config)
	var trace []m.Event
	var grantSequence int64
	var grantID string
	atoms := map[[2]int64]b.Value{}
	emit := func(f Fact, kind string, d b.Obligation, actor b.Agent, value b.Value, mapping string) m.Event {
		var chosen m.Transition
		count := 0
		for _, edge := range m.Successors(config, state) {
			e := edge.Event
			if e.Kind == kind && e.Obligation == d && e.Actor == actor && !e.Duplicate && (value == b.NoValue || e.Result == value) {
				chosen = edge
				count++
			}
		}
		need(count == 1, "unique_forward_edge", f.ID, fmt.Sprintf("%s/%d actor %d has %d edges", kind, d, actor, count))
		e := chosen.Event
		era := uint8(0)
		if state.Revoked {
			era = 1
		}
		need(e.Activation == 0 && e.Revision == era && e.Violations == 0 &&
			(d == b.NoObligation || e.Binding == m.Binding(config, d)), "model_event_descriptor", f.ID, kind)
		state = chosen.State
		trace = append(trace, e)
		out.Events = append(out.Events, ProjectedEvent{e, f.ID, f.OperationID, mapping, f.Order})
		return e
	}
	for _, f := range v.OrderedFacts {
		order = f.Order
		switch f.Kind {
		case "execution":
			h := rows.history[historyKey(f.Run, f.Sequence)]
			for _, i := range h.Intents {
				if d, ok := requestTo[i.ID]; ok {
					eventFact := f
					eventFact.OperationID = i.ID
					emit(eventFact, "commit", d, b.Agent(d), b.NoValue, "retained original commitment")
				}
				if d, ok := replyTo[i.ID]; ok {
					eventFact := f
					eventFact.OperationID = i.ID
					emit(eventFact, "send-result", d, b.Agent(1-d), b.NoValue, "result-bearing reply emission")
				}
			}
			if d, ok := replyObservation[h.Selected]; ok {
				eventFact := f
				eventFact.OperationID = h.Selected
				emit(eventFact, "consume", d, b.Agent(d), b.NoValue, "principal selected the authenticated reply")
			}
		case "authorization":
			if rows.admitted[f.OperationID] == "" {
				continue
			}
			a := rows.auth[f.OperationID]
			need(a.Allowed && a.Source.Run == f.Run && a.Source.Epoch == 1 && a.Source.Revision == 1,
				"admitted_transport_source", f.OperationID, "")
			if d, ok := requestTo[f.OperationID]; ok {
				emit(f, "send-request", d, b.Agent(d), b.NoValue, "accepted admission places send at source authorization")
				emit(f, "admit", d, b.Agent(1-d), b.NoValue, "accepted original admission")
			}
			if d, ok := replyTo[f.OperationID]; ok {
				emit(f, "receive-result", d, b.Agent(d), b.NoValue, "accepted reply admission")
				emit(f, "record", d, b.Agent(d), b.NoValue, "retained reply observation")
			}
		case "lifecycle":
			l := rows.life[f.OperationID]
			a := rows.auth[f.OperationID]
			checkRevocationControl(f, l, a, entry, out.Roots[b.BA], rootIDs[b.A], v.Namespace)
			var e m.Event
			if l.Request.Operation == "DelegateAuthority" {
				need(grantID == "", "duplicate_grant", f.OperationID, "")
				grantID, grantSequence = f.OperationID, l.DecisionSequence
				e = emit(f, "retain-grant", b.BA, b.B, b.NoValue, "accepted request-bound lifecycle issuance")
				want := m.GrantRecord{
					ID: revocationGrantName, Parent: out.Roots[b.BA].Model, Root: out.Roots[b.BA].Model,
					Issuer: b.B, IssuerGeneration: 1, Grantee: b.A, GranteeGeneration: 1,
					Permissions: m.Set(m.BPermission), Binding: m.Binding(config, b.BA),
					NoRedelegation: true, Accepted: true,
				}
				need(e.HasGrant && e.Grant == want, "model_grant_descriptor", f.OperationID, "")
			} else {
				emit(f, "revoke", b.NoObligation, b.B, b.NoValue, "accepted exact-root lifecycle revocation")
			}
			era := uint8(0)
			if state.Revoked {
				era = 1
			}
			out.Controls = append(out.Controls, RevocationControl{f, l, a.Source, out.Roots[b.BA].Model, era})
		case "reservation":
			d, ok := effectTo[f.OperationID]
			need(ok, "unmapped_reservation", f.OperationID, "")
			p := &out.Principals[d]
			a := rows.auth[f.OperationID]
			checkRevocationDecision(f, d, p, rows, effects[f.OperationID], out.Roots[d], rootIDs[1-d], grantID, grantSequence)
			e := emit(f, "authorize", d, b.Agent(1-d), b.NoValue, "first decision at reservation prefix")
			compareRevocationDecision(e, d, p, out.Roots[d])
			p.AuthorizationSource = a.Source
		case "authority_admission":
			d, ok := effectTo[f.OperationID]
			need(ok, "unmapped_gate", f.OperationID, "")
			p := &out.Principals[d]
			g := f.AuthorityAdmission
			need(g != nil && p.Reservation != nil && p.Reservation.Allowed && len(p.Gates) == 0 &&
				g.Source == p.AuthorizationSource && g.AttemptDigest == p.Reservation.AttemptDigest &&
				g.Ordinal == 1 && g.Version == bindingVersion && g.Namespace == v.Namespace,
				"exact_fresh_gate_binding", f.OperationID, "")
			kind := "dispatch-admit"
			if !g.Allowed {
				kind = "dispatch-refuse"
			}
			e := emit(f, kind, d, b.Agent(1-d), b.NoValue, "fresh gate on the immutable selected path")
			need(e.Decision == m.Allowed && e.Authority == revocationModelWitness(d, out.Roots[d], e.Binding.Permission, true) &&
				(g.Allowed && g.Reason == "authorized" && e.Refusal == "" ||
					!g.Allowed && g.Reason == "authority_revoked" && e.Refusal == "selected-path-revoked"),
				"model_fresh_disposition", f.OperationID, "")
			p.Gates = append(p.Gates, f)
		case "external_result":
			d, ok := effectTo[f.OperationID]
			need(ok, "unmapped_external_result", f.OperationID, "")
			p := &out.Principals[d]
			receipt, exists := rows.receipts[f.OperationID]
			need(exists && p.Exposure != nil && p.Exposure.Contacts == 1 && p.Exposure.Applications == 1 &&
				p.Exposure.Endpoint == out.Roots[d].Concrete.Permission.Target && len(p.Gates) == 1 &&
				p.Gates[0].AuthorityAdmission.Allowed && slices.Equal(p.Exposure.AllowedGates, []int64{p.Gates[0].AuthorityAdmission.Ordinal}),
				"atomic_apply_evidence", f.OperationID, "")
			content := [2]int64{receipt.Mutation, receipt.Total}
			value := atoms[content]
			if value == b.NoValue {
				value = b.R1
				if len(atoms) != 0 {
					value = b.R2
				}
				atoms[content] = value
			}
			e := emit(f, "apply", d, b.Agent(1-d), value, "unique contact and application joined to retained result")
			need(e.Result == value && e.Decision == m.Allowed &&
				e.Authority == revocationModelWitness(d, out.Roots[d], e.Binding.Permission, true),
				"model_apply_descriptor", f.OperationID, "")
			p.ResultFact, p.Receipt = &f, &receipt
		}
	}
	out.State = state
	audit, err := m.AuditTrace(config, trace)
	out.Audit = &audit
	need(err == nil && audit.ProtocolValid && audit.AuthoritySafe, "independent_trace_audit", phase.Path, fmt.Sprint(err))
	for d := b.AB; d <= b.BA; d++ {
		p, o := out.Principals[d], state.Base.Obligations[d]
		need(o.Committed == p.Original.Committed && o.Admitted == p.Original.Admitted &&
			o.Applied == (p.Receipt != nil) && o.Consumed == p.Original.ReplySelected && audit.Paid[d] == o.Consumed,
			"model_principal_custody", p.Original.ID, "")
		if state.Authority[d].Permit && !o.Applied {
			out.Unresolved = append(out.Unresolved, d)
		}
	}
	out.Supported = true
	return out
}

func revocationCompositionRoots(raw RawSnapshot) (ids [2]string) {
	need(len(raw.Tables["runs"]) == 2 && len(raw.Tables["subjects"]) == 2 &&
		len(raw.Tables["composition_roots"]) == 2 && len(raw.Tables["composition_static_members"]) == 2,
		"two_static_principals", raw.Namespace, "")
	for _, r := range raw.Tables["runs"] {
		id := RunID{textCell(r, "subject"), intCell(r, "generation")}
		need(slices.Contains(revocationRuns[:], id) && intCell(r, "epoch") == 1, "runtime_activation_outside_domain", id.Subject, "")
	}
	for _, r := range raw.Tables["subjects"] {
		id := RunID{textCell(r, "id"), intCell(r, "generation")}
		need(slices.Contains(revocationRuns[:], id) && intCell(r, "revision") == 1, "runtime_revision_outside_domain", id.Subject, "")
	}
	for _, r := range raw.Tables["history"] {
		need(intCell(r, "epoch") == 1 && intCell(r, "source_revision") == 1,
			"runtime_history_activation_outside_domain", textCell(r, "subject"), "")
	}
	for _, row := range raw.Tables["composition_roots"] {
		r := decodeCell[Root](row, "descriptor")
		i := slices.Index(revocationRuns[:], r.Run)
		need(i >= 0 && slices.Equal(r.Members, []RunID{r.Run}) && r.ID != "", "exact_composition_membership", r.ID, "")
		ids[i] = r.ID
	}
	need(ids[0] != "" && ids[1] != "" && ids[0] != ids[1], "exact_composition_membership", raw.Namespace, "")
	return ids
}

func revocationRoots(raw RawSnapshot, entry revocationManifest, files map[string][]byte, holderRootID string) (roots [2]RevocationRoot) {
	rec := entry.Manifest.Reciprocal
	declared := [2]AuthorityRoot{rec.ABServiceRoot, rec.BAServiceRoot}
	need(declared[0] != declared[1], "aliased_exact_roots", raw.Namespace, "")
	readyA, readyB := readRevocationReady("a", files), readRevocationReady("b", files)
	need(readyA.URL != readyB.URL, "aliased_endpoint_resources", raw.Namespace, "")
	targets := [2]string{readyA.URL, readyB.URL}
	permissions := decodeCell[[]Permission](raw.Tables["composition_permissions"][0], "descriptor")
	external := 0
	for _, p := range permissions {
		if p.Kind == "external" {
			external++
		}
	}
	need(external == 2, "exact_external_permission_count", raw.Namespace, "")
	for d := b.AB; d <= b.BA; d++ {
		p := Permission{Kind: "external", Source: "b", Target: targets[d], Operation: "apply", Revision: 1,
			Obligation: [2]string{"ab-service", "ba-service"}[d]}
		expect := AuthorityRoot{revocationRuns[b.B], revocationRuns[b.B], p}
		need(declared[d] == expect && slices.Contains(permissions, p), "exact_root_mapping", raw.Namespace, p.Obligation)
		roots[d] = RevocationRoot{expect, holderRootID, revocationRootNames[d]}
	}
	return roots
}

type revocationRows struct {
	intents      map[string]Intent
	history      map[string]history
	admitted     map[string]string
	auth         map[string]Authorization
	life         map[string]LifecycleResult
	reservations map[string]Reservation
	attempts     map[string]ExternalAttempt
	charges      map[string]Charge
	receipts     map[string]Receipt
}

func indexRevocationRows(raw RawSnapshot) revocationRows {
	r := revocationRows{
		intents: map[string]Intent{}, history: map[string]history{}, admitted: map[string]string{},
		auth: map[string]Authorization{}, life: map[string]LifecycleResult{}, reservations: map[string]Reservation{},
		attempts: map[string]ExternalAttempt{}, charges: map[string]Charge{}, receipts: map[string]Receipt{},
	}
	for _, row := range raw.Tables["intents"] {
		i := decodeCell[Intent](row, "descriptor")
		r.intents[i.ID] = i
	}
	for _, row := range raw.Tables["history"] {
		run := RunID{textCell(row, "subject"), intCell(row, "generation")}
		sequence := intCell(row, "sequence")
		r.history[historyKey(run, sequence)] = history{Run: run, Sequence: sequence, Selected: nullableText(row, "selected_observation"),
			Epoch: intCell(row, "epoch"), Revision: intCell(row, "source_revision"), Intents: decodeCell[[]Intent](row, "intents")}
	}
	for _, row := range raw.Tables["admissions"] {
		a := decodeCell[Admission](row, "descriptor")
		d := decodeCell[AdmissionResult](row, "result")
		if a.Kind == "message" && d.Accepted {
			i := r.intents[a.IntentID]
			need(i.ID != "" && d.Bound != nil && *d.Bound == (RunID{i.Target, i.Generation}) &&
				d.ObservationID == identity("observation", "message", i.ID) && r.admitted[i.ID] == "",
				"exact_transport_admission", a.IntentID, "")
			r.admitted[i.ID] = textCell(row, "id")
		}
	}
	for _, row := range raw.Tables["attempts"] {
		a := decodeCell[Authorization](row, "result")
		r.auth[a.IntentID] = a
	}
	for _, row := range raw.Tables["lifecycle_results"] {
		l := decodeCell[LifecycleResult](row, "result")
		r.life[l.Request.Intent.ID] = l
	}
	for _, row := range raw.Tables["composition_reservations"] {
		d := decodeCell[Reservation](row, "descriptor")
		r.reservations[d.OperationID] = d
	}
	for _, row := range raw.Tables["external_attempts"] {
		a := decodeCell[ExternalAttempt](row, "descriptor")
		r.attempts[a.Intent.ID] = a
	}
	for _, row := range raw.Tables["composition_charges"] {
		c := decodeCell[Charge](row, "descriptor")
		need(c.Kind == "effect", "non_effect_debit_outside_domain", c.OperationID, "")
		r.charges[c.OperationID] = c
	}
	for _, row := range raw.Tables["external_receipts"] {
		r.receipts[textCell(row, "intent_id")] = decodeCell[Receipt](row, "descriptor")
	}
	return r
}

func checkRevocationControl(f Fact, l LifecycleResult, a Authorization, entry revocationManifest, root RevocationRoot, recipientRootID, namespace string) {
	i := l.Request.Intent
	need(l.Accepted && l.Authority != nil && l.Authority.Version == bindingVersion && l.Authority.Namespace == namespace,
		"accepted_authority_control_required", f.OperationID, l.Reason)
	need(f.Kind == "lifecycle" && i.Run == revocationRuns[b.B] && i.ID == f.OperationID &&
		l.Before != nil && f.Run == l.Before.Run && f.Sequence == l.Before.Sequence &&
		a.Allowed && a.Source == (Source{i.Run, 1, i.Sequence, 1}) &&
		l.Request.SourceRevision == 1 && l.Request.ExpectedRevision == 1 && l.Successor == 0 && equal(l.Before, l.After),
		"exact_authority_lifecycle_source", f.OperationID, "")
	need(l.Authority.Root == root.Concrete && l.Authority.RootID == root.Composition,
		"exact_authority_lifecycle_root", f.OperationID, l.Request.Operation)
	rec := entry.Manifest.Reciprocal
	switch l.Request.Operation {
	case "DelegateAuthority":
		var claim Claim
		want := Claim{bindingVersion, root.Concrete, revocationRuns[b.A], revocationRuns[b.A], 1, rec.BARequest, rec.BAService}
		need(strictJSON([]byte(i.Payload["claim"].Text), &claim) == nil && claim == want && rec.Claim == want &&
			l.Reason == "delegated" && f.OperationID == rec.Grant && i.Target == "a" && i.Generation == 1 &&
			l.Authority.RecipientRootID == recipientRootID && l.DecisionSequence == 1,
			"exact_grant_binding", f.OperationID, "")
	case "RevokeAuthority":
		var revoke Revoke
		need(strictJSON([]byte(i.Payload["claim"].Text), &revoke) == nil && revoke == (Revoke{bindingVersion, root.Concrete}) &&
			l.Reason == "revoked" && i.Target == "b" && i.Generation == 1 && l.Authority.RecipientRootID == "" && l.DecisionSequence == 2,
			"exact_revoke_target", f.OperationID, l.Reason)
	default:
		fail("lifecycle_operation_outside_domain", f.OperationID, l.Request.Operation)
	}
}

func checkRevocationDecision(f Fact, d b.Obligation, p *RevocationPrincipal, rows revocationRows, effect EffectReport,
	root RevocationRoot, servingRootID, grantID string, grantSequence int64) {
	i := rows.intents[f.OperationID]
	a := rows.auth[f.OperationID]
	attempt := rows.attempts[f.OperationID]
	r := rows.reservations[f.OperationID]
	w := effect.Witness
	use := Use{Version: bindingVersion, Path: "own", Root: root.Concrete}
	if d == b.BA {
		need(grantID != "", "service_requires_accepted_grant", f.OperationID, "")
		use.Path, use.GrantID = "delegated", grantID
	}
	need(i.Authority != nil && *i.Authority == use && w != nil && w.Use != nil && *w.Use == use &&
		w.RootID == root.Composition && w.ServingRootID == servingRootID &&
		i.Target == root.Concrete.Permission.Target && i.Obligation == root.Concrete.Permission.Obligation &&
		i.Run == revocationRuns[1-d] && i.Run == f.Run && i.Sequence == f.Sequence,
		"exact_selected_authority_path", f.OperationID, "")
	need(a.Allowed && a.Source == (Source{f.Run, 1, f.Sequence, 1}) &&
		attempt.Request.Epoch == a.Source.Epoch && attempt.Request.Authorization == a.ID &&
		r.OperationID == f.OperationID && r.RootID == servingRootID && r.AttemptDigest != "" &&
		r.Allowed == w.Allowed && r.Allowed == effect.Allowed && r.Reason == w.Reason && r.Reason == effect.Reason,
		"first_decision_or_resource_refusal_outside_domain", f.OperationID, effect.Reason)
	need(p.RequestSelection != nil && p.RequestSelection.Order < f.Order &&
		p.RequestSelection.Sequence == p.Original.RequestSelection && p.RequestSelection.Sequence < i.Sequence &&
		p.RequestAdmissionID != "", "exact_request_selection_before_service", f.OperationID, "")
	if d == b.BA && w.Allowed {
		need(w.OriginalRequestID == p.Original.ID && w.ServiceOperationID == f.OperationID &&
			w.RequestObservationID == p.RequestObservationID && w.RequestConsumedSequence == p.RequestSelection.Sequence &&
			w.GrantDecisionSequence == grantSequence, "exact_delegated_selection_binding", f.OperationID, "")
	}
	c, debited := rows.charges[f.OperationID]
	if r.Allowed {
		want := Charge{RootID: servingRootID, OperationID: f.OperationID, Kind: "effect", Effects: 1, DecisionOrder: f.Order}
		need(r.Reason == "authorized" && debited && c == want && effect.ServingCharge == 1,
			"exact_serving_debit", f.OperationID, "")
		p.Debit = &c
	} else {
		need(d == b.BA && r.Reason == "authority_revoked" && !debited && effect.ServingCharge == 0,
			"first_refusal_outside_domain", f.OperationID, r.Reason)
	}
	p.SelectedPath, p.Authority, p.ReservationFact, p.Reservation = &use, w, &f, &r
}

func revocationModelWitness(d b.Obligation, root RevocationRoot, permission m.Permission, allowed bool) m.AuthorityWitness {
	w := m.AuthorityWitness{Principal: b.Agent(1 - d), Generation: 1, Permission: permission}
	if allowed {
		w.Path, w.Root = m.OwnRoot, root.Model
		if d == b.BA {
			w.Path, w.Grant = m.Delegated, revocationGrantName
		}
	}
	return w
}

func compareRevocationDecision(e m.Event, d b.Obligation, p *RevocationPrincipal, root RevocationRoot) {
	allowed := p.Reservation.Allowed
	need(e.Authority == revocationModelWitness(d, root, e.Binding.Permission, allowed) &&
		(allowed && e.Decision == m.Allowed && e.Refusal == "" ||
			!allowed && e.Decision == m.Refused && e.Refusal == "no-covering-authority"),
		"model_first_decision_disagreement", p.ReservationFact.OperationID, p.Reservation.Reason)
}
