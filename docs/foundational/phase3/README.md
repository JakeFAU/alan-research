# Phase 3: the two-agent composition problem

**Question.** Under which interaction contracts do two individually correct,
durable agents form a correct, durable system? Multiple retained computations
and the controlled Phase 2D construction do not answer this general question.

The program starts with two agents, A and B, with retained local state, fixed
capabilities and observable effects. Communication is asynchronous and their
steps interleave independently. Local correctness, edge obligations, system
safety and system progress are separate claims.

| Phase | Research obligation | Exit evidence |
| --- | --- | --- |
| 3A | Define one interrupted delegation and falsify naive composition. | A precise counterexample, impossibility boundary or conditional theorem, with a reproducible independent model. |
| 3B | Test recovery, authority attenuation, cancellation and mutual obligations. | Explicit rely/guarantee contracts; adverse crash and revocation schedules; correspondence to an implemented fragment. |
| 3C | Determine compositional closure. | An induction with explicit interface compatibility and associativity assumptions, or a smallest higher-order counterexample. |

## First slice

[Phase 3A](../phase3a/00-charter.md) isolates responsibility for retaining one
authentic result. If B may release its last result when A acknowledges volatile
delivery, can A recover after interruption and consume that same result?
The alternative acknowledgment certifies A's durable record, before ordinary
consumption. An independent finite model compares these two contracts.

This is an abstract interaction experiment. It neither installs a distributed
runtime nor alleges that the present Store discards receipts. Phase 1K already
retains ordinary-message/ACK lifetime counterexamples. Phase 2B already separates
local drain from external completion, and Phase 2D already separates parent
decision from resolver completion. Those results remain prior evidence.

## Progression gates

Before calling Phase 3B established, attack stable identity across recovery and
new invocation, request/result/consumption binding, grant attenuation and shared
authority conservation, cancellation versus already admitted work, B's own
effect/receipt crash window, and two outstanding obligations in opposite
directions. A waiting cycle is not completion; local quiescence is not global
quiescence. Identify in-flight messages and unpaid obligations explicitly.

Before claiming closure for more agents, define the composition operator and
observable equivalence, show assumptions survive hiding and regrouping, and
test a three-agent cycle and shared-authority fan-out. Two-agent enumeration
alone proves neither associativity nor arbitrary-size correctness. No claim
about reduced verification cost is made until an actual compositional proof
supports it.
