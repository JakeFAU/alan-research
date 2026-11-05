# Phase 3B.4 independent qualification runbook

Use the executable pair and per-file inputs in [the freeze manifest](evidence/task5/freeze/pair.json), not a later HEAD merely sharing a branch name. The qualified code is research `0c303c3d43d6582fed0a9caac24f50e71d844a19` and runtime `fc60391f48b46a8c8635afc46882db8338b6c93b`. Evidence-only commits are separately identified. No old archive is overwritten. The source map includes Go/module files, qualification scripts and testdata outside `docs/` and `paper/`; every recorded final command compares this map before and after execution.

## Environment and evidence safety

The recorded host uses Go 1.27.3, darwin/arm64, Staticcheck 2026.2.2 (0.8.2). Use the pinned module graph and a writable offline module cache:

```sh
export GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off
export GOMODCACHE=/tmp/alan-phase3b4-qualification-modcache
export GOCACHE=/tmp/alan-phase3a-go-build
export STATICCHECK_CACHE=/tmp/alan-phase3b4-staticcheck
```

The cache provenance is retained with the evidence. Runtime integration tests use loopback endpoints and subprocess owners. They require execution permission for those facilities. This corpus contains no live model endpoint evaluation. Never open an archived SQLite database in place: copy its database, WAL and SHM together to a disposable directory before any SQLite/Open inspection. The independent checker reads exported JSON and hashes files without opening SQLite.

Use fresh output directories. The command recorder in `evidence/task5/parent-final-gates/alan-phase3b4-gate-recorder.py` records exact cwd/argv, environment, source maps, output digest, exit, elapsed time and watchdog, and refuses to overwrite existing command evidence. A watchdog, build failure or panic is apparatus failure, not semantic rejection.

## Runtime corpus and broad gates

From the frozen runtime repository, run the full ordinary suite and both full integration modes with `-count=1 -timeout=60m`; `-race -tags=integration ./...` includes the ordinary runtime tests under race. The final actual argv, readers and source guards are recorded under [parent-final-gates](evidence/task5/parent-final-gates/).

```sh
go test -json -count=1 -timeout=60m ./...
ALAN_PHASE3B4_EVIDENCE=/tmp/alan-task5-final-normal \
ALAN_PHASE3B3_EVIDENCE=/tmp/alan-task5-final-v1-normal \
  go test -json -count=1 -timeout=60m -tags=integration ./...
ALAN_PHASE3B4_EVIDENCE=/tmp/alan-task5-final-race \
ALAN_PHASE3B3_EVIDENCE=/tmp/alan-task5-final-v1-race \
  go test -json -count=1 -timeout=60m -race -tags=integration ./...
```

Set both old-reader environments to independently rebuilt exact sources before these commands, as recorded in the gate manifests: `ALAN_PHASE3B3_OLD_READER` is the pre-v1 CLI reader, while `ALAN_PHASE3B4_OLD_READER` is the frozen v1 reader for rejecting 15/16. They are different compatibility controls. Their source/build provenance is archived; a binary hash alone is not a source identity.

The promoted-parent test cannot share those full commands' first invocation because Go may execute its package before the durable package has produced the K parents. After each complete corpus, run it explicitly:

```sh
ALAN_PHASE3B4_PUBLIC_PARENTS=/tmp/alan-task5-final-normal \
ALAN_PHASE3B4_PROMOTED_OUTPUT=/tmp/alan-task5-final-promoted-normal \
  go test -json -count=1 -timeout=10m -tags=integration ./internal/executor \
  -run '^TestPhase3B4PromotedPublicParents$'
```

Repeat with race roots and `-race`. These are ancillary internal-writer copies of exact public K01/K02 parents; they do not establish a public adapter-conflict production schedule. Report the initial optional skip and its later executed replacement separately.

Run actual repository build, module verification/tidy, format, Staticcheck and required vet gates in both relevant tag modes. Keep expected inherited research lint diagnostics exact; an arbitrary nonzero lint result is not acceptable. Do not change unrelated source to turn the historical baseline green.

## Independent checker and adversarial campaign

From the frozen research repository:

