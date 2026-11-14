# Phase 3C.2 conditional source argument

This is a source-level mathematical argument for the frozen fixed-chain
relation in `00-contract.md`, supplemented by finite checks. It is not a
mechanized proof, runtime qualification, independent checker result, or theorem
for arbitrary graphs. The new package leaves `contractcomposition` and
`twoagent` unchanged. Fresh independent validation and manuscript reconciliation
belong to subsequent tasks; this document does not mark them complete.

## Premises and observations

A, B and C have fixed generation-one identities, AB and BC have the exact old
request/operation/grant/effect-resource bindings, and each principal's servicing
policy is immutable. Admission and application use the old trusted dispatcher.
Each physical owner can interrupt once and recover; durable state survives.
An application atomically retains its actual opaque result and releases its
original-bound Q unit. Root descriptors and memberships are immutable, updates
to the shared root relation are atomic, and no permanent owner loss, cancellation,
revocation, unknown split effect, or topology mutation occurs. These premises
are not an implementation of a distributed allocator.

Permission is exact request-bound authority; Q is reusable occupancy. The old
counter effect resources remain distinct even when Q is shared. Cumulative E
and other resources are assumed sufficient as in the frozen contract; this
model introduces no E balance or refund transition. Application is still at
most once per original. Therefore releasing Q does not mean undoing expenditure
or granting authority. No claim about resource-cost equivalence follows.

The resource-aware abstraction retains the full named old stage vector and the
two held bits. Its observations are named commit, reserve-capacity, apply and
consume, carrying the original's complete binding. The old external progress
observation deletes reserve-capacity only. The old receipt `release` is a
transport stutter and is unrelated to Q. Completion means AB=D and BC=D.
CanComplete means existence of a finite legal completion suffix, with no
fairness or inevitable completion assertion.

## Conservation and reachable holder causality

Initially no original holds. Only reserve-capacity sets a bit. Its source must
be W, unheld, assigned to the exact root with occupancy strictly below limit;
it increments that root's derived occupancy by one. Apply requires that bit,
changes W to P and clears the bit atomically with result retention. All other
steps retain holds. Thus occupancy never exceeds limit, holds exist only at W,
and no crash, ACK release or consumption silently releases capacity.

Concretely reservation additionally requires a genuinely admitted request and
an up callee. Admission is monotone in the unchanged transport relation. A
reachable held original consequently remains admitted until application. State
validation checks this causal premise, but does not establish global
reachability or receipt authenticity. Authenticity is a history property: the
exported witnesses begin with empty custody and every result comes from a real
application transition. Factored probes are explicitly synthesized and do not
establish such whole-chain histories.

## Why totals lose the completion question

With B serial, early/shared-1 permits genuine prefixes committing, transmitting
and admitting AB and BC, then reserving AB or BC. Both final base records are
identical, including all physical phases, bindings, custody and channels. Each
has occupancy one. Only exact holder identity differs.

When AB holds, its serial guard is false while BC is W. BC cannot reserve an
occupied root. Neither can consume. Transport retries and physical recovery
cannot change those facts; application is the only release. No finite
completion suffix exists. When BC holds, terminal C's guard is true. Applying
BC releases Q; recording and consuming BC requires no Q. AB can then reserve,
apply and consume. The export retains the two prefixes and a full completed
extension of the BC prefix. Equality is checked on captured states, with no
post-state rewriting.

Ready-now permits reserving AB when BC is U. The very next BC commit changes
BC to W and invalidates B's service guard. Admission is not the invalidating
step. The ready-now AB prefix and its BC-holder counterpart again reach equal
base records and equal totals with opposite CanComplete labels. Ready-now
acquisition is not a stable promise of readiness.

## Concrete/resource-aware correspondence

For every reachable concrete source, every old commit and consume projects to
the corresponding named transition. An old apply is retained only when held;
it projects to the new apply, including clearing exactly that hold. Reservation
projects to reserve-capacity. All remaining old transport, duplicate, crash and
recovery steps stutter in stages and holds. Actual result values remain in the
concrete state, never in operational guards as an oracle.

For each action enabled in the **new resource-aware abstraction**, a finite lift
exists from every reachable representative under the premises:

* Commit: recover the caller if needed, then execute its actual commit.
* Reserve: recover the callee; if not admitted, transmit the actual committed
  request (recovering the caller if needed) and admit it; then reserve. All
  preparation stutters in stages and holds, so free capacity and the selected
  policy predicate stay true.
* Apply: its hold implies admission by the causal invariant. Recover the callee
  and perform a real old application, selecting one of the old opaque results,
  retaining that result and clearing only its own hold. No caller receipt or
  WorldResult is read to choose the step.
* Consume: recover the caller. Use its durable receipt if present; otherwise
  record its volatile receipt, receive from the result channel, or recover the
  service owner and transmit its retained actual result. Durable-ACK custody
  ensures a retained source when no caller-side source exists. Record and
  consume without a capacity reservation or reapplication.

Preparation changes no other stage and no hold. This gives the exact target
projection at the final step. Project a concrete finite completion path and
delete stutters for one direction; successively lift each action of a finite
abstract completion path for the other. Under the premises, concrete and
resource-aware abstract CanComplete coincide. This does not assert strong
step equivalence, cost equivalence or divergence/fairness transfer. The probes
check legal adjacency, full projection and oracle independence on their stated
domain; they do not quantify over the whole concrete product.

## Stable-ready sufficiency and its cost

Inductive invariant: each reachable held original has a service guard that is
true and remains true until it applies. Reservation establishes either an
immutable servicing callee or all that callee's outgoing originals D. D is
absorbing. Another reservation, transport step, interruption or recovery changes
neither fact. Application removes its own hold. A commit cannot invalidate a
serial AB holder because BC was already D; C has no outgoing original.

From any reachable state, apply existing holders after recovery, retaining
results and releasing their units. Their readiness is preserved during this
finite sequence. Consume applied work using actual durable custody without Q.
Then finish BC before AB. If BC is U, commit it; if W, reserve the now-free unit
and apply; if P, consume. C's terminal predicate is vacuous. After BC=D, serial
B satisfies stable-ready for AB. The same finite construction works for shared
limit one, shared limit two and separate limit-one roots. It is existential and
conditional, not an inevitable schedule, necessity result or optimal policy.

An early/old trace can apply AB before BC commits. Stable-ready with serial B
cannot acquire AB while BC=U. Completing BC first introduces visible BC
commit/apply/consume events, which cannot be hidden as preparation. Therefore
the stronger policy restricts old progress traces. There is no old-q universal
backward lift, unrestricted trace equivalence, or free repair. The old-trace
witness is retained precisely to expose this restriction.

## What joins and what does not

Leaves retain the old owner/port/original discipline plus layout, policy, exact
memberships and root descriptors. An owned A/B must retain its outgoing original;
C has none. A join checks both fragments, identifies equal root IDs once, and
rejects conflicting limits, duplicated ownership/originals, missing/extra
memberships, copied pools and incompatible layouts/policies. Exact expected
root records are derived from the declared layout; no numeric value is clamped.
Close requires the exact chain and configuration. Both parenthesizations retain
the same records and produce the same whole-chain relation.

This joins static declarations, not independently executed histories. Dynamic
holds inhabit one shared state/relation. Two local histories with locally legal
occupancy cannot be unioned to allocate two units of a shared limit-one root.
No open-component substitution theorem, arbitrary-topology associativity theorem,
or smaller constant-size progress interface is claimed.
