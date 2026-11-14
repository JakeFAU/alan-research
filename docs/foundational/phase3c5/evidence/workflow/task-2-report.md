# Task 2 report

Status: DONE (independent review remains with controller).

## Scope and source freeze

Implemented the independent Python auditor, ten-test qualification suite, separate
fresh-directory capture helper, independent-audit account and complete Task 2
archives. Producer/contract/schema/manuscript files were not modified. Task 1's
capture helper was reused unchanged, with its exact hash recorded.

- Auditor SHA-256: `8944f418c39f21db4b4ad2c38b70b5efe039957d8abbe94cff79f6dc8bbd4235`
- Tests SHA-256: `6bdeaf38227e06a34c2679a40e70ce775a250e6020945d2e6fbc2cacc22c874a`
- Task 2 capture helper SHA-256: `eff7f2630ff28071bd69e0fa8472fb5ccae457e255a3f9ff6d4662934d269f33`
- Full reuse/source identities: `docs/foundational/phase3c5/evidence/task2/qualification-1/checker-sources.json`.

No producer predicates or prior verdicts are oracles. New inclusion/reduction/
coverage/hit logic is independent. Unchanged C4 automaton/product/path utilities
and C3/C2 descriptor/strict-parser utilities are explicitly pinned before import;
the account explains their precise boundaries. All seven inherited monitors,
72 plants and 324 controllers are independently reconstructed from those Python
constructors before retaining inherited serialization.

## Actual verification

1. `python3 docs/foundational/phase3c5/tools/audit.py docs/foundational/phase3c5/evidence/task1/capture-3/export.json.gz`: accepted after the preserved initial auditor correction.
2. `python3 -m unittest discover -s docs/foundational/phase3c5/tools -p 'test_*.py' -v`: ten tests pass. Final captured run 32.094 seconds. Includes all 256 families / 65,536 ordered pairs over a three-word universe for the finite algebra laws.
3. `python3 docs/foundational/phase3c5/tools/capture_audit.py docs/foundational/phase3c5/evidence/task2/qualification-1`: exit 0; fresh producer capture, ten tests, two valid CLI inputs and 56 exact hostile CLI inputs. All expected outcomes. No initially accepted hostile inputs.
4. All recorded archive/input hashes and source hashes verified. Nested capture manifests verified.
5. Repeating the capture command exits 1 with `refusing existing capture directory`; all 311 files remain unchanged, recorded in `fresh-directory-refusal.json`.

Fresh export raw identity is exactly capture3/capture4's
`ab3cbdfda577a3ae804eb76a94a8c7ac665cf00c69f992ad6804b6e0d0e27b8a`,
14,097,260 bytes. Producer source bindings remain unchanged. The fresh producer
capture includes module-relative go-list source roster equality through the
unchanged Task 1 helper. No Go suites, root checks, old campaigns or SQLite runs
were performed by Task 2.

## Measurements and preserved negative outcomes

3,564 products; 45,800 states; 59,492 arcs; 22,716 CanComplete states; 1,234
initially completable products; 760 nonblocking products. 36 comparisons;
1,442 residual comparisons; 5,208 language graph nodes; 146 profile graph nodes.
Equal initial/FULL/MIN/MAX/J pairs: 13/6/7/13/7. All 2,772 J substitution rows
retained: 844 nonblocking, 1,928 blocking.

81 local failed targets; 28 H empty; largest H has 3 words. No singleton-epsilon H
occurs in the roster, so dedicated tests establish H empty versus H containing
epsilon separately. Both new monitors retain 52 nonblocking/20 blocking plant
contexts and 324 blocking plant/controller contexts each. RF/doomed MAX failure,
MIN initial-language loss, J nonnecessity, local-hit-not-NB and legal-trace
controls are independently rebuilt. Local H is not described as a global NB
context.

## Preserved failure

The initial auditor mistakenly kept `Use` in `Input.Sources`, contrary to the
accepted Path/SHA256-only schema. It rejected valid input with
`REJECT: export.Input.Sources[0]: keys mismatch`. Original auditor bytes and
exact result are retained in `evidence/task2/red-initial/`. Correction changed
only the auditor serialization. No producer defect found, no evidence changed,
no attack accepted before repair.

## Self-review and limits

Reviewed source changes, staged status/stat and complete staged diff. Generated
capture records and exact binary-input archives were read and checked through
nested manifests, source hash comparisons and lossless input hash recovery.
The commit contains only Task 2 files. The local report is workflow metadata
under ignored `.superpowers/`; it is not forced into the repository.

The checker remains finite qualification under the frozen semantic assumptions;
its product implementation reuses independently written predecessor Python.
No unrestricted theorem, second implementation of those Python utilities,
distributed/split-effect guarantee, or whole-agent extraction is established.
Controller owns independent review, manuscript reconciliation and delivery.

Primary account: `docs/foundational/phase3c5/03-independent-audit.md`.
Evidence root: `docs/foundational/phase3c5/evidence/task2/`.

Commit: `b079c67db3298688f32f2c5c8218a01f81d45601` (`Audit finite residual antichains independently`). Post-commit working tree is clean.

# Repair round 1 — executable source binding

The initial independent review returned **REJECT**, with one Important/P2
finding. Its successful stale-bytecode substitution disproved the original
source-to-execution guarantee. The earlier report's “no initially accepted
hostile inputs” applies only to the 56 JSON/wire inputs; it does not cover the
separate successful cached-code provenance attack. There is no evidence that
qualification-1 actually used stale code or that its finite conclusions were
wrong. Its executable identity was insufficiently guaranteed.

## Reproduced failure and preservation

Before editing, ran:

