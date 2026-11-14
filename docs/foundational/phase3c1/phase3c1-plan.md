# Phase 3C.1 Contract Preservation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Test whether a named progress contract preserves original obligations and existential completion when a verified pair is composed with a third agent.

**Architecture:** A small Go progress model and explicit compatible joins are related by a finite-path argument to an adapter over unchanged durable-ACK transport semantics. Independent Python graph/history validation challenges exports; preserved hiding counterexamples and manuscript reconciliation bound the conclusion.

**Tech Stack:** Existing Go module, Python standard library, cached Tectonic/Poppler; no new dependencies.

**Spec:** `docs/foundational/phase3c1/00-contract.md`.

## Global Constraints

- Research-model changes only; no production/runtime changes, new dependencies or mutation of old models/evidence.
- Only the fixed chain/cycle grammar, all eight service masks, fixed generation/authority, disjoint resources and atomic application/receipt are qualified.
- Preserve all original debt, endpoint policy incidence and visible application labels through grouping; abstract enablement is not concrete immediate readiness.
- Distinguish universal source arguments, complete small abstract checks, factored edge checks, concrete witnesses and unestablished runtime/distributed claims.
- Preserve every failed candidate and reviewer REJECT; no weakening evidence to obtain acceptance.
- Never open retained SQLite or read/change runtime `.envrc`; use `/tmp/alan-phase3c1/` for scratch.
- One implementation worker at a time, no worker subagents; separate task reviews and one final whole-branch review; push a draft PR and do not merge.

## Review Focus

- A grouped summary drops an internal original or substitutes a single service policy; genuine R0/R1 witnesses must expose it.
- An open port is treated as an idle callee, or overlapping owners/aliased resources are accepted; malformed compositions must reject.
- The lifting argument uses a favorable representative or remote/oracle knowledge; check arbitrary old edge states and actual custody, including loss before recording and after ACK/release.
- The macro DAG is confused with inevitable progress or physical readiness; retain retry divergence and down/unadmitted distinctions.
- Go and Python agree only because they share a transition implementation or trust reported counts; rebuild independent expected relations and attack the exported evidence.

---

### Task 1: Implement progress contracts, concrete lifting and candidate counterexamples

**Files:** Create `research/contractcomposition/{model,compose,concrete,explore}.go`, corresponding `_test.go` files, `README.md`, `cmd/export/{main.go,main_test.go}`; create `docs/foundational/phase3c1/01-argument.md`, `02-evidence-schema.md`, `evidence/task1/`.

**Interfaces:** Consume the contract and advisory `/tmp/alan-phase3c1/semantics-recon.md`, `tooling-recon.md`. Produce deterministic `go run ./research/contractcomposition/cmd/export` JSON on stdout and a documented versioned schema for Task 2. Freeze exact named config/state/action/graph/witness/claim fields in `02-evidence-schema.md`; include explicit identities and complete graph arcs. The concrete adapter may reuse only unchanged `research/twoagent` APIs; all other new behavior lives in the new package.

- [ ] Define `Config`, `State`, `Action`, `Component` and concrete adapter records with constructor validation. Use U/W/P/D stages, all eight masks and exact named bindings. Write failing semantic tests before behavior: hidden-debt completion collision, hidden-wait readiness collision, closed-port/owner/resource rejection, grouped equality and serial-cycle progress obstruction. Preserve RED logs.
- [ ] Implement pure progress successors, Complete, named service guards, compatible leaf/join/close, full small-graph BFS and reverse CanComplete. Export flat/left/right comparisons, all configs, deterministic witnesses and dependency-certificate checks. Counts are computed, never prefilled from expectations.
- [ ] Implement concrete transport transition/projection and finite-step lifting using actual retained data and legal transitions. Test every old durable-ACK reachable edge state against relevant phase/policy/availability cases; check old independent trace ledger on projected witnesses, unchanged world-oracle noninterference, pre-record loss, ACK/release custody, no reapplication, exact bindings and concrete retry stutter. Record precise coverage; do not claim full microstate enumeration.
- [ ] Write the transition-by-transition projection/lifting argument and limits. Keep genuine R0/R1 histories and the all-serial compatibility control. Freeze the evidence schema and sample export for the independent auditor.
- [ ] Run `go test -count=1 ./research/contractcomposition/...`, `go test -race -count=1 ./research/contractcomposition/...`, scoped Staticcheck and formatting. Inspect full diff/status, commit assigned paths only, and write the task report with exact commands, failures, claims and gaps. Coordinate index use with controller.

