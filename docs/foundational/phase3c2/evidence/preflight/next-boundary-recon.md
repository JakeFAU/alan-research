# Next-boundary reconnaissance: Phase 3C.2

Status: read-only source reconnaissance and proposed argument; no new model, test, enumeration, runtime campaign, branch/index operation, SQLite open, or `.envrc` access. No primary-literature claim or novelty claim is made. Counts below are attributed to the existing packet, not rerun. The requested merged-base identity was supplied by the parent; this report does not independently verify GitHub delivery.

## Recommendation

Choose **shared authority with identity-indexed reservations and service interference**, with a precise decision about the kind of authority before implementation. The strongest small candidate is a reusable, finite dispatch-capacity permit shared between two exact services. It exposes a safe, acyclic service system whose completion depends on which original holds the permit. Prove/refute an abstraction and composition contract for that interference; do not make the packet a general resource scheduler.

This is a new semantic premise. Phase 3B.2's permissions are sets of complete operation/resource/scope tuples, not consumable credits or concurrent-dispatch slots. A reusable slot, a permanently consumed credit, and permission to mutate a shared counter are three different algebras. This report recommends the first because it yields a nontrivial interaction with the existing service gate while retaining atomic effects. If the intended research question specifically means permanently consumed authority, use that interpretation explicitly instead; the proposed hold/wait witness below does not establish a result for it.

Do not title the result unrestricted composition, substitution, arbitrary topology, or linear capability conservation. Treat it as the next bounded interference packet.

## What already exists

Inspected sources: `AGENTS.md`; `docs/foundational/phase3/README.md`; `docs/foundational/phase3c1/{00-contract,01-argument,04-results}.md`; `research/contractcomposition/{model,compose,concrete}.go`; `docs/foundational/phase3b2/01-contract.md`; `research/authorityservice/model.go`; and the closure references found by targeted `rg`.

Phase 3C.1 establishes a conditional source argument, supplemented by checks, for the fixed three-owner chain/cycle. Every named original retains its full immutable binding and U/W/P/D stage; every owner retains its physical availability and service policy. Concrete transport steps project to named progress actions or stutters. Every enabled abstract action has a finite concrete lift from every reachable concrete representative, without advancing another original or reapplying an effect. This yields equivalence of existential finite completion under the stated durable-custody, recoverability, atomic application/retention and disjoint-sufficient-resource premises.

The grammar is literally enforced in code: `Config.Validate` compares with `NewConfig` for A/B/C chain or cycle; fixed arrays/slots encode AB/BC/CA. `Component.Validate` requires one outgoing obligation per caller and distinct resources. `Join` is checked canonical union; `Close` rejects unresolved ports. There is no replacement operation or quantified contextual substitution relation.

The old R0 attack already rejects erasing internal debt/visible applications. R1 already rejects debt/policy bits for immediate readiness. The all-serial cycle already shows circular service promises can prevent completion despite intact custody and authority. Recreating those examples with renamed nodes adds little.

Phase 3B.2 already distinguishes consent, exact rooted grants, tuple attenuation, and invented Cartesian permissions. Its contract explicitly says that shared counter values/interference are abstracted away and that no shared-resource algebra or isolation theorem follows. Therefore set union is not an unproved consumable-authority protocol. A double-spend example refutes a newly posited copying rule, not existing Phase 3B.2.

## Comparing the three directions

| Direction | Real new obligation | Assessment |
| --- | --- | --- |
| Arbitrary finite topology lifting under unchanged premises | Generalize local transport projection/lifting and service-dependency induction to named finite incidence maps; handle many outgoing edges with an explicit waiting rule. | Useful theorem hygiene, but much is a parametric restatement of the existing local proof. A larger graph demo adds no theorem. Outdegree >1 changes waiting and must not be smuggled in as indexing work. |
| Open-component substitution | Specify ports, accepted environment actions, observations/refusals, latent debt, and a relation between genuinely different components; prove it is preserved by admitted contexts. | Substantive, but currently there are no alternative components or adequate environment semantics. Defining enough can make this a much larger semantic project. Reusing identical retained unions establishes no substitution. |
| Shared authority interference | Specify the authority algebra and shared-state ownership; show why independently valid contracts do or do not compose when authority affects progress. | Best small next attack, directly addresses the named roadmap gap, and can preserve the existing custody machinery. Needs a counterexample beyond duplicating one credit. |

General finite-graph lifting can remain a subsequent theorem packet after the interaction contract is known. Open substitution should be attempted once there is an actual pair of different implementations/contracts whose equivalence is worth preserving.

## Why the obvious overspend witness is insufficient

A one-credit root copied to B and C followed by two irreversible spends immediately violates conservation. It is a useful negative control and distinguishes set membership from a resource algebra, but merely discovering that 1+1 exceeds 1 does not answer the difficult composition question. A positive construction that assigns nonoverlapping credits to fixed requests and checks their sum is also largely arithmetic unless it confronts admission, reservation lifetime, recovery, or completion.

