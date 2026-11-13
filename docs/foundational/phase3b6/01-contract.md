# Phase 3B.6: Effect observation across one owner loss

Status: independently reviewed contract, before implementation and qualification. Base research commit `df840103aa3a85e33b4b95ddcf2fe3df071cb024`. Phase 3B.5 PR #46 was open and unmerged at preflight on 2026-11-13. This packet branches from it and does not authorize a merge.

## Question and contribution

For one request-bound service in a two-principal computation, what distinguishes a historical admission, an endpoint application, and a locally retained result after one owner-process loss and exact-root revocation? Can a refused fresh attempt resolve the older original's outcome?

The intended contribution is a finite observation partition and a separately checked, source-conditioned relation to sealed historical captures. Phase 2B/2C already establish lost-result indistinguishability and provisional absence; Phase 3B.5 already distinguishes historical allowance from current authority. This packet studies their intersection without claiming a new general impossibility theorem, a whole-pair recovery theorem, or arbitrary-agent composition.

## Frozen scope

- A new opt-in `research/effectobservation` model represents one selected BA service boundary, not both reciprocal computations. A and B share one Store owner. One death clears that owner's volatile state once. Endpoint and durable Store survive.
- The initial cut already retains an accepted exact request-bound grant, first allow, one debit and one allowed gate. Its original, source generation/revision, grant, exact b root, original activation and wire binding are fixed. The distinct B-owned c root is a control coordinate. Both roots have the same owning principal and different complete permissions.
- One successful gate, at most one launched request, one ingress, one application, one local retained result, one principal consumption, one b revocation, one owner loss/reopen, and optionally one fresh refused gate after reopen. No successful re-exposure, new grant, read/import authority, extra owner, multihop or arbitrary endpoint behavior.
- Application returns either of two opaque authentic values. The endpoint is honest and deterministic conditional on the selected value; the finite branch represents an unknown external result. Result values are not invented from a default, allowance, gate or refusal.
- No production runtime changes or new runtime campaign. The actual idempotent endpoint fixture is a historical anchor; this one-contact model does not establish a non-idempotent retry theorem.

## Operational boundary

Use separate transitions for launch, ingress, application, authentic successful response arrival, local durable retention and principal consumption. Unsuccessful HTTP completion is hidden in this abstraction and cannot create result custody. Launch transfers the packet beyond the owner's volatile opportunity. Ingress and application remain possible after loss and revocation. Application does not set retained result or consumption.

The owner has three stages: original live owner, stopped, reopened live owner. Reopen advances both principals' activation together; it does not change generation, revision, original wire, first allow, debit or gate history. The original response belongs to the original owner incarnation. Losing that owner clears its permit and any unretained response, and cannot deliver a late response to the reopened owner. A retained result survives loss and can be consumed after reopen. No lookup/import transition exists.

Revocation affects exactly b, preserving c and all historical decisions. It does not cancel an already launched packet, retract an application, erase authentic evidence or disable a surviving pre-revocation opportunity. Fresh refusal is enabled only after reopen when b is revoked and the original has no retained result; it changes only the recorded fresh-gate disposition. An applied-but-unretained original can therefore receive a fresh refusal.

Every transition carries a complete descriptor sufficient for an independently interpreted event ledger: kind, original binding, actor/owner incarnation where applicable, targeted root and authentic value where applicable. Values irrelevant to an event must be rejected rather than silently ignored. The finite model may use fixed symbolic bindings, but the report must declare their meaning and not claim they validate concrete cryptographic provenance.

## Observations and falsifiers

`O_store` contains the fixed binding, accepted grant/allow/debit/first-gate history, owner stage/current activation, b/c authority, fresh-refusal history, retained result and consumed result. It excludes permit, launch, ingress, world application and unretained response. The owner-stage evidence is supplied by the experimental death/Open protocol, not inferred from the Store alone.

`O_dossier` extends this with complete identity-bound ingress/application evidence at a declared cut. It does not expose an unobserved in-transit launch. A missing/partial endpoint observation is unknown or unsupported, never a zero. A negative observation states absence at the cut and cannot establish future nonapplication when a launched packet may remain outside the owner. Positive application evidence does not become local result custody.

Enumerate the full finite reachable graph, audit every outgoing edge and merged target with a separately written ledger, and independently group states by both observations. Preserve shortest qualifying witness pairs relative to the post-admission initial cut:

1. Same `O_store`, reopened/revoked/fresh-refused, no local result, with zero versus one application.
2. Same local evidence with two distinct authentic unretained world values.
3. No application at an observed cut, followed by ingress/application after owner loss and revocation.
4. Authentic result retained before death, then consumed after reopen despite b revocation.

Do not reuse Phase 3B.5's atomic `Applied` completion predicate. A state with world application and no surviving result is a counterexample to importing that predicate into this new domain, not a refutation of its original theorem. No fairness, universal completion, reverse simulation or general progress theorem is claimed here.

## Concrete evidence relation

Add `CheckEffectObservationDirectory(dir string) (EffectObservationCaseReport, error)` and a new `research/delegationcorrespondence/cmd/effectobservation -case DIR`. Preserve the complete legacy report beside new, independently named verdicts for raw validity, continuity, physical coverage, protocol/cut attribution and model support. Preserve every old API, report and exclusion.

The initial concrete domain is exactly the two sealed producer protocols below in normal and race modes. Each input must match the pinned producer commit/tree, relevant source-file digests, case/profile/schedule, ordered capture declarations and required identity-bound sidecars. Producer identities are necessary but do not turn editable exports into cryptographic proof; the accepted relation assumes the sealed archive and honest pinned producer. A semantically coherent rehashed dossier is a checker attack, not a newly trusted archive.

