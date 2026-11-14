# Phase 3C.1: named progress contracts and finite completion

## Status, domain, and observation

This is a direct source-level mathematical argument for the model in
`research/contractcomposition`, supplemented by complete small abstract checks,
factored concrete probes, and concrete histories. It is not a mechanized proof,
a full concrete-product exploration, or a runtime refinement. The independent
export audit and manuscript reconciliation are separate tasks.

The domain is exactly the fixed three-owner chain A→B→C and cycle A→B→C→A,
with originals AB, BC, and (cycle only) CA. Every owner has generation 1 and at
most one outgoing original; every original has its immutable full binding,
fixed grant, operation 1, and distinct sufficient resource. The eight policy
masks assign serial/servicing to A/B/C separately. Even C's redundant chain
policy remains represented. Each physical owner can crash once and recover;
crashing erases only its volatile caller receipt. Atomic application retains a
callee receipt and permanent application tombstone. Channels and durable data
survive. These are premises, not conclusions about production ALAN.

An observation is the sequence of `(original, full binding, commit/apply/consume)`
events. It abstracts result values, transport, interruption, recovery, recording,
ACK and release to stutters. Internal applications are still visible world
applications. Complete means **every declared original is consumed**, including
latent internal originals, regardless of channel drain, retained callee copies,
or physical availability. R1 separately tests *immediate concrete application
readiness*; the positive progress relation does not claim to preserve that
readiness from stages alone.

For a concrete state s, q(s) gives U if uncommitted, W if committed/unapplied,
P if applied/unconsumed, and D if consumed. Causal invariants make these cases
exhaustive on reachable states. Waiting(i) is exactly W or P on i's outgoing
original; an actual owner without an outgoing original is not waiting. An
unresolved port is never treated as an idle owner.

## Concrete construction and invariant argument

`ConcreteSuccessors` calls the unchanged durable-ACK `twoagent.Successors` after
renaming one original to the old A/B binding. It discards that edge copy's
physical crash/recovery events, filters each remaining event by its actual
actor's availability, and filters Apply by the actual callee's named policy and
outgoing stage. It then supplies one crash/recovery transition per physical
owner. This is shared implementation, not an independent transport oracle.
The old `twoagent.CheckTrace` causal ledger is separately implemented and is
run on each renamed original of every retained concrete history.

The inductive invariants used below are:

1. Binding, commitment, admission, application tombstone, recording and
   consumption are immutable or monotone as appropriate. Admission requires a
   committed request; application requires admission; consumption requires a
   durable authentic recorded result.
2. Apply atomically sets the application tombstone, retained result and
   inaccessible world fact to the same actual value. No later Apply is enabled
   for that original. Copies are sourced only from retained receipt, channel,
   volatile receipt, or durable receipt; no operation obtains bytes from the
   world fact.
3. If applied, either a callee receipt is retained or the caller has durably
   recorded it. The only deletion of a retained result is Release, which needs
   an authentic identity-only ACK. That ACK originates after durable recording,
   and recording survives crashes. Repeated Release after duplicate ACKs is
   harmless to this invariant.
4. A down owner has already used its one crash and has no volatile caller
   receipt. Recovery changes activation 0→1 and availability, not generation,
   binding, policy, durable data, or progress. There is a distinct availability
   record for A, B and C.

These follow by induction over the following transition classes. Constructor
and retained-record validation check the representational premises; validation
is not a decision procedure for global reachability. The executable lifting
algorithm does not consult `WorldResult`, including during source validation.
Authenticity relative to the actual historical result is checked by the causal
ledger, rather than reconstructed from an oracle field.

