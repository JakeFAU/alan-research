# Phase 3B.4 explicit request binding implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development to implement this plan task by task. The controller dispatches workers and reviewers; workers do not dispatch subagents.

**Goal:** Resolve whether exact retained request selection can replace v1 instruction adjacency across declared finite work and same-Store recovery without expanding authority.

**Architecture:** First reproduce the unchanged v1 obstruction and freeze a smallest-evidence contract. Implement a separately versioned historical-selection predicate in the runtime, with an independently specified relation and raw checker in research. Preserve the prior model exports, old contracts, negative histories and unresolved correspondence boundaries.

**Tech Stack:** Existing Go 1.27.3 modules, SQLite through existing dependency, public retained programs, controlled loopback endpoints, acknowledged subprocess cuts, Go independent models/checkers and the existing manuscript build.

**Spec:** `docs/foundational/phase3b4/00-brief.md` (verbatim user packet, SHA256 a51a078bcaa6a8dbc8fe9b6b0bbbb9518c62a5db7965a5f9da95edb5db31af72), repository AGENTS.md, and Task 1's independently reviewed contract.

## Global Constraints

- Execute through independently reviewed draft delivery, not merely a plan or first passing example. Do not merge or release. Ordinary decisions within scope are delegated; preserve and explain genuine obstructions rather than manufacturing success.
- No new dependencies, paid services, distributed ownership, multihop delegation or additional agents.
- Preserve formats 13/14 and alan/request-authority/v1 under their original interpretation.
- Preserve all old deterministic exports byte-for-byte.
- No direct SQL, internal-import escape hatch or host-injected consumption may establish a positive history.
- Replay remains effect-free. Open validates before renewing authority and retains the existing failed-open guarantees.
- Generic unknown mutations remain unresolved without automatic retransmission.
- No new bibliography entries unless the argument actually uses them.
- Do not merge, release, force-push shared history or change unrelated branches.
- Research worktree: `/Users/jacob/.codex/worktrees/phase3b4-research/alan-research`, base `9ddfe29a2a5fa43a06656514af236401cc07914b`.
- Runtime worktree: `/private/tmp/alan-phase3b3-runtime`, base `1fdd4c9e1720bf60b344914bcbc3d0938ba8d678`. Both use `codex/phase3b4-explicit-request-binding`. The old directory name is intentional reuse of a clean isolated checkout. Leave original checkouts, `.envrc`, and the active performance worktree untouched.
- Use explicit working directories. Never open archived SQLite files in place; copy DB/WAL/SHM to disposable storage before SQLite inspection. Archive failures as failures, not policy rejections.
- Each substantive task ends in a narrow commit, self-review report and independent spec/quality review. Record exact commands, source identity and evidence paths. Do not claim tests not actually run. The controller performs final whole-branch review and PR delivery.

## Review Focus

- A second authentic same-peer request with identical payload and obligation must not replace the grant's exact original (Tasks 2, 4, 5).
- A legitimate selection occurring after service commitment cannot retroactively authorize that service (Tasks 1, 3, 5).
- Current revocation and owner epoch checks must remain effective after historical binding succeeds (Tasks 3–5).
- Profile/schema/nested-descriptor disagreement must reject without mutating the Store; old readers must reject new formats (Tasks 3, 5).
- Incomplete causal prefixes must receive a bounded diagnostic instead of panic or invented missing evidence (Task 5).

---

### Task 1: Reproduce v1 and specify the candidate

**Files:** Create research `01-charter.md`, `02-contract.md`, `evidence/baseline/`, `evidence/construction/` under `docs/foundational/phase3b4/`. Create runtime `internal/durabletest/phase3b4_construction_test.go` only if needed for the public-construction witness. Retain preflight reports. Do not modify production authorization, principal oracle or old exports in this task.

**Interfaces:** Consume public `delegation.Reciprocal`, `runtime.Create/Open/Run/Inspect/Replay`, immutable grant/use records and existing Phase 3B.3 exporter. Produce an exact candidate contract (wire version/format/public example signature, missing-versus-invalid handling, causal prefix, index/transaction requirements) and genuine supported/refused/padded construction exports. New representation is not fixed until evidence and review; first evaluate derivation from grant.OriginalRequestID. Candidate v2 numbers 15/16 appear free but verify allocation before reserving them. Keep `RequestAuthorityVersion` as v1; proposed additional public constant `RequestAuthorityBindingVersion = "alan/request-authority/v2"` requires review.

