# Frozen Phase 3C.2 evidence schema: `alan/phase3c2/v1`

The command, run from repository root, is:

```sh
go run ./research/capacitycomposition/cmd/export
```

It accepts no arguments; stdout is one deterministic, indented UTF-8 JSON object
followed by a newline. Errors are written to stderr with exit status 1. It reads
only the declared source files, runs finite exploration, constructs witnesses,
and executes factored probes. It does not read SQLite, `.envrc`, runtime state,
or live endpoints. All object keys below use exact exported Go capitalization.
No fields may be silently omitted, defaulted, coerced, or added by a checker.
Unless explicitly nullable below, objects/arrays must have the stated type.
Counts and indexes are JSON integers, never Booleans or floats. Strings and
Booleans are strict types. Array order is significant. Reject unknown keys,
malformed records, duplicated identifiers/arcs, invalid indexes and missing or
extra roster members. An integer root limit must match the layout exactly.

## Top-level object and provenance

`Evidence` has exactly:

| Key | Type / meaning |
| --- | --- |
| `Schema` | string, exactly `alan/phase3c2/v1` |
| `Sources` | array of `SourceHash`, strictly lexicographically sorted by `Path` |
| `Graphs` | array of 72 complete `Graph` objects |
| `Groupings` | array of 72 `Grouping` objects, same roster order |
| `Witnesses` | array of the 10 `Witness` objects listed below |
| `Probes` | `ProbeSummary`, measured operational probe results |

`SourceHash` has `Path` (repository-relative slash-separated string) and
`SHA256` (64 lowercase hexadecimal characters). The exact source set is
`go.mod`, `go.sum`, `docs/foundational/phase3c2/{00-contract,01-argument,02-evidence-schema}.md`,
and every `.go` file recursively under `research/capacitycomposition`,
`research/contractcomposition`, and `research/twoagent`, including tests and
exporter commands. README/log/export/report files are not in this source set.
Hashes bind the entire packet, including the probe summaries, to those bytes.
They are not an independent validation result. The standalone `source-hashes.txt`
archive repeats these records; it is not an additional field.

## Roster and static types

Roster order is the nested product: layout in `shared-1`, `shared-2`,
`separate-1`; policy in `early`, `ready-now`, `stable-ready`; mask integers 0
through 7. Config ID is `layout/policy/mask`, with decimal mask and no padding.

`Config`: `Base` (`BaseConfig`), `Layout` (string), `Policy` (string), `Roots`
(array of `Root`), `Memberships` (array of `Membership`).

`BaseConfig`: `Topology` (string `chain`), `ServiceMask` (integer 0..7), `Agents`
(array of exactly A,B,C `Agent` records), `Originals` (array of exactly AB,BC
`Binding` records). Agent order maps masks to bits 0,1,2.

