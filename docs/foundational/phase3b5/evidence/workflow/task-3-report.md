# Task 3 self-review

Status: implementation and scoped verification complete, pending controller's independent review. Base `47d4e4a9e546c34e91f41087810f5b99937064a5`. Commit `331186c7b44cdc18b9f7a0970992fe00ce89714b`.

## Question and result

The opt-in projection checks whether the sealed two-principal histories can witness exact-root revocation while retaining historical authorization and debit. Both sealed modes support N05, N06, N07 and P02-v2 and explicitly exclude the genuine applied-unknown history. Each supported cut is an independently audited forward trace. The producer is historical synthetic commit `5b4662fd682a220ca3065b15e86672b7e56567e7`, tree `44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`, with archived Go 1.27.3 and respective race-child flags. This was a new analysis of old bytes. No runtime process or SQLite connection was started.

| Case in each mode | Raw | Physical | Continuity | Legacy projection | New projection | Concrete endpoint outcome |
| --- | --- | --- | --- | --- | --- | --- |
| N05 | Valid | Valid | Pass | Unsupported | Supported, 7/8/9 events | 0 contacts, 0 applications; grant 7, revoke 17, first refusal 21 without debit/gate/result |
| N06 | Valid | Valid | Pass | Unsupported | Supported, 8/9/10 events | 0 contacts, 0 applications; grant 8, first allow/debit 16, revoke 20, refused gate 22, both unpaid |
| N07 | Valid | Valid | Pass | Unsupported | Supported, 9/10/22 events | 2 contacts/applications, old BA permit 17 survives revoke 21; BA result 23 and consume 28; c-backed AB 33/34/35 and consume 40 |
| P02-v2 | Valid | Valid | Pass | Supported | Supported, 21 events | 2 contacts/applications, both consume; BA 16/17/18, AB 28/29/30 |
| applied-unknown-revoked-before-fresh-gate | Valid | Valid | Pass | Unsupported | Excluded by whole_store_loss and applied_without_retained_result; 5 capture locators, zero events | 1 contact/application and 1 acknowledged Store loss, no retained result |

N06 sidecar preserves the B2 lifecycle grant, B1 original, A2 observation selection, A4 attempt, original debit on A's serving composition root, B7 revoke, sole gate at source epoch 1/revision 1, and full typed BA root. N07 carries unresolved BA opportunities at cuts ending 17 and 22, no Apply at either, and two independent audited Apply/consume joins at final cut 48. The complete BA and AB roots both use B's composition membership and differ in the permission tuple. Epoch 1 corresponds to the model's initial activation 0. Model event revision is the pre-event revocation era; the lifecycle sidecar records the resulting era independently of runtime ownership revision.

## Verification and retained evidence

All Go commands ran with `GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b5/go-build STATICCHECK_CACHE=/tmp/alan-phase3b5/staticcheck`, using `ALAN_REVOCATION_CORPUS=/tmp/alan-phase3b5/corpus` unless stated. Evidence below is relative to `docs/foundational/phase3b5/evidence/task3/`.

