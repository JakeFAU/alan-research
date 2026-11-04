# Phase 3B.1: reciprocal obligations, custody and progress

## Question

If two opposite-direction delegations each obey Phase 3A's durable-ACK contract,
does their composition preserve custody and admit a completion schedule when
each agent can wait on the other? Separate preservation of evidence from service
dependencies: a safe pair can still owe results that neither can produce.

The user merged Phase 3A and authorized continuing the research program. This
first Phase 3B slice investigates mutual obligations and recovery of their
retained control. It does not close Phase 3B's authority or cancellation work.

## Approach

Use an independent finite model with two **actual** agents, not four independent
caller/callee processes. Each has one outgoing obligation and serves the other
agent's incoming obligation. A physical interruption affects both roles of that
agent. Compare serial incoming service with service allowed while an outgoing
result remains unresolved, including the two asymmetric placements.

Add a separate public-API, same-Store program witness. Both ordinary roots first
send requests. Serial programs then wait for their own replies before selecting
the peer request; a servicing program selects and answers the peer request first.
Both end at explicit unsupplied hold inputs. The observation is authenticated
result consumption, not root termination or world-effect completion.

This approach tests the next boundary without adding cross-store transport or
privileged runtime roles. Implementing distributed ownership now would combine
several unresolved assumptions. Merely running more independent one-edge models
would fail to represent the shared two-agent service dependency.

## Deliverables

- A frozen reciprocal transition system, with stable edge identities and shared
  per-agent availability, at most one interruption of each agent, and strong ACK
  on both edges.
- A minimal dependency-trap witness, distinct from a drained admitted-wait
  witness and a retained wait after both agents recover.
- A conditional custody projection argument and an attack on completion: test
  whether a single agent able to serve while waiting suffices for this pair.
- Independent trace/edge checks, exact finite graph evidence, and a bounded
  public-runtime realization of the request-order obstruction.
- Synchronized new manuscript content and explicit remaining obligations.

## Boundaries inherited and still open

The model fixes disjoint grants/resources, one request per direction, two opaque
result values per request, intact durable stores and atomic application-plus-
receipt retention at the callee. It allows asynchronous delay, retries and
duplication. It does not model cancellation, revocation, authority transfer,
storage rollback, non-atomic external effects, unbounded crashes or resource
exhaustion. These are premises or deferred obligations, not newly proved facts.

The runtime witness uses one SQLite Store and ordinary authored programs. Its
whole-Store close/reopen affects both roots; it is not independent-owner or
A-only process-crash qualification. Pure request/reply ordering does not qualify
the model's external application/receipt atomicity assumption.

Prior evidence: [Phase 3A](../phase3a/03-results.md) established the custody
boundary; [Phase 1K](../../../research/communication/README.md) retained different
source-lifetime/ACK counterexamples; [Phase 2D](../phase2d/02-frame-and-completion.md)
separates decision, settlement, publication and resolver completion. This slice
preserves those distinctions and targets a reciprocal application-service wait.
