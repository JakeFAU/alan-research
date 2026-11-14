# Phase 3C.2 builder sanity check

**Verdict: APPROVE for implementation of the proposed research packet.**

This is a read-only preflight assessment of implementability and the proposed arguments, not final review, executed qualification, or approval of an implementation that does not yet exist. I found no blocking inconsistency in the contract or plan. The negative pair and the conditional stable-ready result follow from the specified finite-chain relation. Their intended concrete witnesses remain to be constructed and checked.

## Scope and evidence inspected

Read `AGENTS.md`, `docs/foundational/phase3c2/00-contract.md`, and `docs/foundational/phase3c2/phase3c2-plan.md`. Inspected the unchanged `research/contractcomposition/{model,concrete,compose}.go`, relevant concrete tests, and `research/twoagent/model.go`. No tests, graph enumeration, runtime execution, SQLite access, `.envrc` access, implementation, index changes, branch changes, or subagents were used. Only this temporary report was written.

Research question: does retaining only occupancy lose the distinction between a blocking upstream reservation and a completable downstream reservation, and can a monotone readiness predicate restore existential finite completion in the fixed chain? The semantic boundary is original-bound reusable reservation added to the existing stage/custody relation; atomic application/result retention/release, immutable policies/topology, and recoverable physical owners are essential premises.

## Negative pair is constructible

For B serial and shared limit one, start from genuine empty initial custody and execute:

1. AB commit, AB send-request, AB admit.
2. BC commit, BC send-request, BC admit.
3. Branch into AB reserve-capacity versus BC reserve-capacity.

Both branch endpoints have stage vector `(W,W)`, identical full base transport and availability, and occupancy one. Reservation changes only the holder bits. The common base has no result bytes, so opaque-result choice creates no hidden difference in this comparison.

In the AB-holder branch, AB cannot apply while BC is W; BC cannot reserve the occupied root. Neither original can consume because neither is applied. Retransmission, deduplication, interruption and recovery cannot change either stage or release the hold. This establishes absence of a finite completion suffix in the specified relation, even though microsteps need not stop.

In the BC-holder branch, C has no outgoing obligation, hence its serial service guard is true as well as its servicing guard. Apply BC, route/record/consume its actual result, then reserve/apply/consume AB. BC consumption needs no slot. This gives a finite completion suffix.

The underlying guard is explicit in `research/contractcomposition/model.go:134-152`: waiting means the outgoing original is W or P. Concrete admission is not blocked by that guard; only application is filtered by it in `research/contractcomposition/concrete.go:163-185`. Thus the prefixes do not accidentally require B to service AB before BC admission.

## Ready-now attack and old-trace restriction follow

Ready-now can reserve AB after AB admission while BC is U, because B is not yet waiting. Committing and admitting BC then produces the noncompletable AB-holder endpoint above. A matched BC-holder endpoint is possible: commit/admit both before reserving BC. The final full base transport and physical states can again be identical without editing any captured state. Readiness invalidates at BC commit, not merely at admission; the witness should expose that exact step.

An old/early trace can commit/admit/reserve/apply AB while BC remains U. Stable-ready with B serial cannot reserve AB in that state; it requires BC D. Under the old progress observation, reaching BC D contributes visible BC commit/apply/consume events. Those events cannot be hidden as transport preparation. The new policy therefore cannot reproduce that old progress trace by erasing reserve events. The proposed exclusion of old-abstraction backward lifting is necessary and correctly stated.

## Positive argument is coherent under the stated premises

The useful invariant is stronger and more precise than general eventual service: every held original has its unchanged service guard true, and that truth remains true until its application. For stable-ready, acquisition establishes either immutable servicing policy or all the callee's outgoing originals D. D is absorbing; neither transport nor recovery changes stages or policies. Another acquisition does not affect either premise. Application removes its own hold atomically. Commit cannot invalidate readiness of a serial held AB because BC was already D; C has no outgoing original. This is a direct induction over the declared transitions.

Every existing holder can consequently apply after recovery of its callee, releasing its slot and retaining its genuine result. Finish those holders, then complete any unfinished BC before any unfinished AB, consuming already-applied results without a reservation. BC can acquire whenever it is W and not held once existing holders have applied. After BC D, a serial B satisfies stable-ready for AB. This construction covers shared-one, shared-two and separate-one without assuming fairness or inevitable scheduling. If several holders exist, finishing each is finite; their readiness is not lost when another holder applies or consumes.

The transport lifting argument can be made locally for each new abstract action:

