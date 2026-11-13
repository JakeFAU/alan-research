# Phase 3B.6 Task 3 independent review

Spec compliance: **APPROVE**.

Task code/evidence quality: **APPROVE**, with one nonblocking Minor wording finding below.

Reviewed base `77f897e0c8bde7c1b3ea736a19f26029debb588e` through head `b7ed93186a16a0a52637de7b7e3a13a5fd37e523`. This verdict is for Task 3 only. It is not whole-branch approval, proof of archive authenticity, or permission to merge. The current contract clarifications and reconnaissance erratum were read as the requested review inputs, rather than treating prior controller rulings or approval as evidence of correctness.

## Findings

No Critical or Important findings.

**Minor — archive/dossier count in serialized trust statement.** `research/delegationcorrespondence/effect_protocol.go:75` says “four unsigned sealed archives.” There are two owning archives (`raw-normal.tar.gz` and `raw-race.tar.gz`) containing the four selected dossiers, as the source index and correspondence document correctly describe. This does not affect acceptance or provenance checking, but the serialized report misstates its evidence inventory. Correct the phrase to “two unsigned sealed archives containing four selected dossiers” when updating this text, and refresh affected serialized report evidence/hashes if changed. This editorial discrepancy does not block Task 3 approval.

## Specification and semantic checks

- The new API and command are opt-in. The complete legacy report is retained, and raw, continuity, physical, protocol and support verdicts remain separate (`effect_observation.go:78–136`). All sixteen declared captures remain present; the four pre-admission loss prefixes have explicit outside-domain reasons. The task diff adds only task-specific checker/command/tests, anchors and research evidence; no pre-existing checker/model/runtime source is edited.
- Producer selection checks pinned commit, source manifests, exact sealed fixture manifest and every capture's program hash. It requires the declared profile, schedule and exact capture order. The archive-member comparison is informational and does not substitute for semantic checks (`effect_protocol.go:34–118`). The harmless re-encoding positive and semantic-only hostile assertions exercise that distinction. Unsigned archive provenance and honest pinned production remain premises, not consequences of editable hashes.
- Loss certification binds endpoint ack phase, PID/token/Store, exact wire key/authorization/source and unique committed contact. Runtime ack must identify the selected sole in-flight authorization and match kill/Open/capture PID plus the separately declared sealed process anchor (`effect_protocol.go:120–180`). Open does not contain Store/token; neither the code nor the explanatory document claims otherwise. I checked the source extracts for process pipe configuration, ack validation, controller order, kill confirmation, endpoint commit/journal/barrier/response and caller response validation. They support both acknowledgments before confirmed death and no successful authentic receipt before the post-death release. They do not establish uninterrupted occupancy/transport or ack emission order, and the certificates explicitly preserve those limits.
- Exact renewal checks the entire table-name set, all non-run rows, both runs and every run field except the exact epoch 2→3 change (`effect_protocol.go:182–214`). Focused inspection of unchanged `external.go:397–418` and `external.go:508–514` confirms the erratum: inherited validation already checks inline stopped→reopened renewal and rejects one-principal renewal. The helper mutant is correctly characterized as helper-level evidence, not an end-to-end acceptance gap. Settings/Roots equality is not imported into this contract.
- The mapper checks the exact BA request/service distinction, grant/root/path, authorization/wire, reservation/debit, first allowed gate, exact b control and fresh current-epoch refusal. Lifecycle control count and exact root checks prevent an additional c revocation. Original attempt and historical authorization/debit are preserved across captures (`effect_mapping.go:7–143`). The source incarnation is mapped explicitly rather than rewriting immutable wire coordinates.
- Local retention and actual principal consumption are joined to their own raw evidence and ordered facts (`effect_mapping.go:121–138`). I inspected the unchanged custody settlement checks (`custody.go:165–210`) to resolve the specific concern that counting settlements might accept a lifecycle settlement as an external receipt. Raw validation independently joins the external attempt, authentic retained receipt, disposition and resource-receipt observation; lifecycle records have a different checked form.
- N07's zero-contact intermediate cuts follow the synchronous callback/admission/worker-launch source order. Final receipt corroboration is separately labeled. Final-only endpoint truth is not projected backward into the two intermediate Store/Dossier observations (`effect_mapping.go:177–219`, `effect_observation.go:185–235`). A single successful gate and unique endpoint contact provide attribution; no endpoint gate ordinal is invented.
- Witnesses extend coherently and label prerequisite launch/response steps. Their order is existential, not a recovered physical chronology. Both operational observations and independent ledger projections must match the concrete projection (`effect_observation.go:149–207`). The named dependency check of `research/effectobservation/oracle.go` confirms that `AuditTrace` and its projections interpret ledger state without calling operational successor/observation/authority helpers. This is a focused integration check, not a repeat review of all Task 1/2 graph evidence.

## Evidence and review work

