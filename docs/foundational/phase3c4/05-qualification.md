# Phase 3C.4 qualification and provenance

## Source boundary

Merged baseline: `c369e25c8db268d418ab7c11573d46dbc1d9425f`, identical in tree to
delivered `fc3483aa3a0a1ce0d8ff7110dd225815d11d855b`. Frozen final Go source:
`9b762fb2b3e5959997f05b40089d17b05110ea8d`. Task 1 evidence correction:
`5b0496f9de78149df69937e052297af377a82dc2`. Accepted independent Task 2 and Task 3
base: `92c5d55a95dbcc08525deee5619f51fce47b7789`.

Task 3 changes results, qualification, roadmap, manuscript/PDF and its verification
artifacts. Frozen contract/argument/schema, Task 1/2 sources/evidence, Go, old
models, runtime, modules and old evidence remain unchanged. The controller owns
review and delivery records; independent acceptance does not mean merged.

## Fresh identity checks and historical test attribution

[Baseline inputs](evidence/qualification/baseline-inputs.json) bind all 683 tracked
Go/module/workspace inputs to actual merged Git blobs and current bytes.
[New Go inputs](evidence/qualification/new-go-inputs.json) bind all 12 added Go
files, including tests, to the frozen final Go commit.
[Embedded inputs](evidence/qualification/embedded-inputs.json) bind all five
tracked go:embed targets to merged blobs. Task 3 rechecks every identity; this
supplements source attribution without expanding historical test coverage.

[Package comparison](evidence/qualification/package-comparison.json) compares the
fresh 78-package roster with the **actual merged prior gzip blob**: 76 inherited
packages plus residualinterfaces and its exporter, no removals. As documented in
[Phase 3C.3 qualification](../phase3c3/05-qualification.md), the inherited older
70-package ordinary run had 57 passing packages and 13 no-test skips, followed
by scoped model-package runs. They were not a single 76- or 78-package test run.
No fresh whole-root ordinary/race/integration suite or global lint-clean status
is claimed. Prior global Staticcheck qualification retains 13 ordinary and 31
integration-tag inherited diagnostics. No go vet was invoked for this packet.

## Actually executed evidence

| Scope | Result and exact attribution |
| --- | --- |
| Task 1 ordinary and race tests before generic label correction | Six top-level tests passed in each mode; preserved in `evidence/task1/pre-strictness-correction/ordinary-execution.json` and `race-execution.json` |
| Task 1 final affected ordinary/race reruns | `go test [-race] -v -timeout 30m -run TestCompleteEvidenceAndWitnesses ./research/residualinterfaces/...`; one affected test passed in each mode; exporter has no tests; final `ordinary-execution.json` / `race-execution.json` |
| Task 1 final scoped Staticcheck | `staticcheck ./research/residualinterfaces/...`; exit 0, empty output; `staticcheck-execution.json` |
| Task 1 final exports and Task 2 fresh export | All three raw captures identical, 8,271,854 bytes, SHA-256 `b40f1d5ac0aea0f43999502ea318cb128f545eb8fb1fd129451915a5aa84bf1f` |
| Task 2 tests | Nine methods pass, including 39 semantic mutation subtests and seven malformed JSON controls; `evidence/task2/qualification-1/tests.json` |
| Task 2 exact CLI executions | 50 total: two valid acceptances, 48 hostile rejections; nine byte-level controls plus 39 semantic attacks; `cli-commands.json` |
| Controller `go build ./...` | Exit 0, empty output, 15.277s; `evidence/qualification/build/` |
| Controller `go mod verify` | Exit 0, `all modules verified`; `module-verify/` |
| Controller `go mod tidy -diff` | Exit 0, empty output; `module-tidy-diff/` |
| Controller `go list ./...` | Exit 0, 78 packages; `package-list/` |
| Task 3 verifier | Source/artifact hashes, baseline Git identities, actual counts, archive inputs, manuscript alignment and new links; `evidence/qualification/verification.json` |
| Task 3 cached PDF builds and extraction | Baseline/current Tectonic, pdfinfo and pdftotext exit 0; `evidence/manuscript/build-commands.json` |

Root checks were run by the controller on frozen Go head `9b762fb2`; Task 3
archives their exact job/result/output records, not a claimed new execution.
Historical and scoped Go suites and the independent audit were not repeated.

Original/pre-correction dependency enumerations exited 0 but contained nonfatal
cache warnings. The first Task 1 report incorrectly claimed clean output; the
independent review **REJECT** and initial report are preserved. Evidence-only
repair retained both originals, added a clean resolved enumeration with identical
stdout, and obtained scoped **APPROVE**. Final source and capture hashes were
unchanged. `evidence/task1/go-list-deps-resolved-verification.json` records 109
dependency packages, zero package errors, nine matching compiled local inputs and
40 checked source records. This is not a new final full scoped test invocation.

