# Task 1 report — Phase 3C.2 model, adapter and frozen export

Commit: `60d46a43dd32caa07f0b814e6741213f506a9cb6` —
`Model shared-capacity reservation contracts`. All 33 committed files are
Task 1 scope. Post-commit task-owned diff is empty; only controller-owned
`evidence/preflight/` remains untracked. This ignored report is handed off
separately for the controller to archive.

## Scope and result

Base: `072cf44cc89879d60bf49233d3d4dc541209df0d`, branch
`codex/phase3c2-shared-capacity`. Research-only new package and evidence. No
runtime, dependency, historical model/evidence, SQLite or `.envrc` changes or
access. No subagents, push, PR or merge. The controller's preflight evidence is
outside this task's staging scope.

Question: can exact authorization and capacity conservation coexist with loss
of finite completion when the fixed AB/BC chain shares a reusable slot, and
which retained holder identity/readiness conditions preserve that question?

The finite model establishes the requested negative pairs within its domain:
early/shared-1 and ready-now/shared-1 each have a noncompletable AB-holder state
for each serial-B mask. Genuine concrete prefixes have equal entire base
transport/physical state and equal occupancy but opposite completion labels
when the holder changes. Ready-now loses readiness at BC commit. Stable-ready
has a finite completion suffix from every reachable abstract state in the
roster; shared-two, separate-one and servicing-B controls also pass. The old
AB-before-BC-commit trace remains an explicit restriction witness.

The universal concrete/resource-aware correspondence and stable-ready
completion arguments are conditional source proofs, not mechanized proofs or
finite enumeration of the entire concrete product. The new resource-aware
abstraction is the lifting target. Old-q universal lifting is expressly not
claimed. Permission, cumulative E and reusable Q remain distinct; E sufficiency
is a premise, with no new E counter/refund semantics.

## Implementation

* `model.go`: strict fixed-domain config/state validation, immutable root
  descriptors and exact memberships, old named stages plus original-held bits,
  three reservation policies and atomic application/retention/capacity release.
  Invalid numeric capacities are rejected, never coerced. A private base
  template avoids allocating a complete config for every validation call.
* `compose.go`: reuse of unchanged owner/port/binding joins, additional exact
  root/membership/policy validation, two regroupings and strict closure.
  Owned A/B originals cannot be silently omitted. Shared references identify
  one pool; joins do not add limits or merge dynamic histories.
* `concrete.go`: wrapper around the unchanged concrete transport relation;
  apply is filtered on an existing admitted hold and clears it atomically.
  Reservation uses actual admission/up callee. Finite lifts recover owners and
  route real custody, with no other original/hold progress during preparation.
  Guards/lifts never read WorldResult.
* `explore.go`: complete deterministic BFS with every arc and independent
  reverse traversal of that graph for existential completion labels.
* `evidence.go`: ten genuine initial-custody witness histories, exact adjacency
  checks, old erased-capacity custody-ledger checks, and full grouping records.
* `probe.go`: explicitly synthesized factored domain, full projection, every
  concrete successor, every enabled abstract action, legal lift adjacency,
  no unrelated progress, and WorldResult-erasure independence comparisons.
* Exporter: deterministic `alan/phase3c2/v1` JSON, all records and source SHA256s;
  strict zero-argument CLI. Encoder smoke tests avoid repeating expensive probes.
* Argument/schema/README: exact conditional claims, API failure conventions,
  complete schema, witness roster and precise operational-probe limitations.

## RED/GREEN evidence and commands

All Go commands used `GOCACHE=/tmp/alan-phase3c2/gocache`. Retained logs are under
`docs/foundational/phase3c2/evidence/task1/`.

Before implementation of each boundary, tests ran against intentionally inert
compiling stubs. The retained RED logs are genuine command output:

