# Known Adversarial Blind Spots (Read-Only)

> **For Agents:** Read-only reference of failure shapes that previous adversarial passes missed in ALAN.
> - **Do NOT edit, append to, or rewrite this file** unless explicitly instructed by the user.
> - **Not rules or required tests:** Most entries will **not** apply to your current task. Skim when designing an adversarial attack or review; ignore any entry that does not fit the code in front of you.

---

## 1. Error Return in a Shared Selector / Queue Scan (Head-of-Line Wedge)

- **What we missed:** Testing a single failing run or bad input, observing the expected error return, and stopping—without placing a healthy run or valid input *behind* it and running a second cycle.
- **Code shape (`internal/executor/runner.go` `selectReady`):**
  ```go
  for _, id := range ids { // or for _, obs := range observations
      if _, err := Reduce(p, proj, obs, ...); err != nil {
          if errors.Is(err, ErrIneligible) { continue }
          return nil, err // <-- aborts entire scan; cursor resets on next Run()
      }
  }
  ```
- **Why it bites:** If `a-fault` (or `obs1`) stays `"open"` / unconsumed at the front of iteration order and returns a non-`ErrIneligible` error (e.g., `"end"` output schema violation or `"resource limit"`), every subsequent `Store.Run` call hits index `0` and aborts with `0` steps, starving `z-healthy` (or `obs2`) forever.
- **Probe when relevant:** Put a healthy item *after* the faulting item in scan order and call the runner twice.

---

## 2. Existential (`∃`) vs. Universal (`∀`) Static Validation (`author.Compile`)

- **What we missed:** Testing disjoint schemas (`[a] -> [b]` rejected) and valid constants (`"a" -> [a]` accepted), while missing partial overlap (`[a, b] -> [a]`), unconstrained-to-constrained variable copies (`string -> [a]`, `MaxBytes: 100 -> MaxBytes: 5`), and capability-tag erasure (`request:invocation:* -> string`).
- **Code shape (`author/compile.go` `checkOperand`):**
  ```go
  candidates := src.Enum
  if len(candidates) == 0 { candidates = want.Enum } // borrows target enum when src is unconstrained!
  for _, c := range candidates {
      if src.Validate(c) == nil && want.Validate(c) == nil {
          compatible = true; break // proves src ∩ want != ∅, NOT src ⊆ want
      }
  }
  // Also: variable-to-variable string copy never checks src.MaxBytes <= want.MaxBytes
  ```
- **Why it bites:** Proving *some* value can succeed at runtime lets wider or unconstrained variables flow into narrower sinks and fail only at runtime `"end"`. Similarly, collapsing `request:*` string tags to `"string"` and exempting `"string"` in `await` lets arbitrary strings compile as invocation handles and deadlock at quiescence.
- **Probe when relevant:** Test wider-to-narrower (`[allowed, forbidden] -> [allowed]`, unconstrained `string -> [allowed]`, larger `MaxBytes ->` smaller `MaxBytes`, raw `string -> AwaitInvocation`).

---

## 3. Silent Constraint Erasure at CFG Branch Joins (`author.Compile`)

- **What we missed:** Testing single-path straight-line programs or branch joins where one branch is valid, without testing branch joins where incoming edges disagree on constants or field metadata.
- **Code shape (`author/compile.go` CFG merge):**
  ```go
  for name, val := range consts[target] {
      if nv, ok := nextConst[name]; !ok || !val.Equal(nv) {
          delete(consts[target], name) // <-- constraint forgotten instead of joined/widened
      }
  }
  for name, f := range declFields[target] {
      if nf, ok := nextField[name]; !ok || !equalField(f, nf) {
          delete(declFields[target], name) // <-- downstream checkOperand sees neither const nor field!
      }
  }
  ```
- **Why it bites:** If `branch1` sets `out = "invalid_1"` and `branch2` sets `out = "invalid_2"`, the join deletes `out` from `consts` and `declFields`. Downstream `checkOperand` finds no entry in either map and returns `nil`, accepting a program where *every* branch is illegal.
- **Probe when relevant:** Construct a 2-branch `IfEqual` that assigns two different invalid values (or differently bounded fields) to the same variable before a constrained sink or `Success`.

---

## 4. Cross-Gate / Cross-Query Predicate Asymmetry in Multi-Table State Machines

- **What we missed:** Testing awaited or all-accepted child invocations, without checking that *every* gate (`AuthorizeIntent`, `prepareAdapterWork`, `prepareCompositionScope`, `SettleRoot`, `readRootCompletions`) uses the same predicate for **unawaited** or **rejected** children.
- **Code shapes (`internal/executor/invocation*.go`, `composition_*.go`):**
  - *Split-brain on parent exit with unawaited accepted child:*
    - `compositionScopeClosed`: `if root.Projection.Status != "open" { return true }` -> blocks child `prepareAdapterWork` (`Reason: "scope_closed"`).
    - `AuthorizeIntent`: still returns `Allowed: true` for the child.
    - `prepareCompositionScope`: sees `child.Accepted && publications == 0` -> blocks `PrepareRoot` with `Pending: "descendant_outcome"`.
  - *Settlement vs. inspection drift on rejected child (`!child.Accepted`):*
    - `prepareCompositionScope`: checks `composition_publications` only `if child.Accepted` -> `SettleRoot` succeeds.
    - `readRootCompletions`: checks `composition_publications` for *all* `"invocation"` intents without checking `Accepted` -> `PendingInvocations` stays non-empty forever.
- **Probe when relevant:** When multiple functions or SQL queries inspect the same lifecycle/composition entity, compare their exact `WHERE` clauses and guard conditions under non-happy states (parent finished before child; child admission rejected).

---

## 5. Partial State-Tuple Updates & Manuscript-vs-Executor Drift

- **What we missed:** Checking the primary boolean flag of a state transition without checking companion version/revision fields, or assuming features tested in `research/` packages are wired into `internal/executor`.
- **Examples:**
  - `decideTarget` (`internal/executor/lifecycle.go`) sets `after.Retired = true` on `RetireSubject` without incrementing `subjects.revision`.
  - `paper/alan-dissertation.md` §§6.11–6.15 and `research/authorityservice` model `DelegateAuthority`, `RevokeAuthority`, and storage formats 13–16, while `internal/executor/lifecycle.go` rejects those operations (`"unknown lifecycle operation"`) and implements schema versions 3–8.
- **Probe when relevant:** Check every field of the mutated struct/row on authority transitions, and verify claims directly against `internal/executor` rather than sibling `research/` packages.
