# Phase 3C.4 evidence schema: `alan/phase3c4/v1`

`Export(io.Writer) error`, invoked by `go run ./research/residualinterfaces/cmd/export`
from the repository root, emits one deterministic JSON object plus newline.
All listed fields are required, no additional semantic fields or records are
allowed. Booleans are JSON booleans, integer counts/indices are JSON integers;
strings, arrays, objects and null have exactly the indicated types. Empty arrays
are `[]`, not null, except descriptor `Tuple` inherited as null and `Mismatch`
explicitly null on equality. No timestamps or host paths enter semantic bytes.

## Input and descriptors

Top level: `Schema`, `Input`, `Alphabet`, `Monitors`, `Plants`, `Controllers`,
`Comparisons`, `Products`, `Substitutions`, `CoreOutcomes`, `Negative`, `Loop`,
`Strictness`, `Counts`, `Sources`.

`Input`: `Path`, `CompressedBytes`, `CompressedSHA256`, `RawBytes`, `RawSHA256`,
`Schema`, `Sources`. Schema is `alan/phase3c3/v1`; path and raw identity are the
contract values, compressed SHA-256 is
`86a9d22f9808a67dd0066b6b77dd2b3c01a2125d10863c6fbe1369097f6624c8`.
`Sources` is the exact predecessor roster of 27 `{Path,SHA256}` records, each
verified against current exact source bytes. Old outcome fields are not imported.
The compressed/raw binding is authenticated before selecting descriptors; this
is a trusted pinned fixture decoder, not a public arbitrary-input API.

Alphabet, signature, automaton, state and edge field names and meanings are
exactly those of `docs/foundational/phase3c3/02-evidence-schema.md` (Alphabet:
`ID,Action,Actor,Root`; automaton: `ID,Role,Signature,Initial,States,Edges`; state:
`ID,Tuple,Marked`; edge: `From,Label,To`). Preserve the exact inherited bound action
and signature objects, including all roots, memberships, original bindings,
agents, generations, layout, policy and service mask. No normalization erases
identity or root distinctions. Five inherited M descriptors are followed by
release-first-duplicate and release-first-loop as specified in the contract.
The new monitors have four all-marked states in decimal order. Their edges are
ordered by state ID for the five unrestricted loops, then 0-r-1, 1-a-2, then
duplicate's 1-a-3 or loop's 1-tau-3,3-tau-1,3-a-2, then r/a/c loops at state 2
(and state 3 for duplicate). The alphabet remains the eight predecessor labels
in AB/BC and commit/reserve-capacity/apply/consume order; tau is private.
Plants and controllers retain their exact inherited order and all descriptors,
including unmarked and unreachable D states. No external M/D action is vetoed or
changes private state. Monitor rebinding uses the selected plant's entire exact
signature. Generic controls use Role `generic` and Signature `{}` only.

## Language certificates and profile comparison

`Comparisons` contains all 21 unordered monitor pairs in monitor-order nested
loops. A profile comparison has `Left`, `Right`, `InitialLanguage`, `Residuals`,
`Pairs`, `Equal`, `Word`, `Mismatch`. Left and Right identify the automata.

A language comparison (InitialLanguage and every residual Language) has `Left`,
`Right`, `Equal`, `Word`, `LeftAccept`, `RightAccept`, `Pairs`. Left/Right remain
automaton IDs; the selected starting state is given by the enclosing residual
record, or the automaton Initial for InitialLanguage. Pairs is the complete
breadth-first reachable epsilon-closed subset-pair graph. Each node has `Left`,
`Right` (sorted state ID arrays), `LeftAccept`, `RightAccept`, `Arcs`; each arc
has `Label`, `To` (zero-based pair index). Every node has all eight alphabet arcs
in alphabet order, including transitions to empty subsets. Initial index is zero.
Acceptance means the subset contains at least one mark. Equal is true iff no
node disagrees. On inequality, Word is the first shortest acceptance separator
(BFS with alphabet-order tie break), with actual acceptance booleans. On equality,
Word is `[]` and both top-level acceptance fields are false sentinels; they do
not describe epsilon acceptance. Per-node booleans always describe actual marks.

`Residuals` has every declared left-state/right-state cross pair in descriptor
state order. Each has `LeftState`, `RightState`, `Language` (complete comparison).
Unreachable descriptor states are included. No class minimization is assumed.

Profile `Pairs` is the complete BFS weak-reach subset-pair graph from the initial
states, with all eight arcs at every node. Nodes have `Left`, `Right`,
`LeftCoverage`, `RightCoverage`, `Equal`, `Arcs`. Coverage arrays follow sorted
state order and contain `{State,Matches}`. Matches lists **all** opposite reached
states with equal residual language, in sorted order. Equal requires every
coverage list on both sides to be nonempty. Both empty subsets compare equal;
one empty subset versus a reachable empty-language state does not. Multiplicity
is ignored through coverage, so one-to-many matches are valid.

The top Equal is true iff every profile node agrees. Word is the first shortest
mismatch prefix by BFS/alphabet order. On equality Word is `[]` and Mismatch is
null. Otherwise Mismatch is `{Side,State,Separators}`. Side is `left` if any left
state is unmatched, otherwise `right`; State is the first sorted unmatched state
on that side. Separators lists every opposite reached state in sorted order.
Each record has `OppositeState`, `Residual` (index into Residuals), `Word`,
`LeftAccept`, `RightAccept`, copying that unequal residual certificate's shortest
separator. Left/right acceptance orientation stays the comparison orientation,
including when the unmatched side is right. An empty opposite subset has exactly
`Separators: []`. No separator against the opposite union is required.