| Concrete class | Progress projection and invariant preservation |
| --- | --- |
| Commit | U→W on that named original; caller up, no prior commitment. |
| SendRequest | Stutter; sources an immutable committed request, including retries after admission/application/consumption. |
| Admit / Deduplicate, both delivery variants | Stutter; actual callee up, authentic queued request; only admission fact grows. Retaining a transport copy is allowed. |
| Apply, either opaque value | W→P; callee up/admitted and exact serial guard holds; tombstone/result retention are atomic, so custody starts and application cannot repeat. |
| SendResult | Stutter; actual callee up and retained receipt supplies bytes. It does not inspect caller volatile/recorded state to suppress retries. |
| ReceiveResult, both delivery variants | Stutter; caller up and authentic result channel supplies volatile bytes; the duplicate variant retains the channel. |
| Record | Stutter; caller up, genuine volatile bytes become durable; no record replacement. |
| Consume | P→D; caller up and durable recorded receipt supplies the result; no repeated consumption. |
| SendACK | Stutter; caller up and durable record already exists; ACK carries identity only. |
| Release, both delivery variants | Stutter; callee up and authentic queued ACK justifies deletion; caller durable custody survives; tombstone remains. |
| Crash | Stutter; affects all roles of that owner and clears only its outgoing volatile receipt. Durable custody and all progress registers survive. |
| Recover | Stutter; restores the owner's roles with the next activation; no logical identity or stage changes. |

On each original, a physical crash/recovery of its caller is exactly the old
edge crash/recovery; a callee interruption is a stutter in that old edge
projection. Filtering the edge relation removes transitions, never creates
unauthorized transitions. This establishes the stated safety induction for the
adapter without assuming the full product has been enumerated.

## Universal representative lifting

**Proposition (conditional, fixed fragment).** For every reachable concrete
state s and every abstract action q(s)→r, there exists a finite legal concrete
path from s with progress stutters before one final action of exactly the same
name/binding, yielding r. No other original changes progress; no effect is
reapplied. This quantifies over every reachable representative, not just a
chosen canonical transport state.

Choose a suffix with no further crashes. Any down endpoint can recover. Recovery
clears no durable fact and cannot change Waiting. Existing channels and
custody are used; no global drain or retry suppression is required.

* **U→W (commit e).** Recover the caller if down, then Commit e. Its commitment
  is absent by U, and no serial service gate restricts caller commitment.
* **W→P (apply e).** Recover the caller and callee as necessary. If already
  admitted, retain that fact. Otherwise, use the queued request if present, or
  SendRequest from the durable commitment and then Admit without retaining a
  duplicate. The callee's outgoing progress did not change during these steps,
  so the abstractly enabled serial guard still holds. Application is absent by
  W. Apply once with either allowed result, retaining it atomically. There are
  at most two recoveries, a send, an admission, and one application. Admission
  and availability are additional concrete conditions: abstract eligibility is
  **not** immediate concrete readiness.
* **P→D (consume e).** Recover the caller if down. If already durably recorded,
  Consume directly. Otherwise, Record a present volatile receipt. If neither
  record nor volatile receipt exists, receive the queued result if present.
  If the result channel is empty, custody invariant 3 gives a retained callee
  result; recover the callee if down, SendResult from that actual retained copy,
  receive it without keeping a duplicate, Record, and Consume. Every preparatory
  step is serviceable while an owner waits outgoing. There are at most two
  recoveries, a send, a receive, a record, and one consumption. Apply is never
  used in this case; neither a missing record nor a prior release licenses
  reapplication.

Recovery preserves the other originals' retained records. Volatile loss belongs
to the earlier crash, not the recovery or the constructed suffix. The
construction never requires another original to commit, apply or consume.
Sharing an actor across its incoming and outgoing roles does not invalidate
these stutters;
separate availability records and intact durable custody are essential.

Conversely, the table gives the projection of **every** reachable concrete
step to its named abstract step or a stutter. Removing stutters from a finite
concrete completion suffix yields an abstract completion path. Repeatedly
lifting each action in a finite abstract completion path, starting at the
actual representative reached by the previous lift, yields a finite concrete
completion suffix. Since Complete is the same conjunction of all named D
registers, this proves, under the premises:

`CanComplete_concrete(s) iff CanComplete_abstract(q(s))`.

This is existential finite reachability. An execution may keep retrying forever;
there is no fairness transfer, inevitable progress, divergence equivalence,
matching crash schedule, arbitrary-value trace equivalence, or global-quiescence
claim. The retained retry witness has an exact concrete two-step cycle.