| Command | Exit | Evidence and result |
| --- | --- | --- |
| `go test ./research/delegationcorrespondence -run '^TestRevocationSealed$' -count=1` before new API | 1 | `red.log.gz`, missing API build RED |
| Same command with API forwarding the legacy report only | 1 | `red-semantic.log.gz`, genuine sealed positive cases rejected by missing opt-in relation and exclusion reasons absent |
| Same command after initial projector | 1 | `first.log.gz`, implementation incorrectly used the grant lifecycle fact's target locator as issuer locator. Compared retained lifecycle/authorization coordinates, corrected separate joins, then reran exit 0 |
| `go test ./research/delegationcorrespondence ./research/delegationcorrespondence/cmd/revocation -run '^TestRevocation' -count=1 -v`, first hostile-data pass | 1 | `attack-first.log.gz`, wrong-target attack had a stale lifecycle settlement mirror, so raw rejection preceded the intended mapping predicate. This was an attack-construction failure and no semantic kill. Rewrote the settlement mirror and reran |
| Same focused command, final | 0 | `green.log.gz`, sealed both modes, scoped and raw attacks, endpoint requirements, ordered continuity, CLI exit/writer checks |
| `go test ./research/delegationcorrespondence -run '^TestBindingDenseAttacks$' -count=1 -v` | 0 | `legacy-attacks.log.gz`, M04 retained, genuine preliminary N02/N03 provide the earlier actual future-selection and late-grant controls |
| `python3 docs/foundational/phase3b5/evidence/task3/prepare_corpus.py /tmp/alan-phase3b5/corpus` | 0 | `archive-verification.json`, both sealed archive SHA-256 values, 232 selected archive members per mode checked, declared extracted bytes checked against both digest layers |
| Same script with `/tmp/alan-phase3b5/task3/reproduction --prepare`, then `ALAN_REVOCATION_CORPUS=/tmp/alan-phase3b5/task3/reproduction go test ./research/delegationcorrespondence -run '^TestRevocationSealed$' -count=1 -v` | 0/0 | `reproduced.log.gz`, independent fresh declared-input extraction and sealed tests |
| `go build -o /tmp/alan-phase3b5/task3/old-check ./research/delegationcorrespondence/cmd/check` before any non-test source edit at base | 0 | Task baseline binary |
| `go build -o /tmp/alan-phase3b5/task3/check ./research/delegationcorrespondence/cmd/check` and `go build -o /tmp/alan-phase3b5/task3/revocation ./research/delegationcorrespondence/cmd/revocation` | 0/0 | Current binaries used for legacy and new deterministic checks |
| `python3 docs/foundational/phase3b5/evidence/task3/compare_reports.py /tmp/alan-phase3b5/corpus /tmp/alan-phase3b5/task3/old-check /tmp/alan-phase3b5/task3/check /tmp/alan-phase3b5/task3/revocation` | 0 | `verdicts.json`, 10 report pairs byte-equal, 10 old/current legacy report pairs byte-equal, embedded legacy objects equal. New CLI exits 0 for the 8 supported cases and 1 for the 2 exclusions. Legacy CLI exits 0 for all 10 |
| `python3 docs/foundational/phase3b5/evidence/task3/qualify_mutants.py /tmp/alan-phase3b5/task3/mutants docs/foundational/phase3b5/evidence/task3` | 0 | `mutants.json` gives exact build/assertion commands and hashes; both overlay builds exit 0, both assertions exit 1 at their intended semantic assertions |
| `staticcheck ./research/delegationcorrespondence ./research/delegationcorrespondence/cmd/revocation` | 0 | `staticcheck.log.gz`, no diagnostics |
| `git diff --check`, `gofmt -l` on the six changed/new Go files, Python `ast.parse` on the three evidence scripts | 0 | No whitespace/format/syntax findings. `source-sha256.json` pins Go sources in authorityservice, mutualobligations, delegationcorrespondence, module files and scripts |

The final hostile matrix reaches separate boundaries. Swapped and aliased complete roots preserve raw, physical and continuity checks and fail at `exact_root_mapping` or `aliased_exact_roots`. A wrong-target c-root revocation is a raw-valid N07 intermediate cut with rehashed program, intents, source/lifecycle request/result, history, settlement and fact mirrors. It fails at `exact_authority_lifecycle_root`, order 21. The compiled collapse-revoke-root mutant falsely accepts that cut. Missing journal/snapshot evidence fails at `required_endpoint_artifact_missing` even where the legacy external report stays valid. Missing snapshot application coverage fails at `endpoint_snapshot_coverage`. Swapped capture declarations fail continuity at `capture_order`. Rehashed N06 attacks fail at `lifecycle_before_source` order 7, `first_authority_policy` order 20 when revocation is moved ahead of first allow, `fresh_gate_policy` order 22, `activation_bounds` order 22 when a positive epoch-1 gate is stale against epoch-2 source history, `result_without_allowed_gate` order 23, and `first_authority_policy` order 16 for future/wrong-original witness coordinates. The compiled omit-current-gate-revocation mutant makes that forged fresh allow raw-valid. It remains outside model support, which is a separate verdict.

## Changes and self-review

- `revocation.go` and the `cmd/revocation` command add the opt-in mapper, separate report and evidence-aware gates. The only legacy change is factoring the private directory load/check to share the same checked bytes. Ten before/after byte comparisons and the retained legacy attacks cover this boundary.
- Focused tests exercise both sealed modes and rehashed attacks. The successful raw c-target attack isolates the new mapper from independent raw rejection. Compiled mutants make the two scope/freshness predicates falsifiable.
- `04-correspondence.md`, `05-runbook.md`, and task3 evidence document and reproduce the scoped result. No production runtime, manuscript or broader model change is in this task.
- Ruling: require both ready records in addition to both journals/snapshot because they bind each label to the tuple's endpoint. The sealed cases meet this requirement. Endpoint data remains unsigned and contextual. Source digest verification supplies byte identity and carries the archived producer claim.
- Ruling: rejected cases retain capture locators with zero projected events. Unsupported whole-Store/unknown histories cannot gain a synthetic atomic Apply. Physical checks precede support, and a raw invalidity is still exposed by the original verdict and top-level RawValid.
- Ruling: the stale-epoch witness uses a rehashed positive old epoch. The wrong-target history is an intermediate cut checked directly by the mapper after raw validation because a complete N07 history would need to change its subsequent AB service decision. This isolates the specific root mapping obligation without inventing a supported full history.
- Residual limits: only the declared two-principal, two-tuple, initial-activation relation, at most one accepted request-bound grant/revocation, single exposure for known results, no whole-Store crash relation, no unknown final outcome, no runtime atomicity, no reverse simulation, no general progress conclusion. Transport placement is causal and the endpoint wire has no gate ordinal. No general permission-algebra homomorphism is established.
- The controller owns the full ordinary/race/lint/build/module qualification and original48 regression. The scoped tests used the two sealed Phase 3B.4 raw modes with the supplied environment. Independent review remains pending.

