# Phase 3B.5 Grant Revocation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Qualify separate accepted delegation followed by exact-root revocation in a bounded two-agent model and a scoped projection of sealed runtime evidence.

**Architecture:** Add one opt-in configuration to the existing finite authority model while preserving its old exports. Extend the separately implemented ledger and reverse-reachability query, then add a distinct concrete projection with mandatory physical evidence. Reconcile bounded findings and retain obstructions.

**Tech Stack:** Existing Go 1.27.3 module, independent Go models/checkers, sealed JSON/JSONL archives, existing LaTeX/PDF toolchain. No new dependencies.

**Spec:** `docs/foundational/phase3b5/01-contract.md` and repository AGENTS.md.

## Global Constraints

- Deliver an independently reviewed draft research PR; do not merge or release.
- Preserve the old sixteen configurations, State/Event JSON shapes, enum values and all four historical deterministic exports byte-for-byte. New behavior is opt-in.
- Revocation targets one complete authority root, not a GrantID or composition budget root. BA uses b; the new configuration's AB service uses distinct B-owned c.
- First decisions and issuance outcomes are historical. Fresh admission checks present authority; old opportunities survive revocation but not server crash. Authentic retained results remain consumable.
- No new dependencies, runtime production changes, additional agents/owners, multihop, regrant, non-atomic model effects or unrelated refactors.
- Independent ledger/checker must not import runtime authorization, reducer or Replay or call operational authority helpers. No invented facts, endpoint contacts or principal consumptions.
- Raw validity, external evidence, continuity and model support are separate verdicts. Legacy projection remains unchanged. Unknown outcomes and whole-Store loss stay outside the new relation.
- Work only in `/Users/jacob/.codex/worktrees/phase3b4-research/alan-research` on `codex/phase3b5-grant-revocation`, base `4fea6c31eb8f438c0dea414aa5f8d7d0ecb5e8a4`. Do not touch other checkouts or `.envrc`, and never open archived SQLite in place.
- Use explicit workdirs and offline Go environment: `GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b5/go-build`.
- One implementation agent at a time; each task commits, self-reviews and receives independent spec/quality review. Report actual commands/exits and preserve failed attempts. No worker spawns subagents.

## Review Focus

- Root aliasing must not turn N07's distinct AB authority into revoked BA authority (Tasks 1–3).
- Rejected late issuance cannot erase the request or rewrite an old refusal (Tasks 1–2).
- A revoked old permit must survive until application or server crash; a stale permit cannot return after recovery (Tasks 1–2).
- Missing endpoint evidence and allowed-gate/no-result histories must never become synthetic Apply or zero-contact claims (Task 3).
- A final revoked set cannot retroactively invalidate a lawful first decision; a grant or outcome cannot move to a convenient causal position (Tasks 2–3).

---

### Task 1: Add the opt-in exact-root revocation model

**Files:** Modify `research/authorityservice/model.go`, `explore.go`, `cmd/export/main.go` only as needed; create `revocation_test.go` and command tests; write `docs/foundational/phase3b5/02-model.md` and `evidence/task1/`.

**Interfaces:** Produce `GrantRevocationConfig() Config` returning `{Policy:"S2", Scenario:"SeparateGrantRevocation"}`. `Initial`, `Binding`, `Successors`, `Explore` accept that opt-in configuration. Append `ABOwnPermission` after existing tuples with permission `(increment,counter-AB,blue)`, root `B/g1:c` only available in this configuration. `Configurations()` and default `export -config all` stay unchanged. `export -config S2/SeparateGrantRevocation` selects the new graph. Independent oracle support is Task 2; do not claim it here.

