# Phase 3B.3 reproduction

Use the research and public runtime repositories separately. The research
contract is frozen at `360cae66f5ffdd6e7a95a568552df5522122da21`; its final bytes
have SHA256 `f19d9c13f84c17496fabef91cb7575cf52972ff228c0b69e063cdc1c99100e12`.
The runtime's final executable source is
`759df6e3d2505f27e1287b2cfeb07bc4290afc0a`. The research executable source is
`7894bfae268e6dc0f5418bda3c6022f3a8a55540`. Exact command records and final
checker outputs are in [research-frozen](evidence/research-frozen/README.md),
with all boundaries in the [evidence map](06-evidence-index.md). Later delivery
commits do not replace those executable identities.

These commands use the existing Go 1.27.3 toolchain and cached dependencies.
There are no new dependencies, live model calls or paid endpoints. Runtime
qualification requires loopback listeners and child processes. Use fresh output
roots, never an extracted original evidence directory.

```sh
export GOTOOLCHAIN=local
export GOCACHE=/tmp/alan-phase3b3-reproduction-go
export STATICCHECK_CACHE=/tmp/alan-phase3b3-reproduction-staticcheck
export GOPROXY=off
export GOSUMDB=off
```

## Research regression and original model exports

From the research repository:

```sh
ALAN_TASK3_RAW=/tmp/alan-phase3b3-runtime-frozen/normal ALAN_TASK3_PROGRESS=1 \
  go test -json -count=1 -timeout=60m ./...
ALAN_TASK3_RAW=/tmp/alan-phase3b3-runtime-frozen/race ALAN_TASK3_PROGRESS=1 \
  go test -json -race -count=1 ./research/twoagent/... ./research/mutualobligations/... ./research/authorityservice/... ./research/delegationcorrespondence/...
ALAN_PHASE3B_RUNTIME_EXPORT=/tmp/alan-phase3b3-prior-reciprocal-normal \
  go test -json -tags=integration -count=1 ./examples/communication/mutual
ALAN_PHASE3B_RUNTIME_EXPORT=/tmp/alan-phase3b3-prior-reciprocal-race \
  go test -json -race -tags=integration -count=1 ./examples/communication/mutual
ALAN_PHASE2_EVIDENCE=/tmp/alan-phase3b3-prior-authority-normal \
  go test -json -tags=integration -count=1 ./internal/durabletest -run '^TestPhase3B2'
ALAN_PHASE2_EVIDENCE=/tmp/alan-phase3b3-prior-authority-race \
  go test -json -race -tags=integration -count=1 ./internal/durabletest -run '^TestPhase3B2'
staticcheck -tags=integration ./research/twoagent/... ./research/mutualobligations/... ./research/authorityservice/... ./research/delegationcorrespondence/... ./examples/communication/authority ./examples/communication/mutual
staticcheck ./internal/durabletest
staticcheck -tags=integration ./internal/durabletest
```

The last command has six pre-existing diagnostics, retained with the unchanged
base-input comparison in `evidence/research-baseline`. It is an expected nonzero
lint result, not a green gate. Compare the exact diagnostics rather than accepting
any nonzero result. Ordinary root-module tests do not include integration-only
checks or separately declared consumer modules. Count package-level no-test
events separately from named skipped tests.

Build and compare the unchanged deterministic model exports:

```sh
go build -buildvcs=false -o /tmp/alan-3a-export ./research/twoagent/cmd/export
/tmp/alan-3a-export > /tmp/alan-3a-model.json
cmp docs/foundational/phase3a/evidence/model.json /tmp/alan-3a-model.json
go build -buildvcs=false -o /tmp/alan-3b1-export ./research/mutualobligations/cmd/export
/tmp/alan-3b1-export > /tmp/alan-3b1-model.json
cmp docs/foundational/phase3b/evidence/model.json /tmp/alan-3b1-model.json
go build -buildvcs=false -o /tmp/alan-3b2-export ./research/authorityservice/cmd/export
/tmp/alan-3b2-export > /tmp/alan-3b2-model.json
gzip -dc docs/foundational/phase3b2/evidence/model.json.gz > /tmp/alan-3b2-archived-model.json
cmp /tmp/alan-3b2-model.json /tmp/alan-3b2-archived-model.json
go build -buildvcs=false -o /tmp/alan-3b2-audit ./research/authorityservice/cmd/audit
/tmp/alan-3b2-audit > /tmp/alan-3b2-audit.json
gzip -dc docs/foundational/phase3b2/evidence/audit.json.gz > /tmp/alan-3b2-archived-audit.json
cmp /tmp/alan-3b2-audit.json /tmp/alan-3b2-archived-audit.json
```

