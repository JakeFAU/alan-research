# Task 3 — versioned historical request binding

Status: **DONE_WITH_CONCERNS**, ready for independent runtime/recovery and compatibility review. The runtime implementation and focused qualification are complete; this is not completion of the Task4/5 campaign or a general correctness proof.

## Source and delivery

Runtime worktree `/private/tmp/alan-phase3b3-runtime`, branch `codex/phase3b4-explicit-request-binding`, base `3f8307664b6543331a41c943e91fb5bac8b376ef` (production baseline `1fdd4c9e1720bf60b344914bcbc3d0938ba8d678`). Local commits, not pushed:

- `158bc91c7ff8ec6bee01388d4cf6aacdb9ee7d08` — Add opt-in historical request binding authority.
- `91e8c57147b19a19136c0db2a84c317755e71627` — Add finite request binding delegation example.
- `5b68180bc18a5c7ee3f5be40d77b251e054631a9` — Verify retained binding results survive revocation (test only).

Final runtime HEAD is `5b68180bc18a5c7ee3f5be40d77b251e054631a9`; worktree clean. Final frozen gates below ran on that exact committed source. [Source inventory](frozen-source.json) hashes every tracked Go/module file. Its canonical source-hash-map digest is `c45638fab55f385dfd45aeff59f0a8bce67203f097e297d6aa2ba0da420d5ccb`; production-only map digest is `88e2f0fa98d0c1bcfd0d201d061d3932bf49910c38b799815562ab694fa231d2`. These are map digests, not Git commit hashes.

Consumed reviewed research contract `02-contract.md` SHA256 `3b3662639f22adf9bfe3c2e5297afc7b5b6f4a389c6e8be4cb6a2f7deca9f3f3`, relation `04-model.md` `09f715a6f4aa7f2ac4be3dbfd7444e5495d7947d72867942215698906b659f37`, and matrix `05-matrix.md` `32367cd857faedc0410f5e9a03f9b6a10b30407ecd8a903e35281e0567c69a31`. Their semantics/status bytes were not changed by this task. The controller supplied [independent relation approval](../reviews/task2-r1/review.md), retained byte-for-byte with its [scope/source manifest](../reviews/task2-r1/manifest.json); no new independent reviewer was spawned by this implementer.

## Question, boundary and result

Can a new opt-in policy authorize the exact service named by a genuine grant after the grant's exact original request was selected earlier, even with ready/local reductions intervening, without reinterpreting v1 refusals or treating historical selection as fresh authority?

Focused implementation evidence supports that bounded behavior. V1 keeps `RequestAuthorityVersion` and formats13/14. V2 uses `RequestAuthorityBindingVersion` and formats15/16. The v1 selection evaluator remains the old adjacent-selection branch. V2 derives q from the selected grant, derives o from q's canonical logical ID, checks same-Run positive c strictly before e.Sequence, and validates history/consumption/revision, actual receive reduction, committing-intent mirrors, accepted exact-message admissions and source authorization/fact ordering. Lifecycle decision sequence, local reduction sequence and composition fact order are never compared across sorts.

Admission has no invented global ordinal. The genuine atomic admission/selection writers provide the causal prerequisite; exact descriptors, source authorization and selecting execution anchor it. Duplicate arbitrary nonempty UTF-8 transport IDs are supported. The new query is exactly `SELECT id FROM admissions WHERE observation_id=? ORDER BY id`; all matched supplying admissions are checked with their actual transport IDs. No ACK/latest-receive/application-register shortcut is used.

The v2-only non-unique `admissions_observation_id(observation_id)` index is created only at provisioning and checked by readers. Missing/incorrect shape rejects rather than repairing on Open. The observed plan is `SEARCH admissions USING INDEX admissions_observation_id (observation_id=?)`; [query-plan/old-reader output](index-and-old-reader.log). This does not imply total constant-time authorization: inherited prefix/projection work and the actual receive recheck remain.

The historical receive recheck uses canonical typed program/predecessor/successor/observation byte sizes and the actual empty receive-intent output; it does not impose default operational limits on an already valid larger history. Commit and whole-ledger validation retain caller Limits. A >1MiB state under 2MiB limits and a compact raw artifact with omitted zero fields passed.

