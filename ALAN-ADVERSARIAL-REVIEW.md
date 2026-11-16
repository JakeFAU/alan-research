# ALAN: Independent Adversarial Architecture and Code Review

- **Repository:** `github.com/JakeFAU/alan-research`
- **Reviewed Commit:** `aa7a1d2181459f7e40803b4fb9709c59af39b796` (`main`, clean working tree prior to creating this report)
- **Toolchain & Platform:** `go version go1.27.4 darwin/arm64` (`go.mod` specifies `go 1.27.3`; sole module dependency `github.com/mattn/go-sqlite3 v1.14.52`)
- **Review Scope & Constraints:** Production source, existing tests, and documentation were inspected without modification. Independent adversarial counterexamples were executed using disposable Go compiler `-overlay` test packages and removed after capturing output.

---

## Executive Verdict

**Verdict: C. Significant correctness defects**

ALAN (`github.com/JakeFAU/alan-research` at commit `aa7a1d2181459f7e40803b4fb9709c59af39b796`) establishes a credible local deterministic reducer over SQLite WAL storage, clearly separating intent commitment, authorization, transport, receipt recording, and consumption. However, adversarial code inspection and live test reproduction demonstrate **significant correctness, liveness, and soundness defects** across the compiler-to-runtime boundary and hierarchical composition lifecycle that invalidate key guarantees even when existing package suites pass.

Specifically, four defect clusters were reproduced via minimal adversarial tests and two structural divergences confirmed in source:

1. **Total `Store.Run` starvation and observation poison-pill deadlocks:** Any non-`ErrIneligible` `Reduce` error inside `selectReady` (such as an `"end"` output-schema rejection or an oversized input observation on a `"receive"` PC) aborts `Store.Run` immediately before advancing `last` or testing remaining observations/runs, permanently wedging every healthy run in the store.
2. **Unsound `author.Compile` static schema checking:** `checkOperand` ignores `MaxBytes` on variable-to-variable copies, tests existential enum intersection (`src ∩ want ≠ ∅`) instead of subset inclusion (`src ⊆ want`), and CFG branch joins delete conflicting field schemas—allowing statically illegal programs to pass compilation and trigger defect (1) at runtime.
3. **Composition scope-closure split-brain:** When a root run terminates (`Status != "open"`) with an unawaited admitted child, `AuthorizeIntent` still permits child work while `compositionScopeClosed` blocks `prepareAdapterWork` (`scope_closed`), stranding `PrepareRoot` at `Pending: descendant_outcome`. Meanwhile, rejected unpublished child invocations let `SettleRoot` succeed while `RootCompletion.PendingInvocations` reports them pending forever.
4. **Unchecked dynamic `AwaitInvocation` strings and manuscript divergence:** On the reviewed `alan-research` commit `aa7a1d21`, `author.Compile` downgrades `request:invocation:*` handles to raw `"string"`, allowing bogus strings to strand open runs at quiescence. That checkout also lacks `DelegateAuthority`, `RevokeAuthority`, storage formats 13–16, and `RetireSubject` revision increments.

---

## System Baseline and Fundamental Invariant Inventory

### Observed Qualification Runs on the Reviewed Commit

Executed directly on commit `aa7a1d2181459f7e40803b4fb9709c59af39b796`:

| Command | Scope Completed | Result |
| :--- | :--- | :--- |
| `go test -race ./internal/executor` | Full package | `ok github.com/JakeFAU/alan-research/internal/executor 24.798s` |
| `go test -race -tags=integration -run '^(TestComposition\|TestAdapter\|TestCancellation)' ./internal/executor` | Filtered integration tests | `ok github.com/JakeFAU/alan-research/internal/executor 12.472s` |
| `go test -race ./schema/... ./effects/... ./compose/... ./author/... ./runtime/... ./adapters/... ./internal/boundary/... ./internal/fixture/... ./research/...` | **Partial (canceled during `./research/...`)** | Passed: `schema` (`1.267s`), `effects` (`1.181s`), `compose` (`9.818s`), `author` (`21.150s`), `runtime` (`[no test files]`), `adapters/lookup` (`1.267s`), `adapters/openai` (`1.140s`), `adapters/receiptquery` (`1.197s`), `internal/boundary` (`1.152s`), `internal/fixture` (`1.253s`), `research/authorityservice` (`112.996s`), `research/authorityservice/cmd/audit` (`76.364s`), `research/authorityservice/cmd/export` (`24.315s`), `research/authorityservice/cmd/revocation` (`33.215s`). Remaining `./research/...` subpackages incomplete. |
| `go test ./cmd/... ./examples/...` | **Partial (canceled before `./examples/...` finished)** | Passed: `cmd/alan` (`0.376s`), `cmd/alan-compose` (`1.371s`), `cmd/alan-evidence` (`0.180s`), `cmd/alan-investigate` (`1.850s`), `cmd/alan-resource` (`0.189s`). `./examples/...` incomplete. |

### Fundamental Invariant Inventory

1. **Single-host writer serialization and SQLite WAL durability (`internal/executor/store.go`):**
   `acquire` takes an OS advisory lock via `syscall.Flock` on `path + ".lock"` (`store.go:159-180`), and `connect` opens SQLite with `SetMaxOpenConns(1)`, `WAL`, `_synchronous=FULL`, `_foreign_keys=on`, and `_busy_timeout=5000`. Committed transitions append rows to `history` (`store.go:107, 831`).
2. **Explicit storage schema versioning (`internal/executor/store.go:72`, `composition_types.go:13`, `reconciliation_profile.go:13`, `receipt_conflict.go:12`):**
   Base `SchemaVersion = 3`, `compositionSchemaVersion = 4`, `reconciliationSchemaVersion = 5`, plus `receiptConflictSchemaOffset = 3` promoting schemas to versions 6, 7, and 8 via `ensureReceiptConflictSchema`.
3. **Deterministic bounded local reduction (`internal/executor/program.go:290-495`, `collections.go:180-224`, `author_metadata.go:120-281`):**
   `Reduce` advances a program counter, validates output schemas on `"end"` with `Outcome == "succeeded"` (`validateAuthorOutput`), cleans collection cursors (`cleanCursors`), enforces `MaxOutputBytes` (`checkOutputSize`), and checks projection limits (`validateProjection`). Structured collection instructions (`foreach`, `next`, `loop`, `append`) execute inside `reduceCollection`.
