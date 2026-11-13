# Phase 3B.7 matched concrete effect histories

The two sealed single-host pairs each exhibit opposite endpoint application
facts with identical declared local observations. The normal and race builds
each executed one pair from a common closed seed copied into separate paths
before execution. The branches ran sequentially with separate endpoint
processes and databases at the same literal URLs within each pair. Both
principals shared a single SQLite Store owner. The checker accepted the pairs
only after a post-capture correction to its authorization-event validation.
[Corrected reports](evidence/qualification/corrected-checks.json) ·
[Initial rejections](evidence/qualification/initial-checks.json)

## Observed cuts and world facts

| Observation within each pair | Normal | Race |
| --- | --- | --- |
| Complete seed file sets | Same lengths and hashes | Same lengths and hashes |
| Complete post-open, stopped, reopened and final `inspect.json` | All four byte and structure matches | All four byte and structure matches |
| Complete final `run.json` | Byte and structure match | Byte and structure match |
| Tables and rows in each final export | 41 tables, 94 rows | 41 tables, 94 rows |
| Zero branch selected contacts / applications | 0 / 0 | 0 / 0 |
| Applied branch selected contacts / applications | 1 / 1 | 1 / 1 |
| Applied branch endpoint receipt | Idempotent, delta 1, mutation 1, total 1 | Idempotent, delta 1, mutation 1, total 1 |
| Phase 3B.6 `StoreView` at stopped, reopened and final | All three match | All three match |
| Sealed archive members | 122 | 122 |

Each pair contains one shared seed export and four exports per branch, for
nine distinct export artifacts per pair and eighteen across the experiment.
The seed-copy prerequisite compares the complete seed file sets by hash. The
four later comparisons directly cover the full export, including `Settings`,
`Roots` and every table, row, column, descriptor and nested value. No
identity renaming, field deletion, table filtering or post-hoc reordering was
applied. These are within-pair comparisons.
[Capture summary](evidence/qualification/capture-summary-initial.json) ·
[Comparison](03-comparison.md)

The initial open advances both principals from epoch 1 to 2. At the stopped
cut A is at sequence 4 and B at sequence 5, both at epoch 2. Reopening advances
both to epoch 3. The final sequences are A at 4 and B at 8. The selected BA
service retains its request-bound grant, historical first allowance, one
serving-root debit and original allowed gate. Exact revocation targets root
`b`, the fresh gate refuses, and the distinct root `c` survives. Neither
branch retains a selected authentic result, effect settlement or consumption.
[Corrected normal report](evidence/qualification/corrected-commands/check-normal-corrected-1/output.log.gz) ·
[Corrected race report](evidence/qualification/corrected-commands/check-race-corrected-1/output.log.gz)

## Causal interpretation and model relation

In both branches the worker reached the acknowledged pre-call hold after
admission, and the scheduler acknowledged its own in-flight hold with the
selected authorization. The zero branch held both through confirmed owner
death. Complete, process-bound endpoint journals and snapshots establish zero
contacts and applications. In the applied branch only the worker was released
before death. The exact endpoint B post-commit acknowledgment places the
selected application before death. Successful receipt was withheld through
death. The recovery refusal supplies a fresh gate disposition while the
original application fact differs across the pair.
[Protocol](02-protocol.md) · [Sealed normal and race archives](evidence/capture/)

At the stopped, reopened and final cuts, the Phase 3B.6 model prefixes have
lengths 1, 2 and 4 for zero and 4, 5 and 7 for applied. Each maps concrete
epoch 2 to model activation 1 and epoch 3 to activation 2. The selected
applied trace infers launch as a prerequisite to its witnessed ingress and
application. The independent ledger audits the existential prefixes, with
the seed and post-open cuts outside its post-admission domain. The complete
export and final report equality are direct comparisons of the archived
bytes, independent of the smaller model projection.
[Phase 3B.6 finite model and historical relation](../phase3b6/07-results.md) ·
[Comparison](03-comparison.md)

If a decision takes only the declared local observation as input, it receives
the same input in each branch. The selected endpoint application fact
differs, so such a decision cannot recover it for both histories. The
source-conditioned part of this conclusion uses the honest pinned producer
and endpoint, process identity and barrier ordering, and complete endpoint
coverage. Endpoint and process evidence in the full dossiers distinguish
the branches. The checker uses the archived SQLite files as opaque bytes and
relies on the frozen procedure for their connection to the exports. The
source and binary metadata are unsigned. The authority checks reuse the
existing independent raw checks through a narrow adapter.
[Qualification](06-qualification.md)

## Scope

The observed contrast covers one selected effect with a native idempotent
endpoint, one owner loss in each branch, and the prescribed reopen,
revocation and refusal sequence. Phase 3B.6's historical applied branch and
finite observation result retain their prior scope. The present result
supplies a genuine held pre-call counterpart with matching declared local
exports. Independent provisioning, concurrent cloned owners, distributed
crash recovery, non-idempotent endpoint behavior, general recovery
equivalence, universal liveness and closure beyond two agents remain open.
