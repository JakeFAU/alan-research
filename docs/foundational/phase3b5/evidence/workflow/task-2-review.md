### Spec Compliance

- ✅ Spec compliant for Task 2 at `0e8d8f874b5993b62692aee1b872ef7619a93422`, reviewed against base `e6e53deb`. No missing, extra, or misunderstood requirement found in this task's scope.
- ✅ Independent ledger: the new binding/root declarations and exact b-only current-validity rule are local to the checker (`research/authorityservice/oracle.go:101`, `:134`, `:155`). Grant identity, structural validation of rejected claims, and issuance-time validity are separated (`:219`, `:236`). Event interpretation uses the independent ledger methods, not operational authorization/binding/current-validity helpers.
- ✅ Separate issuance remains forbidden on commitment and is restricted to the committed BA original and one attempt (`research/authorityservice/oracle.go:384`, `:425`). First decisions are immutable; fresh admission rechecks the selected authority; application uses an existing activation-bound opportunity; crash clears it (`:341`, `:381`, `:445`, `:464`, `:481`). Revocation leaves historical grants and custody untouched (`:355`).
- ✅ The complete audit checks each shortest prefix, every outgoing arc, and equality of independent ledgers at merged targets (`research/authorityservice/oracle.go:570`, `:597`). The added configuration flag is constant across a graph; omitting that flag from target comparison does not omit mutable state.
- ✅ Reverse reachability is computed from independently enumerated incoming edges and safe complete seeds, then compared against the candidate at every reachable state (`research/delegationcorrespondence/progress.go:75`, `:92`, `:140`). AB refusal counting, all-edge BA refusal preservation, and independently computed SCCs are present (`:154`, `:195`, `:256`). Temporal crash/application witnesses use an actual qualifying edge after a shortest source prefix, so state merging does not erase their required ordering (`:208`).
- ✅ Tests retain lawful old-permit and recovery controls, malformed identity/ancestry/tuple/binding/generation claims, stale activation, exact b/c scope, sticky refusal, fresh refusal, and both authentic consumptions (`research/authorityservice/revocation_corpus.go:39`, `:79`; `research/authorityservice/revocation_oracle_test.go:44`). Progress tests require the crash counterexample and successful custody suffixes (`research/delegationcorrespondence/revocation_progress_test.go:9`).
- ✅ New command selection is opt-in, historical progress fields are omitted for old configurations, and default audit/export commands are untouched (`research/delegationcorrespondence/cmd/progress/main.go:16`; `research/delegationcorrespondence/progress.go:29`; `research/authorityservice/cmd/revocation/main.go:21`). Archived-export comparisons check decoded bytes against the original four archives, not only two current runs (`docs/foundational/phase3b5/evidence/task2/qualify_exports.py:38`). Default progress is separately compared against a baseline-source build (`docs/foundational/phase3b5/evidence/task2/compare_old_progress.py:35`).
- ⚠️ Later-task runtime correspondence, whole-branch qualification, manuscript reconciliation/PDF review, and final delivery are outside this task gate. Task 2 explicitly leaves them outstanding (`docs/foundational/phase3b5/03-authority-progress.md:81`).

### Strengths

- The checker validates the structure of a rejected grant independently before checking its recorded acceptance against current authority. Revocation cannot make a malformed grant appear acceptable merely because it was rejected (`research/authorityservice/oracle.go:236`).
- Mutation qualification requires successful compilation, exit 1 from the targeted test, the expected assertion text, and no panic. Positive c-root and crash-recovery controls deliberately detect an overrestrictive checker (`docs/foundational/phase3b5/evidence/task2/qualify_mutants.py:20`, `:60`; `research/authorityservice/revocation_oracle_test.go:9`, `:84`).
- The preserved temporal-cut RED result documents a real witness-extraction defect and the final code fixes extraction without changing the predicate or transition relation (`docs/foundational/phase3b5/03-authority-progress.md:71`; `research/delegationcorrespondence/progress.go:208`).
- Necessity/sufficiency are explicitly source-level sketches restricted to reachable finite states. Existential completion, retry SCCs, atomic effect assumptions, and unsupported runtime/distributed generalizations remain distinct (`docs/foundational/phase3b5/03-authority-progress.md:27`, `:49`, `:53`).

### Issues

- Critical: none found.
- Important: none found.
- Minor: none raised.

### Checks and Limits

- Read the supplied complete diff once in sequential sections; no git commands or test/build/lint suites were rerun. No new execution probe was needed.
- Focused context reads were required because the diff cut off `grantPath`, `accept`, and `auditGraph`: checked the named risks of inline/duplicate issuance admission, historical decision versus current permission, crash-permit lifetime, and omitted all-edge/merge checks in `research/authorityservice/oracle.go:162`, `:312`, `:568`. Read the cut-off graph construction/reverse-search portion of `QueryProgress` for the named risk of importing the operational explorer's answer (`research/delegationcorrespondence/progress.go:75`). No broader source crawl was performed.
- Read retained test, RED, mutation and lint streams. Final focused/progress/race streams contain passing package results and no warning, panic, failure or skip markers; Staticcheck's retained stream is empty. These are inspected implementer-run artifacts, not newly executed tests (`docs/foundational/phase3b5/evidence/task2/checks.json:1`).
- Decompressed the two new exports and independently hashed their bytes: audit `a5f9f8fac9faaceed3dc4af2c213565383b7954f75b271e7426807c20e0491ad`; progress `faf6ff19137e9c5ac9e22b97a4835699fc151323eb972fb5538c4a277a65db5a`. Parsed audit counts, 34 controls, all six temporal cuts, suffix lengths, paid-original flags, SCC counts and zero predicate mismatches agree with the note and recorded manifests (`docs/foundational/phase3b5/evidence/task2/exports.json:175`).
- Independently checked every current source against `source-sha256.json`, each archived original/decoded export hash against `exports.json`, and each retained mutation build/assertion stream against `mutants.json`: no mismatch. No binary was rebuilt or rerun during this review.
- Finite quotient edge consistency is not independent reconstruction of the entire operational graph or enumeration of every cyclic history. Passing bounded evidence and the source-level sketches do not establish unrestricted progress or runtime correspondence (`docs/foundational/phase3b5/03-authority-progress.md:7`, `:13`, `:59`).

### Assessment

**Task quality: Approved.** The implementation is narrow, preserves historical interfaces/exports, and tests the intended semantic distinctions with independent event interpretation and meaningful compiled mutants. No actionable defect was found within this task-scoped review.
