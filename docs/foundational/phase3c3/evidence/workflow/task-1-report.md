# Task 1 report: boundary automata, substitution attack and captures

## Status and commits

Task 1 implementation is complete and locally qualified, awaiting independent
review and Task 2 reconstruction. Implementation/evidence commit:
`c8230eb2ed2bf3cd45635f72d5eef84bcc7372fa`.

Base used: `e7408eb493c9e2703c9347a565b8ac58d4de9202` on
`codex/phase3c3-component-interfaces`. This local handoff report is excluded by
the configured ignore rule for `.superpowers/sdd/phase3c3-plan`; an ordinary
`git add` declined it, so no report commit was created. The report remains at
the requested path. No implementation was amended or rewritten. No push, PR,
merge, manuscript update or independent-review claim is made by this worker.

## Question and result

The investigated negative claim was that the full inherited concrete state,
exact holdings, configuration/capabilities and static service/release rules
suffice to determine completion after replacing B's continuation monitor. It is
falsified: histories 1/2 share those exact fields at the genuine four-event cut,
with opposite residual product CanComplete. All five M fixtures have every state
marked, so private acceptance cannot manufacture the failure.

The positive result is a conditional source argument: under exact signature
agreement, shared visible synchronization, independent private state and
conjunctive marking, marked finite visible language of a product is the
intersection of factor languages (both directions supplied). Equal languages of
pointed monitors therefore preserve **initial existential possible completion**
in every compatible finite context under these assumptions. A corresponding
statement for residuals requires equality of those residual languages.

Finite evidence checks the frozen roster, exact NFA languages and grouping. It
also falsifies extension to every-reachable-state nonblocking: release-first
and release-first-doomed have equal initial languages, while the latter has a
reachable debt-holding doomed residual. Both its successful and doomed branches
are genuine captured histories. Stable-ready's old every-state result is not
imported after imposing monitor restrictions.

## Implementation and interfaces

Created `research/componentinterfaces/` with automaton/signature/role validation,
exact epsilon-closure/subset-pair comparison, full synchronized/private product
BFS, reverse CanComplete, nonblocking classification, inherited K adapter, exact
M and D fixtures, concrete wrapper/captures, deterministic exporter and focused
tests. The sole exported function is `Export(io.Writer) error`.

Created `docs/foundational/phase3c3/01-argument.md`, the complete versioned
`02-evidence-schema.md`, and `evidence/task1/`. New declarations are unexported
except Export. Label IDs map to full old exact actions via the signature/alphabet;
tau is separate and has no forged empty-binding visible action. Product source
keys are explicit flat tuples and compact JSON IDs. Grouping archives both full
parenthesized relations in addition to each flat graph, not a counts-only result.

Concrete visible steps use unchanged `capacitycomposition.Lift`; preparation
retains M. Each resulting step is checked against the real wrapped successor
relation. Private steps require B up; a requested tau from down B first uses
actual recovery, preserving retained PC. Captured post-states are never repaired.
No production code, dependency, old model, old evidence or manuscript was changed.

## Measured evidence

All totals below are configuration sums, not counts of unique semantic states
across configurations. Grouping copies are excluded from product totals.

| Roster | Products | States | Arcs | CanComplete states | Initial CanComplete | Nonblocking products |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| 72 K × 5 M | 360 | 4,868 | 6,260 | 4,672 | 304 | 236 |
| Core K × 324 D × 5 M | 1,620 | 14,676 | 17,160 | 4,854 | 318 | 156 |
| Total | 1,980 | 19,544 | 23,420 | 9,526 | 622 | 392 |

- 72 exact K descriptors, five exact M descriptors, all 324 D descriptors,
  including all declared unreachable D states and empty-mark controllers.
- Ten exact language comparisons, three equal pairs, 43 complete subset-pair
  states total. Inequalities retain shortest words and actual acceptance checks.
- 1,620 grouping comparisons, both full parenthesizations matching each flat
  relation after explicit tuple normalization.
- 1,188 equal-language pair/context substitutions: 396 successful and 792 blocked.
  These compare initial CanComplete, never arbitrary residuals.
- Eight genuine concrete histories: 123 events. Per-history event counts are
  5, 18, 18, 18, 18, 21, 18, 7. Histories 1/8 end noncompletable; all other six
  end complete. History 8 retains PC 3 after recovery.

Core K product facts: free has 24 states/34 arcs and initial completion but is
not nonblocking; commit-first has 4/3 and no initial completion; release-first
has 13/16 and is nonblocking; release-first-tau has 14/17 and is nonblocking;
release-first-doomed has 14/17, initial completion, and is not nonblocking.

## RED and GREEN evidence

Every execution record below retains exact argv, explicitly set environment,
exit code, stdout and stderr. The export records refer to the exact raw stdout
artifact via its decompressed gzip hash. Go was `go1.27.4 darwin/arm64`;
Staticcheck was `2026.2.2 (0.8.2)`. The module's existing Go directive is unchanged.

