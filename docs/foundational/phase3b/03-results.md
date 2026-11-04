# Phase 3B.1: custody composes; completion needs service

## Claim under investigation

For two reciprocal obligations, do durable local facts and Phase 3A's durable-ACK
contract imply that both original results can eventually be consumed? The answer
depends on the service contract. Custody is preserved even when both agents enter
a wait from which neither original can be applied.

The local guarantees under attack are safety guarantees: stable original
identity, retained committed state, authentic results, at-most-once application
and consumption, and durable result custody. They do not already promise service
of an incoming request while awaiting an outgoing result. Adding that progress
guarantee to the premises would change the claim. The counterexample does not
falsify an assume/guarantee argument whose service premises actually hold.

## Conditional two-edge result

**Domain.** Consider exactly the transition system in the [contract](01-contract.md):
two agents A and B, two immutable original obligations d_AB and d_BA, fixed disjoint
grants and resources, authentic asynchronous one-slot channels with retries and
duplicate delivery, durable ACK on both edges, intact stores, atomic application
and receipt retention, and at most one recoverable crash of each actual agent.
Crashing an agent disables both its caller and callee roles. The two possible
opaque results on each edge are not available through any oracle operation.

Let `Applied(e)` and `Consumed(e)` be durable facts for an original edge. Let
`Waiting(i)` mean its outgoing request is committed and not consumed. A serial
agent can apply an incoming request only when not Waiting; a servicing agent can
apply it while Waiting. This gate does not suppress admission, result transmission,
recording, consumption or ACK handling by an up agent. Define:

```text
Complete = Consumed(d_AB) AND Consumed(d_BA)
Trap = Committed(d_AB) AND Committed(d_BA)
       AND NOT Applied(d_AB) AND NOT Applied(d_BA)
CanComplete(s) = some finite legal suffix from s reaches Complete
```

**Proposition (conditional, with a non-mechanized proof).** For every reachable
state in this domain:

1. Each applied edge retains its authentic result at its callee or in its caller's
   durable record, in all four service profiles.
2. With two serial agents, `NOT CanComplete(s)` holds exactly when `Trap(s)` holds.
3. With either agent servicing, `CanComplete(s)` holds.

This is a proposition about this declared pair and these four policies. It is
not a theorem for arbitrary programs, general recovery, or arbitrary-size
composition. The finite exploration checks the frozen realization separately
from the argument below.

### Custody by edge projection

Project a joint trace onto either original edge, renaming its fixed identities
consistently to the Phase 3A durable-ACK identities. A crash or recovery of that
edge's caller maps to the corresponding Phase 3A transition and clears only its
volatile receipt. The other agent's crash or recovery, and work on the other
edge, project to stutters. The extra callee-availability and serial-service guards
only remove possible transitions from the projected one-edge relation. They
cannot introduce an invalid edge transition.

The Phase 3A custody induction therefore applies to each projected prefix.
Conjoining the two invariants gives the pair's custody invariant. This is a
local safety argument through two explicit interfaces. It does not require
re-proving custody over all global histories once those projection obligations
are established. Removing possible transitions, however, need not preserve a
path to completion: that is exactly the obstruction below.

### Why the serial trap cannot complete

Consumption requires the original application and an authentic result. In Trap,
neither original was applied, so neither was consumed and both agents are Waiting.
Suppose there were a first subsequent application. Its serial callee would need
to have consumed its own outgoing result; that requires the opposite application
to have happened already, contradicting firstness. Commitment and application
facts survive crashes and recovery, which cannot open either guard. Hence no
completion suffix exists from Trap.

The shortest Trap prefix commits both original requests, before transmission.
The separate admitted-wait witness sends and admits both requests and leaves
both agents up with empty channels. Retrying requests or delivering duplicates
can still occur. This is a completion obstruction, not a zero-outdegree state or
a claim that the transition system has stopped.

### Why there is a completion suffix outside the serial trap

First recover any down agent and choose a suffix with no additional crashes;
both choices are permitted by this failure model. For an applied original,
custody guarantees either a caller record ready for consumption or a callee
receipt that can be transmitted, received and recorded before consumption.
These operations are not blocked by the serial-service guard. Choose deliveries
that remove channel copies; retries need not starve useful steps.

