# Task 1 — Operational model and focused qualification

## Status and source

Complete. Branch `codex/phase3b5-grant-revocation`, base `8971eb11a2b6b9cca5a16bdedb47fc1c43838395`. Implementation, tests, proof sketches and evidence are in `2107c44c16eb4a8cb7f67122d734cd995cb953ec`, tree `070861db94c04820cce482249edb3e98b0e69014`. Verified the executable source hashes from `docs/foundational/phase3b5/evidence/task1/source-sha256.txt` again after the commit. No remote operations. Controller plan/corpus work was outside the staged paths.

## Result

The exact opt-in configuration `S2/SeparateGrantRevocation` uses separate, single-attempt, request-bound BA issuance and exact b-root revocation, with AB under independent `B/g1:c`. The selected historical first path and issuance status persist. A new permit cannot be admitted through revoked b. An existing BA permit can apply exactly once after revoke or disappear on A's crash. The AB positive witness consumes a different opaque authentic result from BA, and explicitly observes B's original serial service guard before BA is consumed. Durable consumption survives revocation, receipt release and both agents' recoveries.

The full exploration has 229,834 states and 1,807,246 arcs, with zero operational, descriptor, custody or authority violations. There are 32,768 complete states, 203,833 states with a safe completion suffix, 26,001 safe states with no completion suffix, and 188,614 cyclic states across 75,070 cyclic SCCs. The candidate completion predicate agreed at every state with the operational explorer's reverse reachability. The shortest universal completion counterexample is immediate b revocation. The preserved witness report includes early refusal, sticky refusal after accepted issuance, revocation between authorization and admission, a revoked BA application, the crash of A holding a revoked BA opportunity, and an unfair retry cycle. Both new witness cuts were checked against full successor event descriptors and the exact resulting state.

## Commands, direct exits and evidence

Go environment used for tests/build and Staticcheck:

```sh
export GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off
export GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b5/go-build
```

`go version` exited 0 with `go version go1.27.3 darwin/arm64`. `staticcheck -version` exited 0 with `staticcheck 2026.2.2 (0.8.2)`. The table's relative evidence filenames are under `docs/foundational/phase3b5/evidence/task1/`; test and Staticcheck streams were redirected together to the corresponding `.log` then archived with `gzip -n`.

| Command | Exit | Evidence |
| --- | ---: | --- |
| `go test -count=1 -v ./research/authorityservice/cmd/export -run '^TestExportSeparateGrantRevocation$'` | 1 | `red-export.log.gz`. New configuration gave exit 2, unknown configuration, before opt-in routing. |
| `go test -count=1 -v ./research/authorityservice -run '^TestSeparateGrantRevocation'` | 1 | `red-model.log.gz`. With configuration admitted to avoid an invalid-config panic, the tests reached real transitions and observed missing separate grant/revoke and wrong AB tuple. No-inline already passed and was retained as a guard. |
| `go test -count=1 -v ./research/authorityservice -run '^TestSeparateGrantRevocation'` | 1 | `red-cuts.log.gz`. With transitions implemented, the two missing named witness cuts were the only failures. |
| `go test -count=1 -v ./research/authorityservice ./research/authorityservice/cmd/export -run '^Test(SeparateGrantRevocation\|ExportSeparateGrantRevocation)'` | 0 | `green.log.gz`. All focused model and command tests, including full new exploration. The backslash in this table escapes Markdown's pipe delimiter. |
| `go test -count=1 -v ./research/authorityservice ./research/authorityservice/cmd/export -run '^Test(ApplyRequiresFreshAuthorityAdmission\|DecisiveCandidate\|CompleteGraphs\|RevocationTemporalBoundaries\|LateGrantCannotRewriteRefusal\|ReleasedCustodySurvivesBothCrashes\|ExitDisposition\|OutputFailureIsApparatusError\|DeterministicExport)$'` | 0 | `old-model-tests.log.gz`. Includes all sixteen old complete graphs. |
| `go build -buildvcs=false -o /tmp/alan-phase3b5/task1-export ./research/authorityservice/cmd/export` | 0 | Built binary at the indicated path. |
| `staticcheck ./research/authorityservice ./research/authorityservice/cmd/export` | 1 | `staticcheck-cache-blocked.log.gz`. Default Darwin cache blocked by the sandbox. |
| `XDG_CACHE_HOME=/tmp/alan-phase3b5/task1-cache staticcheck ./research/authorityservice ./research/authorityservice/cmd/export` | 1 | `staticcheck-xdg-cache-blocked.log.gz`. Darwin still selected the default cache. |
| `STATICCHECK_CACHE=/tmp/alan-phase3b5/task1-cache staticcheck ./research/authorityservice ./research/authorityservice/cmd/export` | 0 | `staticcheck.log.gz` is empty. |

