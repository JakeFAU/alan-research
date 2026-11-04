# Task 3 raw-policy re-review — fix round 1

**Specification verdict: APPROVE within the scoped raw-policy re-review.**

**Quality verdict: APPROVE within that scope.** R1, R2 and R3 are **ADDRESSED**. The recipient-mutation coverage/name correction and verified-byte reuse concern are also **ADDRESSED**. No new blocking breakage was found in the six-path fix diff. The previous REJECT report and exact adverse evidence remain preserved; this review supersedes its unresolved status for this corrected candidate, without deleting or waiving the original findings.

Candidate identity: 23-file source manifest SHA256 `036e16e86dd06c63ad98342c98676409a5b0d8316bf4612c0e3794fcd70f00b3`; fix diff SHA256 `2a1b86a9bd2dc3f97a18e238317424adf4c2c76a3b48275a2ec117dfd5642893`; qualified checker binary SHA256 `8f3f5d19223546a20991f105624a70db60959479a255aeb56e3aa961641a176d`. The rejected candidate manifest remains `9ca9c230e7ebbe3185474fc5a52db07fda36af0a764f006733032bbe9d5477c6` on research base `31c8af68e818f97a842e790669fbc8f8c98dc8bd`.

I read the appended round report, the exact 944-line fix diff and source manifest, and inspected the actual corrected probe results, raw mutations, corpus outputs, command records and RED/GREEN logs under `/tmp/alan-phase3b3-task3-evidence/fix-round1/`. Review was limited to the six changed/new paths. Unchanged model/projection/progress/runtime code was not reopened. Parent verification of all current source hashes and all 85 evidence-manifest artifacts is attributed to the parent, not represented as an additional hash audit by this reviewer.

## Finding outcomes

### R1 — accepted message/input permission: ADDRESSED

`custody.go:129–141` now constructs the exact six-field original permission for an accepted ordinary message using its committed source revision, and requires equality with a retained original permission. This occurs independently of the existing source-authorization and observation-identity checks. Accepted local input uses revision zero.

The tests at `review_fixes_test.go:19–64` remove each of AB request, BA request, ready, BA reply and AB reply from both permission mirrors while preserving successful histories/results; they require `admission_permission`. The local-input test distinguishes an accepted revision-zero route from a wrong-revision route. The retained RED log shows these controls failing before correction.

The exact original reviewer mutant remains SHA256 `304ad2640f7664a6916fe4126d317e4ab583531ccf72991295ba1fb90984b4c5`. Actual `review-request-route.log` now reports `Valid=false`, edge `admission_permission`, and identifies the missing B→A `ba-request` tuple; `commands.json` records CLI exit1. This rejects the policy contradiction rather than merely detecting unequal permission mirrors.

### R2 — immutable source/authority continuity: ADDRESSED

`external.go:334–357` compares complete immutable metadata, program, subject, permission, origin, composition deployment/profile/root/static-membership sets, as well as the fixed run IDs/program bindings/initial projections. `CheckContinuity` invokes it before the existing append-only protocol-row checks. Sequence, dynamic projection and epoch are not mistaken for immutable source fields.

`TestImmutableSourceAndAuthorityContinuity` at `review_fixes_test.go:66` keeps each permission/program/deployment mutant individually raw-valid, then requires continuity rejection. The original copied reviewer probe in `review-continuity.log` retains `InitialValid=true` and `FinalValid=true` but now returns `immutable_source_changed`, detail `permissions`. Thus the fix checks cross-prefix conservation instead of suppressing a difficult prefix or changing its raw validity.

The revised phase ordering distinguishes ordinary final-before-reopened captures from process stopped-before-recovered-final captures. The actual normal/race corpus outputs have no continuity failure for their valid supplied histories.

### R3 — actual owner epoch and renewal: ADDRESSED

`external.go:209,225–236` now binds each process-scoped boundary's source RunID and epoch to the corresponding owner context before joining its gate ordinal. A supplied nonzero event PID must agree with the owner. `checkOwnerEvidence` at389 validates Open namespace/success/PID and the complete actual post-open run set, including program, revision, sequence and projection. `Owners` and `OwnerBindings` retain the paths/hashes and joins used.

