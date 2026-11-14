# Independent Phase 3C.4 audit

## Question and result

Does the complete frozen export faithfully represent residual-language profiles,
all synchronized product transitions and nonblocking classifications under the
finite contract? The independent Python reconstruction agrees with every exported
field. This is bounded model evidence supporting the separate conditional argument,
not a mechanized theorem or a runtime/refinement qualification.

The reconstructed roster has 72 plants, seven monitors, 324 controllers and 2,772
products, with 29,156 reachable states and 35,672 relation arcs. Of the product
states, 14,410 can complete; 886 products can initially complete and 656 are
nonblocking. The 21 monitor comparisons include 10 equal initial languages and
six equal profiles, 223 cross-state residual decisions, 850 language subset-pair
nodes and 89 profile subset-pair nodes. All 2,376 profile-equal substitution rows
are retained: 792 nonblocking and 1,584 blocking. The generic strictness control
is separate from these counts.

## Independence and reused semantics

`tools/audit.py` freshly constructs epsilon closures, complete deterministic
subset-pair searches, shortest separating words, bidirectional profile coverage,
every opposite-state separator, synchronous products, reverse marked reachability,
context substitution rows, negative paths and the private lasso with completion
escapes. It does not call or translate Go operational predicates, execute the
predecessor audit, or use producer outcomes/counts as an oracle. Exact typed
recursive comparison checks the entire finite schema, including all keys,
ordering, marks, arcs, certificate bits and redundant counts. JSON duplicate
keys and nonfinite constants reject during decoding; bool/int/float substitutions
reject without Python's bool/int equality shortcut.

The unchanged Phase 3C.3 Python plant, monitor, controller and descriptor
constructors are reused and attributed to
`docs/foundational/phase3c3/tools/audit.py`. Their unchanged Phase 3C.2 dependency
supplies configuration/action constructors and strict JSON utilities. Both Python
files have hard-coded verified SHA-256 identities checked **before import** and
are additionally included in the checker source manifest. No earlier Python
product, language comparison, grouping, history replay or outcome is invoked.
This reuse is a declared dependency, not a claim of independence from all earlier
model definitions.

The predecessor compressed hash is
`86a9d22f9808a67dd0066b6b77dd2b3c01a2125d10863c6fbe1369097f6624c8`;
its 25,384,964 decompressed bytes have hash
`a073a22fe8ac62e05f2d15b0f5146b19db9606e38ce4619261e7e4108e27269d`.
All 27 predecessor source records match current bytes. Reconstructed inherited
descriptor relations are compared after sorting states/edges, because the earlier
Python constructors use a different permissible BFS sibling order. Only after
that comparison is the authenticated predecessor descriptor serialization reused.
New product and certificate BFS order remains strict. Imported configurations,
roots, memberships, generations and bound actions are never erased or normalized
away. Monitor rebinding uses the complete selected plant signature.

The producer manifest is checked against actual file hashes. Its nine compiled
local sources also match an independent read of Task 1's saved clean
`go-list-deps-resolved.json.gz`; this is an attributed existing enumeration,
not a new `go list` run. Go source bytes are hashed, not used for semantics.

## Adverse cases and falsification boundaries

The release-first/doomed pair has equal initial language and unequal profiles.
Its shortest monitor mismatch is reserve AB; its actual core-product witness
first commits AB and then reserves. The doomed path subsequently takes tau into
an unpaid residual without a completion suffix. The two core products have
opposite nonblocking classifications despite both initially permitting completion.
A forged union-based comparison is rejected.

Duplicate residual states are accepted through one-to-many coverage, not a
bijection. A direct generic test has a residual language `{a,b}` opposed to
separate residuals `{a}` and `{b}` after a common prefix. The checker retains the
different required separators; no suffix distinguishes that residual from the
opposite union. Empty reach sets remain different from reached empty-language
residuals. The unmarked strictness control's profile mismatch cannot distinguish
contextual nonblocking: conjunction prevents every product state from marking,
including in any pointed compatible context.

