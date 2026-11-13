# Phase 3B.6 qualification

## Frozen inputs

The affected research qualification uses commit `e0615c1300baef94c077fea8168c729e6f6434ed` and tree `a742cc32338d939ebd8f8e53483530ebc33b7e9c`. Its 612-input manifest covers executable sources, fixtures and evidence tools. The manifest SHA-256 is `c7469fb8b970668c03076f4c2a4e89fac7a0f8eed1b81bdcf34aa3c27be1e2c2`. The environment uses Go 1.27.4 darwin/arm64 and Staticcheck 2026.2.2 (0.8.2), with offline module and cache settings in the command records. Later narrative and PDF edits fall outside the executable manifest. [Freeze record](evidence/task4/freeze.json).

The unchanged ordinary group started from `77f897e0c8bde7c1b3ea736a19f26029debb588e`. Dependency selection excluded any packages that depend on the changed research packages. All 528 inputs for that group matched the final freeze, and the final `go list ./...` inventory contains exactly the 58 unchanged and 10 affected packages, with no overlap. A focused dynamic-read scan identified provenance emission that could record excluded research source hashes without affecting the commands' assertions. The qualification uses two ordinary invocations, each against its recorded input set. [Package selection](evidence/task4/ordinary-package-selection.json) and [dynamic reads](evidence/task4/ordinary-dynamic-read-review.md).

The historical runtime producer is `5b4662fd682a220ca3065b15e86672b7e56567e7`, tree `44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`, produced with Go 1.27.3. Phase 3B.6 checks retained runtime evidence from both producer modes. Task 3 verified the two owning archive byte streams, 220 selected dossier members, 17 complete-source hashes and 31 literal source extracts. Archive verification hashes SQLite members as opaque bytes. The post-qualification source and archive recheck exited 0 in 1.143 seconds. [Protocol and corpus provenance](04-correspondence.md), [source manifest](evidence/task3/protocol-source/index.json), [historical toolchain](../phase3b4/evidence/prospective-main/tools-and-readers.json), and [post-qualification recheck](evidence/task4/source-archive-verification.json).

## Commands and recorded results

The controller verified all 612 frozen inputs, the 528 earlier input hashes, 14 closed command records and the complete package inventory. Formatting of the 21 changed Go files exited 0 with empty output. [Final verification](evidence/task4/verification.json).

| Job | Direct exit | Elapsed (s) | Recorded result |
| --- | ---: | ---: | --- |
| `go test -json -count=1 -timeout=60m` with the 58 unchanged packages | 0 | 1,470.108 | 2,383 test and subtest pass events in 48 tested packages. Ten packages contain no tests. Three named skips. |
| `go test -json -count=1 -timeout=60m ./research/effectobservation/... ./research/delegationcorrespondence/...` | 0 | 31.571 | 419 pass events in seven tested packages. Three packages contain no tests. No named skips. |
| Same two affected package trees with `-race -timeout=15m` | 0 | 263.641 | 419 pass events in seven tested packages. Three packages contain no tests. No named skips. |
| `staticcheck ./research/effectobservation/... ./research/delegationcorrespondence/...` | 0 | 1.099 | Empty output. |
| `staticcheck ./...` | 1 | 3.106 | Exactly the baseline's 13 ordinary diagnostics, byte for byte. |
| `staticcheck -tags=integration ./...` | 1 | 3.408 | Exactly the baseline's 31 integration diagnostics, byte for byte. |
| `go build ./...` | 0 | 2.228 | Empty output. |
| `go mod verify` | 0 | 0.074 | All modules verified. |
| `go mod tidy -diff` | 0 | 0.074 | Empty diff. |
| New model export | 0 | 0.080 | Reproduces the retained graph. |
| New independent audit export | 0 | 0.076 | Reproduces the retained audit. |
| Frozen historical export runner | 0 | 30.081 | Executes two runs of each legacy export and the Phase 3B.5 audit and progress commands. |
| Phase 3B.5 revocation model export | 0 | 1.105 | Reproduces the retained `S2/SeparateGrantRevocation` model. |
| Prior default correspondence progress command | 0 | 2.147 | Reproduces the retained default progress output. |

