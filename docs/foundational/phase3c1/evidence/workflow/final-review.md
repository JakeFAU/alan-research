# Phase 3C.1 whole-branch independent review

Reviewed base `1456a0000a0e14bb3dd105026189717c71ab905b` through head
`dce9de0a7b5795dd646dac58526e70729bc6231b`, four commits, using the supplied
`review-1456a000..dce9de0a.diff`. This is the broad final technical review of
the Contract Preservation Under Composition packet. It is not merge authorization.

## Question and scope

Does retaining named original progress, immutable bindings and physical-owner
policy incidence preserve the declared observations and existential completion
through compatible grouping in the fixed chain/cycle, under authentic transport,
recoverable crashes and atomic application/result retention?

I read the repository AGENTS, final-review instructions, contract, plan, progress
ledger/rulings, all three author reports and independent task reviews, and the
requested code-reviewer template. I reviewed the model, tests, exporter,
independent checker, evidence procedures, argument and manuscript changes.
Generated exports/manifests were inspected with read-only structural and digest
probes instead of dumping whole graphs. No prior reviewer REJECT appears in the
provided task reports or ledger; preserved implementation/auditor RED failures
remain distinct from review decisions.

## Strengths and load-bearing argument assessment

- **Projection covers the concrete classes.** `concrete.go:163` filters the
  unchanged edge relation by the actual actor's availability and applies the
  serial gate only to application. Commit, Apply and Consume are the only
  progress changes; transport, recording, ACK/release and physical interruption
  stutter. The one physical-owner record affects both incoming and outgoing
  roles. Each original's caller crash/recovery projects to the old edge relation;
  callee interruption is an edge stutter. The fixed full binding is retained
  outside the old-edge renaming. This supports the transition-class argument in
  `docs/foundational/phase3c1/01-argument.md:36` without treating reuse of the
  reducer as independent evidence.
- **Custody is an inductive premise actually preserved by the imported rules.**
  Application sets the permanent tombstone and callee receipt atomically.
  Send/receive/record use retained copies; ACK requires durable recording;
  release can erase callee custody only after that recording. Caller crashes
  erase only volatile custody. Duplicate requests/results/ACKs and repeated
  release do not permit a second application or destroy the durable custody
  needed for consumption. I checked the unchanged `research/twoagent/model.go`
  and separately implemented `oracle.go` for these named risks, rather than
  relying on the new adapter or test agreement alone.
- **The lifting clause is universal over reachable representatives.**
  `research/contractcomposition/concrete.go:214` and
  `docs/foundational/phase3c1/01-argument.md:98` handle the necessary cases:
  U needs caller recovery and commit; W uses existing admission or a queued/
  retained committed request, then application; P uses durable recording,
  volatile receipt, queued receipt, or actual retained callee custody, then
  consumption. Occupied channels do not require draining unrelated originals.
  Recoveries do not change outgoing stages; the serial gate therefore stays
  true during the W preparation. Receipt plumbing remains serviceable while
  waiting. No case needs another original to advance or an effect to reapply.
  Neither source validation nor lifting uses WorldResult to acquire bytes or
  choose work. Reachability is a stated premise, not falsely decided by the
  representation validator. I found no missing reachable custody case.
- **Both completion-path directions are present.** Finite concrete paths project
  after deleting stutters. Each abstract path lifts from the representative
  actually reached by its preceding lift. Complete is the conjunction over all
  original D registers, so the stated CanComplete equivalence follows under
  these premises. The result does not assert immediate readiness, physical
  schedule matching, channel drain or inevitable progress. A genuine concrete
  retry loop remains in the evidence.
- **Grouping retains exactly the facts the argument reads.**
  `research/contractcomposition/compose.go:34`, `:87`, `:126` validate owners,
  exact bindings, resource/outgoing uniqueness, endpoint descriptors and closed
  ports. Joining does not erase internal originals or combine availability or
  policy records. The fixed leaf/join/close grammar yields the same named closed
  configuration under both parenthesizations; the previously checked projection
  and lifting conditions depend on those retained coordinates. Canonical union
  equality is correctly only the algebraic part of the result. This does not
  justify arbitrary open-context replacement or a smaller constant-size summary.
