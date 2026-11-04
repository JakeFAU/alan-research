# Phase 3B.2 frozen authority contract

## Baseline and research question

The user supplied [the milestone brief](00-brief.md) and authorized research,
executable witnesses, independent review, manuscript reconciliation and a draft
PR, without merge or release. PR #42 is merged at
`1d41d1e20b579443c98a6caa53e51036e55f2d82`; its tree equals tested head
`e890a301d4c8cfead76f67676e00b6b4a2bd44e4`. The Phase 3B.1 manifest's 36 source
and 161 evidence records were verified before starting. Its models, positive
controls, 144-state serial obstruction and archived evidence remain unchanged.

Question: which explicit authority contract permits the two original services
to complete without granting a principal a permission it did not receive, or
inventing a permission tuple that no root possessed? Safety and existence of a
completion suffix are separate targets.

## Scope and chosen construction

Extend `research/mutualobligations` by embedding its unchanged State and using its
actual Successors for request/custody/recovery steps. New authority steps project
to stutters; old Apply is additionally gated by the new authority/admission state.
Do not treat the old fixed Binding.Grant string as evidence of the new permission.
The new authority binding names the same original caller/callee/generations and
operation sequence, plus a scenario-specific complete permission tuple. This
explicit mapping is a custody projection, not a proof of authorization.

For BA=b the authority operation names counter-B although the inherited BA
custody binding names counter-A. This is explicit symbol renaming. The Explicit
case places both operations on counter-B, so it does not retain the earlier
disjoint-resource assumption: its per-original opaque results abstract away
shared counter values and interference. No shared-resource algebra or isolation
theorem follows from this custody projection.

The main extension fixes A servicing and B serial, the previously qualified
asymmetric service placement. It retains two originals, two opaque result values
per original, three one-slot channels per edge, retry/duplicate cycles and one
recoverable crash per actual agent. Original all-four-policy regressions remain
independent controls. No graph depth cutoff or randomized sampling is permitted.

A one-hop grant is enough to expose the first boundary. Longer returned claims,
rootless cycles and scope mixing are independent malformed-ledger controls; no
valid arbitrary redelegation protocol is claimed. Revocation has at most one
monotone change to the selected B root permission. Cancellation propagation,
owner transfer and independent durable owners remain outside this slice.

Alternative considered: implementing distributed grants before qualification.
It would conflate missing runtime representation with the research result, so
this slice instead records the exact implementation boundary. Alternative:
irrevocable one-operation permits at the first authorization. That is coherent
but differs from the current fresh physical-admission gate; the chosen model
rechecks validity at each new admission.

## Complete permission tuples

Use the finite permissions below, interpreted as complete triples rather than
independent operation/resource/scope labels:

| Name | Operation | Resource | Scope | Initial holders |
| --- | --- | --- | --- | --- |
| a | increment | counter-A | red | A |
| b | increment | counter-B | blue | B |
| r | read | status | shared | A and B |
| x | increment | counter-A | blue | nobody |

Thus C_A = {a,r} and C_B = {b,r}. x exists in the operation vocabulary but is not
a capability anyone initially owns. A grant's permissions form a set of these
complete triples. Union means union of triples. It never means a Cartesian
product of the sets of operations, resources and scopes. Request permission is
the original communication relation; it does not imply execution permission.

An explicit grant retains its identity, issuer/grantee logical identities and
generations, original root/parent, complete permission set, exact original service
request and operation, authority revision, and no-redelegation disposition.
The grant can restrict a root's permission set and is additionally narrowed to
one original operation. A capability-shaped message/result is data, not a grant.

## Candidate policies and independent safety policy

S1 accepts only the serving principal's own currently valid root permission.
S2 accepts S1 authority or an explicit rooted, non-amplifying, request/operation-
bound grant from the requester. Both require a complete tuple to be covered by
one valid authorizing path.

Freeze S3 as S2 plus an ambient fallback that borrows a requester's current root
permission by identity, without an explicit grant for this service. S3 is a
candidate operational rule checked against the same explicit-consent safety
policy as S1/S2. It may preserve the global tuple union while violating principal
or service-binding authority. Its rejection under that policy does not establish
that all possible identity-based delegation policies are invalid.

A separately named Cartesian-composition mutant accepts a tuple by mixing
components from valid roots/delegates. It is a negative control for x, not a
silent change to the meaning of S3. Report ambient consent violations, missing
origin, binding/attenuation violations and tuple invention distinctly.

## Temporal boundaries

Each original has an immutable first authorization decision: undecided, allowed
with its exact authority witness, or refused with the violated condition. It is
not rewritten by later grant retention, revocation or recovery. The authorization
boundary follows request admission and checks full identities, root ancestry,
attenuation, operation binding and validity at that decision.

