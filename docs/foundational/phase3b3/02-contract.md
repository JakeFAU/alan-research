# Phase 3B.3 frozen request-bound authority contract

Architecture-only source review, 2026-11-04. Runtime inspected: `/private/tmp/alan-phase3b3-runtime`, HEAD `af6943a5805a9a6c0e8fe77f804264f4765d07ec`. Requirements, runtime gap review, and correspondence design were read. The builder review itself made no source edits and ran no tests. Contract version 2 corrects the rejected version-1 literal-root construction before production implementation; the failed contract, source and RED log are preserved separately. This is intended behavior to attack and qualify, not a claim of implementation correctness. The user packet and verified gap govern its scope.

## Decision

A narrow implementation remains feasible. Use **explicit `DelegateAuthority` and `RevokeAuthority` operations in the existing lifecycle intent/result protocol**, with a separate authority policy branch. Reuse its immutable request/result rows, source authorization, lifecycle decision sequence, ordinary runner delivery, settlement, and composition fact ordering. Do not reuse `lifecycleGrant` to authorize these operations: a broad lifecycle permission is neither necessary nor sufficient for delegating an effect permission. The legitimate source effect tuple itself is the root. No new scheduler, accounting ledger, principal class, or transitive grant graph is necessary.

Record each physical gate as a **typed primary `authority_admission` record in existing append-only `composition_facts`**. The existing `adapter_dispatch_admissions` cannot represent repeated raw external exposures, refusals or explicit activations, but that does not make a new table necessary. A new fact subtype provides the needed primary evidence under the existing `decision_order`. This deliberately extends the current convention that each composition fact mirrors another source row; that extension is explicit, version-gated and independently validated. No second table or chronology is introduced. See the pre-freeze challenge addendum below.

Support both existing external and adapter effect decisions with one common authority witness predicate. This does not require compiler handle APIs. Use raw external retained programs for reciprocal progress, and existing authored adapter programs for the meaningful resource/operation Cartesian control and generic unknown-outcome control. B can be a raw authority-issuing program in that control; A can be an ordinary annotated adapter program. No program needs to mix new authority control with Author metadata for the minimum corpus. `adapterRequest` currently rejects a program without Author metadata (`adapters.go:79–83`); do not claim raw unannotated adapters work.

This is a small semantic primitive with several integration points. Touching replay, inspection, and version routing is necessary durable-protocol work, not evidence of a prohibited second authority subsystem.

## 1. Explicit supported domain and profile

Frozen exact opt-in:

```go
// Existing CompositionProvision gets one optional field.
RequestAuthority string `json:",omitempty"`
// The sole supported nonempty value is:
const RequestAuthorityVersion = "alan/request-authority/v1"
```

The selected storage formats are **13** (request authority) and **14** (request authority plus conflicting-receipt journal). Both use the existing `alan/composition/v1` profile and schema-4 composition protocol underneath, just as 9–12 route to that protocol. `protocolSchemaNumber(13|14)=4`; `hasRequestAuthoritySchema(13|14)=true`; `hasReceiptConflictSchema(14)=true`; conflict promotion is exactly 13→14. Neither version reports the legacy dispatch-admission or continuation feature merely by numeric ordering. These numbers are free in the inspected baseline; reserve them before edits and reject a collision if another branch has since assigned them.

For this slice require composition provision with two existing static generation-one principals, each included in an existing composition root, adequate existing resource limits, no templates/dynamic members, no Reconciliation, no Continuations, and `RecordDispatchAdmissions=false`. The new composition fact subtype is the gate evidence for this opt-in profile. Reject combinations instead of implementing an unqualified cross-product of features. Keep the normal composition invocation tables present if existing creation expects them; no invocation is admissible without a template. Use existing effect reservations and charges unchanged. A delegated effect consumes **the serving principal's existing root allowance**, exactly once per original allowed reservation. Grants, revocations, rechecks, and repeat idempotent contacts neither reserve effects nor refund expenditure.

