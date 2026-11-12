## Status and commit

Task 2 complete and self-reviewed, ready for independent spec/quality review. Branch `codex/phase3b5-grant-revocation`, Task 2 base `e6e53deb`. Commit `0e8d8f874b5993b62692aee1b872ef7619a93422` (`Audit separate-grant revocation and finite completion independently`), tree `a75325bb7a08f4ccc70b98f947652834dfd9c658`.

Only the explicit Task 2 code/test/note/evidence paths were staged and committed. The controller's `docs/foundational/phase3b5/phase3b5-plan.md` edit remains unstaged. This report is under the existing ignored task-plan directory and is intentionally not force-added. No push, PR, manuscript edit, runtime modification, extraction of sealed corpora, production dependency change, or subagent invocation occurred.

## Result

The independent ledger accepts the new opt-in configuration, declares its complete AB binding and B-owned c root independently, checks frozen separate BA issuance including structurally valid rejected claims, preserves historical first decisions and permits across b revocation, checks current roots only for fresh admission, and clears permits on server crash. Existing operational authority/binding/selection/current-validity helpers are not called by event interpretation. Corpus generation and graph enumeration use operational successors to supply candidate events; independence is the interpretation of those events, not a second implementation of the entire operational relation.

`AuditGraph` checks 229,834 shortest prefixes, every one of 1,807,246 arcs, and independent ledger equality at merged targets. Zero unsafe prefixes/edges; paid-original counts `[91520,82304]`. `GrantRevocationOracleCorpus` retains 34 complete history controls: seven lawful controls and 27 hostile modifications, including both authentic post-revocation consumptions, old-permit survival, crash-cleared opportunity, unchanged c authority, sticky first refusal, grant identity/ancestry/tuple/binding/generations/revision, and fresh gate versus historical first allow.

Independent reverse reachability matches `Unsafe==0 && BA.Decision!=Refused && (!Revoked || BA.Applied || BA.Permit)` at every reachable state. Safe complete states: 32,768. Safe-completable: 203,833. Noncompletable: 26,001. Predicate mismatches: 0. AB refusal states: 0. All 57,118 arcs from a BA-refused source preserve refusal. Independent SCC calculation gives 188,614 cyclic states in 75,070 cyclic SCCs. Maximum shortest suffix is 23 steps; no incomplete dead ends, which does not imply inevitable completion.

Shortest immediate revocation counterexample has one event. Early refusal and later retained accepted grant cuts have four/five events. Revoked-after-allow-before-permit cut has six. Revoked permit surviving has seven plus a 15-step successful custody suffix. Revoked-permit crash and application have eight-event prefixes, respectively no completion suffix and a 14-step successful suffix. Each selected prefix/suffix is ledger-audited. Temporal cuts use shortest source prefixes followed by actual qualifying arcs because merged targets' own shortest prefixes can omit the required temporal ordering.

The candidate predicate never failed; no theory/model claim was adjusted to hide an adverse state. The code includes preservation of an initial predicate mismatch witness if one occurs in a future run. Source-level necessity/sufficiency sketches and explicit reachable-state/environmental assumptions are in `docs/foundational/phase3b5/03-authority-progress.md`.

## Verification environment and artifacts

All Go work used:

```sh
GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off
GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b5/go-build
STATICCHECK_CACHE=/tmp/alan-phase3b5/staticcheck
```

Task evidence root: `docs/foundational/phase3b5/evidence/task2/`. `checks.json`, `exports.json`, `mutants.json`, `old-progress.json` and `red-reproduced.json` give argument arrays, explicit working directories, direct process exits and outputs. All test/build/mutant stream artifacts are gzip. `source-sha256.json` pins final Go inputs for the affected model, checker, correspondence and mutual-obligation packages plus module files. It was verified against the working files after qualification. Source-format changes were completed before the final export reruns and final mutant copies. The last test-only adjustment reordered the accepted/rejected issuance assertions so the issuance mutant visibly fails by accepting a forgery; that targeted test passed again afterward.

### RED/GREEN and retained failures

