# Immediate-selection narrowing: retained runtime refusal versus model allow

Independent read-only audit, 2026-11-04. **The proposed inference is confirmed.** The supplied ordinary-runtime control has a genuine admitted BA request, a genuine accepted grant covering the exact A4 service, no retained revocation, and both originals unpaid. Its first service authority decision refuses because the immediately preceding A reduction selected grant-ready instead of BA. The unchanged S2/LateGrant model allows delegated authorization after the corresponding admitted-request/retained-grant prefix. This is a concrete difference between the model and the deliberately narrower runtime interaction contract, not an authority-amplification counterexample or evidence of full correspondence.

No repository edits, subagents, commits, runtime executions or SQLite opens were performed. JSON evidence was read directly. This audit does not certify a general raw evidence checker or a runtime progress theorem.

## Exact retained witness

Input: `/tmp/alan-phase3b3-runtime-development/controls/control-wrong-prior-selection-1914204028/final/inspect.json`, SHA256 `60897acf3d7cc14ab351e23314ddb656e1b68801bbeaa772521623bb79dc0327`. Namespace `18b7dbe0e0946ecb6c270c7e51f8418779b530c5f863bd638e2cbacba40993bb`, storage 13.

| Role | Actual immutable ID |
| --- | --- |
| A→B original request | `["a",1,1,0]` |
| B→A original request | `["b",1,1,0]` |
| B's committed DelegateAuthority operation / accepted grant | `["b",1,2,0]` |
| B's grant-ready message | `["b",1,5,0]` |
| A's exact granted service operation | `["a",1,4,0]` |

The targeted scratch checker joins every emitted intent to its actual same-generation history row and canonical index identity; joins message admissions to their original intents and generation-bound observations; verifies the original permission row, exact claim/use/service tuple and root membership; reconstructs root hashes from namespace/logical Run; and checks the chronological facts and selection mismatch. It imports no runtime helper and reads no SQLite. It is a targeted assertion script for these supplied bytes, not a hostile-input-complete independent validator.

Observed chronology from `composition_facts`:

| Global order | Retained event |
| ---: | --- |
| 1 | A sequence 1 commits AB original. |
| 3 | B sequence 1 commits BA original. |
| 5 | A sequence 2 selects BA original's admitted observation. |
| 8 | DelegateAuthority result is accepted, reason `delegated`; decision sequence 1. |
| 13 | A sequence 3 selects grant-ready's observation. |
| 14 | A sequence 4 commits the exact service named in the grant. |
| 15 | Existing source authorization for that service allows. |
| 16 | Immutable external permission attempt/reservation refuses `authority_request_not_admitted`. |

BA's observation ID is `["observation","message","[\"b\",1,1,0]"]`; its `consumed_sequence` is 2 and A history sequence 2 selects precisely that ID. The ready observation is `["observation","message","[\"b\",1,5,0]"]`, consumed at A sequence 3. A's service is ordinal 0 of sequence 4, so the immediate predecessor is ready, not BA. A's sequence-4 state still retains `peer_request` with the real BA ID; the failure is not absence of that earlier request value.

The grant's Root.Run and Root.RootRun are both `b/g1`; Recipient and RecipientRoot are `a/g1`; recipient revision is 1. The complete original B permission is `(external, b, http://127.0.0.1:54097/apply, apply, revision 1, ba-service)`. The service selector uses this exact root and grant ID. OriginalRequestID is the genuine BA commitment and ServiceOperationID is exactly A4. The lifecycle request equals its committed intent, and its accepted result leaves Before=After with A open at sequence 2. Grant and serving root IDs match the namespace-derived real roots and their membership rows.

There is exactly one lifecycle result, for DelegateAuthority; no RevokeAuthority intent/result appears. Both principals retain generation/revision 1 and are unretired; no accepted generation/retirement mutation exists in this snapshot. A and B are still open immediately before the service decision. Thus neither revocation, wrong service locator, changed generation, absent original permission, nor a closed-source first authorization explains this refusal.

Both genuine originals were committed before the first service attempt (orders 1 and 3 precede 14/16). AB's message observation is present but unselected. BA's request is selected, but request selection does not pay BA's result obligation. There are no BA/AB reply intents, no external receipts, no physical-admission facts and no effect charges. The retained server receipt observation is a denial (`accepted=false`, `grant_denied`), not an authoritative successful service result. Accordingly **both original service obligations are unpaid**; no ready/lifecycle/denial selection substitutes for their authentic result consumption.

The separately read `endpoints.json` reports zero contacts/applications at both endpoints. `run.json` reports Error=`<nil>`, Quiescent=true, A failed and B still open. These are prior control-run artifacts, not commands rerun in this audit. Runtime quiescence here plainly does not establish the model's completion condition.

## Actual model witness and the disagreement

`/tmp/alan-phase3b3-selection-model.go` starts from unchanged `Initial(S2/LateGrant)` and selects unique actual `Successors` events, with Duplicate=false. It generates:

```text
commit AB → send-request AB → admit AB →
commit BA → send-request BA → admit BA → retain-grant BA →
authorize BA (Allowed, Delegated)
```

This is a generated model witness, **not a transcription of unrecorded runtime transport timing**. The snapshot fixes the relevant commitments, BA selection before grant, grant before A4 and lack of service payment. It does not itself put every transport send/admission into the composition fact total order. The chosen AB send/admit placement is an admissible witness linearization of the available causal evidence; its exact wall-clock placement is not asserted. No missing runtime authority/receipt was supplied from this model witness.

