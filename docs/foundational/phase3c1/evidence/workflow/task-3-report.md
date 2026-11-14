# Task 3 report: result, qualification and manuscript reconciliation

Status: implementation complete, ready for independent task review. No
whole-branch approval, push, PR creation, merge or release is claimed here.

Task base: `7b804df8b0dbbd160a838cdee8681558e69dc200`.
Commit: `dce9de0a7b5795dd646dac58526e70729bc6231b` (Reconcile contract preservation results and manuscript).

## Question, claim and semantic challenge

The question is whether the fixed A/B/C chain/cycle can retain an existential
completion contract through compatible grouping without erasing debt, application
labels or named service dependencies. I read the task brief/context, contract,
argument, frozen schema, independent audit and both task review reports. I checked
the mathematical steps against `research/contractcomposition/concrete.go`,
`compose.go`, `model.go` and unchanged `research/twoagent/model.go`, including
request/result occupied-channel guards, duplicate deliveries, durable ACK release,
caller crash/recovery, actual actor filtering and the named serial apply guard.

No source defect or incomplete case requiring a model/argument edit was found.
In particular:

- U->W requires only caller recovery and absent commitment.
- W->P can recover needed endpoints and admit from a queued request or durable
  commitment. No preparatory step changes the callee's outgoing progress, so
  the same abstract guard still holds. Immediate admission/upness are separate.
- P->D uses durable recording directly, or volatile/channel/retained custody;
  release requires an ACK after durable recording. It never applies again.
- Every concrete class projects to one named progress step or a stutter; repeated
  lifting starts from the actual representative reached by its predecessor.
  Complete is the same conjunction over every declared original.
- Grouping only preserves this because all debt, full bindings, policies/ports
  and physical owners survive. Canonical union equality alone is insufficient.
- Acyclicity is sufficient for every reachable state in a fixed profile, not
  necessary for one state's completion. The all-serial initial state can
  complete, but W/W/W cannot.

The report/manuscript separate conditional source argument, complete small
abstract checks, factored concrete probes, independent selected histories and
unestablished runtime/distributed claims. R0 is genuine information loss under
satisfied premises. R1 is specifically a readiness-summary failure and is not
promoted to a counterexample to a weaker trace observation. Permanent loss,
split application/retention, erased custody, changed bindings and shared-resource
depletion violate the claimed domain. No shortest-prefix or global topology
minimality claim was added.

## Changes and ownership

61 staged paths were inspected and committed within the assigned scope:

- New `04-results.md`, `05-qualification.md`, `06-reviews.md`.
- New `evidence/qualification/`: actual controller baseline/root command records,
  source/package accounting, deterministic gzip logs and archive procedure.
- New `evidence/manuscript/`: exact source/PDF/render hashes, baseline/current
  build logs/text, source-check notes and a read-only alignment/integrity verifier.
- `docs/foundational/phase3/README.md` current status.
- `paper/alan-dissertation.md`, `.tex`, `.pdf` and `references.bib`.

No edit to `01-argument.md`, schema, producer/auditor source, old executable or
historical evidence was needed. No dependency, SQLite, `.envrc`, production or
runtime change/access occurred. No subagent or self-dispatched reviewer was used.
The controller owns final workflow/reviewer archive, delivery metadata and PR.

## Fresh Task 3 checks actually executed

1. Ran the PDF skill marker successfully exactly once:
   `node .../pdf/container_tools/mark_artifact_operation_started.mjs
   --operation-kind edit --expected-output-count 1 --output-format pdf`.
2. Ran `/tmp/alan-phase3c1/task3/archive.py` to copy/hash controller evidence and
   verify current inputs. All 598 original frozen hashes match. All 637 current
   Go/module hashes match; the extra files are 10 new Go source/tests and 29
   preexisting nested-module files independently equal to base Git blobs.
   All 17 Task 1 source hashes, 22 Task 2 source hashes, 84 retained-file hashes
   and 36 raw/packed gzip artifact digests match. Gzip metadata is deterministic.
3. Parsed baseline JSONL and clean package-list output: 70 baseline terminal
   events, 57 pass/13 no-test skips, no failures; exactly 72 current packages =
   old 70 + model/export. Parsed final Task 1 normal/race logs and exit metadata:
   both new packages pass both modes. This is verification of retained executions,
   not a new suite run. No suite was rerun by Task 3.
