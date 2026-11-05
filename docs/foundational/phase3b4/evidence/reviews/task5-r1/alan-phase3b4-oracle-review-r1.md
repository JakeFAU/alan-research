# Phase 3B.4 Task 5 independent raw-evidence and compatibility review

## Review status

APPROVE — independent Task 5 specification and code-quality review. Critical: 0; Important: 0; Minor: 0. The final ordinary research gate and both evidence-only deltas are closed and reviewed. This approval applies only to Task 5's bounded qualification and the exact heads below; it is not approval to merge/release or a claim that Task 6 is complete.

Research base: `c122de0d5f51210b60d6b8af878a7ccd8a3f8c01`.
Research executable: `0c303c3d43d6582fed0a9caac24f50e71d844a19`.
Runtime executable: `fc60391f48b46a8c8635afc46882db8338b6c93b`.
Final research evidence head: `b940f0d2a30f3d0362cc24ea2d7c0bbfdf4d74cc`.
Final runtime evidence head: `7c6d092aa293162c357a4e3749b8abbff21b8af2`.

Read the Task 5 brief before the implementer report, then the frozen contract, model relation and matrix. Reviewed the provided full/code diffs, using targeted source context only to resolve the changed predicates and their enforcement. Both checkouts and Git state remained read-only. No subagent, source fix, broad test run, race rerun, archived SQLite opening, commit or push was performed. Reviewer scripts and this report are confined to `/tmp`.

## Specification verdict

APPROVE. The frozen Task 5 implementation and final sealed evidence satisfy the reviewed independent-oracle boundaries. Final gate and packaging checks have closed. The manifest-commit chronology differs from the literal plan-brief wording and is explicitly disclosed below; under the controller's stated packet-authority interpretation, this is not an unresolved qualification defect. The supported conclusion is finite, conditional qualification on the named source pair, not general refinement, inevitable completion, hostile-archive authenticity or distributed ownership.

The three verdict domains stay distinct: raw retained-policy/reduction validity, independent endpoint contact/application evidence, and supported projection to the old model. Principal custody is explicit for the two-original cases; the one-original authored adapter controls deliberately do not claim two-principal custody or supported correspondence. Missing Store construction at N10h remains no raw verdict, rather than a policy refusal. Unsupported projection does not excuse corrupt raw evidence.

## Findings

### Critical: 0

None.

### Important: 0

None.

### Minor: 0

None. The former final gate/evidence-delta prerequisites have been met.

## Verified strengths and decisive evidence

1. **Exact historical binding, not convenient latest-message selection.** `research/delegationcorrespondence/binding.go:6` derives q from the immutable grant; lines 9–24 enforce the canonical observation, strict earlier consumption, actual same-run history, revision and source/admission joins. The ordered interpreter independently validates real reduction and source availability. `check.go:740` retains accepted-grant-before-decision and exact service identity; `check.go:755` selects v2 while lines 760–779 preserve v1 adjacency. The removed duplicate delegated c<e check does not weaken enforcement: the authoritative delegated predicate remains at `binding.go:11`; delegated custody checks the validated witness's original at line 64; the distinct own-path check remains at line 67. The compiled M01 result demonstrates that this is now the intended nonredundant historical predicate, rather than a redundant deletion falsely counted as a kill.

2. **Actual q and q2 evidence.** In final-normal P10, the genuine B1 original is selected at A2, the genuine B6 auxiliary at A4, and the allowed service is A5 with witness naming B1/A2. In N01, B1 is admitted but unconsumed; ready B5 is selected at A2 and B6 at A3. Actual A4 remains refused with `authority_request_not_selected` and empty request observation/consumption witness fields. Both q identities arise from retained emitted message/program rows. Explicit custody reports P10's two principals consumed and q2 separately unpaid; N01 has no endpoint contact/application. These are actual distinct IDs, not suffix or payload substitutes.

3. **Future selection and current revocation are different boundaries.** N02 retains actual service A2 and q selection A3; its first refusal has empty selection witness fields, preserved across later captures. N06 retains lawful first reservation/allow at composition order 16, genuine revoke at order 20, and a fresh gate at order 22 refused with `authority_revoked`. Its endpoint logs contain no application. The M01 and M04 attacks preserve these source coordinates and change the complete required dependent mirrors; rejection occurs at `selection_before_service` and `fresh_gate_policy`, respectively.

4. **Generation falsifier is a real target2 message.** Final-normal M03's B6 committed intent and admission descriptor request A generation 2. Genuine result is Bound=A1, Accepted=false, `target_generation_changed`, no observation ID. The attack retains target2 while changing the result to Accepted=true/Bound=A1 and adding the matching observation/index. Sound rejection is `admission_binding`. The isolated M03 removes precisely the requested-target versus bound-generation join and falsely accepts. No generation2 principal, root or unsupported generation writer is invented.

5. **Exact custody survives equal totals.** The genuine auxiliary case has A5 result `(Mutation=1, Total=1)` for B1 and A13 result `(Mutation=2, Total=1)` for separately granted B6. The hostile A16 reply still names B1 but carries A13's genuine result lineage. Retained receipts remain unchanged. `binding.go:90` rejects the exact service mismatch at `reply_result_custody`; the attack test requires ordinary raw authority/reduction to pass first. A legitimate separate q2 grant is not itself classified as an attack.