Together the two ordinary groups record 2,802 pass events, 55 tested packages and 13 packages containing no tests. The only named skips in the final ordinary groups are optional artifact exports from the unchanged group. These are distinct from packages containing no tests. [Command record and log](evidence/task4/qualification-unaffected/).

| Named skip | Reason |
| --- | --- |
| `examples/communication/TestCommunicationArtifactExport` | Explicit artifact export was not requested. |
| `examples/join/TestJoinArtifactExport` | `ALAN_JOIN_EXPORT` was unset. |
| `examples/substitution/a/TestCandidateAExport` | `ALAN_PHASE1I_A_OUTPUT` was unset. |

The affected jobs enable the opt-in Phase 3B.6 and Phase 3B.5 two-mode corpora and the applicable earlier raw, promoted and progress regressions. The historical inputs differ by ordinary/race producer mode where recorded. The command records preserve exact arguments, environment, direct exits, elapsed times and log hashes. [Concrete runbook](05-runbook.md), [unchanged command records](evidence/task4/qualification-unaffected/) and [affected command records](evidence/task4/qualification-final/).

## Prior exports and scoped checks

Each of the four archived historical exports matched its retained bytes on both frozen-source runs. The default authority export preserves its sixteen configurations. The Phase 3B.5 model, audit and progress exports match their retained bytes. The default correspondence progress output matches its historical byte count and recorded digest. The new graph and audit exports match the retained Task 1 and Task 2 bytes. [Frozen output comparisons](evidence/task4/exports/exports.json).

Task 1 enumerated the graph twice with byte-identical output. Task 2 ran the independent audit and retained compiled RED/GREEN evidence and four failing semantic mutants. Task 3 reproduced ten pairs of old revocation and raw legacy reports from frozen Phase 3B.5 source and checked them against archived digests before comparing bytes to current outputs. It also exported the four new case reports twice and checked hostile inputs against the causal and renewal layers. The causal mutant fails an end-to-end assertion and the exact-renewal mutant fails direct helper assertions. [Task 1](evidence/task1/) · [Task 2](evidence/task2/) · [Task 3](evidence/task3/).

## Apparatus and intentionally failing checks

The Task 1 exporter test initially used an unavailable standard-library helper and failed to compile. A local erroring writer made the test executable. Initial Task 1 and Task 3 Staticcheck runs failed on the default cache write. The recorded successful runs use a writable `STATICCHECK_CACHE`. Task 2 preserved initial whitespace failures for the nested mutant patch files and corrected their representation before qualification. Task 3's initial RED command wrapper obscured the test exit. A direct reproduction recorded exit 1 with the intended semantic failures. Its initial frozen-source extraction used a Python tar API unavailable on this host, and the replacement checks paths before extracting regular members. The investigation log records early candidate and hostile-fixture failures separately from the retained assertions. [Task 1 verification](evidence/task1/verification.md) · [Task 2 verification](evidence/task2/verification.md) · [Task 3 investigation](evidence/task3/investigation.md).

The earlier short-mode targeted baseline is separate from the final ordinary and race groups. Its 11 named skips and direct exit are in the baseline record. Compiled RED and mutant assertion failures test whether the intended guards are exercised. [Targeted baseline](evidence/preflight/baseline.json) and [review dispositions](08-reviews.md).

## Manuscript and PDF

The Markdown and LaTeX Abstract, Phase 3B.5/3B.6 subsections and Conclusion reconcile the bounded claim. Section 7 is an inherited Markdown-only section. The PDF build and visual inspection follow independent manuscript source review. [PDF evidence](evidence/task4/pdf/).