At the model state immediately before authorization, both originals are committed and unconsumed, BA is admitted, the grant is retained, Revoked=false and Unsafe=0. The sole enabled BA authorize successor is Allowed with Path=Delegated, root `B/g1:b`, grant `grant:BA:b`, original BA and generation-one parties. There is no corresponding Refused BA authorize successor. The exact generated descriptors, before/after states and candidate list are retained in `/tmp/alan-phase3b3-selection-model.json`.

Unchanged `AuditTrace` independently returned:

```text
ProtocolValid=true; AuthoritySafe=true; Conditions=[];
Paid=[false,false]; Unpaid=[true,true]; EventsChecked=8.
```

The proposed coarse map makes earlier request selection, grant lifecycle response selection and ready selection stutter relative to the two model originals. The model retains admission and exact service/grant identity; it has no “latest selected observation” coordinate. Therefore the ready selection cannot invalidate this model grant or disable the model authorize transition. The actual runtime refusal at order 16 is the **first observable mismatch**, not an invalid model history caused by missing consent.

`request_authority.go:445–486` explains the runtime boundary: it loads `history[service.Sequence-1].selected_observation`, validates that selected observation, then requires it to be the admitted observation of the grant's OriginalRequestID. The selected ready message passes the same-source/generation checks but is not BA's admitted observation, producing `authority_request_not_admitted`. The raw reason must be retained exactly; its interpretation here is “immediate predecessor is not this original request's admitted observation,” not “the request was never admitted.”

The runtime is enforcing frozen contract §3's explicit immediate-next-reduction restriction. There is no new production semantic change in this control: `phase3b3_controls_test.go:182–184` swaps A Code[1] and Code[2] in the ordinary public input program. Positive `examples/delegation/reciprocal.go:150–155` instead orders grant-ready selection, BA original selection, immediate BA service. The positive ordering preserves the current predecessor check. This audit read that positive source; it did not rerun a positive runtime execution.

## Consequence for Task 3 and claims

Preserve a named discrepancy such as `immediate_selection_narrows_model_admission` with separate fields:

- Model proposed BA authorization: generated, protocol valid, authority safe, Allowed/Delegated.
- Actual runtime authority decision: Refused, raw reason `authority_request_not_admitted`.
- Runtime contract verdict: legitimate refusal under its declared immediate-next-reduction rule.
- Correspondence verdict under unrestricted ready-selection stuttering: mismatch.
- Restricted positive-domain membership: false, because the latest predecessor selection is ready rather than the original.
- Actual custody: both originals unpaid; zero physical exposures in the supplied endpoint evidence.

Do not relabel this model history “invalid missing authority” or claim model/runtime agreement merely because the negative runtime test expected refusal. Conversely, this does not require weakening the runtime check: the minimum implementation deliberately does not infer arbitrary retained dataflow. Report the actual narrower domain and preserve the counterexample alongside the positive handshake. If a future equivalence treats the extra ready selection as harmless stuttering for every runtime program, this witness falsifies that claim; if a future theorem restricts to exact immediate-request-selection programs, this case is explicitly outside its hypothesis.

The model also abstracts tuple/result vocabularies. This witness checks the BA authority chronology and original identities; it does not establish a complete physical-resource isomorphism or arbitrary two-agent runtime reachability. No completion suffix was executed in this runtime control, and no whole-runtime progress theorem is inferred from the finite model's completion reachability.

## Reproduction, source identity and digests

Executed from research checkout, Go `go1.27.3 darwin/arm64`; both commands exited 0:

```sh
python3 /tmp/alan-phase3b3-selection-raw.py > /tmp/alan-phase3b3-selection-raw.json
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3b3-progress-go-build go run /tmp/alan-phase3b3-selection-model.go > /tmp/alan-phase3b3-selection-model.json
```

`/tmp/alan-phase3b3-selection-raw.json` records exact joined IDs, composition orders, observation consumption sequences, claim, runtime verdict, original program digest, no-revocation check, custody status, and raw-cell SHA256 references. `/tmp/alan-phase3b3-selection-sha256.txt` records complete SHA256 values for both scripts/results, input snapshot, endpoint/run artifacts and inspected runtime sources.

Research source used is unchanged `360cae66f5ffdd6e7a95a568552df5522122da21` authorityservice model/oracle; hashes are also retained in the earlier progress audit. Runtime checkout HEAD observed `347094a7b18d351c6ad18df526ef6133777bca44`, with the public example and qualification source untracked at inspection. Therefore HEAD alone does not identify the control source. Exact source hashes observed:

```text
04c151eb375ed9d2fb5bed9598ad6602bb06b84dc66d7d3f6cad1c7f19e641ac  internal/durabletest/phase3b3_controls_test.go
c23032d29b2d3c39972bc158781ec183fc0940b4401b1c29f13c4b1f8d3c2387  examples/delegation/reciprocal.go
0428f5fe32453d08e2d6388f723dda76f8edf42cd8c276489938b05bbb6aec38  internal/executor/request_authority.go
335cb8953c4bad5462caf32660ca4d4481aca5b0efad624fe5341a19c5bdb694  /tmp/alan-phase3b3-selection-raw.py
15ddcdb7f3124497f9d3efdf275b7650cd810e1dc2775edc2f4701db254bc7ba  /tmp/alan-phase3b3-selection-raw.json
85d2435b10a46f1a7e7ff2dc9cdb4dfce2755591b8e15873f5f0d835779a49ad  /tmp/alan-phase3b3-selection-model.go
833586b3a3568c1684b7726bcb4b54b4973be025955f65fb9d3d2e640108d8aa  /tmp/alan-phase3b3-selection-model.json
```

No Replay/Open/Inspect method, old audit suite, SQLite reader, endpoint or runtime was invoked in this audit. Existing exported JSON is evidence of the earlier run, with its provenance bounded by its supplied bytes and source context.
