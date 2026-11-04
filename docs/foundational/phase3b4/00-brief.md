ALAN Research Phase 3B.4
Explicit Request Binding Across Intervening Computation and Recovery
Assignment
Goal: Determine whether the two-agent runtime can preserve a delegated service's explicit binding to an authentic earlier request selection across permitted intervening computation and process recovery, without weakening complete-tuple attenuation, exact operation binding, or current admission checks.

Architecture: Retain the Phase 3B.3 two-agent, one-hop, same-Store experiment. Replace instruction adjacency only in a separately identified candidate contract with a verifiable relationship to the original request-selection event. Prefer existing retained facts over new state. Keep the old contract, model exports, refusals and counterexamples intact.

Tech stack: Existing Go runtime, public author/runtime/effects APIs, SQLite, controlled endpoint, subprocess barriers, independent model and raw-evidence checkers. Discover actual versions before running commands. No new dependencies, paid services, distributed ownership, multihop delegation or additional agents.

Specification: This assignment, the verified Phase 3B.3 packet, and the contract frozen in Task 1. Read each repository's governing instructions. Execute through independently reviewed draft delivery, not merely a plan or first passing example. Do not merge or release. Ordinary decisions within scope are delegated; preserve and explain genuine obstructions rather than manufacturing success.
1. Starting evidence and exact question
The research lives in github.com/JakeFAU/alan-research; production-runtime changes belong in github.com/JakeFAU/alan.
Anchor	Reviewed identity
Research PR #44 delivery	d365c7e8fbc2985b998936cc162f322f6a8f2367
Research executable qualification	7894bfae268e6dc0f5418bda3c6022f3a8a55540
Runtime PR #10 reported delivery	91963d5a81e4ce16f424aae36cc4fb382a56e2cd
Runtime executable qualification	759df6e3d2505f27e1287b2cfeb07bc4290afc0a


These are historical anchors, not assertions about current main or later merge contents. Fetch both repositories, verify ancestry and PR status, and inspect the actual current code before choosing bases. Use isolated worktrees; do not mix this work with bounded-series performance work.

Read the Phase 3B.3 contract, results, correspondence checker, runtime qualification and especially:
- docs/foundational/phase3b3/02-contract.md
- docs/foundational/phase3b3/04-results.md
- docs/foundational/phase3b3/evidence/selection-narrowing/review.md
- research/delegationcorrespondence/
- Runtime internal/executor/request_authority.go
- Runtime examples/delegation/reciprocal.go
- Runtime internal/durabletest/phase3b3_controls_test.go
- The current public inspection, Replay, Open, dispatch and protocol-version paths they use.
The preserved obstruction
Under alan/request-authority/v1, the serving effect must immediately follow the committed reduction that selected its grant's original request. The runtime checks history[service.Sequence-1]. It already records the request observation and consumption sequence in the authority witness.

Phase 3B.3 preserves this genuine history:
A selects the authentic original request
B's exact delegated grant is retained
A selects the grant-ready message
A emits the exact service operation named in the grant
First service decision refuses: authority_request_not_admitted
The request was admitted and selected; the refusal reflects the immediate-predecessor rule, not absence of the request. The unchanged LateGrant model permits authorization at the corresponding admitted-request/retained-grant prefix. The supported runtime ordering selects ready first, then the original request, then emits the service.

This run asks whether explicit historical binding can replace that adjacency requirement in a declared extension. It does not declare the old implementation incorrect under its old contract. An unchanged v1 store must keep its original meaning.
2. Nonnegotiable distinctions
Selection provenance is not current authority. The fact that A selected request q can remain true after B revokes permission. That historical fact must survive while a fresh dispatch is refused.

A retained grant is not physical admission. Recovery may preserve the logical grant and original decision/debit, but a new physical opportunity still requires the current-owner gate and the same selected authority path.

