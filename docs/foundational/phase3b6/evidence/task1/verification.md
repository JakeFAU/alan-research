# Task 1 operational qualification

Implementation base `e339c2db337cd5b5a93990b8a42160994a3ec2f4`, tree `ef5a04949ca26f7021f2fc3128bc856599cf8362`. The [source manifest](source-sha256.txt) identifies the exact implementation and test bytes. [SHA256SUMS](SHA256SUMS) identifies the export, summary and compressed logs. The decompressed graph SHA-256 is `142ae6cd3ce9622ecd535d15ad67b923983d8991e357da596858c86239be8d15`.

Commands ran from the repository root with Go 1.27.4 darwin/arm64 and Staticcheck 2026.2.2 (0.8.2). Go commands used:

```sh
export GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off
export GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b6/go-build
```

| Command | Exit and evidence |
| --- | --- |
| `go test -count=1 ./research/effectobservation -run '^(TestSeparateStages\|TestLossPreservesExternalPacket\|TestLateResponseCannotReachReopenedOwner\|TestExactRootAndHistoricalAuthority\|TestRefusalAfterUnretainedApplication\|TestRetainedResultSurvivesLoss)$'` with compiled transition/projection stubs | 1, [six failing assertions](red.log.gz) for missing launch, revoke or loss |
| Same focused command with operational transitions and projections | 0, [green log](green.log.gz) |
| `go test -count=1 ./research/effectobservation/...`, initial exporter test | 1, [build diagnostic](export-test-build-failure.log.gz) for unavailable `iotest.ErrWriter` |
| `go test -count=1 ./research/effectobservation/...` with local writer returning `io.ErrClosedPipe` | 0, [suite log](tests.log.gz) |
| `go run ./research/effectobservation/cmd/export > /tmp/alan-phase3b6/task1-model-a.json` | 0 |
| `go run ./research/effectobservation/cmd/export > /tmp/alan-phase3b6/task1-model-b.json` | 0 |
| `cmp /tmp/alan-phase3b6/task1-model-{a,b}.json` | 0, exact bytes matched |
| `gzip -n -c /tmp/alan-phase3b6/task1-model-a.json > docs/foundational/phase3b6/evidence/task1/graph.json.gz` | 0, [complete graph](graph.json.gz) |
| `gofmt -l research/effectobservation` | 0, empty |
| `staticcheck ./research/effectobservation/...` | 1, [default cache denied by sandbox](staticcheck-default-cache.log.gz) |
| `STATICCHECK_CACHE=/tmp/alan-phase3b6/staticcheck staticcheck ./research/effectobservation/...` | 0, empty |

For the summary, count all graph states and outgoing arcs and group arcs by event kind. Use the first BFS state for each qualifying application value among fresh-refused states. Select the first stopped, revoked, launched state with no ingress, then append its ingress and first application arcs. For the retained control, take the consume arc of the first reopened/revoked BFS state that offers it, preserving the shorter source prefix and final consume descriptor. [Summary JSON](summary.json) preserves those paths with full descriptors.