Current gates retain original operation/attempt/reservation/debit reconstruction and existing lock/worker-release lifetime. Selection adds no authority, budget, scheduler, registry, mutable request cache or dependency. Tests preserve missing-grant/selection refusals across later selection/grant and Replay/Open. A completed delegator before service still permits the original decision/gate; explicit revocation blocks fresh gates. A genuine previously retained result remains selectable after revocation.

## RED/GREEN evidence

[Evidence manifest](manifest.json) retains original development logs, bounded source mutants, exact reproduced RED commands/exits/source hashes and final GREEN outputs.

1. Profile RED: the literal v2 profile failed Create with `unsupported version`, not a missing symbol. Version routing made it pass.
2. Semantic adjacency RED: a genuinely provisioned v2 request→ready→service history reached `authority_request_not_admitted`. The v1 control retained that refusal. The new exact selection predicate made v2 allow and retain the original c/o.
3. Committing-history RED: removing only q's or e's history emitted-intent mirror still allowed direct v2 predicate evaluation, although whole Replay already rejected the ledger. The added exact index/revision/history-intent join rejects both at the binding boundary.

The initial logs are `profile-red.log`, `selection-red.log`, `commit-join-red.log`; GREEN logs are `selection-green.log`, `commit-join-green.log` and final frozen suite logs. Initial exploratory working-tree hashes were not captured at every edit. To retain precisely reproducible causal evidence, `reproduce_red.py` replayed all three guard reversions in a disposable final-source copy; all compiled and failed the intended assertions with exit1. `red-reproduction.json` contains exact per-case source hashes, argv and mutations. These are explicit guard-reversion reproductions, not the Task5 independent oracle/source-mutant campaign. Additional public/compatibility/policy tests are regression qualification; no pre-implementation RED is claimed for every supplementary case.

Apparatus failures were not treated as semantic RED: an initially wrong SQL column in one malformed-copy test; a missing test import (`integration.log`); Staticcheck's sandbox-denied default cache (`staticcheck.log`, corrected with `STATICCHECK_CACHE`); and the first retained-result control addressing the helper's old `pending` register rather than its new revoke destination. That last control was corrected in the provision, not by retained SQL repair; the apparatus logs remain separately named.

## Changes and focused coverage

Production changes are limited to authority version routing/validation, formats/conflict promotion, profile/inspection agreement, new exact selection/index helper, and the public constant. Existing wire fields and envelope digest domains remain unchanged. The public `examples/delegation/RequestBinding` implements the exact frozen options/manifest interface with raw artifacts, S2 lifecycle-issued grants, actual service/reply locators, full program hashes and two principal request IDs. `Reciprocal` source/output remains v1. Local lengths0/1/3/8/32 provision deterministically; automatic endpoint tests cover adjacent/nonadjacent order, local1/3 and both three-reduction forward branches. The branch family uses only set/branch instructions and keeps the reviewed PC paths. README and executable example document the boundary.

Focused tests cover q versus equal-payload q2 (including q admitted but unselected), delayed decision after a future c, decision before selection followed by real later c, sticky missing grant, arbitrary duplicate transports/index plan, exact root/holder/recipient/generation/revision/service/action/path controls, missing selector/no fallback/result-as-grant, forged wire/attempt/reservation agreement, missing original debit, revoked gate/reopen, completion-before-service, retained result consumption after revoke, malformed index/metadata/PRAGMA/profile/consumption/admission, unchanged failed-Open Store bytes, and receipt-conflict promotion/readers. These are targeted controls, not all dense matrix substitutions.

Compatibility distinction: unchanged `TestPhase3B3CartesianOperations/r1-add-one` exercises the exact authored physical endpoint case under13; new `TestPhase3B4AuthoredAdapterCompatibility` uses its genuine `author.Compile` construction, changing only coherent authority inputs/hashes before Create, under15. Both pass exact service/result selection and real endpoint application checks. A separate authored fixture exercises genuine13→14 and15→16 writer promotion, Inspect/Replay/Open and original result consumption after revoke. **Exact K01/K02 public-program conflict-format counterparts14/16 remain Task4/5 obligations.** No raw-only artifact restriction was introduced.

The frozen old reader remains unchanged at SHA256 `f08fc613f6411adb005adb9bf67409bb456eff69cd39599fa9652957e1a1db36`, source `1fdd4c9e1720bf60b344914bcbc3d0938ba8d678`. On complete disposable copies only, invocation `old-reader.test -test.run=^TestPhase3B3Owner$ -test.v -test.timeout=2m`, with `ALAN_PHASE3B3_CHILD=1`, `ALAN_PHASE2_DB=<copy>`, `ALAN_PHASE3B3_OUTPUT=<fresh directory>` and inherited ALAN variables removed, rejects15/16 with exact `unsupported version` in `open.json`; copied database bytes are unchanged. No archived SQLite original was opened.

