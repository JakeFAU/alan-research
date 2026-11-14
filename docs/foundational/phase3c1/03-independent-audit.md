# Phase 3C.1 independent evidence audit

## Question and result

Does the fixed exported evidence preserve every named progress register, service
dependency and concrete receipt obligation through grouping, and do its reported
progress verdicts follow from the declared semantics? The independent Python
audit accepts the genuine Task 1 export and rejects 35 separately retained hostile
inputs. This establishes agreement for the complete small abstract graphs and
seven selected concrete histories. It does not establish whole-product concrete
exhaustiveness, a machine-checked lifting proof or runtime/distributed refinement.

The two freshly generated exports are byte-identical to each other and to the
immutable Task 1 sample: 1,584,501 raw bytes, SHA-256
`48b4cfb4cc10b7dde4979874f669014f4863e1e05e48737ef4561f119dfc3602`.
The retained full export is `evidence/task2/export.json.gz`; gzip metadata uses
mtime zero and no source filename. Each hostile input has its own compressed
file and rejection log under `evidence/task2/hostile/`. None replaces or edits a
genuine record.

## Independent method and boundary

`tools/audit.py` uses Python 3.9.6 and only the standard library. It neither imports
nor executes Go. Export regeneration is a separate capture command. The author
read the frozen contract/schema and genuine exported data, not the new producer
implementation. The controller supplied a declarative clarification of the
unchanged imported durable-ACK protocol after consulting its old source; that
clarification is retained as `evidence/task2/protocol-interface.md` and hashed in
the manifest. In particular it specifies occupied-channel guards, repeatable
ACK release, post-consumption receipt delivery and unchanged volatile custody
on recording/consumption. This shared protocol specification is an explicit
independence boundary: separate implementations do not prove the specification
itself correct.

The audit builds the finite register lattice independently, selects reachable
states from U, reconstructs all named transitions and computes a predecessor
fixed point from Complete. It checks exact node state sets and every outgoing
arc, thus also every incoming arc, rather than trusting the exported predecessor
tree. It separately checks legal BFS parent/action/depth witnesses, Complete,
CanComplete, all graph totals and the specified prerequisite-first dependency
ordering. A non-tree incoming-edge deletion is a retained negative control.

Leaf construction and compatible union retain each physical owner, open endpoint,
full binding, resource and outgoing register. Both parenthesizations are built
again and their closed state/arc relations are compared with the flat relation.
The checker rejects missing internal debt, erased ports, policy/generation
conflicts and owner/resource aliases. The fixed fixture roster and full bindings
are checked exactly; these are checks for the declared chain/cycle grammar,
not an unrestricted component-composition API qualification.

Concrete replay starts from independently constructed empty custody and physical
availability. A declarative protocol table determines legal events and effects;
a per-original causal ledger separately records actual application, caller
recording and ACK origin. Actors obtain result bytes only from retained data,
channels or their recorded/volatile receipts. The historical `WorldResult` is
written at application and compared as part of the reported post-state; it
never supplies bytes or enables an actor. All roles consult the same physical
owner record. Every reported post-state and progress projection is compared
with independently replayed state, including bindings when result bytes coincide.
The auditor recomputes concrete CA readiness for R1 with availability, admission
and the named outgoing dependency, rather than confusing macro eligibility
with immediate concrete readiness.

## Executed outcomes

| Fixture class | Profiles | States | Arcs | Can complete | Cannot complete |
| --- | ---: | ---: | ---: | ---: | ---: |
| Chain | 8 | 128 | 184 | 128 | 0 |
| Cycle | 8 | 511 | 1,053 | 510 | 1 |
| Total | 16 | 639 | 1,237 | 638 | 1 |

Thirty-two reconstructed grouped graphs have the same full relations as their
respective flat graphs. All dependency certificates match; every acyclic
profile is completable at every reachable abstract state. The cyclic profile
is not classified as universally blocked: 62 of its 63 states can complete.
Only its W/W/W state cannot. This is failure to discharge circular service
assumptions, with valid identity and custody, not an authority/custody failure.

All seven witnesses, containing 54 microsteps, replay successfully:

- R0 has equal full BC transport and physical availability, with hidden AB at W
  versus D and global completion false versus true. Both can complete. The
  continuation adds the hidden AB application and consumption visibly.
