# Task 1 contract/construction review — revision 1

## Spec compliance

**REJECT — one Important contract-precision issue.** The baseline and public-construction requirements are supported by the inspected fixture and actual archived JSON. The exact version boundary is not ready to freeze: `docs/foundational/phase3b4/02-contract.md:27` conflates Store metadata semantics with supported artifact/reducer semantics. Correct this before using the document as implementation authority. This is a contract/construction gate only; no implementation approval is claimed.

Reviewed research base `9ddfe29a2a5fa43a06656514af236401cc07914b` → head `7a6d6dad839442c59fc95b26af3e07c4ab9e7d1c`, and runtime base `1fdd4c9e1720bf60b344914bcbc3d0938ba8d678` → head `3f8307664b6543331a41c943e91fb5bac8b376ef`, through the supplied diff packages. Research paths below are relative to `/Users/jacob/.codex/worktrees/phase3b4-research/alan-research`; runtime paths are relative to `/private/tmp/alan-phase3b3-runtime`.

## Strengths

- The representation is minimal and starts from the exact grant's `OriginalRequestID`, not a mutable register or conveniently selected message. Existing observation/consumption/history coordinates are sufficient for this candidate; rejected alternatives and their tradeoffs are explicit (`02-contract.md:5–11`, under `docs/foundational/phase3b4`).
- Strict selection-before-service commitment is distinguished from grant-before-decision, including delayed decisions, future selections, missing evidence, corruption and immutable denial witnesses (`02-contract.md:41–69`). The contract does not require grant acceptance before selection, and does not allow later selection to repair an earlier service.
- Current permission, revocation, owner gate, exact path, receipt consumption and generic unknown-effect policy remain separate from historical selection (`02-contract.md:73–79`). No production authorization code or principal oracle is changed in the runtime diff.
- The nonunique admission-observation index has an explicit query and validation shape; arbitrary transport IDs and absent ACKs are addressed. Inherited prefix scans are not misrepresented as constant-time work (`02-contract.md:81–92`).
- The construction fixture changes claims, dependent locators, branch targets and program hashes before public `Create`, then checks a real accepted grant, request selection, service identity, refusal layer and contacts (`internal/durabletest/phase3b4_construction_test.go:24–118,140–239`). Both branch outcomes distinguish executed reductions from source PCs. No retained-row repair or private authority writer occurs in this change.
- The matrix distinguishes same-peer identity substitution, register overwrite, false acceptance versus false refusal, source/gate revocation, future selection, sticky denial, wrong service locators and two actual owner renewals (`02-contract.md:127–182`). The final fixture/race logs contain six passing leaves; earlier compile and sandbox failures are explicitly retained as apparatus/development failures, not semantic evidence (`evidence/construction/construction-final.log:1–22`, `construction-race.log:1–22`, `01-charter.md:41`).

## Issues

### Critical

None found in this task-scoped review.

### Important — distinguish metadata semantics from supported program semantics before freezing

**`docs/foundational/phase3b4/02-contract.md:27`** says: “All four formats retain ... reducer `alan/1d/v1`.” That is accurate for the Store's fixed `metadata.semantics` value, but not an accurate exhaustive program/reducer compatibility contract. The unchanged source accepts `alan/1d/v1`, `alan/1e/v1`, and `alan/1g/v1` (`internal/executor/author_metadata.go:10–12,44–55`); artifact insertion and loading use that predicate, then check each artifact's own retained semantics (`internal/executor/store.go:519–544,561–580`). Store metadata is separately written and checked against `SemanticVersion` (`store.go:387,993`). Existing Phase 3B.3 adapter construction already supplies an authored artifact under the authority profile (`internal/durabletest/phase3b3_adapter_test.go:99–118`).

Because this document is the exact interface authority for the next implementation task, the present sentence can be implemented as an unintended `alan/1d/v1`-only restriction, contradicting preservation of formats 13/14 and surrounding adapter policy. This is not a demonstrated production regression: Task 1 has no such production edit. It is a blocking ambiguity in the proposed frozen compatibility contract.

**Required correction:** name the fixed Store metadata value explicitly; separately preserve the existing per-artifact semantics and validation rules, including supported authored/collection versions, without implying that every feature combination is newly enabled. State that the new public example uses `alan/1d/v1` if that is its intended narrower fixture domain. Carry an existing authored-adapter old-profile acceptance regression into later compatibility qualification, and specify its v2 disposition. A wording correction and corresponding exact compatibility requirement suffice for this Task 1 gate; do not add production code here.

### Minor

None requiring change. The transport/index domain is specified adequately; later tests should exercise arbitrary transport IDs and repeated accepted transports for one logical observation as that specification requires.

## Named source risks checked outside the diff

1. **Selection-row revision equality versus the existing recipient-revision policy.** Checked `evaluateEffectAuthority` (`internal/executor/request_authority.go:435–461`) and issuance (`:349–367`): v1 explicitly matches grant recipient revision to service source revision, and does not separately compare the selected history row's revision. Checked production subject mutation sites: current lifecycle transitions update generation or retirement (`internal/executor/lifecycle.go:339–344`); no production writer changing subject revision was found. The candidate's positive domain preserves logical identity/revision (`02-contract.md:11,44,79,150`). Within this declared static fixed-revision domain, the new comparison is redundant with authentic validated reduction provenance, not a demonstrated false refusal or changed lifecycle policy. Do not extrapolate this conclusion to a future revision-changing API; that would require a new policy decision and test.
2. **Reducer compatibility.** Checked the exact `supportedSemantics`, metadata, artifact creation/load, and existing authored-adapter call sites listed in the Important finding. No test rerun was needed to establish the distinction.

