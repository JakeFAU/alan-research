# Phase 3B.5 model reconnaissance

Read-only source inspection on `codex/phase3b5-grant-revocation`; initial `git status --short` was empty. No repository files edited and no tests or enumeration executed. This is a proposed extension, not a qualified result. Updated after relation reconnaissance identified runtime revocation's exact root scope.

## Question and smallest semantic boundary

Can one separately retained, accepted, request-bound B→A grant retain its lawful historical decision after **revocation of its exact originating authority root**, while every fresh physical admission on that path is refused? Preserve an unrelated B-owned AB service root, original custody/retries, and atomic application-plus-receipt.

Do not call this grant-ID revocation. The relation agent reports runtime `RevocationClaim` contains Version + exact AuthorityRoot, not GrantID; the root includes Run, RootRun and a complete permission tuple. N06/N07's revoked BA-service root differs from B's AB-service root even though they share a composition root. This is a correction to my initial hypothetical grant-ID proposal. That proposal is not recommended: it would introduce an API/semantics the actual bounded runtime does not use.

## Inspected source facts

- `research/authorityservice/model.go:13` `Configurations()` returns exactly 16 configs; `valid` currently accepts only that list. `Config` has two strings.
- `State` (`:186`) embeds unchanged `mutualobligations.State`, two `AuthorityState`s, one `GrantStatus`, one `Revoked` bit and monotone unsafe flags. That bit means B/g1:b revoked, not grant-ID revoked.
- `GrantRecord` (`:122`) retains full original binding, tuple bitset, root/parent, issuer/grantee generations, revision, no-redelegation and issuance `Accepted`. Accepted/rejected issuance is immutable and distinct from present validity.
- `current` (`:244`), `pathCurrent` (`:267`) and `choose` (`:273`) make old Revocable invalidate both AB's B-root path and BA's grant path. This is correct for its frozen one-b-root experiment, but cannot represent N07's surviving distinct AB root.
- `Successors` (`:341`) uses actual base successors. Only Apply is additionally gated. Server crash clears its permit; first decisions remain immutable. Authority steps stutter under base projection.
- Existing LateGrant retains a separate grant after BA commit and allows the losing authorization-before-issuance order, but cannot revoke. Existing Revocable permits root revoke, including before commitment, but issues inline on BA commit and retains rejected inline issuance. A combined separate issuance/revoke scenario also needs distinct complete B tuples to match the intended two runtime root scopes.
- `Event` (`:197`) is comparable and fully serialized. `Event.Revision` is source revocation era 0/1; grant/witness revisions remain historical. The relation agent notes runtime holder revision remains 1 across this event, so model revocation era must not be equated with runtime holder revision.
- `explore.go` reconstructs compact events with `describe`; descriptors must round-trip exactly. BFS successor order determines shortest witnesses and export bytes.
- `oracle.go` independently defines configs/bindings/roots and ancestry. `grantPath` (`:149`) distinguishes current validation from ancestry; `issue` (`:206`) checks issuance-time acceptance; `accept` (`:284`) retains first witness and validates fresh opportunities; `auditGraph` (`:537`) compares independent ledgers at every merged target.
- `cmd/export/main.go` defaults to `-config all` over Configurations and marshals report/state/event fields. Changing old enumeration/order, fields or cuts can change old bytes. Existing deterministic export test checks two current outputs for S2/Missing, not historical bytes.
- `model_test.go:130` tests old root-revoke temporal boundaries, `:189` sticky late refusal, `:204` released custody across both crashes. None covers separate issuance + independently scoped sibling root.
- `phase3b2/01-contract.md` fixes first-decision immutability, admission-time validity of the same path, surviving old permits, crash clearing, distinct custody and atomic Apply+receipt. Its root target is explicitly B's b.
- `phase3b4/04-model.md:78–98` and `08-results.md:45–57` exclude separate grant+revocation from prior correspondence; also excluded are unknown outcomes, repeated exposures, whole-Store death→two-agent crashes and full reverse simulation. New model enumeration alone resolves none of those unrelated exclusions.

## Recommended minimal extension

Prefer **one opt-in scenario inside authorityservice**, with its own extension enumerator and export/audit entry point. Preserve Configurations' old sixteen values/order and old commands' default/all scope exactly. Proposed name `S2/SeparateGrantRevocation` (or root-explicit name chosen by the parent).

