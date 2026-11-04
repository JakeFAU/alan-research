# Task 2 implementation report

## Result and scope

Implemented the public same-Store request-order witness in
`examples/communication/mutual/` only: `deployment.go`, `deployment_test.go`,
`runtime_test.go`, and `README.md`. Existing runtime, author, and Phase 3A sources
were not edited. No commits or subagents were used.

`New` constructs two ordinary static roots with reciprocal request/reply ports,
separate original request identities and values, reply identity/value checks,
and final unsupplied local hold inputs. The four profiles are serial/serial,
A servicing, B servicing, and both servicing. There are no external effects,
new runtime operations, invocation/composition roots, or privileged roles.

The public prefix executes actual lowered instructions until each original
request intent appears. Each root commits its `validated=false` initialization
and then its Send, so the prefix is four commits, not two. Both requests are
then publicly delivered and acknowledged before ordinary Run starts. Both Send
orders and both delivery orders are exercised for every profile.

## Executed evidence

The assertion-level RED used a constructor stub returning an empty provision;
the adverse runtime test failed with `two ordinary roots required, got 0`.
The retained proper command exited 1; `/tmp/alan-phase3b-runtime-red.log` records
the failure. The final source passes every focused gate below.

All Go commands used:

```sh
export GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build
export GOPROXY=off GOSUMDB=off
export STATICCHECK_CACHE=/tmp/alan-phase3b-staticcheck
```

| Command | Exit | Log |
| --- | ---: | --- |
| `go test -tags=integration ./examples/communication/mutual -run '^TestReciprocalRuntime/serial_serial/send_a/deliver_a$' -count=1` (stub) | 1 | `/tmp/alan-phase3b-runtime-red.log` |
| Same focused serial command after construction | 0 | `/tmp/alan-phase3b-runtime-serial.log` |
| `go test ./examples/communication/mutual -count=1` | 0 | `/tmp/alan-phase3b-runtime-constructor.log` |
| `go build ./examples/communication/mutual` | 0 | `/tmp/alan-phase3b-runtime-build.log` |
| `ALAN_PHASE3B_RUNTIME_EXPORT=/tmp/alan-phase3b-runtime-normal go test -tags=integration ./examples/communication/mutual -count=1 -v` | 0 | `/tmp/alan-phase3b-runtime-normal.log` |
| `ALAN_PHASE3B_RUNTIME_EXPORT=/tmp/alan-phase3b-runtime-race go test -race -tags=integration ./examples/communication/mutual -count=1 -v` | 0 | `/tmp/alan-phase3b-runtime-race.log` |
| `staticcheck -tags=integration ./examples/communication/mutual` | 0 | `/tmp/alan-phase3b-runtime-staticcheck.log` |
| `python3 /tmp/alan-phase3b-runtime-audit.py` | 0 | `/tmp/alan-phase3b-runtime-audit.log` |

No full repository suite was rerun, following the task's focused verification
instruction. The normal and race profiles each execute 16 runtime cases and
four valid constructor profiles plus an unknown-profile rejection.

Observed per case:

| Profile | Cases per normal/race run | Admitted requests | Consumed requests | Admitted/consumed results | Unpaid application obligations |
| --- | ---: | ---: | ---: | ---: | ---: |
| serial/serial | 4 | 2 | 0 | 0 / 0 | 2 |
| A servicing | 4 | 2 | 2 | 2 / 2 | 0 |
| B servicing | 4 | 2 | 2 | 2 / 2 | 0 |
| both servicing | 4 | 2 | 2 | 2 / 2 | 0 |

Serial Run makes zero steps and dispatches. The repeated and reopened runs also
make zero steps and dispatches. Both original admitted peer requests are
explicitly proposed to their roots' current Receive and rejected; the complete
retained table snapshots remain unchanged after these probes. Reopen increments
both activation epochs by exactly one, while original IDs, program identities,
sequences, projections, observations, histories, admissions, requests, subjects,
and delivery acknowledgments remain stable. Replay agrees with the actual
retained RunStates, including the new epochs.

