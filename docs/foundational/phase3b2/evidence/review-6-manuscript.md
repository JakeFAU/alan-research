# Scope 6 independent review: manuscript claims and reproducibility

**Disposition: APPROVE for the current Task 4 manuscript, results, runbook and retained-evidence scope. Task 4 specification and quality: APPROVE.** No blocking semantic or artifact finding remains in this scope. This does not approve an unfinished ordinary repository suite, the final delivery manifest, or the whole branch. Those are explicit remaining gates below; no pending test is counted as a pass.

## Scope and source identity

Read the frozen brief and contract, results and runbook, supplied Task 4 patch, actual edited manuscript section/Abstract/Conclusion/BibTeX, runtime-archive report, command/source/model/runtime manifests, compressed model and independent ledger exports, and all five previous scoped reviews. The R4 REJECT, exact retained generation counterexample, and subsequent R4.1 ADDRESSED/APPROVE remain distinct historical facts. I did not waive that finding based on another reviewer's agreement.

The executable source is `f75e26984c84fe649f26d5e4eaf4ed70e7506e46`, based on `1d41d1e20b579443c98a6caa53e51036e55f2d82`. My independent artifact probe compared **all 563** frozen manifest records to their Git objects at that executable commit and the current source bytes. Every size, SHA-256 and byte comparison passed. This does not imply that historical snapshots or separate consumer modules were executed.

The post-run Python utility has SHA-256 `1bd2650abee631047945aeeae0d67c73a47776009e313110fef0868e739e8465`; its reproduced output has SHA-256 `f92c8e7ea60f04934df42ec6370609eb140ece7c83cc1c843ae8c363f4b3015f`. It is correctly identified separately from the frozen Go execution source. The runtime index digest is `6200a8f7ad0c1547213f8279f4369e23d2f93f4d64060b1ce4b5e7a90f1177ce`.

Reviewed results/runbook digests were respectively `0d2dd2556075cad5a5880d106badaf5b3a3cbc3a165d5a482a52b752f23311cd` and `90c02250644a6ac7bcf42c6e14a022e799c4862b0ce67c9cb2443e73059aaae6`. These documents still correctly report the root suite and final integration review as pending. Later delivery-status edits must preserve the semantic boundaries reviewed here.

## Claims checked

