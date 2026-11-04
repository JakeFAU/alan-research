# Phase 3B.4 — frozen prospective covering family

Status: **predictions frozen before candidate runtime or new model enumeration**. All v2 rows below are unexecuted. This is a covering family, not a Cartesian product or a claimed count of passing tests. IDs refine [contract §8](02-contract.md) without changing it. [Relation](04-model.md) separates S/H history, current V/R, restricted C, physical evidence and custody. Source-level sketches are not candidate runtime verification. The prior six Task 1 construction runs used v1 only.

## Shared bounds, identity declaration and verdict channels

Two initial logical principals A=(a,1), B=(b,1), original fixed roots/revisions1, one hop, one trusted Store owner. Principal IDs are AB=`["a",1,1,0]` and BA/q=`["b",1,1,0]`. Standard grant=`["b",1,2,0]`, ready=`["b",1,5,0]`; with n inserted A reductions, BA service=`["a",1,4+n,0]`, planned successful reply=`["a",1,7+n,0]`; AB service=`["b",1,9,0]`, reply=`["b",1,12,0]`. Each declaration is verified against genuine program/commit facts, not accepted as authority. Namespace, all program hashes, roots/permissions/revisions and complete IDs remain in every record.

Positive family: raw1d public RequestBinding construction, v2 formats15/16 with exact nested v2 authority tags; new provision for each changed program, all claims and PCs corrected before hash/Create. S2Handshake, TerminalHandshake=true, two distinct loopback idempotent endpoints, same complete service tuples throughout a case. Standard byte limits: each MaxProgramBytes/MaxDataBytes/MaxObservationBytes/MaxOutputBytes=1MiB; root MaxEffects1/MaxInFlight1 each; Run MaxSteps256/MaxDispatch256/MaxInFlight2/ExternalTimeout5s and15s outer context. Record measured serialized bytes and counters, not assumed headroom. Run ceilings are per invocation. Controlled schedules retain their actual public transitions/cut events and elapsed-time limits separately. No live model endpoint or new dependency.

Admission records have no global timestamp/order field: check genuine causal admission prerequisites and exact descriptor/observation joins, retain actual arbitrary transport IDs, and report a justified transport linearization rather than inventing chronology. Later duplicate admissions need not precede the original selection. GrantDecisionSequence is a lifecycle sequence, distinct from composition fact order and A/B reduction sequences.

Record separately for every leaf: valid raw history versus corruption/version error; historical S/H; nested permission/first reservation/debit; fresh gate decision; independent endpoint contacts/applications; authentic principal AB/BA result consumptions; auxiliary q2 status; ACK drain; terminal A/B; Run disposition/quiescence; supported C projection or explicit unsupported reason. Safe denial never counts as completion. For a well-formed policy denial, raw validity is expected true even when progress fails. For hostile retained contradiction, raw validity is expected false. Do not turn missing capture, panic or an unrelated earlier rejection into a target-policy result.

## Semantic covering leaves

Unless a row overrides it, use n=0, request→ready order, grant before first decision, no revocation/owner loss and sufficient V/R. P01 and P02 each compare freshly provisioned v1/v2; all remaining new rows use v2. P03 is an alias of the candidate P02 run, not an additional counted experiment. P08/P09 isolate ordering with controlled transitions; automatic P02 remains mandatory independently.

