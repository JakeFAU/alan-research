package delegationcorrespondence

import (
	"encoding/json"
	"fmt"
	"slices"

	m "github.com/JakeFAU/alan-research/research/effectobservation"
)

// EffectObservationCaseReport preserves legacy evidence and separates each new
// verdict. Supported is conditional on the declared unsigned producer anchors.
type EffectObservationCaseReport struct {
	Legacy                                              CaseReport
	RawValid, ContinuityValid, ProtocolValid, Supported bool
	Physical                                            ExternalReport
	Producer                                            EffectObservationProducer
	Reasons                                             []string
	Captures                                            []EffectObservationCapture
	Limitations                                         []string
}

type EffectObservationProducer struct {
	Commit, Tree, Test, GoVersion                   string
	Race                                            bool
	SourceSHA256, BindingSourceSHA256               string
	ArchiveSHA256, MemberPrefix, SealedDigestSHA256 string
	SealedBytesMatch                                bool
	SourceFiles                                     map[string]string
	Trust                                           string
}

type EffectObservationCapture struct {
	Path, SHA256 string
	Supported    bool
	Reason       string
	Certificate  EffectCausalCertificate
	Binding      *EffectConcreteBinding
	Store        *m.StoreView
	Dossier      *m.DossierView
	Witness      []EffectWitnessStep
	Audit        *m.Ledger
}

type EffectCausalCertificate struct {
	Kind                                           string
	Artifacts                                      map[string]string
	Constraints, InferredPrerequisites, Unobserved []string
	EndpointAtCut                                  string
	FinalCorroboration                             *Receipt
	OriginalEpoch, CurrentEpoch                    int64
	EpochMap                                       map[int64]uint8
}

type EffectConcreteBinding struct {
	Original      Original
	ServiceID     string
	Roots         [2]RevocationRoot
	Grant         LifecycleResult
	Authorization Authorization
	Attempt       ExternalAttempt
	Reservation   Reservation
	Debit         Charge
	Gates         []Fact
	Revoke        *LifecycleResult
	Receipt       *Receipt
	Consumption   *Fact
}

type EffectWitnessStep struct {
	Event m.Event
	Basis string
}

// CheckEffectObservationDirectory reads declared JSON evidence and hashes binary
// members as opaque bytes. It never opens SQLite or executes the runtime.
func CheckEffectObservationDirectory(dir string) (out EffectObservationCaseReport, err error) {
	legacy, files, err := checkBindingDirectory(dir)
	out.Legacy = legacy
	out.Physical = legacy.External
	out.Reasons = []string{}
	out.Captures = []EffectObservationCapture{}
	out.Limitations = []string{
		"Unsigned sealed archives and honest pinned producer are external trust assumptions; internally rehashed bytes alone are not provenance.",
		"Existential abstract witness with inferred prerequisites; not recovered complete chronology, simulation, fairness or completion.",
		"One selected BA service; AB/c is a concrete control outside completion claims. No matched v2 equal-local zero/one pair or new runtime campaign.",
		"Absence at a certified cut is not future nonapplication; final endpoint truth is not local result custody or import authority.",
	}
	if err != nil {
		return out, err
	}
	for _, p := range legacy.Phases {
		out.Captures = append(out.Captures, EffectObservationCapture{Path: p.Path, SHA256: p.SHA256, Reason: "prerequisite_not_certified"})
	}
	out.RawValid = len(legacy.Phases) > 0
	for _, p := range legacy.Phases {
		out.RawValid = out.RawValid && p.Verdict.Valid && p.Verdict.Custody.Declared
	}
	out.ContinuityValid = legacy.ContinuityFailure == nil
	checkRevocationPhysical(&out.Physical, files)
	if !out.RawValid {
		out.Reasons = append(out.Reasons, "raw_binding_not_valid")
	}
	if !out.ContinuityValid {
		out.Reasons = append(out.Reasons, "capture_continuity_failed")
	}
	if !out.Physical.Valid {
		out.Reasons = append(out.Reasons, "required_physical_evidence_not_valid")
	}
	if len(out.Reasons) > 0 {
		return out, nil
	}
	defer func() {
		if p := recover(); p != nil {
			f, ok := p.(Failure)
			if !ok {
				panic(p)
			}
			out.Reasons = append(out.Reasons, f.Edge+": "+f.ID)
			out.Supported = false
		}
	}()
	var entry revocationManifest
	need(json.Unmarshal(files["case.json"], &entry) == nil, "effect_case_encoding", "case.json", "")
	loss := legacy.Name == effectLossName
	out.Producer = checkEffectProducer(files, legacy)
	checkEffectCaptures(files, entry, legacy, loss)
	if loss {
		checkEffectCausalAcks(files, entry, legacy, out.Producer)
		checkEffectExactRenewal(files["stopped/inspect.json"], files["reopened/inspect.json"])
	} else {
		checkEffectCallback(files, entry)
	}
	// Protocol certification includes exact cut contents, not just labels and acks.
	bindings := make([]EffectConcreteBinding, len(legacy.Phases))
	for i, p := range legacy.Phases {
		if loss && i < 2 {
			out.Captures[i].Reason = "outside_domain_pre_admission"
			continue
		}
		bindings[i] = checkEffectBinding(p, files, entry, loss, i)
	}
	out.ProtocolValid = true
	state := m.Initial()
	events := []m.Event{}
	steps := []EffectWitnessStep{}
	emit := func(kind m.EventKind, basis string) {
		found := false
		for _, edge := range m.Successors(state) {
			if edge.Event.Kind != kind || (edge.Event.Value != m.NoResult && edge.Event.Value != m.Value1) {
				continue
			}
			need(!found, "effect_ambiguous_model_edge", string(kind), "")
			found = true
			state = edge.State
			events = append(events, edge.Event)
			steps = append(steps, EffectWitnessStep{Event: edge.Event, Basis: basis})
		}
		need(found, "effect_missing_model_edge", string(kind), "")
	}
	for i, p := range legacy.Phases {
		if loss && i < 2 {
			continue
		}
		if loss {
			switch i {
			case 2:
				emit(m.Launch, "inferred prerequisite of matched endpoint ingress; before application and loss")
				emit(m.Ingress, "matched endpoint journal; before post-commit acknowledgment")
				emit(m.Apply, "matched committed receipt and post-commit acknowledgment before confirmed kill; authentic value mapped to v1")
				emit(m.LoseOwner, "both acknowledgments consumed before confirmed SIGKILL under pinned controller")
			case 3:
				emit(m.Reopen, "exact both-principal renewal; relative activation 2")
			case 4:
				emit(m.RevokeB, "accepted exact b revoke after reopened capture")
				emit(m.FreshRefuse, "sole fresh current-epoch gate after revoke; does not settle original")
			}
		} else {
			switch i {
			case 1:
				emit(m.RevokeB, "accepted exact b revoke inside synchronous callback")
			case 2:
				emit(m.Launch, "inferred prerequisite; source callback returned after both intermediate cuts")
				emit(m.Ingress, "final matched endpoint journal")
				emit(m.Apply, "final matched authentic application; value mapped to v1")
				emit(m.Response, "inferred authentic response prerequisite of retained external receipt")
				emit(m.Record, "retained external receipt and matching effect settlement")
				emit(m.Consume, "actual principal reply selection and paid custody; intervening reply transport hidden")
			}
		}
		binding := bindings[i]
		store, dossier := projectEffectObservation(binding, loss, i)
		need(m.ObserveStore(state) == store && m.ObserveDossier(state) == dossier, "effect_projection_disagreement", p.Path, "")
		ledger, auditErr := m.AuditTrace(events)
		need(auditErr == nil, "effect_independent_audit", p.Path, fmt.Sprint(auditErr))
		need(ledger.StoreView() == store && ledger.DossierView() == dossier, "effect_independent_observation", p.Path, "")
		out.Captures[i] = EffectObservationCapture{Path: p.Path, SHA256: p.SHA256, Supported: true,
			Certificate: effectCertificate(files, binding, loss, i), Binding: &binding, Store: &store, Dossier: &dossier,
			Witness: slices.Clone(steps), Audit: &ledger}
	}
	out.Supported = true
	return out, nil
}

