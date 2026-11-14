# Phase 3C.1 qualification and provenance

## Source boundary and attribution

The investigation base is the verified Phase 3B closure merge
`1456a0000a0e14bb3dd105026189717c71ab905b`. Task 1's final model/export head is
`5f907cda569d97ec5c08923b35e77de5b3df06ce`; Task 2 and this reconciliation's base
is `7b804df8b0dbbd160a838cdee8681558e69dc200`. No old executable, runtime source,
module/dependency or historical evidence bytes changed. The new Go code is
confined to `research/contractcomposition` and its exporter. Task 3 changes
only results, qualification/review records, roadmap and manuscript artifacts.

The controller ran the old-package baseline and root checks during this Phase
3C.1 task. Task 1 ran the added-package checks; Task 2 ran the independent audit
and export capture. Task 3 inspected their retained outcomes and verified
current input hashes; it did not rerun these unchanged-source suites or claim
that it authored their executions. The Phase 3B.7 global-lint logs are historical
comparison evidence. No new runtime campaign, SQLite access or runtime `.envrc`
access occurred.

The [archive manifest](evidence/qualification/archive-manifest.json) records raw
and retained SHA-256 digests and original scratch paths. Compressed logs have
mtime zero and no filename header. Commands, environment overrides, exact
package argv, exit records, source manifests and original outputs are retained.
The archive excludes caches and temporary runner state. Tool versions were
Go 1.27.4 darwin/arm64, module language 1.27.3, Python 3.9.6, Staticcheck
2026.2.2 (0.8.2), and Tectonic 0.17.0.

## Package and source accounting

The baseline command was `go test -json -count=1 -timeout=60m` with the fixed
70-package argv in [baseline inputs](evidence/qualification/baseline/inputs.json),
run before new Go sources existed. Exit 0, 1,517.651 seconds. All 70 packages
have terminal events: 57 pass and 13 no-test skips, zero failures or missing
packages. It was not a run of a changing `./...` package set.

The clean current package enumeration has 72 packages, exactly those 70 plus
`research/contractcomposition` and `research/contractcomposition/cmd/export`.
Both new packages pass Task 1's final ordinary and race runs. Their final logs
are [normal](evidence/task1/normal.log) and [race](evidence/task1/race.log), with
argv/environment/exit records in the
[source manifest](evidence/task1/source-manifest.json). This combined coverage
is an ordinary qualification of the current package roster, not a fresh single
72-package invocation or a whole-repository race/integration run.

All 598 files frozen before baseline match current hashes. The later 637-file
manifest additionally covers 10 new model/test files and 29 preexisting nested
module files omitted from the initial manifest. Those 29 were separately
compared with the base Git blobs and match. They were **not** all part of the
original 598-file freeze. The [input comparison](evidence/qualification/input-comparison.json),
[package coverage](evidence/qualification/package-coverage.json) and
[Task 3 verification](evidence/qualification/task3-verification.json) record
these distinct checks. Task 3 also verified all 17 Task 1 source digests, 22 Task
2 source digests, 84 Task 2 retained-file digests and 36 compressed-artifact
raw/packed digests. No source recheck was needed after documentation-only edits.

## Executed command results

| Command/scope | Result | Attribution / retained evidence |
| --- | --- | --- |
| Baseline fixed old-package ordinary suite | Exit 0; 57 pass, 13 no-test skips | Controller; `evidence/qualification/baseline/` |
| `go test -count=1 -v ./research/contractcomposition/...` | Both new packages pass | Task 1 final `normal.log` |
| `go test -race -count=1 ./research/contractcomposition/...` | Both new packages pass | Task 1 final `race.log` |
| `go build ./...` | Exit 0; clean writable-cache run | Controller; `commands/build-cache-writable/` |
| `go mod verify` | Exit 0; all modules verified | Controller; `commands/module-verify/` |
| `go mod tidy -diff` | Exit 0; no diff | Controller; `commands/module-tidy-diff/` |
| `staticcheck ./research/contractcomposition/...` | Exit 0, no diagnostics | Controller; `commands/scoped-staticcheck/` and Task 1 final log |
| `staticcheck ./...` | Exit 1; 13 inherited diagnostics, exact byte match | Controller; `commands/global-staticcheck/` |
| `staticcheck -tags=integration ./...` | Exit 1; 31 inherited diagnostics, exact byte match | Controller; `commands/global-staticcheck-integration/` |
| Python independent audit and unittests | Genuine accepted, 8 test methods pass; 35 hostile inputs rejected | Task 2 `evidence/task2/manifest.json` and rejection report |
| Two fresh Go exports | Byte-identical to each other and frozen sample | Task 2 manifest; raw SHA-256 recorded in results |