## Compatible grouping and the dependency certificate

The allowed context grammar consists only of the three named leaves A, B, C,
compatible binary joins of disjoint owner sets, and closure to one of the two
fixed topologies. Leaves own their outgoing original and expose callee ports
with exact generation and service-policy descriptors. Join checks unique
owners/originals/outgoing incidence, exact immutable bindings, disjoint
resources, and agreement of endpoint descriptors. Closure rejects unresolved
ports and any owner/topology outside the fixed grammar. The exported intermediate
AB and BC components remain open where appropriate. No operation silently adds
an absent owner or collapses two actors into one availability/policy record.

Join is canonical compatible union. It retains every original even when its
endpoints become internal; it does not hide commit/application/consumption
labels. Thus `(A join B) join C` and `A join (B join C)` have the same named
registers, policy incidence, bindings and ports. The concrete transition filters,
projection, custody invariant and lifting cases depend only on those retained
facts, not the syntax of grouping. They therefore continue to apply after
compatible joins and closure. Canonical union equality by itself would not
establish the transport abstraction; the preceding universal argument is the
substantive additional obligation.

A dependency e→f means e's callee is serial and f is that callee's outgoing
original. An acyclic dependency certificate orders prerequisites before
dependents. To complete from any reachable state, use this order and finish each
original: commit if U, apply if W, consume if P, leave D paid. When applying,
every serial outgoing prerequisite is already D, so its gate is open. The
lifting proposition supplies the concrete plumbing. This proves acyclicity is
sufficient for all reachable states; it is not necessary for an individual
state. In particular an all-serial cycle's initial state can complete with an
appropriate ordering even though its static certificate is cyclic.

Every chain mask has an acyclic certificate. In the cycle, any servicing owner
removes the dependency edge through that owner. The all-serial state WWW is
trapped in progress and cannot complete. Pairwise conditional promises cannot
be discharged circularly. This is an assumption-compatibility control with
intact identity/authority/custody, not a custody failure.

## Retained counterexamples and limits of the probes

R0 compares an all-servicing chain after BC has completed with AB still W,
against its continuation after AB also completes. The complete BC transport
record and all physical availability records are equal. Both states can
complete, but Complete and AB's visible application/consumption history differ.
Erasing latent debt therefore loses the declared observations without any
circular wait or failed premise.

R1 compares an A-serial, B/C-servicing cycle with CA admitted and A up, keeping
all BC/CA transport fields equal. AB is U versus W, so both have unpaid AB debt,
but CA's immediate Apply is enabled only in the U case. A hidden-debt Boolean
or a component policy Boolean does not determine that declared readiness
observation. This pair alone does not refute a weaker relation omitting
readiness. Named AB stage and A's policy retain the missing fact; actual
readiness still also requires admission and availability.

The histories are deterministic constructive examples, not assertions of
shortest transport prefixes or minimal topology. The fixed grammar contains
only these three-owner examples; no larger minimality domain was searched.

The factored test takes every one of the **311 actually enumerated old
single-edge durable-ACK states**, including both results and caller phases.
For each of AB/BC/CA in the cycle, it crosses these with the callee's four
outgoing stages, two service policies, and three physical availability phases.
It also checks terminal chain BC with no outgoing C original, both C policies,
and three C phases. The other edge states are canonical durable representatives.
The executed test reports **24,258 factored representatives, 184,256 concrete
successor arcs, and 14,058 enabled named-edge lifts**. These include repeated
or overlapping concrete contexts; they are not unique reachable-product counts.
The 7,464-representative earlier log is preserved as an intermediate test run,
not the final coverage total. `WorldResult` erasure checks successor
noninterference; a separate test checks the lifting algorithm's noninterference.

Finite checks support the source argument but do not replace its universal
representative clause. Permanent owner loss, split application/retention,
deleted durable custody, forged bindings or depleted shared resources violate
its premises. Dynamic identities/grants, cancellation, arbitrary topologies and
program contexts, shared-authority fan-out, runtime recovery refinement and
reduced verification-cost claims remain unestablished. Phase 3C remains open.