- R1 has equal full BC/CA transport and physical availability, with hidden AB at
  U versus W. Both carry unpaid AB debt, but concrete CA application readiness
  changes from enabled to disabled. This is a readiness distinction, not a
  refutation of every weaker trace relation.
- The all-serial control reaches W/W/W with no abstract successor/completion
  suffix, while obeying concrete custody and identity rules.
- Retry states 3 and 5 are exactly equal, separated by legal request send and
  deduplication. This supplies one repeatable concrete divergence path.
- The pre-record-loss history loses the actual caller volatile receipt, recovers,
  retransmits from callee custody, records, sends ACK, releases and consumes.
  No world-oracle reconstruction or repeated application is used.

The eight unittest methods pass, including 35 hostile full-export controls and
additional direct ledger/join checks and CLI probes for malformed JSON,
duplicate object keys and NaN in an integer field. The 35 CLI invocations all exit 1 with the
specific expected rejection reason; the genuine CLI exits 0. Reasons and raw
input digests are in `evidence/task2/hostile-rejections.json`. Controls cover schema
shape/types/version, missing profiles, forged bindings, missing/extra/non-tree
incoming arcs, wrong serial guard/completion/progress/count/BFS/dependency,
erased hidden debt/ports/owners, forged grouped closure, invented initial or
post-state facts, missing visible projection, unavailable actors, premature
release, repeated application, fabricated receipt bytes, malformed collisions,
wrong retry endpoints, erased loss history and overstated evidence category.

The source argument's universal representative lifting obligation remains a
source argument. Task 1's factored Go checks are separate evidence, not rerun
by this auditor and not reclassified as whole-product exploration. Missing
recoverability, split application/retention, erased custody, changed bindings
or shared-resource depletion violate the declared premises; the present audit
does not make a preservation claim for those environments. Claim IDs/statuses
and nonempty statements/evidence/limits are checked structurally. The actual
prose and its source support require review; arbitrary prose truth is not
certified by a JSON parser.

## Reproduction and retained process evidence

From the repository root:

```sh
mkdir -p /tmp/alan-phase3c1/task2
python3 -m unittest discover -s docs/foundational/phase3c1/tools -p 'test_*.py'
GOCACHE=/tmp/alan-phase3c1/task1/go-cache go run ./research/contractcomposition/cmd/export > /tmp/alan-phase3c1/task2/fresh.json
python3 docs/foundational/phase3c1/tools/audit.py /tmp/alan-phase3c1/task2/fresh.json
python3 docs/foundational/phase3c1/evidence/task2/capture.py
```

The capture script records actual argv, environment overrides, exit codes,
outputs, tool versions, raw/compressed artifact digests and source hashes in
`evidence/task2/manifest.json`. It reruns two exports, requires byte equality
with the frozen sample, audits the genuine record, runs unittests and checks
every separate hostile record through the CLI. Its Go cache is the completed
Task 1 cache; scratch JSON stays under `/tmp/alan-phase3c1/`. The manifest does
not hash itself. Source hashes include the contract, schema, argument, producer,
old imported protocol sources, independent checker/tests and protocol interface.

The first test run preceded validator implementation and produced 36 assertion
failures against an accept-all stub (`red-initial.log`). Its inner unittest
exit was not separately captured by that initial shell wrapper; the failure
output is retained without inventing an exit record. The next run had one
expected-reason mismatch: deleting recovery also left a forged post-state,
which correctly failed before the intended down-actor event. That hostile
fixture was corrected to preserve the down state, so the availability guard
is now tested directly (`green-attempt-1.log`). A later test-first control
exposed that an entirely erased loss/release history was still accepted;
`red-loss-witness.log` records those two failing assertions, exit 1. The auditor
then added explicit evidence-presence checks for that named witness. The final
`unittest.log` records all eight test methods passing. These failed candidates
remain historical process evidence, not current failures or fabricated producer
defects.

No producer/schema defects were found by this audit. No runtime or old-model
source, old evidence, manuscript, dependencies or retained SQLite was modified;
no runtime `.envrc` was accessed. Manuscript reconciliation, root qualification
and independent task/whole-branch review remain controller obligations.
