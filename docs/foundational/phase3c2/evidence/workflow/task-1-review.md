# Task 1 independent review

Reviewed base `072cf44cc89879d60bf49233d3d4dc541209df0d` to head
`60d46a43dd32caa07f0b814e6741213f506a9cb6`, using the supplied complete diff,
Task 1 brief/report, frozen `00-contract.md`, and repository AGENTS.md.

**Spec compliance: ✅ for the declared Task 1 boundary.**

**Task quality: APPROVE, with one nonblocking documentation finding.**

This approves the model, concrete adapter, source argument and frozen export.
It does not approve Phase 3C.2 closure, an independent checker result, manuscript
reconciliation, or universal behavior beyond the source argument's premises.

## Question and result

The reviewed question is whether exact authorization and reusable capacity
conservation can coexist with loss of finite completion in the fixed AB/BC
chain, and whether original identity and stable reservation readiness preserve
the relevant completion distinction. The implementation and stated conditional
argument agree on that boundary. No blocking semantic or evidence defect was
found in this review.

## Spec checks

| Requirement | Assessment |
| --- | --- |
| Research-only extension; old model/runtime/dependencies/evidence unchanged | ✅ Supplied complete diff adds only the new package, argument/schema and Task 1 evidence. |
| Fixed generation-one chain, all eight masks, three layouts/policies | ✅ Constructors and exact config validation preserve these records; export contains the full 72-member roster. |
| Permission, cumulative E and Q distinct | ✅ Full old bindings persist; E sufficiency is an explicit premise, not an implemented refund/accounting claim. |
| Original-held identity, exact roots/memberships, no cap coercion | ✅ Validation retains the exact layout records; occupancy is derived from holders. |
| Atomic application/result retention/Q release; no crash refund | ✅ Only a permitted application clears a hold; old recovery/transport steps retain it. |
| Complete finite graph and existential CanComplete | ✅ BFS retains all arcs, including existing-node targets, and reverse traversal begins at both-D nodes. |
| Required adverse holder pair, ready-now invalidation and controls | ✅ Genuine paired histories are constructed through successors; full base states and pool totals are compared. Tests cover controls throughout the roster. |
| Stable-ready restriction on the old AB-first trace | ✅ Concrete witness, abstract negative test and argument retain the visible BC-progress cost. |
| Concrete finite lifting, physical availability, actual custody | ✅ Preparation recovers/routs custody without stage/hold progress; apply requires admission and a hold; consume uses actual retained/durable receipt paths. |
| No operational WorldResult oracle | ✅ New guards/lift do not inspect it; checked delegated validation/guard/successor code likewise has no oracle decision. |
| Explicit factored probe domain, every enabled action, projection/adjacency | ✅ Domain and checks are explicit in probe.go and the schema; exhaustive concrete-product coverage is expressly disclaimed. |
| Two regroupings and invalid joins | ✅ Static ownership/ports and exact capacity memberships are checked, equal root references are identified once, and both closed relations agree. |
| Frozen deterministic export, genuine full witnesses, provenance | ✅ Schema documents the records and task boundary; archive integrity and all 32 source hashes independently matched in this review. |
| RED evidence and required verification | ✅ Retained logs contain assertion failures and later ordinary/race/Staticcheck/export success; report explains the RED shell-status limitation and cache failure. Those executions are attributed to the implementation task. |
| Independent checker/manuscript/PDF | ✅ Correctly left as later-task obligations rather than claimed complete here. |

## Strengths

- `research/capacitycomposition/model.go:68` validates complete fixed records
  rather than accepting aggregate-equivalent resource configurations. The
  reservation relation preserves the named original and its full binding.
- `research/capacitycomposition/concrete.go:50` wraps the existing transport
  relation with one explicit application restriction and atomic hold clearing.
  Capacity does not become a prerequisite for receipt consumption.
