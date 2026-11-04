# Phase 3B.4 candidate request-binding contract

**Status: independently APPROVED contract, revision 2 (2026-11-04); candidate implementation unqualified.** Revision 1 was rejected for conflating Store metadata semantics with supported program semantics; the [exact rejected revision and independent review](evidence/reviews/task1-r1/manifest.json) are retained. This correction requires scoped re-review. This is the exact candidate to attack and independently review before production or principal oracle construction. It is not an implemented guarantee. The user [packet](00-brief.md), [v1 contract](../phase3b3/02-contract.md) and [genuine construction gate](01-charter.md) bound the proposal. Rejected alternatives and preflight reports remain retained.

## 1. Selected representation and domain

Select representation 1: derive the exact original selection from the accepted grant's existing `DelegationClaim.OriginalRequestID`. The service already fixes `AuthorityUse.GrantID`; that immutable grant fixes q and exact e. Derive o from q's logical ID, load its consumption coordinate c, then the real same-generation history row. Existing `EffectAuthority.RequestObservationID` and `RequestConsumedSequence` record the result. No additional public selection field, binding operation, authority table or mutable last-request register is required.

Representation 2 (caller-supplied retained selection reference) redundantly encodes a uniquely derivable fact and adds stale/forged-reference and construction obligations. It does not remove admission validation. Representation 3 (new binding operation/object) adds transitions, state, sequence-locator shifts and accounting vocabulary without a demonstrated expressive gap. Either could become relevant in another semantics allowing multiple independently selected observations of one logical message; this runtime has one deterministic logical observation consumed at most once. These alternatives are rejected for this slice, not universally impossible.

Preserve two initial generation-one logical principals, A servicing/B serial, one-hop grants, original fixed roots and one trusted Store owner. No templates/dynamic members, reconciliation, continuations or legacy dispatch-admission feature combination is newly enabled. Existing bounds, no-fallback, root debit, immutable first decision and source closure rules remain. Historical binding does not prove the payload was computed from q or satisfies q's business objective.

## 2. Exact version and wire boundary

At inspected runtime base `1fdd4c9e1720bf60b344914bcbc3d0938ba8d678`, constants allocate physical formats 3–14 only. Reserve candidate 15/16, subject to this review and recheck before production:

| Physical format | RequestAuthority profile | Meaning | Conflict promotion |
|---|---|---|---|
| 13 / 14 | `alan/request-authority/v1` | unchanged immediate predecessor policy | 13→14 |
| 15 / 16 | `alan/request-authority/v2` | exact earlier selection policy below | 15→16 |

```go
const RequestAuthorityVersion = "alan/request-authority/v1" // unchanged
const RequestAuthorityBindingVersion = "alan/request-authority/v2" // new public opt-in
```

`CompositionProvision.RequestAuthority` must equal the exact version mapped by metadata schema and `PRAGMA user_version`. All four formats retain composition protocol `alan/composition/v1` / protocol schema 4. The fixed Store `metadata.semantics` value remains `SemanticVersion = "alan/1d/v1"`; this metadata value is not an exhaustive restriction on per-artifact program/reducer semantics. Formats 15/16 do not imply continuation/reconciliation/legacy dispatch features by numeric ordering. Receipt conflicts are present only in 14/16; promotion preserves the authority interpretation.

Preserve the existing per-artifact supported semantics and validation rules: `alan/1d/v1` (raw program), `alan/1e/v1` (authored program with `alan/author/v1` metadata), and `alan/1g/v1` (collection/authored program with `alan/author/v2` metadata). `supportedSemantics`, artifact byte/hash and retained-semantics agreement, `validateMetadata`, instruction/collection validation, authority-profile validation and existing deployment/feature restrictions remain applicable. Neither formats 15/16 nor this compatibility list enables an otherwise rejected feature combination, invocation/template configuration or authority selector. The new public `RequestBinding` example is narrower: its artifacts use `alan/1d/v1`; it does not narrow the runtime's accepted artifact semantics.

