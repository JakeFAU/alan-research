# Task 3 review — 331186c7b44cdc18b9f7a0970992fe00ce89714b

## Spec compliance

- ❌ One evidence-output defect needs correction. The core scoped correspondence implements the requested semantic checks, but the deterministic comparison summary misreports the number of unresolved opportunities (`docs/foundational/phase3b5/evidence/task3/compare_reports.py:55`; finding I1 below).
- ⚠️ This is the Task 3 gate only. The controller-owned full ordinary/race qualification, original48 regression, manuscript synchronization, PDF verification, and whole-branch review remain outside this diff. The task report explicitly reserves those checks for the controller; they are not treated as completed by this review.

## Strengths and semantic checks

- The directory API retains the old report, shares its checked input bytes, and separates raw binding, physical validity, continuity, and model support (`research/delegationcorrespondence/binding_files.go:14`, `research/delegationcorrespondence/revocation.go:104`). The legacy refactor changes return plumbing rather than validation policy. The comparison script checks old/current legacy bytes and embedded legacy JSON, as well as repeat opt-in output (`docs/foundational/phase3b5/evidence/task3/compare_reports.py:34`).
- Complete authority roots are checked against B's actual permission tuples and two distinct endpoint identities; the BA/b and AB/c map remains injective (`research/delegationcorrespondence/revocation.go:568`). Accepted lifecycle controls retain their actual source and target-state coordinates, and the accepted revoke must name the exact BA root (`research/delegationcorrespondence/revocation.go:659`, `:669`).
- First authorization is mapped at the reservation prefix, with exact selector, request-selection, grant sequence, serving root, and immutable effect-charge joins (`research/delegationcorrespondence/revocation.go:689`, `:723`, `:743`). The fresh gate separately checks source, attempt digest, ordinal, and present disposition (`research/delegationcorrespondence/revocation.go:489`). Genuine N06 and the reordered-revocation/forged-gate attacks distinguish historical allowance from fresh admission (`research/delegationcorrespondence/revocation_test.go:98`, `research/delegationcorrespondence/revocation_attacks_test.go:104`).
- Both ready-bound journals and the complete snapshot are mandatory, with ordered requests/application coverage and duplicate committed-contact rejection (`research/delegationcorrespondence/revocation.go:227`, `:263`). The existing external checker supplies actual wire, endpoint, receipt, and allowed-gate joins; the new mapper requires a unique physical application before emitting Apply (`research/delegationcorrespondence/revocation.go:508`).
- Every projected event is selected from exactly one enabled operational edge and the complete trace is independently audited; custody and paid status are compared afterward (`research/delegationcorrespondence/revocation.go:398`, `:526`). Intermediate N07 cuts retain a single permit without Apply in the Go value, and the genuine unknown/loss dossier produces no events (`research/delegationcorrespondence/revocation_test.go:79`, `:118`).
- The data attacks exercise semantic predicates after rehashing, rather than stopping at digest failure. The wrong-target attack preserves raw validity and reaches the exact-root check; its compiled mutant falsely accepts the cut. The gate mutant becomes raw-valid while model support still refuses it, correctly preserving the distinction between layers (`research/delegationcorrespondence/revocation_attacks_test.go:239`, `docs/foundational/phase3b5/evidence/task3/qualify_mutants.py:16`, `:24`).
- Claims remain bounded to declared two-principal histories and explicitly exclude whole-Store recovery, unknown results, repeated exposures, reverse simulation, and general progress/distributed claims (`docs/foundational/phase3b5/04-correspondence.md:33`).

## Scoped choices