If at least one original has applied, consume its result at its caller. This
opens that caller's serial incoming-service gate. Commit the other original if
needed, send and admit it if needed, and apply it at the now-enabled callee.
Then transmit, record and consume the remaining result. Already completed steps
are omitted. Complete is reached without reapplying either original.

If neither original has applied but the state is not Trap, at least one original
is uncommitted. Keep that agent's outgoing request uncommitted until it serves
the other agent's request, committing and admitting that other request if needed.
The serving agent is not Waiting, so this application is legal. The preceding
case then supplies a completion suffix. These cases exhaust reachable states
outside Trap and prove the converse.

### Why one servicing agent suffices

Let A be servicing and B serial; the symmetric case exchanges their names.
Recover down agents and commit B's outgoing original if needed. A can admit and
apply that request regardless of A's outstanding outgoing obligation. If it was
already applied, custody supplies its result without another application. Let B
record and consume that result. B's serial gate now opens, so A's original can
be committed, admitted, applied at B and consumed at A. If any of these facts
already holds, omit the corresponding step. The same argument also covers two
servicing agents. Thus one servicing agent suffices within this two-edge policy
space; this is not a necessity claim for every possible protocol.

These constructions prove existence of a finite suffix. Inevitable progress
requires eventual recovery and actual service of the needed commit, request,
application, result, record and consumption actions. An enabled useful action
can otherwise be postponed forever by retries. No fairness or probability claim
is inferred from reverse graph reachability.

## Finite-model evidence

The [deterministic export](evidence/model.json) enumerates the complete reachable
graphs without a depth cutoff. Reverse reachability measures existence of a
completion suffix; every repeated-state arc remains in the graph.

| Policy of A / B | States | Arcs | States without a completion suffix | Trap states | States with an applied original |
| --- | ---: | ---: | ---: | ---: | ---: |
| serial / serial | 68,497 | 511,552 | 144 | 144 | 68,272 |
| servicing / serial | 96,721 | 750,112 | 0 | 144 | 96,496 |
| serial / servicing | 96,721 | 750,112 | 0 | 144 | 96,496 |
| servicing / servicing | 96,721 | 751,744 | 0 | 144 | 96,496 |

All four graphs have zero local-state, durable-transition and custody violations,
and 16,384 Complete states each. Trap is merely a state predicate in the servicing
profiles; it is an obstruction only under serial/serial service. There are 45,824
states with a released receipt in serial/serial and 63,232 in each other profile.
Custody is checked per original throughout these post-application states, rather
than inferred from its vacuous truth at a no-application wait.

The three shortest serial witnesses have 2, 6 and 10 events respectively. The
six-event admitted-wait trace is:

```text
A.commit(AB); A.send-request(AB); B.admit(AB);
B.commit(BA); B.send-request(BA); A.admit(BA)
```

Its recovered variant appends A.crash, A.recover, B.crash and B.recover. These
lengths are shortest for the declared predicates from this initial state and in
this transition alphabet, not globally minimal protocol counterexamples. A
separate two-event send-request/deduplicate cycle outside CannotComplete can be
repeated without consumption even though a completion suffix exists.

The independent joint causal ledger checks all 2,763,520 graph arcs and their
derived targets. The separate Phase 3A projection checks both edges on each arc:
5,527,040 transition-or-stutter checks across 358,660 states. Malformed trace
controls challenge edge/actor/generation/grant binding, fabricated receipts,
premature ACK/consumption, repeated application, service while serially blocked,
and both roles of a down agent. Oracle perturbations challenge dependence on
inaccessible actual-result facts. These are independent implementation checks
of the declared model, not runtime refinement or a mechanized proof.

Normal and race tests and Staticcheck pass. The rejected naive completion claim
is preserved as an [intentional exit 1](evidence/naive-completion.stderr), separate
from passing regression tests and apparatus-error exit 2. See the
[implementation report](evidence/model-implementation-report.md),
[review](evidence/model-review.md) and [runbook](04-runbook.md) for exact gates,
commands and evidence identity.

## Public-runtime witness and quiescence

The separate fixture uses two ordinary authored roots and reciprocal request/reply
ports in one Store. Both initial Sends are committed and both requests admitted
and delivery-acknowledged through public APIs. A serial root waits for its own
reply before selecting the peer's request; a servicing root answers the peer
first. Final unsupplied hold Receives keep both roots open, separating this
experiment from Phase 1K's source-closure/ACK lifetime counterexample.