`checkOwnerRenewal` at371 requires an open run to advance exactly one epoch and a terminal run to retain its epoch, with unchanged sequence/projection and immutable source. Fresh-owner context requires the stopped predecessor. Standalone `reopened/inspect.json` checks use stopped when present, otherwise the actual ordinary final predecessor; this matches the captured layouts rather than inventing an absent phase.

Actual `review-epoch.log` preserves `RawValid=true` for the historically bounded epoch2 gate, but reports `ExternalValid=false`, edge `boundary_owner_epoch`, while retaining fresh-owner epoch3. The original Request.Epoch remains distinct from actual owner/gate epoch. Tests at `review_fixes_test.go:159–252` include the stale gate/event pair, false non-renewal, false terminal renewal and missing fresh-owner post-open snapshot. The prior RED log retains their pre-fix failures.

I specifically inspected both corrected normal and race outputs for the two cases without initial snapshots. `generic-unknown` and `applied-unknown-revoked-before-fresh-gate` report first-owner `RenewalChecked=false`, empty prior path/hash and the explicit limitation; supplied Open/post-open binding remains checked. Generic unknown's fresh owner reports renewal checked from `stopped/inspect.json`. The controlled applied-unknown/revoked path uses its separate stopped→reopened inspection check and does not invent a fresh-owner process record. Cut10 retains A epoch2 and renews open B to epoch3. These qualifications are accurate and do not weaken the original owner-context finding.

### Actual recipient mutation/name: ADDRESSED

`mutations.go:53` adds a genuine recipient-subject mutation; the old holder mutation is explicitly named `altered-holder` at128. I inspected `raw-mutations.json`: the recipient changes a/1→b/1 in the lifecycle claim and matching decision request; the holder/root and external receipt rows are unchanged. Exact semantic field, primary-key/column coordinates and complete before/after cells are retained.

The new control correctly reports `lifecycle_origin` at order0, because its changed request contradicts the committed original intent. The report and README expressly describe original-source integrity evidence, not a claim to have reached the later recipient-policy branch. Existing authentic recipient refusal controls remain separate. This fulfills the requested correction without relabeling the older holder control or overstating the new rejection edge.

### Verified-byte reuse: ADDRESSED

`files.go:125–142` caches JSON/JSONL bytes only after the normal read and digest check. All later projection, external, continuity and source processing uses those cached bytes (`files.go:174–187`). The two discarded-error rereads are removed. Initial read/path/digest failures still propagate as apparatus errors. The inspected ordinary/race runs include the existing relocated-corpus test through this path.

## Verification evidence and remaining limits

The actual round command records contain 12 commands with their expected exits. I inspected the ordinary/race logs and counted 58 passing test/subtest events each, zero skips/failures and no race warning. Staticcheck exits0. Both complete corpus outputs retain 48 cases, five invalid final hostile histories, no external failure among raw-valid cases, no continuity failure and the sole named `immediate_selection_narrows_model_admission` mismatch. Unqualified CLI checks remain exit1 and explicit-expectation checks exit0; expected findings are not relabeled as agreement. The pre-context-layout failures remain separately preserved and are not counted as successful frozen qualification.

No new execution probe was necessary: the retained reruns of the exact three reviewer falsifiers, substantive new negatives and actual positive corpus outputs resolve the scoped doubts. I did not rerun a Go/model/race suite, open SQLite, modify a source/index/HEAD or archived witness, or spawn subagents. Only this re-review report was written.

Approval remains bounded to this corrected research checker and its declared generation-one, two-static-root profile13 evidence domain. Missing initial captures do not prove initial renewal. Unsigned source/owner/endpoint files do not establish authenticity or exposure completeness, and an immutable endpoint wire does not identify the current gate ordinal. The recipient mutation's first edge is source integrity, not a full recipient-policy mutant. Whole-Store model crash correspondence, profile14/general-program completeness, whole-module frozen regression, manuscript reconciliation and final delivery remain separate obligations. No broader theorem or runtime approval is implied by this re-review.
