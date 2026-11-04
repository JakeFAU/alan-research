# Independent review 3: safety versus completion reachability

Disposition: **APPROVE** for the frozen model's conditional completion argument and separation of safe/raw completion. No blocking finding in this scope. This approval does not certify final immutable qualification, runtime correspondence, the runtime evidence checker, manuscript/PDF synchronization, or the whole milestone.

## Reviewed scope and identities

Read the brief, frozen contract, current results, operational `model.go` and `explore.go`, retained-history `oracle.go`, unchanged mutual-obligation transition model, relevant tests, task-1 review patch and task-1/task-2 development reports. The argument is confined to the sixteen configurations, two original obligations, A servicing/B serial, fixed roots/generation, one possible one-hop grant, opaque per-original results, bounded physical interruptions, and atomic application plus retained receipt. Source hashes at review:

| File | SHA256 |
| --- | --- |
| `research/authorityservice/model.go` | `11ec7401ade2c4b5b942b70929cc5df8a8f47da4d072a0b5d980b86e19034dbe` |
| `research/authorityservice/explore.go` | `3916f87834709dcd7c0d55b2470daf6317fb9ab2ac637a01bcc0723a7eb0e26e` |
| `research/authorityservice/oracle.go` | `06bbe53e48084a61b3b16fec8a1f096a4e35e3bb607fc58a678ea2f33f6b01dd` |
| `research/mutualobligations/model.go` | `fe31ad17d0633977d51fb8345940f4868d21d08758dda5f7e47c4af5e32da3d8` |
| `docs/foundational/phase3b2/01-contract.md` | `be7987bd4d7b0d9453d7995148a2a8464acd1fc5ea42aef6b1eb286515e21772` |
| `docs/foundational/phase3b2/03-results.md` | `dba9df1211e6ce76556180831b79e8b5fd1a2cc7b55d13b2d03e3f691f361b6b` |

## Constructive argument and new probe

The proof sketch can be made constructive without assuming fairness. Recover each currently down agent, then choose no further crash or revocation. Finish BA first: commit if needed, retain the late grant before any still-undecided authorization when that scenario permits it, transport/admit the request, authorize/admit/apply if not already applied, and obtain a recorded result and consume it. Then finish AB by the same procedure. A's servicing rule permits the first application despite its outgoing wait; BA consumption opens B's serial application guard. An already recorded result can be consumed directly. If no recorded copy exists after an application, release could not have removed the retained copy: release requires a recorded caller copy. A full result channel already contains the authentic result, while request/ACK occupancy cannot obstruct these independent channels. This covers partial progress, released receipts, retained permits and either physical interruption.

A new reviewer-authored executable constructed that suffix from every state classified as completable, using only exported successor transitions to execute chosen steps. It does not call `Explore`, reverse reachability, operational authority predicates, or the ledger validator. It also checks that the proposed noncompletable set contains no Complete state and is closed under every successor. This gives a different progress check from the implementation's predecessor traversal. It still shares the operational transition generator; it is not an independently specified transition model or a mechanized general proof.

Across all 1,198,963 reachable states, the constructive counts agree with the development table. The longest constructed suffix has 22 events, except LateGrant's 23. These are bounds on this chosen finite-model schedule, not minimal lengths, runtime latencies, or universal scheduler bounds. No previously completed suite was rerun.

## Challenges and conclusions

