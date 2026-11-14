# Phase 3C.5 results: completion languages and residual antichains

## Question and conditional result

Can redundant residual completion alternatives be discarded while retaining
initial possible completion and contextual nonblocking under the same product?
Yes, under the [frozen contract](00-contract.md): J consists of the initial
marked language and the inclusion-minimal residual-language family at every
visible word. It strictly coarsens FULL on plant-relevant fixtures and is
closed under the declared finite synchronized product. This is a conditional
[source argument](01-argument.md), supported separately by exact finite
reconstruction, not a mechanized proof or runtime refinement.

Finite pointed marked epsilon-NFAs have exact compatible signatures, a common
visible alphabet, synchronized visible steps, independent private steps and
conjunctive marks. Shared authority, effects and Q stay in one unchanged whole
K. This replaces finite B continuation monitors, not independently executing
agents. No priorities, fairness, inevitability or divergence preservation is
assumed. Nonblocking means every reachable state has some finite marked suffix.

Every member of a finite family contains a minimal member. Therefore minima
preserve all suffix-hit tests; reduction commutes with pairwise intersection.
Cartesian weak reachability and residual intersection transfer this law to the
product. Initial marked languages intersect under product. These facts justify
J closure and preservation of initial CanComplete and nonblocking, with
structural induction only over the declared grammar containing one K.

For finite F and G, every finite H hitting all F members also hits all G members
iff every G member contains an F member. When coverage fails at target K,
choose one word in each source L minus K; their finite set hits every source
and misses K. For empty F, H is empty. Mutual coverage is exactly equality of
minimal families. Finiteness is essential. This is **local** necessity for
finite-suffix hit tests, not necessity for global contextual observations.

## Strictness and preserved falsifications

- `minimal-base` and `minimal-augmented` have equal initial languages and MIN
  at every word but different FULL at epsilon. The additional L4 lies properly
  between L1 and L0. Both differences intersect actual legal core-K traces:
  `commit:BC` and `commit:AB, reserve-capacity:AB, commit:BC`. The legal-trace
  control alone marks every plant state; the paid-completion products are
  unchanged. Each new monitor has 52 nonblocking and 20 blocking K contexts,
  and all 324 K/D contexts remain blocking. The rejected plant-disabled
  alternative remains in [preflight evidence](evidence/preflight/manifest.json).
- MAX discards a dominated bad residual. Release-first and release-first-doomed
  have equal initial language and all-word MAX but opposite core nonblocking.
  Actual healthy/doomed core paths remain in the export.
- Generic MIN loses initial language: a marked looping initial state with a
  private edge to an unmarked looping state and a single unmarked looping
  state both have MIN={empty language} at every word. Their initial languages
  are Sigma* and empty. J distinguishes them. All-marked leaves do not justify
  dropping the initial component for composites with nontrivial K marks.
- Generic J is not necessary: unmarked no-edge X and single-a-edge Y have
  different MIN after a but no marks in any pointed conjunctive-mark product.
  Initial CanComplete and nonblocking are both false in every such context.
  These generic examples are outside the all-marked ALAN roster.
- A local H trie is not a global nonblocking distinguisher. With one a edge
  to a marked state and another to an unmarked dead state, H={a} gives an
  initial hit but a blocking product. Empty H differs from H containing epsilon.

## Exact measured evidence

| Quantity | Count |
| --- | ---: |
| Monitors / plants / controllers | 9 / 72 / 324 |
| Unordered comparisons | 36 |
| Equal initial / FULL / MIN / MAX / J pairs | 13 / 6 / 7 / 13 / 7 |
| Residual comparisons / language graph nodes / profile graph nodes | 1,442 / 5,208 / 146 |
| Products / states / arcs | 3,564 / 45,800 / 59,492 |
| CanComplete states / initially completable products | 22,716 / 1,234 |
| Nonblocking products | 760 |
| J-equal substitution rows | 2,772 |
| Nonblocking / blocking substitution rows | 844 / 1,928 |
| Failed local coverage targets / empty H / maximum H words | 81 / 28 / 3 |
| Independent finite algebra families / ordered family pairs | 256 / 65,536 |

Full epsilon-closed subset graphs, shortest directional witnesses, equivalence
classes, strict dominance, local H source hits/target misses, products and both
true and false substitution outcomes are retained. No bounded-word sampling
or counts-only oracle replaces reconstruction. See the
[accepted final producer capture](evidence/task1/capture-3/artifact-manifest.json),
[measured summary](evidence/task1/measured-summary.json) and
[independent audit](03-independent-audit.md).

## Qualification and limits

Accepted Python evidence is [qualification-2](evidence/task2/qualification-2/results.json):
13 passing tests, nine poisoned-cache boundary controls, two valid CLI
acceptances and 56 JSON/wire rejections. All 59 child source manifests match
nine declared buffers. The independent algebra and product logic is independent
of Go; unchanged predecessor Python descriptor/automaton utilities are reused
with explicit pins, not independently reimplemented a second time.

The initial independent review's successful cached-code substitution is retained
with pre-fix snapshots. It disproved unconditional source-to-execution binding;
it does not show qualification-1 actually ran stale code. The source-only runner
now compiles recorded buffers throughout its declared import chain. Python and
stdlib remain trusted. Historical standalone predecessor entrypoints are not
retroactively repaired. The earlier valid-input serialization rejection is a
separate preserved auditor error. Producer captures 1/2 retain cache warnings;
final captures 3/4 and Task 2's original fresh export have identical semantic bytes.

[Qualification](05-qualification.md) records source/compiler identities and
attributed checks. No new runtime campaign, SQLite opening, physical history,
old exhaustive transport probe, or whole-root suite was run. Ware and Malik's
[primary conflict-preorder paper](https://arxiv.org/abs/1108.0230) is background;
its termination-event/partial-alphabet theorem is not this packet's proof.
No minimal automaton, weakest global interface, optimality or cheaper
verification follows. Full-agent substitution, runtime refinement, dynamic
authority, distributed ownership, effects, exact costs, refusals, divergence,
general Phase 3C closure and Theorem 4.8 remain open.