Without RequestAuthority, old creation, operation recognition, JSON bytes, authority predicates, and schemas 3–12 retain their old semantics. All added fields use `omitempty`. The authority operations, selection fields, typed witness fields, and new fact subtype are forbidden when the feature is absent. Unknown feature strings fail closed. Old reader rejection follows from storage 13/14 and unknown descriptor fields; test it, do not assume it. `SemanticVersion` for the reducer need not change: optional literal retained fields are copied as data and gated by the Store protocol. Program validation, creation and replay must reject their use outside the opt-in profile.

## 2. Exact representations

Names below fix the implementation interface; they are not claims that these public types already exist. Public aliases belong in runtime. No host method may mint a grant from an uncommitted claim; the only grant writer is the ordinary committed lifecycle path.

```go
type AuthorityRoot struct {
    Run        RunID       // original holder and logical generation
    RootRun    RunID       // exact logical locator: existing ResourceRoot.Run
    Permission Permission  // exact six-field original permissions row
}
type DelegationClaim struct {
    Version            string // RequestAuthorityVersion
    Root               AuthorityRoot
    Recipient          RunID
    RecipientRoot      RunID  // exact existing recipient ResourceRoot.Run locator
    RecipientRevision  int64
    OriginalRequestID  string
    ServiceOperationID string
}
type RevocationClaim struct {
    Version string
    Root    AuthorityRoot
}
type AuthorityUse struct {
    Version string
    Path    string // "own" or "delegated"; all other values invalid
    Root    AuthorityRoot
    GrantID string // empty for own; actual issuer lifecycle Intent.ID for delegated
}
```

`Root.RootRun` and `RecipientRoot` are nonempty, generation-bound logical locators, not wildcards or hashed IDs. The Store resolves each against the immutable retained composition profile in its own namespace, checks exact holder/recipient membership, and derives the canonical namespace-bound IDs. Holder `Root.Run` and resource-root locator `Root.RootRun` are distinct coordinates, even when this fixture sets them equal. A root locator alone confers no permission. `Permission.Source` must equal `Root.Run.Subject`; Permission.Revision binds that holder revision. Do not confuse this holder provenance with the action tuple after changing recipient. Attenuation is exact equality of the action component `(Kind, Target, Operation, Obligation)` plus the root scope, with a singleton grant subset of that issuer's original relation. No wildcard semantics, separate dimension sets, delegated parents, redelegation flag that can be enabled, or implicit fallback. Every grant is permanently non-redelegable by protocol definition.

For adapter operations the required action is `(adapter, intent.Target, intent.Obligation, intent.Obligation)`, exactly as current adapter permission lookup. The fixture's immutable adapter deployment maps a capability to one fixed physical resource endpoint; the operation/obligation distinguishes actual allowed actions, with trusted adapters checked to implement those distinctions. A claimed resource string in payload does not change that binding.

For external operations the required action is `(external, intent.Target, "apply", intent.Obligation)`. Version 13 requires exact obligation equality; legacy external lookup currently omits obligation, and remains unchanged outside opt-in. An external endpoint URL is the physical resource scope, and `apply` is the actual operation. Do not pretend two arbitrary obligation labels create two different physical operations. Use the adapter corpus for the mandatory `(resource1, operation2)` cross-pair unless a separate endpoint fixture actually implements and validates distinct physical operations. External delta/profile bytes remain bound to the exact service operation and retained attempt, not a general range permission.

Add `Authority *AuthorityUse` (`omitempty`) to retained `Effect` and `Intent`, copied unchanged by emit. This is control metadata outside Payload, so adapter input schema and external wire bytes remain unchanged. Literal raw provision programs can encode it; no compiler Send/Receive-handle API is needed. The existing `emit` Dst already receives its canonical original ID. An absent selector on an opt-in external/adapter operation results in a retained `authority_missing` denial; it never searches requests, grant-shaped payloads, results or another principal's permission. An invalid encoding fails validation; a well-formed but false claim gets a retained refusal.

Delegate payload is exactly `operation=DelegateAuthority`, `revision=RecipientRevision`, and `claim=<canonical JSON DelegationClaim>`. Intent target/generation equal Recipient, and its obligation identifies its ordinary lifecycle response. Revoke uses `operation=RevokeAuthority`, target/generation equal Root.Run, revision=Root.Permission.Revision, and a canonical RevocationClaim. Canonical JSON round-trip comparison plus strict typed decoding rejects duplicate, unknown, missing required or differently represented fields; bounded claim size uses existing program/output/data limits. A claim cannot supply its own grant identity.