A more useful witness keeps conservation, consent, exact identity, and durable custody valid throughout, but makes an aggregate authority summary hide the fact that changes completion. This separates authority safety from compositional progress.

## Proposed minimal hard witness

Use three owners A, B, C and two originals e=A→B and f=B→C (the old chain). B is serial; C services incoming work. The authority topology is a fan-out: one root R at A issues two fixed, exact, non-redelegable grants, to B for e and C for f, sharing one dispatch-capacity token t. The request topology and authority topology are deliberately different. Issuer A's approval of B's request to C must be explicitly part of this new grant contract; it is not silently inherited from the old requester-only grant model.

Alternative if a literal fork in the request topology is mandatory: add A→C as a separately named original, retaining its debt, and explain why it is a fixture requirement. It is not needed for the two-service interference witness and must not be presented as globally minimal.

Token semantics:

- Exactly one token exists. Free or held by exactly one fully bound original; reserve and release update that authoritative durable record atomically.
- Reserve requires a genuine admitted request, valid exact rooted grant, and free token. It may occur while the callee is serially waiting; this is the candidate policy under attack.
- Apply requires that original's token plus the old service guard, and atomically applies, retains the result/tombstone, and returns the token. No split external-effect window is introduced.
- Crash/recovery preserve the original-bound reservation, custody and binding. There is no automatic release, timeout, revocation, preemption, cancellation, transfer or minting. A recovered server may continue an unapplied original using its durable reservation under the same atomic-effect premise. State this explicitly rather than reusing the old volatile-admission rule.
- Receipt transport, durable recording, consumption, ACK and release remain serviceable while waiting and do not require this token. Physical availability stays separate for A/B/C.

After committing and admitting e and f, both are W. Two legal prefixes differ only in which service reserves t:

1. **t held by e:** e cannot apply because B is waiting on f. f cannot reserve because e holds t. No result exists to consume. Transport retry and recoverable interruption do not release t. No completion suffix exists.
2. **t held by f:** C can apply f, returning t; B can consume f; B's service gate opens; e reserves/applies/consumes. A completion suffix exists.

The shared root/grant descriptors, named progress vector WW, physical availability, and aggregate accounting (`free=0`, `held=1`, `completed-applications=0`) agree. Only the reservation-to-original association differs. Each history obeys conservation and consent. The static service dependency e→f is acyclic in both states, yet it is insufficient after resource sharing.

This directly falsifies an aggregate interface for **CanComplete**, stronger than an immediate-readiness-only collision. The missing datum is an authority reservation dependency. In the bad state, e waits for f's service progress and f waits for e's authority release. The issue is the interaction of two kinds of dependency, not an overspend.

These are proposed finite paths and a source-level reasoning sketch. They have not been executed or independently checked. Minimality is only that two pending services and one shared token suffice for this declared construction; no search over all authority protocols/topologies was performed.

## Conditional positive target

Primary target: an authority-aware abstraction that represents each original's full binding/stage and the exact shared token holder, together with the root/grant facts and per-owner policies. Include explicit reserve actions in the abstract relation; they are silent with respect to the old commit/apply/consume observation but change retained abstract state. Application returns the token as specified.

Prove (or falsify) that every concrete action projects to the declared abstract relation and that each enabled abstract action can be lifted from every reachable representative with only appropriate transport/recovery preparation. The result sought is:

`CanComplete_concrete(s) iff CanComplete_augmented(q(s))`.

It must correctly classify the safe noncompletable state; a positive abstraction result does not require engineering it away. Universal representative quantification matters: canonical states alone will miss which request currently owns the reservation. The proof must also check that transport preparation cannot change another original's stage or acquire/transfer the token silently in a way omitted by the abstraction.

Composition must identify one authoritative root/token record, combine nonduplicated exact grant/claim records, and preserve holder-to-original incidence even when endpoints become internal. An open component cannot independently promise use of a token already owned by its context. Either carry an explicit root/authority port and environment transitions, or restrict joining to components whose shared-authority synchronization is specified. Simply unioning duplicate local copies of the root record is not composition.

Optional, separately labeled progress theorem: under a stronger admission discipline, reserve e only after B's own prerequisite f is D; reserve f whenever admitted and the token is free. Because B directly retains consumption of f, this guard need not read an inaccessible global oracle. Then complete f, release the token, complete e. This narrows allowed traces when f is U; do not claim it preserves the old relation or is the unique/best fix. Do not add this policy unless necessary to answer the research question. The corrected abstraction plus preserved negative witness is already a meaningful result.

An alternative guard checking merely `!Waiting(B)` at reserve is insufficient: B can reserve e while f is U and subsequently commit f, forming the same obstruction. Any positive policy needs stability of its readiness premise (e.g. prerequisite already D or an explicit commitment restriction), not a one-time local readiness check.

## Concrete packet acceptance boundaries

