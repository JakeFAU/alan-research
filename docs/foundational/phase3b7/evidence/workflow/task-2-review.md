# Task 2 independent review

Spec: **REJECT**

Quality: **REJECT**

Reviewed base `f238b778ac786d3efbd86f94e98068439b9e4ed2` through head
`eca3b3bee03b32c7bfffc65cba124d3ba638446c`, the supplied diff, Task 2 brief and
report, contract, native producer protocol, added checker/adapter/tests, and
retained verification evidence. HEAD matched the reviewed head. This review
does not qualify a genuine pair or reopen the separately approved producer.

## Strengths

- Complete outer JSON comparison preserves fields, ordered arrays, number
  tokens, descriptor strings, and null/empty distinctions. Difference paths
  are deterministic and unambiguous. Both exact-renewal checks include
  Settings and Roots.
- The report separates seals, source, raw authority, causal protocol, endpoint
  world evidence, full equality, final-report equality, and model mapping.
  Producer completion does not replace those checks.
- The authority adapter is narrow, additive, and accurately labeled as reuse.
  The retained four historical normal/race byte comparisons support unchanged
  old report output; old pinned entry points are unchanged in this diff.
- Source-conditioned database/export linkage and unsigned producer trust are
  disclosed. Synthetic fixtures are clearly distinguished from genuine
  execution evidence; the shared seed and four branch exports match the
  intended protocol.
- Retained ordinary/race/adapter logs report passing tests, and the retained
  Staticcheck log is empty. Meaningful comparison and causal RED failures are
  retained. Those results do not cover the additional attacks below.

## Actionable findings

### 1. [P1] Require explicit zero-world snapshot coverage fields

Location: `research/matchedcapture/world.go:66-70`.

`endpointSnapshot` decoding defaults absent counters to zero and absent arrays
to nil, and the cardinality check treats both as complete empty coverage.
Deleting **Requests, Applications, Total, and Contacts** from both
`zero/endpoint-b/stopped.json` and `zero/endpoint-b/final.json`, then resealing
the synthetic dossier, still returns `World.Valid=true` and `Supported=true`.
The outer three-field check does not validate the snapshot's required fields.
This permits missing application evidence to be certified as zero, contrary
to the explicit complete-endpoint-coverage prerequisite. Require every native
snapshot member and its expected JSON type before evaluating cardinalities;
do not let missing/null counters stand for zero. Preserve the native empty
array representation. Add rehashed deletion/null controls for both endpoints
and cuts while retaining lawful zero and one controls.

### 2. [P2] Validate the complete final Run report schema

Location: `research/matchedcapture/store.go:61-65`.

Final reports are decoded into maps and checked only for `Error` and
`Report.UnknownExternal`. Replacing **both** final `run.json` files with
`{"Error":"<nil>","Report":{"UnknownExternal":["[\"a\",1,4,0]"]}}`
still returns `Raw.Valid=true`, complete report equality, and
`Supported=true`. Equality correctly compares the supplied bytes, but neither
side is a complete native Run report. This loses the declared complete-report
observation while certifying support. Validate all required native report
fields and types, rejecting unsupported/malformed schemas, without filtering
or rewriting either input to the equality comparator. Add a coherently
truncated-both-sides control.

### 3. [P2] Reject contradictory native authorization-boundary coordinates

Location: `research/matchedcapture/protocol.go:200-202`.

The authorization boundary is joined only by Intent and Attempt. Changing its
Subject to `"wrong"` and Epoch/Sequence to `999` in the zero owner's native
`boundaries.jsonl`, leaving the fixed selected raw request and pre-call event
unchanged, still returns `Protocol.Valid=true` and `Supported=true` after
resealing. Thus the claimed boundary corroboration accepts a conflicting
selected authorization identity. Validate this event's complete native
coordinate/selector schema and join its emitted selected coordinates to the
retained request. Keep any distinction between passive native events and
inherited enriched scheduler acknowledgments explicit; do not assume their
InFlight representations are identical. Add a rehashed conflicting-coordinate
control alongside the existing acknowledgment tests.

### 4. [P2] Check the effective test selection, not occurrence of a flag

Location: `research/matchedcapture/source.go:122-128`.

The source check accepts the expected `-test.run` argument anywhere in argv.
Appending `-test.run=^OtherTest$` after the expected selector still returns
`Source.Valid=true` and `Supported=true`. Go flag parsing uses the later
assignment, so this metadata does not describe selection of the frozen capture
test. Reject repeated/conflicting selectors or validate the effective parsed
selection, including argument termination and supported flag spellings. Add
an override control; the existing substring-only negative test does not cover
this case.

## Focused verification performed by this review

One offline checker-only probe used a Go overlay to add a temporary test file
without changing repository source, index, or branch. Each case started from
the existing lawful synthetic fixture, made the stated mutation, and resealed
the scratch fixture. No genuine archive was edited or opened through SQLite;
no producer/runtime experiment ran. No broad suite, race suite, or lint suite
was rerun.

```sh
env GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off \
  GOMODCACHE=/tmp/alan-phase3b5/modcache \
  GOCACHE=/tmp/alan-phase3b7/go-build \
  go test -overlay=/tmp/phase3b7-review-probes/overlay.json \
  -run '^TestIndependentReviewProbes$' -v ./research/matchedcapture
```

Exit status: **1**, because all four rejection assertions failed. Every case
printed `Supported=true Seal=true Source=true Raw=true Protocol=true
World=true`. Package duration: 0.634s. Probe source:
`/tmp/phase3b7-review-probes/probe_test.go`; full output:
`/tmp/phase3b7-review-probes/result.log`. These scratch probes are review
evidence, not changes to the submitted implementation.

Correct the four acceptance gaps and obtain independent re-review before
freezing Task 2 for genuine qualification. The existing REJECT findings remain
part of the review history; a subsequent review must explicitly resolve them.
