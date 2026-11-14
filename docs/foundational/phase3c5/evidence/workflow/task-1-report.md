# Task 1 report — finite residual antichains

Status: DONE (independent review remains controller-owned).
Base HEAD: `f4c4fce7459e490f1632a85ac6b87462688e0372`.
Commit: `a036b4a6133f1285fd52dc52db2eb782a3ae8439` — Add finite residual antichain model and complete evidence.
Post-commit `git status --short`: clean (the ignored workflow report is intentionally outside the commit).

## Question and conditional result

Can inclusion-minimal completion alternatives, paired with initial marked
language, replace FULL residual profiles for the declared finite synchronized
grammar? The source argument gives finite-family coinitiality, exact hit-test
preservation, the product/minima law, conditional nonblocking and J closure, and
the exact local finite-H implication characterization in both directions.
It is not a mechanized proof, runtime refinement, novelty claim, speedup claim
or global contextual necessity theorem. Generic marked controls remain distinct
from the all-marked monitor roster, and one whole unchanged K owns effects,
shared authority and Q.

## Changes and interface

All new Go code is confined to `research/residualantichains/` and its exporter.
The sole public function is `Export(io.Writer) error`, called from repository
root with `go run ./research/residualantichains/cmd/export`. It emits deterministic
`alan/phase3c5/v1` JSON. The package contains focused automaton, input, language,
antichain, profile, product, witness and evidence files plus five focused test
files and README. Small unchanged predecessor automaton/product mechanics were
copied into this package rather than importing or executing predecessor code;
new directed inclusion, class reduction and local coverage/hit logic is local.

Created `01-argument.md`, `02-evidence-schema.md`, `tools/capture.py` and the
complete `evidence/task1/` packet. Schema freezes graph/record ordering,
canonical classes and dominators, every local failure H, null versus epsilon,
all products, controls, source roles, counts, and capture normalization for Task 2.
No existing model/runtime/module/test/evidence file was modified. No manuscript,
root checks, Python auditor or later-task implementation was attempted.

## Exact inputs and source binding

The Phase 3C.4 pinned gzip was checked as length 157341 and SHA-256
`35115c3104cddd8a789c80a8c0e73eb7c7f173632af512776ba8502b61b9afc9`;
raw length 8271854 and SHA-256
`b40f1d5ac0aea0f43999502ea318cb128f545eb8fb1fd129451915a5aa84bf1f`.
Its schema and all 40 inherited source/provenance entries are verified.
Only descriptors/alphabet/signatures and provenance are decoded; old outcomes
are not oracles. The rejected preflight fixture is not instantiated.

The semantic manifest has 54 actual read inputs, including 10 compiled-local
Go files. Go-list's 109 dependency-package records are retained unmodified.
Both final captures have exit 0, empty stderr, no Error/DepsErrors/Incomplete
package records, and exact module-relative compiled roster agreement. Normalization
uses package.Module.Dir and Module.Path, not equality to this checkout path.
Module files and every inherited source hash are bound distinctly from compiled
producer sources. `final-dependency-validation.json` records those checks.

Final freeze: `evidence/task1/source-freeze-final.json`, SHA-256
`f3d0e4868d22ee7eb9a48cbba2eaa4159ee293ccd058912793f1a44c07c1033c`.
The Go, test, argument and schema hashes are unchanged since the initial freeze;
only the capture helper's cache environment changed, as detailed below.

## Measured results

- 9 monitors, 72 plants, 324 controllers, 36 unordered comparisons.
- 3564 complete products: 45800 nodes, 59492 arcs, 22716 CanComplete nodes;
  1234 initially possible completions and 760 nonblocking products.
- Equal initial languages: 13 pairs; equal FULL: 6; equal MIN: 7;
  equal MAX: 13; equal J: 7. All 13 MAX cases remain explicitly retained.
- 1442 residual certificates counting cross and repeated within-monitor pairs;
  5208 complete language subset pairs, 146 complete profile subset pairs.
- 2772 J-equal substitution rows across every one of 396 contexts: 844 both
  nonblocking, 1928 both blocking. Initial-possibility flags also agree in every row.
- 81 failed local coverage targets, including 28 empty H; largest H has 3 words.
  Every source-class shortest outside word and every hit/miss check is retained.
- Each original new fixture has 52 nonblocking and 20 blocking K contexts;
  all 324 K/D contexts are blocking. None were discarded to improve the result.
- minimal-base/minimal-augmented: initial equality, FULL inequality, all-word
  MIN/MAX/J equality. Both proper residual differences survive legal core-K traces:
  L4 minus L1 has shortest legal word `[commit:BC]`; L0 minus L4 has shortest
  legal word `[commit:AB,reserve-capacity:AB,commit:BC]`. These use only the separate
  all-marked legal-trace control, not modified paid-completion plant marks.
- RF/doomed MAX failure retains actual core paths to healthy M=1 and doomed M=3
  at the same unpaid K state. Generic MIN-only loss, J nonnecessity, multiword H
  and local-hit/global-NB failure have explicit finite trie products and stay
  outside the nine-monitor totals.

