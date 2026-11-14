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

## Current status

[The Phase 3B closure assessment](../phase3b-closure/01-assessment.md) closes the bounded investigation against the original gates below. It records contracts, adverse witnesses, implemented fragments and explicit obstructions, without claiming a unified two-agent correctness theorem. Positive cross-invocation delegation, delegated cancellation-after-admission custody and whole-pair split-effect completion remain unestablished. [Phase 3C.1](../phase3c1/04-results.md) now gives a conditional source argument for possible-completion preservation under compatible grouping in a fixed three-agent chain/cycle, with complete small abstract checks, factored transport probes and concrete hiding failures. It retains every named debt, binding and endpoint policy dependency. This does not establish shared-authority fan-out, arbitrary contexts/topologies, runtime refinement or general Phase 3C closure.

[Phase 3C.2](../phase3c2/04-results.md) adds the fixed-chain shared-capacity boundary. Equal occupancy and exact authorization can conceal opposite finite-completion possibilities; ready-now reservation can lose readiness after acquisition. Stable-ready is conditionally sufficient under fixed policies, recoverability and atomic application/result-retention/Q-release, but excludes old AB-before-BC progress traces. Static joins retain one root identity and exact original memberships; they do not union independent dynamic holdings. The 72 checked abstract configurations, synthesized factored probes and ten captured histories have distinct evidence scopes. Permission, cumulative expenditure and reusable capacity remain separate. This is research-model evidence, with no new runtime qualification or general Phase 3C closure.

[Phase 3C.3](../phase3c3/04-results.md) falsifies the sufficiency of a full capacity/custody snapshot when a continuation monitor constrains future ordering. Under exact signatures, common visible synchronization, independent private state and conjunctive marking, equal marked visible languages preserve initial existential possible completion. Equal languages can still admit a doomed residual; the old stable-ready progress result cannot be imported after removing continuation actions. The unchanged whole plant retains all agent service semantics: this is monitor substitution, not complete-agent replacement. Independent reconstruction covers 1,980 finite products, ten exact language comparisons, 1,620 triple regroupings and eight concrete model histories. No runtime refinement, general nonblocking, full abstraction, optimality or verification-cost reduction is claimed. General Phase 3C remains open.

[Phase 3C.4](../phase3c4/04-results.md) retains trace-indexed sets of individual residual marked languages. Their equality preserves every-reachable-state nonblocking and is closed under the declared compatible finite synchronized product; derivative unions are insufficient. The bounded reconstruction covers 2,772 products, 21 monitor comparisons and 2,376 profile-equal substitution rows, including blocking contexts. Generic strictness lies outside the all-marked ALAN roster. An escapable abstract private cycle separates nonblocking from inevitable completion, without establishing its realization by the strongly terminating runtime instruction grammar. This remains monitor substitution over one unchanged K, with no new concrete/runtime qualification. General Phase 3C and Theorem 4.8 remain open.

[Phase 3C.5](../phase3c5/04-results.md) strictly coarsens FULL with J: initial marked language plus inclusion-minimal residual languages. Under the same finite synchronized grammar, J preserves initial possible completion and nonblocking and is closed under product. Finite suffix-hit implication has an exact local coverage characterization, not a global necessity theorem. Plant-relevant strictness, MAX failure, generic MIN loss and J nonnecessity remain distinct controls. Independent reconstruction covers 3,564 products and all 2,772 J-equal substitution rows, including 1,928 blocking rows; both new monitors retain all 324 blocking K/D contexts. Corrected source-only Python qualification follows a preserved successful bytecode-cache attack. No weaker optimal global interface, cheaper verification, runtime refinement or general Phase 3C closure follows. Theorem 4.8 remains open.

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

The [Phase 3B.1 reciprocal-obligation slice](../phase3b/00-charter.md) tests two
actual agents with shared per-agent availability and durable ACK on both edges.
It separates custody from the existence of a completion suffix, compares all
four placements of serial versus servicing policy, and adds a bounded public
same-Store witness. Its [conditional argument](../phase3b/03-results.md) does not
close the rest of Phase 3B.

The [Phase 3B.2 authority slice](../phase3b2/03-results.md) adds asymmetric complete-tuple permissions, rooted original-bound delegation, immutable decisions and fresh-admission validity. It separates consent violations from permission invention, and authority safety from completion reachability. Its runtime witnesses qualify static permissions and retained refusals; a general request-bound delegation primitive was absent at that frozen Phase 3B.2 boundary. [Phase 3B.3](../phase3b3/04-results.md) and [Phase 3B.4](../phase3b4/08-results.md) subsequently implemented scoped one-hop request-bound versions, with their selection-order and identity limitations retained.

Before calling Phase 3B established, attack stable identity across recovery and
new invocation, request/result/consumption binding, grant attenuation and shared
authority conservation, cancellation versus already admitted work, B's own
effect/receipt crash window, and two outstanding obligations in opposite
directions. A waiting cycle is not completion; neither local nor global
quiescence by itself certifies paid application obligations. Identify in-flight
messages and unpaid obligations explicitly, and define distributed quiescence
separately before claiming a local-to-global implication.

Before claiming closure for more agents, define the composition operator and
observable equivalence, show assumptions survive hiding and regrouping, and
test a three-agent cycle and shared-authority fan-out. Two-agent enumeration
alone proves neither associativity nor arbitrary-size correctness. No claim
about reduced verification cost is made until an actual compositional proof
supports it.
