# Task 3 report: qualification and manuscript reconciliation

Status: DONE, ready for independent review. Controller owns review, workflow and delivery.
Base: `208e1411dd9074f4d97eb9fbba4fb9b8755b3236`.

## Question and result

Reconciled the accepted evidence for J = initial marked language plus minimal
residual profile. The manuscript/results preserve strict FULL coarsening with
legal-plant relevance; product closure preserving initial completion and
nonblocking; finite suffix-hit coverage as local necessity only; generic MIN
loss; MAX failure; generic J nonnecessity; all blocked new-fixture K/D contexts.
No weakest global interface, verification-cost reduction, runtime refinement or
general Phase 3C closure claim was added. Conditional source argument, exact
finite reconstruction, historical qualification and current identity checks
remain distinct.

## Changes

- New `docs/foundational/phase3c5/04-results.md`, `05-qualification.md`,
  `tools/verify.py`, `evidence/qualification/` and `evidence/manuscript/`.
- Updated Phase 3 roadmap, Markdown/TeX section 6.20/6.15 respectively,
  Abstract and Conclusion. Rebuilt `paper/alan-dissertation.pdf`.
- Added Ware/Malik (2011) bibliography entry after checking the primary arXiv
  abstract metadata. Controller's deeper primary inspection remains attributed;
  no complete proof or imported theorem is claimed.
- All 760 accepted Task 1/2 source/archive files, all old Go/model/runtime/module
  inputs, contract, argument, schema and checker sources verified unchanged.
- Controller root checks archived byte-for-byte. Eight required stdout/stderr
  `.log` evidence files needed targeted force-add due the generic ignore rule;
  complete staged manifest verification confirms they are included.

## Actually executed checks and attribution

1. `python3 docs/foundational/phase3c5/tools/verify.py > docs/foundational/phase3c5/evidence/qualification/verification.json` — exit 0 / PASS.
   Rechecks 695 baseline Go/module inputs and complete Git roster; five embed
   targets; 15 added Go files; 760 accepted source/archive files; 54 producer read
   sources; ten compiled local inputs across 109 saved dependency packages;
   exact final export identities; predecessor raw/compressed identity; manifests,
   exact hostile inputs, all child outputs and 59 nine-source buffer reports.
   Historical qualification-1 resolves three changed C5 identities to pre-fix
   snapshots and their original Git blobs, rather than repaired current files.
   Current checks establish source/archive consistency, not retroactive execution.
2. Actual historical 78-package roster vs fresh controller 80-package roster:
   only residualantichains and exporter added, no removals. No new whole-root
   test run. Controller build/module-verify/tidy-diff/package-list ran at reviewed
   Go head and have exit 0; they were archived and verified, not rerun.
3. Source-only qualification-2 is accepted evidence: 13 tests, nine cache
   boundaries, two valid acceptances and 56 JSON/wire rejections. No unchanged
   checker or Go suite was rerun by Task 3. Qualification-1 and successful
   pre-fix cache substitution remain historical and immutable.
4. `python3 docs/foundational/phase3c5/evidence/manuscript/build_pdf.py` — final
   baseline/current cached Tectonic, pdfinfo and pdftotext all exit 0. The final
   source-only theorem reference correction uses the TeX label (rendered 4.5)
   corresponding to Markdown 4.8, preserving the inherited numbering difference.
5. Bundled Python ran `evidence/manuscript/render_pdf.py`; all 54 final pages
   rendered at 100 dpi with successful commands. Viewed 14 full-resolution 2x2
   contact sheets and re-viewed final page 51 after the theorem-reference edit.
6. `git diff --cached --check` — exit 0. Reviewed status/stat, authoring/script
   diffs and the complete staged surface. Every artifact-manifest file is
   tracked/staged with exact expected bytes. Complete binary diff retained at
   `/tmp/alan-phase3c5/task3-complete-staged.diff` (1,385,382 bytes; SHA-256
   `520c5a558dba666b2e9b3d0944f2df27aac141b58768f7a1d39ecec95c97e358`).
   Staged 53 files, all inside the Task 3 scope; no unrelated files.

## PDF measurements and visual coverage

Fresh baseline 51 pages, current 54. Both final logs have the same seven
warnings: four underfull, one overfull and two math-font warnings. All changed
and continuation pages 1–54 inspected; no new clipping, overlap or illegibility.
The expanded Abstract/Contents introduces a mostly blank Contents continuation
on page 3. This is recorded rather than concealed; no unrelated layout rewrite.
Inherited Markdown-only Section 7 and numbering differences remain disclosed.
Eleven instruction forms and open general congruence remain intact.

New section: PDF pages **49–51**. Full-page images for controller spot-check:

- `/tmp/alan-phase3c5/task3-pdf/renders/page-49.png`
- `/tmp/alan-phase3c5/task3-pdf/renders/page-50.png`
- `/tmp/alan-phase3c5/task3-pdf/renders/page-51.png`

Conclusion continuations: pages 51–53, with images at the same folder using
`page-51.png`, `page-52.png`, `page-53.png`. Bibliography page 54 contains Ware/Malik.
Abstract/Contents pages 1–3: `page-01.png`, `page-02.png`, `page-03.png`.
Final PDF: `/Users/jacob/.codex/worktrees/phase3b4-research/alan-research/paper/alan-dissertation.pdf`.

## Failures preserved and remaining concerns

- First current PDF build failed at inline `T_F(H)` lacking math delimiters.
  Exact command, failed source and log retained under manuscript build-attempt-1.
  Corrected only the new inline math; final build succeeded.
