# Reciprocal result custody and the serial service trap

This independent finite model asks whether two reciprocal Phase 3A durable-ACK
obligations preserve custody and admit joint consumption when a waiting agent
may postpone its peer's application. It models two actual agents, with exactly
one shared availability record per agent, rather than four separately available
roles. It imports no production runtime. The frozen `research/twoagent` package
is imported only by projection tests and remains unchanged.

## Contract and API

`Initial() State`, `Successors(Profile, State) []Edge`, `Explore(Profile) Report`
and `CheckTrace(Profile, []Event) error` are the principal interfaces. `Profiles`
returns serial/serial, servicing/serial, serial/servicing and servicing/servicing
in that deterministic order. The first policy belongs to A, the second to B.
`Successors` and `Explore` panic on an unknown profile; `CheckTrace` rejects it.

Agents A/B are numeric 0/1; obligations AB/BA are 0/1. Physical crash/recovery
events use `NoObligation=2` and an empty binding. Values are `NoValue=0`, R1=1,
R2=2. Event kinds are strings. Each edge event explicitly carries its actor and
full immutable `Binding`: caller/callee, both original generations, operation,
request, grant and resource. Both operations have local sequence 1 and use the
same opaque value alphabet; their original identities, grants and resources
remain distinct. These bindings occur in events, not repeated graph state.

`State` contains `[2]Availability` (Down, Crashed, Activation) and `[2]EdgeState`.
Each edge stores committed/admitted/applied/consumed markers; volatile, recorded,
retained and inaccessible actual result values; and one-slot request, result and
identity-only ACK channels. Booleans and uint8 enums keep states compact and
comparable. The graph uses uint32 IDs/arc offsets, int32 BFS parents and uint16
internal event codes; the fixed experiment fits those representations. Reports
export counts and finite shortest witnesses, never the full graph or unbounded
state histories.

A crash disables both roles of its physical agent, clears only its outgoing
volatile receipt and preserves both durable roles and all channels. Each agent
may crash once and recover once; recovery changes activation 0 to 1 without
changing original generation or authority. Every responsible actor must be up,
including admission, incoming application, result sending and ACK release.
Application and retained callee receipt creation are atomic.

Only incoming **Apply** is gated. A serial agent's gate reads exactly its own
outgoing `Committed && !Consumed`. A servicing agent may apply while waiting.
Admission, transmission, durable recording, consumption and ACK processing
remain enabled by their ordinary prerequisites during waits. Requests can be
retried and deliveries can retain a duplicate channel copy. Application and
consumption remain once-only; an applied tombstone cannot recreate a released
receipt. A strong ACK requires durable caller recording and carries no result
bytes or digest. The actual-result oracle is never a guard or receipt source.

## Finite graph result

BFS enumerates every reachable state and arc without a depth cutoff. Reverse
reachability computes existence of a suffix consuming **both** originals.
Complete does not require termination, channel drain, ACK release or quiescence.
Tarjan SCC analysis retains retry/duplicate cycles, including self-loops.

| Profile | States | Arcs | CannotComplete | Trap | Applied states | Released states |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| serial/serial | 68,497 | 511,552 | 144 | 144 | 68,272 | 45,824 |
| servicing/serial | 96,721 | 750,112 | 0 | 144 | 96,496 | 63,232 |
| serial/servicing | 96,721 | 750,112 | 0 | 144 | 96,496 | 63,232 |
| servicing/servicing | 96,721 | 751,744 | 0 | 144 | 96,496 | 63,232 |

Every profile has 16,384 Complete states. Applied/released counts count states
with at least one applied/released original. Every original is checked for
post-application custody in every reachable state: it retains a callee receipt
or a durable caller record. All local safety, durable-transition and custody
violation counts are zero. Tests also exercise release followed by simultaneous
down states and recovery of both agents, then authenticated consumption; the
minimal no-application trap alone would give no custody qualification.

For serial/serial, `CannotComplete iff Trap`, where Trap means both originals
committed and neither applied. The first incoming application would require
consumption of the other unapplied original, which lacks its causal receipt
chain. Recovery preserves those wait facts. The shortest Trap is two commits.
The distinct six-event witness has both requests admitted, both agents up and
all channels empty. Its ten-event recovered variant retains that same wait
across both permitted crashes/recoveries. These states still have retry cycles;
Trap is not zero outdegree or transition quiescence.