- [x] Read the contract and `/tmp/alan-phase3b5-model-recon.md`; read relevant source/tests. Add named failing tests for no inline grant, accepted/rejected separate issuance, preserved first allow/fresh refusal, pre-revocation permit applying after revoke, A crash/recovery clearing permit, late-grant sticky refusal, distinct AB root with genuine reciprocal completion, and immutable consumed results.
- [x] Run focused new tests and record genuine RED evidence. Implement only configuration routing, distinct tuple/root, separate issuance status and revoke enablement. Keep old numeric values, successor ordering and state/event fields stable.
- [x] Add new-scenario-only witness cuts that explicitly identify BA's application after revocation and revoked crash opportunity; avoid the generic cut accidentally selecting AB/c application. Test complete edge descriptor reconstruction and expected target original. Preserve old cuts/output.
- [x] Run focused GREEN, old authorityservice model tests and new complete `Explore`. Preserve shortest failing universal-completion witness, safe counterexample states and retry cycles. Record counts only after actual execution. Write transition/invariant proof sketches and limits in `02-model.md`.
- [x] Compare old default model export with Phase 3B.2 archived model bytes, and run new export twice for determinism. Run formatting/diff checks. Commit task files/evidence only; report exact command results and source identity. No oracle changes, manuscript or production runtime edits.

### Task 2: Independently audit and characterize progress

**Files:** Modify `research/authorityservice/oracle.go`, add `revocation_oracle_test.go` and narrowly scoped corpus helpers; create `research/authorityservice/cmd/revocation/main.go` with tests; extend `research/delegationcorrespondence/progress.go` and progress tests/command selection; write `docs/foundational/phase3b5/03-authority-progress.md` and `evidence/task2/`.

**Interfaces:** Consume Task 1 `GrantRevocationConfig`. Existing `AuditTrace`, `AuditGraph` accept it using independently declared roots/bindings/transitions. Produce `GrantRevocationOracleCorpus() (OracleCorpusReport,error)` and `cmd/revocation` deterministic JSON containing the new graph audit and corpus; existing `cmd/audit` remains unchanged. `QueryProgress(config)` additionally accepts only the new opt-in config; its old outputs remain unchanged.

- [x] Read the contract and operational interface, then independently implement new ledger configuration/binding/root support. Add RED/GREEN positive and hostile trace tests for immutable first decisions, exact b/c roots, grant identity/ancestry/complete tuple/original/generation, rejected issuance after revoke, current admission, stale epoch, and surviving permit. Do not call model authority helpers from the ledger.
- [x] Extend reverse reachability with the contract's candidate predicate; test every reachable state against computed safe-suffix reachability and audit witnesses independently. Preserve any mismatch before revising a candidate. Explicitly test AB refusal unreachability and sticky BA refusals, one revoked permit crash counterexample and successful post-revocation custody suffix.
- [x] Run complete new `AuditGraph`, deterministic corpus and progress twice; preserve states/arcs/hashes, shortest witnesses, SCC boundaries and exact existential interpretation. Give necessity/sufficiency proof sketches restricted to reachable finite states in `03-authority-progress.md`.
- [x] Run compiled independent-ledger mutants from disposable source copies: omit fresh current-root check, revoke c with b, accept post-revocation issuance, and retain crash permit. Each must build and be killed by its targeted assertion; preserve failures separately. Retain positive old-permit controls so an overrestrictive checker is detectable.
- [x] Reproduce all four old deterministic exports against archived bytes; old current-vs-current determinism alone is insufficient. Run focused tests, format and scoped Staticcheck. Commit task files/evidence; self-review and report commands, exits, exact semantics and limitations.

### Task 3: Map sealed exact-root revocation histories

**Files:** Create `research/delegationcorrespondence/revocation.go`, focused sibling tests, `cmd/revocation/main.go` and tests; reuse or minimally factor `project.go` only with byte-preserving legacy regressions. Add compact JSON fixtures/provenance only if needed. Write `docs/foundational/phase3b5/04-correspondence.md`, `05-runbook.md`, and `evidence/task3/`.

**Interfaces:** Produce `CheckRevocationDirectory(dir string) (RevocationCaseReport,error)` and `cmd/revocation -case DIR`. The new report contains the legacy `CaseReport`, new per-capture projections, explicit sidecar root/identity/epoch/debit/gate joins, final support/reasons and independent model audits. Consume existing `CheckBindingRaw`, `CheckExternal`, `CheckContinuity` and Task 1/2 model/oracle APIs. Legacy APIs/outputs remain unchanged.

