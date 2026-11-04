# Independent public-runtime review

Reviewer: a fresh `mutual_runtime_review` agent, read-only, supplied the task
brief, implementation report/diff, frozen contract and raw artifacts. This is
an independent automated review, not external peer review or compiler proof.

**Verdict:** bounded specification met; no critical or important findings.
The reviewer read all four fixture files and owning public Commit, DeliverMessage,
AcknowledgeMessage, Open and runner-report implementations. The declared matrix
has exactly two ordinary roots and two original request/result obligations,
with 16 cases per normal/race run.

A separate read-only raw-fact scan covered all 32 exported cases: exact case
matrix, two roots, four prefix commits, exact program hashes/bytes, accepted
admission-to-intent-to-observation links including generations, unique selected
histories, response identity/value causality, serial current-Receive mismatch,
specific refusal reasons, final holds, and reopen stability. Servicing cases
contain 16 further steps and two dispatches. Serial reopened Run leaves every
inspected table unchanged; reopening changes the run epochs in these artifacts.
Quiescent-but-unpaid is supported by actual unconsumed requests and absent
result-consumption links, independently of report bookkeeping.

Two minor assertion findings were closed:

- A refused selection previously accepted any error. It now requires the
  current public `ineligible observation` error; the facade exposes no sentinel
  for this condition.
- Correlation checks now include originating intent source/target generations
  and selected-history generation as well as subject/sequence.

The implementer regenerated all exports, normal/race logs and the raw-fact audit,
and reran focused Staticcheck successfully. **Scoped re-review approved** the
fixes, verified that only those assertions changed, checked all archived exports
against their regenerated counterparts, and confirmed the 32-case audit and
final source hash `4b48dc9dcb1c8b82bb006dfcf2139457c5e41f58b4ab3d19419472a01abfceba`.
No task finding remains open. The parent additionally ran the portable archived
audit and compared its result byte-for-byte with `runtime-audit.json`.

No Go suites were rerun by the reviewer. Empty build/lint output alone is not
exit-status evidence; the implementation report records the observed statuses,
and the command manifest distinguishes their provenance.

Excluded explicitly: all-schedule/fair progress, independent-owner or process-kill
recovery, model receipt-release ACK correspondence, external mutation atomicity,
retries, full compiler/model refinement and root termination. The clean whole-Store
reopen renews both roots; both final local hold inputs intentionally remain unpaid.
