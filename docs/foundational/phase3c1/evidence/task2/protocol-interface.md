# Durable-ACK protocol interface clarification for Task2

Controller inspected unchanged research/twoagent/model.go Successors on this task. The Phase3C1 contract explicitly imports this protocol. These are interface facts, not new premises or a claim that the new producer is independently verified. The new Python ledger should be derived separately from these declarative facts plus the frozen contract/schema; do not copy the Go branch structure. All fields not named by an effect stay unchanged. New physical-owner availability/serial gates and renamed full bindings are governed by 00-contract/02-schema.

All caller events require the actual caller up. All callee events require the actual callee up. Incoming Apply alone also requires that actual callee servicing or not waiting on its outgoing original. Result-bearing event bytes equal the listed source, and events do not obtain bytes from WorldResult. Pure identity events have empty result. Duplicate delivery is permitted only for admit/deduplicate/receive-result/release and retains its channel; otherwise the consumed channel is emptied.

| Event | Additional precondition | Effect / result source |
| --- | --- | --- |
| commit | Not committed | Set committed. |
| send-request | Committed and request channel empty | Set request channel. No condition on prior admission/application/consumption. |
| admit | Request channel occupied and not previously admitted | Set admitted; deliver request with optional retained channel. |
| deduplicate | Request channel occupied and already admitted | Preserve admitted; deliver request with optional retained channel. |
| apply | Admitted and never applied | Choose r1/r2, atomically set Applied, Retained and historical WorldResult; result is that chosen value. |
| send-result | Retained result nonempty and result channel empty | Fill result channel from Retained. No remote caller volatile/recorded check. |
| receive-result | Result channel nonempty | Set/overwrite Volatile from that actual channel, whether already recorded/consumed or not; optional retained channel. |
| record | Volatile nonempty and Recorded empty | Set Recorded from Volatile; Volatile unchanged. No replacement recording. |
| consume | Recorded nonempty and not consumed | Set Consumed; Recorded and Volatile unchanged. |
| send-ack | Recorded nonempty and ACK channel empty | Fill identity-only ACK channel; no consumption prerequisite. |
| release | ACK channel nonempty | Clear Retained; optional retained ACK channel. May repeat after retained result was already cleared. Tombstone unchanged. |

Crash occurs only while up and not previously crashed; set Down/Crashed, retain Activation0 and clear only the owner's outgoing Volatile. Recover occurs only while down; clear Down and set Activation1 while Crashed stays true. Every actor has exactly one such physical record, shared by its caller/callee roles; no separate per-edge crash event is inserted.

The selected genuine histories alone do not define protocol guards. A validator may additionally check history-level causal facts (one authentic application, ACK from recorded custody, exact binding), independently of the producer's state flags. Record any further uncertainty explicitly rather than accepting events by observed producer deltas alone.
