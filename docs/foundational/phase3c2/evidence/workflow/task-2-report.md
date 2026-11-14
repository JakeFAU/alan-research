# Task 2 report — independent Phase 3C.2 evidence audit

Status: DONE.

Commit: `a24f8ff1c6aac22a0cffcffc2c614ccfca819c71` — `Audit shared-capacity evidence independently`.

Base: `60d46a43dd32caa07f0b814e6741213f506a9cb6`.
Branch: `codex/phase3c2-shared-capacity`.
Workdir: `/Users/jacob/.codex/worktrees/phase3b4-research/alan-research`.

## Question and bounded result

Independently reconstruct the complete frozen resource-aware abstract relation,
static grouping/root identity and real-history custody evidence. The checker
accepts 72 complete graphs: 1,632 states, 2,412 arcs, 1,624 CanComplete and 8
CannotComplete labels; 72 exact grouping records; ten genuine exported
histories with 161 full concrete microsteps. No material specification/schema
contradiction was encountered.

All stable-ready abstract states, and the shared-2/separate-1/servicing-B
controls, have a finite completion suffix. The checker verifies captured
holder-pair full Base equality and equal occupancy with opposite completion,
ready-now invalidation at BC commit, stable-ready rejection of the old erased
AB-before-BC trace, holder preservation under crash, causal durable ACK custody,
both result values and complete post-state equality.

## Files and archive

Created `docs/foundational/phase3c2/tools/audit.py`, `tools/test_audit.py`,
`03-independent-audit.md`, and `evidence/task2/` only: 198 new tracked Task 2 files.
The latter contains capture tooling, a byte-identical fresh export, 43 complete
compressed forged inputs and their commands/stdout/stderr, source hashes,
non-tree-control measurements, Python checker-test fixtures, RED/final test
logs, and an archive manifest. Each full file is enumerated by the manifest;
the 201 entries include the new files except the manifest itself plus four
unchanged shared-specification inputs.

Raw fresh export SHA256:
`a7cf316219b30cdc83d455b95a16549c022d55aa422b642b8f5976315e0c0fd2`.
Fresh export executed once, exit 0, 134.308554 seconds. Raw bytes: 4,514,212.
The checker accepts plain JSON and deterministic gzip with identical summaries.
The producer CLI remains plain JSON.

The non-tree control removes source node 2's AB commit arc to node 3 in the
first graph. All 24 states and all BFS predecessor arcs survive. The reported
arc count is corrected from 34 to 33, matching actual remaining arcs. Rejection
is specifically `Graphs[0].Nodes[2].Arcs: length mismatch`, not a final count
mismatch. Exact mutation and measured proof are retained.

## Independence and limits

No Go operational guard/source implementation was read as an oracle, imported,
or invoked by the checker. It hashes producer bytes only. Shared specifications
are this packet's frozen contract/argument/schema, prior Phase 3C.1
contract/argument/schema and Phase 3A durable-ACK contract. The archive binds
those clarification documents by hash. Algorithms are separately implemented;
the specification was shared, not independently discovered.

Operational probes remain source-bound reported results: 311 old edge states,
1,854,036 representatives, 11,972,664 micro-arcs, 2,451,168 lifts, 4,833,504 lift
microsteps. Their domain/roster/types/totals are structurally checked. Python
does not independently re-enumerate that probe domain or certify the Go lifting
and oracle-noninterference checks. The summary explicitly says so.

The 144 Python-authored BC-first histories (7,200 steps) are checker test fixtures,
not additional independent semantic evidence. They do not consume producer
histories but share the Python transition rules with replay. They assert
hand-specified completion, actual recorded values, released holdings/custody and
three recovered activations. Primary independent evidence is the complete
abstract reconstruction and replay of ten captured Go histories.

This is not runtime/distributed qualification, whole-product enumeration,
mechanized proof, fairness transfer, cost equivalence or arbitrary-context
substitution. Manuscript/roadmap/PDF reconciliation and broader qualification
remain Task 3 responsibilities. No production/dependency/old evidence edits,
SQLite opens, environment configuration reads or live calls occurred.

## Verification and historical failures

The final unittest command passed 5 tests (6.453 seconds), including 43 hostile
subcases and 144 generated-history subcases. Every hostile CLI run exited 1
with its intended reason and retained exact input/output. Genuine plain/gzip
audits exited 0. All 32 producer source hashes match current files.

