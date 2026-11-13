# Phase 3B.6 design review

Verdict: **REJECT pending one narrow temporal-contract correction**.

Reviewed `docs/foundational/phase3b6/01-contract.md`, the two reconnaissance reports, owning research checker/model code, the frozen producer source at commit `5b4662fd682a220ca3065b15e86672b7e56567e7` in `/Users/jacob/dev/personal/alan`, and selected extracted normal-mode acknowledgment/cut artifacts. No implementation, tests, baseline reruns, runtime campaign, archived SQLite opens, or commits were performed. This is design readiness review, not executable qualification or a new archive-integrity audit.

## Required correction

### R1 / P2: Do not certify that the endpoint barrier remained held until owner death

Location: `docs/foundational/phase3b6/01-contract.md:50`, sentence “The endpoint response was held through owner death.”

The source establishes a narrower fact. In the frozen producer:

- `internal/durabletest/phase3b4_ancillary_test.go:180–244` consumes the endpoint post-commit acknowledgment and the runtime acknowledgment, kills/waits for the child, captures stopped state, and only then writes the endpoint release byte.
- `internal/testcmd/alan-resource/main.go:213–233` emits the committed journal entry, invokes the post-commit barrier using the HTTP request cancellation context, and emits the successful receipt only when that barrier returns without error. An error instead takes the HTTP error branch.
- `internal/boundary/boundary.go:110–139` permits the barrier wait to end on `ctx.Err()` before any release byte. It is not an unconditional wait-for-controller-release.
- `internal/durabletest/phase3b4_cuts_test.go` configures the first owner's Run with a five-second external timeout; `internal/executor/external.go` passes this to the HTTP client. The controller does not archive a synchronized duration proof that rules out timeout/cancellation before SIGKILL. Its runtime acknowledgment is evidence at emission, not proof that every transport remains in flight continuously until death.

Thus the pinned protocol proves **application before loss and no successful authentic receipt delivery before loss**, not that the endpoint remained paused or emitted no HTTP response of any kind until loss. A request timeout can cause an HTTP error path before kill while preserving all the relevant local/world observations. This ambiguity is harmless to the proposed result if explicitly abstracted away; it becomes a false temporal claim if a causal certificate records `barrier-still-held-at-kill` as established history.

Requested replacement or equivalent:

> The controller releases the endpoint only after acknowledged owner death. Before that release the endpoint cannot emit a successful authentic receipt; request cancellation may terminate the post-commit wait with an error. The dossier does not establish continuous barrier occupancy or transport liveness until death.

Apply the distinction to future causal certificates: certify application-before-loss and absence of authentic local result; hide timeout/error transport termination as an unobserved detail. Do not add a new runtime experiment or enlarge the model merely to distinguish this irrelevant error path.

## Scope and soundness assessment

Apart from R1, the design is appropriately small and implementable. It retains the previous atomic model's domain and outputs, restricts the new graph to one admitted original/service, and treats the shared owner loss as one transition. Splitting endpoint application from local retention is necessary; the old `Applied` progress predicate relies on atomic custody in `research/authorityservice/model.go` / `research/delegationcorrespondence/progress.go` and cannot be transported into the new model.

The N07 prelaunch inference is justified by actual source sequencing: the `authority.admission/allowed` observer executes synchronously inside `admitDispatch` (`internal/executor/dispatch.go:57–83`); its callback captures before/after revocation (`internal/durabletest/phase3b4_controlled_test.go:281–294`); the worker starts only after `admitDispatch` returns (`internal/executor/runner.go:219–240`). The passive boundary artifact has empty process fields, so this must remain producer-conditioned evidence, as the contract already requires. Final journals alone would not justify those earlier zeros.

Application-before-loss is justified independently of any relative ordering between the two acknowledgment emissions: resource application precedes its acknowledgment; both acknowledgments are received before the controller kills and waits. The contract correctly forbids inventing a cross-process total order. The selected normal artifacts name the matching endpoint post-commit source and sole runtime in-flight authorization; `kill.json` records the runtime PID and `signal: killed`. This review does not claim comprehensive sidecar validation, which remains implementation/qualification work.

The observation definitions avoid the main leak: endpoint application and unretained result values are excluded from the local view. Owner-stage/current-activation provenance is explicitly protocol-supplied; consequently the new theorem must describe equality under this declared local projection, not equality of all raw Store bytes or of all runtime telemetry. The distinction between a negative observation at a cut and permanent nonapplication is sound. The new model's equal-observation pair remains a finite-model result, not a matched v2 empirical pair.

## Nonblocking implementation checks

1. Keep `ResponseArrival` specifically an authentic successful result arrival. Unsuccessful HTTP completion is an invisible event in this abstraction; R1 must not accidentally turn such a completion into retained result or an extra model application.
2. Keep source-conditioned causal constraints separate from one chosen executable witness path. In particular, exact ordering among loss-independent/revocation-independent steps is not recovered chronology merely because a canonical model path chooses it.
3. Validate the full concrete raw/custody evidence while projecting only BA into this model. N07 has genuine AB endpoint activity; the one-contact restriction is per selected BA service, and AB remains a concrete control, not an accidental second modeled contact or a completion theorem.
4. Record archive membership/trust separately from semantic hostile-fixture validity. Coherently rehashed attack rejection should exercise the intended causal/custody check rather than count a mere manifest mismatch as a semantic mutant kill.
5. Preserve unknown/missing endpoint evidence as unknown/unsupported; a final complete dossier does not itself populate earlier cut observations. The current old checker permits optional endpoint artifacts and uses final physical exposures, so new cut attribution must remain explicit.

No broader design rewrite is requested. Correct R1 and obtain focused re-review before implementation.

## Focused R1 re-review

R1: **ADDRESSED**. Revised design verdict: **APPROVE**.

Read the revised contract at lines 21 and 50 and the matching Task 3 certification requirement in `docs/foundational/phase3b6/phase3b6-plan.md:65`. Response arrival now means authentic successful receipt arrival; unsuccessful HTTP completion is explicitly hidden and cannot establish custody. The historical certificate is limited to successful receipt withholding before the controller's post-death release, while continuous barrier occupancy and transport liveness until death are explicitly unestablished. This matches the pinned source behavior identified in R1 and leaves the intended application-before-loss/local-result-unknown argument intact.

This approval addresses design readiness after the narrow correction. It does not approve future implementation, assert test results, or broaden the proposed finite/source-conditioned claims. No tests, runtime campaigns, or subagents were used for this focused re-review; only this review report was appended.
