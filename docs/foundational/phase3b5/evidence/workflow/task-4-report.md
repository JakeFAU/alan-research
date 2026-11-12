# Task 4 worker report

## Status and exact change set

Narrative and manuscript source work is committed and ready for independent task review. Base and frozen executable commit `4ad1f9ef37c57375d4a39eafc9a5fcd91bc0425d`. Task commit `4cabef20bf48453cba9dad5ac03a7d9e54a19d77`, tree `0965a261085d1acdfb7051cd0ba2277d08d7cdff`. The verified commit paths are:

- `docs/foundational/phase3b5/06-qualification.md`
- `docs/foundational/phase3b5/07-results.md`
- `docs/foundational/phase3b5/08-reviews.md`
- `paper/alan-dissertation.md`
- `paper/alan-dissertation.tex`

The controller's unstaged plan edit remains outside the commit. No executable files, evidence scripts, references, PDF, runtime repository, or remote state were modified. The reports and records linked under `evidence/workflow/` and `evidence/task4/` will be staged by the controller.

## Claim and evidence categories

- Exhaustive finite safety and safe-suffix characterization on `S2/SeparateGrantRevocation` has 229,834 reachable states and 1,807,246 arcs. The new subsection distinguishes operational graph enumeration, independent descriptor interpretation, and independently computed reverse reachability. The supporting source-level necessity/sufficiency argument remains a sketch under the reachable-state invariants.
- The retained one-step universal-completion counterexample is immediate b revocation. Other losing cuts isolate immutable early refusal, revocation before a fresh opportunity, and a crash clearing an old revoked-path permit. Retry cycles keep existential completion separate from inevitable progress.
- The conditional forward relation checks eight supported final cases and twenty captures across both sealed modes. N06 retains its first allow and debit and refuses its only fresh gate with zero contact/application in the declared evidence. N07 retains one unresolved BA opportunity at each intermediate cut. In the final cut, it consumes the BA result and fresh c-backed AB result. A collapse of the complete roots under one revocation would block that observed AB admission. The independent ledger's c-with-b mutant and the separate concrete wrong-target attack expose distinct scope errors.
- The applied-unknown case supplies a retained exclusion after acknowledged whole-Store loss. It contributes no projection events. The evidence does not establish runtime application/receipt atomicity, reverse simulation, model crash correspondence after whole-Store loss, general permission-algebra equivalence, distributed revocation or composition beyond two agents.
- The affected Abstract and Conclusion and new subsection are synchronized. The Phase 3B.4 nonrevoking domain is retained as historical. The new body is MD 6.13 and expected PDF 6.8. The pre-existing MD-only Section 7 is byte-identical and remains outside this reconciliation. The next bounded question concerns evidence distinguishing admitted opportunity from applied-but-unrecorded result after one acknowledged process loss under a fixed observation contract.

## Direct worker checks

The PDF skill marker ran immediately before first manuscript authoring, exactly once:

```sh
/Users/jacob/.cache/codex-runtimes/codex-primary-runtime/dependencies/node/bin/node /Users/jacob/.codex/plugins/cache/openai-primary-runtime/pdf/26.1111.10934/skills/pdf/container_tools/mark_artifact_operation_started.mjs --operation-kind edit --expected-output-count 1 --output-format pdf
```

Exit 0, empty output, 0.735395291 seconds. The following manuscript authoring command also exited 0. I notified the controller. The controller previously checked Tectonic 0.17.0 and Poppler 26.10.0 and owns PDF build, render and inspection after substantive review.

