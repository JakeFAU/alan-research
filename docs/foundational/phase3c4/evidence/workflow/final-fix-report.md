# Phase 3C.4 final review clarification

Status: DONE. Reviewed base `7c90a10adba46433653de8509fa7f68333127b6b`; fix commit `0bc7207efd8bd31c65a4ea18b04fb6fb50583466` (Clarify Phase 3C.4 archive reproduction).

The final review's nonblocking P3 applies to the historical Task 1 driver: its destination is the script's parent directory and its writes can replace historical records and captures before a run succeeds or fails. The controller elected to preserve the archived executable and evidence. I added `evidence/task1/README.md` to classify the driver and warn against execution in place. I expanded `05-qualification.md` to distinguish read-only archive checks, new Task 2 checker/export/CLI captures, and separate scoped Go ordinary/race/Staticcheck commands. Verified from source that `tools/capture.py` creates the named attempt with `attempt.mkdir()` before writing results and uses a second scratch directory. The guidance names the scratch parent prerequisite, requires a fresh name in both locations, and discloses the saved compiler enumeration's absolute checkout paths. The archived Task 1 driver still has no overwrite protection.

The current `evidence/qualification/artifact-manifest.json` changes precisely one prior binding (the qualification narrative) and adds one (the Task 1 README). The 5,422-byte pre-change artifact manifest is retained in `evidence/final-review-fix/before-artifact-manifest.json`, identical to the reviewed Git blob; SHA-256 `6c32d2212131aaf4898f9cc3fc368742d10b941945bb90a115aaceb0c3280ffc`. The same directory's `result.json` and `verification.json` retain the actual fresh invocation and exact stdout, separately from historical Task 3 verification. No earlier command record or result was relabeled or changed.

Verification performed:

- `python3 docs/foundational/phase3c4/tools/verify.py` exited 0 with empty stderr and status PASS. It checked 110 accepted Task 1/2 files, 57 Task 2 archive entries, manuscript alignment, and 15 local links. Machine comparison confirmed the output equals the historical verification result except for the two added links.
- An inline Python identity check compared the before-manifest with `git show` at the reviewed base, checked the precise binding delta and refreshed hashes, and verified the capture record's three file bindings.
- Inspected complete staged status, diff stat and diff before commit; `git diff --cached --check` and `git diff --check HEAD^ HEAD` passed. Post-commit path checks showed exactly the six authorized paths and no change to any existing Task 1/2 file, tool, manuscript, manuscript artifact, or research source.
- Did not run the historical driver, fresh Task 2 capture, audit CLI, Go test/race suites, Staticcheck, model, runtime, SQLite, or PDF builds/renders. This was a documentation and evidence identity change; no new test or semantic result is claimed.

Self-review found no out-of-scope changes. The controller's concurrent uncommitted `06-reviews.md` edit was kept out of staging and the commit. No push or merge occurred.
