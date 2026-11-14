# Phase 3C.2 independent evidence audit

## Question and result

Does the exported finite packet contain the complete resource-aware relation,
exact static root/holder identity, and genuine concrete histories supporting the
stated completion and custody witnesses?

The separate Python standard-library checker reconstructs all 72 abstract
configuration graphs: **1,632 states, 2,412 arcs, 1,624 CanComplete labels and 8
CannotComplete labels**. All stable-ready states have a finite completion suffix;
the shared-limit-two, separate-limit-one and servicing-B controls also do.
It reconstructs all 72 static regrouping records and replays the ten exported
histories, **161 concrete microsteps**, from empty initial custody.

This is bounded finite validation and constructive history checking. It does
not mechanize the conditional source argument, establish runtime qualification,
prove distributed allocation, exhaust the concrete product, or establish
fairness, trace equivalence or arbitrary-context substitution.

## Independent algorithms and shared specification

`tools/audit.py` is a new declarative implementation. It neither imports nor
calls Go operational guards, successors, projection, lifting, exploration or
runtime Replay. Producer source files are read only as uninterpreted bytes for
SHA256 binding. The implementation was derived from this packet's frozen
`00-contract.md`, `01-argument.md` and `02-evidence-schema.md`, with prior custody
semantics clarified by `../phase3a/01-contract.md` and the Phase 3C.1
contract/argument/schema. Those shared specifications are archived by hash in
the Task 2 manifest. Independent algorithms are not an independently discovered
specification.

The checker constructs the full BFS relation using named original stages and
held bits, including every arc to an already discovered state. A reverse
adjacency traversal computes completion reachability. It compares every state,
action, target, parent, label and count with strict recursive type equality.
Neither the exported BFS tree nor exported CanComplete flags serve as an
oracle. Component records are independently reconstructed from exact owner
sets, ports, outgoing bindings, layout and root memberships; both closed
parenthesizations must equal that same complete configuration/relation.

Concrete replay begins with a constructed empty state. Each event is checked
against predecessor custody, physical availability, actor, exact binding,
service readiness and capacity. Python produces the complete next state and
compares every field with the exported post-state. Only an actual `apply` event
creates result bytes and historical WorldResult; transmission, reception,
recording and consumption obtain bytes from their causal custody source.
The historical field is never used to enable an action or synthesize a receipt.
Crash/recovery cannot alter holds, durable data or logical identity.

The audit compares entire Base records and root totals for the two holder pairs,
checks prefix extensions, ready-now invalidation at BC commitment, the retained
crashed AB holder, duplicate request/result/ACK delivery, both pre-record losses,
release-before-consumption ordering and complete r1/r2 history correspondence.
For the old early trace, it computes stable-ready trace membership with an
epsilon closure over reserve-capacity only. The erased old trace is rejected;
visible BC progress cannot be silently used as preparation.

## Operational reports remain a separate evidence category

The exported factored probe summary reports **311 old edge states, 1,854,036
representatives, 11,972,664 micro-arcs, 2,451,168 lifts and 4,833,504 lift
microsteps**. The Python audit checks exact domain/roster, strict count types,
row totals and source identities. **It does not independently re-enumerate
these operational probes.** Its JSON output explicitly sets
`operational_probes_independently_reenumerated` to false. Source hashing binds
these reports to named producer bytes; it is not independent corroboration of
the reported probe counts or Go lifting/noninterference behavior.

The fresh exporter invocation lives in separate `evidence/task2/capture.py`
capture tooling. The checker itself never launches it. The captured plain JSON
stdout is byte-identical to the reviewed Task 1 capture, raw SHA256
`a7cf316219b30cdc83d455b95a16549c022d55aa422b642b8f5976315e0c0fd2`.
It is archived as `evidence/task2/export.json.gz` with deterministic gzip
packaging; the CLI continues to emit plain JSON. Both formats yield identical
audit summaries.

## Adverse controls and independently authored histories

The unittest suite and capture tooling retain **43** exact full hostile inputs,
rejection stderr, stdout, commands, exits and intended reasons in
`evidence/task2/hostile/` and `hostile-index.json`. These are **forged evidence
controls**, not operational counterexamples or historical executions.

Controls challenge omitted non-tree arcs, incorrect labels/parents/projections,
aggregate or changed holders, invented memberships, conflicting limits, copied
pools, overflow, duplicate reservation/application, premature capacity or
receipt release, missing actual application, forged receipt/binding/activation,
wrong/down actors, premature ACK/consumption, invented initial/post-state facts,
unknown fields, missing configurations, null arrays, duplicate JSON keys,
nonfinite numbers, Boolean/float counts and source/probe metadata. A removed
non-tree arc has its graph arc count reduced as well: superficial arc totals
and every BFS predecessor still agree, but full-relation comparison rejects it.
Two additional controls preserve legal replay while falsifying the specified
paired-history/custody-order evidence; they test witness meaning beyond legality.

As checker tests, a separately authored BC-first schedule generates **144 Python histories**
(all 72 configurations × both values) without reading exported events or states.
These 7,200 generated microsteps are test fixtures, not additional independent
semantic evidence. Each completes with physical recovery of all three owners, duplicate deliveries,
actual durable receipt custody and retries after consumption/release. These are
archived in `python-generated-histories.json.gz`. Generation and replay share
the Python transition implementation: they are independent of Go, not a second
independent oracle for the Python checker. The hostile evidence tests and
hand-specified terminal assertions provide additional checks on that checker.

The initial missing-checker acceptance failure and two subsequently exposed
witness-contract omissions are retained as historical RED logs. The final test
log records current success; those RED logs are not claims of current defects.
No production code, dependencies, old evidence, SQLite stores or environment
configuration were changed or opened.

## Reproduction and archive verification

From the repository root:

```sh
PYTHONPYCACHEPREFIX=/tmp/alan-phase3c2/pycache python3 -m unittest discover -s docs/foundational/phase3c2/tools -p 'test_*.py' -v
python3 docs/foundational/phase3c2/tools/audit.py docs/foundational/phase3c2/evidence/task2/export.json.gz
PYTHONPYCACHEPREFIX=/tmp/alan-phase3c2/pycache python3 docs/foundational/phase3c2/evidence/task2/capture.py
PYTHONPYCACHEPREFIX=/tmp/alan-phase3c2/pycache python3 docs/foundational/phase3c2/evidence/task2/capture.py --manifest
PYTHONPYCACHEPREFIX=/tmp/alan-phase3c2/pycache python3 docs/foundational/phase3c2/evidence/task2/capture.py --verify
```

Capture invokes `go run ./research/capacitycomposition/cmd/export` with
`GOCACHE=/tmp/alan-phase3c2/gocache`. It records every command and output;
`--skip-export` reuses the already captured scratch stdout when only the
remaining capture stages need rerunning. Manifest verification checks packed
and uncompressed artifact hashes, Task 1 archive hashes, exact source-hash
lists and current producer source bytes. The manifest excludes itself.

Task 3 owns manuscript/roadmap/PDF reconciliation and repository-wide
qualification. This audit changes no mathematical or operational semantics and
does not independently assert whole-repository verification.
