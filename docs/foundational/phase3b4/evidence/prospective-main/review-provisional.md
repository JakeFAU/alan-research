# Phase 3B.4 external-main drift addendum — source assessment, evidence pending

## Scope and status

This resumes the same independent whole-branch review for one scoped integration question. It does not replace `/tmp/alan-phase3b4-whole-review-r1.md` or repeat that broad review. **Source assessment: no findings (0 Critical / 0 Important / 0 Minor). Final prospective APPROVE is withheld pending the closed qualification outputs.**

Runtime main advanced after the original approval from `1fdd4c9e1720bf60b344914bcbc3d0938ba8d678` to `926a899b651fccc48be52a054612a92afcf1ddfd`. The conflict-free prospective combination with Phase3B.4 runtime delivery `9dbcf3b9ee3047ab280ac150b05a184d44832007` is detached synthetic commit `5b4662fd682a220ca3065b15e86672b7e56567e7`, tree `44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`, at `/private/tmp/alan-phase3b4-prospective-runtime`. I independently read its HEAD/tree and clean status. Identity records declare 351 runtime executable inputs. Original scientific qualification remains the separately frozen pair; this prospective tree must not replace its identity retroactively.

Read first: `/tmp/alan-phase3b4-prospective-qualification/{identity.json,plan.json,progress.md,main-delta-review.diff}`. Then checked the actual committed root-source delta and surrounding call sites, membership query, Replay/Open transactions, finishRun lock, new helper tests, and the competing-debt test/helper lifecycle. No source/Git mutation, database opening, suite rerun or delegation was performed.

## Strengths and interaction assessment

- `internal/executor/composition_inspect.go:39` introduces only a function-local membership map. Its complete `RunID` key includes generation; membership is reused only after a successful query, and map presence distinguishes legitimate empty membership from an unseen key. It is lazily populated at the first original traversal position, does not escape the call and introduces no retained state.
- The stable-cut premise is supported by both production callers: `finishRun` holds `Store.mu` through the report (`runner.go:497`), while `validateCompositionStore` is reached under the provisioning/renewal or read-only validation transactions. Replay and Open preserve validation/activation order. No helper invoked here writes membership. Phase3B.4's fixed two-root membership makes this interaction narrower still.
- Root membership is not selection provenance or authority. The optimization does not cache a permission, grant, witness, current revocation, attempt, debit, receipt or complete report. Phase3B.4 q/o/c/g/e checks, first decision, fresh gate, v1/v2/15/16 interpretation and conflict promotion remain unchanged. Later report calls reconstruct late results/settlements and membership afresh.
- New root-report tests independently cover complete semantic reports, repeated/absent membership, full-generation keys, distinct stores/calls, retained SQL/cancellation errors, real dynamic child membership and late receipts/conflicts. Synthetic generation/reparenting helpers are explicitly labeled and do not pretend to be public positive histories.
- `internal/durabletest/competing_debt_test.go:38` now acknowledges and parks the exact adopter at `runner.received/pre-account` before injecting its settlement fault. The test retains both genuine response boundaries, distinct admitted reads, revocation, original/adopter/duplicate settlement assertions, sticky error checks and fresh-owner exactly-once accounting. Cleanup releases both pauses before cancelling/joining Run. This is a test scheduling correction, not a new production settlement policy or SQL-based positive Phase3B.4 bootstrap.

## Issues

Critical: none found. Important: none found. Minor: none found.

## Required evidence before prospective closure

1. Closed runtime ordinary, integration and integration-race gates on the exact detached identity, plus build/module/format/Staticcheck/required vet results and any skips/failures recorded accurately.
2. Fresh Phase3B.4 normal/race corpora and exact public K01/K02-parent promotions, followed by the unchanged independent research oracle in both modes. These must retain actual causal/raw/endpoint/continuity outcomes rather than infer compatibility from the small textual delta.
3. Before/after source/input guards for all 351 prospective runtime inputs and 562 frozen research executable inputs; explicit current scientific pair versus prospective identity in the final evidence.
4. Exact evidence paths/hashes and a final output summary sufficient to check the formerly reviewed counts/claims and relevant recovery/current-gate behaviors. No new universal or performance claim follows from this integration check.

The supplied plan's top-level generic `scope` sentence mentions full ordinary tests in both repositories and old exports. Its explicit prospective progress note and job list instead correctly omit an unchanged research full suite/model-export rerun. Actual command records, not that reused generic sentence, must determine reported work. This is a reporting interpretation, not a source finding or request for unnecessary repeated suites.

## Declined to judge for this addendum

- Independent P2 performance measurements, forecasting and evidence-publication claims in incoming main's docs/nested tooling: not modified by Phase3B.4 and outside this authority/recovery integration question; this report does not independently approve them.
- Exact outcome equivalence for faults/cancellation occurring only at a deliberately eliminated redundant SQL query: not claimed by P2's stable-call contract; retained query/cancellation/error boundaries and authority gates remain the relevant integration checks.
- Unsupported direct SQLite membership writes during a stable report, distributed ownership or persisted cache reuse: outside the cooperative owner/snapshot contract and not introduced by this delta.
- Repeating unchanged research suites/models or original broad claim review: explicitly outside this scoped addendum; previous bounded approval and exclusions remain preserved.
- Future push, PR state, actual merge/rebase or release: not completed or authorized by this review.

## Provisional assessment

No source-level obstacle to prospective qualification was found. Final spec/quality approval will be stated only after inspecting closed exact-tree evidence. Neither original branch/main is changed by this detached check.