Compatibility qualification must preserve the existing authored-adapter acceptance case `TestPhase3B3CartesianOperations/r1-add-one` in `internal/durabletest/phase3b3_adapter_test.go`, which constructs A through `author.Compile` with artifact semantics `author.SemanticVersion = "alan/1e/v1"` under the v1 authority profile. Run the unchanged v1 case in normal and race qualification and verify Create, inspection, Replay and Open accept its genuine retained artifacts and history under formats 13/14, including the format-14 receipt-conflict variant when promoted through the existing writer. For a newly provisioned v2 counterpart, keep the authored semantics/metadata, exact adapter tuple, adjacent request/service order and operation locators; change the Store authority profile and all input authority claim/use versions coherently to v2 before provisioning, rehash changed program bytes, and require the candidate writer to emit v2 decisions/witnesses/gates. Do not edit retained descriptors. Expected v2 disposition: accept the valid authored artifact and authorize the same exact service under formats 15/16, with the genuine result retained/consumed and one endpoint application; do not reject solely because its semantics is `alan/1e/v1`. Invalid metadata/semantics pairings and unsupported feature combinations must still reject. This is a future compatibility requirement, not a new execution result or a claim that every collection program is supported by the two-principal fixture.

No fields are added to `DelegationClaim`, `RevocationClaim`, `AuthorityUse`, `AuthorityDecision`, `EffectAuthority` or `AuthorityAdmission`. Every interpreted claim/revoke, retained use, lifecycle authority decision, effect witness, and physical gate descriptor must carry the exact Store authority version. Semantic version equality is checked at every interpretation boundary, not just 'feature enabled'. Literal program control metadata must match the chosen profile before provisioning; lifecycle claim values must match when interpreted as authority operations, including computed payloads. Ordinary uninterpreted application strings are not grants. Unknown/mixed authority versions fail closed with a version error; retained contradictory descriptors also fail ledger validation.

Unchanged structural composition fact ID/digest domain strings such as `alan/fact/authority-admission/v1` and `alan/authority-admission/v1` remain envelope encoding versions, not authority-policy versions. The canonical typed descriptor includes `Version=v2` in its digest. Do not globally replace every `/v1` string. All old field omission/recognition behavior and original version interpretation remain intact.

No in-place migration from 13/14 to 15/16. New behavior requires a newly provisioned v2 Store, never reinterpretation of an old refusal. New readers validate/replay/open old histories with old semantics and unchanged witness bytes; old readers reject 15/16 before owner renewal. Test metadata/PRAGMA-only, profile-only, nested claim/use/decision/witness/gate and coherent multi-field downgrades. A v2 nonadjacent allow rehashed as v1 must fail old-policy validation. Full unsigned policy-consistent history rewriting is outside authenticity guarantees; hashes enforce consistency, not proof of an untampered past.

## 3. Exact causal predicate

Let q be the original request intent; g the accepted `DelegateAuthority` lifecycle result selected by e's `AuthorityUse.GrantID`; o the accepted ordinary-message observation of q; c the selecting reduction; e the actual external or adapter service intent; P the authority decision prefix. Prefix means existing ordered composition facts plus their validated causal source records, not final-state presence.

The inherited checks run in their existing order: original effect commitment/source authorization, selected root validity and complete action tuple, own-path owner check or delegated grant existence/acceptance, exact service/recipient/root/revision comparison. Only the delegated selection subpredicate changes under v2:

