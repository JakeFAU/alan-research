# Independent finite antichain audit

## Question and result

Does the accepted `alan/phase3c5/v1` export exactly realize the frozen finite
contract, including directional inclusion, minimal residual families, local
finite suffix tests, and all synchronized products? The independent Python
reconstruction agrees with the captured export on every typed field. This is
finite executable qualification; it is not a proof of unrestricted contextual
necessity, distributed recovery, or whole-agent substitution.

`tools/audit.py PATH` accepts plain or gzip JSON. It constructs its own expected
export and performs complete typed structural comparison. Neither a supplied
verdict nor a supplied count supplies an oracle. Duplicate keys, extra keys,
missing records, mismatched array order, null/list confusion and bool/int
confusion are rejected. Diagnostics identify the mismatching structural path.

## Independence and exact reuse

No Go operational predicate is imported, translated, or called by the auditor.
New directional inclusion checks scan complete epsilon-closed subset graphs for
left-accept/right-reject and its reverse, retaining independent shortest words.
New code constructs language equality classes, all strict inclusion relations,
MIN/MAX families, complete profile comparisons, both coverage directions and
every local finite H. Actual automaton acceptance checks every H source hit and
target miss. The initial language remains an independent component of J.

The following unchanged Python modules are checked against these SHA-256 values
and compiled from those same byte buffers, including every transitive executable
dependency. A private importlib facade redirects the unchanged predecessors’
explicit loaders to these buffers; it does not consult Python bytecode caches:

| Module | SHA-256 | Used boundary |
| --- | --- | --- |
| `phase3c4/tools/audit.py` | `163fb14ea65af59a507e52b88c188c0b4fe2df65b40e8d2d5c7aef7f076fc88e` | Seven monitor constructors, epsilon-NFA closure/transition and complete subset traversal, synchronized product/reverse reachability, shortest product paths |
| `phase3c3/tools/audit.py` | `287be61a4e632b37510d401a658be7fa3fea465697d5586b43f07a58b2b7ad75` | Descriptor/alphabet, declarative plant and controller constructors, descriptor normalization |
| `phase3c2/tools/audit.py` | `1dc485fb2e3f0594b6122daa3a63ec05a12c8f5a6b9988afe5d2a775a768f193` | Configuration/action bindings, strict JSON parsing and typed structural equality |

Paths in that table are relative to `docs/foundational/`. No predecessor
`reconstruct`, audit campaign, product verdict, comparison verdict, or old
runtime probe is executed. Descriptor constructors independently check all
seven inherited monitors, 72 plants and 324 controllers against the pinned
predecessor before preserving their serialization. The predecessor compressed
and raw lengths/hashes and all 40 inherited source identities are checked.
The new 54-entry producer source manifest is reconstructed from current files.
The initial Task 2 capture used the unchanged Task 1 helper to check the
module-relative compiled roster against recorded `go list -deps -json`; its exact
source identity is included in the historical `qualification-1` manifest. The
Python-only repair reuses that unchanged producer evidence and does not rerun Go
commands or modify Task 1’s helper.

Every product is independently rebuilt using those unchanged Python automaton
utilities: all relation arcs, private moves, conjunctive marks, reverse
CanComplete labels and initial/nonblocking observations. Legal-trace controls,
shortest adjacent negative paths, generic controls and finite suffix tries are
also rebuilt. This reuse remains independent of the Go producer, but is not a
second independent implementation of the predecessor Python utilities.

## Measured finite evidence

The reconstructed roster has nine monitors and 36 pair comparisons. There are
1,442 residual comparisons, 5,208 language subset-pair nodes, and 146 profile
subset-pair nodes. Equal-pair totals are 13 initial languages, six FULL profiles,
seven MIN profiles, 13 MAX profiles and seven J profiles.

The 3,564 products contain 45,800 states and 59,492 arcs; 22,716 states can
complete. Initial completion holds for 1,234 products and nonblocking for 760.
All 2,772 J-equal substitution rows remain: 844 nonblocking and 1,928 blocking.
Each new monitor has 52 nonblocking and 20 blocking plant contexts; all 324
plant/controller contexts for each remain blocking.

The independent profile reconstruction finds 81 failed local coverage targets,
including 28 with H empty; the largest H has three words. The roster happens to
contain no singleton-epsilon H. Dedicated controls separately check H empty
versus H containing epsilon, directional epsilon counterwords, equivalent
states, strict dominance and tau cycles. The multiple-word generic H remains
`{[apply:AB], [reserve-capacity:AB]}`. These are local separately rooted hit
tests, not global nonblocking distinguishers. The explicit local-hit product
has initial completion true and nonblocking false.

For the legal-trace control, the shortest words proving the two proper
inclusions are `[commit:BC]` and
`[commit:AB, reserve-capacity:AB, commit:BC]`. Every control plant state is marked
only for that control. These words establish legal trace differences, not paid
completion. The RF/doomed MAX counterexample, generic MIN initial-language
loss, and generic J nonnecessity remain intact.

The independent set-theoretic test enumerates all 256 families of subsets of a
three-word universe and all 65,536 ordered family pairs. It checks the product
minima law, finite-hit implication/coverage equivalence, mutual coverage/minima
equality, every-family-member coverage and hit preservation by minima. This is
an exhaustive finite algebra check, separate from the conditional source proof.

