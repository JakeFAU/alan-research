# Phase 3B.6 split effect-observation model

The finite model asks whether one accepted request-bound grant, one historical allowance and one allowed gate distinguish external application from recoverable local custody after exact-root revocation and one shared owner-process loss. It starts after admission of the selected BA service and before its worker launches. A serves B's original request. The endpoint and durable Store survive the owner process. The second B-owned permission c serves as an authority control.

## State and event schema

All states, bindings, events and observations are comparable Go values. JSON uses the tagged snake-case field names in the model. Result symbols are empty for absent, `v1` and `v2` for two opaque authentic values. These symbols carry no runtime receipt bytes.

| State fields | Meaning |
| --- | --- |
| `Binding` | `Original` is `BA/1`, `SourceGeneration` and `SourceRevision` are 1, `Grant` is `grant:BA:b`, `Root` is `B/g1:b`, `OriginalActivation` is 1, and `Wire` is `wire:BA:1`. The symbolic bundle abstracts the correlation between the original request and the distinct serving effect operation. |
| `AcceptedGrant`, `FirstAllow`, `Debit`, `FirstGateAllowed` | Immutable accepted grant, historical allow, one debit and one historical allowed physical gate. Each boolean is true in the initial cut. |
| `Owner`, `Activation` | `original`, `stopped`, or `reopened`. Both principals share activation 1 until reopen, including the retained activation while stopped. Both advance to 2 at reopen. |
| `BRevoked`, `ControlRoot` | Current b revocation and the continuously valid c coordinate. In this single-revocation profile b denotes the complete permission `(increment,counter-B,blue)` and c denotes `(increment,counter-AB,blue)`. This fixed distinction carries no general permission algebra. |
| `FreshRefused` | Whether the reopened owner has recorded one refused fresh gate on b. |
| `Permit` | Volatile launch opportunity held by A in the original incarnation. Initially true. |
| `Launched`, `Ingress` | Historical packet transfer and endpoint ingress. A launched packet may remain outside the owner before ingress. |
| `Applied`, `Response` | Authentic world result and authentic successful arrival in the original owner's volatile custody. |
| `Retained`, `Consumed` | Authentic result durably recorded by A and its once-only consumption by B. Consumption preserves the retained result. |

Every event has `Kind`, the full fixed `Binding`, `Actor`, `Incarnation`, `Root` and `Value`. `Incarnation` names the shared activation used by a local event, including the lost incarnation for loss and the newly installed incarnation for reopen. Endpoint events have incarnation 0. The separate event root names an explicitly targeted authority action. The binding root still identifies the original admission path on every event. An empty field denotes an inapplicable target root or value. Operational replay compares every field against an enabled descriptor and rejects malformed or disabled events.

## Transition relation

The table order is the deterministic successor order. Application enumerates `v1` before `v2`. A live owner is either original or reopened.

| Kind | Actor | Guard and change |
| --- | --- | --- |
| `launch` | A | Original owner holds the permit. Clear the permit and mark the packet launched. |
| `ingress` | endpoint | Packet launched, ingress absent. Mark one ingress. |
| `apply` | endpoint | Ingress present, application absent. Set only the selected authentic world value and carry it in the descriptor. |
| `response` | A | Original owner, application present, response and retention absent. Install the authentic world value in volatile response custody. |
| `record` | A | Original owner holds a response. Transfer its value to durable retention, clearing volatile custody. |
| `consume` | B | Live owner, retained value present, no consumption. Mark that value consumed. |
| `revoke-b` | B | Live owner, b current. Target b and mark only its revocation. |
| `lose-owner` | owner | Original owner. Enter stopped and clear permit and unretained response. |
| `reopen` | owner | Stopped owner. Enter reopened and advance shared activation to 2. |
| `fresh-refuse` | A | Reopened owner, b revoked, retention and prior fresh refusal absent. Target b and mark only fresh refusal. |

Only application can introduce a value. Each later step copies that value from the appropriate custody coordinate, and record requires a successful response in the original incarnation. Loss clears that response, so a late response cannot establish custody in the reopened owner. Unsuccessful HTTP completion is hidden and cannot create custody. Retention and consumption are durable. No rule uses endpoint evidence to import a result.

Revocation changes current b authority without affecting the accepted historical path, c, or any prior opportunity or result. Loss leaves the launched packet and any recorded ingress/application unchanged. Ingress and application guards depend only on the packet's external history. A fresh refusal can therefore occur while an original is applied without a retained result, or while its packet may still apply. The refusal cannot reveal whether that application occurred. These arguments depend on the fixed binding, one gate, one launched packet, honest endpoint, and one owner loss.

## Observations and operational enumeration

`ObserveStore` projects binding, admission history, owner stage and activation, b/c authority, fresh refusal, retained result and consumption. The death/Open protocol supplies owner-stage evidence. This is equality under the declared projection. `ObserveDossier` extends it with ingress and the authentic applied value at a cut with complete endpoint evidence. Both hide the permit, in-transit launch and unretained response. A dossier with missing or partial endpoint evidence requires an unknown or unsupported verdict outside this model projection.

Complete breadth-first enumeration yielded **63 states and 104 arcs**. The longest shortest prefix is 9 events from the post-admission cut. `Graph.States` gives state IDs, `Graph.Arcs[id]` gives each outgoing event and its `To` state ID, and `Graph.Prefixes[id]` gives one shortest trace. The [compressed graph](evidence/task1/graph.json.gz) contains all three arrays. The [summary](evidence/task1/summary.json) records arc counts and complete descriptors for the boundary witnesses.

| Witness | State IDs | Event lengths |
| --- | --- | --- |
| Reopened, revoked and fresh-refused, equal Store projections, no application versus `v1` | 25, 49 | 4, 7 |
| Same local projection with distinct unretained authentic values `v1`, `v2` | 49, 54 | 7, 7 |
| Revoked and stopped with no ingress or application, then ingress and application while stopped | 13 → 22 → 28 | 3, 4, 5 |
| Retained before loss and consumed after reopen with b revoked | 61 | 9 |

The first two pairs select the shortest qualifying state prefixes. The third extends the shortest stopped/revoked cut that still has a packet awaiting ingress, with the two required endpoint steps. That cut's negative evidence is valid at the cut and compatible with later application. The last witness selects an incoming consumption arc from a reopened, revoked source with minimum BFS distance. A shortest prefix to the same target can consume before loss and so would not witness consumption after reopen.

The state and arc checks cover custody, immutable historical decisions, stage separation and shortest-prefix reconstruction. This is an operational enumeration. Independent ledger interpretation, observation partitioning and the sealed-runtime evidence relation are separate work. The finite domain has no retry, lookup/import, repeated ingress, extra owner, separate principal crashes, or reciprocal service completion. Applied-but-unretained cuts exceed the atomic Phase 3B.5 predicate's domain. The enumeration establishes no fairness, universal completion or reverse simulation result.
