# Phase 3B.3 request-bound delegation implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development.

**Goal:** Determine whether a small retained authority primitive realizes the
frozen two-agent S2 contract in the actual public runtime.

**Architecture:** Keep the existing forward model unchanged, add independently
checked model/runtime correspondence evidence, and implement only the authority
protocol needed by the experiment in the public runtime repository. Preserve
negative existing-API witnesses before changing production semantics. Authority
safety and completion reachability are separately qualified.

**Tech Stack:** Existing Go 1.27.3, SQLite, ordinary retained programs, local
deterministic endpoints, Tectonic and the existing manuscript sources.

**Spec:** [user packet](00-brief.md), [verified gap](01-runtime-gap.md), and the
[frozen contract](02-contract.md). No production implementation begins until
the exact contract and its falsifiers have been written and builder-checked.

**Execution status:** Tasks 1 and 2 are complete. The public-runtime baseline
source is `8319a9067c16aaf5e3ecccb3d0c62f5c536e6140`; the reviewed primitive is
`4f6f5f2500a48b2ca01ee7b09a823fa63656b7ab`. Its evidence is retained through
`6352969` in runtime `docs/qualification/phase3b3`, including independent reviews
and both rejected contract/candidate corrections. Task 4 precedes Task 3 so the
independent checker consumes the actual raw runtime export. Neither the full
process-cut qualification nor model/runtime correspondence is established yet.

## Global Constraints

- Exactly two logical agents, established A-servicing/B-serial service policy.
- Request-bound one-hop delegation only; complete permission-tuple containment.
- Retained original request, operation and logical generation identities.
- Fresh physical admissions recheck the same selected authority path.
- No ambient borrowing, automatic redelegation, second scheduler or budget ledger.
- Preserve first refusals, unresolved external outcomes and existing expenditure.
- Runtime implementation/qualification belongs in `github.com/JakeFAU/alan`.
- Preserve all earlier research model, fixture, archive and rejected-candidate files.
- Source qualification must use immutable commits; no merge or release.
- Use Staticcheck; preserve pre-existing diagnostics separately.

## Review Focus

- Equal payloads with distinct request/operation/generation identities must not share authority.
- An allowed old decision does not permit a fresh physical admission after revocation.
- A surviving original-bound grant is distinct from an obsolete owner admission opportunity.
- A public runtime replay result is not an independent authority-policy oracle.
- Missing result knowledge, scheduling wait and insufficient authority must remain distinguishable.

## Task 1: Current-runtime realization attempt

**Files:** runtime `examples/delegation/baseline.go`, its focused test, and
`internal/durabletest/phase3b3_baseline_test.go`; research evidence and gap notes.

**Interfaces:** A public `Baseline(scenario BaselineScenario) (runtime.Provision,
error)` builds two ordinary authored roots. Scenarios are `local`, `missing`,
`claim`, `cartesian`, and `result_launder`. Fixed capability identifiers encode
concrete resource/scope tuples; adapters are supplied separately by the caller.
Tests export raw provision, inspections, Replay outcomes, selected original
result custody and independent endpoint evidence under `ALAN_PHASE3B3_EVIDENCE`.

- [x] Build the smallest ordinary public-API reciprocal fixture with A servicing and B serial; no production edits or broad A permission in S2 attempts.
- [x] Run S1 and missing/claim/Cartesian/result-laundering attempts on the unchanged public runtime; retain the genuine refused permission decision and zero unauthorized contacts.
- [x] Demonstrate that the desired current-API S2 completion assertion fails for the expected authority reason, with a passing S1 apparatus control.
- [x] Review the fixture and preserved negative evidence before relying on it.

## Task 2: Freeze and implement the minimal authority protocol

**Files:** runtime new `internal/executor/request_authority.go`,
`request_authority_admission.go` and focused `request_authority*_test.go`;
existing program/lifecycle/effect/dispatch/composition/format/replay/inspection
integration files listed exactly in contract section7; public aliases in
`runtime/runtime.go` and `runtime/inspection.go`. The three baseline files stay
unchanged. No research production executor changes.

