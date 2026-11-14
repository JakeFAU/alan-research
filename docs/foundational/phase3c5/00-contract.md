# Phase 3C.5: completion languages and residual antichains

## Question and boundary

Can the full residual-language profile of Phase 3C.4 discard redundant
completion alternatives while preserving both initial possible completion and
contextual nonblocking, and remain closed under the same composition operator?
What exact local sense of necessity can be justified for the reduced profile?

Base: merged Phase 3C.4 `347d93ba25a5f7a10f75144b1a44107be09a2f39`, with the
same tree as delivered `358fb222dec0e00148b3f2a804531c5b31202342`.
This packet replaces finite B continuation monitors over the unchanged whole
Phase 3C.2 plant K. It does not extract an entire independently executing agent.
New Go code is confined to `research/residualantichains/` and its exporter.
Do not modify old model, runtime, module, test or evidence files.

Controller choice under the user's continuing research authorization: retain
the initial marked language and the inclusion-minimal residual languages at
every visible prefix. Full profiles retain redundant supersets; maximal
profiles discard the hardest completion requirements. A general conflict
preorder decision procedure is a larger question. This packet does not claim
the weakest global contextual relation, a minimal automaton, a faster decision
procedure, or a new general process-theory result.

## Semantics and mathematical obligations

Use finite pointed marked epsilon-NFAs with exact compatible signatures and a
common visible alphabet. All visible steps synchronize all factors; private
tau changes only its factor's private state without inspecting another factor.
Marks are conjunctive; finite-path acceptance includes epsilon and leading and
trailing tau. No event priorities, fairness, inevitability or divergence
preservation are assumed. Shared authority, effects and Q remain in one
unchanged K. These restrictions are load-bearing hypotheses.

Retain Phase 3C.4's weak reach set R_X(w), individual residual marked language
L_X(x), and profile P_X(w) = {L_X(x) : x in R_X(w)}. For a finite family F of
languages define A(F) as its distinct inclusion-minimal members: collapse
equivalent languages, then discard proper supersets. Retain the empty family
and a family containing the empty language as different objects. Define

    A_X(w) = A(P_X(w)),
    J(X) = (Lm(X), w ↦ A_X(w)).

Prove as conditional source arguments, with exact quantifiers and both
directions where appropriate:

1. Every member of finite F contains a member of A(F). For every suffix
   language H, all F members intersect H iff all A(F) members intersect H.
2. With F tensor G = {L intersect H : L in F, H in G},
   A(F tensor G) = A(A(F) tensor A(G)). Empty intersections remain members.
3. Using the inherited Cartesian reachability and residual intersection
   lemmas, equality of A at every word preserves nonblocking in every
   compatible finite context and is itself closed under product.
4. Lm(X||Y) = Lm(X) intersect Lm(Y), so J equality preserves initial
   CanComplete and nonblocking and is closed under the same product. State
   induction only for this finite synchronized grammar, with one K.
5. For T_F(H) := every L in F intersects H, and finite F,G:

       [for every finite H, T_F(H) implies T_G(H)]
         iff [for every K in G, some L in F satisfies L subseteq K].

   When coverage fails at K, choose one word in L minus K for every source
   language L; their finite set H hits every source and misses K. An empty
   source family yields H empty. Mutual coverage is exactly equality of
   minimal families. This characterizes local finite-suffix hit tests, not
   whole-system contextual necessity. Explain the finite-family hypothesis.

The MIN profile alone loses initial marked language in the generic marked
domain. Preserve the counterexample: X has marked initial s, unmarked d,
tau s->d, and every visible self-loop at both states; Y is one unmarked state
with every visible self-loop. Their MIN profiles are {empty language} at every
word but their initial languages are Sigma* and empty. J distinguishes them.
All-marked leaf monitors make the initial component redundant, but composites
with K have nontrivial conjunctive marks, so that shortcut cannot govern the
induction. Do not silently narrow the theorem to leaves.

J is not necessary for general global contextual observations. Retain the
unmarked no-edge X / single-a-edge Y example: equal empty initial languages,
different MIN profiles after a, and both initial CanComplete and nonblocking
false in every pointed conjunctive-mark context. This generic example is
outside the all-marked ALAN roster.

## Frozen input, fixtures and comparisons

The only public Go boundary is `Export(io.Writer) error`, invoked from the
repository root by `go run ./research/residualantichains/cmd/export`.
Emit deterministic schema `alan/phase3c5/v1`, no timestamps or temporary paths
in semantic output. Internal representations are builder-owned.

Read `docs/foundational/phase3c4/evidence/task1/export.json.gz`. Verify schema
`alan/phase3c4/v1`, compressed length 157341 and SHA-256
`35115c3104cddd8a789c80a8c0e73eb7c7f173632af512776ba8502b61b9afc9`, raw length
8271854 and SHA-256
`b40f1d5ac0aea0f43999502ea318cb128f545eb8fb1fd129451915a5aa84bf1f`.
Verify all 40 predecessor source/provenance records against current files.
Import only its alphabet, signatures, seven M, 72 K and 324 D descriptors.
Old products, comparisons, verdicts and counts are not oracles. Preserve exact
roots, memberships and action bindings; explicit rebinding uses the whole
inherited plant signature. No dependency or general external-input API.

