# Phase 3B.7 Task 2 implementation report

## Result and scope

Implemented the independent offline paired-history checker and opt-in command
against Task 1's actual `alan/phase3b7-matched/v1` schema. The result is tested
checker behavior on synthetic fixtures plus unchanged historical authority
exports. No genuine Phase 3B.7 pair was run. No production runtime, producer
source, `.envrc`, dependency, old protocol entry point, or genuine archive was
modified. No SQLite database was opened through a database API.

Base: `f238b778ac786d3efbd86f94e98068439b9e4ed2` on
`codex/phase3b7-matched-capture`. The Task 1 source anchor is runtime commit
`5b4662fd682a220ca3065b15e86672b7e56567e7`, tree
`44af03133ac7ba75fb6f6f8dc36629ed4cf52d98`. The exact embedded Task 1 preparation
manifest has 637 runtime members, four overlays and the preparation script;
its decompressed SHA-256 is
`557e832fc84e7cbdb47d1d4767938f80ddd1ca1d77c8c647cb665a3436546994`.

## Interfaces and decisions

- `research/matchedcapture.CheckDirectory(dir string) (Report, error)` reads
  JSON/JSONL and opaque member hashes. `CompareJSON` implements strict whole
  JSON comparison. `cmd/check -pair DIR` emits deterministic indented JSON;
  exit 0 is supported, 1 is a semantic rejection, 2 is usage/I/O failure.
- `Report` independently exposes `Seal`, `Source`, `Raw`, `Protocol`, `World`,
  `Equality`, `ReportEquality` and `Model`, together with every exact input
  digest, raw branch verdict/bindings, endpoint facts and model witness/ledger.
  Producer status never supplies a semantic verdict.
- Equality includes every outer field, Settings, Roots, all tables/rows and
  descriptor strings. Object order/whitespace are ignored; arrays, number
  tokens, strings, null/empty and multiplicity remain exact. Duplicate keys
  fail closed. Difference paths escape punctuation-bearing keys. Existing
  malformed byte streams can still be byte-equal; missing exports cannot.
- One shared seed export plus the pre-execution byte-copy manifest is the seed
  prerequisite. Four subsequent complete exports per branch are compared.
  No separate initial branch export is invented.
- Full outer exact renewal checks cover both principals at 1→2 and 2→3.
  Native seed manifest/provision/program/run/subject/permission/root joins,
  source anchors, exact process and wire identities, original debit, b revoke,
  fresh refusal, c survival and absent selected custody are checked.
- `delegationcorrespondence.CheckMatchedBinding` narrowly exposes unchanged
  raw binding/continuity and selected-effect checks. It creates no historical
  case wrapper and changes no old anchored protocol. This reuse is independent
  of the runtime producer, not a second implementation of the authority
  checker. Targeted adapter tests compare its raw verdicts to direct calls.
- The model maps stopped/reopened/final only, using the published StoreView
  domain and separate AuditTrace ledger. World validation checks exact
  endpoint evidence. At-cut certification additionally requires the source
  and causal protocol verdicts. Projected equality cannot certify full equality.
- Unknown schema fails closed. The frozen raw schema lacks ReceiptEvidence;
  if such a field appears, whole-export equality retains it and raw validation
  rejects the unsupported schema. It is never removed to appease a helper.

## Evidence and commands

Build/test commands used these offline Go variables; Staticcheck also used
the indicated cache:

```sh
env GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off \
  GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b7/go-build \
  STATICCHECK_CACHE=/tmp/alan-phase3b7/staticcheck
```

Final verification commands all exited 0:

```sh
go test -count=1 ./research/matchedcapture/...
go test -race -count=1 ./research/matchedcapture/...
# Add ALAN_EFFECT_OBSERVATION_CORPUS=/tmp/alan-phase3b5/corpus to the prefix:
go test -race -count=1 ./research/delegationcorrespondence -run '^TestMatchedAdapter$'
staticcheck ./research/matchedcapture/... ./research/delegationcorrespondence
```

The final ordinary suite took 8.350s for the checker and 0.111s for command
tests. The final race suite took 88.616s and 1.130s; adapter race checks took
3.264s. Staticcheck emitted no diagnostics. The command tests build the actual
checker and require two synthetic positive invocations to emit identical
report bytes. All package tests used Go 1.27.4 on darwin/arm64. Logs are in
`docs/foundational/phase3b7/evidence/task2/{green-normal,green-race,adapter-race,staticcheck}.log.gz`.
`gofmt -l` on the new Go files and `git diff --check` emitted nothing and exited 0.
The embedded source copy was byte-compared to Task 1 and its decompressed
manifest hash/counts were checked (`source-anchor.json`).

The synthetic lawful controls use historical JSON from
`/tmp/alan-phase3b5/corpus/raw-normal/applied-unknown-revoked-1243056310`, copied
into scratch fixtures with invented process transcripts and endpoint worlds.
The opaque synthetic database files explicitly say they are not SQLite. No
historical archive is resealed. The corpus provenance and source anchor are
recorded in `research/matchedcapture/testdata/README.md`.

