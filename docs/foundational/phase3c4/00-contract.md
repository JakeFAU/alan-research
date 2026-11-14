# Phase 3C.4: residual interfaces and nonblocking substitution

## Question and chosen boundary

Phase 3C.3 preserves initial existential CanComplete while admitting a doomed
residual. What additional behavioral information is sufficient to preserve the
existence of a completion suffix from **every reachable state** after replacement?
Call that property nonblocking; it does not mean that every execution finishes.

This packet replaces finite B continuation monitors over the unchanged whole
Phase 3C.2 plant K. It does not extract an entire independently executing agent.
Base: merged Phase 3C.3 `c369e25c8db268d418ab7c11573d46dbc1d9425f`, whose tree
equals delivered `fc3483aa3a0a1ce0d8ff7110dd225815d11d855b`.

Controller choice under the continuing research authorization: preserve the set
of individual residual marked languages after each visible word. Initial
language equality, even reformulated as equality of the union after every word,
forgets internal alternatives. A bisimulation could supply a stronger sufficient
rule but adds structural matching obligations. General conflict-equivalence
decision/minimization is a larger problem. The selected interface has a direct
finite decision procedure and a short conditional argument. It is not claimed
necessary, weakest, state-minimal, optimal or cheaper to verify.

## Semantics and mathematical obligations

Use finite pointed marked epsilon-NFAs. All compatible factors have the same
exact signature and visible alphabet; every visible step synchronizes all
factors. Private tau changes one factor's private state only, without reading
or writing another factor. Marks are conjunctive. Acceptance uses a finite
path, with tau erased, including leading/trailing tau and the empty word.
No event priorities, fairness, inevitability or divergence preservation are
assumed. Shared authority, effects and Q remain in the one unchanged K.

For visible word w, let R_X(w) be every state weakly reachable from X's initial
state along w, including trailing tau. Let L_X(x) be the marked visible language
from individual state x, including its own private closure. Define
P_X(w) = { L_X(x) : x in R_X(w) }, a set of languages, not their union or a
multiset. An unreachable word gives the empty profile; a reachable dead residual
gives a profile containing the empty language. These cases must remain distinct.

The source argument must establish, with both directions and exact quantifiers:

1. The union of P_X(w) is the left derivative w^-1 Lm(X). Comparing these unions
   for every word therefore adds nothing to initial marked-language equality.
2. R_(X||Y)(w) = R_X(w) x R_Y(w), by projection and interleaving finite private
   segments. Private independence and full visible synchronization are essential.
3. L_(X||Y)(x,y) = L_X(x) intersect L_Y(y), using the inherited finite-path
   intersection argument at the particular residual states.
4. X||Y is nonblocking iff for every visible w, every L in P_X(w) and H in
   P_Y(w) have nonempty intersection. Quantification over empty reach sets is
   vacuous; reachable states with an empty residual language are not removed.
5. Equality P_M(w)=P_N(w) for every w preserves nonblocking in every compatible
   finite context X. Further, P_(X||M)(w) consists exactly of pairwise language
   intersections from the factor profiles, so profile equality is closed under
   this product. State the resulting induction only for this compatible finite
   synchronized-automata grammar, with exactly one shared K when instantiated.

This is a conditional source argument, not a mechanized proof or a new general
process-theory theorem. It does not establish manuscript Theorem 4.8, whose
stronger ALAN context, resource, refusal, effect and divergence obligations stay
open. Relate the question to established contextual nonblocking/conflict
equivalence without claiming to decide its weakest relation or importing a
priority-based abstraction theorem.

## Immutable predecessor input and new implementation

New Go code is confined to `research/residualinterfaces/` and its exporter.
Do not modify old model, runtime, module, test or evidence files. The sole new
public boundary is `Export(io.Writer) error`, invoked from the repository root
by `go run ./research/residualinterfaces/cmd/export` with JSON on stdout and
diagnostics on stderr. Internal names and representations are builder-owned.

Read the pinned predecessor capture at
`docs/foundational/phase3c3/evidence/task1/export.json.gz`. Its decompressed
length is 25,384,964 bytes and SHA-256 is
`a073a22fe8ac62e05f2d15b0f5146b19db9606e38ce4619261e7e4108e27269d`.
Verify compressed and raw identities and schema `alan/phase3c3/v1`; bind the
compressed identity freshly from the repository rather than inventing it.
Verify its 27 declared source hashes against current files. Explicitly identify
the imported alphabet, 72 plants, five monitors and 324 controllers. Ignore old
Products, Comparisons, Groupings, Substitutions and outcome/count fields as
oracles: construct new product graphs and reverse reachability from automata.
Retain exact signatures, roots, memberships and bound action descriptors.
Rebinding monitor syntax to another inherited configuration is explicit and
uses that plant's entire signature, as in Phase 3C.3.

The new source manifest binds every actual producer input: compiled local
non-test Go sources, module files, this contract, argument and schema, pinned
predecessor bytes and any predecessor source files read for verification.
Distinguish compiled operational dependencies from hash-checked provenance.
Cross-check the actual compiled local file roster with `go list -deps -json`.
No new dependency or general external-automaton parsing API is introduced.

## Frozen monitor and context roster

Retain the five predecessor monitors without changing their behavior. Add two
all-marked M monitors, each with states 0,1,2,3 and initial state 0. Write
r=reserve-capacity:AB, a=apply:AB, c=commit:BC and O=the other five visible labels.
Every state has an O self-loop for each O label; unlisted r/a/c edges are absent.

