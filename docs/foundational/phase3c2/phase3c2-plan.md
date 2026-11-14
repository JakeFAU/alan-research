# Phase 3C.2 Shared Capacity Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Test and characterize the loss and recovery of finite completion under shared reusable capacity in the fixed A→B→C chain.

**Architecture:** A new research package wraps the unchanged Phase 3C.1 custody model with original-bound capacity reservations. A separate Python checker reconstructs exported finite graphs and concrete histories. Source arguments, falsifying witnesses and qualification boundaries are reconciled in both manuscript forms.

**Tech Stack:** Existing Go module/toolchain, Python standard library, existing Tectonic/Poppler; no dependency changes.

**Spec:** `docs/foundational/phase3c2/00-contract.md`

## Global Constraints

- Research-only extension; existing runtime and historical model/evidence remain unchanged.
- Fixed generation-one A/B/C chain AB/BC, all eight policy masks, three layouts and three reservation policies from the spec.
- Permission provenance, cumulative expenditure E, and reusable capacity Q are distinct.
- Application, authentic result retention and Q release are atomic; no early release or crash refund.
- Keep all named obligations, full bindings, root identities/memberships and per-owner availability.
- No WorldResult reads in operational guards/lifting; the independent audit does not call Go semantics or runtime Replay.
- CanComplete is existential finite completion; no arbitrary topology, runtime, fairness, split-effect or old-trace-equivalence claim.
- Do not open retained SQLite databases or `.envrc`; do not modify dependencies, production code or old evidence.
- Keep Markdown/LaTeX mathematical content aligned; rebuild and inspect PDF; disclose inherited divergence/warnings.
- Use Staticcheck rather than go vet. Preserve failures and distinguish fresh checks from historical source-bound reuse.

## Review Focus

- A copied root description must not mint capacity; Task 1 rejects aliases with conflicting limits and Task 2 rejects cloned-pool evidence.
- Reservation-time readiness can be invalidated by a later commitment; Task 1 retains that trace and Task 2 reconstructs it.
- Stable-ready removes old legal traces; Tasks 1/2 preserve the control and Task 3 states the proof boundary.
- A recovered holder must not lose capacity ownership or invent a receipt; Tasks 1/2 test concrete recovery and hostile histories.
- Partial graphs can manufacture completion results; Task 2 compares every independently regenerated arc, including non-tree arcs.

---

### Task 1: Model, concrete adapter and frozen export

**Files:**
- Create: `research/capacitycomposition/{model,compose,concrete,explore}.go` and focused `_test.go` files.
- Create: `research/capacitycomposition/cmd/export/{main,main_test}.go` and package README.
- Create: `docs/foundational/phase3c2/{01-argument,02-evidence-schema}.md`.
- Create: `docs/foundational/phase3c2/evidence/task1/` (raw export, source hashes, RED/verification logs).

**Interfaces:**
- Consume unchanged `contractcomposition.NewConfig("chain",mask)`, its types, component functions and concrete transport relation.
- Produce package `capacitycomposition`: `NewConfig(layout, policy string, mask uint8) (Config,error)`, `Initial(Config) State`, `Successors(Config,State) []Step`, `Explore(Config) Graph`, `NewConcrete(Config) (ConcreteState,error)`, `Project(Config,ConcreteState) State`, `ConcreteSuccessors(Config,ConcreteState) []MicroStep`, `Lift(Config,ConcreteState,Action) ([]MicroStep,error)`.
- Produce `Leaf(Config,string) (Component,error)`, `Join(Component,Component) (Component,error)`, `Close(Component) (Config,error)`, validated full root/membership records, and deterministic JSON export via `go run ./research/capacitycomposition/cmd/export`.
- JSON schema version is `alan/phase3c2/v1`; document exact keys/types, roster, witness/probe meanings and CLI output in the schema. Task 2 consumes this schema and checked export, not Go implementation helpers.

- [x] Write failing tests for exact root accounting, aggregate-equal/different-CanComplete holders, ready-now invalidation, stable-ready every-reachable completion, and the excluded old AB-apply trace. Execute and retain RED output before implementation.
- [x] Implement the finite relation and complete BFS/reverse reachability, rejecting malformed configs/states and keeping all arcs. Include shared-two/separate-one and all-servicing controls, two regroupings and invalid joins. No cap values or counts may be silently coerced.
- [x] Implement concrete wrapping and finite-path lifting with all old transport/crash behavior intact. Test every enabled abstract action from declared factored custody/availability representatives; enumerate the representative domain explicitly and validate full projection, legal adjacency and no unrelated progress. Include both real opaque result values, retry/ACK behavior and crash-held witnesses. Report that this is not exhaustive concrete-product enumeration.
- [x] Construct genuine witness prefixes/suffixes from the transition relation for both holders and ready-now invalidation; compare full base states and pool totals. Keep the old-trace restriction witness. Emit full post-states/events and source-bound probe summaries.
- [x] Write conditional source proof and schema. Explain what composition joins, what it cannot join, and why the universal lifting claim refers to the new resource-aware abstraction.
- [x] Run `go test ./research/capacitycomposition/...`, `go test -race ./research/capacitycomposition/...`, `staticcheck ./research/capacitycomposition/...`, deterministic export comparison, formatting/diff checks. Inspect full diff and commit narrowly with imperative subject. Return exact command logs, coverage counts and concerns.

