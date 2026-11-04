# Task 1 report — reciprocal finite model

## Question and result

Does the fixed reciprocal durable-ACK composition preserve exact-edge custody,
and which service policies admit joint consumption from arbitrary reachable
prefixes? The implemented finite graph supports the conditional custody claim
and falsifies unrestricted pair completion under serial/serial service.

All four complete graphs have zero local-state, durable-transition and custody
violations. Serial/serial has `CannotComplete iff Trap`, exactly 144 states.
Either asymmetric service placement, and both servicing agents, have a
completion suffix at every reachable state. This is existential reachability,
not inevitable completion or a scheduling probability. Retry/duplicate cycles
remain in the graph; no exploration depth cutoff is used.

| Profile | States | Arcs | Cyclic states | Cyclic SCCs | CanComplete | CannotComplete | Trap |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| serial/serial | 68,497 | 511,552 | 56,580 | 22,384 | 68,353 | 144 | 144 |
| servicing/serial | 96,721 | 750,112 | 80,388 | 31,408 | 96,721 | 0 | 144 |
| serial/servicing | 96,721 | 750,112 | 80,388 | 31,408 | 96,721 | 0 | 144 |
| servicing/servicing | 96,721 | 751,744 | 80,388 | 31,408 | 96,721 | 0 | 144 |

Every profile has 16,384 Complete states. Nonvacuous custody counts are 68,272
applied / 45,824 released states for serial/serial and 96,496 applied / 63,232
released states for each other profile. These count states with at least one
original applied/released; the custody predicate checks each original. Targeted
traces release a retained receipt after durable ACK, put both physical agents
down, recover both and consume the surviving authentic durable result.

Serial shortest witnesses have lengths 2 (Trap), 6 (both admitted, every channel
empty, both up), and 10 (same wait after both crashes/recoveries). These are
separate exported witnesses. The recovered witness preserves the admitted wait
before the interruptions. Trap still has retries and is not transition
quiescence. Each profile also exports a repeatable unconsumed retry cycle with a
completion suffix and shortest completion/post-application/release witnesses.

## Exact implementation and qualification

Owned files, all newly created:

- `research/mutualobligations/model.go`
- `research/mutualobligations/explore.go`
- `research/mutualobligations/oracle.go`
- `research/mutualobligations/model_test.go`
- `research/mutualobligations/README.md`
- `research/mutualobligations/cmd/export/main.go`
- `research/mutualobligations/cmd/export/main_test.go`
- this report

Principal APIs: `Initial() State`, `Successors(Profile, State) []Edge`,
`Explore(Profile) Report`, `CheckTrace(Profile, []Event) error`. Helpers expose
Profiles, FixedBinding, Complete, Trap and Recoverable. A State has exactly two
actual Availability records and two compact EdgeState records. A physical
interruption disables both roles of that agent and clears only its caller-side
volatile receipt. Each original's identity, operation, generations, request,
grant and resource is fixed in its explicit event Binding. Both originals have
operation sequence 1 and may produce identical opaque results without merging
provenance. Strong ACK is enforced separately on each edge.

Only incoming Apply uses the executing agent's own outgoing Committed and
Consumed facts. No admission, result transmission, recording, consumption or ACK
operation is blocked merely by waiting. Every role requires its actual actor up.
The actual result is retained in an inaccessible per-edge oracle field and never
read to enable work or generate agent bytes.

Graph states are compact comparable structs, graph arcs use uint32 IDs and
uint16 event codes, and BFS stores only one predecessor per state. Report JSON
contains counts and finite shortest witnesses, not the full graph. Reverse CSR
completion reachability and Tarjan SCC analysis retain every arc. No new
module dependencies or production-runtime imports were added.

The independent joint causal ledger does not call Initial, Successors, Explore
or the local safety checks. Every BFS predecessor is independently checked once;
every reachable outgoing event is then checked against that ledger and its exact
independently derived target is compared to the graph target. Graph closure and
shortest predecessor distances are checked separately. Adversaries exercise
every caller and callee action while down in both directions, as well as wrong
actors, originals, grants/resources, generations, requests and endpoint binding;
fully rebound cross-edge receipt spoofing; receipt forgery; premature ACK and
consumption; deleted-receipt resend; duplicate application/consumption and ACK
payload leakage. Equal-valued original results have explicit positive controls.

Projection is a separate test against the unchanged actual Phase 3A durable-ACK
adjacency, indexed once. Each fixed reciprocal binding is consistently renamed
to Phase 3A's fixed original binding. The caller's interruption maps to old
crash/recovery; the callee's interruption and opposite-edge work stutter. Every
new arc is checked in both edge projections, independently of the new transition
function and joint ledger:

| Profile | Nonstutter old transitions | Stutters | Total |
| --- | ---: | ---: | ---: |
| serial/serial | 409,284 | 613,820 | 1,023,104 |
| servicing/serial | 604,708 | 895,516 | 1,500,224 |
| serial/servicing | 604,708 | 895,516 | 1,500,224 |
| servicing/servicing | 606,340 | 897,148 | 1,503,488 |

