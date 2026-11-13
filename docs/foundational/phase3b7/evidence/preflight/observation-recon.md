# Phase 3B.7 observation reconnaissance

Read-only reconnaissance, 2026-11-13. No runtime experiment, database open, source/index/branch mutation, or subordinate agent. This file is the only output artifact written by this task.

## Question and result

Can a genuinely captured admitted-but-uncalled v2 execution, killed once and reopened before exact b revocation and fresh refusal, have the same declared local observation as the sealed applied-but-unrecorded execution? The finite model says its deliberately small projection admits that pair. The existing sealed runtime corpus does not yet establish it. The proposed experiment should distinguish **model-projection equality**, **full exported local-store equality modulo declared identity renaming**, **protocol validity**, and **opposite world facts** as four separate verdicts.

The main risk is silently identifying the model projection with the complete Store. A model match may coexist with differences in retained program counters, projections, sequences, history selections, message acknowledgments, resource descriptors, composition-fact order, pending deliveries, or run status. Such a difference falsifies the proposed richer local match; it does not refute the finite theorem about a smaller projection. A failed richer match must remain visible rather than being repaired by deleting the differing coordinate.

## Materials inspected and direct checks

Read AGENTS.md; Phase 3B.6 contract, model, observations, correspondence and results; `research/effectobservation/model.go` and `observations.go`; `research/delegationcorrespondence/effect_observation.go`, `effect_mapping.go`, `effect_protocol.go`, relevant portions of `check.go` and `wire.go`; retained Phase 3B.6 producer-source extracts. The using-superpowers skill explicitly instructs dispatched subagents to ignore its workflow.

Read sealed JSON/JSONL from these existing extractions:

- `/tmp/alan-phase3b5/corpus/raw-normal/applied-unknown-revoked-1243056310`
- `/tmp/alan-phase3b5/corpus/raw-race/applied-unknown-revoked-2740959369`

Freshly SHA-256 checked both owning archive byte streams against `docs/foundational/phase3b4/evidence/prospective-main/archives.json`, and all 64 declared regular members of each selected loss dossier against that manifest. Both computed archive comparisons and all 128 member comparisons passed. SQLite, WAL and SHM members were read only as opaque bytes for hashing. No SQL connection, PRAGMA, integrity check, replay against a database, or runtime test ran.

The hash check used Python standard-library `pathlib`, `hashlib`, and `json`; it compared archive SHA-256, each selected member length, and each member SHA-256. Further read-only JSON inspection enumerated row counts, per-field differences, run coordinates, ordered composition facts, and selected history/acknowledgment rows. No fresh Go tests or checker exports ran. Source-extract content was read as retained evidence; this task did not reverify it against frozen Git blobs, so Phase 3B.6's prior verification is attributed as prior evidence.

## Established historical local cut

Both historical modes use request-authority v2 and the S2 terminal-handshake fixture. Selected BA principal request is `["b",1,1,0]`; its distinct serving effect is `["a",1,4,0]`. Grant is `["b",1,2,0]`, ready message `["b",1,5,0]`, and controlled revoke `["b",1,7,0]`. Do not collapse principal original, effect, authorization, wire key, or composition-root identities.

In the inspected normal dossier:

| Cut | A epoch / sequence / PC | B epoch / sequence / PC | Composition facts |
| --- | --- | --- | --- |
| stopped | 2 / 4 / 4 | 2 / 5 / 5 | 17 |
| reopened | 3 / 4 / 4 | 3 / 5 / 5 | 17 |
| final | 3 / 4 / 4 | 3 / 8 / 8 | 23 |

Final inspection contains 41 tables and 94 rows. Final has one external attempt, zero external receipts, one reservation, one effect charge, two lifecycle results, two settlements (grant/revoke only), two open runs, twelve history rows, six observations/admissions/attempts/intents, three message acknowledgments, and no composition result/reply. Fact 17 is allowed gate ordinal 1, fact 21 the accepted exact revoke, fact 23 refused gate ordinal 2. The original wire and historical source remain epoch 2; fresh refusal is epoch 3. `run.json` reports Steps 0, Dispatches 1, Quiescent true, one unresolved/unknown BA external operation, one pending external delivery, and no unknown adapters. Neither principal is paid.

The complete b root includes owner/run, root-run, permission kind/source/target/operation/revision/obligation. Its target is endpoint B and obligation `ba-service`; distinct c targets endpoint A with obligation `ab-service`. They share owner B but cannot be aliased. The inspected normal program IDs depend on concrete endpoint URLs. Composition root IDs depend on namespace and run. Both participate recursively in retained descriptors and hashes.