6. **Real reduction is checked independently.** `custody.go:56` now compares every receive coordinate exactly; non-receive selection attacks cover emit/set/branch/end. The supplemental isolated-runtime probe extracts actual prior projection, program and receipt from the genuine/attack exports and calls runtime Reduce separately from the oracle: genuine receipt selection succeeds, empty Operation returns `ErrIneligible`. Its source, input hashes, isolated-copy provenance and log are retained. The oracle imports no runtime authorization/executor/Reduce/Replay implementation. This runtime comparison does not establish a positive service history through an internal API; it is a separate reduction-semantic control.

7. **All intended compiled mutants are genuine semantic kills.** Independently verified all 2,802 content-addressed mutant file mappings and 562 distinct blob hashes, each recorded source map and log hash, and that every mutant changes only its declared source file versus the frozen research map. All five builds exit 0. M01, M02 acceptance direction, M03, M04 and the restored wildcard each fail their intended assertion with `valid=true`, no panic. M02 also falsely refuses genuine P10 at `first_authority_policy`, while its sound parent succeeds. These are not compilation, parser, unrelated-guard or missing-fixture failures.

8. **Lifecycle and owner continuity preserve history.** Ordinary lifecycle checking uses the ordered historical subject/projection before reconciling current retirement and closure evidence. Inspected N10c/N10d raw facts: original authorization/grant facts precede actual CancelRun/RetireSubject results, with later service boundary refusal. `external.go:304` enforces declared capture order; lines 328–350 preserve immutable causal rows and monotone observation consumption. C07's actual sequence is epochs 1→2 (initial Open), 2→3 (first recovery), 3→4 (second recovery). A's cut sequences are 2 and 6. Both distinct process-loss records have matching PID, `signal: killed`, post-commit acknowledgment and token. The three renewal comparisons are separate; the test explicitly covers missing initial pre-Open evidence as unchecked and missing second kill acknowledgment as failure.

9. **Compatibility is evidenced, with ancillary scope explicit.** K01's actual authored A artifact is `alan/1e/v1` in format13; K02's is `alan/1e/v1` in format15. Both retain exact adapter service A3 after original selection A2, result consumption and one endpoint ingress/commit. Runtime normal/race logs contain actual passes for unchanged `TestPhase3B3CartesianOperations/r1-add-one`, authored conflict profiles, old-reader controls, duplicate transport/index checks and malformed-Open index/metadata/PRAGMA/profile/observation/admission variants. Independently compared each promoted14/16 raw parent byte-for-byte with its public K01/K02 parent capture, and every preexisting nonmetadata table in the post-Open/promoted pair. Only format metadata and appended receipt-conflict evidence change. This proves no public conflict-production schedule and no second application; the report states that limitation.

10. **Source and evidence provenance is concrete.** Independently hashed all 562 research and 350 runtime source archive entries against freeze manifests. All six independent campaign command records have before=after maps equal to the frozen research map, matching log hashes and exit0. Independently checked all 3,087 per-case digest entries in each final corpus (including complete archived database companions as bytes, never SQLite opened). Recounted raw case manifests and endpoint JSONL records: each mode has 81 cases, 80 Stores, 163 captures, 11 acknowledged losses, 89 ingress and 89 committed events. These agree with the claimed bounded counts.

11. **Old exports and costs are not inflated into proofs.** Independently compared both new outputs for all four old model/audit exports with the actual preserved files, decoding the two gzip files: byte-equal in every case. Recomputed all min/median/max timing distributions from the 40 cost observations, five per n=0/1/8/32 per mode. Allocations remain explicitly unmeasured, runner/boundary overhead is identified, and no total constant-time or speedup conclusion is claimed. Global research Staticcheck remains inherited non-green: compared the actual archived baseline to final logs byte-for-byte, 13 ordinary and 31 integration lines. The changed checker package is independently recorded clean; inherited findings are not silently waived as a clean repository.

## Verification actually performed by this reviewer

- Read-only diff/source/contract/evidence inspection using `sed`, `rg`, `nl`, and Python `json`, `tarfile`, `hashlib`, `gzip`, `statistics`.
- `/tmp/alan-phase3b4-oracle-evidence-check.py` completed with exit0; log at `/tmp/alan-phase3b4-oracle-evidence-check.log` records archive/digest/mutant verification and decisive raw coordinates.
- `/tmp/alan-phase3b4-oracle-secondary-check.py` completed with exit0; log at `/tmp/alan-phase3b4-oracle-secondary-check.log` records exact freeze-map agreement, mutant changed-file isolation, old export records, promoted parents, M03, lifecycle and equal-total custody checks.
- Additional inline Python independently verified old export byte equality, exact inherited lint logs, raw case/endpoint counts, and cost summary arithmetic. One reviewer counting attempt initially assumed `Captures` was always a list; N10h's actual null no-Store value produced a TypeError. The counting script was corrected to treat absent captures as empty and completed successfully. This was reviewer tooling, not a source/test failure or negative policy evidence.
- Existing actual test executions are attributed to the archived controller/implementer command records and logs; this reviewer did not rerun broad suites.

