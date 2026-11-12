# SDD ledger — plan: docs/foundational/phase3b5/phase3b5-plan.md

## Objective and verified setup

Continue merged Phase 3B.4 with separately retained grant plus exact-root revocation, through reviewed draft research PR. Contract: docs/foundational/phase3b5/01-contract.md. Branch codex/phase3b5-grant-revocation; research base4fea6c31; runtime committed mainc72c396d/tree52f4d7c9. Existing clean linked worktree intentionally reused. Prior PR45/13 merged; no runtime source change planned. Preflight provenance and read-only recon archived under evidence/preflight. Baseline targeted JSON has four tested-package passes and no failures; optional evidence tests skip; expired process handle loses direct process exit, explicitly recorded.

## Rulings (chronological)

Ruling: Execute the narrowly proposed next question through draft PR without another design/plan permission prompt — the user explicitly asked to continue the authorized research and AGENTS supplies execution/finishing preferences; developer instructions prioritize carrying out that scope — if the chosen slice differs from intent, research-only isolated commits are reviewable and reversible.
Ruling: Reuse the clean linked worktree with a new branch from verified merged main — current worktree already provides isolation, and no cleanup of unrelated work is needed — if path naming confuses later work, manifests and branch identity must disambiguate it.
Ruling: Model exact-root revocation with distinct b and c complete tuples, not grant-ID or whole-composition revocation — the owning runtime API and sealed N07 distinguish those roots — if the abstraction loses a relevant property, the scoped correspondence must be narrowed or rejected.
Ruling: Add one opt-in authorityservice scenario and preserve old sixteen exports, rather than duplicate a model package — existing state/event shapes express the finite question with less machinery — compatibility costs four archived byte regressions and independent config-specific oracle support.
Ruling: Reanalyze sealed normal/race cases without rerunning unchanged production runtime — runtime tree equivalence and immutable corpus are available, and the research gap is the model/relation — this establishes conditional analysis of historical observations, not fresh runtime execution or broader schedules.

## Preflight conflict scan

| Tasks | Shared file/interface or internal consistency | Finding and resolution |
|---|---|---|
| 1 | model/Explore/export tests versus oracle scope | New config is operational-only until Task2; Task1 does not claim independent audit. Old tests remain valid. |
| 2 | independent ledger/progress/compiled mutants | Ledger independently specifies bindings; progress uses public successors but reverse reachability and audits are separate evidence. No universal-history claim. |
| 3 | per-capture mapping versus final physical evidence | Intermediate cuts retain opportunity without Apply; final dossier supports endpoint claims only at declared boundary. Missing physical evidence cannot imply zero contact. |
| 4 | qualified source versus later paper/review edits | Freeze executable source; documentation-only commits do not change qualified inputs; any code fix reruns affected gates. |
| 1–2 | GrantRevocationConfig, tuple/root, Event shape | Exact name fixed; new root configuration-specific in both implementations; default old exports remain unchanged. |
| 1–3 | Successors and explicit BA/c mapping | Fixed injective whole-root mapping; actual events selected uniquely and complete descriptors compared; model era differs from runtime holder revision. |
| 1–4 | finite counts and historical claims | Only actual graph results qualify; prior historical sections retained, no predicted counts. |
| 2–3 | AuditTrace/new config, QueryProgress | New projection consumes independently audited event vocabulary; progress is not inferred from one runtime trace. |
| 2–4 | audit/progress results and proof sketch | Explicit reachable-state quantification and existential scheduling; no runtime liveness theorem. |
| 3–4 | raw/physical/support verdicts and corpus provenance | Separate verdicts and original producer identity preserved; no runtime execution claim. |

## Task status

Task1: in progress; BASE8971eb11. Implementer /root/grant_revocation_model. Baseline log gzip committed before dispatch.
Task2: pending.
Task3: pending.
Task4: pending.