Do **not** require namespace in the program claim: Create generates it randomly, so precompiled literal programs cannot know it. Derive namespace into the retained authority decision from the owning Store. Validate it against metadata at every read; it is not caller input. The same constraint applies indirectly to hashed resource-root IDs: programs name exact logical root locators instead. Decisions retain canonical IDs derived from namespace and those locators, and membership validation separately establishes whether the claimed roots actually exist.

```go
type AuthorityDecision struct {
    Version, Namespace string
    Root AuthorityRoot
    RootID, RecipientRootID string // derived namespace-bound root candidates
    // Existing Request holds the exact original claim bytes and source.
    // Existing lifecycle Accepted/Reason holds the single decision.
}
// Add *AuthorityDecision Authority, omitempty, to LifecycleResult.

type EffectAuthority struct {
    Version, Namespace string
    Use *AuthorityUse // nil only for an absent selector refusal
    RootID, ServingRootID string // derived namespace-bound root candidates
    Allowed bool
    Reason string
    // For delegated path only, copied from the actual accepted grant:
    OriginalRequestID, ServiceOperationID string
    RequestObservationID string
    RequestConsumedSequence int64
    GrantDecisionSequence int64
}
// Add *EffectAuthority Authority, omitempty, to ExternalAttempt/AdapterAttempt.
```

For every syntactically valid claim/selector, derive `RootID = compositionRootID(namespace, Root.RootRun)` even if the locator or membership is false; such a candidate ID is not evidence that a root exists. Delegate decisions similarly derive `RecipientRootID`; revoke decisions leave that field empty because there is no recipient claim. Effect witnesses derive `ServingRootID` from the actual serving member's retained root. A missing selector alone has an empty RootID. The typed decision/witness validator recomputes all these fields and checks membership as part of the accepted/refused policy, rather than interpreting any candidate ID as authority. Primary gate outer RootIDs contain only actual retained serving/issuer membership roots, never a nonexistent candidate root.

The existing attempt `Allowed/Reason` is still the combined source/authority/resource decision. The nested witness keeps permission allowance distinct from a subsequent resource refusal. A permission witness is immutable even when reservation is denied. Denial witnesses retain the requested selector and exact failed edge; no fake root or grant is invented when the selector is absent. Use is a pointer as shown: nil means an absent selector and is legal only with the retained `authority_missing` refusal. A non-nil selector never silently changes path.

## 3. Issuance, service binding, and temporal semantics

At the first Delegate decision, in the existing owned transaction:

1. Reload the exact committed lifecycle intent and require its real RunID to equal Root.Run. Its retained source revision must equal Permission.Revision. A caller-supplied LifecycleRequest cannot substitute for the stored request.
2. Require source authorization and **current** source open/generation/revision/lifecycle validity at issuance, not merely an old sticky source allow. Check the root is that static holder's original provision root and exact complete Permission row; never resolve Root through a grant.
3. Require the root tuple not revoked before this decision. Require exact recipient current RunID, revision and root identity; recipient must be open for issuance.
4. Reload OriginalRequestID: authentic committed message, actual source Root.Run, actual target and generation Recipient. It need not yet be authorized or delivered: issuance establishes consent for an authentic committed request, and the later service decision requires its legitimate accepted admission and selected observation. Requiring an already retained request authorization here would introduce a hidden delivery-order prerequisite. A payload field named request_id is not evidence.
5. Parse ServiceOperationID as exact canonical `[recipientSubject,recipientGeneration,sequence,0]`, with positive sequence, no alternate encoding or ordinal. It may be future; its binding is checked again when the real effect exists.
6. Retain accepted/refused LifecycleResult and its composition lifecycle fact atomically. Before and After target states are equal for these operations. No permission row or subject revision changes. Existing repeated DecideLifecycle returns the same decision, never reevaluates it into a different answer.

