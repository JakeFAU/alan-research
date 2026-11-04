# Independent review 2: model completeness and witness minimality

Disposition: **APPROVE** for the frozen Task 1 operational model and graph-analysis scope. No blocking finding. This does not approve the separate ledger, runtime, manuscript/PDF, final source qualification, or milestone as a whole.

## Reviewed boundary

Read the task brief, development report, supplied full task-1 patch, frozen contract, operational sources/tests/exporter/README, current results argument, and unchanged `research/mutualobligations/model.go`. The applicable enumeration domain is exactly S1/S2/S3 × Local/Explicit/Missing/ScopeMix plus S2/Shared, S2/LateGrant, S2/Revocable, and Cartesian/ScopeMix: sixteen graphs. A services while waiting; B is serial. There are two originals, two opaque result values per original, fixed logical generation 1, three one-slot channels per original, at most one recoverable crash per physical agent, one fixed possible explicit grant, and at most one B:b revocation.

Reviewed source SHA256:

- `research/authorityservice/model.go`: `11ec7401ade2c4b5b942b70929cc5df8a8f47da4d072a0b5d980b86e19034dbe`
- `research/authorityservice/explore.go`: `3916f87834709dcd7c0d55b2470daf6317fb9ab2ac637a01bcc0723a7eb0e26e`
- `research/authorityservice/model_test.go`: `c8f1cc61eee61ef4a27f6f26402c86826c0c1e29c83096c37ca1fb6b3048abeb`

## Findings and challenged failure modes

1. **Inherited schedules retained.** Source inspection and a new exhaustive projection probe checked every reachable state's base successors against the extension's successors. Each base transition appears exactly once, except Apply without the required new permit, which appears zero times. Added transitions stutter in the base state. Both physical agents retain crash and recovery transitions; crashing an agent clears the incoming service permit while the unchanged base step clears that agent's outgoing volatile result. Neither direction's availability was accidentally modeled as two independent role processes.
2. **No compressed-graph loss detected.** Full states are comparable map keys; the BFS grows until its queue is exhausted, with no depth cutoff. Every generated arc is retained, including duplicate/self-loop arcs. The largest observed outdegree is 17, below the uint8 count capacity; state/arc/parent indexes are also comfortably within their selected integer widths. Event-code fields occupy separate bits and every generated event is compared with its reconstruction by Explore. Fixed configuration plus state determines all reconstructed grant and decision descriptors in this deliberately fixed universe. That representation does not establish adequacy for arbitrary dynamically varying grant records.
3. **Reverse completion and cycles independently cross-checked.** The new probe uses full-state adjacency rather than compact event codes, reconstructs predecessor reachability, and uses Kosaraju instead of the implementation's Tarjan SCC traversal. State count, arc count, raw completable count, safely completable count, cyclic-state count and cyclic-component count agree in all sixteen graphs. Retry cycles remain present, so existential completion does not imply inevitable completion.
4. **Unsafe prefixes cannot merge away.** Unsafe is part of the state key and is monotone over every inspected arc. Safe reverse traversal begins only at safe complete states and never traverses an unsafe predecessor. Both unsafe graphs have raw completion from every state but zero safely completable states, including their initially safe prefixes: eventual BA completion requires crossing the unsafe boundary. This is the correct distinction between a safe prefix and a safe completion suffix.
5. **Minimal witnesses are vocabulary-relative and correct.** Independent BFS minima in S3/Missing and Cartesian/ScopeMix are 4 events for unsafe authorization, 5 for unsafe admission and 6 for unsafe application, equal to the exported named-cut traces. A lower-bound argument is also available: request commitment precedes send, send precedes request admission, admission precedes first authorization, authorization precedes physical admission, and Apply requires that admission. BA alone realizes these bounds because A services while waiting. Crashes, retries, or an AB prefix only add steps. These are minimal authority-boundary witnesses in this model, not minimal histories containing both reciprocal commitments, and not a metric-independent theorem about runtime traces.
6. **Blocking schedules remain visible.** The S2/LateGrant graph retains 3,852 raw-blocked states and S2/Revocable retains 26,883. Inline grant issuance after revocation still permits the request commitment, recording rejected issuance. Late retention does not replace a refused first decision. Pre-revocation admission survives revocation until application or server crash; subsequent fresh admission rechecks current authority. The target tests address these distinctions, and the exhaustive probe retains all corresponding base schedules.
7. **Previous baseline unchanged.** `git diff --exit-code 1d41d1e20b579443c98a6caa53e51036e55f2d82 -- research/mutualobligations docs/foundational/phase3b1` returned 0 with no diff. The unchanged mutualobligations README records 68,497 serial/serial states, with 144 unable to complete. I did not rerun that prior suite in this review: its result is prior evidence, and final regression qualification remains the parent's task. The enlarged sixteen graphs are not claimed to reproduce the old 144-state count.

