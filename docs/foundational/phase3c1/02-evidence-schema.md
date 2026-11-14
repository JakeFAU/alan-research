# Phase 3C.1 evidence interface, version 1

Producer: `go run ./research/contractcomposition/cmd/export` from the repository
root. UTF-8 JSON, one object, trailing newline. `Schema` is exactly
`alan.phase3c1.progress.v1`. Field names below are case-sensitive Go-exported
names. Every listed field is present; no implicit identity renaming, omitted
progress register, or unknown null sentinel is permitted. Arrays have the
specified deterministic order. This document freezes the Task 2 input.

The independent auditor must derive expected transitions, completion, reverse
reachability, policy dependencies and concrete causal events itself. Producer
flags/counts/claims are assertions to check, not an oracle. It must retain and
reject its own hostile mutations without editing the genuine export.

## Top-level and records

| Record | Exact fields and types |
| --- | --- |
| Evidence | `Schema:string`, `Graphs:Graph[]`, `Groupings:Grouping[]`, `Witnesses:Witness[]`, `Collisions:Collision[]`, `Retry:RetryCertificate`, `Claims:Claim[]` |
| Config | `Topology:string`, `ServiceMask:integer`, `Agents:Agent[]`, `Originals:Binding[]` |
| Agent | `ID:string`, `Generation:integer`, `Servicing:boolean` |
| Binding | `Original:string`, `Caller:string`, `CallerGeneration:integer`, `Operation:integer`, `Callee:string`, `CalleeGeneration:integer`, `Request:string`, `Grant:string`, `Resource:string` |
| State | `Stages:string[3]` |
| Action | `Kind:string`, `Original:string`, `Binding:Binding` |
| Arc | `Action:Action`, `To:integer` |
| Node | `State:State`, `Arcs:Arc[]`, `Complete:boolean`, `CanComplete:boolean`, `Parent:integer`, `Via:Action or null` |
| Dependency | `Original:string`, `Requires:string` |
| Certificate | `Dependencies:Dependency[]`, `Acyclic:boolean`, `Order:string[]` |
| Graph | `Config:Config`, `Nodes:Node[]`, `States:integer`, `Arcs:integer`, `CanComplete:integer`, `CannotComplete:integer`, `Certificate:Certificate` |
| Component | `Owners:Agent[]`, `Ports:Agent[]`, `Originals:Binding[]` |
| Grouping | `ConfigID:string`, `A:Component`, `B:Component`, `C:Component`, `AB:Component`, `BC:Component`, `Left:Component`, `Right:Component`, `Flat:Config`, `LeftClosed:Config`, `RightClosed:Config`, `Equal:boolean` |
| Availability | `Down:boolean`, `Crashed:boolean`, `Activation:integer` |
| Transport | `Binding:Binding`, `Committed:boolean`, `Admitted:boolean`, `Applied:boolean`, `Consumed:boolean`, `Volatile:string`, `Recorded:string`, `Retained:string`, `RequestChannel:boolean`, `ResultChannel:string`, `ACKChannel:boolean`, `WorldResult:string` |
| ConcreteState | `Agents:Availability[3]`, `Originals:Transport[3]` |
| Event | `Kind:string`, `Actor:string`, `Original:string`, `Binding:Binding`, `Result:string`, `Duplicate:boolean` |
| MicroStep | `Event:Event`, `State:ConcreteState` |
| Witness | `Name:string`, `Config:Config`, `Initial:ConcreteState`, `Steps:MicroStep[]`, `Projection:Action[]`, `Complete:boolean`, `CanComplete:boolean` |
| Collision | `Rule:string`, `Left:string`, `Right:string`, `HiddenOriginal:string`, `ExposedOriginals:string[]`, `Observation:string` |
| RetryCertificate | `Witness:string`, `FromState:integer`, `ToState:integer` |
| Claim | `ID:string`, `Status:string`, `Statement:string`, `Evidence:string[]`, `Limits:string[]` |

Empty arrays are `[]`, not null. `Via` is null only for the initial node.
`WorldResult` is an inaccessible historical oracle field, never an actor's
receipt source. It is exported to check authenticity against a concrete
history; concrete stepping/lifting must not use it to enable work or obtain
result bytes. Values are `""` (absent), `"r1"`, or `"r2"`; ACKs have no payload.

## Fixed identities and ordering

`Graphs` and `Groupings` contain exactly 16 entries: chain masks 0 through 7,
then cycle masks 0 through 7. Bit 0 services A, bit 1 B, bit 2 C. Config agents
are A/B/C in that order, all generation 1. Components sort owners and ports by
A/B/C; bindings sort AB/BC/CA. The chain declares exactly AB, BC; the cycle adds
CA. `ConfigID` is `chain/<mask>` or `cycle/<mask>` in decimal.

Bindings are exact, not just labels:

| Original | Caller → callee | Request | Grant | Resource |
| --- | --- | --- | --- | --- |
| AB | A → B | `increment` | `grant:A/g1->B/g1:counter-B` | `counter-B` |
| BC | B → C | `increment` | `grant:B/g1->C/g1:counter-C` | `counter-C` |
| CA | C → A | `increment` | `grant:C/g1->A/g1:counter-A` | `counter-A` |

Every binding has both generations 1 and Operation 1. Equal local operation
numbers and equal result bytes cannot identify an original. An empty binding
(all string fields empty, integers zero) occurs **only** in physical
crash/recovery events and the absent chain CA transport slot.

