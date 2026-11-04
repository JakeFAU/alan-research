# Phase 3B.3 independent progress and obligation audit

Read-only research audit, 2026-11-04. Research HEAD observed `360cae66f5ffdd6e7a95a568552df5522122da21`. Read packet §§11,13,17,19–20, frozen contract v3, Task 3 brief and `/tmp/alan-phase3b3-corpus-interface.md`; inspected unchanged authorityservice, its independent ledger, and inherited mutualobligations model. No runtime implementation/evidence is certified by this audit. No repository source/model edit, commit, subagent or archived SQLite open was performed. Scratch source and output are under `/tmp`.

## Result

The proposed LateGrant conditional completion query is true on the entire unchanged finite graph, and can be strengthened to an exact characterization. It is not a forward-invariant scheduling guarantee and does not imply completion from every runtime state. The proposed correspondence domain restrictions are necessary. In particular separate-grant-plus-revocation, unknown physical outcomes, repeated contacts, resource failures and handshake failures must remain visible outside the direct model mapping where applicable.

### Exact LateGrant query

Let `R` be all states reachable from `Initial(S2/LateGrant)` by unchanged `Successors`; let `C(s)` be unchanged `Complete(s)` (both authentic original results consumed). Let `SafeReach(s)` mean there is a finite path from `s` to a `C` state on which every state has `Unsafe == 0`, including `s`.

The following equivalence was exhaustively checked, without excluding any reachable state:

```text
for every s in R:
  SafeReach(s) iff s.Authority[BA].Decision != Refused
```

| Quantity | Observed |
| --- | ---: |
| Reachable states / all arcs, including loops | 111,065 / 840,616 |
| Safe complete states | 16,384 |
| Safe-completable states | 107,213 |
| States with BA first decision not Refused | 107,213 |
| Predicate/reverse-reachability mismatches | 0 |
| Noncompletable states with BA Refused, no grant yet | 1,926 |
| Noncompletable states with BA Refused, retained grant | 1,926 |
| Maximum shortest safe-completion distance among qualifying states | 23 edges |

The distance is an existential shortest-path bound for this graph, not a scheduler-time bound. The predicate deliberately includes states with no grant yet: it states that future grant retention remains possible, not that covering authority already exists at that prefix. A sufficient-authority formulation must not silently substitute “grant retained now” for “no irreversible first refusal yet.”

Minimal falsification of unconditional completion: `commit BA → send-request BA → admit BA → authorize BA(Refused)`. Appending `retain-grant BA` does not restore completion. Prepending `commit AB` gives the same counterexample with both originals outstanding. The exact JSON shortest counterexample is retained in the probe output. Crash/recovery cannot erase the first decision; the existing `late-grant-sticky-refusal` oracle control covers that extension.

The qualifying set is not forward closed: before a grant, an enabled BA authorize transition leaves it by refusing. The handshake must enforce grant retention before the first actual service permission decision. The graph equivalence alone does not establish that an ordinary runtime program enforces this ordering. Conversely, permanently suppressing the refusal transition would modify the old graph; do not do that.

### Exact revocable query and adverse control

On unchanged S2/Revocable, the LateGrant predicate is insufficient: **22,887 states** have no BA first refusal but cannot safely complete. The shortest such counterexample is a single `revoke B/b` from the initial state. Neither original need have attempted authorization yet.

The following state predicate exactly matched safe reverse completion reachability:

```go
viable := s.Unsafe == 0
for d := AB; d <= BA; d++ {
    viable = viable && s.Authority[d].Decision != Refused &&
        (!s.Revoked || s.Base.Obligations[d].Applied || s.Authority[d].Permit)
}
```

| Quantity | Observed |
| --- | ---: |
| Reachable states / all arcs | 220,564 / 1,740,832 |
| Safe complete states | 32,768 |
| Safe-completable states / predicate states | 193,681 / 193,681 |
| Predicate/reverse-reachability mismatches | 0 |
| Maximum shortest suffix for completable no-BA-refusal states | 22 edges |

This predicate is specific to this frozen graph. AB and BA both require B's `b` root here. Revoking that root can block B's own-root AB action as well as A's delegated BA action. A pre-revocation permit survives revocation; a server crash destroys that volatile permit. Optional crashes can be avoided by an existential suffix. Applied originals retain recoverable custody under this model's atomic application-plus-receipt assumption. Neither B/r nor A/r is exercised by these two original service bindings, so the graph predicate alone does not test unrelated-runtime-action completion after B/b revocation.

