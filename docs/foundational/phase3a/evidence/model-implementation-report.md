# Task 1 report: independent two-agent model

Status: complete within assigned model scope. Source finalized for parent review;
no commits made. Only `research/twoagent/**` and this report were edited.

## Question and result

The abstract implication “local durable safety plus B's at-most-once original
request application implies A's useful recovery” is falsified in the delivery-ACK
profile. Both actual opaque results have ten-event shortest recovered-loss
witnesses; their complete recoverable states are identical after erasing only
`WorldResult`. The B tombstone and ACK contain no result bytes/digest.

Durable-ACK preserves the custody invariant on every reachable state of the
frozen model. Every such state has some path to A's durable authentic once-only
consumption. This is not joint termination, B release, transport drain, global
quiescence or universal eventual progress. The README supplies separate
source-independent authenticity, custody induction, two-world loss, existential
completion and additional service obligations. There is no runtime adapter,
production import, authority attenuation qualification or new Store bug claim.

## Files and API

- `model.go`: fixed identity/grant, two values, one A-only crash, separate
  durable/volatile receipt, consumption marker and one-slot asynchronous
  request/result/identity-only ACK channels. B runs while A is down. Delivery can
  retain a channel copy, giving explicit duplication cycles. WorldResult is
  assigned at application only, never used for agent enabling or receipt sources.
- `explore.go`: complete BFS, all graph edges retained, shortest traces, separate
  local state/transition safety and custody counts, reverse completion
  reachability and Tarjan SCC cycle counts.
- `oracle.go`: independent causal event ledger; no calls to model transitions,
  explorer, Initial or model safety checker.
- `model_test.go`: behavioral graph/counterexample tests, every-edge oracle
  prefixes, 14 negative trace cases, asynchronous B/late receipt, genuine retry
  and resend, durable record/consume crash boundaries, oracle and B-private-state
  noninterference, transition checker negative cases, graph closure/shortest
  paths, deterministic output and explicit unfair retry cycle.
- `cmd/export/main.go`: deterministic JSON export and
  `-assert-naive-recovery` expected-failure mode (exit 1 with actual witnesses).
- `README.md`: domain, exact field/event meanings, finite envelope, count table,
  conditional arguments, open boundaries and runnable commands.

Required API: `Initial() State`, `Successors(Profile,State) []Edge`,
`Explore(Profile) Report`, `CheckTrace(Profile,[]Event) error`.
`Report.Nodes` exposes the graph in Go but is omitted from JSON;
`Report.Trace(nodeID)` returns the shortest event list (including an explicit
empty slice at the initial node). JSON includes `ShortestCustodyViolation`,
`LossWitnesses`, `Indistinguishable`, all counts and fixed contract version.
Empty `LossWitnesses` is `[]`; optional witness/pair pointers remain `null`.

## Final graph evidence

| Profile | States | Edges | Cyclic states / SCCs | Local state / edge failures | Custody violations | Can / cannot complete |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| delivery-ACK | 407 | 1,848 | 361 / 155 | 0 / 0 | 64 | 391 / 16 |
| durable-ACK | 311 | 1,372 | 273 / 115 | 0 / 0 | 0 | 311 / 0 |

Shortest custody failure: 8 events. Both shortest recovered-loss traces: 10
(the canonical prose order, r1/r2): commit/send/admit/apply/send-result/receive/
send-ACK/release/crash/recover. Both loss nodes have `CanComplete=false`; the
pair has `RecoverableEqual=true`. BFS has no depth or visit cutoff. Cycles are
retained; SCC counts are genuine cycles, not reachability merges.

## Commands and actual evidence

Toolchain observed: `go version go1.27.3 darwin/arm64`;
`staticcheck 2026.2.2 (0.8.2)`.

All Go invocations used:

```sh
GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build GOPROXY=off GOSUMDB=off
```

Executed:

```sh
go test ./research/twoagent/... # initial assertion RED against inert skeleton
# -> /tmp/alan-phase3a-model-red.log, exit 1

go test -run TestLocalTransitionCheckerRejectsLostFacts ./research/twoagent
# -> /tmp/alan-phase3a-model-transition-red.log, exit 1 against inert checker

go test -v -count=1 ./research/twoagent/...
# -> /tmp/alan-phase3a-model-normal.log, exit 0 (final source)

go test -race -count=1 -v ./research/twoagent/...
# -> /tmp/alan-phase3a-model-race.log, exit 0 (final source)

STATICCHECK_CACHE=/tmp/alan-phase3a-staticcheck staticcheck ./research/twoagent/...
# -> /tmp/alan-phase3a-model-staticcheck.log, exit 0, empty output (final source)

GOMODCACHE=/tmp/alan-phase3a-go-mod go build -o /tmp/alan-phase3a-model-export ./research/twoagent/cmd/export
# -> /tmp/alan-phase3a-model-build.log, exit 0, empty output

/tmp/alan-phase3a-model-export > /tmp/alan-phase3a-model-export.json
/tmp/alan-phase3a-model-export > /tmp/alan-phase3a-model-export-again.json
cmp /tmp/alan-phase3a-model-export.json /tmp/alan-phase3a-model-export-again.json
# exit 0: byte-identical repeated export

/tmp/alan-phase3a-model-export -assert-naive-recovery > /tmp/alan-phase3a-model-naive-red.json 2> /tmp/alan-phase3a-model-naive-red.log
# exit 1, preserved in /tmp/alan-phase3a-model-naive-red.status
# stderr: FALSIFIED: local safety plus original-request deduplication does not
# imply useful recovery; see independently checked recovered-loss witnesses
```

Final repeated export SHA-256:
`b380e7653a843b03ad9fe036f8ec5ca03a966cb1151b89e8a3e36a3177baf385`.
Hash file: `/tmp/alan-phase3a-model-export.sha256`.

An initial Staticcheck attempt could not write its default cache outside the
sandbox; setting `STATICCHECK_CACHE` to the path above resolved it. An initial
successful Go build printed a default module-stat-cache write warning; the final
build uses a writable `GOMODCACHE` and its retained log is empty. No network,
sockets, external dependencies or paid endpoints were needed. Parent owns root
suite verification; I did not repeat it.

## Remaining boundaries

This model assumes a fixed valid grant and atomic B application/result storage.
It does not establish current Store correspondence, B/simultaneous-crash safety,
external-effect receipt recovery, cancellation/authority dynamics, two-way waits,
resource fairness, arbitrary numbers of agents or composition laws. Parent owns
manuscript/document synchronization and archival of final evidence. No unresolved
model test/lint failure remains.
