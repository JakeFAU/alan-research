# Phase 3C.3 independent interface and history audit

## Question and result

Do the retained finite automata, products and concrete model histories support
the stated distinction between initial completion-preserving substitution and
preservation of every residual's ability to complete?

**Bounded result:** the independent Python reconstruction agrees with the
frozen producer export across the entire declared abstract roster and all eight
concrete histories. The common-cut counterexample is genuine in the model:
commit-first and release-first have identical complete old state, signature and
holdings at the cut, yet opposite possible-completion labels. The full product
relation determines these labels, not a stalled sample schedule.

The three release-first variants have equal successful visible languages. The
doomed variant nevertheless has a reachable bad residual in the core context.
Initial language equality therefore supports the stated initial existential
completion comparison; it does not establish residual equivalence or universal
nonblocking. This audit is finite model evidence, not runtime qualification,
mechanized proof or complete-agent substitution.

## Independence and method

`tools/audit.py` consumes the contract and frozen schema. It does not import,
execute or mechanically translate Go transition, lifting or explorer helpers.
The producer is invoked once solely to capture bytes. Go dependency source
files are hashed without reading their operational definitions as an oracle.

The new Python code independently constructs:

- The plant's stage/hold relation from immutable service and reservation rules,
  with each original's root occupancy derived from exact membership.
- All five monitor fixtures and all 324 role-restricted controller descriptors,
  including unreachable controller states and every external-action self-loop.
- Complete synchronized products, all relation arcs, independent private steps,
  conjunctive marks, and reverse reachability from marks.
- Both actual binary parenthesizations for every triple case, with explicit
  flattening of intermediate tuple states. Equality includes arcs, marks and
  reverse labels; it is not equality of totals alone.
- Epsilon closures and the complete reachable subset-pair graphs. Their BFS
  supplies the shortest alphabet-order acceptance separator for unequal pairs.
  Tests separately exercise epsilon cycles, nondeterminism, empty words and
  empty subsets, and enumerate every shorter word for the retained separators.
- Every equal-language substitution row in both rosters, including blocked
  contexts. The comparison is explicitly of initial possible completion.

The unchanged checker at `../phase3c2/tools/audit.py` is reused for inherited
configuration/binding construction, exact JSON typing, and concrete transport /
physical custody replay. Its SHA-256 is
`1dc485fb2e3f0594b6122daa3a63ec05a12c8f5a6b9988afe5d2a775a768f193`.
This is attributed historical Python source reuse, not an independently
invented specification or a fresh rerun of the old full qualification campaign.
The new plant reconstruction does not call that checker's abstract transition
function; inherited concrete replay retains its existing capacity guards.

Concrete replay starts with genuinely empty initial custody. Each event derives
its full post-state from the preceding state, including actual application
result provenance, physical eligibility, exact holds and retained monitor PC.
New private steps require B up and can change only PC. Every projected visible
or private transition is checked against the independently reconstructed
product. Transport and physical events must stutter in that projection. Cuts
are recomputed from actual replayed states. Recovery schedules, common-cut
identity, completed suffixes and the blocked final states are checked explicitly.
The eight histories are model executions, not endpoint or runtime captures.

## Measured finite totals

These are configuration sums from independent reconstruction. Grouping graphs
are checked separately and are not added again to the product totals.

| Measure | Reconstructed value |
| --- | ---: |
| Plants / monitors / controllers | 72 / 5 / 324 |
| Products | 1,980 |
| Product states / arcs | 19,544 / 23,420 |
| Product states with a completion suffix | 9,526 |
| Products initially able to complete | 622 |
| Products nonblocking at every reachable state | 392 |
| Language comparisons / equal pairs | 10 / 3 |
| Reachable subset pairs across comparisons | 43 |
| Triple regroupings, each with two reconstructed parenthesizations | 1,620 |
| Substitution rows | 1,188 |
| Successful / blocked substitution rows | 396 / 792 |
| Concrete histories / events | 8 / 123 |

The equal pairs are release-first versus release-first-tau, release-first versus
release-first-doomed, and release-first-tau versus release-first-doomed. Successful
and blocked substitution totals count pair/context rows, not unique contexts.

## Captures and hostile controls

`evidence/task2/capture.py` records argv, cwd, explicit environment overrides,
stdout/stderr and exit code. Each output is created exclusively; rerunning the
capture script refuses to overwrite the first capture. The fresh source-bound
export is byte-identical to Task 1: raw SHA-256
`a073a22fe8ac62e05f2d15b0f5146b19db9606e38ce4619261e7e4108e27269d`,
25,384,964 bytes. The retained gzip reproduces the exact accepted plain bytes;
the plain CLI input is also staged under `/tmp/alan-phase3c3`.

The seven unittest cases include 53 semantic hostile mutations. Separate CLI
captures cover those same 53 inputs plus duplicate JSON keys and nonfinite JSON
numbers: **55 exact rejected inputs**. Both valid plain and gzip CLI inputs
are accepted. No hostile input was initially accepted. The initial missing-checker
failure and first successful hostile unittest output remain in the evidence.

The controls challenge omitted plant gates, forged language equality, invalid or
nonshortest distinguishing words, corrupted subset acceptance/arcs, omitted tau,
forged marks, external-action vetoes, hidden effects, incorrect actors/bindings,
wrong roots and signatures, omitted/duplicated descriptors and edges, false
product/grouping/substitution completion, invented receipts, application or
recovery post-states, illegal private steps, wrong PC/activation/hold recovery,
forged cuts, unknown fields, invalid types and forged source binding.

The non-tree-arc attack retains every original BFS discovery-parent edge and
all original reachable nodes. It corrects the graph arc count and aggregate
product-arc count. `non-tree-control-check.json` retains those parent edges and
the exact omitted arc. The checker rejects the missing relation arc even though
reachability parents and superficial counts remain consistent.

Evidence indexes and binding:

- `hostile-index.json`: exact input hashes, rejection exits and reasons.
- `accepted-index.json`, `audit-plain.command.json`, `audit-gzip.command.json`:
  accepted input identities and command results.
- `fresh-export.command.json`, `fresh-export-comparison.json`: fresh capture and
  byte comparison against the retained Task 1 capture.
- `unittest.command.json`, `independent-counts.json`: test output and reconstructed
  totals.
- `checker-sources.json`: the new checker, tests, capture script and reused
  Phase 3C.2 Python source hashes.
- `manifest.json`: SHA-256 and size of every retained Task 2 evidence artifact.

## Scope and remaining obligations

This audit establishes agreement over the finite roster and verifies these
specific concrete histories. It does not exhaust the concrete product or prove
the source-level finite-lifting argument for every concrete representative.
The conditional language-intersection argument in `01-argument.md` remains a
separate mathematical/source argument under the declared synchronization,
private-state, signature and marking assumptions. The checker and tests are
fallible implementations; agreement is evidence, not universal proof.

No production code, runtime execution, live endpoint, full-repository test run,
old million-probe campaign, arbitrary topology, dynamic authority, split-effect
unknown outcome or distributed allocator was exercised here. This task changes
only audit tools and audit evidence. Task 3 retains responsibility for manuscript,
roadmap, PDF, scoped Go qualification and whole-packet reconciliation. This audit
does not claim general Phase 3C closure.