The command directories above are relative to
`evidence/qualification/`. Global lint is **not clean**. The
[exact comparison](evidence/qualification/staticcheck-comparison.json) has no
added or removed diagnostic bytes against Phase 3B.7's 13/31 logs. The initial
root build also exited 0 but emitted a metadata-cache permission diagnostic;
the subsequent writable-cache execution is clean. Initial combined package-list
stdout/stderr likewise included a cache diagnostic; only the subsequent clean
list supplied package accounting. Both original command outcomes remain
retained, and neither warning is miscounted as a package. Task 1's earlier
Staticcheck cache failure and Task 2's failed candidates also remain in their
original task archives.

No new whole-repository integration/race suite or old model export campaign was
run. Unchanged inputs and the listed scope justify reuse, not a broader claim.

## Manuscript, bibliography and PDF

Markdown and LaTeX add matching Phase 3C.1 content, the conditional completion
formula and its quantifiers, R0/R1, the distinction between abstract eligibility
and concrete readiness, grouping grammar, acyclic sufficiency, measured counts
and remaining boundaries. Abstract and Conclusion summarize the current result.
The earlier Phase 3B closure wording is marked as that historical boundary.
The [primary citation check](evidence/manuscript/prior-art.md) discloses the
publication/author-copy pagination discrepancy and the absence of a verified
DOI. This prior art is background only.

The PDF operation marker ran successfully once with `--operation-kind edit
--expected-output-count 1 --output-format pdf`. Baseline and current source
copies were built from scratch with `tectonic --only-cached --keep-logs` under
`/tmp/alan-phase3c1/task3/{base,current}/`; both exited 0. Poppler extracted text
and rendered the current output for visual inspection. The final repository
artifact is `paper/alan-dissertation.pdf`, rather than a generic output folder.
[PDF evidence](evidence/manuscript/manifest.json) records sources, command
outputs, build logs, text, PDF and render digests, comparison and inspection.

The current PDF has 46 pages versus 45 in the baseline rebuild. Same-index text
changes occur on pages 1, 2, 4, 5, 11, 14, 22, 24, 25, 32, 34, 42-46, including
bibliography renumbering and the contents page. Visual inspection includes all
these pages and their continuations: 1-6, 11-12, 14-15, 22-26, 32-35 and 42-46.
The exact rendering disposition is retained with the PDF manifest.

Both final TeX logs have the same **four underfull hbox warnings, one overfull
hbox warning (5.2184pt), and two math-mode font warnings**. The earlier planning
shorthand of “five overfull warnings” was inaccurate: five is the combined
hbox count, not five overfull boxes. Current/baseline warning contexts agree;
there are no new warnings from the added section or bibliography. Existing
Markdown-only Section 7, different rendered section/theorem numbering, and
format-specific markup remain inherited discrepancies. This task aligns its
changed mathematical content; it does not claim global manuscript equivalence.

The retained read-only check is reproducible with:

```sh
python3 docs/foundational/phase3c1/evidence/manuscript/verify.py
```

Its [actual result](evidence/manuscript/verification.json) verifies the 38
qualification archive files, eight compressed manuscript artifacts, current
paper hashes, normalized exact new-section equality, changed Abstract and
Conclusion prose equality, and local documentation links. The artifact operation
and source/hash verification are fresh Task 3 work; they do not rerun model or
runtime suites.

## Acceptance boundary

The results combine a conditional source argument, complete small abstract
checks, factored concrete probes, and selected concrete witnesses. No count or
review substitutes for the universal representative proof obligation. The
independent task and final whole-branch review remain separate gates recorded
in [review status](06-reviews.md); delivery is controlled separately. Shared
resource/authority fan-out, arbitrary contexts, changing authority, split
effects and distributed runtime refinement remain outside this qualification.
