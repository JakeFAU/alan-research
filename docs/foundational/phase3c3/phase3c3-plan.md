# Phase 3C.3 Component Interfaces Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Falsify snapshot-only substitution and characterize a successful-boundary-language interface for finite continuation monitors.

**Architecture:** Reuse the unchanged Phase 3C.2 plant and concrete lift through a new research wrapper. A finite automaton layer compares marked languages exactly and composes role-checked continuation monitors. An independent Python checker reconstructs the frozen domain and genuine model histories.

**Tech Stack:** Existing Go module and standard library; Python standard library; existing Tectonic/Poppler manuscript tooling. No new dependencies.

**Spec:** [00-contract.md](00-contract.md).

## Global Constraints

- New code is confined to `research/componentinterfaces/`; records and Python auditor to `docs/foundational/phase3c3/`.
- No production code, dependencies, old models or historical evidence changes.
- The unchanged Phase 3C.2 relation K retains the full A/B/C service, authority and capacity semantics. This is continuation-monitor substitution, not replacement of an entire agent.
- All eight old named commit/reserve-capacity/apply/consume labels and their exact bindings remain visible. Only private bookkeeping tau is hidden.
- The five M fixtures and 324 two-state A/C controllers are exactly those in the contract. Every state of the five M fixtures is marked.
- Complete means both originals D and all factor marks; CanComplete is existential finite completion. It is distinct from every-reachable-state nonblocking.
- Product signatures retain one shared Q relation; monitors cannot change authority, stages, holds or external effects independently.
- Preserve exact raw captures, rejected inputs, source hashes and prior failures. Passing bounded checks do not prove the general source argument.
- Use scoped ordinary/race Go checks and Staticcheck; Python independent reconstruction; fresh build/module checks. Historical qualification may be reused only with exact source identity and explicit attribution.
- Reconcile changed Markdown/LaTeX mathematical content, bibliography implications and rebuilt PDF. General Phase 3C remains open.

## Review Focus

- Tau cycles and empty visible words can defeat a finite-prefix comparison: Task 1 exact equivalence tests and Task 2 independent closure/pair search.
- All-marked monitors can preserve initial language yet introduce a bad residual: Task 1 doomed-branch test; Task 2 complete product reverse reachability.
- A role-invalid monitor could veto someone else's action: Tasks 1/2 reject missing external self-loops and mismatched signatures.
- Loss of a non-tree arc can preserve BFS parents and superficially valid counts: Task 2 retains a hostile input correcting those counts and requires rejection.
- Model histories may accidentally reset private state on recovery or forge equality: Tasks 1/2 replay from empty custody, verify common cuts and retained doomed state through interruption.

---

### Task 1: Boundary automata, substitution attack and model captures

**Files:**
- Create `research/componentinterfaces/automaton.go`, `language.go`, `product.go`, `plant.go`, `concrete.go`, `evidence.go` and focused `_test.go` files; a small `fixtures.go` or `capture.go` split is permitted if responsibilities require it.
- Create `research/componentinterfaces/cmd/export/main.go` and `README.md`.
- Create `docs/foundational/phase3c3/01-argument.md`, `02-evidence-schema.md`, and `evidence/task1/`.

**Interfaces:**
- Consumes unchanged `capacitycomposition.NewConfig`, `Initial`, `Successors`, `Complete`, `NewConcrete`, `ConcreteSuccessors`, `Project`, `Lift`; `contractcomposition.Action/Binding` supplies exact visible bindings.
- Produces `Export(w io.Writer) error`, called by `go run ./research/componentinterfaces/cmd/export`; stdout is deterministic versioned JSON, diagnostics go to stderr. No timestamps or temporary paths in semantic export.
- Produces a complete versioned schema in `02-evidence-schema.md` for Task 2: exact alphabet/signatures, five monitor descriptors, 72 plant descriptors, 324 D descriptors, 1,980 full product graphs, all ten language comparisons, normalized grouping checks, eight concrete histories with event/post-state/cut metadata, counts, and producer source manifest. Full declared arcs and markings are retained; no counts-only certificate.
- Internal Go types/function names are the implementer's responsibility; keep them focused and expose only the exporter boundary needed downstream.

- [x] Read the spec, owning old functions, relevant Go/testing skills and write focused failing tests before implementation. Retain raw RED command/output/exit; distinguish compile failures from falsified semantic hypotheses.
- [x] Implement strict automaton/signature/role validation, shared-visible/private-tau product, complete BFS/reverse CanComplete and every-reachable nonblocking classification. Keep source-state keys stable and grouping tuples explicitly normalizable.
- [x] Implement exact epsilon-NFA language comparison by closure/subset-pair exploration, including shortest distinguishing words and acceptance checks. Tests cover epsilon acceptance mismatch, tau cycle, empty subset, nondeterministic equivalent branches and changed post-prefix behavior beyond a shallow sample.
- [x] Construct the five exact M fixtures and 324 exact D fixtures. Enumerate the contract's 360 + 1,620 full products and all ten language comparisons. Check both parenthesizations against the flat triple relation for the 1,620 D cases. Retain both successful and blocked equal-language substitutions.
- [x] Implement the concrete M wrapper and eight genuine captured histories, using old custody preparation and real successor steps. Test common-cut full old state/holds/config equality with opposite product CanComplete; positive suffix; shared-two/separate/servicing controls; private B-up requirement; retained PC after crash/recovery; successful and doomed branches. No post-capture state rewriting.
- [x] Write the source argument: exact product-language intersection in both directions; conditional substitution for pointed components; ALAN instantiation and concrete projection/lifting premises; source proof versus finite evidence; nonblocking counterexample; no necessity/full abstraction for the role-restricted grammar. Preserve the stable-ready restriction boundary.
- [x] Freeze schema and capture the complete deterministic export twice from final producer source. Bind all actual producer inputs, including imported old model files and relevant Go module files. Record actual measured totals; do not invent expected graph sizes.
- [x] Run `go test -v -timeout 30m ./research/componentinterfaces/...`, the same command with `-race`, and `staticcheck ./research/componentinterfaces/...` with writable task caches. Capture exact argv/environment/exit/output. After checks pass, do not rerun unchanged suites merely for repetition.
- [x] Inspect full task diff and status, commit only this task, and write the task report with semantics, evidence, counts, limitations and all concerns. Leave root-wide checks/manuscript to controller/Task 3.