At the first effect authority decision, require the exact selector. Own means the actor's own legitimate unrevoked original row; delegated means the actual accepted issuer lifecycle decision identified by GrantID. Require the grant's root, recipient RunID/revision/root, exact real effect ID, and action tuple to match. Recheck current root validity. No own→delegated or delegated→own fallback is allowed after selection or refusal, even if another valid path exists.

Causal link: follow OriginalRequestID's accepted message admission to the actual generation-bound recipient observation. Require the recipient history row at sequence `service.Sequence-1` to have selected precisely that observation, its `consumed_sequence` to agree, and the service effect to be ordinal zero in the immediately following committed reduction. Bind all stored descriptor/index identities. This is an explicitly **immediate-next-reduction service contract**. It does not infer arbitrary dataflow from some earlier receive. Future-ID prediction is safe only because the later authentic intent and the retained receive are checked; predicting a locator alone cannot authorize a different operation.

Revocation is monotone for the exact AuthorityRoot. Only that legitimate holder's ordinary RevokeAuthority intent can revoke it. Source must be currently valid at the decision. A duplicate root revocation may return accepted `already_revoked` with no extra semantic effect; it cannot un-revoke. A different root tuple remains valid, including the same holder's r permission. Accepted grants remain retained but become unusable for new authority decisions/admissions. Revocation neither closes principals nor changes prior allows/results or accounting.

Root validity at decision/admission means exact original holder RunID/current generation, original revision, static root membership, non-retired holder, exact permission, and no earlier exact-root revocation. Serving source openness/scope checks remain the existing final gate checks. **Normal delegator completion alone does not revoke its issued grants.** Issuance required a live source, but later grants express retained consent; requiring the issuer to remain open forever would add a model constraint and can break a completion suffix. Retiring/changing the issuer's generation or revision invalidates the root; the experiment does not add cancellation propagation. Explicit root revocation remains the qualified temporal control. If project review instead requires root-scope closure to invalidate issued consent, label that additional restriction and test its progress consequence before claiming correspondence.

## 4. Every fresh physical gate

Add no table. Versions 13/14 select the existing composition DDL with only `'authority_admission'` added to its closed kind CHECK, following the existing continuation/reconciliation DDL pattern. Extend `compositionFact` with an optional `*AuthorityAdmission` field, forbidden on all other kinds and required on this kind:

```go
type AuthorityAdmission struct {
    Version, Namespace string
    Ordinal int64              // contiguous per original, beginning at 1
    Source Authority           // actual Run/Epoch/Sequence/Revision at this gate
    AttemptDigest string       // exact immutable attempt including selected witness
    Allowed bool
    Reason string
}
// compositionFact.AuthorityAdmission *AuthorityAdmission `json:",omitempty"`
```

Outer fact fields are exact and non-overloaded: Kind=`authority_admission`; Run=descriptor.Source.Run; Sequence=descriptor.Source.Sequence (actual current reduction prefix, not original sequence and not ordinal); OperationID=original effect ID; RootIDs=sorted unique serving root and selected issuer root when delegated; ID=`compositionIdentity("alan/fact/authority-admission/v1", [Namespace, OperationID, Ordinal])`; RecordDigest=`compositionIdentity("alan/authority-admission/v1", descriptor)`; Order=the existing next composition decision_order. Version remains compositionFactVersion under the explicit new storage/profile gate. Descriptor Version is RequestAuthorityVersion. The state preceding this gate is precisely facts with Order smaller than this fact, so no AfterFact or second sequence is stored.

AttemptDigest binds exact prepared adapter request or raw external Request and its selected permission witness. Neither it nor a namespace is supplied by the program. The tuple is resolved through that exact attempt, never independently mixed labels. Per-original ordinals distinguish repeat idempotent contacts; all admissions/refusals share the existing total order with revocation, authorization, reservation and result events.

Extend registerDispatch under existing `mu` and `ownerMu`: reload attempt; require original allowed reservation; evaluate current source gate and the **same selected** root/grant; capture current activation; append and read back the exact allowed/refused composition fact in the same owned transaction; commit before releasing a physical opportunity. Its existing composition decision_order distinguishes revoke-before-gate from revoke-after-gate without wall-clock assumptions. Gate failure or journal write failure releases no worker. No admission row is a resumable permit: reopen loses all volatile opportunities and a new exposure obtains a new ordinal/current-epoch check. No grant is created on reopen. Do not call a retained admission proof of an actual endpoint contact. A new kind-specific append path must validate an actual current gate, not accept a caller-assembled fact or require a nonexistent duplicate source table.