## Progress, quiescence and retained obligations

`Complete` means exactly `Base.Obligations[AB].Consumed && Base.Obligations[BA].Consumed` (`mutualobligations/model.go:136`). It does not mean absence of outstanding ACKs, absence of buffered duplicate messages, drained transports, closed principals, an exited runner or a quiescent global runtime. A primary authority gate, an endpoint application, a retained external receipt, a settled server result, a ready message and a lifecycle outcome are each insufficient to mark an original service paid.

Both enumerated graphs have **zero incomplete states with no outgoing transition**. Every complete state also has outgoing transitions (16,384 and 32,768 respectively). Thus a graph-outdegree deadlock/quiescence test would miss all authority-blocked originals and cannot serve as a completion test. Retry/deduplication cycles remain possible. A blocked original may coexist with indefinitely enabled protocol bookkeeping. This is not evidence of a scheduling defect.

The independent ledger's `Paid` is actual original consumption; `Unpaid` is its negation, even for a not-yet-committed original. A report should distinguish absent/not-committed, committed-unpaid, first-refused, fresh-admission-refused, unknown-outcome, retained-result, and consumed. Do not describe both originals as outstanding unless both genuine request commitments are present.

The existential result supplies no inevitable completion assertion, no fairness theorem, and no proof of ordinary runtime scheduler reachability. `Explore.UnfairRetry` supplies the existing refusal-free retry starvation control; its name is intentional, so do not call it a counterexample to an unspecified fair-scheduler property. Runtime `Run` returning or a step/time limit expiring is not model deadlock proof. Positive runtime qualification must show an actual ordinary execution with both original requests committed, BA served by A while AB remains unpaid, and B selecting BA's authentic response before AB service.

## Mapping findings and implementation requirements

1. **Separate grant is LateGrant.** `model.go:335,408,414` separates inline grant configurations from LateGrant and Revocable. BA request commitment and accepted Delegate decision are distinct runtime prefixes. Map the latter, not the ready message or expected case, to `retain-grant`. A rejected Delegate has no direct nonrevoking LateGrant issuance edge. A separate accepted grant followed by revocation has no direct unchanged configuration. Keep `separate_grant_with_revocation_not_in_frozen_configuration`; temporal raw validation can still succeed independently.

2. **Permission and resource decisions must remain distinct.** The permission verdict in immutable `EffectAuthority` can map to model `authorize` only after the authentic request admission and causal selection are established. Source Authorization is not this event. If permission allows but resource reservation fails, the model has no budget-refusal transition. Mapping only the nested permission allow is at most an authority-prefix comparison; preserve the unmodelled resource blocker and do not infer a suffix. Likewise absence of a selector can be an honest runtime refusal even where the model would choose its own root: the explicit-selector runtime policy has a narrower action domain than the model's implicit choice. Treat malformed/wrong selector and binding controls as invalid proposed model histories or independently enforced runtime controls, not arbitrary forward-generated model refusal edges.

3. **Physical opportunities are not applications.** A gate fact maps to dispatch-admit only when its original, retained first permission decision, activation, root and grant all match. A gate is not evidence of endpoint contact. The model consumes one outstanding permit at an atomic apply-plus-retained-receipt edge. It has no applied-with-no-retained-receipt state before legitimate later ACK release. Therefore called-but-no-local-result prefixes, generic unknown adapters and same-epoch repeated idempotent contacts require explicit unsupported/domain labels or a separately justified abstraction. They cannot be filled in from the expected model trace. A repeated admission with no intervening apply/crash has no enabled direct model transition while Permit is true. After model Apply, another admission is also disabled.

4. **Crashes need a declared projection.** The model permits at most one recoverable crash per physical agent, with activations 0→1; runtime epochs and repeated owner reopens have a larger domain. Same-Store process loss affects both open activations. Mapping one actual shared crash/reopen into two per-agent model crash/recover transitions can be an explicitly justified expansion with no intervening service; it must not silently imply independent-owner recoverability or invent one model crash to explain every repeat exposure. Actual empty-reopen epoch facts require controller/boundary evidence.

5. **Runtime handshake stutters preserve additional obligations.** B's lifecycle-response selection and grant-ready transport can be stutters for the two-original model projection, but the real prefix may still await them. A denied/lost ready message, unsupported admission, depleted resource, closed source, or stalled ordinary program can remove a runtime suffix even when its projected state satisfies the finite predicate. Preserve the exact ready/lifecycle/request identities and statuses; do not return a global `CanComplete=true` from the model alone.

