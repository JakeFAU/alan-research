# Task 3 repair round 1 — fbcedaee..a49ee5b0

## Finding verdict

**R1 / P2 — ADDRESSED.** `docs/foundational/phase3c5/tools/build_pdf.py:13–22,25–53` and `docs/foundational/phase3c5/tools/render_pdf.py:12–21,29–64` allocate the caller-supplied output with exclusive `mkdir()` before capture writes. Existing files, directories and dangling symlinks reach the refusal path. Build artifacts, extracted text, commands and renders are written beneath the newly allocated directories. `render_pdf.py:43–46,65` explicitly records that visual inspection was not performed. Both historical entrypoints under `evidence/manuscript/`, line 2, refuse immediately.

`docs/foundational/phase3c5/tools/verify.py:38–46,59–73,171–191` restricts historical resolution to the original manuscript and qualification manifests and the four declared source paths. The verifier compares the snapshots to their original Git blobs and verifies current repair files through the separate repair manifest. The four inert snapshots also match the old-side blob IDs in the supplied patch. The delivered PDF, manuscript sources and original capture records are outside the repair diff, and their live hashes match the protected before/after records.

## Checks

- Parsed all 14 repair JSON files and independently checked 31 repair artifact identities, the original manifest bindings, and 779 protected hashes. Cross-checked all eight compressed manuscript archives against the render manifest and their decompressed bytes against the build manifest. Cross-checked the 17 build and 67 render identities against the saved output rosters and manifest hashes.
- Archived commands show two fresh successes and eight exit-1 refusals, including the retired paths and existing directory, file and dangling-symlink controls for both active helpers. All eight archived inner build commands and 54 page-render commands report exit 0. The fresh captures record no visual acceptance.
- This review checked source and archived evidence. It did not run the build, render, model suites, or the full repository verifier.

## New breakage in the fix diff

None found in scope.

## Out-of-scope observations

R2 (sparse Contents continuation) and the seven inherited warnings remain deferred.

## Verdict

**APPROVE — R1 addressed.** No new blocking finding in the repair diff.
