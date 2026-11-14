# Conditional argument: completion languages and residual antichains

## Domain and observations

This is a source argument, not a mechanized proof or runtime refinement. Fix a
common finite visible alphabet Sigma and finite pointed marked epsilon-NFAs
with exactly compatible signatures. Every visible event synchronizes all factors;
a tau step changes only its factor's private state and cannot inspect other
factors. Marks are conjunctive. All acceptance uses finite paths, including
zero steps and leading/trailing tau. The grammar is finite synchronized products
with one unchanged whole K; effects, shared authority and Q remain inside K.
No priority, fairness, inevitability, divergence preservation or unrestricted
agent substitution is assumed.

Let R_X(w) be the epsilon-closed set of states reachable by visible word w.
Let L_X(x) be the marked finite-word language from state x, and let
P_X(w) = {L_X(x) : x in R_X(w)}. These are finite *sets of languages*, even when
the languages are infinite. Empty reach gives the empty family; a reachable
state with no marked continuation gives a family containing the empty language.
The two are different. Define A(F) by first identifying equal languages and
then retaining precisely the inclusion-minimal members. Empty intersections
are never dropped. Define J(X) = (Lm(X), w -> A(P_X(w))). Equality is quantified
over every finite visible word, including epsilon and unreachable words.

CanComplete at a state means existence of a finite marked path; nonblocking
means every reachable state has CanComplete. Neither implies inevitability.
The graph certificates decide finite automata exactly by full subset exploration,
not by a chosen maximum word length.

## Finite coinitiality and hit preservation

For every finite family F and every L in F, there is M in A(F) with M subseteq L.
Starting at L, descend along a proper included member while possible. Finiteness
precludes an infinite strict descent; the resulting member is minimal in F.
Equivalently one can choose a minimal member of the nonempty finite subfamily
{N in F : N subseteq L}. Equivalent duplicates are collapsed first.

For *every* suffix language H (finite or infinite), define
T_F(H) iff every L in F has L intersect H nonempty. Then
T_F(H) iff T_A(F)(H). The forward implication restricts the quantified family.
For the reverse, choose the contained minimal member for each L and use its hit
as a hit of L. When F is empty both predicates are vacuously true; when F
contains the empty language both are false, since that language is minimal.
Finiteness is essential: an infinite strictly descending family can have no
minimal member, making its reduced hit test vacuous.

## Product reduction law

Write F tensor G = {L intersect H : L in F, H in G}. For finite F,G,

    A(F tensor G) = A(A(F) tensor A(G)).

If either factor family is empty, both sides are empty. Otherwise every member
of F tensor G contains a member of A(F) tensor A(G), by choosing contained
minima of both factors. The latter family is also a subfamily of the former.
For a general finite D subseteq C coinitial in C, their minima coincide: if
c is minimal in C, a contained d in D must equal c; if d is minimal in D,
a contained minimal c in C contains a d' in D, forcing d'=c=d. Apply this fact.
A minimal intersection may be empty; retaining it is necessary for the law and
for truthful nonblocking tests.

## Reachability, residual intersections, and closure

Under the stated synchronization/private-state hypotheses, for every word w,

    R_(X||Y)(w) = R_X(w) x R_Y(w),
    L_(X||Y)(x,y) = L_X(x) intersect L_Y(y).

Projection proves one direction of both claims. Conversely, split each chosen
finite factor path into private segments around the identical sequence of
visible events. Interleave the independent private segments, synchronizing
corresponding visible steps. This constructs a product path, including leading
and trailing private segments. Conjunctive marks give residual intersection.
This construction fails without the private independence and synchronization
hypotheses; it says nothing about infinite divergent runs.

