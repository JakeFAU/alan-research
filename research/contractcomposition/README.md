# Named progress contracts (Phase 3C.1)

This fresh research package studies the fixed chain A→B→C and cycle A→B→C→A
under all eight serial/servicing masks. It retains named U/W/P/D original
registers, bindings and physical-node policy incidence through grouping. It is
not a production runtime, a full microstate-product exploration, or arbitrary
ALAN congruence.

The concrete adapter reuses only unchanged `research/twoagent` APIs. It preserves
request/result/ACK retries and duplicate deliveries. Every physical owner has
its own availability; serial waiting gates only incoming application. Lifting
uses retained requests and durable result custody, never the external oracle.
The old independent trace ledger checks each concrete witness's renamed edges.

Run from the module root:

```sh
go test -count=1 ./research/contractcomposition/...
go test -race -count=1 ./research/contractcomposition/...
staticcheck ./research/contractcomposition/...
go run ./research/contractcomposition/cmd/export
```

The exporter writes deterministic JSON on stdout. The complete abstract graphs
contain every arc, not just BFS predecessor trees. It includes all grouping
interfaces, concrete R0/R1 histories, the serial-cycle control, pre-record loss
and release, and an actual retry loop. Counts are computed from exploration.
The versioned field contract is
`docs/foundational/phase3c1/02-evidence-schema.md`; the direct projection/lifting
argument and exact factored coverage are in `01-argument.md`.

`Config.Validate`, `Component.Validate` and `ConcreteState.Validate` check the
fixed representation/premises. The latter inspects retained records, not the
world oracle, and does not decide global reachability. `Successors`,
`ConcreteSuccessors`, `Project`, `ProjectEdge`, `Complete` and `Dependencies`
operate on declared, valid states/configurations; use constructors and checked
witnesses at input boundaries. `Lift` validates the source's retained premises
and action, but its universal guarantee still assumes a reachable source.
`ServiceGuard` reports only policy eligibility; admission and callee availability
are separate concrete conditions.

`TestEveryOldEdgeRepresentativeProjectionAndLifting` checks all old durable-ACK
single-edge states with each named cycle edge's relevant callee stages, policies
and availability, plus terminal chain BC. These are factored probes with
canonical surrounding edges, not exhaustive three-edge transport reachability.