Use existing State/Event shapes and `Revoked` bit. Append one new complete permission atom **after all existing Tuple values**, e.g. `ABOwnPermission = (increment,counter-AB,blue)`, with B-owned root `B/g1:c`. In this scenario only, Binding(AB) uses c; Binding(BA) stays b, and the one B→A grant stays bound exactly to BA/b. Existing `current` naturally leaves B/g1:c current when B/g1:b is revoked. Extend root lookup and the oracle's independently declared roots/tuples/bindings. Never modify old numeric values or old configuration bindings. Do not reuse Shared r: that would turn an increment service into a read just to save a tuple.

The new tuple is an **opaque complete action abstraction**, not a claim to preserve runtime endpoint/obligation strings or a permission-algebra isomorphism. The relation sidecar must retain both exact runtime roots, full tuples, operation IDs and revoke fact ID. One fixed mapping sends exact BA root to b and distinct AB own root to c. It must check that the accepted actual revocation targets the b-mapped root; sharing a composition RootID is insufficient.

Transitions:

1. BA Commit carries no inline grant for new scenario; it remains enabled after root revocation.
2. Separate retain-grant is enabled after BA commitment, while issuer B is live, and only when Grant==NoGrant. Accepted iff B/g1:b is live at issuance. After revocation retain a rejected claim; do not suppress ordinary request commitment. No renewal or second grant ID.
3. Existing revoke transition is enabled for this scenario once, by live B, before or after issuance. It changes current validity of exact BA/b root and derived grant, while leaving AB/c root valid. It changes no historical issuance, first decision, permit or custody fact.
4. Current delegated admission requires accepted grant and unrevoked b. Existing own-root selection discovers c for AB. Apply consumes an existing permit without rechecking current validity.
5. Refusal before grant or after revoke stays sticky; recovery/late issuance cannot rewrite it.

Implementation touches likely limited to model config/binding/root/tuple helpers, retain/revoke enablement and rejected separate issuance, independent oracle counterparts, new scenario explorer cuts/tests, new exporter selection. No added state or event fields, new event kind, second revocation bit, grant-ID revocation protocol, or generic config flags are needed. Existing `grant(c,s)` already sets Accepted from Revoked; separate retain branch currently unconditionally sets RetainedGrant and needs scenario-scoped rejected status after revoke. Old LateGrant never reaches Revoked so its bytes/graph stay the same.

Alternative: separate `research/grantrevocation` package keeps frozen source/API mechanically untouched and allows dedicated types. Its cost is duplicated transition/graph/oracle machinery or refactoring frozen sources into generic abstractions. A thin wrapper around old LateGrant successors must still replace AB bindings/authority witnesses and rejected issuance, while describe is private. For a single scenario, opt-in reuse is less machinery if exact historical export compatibility is mandatory. Use separate package only if unchanged source bytes, beyond unchanged behavior/exports, is an explicit requirement.

## Frozen bounds

A and B only; logical generations fixed at 1; A servicing/B serial; exactly AB/BA originals; original complete tuples plus one independently scoped B-owned complete c tuple for the new scenario; one exact no-redelegation B→A grant on BA operation1; one separate issuance attempt accepted/rejected; at most one exact b-root revoke; one recoverable crash per physical agent; two opaque values per original; original one-slot channels and retry cycles. No timers, renewal/regrant, grant-only revoke, independent owners, arbitrary multihop/roots, third originals or non-atomic effect window. Full graph enumeration without cutoff; no predicted counts.

## Invariants and targeted falsifiers

1. **Historical preservation:** first Decision/Path/Refusal and Accepted/rejected issuance never change; root revocation monotone. Lawful allow remains historical after revoke.
2. **Exact scope:** b revoke disables fresh delegated BA paths but preserves AB/c own-root authorization/admission. Existing r roots remain independent in legacy controls. Mutant equating c and b must fail the surviving-AB control.
3. **Fresh admission:** same selected historical allowed path plus current b validity; unique opportunity tied to serving activation; Apply consumes it once.
4. **Old permit:** revocation does not cancel an existing opportunity. Mutant rechecking present validity at Apply must fail the positive permit-before-revoke trace.
5. **Crash:** server crash clears its volatile permit; recovery increments only that agent. New admission after recovery checks current validity, while original accepted grant/decision and durable result remain.
6. **Custody:** revoke/crash cannot manufacture, reinterpret or consume results; no payload becomes authority; completion means both authentic original consumptions.
7. **Independent audit:** old/new config distinction, exact complete tuple, root/parent/identity/generation/original/operation binding, revoke order/actor/revision, rejected issuance, first decisions, permit epochs and authentic result custody checked without operational helpers. Merged graph states need equal independent ledger.

