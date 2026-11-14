# A conditional continuation-interface argument

## Question and observation boundary

The old full capacity/custody snapshot does not identify the future transition
relation of B's continuation monitor. The question is whether an explicit
successful visible language supplies a sufficient replacement interface for
initial *possible* completion. This packet replaces only M: the unchanged K
still owns all A/B/C service guards, exact authority bindings, stages, holdings,
and capacity transitions. It does not factor an entire B agent out of K.

Fix a finite common alphabet Σ of the eight exact bound actions. Fix the
immutable signature (full old configuration, exact originals, root descriptors
and original memberships). All products below require exact signature equality;
a product never creates an additional pool. Factors have disjoint private state.
Each finite marked epsilon-NFA has one initial state, explicit states, marks and
edges. A visible transition synchronizes across every factor. A private tau
advances only one factor and reads/writes no other's state. Product marks are
conjunction. A finite word is accepted when some finite path ends marked after
erasing only tau, including any finite leading/trailing tau. Divergence is not
a successful finite path and receives no special acceptance.

K's atomic application remains visible under every grouping. K is marked only
when both originals are D. Every M state is marked, including the doomed PC 3:
M cannot create the negative result merely by refusing private acceptance.
Opaque result values retain the inherited abstraction; no payload equivalence
is asserted. Durable receipt ACK/release and reusable Q release remain distinct.

## Finite-path intersection, both directions

**Proposition (conditional source argument).** For compatible pointed automata
X and Y under these rules, `Lm(X || Y) = Lm(X) ∩ Lm(Y)`.

*Forward inclusion.* Take a successful finite product path with visible word w.
Project its states and transitions onto X. Each synchronized visible step is an
X transition with the same exact label; each X-private step is its own tau;
Y-private steps become omitted stutters. The projected finite path starts at X's
initial state, ends at an X mark by conjunctive product marking, and has visible
word w. The same construction for Y shows w belongs to both languages. This
argument includes the empty word and trailing private steps.

*Reverse inclusion.* Take finite accepting X and Y paths for the same visible
word `a1...an`. Decompose each path into its finite tau-only leading segment,
first visible step, intervening tau-only segments, subsequent visible steps,
and finite trailing tau segment. At the product initial state, interleave the
two leading private segments in either order. Their independence guarantees
that each stays enabled. Both first visible steps now have the same label, so
synchronize them. Repeat for each intervening private segment and visible
letter, then execute both trailing private segments. The resulting finite
product path ends in two marks and accepts the original word. The same
construction with only the leading/trailing segments covers n=0. No fairness,
termination of all executions, determinism, or bound on private path length is
needed: the witnesses are already finite. Nondeterministic path choices are
existential and can be chosen separately for this same word.

**Sufficient replacement.** Fix any compatible finite pointed context K. If
`Lm(M)=Lm(N)`, intersection gives `Lm(K||M)=Lm(K||N)`, hence equal initial
CanComplete (nonempty successful language). The quantifier is over compatible
contexts obeying the stated product assumptions. It applies also with `K||D`
as the context. For *residual* pointed components, the same result holds only
if those particular residual languages are equal. Equality of the initial
languages does not identify residual languages after arbitrary matched histories.

**Grouping.** The flat triple relation has synchronized visible steps in all
three factors, and each private step in exactly one. Flattening `((x,y),z)` and
`(x,(y,z))` to `(x,y,z)` preserves every transition in both directions and the
conjunction of all three marks. This explicit state bijection is stronger than
merely comparing final languages. There remains exactly one K, so grouping
cannot sum, clone or independently allocate Q. The implementation checks the
complete finite flat and both grouped relations after tuple normalization.

This is elementary finite automata reasoning instantiated at the declared
boundary, not a new unrestricted process-theory result. If a larger grammar
contains exact finite-word testers, an unequal-language distinguishing word
can be used to separate components by initial possible completion. The present
role-restricted A/C controllers are not asserted to realize those testers.
Necessity, full abstraction or a coarsest interface for ALAN contexts is **not**
established. No allocation-policy optimality or performance benefit follows.

## Exact language decision

For a finite NFA, epsilon closure is the finite reachability set under tau.
Start from the epsilon closures of the pointed initial states. A visible letter
maps each subset to the epsilon closure of all outgoing targets for that
letter. A subset accepts exactly when it intersects the marked set. Explore
all reachable pairs of these subsets, including empty subsets. Equal acceptance
at every reachable pair is sufficient and necessary for equal finite languages:
induction on word length gives precisely the pair reached by each word; every
explored pair has a witnessing word. A BFS first acceptance mismatch supplies
a shortest distinguishing word, including epsilon. The exporter retains every
pair and every letter arc, even after finding a mismatch, plus the actual
acceptance checks. Tests cover tau cycles, nondeterminism, empty subsets,
epsilon mismatch and a mismatch past a shallow prefix sample. Independent
reconstruction, rather than the producer's Boolean, is still required for the
qualification packet.

