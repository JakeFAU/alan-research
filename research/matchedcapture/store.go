package matchedcapture

import (
	"encoding/json"
	"fmt"
	"strings"

	dc "github.com/JakeFAU/alan-research/research/delegationcorrespondence"
	m "github.com/JakeFAU/alan-research/research/effectobservation"
)

func checkStore(f map[string][]byte, branch string) dc.MatchedBindingReport {
	local := map[string][]byte{"initial/inspect.json": f["seed/initial/inspect.json"]}
	for p, b := range f {
		if strings.HasPrefix(p, branch+"/") {
			local[strings.TrimPrefix(p, branch+"/")] = b
		}
	}
	for _, label := range []string{"a", "b"} {
		local["endpoint-"+label+"-ready.json"] = f[branch+"/endpoint-"+label+"/ready.json"]
	}
	for _, cut := range []string{"initial", "first-owner/post-open", "stopped", "reopened", "final"} {
		p := cut + "/inspect.json"
		outer := read[map[string]json.RawMessage](local, p)
		for _, key := range []string{"Namespace", "KeyVersion", "Settings", "Roots", "Tables"} {
			_, ok := outer[key]
			require(ok, p, "missing complete export field "+key)
		}
		// Unchanged raw validator deliberately fails closed on unsupported schema,
		// including any newly introduced ReceiptEvidence. Comparison still retains it.
		raw := decode[dc.RawSnapshot](local[p], p)
		checkOuter(raw, cut, p)
		epoch, a, b := int64(2), int64(4), int64(5)
		switch cut {
		case "initial":
			epoch, a, b = 1, 0, 0
		case "first-owner/post-open":
			a, b = 0, 0
		case "reopened":
			epoch = 3
		case "final":
			epoch, b = 3, 8
		}
		require(len(raw.Tables["runs"]) == 2, p, "two principals required")
		seen := map[string]bool{}
		for _, r := range raw.Tables["runs"] {
			s := decode[string](r["subject"], p)
			want := a
			if s == "b" {
				want = b
			}
			require((s == "a" || s == "b") && !seen[s] && decode[int64](r["generation"], p) == 1 && decode[int64](r["epoch"], p) == epoch && decode[int64](r["sequence"], p) == want, p, "principal coordinates")
			seen[s] = true
		}
	}
	exactRenewal(local["initial/inspect.json"], local["first-owner/post-open/inspect.json"], 1, branch+"/first-owner/post-open")
	exactRenewal(local["stopped/inspect.json"], local["reopened/inspect.json"], 2, branch+"/reopened")
	report := dc.CheckMatchedBinding(local, f["seed/manifest.json"])
	require(report.Failure == nil, branch, "raw binding: "+fmt.Sprint(report.Failure))
	require(len(report.Bindings) == 3, branch, "missing selected cut bindings")
	run := read[map[string]json.RawMessage](local, "run.json")
	require(decode[string](run["Error"], "run.json.Error") == "<nil>", branch+"/run.json", "run failure")
	body := decode[map[string]json.RawMessage](run["Report"], "run.json.Report")
	unknown := decode[[]string](body["UnknownExternal"], "run.json.Report.UnknownExternal")
	require(same(unknown, []string{service}), branch+"/run.json", "unknown external identity")
	return report
}
func exactRenewal(a, b []byte, epoch int64, path string) {
	x, e := parseJSON(a)
	require(e == nil, path, fmt.Sprint(e))
	before, ok := x.(map[string]any)
	require(ok, path, "outer object")
	tables, ok := before["Tables"].(map[string]any)
	require(ok, path, "tables")
	runs, ok := tables["runs"].([]any)
	require(ok && len(runs) == 2, path, "runs")
	for _, v := range runs {
		r, ok := v.(map[string]any)
		require(ok, path, "run row")
		require(r["epoch"] == json.Number(fmt.Sprint(epoch)), path, "original activation")
		r["epoch"] = json.Number(fmt.Sprint(epoch + 1))
	}
	c, e := CompareJSON(jsonBytes(before), b)
	require(e == nil && c.StructuralEqual, path, "non-exact complete renewal at "+strings.Join(c.Paths, ", "))
}
func selectedRequest(f map[string][]byte, branch string) dc.WireRequest {
	raw := read[dc.RawSnapshot](f, branch+"/stopped/inspect.json")
	var found []dc.WireRequest
	for _, r := range raw.Tables["external_attempts"] {
		if decode[string](r["intent_id"], "external_attempts") == service {
			a := decode[dc.ExternalAttempt]([]byte(decode[string](r["descriptor"], "external_attempts")), "external_attempts.descriptor")
			require(a.Allowed && a.Intent.ID == service && a.AuthorizationID == a.Request.Authorization, branch, "selected allowed request")
			found = append(found, a.Request)
		}
	}
	require(len(found) == 1, branch, "one selected request")
	q := found[0]
	require(q.Subject == "a" && q.Generation == 1 && q.Epoch == 2 && q.Sequence == 4 && q.Delta == 1 && q.Profile == "idempotent" && q.Key != "" && q.Authorization != "", branch, "selected wire coordinates")
	return q
}
func checkOpen(f map[string][]byte, p, cut string, id event) {
	o := read[opened](f, p)
	r := read[dc.RawSnapshot](f, cut+"/inspect.json")
	require(o.PID == id.PID && o.Scenario == id.Scenario && o.Token == id.Token && o.Store == id.Store && o.Error == "<nil>" && o.Namespace == r.Namespace && len(o.Runs) == 2, p, "process/open/capture join")
	for _, row := range r.Tables["runs"] {
		subject := decode[string](row["subject"], p)
		v, ok := o.Runs[subject]
		require(ok, p, "missing principal")
		projection := decode[dc.ProjectionState]([]byte(decode[string](row["projection"], p)), p)
		require(v.ID == (dc.RunID{Subject: subject, Generation: 1}) && v.ProgramID == decode[string](row["program_id"], p) && v.Epoch == decode[int64](row["epoch"], p) && v.Sequence == decode[int64](row["sequence"], p) && v.Revision == 1 && same(v.Projection, projection), p, "run/capture join")
	}
}