Historical success here means one endpoint B contact/application with no local receipt. Endpoint A is empty. Both process acknowledgments precede confirmed SIGKILL by the pinned controller protocol. Endpoint B's post-commit acknowledgment binds the exact key/authorization/source and ready identity. Reopen advances both activations once. Controller revocation precedes the fresh refusal. The two acknowledgment emission times have no certified total order. The final report's one dispatch consideration is not a physical retry.

## Three comparison domains to declare before capture

### 1. Existing model observation

The exact existing StoreView comprises symbolic original/source generation/source revision/grant/root/original activation/wire binding; accepted grant; historical first allow; one debit; first gate allowed; owner stage/current activation; b revoked; distinct valid c; fresh refusal; retained result; consumption. It excludes launch, ingress, application and unretained response. Owner stage comes from death/Open protocol, not inspection alone. Actual epoch 2 maps to original activation 1; actual epoch 3 to reopened activation 2. Keep this map fixed for the new paired protocol.

The model's symbolic root tuple is an abstraction; it is not the runtime's literal permission descriptor. An equal StoreView is a projected match and must not be described as byte-identical Store or unrestricted local indistinguishability.

### 2. Full exported local state modulo identity renaming

Recommend comparing parsed `inspect.json` at stopped, reopened and final, including all 41 table names, every row/column, nested descriptors, initial/current projections, ordered histories/facts, complete permissions, and derived Roots summaries. Include Settings separately, requiring equality for the controlled experiment. Settings/Roots were explicitly outside Phase 3B.6's renewal helper; adopting this stronger new comparison must not retroactively describe old checks as defective.

Also compare final `run.json` after the same renaming. It is a runtime observation available to an observer and may distinguish two executions even when the retained table subset matches. Require row multiplicity, order coordinates and zero-valued fields to match; never replace a full comparison by counts. Sort unordered SQL-export row collections by schema-defined normalized keys when necessary, but retain sequence/decision-order fields and ordered nested arrays. Normal/race `composition_static_members` order differs because sorting raw root IDs changes row order; joining by normalized run avoids falsely claiming a semantic subject change.

Opaque SQLite byte equality is neither required nor sufficient: page/WAL layout and random namespace create differences. State the result as equality of **the complete declared JSON export**, not equality of every hidden SQLite bit. Hash binary artifacts for custody only.

### 3. Protocol and forensic evidence

Case/source/protocol declarations, PID/token/path identities, kill/Open records, pre-call versus in-flight/post-commit acknowledgments, boundary journals, and endpoint dossiers establish provenance, causal cuts and world facts. They are expected to differ and are not the equal-local observation. In particular the boundary journals themselves distinguish pre-call kill from post-commit kill. If an observer is permitted this entire protocol dossier, the histories are distinguishable by construction. Declare this exclusion expressly; merely calling both collections "local evidence" would be false.

The endpoint dossier must positively certify zero versus one at each declared cut under complete process-bound coverage. Missing JSONL is unsupported, not zero. A present empty journal requires ready identity and complete interval/snapshot evidence. Zero at the final cut is a fact at that cut; a general future-nonapplication theorem requires additional reasoning about outstanding packets.

## Permitted identity renaming

Fresh runtime.Create namespaces, ephemeral endpoint addresses, process IDs/tokens and temporary paths differ across genuine runs. Preserve raw bytes. The comparison creates a separate derived view with a typed, injective role map and a complete mapping manifest.

1. Map the one Store namespace to a shared symbol after verifying all metadata and identity-bound uses agree within each capture chain. Do not map arbitrary unrelated strings.
2. Map endpoint A and endpoint B URLs to two distinct symbols, only after ready/process/provision/permission/wire joins establish their roles. Preserve scheme, method/path semantics, all other permission fields, and endpoint distinction. Ports are deployment coordinates, not permission authority to be erased without those joins.
3. Validate each program ID against its raw retained program bytes first. Compare actual program structure after only permitted endpoint-address substitutions. Then derive a symbolic program identity from normalized program structure. Never pair arbitrary hashes merely because two runs each have two programs.
4. Validate composition-root IDs against namespace and exact run (`digest("alan/resource-profile/v1", [namespace, run])` in the inspected checker), then derive normalized root IDs. Root IDs may differ by namespace; root membership, limits and authority are not renameable.
5. Validate full original wire keys, attempt digests, fact IDs and record digests against their owning raw structures first; compare their normalized dependency structures or recomputed canonical identities. These hashes change transitively with namespace, roots, URLs and programs. Never drop all `Digest`, `ID`, or 64-hex fields. That would hide contradictory provenance and content.
6. For protocol validation only, join process identities/paths within each dossier before renaming process roles. Preserve owner-versus-endpoint separation, original-owner-versus-reopened-owner distinction, and one-loss ordering. These fields do not belong to the Store equality projection unless explicitly added.