Maintain existing direct-dispatch and runner coverage; both must enter this gate. An idempotent external retry can have multiple allowed ordinals and the same original key/attempt. A revoked root can yield a later refused ordinal. Generic adapters still do not automatically retry an allowed original with no result. Preserve the first authority allow, debit, unknown state and unpaid obligation on refusal; do not overwrite the original attempt or synthesize an external result.

Returned results bind the same original attempt as today and require a prior allowed journal entry at an earlier causal prefix. Local completion carries its gate ordinal for acknowledged boundary logs. The frozen minimum stores no new result-arrival field: retained results establish original-operation authority and an earlier valid exposure, while the independent contact/boundary trace establishes which local exposure returned those bytes. Do not claim the Store alone uniquely identifies the returning exposure. Do not alter the endpoint's immutable key or request.Epoch on retry. The endpoint request.Epoch is the original authorization epoch today, not the current physical activation. The new admission descriptor is the latter's authority evidence. Independent external-contact evidence is still needed to distinguish admitted-but-not-called from called-but-unknown and to map the observed receipt/application to actual transport history. A result surviving revocation is consumed using its original identity without a fresh mutation admission.

Replay and typed inspection must reconstruct validity at **historical** lifecycle/reservation/gate prefixes. Looking at today's revoked set when validating an old allow incorrectly destroys legal retained history. Conversely, checking only a final permission set misses revoke-before-admission. Use the composition prefix for cross-root chronology; the lifecycle decision sequence alone orders only lifecycle results. Ensure authority accepted/refused decisions and all gate records appear in typed inspection, alongside pending/unknown operations and unchanged expenditure.

## 5. Positive handshake and ordinary scheduling

Do not rely on the runner delivering a grant before a request just because of sort order. Use this ordinary two-agent handshake:

1. A emits AB, then waits for a distinct `grant-ready` message before selecting BA.
2. B emits the genuine BA request, then DelegateAuthority naming that request and A's predeclared exact future service locator. B waits for its own retained lifecycle outcome.
3. Only after a matched **accepted** lifecycle result does B emit grant-ready to A. That message merely controls when A chooses to serve; it grants no authority.
4. A selects grant-ready, then selects the genuine BA request, then emits its service effect in the immediate next reduction with the exact delegated selector. A continues servicing while AB remains unpaid.
5. B remains serial for reciprocal service: it consumes the authentic BA reply/result before selecting AB and emitting its own-root service. Both original obligations must actually be consumed.

B's request can be admitted/queued before its grant exists; A's program avoids making a first effect decision until acceptance. Crash/reopen may expose those staged states. They are implementation intermediates, not literally the model's atomic request-plus-grant event. For the initial implementation, freeze the direct chronological mapping to the unchanged **S2/LateGrant** oracle with grant retention before first permission decision. It is an authority-safe subset of the one-hop S2 contract, not a literal runtime realization of S2/Explicit's inline atomic Commit. Preserve every intermediate prefix and keep the sixteen old graphs unchanged. A later macrostep proof could justify an abstract S2/Explicit projection, but no such proof is assumed here. Report this narrowing/disagreement visibly in results and manuscript. A refused/lost grant-ready does not count as authority-safe complete. Independently exercise the no-handshake late-grant case: first denial must remain after later accepted grant and reopen.

No new scheduler or delivery prioritization is needed. A raw Program's existing emit Dst can retain its emitted request ID; because building arbitrary typed claim JSON from dynamic pieces is not currently an author convenience, the minimum finite fixture predeclares exact sequence locators and validates all referents. This is usable as low-level public Artifact JSON but brittle under program edits. Document that limitation rather than claiming an ergonomic general API. Dynamic construction, rich service bodies, and compiler handle support are future work.

## 6. Recovery/progress envelope

