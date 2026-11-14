# Phase 3C.1 semantic reconnaissance

Read-only inspection at HEAD `1456a0000a0e14bb3dd105026189717c71ab905b`, 2026-11-14. No model/test/capture executed; no SQLite opened; no repository file edited. The numerical results below are historical source claims, not new qualification. This document proposes definitions and proof obligations, not an established result.

## Recommendation

Study **preservation of custody and existential completion by a progress contract automaton**, with an explicit hidden-obligation register. Begin with the chain A→B→C, then close CA to obtain A→B→C→A. Keep one immutable outgoing original per node, fixed generations and disjoint fixed resources, durable ACK, exact endpoint identities, serial/servicing policy, recoverable bounded crashes, and atomic application plus callee retention. Do not collapse an AB component into an ordinary single agent with one service-policy bit.

The useful research question is whether AB's exported contract can replace its transport implementation when composed with C, preserving precisely specified observations and CanComplete. Raw product associativity is not the claim. Hiding should remove transport bookkeeping while retaining dependencies, world applications and unpaid original obligations.

## Source facts and boundaries

* `docs/foundational/phase3/README.md`: 3C requires an induction with interface/associativity assumptions or a smallest higher-order counterexample. Hiding/regrouping, a three-agent cycle and shared-authority fan-out remain open.
* `docs/foundational/phase3b/01-contract.md`: serial service gates only incoming Apply, reading the executing node's outgoing Committed && !Consumed. Receipt transport, recording/consumption and ACK remain serviceable. Complete means both original consumptions, not termination, drain or quiescence.
* `docs/foundational/phase3b/03-results.md`: edge projection transfers custody but not progress; published pair graphs contain 68,497–96,721 states and 511,552–751,744 arcs. One servicing endpoint is sufficient for that exact reciprocal pair. These historical counts are not estimates for the new graph.
* `research/mutualobligations/model.go:73`: exact binding; `:97`: one physical availability record per node; `:105`: EdgeState; `:150`: Successors. Apply sets Applied/Retained/WorldResult together. Crash disables both roles and erases only that caller's volatile result.
* `research/mutualobligations/model_test.go:137` and `:156`: projection to the actual frozen `research/twoagent` transition graph. `research/mutualobligations/oracle.go:20`: independent causal ledger. Neither merely asserts that a custody predicate is true.
* `research/mutualobligations/explore.go:113`: exhaustive BFS and reverse reachability; repeated-state arcs are retained. Its uint16 event packing assumes two agents and must not be silently reused for three.
* `docs/foundational/phase3b-closure/01-assessment.md`, “Atomic versus split domains”: the 3B.6/3B.7 split-effect experiments cannot strengthen the atomic pair's completion theorem. Their runtime process-loss evidence is not an independent-owner refinement of this model.
* `paper/alan-dissertation.md`, §§6.9 and final Phase 3B closure discussion (around lines 995–1003 and 1112): current prose carefully limits the custody proof pattern and leaves compositional closure open. New positive claims require matching Markdown/LaTeX treatment; bibliography need not change absent new literature claims.

## A small exact semantic boundary worth attempting

Fix a finite named set E of original obligations. An edge e has caller src(e), callee dst(e), full immutable binding and an inaccessible authentic result fact. For the initial scope each node has at most one outgoing original. Define its four **progress stages**:

* U: uncommitted;
* W: committed and not applied;
* P: applied and not consumed;
* D: consumed.

The concrete causal invariants make these cases exhaustive. Importantly P includes recorded-but-unconsumed, callee-retained, and receipt-in-transit substates; custody, not arbitrary reconstruction, supplies progress between them.

For node i, Waiting(i) iff its unique outgoing edge is W or P. If there is no outgoing edge, Waiting is false. In an all-up macro progress system:

* U(e)→W(e): commit;
* W(e)→P(e): permitted iff dst(e) is servicing or not Waiting(dst(e));
* P(e)→D(e): consume;
* Complete iff all named originals, including hidden internal originals, are D.

The W→P macro includes finite request/admission plumbing before the actual visible Apply. P→D includes finite authentic receipt plumbing. It must **not** redefine external application as a silent action: each W→P carries the edge's observable world-application label. Each P→D carries its exact consumption label. Contract ports retain binding, endpoint policy and original identity; hiding internal communication does not remove world observations or obligations from Complete.

A component AB exports a small automaton whose internal progress stage for AB is latent state, while boundary obligations BC (and later CA) remain ports. Retaining one four-valued hidden progress register is a meaningful abstraction of receipt buffers, ACK copies, result choices, crash bookkeeping and admissions. It is not raw-product associativity, and it does not claim constant-size summaries for arbitrary components.

