### Spec Compliance

- **REJECT — issues found.** Task 3's exact historical-binding implementation is substantially present, but original v1 program acceptance is narrowed and a malformed v2 index can pass Open despite defeating the required exact indexed lookup. These are independently reproduced below.
- Reviewed base `3f8307664b6543331a41c943e91fb5bac8b376ef` → head `5b68180bc18a5c7ee3f5be40d77b251e054631a9`, using the supplied 2,473-line diff in one sequential pass. Read the task brief, implementer report, contract, relation/matrix premises, source manifest and retained test/command evidence. The contract/relation are requirements, not evidence of implementation correctness.
- ⚠️ Tasks 4/5 remain responsible for the exact **K01 unchanged public authored `TestPhase3B3CartesianOperations/r1-add-one` promoted to format14**, and **K02 its coherent v2 public counterpart promoted to format16**, with genuine writer conflict promotion, exact service/result/endpoint custody and Create/Inspect/Replay/Open normal/race evidence. `internal/executor/request_binding_adapter_test.go:17` covers a different authored fixture under13/14/15/16; `internal/durabletest/phase3b4_binding_test.go:16` supplies the exact v2 authored public counterpart under15. The former does not discharge the missing promoted exact public-program cases.
- ⚠️ Broader dense raw/custody attacks, coherent multi-descriptor downgrades, independent oracle/source mutants, acknowledged SIGKILL corpus, endpoint continuity, cost campaign and final frozen-pair qualification remain Tasks4/5. Manuscript reconciliation remains Task6. These are not represented here as completed or as Task3 blocking coverage expansion.

### Strengths

- `internal/executor/request_binding.go:16`: selection starts from the selected grant's exact original ID and canonical observation ID; it joins consumption/history and strict same-run c<e, authentic admissions with actual transport IDs, source authorization, commitment mirrors and a real receive reduction. This avoids a latest-receive/register/ACK shortcut.
- `internal/executor/request_authority.go:464`: the v2 selection branch is distinct from the retained adjacent-selection branch; grant/service exactness and first-decision reconstruction remain in the existing authority path. No new authority cache, budget ledger, scheduler, dependency or worker-release protocol appears in the diff.
- `internal/executor/request_binding_test.go:120`: the future-selection cases check both delayed decision and sticky prior refusal across Replay/Open. `internal/executor/request_binding_policy_test.go:38` and `internal/executor/request_binding_adapter_test.go:17` separately exercise revoked fresh opportunities and consumption of a previously retained genuine adapter result after revoke.
- `examples/delegation/request_binding_run_test.go:15`: automatic public programs cover adjacent/nonadjacent order, local1/3 and both finite branches using actual endpoints, exact result selection and Replay/Open. The constructor emits genuine lifecycle grants and adjusts locators before provisioning.
- Retained evidence: all315 source-hash entries match the current runtime checkout. Frozen ordinary/integration/race logs contain package success results; Staticcheck, vet, tidy, format and diff checks have no diagnostics. The three retained RED logs demonstrate unsupported profile, adjacency refusal and absent committing-intent mirror respectively. Guard-reversion reproductions are correctly distinguished from original chronology and future independent mutants.

### Issues

#### Critical (Must Fix)

- None found in this task-scoped review.

#### Important (Should Fix)

1. **[P1] Preserve v1 dormant program interpretation — `internal/executor/request_authority.go:578–582`.** The new literal-claim preflight runs for every nonempty authority version, including v1. Before this change, program loading did not interpret lifecycle claim payload strings; authority claims were validated when their instruction produced an authority operation. A previously valid v1 Store can therefore contain unreachable or not-yet-executed literal claim text that this new reader rejects before recovery. This changes original formats13/14 semantics even though the adjacent-selection evaluator is unchanged.

   Focused public-API witness: `/tmp/alan-phase3b4-review-v1.go` builds ordinary v1 `Reciprocal`, appends an unreachable lifecycle emit after the existing terminal suffix, puts `DelegateAuthority` and `uninterpreted dormant text` in literal payload fields, and rehashes the artifact before Create. Exact base source reconstructed in disposable `/tmp/alan-phase3b4-review-base` by copying Go/module sources and reversing the supplied diff returns `Create error=<nil>` and `Replay error=<nil>`. Head returns `Create error=unsupported version`. Head Open of the copied genuine base Store also returns `open error=unsupported version`. No retained SQL was edited to establish this positive base history. The old validator in the diff and the focused unchanged `program.go:132–192` check confirm the causal change.

   Restrict the new program-load literal preflight to v2 (or provide a version-specific equivalent preserving the old v1 acceptance rules), while retaining exact version validation when v1 authority operations are actually interpreted. Add a genuine base-produced dormant-claim compatibility fixture and require new Replay/inspection/Open acceptance without rewriting it. If the frozen contract's generic literal-validation requirement is interpreted as applying retroactively to v1, that clause conflicts with its explicit v1 preservation requirement and needs a recorded clarification; it cannot silently override old semantics.

