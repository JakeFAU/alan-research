# Phase 3C.3: completion-preserving continuation interfaces

## Question and chosen scope

Can equal obligations, exact resource holdings and static release dependencies
support replacement of a component whose continuation controls future ordering?
If not, which behavioral interface preserves possible completion in a declared
context grammar?

The user authorized this next experiment after the merged Phase 3C.2 assessment.
The controller chooses a research-only continuation-monitor slice. The unchanged
Phase 3C.2 relation K retains the full A/B/C service, authority and capacity
semantics. A replaceable finite monitor M constrains B's continuation choices.
This is **continuation-monitor substitution**, not replacement of an entire
agent. K still contains B's service contract; no open-agent factorization or
runtime refinement is claimed.

Base: `ab01204a0b5121b81379043576ee7bbe46118a19` (merged PR 51), tree-identical
to the delivered `6576a20395716253631327544def5435211f1822`.

## Approaches and decision

1. Extend the ownership snapshot with more static dependency fields. This is a
   useful negative candidate, but cannot silently encode the complete private
   continuation and still count as the same snapshot.
2. Retain a finite boundary transition contract and compare its successful
   visible languages. This is chosen: it makes future ordering explicit and
   supports a short, conditional substitution argument.
3. Factor complete independently executing agents or synthesize an optimal
   allocator. Those introduce additional ownership, custody and policy questions
   and remain outside this packet.

The choice and context bounds below are controller decisions within the user's
authorization. No allocation-policy optimality or verification speedup is sought.

## Inherited premises and immutable signature

Use the fixed A→B→C chain, originals AB and BC, generation one, exact bindings,
eight service masks, three layouts (`shared-1`, `shared-2`, `separate-1`) and
three reservation policies (`early`, `ready-now`, `stable-ready`) from Phase 3C.2.
Permission provenance, sufficient cumulative expenditure E and reusable Q remain
distinct. Application, authentic result retention and Q release are atomic.
Root identities, limits, original memberships, service and reservation policies
are immutable. Physical owners may interrupt and recover; durable custody
survives. No cancellation, changing authority, permanent loss, split unknown
effects or distributed allocator is added.

The common signature carries the exact old configuration, original bindings,
root descriptors and memberships. Joining monitors never creates or sums pools.
Only K changes the old stages and held bits. Compatibility rejects missing,
extra, forged or conflicting signature records before composition.

## Visible labels and private state

The common visible alphabet has eight exact actions: `commit`,
`reserve-capacity`, `apply`, `consume` for each of AB and BC. Each action retains
its full old binding; actor and root membership are determined by the immutable
signature. Every application remains visible, including grouped internal
endpoints. Opaque result values retain the old abstraction; this is not payload
equivalence. Receipt ACK/release remains distinct from Q release.

A finite marked epsilon-NFA has explicit states, one initial state, marked
states and a complete edge list. `tau` is a separate private label, never a
forged visible action with an empty binding. A tau changes only that automaton's
private state. It cannot read hidden K state, write shared state, obtain
authority, release Q or apply an effect. All visible events synchronize across
the common alphabet; private steps interleave independently. Product success is
the conjunction of factor marks. K is marked exactly when both originals are D.

B owns commit/consume BC and reserve/apply AB. The five fixed M fixtures below
restrict only reserve AB, apply AB and commit BC; consume BC and the four A/C
events always have self-loops. The A/C controller D fixtures restrict only
commit AB and apply BC; their other six labels always have self-loops. This is
an explicitly limited role grammar. Missing externally owned transitions are
rejected, not interpreted as permission to veto another owner.

## Frozen monitor fixtures

Write r = reserve AB, a = apply AB, c = commit BC; O is the other five labels.
Every state of the five M fixtures is marked: private acceptance cannot invent
the capacity failure. Unlisted r/a/c moves are disabled; O self-loops at every
state. The following names and graphs are fixed:

| ID | States / initial | Additional transitions |
| --- | --- | --- |
| `free` | {0} / 0 | r,a,c self-loop at 0 |
| `commit-first` | {0,1,2} / 0 | 0-r→1, 1-c→2; r,a,c self-loop at 2 |
| `release-first` | {0,1,2} / 0 | 0-r→1, 1-a→2; r,a,c self-loop at 2 |
| `release-first-tau` | {0,1,2,3} / 0 | 0-r→1, 1-tau→3, 3-a→2; r,a,c self-loop at 2 |
| `release-first-doomed` | {0,1,2,3} / 0 | release-first edges plus 1-tau→3; state 3 has only O self-loops |

These are new research continuation restrictions, not behavior asserted of the
runtime. In particular stable-ready's old every-reachable-state result cannot
be imported unchanged after removing actions through M.

## Negative candidate and concrete witnesses

The snapshot candidate retains the entire old concrete state, exact holds,
full configuration, capabilities, declared original inventory and static
service/release rules. It forgets M's future transition relation. Compare
commit-first and release-first in `early/shared-1/mask=0` after the same genuine
prefix committing, transmitting/admitting and reserving AB. At this cut AB is
W and held, BC is U and unheld, and AB's static service guard is true. Capture
both histories from the actual initial state; do not rewrite post-states to
make them match.

Commit-first can only commit BC next among progress actions; thereafter AB
cannot apply while BC waits, and BC cannot reserve Q. The full product graph,
not one stalled schedule, determines whether the cut lacks any completion
suffix. Release-first can apply AB while BC=U and then complete both originals.
Both products must compare equal under the candidate snapshot at the cut.

