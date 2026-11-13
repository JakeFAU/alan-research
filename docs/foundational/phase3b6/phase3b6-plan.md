# Phase 3B.6 Effect Observation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Distinguish historical admission, world application and recoverable local evidence across one shared owner loss and exact-root revocation.

**Architecture:** A new finite single-service boundary model separates transport/application/retention, with an independently interpreted event ledger and observation partitions. A separate opt-in checker relates two sealed historical protocols to that model through explicit causal certificates. Existing two-agent atomic semantics and reports remain untouched.

**Tech Stack:** Existing Go module, Go 1.27.4, Staticcheck 2026.2.2, sealed JSON/JSONL, existing LaTeX/PDF tools. No new dependencies.

**Spec:** `docs/foundational/phase3b6/01-contract.md` and repository AGENTS.md.

## Global Constraints

- Deliver an independently reviewed draft research PR; do not merge or release.
- Preserve old sixteen authority-service configurations/order, State/Event JSON/enums, four historical deterministic exports, default command scope, old ledger/progress outputs, Phase 3B.5 model/audit/progress exports, and all ten old revocation concrete report bytes. New behavior is opt-in.
- One selected BA service boundary in a two-principal context, one shared Store-owner loss/reopen, one allowed physical gate/contact/application, one exact b revocation and optional refused fresh gate. Distinct c survives. No whole-pair recovery or general progress theorem.
- Keep historical allow/debit/gate, volatile opportunity/response, external packet/application and local retention/consumption separate. Endpoint evidence is not an import authority or recoverable result.
- No production runtime changes, new runtime campaign, dependencies, automatic retransmission or unrelated refactors. Never open archived SQLite; hash it as bytes.
- Independent ledger must not call operational successor, authority, transition, classification or observation helpers. Raw, continuity, physical, protocol/cut attribution and model support verdicts stay distinct.
- Work only in `/Users/jacob/.codex/worktrees/phase3b4-research/alan-research`, branch `codex/phase3b6-effect-observation`, base `df840103aa3a85e33b4b95ddcf2fe3df071cb024`. Do not edit other checkouts or `.envrc`.
- Offline environment: `GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b6/go-build`; use `/tmp/alan-phase3b6/` for temporary data.
- One implementation agent at a time; each task commits, self-reviews and receives independent spec/quality review. No worker subagents. Preserve real failures, exact commands/exits and source identities.

## Review Focus

- A launched packet can ingress and apply after loss even when no ingress was observed at the earlier cut (Tasks 1–2).
- A response destined for the dead incarnation cannot become a result of the reopened owner (Tasks 1–2).
- Endpoint application-before-loss requires matched acknowledgments and pinned producer ordering; final journal truth alone is insufficient (Task 3).
- Renewal must cover both principals and every unchanged durable field/table; a PID or epoch change alone proves neither loss nor exact renewal (Task 3).
- Model witness choice is not observed chronology, and a modeled equal-observation pair is not an experimentally captured v2 pair (Tasks 3–4).

---

### Task 1: Define and enumerate the split effect boundary

**Files:** Create `research/effectobservation/model.go`, `explore.go`, focused tests, `cmd/export/main.go` and tests; `docs/foundational/phase3b6/02-model.md`, `evidence/task1/`.

**Interfaces:** Produce `Initial() State`, `Successors(State) []Edge`, `ObserveStore(State) StoreView`, `ObserveDossier(State) DossierView`, `Explore() Graph`. `Edge` carries a complete `Event` and target `State`; `Graph` exposes deterministically BFS-ordered `States`, indexed outgoing `Arcs`, and one shortest `Prefixes` event sequence per state. States and observations are comparable value structs. Use distinct typed owner stage, event kind and authentic result symbols; no JSON byte-slice/base64 result encoding. Event descriptors include fixed binding, current actor/incarnation, exact root and value as applicable. Declare full field meanings in the model note. New exporter takes no arguments and emits complete deterministic JSON graph, checking writes and usage errors. Oracle APIs arrive in Task 2.