A bound request is not arbitrary dataflow proof. This slice establishes that an explicitly identified service is permitted for an authentic original request. It does not prove the service payload was mathematically computed from that request or that the service satisfies its business objective.

Completion reachability is not inevitable termination. A safe completion path, an actual automatic run, result consumption, ACK completion and whole-deployment quiescence are separate results.

Instruction insertion changes identities. Existing operation IDs include committed reduction sequence. Inserting a step can change the exact future service ID even when its payload is unchanged. Do not mistake a stale service locator for a failed request-binding extension, or normalize that difference away.
Task 1: Reproduce the obstruction and freeze the candidate contract
Baseline and construction gate
Reproduce the original supported and refused orderings on the selected baseline using genuine ordinary runtime operations. Preserve the final authority decision, original selection, grant, service identity, unpaid obligations and independent endpoint contacts. Keep the prior model and its audit unchanged.

Use the equal-length request/ready swap as the first discriminating pair: it preserves the service's sequence position and isolates adjacency. For step-padding experiments, each newly provisioned program must obtain a genuine grant for its actual future service operation. Never edit a retained grant, program or history to repair a changed locator. If public construction cannot name that future operation without unsupported inspection or fabrication, record the separate authoring obstruction.
Choose the smallest representation
Compare these approaches before freezing the design:
1. Derive the exact earlier selection from the grant's existing OriginalRequestID, accepted admission, recipient observation and consumption/history records.
2. Require an explicit selection reference in retained service metadata and validate it against those same original records.
3. Add a new binding operation or retained object only if the first two demonstrably cannot express or validate the contract.
Prefer the first sufficient approach. Do not create a new capability registry or scan for any convenient prior message from the grantor. The lookup must start from the exact original request named by the accepted grant. A caller-supplied reference is evidence to verify, not authority to trust.

Freeze the wire/public interface, interpretation of missing and invalid bindings, transaction boundaries, compatibility policy and test matrix before production implementation or principal oracle construction. Obtain an independent contract/construction review. Keep rejected alternatives and their reasons.
Required causal predicate
Use names equivalent to q for the original request, g for the retained grant, o for the admitted request observation, c for its selecting reduction and e for the actual service. The candidate must establish:
q was genuinely committed and legitimately admitted to A's logical generation;
c is A's real committed selection of o, and o is the observation of q;
c precedes e's committed reduction in that same logical generation;
g was accepted before the relevant authority decision and names q and exact e;
g's holder, recipient, roots, revisions and complete permission tuple match;
the required authority is valid at the decision and every fresh physical admission.
Define the corresponding historical-prefix checks for Replay and inspection. A request selected only after service commitment cannot retroactively authorize that service. A current grant or result cannot justify a historically earlier decision.

Retain existing rules for original refusals, root validity, normal delegator completion, retirement, generation/revision changes, revocation, root debits and no authority-path fallback. Do not change these policies while studying adjacency.

Task exit: a genuine baseline witness, a publicly constructible candidate program, an exact reviewed contract, a compatibility decision and prospective falsifiers. No implementation success is claimed yet.
Task 2: Define permitted intervening work and the correspondence claim
The principal positive domain remains two initial logical agents, A servicing/B serial, one-hop grants and one trusted Store owner. Preserve the original reciprocal obligations and endpoint contract.

Begin with finite stretches of local set/copy and forward branch instructions, grant-readiness selection, and same-Store owner loss/reopen. These are candidate permitted steps, not automatically harmless transitions. A branch that changes the granted service, a receive that supplies different authority, resource exhaustion, closure or revocation can change eligibility.

Freeze which facts a permitted stretch preserves: original request selection, logical identities, grant and service binding, permission tuple, program identity, adequate resource headroom and source eligibility. Owner reopening may change physical epochs; it must not change logical recipient identity or mint a grant. Another selected message must not silently replace the service's explicitly chosen original.
Three separate obligations
Historical binding persistence: Establish whether the authentic selection-to-service relationship survives permitted committed extensions and reopening. A crash before selection commit is not a retained selection. Contradictory or missing evidence must not be reconstructed from the desired outcome.

