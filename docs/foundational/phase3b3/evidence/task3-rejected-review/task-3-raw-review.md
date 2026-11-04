# Task 3 independent raw-policy review

**Specification verdict: REJECT.** Three focused probes show accepted evidence that violates the stated raw-policy/recovery contracts. The existing corpus successes and five hostile rejections do not discharge these missing checks.

**Quality verdict: REJECT.** The package has a useful independent structure and preserves the named model mismatch, but the missing negative tests permit false positive validity/custody/recovery claims. Correct and independently re-review the findings below before qualifying this checker. This verdict does not reject the unchanged finite model or the previously reviewed runtime primitive.

Candidate reviewed: the 22-file new `research/delegationcorrespondence` package/fixtures represented by `task-3-review.diff` and `task-3-review-source.json`, over research `31c8af68e818f97a842e790669fbc8f8c98dc8bd`. Raw runtime evidence is source `759df6e3d2505f27e1287b2cfeb07bc4290afc0a`. I read the raw-review brief, Task3 brief/report, candidate source/diff, and actual retained normal/race evidence. Parent-owned manuscript changes are outside this review.

## Blocking findings

### R1 — P2: reconstruct accepted message permission instead of trusting admission/authorization flags

Location: `research/delegationcorrespondence/custody.go:125–135`, especially the `message_origin` predicate at127; source authorization handling in `check.go:514–519` is a separate source-liveness boundary.

`validateAdmissions` verifies that a message admission agrees with the committed message and its observation, then accepts the stored `Accepted` result. It does not require the exact original communication permission `(message, source, target, message, source revision, obligation)`. The later delegated-service check sees an accepted selected request and accepts its authority chain. Original/current permission-table equality is insufficient: both can consistently omit the necessary route.

Focused falsifier: from exact `testdata/positive.json`, remove only the B→A `ba-request` permission from **both** `permissions` and the canonical `composition_permissions` descriptor. Preserve all programs, histories, grant decisions, effect decisions, gates, receipts and selections. The final qualified checker binary returns **exit0, Valid=true, Failure=null, Custody.Complete=true**. No stale hash or malformed descriptor is involved.

Evidence:

- `/tmp/alan-phase3b3-task3-raw-review/missing-request-route.json`
- `missing-request-route-verdict.json`, `missing-request-route-command.json`, and `missing-request-route.stderr` in the same directory.
- Original SHA256 `e371654a8b8b9def3ab6f29edfd067c7bc418f7c9e58a72da23fed2ff8296d18`; mutant SHA256 `304ad2640f7664a6916fe4126d317e4ab583531ccf72991295ba1fb90984b4c5`.

This is a policy contradiction within the input, not a promise to authenticate unsigned evidence. To resolve the named question I additionally inspected owning runtime source `internal/executor/messages.go:466–477,521–535`: ordinary admission calls `admissionPermitted`, queries the exact six-field permission and refuses `route_denied` when absent. No runtime experiment was rerun.

Required correction: independently validate accepted ordinary message routes using retained original permission tuples and source revision, including request, ready and result-bearing reply routes. Apply the analogous declared local-input rule where claimed in scope; keep source authorization distinct from admission permission. Add a negative test that changes both permission mirrors while preserving all successful result rows, and require a semantic admission-policy failure rather than a table-difference failure.

### R2 — P2: captured-prefix continuity omits original authority/source tables

Location: `research/delegationcorrespondence/external.go:329–339`.

`CheckContinuity` compares namespace and an append-only subset of operation/history rows. It does not compare the original/current permissions, program source, deployment, profile, static roots/membership or origin tables. Thus two individually valid snapshots with the same namespace may present different purported immutable original authority and still be certified continuous. This fails the specified no-authority-amplification/same-original-Store recovery check.

