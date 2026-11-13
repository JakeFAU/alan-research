# Phase 3B.7 qualification

## Capture freeze and custody

The controller froze 627 tracked acceptance and executable inputs at research
commit `193039406557cb37be815d5f846b22a75d48fde6`. The runtime was pinned to
`5b4662fd682a220ca3065b15e86672b7e56567e7` (tree
`44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`) with four test-only overlays.
The source inventory checks cover 637 runtime members and four overlays
before and after capture. The ordinary and race capture commands ran once
each using Go 1.27.4 on darwin/arm64 and both exited 0. The exact source,
arguments, environment, producer status, toolchain and executable
identities are in the freezes, command records and archived `source.json`.
[Initial freeze](evidence/qualification/freeze.json) ·
[Initial checker identity](evidence/qualification/checker-before.json) ·
[Producer before](evidence/qualification/producer-before.json) ·
[Producer after](evidence/qualification/producer-after.json)

| Mode | Capture wall time | Archive members / bytes | Archive SHA-256 |
| --- | ---: | ---: | --- |
| Normal | 4.036 s | 122 / 369,092 | `0d2b423ea6f611be28595ac71d973507c7ccf6d2444e68fc05c28af61fdf4b65` |
| Race | 11.622 s | 122 / 369,260 | `19be235cec79f5555eaaa4d92263af884f36c9fa4ef45ff814f81de9bc4c6292` |

The archive command verified each regular member's path, mode, length and
SHA-256 against the stopped raw directory. Each archive was sealed before
the initial checker ran on that pair. The outer records include the producer member
seals. The wall times are single command measurements.
[Initial commands](evidence/qualification/initial-commands/index.json) ·
[Normal archive inventory](evidence/capture/normal-v1/archive.json) ·
[Race archive inventory](evidence/capture/race-v1/archive.json)

## Initial rejection and corrected acceptance

The initial checker rejected both genuine pairs on both invocations per
mode. The seal, source, raw and world verdicts passed and the complete
exports and final reports matched exactly. The protocol verdict rejected the
selected native authorization coordinates in both branches. It required
ordinal 0 from the event constructor. The frozen Run transport wrapper copies
the admitted worker's ordinal into that event before emitting it. All four
actual branches emit ordinal 1. The synthetic fixture and first re-review
missed the wrapper. The initial reports are preserved and the original
protocol rejection is an apparatus finding.
[Initial checks](evidence/qualification/initial-checks.json) ·
[Root-cause source extracts](evidence/task2/review-r2/root-cause.json)

The corrected checker at `90724cce9082f1ca978d785efa743e851a344617` binds the
passive authorization event to the worker's exact ordinal 1. The native
excerpts from all four branches pass and synthetic ordinals 0, 2 and -1
reject. The independent post-capture re-review approved this correction.
The corrected freeze contains 629 inputs. Its six-path delta consists of
the checker change, its comparison documentation and test/fixture changes.
The producer, contract, equality comparator, raw data and sealed archives
were held fixed. There was no recapture.
[Corrected freeze](evidence/qualification/freeze-corrected.json) ·
[Six-path delta](evidence/qualification/freeze-correction-delta.json) ·
[Native input provenance](evidence/task2/review-r2/native-inputs.json) ·
[Reviews](07-reviews.md)

| Mode | Initial exits | Corrected exits | Corrected repeat report SHA-256 |
| --- | --- | --- | --- |
| Normal | 1, 1 | 0, 0 | `edf131b32f95075db836c86d16992f779fe03a18fbde327acba71dfe56aaa3e3` |
| Race | 1, 1 | 0, 0 | `4d7326b200b8bacc679eb0f3c1a57e16b3bab6628756855a82cb00a1e820c38f` |

Repeated reports are byte-identical within each mode and checker revision.
The corrected checks pass seal, source, raw consistency, protocol,
complete export and final-report equality, model mapping and endpoint world
verdicts. Checker identities were checked before and after their use. The
archive and its members were rechecked after acceptance. The source
provenance and executable metadata are unsigned, and database-to-export
fidelity rests on the pinned producer's copy and export procedure.
[Initial checks](evidence/qualification/initial-checks.json) ·
[Corrected commands](evidence/qualification/corrected-commands/index.json) ·
[Corrected checks](evidence/qualification/corrected-checks.json) ·
[Initial checker recheck](evidence/qualification/checker-after-final.json) ·
[Corrected checker recheck](evidence/qualification/checker-corrected-after.json)

## Executable checks and compatibility

The producer's ordinary and race barrier-control suites each passed nine
controls, and its scoped integration Staticcheck emitted zero diagnostics.
The checker and command suites passed focused ordinary and race tests and
the authority adapter passed its focused race check. The post-capture
ordinal change has its own native-input RED/GREEN evidence and focused
verification. Earlier checker controls cover malformed and incomplete
inputs, chronology, identity joins and richer local differences invisible
to the model projection.
[Producer checks](evidence/task1/) · [Checker checks](evidence/task2/)

Four fresh historical effect-report comparisons match byte for byte across
the unchanged legacy entry points. They cover N07 and the applied-unknown
loss case in normal and race modes. These are rechecks of sealed historical
inputs. The full repository ordinary qualification reported for Phase 3B.6
is historical. Fresh Phase 3B.7 work uses the affected package and historical
corpus scope.
[Four byte regressions](evidence/task2/historical-byte-regression.json) ·
[Phase 3B.6 qualification](../phase3b6/06-qualification.md)

The affected corpus suites cover `./research/effectobservation/...` and
`./research/delegationcorrespondence/...` with the historical corpus inputs.
Both ordinary and race commands exited 0. Each recorded 422 test and subtest
pass events, seven tested packages, three packages without tests and no
named skips. The ordinary and race command wall times were 32.981 s and
267.564 s, respectively.
[Verification summary](evidence/qualification/verification-summary.json) ·
[Exact commands and logs](evidence/qualification/verification-commands/index.json)

`go build ./...`, `go mod verify`, `go mod tidy -diff` and scoped
Staticcheck passed. The initial formatting selection included a retained
historical reviewer probe. The corrected current-source formatting check
had empty output, with no intervening source edits. Global Staticcheck
returned exit 1 with 13 ordinary and 31 integration diagnostics, both
outputs byte-identical to the Phase 3B.6 baseline. The package verifier
confirmed all 629 corrected freeze inputs, both 122-member archives and
the hashes and outputs of 23 recorded commands.
[Verification summary](evidence/qualification/verification-summary.json) ·
[Package verification](evidence/qualification/package/verification-before-review.json)

## Manuscript

The edited Abstract, Phase 3B.6 closing paragraph, new Phase 3B.7 subsection
and Conclusion are synchronized between Markdown and LaTeX. The inherited
Markdown-only Section 7 remains absent from LaTeX and PDF. The bibliography
and existing TeX layout warnings are inherited.
[Build and visual verification](evidence/manuscript/visual-review.md)
