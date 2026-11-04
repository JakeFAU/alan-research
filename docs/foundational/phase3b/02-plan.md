# Phase 3B.1 implementation plan

> **For agentic workers:** Use superpowers:subagent-driven-development for the independent model and public-runtime witness, with independent review of both.

**Goal:** Separate reciprocal custody safety from completion, and qualify the smallest service condition for the declared pair.

**Architecture:** A new compact independent Go model with two shared agent-availability records and two reciprocal edge records; an independent joint trace checker and projection checks against unchanged Phase 3A. A separate ordinary public-API runtime fixture exhibits only the control-order boundary.

**Tech stack:** Existing Go 1.27.3, standard library for the model, existing author/runtime/SQLite for the fixture, existing Tectonic manuscript build.

**Spec:** `docs/foundational/phase3b/00-charter.md` and `01-contract.md`.

## Global constraints

- Exactly two actual agents, one immutable request per direction, two opaque result values per edge, at most one crash per agent, four fixed service profiles.
- Preserve Phase 3A source and evidence; no production runtime changes or new dependencies.
- Fixed disjoint grants/resources and atomic application/receipt retention are premises; dynamic authority, cancellation and distributed-owner correspondence remain open.
- Safety, existential completion, fair progress, root termination and global quiescence are distinct claims.
- Use complete finite reachability without depth cutoffs; Staticcheck is the lint gate.

## Review focus

An agent crashing only its caller role; wrong-edge receipt acceptance when values
coincide; conflating retry-enabled traps with terminal graph nodes; treating
admission/dispatch ACK as consumed application work; lifting liveness through a
projection that removes enabled work. Each is addressed in its owning task.

## Task 1: reciprocal model and independent evidence

**Files:** Create `research/mutualobligations/{model.go,explore.go,oracle.go,model_test.go,README.md}` and `cmd/export/main.go` beneath that package.

**Interfaces:** `Initial() State`, `Successors(Profile, State) []Edge`, `Explore(Profile) Report`, `CheckTrace(Profile, []Event) error`. Choose documented compact field types; do not store unbounded histories in graph states. Export deterministic JSON for all four profiles and shortest witnesses.

- [x] Write behavioral tests first and observe genuine RED; reproduce the naive completion claim as an intentional failure outside the passing suite.
- [x] Implement shared agent availability, exact reciprocal binding, durable ACK, service guards, all crash/recovery cuts and duplicate/retry cycles.
- [x] Enumerate each complete graph, check local safety/custody, reverse completion reachability, Trap iff CannotComplete for serial/serial, and positive asymmetric profiles. Retain shortest trap, admitted/drained wait and recovered-wait traces plus nonvacuous applied states.
- [x] Independently check legal traces, including actor/edge bindings, service gating and both roles while down. Reject malformed correlation/grants, receipt forgery, premature ACK, duplicate application/consumption and illegal service.
- [x] Project every reachable transition onto Phase 3A or stutter, and test oracle noninterference. Keep this check separate from the joint checker and positive completion argument.
- [x] Run normal/race/Staticcheck, repeat deterministic export, and report exact counts, limits and source identities.

## Task 2: public same-Store wait witness

**Files:** Create `examples/communication/mutual/deployment.go`, `runtime_test.go`, and `README.md` (add focused construction tests only if needed).

**Interfaces:** Public author/runtime/schema only; two roots and reciprocal request/result ports. Record inspection/history evidence through public APIs. Exact helper signatures are the implementer's choice.

- [x] Read the actual APIs and existing examples; write an adverse serial waiting test before constructing the deployment.
- [x] Establish both initial Sends, then both request admissions/ACKs with public APIs. Exercise both send orders and delivery orders across serial, either asymmetric, and both servicing profiles.
- [x] Keep roots open at explicit final holds; check actual request/result consumption and validated results, not terminal outcomes or report bookkeeping alone.
- [x] Check repeated Run and whole-Store close/open/run/replay at the serial cut, including stable identities/sequences and correctly changed activation epochs.
- [x] Retain raw program/provision/inspection evidence and compact observations for the declared cases; run focused normal/integration/race checks and Staticcheck as applicable.

## Task 3: argument, manuscript and delivery

**Files:** Create `docs/foundational/phase3b/{03-results.md,04-runbook.md,evidence/*}`; update the Phase 3 index, root research index and both manuscript representations/PDF (BibTeX only for an actually checked new source).

- [x] Independently review model and runtime tasks and attack the conditional proofs; preserve concrete findings and limits.
- [x] Reconcile the model counts, exact runtime observations, proof arguments and remaining authority/cancellation/owner gaps without claiming complete Phase 3B closure.
- [x] Archive commands, test output, source hashes, deterministic export, raw runtime records and expected failures.
- [x] Synchronize the added mathematical content, rebuild and inspect affected PDF pages, and complete whole-change review.
- [x] Finish ordinary-root verification plus focused final new-package checks and package coverage; inspect every failure and avoid repeating unchanged expensive suites without cause.

**Delivery:** Commit, push and open the next PR under the user's finishing preference; do not merge it automatically.
