# Phase 3B closure-assessment contract

## Question and authorized scope

Has the accumulated Phase 3B work answered its original research obligations well enough to close the bounded investigation, and which obligations remain for Phase 3C or a later extension? The user approved a consolidated goal-to-evidence assessment after merging Phase 3B.7. The verdict is an audit conclusion to establish, not a predetermined claim that the phase passed.

The original authority for the questions and progression gates is `../phase3/README.md`, read at base `7ae2c5ddd5e18d2907dd91fcefc896456a25562d`. The original discussion distinguishes local correctness, system safety, progress and compositional closure. Preserve those distinctions. The repository calls the completed matched-capture packet Phase 3B.7; do not invent a Phase 4 charter.

## Closure rule

A research obligation may be answered by a proved statement under explicit assumptions, a bounded model result, an adverse executable witness, or a precisely characterized obstruction, according to the evidence actually present. A passing implementation example does not prove a general theorem. A known limitation is not automatically an unfinished experiment; conversely, an original load-bearing obligation cannot be erased by relabeling it as out of scope.

Assess separately:

1. Whether the original Phase 3B exit evidence exists: explicit rely/guarantee contracts, adverse crash/revocation schedules and correspondence to an implemented fragment.
2. Which progression gates have been attacked: stable identity through recovery and new invocation; request/result/consumption binding; attenuation and shared-authority conservation; cancellation versus already admitted work; the server effect/receipt crash window; reciprocal outstanding obligations; in-flight/unpaid obligations and quiescence.
3. What the evidence does not establish: a unified model/runtime refinement, universal recovery/liveness, distributed ownership, or compositional closure. The original Phase 3C operator/equivalence, regrouping/hiding, three-agent cycle and shared-authority fan-out obligations remain distinct.

Each matrix row must name the original obligation, a specific source/artifact, its evidence tier and assumptions, the adverse result or counterexample where applicable, and a closure disposition. Choose a clear overall verdict: closed within an explicitly stated domain, not closed with named necessary gaps, or unresolved with the smallest missing evidence. Do not use unsupported generalizations as an unbounded backlog.

## Evidence standards

- Inspect owning source and machine-readable records where they can change a claim; prior prose and reviewer approval alone are not proof.
- Attribute historical experiments and tests to their pinned packets. New read-only inspection and hash verification are audit evidence, not a new runtime campaign or test-suite pass.
- Preserve the non-atomic effect/result boundary. The atomic application/receipt assumptions of earlier models cannot be silently combined with Phase 3B.6/3B.7's split stages to assert a stronger theorem.
- Preserve every known counterexample, failed candidate, prior REJECT, post-capture checker correction and source-conditioned custody limitation relevant to the assessment.
- The native Phase 3B.7 endpoint is idempotent; its branches share a pre-execution seed and one Store owner. Full declared exports/reports agree, while the full forensic dossiers differ.
- Never open retained SQLite evidence, modify historical captures, or read/change the runtime checkout's `.envrc`. No new captures or model executions are needed merely to restate already checked results.

## Deliverable

Create one central assessment with the complete obligation matrix, the closure verdict and exact remaining boundaries. Update the Phase 3 roadmap's stale current-status wording and links without rewriting original exit criteria. Reconcile only the manuscript passages affected by the closure interpretation in Markdown and LaTeX; rebuild and inspect the PDF, and change BibTeX only if actual references change. Preserve the inherited Markdown-only Section 7 and disclose it.

Include an audit verification record, exact input references/hashes, independent task and final reviews, and a draft PR. No production/runtime/model/checker behavior, dependencies or old evidence changes; no merge or release. The research worktree is reused on `codex/phase3b-closure-assessment` from the merged base above. Runtime builds, if inspection needs any temporary material, belong under `/tmp/alan-phase3b-closure/`.
