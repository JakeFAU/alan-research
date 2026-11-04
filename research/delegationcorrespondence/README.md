# Independent request-delegation correspondence experiment

This research-only package checks the frozen Phase 3B.3 generation-one,
two-static-root request-authority fixture. It imports no public runtime,
executor, authorization helper, database driver or Replay implementation.
`CheckRaw` does not invoke the model. The separate generator, projection and
progress query use unchanged public `authorityservice` APIs.

The independent checker validates profile13 raw JSON, original deployment and
complete permission tuples, real root membership, committed program-produced
intents, lifecycle grant/revocation decisions, immediate original-request
selection, immutable reservations, original serving-root charges, historical
fresh gates, result/settlement custody, and actual caller selections. It
reconstructs mirrored composition facts from source rows; primary gates are
checked against prior policy, not their own claimed allow bit or hash.
Unknown tables/columns, duplicate primary keys, duplicate JSON keys,
noncanonical typed descriptors and mismatched indices fail closed. JSON
integers remain exact int64 values. The actual hostile Tables-only export is
sparse: omitted known tables mean zero exported rows. Ordinary Inspect exports
require the complete table set. Neither path fills a missing row from a model
expectation.

Accepted request, ready and result-bearing reply admissions independently
require their exact original six-field communication permission, including the
committed source revision. Local input requires its exact revision-zero route.
Source authorization remains a separate liveness check.

Selections also require their retained causal predecessors at the actual
composition prefix: message authorization, lifecycle decision, original
reservation and allowed effect result, or terminal source history. Final
observation/settlement existence alone cannot justify an earlier selection.
There is no invented global admission timestamp; local inputs remain external
input assumptions without a source composition fact.

The small instruction interpreter independently checks the actual `emit`,
`receive`, `branch`, `set` and `end` slice. It is not a proof of the full reducer
or author compiler. Adapter preparation remains a declared trusted boundary:
the exported deployment maps capability to endpoint; add-one/add-ten select
actual delta1/delta10 despite input delta999. Prepared/result hashes, request
identity, endpoint wire, original receipt and returned values are checked.
This package does not qualify conflict schema14 or general program/lifecycle
feature combinations.

## Independent result layers

- `Final` / each phase `Verdict`: raw structural, authority and original-custody
  checks; an honest refusal is a valid history.
- `External`: endpoint ingress/application evidence, local gate/result boundary
  ordinals and matching Open/post-open activation snapshots. Contacts are not
  applications. Unsigned retained exports do not prove contact completeness,
  and the immutable wire has no fresh gate ordinal.
- `ContinuityFailure`: immutable history/attempt/grant/reservation/debit/fact
  preservation across the supplied prefix snapshots, plus exact original and
  current permission, program, deployment, profile, static root/member and
  origin sets. Run source/initial-state fields remain fixed while execution,
  activation epochs and observation consumption may advance.
- `Projection`: raw-driven model events, evidence references, explicit action
  and result abstraction, model AuditTrace, stutters, unsupported reasons and
  mismatch. `Scope` is explicit. A valid final logical service subsequence is
  not a chronological crash/recovery correspondence result.
- `Outcomes`: actual supplied Run, Replay, Open and typed/ordinary Inspect
  reports, retained separately. These are not the checker's authority oracle.

The known `immediate_selection_narrows_model_admission` discrepancy remains:
the raw immediate-next-reduction contract legitimately refuses after BA then
ready selection, while unchanged LateGrant allows delegated service. The
reversed service laundering control is outside the model's A-servicing/B-serial
program order. Separate-grant revocation, unknown atomic application/receipt
prefixes, repeated exposures and rejected grant issuance are explicit
unsupported mappings. For shared-process cuts, no paired per-agent crash
expansion is proved; a final service-subsequence audit remains visible while
full recovery projection is marked unsupported. No model crashes are invented.

Process-scoped boundary source runs and epochs must agree with that owner's
actual Open and post-open snapshot. `Owners` preserves their paths, hashes,
process IDs and run state; `OwnerBindings` ties ordinal events to that context.
Captured Open renewal increments open runs once and leaves terminal epochs
unchanged. Two fixtures lack an initial pre-Open snapshot: their first owner's
`RenewalChecked` is false and the missing evidence is explicit, while its
Open/post-open binding and later stopped-to-reopened renewal are checked.
Ordinary controls capture final execution before reopened inspection; process
cuts capture stopped execution before reopened and final recovery evidence.
These context checks do not prove a whole-Store model crash expansion.

