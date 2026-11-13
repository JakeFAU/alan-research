# Effect observation runbook

Work from the research repository. Use the verified extraction produced by the existing Phase 3B.5 corpus preparation process; do not open archived SQLite. The new CLI is opt-in and leaves old commands unchanged.

```sh
export GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off
export GOMODCACHE=/tmp/alan-phase3b5/modcache
export GOCACHE=/tmp/alan-phase3b6/go-build
export STATICCHECK_CACHE=/tmp/alan-phase3b6/staticcheck
export ALAN_EFFECT_OBSERVATION_CORPUS=/tmp/alan-phase3b5/corpus
export ALAN_REVOCATION_CORPUS=/tmp/alan-phase3b5/corpus

go test -count=1 -v ./research/delegationcorrespondence ./research/delegationcorrespondence/cmd/effectobservation -run '^TestEffectObservation'
staticcheck ./research/delegationcorrespondence/...
go run ./research/delegationcorrespondence/cmd/effectobservation -case /tmp/alan-phase3b5/corpus/raw-normal/N07-3888622734
```

The command emits the full JSON report, exits 0 for supported input, 1 for a checked but unsupported relation, and 2 for usage, read/digest or write failure. An unsupported relation need not imply invalid raw evidence. `ALAN_EFFECT_OBSERVATION_CORPUS` enables the new optional corpus tests; absent it, those tests explicitly skip. Production code has no extraction path or environment dependency. The old environment variable independently enables old tests.

The four selected directories are resolved from `manifest.json`. Current sealed members are normal loss `applied-unknown-revoked-1243056310`, normal N07 `N07-3888622734`, race loss `applied-unknown-revoked-2740959369`, and race N07 `N07-3721806614`.

Reproduce focused exports, source verification and compatibility evidence:

```sh
python3 docs/foundational/phase3b6/evidence/task3/verify_task3.py
python3 docs/foundational/phase3b6/evidence/task3/qualify_mutants.py "$PWD" docs/foundational/phase3b6/evidence/task3/mutants /tmp/alan-phase3b5/corpus
```

`verify_task3.py` records commands, direct exit statuses, elapsed times and environment; tests command success/unsupported/usage/output failure, builds the three relevant commands, runs scoped Staticcheck, checks formatting, exports each new report twice, and compares ten old report/raw pairs with frozen Phase 3B.5 source and archived hashes. `verify_source.py` reads the frozen runtime Git blobs at `/Users/jacob/dev/personal/alan`, owning archives and the existing extraction. Its hardcoded historical verification routes are tooling context, not production checker paths. It never executes runtime code or opens a database. `qualify_mutants.py` compiles disposable source copies and distinguishes the end-to-end causal mutant from the helper-level exact-renewal mutant.

See `evidence/task3/verdicts.json`, `legacy-comparison.json`, `commands.json`, `source-sha256.json`, compressed complete reports, mutant logs and `SHA256SUMS`. `investigation.md` records failed apparatus/candidate tests rather than silently treating them as qualification. The controller separately runs broad tests/race/build/module checks and manuscript/PDF reconciliation on the reviewed source freeze.