1. **Positive assumptions suffice.** Every Local state and every S2 Explicit/Shared state admits the proposed safe suffix. Inline explicit issuance occurs before any authorization of that original. Under nonrevocation it cannot be rejected, so no hidden refusal defeats the positive theorem. S3/Explicit also has the same positive behavior under its selected explicit path; the prose theorem need not claim this additional case.
2. **Late-grant refusal exactly explains blocking.** In S2/LateGrant, raw completion is possible exactly when BA's first decision is not Refused. There are 3,852 states outside that set. A four-event BA commit/send/admit/authorize prefix already fixes the refusal; later grant retention cannot change it. The probe checks that the blocked set stays closed through grant retention, recovery, retries and every other transition. Conversely, choosing grant retention before a still-undecided authorization constructs completion from every remaining state.
3. **Revocation blocks before refusal is recorded.** In reachable S2/Revocable states, completion is possible exactly when either no revocation has occurred, or both originals are already Applied or still hold a Permit. There are 26,883 states failing that condition. The initial one-event Revoke prefix already blocks eventual completion; a refused authorization or rejected issuance is later evidence of that restriction. No smaller nonempty revocation witness is possible, and the initial state remains completable.
4. **An admitted opportunity is per original, not a global completion reservation.** Revoking after BA admission preserves BA's ability to apply. It does not preserve AB's future admission if AB lacks its own opportunity. A crash of the relevant server clears its unapplied opportunity and can change a completable revoked state into a blocked one. Caller-only death does not clear the server's opportunity. The exact predicate above captures both originals and avoids an unjustified inference from one preserved opportunity to global progress. The existing prose says admitted work survives revocation; it does not make that broader inference.
5. **Safe prefixes and safe completions remain distinct.** `Unsafe` is retained in the state key and monotone; safe reverse reachability seeds only safe Complete states and excludes unsafe predecessors. Consequently it means completion whose entire retained prefix and suffix are safe, not merely future-safe behavior after an earlier violation. S3/Missing and Cartesian/ScopeMix can complete operationally from every reachable state, yet no safe Complete state exists: BA consumption requires application, application requires admission and an allowed decision, and that necessary decision crosses Ambient or Cartesian respectively. A previously unsafe Complete state has a zero-length raw suffix but is still not safely completable. Safe initial prefixes in these cases do not rescue safe completion.
6. **Custody is not payment.** `Complete` uses both originals' actual Consumed flags. The ledger's Paid/Unpaid projection independently uses consumption; refusal, authorization, application, publication, release and recording alone do not pay either original. Released custody remains useful for progress because caller recording precedes release, not because release itself means completion. The shared-resource case retains per-original opaque authenticity, not counter arithmetic or isolation.
7. **Inevitability is correctly unclaimed.** A nonterminal state after request admission can repeat send-request/deduplicate indefinitely, delaying enabled authorization or consumption. Bounded crash counts do not bound these retry cycles. Constructing one suffix proves existence only. No fairness assumption is needed to choose that suffix, and no fair-scheduler termination theorem follows from this review.

## Remaining boundaries

The proposition in `03-results.md` is a conditional source-level argument for this finite relation, supported by bounded enumeration and the constructive probe. It is not a claim about arbitrary grants, multihop delegation, all service placements, hostile state, distributed owners, or unrestricted crash schedules. Recovery is enabled for down agents, transport and custody steps can be chosen, and application/receipt retention is atomic. Generic runtime unknown-outcome handling can preserve safety while losing this modeled completion path. The separately reported runtime-checker history-generation finding R4.1 is outside this model-progress approval and is not waived by it. Final source freeze and immutable evidence remain required delivery work.

## Executed command and retained evidence

Executed once, exit 0:

```sh
GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build GOPROXY=off GOSUMDB=off go run /tmp/phase3b2-review-progress.go > /tmp/phase3b2-review-progress.log
```

Probe SHA256: `5f9e3b5ff72cd635519a0956cbe341e683b151b59bb216725931635e2072f60f`.
Log SHA256: `11356f655b8510d728726591b2a81c1677a8112b52dd6e1e4bf3e1a68222963b`.
Exact source and output are embedded below so reproduction does not require those temporary files.

### Probe source

