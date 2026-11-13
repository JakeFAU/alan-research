# Phase 3B.7 — Matched Concrete Effect Histories

## Question and status

Can two genuine executions have opposite selected endpoint application facts
while preserving the same declared local observation after admission, owner
loss, reopening, exact-root revocation and fresh refusal? This contract is
written before the new capture campaign. It specifies an experiment, not a
result or a general distributed impossibility theorem.

The user requests the missing admitted-but-unapplied counterpart to Phase
3B.6's historical applied-but-unrecorded execution. The controller chooses a
fresh pair from a common closed seed to remove incidental identity differences;
this choice is an experimental design decision, not a claim of user approval of
this document. Existing Phase 3B.6 artifacts and interpretations remain intact.

## Domain and alternatives

Use the frozen runtime commit `5b4662fd682a220ca3065b15e86672b7e56567e7`
(tree `44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`) plus separately sealed
test-only producer source. The S2 terminal-handshake fixture has two principals
in one retained SQLite Store, request-authority v2, the selected BA service,
distinct roots b and c, and the existing controlled exact-b revocation program.
No production runtime, dependencies, automatic retry, or authority semantics
change. One Store-owner loss interrupts both principals; this is not separate
distributed ownership.

Comparing a new zero execution directly to the historical applied execution
would require a substantial typed renaming of namespaces, URLs, program hashes,
root IDs and recursively derived digests. Comparing only the finite model's
small projection would miss other retained differences. Instead, prepare one
closed, pre-execution Store and copy its complete SQLite file set into isolated
branch paths before either branch executes. Run branches sequentially against
separate fresh endpoint stores/processes at the same two logical URLs. Record
and verify the common seed hashes. This is a forked counterfactual fixture,
not two independently provisioned deployments or concurrent cloned ownership.
Never operate both Store copies against one live endpoint state.

## Fixed protocol

1. Start isolated endpoint A and B processes, record ready identity and URLs,
   provision and close the common v2 seed, and seal it. Copy the full closed
   seed file set to the two branch paths before either owner starts. Require
   both copies to equal the seed bytes before execution.
2. In each branch open once in an owner child. Preserve its post-open capture.
   The selected worker records and fsyncs `dispatch.admitted/pre-call`, after
   the durable allow/authorization/debit, then blocks on a separate inherited
   release pipe. The same hook executes in both branches before any call.
3. The scheduler reaches `runner.waiting/in-flight`; its observer first waits
   for the selected pre-call notification without holding the event-log mutex,
   then enters the existing acknowledged scheduler barrier. Its sole in-flight
   authorization must match the selected worker. Pass `external.authorized`
   through the existing observer so that in-flight tracking remains accurate.
   Keep the scheduler barrier unreleased until confirmed SIGKILL in both runs.
   Its owner lifetime must exclude deadline/cancellation/EOF release before
   death; retain abort evidence and classify early return as apparatus failure.
4. Zero: never release the worker before confirmed owner death. Applied:
   release only the worker after the scheduler acknowledgment; await the exact
   endpoint B `resource.applied/post-commit` acknowledgment, with successful
   response withheld, then kill and wait for confirmed owner death. Record
   protocol actions durably. Release endpoint B only after death. Context
   cancellation/pipe EOF must fail closed, never trigger a call.
5. Capture stopped Store bytes and exported observations without reopening the
   original owner. Open once through the public runtime, capture renewed state,
   execute exact-b revocation and one fresh admission consideration, then
   close. The historical allow/debit survives; the fresh gate is refused; c
   survives; no authentic selected result, effect settlement or consumption is
   retained. Expected stopped run coordinates are A sequence 4 and B sequence
   5, both at epoch 2; reopened epoch 3; final B sequence 8 and A sequence 4.
6. Obtain complete endpoint snapshots/journals, terminate and wait for all
   endpoint processes, archive their closed database file sets and final local
   captures. Before the second branch, bind new isolated endpoints to exactly
   the original URLs. Rebinding failure is an apparatus failure, not permission
   to substitute identities. No surviving owner or transport may cross runs.

