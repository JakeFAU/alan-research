# Bounded authority and reciprocal service

This package extends the unchanged `research/mutualobligations` state and actual
successor relation. It qualifies the frozen Phase 3B.2 contract in
`docs/foundational/phase3b2/01-contract.md`. It is an executable finite research
model, not a runtime delegation implementation or an unrestricted theorem.

## Boundary and representation

The profile is A servicing/B serial. Both original obligations retain two opaque
result values, all three one-slot channels, retry and duplicate-delivery cycles,
and one recoverable crash per actual physical agent. `State.Base` is exactly the
base state. Authority transitions stutter under that projection. Only inherited
Apply is additionally gated; no other base transition is removed. Crash clears
volatile admission opportunities for the crashing server, while the unchanged
base transition clears its caller-side volatile result.

Permissions are complete triples: a=(increment,counter-A,red),
b=(increment,counter-B,blue), r=(read,status,shared), and
x=(increment,counter-A,blue). A initially holds a,r; B holds b,r. `PermissionSet`
is a bitset of complete tuples, never a product of component sets. x has no root.
The base `Binding.Grant` strings are not evidence. `ServiceBinding` explicitly
maps the original caller/callee generations and operation to the scenario's
permission. BA=b renames the old counter-A custody operation; Explicit therefore
abstracts shared-counter interference with per-original opaque results.

Every `Event` is comparable and serialized with its original service binding,
actor, physical activation, root revision, custody result/duplicate flag, and
applicable grant or authority descriptor. `Activation` and `Revision` describe
the source state: Recover changes the selected activation from 0 to 1, and
Revoke changes B's b-root revision from 0 to 1. Unrelated root witnesses keep
revision 0. `HasGrant` distinguishes an issuance claim from the zero struct.
Inline issuance is retained on BA Commit even when rejected after revocation.

`GrantRecord` retains ID, root/parent, issuer/grantee generations, complete tuple
set, exact original binding, revision, no-redelegation and issuance disposition.
The single valid grant is B/g1:b → A/g1 for BA operation 1. ScopeMix retains that
same b grant for an operation requesting x: a valid grant need not satisfy the
request. No arbitrary redelegation protocol is implemented.

The first authorization decision is immutable. S1 uses server roots; S2 adds
an explicit grant; S3 adds ambient requester-root borrowing. All are evaluated
against the same explicit-consent safety policy. Ambient borrowing marks
`ConsentViolation`. The separately named Cartesian mutant combines a's resource
with b's scope, marking `TupleInvention`. Unsafe history is monotone and survives
recovery, result retention and completion. The operational predicates used for
these counts are not an independent oracle.

Fresh DispatchAdmit requires the original allow and current validity of its
selected path. An outstanding opportunity prevents duplicate admission, survives
later revocation, is cleared by a server crash, and is consumed by Apply.
Revocation changes only B's b root, not either r root. Admission refusal is
retained separately from the immutable first authorization decision. This model
retains the base atomic application-plus-receipt assumption; it does not describe
the generic runtime's unknown-outcome recovery after an admitted attempt.

## Complete enumeration

`Configurations` lists exactly 16 frozen graphs. `Explore` performs uncapped BFS,
retains every arc including self-loops, computes raw and safe reverse completion
reachability, and uses Tarjan SCC analysis for all cyclic states/components.
Complete means both authentic original results consumed. Safe completion also
requires every prefix and suffix to remain authority-safe. Reachability is
existential; no scheduler fairness, inevitable completion or finite history
count is claimed.

| Configurations | States each | Arcs each | Raw completable | Safely completable |
| --- | ---: | ---: | ---: | ---: |
| S1/Local; S2/Local,Explicit,Shared; S3/Local,Explicit | 103,361 | 793,736 | 103,361 | 103,361 |
| S1/Explicit,Missing,ScopeMix; S2/Missing,ScopeMix; S3/ScopeMix | 6,741 | 37,434 | 0 | 0 |
| S3/Missing; Cartesian/ScopeMix | 103,361 | 793,736 | 103,361 | 0 |
| S2/LateGrant | 111,065 | 840,616 | 107,213 | 107,213 |
| S2/Revocable | 220,564 | 1,740,832 | 193,681 | 193,681 |

The frozen-source qualification counts total 1,198,963 states and 9,155,940 arcs.
Each unsafe configuration has 98,546 unsafe states and 5,776 locally violating
transitions. SafeStates and SafeTransitions require an entirely safe prefix;
ViolatingTransitions counts the individual unsafe authorization/admission/apply
steps even if their source already has an unsafe prefix. All operational
invariant and custody checks were zero. Qualification uses Go source commit `f75e26984c84fe649f26d5e4eaf4ed70e7506e46`;
see the Phase 3B.2 results and evidence manifests for exact command scope and
the explicitly non-green pre-existing integration Staticcheck diagnostics.

The decisive shortest BA prefix is Commit → SendRequest → Admit → Authorize:
S1/Explicit and S2/Missing refuse; S2/Explicit allows the explicit grant;
S3/Missing allows ambient borrowing and violates consent. S3/ScopeMix refuses;
the separate Cartesian mutant invents x. Reports preserve shortest completion,
noncompletion, unsafe, temporal and custody witnesses, plus an unfair retry cycle.
Grant/refusal/revocation/admission/application/release/crash cuts remain visible.

The internal graph uses compact event codes and reconstructs full descriptors
from source states. Same-package independent checker tests can access
`explore(c)`, `g.nodes`, `g.edges(i)`, `g.trace(i)` and
`arc.event.event(c, sourceState)`. This supports cached-prefix checking without
serializing entire graphs or repeatedly materializing large histories. Finite
state merging is not a claim that only finitely many cyclic histories exist.

## Reproduction and controls

Use the repository's pinned local toolchain and existing dependencies:

```sh
export GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build GOPROXY=off GOSUMDB=off
export STATICCHECK_CACHE=/tmp/alan-phase3b2-staticcheck
go test ./research/authorityservice/...
go test -race ./research/authorityservice/...
staticcheck ./research/authorityservice/...
go build -o /tmp/alan-phase3b2-model-export ./research/authorityservice/cmd/export
/tmp/alan-phase3b2-model-export -config all
/tmp/alan-phase3b2-model-export -config S3/Missing -assert-authority-safe
/tmp/alan-phase3b2-model-export -config Cartesian/ScopeMix -assert-authority-safe
/tmp/alan-phase3b2-model-export -config S2/Missing -assert-universal-completion
```

The last three commands intentionally exit 1 after writing counterexample JSON.
Invalid arguments, invariant failures and output errors exit 2. Ordinary exports
exit 0. The initial behavioral RED reproduced base Apply without any new
authorization/admission; the completed extension blocks it. Targeted tests also
cover all revocation cuts, pre-revocation opportunity survival, server-crash
clearing and re-admission, sticky late refusal, rejected inline issuance without
suppressed commitment, shared-root independence, and released custody surviving
both physical crashes before both originals consume their results.

No production package or Phase 3A/3B.1 source is modified. The independent ledger, runtime qualification and manuscript are documented in
the Phase 3B.2 results. Model success does not establish that the public runtime
has a request-bound grant primitive.