```sh
ALAN_PHASE3B4_CHECK_ROOT=/tmp/alan-task5-final-normal \
ALAN_PHASE3B4_PROMOTED_CHECK_ROOT=/tmp/alan-task5-final-promoted-normal \
ALAN_PHASE3B4_CHECK_OUTPUT=/tmp/task5-independent-normal \
ALAN_PHASE3B4_ATTACK_OUTPUT=/tmp/task5-attacks-normal \
  go test -json -count=1 -timeout=5m ./research/delegationcorrespondence \
  -run '^TestBinding'
```

Repeat against the race corpus with `-race`. When `ALAN_PHASE3B4_CHECK_ROOT` is set, every attack resolves its genuine parent by exact case name in that new corpus. The committed compressed testdata is development regression input, never substituted for a missing final parent. N10h's failed Create is explicitly no-Store/no-raw-verdict.

For a reviewable single case:

```sh
go run ./research/delegationcorrespondence/cmd/check \
  -binding-case /absolute/extracted/case-directory
```

The report separates raw validity, exact principal/auxiliary custody, physical endpoint evidence, ordered owner continuity and restricted model projection. Unsupported projection cannot repair invalid raw evidence. `CheckRaw`/`-input` without a principal declaration checks v2 raw policy but explicitly leaves principal custody unchecked; use the case interface for the full result. The old v1 corpus CLI and expectations remain available unchanged.

Run checker-code mutations only in isolated source copies:

```sh
python3 research/delegationcorrespondence/qualify_mutants.py \
  /absolute/frozen/research/source /tmp/task5-new-mutants \
  /tmp/alan-task5-final-normal
```

Each mutant must build. M01 must accept the future-selection false allow; M02 must both falsely refuse genuine P10 and accept the q2-only N01 false allow; M03 must accept the real target2 message's contradictory A1 admission; M04 must accept the false fresh gate after genuine revocation. The restored wildcard interpreter is a separate defect-control mutation. An unrelated early guard, parser error, missing fixture, panic, or build failure does not qualify as a kill. Exact changed predicates and assertion logs are retained.

The equal-total custody attack preserves genuine prior q2 service/result lineage and all ordinary reduction/admission mirrors. It must pass raw authority/reduction while the explicitly declared principal association rejects the q reply carrying q2's result. A genuinely matching separate q2 grant/service is lawful and is never labeled an invalid substitution by itself.

## Prior regressions and unchanged models

The final full research ordinary suite uses both the new corpus and the preserved old48-case corpus. `ALAN_TASK3_RAW` points to the historical v1 roots, **not** the new81-case root; some prior tests read their original absolute `Directory` entries. The preserved extraction/preflight records describe those paths. Use `ALAN_TASK3_PROGRESS=1` for the expensive prior progress audit. Run the affected research model/checker race packages, mutual integration/race and Phase3B2 durable integration/race according to the exact recorded commands and the [prior runbook](../phase3b3/05-runbook.md).

Build all four unchanged export/audit commands with `-buildvcs=false`. Execute each twice into new files and compare each byte-for-byte with the preserved export, decoding gzip for phase3b2. `evidence/task5/apparatus/export-twice.py` records this exact comparison. These finite graph/export regressions do not prove runtime semantics or enumerate arbitrary histories.

## Costs and archives

Each mode contains five actual observations at each n=0,1,8,32. Preserve constructor, first-decision-window, fresh-gate-window, Replay and disposable-copy Open durations, the run/byte/debit sidecars and exact query plan. The first two authority windows include runner/boundary overhead. Allocations were not measured. No speedup, asymptotic whole-operation bound, or performance threshold follows; inherited prefix/projection work remains distinct from the exact indexed admission lookup.

`artifacts.json` hashes every packaged evidence file. Raw corpus archives preserve complete DB/WAL/SHM companions and JSON, without SQLite opening during packaging. Mutant and preliminary apparatus archives use content-addressed deduplication: `files.json` maps original relative paths to `sha256/<digest>` blobs, whose exact bytes reconstruct each isolated source and log. Archive-format sidecars document this layout. Keep preliminary apparatus failures separate from final exact-source results and preserve both.

## Concurrent-main integration reproduction

[The prospective-main addendum](evidence/prospective-main/README.md) records the later exact detached combination, all 15 actual commands, fresh corpus roots, source maps, archive/member hashes and independent closure. Its `plan.json` job list and `summary.json` identify the checks actually repeated. Original scientific qualification and prospective integration qualification have distinct identities; recreate the stated parent combination rather than silently testing a later main. Permit loopback listeners/subprocesses and use fresh output paths.
