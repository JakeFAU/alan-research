# Task 1 independent review

## Spec compliance — APPROVE

Reviewed base `f4c4fce7459e490f1632a85ac6b87462688e0372` to head `a036b4a6133f1285fd52dc52db2eb782a3ae8439` against the task brief, copied global constraints, binding contract, and implementer report. No missing, extra, or misunderstood Task 1 requirement requiring correction was found.

- `research/residualantichains/language.go:79` completes the epsilon-closed subset graph after finding failures, tests each inclusion direction separately, and retains deterministic BFS witnesses. `language_test.go:18`, `:28`, and `:54` exercise proper inclusion, epsilon, tau cycles, nondeterminism, both directions, and a 13-event counterword.
- `research/residualantichains/antichain.go:41` collapses equivalent classes before strict MIN/MAX filtering. `:126` computes source-subset-target coverage and a shortest outside word for each source class on every uncovered target. `profile.go:32` checks the entire weak-reach graph and retains initial/FULL/MIN/MAX/J separately. `profile_test.go:24` and `:58` cover duplicate minima, empty family versus empty language, empty H versus epsilon, and multiword/directional coverage.
- `research/residualantichains/fixtures.go:5` implements the exact original two fixtures. `witness.go:102` clones the legal-trace control plant's states before changing marks, preserving the original paid-completion plants. The captured last comparison has FULL=false, MIN/MAX/J=true, equal initial language, and FULLWord=epsilon. The legal-control shortest proper differences are `[commit:BC]` and `[commit:AB,reserve-capacity:AB,commit:BC]`.
- `research/residualantichains/evidence.go:61` reconstructs all 3564 products and all 396 contexts for every J-equal pair. `product.go:35` retains full reachable relation arcs and computes CanComplete by reverse reachability. `evidence_test.go:11` retains blocked contexts, actual RF/doomed paths, and generic control observations. The archive contains 3564 unique products and 2772 unique substitution rows with matching initial/NB observations on both sides.
- `research/residualantichains/witness.go:72`, `:135`, `:171`, and `:191` preserve the RF/doomed, generic MIN-loss, generic J-nonnecessity, multiword-trie, and local-hit/global-NB controls. `antichain_test.go:8` uses direct finite-set intersections/hit predicates for exhaustive small-family checks; `product_test.go:8` checks 128 finite product contexts.
- `docs/foundational/phase3c5/01-argument.md:29` gives finite coinitiality and hit preservation; `:45` gives the product/minima law; `:60` makes Cartesian weak reach and residual intersection conditional on synchronized visible events and independent private steps; `:90` gives both directions of local finite-H characterization. Generic versus all-marked boundaries, finite-family assumptions, local versus global necessity, and lack of runtime/mechanized/novelty/speedup claims remain explicit.
- `research/residualantichains/input.go:48` authenticates the pinned predecessor before selecting descriptors and Path/SHA256 provenance. `:108` binds actual read files and distinguishes compiled inputs from predecessor provenance. `evidence.go:185` is the sole public function. `docs/foundational/phase3c5/tools/capture.py:113` uses exclusive directory creation before commands or output writes; the retained overwrite record reports unchanged content, size, and mtimes.
- `docs/foundational/phase3c5/evidence/task1/red/semantic-red.json:1` preserves actual semantic assertion failures. Focused extraction of `red/language.go.gz` confirms its insufficient baseline equated both inclusion flags with equality; the RED was not a missing-symbol failure. `false-hypotheses/` is separately labeled, and `evidence_test.go:191` keeps intentional failures out of normal passing qualification.
- The diff adds only the Task 1 package, argument/schema, helper, and evidence. Existing model/runtime/module/test/evidence files and the manuscript are unchanged. Independent audit, root checks, and manuscript reconciliation are later-task obligations, not missing Task 1 deliverables.

## Strengths

- `research/residualantichains/antichain.go:126` retains all local failures rather than only a convenient prefix, and recomputes actual hit/miss observations for their finite H sets.
- `docs/foundational/phase3c5/02-evidence-schema.md:48` documents exact graph ordering, witness sentinels, class representatives, source roles, and generic/control separation sufficiently for independent reconstruction.
- `docs/foundational/phase3c5/evidence/task1/README.md:3` explicitly selects clean final captures while preserving earlier diagnostic captures; it does not erase the module-cache warning or relabel it as pristine evidence.

## Issues

### Critical

None found.

### Important

None found.

### Minor

No outstanding change requested. Historical diagnostic noise is present in `docs/foundational/phase3c5/evidence/task1/capture-1/go-list-deps.command.json:15` and the corresponding capture-2 record. The warning is preserved and resolved by the cache override; focused inspection confirmed captures 3/4 have empty stderr and no package Error, DepsErrors, or Incomplete records. It is not an unresolved qualification defect.

## Review checks and limits

- Read the complete textual diff in bounded passes. One tool display truncated repeated capture metadata; all 49 added JSON records were then parsed mechanically from the complete diff, and artifact/source checks below verified their identities. No changed source file was separately reread for a second code-review pass.
- Named risk: binary capture/source authenticity. A read-only Python extraction verified all 70 artifact-manifest entries against current bytes and sizes, byte equality of all four exports and stdout archives, all 54 semantic source hashes, all 40 inherited Path/SHA256 entries, and all 19 final freeze identities. Imported alphabet, seven monitors, 72 plants, and 324 controllers exactly match the predecessor descriptors. Raw export identity is 14097260 bytes, SHA-256 `ab3cbdfda577a3ae804eb76a94a8c7ac665cf00c69f992ad6804b6e0d0e27b8a`.
- Named risk: compiled-source provenance being bound to an incidental absolute checkout. Parsed the archived final dependency streams and independently normalized package directories relative to their recorded Module.Dir. Both streams contain 109 packages, no package errors, and the same ten compiled local files as the semantic manifest. No dependency command was rerun.
- Named risk: proof/evidence correspondence for strict reduction and plant marking. Focused binary extraction checked the strict-fixture flags, legal-control differences, unchanged original descriptor marks, unique product/substitution cardinalities, and pairwise substitution observations. This was not a substitute for the later independent full graph auditor.
- Reviewer-check diagnostic: the first provenance comparison incorrectly compared the predecessor's full records, which additionally contain Use, with the documented Path/SHA256 projection. That assertion failed; the corrected projection comparison passed. This was a review-script expectation error, not a producer failure or edited evidence.
- Read the existing ordinary/race/Staticcheck command records and their clean results. Did not rerun those suites, any exporter, the runtime, or SQLite. No source, semantic evidence, index, or Git state was modified. The only write is this requested review report.
- Cannot verify execution chronology merely from a diff or hash record. Archived commands/results are retained evidence, while current hash agreement establishes correspondence to reviewed bytes. Later independent hostile-input audit and root/manuscript qualification remain controller-owned.

## Quality — APPROVE

The finite-model implementation follows the stated semantic boundary, separates concerns into small files, and supplies meaningful negative and finite-set checks. Source arguments and measured certificates are consistently bounded; no blocking correctness or evidence-quality defect was found in this task review.