### Task 2: Independent interface and history audit

**Files:**
- Create `docs/foundational/phase3c3/tools/audit.py`, `test_audit.py`, `03-independent-audit.md`, and `evidence/task2/`.

**Interfaces:**
- Consumes Task 1's frozen `02-evidence-schema.md` and exporter stdout, not Go helpers or operational predicates.
- Produces CLI `python3 docs/foundational/phase3c3/tools/audit.py PATH` accepting plain JSON or gzip. Valid complete export exits 0; malformed/false evidence exits nonzero with a precise reason.
- Produces exact captured accepted/rejected inputs, command records, source/archive manifest and an audit report for Task 3. Fresh Go export may be invoked to capture bytes but cannot serve as the independent transition oracle.

- [x] Read contract/schema first; implement a separate declarative plant, fixture generator, complete synchronous product, epsilon-language comparison, marking/reverse reachability and concrete replay. Do not import, translate mechanically from, or call Go guards; use the shared written semantics. Reject duplicate/unknown records and invalid types, including Python bool-as-int accidents.
- [x] Independently regenerate all 1,980 graphs with all arcs, marks and CanComplete values, ten comparisons, 324 descriptors and normalized parenthesizations. Check distinguishing words by independent acceptance and shortestness as the schema promises. Validate equal-language substitution in the entire declared roster.
- [x] Replay all eight concrete histories from genuine empty initial custody. Recompute every post-state, actual-result provenance, monitor step/role, physical eligibility and cut equality. Tie concrete projections to product states and verify positive suffixes/negative graph labels. No hidden Q changes, result invention or PC reset.
- [x] Write positive and hostile tests that cover every Review Focus and all hostile categories in the contract. At least one omitted non-tree-arc attack must retain reachability parents and correct superficial counts. Include forged equivalence, tau/mark changes, external veto, signature/binding/root changes, fabricated completion and altered common-cut/recovery evidence. Keep rejected fixtures exact and reasons meaningful; preserve any initially accepted attack before fixing checker.
- [x] Capture a fresh producer export and confirm source-bound byte equality with Task 1. Run `python3 -m unittest discover -s docs/foundational/phase3c3/tools -p 'test_*.py' -v` plus valid plain/gzip and every hostile CLI input, retaining argv, outputs and exit status. Hash every accepted/rejected input and checker source.
- [x] Write audit scope and measured findings. Explicitly separate independent abstract reconstruction and concrete replay from the conditional source proof and any unexecuted runtime claims. Inspect full diff and commit only Task 2; report exact paths/commits/checks and concerns.

### Task 3: Qualification, manuscript and bounded research result

**Files:**
- Create `docs/foundational/phase3c3/04-results.md`, `05-qualification.md`, `evidence/qualification/`, `evidence/manuscript/` and an artifact verifier.
- Modify `docs/foundational/phase3/README.md`, `paper/alan-dissertation.md`, `paper/alan-dissertation.tex`, `paper/references.bib` as warranted, and rebuild `paper/alan-dissertation.pdf`.

**Interfaces:**
- Consumes accepted Task 1/2 reports/schema/captures and controller handoff `/tmp/alan-phase3c3/task3-handoff.md`, including `/tmp/alan-phase3c3/baseline-inputs.json` and freshly captured build/module/package checks.
- Produces source-bound qualification and manuscript/PDF records, verified alignment and scoped result. Reviewer status/delivery metadata remain controller-owned.

- [x] Read the accepted evidence and distinguish the negative snapshot result, conditional initial-CanComplete theorem, concrete-model correspondence and doomed-residual boundary. Check every manuscript implication against owning source rather than treating old prose as authority.
- [x] Archive controller checks with exact bytes and manifests. Verify all 669 captured merged-baseline Go/module inputs unchanged; link source-identical historical test/lint coverage separately from fresh scoped evidence. Check current package roster is inherited 74 plus new packages, explaining any difference. No fresh whole-repo suite/global lint claim.
- [x] Write concise results/qualification with actual counts, unchanged-source identities, original RED and hostile evidence, no full-agent/runtime/optimality/verification-cost claims. Clarify that ordinary monitors are all marked and additional continuation restrictions invalidate importing the old stable-ready progress theorem.
- [x] Read/apply the PDF skill and coordinate its one-time operation marker with controller before source authoring. Align new Markdown/LaTeX mathematical content and changed Abstract/Conclusion. Use the controller's inspected Tang/Moor primary source as background, with exact metadata; do not imply a new general automata theorem. Preserve inherited manuscript differences explicitly.
- [x] Build baseline/current PDF with existing cached tooling, inspect changed pages and continuations visually, compare warnings, extract text, and record hashes. No unrelated global paper repair or normalization. Verify citations and local links in changed content.
- [x] Provide a deterministic artifact/source/alignment verifier and execute it. Inspect complete task diff, confirm only intended files/evidence/PDF are included, commit Task 3, and report exact checks, remaining limitations and review pointers.