These comparisons retain the sixteen authority configurations, prior 33 history
controls, four ancestry inputs and prior custody/service findings. They do not
enumerate infinitely many cyclic histories or qualify the runtime by themselves.
Expected-failure controls remain documented in the earlier phase runbooks.
For example, `-config S2/LateGrant -assert-universal-completion` on the authority
export must exit 1, as must the S3/Missing and Cartesian/ScopeMix authority-safety
assertions. Exit 2 is an apparatus/argument failure. Positive S2/Explicit with
both assertions must exit zero. Never replace a retained RED with a green test
that silently changes its claim.

## Independent raw checking, projection and progress

The research package does not call runtime authorization or Replay. Its raw
oracle operates on exported JSON before the separate model projection. The
final frozen runtime corpus is the input, including the five malformed copies
and the actual selection-order counterexample. The model generator and progress
query use the unchanged earlier model APIs.

```sh
go build -buildvcs=false -o /tmp/alan-3b3-check ./research/delegationcorrespondence/cmd/check
go build -buildvcs=false -o /tmp/alan-3b3-corpus ./research/delegationcorrespondence/cmd/corpus
go build -buildvcs=false -o /tmp/alan-3b3-progress ./research/delegationcorrespondence/cmd/progress
go build -buildvcs=false -o /tmp/alan-3b3-mutations ./research/delegationcorrespondence/cmd/mutations
/tmp/alan-3b3-corpus > /tmp/alan-3b3-model-corpus.json
/tmp/alan-3b3-progress > /tmp/alan-3b3-progress.json
/tmp/alan-3b3-mutations research/delegationcorrespondence/testdata/positive.json \
  > /tmp/alan-3b3-raw-mutations.json
/tmp/alan-3b3-check -root /tmp/alan-phase3b3-runtime-frozen/normal \
  > /tmp/alan-3b3-normal-unqualified.json
/tmp/alan-3b3-check -root /tmp/alan-phase3b3-runtime-frozen/normal \
  -expect docs/foundational/phase3b3/evidence/task3-development/expectations.json \
  > /tmp/alan-3b3-normal-expected.json
```

The unqualified full-corpus command must exit 1, retaining all raw failures and
the named mismatch. The explicit 48-case expectation map gives exit 0 only when
the actual findings match it. Parsed JSON from the two invocations must be
identical. Repeat both commands for `race`. Exit 2 is an argument, input-access,
outer-encoding or digest apparatus failure; it is not a semantic falsifier.
Use the built CLI when asserting exact exit codes because `go run` wraps them.

For a relocated archive, pass `-root /absolute/extracted/runtime-frozen/normal`
and `-original-prefix /tmp/alan-phase3b3-runtime-frozen/normal`, replacing both
suffixes with `race` for that corpus. The loader verifies digests before policy
checking and uses those verified bytes. It neither opens SQLite nor repairs
the retained evidence. The positive fixture in the package is byte-identical
to the frozen ordinary S2 final export.

`ALAN_TASK3_RAW` and `ALAN_TASK3_PROGRESS=1` are necessary to execute the full
external-corpus and exhaustive-progress tests; an ordinary invocation without
them deliberately skips those tests. The standalone CLIs separately retain the
22-case/eight-class model corpus, all fifteen raw mutations and full progress
results. Development/review falsifiers remain archived separately from final
committed-source qualification. The whole-Store crash projection is unsupported
even when a final service subsequence passes its ledger.

The mutation generator requires the complete S2 positive parent used above,
including its grant, effects and receipts. Raw validity alone does not imply
those rows exist. A final review found that giving it the raw-valid cut-01
initial snapshot currently exits 2 with an index panic rather than a clear
precondition error. This diagnostic limitation is retained and deferred; use
`check -input` to check arbitrary captured prefixes. It does not affect the
fifteen mutations generated from the qualified complete parent.

## Public runtime qualification

From a clean runtime checkout of the executable freeze, build the historical old
reader from `af6943a5805a9a6c0e8fe77f804264f4765d07ec` in a separate checkout.
In that old checkout run `go build -buildvcs=false -o /tmp/alan-phase3b3-old-reader-rebuilt ./cmd/alan`, then set `ALAN_PHASE3B3_OLD_READER` to the resulting
absolute binary path in the current runtime test environment. Do not run the
new source while claiming old-reader compatibility. The runtime's
`docs/qualification/phase3b3/README.md` and `task4-frozen/README.md` contain the
old-source/overlay instructions and exact final command records.

