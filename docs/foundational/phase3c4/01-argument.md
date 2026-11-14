# Residual marked-language profiles and contextual nonblocking

## Domain and status

This is a conditional mathematical source argument for finite pointed marked
epsilon-NFAs, not a mechanized proof, an unrestricted ALAN congruence theorem,
or a runtime refinement result. The operational assumptions are explicit: factors
have a common visible alphabet and exact compatible signature; all visible events
synchronize all factors; a private tau step changes only one factor and cannot
inspect or alter another factor; product marking is conjunction. Acceptance means
a finite path to a mark after erasing tau, including leading and trailing tau and
the empty word. No priorities, fairness, inevitable completion, or preservation
of divergence is imposed. Finiteness is needed for the decision procedure; the
finite-path constructions below require private independence and synchronization.

For a pointed automaton X and visible word w, R_X(w) contains all states reachable
along a finite path projecting to w, including trailing private steps. L_X(x)
contains the visible words of finite paths from state x to a mark. Define
P_X(w) = {L_X(x) | x in R_X(w)}. This is a set of languages, not a multiset and
not their union. The empty profile and a profile containing the empty language
are different. Nonblocking means every reachable state has some finite marked
continuation. It says nothing about whether every maximal execution completes.

## 1. Derivative union

For every X, w, and v, v belongs to the union of P_X(w) iff there is a state x
weakly reachable along w from the initial state and a marked path from x with
visible word v. Concatenating those finite paths gives a marked path with visible
word wv, proving forward inclusion in w^-1 Lm(X). Conversely, split any finite
marked path with visible word wv immediately after its |w|-th visible event
(before any remaining private steps). If w is empty, split at the initial state.
The split state lies in R_X(w), and the remainder accepts v. Thus the union is
exactly the derivative. Initial marked-language equality implies equality of
these unions at every word, and union equality at every word implies initial
language equality by choosing w=epsilon. Derivative unions therefore add no
information capable of detecting a hidden doomed alternative.

## 2. Cartesian weak reachability

For all compatible X,Y and w,
R_(X||Y)(w) = R_X(w) × R_Y(w).
Project a product path to each coordinate, deleting the other coordinate's
private steps, to obtain both factor paths: this proves left-to-right inclusion.
For the reverse inclusion, choose finite factor paths witnessing each coordinate
in the right side. Break both paths into private segments separated by the
successive visible letters of w. Interleave the first private segments, synchronize
the first visible letter, and repeat, finally interleaving the trailing segments.
Private independence makes each interleaving legal; full synchronization ensures
the visible letters match. This produces the desired product endpoint. If either
factor reach set is empty, projection rules out a product path.

## 3. Residual intersection

For all states x,y of compatible factors,
L_(X||Y)(x,y) = L_X(x) intersect L_Y(y).
Projection of a finite accepting product path yields factor accepting paths
with the same visible word and marked endpoints. Conversely, accepting factor
paths with a common word can be interleaved segment by segment as above; their
endpoint is marked by conjunction. This applies to arbitrary residual states,
not only the initial states, and includes epsilon and tau cycles traversed a
finite number of times. Infinite paths are not acceptance witnesses.

## 4. Exact nonblocking criterion

X||Y is nonblocking iff for every finite visible word w, every L in P_X(w),
and every H in P_Y(w), the intersection L intersect H is nonempty.
For the forward implication, choose states x,y realizing L,H. Cartesian weak
reachability makes (x,y) reachable. Its completion suffix exists by nonblocking
and belongs to the intersection by the residual equation. For the reverse
implication, every reachable product state has some finite visible prefix w;
its coordinate languages occur in the two profiles. Their nonempty intersection
supplies a marked continuation. Quantification over an empty profile is vacuous,
because no product state exists for that word. A reachable state with empty
residual language remains an element of the profile and cannot be discarded.

## 5. Substitution and closure

Suppose for every visible w, P_M(w)=P_N(w). For every compatible finite context X,
the criterion in section 4 is then literally the same for X||M and X||N, in both
directions, so their nonblocking classifications agree. Furthermore, combining
sections 2 and 3 gives, for every w,

P_(X||M)(w) = {H intersect L | H in P_X(w), L in P_M(w)}.