Thus P_(X||Y)(w)=P_X(w) tensor P_Y(w). If for every w,
A(P_X(w))=A(P_X'(w)), the reduction law gives equality of the reduced profiles
of X||Y and X'||Y for every compatible finite Y. The same holds for replacing
both factors with reduced-profile-equal factors.

For each w, nonblocking of the reached product states means each L in P_X(w)
intersects each H in P_Y(w). Applying hit preservation separately for each H
shows that reduced-profile equality preserves this condition. Quantifying over
all words establishes preservation of nonblocking in every compatible finite
context. Private steps are covered because each reach set is epsilon-closed.

Initial language satisfies Lm(X||Y)=Lm(X) intersect Lm(Y), by the same path
argument at the designated roots. Therefore J equality preserves both initial
CanComplete (nonemptiness of Lm) and nonblocking, and is closed under product.
Structural induction applies to the declared finite synchronized grammar with
one unchanged K. It is not an induction over arbitrary runtime programs,
authority changes, distributed ownership, or multiple interacting K copies.

## Exact local characterization and its limits

For all finite language families F,G:

    [for every finite H, T_F(H) implies T_G(H)]
      iff [for every K in G, some L in F satisfies L subseteq K].

For sufficiency, choose such an L for each K. A source hit lies in K, establishing
the target hit. This also holds for infinite H. For necessity, suppose coverage
fails at K in G. For every L in F choose a word u_L in L minus K, which exists
because L is not a subset of K. Let H={u_L : L in F}. Finiteness of F makes H
finite, and H hits every source member while missing K. If F is empty, choose
H empty: the source predicate is vacuously true and the target predicate false.
If F contains the empty language, coverage cannot fail at any K, consistently
with the source predicate being false for every H. H empty and H={epsilon}
are distinct constructions, represented differently in the evidence.

Mutual coverage is equivalent to equality of minimal families. If minima agree,
coinitiality supplies coverage in both directions. Conversely, for any minimal
f in F choose g in G contained in f and f' in F contained in g. Minimality
forces f'=g=f. It is minimal in G too, since any proper included G member would
contain an F member properly included in f. Reverse the argument for G.
Finiteness is used for minimal families and for the finite separating H; these
claims do not assert their unrestricted infinite-family versions.

This is an exact characterization of *local finite-suffix hit tests*. A trie
for H, synchronized with one separately rooted residual, tests initial possible
completion. It does not in general distinguish global nonblocking: q-a->marked f
and q-a->unmarked dead d intersect H={a}, but the product can reach the dead d
branch and is blocking. The evidence preserves this finite trie product.
Local necessity therefore does not imply global contextual necessity.

## Strict reduction and adverse controls

The original contract fixtures are used verbatim. All states are marked.
In minimal-augmented, L1 is properly included in L4, properly included in L0:
state 4 adds c to the choices of state 1; state 0's direct r/a/c transitions to
unrestricted state 3 include those of state 4 and add a. The five other events
loop everywhere. At epsilon the base family contains L0,L1; the augmented
family additionally contains L4. Its redundancy disappears under MIN. After
any visible word, extra states from branch 4 are already present from the
base branches except for the redundant state 4 while only unrestricted events
have occurred. Initial languages agree. Full graph exploration certifies these
all-word claims and strictness without relying on this informal case split.

Both differences are tested against actual legal core-K traces by marking every
core-K state in a separately identified *control* plant, then comparing separately
rooted products at 1 versus 4 and 4 versus 0. This changes only the legal-trace
control, not the original 3564 products. It establishes plant relevance of the
differences, not paid completion or a runtime execution. The rejected alternative
in preflight remains design history, outside the monitor roster and qualification.

MAX is unsound for nonblocking: release-first and release-first-doomed have equal
initial languages and equal all-word MAX profiles, but opposite core nonblocking.
The doomed residual at state 3 is a proper sublanguage of state 1 and is removed
by MAX. Actual core product paths reach the healthy and doomed unpaid states.

MIN alone is insufficient for initial language in the generic marked domain:
X starts at marked s, with tau to unmarked d, and every visible event loops at
both states; Y is unmarked d with all loops. Every MIN profile is {empty language},
while their initial languages are Sigma* and empty. J separates them. For an
all-marked leaf, each reachable state contains epsilon; empty/nonempty reduced
profiles determine its trace/marked language. That leaf shortcut cannot govern
composition with K, whose conjunctive marks are nontrivial.

J is sufficient, not necessary for the two global contextual observations in the
generic domain: unmarked no-edge X and unmarked single-a-edge Y have empty initial
languages and unequal MIN profiles after a. In every pointed conjunctively marked
context both products have no marks, so initial CanComplete and nonblocking are
both false. This is a structural argument for all such contexts, separate from
the explicit example trie products. Neither generic example belongs to the
all-marked ALAN roster.

No novelty, minimal automaton, weakest global contextual relation, speedup,
runtime refinement, exact-cost/refusal/divergence equivalence, or broader
Phase 3C closure is claimed. The independent auditor and manuscript reconciliation
are later tasks; this packet's Go tests and source arguments do not substitute
for them.
