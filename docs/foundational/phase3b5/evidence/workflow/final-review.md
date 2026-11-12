### Strengths

- The operational change is narrow and opt-in. Separate issuance, immutable first decisions, current admission checks, surviving pre-revocation permits, crash-cleared permits, and distinct b/c roots agree with the contract.
- The ledger independently interprets authority and custody descriptors. Reverse reachability and SCC analysis remain distinct from the operational explorer; the documentation accurately identifies their shared transition-graph input.
- The concrete projector checks complete roots, lifecycle targets, historical decisions, fresh gates, principal selection, debit identity, and physical evidence. Intermediate N07 captures retain an opportunity without inventing Apply.
- Adverse evidence is substantive: compiled mutants fail targeted assertions, universal-completion counterexamples remain retained, and the corrected unresolved-opportunity serialization has direct JSON-level assertions.
- Manuscript changes distinguish finite enumeration, source-level proof sketches, and conditional historical trace checks. They preserve the unknown-after-loss exclusion and avoid claiming runtime atomicity, reverse simulation, or distributed closure.

### Review coverage and evidence

Reviewed `4fea6c31eb8f438c0dea414aa5f8d7d0ecb5e8a4..2b7be7d6612e8770e6638e5026e3ff0b246052f2` in passes: contract/plan/rulings; changed Go implementations and focused tests; concrete relation and supporting physical-evidence checks; manuscript/documentation changes; qualification, compatibility, mutant, and artifact evidence.

Read-only verification during this review established:

- All 582 frozen executable-source hashes match the checkout.
- All ten retained qualification streams match their recorded compressed/uncompressed hashes and sizes.
- Parsing the actual test streams confirms 2,633 ordinary test/subtest passes, three named optional-export skips, and no failures; affected race evidence contains 422 passes, no named skips, and no failures.
- Retained model/audit/progress outputs agree on 229,834 states and 1,807,246 arcs, with 203,833 safe-completable states and zero predicate mismatches.
- All four historical archive digests match the recorded fresh export digests. Final retained new audit/progress bytes equal Task 2's outputs.
- All 47 Task 4 manifest entries match. PDF archive members and installed manuscript/PDF hashes match the recorded build inputs and output.
- Actual N06/N07 normal-mode capture facts confirm the stated authorization/revocation/gate/result ordering. Their declared capture hashes match; endpoint snapshots show the reported zero and two total applications respectively.

Outside the diff, I inspected the existing ledger/custody logic and `CheckExternal` joins to check integration assumptions, plus the selected N06/N07 raw JSON for chronology and physical-evidence claims. I did not open SQLite, rerun suites, enumerate graphs, rebuild, or alter checkout/index/branch state. PDF visual findings remain attributed to the retained controller inspection.

### Issues

Critical: none found.

Important: none found.

Minor: none found.

### Declined to judge

- General distributed revocation, arbitrary agent counts, regrant/multihop, and permission-algebra equivalence: excluded from the frozen model and claimed relation.
- Whole-Store loss as modeled A/B crashes, unknown/applied-unrecorded effects, repeated exposures, and runtime application/receipt atomicity: expressly excluded; the genuine loss case remains unsupported.
- Cryptographic authenticity or resistance to coherent rewriting of unsigned exports: manifests establish byte identity, not producer truth.
- Repair of inherited global lint diagnostics and the pre-existing Markdown-only manuscript section: disclosed baseline conditions outside this cohesive change.
- A new runtime campaign: none is claimed or required for this historical-corpus analysis.

### Recommendations

Complete the controller-owned archival of this review and final workflow status, then proceed with the planned draft research PR. No implementation correction is requested.

### Assessment

**Ready to merge? Yes, on the reviewed technical scope.**

The implementation, evidence, and claims are consistent with the bounded contract, with no outstanding review finding. Delivery remains the authorized draft-PR workflow; this assessment does not constitute merge or release authorization.