### Task 2: Independently audit the exported graphs and concrete witnesses

**Files:** Create `docs/foundational/phase3c1/tools/{audit.py,test_audit.py}`, `03-independent-audit.md`, `evidence/task2/`; consume frozen Task 1 schema/exports without changing them.

**Interfaces:** Read `00-contract.md` and `02-evidence-schema.md`; independently encode Python successor/guard, compatibility and concrete ledger checks without importing/calling Go reducers. `python3 docs/foundational/phase3c1/tools/audit.py INPUT.json` prints a deterministic audit result and exits nonzero on rejection. Tests run via `python3 -m unittest discover -s docs/foundational/phase3c1/tools -p 'test_*.py'`.

- [ ] Write failing tests for each named hostile control before the validator accepts/rejects it. Validate schema, exact config roster and IDs/bindings, reachability from initial state, the complete expected transition relation including incoming arcs, independently computed Complete/CanComplete, both grouped constructions, dependency claims and all reported counts.
- [ ] Replay all concrete witness events independently from initial state, enforce shared availability and receipt custody, compare reported states/projections and recompute R0/R1 predicates. Verify retry-cycle and assumption classification. Reject forged binding, missing/extra arcs, dropped hidden debt, wrong guard/progress verdict, premature release and repeated application for specific reasons. Genuine records remain unchanged.
- [ ] Run fresh export and audit; retain compressed full graphs, command logs, raw digest manifest and hostile rejection report. Explicitly distinguish abstract exhaustiveness, factored Go checks and selected concrete replay. Any substantiated model defect returns to its owner and must be fixed/re-reviewed; do not compensate in the oracle.
- [ ] Run Python tests and inspect all changes; commit only assigned paths and write a report for independent task review.

### Task 3: Reconcile the result, qualification and manuscript

**Files:** Create `docs/foundational/phase3c1/{04-results,05-qualification,06-reviews}.md`, evidence/manuscript and qualification records; refine `01-argument.md` only as required by actual results; update `docs/foundational/phase3/README.md`, `paper/alan-dissertation.{md,tex,pdf}`, `paper/references.bib` if cited. Controller owns final workflow/delivery metadata.

**Interfaces:** Consume reviewed Go model, independent audit and controller baseline records. State conditional mathematical arguments separately from measured finite results, source inspection and concrete witness checks; neither reviewers nor counts prove a broader theorem.

- [ ] Challenge every positive proof step against actual code/evidence; narrow any incomplete claim. State genuine hiding failure versus violated assumptions, exact observation relation, grouping limits, the acyclic sufficient condition's quantifiers and remaining Phase 3C obligations.
- [ ] Reconcile manuscript and roadmap without changing historical results. Check the primary prior-art source before adding a precise BibTeX citation. Read PDF skill before this new operation; build from scratch copies, inspect changed pages and continuations, retain source/log/render hashes and inherited discrepancies.
- [ ] Record the already-run ordinary suite over the frozen old packages plus fresh added-package tests; verify total package coverage and no old executable changes. Run root `go build ./...`, `go mod verify`, `go mod tidy -diff`, and scoped Staticcheck; compare any global lint outputs with the inherited 13/31 diagnostics. No unnecessary runtime campaign or repeated full suite.
- [ ] Verify artifact hashes, local links, schema/audit agreement, paper synchronization, complete scope and whitespace; inspect full diff and commit assigned paths. Report exact results and limitations for task review.
- [ ] Controller obtains final whole-branch review, addresses findings with scoped re-review, archives exact workflow, pushes/attaches draft PR and verifies delivery. No merge.
