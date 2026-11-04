# Frozen reciprocal-obligation contract

## Two agents and two identities

Agents A and B each have a shared physical availability record: up/down, whether
their one permitted crash occurred, and activation 0/1. There are exactly two
immutable original obligations, d_AB and d_BA. Each binds caller/callee identities,
their original generations, operation sequence, request and fixed grant/resource.
The requests may have equal local sequence numbers and their opaque result values
may coincide; the full edge identity still distinguishes them.

Each edge retains its caller's committed request, received volatile bytes,
durably recorded result and once-only consumption marker; its callee's durable
admission, once-only application tombstone and retained result; and separate
one-slot request, result and identity-only ACK channels. The actual result is an
inaccessible per-edge oracle fact. Neither agent can read an oracle fact to
enable work or reconstruct a deleted receipt.

An A crash disables A's caller actions on d_AB and callee actions on d_BA, erases
only A's volatile received result, and preserves both durable roles and channels.
B is symmetric. Recovery changes physical activation, not logical identity,
consumption, application, grants or the service policy.

## Service contract

Let Waiting(i) mean that agent i has committed its outgoing request and has not
consumed that request's result. Under **serial** service, i may apply an admitted
incoming request only when not Waiting(i). Under **servicing** policy, i may
apply it while Waiting(i). This guard reads only i's own retained outgoing state.

The guard restricts incoming **application only**. Request admission, receipt
transmission, result recording/consumption and ACK processing remain serviceable
while waiting, provided the responsible agent is up. No operation by either role
of a down agent is enabled; channel contents persist. Application and callee
receipt retention are atomic. Retries preserve original identity and never
reapply an admitted/applied request. A caller sends a release-permitting ACK only
after durable recording; receipt bodies and digests are absent from ACKs and
released tombstones. Recorded caller results remain retained after consumption.

Explore four profiles: serial/serial, servicing/serial, serial/servicing, and
servicing/servicing. All have the same channels, fixed grants and recovery rules.
As in Phase 3A, duplicate delivery may retain a channel copy. Keep retry and
duplicate cycles in the reachable graph; there is no exploration depth cutoff.

## Claims to attack

**Local safety and custody.** Every original applies at most once, every selected
result is authentic for that exact edge, every consumption occurs once, committed
facts survive crashes and fixed authority never changes. After each application,
the callee retains its result or the caller has durably recorded it. Exercise
post-application states so this check is not merely vacuous in the initial wait.

**Completion target.** Complete means both original results were durably consumed.
It does not mean that both agents terminated, ACKs drained, receipts were released,
or the system is globally quiescent. CanComplete is existential graph reachability
of Complete; no scheduling probability or inevitable completion is inferred.

**Serial trap.** Define

    Trap = committed_AB AND committed_BA AND NOT applied_AB AND NOT applied_BA.

Attack the sharper conjecture that, in the serial/serial reachable graph,
CannotComplete is exactly Trap. A trace to Trap may precede message delivery.
Retain a separate witness with both agents up, both requests admitted and all
channels empty. Also retain such a wait after both allowed crashes and recoveries.
Do not call any of these a zero-outdegree graph state: original requests may be
retried forever and environment crashes/recovery may still occur.

**Sufficient service.** Test whether either asymmetric profile has CanComplete
at every reachable state. If so, both servicing agents are stronger than necessary
for this two-edge obstruction. Retain an unfair retry cycle in a state that can
complete, to distinguish enabled work postponed forever from an unresolvable
dependency. Eventual progress needs a separate recovery/service premise.

## Source-level argument obligations

Project either edge onto the unchanged Phase 3A durable-ACK model, renaming its
fixed endpoint/obligation identities consistently. The caller's crash/recovery
maps to Phase 3A; callee crash/recovery and opposite-edge work project to stutters.
Shared availability and serial service can remove possible steps, not add invalid
edge steps. Check this on every reachable edge using the actual frozen Phase 3A
transition relation, independently of the new joint transition implementation.
This projects custody safety; removal of enabled work cannot transfer completion
reachability or liveness automatically.

For serial Trap, either first application would require consumption of the other
unapplied original. The causal result chain forbids that, and recovery preserves
the guard facts. For one servicing agent, serving its incoming request lets the
other consume, opening that other's serial gate. Complete this argument from
arbitrary reachable prefixes, including retained receipts and down agents, under
the declared finite recoverable-crash and service assumptions.

## Public-runtime correspondence boundary

Use two ordinary authored roots with separate reciprocal request/reply ports.
Establish both original Sends and admit/acknowledge both requests through public
Store APIs before running. In the serial profile, neither current Receive matches
its admitted peer request. In a servicing profile, one root processes that request
before waiting for its own result. Keep both sources open at a final hold Receive
to separate the obstruction from Phase 1K's Send/End lifetime race.

Measure admitted request observations and actual selected consumption links in
history. Derive unpaid application obligations from the frozen request/reply
contract; a message delivery acknowledgment alone does not pay one. RunReport's
Quiescent, PendingDeliveries and UnresolvedIntents are bookkeeping observations,
not a substitute for this application-obligation projection. Repeated Run and
whole-Store close/open/replay must preserve the serial wait. This is a bounded
control-order witness, not full small-step refinement or independent-owner
recovery, and neither runtime profile performs an external mutation.
