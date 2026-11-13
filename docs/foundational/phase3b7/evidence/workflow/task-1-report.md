# Phase 3B.7 Task 1 implementation report

## Scope

Implemented the test-only producer for the frozen runtime commit
`5b4662fd682a220ca3065b15e86672b7e56567e7`, tree
`44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`. The research branch began at
`62dc6c299007abd40d3729827489c7cfb6e68f49`. The original runtime Git repository was accessed read-only. Environment files
were untouched, including `.envrc`. No matched capture or checker was run.

## Files and interfaces

- `research/matchedcapture/producer/prepare.py LOCAL_RUNTIME_REPOSITORY FRESH_DESTINATION`
  prepares an independent detached local clone and emits the full source
  inventory. It verifies each tracked byte sequence against the frozen Git
  blob, verifies file modes and records SHA-256. It installs four new
  integration-tagged test files from the adjacent `.go.txt` sources.
- `TestPhase3B7MatchedCapture` requires a fresh
  `ALAN_PHASE3B7_EVIDENCE` directory. It writes the pair index, full
  member-length/hash manifest, source/build identity, shared seed, and both
  branches' complete captures, protocol logs and endpoint dossiers.
- `TestPhase3B7Owner` is the subprocess. Scheduler ack/release uses fd3/fd4;
  the independent worker ack/release uses fd5/fd6. The owner context cannot
  time out or inherit cancellation. Any scheduler barrier return syncs abort
  evidence and exits 77 before the Run defer can account completion. Worker
  EOF, invalid commands and canceled context cannot permit a call.
- The controller verifies both acknowledgments, retains the exact release
  writes and wait statuses, and holds the scheduler pipe through confirmed
  SIGKILL. Failed branch cleanup reaps owner before releasing any pipe.
  Sequential branch endpoints use the same literal URLs and fresh files,
  with ready/config/command/token/PID, unfiltered journal, stopped snapshot,
  reaped database and final snapshot.
- `docs/foundational/phase3b7/02-protocol.md` is the complete schema, field
  and order reference for the checker. Producer status reports capture
  completion. The comparison receives direct public `inspect.json` exports
  without removing or projecting fields and checks semantic conclusions
  independently. No historical case wrapper is fabricated.

## Commands and results

Verified from the research checkout:

```sh
python3 research/matchedcapture/producer/prepare.py \
  /Users/jacob/dev/personal/alan /tmp/alan-phase3b7/producer-verified
git -C /tmp/alan-phase3b7/producer-verified diff --exit-code
gofmt -l research/matchedcapture/producer/*.go.txt
```

Preparation exited 0. The clone has 637 tracked runtime members and exactly
four installed overlay files plus the generated source manifest as untracked
changes. The tracked diff and gofmt file list were empty. The emitted
`phase3b7-producer-source.json` SHA-256 is
`557e832fc84e7cbdb47d1d4767938f80ddd1ca1d77c8c647cb665a3436546994`.
The copy in evidence is gzip-compressed with deterministic header fields.

Both GREEN commands used `/tmp/alan-phase3b7/producer-verified`:

```sh
env GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off \
  GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b7/go-build \
  go test -json -count=1 -tags=integration ./internal/durabletest \
  -run '^TestPhase3B7(Barrier|Scheduler)Controls$'
env GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off \
  GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b7/go-build \
  go test -race -json -count=1 -tags=integration ./internal/durabletest \
  -run '^TestPhase3B7(Barrier|Scheduler)Controls$'
env GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off \
  GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b7/go-build \
  STATICCHECK_CACHE=/tmp/alan-phase3b7/staticcheck \
  staticcheck -tags=integration -f json ./internal/durabletest
env GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off \
  GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b7/go-build \
  go test -c -tags=integration \
  -o /tmp/alan-phase3b7/producer-verified-controls.test ./internal/durabletest
go version -m /tmp/alan-phase3b7/producer-verified-controls.test
```

Each command exited 0. Go 1.27.4 on darwin/arm64 used cgo and the pinned
`go-sqlite3` v1.14.52. The nine normal and race controls passed, covering
explicit worker release, EOF, invalid byte, canceled context with release
pending, held worker followed by cancellation, eight acknowledgment identity
mutations within one subtest, scheduler confirmed SIGKILL, unexpected scheduler
release and scheduler EOF. Staticcheck emitted zero diagnostics. The initial
Staticcheck attempt using `XDG_CACHE_HOME` failed because the Darwin cache
resolver still used the home cache; `STATICCHECK_CACHE` corrected the command.

For semantic RED, the same offline variables were used in
`/tmp/alan-phase3b7/producer-dev` with two temporary test overlays, now
preserved as `red-controls.go.txt` and `red-owner.go.txt`. The permissive
implementation returned nil without checking the worker release pipe and
accepted every acknowledgment. Command:

```sh
env GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off \
  GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b7/go-build \
  go test -json -count=1 -tags=integration ./internal/durabletest \
  -run '^TestPhase3B7BarrierControls$'
```

Exit 1 after test execution. Five assertions rejected EOF, a wrong byte, a
canceled context despite a release byte, early escape with no release, and a
mismatched worker Store identity. An initial unused-import build failure was
corrected before the retained RED. The finished code passes those same
assertions. Raw test and build evidence is in
`docs/foundational/phase3b7/evidence/task1/`.

## Decisions and unresolved boundaries

- Both endpoints are SIGKILLed and reaped, matching the frozen fixture's
  shutdown strategy. The producer preserves their complete reaped SQLite file
  sets, including WAL and SHM when present, and inspects disposable copies.
- Normal test failures execute the branch cleanup and pair seal. An external
  process termination or test runtime panic that prevents Go cleanup may
  leave raw files without a complete pair seal; the outer qualification
  driver must preserve and classify such a failed attempt. The checker must
  require the archive seal for a success verdict.
- The applied endpoint may exit its response callback on request cancellation
  after application. The retained record certifies when the controller sent
  its release byte, with the successful receipt gated through confirmed owner
  death. It does not certify continuous endpoint callback occupancy.
- Complete matched capture, independent verification, archive qualification
  and synchronized manuscript effects remain dependent on the next tasks.
