# Phase 3B Closure Assessment Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Determine whether Phase 3B's bounded investigation can close by mapping every original obligation to exact evidence and remaining limits.

**Architecture:** Audit the existing Phase 3A/3B.1–3B.7 record, then produce one assessment and synchronized research-status text. Qualification is read-only evidence inspection plus document/PDF checks, with no new runtime or model behavior.

**Tech Stack:** Existing Markdown/LaTeX/BibTeX, Python standard library for audit records, cached Tectonic/Poppler; no new dependencies.

**Spec:** `docs/foundational/phase3b-closure/00-contract.md`.

## Global Constraints

- Verdict is evidence-driven, not predetermined. Original unmet requirements cannot be hidden as new exclusions, and unsupported generalizations do not automatically require more experiments.
- No production/runtime/model/checker behavior, dependencies or historical evidence changes. No new runtime/model campaign or test-suite pass may be claimed.
- Preserve atomic versus split-effect domains, source-conditioned correspondence, original failures/reviews and the Phase 3B.7 post-capture checker correction.
- Never open retained SQLite or read/change the runtime `.envrc`.
- Use existing research worktree on `codex/phase3b-closure-assessment`, base `7ae2c5ddd5e18d2907dd91fcefc896456a25562d`; temporary outputs under `/tmp/alan-phase3b-closure/`.
- One implementation agent at a time; no worker subagents. Separate task and whole-branch review; push a draft PR and do not merge.

## Review Focus

- Missing positive new-invocation or ownership transitions must be distinguished from tested stable identities and explicit unsupported-operation refusals.
- An atomic application/receipt model plus a split-effect witness must not be presented as a single proven refinement.
- Safety, possible completion, fairness-dependent liveness and quiescence must remain distinct.
- Stale roadmap prose must not erase the original gates or retroactively strengthen historical evidence.
- A bounded closure verdict must state its domain while assigning each actual unresolved obligation precisely.

---

### Task 1: Audit original gates and reconcile the closure verdict

**Files:** Create `docs/foundational/phase3b-closure/01-assessment.md`, `02-verification.md`, `03-reviews.md`, `evidence/`; update `docs/foundational/phase3/README.md`, `paper/alan-dissertation.md`, `.tex`, `.pdf`. Bibliography only if needed. Controller owns final workflow/archive/delivery metadata.

**Interfaces:** Consume the contract and read-only reconnaissance at `/tmp/alan-phase3b-closure/early-evidence.md` and `late-evidence.md`, verifying decisive claims against owning sources. Produce a complete evidence matrix and an unambiguous bounded closure disposition; keep any unresolved core gap visible. Controller and reviewer consume exact audit source hashes, commands and report.

- [ ] Read original exit criteria, progression gates and the relevant source-owned contracts, results and machine-readable evidence. Record sources inspected and exact semantic disagreements rather than relying on recon conclusions alone.
- [ ] Write the assessment with a row for every original gate, explicit rely/guarantee assumptions, evidence tier, negative witness and disposition. Test the proposed closure against the five review-focus hazards. State whether any original mandatory gap blocks closure, and distinguish it from Phase 3C and later extensions.
- [ ] Update stale roadmap status without rewriting the original gates. Synchronize the manuscript's affected claims and status; preserve theorem numbering, historical result qualifications, bibliography and inherited divergence. Read PDF skill before the new artifact operation; build and visually inspect changed pages/continuations, recording hashes/warnings.
- [ ] Verify local links, matrix coverage, input existence/hash stability, MD/TeX alignment and unchanged executable/old evidence inputs. Record this as an audit, not new runtime qualification. Inspect full diff, commit narrowly and report actual commands, conclusion and remaining risks.
- [ ] Controller obtains independent task spec/quality review and one broad final review, fixes substantiated findings with re-review, archives workflow and pushes/attaches a draft PR. No merge. Final answer states closure decision, evidence, changes and exact remaining obligations.