```go
package main

import (
 "fmt"
 a "github.com/JakeFAU/alan-research/research/authorityservice"
 b "github.com/JakeFAU/alan-research/research/mutualobligations"
)

// Independent closed-form candidate; not graph reverse reachability.
func possible(c a.Config, s a.State) bool {
 switch c.Scenario {
 case "Local", "Shared": return true
 case "Explicit": return c.Policy != "S1"
 case "Missing": return c.Policy == "S3"
 case "ScopeMix": return c.Policy == "Cartesian"
 case "LateGrant": return s.Authority[b.BA].Decision != a.Refused
 case "Revocable":
  if !s.Revoked { return true }
  for d:=b.AB; d<=b.BA; d++ {
   if !s.Base.Obligations[d].Applied && !s.Authority[d].Permit { return false }
  }
  return true
 }
 panic("configuration")
}

// Concrete BA-then-AB strategy, no graph search and no reverse labels.
func finish(c a.Config, s a.State) (a.State, int) {
 count:=0
 step:=func(kind string, d b.Obligation, actor b.Agent) {
  for _,tr:=range a.Successors(c,s) {
   e:=tr.Event
   if e.Kind==kind && e.Obligation==d && e.Actor==actor && !e.Duplicate {
    s=tr.State; count++; return
   }
  }
  panic(fmt.Sprintf("missing %s/%d/%d in %s %+v",kind,d,actor,c.Name(),s))
 }
 for actor:=b.A; actor<=b.B; actor++ {
  if s.Base.Agents[actor].Down {step("recover",b.NoObligation,actor)}
 }
 for _,d:=range []b.Obligation{b.BA,b.AB} {
  caller,server:=b.Agent(d),b.Agent(1-d)
  if s.Base.Obligations[d].Consumed {continue}
  if !s.Base.Obligations[d].Committed {step("commit",d,caller)}
  if c.Scenario=="LateGrant" && d==b.BA && s.Grant==a.NoGrant {step("retain-grant",d,caller)}
  if !s.Base.Obligations[d].Admitted {
   if !s.Base.Obligations[d].RequestChannel {step("send-request",d,caller)}
   step("admit",d,server)
  }
  if !s.Base.Obligations[d].Applied {
   if s.Authority[d].Decision==a.Undecided {step("authorize",d,server)}
   if !s.Authority[d].Permit {step("dispatch-admit",d,server)}
   step("apply",d,server)
  }
  if s.Base.Obligations[d].Recorded==b.NoValue {
   if s.Base.Obligations[d].Volatile==b.NoValue {
    if s.Base.Obligations[d].ResultChannel==b.NoValue {step("send-result",d,server)}
    step("receive-result",d,caller)
   }
   step("record",d,caller)
  }
  step("consume",d,caller)
 }
 if !a.Complete(s) {panic("constructor not complete")}
 return s,count
}
func main() {
 for _,c:=range a.Configurations() {
  states:=[]a.State{a.Initial(c)}
  seen:=map[a.State]bool{states[0]:true}
  raw,safe,blocked,maxSteps:=0,0,0,0
  for i:=0;i<len(states);i++ {
   s:=states[i]; can:=possible(c,s)
   for _,tr:=range a.Successors(c,s) {
    if !can && possible(c,tr.State) {panic("blocked set not forward closed")}
    if s.Unsafe & tr.State.Unsafe != s.Unsafe {panic("unsafe history forgotten")}
    if !seen[tr.State] {seen[tr.State]=true; states=append(states,tr.State)}
   }
   if !can {
    if a.Complete(s) {panic("blocked complete state")}
    blocked++; continue
   }
   end,n:=finish(c,s); raw++
   if end.Unsafe==0 {safe++}
   if n>maxSteps {maxSteps=n}
  }
  fmt.Printf("%s states=%d constructed=%d safe_constructed=%d closed_blocked=%d max_suffix=%d\n",c.Name(),len(states),raw,safe,blocked,maxSteps)
 }
}
```

### Probe output

```text
S1/Local states=103361 constructed=103361 safe_constructed=103361 closed_blocked=0 max_suffix=22
S1/Explicit states=6741 constructed=0 safe_constructed=0 closed_blocked=6741 max_suffix=0
S1/Missing states=6741 constructed=0 safe_constructed=0 closed_blocked=6741 max_suffix=0
S1/ScopeMix states=6741 constructed=0 safe_constructed=0 closed_blocked=6741 max_suffix=0
S2/Local states=103361 constructed=103361 safe_constructed=103361 closed_blocked=0 max_suffix=22
S2/Explicit states=103361 constructed=103361 safe_constructed=103361 closed_blocked=0 max_suffix=22
S2/Missing states=6741 constructed=0 safe_constructed=0 closed_blocked=6741 max_suffix=0
S2/ScopeMix states=6741 constructed=0 safe_constructed=0 closed_blocked=6741 max_suffix=0
S3/Local states=103361 constructed=103361 safe_constructed=103361 closed_blocked=0 max_suffix=22
S3/Explicit states=103361 constructed=103361 safe_constructed=103361 closed_blocked=0 max_suffix=22
S3/Missing states=103361 constructed=103361 safe_constructed=0 closed_blocked=0 max_suffix=22
S3/ScopeMix states=6741 constructed=0 safe_constructed=0 closed_blocked=6741 max_suffix=0
S2/Shared states=103361 constructed=103361 safe_constructed=103361 closed_blocked=0 max_suffix=22
S2/LateGrant states=111065 constructed=107213 safe_constructed=107213 closed_blocked=3852 max_suffix=23
S2/Revocable states=220564 constructed=193681 safe_constructed=193681 closed_blocked=26883 max_suffix=22
Cartesian/ScopeMix states=103361 constructed=103361 safe_constructed=0 closed_blocked=0 max_suffix=22
```
