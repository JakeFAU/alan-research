# Final research qualification

Executable source: `7894bfae268e6dc0f5418bda3c6022f3a8a55540`.
The public-runtime inputs are the ordinary/race archives qualified at
`759df6e3d2505f27e1287b2cfeb07bc4290afc0a`; this directory does not replace
their independent endpoint or Store evidence.

`commands.json` records all 42 commands, exact environment additions, exit
statuses and log hashes. Fourteen deliberately return nonzero: seven prior
model RED controls, two unqualified corpus reports, four killed checker-source
mutations and the six pre-existing integration Staticcheck diagnostics. Every
actual status matched its declaration. Read `verification-summary.json` for
counts and findings, `completion.json` for the source/lint conservation checks,
and the command logs for the underlying evidence.

`manifest.json` binds the archive and all 1,375 members under `research-frozen/`.
The archive contains the exact controller, source/fixture hashes, prior-model
comparisons, raw reports, full progress queries, generated corpus, mutations,
test logs, emitted prior regression evidence and checker-mutant source copies.
Eight rebuildable binaries are omitted; their hashes, exact build commands and
source identity remain recorded. A rebuilt binary at another path is not
automatically byte-identical. The source-guard preflight rejected a controller
SHA typo before running qualification; its record is retained as an apparatus
error, not a semantic falsifier.

Verify SHA256 for the archive and every listed member before extraction or use.
Do not open archived SQLite files, even read-only; use raw JSON or disposable
complete DB/WAL/SHM copies. The wrapper was verified without opening SQLite.
Preserve original absolute paths in digest-bound runtime indices and pass the
checker's explicit `-original-prefix` relocation mapping. The [runbook](../../05-runbook.md)
describes reconstruction and the [results](../../04-results.md) separate raw
validity, supported projection, the retained model/runtime mismatch and unknown
outcomes. Expected-outcome files alter exit status only, never those verdicts.