4. **Separated intent, authorization, dispatch, transport, receipt, and settlement stages (`internal/executor/messages.go`, `external.go`, `adapters.go`, `dispatch.go`, `receipt_conflict.go`):**
   Internal messaging uses `AdmitInput`, `AuthorizeIntent`, `DeliverMessage`, `AcknowledgeMessage`. External and adapter side effects stage through `prepareExternal` / `prepareExternalWork` / `transportExternal` / `exchangeExternal` and `prepareAdapterWork`, with ownership pinning in `DispatchAuthority`, `pinDispatchOwner`, `admitDispatch`, `registerDispatch`, and conflict tracking via `commitReceiptConflict`.
5. **Typed identity, run sequence, and lifecycle gating (`internal/executor/lifecycle.go`):**
   Runs carry authority/sequence tokens (`run.Authority()`, `run.Sequence`) and subject retirement/run lifecycle transitions execute through `loadLifecycleRequest`, `decideTarget`, `DecideLifecycle`, and `lifecycleOutcome`.
6. **Staged composition admissibility, closure, reservation, and root publication (`internal/executor/invocation.go`, `invocation_outcome.go`, `invocation_closure.go`, `composition_reservation.go`, `composition_root.go`, `composition_prefix.go`, `composition_inspect.go`):**
   Invocations admit via `AdmitInvocation`, freeze scopes via `freezeCompositionScope` / `prepareInvocation` / `prepareCompositionScope`, track closures via `recordCompositionLifecycleClosure` / `recordCompositionRetirement`, reserve resources via `recordCompositionReservation`, and publish root settlement via `SettleRoot` (`composition_root_publications`).
7. **Compile-time static wiring and validation (`author/author.go`, `author/application.go`, `author/invocations.go`, `author/compile.go`):**
   `Application.Provision` matches declared ports, generates permissions, and retains composition templates (`retainInvocations`), while `Compile` lowers steps (`Set`, `IfEqual`, `Send`, `Receive`, `Request`, `Await`, `Success`, `Failure`, `Append`, `ForEach`, `Repeat`) and checks operand/field compatibility via `checkOperand`.
8. **Non-reissuing internal ledger replay (`internal/executor/replay.go`, `composition_replay.go`):**
   `Replay` (`validateLedger`, `validateOrdinaryLedger`, `readHistory`, `readObservations`, `readIntents`) and `composition_replay.go` verify recorded ledger transitions and composition relational facts without reissuing external effects.

---

## Verified Findings

Ranked by severity and blast radius.

### Finding 1: `selectReady` error-abort creates permanent cross-run `Store.Run` starvation and per-run `"receive"` observation head-of-line deadlock

- **Claim violated:** Fair / bounded scheduler liveness across independent runs in a `Store`; single-fault isolation so one malformed run or ineligible/oversized observation cannot permanently freeze unrelated healthy runs or subsequent valid observations.
- **Relevant source:**
  - `internal/executor/runner.go:113, 233-236, 318-326, 336-337, 358-378` (`Store.Run`, `selectReady`)
  - `internal/executor/program.go:470-485` (`Reduce`, `case "end":`)
- **Counterexample:**
  1. **Cross-run starvation:** Two independent open runs exist in one `Store`: `a-fault` (at `PC = 1`, instruction `"end"` with `Outcome = "succeeded"`, holding a string exceeding output `MaxBytes`) and `z-healthy` (at `PC = 1`, valid instruction ready to finish). `Store.Run` calls `selectReady(ctx, last)` with `last == nil`, starting at index `0` (`a-fault`). `selectReady` calls `Reduce(p, run.Projection, nil, id, run.Sequence+1, s.limits)`. At `program.go:470-474`, `validateAuthorOutput` fails with `output schema: field "short": string outside bound or enum`. Instead of skipping `a-fault` or committing a terminal fault transition, `selectReady` executes `return nil, err` (`runner.go:361`), immediately aborting `Store.Run`. On every subsequent `Store.Run` invocation, `last` resets to `nil`, so `selectReady` again inspects `a-fault` first and immediately aborts with `steps=0`—permanently starving `z-healthy`.
  2. **Per-run observation poison pill:** A run `receiver` sits at `"receive"` (`PC = 0`) with two unconsumed observations in the database: `obs1` (`1-oversize-when-bound`, oversized relative to step limits) followed by `obs2` (`2-fits`, valid). `selectReady` iterates `observationIDs` (`runner.go:363-378`). Only `errors.Is(err, ErrIneligible)` continues (`runner.go:371-373`); `Reduce` on `obs1` returns `"resource limit"`, which hits `return nil, err` at `runner.go:374-376`. `obs2` is never tested, and `Store.Run` aborts with `steps=0`.
- **Reproduction:** Executed via disposable Go `-overlay` tests (`-race`, `-tags=integration`) against unmodified repo packages:
  - `go test -v -race -overlay=/tmp/alan_overlay.json -run '^TestAdversarial' ./author` (`MaxBytesHoleAndPermanentStoreRunStarvation`, paired with Finding 2)
  - `go test -v -race -tags=integration -overlay=/tmp/alan_overlay.json -run '^TestAdversarial' ./internal/executor` (`ObservationHeadOfLinePoisonPill`)
- **Observed versus expected:**
  - *Observed (cross-run starvation):* Attempt 1 returns `runErr=output schema: field "short": string outside bound or enum`, `steps=2`, `quiescent=false`, leaving `a-fault=(status:open pc:1 seq:1)` and `z-healthy=(status:open pc:1 seq:1)`. Attempts 2 and 3 return the identical error with `steps=0` and no progress on `z-healthy`.
  - *Observed (HOL poison pill):* `HOL test: runErr=resource limit steps=0 receiver.Status=open receiver.PC=0 obs1=["observation","input","1-oversize-when-bound"] obs2=["observation","input","2-fits"]`.
  - *Expected:* Either the faulty run/observation is durably quarantined, rejected, or transitioned to a failed/ineligible state, or `selectReady` isolates the failing run/observation rather than aborting the scheduler across all runs in the `Store`.
