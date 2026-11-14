# Spec compliance: REJECT

Task 2, base `a036b4a6`, head `b079c67db3298688f32f2c5c8218a01f81d45601`.

The finite reconstruction and captured adversarial matrix match the requested scope. One Important finding prevents acceptance of the exact executable-source binding: hashing Python source before importing it does not ensure that those bytes are executed. A focused experiment reproduced execution of different cached code while all three predecessor source pins passed.

# Strengths

- `docs/foundational/phase3c5/tools/audit.py:63`: inclusion uses the correct directional acceptance predicates, retains the complete subset graph, and distinguishes missing witnesses from epsilon. Both directions use alphabet-order BFS witnesses supplied by the pinned subset utility.
- `docs/foundational/phase3c5/tools/audit.py:90`: equivalent residuals are collapsed before strict MIN/MAX filtering; all strict relations are retained. `:118` constructs every target's coverage and finite H, then checks actual source acceptance and target rejection. `:141` compares complete reachable profile graphs, including empty subsets.
- `docs/foundational/phase3c5/tools/audit.py:218`: compressed/raw predecessor identities and current inherited source hashes are checked. `:247` independently reconstructs all inherited descriptors before retaining pinned serialization. No predecessor verdict is used as an oracle.
- `docs/foundational/phase3c5/tools/audit.py:262`: all products, J substitution rows, core outcomes, negative paths and legal-trace controls are rebuilt. `:189` separately reconstructs generic controls and finite suffix tries. The old Python product utility implements synchronized transitions, private tau and reverse marked reachability; it does not call Go or consume old products.
- `docs/foundational/phase3c5/tools/test_audit.py:78`: the semantic attacks cover the requested falsifications, including missing non-tree arcs with repaired counts. `:210` exercises epsilon, empty language versus empty family, equivalent states and tau cycles; `:229` separately checks the finite algebra over all 256 families and 65,536 ordered pairs.
- `docs/foundational/phase3c5/tools/capture_audit.py:65`: exact wire inputs, output streams, hashes, exit codes and diagnostic expectations are retained for both valid formats and every hostile CLI case. `:104` refuses an existing output directory before capture writes. `docs/foundational/phase3c5/03-independent-audit.md:113` accurately distinguishes the preserved initial auditor error from a producer counterexample.

# Issues

## Important / P2 — verified source bytes are not necessarily executed

**Primary location:** `docs/foundational/phase3c5/tools/audit.py:29–36`.

The code hashes each `.py` file and then executes a `SourceFileLoader` through `spec_from_file_location(...).loader.exec_module(...)`. That loader may use a pre-existing `.pyc` whose timestamp and size match the source. `sys.dont_write_bytecode = True`, `PYTHONDONTWRITEBYTECODE=1`, and `-B` prevent writing caches; they do not prevent reading them. Consequently the predecessor source hashes and `checker-sources.json` can accurately describe the disk files while a different automaton, product, strict-comparison or descriptor implementation actually executes. This violates the task's source-bound executable reuse claim, including `docs/foundational/phase3c5/03-independent-audit.md:28–35`.

**Reproduced:** In a temporary tree containing exact copies of the four auditor source files, I compiled a same-length alteration of the C4 module's `ALPHABET` assignment, restored the original source bytes and mtime, asserted its published SHA-256, and imported C5 under `python3 -B`. The import exited zero and reported `ALPHABET: ('CACHE',)`. All predecessor source-pin checks passed. The experiment executed no full reconstruction or test suite and changed no repository source/evidence. This establishes a source-to-execution mismatch, not that the existing capture actually used a stale cache.

**Affected import boundaries:**

- `docs/foundational/phase3c5/tools/audit.py:34–36`: C5 loads C4.
- `docs/foundational/phase3c4/tools/audit.py:27–29`: reused C4 loads C3.
- `docs/foundational/phase3c3/tools/audit.py:20–22`: reused C3 loads C2.
- `docs/foundational/phase3c5/tools/capture_audit.py:15–16`: capture imports local `audit` and `test_audit` before recording source identities at `:42–44`.
- `docs/foundational/phase3c5/tools/test_audit.py:5`: test module imports local `audit`; unittest discovery itself normally imports `test_audit` through the standard cache-capable machinery. The capture's `PYTHONDONTWRITEBYTECODE` override at `capture_audit.py:22` does not close these boundaries.

