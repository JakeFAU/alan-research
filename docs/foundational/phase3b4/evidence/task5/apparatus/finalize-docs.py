import hashlib,json,pathlib,shutil,subprocess
repo=pathlib.Path('/Users/jacob/.codex/worktrees/phase3b4-research/alan-research');runtime=pathlib.Path('/private/tmp/alan-phase3b3-runtime');out=repo/'docs/foundational/phase3b4/evidence/task5';work=pathlib.Path('/tmp/alan-task5');parent=pathlib.Path('/tmp/alan-phase3b4-parent-final-gates')
required=['research-ordinary','research-model-checker-race','runtime-ordinary','runtime-integration','runtime-integration-race','runtime-promoted-normal','runtime-promoted-race','research-prior-authority-normal','research-prior-authority-race','research-prior-reciprocal-normal','research-prior-reciprocal-race']
for name in required:
 d=json.loads((parent/(name+'.json')).read_text());assert d['exit']==0 and d['source_unchanged'],name
subprocess.run(['python3',str(work/'summarize-gates.py')],check=True)
gates=json.loads((work/'gate-summary.json').read_text());assert all(g['source_unchanged'] for g in gates)
assert all(g['exit']==0 or g['name'] in ['research-staticcheck','research-staticcheck-integration'] for g in gates)
shutil.copytree(parent,out/'parent-final-gates',dirs_exist_ok=True);shutil.copy(work/'gate-summary.json',out);shutil.copy(work/'mutant-source-after.json',out)
for name in ['package-final.py','verify-archives.py','summarize-gates.py','finalize-docs.py','doc-finalization-first-failure.txt']:shutil.copy(work/name,out/'apparatus')
shutil.copy(work/'evidence-readme.md',out/'README.md')
shutil.copy('/tmp/alan-phase3b4-summarize-parent-gates.py',out/'parent-final-gates')
summary=json.loads((out/'verdict-summary.json').read_text());cost=json.loads((out/'cost-summary.json').read_text())
counts='| Mode | Cases / Stores | Raw captures valid | Both principals complete | Supported final C | Contacts / applications | Acknowledged losses |\n| --- | --- | --- | --- | --- | --- | --- |\n'
for mode,d in summary.items():counts+=f"| {mode} | {d['Cases']} / {d['RawEvaluated']} | {d['Captures']} / {d['Captures']} | {d['PrincipalComplete']} | {d['SupportedModelProjection']} | {d['Contacts']} / {d['Applications']} | {d['AcknowledgedLosses']} |\n"
costtable='Median milliseconds (five observations per cell); exact min/median/max and raw observations are retained in the cost JSON files.\n\n| Mode | n | Construction | First decision window | Fresh gate window | Replay | Open |\n| --- | --- | --- | --- | --- | --- | --- |\n'
for mode,groups in cost.items():
 for n,d in groups.items():costtable+='| '+mode+' | '+n+' | '+' | '.join(f"{d['Nanoseconds'][k]['median']/1e6:.3f}" for k in ['construction','first_decision_window','fresh_gate_window','Replay','Open'])+' |\n'
gatesummary='Full ordinary tests pass in both repositories; full runtime integration and integration-race pass; affected research model/checker race and prior reciprocal/authority integration/race pass. Both repositories pass build, module verification/tidy-diff, format and vet (ordinary/integration). Runtime Staticcheck passes in both modes; changed research checker Staticcheck is clean. Global research Staticcheck remains exit1 with exactly the inherited13 ordinary and31 integration diagnostics, byte-identical to the baseline. Both final build logs are empty, resolving the earlier module-cache warning. Every final command retains unchanged source/input maps.'
skips='The named optional test skips are `TestPhase3B4PromotedPublicParents` in each first full runtime integration invocation (satisfied by later explicit normal/race passes), and `TestPhase3B3BaselineS2CompletionExpectedFailure` in those two full suites (historical opt-in negative control, not enabled or claimed as a final result). The full ordinary research suite also skips three opt-in artifact exporters: `TestCommunicationArtifactExport` (explicit artifact export not requested), `TestJoinArtifactExport` (`ALAN_JOIN_EXPORT` unset), and `TestCandidateAExport` (`ALAN_PHASE1I_A_OUTPUT` unset). These exports are outside the required Task5 model-export set. No independent Binding, affected research race or prior regression test is skipped. Package-level `[no test files]` events are separately recorded in `gate-summary.json`; they are not omitted test executions.'
# This assertion prevents prose from silently omitting new named skips.
for g in gates:
 for s in g['test_skips']:assert s['test'] in ['TestPhase3B4PromotedPublicParents','TestPhase3B3BaselineS2CompletionExpectedFailure','TestCommunicationArtifactExport','TestJoinArtifactExport','TestCandidateAExport'],(g['name'],s)
table='| Executed final gate | Exit | Elapsed seconds |\n| --- | --- | --- |\n'
for g in gates:table+=f"| {g['name']} | {g['exit']} | {g['elapsed_seconds']:.3f} |\n"
q=(work/'qualification-draft.md').read_text().replace('FINAL_COUNTS',counts).replace('FINAL_GATES',gatesummary+'\n\n'+skips).replace('FINAL_COSTS',costtable)
(repo/'docs/foundational/phase3b4/06-qualification.md').write_text(q)
r=(work/'report-prefix.md').read_text();r=r.replace('Status: final campaign in progress. Executable sources are frozen; evidence-only packaging follows. This file is deliberately untracked.','Status: implementation and final exact-source qualification complete; controller independent review and Task6 manuscript reconciliation remain separate. This file is deliberately untracked.')
r=r.replace('The final normal corpus currently passes:', 'Both final normal and race corpora pass, each with:').replace('Race/broad gate completion and final archives are still pending.','All final broad gates have completed; exact inherited lint remains separately non-green. Complete raw archives and command records are retained.')
r+='\n'+counts+'\n'+(work/'report-tail.md').read_text().replace('FINAL_GATE_TABLE',gatesummary+'\n\n'+table).replace('FINAL_SKIPS',skips).replace('FINAL_COST_TABLE',costtable)
(repo/'.superpowers/sdd/03-plan/task-5-report.md').write_text(r);(out/'report.md').write_text(r.replace('This file is deliberately untracked.','This is the tracked copy of the untracked implementation handoff report.'))
rt=(work/'runtime-task5-draft.md').read_text().replace('FINAL_COUNTS',counts).replace('FINAL_GATES',gatesummary+'\n\n'+skips)
(runtime/'docs/qualification/phase3b4/task5.md').write_text(rt)
rtdata={'QualifiedPair':json.loads((out/'freeze/pair.json').read_text()),'Counts':summary,'CostObservationsPerMode':20,'Gates':gates,'ResearchEvidenceDirectory':'docs/foundational/phase3b4/evidence/task5','NoRuntimeExecutableChanges':True,'IndependentReview':'controller pending','ManuscriptReconciliation':'Task6 pending'}
(runtime/'docs/qualification/phase3b4/task5-summary.json').write_text(json.dumps(rtdata,indent=2)+'\n')
(out/'artifacts.json').write_text(json.dumps({str(p.relative_to(out)):hashlib.sha256(p.read_bytes()).hexdigest() for p in sorted(out.rglob('*')) if p.is_file() and p.name!='artifacts.json'},indent=2)+'\n')
print('Finalized',len(gates),'command records; documentation and report written; no executable source edits')
