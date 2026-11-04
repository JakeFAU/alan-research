# Phase 3B.1 reproduction

Run from the repository root with Go 1.27.3. The new model uses the standard
library; its projection tests also import the unchanged Phase 3A model. The
runtime fixture uses existing public author/runtime APIs and existing SQLite
dependencies. No model credentials, paid endpoints or new dependencies are used.
Offline commands below require the repository's dependencies already cached.

```sh
export GOTOOLCHAIN=local
export GOCACHE=/tmp/alan-phase3b-reproduction-cache
export STATICCHECK_CACHE=/tmp/alan-phase3b-reproduction-staticcheck
export GOPROXY=off
export GOSUMDB=off

go test -count=1 -v ./research/mutualobligations/...
go test -race -count=1 -v ./research/mutualobligations/...
staticcheck ./research/mutualobligations/...
go build -buildvcs=false -o /tmp/alan-phase3b-export ./research/mutualobligations/cmd/export
/tmp/alan-phase3b-export > /tmp/alan-phase3b-model.json
cmp docs/foundational/phase3b/evidence/model.json /tmp/alan-phase3b-model.json
```

The default export includes all four complete finite graphs' counts and shortest
witnesses; it does not serialize every graph state. The tests check graph closure,
the independent causal ledger, every edge projection to the frozen Phase 3A
relation, oracle noninterference, adverse traces and CLI exit statuses. No search
depth cutoff or probabilistic scheduler is used.

## Rejected completion claim

```sh
/tmp/alan-phase3b-export --profile serial/serial --assert-naive-completion \
  > /tmp/alan-phase3b-naive.json 2> /tmp/alan-phase3b-naive.stderr
```

Expected exit status **1**, with `naive completion claim falsified` on stderr.
This means the selected graph contains a reachable state without any completion
suffix. Exit 2 signals an apparatus/argument/output error; exit 0 does not
reproduce the serial counterexample. An asymmetric profile with the same flag
returns 0. The shortest Trap/admitted-wait/recovered-wait witnesses have 2/6/10
events, and exported unfair retry cycles distinguish possible completion from
inevitable progress.

## Public-runtime cases and raw evidence

```sh
go test -count=1 ./examples/communication/mutual
go build ./examples/communication/mutual
ALAN_PHASE3B_RUNTIME_EXPORT=/tmp/alan-phase3b-runtime-normal \
  go test -tags=integration -count=1 -v ./examples/communication/mutual
ALAN_PHASE3B_RUNTIME_EXPORT=/tmp/alan-phase3b-runtime-race \
  go test -race -tags=integration -count=1 -v ./examples/communication/mutual
staticcheck -tags=integration ./examples/communication/mutual
```

Each run executes 16 cases (four profiles by two Send orders by two request
delivery orders). Export is optional and produces an exact program pair, compact
summary and `evidence.json.gz` per case. Gzip headers omit timestamps; Store
namespaces differ between fresh runs, so the runtime evidence is not claimed
byte-identical. The model export is deterministic.

Recheck the committed raw runtime artifacts independently of the runner:

```sh
python3 docs/foundational/phase3b/evidence/runtime-audit.py \
  docs/foundational/phase3b/evidence /tmp/alan-phase3b-audit.json
cmp docs/foundational/phase3b/evidence/runtime-audit.json /tmp/alan-phase3b-audit.json
```

The audit accepts an evidence root containing `runtime-normal/` and `runtime-race/`
as its first argument and output JSON as its second. The test's native export
paths above can be copied beneath such a root for a new audit. It checks actual
observations and selected history; delivery ACKs and quiescence alone do not pay
application obligations. A serial run is quiescent while both originals remain
unpaid. Whole-Store reopen renews both roots, and final holds keep them open.

## Repository and manuscript

The ordinary repository gate is:

```sh
GOTOOLCHAIN=local go test -count=1 -timeout=60m ./...
```

It includes existing long-running adaptive tests and requires loopback listeners
and child processes for some packages. This slice started that suite on the
merged base before creating the new packages, then checked the unchanged-source
condition and covered every added package separately. The package-coverage
artifact records the union against the final `go list ./...` inventory. This is
ordinary root-module coverage, not the entire integration-tagged/race suite or
the separate example-consumer modules. Focused integration/race coverage above
qualifies the new fixture; no production package or consumer module changed.

From `paper/`, run `tectonic --keep-logs alan-dissertation.tex`. The manuscript
review records affected rendered pages, warnings and the scope of Markdown/TeX
synchronization. No new literature attribution is introduced in this slice;
`paper/references.bib` remains unchanged.

`evidence/source.json` binds source files, prior-model inputs, module declarations,
research documents and manuscript inputs by SHA-256. `commands.json` records
actual command statuses and log provenance. Raw observations and qualification
counts do not establish dynamic authority, cancellation, distributed ownership,
external-effect atomicity or arbitrary-size compositional closure; read the
[results](03-results.md) for the precise proposition and remaining obligations.