## Cannot verify / limits

- This review does not prove universal semantic refinement, complete adversarial coverage, arbitrary contexts, N-agent/distributed ownership, unlimited-loss liveness, compiler correctness, or authenticity against coherent unsigned full-history rewrites.
- Admission has no recorded global timestamp; the supported order is the documented causal linearization under genuine writer/selection assumptions. Endpoint records are contextual evidence; unique contact-to-gate attribution is not generally recoverable from the wire alone.
- The independent raw checker deliberately covers the retained finite instruction/artifact slice. K01/K02 substantiate authored compatibility in that slice, not all author/collection constructs.
- Formats14/16 promotion remains ancillary internal-writer compatibility. No public conflict-production schedule is claimed. Unsupported generation/revision-changing public construction remains an obstruction, not silently filled by SQL or host-injected consumption.
- Task 6 manuscript/LaTeX/PDF reconciliation is a separate outstanding task. Task 5 does not claim that the entire paper already reflects this new qualification.

## Code-quality verdict

APPROVE for the declared research slice: explicit version selection, independent transition checks, narrow prerequisite diagnostics, clear separation of raw/custody/physical/model reports, complete intended mutant semantics, and retained failed apparatus attempts. No new dependency, runtime production mutation, unrelated refactor or old-model export change is present in the Task 5 diff. No actionable source-quality defect was found. Both final quality and specification verdicts are APPROVE after the final gate/delta checks below.


## Final closure and evidence-only delta review

The provisional report is preserved separately at `/tmp/alan-phase3b4-oracle-review-r1-provisional.md`, SHA256 `d7268d842cd824bfc3ce79a2a170ebc1f17a3ab2ec637e7337904759516ba097`. This final revision does not erase the earlier hold.

- Inspected final research delta `0c303c3d..b940f0d2` and runtime delta `fc60391f..7c6d092a`, plus the completed report, qualification, runbook and gate summaries. Research delta is exactly 150 documentation/evidence files beneath `docs/foundational/phase3b4/`; runtime delta is only `docs/qualification/phase3b4/task5.md` and `task5-summary.json`. Both Git worktrees were clean at the verified final heads. All 562/350 current executable/input hashes still match the pre-run freeze; no source requalification is needed for these evidence-only additions.
- Independently validated all 147 `artifacts.json` entries and exact complete coverage of the final evidence tree. Independently validated every one of 28 controller command records, log hashes, expected exit, no-watchdog status, before/after maps against the applicable frozen source map, and absence of JSON test-failure events. The two nonzero exits are precisely the disclosed inherited research Staticcheck baselines, already checked byte-for-byte; they are not recategorized as clean lint.
- Full ordinary research command is `go test -json -count=1 -timeout=60m ./...`, with the explicit historical48-case input, new binding root, promoted root and progress audit environment recorded. Exit0, elapsed1517.3825709819794 seconds, 48 package-pass events, 2528 test-pass events; source maps unchanged. Its log SHA256 is `4586244ab2921822fcddee2e8c08a0bac0aaabe026d2c8bfca6e2395ec4629f6`. The three named ordinary skips are the opt-in `TestCommunicationArtifactExport`, `TestJoinArtifactExport` and `TestCandidateAExport`; final docs report them and distinguish the four required old-model outputs actually reproduced twice. Initial promoted-parent skips are covered by the later explicit normal/race commands; the historical baseline expected-failure control remains an explicit skip, never claimed as a new result. No independent binding, affected research race or prior-regression skip is hidden.
- The runtime companion JSON is exactly equal to the verified research qualified-pair object, verdict counts and all34 gate summaries (six independent campaign plus28 controller gates). Its source-free evidence nature and bounded claim wording are accurate.
- Actual chronology is now explicit in the final report: executable heads were committed clean before qualification; freeze manifests/archive files were generated and independently checked before qualification, then first committed with final evidence. The Task5 brief's literal wording also asked to commit the manifest files before qualification. The controller clarified that the governing user packet requires the exact source freeze and contemporaneous pre/post hashes, not an earlier Git commit of those manifest files. I accept that disclosed plan-wording deviation for this review because the final commands independently embed the same frozen source maps and the evidence is not claimed to have an earlier manifest commit. This does not establish authenticity against fully coherent unsigned reconstruction, which remains outside scope.
- The documented raw-log blank line at EOF is preserved byte-for-byte; normalizing original test output to silence a whitespace check would damage provenance. Final documentation contains no premature independent approval and leaves manuscript reconciliation to Task6.

Reviewer final-delta command `/tmp/alan-phase3b4-oracle-final-delta-check.py` completed exit0; its log `/tmp/alan-phase3b4-oracle-final-delta-check.log` records exact scope, artifact coverage, current freeze agreement, ordinary completion and28 command checks. An additional read-only Python comparison validated the runtime companion's exact34-gate correspondence. No unchanged broad suite was rerun.