Each fresh DispatchAdmit requires the original allow and current validity of
that same authority path. This is distinct from the inherited Admit of a request.
Its opportunity belongs to the serving physical activation. Application requires
that opportunity and the inherited service guard. A later revocation prevents
fresh admissions but does not retroactively erase an already admitted opportunity,
completed application or retained result. Crash clears the affected server's
volatile opportunities, preserves grant/decision/refusal facts, and changes
activation on recovery. Re-admission must recheck current validity. Keep enough
retained admission/refusal evidence to check order without inventing a new grant.
Revoke affects B's b root and grants derived from it, not B's shared r root.
DispatchAdmit is disabled while the same original has an outstanding opportunity;
Apply consumes that opportunity. Both initial physical activations may be 0;
recovery increments only the recovered agent's activation as in the base model.

The model retains Phase 3B.1's atomic application-plus-receipt premise. There is
no lost external receipt between those steps. A crashed, not-yet-applied abstract
opportunity can be admitted anew only under this premise and current authority.
This is explicitly not the generic adapter runtime's unknown-outcome recovery
rule; the runtime qualification must preserve that disagreement.

In the positive delegated scenario the exact grant is retained with request
commitment, before any authorization attempt. Grant issuance itself requires a
valid parent at that moment; a rejected grant claim confers no permission. The
separate late-grant scenario permits authorization first and preserves its
permanent refusal. Do not remove losing authorization transitions to create a
universal progress result. A revoked root can legitimately make completion
impossible even though safe custody remains intact.
If b was revoked before inline-grant commitment, the request still commits with
a retained rejected issuance and no valid grant. Do not suppress request Commit.

## Minimal experiments and observables

Main cases under S1/S2/S3:

1. Local: AB requests b at B; BA requests a at A. No delegation required.
2. Explicit: AB requests b; BA requests b with an exact B-to-A grant for BA.
3. Missing: same original operations as Explicit, with the grant removed.
4. Scope mixing: AB requests b; BA requests x; BA's commitment retains an otherwise
   legitimate grant for b to A, which must not combine with a to authorize x.
   The grant binds the original operation identity, but its tuple set does not
   cover the requested x. Grant validity does not imply request satisfiability.

Additional bounded controls: shared-r service; late grant and sticky refusal;
one root revocation before authorization, between authorization/admission,
after admission and after application; the named Cartesian mutant; and oracle
mutations of grant parents, generations, request/operation bindings, tuple sets,
admission epochs and result-as-authority claims. Freeze any added cases visibly
before their final enumeration; do not silently expand the capability universe.

Complete retains its previous meaning: both original authentic results actually
consumed. A refusal is retained evidence but does not pay an original result
obligation. Report raw operational completion reachability separately from paths
whose entire prefix/suffix is authority-safe. Keep an irreversible unsafe-history
marker or equivalent retained violation evidence so state merging cannot erase
an earlier illegal authorization/admission/application. Report finite states,
arcs, violating transitions and shortest witnesses, not a finite total count of
histories in a graph containing retry cycles.

Custody checks identify each original request, server, authorization, application,
retained result, publication, delivery, recording and actual consumption. At least
one nonvacuous applied/released-receipt control must survive interruption.

## Independent validation and runtime realization

The independent ledger validator consumes event histories and walks retained
authority ancestry with a visited set. It must not call the operational model's
Successors, authorization predicates or transition validator. It independently
checks root origin, identity adjacency, subset containment, no redelegation,
service binding, issuance/authorization/revocation order, admission activation,
and the original application/custody chain. Mutations must preserve superficial
outcomes/counters while changing decisive authority edges.

Current source inspection establishes that public Permission is a static tuple
of kind/source/target/operation/revision/obligation. Source disposition,
adapter permission/reservation and physical admission are separate boundaries.
There is no generic retained request-bound B-to-A grant with an attenuation
chain. ReceiptGrant is specialized receipt-query/adoption authority. S2 therefore
has no claimed direct runtime realization here; do not add a broad grant system.

After a decisive model candidate, add a separate extension of the same-Store
ordinary authored fixture while preserving the previous files. Use distinct
fixed adapter configurations for concrete resource/scope authority. Exercise
legitimate local service, genuine grant-denied refusal, untrusted claim payloads,
and real SIGKILL/fresh-owner recovery. Removing an unwanted permission before
Create is public provisioning, not database editing. Result/refusal observations
must never be recast as delegated authority or successful service.

An allowed generic adapter attempt without a retained result is not retransmitted
by the current runtime. Kills after reservation but before call may therefore
reopen safely blocked even when the test endpoint saw zero contacts. Preserve
this fact, and distinguish source authorization, effect permission, physical
admission, result retention, settlement and consumption in every cut. Exact
delegation-retention cuts are inapplicable when the primitive is absent; static
provision retention is not a substitute. No independent-owner recovery is claimed.
