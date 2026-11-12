### Spec Compliance

- ✅ Spec compliant for Task 1. The opt-in constructor and validity route preserve the original configuration list; `ABOwnPermission` is appended and only the new configuration selects it (`research/authorityservice/model.go:28`, `:57`, `:106`). Export selection preserves default scope (`research/authorityservice/cmd/export/main.go:30`).
- ✅ Separate issuance is commitment-gated, live-B-only and single-attempt, with immutable accepted/rejected status; revocation changes only `Revoked` (`research/authorityservice/model.go:427`). The tests cover both issuance outcomes, refusal history, surviving and crashed permits, genuine reciprocal completion, and durable authentic consumption (`research/authorityservice/revocation_test.go:27`, `:76`, `:93`, `:107`, `:126`, `:144`, `:189`).
- ✅ New witness cuts explicitly select BA application and A's revoked BA permit crash; legacy cut behavior is guarded by the new configuration check (`research/authorityservice/explore.go:215`, `:237`). Tests replay complete event descriptors and compare exact target states (`research/authorityservice/revocation_test.go:241`).
- ⚠️ Independent ledger/reachability qualification, runtime correspondence, all-four historical export regression, and manuscript reconciliation remain later-task obligations. Task 1 correctly identifies these limits rather than claiming them (`docs/foundational/phase3b5/02-model.md:63`). The Task 1 default authority-service export comparison is recorded in `docs/foundational/phase3b5/evidence/task1/export-sha256.txt:1`; I did not regenerate exports.

### Strengths

- The implementation is a small extension of the existing transition system. Issuance status and revocation remain separate facts; no new state/event fields, dependency, runtime change, or custody transition is introduced (`research/authorityservice/model.go:427`).
- The AB positive witness checks the distinct complete c-root permission, observes the existing B serial guard both before and after BA recording, then consumes different authentic result values for the two originals (`research/authorityservice/revocation_test.go:144`). This tests the intended semantic distinction rather than merely allowing an AB authorization.
- The adverse histories discriminate historical authorization from current admission and pre-revocation permits from post-crash opportunities (`research/authorityservice/revocation_test.go:76`, `:107`). Late accepted issuance cannot repair a prior refusal (`:126`).
- Complete exploration checks the candidate predicate at every reachable state while preserving losing safe states and retry cycles (`research/authorityservice/revocation_test.go:215`). The documentation labels its argument a source-level sketch and explicitly bounds the finite relation (`docs/foundational/phase3b5/02-model.md:15`, `:25`, `:63`).

### Issues

#### Critical (Must Fix)

- None found.

#### Important (Should Fix)

- None found.

#### Minor (Nice to Have)

- None found.

### Assessment

**Task quality:** Approved.

**Reasoning:** The opt-in changes preserve the original semantics for old bindings and configurations while adding the requested separate-grant/exact-root temporal behavior. Tests target meaningful positive and negative boundaries, including the unchanged serial guard and complete witness reconstruction; no blocking correctness or maintainability issue was found.

### Evidence and review limits

- Reviewed `.superpowers/sdd/phase3b5-plan/review-8971eb11..2107c44c.diff` once as the change view, against the task brief, global constraints, contract, and implementer report. Review identity: base `8971eb11`, head `2107c44c16eb4a8cb7f67122d734cd995cb953ec`.
- The model diff cuts off `Successors` and descriptor/authority functions needed to evaluate the changed issuance rule. I read `research/authorityservice/model.go:180–445` to complete that context. Named risk: recomputing current grant validity could accidentally rewrite historical descriptors, or inherited apply/crash handling could invalidate the claimed permit boundary. The inspected `describe` emits a grant only at issuance, selected witnesses remain stable, apply consumes a permit, crash clears the server's permit, and revocation only flips `Revoked` (`model.go:315`, `:372`, `:427`). No broader code crawl was performed.
- Decompressed and inspected retained RED, cut-RED, GREEN, old-model and Staticcheck logs, plus the exported report. The GREEN log and report agree on 229,834 states and 1,807,246 arcs, with zero reported operational/custody/authority violations. The new graph and descriptor checks share operational machinery and are not independent semantic validation (`docs/foundational/phase3b5/02-model.md:29`, `:63`).
- Preserved Staticcheck cache failures are environmental failed attempts, followed by a documented successful cache override and an empty final `staticcheck.log.gz`. The passing test logs contain ordinary verbose test evidence; no unresolved warning was identified. These are inspected historical results, not commands rerun by this reviewer.
- No tests, suites, exports, Git commands, or source mutations were performed. Only this requested review report was written. No focused execution was necessary because inspection exposed no unanswered behavioral doubt requiring a new probe.
