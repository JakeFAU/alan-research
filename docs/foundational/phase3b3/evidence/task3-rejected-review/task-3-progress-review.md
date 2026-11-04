# Task 3 independent projection, corpus and progress review

## Verdicts

**Specification: APPROVE for the assigned Task 3 slice.** The frozen candidate satisfies the specified bounded corpus/projection/progress requirements. No blocking specification defect was found in the inspected source or retained evidence. This is not approval of the complete branch, manuscript reconciliation, final committed-source qualification, or arbitrary runtime correspondence.

**Quality: APPROVE for the assigned slice.** The generator uses actual unchanged model successors, the raw checker is upstream of projection, exhaustive progress enumeration retains losing states, and the independent ledger check is applied to the actual constructed completion witness. No actionable P0–P2 finding was established. The ignored `AuditTrace` error discussed below does not create a successful false suffix audit in this frozen implementation.

## Reviewed identities and boundaries

Read `task-3-progress-review-brief.md` first, then `task-3-brief.md`, `task-3-report.md`, the relevant file sections of `task-3-review.diff`, and the actual corresponding source. The source reads focused on `corpus.go`, `project.go`, `progress.go`, their tests/commands/README and `files.go`'s projection/CLI interface. Read `mutations.go` for mutation ancestry, and the relevant `check.go`/`custody.go` interfaces for derivation of original identity, serial service, selection mismatch and causal availability. The other independent reviewer owns complete raw policy/loader/external checking.

The research base is `31c8af68e818f97a842e790669fbc8f8c98dc8bd`; raw runtime source is `759df6e3d2505f27e1287b2cfeb07bc4290afc0a`. I independently computed SHA256 for all 22 current files against `/tmp/alan-phase3b3-task3-evidence/final/source.json`: no mismatch. The review manifest `task-3-review-source.json` lists the same source identities. I independently checked the byte count and SHA256 of all 153 entries in `final/manifest.json`: no mismatch. Thus the inspected source and inspected candidate-run evidence agree on their advertised identities; HEAD alone is not used to identify uncommitted checker files.

Read frozen contract v3, the interface report, progress-design/review.md and selection-narrowing/review.md for requirements and previously identified risks. These earlier audits were not substituted for checking final candidate output.

Read unchanged sources outside the candidate only for named risks:

- `research/authorityservice/oracle.go:48–71`: whether an ignored returned error could leave a successful audit flag set.
- `research/authorityservice/oracle_corpus.go:38–148`: whether selected mutation parents actually describe the source histories, rather than unrelated generic traces.
- SHA256 reads of `authorityservice/{model,explore,oracle,oracle_corpus}.go`: identities match the frozen worker report. No old model was edited or redefined.

No runtime API, Replay, SQLite reader, archived SQLite database, graph suite or whole-module test was invoked. No subagent was spawned. The only repository write in this review is this report. Inspection commands were `cat`, `sed`, `rg`, `nl`, and inline Python JSON/hash comparisons. No new execution probe was needed after source and retained evidence answered the specified questions.

## Specification assessment

### Positive correspondence and actual raw originals

The actual normal S2 input is `/tmp/alan-phase3b3-runtime-frozen/normal/ordinary-s2-handshake-628745432/final/inspect.json`, SHA256 `e371654a8b8b9def3ab6f29edfd067c7bc418f7c9e58a72da23fed2ff8296d18`. I read its original permission tuples, history selections, attempts and ordered facts, and compared them with `final/normal-expected.json`.

AB `["a",1,1,0]` and BA `["b",1,1,0]` commit at orders 1 and 3. B's accepted grant `["b",1,2,0]` is retained at order 7. A selects ready at local sequence 2 and the authentic BA request at sequence 3. Its BA service `["a",1,4,0]` reserves delegated authority at order 16, receives its gate/result at 17/18, and emits BA reply `["a",1,7,0]` at order 21. B consumes that reply at order 23 before its AB service `["b",1,9,0]` at reservation order 28. A consumes the original AB reply at order 35. The raw permission rows genuinely contain two distinct B service endpoint/obligation tuples. Both originals remain outstanding when BA service starts.