2. **[P2] Validate the index collation — `internal/executor/request_binding.go:206–212`.** `pragma_index_info` reports the column name/count but omits collation. Consequently `CREATE INDEX admissions_observation_id ON admissions(observation_id COLLATE NOCASE)` passes the new shape checks. The query uses the column's BINARY equality, so this replacement cannot provide the required exact lookup. It violates contract§6's malformed-index rejection and no-new-admissions-scan requirement.

   Focused witness: `/tmp/alan-phase3b4-review-index.go` genuinely provisioned a fresh v2 Store through public `RequestBinding`/Create. Only its index was replaced on that disposable Store. The same metadata queries used by the validator returned `[('admissions_observation_id',0,'c',0)]` and `[(1,'observation_id')]`; `pragma_index_xinfo` exposed `NOCASE`. Python SQLite's `EXPLAIN QUERY PLAN SELECT id FROM admissions WHERE observation_id=? ORDER BY id` returned `SCAN admissions USING INDEX sqlite_autoindex_admissions_1`. Head's actual public Open returned `open error=<nil>`. This is an accepted malformed-format witness, not a claim of unauthorized service or a timing benchmark.

   Validate the key column's intended BINARY collation using `pragma_index_xinfo` (and the expected key shape), and add the NOCASE replacement to the failed-open controls at `request_binding_test.go:284`. Require rejection before epoch renewal and unchanged immutable bytes. Retain the valid default index query-plan control.

#### Minor (Nice to Have)

- **[P3] Build evidence includes a cache-write warning — `docs/foundational/phase3b4/evidence/task3/frozen-build.log:1` in the research worktree.** `go build` exited0 but reported `go: writing stat cache ... operation not permitted`. This is an apparatus/cache warning, not a failed build or evidence of a runtime semantic defect. Preserve and disclose it; use a writable cache for the next required qualification. No unchanged suite rerun was performed merely to remove the warning.

### Exact Checks and Scope

- Read supplied diff once in sequential chunks; no checkout, index, branch or Git-command mutation. A later mechanical diff line-number mapping located the cited guards without re-reviewing changed files.
- Concrete outside-read risk: whether existing general program validation already rejected the dormant payload. Read only `internal/executor/program.go:132–230`; it does not interpret these literal lifecycle claims. Read the public Reciprocal wire/interface definitions to construct the focused public provision. Initial guessed `validate.go`/`reduce.go` paths were absent; no conclusion depended on them.
- Frozen-evidence audit: read `commands.json`, `frozen-source.json`, all named frozen command outputs, `index-and-old-reader.log`, the three original RED logs and the reproduction metadata. Hash-compared all315 listed runtime files; zero mismatches. Recorded ordinary/integration/race results are attributed to retained implementer evidence, not rerun by this reviewer.
- Focused commands used `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3a-go-build go run /tmp/alan-phase3b4-review-index.go create|open <disposable-db>` and the analogous `/tmp/alan-phase3b4-review-v1.go <disposable-db>`. The v1 comparator ran once in the reconstructed base tree and once in the unchanged head tree. No package suites, repeated stress tests or network endpoints were run.
- Every database used by these probes was freshly created under `/tmp`; the base v1 Store was copied with available DB/WAL/SHM components before head Open. No archived SQLite file was opened, and no production/source checkout was changed. No subagents were dispatched.
- No new protocol or contact code appears in the diff, but this review does not independently prove all unchanged dispatch-lock, generic-unknown-recovery or external-effect-free replay behavior. The retained suite and task-scoped source inspection are bounded evidence.

### Assessment

**Task quality: Needs fixes. Verdict: REJECT.**

The exact-selection implementation and bounded tests are coherent, but the two reproduced boundary failures prevent Task3 approval. Correct them narrowly and obtain re-review; retain the unqualified Task4/5 obligations above.
