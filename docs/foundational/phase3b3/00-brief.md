# ALAN Research Phase 3B.3
## Request-Bound Delegation: Model-to-Runtime Correspondence

### Mission

Determine whether ALAN's actual runtime can realize the authority contract established in Phase 3B.2 while preserving the completion reachability established in Phase 3B.1.

Phase 3B.2 found that reciprocal service can preserve authority safety under specific conditions:

- Every application has retained authority rooted in legitimate permissions.
- Delegation attenuates complete permission tuples.
- Grants bind to the original request and logical generations.
- Fresh physical admissions recheck grant validity.
- Results and permission to communicate do not confer authority.
- Completion remains reachable when sufficient covering authority is available at the necessary decisions and admissions.

The finite model supports these claims within its declared bounds. The production runtime does not yet provide general request-bound S2 delegation.

**Central research question: Can an ordinary ALAN computation implement safe, progress-capable reciprocal delegation using a small retained authority primitive, without introducing privileged agent behavior or a parallel authorization subsystem?**

Do not assume the model-to-runtime correspondence holds.

Attempt to falsify it.

Complete the research through executable witnesses, independent validation, manuscript updates, frozen verification, and a draft PR.

Do not merge or release.

---

## 1. Establish the verified baseline

Research repository:

`github.com/JakeFAU/alan-research`

Relevant prior work:

- Phase 3B.1, draft PR #42.
- Phase 3B.2, draft PR #43.
- Phase 3B.2 base: `1d41d1e2`.
- Frozen executable qualification: `f75e2698`.
- Phase 3B.2 delivery: `d82a69c7`.

Runtime repository:

`github.com/JakeFAU/alan`

Verify current branches, commits, existing runtime delegation capabilities, and the relationship between the research fixtures and public runtime APIs.

Do not assume either research PR has merged.

Use an appropriate clean branch derived from the verified Phase 3B.2 state. Preserve prior witnesses and results.

Read the exact Phase 3B.2 authority contract, the sixteen finite graphs, the independent authority ledger, the Cartesian-permission mutant, the ambient-borrowing counterexample, and the runtime limitation recorded in the results.

Identify precisely which modeled authority transitions cannot currently be represented by the runtime.

Do not begin implementation until this gap is documented.

## 2. Freeze a minimal runtime authority contract

Define the smallest additional authority mechanism required for request-bound one-hop delegation.

The proposed mechanism must express:

1. Original authority held by a delegating principal.
2. The complete permission tuple being delegated.
3. The receiving principal.
4. The original service request.
5. The logical generations involved.
6. The permitted operation and resource scope.
7. The authorization decision and its retained witness.
8. The validity conditions for each fresh physical admission.
9. The conditions under which the delegated authority expires or becomes inapplicable.

Delegation must never manufacture a permission that is absent from the legitimate source authority.

Represent permission as a relation over complete tuples, rather than independently combining permitted values along each dimension.

For a delegated permission set \(D\) derived from source authority \(P\), require:

\[
D \subseteq P
\]

where containment means complete permission-tuple containment under the declared scope semantics.

The source authority itself must be valid at the relevant decision boundary.

A request, message, retained result, or invocation handle must not independently confer authority.

Freeze the contract, assumptions, and falsifiers before constructing the runtime implementation.

## 3. Separate the research model from the implementation

Maintain three distinguishable components:

**Specification:** The frozen operational authority contract.

**Model:** The finite transition system that determines which histories are permitted.

**Runtime:** The actual ALAN implementation that retains and executes computations.

Do not use the runtime implementation as its own specification.

Do not rewrite the reference model merely to match convenient runtime behavior.

When the runtime rejects a model-valid history, determine whether:

- The model omitted a relevant constraint.
- The runtime lacks an expressive primitive.
- The runtime incorrectly rejects legitimate authority.
- The witness does not actually correspond to the modeled execution.

When the runtime accepts a model-invalid history, treat that as a potential authority defect until independently resolved.

Preserve disagreements as evidence.

## 4. First attempt realization without changing the runtime

Before adding a primitive, try to express S2 using existing public ALAN facilities.

Use:

- Ordinary retained agents.
- Existing invocation identities.
- Explicit authority grants.
- Existing composition relationships.
- Original operation identities.
- Retained authorization and admission evidence.
- Existing recovery and inspection APIs.

Construct the smallest positive reciprocal-service scenario requiring a capability that the serving agent does not initially possess.

If an ordinary explicit grant can already satisfy the frozen contract, demonstrate it.

If it cannot, preserve the minimal obstruction.

Do not substitute a broadly authorized server for genuine delegated authority. That would make the example work while avoiding the research question.

Do not equate passing a capability identifier in a message with granting authority.

## 5. Specify the minimal missing primitive

If existing facilities are insufficient, propose the smallest reusable authority primitive.

Conceptually, it might permit:

\[
\operatorname{Delegate}
(
\text{source},
\text{recipient},
\text{request},
\text{permission},
\text{validity}
)
\]

This notation is illustrative, not a required API.

The primitive must preserve a verifiable chain from legitimate original authority to the delegated operation.

It must not introduce:

- Ambient authority inheritance.
- Automatically transitive delegation.
- General multihop delegation.
- Unbounded capability propagation.
- Implicit delegation through messages.
- A second resource-accounting system.
- A separate scheduler.
- A privileged agent category.

A grant is not a physical execution.

A grant is not a reservation.

A grant is not an irrevocable promise that a future operation may execute.

Under this slice's selected temporal contract, fresh physical admissions recheck the relevant authority.

Preserve distinct decision boundaries for source authorization, delegation, reservation, admission, application, and result consumption.

## 6. Implement only what the experiment requires

If a runtime change is necessary, keep it narrowly scoped and independently reviewable.

Prefer existing authority and operation protocols.

Use stable original identities and retained evidence.

Avoid introducing a new persistence format unless unavoidable.

If a schema change is unavoidable, document exactly why the existing store cannot represent the required facts and qualify backward compatibility separately.

Do not modify runtime behavior merely to make the experimental fixture easier.

If a safe realization requires a substantial new authority subsystem, stop the implementation expansion, preserve the obstruction, and complete the slice as a research result with a proposed follow-up contract.

Do not turn a bounded experiment into an open-ended runtime rewrite.

## 7. Reproduce the decisive Phase 3B.2 cases

Use exactly two logical agents.

Maintain the established A-servicing/B-serial scheduling contract.

Exercise three service policies:

### S1: Server-owned authority

A serves B using only authority A legitimately possesses.

No delegation is required.

### S2: Explicit attenuated delegation

B delegates an operation-bound subset of its legitimate authority to A.

A completes the service using the retained delegated grant.

### S3: Ambient borrowing

A attempts to exercise authority inferred from B's request or identity without legitimate delegation.

The runtime must reject the unauthorized operation.

Do not implement S3 as a supported runtime feature merely to demonstrate its danger. Construct it as an adversarial attempted action.

Preserve the original Phase 3B.2 positive and negative witnesses.

## 8. Test the Cartesian authority mutant

This is a mandatory negative test.

Construct legitimate permissions such as:

\[
P=\{(r_1,o_1),(r_2,o_2)\}
\]

Attempt an operation requiring:

\[
(r_1,o_2)
\]

The latter permission must remain unauthorized unless it was separately present in the legitimate authority relation.

The test must fail if an implementation independently combines resource scopes and operation permissions into a Cartesian product.

Preserve the full tuple at every relevant authority boundary.

Check the original grant, delegated grant, operation authorization, physical admission, and retained result.

No final-state capability-set comparison alone is sufficient.

## 9. Bind delegation to causal identity

Test that delegation cannot be transferred between:

- Different original service requests.
- Different caller generations.
- Different serving generations.
- Different root authority scopes.
- Different resources.
- Different operations with superficially equal payloads.
- Original and retried invocations.
- Successive owner activations.

A grant must remain linked to the original request that justified it.

An identical payload is not evidence of identical authority.

An old grant must not become usable merely because the new request has matching parameters.

Do not let reopening generate a new grant implicitly.

## 10. Temporal authority and revocation

Retain Phase 3B.2's chosen contract: fresh physical admissions recheck validity.

