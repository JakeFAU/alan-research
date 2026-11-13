# Task 1 independent spec and quality review

Reviewed base `e339c2db337cd5b5a93990b8a42160994a3ec2f4` through task head `9a28d2ca4a4e031f1c98860154bbf0ed4aa7fb6e` using the complete supplied 1,704-line diff, task brief/report, reviewer instructions and `docs/foundational/phase3b6/01-contract.md`.

Spec compliance: **APPROVE**.

Task code quality: **APPROVE**.

These verdicts cover Task 1's operational model, exporter, tests and evidence. They do not certify the future independent ledger, concrete evidence relation, manuscript reconciliation or whole-branch qualification.

## Findings

No Critical, Important or Minor actionable findings identified. No required correction for this task.

## Evidence and strengths

- `research/effectobservation/model.go:135`: the reachable-domain transition relation keeps launch, ingress, application, volatile response, durable record and consumption separate. Loss clears only the permit and volatile response while preserving packet/world history and retained/consumed results. Endpoint ingress/application remain possible after loss and exact b revocation; current revocation does not rewrite admission or c. Reopen is once-only and advances the shared activation. Fresh refusal cannot create result custody or another contact.
- `research/effectobservation/model.go:208` and `:218`: Store and dossier projections match the declared boundary. The local projection omits transport/world/volatile response; the dossier adds ingress/application without importing a result. Fixed comparable binding and typed result/owner/event symbols are suitable for a separate interpreter.
- `research/effectobservation/explore.go:18`: deterministic successor order and first-discovery BFS produce stable indices and copied shortest prefixes without prefix aliasing. `:47` compares complete event values, including otherwise irrelevant fields, so replay rejects malformed descriptors. Its explicit operational status avoids presenting shared transition code as an independent oracle.
- `research/effectobservation/explore_test.go:10` checks each enumerated state/edge, immutable history, stage separation, state uniqueness, first discovery ordering and shortest-prefix reconstruction. `:163` mutates every descriptor coordinate on every outgoing arc. The six behavioral tests exercise both authentic values, late responses, retained recovery, revocation and refused gates.
- `research/effectobservation/cmd/export/main.go:17`: argument errors return 2, JSON write errors return 1, and the complete graph is exported deterministically. Focused command tests cover usage, write failure, complete JSON round-trip and byte equality.
- `docs/foundational/phase3b6/02-model.md:57` correctly distinguishes a shortest state prefix from a history-qualified witness. The retained control explicitly ends in consumption after reopen; merely reaching the same final state would not establish that history. `:59` explicitly leaves independent interpretation/partitioning and sealed-runtime correspondence pending and avoids a general progress claim.

## Artifact inspection performed

Read all compressed task logs. `red.log.gz` contains six ordinary behavioral failures in compiled tests, rather than only missing-symbol diagnostics. GREEN and complete-package logs contain success records. The initial exporter compilation failure and denied default Staticcheck cache are preserved and distinguished from reported successful later checks in `verification.md`.

One focused read-only Python inspection addressed the concrete risk that the hand-selected summary paths, especially the retained control, might not be paths in the archived graph. It decoded only task JSON/gzip evidence, traversed archived arcs, and checked all archived prefixes against graph distances. It exited 0: 63 reachable states, 104 arcs, maximum shortest distance 9; every archived prefix reconstructs its indexed state at minimum graph distance. All summary descriptors identify unique outgoing arcs and reach their reported state IDs.

The inspected boundary paths are:

- Equal Store projection, no application versus v1: `0→2→7→15→25` and `0→1→4→9→17→28→39→49` (4 and 7 events).
- Distinct unretained values: state 49 versus `0→1→4→10→20→32→43→54`, both 7 events. The states differ in applied value but agree on all declared Store coordinates.
- Provisional negative cut: `0→1→5→13→22→28`, with owner stopped and b revoked at 13, ingress at 22 and application at 28.
- Retention before loss and consumption after reopen: `0→1→4→9→16→26→37→47→57→61`; record occurs before loss, and the last consume descriptor uses incarnation 2 with b already revoked.

This inspection checks the archived graph's linkage and supplied witnesses, not independent semantic validity of every operational edge. No Go tests, exports, formatting or Staticcheck checks were rerun. No old source was inspected, no archived SQLite was opened, and no repository/index/branch state was changed other than writing this requested review report.

## Cannot verify and controller obligations

No Task 1 requirement remains unverifiable to the degree required for this task-scoped review. The verification commands/exits are attributed to the implementer's retained evidence and report; the original transient stub source and successful empty-output formatting/Staticcheck invocations were not independently reproduced here. The supplied manifests identify final implementation/test and artifact bytes; this review does not authenticate execution chronology.

Before whole-branch delivery, the controller must still establish the independent ledger's separation from operational helpers, all-edge/merged-target auditing and independently derived observation partitions; concrete causal/physical/cut attribution; compatibility of historical outputs; and final manuscript/qualification/review requirements. Those are explicitly later tasks, not missing Task 1 implementation.