Ruling: Permit extraction of declared binary members for hash verification without opening SQLite — the existing independent directory checker hashes every declared member, so JSON-only extraction contradicts its integrity contract — the cost is disposable disk space; no database inspection or mutation is authorized. Task3 plan corrected before dispatch.

Corpus extraction complete: six archive digests and 10,634 declared regular members verified; raw normal/race, v1 normal/race and promoted normal/race under /tmp/alan-phase3b5/corpus. JSON read and binary hash only; no SQLite opens. Manifest evidence/preflight/corpus-extraction.json.

Old corpus apparatus: newer v1 archive contains54 cases, so original48 regression corpus extracted at its unchanged absolute prefix from runtime docs/qualification/phase3b3/task4-frozen/evidence.tar.gz. Archive SHA0c8c48f2, all3784 members verified; normal/race48 each. See evidence/preflight/old48-extraction.json. No byte rewrites or SQLite opens.

Task1 operational preliminary report (not independent result): 229834 states/1807246 arcs; 203833 safe-completable,26001 safe noncompletable; old default authority export byte match reported. Final evidence/report awaited. Staticcheck requires STATICCHECK_CACHE explicitly on Darwin; use /tmp/alan-phase3b5/staticcheck in subsequent commands.

Task1: implementation2107c44c complete; report and archived model/logs read by controller, byte hash and actual graph counts confirmed. Independent reviewer /root/review_grant_revocation_model; package review-8971eb11..2107c44c.diff. Gate pending.

Task1: complete (commits8971eb11..2107c44c, review clean). Spec compliant/quality approved by /root/review_grant_revocation_model; no findings. Task2 BASEe6e53deb; pending dispatch.

Task1 reviewer limitation resolved by scope: independent ledger/progress owns Task2, concrete relation Task3, all-four regression Task2, manuscript/full qualification Task4; none was a missing Task1 requirement. Task2 implementer /root/grant_revocation_oracle, BASEe6e53deb, task-2-brief.md.

Task2 preliminary: independent audit counts229834 prefixes/1807246 arcs match operational graph; reverse reachability/SCC agreement reported. A witness-selection apparatus error selected a merged target shortest prefix lacking the desired earlier permit; retained RED and corrected to qualifying actual arcs. Await committed report and review; no result promoted yet.

Task2 implementation0e8d8f87, BASEe6e53deb; report and source-hash/tests/export/mutant evidence inspected. Review pending /root/review_grant_revocation_oracle with review-e6e53deb..0e8d8f87.diff. New predicates/SCC/audit match; all4 compiled mutants assertion-killed, all4 archived exports and baseline-source old progress byte-equal. Preserved RED and wrapper/cache apparatus failures remain explicit.

Task2: complete (commitse6e53deb..0e8d8f87, review clean). Reviewer /root/review_grant_revocation_oracle: spec/quality approved, no findings; finite-model boundary explicitly retained. Concrete relation and final qualification owned by later tasks.

Task3: in progress; BASE47d4e4a9; implementer /root/grant_revocation_correspondence. Brief task-3-brief.md plus prepared corpus context. Tasks1/2 approved; production runtime unchanged. Task2 reviewer later-task limitations mapped to Task3/4, no unresolved task gap.

Task3 preliminary: both-mode selected cases/projected cuts GREEN and genuine unknown-after-loss excluded; adversarial mapping/CLI/determinism/legacy work ongoing. Controller qualification plan prepared at /tmp/alan-phase3b5/qualification-plan.json; NOT RUN before Task3 review/freeze. Prior inspected full ordinary duration1517s, so use repository60m timeout with4200s outer watchdog; parallel docs/review only after source freeze.

Task3 implementation331186c7, BASE47d4e4a9; report, all63 input hashes, GREEN/reproduction/legacy attack/empty lint logs inspected. Gate pending reviewer /root/review_grant_revocation_correspondence, package review-47d4e4a9..331186c7.diff. Eight supported final cases/twenty captures across two modes; two unknown-after-loss exclusions retain no events. Legacy10 dossiers byte-equal. Controller first JSON summary assumed mutant list but file is object; inspection helper raised AttributeError after source/log checks and was corrected. No qualification command was affected. Full final qualification NOT YET RUN.