- [ ] Read the user packet Task 1, previous contract/results/selection-narrowing review, current authority/version/public fixture paths and both retained preflight reports.
- [ ] Record current baseline provenance and actual PR head/merge tree differences. Reproduce `TestPhase3B3ReciprocalOrdinary/s2-handshake` and `TestPhase3B3ReciprocalBindings/wrong-prior-selection` with `ALAN_PHASE3B3_EVIDENCE` pointing to a fresh Phase 3B.4 directory. Preserve original IDs, decision reason, grant, actual selection, unpaid obligations and endpoint contacts. Assert supported completion and genuine refusal specifically due to adjacency, with unchanged equal service sequence.
- [ ] Establish public pre-provision construction for inserted lengths 0, 1 and 3 and both finite forward branch outcomes. Genuine claims must name the actual service sequence; update every dependent locator before Create and hash the actual artifacts. No retained-row repair. A v1 refusal at the intended selection layer establishes construction, not success of the candidate. Keep baseline source and fixture source identities distinct.
- [ ] Freeze the candidate contract as review-pending: exact q/g/o/c/e joins, same logical generation, strict 0<c<e.Sequence, grant before decision, real reduction and admission, source permission and authorization, roots/revisions/tuple/selected path/current gate, immutable prefix order and refusal stickiness. Prefer existing witness fields. Specify strict semantic-version equality through profile/schema/claims/uses/witnesses/gates. Distinguish natural absence/refusal from corrupted retained evidence. Compare all three representations and explain rejected alternatives.
- [ ] Decide whether a v2-only index on existing admissions(observation_id) is needed; no new authority table or last-request cache. State inherited prefix scans separately from the newly added lookup. Freeze cost guard lengths 0,1,8,32 and measured operations before execution.
- [ ] Specify the exact program/cut/attack matrix from packet Tasks 2–5, predicted outcomes and current public signature for the new example. Preserve ambiguity as an explicit obstruction if public construction genuinely fails. Build an old-reader test binary from the unchanged selected baseline for later compatibility checks and record its hash.
- [ ] Run focused baseline/construction tests, formatting and diff checks; archive command/exits/source hashes. Commit only Task 1 fixtures/docs/evidence plus initial brief/plan/preflight docs. Report contract status review-pending; controller obtains independent contract/construction review before production or principal oracle.

### Task 2: Freeze the relation and bounded claims

**Files:** Create research `docs/foundational/phase3b4/04-model.md`, `05-matrix.md`, and any new auxiliary model under `research/requestbinding/` with focused tests. No edits to previous model exports, production runtime or principal raw oracle.

**Interfaces:** Consume reviewed `02-contract.md` and genuine Task 1 exports. Produce a source-level H/A/C argument and an explicitly bounded optional selection/intervening-work transition model with deterministic independent export. Existing Phase 3B.2 model and audit remain unchanged; new relation names exact principal original request IDs and keeps auxiliary requests separate.

- [ ] Define historical selection S, selection-to-service binding H, current validity V/resource headroom R, authority preservation and restricted correspondence separately. Give quantified domains and proof sketches with each premise used. Show future selection fails and historical selection survives revocation while gate eligibility does not.
- [ ] Define allowed set/copy scratch writes, finite forward branch paths preserving service operands, readiness receives with delivery/selection premises and current eligibility, logical generation identity and program identity within an execution. Instruction insertions define fresh provisions and an explicit identity correspondence, not byte-equal histories. Steps/costs/divergence and pending results/ACKs remain observable in a sidecar.
- [ ] Freeze finite semantic family lengths 0,1,3, both branch outcomes, request/ready orders, selected-before/after-grant schedules and owner-loss counts 0,1,2 as a covering family, not an uncontrolled Cartesian product. Cost guard 0,1,8,32 is separate. Specify each family and expected verdict before new runtime confirmation.
- [ ] Build an auxiliary model only where it adds an independent falsifiable check; if used, TDD pins q-selected-before-e, later same-peer selection preserving q, sticky early refusal, revocation at gate and source-level induction boundaries. Relate explicitly to old LateGrant admission vocabulary under nonrevoking known-result single-exposure premises. Do not project whole-Store death as per-agent crashes.
- [ ] Freeze raw/custody principal IDs and auxiliary-message handling requirements. Existing custody suffix/OR matching is insufficient for two same-peer originals: require exact original-request association, not payload/obligation identity. A genuinely distinct grant for q-prime is not inherently forbidden; negative substitutions must violate a frozen authentic binding.
- [ ] Predeclare checker-code mutants: omit c<e, latest-receive instead of exact bound q, omit a generation comparison, skip current revocation at gate. Predict concrete false acceptance or false refusal killed at assertion, not build error or panic. Record all remaining old model mismatch categories.
- [ ] Verify deterministic output and focused tests; commit model/argument/matrix. Controller obtains independent relation/model review before principal oracle construction or treating runtime runs as confirmation.

