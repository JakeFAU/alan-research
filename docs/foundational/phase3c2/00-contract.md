# Phase 3C.2: Shared Capacity and Stable Reservation Contracts

Frozen research contract, 2026-11-14. Base: merged Phase 3C.1,
`69534188f796efab2e866750f6a76e2c7a3a1706`. This is a research-only extension;
the existing runtime and historical model/evidence remain unchanged.

## Question and smallest boundary

Can two exactly authorized services preserve capacity conservation yet lose
finite completion when they share a reusable slot? What information and
reservation discipline must a composition contract retain to preserve that
completion question across transport, recovery and regrouping?

The fixed service topology is A→B→C, with originals AB and BC. All eight
fixed serial/servicing policy masks are retained. Permission provenance,
cumulative expenditure E, and reusable capacity Q are distinct. The original
Phase 3C.1 full bindings, distinct effect resources and permanent covering
authority remain unchanged; E and other resources are sufficient. This does
not introduce a new grant protocol or imply that sharing a permission ancestor
shares a budget. Phase 3B.3 charges the serving principal's actual resource
root (`../phase3b3/02-contract.md`). Phase 1G already excludes waiting logical
parents from physical slots (`../phase1g/01-composition-contract.md`). The
adverse candidate here is not an allegation of a runtime defect or a newly
discovered general deadlock phenomenon.

## Finite domain and state

Use the unchanged Phase 3C.1 chain and transport relation, with fixed
generation-one A/B/C, AB/BC exact bindings, U/W/P/D stages, both opaque result
values, and one recoverable interruption per physical owner. Add a separate
capacity configuration:

* `shared-1`: AB and BC reference the single root `Q`, limit 1.
* `shared-2`: AB and BC reference the single root `Q`, limit 2.
* `separate-1`: AB references `Q-AB`, BC references `Q-BC`, each limit 1.

Roots have immutable identity and limit. Original-to-root membership is exact
and retained. Static components carry references to these roots; joining two
references does not allocate a second pool. Atomic updates to one trusted
authoritative root record are a premise, not a distributed allocator result.

The abstract state is the old named stage vector plus a Boolean held bit for
each declared original. A held bit means that exact original owns one unit of
its declared Q. Only W may hold; occupancy is derived by membership, never
stored as a substitute for holder identity. Every root's occupancy is at most
its limit. Concrete state adds the same bits to the complete unchanged
transport state. Recovery retains bits and original identity, while the old
availability/activation and volatile-result rules remain in force.

## Transitions and observations

Commit is unchanged: U→W, without needing a slot. Consumption is unchanged:
P→D using the actual durable result, without needing a slot.

`reserve-capacity` requires W, no existing hold, free capacity at the exact
root, and the selected policy below. Concretely it additionally requires a
genuinely admitted request and an up callee. It binds the hold to the original,
not a physical activation. Application requires that hold and the unchanged
serial/servicing service guard, applies once, retains its actual result and
**releases Q atomically in that same transition**. There is no early release,
timeout, preemption, refund of E, or release merely because of interruption.
No new separate release event is introduced. Existing receipt-custody `release`
retains its old meaning and is unrelated to Q.

Three reservation policies are compared:

* `early`: any admitted W original may reserve.
* `ready-now`: additionally require the old service guard at reservation time.
* `stable-ready`: require that the callee is permanently servicing, or that
  every outgoing original of that callee in this fixed chain is already D
  (vacuously true for C). Topology and policies cannot change and D is absorbing.

All 3 layouts × 3 policies × 8 masks form the finite roster. Counts of explored
states/arcs and representative probes must be measured, not predicted. The
resource-aware abstraction observes named commit, reserve-capacity, apply and
consume transitions. The old external progress observation erases only reserve
events. Completion is **both originals D**, not quiescence or empty channels.
`CanComplete` means existence of a finite legal completion suffix; no fairness
or inevitable completion is asserted.

## Required attacks and controls

1. For `early/shared-1`, B serial: construct two genuine concrete prefixes with
   both originals W/admitted, equal full base transport/physical state and equal
   pool totals, but different holders. Holder AB prevents completion; holder BC
   permits completion. Establish equality by comparing captured states, not by
   rewriting evidence. Conservation, authority and custody hold in both.
