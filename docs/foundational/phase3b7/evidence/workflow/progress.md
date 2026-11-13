# SDD ledger — plan: docs/foundational/phase3b7/phase3b7-plan.md

## Identity and authorization

- Research branch `codex/phase3b7-matched-capture`, base `d8645d739445d4b09c0d38dca4454ca6297fdf22`, existing isolated worktree `/Users/jacob/.codex/worktrees/phase3b4-research/alan-research`.
- User explicitly asks continuation and the genuine matched zero/application experiment. AGENTS prefers subagent-driven execution and push/draft PR. No merge/release authorization.
- PR46 and PR47 verified merged; current main's tree equals previous Phase3B6 delivered tree `aa0fb1acc3acb037d488b40b94bfdcdf25c49b3c`.
- Runtime source is read-only at `/Users/jacob/dev/personal/alan`; frozen commit/tree are in contract. Never inspect/change untracked `.envrc`.
- No new runtime capture has run. Recon reports in `/tmp/alan-phase3b7/`; design reviewer `/root/matched_design_review` currently running.

## Rulings

Ruling: Execute the specified bounded experiment without another design-permission round — the user has supplied the goal and asked to continue, and developer instructions require routine choices within that authorization — an unsuitable protocol choice would cost rework, not an unapproved deployment.

Ruling: Use one sealed closed seed and sequential isolated endpoints with identical URLs for a fresh pair — this removes incidental namespace/digest renaming while preserving genuine causal histories — the result is limited to forked initialization and does not cover independent provisioning.

Ruling: Compare full exported JSON structurally with all arrays ordered, and report byte equality separately — this predeclares an observation stronger than the finite model's projection — incidental export-order differences could falsify the richer match and must be retained.

Ruling: Preserve the current isolated research worktree and create only disposable runtime builds under /tmp — source inspection and capture need no production checkout edits — reproduction requires the pinned Git objects and offline module cache.

## Preflight conflict scan

| Scope | Produced / consumed or internal check | Finding |
|---|---|---|
| Task 1 / Task 2 | Versioned raw pair schema, source/capture/control seals consumed by checker | Task 1 report must fix concrete schema before Task 2; old anchors unchanged. |
| Task 1 / Task 3 | Frozen producer, no Task 1 qualification run; Task 3 runs normal/race pairs | Consistent; barrier control tests are separate development tests. |
| Task 2 / Task 3 | Reviewed comparison and freeze before genuine captures | Consistent; checker fixes after evidence require a new source qualification, never raw edits. |
| Task 1 | Test-only overlay, closed seed, dual holds, public recovery tail | Consistent; no production source or old campaign modifications. |
| Task 2 | Independent raw validation versus producer assertions; synthetic test corpus | Consistent; errors and unsupported domain stay separate from equality verdicts. |
| Task 3 | Controller captures/tests; worker evidence/manuscript; review and delivery | Consistent; no broad suite rerun claimed unless executed; historical full-suite evidence attributed. |

## Tool preflight

Observed: Go 1.27.4 darwin/arm64; Staticcheck 2026.2.2 (0.8.2); Python 3.9.6; Tectonic 0.17.0; Poppler 26.11.0. Existing offline module cache `/tmp/alan-phase3b5/modcache` present. No tool install or dependency mutation.

## Tasks

- Task 1 complete: commit 7ea8898b, independent spec/quality APPROVE.
- Task 2 reopened after genuine capture: fix90724cce committed; /root/matched_checker_review performing round2 scoped review. Original REJECT/APPROVE and initial genuine-check failures retained.
- Task 3 active, base6dc16ecf: exactly one normal and one race pair captured/sealed; strict export/report byte equality and zero/one world evidence pass, full checker awaits round2 correction/re-review. No recapture.

## Design gate and baseline

- Design reviewer `/root/matched_design_review`: APPROVE for implementation; report archived under evidence/preflight. Added explicit scheduler-lifetime and database/export trust limits to contract. No user clarification required.
- Baseline focused model/correspondence suite exited 0: 228 test/subtest passes, seven package passes, three no-test package skips. Optional external-corpus tests were not enabled; do not equate this with the historical 419-event corpus qualification. Full log/command archived in preflight.

## Task 1 dispatch

- Base `62dc6c299007abd40d3729827489c7cfb6e68f49`; implementer `/root/matched_producer` (gpt-6.1-sol, high). Brief and report paths are in this workspace. Task1 scope: producer and focused controls only; no genuine campaign.

- Task1 checkpoint: permissive barrier-control implementation compiled and failed meaningful EOF/invalid-byte/cancellation/unreleased/identity assertions; corrected controls and owner/process/seed/recovery code compile in disposable runtime checkout. This is implementer-reported pending review, not genuine capture evidence.
- Controller prepared uncommitted temp command recorder `/tmp/alan-phase3b7/run-command.py`; Python compile check passed. This will be reviewed as Task3 tooling before use.

