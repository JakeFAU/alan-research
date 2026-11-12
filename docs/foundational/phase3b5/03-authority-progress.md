# Phase 3B.5 independent authority and progress audit

The question is whether the separately retained BA grant preserves historical authorization while exact b-root revocation blocks fresh admission, leaves a prior opportunity usable until its server crashes, and leaves AB's distinct c root current. Task 2 qualifies this question within the frozen two-agent, generation-one finite model. It does not establish a runtime correspondence or a distributed revocation theorem.

## Independent ledger boundary

`AuditTrace` and `AuditGraph` accept exactly the new opt-in `S2/SeparateGrantRevocation` in addition to their historical configurations. The ledger declares the AB binding `(increment,counter-AB,blue)` and B-owned `B/g1:c` itself. It does not call the operational binding, authority selection, authority validity, grant construction or successor predicates when interpreting an event. The graph enumerator and corpus generator use operational successors to supply events; they are not independent implementations of the graph. The ledger shares descriptor types and enum values, not their operational derivations.

The ledger checks one frozen `grant:BA:b` claim, its complete binding and tuple, identity adjacency and generation, exact ancestor/root, revision and no-redelegation condition. Rejected issuance must still have the correct frozen structural claim. Acceptance is evaluated separately against issuance-time current authority. Revocation can occur once, by live B, in the correct event era, without an original or authority payload. It changes only current b validity. It leaves grants, decisions, selected paths, custody and opportunities intact.

First decisions remain immutable. Fresh admission checks the original historical path against current roots; application checks its activation-bound opportunity without repeating that current-root test. Thus an admitted BA opportunity may apply after revocation. A's crash removes that opportunity; recovery cannot recreate it. Authentic results remain recordable and consumable. The c root can authorize and admit AB after b revocation, but B's inherited serial service guard still prevents AB application while B waits for BA consumption.

The complete [independent audit and hostile corpus](evidence/task2/revocation-audit.json.gz) contains 229,834 checked shortest prefixes and 1,807,246 checked outgoing arcs, with zero unsafe prefixes or edges. At each merged graph target the independent ledgers agree, including originals, agents, immutable grants, revocation and violation state. Paid-original prefix counts are 91,520 for AB and 82,304 for BA. These are finite quotient checks, not enumeration of every cyclic event history.

The 34 retained history controls comprise seven lawful controls and 27 hostile modifications. They discriminate historical first allow from fresh refusal, rejected issuance, sticky first refusal, b/c scope, grant identity/ancestry/full binding/tuple/generations, revocation actor/era, stale activation, and authentic post-revocation reciprocal consumption. Positive old-permit and current-path recovery controls detect an overrestrictive checker.

Four disposable-source independent-ledger mutants compiled successfully. Each exited 1 at its targeted test assertion, with no panic or compilation failure counted as a kill: omitted current-root admission check; c revoked together with b; issuance accepted after revocation; and permit retained across crash. The [mutant manifest](evidence/task2/mutants.json) records exact substitutions, copied-source hashes, commands and exits; each build and failure stream is retained separately. The issuance mutant is killed by accepting the forged accepted issuance; the c and crash mutants are killed by rejecting the corresponding lawful controls.

## Exact finite safe-suffix predicate

For every reachable state `s` in this configuration, the independently computed reverse-reachability set agrees with:

```
s.Unsafe == 0 && s.Authority[BA].Decision != Refused &&
    (!s.Revoked || s.Base.Obligations[BA].Applied || s.Authority[BA].Permit)
```

This means existence of a finite safe suffix that consumes both authentic originals. It does not require fairness, promise completion of every execution, or quantify over arbitrary manually constructed states. The suffix can choose not to crash or revoke again. The model includes at most one recoverable crash per agent, one retained original per direction, one grant attempt, one exact-root revocation, one-slot channels and atomic effect application plus receipt. No debit ledger, unknown effect outcome, whole-Store loss, regrant, multihop or additional agent is modeled.

The [independent progress report](evidence/task2/revocation-progress.json.gz) enumerates operational successors afresh and computes its own incoming-edge reverse search. It does not use the operational explorer's completion flags or SCC function. Every reachable state is compared with the candidate predicate. BA first-refusal preservation is checked on every outgoing arc whose source has that refusal; AB refusal is counted over all reachable states. All retained prefix/suffix witnesses are independently ledger-audited. An initial mismatch would be retained as `predicate-mismatch`; none occurred, so no candidate revision was made.

| Quantity | Count |
| --- | ---: |
| States / arcs | 229,834 / 1,807,246 |
| Safe complete states | 32,768 |
| States with a safe completion suffix | 203,833 |
| States without a safe completion suffix | 26,001 |
| Predicate mismatches | 0 |
| AB first-refusal states | 0 |
| Arcs preserving an existing BA first refusal | 57,118 |
| Maximum shortest completion suffix | 23 steps |
| Incomplete dead ends | 0 |
| States in cyclic SCCs / cyclic SCCs | 188,614 / 75,070 |