Either asymmetric policy is sufficient for existential completion from every
reachable prefix in this fixed graph. A servicing agent can pay its incoming
obligation; the peer then consumes, opening that peer's serial gate. Durable
receipts/records allow this argument to resume after interrupted prefixes.
Finite recoverable interruptions and eventual service of the needed steps are
premises for a scheduled completion argument. No universal liveness or fairness
claim follows from reverse reachability: every profile exports a repeatable,
unconsumed request retry cycle in a state that also has a completion suffix.

## Independent qualification

The joint `CheckTrace` is a separately implemented causal ledger. It does not
call Initial, Successors, Explore or local safety predicates. Its own frozen
binding table checks exact identities/resources, responsible roles, shared
availability, serial service, transport copies, durable ACK and receipt custody.
Tests build each BFS node's ledger from its checked predecessor, then validate
every reachable outgoing arc and compare the independently derived target state.
This avoids repeatedly reconstructing full prefixes. Adversaries cover every
caller/callee action while down in both directions; correlation, generation,
endpoint, grant/resource and actor forgery; fully rebound cross-edge receipt
spoofing; early ACK/consumption; receipt forgery/deletion and repeated application
or consumption. Both identical opaque results and distinct per-edge choices are
reachable without merging provenance.

A separate check projects each edge onto the actual frozen Phase 3A durable-ACK
transition graph or stutter. Fixed reciprocal bindings are consistently renamed
to Phase 3A's single fixed binding. The caller's physical interruption projects
to its crash/recovery; the callee's interruption and opposite-edge work stutter.
Both projections are checked for every joint arc. The old adjacency is explored
and indexed once per test; it is not repeatedly explored for each new state.
Projection checks total 1,023,104 for serial/serial; 1,500,224 for each asymmetric
profile; and 1,503,488 for servicing/servicing. Normalizing bindings cannot itself
qualify authority: the independent joint checker supplies that exact-edge check.
Removing enabled steps preserves this local safety projection, not completion.

A further exhaustive perturbation independently replaces either actual-result
oracle with NoValue/R1/R2 in each reachable state and compares all enabled events
and recoverable targets. NoValue is an absent-oracle probe, potentially malformed
for an applied state, not a third authentic result. This also detects a guard
that reads only oracle presence. No oracle fact enables work or supplies
agent-visible data.
Concurrent-reader tests exercise shared inputs and caller-owned output slices
under the race detector. Passing bounded checks is evidence for this frozen
model, not a proof about an unbounded or distributed implementation.

## Reproduction and exit codes

From the repository root, with the pinned local toolchain:

```sh
export GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build GOPROXY=off GOSUMDB=off
export STATICCHECK_CACHE=/tmp/alan-phase3b-staticcheck
go test -count=1 -v ./research/mutualobligations/...
go test -race -count=1 -v ./research/mutualobligations/...
staticcheck ./research/mutualobligations/...
go build -o /tmp/alan-phase3b-model-export ./research/mutualobligations/cmd/export
/tmp/alan-phase3b-model-export > /tmp/alan-phase3b-model-export-1.json
/tmp/alan-phase3b-model-export > /tmp/alan-phase3b-model-export-2.json
cmp /tmp/alan-phase3b-model-export-1.json /tmp/alan-phase3b-model-export-2.json
/tmp/alan-phase3b-model-export --profile serial/serial --assert-naive-completion
```

The default export emits all four reports; `--profile` selects one named profile.
Normal success is exit 0. `--assert-naive-completion` intentionally returns exit 1
when a selected graph contains CannotComplete states, preserving JSON evidence
on stdout and a counterexample diagnostic on stderr. Invalid arguments,
invariants, checked witnesses or output failures return exit 2. Thus the naive
unrestricted pair-completion claim is attacked outside the passing test suite.
The CLI exit-code behavior also has passing regression tests.

The behavioral RED used an empty explorer stub and failed on the missing finite
graph/cycles; it was not an assertion of an unsupported universal theorem. Task
logs and deterministic exports use `/tmp/alan-phase3b-model-*`; source/evidence
archiving and manuscript synchronization belong to the parent research slice.

This fixed model assumes intact durable stores, fixed authority and disjoint
resources, one original each direction, trusted transport binding and atomic
application-plus-receipt retention. Cancellation/revocation, ownership transfer,
rollback, non-atomic effects, unbounded crashes and resource exhaustion remain
outside its domain. It does not establish production-runtime refinement,
distributed ownership or unrestricted compositionality.
