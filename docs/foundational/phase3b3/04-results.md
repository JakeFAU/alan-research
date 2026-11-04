# Phase 3B.3: request-bound delegation in the public runtime

Status: the authority primitive is reviewed and frozen; reciprocal-result,
process-cut and independent model/runtime qualification are in progress. This
document records completed evidence without anticipating those results.

## Question and the baseline obstruction

Can the actual runtime let A serve B's original request using one explicitly
delegated B permission, while preserving the original request, service and
result across interruption? The [contract](02-contract.md) separates authority
safety from completion reachability. It fixes two initial logical agents,
A servicing and B serial, one-hop delegation and one shared Store.

At public runtime `af6943a5805a9a6c0e8fe77f804264f4765d07ec`, the relevant
permission lookup uses the serving subject's static complete permission row.
The baseline fixture at `8319a9067c16aaf5e3ecccb3d0c62f5c536e6140` retains a
genuine source-allowed, permission-refused service when A lacks B's permission.
A claim in a message and a prior result do not supply a new permission.
The S1 control completes with actual independently observed effects. The
deliberate S2-completion assertion fails for the missing-authority reason.
Normal and race evidence was regenerated on the frozen baseline source.

This obstructs the specified service by A with the old API; it is not an
impossibility theorem for every possible encoding or a reason to substitute
a broadly authorized A. The baseline Cartesian labels exercise exact static
tuple lookup but use the same physical operation. The new qualification must
also distinguish actual physical operations.

## Primitive and its conditional contract

The reviewed runtime source is
`4f6f5f2500a48b2ca01ee7b09a823fa63656b7ab`. An opt-in composition profile in
storage formats 13/14 interprets two ordinary committed lifecycle operations:
`DelegateAuthority` and `RevokeAuthority`. They share the existing lifecycle
and composition journals, scheduler and resource accounting. No new authority
table, host-only grant API, scheduler or budget ledger is introduced.

A grant names a legitimate original holder and full permission row, the
holder's logical root locator, the recipient and its root/revision, an authentic
original message, and the exact future service operation. The service explicitly
selects its own or delegated authority path. Delegated service must immediately
follow selection of that grant's authentic original message. Root locators
are checked against actual membership; merely computing a namespace-bound root
hash does not establish that the root exists or holds the permission.

Grant retention survives reopening of the same logical computation. A physical
admission opportunity does not. Each fresh admission reconstructs the original
request, verifies the retained first decision and original debit, and rechecks
the same selected authority path at its current ordered prefix. The runtime
retains a primary `authority_admission` fact before releasing the worker.
Admission ordinals and current activation epochs remain distinct from the
original immutable wire request and authorization epoch. A refusal remains a
refusal after a later grant; reopening neither regrants nor charges again.

Revocation is monotone and names one original holder/root/permission tuple.
It does not erase previously admitted opportunities or results. A selected
invalid delegation cannot fall back to another authority path. Result
possession, request payload equality and another generation's authority do not
supply a grant. Normal completion of the delegator does not itself revoke the
grant. Retirement, changed logical generation/revision and explicit revocation
remain separate validity conditions.

These are implemented and tested protocol rules under the stated trusted
runtime/Store and adapter boundaries, not a formal proof of a minimum possible
implementation or a distributed authorization theorem. Opaque adapter wire
preparation remains trusted ABI output. The Store's admission journal proves
an opportunity under its retained policy; independent contact evidence is
needed to establish a call or application.

## Falsified candidates and corrections

The rejected literal-root contract required a hash that includes the Store
namespace before Create generates that namespace. A concrete construction
failure led to checked logical root locators; the rejected contract and RED
test are preserved.

Independent implementation review then rejected three findings. The candidate
inappropriately rejected legacy inert lifecycle-name text; typed inspection
accepted forbidden authority witnesses nested in legacy attempts; and the live
gate trusted an attempt/reservation pair that could be coherently rewritten
away from the original request. The last finding released an opportunity in a
malformed-Store probe, without making a physical call. It is not reported as an
observed unauthorized application.

The correction preserves legacy inert names without interpreting them as
authority, rejects both nested witness kinds outside the opt-in profile, and
reconstructs the derivable original request inside the admission transaction.
All sixteen expanded original-request mutation controls reproduce the old gate
failure and reject after correction. Both independent scoped re-reviews return
APPROVE. Rejected source, tests, raw witnesses and original review findings
remain retained; approval does not erase their history.

## A model-valid authorization the runtime refuses

