# Phase 3B.5 qualification

## Frozen inputs

Research qualification runs from executable commit `4ad1f9ef37c57375d4a39eafc9a5fcd91bc0425d`, tree `70cfd58bc56f9b56f771055481a56fc224877c21`, against research main `4fea6c31eb8f438c0dea414aa5f8d7d0ecb5e8a4`. The 582-file executable manifest covers tracked Go, module and testdata inputs and evidence Python scripts. Its SHA-256 is `51ab5216312d59b65987c239b5ff66cbbc3cbd7b44b9fb9cbcc40d80c6e42936`. Later prose and PDF changes fall outside this input set. Toolchain versions are Go 1.27.3 on darwin/arm64, Staticcheck 2026.2.2 (0.8.2), and Python 3.9.6.

Runtime main was inspected at `c72c396df09c55b0be1a91b2557dcec2ef6457a8`, tree `52f4d7c9a2f573e78112be0c1f07b26834918603`. The sealed captures retain producer commit `5b4662fd682a220ca3065b15e86672b7e56567e7`, tree `44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`. This qualification executes research against historical bytes. The runtime remains unchanged.

The [corpus extraction](evidence/preflight/corpus-extraction.json) verifies six sealed archive digests and 10,634 declared regular members. The [original 48-case extraction](evidence/preflight/old48-extraction.json) verifies a separate earlier archive and 3,784 members for the normal/race regressions. [Task 3 archive verification](evidence/task3/archive-verification.json) independently checks both selected producer-mode archives and 232 members per mode. SQLite files are hashed as bytes. The [concrete runbook](05-runbook.md) describes reproduction from a fresh extraction.

## Frozen-source runs

The ten closed command records preserve the arguments, working directory, offline environment, timeout, direct exit, log digest and executable-manifest check. The final verifier confirmed all 582 frozen input hashes, all recorded command logs, formatting of changed Go files and the branch whitespace check. The four historical exports matched their archives exactly. The new model, audit and progress outputs matched the retained Task 1 and Task 2 bytes. [Qualification records](evidence/task4/qualification/), [freeze record](evidence/task4/freeze.json) and [final verification](evidence/task4/verification.json).

| Job | Recorded disposition |
| --- | --- |
| `go test -json -count=1 -timeout=60m ./...` | Exit 0 in 1,505.540 seconds. 2,633 test and subtest pass events across 51 tested packages, with 13 packages containing no tests. The job enables the earlier normal-mode inputs, prior progress checks and the two-mode opt-in corpus. |
| `go test -json -count=1 -timeout=15m -race ./research/authorityservice/... ./research/delegationcorrespondence/...` | Exit 0, 422 passing test events, seven packages passed and three had no tests. No named tests skipped. The job enables the earlier race-mode inputs and the two-mode opt-in corpus. |
| `staticcheck ./research/authorityservice/... ./research/delegationcorrespondence/...` | Exit 0, empty log. |
| `staticcheck ./...` and `staticcheck -tags=integration ./...` | Exit 1 for each. Logs exactly match the respective baselines with 13 and 31 inherited diagnostics. |
| `go build ./...` | Exit 0. |
| `go mod verify` | Exit 0. |
| `go mod tidy -diff` | Exit 0, empty log. |
| `go run -buildvcs=false ./research/authorityservice/cmd/export -config S2/SeparateGrantRevocation -assert-authority-safe` | Exit 0. The output hash matches Task 1's retained model bytes. |
| `qualify_exports.py` against frozen source | Exit 0. Two runs of each historical export matched its archived bytes. New audit and progress exports are deterministic across two runs each. |

The ordinary run recorded exactly three named skips, all optional artifact exports. [Ordinary command record](evidence/task4/qualification/research-ordinary.json) and [skip output](evidence/task4/verification.json).

| Named skip | Recorded reason |
| --- | --- |
| `examples/communication/TestCommunicationArtifactExport` | Explicit artifact export was not requested. |
| `examples/join/TestJoinArtifactExport` | `ALAN_JOIN_EXPORT` was unset. |
| `examples/substitution/a/TestCandidateAExport` | `ALAN_PHASE1I_A_OUTPUT` is required to export the first retained bytes. |

## Manuscript PDF

Cached Tectonic 0.17.0 built the independently reviewed manuscript at `4cabef20bf48453cba9dad5ac03a7d9e54a19d77`, exit 0. The installed 42-page PDF has SHA-256 `83bfaa590ad6b3cbd61ecf32043efdce93ed7ac75805d4feee1e625911de5bad`. The build record pins the Markdown, TeX and bibliography hashes. The controller inspected rendered pages 1, 2, 9, 28, 29 and 38–42, including new PDF Section 6.8 on pages 39–40 and Conclusion Section 7 on pages 40–41. Five unique box-warning patterns come from unchanged source passages. [Build record](evidence/task4/pdf/build-record.json), [visual inspection](evidence/task4/pdf/visual-review.md) and [archive member hashes](evidence/task4/pdf/manifest.json).

## Reviewed scoped evidence

- The model export preserves an explicit exit 1 falsification of universal completion, with the same JSON bytes as the successful safety assertion. All sixteen historical configurations remain in the default exporter. [Operational commands and exports](evidence/task1/).
- The independent ledger and progress checks audit the complete finite graph and retain 34 history controls. Four mutants compile and fail their targeted semantic assertions. The four historical exports match actual archives, and the old progress output matches a baseline-source build. [Graph and comparison records](evidence/task2/).
- The concrete checker covers both sealed modes, verifies declared identities and preserves all ten legacy outputs byte for byte. The c-target and fresh-gate mutants compile and reach their intended assertion failures. Serialized opportunity counts are checked in the active verdicts. [Scoped command evidence](evidence/task3/) and [active verdicts](evidence/task3/verdicts.json).

## Apparatus failures and intentional negative controls

Task 1's first two Staticcheck attempts hit the Darwin cache restriction. Setting `STATICCHECK_CACHE` to the disposable cache resolved it. Task 2 re-ran the ledger RED after a zsh wrapper obscured the direct exit. Its first witness selector missed required temporal arcs at merged targets, and the final selector uses qualifying arcs. Its first old-progress archive extraction called a Python API unavailable on the host. A checked regular-member extraction supplied the successful baseline comparison. Task 3's first projector used the lifecycle target locator for the issuer join. The first wrong-root attack failed raw validation because a mirrored settlement was stale. Corrected joins and a coherent attack isolate the intended assertion. The review separately preserved the incorrect initial serialization of unresolved opportunities. [Task evidence](evidence/) and [review index](08-reviews.md).

The preflight targeted baseline retained structured test output after its process handle expired. The direct process exit is unrecorded for that baseline. The controller's first recorder syntax check was blocked by the default Python cache path. It passed with `PYTHONPYCACHEPREFIX` redirected to the disposable directory. The final recorder records direct exits. Historical RED runs, intentionally failing mutant assertions and the universal-completion counterexample carry their own dispositions.
