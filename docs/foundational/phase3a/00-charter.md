# Phase 3A: interrupted delegation and result custody

## Question and intended result

Does preservation of each agent's committed state, together with B applying
one original delegation at most once, imply that A can recover its actual
result? Try to falsify that implication at the result-acknowledgment boundary.

The intended deliverable is a minimal trace **in a declared transition system**
and a small conditional custody argument. Counts establish only exhaustive
reachability in that finite model. The existing ALAN runtime remains unchanged.

## Approaches and chosen scope

1. Extend the single-host Store with a cross-store transport. This would mix
   protocol research with implementation and authority decisions too early.
2. Run an independent finite interaction model, then identify implementation
   correspondence obligations. This is the selected first slice: it exposes
   the smallest missing edge contract without assuming a distributed owner.
3. Begin with arbitrary mutual obligations. This adds wait cycles before the
   one-direction recovery boundary is understood; it remains the next slice.

The user supplied the research program and authorized starting it. The concrete
custody experiment and finite bounds are implementation choices made here.

## Local contracts under attack

A preserves committed delegation identity, any durably recorded result, and
its once-only consumption marker across its one modeled crash. It never consumes
a fabricated result or changes authority. Volatile receipt bytes are not
promised durable. B applies the fixed authorized request once, atomically
retains its result, and keeps a deduplication tombstone even after release.
B's weak contract allows releasing the result on an acknowledgment of delivery.
Neither local contract alone promises continuous durable custody across the edge.

The composite obligation is stronger: after B has acted, A can recover and
consume the authentic original result without repeating the mutation. Safety
does not itself imply progress; useful recovery needs recoverable evidence and
service. Unknown is an honest outcome, but does not meet this useful-recovery
obligation.

## Scope and prior evidence

- Exactly two logical agents, one stable delegation, one immutable request,
  one application, two possible opaque authentic result values, at most one
  A-only interruption, and separate local durable stores in the abstract model.
- B's application and result recording are atomic **by assumption**. The model
  does not solve the external effect/receipt gap investigated in Phase 2.
- A fixed valid grant authorizes B's one operation. Message possession conveys
  no new capabilities. Dynamic delegation, overlap and cancellation remain 3B.
- Honest authenticated transport may delay and duplicate; no forgery, store
  rollback, permanent storage failure, receipt lookup or receipt regeneration.
- No claim of current Store refinement, process-crash qualification, exact-cost
  equivalence, general distributed recovery or arbitrary-size composition.

Prior sources: [Phase 1K communication](../../../research/communication/README.md),
[Phase 2B recovery](../phase2b/04-recovery-and-world.md),
[Phase 2D conclusion](../phase2d/09-decisions-and-next.md), and the manuscript's
Phase 2 recovery discussion.
The new boundary is transfer of result custody between independently retained
agents, not rediscovery of the existing source-lifetime or external-receipt gap.