Authority preservation: Under the additional current-validity premises, establish that permitted intervening work does not by itself invalidate a correctly bound service or expand its authority. Keep revocation outside the positive invariance premise and inside negative testing.

Restricted progress/correspondence: Identify the precise prefixes and event vocabulary for which the candidate restores the formerly refused model authorization and supports both authentic result consumptions. Do not claim all model histories are implemented merely because this one obstruction is removed.

Use a source-level argument where possible and separately label finite enumeration and concrete executions. Preserve all old deterministic exports byte-for-byte. A new auxiliary model may add explicit selection and intervening steps; relate it to the old model rather than silently modifying the old acceptance rule.

Keep any extra same-peer request used for binding attacks explicitly separate from the two original service obligations; it is not a third principal or an undeclared expansion of the progress theorem.

Bound the new state exploration prospectively. Small families with zero, one and multiple intervening steps, alternative readiness orderings and up to two owner losses are preferable to an uncontrolled Cartesian expansion. An induction over declared finite stretches may support a broader statement; enumeration alone does not.

A projection may hide an internal step only after its effect on observations, eligibility, resources and divergence is justified. Keep visible costs in a sidecar or the relation itself. Do not equate whole-Store process death with independent A/B crashes without an explicit argument. Where that map is unavailable, report concrete recovery correctness separately from model correspondence.

Task exit: an independently reviewed relation, preserved baseline mismatch, finite bounds, predicted verdicts and mutation controls before treating runtime results as confirmation.
Task 3: Implement the minimal versioned runtime extension
Implement in alan; keep independent models, scientific arguments and raw auditors in alan-research. Reuse journals, authority machinery, resource accounting and dispatch lifetimes. Do not introduce an agent-specific scheduler, host-issued grant shortcut or second budget ledger.

Write a failing test for each changed behavior, confirm the failure is at the intended causal layer, then make the smallest correction. Commit independently testable increments. The new public example must use public packages and ordinary retained programs. No direct SQL, internal-import escape hatch or host-injected consumption may establish a positive history.
Versioning is part of correctness
The candidate admits histories v1 deliberately refuses. Give that semantic difference a retained, unambiguous version/profile boundary. Preserve formats 13/14 and alan/request-authority/v1 under their original interpretation. Discover currently allocated version numbers; do not assume the next numbers are free.

Do not add a table merely because the semantics changes, but do not avoid necessary versioning merely because the physical schema can be reused. Verify old readers reject unsupported new semantics, new readers preserve old histories, and profile/schema/descriptor downgrades cannot hide the new behavior. Do not reinterpret a previously retained refusal as an allow. Newly enabled behavior requires a newly provisioned supported execution, not mutation of an old refused original.
Durable validation and admission
The recorded witness must identify the actual request observation and selecting reduction. It must survive further receives, overwritten application registers and reopening without relying on a mutable “last request” register or an in-memory cache. Whether metadata explicitly carries that reference or the runtime derives it is the frozen design decision.

Require the same interpretation in first authority decisions, fresh admissions, historical validation, typed inspection, Replay and Open. Replay remains effect-free. Open validates before renewing authority and retains the existing failed-open guarantees. At a fresh gate, reconstruct the original operation and selected path; do not trust a mutually consistent but forged attempt/reservation pair.

Keep one grant bound to one original service operation. This is not automatically one physical contact: any existing idempotent recovery entitlement remains separately governed and freshly admitted. A second distinct service cannot reuse the grant merely because its payload matches. Generic unknown mutations remain unresolved without automatic retransmission.

Retained results remain consumable under the existing rules after authority changes; that is not permission for another physical action. A new grant after a sticky original refusal cannot rehabilitate that refused operation.

