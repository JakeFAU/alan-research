# Task 1 independent implementation review

## Verdicts

- **Spec compliance: PASS within the frozen abstract-model scope.** No Critical or Important mismatch found against `docs/foundational/phase3a/00-charter.md`, `01-contract.md`, the task brief, or the supplied AGENTS.md requirements.
- **Code quality: PASS.** No Critical or Important implementation defect found. The model, explorer, checker and exporter are small, explicit and cohesive; no production-runtime package or external dependency is introduced.

## Basis

Reviewed the task brief, implementation report, complete task diff, all six `research/twoagent` files and both authoritative scope documents. This review used the Go testing/review guidance. Source was frozen for review. I made no model/source changes and did not duplicate the implementer's reported normal/race/Staticcheck runs. No suspected defect warranted a targeted execution probe.

The transition system fixes exactly A/B, one request/correlation and grant, two application outcomes and at most one A crash. Receipt delivery, durable recording, consumption, ACK transmission and B release are separate transitions. Request retries retain identity; application and consumption guards prevent repeated action. Crash removes only A's volatile receipt and recovery changes activation (`model.go:103-147`, `165-196`).

B's guards read its own retained state and channel occupancy, not A's private state. Its admission/application/transmission/release work remains enabled while A is down (`model.go:149-187`). `WorldResult` is written on application and otherwise merely carried with the state; it supplies neither guards nor agent-visible bytes. Recovery does not consult it. `Recoverable` erases only that field (`model.go:87-91`, `143-177`). Therefore equality of the paired recovered states compares the complete declared recoverable state, not a favorable subset.

The explorer exhausts an expanding BFS queue without an execution/depth cutoff, records edges even when their targets were previously visited, and reconstructs shortest predecessor paths. Reverse reachability starts at A's consumption marker and assigns each predecessor once; it does not require B release or empty channels (`explore.go:55-98`). Custody and recovered absence of all receipt copies are separate predicates (`explore.go:99-122`). The Tarjan implementation correctly distinguishes a nontrivial SCC or singleton self-loop from an acyclic merge. Thus the reported meanings of graph closure, cyclic-state counts, shortest witnesses and existential completion match the implementation. Shortestness is relative to this declared event alphabet.

## Checker independence and negative evidence

`CheckTrace` independently maintains causal facts and queued/retained copies; it does not invoke `Successors`, `Explore`, `Initial` or either safety checker (`oracle.go:5-15`). It rejects a send after release by checking retained bytes rather than consulting historical `actual` (`oracle.go:59-70`, `98-105`). The historical actual result is used to validate authenticity, not to generate a receipt. ACK events are identity-only, and the durable profile checks prior durable recording (`oracle.go:24-30`, `79-97`).

The tests exercise all graph edges appended to independently checked shortest prefixes, both pre/post-recovery receipt cases, genuine resend/deduplication, record/consume crash boundaries, oracle noninterference, B-private-state independence and an unconsumed retry cycle. The 14 malformed trace cases cover the important claimed rejection classes, including wrong value/correlation/grant, early durable ACK, duplicate application/consumption and resend without retained custody. Committed-fact corruption tests also challenge the transition checker. These are substantive checks rather than relying solely on pinned graph totals.

## Limits and uncertainties

The causal checker is a separate implementation of the same declared contract, not an independently derived specification. It validates event legality; it does not compare its final internal ledger with every exported model-state field. The graph-closure test uses the same transition generator as exploration. Source inspection and separate state/edge checks support the bounded conclusions, but these mechanisms do not establish compiler/Store correspondence or immunity to a shared misunderstanding of the contract.

The numerical counts and actual tool outcomes were reported by the implementer; this review did not independently rerun them. The current source and tests support their intended interpretation. There is no remaining concrete implementation blocker in the reviewed slice. Manuscript synchronization and archival evidence remain with the parent task; broader crash, authority, fairness and composition obligations remain explicitly open.