- Permissions are complete tuples. S1 uses server roots; S2 additionally uses the exact one-hop rooted, original/generation-bound grant; S3 adds ambient requester-root borrowing. The explicit-consent safety policy is visible. Ambient borrowing can violate consent without inventing a globally unheld tuple. The separately named Cartesian mutant invents x. The paper neither conflates these failures nor condemns every possible implicit authority policy.
- The safety argument is an induction over the declared finite transition relation under authentic identity/channel and intact-storage assumptions. It is labelled a conditional source-level argument. Immutable first decisions, validity of the same retained path at fresh admission, activation-bound opportunities, revocation after admission and server-crash loss of opportunities are distinguished. The custody projection removes applications and maps new authority bookkeeping to stutters. Shared b results remain opaque, so no counter arithmetic or shared-resource isolation is inferred.
- The positive completion argument assumes covering authority before the first decision and at needed admissions. It constructs one BA-then-AB completion suffix, with recovery and no further interruptions chosen. It is not scheduler inevitability. LateGrant is completable exactly when BA's decision is not Refused; Revocable is completable exactly when not revoked or **both** originals are Applied or retain their own Permit. The complement is stated to be closed/nonterminal within this relation. This matches the prior independent constructive/closed-set review, not a universal recovery theorem.
- Four/five/six-event forbidden authorization/admission/application prefixes are distinguished from zero/four/one-event shortest missing-authority/late-refusal/revocation completion obstructions. A shortest authority failure needs only BA, not both reciprocal waits. No claim relies on mislabelling those different minima.
- All sixteen results-table rows match the compressed operational export exactly, including safe/unsafe prefix, locally violating arcs, raw completable, safely completable and raw blocked counts. Totals are 1,198,963 states and 9,155,940 arcs. Six positive configurations have 103,361 states each; six insufficient-authority configurations have 6,741 each. The two unsafe cases have 98,546 unsafe-prefix states and 5,776 locally violating arcs each. LateGrant and Revocable retain 3,852 and 26,883 blocked states. No finite history count is inferred from cyclic graphs.
- Ledger totals independently summed from the retained export are 1,198,963 BFS prefixes, 9,155,940 outgoing arcs, 197,092 unsafe prefixes and 1,536,272 unsafe arc extensions. These are not the 11,552 locally violating arcs. The export has 33 history inputs and four ancestry inputs; the eighteen malformed complete histories and four equal/different-value pairings are retained. Rootless/returned chains are hostile inputs, not a qualified arbitrary multihop operational protocol.
- Five runtime scenarios and nine checked-SIGKILL cuts per normal/race run are supported by retained records. Their public runtime authority is static provisioning, not general S2 delegation. Grant-retained cuts are explicitly inapplicable. Reservation-without-result and physical-admission-before-call preserve exact unknown identities, zero contacts and unpaid obligations; the atomic model can re-admit where the generic runtime conservatively cannot. No general model/runtime refinement, independent-owner recovery, N-agent closure or unrestricted effect guarantee is claimed.
- The Abstract and Conclusion preserve these limitations. The new Markdown and TeX authority sections agree on the propositions, counts and predicates on direct comparison; the parent's format-normalization record includes Abstract and Conclusion. The PDF digest matches the parent's inspection record, which reports 37 pages and visual inspection of pages 1–3 and 34–37. I did not independently rerender or visually inspect the PDF, and do not claim to have done so. Existing warnings remain disclosed.
- I independently opened the [primary RFC 2693 source](https://datatracker.ietf.org/doc/html/rfc2693#section-6.3). Its September 1999 Experimental status, authors and title agree with BibTeX; §§6.1–6.3 support the limited attribution of issuer/subject linkage, delegation permission and authorization/validity intersection. No novel-attenuation or SPKI-implementation claim is made.

## New artifact checks actually executed

Executed from repository root, exit 0:

```sh
python3 docs/foundational/phase3b2/evidence/runtime-audit.py docs/foundational/phase3b2/evidence /tmp/phase3b2-review6-runtime-audit.json
cmp docs/foundational/phase3b2/evidence/runtime-audit.json /tmp/phase3b2-review6-runtime-audit.json
python3 /tmp/phase3b2-review6-probe.py
```

The first two commands regenerated the full bounded archive audit byte-for-byte. This checked both archives' member/digest/source/kill joins, runtime cases/cuts, 115 local SQLite backup reconciliations per run, 46 retained read-only Replay records per run, endpoint provenance links and retained negative controls. It did not rerun Replay, process kills or Go tests.

The third command is the independent probe retained below. It rebuilt each 521-member deterministic tar.gz from its contained bytes and newly constructed normalized TarInfo records, reproducing both archive files exactly. It also checked all command artifact hashes and their decompressed output hashes, all frozen Git objects, all table rows and model-summary fields, ledger totals and the inspected PDF digest. Both fresh negative controls changed only the separate first kill record's PID; the archive checker rejected each with `actual DURABLE kill records mismatch`. These distinguish checked logged signal evidence from a self-asserted crash JSON. No archive or database bytes were changed.

## Findings and final gate conditions

**R6.1 — baseline command source attribution: ADDRESSED.** During review, the baseline Staticcheck command initially named the task's frozen source even though its cwd was the merged-base materialization. I reported that inconsistency. The parent corrected `source_commit` to the actual base and retained the task Go identity separately, with a checked materialization record. My probe verified the corrected base identity and all current command artifacts/digests. The six integration Staticcheck diagnostics remain explicitly non-green and unchanged on the base; this review does not suppress them.

The runbook's `evidence/source.json` is an intended final delivery manifest, not yet present at this scoped review. The parent explicitly confirmed it will be produced after the root suite and reviews. Approval here is not a claim that this remaining delivery artifact exists. Before PR completion: record the ordinary root suite's actual exit and retained output; create and validate the final manifest and all paths it references; update pending delivery statements consistently; retain this scoped review alongside the previous reviews; and obtain whole-change review. If a later result falsifies a claim or changes executable bytes, revise/requalify the affected claim instead of extending this approval silently.

No production or test source was edited, no full suite was duplicated, and no repository semantic fix was made in this review. Agreement among reviewers and archive checks is evidence within the stated assumptions, not an independent universal theorem.

## Independent probe source and output

Save the following as `/tmp/phase3b2-review6-probe.py`; it reads repository artifacts, builds archives in memory, and uses scratch read-only SQLite copies only through the already documented utility. Importing that utility may produce a disposable `__pycache__`; the review removed its generated cache before finishing.

```python
import copy,gzip,hashlib,importlib.util,io,json,pathlib,re,subprocess,tarfile
p=pathlib.Path('docs/foundational/phase3b2/evidence');sha=lambda b:hashlib.sha256(b).hexdigest()
f=json.loads((p/'frozen-source.json').read_text());commit=f['ExecutableSourceCommit']
for x in f['Files']:
 b=subprocess.check_output(['git','show',commit+':'+x['path']]);assert len(b)==x['bytes'] and sha(b)==x['sha256'];assert pathlib.Path(x['path']).read_bytes()==b
print('PASS 563 manifest source records equal frozen Git objects and current bytes')
m=json.loads(gzip.decompress((p/'model.json.gz').read_bytes()));s=json.loads((p/'model-summary.json').read_text());a=json.loads(gzip.decompress((p/'audit.json.gz').read_bytes()))
assert len(m)==16 and sum(x['States'] for x in m)==1198963 and sum(x['Arcs'] for x in m)==9155940
for x,y in zip(m,s['Configurations']):
 for k,v in y.items():
  if k=='WitnessLengths':assert all(vv==(len(x[kk]['Trace']) if x[kk] is not None else None) for kk,vv in v.items())
  elif k=='CutLengths':assert all(vv==len(x['Cuts'][kk]['Trace']) for kk,vv in v.items())
  else:assert x[k]==v
rows=[z for z in pathlib.Path('docs/foundational/phase3b2/03-results.md').read_text().splitlines() if re.match(r'\| (S[123]/|Cartesian/)',z)]
for row,x in zip(rows,m):
 parts=[v.strip() for v in row.split('|')[1:-1]];assert parts[0]==x['Config']['Policy']+'/'+x['Config']['Scenario'];assert list(map(lambda z:int(z.replace(',','')),parts[1:]))==[x[k] for k in ['States','Arcs','SafeStates','UnsafeStates','ViolatingTransitions','CanComplete','CanCompleteSafely','CannotComplete']]
assert len(rows)==16
assert [sum(x[k] for x in a['Graphs']) for k in ['Prefixes','Edges','UnsafePrefixes','UnsafeEdges']]==[1198963,9155940,197092,1536272]
assert len(a['Controls']['Histories'])==33 and len(a['Controls']['Chains'])==4
print('PASS model summary, all 16 documentation rows, 33 histories/4 chains and ledger totals')
c=json.loads((p/'commands.json').read_text())
for x in c['Commands']:
 b=(p/x['artifact']).read_bytes();assert sha(b)==x['artifact_sha256'];raw=gzip.decompress(b) if x['artifact'].endswith('.gz') else b;assert sha(raw)==x['raw_log_sha256']
print('PASS all command artifacts/digests, including decompressed output identities')
assert next(x for x in c['Commands'] if x['name']=='staticcheck-runtime-baseline')['source_commit']==f['BaseCommit']
assert sha(pathlib.Path('paper/alan-dissertation.pdf').read_bytes())==json.loads((p/'pdf-inspection.json').read_text())['sha256']
print('PASS corrected baseline identity and inspected PDF digest')
spec=importlib.util.spec_from_file_location('auditor',p/'runtime-audit.py');mod=importlib.util.module_from_spec(spec);spec.loader.exec_module(mod)
for run in ['normal','race']:
 original=(p/('runtime-'+run+'.tar.gz')).read_bytes();out=io.BytesIO();raw={}
 with tarfile.open(fileobj=io.BytesIO(original),mode='r:gz') as tf:
  members=tf.getmembers()
  with gzip.GzipFile(filename='',fileobj=out,mode='wb',mtime=0) as gz:
   with tarfile.open(fileobj=gz,mode='w',format=tarfile.PAX_FORMAT) as rebuilt:
    for x in members:
     data=tf.extractfile(x).read();raw[x.name]=data
     info=tarfile.TarInfo(x.name);info.size=len(data);info.mode=0o644;rebuilt.addfile(info,io.BytesIO(data))
 assert out.getvalue()==original
 kills=json.loads((p/('kill-records-'+run+'.json')).read_text());kills[0]['pid']+=1000000
 try:mod.audit_run(run,raw,kills,{x['path']:x for x in f['Files']})
 except ValueError as e:print('PASS '+run+' deterministic archive bytes and forged kill PID rejected: '+str(e))
 else:raise AssertionError('forged kill accepted')
```

Actual output:

```text
PASS 563 manifest source records equal frozen Git objects and current bytes
PASS model summary, all 16 documentation rows, 33 histories/4 chains and ledger totals
PASS all command artifacts/digests, including decompressed output identities
PASS corrected baseline identity and inspected PDF digest
PASS normal deterministic archive bytes and forged kill PID rejected: actual DURABLE kill records mismatch
PASS race deterministic archive bytes and forged kill PID rejected: actual DURABLE kill records mismatch
```
