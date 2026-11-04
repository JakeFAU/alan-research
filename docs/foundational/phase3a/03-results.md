# Phase 3A argument and evidence boundary

## Claim under investigation

The naive implication is: A preserves its committed local state; B preserves its
application tombstone and applies a delegated request at most once; therefore A
can recover after interruption and consume B's actual result. This implication
omits a contract about who must retain the result during transfer.

The local contracts here are deliberately weak and explicit. They are not a
definition of every possible individually correct agent. In particular, B may
release a result on a transport-delivery ACK, and A promises durability only for
committed input. Calling either one already correct for the global useful-recovery
obligation would assume the disputed composition condition.

## Counterexample and indistinguishability argument

For one fixed delegation d, B applies the authorized work once and retains an
opaque authentic result r. A receives r in volatile memory and sends delivery
ACK. B receives the ACK and releases its receipt, retaining a deduplication
tombstone that identifies d but contains neither r nor a digest of r. With no
second result in flight, A crashes before recording r. Recovery preserves d and
the tombstone but cannot recover r. Resending the original request does not
reapply B's work and cannot reconstruct the released result.

Both local contracts hold. The effect happened once, the identities remain
stable, and no committed A state was lost. Useful recovery nevertheless fails.
Withholding the result is safe; manufacturing it or initiating another mutation
does not satisfy the stated obligation.

Choose two possible result values r1 != r2 and execute the same schedule. After
recovery, the agents' retained and volatile state and the channel contents agree;
only an inaccessible world fact differs. No recovery procedure restricted to
these observations can guarantee the actual authentic result in both worlds.
The argument requires opaque non-reconstructible values and no new observation
channel. An endpoint lookup, another archive, or a B tombstone retaining the
result changes the premises. This is a small indistinguishability argument for
the frozen contract, not a new general distributed impossibility theorem.

## Conditional custody lemma

For the fixed d and actual result r, assume:

1. B establishes a durable copy D_B = r atomically with its one application.
2. A records only an authenticated result for d; a commit D_A = r survives A's
   interruption. Consumption is once-only and retained, and D_A is not deleted.
3. A sends the ACK that permits release only after committing D_A = r.
4. B releases D_B only on that authentic, d-bound ACK. Messages do not confer
   authority; there is no rollback or other receipt deletion.

Then every reachable prefix after application satisfies D_B = r or D_A = r.
This ACK discipline is sufficient under these premises; it is not claimed to be
the only possible protocol. Loss of this invariant can precede unrecoverability:
A may still hold a volatile receipt. The negative recovery witness additionally
requires losing every accessible copy, including any result in flight.

**Proof by transition cases.** Application establishes the left disjunct. Sending,
delivery, repeated original requests and A's interruption/recovery do not delete
durable receipts. A's record establishes the right disjunct, and its consumption
preserves it. Before B can release, the ACK's causal precondition establishes
the right disjunct, which persists. Thus release preserves the invariant.
These are local preservation properties connected by an edge precondition;
they do not require examining a complete execution history when applying the
lemma to a component that already meets those contracts.

The lemma establishes custody safety. The finite graph additionally checks
whether a completion suffix exists from each reachable state. Such an
existential suffix does not establish inevitable termination. Eventual recovery,
service of request/application, receipt retransmission/delivery, record and
consumption, adequate storage, and an eligible live continuation are additional
premises for useful progress. A permanently delayed agent may remain pending.

## Relationship to existing ALAN evidence

The existing Store is stronger at local message admission than the weak model:
[`admit`](../../../internal/executor/messages.go) inserts the accepted observation
in the admission transaction, and [`Commit`](../../../internal/executor/store.go)
updates consumption, successor state, history and intents in one transaction.
The Store also excludes a second owner and renews open runs during recovery.
A process kill of that owner is not the modeled A-only interruption with a
separately executing B owner. No current-runtime defect is claimed.

Phase 1K already preserves source-closure and ACK-service counterexamples.
Phase 2B already distinguishes local drain from delayed external application.
Phase 2D already separates parent decision, original settlement, publication
and resolver completion. This experiment instead attacks B's voluntary release
of a receipt that was already durable. Its assumed atomic application/record
step does not remove Phase 2's external effect/receipt impossibility boundary.

The distinction between delivery acknowledgment and application responsibility
is established systems knowledge; see Saltzer, Reed and Clark's discussion of
[delivery guarantees](https://web.mit.edu/6.1800/www/readings/papers/endtoend.pdf)
(1984, section "Delivery guarantees", fifth PDF page). The contribution here
is its explicit use as a falsifiable ALAN interaction contract, not novelty of
the underlying principle. The primary paper and publication metadata were
checked on 2026-11-04 UTC.

## Qualification status

The [deterministic export](evidence/model.json) enumerates the entire reachable
graph of each frozen profile, including retry/duplicate cycles, without a depth
cutoff. These counts are **model evidence**, not process-crash measurements.

| Profile | States | Edges | Local state / transition violations | Custody violations | States without a path to A consumption |
| --- | ---: | ---: | ---: | ---: | ---: |
| Delivery ACK | 407 | 1,848 | 0 / 0 | 64 | 16 |
| Durable ACK | 311 | 1,372 | 0 / 0 | 0 | 0 |

Both result values have a shortest ten-event recovered-loss witness:

```text
A.commit; A.send-request; B.admit; B.apply(r); B.send-result(r);
A.receive-result(r); A.send-ack; B.release; A.crash; A.recover
```

The shortest custody failure is its eight-event prefix ending in release. Both
recovered-loss endpoints have no completion path and exactly equal recoverable
states after erasing only the inaccessible `WorldResult`. These are shortest
paths from this initial state to these predicates in this transition alphabet;
no globally minimal protocol counterexample is claimed.

Normal and race tests and Staticcheck passed on the new model. The independent
trace checker validates every reachable edge appended to its BFS predecessor
trace and rejects fourteen malformed trace cases. Tests separately check durable
facts across edges, B's independence from A's private state, and noninterference
of the oracle fact with agent transitions. A concrete unfair retry cycle shows
why completion reachability is not inevitable progress. The rejected naive
hypothesis is reproducible as an intentional exit 1, outside the passing suite.

The [argument review](evidence/argument-review.md) found no blocking logical
defect in the scoped induction/indistinguishability arguments and identified the
qualification boundaries above. [Implementation review](evidence/model-review.md)
and [whole-change review](evidence/whole-change-review.md) found no blocking
defects. The latter's pending repository verification gate was subsequently
[closed](evidence/verification-completion.md): the original suite and new-model
tests cover all 50 current root-module packages, with 40 passing test packages
and 10 built without tests. See the [runbook](04-runbook.md),
[implementation report](evidence/model-implementation-report.md), source hashes
and command records for evidence provenance and exact checks.

## Remaining obligations

Phase 3A begins with this one edge. A process-level correspondence for two durable
owners, B interruption, dynamic delegation and capability attenuation, revocation
versus completed effects, mutually outstanding obligations, global quiescence,
shared resource conservation and arbitrary-size closure remain open. The custody
lemma alone does not solve any of those obligations or establish associativity.