The loop witness retains an actual unpaid prefix, two tau edges closing the cycle
1→3→1, and shortest marked escapes from both cycle states. Both cycle states are
unmarked but can complete. This refutes nonblocking-implies-inevitable-completion
in the abstract finite-monitor grammar. It makes no fairness, new physical
history, crash-capture or compiler/runtime realization claim. Validated ALAN local
programs' strongly terminating control grammar is not shown to realize this loop.

## Executed qualification and archives

All fresh records are under `evidence/task2/qualification-1/`:

- `tests.json`: `python3 -m unittest discover -s docs/foundational/phase3c4/tools -p 'test_*.py' -v` passed all nine test methods, including 39 semantic mutation subtests and seven malformed JSON decoding controls.
- `cli-commands.json`: two exact valid CLI executions (plain and gzip) accepted; all 48 hostile CLI executions rejected with expected diagnostic substrings. The hostile corpus contains the 39 semantic mutations plus nine byte-level controls. Every exact input is retained, including invalid UTF-8. The plain valid file's exact raw bytes are preserved by `export-third.json.gz`; its original scratch path and raw hash are recorded.
- Semantic attacks include a removed non-tree product edge with discovery parents preserved and counts repaired, union-for-profile forgery, missing coverage/matches/separators, bad shortest words and acceptance bits, missing language arcs, dead/unreachable confusion, altered tau/marks, external vetoes, identity/generation/root changes, invented completion/nonblocking, false-context omission, invalid lasso/debt/escape/negative paths, source forgery/omission, unknown keys, duplicate records and strict type attacks.
- `export-third-command.json` and `fresh-equality.json`: fresh `go run ./research/residualinterfaces/cmd/export` exited 0 with empty stderr. All 8,271,854 raw bytes equal Task 1 exactly, SHA-256 `b40f1d5ac0aea0f43999502ea318cb128f545eb8fb1fd129451915a5aa84bf1f`; source manifests also match exactly.
- `checker-sources.json`, `environment.json` and `archive-manifest.json` bind checker/dependency sources, Python/platform and explicit cache settings, exact inputs, stdout/stderr, argv and exits. Capture requires a new attempt directory and never overwrites an earlier attempt.

No hostile input was initially accepted. Two preparation failures remain visible:
`attempt-1-audit.py` and `attempt-1.json` preserve the valid export's initial
rejection due to inherited descriptor ordering; `attempt-2-test_audit.py` and
`attempt-2-test-failure.txt` preserve a test fixture construction error using a
nonexistent identity path. The latter failed before passing an input to the
checker; it was corrected to the actual `Base.Originals[0]` field. Neither is
reported as a semantic counterexample or a successful hostile acceptance.

Reproduce the checker with:

```sh
python3 docs/foundational/phase3c4/tools/audit.py docs/foundational/phase3c4/evidence/task2/qualification-1/export-third.json.gz
python3 -m unittest discover -s docs/foundational/phase3c4/tools -p 'test_*.py' -v
```

`tools/capture.py NEW_ATTEMPT_NAME` records a fresh complete qualification with
explicit `/tmp/alan-phase3c4/` caches; the supplied name must not exist. Existing
CLI records provide the command and exact archived input for each individual
hostile reproduction.

## Limits and delivery boundaries

Source argument, producer execution, independent qualification, review and
integration remain separate. This audit does not establish universal ALAN
congruence, necessity/minimality of profiles, divergence preservation, dynamic
shared authority, full-agent substitution, distributed ownership or split effects.
It repeats no Go test/race/Staticcheck suite, whole-root test campaign, old
transport probes, concrete runtime/endpoint history or SQLite operation. Earlier
Go and root qualification belongs to Task 1/controller evidence. Manuscript
reconciliation and independent review belong to their subsequent tasks.