1. Initial `go test ./research/authorityservice -run '^TestRevocationOracle' -count=1 -v` failed at real assertions because the independent ledger rejected the opt-in config. `red.log.gz` preserves the stream. Its first shell wrapper mistakenly assigned zsh's read-only `status`, so the wrapper exited 1 without separately capturing the Go exit. This was corrected, and the same initial positive requirements were rerun in a disposable copy of Task 2 base: direct Go exit 1, no build failure, `unknown-frozen-configuration`; see `red-reproduced.json` and `red-reproduced.log.gz`.
2. Initial `go test ./research/delegationcorrespondence -run '^TestGrantRevocationProgress$' -v -count=1`: direct exit 1, `unsupported progress query S2/SeparateGrantRevocation`; `progress-red.log.gz`.
3. Initial post-implementation progress test found missing revoked-permit crash/application cuts, despite correct predicate/count/SCC results: direct exit 1; `progress-cuts-red.log.gz`. Target BFS parents can differ from the relevant temporal arc. Corrected selection scans actual arcs and audits the resulting trace; no graph or semantic change.
4. `green-initial.log.gz` records initial focused ledger success with the graph intentionally skipped by `-short`. `green.log.gz` records the first successful complete new ledger audit, progress query and command tests. That run skipped the environment-gated historical progress queries; final verification enabled them.
5. The first disposable baseline-progress extraction tried Python's newer `tarfile.extractall(filter=...)` on an older available Python and raised `TypeError`. The compound shell continued into a passing Go test and returned 0, so that shell result is not reported as a successful extraction. The retained `compare_old_progress.py` instead validates and writes regular members individually, and its final run exited 0 with exact-byte equality.

### Final scoped tests and lint

The following direct commands all exited 0. Exact streams and argv are recorded in `checks.json`:

```sh
go test -count=1 -v ./research/authorityservice ./research/authorityservice/cmd/revocation ./research/authorityservice/cmd/audit ./research/authorityservice/cmd/export

ALAN_TASK3_PROGRESS=1 go test -count=1 -v ./research/delegationcorrespondence ./research/delegationcorrespondence/cmd/progress -run 'TestGrantRevocationProgress|TestProgress'

go test -race -count=1 ./research/authorityservice ./research/delegationcorrespondence -run 'TestRevocationOracle|TestGrantRevocationProgress|TestProgressCycles'

staticcheck ./research/authorityservice ./research/authorityservice/cmd/revocation ./research/delegationcorrespondence ./research/delegationcorrespondence/cmd/progress

go test ./research/authorityservice -run '^TestRevocationOracleRejectedIssuance$' -count=1 -v
```

The first four logs are `focused.log.gz`, `progress.log.gz`, `race.log.gz`, `staticcheck.log.gz`; the final focused assertion-order check is `issuance-final.log.gz`. New command tests cover exact opt-in selection, invalid arguments and output-write failure. Old progress tests passed with the environment flag enabled. No unrelated broad suite or final whole-branch suite was run; those remain controller-owned.

### Compiled checker mutants

```sh
python3 docs/foundational/phase3b5/evidence/task2/qualify_mutants.py . /tmp/alan-phase3b5/task2-mutants-final docs/foundational/phase3b5/evidence/task2
```

Wrapper exit 0. Each mutant `go test -c` exit 0; each targeted compiled test execution exit 1, with its expected assertion and no panic. `mutants.json` records exact source replacements and copied-source hashes. Each build/failure log is separate.

- Omit fresh current-root check → forged fresh revoked admission is accepted, targeted assertion kills.
- Revoke c together with b → lawful AB c-root trace is rejected, targeted assertion kills.
- Accept issuance after revocation → forged accepted issuance is accepted, targeted assertion kills.
- Retain crash permit → lawful recovery/fresh-refusal trace is rejected because stale permit remains, targeted assertion kills.

The positive old-permit application and current-path recovery/readmission controls pass in the unmodified ledger, detecting overrestriction. No compilation error or panic was counted as mutation evidence.

### Deterministic exports and archived compatibility