| Record under evidence/task1 | Actual command / result |
| --- | --- |
| `red-01.json` | `go test -v ./research/componentinterfaces/...`: exit 1, expected missing model symbols; compile failure, **not** semantic falsification |
| `check-02.json` | Same focused command: exit 0 after automata/product/concrete implementation |
| `red-03-export.json` | Same command: exit 1, missing exporter/evidence symbols; compile failure, **not** semantic falsification |
| `red-04-false-hypotheses.json` | `go test -v -run TestAttack ./research/componentinterfaces`: exit 1, both deliberately false semantic hypotheses falsified |
| `green-test.json` | `go test -v -timeout 30m ./research/componentinterfaces/...`: exit 0; package 2.079 seconds |
| `green-race.json` | `go test -v -timeout 30m -race ./research/componentinterfaces/...`: exit 0; package 20.877 seconds |
| `green-staticcheck.json` | `staticcheck ./research/componentinterfaces/...`: exit 0, no diagnostics |
| `source-dependencies.json` | `go list -deps -json ./research/componentinterfaces/cmd/export`: exit 0, no diagnostics, exact compiled local source set retained |
| `source-binding-check.json` | Exact stored `python3 -c` program: exit 0; all source hashes current, compiled local Go set exactly bound, every monitor state marked, artifact raw hash/length match |
| `artifact-manifest.json` | Two `go run ./research/componentinterfaces/cmd/export` executions: both exit 0, no stderr, byte-identical stdout |
| `toolchain-*.json` | Actual go/staticcheck versions, relevant Go environment, and clean `git diff --cached --check` |

Go commands used `GOCACHE=/tmp/alan-phase3c3/gocache`. Staticcheck's successful
run additionally used `STATICCHECK_CACHE=/tmp/alan-phase3c3/staticcheck-cache`.
The clean dependency enumeration additionally used
`GOMODCACHE=/tmp/alan-phase3c3/gomodcache`.

The first Staticcheck attempt set XDG_CACHE_HOME, which macOS's cache resolution
did not honor, and failed on an unwritable normal cache. Its exact failure is
retained in `failed-staticcheck-cache.json`; explicit STATICCHECK_CACHE fixed it.
The first go-list enumeration exited 0 but emitted a normal-module-cache write
warning; it remains in `source-dependencies-cache-warning.json`. The explicit
GOMODCACHE rerun is clean. These are tooling/cache outcomes, not semantic failures.

The false hypothesis source is retained as `false-hypotheses_test.go.txt`, with
reproduction instructions in README. It was moved out of the ordinary suite
because it asserts claims the experiment disproves; the passing suite explicitly
checks the counterexample behavior. No desired invariant test was weakened.
Initial model/export REDs were compilation failures and are reported separately.

The final ordinary/race/Staticcheck checks followed the final Go edits. No
unchanged suite was rerun for repetition. Old million-case probes and the
whole-repository suite were not rerun, per controller scope. Controller owns
fresh root build/module checks. Independent Python reconstruction is pending.

## Deterministic artifact and source boundary

The final exporter ran twice, taking 0.863 and 0.569 seconds respectively; these
are capture durations, not controlled performance measurements or speedup claims.
There were no long-running captures. Raw stdout: 25,384,964 bytes, SHA-256
`a073a22fe8ac62e05f2d15b0f5146b19db9606e38ce4619261e7e4108e27269d`.
Deterministic gzip: 403,411 bytes, SHA-256
`86a9d22f9808a67dd0066b6b77dd2b3c01a2125d10863c6fbe1369097f6624c8`.
Zero gzip mtime, empty filename. Semantic stdout has no timestamp or temporary
path. Mutable plan, report and review metadata do not affect its source hashes.

Sources binds 27 inputs: all 22 actual non-test local Go files compiled for the
producer, root go.mod/go.sum, frozen contract, source argument and schema. The
22-file set was compared to actual `go list -deps -json` package GoFiles, not
assumed from imports alone. Old imported packages' unused compiled declarations
are included. Toolchain version/environment is recorded separately; standard
library binaries are not individually hashed. Source hashes remained current
through the implementation commit; this report is deliberately outside the
semantic manifest.

## Self-review and remaining obligations

Inspected full task source/documentation diffs, test assertions, staged file
scope and retained evidence records; verified archive/source hashes and complete
compiled source coverage. Git diff whitespace check passed. No unrelated changes
were staged. The task commit does not claim independent review approval.

The proof is conditional source reasoning; finite graph/language checks and
concrete model histories are different evidence categories. Go adjacency checks
share owning model implementations. The independent checker must reconstruct
complete relations and reject hostile inputs rather than trust producer Booleans,
counts, grouping flags or source-level self-review. It has a complete frozen
schema and artifact to consume.

No exhaustive concrete-product claim, runtime refinement, complete-agent
replacement, arbitrary topology, dynamic authority, split/unknown effects,
distributed allocator, universal nonblocking, necessity/full abstraction for the
role-restricted ALAN grammar, optimal policy or general Phase 3C closure follows.
Task 3/controller still own manuscript/roadmap synchronization, literature/PDF
checks, root verification, independent review and delivery. No unresolved Task 1
implementation blocker remains; the cache incidents above are preserved and
resolved. No subagent/reviewer was dispatched by this worker.