## Tests and actual failures

Initial semantic RED command (record plus source snapshots in `red/`):
`GOCACHE=/tmp/alan-phase3c5/go-cache go test -v -run 'TestDirectedInclusion|TestStrictReduction|TestGenericMINInformationLoss' ./research/residualantichains`.
Exit 1, with all three intended semantic assertion failures against the initial
insufficient equality/FULL baseline; no missing-symbol failure was used.

Focused implementation checks then passed. Final fresh qualification is retained
in `qualification/` with exact commands, output and exit:

1. `go test -v -count=1 -timeout 30m ./research/residualantichains/...` — exit 0.
2. `go test -race -v -count=1 -timeout 30m ./research/residualantichains/...` — exit 0.
3. `staticcheck ./research/residualantichains/...` — exit 0, empty stdout/stderr.

These use `GOCACHE=/tmp/alan-phase3c5/go-cache` and
`STATICCHECK_CACHE=/tmp/alan-phase3c5/staticcheck-cache`. Twelve top-level tests
pass; one deliberate false-hypothesis test is gated/skipped in normal qualification,
with five explicit compatibility subtests. Boundary checks cover epsilon,
leading/trailing tau, cycles, nondeterminism, empty subsets/families/languages,
both directions, strict/non-strict dominance, duplicate classes and a 13-event
counterword beyond short sampling. Finite-set checks exhaust all 256 families
over 8 languages on {epsilon,r,a}, all 2048 family/suffix hit cases, all 65536
family pairs with eight H tests each, plus 64 actual trie products and 128 small
context/product-closure cases. The complete ALAN roster is rebuilt in tests.

Three separately captured intentional failures (`false-hypotheses/`) run the
final `TestFalseHypotheses` with `ALAN_FALSE_HYPOTHESIS` set to
`MIN-preserves-initial`, `MAX-preserves-NB`, and `local-hit-implies-NB`; each exits
1 with its expected semantic falsification. Normal tests assert those same
counterexamples successfully. No other test failures occurred after the initial RED.

## Final captures and diagnostic preservation

Use `evidence/task1/capture-3/export.json.gz` as the final fixture and `capture-4/`
as its fresh repetition. The helper was invoked with each fresh directory and
refused an existing directory without changing file hashes, lengths or mtimes.
Final captures contain exact commands/results, gzip stdout, the conventional
export gzip, module/root-aware dependency evidence, normalized compiled roster,
all read-input hashes, environment and artifact manifests.

Both final exports and the two earlier ones have identical semantic bytes:

- raw bytes 14097260; SHA-256
  `ab3cbdfda577a3ae804eb76a94a8c7ac665cf00c69f992ad6804b6e0d0e27b8a`;
- deterministic gzip bytes 324201; SHA-256
  `5b9bfc9e9f63ab056013671eef4c42705c40732d6ec3d4d3558a75269770a5f4`.

`final-capture-equality.json` and `final-overwrite-refusal.json` retain the checks.
Captures 1/2 remain unchanged: go-list exited 0 but emitted a nonfatal default
module stat-cache permission warning. Complete records and the original helper
snapshot are retained. The helper now sets:

```
GOCACHE=/tmp/alan-phase3c5/go-cache
GOMODCACHE=/tmp/alan-phase3c5/mod-cache
STATICCHECK_CACHE=/tmp/alan-phase3c5/staticcheck-cache
```

Go is installed locally at `/opt/homebrew/Cellar/go/1.27.4/libexec` (go1.27.4).
The new module cache was sufficient for this stdlib-only exporter; broader root
checks may need the pinned SQLite dependency seeded from existing cache. No root
checks were run here. `toolchain.json` records actual tool/environment output.

The helper is not read or hashed by the Go producer, by explicit schema design;
its bytes are separately bound in capture environment records and the final
freeze. Therefore changing only its writable cache override leaves semantic
bytes and all embedded source hashes identical. The original helper hash remains
in `source-freeze.json`; the final helper hash is in `source-freeze-final.json`.
No semantic artifact was edited. No unchanged Go suite was repeated after this
capture-only fix. Clean captures have empty go-list stderr and no package errors.

## Self-review and remaining obligations

Reviewed the source, complete staged diff/status/stat, schema against exported
fields, source hashes and retained command diagnostics. `git diff --cached --check`
passed. Only Task 1 files are staged; no unrelated/generated runtime files or
investigation instrumentation are included. Evidence captures and RED snapshots
are intentional research artifacts. Independent review, independent Python hostile
input audit, root qualification and manuscript/PDF reconciliation remain later
controller-owned work. There is no claim that this task establishes runtime
refinement, distributed recovery, effectful full-agent substitution, global
contextual necessity, refusal/divergence equivalence or unrestricted Phase 3C.

Controller handoff note: the controller reports the pinned SQLite module has now
been seeded into the writable module cache, with a record at
`/tmp/alan-phase3c5/module-cache-seed.json`. This preparation was not performed
or qualified by the Task 1 builder.