## Limits of this approval

The enumeration is exhaustive over the stated reachable relation, not all permissions, policies, service placements, graph histories, or adversarial ledger mutations. Shared and revocation are separate configurations; the shared-root-after-b-revocation test explicitly injects that state and is a predicate control, not a seventeenth enumerated combined scenario. Fixed valid grant descriptors and their chosen path priority are construction assumptions. The operational violation labels are not independent authority checking. The new probe also shares Successors with the implementation and does not validate semantic completeness against arbitrary external systems; its independent contribution is transition projection and graph-analysis verification. Authority ancestry/malformed histories remain the separately authored ledger's responsibility.

The atomic Apply/receipt abstraction, shared-resource opaque values, lack of arbitrary redelegation, and single-host/runtime correspondence limits are stated in the contract/results. No universal theorem, generic recovery equivalence, fair-scheduler termination, or N-agent result follows from these graph counts.

## New decisive experiment

Executed, exit 0:

```sh
GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build GOPROXY=off GOSUMDB=off go run /tmp/alan-phase3b2-review2.go > /tmp/alan-phase3b2-review2.log
```

The probe enumerated 1,198,963 states and 9,155,940 arcs across the sixteen graphs. It did not modify operational source. Existing normal/race/Staticcheck results are attributed to the development report, not rerun or claimed as new review evidence.

Probe SHA256: `f2438ee9731ba25896ddd23f29693c6e0c50463610505dd5497a9b35745405d3`.
Log SHA256: `b4c55433d96cf98541a36988e32fc30d530e7ffc579a959f50c63143c5f803e1`.
The exact probe and output follow so this evidence does not depend on temporary-file retention.

### Probe source

