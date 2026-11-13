# Whole-branch review instructions

Read the provided whole-branch diff package and the Phase 3B.6 contract. This is the single broad final review of the new research package, integration, manuscript and evidence. Use the Senior Code Reviewer procedure in the requesting-code-review skill, including explicit strengths, Critical/Important/Minor findings, assessment and a list of behaviors declined for review with reasons. Return explicit APPROVE or REJECT for specification and quality; prior reviews are evidence, not an instruction to agree.

The ledger is `progress.md` beside this file. It has no parked findings or deferred-minor items. Preserve the initial rejected design claim and Task 3 correction in assessing the final account. Task 4 implementation, qualification, source review and PDF inspection are complete. The plan's remaining closure boxes include this very review, archival of its verdict and draft-PR delivery; those subsequent actions are not yet represented as done.

Read-only review: do not change checkout, index, branch or source; write only the requested review report and disposable read-only probe artifacts under `/tmp/alan-phase3b6/`. Do not spawn subagents. Do not rerun tests already evidenced on identical source. If an unanswered concrete correctness doubt requires execution, run only the focused check needed and report its exact result. Do not open archived SQLite. Review the complete branch diff in bounded chunks; absence from truncated tool output is not absence from the source file.

Use `docs/foundational/phase3b6/evidence/workflow/` for preserved task briefs, implementer reports and independent task reviews. `evidence/task4/verification.json` binds the final executable inputs and grouped qualification. `evidence/task4/pdf/` binds the source-reviewed PDF and inspected renders. Read owning source for consequential claims rather than relying on a prior report alone. State what you actually checked and remaining uncertainties without promoting finite or source-conditioned evidence into a general theorem.

Write the complete review to the supplied report path; return only explicit verdicts, path and a concise findings summary. Include the reviewed base/head and any focused probes. No merge, push or external messaging.

## Binding global constraints

- Deliver an independently reviewed draft research PR; do not merge or release.
- Preserve old sixteen authority-service configurations/order, State/Event JSON/enums, four historical deterministic exports, default command scope, old ledger/progress outputs, Phase 3B.5 model/audit/progress exports, and all ten old revocation concrete report bytes. New behavior is opt-in.
- One selected BA service boundary in a two-principal context, one shared Store-owner loss/reopen, one allowed physical gate/contact/application, one exact b revocation and optional refused fresh gate. Distinct c survives. No whole-pair recovery or general progress theorem.
- Keep historical allow/debit/gate, volatile opportunity/response, external packet/application and local retention/consumption separate. Endpoint evidence is not an import authority or recoverable result.
- No production runtime changes, new runtime campaign, dependencies, automatic retransmission or unrelated refactors. Never open archived SQLite; hash it as bytes.
- Independent ledger must not call operational successor, authority, transition, classification or observation helpers. Raw, continuity, physical, protocol/cut attribution and model support verdicts stay distinct.
- Work only in `/Users/jacob/.codex/worktrees/phase3b4-research/alan-research`, branch `codex/phase3b6-effect-observation`, base `df840103aa3a85e33b4b95ddcf2fe3df071cb024`. Do not edit other checkouts or `.envrc`.
- Offline environment: `GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b6/go-build`; use `/tmp/alan-phase3b6/` for temporary data.
- One implementation agent at a time; each task commits, self-reviews and receives independent spec/quality review. No worker subagents. Preserve real failures, exact commands/exits and source identities.