Use the existing idempotent external protocol for the S2 progress witness: stable original key, trusted loopback endpoint, atomic endpoint ledger/deduplication, eventual successful response when required authority remains valid, enough existing root budget, and the two-agent A-servicing/B-serial program. The endpoint may apply once and return the original receipt on a later authorized contact. This does not prove zero duplicate contacts, nor generic exactly-once effects. Preserve all contact attempts separately from applications.

After revocation an unknown external original cannot be re-exposed even if deduplication might make the next call harmless: a fresh physical admission requires current authority. An already applied unknown original may therefore remain unpaid. After a retained result, no reissue is needed. Generic adapter allowed/no-result states remain honestly unknown and can lack a completion suffix. Same-Store reopen renews both open activations; no independent-owner recovery is established.

The old finite model's every-reachable-state completion reachability remains a model fact under its assumptions. Do not lift it to every runtime state by hiding unknown effects, first refusal, staged grants, missed ready messages, or depleted resources as harmless stuttering. Classify each differential mapping and each SIGKILL cut explicitly.

## 7. File touch map

Existing paths verified in this checkout; new files proposed:

- `internal/executor/request_authority.go` (new): typed claims/root/use/decision, canonical checks, issuance/revoke policy, whole-tuple effect authority, historical prefix helpers. Keep production helper independent from the research checker.
- `internal/executor/request_authority_admission.go` (new): typed primary composition admission fact, gate writer, ordinal indexing and historical validation; no table.
- `internal/executor/program.go`: optional literal Authority metadata on Effect/Intent; deterministic copy and validation. No new scheduling instruction.
- `internal/executor/lifecycle.go`, `lifecycle_replay.go`: gated operation recognition, separate policy branch, typed optional decision, immutable no-state-change results, historical revalidation. Keep old Retry/Cancel/Abandon/Retire code path unchanged.
- `internal/executor/adapters.go`, `external.go`: selected authority witness in attempts and composite first decision; old predicates remain for old formats; validation uses historical authority.
- `internal/executor/dispatch.go`: gate current-root recheck/journal before worker release; both adapter and external paths. `runner.go` only if plumbing gate ordinal/boundary evidence requires it; no new queue/policy.
- `internal/executor/composition_types.go`, `composition_store.go`, `store.go`, `receipt_conflict.go`: feature/profile/version13/14 routing, exact allowed combination, opt-in composition kind CHECK extension, open/create flags and 13→14 conflict promotion.
- `internal/executor/composition_facts.go`, `composition_prefix.go`, `composition_replay.go`: typed lifecycle facts and immutable attempt digest integration, prefix authority/gate interpretation, no second chronology. `composition_reservation.go` should need no semantic accounting change.
- `internal/executor/inspect.go`, `inspection_execution.go`, `inspection_deployment.go`, `inspection_history.go`, `inspection_effects.go`, `inspection_effect_order.go`, `inspection_causality_authority.go`, `inspection_causality_chronology.go`: feature/fact decoding, types and independent structural/temporal consistency. Exact subset determined by compile/reference search, not a mandate to touch each.
- `runtime/runtime.go`, `runtime/inspection.go`: public aliases and views. `author/*` is not required for the minimum raw-external/annotated-A-control corpus; do not add an undeclared-kind bypass.
- New focused runtime/executor tests: authority policy, causal selection, gate order, reopen/corruption, format compatibility, public consumer and real SIGKILL witness. Boundary hooks can add authority-specific acknowledged phases under `internal/boundary` and lifecycle/dispatch without sleeps.
- Research-only specification/model/extractor/checker/witness/manuscript changes stay in the research repository; runtime qualifications and production source stay in public runtime.

## 8. Required attacks before accepting implementation