```go
package main
import (
 "fmt"
 "os"
 a "github.com/JakeFAU/alan-research/research/authorityservice"
 b "github.com/JakeFAU/alan-research/research/mutualobligations"
)
func must(ok bool, msg string) {if !ok {panic(msg)}}
func main(){
 must(len(a.Configurations())==16,"configuration count")
 for _,c:=range a.Configurations(){check(c)}
}
func check(c a.Config){
 states:=[]a.State{a.Initial(c)}; ids:=map[a.State]int{states[0]:0}; depth:=[]int{0}; edges:=[][]int{}; incoming:=[][]int{nil}; maxOut:=0; arcs:=0; minima:=map[string]int{}; crashes:=[2]int{}; recoveries:=[2]int{}
 for i:=0;i<len(states);i++ {
  s:=states[i]; next:=a.Successors(c,s); baseNext:=b.Successors(b.ServicingSerial,s.Base); seenBase:=make([]int,len(baseNext)); out:=[]int{}; if len(next)>maxOut {maxOut=len(next)}
  for _,tr:=range next {
   e:=tr.Event; matched:=false
   for j,bt:=range baseNext {
    if e.Kind==string(bt.Event.Kind)&&e.Actor==bt.Event.Actor&&e.Obligation==bt.Event.Obligation&&e.Result==bt.Event.Result&&e.Duplicate==bt.Event.Duplicate&&tr.State.Base==bt.State {seenBase[j]++; matched=true}
   }
   if !matched {must(tr.State.Base==s.Base,"new authority transition changes projection")}
   if e.Kind=="crash" {crashes[e.Actor]++; must(!tr.State.Authority[1-e.Actor].Permit,"crash kept server permit")}
   if e.Kind=="recover" {recoveries[e.Actor]++}
   must(s.Unsafe&tr.State.Unsafe==s.Unsafe,"unsafe prefix erased")
   if e.Violations!=0 {d:=depth[i]+1; if old,ok:=minima[e.Kind];!ok||d<old {minima[e.Kind]=d}}
   j,ok:=ids[tr.State]; if !ok {j=len(states);ids[tr.State]=j;states=append(states,tr.State);depth=append(depth,depth[i]+1);incoming=append(incoming,nil)}
   out=append(out,j);incoming[j]=append(incoming[j],i);arcs++
  }
  for j,bt:=range baseNext {want:=1;if bt.Event.Kind==b.Apply&&!s.Authority[bt.Event.Obligation].Permit {want=0};must(seenBase[j]==want,fmt.Sprintf("base edge mismatch %s/%d/%s",c.Name(),i,bt.Event.Kind))}
  edges=append(edges,out)
 }
 // Kosaraju independently checks the compact export's Tarjan counts.
 visited:=make([]bool,len(states));order:=[]int{};var dfs func(int);dfs=func(v int){if visited[v]{return};visited[v]=true;for _,w:=range edges[v]{dfs(w)};order=append(order,v)};for i:=range states{dfs(i)}
 visited=make([]bool,len(states));cycleStates,cycleComponents:=0,0;var back func(int)int;back=func(v int)int{if visited[v]{return 0};visited[v]=true;n:=1;for _,w:=range incoming[v]{n+=back(w)};return n}
 for k:=len(order)-1;k>=0;k--{v:=order[k];if visited[v]{continue};n:=back(v);cyclic:=n>1;if n==1 {for _,w:=range edges[v]{if w==v {cyclic=true}}};if cyclic{cycleStates+=n;cycleComponents++}}
 // Boolean predecessor closure over full adjacency, independent of compact codes.
 reach:=func(safe bool)int{v:=make([]bool,len(states));q:=[]int{};for i,s:=range states{if a.Complete(s)&&(!safe||s.Unsafe==0){v[i]=true;q=append(q,i)}};for k:=0;k<len(q);k++{for _,p:=range incoming[q[k]]{if !v[p]&&(!safe||states[p].Unsafe==0){v[p]=true;q=append(q,p)}}};return len(q)}
 raw,safe:=reach(false),reach(true);r:=a.Explore(c)
 must(r.States==len(states)&&r.Arcs==arcs&&r.CanComplete==raw&&r.CanCompleteSafely==safe&&r.CyclicStates==cycleStates&&r.CyclicComponents==cycleComponents,"independent graph mismatch")
 for _,kind:=range []string{"authorize","dispatch-admit","apply"}{if d,ok:=minima[kind];ok{want:=map[string]int{"authorize":4,"dispatch-admit":5,"apply":6}[kind];must(d==want,"nonminimal witness");must(len(r.Cuts["unsafe-"+kind].Trace)==d,"export witness not minimal")}}
 must(maxOut<256,"compact count overflow");must(crashes[0]>0&&crashes[1]>0&&recoveries[0]>0&&recoveries[1]>0,"missing physical crash/recovery")
 fmt.Fprintf(os.Stdout,"%s states=%d arcs=%d raw=%d safe=%d cyclic_states=%d cyclic_components=%d max_out=%d crashes=%v recoveries=%v unsafe_minima=%v PASS\n",c.Name(),len(states),arcs,raw,safe,cycleStates,cycleComponents,maxOut,crashes,recoveries,minima)
}
```

