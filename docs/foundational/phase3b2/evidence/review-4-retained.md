# Independent review 4: retained authority and recovery evidence

Verdict: **REJECT** pending the runtime inspection-oracle correction below.
Task 2's independent finite ledger is **APPROVE within its declared finite scope**;
this overall rejection concerns the runtime retained-custody qualification.

## Scope and method

Reviewed the frozen `01-contract.md`, Task 2 brief/report/review patch and actual
oracle, corpus and tests; Task 3 report/review patch and actual current runtime
test source. Reviewed retained-authority origin, complete tuple attenuation,
identity/generation/service binding, first decisions, admission/revocation order,
activation recovery, unsafe candidate histories, custody, graph merging and
runtime unknown/corruption assertions. Source inspection and the two new probes
below are this review's evidence. Previously reported full graph, race, runtime
and Staticcheck commands were not independently rerun here.

## Blocking finding R4.1: consumption history is not bound to logical generation

`internal/durabletest/phase3b2_runtime_test.go:599-605` (`selectedBefore`) matches
observation ID, subject and sequence but omits the history row's generation.
The final paid-result selection at lines 210-217 omits generation too. Therefore
histories for a different logical run can establish original request consumption,
adapter result consumption and final original-result consumption.

Executed a fresh positive Local fixture using the existing real endpoint and
child-process helpers, first established `phase3b2Paid == 2`, then changed only
`generation` from 1 to 2 in all six history rows with a selected observation.
No observation, consumed sequence, intent, authority, result, settlement, payload,
endpoint value or other history field was changed. The six affected selections
were A sequences 5, 7, 11 and B sequences 5, 11, 13. The resulting inspection has
no generation-1 history selection for these observations. Nevertheless:

```
changed 6 history generations, retained paid=2 chain=<nil>
```

This is an accepted missing-edge control preserving final result values and
counters. It falsifies the claimed independent exact original consumption chain;
it does not establish a production runtime vulnerability. Runtime Replay was
not asked to accept this malformed inspection and may reject it independently.

Correction required: match the complete `RunID` (subject and logical generation)
for both intermediate and final selections, with a retained negative control
that preserves all result and observation rows. Re-review the correction before
marking this scope approved. A helper returning an ordinary error for the full
paid-chain check would permit negative testing of final consumption without
relying on `testing.T.Fatal`.

Probe and actual before/after inspections are retained temporarily at:

- `/tmp/alan-phase3b2-review4-runtime_test.go`
- `/tmp/alan-phase3b2-review4-overlay.json`
- `/tmp/alan-phase3b2-review4-before.json`
- `/tmp/alan-phase3b2-review4-after.json`
- `/tmp/alan-phase3b2-review4-runtime.log`

Executed (using the same local-toolchain/offline/cache environment as Task 3):
`go test -overlay=/tmp/alan-phase3b2-review4-overlay.json -tags=integration ./internal/durabletest -run '^TestReview4HistoryGeneration$' -count=1 -v`.
Initial sandbox execution could not bind the loopback endpoint. The authorized
escalated run executed successfully and demonstrated the bad acceptance; a second
run retained the exact before/after input JSON. These are review witnesses, not
final frozen-source qualification artifacts. The overlay added no repository
source and the positive database was not edited.

## Task 2 findings and independently attempted attack

`AuditTrace` and its ledger interpretation do not invoke operational successors,
authorization predicates, binding construction or transition guards. Literal
roots and tuple interpretation are independently represented. Grant paths use a
visited set and validate ancestry, issuer/grantee generations, full binding,
subset containment, revision and no-redelegation. Issuance evaluates parent
validity before retaining acceptance. First authorization cannot be rewritten;
fresh admission rechecks the retained selected path, while an already admitted
opportunity may apply after revocation. Crash clears the server's opportunity and
increments only its activation on recovery. Unknown runtime adapters are not
silently treated as this model's atomic application-plus-receipt transitions.

Candidate-invalid histories terminate without safety certification. Ambient and
Cartesian candidate-valid histories retain distinct irreversible safety flags
while custody validation continues. `auditGraph` checks one BFS prefix per node,
every outgoing edge, and equality of the entire independent ledger at each
merged target (including immutable grants, decisions, epochs and violations).
The wrapper still shares graph generation and event serialization; the report
explicitly excludes independent enumeration of the infinite history language.
There is no stronger universal or compiler-correctness conclusion.