Initial RED: the acceptance test failed because audit.py did not yet exist
(exit 1 overall, subprocess exit 2), retained in `red.log`. Later the legal
but non-corresponding r2-history and release-after-consumption controls were
initially accepted; two failing assertions (exit 1) are retained in
`red-witness-contract.log`. Full-history value symmetry and release-order
requirements corrected those omissions; final controls reject both. There are
no remaining failing tests from this task.

Capture orchestration executed:

```sh
PYTHONPYCACHEPREFIX=/tmp/alan-phase3c2/pycache python3 docs/foundational/phase3c2/evidence/task2/capture.py
PYTHONPYCACHEPREFIX=/tmp/alan-phase3c2/pycache python3 docs/foundational/phase3c2/evidence/task2/capture.py --skip-export
PYTHONPYCACHEPREFIX=/tmp/alan-phase3c2/pycache python3 docs/foundational/phase3c2/evidence/task2/capture.py --manifest
PYTHONPYCACHEPREFIX=/tmp/alan-phase3c2/pycache python3 docs/foundational/phase3c2/evidence/task2/capture.py --verify
```

All exited 0. The second invocation tested final capture metadata/non-tree
checks while reusing the already successful unchanged fresh exporter stdout.
The first two outputs were:

```json
{"fresh_export_sha256": "a7cf316219b30cdc83d455b95a16549c022d55aa422b642b8f5976315e0c0fd2", "generated_histories": 144, "hostile_controls": 43}
```

Final manifest verification output:

```json
{"verified_source_hashes": 32, "verified_task1_raw_sha256": "a7cf316219b30cdc83d455b95a16549c022d55aa422b642b8f5976315e0c0fd2", "verified_task2_artifacts": 201}
```

`git diff --cached --check` exited 0. A complete inventory/blob comparison
verified all 201 manifest records against staged or unchanged Git blobs,
including packed/raw hashes; only the 198 Task 2 files are staged. Every
command metadata exit and referenced output path was checked.

Controller-owned preflight evidence remains untracked/unstaged, and the
controller's phase3c2-plan.md changes remain unstaged. The first ordinary git
add failed on sandbox index.lock permissions; narrowly scoped escalation
succeeded. A diagnostic ps command was denied by the sandbox; it was not needed
for any result. No automatic approval rejection remains unresolved.

## Self-review

Read complete checker/test/capture/document changes and inspected the staged
artifact inventory, logs, hostile-index, command metadata, non-tree control,
source hashes and archive hashes. Checked every manifest file against the Git
index to catch omitted ignored logs; explicitly staged only task-owned logs.
No unrelated files, caches, temporary investigations, profiler outputs, source
changes or old evidence enter this commit. Scratch/caches live under
`/tmp/alan-phase3c2/`; initial Python caches in tools were removed.

Reviewed strict unknown-field and duplicate-key/type handling, deterministic BFS
parent/all-arc correspondence, reverse completion, exact roots/memberships,
causal Apply provenance, full post-state equality, actor/availability, ACK
release, static grouping and stable-ready erased-trace membership. The two
self-review findings above were corrected and retested. No material concerns
remain beyond the declared independent-check boundary and pending independent
review of this task.

Root build/module checks and broad Go qualification were not duplicated here;
the controller owns those captures. No subagents, PR, push or merge performed.

## Exact retained subprocess commands and outputs

All commands below ran in the workdir above. Capture set
`GOCACHE=/tmp/alan-phase3c2/gocache` and
`PYTHONPYCACHEPREFIX=/tmp/alan-phase3c2/pycache`; it did not read .envrc.

### audit-gzip.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/export.json.gz
```

Exit 0; 0.115993 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/audit-gzip.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/audit-gzip.stderr.log`.

```text
{
  "arcs": 2412,
  "can_complete": 1624,
  "cannot_complete": 8,
  "genuine_histories": 10,
  "independent_abstract_graphs": 72,
  "operational_probes_independently_reenumerated": false,
  "replayed_microsteps": 161,
  "reported_operational_probes": {
    "Configurations": 72,
    "Domain": "factored-old-edge-v1",
    "LiftMicroSteps": 4833504,
    "Lifts": 2451168,
    "MicroArcs": 11972664,
    "OldEdgeStates": 311,
    "Representatives": 1854036
  },
  "schema": "alan/phase3c2/v1",
  "source_hashes_verified": 32,
  "states": 1632,
  "static_groupings": 72
}
```

