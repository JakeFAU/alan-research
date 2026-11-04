# Phase 3A implementation plan

> **For agentic workers:** Use superpowers:subagent-driven-development for the independent model task and its review. Track the checks below.

**Goal:** Falsify naive two-agent useful recovery and qualify a conditional result-custody contract.

**Architecture:** A dependency-free finite Go model separate from the production executor, paired with an independent trace checker and deterministic JSON evidence. Reconcile its bounded conclusions with both manuscript representations.

**Tech stack:** Existing Go 1.27.3, standard library, Python 3 for evidence checks, existing Tectonic PDF build.

**Spec:** `docs/foundational/phase3a/01-contract.md` and `00-charter.md`.

## Global constraints

- Exactly two logical agents, one stable delegation, one immutable request, one application, two possible opaque authentic result values, at most one A-only interruption.
- No production-runtime changes, dependencies, paid endpoints or distributed implementation claims.
- Freeze local contracts and retain a shortest counterexample in the declared transition alphabet; do not claim universal minimality.
- Report model safety, completion reachability, assumptions and implementation correspondence separately.
- Use Staticcheck for lint; do not substitute a successful test for a proof.

## Review focus

Lost last receipt despite locally legal release; response channel surviving A's crash;
B progress during A downtime; durable consumption surviving recovery; invalid
trace/payload accepted by a checker that merely shares the transition function.

## Task 1: model, independent oracle and evidence

**Files:** Create `research/twoagent/{model.go,explore.go,oracle.go,model_test.go,README.md}` and `research/twoagent/cmd/export/main.go`.

**Interfaces:** `Initial() State`, `Successors(Profile, State) []Edge`, `Explore(Profile) Report`, `CheckTrace(Profile, []Event) error`. Final field names are the model implementer's choice and must be documented. The export command writes deterministic JSON for both profiles, shortest loss witness, paired indistinguishable worlds, and safety/reachability counts.

- [x] Add behavioral tests and observe the naive useful-recovery hypothesis fail on an actual trace; preserve its expected-failure reproduction outside the passing suite.
- [x] Implement the fixed finite transition system with stable identity, distinct record/consume/ACK, one A crash, and asynchronous B/channel steps.
- [x] Check local safety separately from custody; enumerate the complete reachable graph without arbitrary depth limits; record shortest witness using BFS and reverse-reachability completion checks.
- [x] Independently check legal trace prefixes without calling `Successors`; reject fabricated values, changed correlation, early durable ACK, duplicate application and consume-before-record. Exercise receipt delivery both before and after recovery.
- [x] Export deterministic evidence and reproduce it; run focused normal/race tests and Staticcheck.

## Task 2: research reconciliation and independent review

**Files:** Create `docs/foundational/phase3a/{03-results.md,04-runbook.md,evidence/*}`; modify `paper/alan-dissertation.{md,tex}`, its PDF and the root research index.

**Interfaces:** Consume Task 1's exact graph, traces and command output; no handwritten qualification counts.

- [x] Review the model against the frozen local and edge contracts and try to falsify the conditional proof.
- [x] Archive source hashes, exact commands, test output, JSON, the rejected hypothesis and review findings.
- [x] State the paired-world argument and custody induction with all premises, and separate previous findings and current model evidence.
- [x] Synchronize Markdown/LaTeX; add the checked Saltzer/Reed/Clark reference to BibTeX; rebuild and inspect the affected PDF pages.
- [x] Run the ordinary repository suite, focused race/lint, artifact/link checks and an independent whole-change review.

**Delivery:** Commit, push the task branch and open a pull request under the repository's finishing preference. Git and the pull request record delivery; Phase 3B/3C remain open.
