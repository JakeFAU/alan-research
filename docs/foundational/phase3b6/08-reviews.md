# Phase 3B.6 reviews and scope decisions

| Reviewed work | Reviewed head | Independent disposition |
| --- | --- | --- |
| Design and contract | [Design review record](evidence/preflight/design-review.md) | Initial REJECT on the claimed continuous barrier occupancy. Focused rereview approved the revised causal claim. |
| Task 1 operational model | `9a28d2ca4a4e031f1c98860154bbf0ed4aa7fb6e` | APPROVE for specification and quality, with no actionable findings. [Review](evidence/workflow/task-1-review.md). |
| Task 2 independent ledger and observations | `77f897e0c8bde7c1b3ea736a19f26029debb588e` | APPROVE for specification and quality, with no actionable findings. [Review](evidence/workflow/task-2-review.md). |
| Task 3 source-conditioned witness relation | `b7ed93186a16a0a52637de7b7e3a13a5fd37e523`, followed by correction `566b5e3ec1a2b23f6c66bc025a0e0ad447a3b05f` | APPROVE for specification and quality. The scoped rereview confirmed the minor inventory correction and identified no new findings. [Review and rereview](evidence/workflow/task-3-review.md). |
| Task 4 manuscript sources and grouped qualification | `c3665e68e706065126c2ecb2971a04fbb72ebd61` | APPROVE for specification and quality, with no actionable findings. This verdict precedes the PDF and whole-branch review. [Review](evidence/workflow/task-4-review.md). |
| Final whole-branch review | `8450871e0b0302681ce8b4af9faf505883b203aa` | APPROVE for specification and quality, with no actionable findings. Ready for the draft research PR; no merge or release authorization. [Final review](evidence/workflow/final-review.md). |

The initial design inferred continuous endpoint barrier occupancy until owner death. The pinned endpoint can exit its wait with a cancellation or timeout error before death. The accepted claim is that the endpoint commits before its acknowledgment and cannot emit a successful authentic receipt before the controller's post-death release. Both acknowledgments precede death. Their relative order remains unobserved.

The Task 3 serialized trust statement originally counted four archives. Two unsigned archives hold the four selected dossiers. The correction refreshes the serialized reports and corresponding digests. The review records no change to semantic acceptance. [Correction and report comparison](evidence/task3/review-i1/).

## Material scope decisions

| Decision | Evidence boundary |
| --- | --- |
| Model one selected BA service across one owner loss and reanalyze the loss and N07 protocols. | Full reciprocal progress and an experimentally matched v2 ambiguity pair remain open. |
| Use a declared sealed-dossier process anchor for the loss runtime acknowledgment's Store and token. | Child Open carries neither field. This external premise confines the relation to the declared dossiers and supplies no general process authentication. |
| Qualify unchanged ordinary packages during Task 3 and affected packages after the reviewed freeze. | The complete package inventory and dependency/source comparisons close coverage across two invocations. An executable or relevant input change requires requalifying the affected work. |
| Exercise the causal mutant end to end and the renewal mutant at the exact-renewal helper. | Inherited raw and physical checks already reject the malformed renewal dossiers. The helper result tests the specified comparison of run and table state. Inspection Settings and Roots are outside that contract. |
| Preserve the existing manuscript representation gap. | The edited Abstract, subsections and Conclusion align across Markdown and LaTeX. The inherited Markdown-only Section 7 remains absent from the TeX and PDF. |

The controller built and inspected the PDF, preserving the inherited layout warnings and Markdown-only section gap. The final reviewer independently checked the changed code, source/evidence joins, qualification records and selected PDF pages. [PDF record](evidence/task4/pdf/visual-review.md) and [workflow records](evidence/workflow/). No parked or deferred finding remains.

The controller retains each unassessed domain from the final review within the already declared contract. These are limits on claims, not waived correctness findings:

| Domain not assessed by final review | Controller disposition |
| --- | --- |
| Cryptographic archive authenticity and actual execution chronology | Honest pinned production and sealed provenance remain premises; the process-anchor ruling supplies no authentication theorem. |
| Matched v2 equal-local zero/one runtime pair | Remains the next experiment under the one-service scope ruling. |
| General recovery, pair completion, arbitrary-agent composition, distributed ownership/revocation, successful retry and forensic import | Remain outside the frozen contract; no result is asserted for them. |
| Arbitrarily constructed model/ledger values or arbitrary producer protocols | The documented reachable, post-admission domain and two pinned protocols remain the API/checker domains. |
| Universal filesystem/environment noninterference | The grouped-qualification ruling establishes the recorded source/dependency coverage and focused dynamic-read boundary only. |
| Entire inherited dissertation and all unchanged rendered pages | The edited passages and named renders are checked; the older representation gap and warnings remain disclosed. |
| Live GitHub parent and PR state | Controller rechecks these during delivery; this is an external-state check rather than a source-review conclusion. |
