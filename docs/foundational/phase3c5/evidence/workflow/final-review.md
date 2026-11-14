# Phase 3C.5 final independent review

Base: `347d93ba25a5f7a10f75144b1a44107be09a2f39`.
Head: `bbf58bf49407480a8de747f5c1b18c42b05b7074` (seven commits).
Review date: 2026-11-14.

## Strengths and plan alignment

- The source argument in `docs/foundational/phase3c5/01-argument.md:32` correctly uses finite coinitiality to preserve suffix hits. Its pairwise-intersection reduction law at line 49 retains empty intersections and handles empty families. The Cartesian weak-reach and residual-intersection argument at line 64 establishes both directions under synchronized visible steps and independent private steps. Pairing the reduced profile with the initial marked language closes the stated initial-completion gap. The induction stays within the declared finite grammar with one unchanged K.
- The local equivalence at `01-argument.md:97` has the correct coverage direction. If a target contains no source language, choosing one word outside that target from each source produces the finite separating H. Empty F, an empty-language member, and H containing epsilon are treated separately. The mutual-coverage argument is sound for the stated finite families. The text does not promote this local characterization to global contextual necessity.
- The Go model confines operational changes to `research/residualantichains/`. Directed inclusion retains complete epsilon-closed subset graphs after a mismatch, records both shortest directional witnesses, and uses alphabet-order BFS. Equality classes precede strict MIN/MAX filtering; coverage uses source inclusion in target. Product construction retains private transitions, synchronized Cartesian visible transitions and reverse marked reachability. The legal-trace control copies the plant state slice before altering marks, leaving paid-completion products intact.
- The original fixture and adverse controls remain substantive. The two proper differences have legal core-K witnesses, while each new monitor retains 52 nonblocking and 20 blocking K cases and all 324 blocking K/D cases. The MAX counterexample, generic MIN initial-language loss, generic J nonnecessity and local-hit/global-NB failure are distinct. Tests exercise the intended semantic boundaries, including finite-set oracles, empty families, duplicate language classes, tau cycles, multiple-word H and complete relation arcs.
- The Python auditor independently reconstructs the finite result using attributed predecessor descriptor/automaton utilities. It does not consume Go predicates or predecessor verdicts. Strict recursive typed equality checks every expected field, including graph arcs and witnesses; counts are additional checks. The acknowledged shared Python utilities remain a real independence boundary.
- Changed manuscript mathematics, Abstract and Conclusion align between Markdown and LaTeX. The source verifier also preserves the eleven-form grammar and open congruence obligation. The new bibliography metadata agrees with the inspected [primary arXiv record](https://arxiv.org/abs/1108.0230); that work is used as background, with no imported proof or novelty claim.

## Previously rejected findings

Both earlier REJECT findings remain preserved and are addressed at this head.

1. **Task 2 executable-source binding — addressed.** `tools/source_loading.py:29` compiles the retained bytes that supply the source hashes. Its private importlib facade and local import hook cover the inherited C4/C3/C2 loaders and declared local modules. The accepted direct-script bootstrap and source-only test runner use that bundle. Captured child manifests match the nine current source buffers. The nine poisoned-cache controls demonstrate the original stock-loader vulnerability and the intended source-only result. Qualification-1 and pre-fix snapshots remain historical; no claim establishes that the original captured run actually used stale code.
2. **Task 3 evidence/PDF overwrite — addressed.** `tools/build_pdf.py:13` and `tools/render_pdf.py:12` allocate fresh output directories before capture writes. They write beneath those directories and do not install the generated PDF or publish accepted evidence. Both retired executable locations refuse immediately. Historical script bytes are inert snapshots bound to their original Git blobs. The retained qualification records two fresh successes, eight refusals, and preservation of 779 protected files. Source inspection supports that correction; the semantic model and delivered PDF were not changed by the repair.

Paths beginning `tools/` in this report are relative to `docs/foundational/phase3c5/`.

## Issues and explicit triage

### Critical

None found.

### Important

None found. No further correction or qualification rerun is required by this review.

### Minor

1. **P3 — sparse Contents continuation.** `paper/alan-dissertation.tex:65` and `paper/alan-dissertation.tex:864` expand the Abstract/Contents enough that PDF page 3 contains only the new subsection and Conclusion entries. I directly inspected the existing page-3 render and confirmed the large blank area. This is presentation debt, with no missing or clipped content. **Disposition: defer; nonblocking.** A later scoped Contents-spacing adjustment is optional and should receive its own PDF comparison.
2. **P3 — inherited TeX warnings.** `docs/foundational/phase3c5/evidence/manuscript/comparison.json:7` records the same seven warnings in baseline/current: four underfull boxes, one overfull box and two math-font warnings. **Disposition: defer; nonblocking inherited debt.** No new warning was introduced, and unrelated manuscript normalization is unnecessary for this research slice.

The inherited Markdown-only section and theorem-numbering differences remain disclosed. In particular the new TeX reference correctly resolves to the existing PDF congruence theorem numbered 4.5, corresponding to Markdown 4.8; the packet does not claim global representation normalization.

## Checks actually performed

- Read the frozen contract, plan, source argument, schema, independent-audit/results/qualification accounts, controller decisions, workflow ledger and original/re-review findings. Read the current Go implementation and tests, active Python checker/capture/source-loader/test/PDF/verifier code, relevant reused C4/C3/C2 utilities, and manuscript/roadmap diff. Reviewed historical reconnaissance for the original fixture versus the rejected plant-disabled alternative.
- Verified the supplied review package's diff portion equals `git diff -U10 347d93ba..bbf58bf4` exactly. Verified all 872 changed worktree files equal their reviewed-head Git blobs. Mechanically parsed all 355 changed JSON files and decompressed all 434 changed gzip files successfully. This archive decoding is an integrity check, not execution of historical source or a semantic oracle.
- Ran `python3 -B docs/foundational/phase3c5/tools/verify.py`: **PASS**, exit 0. It checked 695 baseline inputs, five embedded inputs, 15 new Go inputs, 760 frozen accepted files, 54 producer source entries, ten compiled local files, recorded module-relative dependency identities, historical snapshots, accepted child records, source/manuscript alignment and artifact identities. It checked the actual 78-to-80 package delta and preserved root-command records without running those commands again.
- Verified raw export identity and capture equality: **14,097,260 bytes**, SHA-256 `ab3cbdfda577a3ae804eb76a94a8c7ac665cf00c69f992ad6804b6e0d0e27b8a`. The verifier cross-checks the 3,564 products, 45,800 states, 59,492 arcs and 2,772 substitution rows, including all 1,928 blocking rows. This is archived-result checking; the independent full reconstruction was not rerun in this review.
- Independently checked all 27 workflow-manifest entries and 31 preflight-manifest entries, respecting whether the original source was already gzip-compressed. Initial reviewer helper assumptions about review-package headers and preflight compression layers were corrected before these passing comparisons; they were not repository defects.
- Inspected the archived ordinary/race/Staticcheck, semantic RED and source-only Python command/results. The current source-bound Python qualification records 13 passing tests, two valid CLI acceptances and 56 hostile rejections. The passing verifier checked all 59 child source-buffer reports and retained output identities. These are attributed earlier executions, not fresh reviewer test runs.
- Checked all 54 original page-render identities against the manuscript manifest; directly viewed pages 3, 49, 50 and 51. The new mathematical section and continuation were readable without observed clipping or overlap. Full 54-page visual acceptance remains attributed to the earlier recorded inspection. No build, rendering, PDF installation or artifact operation marker was executed here.
- Confirmed only the packet, new research package, roadmap and manuscript files differ from the base. No production runtime, old model, dependency/module, predecessor test or predecessor evidence file changed. No Go suite, full Python suite, model campaign, runtime execution, SQLite open or physical-history replay was rerun.

## Declined to judge

- Unrestricted ALAN contextual congruence, arbitrary full-agent substitution, distributed ownership, dynamic authority, split-effect recovery, exact costs, refusals, fairness and divergence preservation: explicitly outside this finite monitor-substitution contract and still open.
- Globally weakest contextual equivalence, minimal automata, verification speedup, and unrestricted infinite-family extensions: no such result is claimed; the local finite-H characterization cannot establish them.
- Arbitrary externally supplied Go automata or general-purpose Python sandbox behavior: the Go boundary exports pinned builder-owned data, and the source bundle explicitly trusts Python/stdlib and only governs declared modules and entrypoints.
- Retroactive source-only guarantees for historical checker executions or standalone predecessor entrypoints: the original cache exposure and unchanged entrypoints are explicitly retained as historical limits.
- Fresh whole-root test/race/integration qualification and global Staticcheck cleanliness: not part of this source-scoped change; earlier results and inherited diagnostics remain attributed, and no such fresh claim is made.
- Global manuscript proofreading, inherited representation normalization and proof-checking the complete Ware/Malik paper: outside this packet; the changed text, new reference metadata and scoped visual evidence were checked.

## Assessment

**APPROVE.** The whole branch satisfies the inspected Phase 3C.5 contract with a conditional source argument and exact finite qualification, while retaining its counterexamples, blocked cases and open obligations. Both important review findings are corrected without rewriting their historical evidence; only the two explicitly deferred presentation issues remain.

**Ready for the controller's draft-PR delivery: yes.** This verdict is neither a mechanized proof nor permission to merge automatically. The only reviewer-authored file is this report; the index, HEAD, branch, sources and evidence were not modified.