4. Compared actual root global Staticcheck output bytes to Phase 3B.7: exact
   equality, 13 ordinary/31 integration findings, both exit1. Inspected root
   build/module/scoped lint job/result/output: clean writable-cache build exit0,
   module verify/tidy-diff/scoped Staticcheck exit0. Initial environmental cache
   warnings remain archived separately. These commands were controller-run on
   unchanged Go/module sources, not rerun here.
5. Opened both primary Abadi/Lamport sources with web tools; read the introduction
   and Section5.2 theorem hypotheses. Published Microsoft metadata is TOPLAS
   17(3), May1995, 507-534. Author revision front matter is3Nov1995 and footer
   says507-533. The discrepancy is recorded rather than silently reconciled.
   No DOI was verified/added; citation is background only.
6. Built baseline copies from Task3 base and current TeX/BibTeX in fresh scratch
   directories using `tectonic --only-cached --keep-logs /tmp/.../alan-dissertation.tex`.
   Both exit0. Ran `pdftotext -layout` on both, `pdfinfo` on current, and
   `pdftoppm -r 85 -png` on current. Retained source, log/text, PDF and render
   digests. Copied the inspected current PDF to the repository final path.
7. Opened/visually inspected all24 selected current PNG pages:1-6,11-12,14-15,
   22-26,32-35,42-46. This includes all16 same-index changed text pages and their
   continuations, with new prose, formula, bibliography and renumbered citations.
   Current46pages, baseline45. No new clipping, overlap or missing glyphs;
   inherited page32 minor margin overhang remains readable/unclipped.
8. Ran `python3 docs/foundational/phase3c1/evidence/manuscript/verify.py`, exit0.
   Final retained result checks38 qualification archive files,8 compressed
   manuscript artifacts,current paper hashes,exact new-section text equality
   after markup normalization,changed Abstract/Conclusion prose equality,and31
   local documentation links. It performs no build/test or write itself.
9. Ran `git diff --check` and `git diff --cached --check`, exit0. Inspected git
   status/stat, the full changed prose and script diffs, all generated manifest
   records by structural/hash checks, and PDF visuals. A final scoped-index
   probe checked every staged byte against disk, every path against ownership,
   all staged JSON parsing/gzip validity and deterministic metadata, and every
   archived command's argv/result/output. Exactly61paths, no unrelated file.

## Failures / deviations retained honestly

The initial archival helper parsed Go's `ok  \t` lines using a strict `ok\t`
prefix and failed its scratch assertion. I inspected the actual logs, corrected
the parser to whitespace-separated tokens and reran it successfully. The test
logs/results were never edited and no model/test failure was hidden.

An optional contact-sheet helper failed because system Python has no Pillow.
No package was installed. Poppler had already rendered every page; I used
`view_image` directly on the24 individual PNGs. The PDF manifest notes this
scratch-only failure. No document edit followed the visual inspection.

## Manuscript consistency and inherited limits

Both sources carry the same new semantic content, completion formula/quantifiers,
R0/R1, grammar and conditional progress argument. Abstract and Conclusion reflect
the new bounded result. The prior closure subsection is accurately historical.
Acyclicity quantifiers, all original debt and visible applications, abstract versus
immediate readiness and evidence categories are explicit. The broader Phase3C
remains open.

Actual final-engine warning inventory is identical in baseline/current:
4 underfull hbox,1 overfull hbox(5.2184pt),2 invalid-small-in-math font warnings.
This corrects the context shorthand “five overfull warnings”; it does not conceal
or fix historical warnings. Existing Markdown-only Section7 and section/theorem
numbering differences remain. Changed content is aligned; global equivalence is
not claimed. PDF SHA256:
`15526232e2794655608173177662a26f774fcff30d6972ab45dd6ba3c6ae7029`.

## Self-review and remaining obligations

Scope and evidence attribution are cohesive. No theorem is inferred from counts
or review approval. No runtime refinement, general congruence, fairness/divergence
transfer or arbitrary-value equivalence is claimed. Factored counts remain
nonunique contexts, not concrete-product reachability. Prior failures and hostile
inputs are unchanged. Global lint remains inherited13/31, not clean.

Task1 minor M1 remains open for controller triage; Task3 does not change tested Go
source to address a nonblocking diagnostic request. Task3 independent task review
and final whole-branch review remain pending; `06-reviews.md` states that truthfully.
The controller must archive actual review outcomes, preserve any REJECT/correction,
finalize workflow/delivery records and create/verify the draft PR without merging.