- [x] Read the contract and preflight recon. Add `TestSeparateStages`, `TestLossPreservesExternalPacket`, `TestLateResponseCannotReachReopenedOwner`, `TestExactRootAndHistoricalAuthority`, `TestRefusalAfterUnretainedApplication`, and `TestRetainedResultSurvivesLoss`. Record genuine failing assertions before production implementation; a missing-symbol compilation failure alone is insufficient RED evidence.
- [x] Implement exactly the post-admission finite domain: launch, ingress, apply either authentic value, response, record, consume, revoke-b, one owner loss/reopen and fresh refusal. Reject malformed descriptors in any trace application helper. Keep fixed first allow/debit/grant and c control explicit. No retry, lookup, atomic Apply/Retained, per-agent crash expansion or old-package changes.
- [x] Test all reachable state invariants, transition-stage separation and deterministic BFS/shortest-prefix reconstruction in `TestCompleteGraph`. Enumerate once and retain graph/summary hashes plus boundary witnesses; label counts only after execution. Write conditional transition arguments and what remains outside the model.
- [x] Run focused package and command tests, graph export twice, formatting and scoped Staticcheck. Inspect full diff and commit task-only files/evidence. Report exact API schemas for the next worker, commands/exits, observed counts, source identities and limitations. Do not claim independent audit yet.

### Task 2: Independently audit observations and preserve falsifiers

**Files:** Create `research/effectobservation/oracle.go`, `observations.go`, focused tests/corpus helpers, `cmd/audit/main.go` and tests; `docs/foundational/phase3b6/03-observations.md`, `evidence/task2/`.

**Interfaces:** Consume Task 1 descriptors/graph without operational decision helpers in the ledger. Produce `AuditTrace(events []Event) (Ledger, error)`, `AuditGraph(graph Graph) (AuditReport, error)`. The separately represented ledger exposes independently derived Store/Dossier views for comparison. `AuditReport` includes checked states/arcs/merged targets, independent observation-class counts, and shortest qualifying witnesses/pairs with exact declared search criterion. New no-argument audit command emits this report and hostile-control results deterministically. Task 3 can call `AuditTrace` and compare its independently derived views with the operational witness cut.

- [x] Add RED/GREEN tests for an independently interpreted valid event prefix and each malformed/stale/wrong-root/wrong-value transition. Validate irrelevant event fields, fixed identities, conservation, old permit survival, external packet survival, absent late-response delivery and authentic local custody without calling operational helpers.
- [x] Audit every outgoing graph edge, shortest prefix and merged-target ledger equality. Independently derive/group observations and compare with model projections in `TestIndependentObservationPartition`; retain shortest revoked/reopened/refused equal-Store zero/one-application pair, distinct unretained-value pair, post-loss negative-cut/late-application path, and retained-result recovery control. Specify shortestness relative to the initial cut and tie-breaking. Keep absence-at-cut separate from future nonapplication; no completion theorem.
- [x] Add semantic hostile traces for gate/contact/application/retention conflation, loss erasing packet/world truth, revived old response/permit, b/c aliasing, rewritten historical allowance, fresh allowed retry and forensic-result import. Run compiled independent-ledger mutants from disposable copies for at least packet deletion on loss, application-as-retention, old-response acceptance after reopen, and c revocation. Each must compile and be killed by a meaningful targeted assertion, with lawful controls retained.
- [x] Run complete graph audit and exports twice, focused tests, formatting and scoped Staticcheck. Preserve commands, hashes, assertion failures and any discovered counterexample. Commit task-only files and self-review. Document conditional observation arguments separately from finite checks and old W3/Phase2B/2C results.

### Task 3: Certify the sealed loss and pretransport cuts

**Files:** Create `research/delegationcorrespondence/effect_observation.go`, narrowly separated protocol/mapping helpers and tests, `cmd/effectobservation/main.go` and tests; `docs/foundational/phase3b6/04-correspondence.md`, `05-runbook.md`, `evidence/task3/`. Reuse existing raw/physical helpers without changing legacy results.

**Interfaces:** Produce `CheckEffectObservationDirectory(dir string) (EffectObservationCaseReport,error)` and command `-case DIR`. Preserve full legacy `CaseReport`; separately report raw/continuity/physical/protocol verdicts, explicit unsupported reasons, per-cut causal certificate, projected observation, existential model event witness and independent audit. Consume Task 1/2 model and ledger. Domain: exactly sealed loss and N07 normal/race; pre-admission loss prefixes remain outside-domain entries.

