# Phase 3B.4 charter and construction gate

Status: Task 1 construction demonstrated; candidate contract **review-pending**. No candidate production predicate, principal oracle or model extension has been implemented. This packet is not an approval or a claim that the new semantics works.

Question: Can a service retain an exact binding to the authentic request it selected before intervening computation and same-Store recovery, while every authority decision and fresh physical admission still checks the original complete permission tuple, exact service, selected path and current authority?

The smallest boundary is the delegated request-selection predicate and its retained semantic version. Two initial generation-one principals, A servicing/B serial, one hop, one trusted Store owner, existing lifecycle issuance, source authorization, budgets, dispatch lifetimes and endpoint protocol remain the experiment. No new dependency, authority registry, scheduler, owner topology or generic dataflow claim is justified.

## Inspected bases and provenance

| Repository | Selected base | Tree |
|---|---|---|
| Research | `9ddfe29a2a5fa43a06656514af236401cc07914b` | `a566a5eb31335ae68ea9c9073a1f56b5102a4a11` |
| Runtime | `1fdd4c9e1720bf60b344914bcbc3d0938ba8d678` | `24b1cbaf259fe7e211c5464242daf65ad5dcbb0f` |

The controller's retained GitHub metadata reports research PR44 merged at 2026-11-04T21:26:36Z, head `d365c7e8fbc2985b998936cc162f322f6a8f2367`; runtime PR10 merged at 21:26:49Z, actual final head `9bc02caa114ab1341524849d580dc7ed728e44f0`. Each actual head tree equals its selected merge tree. The head and historical qualified commits are not ancestors of the selected bases; do not infer ancestry from tree equality. Local `git diff`, tree IDs and `merge-base --is-ancestor` outcomes are retained in [baseline provenance](evidence/baseline/provenance/).

Runtime's historical reported delivery `91963d5a81e4ce16f424aae36cc4fb382a56e2cd` differs from the selected base in 41 files, chiefly previously qualified series-scaling evidence and tests. The only production delta is the immutable `invocationProtocol` boolean in `compositionProfile`, initialized at profile retention and used by `invocationEnabled`; the request-authority predicate is unchanged. This audit inspected that full production diff. Research historical delivery has the same tree as its selected base; historical executable qualification `7894bfae268e6dc0f5418bda3c6022f3a8a55540` differs in later documentation/manuscript/evidence, as retained in the diff inventory. These are baseline differences, not Phase 3B.4 changes.

The immutable baseline old-reader binary was built before fixture edits at `/tmp/alan-phase3b4-task1/compatibility/old-reader.test`. [Build provenance and SHA256](evidence/baseline/baseline-command.json) record Go 1.27.3, selected commit/tree, the clean pre-fixture status and exact build command. It supports `TestPhase3B3Owner` with `ALAN_PHASE3B3_CHILD=1`, `ALAN_PHASE2_DB=<disposable copy>` and `ALAN_PHASE3B3_OUTPUT=<new directory>`; later compatibility tests must never point it at an archive in place. The binary is retained locally, not committed; reproduce from the exact base at an isolated checkout. Build path/toolchain affect the hash and must be recorded on rebuild.

## Executed baseline and construction evidence

The unchanged tests `TestPhase3B3ReciprocalOrdinary/s2-handshake` and `TestPhase3B3ReciprocalBindings/wrong-prior-selection` passed with fresh exports. Both service IDs are `["a",1,4,0]`, grant `["b",1,2,0]`, original request `["b",1,1,0]`, ready `["b",1,5,0]`. The supported program selects ready at A2 and q at A3, then completes both actual services, both result consumptions and both terminal principals; each independent endpoint records one contact. The swap selects q at A2 and ready at A3; its first source authorization allows but nested/combined authority refuses `authority_request_not_admitted`. The genuine grant matches A4 and the original request. Both principal results remain unpaid, with zero charges, no gate and zero endpoint contacts. Quiescence in that refused run is not completion.

[Targeted baseline assertions](evidence/baseline/assert-baseline.py) checked these actual exported identities, admissions, selecting history rows, claim, reasons, paid states, charges and endpoints. This small assertion is specific to the pair and consumes public inspection; it is not the independent Phase 3B.4 raw-policy oracle. Full original records and SQLite stopped bytes are in the archive. Old model exports were not changed; the controller independently checked their compressed and decoded hashes in [the preservation manifest](evidence/baseline/preserved-model-manifest.json).

Runtime fixture commit `3f83076` adds only `internal/durabletest/phase3b4_construction_test.go`. Tests ran on the exact file hash retained in [construction verification](evidence/construction/verification.json), against the unchanged production base; exporter HEAD therefore records the base plus that separately hashed fixture, not a claim of a clean commit before it existed. The later commit preserves the tested bytes.

| New public input | Inserted committed reductions | Actual service / planned successful reply | Actual v1 result |
|---|---:|---|---|
| Original adjacent | 0 | A4 / A7 | completed |
| q → ready → service | 0 | A4 / A7 | `authority_request_not_admitted` |
| q → ready → set → service | 1 | A5 / A8 | `authority_request_not_selected` |
| q → ready → set/copy/set → service | 3 | A7 / A10 | `authority_request_not_selected` |
| Forward true path | 3 | A7 / A10 | `authority_request_not_selected` |
| Forward false path | 3 | A7 / A10 | `authority_request_not_selected` |

A notation means canonical `["a",1,sequence,0]`; the archived manifests retain full IDs. Copy is the existing `set` with a `Ref` operand, not a new opcode. Each branch path executes branch, set and branch, while occupying five code positions; the observed scratch value distinguishes both outcomes. Service sequence counts executed reductions, not code position. Successful reply locators are construction references; refusals do not invent emitted replies.

Every program is rebuilt and hashed before `Create`; the accepted lifecycle claim, both message claim copies, downstream branch PCs and dependent manifest references use the actual new locators. No retained rows are edited. Tests assert the actual service's source authorization, genuine exact accepted grant, real q/ready selections, expected refusal reason, unpaid obligations and endpoint contacts. Six cases passed in normal and race modes; scoped integration Staticcheck and formatting/diff checks passed. First sandbox baseline failed because loopback bind was denied; a fixture compile error (`int` versus `int64` endpoint count) preceded its first execution. Both are retained apparatus/development failures, not semantic RED evidence. No TDD production change occurs in Task 1.

## Claims and remaining work

Public pre-provision construction is established for these finite instances. V1 refuses the intended nonadjacent cases at the selection layer. Derivation from existing grant/admission/observation/history fields is a source-supported candidate, not an implemented or proved result. Independent contract/construction review must precede production and principal oracle work; relation/model review remains separate.

The [contract](02-contract.md) predeclares the positive domain, version mapping, exact lookup, missing/corrupt distinction, prospective cuts/attacks and cost guard. It separates historical binding, current authority and restricted progress. Existing `readAuthorityPrefix` scans remain outside any constant-time claim. Whole-Store crash mapping, separate-grant revocation projection, repeated contacts and applied-without-local-result correspondence remain open.

Manuscript implications were checked in Abstract, §6.11 and Conclusion in both mathematical representations. They still describe v1 adjacency accurately. Task 1 adds no production semantic change or new theorem, so the existing manuscript and bibliography remain unchanged. Task 6 must reconcile a qualified candidate result in Markdown/LaTeX and rebuild/inspect the matching PDF; no early candidate-success statement is inserted here.
