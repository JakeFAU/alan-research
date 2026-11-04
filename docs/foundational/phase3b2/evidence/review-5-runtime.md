# Scope 5 independent review: runtime correspondence and Task 3

Disposition: **APPROVE** within scope 5 for Task 3's runtime correspondence, specification and code quality, bounded to the inspected development source and reported development execution below. This does not approve the whole task: the separate scope 4 reviewer has a generation/history missing-edge REJECT under correction and must independently re-review it. This is not approval of final immutable-source qualification or delivery. No production source or positive database was changed by this review. No new runtime experiment was necessary, and no test execution is claimed as reviewer-run.

## Reviewed question and source

Does the ordinary same-Store fixture preserve principal-local static permissions and original-result custody across real process interruption, while exposing the exact missing S2 representation and the runtime's conservative unknown-outcome boundary?

I read the Task 3 brief, report and review patch; all new `examples/communication/authority/` files; `internal/durabletest/phase3b2_runtime_test.go`; the frozen `01-contract.md`; draft `03-results.md` and `04-runbook.md`; and the owning source in `internal/executor/{messages,adapters,dispatch,runner,store,reconciliation_profile}.go`, `internal/boundary/boundary.go`, `examples/effectboundary/adapter.go`, and existing durable process/backup/replay helpers.

The current four example files match the supplied review patch. The runtime test adds, beyond that patch, incoming ordinary-message admission-chain verification, permission revision 1, and exact source authorization ID/reason/revision checks. Those additions were inspected directly. Reviewed SHA-256 values:

| File | SHA-256 |
| --- | --- |
| `examples/communication/authority/adapter.go` | `a32c2b955edc94e18d655aae3412526eae6adbb86aae20ced6bd7f1f1b0ea27b` |
| `examples/communication/authority/deployment.go` | `1018b6d798f028a62e2b5b6cdb3c4601ce51fe051ed0153f69f68fa41d4c2257` |
| `examples/communication/authority/deployment_test.go` | `7988eb5d35804e9f3ed81ddf770e9dd9f1aa8183e262bb3199691752da8eab42` |
| `examples/communication/authority/README.md` | `0b78370480f2d6ce34eebf7974faa61f5c736abc5f7767822b0926dddec711be` |
| `internal/durabletest/phase3b2_runtime_test.go` | `2d1ea47e0a340bd44901c2fce6206349c23bc5195a5b99f9a500065b54a3ed76` |

## Specification and semantic findings

- Construction uses public `author.Compile`, `Application.Provision` and runtime Create/Commit/Deliver/Acknowledge/Open/Run. Exactly two roots exist; A's unwanted automatically generated permission is removed before Create, and A's own permission remains. The only positive SQL-like edits are absent: actual DELETE/UPDATE statements are confined to disposable backups for rejection controls.
- The trusted registry maps each whole capability to a fixed endpoint and `Prepared.Model` resource/scope. Payload claims and authentic consumed result values remain data. The runtime checks static kind/source/target/operation/revision/obligation permission tuples, not an independently interpreted resource/scope algebra. This is accurately bounded in the fixture README and contract.
- Owning source confirms the cut order: immutable source disposition in `AuthorizeIntent`; exact static `adapterGrant` plus retained `adapter_attempts`; `admitDispatch`'s `dispatch.consider/pre-admission`; volatile registration; then `adapter.authorized/pre-transport` inside transport. The historically named adapter boundary occurs after physical admission. Receipt retention, settlement and authored consumption are separate.
- Nine table entries really invoke a child, wait for an identity-checked acknowledged barrier, call SIGKILL and check the actual signal disposition, and exec a distinct fresh owner. The before-authorization, after-refusal and mutual-wait manual boundaries only observe public-runtime state. Both retained run epochs must increment by exactly one. Earlier programs, subjects and permissions must remain identical, and existing source decisions/adapter attempts/results/settlements must survive as exact rows.
- For the two allowed/no-result cuts, the endpoint has zero contacts, but this external fact is never installed as a runtime receipt. Fresh Run must report exactly the interrupted original in `UnknownAdapters`, keep results/settlements absent, preserve both unpaid originals and leave both worlds unchanged. The source's generic existing-allowed-attempt branch returns without transport. Known receipts can settle and consume without reissuing the completed operation.
- The paid-result oracle checks actual history selection, not a paid flag alone. It follows the original incoming ordinary-message intent/admission and consumption before the service operation; exact source authorization and principal permission; allowed prepared attempt; successful candidate/result and digest; exact settlement and consumed adapter observation before reply; reply admission; and actual original-result selection. The separate endpoint oracle matches the exact operation key, mutation, total and delta against one application in the fixed configured world's ledger. Mutations preserve superficial paid observations/result counts while breaking a decisive chain edge.
- `Permission` has no generic grant identity, issuer/grantee ancestry, attenuation set, original-service binding or retained revocation chain. `ReceiptGrant` is specifically a resolver/original/endpoint receipt-query/adoption construct. Therefore the missing S2 primitive is an identified representation boundary; no broad delegation implementation or substitute static grant is smuggled into the fixture. Delegation-retention cuts are correctly inapplicable.
- The model's atomic Apply-plus-receipt assumption cannot establish generic runtime completion: the runtime can safely remain unpaid after a pre-call reservation loss. Draft results explicitly retain this disagreement and disclaim full model/runtime refinement, distributed ownership, arbitrary delegation and shared-resource isolation.

## Preserved finding R5-1: checked kill evidence export

The legacy kill helper's `emit` does **not** consult `ALAN_EVIDENCE`. It writes `evidence-*.json` under the `binaries(t)` directory created as `/tmp/alan-task4-evidence-*`, and emits the complete record through `t.Logf("DURABLE %s", raw)`. In contrast, `phaseEmit` honors `ALAN_PHASE2_EVIDENCE`. Merely setting both environment variables does not copy the checked kill record into the phase export directory.

The runbook briefly asserted otherwise during review; this was a documentation/evidence-packaging finding, not a runtime transition defect. **Resolved on source reinspection:** `04-runbook.md` now accurately requires full verbose output, explains the legacy temporary files, and states that final archival extracts the complete `DURABLE` records. Final qualification must preserve the complete verbose command output containing each checked `DURABLE` SIGKILL record, or explicitly archive those legacy JSON files, and describe that route accurately. A crash-case JSON string saying `signal: SIGKILL` alone is not the checked process disposition. The parent was notified before final source freeze. This review does not approve a claim that final evidence already exists.

## Execution evidence and remaining obligations

I inspected the actual development logs at `/tmp/alan-phase3b2-runtime-final-{unit,integration,race}.log`: focused unit passes, integration pass in 13.317s, and race/integration pass in 45.481s. These are prior implementer executions, not reviewer executions, and the earlier race run does not by itself bind the subsequent deterministic oracle additions. The retained RED log fails because the new package has no implementation yet; it establishes a missing-construction baseline, not a runtime semantic counterexample.

The inspected integration Staticcheck log reports SA4006 in unchanged `phase2_crash_test.go:262` and ST1013 in unchanged `phase2d_adversarial_world_test.go:95,100,106,112,123`. Its gate is explicitly non-green; no test or diagnostic was weakened. The report separately records the normal Staticcheck pass.

Final source freeze, immutable-source normal/race execution, complete command logs, raw coherent archives and source/evidence manifest validation remain parent-owned work. No additional blocking finding remains within scope 5. Overall Task 3 remains subject to the scope 4 REJECT and independent re-review, plus final delivery obligations. R5-1's documentation is corrected; its final archive condition remains to be verified.
