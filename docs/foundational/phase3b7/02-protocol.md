# Phase 3B.7 producer protocol

The raw archive format is `alan/phase3b7-matched/v1`. A single invocation creates
one fresh output directory, processes zero before applied, and uses one closed
runtime seed for both. Each branch has its own endpoints, owner and working
database paths. The controller keeps all archived SQLite files unopened.
Endpoint and runtime exports use their working files or disposable copies.

## Preparation and build identity

From the research checkout, run
`python3 research/matchedcapture/producer/prepare.py <local-runtime-repository> <fresh-checkout>`.
The command verifies the frozen runtime commit and tree, checks out the commit
detached in an independent local clone, verifies every tracked file's mode and
Git blob identity, then inserts the four named integration test files. The
original repository is read only. `phase3b7-producer-source.json` in the
prepared checkout is an inventory of every tracked runtime file, every
installed overlay and the preparation script. It uses these fields:

| Field | Contents |
|---|---|
| `version` | Archive protocol version |
| `runtime_commit`, `runtime_tree` | Full Git object IDs |
| `runtime_files` | Map of repository-relative paths to `{mode, git_blob, bytes, sha256}` |
| `overlays` | Map of installed repository-relative paths to `{source, bytes, sha256}`. The source is the corresponding producer filename ending in `.go.txt` |
| `preparation_script` | `{bytes, sha256}` of `prepare.py` |

Run the Go test in the prepared checkout. Set `ALAN_PHASE3B7_EVIDENCE` to a
fresh path whose parent exists. The ordinary and race invocations use the same
prepared sources and distinct raw roots.

```sh
env GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off \
  GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b7/go-build \
  ALAN_PHASE3B7_EVIDENCE=/absolute/fresh-raw-normal \
  go test -json -count=1 -tags=integration -timeout=10m ./internal/durabletest \
  -run '^TestPhase3B7MatchedCapture$'
```

For the race invocation add `-race`. The existing endpoint binary builder
propagates race mode. Keep the command, exit status, full test output, and
frozen comparison and contract identities alongside the sealed raw directory.
The raw `source.json` adds `version`, `prepared` (the preparation manifest),
`go_version`, `goos`, `goarch`, `race`, `environment` (a fixed allowlist),
`test_argv`, `build_info`, `owner_binary`, and `endpoint_binary`. The binary
fields are length and SHA-256 records. The owner is a subprocess of the same
test executable. The historical endpoint uses the default HTTP fixture.

## Archive and seed

- `pair.json` declares `version`, `status` (`complete` or `failed`), `mode`
  (`normal` or `race`), `controller_pid`, and `completed_branches` (an ordered
  array containing zero and then applied on successful capture). This is
  producer status. The independent comparison determines equality, custody
  and world facts.
- `members.json` declares `version` and `files`, a map of forward-slash
  relative member paths to `{bytes, sha256}`. It covers every regular file
  including `pair.json` and excludes itself.
- `seed/provision.json` is the public provision. `seed/manifest.json` is the
  reciprocal `delegation.Manifest` with the v2 claim and exact revocation
  fixture. `seed/initial/` is the sealed capture of the closed seed.
  `seed/copies.json` contains `seed`, `zero`, and `applied` maps of
  `runtime.sqlite` and any present `-wal` and `-shm` companions to length
  and SHA-256 records, written before either owner starts.
- Each of `zero/` and `applied/` contains `protocol.jsonl`,
  `first-owner/open.json`, `first-owner/post-open/`, `stopped/`, `reopened/`,
  `recovery-open.json`, `controlled-revocation.json`, `run.json`, `final/`,
  and `endpoint-a/` and `endpoint-b/` dossiers. Each SQLite runtime capture
  has the frozen `bindingCapture` layout, including the complete
  `inspect.json`, typed inspection, replay, outcomes, history, raw SQLite
  file set and request-binding sidecars. The inspection JSON is the direct
  public runtime export. The source's defined schema controls whether
  optional top-level fields are present.
- Both owner open files declare `PID`, `Scenario`, `Token`, `Store`,
  `Namespace`, `Runs` (map by subject), and `Error`. `run.json` has
  `Report` and `Error` from the recovery Run call.
- `first-owner/boundaries.jsonl` retains every original
  `boundary.Event`, enriched with `Scenario`, `Token`, `Store`, and
  `PID`. `first-owner/pre-call.json` is the exact worker acknowledgment.
  Applied also has `first-owner/worker-release.json` containing the same
  event after receipt of the release byte and before transport. The owner
  keeps `process.log`; its `abort.json` records `stage`, `event` and
  `error` on barrier failure.
