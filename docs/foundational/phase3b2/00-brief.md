# ALAN Research Phase 3B.2: Authority Attenuation Under Reciprocal Service

## Mission

Investigate whether two independently authorized, durably retained agents can serve one another while waiting without amplifying, laundering, or ambiguously transferring authority.

Phase 3B.1 established that:

- Two agents can retain authentic result custody while making no progress.
- Under the declared serial/serial contract, 144 reachable states have no completion path.
- Allowing either agent to serve while waiting restores a completion path from every reachable state.
- The runtime reproduces the mutual wait across reopening, including quiescence with two unpaid obligations.

Phase 3B.2 asks whether this liveness improvement preserves authority safety.

**Research question: Can reciprocal service remain both progress-capable and authority-safe when the participating agents possess different capabilities?**

The objective is a bounded result: a counterexample, a conditional theorem, or a precisely characterized unresolved boundary.

Do not presume authority amplification exists. Attempt to discover it.

Do not generalize to N agents.

Complete the research, executable witnesses, independent review, manuscript update, and draft PR.

Do not merge or release.

## 1. Verify the research baseline

Repository: `github.com/JakeFAU/alan-research`

Relevant previous work:

- Phase 3B.1 draft PR #42.
- Two-agent reciprocal-service model.
- Serial/serial deadlock witness.
- Reciprocal-service reachability checks.
- Same-Store runtime fixture.
- Existing retained custody, authority and recovery contracts.

Fetch and verify current `main`, PR #42, its tested commit, and outstanding work.

Do not assume PR #42 is merged.

If it remains unmerged, derive the new slice from the verified Phase 3B.1 source without modifying or prematurely merging that PR.

Use a clean branch or worktree.

Read the Phase 3B.1 frozen assumptions, model, witness manifest and results before proposing new semantics.

Preserve its negative and positive controls unchanged.

## 2. Freeze the authority model before experimenting

Define a finite universe of capabilities.

Begin with two agents:

\[
A,\quad B
\]

Give them asymmetric authority.

For example:

\[
C_A=\{a,r\}
\]

\[
C_B=\{b,r\}
\]

where:

- \(a\) is a capability available initially only to A.
- \(b\) is a capability available initially only to B.
- \(r\) is a shared, narrowly scoped capability.

These are illustrative symbols, not required runtime names.

Model capabilities as permissions over concrete operations, resources and scopes, not merely unstructured labels.

Distinguish:

1. Authority originally held by the serving agent.
2. Authority held by the requesting agent.
3. Authority explicitly delegated for one service.
4. Authority retained from an earlier authorized operation.
5. Authority revoked before an attempted operation.
6. Authority that was never validly granted.

Identify the exact decision boundary at which authority is checked.

A retained request does not automatically establish authority for every future action taken to satisfy it.

A service request must not confer the requester's entire capability set by implication.

## 3. Define three candidate service contracts

Investigate three deliberately distinct rules.

### Contract S1: Server-owned authority

When A serves B, A may exercise only authority A already possesses.

B's request supplies data and an objective, but no additional capability.

This is the minimal baseline.

### Contract S2: Explicit attenuated delegation

B may delegate a subset of its authority to A for a particular service.

The delegated grant must be:

- Explicit.
- Scope-bounded.
- Bound to a specific request or operation.
- Traceable to its original authority.
- Non-amplifying.
- Durable across recovery.
- Subject to the declared revocation semantics.

A's effective service authority must not exceed the permitted combination of its own authority and valid delegated authority.

### Contract S3: Ambient or implicit delegation

The serving agent may use authority inferred from the requester's identity, retained context, or reachable capabilities without an explicit operation-bound delegation.

Treat S3 as a candidate for falsification, not a recommended design.

Investigate whether S3 admits histories that S1 or S2 reject.

Do not assume every implicit delegation is invalid unless the frozen policy says so.

The purpose is to identify the exact authority contract needed for safe composition.

## 4. Define safety and progress separately

### Authority safety

Every externally meaningful operation must have an independently checkable authorization chain.

For each operation \(e\), require a valid witness establishing:

- The principal performing the operation.
- The originating request.
- The capability exercised.
- The resource and scope.
- The source of authority.
- The delegation chain, if any.
- The relevant authority epoch or revision.
- The decision boundary.
- The original retained authorization decision.

No valid execution may create authority solely through cyclic delegation, repeated service, or recovery.

### Progress

Retain Phase 3B.1's completion-reachability criterion:

\[
\forall s\in R,\quad
\exists\text{ a valid path }s\leadsto\text{Complete}
\]

Evaluate it separately for each service contract.

Do not confuse existence of a completion path with inevitable completion under all schedulers.

If preserving authority safety removes all completion paths from some reachable states, record that result.

Do not weaken the authority contract merely to recover liveness.

The research goal is to characterize the feasible intersection of safety and progress.

## 5. Start with the smallest counterexample

Begin with exactly two agents, one reciprocal dependency, and a minimal capability universe.