| Monitor | Additional edges |
| --- | --- |
| release-first-duplicate | 0-r->1; 1-a->2; 1-a->3; r/a/c self-loops at both 2 and 3 |
| release-first-loop | 0-r->1; 1-a->2; 1-tau->3; 3-tau->1; 3-a->2; r/a/c self-loops at 2 |

Monitors restrict only B's reserve/apply AB and commit BC; all other visible
labels self-loop without changing private control. Context D retains the exact
324 predecessor two-state A/C controllers, including empty markings and
unreachable states. Do not discard awkward descriptors to improve counts.

Enumerate 72 inherited K configurations x seven M = 504 products, and core
early/shared-1/mask0 K x 324 D x seven M = 2,268 products: 2,772 complete graphs.
Keep every reachable node, mark, named arc and CanComplete label. Nonblocking
means all reachable nodes have CanComplete=true. Counts other than the frozen
roster are measured from actual runs. Preserve blocked contexts too.

## Exact interface decision and witnesses

Compare all 21 unordered M pairs. For every cross-pair of their declared states,
decide equality of residual marked languages by complete epsilon-closure/subset
pair search, including empty subsets, tau cycles and nondeterminism. Retain
complete comparison relations and shortest acceptance separators when unequal.

Then explore all reachable pairs of weak reach-state subsets from the monitor
initial states, retaining every visible arc even after finding a mismatch.
Their profiles agree exactly when each state on each side has a residual-language
equal partner on the other side. This is bidirectional coverage, not a bijection
and not a comparison of state counts. A mismatch retains a shortest visible
prefix, an unmatched side/state and its residual separator against **each**
opposite reached state. There need not be a single suffix separating that state
from the opposite union. Empty opposite subsets are handled explicitly.

Generic algorithm controls include unmarked X with one initial state and no
edges, and unmarked Y with two states and one a-edge from its initial state.
After a their profiles are empty versus containing the empty language. Both
remain blocking in every pointed context because neither can ever mark. Retain
this small strictness example separately from the all-marked ALAN roster: a
profile mismatch is not automatically a contextual nonblocking separator.

Record initial-language equality separately from profile equality. For every
profile-equal pair and every one of the 396 declared K or K/D context descriptors,
record both product nonblocking classifications, retaining both true and false
cases. All seven monitors are compared in core K, so language-equal/profile-unequal
cases and their global outcomes remain visible. Do not infer necessity from a
profile mismatch or call equal finite-roster outcomes universal equivalence.

Retain three distinct falsification boundaries:

- The inherited release-first/doomed pair has equal initial languages and hence
  equal derivative unions at every word, yet opposite core-product nonblocking.
  The shortest monitor profile prefix r is distinct from a real K path, which
  must first commit AB. Retain actual product paths and the unmatched residual.
- The duplicate monitor tests one-to-many coverage of equal residual languages.
  Multiplicity is not observable in this profile.
- In core K, retain a real prefix reaching unpaid held debt in the loop monitor,
  a tau cycle 1->3->1, and an actual finite completion escape. Every cycle state
  has a completion suffix but an infinite private execution is possible. This
  refutes nonblocking-implies-inevitable-completion without making a fairness
  claim. These are abstract paths/lassos, not new physical crash captures.

## Evidence and independent challenge

Task 1 freezes `01-argument.md` and `02-evidence-schema.md`, captures the complete
deterministic new export twice from final source, and retains exact bytes,
source hashes, commands, exits and diagnostic failures. Semantic falsification
records remain separate from initial missing-symbol compilation REDs.

Task 2 uses new Python profile/language/witness checks independent of Go. It may
reuse the unchanged Phase 3C.3 Python plant/monitor/controller constructors and
their earlier Python dependencies with explicit hash binding and attribution.
It must reconstruct the full new product relation, reverse reachability and
interface certificates; no Go predicates or producer verdicts serve as oracle.
Keep exact accepted and rejected inputs, including a missing non-tree arc with
discovery parents and superficial counts repaired, forged coverage/union
substitution, omitted unmatched-state separators, bad shortest prefixes,
tau/mark changes, dead-vs-unreachable confusion, external vetoes, signature/root
changes, fabricated nonblocking, invalid lasso/escape and JSON type/key attacks.
Any initially accepted hostile input remains retained before repair.

Fresh scoped ordinary/race Go tests, scoped Staticcheck, checker tests, complete
valid/hostile CLI checks and root build/module/package checks qualify this slice.
Do not repeat unchanged suites merely for repetition. Source-identical historical
qualification remains explicitly attributed; no fresh whole-root suite or global
lint-clean claim follows. No new concrete histories, endpoint/runtime campaign,
SQLite opening or old exhaustive transport-probe execution is required or claimed.

## Manuscript, requested corrections and delivery

Align new Markdown/LaTeX mathematics, Abstract/Conclusion and roadmap with the
actual result; check bibliography implications and rebuild/inspect the PDF.
Also correct the user's verified editorial findings: the displayed grammar in
Markdown 3.1.2 lists eleven forms (Set, Add, Branch, Jump, Receive, Emit, End,
ForEach, Next, Loop, Append), not ten; Markdown 3.4 must no longer claim that the
general congruence failures have been resolved by a proof. Check equivalent
LaTeX locations, adding concise aligned clarification where its condensed
presentation has no corresponding sentence. Do not claim these edits resolve
all inherited manuscript differences or the open Theorem 4.8.

Preserve source/evidence/review/delivery boundaries and rejected reviews. Deliver
a narrowly scoped draft PR with independent task and whole-branch review under
standing authorization; retain the worktree and do not merge. General Phase 3C,
full-agent substitution, dynamic authority, distributed ownership and split
effects remain open.