| Archive | Command | Observed RED |
| --- | --- | --- |
| `red-model.log` | `go test ./research/capacitycomposition/...` | Five contract tests failed: missing exact root, missing commits, empty graph and excluded-trace setup. |
| `red-compose.log` | `go test ./research/capacitycomposition/...` | Regrouping mismatch and omitted leaf membership. |
| `red-concrete.log` | `go test ./research/capacitycomposition/...` | Four concrete tests failed on missing actual commit. |
| `red-evidence.log` | `go test ./research/capacitycomposition/...` | Missing witness histories and unimplemented factored probe runner. |
| `red-export.log` | `go test ./research/capacitycomposition/cmd/export` | Empty JSON output and ignored CLI argument. |

Those shell capture commands printed their saved logs afterwards, so the shell
wrapper's final status was the successful `cat`, not the Go test status. The
archived Go output itself records FAIL. An earlier composition-stub run exposed
a test setup indexing panic; the setup guard was corrected before the retained
assertion-based RED capture. A transitional compile attempt after the evidence
stub was replaced failed until `probe.go` was added; it is not a semantic result.

Targeted GREEN checks were run after the model, composition, concrete and
witness implementations. A targeted unchanged predecessor check,
`go test -v ./research/contractcomposition -run '^TestEveryOldEdgeRepresentativeProjectionAndLifting$'`,
passed and measured its existing 311 old states, 24,258 factored representatives,
184,256 micro arcs and 14,058 named-edge lifts. This was not a full repository
suite. `initial-probes.log` preserves the first new full probe pass (222.22s), run as
`go test -v ./research/capacitycomposition -run '^TestFactoredRepresentative'`.
After that run, exact configuration validation was changed to avoid repeated
allocating deep equality and failure diagnostics were expanded. The final
ordinary probe rerun was therefore required by a source change. Domain and
expected semantic conditions were unchanged; no performance result is claimed.

Final verification:

* `go test -v -timeout 30m ./research/capacitycomposition/...` — PASS, exit 0,
  116.511 seconds for the recorded wrapper; all 17 tests across two packages.
  Full log: `test.log`.
* `go test -v -race -timeout 30m ./research/capacitycomposition/...` — PASS,
  exit 0, 632.504 seconds for the recorded wrapper; all 17 tests, including the
  full probe domain (627.29s), with no race reports. Exact command and result
  are retained in `race.log`. Timeout is explicit to
  accommodate the unchanged full probe domain under race instrumentation.
* `env STATICCHECK_CACHE=/tmp/alan-phase3c2/staticcheck-cache staticcheck ./research/capacitycomposition/...`
  — PASS, exit 0, no diagnostics (`staticcheck.log`). The initial invocation
  failed only because Staticcheck tried its default cache outside writable
  roots; that failure is preserved in `staticcheck-cache-failure.log`.
* Two `go run ./research/capacitycomposition/cmd/export` captures and `cmp -s`
  — PASS, both exit 0 (150.045s and 149.243s), comparison exit 0;
  exact commands/results in `export.log`. The checked archive is
  `export.json.gz`, deterministic gzip of the exact plain stdout. Both raw
  captures remain under `/tmp/alan-phase3c2/`; `artifact-manifest.json` retains
  raw and packed SHA256 and byte lengths. The schema/Sources are unchanged.
* `gofmt -l research/capacitycomposition` and `git diff --cached --check` —
  PASS, empty diagnostics. Complete new source/test/documents and the generated
  payload were reviewed; exact source identity, graph/grouping counts and matched
  witness prefixes were checked (`self-review.log`). The final staged list
  excludes the controller-owned preflight files and ignored handoff report.

Environment: Go 1.27.4 darwin/arm64 and Staticcheck 2026.2.2 (0.8.2), with exact
version outputs archived. No `go vet` command was invoked. No whole-repository
runtime suite was run; the explicit task forbids SQLite access. Historical
suite reuse, if claimed by the controller, must retain its own verified
source/dependency identity and historical attribution.

## Measured coverage and artifacts

Complete finite abstraction: **72 graphs, 1,632 states, 2,412 arcs, 1,624
CanComplete states, 8 noncompletable states**. Each failing graph has exactly
one noncompletable state: shared-1, early or ready-now, masks 0/1/4/5. All 24
stable-ready graphs have every reachable state completable. Both parenthesizations
agree with the flat relation in all 72 configurations. `graph-totals.log` retains
every row, also present in the raw export.

