# Preserved historical false acceptance R4.1

The independent review changed only six history rows' logical generation from 1 to 2. Observation rows, consumed markers, results and other retained data remained unchanged. The pre-fix inspection oracle still reported two paid originals. This is a defect in the research evidence validator, not a demonstrated production runtime authorization defect.

`before.json.gz` and `after.json.gz` are the exact inspections from the successful review probe. `pre-fix-runtime-test.go.txt` retains the exact implementation that accepted the mutation (SHA256 `2d1ea47e0a340bd44901c2fce6206349c23bc5195a5b99f9a500065b54a3ed76`). The manifest binds those inputs and the probe. Final source has regressions for both all-selected-history and final-consumption-only generation mutations; the original REJECT remains in the adjacent review record.

To reproduce historical acceptance without changing the positive database or checkout, create a Go overlay mapping the repository's `internal/durabletest/phase3b2_runtime_test.go` to `pre-fix-runtime-test.go.txt`, and a new virtual `internal/durabletest/review4_probe_test.go` to `runtime-probe.go.txt`. Use absolute paths in the JSON `Replace` object. Then run:

```sh
GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build GOPROXY=off GOSUMDB=off \
  go test -overlay=/tmp/phase3b2-historical-overlay.json -tags=integration \
  ./internal/durabletest -run '^TestReview4HistoryGeneration$' -count=1 -v
```

The archived probe writes temporary before/after inspection files at its declared `/tmp` paths. Its log records the six changed selections and erroneous `paid=2 chain=<nil>` outcome. Historical acceptance is the expected failure witness; it must not be interpreted as final passing qualification.