Construct paired cases:

**Positive control:** A serves B using authority A legitimately possesses.

**Delegation control:** B explicitly delegates a narrowly scoped capability to A, which exercises only that grant.

**Negative control:** A performs an operation requiring B's capability without a legitimate delegation.

**Cycle control:** A and B exchange authority claims that appear locally plausible but lack a valid original granting path.

**Recovery control:** An authorization decision is committed, the owner dies, and reopening preserves the same authority without inventing a fresh grant.

Do not introduce more agents, more capabilities, or more protocol states until the minimal cases are understood.

Preserve the smallest decisive histories.

## 6. Investigate authority laundering

Try to construct executions in which an unauthorized action appears authorized because its origin has been obscured by composition.

Examples to investigate:

- A requests work from B, which returns a capability-bearing result that A treats as authority.
- B requests an operation from A that A could not authorize for the relevant resource scope.
- A and B delegate the same capability back and forth until its original scope is lost.
- A service operation is rebound to a different request after reopening.
- A previously valid grant is reused for a different operation.
- An authorization from one generation is incorrectly applied to another.
- A refusal is followed by a later service that appears to satisfy the refused operation.
- Two individually valid grants combine into an operation outside either grant's permitted scope.

The last case is particularly important.

Investigate whether authority composition is defined by a union, an intersection, an explicit delegation rule, or some other constrained operation.

Do not choose the rule based on convenience.

The semantics must distinguish permission to request an operation from permission to execute it.

## 7. Model cancellation and revocation only as needed

This slice is not a general cancellation study.

However, authority attenuation cannot be evaluated meaningfully without a minimal temporal rule for grant validity.

Use the narrowest existing revocation semantics necessary to establish whether a grant is valid at its authorization boundary.

Preserve the distinction between:

- Revoking authority for a future dispatch.
- An operation already authorized and committed.
- A physical effect already admitted.
- An external outcome that remains unknown.

Do not claim that revocation undoes previously executed external effects.

Do not add general cancellation propagation, distributed revocation or cross-owner authority management.

Those remain separate research questions.

## 8. Exhaustive bounded model

Extend the Phase 3B.1 finite transition model rather than replacing it.

Introduce:

- Initial capability ownership.
- Scoped delegation.
- Service request identity.
- Retained authorization evidence.
- Validity and revocation boundaries.
- Reciprocal service while waiting.
- Process interruption and recovery where represented by the existing model.

Enumerate reachable states for S1, S2 and S3 under explicitly stated bounds.

For each contract, report:

- Number of reachable states.
- Authority-safe reachable states.
- Unsafe histories, if any.
- States from which completion remains reachable.
- States from which completion is impossible.
- Minimal counterexample length.
- Additional assumptions required for the result.

Count models and histories precisely. Do not compare raw state counts across contracts without accounting for differences in their state spaces.

Retain Phase 3B.1's original 144-state serial/serial baseline as an independent regression, not as a number the new model must reproduce after extending its state vocabulary.

## 9. Independent validation

Use at least two genuinely independent methods.

One should enumerate the operational state transitions directly.

The other should validate retained authority chains and service relationships from the resulting histories without calling the production transition validator.

The independent checker must establish:

- Grant origin.
- Delegation attenuation.
- Scope containment.
- Operation binding.
- Authorization order.
- Retained identity.
- Absence of unsupported authority amplification.

Include adversarial mutations that preserve superficial counters, identities or terminal outcomes while changing a decisive authorization edge.

A checker that only compares final capability sets is insufficient.

Retain both accepted and rejected histories.

If two checkers agree because they share the same incorrect semantic assumption, record the limitation rather than calling their agreement independent proof.

## 10. Same-Store runtime witness

Extend the existing two-agent same-Store runtime fixture only after the model yields a decisive candidate.

Demonstrate at least:

1. Reciprocal service under legitimate local authority.
2. Reciprocal service using a valid attenuated delegation, if existing runtime primitives support it.
3. Rejection of an unauthorized service attempt.
4. A real process death after authorization but before service completion.
5. Successful reopening without authority amplification.
6. Retained unfinished obligations when progress is legitimately blocked.

Do not synthesize a successful witness by editing the database.

Use ordinary runtime operations for positive cases and disposable corrupted copies for adversarial cases.

If the current runtime lacks a delegation primitive required to realize the model, do not quietly implement a broad new authority system.

Instead:

- Identify the exact missing primitive.
- Demonstrate the model result independently.
- Record the runtime realization boundary.
- Recommend a narrowly scoped follow-up implementation if justified.

A research slice may succeed with a counterexample and a proven implementation obstruction.

## 11. Preserve the custody result

Authority experiments must not weaken Phase 3B.1's retained custody guarantees.

For every service operation, identify:

- Original request.
- Serving agent.
- Original authorized operation.
- Retained result.
- Publication.
- Delivery.
- Actual consumption.
- Remaining obligation, if any.

A completed service must not be inferred merely from authentic result custody.

A refused service must not be treated as completed work.

