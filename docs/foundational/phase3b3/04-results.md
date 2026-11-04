# Phase 3B.3: request-bound delegation in the public runtime

Status: runtime implementation, reciprocal-result custody, process-cut
qualification and the independent model/runtime checker are reviewed and frozen.
Final research qualification and whole-branch review are in progress. This
document separates completed evidence from those outstanding gates.

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
tuple lookup but use the same physical operation. The new qualification distinguishes actual physical operations as described
below.

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

The separate correspondence package repeats the exhaustive query and checks
actual longest witnesses with the unchanged independent authority ledger. Its
reviewed result agrees with the design audit: LateGrant has a maximum shortest
suffix of 23 edges, Revocable 22. The respective complete audited histories
contain 25 and 24 events, including two prior crashes. The Revocable predicate
checks both first decisions and requires every unapplied original to retain its
own admitted opportunity after revocation. This is equivalent to the earlier
paper's characterization on the reachable graph; the prior model and its
formula have not been rewritten. Final committed-source execution is recorded
separately below.

## Runtime result and process boundaries

The final runtime source is `759df6e3d2505f27e1287b2cfeb07bc4290afc0a`.
The `Reciprocal` fixture contains exactly two ordinary authored roots: A can
service while its original is outstanding, while B waits for BA's authentic
reply before servicing AB. Both original commitments precede A's service.
The positive order is accepted grant, ready selection, original BA selection,
then the exact delegated effect. Endpoint result, server receipt selection,
reply body and requester selection are joined by full logical RunID, original
message, operation, observation and consumption sequence. Equal opaque payloads
do not collapse those identities. Each service spends one debit on its actual
serving root. No broad permission is added to A.

Each ordinary/race corpus contains 48 indexed cases: four reciprocal scenarios,
fifteen binding controls, eleven process cuts, four physical Cartesian controls,
three revocation timings, sticky early refusal followed by a late grant, actual
prior-result laundering, serial/serial waiting, generic unknown and revoked
applied-unknown cases, and a mutation base plus five malformed copies. Counts
include negative and unsupported-model histories; they are not 48 successful
completions. The Cartesian endpoints implement distinct trusted add-one/add-ten
operations. A hostile payload delta cannot select either operation; both real
permissions work separately and cross pairs retain refusal with zero
application, gate and debit.

| Cut | Actual retained boundary / recovery observation |
| --- | --- |
| 01 | Grant execution staged, before its commit; recovery executes the original program to obtain the grant |
| 02 | Accepted grant decision retained; original logical binding survives reopening |
| 03 | Service committed before source authorization |
| 04 / 05 | Coincident after-decision / before-gate hook; original allow and debit remain, fresh admission uses the new owner epoch |
| 06a | Gate retained before endpoint contact; one old unused opportunity plus one fresh opportunity, one contact/application |
| 06b | Endpoint application before local result retention; two contacts and two gates, one stable-key application |
| 07 | Result recorded before settlement; recovery uses that result without another gate/contact |
| 08 | Result settled before authored selection; recovery preserves the exact result chain |
| 09 | A waits on its in-flight receipt and B on BA's reply; two contacts and two gates, one application after recovery |
| 10 | A terminal, both original results already consumed; B waits for the terminal observation, so only B's activation renews |

Ten requested labels yield eleven actual acknowledged SIGKILL witnesses because
06 has two distinct subcases. Labels 04 and 05 deliberately share the available
runtime hook. Every process cut records the barrier acknowledgment, checked
signal/exit, distinct fresh owner PID, same original Store, coherent post-Open
snapshot before Run and independent endpoint ingress/application evidence.
Cut 09 is an actual receipt/reply wait, not the separate serial/serial dependency
obstruction. Cut 10 does not leave either original result unpaid.

A gate is an opportunity, not evidence of contact. An allowed pre-revocation
opportunity may apply afterward. Revocation before the fresh gate blocks that
path while preserving the original allow and debit; it does not revoke B's
other, distinct service root. Recorded results remain consumable after
revocation. The generic adapter applied-unknown control makes no fresh call and
remains `UnknownAdapters`. The idempotent applied-unknown/revoked control retains
one prior application, refuses a second gate, and remains `UnknownExternal`
without an invented result or reply. Neither is relabeled a completion.

