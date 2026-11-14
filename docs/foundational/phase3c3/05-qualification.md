# Phase 3C.3 qualification and provenance

## Source boundary and attribution

Merged baseline: `ab01204a0b5121b81379043576ee7bbe46118a19`. Accepted Go producer:
`c8230eb2ed2bf3cd45635f72d5eef84bcc7372fa`; accepted independent audit and Task 3
input: `c4b4e503a1f7e135ea25e06fc918c008e33f0100`. New Go code is confined to
componentinterfaces and its exporter. Task 3 changes documentation, bibliography,
PDF and qualification/verification artifacts. Frozen producer sources, contract,
argument, schema, prior models, dependencies and evidence remain unchanged.

Task 1 ran scoped Go tests/race/Staticcheck and exact exports; Task 2 reconstructed
the full abstract roster, replayed eight histories and challenged the checker.
The controller ran fresh root build/module/package checks. Task 3 archives their
exact records, verifies current source identity and builds/inspects the manuscript.
No unchanged Go suite, global lint or million-case probe campaign was repeated.

## Historical qualification and fresh input identity

[Baseline inputs](evidence/qualification/baseline-inputs.json) capture all 669
tracked merged-baseline Go/module/workspace inputs. Each was freshly compared to
both current bytes and its merged-baseline Git blob: all match. This complete
snapshot does not retroactively expand historical test input manifests.
The [Phase 3C.2 qualification](../phase3c2/05-qualification.md) preserves the
original 637-entry manifest, separately identified 20 preexisting files and
12 then-new capacity files, with their actual qualification boundaries.
Nested modules and archived fixtures are not implied to be root ./... coverage.

The [package comparison](evidence/qualification/package-comparison.json) confirms
76 = 74 inherited + componentinterfaces and its exporter, with no removals.
The older 70-package ordinary run (57 pass, 13 no-test skips) plus subsequent
scoped model-package runs form the inherited coverage; they were not a single
74- or 76-package invocation. Source-identical inherited coverage and fresh
scoped checks support combined ordinary package coverage with that attribution.
No fresh whole-root suite, whole-root race/integration run or runtime campaign
is claimed. Global Staticcheck is not newly checked or claimed clean: prior
qualification retains 13 ordinary and 31 integration-tag inherited diagnostics.
No go vet was invoked for this packet.

## Executed records reused without repetition

| Scope / command | Recorded result | Attribution |
| --- | --- | --- |
| `go test -v -timeout 30m ./research/componentinterfaces/...` | Nine top-level tests pass; exporter has no tests; exit 0 | Task 1 green-test.json |
| `go test -v -timeout 30m -race ./research/componentinterfaces/...` | Nine top-level tests pass; no race report; exit 0 | Task 1 green-race.json |
| `staticcheck ./research/componentinterfaces/...` | No diagnostics; exit 0 | Task 1 green-staticcheck.json |
| `go run ./research/componentinterfaces/cmd/export` | Two Task 1 captures and Task 2 fresh capture byte-identical | Task 1 artifact-manifest.json and Task 2 fresh-export-comparison.json |
| `python3 -m unittest discover -s docs/foundational/phase3c3/tools -p 'test_*.py' -v` | Seven methods pass, including 53 semantic hostile controls | Task 2 unittest.command.json |
| Python audit CLI on plain/gzip and exact hostile inputs | Both valid encodings accepted; all 55 hostile inputs rejected | Task 2 accepted-index.json / hostile-index.json |
| `go build ./...` | Exit 0, empty output, 13.612s | Controller build-cache-writable |
| `go mod verify` | Exit 0, all modules verified | Controller module-verify |
| `go mod tidy -diff` | Exit 0, empty output | Controller module-tidy-diff |
| `go list ./...` | Exit 0, clean 76-package roster | Controller package-list-cache-writable |

Exact argv, cwd/source head where recorded, environment, output and exits remain
in the captures. Controller output.log files are archived as deterministic gzip;
decompression preserves the exact bytes, including empty output. Its
[runner](evidence/qualification/run_check.py) and source-bound job/result records
are archived unchanged. Task 1's initial Staticcheck cache failure and initial
successful dependency enumeration with a cache warning remain distinct from
clean reruns. Compilation REDs remain distinct from the semantic negative tests.
Task 2's initial missing-checker RED remains retained; no hostile input was
initially accepted.

The exporter binds 27 inputs: 22 compiled local non-test Go sources, root
module files and three frozen semantic documents. Task 3 verifies those and all
125 Task 2 artifact records plus four checker sources. Task 1/Task 2 raw exports
match at 25,384,964 bytes, SHA-256
`a073a22fe8ac62e05f2d15b0f5146b19db9606e38ce4619261e7e4108e27269d`.
The old independent Python custody checker is reused unchanged and explicitly
attributed, not presented as an entirely new oracle or an old full probe rerun.

## Manuscript and PDF

Both source forms add the same Phase 3C.3 content; changed Abstract/Conclusion
match. The [primary literature note](evidence/manuscript/primary-literature.md)
attributes publisher full-text/metadata inspection to the controller: Tang and
Moor, Discrete Event Dynamic Systems 34, 125–161 (2024), DOI
10.1007/s10626-024-00394-2. This is background on composition/nonblockingness;
the packet neither instantiates their priority/conflict abstraction nor claims
a new general automata theorem. Failed De Nicola/Hennessy access attempts yield
no added verified-text attribution.

The PDF operation marker ran exactly once, edit/count1/pdf, immediately before
first manuscript authoring. The [manuscript manifest](evidence/manuscript/manifest.json)
binds actual baseline/current cached Tectonic builds, extracted text, page counts,
warning comparison, render hashes and visual inspection. Inherited Markdown-only
Section 7, section/theorem numbering and format-specific bibliography differences
remain explicit; only changed/new mathematical content is asserted aligned.
No global manuscript normalization or equality claim is made.

Fresh rebuilds have 48 baseline and 50 current pages. Warning lines match exactly:
four underfull hboxes, one 5.2184pt overfull hbox and two math-mode font warnings.
All 11 same-index changed pages (1, 2, 4, 5, 24, 45-50) and their continuations
were inspected: 16 pages total (1-6, 23-25, 44-50). An initial two-line Abstract
spill was corrected by shortening the new summary; final page 1 contains the
whole Abstract. No new visual defects were found. Existing warnings and old
code-block pagination remain disclosed. The tracked deliverable is
`paper/alan-dissertation.pdf`. Cached tooling added no dependencies.

Reproduce artifact/source/alignment verification from the repository root:

```sh
python3 docs/foundational/phase3c3/evidence/manuscript/verify.py
```

This verifier reads hashes, exact archives, counts, source alignment and local
links; it does not rerun the accepted producer/auditor or claim to mechanize the
conditional proof. The [results](04-results.md) retain the initial/residual and
model/runtime boundaries. Review and delivery status are controller-owned and
separate from these executed checks. General Phase 3C remains open.