- [x] Read the contract and `/tmp/alan-phase3b5-relation-recon.md`. Verify sealed archive and selected member hashes, extract declared artifacts only to disposable storage (hash binary database members as bytes; never open them with SQLite). Add RED tests using genuine N05/N06/N07/P02-v2 normal/race cases and outside-domain applied-unknown-revoked-before-fresh-gate. Preserve historic producer/source identity; no new runtime execution claim.
- [x] Require both endpoint journals, endpoint snapshot, declared capture digests, exact principals, independently valid raw binding and continuity. Define a stable exact BA→b and distinct B AB→c map with complete roots. Accepted grant/revoke originate only at their own lifecycle facts. Map raw first decisions, selected path/root/grant and fresh disposition; compare all explicit descriptors. Sidecar preserves runtime revision separately from model revocation era and concrete selection/cost identity.
- [x] Project each authentic capture in fact order through uniquely enabled edges and independent `AuditTrace`. Intermediate permitted-but-unrecorded cuts retain unresolved opportunity and no Apply. Final known-result Apply requires matching single contact/application and retained receipt. N06 final proves retained allow/debit plus fresh refusal and zero contacts under physical evidence; N07 final proves old opportunity/custody plus unaffected AB authority. Preserve old unsupported projections beside new scoped support.
- [x] Add coherent mapping/data attacks: swapped or aliased complete roots, wrong revoke target, moved grant/revoke fact, first allow evaluated at final revocation, fresh allowed gate after revoke, result synthesized from refusal, future/wrong selected request, stale epoch, missing one journal/snapshot and applied-unrecorded outside-domain. Rehash disposable evidence where needed so tests reach the intended semantic predicate. Record exact rejection or declared unsupported reason, without excusing raw invalidity as unsupported.
- [x] Run new projection twice for determinism on both sealed modes, focused positive/negative tests, legacy projection regression on P02 and unchanged old N06/N07 refusal. Run at least one compiled projector scope/eligibility mutant with meaningful false-accept/false-reject assertion, and retain existing concrete gate mutant evidence. Commit files/evidence and self-review. Report every support/exclusion verdict and all residual assumptions.

### Task 4: Qualify, reconcile and prepare delivery

**Files:** Create `docs/foundational/phase3b5/06-qualification.md`, `07-results.md`, `08-reviews.md` and evidence manifests; modify `paper/alan-dissertation.md`, `.tex`, rebuild `.pdf`, and update bibliography only if actually used. Controller owns final qualification/review/delivery coordination.

**Interfaces:** Consume reviewed Tasks 1–3 and exact executable source. Produce bounded result/counterexample statements, reproducibility manifests, synchronized manuscript and reviewed PDF. No production runtime companion PR unless a separately justified concrete defect changes scope.

- [x] Freeze/commit executable source and capture inputs/toolchain manifests. Run full ordinary research tests, affected race tests, complete model/audit/progress and both-mode scoped corpus, relevant old corpus regressions, build/module/format/scoped Staticcheck. Record inherited global lint separately. Inventory all named skips and apparatus failures. Requalify affected checks after executable fixes.
- [x] State precise result categories: finite graph safety/progress characterization, source-level proof sketch, conditional forward trace relation, actual sealed evidence, shortest universal-completion counterexamples, and unresolved wider claims. Reconcile Abstract/body/Conclusion in Markdown and LaTeX; preserve older phase findings as historical. Propose one next bounded question only.
- [x] Obtain task spec/quality review of evidence/manuscript; then controller runs one broad whole-branch review and scoped fixes as necessary. Apply current PDF skill, build only after substantive review, inspect changed pages/continuations visually, and record PDF/source hashes. Never label an unexecuted check passed.
- [ ] Archive all task briefs/reports/reviews, ledger and rulings. Controller verifies clean intended diff, pushes the task branch and opens/attaches a draft research PR against main. No merge, release, force push or unrelated cleanup. Final report includes question, bounded result, tests, changed artifacts, remaining obligations and material rulings with their costs.
