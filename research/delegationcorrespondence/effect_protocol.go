package delegationcorrespondence

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"slices"
	"strings"
)

const effectLossName = "applied-unknown-revoked-before-fresh-gate"

// These are external, unsigned historical evidence anchors, not supplied by the
// edited dossier. Open has no Store/token fields; those identities come from the
// verified sealed archive. Semantic ack fields are deliberately not hash-gated.
//
//go:embed testdata/effect-observation/anchors.json
var effectAnchorBytes []byte

type effectAnchors struct {
	Commit, Tree string
	SourceFiles  map[string]string
	Cases        []effectCaseAnchor
}
type effectCaseAnchor struct {
	Name                                                                              string
	Race                                                                              bool
	ManifestSHA256, DigestSHA256, ArchiveSHA256, MemberPrefix, CanonicalMembersSHA256 string
	SourceHashes, BindingHashes                                                       map[string]string
	RuntimeIdentity                                                                   boundaryEvent
}

func checkEffectProducer(files map[string][]byte, legacy CaseReport) EffectObservationProducer {
	var anchors effectAnchors
	need(strictJSON(effectAnchorBytes, &anchors) == nil, "effect_builtin_anchor", "", "")
	source := effectJSON[struct {
		HEAD, GoVersion, Test, GOOS, GOARCH string
		RaceChildren                        bool
		FilesSHA256                         map[string]string
	}](files, "source.json")
	wantTest := "TestPhase3B4CampaignRevocation/N07"
	if legacy.Name == effectLossName {
		wantTest = "TestPhase3B4CampaignAppliedUnknownRevoked"
	}
	need((legacy.Name == "N07" || legacy.Name == effectLossName) && strings.TrimSpace(source.HEAD) == anchors.Commit && source.Test == wantTest &&
		source.GoVersion == "go1.27.3" && source.GOOS == "darwin" && source.GOARCH == "arm64", "effect_producer_identity", "source.json", "")
	var anchor effectCaseAnchor
	for _, a := range anchors.Cases {
		if a.Name == legacy.Name && a.Race == source.RaceChildren {
			anchor = a
		}
	}
	need(anchor.Name != "" && equal(source.FilesSHA256, anchor.SourceHashes), "effect_producer_source", "source.json", "")
	binding := effectJSON[map[string]string](files, "binding-source.json")
	need(equal(binding, anchor.BindingHashes), "effect_producer_source", "binding-source.json", "")
	need(sha(files["manifest.json"]) == anchor.ManifestSHA256, "effect_sealed_manifest", "manifest.json", "")
	manifest := effectJSON[struct{ ProgramIDs map[string]string }](files, "manifest.json")
	need(len(manifest.ProgramIDs) == 2, "effect_manifest_programs", "manifest.json", "")
	for _, phase := range legacy.Phases {
		raw := effectRaw(files, phase.Path)
		for _, row := range raw.Tables["runs"] {
			subject := textCell(row, "subject")
			need(manifest.ProgramIDs[subject] != "" && textCell(row, "program_id") == manifest.ProgramIDs[subject],
				"effect_manifest_programs", phase.Path, subject)
		}
	}

	entry := effectJSON[struct{ Manifest json.RawMessage }](files, "case.json")
	need(bytes.Equal(effectCanonicalJSON(entry.Manifest), effectCanonicalJSON(files["manifest.json"])), "effect_manifest_join", "case.json", "")
	// This is informational membership, separate from semantic acceptance of
	// internally rehashed attacks. Binary artifacts were already byte-hashed.
	return EffectObservationProducer{Commit: anchors.Commit, Tree: anchors.Tree, Test: source.Test, GoVersion: source.GoVersion,
		Race: source.RaceChildren, SourceSHA256: sha(files["source.json"]), BindingSourceSHA256: sha(files["binding-source.json"]),
		ArchiveSHA256: anchor.ArchiveSHA256, MemberPrefix: anchor.MemberPrefix, SealedDigestSHA256: anchor.DigestSHA256,
		SealedBytesMatch: sha(bytesJSON(legacy.Files)) == anchor.CanonicalMembersSHA256, SourceFiles: anchors.SourceFiles, Trust: "Pinned commit/tree/source and two unsigned sealed archives containing four selected dossiers are external assumptions. Runtime token/Store anchor is not an Open field or general process authentication."}
}

