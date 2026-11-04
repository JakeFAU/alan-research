# Review 1 — authority semantics, attenuation and Task 1 quality

## Disposition

**Task 1 specification: APPROVE. Task 1 code quality: APPROVE.**
No blocking finding or requested source correction was identified within this
review's frozen scope. This is independent scoped review, not approval of the
whole milestone, runtime correspondence, independent oracle, or final manuscript.
Severity/file/line/reproduction findings: none requiring correction.

Reviewed the task brief, implementation report, complete six-file Task 1 patch,
`01-contract.md`, implementation plan, draft `03-results.md`, and the unchanged
base state/successor source. A direct extraction of all six added source files
from `task-1-review.patch` matched their working-tree bytes exactly. Their SHA256
hashes match the Task 1 report; in particular model.go is
`11ec7401ade2c4b5b942b70929cc5df8a8f47da4d072a0b5d980b86e19034dbe` and
explore.go is `3916f87834709dcd7c0d55b2470daf6317fb9ab2ac637a01bcc0723a7eb0e26e`.
No production or model source was edited in this review.

## Authority argument checked

- Root ownership is explicit in `model.go:224–244`. The only retained accepted
  delegation constructor (`245–248`) fixes B/g1:b as both parent and root,
  B/g1 as issuer, A/g1 as grantee, complete tuple b, original BA binding,
  revision zero and no redelegation. Acceptance tests current parent validity.
  A rejected issuance still accompanies the original Commit (`367–371`).
- The closed reachable-state representation stores grant status and the selected
  path, reconstructing exact immutable descriptors from configuration. This is
  adequate for these sixteen fixed graphs: there is no input transition that can
  mutate the parent, grant set, original, logical generation or operation. It is
  not a general grant-validation API, arbitrary-state validator, or multihop
  implementation. Malformed external histories require the separately reviewed
  ledger. The safety argument must retain this boundary.
- Selection (`273–291`) first tests a full server-owned tuple. The delegated
  branch requires BA, accepted grant status, exact b and current validity. There
  is no combination of component fields in S1/S2/S3. ScopeMix can therefore
  retain a legitimate b grant without satisfying x. The explicitly named
  Cartesian mutant alone chooses that path. Ambient borrowing and tuple
  invention have different violation bits (`293–300`). Neither an inherited
  `Binding.Grant` string nor a result becomes a grant constructor.
- Authorization is enabled only after original request admission and only from
  Undecided (`378–394`). The retained decision/path/refusal is never replaced.
  LateGrant changes only grant status. Exact original binding/generations come
  from the fixed scenario constructor, not the physical activation.
- A fresh admission uses the retained allow and current validity of that same
  path (`395–405`); root revocation removes b and its descendant eligibility but
  does not remove r (`244`). The witness's root revision remains its historical
  zero; Event.Revision reports the current source-state b revision. These are
  intentionally different after revocation. Admission records the serving
  activation; each physical crash clears only the permit it served (`362–365`);
  Apply consumes its permit (`355–360`). An existing opportunity survives
  revocation and a requester crash, but not a server crash. Generic unknown
  adapter outcomes are outside this atomic application/receipt model.
- Base successors are actually invoked (`351`) and only Apply is filtered.
  Authority transitions project to stutters; original custody is not paid by
  refusal. Unsafe bits accumulate monotonically (`348`) and safe reverse
  reachability excludes any unsafe-prefix state. Completion cannot erase an
  earlier violation.

I attempted to shorten the four-event authority counterexample. Within the
frozen vocabulary, any first unsafe authorization requires admitted BA, which
requires SendRequest, which requires Commit. The authorization event is thus
at least fourth. Initial state, commitment (including the legitimate b grant
for ScopeMix), transport, crashes and revocation do not themselves perform the
ambient authorization or invent x. Physical admission and application add
further events. This establishes minimality of that authority cut in this
transition vocabulary, not a general minimum for all delegation systems.

## Additional adverse probes actually executed

The existing complete-graph, race and Staticcheck results are attributed to the
Task 1 report; I did not rerun that suite or treat its counts as proof. I ran a
new public-API probe, exit 0, addressing different physical-role and original-
identity edges:

1. BA commit/send/admit/authorize/dispatch-admit; revoke B:b; crash B; Apply BA
   at A. The already admitted A opportunity survives the requester crash and
   remains authority-safe. Apply consumes it.
2. Continue by recovering B and committing/sending/admitting/authorizing AB.
   AB is refused; its fresh admission and application are absent. The retained
   BA grant and equal operation number 1 cannot authorize the other original.
3. In a fresh Revocable graph, admit BA; crash B; crash A; recover A. Apply is
   absent before fresh admission. Re-admission is at activation 1 and Apply is
   safe while the issuer B is still down, because no revocation occurred.
4. Across all fourteen configurations other than the two intentional unsafe
   controls, BA cannot Apply immediately after request admission before any
   authority decision.

Exact command:

```sh
GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build GOPROXY=off GOSUMDB=off go run /tmp/alan-phase3b2-review1-probe.go
```

Output:

```text
PASS remote opportunity: BA commit/send/admit/authorize/dispatch-admit; revoke B:b; crash B; apply BA at A
PASS original/revision isolation: recover B; AB commit/send/admit/authorize refuses; no admission/apply
PASS grant persistence: BA admitted; crash B; crash A; recover A; no apply; re-admit activation1; apply while B remains down
PASS no pre-authorization application across the 14 nonnegative configurations
```

The complete probe is retained below so the temporary path is not an evidence
dependency. Save this block at the command's path to reproduce it.

```go
package main
import (
 "fmt"
 m "github.com/JakeFAU/alan-research/research/authorityservice"
 b "github.com/JakeFAU/alan-research/research/mutualobligations"
)
func step(c m.Config,s m.State,k string,d b.Obligation,a b.Agent) m.State {
 for _,tr:=range m.Successors(c,s) { if tr.Event.Kind==k && tr.Event.Obligation==d && tr.Event.Actor==a { return tr.State } }; panic(fmt.Sprintf("missing %s/%d/%d",k,d,a))
}
func no(c m.Config,s m.State,k string,d b.Obligation) { for _,tr:=range m.Successors(c,s) { if tr.Event.Kind==k && tr.Event.Obligation==d { panic("unexpected "+k) } } }
func require(ok bool,label string) { if !ok {panic(label)} }
func ba(c m.Config) m.State { s:=m.Initial(c); for _,k:=range []string{"commit","send-request","admit","authorize","dispatch-admit"} {a:=b.A;if k=="commit"||k=="send-request" {a=b.B};s=step(c,s,k,b.BA,a)};return s }
func main() {
 c:=m.Config{Policy:"S2",Scenario:"Revocable"}
 s:=ba(c)
 s=step(c,s,"revoke",b.NoObligation,b.B)
 s=step(c,s,"crash",b.NoObligation,b.B)
 require(s.Authority[b.BA].Permit,"requester crash removed remote permit")
 s=step(c,s,"apply",b.BA,b.A)
 require(s.Base.Obligations[b.BA].Applied && !s.Authority[b.BA].Permit && s.Unsafe==0,"application after requester crash failed")
 fmt.Println("PASS remote opportunity: BA commit/send/admit/authorize/dispatch-admit; revoke B:b; crash B; apply BA at A")
 s=step(c,s,"recover",b.NoObligation,b.B)
 for _,k:=range []string{"commit","send-request","admit","authorize"} { a:=b.B;if k=="commit"||k=="send-request" {a=b.A}; s=step(c,s,k,b.AB,a) }
 require(s.Authority[b.AB].Decision==m.Refused,"BA grant leaked into AB after B recovery")
 no(c,s,"dispatch-admit",b.AB);no(c,s,"apply",b.AB)
 fmt.Println("PASS original/revision isolation: recover B; AB commit/send/admit/authorize refuses; no admission/apply")
 s=ba(c)
 s=step(c,s,"crash",b.NoObligation,b.B)
 s=step(c,s,"crash",b.NoObligation,b.A)
 require(!s.Authority[b.BA].Permit,"server crash failed to clear permit")
 s=step(c,s,"recover",b.NoObligation,b.A)
 no(c,s,"apply",b.BA)
 s=step(c,s,"dispatch-admit",b.BA,b.A)
 require(s.Authority[b.BA].PermitActivation==1,"new admission stale")
 s=step(c,s,"apply",b.BA,b.A)
 require(s.Unsafe==0,"unrevoked grant did not survive issuer outage")
 fmt.Println("PASS grant persistence: BA admitted; crash B; crash A; recover A; no apply; re-admit activation1; apply while B remains down")
 for _,c:=range m.Configurations() {
  if c.Policy=="Cartesian" || c.Policy=="S3"&&c.Scenario=="Missing" { continue }
  s:=m.Initial(c)
  for _,k:=range []string{"commit","send-request","admit"} { a:=b.A;if k!="admit" {a=b.B};s=step(c,s,k,b.BA,a) }
  no(c,s,"apply",b.BA)
 }
 fmt.Println("PASS no pre-authorization application across the 14 nonnegative configurations")
}
```

## Code quality and limits

Checked compact event encoding/roundtrip, BFS predecessor traces, all-arc storage,
reverse completion reachability, Tarjan cycle accounting, CLI exit separation,
writer errors and deterministic JSON ordering. No depth cutoff was found. The
integer widths fit the fixed alphabet and frozen graph sizes; they are not a
capacity guarantee for arbitrary future capability universes. Struct values are
comparable and copied by value, so this model has no mutable grant-set aliasing
or concurrent lifecycle mechanism to audit. Tests include genuine behavioral
and temporal checks, not just expected count snapshots.

The operational invariant checks and the selected path's violation classifier
are shared implementation machinery, explicitly labelled as such. This review
has not independently validated every generated history, hostile retained-state
mutation, or the oracle's rejection coverage. It does not establish arbitrary
redelegation, distributed revocation, independent ownership, shared-counter
isolation, fair-scheduler termination or generic runtime recovery equivalence.
The draft conditional argument stays within the fixed-root, authentic-principal,
intact-state and atomic-receipt assumptions. Final source binding, independent
ledger coverage, runtime evidence and manuscript reconciliation remain separate
milestone obligations rather than failures of this scoped Task 1 review.