## Complete products and substitution

Products follows 72 K outer × seven M inner (504), then 324 D outer × seven M
inner over core early/shared-1/mask0 K (2,268), total 2,772. Graph schema is the
predecessor graph schema: `ID,Factors,Initial,Nodes,States,Arcs,CanComplete,Nonblocking`.
Nodes have `ID,Tuple,Marked,CanComplete,Arcs`; arcs have `Label,To` (target node ID).
Tuple order is K,M[,D], IDs compact JSON arrays of those atomic state IDs.
Every reachable state and every relation arc is included. BFS discovers private
successors in factor order with sorted targets, then visible successors in
alphabet order with lexicographic Cartesian target enumeration. Stored arcs are
sorted lexicographically by (Label,To) and deduplicated. Marks are conjunction;
CanComplete is reverse reachability from all marked states, including each mark
itself. Graph CanComplete is a count; Nonblocking means every reachable node
CanComplete. Counts are redundant measured checks, never substitutes for graphs.

Substitutions contains every profile-equal pair (comparison order) across all
396 contexts (72 K then 324 core K/D, inherited order). Each row has
`Left,Right,Context,ProductLeft,ProductRight,NonblockingLeft,NonblockingRight`.
Context ID joins K and optional D with ` || `; products retain K,M[,D] order.
Both nonblocking and blocking rows must remain. CoreOutcomes has seven records
in monitor order, each `Monitor,Product,InitialCanComplete,Nonblocking`.
This retains the difference between initial possible completion and global
nonblocking even for initial-language-equal/profile-unequal pairs.

## Witnesses and separate generic control

A path has `Start,Steps,End`. Steps have `From,Label,To`; every step must be an
actual edge of its referenced product and adjacent endpoints must agree.
A shortest path uses BFS with stored lexicographic (Label,To) arc order.

Negative has `Comparison` (release-first/doomed profile index), `LeftProduct`,
`RightProduct`, `LeftPath`, `RightPath`. The left path is the shortest path from
initial to tuple `[W,U,true,false,1]` (shown here as two conceptual components:
K key `W,U,true,false`, M key `1`); right ends at the corresponding M key `3`.
Both first commit then reserve AB; right then takes private tau. The shortest
monitor mismatch Word (reserve alone) remains in its separate comparison. The
left endpoint can complete, the right cannot. The core products have opposite
Nonblocking and equal initial CanComplete.

Loop has `Product,Prefix,Cycle,Escapes`. Product identifies core K/loop. Prefix
is the shortest path to K key `W,U,true,false`, M key `1`. Cycle is exactly the
two tau edges 1->3->1 at that unchanged K state. Escapes contains shortest marked
paths from cycle state 1 then state 3. Both cycle states have unpaid held AB debt,
are unmarked, and have CanComplete=true. End of Cycle equals its Start and prefix
End. These are finite abstract paths describing an infinite possible repetition,
not runtime/crash captures or a fairness result.

Strictness has `Left,Right,Comparison`, separate from Monitors/Products/counts.
Left is generic-unmarked-X, one unmarked state 0 and no edges. Right is
generic-unmarked-Y, unmarked states 0,1 and edge 0-apply:AB->1. Both Initial 0,
Tuple null, Role generic, Signature `{}`. Comparison is the complete profile
comparison above, including its equal initial-language certificate and its
empty-opposite separator list at the shortest mismatch. Both automata are always
blocking in pointed compatible contexts because they have no marks. This
strictness example is outside the all-marked ALAN roster.

## Counts, source binding and capture

Counts fields: `Plants,Monitors,Controllers,Products,ProductStates,ProductArcs,
ProductCanComplete,InitialCanComplete,NonblockingProducts,Comparisons,
EqualLanguages,EqualProfiles,ResidualComparisons,LanguageSubsetPairs,
ProfileSubsetPairs,Substitutions,NonblockingSubstitutions,BlockingSubstitutions`.
Product totals sum the 2,772 declared graphs. LanguageSubsetPairs sums initial
and all cross-state certificates across the 21 ALAN comparisons; ProfileSubsetPairs
sums their weak subset graphs. Generic strictness is excluded. Substitution totals
count rows with both booleans true or both false; unequal booleans are invalid.
All non-roster counts are measured, not assumptions.

Sources is a path-sorted array `{Path,SHA256,Use}`. Use is one of `compiled producer`,
`module`, `hash-checked specification`, `hash-checked predecessor provenance`,
`decoded pinned automata`. All local non-test Go files actually compiled by the
exporter, including main, are bound; root module files, the three Phase 3C.4
contract/argument/schema files, pinned compressed bytes and all 27 checked old
sources are bound. Shared go.mod/go.sum receive Use module. Old code is only
hash-checked, not imported or executed. Exact `go list -deps -json` evidence
cross-checks the compiled-local roster separately. Tests, README and mutable
reports/logs/manifests are not semantic producer inputs. Artifact manifests bind
raw and deterministic gzip bytes and two identical captures. Commands, explicit
cache environment, stdout/stderr and exits are retained separately, including
initial missing-symbol RED and any resolved diagnostics. Deliberate semantic
hypothesis failures are distinct records from those compilation failures and
from passing tests which assert counterexamples.
