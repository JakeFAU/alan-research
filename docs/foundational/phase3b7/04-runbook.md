# Phase 3B.7 runbook

## Recheck the two sealed pairs

Use the normal and race [archive records](evidence/capture/) and the accepted
checker at research commit `90724cce9082f1ca978d785efa743e851a344617`. The
[comparison contract](03-comparison.md) specifies the checker's inputs and
separate verdicts. The exact original and corrected command records are under
[initial-commands](evidence/qualification/initial-commands/index.json) and
[corrected-commands](evidence/qualification/corrected-commands/index.json).

From the research repository, set `scratch` to a new absolute scratch directory
and `modcache` to the populated offline Go module cache. Verify each
`raw.tar.gz` against its adjacent `archive.json` by length and SHA-256 before
extracting it. The archive record also lists each member's name, mode, byte
length and digest. Extract to separate empty directories. Build from the
corrected freeze and keep reports outside the extracted raw directories.

```sh
phase=docs/foundational/phase3b7
mkdir -p "$scratch/normal-v1" "$scratch/race-v1"
for pair in normal-v1 race-v1; do
    tar -xzf "$phase/evidence/capture/$pair/raw.tar.gz" -C "$scratch/$pair" || exit
done
env GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off \
    GOMODCACHE="$modcache" GOCACHE="$scratch/go-build" \
    go build -o "$scratch/matched-check" ./research/matchedcapture/cmd/check || exit
for pair in normal-v1 race-v1; do
    "$scratch/matched-check" -pair "$scratch/$pair" > "$scratch/$pair.1.json" || exit
    "$scratch/matched-check" -pair "$scratch/$pair" > "$scratch/$pair.2.json" || exit
    cmp "$scratch/$pair.1.json" "$scratch/$pair.2.json" || exit
done
```

Expected exits are 0 for both corrected checks of each pair. Compare their
report digests with [corrected-checks.json](evidence/qualification/corrected-checks.json).
The full reports are gzip-compressed command outputs alongside the command
records. To reproduce the preserved pre-capture rejection, build the checker
from `193039406557cb37be815d5f846b22a75d48fde6` in a separate checkout and
run it twice on each of the same directories. Expected exit is 1, with the
authorization-coordinate rejection recorded in
[initial-checks.json](evidence/qualification/initial-checks.json).
Keep those reports separate. Leave all SQLite members unopened and read their
bytes only for custody checks.

## Reproduce the experimental protocol

The existing archives are sufficient for the offline check above. A fresh
execution creates a separate dossier with new process and Store identities.
Use [the pinned preparation procedure](02-protocol.md#preparation-and-build-identity)
to clone the runtime detached at `5b4662fd682a220ca3065b15e86672b7e56567e7`
and add its four named test overlays. Record the Go toolchain, source manifest,
owner and endpoint binary identities, and both source and checker freezes.
The recorded experiment used Go 1.27.4 on darwin/arm64, an offline module
cache and the integration test driver. The complete build and test arguments
and environment are in the individual `job.json` records.

In the prepared checkout run the following once with a fresh
`ALAN_PHASE3B7_EVIDENCE` output directory for each mode. Add `-race` for the
race mode. Run zero before applied as the producer specifies.

```sh
go test -json -count=1 -tags=integration -timeout=10m ./internal/durabletest \
    -run '^TestPhase3B7MatchedCapture$'
```

Preserve the exit and output and archive the stopped raw directories with
`python3 "$research/$phase/evidence/qualification/archive-capture.py" "$raw" "$archive_dir"`,
where `research` is the absolute research checkout path, `raw` is the
stopped capture root and `archive_dir` is a fresh directory outside it.
The archiver inventories the raw tree, verifies it stayed unchanged, and
checks every archived member against that inventory. Seal each attempt before
checking it. Retain producer failures, unsupported checker reports and
equality differences as observed. A later apparatus correction requires its
own freeze and evidence.
