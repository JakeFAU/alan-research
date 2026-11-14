# Phase 3C.1 delivery

[Draft PR #50](https://github.com/JakeFAU/alan-research/pull/50) was opened on
14 November 2026 against `main` and attached to the Codex chat. The first
post-creation GitHub check reported open, draft, mergeable and unmerged, with
head `46f4e57bbed7cd517006b02ac001f21380178241` and base
`1456a0000a0e14bb3dd105026189717c71ab905b`. This is a recorded delivery snapshot,
not a claim that later PR state cannot change.

## Source and review boundary

- Model/export qualification: `5f907cda569d97ec5c08923b35e77de5b3df06ce`.
- Independent auditor: `7b804df8b0dbbd160a838cdee8681558e69dc200`.
- Reconciled and independently reviewed semantic/manuscript head:
  `dce9de0a7b5795dd646dac58526e70729bc6231b`.
- Exact review archive and status: `46f4e57bbed7cd517006b02ac001f21380178241`.

All three task gates returned spec/quality APPROVE. The final technical review
returned APPROVE with no blocking findings. [Review status](06-reviews.md)
links the exact reports. M1 remains explicitly deferred and unfixed: a future
factored-test failure may require instrumentation to identify its context.
No source, assertion, theorem premise or evidence was weakened to obtain approval.

Changes after the reviewed semantic head only archive the workflow, update
review/plan status and record delivery. They do not change the model, auditor,
manuscript, PDF or qualified experimental records. The final workflow archive
preserves all three controller rulings and their stated cost if wrong. Two
original extracted briefs use deterministic gzip so their trailing blank lines
survive exactly without creating text-diff whitespace errors.

## Qualification and limits

[Results](04-results.md) and [qualification](05-qualification.md) contain the
conditional source argument, 639-state/1,237-arc abstract checks, factored probes,
seven independently checked model histories and 35 hostile-input rejections.
The frozen 70-package ordinary baseline and both new-package ordinary/race
runs cover the current 72-package roster. Root build/module checks and scoped
Staticcheck pass; global Staticcheck retains the exact inherited 13/31 findings.
The manuscript/PDF, source hashes and archive bytes were checked before push.

This packet establishes no runtime refinement or unrestricted compositional
closure. Phase 3C remains open. No merge or release was performed. The named
branch and worktree are preserved for review.
