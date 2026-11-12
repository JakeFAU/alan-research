# Concrete exact-root revocation correspondence

The opt-in relation maps ordered captures of N05, N06, N07 and P02-v2 from both sealed Phase 3B.4 producer modes into `S2/SeparateGrantRevocation`. The executable checks a finite forward trace at each capture and audits every projected event with the independent model ledger. The files retain the original producer commit `5b4662fd682a220ca3065b15e86672b7e56567e7` and tree `44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`. The analysis reads existing exports. It establishes no new runtime campaign.

## Relation

`CheckRevocationDirectory` retains the complete legacy `CaseReport`. The report separately exposes the validity of raw principal binding across captures, continuity, physical evidence, and support in the opt-in model. Reaching model support requires valid raw binding, verified declared bytes, an append-only capture chain, and both endpoint journals, their ready records and the matching complete endpoint snapshot. Endpoint zeroes are established by those artifacts together with absence of an allowed gate. The exported Store and endpoint records are unsigned. The corpus manifest anchors the input byte identities externally.

The mapper checks two generation-one principals and separate, singleton composition memberships. Both complete authority roots are held by B. Each is checked against the actual permission rows, the declared complete root, and the corresponding endpoint identity. It maps the delegated `ba-service` tuple at endpoint B to `B/g1:b` and the B-owned `ab-service` tuple at endpoint A to `B/g1:c`. The two tuples have different obligations and resources despite sharing a holder and composition membership. This finite mapping is injective. It supplies no general permission-algebra correspondence.

Original commitments, reply emissions and actual selecting executions supply custody transitions. Authenticated accepted message admissions permit a causal expansion at source authorization. This places send/admit or receive/record at an eligible fact and makes no measurement of wall-clock transport order. Accepted delegation and exact-root revocation map at their own lifecycle facts. The grant lifecycle fact locates the recipient's state, while the separate source authorization locates B's operation. Grant-ready transport and lifecycle-outcome consumption add no grant or revocation.

The first decision maps at its reservation prefix. The projector verifies the raw selector and witness, exact original, observation, selecting sequence, grant and service operation, and the serving composition root's immutable debit. An early exact-root refusal maps to model refusal with no covering path while preserving the concrete attempted selector. A fresh gate must use the same selected path, original attempt digest, authorization source, and ordinal. The mapped gate disposition must agree with the concrete gate. The report carries the runtime revision and epoch alongside the model event revision and activation. For the observed live-owner cuts, runtime epoch 1 maps to initial model activation 0. The model event revision is a revocation era. Runtime ownership revision remains 1.

A pre-existing allowed gate with a missing retained result projects a permit at an intermediate cut. Only a retained result fact with one matched allowed gate, contact, committed application and receipt permits atomic model Apply. Mutation and total values are abstracted to two opaque atoms while concrete operation and receipt identities stay in the sidecar. The endpoint wire lacks a gate ordinal. Attribution here relies on the unique allowed gate and exact original wire. It does not establish atomicity of runtime application and receipt recording. Only actual selected authenticated replies count as principal consumption.

## Observed verdicts

Each row applies to both sealed normal and race modes. Raw binding, physical evidence, and capture continuity pass for all five cases. The zero-contact rows include independently checked zero applications.

| Case | Retained ordering and boundary | Legacy support | Opt-in support |
| --- | --- | --- | --- |
| N05 | Grant 7, revoke 17, first refusal 21. No debit, gate, contact, or result. | No | Yes |
| N06 | Grant 8, first allow and debit 16, revoke 20, sole fresh refusal 22. Both principals remain unpaid. | No | Yes |
| N07 | Grant 8, allow 16, old opportunity 17, revoke 21, BA result 23 and consumption 28. AB uses c at 33/34/35 and is consumed at 40. Two contacts and two applications. | No | Yes |
| P02-v2 | No revocation. BA allows/gates/applies at 16/17/18 and AB at 28/29/30. Both principals consume authentic replies. | Yes | Yes |
| Applied unknown, revoked before a fresh gate | One contact and application, acknowledged whole-Store loss, then fresh refusal without a retained result. | No | No |

N06 cuts end at composition orders 16, 21, and 22. The last two carry the original allowance and debit. Only the final cut carries the fresh refusal. N07 cuts end at 17, 22, and 48. The two intermediate cuts carry an unresolved BA opportunity and no Apply. The final cut carries both applications and consumptions. The applied-unknown case remains outside the relation with explicit whole-Store loss and applied-without-retained-result reasons. It supplies no projection events.

Rehashed attacks reject altered capture ordering, aliased and swapped roots, missing physical artifacts, missing snapshot application coverage, moved grant/revoke facts, forged fresh allowance, a stale positive epoch, fabricated result after a refused gate, and forged request-selection joins. A coherently changed N07 intermediate capture remains raw-valid when it revokes the complete c root instead of b. The mapper rejects that capture at the exact lifecycle-root join. A compiled mutant that compares the shared composition root in place of the complete target incorrectly accepts it. The independent concrete gate mutant incorrectly accepts a fresh gate after the authentic N06 revocation.

The correspondence is scoped to these declared bounded histories. General rejected issuance, duplicate revocations, other root and principal maps, resource refusals, repeated physical exposures, post-loss recovery, and unknown final outcomes remain outside the claimed relation. No model crash correspondence, reverse simulation, arbitrary-state runtime completion, or distributed closure follows. The observed result controls establish only that the two principals consume authentic results in their retained traces.