- **Impact:** Critical liveness and cross-run isolation failure; a single bad run or poison observation imposes total Denial of Service on an entire `Store`.
- **Confidence:** Reproduced.

---

### Finding 2: `author.Compile` static schema checking omits variable `MaxBytes` checks, uses existential enum intersection (`src ∩ want ≠ ∅`), and erases field constraints across CFG joins

- **Claim violated:** Compile-time schema/type soundness (`author.Compile` rejecting ill-typed or out-of-bounds data flows statically so compiled programs do not fail `validateAuthorOutput` at runtime).
- **Relevant source:**
  - `author/compile.go:428-466` (`checkOperand`)
  - `author/compile.go:473-481` (`nextField` update on assignment)
  - `author/compile.go:807-816` (CFG branch-join pruning of `consts` and `declFields`)
- **Counterexample:**
  Three distinct static-analysis holes allow invalid data to pass `author.Compile`:
  1. **Unchecked `MaxBytes` across string variables:** In `checkOperand` (`author/compile.go:434-465`), when `o` is a variable reference (`o.Ref != "" && o.Literal == nil && o.Offset == nil`) and `want.Kind == "string"`, `src` is fetched from `mField[o.Ref]`. If both `src` and `want` have `Kind == "string"` and neither specifies `Enum` (`len(candidates) == 0`), `src.MaxBytes` is never compared against `want.MaxBytes`. Copying a larger input string into a smaller constrained `out` field (`short`) compiles cleanly and explodes at runtime `"end"` (`output schema: field "short": string outside bound or enum`), triggering Finding 1.
  2. **Existential enum check (`∃ c ∈ candidates` instead of `∀` / subset):** When `len(candidates) > 0`, `checkOperand` loops over `candidates` and sets `compatible = true; break` on the **first** candidate satisfying both `src.Validate` and `want.Validate` (`author/compile.go:445-456`). Thus:
     - Assigning enum `[allowed, forbidden]` to `[allowed]` passes because `"allowed"` satisfies both;
     - Assigning an unconstrained `string` (`src.Enum == nil`) to `[allowed]` borrows `want.Enum` as `candidates`, tests `"allowed"`, sets `compatible = true`, and passes cleanly!
  3. **CFG join constraint erasure:** At branch joins (`author/compile.go:807-816`), whenever incoming edges disagree on a constant value in `consts[target]` or on `Kind`/`MaxBytes`/`Enum` in `declFields[target]`, the compiler deletes that variable's entry from `consts[target]` and `declFields[target]`. Subsequent `checkOperand` finds neither a constant (`evalConst`) nor a field entry (`mField[o.Ref]`) and silently accepts downstream output/send assignments even when **every** branch assigned a value that violates `out`'s schema.
- **Reproduction:**
  - `go test -v -race -overlay=/tmp/alan_overlay.json -run '^TestAdversarial' ./author` (`MaxBytesHoleAndPermanentStoreRunStarvation`, `ExistentialEnumHole`, `CFGBranchJoinDropsConstraints`)
- **Observed versus expected:**
  - *Observed:* `author.Compile` succeeded without error in all three cases: larger-`MaxBytes` to smaller-`MaxBytes`, `[allowed,forbidden]->[allowed]`, unconstrained `string->[allowed]`, and two-branch joins where every branch assigns an invalid literal to `out`. At runtime, `MaxBytesHoleAndPermanentStoreRunStarvation` aborted `Store.Run` with `output schema: field "short": string outside bound or enum`.
  - *Expected:* `author.Compile` must check universal schema subsumption (`src.MaxBytes <= want.MaxBytes` for unbounded inputs, `src.Enum ⊆ want.Enum` requiring non-empty bounded `src` when `want` is enum-restricted) and join `consts`/`declFields` over least-upper-bound constraints (or conservatively reject missing/untracked bounds) rather than deleting constraints at CFG merges.
- **Impact:** High — destroys compile-time schema soundness and serves as an application-level amplifier for Finding 1 (`Store.Run` permanent wedge).
- **Confidence:** Reproduced.

---

### Finding 3: Root run completion closes the composition scope while admitted unawaited children remain unpublished, causing an authorization–execution deadlock

- **Claim violated:** Composition lifecycle causality and scope closure consistency ("no new mutation authorized after closure"; admitted composition descendants must either remain able to publish their outcome or be transitioned to closure consistently with parent settlement).
- **Relevant source:**
  - `internal/executor/invocation.go:170-185` (`compositionScopeClosed`)
  - `internal/executor/composition_prefix.go:15-17` (`compositionPrefix.closed`)
  - `internal/executor/invocation_outcome.go:379-407` (`prepareCompositionScope`)
- **Counterexample:**
  A root run emits a child invocation and reaches `"end"` (`Status = "succeeded"`) without awaiting that child's completion. Because the root run is no longer `"open"`, both `compositionScopeClosed` (`invocation.go:175-177`, `if r.Projection.Status != "open" { return true, nil }`) and `compositionPrefix.closed` (`composition_prefix.go:15-17`) report the scope as **closed**. Meanwhile, the child invocation was already admitted (`child.Accepted == true`) and has no row in `composition_publications`:
  1. `prepareCompositionScope` (`invocation_outcome.go:389-406`) sees `child.Accepted == true` and `composition_publications` count `0`, returning `Pending: "descendant_outcome"` (`Ready: false`) for `PrepareRoot`.
  2. Calling `AuthorizeIntent` for the child still returns `Allowed = true`.
  3. Calling `prepareAdapterWork` for the child's work rejects with `Allowed = false`, `Reason = "scope_closed"`, `work = <nil>`.
- **Reproduction:**
  - `go test -v -race -tags=integration -overlay=/tmp/alan_overlay.json -run '^TestAdversarial' ./internal/executor` (`RootRunSucceededClosesScopeForUnawaitedChild`)
- **Observed versus expected:**
  - *Observed:* `Unawaited child after root.Status=succeeded: PrepareRoot=(Ready:false Pending:descendant_outcome) child.AuthorizeIntent.Allowed=true child.AdapterAttempt=(Allowed:false Reason:scope_closed) work=<nil>`.
  - *Expected:* Either scope closure (`Status != "open"` on unawaited completion) transitions pending children to closure/rejection consistently across both `AuthorizeIntent` and `prepareCompositionScope` (so the root can settle or fail cleanly), or accepted children created prior to root completion remain runnable to publication.