- **The dependency claim has the right quantifiers.** Processing outgoing
  prerequisites before serial incoming dependents gives a completion suffix
  from every reachable state when the static graph is acyclic. This is a
  sufficient condition, not necessity for a particular state. The all-serial
  initial state can complete; W/W/W cannot. The source and manuscript keep that
  compatibility obstruction separate from an authority or custody violation.
- **R0 and R1 remain genuine, bounded attacks.** The seven event histories are
  legal under the inspected protocol. R0's continuation changes only AB after
  BC completes, with identical full exposed BC transport and all availability
  records. Both states can complete, while Complete and visible internal
  application/consumption differ. R1 keeps full BC/CA transport and availability
  equal, with CA admitted and A up, while AB U versus W changes immediate CA
  readiness. It refutes the stated readiness summary only; no weaker
  readiness-free trace relation is declared falsified by this pair. Neither
  witness depends on a failed premise or an unsubstantiated minimality claim.
- **The independent validator has meaningful coverage and a visible boundary.**
  `docs/foundational/phase3c1/tools/audit.py:164` compares complete expected
  state/arc relations, including non-tree incoming arcs, and independently
  classifies Complete/CanComplete. `:229` rebuilds both grouped configurations;
  `:291` maintains shared availability and actual causal custody from initial
  state; `:355` compares every selected post-state/projection; `:393` recomputes
  the collision, retry and loss-history predicates. The code does not call Go
  reducers. Its shared protocol specification and selected-history scope are
  disclosed, and prose truth/universal lifting are not presented as parser
  guarantees. The missing-loss-history RED and corrected unavailable-actor
  control are preserved.
- **Reconciliation matches the result.** The new manuscript section, Abstract,
  Conclusion and roadmap retain the fixed domain, full debt, static authority,
  atomicity, authentic custody, existential progress and evidence categories.
  Complete abstract checks, factored edge probes and selected concrete replay
  are not conflated with runtime refinement. The citation is explicitly
  background rather than an unverified instantiation of a general theorem.
  Earlier Phase 3B boundaries remain visible and Phase 3C remains open.

## Findings

### Critical

None found.

### Important

None found.

### Minor — M1 remains open, nonblocking

**File:** `research/contractcomposition/concrete_test.go:217` and `:239`.

The projection and lifting assertions inside the 24,258-context loop report
generic failures without topology, target original, service policy, availability,
old source state or offending transition. A future regression would fail the
test but require instrumentation to isolate its witness. Add these values to
the failure messages, or introduce named subtests at an appropriate outer
boundary. The assertions themselves are substantive and fail on the relevant
semantic violation; their diagnostic weakness does not undermine this packet's
present results or proof argument.

**Final-review triage:** confirmed as a real minor diagnostic improvement;
recommend explicit deferral to a narrow follow-up. It is **unfixed**, not rejected
or silently resolved, and is not a condition of this technical approval. The
controller should preserve that disposition in the final review archive.

## Actual checks performed

1. Read the supplied whole-branch diff in source/document chunks, including all
   changed Go/Python implementation and tests, contract/plan, argument/schema,
   results/qualification, preflight alternatives and corrections, archive/capture
   procedures, and changed Markdown/LaTeX/BibTeX content. Large combined outputs
   that truncated were followed by narrower reads of the affected material.
   Applied the requested review workflow and Go test/code-clarity guidance;
   did not dispatch subagents.
2. Read-only Git checks confirmed a clean checkout at the specified head and
   exactly the four declared commits. `git diff --stat BASE..HEAD` showed 190
   paths confined to the new model, Phase 3C.1 packet, roadmap and manuscript.
   `git diff --check BASE..HEAD` exited 0. A Python hunk probe checked all 130
   textual diff files against current head lines and base removed/context lines.
   No source, branch, index or HEAD mutation occurred.