func checkEffectCaptures(files map[string][]byte, entry revocationManifest, legacy CaseReport, loss bool) {
	c := effectJSON[BindingCase](files, "case.json")
	dirs := []string{"before-revoke", "after-revoke", "final"}
	schedule := "controlled-public-revocation-at-boundary"
	if loss {
		dirs = []string{"initial", "first-owner/post-open", "stopped", "reopened", "final"}
		schedule = "ancillary-v2-process-cut-then-controlled-retained-protocol"
	}
	need(c.Profile == bindingVersion && c.Schedule == schedule && len(c.AuxiliaryRequests) == 0 && len(c.Captures) == len(dirs) &&
		legacy.Files["case.json"] == legacy.CaseSHA256, "effect_protocol_declaration", "case.json", "")
	previous := ""
	for i, d := range dirs {
		cap := c.Captures[i]
		need(cap.Directory == d && cap.Label == d && cap.Ordinal == i+1 && cap.Predecessor == previous && cap.OwnerPID > 0 &&
			cap.AcknowledgedSIGKILL == (loss && i == 2), "effect_capture_order", d, "")
		previous = d
		if !loss {
			need(cap.OwnerPID == c.Captures[0].OwnerPID, "effect_owner_identity", d, "")
		}
	}
	losses := 0
	for name := range files {
		if name == "kill.json" || (strings.HasPrefix(name, "loss-") && strings.HasSuffix(name, ".json")) {
			losses++
		}
	}
	wantLosses := 0
	if loss {
		wantLosses = 1
	}
	need(losses == wantLosses && legacy.External.AcknowledgedLosses == wantLosses, "effect_single_owner_loss", "case.json", "")
	if loss {
		caps := c.Captures
		need(caps[0].OwnerPID == caps[3].OwnerPID && caps[3].OwnerPID == caps[4].OwnerPID && caps[1].OwnerPID == caps[2].OwnerPID &&
			caps[0].OwnerPID != caps[1].OwnerPID, "effect_owner_identity", "case.json", "")
	}
	rec := entry.Manifest.Reciprocal
	need(rec.A == revocationRuns[0] && rec.B == revocationRuns[1] && entry.PrincipalRequests == entry.Manifest.PrincipalRequests &&
		entry.PrincipalRequests == [2]string{rec.ABRequest, rec.BARequest} && rec.BARequest != rec.BAService, "effect_principal_identity", "case.json", "")
}

func checkEffectCausalAcks(files map[string][]byte, entry revocationManifest, legacy CaseReport, producer EffectObservationProducer) {
	acks := effectJSON[map[string]boundaryEvent](files, "acknowledgments.json")
	endpoint, runtime := acks["Endpoint"], acks["Runtime"]
	need(len(acks) == 2, "effect_causal_ack_pair", "acknowledgments.json", "")
	ready := readRevocationReady("b", files)
	stopped := effectRaw(files, "stopped/inspect.json")
	rows := indexRevocationRows(stopped)
	attempt := rows.attempts[entry.Manifest.Reciprocal.BAService]
	wire := attempt.Request
	need(endpoint.Boundary == "resource.applied" && endpoint.Phase == "post-commit" && endpoint.Scenario == producer.Test &&
		endpoint.PID == ready.PID && endpoint.Token == ready.Token && endpoint.Store == ready.Store && endpoint.PID > 0 && endpoint.Token != "" && endpoint.Store != "",
		"effect_endpoint_causal_ack", "Endpoint", "")
	need(endpoint.Subject == wire.Subject && endpoint.Generation == wire.Generation && endpoint.Epoch == wire.Epoch && endpoint.Sequence == wire.Sequence &&
		endpoint.Intent == wire.Key && endpoint.Attempt == wire.Authorization && len(endpoint.InFlight) == 0 && endpoint.AuthorityAdmissionOrdinal == 0,
		"effect_endpoint_ack_binding", "Endpoint", "")
	kill := effectJSON[struct {
		PID  int
		Wait string
	}](files, "kill.json")
	opened := effectJSON[struct {
		PID              int
		Namespace, Error string
		Runs             map[string]OwnerRun
	}](files, "first-owner/open.json")
	need(kill.Wait == "signal: killed" && kill.PID == entry.Captures[2].OwnerPID && kill.PID == opened.PID && runtime.PID == kill.PID &&
		runtime.PID != endpoint.PID && opened.Namespace == stopped.Namespace && opened.Error == "<nil>", "effect_confirmed_kill", "kill.json", "")
	need(runtime.Boundary == "runner.waiting" && runtime.Phase == "in-flight" && runtime.Scenario == producer.Test && runtime.Subject == "" &&
		runtime.Generation == 0 && runtime.Epoch == 0 && runtime.Sequence == 0 && runtime.Intent == "" && runtime.Attempt == "" && runtime.AuthorityAdmissionOrdinal == 0 &&
		slices.Equal(runtime.InFlight, []string{wire.Authorization}), "effect_runtime_inflight_ack", "Runtime", "")
	var anchors effectAnchors
	need(strictJSON(effectAnchorBytes, &anchors) == nil, "effect_builtin_anchor", "", "")
	for _, a := range anchors.Cases {
		if a.Name == legacy.Name && a.Race == producer.Race {
			want := a.RuntimeIdentity
			need(runtime.Token == want.Token && runtime.Store == want.Store && runtime.PID == want.PID && runtime.Scenario == want.Scenario,
				"effect_sealed_process_anchor", "Runtime", "")
		}
	}
	// The full snapshot/Open run join is already checked by the legacy external
	// checker, including initial->post-open renewal. Require that evidence here.
	owner, ok := legacy.External.Owners["first-owner"]
	need(ok && owner.PID == runtime.PID && owner.RenewalChecked && len(opened.Runs) == 2, "effect_owner_open_join", "first-owner/open.json", "")
	for _, r := range opened.Runs {
		need(r.Epoch == 2 && r.ID.Generation == 1, "effect_owner_original_epoch", r.ID.Subject, "")
	}
	// Exactly one matched ingress and commit; no fabricated gate ordinal in wire.
	checkEffectEndpointPair(files, wire)
}