- **Impact:** High — leaves the composition tree in a permanent deadlock and exhibits split-brain between `AuthorizeIntent` (`Allowed=true`) and `prepareAdapterWork` (`Reason=scope_closed`).
- **Confidence:** Reproduced.

---

### Finding 4: `SettleRoot` succeeds on rejected unpublished children while `RootCompletion.PendingInvocations` reports them pending indefinitely

- **Claim violated:** Consistent persisted truth between root settlement (`SettleRoot`) and composition completion inspection (`RootCompletion` / `readRootCompletions`).
- **Relevant source:**
  - `internal/executor/invocation_outcome.go:379-407` (`prepareCompositionScope`)
  - `internal/executor/composition_root.go:75` (`SettleRoot`)
  - `internal/executor/composition_inspect.go:77-84` (`readRootCompletions`)
- **Counterexample:**
  A frozen composition scope contains child invocation intents whose admissions were recorded as rejected (`!child.Accepted`) with no `composition_publications` rows.
  - `prepareCompositionScope` (`invocation_outcome.go:389-406`) checks `compositionDescendantsPending` and `composition_publications` **only** inside `if child.Accepted { ... }`, skipping rejected children. `SettleRoot` subsequently writes a succeeded publication into `composition_root_publications` (`composition_root.go:75`).
  - `readRootCompletions` (`composition_inspect.go:77-84`) scans `"invocation"` intents and appends `operation` to `PendingInvocations` whenever `SELECT count(*) FROM composition_publications WHERE operation_id=?` equals `0`, ignoring whether the invocation admission was rejected (`Accepted == false`).
- **Reproduction:**
  - `go test -v -race -tags=integration -overlay=/tmp/alan_overlay.json -run '^TestAdversarial' ./internal/executor` (`SettleRootSucceedsWhileRejectedChildIsUnpublished`)
- **Observed versus expected:**
  - *Observed:* `SettleRoot with rejected unpublished children: Ready=true Outcome.Status=succeeded Outcome.Reason=member_succeeded RootCompletion.PendingInvocations=[["root",1,1,0] ["root",1,2,0]]`.
  - *Expected:* Either rejected admissions publish/record a terminal outcome or `readRootCompletions` excludes non-accepted/settled invocations so `RootCompletion.PendingInvocations` is empty once the root settles.
- **Impact:** Medium — causes durable state inconsistency between settled root outcome (`succeeded`) and inspection views (`PendingInvocations` non-empty).
- **Confidence:** Reproduced.

---

### Finding 5: `author.Compile` downgrades `request:invocation:*` handles to raw `"string"`, allowing unvalidated dynamic `AwaitInvocation` calls that strand open runs at quiescence

- **Claim violated:** Request/await handle type safety (`Await` / `AwaitInvocation` matching a valid handle returned by a corresponding request/invocation) and absence of ill-typed compiled programs silently stranding runs open.
- **Relevant source:**
  - `author/compile.go:531-533`
  - `author/compile.go:739-747`
- **Counterexample:**
  In `author/compile.go:531-533`, any string field whose tracked kind starts with `"request:"` is erased back to `"string"`, and `author/compile.go:745` explicitly allows `await` whenever `kind == "string" && strings.HasPrefix(effect, "invocation:")` instead of requiring `kind == "request:"+effect`. Passing an arbitrary constant or unrelated string variable to `author.AwaitInvocation` compiles cleanly and leaves the root run suspended forever in `"receive"` (`Status=open`) while `Store.Run` reports quiescence (`Quiescent=true`).
- **Reproduction:**
  - `go test -v -race -overlay=/tmp/alan_overlay.json -run '^TestAdversarialDynamicAwaitInvocationDeadlock$' ./author`
- **Observed versus expected:**
  - *Observed:* `Dynamic AwaitInvocation with bogus string var: Quiescent=true Steps=5 Dispatches=3 OpenRuns=[{root 1}] rootStatus=open rootPC=2` (`ok github.com/JakeFAU/alan-research/author 1.324s`).
  - *Expected:* `author.Compile` must preserve distinct handle capability/typing (`request:invocation:<id>` or explicit typed invocation-handle union across dynamic `Invoke` sites) rather than accepting arbitrary `"string"` values in `AwaitInvocation`.
- **Impact:** Medium — disables static correlation safety for dynamic invocations and silently produces orphaned open runs at runtime quiescence.
- **Confidence:** Reproduced.

---

### Finding 6: The reviewed research checkout lacks authority lifecycle operations, formats 13–16, and a `RetireSubject` revision increment

- **Claim violated:** Distinct `authority-revision` evolution on authority mutations (`AGENTS.md` §4 / manuscript §6) and documented `DelegateAuthority`, `RevokeAuthority`, and storage formats 13–16 (`paper/alan-dissertation.md` §§6.11–6.15, lines 1027–1102).
- **Relevant source:**
  - `internal/executor/lifecycle.go:69-73, 205-223` (`loadLifecycleRequest`, `decideTarget`)
  - `internal/executor/store.go:72`, `composition_types.go:13`, `reconciliation_profile.go:13`, `receipt_conflict.go:12`
  - `paper/alan-dissertation.md:1027-1102`
- **Counterexample:**
  1. `decideTarget` (`lifecycle.go:205-223`) handles `CancelRun`, `AbandonRun`, and `RetireSubject` by setting `after.Retired = true` while leaving `subjects.revision` unincremented.
  2. `loadLifecycleRequest` (`lifecycle.go:69-73`) accepts only `"Retry"`, `"CancelRun"`, `"AbandonRun"`, and `"RetireSubject"`, rejecting `"DelegateAuthority"` and `"RevokeAuthority"` with `fmt.Errorf("unknown lifecycle operation")`.
  3. Executor schemas implement base `SchemaVersion = 3`, `compositionSchemaVersion = 4`, `reconciliationSchemaVersion = 5`, plus `receiptConflictSchemaOffset = 3` (versions 6, 7, 8). Formats 13–16 described in `paper/alan-dissertation.md` §§6.11–6.15 are not present in `internal/executor`.