2. For `ready-now/shared-1`, reserve AB while BC=U, then commit/admit BC. The
   guard was true at acquisition and is false afterwards; preserve the legal
   noncompletable prefix and a matched BC-holder counterpart where possible.
3. Preserve a legal old/early trace applying AB before BC is committed when B
   is serial. Show why `stable-ready` cannot lift that old action without first
   advancing BC. The stronger policy restricts old traces; do not claim old-q
   backward lifting, unrestricted trace equivalence or a free repair.
4. Shared limit two, separate limit-one roots, and permanently servicing B are
   controls. All stable-ready reachable abstract states must have a completion
   suffix; a failure is evidence, not a reason to weaken the test.
5. Concrete recovery/custody witnesses preserve exact holder identity, both
   opaque results, duplicate messages, ACK discipline, independent physical
   availability and genuine durable receipts. Guards/lifting must never inspect
   WorldResult; no fabricated initial receipt or hidden post-state repair.

## Positive targets and composition boundary

Provide a source argument, separately labeled from finite checks, that every
concrete step projects to a resource-aware step or stutter, and each enabled
resource-aware abstract action has a finite lift from **every reachable
concrete representative**. Reserve may prepare actual request custody and
recover owners; preparation changes neither another original's stage nor any
hold. Apply requires an already admitted held original. Consume follows actual
retained custody independently of capacity. This supports equivalence of
concrete and resource-aware abstract CanComplete under the frozen premises.
Finite probes supplement this argument; they do not exhaust the concrete product.

For stable-ready, prove the held-work readiness invariant by induction. Finish
any existing holder, then finish BC and AB in dependency order, recovering
owners and routing real receipts as needed. Every reachable representative
has such a finite suffix under the atomic application/retention/release and
recoverability premises. This is a fixed-chain conditional source proof, not
mechanized proof, a runtime qualification or a theorem for arbitrary graphs.

Reuse the base owner/port/original component discipline. Add immutable root
descriptors, exact original memberships and the policy. Join identifies equal
root references, rejects same-ID conflicting limits, duplicate originals or
owners, missing/extra memberships and incompatible policies. Close accepts
exactly the declared chain/layout/mask/policy. Both parenthesizations must yield
the same retained configuration and the same resource-aware relation. Dynamic
holds live in one shared relation; arbitrary independently generated histories
cannot be unioned just because each local occupancy is legal. No open-component
substitution theorem or smaller constant-size interface is asserted.

## Evidence and exclusions

Export a versioned complete finite graph roster, all arcs, completion labels,
both regroupings, genuine concrete witnesses and precise representative-probe
coverage. Freeze schema before the independent checker. A separate Python
standard-library checker must regenerate the abstract relation and reverse
completion classification, replay concrete witnesses from the true initial
state with declarative rules, and check regrouping/root identity. It may share
the frozen specification but may not call/import operational Go guards,
successors, projection, lifting, graph explorer or runtime Replay.

Retain hostile-input rejection evidence, including a removed non-tree arc,
wrong CanComplete, aggregate-only/changed holder, invented root membership,
conflicting same-ID limits, copied shared pools, capacity overflow, duplicate
reservation, early release, wrong actor/activation, forged binding/receipt and
post-state invention. Reject malformed schema, missing roster members and
Boolean-as-integer counts. Hostile records are forged controls, not operational
histories. Distinguish graph totals, probe totals and witness counts.

Exclude dynamic membership/generations/policies, cycles, arbitrary topology/N,
revocation/cancellation, permanent owner loss, split unknown effects, distributed
arbitration, fairness transfer and resource-cost equivalence. Preserve the
earlier packet's claims under its own assumptions. Update both manuscript forms,
roadmap, references if changed, and rebuilt PDF with exact qualification scope.
New tests/race/scoped Staticcheck and independent audit are fresh. Full unchanged
repository tests may be attributed to retained Phase 3C.1 evidence only after
verifying exact source/dependency identity; report that reuse explicitly.
