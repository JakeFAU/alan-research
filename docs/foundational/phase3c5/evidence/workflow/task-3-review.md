# Task 3 review — 208e1411..fbcedaee

## Spec compliance — REJECT

- The mathematics, evidence attribution, measured results, changed manuscript alignment and bibliography satisfy the inspected requirements. The manuscript capture/reproduction helpers do not satisfy the fresh-output/no-overwrite contract; finding R1 requires correction.
- This is a Task 3 gate, not approval of the whole branch or of unrestricted semantic claims. The full Go/Python model suites and formal argument were not rerun or re-proved.

## Strengths

- `paper/alan-dissertation.tex:866–892` retains finite pointed marked automata, independent private steps, common synchronized alphabet, conjunctive marks and one unchanged K. It distinguishes local finite-suffix necessity from global necessity and preserves the MIN, MAX and generic J counterexamples. The product/minima argument and mutual-coverage reasoning are coherent under the declared finite-family hypotheses.
- `docs/foundational/phase3c5/04-results.md:35–58` retains legal-plant relevance, all blocked fixture contexts and the local-hit/global-nonblocking distinction. The control that marks all plant states is explicitly separate from paid completion.
- `docs/foundational/phase3c5/05-qualification.md:60–74` and `tools/verify.py:109–133` distinguish accepted source-only qualification from historical qualification-1 and resolve its changed inputs to pre-fix snapshots. No retroactive execution claim is made.
- `docs/foundational/phase3c5/tools/verify.py:56–108` checks baseline and accepted source identities, actual package rosters, producer archive identities and compiler paths relative to the recorded module root. `:139–163` checks changed manuscript alignment, links, archived PDF inputs/results and operation marker. Its stated identity/archive scope does not pretend to rerun the independent semantic oracle.
- `paper/references.bib:305–314` matches the authors, title, year, EPTCS volume/pages and DOI on the inspected [primary arXiv record](https://arxiv.org/abs/1108.0230). The manuscript uses the work only as background.

## Issues

### Important — R1 / P2: reproduction overwrites frozen manuscript evidence and the delivered PDF

- `docs/foundational/phase3c5/evidence/manuscript/render_pdf.py:5–12` sets its output directory to the committed evidence directory and immediately overwrites `build-commands.json` and the archived baseline/current logs, extracted text and pdfinfo. Lines 16, 21 and 26 also overwrite the baseline PDF archive, comparison and render-command record. These writes occur before rendering has successfully completed and without any fresh-directory guard.
- `docs/foundational/phase3c5/evidence/manuscript/build_pdf.py:9,20,27–34` reuses a fixed scratch directory with `exist_ok=True`, replaces earlier attempt records/sources and unconditionally copies a newly timestamped PDF over `paper/alan-dissertation.pdf`. A rerun therefore changes the delivered artifact and subsequently replaces the evidence describing the original reviewed build; the retained manifest and visual-review statement then describe different bytes until manually repaired.
- This conflicts with `docs/foundational/phase3c5/00-contract.md` under “Capture, independent challenge and qualification”: capture helpers must live in `tools/`, require a fresh output directory and refuse overwrite. It also undermines Task 3's requirement to retain failures and archive attributable commands/exits/hashes. This is a concrete reproduction behavior, not evidence that the accepted model captures were altered.
- Minimal correction: preserve these exact historical script bytes and original command records as historical execution evidence; add safe active helpers under `tools/` with explicit fresh output paths and refusal before any write when outputs already exist. Build/render only into the new run directory. Installing a rebuilt PDF into `paper/` or publishing a new accepted evidence set must be an explicit, separate authoring action. Ensure old executable locations cannot silently perform the unsafe behavior (for example, retain their original bytes under an explicitly historical inert snapshot name and leave a refusing entrypoint). Update links/manifests and qualify only fresh-output success, existing-output refusal, and unchanged committed evidence/PDF hashes; no model or unchanged suite rerun is needed.

### Minor — R2 / P3: sparse Contents continuation

- `paper/alan-dissertation.tex:65` (expanded Abstract) and `:864` (new subsection entry) push the last two Contents entries onto PDF page 3. Direct inspection of pages 2–3 confirms page 3 contains only the Phase 3C.5 entry and Conclusion, leaving most of the page blank. This is a presentation regression, accurately disclosed in the manuscript manifest, with no clipping or loss of content. A small scoped Contents-spacing adjustment is optional; this issue alone does not block approval.

### Minor — inherited warning debt

- `docs/foundational/phase3c5/evidence/manuscript/comparison.json:7–25` records seven warnings in each final build: four underfull, one overfull and two math-font warnings. Independently decompressed final logs confirm seven each. They are inherited, not newly introduced Task 3 warnings; no unrelated normalization is requested.

## Focused checks and review limits

- Read the supplied complete review package in passes: core prose/code directly and all 23 added JSON artifacts mechanically. Verified 766 current repository-resolving manifest identity entries. An initial reviewer check incorrectly compared baseline PDF-source records with current files; corrected the check to separate historical baseline identities from current ones before reporting the passing current-identity result. No repository Git command was run.
- Compared the new section, Abstract and Conclusion using the reviewed normalization logic: all three align. The normalization explicitly maps TeX's congruence label to Markdown Theorem 4.8; PDF page 51 correctly renders the inherited TeX numbering as Theorem 4.5.
- For the named accepted-versus-historical execution risk, independently checked nine current checker source hashes, all 59 saved child source-buffer reports and both saved output streams, plus the three historical pre-fix source snapshots and unchanged historical sources. No child command or test was rerun.
- For the named PDF risk, read the PDF skill and inspected existing rendered pages 1–3 and 49–54, including the full new mathematical section, Conclusion continuation and bibliography. No clipping/overlap or illegible new mathematics was observed. Reviewed the final build/render records and archived final logs; did not rebuild or rerender.
- For the named preservation risk, traced every write in both manuscript scripts statically. Did not execute either unsafe helper or mutate any evidence. The only review-created file is this report.
- Historical Git-blob checks and the complete source/archive verifier's recorded successful run remain attributed to the implementer/controller; this review inspected their implementation and available artifact identities without rerunning Git-backed verification. No claim of independently repeated whole-root qualification is made.

## Quality — REJECT / Needs fixes

R1 is the sole blocking finding. The research reconciliation is precise and the inspected evidence supports its bounded claims, but the reproduction path must preserve the original accepted artifacts and execution history before Task 3 can receive APPROVE.
