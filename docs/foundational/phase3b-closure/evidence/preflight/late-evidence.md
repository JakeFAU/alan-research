# Phase 3B.5–3B.7 read-only closure reconnaissance

Inspected HEAD: `7ae2c5ddd5e18d2907dd91fcefc896456a25562d`. No tests, captures, checker executions, database opens, runtime edits, historical evidence edits, or commits performed. Counts below are archived results, except explicitly described direct JSON/byte comparisons made during this assessment.

## Original gate mapping and assessment

Owning gate: `docs/foundational/phase3/README.md`, “Progression gates”: attack cancellation versus already admitted work; B's effect/receipt window; two opposite outstanding obligations; distinguish completion from waiting/quiescence; identify unpaid/in-flight work; provide explicit contracts, adverse schedules, and correspondence to an implemented fragment. The README reserves composition operator/equivalence, hiding/regrouping, three-agent cycle, and authority fan-out for 3C. Its distributed-quiescence definition requirement is conditional on claiming a local-to-global implication.

| Original obligation | Late evidence | Assessment and limit |
|---|---|---|
| Cancellation/revocation versus already admitted work | 3B.5 exact b-root revocation; N06 first allow/debit survives but fresh gate refuses without contact; N07 old opportunity applies after revoke and distinct c-root remains usable. 3B.6 models loss clearing old permits without deleting launched packets. | Attacked with finite graph and selected runtime traces. Exact-root revocation is not global cancellation, grant-ID revocation, or composition-root closure. |
| B's effect/receipt crash window | 3B.6 splits launch, ingress, application, response, retention, consumption. Historical applied-unknown loss has application before death and no retained selected result. 3B.7 supplies matched pre-call zero branch. | Addressed for one selected BA service and shared Store owner; does not establish whole-pair split-effect recovery or arbitrary endpoint/retry semantics. |
| Local evidence versus world truth | 3B.6 equal-Store finite witnesses; 3B.7 normal/race full-export and final-report matching pairs with zero versus one application. | Finite observation argument plus source-conditioned concrete experimental witnesses. Endpoint evidence does not itself become authentic local custody. |
| Completion versus waiting/quiescence | 3B.5 exact finite safe-completion-suffix predicate, retained losing states and retry cycles. 3B.7 genuine final reports are locally quiescent with open runs, unpaid/unresolved work and unknown effect outcome. | Positive progress is existential and atomic-domain only. Local quiescence demonstrably does not certify paid completion. No distributed local-to-global implication is asserted. |
| Implemented-fragment correspondence | 3B.5 eight supported cases/twenty captures; 3B.6 twelve supported cuts out of sixteen with pre-admission exclusions; 3B.7 actual normal/race matched pairs. | Meets “correspondence to an implemented fragment” at the stated forward/existential/source-conditioned tier, not reverse simulation or all-runtime-state refinement. |

I found no late-slice evidence that an original Phase 3B exit condition necessarily remains unmet. This is a scoped finding, not the overall closure verdict: earlier identity, invocation, binding, attenuation/conservation and reciprocal-service gates require the other slices. A closure statement can summarize a set of bounded positive/negative results. It cannot claim a unified unrestricted composition theorem. Requiring universal liveness, distributed recovery or arbitrary-agent closure would add requirements beyond the original Phase 3B gate.

## 3B.5: finite atomic reciprocal domain

Primary documents: `docs/foundational/phase3b5/{01-contract,04-correspondence,07-results}.md`.

Archived `evidence/task2/revocation-progress.json.gz` directly reports 229,834 states, 1,807,246 arcs, 32,768 safe complete states, 203,833 states with a safe completion suffix, and zero predicate mismatches. Thus 26,001 reachable safe states lack such a suffix. `evidence/task2/revocation-audit.json.gz` directly reports 229,834 audited prefixes, 1,807,246 edges, zero unsafe prefixes/edges. The exact predicate is:

`Unsafe==0 && BA.Decision!=Refused && (!Revoked || BA.Applied || BA.Permit)`.

Source inspected: `research/authorityservice/model.go` (`GrantRevocationConfig`, `Successors`), `research/authorityservice/oracle.go` selected audit/authority transitions, `research/delegationcorrespondence/progress.go` (`QueryProgress`), and `research/mutualobligations/model.go` (`Complete`, Apply). Important source confirmation: inherited Apply sets Applied, Retained and WorldResult in the same transition; Complete requires both original consumptions. Crash clears the server permit, revoke changes the current b root, and fresh admission checks current selected authority. Progress search uses operational successors for the graph, then separate reverse reachability; this is not a wholly independent graph generator or an enumeration of all cyclic histories.

`docs/foundational/phase3b5/evidence/task3/verdicts.json` records the eight supported N05/N06/N07/P02-v2 final cases across normal/race and twenty supported captures. Both applied-unknown-revoked cases remain raw/physical-valid but unsupported for the atomic projection, with reasons `whole_store_loss` and `applied_without_retained_result`. N07 intermediate projections retain unresolved BA permits; final atomic Apply requires a matched authentic retained result and single physical application. Source `research/delegationcorrespondence/revocation.go` preserves separate raw, physical and support verdicts and exact b/c root mapping.

Assumptions: generation one, one original each way, two opaque results, one-slot channels, at most one recoverable crash per agent, one no-redelegation grant attempt, one exact b-root revoke, A servicing while waiting and B serial, no separate model debit ledger. Runtime side retains exact debits/roots/bindings as checked sidecars. No whole-Store-loss or split-effect projection follows.

## 3B.6: one selected service, split effects and observation

Primary documents: `docs/foundational/phase3b6/{01-contract,04-correspondence,07-results}.md`.