A new deletion probe used `OracleCorpus`'s first complete R1/R1 control and
removed each of its 20 events in turn. All 19 deletions leaving a subsequent
required boundary were rejected, including missing commit, send, admission,
authorization, dispatch admission, application, publication, delivery, recording
and the BA consumption needed by B's serial service. Deleting the final AB
consumption correctly left a valid partial history with only BA paid. Thus no
false two-original completion survived this missing-edge attack. Probe:
`/tmp/alan-phase3b2-review4-model.go`, executed with `go run` in the report's
local/offline/cache environment.

The retained prior export `/tmp/authority-oracle-first.json` sums to 1,198,963
prefixes, 9,155,940 edges, 197,092 unsafe prefixes and 1,536,272 unsafe extensions,
with 33 histories and 4 chain controls, agreeing with Task 2's report. This
arithmetic check is not a fresh graph enumeration. The 18 malformed complete
histories preserve both R1 applications; controls distinguish equal values from
original provenance. Longer rootless/return chains are expressly ledger inputs,
not generated arbitrary-delegation histories. No additional Task 2 blocker found.

## Runtime evidence accepted with the stated blocker

The current source checks nine real process cuts and does not substitute reopen
for process interruption. It distinguishes source authorization, reservation,
physical admission, result retention and settlement. Allowed-without-result cuts
require zero endpoint contacts at interruption, exact unknown original identity,
no invented result/settlement, zero payment and unchanged endpoint snapshots
following recovery. Known retained results are consumed without physical reissue.
Existing immutable facts are compared across recovery and both physical epochs
must advance. No S2 rooted grant implementation or independent-owner recovery is
claimed.

The inspection oracle checks substantial authority/result/settlement/request
links, and the endpoint oracle independently matches result keys/values to
concrete resource/scope application ledgers. Existing negative controls retain
surface outcomes while removing request selection, changing a permission tuple
or substituting a successful result. Disposable-copy Replay corruption controls
are useful consistency checks, but run/replay agreement does not independently
validate the compiler. Their presence does not repair R4.1.

Final immutable source/evidence manifests and broad qualification remain the
parent integration task after all review corrections and source freeze.

## R4.1 scoped re-review disposition

**R4.1 ADDRESSED. Current scoped verdict: APPROVE.** The original REJECT and its
counterexample above remain historical evidence and are not erased.

Reviewed `.superpowers/sdd/02-plan/task-3-r4-fix.patch`, the appended Task 3
correction report, and the corrected runtime test source. Verified the latter's
SHA-256 is
`ab041fb061de1958ec24341ee5f3621d542c0a24eff2c28af47d89fec0c976da`.
The correction binds intermediate selected observations and history rows to the
complete RunID and exactly one prior selection. The new error-returning final
paid validator also requires matching logical generations; ordinary test callers
still invoke both the intermediate chain and final paid checks. Reply-origin
matching is strengthened consistently. No new breakage was found in this fix's
scope.

Independently reran the original archived-inspection attack against the corrected
checker through a temporary Go overlay, with no endpoint/process rerun needed.
The loader preserves SQLite integer column types when reading the retained JSON.
The unchanged positive inspection still passes both checks with paid=2. The exact
original six-generation mutation now fails both intermediate and final checks,
with validated paid=0. A second probe derived from that same positive inspection
moves only the two final consumption history rows: the intermediate chain still
passes and final paid validation rejects with paid=0. Thus the second check is
not passing merely because an earlier request/adapter edge already failed.

Command actually executed with the local/offline/cache environment:

```
go test -overlay=/tmp/alan-phase3b2-review4-fixed-overlay.json -tags=integration ./internal/durabletest -run '^TestReview4HistoricalWitnessFixed$' -count=1 -v
```

Result: PASS; package reported 0.235s. The all-six input produced
`service is not linked to consumed original request/result` and
`consumption lacks exact run/history selection`; the final-only input produced
the latter error while retaining intermediate-chain success. Probe source is
`/tmp/alan-phase3b2-review4-fixed_test.go`. The original before/after inspections,
pre-fix source and reproduction are now preserved in
`docs/foundational/phase3b2/evidence/review4-counterexample/`.

The implementer's reported 12.911s integration and 43.724s race runs were read,
not independently repeated in this scoped re-review. Approval concerns R4.1's
correction and the originally reviewed retained-evidence scope; final immutable
source/evidence generation and qualification remain with parent integration.
