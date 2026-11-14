# Task 3 independent review

Reviewed base `a24f8ff1c6aac22a0cffcffc2c614ccfca819c71` through head
`fcb2a6d53b5a0240dda3d530868d803bea77b8b0`.

**Spec compliance: ✅. Quality verdict: APPROVE.**

No blocking or nonblocking actionable findings were identified in this Task 3
change. This approves the manuscript/qualification change within its declared
scope; it does not independently prove the model, qualify a runtime, approve
an unreviewed later change, or authorize merging.

## Scope and method

Read AGENTS.md, Task 3 brief/report and the frozen contract. Read the complete
2,432-line supplied review package in bounded chunks. One large final chunk was
truncated by the output transport; reread only the affected middle context.
Inspected the unchanged `01-argument.md` for the concrete risk that the new
universal correspondence/sufficiency wording might exceed its source argument.
Inspected retained export/source identity and race records for the concrete risks
of numerical drift, invalid historical qualification reuse and the M1 timeout
correction. Read the PDF skill and performed only read-only artifact checks.

## Specification assessment

| Requirement | Result and evidence |
| --- | --- |
| Research-only scope | ✅ Diff consists of results/qualification, manuscript/PDF/bibliography, roadmap and README. No production, old model, dependency or historical-evidence edit. |
| Frozen finite domain | ✅ Fixed A/B/C chain, generation-one inherited bindings, eight masks, three layouts and three policies are preserved. New prose does not broaden to cycles or arbitrary topology. |
| Permission/E/Q separation | ✅ Permanent exact grants and sufficient E remain premises. Root memberships and original-bound held bits determine occupancy; Q release is atomic with application and actual result retention. Interruption and ACK release do not refund capacity. |
| Negative boundary | ✅ Matched holders, ready-now invalidation at BC commit and excluded AB-before-BC-commit trace retain their distinct meanings. Stable-ready is expressly sufficient, with necessity/optimality unclaimed. |
| Positive claim | ✅ Correspondence targets the new resource-aware abstraction; CanComplete is existential finite completion. Recoverable owners, genuine custody, fixed policies, absorbing D and atomic application/retention/release are visible premises. No old-q backward lifting, fairness, divergence or exact-cost equivalence is asserted. |
| Composition | ✅ Static joins preserve named obligations/bindings/root identity/membership/policies; independent dynamic histories are not unioned. Per-owner recovery is retained, with no combined owner availability assumption. |
| Evidence counts | ✅ Graph totals are configuration sums, probes include target-selection multiplicity, selected captured histories differ from generated checker fixtures, and the Python probe-summary check is explicitly structural. |
| Fresh/historical attribution | ✅ New scoped ordinary/race/Staticcheck/audit work is attributed to Tasks 1/2; root checks to the controller; unchanged historical package evidence is not called a fresh 74-package run. Supplemental files outside the old 637-input freeze are separately identified. |
| Manuscript alignment | ✅ New section is identical after the two presentation substitutions; changed Abstract/Conclusion paragraphs are exactly equal. New bibliography entry agrees with the recorded primary-source inspection. Inherited Section 7/numbering/content differences remain disclosed. |
| PDF | ✅ Tracked PDF matches manifest and retained text extraction; eight compressed artifacts match both raw and packed hashes. Actual archived warning lines match the disclosed baseline/current lists. Existing inspected render bytes match manifest. |
| Task 1 M1 | ✅ README now uses `-race -timeout 30m`, matching the retained passing race command. README is outside the 32 producer source hashes. No unnecessary race rerun. |

## Strengths

The claim hierarchy is unusually clear: source argument, complete finite graph
roster, factored operational probes, independent captured-history replay and
checker self-tests are not conflated. The manuscript makes the old-trace cost of
stable-ready explicit instead of presenting it as a behavior-preserving repair.
The qualification packet also preserves initial cache warnings, existing lint
findings, supplemental historical-input accounting and prior failures.

## Checks actually performed in this review

All commands below were read-only except writing this review report.