The temporal revocation controls use explicit public retained transitions at
an acknowledged gate barrier, separately identified from ordinary automatic
Run. They do not prove automatic reachability of that interleaving under the
single synchronous runner. No permission SQL edits or host-minted grants create
these authority histories. Only the separately named malformed-Store controls
edit storage, and their untrusted raw exports remain available independently of
Inspect/Replay acceptance.

The runtime services select two distinct B permissions/endpoints. Mapping each
whole action to the model's B/b is a many-to-one abstraction for selected
nonrevoking traces. It does not preserve the permission algebra, shared-counter
arithmetic or revocation scope. Separate-grant-plus-revocation and
applied-without-local-receipt/duplicate-contact prefixes remain outside the
frozen model's direct chronological domain. The raw checker must still validate
their actual temporal authority and custody facts; unsupported model mapping is
not a waiver of a raw policy failure.

Whole-Store process interruption also lacks an established chronological
projection to independently crashing modeled agents. A final service subsequence
can pass its abstract ledger without proving that the intervening process
failure maps to a pair of model crash/recovery steps. The checker labels this
limit explicitly; it inserts no invented crash steps to make the correspondence
appear complete.

## Independent checker and preserved false acceptances

The new research-only package reads the exported source, permission tuples,
deployment, roots, original history, ordered decisions, gates and observations.
Its raw policy oracle imports neither the runtime/executor nor its authorization,
Replay or SQLite implementation. A separate projection invokes unchanged model
successors only after raw validation. Five authored instruction forms are in
scope: emit, receive, branch, set and end. This is not a general compiler or
reducer proof. Exact integers, duplicate keys, primary keys, descriptor encodings
and the declared table/column sets are checked before policy reconstruction.

The checker package contains 22 selected/generated model cases in the eight required
classes, retains the earlier 33 histories and four ancestry inputs, and derives
the positive BA delegation from actual model transitions. Its runtime reports
separate raw validity, external evidence, prefix continuity, supported projection,
named mismatch and supplied runtime-reader dispositions. Honest authority refusal
can be a valid raw history. Unsupported model projection cannot excuse a raw
policy contradiction. An explicit expected-outcome file controls the CLI exit
status only; it does not change the reported verdicts or erase the known mismatch.

Review found five substantive checker false acceptances across development.
The parent's probes removed a current subject while preserving history, and
placed receipt selection before its result fact. Corrections require the exact
current subject set and the actual selected observation's causal predecessors
at the current composition prefix, with separate message, lifecycle, receipt
and terminal controls. A final snapshot containing a result cannot justify an
earlier selection by itself.

The independent raw-policy review then reproduced three more gaps: accepted
messages without their original route permission, changed original authority
between individually valid captured prefixes, and a stale gate/event epoch
contradicting the fresh owner's Open evidence. All exact adverse inputs, probe
programs, false-positive outputs and the rejected source are retained in
`evidence/task3-development` and `evidence/task3-rejected-review`. These findings
concern missing joins inside the supplied evidence, not cryptographic
authentication of unsigned exports. The corrected source requires exact original
route permission, conservation of immutable source/authority sets across captured
prefixes, and agreement between each supplied owner's actual Open/post-open
state and its boundary/gate epoch. Captured reopening increments each open run's
epoch once and preserves terminal epochs. The original three probes now fail at
`admission_permission`, `immutable_source_changed` and `boundary_owner_epoch`.
The continuity and epoch probes preserve their individually raw-valid snapshots;
their contradictions arise at the cross-evidence joins.

The old holder mutation is now named accurately, with a separate actual recipient
mutation retaining exact before/after cells and external receipt rows. All fifteen
raw mutants reject. The recipient control first fails original-source integrity;
it is not presented as having reached a later recipient-policy branch. The loader
also reuses its digest-verified JSON/JSONL bytes throughout the checks, eliminating
discarded reread errors. Both independent Task 3 reviews returned APPROVE after
correction; the original REJECT and failed intermediate layout assumptions remain
preserved in `evidence/task3-round1` and its predecessor archives.

