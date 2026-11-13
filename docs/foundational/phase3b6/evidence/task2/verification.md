# Task 2 verification

Implementation base `9a28d2ca4a4e031f1c98860154bbf0ed4aa7fb6e`. Exact final Go sources, including unchanged operational inputs, are identified in [source-sha256.txt](source-sha256.txt). [commands.json](commands.json) records the final commands, exit codes, output hashes and environment. [SHA256SUMS](SHA256SUMS) covers evidence files other than itself. No broad repository suite was run by this worker; the controller owns that qualification.

All Go commands used:

```sh
export GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off
export GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b6/go-build
export STATICCHECK_CACHE=/tmp/alan-phase3b6/staticcheck
export PYTHONPYCACHEPREFIX=/tmp/alan-phase3b6/python-cache
```

| Command/check | Exit/result | Evidence |
| --- | --- | --- |
| `go test -count=1 ./research/effectobservation -run '^TestLedger'` against compiled `AuditTrace` stub returning empty ledger/nil error | 1; 8 lawful and 14 hostile assertions failed | `red.log.gz`; exact initial sources `red-oracle.go.gz`, `red-oracle_test.go.gz` |
| Same initial focused tests after implementing independent ledger | 0 | `green.log.gz` |
| `go test -count=1 -v ./research/effectobservation/...` with all final tests | 0 | `tests.log.gz` |
| `go run ./research/effectobservation/cmd/audit` twice | 0 both; byte equality checked | `audit.json.gz`, `commands.json` |
| `go run ./research/effectobservation/cmd/export` twice | 0 both; byte equality checked | `graph.json.gz`, `commands.json` |
| `gofmt -l research/effectobservation` | 0, empty output | `format.log.gz` |
| `staticcheck ./research/effectobservation/...` with writable cache | 0, empty output | `staticcheck.log.gz` |
| `go version`; `staticcheck -version` | 0; Go 1.27.4 darwin/arm64; Staticcheck 2026.2.2 (0.8.2) | respective compressed version logs |
| `git diff --cached --check`, initial evidence patches | 2; unified-diff context lines triggered space-before-tab diagnostics; changed patches to zero context | `diff-check-initial.log.gz` |
| `git diff --cached --check`, final task diff | 0, empty output | executed before commit |
| `python3 docs/foundational/phase3b6/evidence/task2/run-mutants.py` | 0; all required semantic mutations compiled and failed targeted assertions | `mutants.json`, four `.patch` and `.log.gz` pairs |

The two uncompressed audit exports have the same SHA-256 recorded in `commands.json`; likewise the graph exports. The graph digest is unchanged from Task 1 (`142ae6cd3ce9622ecd535d15ad67b923983d8991e357da596858c86239be8d15`). Compression uses zero mtime. Output equality was checked directly in Python, not inferred from hashes.

The mutant runner makes disposable module copies under `/tmp/alan-phase3b6`, runs `go test -c -o <disposable>/ledger.test ./research/effectobservation`, then invokes that binary with the exact `-test.run` regex in its manifest/log. It first compiles the unmodified source and verifies that all four targeted tests pass. Every changed source then compiles with exit 0 and its targeted test exits 1 with an ordinary assertion failure. Packet deletion produces an invalid rejection of lawful ingress; application-as-retention changes the observed custody; response acceptance after reopen produces an invalid acceptance; c revocation breaks the surviving control root. These are semantic failures, not compile errors or panics. Exact source replacements and mutant hashes are retained.

RED can be reproduced in a disposable copy of Task 1's package by restoring the two compressed RED Go sources and running its recorded command. The final package adds descriptor and graph tests/controls and does not retain stubs in executable sources. No accidental compile failure was counted as RED, and no discovered model counterexample was suppressed. The applied/unretained and delayed-application witnesses are the intended boundary counterexamples and remain in the exported audit.

The initial diff whitespace check diagnosed the leading context marker before Go indentation inside saved `.patch` artifacts. The final runner emits zero-context patches, preserving the exact semantic replacements and source hashes; all four compiled mutant checks were rerun successfully after that artifact-format correction. No Go source changed for this correction.