- Task1 published concrete v1 schema in `02-protocol.md`; source-conditioned closed seed/dual barrier shape matches contract. Controller refined Task2 file map by responsibility before dispatch; no semantic scope change.

- Task1 committed `7ea8898b4ad9df31713d4061a691ac26ad78723b`; independent reviewer `/root/matched_producer_review` (gpt-6-astra, high) running against exact range 62dc6c29..7ea8898b. Controller verified current overlay hashes equal tested prepared-source manifest; no actual capture yet. Scratch report corrected to state environment files, including .envrc, were untouched (no contents read).

Task 1: complete (commits 62dc6c29..7ea8898b, review clean). Independent reviewer `/root/matched_producer_review` explicitly APPROVE spec and quality, no findings. Actual captured pairs still absent. Controller verified nine named controls / eleven pass events per normal/race mode and empty Staticcheck output.

- Checker file-map refinement committed separately at `f238b778ac786d3efbd86f94e98068439b9e4ed2`. Task2 base is current HEAD f238b778; brief generated, next worker pending.

## Task 2 dispatch

- Base `f238b778ac786d3efbd86f94e98068439b9e4ed2`; implementer `/root/matched_checker` (gpt-6-astra, high), report `task-2-report.md`. Synthetic checker controls only, no genuine matched captures. Producer source unchanged and approved; no parked findings.

Ruling: Expose unchanged correspondence binding/continuity helpers through a narrow new adapter — the existing independent raw validator already checks the full authority tuple, and duplicating it would introduce drift — this reuses prior authority-checking logic, so the new experiment is not an independent reimplementation of that validator; new outer-state, source, protocol and world checks remain separate. Existing pinned entry points and report bytes must remain unchanged.

- Controller temp tooling controls: command recorder preserved direct exit 7 (explicit expected-exit control); archive helper retained all 256 opaque byte values and both declared synthetic members, with raw unchanged. These are labeled apparatus controls, not runtime experiment evidence. Scripts remain under /tmp until Task3 packaging/review.

- Task2 checkpoint: implementer reports normal/adapter race/Staticcheck pass and four historical effect-report byte regressions unchanged. Full checker race still running; no commit or independent review yet. Malformed-input cases found and fixed include Settings/Roots, duplicate seed sources, wrong test selection, and nil-dependency/empty-argv panics; failed logs will be retained. No runtime/producer edits or real capture.

- Task2 committed `eca3b3bee03b32c7bfffc65cba124d3ba638446c`. Reviewer `/root/matched_checker_review` (gpt-6-astra, high) is checking exact f238b778..eca3b3be range. Task2 report records ordinary/race/Staticcheck results, four old effect reports byte-identical, meaningful malformed-input failures/fixes, and no genuine capture. Controller prepared `/tmp/alan-phase3b7/freeze-inputs.py` for later exact source/acceptance inventory; only Python compile check so far, no freeze asserted.

## Task 2 review rejection and fix round 1

- Reviewer `/root/matched_checker_review` REJECTED spec and quality at eca3b3be. Four coherently resealed synthetic probes returned Supported=true: missing zero snapshot coverage fields; both final reports truncated to Error/UnknownExternal; contradictory native authorization Subject/Epoch/Sequence; later overriding test-run flag. Full report `task-2-review.md`; probes `/tmp/phase3b7-review-probes/`. No genuine experiment ran.
- Fix round1/5 dispatched to original `/root/matched_checker`; base eca3b3bee03b32c7bfffc65cba124d3ba638446c. All four findings required, original rejected commit and review retained. Dependent qualification remains gated.

- Fix round1 committed `6dc16ecf36441a99005a62eba4cee2418f5f5d72`; implementer report has original probe pass, full checker/command ordinary and race pass (216.329s checker), scoped Staticcheck empty. Native schemas verified against frozen source; zero uses Requests:null and Applications:[]. Initial test syntax and wrong-path command failures retained. Scoped reviewer `/root/matched_checker_review` now checks eca3b3be..6dc16ecf, report target `task-2-review-r1.md`.

Task 2: fix round 1/5 (4 addressed, 0 open — complete endpoint snapshots; complete Run schema; native authorization identity; effective test selector; commits eca3b3be..6dc16ecf).

Task 2: complete (commits f238b778..6dc16ecf, review clean after retained REJECT/correction). Reviewer `/root/matched_checker_review` explicitly APPROVE spec/quality in `task-2-review-r1.md`; no out-of-scope findings.

## Task 3 start

- Base `6dc16ecf36441a99005a62eba4cee2418f5f5d72`. Controller owns freeze/captures/qualification and tools; a later fresh worker owns results/manuscript from actual evidence. All required precursor review gates passed.