- Commit: recover its caller if needed and commit.
- Reserve: recover the needed endpoints, send a request if absent, admit if needed, then reserve. These preparations do not change stages or holds, so the policy predicate and available capacity survive.
- Apply: a reachable hold implies admitted W work. Recover the callee and apply, transferring the actual chosen result into durable custody while clearing that original's hold.
- Consume: recover the caller; use its recorded result, or its volatile result, or actual channel/retained service custody. Durable-ACK custody guarantees an available retained source if none of the caller-side sources exists. Route, record and consume without reserving.

Every preparation above is an abstract stutter. Each completed action has the exact named resource-aware successor, which supports inductively lifting a finite abstract completion path. Conversely, projecting a concrete finite completion path and removing stutters yields an abstract completion path. These are source-level arguments conditional on reachable representatives; malformed but superficially valid states are not covered. The contract correctly reserves independent finite checks as supporting evidence rather than identifying them with this argument.

## Implementation risks to resolve within the existing plan

These are nonblocking implementation obligations, not requests to expand scope.

1. **Do not reuse old lifting unchanged.** `contractcomposition.Lift` checks the old action set and can prepare admission as part of old apply (`concrete.go:214-269`). `NewAction` rejects reserve-capacity (`model.go:213-218`). The wrapper must own enablement/lifting for the new relation and ensure final apply has a pre-existing hold. Reusing unchanged transport successors is feasible: preserve non-apply steps, filter apply on a hold, clear only that hold on apply, and append reserve transitions. Candidate states from the old successor relation should remain the source of authentic result transitions.

2. **Validate hold/admission causality, not just W and occupancy.** A concrete state with a W holder but no admitted request is structurally possible to construct in Go and violates the lifting premise. State validation should reject it; real witnesses must start from the true initial state. The invariant is easy to preserve because existing admission is monotone. Source proofs must explicitly mention this invariant.

3. **Do not inherit factored probes as reachable histories.** Existing `TestEveryOldEdgeRepresentativeProjectionAndLifting` constructs canonical other-edge records directly. Those are probes, not witnessed whole-chain histories. The new test domain needs exact stage/hold/admission consistency, the stable-ready reachable-state constraint, and physical-owner consistency when B acts as both AB callee and BC caller. Synthesized probe states should be labeled as such; no synthesized receipt should appear in exported operational histories. The plan already requires genuine witness prefixes separately.

4. **Keep root identity separate from the old effect resource.** Existing component validation rejects duplicated `Binding.Resource`, which represents the distinct counter effect resources. Shared Q must be an additional descriptor/membership relation. Reusing `Binding.Resource` to encode Q would break the required unchanged bindings and would cause valid shared-Q joins to fail. Validate membership coverage against the component's declared originals and canonicalize root descriptors by identity; joining references must not add limits.

5. **Freeze precise corruption semantics.** The old concrete Event has an actor but no event activation field; activation is in per-owner availability. A wrong-activation hostile control can mutate a post-state and be rejected by declarative replay. If the new schema adds an event activation, its exact binding/validation rules must be specified. Neither interpretation establishes a broader stale-message fencing theorem, which this packet does not claim.

6. **Oracle independence includes lifting and validation.** Preserving WorldResult in full captured states is appropriate for audit of provenance, but operational checks must not read it. The old tests erase this field to compare successor/lift behavior. Carry that independent-observation test to the wrapper, including reserve and capacity-aware apply, and ensure the Python checker distinguishes actual prior apply from a claimed final receipt.

7. **Complete finite graphs need complete arc comparison.** A reverse reachability solver over a truncated exported tree can label some states plausibly while missing legal transitions. The plan's independent regeneration and removed-non-tree-arc hostile input directly address this. Keep exact identity/duplicate checks for nodes/arcs as well as counts; counts alone do not discharge completeness.

## Remaining obligations and bounds

No graph totals, probe counts, actual witness records, test results, checker results, manuscript alignment, PDF verification, or inherited-suite source identity were established by this review. All remain implementation/qualification work in the three tasks. The unchanged-source reuse provision is conditional and should not be reported as a new full-suite run.

The contract's stable-ready repair is sufficient for this fixed chain and changes the observable old trace set. It does not establish necessity, optimal permissiveness, arbitrary-graph completion, distributed allocation, fairness, or runtime correctness. No internal requirement forces such a stronger claim. Subject to those explicit limits, the planned packet is implementable and its central source-level conclusions are defensible.
