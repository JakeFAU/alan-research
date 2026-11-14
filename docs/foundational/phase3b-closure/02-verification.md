# Phase 3B closure audit verification

This is an audit-only task. No production/runtime/model/checker behavior or dependencies changed; no model enumeration, Go test, capture campaign or historical acceptance checker was rerun. Retained SQLite was never opened. Historical numerical/qualification claims remain attributed to their owning packets.

## Commands and results

Commands below ran from the research worktree unless an absolute scratch path is shown.

- Source/contract inspection used `rg`, `sed`, `cat` and Python JSON/tar reads. Decisive retained text members and 47 owning files are bound by [inspected-inputs.json](evidence/inspected-inputs.json); the roadmap hash is explicitly from merged base `7ae2c5ddd5e18d2907dd91fcefc896456a25562d`, before the status edit. The [inspection script](evidence/inspect-inputs.py) reproducibly reads JSON/text only. It recorded four actual CancelRun traces, both Retry refusals, original serial unpaid work, model report counts and two matched pairs. This extraction is not a rerun of their historical semantic validators.
- `python3 docs/foundational/phase3b-closure/evidence/inspect-inputs.py`: exit 0; 47 files, 26 archive text members. All four full export pairs and the complete run report pair are byte-identical within each normal/race archive. [Extracted facts](evidence/inspected-facts.json).
- `python3 docs/foundational/phase3b7/evidence/qualification/freeze-inputs.py . --check docs/foundational/phase3b7/evidence/qualification/freeze-corrected.json`: exit 0; all **629 executable/acceptance inputs unchanged**. This is an input-stability check, not new qualification.
- `python3 docs/foundational/phase3b-closure/evidence/verify-audit.py`: checks all manifested inputs, original exit/gate text, all twelve matrix rows, local target existence, identical added MD/TeX prose, unchanged bibliography/Markdown-only Section 7, document-only tracked scope and the selected raw temporal/quiescence assertions. [Exact results](evidence/audit-checks.json). Link checking covers file/directory existence, not semantic validation of every linked historical artifact or Markdown anchor.
- `git diff --check`: whitespace check. Full assigned textual diff and final stat/status reviewed before the narrow author commit. PDF verified separately by rendering.

Exploratory reads initially used an incorrect contract filename, a nonexistent qualification filename and a nonexistent bundled `pdftotext` override path. JSON exploration initially assumed all events were objects and used `lifecycle_outcomes` instead of `lifecycle_results`; a later summary mistakenly looked for the result inside the request descriptor. These commands failed without modifying historical evidence. Corrected reads use mixed event types and the actual lifecycle `result` column; the final script/checks succeed. No experimental result was discarded or rerun to obtain acceptance.

## Manuscript

Added two closure paragraphs at the end of the existing matched-history subsection and a short Conclusion status paragraph, identically in Markdown/LaTeX. No theorem, subsection or bibliography numbering was changed. Roadmap status now links the assessment and dates the old primitive-absence statement to the Phase 3B.2 freeze; its original question, exit criteria and final gate paragraphs are unchanged.

Cached Tectonic rebuilt the PDF successfully. It has 45 pages; changed layout-text pages are 42–45. Pages 41–45 were rendered and visually inspected, including subsection/Conclusion and reference continuations. The same five inherited box warnings remain on unchanged lines 214–215 and 558–559. No new warning was introduced. [Visual/build record and hashes](evidence/manuscript/visual-review.md).

The inherited Markdown Section 7 “Formal Falsification” remains byte-identical and absent from LaTeX/PDF. Existing Markdown 6.15 versus PDF 6.10 and Markdown Conclusion 8 versus PDF 7 numbering remains. This is affected-text synchronization, not a claim of whole-manuscript parity.

## Scope and interpretation checks

Self-review separately considered all five hazards: unsupported invocation transitions, cancellation/revocation conflation, actual caller/server/endpoint roles, atomic versus split-effect domains, and quiescence/unpaid work. The central assessment leaves each limitation visible and distinguishes original Phase 3C requirements from optional extensions. The supported closure is a research disposition over named fragments, not a positive cross-product qualification. Independent reviewers and delivery remain separately tracked in [03-reviews.md](03-reviews.md).