## Final verification

Go `go1.27.3 darwin/arm64`; Staticcheck `2026.2.2 (0.8.2)`. All final commands/exits/output names and exact source hashes are in [commands.json](commands.json). Tests set `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3a-go-build`. Loopback/subprocess tests ran with authorized escalated execution. Small-gate commands additionally set `GOWORK=off` and `STATICCHECK_CACHE=/tmp/alan-phase3b4-staticcheck`.

All exited0 on final HEAD:

```sh
go test -count=1 -timeout=60m ./...
go test -count=1 -timeout=60m -tags integration ./internal/executor ./internal/durabletest ./examples/delegation -run '^(TestRequestBinding|TestRequestAuthority|TestPhase3B4AuthoredAdapterCompatibility|TestPhase3B3CartesianOperations)'
go test -count=1 -timeout=60m -race -tags integration ./internal/executor ./internal/durabletest ./examples/delegation -run '^(TestRequestBinding|TestRequestAuthority|TestPhase3B4AuthoredAdapterCompatibility|TestPhase3B3CartesianOperations)'
go build ./...
staticcheck ./...
staticcheck -tags integration ./internal/executor ./internal/durabletest ./examples/delegation ./runtime
go vet ./...
go vet -tags integration ./...
go mod verify
go mod tidy -diff
gofmt -l adapters author cmd compose effects examples internal runtime schema
git diff --check
```

Focused integration/race also set `ALAN_PHASE3B4_OLD_READER=/tmp/alan-phase3b4-task1/compatibility/old-reader.test`. Format/tidy output is empty. The build exited0 but emitted a nonfatal module stat-cache write warning (`operation not permitted`) under the sandbox; its exact warning remains in `frozen-build.log`. This is not a silent build log or a semantic failure. Relevant logs are `frozen-ordinary.log`, `frozen-integration.log`, `frozen-race.log` and the other `frozen-*.log` files. Earlier superseded suites remain identifiable separately.

## Self-review and remaining obligations

Complete tracked/new file diffs and staged scopes were inspected before each commit. No unrelated source, dependencies, generated outputs, instrumentation, secrets or archived fixtures entered runtime commits; no push/merge occurred. The extra q/e commitment guard and reviewed set/branch-only construction were self-review corrections, followed by final-source requalification. Selection/source/root/gate coordinates remain separate.

Independent runtime/recovery and compatibility review is still required before Task4. Tasks4/5 retain all full raw/custody corpus, coherent multi-descriptor downgrades, all matrix leaf combinations, acknowledged SIGKILL cuts, exact K01/K02 promoted public programs, independent endpoint/capture continuity, raw-checker/source-mutant qualification, five-observation length0/1/8/32 cost campaign, and final full shared-validator integration/race on the frozen research/runtime pair. No principal raw oracle was implemented here. The source has no live subject-revision updater; N10 must preserve/report that construction boundary rather than fabricate one.

The manuscript remains the historical qualified v1 account; this task does not claim candidate research qualification or update its three representations/PDF. Task6 must reconcile the new qualified result only after the campaign/reviews. Full unsigned history rewriting, general dynamic service locators, unrestricted contexts, distributed ownership and general termination remain outside this result.

## Evidence retention

The Task3 logs, exact command/source records, guard-reversion source snapshots and implementation diffs are intentionally retained, including `.log` files otherwise ignored by the repository. The exact Task2 relation review is retained separately under `evidence/reviews/task2-r1`; its approval scope is unchanged. A tracked copy of this report is `docs/foundational/phase3b4/evidence/task3/report.md`.

The retention-only checks verified every original manifest hash, all12 final command records and all315 frozen source-file hashes. No credential-pattern candidates were found, and the complete staged scope contains only the requested evidence/review. Raw `git diff --cached --check` reports167 diagnostics solely inside literal context lines of archived `runtime-commit.diff`; those evidence bytes remain unchanged. Both archived diffs parse with `git apply --numstat` without application, and the whitespace check outside that patch data passes. [Retention checks](retention-checks.json) record the exact commands/dispositions. No runtime/model tests were added or rerun during evidence retention.