Task 1 compilation REDs and generic strictness-label failure are distinct from
semantic falsifications. Task 2 retains its initial valid-export rejection due
to constructor ordering and its invalid fixture path failure. No hostile input
was initially accepted; the full 57-entry Task 2 archive manifest remains checked.

## Manuscript, bibliography and PDF

The [primary-literature note](evidence/manuscript/primary-literature.md) records
the controller's publisher-text inspection of Tang and Moor, Discrete Event
Dynamic Systems 34, 125–161 (2024), DOI 10.1007/s10626-024-00394-2. Existing metadata
matches; no BibTeX edit is warranted. Priority shaping and stronger abstraction
rules are not instantiated. No additional source inspection is claimed here.

The PDF marker executed once, edit/count1/pdf, immediately before source authoring.
Fresh baseline/current cached builds produced 50/51 pages. Each has seven warning
lines: four underfull boxes, one 5.2184pt overfull box and two inherited math-font
warnings. Raw warning lines differ in source line numbers; normalized diagnostic
text matches. No new warning category or magnitude appears. All 44 changed pages
(1, 2, 10–51) and adjacent pages 3 and 9 were rendered at 100 dpi and visually
inspected in full-resolution contact sheets. The new section is on pages 47–48;
the revised Conclusion continues through page 50. No new clipping, overlap or
illegibility was observed. This is scoped visual review, not global proofreading.

Both new sections, changed Abstract and changed Conclusion are mechanically
aligned after removing only markup differences. The user's corrections retain
version-specific instruction admission and the open general congruence claim.
TeX references `thm:congruence` (rendered Theorem 4.5); Markdown retains historical
Theorem 4.8 and points to Section 4.5. The inherited Markdown-only Section 7,
condensed TeX sections, numbering and bibliography presentation differences are
not normalized or claimed fully reconciled.

The first contact-sheet attempt with system Python failed because Pillow was
absent; the bundled Python succeeded. No PDF build failed. Self-review removed
accidental inline math delimiters inside new display equations before building.
The first verifier run omitted initial-language decision nodes from its recount;
the corrected run then rejected its own stale source hash, as intended by the
integrity check. Both failed invocations and exact verifier sources are retained
under `evidence/qualification/verifier-attempt-{1,2}/`; the count logic and final
source manifest were corrected without changing producer/checker evidence.
Build/extraction/render commands, PDF/text/log hashes, visual coverage and the
contact-sheet preparation failure are retained in the [manuscript manifest](evidence/manuscript/manifest.json).

## Checking the archive and recording a new attempt

From the repository root, the deterministic verifier and standalone audit CLI
can check existing evidence without writing to the archive:

```sh
python3 docs/foundational/phase3c4/tools/verify.py
python3 docs/foundational/phase3c4/tools/audit.py docs/foundational/phase3c4/evidence/task2/qualification-1/export-third.json.gz
```

The verifier reads evidence and sources; it neither imports runtime modules nor
opens SQLite. The [final review clarification](evidence/final-review-fix/) retains
the reviewed qualification manifest and a fresh verifier invocation. Earlier
qualification runs retain their original source and execution attribution.

The [Task 1 archive README](evidence/task1/README.md) identifies `qualify.py` as
the historical capture driver. Do not rerun it in place: it still writes to its
own archive and can replace existing evidence, even if the run fails partway through.

For a new checker test, export, and CLI qualification, use the existing capture
script from the repository root:

```sh
mkdir -p /tmp/alan-phase3c4/task2
python3 docs/foundational/phase3c4/tools/capture.py NEW_ATTEMPT_NAME
```

Replace `NEW_ATTEMPT_NAME` with a single new directory name, unused under both
`evidence/task2/` and `/tmp/alan-phase3c4/task2/`. The script creates the archive
directory without allowing reuse before it writes any results. Its compiler
roster check uses the absolute paths in the saved Task 1 enumeration, so run it
at the original checkout path for that comparison to match.

The capture command does not run the Go ordinary or race suites or Staticcheck.
For a new full scoped run of those checks, run them separately:

```sh
go test -v -timeout 30m ./research/residualinterfaces/...
go test -race -v -timeout 30m ./research/residualinterfaces/...
staticcheck ./research/residualinterfaces/...
```

The historical final ordinary and race reruns used the narrower test selection
shown above. Rebuilding and rendering the manuscript are separate and use the
archived scripts. No new concrete/runtime qualification, old exhaustive probe
execution, general congruence proof or delivery/integration claim follows.