| ID | Fixed variation and actual coordinates | Predicted outcome and falsifier |
|---|---|---|
| P01 | ready at A2, q at A3, e=A4; v1 and v2 | Both authorize and both authentic principal results consumed. |
| P02 / P03 | q at A2, ready at A3, e=A4; n=0; v1 and v2 automatic Run | Preserve v1 `authority_request_not_admitted`; candidate allow and automatic reciprocal completion. V1 denial is not candidate RED/GREEN evidence. |
| P04 | qA2, readyA3, literal scratch setA4, e=A5; n=1 | Allow; replyA8, AB B9/B12; both consumptions. |
| P05 | qA2, readyA3, set literalA4, set/Ref copyA5, set literalA6, e=A7; n=3 | Allow; replyA10; protected operands unchanged, no extra effect debit. |
| P06 | n=3 executed, ForwardBranch=true, LocalSteps=0; PCs3→4→5→8 | A4/A5/A6 then eA7, scratch1, replyA10; allow and both consumptions. |
| P07 | Same, ForwardBranch=false; PCs3→6→7→8 | A4/A5/A6 then eA7, scratch2, replyA10; same authority, different visible path. |
| P08 | Control B's lifecycle issuance/acceptance after q selection A2; ready follows genuine acceptance; eA4 | Allow; original q commitment precedes issuance, c precedes e. |
| P09 | Genuine grant accepted before q selection A2; then readyA3/eA4 | Allow; c need not precede grant. Retain exact composition fact order in both P08/P09. |
| P10 | qA2→readyA3→q2A4→eA5, matching g(q,eA5) | Allow bound to q/cA2 after later equal-payload same-peer selection. Both principals paid; auxiliary q2 separately unpaid unless genuinely served. See explicit map below. |
| P11 | qA2→readyA3→`set peer_request="overwritten"` at A4→eA5; service/reply live operands unchanged | Allow bound to q/cA2; replyA8 and both actual result consumptions. Any dependence on dead application register falsifies intended binding. |

For P10, insert ordinary q2 at B6 after readyB5, using the same source/target/obligation and byte-equal payload as q but a distinct full original ID `["b",1,6,0]`. A receives q2 after q and ready, binding only dead peer_request/peer_claim (or declared scratch) fields. The byte-equal payload may contain q's textual request_id; it does not make q2 equal to q. Rebuild B's shifted suffix: BA reply selectionB7, paid branchB8, AB selectionB9, AB serviceB10=`["b",1,10,0]`, resultB11, branchB12, AB replyB13=`["b",1,13,0]`; BA service/replyA5/A8. g remains B2 and names q/B1 and exact eA5. The principal domain remains AB/A1 and BA/B1; q2/B6 is an explicitly separate auxiliary obligation. Ordinary automatic P10 may select q first by program/schedule only if demonstrated; controlled exact observation selection is labeled, never attributed to automatic Run. No new C claim for arbitrary extra-original histories follows from the raw/custody success.

| ID | Controlled negative leaf | Predicted actual layer |
|---|---|---|
| N01 | Same q2/B6 construction; q admitted but unselected, readyA2 then only q2 selectedA3; g names q/eA4. B service/reply shifts toB10/B13. | `authority_request_not_selected`, no BA exposure. Keep q2 and q admission distinct. |
| N02 | Program emits genuine eA2 before selecting qA3, later readyA4; exact grant names eA2 and is accepted before delayed D after q selection. | `authority_request_not_selected`; empty request witness fields; later D cannot turn cA3 into c<eA2. |
| N03 | eA2 emitted after AB originalA1; D before any genuine grant; later grant naming q/eA2 accepted. | Immutable `authority_grant_missing`; no later repair/redebit. Any later selection cannot change original reason. |
| N04 | Grant accepted; eA2 and D before q selectionA3, then genuine later selection. | Immutable `authority_request_not_selected`; empty original selection witness persists through Replay/Open. |
| N05 | Genuine q selected, then revoke before first decision (n=0) | S and historical binding coordinates survive; current authority refuses at inherited revocation boundary, no exposure. |
| N06 | Genuine allow/debit then revoke before any fresh gate (n=0) | Original allow/debit retained; fresh gate refuses, no new contact. |
| N07 | Genuine gate before revoke, known result retained and consumed after revoke | Old opportunity/result retains meaning; no fresh permission. Endpoint application separately proved. Outside nonrevoking C. |
| N08 | New n=1 program emits actual A5 but genuine grant still names old A4 | `authority_service_mismatch`, not request-selection failure. |
| N09a | Unequal forward path reaches textual service after2 versus3 reductions; fixed grant forA7; choose2-step path yieldingA6 | Exact service mismatch; do not normalize A6 toA7. |
| N09b | Inserted write changes live target/obligation operand with grant unchanged | Actual action/service validation refuses at its inherited boundary; not a permitted scratch step. |
| N09c | Forward branch skips designated service | Outside progress family; no invented service, no completion claim. |
| N10a/b | Insufficient relevant MaxEffects / occupied MaxInFlight | Resource boundary refuses/blocks as applicable; distinguish permission allow from resource nonallow. Use genuine occupied reservation for in-flight control. |
| N10c/d/e/f | Source/scope closure; retirement; changed holder/recipient generation; changed revision before fresh decision/gate | Inherited source/authority boundary rejects current opportunity; historical S remains about old generation. Each is a separate leaf, no invented TTL. |
| N10g/h | Run MaxSteps=3 before required A4 service / explicit byte limit below a real intermediate projection or program | Actual Run/byte validation error, not binding refusal; retain measured limit and offending bytes/counter. |
| N11a/b/c/d | Cartesian tuple mixing / ambient borrowing / result-as-grant / attempted own-path fallback | Prior negative policy preserved; no invented tuple, consent, grant or alternative path. Each separates syntactic input error from well-formed claim refusal. |
| N12 | Second distinct same-payload service selects first service's grant | Exact operation mismatch; no second use inferred from q selection. |

