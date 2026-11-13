# Task 2 report

Status: DONE pending fresh independent task review.

Base `9a28d2ca4a4e031f1c98860154bbf0ed4aa7fb6e`.
Task commit `77f897e0c8bde7c1b3ea736a19f26029debb588e`, tree `a6858fbd877081003678fb98351359d6ce509d24`.
Only the controller's pre-existing Task 1 plan-checkbox edits remain dirty. No model, other task, production, dependency, paper or old-package files were changed.

## Question and result

The independent ledger checks whether historical admission/current refusal, external truth and local authentic custody stay separated across one loss/reopen and exact b revocation. It agrees with all 63 operational states and 104 arcs, including 33 multiply reached targets and 42 additional incoming arcs. It independently derives 31 Store classes and 52 dossier classes. These counts were computed by this worker's executed audit before receiving the controller's matching independent count.

The report preserves shortest qualifying pairs 25/49 (lengths 4/7, zero/v1) and 49/54 (7/7, v1/v2), negative cut 13 followed by ingress/application through 22 to 28 (3/4/5), and recovery target 61 via a final consume from reopened/revoked source 57 (length 9). It inspects qualifying arcs rather than assuming a shortest target prefix has the right chronology. The report states exact pair/individual minimization and event-order tie-breaking.

The audit checks descriptor completeness, every prefix/edge and merged ledger; separately enumerates valid descriptor shapes to reject missing enabled edges; rejects duplicate outgoing events/states and malformed targets/prefixes; compares independent observations with operational projections. The interpreter never calls operational initial/successor/replay/authority/classification/projection functions. The equality comparison is outside the ledger interpreter.

## Interfaces for Task 3

- `AuditTrace(events []Event) (Ledger, error)` reconstructs a ledger from the fixed admitted cut; errors wrap the failing event index/kind and reason.
- `Ledger.StoreView() StoreView` and `Ledger.DossierView() DossierView` independently derive the existing comparable observation types. Task 3 may compare these to the operational witness cut without importing operational decision logic into the ledger.
- `Ledger` is comparable and separately represented: `PacketSent`, `Contact`, `Opportunity`, `Dead`, `Restarted`, `Revoked`, `CValid`, `Refused` are booleans; `World`, `Inbox`, `Stored`, `Used` are `Result`. It represents a semantic cut, not an append-only history (loss forgets old volatile response custody). Constructing one directly does not validate it.
- `AuditGraph(graph Graph) (AuditReport, error)` returns checked state/arc/prefix counts, merged-target/additional-arc counts, Store/dossier class counts, declared search criterion, pairs `ZeroOne`, `DistinctValues`, `DelayedApplication{Cut,Later}`, and `RetainedRecovery`.
- `Witness{StateID int, Events []Event, Store StoreView, Dossier DossierView}`.
- `AuditControls() ([]ControlResult,error)` runs 9 lawful and 11 hostile controls; the new no-argument `cmd/audit` emits graph report and controls as deterministic JSON. Usage/write errors are tested.

## Executed evidence

Full command/environment/exit metadata is in `docs/foundational/phase3b6/evidence/task2/commands.json` and `verification.md`; source identities in `source-sha256.txt`; all evidence-file digests in `SHA256SUMS`.

- Genuine initial RED: compiled empty-ledger/nil-error stub failed 8 lawful and 14 hostile assertions. Original RED source/test bytes and log are compressed in evidence. Initial focused GREEN passed after ledger implementation.
- Final `go test -count=1 -v ./research/effectobservation/...` passed, including all descriptor mutations, state-pair partition comparisons, malformed graph tests and command tests.
- `go run ./research/effectobservation/cmd/audit` twice and `cmd/export` twice passed; direct byte equality confirmed for each pair. Uncompressed audit SHA-256 `20a5479244fc69f1655af2535d7f2adb3aac451136e76c49da1eeb3047fc26bd`; graph unchanged from Task 1, `142ae6cd3ce9622ecd535d15ad67b923983d8991e357da596858c86239be8d15`.
- Four disposable independent-ledger mutants compiled successfully and were killed by targeted semantic assertions: packet deletion on loss; application-as-retention; old-response acceptance after reopen; c revocation. All four targeted tests passed on the unmodified compiled binary first. Exact replacements, mutant source hashes, command regexes and ordinary failing assertions are preserved; runner was rerun after artifact formatting correction.
- `gofmt -l research/effectobservation` and scoped Staticcheck passed with empty output. Go 1.27.4 darwin/arm64 and Staticcheck 2026.2.2 (0.8.2), no drift.
- Initial staged-diff check flagged whitespace in unified patch context markers. Preserved its exit-2 diagnostic; changed saved mutation patches to zero-context without changing Go sources, reran mutant runner, and final `git diff --cached --check` passed.
- All 26 intended evidence members were checked as tracked; every gzip decompressed; source/evidence hashes verified before commit and evidence hashes reverified afterward. Scoped Git escalation was used for final staging/commit after sandbox denied metadata write. No automatic-approval rejection occurred.

No broad, race, build/module, legacy-export or concrete-corpus qualification was performed by this worker; the controller owns those checks.

## Self-review and remaining limits

Inspected complete staged code/docs/evidence diff, gzip contents, manifests, task-only scope, untouched model and unstaged controller plan diff. Review focused on forbidden helper reuse, all-edge completeness, field validation, vanished volatile history at merged states, retained-custody recovery, and history-qualified shortest witnesses. No operational-model defect was found. Four required semantic mutants fail assertions rather than compile or panic failures. Preserved lawful controls detect overrestriction.

The observation note distinguishes finite checks and conditional arguments from older W3/Phase 2B/2C results. Complete dossier evidence is a precondition; missing evidence is not zero. Positive endpoint evidence is not import authority; a refused fresh gate does not decide the original world result. The finite pair is not a matched sealed v2 runtime pair. No completion theorem, new runtime campaign, unrestricted impossibility result or whole-pair recovery result is claimed.

Task 3's source-conditioned concrete relation, controller broad qualification, Task 4 manuscript reconciliation and independent review remain outstanding. The ignored workflow report is left for controller archiving; no controller ledger or plan file was edited by this worker.