## Review fix round 1 — serialized unresolved opportunities

Finding I1 at `331186c7b44cdc18b9f7a0970992fe00ce89714b` identified an incorrect archived cardinality. The Go projection retained one BA opportunity, but its byte-like obligation slice serialized to `"AQ=="` and the old summary emitted four BA entries. The corrected opt-in report serializes `Unresolved` as an array of direction labels. Cuts without an unresolved opportunity emit `[]`, including excluded capture locators. The model State/Event structures are untouched. The strict comparator checks the three N07 capture paths in order and requires the exact opportunity lists `["BA"]`, `["BA"]`, and `[]` in each sealed mode. Every other selected capture must contain an empty array.

Historical verdicts, comparison script, source hashes, the complete review and direct old CLI observations are preserved in `docs/foundational/phase3b5/evidence/task3/review-i1/`. Its provenance file binds the copies to their SHA-256 and the rejected implementation's source commit. Direct old CLI commands exited 0 and showed `"AQ=="` at both N07 intermediate cuts in each mode.

All Go commands used the same offline environment and corpus as the initial Task 3 checks. Evidence paths below are relative to `docs/foundational/phase3b5/evidence/task3/`.

| Covering command | Exit | Result |
| --- | --- | --- |
| `go test ./research/delegationcorrespondence/cmd/revocation -run '^TestRevocationCommand$' -count=1 -v` with new assertion before DTO change | 1 | `review-i1/cli-red.log.gz`, unmarshaling the N07 base64 string into a string list fails |
| `python3 docs/foundational/phase3b5/evidence/task3/compare_reports.py /tmp/alan-phase3b5/corpus /tmp/alan-phase3b5/task3/old-check /tmp/alan-phase3b5/task3/check /tmp/alan-phase3b5/task3/i1/revocation-base` with strict comparator before DTO change | 1 | `review-i1/summary-red.log.gz`, comparator rejects the old null field at the first capture; the CLI RED independently reaches the actual N07 base64 defect |
| `go test ./research/delegationcorrespondence ./research/delegationcorrespondence/cmd/revocation -run '^TestRevocation' -count=1 -v` | 0 | `review-i1/green.log.gz`, covering Go projection and serialized CLI assertions |
| `go build -o /tmp/alan-phase3b5/task3/i1/revocation ./research/delegationcorrespondence/cmd/revocation` and `go build -o /tmp/alan-phase3b5/task3/i1/check ./research/delegationcorrespondence/cmd/check` | 0/0 | Fixed CLI and current legacy comparator binaries |
| `python3 docs/foundational/phase3b5/evidence/task3/compare_reports.py /tmp/alan-phase3b5/corpus /tmp/alan-phase3b5/task3/old-check /tmp/alan-phase3b5/task3/i1/check /tmp/alan-phase3b5/task3/i1/revocation` | 0 | Regenerated `verdicts.json`, 10 deterministic pairs, 10 byte-equal legacy pairs, exact N07 cardinalities both modes |
| `python3 /tmp/alan-phase3b5/task3/i1/verify-delta.py` | 0 | `review-i1/delta-verification.json` records commands, exits and report hashes for direct pre/post CLI output comparisons. Removing only the per-capture `Unresolved` field leaves all decoded JSON equal on all ten dossiers. Each corrected CLI list matches the regenerated summary |
| `python3 docs/foundational/phase3b5/evidence/task3/qualify_mutants.py /tmp/alan-phase3b5/task3/i1/mutants docs/foundational/phase3b5/evidence/task3` | 0 | Regenerated `mutants.json`; both builds exit 0 and both intended assertions exit 1 |
| `staticcheck ./research/delegationcorrespondence ./research/delegationcorrespondence/cmd/revocation` | 0 | `review-i1/staticcheck.log.gz`, empty output |
| `gofmt -l` on the three touched Go files, `git diff --cached --check`, Python `ast.parse` on the active evidence scripts, regenerated source SHA-256 and reverified current source, archived copies and mutant logs against their declared hashes | 0 | No format or hash discrepancies |

The original support/exclusion matrix and all audited model events are unchanged. The controller-owned broader qualification was not run during this fix.

Fix commit `4ad1f9ef37c57375d4a39eafc9a5fcd91bc0425d`. Clean worktree after commit.