- `applied-unknown-revoked-before-fresh-gate`: five captures. `initial` and `first-owner/post-open` precede the model's admitted cut and remain explicit outside-domain prefixes. `stopped`, `reopened`, and `final` may enter the new relation. Endpoint post-commit and runtime in-flight acknowledgments, both consumed before acknowledged SIGKILL under the pinned producer, establish application-before-loss. They have no invented total order between them. The controller releases the endpoint only after acknowledged owner death. Before release the endpoint cannot emit a successful authentic receipt; request cancellation may instead terminate its wait with an error. Continuous barrier occupancy and transport liveness until death are not established. Both principals renew exactly once from stopped to reopened, with all non-run tables and all other run fields identical. Final exact b revocation and refused epoch-current gate preserve one immutable debit/first allowance and no local BA result or principal consumption. One scheduler dispatch after reopen is compatible with zero additional physical ingress; do not report zero dispatches.
- `N07`: three captures. The synchronous allowed-gate callback captures before/after-revoke before the worker can launch. This producer ordering, not final journals alone, supports zero contact/application at those two cuts. Final BA application, retained result and actual principal consumption are joined to their own evidence; the separate AB/c completion remains a concrete control outside this model's completion claims.

Producer commit `5b4662fd682a220ca3065b15e86672b7e56567e7`, tree `44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`. Owning archives are `docs/foundational/phase3b4/evidence/prospective-main/{archives.json,raw-normal.tar.gz,raw-race.tar.gz}`. Inspect the pinned source for temporal semantics. Preserve source extracts/hashes and causal certificates, not merely a copied conclusion.

Raw checks must validate all declared captures and request custody. Physical checks require both ready-bound journals and a coherent final snapshot with full interval coverage. The BA original request and serving effect operation are different identities. Join exact request/wire/key/source/authorization, grant/path/root, reservation/debit, historical gate and current refused gate. Endpoint wire has no gate ordinal; uniqueness is supplied by the one-successful-gate/at-most-one-contact restriction, not a fabricated field.

For the loss profile additionally bind endpoint ack phase/PID/token/Store/key/authorization/source to ready, ingress and committed receipt; bind runtime ack phase/PID/token/Store and sole in-flight authorization to child Open and kill/wait; require one declared loss; validate both principals' exact renewal; bind controlled revoke to its accepted lifecycle fact and exact b tuple, distinct from c. Lifecycle settlements cannot count as external-result settlements. Final-only evidence must not leak into earlier cuts.

A model witness is an existential abstract path consistent with the certified causal constraints and projected cut; it is not a recovered complete physical history. Any unobserved ordering must be labeled as a witness choice. Inferred launch/ingress from a matched committed application are prerequisite steps under the producer/endpoint contract. Report these separately from directly retained local facts. Prefix witnesses for ordered supported cuts must be coherent extensions. Earlier outside-domain prefixes are not silently deleted from the report.

v1 pre-call/applied-no-receipt controls and Phase 2B/2C U0/U1 remain historical comparisons only. They do not establish a matched v2 equal-local-record pair. The new finite model may show such a pair while the sealed v2 corpus supplies only its applied-world branch. New runtime evidence would be required for the matched v2 claim.

## Independent checks and compatibility

The independent event ledger must not call operational successors, authority helpers, transition/classification helpers or operational observation projections. Sharing descriptor types is allowed. Audit every graph edge, shortest state prefix and merged target; compare independently derived observation partitions. Finite graph coverage is not all cyclic-history enumeration.

Include lawful controls and attacks for gate-as-contact, contact-as-application, application-as-retention, loss deleting pending packet/world truth, revived old permit/response, b/c aliasing, retroactive first refusal, fresh allowed retry, forensic-result import, definitive absence, missing artifact as zero, wrong/missing causal acknowledgment, only-one-principal renewal, wrong exact revoke, stale fresh gate, duplicate physical contacts and forged result/custody. Compiled checker mutants must build and fail meaningful assertions, not crash or fail to compile. The causal-ack mutant exercises end-to-end false acceptance. Exact table/run renewal overlaps inherited raw/physical validation, so its mutant is tested at the helper boundary; that is not a new end-to-end acceptance-gap claim. Inspection Settings/Roots equality is outside this renewal contract.

Preserve old sixteen authority-service configurations/order, State/Event JSON/enums, four historical deterministic exports, default command scope, old ledger/progress outputs, Phase 3B.5 model/audit/progress exports, and all ten old revocation concrete report bytes. New behavior is opt-in. No dependencies, unrelated refactors or automatic retransmission.

## Qualification and delivery

Use Go 1.27.4 darwin/arm64 and Staticcheck 2026.2.2 (0.8.2) observed at preflight; record any drift. Prior Go 1.27.3 qualification is historical. Run fresh full ordinary tests, affected race tests, build/module/format checks, scoped Staticcheck, deterministic model/audit exports and both sealed modes. Record inherited global lint diagnostics independently. Freeze exact executable inputs and preserve genuine RED, hostile, mutant and qualification evidence.

Keep Markdown/LaTeX claims, assumptions and qualification aligned and rebuild/inspect the PDF. No new literature citation is required unless actually introduced and checked. Disclose the inherited Markdown-only section divergence rather than claiming global manuscript parity. Obtain independent per-task review and one final whole-branch review. Deliver a draft PR against the parent branch while PR #46 remains unmerged; recheck before delivery. No merge/release.

Work only in the current isolated research worktree. Read the runtime source as needed; do not edit it, other checkouts or `.envrc`. Never open archived SQLite; hash it as bytes. Offline Go environment: `GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b6/go-build`. One implementation agent at a time, fresh task reviewers, no worker subagents. Controller decisions remain reviewable and are recorded as rulings, not presented as user approval of an unpublished design.