Task3 review rejected I1 (Important): []b.Obligation serializes as base64AQ==; compare_reports.py read four characters as four BA opportunities. Core state/trace hasone. Fix round1 begins from331186c7; original implementer resumes; task-3-fix-1.md carries verbatim finding. No full qualification started.

Ruling: Require endpoint ready records as well as journals/snapshot — they bind endpoint labels and complete tuples, and reviewed sealed cases provide them — dossiers lacking ready records remain outside this relation.
Ruling: Keep excluded dossiers as capture locators with no model events — no chronological loss/unknown relation has been established — supported prefixes inside excluded dossiers require a later separately scoped analysis.
Ruling: Use a raw-valid intermediate cut for the wrong-root attack — a complete changed N07 history would also need a changed later AB decision — the mutant result qualifies that mapping predicate, not an altered full runtime campaign.

Task3 scope rulings accepted by reviewer; no semantic-chain defect beyond I1 evidence summary. Fix round1 in progress /root/grant_revocation_correspondence, BASE331186c7. Explicit AB/BA report labels and serialized cardinality checks chosen; old model State/Event JSON unchanged. Whole-branch qualification plan remains pending.

Task3: complete (commits47d4e4a9..4ad1f9ef, review clean after fix round1/5). Scoped rereviewer /root/rereview_revocation_i1 verified I1 addressed, exact serialized arrays, historical hashes, no new breakage. Report task-3-rereview-1.md.

Task4: in progress. Executable source frozen at4ad1f9ef37c57375d4a39eafc9a5fcd91bc0425d; controller starts full qualification, manuscript worker owns docs/source only. Research origin/main freshly fetched remains4fea6c31; no drift.

Task4 manuscript worker /root/grant_revocation_manuscript committed4cabef20 (BASE4ad1f9ef). Report task-4-report.md. Narrative reviewed inputs and nine closed qualification records; full ordinary explicitly pending. PDF marker executed once successfully before authoring; controller must not repeat. Task reviewer /root/review_revocation_manuscript, package review-4ad1f9ef..4cabef20.diff. Controller-owned test/archive/PDF/delivery gates remain open; task not complete.

Qualification partial: affected race exit0 in247.654s,422 passing test/subtest events,7 tested packages/3 no-test packages,zero named skips; scoped lint/build/module/model/exports exit0; global lint13/31 diagnostics exactly byte-equal to baseline. Full ordinary remains running, no claimed exit. Source582 files unchanged.

Task4 substantive review clean at4cabef20: task-4-review.md. All cannot-verify items are controller gates, still open as applicable. PDF then built exit0 from reviewed source (before whole-branch review),42pages SHA83bfaa590ad6b3cbd61ecf32043efdce93ed7ac75805d4feee1e625911de5bad. Viewed pages1,2,9,28,29,38–42, no new layout defect. Five inherited warning patterns/source ranges unchanged; details archived evidence/task4/pdf. Any subsequent source correction requires rebuild.

Final qualification closed: runner process exit0; all10 records matched. Ordinary exit0,2633 test/subtest pass events,51 tested packages/13 no-test packages,3 named optional artifact-export skips; race422 events,7 tested/3 no-test packages,0 named skips. Final verifier exit0 checked all582 source hashes, every closed log hash/size, formats, old4 exports and Task1/2 output byte equality. Evidence copied to task4/qualification and verification.json. Current manuscript/PDF unchanged from reviewed source.

Task4: complete (commits4ad1f9ef..836cbb1f, review clean for substantive source; controller verification/artifact gates satisfied). Final ordinary/skip/log/source and PDF requirements are resolved by the retained actual command records and visual inspection. Worker factual closure836cbb1f checked all10 closed logs and hashes. Whole-branch review, final archive closure and PR delivery are branch-finish gates; no implementation findings parked.