- `cat` of the four assigned inputs; `wc -l`, `rg -n '^diff --git'`, and bounded
  `sed -n` reads of the supplied complete diff.
- `cat docs/foundational/phase3c2/01-argument.md`: compared the new source claims
  against the named projection, lifting, stable-readiness and joining arguments.
- Python standard-library `gzip.decompress` plus SHA256 checks on all eight
  manuscript compressed artifacts: packed hash, raw hash and raw byte count
  match. All four current paper-file hashes match the manuscript manifest.
- `pdftotext -layout paper/alan-dissertation.pdf -` via Python subprocess:
  exact output bytes equal decompressed `current-text.txt.gz`.
- Extracted actual `Underfull`, `Overfull` and `LaTeX Font Warning:` lines from
  both compressed TeX logs: both match their seven-line manifest lists (four
  underfull, one overfull and two font warnings).
- Decompressed retained Task 1 export, without invoking producer/checker:
  4,514,212 raw bytes and SHA256
  `a7cf316219b30cdc83d455b95a16549c022d55aa422b642b8f5976315e0c0fd2`;
  72 graphs, 1,632 nodes, 2,412 arcs, 1,624 CanComplete, eight CannotComplete;
  72 grouping records all report Equal; ten histories contain 161 steps.
  Noncompletion profiles are exactly early/ready-now shared-1 masks 0,1,4,5.
  Probe metadata matches the new reported figures.
- SHA256 checked all 32 retained producer paths against current bytes.
- SHA256 compared 637 historical manifest inputs and 20 supplemental inputs
  with current files and `git show 69534188f796efab2e866750f6a76e2c7a3a1706:<path>`.
  All match. Also compared each supplemental input with
  `git show 1456a0000a0e14bb3dd105026189717c71ab905b:<path>`; all 20 match.
  No SQLite contents were accessed by these checks.
- Compared the README correction to the first line of retained `race.log`:
  `GOCACHE=/tmp/alan-phase3c2/gocache go test -v -race -timeout 30m ./research/capacitycomposition/...`.
- Python compared the new Markdown/TeX section after exactly the citation and
  results-link markup substitutions; equal. Changed Abstract and Conclusion
  paragraphs compare byte-for-byte equal.
- SHA256 checked all 24 existing `/tmp/alan-phase3c2/task3/renders/page-NN.png`
  files against the manuscript render manifest. Viewed existing contact sheets
  1, 5 and 6 (pages 1-4, 34-35, 43-48). The new section, Abstract, Conclusion
  and bibliography have readable text and intact layout; no clipping, overlap
  or missing glyph was observed. No PDF rebuild or new render was performed.

Two exploratory reader snippets initially used the wrong archive key/type
(`base` instead of `baseline`; `len()` on an integer arc count), and a profile
print used absent `Config.ID`. These read-only inspection scripts failed before
those specific print/check operations, were corrected, and produced the results
above. They were reviewer inspection errors, not artifact/test failures; no
repository bytes were changed to resolve them.

## Declined-to-verify and remaining limits

Did not rerun passing ordinary/race/Staticcheck/export/audit/checker/build/module
suites, the packet's full verification script, or any PDF rebuild. Did not
independently re-enumerate the million-case operational probe domain or entire
concrete product, replay captured histories again, mechanize the proof, or
review unchanged operational source anew. Those remain prior Task 1/2 review
and evidence inputs. Did not independently browse the background article;
metadata checking here is against the clearly attributed controller note and
aligned manuscript/BibTeX entry. Did not independently inspect the remaining
12 rendered pages visually; their inspection is attributed to Task 3 and their
render hashes were checked here. No source/index/branch mutations, SQLite or
.envrc access, dependency changes, subagents, push, PR or merge occurred.

The verifier's recorded booleans for visual coverage/warning equality are
attestations rather than independent visual or log analyses; this review
resolved the relevant binary opacity directly as described above. Approval
retains the explicitly disclosed nonmechanized proof, historical evidence,
fixed-chain/atomic/recoverable scope and open general Phase 3C obligations.
