# Reciprocal request-order witness

`New` authors two ordinary retained roots in one Store, with one request each in
opposite directions. Each root echoes the admitted peer request into a reply and
checks its own reply's request identity and value. The four profiles change only
whether each root serves its peer before awaiting its own reply. No runtime or
author operations are added.

The integration driver first commits both original requests and then publicly
delivers and acknowledges both. It executes actual lowered instructions until
the request intent appears: initialization precedes each Send, so one host Commit
per source-level Send would establish the wrong cut. Both Send orders and both
request-delivery orders are exercised for each profile.

| Profile | Cases | Requests consumed | Results consumed | Unpaid obligations |
| --- | ---: | ---: | ---: | ---: |
| serial/serial | 4 | 0 | 0 | 2 |
| A servicing, B serial | 4 | 2 | 2 | 0 |
| A serial, B servicing | 4 | 2 | 2 | 0 |
| both servicing | 4 | 2 | 2 | 0 |

These are bounded executable observations under the declared driver and ordinary
runner, not universal scheduling claims. Each source stays open at an explicit
unsupplied local hold after validation and service. Consumed replies discharge
the experiment's application obligations; root termination is not the target.
This separates the obstruction from the earlier Send/End authorization race.

In the serial case both request observations are admitted but neither matches
its root's current result Receive. Explicit selection attempts are rejected
without changing retained state. The Store is quiescent under the runner's
semantics, with no pending deliveries and no `UnresolvedIntents`, while two
application obligations remain unpaid. Quiescence is not completion. The model's
retry cycles are a separate abstraction: these programs do not issue retries.

Repeated Run and whole-Store close/open/Run retain the serial wait, original
identities, sequences, program bytes, requests, and consumption facts. Opening
renews **both** roots' activation epochs. Replay checks historical internal
consistency and agrees with the retained original roots. This does not qualify
independent-owner or A-only process recovery, full model/runtime refinement, or
external application/receipt atomicity. Request selection and reply commitment
are separate ordinary commits. The runtime's delivery acknowledgments are not
the model's application-level receipt-release ACK protocol.

Reproduce the focused qualification from the repository root:

```sh
export GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build
export GOPROXY=off GOSUMDB=off
export STATICCHECK_CACHE=/tmp/alan-phase3b-staticcheck
go test ./examples/communication/mutual -count=1
go build ./examples/communication/mutual
ALAN_PHASE3B_RUNTIME_EXPORT=/tmp/alan-phase3b-runtime-normal \
  go test -tags=integration ./examples/communication/mutual -count=1 -v
ALAN_PHASE3B_RUNTIME_EXPORT=/tmp/alan-phase3b-runtime-race \
  go test -race -tags=integration ./examples/communication/mutual -count=1 -v
staticcheck -tags=integration ./examples/communication/mutual
```

Export is optional. Each of the sixteen case directories contains exact public
program JSON, `summary.json`, and deterministic-header `evidence.json.gz` with
the provision, initial commit pre/post states, original request admissions,
complete public inspections (including observations, intent descriptors, and
selected histories), runner reports, refused serial selections, and replay.
Program bytes also appear base64-encoded in the provision. Store namespaces vary
between independent runs; gzip has no timestamp, but exports are not claimed
byte-identical across fresh Stores. All counts can be checked from the raw
admission/observation/history facts; delivery ACKs alone do not discharge them.

The frozen scope is in [the Phase 3B contract](../../../docs/foundational/phase3b/01-contract.md).