Focused falsifier: use `initial/inspect.json` and `final/inspect.json` from actual normal `cut-06b-applied-no-receipt-3686630247`. Add a new B external permission to the final `permissions` table and final canonical `composition_permissions` descriptor only. Existing execution rows/results are unchanged. Both `CheckRaw` calls return valid and `CheckContinuity` returns **nil**.

Evidence/probe:

- `/tmp/alan-phase3b3-task3-raw-review/continuity.go`
- `changed-permissions-initial.json`, `changed-permissions-final.json`, `continuity-result.json` in that directory.
- Result: `InitialValid=true`, `FinalValid=true`, `ContinuityFailure=null`.

Required correction: compare the immutable source/deployment/permission/root/origin sets across every supplied valid prefix, separately from append-only protocol records and allowed run-epoch/consumption changes. Preserve explicit exceptions only where the bounded contract actually permits them. Add the above authority-addition negative and at least a program/deployment drift control; do not discard the contradictory prefix or classify it as unsupported model projection.

### R3 — P2: fresh-owner boundary epochs are not joined to the owner's actual Open epoch

Locations: `research/delegationcorrespondence/external.go:224–278`; `check.go:559–569` intentionally provides only retained historical activation bounds.

The boundary loop checks an event's epoch against the gate descriptor. The separate owner loop checks Open against its post-open snapshot. Nothing joins these two checks. A self-consistent stale gate/event pair therefore passes despite contradicting the supplied current owner. Historical bounds alone cannot reject an epoch between the old reduction and the next new-epoch reduction; that is exactly why the contract requires external Open context for these claims.

Focused falsifier on the same actual applied-unknown cut: change only A4 ordinal2's `AuthorityAdmission.Source.Epoch` from3 to2, recompute its `RecordDigest`, and change the corresponding ordinal2 events in `fresh-owner/boundaries.jsonl` to epoch2. Leave fresh-owner `open.json` and `post-open/inspect.json` at epoch3, leave original Request.Epoch2 unchanged, and preserve the subsequent epoch3 history and endpoint evidence. Both `CheckRaw` and `CheckExternal` return **valid**; the report itself still lists `OwnerEpochs["fresh-owner"]["a"]=3` while accepting returned ordinal2 at epoch2.

Evidence/probe:

- `/tmp/alan-phase3b3-task3-raw-review/epoch.go`
- `stale-current-epoch.json`, `stale-current-epoch-boundaries.jsonl`, `epoch-result.json` in that directory.
- Result: `RawValid=true`, `ExternalValid=true`, both failures null.

Required correction: for process-scoped boundary files, bind the event/gate source RunID and epoch to that owner's actual Open/post-open state and retain the known process/source context. Check the allowed renewal relation between stopped and reopened runs, with terminal runs not renewed. Keep original wire epoch separate. Reject contradictory supplied owner evidence even when local descriptor hashes and event/gate equality agree. This does not require proving authenticity or complete contact-to-gate attribution from unsigned Store bytes.

## Parent questions and additional corrections

**Recipient mutation label/coverage:** Confirmed. `mutations.go:100–101` changes `Authority.Use.Root.Run.Subject`, the original holder, not the recipient. In actual `final/raw-mutations.json`, that case changes holder b/1→a/1 and leaves the entire lifecycle-results table unchanged. Its `first_authority_policy` rejection is useful holder/witness-consistency evidence, but cannot be described as an altered-recipient retained-history mutation. The frozen corpus does contain real recipient-generation/root/revision input controls whose authentic refused decisions are independently accepted as valid histories; that is a different test. Correct the generated name/report and supply a real recipient mutation with exact before/after coordinates and retained result evidence before claiming the Task3 required mutation set is complete. This is a coverage/reporting correction in addition to R1–R3; I did not find a recipient-policy false acceptance in the inspected code.