1. Require g's lifecycle decision fact in P before this decision; authenticate g from its original committed lifecycle intent, real issuer/source revision and issuance policy. g may be accepted before or after c. Original q must precede issuance; service commitment may precede grant acceptance provided its first decision has not happened and all decision premises hold. An earlier refused decision remains refused.
2. Load q by `intents.id = g.OriginalRequestID`, never by payload, peer, most recent message or obligation label. Validate canonical ID/index/descriptor equality, actual committing reduction, ordinary message kind, issuer logical Run and revision, exact target recipient Run, source permission and source authorization. The grant holder/root, recipient/root/revision and exact full permission row must match g, e's selector, original deployment and current decision policy. Resolve root locators through namespace and actual membership.
3. Derive `o.ID = identity("observation", "message", q.ID)` using the unchanged canonical identity encoding. Load observations by that primary key. Require exact descriptor/index ID, subject/generation, recipient Run, `ordinary-message` provenance, issuer source/generation and content equal to the authentic q-derived descriptor. Its correlation `Operation` field alone is not original identity.
4. Require a real consumed sequence c satisfying **`0 < c < e.Sequence`** in e's same logical Run. c is compared with service commitment, not decision time. Load `history` by `(e.Run.Subject, e.Run.Generation, c)`; require `selected_observation == o.ID`, matching observation consumption and the true validated program reduction. Its source revision must equal the preserved recipient revision. Validate its execution fact before e's execution fact and both before P; validate admission and source authorization as causal predecessors of c. No current/future history row may justify a past prefix.
5. Join genuine accepted admission(s) by o.ID. Each supplying admission must match q's exact committed message descriptor using its **actual** transport ID, trusted domain, logical ID and original intent; `Accepted=true`, bound recipient exactly e.Run and `ObservationID=o.ID`. Recompute `observationFor` and indexed/JSON identities. Public transport IDs may be any valid nonempty UTF-8; do not assume the automatic runner's canonical transport name. Multiple valid transport attempts can share one logical observation. ACK is separate and is not an admission/selection premise.
6. Require q's real allowed source authorization, route tuple and causal prefix validity, not a final-state permission invented after selection. Keep e ordinal zero, g.ServiceOperationID exactly e.ID, holder/recipient generations/revisions, both roots, selected path and complete `(Kind, Target, Operation, Obligation)` attenuation. External operation is `apply`; adapter operation is its obligation, as in v1. No Cartesian mixing, grant from a result, ambient borrowing, or path fallback.
7. Set existing witness `OriginalRequestID=q.ID`, `ServiceOperationID=e.ID`, `GrantDecisionSequence=g.DecisionSequence`; on binding success set `RequestObservationID=o.ID` and `RequestConsumedSequence=c`. Namespace, selected Use, RootID and ServingRootID follow the unchanged derivation. Permission allowance is separate from resource reservation allowance.

A future authentic selection (`c >= e.Sequence`) fails even when committed before a delayed decision. Replaying that refused decision after yet more computation must retain the same reason and witness. A real selection in another generation cannot satisfy this predicate. Prefix validation must not infer that selection occurred merely from endpoint output, application registers or an observation admitted but unconsumed.

## 4. Absence, false claims and corruption

Keep v1 reason vocabulary/precedence exactly unchanged. For v2, structural and historical validation first distinguishes valid incomplete prefixes from contradiction:

| Input at the relevant prefix | Candidate disposition |
|---|---|
| No selector | retained `authority_missing` refusal |
| Selected grant absent from P | retained `authority_grant_missing` refusal |
| Existing selected lifecycle result refused/not a grant | inherited `authority_grant_refused` refusal |
| Valid selected grant but wrong exact service/recipient/root/revision | inherited `authority_service_mismatch` refusal |
| q not yet admitted, or admitted but not selected before e; q2 only selected; genuine c≥e | retained `authority_request_not_selected` refusal |
| q's claimed selection exists but observation/history/consumption/descriptors disagree, lacks genuine admission/source authorization, or violates causal order | corruption error; never retain a successful denial as a repair |
| Syntactically invalid/unknown/mixed authority versions | validation/version error, no authority transition |
| Valid historical binding but currently revoked/stale root/source/closed scope or exhausted budget | inherited authority/source/resource refusal at its actual boundary |

For an absent or future selection, leave request-observation/consumption witness fields empty/zero. Populate grant-bound original/service/decision coordinates only after authentic matching grant checks as in v1. Later delivery/selection cannot change a historical refusal into `authority_request_not_admitted` or fill its former empty witness. For a valid ledger an actual c<e selection already has its accepted admission and route authorization; their disappearance is corruption, not natural absence. Failures in general retained-ledger integrity are not hidden by an earlier policy refusal.

Malformed public claims are rejected by input validation; well-formed false claims retain the existing issuance/service refusal semantics. Never classify a malformed archive, checker panic, missing capture or unsupported projection as a valid policy denial.

## 5. Transactions, current admission and recovery

Reuse the owned SQLite transactions. Admission commits observation and accepted admission atomically. Selection commits consumption, history, successor, emitted intents and execution fact atomically. Lifecycle decision/result and its fact are atomic. First attempt/witness, reservation, debit and chronological fact retain their existing atomicity. New code adds no authoritative cache or mutable binding record.

A new physical opportunity reconstructs the immutable original operation and attempt, original allowed decision and debit, then rechecks the **same** selected grant/root with current holder/recipient generation/revision/membership, source openness, scope, revocation and owner epoch. Attempt/reservation agreement alone cannot authorize a forged operation. Existing dispatch locks and durable gate-before-worker-release remain. Epoch renewal is physical; it does not change q/o/c/g/e logical IDs, mint a grant or authorize a new contact from an old volatile gate.