6. **Custody event expansion must be explicit.** Runtime durable message admission can establish both delivery and durable local observation in one retained transaction, while the model separates receive-result from record and volatile loss. If projection expands one raw transaction into several model events, identify that expansion and bind all expanded events to the same actual raw evidence. Do not invent a missing actual request/result publication, receipt, record or consumption from the desired model path. Result bytes alone cannot identify an original when results are equal. Both actual caller selections, with exact generation-bound joins, define runtime service completion.

7. **Opaque result abstraction stays bounded.** In Explicit/LateGrant both logical services can touch the same B resource, but the frozen model represents per-original opaque R1/R2 values rather than arbitrary shared-counter interference. Freeze an equality-preserving per-case result abstraction and complete original-ID mapping. A whole-action mapping for adapter Cartesian controls is a declared relation, not coordinate-wise equality with the model's increment/read/red/blue vocabulary.

8. **Negative enforcement is not an accepted invalid history.** Preserve proposed mutated model events, first failed edge, actual ordinary refusal or hostile-store rejection, raw checker, Replay, Open, Inspect and contacts separately. `unsupported` is neither correspondence agreement nor an authority violation. The eight corpus categories do not establish eight ordinary-runtime trace classes. Existing 33-history/four-chain controls and sixteen graphs remain unchanged; this audit reran only the two relevant graphs.

## Executed evidence and reproduction

Scratch files:

- `/tmp/alan-phase3b3-progress-probe.go`: independent BFS using only public Initial/Successors; stores every arc; computes safe reverse reachability from Complete and compares the two predicates against every state; retains shortest counterexample events. It does not call the runtime or alter the model.
- `/tmp/alan-phase3b3-progress-probe.json`: both count/query reports and counterexamples, JSON array.
- `/tmp/alan-phase3b3-progress-ledger.go`: calls unchanged public AuditGraph for these same two configurations.
- `/tmp/alan-phase3b3-progress-ledger.json`: independent ledger reports.

Executed from the research checkout, Go `go1.27.3 darwin/arm64`, exit 0 for both final commands:

```sh
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3b3-progress-go-build go run /tmp/alan-phase3b3-progress-probe.go > /tmp/alan-phase3b3-progress-probe.json
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/alan-phase3b3-progress-go-build go run /tmp/alan-phase3b3-progress-ledger.go > /tmp/alan-phase3b3-progress-ledger.json
```

AuditGraph independently accepted every shortest prefix and outgoing arc: LateGrant 111,065 prefixes / 840,616 edges; Revocable 220,564 / 1,740,832; both had zero unsafe prefixes/edges. This checks the unchanged finite quotient and the separate ledger, not the new raw-runtime checker. No race, runtime, SIGKILL, full sixteen-graph, manuscript build or whole-repository tests were run in this audit.

For Task 3, retain a shortest suffix by storing the next edge when reverse BFS marks a predecessor. Replay the original generated prefix plus suffix through AuditTrace; require both Paid values and AuthoritySafe. Record all excluded predicate-false states/counts, not just qualifying ones. Test unconditional LateGrant completion and the weaker revocable predicate as deliberately failing controls; the strengthened queries should pass. Do not mutate Successors, derive a suffix from an expected final value, or use final runtime state alone as evidence of chronology.

Observed SHA256:

```text
11ec7401ade2c4b5b942b70929cc5df8a8f47da4d072a0b5d980b86e19034dbe  research/authorityservice/model.go
3916f87834709dcd7c0d55b2470daf6317fb9ab2ac637a01bcc0723a7eb0e26e  research/authorityservice/explore.go
06bbe53e48084a61b3b16fec8a1f096a4e35e3bb607fc58a678ea2f33f6b01dd  research/authorityservice/oracle.go
fe31ad17d0633977d51fb8345940f4868d21d08758dda5f7e47c4af5e32da3d8  research/mutualobligations/model.go
f19d9c13f84c17496fabef91cb7575cf52972ff228c0b69e063cdc1c99100e12  docs/foundational/phase3b3/02-contract.md
```

Navigation-only gopls attempts failed because its persistent cache path under `~/Library/Caches/gopls` was not writable; source inspection and the successful Go probes supplied the findings. A first attempted phase3b1 documentation path did not exist; no archived evidence was opened through that path. At audit end an unrelated parent-owned modification to `docs/foundational/phase3b3/03-plan.md` was visible; this audit made no repository edits. This report is an evidence/query design audit, not approval of Task 3 code that has not yet been supplied.
