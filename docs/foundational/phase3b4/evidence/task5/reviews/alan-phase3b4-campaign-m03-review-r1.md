# Scoped M03 producer completion review — APPROVE

## Verdict and reviewed boundary

**Spec compliance: APPROVE. Task quality: APPROVE. Overall: APPROVE.**

The two requested points are addressed: the missing genuine M03 producer parent is now retained in normal and race modes, and the P3 supplied-parent overclaim is corrected. Findings addressed: **2/2**. Unaddressed scoped findings: **0**. New findings: **0 critical, 0 important, 0 minor**. No new scoped breakage was found.

This review is restricted to runtime `2b8f6e19a013d2b918146a7193451843004d2f86..fc60391f48b46a8c8635afc46882db8338b6c93b` and research `5d7473735d550ac1550228e9991863c3601afde7..c122de0d5f51210b60d6b8af878a7ccd8a3f8c01`, plus the evidence necessary to check those additions. The original Task4 campaign remains independently APPROVED by `/tmp/alan-phase3b4-campaign-review-r1.md`; this is not a broad re-review, Task5 oracle/mutant verdict, full frozen-pair qualification, manuscript approval, merge approval, or claim of general compositionality.

Requirements read first: research `.superpowers/sdd/03-plan/task-4-brief.md`, including its explicit completion follow-up; then the appended `.superpowers/sdd/03-plan/task-4-report.md`. I checked frozen `docs/foundational/phase3b4/05-matrix.md:67` for the exact M03 parent and remaining contradictory-admission/mutant obligation.

## Finding 1 — missing genuine target-generation-two M03 parent: addressed

The original finding concerned a real message's requested target generation, not a generation field in a delegation claim. The new fixture satisfies that distinction.

- `internal/durabletest/phase3b4_generation_parent_test.go:22–31` constructs the ordinary public program before Create. It uses the existing P10 insertion, changes B's actual PC5 message effect to literal generation2, recomputes the retained program through the existing editor, and calls public `runtime.Create`. It does not create A2 or change production authority policy.
- `:36–46` advances the public controlled prefix using ordinary public Commit, lifecycle and message APIs, then calls public `AuthorizeIntent` and `DeliverMessage` for the exact q2 ID `["b",1,6,0]`. The unchanged helpers were inspected to confirm that these calls do not use direct SQL, internal-import escapes, retained repair or injected consumption to establish the history.
- `:54–78` checks exact target/generation, same payload/obligation as q, the actual B6/PC5 emit reduction, one retained q2 admission, and no q2 observation. `:81–97` separately checks the unchanged root/subject set, only A1/B1, open/unpaid principals, absent planned principal service intents, and zero endpoint contacts/applications. `:100` gives the new case the frozen campaign's `TestPhase3B4Campaign` prefix and records its explicitly controlled/outside-C classification and ordered captures.

I independently checked the actual retained records in both modes, rather than treating these assertions or the producer report as evidence by themselves:

| Mode | Namespace | Case SHA256 |
|---|---|---|
| Normal | `8a779d2ae382c51dc8992f74e5ac64ea63964a22696ad28275c8d7d350fb034f` | `62e3e9f602551984e3d7e630b1a851399a0c1d3a457f4ffd4d78a0b2c0b4e1ad` |
| Race | `205428d39443e86b1a8d2f3febc3ce20dd751f233bf11a6d1f77fefa1e54fd6a` | `0dca65ce497732d5eb35f7af206c1ed1ebecd93037b7983000e376714a5281e4` |

Each namespace contains a genuine source-B1 intent q2, source sequence6/ordinal0/revision1, kind message, target `a`, requested generation2. B5's actual successor PC is5; B6 emits exactly that intent and reaches PC6. The raw history intent, raw intent descriptor, typed operation, actual retained program PC5 effect and hash-addressed exported program bytes agree. q2 shares q's payload/obligation but has a distinct full operation ID; q remains the original principal BA request.

The q2 source-authorization attempt contains that exact intent descriptor and revision1. Its result matches `auxiliary-delivery.json`: `Allowed=true`, `Reason=authorized`, `IntentID=q2`, `Source={Run:B1, Epoch:1, Sequence:6, Revision:1}`. The retained admission descriptor names that same full intent/logical ID, source B1, target `a`, requested generation2, exact obligation and payload. Its retained result is `Accepted=false`, `Reason=target_generation_changed`, `Bound=A1`, empty `ObservationID`, and a null raw observation index. There is one q2 admission and no q2 observation. Bound=A1 identifies the current candidate run considered by a refusal; it does not establish accepted A1 delivery.

Both ordered captures already contain the emitted q2. The final snapshot adds only its source authorization attempt, one matching composition authorization fact, and the refused admission. All other raw tables are identical across that boundary, including history, programs, subjects, roots, runs, observations and charges. Each final capture contains exactly two generation1 runs and two corresponding roots. Both principals remain open/unpaid; neither planned principal service exists; both independent endpoints have zero contacts and applications. Each captured Inspect/InspectExecution/Replay outcome is successful.

I also checked the prospective principal locators against the actual retained programs and manifest: BA serviceA5/replyA8 and AB serviceB10/replyB13. The principal claim still names genuine q/B1 and planned BA serviceA5. The added auxiliary is not relabeled as a principal or treated as paid. This is a sufficient genuine public parent for the next M03 construction; it does not execute that hostile construction.

## Finding 2 — P3 supplied-parent overclaim: addressed

