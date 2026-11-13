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

## Review round 1 corrections

Independent review rejected `eca3b3bee03b32c7bfffc65cba124d3ba638446c`
with four reproduced acceptance gaps. All four are addressed in this follow-up;
the original REJECT remains in `task-2-review.md`. Original probe source and
output were copied byte-for-byte from `/tmp/phase3b7-review-probes/` to
`evidence/task2/review-r1/probe_test.go.txt` and `result.log`; the original
`overlay.json` is also retained. These copies were not edited. This correction is subject to independent re-review, not yet approval.

The owning source was inspected in the Task 1 frozen prepared checkout
`/tmp/alan-phase3b7/producer-verified`, with the same runtime commit/tree and
complete preparation manifest identified above. Relevant native definitions
are `internal/fixture/resource.go:35` (Snapshot), `:218` (Inspect),
`internal/executor/runner.go:27` (RunReport), `:498` (finishRun),
`internal/executor/program.go:81` (Intent),
`internal/executor/composition_inspect.go:11` (RootCompletion),
`internal/boundary/boundary.go:17` (Event), and
`internal/executor/external.go:174` (the emitted authorization event).
Go 1.27.4's `src/flag/flag.go:1075` was also inspected for termination,
flag spellings, separate values, and repeated assignment behavior.

Corrections:

- Endpoint exports now validate all eleven ordinary snapshot fields and native
  types, including required zero-valued fields. Both endpoints and both cuts
  are checked. Native zero `Requests:null`, `Applications:[]`, and explicit
  integer counters remain valid. Applications:null and missing fields fail.
- Both final successful Run reports require all ten native fields, typed nested
  records, non-null report arrays, and nonnegative counters. Their Roots join
  the final raw export. Complete-schema validation compares a typed round trip
  only as a validation step; the original entire JSON inputs still determine
  full report equality. Equal truncated reports therefore remain equal but
  fail Raw and Supported.
- The passive authorization boundary joins the selected retained wire's entire
  Subject/Generation/Epoch/Sequence and Intent/Attempt, requires pre-transport,
  empty optional root/admission selectors, and nil InFlight. The inherited
  scheduler acknowledgment's enriched InFlight remains separately checked.
  The synthetic fixture's previously empty authorization Phase was corrected
  to the owning source's pre-transport value; producer source is unchanged.
- Source argv is parsed with Go's flag parser. Exactly one exact test.run
  selector and effective count one are required. Single/double dash and
  separate/equal-value spellings are supported. Repeated selectors, effective
  overrides, selector text used as another option's value, positional tails,
  termination before selection, skip flags, and other unsupported flags fail.
  The supported driver flags are listed in `03-comparison.md`.

`review_r1_test.go` adds resealed deletion/null/wrong-type/contradiction/override
controls, coherent both-cut zero deletions, identically truncated reports,
nested report omissions, and lawful flag forms. It also directly probes each
ordinary snapshot member's presence/type. Existing normal/race synthetic
positives and original causal, equality, malformed-input, and command tests
remain in the full suite. No genuine pair ran or archive was rewritten; no
SQLite API was used. Old anchored authority entry points and their source are
unchanged in this correction. The earlier four historical output-byte
regressions remain prior evidence and were not rerun for these isolated fixes.

All Go commands below used the offline environment prefix shown in the
original verification section (including STATICCHECK_CACHE for Staticcheck).
Complete outputs are retained under `evidence/task2/review-r1/`:

| Command | Exit | Output |
| --- | --- | --- |
| `go test -count=1 ./research/matchedcapture -run '^TestReviewR1'` (first test-writing attempt) | 1 | `red-setup-error.log`: missing test-file brace; no tests ran |
| Same command after correcting test syntax, before checker fixes | 1 | `red.log`: semantic RED failures in all five initial control groups, 7.815s |
| Same command after checker fixes | 0 | `green-first.log`: initial new controls, 7.857s |
| `go test -count=1 -v ./research/matchedcapture ./cmd/check` | 1 | `normal-wrong-command-path.log`: actual checker passed in 17.225s; nonexistent root command path failed setup |
| `staticcheck ./research/matchedcapture/... ./cmd/check` | 1 | `staticcheck-wrong-command-path.log`: nonexistent root command path |
| `go test -overlay=/tmp/phase3b7-review-probes/overlay.json -run '^TestIndependentReviewProbes$' -v ./research/matchedcapture` | 0 | `probes-green.log`: all four original attacks rejected, 0.554s |
| `go test -count=1 -v ./research/matchedcapture/...` | 0 | `green-normal.log`: checker 18.005s; actual nested command 0.133s |
| `staticcheck ./research/matchedcapture/...` | 0 | `staticcheck.log`: empty, no diagnostics |
| `go test -race -count=1 -v ./research/matchedcapture/...` | 0 | `green-race.log`: checker 216.329s, nested command 1.126s; no race diagnostics |

Scoped normal/race suites include every test in `check_test.go`,
`compare_test.go`, `review_r1_test.go`, and the actual nested command package.
Staticcheck covers the same affected package tree. Broad authority/runtime
verification was not repeated because those sources are unchanged. Formatting
(`gofmt -l research/matchedcapture`) and `git diff --check` emitted no output
and exited 0. Byte comparisons confirmed preserved review probe/result files
against their originals; SHA-256 values are respectively
`93d7dcd40f2e2300938a8f72641ffee248701ca4401992417730f651bbe2f674`
and `45fdd362a1272b3843863ee72128093852a6bff7ba5608b3aefd9d55b9319beb`.

Full follow-up source/test/documentation/evidence diff was self-reviewed.
All four reproduced findings now have passing rejection controls, while lawful
synthetic positives pass. Independent re-review and genuine qualification
remain outstanding; all source-conditioned and single-host limits above
remain in force. There is no new manuscript empirical claim in this fix.
