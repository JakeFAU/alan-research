# Phase 3B.5 operational model

The opt-in `S2/SeparateGrantRevocation` configuration tests whether one retained, separately issued grant can support a historical allow while revocation bars every fresh admission on its path. AB uses B's own `B/g1:c` root for `(increment,counter-AB,blue)`. BA uses the original `B/g1:b` root for `(increment,counter-B,blue)` and the request-bound descendant `grant:BA:b`. The two roots name complete distinct action tuples.

## Transitions and preservation sketch

BA commitment enters ordinary request custody without issuing a grant. Live B can retain one issuance attempt after BA commitment. The claim binds B/g1 to A/g1 for BA operation 1, carries only the complete b tuple, names b as parent and root, and forbids redelegation. The retained issuance is accepted when b is current and rejected when b was revoked before the attempt. Once issued, its status is immutable. Live B can revoke b once, including before BA commitment. This transition changes only `Revoked` and leaves the ordinary commitment and delivery paths available.

At BA's first authorization, a covering path requires the exact accepted BA grant and current b. An early refusal persists across grant issuance, revocation and recovery. A historical allow persists across revocation too. Fresh physical admission uses that same selected path with current validity and the serving agent's current activation. No alternate path can repair a refused decision. AB's own c root is current regardless of b's revocation. The inherited serial guard still prevents B from applying AB while B waits to consume its already committed BA result.

Once admitted, the volatile permit is activation-bound. Application requires that permit and consumes it in the same modeled step that creates the effect and durable receipt. This rule admits a BA application after b revocation if A's old permit survives. A crash clears the serving agent's permit. Recovery changes the activation and cannot reconstruct it. A fresh attempt on revoked b then records refusal while preserving the first allow. The event's revocation era describes the source state. The grant and authorization witness keep their historical revision descriptors.

Only application introduces a result. The inherited send, receive, record, consume and release transitions carry the authentic result unchanged. Release requires an acknowledgement sent after durable recording. Thus a crash can clear a volatile result, but an already applied original still has a retained or durably recorded result. Revocation has no custody transition, so either authentic result remains consumable under the existing protocol. The same transition inspection shows that neither payloads nor results can create a grant or a root.

These sketches concern two generation-one agents with one crash each, one original per direction, one b-root revocation, one issuance attempt, one-slot channels, two opaque result alternatives per original, and atomic application plus receipt. The model has no separate debit ledger. The code's first authorization and activation-bound admission are the boundaries checked here.

## Safe completion in the finite graph

For every reachable state in this configuration the candidate condition is

```
Unsafe == 0 && BA.Decision != Refused && (!Revoked || BA.Applied || BA.Permit)
```

Necessity follows from monotone unsafety and immutable first refusal. With b revoked and BA unapplied and unpermitted, neither issuance nor admission can restore a valid permit. The only application rule needs a permit. An incomplete BA original therefore prevents reciprocal completion.

For sufficiency, choose a suffix with no further crashes or revocation and recover either down agent. If BA has a retained application, deliver, record and consume its result as needed. If revoked BA has a permit and has not applied, its server A is live in the permit activation and can apply because A services incoming work while waiting. If b remains current, commit BA and retain the grant where needed. An accepted BA grant precedes any existing allow. With no first refusal, deliver and authorize the original where needed, admit through b, and apply. One-slot channels carry the same original or its result, so an occupied slot can advance the needed delivery. In each case, consuming BA clears B's serial wait. Complete any outstanding AB steps under c and consume its authentic result. This construction exhibits one suffix using the necessary recoveries and deliveries.

Complete breadth-first enumeration and the explorer's reverse reachability agreed with the predicate at every reachable state. AB first refusal was unreachable. Operational checks compared each compact arc's reconstruction to its complete original descriptor and found no mismatch.

| Measure | Count |
| --- | ---: |
| Reachable states | 229,834 |
| Arcs | 1,807,246 |
| Unsafe states | 0 |
| Operational or custody violations | 0 |
| Complete states | 32,768 |
| States with a safe completion suffix | 203,833 |
| Safe states without a completion suffix | 26,001 |
| States in cyclic strongly connected components | 188,614 |
| Cyclic strongly connected components | 75,070 |

The [exported report](evidence/task1/model.json.gz) preserves these shortest prefixes and a retry cycle. State IDs are the report's BFS indices.

| Witness | State | Events | Safe completion suffix |
| --- | ---: | ---: | --- |
| Immediate b revocation | 5 | 1 | No |
| BA first refusal before issuance | 99 | 4 | No |
| Later accepted grant after BA first refusal | 222 | 5 | No |
| b revocation after BA allow, before admission | 445 | 6 | No |
| Fresh BA refusal after historical allow | 783 | 7 | No |
| b revocation after BA permit | 781 | 7 | Yes |
| A crashes with that revoked BA opportunity | 778 | 8 | No |
| BA applies using that revoked-path opportunity | 1258 | 8 | Yes |
| Shortest reciprocal completion | 66302 | 21 | Yes |

The BA application and A crash cuts both start their last step with b revoked and BA's permit present. The former ends in BA's authentic application. The latter has no event obligation because a crash targets physical agent A and clears its BA service opportunity. Witness tests replay both complete descriptors and compare the resulting states. The retained unfair retry witness returns to the same state through send-request and deduplication while the state still has a completion suffix.

## Verification boundary

The historical default authority-service export matched the decompressed Phase 3B.2 archive byte-for-byte. Two opt-in runs also matched byte-for-byte. The safety assertion exited 0 and the universal-completion assertion exited 1 with the same report preserved. The [command evidence](evidence/task1) includes focused pre-implementation failures, successful tests, and SHA-256 manifests.

This result is a source-level sketch and complete operational enumeration of the stated finite relation. The independent ledger, independent reachability audit and concrete runtime mapping are separate qualification work. Whole-Store process loss, repeated physical exposure, non-atomic application and unknown external outcomes require other models. The result does not establish distributed revocation, permission-algebra equivalence, or completion under arbitrary schedules.
