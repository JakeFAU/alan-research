# Phase 3B.2 authority attenuation implementation plan

> **For agentic workers:** Use superpowers:subagent-driven-development for the independent tasks and scoped reviews. The user explicitly authorized completing this milestone and a draft PR without merge or release.

**Goal:** Determine when reciprocal service preserves rooted, non-amplifying authority and completion reachability.

**Architecture:** Extend the unchanged reciprocal finite model with explicit grant/decision/admission records. Validate serialized histories with an independently authored authority ledger. After a decisive model candidate, add an ordinary same-Store fixture and real SIGKILL qualification, preserving missing-delegation and unknown-outcome implementation boundaries.

**Tech stack:** Existing Go 1.27.3, standard library finite model, existing author/runtime/effects/SQLite and subprocess barriers, Tectonic manuscript build.

**Spec:** `00-brief.md` and frozen `01-contract.md` in this directory.

## Global constraints

- Exactly two agents and two original service obligations; A servicing/B serial in the extension. Keep all Phase 3A/3B.1 files and archives unchanged.
- Complete operation/resource/scope tuples a,b,r,x; roots A={a,r}, B={b,r}; x initially unheld. One-hop explicit grants, no redelegation; at most one b-root revocation and one recoverable crash per actual agent.
- S1 server roots; S2 explicit attenuated original-bound delegation; S3 S2 plus ambient requester-root borrowing. Separate Cartesian mutant. Never conflate principal/consent violation with invention outside the global tuple union.
- Original decisions/refusals persist. Fresh physical admissions recheck validity; existing opportunity survives later revocation but is cleared by crash and consumed by Apply. Atomic application/receipt is a model premise, not generic adapter recovery behavior.
- No production engine, broad delegation framework, new dependencies, N-agent result, independent-owner recovery, or unknown-effect fabrication.
- Final executable qualification uses an immutable committed source revision; later evidence/manuscript commits must preserve those source hashes.

## Review focus

- The old Binding.Grant string or authentic result mistakenly serving as new authority evidence.
- Rootless ancestry, invalid cross-generation/request/operation reuse, and componentwise permission mixing.
- A historical allow bypassing current admission validity, or a late grant erasing a refusal.
- Revocation or missing authority silently excluded to manufacture universal completion.
- Source disposition, effect permission, physical admission and result custody conflated in runtime crash labels.

## Task 1: operational authority extension and decisive histories

**Files:** Create `research/authorityservice/model.go`, `explore.go`, `model_test.go`, `README.md`, and `cmd/export/{main.go,main_test.go}`. Preserve base package; no production imports.

**Interfaces:** Publish comparable `Config`, `State`, `Event`, `Transition` and report types. Functions `Configurations() []Config`, `Initial(Config) State`, `Successors(Config, State) []Transition`, `Explore(Config) Report`. Freeze serialized event/grant/authorization vocabulary before Task 2. Include full authority descriptors in events, not only final capability bits.

- [ ] Behavioral RED, then minimal local/explicit/missing cases. Report a concrete accepted/rejected paired history before the runtime implementation begins.
- [ ] Use actual base transitions; gate Apply; add grant/authorization/fresh admission/revocation events and physical-crash admission clearing. Preserve caller/callee availability and receipt custody.
- [ ] Enumerate S1/S2/S3 for Local/Explicit/Missing/ScopeMix, plus explicitly listed Shared, S2 LateGrant, S2 Revocable and Cartesian ScopeMix controls. Keep every reachable arc/cycle and reverse completion reachability; no depth cutoff.
- [ ] Retain unsafe-prefix evidence so graph merging cannot erase violations. Report authority-safe state/transition counts, raw and authority-safe completion reachability, shortest bad authorization/admission/application traces, and positive/temporal/custody traces. Infinite history families must not be reported as finite history counts.
- [ ] Positive S2 has its grant before first authorization; separate late-grant refusal remains permanent. Root revocation must not disable request commitment or rewrite historical outcomes.
- [ ] Produce deterministic machine-readable exports, genuine expected-failure controls with distinct apparatus errors, targeted Go tests/race/Staticcheck, and an exact source/report freeze.

## Task 2: independent retained-authority validator

