# Phase 3A: a finite two-agent result-custody model

This package asks whether preserving each agent's committed state and applying
one original request at most once is sufficient for useful recovery. It is an
independent **abstract contract**, with no production imports or runtime adapter.
It does not expose a current ALAN Store defect or establish distributed Store
refinement. See the [charter](../../docs/foundational/phase3a/00-charter.md) and
[frozen contract](../../docs/foundational/phase3a/01-contract.md).

The delivery-ACK profile falsifies the proposed implication; the durable-ACK
profile preserves custody within this complete finite transition system. A
conditional induction argument below is separate from these graph counts.

## Domain, identities and state

Exactly one immutable original request has identity `FixedDelegation()`: caller
`A`, caller generation 1, operation sequence 1, callee `B`, callee generation 1,
and request `increment-counter`. `FixedGrant` is a fixed valid authorization
premise. Its preservation checks are syntactic checks, **not** authority
attenuation, conservation, cancellation or revocation qualification. No message
creates authority. A recovery changes `Activation` from 0 to 1, preserving the
original request identity; retries are sends of that same original request.

`Initial() State` has this identity/grant, with all other fields zero. State is
comparable, so the complete state, including the inaccessible world fact,
identifies a graph node. The fields are:

| Field | Meaning |
| --- | --- |
| `Delegation`, `Grant` | Immutable original correlation/request and fixed valid grant |
| `Committed` | A's durable original request commitment |
| `ADown`, `Crashed`, `Activation` | A availability, one-crash marker, physical activation 0 or 1 |
| `AVolatile` | Received result bytes erased by A's crash |
| `ARecorded` | Durably recorded result; retained after consumption |
| `AConsumed` | Durable once-only consumption marker |
| `BAdmitted` | Durable admission/deduplication identity |
| `Applied` | Durable application tombstone; retains no result or digest |
| `BResult` | B's optional durable authentic result receipt |
| `RequestChannel` | One-slot authenticated original-request channel |
| `ResultChannel` | One-slot authenticated result channel, containing result bytes |
| `ACKChannel` | One-slot authenticated identity-only ACK; no result or digest |
| `WorldResult` | Inaccessible oracle fact recording the actual application result |

Results are `""` (absent), `r1`, or `r2`. Application nondeterministically chooses
one of the two authentic opaque values and atomically retains it in `BResult`.
Atomic application/receipt recording is an assumption; no external-effect gap is
modeled here. `WorldResult` never enables an agent transition and never supplies
agent-visible result bytes. The event history and report are researcher evidence,
not an extra retained store available to recovery.

`Recoverable(s)` erases **only** `WorldResult`. It preserves all agent, channel,
identity, grant and activation fields, including any volatile bytes. This gives a
strong equality test of the paired recovered-loss states. It is not a general
observation-equivalence definition for a distributed implementation.

## Transition relation

`Successors(Profile, State) []Edge` returns every enabled transition. Each `Edge`
contains an `Event` and target `State`. Events bind `Kind`, `Delegation`, `Grant`,
and (only for value-bearing operations) `Result`. `Duplicate` is legal only on
request admission/deduplication, result receipt and ACK delivery: delivery with
`Duplicate=true` retains a copy in the channel. Repetition permits arbitrary
finite duplication through finite-state cycles, without a queue length cutoff.
Delay is represented by postponing a step; channels do not lose messages.

| Event | Preconditions and state change |
| --- | --- |
| `commit` | A up, uncommitted: retain original request |
| `send-request` | A up, committed, request slot empty: enqueue original request |
| `admit` / `deduplicate` | Request present: first retain admission, or recognize original identity; consume slot unless duplicate |
| `apply` | B admitted, no application tombstone: apply once, atomically retain r1 or r2 and tombstone |
| `send-result` | B retains receipt, result slot empty: copy B's retained bytes to slot |
| `receive-result` | A up, result present: copy slot to volatile receipt; consume slot unless duplicate |
| `record` | A up, volatile receipt, no durable receipt: durably retain those bytes |
| `consume` | A up, durable receipt, no marker: consume that receipt once and retain marker |
| `send-ack` | A up, ACK slot empty: delivery-ACK needs volatile **or durable** receipt; durable-ACK needs durable receipt |
| `release` | ACK present: B deletes retained result, preserves tombstone; consume ACK unless duplicate |
| `crash` | A up, no earlier crash: erase only A volatile receipt; mark A down |
| `recover` | A down: make A available at activation 1 |

B admission, application, result transmission and ACK processing do not inspect
A's private state or availability. They can all run while A is down. There are
no B crashes, new requests, fresh generations, dynamic grants or authority races.

## Complete finite graph and independent checking

`Explore(Profile) Report` performs BFS until the work queue is exhausted. There
is no depth, trace-length, visit or execution cutoff. Every outgoing edge is
retained, including self-loops and edges to already visited states. Reverse
reachability from consumption nodes computes `CanComplete`: some finite service
schedule reaches A's durable authentic once-only consumption. This does not mean
joint agent termination, B receipt release, transport drain or global quiescence. Tarjan SCCs count genuine cyclic states
and components, rather than counting all graph merges as cycles.