The reviewed source is `7894bfae268e6dc0f5418bda3c6022f3a8a55540`.
Candidate ordinary/race qualification each had 58 passing test/subtest events,
zero skips/failures and clean package Staticcheck. Each complete raw corpus has
48 cases, 117 captured phases and 1,761 verified artifact digests: 43 valid final
raw histories and five hostile rejections. There are no external or continuity
failures for the valid histories. The only named correspondence discrepancy is
`immediate_selection_narrows_model_admission`; no projection apparatus error is
hidden as an unsupported domain. Final committed-source verification below is
separate from these candidate runs.

Two controls lack an initial pre-Open snapshot. Their first-owner report explicitly
sets `RenewalChecked=false`; it checks supplied Open/post-open and gate agreement
without claiming an unobserved initial renewal. Subsequent stopped-to-reopened
renewal remains checked. These are internally consistent, unsigned contextual
records. Neither Store bytes nor this checker establish authenticity, complete
contact attribution, arbitrary-program validity or full hostile-history detection.

## Frozen runtime verification

The final pipeline bound 304 Go/module inputs before and after execution, with
Go 1.27.3 darwin/arm64 and Staticcheck 2026.2.2 (0.8.2). All eight commands exited
zero. The full ordinary module reported 1,265 test/subtest pass events across
18 packages plus three packages without tests. Core integration reported 1,597
pass events across four packages; focused authority race reported 105. Each
new reciprocal run reported 53 pass events and 48 indexed cases, with no skips.
The baseline-on-current-source check reported seven pass events and one deliberate
opt-in expected-failure skip. Historical baseline RED remains separate evidence.
Full runtime Staticcheck passed in ordinary and integration modes. Event counts
include subtests and do not count independent experiments.

The final archive retains 3,784 files. All member hashes and all 1,761 per-case
artifact digests in each corpus were checked. Executable source remains the
freeze above; evidence delivery is runtime `2b8bb82`, under
`docs/qualification/phase3b3/task4-frozen`. Development failures, rejected
candidates and scoped independent review reports are retained alongside it.
Independent recovery/admission and binding/custody reviews both returned
APPROVE. They checked actual witnesses rather than treating passing tests or
reviewer agreement as a proof. The subsequent whole-change review also returned
specification, quality and overall APPROVE within the declared bounded scope.

Runtime main later advanced to `76031c34c1d0eeff325e1a8e9ca4b6e2fa2b6315`
with the separately qualified composition-mode cache. A clean prospective
combined tree was checked without merging or rewriting either branch. Core
integration passed 1,611 test events, authority/mode race passed 119, all 48
reciprocal race cases passed, and root-module integration Staticcheck was clean.
The initial archive-only attempt failed the harness's Git-HEAD export requirement;
that apparatus failure and its raw outputs remain separate from the corrected
same-tree qualification. Runtime `current-main` evidence commit `ed2d086` records
the exact input commits, tree and disposable source identity. This supplement does
not replace the original frozen runtime experiment or its research checker inputs.

Ten genuine old-created formats 3–12 fixtures separately passed 80 primary
old/current Replay, Inspect, InspectExecution and Open calls. All 60 paired
exports match byte for byte. Open changes only the expected epochs for still-open
runs. Formats 5/8 enable reconciliation but perform no receipt-read operation;
formats 9–12 have no open runs at this cut. This is finite compatibility evidence,
not universal equivalence or fresh-activation coverage for every format.
Frozen old-reader tests reject formats 13/14. The reciprocal process cuts use
format 13; no format-14 process qualification is claimed.

The original baseline review changed SQLite SHM bytes during a read-only open;
that custody incident and reconstructed bytes are retained separately. The
definitive baseline was regenerated afterward. Later archive checks hash file
bytes and use disposable complete copies for SQLite access. The six pre-existing
research integration-Staticcheck diagnostics are retained in
`evidence/research-baseline`; they are distinct from the clean runtime lint runs.

## Frozen research verification