```sh
python3 docs/foundational/phase3b5/evidence/task2/qualify_exports.py . /tmp/alan-phase3b5/task2-exports-final docs/foundational/phase3b5/evidence/task2

python3 docs/foundational/phase3b5/evidence/task2/compare_old_progress.py . /tmp/alan-phase3b5/task2-old-progress-verified /tmp/alan-phase3b5/task2-exports-final/revocation-progress docs/foundational/phase3b5/evidence/task2
```

Both wrapper exits 0. Every enclosed `go build -buildvcs=false` and binary invocation exit 0. Each of the four historical exports ran twice and matched the actual archived bytes, rather than merely current-vs-current outputs:

| Export | Original archive | Decoded/current SHA-256 |
| --- | --- | --- |
| twoagent | `docs/foundational/phase3a/evidence/model.json` | `b380e7653a843b03ad9fe036f8ec5ca03a966cb1151b89e8a3e36a3177baf385` |
| mutualobligations | `docs/foundational/phase3b/evidence/model.json` | `28718bb9d29cd31eb870d05bad63360175abf105d1ee0724f96c36a6e7b43533` |
| authorityservice default | `docs/foundational/phase3b2/evidence/model.json.gz` | `180dc509fcf2a017a241e0c3cfb2a22a7f4f353a0731737c897d053d5feed46e` |
| default independent audit | `docs/foundational/phase3b2/evidence/audit.json.gz` | `a8e2c4f2f5395c3992c12609afdc89948c97347e1c8221942f8574647c0b45c1` |

New audit/corpus output ran twice, equal SHA-256 `a5f9f8fac9faaceed3dc4af2c213565383b7954f75b271e7426807c20e0491ad`. New progress ran twice, equal SHA-256 `faf6ff19137e9c5ac9e22b97a4835699fc151323eb972fb5538c4a277a65db5a`. The preserved artifacts are `revocation-audit.json.gz` and `revocation-progress.json.gz`.

Default progress bytes also match a build from `git archive e6e53deb` at SHA-256 `0dd6ff088c598ec5f82f24bd8d1e2f0418514b1e1f4d406ce80c602db1323d5b`. This is separately labeled a baseline-source comparison because there is no prior archived-export pointer for it. JSON fields added to progress are opt-in/omitted for old configurations.

### Formatting, evidence integrity and commit

`gofmt -l` on all nine changed/new Go files produced no output, exit 0. `git diff --check` and `git diff --cached --check` exited 0. Python gzip decompression of every Task 2 `.gz` and verification of every `source-sha256.json` entry exited 0. The explicit scoped `git add` and commit were run with approved metadata-write escalation and exited 0. The commit path list was inspected afterward; it contains only Task 2 files. `git status --short` shows only the controller's unstaged plan edit.

## Self-review and concerns

I checked the ledger diff against the frozen new contract and verified no invocation of operational authority helpers was introduced. Only the new scenario adds c to ledger roots; historical root/config bindings remain unchanged. Grant structure and issuance-time current validity are separate checks. First decisions and grant outcomes are never revalidated against later revocation. Apply consumes an existing activation-bound opportunity; it does not recheck current b. The ledger's merge equality remains a finite-state consistency check and does not imply exhaustive cyclic-history coverage.

I checked the independent progress queue, safety restriction, SCC passes and shortest witness construction. The progress code calls the operational successor function to enumerate the finite relation, but computes reachability and SCCs independently of `explore.go`; it does not import its answer. Candidate validation covers all states. Temporal cut extraction explicitly handles merged targets and audits selected events, and successful suffix controls require both paid originals. Retry SCCs and the losing immediate-revocation state remain visible.

No known unresolved Task 2 semantic defect. The proof remains a source-level sketch supported by exhaustive finite-state checks, not a mechanized or unrestricted theorem. No runtime trace mapping, distributed behavior, permission-algebra isomorphism, unknown-effect recovery, arbitrary runtime-state completion or inevitable progress claim is made. Concrete correspondence, final full research/race/build/module qualification, manuscript reconciliation and independent reviews remain controller-owned later work. No final PR/delivery claim is made here.