- **Reproduction:** Direct inspection of `internal/executor/lifecycle.go:69-73, 205-223`, schema constants in `internal/executor/`, and `paper/alan-dissertation.md:1027-1102`.
- **Observed versus expected:**
  - *Observed:* `RetireSubject` mutates retirement status without advancing `Revision`, while delegation/revocation lifecycle operations fail with `"unknown lifecycle operation"`.
  - *Expected:* Either the code implements `DelegateAuthority`, `RevokeAuthority`, revision bumps on `RetireSubject`, and formats 13–16, or the manuscript (`paper/alan-dissertation.md` and `AGENTS.md` §4) clearly scopes commit `aa7a1d21` to schema versions 3–8 without delegation/revocation support.
- **Impact:** Medium — paper-to-implementation mismatch and stale `Revision` metadata across `RetireSubject` transitions.
- **Confidence:** Source-proven.

---

## Unverified Risks

These credible architectural or operational risks were not reproduced as failing executions during this review and require targeted evidence to resolve:

1. **Pre-crash / post-HTTP external side-effect delivery without endpoint epoch fencing (`internal/executor/external.go`, `adapters.go`, `dispatch.go`, `store.go`):**
   `Store` uses single-host `syscall.Flock` + SQLite `WAL`/`_synchronous=FULL` to avoid concurrent local DB writers, and `external.go` / `dispatch.go` separate intent preparation from transport and receipt commit. However, if a process crashes *after* sending an HTTP request (or while an outbound call is in flight to an external HTTP endpoint that ignores/lacks idempotency keys) and restarts under a new `Store` process/dispatch owner, no protocol mechanism observed here prevents physical duplicate requests at non-idempotent endpoints, nor proves how recovery avoids spurious divergence across every crash window.
   - *Evidence needed:* Crash-injection fault tests killing the executor subprocess (`SIGKILL`) immediately before/after socket write, immediately after HTTP response receipt before SQLite commit, and across DB reopen with an external request counter recording physical HTTP deliveries.
2. **Whole-run `Replay` memory and quadratic cost on maximum-length runs (`internal/executor/replay.go`, `store.go`):**
   `readHistory` reads `history` rows for replay validation, while bounded loops allow long execution traces. Whether maximum-step, large-projection histories (`MaxOutputBytes` / structured collections across iterations) hit OOM or unacceptable validation latency during `Replay` was not benchmarked.
   - *Evidence needed:* Adversarial benchmark / fuzz test executing maximum-permitted steps and collection cardinalities followed by full `Replay` and memory profiling.
3. **Live `adapters/openai` protocol edge cases vs. synthetic fixtures (`adapters/openai/openai.go`, `adapters/lookup`, `adapters/receiptquery`):**
   Existing adapter tests run against deterministic local fixtures (`adapters/openai` passed in `1.140s` offline). Real OpenAI HTTP behaviors (partial chunked responses, rate-limit `429` retry headers, schema drift, proxy truncation, timeouts after request receipt) were not exercised against live endpoints.
   - *Evidence needed:* Fault-injecting reverse-proxy HTTP tests covering chunk drops, mid-stream TCP resets, `429`/`5xx` malformed JSON bodies, and delayed responses racing cancellation.
4. **Incomplete subpackages in this review session (`./research/...`, `./examples/...`):**
   Although `research/authorityservice` and its `audit`, `export`, and `revocation` commands passed under `-race`, remaining subpackages under `./research/...` and all packages under `./examples/...` were canceled mid-run and remain unverified at runtime in this review.
   - *Evidence needed:* Full completion of `go test -race ./...` and `go test -tags=integration ./...` without cancellation.

---

## Claims-Versus-Evidence Matrix

| Guarantee / Claim | Implementation Mechanism | Existing Qualification Evidence | Independent Adversarial Test / Analysis | Verdict |
| :--- | :--- | :--- | :--- | :--- |
| **1. Single-writer DB serialization & WAL durability** | `syscall.Flock` on `path+".lock"`, `SetMaxOpenConns(1)`, `WAL`, `_synchronous=FULL`, `_foreign_keys=on`, `history` table (`store.go`) | `internal/executor` `-race` tests pass (`24.798s`) | Verified lock & SQLite pragma configuration in `store.go`; excludes multi-host / shared network FS assumptions (`README.md`) | **Supported (single-host)** |
| **2. Local deterministic step reduction & output bounds** | Pure `Reduce` validating `"end"` output via `validateAuthorOutput`, `checkOutputSize`, `validateProjection` (`program.go:470-485`) | `internal/executor` unit tests pass | `Reduce` catches invalid outputs cleanly locally, **but** its error return crashes scheduler liveness in `selectReady` | **Supported (locally)** / **Falsified (scheduler impact)** |
| **3. Multi-run scheduler liveness & fault isolation** | `Store.Run` calls `selectReady(ctx, last)` cycling over `ids` (`runner.go:318-378`) | Existing runner tests pass with valid inputs | `MaxBytesHoleAndPermanentStoreRunStarvation`: faulty run at index 0 wedges `z-healthy` forever (`steps=0` on repeats) | **Falsified** |
| **4. Observation dequeueing / resilience against bad observations** | `selectReady` scans unconsumed observations for `"receive"` PC (`runner.go:363-378`) | Existing runner tests pass | `ObservationHeadOfLinePoisonPill`: `obs1` returning `"resource limit"` blocks `obs2` forever | **Falsified** |
| **5. Static schema soundness in `author.Compile`** | `checkOperand`, `nextField`, CFG branch-join pruning (`author/compile.go:428-481, 807-816`) | `author` tests pass (`21.150s`) | Three reproduced holes: `MaxBytes` copy bypass, existential `src ∩ want` enum pass, CFG join field deletion | **Falsified** |
| **6. Type-safe request/await correlation** | `read` / `await` handle check (`author/compile.go:531-533, 739-747`) | `author` & `compose` tests pass | `TestAdversarialDynamicAwaitInvocationDeadlock`: bogus `"string"` accepted by `AwaitInvocation`, leaves `root` open at quiescence | **Falsified** |
| **7. Nested composition parent–child lifecycle & closure** | `compositionScopeClosed` (`invocation.go:170-185`), `compositionPrefix.closed`, `prepareCompositionScope` (`invocation_outcome.go:379-407`) | `TestComposition*` `-race -tags=integration` passes (`12.472s`) | `RootRunSucceededClosesScopeForUnawaitedChild`: root close leaves accepted child authorized yet adapter-blocked (`scope_closed`) and root `Pending: descendant_outcome` | **Falsified** |
| **8. Settlement vs. inspection consistency (`SettleRoot` vs `RootCompletion`)** | `prepareCompositionScope` checks `if child.Accepted`; `SettleRoot` inserts root pub; `readRootCompletions` checks `composition_publications` (`composition_inspect.go:77-84`) | `TestComposition*` passes | `SettleRootSucceedsWhileRejectedChildIsUnpublished`: root settles `"succeeded"` while `PendingInvocations` stays non-empty | **Falsified** |
| **9. Bounded nested collection iteration** | `reduceCollection` (`collections.go:180-224`) and metadata validation (`author_metadata.go`) | `author` / `internal/executor` collection tests pass | Inspected state-transition counters & cursor cleanup (`cleanCursors`); no unbounded escape found in reviewed paths | **Supported** |
| **10. Authority delegation, revocation, & formats 13–16** | `lifecycle.go:69-73, 205-223`, schema versions 3–8 (`store.go`, `composition_types.go`, `reconciliation_profile.go`, `receipt_conflict.go`) | `research/authorityservice` external service tests pass (`112.996s`), not integrated into `internal/executor` | `DelegateAuthority`/`RevokeAuthority` return `"unknown lifecycle operation"`; `RetireSubject` does not bump `Revision`; formats 13–16 absent | **Falsified (paper vs. code)** |
| **11. History-only consistency replay without re-execution** | `Replay` (`replay.go`, `composition_replay.go`) inspects persisted `history`, observations, intents | `internal/executor` replay tests pass | `Replay` validates internal consistency without external transport, **but** shares `Reduce` as its own oracle | **Supported (bounded internal consistency only)** |
| **12. Crash-safe non-idempotent external effect deduplication** | Separate intent commit, `DispatchAuthority`, `external.go`, `commitReceiptConflict` | Synthetic adapter/dispatch integration tests pass | Not tested against live `SIGKILL` mid-socket-write or non-idempotent external HTTP endpoints across restarts | **Unverified** |

