# Phase 3C.2 qualification and provenance

## Source boundary and attribution

Merged base: `69534188f796efab2e866750f6a76e2c7a3a1706`; Task 1 model/export:
`60d46a43dd32caa07f0b814e6741213f506a9cb6`; Task 2 independent audit and Task 3
input head: `a24f8ff1c6aac22a0cffcffc2c614ccfca819c71`. New Go inputs are confined
to 12 files in capacitycomposition and its exporter. Task 3 changes results,
qualification, README reproduction, roadmap and manuscript artifacts only.
No old source/dependency/evidence bytes, runtime behavior or SQLite stores were
changed; no SQLite or environment configuration was opened.

Task 1 executed scoped model tests/race/Staticcheck and exports; Task 2 executed
independent audits, checker tests, hostile controls and a fresh export. Task 3
reused their reviewed unchanged-source captures. The controller executed root
build/module/package checks in this task. Task 3 freshly verified their retained
bytes and source identity, built and inspected manuscript artifacts and checked
alignment/links. This attribution does not imply a fresh whole-root test run.

## Historical package qualification and exact input identity

[Input identity](evidence/qualification/input-identity.json) freshly compares
all 637 entries in Phase 3C.1's retained Go/dependency manifest against both
current bytes and merged-base Git blobs: all match. A broader tracked-file
check found 20 additional preexisting Go/module files outside that manifest.
They also match merged-base blobs and the original historical baseline commit
1456a0000a0e14bb3dd105026189717c71ab905b (20/20 present and identical). Twelve are under examples/evidence: ten
root-package Go inputs and two testdata consumer files. Three are an archived
Phase 1K public-consumer nested module (go.mod, go.sum, public_test.go); five
are archived Phase 3B.4 apparatus/fixture files. Their exact paths and hashes
are retained separately. They were not part of the original 637-file freeze;
this supplemental Git identity check does not retroactively place them there.
No unexplained input drift was found.

The inherited [Phase 3C.1 qualification](../phase3c1/05-qualification.md) combines
a fixed 70-package ordinary run (57 pass, 13 no-test skips, exit 0; 1,517.651s)
with its two new model/exporter package runs. It was not one 72-package command.
The clean current [roster comparison](evidence/qualification/package-comparison.json)
is exactly 74 = 72 inherited + capacitycomposition and its exporter, no removals.
The two added packages pass the new Task 1 ordinary/race checks. Combined
ordinary package coverage is supported by unchanged inputs and fresh scoped
checks; there was no fresh single 74-package test invocation, whole-root race
suite, integration run or runtime campaign. Nested modules are not implied to
be part of root ./... test coverage.

## Executed checks, reused without repetition

| Scope / exact command | Result | Attribution / archive |
| --- | --- | --- |
| `go test -v -timeout 30m ./research/capacitycomposition/...` | 17 tests across two packages pass; exit 0 | Task 1 `evidence/task1/test.log` |
| `go test -v -race -timeout 30m ./research/capacitycomposition/...` | 17 tests pass; no race report; exit 0 | Task 1 `evidence/task1/race.log` |
| `env STATICCHECK_CACHE=/tmp/alan-phase3c2/staticcheck-cache staticcheck ./research/capacitycomposition/...` | Exit 0, no diagnostics | Task 1 `evidence/task1/staticcheck.log`; initial cache failure retained separately |
| `go run ./research/capacitycomposition/cmd/export` | Two Task 1 outputs and Task 2 fresh output byte-identical | Task 1 export.log and Task 2 fresh-export.command.json |
| `python3 -m unittest discover -s docs/foundational/phase3c2/tools -p 'test_*.py' -v` | 5 methods pass, with 43 hostile controls and 144 generated-history fixtures | Task 2 unittest command/stdout/stderr |
| Python audit, plain JSON and gzip | 72 complete graphs and 10 captured histories accepted; all 43 hostile CLI inputs reject as intended | Task 2 manifest and captured command records |
| `go build ./...` | Exit 0; approved writable-cache rerun clean | Controller `evidence/qualification/build-cache-writable/` |
| `go mod verify` | Exit 0; all modules verified | Controller `evidence/qualification/module-verify/` |
| `go mod tidy -diff` | Exit 0, no diff | Controller `evidence/qualification/module-tidy-diff/` |
| `go list ./...` | Exit 0; clean 74-package roster | Controller `evidence/qualification/package-list-cache-writable/` |

