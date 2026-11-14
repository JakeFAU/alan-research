# Completion-preserving continuation interfaces

This research-only finite model composes the unchanged Phase 3C.2 K relation
with five explicit B continuation monitors and 324 bounded A/C controllers.
It tests initial existential completion, exact successful visible languages,
and normalized grouping. It separately retains every-state CanComplete and
nonblocking labels. The equal-language doomed monitor demonstrates why those
claims must remain distinct.

From the repository root:

```sh
GOCACHE=/tmp/alan-phase3c3/gocache go test -v -timeout 30m ./research/componentinterfaces/...
GOCACHE=/tmp/alan-phase3c3/gocache go test -v -timeout 30m -race ./research/componentinterfaces/...
GOCACHE=/tmp/alan-phase3c3/gocache STATICCHECK_CACHE=/tmp/alan-phase3c3/staticcheck-cache staticcheck ./research/componentinterfaces/...
GOCACHE=/tmp/alan-phase3c3/gocache go run ./research/componentinterfaces/cmd/export > /tmp/alan-phase3c3/export.json
```

`Export(io.Writer) error` is the only exported package function. CLI takes no
arguments and emits deterministic full JSON with source hashes. It expects the
repository root as working directory and returns diagnostics on stderr. All
other model functions are internal, consume constructed descriptors, and are
not an untrusted JSON parser. Exact signature/role/automaton validation precedes
product construction. Language comparison consumes the generated automata or
explicit test NFAs; it is not an input-validation boundary.

The export retains all 1,980 reachable product graphs, 1,620 full pairs of
grouped graphs, ten complete epsilon-subset-pair comparisons, all descriptors,
equal-language substitution rows (successful and blocked), and eight genuine
concrete model histories. Histories use inherited custody preparation and real
successor steps; they are not runtime or external endpoint captures. No old
model, production code or dependency changes.

See [contract](../../docs/foundational/phase3c3/00-contract.md),
[source argument](../../docs/foundational/phase3c3/01-argument.md) and
[exact schema](../../docs/foundational/phase3c3/02-evidence-schema.md).
The source argument is conditional on the specified private-state and product
rules; finite tests are not proofs for arbitrary ALAN contexts. In particular,
this is not complete-agent substitution, general nonblocking or full abstraction.

Task 1 archives raw JSON stdout as deterministic `export.json.gz` (zero mtime,
empty filename), with byte lengths and SHA-256 in `artifact-manifest.json`.
Decompress with `gzip -dc docs/foundational/phase3c3/evidence/task1/export.json.gz`.
The two false hypotheses are retained as `false-hypotheses_test.go.txt` beside
their raw RED execution. They are deliberately excluded from the passing suite;
the suite instead asserts their concrete counterexamples. To reproduce the
negative run, copy that file temporarily to this package's `attack_test.go`,
run `go test -v -run TestAttack ./research/componentinterfaces`, then remove the
temporary copy. The expected exit is 1 for both falsified hypotheses; the
initial missing-symbol RED logs are separately identified as compile failures.
