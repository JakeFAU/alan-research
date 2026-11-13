# Task review procedure

Read the provided task brief, implementer report and diff package; the package contains the full relevant diff. This is task-scoped spec-compliance AND code-quality review. Broad whole-branch review happens later. Do not pre-judge findings based on implementer rationale. Read each diff chunk once; inspect unchanged code only for a named concrete risk and document that focused check. No checkout/index/branch mutations and no subagents. Do not rerun tests already evidenced on the same code. If an unanswered concrete doubt requires execution, run only that focused check; never broad/race/repeated suites. Evidence absent from a truncated output is not absent from its source file.

Return APPROVE or REJECT explicitly for both spec and task quality. Name Critical/Important/Minor findings with file:line, technical consequence and required correction; note strengths with evidence. A cannot-verify requirement is an explicit gap, not an implicit rejection. Report what a controller must check. Write the full review at the supplied report path; return only verdicts, path and concise findings.

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
