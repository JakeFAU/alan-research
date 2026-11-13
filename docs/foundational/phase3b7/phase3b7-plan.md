# Phase 3B.7 Matched Capture Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Produce and independently check a genuine zero/one application pair with equal declared recovered local observations, or retain the exact falsifying difference.

**Architecture:** A test-only overlay runs the frozen runtime from a disposable local checkout. A common closed seed, fixed logical URLs and separate worker/scheduler barriers produce two isolated histories. An opt-in raw-data checker validates protocol and full local equality independently; the existing finite model remains a separately identified projection.

**Tech Stack:** Existing Go module/toolchain and Staticcheck, Python standard library for custody tooling, existing Tectonic/Poppler manuscript tools. No new dependencies.

**Spec:** `docs/foundational/phase3b7/01-contract.md` and repository `AGENTS.md`.

## Global Constraints

- Deliver an independently reviewed draft research PR; do not merge or release.
- No production runtime, dependencies, automatic retry, or authority semantics change.
- Existing Phase 3B.6 artifacts and interpretations remain intact. Old checkers' pinned producer protocols remain unchanged.
- No later pass rewrites raw artifacts to create equality. Opening the archived evidence SQLite files is prohibited.
- No identity renaming, field deletion, table filtering or post-hoc reordering is permitted. Array order, multiplicity, numbers, strings and null/empty distinctions remain significant.
- Freeze producer/comparison sources and contract hashes before the qualification pairs. Preserve every attempted qualification run, including mismatch or apparatus failure.
- Work in the existing research worktree on `codex/phase3b7-matched-capture`, base `d8645d739445d4b09c0d38dca4454ca6297fdf22`; do not edit the runtime checkout or `.envrc`.
- Use `/tmp/alan-phase3b7/` for disposable runtime builds and tests. Offline Go environment: `GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b7/go-build`.
- One implementation agent at a time; no worker subagents. Each task receives independent spec and quality review. Preserve actual commands, exits, source identities and failed attempts.

## Review Focus

- A worker barrier alone leaves scheduler progress uncontrolled; both must remain held at the declared common cut (Task 1).
- Copying only SQLite's main file while WAL is live creates a false seed/capture; use complete closed file sets and verify bytes (Task 1).
- Identical small model projections can hide different full local state; compare complete exports without pruning (Task 2).
- An empty or final-only endpoint artifact does not certify zero at an earlier cut; use process-bound coverage and causal holds (Task 2).
- Both-pair success is bounded evidence, and any failed attempt must survive qualification and manuscript reconciliation (Task 3).

---

### Task 1: Implement the frozen matched-capture producer

**Files:** Create `research/matchedcapture/producer/phase3b7_matched_test.go.txt`, `research/matchedcapture/producer/prepare.py`, `docs/foundational/phase3b7/02-protocol.md`, `evidence/task1/`. Additional focused test-only owner/helper source may be split in the same producer directory when responsibilities warrant it.

**Interfaces:** A preparation command consumes the read-only local runtime Git repository and a new temporary destination, checks the frozen commit/tree, creates a detached local checkout without touching the original, copies the named overlay into `internal/durabletest`, and emits a full producer source identity. New integration entry `TestPhase3B7MatchedCapture` takes `ALAN_PHASE3B7_EVIDENCE` as a fresh output root and emits one complete pair with a versioned `pair.json`, branch protocol records, seed/captures, endpoint dossiers and complete member manifest. `TestPhase3B7Owner` is its owner child. Both normal and race test binaries use identical producer sources. Document precise schema for Task 2; do not fabricate compatibility with old producer anchors.

- [x] Read the contract and `/tmp/alan-phase3b7/capture-recon.md`; inspect the named frozen source. Build the minimal two-barrier owner, retaining fsynced event/control records and exact in-flight identity. Use existing fixture/revoke/inspection helpers without modifying frozen production sources. Fixed-URL endpoints run sequentially with fresh stores and complete waits.
- [x] Implement common closed-seed initialization and full-file copying, pre-call and scheduler handshake, zero hold/applied release, acknowledged SIGKILL, exact renewal/revoke/refuse tail and immutable capture sealing. Record failure evidence before assertions abort; cleanup must terminate/wait without releasing a zero worker into a call. No sleeps as causal synchronization.
- [x] Add focused test-only barrier controls exercising unreleased/released/EOF/context-canceled behavior and identity mismatches. Demonstrate meaningful RED/GREEN evidence for the new fail-closed mechanism before declaring it qualified; do not claim a missing symbol compilation failure as semantic RED.
- [x] Prepare/build/test only the harness controls in the disposable checkout, format overlay sources and run scoped Staticcheck where supported. Do not run `TestPhase3B7MatchedCapture` yet: the checker and contract must be frozen first. Inspect full diff, commit narrowly and report source/interface details, tests, commands and any unresolved risks.

