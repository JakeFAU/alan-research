# Task 1 report: residual interfaces

Status: implemented and qualified within Task 1; independent review/checker and
manuscript/root qualification remain controller-owned. Base was
55153b7e0eff91b9c5df51b95e3c4dddba03d838, with immutable predecessor source base
c369e25c8db268d418ab7c11573d46dbc1d9425f.

## Question and result

The question is whether retaining each individual residual marked language after
every visible prefix suffices to preserve contextual nonblocking. The source
argument establishes a conditional sufficient rule for compatible finite
synchronized epsilon-NFAs with independent private steps and conjunctive marks.
It gives both directions of derivative union, Cartesian weak reachability,
residual intersection, the exact nonblocking criterion, and product-profile
closure. It is not a mechanized proof or Theorem 4.8 closure.

The finite evidence has 21 full monitor profile comparisons, all 223 declared
cross-state language comparisons and their complete subset certificates, plus
all 2,772 products reconstructed from pinned automata. The predecessor products,
comparisons, grouping, histories, outcomes and counts are not decoded as oracles.
Exact trusted fixture bytes and their 27 source hashes authenticate the selected
alphabet, 72 K, five inherited M and 324 D descriptors. Two new all-marked M
fixtures add duplicate residual alternatives and a private cycle.

Six profile-equal pairs are all unordered pairs of release-first,
release-first-tau, release-first-duplicate and release-first-loop. Every one is
checked in 396 contexts, keeping 792 nonblocking and 1,584 blocking rows. The
release-first/doomed pair has equal initial language but unequal profile and
opposite core nonblocking. Its shortest monitor mismatch is reserve:AB, while
the genuine product paths first commit AB. The loop fixture has an unpaid,
unmarked tau cycle with completion escapes from both cycle states. A generic
all-unmarked strictness example is exported separately, outside the all-marked
ALAN monitor domain.

## Changes and implementation boundaries

Only new Task 1 files were added:

- `research/residualinterfaces/`: focused automaton, input, fixture, language,
  profile, product, witness and evidence files; focused tests; sole public
  function `Export(io.Writer) error`; JSON exporter CLI; README.
- `docs/foundational/phase3c4/01-argument.md`, `02-evidence-schema.md`.
- `docs/foundational/phase3c4/evidence/task1/`: two captures, actual command/output/
  exit records, source/artifact/compiled-input manifests, qualification and
  deliberate-hypothesis evaluators.
- This task report.

The local finite-relation representation and ordinary language/product traversal
conventions are adapted from the predecessor package; no predecessor packages
are imported, modified or executed. This is not algorithmically independent
validation. Task 2's Python reconstruction is the independent challenge.
Signatures are opaque authenticated JSON descriptors, copied exactly when M is
rebound; this is deliberately not a public untrusted-automaton decoder.

Controller clarification: reverse-reachability CanComplete/Nonblocking labels
apply to product graphs. Decision-certificate graphs instead retain complete
subset-pair nodes/arcs and actual language acceptance bits or bidirectional
residual coverage. This distinction is specified in the schema. Profile coverage
permits one-to-many matches and explicitly retains empty-vs-dead mismatch and
vacuous opposite separator lists.

## TDD and resolved diagnostics

All records are under `docs/foundational/phase3c4/evidence/task1/`.

1. `red.json`: `go test -run
   TestLanguageBoundaries|TestProfileMultiplicityAndEmptyLanguage
   ./research/residualinterfaces` exited 1 before implementation: undefined
   edge/automaton/state/accepts/labels/compareLanguages. This is missing-symbol
   compilation RED, not semantic falsification.
2. `focused-language.json`: same focused command exited 0 after implementation.
3. `focused-products.json`: focused roster/products/signature/small-context tests
   exited 1. All semantic reconstruction tests passed; the root mutation test
   failed because its literal string replacement did not match indentation in
   the predecessor JSON and therefore changed nothing. Root was Q/Limit 1.
   Corrected the test to decode/change/re-encode the root descriptor.
4. `focused-root-unused-import.json`: first correction command exited 1 due to
   leftover unused bytes import. Removed that test import.
5. `focused-root-corrected.json`: corrected role/signature/root/external-action
   mutation tests exited 0. Neither diagnostic was suppressed or relabeled clean.
6. Controller caught generic strictness using commit:AB instead of the contract's
   a=apply:AB. Added the exact shortest-prefix assertion;
   `strictness-label-red.json` records exit 1 of `go test -v -run
   TestCompleteEvidenceAndWitnesses ./research/residualinterfaces`. Corrected
   fixture/argument/schema and retained the former captures and runs distinctly.
   The affected ordinary/race tests and final captures pass with apply:AB. Also
   corrected historical attribution from physical histories to concrete model
   histories; no physical/runtime qualification is implied.

Tests cover epsilon marking via private closure/cycles, empty subsets, a shortest
separator of length 12, nondeterministic equivalent branches, one-to-many
coverage, dead-vs-unreachable profiles, inherited roster/identity, external-action
veto/state-change rejection, full product graph integrity, core negative and
loop paths/escapes, and product-profile closure for 128 independently constructed
two-state contexts with marks, private steps and nondeterministic visible edges.
They also check Cartesian reachability for epsilon, repeated visible events and
an unreachable label. No tests were weakened to obtain a pass.

