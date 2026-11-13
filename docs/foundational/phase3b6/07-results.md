# Phase 3B.6 results

An accepted grant and an allowed gate establish historical admission. Under the one-service model, the same admission, exact-root revocation and fresh refusal can accompany zero or one application. Complete, identity-bound endpoint evidence can resolve that ambiguity at a declared cut. Continued local execution using the original result requires authentic local custody.

## Conditional observation argument

One selected BA original starts after admission through exact root $b$. Both principals share one Store owner. The endpoint and durable Store survive one owner death. The contract permits at most one launch, ingress, application, retention and consumption, and an optional fresh refusal after reopen. The distinct $c$ root remains valid. Loss clears the original incarnation's volatile permit and unretained response. It preserves packet transfer and world application. No transition imports endpoint evidence or restores the original opportunity after loss. [Model and assumptions](02-model.md).

Let $O_{\mathrm{store}}$ include the binding, admission history, owner stage and activation, current authority, fresh-refusal history, retained value and consumption. The experiment's death/Open protocol supplies the owner stage. Whenever two histories end at the same $O_{\mathrm{store}}$, a decision based only on this observation must give the same answer for both. The finite witnesses below pair identical observations with different application facts and with different authentic world values. A fresh refusal therefore cannot resolve either difference for those histories.

$O_{\mathrm{dossier}}$ adds complete endpoint ingress and application evidence at the cut. Missing or partial endpoint evidence leaves the application unknown or the cut unsupported. A certified zero describes the declared cut. An already launched packet can still ingress and apply after that cut, including while the owner is stopped and $b$ revoked. Positive endpoint evidence identifies a world fact. Recording the value requires a successful authentic response in the original owner, and no forensic import transition exists. A value already durably retained before loss can be consumed after reopen.

## Finite evidence

The independent event ledger audited all 63 reachable states, one shortest prefix per state and all 104 arcs. It checked 33 merged targets and their 42 additional incoming arcs and derived 31 Store observation classes and 52 dossier observation classes. Witness lengths count events from the post-admission initial cut. [Independent audit and shortestness criteria](03-observations.md).

| Witness | Lengths | Result |
| --- | --- | --- |
| Equal local observation with zero or one application | 4 / 7 | Both cuts are reopened, revoked and fresh-refused, with no local result. |
| Equal local observation with distinct authentic world values | 7 / 7 | The two applied values remain unretained. |
| Negative cut with a pending packet | 3, then 2 more | A stopped and revoked cut with zero ingress and application admits subsequent ingress and application. |
| Consumption of a retained result after reopen | 9 | An authentic result survives loss and is consumed despite exact-root revocation. |

The audit exercises nine lawful and eleven hostile controls. Four compiled ledger mutants fail targeted assertions for application as retention, revoking the distinct control root, admitting an old response after reopen and deleting a pending packet on loss. The counts concern the full reachable graph of this finite contract. The observation argument specializes Window $W_3$ and Phase 2B/2C lost-result and provisional-absence precedents to exact-root revocation and fresh refusal. The Phase 3B.5 completion predicate still applies within its atomic application/receipt domain. An applied-but-unretained cut falsifies an extension of that predicate to the split model.

## Relation to historical runtime evidence

Four selected dossiers in two sealed archives retain all 16 declared captures. The separate source-conditioned relation supports 12 cuts and explicitly reports four pre-admission prefixes outside the model. Each supported cut has concrete joins, a causal certificate and an independently audited existential witness. Successive witnesses extend earlier supported prefixes. Runtime producer commit `5b4662fd682a220ca3065b15e86672b7e56567e7`, tree `44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`, and the owning normal/race archives are historical inputs. The relation assumes honest pinned source and sealed unsigned archives. It supplies evidence for the selected service in these dossiers, without recovering a complete physical event order. [Source, causal joins and witnesses](04-correspondence.md).

| Dossier profile, in each producer mode | Captures and witness lengths | Supported conclusion |
| --- | --- | --- |
| Applied-unknown-revoked-before-fresh-gate | Five, with two pre-admission prefixes and three witnesses of lengths 4 / 5 / 7 | The matched endpoint post-commit and runtime in-flight acknowledgments both precede acknowledged owner death. The endpoint cannot emit a successful receipt until the post-death release. Cancellation or error may end the wait earlier. Reopen renews both principals once, followed by exact $b$ revocation and an epoch-current fresh refusal. The original application has no retained BA result or consumption. One post-reopen scheduler dispatch is compatible with zero additional physical ingress. |
| N07 | Three, with witness lengths 0 / 1 / 7 | Synchronous callback ordering certifies zero ingress and application at both intermediate cuts. Final evidence joins one BA application to an authentic retained result and actual principal consumption. The separate AB completion under $c$ is a concrete control. |

Raw validity, continuity, physical coverage, protocol attribution and model support are separate verdicts. Runtime Open supplies no token or Store field. The runtime acknowledgment uses the sealed-dossier process anchor for those identities. The endpoint and runtime acknowledgment order is unobserved. Successful response withholding establishes neither continuous barrier occupancy nor transport liveness. Witness order beyond the certified constraints is a model choice. [Protocol-source index](evidence/task3/protocol-source/index.json).

The old atomic forward relation continues to exclude the genuine applied-unknown case. The new witness relation covers its selected BA service under the stated source and archive assumptions. No new Phase 3B.6 runtime campaign ran. The historical v1 and Phase 2B/2C comparisons and the N07 zero observations do not supply a matched v2 loss pair. General simulation, pair recovery, distributed revocation, arbitrary-agent composition, retry/import authority and global progress remain outside the result.

## Next bounded experiment

Can a matched v2 loss after admission and before the original call followed by the same exact-root revocation and fresh refusal schedule yield a zero-application cut with $O_{\mathrm{store}}$ equal to the applied-world cut? The experiment must certify the selected original and pre-call boundary, obtain complete endpoint evidence at the declared cut, and compare the full declared local projection. An application at the proposed zero cut or a difference in that projection would falsify the proposed match. [Frozen question and qualification domain](01-contract.md).

[Qualification and reproducibility](06-qualification.md) · [Independent reviews](08-reviews.md) · [Runbook](05-runbook.md)
