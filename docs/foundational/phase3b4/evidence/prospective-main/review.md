# Phase 3B.4 prospective integration — final APPROVE

**Final spec/quality verdict: APPROVE for the exact prospective tree and bounded draft-delivery scope. Findings: 0 Critical / 0 Important / 0 Minor.** The original whole-branch report remains unchanged. The initial source assessment below is preserved chronologically; its evidence-pending status is superseded only by the final closure at the end.

## Historical provisional source assessment (preserved)

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


## Final evidence closure

I inspected the closed `summary.json`, all 15 `gates/*.json` records and their matching logs, `driver-results.json`, `tools-and-readers.json`, `fresh-main-and-branch-check.json`, progress chronology, fresh normal/race case records, independent verdicts, promotions and attack verdicts. The exact original prospective identity remains synthetic commit `5b4662fd682a220ca3065b15e86672b7e56567e7`, tree `44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`, combining runtime main `926a899b651fccc48be52a054612a92afcf1ddfd` with Phase3B.4 delivery `9dbcf3b9ee3047ab280ac150b05a184d44832007`.

The review's independent check command was `python3 /tmp/alan-phase3b4-prospective-review-check.py`, exit0. Its exact source and successful log are `/tmp/alan-phase3b4-prospective-review-check.py` and `/tmp/alan-phase3b4-prospective-review-check.log`. It checked file bytes/JSON only; no tests were rerun and no SQLite file was opened.

- All **15 gate record and log hashes match**. Every final gate exits0 without a watchdog expiry. All before/after executable maps agree; I independently hashed every mapped current file against each record: **351 prospective runtime inputs / 562 unchanged research inputs**. The research oracle remains the previously frozen executable source rather than being changed to accommodate the new runtime.
- Runtime ordinary, integration, integration-race, module verification/tidy-diff, build, format, ordinary/integration vet and Staticcheck, exact-public-parent promotions normal/race, and independent binding oracle normal/race are closed. Parsed actual JSON logs contain no failing event. Both integration modes explicitly pass `TestCompetingReadCannotSettleAdopterDebt/original`, `/adopter_query`, `TestRootCompletionsMembershipQueriesOncePerRun` and `TestRootCompletionsFreshLateReceiptSettlementAndConflict`.
- Independently recomputed each mode's **81 cases / 80 Stores / 163 valid captures / 40 principal-complete / 33 supported C / 89 contacts and applications / 11 acknowledged losses** from individual verdict files. All evaluated raw/external/continuity verdicts are valid. Both K01/K02 promotion comparisons are valid in each mode. Initial full-suite promoted-parent skips were satisfied by the separate explicit gates; the historical baseline expected-failure test remains an opt-in skip.
- Rechecked actual **P02-v1, P02-v2, P10, N02, N06 and C07** JSON and all declared file digests for those cases in both modes. P02-v1 preserves the adjacency refusal; P02-v2 joins exact q/admission/o/c/grant/e and retains cA2/eA4, with causal fact orders **3/4/5/14/16/17/18** and genuine endpoint ingress/committed records for A4 and B9. P10 retains cA2 binding at eA5 despite another same-peer selection. N02 preserves cA3/eA2 refusal and a zero consumption witness. N06 preserves first allowance while refusing the fresh gate for revocation.
- C07 preserves two declared acknowledged owner losses (normal PIDs88618/88622; race PIDs93391/93392) and the exact three Open renewal pairs **1→2, 2→3, 3→4** for both principals. At every pair I compared all non-run tables for equality and verified that each run changes only its permitted epoch. The original selection and padded serviceA7 persist. The independently checked external reports separately validate loss acknowledgments and report no continuity failure.
- Fresh exported attack verdicts in both modes reject M01 at `selection_before_service`, M02 at `first_authority_policy`, M03 at `admission_binding`, M04 at `fresh_gate_policy`, A02-filter at `receive_filter`, and equal-total substitution at `reply_result_custody`. This is fresh evidence validation with the unchanged checker, not a claim that compiled-mutant construction was repeated in this addendum.
- Both old-reader binary hashes independently match their provenance records. The fresh fetch/branch record still identifies research main9ddfe29a and runtime main926a899b, with clean reviewed branches and original562/350 executable maps. This review did not itself issue a network fetch.

The first prospective attempt remains an **environmental failure**: its preserved log shows `httptest` loopback bind rejected with `operation not permitted`, exit1, followed by the stopped driver. It is not a semantic RED or passing gate. The subsequent authorized same-source run supplies the closed results above. The provisional review withheld approval pending these results; it was not a rejection for a code finding. Previous whole-branch rejection/correction history remains preserved and is neither waived nor rewritten by this addendum.

## Packaging and factual delivery note

Read the controller's prospective `README.md`: it accurately separates the original scientific pair from this integration check and states that the actual15 jobs govern, with no repeated full research suite, model exports or compiled-mutant run. Its bounded interpretation is consistent with this review. Independently verified all12 compressed archive hashes/sizes against `package/archives.json`:31,967,186 bytes and10,988 manifest members. The controller's packaging record supplies the per-member source-byte comparison; this reviewer did not repeat every member comparison or open any database. This archives already reviewed evidence and introduces no executable delta.

## Final exclusions and disposition

The provisional Declined to judge list still applies. No performance claim, expanded model theorem, new general failure equivalence, public conflict-production schedule or distributed ownership guarantee is approved. Original scientific source identities and counts stay attributed to their original qualification; this separate addendum supports interaction with exactly the stated advanced main. Any later executable main drift needs a fresh scoped assessment and applicable evidence.

No corrective source work or repeated suite is requested. Evidence packaging and future draft PR delivery remain controller responsibilities; this report does not assert that those external actions have finished. The provisional report is separately preserved at `qualification/review-provisional.md`, and the original whole-branch report remains intact.

**Final specification verdict: APPROVE. Final code-quality/integration verdict: APPROVE.** There are **0 Critical, 0 Important and 0 Minor findings** for this exact prospective combination. Source-local membership reuse preserves the reviewed authority/evidence boundaries under the existing stable-cut assumptions, and fresh ordinary/race recovery/compatibility/oracle evidence supports the bounded interaction claim. **This approves reviewed draft delivery only, not merge or release.**