Each product residual supplies an intersection by projection; each intersection
is realized by choosing representative states and applying Cartesian reachability.
Therefore replacing M by N preserves the whole product profile, not only its
nonblocking classification. The symmetric statement follows from commutativity
of intersection. Induction over a finite binary synchronized-product expression
thus permits substitution at any leaf provided every intermediate product obeys
the same signature, alphabet, privacy and conjunctive-marking assumptions.
In the ALAN instantiation there is exactly one shared K; its authority, effects,
resource roots and Q are never duplicated by regrouping. This induction is only
for this compatible finite automata grammar.

## Exact decision and certificates

A finite epsilon closure includes its seed states and follows every private edge.
For any chosen pair of residual states, breadth-first traversal of the complete
pair of determinized subset graphs decides language equality: a node disagrees
iff exactly one subset contains a mark. Every node and visible arc is retained,
including empty subsets and arcs after the first separator. Breadth-first order
and fixed alphabet order give a shortest acceptance separator.

For each pair of monitors, compare all declared cross-state residual languages.
Next traverse the complete pair of weak-reach subset graphs from their initial
states. The profiles at a node agree iff every state on either side has at least
one equal-language partner on the opposite side. This is exactly equality of sets
of languages, without canonical representatives, minimization or a bijection.
At the first mismatch retain its shortest visible prefix and the first unmatched
state (left side first, then sorted state ID). Retain an acceptance separator
against every opposite reached state; an empty opposite set has an empty list.
There need not be a single separator against the union of opposite languages.

## Falsifications and limits

The release-first and release-first-doomed monitors have equal initial marked
languages, hence equal derivative unions. Their profile mismatch first occurs at
reserve-capacity:AB. That monitor word is not a legal core-plant execution: the
real product first commits AB, then reserves it. The doomed monitor can then
privately enter state 3 while AB debt is held and unpaid. Core release-first is
nonblocking whereas the doomed product has a reachable noncompletable state.
This falsifies union sufficiency. The duplicate monitor preserves the same
profile despite multiple equal-language alternatives; state counts cannot be
used as a profile comparison.

The loop monitor is nonblocking in core K, but after commit/reserve AB its private
states 1 and 3 form a tau cycle over unchanged unpaid debt. Both states are
unmarked in the product and have finite completion escapes. Repeating the cycle
forever is an allowed infinite abstract execution; thus nonblocking does not
imply inevitable completion. No fairness or physical crash interpretation is
attached to this lasso. Exact paths and escapes are in the export. This private
cycle is deliberately admitted by the abstract finite-monitor grammar. It is not
established as compilable or admissible in the runtime's strongly terminating,
versioned local instruction grammar with bounded loops/forward control flow.
It does not demonstrate private divergence of validated ALAN local programs.

Profile equality is sufficient, not asserted necessary for contextual nonblocking.
A separate generic strictness control has unmarked X with one state and no edges,
and unmarked Y with two states and one apply:AB edge from 0 to 1. Both have empty
initial marked language and are blocking in every pointed compatible context:
its initial product state exists and no product state can ever mark. Their
profiles differ after apply:AB (empty versus containing the empty language).
This proves strictness on the generic marked-automata domain. X and Y lie outside
the all-marked ALAN monitor roster; this does not prove strictness within that
restricted roster. A profile mismatch alone is not a contextual separator.

The question is contextual nonblocking, often studied as conflict equivalence.
Here the selected residual profile is a directly decidable sufficient interface;
we do not claim to compute the weakest conflict-equivalence relation, a minimal
state representation, or an optimal/cheaper procedure. No priority abstraction
result is imported. Manuscript Theorem 4.8 remains open: arbitrary ALAN contexts,
resource/cost observations, refusal, hidden divergence, dynamic shared authority,
full-agent substitution, distributed ownership and split effects are outside
this result. Exact costs are not preserved or compared.

The new evidence is a finite model reconstruction over pinned predecessor
automata. Historical concrete model histories belong to Phase 3C.3 and retain their
original qualification scope; no fresh runtime, endpoint, SQLite or crash capture
is performed here. The source argument and bounded checks remain distinct from
independent checker/review evidence and delivery status.