Export commands executed with the built binary:

```sh
/tmp/alan-phase3b5/task1-export > /tmp/alan-phase3b5/task1-old-default.json
gzip -dc docs/foundational/phase3b2/evidence/model.json.gz > /tmp/alan-phase3b5/task1-archived-model.json
cmp /tmp/alan-phase3b5/task1-old-default.json /tmp/alan-phase3b5/task1-archived-model.json

/tmp/alan-phase3b5/task1-export -config S2/SeparateGrantRevocation -assert-authority-safe > /tmp/alan-phase3b5/task1-model-first.json
/tmp/alan-phase3b5/task1-export -config S2/SeparateGrantRevocation -assert-authority-safe > /tmp/alan-phase3b5/task1-model-second.json
cmp /tmp/alan-phase3b5/task1-model-first.json /tmp/alan-phase3b5/task1-model-second.json

/tmp/alan-phase3b5/task1-export -config S2/SeparateGrantRevocation -assert-universal-completion > /tmp/alan-phase3b5/task1-model-universal.json 2> docs/foundational/phase3b5/evidence/task1/universal.stderr
cmp /tmp/alan-phase3b5/task1-model-first.json /tmp/alan-phase3b5/task1-model-universal.json
```

Universal completion assertion exited 1 with the explicit falsification diagnostic. All other commands in this export block exited 0. The SHA-256 file `export-sha256.txt` records old default/archive equality at `180dc509fcf2a017a241e0c3cfb2a22a7f4f353a0731737c897d053d5feed46e` and first/second/universal equality at `bb2b3fa469b82a5526691375b474c6e2984d2f532e02fda4e76bb069da4a2f26`. The universal JSON and archive duplicate are represented by the same bytes in `model.json.gz` and the existing Phase 3B.2 archive respectively.

```sh
test -z "$(gofmt -l research/authorityservice/{model,explore,revocation_test}.go research/authorityservice/cmd/export/{main,revocation_test}.go)"
git diff --check -- research/authorityservice docs/foundational/phase3b5/02-model.md docs/foundational/phase3b5/evidence/task1
gzip -t docs/foundational/phase3b5/evidence/task1/*.gz
git diff --cached --check
shasum -a 256 -c docs/foundational/phase3b5/evidence/task1/source-sha256.txt
```

All five checks exited 0. The default-sandbox commit attempt exited 128 because Git's common worktree metadata is outside the writable roots. The same scoped commit with the approved tool escalation exited 0. The index manifest was inspected before commit and the committed path list after it.

## Self-review and remaining qualification

The staged diff was restricted to the opt-in operational branch, its exporter selection and cuts, focused tests and Task 1 evidence. Existing enum values, field sets, config list/order and old successor order were preserved. The archived old authority-service default export supplies a full byte regression. Exact cut names are `ba-application-after-revocation` and `revoked-ba-opportunity-crashed`. The last crash event's obligation is `NoObligation`; actor A and the predecessor BA permit establish which service opportunity it clears. Model revocation era and historical grant/witness revision remain distinct.

No unresolved Task 1 semantic obstruction. The complete graph and descriptor checks share operational source. Independent ledger/all-edge qualification, an independently computed completion predicate check, runtime trace correspondence, all-four export regression, full ordinary/race qualification and manuscript reconciliation remain later tasks. The Task 1 model records no separate debit field and introduces no argument about runtime debit behavior.
