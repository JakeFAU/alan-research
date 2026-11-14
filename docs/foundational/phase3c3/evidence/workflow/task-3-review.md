# Task 3 review

## Spec compliance — REJECT

- ❌ One important mathematical qualification issue remains at `paper/alan-dissertation.md:1160` and `paper/alan-dissertation.tex:811`: the snapshot counterexample is described as establishing necessity of the forgotten transition relation. It establishes insufficiency of this snapshot, not necessity of retaining that entire relation. The rest of the requested task is present and appropriately bounded.
- Reviewed base `c4b4e503a1f7e135ea25e06fc918c008e33f0100` to head `32ffd396b3aa7767dc5908714a36e18cb9266c9d`, using the supplied diff, task brief, global constraints, frozen contract and implementation report. The long diff was processed in focused passes after output truncation; large generated manifests were parsed rather than reproduced line by line. No changed source file was separately reread to duplicate its diff.
- All requested file categories are present: results, qualification, archived controller checks, deterministic verifier, roadmap, both manuscript forms, bibliography and rebuilt PDF. The 46-file package contains no production/model/dependency changes; controller review/delivery records are outside this review.

## Strengths

- `docs/foundational/phase3c3/04-results.md:4-10,28-40,43-66` and `paper/alan-dissertation.md:1156-1174` distinguish snapshot failure, conditional initial existential completion, concrete-model correspondence and the doomed-residual boundary. The language-intersection argument covers projection and reconstruction, independent finite private segments, empty words, nondeterminism, common visible synchronization and conjunctive marking.
- `paper/alan-dissertation.md:1158,1168` and `paper/alan-dissertation.tex:809,819` explicitly preserve all marks, including the doomed sink; distinguish a monitor from complete-agent substitution; retain all eight original visible actions/bindings; and explain why stable-ready's unrestricted progress argument cannot be imported after continuation restrictions.
- `docs/foundational/phase3c3/04-results.md:70-99` keeps configuration sums, language comparisons, triple products, selected histories and hostile controls distinct. Its totals agree with the accepted independent-counts record: 1,980 products, 19,544 states, 23,420 arcs, 9,526 completable states, 622 initially completable and 392 nonblocking products; 1,188 substitution rows, including 396 successful and 792 blocked; eight histories and 123 events.
- `docs/foundational/phase3c3/05-qualification.md:19-36,40-70` correctly attributes source-identical historical ordinary coverage and fresh scoped checks. The prior Phase 3C.2 qualification supports the 637+20+12 input history and 70-package ordinary-run attribution. No fresh combined full suite, clean global lint or runtime campaign is invented.
- `docs/foundational/phase3c3/evidence/manuscript/verify.py:32-107,109-153` supplies useful artifact/source checks and explicit changed-content alignment without pretending to rerun the independent semantic audit or mechanize the argument. `05-qualification.md:76-96` preserves inherited manuscript differences and existing warnings.
- `paper/alan-dissertation.md:1174`, `paper/alan-dissertation.tex:825` and `04-results.md:103-110` explicitly exclude runtime refinement, complete-agent replacement, full abstraction, optimality and verification-cost reduction. General Phase 3C remains open.

## Issues

### Critical

None found.

### Important — P2: Narrow the snapshot necessity claim

`paper/alan-dissertation.md:1160`; `paper/alan-dissertation.tex:811`; rendered PDF page 46:

> The forgotten future transition relation is therefore semantically necessary for this snapshot candidate.

The matched-cut witness demonstrates that the candidate snapshot collapses two completion-relevant futures. It does not show that a literal transition relation is necessary information for an interface. In fact, the following positive result deliberately identifies different transition structures through their marked visible languages, and the role-restricted grammar does not establish a general necessity theorem. This sentence introduces an unsupported necessity claim at the point where the negative result is interpreted.

Replace it in both formats with, for example: **“The snapshot therefore omits a completion-relevant distinction in future ordering.”** Refresh the affected PDF/extraction/render/source hashes and alignment evidence. Do not change the frozen contract, model, captures, tests or argument to accommodate this wording fix. No Go suite or independent semantic audit rerun is needed.

### Minor — inherited diagnostics remain

`docs/foundational/phase3c3/evidence/manuscript/comparison.json:19-38` retains four underfull boxes, one 5.2184pt overfull box and two math-mode font warnings, exactly matching baseline. `05-qualification.md:34-36` preserves earlier global Staticcheck diagnostics. These are disclosed inherited noise, not new regressions or blockers for this task. No unrelated global cleanup is requested.

## Focused checks actually performed

- **Archive/provenance risk:** a read-only Python check compared all 669 unique baseline current-file hashes and all 14 new Go current-file hashes with their manifests. It matched the archived controller job/result bytes and decompressed outputs against the original `/tmp/alan-phase3c3/qualification/` captures for build, module verification, tidy-diff and package list. All matched. This did not rerun any of those commands.
- **Artifact substitution/layout risk:** checked all current paper hashes, all 22 manuscript artifact/supporting-file records, and all 16 available render hashes against the manuscript manifest. All matched. Parsed the supplied initial/final build command records and all 16 rendering records; their recorded exits are zero. Reviewed the verifier's recorded PASS and its source rather than rerunning it or its Git queries.
- **Concrete layout risks:** applied the PDF skill read-only, without running the operation marker. Viewed existing hash-matched renders of pages 1, 2 and 45–48: the revised Abstract fits page 1; contents starts cleanly; the new subsection, equation, cross-page snapshot paragraph, evidence paragraph, Conclusion and bibliography transition have intact margins and readable continuations. This reviewer did not independently view the other ten pages reported by the implementer, and did not rebuild or rerender the unchanged PDF.
- **Historical-coverage inflation risk:** checked the owning unchanged `docs/foundational/phase3c2/05-qualification.md`, which confirms the inherited test boundaries, package progression, input-manifest qualification and lint limitations cited here.
- **Literature-attribution risk:** checked `/tmp/alan-phase3c3/task3-handoff.md` against `evidence/manuscript/primary-literature.md:3-11` and the new bibliography diff. They consistently attribute Tang/Moor publisher inspection to the controller and disclose failed De Nicola/Hennessy access. No independent primary-source access is claimed by this reviewer.
- No Go tests, race suite, lint, model audit, exporter, build, PDF operation marker or Git commands were run. No source/index/HEAD changes were made; the only review output is this requested report.

## Cannot verify from this task diff

- The previously accepted producer and independent checker implementations are unchanged and were not re-reviewed in this task. Their universal correctness cannot be inferred from manifests or counts. This task correctly keeps those finite checks separate from its source argument.
- The 669-file merged-Git comparison is evidenced by the inspected verifier and recorded PASS; this reviewer freshly checked current hashes, not the historical Git blobs. The controller should retain that recorded execution evidence.
- Exactly-once marker execution and human visual inspection of all 16 pages are historical operation claims, not facts recoverable solely from one successful command record. The supplied controller handoff and worker report are consistent; this reviewer did not repeat the marker or all visual inspection.
- Tang/Moor primary-text inspection remains controller-attributed. This reviewer checked local metadata/attribution consistency, not the publisher itself.

## Assessment

**Task quality: Needs fixes. Explicit verdict: REJECT.**

The implementation is narrow, well-qualified and supported by coherent source-bound artifacts. One unsupported necessity sentence must be weakened and its derived manuscript artifacts refreshed before approval; no broader implementation or qualification defect was found.