### Task 2: Implement independent raw comparison and causal validation

**Files:** Create `research/matchedcapture/check.go` (report/orchestration), `protocol.go` (process chronology), `compare.go` (strict complete JSON equality), `store.go` (raw binding/renewal/model projection), `world.go` (endpoint coverage), `source.go` (source/member seals), focused `*_test.go`, `cmd/check/main.go` and tests; `docs/foundational/phase3b7/03-comparison.md`, `evidence/task2/`. Shared wire types may live in `wire.go`; keep responsibilities focused. Add small adapters to the existing correspondence package only if required to expose unchanged validation, with old output byte regressions.

**Interfaces:** `CheckDirectory(dir string) (Report, error)` reads the pair directory without opening databases. The opt-in command accepts `-pair DIR` and prints deterministic JSON; errors exit nonzero. The Report contains independent seal/source/raw/protocol/world/equality/model verdicts, exact input digests and deterministic JSON paths for differences, with no inferred success from case names. Consume Task 1's actual schema documented in its report. Full-export equality covers seed/post-open/stopped/reopened/final; final report equality is separate. Model mapping uses the existing Phase 3B.6 `StoreView` domain and independent ledger for each supported cut.

- [x] Read contract and producer protocol. Write synthetic fixture tests with meaningful RED before implementation for full-field mismatch despite matching small projection, duplicate-key rejection, null/empty distinction and same-applied-pair rejection. No genuine experiment needed for this stage.
- [x] Implement strict complete artifact/member/source checks, canonical JSON comparison preserving all array order, and precise path reporting. Validate the complete selected raw binding and original allowance/debit, process/source/control joins, both-epoch exact renewal, revoke/refusal/c survival and absent authentic custody. Keep every verdict separate; unsupported data fails closed.
- [x] Validate zero through uninterrupted pre-call hold/death and complete endpoint coverage, one through exact post-commit/death/release ordering, and no cross-branch process/store contamination. Add coherently rehashed attacks for every contract-listed falsifier and lawful positive controls. Use raw JSON/JSONL and independent comparison logic, never producer summary verdicts.
- [x] Run focused tests, command tests, scoped race tests and Staticcheck. Preserve unchanged historical report exports when shared adapters are necessary. Inspect full diff, commit narrowly, and report schema/freeze instructions for qualification.

### Task 3: Capture, qualify and reconcile the result

**Files:** Create `docs/foundational/phase3b7/04-runbook.md`, `05-results.md`, `06-qualification.md`, `07-reviews.md`, `evidence/capture/`, `evidence/qualification/`, `evidence/workflow/`; update `paper/alan-dissertation.md`, `.tex`, rebuild `.pdf`. Bibliography only if references change.

**Interfaces:** Controller freezes reviewed Tasks 1–2 and the contract, then executes one normal pair and one race pair. Worker consumes actual frozen raw evidence and independent reports to write the results/manuscript. Controller owns broad verification, archive seals, final whole-branch review and delivery.

- [ ] Controller records exact source/toolchain/commands, builds both producer binaries and runs each prescribed pair once. Seal/archive all raw output and failure evidence before independent checking. No source or raw rewrites during an attempt. Run the independent command on each pair twice and compare deterministic reports. If capture fails, preserve the failed version and determine whether it is apparatus or semantic failure before any revised protocol.
- [ ] Controller performs affected Go tests/race/build/format/Staticcheck, module checks, baseline regressions and source/archive/member checks. The unchanged full repository's existing Phase 3B.6 qualification is attributed as historical; choose any fresh broad tests by dependency impact and disclose scope. Re-run checks only for changed code or unresolved findings.
- [ ] Worker writes measured tables and limitations directly from the actual reports. Reconcile manuscript body and conclusions, clearly distinguishing concrete pair, finite model, source-conditioned protocol argument and general open claims. Preserve every mismatch/failed attempt. Read the PDF skill before the new artifact operation; build and visually inspect changed pages/continuations and record inherited divergence/warnings.
- [ ] Obtain task spec/quality review and one final broad whole-branch review. Correct substantiated findings and obtain explicit re-review. Archive workflow evidence/rulings, inspect final status/stat/full diff, push and create/attach draft PR against current main. Do not merge. Final report gives outcome, evidence, limitations, PR and material rulings/costs.
