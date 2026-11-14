# Phase 3C.3 builder preflight

Verdict: **READY TO PLAN. No blocking contradiction or missing decision found.**

This is a read-only specification and owning-source review, not an implementation
approval, executable qualification, or proof certification. No tests, graph
enumeration, concrete captures, repository edits, or subagents were used.

## Blocking questions, ranked by impact

None. The contract fixes the semantic boundary, adverse witnesses, finite domain,
independent-checker obligations, and excluded generalizations sufficiently for a
sensible engineer to prepare a concrete implementation plan. Representation and
export-format choices can be made during planning without asking the user to
decide ordinary implementation details.

## Checks that support this verdict

### All-marked monitor language equality and nonblocking control

Contract lines 79–90 define all five graphs unambiguously, including O self-loops
in every state and disabled unlisted r/a/c moves. Marking every state makes the
monitor languages prefix closed; it does not make their synchronized products
complete, since K still requires both originals D.

The proposed three equal-language fixtures are consistent by direct finite-path
reasoning. In release-first-tau the private 1→3 step replaces the readiness to
take a with a finite silent preparation. Erasing that step gives exactly the
release-first visible paths. Marking both intermediate states preserves accepted
prefixes. In release-first-doomed every successful ordinary release-first path
remains present. A path choosing 1→3 can subsequently emit only O labels; the
same visible word was already accepted by staying at marked state 1 with its O
self-loops. Thus this extra branch adds no visible accepted word.

The doomed branch nevertheless leaves AB reserved and unapplied while forbidding
both a and c. The core K cannot finish either original from this cut. Conversely,
in the ordinary release-first core product c is disabled until AB has applied;
after a, AB's capacity is released and neither remaining consumption nor the BC
service requires AB to apply again. This supplies a sensible positive control
against which the doomed branch can demonstrate loss of every-reachable-state
nonblocking. These are source-level arguments to qualify, not results of an
executed graph search in this preflight.

The full eight-label language is essential. The monitors can accept words that
K never executes, including repeated actions after state 2. Comparing only K's
feasible traces would be a different interface. The contract explicitly requires
the full alphabet and exact epsilon-NFA language comparison, including empty
words, nondeterminism, cycles, and empty subsets (lines 159–185).

### Synchronous product and theorem scope

Lines 55–75 and 141–156 specify the premises needed for the stated language
intersection argument: a shared exact alphabet, all visible events synchronized,
independent private steps, exact signature agreement, and conjunctive marking.
Finite local paths with the same visible word can be interleaved segment by
segment, including private prefixes and suffixes. Projection supplies the other
direction. Nothing in the contract silently adds hidden shared-state tests to
monitor tau steps.

Associativity is the ordinary same-alphabet synchronized product with one K and
one Q relation. Comparing normalized reachable tuples, arcs, and marks under
both parenthesizations is a concrete finite obligation. It does not require
independently executing agent factorization, independent capacity pools, or
merging separately generated histories. Those broader interpretations are
expressly excluded.

The theorem concerns initial existential completion. Equal-language fixtures
need not have equal completion answers at each individual same-history residual:
the doomed fixture deliberately witnesses that distinction. The instruction to
compare CanComplete for equal-language pairs in both rosters is read consistently
with the explicit initial-state theorem in lines 148–149 and residual caveat in
169–170. Exporting reverse CanComplete labels for every reachable state remains
a separate required graph check. No new user decision is needed to resolve this.

### The 324 role-valid D contexts

There are four independent source/label transition choices: two states times
commit AB and apply BC. Each choice has three possibilities (absent, target 0,
target 1). Four marked subsets then give 3^4 × 4 = 324 descriptors. This is
descriptor counting, not a measured reachable-state or arc total.

The remaining six labels loop at both states, so these controllers cannot veto
B-owned visible actions. They also leave A/C's other two actions unrestricted,
which is a deliberate restriction of the grammar rather than a role mismatch.
Retaining unmarked, unreachable-state, and behaviorally duplicate descriptors is
consistent with the explicitly frozen roster. Marks can exclude completion even
where visible transitions remain possible; conjunctive marking already defines
that behavior. The product count is 324 × 5 = 1,620; the inherited configuration
roster gives 72 × 5 = 360. Actual graph totals must still be measured later.

### Matched concrete cuts and eight histories

The cut after commit AB, real request send/admission, and reserve AB is supported
by the inherited relation. AB is W/held, BC is U/unheld, and B is not waiting on
an outgoing original. The service guard is consequently true. Both monitors
permit exactly the same concrete prefix and reach monitor state 1; all inherited
concrete fields can therefore match without post-state editing. The proposed
snapshot intentionally omits the continuation transition relation.

From the commit-first cut only c can advance progress. It makes BC W and causes
serial B's service guard for AB to become false. With AB holding the sole slot,
BC cannot reserve; there is no apply/consume escape. From the release-first cut,
a is permitted while BC remains U and releases Q, after which BC can finish.
This explains why the full product graph should distinguish the cuts, while
leaving its actual reconstruction and complete reverse labeling to qualification.

The three completion controls also match the inherited guards: limit two and
separate roots permit BC to reserve while AB holds; servicing B (mask 2) permits
AB application despite BC waiting. The crash/recovery histories use only one B
interruption each. The inherited concrete relation supports recovery, and the
new contract explicitly requires private monitor state to survive it. Requiring
B up for private steps is compatible with recovering B before tau; it introduces
no need to weaken physical-step checks or reset monitor control.

Old Lift performs genuine custody preparation before the requested visible
action. The new wrapper must advance M only on that visible action, while
retaining M across old transport and physical steps. The existing APIs support
this extension without editing the old model. Actual legal event sequences,
authentic result production, equality of all cut fields, and retained doomed
control still need independent reconstruction in the eight later captures.

## Owning sources inspected

- `docs/foundational/phase3c3/00-contract.md` (entire contract).
- `docs/foundational/phase3c2/00-contract.md` (inherited finite relation and scope).
- `research/capacitycomposition/model.go`, especially `reservationReady` and
  `Successors` (lines 145 and 167 onward).
- `research/capacitycomposition/concrete.go`, especially `ConcreteSuccessors`
  and `Lift` (lines 51 and 83 onward).
- `research/contractcomposition/model.go`, especially `Waiting`, `ServiceGuard`,
  and `Successors` (service guard starts at line 146).
- `research/contractcomposition/concrete.go`, `ConcreteSuccessors` (line 163
  onward), for actor availability and single-interruption recovery behavior.

No fresh verification of the stated base commit/tree identity, runtime behavior,
historical qualification, or manuscript consistency was performed. Those are
delivery obligations, not preflight blockers exposed by this review.