Application completion is derived from actual matching request/reply observations
and their retained consumption links, including identity and value checks.
`RunReport.Quiescent`, `PendingDeliveries` and `UnresolvedIntents` do not define
this application contract. In particular, the runner excludes an ordinary
message from UnresolvedIntents once its delivery acknowledgment is recorded,
even if the receiver has not consumed its admitted observation
([runner bookkeeping](../../../internal/executor/runner.go)).

A quiescent Store can therefore still owe both application replies. This slice
distinguishes quiescence from successful completion; it does not demonstrate a
failure of a separately specified distributed termination-detection algorithm.
Unlike the retry-enabled abstract Trap, a runner with no eligible instruction
can correctly report quiescence at the serial wait.

The normal and race runs each cover the same 16 cases: four service profiles,
two initial Send orders and two initial request-delivery orders. Per case:

| Policy of A / B | Cases per run | Requests consumed | Results consumed | Unpaid originals |
| --- | ---: | ---: | ---: | ---: |
| serial / serial | 4 | 0 | 0 | 2 |
| servicing / serial | 4 | 2 | 2 | 0 |
| serial / servicing | 4 | 2 | 2 | 0 |
| servicing / servicing | 4 | 2 | 2 | 0 |

Every prefix contains two admitted requests and four real commits: each lowered
program initializes its validation flag before Send. Serial runs make zero
further steps or dispatches. Attempts to select either admitted peer request
fail specifically as an ineligible observation and leave retained tables
unchanged. Repeated Run and whole-Store reopen/Run preserve the wait; both
activation epochs advance by one while original identities, program bytes,
sequences, projections and message/history facts survive. Each servicing run
makes 16 further commits and two result dispatches, validates both replies, and
stops at the two explicit holds. Replay agrees with the retained roots.

The [normal](evidence/runtime-normal) and [race](evidence/runtime-race) archives
retain exact programs, provisions, prefix states, inspections, observation/history
links, admissions, reports, refusals and replay. A [separate audit](evidence/runtime-audit.py)
recalculates counts and checks source/target generations, actual consumption
histories, program hashes and retained recovery facts from all 32 artifacts.
[Independent review](evidence/runtime-review.md) also checked the complete raw
correlation chains, response causality and final Receive locations. The two minor
checker findings were fixed and re-reviewed; no blocking finding remains.

Its scope is a retained control-order witness using pure messages. Whole-Store close/open
renews both roots, not one independently owned agent while its peer keeps running.
It does not realize the model's external application/receipt atomicity or establish
a full transition-by-transition refinement. Runtime delivery ACKs are distinct
from the model's application-level receipt-release ACK protocol. The final holds
are deliberate unfulfilled inputs outside the two original reply obligations.

## Repository and manuscript closure

The ordinary root suite exited 0 on unchanged original sources. Final added-package
tests complete [coverage of all 53 current packages](evidence/package-coverage.json):
43 test packages pass and 10 build without tests. Focused model/runtime race and
integration checks remain distinct from that ordinary coverage. The
[verification record](evidence/verification-completion.md) describes the exact
split and exclusions.

The new manuscript section and conclusion are synchronized in Markdown and
LaTeX; the rebuilt 35-page PDF was [rendered and inspected](evidence/manuscript-review.md).
BibTeX and earlier Phase 3A source/evidence are unchanged. The
[whole-change review](evidence/whole-change-review.md) found no actionable issue
after checking sources, raw records, counts, proofs, hashes and manuscript scope.
Automated review agreement is an additional check, not proof of general composition.

## Remaining obligations

Phase 3B.1 establishes only the declared reciprocal-service slice. Dynamic grant
attenuation, shared-authority conservation, cancellation/revocation races,
non-atomic world effects, independent durable owners, new logical invocations,
unbounded failures and shared resource exhaustion remain unqualified here.
The same-Store result relies on the current authored Receive order and available
resource envelope. Both final holds are intentional; no root-termination result
is claimed.

For Phase 3C, edge custody projection provides a possible safety proof pattern,
not associativity or closure. An additional agent can introduce another service
cycle, shared resources or incompatible authority assumptions. Whether local
contracts discharge those dependencies without rechecking the entire global
history remains the research question.