## Actual causal evidence inspected

Read `final/inspect.json` raw table exports directly from the committed baseline/construction tar archives, plus the corresponding manifests and endpoint JSON. No SQLite archive was opened or extracted. The following are direct joins from those raw records, not copied conclusions from `causal-assertions.json` or the self-report:

| Archived case directory | q selection / e | composition orders selection, grant, e, first decision | observed result |
|---|---|---|---|
| `baseline-allowed/ordinary-s2-handshake-975364978` | A3 / A4 | 13, 7, 14, 16 | authorized; each endpoint 1 contact; both paid |
| `baseline-allowed/control-wrong-prior-selection-1207259051` | A2 / A4 | 5, 8, 14, 16 | `authority_request_not_admitted`; 0 contacts/charges/gates; both unpaid |
| `construction-final/construction-adjacent-715589528` | A3 / A4 | 13, 7, 14, 16 | authorized; each endpoint 1 contact; both paid |
| `construction-final/construction-request-ready-0-1949628549` | A2 / A4 | 5, 8, 14, 16 | `authority_request_not_admitted`; unpaid, no exposure |
| `construction-final/construction-request-ready-1-3116398634` | A2 / A5 | 5, 8, 15, 17 | `authority_request_not_selected`; unpaid, no exposure |
| `construction-final/construction-request-ready-3-3886977122` | A2 / A7 | 5, 8, 17, 19 | `authority_request_not_selected`; unpaid, no exposure |
| `construction-final/construction-forward-true-477812333` | A2 / A7 | 5, 8, 17, 19 | `authority_request_not_selected`; scratch=1; unpaid, no exposure |
| `construction-final/construction-forward-false-4058696548` | A2 / A7 | 5, 8, 17, 19 | `authority_request_not_selected`; scratch=2; unpaid, no exposure |

In all eight snapshots q is `["b",1,1,0]`, g is `["b",1,2,0]`, q has an allowed source authorization and accepted matching admission, its observation's consumed sequence equals the actual selecting history row, and the accepted grant names the actual service. Embedded claim copies equal the manifest claim; retained program bytes hash to their IDs; the actual service occurs in the matching committed history intents. The true branch's A4/A5/A6 successor PCs are 4/5/8; the false branch's are 6/7/8; both execute service at A7 with successor PC9. Thus code position was not substituted for execution sequence. Positive snapshots contain two `authority_admission` facts at orders 17 and 29; refused snapshots contain none.

Raw `inspect.json` SHA256 anchors: baseline supported `6fd7624dc52d36188d8ffae85d06e6833bf631e8c16ca57b1339f04d4e64ec18`; baseline refused `bd200682b56cf423e7e188bd547132b6b39e408cb16658e7c164b799e219a730`; branch true `384c2e9ce41d90437c712fa89ac0e29bf702c698bb86ce1624d9d2d1f782d2d5`; branch false `356903ddfaeaaaa2b75273982260a4af0d9332fdfbaf7d15d6534067b06ac525`.

## Verification and unresolved items

- Performed read-only diff/source inspection and focused Python JSON joins/assertions over eight exported snapshots; these passed. Parsed all 18 added JSON diff hunks, including both manifests (68 and 543 members) and the preserved historical production delta. This is focused evidence inspection, not a new independent policy oracle or a rerun of reported suites.
- The preliminary reviewer inspection script queried a nonexistent `authorizations` table name and another split a quoted embedded diff marker incorrectly; these reviewer apparatus errors were corrected to `attempts` and anchored diff-header parsing. They are not runtime or task failures. Likewise, gate counting was corrected to the actual `authority_admission` spelling before drawing gate conclusions.
- No runtime tests, race suite, model execution, SQLite operations, repository mutation, Git commands or subagents were performed. The sole output file is this report.
- **Gate blocker:** correct the Important version-boundary wording and specify preservation of accepted artifact semantics. Re-review that narrow contract change before production or principal oracle work.
- **Later-task obligations, not missing Task 1 execution:** v2 implementation/RED–GREEN behavior, old-reader rejection of 15/16, new-reader preservation of 13/14 including authored artifacts, malformed/missing index behavior, arbitrary/multiple-transport joins, future/sticky-prefix tests, process cuts, independent oracle and mutant kills, cost guard, and manuscript reconciliation remain unexecuted by this review. Their future completion must not be inferred from this construction evidence.
- Provenance and old-reader build/hash records were inspected as retained records; the binary was not executed and remote PR state was not re-fetched. This gate does not certify whole-repository behavior or complete source-pair qualification.

## Assessment

**Task quality: Needs fixes.** The construction evidence is carefully bounded and supports the finite public-authoring claim. One precise compatibility clarification is needed before the contract can safely govern implementation; no broader redesign or additional construction run is requested.
