# Task 1 independent review

**Spec compliance: APPROVE.**

**Task quality: APPROVE.**

Reviewed base `4b8a0b8b135e314660a2146401f494aad4525cdb` to head
`5f907cda569d97ec5c08923b35e77de5b3df06ce`, using the supplied exact diff
package, task brief, author report, binding review instructions and task-reviewer
template. This is the Task 1 gate, not the final branch review or approval of
the pending independent auditor/manuscript tasks.

## Spec compliance and strengths

- The requested model, composition, concrete adapter, exploration, tests, CLI,
  README, argument, schema and evidence files all appear in the supplied diff.
  Its 35 added paths are confined to the assigned research package and Phase
  3C.1 documentation/evidence. No old-model, runtime, dependency or manuscript
  modification appears in this task diff.
- Exact fixed identities, topology and mask validation are explicit in
  `research/contractcomposition/model.go:66`; the additional validated State
  and Action constructors at lines 203 and 213 reject undeclared records.
  Complete retains every original, and named serial eligibility uses the
  callee's outgoing W/P stage (`model.go:126`).
- Component validation, compatible union and closure retain originals and
  endpoint policies while rejecting duplicate ownership, incompatible
  descriptors, resource/binding deviations and unresolved ports
  (`research/contractcomposition/compose.go:34`, `:86`, `:126`). The flat and
  grouped configurations are explored separately; this is correctly not
  presented as an independent semantic oracle.
- BFS retains every outgoing arc, including arcs whose targets were already
  discovered. Reverse reachability uses the resulting complete incoming
  relation rather than the predecessor tree
  (`research/contractcomposition/explore.go:42`). Dependency construction and
  prerequisite-first ordering at line 90 implement the stated sufficient
  condition without treating cyclicity as a necessary obstruction at each
  state.
- The concrete adapter supplies one availability record and crash/recovery
  transition per physical owner, suppresses all roles of a down actor, and
  applies the serial guard only to Apply
  (`research/contractcomposition/concrete.go:163`). Retained-record validation
  and lifting use actual custody, not WorldResult. Lift's commit, apply and
  consume cases preserve other originals' stages and avoid reapplication
  (`concrete.go:214`). The unchanged twoagent relation supports each preparatory
  send/admit/receive/record/recovery used by those cases.
- The source argument states reachability and custody premises, handles every
  concrete transition class, constructs finite lifts from arbitrary reachable
  representatives, and gives both finite completion-path directions
  (`docs/foundational/phase3c1/01-argument.md:96`). Its grouping and dependency
  argument at line 153 is appropriate to the fixed grammar. I found no missing
  case that invalidates the conditional proposition. This is source-level
  review, not mechanized proof.
- R0/R1 histories, all-serial control, retry cycle and pre-record-loss/release
  history are constructed by legal concrete transitions and passed through the
  old separately implemented causal ledger
  (`research/contractcomposition/explore.go:131`,
  `research/contractcomposition/concrete.go:300`). R1 is explicitly limited to
  the declared readiness observation. Retry divergence is retained rather
  than removed to manufacture a progress result
  (`docs/foundational/phase3c1/01-argument.md:191`).
- The factored test covers all old edge states for AB/BC/CA renamings, callee
  stage/policy/availability cases and terminal chain BC
  (`research/contractcomposition/concrete_test.go:149`). Dedicated tests check
  shared-owner crash effects, independent-ledger negative cases and lifting
  oracle noninterference (`concrete_test.go:291`, `:306`, `:338`). The argument
  correctly labels the reported 24,258 contexts, 184,256 arcs and 14,058 lifts
  as factored checks, not full concrete-product enumeration.
- The versioned schema freezes named fields, identity slots, graph arcs,
  ordering, physical events, witnesses and claim categories
  (`docs/foundational/phase3c1/02-evidence-schema.md:15`). The exported claims
  keep source arguments, bounded graph checks and witnesses distinct
  (`research/contractcomposition/explore.go:310`).

## Issues

### Critical

None found.

### Important

None found.

### Minor

- **M1 — include context in factored-test failure messages.**
  `research/contractcomposition/concrete_test.go:217` and `:239` report generic
  failures inside the 24,258-context loop. These messages omit topology, named
  target, policy, physical phase and source/target states, making a future
  counterexample unnecessarily difficult to reproduce. Add those values to
  the failure messages, or use named subtests at the topology/edge/policy
  boundary. Existing assertions are substantive; this is a diagnostic-quality
  improvement, not an acceptance blocker.