```sh
go test -json -count=1 ./...
go test -json -tags=integration -count=1 ./internal/executor ./runtime ./author ./examples/delegation
go test -json -race -tags=integration -count=1 ./internal/executor ./runtime ./author -run '^TestRequestAuthority'
ALAN_PHASE3B3_EVIDENCE=/tmp/alan-phase3b3-reciprocal-normal \
  go test -json -tags=integration -count=1 ./internal/durabletest \
  -run '^TestPhase3B3(Reciprocal|CartesianOperations|GenericUnknown|ActualResultLaundering|SerialSerial|AppliedUnknownRevoked|RetainedMutations)'
ALAN_PHASE3B3_EVIDENCE=/tmp/alan-phase3b3-reciprocal-race \
  go test -json -race -tags=integration -count=1 ./internal/durabletest \
  -run '^TestPhase3B3(Reciprocal|CartesianOperations|GenericUnknown|ActualResultLaundering|SerialSerial|AppliedUnknownRevoked|RetainedMutations)'
ALAN_PHASE3B3_EVIDENCE=/tmp/alan-phase3b3-baseline-current \
  go test -json -tags=integration -count=1 ./internal/durabletest -run '^TestPhase3B3Baseline'
staticcheck ./...
staticcheck -tags=integration ./...
```

The ordinary and race reciprocal runs each produce 48 indexed cases and 53
passing test/subtest events. Ten cut labels produce eleven SIGKILL witnesses;
04/05 share a hook, and 06 has two distinct subcases. The current baseline
command deliberately skips its opt-in S2-completion expected failure. Reproduce
that actual RED on baseline source `8319a9067c16aaf5e3ecccb3d0c62f5c536e6140`
with the command/environment in `baseline-frozen/commands.json` inside the
historical archive. The endpoint's stable-key idempotency must not be relabeled
as one contact, and a generic unknown must not be relabeled as a receipt.

## Preserved evidence and relocation

Research `evidence/research-frozen/evidence.tar.gz` contains 1,375 members under
`research-frozen/`. Its wrapper manifest records every member and the eight
omitted rebuildable binaries. `freeze.py` is the actual qualification controller;
its historical absolute paths and source guard must be deliberately adapted for
a fresh reproduction, without changing the source identity or expected results.
`commands.json` contains all 42 exact commands, environments, statuses and log
digests. The verification summary distinguishes passing tests, expected semantic
failures, unchanged lint diagnostics and apparatus failures preserved elsewhere.

Runtime `task4-frozen/evidence.tar.gz` has 3,784 members under `runtime-frozen/`.
Verify its SHA256 and every member against the adjacent wrapper manifest before
using it. The ordinary/race `case-index.json` files refer to historical absolute
case directories under `/tmp/alan-phase3b3-runtime-frozen/{normal,race}`. Keep
these digest-bound files unchanged. Supply an explicit original-prefix to
extracted-root mapping to the research loader. Relocation must not change
logical RunIDs, original request IDs, service operation IDs or root identities.

For malformed copies, preserve raw tables even when Replay/Inspect/Open reject.
A policy checker must derive the missing or changed authority edge from source,
permissions, original history and actual prior rows, not from a case's expected
classification. Missing-grant metadata with no original bytes requires joining
the unmutated source case. Runtime policy validity, model correspondence,
unsupported projection domain and runtime reader disposition are separate
columns of evidence.

Never open archived SQLite originals, including read-only opens. Copy every
complete database/WAL/SHM set to disposable storage before using a database
reader. The baseline's historical SHM custody incident and regenerated frozen
run are retained separately. Later byte-level archive checks open no SQLite.
Online post-Open snapshots precede Run and are coherent backups, not arbitrary
copies of active database files.

The packet preserves a model/runtime disagreement: an authentic earlier BA
selection, accepted exact grant and intervening ready selection lead to a
runtime refusal where the unchanged LateGrant model permits authorization.
An expected-refusal test passing does not establish model agreement. Likewise,
separate grant with revocation and applied-without-local-receipt/duplicate-contact
histories must not acquire invented atomic model steps. Read the supported-domain
and stutter records alongside each projected trace.

## Manuscript

From `paper/`, build with `tectonic --keep-logs alan-dissertation.tex`, render the
changed pages with Poppler, and inspect them. Normalize the edited mathematical
paragraphs in Markdown and TeX and compare them. This is a scoped consistency
check; historical sections outside the slice are not silently declared aligned.
No new literature attribution is introduced, and `references.bib` remains
unchanged. Preserve the existing warnings separately from any new layout defect.