Construct controlled histories in which:

1. Delegation is valid at grant time.
2. The serving agent commits a service operation.
3. Source authority is revoked before physical admission.
4. The operation attempts to proceed.

The runtime must enforce the frozen admission rule.

Separately test revocation after a legitimate admission.

Do not claim that revocation reverses a physical action already admitted or applied.

Preserve unresolved external effects.

Do not introduce general cancellation propagation or distributed revocation.

The objective is only to establish whether retained delegated authority remains valid at the relevant boundary.

## 11. Two-agent progress qualification

Reproduce the Phase 3B.1 reciprocal-service scenario with the S2 authority contract.

Under sufficient covering authority, test whether completion remains reachable from every reachable modeled state.

Then remove or revoke a required grant.

The resulting inability to complete may be legitimate.

Distinguish:

- Authority-safe completion.
- Authority-safe waiting.
- Authority-safe refusal.
- Deadlock caused by scheduling.
- Blocked completion caused by insufficient permission.

Do not report an authority refusal as a scheduling defect.

Do not manufacture permissions to make a progress test pass.

Preserve the distinction:

\[
\forall s\in R,\ \exists s\leadsto Complete
\]

versus eventual completion under all fair executions.

Only claim the property actually checked.

## 12. Actual runtime execution and SIGKILL

Demonstrate real same-Store reciprocal service.

Use ordinary ALAN computation paths, not a simulated scheduler standing in for the runtime.

Exercise acknowledged process termination at:

- Before delegation commitment.
- After delegation retention.
- Before service authorization.
- After operation authorization.
- Before physical admission.
- After physical admission with unknown outcome.
- After retained result.
- Before result selection.
- During mutual waiting.
- After one agent completes while the other remains pending.

Reopen the original Store with a fresh owner.

Verify:

- Stable logical identities.
- Preserved grant provenance.
- No authority amplification.
- No unauthorized additional physical admission.
- No refunded expenditure.
- No invented result.
- No disappearance of unresolved obligations.
- No new grant created solely by recovery.
- Original results consumed without reissue where legally possible.

Use acknowledged boundaries rather than sleeps.

Preserve raw retained histories and independent external-contact evidence where physical operations are involved.

Do not claim independent-owner recovery. This remains a same-Store, single-host experiment.

## 13. Differential model-to-runtime testing

Generate bounded histories from the frozen operational model.

Classify each as:

- Valid and complete.
- Valid but incomplete.
- Valid and blocked by authority.
- Invalid due to missing authority.
- Invalid due to scope amplification.
- Invalid due to generation or request mismatch.
- Invalid due to revocation timing.
- Invalid due to impossible causal order.

Replay corresponding histories through the ordinary runtime where expressible.

Use disposable tampered stores for adversarial invalid histories that cannot be produced by legitimate execution.

Compare the model's verdict against:

- Actual execution acceptance or refusal.
- Retained authorization evidence.
- Replay.
- Owner reopen.
- Typed inspection, where available.

Independent checkers must not rely solely on the runtime's authorization helper.

Record every disagreement.

A counterexample that falsifies the model-to-runtime correspondence is a valid research outcome.

Do not hide it behind an amended test expectation.

## 14. Independent validation

Use at least two distinct checks:

1. A forward operational model that constructs authorized histories.
2. An independent retained-history authority checker that derives whether applications have legitimate authority chains.

Additionally, compare the runtime's observable behavior with those verdicts.

Challenge the independent checker using deliberate mutations:

- Missing grant.
- Forged grant origin.
- Altered recipient.
- Wrong original request.
- Wrong logical generation.
- Expanded permission tuple.
- Revoked permission at admission.
- Valid grant reused by a different service.
- Retained result substituted for authorization.
- Equal-payload operation substitution.

Preserve accepted and rejected counterexamples with exact identities and digests.

## 15. Preserve the previous research results

Do not weaken or rewrite:

- Phase 3B.1's serial/serial deadlock result.
- Phase 3B.1's reciprocal-service completion reachability.
- Phase 3B.2's sixteen-graph authority findings.
- The ambient-borrowing counterexample.
- The Cartesian-permission mutant.
- The generation-binding defect discovered during review.
- Existing crash/reopen custody evidence.

