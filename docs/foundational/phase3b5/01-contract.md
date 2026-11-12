# Phase 3B.5 — Separate grants and exact-root revocation

## Question and candidate result

Can a separately accepted request-bound grant and a later exact-root revocation coexist in the bounded two-agent model, preserving the lawful first authorization while refusing fresh physical admission? Can the sealed N06/N07 runtime histories be mapped into that model without inventing effects or broadening revocation to unrelated authority?

This packet continues the merged Phase 3B.4 question. Its initial status is a hypothesis and experiment contract, not a claim that the checks below already passed. The deliverable is a bounded result, retained counterexample, or precisely recorded obstruction, with an independently reviewed draft research PR. No runtime change is initially justified.

## Owning evidence and scope

Research starts at `4fea6c31eb8f438c0dea414aa5f8d7d0ecb5e8a4` (PR #45 merged). Runtime main is `c72c396df09c55b0be1a91b2557dcec2ef6457a8` (PR #13 merged), full tree `52f4d7c9a2f573e78112be0c1f07b26834918603`. That tree matches the previous final prospective delivery tree. Phase 3B.4's sealed prospective-main raw archives were produced by historical synthetic commit `5b4662fd682a220ca3065b15e86672b7e56567e7`, tree `44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`. Keep producer provenance distinct from current source. A new independent analysis of those archives is not a new runtime execution.

Runtime revocation targets `AuthorityRoot{Run, RootRun, Permission{Kind, Source, Target, Operation, Revision, Obligation}}`. It is not grant-ID revocation and not closure of a composition budget root. N07's two B-held service permissions have different complete tuples. Revoking the delegated BA root leaves the distinct AB root current.

## One opt-in finite model

Add exactly `S2/SeparateGrantRevocation`, obtained by `authorityservice.GrantRevocationConfig()`. Preserve the old sixteen `Configurations()` values/order and all four historical deterministic exports byte-for-byte. Preserve existing State/Event JSON shapes, enum values, old bindings, and default command scope. Append the new tuple `ABOwnPermission = (increment,counter-AB,blue)` after existing tuple values; its root is `B/g1:c`, used only by the new configuration. BA remains bound to `B/g1:b` and the single `grant:BA:b`.

The original two-agent reciprocal custody transitions and servicing/serial guard remain unchanged. Generations stay at one; at most one recoverable crash per agent, one retained original each way, two opaque results, one-slot channels, retry cycles, one no-redelegation B→A grant issuance attempt, and one exact b-root revocation. Application plus receipt remains atomic in the model. No regrant, renewal, multihop, extra owners, extra originals, grant-ID revocation, or non-atomic external effect window is introduced.

1. BA commitment carries no inline grant. Separate issuance is possible after commitment while B is live. Issuance before revocation retains accepted status; issuance after revocation retains rejected status. Request commitment is not suppressed by revocation.
2. Revocation by live B changes only current validity of b and its delegated descendant. It never rewrites retained issuance, first decision/path, expenditure, permit, result, or consumption.
3. First authorization is immutable. An early refusal cannot be repaired by a later grant. Fresh admission rechecks the same selected path and current authority, and yields one activation-bound opportunity.
4. A pre-revocation opportunity can apply after revocation. Server crash clears it; recovery cannot resurrect it. A new opportunity after revocation is refused. Already retained authentic results remain consumable.
5. AB's distinct c root remains usable, subject to the unchanged B serial service guard. Payloads and returned results never confer authority.

## Attacks and completion question

Retain shortest counterexamples to universal completion, including early refusal, revocation before a required opportunity, and crash of a previously admitted server after revocation. Preserve retry SCCs; existential completion is not inevitable progress.

Candidate exact predicate on reachable states is `Unsafe==0 && BA.Decision!=Refused && (!Revoked || BA.Applied || BA.Permit)`. AB refusal should be unreachable under its live c root. Enumerate the complete graph without cutoff, compare the predicate to independently computed reverse reachability, and give a source-level necessity/sufficiency argument with each assumption visible. If the predicate fails, retain the witness and correct the claim before presenting a result.

The independent ledger must not call operational successor, binding, authority-selection or current-validity predicates. Audit one shortest prefix per finite state, every outgoing arc, and ledger equality at merged targets. This does not enumerate every cyclic history. Retain hostile traces and compiled checker mutants killed by a meaningful assertion, not a compilation error or panic. Discriminate historical decision validation from present admission validation, b/c scope, grant ancestry/binding, revocation actor/order, rejected issuance, and activation-bound permits.

## Scoped concrete forward relation

Use opt-in `CheckRevocationDirectory(dir string) (RevocationCaseReport, error)` and a separate command. Existing raw verdicts and legacy projection behavior remain unchanged. The new report separates raw validity, physical evidence, continuity, and model support, and retains the legacy report for comparison. Project authentic ordered facts through uniquely enabled model edges, then independently audit the resulting events.

The fixed corpus is N05, N06, N07 and P02-v2 from each sealed normal/race archive. N06 is central: grant order8 → first allow/debit16 → exact-root revoke20 → sole fresh refusal22, zero contacts/applications and no result. N07 discriminates surviving permission: grant8 → allow16 → opportunity17 → revoke21 → BA result23, followed by distinct AB own authorization/gate/result33/34/35 and both authentic consumptions. N05 covers grant7 → revoke17 → first refusal21. P02-v2 preserves the nonrevoking control. These expected coordinates are recon findings, to be checked anew by the executable relation.

Require declared ordered captures, original principal IDs, successful independent raw binding checks, hashed member identity, append-only continuity, both endpoint journals and endpoint snapshot. An absent receipt alone cannot prove zero contact. N06's zero-contact claim requires the external artifacts and no allowed gate. N07's intermediate allowed-gate/no-result cuts may project an unresolved opportunity; neither the opportunity nor the final dossier invents an early application. A final atomic Apply abstraction requires the declared single exposure, actual retained result and matching endpoint application. Exact identities, complete roots/tuples, q/o/c/g/e binding, operation/fact IDs/orders, runtime revision/epoch, debit and gate ordinal stay in a sidecar. Model revocation era is not runtime ownership revision.

A fixed injective mapping sends the exact BA root to b and the distinct B-owned AB root to c. Verify the actual accepted revocation targets the former. No permission-algebra isomorphism is claimed. Source authorization locates a causally admissible send/admit expansion justified by actual accepted transport, not a measured wall-clock delivery time. First decision, path/root/grant and fresh gate disposition must agree, not merely boolean allow.

Whole-Store loss, unknown or applied-unrecorded effects, repeated physical exposures, resource refusal, unsupported roots/duplicate revocations and general rejected runtime grants remain outside this narrow projection. Retain the genuine `applied-unknown-revoked-before-fresh-gate` case as an explicit exclusion. No model crash correspondence, reverse simulation, arbitrary-state runtime completion, general distributed theorem or N-agent closure follows.

## Verification and delivery

No new dependencies, bibliography entries without use, direct SQL evidence repairs, runtime production edits, unrelated branch cleanup or archived SQLite opens. Use Go 1.27.3 and Staticcheck, exact commands and source manifests. Run targeted RED/GREEN tests, complete new graph/audit/progress, old export byte comparisons, concrete and compiled mutant attacks, fresh full ordinary research tests, affected race tests, build/module/format checks and scoped Staticcheck. Disclose optional skips and inherited global lint findings separately. Final qualification pins executable inputs and records any later documentation-only changes.

Reconcile Markdown/LaTeX Abstract, body and Conclusion; rebuild and visually inspect the PDF. Preserve historical Phase 3B.4 restrictions as historical findings. No reference change is expected. Independent task reviews and one whole-branch review precede draft PR delivery. No merge or release.