```sh
python3 -B .superpowers/sdd/phase3c5-plan/task-2-review-cache-repro.py
```

Exit 0; output included `Auditor imported; effective ALPHABET: ('CACHE',)`.
`evidence/task2/fix-1-red/` preserves the exact reviewer script, original command
and output, this replay's command/output, and the pre-fix audit/test/capture
sources. The snapshots were checked byte-for-byte against commit
`b079c67db3298688f32f2c5c8218a01f81d45601`. The original script reads files from
its working tree: replay it against that historical commit, or a temporary
source tree containing the retained C5 auditor and unchanged C4/C3/C2 auditors.
The script itself was not altered to pretend repaired sources reproduce the old
failure. qualification-1 and the original artifact manifest remain unchanged.

## Correction

Added `tools/source_loading.py`. It reads a declared finite set of source
buffers, records their hashes, verifies the three predecessor pins, and compiles
exactly those buffers. Its private importlib facade intercepts the unchanged
C4→C3 and C3→C2 explicit loaders; its local import hook covers declared checker
modules. There is no global importlib monkeypatch or predecessor edit.

The auditor, capture helper and source-only test runner bootstrap the helper
from one read buffer, then re-execute their own recorded source buffers. Local
capture/test imports use that bundle. Captured source records come from those
buffers, not from a later source-file reread. Every child auditor and test runner
emits SOURCE_BUFFERS identities; capture checks equality with its own manifest
before recording a successful outcome. The Python installation and stdlib are
trusted. Arbitrary outside callers choosing their own cache-capable importer
are not authenticated by this API.

Added `tools/run_audit_tests.py` and `test_source_loading.py`. The accepted
qualification command is now:

```sh
python3 docs/foundational/phase3c5/tools/run_audit_tests.py
```

This replaces bare unittest discovery for accepted qualification because the
test modules themselves must be loaded from recorded buffers. The previous
unittest command/result remain historical; no retroactive source-only guarantee
is claimed. The new runner assembles the declared suite from source-loaded
modules. The old standalone C4/C3/C2 entrypoints are unchanged and are not
retroactively repaired or qualified.

The capture helper now reuses the unchanged accepted producer export during
this Python-only repair. No Go exporter, dependency listing, Go suite or root
check was rerun. All producer/inherited source bindings remain checked by the
auditor's independent reconstruction.

## Covering tests and actual qualification

The first source-only development run passed all 13 tests in 34.084s. After
adding child-process identity cross-checks, captured the complete final run:

```sh
python3 docs/foundational/phase3c5/tools/capture_audit.py docs/foundational/phase3c5/evidence/task2/qualification-2
```

Exit 0. The final suite passed all 13 tests in 34.158s; valid plain/gzip inputs
were accepted and all 56 JSON/wire attacks rejected. Nine same-size,
timestamp-valid cache controls cover C4, C3, C2, local audit, test, capture,
source helper, runner and regression-test module. Every control first proves
normal SourceFileLoader executes the poison under `-B`; the source-only loader
then executes the intended source. This is a passing regression with deliberate
successful stock-loader attack controls, not a claim that the attacks fail
under every possible import mechanism.

Two further source tests verify that changing an on-disk file after buffering
cannot change the compiled bytes/recorded identity, and that altered pinned
source is rejected before its code runs. No new test-fixture/preparation errors
or initially accepted JSON/wire attacks occurred in this repair. The prior
Input.Sources serialization error remains separately attributed to the original
auditor's first construction, not conflated with malformed-input rejection.

All 59 final child command source records agree with the nine-entry source
manifest. Their output/input archives and hashes were checked, as were the
historical artifact manifest and pre-fix snapshots. Repeating the corrected
capture command exits 1 with `refusing existing capture directory`; all 298
corrected capture files remain unchanged. Exact records are:

- `evidence/task2/qualification-2/` — final source-only tests and 58 CLI inputs.
- `evidence/task2/fix-1-verification.json` — hash/source/record counts.
- `evidence/task2/fix-1-refusal.json` — refusal command and unchanged files.
- `evidence/task2/fix-1-capture-command.json` — outer command/output/exit.
- `evidence/task2/fix-1-artifact-manifest.json` — historical plus repair archives.

Final source identities:

- Auditor `a2657b3a2212c16fba0f64908a6bf51bfcef716bf4b8c54f777cd553ccf2a86a`
- Capture `140273986f79e70b27a16b4f1b19fbff333645619e5b7847b908852bdf9af677`
- Runner `1e1aa01887f96d9749b3288547cbee808fb0039a197c409e2813dfc1e7ea6279`
- Source helper `69c5788d996b3aa30c0db190cbc3e3cc5372ba25a910687f2e246dad6f3bba1a`
- Semantic tests `15dc11a56e93569959b209cea6c24dee53c5f9d1e7385393b31d22e1f4d485f6`
- Source-loading tests `9d6b52e1d045dc7d3f985aa20452dbd635c0b8ffca78a319d0b73914de600f5c`

## Scope, self-review and status

Source diff review confirms the finite inclusion/reduction/coverage/hit and
product reconstruction algorithms are unchanged. Only executable-loading,
source-identity reporting, local tests/runner, corrected capture and honest
qualification wording changed. No predecessor, Go, root module, manuscript or
old evidence files changed. Historical captures and the successful provenance
attack remain intact. Staged status/stat/full diff and archive identities were
reviewed before committing. Independent re-review remains the controller's gate;
this repair does not claim approval.

Repair commit: `208e1411dd9074f4d97eb9fbba4fb9b8755b3236` (`Bind checker execution to verified Python source buffers`). Post-commit working tree clean. Ready for independent re-review.