---

## Minimal Counterexamples

The standalone Go test files created in `/tmp` via Go compiler `-overlay` were deleted after capture so the working tree remains untouched. Complete execution traces grounded in exact code paths alongside the exact commands and observed outputs are given below.

### Counterexample A: Cross-Run Permanent `Store.Run` Starvation via Compiler `MaxBytes` Bypass

1. **Compile step (`author/compile.go:428-466`):**
   - Program `a-fault` declares an input string field `wide` (large `MaxBytes`) and an output string field `short` (small `MaxBytes`, or strict `Enum`).
   - `a-fault` copies `wide -> short` and returns `Success` (`"end"` with `Outcome = "succeeded"`).
   - In `checkOperand`, `o.Ref == "wide"` with `want.Kind == "string"` and `mField["wide"].Kind == "string"`. Since `len(candidates) == 0`, `checkOperand` returns `nil` without comparing `src.MaxBytes <= want.MaxBytes`. `author.Compile` succeeds.
2. **Runtime step (`internal/executor/runner.go:318-362` & `program.go:470-474`):**
   - Store holds two ready open runs: `a-fault` (at `PC = 1`, `"end"`, carrying an oversized value in `short`) and `z-healthy` (at `PC = 1`, valid `"end"`).
   - **Attempt 1:** `Store.Run` initializes `var last *RunID = nil`, steps `a-fault` to its `"end"` instruction, and calls `selectReady(ctx, last)`. `selectReady` iterates `ids` from `start := 0`, hits `a-fault`, calls `Reduce`, which triggers `validateAuthorOutput(p, s.Data)` (`program.go:472`) -> error `output schema: field "short": string outside bound or enum`. `selectReady` returns `nil, err` at `runner.go:361` before evaluating `z-healthy`.
   - **Attempts 2 & 3:** Each call to `Store.Run` starts again with `last = nil`. `ids[0]` is still `a-fault` (`status=open pc=1 seq=1`). `selectReady` aborts immediately on index `0`, executing `0` steps (`z-healthy` stays `status=open pc=1 seq=1` permanently).

**Executed command & verbatim observed output:**
```text
$ go test -v -race -overlay=/tmp/alan_overlay.json -run '^TestAdversarial' ./author
...
attempt=1 runErr=output schema: field "short": string outside bound or enum steps=2 quiescent=false a-fault=(status:open pc:1 seq:1) z-healthy=(status:open pc:1 seq:1)
attempt=2 runErr=output schema: field "short": string outside bound or enum steps=0 quiescent=false a-fault=(status:open pc:1 seq:1) z-healthy=(status:open pc:1 seq:1)
attempt=3 runErr=output schema: field "short": string outside bound or enum steps=0 quiescent=false a-fault=(status:open pc:1 seq:1) z-healthy=(status:open pc:1 seq:1)
--- PASS: TestAdversarial/MaxBytesHoleAndPermanentStoreRunStarvation (0.03s)
Confirmed: author.Compile accepted both non-subset Enum [allowed,forbidden]->[allowed] and unconstrained string->[allowed]
--- PASS: TestAdversarial/ExistentialEnumHole (0.00s)
Confirmed: author.Compile accepted program where every branch assigns an invalid literal to out
--- PASS: TestAdversarial/CFGBranchJoinDropsConstraints (0.00s)
```

### Counterexample B: Dynamic `AwaitInvocation` String Handle Bypass -> Quiescent Open-Run Leak

1. **Compile step (`author/compile.go:531-533, 739-747`):**
   - `root` constructs or holds an ordinary `string` variable (not obtained from `Invoke`), and executes `AwaitInvocation` targeting an invocation effect (`strings.HasPrefix(effect, "invocation:")`).
   - `author/compile.go:745` permits `await` whenever `kind == "string" && strings.HasPrefix(effect, "invocation:")`.
