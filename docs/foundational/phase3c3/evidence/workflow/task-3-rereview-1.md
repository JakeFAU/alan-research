# Task 3 fix round 1 re-review

## Finding verdicts

- **P2 — Narrow the snapshot necessity claim. ADDRESSED.** The matched-cut witness in `paper/alan-dissertation.md:1160` and `paper/alan-dissertation.tex:811` now concludes, “The snapshot therefore omits a completion-relevant distinction in future ordering.” The two cases agree on the snapshot and yield opposite completion possibilities. The revised conclusion is limited to that distinction, leaving open which representations preserve it. PDF page 46 contains the same sentence.

## New breakage in the fix diff

None. The 13-file fix package from `32ffd396b3aa7767dc5908714a36e18cb9266c9d` to `59c05d5dc9653776a7ad58723453c4e3d0b65e14` changes the sentence in both sources, the PDF, the manuscript manifest and verifier, and adds the fix-round manuscript captures. The frozen contract, model, argument, tests and accepted experiment captures are outside the diff. The original REJECT remains in `task-3-review.md`.

## Focused checks

- Read the verifier changes at `docs/foundational/phase3c3/evidence/manuscript/verify.py:120` and `docs/foundational/phase3c3/evidence/manuscript/verify.py:141`. The verifier checks the narrower aligned source section and compares baseline evidence against the fresh current extraction and log. The retained `fix-round1/verification.json` records exit 0 and alignment PASS. The 19 command records identify one current cached Tectonic build, PDF info and text extraction, and 16 page renders, with exit 0 for each.
- Independently checked all four current paper hashes, the staged build inputs and scratch PDF, all 29 current manuscript archive/supporting-file records, and all 16 current render hashes. The prior archive/supporting-file hashes match the still-retained bytes except for the intentionally updated verifier. A fresh read-only `pdftotext -layout paper/alan-dissertation.pdf -` matched the newly archived extraction byte-for-byte. Comparing the reviewed and current extracts and render hashes identifies page 46 alone. The seven warning lines still match the baseline log.
- Inspected the existing hash-matched page 46 render. The replacement reads cleanly at the end of the continued counterexample paragraph, with the equation, headings, bottom continuation and page number intact. The other 15 inspected-page render hashes match the reviewed version. No build, suite, semantic audit, artifact verifier or PDF operation marker was repeated during re-review.

## Out-of-scope observations

None.

## Verdict

**APPROVE — Task 3 fix round 1.** P2 is addressed. No new Critical or Important finding remains in the scoped diff.