## Final qualification actually executed

The first frozen candidate received the full scoped ordinary/race runs below.
After controller correction of the generic fixture label and historical wording,
only the affected complete-evidence test was rerun ordinarily and with race,
plus scoped Staticcheck, dependency inspection and both final captures. Earlier
full-run records and captures are retained in `pre-strictness-correction/`; final
records are in the Task 1 evidence root. All use
`GOCACHE=/tmp/alan-phase3c4/go-cache` and
`STATICCHECK_CACHE=/tmp/alan-phase3c4/staticcheck-cache`:

- `go test -v -timeout 30m ./research/residualinterfaces/...`: exit 0; six top-level
  tests and five signature subtests pass on the pre-correction candidate; exporter
  package has no tests. Final affected run adds `-run TestCompleteEvidenceAndWitnesses`
  to that command and exits 0 with the corrected exact generic label assertion.
- `go test -race -v -timeout 30m ./research/residualinterfaces/...`: exit 0; same
  tests pass, no race diagnostics on the pre-correction candidate. Final affected
  run adds `-run TestCompleteEvidenceAndWitnesses`; exit 0, no race diagnostics.
- `staticcheck ./research/residualinterfaces/...`: exit 0, empty stdout/stderr.
- `go list -deps -json ./research/residualinterfaces/cmd/export`: exit 0, complete
  output retained gzipped; exact nine compiled local non-test Go files match
  the compiled-source manifest, with no old package dependency.
- `go run ./research/residualinterfaces/cmd/export`: twice, each exit 0, empty
  stderr, raw output byte-identical.
- `python3 docs/foundational/phase3c4/evidence/task1/falsify.py union-sufficiency`:
  actual exit 1, expected rejection of the deliberately false union hypothesis.
- Same evaluator with `inevitable-completion`: actual exit 1, expected rejection
  of the deliberately false abstract inevitability hypothesis. These have their
  own exact command/stdout/stderr/exit logs, separate from TDD and passing tests.

No cache warnings/failures occurred with these writable caches. No whole-root
suite, global lint, old campaign, concrete capture, live model, endpoint, SQLite
open or runtime experiment was run. Controller handles root checks separately.
The evidence driver can reproduce these commands; it checks expected exit 1 for
the two falsification evaluators and ordinary exit 0 for qualification commands.

## Measured output and identities

| Measure | Value |
| --- | ---: |
| K / M / D descriptors | 72 / 7 / 324 |
| Products | 2,772 |
| Reachable product states / arcs | 29,156 / 35,672 |
| CanComplete states | 14,410 |
| Initially completable / nonblocking products | 886 / 656 |
| Comparisons / initial-language equal / profile equal | 21 / 10 / 6 |
| Cross-state residual comparisons | 223 |
| Language subset-pair / profile subset-pair nodes | 850 / 89 |
| Substitution rows / nonblocking / blocking | 2,376 / 792 / 1,584 |

Predecessor gzip SHA-256:
`86a9d22f9808a67dd0066b6b77dd2b3c01a2125d10863c6fbe1369097f6624c8`.
Predecessor raw: 25,384,964 bytes;
`a073a22fe8ac62e05f2d15b0f5146b19db9606e38ce4619261e7e4108e27269d`.
New raw output: 8,271,854 bytes;
`b40f1d5ac0aea0f43999502ea318cb128f545eb8fb1fd129451915a5aa84bf1f`.
Both captures have that raw hash. New deterministic gzip: 157,341 bytes;
`35115c3104cddd8a789c80a8c0e73eb7c7f173632af512776ba8502b61b9afc9`.
The 40-entry source manifest's compact JSON SHA-256:
`185783ec02fe37b77e0dee8bbec5dabb47a3acd8cd45cce5b7233ba654549a31`.
Exact per-file hashes and uses are in `sources.json`; artifact identities and
counts in `artifact-manifest.json`; compiled roster in `compiled-inputs.json`.

## Self-review and remaining obligations

No old source/evidence/module/runtime file was changed. Complete scoped diff,
status, generated artifact identities and source manifest were inspected before
commit. No unrelated/controller files, temporary instrumentation, profiler output
or Python cache artifacts are included.

The lasso is a deliberate abstract finite-monitor fixture. Its compilation or
admission by the strongly terminating runtime local instruction grammar with
bounded loops/forward control flow is unestablished; it does not allege private
divergence of validated ALAN local programs. Generic strictness does not prove
strictness within the all-marked ALAN roster. A mismatching profile is not
necessarily a contextual nonblocking separator.

The new conditional source argument and bounded Go evidence do not establish
weakest-interface necessity, state minimization, verification cost improvements,
exact-cost equivalence, fairness/liveness, arbitrary refusal/divergence contexts,
full-agent substitution, shared-authority or distributed/split-effect claims.
Historical concrete evidence remains attributed to its predecessor packet.
Independent checker/review, manuscript synchronization/PDF, root qualification,
and delivery remain controller-owned and are not claimed complete here.