**Files:** Create `research/authorityservice/oracle.go`, `oracle_test.go` and optional corpus files. Own no operational source. Read Task 1's event vocabulary, then independently implement its interpretation from the frozen contract.

**Interfaces:** `AuditTrace(Config, []Event) (Audit, error)` with explicit protocol validity, authority safety, violated conditions and per-original paid/unpaid status. It must not call Successors, operational authorization/chain predicates or transition validators. Define a strict acceptance helper if useful.

- [ ] Establish root ancestry, complete-tuple containment, exact original/generation/operation binding, issuance order, immutable decisions, revocation timing, activation-bound admission and application/custody relationships independently.
- [ ] Walk parent paths with a visited set; reject rootless cycles, unauthorized return delegation, mutated parents, enlarged tuple sets, stale/rebound originals, result-as-authority and scope-product laundering even when final counters/outcomes match.
- [ ] Validate accepted and rejected generated histories; state exact coverage. Where feasible, cache checked BFS prefixes for every graph edge; otherwise explicitly identify the smaller independently checked witness/case set without claiming exhaustive independent validation.
- [ ] Preserve malformed controls and expected reasons, exercise both result values/equal values without losing provenance, and add focused tests/race/Staticcheck.

## Task 3: public runtime and real process interruption

**Files:** Create `examples/communication/authority/` with ordinary deployment/adapter and construction tests; create focused `internal/durabletest/phase3b2_*_test.go` using existing barriers/backup/process helpers. Keep prior fixtures unchanged.

**Interfaces:** Public author/runtime/effects for positive execution; existing internal boundary only for observation and process synchronization. Exact constructors and registry names are implementer choices documented in the task report. Evidence goes to a task-specific temporary export directory before archiving.

- [ ] Start only after Task 1 yields a decisive paired authority candidate. Write adverse local-authority/refusal checks before construction.
- [ ] Two roots and asymmetric, fixed concrete adapter configurations; legitimate reciprocal service, unauthorized attempt, claim/result laundering, actual consumption and unpaid obligations. Filter missing permissions from Provision before Create; no positive database edits.
- [ ] Explicitly demonstrate the missing S2 primitive and mark delegation-retention cuts inapplicable. Do not substitute static permissions or specialized receipt grants for operation-bound delegation.
- [ ] Genuine SIGKILL + fresh-owner cases: before authorization, after source authorization, after allowed reservation before physical admission, after physical admission before call, after refused attempt, result retained before settlement, settled before consumption, mutual wait, and one original paid/one pending. Label applicable boundaries precisely; keep source authorization distinct from adapter permission.
- [ ] Check original identities/grants/decisions/refusals, fresh epochs, endpoint ingress/application counts and no reissue of completed operations. Generic authorized attempts without results must remain unknown/unpaid after reopen. Preserve raw programs, observations/history, retained permissions/decisions, coherent backups, barriers, kill disposition and endpoint evidence.
- [ ] Disposable corrupted copies or retained-record mutations test rejection; never synthesize positive evidence. Run normal/integration/race/subprocess and Staticcheck gates.

## Task 4: qualification, argument, manuscript and delivery

**Files:** Add `03-results.md`, `04-runbook.md`, evidence and six scoped review records; update Phase 3 index, root research index and new mathematical content in both manuscript representations/PDF. BibTeX changes only for checked new attributions.

- [ ] Independent reviews separately address authority semantics/attenuation; model completeness/minimality; safety versus completion; retained recovery evidence; runtime correspondence/scope; manuscript/reproducibility. Challenge smaller traces/missing edges. Preserve failed claims.
- [ ] Commit and identify immutable executable source. Run ordinary root suite, relevant previous/new normal/race tests, integration and SIGKILL suites, independent chain validation, formatting/Staticcheck, deterministic exports and hash checks against that source. Record exact commands/statuses/skips.
- [ ] Prove only the conditional rootedness/attenuation and progress claims warranted by the model. Distinguish finite enumeration, independent validation, proof argument, runtime observations and missing representation.
- [ ] Archive compact interpretable evidence with source/witness hashes, rebuild PDF, normalize new Markdown/TeX content and visually inspect changed pages.
- [ ] Whole-change review; commit evidence/docs, verify executable hashes unchanged, push and open a draft PR against verified main. Report base and source/delivery commits, exact counts and remaining boundaries. Do not merge or release.
