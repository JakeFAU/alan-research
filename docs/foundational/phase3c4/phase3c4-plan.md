# Phase 3C.4 Residual Interfaces Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Define and challenge a finite behavioral interface preserving contextual nonblocking, with its strict limitations retained.

**Architecture:** A new Go package consumes hash-pinned predecessor automata and decides residual-language profile equality. It reconstructs complete products and captures counterexamples and a private-cycle escape. A separate Python checker reconstructs these relations from the written contract and source-bound predecessor Python constructors. Manuscript qualification includes the user's two editorial corrections.

**Tech Stack:** Existing Go module/standard library, Python standard library, cached Tectonic and Poppler. No new dependency.

**Spec:** [00-contract.md](00-contract.md).

## Global Constraints

- This packet replaces finite B continuation monitors over the unchanged whole Phase 3C.2 plant K. It does not extract an entire independently executing agent.
- New Go code is confined to `research/residualinterfaces/` and its exporter.
- Do not modify old model, runtime, module, test or evidence files.
- No event priorities, fairness, inevitability or divergence preservation are assumed. Shared authority, effects and Q remain in the one unchanged K.
- Monitors restrict only B's reserve/apply AB and commit BC; all other visible labels self-loop without changing private control.
- Nonblocking means all reachable nodes have CanComplete=true.
- This is a conditional source argument, not a mechanized proof or a new general process-theory theorem.
- Do not repeat unchanged suites merely for repetition.
- No new concrete histories, endpoint/runtime campaign, SQLite opening or old exhaustive transport-probe execution is required or claimed.
- General Phase 3C, full-agent substitution, dynamic authority, distributed ownership and split effects remain open.

## Review Focus

- Empty weak reach set versus a reached unmarked dead residual: Tasks 1/2 distinguish empty profile from a profile containing the empty language, including a mismatch with no possible nonblocking separator.
- Multiple reached states with equal residual language: Tasks 1/2 use bidirectional coverage, allowing one-to-many matches; no node bijection or equal state counts.
- Private cycles with an enabled completion escape: Tasks 1/2 retain the actual lasso and suffix; Task 3 distinguishes nonblocking from inevitability and fairness.
- A profile separator may need different suffixes against each opposite residual, and its prefix need not be executable in K: Tasks 1/2 preserve individual language comparisons and real product paths separately.
- Omitted non-tree arcs, altered source/signature identities, unknown JSON fields and bool-as-int values can fake a valid result: Task 2 retains exact hostile inputs and complete reconstruction; Task 3 checks archive/source binding.

---

### Task 1: Residual interface model, argument and complete evidence

**Files:**
- Create `research/residualinterfaces/` with focused Go files for pinned input, finite automata/languages/profiles, products, witnesses and evidence; matching focused tests, `cmd/export/main.go` and `README.md`.
- Create `docs/foundational/phase3c4/01-argument.md`, `02-evidence-schema.md`, and `evidence/task1/`.

**Interfaces:**
- Consumes the exact pinned Phase 3C.3 capture and source manifest in the contract; imported descriptor syntax is the predecessor schema in `docs/foundational/phase3c3/02-evidence-schema.md`.
- Produces `Export(w io.Writer) error`, deterministic JSON stdout via `go run ./research/residualinterfaces/cmd/export` from repository root. No timestamp/temp path in semantic output.
- Produces a complete versioned schema before Task 2: input/source bindings, seven monitor descriptors and exact inherited context descriptors, all 21 profile comparisons with complete cross-state language certificates and weak-subset graphs, all 2,772 complete product graphs, contextual nonblocking substitution rows, negative core comparison, actual private-cycle prefix/cycle/escape and separate generic strictness control. Every graph includes all reachable states, arcs, marks and reverse-reachability labels; no counts-only certificates. Declare canonical ordering and shortest-witness tie-breaking.
- Internal names/types are builder-owned; do not export old-package internals or modify old source to simplify reuse.

