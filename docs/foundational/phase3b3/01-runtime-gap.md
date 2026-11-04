# Phase 3B.3: verified baseline and runtime gap

Status: source audit and experiment specification, before implementation. This
document does not claim that a new delegation primitive has been implemented or
that current-runtime witnesses have passed.

## Question and baseline

Can an ordinary retained computation delegate a legitimate complete permission
tuple to the other existing agent for one authentic original request and service
operation, with fresh-admission validity checks? Separately, does that realization
preserve the finite model's completion reachability under its actual assumptions?

The binding [packet](00-brief.md) authorizes the experiment, narrowly necessary
runtime changes, independent review, manuscript reconciliation and draft PRs,
without merge or release.

GitHub inspection confirmed research PRs #42 and #43 merged. Research `main` is
`8197683c23a3ca774750fc3c2c389c5a64dd38b8`, the merge of #43. Its tree equals
delivery `d82a69c7a576f94599138cc6eddca0b8aeb6d259`. Its Go/module inputs equal
the Phase 3B.2 frozen qualification `f75e26984c84fe649f26d5e4eaf4ed70e7506e46`.
The new research branch starts from that verified merge. Prior models, fixtures,
archives and rejected candidates are preserved.

The separately fetched public runtime `github.com/JakeFAU/alan` is at
`af6943a5805a9a6c0e8fe77f804264f4765d07ec`. Its original checkout has an
untracked `.envrc`, which was not read or changed; implementation uses a clean
isolated checkout. The public runtime was extracted from research source
`003dc9c45ac9b0c1f48565632f5be80f0f05beed` and has since advanced independently
(`docs/extraction.md`). Research fixtures import `alan-research/runtime`, not
the public `alan/runtime`. Their earlier qualification cannot substitute for
testing the current public implementation.

The public runtime adds typed execution inspection, retained dispatch admission,
causal recovery validation and explicit continuations. Both repositories retain
the same exact static adapter-permission lookup. This is a source observation,
not a claim that the full implementations are equivalent.

## Frozen prior contract

The authority specification remains [Phase 3B.2's contract](../phase3b2/01-contract.md),
with the unchanged forward model and independent ledger in
`research/authorityservice`. Its sixteen complete graphs contain 1,198,963 states
and 9,155,940 arcs in the **prior retained qualification**; those numbers are not
new runs. S2/Explicit's completion result requires covering authority before the
first permission decision, validity at every required fresh admission, the fixed
A-servicing/B-serial schedule grammar and atomic modeled application/receipt.

The ambient-borrowing counterexample violates explicit consent while remaining
inside the global tuple union. The separate Cartesian mutant invents an unheld
tuple. Preserve that distinction. The historical generation-binding defect is
an evidence-validator defect, not a demonstrated production-runtime exploit.

## Exact missing representation

Pointers in this section refer to the public runtime at the SHA above.

| Required fact | Existing representation and gap |
| --- | --- |
| Legitimate complete permission | `internal/executor/store.go:Permission` retains kind/source/target/operation/revision/obligation. Fixed trusted endpoint or capability identities may represent concrete resource/scope tuples. No generic scope algebra is needed for this experiment. |
| One-hop B-to-existing-A grant | `internal/executor/adapters.go:adapterGrant` selects only a row whose source is the actual serving subject A. Message claims, results and invocation handles do not introduce a B-rooted alternate path. |
| Dynamic authority | `invocation.go` introduces permissions for newly admitted template members; `composition_binding.go:memberPermissions` derives them from the admitted program catalogue. Handoff introduces directed message routes. Neither grants B's exact root tuple to existing A for an original service. |
| First permission decision | `messages.go:AuthorizeIntent` records source disposition. `AdapterAttempt` separately combines permission with reservation. Source allow is not evidence that A holds the required capability. Neither current record names an S2 ancestry chain. |
| Fresh validity | `dispatch.go:registerDispatch` rechecks the serving generation, revision, status and owning-scope closure. It does not check a selected B root or selective root-tuple revocation. |
| Selective revocation | Existing lifecycle operations close, retire or retry a principal. They do not revoke B's b tuple while leaving B's r tuple valid. |
| Causal request/service binding | `author.Send` does not expose a request handle; `Receive` binds payload fields. A user-supplied request-id string is not the retained message identity. A primitive must verify the actual request admission and selection as well as the exact service intent. |
| Admission evidence | `DispatchAdmission` retains an operation/prepared digest/composition prefix for opted-in scoped adapters. It is not a general grant witness, and must not be assumed present for every fixture. |

All current production writes to the permission relation are provisioning and
admission of new template members or their message handoffs. Ordinary receipt or
selection of B's request cannot add an adapter permission for existing A. The
smallest executable falsifier therefore provisions B with b and A without b,
retains B's authentic request, and has A attempt b. The predicted current-runtime
outcome is source allow followed by retained `grant_denied`, without a physical
contact. A broadly provisioned A is an S1 control, not a successful S2 realization.

This is a precise current-API obstruction. It does **not** prove that a small
integrated primitive is impossible or exceeds the packet's scope. Existing raw
retained programs can emit opaque intents; a retained authority-control protocol
may reuse the existing runner and decision journals. A separate builder review
must settle the smallest representation before production changes.

## Independent progress and identity boundaries

1. Generic allowed adapter attempts without retained results are not retransmitted
   on recovery (`adapters.go:prepareAdapterWork`). This remains true even when an
   external test oracle observed zero contacts. A grant cannot turn that uncertainty
   into a receipt or permission to retry.
2. Existing receipt reconciliation requires the original source to be closed
   (`receipt_read.go:reserveReadLocked`), uses a specialized controlled-counter
   protocol, and does not turn an absent receipt into a new application. Closing
   the waiting serving agent does not preserve its original continuation.
3. The existing raw external **idempotent** path can retry its stable original
   key (`external.go:prepareExternalWork`). An endpoint with a durable atomic
   key/application/receipt ledger is a candidate positive recovery realization.
   This needs its own executed witness and explicit endpoint assumptions; it is
   not a change to generic-adapter recovery.
4. A Store reopen renews all open runs' epochs. It represents a restricted paired
   crash schedule, not independent ownership. Grants remain bound to logical
   identities across that reopen; physical admission opportunities do not transfer.
5. The finite model names its original service operation before execution. A
   runtime primitive must either use a checked deterministic operation locator or
   retain an explicit one-use binding. Equal payloads and arbitrary string labels
   are insufficient. Any immediate-request-selection restriction must be visible.
6. A separate grant commit can lose a race with the first service decision. The
   positive realization must establish a retained handshake or atomic request/grant
   relation; it must not silently remove the late-refusal schedule.

## First experiment and scope gate

Before introducing the primitive, execute ordinary current-runtime S1 service and
the missing-grant, grant-shaped-message and Cartesian substitution attempts.
Retain actual decisions, request/result custody and external contacts. Preserve
the missing-grant refusal as a historical pre-implementation witness.

The candidate primitive must integrate legitimate source provenance, exact tuple
attenuation, original request and operation/generation binding, immutable first
decisions, selective root revocation and fresh physical admission. It must reuse
existing scheduling and resource accounting. A new table is not yet proved
necessary; extending serialized semantics still requires explicit version and
compatibility qualification. Implementation belongs in the runtime repository;
the contract, independent model comparison and interpretation belong here.

Do not stop at the missing-API observation, and do not expand into a parallel
authorization subsystem. If a concrete design or execution fails, retain it and
classify exactly which authority, progress or representability claim it falsifies.