No renaming of subject a/b, logical generation, sequence, ownership revision, authority-admission ordinal, request/service role, decision order, result status/value, debit, operation/obligation, retained data, PC, or program behavior is justified. Across this prescribed protocol the actual epoch sequence should also remain 1→2→3; use the existing relative model epoch map only for StoreView. Do not allow arbitrary epoch shifts to conceal an extra Open. Keep raw integers and null-versus-empty distinctions unless a specific schema equivalence was declared in advance.

A direct normal/race inspection found differences in 21 of 41 final tables before renaming, despite both being the applied branch. Those include programs, histories, runs, permissions, attempts, roots, facts, reservations, grants/revoke and metadata. This demonstrates that a naive bytewise cross-run comparison will fail. It does not establish that every such difference disappears under the proposed mapping; a full canonical comparator remains to be implemented and qualified.

## Minimal independent comparator

Use a new opt-in package/command or standalone research script with a versioned contract. Keep `CheckEffectObservationDirectory`, anchors.json, old commands and old reports unchanged. The old function pins exact manifests, source hashes, producer/test identities and process anchors and assumes one endpoint contact for the loss profile. Broadening it would erase the historical trust boundary.

Recommended pipeline:

1. Read only sealed declared artifacts. Validate archive/member hashes against an external capture manifest, strict JSON (duplicate keys rejected, integral numbers preserved), completeness, schema and exact input-source provenance. New harness source is its own explicitly pinned producer, not historical commit 5b4662f by label inheritance.
2. Independently validate each side's selected original, full grant/root/wire/attempt/reservation/debit, one allowed historical gate, loss acknowledgment, both-principal renewal, exact b revoke, valid c, one epoch-current refused gate, and absence of authentic receipt/effect settlement/consumption. Reuse unchanged lower-level raw validation as a separately named verdict if desired, but equality must not merely compare old model reports.
3. Validate raw referential hashes, construct the limited typed alpha map, produce full normalized local exports and machine-readable JSON-path differences. No write-back to inputs. Report exact input/output hashes and mapping-manifest hash.
4. Independently derive the small StoreView from validated facts, rather than copy fixed values based on a case name. Derive zero/one from complete endpoint evidence and source-conditioned cut certificates separately. No final endpoint record may be projected into an earlier cut without a causal certificate.
5. Report independent booleans for raw/provenance validity, protocol validity, world-zero/world-one, model projection match, stopped/reopened/final richer matches, Settings match, run-report match, and full contract success. Preserve mismatches even if a smaller projection passes. The expected endpoint facts must be unequal; an accidentally matched pair of two applied runs is a failed zero/one experiment.

Capture must be genuine and raw evidence immutable after sealing. Controller naming/source metadata may be prepared by the capture harness itself; a later rewriting pass to make new bytes resemble old anchors is prohibited. A derived comparison view is legitimate only as a separate artifact with raw hashes and the predeclared transformation rules.

## Adverse controls and concrete falsifiers

Lawful controls: same file against itself; independent normal/race historical applied cases (semantic equality is a control to test, not assumed); formatting/key-order re-encoding in explicit synthetic fixtures; a genuinely retained-result case to show custody detection; a role-consistent namespace/endpoint change with every derived digest correctly recomputed.

Required attacks:

- Equal small StoreView but changed PC, current projection data, selected observation, message acknowledgment, fact order, pending delivery or run report. Richer comparison must reject and name the exact path.
- Change a program instruction while mapping its new hash to the old role; reject normalized program inequality.
- Wrong namespace embedded in a nested wire key; stale raw digest; a correct digest of different content; duplicate identity mapping; endpoint A/B swap or b/c alias; reject even if row counts match.
- Extra/missing table, column, row or duplicate row; null/empty confusion; duplicate JSON keys; unexpected integer precision loss. Fail closed.
- Kill before the allowed gate or debit; wrong pre-call authorization/ordinal/source/PID/token; late/uncertified kill; only one principal renewed; extra reopen; changed historical allow; fresh gate still epoch 2; wrong exact root; extra contact; receipt or settlement fabricated from forensic evidence.
- Missing journal/ready/snapshot presented as zero; zero snapshot predating the cut; mismatched process interval; a second contact outside declared interval; unsupported later evidence back-projected into stopped.
- Both sides supplied from the same applied capture; reject lack of genuinely opposite world facts even though local equality succeeds.

Synthetic attacks must be labeled as checker fixtures, written only to separate scratch copies, and never promoted into genuine runtime archives. No mutation should make a malformed fixture accepted solely by updating its self-declared digest manifest. Compiled comparator mutants should demonstrate actual semantic false acceptance, rather than compiler failure.