func checkEffectEndpointPair(files map[string][]byte, wire WireRequest) Receipt {
	ready := readRevocationReady("b", files)
	lines := bytes.Split(bytes.TrimSpace(files["endpoint-b.jsonl"]), []byte("\n"))
	need(len(lines) == 2, "effect_unique_endpoint_contact", "endpoint-b.jsonl", "")
	var ingress, commit endpointEvent
	need(strictJSON(lines[0], &ingress) == nil && strictJSON(lines[1], &commit) == nil, "effect_endpoint_encoding", "endpoint-b.jsonl", "")
	var got WireRequest
	need(ingress.Event == "ingress" && ingress.Contact == 1 && ingress.PID == ready.PID && ingress.Process == ready.Token &&
		ingress.Method == "POST" && ingress.Path == "/apply" && strictJSON([]byte(ingress.Body), &got) == nil && got == wire,
		"effect_exact_ingress", "endpoint-b.jsonl", "")
	need(commit.Event == "committed" && commit.Contact == 1 && commit.PID == ready.PID && commit.Process == ready.Token && commit.Receipt != nil,
		"effect_exact_commit", "endpoint-b.jsonl", "")
	r := *commit.Receipt
	need(r.PID == ready.PID && r.Process == ready.Token && r.Key == wire.Key && r.Profile == wire.Profile && r.Delta == wire.Delta && r.Mutation == 1 && r.Total == wire.Delta,
		"effect_exact_receipt", "endpoint-b.jsonl", "")
	return r
}

func checkEffectExactRenewal(before, after []byte) {
	var a, b RawSnapshot
	need(strictJSON(before, &a) == nil && strictJSON(after, &b) == nil, "effect_renewal_encoding", "reopened", "")
	need(a.Namespace == b.Namespace && a.KeyVersion == b.KeyVersion && len(a.Tables) == len(b.Tables),
		"effect_exact_renewal", "metadata", "")
	for table, rows := range a.Tables {
		next, ok := b.Tables[table]
		need(ok && len(rows) == len(next), "effect_exact_renewal", table, "")
		if table != "runs" {
			need(equal(rows, next), "effect_exact_renewal", table, "")
			continue
		}
		need(len(rows) == 2, "effect_exact_renewal", "runs", "")
		seen := map[string]bool{}
		for i, row := range rows {
			subject := textCell(row, "subject")
			need((subject == "a" || subject == "b") && !seen[subject] && intCell(row, "generation") == 1 && intCell(row, "epoch") == 2,
				"effect_exact_renewal", subject, "")
			seen[subject] = true
			want := Row{}
			for key, value := range row {
				want[key] = value
			}
			want["epoch"] = json.RawMessage("3")
			need(equal(want, next[i]), "effect_exact_renewal", subject, "")
		}
	}
	need(len(a.Tables["runs"]) == 2 && len(b.Tables["runs"]) == 2, "effect_exact_renewal", "runs", "")
}

func checkEffectCallback(files map[string][]byte, entry revocationManifest) {
	ack := effectJSON[boundaryEvent](files, "temporal-boundary.json")
	raw := effectRaw(files, "before-revoke/inspect.json")
	a := indexRevocationRows(raw).attempts[entry.Manifest.Reciprocal.BAService]
	need(ack.Boundary == "authority.admission" && ack.Phase == "allowed" && ack.AuthorityAdmissionOrdinal == 1 &&
		ack.Subject == a.Request.Subject && ack.Generation == a.Request.Generation && ack.Epoch == 1 && ack.Sequence == a.Request.Sequence &&
		ack.Intent == a.Intent.ID && ack.Attempt == a.AuthorizationID && ack.PID == 0 && ack.Scenario == "" && ack.Token == "" && ack.Store == "" && len(ack.InFlight) == 0,
		"effect_synchronous_callback", "temporal-boundary.json", "")
	checkEffectEndpointPair(files, a.Request)
}