Required trace matrix:

- Commit BA → retain accepted grant → send/admit BA → authorize allow → revoke b → dispatch-refuse: preserve acceptance/witness/allow; no BA permit/Apply. Mutant omitted fresh gate revoke check must fail.
- Same through allow → dispatch-admit → revoke b → Apply: one application, permit consumed, no repeat. Then genuine BA return/record/consume opens B's serial guard; AB authorizes/admit/applies under c and both originals consume. This is the decisive N07-shaped full positive control.
- Accepted grant/allow/permit → revoke b → crash A → recover A → fresh refusal. No old epoch permit, admission or Apply. Unrevoked crash/recover control must readmit at activation1.
- Accepted grant → revoke b → first authorization refusal; no second authorization. Also authorize before grant → retain accepted grant → unchanged refusal across recovery. Preserve losing schedules rather than force universal completion.
- Revoke b before BA commit and separate issuance: commit retained, issuance rejected, no delegated authorization. Accepted issuance before revoke remains accepted even though no fresh opportunity is possible.
- After b revoke, AB may authorize/admit through c, but AB Apply still obeys B's unchanged serial guard. Do not mislabel inability to Apply before BA consumption as authority failure.
- Apply BA → revoke b → deliver/record/consume; also revoke after record before consume and after release before both crashes. Results stay payable; grant validity is not a condition for consuming already authentic results.
- Revoke from initial state, at committed-before-issuance, accepted-before-request-admit, after allow/before permit, after permit, after apply and after consume. Actor must be live issuer B; duplicate revoke disabled. Separate rejected claim remains immutable.
- Hostile traces: swapped b/c target mapping; mutated complete tuple/root/parent; bad issuer/grantee generation; wrong original/operation; forged accepted issuance after revoke; stale event revocation era; fresh admit after revoke; reauthorization of original refusal; resurrected crash permit; result-as-authority. Preserve unrelated authentic application outcomes in mutations.

Apparatus: exact old sixteen list/order and archived bytes; new deterministic exporter twice; all new graph edges independently checked; exact compact event reconstruction; all new temporal witness cuts have expected target original (do not let AB/c Apply after b revoke accidentally stand in for the intended BA/pre-admitted-permit witness). The old explorer's generic `application-after-revocation` cut can select either obligation; add scenario-specific cut names/conditions only for the new scenario. Keep expected completion counterexamples, SCC/retry cycles, safe-but-noncompletable states and honest exit1/exit2 distinctions.

## Qualification and remaining boundary

Use pinned Go (go.mod says1.27.3), offline/local toolchain, targeted tests first, affected tests/race and Staticcheck, full new graph and independent all-edge audit, then old sixteen and exact historical export regressions. No suites were run here.

Prior four historical commands/outputs are indexed at `docs/foundational/phase3b4/evidence/task5/apparatus/export-twice.py`, `phase3b4/07-runbook.md:89`, `phase3b4/evidence/task5/report.md:134`: twoagent export, mutualobligations export, authorityservice export defaults, authorityservice audit. Build with `-buildvcs=false`, run twice into new files and compare with archived bytes (decode Phase3B.2 gzip). Current-vs-current determinism alone is inadequate. Do not overwrite archives.

Document the new scope and reconcile affected Markdown/LaTeX manuscript claims. No new reference appears necessary. Keep Phase3B.4's frozen conclusions historically attributed. New model scope can support a newly checked narrow N06/N07 trace mapping; it is not general refinement, permission-algebra equivalence, all-state completion, unknown-outcome recovery or a whole-Store/independent-agent crash correspondence. Exact relation-agent evidence reports N06 accepted grant #8 → first allow #16 → exact root revoke #20 → fresh refusal #22, and N07 accepted grant #8 → allow #16 → permit #17 → revoke #21 → BA result #23 → AB own-root allow/gate/result #33–35. Those fact-order summaries are communicated findings, not independently re-executed in this reconnaissance.
