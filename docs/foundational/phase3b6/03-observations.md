# Phase 3B.6 independent effect-observation audit

The question is whether a historically allowed, subsequently revoked original with a refused fresh gate can distinguish external application from authentic result custody after one shared owner loss. Under the [frozen contract](01-contract.md), the independently interpreted finite graph has equal local observations with zero versus one application and with two distinct unretained results. Complete endpoint evidence distinguishes those cuts but does not supply a locally recoverable result.

## Independent interpretation and audited boundary

`research/effectobservation/oracle.go` implements `AuditTrace([]Event) (Ledger, error)` without calling `Initial`, `Successors`, operational replay, authority/transition/classification helpers, or either operational observation projection. It shares descriptor types and symbols. The ledger separately records packet transfer, contact, world value, volatile inbox, stored/used values, opportunity, death/restart, b revocation, valid c and fresh refusal. Its initial admission constants and complete fixed binding are independently specified. Every event checks actor, incarnation, root, value, every binding field, and irrelevant fields before checking its custody/lifecycle precondition. `Ledger.StoreView()` and `Ledger.DossierView()` derive observations from those facts.

A ledger represents the semantic cut, not every historical event: loss forgets the old unretained inbox, while preserving packet/world truth and durable result custody. Thus histories which operationally merge must yield equal ledgers even when one had an unretained response before loss. The full event prefixes remain in the witnesses. This audit is not an enumeration of all event histories.

`AuditGraph(Graph) (AuditReport, error)` checks each supplied shortest prefix independently, all state coordinates, every outgoing edge, and ledger equality at every merged target. It rejects out-of-range targets, duplicate states/events, disabled or malformed events, unreachable states, and non-shortest or disconnected prefixes. To detect an incomplete graph, it separately exhausts the valid descriptor shapes over the fixed event kinds, incarnations 0/1/2 and absent/v1/v2 values, checks ledger acceptance and requires every accepted event to be present. It does not use operational successor enumeration to decide which edges are missing. Observation projection comparisons are in the graph audit, outside the ledger interpreter. `TestIndependentObservationPartition` additionally compares equality for every state pair under both independent and operational observations.

This is independence of executable decision paths, not independent authorship of the contract or a proof that both implementations cannot share a specification error. No concrete receipt authenticity or cryptographic provenance follows from fixed symbolic values. The ledger API assumes a post-admission initial cut; constructing a `Ledger` struct directly is not an audit of a trace.

## Finite checks actually executed

The [compressed audit](evidence/task2/audit.json.gz) records all witnesses and 20 named semantic controls. The [graph](evidence/task2/graph.json.gz) preserves its executable input. Both commands were run twice and their respective output bytes matched. Go 1.27.4 darwin/arm64 and Staticcheck 2026.2.2 (0.8.2) matched preflight.

| Audited quantity | Result |
| --- | ---: |
| Reachable states and shortest prefixes | 63 each |
| Outgoing arcs, including all merged edges | 104 |
| Targets with more than one incoming arc | 33 |
| Additional incoming arcs beyond one per noninitial state | 42 |
| Independently derived Store observation classes | 31 |
| Independently derived dossier observation classes | 52 |
| Command controls | 9 lawful, 11 hostile; all expected outcomes |
| Compiled independent-ledger semantic mutants | 4 compiled; all killed by targeted assertions |

The evidence is finite model checking and tested implementation agreement. No universal completion, fairness, forward/reverse recovery equivalence or arbitrary distributed composition result follows from these counts.

## Shortest qualifying witnesses

All lengths count events from the post-admission initial cut. Pair selection minimizes the **sum** of the two lengths, then left length, then lexicographic left/right event paths. Individual witnesses minimize length, then event-path order. The fixed order is launch, ingress, apply(v1), apply(v2), response, record, consume, revoke-b, lose-owner, reopen, fresh-refuse; result-bearing v1 precedes v2. Canonical paths are recomputed with breadth-first search, independent of supplied arc order. The shortestness claim concerns this finite graph and these precise predicates.

