# Task 2 independent re-review — round 1

Spec: **APPROVE**

Quality: **APPROVE**

Reviewed correction base `eca3b3bee03b32c7bfffc65cba124d3ba638446c` through
head `6dc16ecf36441a99005a62eba4cee2418f5f5d72`. HEAD matched that head.
This is a scoped re-review of the four original findings and breakage introduced
by their fixes. The original `task-2-review.md` remains intact, including its
REJECT verdicts and reproduced counterexamples. This approval resolves those
findings at the new head; it does not qualify a genuine matched pair.

## Original findings

| Finding | Status | Basis |
| --- | --- | --- |
| 1. Complete zero-world snapshot coverage | **ADDRESSED** | Both cuts of both endpoints now use complete native schema validation before cardinality checks. Missing fields and null scalar counters cannot acquire zero defaults; Applications must be an array. Native empty Requests may remain explicit null. |
| 2. Complete final Run report | **ADDRESSED** | The ten native report fields and typed nested records are validated, successful report arrays must be non-null, counters must be nonnegative, and Roots join the final export. Identically truncated reports remain structurally equal but fail Raw and Supported. |
| 3. Passive authorization coordinates/selectors | **ADDRESSED** | The complete native event schema is checked; Subject, Generation, Epoch, Sequence, Intent and Attempt join the selected wire, with pre-transport phase, empty optional root/admission selectors, and nil passive InFlight. Enriched scheduler acknowledgments remain separate. |
| 4. Effective unique test selection | **ADDRESSED** | Go's flag parser interprets the supported driver flags. Exactly one exact capture selector and effective count one are required. Overrides, duplicate selectors, positional tails, selectors consumed as another flag's value, unsupported flags, and termination before selection fail. Supported single/double-dash and separate/equal-value forms are covered. |

For finding 1, inspected the frozen producer's `internal/fixture/resource.go`
Snapshot and Inspect definitions. The correction properly preserves the native
empty Requests representation while requiring explicit coverage fields; it does
not impose an invented non-null Requests requirement. The two optional fields
for the other endpoint profile remain outside this frozen ordinary-endpoint
domain.

For finding 2, compared the added report types to the frozen native RunReport,
finishRun, Intent, and RootCompletion definitions and the reused typed intent
representation. The complete-schema round trip is validation only. It does
not replace, delete fields from, or normalize the input used by complete
cross-branch equality. Moving rootView to package scope preserves the existing
outer-export checks.

For finding 3, inspected the frozen native Event definition and the actual
external.authorized emission in `internal/executor/external.go`. They support
the added coordinate checks and pre-transport phase. Correcting the synthetic
fixture's prior empty Phase matches the owning source; it does not alter the
producer or weaken an invariant.

For finding 4, checked the actual parser configuration and lawful/hostile argv
controls. Use of the standard flag parser accounts for effective values and
argument consumption; the added unique-selector count rejects repeated
assignments even when the final selector would otherwise be acceptable. The
driver flag subset and rejection policy are documented.

## Evidence inspected

The unchanged original probe source, original failing output, and overlay were
preserved under `docs/foundational/phase3b7/evidence/task2/review-r1/`. This
review byte-compared each against its original `/tmp/phase3b7-review-probes/`
file: all matched. Probe source SHA-256:
`93d7dcd40f2e2300938a8f72641ffee248701ca4401992417730f651bbe2f674`;
original failing output SHA-256:
`45fdd362a1272b3843863ee72128093852a6bff7ba5608b3aefd9d55b9319beb`.

The retained `probes-green.log` shows all four original attacks rejected at the
appropriate stage: missing endpoint fields fail World; conflicting authorization
coordinates fail Protocol; truncated reports fail Raw; overridden selection
fails Source. Every case reports Supported=false while retaining a valid seal.

Inspected the added resealed deletion/null/type/coordinate/flag controls,
including coherent two-cut deletion and identical two-branch report truncation.
The retained complete ordinary and race logs show lawful synthetic positives,
original controls, new controls, and command tests passing with no failure or
race markers. Recorded checker times are 18.005s ordinary and 216.329s race;
command times are 0.133s and 1.126s. The retained scoped Staticcheck log is
empty. Setup errors and incorrect command-path attempts are retained and
distinguished from successful verification in the correction report.

No test, race, lint, runtime, or producer command was rerun during this
re-review. No archived SQLite file was opened, no genuine capture occurred,
and no repository source, index, or branch was changed by this reviewer.
Only this new review report was written.

## New breakage and remaining boundary

No actionable new breakage found in the correction diff. No additional
out-of-scope finding was identified for final review. The fixes remain confined
to the offline checker, synthetic controls, documentation, and retained review
evidence; old authority entry points and complete equality logic are unchanged.

Task 2 is approved for the stated synthetic/offline checker scope. Source and
binary metadata remain unsigned; database/export correspondence remains
conditioned on the frozen producer. Genuine normal/race qualification,
preservation of every attempt, and subsequent empirical/manuscript
reconciliation remain separate obligations before claiming a matched concrete
witness.