Frozen concrete-history roster (eight):

1. commit-first, early/shared-1/mask0: common cut then blocked BC commitment.
2. release-first, early/shared-1/mask0: same cut and genuine completion suffix.
3. commit-first, early/shared-2/mask0: completion control.
4. commit-first, early/separate-1/mask0: completion control.
5. commit-first, early/shared-1/mask2: servicing-B completion control.
6. release-first-tau, early/shared-1/mask0: interrupt B after reserve, recover,
   execute private bookkeeping, and complete; retained control does not reset.
7. release-first-doomed, early/shared-1/mask0: choose successful branch and finish.
8. release-first-doomed, early/shared-1/mask0: take doomed tau after reserve,
   interrupt/recover B, retain the doomed control state and inability to finish.

Extend the old concrete relation only by M: old transport stutters retain M;
visible progress must match an M edge; a private step requires B up and changes
only M. The old Lift prepares actual custody without advancing M, then takes
the permitted visible step. Recover B before a requested tau when needed.
Every captured event must be a real successor; all results originate in actual
model application. These remain concrete **model** histories, not runtime or
endpoint captures. No exhaustive concrete-product claim is made.

## Positive interface and argument obligations

For a pointed finite automaton M, let Lm(M) be the finite visible words that
have a path to a marked state after erasing only private tau steps. The
interface is an automaton presentation of this language plus its immutable
signature and role contract, not a CanComplete bit or a finite trace sample.

Under exact signature agreement, the common synchronization alphabet, independent
private state and conjunctive marking:

`Lm(K || M) = Lm(K) ∩ Lm(M)`.

Prove both directions by projecting a finite successful product path and by
interleaving finite private segments around matching visible steps. Consequently
Lm(M)=Lm(N) implies equal initial existential CanComplete in every finite
compatible K, including each fixed-chain K combined with a role-valid A/C D.
Associative grouping follows for this same synchronized product, with one K
and hence one Q relation; it is not a union of independently allocated histories.

The elementary language argument is established automata reasoning applied to
this declared boundary, not a new general process-theory theorem. Necessity can
be stated only for a larger grammar containing all exact finite-word testers;
do not claim that the role-restricted ALAN contexts realize those testers.
No coarsest practical interface, optimal policy or cost reduction follows.

Language equality is decided exactly by epsilon closures and reachable pairs
of subsets, with a shortest acceptance-distinguishing word on inequality. It
must account for empty words, nondeterminism, tau cycles and empty subsets.
Finite sampled prefix agreement is not equality. Validate complete certificates
or reconstruct them independently rather than trusting a producer Boolean.

Release-first and release-first-tau must have equal languages despite private
control differences. Release-first-doomed also has the same language, but in
the core ALAN context reaches a state with no completion suffix. Thus the
positive theorem does not preserve every-reachable-state nonblocking, bad
branching, refusals, fairness, inevitability, divergence or exact cost. Initial
language equality is not equality of every residual after the same history.

## Frozen abstract qualification domain

- The 72 inherited K configurations crossed with all five M fixtures: 360
  complete reachable product graphs, all arcs and reverse CanComplete labels.
- The core K = early/shared-1/mask0 crossed with 324 A/C D controllers and all
  five M fixtures: 1,620 complete graphs. D has states {0,1}, initial 0; each
  of four source/label pairs for commit AB and apply BC independently chooses
  absent, target0 or target1, and any of four marked subsets. The other six
  labels self-loop at both states. No tau. All 324 descriptors are retained,
  including contexts with no completion. This is the entire declared bounded
  roster, not all contexts or all ALAN programs.
- Exact language comparisons for all ten unordered pairs of the five M
  fixtures over the full eight-label alphabet. Preserve distinguishing words
  and their acceptance checks for unequal pairs.
- For all 1,620 D cases, compare flat triple-product states/arcs/marks with
  both parenthesizations after explicit tuple normalization. Compare CanComplete
  for the equal-language pairs in both rosters; retain blocked as well as
  completable contexts. Report measured totals as configuration sums.

Independent Python reconstructs this entire abstract domain and the eight
concrete histories from the specification, without importing/calling Go guards.
Hostile controls include omitted non-tree arcs with corrected counts, changed
gates, forged equivalence, distinguishing-word errors, omitted tau, a forged
mark, hidden/forged effect label, wrong signature/root, external-action veto,
false completion, illegal physical/private steps and altered common-cut fields.
Keep every exact hostile input, expected rejection and command result. Failed
checker attempts remain evidence. Captures and evidence are never edited to
manufacture the desired witness.

## Delivery and limits

New code is confined to `research/componentinterfaces/`; records and Python
auditor to `docs/foundational/phase3c3/`. No production code, dependencies, old
models or historical evidence changes. Use Go standard library and Python
standard library, exact captured commands/source hashes and scoped tests/race/
Staticcheck. Reuse historical source-identical qualification with explicit
attribution; do not call it a fresh whole-repository run. Build/module checks,
changed manuscript alignment and PDF inspection are required.

Update the roadmap, Markdown, LaTeX and checked bibliography/PDF implications.
Record proof, finite checks, concrete model captures, prior evidence and review
as separate categories. The current packet does not establish complete-agent
substitution, runtime refinement, arbitrary topology, dynamic authority,
distributed allocation, split effects, full abstraction for ALAN contexts,
universal nonblocking or general Phase 3C closure.