Each servicing Run performs 16 further ordinary commits and two result
dispatches. Both requests and replies have actual unique selected histories;
both roots set `validated=true` after checking their own result identity and
value, and remain open at an unsupplied local input Receive. No root terminates.

The serial runner reports quiescence, zero pending deliveries and zero
UnresolvedIntents despite two unpaid application obligations. This is a genuine
quiescent Store under runner semantics. Obligation counts are computed from
authentic request/result observations and their consumption/selected-history
links, independently of the runner bookkeeping.

## Raw artifacts for parent review/archive

- `/tmp/alan-phase3b-runtime-normal/`: all 16 normal cases.
- `/tmp/alan-phase3b-runtime-race/`: all 16 race cases.
- `/tmp/alan-phase3b-runtime-audit.py`: separate raw-fact recalculation.
- `/tmp/alan-phase3b-runtime-audit.json`: 32 checked cases, 8 serial and 24
  servicing, with raw counts, history/acknowledgment counts, report metrics and
  compressed evidence SHA-256 digests.
- The gate logs listed above.

Each case directory includes two exact `program-<id>.json` files, `summary.json`
and `evidence.json.gz`. The compressed JSON includes the full public provision,
four prefix commit pre/post states, accepted original request admissions,
complete public inspections/history/observation/intent descriptors, reports,
serial refused selections, and final replay. Serial inspections additionally
record rejected selection, repeat-run, reopen, and reopened-run cuts. Gzip has
no timestamp; fresh Store namespaces still differ, so fresh exports are not
claimed byte-identical.

The Python audit independently recalculates all observation consumption counts
and unpaid request identities from the raw observations and selected histories,
checks exact program SHA-256/raw-file identities, open nonterminal roots, serial
refusals and epoch changes, unchanged serial message/history evidence, and
servicing validation markers. It reports `CasesChecked=32`, `SerialCases=8`,
`ServicingCases=24`; it is not an independent semantics proof.

## Remaining boundaries

This qualifies a bounded request-order obstruction and its two asymmetric
repairs under the declared same-Store driver and ordinary scheduler. It does
not establish all-schedule progress or full model/runtime refinement. The
programs issue no retries, whereas the abstract model retains retry cycles.
Request selection and reply commitment are distinct ordinary commits; no
external mutation or atomic external effect/receipt assumption is qualified.

Whole-Store reopen renews both roots, so this is not independent-owner,
A-only process-crash, or distributed recovery evidence. Replay shares the
historical reducer and checks consistency, not independent compiler correctness.
Runtime delivery acknowledgments are distinct from the model's application-level
receipt-release ACK protocol. Final hold inputs deliberately remain unpaid and
are outside the two request/result application obligations used as the target.

Parent owns manuscript and research-document synchronization and evidence
archival. Source is finalized for parent review; no commit was created.

## Scoped review hardening and final freeze

Closed the two minor runtime-test findings without changing the constructor,
program semantics, case matrix, or expected counts:

1. Serial refused selections now require the exact current public error text
   `ineligible observation`. The runtime facade does not alias the underlying
   ErrIneligible sentinel, so matching that text prevents unrelated storage or
   stale-authority errors from satisfying the assertion.
2. Intent/observation correlation now checks both
   `origin.Run.Generation == observation.SourceGeneration` and
   `origin.Generation == observation.Run.Generation`, alongside the existing
   source/target subjects. Selected histories now also match the observation's
   target generation, alongside subject and consumed sequence.

Reran the normal integration and race integration commands from the table above
against this final test file; both exited 0 with all 16 runtime cases passing.
Focused integration-tagged Staticcheck also exited 0. The normal/race logs and
all 32 raw exports were regenerated at the same paths. The raw-fact audit was
strengthened with the same generation/error checks and rerun, exiting 0 with
32 checked cases, 8 serial and 24 servicing. Constructor/build evidence remains
the prior passing logs because those sources were unchanged.

Final `examples/communication/mutual/runtime_test.go` SHA-256:

```text
4b48dc9dcb1c8b82bb006dfcf2139457c5e41f58b4ab3d19419472a01abfceba
```

This follow-up changed only `runtime_test.go`, this report, and temporary
evidence/audit files. Source is frozen for scoped re-review; no commit was made.