The retained RED commands exited 1 after executing meaningful assertions:

| Log | Command suffix / failed behavior |
| --- | --- |
| `red-comparison.log.gz` | `go test ./research/matchedcapture`; permissive comparison missed hidden state, duplicate key, null/empty, array order, large number, descriptor and ReceiptEvidence differences |
| `red-causal.log.gz` | `go test ./research/matchedcapture -run 'TestResealedAttacks\|TestInvalidSeal\|TestStatusNotSemanticAuthority'`; permissive report accepted rehashed causal/world/raw attacks, invalid seal and producer-status shortcut |
| `red-source-joins.log.gz` | `go test ./research/matchedcapture -run 'TestSeedProvisionJoin\|TestUnexpectedProcessArtifact\|TestEmptySourceArgvFailsClosed'`; seed and extra-process joins were missing, and empty argv caused an index panic |
| `red-malformed-source.log.gz` | `go test ./research/matchedcapture -run 'TestCoherentIdentityAttacks\|TestSourceNullDependency'`; coherent identity controls already rejected; a null dependency caused a nil-pointer panic |
| `red-outer-schema.log.gz` | `go test ./research/matchedcapture -run 'TestProjectedMatch\|TestMalformedOuterSchema'`; equal null Settings/Roots wrongly passed raw validation |
| `red-seed-multiplicity.log.gz` | `go test ./research/matchedcapture -run 'TestSeedSourceMultiplicity\|TestWrongFixtureManifest'`; duplicate seed members and wrong scenario were accepted |
| `red-test-selection.log.gz` | `go test ./research/matchedcapture -run 'TestWrongTestSelection\|TestSyntheticRacePair'`; substring-only test selection admitted a wrong test |
| `red-self-review.log.gz` | `go test ./research/matchedcapture -run 'TestManifestCannot\|TestUnambiguousDifference'`; nonexistent zero-digest replacement member and ambiguous dotted key path were accepted |
| `red-byte-equality.log.gz` | `go test ./research/matchedcapture -run '^TestMalformedBytesStillReportByteEquality$'`; invalid but identical bytes lost byte-equality information |
| `red-missing-export.log.gz` | `go test ./research/matchedcapture -run '^TestMissingExportsDoNotClaimByteEquality$'`; absent files acquired a byte-equality flag after the preceding fix |

An initial development run of the premature-release attack had a test-helper
loop that kept moving the scheduler record until it indexed past the slice.
The helper was corrected to swap once; the retained causal RED was rerun and
contains semantic failures without that test bug. Intermediate GREEN runs
passed after each correction; final suites include all retained controls.

The four historical effect reports (N07 and applied-unknown, normal and race)
were compared byte-for-byte using binaries built from the base and the current
unchanged old entry point. All eight command invocations exited 0, stderr was
empty, and all four report pairs were identical. This is an output-byte
regression, distinct from inspecting the source diff. Exact lengths/hashes are
in `historical-byte-regression.json`; `compare-legacy.py` is retained.

The base binary was built in an independent `/tmp/alan-phase3b7/checker-base`
export made with `git archive f238b778ac786d3efbd86f94e98068439b9e4ed2`:

```sh
go build -o /tmp/alan-phase3b7/old-effect-check ./research/delegationcorrespondence/cmd/effectobservation
# Current research checkout, same offline prefix:
go build -o /tmp/alan-phase3b7/new-effect-check ./research/delegationcorrespondence/cmd/effectobservation
python3 /tmp/alan-phase3b7/compare-legacy.py
```

Both builds and the comparison script exited 0.

## Freeze and remaining obligations

Freeze the complete checker dependency closure (matchedcapture,
delegationcorrespondence, effectobservation), producer/preparation manifest,
contract/protocol/comparison documents, module files and built checker before
qualification. Preserve normal/race first attempts, including failures, and
run the frozen command twice without changing raw input. The comparison
format and qualification procedure are documented in `03-comparison.md`.

This checker assumes honest frozen source execution and unsigned producer
metadata. Binary hashes are retained but are not independently attested to
source. The database/export relation comes from the frozen producer; no
independent SQLite verification is claimed. Model launch is an inferred
prerequisite, not a separately retained launch record. Continuous endpoint
callback occupancy, general distributed recovery, general revocation,
completion/liveness and N-agent closure remain outside the result.

The manuscript's existing Phase 3B.6 discussion and conclusion still describe
the matched concrete pair as an open experiment (`paper/alan-dissertation.md`
around 1098 and 1145, with the corresponding LaTeX passages). Task 2 introduces
no empirical result that changes those statements, so manuscript/PDF changes
remain for qualification reconciliation. Independent Task 2 review and genuine
qualification are still required before claiming a matched concrete witness.

Full staged source/test/documentation diff and binary-fixture provenance were
self-reviewed; the scope contains no unrelated production changes. No push,
amend, merge or genuine matched-capture invocation was performed.