Keep availability as an explicit physical-node vector in the concrete semantics. The macro system can omit it **only for existential suffix reachability**, because every down node can recover and the suffix may choose no further crashes. It cannot then advertise immediate readiness while a node is down, universal liveness, matching crash traces or a single shared failure domain.

### Main abstraction proof obligations

For q mapping concrete edge states to U/W/P/D:

1. Every concrete step maps to the correctly labelled macro step or a progress stutter. Apply's guard reads the same node's retained outgoing stage. Crashes are stutters at this progress level.
2. For each macro step q(s)→r and each reachable concrete representative s, a finite legal concrete suffix realizes it, with intervening progress stutters. Recover agents first; use retained committed requests and durable custody to drive actual channels. No oracle lookup, forgery or retransmission of an applied effect is permitted.
3. q preserves the declared Complete predicate and the identities/order of applications and consumptions.
4. Interface composition preserves those conditions: exact bindings and port incidence match; shared nodes have one availability/policy record; edges/resources remain disjoint; no context reads deleted transport details or adds resource/fairness constraints absent from the contract.

If both directions hold, concrete CanComplete(s) iff macro CanComplete(q(s)) follows by finite path projection and lifting. This is a tractable proof target. It is not full weak bisimulation of runtime histories, and no runtime relation follows automatically.

## Falsifiable naive rule and genuine hiding failure

Naive rule R0: a safe verified component can hide its internal original obligations, retain only complete external edge states plus a component service bit, and define component completion from its exposed originals; a compatible C can then use that summary without changing completion or application observations.

Attack R0 first in an **all-servicing chain**, with no crashes and all stipulated environment premises true. Commit AB and BC. Complete BC's authentic application, receipt recording and consumption while postponing AB application. Call this concrete state s. From the same history, service and consume AB without changing BC's edge state, obtaining t. Then:

* s and t have identical entire exposed BC state, not merely equal counts;
* all custody, identity, at-most-once and atomicity premises hold in both;
* both have a completion suffix, so this is not the old deadlock;
* Complete(s) is false while Complete(t) is true;
* an AB world application and consumption distinguish their retained obligation histories.

A boundary projection that merges s and t cannot preserve Complete. Treating the internal application as τ also loses a world mutation. This is an abstraction-rule failure **inside** the intended environment, independent of serial circular wait, fairness or failed custody. It should be retained even if the repaired rule passes.

A second attack shows that a hidden-debt Boolean alone is insufficient for immediate service semantics. In the cycle let A be serial and B/C servicing. Compare prefixes with identical BC and CA state, including CA admitted, where hidden AB is U versus W. Hidden AB remains unpaid in both, but A's CA application is enabled in the U prefix and blocked in the W prefix. Both states can ultimately complete. A single pair policy or debt bit does not recover the guard. If the claimed observational equivalence omits immediate readiness, this pair by itself is **not** a trace-equivalence counterexample: hidden work could unblock A. State the observation relation first; use readiness/refusal only if explicitly included, or use visible internal Apply order to distinguish the traces.

R0 cannot be salvaged by calling the first witness an unfair run: it compares finite reachable states and the truth of completion, not inevitable progress. Nor can it be salvaged by silently redefining Complete to mean only external debt. Such a definition may be useful, but would be a different guarantee and must also export the internal debt separately.

## Repaired sufficient progress condition

For a frozen topology, build the static dependency graph on obligations: e depends on f exactly when dst(e) is serial and f is that node's outgoing original. A servicing node introduces no such dependency. This is deliberately conservative: U or D can transiently open a serial gate, but the static graph may still list that dependency.

If this dependency graph is acyclic, a completion suffix exists from every reachable state under the stated custody/recovery/plumbing premises. Proof sketch: recover nodes, choose no additional crashes, process dependencies before dependents; an already applied edge can be consumed by custody, an uncommitted edge can be committed, and an unapplied edge has a serviceable callee after its dependencies are D. Finish each edge once. This is a finite constructive argument, not an inference from an enumerated graph.

* Chain AB, BC: dependency AB→BC when B is serial; C has no outgoing edge. The graph is acyclic for all policies.
* Cycle AB, BC, CA: any servicing node breaks the only static dependency cycle. All serial admits a preserved trapped prefix once all originals are committed and unapplied. That last observation is a regression/control; the new result must be the contract-preserving abstraction and hidden-debt failure, not merely rediscovery of this cycle.

For regrouping, keep exact dependency edges (or a proven sufficient transitive boundary summary) and internal cycles/debts. Internal-path elimination must retain paths between external dependencies; deleting an intermediate node is unsound. Proving that elimination order preserves boundary reachability and internal-cycle flags is a real summary-associativity obligation. Avoid introducing this additional quotient until the four-stage automaton has passed its two-way abstraction obligations.

