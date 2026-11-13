# Task3 independent spec and quality review

Use task-3-brief.md, task-3-report.md and task-3-controller-report.md in this directory, plus the exact review package supplied by controller. Task3 base is6dc16ecf36441a99005a62eba4cee2418f5f5d72. Review actual capture/qualification and manuscript claims, including initial checker rejection and reviewed correction. Task1 and Task2 have separate reviews; the correction lies in this range so inspect how its qualification is reported. Final broad branch review follows separately.

## Binding constraints (verbatim)

- Deliver an independently reviewed draft research PR; do not merge or release.
- No production runtime, dependencies, automatic retry, or authority semantics change.
- Existing Phase 3B.6 artifacts and interpretations remain intact. Old checkers' pinned producer protocols remain unchanged.
- No later pass rewrites raw artifacts to create equality. Opening the archived evidence SQLite files is prohibited.
- No identity renaming, field deletion, table filtering or post-hoc reordering is permitted. Array order, multiplicity, numbers, strings and null/empty distinctions remain significant.
- Freeze producer/comparison sources and contract hashes before the qualification pairs. Preserve every attempted qualification run, including mismatch or apparatus failure.
- Work in the existing research worktree on `codex/phase3b7-matched-capture`, base `d8645d739445d4b09c0d38dca4454ca6297fdf22`; do not edit the runtime checkout or `.envrc`.
- Use `/tmp/alan-phase3b7/` for disposable runtime builds and tests. Offline Go environment: `GOTOOLCHAIN=local GOWORK=off GOPROXY=off GOSUMDB=off GOMODCACHE=/tmp/alan-phase3b5/modcache GOCACHE=/tmp/alan-phase3b7/go-build`.
- One implementation agent at a time; no worker subagents. Each task receives independent spec and quality review. Preserve actual commands, exits, source identities and failed attempts.


## Method

Read-only checkout except your new review report. No subagents. Read the supplied diff once in manageable chunks; use focused source/evidence reads for named risks or binary artifacts absent from diff, and record checks. Treat report claims as unverified. Do not rerun already-recorded suites or execute new captures; use a narrow probe only for a concrete unanswered doubt. Never open captured SQLite files. Examine raw tar members as bytes/JSON when needed.

Return explicit spec APPROVE/REJECT and quality APPROVE/REJECT, with evidence locations and severity for each finding. Qualification scope, post-capture classification repair, complete declared local observation, endpoint idempotence, causal protocol, raw custody and source-conditioned limits are central. Check the manuscript's MD/TeX implications and actual PDF inspection evidence. Unresolved general theory remains open. Preserve prior REJECT findings. Write task-3-review.md and return verdicts plus concise findings.
