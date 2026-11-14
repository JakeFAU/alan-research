# Evidence schema `alan/phase3c5/v1`

The sole public Go function is `Export(io.Writer) error`; run it from repository
root with `go run ./research/residualantichains/cmd/export`. It emits one JSON
object and newline with no semantic timestamps or host paths. All fields below
are required. No duplicate keys, added fields, records or alternative types are
part of this schema. Counts and indices are JSON integers, flags JSON booleans.
Empty lists are `[]`; descriptor Tuple preserves inherited null; witness pointers
and successful coverage Failure use explicit null as specified below.

## Inputs and sources

Top level: Schema, Input, Alphabet, Monitors, Plants, Controllers, Comparisons,
Products, Substitutions, CoreOutcomes, Negative, LegalControl, Generics, Counts,
Sources. Input has Path, CompressedBytes, CompressedSHA256, RawBytes, RawSHA256,
Schema, Sources. The path is `docs/foundational/phase3c4/evidence/task1/export.json.gz`;
compressed size/hash are 157341 and
`35115c3104cddd8a789c80a8c0e73eb7c7f173632af512776ba8502b61b9afc9`, raw
size/hash 8271854 and
`b40f1d5ac0aea0f43999502ea318cb128f545eb8fb1fd129451915a5aa84bf1f`.
Schema is `alan/phase3c4/v1`. Sources is its exact ordered 40-entry Path/SHA256
roster, checked against current files. Only alphabet, signatures and the seven
M, 72 K and 324 D descriptors are selected; old verdicts/graphs/counts are not
oracles. Input is trusted pinned data, not a general external-input API.

Alphabet and descriptors retain Phase 3C.4 structure and order. Alphabet entries
are ID, Action (the original object), Actor, Root. Automata are ID, Role,
Signature (exact complete original object), Initial, States, Edges. States are
ID, Tuple, Marked; edges From, Label, To. Signatures retain every inherited root,
membership and action binding. Rebinding a monitor for a product uses the whole
selected plant signature. Tau is private. Labels have the inherited AB then BC,
commit/reserve-capacity/apply/consume ordering.

The two new monitors follow the inherited seven, in minimal-base then
minimal-augmented order, with the exact contract states and transitions. Edges:
state-order five unrestricted loops first; then 0-tau-1; then for r,a,c in that
order, 0-label-3 followed by 3-label-3; then 1-r-2,2-a-3; augmented alone appends
0-tau-4,4-r-2,4-c-3. Decimal state order, all marked, initial 0, Tuple null.
No rejected preflight fixture is in the roster.

Sources is a path-sorted array of Path, SHA256, Use. Use distinguishes compiled
producer, module, hash-checked specification, hash-checked predecessor provenance,
and decoded pinned automata. New non-test Go files and main are compiled sources;
all 40 inherited entries are hash-checked, except module files retain module use.
The three contract/argument/schema files are specifications. The pinned gzip is
actually decoded. All entries are actually read and hashed. No old Go package is
imported. Tests, README, capture helper and mutable reports are not semantic
producer inputs; the capture manifest binds the helper and records separately.

## Exact language certificates

A comparison has Left, Right, Equal, Word, LeftAccept, RightAccept, Pairs,
LeftSubset, RightSubset, LeftOutside, RightOutside. Pairs is the complete BFS
reachable graph of epsilon-closed subset pairs, initial index zero. Each node
has Left, Right (sorted state IDs), LeftAccept, RightAccept, Arcs. Every Arcs list
has all eight labels in alphabet order with Label and integer To index, including
empty-subset transitions. Acceptance means at least one marked state. Exploration
continues after every failure. Leading/trailing tau and epsilon are included.

LeftSubset means left language is included in right; its unique failure condition
is left-accept and right-reject. LeftOutside is null if included, otherwise the
shortest word in left-minus-right, with alphabet-order BFS ties. RightSubset and
RightOutside reverse the direction. Epsilon is `[]`, distinct from null.
Equal is both subset flags. Word is the first shortest XOR witness (either
direction), with its acceptance flags; on equality Word=[] and both top-level
acceptance flags are false sentinels. Per-node flags are always actual.