- A lacks b, B owns b; source Authorization succeeds but missing/ambient selector refuses with zero contacts. B request, B result, or grant-shaped payload cannot change it.
- Authentic accepted B→A grant makes A's exact BA service succeed; remove B's legitimate original tuple and keep every other field unchanged: issuance refuses.
- Exact Cartesian relation `(resource1, operation1),(resource2, operation2)` permits neither `(resource1, operation2)` nor `(resource2, operation1)`. Use actual distinct physical operations, and mutate root claim, grant, selected attempt, admission witness and result linkage, not just final permissions.
- Revoke only B/b; prove B/r remains allowed at issuance/decision/admission and can still complete its real unrelated action.
- Request/recipient/issuer generation, root scope, revision, equal-payload operation, wrong sequence, ordinal, original versus retried invocation, prior selected observation, and missing selection mutations all fail at the correct boundary.
- GrantID references result/message/another grant instead of issuer lifecycle decision: reject. Another service tries same grant: reject. A cannot delegate B-derived authority onward or back.
- Grant accepted after original effect denial cannot repair that original. An explicitly selected invalid grant cannot silently fall back to own authority.
- Revocation before issuance/decision/admission refuses; revocation after released admission does not retroactively invalidate that opportunity; revoke after result retention allows original result selection without transport.
- Crash between allowed gate retention and physical call; crash after physical call with unknown result; reopen under new epoch; new external gate is rechecked and logged, no duplicate reservation, no restored volatile opportunity.
- Current revocation must not cause Replay/Open to reject an historically valid old result; missing earlier revoke/forged prefix must be detected independently.
- Malformed canonical JSON, duplicate keys, descriptor/index mismatch, unknown feature, feature without correct storage, extra authority fields in legacy store, missing gate required by a retained result, repeated ordinals, activation outside retained historical bounds, admission before original reservation, and substituted attempt digest all fail closed. Deletion of a terminal gate without subsequent retained evidence, or a fabricated but policy-consistent gate, additionally requires independent boundary/contact evidence; neither storage layout proves those facts by itself.
- Whole old-format suite plus format-specific 3–12 creation/open/replay/inspection snapshots, 13→14 receipt conflict promotion, and rejection by a frozen old binary. No blanket claim of compatibility from unchanged base DDL.

## 9. Builder decisions that should not remain implicit

The frozen choices are: lifecycle reuse; 13/14 opt-in; singleton exact tuple; typed Effect/Intent selector outside Payload; derived namespace; root-specific monotone revoke; no path fallback; immediate-next-reduction causal link; grant-ready handshake; all physical gate decisions journaled; serving-root accounting; ordinary raw external idempotent progress; genuine adapter Cartesian control; no automatic generic retry.

The remaining encoding choices are settled as follows: pointer `*AuthorityUse` in EffectAuthority represents missing-selector refusal; retained result authority uses its immutable original attempt plus prior valid gate/prefix, while the acknowledged boundary log and independent endpoint trace provide exposure-specific arrival evidence. No new retained result wrapper or remote receipt/key format is introduced. Store-only exposure-specific receipt attribution is not claimed. Separate grant chronology is tested against unchanged S2/LateGrant with the before-first-decision condition; inline S2/Explicit remains a distinct unrealized atomic boundary in this minimum candidate.

The implementation stop criterion is semantic, not line count: stop if these exact checks force a general scope language, transitive capability graph, new scheduler/accounting, generic unknown-result retry, or a wider actor/lifecycle model. No inspected source establishes such an obstruction yet. A bounded reusable authority protocol with a typed primary gate record in the existing composition journal is a justified next implementation attempt, subject to independent attacks and review.


## Pre-freeze challenge addendum: primary composition fact replaces new table

The parent challenged whether the additional table was necessary under packet §6. **It was not established as necessary, and the recommendation is corrected above: reuse composition_facts.** This addendum and the revised Section 4 supersede the original standalone-table design. This is a design correction before implementation, not evidence that either implementation ran.

Actual inspected constraints:

- `composition_store.go:19–30` has a closed kind CHECK assembled through before/after constants; adding one opt-in kind is already an established DDL construction pattern. Versions 13/14 can use `compositionSchemaBeforeRead + ",'authority_admission'" + compositionSchemaAfterRead`. Legacy DDL remains byte-for-byte unchanged. A semantic/storage version is still unavoidable to distinguish the new recognized kind and descriptor semantics; a new table is not.
- `composition_facts.go:107–169` currently validates each fact by constructing its expected shape from a separate retained row, then assigns contiguous decision_order and inserts exact canonical descriptor bytes. The new gate is a **primary record**, not a digest reference to a second copy. Make one explicit branch for this subtype. Its producer reloads the real attempt/current owner/authority and computes the exact gate result under the Store lock, builds the typed record, assigns the next per-original ordinal, then uses the existing global order allocation and canonical insertion/readback. Do not remove source-row reconstruction for existing kinds.
- `composition_replay.go:385–489` constructs an expected map entirely from other tables and rejects count mismatches. A new primary record therefore requires an explicit validation branch, not a simple expected-map append of untrusted self-data. Keep exact reconstruction/count coverage for **all mirrored facts**. Count primary admission facts separately; validate their field shape, canonical ID/digest, exact original intent/attempt/root membership, unique contiguous per-original ordinal, actual source history position/epoch bounds, first allowed reservation, selected grant ancestry, and allowed/refused policy at that precise earlier prefix. Advance state only after validating. Unknown subtype/feature remains an error. Require every retained actual result to have an earlier allowed primary gate for its exact original, and preserve genuine admitted-but-no-result prefixes.
- `inspection_causality_chronology.go:182–195` presently indexes most facts by `(kind,operation)`, which would collapse repeat gates. Add a dedicated gate key `(operation,ordinal)` and a transaction node for **each** gate. Connect original attempt, selected grant and source activation witnesses; maintain shared global fact ordering. `inspection_effect_order.go:54+` needs a direct gate-kind case that recomputes historical policy; do not continue feeding these records through the old adapter-only side-table index.
- `validateCompositionFacts` is the foundation-profile validator and only accepts execution/lifecycle. Do not broaden that legacy profile. The new feature requires composition/v1 and goes through its full protocol validator. Audit call sites to confirm this routing rather than modifying all kind switches indiscriminately.

The independent research checker must interpret the primary gate as recorded runtime evidence and derive its permissibility from prior immutable facts. Hashing the record is an encoding/integrity check, not an independent authorization oracle. It cannot prove an actual physical contact occurred; an extra isolated legal gate inserted into a hostile store may be a policy-consistent fabricated history. A separate table would not solve that: both storage designs are unsigned retained evidence, and neither proves omitted or fabricated external contacts. Bind acknowledged boundary and independent endpoint contact histories to each real execution and state this limit. Likewise, historical activation validation must be honest about reopen intervals with no reduction: existing validateWitnessEpoch bounds epochs between neighboring history rows; it does not alone prove the exact wall-clock location of every empty reopen. The live writer captures actual epoch under owner exclusion, and SIGKILL/owner logs independently test fresh activation behavior.

Comparison: the new table would preserve the old 'every fact mirrors another row' convention, provide a convenient FK and uniqueness key, and leave fact descriptors unchanged. But it duplicates gate representation/order concerns and adds a table solely to hold a primary record that the existing versioned append-only journal can represent. The typed-fact option changes a clearly isolated journal convention, removes the duplicate record and AfterFact join, makes revoke/gate order direct, and stays inside one scheduler/accounting/history protocol. Under the packet's preference for existing mechanisms, **the typed primary fact wins**. No concrete source evidence justifies the extra table for this slice.

## Pre-implementation correction: logical root locators (contract v2)

The first implementation probe falsified v1's literal-root construction before any production edit. `Create` draws a fresh namespace, and `compositionRootID` hashes that namespace with `ResourceRoot.Run`; supplying another Store's exact Root.ID is rejected by ordinary provisioning. Literal Effect.Authority values therefore cannot contain the new Store's hashed root IDs before creation. The original failed test and log are retained in the runtime qualification archive; the rejected contract and implementation report are under `evidence/rejected-contract-v1`. This is a representability defect in the proposed encoding, not an impossibility result for request-bound delegation.

Version 2 replaces program-supplied hashed root IDs with exact `RootRun`/`RecipientRoot` logical locators and retains the derived namespace-bound IDs only in Store-produced decisions/witnesses. Resolution uses the immutable profile and verifies complete permission membership, root membership, generation and revision; empty/wildcard roots do not become valid. No caller-chosen namespace, artifact rewriting, extra table, authority search, or relaxed tuple comparison is introduced. The published protocol string remains v1 because no version-1 protocol implementation or stores have been produced. The contract-document version advances to 2 and preserves its rejected predecessor.