The acyclic criterion is sufficient, not necessary for individual reachable states: an all-serial cycle with an already paid edge may complete. A blanket claim of necessity would repeat a quantifier error. A parametrized finite-topology proof could eventually support arbitrary finite networks in its exact grammar; three-agent enumeration by itself cannot.

## Manageable finite scope and graph size

The coarse all-up progress graph has at most 4^m states: 16 for the chain, 64 for the cycle. With three availability phases per node retained it is at most 4^m·3^n: 1,728 for m=n=3. Eight policy vectors per topology are cheap in this model. These are combinatorial upper bounds for the proposed abstraction, not measured counts.

A full transport copy of EdgeState has six booleans and five three-valued fields, hence a deliberately loose syntactic upper bound 2^6·3^5=15,552 valuations per edge. Three edges and three three-phase nodes yield approximately 10^14 syntactic states before causal restrictions. Historical reachable two-edge counts are far smaller, but adding the third edge could still make whole-graph materialization expensive. Do not promise a tractable full transport enumeration from pair counts. Inspect/instrument state growth with a declared resource limit if later authorized; reaching that limit must report incomplete exploration, never a passing exhaustive claim.

Suggested minimum deliverable:

1. Frozen mathematical contract and observation relation; retain the all-servicing hiding counterexample.
2. Small independent progress model for chain/cycle and all policy placements, preserving explicit hidden progress registers; independent causal/dependency checker.
3. Source-level two-way reduction argument from the durable transport fragment, with targeted concrete lifted witnesses and edge-custody projection obligations. Do not claim exhaustive three-edge microstate validation if only macro enumeration ran.
4. Repaired conditional composition proposition, clearly separate structural proof from bounded checked instances.
5. Manuscript reconciliation, preserved failed R0 output and independent review. No new runtime campaign needed for this research-only slice.

## Assumption failures versus rule failures

Violated premises: split application/retention, permanent owner loss, erased durable custody, spoofed identities, shared resource depletion, changed generation, actual authority revocation, or a context disabling required receipt service. These do not falsify the atomic fixed-authority theorem.

Genuine failures within premises: hiding unpaid debt; mapping multiple physical nodes to one crash domain; erasing a serial node's outgoing dependency; matching ports by local sequence/result value instead of full binding; treating admitted requests as consumed originals; quotienting away observable internal world applications. Tests should label these separately.

Shared-authority fan-out remains outside this slice's disjoint-resource premise. The roadmap still requires its own attack; do not claim Phase 3C closed solely from chain/cycle work.

## Follow-up: strongest narrow positive claim and grouping

Prefer the unchanged retry-enabled concrete semantics plus the four-stage abstraction to simplifying the protocol merely to reduce enumeration. A one-copy or bounded-delivery experiment can be legitimate, but it adds an environmental premise and must not implement sender omniscience by consulting the remote volatile/recorded state to suppress retries. A full transport BFS is not necessary for a sound structural finite-path abstraction argument.

The strongest narrow claim worth attempting is:

> For the fixed chain/cycle fragment, a named progress-contract automaton preserves all-original completion and the serial application-policy guard exactly, and preserves existential completion by two-way finite-path projection/lifting under intact durable custody and recoverable crashes. Compatible composition and hiding preserve these assertions when original debts, physical-node policy incidence and observable world applications remain represented.

This is deliberately weaker than full behavioral quotienting. Stage W→P macro enablement is **not concrete immediate Apply readiness**: concrete readiness additionally needs durable admission and callee availability. Exact readiness needs those fields in the exported certificate. Dropping result values supports only an explicitly value-abstracted progress observation relation, not equality of arbitrary client-visible value traces or payload correctness. Custody authenticity remains a separate exact-binding invariant proved on the concrete edge projections.

Define grouping as a compatible union of globally named obligation registers and physical-node policy records, with port wiring identifying exactly one original at each matched interface. Internal obligation registers become latent automaton state, never disappear. Shared node records must agree rather than be duplicated. Hiding changes labels for internal protocol communication, not world applications or all-original debt accounting. The substantive obligation is that each grouped automaton remains a valid progress abstraction of its actual transport subsystem in every allowed wiring, and that hiding does not alter that relation. Only after this is shown is associativity of the compatible register join useful. Claiming join associativity alone would be tautological.

A dependency certificate can accompany that automaton: each serial incoming edge names its outgoing prerequisite, and the consumer checks acyclicity after joining interfaces. The constructive ordering proves sufficient CanComplete without enumerating the full transport product. This certificate does grow with retained obligations; no constant-size interface, arbitrary-program closure or reduced verification cost is established.
