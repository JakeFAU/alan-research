# Task 3 report — Phase 3C.2 qualification and manuscript reconciliation

Status: DONE.
Commit: `fcb2a6d53b5a0240dda3d530868d803bea77b8b0` — `Qualify shared-capacity results and reconcile manuscript`.
Base: `a24f8ff1c6aac22a0cffcffc2c614ccfca819c71`.
Branch: `codex/phase3c2-shared-capacity`.
Workdir: `/Users/jacob/.codex/worktrees/phase3b4-research/alan-research`.

## Question, result and exact scope

Reconcile the reviewed fixed-chain capacity result with qualification and the
manuscript without expanding it into a runtime or general Phase 3C result.
Added `04-results.md`, `05-qualification.md`, source-bound qualification and
manuscript archives. Updated Phase 3 roadmap, matching manuscript sections,
Abstract and Conclusion, BibTeX background, rebuilt PDF, and README race timeout.
No production/old model/dependency/historical evidence changes; no SQLite or
`.envrc` access. No subagents, push, PR or merge. Controller-owned plan,
preflight, workflow and review-status files remain outside staging.

The manuscript preserves permission/E/Q distinctions, exact holder/root
identity, atomic apply/result-retention/Q-release, recoverable interruption and
fixed-policy premises. Stable-ready has a conditional source sufficiency
argument, while excluding a legal old AB-before-BC-commit progress trace. The
new resource-aware abstraction, not old-q, is the universal finite-lifting
target. Existence of a completion suffix is separate from fairness/inevitability.
Static joins do not union dynamic histories.

Measured counts remain exactly: 72 graphs, 1,632 states, 2,412 arcs, 1,624
CanComplete, eight CannotComplete; 72 grouping records; 311 old edge states,
1,854,036 synthesized representative cases, 11,972,664 arcs, 2,451,168 lifts,
4,833,504 lift microsteps. State/arc graph totals are sums over configurations;
probe counts include multiplicity. Ten captured histories/161 steps are distinct
from 144 Python checker fixtures/7,200 steps. Five test methods and 43 hostile
CLI rejections retain their exact independent-check versus checker-test scope.

## Fresh versus inherited verification

No already-passing model/race/export/audit suite was rerun. README is outside
the export hash set, so M1 was addressed by adding `-timeout 30m` and comparing
to the recorded passing race command. Task 1/2 source remains unchanged.

Fresh Task 3 identity checks:

```sh
python3 docs/foundational/phase3c2/evidence/qualification/capture_inputs.py
```

Exit 0. All 637 Phase 3C.1 manifest entries match both current files and merged
base `69534188f796efab2e866750f6a76e2c7a3a1706`. The broader tracked Go/module
roster identified 20 additional preexisting files outside that freeze: ten
root-package Go inputs and two testdata files under examples/evidence, three
archived Phase 1K public-consumer nested-module inputs, and five archived
Phase 3B.4 apparatus/fixture files. All 20 are present and identical in both
merged base and original historical baseline
`1456a0000a0e14bb3dd105026189717c71ab905b`. Their names/hashes and separate
historical-baseline comparison are in input-identity.json. They are expressly
not claimed part of the original 637-file freeze. Twelve new Go inputs are
all under research/capacitycomposition. All 32 producer hashes and 201 Task 2
artifact hashes match current bytes. No unexplained drift exists.

Historical ordinary qualification is the original fixed 70-package run (57 pass,
13 no-test skips) plus Phase 3C.1's two model/exporter packages. The clean current
roster has 74 packages = 72 inherited + two capacitycomposition packages, no
removals. Current combined ordinary coverage uses unchanged historical inputs
and fresh Task 1 scoped checks. No fresh single 72/74-package invocation,
whole-repository race/integration suite or runtime campaign is claimed.

Controller root checks were already completed and are copied verbatim, with
runner, job/environment/source-head/output/exit records:

| Command | Exit / output | Archive |
| --- | --- | --- |
| `go build ./...`, initial | 0; 186-byte module metadata cache permission warning | qualification/build |
| `go build ./...`, approved writable-cache rerun | 0; empty output | qualification/build-cache-writable |
| `go mod verify` | 0; `all modules verified` | qualification/module-verify |
| `go mod tidy -diff` | 0; empty output | qualification/module-tidy-diff |
| `go list ./...`, initial | 0; warning plus package lines | qualification/package-list |
| `go list ./...`, approved writable-cache rerun | 0; clean 74-package list | qualification/package-list-cache-writable |

These were controller executions in this task, not executions by this worker.
Task 1 scoped ordinary/race/Staticcheck and Task 2 audits/hostile controls are
attributed to their retained reports and logs. Global lint is not claimed
clean: inherited Phase 3C.1 evidence has 13 ordinary and 31 integration-tag
findings; no fresh broad lint rerun occurred. No go vet invocation occurred.

## Manuscript and primary background

Added identical Phase 3C.2 mathematical prose after markup normalization to
Markdown 6.17 and rendered LaTeX 6.12. Updated Abstract/Conclusion paragraphs
are byte-equal in both source forms. Previous disjoint-resource Phase 3C.1
claims remain unchanged. General closure, arbitrary contexts, fan-out,
cancellation, changing grants, split effects, permanent owner loss,
distributed allocation/recovery, cost equivalence and runtime refinement remain
open. New results/qualification documentation links to controller review gates
without claiming approval or delivery.