A residual comparison has LeftState, RightState, Language. The enclosing automata
are rooted at those states; Language still uses their automaton IDs. All declared
states, including unreachable ones, are considered in descriptor nested order.

## Profiles and exact local failures

Comparisons includes all 36 unordered monitor pairs, roster nested order. Each
profile comparison has Left, Right, InitialLanguage, LeftWithin, RightWithin,
Residuals, Pairs, FULL, MIN, MAX, J, FULLWord, MINWord, MAXWord. LeftWithin and
RightWithin contain complete directed comparisons of all ordered pairs within
the corresponding monitor (including self); Residuals contains every cross pair.
Repetition of within certificates across comparisons is deliberate and exact.
InitialLanguage roots at the designated initial states.

Profile Pairs is the same complete weak-reach graph as InitialLanguage.Pairs,
with node fields Left, Right, LeftFamily, RightFamily, LeftToRight, RightToLeft,
FULL, MIN, MAX, Arcs. Nodes, sets and arcs have exactly the language graph's order.
A family has Classes, MIN, MAX. Classes identifies languages among that reached
set; each class has Representative (first sorted state ID), Members (all equal
states), ProperSubsets, ProperSupersets (representatives of *all* strict relations
in class order). Equality is collapsed before filtering. MIN retains precisely
classes with no proper subsets, MAX with no proper supersets. Empty family has
all three lists empty; an empty-language member remains a class.

FULL compares both full class sets by language equality; MIN and MAX compare
those respective representative sets. Top flags quantify all graph nodes.
FULLWord, MINWord, MAXWord are null on all-word equality, otherwise the first
shortest failing visible prefix in alphabet-order BFS (epsilon is []).
J = InitialLanguage.Equal and MIN, independently of FULL or MAX.

LeftToRight tests source left family to target right family; RightToLeft reverses.
There is one coverage record per target class in target order: Target, Sources,
Failure. Sources contains *every* source representative whose language is a
subset of the target, in source order. Failure=null iff Sources is nonempty.
Otherwise Failure has Outside, H, Hits, TargetMiss. Outside has Source, Word for
every source class: the shortest word in that source-minus-target, taken from
the directed certificate. H deduplicates those words in first source order.
Hits has Source, Hit for every source class, checked by actual acceptance of H;
TargetMiss checks no H word accepted by target. Every Hit and TargetMiss is true
on genuine failures. Empty source yields Outside=[],H=[],Hits=[],TargetMiss=true.
H=[[]] instead denotes the singleton epsilon. Every failed target at every node
is retained; no shortest-prefix-only truncation is permitted. These are *local*
finite-suffix hit separators, not global nonblocking contexts.

## Products and substitution rows

Products is all 72 K x nine M, followed by 324 D x nine M using core
`early/shared-1/mask0`, total 3564, with factor order K,M[,D]. Each graph has ID,
Factors, Initial, Nodes, States, Arcs, CanComplete, Nonblocking. Nodes have ID,
Tuple, Marked, CanComplete, Arcs; arcs Label, To (node ID). IDs are compact JSON
arrays of the flattened atomic state IDs. Initial node is first. BFS explores
private successors in factor order, then visible labels in alphabet order with
sorted Cartesian successors; stored arcs sort by (Label,To) and deduplicate.
Every reachable node and relation arc is present. Marks are conjunction;
CanComplete labels come from complete reverse reachability of marks. States,
Arcs and CanComplete are measured counts. Nonblocking iff every node can complete.

Substitutions retains every J-equal pair over all 396 contexts, in comparison
then 72 K/324 K-D order. Fields: Left, Right, Context, ProductLeft, ProductRight,
NonblockingLeft, NonblockingRight, InitialLeft, InitialRight. Context joins K and
optional D with ` || `; product factor order remains K,M[,D]. Both flags must
agree pairwise, but all true *and false* rows remain. CoreOutcomes has Monitor,
Product, InitialCanComplete, Nonblocking for all nine M in roster order.

