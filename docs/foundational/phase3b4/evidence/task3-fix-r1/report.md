# Task3 runtime correction — round1/5

Status: corrections implemented and tested; **independent re-review pending**. The retained [round1 review](review.md) remains REJECT and is not superseded by implementer checks. Runtime base `5b68180bc18a5c7ee3f5be40d77b251e054631a9`; final runtime `9e9c8078a0bd336a33bd1a03c85e3233249af625`. Research evidence base `a9125eaee96a73fc464b7879804aad1cc326fbd7` retains the original Task3 result. No reviewed commit was rewritten.

## P1: v1 dormant claim compatibility

Commit `c03cd09c2cba1d0bee46a6c7bf125817cb001680` restricts the new program-load literal lifecycle claim preflight to v2. V1 retains its historical interpretation boundary: ordinary claim payload strings are not interpreted until an instruction emits an authority operation. Existing `e.Authority` control metadata validation remains exact for both profiles; emitted v1/v2 claim validation is unchanged.

Contract clarification: §2 distinguishes literal program **control metadata** (`e.Authority`) from lifecycle claim **payload values** interpreted at the authority-operation boundary. Its explicit v1 preservation requirement therefore permits dormant malformed or wrong-version claim text in a v1 artifact, including an unreachable suffix. This clarification does not relax actual authority operations or v2 preflight. Frozen contract bytes are unchanged, SHA256 `3b3662639f22adf9bfe3c2e5297afc7b5b6f4a389c6e8be4cb6a2f7deca9f3f3`.

Before producing the compatibility fixture, all 307 Go/module files in the reviewer's reconstructed base tree were compared byte-for-byte with `git show 3f8307664b6543331a41c943e91fb5bac8b376ef:<path>`; no mismatches or additional Go files existed. [Base source map](base-source.json), [exact generator](review-v1.go.txt), and [creation command/output](base-create.json) retain provenance. The unchanged public Reciprocal helper appends an unreachable malformed lifecycle claim before provisioning and rehashes that input. Genuine base Create and Replay both returned nil. No SQL repair or retained program rewriting occurred.

The runtime fixture `internal/executor/testdata/v1-dormant.tar.gz` retains the complete DB/WAL/SHM/lock set; [provenance](v1-dormant-provenance.json) hashes every component and archive. Tests extract complete disposable copies only. New Inspect/Replay/Open accept the original base Store; read-only validation preserves database bytes and Open preserves program records. Current v1 Create also accepts pending malformed/wrong-version text, but its actual emission rejects without changing projection; v2 rejects these literals at Create. Existing canonical-claim, metadata/use version, legacy and format controls pass.

Causal RED: [p1-red.log](p1-red.log) shows the prior Store rejected by Inspect with `unsupported version` and both current v1 Create controls rejected before correction. [p1-green.log](p1-green.log) and [p1-race.log](p1-race.log) pass after the minimal guard change. The exact added RED tests are retained as `red-*.go.txt`, hash-checked against the command source maps.

## P2: exact index key shape

Commit `9e9c8078a0bd336a33bd1a03c85e3233249af625` replaces `pragma_index_info` with `pragma_index_xinfo` and requires exactly one ascending BINARY `observation_id` key. SQLite's auxiliary rowid entry is excluded from key counting. Existing nonunique/nonpartial/user-created checks remain.

Causal RED [p2-red.log](p2-red.log) demonstrates NOCASE replacement acceptance. [p2-key-red.log](p2-key-red.log) additionally catches a descending replacement inconsistent with the declared Create shape. Both controls pass after correction; failed Open preserves database bytes, including activation epoch. [p2-green.log](p2-green.log) retains the valid plan `SEARCH admissions USING INDEX admissions_observation_id (observation_id=?)`. No service-selection or query semantics changed.

The exact reviewer probe sources and all 12 available probe Store components are retained separately in `review-probe-stores.tar.gz` with [component hashes](review-probe-stores.json). They were copied as bytes, never opened during retention.

## Verification on corrected source

[commands.json](commands.json) links every exact command, exit, output and per-file source map. [final-source.json](final-source.json) verifies all 316 Go/module hashes against every final successful gate and final runtime HEAD. Commands set `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3a-go-build GOWORK=off`; Staticcheck also sets `STATICCHECK_CACHE=/tmp/alan-phase3b4-staticcheck`.

All four final focused test variants exited0, using `./internal/executor -run '^(TestRequestBinding|TestRequestAuthority)'`, `-count=1 -timeout=10m`: ordinary, `-race`, `-tags integration`, and `-race -tags integration`. This includes the entire binding test file, new compatibility regression, affected version controls, existing v1 refusal controls and adapter paths. Final ordinary/integration Staticcheck, changed-file gofmt and `git diff --check` exited0 with empty output. Optional frozen-old-reader subtests skip when their environment variable is absent; this round does not claim a new old-reader run.

Two initial broader ordinary/race invocations failed solely because the sandbox denied httptest loopback binding; their logs are retained as `scoped-green.log`/`scoped-race.log`, explicitly apparatus failures, not semantic RED. Both reran successfully with `require_escalated`, as did integration variants.

P3's previous successful build with cache-write warning remains disclosed in the original report/log. No unchanged full suite or build was rerun merely for P3. As scoped by the controller, Task5 retains the final full-source suite/build/module qualification with its prepared writable module cache. The earlier full suite applies to 5b68180, not the corrected HEAD.

## Self-review and remaining obligations

Reviewed both complete runtime diffs, fixture generator/provenance and archive components, and exact final source maps. Runtime is clean after the two narrow commits. No dependencies, scheduler, budgets, principal oracle, frozen contract/relation semantics, or manuscript files changed. Both corrections require independent re-review; Task4 remains held until approval. All prior Task4/5 exact K01/K02 promoted-public-program, corpus/cuts/cost, independent-validator obligations and Task6 manuscript reconciliation remain outstanding.
