# Phase 3B.7 Task 3 results and manuscript report

## Scope and result

Documented the two genuine sealed pairs and their offline checks in
`04-runbook.md`, `05-results.md`, `06-qualification.md` and `07-reviews.md`
under `docs/foundational/phase3b7/`. The results use the frozen and corrected
command records and report exact archive counts, observed coordinates,
application facts and model trace lengths. The initial checker rejected both
pairs twice because it expected the constructor's ordinal 0 after the
transport wrapper had set ordinal 1. The post-capture checker correction
qualifies both original archives. The producer, raw and archive bytes and
observation definition were fixed. No captures or Go suites were rerun by
this worker. Controller's independent task review, whole-branch review,
workflow archive and delivery are pending at handoff.

The concrete result is equal local observations with selected application
counts 0 and 1 within each of two forked pairs. Each pair has one common
closed seed export and four complete byte-equal pairs of later exports.
Each complete final Run report matches its counterpart byte for byte. The
separate Phase 3B.6 existential witnesses use lengths 1/2/4 for zero and
4/5/7 for applied. The source-conditioned protocol argument and the
independent checker remain distinct from the finite model result and direct
export equality. The endpoint was idempotent. The claim is bounded to
forked seeds, two principals sharing a Store owner, sequential isolated
endpoints and the stated cuts. General recovery equivalence, universal
liveness, independent provisioning, concurrent cloned owners,
non-idempotent endpoint behavior, distributed recovery and N-agent closure
remain open.

## Verification performed

- Ran the PDF edit operation marker once before the first source edit.
  Node command used `--operation-kind edit --expected-output-count 1
  --output-format pdf`, exit 0, empty output, tool chunk `c883f5`.
- Synchronized Abstract, Phase 3B.6 closing paragraph, new Phase 3B.7
  subsection and Conclusion in `paper/alan-dissertation.md` and `.tex`.
  A local normalizer confirmed exact normalized equality for the changed
  sections. The Markdown-only Section 7 is inherited and unchanged. The
  bibliography is byte-identical to the base.
- Built from copied sources using Tectonic 0.17.0,
  `tectonic --only-cached --keep-logs alan-dissertation.tex`, exit 0 in
  1.007 s. Installed PDF is byte-identical to build output, 44 A4 pages,
  SHA-256 `7fc751e059f61423dc307e346826dd17099ac3d99c8a8787dd5d48fed17e8f76`.
- Ran Poppler layout text diff against the prior 43-page PDF. Visually
  inspected changed pages 1, 2, 41–44 and boundary page 40 at 110 DPI.
  No clipping, overlap or missing glyphs observed. The five unique inherited
  TeX box warnings and their source line ranges are unchanged.
- Read-only archive and report cross-check verified both 122-member
  inventories, lengths, digests and member modes. Compared four complete
  paired exports and each pair's complete final report directly from
  archive members, then counted nine distinct exports per archive and
  41 tables / 94 rows per final export. Compared all eight recorded
  checker outputs against their recorded digests and initial/corrected
  verdicts. SQLite files were hashed as opaque bytes.
- Extracted both archives into new disposable paths and ran the existing
  corrected checker binary twice per extracted pair. All four invocations
  exited 0 and repeated outputs matched exactly. Normal report digest
  `edf131b32f95075db836c86d16992f779fe03a18fbde327acba71dfe56aaa3e3`.
  Race report digest
  `4d7326b200b8bacc679eb0f3c1a57e16b3bab6628756855a82cb00a1e820c38f`.
  Checked runbook shell syntax and relative links, except the four workflow
  links awaiting controller archive.
- Archived the manuscript sources, build output, logs, rendered changed
  pages and visual/check records under `evidence/manuscript/`. Verified
  all 21 build archive members and compared archived paper inputs and
  PDF to the current workspace. Archive SHA-256
  `a4b1f9869ad393319223b2cebaa64ab69e35afd8b940a410a3c5937d19d96c11`.
- Read controller's final summary and 23-command package verification.
  Affected ordinary and race each record 422 pass events, seven tested
  packages, three no-test packages and no named skips. Build, module
  verification/tidy diff, scoped Staticcheck and current-source format
  check passed. Global ordinary/integration Staticcheck retained the
  identical historical 13/31 diagnostics. Prior repository-wide ordinary
  qualification is historical.

## Remaining verification and delivery

Controller still owns Task 3 independent review, whole-branch review,
workflow record archiving and PR. The PDF skill requests one standalone
output citation for the installed PDF in the final user-facing answer.