- Each endpoint has `process.log`, `ready.json` (`URL`, `Token`, `Store`,
  `PID`), unfiltered `journal.jsonl` (possibly empty),
  `stopped.json` and `final.json` (each `{ready, snapshot, error}`), and
  the reaped `resource.sqlite` file set. The endpoint journal uses the
  frozen `alan-resource` records, with `PID` and `Process` token. Its
  selected `ingress` and `committed` records retain `Contact` plus the
  exact wire request and authentic receipt respectively. The stopped
  snapshot is read after owner death and before any endpoint-response release.
  The final snapshot is read from a disposable complete copy after endpoint
  process death and reaping. There are no additional endpoint instances
  during the interval.
- Failed attempts keep any completed artifacts and cleanup records.
  `failed-working/` contains the reaped working file sets available in a
  failed branch. `failed-seed.sqlite` and companions may also be present.

## Controller transcript

Every line of a branch transcript is an object with `ordinal` (consecutive,
starting at 1), `action`, `role` (`controller`, `owner`, `endpoint-a`,
or `endpoint-b`), `pid`, and `data`. Each line is synced before the
controller proceeds. Success errors are the string `"<nil>"` unless a signal
was expected. The following table describes the full sequence for a complete
branch. Endpoint starts and snapshot/stop pairs occur in A then B order.

| Action | Role | Data |
|---|---|---|
| `start-request` | endpoint | `{config, command}`. PID is 0 until start result. `config` is a boundary event carrying scenario, token, Store, and any barrier selector |
| `start-result` | endpoint | `{error}` and actual PID |
| `ready` | endpoint | `{record, error}`. Record is the ready object |
| `seed-copies` (zero only) | controller | `{file: "seed/copies.json"}` |
| `start-request` | owner | `{config, command}`, PID 0 |
| `start-result` | owner | `{error}` and actual PID |
| `pre-call-ack` | owner | `{event, error}`. Event is the original `dispatch.admitted/pre-call` of the selected BA service |
| `scheduler-ack` | owner | `{event, error}`. Event is inherited `runner.waiting/in-flight` with the sole selected authorization ID |
| `release-request`, `release-result` (applied only) | owner | `{target: "worker", command: "c"}`, then `{target: "worker", bytes: 1, error}` |
| `application-ack` (applied only) | endpoint-b | `{event, error}`. Inherited `resource.applied/post-commit` |
| `kill-request` | owner | `{purpose: "owner-loss", signal}` |
| `wait` | owner | `{purpose: "owner-loss", kill_error, error, signal, exit_code, status}`. SIGKILL has signal `"killed"` and exit code -1 |
| `capture` | controller | `{cut: "stopped"}` |
| `endpoint-snapshot` | both endpoints | `{cut: "stopped", error}` |
| `release-request`, `release-result` (applied only) | endpoint-b | Same release objects with target `"endpoint-response"` |
| `recovery-open` | controller | `{identity, error}`. Identity is a process-bound event with a fresh recovery token and the same branch Store |
| `capture` | controller | `{cut: "reopened"}` |
| `revocation` | controller | `{file: "controlled-revocation.json"}` |
| `recovery-close` | controller | `{run_error, close_error}` |
| `kill-request`, `wait`, `endpoint-snapshot` | each endpoint | Kill/wait purpose `"endpoint-stop"`, snapshot cut `"final"` |
| `capture` | controller | `{cut: "final"}` |
| `branch-complete` | controller | `null` |

Process cleanup can append kill and wait records with purpose `"cleanup"` on
failure. The controller requires exact PID, scenario, Store and token on
acknowledgments. The worker acknowledgment names the selected authorization,
BA service `["a",1,4,0]`, epoch 2, sequence 4 and authority admission ordinal
1. The scheduler acknowledgment carries the same authorization as its only
`InFlight` value. The endpoint application acknowledgment carries the exact
request key in `Intent` and authorization in `Attempt`. Recovered state and
world evidence provide the independent checks of these joins.

The owner forwards `external.authorized` to the inherited observer. The
scheduler reaches the inherited barrier only after the independent worker
record is synced and acknowledged. The controller never sends the scheduler
release. A scheduler observer return writes and syncs abort evidence and exits
the owner with code 77 before deferred Run accounting can execute. The owner
uses a context without a deadline or parent cancellation. The worker
requires an explicit `c` and a live context, with EOF and cancellation
classified as errors. Cleanup keeps the release pipes open until the owner is
reaped. The endpoint's successful receipt is gated by the inherited endpoint
barrier. Its request context can cancel after application, in which case the
handler cannot emit a successful receipt. The post-death release record proves
the ordering of the controller's write, independent of the handler's outcome.
