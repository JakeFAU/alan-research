# Task 3 report

Status: implemented and locally qualified, awaiting controller's fresh independent review. Base `77f897e0c8bde7c1b3ea736a19f26029debb588e`. Task-only implementation commit is identified in the final handoff; controller-owned contract/plan/preflight/task4 edits are excluded.

## Question and result

The implemented conditional relation asks whether the selected BA effect is historically admitted, externally applied, locally retained and actually consumed at each sealed cut across one shared owner loss and exact b revocation. It supports all four intended dossiers, retaining all sixteen captures: twelve supported cuts and four explicit pre-admission outside-domain entries. Loss witness lengths are 4/5/7; N07 lengths are 0/1/7. Every supported concrete projection agrees with operational observation and separately interpreted `AuditTrace` Store/Dossier views.

The two loss branches have one application before death, no locally retained authentic BA result, no BA principal consumption, exact two-principal reopen and final epoch-current refusal after exact b revocation. N07's two intermediate cuts are source-certified before worker launch and have zero ingress/application despite final application; final separately joins authentic receipt, effect settlement and actual principal consumption. This is an existential witness relation under an honest pinned producer, not a recovered total chronology, simulation/completion theorem, new runtime campaign or matched v2 equal-local zero/one pair.

## Files and interfaces

- New checker API `CheckEffectObservationDirectory(dir string) (EffectObservationCaseReport,error)`, with complete legacy report; separately named raw/continuity/physical/protocol/support verdicts; concrete bindings, per-cut causal certificates, observations, witness steps and independent ledger.
- New command `research/delegationcorrespondence/cmd/effectobservation -case DIR`, exits 0/1/2 for supported/checked unsupported/tool or usage failure.
- Focused files: `effect_observation.go`, `effect_protocol.go`, `effect_mapping.go`, corresponding tests, command/tests and `testdata/effect-observation/anchors.json`.
- New notes `04-correspondence.md`, `05-runbook.md` and `evidence/task3/` with literal producer source extracts, commands, deterministic complete reports, hostile/mutant/RED logs, source and artifact hashes.
- No legacy checker/model/ledger files, production runtime, dependencies or manuscript files were edited. Task 4 owns manuscript/PDF reconciliation and broad qualification.

The report's `Producer.SealedBytesMatch` is informational archive membership, distinct from the conditional semantic relation. The production API uses no hardcoded extraction directory. Optional new corpus tests use `ALAN_EFFECT_OBSERVATION_CORPUS`; old tests retain their separate `ALAN_REVOCATION_CORPUS` setting.

## Evidence executed

`evidence/task3/commands.json` records final focused commands, direct statuses, timing and offline/cache environment. Toolchain was Go 1.27.4 darwin/arm64 and Staticcheck 2026.2.2 (0.8.2).

- Independently verified pinned commit/tree, 17 source blobs, 31 literal extracts, 16 binding source hashes and every source.json hash against frozen Git; both raw archives against owning archive metadata; all 220 regular members of four dossiers against archived hashes, sizes and extraction. Binary members were only byte-hashed; no SQLite access or runtime execution.
- Genuine compiled RED stub: four positive assertions and a missing-endpoint-ack false acceptance failed. Original source/log archived; direct reproduction over frozen Phase 3B.5 source exited 1 with semantic failures. Initial shell returned 0 only because it printed the failing Go log afterward; this apparatus detail is disclosed.
- Final `go test -count=1 -v ./research/delegationcorrespondence ./research/delegationcorrespondence/cmd/effectobservation -run '^TestEffectObservation'`: 12 top-level tests and 50 subtests passed, no skips. Subsequent command-only serialization delta passed and source hashes were refreshed.
- Scoped Staticcheck `./research/delegationcorrespondence/...`, plus command-only delta after serialization-test refinement: passed. Formatting passed. Built new command and existing revocation/raw commands.
- All four new reports exported twice with identical bytes, with every cut/hash/verdict retained in `verdicts.json` and compressed complete reports.
- All ten old revocation reports and all ten raw reports matched archived Phase 3B.5 hashes and compared byte-for-byte with reports reconstructed from frozen `df840103aa3a85e33b4b95ddcf2fe3df071cb024`. Individual old report bytes were not archived, so comparison explicitly reconstructs bytes and verifies their archived digests. Legacy loss exclusion remains unchanged.
- Compiled baseline-positive causal mutant: omitted causal ack checks, falsely accepted the rehashed missing-endpoint-ack dossier, and failed the targeted semantic assertion. All genuine positives passed.
- Compiled baseline-positive exact-renewal mutant: omitted helper body, failed seven direct helper assertions covering run omission/addition, one epoch, changed run field, table omission/addition and changed non-run row. This is helper-level false acceptance, not complete-checker end-to-end acceptance; old raw/physical checks already reject those complete hostile dossiers.

Hostile coverage includes ack phase/key/authorization/source/PID/token/Store/scenario, absent/wrong in-flight authorization, kill/wait, multiple declared loss, one-principal renewal, exact helper coverage, moved/aliased revoke, stale fresh gate, duplicate ingress/commit, coherent second ingress returning the same idempotent receipt, missing endpoint coverage, lifecycle-as-external-receipt, forged custody, final-only facts moved into an earlier capture, changed producer/source/schedule and changed fixture program. Required semantic-only attacks assert raw/continuity/physical validity before the new rejection. Other tests preserve their actual lower-layer rejection instead of crediting the new protocol checker.

## Source-conditioned causal boundary and rulings

The loss source consumes both acknowledgments before confirmed SIGKILL. Endpoint commit/journal precede endpoint ack. Successful authentic receipt emission requires a successful post-commit barrier release, performed by the controller after death. Cancellation can return an error earlier; continuous occupancy/transport liveness and acknowledgment emission order are unobserved. The certificate does not strengthen these facts.

Open contains PID/namespace/runs, not runtime Store/token. Runtime Store/token therefore bind to the separately declared sealed-dossier process anchor; they are not falsely attributed to Open or described as general authentication. Endpoint identity independently joins ready/journal/receipt and exact wire. The controller approved this scoped interface interpretation.

The recon/handoff initially missed existing `external.go:508–514` inline reopen renewal. Existing raw schema, continuity and `checkOwnerRenewal` already reject one-principal and table/field shape attacks. The new helper expresses the full exact table/run contract explicitly. Settings/Roots equality was investigated and removed as outside that contract. Controller updated the brief/contract to distinguish end-to-end causal and helper-level renewal mutants. Investigation logs and `investigation.md` preserve apparatus/candidate outcomes.

Concrete epoch2→model1 and reopened epoch3→model2 are explicit; immutable wire is never rewritten. The normal N07 epoch1→model1 map is separate. Final endpoint corroboration is labeled separately from at-cut knowledge. No additional physical ingress is inferred from post-reopen scheduler consideration; source reports one dispatch, not zero dispatches.

## Self-review and remaining obligations

Self-review inspected the complete new code/command/test changes, verified anchor/source/generated evidence consistency, and kept final-only observations separate. It added the manifest-to-every-capture program hash join and a coherent rehashed program-change rejection. Full legacy bytes were preserved. No unrelated task/model/ledger/runtime changes were included.

Fresh independent Task 3 review, broad/race/build/module qualification and manuscript/PDF reconciliation remain controller work. No claim of independent approval is made here. Archive exports remain unsigned and provenance assumptions remain material. No whole-pair recovery, arbitrary endpoint, non-idempotent retry, general distributed ownership or unrestricted composition claim follows.