The [selection audit](evidence/selection-narrowing/review.md) preserves a real
ordinary-runtime control with an exact accepted grant, the authentic BA request
already selected by A, and no revocation. A next selects grant-ready and then
commits the exact service operation named by the grant. The first service
permission decision refuses with `authority_request_not_admitted`: the immediate
predecessor is ready rather than BA. This reason does not mean the original
request was never admitted; its actual same-generation selection is retained.

The unchanged LateGrant model permits delegated authorization after the
corresponding admitted-request and retained-grant prefix. A generated eight-event
witness passes its independent AuditTrace, with both originals still unpaid.
The chosen placement of transport events is a permitted linearization of the
available causal evidence, not a claim that the Store records that exact global
transport order. Treating the extra ready selection as unrestricted stuttering
therefore gives a concrete model/runtime mismatch at authorization.

The runtime correctly enforces the deliberately narrower frozen contract.
The positive fixture selects ready first, then the original request, then the
effect. This additional ordering condition must remain in any realization
claim; an expected-refusal test passing does not turn the mismatch into model
agreement. The refused run reports Quiescent with B still open, both original
results unpaid and no endpoint contacts. Runtime quiescence alone does not
establish system completion.

## Finite progress query

An [independent query-design audit](evidence/progress-design/review.md) enumerated
the unchanged S2/LateGrant graph using only its public Initial/Successors
functions. For every reachable state, safe completion is reachable exactly when
BA's first permission decision is not Refused. The comparison covers all
111,065 states and 840,616 arcs: 107,213 states satisfy both sides, with zero
mismatches. Of the 3,852 remaining states, 1,926 have no grant and 1,926 already
retain a grant. Grant retention after an irreversible refusal does not repair
the original operation.

This is an exhaustive finite query, not an unbounded theorem. Its qualifying
set is not forward closed: authorizing BA before a grant can leave it by
refusing. The runtime handshake must separately enforce the required ordering.
The query includes states where a future grant can still arrive; it does not
assert that every qualifying prefix already holds covering authority. The
maximum shortest safe suffix is 23 model edges, an existential distance rather
than a scheduler-time bound.

The same no-BA-refusal predicate fails on S2/Revocable in 22,887 states. Revoking
B/b from the initial state is already a counterexample. The audit retains the
stronger exact predicate and its independent ledger checks. That graph uses
inline issuance, whereas the runtime retains a separate grant; it does not
supply a direct chronological mapping for a separate-grant-plus-revocation
runtime history.

Both queried graphs have zero incomplete states of outdegree zero, and every
complete state still has successors. Completion means both authentic original
results were consumed, not that ACKs are drained, principals are terminal or
the system is globally quiescent. These query-design results do not qualify the
new raw-runtime checker or establish completion from every runtime prefix.

## Primitive qualification and remaining work

The corrected source passed the full ordinary runtime module suite. Separate
three-package ordinary/integration runs and focused authority race runs retained
759, 1,586 and 105 test/subtest pass events respectively, with no failure or skip
events. Those are event counts, not independent experiment counts. Both scoped
Staticcheck commands passed. The ordinary two-agent smoke test establishes real
effects and server receipt selection; its empty acknowledgments do not yet
establish original-requester opaque-result custody.

Ten genuine old-created formats 3–12 fixtures passed 80 primary old/current
Replay, Inspect, InspectExecution and Open calls. All 60 paired exports match
byte for byte. Open changes only the expected epochs for still-open runs.
Formats 5/8 enable reconciliation but perform no receipt-read operation; formats
9–12 have no open runs at this cut. This is finite compatibility evidence, not
universal equivalence or fresh-activation coverage for every format. Frozen old
reader tests separately reject the new 13/14 formats.

Runtime source and qualification are retained in the public runtime repository,
under `docs/qualification/phase3b3`, through evidence commit `6352969`. The
research repository preserves the contract and its rejected predecessors.
The original baseline review changed SQLite SHM bytes during a read-only open;
that custody incident and reconstructed bytes are retained separately. The
definitive baseline was regenerated afterward. Later archive checks hash file
bytes and use disposable complete copies for SQLite access.

Outstanding phase obligations are the full result-custody chain, acknowledged
SIGKILL/reopen matrix, genuinely distinct Cartesian operations, temporal
authority/contact controls, independent raw-history checker and model corpus,
final frozen qualification, manuscript/PDF reconciliation and whole-phase
review. No unrestricted completion, independent-owner recovery, multihop or
arbitrary-agent composition result is claimed.