Factored operational coverage: **311** old durable-ACK edge states, **72**
configurations, **1,854,036** representative cases, **11,972,664** concrete arcs,
**2,451,168** enabled-action lifts and **4,833,504** lift microsteps. Counts
include multiplicity across target selection; they are not unique concrete
product states. The domain crosses every reachable new abstract state, each
target edge's matching old custody states (admitted if held), the other two
owners' independent three-phase availability, and one canonical other-edge
custody record. Both real opaque values occur in the target domain. Every
enabled action is checked, not only the selected target's action.

Genuine witnesses: **10 histories**. Early and ready-now each have an equal-base
AB/BC holder pair and a completed BC-prefix extension; the old-trace exclusion
and a crash-held AB prefix are retained; stable-ready recovery/custody histories
use r1 and r2, actual duplicate requests/results/ACKs, volatile receipt loss,
retained retransmission and recording before ACK release/consumption.

Local hostile tests reject malformed roots/memberships/states, copied shared
pools, inconsistent policies/limits, missing/duplicate owners/originals,
forged permission, invented hold without admission, invented initial receipt,
changed holder, early release, wrong actor/activation, forged result, post-state
invention and duplicate reservation. Full hostile JSON/schema/graph corruption
qualification is Task 2, not claimed here. Test statement coverage percentages
were not measured; the explicit semantic domain counts above are the coverage
reported by this task.

Raw export: 4,514,212 bytes, SHA256
`a7cf316219b30cdc83d455b95a16549c022d55aa422b642b8f5976315e0c0fd2`.
All 32 captured source hashes match the final source files. Deterministic gzip
was verified to decompress to those exact bytes. Packaging does not change the
JSON source manifest or schema.

## Self-review and concerns

Reviewed all new source, tests, argument and schema for the specified invariants;
verified unchanged `contractcomposition`, `twoagent`, `go.mod` and `go.sum`
against the task base. Constructor isolation was checked explicitly: NewConfig
creates fresh old-model slices and never returns the private validation template.
The existing malformed-config test mutates owner generation/policy and original
grant and requires rejection, catching accidental template aliasing. Dynamic holds are never duplicated by composition and
only real Apply clears a held bit. Witness endpoints are captured, not repaired.
Every progress action retains its exact original binding. WorldResult is
preserved for audit but absent from operational decisions. Probe canonical
receipts are clearly labeled synthesized; none appears as a witness initial
receipt. No cap, test condition or domain was weakened to obtain green.

The manuscript's existing §6.16 explicitly assumes sufficient disjoint
resources and excludes shared-resource depletion; this task's new Q relation
is outside that old positive premise. The old claim has not silently been
expanded. Markdown/LaTeX manuscript updates, roadmap, rebuilt PDF, independent
Python audit/hostile controls, and independent review are later controller tasks.
This report does not claim Phase 3C.2 closure or full paper/code reconciliation.

The primary limitation is unchanged: universal lifting rests on the conditional
source argument; probes are not exhaustive concrete-product reachability.
Structural validation does not prove global reachability or receipt authenticity.
Shared-root atomicity, recoverability, fixed identities/topology/policies and
atomic apply/result-retention/release are premises. There is no distributed
allocator, fairness/inevitability, arbitrary topology, cancellation/revocation,
resource-cost equivalence, or old-q unrestricted lifting result.

## File list

New research files:
`research/capacitycomposition/{model,compose,concrete,explore,evidence,probe}.go`,
`research/capacitycomposition/{model,compose,concrete,evidence}_test.go`,
`research/capacitycomposition/cmd/export/{main,main_test}.go`, and package README.

New supporting documents:
`docs/foundational/phase3c2/{01-argument,02-evidence-schema}.md`. This handoff
report remains at its requested ignored local path and is not in the commit.

New task-owned evidence directory:
`docs/foundational/phase3c2/evidence/task1/` contains raw export, source hashes,
five RED logs, initial probe log, normal/race logs, Staticcheck success and
cache-failure logs, environment versions, graph totals and final export/diff
verification logs. The controller-owned `evidence/preflight/` remains outside
this task's commit.