### audit-plain.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /tmp/alan-phase3c2/task2-fresh-export.json
```

Exit 0; 0.120386 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/audit-plain.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/audit-plain.stderr.log`.

```text
{
  "arcs": 2412,
  "can_complete": 1624,
  "cannot_complete": 8,
  "genuine_histories": 10,
  "independent_abstract_graphs": 72,
  "operational_probes_independently_reenumerated": false,
  "replayed_microsteps": 161,
  "reported_operational_probes": {
    "Configurations": 72,
    "Domain": "factored-old-edge-v1",
    "LiftMicroSteps": 4833504,
    "Lifts": 2451168,
    "MicroArcs": 11972664,
    "OldEdgeStates": 311,
    "Representatives": 1854036
  },
  "schema": "alan/phase3c2/v1",
  "source_hashes_verified": 32,
  "states": 1632,
  "static_groupings": 72
}
```

### fresh-export.command.json

```sh
go run ./research/capacitycomposition/cmd/export
```

Exit 0; 134.308554 seconds. stdout: `/tmp/alan-phase3c2/task2-fresh-export.json`; stderr: `docs/foundational/phase3c2/evidence/task2/fresh-export.stderr.log`.

```text
4,514,212 bytes plain CLI JSON; exact archived bytes in export.json.gz; stderr empty.
```

### hostile/actor-down.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/actor-down.json.gz
```

Exit 1; 0.107373 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/actor-down.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/actor-down.stderr.log`.

```text
REJECT: Witnesses[8].Steps[8].Event: actor down
```

### hostile/aggregate-holder.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/aggregate-holder.json.gz
```

Exit 1; 0.048337 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/aggregate-holder.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/aggregate-holder.stderr.log`.

```text
REJECT: Graphs[0].Nodes[0].State.Held: type mismatch
```

### hostile/boolean-count.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/boolean-count.json.gz
```

Exit 1; 0.047748 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/boolean-count.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/boolean-count.stderr.log`.

```text
REJECT: Graphs[0].States: type mismatch
```

### hostile/boolean-generation.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/boolean-generation.json.gz
```

Exit 1; 0.047891 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/boolean-generation.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/boolean-generation.stderr.log`.

```text
REJECT: Graphs[0].Config.Base.Agents[0].Generation: type mismatch
```

### hostile/capacity-overflow.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/capacity-overflow.json.gz
```

Exit 1; 0.099864 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/capacity-overflow.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/capacity-overflow.stderr.log`.

```text
REJECT: Witnesses[0].Steps[6].State.Held[1]: value mismatch
```

### hostile/changed-holder.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/changed-holder.json.gz
```

Exit 1; 0.100055 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/changed-holder.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/changed-holder.stderr.log`.

```text
REJECT: Witnesses[0].Steps[6].State.Held[0]: value mismatch
```

### hostile/conflicting-limit.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/conflicting-limit.json.gz
```

Exit 1; 0.047769 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/conflicting-limit.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/conflicting-limit.stderr.log`.

```text
REJECT: Groupings[0].AB.Roots[0].Limit: value mismatch
```

### hostile/consume-before-record.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/consume-before-record.json.gz
```

Exit 1; 0.106769 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/consume-before-record.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/consume-before-record.stderr.log`.

```text
REJECT: Witnesses[8].Steps[20].Event: consume before durable record
```

### hostile/copied-pool.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/copied-pool.json.gz
```

Exit 1; 0.048812 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/copied-pool.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/copied-pool.stderr.log`.

```text
REJECT: Groupings[0].Left.Roots: length mismatch
```

### hostile/duplicate-application.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/duplicate-application.json.gz
```

Exit 1; 0.10978 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/duplicate-application.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/duplicate-application.stderr.log`.

```text
REJECT: Witnesses[8].Steps[13].Event: illegal apply
```

### hostile/duplicate-json-key.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/duplicate-json-key.json.gz
```

Exit 1; 0.045321 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/duplicate-json-key.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/duplicate-json-key.stderr.log`.

```text
REJECT: duplicate key: Schema
```

### hostile/duplicate-reservation.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/duplicate-reservation.json.gz
```

Exit 1; 0.099229 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/duplicate-reservation.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/duplicate-reservation.stderr.log`.

```text
REJECT: Witnesses[0].Steps[7].Event: illegal reservation
```

### hostile/early-release-poststate.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/early-release-poststate.json.gz
```

Exit 1; 0.110618 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/early-release-poststate.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/early-release-poststate.stderr.log`.