| Witness | State IDs | Lengths | Qualifying facts |
| --- | --- | --- | --- |
| Equal Store, zero/one application | 25 / 49 | 4 / 7 | Both reopened, b revoked, fresh-refused, no retained or consumed result; world absent versus v1. |
| Equal Store, distinct unretained values | 49 / 54 | 7 / 7 | Same local facts, world v1 versus v2. |
| Negative cut followed by late application | 13 → 22 → 28 | 3 → 4 → 5 | Already stopped/revoked, packet pending ingress; append ingress/application while still stopped. |
| Authentic retained-result recovery | 61 | 9 | Result retained before death; final consume arc leaves a reopened/revoked source (57). |

The last witness is selected from qualifying **arcs**, not the stored shortest prefix of its final state. A shortest prefix to state 61 may consume before death; that prefix would not demonstrate consumption after reopen. The audit's chosen path ends in consumption at incarnation 2. Each witness contains complete descriptors and independently derived views.

## Conditional observation arguments

These arguments concern the frozen one-original, one-gate, one-contact, honest-endpoint and one-shared-loss boundary, with no read/import or successful retry transition.

1. The Store projection omits ingress and world application. In the first pair, equal historical grant/allow/debit/gate, current authority, owner protocol evidence, refusal and local custody coexist with different world facts. Any decision based only on that projection must give the same answer to both members, so it cannot determine which application fact holds for this pair. The second pair similarly prevents determining the unretained opaque result from those local facts.
2. At a declared cut, complete identity-bound endpoint evidence can distinguish these alternatives. That is evidence about the external world. Local record requires authentic response custody in the original live incarnation; positive forensic evidence has no transition that imports custody. Losing the owner clears its inbox and opportunity without clearing a launched packet, ingress or application. A reopened owner cannot accept the original response under either the old or new incarnation descriptor.
3. The negative cut states absence **then**. Its surviving packet can subsequently ingress and apply while the owner is stopped and b is revoked. Negative evidence at the cut therefore does not establish future nonapplication. Missing/partial endpoint evidence is outside the complete-dossier projection and must remain unknown or unsupported, never encoded as zero.
4. A result already retained before death survives and can be consumed by the reopened owner despite revocation. This positive control distinguishes durable custody from absence of custody and ensures the ledger does not obtain rejection coverage merely by rejecting all recovery.

An applied-but-unretained cut is outside Phase 3B.5's atomic application/receipt domain. It is a counterexample to extending that atomic completion predicate into this split domain, not a refutation of the old bounded theorem. W3 and Phase 2B/2C already establish lost-result indistinguishability, provisional absence and separately entitled import within their own boundaries. This packet's contribution is their finite intersection with exact-root revocation and fresh refusal; it claims no new unrestricted impossibility theorem. The graph pair is not a matched pair of sealed v2 runtime captures, and no new runtime experiment was run.

## Falsifiers, reproducibility and remaining work

[Verification commands](evidence/task2/verification.md), [source hashes](evidence/task2/source-sha256.txt), [output hashes](evidence/task2/SHA256SUMS), [mutant manifest](evidence/task2/mutants.json), and the [reproducible mutant runner](evidence/task2/run-mutants.py) bind the checks to their inputs. Genuine RED against compiled stubs, GREEN, complete focused tests and each compiled semantic mutant's failing assertion are preserved as gzip logs. Tests reject malformed identities, missing/irrelevant fields, stale epochs, wrong authentic values, repeated events, invented custody, b/c aliasing, rewritten historical state, missing/duplicate arcs, revived old permit/response and allowed retry after reopen. Lawful controls preserve pre-revocation opportunities, external packet/world truth and retained-result recovery.

Task 2 changes only the new model's independent audit, command, tests and evidence. There are no operational model edits, production runtime changes, dependency changes or old-report changes. The sealed concrete evidence relation, repository-wide qualification, manuscript Markdown/LaTeX reconciliation and final review remain later tasks. Existing manuscript claims were inspected for the old atomic boundary and epistemic precedents; this note does not claim that manuscript integration is already complete.
