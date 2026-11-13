# Phase 3B.7 independent comparison

This is an offline, source-conditioned checker for the producer protocol in
[02-protocol.md](02-protocol.md). Task 2 qualifies its behavior with synthetic
fixtures and historical JSON. No genuine matched pair was executed in this
stage, and the synthetic positive control is not evidence of such a pair.

## API and command

`research/matchedcapture.CheckDirectory(dir string) (Report, error)` reads
JSON/JSONL and hashes regular archive members as opaque bytes. It never opens
SQLite through a database driver and never runs the producer. Filesystem
errors return an error; semantic failures remain in the returned report.

```sh
env GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off \
  GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b7/go-build \
  go run ./research/matchedcapture/cmd/check -pair /absolute/sealed-pair
```

For qualification, build and freeze the command binary first, then invoke it
with `-pair`. Exit 0 means `Supported`; exit 1 returns a rejected evidence
report; exit 2 denotes usage, filesystem, or output failure. JSON output uses
fixed indentation and relative input paths, with no current timestamps.

## Separate verdicts

| Field | Meaning |
| --- | --- |
| `Seal` | Exact member-length/SHA-256 coverage; regular members; strict JSON/JSONL; required capture exports and successful export outcomes |
| `Source` | Complete prepared-source manifest equals the externally frozen anchor; supported platform, build/test identity, race mode and offline environment join |
| `Raw` | Seed provision joins, complete export schema, principal coordinates, two exact full-export Open transitions, immutable raw history, original selected binding, grant, allowance/debit, exact revocation, fresh refusal and missing local custody |
| `Protocol` | Full prescribed action sequence, process/Store/token/config/command joins, pre-call and scheduler acknowledgment, release/kill/reap ordering, open exports, owner boundary corroboration, isolated branch processes and rebound URLs |
| `World` | Complete process-bound A/B journals and stopped/final snapshots, exact selected wire and receipt, cardinalities zero/one, closed endpoint bytes present |
| `Equality` | Separate seed-copy prerequisite and complete structural/byte equality for post-open, stopped, reopened and final `inspect.json` |
| `ReportEquality` | Complete final `run.json` comparison, separately reported |
| `Model` | Post-admission projection in the existing Phase 3B.6 `StoreView` domain, independently audited trace, and projected pair equality |

`World.Valid` alone is a check of endpoint evidence. Certifying an application
at the stopped cut additionally requires the source and causal protocol
verdicts. Certifying zero requires the uninterrupted pre-call hold through
confirmed owner death and complete endpoint evidence. `Supported` requires
all verdicts and complete structural equality. Producer `status` is retained
as provenance; it never supplies a semantic verdict. A failed attempt can
still have independently valid raw/equality/world evidence.

`Inputs` records the exact length and SHA-256 of every member, including the
member manifest itself. `RawBranches` retains the unchanged authority-checker
verdicts and selected binding objects. `Worlds` preserves endpoint facts;
`Models` preserves the three supported cuts per branch, their StoreViews,
event witnesses and independent ledgers.

## Complete observation

There is one shared seed export. `Equality["seed-copies"]` checks the complete
sealed seed SQLite file set against all three entries in `seed/copies.json`.
Its byte-equality flag denotes hash equality under the frozen copying
procedure; no separately inspected branch seed exports are invented. The four
subsequent exports per branch are compared directly.

Object order and whitespace are ignored. Array order and multiplicity,
null/empty distinctions, exact strings, embedded descriptor strings, and
number tokens remain significant. Numbers never pass through floating point.
Duplicate object keys, invalid UTF-8, trailing JSON, and unsupported typed
schemas fail closed. Difference paths are deterministic, for example
`$.Tables.runs[0].projection` or `$.Settings.SQLiteVersion`.

The comparator covers the whole outer object, including `Settings`, `Roots`
and any `ReceiptEvidence`. It never constructs a stripped comparison object.
The frozen supported raw schema does not contain `ReceiptEvidence`; if it
appears, complete equality still observes it and raw validation rejects the
unsupported schema. It is not silently dropped to accommodate an old helper.
Both Open checks change exactly the two `runs[].epoch` values (1→2 and 2→3)
in the complete prior export and compare every remaining field exactly.