- `python3 /tmp/alan-phase3b5/check-task4-narrative.py`, exit 0. Compared the normalized affected Abstract, subsection, formula and Conclusion across sources. Checked display delimiters, exact commands in the formula, absence of control bytes, preservation of the earlier nonrevoking boundary, and byte equality of MD Section 7 with the frozen source. Verified all 582 frozen file hashes and the manifest hash `51ab5216312d59b65987c239b5ff66cbbc3cbd7b44b9fb9cbcc40d80c6e42936`. Independently checked all ten active archived dossier verdicts, eight supported final cases, twenty supported captures, N06/N07 physical counts, exact N07 unresolved arrays and the excluded histories' zero projection events.
- Same check verified compressed/uncompressed byte hashes, sizes and direct exits in the nine completed controller records. Affected race exit 0 with 422 passing test events, seven tested packages and zero named skips. Build, module verify, module tidy diff, scoped Staticcheck, opt-in model export and deterministic historical/new exports exit 0. Both global Staticcheck runs exit 1 and match the inherited baseline logs exactly. The ordinary suite record was still open at this handoff.
- An initial `git diff --check` exposed malformed formula text from shell authoring escapes. I repaired the display expression in both sources. Subsequent `git diff --check`, `git diff --cached --check` and `git diff HEAD^ HEAD --check` each exited 0. `git diff --exit-code 4ad1f9ef37c57375d4a39eafc9a5fcd91bc0425d -- paper/references.bib paper/alan-dissertation.pdf` exited 0. Inspected staged and committed paths and confirmed the controller's plan edit is the sole unstaged path at handoff.
- `shasum -a 256` on the five deliverables exited 0. Manuscript hashes are Markdown `a6a8d859236a2a40b6cb60982b3e47baa214714869376480d233439e4f7e5cda` and TeX `c15984c0142bfa636ce5a877158765e832dda117a71372697d8fed28dd05b49c`. Git metadata escalation committed only the five inspected paths and exited 0.

## Remaining controller gates

Record the final full ordinary exit and complete named skip inventory from the runner. Update the pending row and related status in 06. Finish final format and manifest checks, independent Task 4 review and broad whole-branch review. Build and visually inspect the PDF and retain source/PDF digests. Archive workflow and qualification records at the linked paths, verify the clean intended delivery diff, push and open the draft research PR. No merge, release, production change or additional runtime campaign is part of this task.

## Controller artifact completion after source review

Task source review at4cabef20 approved with no findings. Built those exact manuscript inputs with `/opt/homebrew/bin/tectonic --only-cached --keep-logs alan-dissertation.tex` in `/tmp/alan-phase3b5/pdf-final`, exit0. Before/after source hashes match. PDF42 pages, SHA25683bfaa590ad6b3cbd61ecf32043efdce93ed7ac75805d4feee1e625911de5bad. Poppler rendered pages1,2,9,28,29,38–42; controller viewed all ten PNGs and observed no new layout defects. Five unique inherited TeX warning patterns concern unchanged source ranges; exact records and archive are in evidence/task4/pdf. Installed PDF matches the archived bytes. Any later source fix must rebuild the artifact. Full ordinary exit, final evidence/index update and whole-branch review are still pending.

## Qualification and source-review closure

Follow-up status-only commit `836cbb1f5857cc3daf2cbdc603a403a389a12bec`, tree `20905ae36bd79cfa1d21f6da22401537ff0c29c9`, updates only 06-qualification and 08-reviews. The previously pending ordinary command `go test -json -count=1 -timeout=60m ./...` closed with exit 0 in 1,505.540366375 seconds. It recorded 2,633 passing test/subtest events, 51 tested packages and 13 packages without tests. The three named skips were the optional `examples/communication/TestCommunicationArtifactExport` (explicit artifact export unrequested), `examples/join/TestJoinArtifactExport` (`ALAN_JOIN_EXPORT` unset), and `examples/substitution/a/TestCandidateAExport` (`ALAN_PHASE1I_A_OUTPUT` required to export first retained bytes). Affected race closed exit 0 with 422 passing test/subtest events, seven tested packages, three without tests, and zero named skips. Neither suite recorded a failure.

The controller's verifier confirms all ten command records and logs, all 582 frozen hashes, format checks, exact historical export equality, and matching new model/audit/progress bytes. A read-only `python3 -` check exited 0 after independently verifying the ten summary records against the individual manifests and compressed/uncompressed log hashes and byte counts, both suite counts, three exact skip reasons, every frozen source hash, format records, PDF build exit and manuscript source hashes, installed PDF digest, and approval at the reviewed commit. The build and inspection records report cached Tectonic 0.17.0 exit 0, 42 pages and SHA-256 `83bfaa590ad6b3cbd61ecf32043efdce93ed7ac75805d4feee1e625911de5bad`. The controller inspected pages 1, 2, 9, 28, 29 and 38–42. Five unique box-warning patterns came from unchanged source. The source review approved `4cabef20bf48453cba9dad5ac03a7d9e54a19d77` with no findings. No manuscript source or PDF rebuild was needed in this follow-up.

Scoped `git diff --check`, `git diff --cached --check`, and `git diff HEAD^ HEAD --check` exited 0. The source/bibliography `git diff --exit-code` exited 0. Git metadata escalation committed only the two verified status paths. The controller's PDF, evidence and plan changes remain outside the commit. Whole-branch review, final archive checks and draft delivery remain controller gates.
