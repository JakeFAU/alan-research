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
before import, including all transitive executable dependencies:

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
The Task 1 capture helper separately checks the module-relative compiled roster
against recorded `go list -deps -json`; its exact source identity is included in
Task 2's checker manifest. The new helper does not modify Task 1's helper.

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

`evidence/task2/qualification-1/` retains a fresh producer capture, exact CLI
inputs, commands, explicit environment overrides, stdout/stderr, exit statuses,
source identities and an artifact manifest. Raw semantic bytes equal Task 1
capture 3 exactly: 14,097,260 bytes, SHA-256
`ab3cbdfda577a3ae804eb76a94a8c7ac665cf00c69f992ad6804b6e0d0e27b8a`.
The producer source bindings are unchanged. The fresh directory refusal prevents
replacing captured evidence.

The exact CLI matrix contains valid plain and gzip input, 47 semantic attacks,
and nine JSON/encoding attacks. Semantic cases include deleted non-tree arcs
with repaired superficial counts, reversed inclusion, lost equivalent minimum
members, MAX-for-MIN, empty-family confusion, missing/forged H, false hits/misses,
wrong shortest words/prefixes, forged product states/marks/tau, altered signatures
and roots, source hashes/uses/omissions, wrong generic observations, bad actual
paths and discarded blocking substitution rows. Malformed cases include top and
nested duplicate keys, NaN, infinities, overflowing numbers, trailing JSON,
invalid UTF-8 and truncated gzip. Input archives are lossless: decompress
`*.input.gz` once to recover exactly the bytes presented to each CLI, whether
those bytes were plain JSON or gzip JSON. Command records retain the temporary
input path and each result record retains its byte length/hash.

The first auditor run incorrectly retained inherited `Use` fields in
`Input.Sources`; the schema requires Path/SHA256 only. That implementation error
rejected valid input, was corrected against the schema, and is preserved with
the original auditor bytes and result in `evidence/task2/red-initial/`. It was
not a producer counterexample. No semantic evidence was edited to make it pass.
All ten tests passed; both valid CLI inputs were accepted and all 56 hostile
inputs were rejected. No hostile input was initially accepted. Exact results are
recorded in `qualification-1/results.json`. The archive/input hashes were checked,
and a repeated capture command exited 1 while leaving all 311 captured files
unchanged (`fresh-directory-refusal.json`).

Commands:

```sh
python3 docs/foundational/phase3c5/tools/audit.py docs/foundational/phase3c5/evidence/task1/capture-3/export.json.gz
python3 -m unittest discover -s docs/foundational/phase3c5/tools -p 'test_*.py' -v
python3 docs/foundational/phase3c5/tools/capture_audit.py FRESH_OUTPUT_DIRECTORY
```

## Remaining boundary

Exact finite reconstruction can expose a producer mismatch; agreement cannot
establish an unrestricted theorem. The reused descriptor/automaton utilities
and written contract remain shared assumptions. This task ran no Go test suite,
root build/lint campaign, runtime qualification, SQLite open, old exhaustive
transport probe, live model call or physical crash history. It does not qualify
dynamic authority, split effects, arbitrary contexts, divergence, fairness,
refusals, exact costs, distributed ownership or global contextual necessity.
Manuscript reconciliation and independent review belong to the controller.