Native endpoint snapshots require every ordinary snapshot member with its
native type at both cuts for both endpoints. An empty native Requests slice
may be explicit `null`; Applications must be an array, and Total and Contacts
must be explicit integers. Successful final Run reports require all ten native
fields, complete nested records, and root summaries matching the final export.
Schema validation never replaces the original inputs to full equality.

The passive `external.authorized` event must carry the selected wire's complete
Subject/Generation/Epoch/Sequence and Intent/Attempt, `pre-transport` phase,
empty optional root, admission ordinal 1 matching the selected admitted worker,
and native nil InFlight. Run stamps this ordinal in its transport wrapper before
emission; it is not the constructor's zero value. The later
inherited scheduler acknowledgment has its own enriched InFlight check.
Source argv uses Go flag parsing and requires one exact capture selector and
an effective count of one. Single/double dash and separate/equal values are
supported; repeated selectors, trailing positional arguments, and unsupported
flags fail closed. Accepted driver flags are `test.run`, `test.count`,
`test.timeout`, `test.testlogfile`, `test.paniconexit0`, and `test.v`.

## Authority and model boundary

`delegationcorrespondence.CheckMatchedBinding` is a narrow adapter around the
unchanged independent raw binding and continuity checks and the existing
selected-effect binding helper. It accepts the native reciprocal manifest,
not a fabricated historical case. Existing anchored checkers and their pinned
protocols are unchanged. This is independent of the runtime producer, not an
independent reimplementation of all prior authority semantics.

The raw checks preserve the exact selected request, grant/root, wire,
authorization, reservation, original allowance and single serving-root debit.
They check the original epoch-2 gate, both-principal renewal, exact-b revoke,
sole fresh epoch-3 refusal, c's distinct surviving permission, and absent
selected result/settlement/consumption. The outer root summaries are also
joined to their retained roots and charges.

Model mapping covers stopped, reopened and final, with concrete epoch 2 mapped
to model activation 1 and concrete epoch 3 to activation 2. Seed and post-open
are outside the model's post-admission domain. For zero, the trace starts with
owner loss; for applied, launch is an inferred prerequisite before the matched
ingress/application and owner loss. Reopen, revoke-b and fresh-refuse extend
both traces. `AuditTrace` independently checks each event prefix and its local
projection. A projected match cannot substitute for full-export equality.

## Source, archive and qualification freeze

The embedded preparation manifest is an exact copy of Task 1's
`evidence/task1/prepared-source.json.gz`. Its decompressed SHA-256 is
`557e832fc84e7cbdb47d1d4767938f80ddd1ca1d77c8c647cb665a3436546994`.
It contains 637 runtime members from commit
`5b4662fd682a220ca3065b15e86672b7e56567e7`, tree
`44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`, four overlays, and the preparation
script. A dossier cannot replace this anchor by supplying new self-declared
source hashes. Source/build metadata and binary hashes remain unsigned
producer evidence, not independent binary-to-source attestation.

Before either qualification pair, freeze the producer sources, preparation
manifest, contract/protocol/comparison documents, checker sources (including
the reused correspondence and effectobservation packages), Go module files,
and built checker identity. Preserve each attempted normal/race capture and
its exit/output before comparison. Run the frozen checker twice against each
unchanged sealed pair and compare report bytes. Derived reports belong
outside the raw directory. Never open archived SQLite files or reseal a
genuine attempt to repair a mismatch.

The database-to-export connection is supplied by the frozen producer's
copy/export procedure. The checker independently verifies JSON consistency,
protocol joins and opaque byte seals, but does not independently verify
SQLite contents. A successful future pair would remain a bounded forked
single-host witness, not a general recovery, revocation, liveness or
N-agent theorem.

## Synthetic qualification

The fixture copies historical Phase 3B.4 JSON into new scratch directories,
invents distinct process transcripts and endpoint worlds, and uses explicitly
synthetic opaque database bytes. Genuine archives are never rewritten or
opened. Tests cover the lawful normal/race metadata controls, command
repeatability, all contract-listed rehashed falsifiers, coherent identity
rewrites, malformed inputs, source drift, seed multiplicity, and a raw/model
match whose full observation differs. RED and GREEN evidence and historical
byte regressions are retained under `evidence/task2/`.