`evidence/task1/summary.json` records 63 states and 104 arcs. `evidence/task2/audit.json.gz` records 63 independently checked prefixes/states, 104 arcs, 33 merged targets/42 extra incoming arcs, 31 Store classes and 52 dossier classes. Archived witness lengths: equal-local zero/one application 4/7; distinct unretained authentic values 7/7; negative stopped/revoked pending-packet cut length 3 with two-event later ingress/application extension; retained-result recovery consumption length 9.

Source inspected: `research/effectobservation/{model,oracle,observations}.go` and `research/delegationcorrespondence/effect_mapping.go`. Apply changes only external application; response/record/consume are separate. LoseOwner clears opportunity/inbox but preserves packet/world/retained result. No import or successful retry transition exists. `AuditGraph` independently interprets descriptors, checks merged ledgers, and exhausts a separately declared descriptor domain for edge completeness. `StoreView` includes protocol-supplied owner stage, not just intrinsically exported database state.

`evidence/task3/verdicts.json` directly records N07 lengths 0/1/7 and loss lengths 4/5/7 in both modes, with two pre-admission cuts per loss dossier outside the model: twelve supported of sixteen total. These reuse producer `5b4662fd682a220ca3065b15e86672b7e56567e7`, tree `44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`, historical `docs/foundational/phase3b4/evidence/prospective-main/raw-{normal,race}.tar.gz`. No new 3B.6 runtime campaign occurred. Source-conditioned causal certificates establish application-before-loss and successful-receipt withholding; they do not establish a total physical order, continuous callback occupancy, or transport liveness. Authentic committed receipt can support a world fact without local result custody. A cut-certified zero is not future nonapplication if a packet remains launched.

Critical compatibility seam: 3B.5's `Applied` also guarantees retained result; 3B.6's world application does not. The latter supplies a counterexample to extending the former completion predicate to a split-effect domain, while leaving the original predicate intact. The 3B.6 model is one service, not the full reciprocal pair. Do not concatenate these results into a whole-pair recovery/completion theorem.

## 3B.7: operational matched pair, richer declared observation

Primary documents: `docs/foundational/phase3b7/{01-contract,03-comparison,05-results}.md`.

`evidence/qualification/corrected-checks.json` records Supported=true for both normal and race, all independent verdicts true, full-export comparisons true at post-open/stopped/reopened/final, final report equality true, zero/applied contact/application counts 0/0 versus 1/1, and unchanged sealed members. Each `evidence/capture/{normal,race}-v1/raw.tar.gz` has 122 regular dossier members. During this reconnaissance I directly read only JSON members in each archive and confirmed byte equality of all four paired `inspect.json` cuts and complete `run.json`; no database members were opened.

Concrete quiescence witness from `zero/run.json` and identical `applied/run.json` in BOTH archives: `Quiescent=true`, `Steps=0`, `Dispatches=1`; both principals in `OpenRuns`; BA external service `['a',1,4,0]` (actual serialized ID uses JSON double quotes) in UnresolvedIntents/UnknownExternal; pending external delivery; roots remain `Pending='members'`, without outcomes. This is directly useful original-gate evidence: local scheduling quiescence coexists with unpaid selected work and opposite application facts.

The pair starts from copied closed seed bytes, runs sequential branches with separate processes/databases at the same literal endpoint URLs, and shares one Store owner for A/B in each branch. Identical local exports do not mean identical complete forensic dossiers. Process/control/endpoint evidence distinguishes the histories. The endpoint is native idempotent; this is not a non-idempotent retry result.

Source inspected: `research/matchedcapture/{check,compare,world,protocol}.go` and selected `producer/phase3b7_matched_test.go.txt`/source excerpts in the correction record. `CheckDirectory` requires separate Seal/Source/Raw/Protocol/World/Model verdicts plus complete structural equality. Comparator retains every field and array order, rejects duplicate keys, and uses exact number tokens. World checker requires complete process-bound ingress/commit cardinalities and stopped/final snapshots. Protocol checker requires selected pre-call and scheduler acknowledgments; zero holds through confirmed SIGKILL/reap, applied releases worker then awaits post-commit ack, both with successful response release after owner death.

Preserved adverse history: `evidence/qualification/initial-checks.json` initially rejects both pairs because checker expected the wrong passive authorization coordinates. `evidence/task2/review-r2/root-cause.json` preserves source trace: runner stamps authority admission ordinal before transport emission. Current `protocol.go` checks ordinal 1 against the selected worker plus nil native InFlight. Corrected reports accepted the unchanged raw pairs; do not present initial captures as initially accepted or erase rejection.

All operational claims assume honest pinned source, complete process/endpoint coverage and unsigned sealed provenance. JSON-only checker relies on frozen procedure for database-to-export connection; no independent SQLite-content/source-binary attestation follows. Equal input implies any deterministic decision restricted to that observation cannot classify application correctly in both witnessed worlds; this is a conditional observation argument, not an unrestricted distributed impossibility theorem.

## Remaining boundaries and nonclaims

- Earlier slices must establish identity/new-invocation, attenuation and mutual obligations; the late slices alone do not cover all original gates.
- No whole-pair split-effect completion predicate, reverse simulation, arbitrary-runtime-state recovery equivalence, general permission-algebra equivalence, guaranteed eventual progress, retry/import authority, independent-owner distributed correctness or arbitrary-endpoint behavior.
- No local-to-global/distributed-quiescence implication. Its absence is not an unmet conditional gate if the closure assessment states this explicitly.
- Three-agent cycle, shared-authority fan-out, composition/equivalence and hiding/regrouping/associativity are 3C work under the original README.
- This assessment did not rerun historical tests, graph enumeration, checkers, or captures. Historical qualification statements remain attributed to archived evidence.
