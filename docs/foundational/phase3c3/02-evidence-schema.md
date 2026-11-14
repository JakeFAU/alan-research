# Phase 3C.3 evidence schema: `alan/phase3c3/v1`

The producer's sole exported Go API is `Export(io.Writer) error`. Run from the
repository root with `go run ./research/componentinterfaces/cmd/export`. Standard
output is one deterministic JSON object plus a newline; diagnostics use stderr.
All fields below are required; unknown/extra semantic records are invalid. Field
names use the exact capitalization shown. Integers are JSON integers. Arrays
represent complete finite sets/relations unless an order is explicitly defined.
No timestamp, host-specific path or Git commit is part of the semantic export.

## Top level

`Schema`, `Alphabet`, `Monitors`, `Plants`, `Controllers`, `Products`,
`Comparisons`, `Groupings`, `Substitutions`, `Histories`, `Counts`, `Sources`.

`Schema` is `alan/phase3c3/v1`. `Alphabet` has eight records in original order
AB, BC, each with kind order commit, reserve-capacity, apply, consume. Each
record has `ID` (`kind:original`), `Action`, `Actor`, `Root`. `Action` has exact
old `Kind`, `Original`, `Binding`; `Binding` has `Original`, `Caller`,
`CallerGeneration`, `Operation`, `Callee`, `CalleeGeneration`, `Request`, `Grant`,
`Resource`. The root/actor here use the core early/shared-1/mask0 signature.
Other plant roots are derived from their exact `Signature`.

`tau` is a separate edge label and has no alphabet entry or binding.

## Signatures and automata

`Monitors`: the five frozen M descriptors in contract order, instantiated with
the core signature. `Plants`: all 72 K descriptors, layout order shared-1,
shared-2, separate-1; policy order early, ready-now, stable-ready; mask 0..7.
`Controllers`: all 324 D descriptors ordered by code 0..80, then mark mask 0..3.

Each descriptor has `ID`, `Role`, `Signature`, `Initial`, `States`, `Edges`.
`Role` is K, M or D. `Signature` is the complete old capacity `Config`:
`Base` (`Topology`, `ServiceMask`, `Agents`, `Originals`), `Layout`, `Policy`,
`Roots`, `Memberships`. Each agent has `ID`, `Generation`, `Servicing`; each
root has `ID`, `Limit`; membership has `Original`, `Root`. Exact original
bindings are retained; no missing, extra, reordered or conflicting records are
accepted. Joining factors never merges or sums roots. Monitors are instantiated
with the selected plant's exact signature for the 360-case roster; their graphs
remain identical. D uses only the core signature.

A state is `{ID, Tuple, Marked}`. Base descriptors have `Tuple: null`; the ID
is their atomic stable source key. M/D keys are decimal state numbers. K keys
are `stageAB,stageBC,heldAB,heldBC`, where stages are U/W/P/D and booleans are
lowercase true/false (example `W,U,true,false`). K ID is
`policy/layout/maskN`. M IDs are the exact contract names. D ID is `D-CC-M`,
with two-digit zero-padded decimal code CC and mark mask M. Decode four base-3
digits least-significant first in order `(0,commit:AB)`, `(0,apply:BC)`,
`(1,commit:AB)`, `(1,apply:BC)`; 0 absent, 1 target 0, 2 target 1. Mark bit 0
marks state 0; bit 1 marks state 1. All six other labels self-loop.

Each edge is `{From, Label, To}`. Every declared edge and mark is retained,
including unreachable D states. M's five unrestricted labels and D's six
unrestricted labels must be self-loops at every state, with no other target on
those labels. Unknown labels, duplicate edges/states and undeclared endpoints
are invalid. Tau changes only one factor's private state.

## Products and grouping

`Products` has 360 K/M graphs followed by 1,620 core K/M/D graphs (D outer,
M inner). A graph has `ID`, `Factors`, `Initial`, `Nodes`, `States`, `Arcs`,
`CanComplete`, `Nonblocking`. ID joins factor IDs with ` || `; factors retain
order K,M[,D]. Every reachable node and every relation arc is retained. A node
has `ID`, `Tuple`, `Marked`, `CanComplete`, `Arcs`; each arc has `Label`, `To`.
Node ID is the compact JSON encoding of its flat tuple of atomic source keys;
`Initial` and arc `To` refer to these IDs. Tuple components retain K,M,D order.
Nodes are BFS-discovered, starting at index zero; arcs sorted by Label then To.
Duplicate relation arcs are coalesced. Mark is conjunction of factor marks;
CanComplete is reverse reachability from all marked nodes. Graph CanComplete is
a count, States/Arcs are counts, and Nonblocking means every reachable node can
complete. Initial possible completion is `Nodes[0].CanComplete`.

`Groupings` contains 1,620 records `{Product, Left, Right, Equal}`. Product
identifies the flat triple graph. Left/Right are **full** graph objects produced
by `(K||M)||D` and `K||(M||D)`. Their Factors retain the intermediate factor ID,
but all node tuples and IDs are explicitly flattened to the same K,M,D keys.
Equality compares initial state, all normalized nodes, arcs, marks and reverse
labels, independent of BFS node order. A checker must reconstruct both products
rather than trust Equal or counts.