**Interfaces:** Exact public types/fields are frozen in contract sections1–4:
`CompositionProvision.RequestAuthority`, `RequestAuthorityVersion`,
`AuthorityRoot`, `DelegationClaim`, `RevocationClaim`, `AuthorityUse`,
`AuthorityDecision`, `EffectAuthority`, `AuthorityAdmission`; optional
`Effect.Authority`/`Intent.Authority`, lifecycle and effect witness fields.
Authority control uses committed lifecycle operations `DelegateAuthority` and
`RevokeAuthority`, never a host-only grant method. Primary gate facts use kind
`authority_admission`, the existing composition order and exact identities in
section4. Contract v2 uses exact `AuthorityRoot.RootRun` and `DelegationClaim.RecipientRoot` logical locators, with namespace-bound IDs derived into retained decisions and witnesses; rejected v1 and its RED remain archived. Formats13/14 are opt-in; no new table. Expose typed accepted/refused
control and gate evidence through public inspection. Internal helper names may
follow existing package conventions; their policy must implement the frozen
checks, not redefine the contract.

- [x] Add focused `TestRequestAuthority*` failures for real delegated permission, missing/wrong root, full-tuple mismatch, exact request/generation/next-reduction binding, first refusal and revocation-before-admission. Retain failing commands before implementing the corresponding production path.
- [x] Add `TestRequestAuthorityFormats*` for old schemas3–12, unsupported feature combinations, version13 creation and14 conflict promotion, and rejection by the frozen old reader.
- [x] Demonstrate a positive ordinary two-agent Run in focused tests, while using direct public/owned transitions only for deliberately controlled adverse schedules; no alternate scheduler implementation. This smoke test uses empty acknowledgments; Task 4 must establish full requester result custody.
- [x] Resolve builder findings and freeze exact representations, boundaries, versioning and scope in `02-contract.md` before production edits.
- [x] Record precise file/API interfaces in this plan from that freeze; implement through retained ordinary control decisions and existing effect authorization/admission.
- [x] Attack tuple mixing, source ancestry, original binding, sticky refusal, revocation timing, owner activation and malformed storage first.
- [x] Qualify unchanged old-format behavior and explicit new-format validation; no speculative migration or framework expansion. The ten old-created fixtures give finite compatibility evidence, not universal equivalence.
- [x] Independent authority/minimality and attenuation/request-binding review; correct and re-review substantive findings.

## Task 3: Independent correspondence corpus

- [ ] Generate the bounded corpus from unchanged Phase 3B.2 graphs, selecting actual BA delegated witnesses rather than generic AB own-root cuts.
- [ ] Freeze raw-evidence extraction and independent authority/custody checking interfaces from the runtime representation.
- [ ] Preserve all eight required history classifications, mutation controls, unsupported mappings and discrepancies separately.
- [ ] Independently review model/runtime correspondence, progress and retained obligations.

## Task 4: Runtime recovery qualification

- [ ] Exercise real two-agent S1/S2 service and adversarial S3/Cartesian attempts through the actual public runtime.
- [ ] Execute the ten requested acknowledged process-cut boundaries, fresh-owner reopen, endpoint evidence, exact original selection and no invented grants/results.
- [ ] Keep the supported stable-key idempotent endpoint contract separate from generic unknown-adapter blocking and selective revocation obstructions.
- [ ] Independently review recovery, revocation and physical admission.

## Task 5: Freeze, reconcile, review and deliver

- [ ] Freeze source in both repositories; run ordinary, relevant prior-model, integration, race, oracle, formatting and Staticcheck qualification with exact command records.
- [ ] Reconcile Markdown/TeX/BibTeX implications, rebuild PDF and inspect changed pages.
- [ ] Archive interpretable raw evidence, identities, digests and preserved rejected candidates; independently review manuscript/reproducibility and complete whole-change review.
- [ ] Push narrowly scoped branches and open draft research/runtime PRs; verify final identities and clean working state. Do not merge or release.
