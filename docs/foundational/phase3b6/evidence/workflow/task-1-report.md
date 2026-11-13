# Task 1 report

Status: DONE pending independent task review. Implementation base `e339c2db337cd5b5a93990b8a42160994a3ec2f4`, tree `ef5a04949ca26f7021f2fc3128bc856599cf8362`.

Task commit `9a28d2ca4a4e031f1c98860154bbf0ed4aa7fb6e`, tree `70503d1c67ca2a21574bd56a7a38880ccb4d32c7`. Worktree clean after commit.

## Result and evidence

63 reachable states, 104 arcs, maximum shortest prefix length 9. Model proofs, exact descriptor and state meanings, and witnesses are in `docs/foundational/phase3b6/02-model.md`. The zero/application fresh-refusal pair is state 25 versus 49. The distinct unretained values use 49 versus 54. The stopped negative cut at 13 extends through ingress 22 to application 28. The retained/consumed control takes the consume arc 57 → 61 after reopen.

The exact offline environment, commands and observed exits are in `docs/foundational/phase3b6/evidence/task1/verification.md`. Focused RED was exit 1 with six ordinary failed assertions against compiled skeleton transitions. Focused GREEN was exit 0. The exporter test first used a nonexistent standard library helper and failed compilation. A local erroring writer fixed it; the complete model/export suite then exited 0. Both graph commands and `cmp` exited 0. Formatting was empty. Default Staticcheck cache was sandbox-denied, and scoped Staticcheck with the writable cache exited 0. Source and output digests are under `evidence/task1`; the uncompressed graph digest is `142ae6cd3ce9622ecd535d15ad67b923983d8991e357da596858c86239be8d15`.

## Go schemas for the next workers

- `Binding{Original string, SourceGeneration uint8, SourceRevision uint8, Grant string, Root Root, OriginalActivation uint8, Wire string}`. Initial values are `BA/1`, 1, 1, `grant:BA:b`, `RootB`, 1 and `wire:BA:1`.
- `State{Binding Binding, AcceptedGrant bool, FirstAllow bool, Debit uint8, FirstGateAllowed bool, Owner OwnerStage, Activation uint8, BRevoked bool, ControlRoot Root, FreshRefused bool, Permit bool, Launched bool, Ingress bool, Applied Result, Response Result, Retained Result, Consumed Result}`. Initial admission fields are true and debit is 1. Initially original, activation 1, `ControlRoot: RootC`, permit true and every other dynamic field zero.
- `OwnerStage` is a string. `OwnerOriginal="original"`, `OwnerStopped="stopped"`, `OwnerReopened="reopened"`. Stopped keeps the durable historical activation at 1. Reopen advances the shared activation to 2.
- `Result` is a string. `NoResult=""`, `Value1="v1"`, `Value2="v2"`. `Root` is a string. `RootB="B/g1:b"` and `RootC="B/g1:c"`. c is continuously valid in this profile.
- `Event{Kind EventKind, Binding Binding, Actor Actor, Incarnation uint8, Root Root, Value Result}`. All kinds use the entire original binding. `Actor` is a string with `ActorA="A"`, `ActorB="B"`, `ActorEndpoint="endpoint"`, and `ActorOwner="owner"`. Incarnation is the local activation, including new activation on Reopen, or 0 for endpoint events. Separate target root is RootB on RevokeB/FreshRefuse and empty on other kinds. Value is populated only for Apply, Response, Record, Consume. Empty value for other events is mandatory.
- `EventKind` is a string. Transition ordering is `Launch="launch"`, `Ingress="ingress"`, `Apply="apply"` with Value1 before Value2, `Response="response"`, `Record="record"`, `Consume="consume"`, `RevokeB="revoke-b"`, `LoseOwner="lose-owner"`, `Reopen="reopen"`, `FreshRefuse="fresh-refuse"`. ActorA emits launch/response/record/fresh-refuse, ActorB consume/revoke-b, ActorEndpoint ingress/apply, ActorOwner lose-owner/reopen.
- `StoreView` has State's Binding, AcceptedGrant, FirstAllow, Debit, FirstGateAllowed, Owner, Activation, BRevoked, ControlRoot, FreshRefused, Retained, Consumed. `DossierView{StoreView, Ingress bool, Applied Result}`. All of these values are comparable. Observations use no hidden transport or volatile response data.
- `Graph{States []State, Arcs [][]Arc, Prefixes [][]Event}`; `Arc{Event Event, To int}`; `Edge{Event Event, State State}`. `Initial() State`, `Successors(State) []Edge`, `ObserveStore(State) StoreView`, `ObserveDossier(State) DossierView`, `Explore() Graph`, and operational `Replay([]Event) (State, bool)`. The independent oracle must interpret event descriptors itself.

## Scope and self-review

The fixed binding abstracts separate request and serving-operation correlation identities. No cryptographic or concrete provenance validation is implied. Revoke-b needs a live owner, and endpoint events require only the packet's external history. Failed HTTP completion is invisible. Death/Open supplies the owner-stage coordinate in the local observation. Task 1 supplies operational enumeration only. Independent ledger interpretation, the sealed concrete relation, broad repository checks and manuscript synchronization remain subsequent work. No old packages, dependencies or outputs changed.

Self-review inspected the full diff and verified gzip contents and manifests, event descriptors, state/arc ordering, and the last witness's history-sensitive consume arc. The final verification log and independent review are controller-owned.
