# Phase 3B.2: authority attenuation under reciprocal service

**Research question.** Which authority contract lets two durable agents serve one another while waiting, without inventing permissions, and does it preserve completion reachability?

**Result.** Rooted, original-bound, attenuated authority preserves safety in the declared S1/S2 model. Covering authority retained before the first decision and valid at needed admissions also preserves the Phase 3B.1 completion suffix. Missing grants, immutable early refusals and revocation expose safely blocked states. Ambient borrowing violates explicit consent; Cartesian scope mixing separately invents an unheld permission.

**Qualification identity.** Base `1d41d1e20b579443c98a6caa53e51036e55f2d82`; frozen Go/model/runtime-test source `f75e26984c84fe649f26d5e4eaf4ed70e7506e46`. Model, independent checker, race, integration and real process qualification have completed on that source. All six scoped reviews are APPROVE within their stated domains, including explicit re-review after the retained-history correction. The ordinary repository suite passed: 57 packages accounted for, 47 with passing tests and 10 without test files; three optional export tests were skipped. The [whole-change review](evidence/whole-change-review.md) is APPROVE with its final gate conditions now discharged by the ordinary result and manifest verification. Delivery is a draft PR only, with no merge or release.

## Contract and semantic scope

The [frozen contract](01-contract.md) extends the unchanged Phase 3B.1 relation. A services while waiting; B is serial. Both directions retain the original request, operation, authentic result and consumption relationship. There are two original obligations, two opaque results per edge, three one-slot channels per edge and at most one recoverable crash per actual agent. Retries and duplicate delivery remain graph cycles. Application and receipt retention are atomic in this model.

Permissions are complete operation/resource/scope triples:

| Symbol | Operation | Resource | Scope | Initial holder |
| --- | --- | --- | --- | --- |
| a | increment | counter-A | red | A |
| b | increment | counter-B | blue | B |
| r | read | status | shared | A and B |
| x | increment | counter-A | blue | neither |