`docs/foundational/phase3b4/evidence/task4/report.md:65` now identifies the supplied original P02-v2/P10/equal-total/cut/refusal parents, explicitly says the original corpus did not contain M03's real generation2 message, and points to the separate completion. It expressly leaves attacks and checker mutants to Task5.

The completion account at `task4/report.md:95–122` and `task4-completion-m03/report.md:3–30` preserves the prior APPROVE, distinguishes this incomplete outside-C producer prefix from principal completion and a live generation transition, and leaves sound-checker rejection/generation-join-mutant acceptance unclaimed. The README handoff and both manifest trees agree with the corrected narrative. The old review artifacts are retained byte-for-byte rather than rewritten to imply the parent already existed.

## Source, evidence and failure preservation

The runtime commit adds exactly one 101-line integration test file; no production implementation, constructor, dependency, public API, format13/14 interpretation or `alan/request-authority/v1` meaning changes. The research diff consists of the separate completion package and exactly the original Task4 README/report/manifest clarification. No old deterministic export appears among changed paths.

Independent mechanical checks passed:

- **328** Go/module source files match the retained source archive, current read-only source bytes and exact completion commit. All **327** original source-file hashes remain identical; the sole added source is the new fixture.
- All **67** original Task4 evidence files match the preserved pre-completion baseline hashes against research5d. **64** remain byte-identical now; exactly README, report and manifest are intentionally changed. This preserves the original case index, archives, costs, identity maps and source archives without regeneration.
- All **27** completion manifest file hashes, the original Task4 manifest hashes, and all **94** raw archive file members were checked: normal36, race36, apparatus-red22. Both selected case descriptors and their case hashes agree with the index.
- All **5** retained original review/check artifacts match their original `/tmp` bytes.
- All **4** normal/race ordered runtime captures were copied to disposable storage before SQLite inspection. Each is checkpointed and has no archived runtime WAL/SHM; any present companions would be copied as well. Query-only immutable access compared all **41** raw tables per capture against their raw JSON export: **164 table comparisons**. Archived databases were never opened in place. No SQL write created or altered a history.
- Both final focused test command logs report success: normal2.097s and race3.800s. The affected integration Staticcheck log is empty/successful. All **3** successful commands retain exact argv/cwd/environment and the complete **328-file** completion source map. The source map also resolves the case-local provenance records; their recorded base HEAD reflects tests run before committing the new file, not missing source.
- The initial failed probe is preserved as a failure. Its command map reconstructs from the original327 source files plus the archived red fixture. I inspected its actual generation refusal, only A1/B1 and normalized retained roots, as well as the failing `principal/root set changed` log. The faulty comparison was against unnormalized provision roots with empty IDs/version and different ordering; final `:81` compares actual retained roots before/after while separately checking subject and run membership. The failure is not counted as a policy rejection or passing qualification.

## Checks actually run in this review

- `python3 /tmp/alan-task4-m03-scoped-review-check.py` — PASS. Disposable read-only tar/hash/git-blob/JSON/copied-SQLite verification; structured results in `/tmp/alan-task4-m03-scoped-review-check.json`.
- A separate read-only Python inspection of the preserved red archive/log checked the actual refusal and root-normalization assertion cause.
- Runtime `git diff --check 2b8f6e19a013d2b918146a7193451843004d2f86 fc60391f48b46a8c8635afc46882db8338b6c93b` and `gofmt -l internal/durabletest/phase3b4_generation_parent_test.go` — exit0, empty output.
- Research `git diff --check 5d7473735d550ac1550228e9991863c3601afde7 c122de0d5f51210b60d6b8af878a7ccd8a3f8c01` — exit0, empty output.

No Go suite, race suite or Staticcheck command was rerun by this reviewer: the inspected exact-source command records and actual preserved evidence answered the scoped questions. No checkout, index, branch or Git state was modified; no agents were dispatched. Only disposable `/tmp` review artifacts were written.

The review script initially needed corrections for its own assumptions/schema handling: composition authorization facts accompany authorization, copied WAL-mode read-only SQLite needed immutable access for these checkpointed/no-WAL copies, BLOB text required decoding, and Go's successful `ok` line contains padding before the tab. Those were review apparatus corrections, not product/evidence failures. The final check passed without relaxing any requested semantic assertion.

## Outside-scope observations and remaining obligations

No new outside-scope defect was identified. Existing boundaries are carried forward, not re-adjudicated:

1. Task5 must construct the dense contradictory Accepted/Bound=A1 admission and all matching observation/result/index mirrors while retaining the real generation2 intent/program, demonstrate sound-checker rejection at the generation join, and demonstrate the precisely weakened mutant's acceptance. **No Task5 attack or mutant ran in this producer addition.**
2. Task5 still owes the broader independent raw-history/principal/custody oracle, attacks, mandatory mutants, model-support classifications and final frozen-pair verification. The original full80-case campaign was not repeated for this one new fixture; its next complete selector includes the new name.
3. N10e/f's live-generation/revision construction obstructions, promoted14/16's ancillary internal-writer-only compatibility route, and Task6 manuscript reconciliation remain as disclosed. This new message-generation parent does not resolve a live generation transition or change those boundaries.
4. Replay effect-freedom, validation-before-authority-renewal/failed-Open guarantees, and generic unknown-mutation no-retransmission policy have no changed production path here; this scope does not universally re-prove them. No paid service, distributed ownership, multihop delegation, dependency, release or merge is introduced.

**Assessment:** the two requested completion points are satisfied with bounded truthful claims, meaningful causal assertions, separately preserved evidence and the genuine public M03 producer parent. APPROVE this scoped addition.