- **Ready records mandatory: accepted.** They provide the endpoint URL and process/token binding needed to distinguish the two complete permission tuples. This is a justified narrowing met by the sealed cases, documented at `docs/foundational/phase3b5/04-correspondence.md:7` and implemented at `research/delegationcorrespondence/revocation.go:227`.
- **Excluded dossier retains capture locators and zero events: accepted.** Whole-Store/unknown histories have no established relation here. Retaining path/hash/final-cut locators while suppressing a synthetic trace is truthful and retains the independent legacy/raw/physical verdicts (`research/delegationcorrespondence/revocation.go:169`).
- **Wrong-target mutant uses a valid intermediate cut: accepted.** The test first proves independent raw acceptance and proves the genuine BA allowance/gate were mapped before rejection. It isolates the target-root obligation without claiming that an unmodified later AB service would remain valid after a c-root revocation (`research/delegationcorrespondence/revocation_attacks_test.go:239`). This qualifies the mapper predicate, not an altered complete runtime campaign.

## Issues

### Critical

- None found in the scoped review.

### Important

- **I1 — The comparison summary interprets a base64 string as a list of directions.** `RevocationCapture.Unresolved` is `[]b.Obligation` (`research/delegationcorrespondence/revocation.go:41`), whose element type is `uint8` (`research/mutualobligations/model.go:37`). JSON therefore encodes the single BA entry as `"AQ=="`. The comprehension at `docs/foundational/phase3b5/evidence/task3/compare_reports.py:55` iterates four characters, compares each string to integer zero, and reports four BA entries. The retained artifact contains those incorrect multiplicities at `docs/foundational/phase3b5/evidence/task3/verdicts.json:123`, `:135`, `:366`, and `:378`. A single focused invocation of the retained N07 CLI confirmed `Unresolved: 'AQ=='` at both intermediate cuts, exit 0, empty stderr. The model trace itself still has one opportunity, so this does not invalidate the forward trace or N06/N07 support verdicts; it does make the archived evidence summary inaccurate and would mislabel an AB byte as BA too. Correct the new report's JSON representation or decode/validate its actual representation in the comparison script. Add a serialized-output assertion that each N07 intermediate cut contains exactly BA once and the final cut contains none; regenerate the summary and affected manifests while preserving the erroneous artifact as historical failed evidence. Avoid changing existing model State/Event JSON shapes to fix this new report field.

### Minor

- No additional findings.

## Evidence checked and limits

- Read the complete supplied `review-47d4e4a9..331186c7.diff` once in sequential chunks, plus the brief, context, global constraints, task report, contract, and the named reconnaissance note. No changed source file was separately reread.
- Inspected retained `green.log.gz`, `reproduced.log.gz`, `legacy-attacks.log.gz`, both compiled-mutant assertion logs, and `staticcheck.log.gz`. The positive/focused logs report PASS without warnings; mutant failures occur at the intended semantic assertions, not compilation failure or panic; Staticcheck output is empty. Archive/determinism claims were assessed through the checked-in scripts, manifests, and recorded verdicts, not a fresh archive/campaign rerun.
- Named outside-diff risk **endpoint joins might not establish receipt/contact correspondence**: checked `research/delegationcorrespondence/external.go:31–210`. It validates ingress against exact attempt wires and ready URLs, committed receipt content, retained receipts, and allowed-gate counts. The new complete-snapshot checks strengthen those existing prerequisites.
- Named outside-diff risk **unchecked indexing or custody ordering could invalidate the mapper's fixed arrays/table access**: checked `research/delegationcorrespondence/binding.go:26–109` and `research/delegationcorrespondence/check.go:245–270`. Valid declared binding yields exactly two originals in declared order, while the raw checker enforces singleton composition permissions and matches actual permission rows.
- Named outside-diff risk **unresolved JSON cardinality differs from its Go value**: checked `research/mutualobligations/model.go:37` and the absence of a custom obligation marshaler, then ran only `/tmp/alan-phase3b5/task3/revocation -case /tmp/alan-phase3b5/corpus/raw-normal/N07-3888622734` through a small Python JSON inspection. This was the sole executable probe; no suites, race runs, runtime campaigns, SQLite opens, branch/index changes, or source edits were performed.

## Assessment

**Task quality: Needs fixes.** The semantic projection and attacks are cohesive and appropriately bounded. Correct I1 before accepting the retained Task 3 evidence; no other blocking issue was found in this scoped review.
