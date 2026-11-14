# Phase 3C.5 qualification and provenance

## Frozen identities and attribution

Merged baseline: `347d93ba25a5f7a10f75144b1a44107be09a2f39` (same tree as
`358fb222dec0e00148b3f2a804531c5b31202342`). Reviewed Go source:
`a036b4a6133f1285fd52dc52db2eb782a3ae8439`. Initial Task 2:
`b079c67db3298688f32f2c5c8218a01f81d45601`; repaired accepted Task 2 and Task 3
base: `208e1411dd9074f4d97eb9fbba4fb9b8755b3236`.

Task 3 changes results, qualification/verifier, roadmap, manuscript and PDF.
Contract, argument, schema, model/checker sources and old evidence stay immutable.
Controller owns independent review, workflow and delivery. Review approval is
not merge status.

[Baseline inputs](evidence/qualification/baseline-inputs.json) bind all 695
tracked Go/module/workspace files to actual merged Git blobs and current bytes.
[Embedded inputs](evidence/qualification/embedded-inputs.json) bind five targets.
[New inputs](evidence/qualification/new-go-inputs.json) bind 15 Go files including
tests to the reviewed producer commit. The verifier rechecks these identities,
the complete baseline roster and all accepted Task 1/2 tracked bytes.

[Package comparison](evidence/qualification/package-comparison.json) uses the
actual merged historical 78-package `output.log` and fresh 80-package output.
Only residualantichains and its exporter were added; none removed. The failed
preparation guess `output.log.gz`, then correction to the actual tracked name,
is retained. This is package enumeration, not test execution. The module-cache
seed copied existing pinned SQLite source/archive to a writable cache; no
module update or network fetch is claimed.

Historical root qualification remains attributed to its original sources and
commands. Earlier 70-package ordinary qualification had 57 passing packages and
13 no-test skips, followed by scoped added-model runs; it was not one fresh
80-package suite. Prior global Staticcheck's 13 ordinary and 31 integration-tag
diagnostics remain inherited. No fresh whole-root ordinary/race/integration
suite, global Staticcheck-clean result or go vet execution is claimed.

## Executed records, not repeated suites

| Evidence owner and scope | Actual result |
| --- | --- |
| Task 1 scoped ordinary/race | 12 passing top-level tests each; intentional false-hypothesis selector skipped in normal runs |
| Task 1 selected RED controls | Three deliberately false hypotheses fail as intended, retained separately |
| Task 1 scoped Staticcheck | Exit 0, empty output |
| Task 1 final captures 3/4 | Identical raw 14,097,260 bytes; SHA-256 `ab3cbdfda577a3ae804eb76a94a8c7ac665cf00c69f992ad6804b6e0d0e27b8a` |
| Task 2 qualification-1 fresh producer capture | Same raw identity; retained historical Python qualification does not acquire source-only guarantees |
| Task 2 accepted qualification-2 | 13 tests; nine cache-boundary controls; two valid CLI acceptances and 56 JSON/wire rejections |
| Controller `go build ./...` | Exit 0, 13.584 seconds, empty stdout/stderr |
| Controller `go mod verify` | Exit 0, all modules verified |
| Controller `go mod tidy -diff` | Exit 0, empty stdout/stderr |
| Controller `go list ./...` | Exit 0, 80 packages |
| Task 3 verifier | Deterministic identity/archive/alignment/link checks; no model or suite execution |
| Task 3 PDF | Fresh cached baseline/current builds, extraction and visual inspection; records below |

Controller root commands ran on reviewed Go head `a036b4a6`; their
[jobs, results and streams](evidence/qualification/build/job.json) are archived
unchanged, not rerun. The Python repair changed no Go inputs, so these checks
retain their source attribution. Producer manifests distinguish 54 semantic
read inputs from ten compiled local files. The recorded dependency enumeration
is checked using recorded module/root paths, not this checkout's absolute path.
Initial captures 1/2 retain dependency-cache warnings; only captures 3/4 are
final clean captures. See [Task 1 manifest](evidence/task1/artifact-manifest.json).

Accepted checker command: `python3 docs/foundational/phase3c5/tools/run_audit_tests.py`.
Bare unittest discovery remains historical. The
[corrected manifest](evidence/task2/qualification-2/checker-sources.json) binds
nine source buffers and all 59 child reports. The verifier resolves historical
qualification-1's three changed C5 files to
[pre-fix snapshots](evidence/task2/fix-1-red/manifest.json), also matched against
the initial Task 2 commit; unchanged predecessor/helper records still match
current bytes. It does not falsely require old hashes to match repaired files.
All old captures and the successful cache attack remain immutable. The separate
initial serializer failure remains in `evidence/task2/red-initial/`.

## Manuscript/PDF and representation boundary

The new section, Abstract and Conclusion are synchronized in Markdown and TeX.
The eleven-form instruction grammar, historical congruence correction and open
Theorem 4.8 remain. Inherited Markdown-only Section 7 and numbering/formatting
differences are not normalized or claimed resolved. The only new reference is
Ware and Malik (2011), verified against its primary arXiv metadata; the
[inspection record](evidence/manuscript/primary-literature.md) bounds the reading
and explicitly disclaims checking/importing its complete proof.

[Build records](evidence/manuscript/build-commands.json) and
[comparison](evidence/manuscript/comparison.json) retain actual commands, exits,
page/warning measurements and changed-page coverage. The first current build
failed because an inline T_F lacked math delimiters; exact source/log/command
records are preserved before correction. This authoring error did not affect
model evidence. Final PDF measurements and visual coverage are recorded in the
[manuscript manifest](evidence/manuscript/manifest.json), with hashes for source,
PDF, logs and extracted text. Seven inherited final-log warnings remain; no
new final warning is accepted. No global proofreading claim follows.

Run `python3 docs/foundational/phase3c5/tools/verify.py` for the reproducible
source/archive/alignment check. [Recorded result](evidence/qualification/verification.json)
is scoped to these checks. Mathematical and experimental boundaries are in
[results](04-results.md); no new runtime, SQLite, physical-history or distributed
qualification is asserted.