## ALAN instantiation and concrete premises

The old K functions are consumed unchanged. Only K mutates stages/holds;
M restricts reserve/apply AB and commit BC, with self-loops on the other five
labels. D restricts commit AB and apply BC, with self-loops on the other six.
These role restrictions prevent components from vetoing externally owned
actions. The bounded roster is precisely 72 K × 5 M plus core K × 324 D × 5 M,
not arbitrary ALAN programs or topologies.

A wrapped concrete state pairs the **entire** old concrete state with retained
M control. Every old physical/transport stutter retains M; each visible old
progress transition requires the matching M edge. An M tau changes only M and
requires physical B to be up. The abstract tau has no physical-availability
guard. This is consistent with the lifting premise because a down B can take
the inherited recovery stutter first; recovery preserves custody, abstract K,
and M control. The finite lift then executes the tau. This does not grant an
abstract tau the power to recover physical state or claim uninterrupted concrete
availability. No fairness is claimed: existential recovery schedules are enough.

For reachable concrete states, projection sends each wrapped step to the
corresponding K/M edge or a stutter. Conversely, for a visible K/M edge, the old
capacity Lift prepares actual custody with physical/transport stutters before
its requested visible action. Those preparatory steps retain M; the final step
can take the already enabled monitor edge. For a private edge, recover B if
needed as above and execute the private step. These statements rely on the old
reachable-custody/lifting premises and atomic effect/result/Q-release model;
structural Validate alone does not establish receipt authenticity or reachability.
The producer checks each actual wrapped successor of the eight captured paths.
No exhaustive concrete product, new runtime campaign, or runtime-refinement
proof is claimed. The old source-identical qualification is inherited evidence,
not newly rerun million-case probes.

## Falsified snapshot and nonblocking claims

Histories 1 and 2 execute the identical genuine prefix: commit AB, send-request
AB, admit AB, reserve AB. At the cut the entire old concrete state, configuration,
exact holds, capabilities and static rules agree. AB is W and held, BC is U and
unheld, B is up and AB's service guard is true. Both M PCs are 1, but their
future relations differ. Commit-first disables apply AB until commit BC; after
BC commitment AB is blocked by B's serial service guard, while BC cannot acquire
shared Q. No completion suffix exists in the full product graph. Release-first
can apply AB while BC is U, release Q, and complete both originals. The captured
successful suffix and complete reverse-reachability classifications retain this
counterexample without editing post-states. Shared-two, separate pools and
servicing B each supply a genuine commit-first completion control.

Release-first, release-first-tau and release-first-doomed have equal initial
successful languages. The tau variant inserts private bookkeeping between
reserve AB and apply AB; it cannot change visible word acceptance. The doomed
variant adds a branch to an all-marked state that permits only the five O loops.
Words accepted through that branch were already accepted by staying at PC 1
and taking those O loops (all M states marked). It adds no new successful words.
Yet in the core K context its branch after reserving AB prevents both apply AB
and commit BC forever, leaving debt. The full graph retains this reachable
noncompletable residual even though the initial product can complete through
the other branch. Histories 7/8 capture both outcomes; crash/recovery in the
latter preserves PC 3 and its lack of a completion suffix.

Therefore initial equal-language replacement preserves possible completion but
not every-reachable-state nonblocking, branching quality, refusals, fairness,
inevitability, divergence or exact cost. The old stable-ready every-state result
is not silently imported after M removes actions. For example, with nonservicing
B stable-ready can require BC completion before reserve AB, whereas these
restrictive M fixtures may require reserve AB before commit BC. This is a new
restriction-induced circular dependency, not a contradiction in the old K-only
qualification. Universal nonblocking after monitor restriction remains false.

## Evidence categories and remaining obligations

The intersection/substitution/grouping statements above are conditional source
arguments. Full exported graphs, subset-pair decisions and grouping comparisons
are finite model checks within the frozen roster. Histories are concrete model
captures, not endpoint/runtime records. Implementation-side adjacency and tests
share owning model code; they are not an independent oracle. Task 2 must
reconstruct all abstract relations and these histories without calling Go
semantics and retain hostile controls. Independent review is a further challenge,
not proof. Complete-agent replacement, arbitrary topology, dynamic authority,
distributed allocation, split/unknown effects, ALAN full abstraction, universal
nonblocking and general Phase 3C closure remain outside this result.