func projectEffectObservation(b EffectConcreteBinding, loss bool, index int) (m.StoreView, m.DossierView) {
	// Concrete fields have been separately checked; this projection never uses
	// operational successors or final-only local facts for intermediate cuts.
	store := m.StoreView{Binding: m.Binding{Original: "BA/1", SourceGeneration: 1, SourceRevision: 1,
		Grant: "grant:BA:b", Root: m.RootB, OriginalActivation: 1, Wire: "wire:BA:1"},
		AcceptedGrant: b.Grant.Accepted, FirstAllow: b.Authorization.Allowed, Debit: uint8(b.Debit.Effects),
		FirstGateAllowed: b.Gates[0].AuthorityAdmission.Allowed, Owner: m.OwnerOriginal, Activation: 1, ControlRoot: m.RootC,
		BRevoked: b.Revoke != nil}
	if loss {
		store.Owner = m.OwnerStopped
		if index >= 3 {
			store.Owner = m.OwnerReopened
			store.Activation = 2
		}
		store.FreshRefused = len(b.Gates) == 2
	}
	if b.Receipt != nil {
		store.Retained = m.Value1
	}
	if b.Consumption != nil {
		store.Consumed = m.Value1
	}
	dossier := m.DossierView{StoreView: store}
	if loss || index == 2 {
		dossier.Ingress = true
		dossier.Applied = m.Value1
	}
	return store, dossier
}

func effectJSON[T any](files map[string][]byte, name string) T {
	var value T
	var fields map[string]json.RawMessage
	need(strictJSON(files[name], &fields) == nil && json.Unmarshal(files[name], &value) == nil, "effect_artifact_encoding", name, "")
	return value
}

func effectRaw(files map[string][]byte, path string) RawSnapshot {
	return effectJSON[RawSnapshot](files, path)
}

func effectArtifactHashes(files map[string][]byte, names ...string) map[string]string {
	out := map[string]string{}
	for _, name := range names {
		need(len(files[name]) > 0, "effect_certificate_artifact", name, "")
		out[name] = sha(files[name])
	}
	return out
}

// Keep encoding/json explicit at the opaque descriptor boundary.
func effectCanonicalJSON(data []byte) json.RawMessage {
	var value any
	need(strictJSON(data, &value) == nil, "effect_json", "", "")
	return bytesJSON(value)
}