The projection's 21 events preserve these facts and operation IDs. Its BA authorization is Delegated, AB is OwnRoot, and both final ledger Paid values are true. Model event descriptions come from unique matching `Successors`, not handwritten authority witnesses. The independent raw checker derives the permission verdict first; projection cannot supply authority to it. The many-to-one B whole-action abstraction is visibly declared and does not establish permission-algebra equality or shared-resource interference behavior.

Normal S1 input `/tmp/alan-phase3b3-runtime-frozen/normal/ordinary-s1-own-947119016/final/inspect.json`, SHA256 `ccd6367dc6f3c5d28db9d91af787e451e249bee038e888fad85a3870525034c3`, has the actual extra A own permission and no retained delegation operation. Its 20-event OwnRoot projection and both consumed originals are supported by those bytes. This uses the corrected final S1 fixture, not the older unused-grant design note.

### Refusals, mismatch and unsupported domains

The final delayed-selection input is `/tmp/alan-phase3b3-runtime-frozen/normal/control-wrong-prior-selection-3204885184/final/inspect.json`, SHA256 `79bfd70eaf8b13f217a85a4b4d5bf0d10b86e0959dfd408dfc5c53418ec53acb`.

Its BA selection occurs at A sequence 2/order 5; the accepted exact grant is order 8; ready selection replaces the immediate predecessor at A sequence 3/order 13; A4 service reserves at order 16 with `authority_request_not_admitted`. The retained selector is Delegated and names the genuine grant/request/service. The output preserves raw validity, both unpaid originals, and an eight-event model witness ending in Allowed/Delegated with `immediate_selection_narrows_model_admission`. It does not confuse the valid runtime refusal with missing model consent. `project.go` stops at this first disagreement and preserves the model audit separately.

The actual first-refused-then-grant snapshot has SHA256 `8f8b558c079c3adb2e04e84f81d13973cf6c6630fd102fd395daf755a6c3419a`: service refusal at order 9 precedes genuine grant retention at 15. Its eight-event projection keeps the first Refused decision and both originals unpaid. It does not retrospectively authorize the earlier service.

The actual serial/serial final snapshot has SHA256 `710c75cd2b40e42b660af37cef6fd970ae0923f984fc222508aeee0ea733e166`: both originals are committed/admitted, neither selected, and the accepted grant remains. A seven-event safe unpaid service prefix is a supported projection, not a runtime progress claim. Original program histories and supplied Run outcomes remain separate.

I inspected the final outputs for all 48 normal and 48 race cases. Each has 117 captured phases, 43 valid final raw histories and five hostile rejections. Neither output has a projection apparatus Error. Each retains exactly the named selection-narrowing mismatch. Separate-grant revocation, selector refusals, reversed service, unknown result prefixes and repeated physical exposures are explicitly unsupported. Shared-process cut outputs are explicitly marked unsupported for chronological process recovery even when a 21-event final service subsequence audits. No paired crash expansion is asserted and no losing state is removed to claim progress.

The five invalid histories fail at the substantive expected boundaries: deleted grant `lifecycle_settlement`, forged grant `lifecycle_origin`, altered effect selector `reservation`, changed gate attempt `gate_attempt_digest`, and substituted receipt `result_original`. These are raw rejections, not merely unsupported model projections.

### Corpus coverage and mutation provenance

`final/model-corpus.json` contains all eight required classes and 22 selected/generated cases: four complete value pairs; one incomplete prefix with its explicit suffix; one blocked case; four missing-authority, three scope, six request/generation, one revocation and two causal-order controls. `Prior` retains all 33 histories and four ancestry controls.

`corpus.go` commits/admit both originals before retaining the late grant and serving BA. The positive event at index 7 is explicitly checked as BA/Delegated. Complete histories and incomplete suffixes use actual forward transitions. Protocol-valid but consent-unsafe S3/Cartesian controls retain ProtocolValid=true and AuthoritySafe=false with their first unsafe event at index 3; they are not relabeled malformed protocols.

I compared the selected mutations' event arrays with their advertised parents in the retained JSON and checked the old constructor. The ordinary identity/scope controls change the actual `complete-values-1-1` history at their indicated event; late first-decision rewrite extends `late-grant-sticky-refusal`; revoked fresh admission changes the final refusal event of `revoked-before-fresh-admission`; the new missing-grant and swapped-gate controls derive from `reciprocal-1-1`. These parent identities are meaningful. Full raw mutated inputs retain input hashes, changed-cell references and first failed edges; the stronger actual runtime-result-as-grant control remains separate from the synthetic field mutation.