3. Read the unchanged twoagent transition implementation and causal ledger for
   the named occupied-channel, duplicated-delivery, release, actor-availability,
   oracle-access and no-reapplication risks. Source reasoning above is the
   universal argument review; no concrete product enumeration was performed.
4. A read-only digest probe verified 17 Task 1 source hashes, 22 Task 2 source
   hashes, 84 Task 2 retained-file hashes, all 36 Task 2 raw/packed compressed
   artifacts, and the three preflight artifacts. It also verified all 598
   initially frozen old inputs and all 637 current Go/module inputs. The 29
   supplemental preexisting module files were independently compared against
   `git show BASE:path` bytes and matched both manifest and current files.
5. The same probe confirmed the raw Task 1 export and decompressed Task 2 export
   are identical: 1,584,501 bytes, SHA-256
   `48b4cfb4cc10b7dde4979874f669014f4863e1e05e48737ef4561f119dfc3602`.
   All changed generated JSON structures (35 files) and gzip payloads (55 files)
   parsed/decompressed successfully. This includes generated records not printed
   in the prose diff review.
6. An independent read-only structural/relation probe over the retained export
   checked unique states, every expected named U/W/P/D successor against the
   declared guard, exact bindings/targets, full reachability, BFS parent arcs,
   Complete and reverse reachability. It recomputed 639 states and 1,237 arcs;
   only cycle/mask0/W/W/W lacks a completion suffix. This checked the archived
   abstract graph, not a new exporter run or a concrete model campaign.
7. That probe checked the R0/R1 continuation prefixes and identical full exposed
   transport/availability records, inspected all seven actual event sequences,
   and confirmed exact retry states 3 and 5. It checked all 35 hostile raw input
   digests against their retained rejection logs and expected reasons. The
   non-tree incoming-arc deletion retains every BFS parent edge, so its negative
   check is substantively stronger than checking the predecessor tree alone.
   The audit CLI and unittest suite were not rerun.
8. A read-only qualification probe parsed every baseline JSONL event: 70 package
   terminal events, 57 pass, 13 no-test skips, no failure; roster equals the fixed
   argv. The retained clean listing is exactly those 70 plus the two added
   packages. Inspected final ordinary/race logs show both added packages pass,
   with recorded exit zero. Inspected Task 1 RED/GREEN/constructor/oracle logs,
   final 311/24,258/184,256/14,058 factored counts, and Task 2 final audit/test
   outcomes. These are source-bound author/controller executions, not my reruns.
9. Inspected every root job/result/output record and verified its output hash and
   length. Clean writable-cache build, module verify, tidy-diff and scoped lint
   record exit zero. Global Staticcheck records exit one with exactly **13
   ordinary and 31 integration diagnostics**, byte-identical to Phase 3B.7.
   Initial cache-warning build/list outcomes remain archived separately. The
   task is not globally lint-clean and has no new whole-repository race or
   integration qualification.
10. After source inspection, ran the read-only
    `python3 docs/foundational/phase3c1/evidence/manuscript/verify.py`: PASS for
    38 qualification archive files, eight compressed manuscript artifacts, 31
    local links, current paper hashes and normalized new-section/changed
    Abstract/Conclusion alignment. Ran `pdftotext -layout
    paper/alan-dissertation.pdf -` inside a read-only probe: exact match to the
    retained text. Confirmed the 16 changed text pages are covered by the
    recorded 24-page visual-inspection roster. PDF SHA-256 is
    `15526232e2794655608173177662a26f774fcff30d6972ab45dd6ba3c6ae7029`.
11. Parsed baseline/current retained final-engine logs and confirmed identical
    warnings: **four underfull hboxes, one overfull hbox (5.2184pt), and two
    math-mode font warnings**. Current PDF is recorded as 46 pages versus 45
    baseline; the changed text evidence agrees. Inherited Markdown-only Section
    7 and section/theorem numbering differences remain acknowledged. I did not
    rebuild or independently visually inspect the PDF in this final review.

No author Go/Python suite, exporter, runtime capture, build, lint or module
command was rerun. No retained SQLite or runtime `.envrc` was opened. The only
checkout write by this reviewer is this report.