A system may remain safely blocked when no authorized completion path exists.

**Safety takes precedence over manufactured progress.**

## 12. Crash and recovery qualification

Exercise real SIGKILL and fresh-owner reopening at meaningful boundaries:

- Before delegation or authorization.
- After a valid grant is retained.
- After an unauthorized request is refused.
- After service authorization.
- After result retention.
- Before result consumption.
- During mutual waiting.
- After one agent completes while the other remains pending.

Verify that reopening:

- Does not expand authority.
- Does not reissue completed external operations.
- Does not erase refusals.
- Does not manufacture result custody.
- Does not discharge unpaid obligations.
- Preserves original causal identities.

Do not conflate same-Store recovery with independent-owner or distributed recovery.

The latter remains outside this milestone.

## 13. Theoretical deliverable

The paper must state one or more precise propositions, supported by the investigation.

Possible outcomes include:

**Outcome A: Counterexample**

Locally legitimate operations compose into an authority violation under a specified service rule.

Provide the minimal witness and identify the missing global invariant.

**Outcome B: Conditional safety theorem**

Under explicitly bounded delegation and authorization assumptions, reciprocal service does not amplify authority.

State the assumptions and proof scope precisely.

**Outcome C: Safety-progress conflict**

A restrictive authority policy prevents unauthorized service but also eliminates completion paths available under a more permissive rule.

Characterize the trade-off.

**Outcome D: Insufficient evidence**

The model cannot distinguish valid from invalid authority because the retained protocol lacks a necessary witness.

Identify exactly what is missing.

Do not force a positive theorem.

A well-supported negative result is equally valuable.

## 14. Manuscript and evidence

Update the Phase 3 research manuscript and foundational results.

Separate:

- Frozen assumptions.
- Formal model.
- Exhaustive results.
- Independent checks.
- Runtime witnesses.
- Counterexamples.
- Unresolved boundaries.
- Interpretation.

Explicitly distinguish safety, completion reachability, and fair-scheduler termination.

Maintain the existing dissertation structure.

Rebuild the PDF using the established toolchain and inspect the resulting document.

Do not claim a universal theorem from finite enumeration alone.

Do not introduce speculative general-N results.

## 15. Verification

Run the relevant Phase 3A, Phase 3B.1 and new Phase 3B.2 tests.

Include:

- Ordinary repository tests.
- Relevant race-enabled tests.
- Integration and subprocess recovery tests.
- Exhaustive bounded model checks.
- Independent authority-chain validation.
- Staticcheck and formatting.
- Manuscript consistency and PDF build.
- Source and witness identity verification.

Use an immutable source revision for final qualification.

Record toolchain versions, exact commands, exit statuses, skipped tests and unresolved issues.

Preserve source hashes and decisive witness digests.

Do not alter previous results to make the new theory appear consistent.

## 16. Independent review

Obtain separate reviews for:

1. Authority semantics and attenuation.
2. Model completeness and counterexample minimality.
3. Safety versus completion reachability.
4. Retained evidence and recovery.
5. Runtime correspondence and scope.
6. Manuscript claims and reproducibility.

A reviewer should attempt to falsify the central result using a smaller history or a missing authorization edge.

If a counterexample fails under closer scrutiny, retain that finding and revise the claim.

Do not treat reviewer agreement as independent experimental evidence.

## 17. Delivery

Open a draft PR against the verified appropriate base.

Do not merge or release.

Report:

- Base and delivery commits.
- Frozen model and assumptions.
- Exact candidate service contracts.
- Reachable-state counts.
- Authority-safe and unsafe histories.
- Completion-reachability results.
- Minimal positive and negative witnesses.
- Independent validation.
- Same-Store runtime results.
- Crash/reopen evidence.
- Manuscript and PDF status.
- Independent review dispositions.
- Unresolved boundaries.

Preserve decisive evidence in the repository.

Avoid redundant large logs that cannot be independently interpreted.

## 18. Definition of done

This milestone succeeds when it determines, within a precisely declared two-agent model, **whether reciprocal service can preserve authority attenuation while restoring completion reachability**.

Every accepted operation must have a valid retained authority chain.

Every rejected operation must have an identifiable violated condition.

The model must distinguish authentic custody from completed service and legitimate refusal from loss of progress.

No claim about N-agent composition, independent owners, arbitrary distributed systems, or unmodeled external effects.

No new production execution engine.

No broad delegation framework without a separately demonstrated need.

## 19. Final research question

Answer directly:

**What authority contract allows two durably retained agents to serve each other without creating permissions that neither agent legitimately possessed?**

Then answer the harder companion question:

**Does that contract preserve the completion reachability recovered in Phase 3B.1, or does it expose a fundamental safety-progress trade-off?**

If the result is conditional, name the conditions.

If the result is negative, preserve the smallest counterexample.

If the result is unresolved, identify the missing evidence or semantic assumption.

The goal is not to make the architecture look correct. It is to discover exactly when two-agent composition is correct.