- Qualification tools committed at 19303940; exact source/acceptance freeze has 627 inputs. Producer recheck verified all 637 runtime files and four overlays against Task1 manifest; checker built and binary digest recorded in qualification/checker-before.json. About to run first normal v1 pair; no preceding genuine attempt.

## First genuine captures and checker rejection

- Exactly one normal v1 pair and one race v1 pair ran against frozen sources. Producer commands exited0 (4.035817208s normal, 11.621542834s race). Raw archives sealed before independent checking: normal SHA256 0d2b423ea6f611be28595ac71d973507c7ccf6d2444e68fc05c28af61fdf4b65; race19be235cec79f5555eaaa4d92263af884f36c9fa4ef45ff814f81de9bc4c6292.
- Both frozen-checker reports run twice, each direct exit1 and repeated bytes identical. Seal/Source/Raw/World pass; all four raw exports, seed copies and final run reports byte-identical within each pair. Protocol rejects passive external.authorized AuthorityAdmissionOrdinal1 because checker demanded0; Model consequently uncertified. Initial checks summary archived; raw unchanged.
- Root cause read from frozen runner.go:203–215: transport wrapper sets work.event.AuthorityAdmissionOrdinal=op.authority.Ordinal before transportExternal, while Task2 review considered only the event constructor. This is a checker-source-mapping error, not a full-observation mismatch. Original captures, old checker binary/freeze and failed reports must remain unchanged; no recapture. Task2 correction/re-review reopened before final witness claim.

Ruling: Correct the checker's native authorization ordinal from the complete frozen call path and re-check the same sealed captures — the scheduler wrapper supplies the selected gate ordinal1, so accepting0 was a source-mapping error — final acceptance uses a checker corrected after capture, and the initial rejection plus frozen source/binary identity must remain visible; no observation criterion or raw member may change.

- Round2 fix dispatched to original `/root/matched_checker`, base193039406557cb37be815d5f846b22a75d48fde6. No producer/contract/raw changes and no recapture. Exact ordinal1 binding plus zero/wrong-ordinal negatives; controller will freeze a corrected checker separately after review.
- Both real endpoint receipts identify the fixture Profile as idempotent. Do not describe this concrete pair as a non-idempotent endpoint experiment; local zero/one ambiguity under exact-root fresh refusal remains the selected question.

- Round2 fix committed90724cce9082f1ca978d785efa743e851a344617. Source call-path trace/native excerpt positives, exact ordinal1/worker binding, synthetic0/2/-1 negatives; full affected ordinary21.582s, focused boundary race39.622s, scoped Staticcheck. Original binary/raw/reports verified unchanged. Reviewer `/root/matched_checker_review` checking exact19303940..90724cce. Native testdata is embedded under research/matchedcapture/testdata and will be included in corrected freeze.
- Note: Task2 report under this plan's .superpowers scratch is tracked in90724cce; archive its exact final contents into evidence/workflow before removing this plan workspace at finish. No other plan workspace is in scope.

Task 2: complete (commits f238b778..90724cce, review clean after retained original REJECT and post-capture mapping correction). Round2 reviewer APPROVE spec and quality; source trace and exact ordinal1/worker join verified.

- Corrected checker frozen with629 inputs, built and checked each unchanged raw pair twice: four direct exit0 results; repeated report bytes identical. All verdicts valid, all four exports and final reports byte-identical, zero0/0 versus applied1/1 contacts/applications. Model trace lengths zero1/2/4 and applied4/5/7. Raw archives/member sets and both old/new checker binaries unchanged; corrected evidence archived. No recapture.

- Task3 documentation worker /root/matched_results dispatched at90724cce, assigned results/runbook/reviews/manuscript/PDF only. Controller starts affected historical corpus ordinary/race, build/module/format/scoped and baseline global Staticcheck against unchanged corrected freeze. New checker full/focused race evidence is already recorded; no duplicate full new-checker suite.

- Task3 controller qualification complete: ordinary/race historical corpora each422 events,7 tested packages,3 no-test packages,no named skips; build/modules/scoped Staticcheck pass. Global diagnostics13/31 byte-identical to baseline. Current-source format clean; initial selection included an immutable historical probe and is retained. Package consistency verification passed,629 corrected inputs unchanged, two122-member archives and23 command records checked. Source/old-new binary/raw provenance unchanged. Remote main remains originalbase d8645d73; no research branch exists remotely yet.

- Task3 evidence committed9fca8cdb, manuscript/docs committed964d9f2e. Fresh reviewer /root/matched_results_review checking6dc16ecf..964d9f2e with explicit spec/quality verdicts. PDF44pages built/visually checked; original627 and corrected629 source inventories preserved. No qualification code changed after freeze.

Task 3: complete (commits6dc16ecf..964d9f2e, review clean). /root/matched_results_review explicitly APPROVE specification and quality; no actionable defects. Final whole-branch review and draft PR remain controller delivery obligations.
