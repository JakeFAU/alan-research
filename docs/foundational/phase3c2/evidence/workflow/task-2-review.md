# Task 2 independent specification and quality review

Reviewed base `60d46a43dd32caa07f0b814e6741213f506a9cb6` through head
`a24f8ff1c6aac22a0cffcffc2c614ccfca819c71`. HEAD matched the requested head.

**Specification: ✅ PASS. Quality: APPROVE.**

This verdict covers Task 2's independent Python audit, checker tests, audit
documentation and retained evidence. It does not approve Task 3, runtime
qualification, an unrestricted composition theorem, or integration/merge.

## Question and result

Does Task 2 independently reconstruct the complete frozen abstract relation and
replay the actual exported histories, while preserving strict schema, source
provenance and the distinction between independent checks and reported probes?

Yes, within the declared finite chain and evidence boundary. No material
specification disagreement or actionable correctness finding was identified.

## Specification assessment

- ✅ All 72 configurations are constructed independently from the three layouts,
  three reservation policies and eight service masks. Exact named bindings,
  generation-one owner identities, roots, limits and memberships are retained.
- ✅ BFS reconstruction includes every enabled arc, deterministic parents and
  action labels. Reverse adjacency computes existential finite completion;
  exported labels, counts and parent trees are not used as the oracle.
- ✅ Static component records are reconstructed from their owner sets and exact
  outgoing memberships. Both parenthesizations and closed configurations are
  checked structurally against the same fixed configuration. This is the
  specified closed roster check, not a general component substitution theorem.
- ✅ Concrete replay begins with empty custody and all initial owners. Each
  event checks actor, availability, full binding, causal transport state and
  capacity/service conditions; complete generated post-states are compared
  recursively. Application alone creates the actual outcome and retained
  receipt, and clears the exact capacity hold atomically. Crash cannot refund a
  hold; receipt release requires a causally created durable ACK.
- ✅ The checker does not read WorldResult to enable actions or manufacture
  custody, and does not invoke Go semantics or runtime Replay. Its producer-file
  reads are hashes, not a semantic oracle.
- ✅ Holder pairs compare full Base records and exact opposite holders under
  the fixed shared root. Completed histories extend the specified prefixes;
  ready-now invalidation and crash preservation are checked. Stable-ready
  membership of the old erased trace is tested with reserve-only epsilon
  closure, so visible BC progress cannot be hidden as preparation.
- ✅ Recovery witnesses check both actual result values, complete paired-history
  correspondence, duplicate deliveries, both pre-record losses, C's retained
  BC hold, durable custody and release-before-consumption ordering.
- ✅ Strict recursive equality rejects unknown or missing fields, wrong types,
  array order/length errors and Boolean-as-integer substitutions. JSON decoding
  rejects duplicate keys and nonfinite constants. All named hostile classes are
  represented by retained exact forged inputs and reason-specific rejections.
- ✅ Operational probes are explicitly source-bound reports with structural
  roster/type/total checking. The checker output explicitly denies independent
  probe re-enumeration. The 144 Python-generated histories are correctly labeled
  checker fixtures sharing Python transition rules, not additional independent
  semantic evidence.
- ✅ The commit contains only the 198 new Task 2 files. Production, dependencies,
  old evidence and manuscript files are outside its diff. Historical RED logs
  are retained separately from final success. Manuscript work remains Task 3.

## Strengths

The full-relation comparison is substantive: the retained non-tree mutation
removes an arc while adjusting its count and keeping every predecessor arc.
This distinguishes graph reconstruction from tree walking or count checking.

Replay constructs the next state from causal events before comparing the
exported post-state. This prevents a plausible final receipt or completion label
from substituting for a real application, transport and durable-record history.

The audit documentation is precise about independence. Separate algorithms over
a shared frozen specification, source-hashed operational reports and generated
checker fixtures are each given the correct evidentiary status. The legal but
non-corresponding history and custody-order controls also test witness meaning
beyond mere transition legality.

## Findings

**Critical:** none.

**Important:** none.

**Minor:** none requiring a change.

## Review method and focused verification

Read the supplied brief/report, current AGENTS.md, frozen contract and schema,
and the entire supplied 4,684-line diff in bounded, nonoverlapping chunks. The
initial combined input read was truncated; the contract and schema were then
read completely. Changed implementation code was read through the diff once.

One targeted outside-diff lookup checked the inherited transport contract:
whether retries after application/consumption and serviceable receipt handling
agree with the Python replay rules. The Phase 3C.1 schema/contract and Phase 3A
durable-ACK contract support those rules. No Go implementation was read as a
semantic oracle.

Executed `git status --short` and `git rev-parse HEAD`. The existing modified
plan and untracked preflight evidence were left untouched.

Executed one read-only standard-library `python3 -B` archive-integrity script.
It did not run the checker or the test suites. It independently checked:

- all 201 manifest paths, byte lengths and SHA-256 values, including gzip raw
  lengths/hashes and unique manifest paths;
- coverage of all 198 changed files by the manifest, except the manifest itself;
- byte equality of the retained Task 1 and Task 2 exports;
- all 32 producer source hashes against current bytes;
- all 43 hostile input hashes, command exits, stderr linkage and empty stdout;
- the non-tree control's 24 preserved states, corrected 33-arc total versus the
  original 34, and all 23 preserved BFS predecessor arcs.

Fresh script result (exit 0):

```json
{"archive_entries_hash_verified":201,"changed_files_covered":198,"hostile_capture_linkages_verified":43,"non_tree_control_preserves_all_23_parent_arcs":true,"producer_hashes_verified":32,"task1_task2_export_bytes_equal":true}
```

The retained final unittest evidence records five passing tests, including 43
hostile subcases and 144 generated-history subcases. Retained audit outputs
report 72 graphs, 1,632 states, 2,412 arcs, 1,624 CanComplete labels, eight
CannotComplete labels, 72 groupings and ten histories comprising 161 steps.
These are inspected implementation captures, not freshly rerun reviewer tests.

## Declined-to-verify items and remaining obligations

Per the review brief, no already-passing suite, Go exporter, full repository
test, race test or Staticcheck campaign was rerun. No independent re-enumeration
of the factored operational-probe domain, exhaustive concrete-product search,
mechanized proof or Go lifting/noninterference verification was performed.
The documentation already disclaims these stronger conclusions.

Task 3 still owns manuscript/roadmap/PDF reconciliation and broader
qualification. This review does not establish fairness, split-effect recovery,
distributed allocation, arbitrary topology or old-trace equivalence. No retained
SQLite database or environment configuration was opened. Only this review
report was written; no source, evidence, dependency, index or branch was changed.