Negative has Comparison (RF/doomed index), LeftProduct, RightProduct, LeftPath,
RightPath. A path has Start, Steps, End; steps have From, Label, To and must be
actual adjacent product arcs. Paths are shortest BFS with stored arc order.
Left ends at K=`W,U,true,false`, M=1; right at the same K and M=3. Left can
complete and right cannot. The doomed discarded bad residual and its strict
supersets are recorded in that comparison's RightWithin and RightFamily MAX
records at the reserve prefix. Initial and MAX equality and opposite NB are
certified separately; no synthesized runtime history is asserted.

LegalControl has Plant, Monitor, Roots, Products, Inclusions. Plant is an exact
copy of core K except ID suffix `/all-marked-control` and every state marked.
Monitor is minimal-augmented, Roots=["1","4","0"]. Products use this plant and
that monitor separately rooted and renamed `minimal-augmented@root`, with the
plant's exact signature. Inclusions compares product 1 to 4, then 4 to 0:
LeftSubset=true, RightSubset=false in both cases. Complete graphs and shortest
RightOutside certify both proper differences on legal K traces. This control
concerns legal traces, not paid completion; original plant marks are unchanged.

## Generic controls outside the monitor roster

Generics has MINLoss, JNonnecessity, Multiword, LocalHitNotNB. MINLoss and
JNonnecessity each have Left, Right, Comparison, Trie, LeftProduct, RightProduct.
Their generic automata use Signature={} and Role=generic. MINLoss is the contract's
marked s/unmarked d example; JNonnecessity the unmarked no-edge X/single-apply:AB
edge Y example. Trie is epsilon-trie accepting only epsilon. Products give
explicit example observations; the no-marks argument establishes the latter's
all-context claim, not the one trie. Generic descriptors preserve declared order.

Multiword has Automaton, SourceStates=["0","1"], TargetState="2", Coverage, Trie,
Products. Automaton has unmarked 0,1,2 and marked 3, with 0-apply:AB-3,
1-reserve-capacity:AB-3,2-commit:BC-3. Coverage is computed for the specified
source and target families using this automaton's residuals. H is the two-word
set in source order, both singleton words. Trie accepts exactly H; Products is
the three separately rooted copies in 0,1,2 order, renamed `ID@root`, synchronized
with Trie. Initial possible completion is true,true,false.

LocalHitNotNB has X, H=[["apply:AB"]], Trie, Product. X has q unmarked, f marked,
d unmarked and q-apply:AB-f plus q-apply:AB-d. Initial possible completion is
true but NB false. Trie construction starts at unmarked 0 and inserts H words
in recorded order, allocating decimal states on first encountered prefix and
marking precisely completed words. Empty H gives an unmarked no-edge root;
H=[[]] marks the root. Every generic product is complete. No generic graph is
included in roster/product totals.

## Counts and captures

Counts: Plants, Monitors, Controllers, Products, ProductStates, ProductArcs,
ProductCanComplete, InitialCanComplete, NonblockingProducts, Comparisons,
EqualLanguages, EqualFULL, EqualMIN, EqualMAX, EqualJ, ResidualComparisons,
LanguageSubsetPairs, ProfileSubsetPairs, Substitutions, NonblockingSubstitutions,
BlockingSubstitutions. Product totals sum roster graphs only. ResidualComparisons
counts cross and both within lists for every roster comparison; LanguageSubsetPairs
sums their Pairs plus each InitialLanguage.Pairs. ProfileSubsetPairs sums profile
Pairs. Generic and legal controls are excluded. Substitution counts distinguish
both NB true versus both NB false. All counts are measured cross-checks.

`tools/capture.py OUTPUT_DIR` refuses an existing output directory before any
write or command. A fresh directory retains deterministic gzip export bytes,
raw identity, exact command/results, dependency listing and its parsed
module-relative compiled roster, actual read-input identities, helper identity,
and an artifact manifest. Root/module information accompanies go-list evidence;
local source identities normalize relative to Module.Dir and Module.Path rather
than comparing against an incidental capture checkout. Two fresh captures must
have equal raw semantic bytes. No accepted semantic artifact is edited afterward.
RED snapshots and intentionally false hypotheses are separate from final captures
and passing tests. The independent Python auditor belongs to the next task.