### Exhaustive progress and audited suffixes

`progress.go` explores every reachable `State` with unchanged `Initial/Successors` and retains all arcs, including loops. It does not apply the proposed predicate during forward enumeration. Safe reverse BFS is seeded from safe Complete states and traverses safe predecessors; the predicate is compared with that resulting reachability set only afterward. Its stored next edge decreases reverse-BFS distance, so the reconstructed suffix terminates at a complete state. BFS parent events supply the actual prefix; successor indices select actual model events.

The inspected `final/progress.json` gives:

| Query | States | Arcs | Safe-completable | Predicate mismatches | Max shortest suffix |
| --- | ---: | ---: | ---: | ---: | ---: |
| S2/LateGrant | 111,065 | 840,616 | 107,213 | 0 | 23 |
| S2/Revocable | 220,564 | 1,740,832 | 193,681 | 0 | 22 |

LateGrant excludes 1,926 first-refused states without a grant and 1,926 with one; unconditional completion's retained shortest counterexample has four events. Revocable retains 22,887 noncompletable states without BA first refusal; its shortest counterexample is the one-event initial revoke. Both graphs have zero incomplete outdegree-zero states, so no graph dead-end count is substituted for service completion.

The actual longest witnesses include two initial crash events, then suffixes of 23/22 events, respectively. Their independent audits check 25/24 total events and report ProtocolValid=true, AuthoritySafe=true, Paid=[true,true]. Counterexample audits check 4/1 events and remain valid/safe with Paid=[false,false]. These are finite existential graph facts; no fairness, forward invariance, runtime-state completion or resource-independent theorem follows. The Revocable query checks both originals' first decisions plus revoked/applied/permit status. It does not reinterpret or replace the prior 6.10 formula or old graph definitions.

## Quality assessment and specific error-handling question

`progress.go:116` discards the error from `AuditTrace`. In the frozen oracle (`oracle.go:48–71`), every error path sets both ProtocolValid and AuthoritySafe false; `progress.go:157–158` rejects either false flag for the longest suffix. Therefore the ignored error cannot turn a failed longest-suffix independent audit into successful QueryProgress output. The returned counterexample audit remains visible, and both actual counterexample audits are valid in retained evidence. Error text propagation would improve diagnostics, but I do not classify this as a current correctness blocker. The analogous ignored prefix errors in `corpus.go:93` occur only after the full audit succeeded, and this prefix-based ledger cannot have an invalid prefix with a valid full trace.

Projection does compare only the abstract allow/refuse bit directly at reservation time; concrete root/path/binding validation is supplied by the independently derived raw verdict and declared bounded mapping. The actual positive cases above use the intended paths. This is sufficient for this frozen corpus, not evidence of a general mapping of every valid profile13 history or permission algebra.

CLI expectations do not enter `CheckRaw` or model generation. I compared parsed normal/race expected versus unqualified JSON outputs: they are identical for each run. Thus changing expected exits does not rewrite the raw verdict or mismatch. Unsupported projection remains visibly unsupported and is not converted to an invalid raw history.

Parent corrections were checked at their source/evidence boundary: current subjects must exactly cover original principals (`check.go:280` plus existing unique identity validation); actual historical selections invoke `requireObservationPrefix` (`custody.go:50,87–118`). The final parent-control records preserve exit1 at current_subject_set and at the receipt/lifecycle/message/terminal causal boundaries (orders 18/7/11/37). The mutation inputs and verdicts are included in the verified evidence manifest. Full adversarial review of those mechanisms is also within the other reviewer's remit.

## Verification provenance and remaining obligations

`final/commands.json` records 14 candidate commands with the expected exits. The final normal/race causal-control test logs independently contain 41 PASS entries, zero SKIP and zero FAIL entries each. Existing Staticcheck logs and the 131-byte nonfatal build-cache stderr are preserved. These are inspected prior executions, not tests newly run by this reviewer.

This approval does not certify whole-module regressions, all sixteen final graph export/audit comparisons, immutable final committed-source execution, archive extraction, paper/LaTeX consistency, production runtime implementation or release status. Those remain parent-owned. It also does not establish unsigned-history authenticity, external contact completeness, general compiler semantics, arbitrary program substitution or whole-process crash correspondence. No merge or release is authorized by this review.
