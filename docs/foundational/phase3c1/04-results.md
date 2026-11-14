# Phase 3C.1 results: contract preservation under composition

## Question and disposition

Can a two-owner component retain a useful progress contract when a third owner
is attached? **Conditionally, in the fixed chain/cycle grammar.** Retain each
original's full binding and U/W/P/D progress register, each physical owner's
availability and service policy, and every open endpoint dependency. Compatible
joins and closure preserve the named progress relation and the existence of a
finite concrete completion suffix under the stated transport/custody premises.
This is a direct source argument supplemented by finite checks, not a
mechanized theorem, unrestricted congruence or runtime qualification.

The packet also **falsifies** two smaller interface summaries. These failures
are useful results, not exceptions removed from acceptance. The
[contract](00-contract.md), [argument](01-argument.md),
[frozen schema](02-evidence-schema.md), [independent audit](03-independent-audit.md),
[qualification](05-qualification.md) and [review status](06-reviews.md)
separate the claims and evidence.

## What is preserved

The observation is the ordered sequence of original identity, full binding and
commit/application/consumption labels. It abstracts receipt values, transport,
recording, ACK/release and physical interruption/recovery to progress stutters.
Internal world applications remain visible. Complete means every declared
original is D; channel drain, receipt release and physical upness are separate.
This is a value-abstracted progress relation, not arbitrary payload equivalence.

For every reachable concrete representative s, each step projects to a named
abstract transition or a stutter. Conversely, every enabled abstract transition
from q(s) has a finite concrete realization with only stutters before the last
stage change, without advancing another original or repeating an application.
The argument handles all transition classes and uses authentic durable custody,
not an inaccessible world-result field. Down endpoints recover; requests can be
resent from commitments and results routed from actual retained copies. Thus
finite path projection and repeated lifting establish
`CanComplete_concrete(s) iff CanComplete_abstract(q(s))` under the premises.
Abstract eligibility does not assert immediate concrete readiness.

The context grammar is only A/B/C leaves, compatible binary joins of disjoint
owners and closure to chain A->B->C or cycle A->B->C->A. Internal registers and
policy incidence survive both parenthesizations. Union equality is the
algebraic part; the universal representative argument supplies the additional
concrete abstraction obligation. The claims do not hide application labels,
collapse availability, silently fill unresolved ports or erase original debt.

For each profile in this grammar, acyclicity of the static dependency graph is
sufficient for completion from **every reachable state**: finish the outgoing
prerequisites of serial callees first. It is not necessary for completion from
an individual state. The all-serial cycle can complete from its initial state,
but cannot complete from W/W/W. Possible completion is existential; no fairness
transfer, inevitable progress, divergence equivalence or matching physical
schedule is asserted.

## Measured evidence, with categories kept separate

| Evidence | Executed result | Scope |
| --- | --- | --- |
| Complete abstract graphs | Chain: 8 profiles, 128 states, 184 arcs, all 128 completable. Cycle: 8 profiles, 511 states, 1,053 arcs, 510 completable. | All reachable small abstract states/arcs, not a concrete-product enumeration. |
| Independent grouping reconstruction | 32 grouped relations agree with the 16 flat relations. | Both parenthesizations in the frozen grammar; not arbitrary contexts. |
| Factored Go probes | 311 old single-edge states; 24,258 representatives, 184,256 successor arcs, 14,058 enabled named-edge lifts. | Guard/availability crossings, including overlapping contexts and canonical peer states; not unique reachable-product counts. |
| Independent concrete ledger | Seven selected histories accepted. | Exact post-states, event provenance, projection and candidate collisions; not every concrete history. |
| Hostile controls | 35 separately retained records rejected for expected reasons. | Includes missing non-tree incoming arc, forged binding, erased debt, wrong guard, premature release and repeated application. |
| Retry witness | Exact equality of concrete states 3 and 5 after send/deduplicate. | A repeatable two-step divergence witness, not divergence equivalence. |

The totals are 639 abstract states and 1,237 arcs across profiles, with one
noncompletable abstract state. In the all-serial cycle, 62 of 63 states have a
completion suffix; only W/W/W lacks one. These are profile sums, not one merged
graph. Both exported result alternatives remain in the concrete semantics.

The [Go provenance](evidence/task1/source-manifest.json) records actual inputs
and commands. The [Python manifest](evidence/task2/manifest.json) records two
fresh exports byte-identical to each other and the Task 1 sample: 1,584,501 raw
bytes, SHA-256
`48b4cfb4cc10b7dde4979874f669014f4863e1e05e48737ef4561f119dfc3602`.
The Python checker independently regenerates the abstract relation and uses a
declarative concrete ledger. It shares the contract/schema and supplied old
protocol interface, not the new Go transition implementation. Such algorithmic
independence does not prove a shared specification or arbitrary prose true.
Concrete CanComplete verdicts depend on the separately argued abstraction;
there is no independent full concrete-product reachability computation.

## Preserved failures and assumption boundaries

- **R0, erased internal original:** equal full exposed BC transport and equal
  availability, with AB at W versus D after BC completes. Both states can
  complete, but Complete and internal application/consumption history differ.
  All custody premises hold and there is no circular wait. Erasure loses a
  declared observation.
- **R1, hidden-debt/component-policy bit:** equal exposed BC/CA transport and
  availability, CA admitted, A serial/up, with AB at U versus W. Both owe AB,
  but immediate CA application readiness differs. This refutes that summary
  for readiness; it alone does not refute a weaker trace relation that omits
  readiness. Named stage plus admission/availability supplies the missing facts.
- **All-serial W/W/W:** intact custody and authority do not discharge circular
  service promises. This is a compatibility obstruction, not custody loss.
- **Outside the premises:** permanent owner loss, split application/retention,
  deleted durable custody, changed binding and depleted shared resources do not
  falsify a theorem assuming their absence. They remain extension obligations.

R0/R1 are deterministic reachable witnesses; no globally shortest prefix or
minimal topology is claimed. Earlier test candidates and audit failures remain
in their original task evidence. Reviews do not promote these checks into a
universal result. Source inspection during reconciliation found no unmet case
in the lifting construction and required no model or argument change.

## Prior art, manuscript and remaining obligations

Abadi and Lamport's *Conjoining Specifications* supplies established context for
composition, hiding and assumption discharge. It is background, not a general
theorem instantiated by this packet. The
[primary publication record](https://www.microsoft.com/en-us/research/publication/conjoining-specifications/)
gives TOPLAS 17(3), May 1995, 507-534. The consulted
[author revision](https://lamport.azurewebsites.net/pubs/abadi-conjoining.pdf)
is dated 3 November 1995 and prints 507-533 in its article footer. The bibliography
uses the publication metadata and records the author-copy URL; no DOI was
verified or added. See [source-check notes](evidence/manuscript/prior-art.md).

The Markdown/LaTeX manuscript adds the same Phase 3C.1 section, updates the
Abstract and Conclusion, and makes the closure subsection's earlier open-gate
statement historical. The roadmap records this bounded result while preserving
prior evidence. PDF build, alignment, inherited differences and visual checks
are recorded in [qualification](05-qualification.md).

No production/runtime source, dependency, earlier model or old evidence changed;
no new runtime campaign, retained SQLite open or runtime `.envrc` access occurred.
Shared-authority fan-out, arbitrary contexts/topologies, dynamic grants or
identity, cancellation, split effects, independent-owner runtime recovery,
reverse runtime refinement and reduced verification-cost claims remain
unestablished. Phase 3C remains open. [Delivery](07-delivery.md) and
[final review](06-reviews.md) are recorded separately from the research result.