The [final research archive](evidence/research-frozen/README.md) binds source
`7894bfae268e6dc0f5418bda3c6022f3a8a55540`, all 584 Go/module inputs and the
positive raw fixture before and after execution. All 563 base Go/module inputs
remain byte-identical. The 42 command records have their declared exit status;
fourteen are expected nonzero results, not passing assertions. They comprise
seven retained model RED controls, two unqualified corpus verdicts, four killed
checker-source mutants and the unchanged integration-lint baseline.

The full ordinary research suite passed 2,382 test/subtest events across 48
packages. Three unrelated opt-in artifact-export tests were skipped, and fourteen
packages contain no tests. The model/checker race run passed 303 events across
seven packages, with no named test skips. Prior reciprocal integration passed
22 events in each ordinary/race run; prior authority integration passed 20 in
each. The new raw-corpus and exhaustive-progress tests were enabled explicitly
with `ALAN_TASK3_RAW` and `ALAN_TASK3_PROGRESS=1`.

All four previous model/audit exports are byte-identical to their retained
archives. The new commands reproduce 22 model cases in eight classes, fifteen
rejected raw mutants, and the two exhaustive progress predicates above. Both
48-case runtime corpora reproduce 117 captured phases and 1,761 checked digests,
five hostile rejections and the sole named selection-order mismatch. Explicit
expected-outcome and unqualified commands produce identical parsed reports;
only their exit status differs. Four independent source mutations that remove
duplicate-key, attempt-binding, fresh-policy or causal-prefix checking each fail
an actual test assertion on the frozen source, rather than a build error.

Formatting and scoped model/checker/fixture Staticcheck passed. The six inherited
durable integration diagnostics are byte-identical to the pre-change baseline;
whole-research integration lint is not claimed clean. All 1,375 archive members
and command-log hashes were verified without opening SQLite. Eight rebuildable
binaries are omitted with their hashes and exact build records retained. A
mistyped source-identity preflight failed before any qualification command ran;
that controller error is preserved separately from semantic RED evidence.

## Remaining work and limits

The [independent whole-change review](evidence/final-review/review.md) returned
specification, quality and overall APPROVE within the declared bounded scope;
draft delivery remains pending. The checker/corpus and rebuilt manuscript record
the tested protocol, actual cut distinctions and preserved correspondence
counterexample. No unrestricted completion, independent-owner recovery, multihop
or arbitrary-agent composition result is claimed.

The whole-change reviewer found one minor input-diagnostic limitation: the
mutation generator expects a complete S2 parent, and currently exits 2 with an
index panic on a raw-valid incomplete initial snapshot. Its fifteen qualified
controls use the complete positive fixture. This does not affect the raw checker
or qualify the generator for arbitrary valid prefixes. The exact probe is retained;
improving this diagnostic is deferred without changing the frozen executable.

## Direct answers to the research questions

1. **Can ordinary retained computation implement the delegation?** The public
   runtime's scoped opt-in extension demonstrates request-bound one-hop
   delegation with real two-agent service and recovery. The old static API did
   not express this specified service. The demonstrated extension preserves an
   exact request, recipient, service operation and original permission tuple;
   its positive witness also requires ready selection before original selection.
2. **Do the exact Phase 3B.2 conditions suffice for both safety and progress?**
   No unrestricted correspondence result follows. The runtime's additional
   immediate-selection condition produces the retained authorization mismatch.
   Selected positive executions complete safely under the stronger conditions,
   while late refusal, missing authority and revocation retain legitimate unpaid
   obligations. Exhaustive completion reachability is established only for the
   unchanged finite graphs and their stated predicates, not every runtime
   prefix or every fair execution.
3. **What additional primitive was required?** The tested sufficient mechanism
   is a typed committed grant/revocation protocol, an explicit effect authority
   selector and a durable fresh-admission fact with original-bound revalidation.
   It reuses the existing journals, scheduler and charges. This is not a proof
   that no smaller possible encoding exists.
4. **What must change for independent owners, multihop or three agents?**
   Independent owners require an explicit ownership/fencing and transport model
   that orders grant validity, revocation and physical admission across failures.
   Multihop requires a retained attenuation chain and a defined propagation and
   revocation contract. Three agents require new cyclic obligation and shared
   authority/resource counterexamples, followed by a compositional closure
   argument under explicit scheduling and endpoint assumptions. None of those
   extensions is established by this slice.