type ModelCut struct {
	Path   string
	Store  m.StoreView
	Ledger m.Ledger
	Events []m.Event
}

func mapModels(branch string, b dc.MatchedBindingReport) []ModelCut {
	require(len(b.Bindings) == 3, branch, "raw binding prerequisite")
	out := []ModelCut{}
	events := []m.Event{}
	binding := m.Binding{Original: "BA/1", SourceGeneration: 1, SourceRevision: 1, Grant: "grant:BA:b", Root: m.RootB, OriginalActivation: 1, Wire: "wire:BA:1"}
	add := func(kind m.EventKind, actor m.Actor, epoch uint8, root m.Root, value m.Result) {
		events = append(events, m.Event{Kind: kind, Binding: binding, Actor: actor, Incarnation: epoch, Root: root, Value: value})
	}
	if branch == "applied" {
		add(m.Launch, m.ActorA, 1, "", "")
		add(m.Ingress, m.ActorEndpoint, 0, "", "")
		add(m.Apply, m.ActorEndpoint, 0, "", m.Value1)
	}
	add(m.LoseOwner, m.ActorOwner, 1, "", "")
	for i, c := range b.Bindings {
		cut := []string{"stopped", "reopened", "final"}[i]
		owner, epoch := m.OwnerStopped, uint8(1)
		if i > 0 {
			owner, epoch = m.OwnerReopened, 2
		}
		if i == 1 {
			add(m.Reopen, m.ActorOwner, 2, "", "")
		}
		if i == 2 {
			add(m.RevokeB, m.ActorB, 2, m.RootB, "")
			add(m.FreshRefuse, m.ActorA, 2, m.RootB, "")
		}
		view := m.StoreView{Binding: binding, AcceptedGrant: c.Grant.Accepted, FirstAllow: c.Authorization.Allowed, Debit: uint8(c.Debit.Effects), FirstGateAllowed: c.Gates[0].AuthorityAdmission.Allowed, Owner: owner, Activation: epoch, BRevoked: c.Revoke != nil, ControlRoot: m.RootC, FreshRefused: len(c.Gates) == 2}
		ledger, e := m.AuditTrace(events)
		require(e == nil && view == ledger.StoreView(), branch+"/"+cut, "independent ledger/projection disagreement")
		out = append(out, ModelCut{branch + "/" + cut, view, ledger, append([]m.Event{}, events...)})
	}
	return out
}