### Probe output

```text
S1/Local states=103361 arcs=793736 raw=103361 safe=103361 cyclic_states=85796 cyclic_components=33632 max_out=16 crashes=[28585 28585] recoveries=[25519 25519] unsafe_minima=map[] PASS
S1/Explicit states=6741 arcs=37434 raw=0 safe=0 cyclic_states=4846 cyclic_components=2258 max_out=12 crashes=[1855 2261] recoveries=[1687 2219] unsafe_minima=map[] PASS
S1/Missing states=6741 arcs=37434 raw=0 safe=0 cyclic_states=4846 cyclic_components=2258 max_out=12 crashes=[1855 2261] recoveries=[1687 2219] unsafe_minima=map[] PASS
S1/ScopeMix states=6741 arcs=37434 raw=0 safe=0 cyclic_states=4846 cyclic_components=2258 max_out=12 crashes=[1855 2261] recoveries=[1687 2219] unsafe_minima=map[] PASS
S2/Local states=103361 arcs=793736 raw=103361 safe=103361 cyclic_states=85796 cyclic_components=33632 max_out=16 crashes=[28585 28585] recoveries=[25519 25519] unsafe_minima=map[] PASS
S2/Explicit states=103361 arcs=793736 raw=103361 safe=103361 cyclic_states=85796 cyclic_components=33632 max_out=16 crashes=[28585 28585] recoveries=[25519 25519] unsafe_minima=map[] PASS
S2/Missing states=6741 arcs=37434 raw=0 safe=0 cyclic_states=4846 cyclic_components=2258 max_out=12 crashes=[1855 2261] recoveries=[1687 2219] unsafe_minima=map[] PASS
S2/ScopeMix states=6741 arcs=37434 raw=0 safe=0 cyclic_states=4846 cyclic_components=2258 max_out=12 crashes=[1855 2261] recoveries=[1687 2219] unsafe_minima=map[] PASS
S3/Local states=103361 arcs=793736 raw=103361 safe=103361 cyclic_states=85796 cyclic_components=33632 max_out=16 crashes=[28585 28585] recoveries=[25519 25519] unsafe_minima=map[] PASS
S3/Explicit states=103361 arcs=793736 raw=103361 safe=103361 cyclic_states=85796 cyclic_components=33632 max_out=16 crashes=[28585 28585] recoveries=[25519 25519] unsafe_minima=map[] PASS
S3/Missing states=103361 arcs=793736 raw=103361 safe=0 cyclic_states=85796 cyclic_components=33632 max_out=16 crashes=[28585 28585] recoveries=[25519 25519] unsafe_minima=map[apply:6 authorize:4 dispatch-admit:5] PASS
S3/ScopeMix states=6741 arcs=37434 raw=0 safe=0 cyclic_states=4846 cyclic_components=2258 max_out=12 crashes=[1855 2261] recoveries=[1687 2219] unsafe_minima=map[] PASS
S2/Shared states=103361 arcs=793736 raw=103361 safe=103361 cyclic_states=85796 cyclic_components=33632 max_out=16 crashes=[28585 28585] recoveries=[25519 25519] unsafe_minima=map[] PASS
S2/LateGrant states=111065 arcs=840616 raw=107213 safe=107213 cyclic_states=91470 cyclic_components=36234 max_out=16 crashes=[30705 31169] recoveries=[27447 28055] unsafe_minima=map[] PASS
S2/Revocable states=220564 arcs=1740832 raw=193681 safe=193681 cyclic_states=181896 cyclic_components=71960 max_out=17 crashes=[61232 61580] recoveries=[54836 55292] unsafe_minima=map[] PASS
Cartesian/ScopeMix states=103361 arcs=793736 raw=103361 safe=0 cyclic_states=85796 cyclic_components=33632 max_out=16 crashes=[28585 28585] recoveries=[25519 25519] unsafe_minima=map[apply:6 authorize:4 dispatch-admit:5] PASS
```