N02/N03/N04 are narrow controlled programs, not RequestBinding constructor outputs. Their planned successful reply coordinates are irrelevant because their principal service is refused; no uncommitted reply is reported as evidence. Future implementation must publish all program hashes, genuine exact grant/intent facts and event ordering. If public input validation makes one construction impossible, retain that obstruction and revise its test construction through review; do not fabricate a retained parent and call it runtime-produced.

## Retained-evidence attacks and checker-source mutants

A01 splits wrong namespace, holder/recipient generation, root, revision and complete tuple into distinct substitutions. A02 splits wrong q/o/c/e and selected-history/consumption mirrors, missing genuine admission and missing source authorization. A03 uses future grant/selection/source facts at an earlier decision; A04 coherently forges attempt/reservation agreement while contradicting original operation; A05 splits metadata/PRAGMA/profile-only, nested claim/use/decision/witness/gate version mixes and coherent multi-field downgrades. These are attack axes, **not** a Cartesian product. For each applicable leaf start with a named genuine parent, mutate the smallest complete dependent set and recompute all affected descriptors/digests/index mirrors; retain original and mutant plus exact guard reached. Malformed external claims remain validation errors; well-formed false claims use the actual policy-refusal layer. Archive corruption is never a repaired denial.

A02-filter adds a specific actual-reduction falsifier: start from a genuine known-result/receipt-selection parent, coherently change the retained program receipt receive filter Operation to the empty literal and rehash the dependent artifact/deployment/history descriptors without changing the selected receipt's real nonempty Operation. Runtime Reduce compares the empty filter exactly and rejects; the sound independent raw checker must reject at receive_filter. Existing custody.go currently treats this empty value as a wildcard; source inspection predicts a possible false acceptance, but no hostile run or full false-acceptance demonstration has occurred here. Also require selected-observation-on-set/branch/emit/end rejection, preserving genuine receive-only consumption. Keep other filter coordinates/source/generation valid to isolate the Operation guard. Qualification must correct the checker independently rather than importing Reduce, and retain the pre-fix semantic verdict if the attack reaches it. This refines A02's real-reduction premise; no new runtime policy or contract change.

Predeclared source mutants are mandatory future Task5 work; none is executed in Task2. Build success is a prerequisite, not a kill. The assertion must demonstrate the specified wrong semantic verdict with all other relevant joins satisfied; redundant/inapplicable deletion is reported and replaced with a reviewed nonredundant join.