```text
REJECT: Witnesses[7].Steps[7].State.Held[0]: value mismatch
```

### hostile/float-count.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/float-count.json.gz
```

Exit 1; 0.048154 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/float-count.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/float-count.stderr.log`.

```text
REJECT: Graphs[0].Arcs: type mismatch
```

### hostile/forged-binding.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/forged-binding.json.gz
```

Exit 1; 0.099167 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/forged-binding.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/forged-binding.stderr.log`.

```text
REJECT: Witnesses[0].Steps[0].Event.Binding.Grant: value mismatch
```

### hostile/forged-receipt.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/forged-receipt.json.gz
```

Exit 1; 0.105062 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/forged-receipt.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/forged-receipt.stderr.log`.

```text
REJECT: Witnesses[8].Steps[13].Event: result provenance
```

### hostile/invented-initial-receipt.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/invented-initial-receipt.json.gz
```

Exit 1; 0.100693 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/invented-initial-receipt.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/invented-initial-receipt.stderr.log`.

```text
REJECT: Witnesses[0].Initial.Base.Originals[0].Recorded: value mismatch
```

### hostile/invented-membership.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/invented-membership.json.gz
```

Exit 1; 0.047859 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/invented-membership.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/invented-membership.stderr.log`.

```text
REJECT: Groupings[0].A.Memberships[0].Root: value mismatch
```

### hostile/missing-config.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/missing-config.json.gz
```

Exit 1; 0.047057 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/missing-config.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/missing-config.stderr.log`.

```text
REJECT: Graphs: roster
```

### hostile/missing-genuine-apply.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/missing-genuine-apply.json.gz
```

Exit 1; 0.106449 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/missing-genuine-apply.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/missing-genuine-apply.stderr.log`.

```text
REJECT: Witnesses[8].Steps[12].Event: send-result custody/channel
```

### hostile/missing-projection.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/missing-projection.json.gz
```

Exit 1; 0.098403 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/missing-projection.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/missing-projection.stderr.log`.

```text
REJECT: Witnesses[0].Projection: length mismatch
```

### hostile/noncorresponding-r2-history.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/noncorresponding-r2-history.json.gz
```

Exit 1; 0.114364 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/noncorresponding-r2-history.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/noncorresponding-r2-history.stderr.log`.

```text
REJECT: corresponding recovery histories.Steps: length mismatch
```

### hostile/nonfinite-json-number.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/nonfinite-json-number.json.gz
```

Exit 1; 0.023807 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/nonfinite-json-number.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/nonfinite-json-number.stderr.log`.

```text
REJECT: nonfinite JSON number: NaN
```

### hostile/null-array.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/null-array.json.gz
```

Exit 1; 0.049372 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/null-array.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/null-array.stderr.log`.

```text
REJECT: Groupings[0].C.Roots: type mismatch
```

### hostile/poststate-invention.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/poststate-invention.json.gz
```

Exit 1; 0.09904 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/poststate-invention.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/poststate-invention.stderr.log`.

```text
REJECT: Witnesses[0].Steps[0].State.Base.Originals[1].Committed: value mismatch
```

### hostile/premature-ack.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/premature-ack.json.gz
```

Exit 1; 0.108761 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/premature-ack.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/premature-ack.stderr.log`.

```text
REJECT: Witnesses[8].Steps[20].Event: durable ACK
```

### hostile/premature-receipt-release.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/premature-receipt-release.json.gz
```

Exit 1; 0.106563 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/premature-receipt-release.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/premature-receipt-release.stderr.log`.

```text
REJECT: Witnesses[8].Steps[20].Event: release ACK
```

### hostile/probe-boolean.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/probe-boolean.json.gz
```

Exit 1; 0.11605 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/probe-boolean.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/probe-boolean.stderr.log`.

```text
REJECT: Probes.OldEdgeStates: nonnegative integer required
```

### hostile/probe-row-omitted.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/probe-row-omitted.json.gz
```

Exit 1; 0.116466 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/probe-row-omitted.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/probe-row-omitted.stderr.log`.

```text
REJECT: probe roster
```

### hostile/probe-total-forged.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/probe-total-forged.json.gz
```

Exit 1; 0.115081 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/probe-total-forged.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/probe-total-forged.stderr.log`.

```text
REJECT: probe total Lifts
```

### hostile/removed-non-tree-arc.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/removed-non-tree-arc.json.gz
```

Exit 1; 0.048315 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/removed-non-tree-arc.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/removed-non-tree-arc.stderr.log`.