Inspection, Replay, first decision, fresh gate and Open use one version-selected predicate and historical-prefix interpretation. Open validates before owner renewal; failed Open must leave immutable bytes and epochs unchanged. Replay performs no external contact. Agreement among shared helpers is internal consistency evidence, not independent policy proof.

Retain normal delegator completion policy: completion alone does not revoke. Retirement, generation/revision change, scope/source closure and explicit monotone revocation retain their original effects. There is no new wall-clock TTL/lease expiry; 'expiry' controls mean these existing invalidation boundaries. Retained results remain consumable after authority changes; a prior admitted opportunity retains its original meaning. Neither licenses another physical action. One grant binds one original service, not necessarily one idempotent physical contact. Generic unknown effects are not retransmitted automatically.

## 6. Exact admission index and cost guard

Add only a v2-format derived index, created at 15 provisioning and preserved by 15→16:

```sql
CREATE INDEX admissions_observation_id ON admissions(observation_id);
```

This is non-unique: distinct transports may admit the same logical observation. Query `SELECT id FROM admissions WHERE observation_id=? ORDER BY id`, then validate the exact descriptor/result joins above. Observation and history lookups already have primary keys. V2 readers must validate that this named index has the intended single column/non-unique shape; a missing/incorrect index fails format validation rather than silently adding it during Open. Formats 13/14 are untouched. Index entries are derived SQLite state, never authority. No table, last-request cache or scan for a convenient historical message is permitted.

The new lookup must not add an all-history/all-admissions scan. Existing `readAuthorityPrefix`/`authorityFactOrder` composition-prefix scans and projection reconstruction remain; no total O(1) claim is made. Public arbitrary transport IDs prevent deriving a unique admission key from q, and message ACKs may be absent. That is why this one index is selected over a synthetic transport-key shortcut.

Before execution, freeze cost lengths **0, 1, 8, 32**, same toolchain/machine, v2 profile, fixed endpoint/action shape and matching actual service grants. Measure program construction separately; first authority decision; fresh gate; Replay; Open on disposable copies. Record elapsed time, allocations where an available harness exposes them, executed reductions, bytes, admissions count and query-plan evidence showing exact lookup uses this index. Use repeated observations under the same source and document variation, without inventing a threshold or speedup. Compare added lookup separately from inherited prefix work; this is a regression guard, not another optimization campaign. No cost results are claimed in Task 1.

## 7. Frozen public example interface

Add in `examples/delegation/request_binding.go`, using public packages and ordinary retained programs only:

```go
type RequestBindingOptions struct {
    EndpointA, EndpointB string
    LocalSteps int
    ForwardBranch *bool
    ReadyFirst bool
    TerminalHandshake bool
}
type RequestBindingManifest struct {
    Version string // "alan/request-binding-fixture/v1"
    Options RequestBindingOptions
    Reciprocal Manifest
    ServiceSequence, ReplySequence int64
    InsertedReductions int
    ProgramIDs map[string]string // subject -> actual artifact SHA256
    PrincipalRequests [2]string // AB then BA, exact original IDs
}
func RequestBinding(o RequestBindingOptions) (runtime.Provision, RequestBindingManifest, error)
```

The new constructor always opts into v2/S2-handshake. Existing `Reciprocal` and its output remain v1, unchanged. Default `ReadyFirst=false` selects q then ready; true selects ready then q as the adjacent control when no insertion. `LocalSteps` accepts 0..32; negative/larger values fail. Non-nil `ForwardBranch` selects the demonstrated true/false finite path and requires LocalSteps=0. Branch executes three reductions across five code positions; scratch writes are existing set literals and set/Ref copies. Other lengths use alternating literal set/copy in one fixed scratch field. The builder changes every claim copy, exact future service/reply locator, downstream PC and program hash before Create. The one-hop grant remains an ordinary B lifecycle instruction; no host-issued grant shortcut.

Manifest inner Options is S2Handshake with the supplied endpoints/terminal flag. Service sequence is 4+InsertedReductions and planned successful A reply sequence is 7+InsertedReductions. All original full IDs remain exposed. A forward branch cannot derive an operation ID from PC+1. Unsupported unequal-length/dynamic service paths require a genuine actual-locator grant or report an authoring obstruction; this constructor does not claim general dynamic claim synthesis.

