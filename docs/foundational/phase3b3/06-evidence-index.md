# Phase 3B.3 evidence map

Executable source, development candidates, reviews and delivery records have
different roles. A successful expected-outcome command does not turn the named
model mismatch into agreement. Follow the [results](04-results.md) for the claim
and [runbook](05-runbook.md) for the commands.

## Source boundaries

| Boundary | Identity | Role |
| --- | --- | --- |
| Research base | `8197683c23a3ca774750fc3c2c389c5a64dd38b8` | Verified merged Phase 3B.2 state; 563 original Go/module inputs remain unchanged |
| Public runtime base | `af6943a5805a9a6c0e8fe77f804264f4765d07ec` | Existing-API obstruction and frozen old reader |
| Baseline fixture | `8319a9067c16aaf5e3ecccb3d0c62f5c536e6140` | Genuine negative S2 attempt and positive S1 control |
| Contract v3 | `360cae66f5ffdd6e7a95a568552df5522122da21` | Final normative authority contract, SHA256 `f19d9c13f84c17496fabef91cb7575cf52972ff228c0b69e063cdc1c99100e12` |
| Reviewed runtime primitive | `4f6f5f2500a48b2ca01ee7b09a823fa63656b7ab` | Corrected authority implementation and compatibility checks |
| Final runtime executable | `759df6e3d2505f27e1287b2cfeb07bc4290afc0a` | Reciprocal/cut fixtures; original final 48-case ordinary/race evidence |
| Final research executable | `7894bfae268e6dc0f5418bda3c6022f3a8a55540` | Corrected independent checker, corpus, queries and tests |
| Supplemental integration | main `76031c34` + task `c7fe6d8`, tree `bdbbae5f89f0c76d48618b3d8891dc40199ca4ed` | Prospective integration only; neither original branch was merged |

Evidence-only delivery commits do not replace these executable identities. The
draft PR heads are delivery identities; qualification manifests bind the actual
source and commands used.

## Research repository

All paths below are relative to this directory.

| Artifact | What it establishes or preserves |
| --- | --- |
| [00-brief.md](00-brief.md), [01-runtime-gap.md](01-runtime-gap.md), [02-contract.md](02-contract.md) | Original packet, verified public-API gap and frozen contract |
| [Rejected contract v1](evidence/rejected-contract-v1/implementation-obstruction.md) | Namespace-dependent literal roots cannot be supplied before Create in the proposed encoding |
| [Rejected contract v2](evidence/rejected-contract-v2/) | Independent legacy compatibility, typed inspection and original-request gate findings |
| [Selection narrowing](evidence/selection-narrowing/review.md) | Actual retained authorization refusal versus the unchanged model's permitted edge |
| [Progress design](evidence/progress-design/review.md) | Independent exhaustive LateGrant/Revocable query and exact predicate boundaries |
| [Research lint baseline](evidence/research-baseline/command.json) | Six pre-existing integration Staticcheck diagnostics, separately retained |
| [Checker development](evidence/task3-development/manifest.json) | Rejected 22-file checker candidate, initial development failures, parent subject/causal-prefix controls and candidate commands |
| [Independent raw REJECT](evidence/task3-rejected-review/task-3-raw-review.md) | Route, immutable-source continuity and owner-epoch false acceptances with exact raw probe inputs/results |
| [Projection/progress review](evidence/task3-rejected-review/task-3-progress-review.md) | Independent approval of the unchanged model/corpus/projection slice, separately from raw-policy rejection |
| [Checker correction](evidence/task3-round1/review.md) | Scoped re-review approval; 23-file corrected candidate, actual RED/GREEN evidence, fifteen mutations and explicit missing-initial-capture limits |
| [Final research execution](evidence/research-frozen/README.md) | Frozen committed source, all 42 declared outcomes, unchanged old exports, final raw/model/progress verdicts and four killed checker-source mutations |
| [Manuscript interim](evidence/manuscript-interim/review.md) | Earlier scoped manuscript/runtime-reproduction review; not final checker qualification |
| [Manuscript final](evidence/manuscript-final/manifest.json) | Exact final manuscript/PDF, build logs, local MD/TeX alignment, unchanged BibTeX and inspected pages |
| [Whole-change review](evidence/final-review/review.md) | Independent specification/quality/overall APPROVE for the bounded slice; exact minor CLI probe, reviewed document hashes and explicit unproved domains |

Candidate runs above must not be cited as the final committed-source execution.
Review approval is distinct from draft delivery, merge status and proof.

## Runtime repository

The [public runtime repository](https://github.com/JakeFAU/alan) retains its source
and physical evidence under `docs/qualification/phase3b3/`:

| Path | Role |
| --- | --- |
| `baseline-frozen.tar.gz`, `baseline-manifest.json` | Regenerated frozen ordinary/race baseline and explicit S2-completion RED |
| `baseline-development.tar.gz`, `baseline-review.tar.gz` | Development and review history, including the earlier SQLite SHM custody incident |
| `rejected-literal-root/` | Actual construction failure for the first root representation |
| `rejected-candidate-v2/`, `task2-completion/`, `round1/` | Rejected primitive, its concrete failures and independently re-reviewed correction |
| `legacy-matrix/` | Ten old-created formats 3–12, 80 primary calls and 60 byte-equal paired exports; finite coverage only |
| `primitive-provenance-review.md` | Independent source/archive reconstruction check |
| `task4-development/` | Original process-fixture development, apparatus failures and the first selection-narrowing input |
| `task4-frozen/` | Final executable identity, eight actual commands and 3,784-member raw archive |
| `task4-reviews/` | Independent recovery/admission and binding/custody source/raw-witness reviews |
| `current-main/` | Prospective integration with advanced main, initial Git-provenance export failure and corrected same-tree qualification |

## Custody and reconstruction

Each archive has an adjacent manifest. Verify the archive and every member
before using it; a filename containing "green" is not evidence of success.
Original command exit status and the associated claim determine that. Omitted
rebuildable binaries retain their hashes, source and build commands. A rebuilt
binary at another source path is not automatically byte-identical.

Never open an archived SQLite original, even read-only. A reader can change SHM
bytes. Use raw JSON/byte checks or a disposable complete DB/WAL/SHM copy. The
baseline custody incident remains distinct from its regenerated qualification.
Later archive checks did not open SQLite.

For relocated runtime evidence, keep digest-bound absolute paths unchanged and
use the checker's explicit original-prefix mapping. Missing initial captures,
unknown external outcomes, unsupported projections and rejected histories remain
visible. None may be repaired by inventing a receipt, permission or model event.