Task exit: minimal implementation, focused RED/GREEN tests, public construction, explicit historical compatibility and no unreviewed changes to the surrounding authority policies.
Task 4: Run the discriminating program and recovery families
Use genuine writers for positive cases and isolated disposable copies for storage attacks. The following families are required; leaf counts are outcomes, not targets.
Family	Required discrimination
Original adjacent sequence	Remains valid under its original contract and the declared new profile.
Request → ready → exact service	V1 retains its refusal; the candidate authorizes only with verified earlier selection and all other premises. Complete the reciprocal result path, not just authorization.
Local bookkeeping and forward branches	Multiple finite intervening steps preserve binding when the actual service, authority and adequate resources remain as declared. Use genuine matching grants for changed service locators.
Two genuine requests from B	Equal payloads do not identify a request. Intervening selection cannot overwrite the explicit binding to q; substituting the other request or grant is rejected.
Delivery without selection	Authentic admission alone does not satisfy the candidate's selection contract. Forged or future selection cannot supply it.
Grant/readiness order	Selection before grant issuance and selection after issuance are tested where the frozen contract allows them. Service decision before grant remains refused.
Revocation and expiry	Revocation after selection but before a new gate blocks admission; a previously admitted opportunity and retained result keep their actual meanings.
Replay/substitution attacks	Wrong namespace, generation, root, revision, selection sequence, observation, original request, service ID or tuple cannot be laundered into authority.
Cartesian and ambient borrowing	Prior negative boundaries remain: no mixed permission dimensions, no authority from result possession, no fallback to a different authority path.
Completion and waiting	Authentic result custody, both result consumptions, ACK completion, terminal principals and quiescence are measured separately. Safe refusal is not completion.


Genuine process cuts
Use acknowledged SIGKILL and fresh processes on the original Store. Cover request admission before selection; selection before readiness; retained grant and selection after intervening work but before service emission; service emission before its first decision; allowed first decision before a fresh admission; and result retention before consumption. Include two consecutive owner losses and a revoked-before-recovery-gate control.

At each cut record original IDs, selection sequence, grant, service locator, epochs, decision, charges, pending obligations and endpoint contacts. Reopening must neither consume an unselected request nor bind a different one. Retained response recovery must not reissue a completed call. Keep applied-without-local-result and merely-reserved uncertainty distinct.

The whole-process crash does not establish that either logical agent terminated. Separately label automatic Run executions and controlled public-transition schedules. A barrier exposes a cut; it does not prove the unchanged scheduler naturally reaches every proposed ordering. Require an automatic positive example for the ordinary new program, alongside controlled adverse schedules.
Identity comparisons
Within each interrupted execution, preserve original IDs exactly. Across separately provisioned padded programs, declare any correspondence bijection explicitly and preserve full original IDs in the records. Do not call those byte-identical histories. Changed instruction sequence, authority metadata and protocol version are observable differences, not removable noise.

Task exit: the earlier mismatch is either resolved in the explicitly expanded domain with ordinary continuation, or retained as a precisely explained obstruction. Negative controls remain meaningful.
Task 5: Independently check the causal evidence and qualify frozen source
Extend the independent raw-policy checker without importing runtime authorization, Replay or reducer helpers. It must join the original message, route permission, admission, observation, selecting history row, grant, exact effect, first decision, gate and resulting custody. Reconstruct validity at the relevant prefix, not from final-state presence alone.

Keep three verdicts distinct: raw retained-policy validity, physical endpoint evidence and supported model projection. An unsupported projection cannot excuse a raw authority violation. Conversely, valid raw recovery is not full model refinement. Maintain cross-snapshot conservation of immutable program/authority facts and real owner-renewal evidence, including explicit missing-capture limitations.

