# CI and releases

GitHub Actions validates pull requests into `main`, each landed commit, and a
manually selected semantic version release. The workflows use Ubuntu 24.04,
CGO, Go 1.27.4, Staticcheck v0.8.2 and govulncheck v1.9.0. Actions are pinned to
commit SHAs. Go's module and build caches are reused. Every test invocation uses
`-count=1`. Ordinary CI uses controlled fixtures and local processes.

## Automatic checks

| Check | Scope | Job watchdog |
| --- | --- | --- |
| PR checks | Active-source formatting, release-script tests, root module verification/build/vet in ordinary and integration modes, scoped Staticcheck in both modes, smoke packages and three runtime integration tests | 15 minutes |
| Main (pr) | Same checks on the landed SHA | 15 minutes |
| Main (regression) | Seven runtime tests, conflict-receipt crash recovery, and two CLI crash-boundary tests, each normally and with race instrumentation | 30 minutes |
| Main (worktree) | Separate worktree consumer module verification, vet, ordinary and race tests | 15 minutes |
| Main (adaptive) | Separate adaptive consumer module verification, vet, ordinary tests, and one nested/capacity-two/paired race fixture | 30 minutes |
| Main (security) | govulncheck for the ordinary and integration root module and the two consumer modules | 20 minutes |
| Main (lint-report) | Repository-wide ordinary and integration Staticcheck diagnostics | 20 minutes |

`scripts/ci/check.py` contains the exact package lists, named tests, and per-test
watchdogs. The smoke packages are `schema`, `effects`, `compose`, `author`, and
`research/twoagent`. Named selectors are checked for missing tests before running.
The regression scope covers selected crash/replay/unknown-outcome behavior.
These bounded checks provide regression evidence for their selected inputs.

Root `./...` excludes nested modules. Automatic consumer coverage includes the
worktree and adaptive modules documented in the root README. Other consumer
modules and independent Python phase auditors remain in their phase runbooks.
Formatting excludes frozen Go evidence under `docs/`.

Staticcheck is blocking in `schema`, `effects`, `compose`, and
`research/twoagent`, in ordinary and integration modes. `author` remains in
the smoke tests and the complete advisory lint reports. Its ordinary
build currently reports the inherited unused `raceChildren` constant.
Prior qualification recorded
13 ordinary and 31 integration diagnostics globally at its recorded source.
The current complete reports are advisory, and their counts can differ.
Tool execution, loading, and compilation failures still fail the report job.
Each report is retained as JSON Lines and its diagnostic count appears in the
job summary. govulncheck findings and execution failures fail the security job.
govulncheck uses the current Go vulnerability database. Its Go analysis scope
does not establish coverage of vulnerabilities in bundled SQLite C code.

PR runs are cancelled when superseded. Every landed SHA gets its own Main
validation. Logs are retained for 14 days on PRs and 30 days on main.
Watchdogs indicate failure if exceeded. They are initial limits to calibrate
against hosted-runner results.

## Full qualification

In Actions, select **Main validation → Run workflow**. Leave `qualification`
as `none` for the automatic scopes, or explicitly choose `ordinary` or
`integration` to add the corresponding full root suite. Those suites retain
the README's `-timeout=60m` per-package watchdog and a 120-minute job watchdog.
They use a full working tree containing the frozen evidence.

Choose `adaptive-race` to run the complete adaptive consumer race profile with
its documented 60-minute per-package watchdog. Automatic CI runs all ordinary
adaptive consumer cases and the existing race subtest
`TestAdaptivePublicConsumerDirectNested/nested=true/capacity=2/paired=true`.
CI verifies that this exact subtest passed, so a renamed/missing case fails.
The complete local race profile's CLI cases alone took 654 seconds during CI
setup, which motivated keeping the full profile opt-in.

The complete integration race campaign retains its existing README command
and 600-minute per-package watchdog. Run that campaign in a suitable local
environment when needed. Automatic Main validation runs the bounded race
scope listed above. Main validation triggered manually can target a branch
for diagnosis. Release eligibility requires a successful **push-to-main** run.

For a local smoke run, install the pinned Staticcheck version, then run

```sh
CGO_ENABLED=1 GOTOOLCHAIN=local GOWORK=off \
  GOFLAGS='-mod=readonly -buildvcs=false' python3 scripts/ci/check.py pr
```

The other automatic scopes have the same command with `pr` replaced by
`regression`, `worktree`, `adaptive`, `security`, or `lint-report`.

## Publish a version

After Main validation succeeds, select **Release → Run workflow**, choose
`main`, and enter a version such as `v0.1.0` or `v0.2.0-rc.1`. The dispatch's
full SHA is fixed even if another commit subsequently lands on `main`.
The optional `commit` input accepts a full SHA from an earlier successful
push-to-main validation. Leave it empty to release the dispatch commit.

The release script requires the latest post-merge validation for that exact SHA
to be complete and successful. It creates a tag for that commit and publishes
a GitHub Release with generated notes, the full research SHA, and the validating
run and attempt. Prerelease versions are marked as prereleases. Use semantic
versions with a `v` prefix and optional prerelease identifiers. Build metadata
suffixes are excluded from this initial version format.

Publication is serialized. An existing release or a conflicting tag fails.
If publication fails after tag creation, rerunning the workflow can finish
using that same tag and SHA. A retry after main advances must dispatch at the
original SHA using the `commit` input, with the same version. Existing
published versions are never overwritten by the workflow.

The released source is the research tree. The standalone `alan` module has a
separate implementation and no Go module dependency on this repository. Runtime
documentation can pin research citations to the release's full commit SHA.
Research Go imports use the separately documented `go-module` distribution
snapshot because the complete tree exceeds Go's intermediate Git archive limit.

## Repository settings

Once the first PR check has run, require the **PR checks** status check on
`main`, with the branch up to date before merging. Enable GitHub's **Release
immutability** setting before publishing the first version to lock release
tags and assets. These settings are stored by GitHub separately from the
workflow files.

The first hosted run qualifies the selected Ubuntu runner/toolchain combination.
Existing macOS qualification remains attributed to its recorded source and
environment. This CI setup changes qualification automation and release identity.
The manuscript's mathematical claims and executable semantics are unaffected.