func checkSeedProvision(f map[string][]byte) {
	path := "seed/provision.json"
	checkSeedManifest(f)
	p := read[struct {
		Composition struct {
			RequestAuthority string
			Roots            []dc.Root
			Templates        json.RawMessage
		}
		Assets   json.RawMessage
		Programs []struct {
			ID, Semantics string
			Bytes         []byte
		}
		Subjects []dc.Subject
		Runs     []struct {
			ID         dc.RunID
			ProgramID  string
			Projection dc.ProjectionState
		}
		Permissions []dc.Permission
	}](f, path)
	raw := read[dc.RawSnapshot](f, "seed/initial/inspect.json")
	require(p.Composition.RequestAuthority == "alan/request-authority/v2" && len(p.Composition.Roots) == 2 && string(p.Composition.Templates) == "null" && string(p.Assets) == "null", path, "seed composition domain")
	require(len(p.Programs) == 2 && len(raw.Tables["programs"]) == 2 && len(p.Runs) == 2 && len(p.Subjects) == 2, path, "seed principal/program set")
	programSeen := map[string]bool{}
	for _, program := range p.Programs {
		require(!programSeen[program.ID], path, "duplicate seed program")
		programSeen[program.ID] = true
		found := false
		for _, row := range raw.Tables["programs"] {
			if decode[string](row["id"], path) == program.ID {
				require(decode[string](row["semantics"], path) == program.Semantics && decode[string](row["bytes"], path) == string(program.Bytes), path, "seed program bytes")
				found = true
			}
		}
		require(found, path, "seed program identity")
	}
	runSeen := map[dc.RunID]bool{}
	for _, r := range p.Runs {
		require(!runSeen[r.ID], path, "duplicate seed run")
		runSeen[r.ID] = true
		found := false
		for _, row := range raw.Tables["runs"] {
			if decode[string](row["subject"], path) == r.ID.Subject {
				require(r.ID.Generation == 1 && r.ProgramID == decode[string](row["program_id"], path) && same(r.Projection, decode[dc.ProjectionState]([]byte(decode[string](row["initial"], path)), path)), path, "seed run")
				found = true
			}
		}
		require(found, path, "seed run identity")
	}
	subjectSeen := map[string]bool{}
	for _, s := range p.Subjects {
		require(!subjectSeen[s.ID], path, "duplicate seed subject")
		subjectSeen[s.ID] = true
		found := false
		for _, row := range raw.Tables["subjects"] {
			if decode[string](row["id"], path) == s.ID {
				require(s.Generation == decode[int64](row["generation"], path) && s.Owner == decode[string](row["owner"], path) && s.Revision == decode[int64](row["revision"], path) && !s.Retired && !s.RetryCancelled, path, "seed subject")
				found = true
			}
		}
		require(found, path, "seed subject identity")
	}
	require(len(p.Permissions) == len(raw.Tables["permissions"]), path, "seed permission count")
	permissionSeen := map[dc.Permission]bool{}
	for _, p := range p.Permissions {
		require(!permissionSeen[p], path, "duplicate seed permission")
		permissionSeen[p] = true
		found := false
		for _, row := range raw.Tables["permissions"] {
			v := dc.Permission{Kind: decode[string](row["kind"], path), Source: decode[string](row["source"], path), Target: decode[string](row["target"], path), Operation: decode[string](row["operation"], path), Revision: decode[int64](row["revision"], path), Obligation: decode[string](row["obligation"], path)}
			found = found || v == p
		}
		require(found, path, "seed permission identity")
	}
	rootSeen := map[dc.RunID]bool{}
	for _, root := range p.Composition.Roots {
		require(!rootSeen[root.Run], path, "duplicate seed root")
		rootSeen[root.Run] = true
		found := false
		for _, row := range raw.Tables["composition_roots"] {
			v := decode[dc.Root]([]byte(decode[string](row["descriptor"], path)), path)
			if v.Run == root.Run {
				root.ID = v.ID
				root.Version = v.Version
				require(same(root, v), path, "seed root")
				found = true
			}
		}
		require(found, path, "seed root identity")
	}
}

