# Task 2 post-capture independent re-review — round 2

Spec: **APPROVE**

Quality: **APPROVE**

Finding: native authorization admission ordinal falsely rejected — **ADDRESSED**.

Reviewed base `193039406557cb37be815d5f846b22a75d48fde6` through head
`90724cce9082f1ca978d785efa743e851a344617`; HEAD matched the reviewed head.
Scope is this post-capture checker correction and new breakage introduced by
the fix. Original review reports remain intact.

## Corrected interpretation

The round-one review's zero-ordinal interpretation was wrong. I inspected the
event constructor but missed the Run transport wrapper that modifies the event
before emission. That review error is acknowledged here; the old approval is
not evidence that the initial frozen checker accepted a genuine pair.

The frozen call path supports exactly ordinal 1 for this selected admission:

1. `internal/executor/external.go:174` constructs the event.
2. `internal/executor/runner.go:203-215` installs the transport closure, which
   copies `op.authority.Ordinal`, Epoch and Sequence into `work.event` before
   calling `transportExternal`. Lines 223-237 obtain the admitted opportunity
   and pass that same opportunity to the worker transport.
3. `internal/executor/external.go:177-190` reaches the modified event before
   constructing/sending the HTTP request.
4. `internal/durabletest/phase3b7_owner_test.go:144-181` stamps and logs the
   event before forwarding to the inherited observer. This does not populate
   passive InFlight.

I read the owning source and verified all three whole-file lengths/hashes and
recorded excerpts against the prepared checkout and frozen preparation
manifest. The correction is supported by the full call path, rather than
solely by the observed capture values.

## Finding and controls

`research/matchedcapture/protocol.go` now requires the authorization event's
ordinal to equal **1** and the selected worker's ordinal. Existing pre-call
validation already requires that worker's ordinal to be 1 and joins the
selected wire identity. The unchanged raw binding validator independently
requires the original allowed gate's ordinal 1, selected source and attempt
digest. All existing Subject/Generation/Epoch/Sequence, Intent/Attempt, phase,
root, passive InFlight, process identity and ordering checks remain.

The implementation does not accept both zero and one. The synthetic positive
now carries the worker's actual admitted ordinal. Changing the old
wrong-ordinal test from 1 to 2 corrects an invalid expectation; it does not
remove the wrong-ordinal invariant. New branch-isolated and end-to-end resealed
controls reject omitted/default zero, 2 and -1 with Seal and Raw still valid.
The old coordinate, selector, malformed-schema and identity controls remain.
The isolated RED demonstrates that the old helper accepted omitted zero
independently of the other branch's rejection, avoiding a masked negative test.

## Native excerpt custody

I independently decoded `testdata/native-authorization.json.gz` and verified
its decompressed SHA-256:
`19cfd7318d98162b730044c81628570d9c0abe0a0ad04125f40f39121ec87419`.
Its keys are exactly normal/zero, normal/applied, race/zero and race/applied.
For each branch I checked:

- Boundary and pre-call bytes equal their original JSON/JSONL members exactly.
- All three source-member lengths/hashes match both `native-inputs.json` and
  the original pair's member manifest.
- Selected attempt and allowed-gate descriptor strings equal the original
  stopped export strings, without descriptor rewriting.
- The authorization event, pre-call worker and allowed gate all have ordinal
  1; event/worker/wire identity and source coordinates join.

These were read-only JSON/JSONL checks. They establish excerpt custody and the
mapping under review, not a new whole-pair or SQLite qualification.

## Verification and preservation

Inspected both retained RED logs, the full ordinary GREEN log and affected
race GREEN log. Native controls for all four branches, the synthetic positive,
ordinal negatives, and retained authorization controls pass. The recorded
ordinary checker/command times are 21.582s/0.107s; affected race time is 39.622s.
The scoped Staticcheck log is empty. No broad Go tests, race suite or lint was
rerun during this review.

Executed the read-only `verify-preserved.py` (exit 0), and independently read
the four initial reports. The old frozen binary, initial summary, original
excerpt inputs and initial report hashes remain as recorded. Both initial
normal/race report pairs are byte-identical within their mode and retain
Supported=false and Protocol=false. They separately retain valid
Seal/Source/Raw/World verdicts and byte-identical complete paired exports and
final reports. Those independent observations do not turn the initial checker
rejection into a pass.

The correction diff changes no equality algorithm or observation, producer,
source anchor, raw authority helper, contract, or genuine raw evidence. The
comparison documentation changes only the native ordinal mapping. There is no
recapture or new whole-pair result in this correction. This reviewer did not
open SQLite, mutate source/index/branch, rebuild or overwrite the frozen
checker, or run the complete checker against genuine pairs. Only this new
review report was written.

## Decision and next boundary

No actionable new breakage found in this scoped fix. No additional
out-of-scope finding was identified. The ordinal finding is **ADDRESSED**;
the correction is approved for a new, separately identified checker freeze
and rechecking of the same immutable captured pairs.

This is explicitly a **post-capture checker correction**. The initial rejected
reports and earlier reviews remain part of the research record. The controller
must separately freeze/build the corrected checker and retain repeated
reports on the existing raw pairs before claiming its final qualification
verdict. Source-conditioned, unsigned-producer and bounded single-host limits
remain unchanged.