## Qualification and preserved failures

`evidence/task2/qualification-1/` remains the original immutable capture, with a
fresh producer export, the original cache-capable unittest command and all exact
CLI inputs/results. Its raw semantic bytes equal Task 1 capture 3 exactly:
14,097,260 bytes, SHA-256
`ab3cbdfda577a3ae804eb76a94a8c7ac665cf00c69f992ad6804b6e0d0e27b8a`.
It recorded ten passing tests, two accepted valid inputs and 56 rejected
JSON/wire attacks. No JSON/wire attack in that corpus was initially accepted.
That result did not establish executable-source binding.

Independent review subsequently reproduced a distinct successful provenance
attack: a timestamp-valid, same-size C4 bytecode cache substituted
`ALPHABET = ('CACHE',)` while the on-disk predecessor source hashes remained
correct. `-B`, `dont_write_bytecode` and `PYTHONDONTWRITEBYTECODE` prevent cache
writes, not reads. This demonstrated a guarantee gap; it does not establish
that the original qualification actually used stale code or that its finite
semantic results were wrong. The review returned **REJECT** for this finding.

`evidence/task2/fix-1-red/` retains the exact successful reviewer script/output,
a fresh pre-fix replay, original command metadata, and the pre-fix auditor,
test and capture sources from commit
`b079c67db3298688f32f2c5c8218a01f81d45601`. The unchanged original reproduction
reads its current working tree, so reproduce the historical failure against
that commit (or a temporary tree with the retained C5 auditor and unchanged
C4/C3/C2 auditors). Do not interpret running that old script against repaired
sources as a reproduction of the historical experiment.

The correction uses `tools/source_loading.py` to read the declared sources into
byte buffers, hash them, verify predecessor pins, and compile those buffers.
The buffers also supply the recorded source identities. It covers C5→C4,
C4→C3 and C3→C2, local auditor/test/capture imports, the helper, and the new
source-only test runner. Direct CLI entrypoints bootstrap the helper from one
read buffer, then execute their own recorded source buffers. The runner loads
test modules directly from the bundle before constructing the unittest suite.
Each child auditor and runner reports its buffer identities; the capture checks
these against its own records. No global importlib hook or predecessor edit is
needed. The trusted boundary is the Python installation and standard library;
this is not a general sandbox, nor protection for arbitrary external callers
that choose their own cache-capable importer.

`evidence/task2/qualification-2/` is the corrected qualification. It reuses the
unchanged producer export; no Go or root suite is repeated. It records thirteen
passing tests, two accepted valid CLI inputs and all 56 JSON/wire attacks
rejected. The three added tests cover nine cache boundaries, changed pinned
source rejection before execution, and execution of a recorded buffer after
its disk file changes. Every cache case first proves the normal loader executes
the prepared poison under `-B`, then checks that the source-only chain executes
the intended source. These successful stock-loader controls remain distinct
from the corrected loader’s passing result. The source-only runner replaces
bare unittest discovery for accepted qualification; the old command and its
result remain historical and have not acquired this guarantee retroactively.

The exact CLI corpus contains 47 semantic attacks and nine JSON/encoding
attacks. Semantic cases include deleted non-tree arcs with repaired superficial
counts, reversed inclusion, lost equivalent minimum members, MAX-for-MIN,
empty-family confusion, missing/forged H, false hits/misses, wrong shortest
words/prefixes, forged product states/marks/tau, altered signatures and roots,
source hashes/uses/omissions, wrong generic observations, bad actual paths and
discarded blocking substitution rows. Malformed cases include top and nested
duplicate keys, NaN, infinities, overflowing numbers, trailing JSON, invalid
UTF-8 and truncated gzip. Decompress each `*.input.gz` once to recover exactly
the bytes presented to the CLI, whether those bytes were plain or gzip JSON.
Commands, environment overrides, stdout/stderr, exits and byte hashes remain
inspectable. Both capture helpers refuse existing output directories.

The first original auditor run also had a separate serialization mistake:
it retained `Use` fields in `Input.Sources`, while the schema requires
Path/SHA256 only. Its source snapshot and valid-input rejection remain in
`evidence/task2/red-initial/`. That was an auditor implementation error, not a
producer counterexample or a hostile-input rejection. No semantic evidence was
edited to make it pass.

Current commands:

```sh
python3 docs/foundational/phase3c5/tools/audit.py docs/foundational/phase3c5/evidence/task1/capture-3/export.json.gz
python3 docs/foundational/phase3c5/tools/run_audit_tests.py
python3 docs/foundational/phase3c5/tools/capture_audit.py FRESH_OUTPUT_DIRECTORY
```

## Remaining boundary

Exact finite reconstruction can expose a producer mismatch; agreement cannot
establish an unrestricted theorem. The reused descriptor/automaton utilities
and written contract remain shared assumptions. Their historical standalone
entrypoints are unchanged and are not retroactively repaired or qualified by
this task’s source-only loading boundary. This task ran no Go test suite,
root build/lint campaign, runtime qualification, SQLite open, old exhaustive
transport probe, live model call or physical crash history. It does not qualify
dynamic authority, split effects, arbitrary contexts, divergence, fairness,
refusals, exact costs, distributed ownership or global contextual necessity.
Manuscript reconciliation and independent review belong to the controller.