The new Coffman/Elphick/Shoshani background citation uses the original article
metadata previously inspected by the controller (Computing Surveys 3(2), June
1971, 67–78). The copied primary-literature note attributes that inspection and
retains the publisher DOI page's 403 outcome. No new independent web verification
by this worker is claimed; no DOI was inserted. The phenomenon is described as
established resource-wait deadlock, not ALAN novelty.

## PDF operation, build and visual inspection

Read the required PDF skill. Invoked the marker exactly once, immediately before
first manuscript source authoring:

```sh
/Users/jacob/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/bin/node /Users/jacob/.codex/plugins/cache/openai-primary-runtime/pdf/26.1112.11011/skills/pdf/container_tools/mark_artifact_operation_started.mjs --operation-kind edit --expected-output-count 1 --output-format pdf
```

Exit 0, no output. Exact command and placement are retained in
`evidence/manuscript/operation-marker.json`. Final tracked artifact is
`paper/alan-dissertation.pdf`.

Built baseline sources extracted with `git show 69534188...:paper/<file>` and
current sources under `/tmp/alan-phase3c2/task3/{base,current}/`:

```sh
tectonic --only-cached --keep-logs alan-dissertation.tex
pdfinfo /tmp/alan-phase3c2/task3/{base,current}/alan-dissertation.pdf
pdftotext -layout /tmp/alan-phase3c2/task3/{base,current}/alan-dissertation.pdf /tmp/alan-phase3c2/task3/{base,current}/text.txt
```

The brace paths above summarize the separate exact argv arrays retained in
build-commands.json. Both builds/extraction commands exit 0. Initial and second
build command records are retained separately; the final current rebuild
clarifies that sufficiency is established while necessity/optimality are not
claimed. An earlier wording edit explicitly placed Q release atomically with
BC application before consumption. No semantic source change occurred.

Current PDF: 48 pages; freshly rebuilt baseline: 46. Exact warning lines agree:
four underfull hboxes, one 5.2184pt overfull hbox, two math-mode font warnings.
No new warnings. Same-index text changes occur at pages 1, 2, 4, 5, 11, 14, 22,
24, 32, 34, 44-48 (including bibliography renumbering).

Rendered with existing Poppler at 100 dpi, using the bundled Python/Pillow for
contact sheets; no package/dependency installation:

```sh
/Users/jacob/.cache/codex-runtimes/codex-primary-runtime/dependencies/python/bin/python3 /tmp/alan-phase3c2/task3/render_pdf.py
```

All 24 selected pdftoppm commands exited 0; exact argv arrays and render hashes
are retained. Inspected pages 1-6, 11-12, 14-15, 22-25, 32-35 and 43-48. New
section (44-45), Abstract (1), Conclusion (45-46), references (47-48), changed
citation pages and continuations show no new clipping, overlap or missing
glyphs; heading and margin layout is intact. Final affected sheets were inspected
after wording fixes. Inherited Markdown-only Section 7, section/theorem number
differences and format-specific markup remain explicitly disclosed. This is not
a global manuscript-equivalence claim.

## Artifact verification and self-review

```sh
python3 docs/foundational/phase3c2/evidence/manuscript/verify.py
```

Exit 0. Retained verification.json reports PASS: 637 historical manifest inputs,
20 supplemental inputs, 12 new Go inputs, 32 producer inputs, 201 Task 2
artifacts, 23 qualification files, eight compressed manuscript artifacts, 29
local links, normalized new-section equality, changed Abstract/Conclusion
equality, current paper hashes, unchanged warnings, changed-page inspection
coverage, and README race-timeout match. Raw/packed hashes bind the archived
build logs/text; manuscript source, PDF and rendered-page hashes are explicit.
No temporary PDF renders, caches, SQLite files or credentials are committed.

`gofmt -l research/capacitycomposition` and `git diff --check` exit 0 with no
output. Exact tool versions and command outputs are retained in
qualification/task3-commands.json: Go 1.27.4 darwin/arm64, Python 3.9.6,
Tectonic 0.17.0, Poppler 26.11.0, Staticcheck 2026.2.2 (0.8.2).

Self-review covers the entire task-owned diff, docs, manifests/command records,
new manuscript prose and bibliography, alignment, source identity and staged
scope. No old proof/qualification bytes or production source are included.
The inherited 01-argument task-status prose stays untouched because it is in
the frozen export source set; new result/qualification docs give current status.

Transient authoring checks did expose: (1) an initial overly broad assertion
that all nonmanifest Go inputs were new, corrected by separately accounting for
20 historical files; (2) system Python lacked Pillow, so existing bundled Python
was used; (3) the first verification attempted its verification.json link before
that output existed, resolved by capturing stdout to that file. These were
archival/environment setup issues, not semantic failures or changes to reviewed
model evidence. No failed semantic evidence was removed or rewritten.

## Remaining boundaries

Task 3 independent review and final branch review are controller responsibilities.
This worker does not mark them approved. Mathematical statements remain
conditional source arguments plus finite evidence, not mechanized proofs.
Whole concrete-product enumeration, independent Python enumeration of Go probe
cases, universal runtime refinement and general Phase 3C closure remain absent.

## Commit and post-commit state

Committed 52 task-owned files; staged diff/check passed. Explicitly included six
manifest-bound controller output.log files otherwise excluded by the global log
ignore. Every staged byte matched its worktree file; every manifest path was in
the index; no controller file or Go/runtime edit was staged. The initial git add
was blocked by the filesystem sandbox; approved escalated staging and commit
completed successfully. No automatic approval rejection occurred. Post-commit
artifact verification remains PASS. Remaining status contains only controller-owned
phase3c2-plan.md, 06-reviews.md, evidence/preflight/ and evidence/workflow/.