2. **Runtime step:**
   - `root` advances to `PC = 2` (`"receive"` waiting for an observation correlated to the bogus string handle). No invocation can ever produce that observation.
   - `Store.Run` finishes with `Quiescent = true`, leaving `root` stranded in `Status = "open"` at `PC = 2`.

**Executed command & verbatim observed output:**
```text
$ go test -v -race -overlay=/tmp/alan_overlay.json -run '^TestAdversarialDynamicAwaitInvocationDeadlock$' ./author
Dynamic AwaitInvocation with bogus string var: Quiescent=true Steps=5 Dispatches=3 OpenRuns=[{root 1}] rootStatus=open rootPC=2
ok  	github.com/JakeFAU/alan-research/author	1.324s
```

### Counterexample C: `"receive"` Observation Poison Pill + Root Closure / Rejected-Child Split-Brain

1. **HOL Poison Pill (`internal/executor/runner.go:363-378`):**
   - `receiver` waits at `PC = 0` (`"receive"`) with `obs1` (`1-oversize-when-bound`) and `obs2` (`2-fits`) in the DB.
   - `selectReady` calls `Reduce` on `obs1`; `Reduce` returns `"resource limit"` (not `ErrIneligible`).
   - `runner.go:374-376` immediately returns `nil, err`, never reaching `obs2`.
2. **Unawaited Child after Root `"succeeded"` (`invocation.go:170-185`, `invocation_outcome.go:389-406`):**
   - Root emits child invocation and runs to `Status = "succeeded"` without awaiting child.
   - `compositionScopeClosed` returns `true` (`r.Projection.Status != "open"`).
   - `PrepareRoot` hits `if child.Accepted` -> checks `composition_publications` (`0` rows) -> returns `Ready: false, Pending: "descendant_outcome"`.
   - `AuthorizeIntent` on child succeeds (`Allowed: true`), while `prepareAdapterWork` returns `Allowed: false, Reason: "scope_closed", work: <nil>`.
3. **Rejected Unpublished Children (`invocation_outcome.go:389-406`, `composition_root.go:75`, `composition_inspect.go:77-84`):**
   - Rejected child admissions (`!child.Accepted`) have `0` rows in `composition_publications`.
   - `prepareCompositionScope` skips non-accepted children, allowing `SettleRoot` to record `Outcome.Status = "succeeded", Outcome.Reason = "member_succeeded"` (`Ready = true`), while `readRootCompletions` (`composition_inspect.go:77-84`) leaves `[["root",1,1,0], ["root",1,2,0]]` permanently in `PendingInvocations`.

**Executed command & verbatim observed output:**
```text
$ go test -v -race -tags=integration -overlay=/tmp/alan_overlay.json -run '^TestAdversarial' ./internal/executor
HOL test: runErr=resource limit steps=0 receiver.Status=open receiver.PC=0 obs1=["observation","input","1-oversize-when-bound"] obs2=["observation","input","2-fits"]
--- PASS: TestAdversarial/ObservationHeadOfLinePoisonPill (0.02s)
Unawaited child after root.Status=succeeded: PrepareRoot=(Ready:false Pending:descendant_outcome) child.AuthorizeIntent.Allowed=true child.AdapterAttempt=(Allowed:false Reason:scope_closed) work=<nil>
--- PASS: TestAdversarial/RootRunSucceededClosesScopeForUnawaitedChild (0.08s)
SettleRoot with rejected unpublished children: Ready=true Outcome.Status=succeeded Outcome.Reason=member_succeeded RootCompletion.PendingInvocations=[["root",1,1,0] ["root",1,2,0]]
--- PASS: TestAdversarial/SettleRootSucceedsWhileRejectedChildIsUnpublished (0.04s)
ok  	github.com/JakeFAU/alan-research/internal/executor	1.377s
```

---

## Architectural Assessment

### What Is Structurally Sound
The underlying computational design in `internal/executor` separates the right state boundaries:
- A pure finite-step reducer (`Reduce`) isolated from SQLite I/O and external HTTP calls.
- Durable SQLite commits (`WAL`, `_synchronous=FULL`, `syscall.Flock`) appending explicit step transitions to `history` before authorizing external dispatch.
- Explicit composition tables distinguishing admission, scope freezing, closure, reservation, publication, and root settlement instead of conflating them into an implicit in-memory call stack.
- Replay checking persisted state transitions without re-invoking external endpoints.

### Implementation Bugs vs. Deeper Design Issues

| Defect Class | Nature of Issue | Required Remediation Scope |
| :--- | :--- | :--- |
| **`selectReady` aborts `Store.Run` on non-`ErrIneligible` `Reduce` error (Finding 1)** | **Execution-model / scheduler contract defect** | Cannot be solved just by fixing `author.Compile`: any external input observation or projection limit breach at runtime can trigger `Reduce` failure. The execution model must define whether a faulty run transitions durably to failed/quarantined (recording fault status) and whether a malformed/oversized observation is dead-lettered/marked rejected rather than aborting unrelated runs. |
| **`author.Compile` existential enum, missing `MaxBytes` comparison, CFG join constraint deletion (Finding 2)** | **Local compiler implementation defect** | Fixable within `author/compile.go` using sound abstract interpretation: enforce `src.MaxBytes <= want.MaxBytes`, require subset inclusion (`src.Enum ⊆ want.Enum`, disallowed from unbounded strings), and compute lattice joins instead of deleting `consts`/`declFields` keys on branch merge (requiring explicit coercion/validation or rejecting unprovably bounded values). |
| **Root `"succeeded"` closes composition scope while awaiting unawaited accepted children (Finding 3)** | **Composition state-machine design defect** | Requires aligning the composition semantics: either an application/root run cannot transition to `"succeeded"` cleanly while admitted accepted children remain unresolved, or closing a scope on parent completion must deterministically cancel/close/publish those unawaited children so `PrepareRoot`, `AuthorizeIntent`, and `prepareAdapterWork` never contradict one another. |
| **`SettleRoot` vs. `RootCompletion.PendingInvocations` mismatch on rejected children (Finding 4)** | **Local contract / inspection consistency defect** | Fixable cleanly by consulting admission acceptance / scope settlement status in `readRootCompletions` (`composition_inspect.go:77-84`) or recording terminal publication entries for rejected child admissions. |
| **Dynamic `AwaitInvocation` accepts arbitrary `"string"` handles (Finding 5)** | **Authoring language / type-system design gap** | Fixing this changes `author`'s user-facing handle typing: either dynamic invocation handles retain nominal capability typing through variables/collections without collapsing `request:invocation:*` to `"string"`, or `author.Compile` restricts `AwaitInvocation` to provably invocation-originated variables. |
| **`RetireSubject` revision omission & manuscript §§6.11–6.15 divergence (Finding 6)** | **Scope / specification discrepancy** | Increment `Revision` on `RetireSubject` in `decideTarget` if revision monotonicity is intended, and clearly delineate in the manuscript which lifecycle operations (`DelegateAuthority`, `RevokeAuthority`) and formats (`13–16`) live outside `aa7a1d21` (`Schemas 3–8`). |

