# Phase 3B.7 Task 1 independent review

**Specification verdict: APPROVE.**

**Task quality verdict: APPROVE.**

Reviewed base `62dc6c299007abd40d3729827489c7cfb6e68f49` through head
`7ea8898b4ad9df31713d4061a691ac26ad78723b`, the supplied review diff,
Task 1 brief/report, and `docs/foundational/phase3b7/01-contract.md`.
Approval covers the test-only producer, its interface, and focused controls.
It does not certify a matched execution, an independent comparison, or the
later qualification/manuscript deliverables. The controller's unstaged Task 2
file-map refinement was excluded.

## Specification compliance

- **Scope and frozen preparation:** the diff adds only the producer overlay,
  preparation script, protocol documentation, and Task 1 evidence. Preparation
  checks the exact runtime commit/tree, creates an independent detached local
  clone, verifies every tracked blob and mode, and installs only named new
  test files (`research/matchedcapture/producer/prepare.py:12-19,33-101`). No
  runtime/dependency or historical Phase 3B.6 protocol change is included.
- **Common seed and fixed identities:** one closed seed is archived and copied
  to both working branches before execution; complete base/WAL/SHM hashes are
  compared and an unexpected rollback journal is rejected
  (`phase3b7_matched_test.go.txt:92-142`). Branches run sequentially, and the
  applied branch must rebind both original URLs exactly
  (`phase3b7_matched_test.go.txt:188-205,354-359`;
  `phase3b7_process_test.go.txt:206-229`).
- **Independent worker and scheduler holds:** the worker syncs its event and
  acknowledgment before closing the notification channel. The scheduler waits
  outside the log mutex, while `external.authorized` is forwarded to the
  inherited observer. Acknowledgments join exact PID/token/Store/scenario,
  selected service/coordinates/admission, and sole in-flight authorization
  (`phase3b7_owner_test.go.txt:63-84,139-181`). The controller releases only the
  applied worker after validating both acknowledgments
  (`phase3b7_matched_test.go.txt:210-220`).
- **Fail-closed lifetime and cleanup:** the owner excludes parent cancellation
  and deadlines; worker EOF, wrong byte, and cancellation return errors.
  Scheduler return writes abort evidence and exits before Run can unwind.
  SIGKILL is recorded before delivery and checked after wait. Release pipes
  stay open until process reaping, including failed branch cleanup
  (`phase3b7_owner_test.go.txt:33-55,87-99,118-127,164-184`;
  `phase3b7_process_test.go.txt:90-103,173-203`).
- **Recovery tail and custody:** stopped capture precedes the public recovery
  open; exact renewal, prescribed sequence/epoch coordinates, controlled
  revocation, fresh refusal, retained original allow, and absence of local
  receipt/settlement are checked. Endpoint response release follows confirmed
  owner death. Endpoints are killed/reaped before final endpoint copies and
  the final local capture (`phase3b7_matched_test.go.txt:221-273`;
  `phase3b7_process_test.go.txt:233-250`).
- **Raw custody and source provenance:** unique output roots, exclusive new
  JSON files, synced protocol records, binary digests/build information,
  failed-working copies, pair status, and a complete regular-member manifest
  are present (`phase3b7_matched_test.go.txt:41-89,175-185,285-359`). The
  protocol explicitly distinguishes producer completion from independent
  semantic verdicts (`docs/foundational/phase3b7/02-protocol.md:49-97`).
- **Controls and task boundary:** the retained semantic RED rejects real
  permissive behavior, including EOF/cancellation/early escape and identity
  mismatch (`docs/foundational/phase3b7/evidence/task1/red.jsonl:10-37`). Both
  normal and race GREEN logs finish with package success and contain all nine
  named controls (`green-normal.jsonl:1-48`; `green-race.jsonl:1-48`, in the
  same evidence directory). Capture/comparison work is explicitly deferred.

Source names above without a directory prefix are under
`research/matchedcapture/producer/`.

## Strengths

The implementation makes the critical causal distinction explicit: worker
permission to call is independent of scheduler permission to account results.
It validates the identities connecting those two holds rather than relying on
timing. Failure cleanup preserves the zero-worker hold until death and keeps
failed working files before temporary-directory removal. The protocol is
specific enough for an independent JSON-only checker and does not introduce
a filtered local-observation format or fabricate old producer compatibility.

The endpoint callback limitation is accurately disclosed: cancellation may
end the callback, but cannot release a successful receipt before owner death
(`02-protocol.md:142-153`). This is consistent with the frozen endpoint's
post-commit barrier error path; continuous callback occupancy is not claimed.

## Findings

- **Critical:** none found.
- **Important:** none found.
- **Minor:** none requiring correction for Task 1.

## Evidence inspected and limits

No experiments, builds, tests, archived-SQLite opens, branch/index changes, or
source modifications were performed by this review. Commands were limited
to reading/searching the supplied files, reading named frozen source with
`git show`/`git grep`, and parsing/hashing retained JSON/log/source evidence.
This report is the sole review output written.

The decompressed preparation manifest has SHA-256
`557e832fc84e7cbdb47d1d4767938f80ddd1ca1d77c8c647cb665a3436546994`, contains
637 runtime members, and identifies the required commit/tree. All four
overlay length/hash entries match the reviewed files. The retained build
record reports Go 1.27.4, darwin/arm64, cgo, integration tags, and
go-sqlite3 v1.14.52. The empty retained Staticcheck output has no diagnostics;
its exit-zero result is attributed to the implementer's report, not an
independent rerun.

Concrete dependency risks were resolved against runtime commit
`5b4662fd682a220ca3065b15e86672b7e56567e7`: `internal/boundary/boundary.go:76-139`
tracks active authorization and treats EOF as error;
`internal/executor/external.go:183-208` places the pre-call observer before
HTTP transport (the HTTP timeout does not release the worker hold);
`internal/durabletest/phase3b4_evidence_test.go:17-51` and
`phase3b4_controlled_test.go:67-72` export through disposable copies;
`phase3b3_runtime_test.go:89-98,411-453` preserves available companions and
uses a SQLite backup for live captures;
`phase3b3_temporal_test.go:55-92` supplies the exact-root revocation tail;
`internal/testcmd/alan-resource/main.go:47-70,149-218` creates and syncs the
unfiltered journal before the post-commit barrier; and
`internal/durabletest/crash_test.go:32-55` propagates race mode to child
binaries. These are source-inspection results, not new runtime qualification.

The next tasks must still freeze the producer/comparison/contract identities,
preserve every qualification attempt, independently validate full exports
without field deletion or array reordering, qualify opposite endpoint facts,
and reconcile the manuscript with the actual outcome. External termination or
panic that prevents Go cleanup can leave an incomplete seal; the outer driver
must preserve and classify that failure as already disclosed in the report.