### Task 3: Implement the versioned historical-selection predicate

**Files:** Runtime `internal/executor/request_authority.go`, `request_authority_admission.go`, `store.go`, `composition_store.go`, `composition_types.go`, `receipt_conflict.go`, `inspection_execution.go`, `inspection_deployment.go`, `runtime/runtime.go` and focused sibling tests only as required. New public example in `examples/delegation/request_binding.go` with tests and documentation. Use reviewed contract as exact interface authority.

**Interfaces:** Consume reviewed contract/relation. Preserve v1 constant/formats and existing histories. Produce opt-in version/profile and exact historical lookup reused by original decision, admission reconstruction, Replay/inspection/Open; no scheduler or separate budget ledger.

- [ ] Add a focused RED test for new profile plus request→ready→service authorization and for old-v1 refusal preservation. Record causal-layer failure, then add only required version routing/validation and exact selection derivation.
- [ ] Test each changed predicate with genuine provisioned histories and malformed disposable copies: admission-only, selected future observation, wrong generation/q/o/c/e/grant/root/revision/namespace/tuple, no own fallback, sticky refusal, revoked fresh gate and retained-result consumption. Current gates reconstruct original operations and cannot trust forged attempt/reservation pairs.
- [ ] Keep v1 evaluate path unchanged. New formats/profile/nested versions agree exactly; old reader rejects new semantics; new reader preserves v1. Failed Open validates before owner renewal and leaves immutable inputs unchanged. Test conflict promotion under new format and mixed/downgraded tags including coherent descriptor attacks.
- [ ] Add only the reviewed indexed lookup, prove query starts from exact q, and distinguish inherited global-prefix costs. Do not scan for a convenient prior request. No mutable request cache. Keep dispatch lock/worker-release lifetime unchanged.
- [ ] Implement public example contract and finite construction with actual service locators and full manifests. Run focused ordinary/integration/race and Staticcheck, formatting and module checks applicable to changed code. Retain RED/GREEN evidence. Commit independently testable increments and self-review entire task diff. Controller obtains runtime/recovery and compatibility review before Task 4 final evidence.

### Task 4: Exercise ordinary programs and acknowledged process cuts

**Files:** Runtime `internal/durabletest/phase3b4_*_test.go`, `examples/delegation/request_binding.go` only for qualified fixture corrections, runtime qualification docs and fresh evidence exports. Research evidence manifest/index may be prepared but principal oracle belongs to Task 5.

**Interfaces:** Consume reviewed program manifest, versioned runtime and frozen matrix. Produce raw histories, canonical inspection, source manifests, boundaries and independent endpoint databases for automatic and controlled runs, named distinctly. Do not change authority semantics to make a schedule green.

- [ ] Implement/run all positive and negative matrix families: v1/v2 adjacent; request-ready-service; 0/1/3 local set/copy; forward branches; same-peer equal-payload/obligation two requests; selection before/after grant; admission-only; future/forged selection; early sticky refusal; revocation/expiry; tuple/identity substitutions; Cartesian/ambient/result/fallback. Both reciprocal result consumptions, result custody, ACKs, terminals and global quiescence have separate assertions.
- [ ] Implement acknowledged SIGKILL cuts at admission-before-selection, selection-before-ready, grant+selection+padding-before-service, service-before-decision, allow-before-fresh-gate, result-before-consumption, two consecutive owner losses, and revoked-before-recovery-gate. Fresh subprocesses reopen original Store. Preserve all original IDs across cuts and record epochs/charges/gates/pending work/contacts. Unselected does not become selected through Open; completed known-result recovery does not contact again.
- [ ] Keep generic unknown and merely reserved/applied-unrecorded cases distinct. Demonstrate ordinary automatic v2 reciprocal completion, then label controlled adverse schedules separately. Retain model-unsupported recovery rather than inventing per-agent crash steps.
- [ ] Export identity bijections and sidecar costs for fresh padded provisions. Run focused integration/race, prior delegation/custody/recovery regressions and predeclared cost guard. Archive failures and fixes; no fabricated target leaf count. Commit fixtures and qualification machinery. Controller reviews actual causal records and test assertions.