- `research/capacitycomposition/probe.go:63` describes the synthesized domain
  precisely, and `probe.go:141` checks full projection, every enabled action,
  exact adjacency, preparation stuttering and oracle-erasure independence.
  Reported multiplicity is not presented as unique global reachability.
- `research/capacitycomposition/evidence.go:87` builds witnesses from empty
  initial custody using real transitions. Hostile tests attack initial receipts,
  holder substitution, premature release, bindings, actors, activation and
  invented post-states.
- `docs/foundational/phase3c2/01-argument.md:77` distinguishes source reasoning
  from finite checks, states recoverability and atomicity premises, and explains
  both directions of finite completion correspondence. Stable-ready is stated
  as a sufficient restriction, with its loss of old progress traces retained.

## Findings

### Critical

None.

### Important

None.

### Minor

**M1 — Copy the verified explicit timeout into the README race command.**
`research/capacitycomposition/README.md:12` gives `go test -race` without
`-timeout 30m`, although `docs/foundational/phase3c2/evidence/task1/race.log:20`
records 627.29 seconds for the probe test alone. The successful recorded command
at `race.log:1` explicitly raises the timeout. The README command can therefore
time out on the already-recorded environment under Go's default ten-minute
package timeout. Add the verified timeout to that reproducibility example.
This does not invalidate the retained passing race run or the semantic result.

## Review evidence and limits

Read the complete supplied diff in bounded sections. A long repeated probe-log
line caused output truncation; the omitted log/source-manifest segment was then
read separately. No changed source needed an additional full-file read.

Outside-diff inspection was limited to concrete identified risks:

1. The new wrapper delegates validation and transitions. Inspected the relevant
   `contractcomposition` concrete validation/projection/successor functions,
   service guard and component validation, plus `twoagent.Successors`, to check
   that delegated guards did not read WorldResult, that durable custody survives
   interruption, and that physical availability and fixed bindings actually
   support the lift construction.
2. The binary export is opaque in a textual diff. A fresh Python standard-library
   check decompressed it, checked packed/raw byte lengths and SHA256 values,
   compared all 32 exported source hashes to current files, and inspected roster
   sizes and aggregate counts. It read only the declared archive/manifest/source
   paths, with no SQLite or environment files.

Fresh archive check result: 72 graphs, 72 groupings, 10 witnesses; 1,632 states,
2,412 arcs, 1,624 CanComplete and 8 CannotComplete labels. Probe packet totals
are 1,854,036 representative cases, 11,972,664 micro arcs, 2,451,168 lifts and
4,833,504 lift microsteps. These are packet observations, not a fresh operational
probe rerun or independent semantic reconstruction.

No Go suites, race runs, Staticcheck runs or exporter reruns were performed by
this reviewer; their retained passing logs were reviewed instead, as requested.
No focused new semantic probe was necessary after source inspection resolved
the delegated-custody/oracle doubt. No source/index/branch mutation or subagent
was used; this report is the only review output written.

Declined-to-verify items and resolution:

- **Independent abstract re-enumeration and declarative witness replay:** Task 2
  must supply its standalone checker and hostile-input evidence. Go self-checks
  and this source review do not substitute for that work.
- **Whole concrete product or a mechanized universal proof:** not provided by
  this task and not inferred from factored probes. A stronger claim would need
  corresponding proof or coverage; the current conditional source argument is
  appropriately labeled.
- **Paper/LaTeX/PDF consistency and inherited warnings:** later-task changes and
  rebuilt/inspected PDF evidence must resolve these before overall completion.
- **Full unchanged repository suite reuse:** not independently claimed here.
  Any later reuse must bind historical results to exact unchanged source and
  dependency identities.
- **RED execution chronology:** reviewed retained failure output and the
  implementation report's attribution; did not reconstruct preimplementation
  states from history. Stronger archival chronology would require preserving
  those source snapshots and command exit statuses, but it is not needed to
  accept the present Task 1 semantic result.
