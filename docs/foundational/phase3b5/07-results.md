# Phase 3B.5 results

Can a separately accepted request-bound grant preserve a first allow when exact-root revocation bars fresh admission, and can the sealed single-host observations be mapped without broadening that revocation or inventing an effect? The finite model and the scoped concrete analysis answer both questions under the conditions below. [Experiment contract](01-contract.md).

## Finite graph and source-level argument

`S2/SeparateGrantRevocation` has two generation-one agents, one original in each direction, one-slot channels, two opaque results per original, at most one recoverable crash per agent, one grant attempt and one b-root revocation. A serves while waiting. B waits to apply AB until it consumes any committed BA result. BA uses B's complete root b and the non-redelegable request-bound descendant. AB uses B's distinct complete root c. Issuance status and the first decision are retained. Fresh admission checks current authority on the selected path. Revocation preserves an old opportunity, while a server crash clears it. A model application consumes its opportunity and retains the effect result atomically. The model has no separate debit ledger. [Operational model](02-model.md).

| Complete finite graph check | Observed |
| --- | ---: |
| Reachable states and arcs | 229,834 and 1,807,246 |
| Unsafe prefixes or arcs in the independent audit | 0 |
| Complete states | 32,768 |
| States with a finite safe completion suffix | 203,833 |
| Safe states without such a suffix | 26,001 |
| AB first refusals and safe-suffix predicate mismatches | 0 and 0 |
| States in cyclic SCCs and cyclic SCCs | 188,614 and 75,070 |

The independent ledger interprets event descriptors without operational authority helpers. It checks one shortest prefix per state, every outgoing arc and ledger equality at merged targets. The enumerators use operational successors to supply the graph. The reverse search and SCC calculation are independent of the operational explorer's completion and SCC calculations. The checks quantify over the reachable finite quotient. [Independent audit and progress](03-authority-progress.md).

On that set, the following predicate is equivalent to existence of a finite safe suffix consuming both authentic originals.

```
Unsafe == 0 && BA.Decision != Refused && (!Revoked || BA.Applied || BA.Permit)
```

The source-level necessity sketch uses monotone unsafety, immutable refusal, and the lack of any transition creating an opportunity on an already revoked path. For sufficiency, choose a suffix without further crashes or revocation, recover down agents, and advance BA through its retained result, surviving opportunity or still-current path. Reachable custody lets delivery proceed. BA consumption clears B's serial wait, and c remains available for AB. This sketch uses the stated reachable-state and service invariants. It is supported by the complete enumeration.

Immediate b revocation is the one-step shortest universal-completion counterexample. Further shortest losing prefixes isolate early BA refusal at four steps, revocation after allow and before admission at six, and A's crash with a revoked BA opportunity at eight. The crash witness and the eight-step application witness use actual qualifying arcs from shortest source prefixes. The latter admits a 14-step safe completion suffix. Retry cycles also allow an incomplete execution in states with safe suffixes. [Witnesses and graph hashes](03-authority-progress.md).

## Conditional concrete forward relation

The opt-in checker evaluates raw binding, physical evidence, append-only continuity and model support separately. It uses verified, declared ordered captures and both endpoints' ready records, journals and matching complete snapshot. It expands facts through uniquely enabled model events and checks those events with the independent ledger. Mapping the two complete roots to b and c is injective on this corpus. The concrete sidecar retains exact identities, orders, ownership revision, epoch, debit and gate ordinal. The model event era and runtime ownership revision are distinct. Placement of transport steps is causal, using actual accepted admissions. The unsigned exports are anchored to their declared bytes by the archive manifests. [Relation and evidence assumptions](04-correspondence.md).

The observations come from historical sealed normal and race exports. Eight final cases and twenty captures are supported across both modes. Both modes give the same case dispositions.

| Case | Observations in each sealed mode | Opt-in projection |
| --- | --- | --- |
| N05 | Grant 7, revoke 17, first refusal 21. Zero contacts and applications. | Supported |
| N06 | Grant 8, first allow and debit 16, revoke 20, sole fresh refusal 22. Zero contacts and applications, both originals unpaid. | Supported |
| N07 | BA opportunity 17, revoke 21, BA result 23 and consumption 28. Fresh AB authorization, gate and result under c at 33/34/35, with consumption at 40. Two contacts and applications. | Supported |
| P02-v2 | Nonrevoking control. Two contacts and applications. Both original results are consumed. | Supported |
| Applied unknown before a refused fresh gate | One contact and application, acknowledged whole-Store loss and no retained result. | Excluded, zero projection events |

N06's endpoint evidence and absence of an allowed gate establish zero contact in the declared observation interval. A missing receipt alone would leave contact open. The two intermediate N07 cuts each retain one unresolved BA opportunity. The final cut maps each model Apply only where a retained result matches a single allowed gate and actual endpoint application. The wire omits the gate ordinal, so this attribution relies on unique allowance and the exact original wire. The mapping abstracts application plus receipt as atomic inside the model.

N07 separates the two complete roots even though both are B-held and share composition membership. Revoking both through a single collapsed root would reject the observed fresh AB admission. The independent ledger mutant that revoked c along with b failed its assertion by rejecting the lawful c-backed suffix. A separate raw-valid wrong-target attack and compiled mapper mutant discriminate the complete lifecycle-root check from a check of shared composition membership. This finite separation concerns the retained admission in N07. [Concrete verdicts](evidence/task3/verdicts.json) and [checker attacks](evidence/task3/mutants.json).

The scoped relation leaves repeated exposures, resource refusals, general rejected runtime grants, unknown results and whole-Store process loss outside its domain. Runtime application/receipt atomicity, reverse simulation, arbitrary-state runtime completion, general permission-algebra equivalence, distributed recovery and composition beyond two agents remain open. The legacy projection still uses its historical domain. Ten before/after legacy reports were byte-identical. [Qualification](06-qualification.md).

The next bounded question is whether retained evidence across one acknowledged process loss can distinguish an admitted opportunity from an applied mutation whose result was never recorded, under a fixed endpoint observation contract.