---

## Qualification Gaps

Ranked from highest to lowest expected defect-discovery value:

1. **Multi-run mixed healthy/faulty scheduler starvation & poison-observation tests (High):**
   Existing `internal/executor` tests exercise valid runs or isolated single-run failures. None verified that a failing `"end"` output or oversized `"receive"` observation leaves *other* independent open runs or valid queued observations unaffected.
2. **Compiler soundness property/fuzz testing (`author.Compile` vs. runtime `Reduce` oracle) (High):**
   No property test generates two-branch CFGs, narrow `MaxBytes`, narrow `Enum` ports, or non-handle strings paired with arbitrary inputs to check the invariant: *"if `author.Compile` succeeds under declared limits, valid inputs will not cause an `output schema` rejection or unawaited capability deadlock at runtime."*
3. **Full composition edge-matrix tests across all parent–child termination permutations (High):**
   Missing negative tests pairing: (a) parent `"succeeded"` / `"failed"` / canceled before child adapter execution with unawaited accepted children; (b) `PrepareRoot` + `SettleRoot` + `RootCompletion` on rejected child admissions; and (c) `AuthorizeIntent` vs. `prepareAdapterWork` consistency after scope closure.
4. **Independent replay oracle outside `Reduce` (Medium):**
   `Replay` (`internal/executor/replay.go`) invokes the same `Reduce` implementation used during live execution (`runner.go`/`program.go`). A symmetric bug in `Reduce` passes both live execution and `Replay` without detection.
5. **Subprocess `SIGKILL` crash matrix across every commit/dispatch boundary (Medium):**
   Deterministic unit/integration tests verify simulated states, rather than killing an actual worker process across SQLite write, external HTTP write, HTTP response parse, and SQLite receipt commit under real concurrency.
6. **Live / fault-injecting external HTTP tests for `adapters/openai` (Low–Medium):**
   Current tests validate offline fixture behavior (`adapters/openai` `1.140s`); incomplete HTTP streams, `429` backoffs, and timeouts racing run cancellation remain unexercised against realistic transport faults.
7. **Complete baseline test execution for remaining `./research/...` and `./examples/...` packages (Low–Medium):**
   These subpackages were canceled before completion in this review and need an uninterrupted CI/local pass (`go test -race ./...`).

---

## Final Recommendation

**Do not rely on commit `aa7a1d2181459f7e40803b4fb9709c59af39b796` for building multi-tenant, reliability-critical, or higher-level agent applications in its current state.**

The framework is suitable right now **only** as an experimental research prototype on a single host when all authoring inputs and compositions are trusted, well-typed by manual discipline, and tested end-to-end against composition closure paths.

**Guarantees callers MUST NOT rely upon at the reviewed `alan-research` commit `aa7a1d21`:**
- Do **not** rely on one failing run or malformed observation not halting every other unrelated run sharing the same SQLite `Store`.
- Do **not** rely on `author.Compile` passing as proof that `MaxBytes`, `Enum` constraints, CFG-merged variables, or dynamic `AwaitInvocation` handles are type-safe at runtime.
- Do **not** rely on composition scope closure allowing unawaited child invocations to settle cleanly, nor on `RootCompletion.PendingInvocations` matching `SettleRoot` when child invocations are rejected.
- Do **not** rely on that checkout incrementing `subjects.revision` on `RetireSubject`, supporting `DelegateAuthority`/`RevokeAuthority` or schema formats 13–16, or guaranteeing exactly-once external HTTP delivery across arbitrary process crashes against non-idempotent endpoints.

The separately maintained `JakeFAU/alan` runtime at [`c8a39047`](https://github.com/JakeFAU/alan/commit/c8a39047da2cd4214149b7a057f240b1833439e1) documents formats 13/14 for request authority v1 and 15/16 for exact request binding v2. See its [SQLite format contract](https://github.com/JakeFAU/alan/blob/c8a39047da2cd4214149b7a057f240b1833439e1/README.md#L87) and [format constants in the executor](https://github.com/JakeFAU/alan/blob/c8a39047da2cd4214149b7a057f240b1833439e1/internal/executor/request_authority.go#L15-L20).

---

## Final Adversarial Challenge

> **Hypothesis:** *"The system passes every existing qualification test, but an adversarial workload can still violate a foundational execution guarantee."*

**Confirmed.** Every unit and integration test executed on commit `aa7a1d21` in `schema`, `effects`, `compose`, `author`, `adapters/*`, and `internal/executor` passed under `-race`. Yet a minimal two-component `author.Application` exploiting ordinary documented `author` primitives (`Receive`, `Set`, `Success`)—without touching private internals, corrupting SQLite, or bypassing any API check—compiles without error via `author.Compile`, violates output `MaxBytes` (Finding 2) inside `Reduce` on `"end"`, triggers `selectReady`'s unisolated error return (Finding 1), and permanently halts `Store.Run` for an unrelated healthy run in the same `Store` (`steps=0` on every subsequent execution attempt). Simultaneously, on the composition path (`internal/executor`), emitting a child invocation and reaching root `Status = "succeeded"` without `AwaitInvocation` places the composition into a persistent split-brain where `AuthorizeIntent` succeeds, `prepareAdapterWork` denies with `scope_closed`, and `PrepareRoot` remains stuck on `Pending: "descendant_outcome"` (Finding 3).