- [x] Read contract and evidence recon, then independently verify pinned archive/member hashes and relevant frozen producer source from runtime Git. Archive compact source extracts/hashes establishing temporal order and loss/renewal semantics. Read JSON/JSONL only; never SQLite. Add genuine RED positive/negative tests before implementing this relation; preserve legacy loss exclusion.
- [x] Require pinned producer commit/tree/source hashes, exact profile/schedule/manifest and ordered captures. Check every raw binding/custody and continuity; require both endpoint-ready records, journals and complete final snapshot. Join exact BA request versus service identity, root/path/grant, immutable wire/key/auth/reservation/debit, sole allowed gate and any fresh refused gate. Use one-contact uniqueness instead of an invented gate ordinal in endpoint wire.
- [x] For loss, certify both matching barrier acknowledgments before confirmed kill under the pinned producer; verify exact in-flight authorization and endpoint commit binding, successful authentic receipt withheld before the controller release after death (timeout/error completion remains unobserved), one shared loss, exact two-principal renewal with full table/run coverage, accepted exact b revoke after reopen and refused current-epoch gate after revoke. For N07 certify pretransport callback cuts and separately retained final application/result/consumption. Keep endpoint knowledge at cut and final dossier corroboration distinct. Build coherent extending abstract witnesses, marking inferred prerequisites/witness order, and audit independently.
- [x] Add coherent rehashed attacks for missing/wrong ack phase/key/auth/source/process/token, missing runtime in-flight authorization, wrong kill/wait, multiple loss, one-principal or non-exact renewal, moved/aliased revoke, stale fresh gate, duplicate physical ingress/commit, missing endpoint coverage, lifecycle settlement as effect receipt, forged custody, final-only evidence leaked backward and changed producer protocol. Keep gate-only and genuine retained-result positives. Run a compiled causal-ack mutant with an end-to-end raw/physical-valid false-acceptance assertion, and a compiled exact-renewal helper mutant killed by direct renewal assertions; retain malformed table/run dossier tests at their actual inherited rejection layers.
- [x] Export all four reports twice, preserve per-cut verdicts and hashes, and compare all ten old revocation reports with archived Phase3B5 bytes. Run focused tests/command failures, formatting and scoped Staticcheck. Commit only task artifacts. Report no new runtime campaign and no actual matched v2 zero/one local-equality pair.

### Task 4: Qualify and reconcile the bounded result

**Files:** Create `docs/foundational/phase3b6/06-qualification.md`, `07-results.md`, `08-reviews.md`, manifests and `evidence/task4/`; update `paper/alan-dissertation.md`, `.tex`, rebuild `.pdf`; bibliography only if necessary.

**Interfaces:** Consume reviewed Tasks 1–3 and controller's exact-source qualification evidence. Produce synchronized bounded claims/counterexamples, source/corpus/toolchain manifests, reproducible results and visually checked PDF. Controller owns broad verification, final review, workflow archive and delivery.

- [x] Controller freezes executable inputs per group, then runs the full ordinary package set (unchanged group may run during Task3 after dependency-disjointness verification; affected group follows its reviewed freeze), affected race tests, builds/module checks, complete new model/audit/corpus, applicable old corpus regressions, all pinned legacy exports, formatting and scoped Staticcheck. Compare inherited global lint diagnostics separately. Record every named skip and genuine apparatus failure, direct exit and elapsed time. Requalify affected executable changes if needed.
- [x] Worker writes the bounded results and qualification account from actual artifacts, separating abstract observation arguments, finite graph evidence, concrete source-conditioned witness relation and historical runtime evidence. Reconcile affected Abstract/body/Conclusion in Markdown and LaTeX. Preserve prior atomic result and older epistemic precedents; state the missing matched v2 counterpart and one next bounded question.
- [x] Obtain independent task spec/quality review, then one broad whole-branch review with scoped correction/review. Read current PDF skill, build/inspect changed pages and continuations, preserve source/PDF hashes and inherited layout/Markdown-only divergence. Do not assert global parity or unexecuted checks.
- [x] Archive task briefs/reports/reviews, ledger/rulings and qualification manifests; verify every intended evidence member is tracked. Controller inspects final cohesive diff, rechecks PR46 status, pushes branch and opens/attaches a draft PR against the verified parent (or main if merged). No merge/release. Final report states result, tests, limitations, PR link and material rulings/costs.
