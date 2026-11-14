# Phase 3C.2 capacity composition model

Research-only fixed A→B→C extension of unchanged `contractcomposition` and
`twoagent`. Exact original-held capacity is separate from permanent request-bound
permission, cumulative expenditure, and durable receipt custody. Application
retains the actual result and releases capacity atomically.

From the repository root:

```sh
GOCACHE=/tmp/alan-phase3c2/gocache go test ./research/capacitycomposition/...
GOCACHE=/tmp/alan-phase3c2/gocache go test -race ./research/capacitycomposition/...
GOCACHE=/tmp/alan-phase3c2/gocache staticcheck ./research/capacitycomposition/...
GOCACHE=/tmp/alan-phase3c2/gocache go run ./research/capacitycomposition/cmd/export > /tmp/alan-phase3c2/export.json
```

The exporter takes no arguments and emits deterministic indented JSON to stdout;
errors go to stderr with nonzero exit. It checks source hashes relative to the
repository root and executes the full declared factored probe domain. Probe
execution is intentionally substantial; the encoder smoke test does not repeat
it. No database or runtime is opened.

The 72 complete finite graphs contain all arcs and reverse-reachability
CanComplete labels. The ten concrete witness histories begin at genuine empty
custody. Probe records summarize synthesized factored representatives, not the
full concrete product and not genuine whole-chain histories. Invalid input to
error-returning APIs is returned as an error; `Initial`, `Successors`,
`ConcreteSuccessors`, `Complete` and `Explore` panic on invalid declared input.
`Project` is a structural map; validate concrete sources before making semantic
claims about them. `State.Validate`/`ConcreteState.Validate` check structural
premises, not reachability or WorldResult authenticity.

See the frozen [contract](../../docs/foundational/phase3c2/00-contract.md),
[conditional argument](../../docs/foundational/phase3c2/01-argument.md), and
[JSON schema](../../docs/foundational/phase3c2/02-evidence-schema.md). Independent
checker results and manuscript reconciliation are subsequent-task obligations.

Task 1 archives the exact plain stdout bytes as
`docs/foundational/phase3c2/evidence/task1/export.json.gz`, using gzip with zero
mtime and an empty filename header. `artifact-manifest.json` records both raw
and packed SHA256 values and byte lengths. This is archival packaging only;
the frozen schema and CLI stdout remain plain JSON, and the captured `Sources`
records are unchanged. To inspect the raw JSON:

```sh
gzip -dc docs/foundational/phase3c2/evidence/task1/export.json.gz > /tmp/alan-phase3c2/export.json
```