Normalization is a quotient for this conditional safety projection; the joint
checker, not normalization alone, verifies exact reciprocal authority. The
projection does not transfer completion/liveness when service steps are removed.
For every reachable state, independently perturbing either per-edge actual-result
oracle among NoValue/R1/R2 leaves enabled events and recoverable targets unchanged.
NoValue is a malformed/absent-oracle probe for applied states, not a third
authentic result. Comparing each R1/R2 replacement to NoValue additionally
rejects oracle-presence guards. This regression was added after the initial
independent spec/code/proof review; graph semantics and counts are unchanged.
Concurrent reader tests share input states/traces, run per-profile clients in
parallel and mutate caller-owned returned slices, making the race run exercise
shared-input isolation rather than only serial graph traversal.

## Commands actually run and evidence

All Go commands used this environment:

```sh
GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build GOPROXY=off GOSUMDB=off
```

Staticcheck also used `STATICCHECK_CACHE=/tmp/alan-phase3b-staticcheck`.

| Command | Exit | Evidence |
| --- | ---: | --- |
| `go test ./research/mutualobligations` with empty explorer stub | 1 | `/tmp/alan-phase3b-model-red.log` |
| `go test -count=1 -v ./research/mutualobligations/...` | 0 | `/tmp/alan-phase3b-model-normal.log` |
| `go test -race -count=1 -v ./research/mutualobligations/...` | 0 | `/tmp/alan-phase3b-model-race.log` |
| `staticcheck ./research/mutualobligations/...` | 0 | `/tmp/alan-phase3b-model-staticcheck.log` |
| `go build -o /tmp/alan-phase3b-model-export ./research/mutualobligations/cmd/export` | 0 | `/tmp/alan-phase3b-model-export.log` |
| `/tmp/alan-phase3b-model-export` twice | 0, 0 | `/tmp/alan-phase3b-model-export-{1,2}.json` |
| `cmp /tmp/alan-phase3b-model-export-1.json /tmp/alan-phase3b-model-export-2.json` | 0 | `/tmp/alan-phase3b-model-export.log` |
| `/tmp/alan-phase3b-model-export --profile serial/serial --assert-naive-completion` | 1, intentional | `/tmp/alan-phase3b-model-naive.json`, `/tmp/alan-phase3b-model-naive.stderr`, export log |
| `/tmp/alan-phase3b-model-export --profile invalid` | 2 | export log |
| `/tmp/alan-phase3b-model-export --profile servicing/serial --assert-naive-completion` | 0 | export log |
| `gofmt -l research/mutualobligations` | 0, empty | checked in task |
| `git diff --exit-code -- research/twoagent` | 0, empty | checked in task |

The genuine initial behavioral RED failed because the explorer had no reachable
graph/cycles, before model implementation. It is distinct from the later
intentional failure of the naive completion conjecture outside the passing
suite. An initial Staticcheck S1011 finding in a test append loop was corrected;
the original diagnostic is retained in
`/tmp/alan-phase3b-model-staticcheck-initial.log`.

The default deterministic export selects all four profiles; `--profile` selects
one named profile. Normal success is exit 0; `--assert-naive-completion` returns
1 for an actual reachable cannot-complete counterexample and emits JSON on
stdout. Bad arguments, invariant/witness or output errors return 2. Passing CLI
regressions distinguish these statuses.

The repeated JSON SHA-256 is
`28718bb9d29cd31eb870d05bad63360175abf105d1ee0724f96c36a6e7b43533`.
The serial-only intentional-counterexample JSON SHA-256 is
`b8424c2286e02e83993440a1ee7c43d2c02350e844eadc16de07d2fa32b5e8b9`.

## Source status, concerns and remaining obligations

Model sources were finalized before the final normal/race/Staticcheck and repeat
export runs. After the requested absent-oracle regression, the same final normal,
race and Staticcheck commands were rerun successfully (all exit 0), both exports
were regenerated and compared (exit 0), and their hash was also checked against
the pre-regression export. The graph counts and JSON hash are unchanged. These
final logs replace the corresponding paths above. The source manifest was
refreshed after the regression. Sources are ready for the fresh scoped review
and have not been committed.
Exact per-file SHA-256 identities, including frozen Phase 3A files, are in
`/tmp/alan-phase3b-model-source-sha256.txt`. Base HEAD at that freeze was
`8583527a8bc2d0b71a1e9115ded23eb96c8ac86c`, branch
`codex/phase3b-mutual-obligations`. The manifest records uncommitted new model
content rather than treating the base HEAD as its executable identity.

The build emitted a module metadata stat-cache permission warning under the
sandbox, while exiting 0 and producing the working binary. The warning is
preserved verbatim in the export log. No package/test/lint failure remains.
The root suite was not rerun by this task, per the parent's explicit instruction.
No production model changes, no Phase 3A changes, no commits and no dependencies
were introduced.

This task supplies bounded exhaustive graph evidence and a conditional safety
projection. The parent owns review, source/evidence archiving, manuscript
synchronization and the public runtime correspondence boundary. No formal proof
of unrestricted composition, production refinement or distributed ownership is
claimed here. Assumptions remain fixed disjoint authority/resources, intact
storage, one original each way, trusted bindings, atomic application plus callee
receipt retention and at most one recoverable crash per physical agent.
Cancellation/revocation, authority transfer, rollback, non-atomic effects,
unbounded crashes/resource pressure and universal fair-progress arguments remain
outside this finite model.
