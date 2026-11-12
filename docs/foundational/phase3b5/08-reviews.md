# Phase 3B.5 reviews and rulings

| Task | Reviewed source | Disposition |
| --- | --- | --- |
| 1. Operational model and witnesses | `2107c44c16eb4a8cb7f67122d734cd995cb953ec` | Independent spec and quality review approved. [Review](evidence/workflow/task-1-review.md). |
| 2. Independent authority and progress audit | `0e8d8f874b5993b62692aee1b872ef7619a93422` | Independent spec and quality review approved. [Review](evidence/workflow/task-2-review.md). |
| 3. Concrete forward relation | `331186c7b44cdc18b9f7a0970992fe00ce89714b` and fix `4ad1f9ef37c57375d4a39eafc9a5fcd91bc0425d` | Initial review found I1. Scoped rereview approved the fix. [Initial review](evidence/task3/review-i1/review.md) and [rereview](evidence/workflow/task-3-rereview-1.md). |
| 4. Qualification and manuscript reconciliation | `4cabef20bf48453cba9dad5ac03a7d9e54a19d77` with executable inputs frozen at `4ad1f9ef37c57375d4a39eafc9a5fcd91bc0425d` | Independent task review approved the narrative and manuscript sources with no findings. Frozen-source qualification and controller PDF inspection are recorded. [Task review](evidence/workflow/task-4-review.md) and [qualification](06-qualification.md). |

I1 concerned serialized evidence. The initial Go report encoded a single unresolved BA as the byte string `"AQ=="`. The summary then counted its four characters as four BA entries. The corrected report emits direction labels. In each sealed mode the N07 cuts now report `["BA"]`, `["BA"]`, and `[]`. A direct comparison of all ten dossiers found every decoded field outside `Unresolved` unchanged. The historical artifact and red/green assertions preserve the failure and its correction. [I1 provenance and delta](evidence/task3/review-i1/).

## Material scope decisions

| Decision | Reason and cost |
| --- | --- |
| Model the exact BA authority root and a distinct AB root. | The sealed N07 suffix distinguishes two complete B-held roots. The relation requires the complete tuples, their endpoint bindings and the exact lifecycle target. |
| Add an opt-in configuration to the existing model. | The existing State/Event shapes express the finite question. Compatibility requires all four archived-export byte comparisons and a configuration-specific independent interpreter. |
| Require endpoint ready records in addition to both journals and the complete snapshot. | Ready records tie endpoint labels to the complete tuples. The relation excludes dossiers that lack these records. |
| Verify declared binary members in disposable extractions. | The raw directory checker hashes all declared members. Extraction uses temporary disk space and treats SQLite files as opaque bytes. |
| Reanalyze the sealed normal and race captures. | The research question concerns interpretation of those histories. The resulting evidence remains conditional on their declared scope and provides no additional runtime schedules. |
| Freeze executable inputs before the final qualification. | Manuscript, review and delivery artifacts can be edited against that freeze. A change to executable inputs requires affected qualification to be rerun. |
| Preserve the manuscript's historical phase boundaries. | The new abstract, subsection and conclusion are synchronized across Markdown and LaTeX. The Markdown-only Section 7 predates this task and remains a representation gap outside the edited sections. |

The workflow archive retains the task briefs, reports, reviews, decision ledger, and the review packages with member hashes. Whole-branch review and draft PR delivery remain pending. [Workflow records](evidence/workflow/).