A permission set is a set of these complete triples. Combining independently held permissions takes their union; attenuation selects a subset along a rooted grant path. Projecting columns and recombining them is a different, unsafe operation: a and b do not confer x. This follows the declared interpretation that each triple independently permits one indivisible operation, rather than separate rights whose conjunction authorizes a new operation. An application with a different policy must define and qualify that policy explicitly. Rooted delegation and authorization intersection have established precedent; RFC 2693 §§6.1–6.3 defines issuer/subject linkage, delegation permission and authorization/validity intersection. This experiment does not implement SPKI or claim that attenuation is new. [Primary source: Experimental RFC 2693](https://datatracker.ietf.org/doc/html/rfc2693#section-6.3).

S1 uses only the serving principal's roots. S2 also accepts an explicit one-hop grant, rooted in the requester, attenuated to a complete tuple and bound to the exact original service, operation and logical generations. S3 adds borrowing from the requester's roots without that explicit consent. S3 is judged against the frozen explicit-consent contract: its missing-grant witness violates principal/consent binding while remaining inside the union of initially held tuples. The separate Cartesian mutant creates x outside that union. These are different failures; the experiment does not prove that every identity-based delegation policy is invalid.

The first authorization decision is retained and immutable. A later grant cannot rewrite a refusal. Every fresh physical admission rechecks the selected path's current validity and binds one opportunity to the serving activation. Later revocation prevents fresh admissions but does not withdraw an already admitted opportunity or undo an effect. Apply consumes the opportunity; server death clears it. Recovery may seek fresh admission while authority remains valid. Only B's b root and its descendants can be revoked in the revocable case; B's r root remains valid.

## Conditional safety argument

**Proposition (rooted authority and custody).** For the declared S1/S2 transition relation, with authentic principals and channels, intact retained state, the fixed roots above, one-hop non-redelegable grants, exact binding checks, immutable decisions and current-validity checks at fresh admission, every application has a retained original authorization and a valid authority path at its physical admission. Its complete permission tuple belongs to the initial global union. Both original result-custody invariants are preserved.

**Proof argument.** Initially the only authority is the fixed root set. Request commitment and transport create no authority except an explicitly represented issuance, which requires an existing issuer root, exact grantee and original binding, and a subset of the parent's complete tuples. Thus grant creation cannot add a tuple outside that root. S1 selects a server root; S2 can additionally select that retained grant. Authorization checks the selected original and records the decision; no subsequent transition replaces it. Admission requires the original allow, the same rooted path still valid, and the current activation. Apply requires that unconsumed opportunity. Revocation and crash cannot add grants or permissions; revocation changes future eligibility and crash discards opportunities. An opportunity admitted before revocation remains justified at its admission boundary. Authentic results are data, not grant constructors. These facts give an induction over transitions for authority safety.

For custody, erase authority bookkeeping and map each step to the actual Phase 3B.1 transition it wraps, or a stutter. The new gates remove applications; they do not add an application or permit consuming a result without the original chain. The earlier custody induction therefore continues to apply. In the explicit b case both services target counter-B: the projection renames the old resource label but does not prove shared-resource value isolation. Results are opaque per-original values, not a model of counter arithmetic or interference. This is a source-level conditional argument, not a mechanized theorem or a proof of the runtime's external-effect atomicity.

## Completion and its limits

**Proposition (positive-case completion suffix).** In the Local cases and S2's Explicit and Shared cases, every reachable state has an authority-safe completion suffix under the declared nonrevoking, inline-grant policy. This is existential reachability, not termination under every scheduler.

**Argument.** Recover down agents and choose no further crashes. All required roots remain valid; the explicit grant, when needed, is retained with commitment before authorization. Consequently neither first authorization nor fresh admission can fail for lack of the required authority. A can serve B's original while waiting. Retained custody supplies a finite path for B to consume that result; this opens B's serial service guard for A's original. Deliver, record and consume A's result, omitting steps already performed. Finite channel occupancy can be cleared by the existing delivery and duplicate-handling transitions. Receipt release occurs only after the caller's durable copy exists. No additional grant or unknown result is needed. Retry cycles can still postpone this chosen suffix indefinitely.

Missing authority, an immutable refusal reached before a late grant, or revocation before necessary admission breaks the premises. Safe blocking is then a legitimate result. Revocation after admission alone does not invalidate already admitted work, but a subsequent crash can discard the opportunity and make current-authority readmission impossible. The finite results must count those states rather than exclude them to recover a universal completion claim.

Within reachable S2/LateGrant states, completion is possible exactly when BA's first decision is not Refused. Within reachable S2/Revocable states,
$$\mathsf{CanComplete}\iff\neg\mathsf{Revoked}\;\lor\!\bigwedge_{e\in\{AB,BA\}}(\mathsf{Applied}_e\lor\mathsf{Permit}_e).$$
If authority has been revoked, every unapplied original needs its own still-live admitted opportunity. Preserving BA's opportunity alone does not reserve completion of AB. The complement is nonterminal and closed under the frozen transitions; the same BA-then-AB construction completes the other states. Missing-authority cases are already noncompletable at the initial state; the shortest late-refusal and revocation obstructions have four and one events respectively. These are bounded characterizations of this relation.

## Exhaustive graph results

| Configuration | States | Arcs | Safe prefix | Unsafe prefix | Violating arcs | Raw completable | Safely completable | Raw blocked |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| S1/Local | 103,361 | 793,736 | 103,361 | 0 | 0 | 103,361 | 103,361 | 0 |
| S1/Explicit | 6,741 | 37,434 | 6,741 | 0 | 0 | 0 | 0 | 6,741 |
| S1/Missing | 6,741 | 37,434 | 6,741 | 0 | 0 | 0 | 0 | 6,741 |
| S1/ScopeMix | 6,741 | 37,434 | 6,741 | 0 | 0 | 0 | 0 | 6,741 |
| S2/Local | 103,361 | 793,736 | 103,361 | 0 | 0 | 103,361 | 103,361 | 0 |
| S2/Explicit | 103,361 | 793,736 | 103,361 | 0 | 0 | 103,361 | 103,361 | 0 |
| S2/Missing | 6,741 | 37,434 | 6,741 | 0 | 0 | 0 | 0 | 6,741 |
| S2/ScopeMix | 6,741 | 37,434 | 6,741 | 0 | 0 | 0 | 0 | 6,741 |
| S3/Local | 103,361 | 793,736 | 103,361 | 0 | 0 | 103,361 | 103,361 | 0 |
| S3/Explicit | 103,361 | 793,736 | 103,361 | 0 | 0 | 103,361 | 103,361 | 0 |
| S3/Missing | 103,361 | 793,736 | 4,815 | 98,546 | 5,776 | 103,361 | 0 | 0 |
| S3/ScopeMix | 6,741 | 37,434 | 6,741 | 0 | 0 | 0 | 0 | 6,741 |
| S2/Shared | 103,361 | 793,736 | 103,361 | 0 | 0 | 103,361 | 103,361 | 0 |
| S2/LateGrant | 111,065 | 840,616 | 111,065 | 0 | 0 | 107,213 | 107,213 | 3,852 |
| S2/Revocable | 220,564 | 1,740,832 | 220,564 | 0 | 0 | 193,681 | 193,681 | 26,883 |
| Cartesian/ScopeMix | 103,361 | 793,736 | 4,815 | 98,546 | 5,776 | 103,361 | 0 | 0 |

The immutable-source rerun reproduced these counts and byte-identical model/checker exports. The sixteen graphs total 1,198,963 states and 9,155,940 arcs; all operational/custody checks report zero violations. A violating arc denotes an unsafe authorization, admission or application, including arcs whose source already has an unsafe prefix. The two unsafe configurations each have 5,776 such arcs. Counts refer to states and arcs, not the infinitely many histories created by retry cycles. Unsafe-prefix bits are monotone, so reaching a terminal custody state cannot erase an earlier authority violation. Raw completion and completion with a safe entire prefix/suffix are reported separately. The unchanged 144-state serial/serial obstruction is an independent baseline regression, not a target for the enlarged state vocabulary.

The minimal decisive BA histories commit, send, admit and authorize the request. S1 refuses the explicit b request; S2 permits its retained rooted grant and refuses the missing-grant counterpart. S3 permits the missing-grant case through ambient borrowing. The Cartesian mutant permits x from mismatched components. Forbidden authorization, physical admission and application are distinct four/five/six-event cuts. The shortest authority violation needs only BA committed; it does not require both agents to be waiting. Reciprocal obligations determine the completion question, not the minimum length of this interface-authority failure.

## Independent validation

Operational enumeration and an independently authored retained-history ledger are separate methods. The independent checker accepted all 1,198,963 generated BFS-prefix ledgers and all 9,155,940 outgoing arcs, including equality of independent ledgers at merged targets. It marked 197,092 prefixes and 1,536,272 arc extensions as retaining unsafe history; these differ from the 11,552 individual locally violating arcs across the two unsafe configurations. This is finite prefix/edge coverage, not enumeration of the infinite cyclic history language. The [compressed ledger export](evidence/audit.json.gz) retains 33 actual histories and four ancestry inputs. Eighteen malformed complete histories preserve both R1 applications while mutating authority, identity or provenance; four valid value pairings preserve original identity even when results coincide. Rootless cycles and return delegation are hostile ledger inputs, not generated valid multihop protocols. The ledger must inspect ancestry, full tuples, original identities, temporal validity and custody; agreement of final counters or outcome values is insufficient. They still share the frozen semantics and event vocabulary, so agreement is not an independent proof that those assumptions describe every implementation.

## Runtime correspondence and obstruction

The runtime has static exact permissions, retained source authorizations, adapter permission decisions, physical dispatch admission and result custody. They are distinct boundaries. It has no general primitive for a request-bound delegator/grantee chain with attenuated resource/scope authority and retained revocation semantics. Static provisioning and specialized receipt-adoption grants do not implement S2. This is a representational boundary in the inspected API, not an impossibility theorem for every ALAN program or future implementation.

Five ordinary fixture scenarios and nine real SIGKILL/fresh-owner cases passed both normal and race-enabled qualification. Local service consumed both originals with one endpoint application each. Missing authority and grant-shaped data retained genuine `grant_denied` refusals, zero contacts and two unpaid originals. After A consumed a real B result, trying to use that result as authority still refused A, leaving one original unpaid. Serial mutual wait retained two unpaid originals. Each crash renewed both Store activation epochs without changing logical identities, permissions or prior retained decisions. Its resources/scopes are encoded by fixed adapter capability configurations, not inferred from arbitrary request payloads. In particular, generic allowed adapter attempts with no retained result remain unknown after reopening; the runtime does not reissue them merely because this test's endpoint saw no contact. That behavior is stricter than the model's atomic Apply/receipt abstraction and can preserve safety while losing modeled completion.

## Real interruption results

| Cut | What survives / observed recovery |
| --- | --- |
| Committed effect before source authorization | Fresh owner authorizes and consumes both original services. |
| Source authorized, no adapter reservation | Fresh owner completes; this cut is not an adapter permission decision. |
| Allowed reservation, before physical admission | Zero endpoint contacts; exact original remains unknown, no result/settlement, both unpaid. |
| Physical admission, before adapter call | Same conservative unknown/unpaid outcome; no retransmission. |
| Result retained, before settlement | Original result settles and is consumed without repeating its application. |
| Settled, before authored consumption | Original result is consumed without repeating its application. |
| Genuine grant refusal | Refusal survives; no contacts and both originals remain unpaid. |
| Mutual wait | Both admitted requests remain unpaid with no contacts. |
| One original consumed, other source authorized | Fresh owner completes the remaining service; one application per original overall. |
| S2 grant retained | Inapplicable: no general request-bound grant primitive in the public runtime. |

The positive fixture changes no database records. Five disposable inspection mutations preserve visible paid markers and result counts while breaking decisive edges; two separate corrupted database copies are rejected by Replay. Replay is a consistency check, not independent compiler validation. The inspection oracle follows original request consumption, effect intent, source decision and principal permission, prepared attempt, result/digest, settlement, reply and actual generation-bound consumption; endpoint records supply separate application evidence.

Independent review **rejected** the first version of this inspection oracle: it matched selected history by subject/sequence but omitted logical generation. Moving six selections from generation 1 to 2 still appeared to pay two originals. The [exact historical witness and pre-fix source](evidence/review4-counterexample/README.md) remain archived. The corrected oracle rejects both that mutation and a final-consumption-only mutation; scoped re-review explicitly marked R4.1 ADDRESSED and APPROVE. This finding concerned the evidence validator, not a demonstrated production-runtime authorization defect.

## Evidence and remaining obligations

The [runbook](04-runbook.md) gives reproduction commands; [command records](evidence/commands.json), [model summary](evidence/model-summary.json), compressed model/ledger exports and raw runtime archives preserve the evidence. Each normal/race archive has 519 original raw files plus its verbose log and command record. The post-run archive checker reproduces five cases, nine checked kills, five inspection mutations, two Replay corruption controls and 46 read-only Replay records per run, and reconciles 115 local SQLite backups per run with the corresponding retained inspections. These are archive checks, not a third general semantics proof. The frozen-source manifest covers 563 tracked Go/module files, including historical and separate-module source; listing does not imply those separate modules were tested. The root suite ran for 1,511.918 seconds with zero failures; its JSON log and exact optional-export skips are retained in the package-coverage artifact. Focused model/fixture Staticcheck and ordinary durable-test Staticcheck passed. Integration-tagged durable-test Staticcheck returned six diagnostics in unchanged older files; the same command on the merged base reproduced exactly those six diagnostics. This gate is explicitly non-green, with no suppression or unrelated fix.

Separate reviews cover [authority](evidence/review-1-authority.md), [complete graphs and minimality](evidence/review-2-model.md), [safety versus progress](evidence/review-3-progress.md), [retained evidence and correction](evidence/review-4-retained.md), and [runtime correspondence](evidence/review-5-runtime.md). The [manuscript/reproducibility review](evidence/review-6-manuscript.md) is APPROVE, with a baseline-command source-attribution correction preserved. The [whole-change review](evidence/whole-change-review.md) found no further source correction; final verification closes its ordinary-suite and manifest conditions. Review agreement is not experimental proof.

This slice does not establish arbitrary multihop delegation, independently owned stores, distributed revocation, hostile storage, generic unknown-effect recovery, shared-resource isolation, general cancellation, fair-scheduler termination or arbitrary-agent closure. A narrowly justified next implementation question is whether an explicit original-bound grant can be represented and checked at both authorization and fresh physical admission without weakening unknown-outcome handling. That question is not silently answered by adding a broad authority framework here.