I read the complete supplied 5,662-line diff in bounded chunks, including the source extracts, source index, anchors, qualification scripts, command/status records, test/mutant logs, summaries and compatibility records. Compressed historical attack logs were inspected, including the retained constructor panics and the originally mistaken renewal-layer expectation. The RED archive contains the unconditional-support stub and its positive/causal negative tests; its reproduced log records semantic assertion failure with direct exit 1. The historical shell-status limitation is disclosed.

The recorded final focused log contains 12 top-level tests and 50 subtests without skips. The causal mutant compiles, preserves all four genuine positives and falsely accepts the missing-endpoint-ack dossier; the exact-renewal mutant compiles, preserves the positives and fails seven direct helper assertions. Current tests distinguish semantic-only ack/source/schedule/revoke attacks from raw/physical/continuity rejection; the coherent repeated ingress test requires the actual physical rejection edge. No malformed-constructor panic is represented as current qualification success.

No Go tests, race suites, builds, Staticcheck runs or runtime campaigns were rerun in this review. Their execution claims are attributed to the inspected retained command/log evidence. I performed one focused artifact-only Python audit to resolve whether the compressed reports and embedded anchors actually correspond to the recorded evidence:

- All 69 `SHA256SUMS` entries matched their files, and all eight frozen checker/test/anchor source hashes matched current task source.
- All four compressed report byte counts/digests and every summarized capture field matched `verdicts.json`.
- Embedded commit/tree/source hashes matched the source index. All four embedded manifest/digest/source/binding anchors and both runtime process anchors matched JSON read directly from the owning tar archives; binding hashes also matched the index's normalized source paths.

The artifact assertions completed successfully. An appended attempt to display every RED tar member as UTF-8 then failed on a macOS `._` metadata entry; a subsequent read selected the actual `.go` members and displayed both successfully. This was a review display issue, not a failed semantic/evidence assertion. No SQLite was opened. The review made no branch/index changes and touched no controller-owned files.

`verify_task3.py:17–32`, its command records and `legacy-comparison.json` document reconstruction from frozen Phase 3B.5 source, equality with archived old digests, and byte comparison of all ten old revocation plus ten raw reports. The code diff leaves old checker behavior untouched. This is accurately described as comparison with reproduced bytes authenticated by historical digests; original individual report files were not archived. Four new reports were exported twice according to the retained records, with 12 supported and four outside-domain cuts.

## Cannot-verify items and controller obligations

- Unsigned artifacts cannot independently authenticate that the real historical producer executed the claimed process protocol. The source-conditioned relation visibly assumes honest pinned production and sealed-archive provenance. Rehashed semantic acceptance must not be presented as newly authenticated runtime evidence.
- I did not independently rerun the frozen-source/archive verification script or the already evidenced qualification suites. The source verification script and retained result were reviewed; the focused artifact audit above independently checked the embedded-anchor/report joins. Broad current-source ordinary/race/build/module qualification remains Task 4 work.
- Manuscript Markdown/LaTeX/PDF reconciliation, inherited divergence disclosure, whole-branch review and final draft-PR/base status are outside this Task 3 verdict and remain controller obligations.
- There is no actual matched v2 equal-local zero/one application pair, new runtime campaign, general simulation/recovery theorem, arbitrary-agent closure or universal progress claim established here. The correspondence notes and report limitations preserve those boundaries.

## Scoped re-review: inventory wording correction

Reviewed correction `b7ed93186a16a0a52637de7b7e3a13a5fd37e523..566b5e3ec1a2b23f6c66bc025a0e0ad447a3b05f` only, against the supplied delta package and appended implementer report.

Minor inventory finding: **ADDRESSED**.

Spec compliance: **APPROVE**. Task code/evidence quality: **APPROVE**. No new findings in this delta.

`research/delegationcorrespondence/effect_protocol.go:75` now correctly says “two unsigned sealed archives containing four selected dossiers.” The only research-source change is this string literal. I inspected the correction metadata, refresh script, retained command/test/build records and before/after report comparison. Existing command tests and build are recorded as passing; four reports were exported twice by the implementer. I did not rerun those commands or any suite.

A focused read-only artifact comparison independently confirmed that the source and each of the four decompressed reports differ from their preceding immutable Git bytes by exactly one replacement of the reviewed phrase. Prior source/verdict manifests are exact copies of their prior-commit versions; the old source and evidence-manifest hashes match correction metadata; current report hashes/lengths and every capture-summary field match, with all capture summaries unchanged. The current checksum inventory also matches its files. Historical evidence remains attributed to its original source rather than being represented as new executions. No checker guard, model projection, anchor, witness or embedded legacy report changed. No checkout/index edits, SQLite access, runtime execution or subagents were used.

This closes the sole Minor finding. The original Task 3 review's provenance limitations and Task 4/controller obligations remain in force.