| Mutant | Frozen parent and exact unsound change | Required semantic assertion |
|---|---|---|
| M01 omit strict c<e | Genuine N02 with cA3/eA2 and delayed D. In a copied export coherently encode claimed allow/witness/attempt/reservation/debit/gate while preserving genuine program selection order. Remove only strict historical ordering guard from checker. | Sound checker rejects at selection-before-service; mutant accepts same hostile history. A separate remaining causal-order guard killing it makes this deletion redundant, not a qualified mutant. |
| M02 latest receive | P10 and N01, with explicit principal/auxiliary maps. Replace exact grant-original selection/admission join by latest same-peer receive **and its supplying admission**, removing the semantically replaced original association rather than leaving a redundant exact-q check. | P10 gives false refusal of a valid q-bound service if latest q2 replaces q; N01 coherent false-allow export gives false acceptance when only q2 selected. Require both witnesses; neither obligation-suffix custody error nor parser failure counts. |
| M03 omit generation join | Genuine auxiliary message from B1 explicitly requests target A generation2; dense contradictory admission claims Accepted/Bound=A1, matching observation/result/index mirrors. Omit requested target generation versus Bound.Generation comparison. | Sound checker rejects at admission generation binding; mutant accepts contradictory auxiliary admission. Keep actual message intent/program target2, so all other joins are coherent. No generation2 principal/root is introduced. |
| M04 skip current gate revocation | N06 with authentic allow/debit, revoke before fresh gate; rewrite only gate verdict and its dependent digest to allowed/authorized. Skip current revocation specifically at fresh-gate derivation. | Sound checker rejects false fresh gate; mutant accepts. Historical first decision remains lawful; skipping a different original-time check is not this mutant. |

Additional exact-custody falsifier: q/q2 services yield numerically equal real results; a reply names principal q but carries q2's genuine service/result lineage. Sound raw custody rejects exact original-service/reply mismatch. A grant genuinely issued for q2 with authentic earlier q2 selection is allowed on its own merits and is **not** an invalid substitution by itself. Extra paid q2 never repairs unpaid principal q.

Initial, pre-grant, pre-selection, pre-attempt, admitted-only and post-result-unconsumed exports require bounded explicit mutation applicability diagnostics. The retained old incomplete-prefix generator panic is apparatus failure, not negative policy evidence. No mutation generation or principal checker is implemented here.

## Owner-loss covering family (concrete, outside C)

Uninterrupted P02 supplies loss count0. C01–C06 andC08 use exactly one acknowledged owner SIGKILL then a fresh process on the original Store. C07 has exactly two consecutive losses. Only C03/C07 use n=3; others use n=0 nonadjacent q→ready. Cuts are controlled public-transition/barrier schedules, not assertions that automatic Run naturally hits them.

| Cut | Real acknowledged boundary | Frozen recovery assertion |
|---|---|---|
| C01 | q accepted admission, before q selection | Immediately after Open q remains unselected; later real commit alone establishes c. |
| C02 | q selectedA2, before readyA3 | Same q/o/c; after ready then eA4, candidate allow under V/R. |
| C03 | Grant/q selection plus all3 local steps committed, before eA7 | Exact immutable prefix and future reference retained; genuine later eA7, not preexisting fictional e. |
| C04 | eA4 emitted, before first authority decision | Exact e retained; genuine later first decision uses its actual prefix. |
| C05 | First allow/debit, before fresh gate | Old allow/debit retained, new epoch requires current gate before worker release. |
| C06 | Known original result retained, before selection/consumption | Select retained actual result, no completed-call reissue. |
| C07 | First loss after qA2 selection, renew; second loss after3 local steps before eA7, renew again | Both renewals checked separately; q/o/c/grant/program IDs/charges preserved; later genuine eA7. |
| C08 | Revoke after allow, before recovery gate | H/allow/debit retained; fresh gate refuses under current revoked authority. |