**Concrete falsifier:** the fresh zero execution has valid selected BA admission, loss/reopen, b revoke, fresh refusal and empty local receipt, so all model StoreView columns agree; nevertheless final `runs[b].projection`, a `history.selected_observation`, a composition fact's `Order`, or `run.Report.PendingDeliveries` differs after the declared identity map. Then the full local match fails and the raw discrepancy must be retained. A single endpoint B committed row at the purported zero cut independently falsifies the intended world contrast. These are predeclared failure conditions, not reasons to relax the observation after seeing results.

## Recommendation

Predeclare the small model observation and stronger full-export comparison together; run the new producer with the historical S2 handshake and controlled-revoke program unchanged except necessary capture harness plumbing; freeze its source and comparison rules before capture. The handshake supplies B's ready-at-sequence-5 dependency before A's service, reducing but not proving absence of schedule-related local differences. Maintain independent physical/protocol evidence. Claim a matched concrete pair only if the actual archived runs satisfy all declared equalities. Otherwise report the precise counterexample and retain the finite model result at its original scope.

## Follow-up: preferred fresh pair with a common seed and common pre-call barrier

The controller reports runtime reconnaissance identifying `dispatch.admitted/pre-call` immediately before `client.Do`. Conditional on verifying that frozen source and the barrier's synchronous blocking semantics, a fresh/fresh pair is preferable to comparing a new zero run to an old applied run:

- Both branches reach the same exact allowed-gate/debit-bound pre-call event and capture or certify the same local prefix.
- Zero keeps that worker blocked through acknowledged SIGKILL. Source must establish that the blocked hook cannot return or call `client.Do` independently before death; the zero certificate must preserve this causal connection. Complete endpoint coverage corroborates no contact/application.
- Applied releases that same barrier, awaits identity-bound endpoint post-commit acknowledgment, withholds successful receipt through acknowledged SIGKILL, then takes the common reopen/revoke/refuse tail. This removes a choice between two unrelated initial local stop locations.
- Both use one frozen producer binary, identical provision/programs and schedule after owner loss. Distinct branch tags live only in protocol metadata and cannot be inserted into retained program state or semantic requests if exact local equality is intended.

A common **pre-execution** seed Store copied before either branch runs can eliminate cross-run namespace/root/key/digest differences. Fixed concrete endpoint URLs can also eliminate program-ID changes. This permits a much smaller comparator: independently validate each raw side, compare parsed local exports structurally without identity substitution, preserve only schema-defined SQL-row ordering equivalence, and separately certify opposite world outcomes. Retain the broader alpha-comparison design as a separately declared fallback experiment, not a post-hoc response to failed exact equality.

Conditions and scope costs:

1. Seed cloning is fixture initialization, not editing captured evidence. Freeze/hash the seed and producer before execution and record each branch's copy provenance. Use distinct writable destination Store paths. Never rewrite a captured Store or exported JSON to match another branch.
2. Copy a transactionally complete closed seed through a verified supported snapshot/copy procedure. Merely copying a main SQLite file while relevant WAL content exists is not a valid cloning method. This read-only task has not verified which seed-backup procedure is available in the runtime.
3. Common namespace and identities intentionally represent alternative histories of one starting computation. State that result explicitly: this is a controlled forked fixture, not two independently provisioned production deployments, concurrent-owner safety, or a distributed recovery theorem.
4. Endpoints need separate pristine backing state per branch. A fixed URL can be realized by sequential endpoint processes with the same bound address and fresh endpoint Stores. Each ready record binds its actual PID/token/Store. Fully terminate/wait all prior owner/endpoint/transport processes before reusing that address so no lingering packet can contaminate the next branch. Concurrent branches with the same concrete URL and shared endpoint state would invalidate the isolation assumption.
5. PID/token/temporary-path identity still differs in protocol evidence. It remains separately joined inside each branch and outside the declared equal-local observation; do not erase it from raw records. If such identity leaks into retained local tables, exact equality fails and that is reportable evidence.
6. The common seed does not guarantee same local fact order or scheduler state: record the common admitted prefix and compare stopped, reopened and final exports. The shared handshake likely makes a strict match feasible, but this has not been experimentally established.
7. A fresh applied branch needs its own provenance and causal validation. The historical sealed branch remains an unchanged anchor/control; the old anchored checker should continue to reject the new protocol rather than being relabeled as supporting it.

Preferred claimed result, if observed: "Two isolated executions forked from the same sealed pre-execution Store and fixed endpoint configuration have equal declared local exports after shared-owner loss/reopen, exact-root revocation and fresh refusal, while complete endpoint evidence records zero versus one application." Its scope is narrower and its equality evidence stronger than a result obtained through a large unexamined hash-renaming layer.
