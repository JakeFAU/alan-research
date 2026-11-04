# Phase 3A reproduction

Run from a full checkout of this branch. The root module declares Go 1.27.3.
The new model imports only the Go standard library; it neither runs ALAN's
executor nor calls an endpoint. The commands below require no model credentials.

```sh
export GOTOOLCHAIN=local
export GOCACHE=/tmp/alan-phase3a-reproduction-cache
export STATICCHECK_CACHE=/tmp/alan-phase3a-reproduction-staticcheck
export GOPROXY=off
export GOSUMDB=off

go test -count=1 -v ./research/twoagent/...
go test -count=1 -race -v ./research/twoagent/...
staticcheck ./research/twoagent/...
go build -buildvcs=false -o /tmp/alan-phase3a-export ./research/twoagent/cmd/export
/tmp/alan-phase3a-export > /tmp/alan-phase3a-model.json
cmp docs/foundational/phase3a/evidence/model.json /tmp/alan-phase3a-model.json
```

The binary emits both profiles and their finite graph counts. BFS retains all
edges, including repeated-state edges and cycles, and has no depth cutoff.
`CanComplete` counts states from which A's retained consumption is reachable;
it does not mean every schedule terminates, B has released, channels have
drained, or the pair is globally quiescent. Reverse reachability establishes
that property separately from local state and transition safety and custody.

## Reproduce the rejected hypothesis

```sh
/tmp/alan-phase3a-export -assert-naive-recovery \
  > /tmp/alan-phase3a-rejected.json \
  2> /tmp/alan-phase3a-rejected.log
```

**Expected exit status: 1.** The diagnostic must say `FALSIFIED`. This is the
deliberately rejected implication from weak local safety and original-request
deduplication to useful recovery, not a broken passing regression suite.
Unexpected exit 0 or 2 does not reproduce the result. The JSON carries the
shortest recovered-loss witnesses and equal recoverable projections for the
two opaque result values. Each loss trace has ten transitions in this model's
alphabet. The shortest custody-violation prefix is distinct from the later
recovered loss.

`CheckTrace` validates event prefixes using separate causal accounting and does
not call the model transition function. Tests check every reachable edge appended
to its BFS predecessor trace and reject malformed correlation/grants, fabricated
receipt values, premature ACK, duplicate application and consumption, and
consumption before recording. This is an independent implementation check of
the declared model, not a mechanized proof or runtime refinement.

## Repository and manuscript checks

The repository's ordinary command is:

```sh
GOTOOLCHAIN=local go test -count=1 -timeout=60m ./...
```

Existing packages use loopback servers and child processes; an environment
denying those cannot qualify them. The new model itself needs neither. Root
`./...` does not cover separately declared example-consumer modules; no such
module or production implementation changed in this slice.

From `paper/`, rebuild with `tectonic --keep-logs alan-dissertation.tex`.
The [manuscript review](evidence/manuscript-review.md) identifies the inspected
rendered pages and retained warnings. The PDF and both manuscript sources contain
the same added mathematical argument. Existing historical layout differences
were not audited as part of Phase 3A.

## Evidence identity

`evidence/source.json` binds the exact model sources, module declarations and
research/manuscript sources by SHA-256. `evidence/commands.json` records commands,
exit statuses and log paths. These hashes bind the first-slice experiment rather
than requalifying earlier archived phases. `evidence/model.json` is deterministic
for the frozen source; timings in test logs are not deterministic.

Read [results](03-results.md) for assumptions and open obligations. In particular,
an authentic lookup service, retained result digest/body, additional crash, or
dynamic authority protocol changes the model being qualified.