1. Freeze token algebra, authority/request topology, identities, reservation/apply/release boundaries, observations, and composition ownership before implementation.
2. Retain the safe equal-summary/different-CanComplete pair and its exact causal histories; retain the simple overspend as a separately labeled mutant, not the central finding.
3. Independently reconstruct authority conservation and holder identity from event histories; do not invoke the operational grant/transition guards as the oracle.
4. Check all small abstract states under the frozen initial-state reachability relation and representative transport lifting cases. State whether totals are unique states, profile sums or factored probes.
5. Preserve retry divergence, opaque result alternatives and physical-owner crash correlation. Permit no world-result oracle or hidden snapshot repair.
6. A sufficient condition for safe composition is not a proof of inevitable progress. No fairness transfer, arbitrary contexts, arbitrary graphs, real distributed arbitration, dynamic grants, revocation, cancellation or split effects is claimed.
7. Reconcile manuscript/roadmap with the exact result and explain that reusable dispatch capacity is an additional authority algebra. Keep historical tuple-attenuation results distinct.

## Risks / reasons to reject or narrow this recommendation

- If the parent interprets the user's intended fan-out solely as immutable set permissions, this experiment changes the question. Name that choice visibly; do not let the word authority obscure it.
- A single authoritative token state is a semantic primitive, not a distributed implementation. It does not establish consensus, replicated-ledger correctness or failure availability.
- Crash releasing a permit changes the counterexample: an available crash/recovery might restore a completion suffix. Durable reservation preservation is therefore a decisive premise, not incidental bookkeeping.
- Permanently consumed credit semantics would make a one-credit/two-mandatory-spend system noncompletable regardless of holder; that would destroy the stronger witness.
- Releasing on receipt consumption rather than Apply produces a different protocol and requires fresh analysis. Do not silently alter token lifetime to make results pass.
- The mixed dependency example is elementary operational reasoning and should not be marketed as a new general deadlock theorem. Its research value is a precise ALAN composition boundary, a retained counterexample and a sound contract that addresses it.
- No reduction in verification cost follows from this small packet.

## Parent-requested refinement: make the positive completion rule a packet target

The positive candidate should be explicit, alongside the negative aggregate-summary result. In the fixed e=A→B, f=B→C chain, replace unrestricted reservation by:

`ReserveAllowed(e) := admitted(e) && exact_authority(e) && free(Q) && (servicing(callee(e)) || no_outgoing(callee(e)) || all_outgoing_stages(callee(e)) == D)`.

All-outgoing-D is vacuously true for a terminal owner; retaining the cases separately is useful for explaining the premise. Here B is serial and has only f, so B may reserve for e only after f=D. C may reserve for f immediately after admission. Existing Apply still checks its original serial service guard. The additional reserve predicate is monotone in the frozen model: stages never leave D, identities and policies do not change, and no new originals are created. Recovery preserves it. Admission may be prepared without acquiring Q. Thus the rule does not acquire a scarce slot while waiting for a prerequisite to become ready.

Two distinct targets now follow:

1. Universal representative projection/lifting and equality of concrete/augmented-abstract finite-completion reachability for both candidate policies. The permissive policy retains the safe noncompletable witness.
2. For every reachable state of the **monotone-reservation policy** in this exact fixed fragment, a finite concrete completion suffix exists under the custody/recovery premises. This is an existential suffix theorem, not inevitable progress.

Proof skeleton for target 2: choose no further crashes and recover any down endpoints. If f is not D, the reservation invariant means e cannot hold Q. If f is U, commit/admit f; if W, acquire the free token or use f's existing reservation, apply once at servicing C and free Q; if P, route/record the authentic retained result and consume f. Its service/custody steps never require advancing e. Once f=D, Q cannot be held by f because f returned it at Apply and may not reserve again; e is either already D, P with a consumable authentic result, or U/W with Q free or already held by e. Commit/admit/acquire as needed, apply once at B and consume e. Down owners can recover without losing the durable reservation or custody. This covers every reachable representative, not merely a fresh run. Each phase is finite by the old concrete lifting constructions plus a single bounded acquire step. Independent checking must verify all premises used in this argument.

The one-time old service guard is insufficient even if Reserve is forbidden when B is already waiting: reserve e while f=U, then B commits f. The reservation now belongs to a blocked service. This is a second adverse schedule to retain. The alternative discipline forbidding B from committing f while holding e's reservation can also avoid this witness, but changes another observable boundary and need not be implemented in the first packet. Freeze the monotone-reservation rule as the chosen conditional positive candidate, while retaining the alternatives as explicitly unqualified.

### Keep three ledgers/meanings distinct

- R is rooted permission provenance/consent and exact tuple coverage.
- Q is a newly specified shared finite reusable dispatch-capacity account. Both grants explicitly name Q; matching permission ancestry alone does not imply a common capacity account.
- Existing consumable execution budgets E are a separate algebra and are not recycled by releasing Q.

Inspection of `docs/foundational/phase3b3/04-results.md` around its resources/debit discussion records that delegated service debits the serving root, not the issuer's resources. Therefore this packet must not claim that an existing shared permission root already creates a shared budget. The proposed Q is new research semantics, not a runtime behavior discovered in Phase 3B.3. A single authoritative Q update is assumed; no distributed capacity implementation is supplied.