- [ ] Read the contract, predecessor schema/constructors, relevant Go skills and tests. Pin compressed/raw predecessor identity, exact 27 predecessor source hashes and all compiled/hash-only producer inputs. Retain initial failing test command/output/exits separately from semantic falsification.
- [ ] Implement generic finite epsilon-NFA acceptance and exact residual language equality. Tests cover empty-word marking, tau closure/cycles, empty subsets, nondeterministic equivalent branches and a distinguishing suffix beyond a shallow sample.
- [ ] Implement complete weak-subset-pair profile comparison with bidirectional residual-language coverage. Tests require RF/tau/duplicate/loop profile equality; RF/doomed language equality but profile inequality; different state multiplicities accepted; empty profile versus containing-empty-language rejected. Retain shortest mismatch prefix and per-opposite-residual separators, including vacuous opposite sets.
- [ ] Instantiate the exact seven M fixtures and inherited 72 K/324 D roster. Rebuild all 2,772 products from automata, including all arcs and marks; compute CanComplete by reverse reachability. Record every profile-equal pair across all 396 contexts with both nonblocking labels and keep blocked cases. Test exact role, signature, root and visible-effect preservation.
- [ ] Capture genuine abstract paths for the core RF/doomed negative and loop prefix/cycle/completion escape. Assert loop cycle states carry unpaid debt, are unmarked and have a completion suffix; do not call the lasso a crash/runtime capture. Retain the generic unmarked strictness example outside the ALAN roster. Test conditional profile closure under product on small independent automata, including private steps and marks.
- [ ] Write the derivative-union, Cartesian reachability, residual intersection, nonblocking criterion and product-profile closure arguments, with both directions and explicit assumptions. Distinguish generic strictness from the all-marked ALAN domain. Preserve open Theorem 4.8, lack of weakest-interface/cost claims and historical concrete evidence attribution.
- [ ] Freeze argument/schema and capture complete new export twice from final producer sources, with exact bytes/manifests and actual measured totals. Retain deliberate failed hypotheses (union sufficiency; nonblocking implies inevitable completion) as clearly labeled falsification evidence, separate from passing counterexample assertions.
- [ ] Run `go test -v -timeout 30m ./research/residualinterfaces/...`, the same with `-race`, and `staticcheck ./research/residualinterfaces/...` using writable task caches. Verify actual compiled local producer inputs with `go list -deps -json`; preserve any diagnostics and resolved executions accurately.
- [ ] Inspect complete task diff/status and commit only Task 1. Report exact commands/results, paths, source/capture hashes, counts and concerns. Leave root-wide qualification and manuscript to controller/Task 3.

### Task 2: Independent profile, product and witness audit

**Files:**
- Create `docs/foundational/phase3c4/tools/audit.py`, `test_audit.py`, `03-independent-audit.md` and `evidence/task2/`.

**Interfaces:**
- Consumes frozen contract/schema and Task 1 exact export, not Go guards or verdicts.
- May import unchanged `docs/foundational/phase3c3/tools/audit.py` constructors with explicit attribution and source hashes for it and its Python dependencies. New profile/language/reachability/witness logic is independent of Go; old product/outcome fields are not an oracle.
- Produces `python3 docs/foundational/phase3c4/tools/audit.py PATH`, accepting exact plain JSON/gzip. Valid complete evidence exits 0; malformed/false evidence exits nonzero with a specific reason. Produces source/archive manifests, exact accepted/rejected CLI inputs and command records.