- First verifier run failed its report-link check because its first output file
  did not yet exist. Recorded in qualification/verifier-development.json.
  Redirecting stdout creates the report before checking; the verifier itself was
  not weakened and exits nonzero on failure.
- Seven inherited PDF warnings and representation differences remain. General
  Phase 3C, full-agent substitution, dynamic authority, distributed ownership,
  split effects, exact costs, refusals and divergence remain open.
- Python installation/stdlib and reused predecessor utilities remain the stated
  trust boundary. No retroactive guarantee for historical standalone entrypoints.
- Independent Task 3/final review pending; no claim of approval, delivery or merge.

Commit: `fbcedaee554582ec9b2ca5b6ea8f873458516894` (`Reconcile residual antichain qualification and manuscript`). Post-commit working tree clean.

# Repair round 1 — R1 capture preservation

Initial Task 3 review returned **REJECT** for R1/P2. The old manuscript helpers
could overwrite accepted evidence and the delivered PDF. The finding is
substantiated by their write paths; no claim is made that a model capture was
altered. R2 (sparse Contents) and the inherited warnings remain deferred.

## Narrow correction and historical identities

Original build/render helper bytes, verifier and qualification text are retained
unchanged as `.txt` snapshots in
`docs/foundational/phase3c5/evidence/manuscript/fix-1/historical/`, with exact
path/hash mappings in `historical-sources.json` and verification against original
commit `fbcedaee554582ec9b2ca5b6ea8f873458516894`. Original command records,
archives, manifests, visual-review statements and verification output are
unchanged. The old helper entrypoints now immediately exit 1 with REFUSE.

New active `tools/build_pdf.py` and `tools/render_pdf.py` require distinct fresh
output paths and allocate them exclusively before any capture writes. Existing
directories, files and dangling symlinks refuse. Build output and rendering stay
in those new run directories; neither helper installs a PDF into paper or
publishes an accepted evidence set. Such installation is a separate explicit
authoring action. Rendering does not manufacture a visual-review statement.
The verifier resolves historical hashes to explicit snapshots only for the
original manifests and separately checks the corrected current files.

Checked current source SHA-256 identities:

- `tools/build_pdf.py`: `b643200f21defe26f73ae85bce85043eba4bfbc29a4d39037e463a427bf18598`
- `tools/render_pdf.py`: `547772204d7d25ff3a469a9672604d7d517af80a5b693e764f9178d877540b71`
- `tools/check_pdf_capture.py`: `d84675202c427b37e19ccf414b5ab1b2e9f787707bdf6c41b2d6f9333c104b8f`
- `tools/verify.py`: `6ae8ba6925dad3db869c61bb5e00ace2af79cf9dd5910531daca5bad2f22d582`
- Refusing historical build entrypoint: `1be59fcf1786ba08fcff107902f637ef9bb2f320a63492620431608c846e160d`
- Refusing historical render entrypoint: `6b7d480e10cc56711ca1848b23182dca6cc3a0aa02c8bf306ebc7fb8bc97c06f`

## Actually executed covering checks

`python3 docs/foundational/phase3c5/tools/check_pdf_capture.py /tmp/alan-phase3c5/pdf-fix-1-checks`
returned exit 0 and:

```json
{"build_files":17,"existing_build_render_unchanged":true,"fresh_successes":2,"historical_commit":"fbcedaee554582ec9b2ca5b6ea8f873458516894","protected_files":779,"protected_unchanged":true,"refusals":8,"render_files":67,"status":"PASS","visual_review":"Not performed; capture behavior only"}
```

The check performed actual fresh baseline/current cached Tectonic builds,
pdfinfo/pdftotext and 54 page renders. Then both active helpers refused their
existing directories, both historical entrypoints refused, and both active
helpers refused existing-file and dangling-symlink output paths. All eight
refusals exit 1 with REFUSE diagnostics. Existing build/render file sets and
hashes were unchanged; all 779 protected committed manuscript/evidence files
matched both their before/after hashes and original Git blobs. No delivered PDF
was installed. New temporary PDFs have no claimed visual acceptance.

Exact inner argv/stdout/stderr/exits are in `fix-1/commands.json`; fresh build and
render records, source inputs, manifests, logs, text, PDFs and output hash rosters
are archived separately in `fix-1/`. Actual fresh run folders remain under
`/tmp/alan-phase3c5/pdf-fix-1-checks/`. They do not replace the original run folder.

`python3 docs/foundational/phase3c5/tools/verify.py > docs/foundational/phase3c5/evidence/manuscript/fix-1/verification.json`
returned exit 0 / PASS, including the unchanged earlier semantic/source checks,
four historical Task 3 snapshots, 779 protected files, two fresh successes and
eight refusals. This is identity/archive validation, not rerun model semantics.

No Go/model test suite, root check, PDF marker or manuscript layout edit ran.
The sole new qualification is capture behavior and preservation. The delivered
PDF/prose and all original capture records remain byte-identical. No unexpected
failure occurred during repair qualification. Self-review covers source diffs,
archive identities, staged status/stat/full diff and no unrelated paths.
Independent scoped re-review remains pending with the controller.

Repair self-review: `git diff --cached --check` passed. All 33 staged paths are
within the R1 helper/qualification scope. Every repair manifest file matches the
staged bytes. Complete staged diff is
`/tmp/alan-phase3c5/task3-fix1-complete-staged.diff`, 1,805,583 bytes,
SHA-256 `90a8fc753f4aa3e54ed59a3d8c430a8f9e82152bcdc924d688dd88bf3ddc8c2b`.

Repair commit: `a49ee5b046b066deee8b766a25787aa606661179` (`Preserve manuscript captures during PDF reproduction`). Post-commit working tree clean. Ready for scoped independent re-review.