### Task 5: Independently audit and qualify exact executable sources

**Files:** Research `research/delegationcorrespondence/` and command entrypoints/tests as narrowly needed, new `docs/foundational/phase3b4/evidence/`, `06-qualification.md`, `07-runbook.md`. Runtime compact qualification evidence/docs. Separate known mutation diagnostic fix from policy extension in commits/evidence.

**Interfaces:** Consume frozen contract, model relation, manifest/exports and runtime source. Preserve old-v1 parser/policy. Produce three independent verdicts: raw policy, physical endpoint, supported model projection, with explicit unsupported categories. No runtime authorization, reducer or Replay imports in oracle.

- [ ] Add RED/GREEN oracle assertions for strict original-message/route/admission/observation/actual-selection/grant/exact-service/first-decision/current-gate/result custody chain at historical prefixes. Preserve immutable joins across captures and real owner-renewal evidence. Pin principal originals explicitly and account auxiliary requests separately; payload/suffix identity never substitutes for q.
- [ ] Repair incomplete-prefix mutation generator diagnostic narrowly and test no panic. Missing evidence is a diagnostic, never invented. Dense, coherently rehashed wrong generation/selection/future binding and gate revocation attacks must reach meaningful predicate rejection.
- [ ] Run each checker-code mutant from an isolated source copy; record exact changed predicate, expected false acceptance/refusal and assertion failure killing it. Build failure/panic is not mutant detection. Verify raw validity cannot be excused by unsupported projection.
- [ ] Freeze and commit exact research/runtime executable pair, source manifests and old export hashes before final qualification. Record tools/OS/architecture/env/commands/exits/watchdogs, all optional skips and inherited findings separately. Do not mix evidence from earlier sources as final proof.
- [ ] Run full ordinary tests in both repositories, full runtime integration and race required by shared validators, affected research integration/race, prior authority/custody/causal-recovery/version regressions; actual build/module/format/Staticcheck/vet commands where applicable. Hash before/after and rerun affected qualification after executable fixes. Use process sessions and logs rather than silent blocking waits.
- [ ] Archive compact reproducible raw evidence and evidence digests plus exact source pair. Record cost guard observations without extrapolating a performance theorem. Independent raw-evidence/compatibility reviewer inspects actual records and mutant kills; preserve every rejection/correction. Controller also reconciles all earlier reviews before manuscript conclusions.

### Task 6: Reconcile the manuscript and deliver reviewed drafts

**Files:** Research `docs/foundational/phase3b4/08-results.md`, `09-reviews.md`, evidence/review index, `paper/alan-dissertation.md`, `.tex`, rebuilt PDF and `references.bib` only if needed. Runtime public docs/qualification summary. No unrelated source changes.

**Interfaces:** Consume exact executable pair, all reviews and qualification results. Produce bounded final claims, synchronized manuscript/PDF and cross-linked draft PRs. Any correctness source fix returns to affected qualification; documentation-only delivery commits are separate from executable qualification.

- [ ] Answer the packet's six disposition questions explicitly, with proved/proof-sketch/bounded-model/tested/open categories, and recommend only one next question. Preserve old discovery history and remaining whole-Store crash, separate-grant revocation and applied-unrecorded gaps.
- [ ] Synchronize Markdown/LaTeX Abstract, affected body and Conclusion. Check bibliography implications without inventing references. Controller obtains substantive manuscript/whole-change review before final PDF build, then applies PDF skill build and visual inspection to changed pages and continuations. Rebuild after correctness edits.
- [ ] Inspect status/full diff and required checks before final commits. Archive this plan's ledger, briefs, reports and review artifacts without secrets; remove only its scratch workspace when the finishing skill directs.
- [ ] Controller pushes task branches and opens cross-linked draft PRs. State exact bases, executable and delivery commits, versions, scope/invariant/implementation, RED evidence, verification, remaining risks, reviews/rejections/skips and absence of unrelated changes. No automatic merge. Attach both PR artifacts to this chat. If main advanced, qualify an explicit prospective combined tree only when necessary and report its exact identity.