Exact argv, working directories, source heads, environment overrides, elapsed
times, exit statuses and complete output are retained in the named command
records. Go uses GOCACHE=/tmp/alan-phase3c2/gocache. The controller's initial
build and package-list both exited 0 but emitted a module metadata cache
permission warning. Those initial logs remain beside approved clean reruns;
the warning is not a package. The copied runner records how those checks ran.

Global lint is not claimed clean or freshly rerun: Phase 3C.1 retained 13 ordinary
and 31 integration-tag inherited diagnostics, exactly matching its older
comparison. New scoped Staticcheck is clean. No go vet was invoked. Task 1 RED
and cache failure, Task 2 failed control attempts, and the distinction between
Python generated fixtures and independent producer-history replay remain intact.
The README race command now includes the recorded -timeout 30m. README is
outside the 32-file export source set; that documentation correction required
command comparison, not another long race run.

[Archive manifest](evidence/qualification/archive-manifest.json) binds retained
controller checks and verification files; [artifact verification](evidence/manuscript/verification.json)
binds the 201 Task 2 artifacts and all 32 producer inputs to current bytes.
Tool versions in the retained evidence are Go 1.27.4 darwin/arm64, module Go
1.27.3, Python 3.9.6, Staticcheck 2026.2.2 (0.8.2), Tectonic 0.17.0 and Poppler
as recorded in the manuscript command archive.

## Manuscript and PDF

Both manuscript forms add the same Phase 3C.2 section, scope, counts and
limitations; Abstract and Conclusion summarize the current result. Earlier
Phase 3C.1 claims remain in their original disjoint-resource domain. The new
Coffman/Elphick/Shoshani background entry matches the controller-inspected
original article's title, authors, volume/issue, date and pages. Its publisher
DOI page returned 403; no DOI verification is claimed. See the
[primary literature note](evidence/manuscript/primary-literature.md).

The PDF operation marker ran exactly once immediately before manuscript
source authoring, with edit/count 1/pdf. Exact command, baseline/current builds,
warning comparison, text extraction, rendered-page inspection and output hashes
are retained in [the manuscript manifest](evidence/manuscript/manifest.json).
The rebuilt baseline has 46 pages and the current PDF 48. Both logs retain
the identical four underfull hboxes, one 5.2184pt overfull hbox and two math-mode
font warnings; none is new. Same-index text changes occur on pages 1, 2, 4, 5,
11, 14, 22, 24, 32, 34 and 44-48. Visual inspection covers those pages and
continuations: 1-6, 11-12, 14-15, 22-25, 32-35 and 43-48 (24 pages total).
No new visual defects were found. The tracked output is paper/alan-dissertation.pdf. Reproduction uses existing
Tectonic/Poppler, with no added dependencies. All changed mathematical content
is aligned after markup normalization; inherited Markdown-only Section 7,
different rendered section/theorem numbering and format-specific markup remain.
No global manuscript-equivalence claim is made.

The retained verification command is:

```sh
python3 docs/foundational/phase3c2/evidence/manuscript/verify.py
```

## Acceptance boundary

The source argument, finite abstract checks, factored probes and selected real
histories support different scopes. Source hashes and reviewer approval do not
replace mathematical proof or a runtime qualification. See [results](04-results.md)
for explicit open obligations and [review status](06-reviews.md) for separate
independent-review and delivery gates. General Phase 3C remains open.