`State.Stages` slots always mean AB, BC, CA. Present stages are `U`, `W`, `P`,
`D`; the absent chain CA slot is `""`. `ConcreteState.Agents` slots mean A/B/C;
`Originals` slots mean AB/BC/CA. The absent chain CA transport has an empty
binding, all booleans false, and all values empty. It contributes neither debt
nor transitions. There is no wildcard original.

## Graph semantics

Nodes use zero-based array indices as IDs. Node 0 is all U (and absent chain CA),
Parent −1 and Via null. Other Parent values are a BFS predecessor and Via is
its exact Action. Node.Arcs contains **every** outgoing named transition;
retaining a predecessor tree alone is invalid. `To` is a valid node index.
Successors are emitted in original order AB/BC/CA, with at most one progress
action per original at a state. Abstract actions use exactly `commit`, `apply`,
`consume` and the original's full binding.

U→W is commit; W→P is apply if the callee is servicing or its outgoing stage is
neither W nor P; P→D is consume. All other slots remain unchanged. Complete is
all declared originals D. `CanComplete` is reverse reachability from Complete,
not a local absence of outgoing edges, fairness claim, or shortest-path flag.
Graph.States, Arcs, CanComplete and CannotComplete are computed counts.

A dependency is e→f where e's callee is serial and f is that callee's outgoing
original. Dependencies are emitted by original order. Acyclic Order is the
producer's deterministic prerequisite-first ordering, found by repeatedly
scanning originals and adding those whose prerequisites have already been
added. For cyclic graphs Order may be a partial ordering; Acyclic must be false.
The auditor must check edges and order, not accept Acyclic alone.

Groupings retain exact leaves, open AB and BC intermediates, both complete
parenthesizations, and their closed configs. `Equal` is computed from complete
canonical components and separately explored flat/left/right macro graphs.
The full graph is serialized once per config; the auditor reconstructs/checks
the graphs from both closed configs. Algebraic equality is not the concrete
abstraction argument.

## Concrete histories and candidate collisions

Steps stores each event and its **post-state**. Initial precedes step 0.
Physical `crash`/`recover` has a named Actor, empty Original, empty Binding,
empty Result and Duplicate false. Other event kinds are exactly:

`commit`, `send-request`, `admit`, `deduplicate`, `apply`, `send-result`,
`receive-result`, `record`, `consume`, `send-ack`, `release`.

Caller actions are commit/send-request/receive-result/record/consume/send-ack;
callee actions are admit/deduplicate/apply/send-result/release. Every such event
carries its original's exact binding. Duplicate may be true only on
admit/deduplicate/receive-result/release; it retains the transport channel on
delivery. Apply/send-result/receive-result/record/consume carry a result;
other events carry none. The unchanged edge semantics allows retries even
after application/consumption and keeps the application tombstone after release.

Availability phases are `(false,false,0)` initially, `(true,true,0)` after a
crash, and `(false,true,1)` after recovery. Crashing disables all roles of that
physical owner and clears only its outgoing volatile result. Serial waiting
gates only incoming apply, using the actual owner's outgoing committed and
unconsumed facts. Admission and availability remain additional requirements.

Witness names/order:

1. `R0-unpaid`: all-servicing chain, AB W and BC D.
2. `R0-paid`: its continuation, AB D and BC D.
3. `R1-uncommitted`: A-serial/B,C-servicing cycle, AB U, CA admitted W.
4. `R1-waiting`: its continuation committing AB to W.
5. `serial-cycle-control`: all-serial cycle after all three commitments.
6. `retry-divergence`: request send/deduplicate loop after AB admission.
7. `pre-record-loss-and-release`: lost caller volatile receipt, recovery,
   retransmission from retained custody, record/ACK/release/consume.

Projection is the value-abstracted sequence of all commit/apply/consume events,
including internal ones. Witness.Complete and CanComplete classify its final
state. Collision Left/Right name witnesses. R0 compares all physical availability
and complete exposed BC transport equality, but different completion/history;
R1 compares physical availability and complete BC/CA equality but different CA
readiness. R1 does not alone refute a trace relation omitting readiness.

RetryCertificate counts state 0 as Initial and state k as the state after k
steps. Its FromState 3 and ToState 5 in `retry-divergence` are equal concrete
states; the intervening legal send/deduplicate path can repeat forever. This
is an existence witness for concrete divergence, not divergence equivalence.

Claim statuses are `source-argument`, `bounded-exhaustive`,
`bounded-check-and-source-argument`, or `concrete-witnesses`. Statements and
Limits retain the distinction between proof argument, executable checks,
individual witnesses and unestablished claims. These prose claims require
review in addition to structural JSON checking.

## Retained sample, provenance, and command evidence

`evidence/task1/export.json` is the full deterministic sample for Task 2.
`evidence/task1/source-manifest.json` records SHA-256 hashes of the producer,
its tests, this schema, the contract/argument, and the unchanged twoagent source
inputs. It also records the task base commit, tool versions and exact command
argv/exit/log paths. The manifest hashes source inputs, not itself; command logs
and the sample are separate retained files. Intermediate RED/GREEN logs remain
historical evidence of the implementation process, not current outcome claims.

Task 2 owns the independent Python checker and hostile JSON controls, including
forged binding, dropped hidden debt, wrong serial guard/completion, removed/added
arcs, premature release and duplicated application. Task 1 supplies genuine
records and hostile Go tests; it does not preapprove an independent audit.
