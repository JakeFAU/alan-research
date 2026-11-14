# Phase 3C.4 results: residual interfaces and nonblocking substitution

## Question and result

What must a continuation interface retain to preserve a completion suffix from
**every reachable state**, rather than only from the initial state? Under the
[finite synchronized-product contract](00-contract.md), equality of the
trace-indexed sets of individual residual marked languages is sufficient for
contextual nonblocking, and is closed under that product. This is a conditional
[source argument](01-argument.md), supported by independently reconstructed finite
evidence. It is not a mechanized proof or the general ALAN congruence theorem.

For each visible word w, P_X(w) = {L_X(x) : x in R_X(w)}. Its union is exactly
w^-1 Lm(X), so comparing those unions at every word adds nothing to equality of
initial marked languages. Under exact common signatures/alphabet, full visible
synchronization, independent private tau and conjunctive marks, weak product
reachability is Cartesian and each product residual language is an intersection.
Thus nonblocking is precisely nonempty intersection for every pair of residual
languages after every word. Profile equality preserves that condition in every
compatible finite context, and the product profile consists of pairwise
intersections. Induction applies only to this declared finite product grammar,
retaining exactly one shared K in the ALAN instantiation.

This packet replaces finite B continuation monitors over the unchanged whole
Phase 3C.2 plant K. It does not extract an entire independently executing agent.
Authority, effects, debts and Q remain in K. Monitors constrain only B's
reserve/apply AB and commit BC; other visible labels self-loop without changing
private state. There are no priorities, fairness or inevitability premises.

## Measured finite evidence

The [independent audit](03-independent-audit.md) reconstructs all graphs and
certificates from the declared automata; producer outcomes are not its oracle.
Counts are configuration sums, not unique global states or concrete histories.

| Measured item | Count |
| --- | ---: |
| Inherited plants / monitors / inherited controllers | 72 / 7 / 324 |
| K/M products / core K/D/M products | 504 / 2,268 |
| Complete product graphs | 2,772 |
| Reachable product states / arcs | 29,156 / 35,672 |
| States with a completion suffix | 14,410 |
| Initially completable / nonblocking products | 886 / 656 |
| Unordered monitor pairs | 21 |
| Equal initial languages / equal profiles | 10 / 6 |
| Cross-state residual decisions | 223 |
| Language / profile subset-pair nodes | 850 / 89 |
| Profile-equal pair/context substitution rows | 2,376 |
| Nonblocking / blocking substitution rows | 792 / 1,584 |

Every profile-equal pair is checked in all 396 declared contexts; blocked
contexts are retained. Finite-roster agreement is distinct from the conditional
universal-context argument. There is no necessity, minimality, optimality or
reduced-verification-cost claim.

## Retained falsifications and adverse boundaries

The release-first/doomed pair has equal initial languages and derivative unions
but opposite core-product nonblocking. The shortest monitor mismatch, reserve
AB, is not itself a K execution: the real product first commits AB, then
reserves, then enters the doomed private residual with held unpaid debt. Exact
paths and per-opposite-residual separators are retained in the export.

The duplicate monitor validates one-to-many language coverage. Profiles are
sets, not multisets or node bijections. The checker separately exercises a
residual with language {a,b} opposed to residuals {a} and {b}; different suffixes
are needed against each opposite residual, and none distinguishes their union.

The generic unmarked X/Y control distinguishes an unreachable word (empty
profile) from a reached dead residual (profile containing the empty language).
Both are blocking in every pointed compatible context because neither can mark.
It proves generic strictness of the sufficient profile condition, **outside the
all-marked ALAN roster**; it does not establish strictness within that roster or
make every profile mismatch a contextual nonblocking separator.

The loop monitor is nonblocking in core K but has a real prefix to held unpaid
AB debt, a private 1-to-3-to-1 cycle, and finite marked escapes from both cycle
states. An infinite repetition is possible in the abstract finite-monitor
semantics. This refutes nonblocking-implies-inevitable-completion without a
fairness claim. The cycle is not established as compilable/admitted by the
runtime's strongly terminating, versioned local instruction grammar with bounded
loops and forward control. It is not a validated runtime private-divergence bug.

Semantic RED records for union sufficiency and inevitability are retained apart
from compilation/fixture failures. The independent checker accepts two exact
valid CLI inputs and rejects all 48 hostile inputs, including omission of a
non-tree arc after superficial count repair, forged profile coverage, altered
source/signature identities, invalid lasso/escape and strict JSON attacks.
Its nine test methods pass. Initial positive ordering and test-construction
failures remain archived; no hostile input was initially accepted.

## Source, manuscript and open obligations

The three final exports have identical 8,271,854 raw bytes, SHA-256
`b40f1d5ac0aea0f43999502ea318cb128f545eb8fb1fd129451915a5aa84bf1f`.
The producer binds 40 source/provenance records, including nine compiled local
non-test Go sources and the pinned predecessor. The checker binds five sources,
including unchanged Phase 3C.3/3C.2 Python constructors reused with attribution.
The [qualification record](05-qualification.md) separates Task 1 scoped tests,
Task 2 audit, controller root checks, historical coverage and Task 3 verification.
It preserves the original dependency-cache-warning review REJECT and subsequent
evidence-only repair/APPROVE through the controller-owned review record.

Markdown and LaTeX add aligned Phase 3C.4 content and Abstract/Conclusion updates.
The eleven-form grammar count is corrected with version-specific admission
retained. Markdown 3.4 now motivates the proposed Section 4.5 relation without
asserting general congruence has been proved; condensed TeX contains aligned
clarifications and references its existing theorem label. Bibliography entries
are unchanged: existing Tang/Moor 2024 is background, not an imported priority
abstraction theorem. The PDF is rebuilt and inspected; inherited format and
numbering differences remain disclosed.

No new concrete histories, endpoint/runtime campaign, SQLite opening or old
exhaustive transport probes occurred. Historical captures retain their original
scope. General Phase 3C and Theorem 4.8 remain open, including full-agent
substitution, arbitrary ALAN contexts, dynamic/shared authority, refusals,
resource-cost and divergence preservation, distributed ownership, split effects
and runtime refinement. Independent review and delivery are separate obligations;
this result record does not certify integration.