Add two all-marked B monitors `minimal-base` and `minimal-augmented`. Let
r=reserve-capacity:AB, a=apply:AB, c=commit:BC; O is the other five labels.
Every state has each O self-loop. Unlisted r/a/c edges are absent.

| Monitor | States / initial | Additional edges |
| --- | --- | --- |
| minimal-base | 0,1,2,3 / 0 | 0-tau->1; 0-r->3; 0-a->3; 0-c->3; 1-r->2; 2-a->3; r/a/c loops at 3 |
| minimal-augmented | 0,1,2,3,4 / 0 | all base edges; 0-tau->4; 4-r->2; 4-c->3 |

Expected source argument: L1 is properly contained in L4, which is properly
contained in L0. FULL profiles differ at epsilon, but MIN and initial language
agree at every word. Both proper differences intersect actual legal core-K
traces; certify this by marking every core-K state for this control only.
This control concerns legal traces, not paid completion. Design reconnaissance
found 52 nonblocking and 20 blocking K contexts per new monitor, with all 324
K/D contexts blocking. Preserve those blocked outcomes. A rejected alternative
fixture had more productive contexts but only plant-disabled differences;
archive the design evidence rather than substituting it to improve counts.

Enumerate nine M, all 36 unordered pairs, and all 3564 products: 72 K x nine M
plus core early/shared-1/mask0 K x 324 D x nine M. Compare initial language,
FULL, MIN, MAX and J separately. Keep all existing monitors and awkward
contexts. Nonblocking means all reachable nodes have CanComplete=true.
All other counts are measured, not declared expectations.

## Exact evidence and adverse controls

Retain complete epsilon-closed subset-pair graphs for directed language
inclusion, both within and across monitors. Inclusion fails at left-accept
and right-reject, not XOR. Retain shortest directional counterwords in fixed
alphabet order; distinguish absent witness from epsilon. Complete the graph
after a failure. Compare all weak-reach subset pairs at every visible word,
including empty subsets. Canonical representatives collapse equal languages
before strict MIN/MAX filtering; retain classes and proper dominators.

At each failed local coverage retain an uncovered target residual, shortest
word outside it for every source language class, and the deduplicated finite
H. Check every source hit and target miss, including H empty. A suffix-set
trie only tests initial CanComplete of each separately rooted residual; it
does not automatically distinguish global nonblocking. Preserve the explicit
control q-a->marked f and q-a->unmarked dead d with H={a}: initial hit true,
whole product nonblocking false.

Retain the RF/doomed MAX negative: all-word MAX and initial language equality,
opposite core nonblocking, the dominated bad residual, and actual core paths.
Test multiple-word H (F={{a},{r}}, K={c}, H={a,r}), inclusion direction,
epsilon, empty family, empty language, equivalent duplicate states, tau cycles
and strict versus non-strict dominance. Include finite independent algebra
checks of the product/minima law. Generic marked controls stay separate from
the all-marked monitor roster. Do not call local separators global contexts.

Keep complete product states/arcs/marks/reverse CanComplete labels and both
true and false J-equal substitution rows for all 396 contexts. Use exact
automata reconstruction, not shallow word sampling or counts-only checks.

## Capture, independent challenge and qualification

Source manifests bind every actual producer input, distinguishing compiled
Go files from hash-only provenance. Cross-check `go list -deps -json` using
recorded module/root information and repository-relative source identities,
not an accidental dependence on the capture checkout's absolute path.
Freeze argument/schema, capture the final deterministic export twice, and
keep byte identities, commands, outputs and exits. Capture helpers live in
`tools/`, require a fresh output directory, and refuse to overwrite captures.
Do not edit semantic evidence after capture. Preserve initial RED failures
and semantic falsifications separately.

The Python auditor reconstructs descriptors, inclusion, reduction, suffix
witnesses, all products and observations independently from Go. It may reuse
unchanged predecessor Python descriptor and automaton utilities, with exact
hash attribution; old outcome fields and Go predicates are never oracles.
New antichain/coverage/hit logic is independent. Retain exact valid and hostile
plain/gzip CLI inputs and results, including missing non-tree arcs with
superficial counts repaired, reversed inclusion, lost equivalent minima,
MAX substituted for MIN, empty-family confusion, missing/forged suffixes,
bad shortest witnesses, fabricated products/marks, source/signature changes,
duplicate keys and strict JSON type attacks. Preserve initially accepted
attacks before repairing the checker.

Run fresh scoped ordinary/race Go tests, scoped Staticcheck, checker tests,
valid/hostile CLI checks and root build/module/package checks. Do not repeat
unchanged suites merely for repetition. Attribute source-identical historical
qualification; no fresh whole-root suite or global lint-clean claim follows.
No runtime campaign, SQLite opening, old exhaustive transport-probe execution,
distributed/split-effect qualification or new physical history is required.

## Reconciliation and delivery

Synchronize the changed Markdown/LaTeX mathematics, Abstract/Conclusion and
roadmap; add verified primary background references only as warranted.
Rebuild and inspect the PDF, comparing fresh baseline warnings and pages.
Preserve the eleven-form count, corrected historical congruence wording, and
open Theorem 4.8. Do not claim all inherited manuscript differences resolved.
Archive independent review findings, corrections, workflow and source checks.
Deliver a narrow draft PR, retain the worktree and do not merge automatically.
General Phase 3C, full-agent substitution, dynamic authority, distributed
ownership, effects, exact costs, refusals and divergence remain open.
