# Reproduce the concrete correspondence

Run from the research repository root with Go 1.27.3 and the repository's dependencies already cached. The commands read sealed Phase 3B.4 artifacts and write disposable copies. Database members are hashed as byte strings. No SQLite connection or runtime campaign is required.

```sh
export GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off
export GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b5/go-build
export STATICCHECK_CACHE=/tmp/alan-phase3b5/staticcheck
export ALAN_REVOCATION_CORPUS=/tmp/alan-phase3b5-revocation-reproduction
evidence=docs/foundational/phase3b5/evidence/task3
python3 "$evidence/prepare_corpus.py" "$ALAN_REVOCATION_CORPUS" --prepare
go test ./research/delegationcorrespondence ./research/delegationcorrespondence/cmd/revocation \
  -run '^TestRevocation' -count=1 -v
go test ./research/delegationcorrespondence -run '^TestBindingDenseAttacks$' -count=1 -v
```

Use a fresh extraction path with `--prepare`. Omit it to verify the declared files in an existing extraction. The manifest produced by preparation is the case locator required by the corpus-backed tests. Without that environment variable, those tests skip and the CLI usage test still runs.

The opt-in command uses the separately named correspondence executable.

```sh
scratch=$(mktemp -d /tmp/alan-revocation.XXXXXXXX)
go build -o "$scratch/revocation" ./research/delegationcorrespondence/cmd/revocation
"$scratch/revocation" -case "$ALAN_REVOCATION_CORPUS/raw-normal/N06-4077261617"
```

Exit 0 means supported. Exit 1 still emits the verdict and means a semantic or scope check prevents support. Exit 2 is a usage, input apparatus, or output error. Inspect `RawValid`, `Physical`, `Legacy.ContinuityFailure`, and `Reasons` before interpreting `Supported`. The final capture carries the raw fact identifiers and projection events used by the independent ledger audit.

Compare both modes and the unchanged legacy report against Task 3's baseline. Use the same extraction path for both binaries since the legacy report retains that path.

```sh
mkdir "$scratch/base"
git archive 47d4e4a9e546c34e91f41087810f5b99937064a5 | tar -xf - -C "$scratch/base"
(cd "$scratch/base" && go build -o "$scratch/check-base" ./research/delegationcorrespondence/cmd/check)
go build -o "$scratch/check" ./research/delegationcorrespondence/cmd/check
python3 "$evidence/compare_reports.py" "$ALAN_REVOCATION_CORPUS" \
  "$scratch/check-base" "$scratch/check" "$scratch/revocation"
python3 "$evidence/qualify_mutants.py" "$scratch/mutants" "$scratch/mutant-logs"
staticcheck ./research/delegationcorrespondence ./research/delegationcorrespondence/cmd/revocation
```

The comparison runs the opt-in report twice for each case and compares bytes, then compares the old and new legacy command bytes and checks that the embedded legacy JSON agrees. Mutants compile in disposable Go overlays. The script requires each mutant binary to fail the corresponding semantic assertion and writes its exact commands, source hashes, exits and compressed logs.

The recorded extraction hashes, comparison verdicts, and assertion evidence are under [evidence/task3](evidence/task3/). Both the archive verification and the reports retain historical producer provenance. Report byte hashes include the recorded directory and are comparable on the same extraction path.