Within one execution namespace/programs/logical IDs are conserved. Across separately provisioned programs provide an explicit role correspondence, not byte equivalence: AB request→AB request, BA request→BA request, grant→grant, ready→ready, A4→A(4+n), A7→A(7+n), unchanged B9/B12; both actual programs and namespaces remain visible. Bind the two principal caller/callee/service roles independently against committed facts. Equal-payload auxiliary q2 is separate from the principal pair and never inferred by obligation suffix.

## 8. Prospective program, cut and attack matrix

These are predictions, not executed v2 results. Task 2 refines the relation and exact bounded model before treating any runtime run as confirmation. Use a covering family, not a Cartesian product.

| ID | Genuine program/schedule | Expected v2 result / discriminator |
|---|---|---|
| P01 | Original adjacent, v1 and v2 | both authorize and complete reciprocal custody |
| P02 | q→ready→e, equal A4 position, v1 and v2 | preserve v1 refusal; v2 authorizes and automatic Run completes |
| P03–P05 | q→ready, local lengths 0/1/3, genuine actual grants | authorize; both real results consumed; separate ACK/terminal/quiescence measures |
| P06–P07 | Both equal-length forward branch outcomes | authorize exact A7; demonstrate actual path |
| P08–P09 | q selection before and after grant issuance, grant before decision | both authorize with authentic prefix |
| P10 | Genuine q then equal-payload/same-obligation q2, g names q/e | permit q binding after later receive/overwritten register; q2 remains auxiliary |
| P11 | After authentic q selection, overwrite application `peer_request` while preserving actual service and reply operands | authorize through retained q, not application register; exact original/result custody still checked |
| N01 | q admitted but unselected, only q2 selected | `authority_request_not_selected`; no exposure |
| N02 | e committed before authentic q selection; decision delayed until after c | `authority_request_not_selected`, stable through Replay/Open |
| N03 | First decision before grant, later genuine grant/reopen | original missing-grant refusal remains sticky |
| N04 | First decision before selection, later genuine selection/reopen | original not-selected refusal remains sticky |
| N05 | Revoke after selection before first decision | historical selection survives; current authority refuses |
| N06 | Allow then revoke before fresh/recovery gate | original allow/debit retained; fresh gate refuses |
| N07 | Genuine gate before revoke; retained result consumed after revoke | old opportunity/result retains meaning; no extra permission |
| N08 | Stale old A4 grant for newly padded A5 | exact service mismatch, not selection failure |
| N09 | Unequal path or branch changes service operands/skips service | mismatch/refusal or declared outside progress; no 'harmless step' label |
| N10 | Exhausted resources/closure/retirement/generation/revision | inherited actual boundary refuses; binding fact not erased |
| N11 | Cartesian tuple mixing, ambient borrowing, result-as-grant, own fallback | previous negative policy preserved |
| N12 | A second distinct service attempts same grant | exact operation binding refuses even if payload equal |
| A01 | Wrong namespace, holder/recipient generation/root/revision/tuple | reject contradiction or retain well-formed false-claim refusal at actual layer |
| A02 | Forged c, wrong o/q/service/selected history; missing admission/source auth | dense coherent attack must reach causal rejection |
| A03 | Future final grant/result used for earlier decision | historical prefix rejects; no final-state laundering |
| A04 | Coherently forged attempt/reservation with original operation mismatch | fresh gate reload rejects original mismatch |
| A05 | Metadata/PRAGMA/profile/nested descriptor downgrade or version mix | reader rejects before renewal; no mutation |
| K01 | Existing authored-adapter `TestPhase3B3CartesianOperations/r1-add-one`, v1 formats 13/14 | preserve artifact acceptance, exact allowed service/result and Create/Inspect/Replay/Open compatibility; normal and race qualification |
| K02 | Newly provisioned coherently versioned v2 counterpart of K01, formats 15/16 | accept `alan/1e/v1` authored artifact; authorize same exact service; preserve metadata/feature rejections, as specified in §2 |

The positive stretch argument permits writes only where their read/write footprint preserves the exact service, selected authority, reply/custody operands and resource/source conditions. P11 separately tests a dead application request register; it is not permission to overwrite a live service input. Both designated principal results may be consumed while auxiliary q2 remains unpaid, even if Run is quiescent; report all three states separately.