Message admissions do not have their own composition order. A projection uses
a causal linearization: a real accepted transport is placed after its source
authorization and before its actual selected observation. A durable reply
admission expands to model receive-result plus record, with both events tied
to the same actual original. Distinct original IDs never merge when numerical
results happen to match. Two physical B service tuples are a declared
many-to-one whole-action abstraction for the nonrevoking trace, not equality
of permission algebras, shared-resource interference or revoke behavior.

## Reproduce

Use the repository's pinned local Go toolchain. Raw roots are immutable exports;
no command opens SQLite. `-original-prefix` maps preserved absolute index and
case paths to an extracted root without changing their bytes/digests.

```sh
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3b3-task3-build \
  go run ./research/delegationcorrespondence/cmd/corpus > corpus.json

GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3b3-task3-build \
  go run ./research/delegationcorrespondence/cmd/progress > progress.json

GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3b3-task3-build \
  go build -o /tmp/delegation-check ./research/delegationcorrespondence/cmd/check

/tmp/delegation-check -root EXTRACTED_ROOT -original-prefix ORIGINAL_ROOT > check.json
/tmp/delegation-check -input RAW_SNAPSHOT.json > snapshot-check.json
/tmp/delegation-check -root EXTRACTED_ROOT -original-prefix ORIGINAL_ROOT \
  -expect expectations.json > qualification.json
```

Checker exit codes: **0** means actual results match explicitly supplied
expectations (or every input is raw-valid with no mismatch when no expectations
are supplied); **1** means semantic rejection/discrepancy; **2** means usage,
malformed outer JSON, inaccessible evidence, digest disagreement or another
input-apparatus error. `go run` wraps nonzero program exits, so use the built
binary when checking exact exit codes. An expectation map must cover every
case exactly and contains `{ "RawValid": bool, "Mismatch": string }`.
Expectations never affect policy or projection. Expected mismatch remains
mismatch in JSON. No verdict is inferred from case name, directory or
`Classification`; that supplied label is retained as `ClaimedClassification`.

The unqualified full frozen-corpus check intentionally exits1: it contains five
invalid hostile histories and the named model/runtime mismatch. A complete
explicit expectation map can qualify those findings with exit0 without
turning them into correspondence agreement. Both invocations must be retained.

```sh
ALAN_TASK3_RAW=/tmp/alan-phase3b3-runtime-frozen/normal ALAN_TASK3_PROGRESS=1 \
  GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3b3-task3-build \
  go test ./research/delegationcorrespondence/... -count=1 -v

ALAN_TASK3_RAW=/tmp/alan-phase3b3-runtime-frozen/race ALAN_TASK3_PROGRESS=1 \
  GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3b3-task3-build \
  go test -race ./research/delegationcorrespondence/... -count=1 -v
```

Without those environment variables, only the external-corpus and exhaustive
progress tests skip; the committed positive fixture, malformed inputs, fifteen
retained-input mutants, corpus generation, strict encoding and CLI exit tests
still execute. The positive fixture is exact exported JSON, not an expected
model reconstruction. `cmd/mutations RAW_SNAPSHOT` emits full mutated inputs,
original/mutated hashes, changed-cell references and first failed edges.

`altered-holder` changes the selected root's holder. The separate
`altered-recipient` changes the retained lifecycle claim's recipient and its
decision copy while preserving committed origins and external receipt rows; it fails
original-source validation. Semantic field paths and complete before/after cell
bytes distinguish these controls. Corpus checking caches the already
digest-verified JSON/JSONL bytes for every later checker layer and never rereads
those files without error handling.

`GenerateCorpus` retains all eight categories and the unchanged 33-history /
four-ancestry prior corpus. `QueryProgress` enumerates every reachable state and
arc, performs safe reverse reachability, preserves excluded-state counts and
counterexamples, and audits an actual longest shortest completion suffix. The
LateGrant equivalence is conditional on no permanent BA first refusal; the
Revocable predicate additionally checks both originals' permit/applied state.
These are existential finite-graph results, not a runtime liveness theorem.
The old sixteen graphs and their archived results remain untouched; this new
query is not a replacement for their existing full audit commands.