`Agent`: `ID` (A/B/C string), `Generation` (integer 1), `Servicing` (Boolean
matching that owner's mask bit).

`Binding` has exactly these fields:

| Key | Type and exact AB / BC values |
| --- | --- |
| `Original` | string `AB` / `BC` |
| `Caller` | string `A` / `B` |
| `CallerGeneration` | integer 1 |
| `Operation` | integer 1 |
| `Callee` | string `B` / `C` |
| `CalleeGeneration` | integer 1 |
| `Request` | string `increment` |
| `Grant` | string `grant:A/g1->B/g1:counter-B` / `grant:B/g1->C/g1:counter-C` |
| `Resource` | string `counter-B` / `counter-C` |

A physical crash/recover event alone uses a zero `Binding`: every string field
empty and every integer field zero. It cannot be used for original events.

`Root`: `ID` (string), `Limit` (integer). `Membership`: `Original` (string),
`Root` (string). For shared-1/shared-2, `Roots` is exactly `[{ID:"Q",Limit:1}]`
or `[{ID:"Q",Limit:2}]`, and memberships are AB→Q, BC→Q in that order.
For separate-1, roots are Q-AB/1 then Q-BC/1; memberships are AB→Q-AB then
BC→Q-BC. Shared references identify one record; copying or summing them is
invalid. These roots are additional to the unchanged effect `Resource` fields.

## Abstract states, actions and complete graphs

`BaseState`: `Stages`, an array of exactly three strings, positions AB,BC,CA.
The first two are U/W/P/D; the third is exactly empty for this chain.

`State`: `Base` (`BaseState`), `Held` (array of exactly two Booleans in AB,BC
order). Only W can hold. Each root's occupancy is the number of true bits
assigned by exact membership, bounded by its limit. Aggregate occupancy is
never a replacement for `Held` and is not an exported state field.

`Action`: `Kind` (one of `commit`, `reserve-capacity`, `apply`, `consume`),
`Original` (AB/BC string), `Binding` (the exact matching binding).

`Graph`: `Config` (`Config`), `Nodes` (array of `Node`), `States`, `Arcs`,
`CanComplete`, `CannotComplete` (nonnegative integers). Counts equal the actual
node/arc/label totals, not estimates.

`Node`: `State` (`State`), `Arcs` (array of `Arc`, including every non-tree arc),
`Complete` (Boolean, both originals D), `CanComplete` (Boolean, reverse
reachability from complete nodes), `Parent` (integer), `Via` (`Action` or null).
Node index is its array position. Index zero is the genuine abstract initial
state with `Parent=-1`, `Via=null`; every other node has its deterministic BFS
parent and action. `Arc`: `Action` (`Action`), `To` (valid node index).

Successor order is old transitions in original order AB,BC (commit, held apply,
or consume as enabled), followed by capacity reservations in AB,BC order. BFS
assigns IDs on first discovery; all arcs are retained. The checker must derive
the exact complete relation, not just walk the exported parent tree or verify
counts. It must independently recompute completion labels.

## Composition records

`BaseComponent`: `Owners`, `Ports` (arrays of `Agent`), `Originals` (array of
`Binding`). All follow named A/B/C or AB/BC order. Ports represent the declared
callee until owned. Every owned A or B has its outgoing original; C has none.

`Component`: `Base` (`BaseComponent`), `Layout` (string), `Policy` (string),
`Roots` (array of exactly the root descriptors referenced by its originals,
sorted by ID), `Memberships` (array of exactly those originals' memberships in
AB/BC order). Empty arrays are `[]`, not null. An owned terminal leaf is not an
empty component; it retains owner C and the layout/policy with empty roots and
memberships. Joining identifies root IDs, never adds limits. Check all input,
intermediate and closed records, not only the Boolean equality field.

`Grouping` has `ConfigID` (string); `A`, `B`, `C`, `AB`, `BC`, `Left`, `Right`
(each a `Component`); `Flat`, `LeftClosed`, `RightClosed` (each a `Config`);
`Equal` (Boolean true). AB is Join(A,B), BC is Join(B,C), Left is Join(AB,C),
Right is Join(A,BC). Closed configs equal Flat exactly; their complete abstract
relations agree. These are static declarations, not unions of local histories.

## Concrete states and events

`ConcreteState`: `Base` (`BaseConcreteState`), `Held` (two Booleans).
`BaseConcreteState`: `Agents` (array of three `Availability` records), `Originals`
(array of three `Transport` records). CA is a fully zero transport record.

`Availability`: `Down`, `Crashed` (Booleans), `Activation` (integer).
Allowed triples are (false,false,0), (true,true,0), (false,true,1). One physical
record controls both roles of B. Crash loses only that owner's outgoing volatile
receipt and retains holdings, admission, durable receipts and channels.

`Transport` keys are:

* `Binding`: exact `Binding`, or zero for absent CA.
* `Committed`, `Admitted`, `Applied`, `Consumed`: Booleans.
* `Volatile`, `Recorded`, `Retained`, `ResultChannel`, `WorldResult`: value
  strings, each empty, `r1`, or `r2`.
* `RequestChannel`, `ACKChannel`: Booleans.

WorldResult records the actual outcome assigned by application for independent
historical audit. It is not accessible to guards or lifting. Witness replay
must check each claimed receipt against its preceding actual apply, not accept
a syntactically valid final receipt as proof. For an absent CA, all flags are
false and all value strings empty.

`Event` has `Kind` (string), `Actor` (A/B/C string), `Original` (AB/BC or empty),
`Binding` (`Binding`), `Result` (empty/r1/r2 string), `Duplicate` (Boolean).
No event activation field is added; activation is bound by the full post-state
and legal physical transition. A forged activation in a post-state is invalid.

Kinds: `commit`, `send-request`, `admit`, `deduplicate`, `reserve-capacity`,
`apply`, `send-result`, `receive-result`, `record`, `consume`, `send-ack`,
`release`, `crash`, `recover`. Actor is caller for commit/send-request/
receive-result/record/consume/send-ack and callee for admit/deduplicate/
reserve-capacity/apply/send-result/release. Physical events use their physical
owner with empty Original, zero Binding, empty Result and Duplicate=false.
Result is the actual value for apply/send-result/receive-result/record/consume;
it is empty for other kinds. Duplicate=true is possible only for admit,
deduplicate, receive-result and receipt release; it preserves the delivered
channel copy. Request retry remains legal after admission or application.
`release` requires an ACK derived from durable recording and releases retained
receipt custody only. Apply requires a held admitted W original, the service
guard and an up callee; it retains its actual result and clears its hold in one
step. There is no standalone capacity-release kind.

`MicroStep`: `Event` (`Event`), `State` (full post-`ConcreteState`).
All unmodified fields must be identical to the predecessor. Invalid actor,
binding, result, capacity, activation or invented post-state changes must fail.

## Genuine witness roster

`Witness`: `Name` (string), `Config` (`Config`), `Initial` (`ConcreteState`),
`Steps` (array of `MicroStep`), `Projection` (array of `Action`), `Complete`,
`CanComplete` (Booleans). Initial must be exactly NewConcrete: empty custody,
empty channels, unheld originals, all owners initial/up, and exact bindings.
Projection includes exactly the nonstuttering resource-aware transitions.
Endpoint labels match the corresponding independently reconstructed abstract
node. No witness starts from a synthesized receipt or probe representative.

In order, all use shared-1 and mask 0:

1. `early-AB-holder`: both W/admitted, AB holds, CanComplete=false.
2. `early-BC-holder`: identical full base to 1, BC holds, CanComplete=true.
3. `early-BC-completed`: exact prefix 2 extended to both D.
4. `ready-now-AB-holder`: reserve AB with BC=U, then commit/admit BC;
   CanComplete=false. Step index 3 reserves; index 4 commits BC and invalidates
   readiness. Indexes are zero-based.
5. `ready-now-BC-holder`: same full base as 4, BC holds, CanComplete=true.
6. `ready-now-BC-completed`: exact prefix 5 extended to both D.
7. `old-AB-before-BC-commit`: early policy applies AB while BC=U. Its old
   progress projection cannot be lifted by stable-ready without BC progress.
8. `early-AB-holder-crashed`: exact prefix 1 plus B crash; hold survives,
   CanComplete=false despite recoverability.
9. `stable-recovery-custody-r1`: stable-ready; BC holder survives C interruption;
   BC then AB genuinely apply r1, duplicate request/result/ACK delivery occurs,
   each caller loses a pre-record volatile receipt, recovers, records actual
   retained custody, releases via durable ACK, then consumes; both D.
10. `stable-recovery-custody-r2`: the corresponding actual r2 history.

For pairs 1/2 and 4/5 compare the **entire** captured Base records and the derived
root totals, with different Held and opposite CanComplete. Completed histories
must extend their named prefixes byte-for-byte at the event/state level. These
are constructive histories, not claims of globally shortest witnesses.

## Probe domain and measured summaries

`ProbeSummary`: `Domain` (string `factored-old-edge-v1`), `OldEdgeStates`,
`Configurations`, `Representatives`, `MicroArcs`, `Lifts`, `LiftMicroSteps`
(nonnegative integers), `Rows` (array of 72 `ProbeRow` in roster order).
`ProbeRow`: `ConfigID` (string), `Representatives`, `MicroArcs`, `Lifts`,
`LiftMicroSteps` (nonnegative integers). Totals sum rows; Configurations equals
72; OldEdgeStates is the measured old durable-ACK state count.

Domain enumeration is explicit:

1. Every config and every reachable resource-aware abstract state in its graph.
2. Each target original AB and BC, in that order.
3. Every state of unchanged twoagent DurableACK BFS with matching target stage;
   a held target also requires old admission. All transport fields and its
   caller's physical phase are retained, including both opaque results.
4. Each of the other two physical owners independently takes all three phases
   (initial, down, recovered) in that order.
5. The other original is synthesized canonically: U has no committed/admitted
   flags; W is committed/admitted without receipt; P additionally has applied,
   retained r1 and WorldResult r1; D additionally has recorded r1 and consumed.
   No volatile receipt/channel is synthesized for that canonical original.
   Held bits are exactly the reachable abstract state's bits.

Counts are enumeration cases, with multiplicity across target selection; they
are not distinct global concrete states or histories. Each representative is
validated, its complete projection compared with the prescribed abstract state,
every concrete successor checked to stutter or match a named full abstract
transition, and **every enabled abstract action** lifted. Every lift microstep
is checked for legal adjacency and full projection (source until final target).
Successor and lift choices are compared after erasing only WorldResult. No
unrelated stage/hold progress is allowed during preparation.

These are source-bound **reported operational results**, not independent probe
re-enumeration and not exhaustive concrete-product coverage. An independent
checker may structurally validate roster/totals/domain and independently replay
the genuine witnesses; it must not imply it reran this probe domain unless it
actually did so. Abstract graph re-enumeration is a separate mandatory check.

## Hostile controls and independent-check boundary

Forged controls are not operational histories. Reject removed non-tree arcs,
wrong completion labels, aggregate-only or changed holders, invented membership,
conflicting same-ID limits, copied pools, capacity overflow, duplicate
reservation, early release, wrong actors/activations, forged bindings/receipts,
post-state invention, malformed schema, omitted configs and Boolean counts.
The independent checker consumes this JSON/specification; it must not call Go
operational guards, successor functions, projection, lifting, explorer or runtime
Replay. Hash agreement and shared-implementation checks are not independence.
