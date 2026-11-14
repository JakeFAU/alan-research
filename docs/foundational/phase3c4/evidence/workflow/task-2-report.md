# Task 2 implementation report

Status: DONE; independent task review remains controller-owned.

## Implemented

Created `docs/foundational/phase3c4/tools/audit.py`, `test_audit.py`, and
`capture.py`; the complete narrative is `docs/foundational/phase3c4/03-independent-audit.md`.
Exact preserved inputs, command records, source manifests and prior failed
preparation attempts are in `docs/foundational/phase3c4/evidence/task2/`.

The checker independently reconstructs finite language/profile certificates,
all 2,772 complete product graphs, reverse reachability, all 2,376 substitution
rows, negative paths, private lasso/escapes and generic strictness control.
Typed complete structural equality validates every key/node/arc/coverage/acceptance
bit and shortest witness; no exported verdict/count serves as an oracle.
Only unchanged source-bound Phase 3C.3 descriptor constructors and Phase 3C.2
configuration/action/JSON utilities plus predecessor provenance helpers are reused.
Their exact Python hashes are verified before import. No Go operational source
was read or translated. All 27 predecessor provenance hashes and pinned raw/gzip
identities are verified; inherited descriptor relations are independently checked
before preserving authenticated predecessor serialization.

## Executed evidence

`PYTHONDONTWRITEBYTECODE=1 python3 docs/foundational/phase3c4/tools/capture.py qualification-1`
completed with exit 0, creating a never-overwritten evidence directory.
Its records include these actual subprocess commands:

- `python3 -m unittest discover -s docs/foundational/phase3c4/tools -p 'test_*.py' -v`: 9/9 test methods passed; 39 semantic hostile subtests and seven malformed-JSON controls included. Final run 15.691 seconds.
- `go run ./research/residualinterfaces/cmd/export`: exit 0, empty stderr, fresh 8,271,854 raw bytes exactly equal Task 1; raw SHA-256 `b40f1d5ac0aea0f43999502ea318cb128f545eb8fb1fd129451915a5aa84bf1f`; source manifests exactly equal.
- `python3 docs/foundational/phase3c4/tools/audit.py PATH`: 50 standalone CLI executions, two valid inputs (plain/gzip) accepted, 48 exact hostile inputs rejected with the specified diagnostics. All exact hostile bytes, including invalid UTF-8 and malformed JSON, remain archived. No initially accepted hostile input occurred.

Clean Task 1 saved compiler enumeration was independently parsed and its nine
compiled local Go files matched the current producer manifest. This was not a
new enumeration run. No Go test/race/Staticcheck suite, whole-root suite, old
campaign, SQLite open or physical history capture was repeated.

Measured counts: 72 plants, 7 monitors, 324 controllers, 2,772 products;
29,156 states, 35,672 arcs, 14,410 CanComplete states; 886 initially completable
and 656 nonblocking products. 21 comparisons, 10 equal initial languages,
6 equal profiles, 223 residual comparisons, 850 language subset pairs,
89 profile subset pairs. 2,376 substitution rows: 792 nonblocking, 1,584 blocking.

Subsequent explicit archive/source inspection verified every archive manifest
hash, every checker source hash, all 50 input raw/compressed hashes and command
exits/diagnostics. The removed non-tree edge is
`["U,W,false,false","0"] --commit:AB--> ["W,W,false,false","0"]`;
actual alphabet-ordered discovery parents remain present and counts are repaired.
Every witness edge and endpoint was checked against its full graph; negative
endpoints have unpaid AB debt and opposite CanComplete flags.
`git diff --cached --check` passed.

## Preserved failures and TDD status

The task did not require a TDD-first sequence; do not describe the initial
preparation failures as semantic hypothesis REDs. The first positive CLI attempt
rejected because inherited Python constructor BFS sibling order differed from the
pinned inherited ordering. Its exact checker and diagnostic are retained as
`attempt-1-audit.py` / `attempt-1.json`; independent relation normalization fixes
this without relaxing new product/certificate ordering. The first test run had
one fixture-construction error, `KeyError: AB`; the old test source and diagnostic
summary remain in `attempt-2-test_audit.py` / `attempt-2-test-failure.txt`.
Correcting the fixture to `Signature.Base.Originals[0].CallerGeneration` made
that attack execute and reject. No hostile evidence was accepted before repair.

## Self-review and scope

Inspected the staged status, full code/report diff, preserved source deltas,
all CLI diagnostics and complete archive content/hash correspondence. Only Task 2
paths are staged. The controller's concurrently created `06-reviews.md` is
untracked and excluded. No prior evidence, Go source, runtime code, paper,
ledger or root qualification metadata was changed.

The audit remains a finite independent reconstruction with explicit constructor
reuse. It does not mechanize the conditional source argument or establish
unrestricted ALAN congruence, minimality/necessity, full-agent substitution,
distributed ownership, split effects, fairness or compiler realization of the
abstract private loop. Manuscript synchronization is Task 3; review and delivery
are controller responsibilities. No unresolved implementation concern was found.

Commit created: `92c5d55a` — Audit residual profiles and nonblocking evidence independently.
Post-commit status contains only controller-owned untracked `docs/foundational/phase3c4/06-reviews.md`; all Task 2 paths are committed.