SCC counts come from an independent two-pass iterative DFS. A component is cyclic if it has multiple states or a self-loop. Retry SCCs preserve the distinction between possible progress and inevitable progress; the historical operational retry witness remains in Task 1's report. Absence of an incomplete dead end does not imply completion: losing states can still retry forever.

The deterministic BFS state JSONL hash is `0e9c40aceeaae702b576026164a64b851f9aedfef299fe5296e272b3968e60de`; the arc JSONL hash is `d6a91fb65a652913a9adf55b5054ab8fb9a87b1cd5565e16da501296a1b93f9f`. Arc records name source/target BFS index and successor index. The independent counts and SCCs agree with Task 1's operational claims.

### Necessity sketch, restricted to reachable states

Unsafety is monotone, so an unsafe state has no safe completion. A first BA refusal is immutable; application requires an allowed decision's opportunity, so such a refusal precludes BA consumption. If b is revoked and BA has neither applied nor retained a permit, no future transition can create a lawful BA opportunity: later issuance is rejected, admission requires current b, and revocation cannot be undone. Application cannot occur without an opportunity, so no authentic result can be consumed for that original. Completion therefore implies the predicate.

### Sufficiency sketch, restricted to reachable states

Choose a suffix with no further crashes or revocation and recover any down agent. Reachability supplies the custody and activation invariants checked on the graph: a permit belongs to a live server activation; an applied original retains its authentic receipt or the caller's durable record; a consumed original already discharged its obligation.

If BA has applied, advance its retained result through delivery, recording and consumption as needed. If revoked BA has an unapplied permit, A is live in that permit's activation and can apply immediately: A's servicing policy allows this while awaiting its own outgoing result. If b is current, no rejected BA issuance is reachable, because rejection requires irrevocable prior b revocation. Commit BA if necessary, retain its accepted grant before any outstanding first authorization, deliver and admit the request, and authorize/admit/apply it. An existing allowed decision already selected the same accepted grant. The predicate excludes an unrecoverable early refusal. At every such step, an occupied channel contains the same original or its authentic result and can be drained in the needed direction.

Consume BA to clear B's serial service guard. Then complete any outstanding AB work under c, which b revocation never invalidates. Its first refusal is unreachable because the exact c permission is always current. Retained authentic results can be consumed after recovery or receipt release using the inherited durable custody rules. This constructs a finite safe suffix. It is a source-level proof sketch under the reachable-state invariants, supported by complete finite enumeration, not a mechanized general theorem.

## Shortest adverse and successful witnesses

| Cut | Prefix steps | Safe suffix steps | Interpretation |
| --- | ---: | ---: | --- |
| Immediate b revocation | 1 | none | Shortest failure of universal completion |
| Early BA refusal | 4 | none | Grant not yet retained |
| Accepted grant after BA refusal | 5 | none | First decision stays refused |
| Revoke after BA allow, before permit | 6 | none | Historical allow is insufficient |
| Revoked BA permit survives | 7 | 15 | Old opportunity supports reciprocal consumption |
| A crashes holding revoked BA permit | 8 | none | Crash removes the sole remaining opportunity |
| BA applies with revoked-path permit | 8 | 14 | Authentic result and c-root suffix complete |

The first temporal-cut implementation searched only each target's shortest BFS prefix and missed the revoked-permit crash/application edges, because those targets merged earlier prefixes with different orderings. That failure is retained in `progress-cuts-red.log.gz`. The correction scans actual arcs, retaining a shortest source prefix followed by the qualifying edge, then the independently computed target suffix. This corrected witness selection did not change the graph, candidate predicate or transition semantics.

## Reproduction and compatibility

Use Go 1.27.3 with `GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b5/go-build STATICCHECK_CACHE=/tmp/alan-phase3b5/staticcheck`. Task evidence contains exact executable commands and direct exits in `checks.json`, `exports.json`, `mutants.json` and `old-progress.json`; standalone Python scripts reproduce the export, mutation and historical-progress comparisons using fresh disposable output directories.

New `authorityservice/cmd/revocation` and `delegationcorrespondence/cmd/progress -config S2/SeparateGrantRevocation` outputs matched byte-for-byte on two runs. All four historical exports were each run twice and matched the original archived bytes: twoagent, mutualobligations, authorityservice default export and authorityservice default audit. The default progress command also matched bytes generated from Task 2 base `e6e53deb`; it has no historical archived-export pointer, so that source-based comparison is reported separately.

Focused tests, full new all-edge/progress checks, affected race tests, command argument/output-error tests, and scoped Staticcheck passed. The initial missing opt-in acceptance/progress failures and the temporal-cut failure remain archived. The first focused run skipped the old environment-gated progress tests; the later run set `ALAN_TASK3_PROGRESS=1` and passed both historical progress queries. No new dependency, operational model change, State/Event shape change, manuscript edit or runtime evidence claim is part of Task 2. Concrete correspondence, whole-branch qualification, manuscript reconciliation and independent review remain separate tasks.