**Required correction:** Execute the exact verified source buffers, bypassing bytecode caches throughout the predecessor import chain and the local capture/test imports. Make the source identity bind the bytes compiled, rather than hashing and subsequently reopening through a cache-capable loader. Keep predecessor files unchanged as required by the contract; implement source-only loading/isolation in the new task-owned code. A source-only outer import that still permits cached transitive imports is insufficient. Mere `-B`, setting `dont_write_bytecode`, or repeating the existing hashes is insufficient.

**Required regression/evidence:** Add focused temporary-copy controls with timestamp-valid same-size stale bytecode at each inherited boundary and the local capture/test boundaries, asserting that the intended source semantics execute (or that the input is rejected before execution). Retain the reproduced failure. Preserve the current qualification directory; capture corrected checker/tests/helper identities and fresh results in a new directory, with a precise account of what was rerun. Unchanged Go suites need not be repeated for this Python repair. Re-review the actual corrected diff and evidence before accepting Task 2.

**Reproduction artifacts:**

- `.superpowers/sdd/phase3c5-plan/task-2-review-cache-repro.py`: exact successful experiment source, originally supplied to `python3 -B -` on stdin.
- `.superpowers/sdd/phase3c5-plan/task-2-review-cache-repro.command.json`: original invocation, replay argv, cwd, source SHA-256 and exit.
- `.superpowers/sdd/phase3c5-plan/task-2-review-cache-repro.stdout`: exact observed successful experiment output.
- `.superpowers/sdd/phase3c5-plan/task-2-review-cache-repro.stderr`: empty output.

Replay from the repository root with `python3 -B .superpowers/sdd/phase3c5-plan/task-2-review-cache-repro.py`. The script confines temporary sources and caches to its temporary directory. Initial harness attempts encountered the host's default cache write restriction and then a `/var` versus `/private/var` cache-path alias; the retained successful script explicitly isolates the cache and resolves the temporary path.

## Critical

None found.

## Minor

None raised.

# Checks and review boundaries

- Read the complete 9,876-line, 320-file diff in bounded passes. Read all newly introduced executable source and the audit account. Mechanically parsed all 134 JSON additions, compared their bytes with the working-tree copies, and decoded all 181 binary gzip additions. The retained initial auditor differs from the final auditor only in the documented `Input.Sources` correction at `audit.py:232`; its recorded hash matches.
- Checked 634 manifest entries across nested manifests, all 62 command records' retained stream identities, all 58 exact CLI inputs/results, and 61 checker/producer source identities. Both accepted inputs have exit zero; all 56 hostile inputs have nonzero exits and the recorded specific diagnostics. All 47 semantic hostile wire inputs differ from the accepted export under type-sensitive JSON comparison. The first comparison used ordinary Python equality and correctly exposed its float/int ambiguity; the completed comparison used canonical JSON to preserve types.
- Fresh captured raw semantic bytes equal Task 1 capture 3 exactly: 14,097,260 bytes, SHA-256 `ab3cbdfda577a3ae804eb76a94a8c7ac665cf00c69f992ad6804b6e0d0e27b8a`. The recorded result counts match the export and the independent-measurements record. No failed outcomes or accepted hostiles are present in the captured matrix.
- Read the captured unittest output: all ten tests passed in the recorded 32.094 seconds, without warnings or traceback. This is attribution to the retained run, not a fresh reviewer run. No Go suite, full audit, full Python suite, or old campaign was rerun.
- Named outside-diff risk: independence of reconstruction and strict validation. Checked the reused C4 `Automaton`, `subsets`, `product`, `shortest` and monitor constructors, C3 descriptor/plant/controller constructors and normalization, and C2 strict JSON loader/comparator. These inspectably provide independent finite operations; typed equality rejects bool/int and extra/missing keys, and duplicate keys/nonfinite constants are rejected by the loader.
- Named outside-diff risk: exact capture/source identity. Checked the unchanged Task 1 capture helper and the archived module-relative compiled roster. The helper compares the actual module-relative Go source list to the manifest and hashes the recorded producer inputs. This does not cure the Python execution-binding finding above.
- Named outside-diff risk: executable provenance. Checked the full inherited import chain and ran only the temporary-copy cache experiment described above. There is no evidence from this review that the existing archived semantic results were wrong; their claimed executable identity is insufficiently guaranteed.
- The source argument, manuscript reconciliation and broader Phase 3C claims are outside this task gate. The account explicitly retains those boundaries at `docs/foundational/phase3c5/03-independent-audit.md:132–141`.

# Task quality: REJECT

The finite auditor, adversarial coverage and retained archive are otherwise coherent and narrowly scoped. Exact executable provenance is load-bearing for this task; the demonstrated mismatch must be repaired and independently re-reviewed before approval.
