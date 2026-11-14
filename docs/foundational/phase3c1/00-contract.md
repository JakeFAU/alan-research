# Phase 3C.1: Contract Preservation Under Composition

## Question and authorization

Can a pair expose a reusable progress contract when a third agent is attached, without losing internal obligations or importing assumptions that the new environment does not satisfy? The user approved this packet after the Phase 3B closure assessment. The result may be a conditional preservation argument, a falsified candidate rule, or a precisely identified remaining obligation; a positive verdict is not predetermined.

Base: `1456a0000a0e14bb3dd105026189717c71ab905b` (merged PR #49). Reuse the isolated research worktree on `codex/phase3c1-contract-preservation`. This is a research-model packet, not a production/runtime change or a new runtime capture campaign.

## Chosen approach and alternatives

Use a four-stage progress abstraction of the unchanged durable-ACK transport rules, with explicit original identities and endpoint service dependencies. Check the complete small abstract graphs, retain concrete transport witnesses, and give a separate finite-path projection/lifting argument. This directly tests information hiding without restricting retransmission using undeclared remote knowledge.

A full three-agent transport-state enumeration would retain unnecessary transport combinations for the first question; no such exhaustive result is required or claimed. A one-copy delivery restriction would change the environmental premises and could accidentally give a sender knowledge of remote receipt state; it is not the selected approach. A three-agent runtime demo would not establish contract preservation. These are design choices within the approved research direction.

## Fixed domain

The only closed topologies qualified here are the chain `A -> B -> C` with originals AB and BC, and the cycle adding CA. Agents have generation 1, one immutable outgoing original at most, one physical availability record each, and a fixed serial/servicing policy. Explore all eight service-policy masks in both topologies, including redundant terminal-C masks. Each original has a full distinct binding: caller, callee, both generations, local operation 1, request, explicit fixed grant and distinct resource. Result values may coincide across originals; identity must not collapse to value or local sequence.

Each agent may crash once and recover; recovery changes activation, not logical identity. Crash disables all of its roles and clears only its volatile caller receipt. Channels and durable data survive. Durable request/admission, application tombstones, callee result custody, caller recording/consumption, and identity-only ACKs follow the existing `research/twoagent` durable-ACK relation, renamed per original. Shared callee availability and the serial service gate restrict enabled steps. No actor reads the external oracle result to enable work or reconstruct a missing receipt.

Application and callee result retention are atomic. Request/result/ACK retransmission and duplicate-delivery variants remain in the concrete model. Release requires an ACK originating after caller durable recording. Receipt handling and consumption remain serviceable while the agent waits on its outgoing original. Resources are disjoint and sufficient; communication is authentic; logical identities, grants and topology do not change.

The concrete adapter may call the unchanged `twoagent.Successors` relation with explicit renaming, actor-availability filtering and the service guard. That reuse is not an independent oracle. Independent evidence must include the separately implemented old trace ledger and a new Python check of the exported abstract graphs and selected concrete histories. No whole three-agent microstate enumeration is implied by these checks.

## Progress contract and observations

For each original e, retain its immutable binding and a stage:

| Stage | Concrete meaning |
| --- | --- |
| U | Not committed. |
| W | Committed, not applied. |
| P | Applied, not consumed. |
| D | Consumed. |

For agent i, `Waiting(i)` holds exactly when its outgoing stage is W or P; an agent without an outgoing original is not waiting. The abstract transitions are U->W (commit), W->P (apply), and P->D (consume). Apply is permitted when its callee is servicing or is not waiting. Complete means every declared original is D, including originals internal to a grouped component.

The observation retains exact original identities and the order of commit/application/consumption events, while abstracting receipt values, transport steps, physical interruption and recovery to progress stutters. A world application remains visible even when both endpoints are inside a component. This is a value-abstracted progress observation, not arbitrary payload or runtime-history equivalence. The abstract service guard is not concrete immediate readiness: admission and physical availability are additional concrete conditions.

## Composition and hiding

Leaf contracts own their named physical agent and its outgoing original. A join is permitted only for disjoint owner sets, unique original identities, compatible endpoint generations/bindings/policies, at most one outgoing original per owner, and distinct resource identities. Reject conflicting bindings, duplicated ownership/originals, unresolved ports at closed-system construction, and resource aliasing. Open ports remain explicit; an absent external agent must never silently count as a non-waiting agent.

Joining retains all named progress registers and node-policy incidence. An original whose endpoints become internal remains latent contract state and remains part of debt accounting. Grouping must not replace two physical agents with one availability record or one service-policy Boolean. Compare the direct flat configuration with `(A join B) join C` and `A join (B join C)` for each fixture/profile after canonical identity alignment. The algebraic join comparison alone is not the substantive theorem: it must be accompanied by the abstraction obligations below.

## Candidate failures to retain

**R0: erase internal originals.** In the all-servicing chain, finish BC while AB remains committed/unapplied, then finish AB without changing BC. The two states must have identical full exposed BC transport state and physical availability, while global Complete and the hidden AB application/consumption history differ. Each state must be reachable by actual concrete steps and satisfy all declared custody/environment premises. Both can complete: this is an information-loss counterexample without a waiting cycle. Report the exact scope of any minimality claim.

**R1: export only hidden-debt presence or a component service bit.** In a cycle with A serial and B/C servicing, keep BC and CA equal, with CA admitted and A up, and compare hidden AB at U versus W. Both carry unpaid AB debt, but A's concrete CA application readiness differs. This refutes that summary for the declared readiness observation; it does not by itself refute a weaker trace relation that omits readiness. Retaining A's named outgoing dependency repairs this missing fact.

**Assumption-compatibility control.** Retain the all-serial cycle after all originals are W. Its inability to complete is not an authority/custody failure. Pairwise conditional service promises cannot be discharged circularly. Classify missing recoverability, split application/retention, erased custody, changed binding or shared-resource depletion as violated premises, rather than calling them counterexamples to a theorem that assumes those facts.

## Positive proof obligations

The target is a proposition for the fixed chain/cycle grammar, not general ALAN congruence:

1. Every reachable concrete transition projects to a named abstract transition or a progress stutter. Exact custody/identity safety is a separate concrete invariant, not encoded by a stage alone.
2. For every abstract transition from q(s), with s a reachable concrete representative, there is a finite legal concrete path realizing it with only progress stutters before the final stage change. Recover necessary endpoints; use retained requests for admission and actual durable custody for receipt delivery/recording/consumption. No other original must advance and no effect may be reapplied. The suffix can choose no further crashes.
3. q preserves Complete. Finite path projection and lifting therefore establish `CanComplete_concrete(s) iff CanComplete_abstract(q(s))` under the stated premises. This is existence of a completion suffix, not inevitable progress, fairness preservation, divergence equivalence, matching physical schedules or global quiescence.
4. Compatible grouping/hiding preserves those obligations by retaining all original registers, bindings and policy dependencies. Identify the allowed context grammar explicitly. A static dependency edge e->f names the outgoing original f of e's serial callee. Acyclicity is a sufficient condition for completion from every reachable state: finish dependencies first, using custody to consume applied originals. Do not assert necessity for a particular state.

The argument must inspect every transition class and both path directions. Executable probes may cover every old single-edge reachable state combined with the service-guard cases, but those are factored checks, not enumeration of the full product. A hidden-debt collision and a grouped trace comparison are not substitutes for the universal representative clause in item 2. Any unproved step must remain labeled as such.

## Evidence and negative controls

Keep the complete small abstract node/arc graphs, including all incoming arcs, reverse reachability classifications, concrete microstep witnesses and their projections, the two failed candidate summaries, compatible grouping comparisons, dependency checks, and a retained concrete retry/divergence witness. Expected results are hypotheses until commands run. Include hostile records for forged edge binding, dropped hidden debt, wrong serial guard, wrong completion classification, removed/added arcs, and invalid receipt release/application duplication. The independent checker must reject these controls for specific reasons without modifying the genuine evidence to make it pass.

Keep historical findings separate from newly executed tests and model checks. Record exact source hashes, command argv/exit codes and output artifacts. Baseline ordinary tests run once against the 70 existing packages while new work proceeds; qualify all added packages separately and verify the old executable inputs remain unchanged. Run affected ordinary/race tests, Staticcheck, root build and module checks as appropriate. Preserve inherited global Staticcheck diagnostics without weakening them. No new dependencies, no retained SQLite opens and no runtime `.envrc` access.

## Prior art and final scope

Abadi and Lamport, *Conjoining Specifications*, TOPLAS 17(3), 1995, pp. 507–534, supplies established composition/hiding and assumption-discharge context. Primary PDF: <https://lamport.azurewebsites.net/pubs/abadi-conjoining.pdf>; publication record: <https://www.microsoft.com/en-us/research/publication/conjoining-specifications/>. The packet's direct fragment argument must stand on its own; do not claim the paper's general theorem was mechanized or all of its hypotheses discharged here. Add a checked bibliography entry only when citing it in the manuscript.

Reconcile the roadmap and affected Markdown/LaTeX manuscript text, bibliography if cited, and rebuild/visually inspect the PDF. Preserve the inherited Markdown-only Section 7 and warnings. Deliver a draft PR with independent task and final reviews; no merge or release.

Shared-authority fan-out, arbitrary contexts/topologies, dynamic identities/grants, cancellation, split effects, independent-owner runtime recovery, reverse runtime refinement and reduced verification-cost claims remain outside this packet. Phase 3C as a whole remains open even if this packet succeeds.