**Discarded reread errors:** `files.go:149,154` rereads verified bytes with discarded `os.ReadFile` errors. Under the explicitly immutable-input boundary, no input drift was observed and this is not the cause of the three false acceptances. It is still avoidable integrity/error-classification debt: later checks can use bytes different from those named in `report.Files`, and an I/O error can become a semantic failure or missing optional source/boundary evidence rather than the documented apparatus exit2. Cache the already verified JSON/JSONL bytes and use them for projection/external/continuity/source processing, or propagate/reverify reread failures. I did not inject a concurrent filesystem mutation outside the immutable-input assumption, and I do not count this concern as an additional reproduced blocker.

## Existing evidence and corrections verified

I inspected raw source/reader outcomes directly, rather than accepting `Classification` labels or only the implementation report:

- Normal `ordinary-s2-handshake-628745432`: Run reports both roots succeeded, no unresolved/unknown originals; raw `final/inspect.json` exactly matches the committed positive fixture SHA256 above; Inspect/InspectExecution/Replay outcomes are nil. The independent result has two effects and complete custody. Existing normal/race expected outputs retain 48 cases and five invalid final raw histories each.
- Actual frozen missing-grant, forged-grant-origin, effect-use-result, gate-attempt-digest and result-link hostile files retain reader errors separately. Independent first failed edges remain respectively `lifecycle_settlement`, `lifecycle_origin`, `reservation`, `gate_attempt_digest`, and `result_original`. They are not rejected merely for sparse Tables-only shape.
- Actual recipient-generation control is an honest runtime grant refusal, with B failed and A open; raw checking accepts that valid refusal history. Actual applied-unknown-revoked case remains unresolved in `UnknownExternal`, with old application and refused fresh gate; unsupported projection does not erase the raw state.
- The final parent's missing-current-subject input is rejected at `current_subject_set`. The future-receipt-selection input is rejected at `observation_before_effect_result`, order18. The source now requires exact current-subject cardinality and selected observation predecessors at the current fact prefix. Existing retained causal controls separately cover receipt, lifecycle, message and terminal availability; their final outputs reject at the stated specific edges. These earlier findings are resolved in this candidate and are not being reopened by R1–R3.
- The genuine `control-wrong-prior-selection` discrepancy remains `immediate_selection_narrows_model_admission` in both normal/race outputs. Raw refusal validity remains separate from the unchanged model's allowed delegated edge. Separate-grant revocation, unknown external receipt prefixes, repeated exposures and whole-Store crash projection remain explicitly unsupported; none of those classifications excuses a raw-policy failure.

The checker uses independent DTOs, a five-instruction interpreter, raw original permission/root membership and fact-prefix policy reconstruction. It does not invoke runtime authorization or SQLite. Those are sound independence choices; the new probes show specific missing joins rather than a wholesale reliance on runtime verdicts. Detailed generated graph/progress semantics remain the other reviewer's responsibility.

## Executed checks and limits

Only the three named focused falsifiers above were executed during this review. The first used the worker's frozen final checker binary:

```sh
/tmp/alan-phase3b3-task3-evidence/final/check -input /tmp/alan-phase3b3-task3-raw-review/missing-request-route.json
```

The other two called the candidate public checker APIs from temporary standalone probe programs:

```sh
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3b3-task3-build go run /tmp/alan-phase3b3-task3-raw-review/continuity.go
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3b3-task3-build go run /tmp/alan-phase3b3-task3-raw-review/epoch.go
```

Both probe commands exited0 and printed the false-acceptance results preserved above. No broad Go/race/model suite was repeated. Existing final qualification logs/results were inspected and attributed to those prior executions. I opened no SQLite database, mutated no archived evidence/index/source/branch/HEAD, and spawned no subagents. Only this review and temporary probe inputs/programs/results were written.

These are bounded checker findings, not runtime defects, distributed-system claims or a requirement for general hostile-history completeness. Requalification must retain the rejected candidate and exact adverse input bytes, add meaningful negatives, rerun appropriate candidate checks, and obtain independent re-review before this raw-policy task is approved.