## Exact languages and substitution

`Comparisons` contains all ten unordered M pairs in contract order. Each record
has `Left`, `Right`, `Equal`, `Word`, `LeftAccept`, `RightAccept`, `Pairs`.
Pairs is the complete BFS of epsilon-closed subset pairs, with all eight arcs
per pair, even after inequality is found. A pair node has `Left`, `Right`
(sorted arrays of state IDs, including `[]` for empty subsets), `LeftAccept`,
`RightAccept`, `Arcs`. Each arc has `Label`, `To` (zero-based pair-node index).
Initial pair index is zero. On inequality Word is a shortest acceptance
separator (BFS, alphabet-order tie break), and LeftAccept/RightAccept are its
actual acceptance outcomes. On equality Word is `[]` and both top-level
acceptance fields are false sentinels, not claims about epsilon acceptance.
Pair-node acceptance fields always denote actual acceptance.

`Substitutions` contains one check per equal-language pair per context (72 K
contexts and 324 core K/D contexts). A record has `Left`, `Right`, `Context`,
`ProductLeft`, `ProductRight`, `CanCompleteLeft`, `CanCompleteRight`.
Context is K ID for the first roster and `K ID || D ID` for the second.
The booleans refer to the **initial** product nodes. Both blocked and successful
contexts remain. No assertion compares arbitrary residuals or nonblocking.

## Concrete histories

`Histories` contains the eight contract histories in roster order, with IDs
`history-1` through `history-8`. Each has `ID`, `Config`, `Monitor`, `Initial`,
`Steps`, `Cuts`, `Complete`, `CanComplete`. Config is the full old configuration;
Monitor references the corresponding M descriptor, instantiated with Config.
A concrete state has `Old`, `Monitor` (retained PC as decimal string). Old is the
unchanged capacity ConcreteState: `Base` and `Held` (two booleans). Base has
`Agents` (three availability records) and `Originals` (three transports, third
unused). Availability fields: `Down`, `Crashed`, `Activation`. Transport fields:
`Binding`, `Committed`, `Admitted`, `Applied`, `Consumed`, `Volatile`, `Recorded`,
`Retained`, `RequestChannel`, `ResultChannel`, `ACKChannel`, `WorldResult`.
Values are empty string, r1 or r2. Empty undeclared bindings have all strings
empty and numbers zero. Initial is captured from NewConcrete with PC zero.

A step is `{Event, Private, State}`; State is the complete actual post-state.
Event uses old fields `Kind`, `Actor`, `Original`, `Binding`, `Result`,
`Duplicate`. Private tau uses Kind `tau`, Actor B, empty Original/Binding/Result,
Duplicate false, Private true. It requires B up and changes only Monitor.
All inherited events have Private false. Commit/reserve/apply/consume advance
M on the matching visible label; transport/physical stutters retain M.
Application's r1/r2 is genuinely generated by the old relation; histories
choose the old Lift's deterministic first legal result. No receipt is invented.

Cuts have `Name`, `AfterSteps` (number of events already taken, zero is Initial),
`Node` (matching K/M graph key), `Complete`, `CanComplete`. First cut is
`after-reserve-AB` for every history. History 6 adds `B-down-before-private` and
`private-after-recovery`; history 8 adds `doomed-before-crash` and
`doomed-after-recovery`. All end in `final`. Histories 1/2 have equal Config and
entire Old state at their first cuts, but opposite CanComplete. History 1 ends
after BC commitment. History 8 ends in doomed state 3 after crash/recovery.
Other histories end marked. Cut and final labels must be independently derived
from the full product relation, never inferred from a stalled sample history.

## Counts and source binding

`Counts` has `Plants`, `Monitors`, `Controllers`, `Products`, `ProductStates`,
`ProductArcs`, `ProductCanComplete`, `InitialCanComplete`, `NonblockingProducts`,
`Comparisons`, `EqualComparisons`, `SubsetPairs`, `Groupings`, `Substitutions`,
`SuccessfulSubstitutions`, `BlockedSubstitutions`, `Histories`, `ConcreteSteps`.
Graph state/arc sums count each configuration separately; grouping graphs are
not added again. Successful/blocked substitution totals count equal-language
pair/context rows, not unique contexts. Measured values are not schema axioms.

`Sources` is an ordered list `{Path, SHA256}` of actual Go producer dependency
source files and relevant module files. Paths are repository-relative, hashes
lowercase SHA-256 of exact bytes. It includes every non-test Go file in the
producer and imported local model packages (including unused declarations in
compiled source files), root go.mod/go.sum, and the frozen contract, argument
and this schema. Tests and mutable execution/review metadata are excluded.
The separate artifact manifest binds exact raw and deterministically gzipped
export bytes and records the second identical capture. Verification logs retain
exact argv, explicit environment, stdout/stderr and exit codes.
