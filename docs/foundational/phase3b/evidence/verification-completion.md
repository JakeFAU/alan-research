# Verification coverage

The ordinary root suite began on merged base
`8583527a8bc2d0b71a1e9115ded23eb96c8ac86c`, before the new packages were created:

```sh
GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build GOPROXY=off GOSUMDB=off \
  go test -count=1 -timeout=60m ./...
```

It exited 0. The existing adaptive package took 1,477.571 seconds; no timeout or
test failure was hidden. Original Go/module sources remained unchanged. Final
ordinary tests for both new model packages and the new fixture constructor
complete coverage of the final `go list ./...` inventory: **53 packages**, with
**43 passing test packages** and **10 built without test files**. The exact union
and log path for every package are in `package-coverage.json`.

Focused new-model race and Staticcheck gates passed. The runtime fixture's normal
and race integration runs each pass all 16 profile/order cases; its focused
Staticcheck gate passed. The independent portable raw-evidence audit passes all
32 archived cases and exactly reproduces the committed audit JSON. Deterministic
model exports match. The naive serial completion assertion intentionally exits
1 and apparatus errors are separately tested as exit 2.

The package inventory and model build emitted an existing module metadata-cache
write warning under the sandbox while exiting 0. The model build warning is
retained in `model-export.log`; it did not prevent creation/use of the exporter.
The root ordinary suite itself ran with the permissions needed for existing
loopback/child-process tests.

This coverage does not stand for a whole-repository integration/race run or the
separately declared consumer modules. No production package, module dependency,
Phase 3A source/evidence or consumer module changed. Task-specific reviewers
inspected the frozen sources/evidence and their minor assertion fixes. The
whole-change review is recorded separately.