A genuinely issued g(q2,e), with authentic earlier q2 selection and all other premises, is not inherently invalid. A substitution attack must violate an actual original binding or retained consistency. Raw policy, independent endpoint and supported-model verdicts remain separate; an unsupported model projection cannot excuse raw false acceptance.

| Cut ID | Acknowledged SIGKILL boundary, then fresh process on original Store | Expected retained fact / recovery assertion |
|---|---|---|
| C01 | q admission before selection | q stays admitted/unselected immediately after Open; only real later commit selects |
| C02 | q selection before ready selection | preserve exact q/o/c; ready then valid service |
| C03 | grant and selection, after intervening work, before e emission | preserve history and actual planned locator; genuine later emission |
| C04 | e emission before first authority decision | preserve exact e, evaluate at later real decision prefix |
| C05 | first allow before fresh gate | immutable allow/debit retained; new owner/current gate required |
| C06 | retained result before consumption | original result retained/consumed; no completed-call reissue |
| C07 | two consecutive owner losses after selection and after work | both actual renewals preserve q/o/c/g/e and old charges |
| C08 | revoked before recovery gate | historical binding/old allow survives, fresh exposure refused |

C03/C07 use length3; C01/C02/C04/C05/C06/C08 use the zero-local-step nonadjacent program. Ordinary automatic positive P02 is mandatory; controlled public-transition schedules and barriers are labeled separately. Do not claim the unchanged scheduler naturally reaches a controlled cut. Keep merely-reserved versus applied-without-local-result states separate; idempotent repeats versus generic unknown non-retransmission remain separate endpoint contracts.

Each interrupted execution carries an ordered capture list: ordinal, predecessor capture, event/cut label, owner PID/identity, acknowledged signal, pre-Open snapshot, post-Open snapshot, epoch transitions and digests. C07 must show **both** renewals, not fixed filename ordering or one first/fresh pair. Missing pre-Open evidence explicitly leaves that renewal unchecked. Every capture retains original IDs, q/o/c/grant/service, authority/charges/pending obligations, programs/profile/roots/permissions, and independent contact/application ledgers. Preserve immutable-prefix rows and single assignment of consumption.

## 9. Independent evidence and qualification obligations

No principal oracle is constructed in Task 1. Task 5's independent raw checker must reconstruct original intent→route→admission→observation→real c→g→exact e→historical first decision→current gate→result/reply custody without runtime authorization/Replay/reducer imports. Its declared principal domain must fix exact AB/BA caller/callee/service IDs and independently check them; q2 remains separately visible. Own-path association needs the actual request/service/reply program chain, not the old obligation-suffix OR shortcut.

Retain arbitrary public transport IDs in candidate raw joins, or explicitly report an automatic-transport-only unsupported checker domain. No invalidity inference follows merely from leaving that declared subset. Freeze source model exports unchanged. Restrict old-model correspondence to justified nonrevoking known-result single-exposure prefixes; whole-Store deaths, separate-grant revocation and applied-unrecorded/repeated-exposure gaps remain explicit.

Before final confirmation, mutate exported evidence **and** checker source. Required mutants: delete strict c<e (dense coherent false-allow from genuine N02 parent); replace exact q lookup with latest same-peer receive/admission (P10 false-refusal and N01 false-acceptance); omit a nonredundant generation join (e.g. requested admission generation versus Bound.Generation); skip current revocation specifically at fresh gate (coherent N06 false-allow). Each applicable mutant must exhibit a semantic wrong verdict detected by assertion, not compilation error, panic or an unrelated earlier guard. Record redundant/inapplicable mutants honestly and replace them with an applicable one preserving the intended hazard.

Exercise initial/pre-grant/pre-selection/pre-attempt/admitted-only/post-result-unconsumed prefixes. Mutation tooling must return bounded explicit applicability/diagnostics, not panic. Preserve the earlier diagnostic failure separately from any narrow tooling repair. Exact source pairs, commands/exits/toolchains and hashes precede final qualifications; later executable fixes invalidate affected evidence. Required final normal/integration/race/build/module/Staticcheck and repository gates, independent contract/model/runtime/recovery/raw/compatibility reviews, manuscript reconciliation and cross-linked draft delivery remain Tasks 2–6. This contract makes no general termination, dataflow, distributed, N-agent or unrestricted refinement claim.