Capture list order is explicit: ordinal, predecessor capture, event/cut label, owner PID/identity, acknowledged SIGKILL, pre-Open bytes, post-Open bytes, epoch transition and hashes. Every renewal requires its own before/after pair; missing pre-Open data is RenewalUnchecked, not inferred success. Keep all original identity, program/profile/root/permission rows, real c and consumption, attempts/reservations/debits, pending results/ACKs and independent contact/application records. Separately report merely reserved versus physically applied without local result. Idempotent repeated contact/single application and generic unknown non-retransmission are different endpoint contracts; no paired-agent crash mapping is inferred.

## Compatibility and cost guards

K01 preserves unchanged authored `TestPhase3B3CartesianOperations/r1-add-one` under v1 formats13/14 (receipt-conflict promotion by genuine writer). K02 newly provisions its coherently versioned v2 counterpart under15/16 with artifact semantics `alan/1e/v1`/author-v1 metadata; input versions and hashes change before Create, retained records never repaired. Both require exact tuple/service/result, one endpoint application and Create/Inspect/Replay/Open acceptance in normal/race qualification. Raw1d-only rejection is a failure. Invalid semantics/metadata and unsupported feature combinations still reject. Old reader must reject15/16 before epoch renewal; old13/14 witness bytes/meaning remain. Coherent v2 nonadjacent allow downgraded to v1 must fail old adjacency validation. Full unsigned history rewriting remains outside authenticity.

Cost-only family lengths **0,1,8,32** uses same machine/toolchain/frozen source, v2 profile, request→ready, same action shape, exact corresponding serviceA(4+n)/replyA(7+n), no owner loss or extra q2. Perform5 independent observations per length for construction, first decision, fresh gate, Replay and Open (Open on disposable copies; setup outside measured operation). Record all raw observations/variation; no arbitrary threshold or speedup is predicted. Record allocations only where measured by the harness, reductions/bytes/admission counts, and query-plan evidence for exact nonunique `admissions_observation_id(observation_id)` lookup. The added lookup must not scan all history/admissions. Inherited composition-prefix/projection work remains separate and may scale; no total O(1) claim. Cost guards do not expand semantic claims to every combination at lengths8/32.

## Verification and qualification handoff

Task2 runs byte/hash preservation, review-freeze identity, local-link and whitespace checks only. No new Go model code means no TDD RED/GREEN and no new deterministic model export. The controller explicitly deferred heavy unchanged suites to final frozen-pair Task5 qualification. Existing exports must remain byte-identical now, with decoded gzip hashes also checked. The future final qualification must rerun deterministic export/audit twice into fresh temporary paths and compare both with the preserved decoded exports; it must never overwrite historical files.

Known existing commands (not run by Task2), from research README and model README:

```sh
export GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3a-go-build
# Each repository, from its exact frozen source and with recorded cwd:
go test -count=1 -timeout=60m ./...
# Research affected/model regression and final unchanged export:
go test -count=1 ./research/authorityservice/... ./research/delegationcorrespondence/...
go run ./research/authorityservice/cmd/export -config all
go run ./research/authorityservice/cmd/audit
# Required research gates, scope adjusted to the actual final edits:
go test -count=1 -timeout=60m -tags integration ./...
go test -count=1 -timeout=600m -race -tags integration ./...
staticcheck ./research/authorityservice/... ./research/delegationcorrespondence/...
go vet ./...
go vet -tags integration ./...
go mod verify
```

Record source/toolchain/env/argv/exits and output destinations; fixture export environment must point to fresh roots and old-reader compatibility must use the frozen old binary. Runtime final full integration/race/build/module/Staticcheck and inherited vet/format gates use its actual documented commands (`docs/extraction.md`), not invented test names. New candidate test selectors/capture/checker CLIs do not yet exist and must be published with source before their qualification; never claim a future command ran. Loopback tests require authorized execution outside the restricted sandbox.

Exit condition for this task is an independently reviewable relation and predeclared falsifiers. Candidate implementation, automatic completion, model trace checks, mutations, whole-Store continuity, compatibility, cost measurements and manuscript reconciliation remain later obligations. Subsequent semantic changes to these predictions require an explicit revision and retained previous version before requalification.
