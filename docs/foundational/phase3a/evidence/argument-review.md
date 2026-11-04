# Phase 3A independent argument review

Scope: read `AGENTS.md`, `00-charter.md`, `01-contract.md` and `03-results.md`.
No implementation inspection, model execution or tests. This is an argument
review, not qualification of the unfinished experiment.

**Spec verdict: acceptable for the declared first slice.** The two profiles
isolate one missing edge contract while keeping identity, authorization and
local durability fixed. The two opaque outputs and one A-only interruption
suffice for the proposed negative witness. The single-slot channels and fixed
delegation do not substantiate unrestricted transport or composition claims.

**Argument verdict: sound conditional counterexample and custody induction;
no blocking logical defect found.** Executable correspondence, actual shortest
trace and finite-graph results remain unreviewed and unestablished by this review.

## Important findings and qualification conditions

1. **The falsified implication must retain its weak antecedent.**
   `03-results.md:5-14` correctly distinguishes preservation of committed local
   state and at-most-once application from end-to-end recovery correctness.
   The witness does not show that two agents individually satisfying the latter
   obligation cease to satisfy it under composition. Describing it merely as
   “two correct durable agents fail” would lose the essential contract boundary.
   B's permission to release on delivery is a modeled choice, not an established
   current-runtime behavior. The supplied documents explicitly disclaim a
   current-runtime defect; this review does not independently verify their Store
   correspondence statements.

2. **Loss of custody and unrecoverability are different witnesses.**
   B's release already falsifies `D_B = r OR D_A = r` while A may still hold r
   in memory. Useful recovery fails only after the last accessible copy is lost.
   The latter witness therefore needs the empty result channel, erased volatile
   receipt, absent durable A receipt and non-reconstructive B tombstone described
   in `03-results.md:18-38`. A custody violation alone is insufficient evidence
   of irrecoverability. Preserve both predicates when reporting model results.

3. **Indistinguishability must cover the complete recovery interface.**
   The two executions must agree on every accessible field and future observation
   source: control state, tombstones, consumption markers, channels, request,
   identities and any retained metadata. Removing only an oracle field is valid
   only if neither transitions nor recovery can read it. Result-dependent hashes,
   logs, status flags or an external lookup can distinguish worlds. Under the
   explicitly stated opaque-result/no-new-observation assumptions, the argument
   is valid: the same observable recovery behavior cannot guarantee two distinct
   correct results. An honest unknown remains safe but misses useful recovery.

4. **The induction is sufficient, not a necessity or implementation theorem.**
   `03-results.md:42-60` handles the destructive transition correctly: authenticated
   ACK provenance supplies A's earlier durable commit, and monotonic retention
   makes that fact persist until B releases. This is not circular when the
   ACK-send guard is separately realizable using A's own durable state and B
   uses only the received ACK. A model that lets B inspect A's private state or
   directly guards release with the custody invariant would instead assume the
   result operationally. Other protocols can preserve recovery without this
   specific ACK discipline; no minimal/necessary protocol claim follows.

5. **Finite evidence must remain separated from the proof sketch.**
   Exhaustive reachability can qualify the frozen finite graph. A shortest
   counterexample is shortest only for its initial state, failure predicate and
   transition granularity. Two values suffice for the negative argument; they
   do not establish payload-parametric implementation correctness. Existence
   of a consumption suffix from every state is not universal fair termination.
   The documents correctly leave scheduler fairness, resource adequacy,
   additional crashes, dynamic authority and arbitrary composition open.

## Prior knowledge and remaining scope

The delivery-versus-application acknowledgment distinction is established:
Saltzer, Reed and Clark explicitly discuss that boundary in “Delivery
guarantees,” fifth PDF page of [End-to-End Arguments in System
Design](https://web.mit.edu/6.1800/www/readings/papers/endtoend.pdf), inspected
during this review. The proposed contribution is an ALAN-specific falsifiable
contract and retained witness, not a novel general impossibility theorem.

No counts, runtime qualification, formal mechanization, B-effect atomicity
implementation, liveness theorem, distributed refinement or manuscript
consistency check is established by this review. Atomic B application/receipt
recording and durable-store preservation are premises throughout.