// Check the public summary schema and its joins to raw rows. The old raw
// authority validator deliberately does not interpret these outer summaries.
func checkOuter(raw dc.RawSnapshot, cut, path string) {
	require(raw.Namespace != "" && raw.KeyVersion == 2, path, "complete namespace/key identity")
	settings := decode[struct {
		JournalMode, SQLiteVersion, SQLiteSourceID         string
		Synchronous, ForeignKeys, BusyTimeoutMS, QueryOnly int
	}](raw.Settings, path+".Settings")
	require(settings.JournalMode == "wal" && settings.SQLiteVersion != "" && settings.SQLiteSourceID != "" && settings.Synchronous == 2 && settings.ForeignKeys == 1 && settings.BusyTimeoutMS == 5000 && settings.QueryOnly == 1, path, "inspection settings")
	type rootView struct {
		RootID                                                  string
		Run                                                     dc.RunID
		LocalStatus                                             string
		Limits                                                  dc.Limits
		Charges                                                 dc.Charge
		Outcome                                                 json.RawMessage
		Pending                                                 string
		PendingInvocations, UnresolvedExternal, UnknownExternal []string
	}
	roots := decode[[]rootView](raw.Roots, path+".Roots")
	require(len(roots) == 2 && len(raw.Tables["composition_roots"]) == 2, path, "complete root summaries")
	seen := map[string]bool{}
	for _, r := range roots {
		require(!seen[r.RootID] && r.RootID != "" && r.LocalStatus == "open" && r.Pending == "members" && string(r.Outcome) == "null" && r.PendingInvocations != nil && len(r.PendingInvocations) == 0, path, "root summary domain")
		seen[r.RootID] = true
		found := false
		for _, row := range raw.Tables["composition_roots"] {
			v := decode[dc.Root]([]byte(decode[string](row["descriptor"], path)), path)
			if v.ID == r.RootID {
				require(v.Run == r.Run && v.Limits == r.Limits, path, "root descriptor join")
				found = true
			}
		}
		require(found, path, "root identity")
		effects := int64(0)
		for _, row := range raw.Tables["composition_charges"] {
			c := decode[dc.Charge]([]byte(decode[string](row["descriptor"], path)), path)
			if c.RootID == r.RootID {
				effects += c.Effects
			}
		}
		require(r.Charges == (dc.Charge{RootID: r.RootID, Effects: effects}), path, "root aggregate charges")
		want := []string{}
		if cut != "initial" && cut != "first-owner/post-open" && r.Run.Subject == "a" {
			want = append(want, service)
		}
		require(same(r.UnresolvedExternal, want) && same(r.UnknownExternal, want), path, "root unresolved external custody")
	}
}

func checkSeedManifest(f map[string][]byte) {
	path := "seed/manifest.json"
	m := read[struct {
		Version string
		Options struct {
			EndpointA, EndpointB, Scenario string
			TerminalHandshake              bool
		}
		A, B                                                                       dc.RunID
		ABRequest, BARequest, Grant, Ready, BAService, BAReply, ABService, ABReply string
		BAServiceRoot, ABServiceRoot                                               dc.AuthorityRoot
		Claim                                                                      dc.Claim
	}](f, path)
	require(m.Version == "alan/reciprocal-fixture/v1" && m.Options.Scenario == "s2-handshake" && m.Options.TerminalHandshake, path, "frozen reciprocal fixture")
	require(m.A == (dc.RunID{Subject: "a", Generation: 1}) && m.B == (dc.RunID{Subject: "b", Generation: 1}) && m.ABRequest == `["a",1,1,0]` && m.BARequest == `["b",1,1,0]` && m.Grant == `["b",1,2,0]` && m.Ready == `["b",1,5,0]` && m.BAService == service && m.BAReply == `["a",1,7,0]` && m.ABService == `["b",1,12,0]` && m.ABReply == `["b",1,15,0]`, path, "frozen retained instruction locators")
	for _, branch := range []string{"zero", "applied"} {
		a := read[ready](f, branch+"/endpoint-a/ready.json")
		b := read[ready](f, branch+"/endpoint-b/ready.json")
		require(a.URL == m.Options.EndpointA && b.URL == m.Options.EndpointB && m.ABServiceRoot.Permission.Target == a.URL && m.BAServiceRoot.Permission.Target == b.URL, path, "manifest endpoint joins")
	}
}