If the runtime experiment requires a narrower set of assumptions, state the narrowing explicitly.

Do not silently generalize prior finite-model results.

## 16. Manuscript and research artifacts

Update the dissertation with:

- The exact S2 runtime realization contract.
- The minimal missing primitive, if one was necessary.
- The correspondence claim being tested.
- The supported and unsupported operational histories.
- The positive reciprocal-service witness.
- Authority-laundering negative controls.
- The SIGKILL/reopen evidence.
- Any model/runtime disagreement.
- The remaining limits of the result.

Separate theoretical findings from engineering implementation results.

An executable example is not a theorem.

A finite exhaustive model is not a proof of arbitrary composition.

Rebuild and inspect the PDF.

Keep all claims limited to two agents, the declared service policy, one-hop delegation, and the tested same-Store setting.

## 17. Verification and review

Run the relevant Phase 3A, 3B.1 and 3B.2 regression suites, plus the new 3B.3 tests.

Include:

- Ordinary tests.
- Focused integration tests.
- Race-enabled recovery tests.
- Exhaustive bounded model checks.
- Independent authority-ledger checks.
- Replay and inspection checks.
- Formatting and Staticcheck.
- Manuscript consistency.
- PDF build and inspection.

Use a frozen source revision for final verification.

Record exact commands, toolchain versions, source identities, test counts, skipped tests and failures.

Preserve pre-existing Staticcheck findings separately. Do not call a non-green whole-repository run green.

Obtain independent reviews covering:

1. Authority semantics and primitive minimality.
2. Model-to-runtime correspondence.
3. Capability attenuation and request binding.
4. Recovery, revocation and physical-admission boundaries.
5. Progress and retained obligations.
6. Manuscript accuracy and reproducibility.

Reviewers must inspect source and retained witnesses.

Preserve substantive rejected candidates and their corrections.

## 18. Delivery and scope control

Deliver a draft research PR against the verified appropriate base.

Do not merge or release.

If a runtime implementation change is required, keep its source and qualification in the runtime repository and its theoretical analysis in the research repository.

Do not mix unrelated changes into either repository.

Report separately:

- Research PR and commit.
- Runtime PR and commit, if applicable.
- Frozen authority contract.
- Missing primitive or existing mechanism used.
- Model results.
- Runtime realization results.
- Positive and negative witnesses.
- SIGKILL/reopen findings.
- Correspondence failures and resolutions.
- Independent review findings.
- Remaining assumptions.

Do not broaden into multihop delegation, arbitrary N-agent composition, distributed owners, non-atomic external effects, or generalized cancellation.

Those require separate studies.

## 19. Definition of done

The slice succeeds if it establishes one of these outcomes:

**Positive realization:** The actual ALAN runtime implements the frozen S2 contract and reproduces the bounded two-agent reciprocal-service behavior without authority amplification.

**Negative realization:** A concrete retained execution demonstrates that the existing or proposed runtime mechanism violates the model's authority contract.

**Implementation obstruction:** The required authority contract cannot be expressed without a larger primitive than the scoped experiment permits, and that obstruction is demonstrated precisely.

For a positive result, require:

- Authentic retained grant provenance.
- Complete permission-tuple attenuation.
- Original request and generation binding.
- Fresh-admission validity checks.
- Correct rejection of ambient borrowing and Cartesian amplification.
- Progress under the declared sufficient-authority assumptions.
- Recovery without invented grants or duplicate physical operations.
- Agreement between independently checked model and runtime histories.

## 20. Final research questions

Answer directly:

**Can ALAN's real two-agent execution model implement request-bound, one-hop attenuated delegation using its ordinary retained-computation machinery?**

**Does the implementation preserve both authority safety and completion reachability under the exact Phase 3B.2 conditions?**

**What is the smallest additional runtime primitive required, if any?**

**Which assumptions would have to change before this result could extend to independent owners, multihop delegation, or three agents?**

Do not answer the last question by claiming those extensions already work.

The purpose of Phase 3B.3 is to determine whether the mathematical authority contract survives contact with a real durable runtime.