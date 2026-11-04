# Contract-v3 scoped confirmation

**Verdict: APPROVE** for the contract text at `.superpowers/sdd/03-plan/contract-v3-draft.md` with SHA-256:

`c0b1aa53af1331f31d2445d86717a1702293d98272df4eab6eec464316bc9fff`

This is a document-only confirmation of the two requested corrections following the substantive contract-delta assessment in `task-2-semantics-review.md`.

1. Section 4 now names only the existing `DispatchExternal` API, preserves its complete-composition `ErrDispatchMode` refusal, and requires every supported physical route to pass the gate. The nonexistent `DispatchAdapters` reference is removed. No new direct transport route is authorized.
2. The v3 correction addendum now attributes computed DelegateAuthority to the standalone af6943a public probe, and literal/computed coverage for both DelegateAuthority and RevokeAuthority to the separate four-case cross-version regression. This matches the previously inspected source and retained execution evidence.

I inspected the complete v2-to-current-draft diff and computed the above SHA with `shasum -a 256`. The substantive distinctions remain intact: legacy inert payload strings preserve their old retention and unknown-operation behavior, while actual authority metadata, accepted authority decisions, and gate facts remain feature-bound.

No tests were executed, no implementation source was edited or re-reviewed, and no subagents were spawned for this confirmation. The candidate specification and code-quality **REJECT** verdicts in the earlier review remain open; this approval neither resolves those findings nor approves Task 2 implementation or the whole phase.