### Task 2: Independent audit and adverse evidence

**Files:**
- Create: `docs/foundational/phase3c2/tools/{audit,test_audit}.py`.
- Create: `docs/foundational/phase3c2/03-independent-audit.md`.
- Create: `docs/foundational/phase3c2/evidence/task2/` (capture tooling, immutable export, exact hostile records/rejections, manifest, logs).

**Interfaces:**
- Consume Task 1's `02-evidence-schema.md`, `01-argument.md`, exported `alan/phase3c2/v1` document and declared CLI.
- Produce `python3 docs/foundational/phase3c2/tools/audit.py EXPORT` with nonzero failure for malformed/incorrect evidence, deterministic summary JSON on success.
- Produce standard-library unittest suite and reproducible hostile generation/capture commands; keep all original evidence and source hashes.

- [x] Implement from the frozen spec/schema a separate declarative graph generator, reverse completion solver, root/grouping validation and concrete trace replay. Do not import or shell out to Go helpers or reuse their guards; independent algorithmic validation is not an independently discovered specification.
- [x] Regenerate all roster graphs/all arcs and compare exact complete labels, parent reachability, CanComplete and counts. Reconstruct holder-equality witnesses and legal ready-now invalidation; verify the stable-ready trace restriction under the stated erased observation.
- [x] Replay actual histories from empty initial custody with actor/availability, exact binding, root membership, reservation/apply/release and result/ACK constraints. Do not infer receipt truth from claimed labels or final states. Validate every provided concrete post-state and projected action.
- [x] Add hostile tests for all spec-listed corruption classes, plus strict schema/duplicate keys/unknown fields/types. Each mutation must be rejected for its intended reason; preserve exact mutated input and output, not only a test name. Demonstrate all-arc checking by removal of a non-tree arc. Do not call a record mutation an operational counterexample.
- [x] Run `python3 -m unittest discover -s docs/foundational/phase3c2/tools -p 'test_*.py' -v`, audit a fresh Go export, and verify retained manifests. Inspect full diff and commit. Return exact evidence paths, measured totals, independence limits and any spec/model disagreements.

### Task 3: Qualification and manuscript reconciliation

**Files:**
- Create: `docs/foundational/phase3c2/{04-results,05-qualification}.md`, `evidence/qualification/`, `evidence/manuscript/`.
- Modify: `docs/foundational/phase3/README.md`, `paper/alan-dissertation.{md,tex,pdf}`; update `paper/references.bib` only if verified citations change.
- Modify: `research/capacitycomposition/README.md` to address Task 1 review M1: copy the verified `-timeout 30m` into the race reproduction command.

**Interfaces:**
- Consume reviewed model/export/checker/proof and exact captured evidence from Tasks 1/2.
- Produce source-bound qualification manifest, result table with measured counts, scoped claims, aligned manuscript text and rebuilt/render-inspected PDF.

- [x] Inspect actual source/evidence and verify export/audit hashes against final source. Reuse this task's reviewed unchanged-source scoped ordinary/race/Staticcheck and audit runs; rerun only for fixes or unresolved concerns. Run root `go build ./...`, `go mod verify`, `go mod tidy -diff`, and formatter/diff checks; record actual commands, versions, exits and outputs. Existing broad lint findings are historical, not silently green.
- [x] Correct the README race command to use the recorded explicit timeout. This documentation-only correction needs comparison to the passing command log, not a new long race run; README is outside the export source-hash set.
- [x] Verify all prior Go/dependency inputs against retained Phase 3C.1 manifests and merged base. Attribute its unchanged full-suite evidence explicitly instead of claiming a new full-suite run. Any changed old input invalidates reuse and requires appropriate fresh checks. Reconcile `go list ./...` package coverage with inherited plus new packages; no hidden missing packages.
- [x] Write results distinguishing elementary resource-wait counterexample, source proof, finite checked roster, factored concrete probes and selected histories. Explain stable-ready's cost and the new-versus-old abstraction boundary. Preserve open obligations and current runtime distinction. Cite primary background accurately if needed; do not invent metadata or novelty.
- [x] Update both manuscript forms and roadmap with the same claims, assumptions, counts, controls and exclusions; preserve prior theorem scope. Use the PDF skill and existing build tools to rebuild, extract text and render/inspect changed pages. Record inherited differences/warnings and check for new ones. No wholesale manuscript repair.
- [x] Archive source-bound manifests/logs without temporary artifacts or secrets. Inspect complete diff, run applicable final checks and commit narrowly. Return exact fresh/historical qualification split, PDF evidence and limitations.