## Actual checks performed

1. Read the supplied brief, report and review instructions; read the required
   task-reviewer template. Applied the Go test-review and code-clarity guidance
   within the explicit read-only/no-subagent scope.
2. Read every nongenerated diff hunk in manageable chunks, including the
   complete argument, schema, source, tests, README, command logs and manifest.
   The first large tool result truncated part of the schema; reread that
   exact schema diff range. No changed source file was independently reread.
3. Read `docs/foundational/phase3c1/00-contract.md` to resolve the precise
   positive-proposition and candidate-observation requirements.
4. For the named risk that the adapter's universal lifting relies on an old
   transport guard or retry rule it does not actually satisfy, read unchanged
   `research/twoagent/model.go`. It permits the requisite retained-data steps,
   retains retransmission/duplicate variants, gates caller actions on caller
   availability, and leaves callee work independent of caller availability.
5. For the named risk that the purported independent historical ledger shares
   producer transition logic or invents custody, read unchanged
   `research/twoagent/oracle.go`. It maintains its own causal variables, rejects
   repeated application/premature release and does not call Successors.
6. Ran one read-only `python3 - <<'PY'` structural/provenance probe over added
   records reconstructed from the supplied diff. It verified all 17 source
   SHA-256 entries (changed inputs from the diff; unchanged inputs from disk),
   every manifest command-output path, the 1,584,501-byte export's declared
   SHA-256, the exact 16-config roster, unique graph states, valid arc targets,
   predecessor/Via membership, and consistency of serialized aggregate counts.
   It also checked the seven witness records, equality of exposed full
   transport/availability in both candidate pairs, the exact retained retry
   state equality, and that the only nulls are the 16 initial-node Via fields.
   All assertions passed. This probe did not independently regenerate
   successors or perform Task 2's semantic audit.
7. Inspected the retained final normal/race outputs and RED/GREEN logs from
   the diff. Normal and race logs report success for both packages; final
   formatting and Staticcheck logs are empty, with exit zero recorded in the
   manifest. The historical Staticcheck cache-permission failure is explicitly
   preserved and followed by a successful command with a writable cache. No
   unexplained warning appears in the final test logs. These are inspected
   author-run results, not reviewer-executed suites.
8. No Go tests, race suite, exporter regeneration, runtime campaign, SQLite
   access, `.envrc` access, index mutation, branch mutation or subagent occurred.
   Only this report was written. One initial read-only `git diff --stat BASE
   HEAD` duplicated the supplied package's stat before the template's no-Git
   instruction was read; no subsequent Git command was run. This is a review
   procedure deviation, not a code or evidence change.

## Cannot verify / controller disposition

- The report's temporal TDD claim cannot be independently established from a
  final commit and retained RED/GREEN logs alone. Those logs exist and show the
  claimed semantic failures, but intermediate source snapshots are not present
  in this diff. No fabricated chronology is inferred or alleged.
- The manifest and successful logs support the author-run checks; this review
  does not independently authenticate their execution environment or rerun
  them. The recorded source hashes do match the reviewed inputs. A separate
  historical Git comparison of old sources to the base was not repeated.
- The independent Python semantic audit and hostile JSON rejection controls
  are Task 2 obligations. The structural probe here cannot substitute for
  those checks, and this approval does not preapprove their outcome.
- Manuscript synchronization, PDF build/inspection, the root build/module
  checks, global Staticcheck accounting, the reported frozen 70-package
  baseline and final whole-branch review remain controller/Task 3 obligations.
  They are neither verified nor claimed complete by this Task 1 review.
- Full concrete-product enumeration, mechanized proof, arbitrary contexts,
  fairness/divergence equivalence and runtime/distributed refinement are
  explicitly outside this task's claim. They were not verified and are not
  prerequisites silently discharged by this approval.

## Assessment

**Task quality: APPROVE.** The implementation realizes the stated fixed
progress contract with explicit custody and identity boundaries, preserves
counterexamples and retry behavior, and supplies an inspectable conditional
lifting argument. The minor diagnostic issue does not undermine the tested
properties or the source argument; independent export checking and final
manuscript/branch qualification remain necessary subsequent gates.