## Declined to judge / explicit limits for controller disposition

- Mechanized proof and exhaustive reachable concrete-product correctness: the
  packet supplies a conditional source argument, factored probes and selected
  histories. I assessed that argument, but neither the counts nor this review
  constitute mechanization or full-product enumeration.
- Arbitrary topology/program contexts, shared-authority fan-out, constant-size
  interfaces and reduced verification cost: the leaf/join/close grammar and
  disjoint fixed-resource premises deliberately exclude these extensions.
- Dynamic identities/generations/grants, cancellation, revocation, shared-resource
  depletion, non-atomic application/retention, permanent owner loss and deleted
  durable custody: these violate or extend the fixed premises and need their own
  models and adverse witnesses, not a reinterpretation of this approval.
- Runtime/distributed refinement, independent-owner recovery, reverse runtime
  refinement, cross-invocation delegation, delegated cancellation-after-admission
  custody and whole-pair split-effect completion: no new runtime campaign or
  correspondence theorem is delivered here; the inherited Phase 3B limitations
  remain and this packet does not discharge them.
- Universal liveness, fairness transfer, divergence equivalence, matching physical
  schedules, global quiescence, arbitrary result-value equivalence and exact
  resource-cost preservation: the observation and theorem are value-abstracted
  finite progress/possible completion, with an explicit retry divergence witness.
- Generic malformed-input hardening of internal model APIs: the README declares
  valid config/state preconditions for raw successors/projections; checked
  constructors, Lift and witness/audit boundaries enforce the scoped interface.
  This is not a production general-purpose composition API.
- Historical wall-clock execution, environment authenticity and exact TDD edit
  chronology: current hashes and retained command outcomes are consistent, but
  intermediate source snapshots and the initial Python RED inner exit metadata
  are not all retained. No independent rerun or historical authentication is
  claimed; missing metadata is disclosed rather than fabricated.
- The independent auditor author's historical reading behavior: source shows no
  producer imports/reducer calls and separate algorithms, but noninspection is
  process testimony. The shared contract/schema/old protocol specification is
  an explicit limit on independence, not independently proven semantics.
- A fresh third primary-literature verification or DOI check: I reviewed the
  bibliography change, source-check notes and Task 3 author's/reviewer's actual
  primary-source checks. I did not reopen the external sources. The published
  versus author-revision pagination distinction remains recorded, and no new
  general-theorem application or verified DOI is claimed.
- A second full PDF rebuild or visual review: source/hash/text/warning agreement
  was checked here; visual qualification is attributed to the author's 24-page
  inspection and Task 3 reviewer's direct pages 42–44 inspection. The inherited
  page-32 overhang, font/box warnings and global Markdown/LaTeX differences were
  neither repaired nor silently certified away.
- Repair of the inherited 13/31 global Staticcheck diagnostics and wholesale
  unrelated formatting/test-style changes: no old executable changed, inherited
  diagnostics match exactly, and such cleanup is outside this research slice.
  M1 is separately retained above as the specific useful diagnostic follow-up.
- Unchanged manuscript theorems, historical model/runtime results and discarded
  preflight one-copy product bounds: the final claims preserve their historical
  status; this is not a new audit or execution of every earlier research phase
  or an acceptance of discarded alternatives as measured results.
- Controller-owned final workflow/review archival, Task 3/final-gate status
  updates, push and verified draft-PR delivery: these remain after the reviewed
  head. This review approves no unseen follow-up source change and supplies no
  merge or release authorization.

## Recommendations and assessment

Preserve M1 as an explicitly deferred nonblocking item; archive the exact task
and final reviews and update the pending status in `06-reviews.md` during the
controller's delivery pass. Keep the conditional proof, finite checks, selected
history evidence and inherited limitations distinct in the draft PR.

**Verdict: APPROVE.** No blocking defect was found in the reviewed branch's
fixed-fragment argument, implementation, independent audit or reconciliation.
The only actionable finding is minor M1, with the explicit deferral above.
Ready for the controller's draft-PR delivery steps; this verdict is not merge
authorization and does not establish general Phase 3C closure.