Mutate both exported evidence and checker code. At minimum, force false acceptance by deleting the selection-before-service check, using the latest receive instead of the bound one, omitting a generation comparison, and skipping current revocation at admission; the qualification must detect each applicable mutant. Preserve dense, coherently rehashed attacks rather than relying only on gaps or broken hashes. Test incomplete prefixes: a checker crash is an apparatus failure, not a valid rejection. Repair the known mutation-generator diagnostic only where needed for these inputs and record that tooling change separately.

Run affected ordinary, integration and race suites in both repositories; full ordinary suites in both; full runtime integration/race where required by the touched shared validators; prior delegation, custody, causal-recovery and compatibility regressions; and applicable build, vet, formatting, module and Staticcheck gates. Use actual repository commands and documented watchdogs. Record optional export skips and inherited findings separately from new failures.

Freeze exact executable source pairs before final qualification. Record hashes before/after, toolchains, command arguments, exits and evidence digests. A later source fix invalidates affected qualifications; documentation-only delivery commits must be identified separately. Cross-linked PRs must identify the tested pair, and any prospective combined tree with concurrently advancing main needs its own applicable checks.

Check construction/validation cost for a few predeclared intervening-prefix lengths. Avoid introducing a full-history scan per authority check where exact indexed identities suffice. This is a regression guard, not authorization for another performance optimization campaign.

Obtain independent contract/model, runtime/recovery, raw-evidence/compatibility and whole-change reviews. Reviewers must inspect actual causal records and challenge false acceptance as well as false refusal. Preserve rejections and corrections. Agent-review agreement is not human review or a proof by itself.

Task exit: reproducible frozen-source evidence, explicit supported and unsupported domains, and no concealed failing apparatus or unqualified source changes.
Task 6: Conclude the slice and deliver the research/runtime pair
Keep a concise packet under docs/foundational/phase3b4/: charter, frozen contract, model/argument, correspondence results, qualification, runbook and review/evidence index. Add matching runtime public documentation and compact qualification records. Preserve older sections and witnesses; keep the abstract and conclusion current without rewriting the discovery history.

Synchronize Markdown and LaTeX. Build the manuscript after substantive review and visually inspect changed pages and their continuations. If a correctness edit follows the build, produce a matching final PDF or explicitly disclose the mismatch; do not present a known stale document as current. No new bibliography entries unless the argument actually uses them.

Open cross-linked draft PRs, one per changed repository, with exact bases, executable commits, delivery commits, version decisions, independent reviews, failures, skips and the dependency order. Do not merge, release, force-push shared history or change unrelated branches. Do not create an empty companion PR if no change to that repository is justified.
Final disposition
Answer these questions directly:
1. Did explicit retained binding remove the immediate-selection restriction for the declared intervening programs?
2. What existing evidence sufficed, and what additional representation or versioning was necessary?
3. Did the new ordinary program actually complete reciprocal service without broader authority?
4. Which authority-changing or binding-changing steps still invalidate progress, and why?
5. Does the claim cover a source-level family of finite stretches, finite instances only, or a demonstrated obstruction?
6. Which model/runtime correspondence gaps remain, especially whole-Store crashes, separate-grant revocation and applied-but-unrecorded effects?
A positive result is request binding stable across the declared intervening computation and recovery, with authority still checked at the correct boundaries. It is not arbitrary program equivalence, general dataflow provenance, inevitable termination or N-agent closure.

A negative result is acceptable if it preserves the smallest genuine obstruction and identifies the missing premise or capability. Exhaust the frozen positive and negative families; do not stop at one convenient success or failure.
Convergence rule
Close this slice when the adjacency question is resolved at its declared scope. Do not remove unrelated restrictions merely to obtain a broader theorem. Recommend one next research question only after stating the remaining gaps. Independent owners, multihop delegation, general N-agent composition and a broader non-atomic-effect correspondence study remain separate assignments.

The intended outcome is that an author may receive a request, perform declared intervening work, survive interruption, and execute the exact authorized service without relying on accidental instruction adjacency. The runtime, not an application convention, must justify that binding.