- [ ] Read contract/schema before implementation. Reconstruct all inherited/new descriptors from written definitions and source-bound Python constructors; check pinned predecessor compressed/raw bytes and source manifest. Implement strict types/keys, duplicate detection and finite schema normalization without bool/int equivalence.
- [ ] Independently reconstruct all cross-state language decisions, complete profile subset graphs, shortest mismatch prefixes and every unmatched-state separator. Reconstruct all 2,772 product graphs, every CanComplete/nonblocking label and contextual substitution row. Verify true and false contexts without treating a mismatch as a necessary contextual separator.
- [ ] Independently validate each negative path, actual loop prefix/cycle/escape and generic strictness example. Verify cycle closure, unpaid debt, marks and complete escape. No post-capture editing, no new physical-history claim.
- [ ] Write positive and hostile tests for all Review Focus items and contract attack categories. Include repaired-count/missing-non-tree-edge input, union-for-profile forgery, missing coverage/separator, bad shortest prefix, dead/unreachable confusion, forged marks/tau, external veto, identity/root mismatch, invented nonblocking and invalid lasso/escape. Retain JSON duplicate-key/nonfinite/type controls. Preserve every initially accepted attack before repair.
- [ ] Capture a third fresh producer export and compare exact bytes/source identities. Run `python3 -m unittest discover -s docs/foundational/phase3c4/tools -p 'test_*.py' -v`, valid plain/gzip audit and every exact hostile CLI input. Archive argv/environment/output/exits and hash all inputs/checker dependencies. Do not replace CLI checks by test-only assertions.
- [ ] Report independent versus reused semantics, actual counts, falsifications and limits. Inspect full task diff, commit only Task 2, and provide the task report for review/qualification.

### Task 3: Qualification, manuscript reconciliation and requested corrections

**Files:**
- Create `docs/foundational/phase3c4/04-results.md`, `05-qualification.md`, `evidence/qualification/`, `evidence/manuscript/` and an artifact/source/alignment verifier.
- Modify `docs/foundational/phase3/README.md`, `paper/alan-dissertation.md`, `paper/alan-dissertation.tex`, `paper/references.bib` only as warranted, and rebuild `paper/alan-dissertation.pdf`.

**Interfaces:**
- Consumes accepted Task 1/2 records and controller handoff `/tmp/alan-phase3c4/task3-handoff.md`; baseline `/tmp/alan-phase3c4/baseline-inputs.json` has 683 tracked inputs already matched to current bytes and actual merged Git blobs. Fresh root checks will be captured under `/tmp/alan-phase3c4/qualification/` after Task 1 source freezes.
- Produces scoped results, qualification and synchronized changed manuscript content/PDF; reviewer/delivery records remain controller-owned.

- [ ] Check new claims against owning source/evidence: derivative unions insufficient, profile equality sufficient for contextual nonblocking and closed only under the declared product, generic strictness, and nonblocking versus inevitable completion. Distinguish new abstract evidence from historical concrete captures and unestablished full-agent/runtime behavior.
- [ ] Archive exact controller build/module/package checks. Recheck all 683 baseline inputs and newly accepted Go/checker/source identities, including pinned predecessor bindings. Compare fresh package roster with actual merged prior roster, not a guessed full-suite count. Attribute historical tests/lint accurately; no fresh whole-root/global-clean claim.
- [ ] Write results/qualification with actual measured counts, retained RED/hostile evidence, source identity and limitations. Check bibliography implications against controller's inspected primary-literature note; no invented citation or stronger borrowed theorem claim.
- [ ] Read/apply PDF skill and coordinate its one-time operation marker with controller before source authoring. Align new section and changed Abstract/Conclusion in Markdown and LaTeX. Correct the eleven-form count and the overstatement in Markdown 3.4; inspect equivalent condensed TeX locations and include concise aligned clarification. Keep historical theorem numbering/open status and disclose remaining inherited format differences.
- [ ] Build baseline/current PDF with cached tooling, extract text, compare pages/warnings and visually inspect all changed pages plus relevant continuations. Record hashes, commands/exits and scoped visual coverage; preserve any initial failures/repairs. No unrelated global normalization.
- [ ] Execute the deterministic source/artifact/alignment verifier; check new local links, corrected count and congruence wording, exact input hashes and actual evidence counts. Inspect complete task diff, commit only Task 3, and report remaining obligations and review inputs.