`Report.Nodes` exposes the complete graph in Go and is omitted from JSON. Each
node has its state, outgoing event/target edges, BFS parent/event and completion
flag. `Report.Trace(nodeID)` reconstructs a shortest path. JSON contains these
counts and witnesses (field names retain their Go capitalization):

- `States`, `Edges`, `CyclicStates`, `CyclicComponents`;
- `LocalSafetyViolations` (state checks) and `LocalTransitionViolations`
  (committed-fact preservation, single application and single consumption);
- `CustodyViolations`, `CanComplete`, `CannotComplete`;
- `ShortestCustodyViolation`, one shortest `LossWitnesses` entry per authentic
  result, and `Indistinguishable` with both endpoints/traces and exact equality.

`CheckTrace(Profile, []Event) error` is a separately implemented causal ledger.
It calls neither the transition generator nor the explorer/state checker. It
checks committed request sources, admission/deduplication, single application,
retained receipt sources, authentic channel copies, receipt/record/consumption
bindings and profile-specific ACK causality. In particular, a past application
fact cannot justify resending a receipt after deletion. Tests check every graph
edge appended to its source's BFS trace, targeted adversarial traces, local
transition preservation, oracle noninterference and B's independence from A's
private fields. This is independent implementation checking of this contract,
not an independent specification of ALAN's compiler or runtime.

## Result and conditional obligations

The frozen model currently enumerates:

| Profile | States | Edges | Cyclic states / SCCs | Local state / edge violations | Custody violations | Can / cannot complete |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| delivery-ACK | 407 | 1,848 | 361 / 155 | 0 / 0 | 64 | 391 / 16 |
| durable-ACK | 311 | 1,372 | 273 / 115 | 0 / 0 | 0 | 311 / 0 |

Both authentic values have a shortest recovered-loss witness of ten events:

```text
commit; send-request; admit; apply(r); send-result(r); receive-result(r);
send-ack; release; crash; recover
```

The shortest custody failure is the eight-event prefix ending in release.
Custody is already absent while A still has volatile bytes. After crash and
recovery no durable, volatile or in-flight result remains. The two final states
are equal after erasing only the inaccessible oracle fact; identity-only
application/deduplication markers cannot distinguish their results.

These source-independent obligations explain the bounded result:

1. **Authenticity premise:** the only new result originates from B's one
   authorized atomic application; messages and records copy authentic bytes for
   the same original identity. There is no authenticated receipt lookup,
   regeneration, external result storage or transport forgery.
2. **Durable custody induction:** after application establish
   `BResult = r OR ARecorded = r`. All steps except release preserve its retained
   copies. Durable-ACK release is causally after recording the authentic result
   in A's durable store; monotone retention preserves that fact through ACK
   delay/duplication and A's crash. Thus release also preserves custody.
3. **Useful-recovery counterexample:** delivery-ACK allows release before that
   durable recording. The paired recovered-loss states have identical available
   evidence but different actual opaque results. An evidence-respecting recovery
   using only that state cannot guarantee returning the actual result in both
   worlds. Guessing supplies no evidence; withholding remains safe but does not
   satisfy useful recovery. This is a two-world argument under these premises,
   not a new general distributed impossibility theorem.
4. **Existential completion:** if A is down, recover it. If uncommitted, commit;
   if B is not admitted, send/deliver the original request; if unapplied, apply.
   Custody then supplies either A's durable result for consumption or B's result
   for transmission, receipt, recording and consumption. One-slot occupancy is
   serviced by delivery. These choices describe a finite completion schedule;
   the complete graph independently checks its existence at every durable-ACK
   state.
5. **Liveness is additional:** unrestricted schedules can indefinitely repeat
   original-request retries/deduplication while postponing application. Eventual
   recovery and service of enabled request, application, retransmission,
   delivery, recording and consumption obligations are extra assumptions.
   Neither reverse reachability nor custody alone proves inevitable completion.

General Store correspondence, actual process crashes, non-atomic B receipts,
B/simultaneous crashes, cancellation, grant dynamics, opposite-direction waits,
resource fairness, arbitrary agent counts and composition laws remain open.

## Reproduction

From the repository root, using the installed Go 1.27.3 toolchain:

```sh
export GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build GOPROXY=off GOSUMDB=off
go test -count=1 -v ./research/twoagent/...
go test -race -count=1 -v ./research/twoagent/...
STATICCHECK_CACHE=/tmp/alan-phase3a-staticcheck staticcheck ./research/twoagent/...
go build -o /tmp/alan-phase3a-model-export ./research/twoagent/cmd/export
/tmp/alan-phase3a-model-export > /tmp/phase3a-evidence.json
/tmp/alan-phase3a-model-export > /tmp/phase3a-evidence-again.json
cmp /tmp/phase3a-evidence.json /tmp/phase3a-evidence-again.json
/tmp/alan-phase3a-model-export -assert-naive-recovery > /tmp/phase3a-naive-red.json
```

The last command **intentionally exits 1** and explains the falsification on
stderr; its stdout preserves the same complete evidence JSON. It is separate
from the passing regression suite, which asserts that the counterexample exists
and is independently valid. No sockets, paid endpoints, third-party dependencies
or production packages are used. Export order is fixed by profile, successor
order, BFS discovery order and struct fields; map iteration never determines
exported ordering. Counts describe the frozen transition system, not executions
of the current Store.