The normal and race builds each execute one prescribed pair, zero first. These
are two bounded repetitions, not a statistical study. Development probes are
labeled separately. Freeze producer/comparison sources and contract hashes
before the qualification pairs. Preserve every attempted qualification run,
including mismatch or apparatus failure; never retry until a desired match or
discard a failed history. Any necessary protocol correction creates a new
version and retains the rejected attempt.

## Declared local observation and falsifiers

For stopped, reopened and final cuts compare the complete `inspect.json`
structure: namespace, key version, Settings, Roots, every table, row, column,
descriptor and nested value. JSON object key order/whitespace are immaterial;
array order, multiplicity, numbers, strings and null/empty distinctions remain
significant. Reject duplicate keys and malformed schema. No identity renaming,
field deletion, table filtering or post-hoc reordering is permitted. Report
byte equality separately from structural equality. Also compare complete final
`run.json` reports. Seed and post-open equality are prerequisites.

Derive and compare the existing Phase 3B.6 `StoreView` separately, using its
published projection and relative activation mapping. A projected match cannot
substitute for complete-export equality. A differing PC, projection, history,
fact order, pending delivery or run report falsifies the richer match. Any B
application in the zero branch falsifies the proposed world contrast. Neither
failure is repaired by weakening this observation after capture.

Process PIDs/tokens, paths, barrier/control journals and endpoint evidence are
provenance/causal evidence outside this declared local observation. They must
distinguish the histories. The claim is not equality of the entire forensic
dossier, physical SQLite layout, timing, or every possible observer. Endpoint
evidence cannot supply an authentic missing result to the reopened Store.

## Custody and independent checking

The producer writes raw captures as they occur, then seals regular members with
lengths and SHA-256 and archives the complete dossier. No later pass rewrites
raw artifacts to create equality. Derived reports live separately and identify
their exact inputs. Capturing via a disposable complete copy is allowed;
opening the archived evidence SQLite files is prohibited. The independent
checker reads JSON/JSONL and hashes database files only as opaque bytes. The
frozen capture procedure supplies the database-to-export connection; the
JSON-only checker does not independently revalidate SQLite contents.

Separate verdicts cover seal/provenance validity, local raw consistency,
protocol/cut validity, full equality at each cut, report equality, model
projection equality and opposite world facts. Positively certify zero from a
held pre-call through confirmed death and complete empty process-bound endpoint
coverage; missing journals or snapshots do not mean zero. Certify one from the
matched post-commit acknowledgment, complete journal, exact binding and owner
death before successful receipt release. Check original allow/debit, complete
request/root/grant/wire/authorization identity, exact two-principal renewal,
exact-b revocation, current-epoch refusal and absent selected local custody.
Old checkers' pinned producer protocols remain unchanged.

Synthetic hostile controls live in separate scratch fixtures, never genuine
archives. At minimum reject rehashed wrong/missing pre-call or scheduler
acknowledgment, premature worker release, wrong kill/wait, incomplete endpoint
coverage, a second contact, same-applied input on both sides, altered durable
state hidden by an equal small projection, non-exact renewal, stale refusal,
wrong root, and invented local receipt. An independently implemented comparison
must expose paths of differences and must not trust producer verdict labels.

## Deliverable and limits

Deliver a reviewed draft research PR with source/test overlay, reproducible
capture instructions, sealed genuine pairs or preserved failures, independent
checks and adversarial controls, and synchronized manuscript implications/PDF.
No merge or release. No new literature claim or bibliography change is needed
unless a source is actually introduced and verified.

Even a successful pair establishes a source-conditioned experimental witness
for this forked single-host fixture and this declared observation. It does not
establish a general distributed impossibility theorem, unrestricted recovery
equivalence, arbitrary revocation correctness, liveness or N-agent closure.