```text
REJECT: Graphs[0].Nodes[2].Arcs: length mismatch
```

### hostile/source-forged.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/source-forged.json.gz
```

Exit 1; 0.046502 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/source-forged.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/source-forged.stderr.log`.

```text
REJECT: Sources[0].SHA256: value mismatch
```

### hostile/standalone-capacity-release.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/standalone-capacity-release.json.gz
```

Exit 1; 0.099425 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/standalone-capacity-release.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/standalone-capacity-release.stderr.log`.

```text
REJECT: Witnesses[0].Steps[6].Event: event kind
```

### hostile/unknown-event-field.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/unknown-event-field.json.gz
```

Exit 1; 0.097288 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/unknown-event-field.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/unknown-event-field.stderr.log`.

```text
REJECT: Witnesses[0].Steps[0].Event: keys mismatch
```

### hostile/unknown-nested-field.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/unknown-nested-field.json.gz
```

Exit 1; 0.04814 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/unknown-nested-field.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/unknown-nested-field.stderr.log`.

```text
REJECT: Graphs[0].Nodes[0].State: keys mismatch
```

### hostile/unknown-top-field.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/unknown-top-field.json.gz
```

Exit 1; 0.046576 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/unknown-top-field.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/unknown-top-field.stderr.log`.

```text
REJECT: Evidence: keys mismatch
```

### hostile/wrong-activation.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/wrong-activation.json.gz
```

Exit 1; 0.106701 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-activation.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-activation.stderr.log`.

```text
REJECT: Witnesses[7].Steps[7].State.Base.Agents[1].Activation: value mismatch
```

### hostile/wrong-actor.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/wrong-actor.json.gz
```

Exit 1; 0.099331 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-actor.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-actor.stderr.log`.

```text
REJECT: Witnesses[0].Steps[0].Event: actor
```

### hostile/wrong-completion.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/wrong-completion.json.gz
```

Exit 1; 0.050513 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-completion.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-completion.stderr.log`.

```text
REJECT: Graphs[0].Nodes[0].CanComplete: value mismatch
```

### hostile/wrong-final-label.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/wrong-final-label.json.gz
```

Exit 1; 0.097767 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-final-label.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-final-label.stderr.log`.

```text
REJECT: Witnesses[0].CanComplete: value mismatch
```

### hostile/wrong-parent.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/wrong-parent.json.gz
```

Exit 1; 0.057372 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-parent.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-parent.stderr.log`.

```text
REJECT: Graphs[0].Nodes[1].Parent: value mismatch
```

### hostile/wrong-recovery-custody-order.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/wrong-recovery-custody-order.json.gz
```

Exit 1; 0.116615 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-recovery-custody-order.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-recovery-custody-order.stderr.log`.

```text
REJECT: recovery release before consumption
```

### hostile/wrong-schema.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/tools/audit.py /Users/jacob/.codex/worktrees/phase3b4-research/alan-research/docs/foundational/phase3c2/evidence/task2/hostile/wrong-schema.json.gz
```

Exit 1; 0.045912 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-schema.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/hostile/wrong-schema.stderr.log`.

```text
REJECT: Schema: value mismatch
```

### unittest.command.json

```sh
/Library/Developer/CommandLineTools/usr/bin/python3 -m unittest discover -s docs/foundational/phase3c2/tools -p 'test_*.py' -v
```

Exit 0; 6.486153 seconds. stdout: `docs/foundational/phase3c2/evidence/task2/unittest.stdout.log`; stderr: `docs/foundational/phase3c2/evidence/task2/unittest.stderr.log`.

```text
test_genuine_export_is_independently_accepted (test_audit.AuditAcceptance) ... ok
test_hostile_records_rejected_for_specific_reasons (test_audit.HostileEvidence) ... ok
test_plain_json_and_gzip_have_identical_meaning (test_audit.HostileEvidence) ... ok
test_completion_under_every_configuration_and_both_results (test_audit.IndependentlyGeneratedHistories) ... ok
test_local_transport_guards_reject_impossible_events (test_audit.IndependentlyGeneratedHistories) ... ok

----------------------------------------------------------------------
Ran 5 tests in 6.453s

OK
```

## Post-commit status

Commit succeeded (exit 0). `git status --short` contains only the controller-owned
modified phase3c2-plan.md and untracked evidence/preflight directory. The index
is empty. No push, PR or merge was performed.
