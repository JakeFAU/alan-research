# Phase 3B.3 contract-v2 logical root locator review

**Verdict: APPROVE the proposed contract-v2 representation correction.** This is design approval only, not approval of an implementation, completed qualification, or a PR. No production source was edited and no tests were run by this reviewer.

## Evidence inspected

- Frozen research `docs/foundational/phase3b3/02-contract.md` (version 1).
- Research `.superpowers/sdd/03-plan/task-2-report.md`.
- Runtime `internal/executor/request_authority_test.go` and `/tmp/alan-phase3b3-literal-root-red.log`. The retained prior run reports one intended failure, `root identity mismatch`; this review did not rerun it.
- Actual runtime `composition_types.go:137–139`, `composition_store.go:93–175`, `store.go:328` and the root/profile validation paths at base `af6943a5805a9a6c0e8fe77f804264f4765d07ec`.

The original design correctly removed a literal Namespace requirement but missed the same dependency through literal RootID. Root IDs are `compositionRootID(namespace, ResourceRoot.Run)`; Create chooses namespace internally; provision rejects a nonempty ID that does not equal the derived identity. The failed second Create is a concrete falsifier of reusing that literal ID under ordinary construction, consistent with the source. It is not a general impossibility theorem about delegation or every possible program. Preserve the rejected version-1 wording, test source and failing log before replacing its assertion or fixture.

## Why the correction is representable

Use these exact logical locators in program-authored claims and selectors:

```go
type AuthorityRoot struct {
    Run        RunID      // permission holder; not necessarily resource root owner
    RootRun    RunID      // exact ResourceRoot.Run
    Permission Permission
}
// DelegationClaim.RecipientRoot RunID replaces RecipientRootID string.
```

Both RunIDs and their generations are known when the ordinary provision and raw program artifacts are constructed. The Store already derives the opaque identity from namespace and the resource-root RunID. There is no need to choose the namespace, rewrite artifacts after Create, add a runtime interpolation language, use an empty wildcard, or broaden authority.

The actual composition provision enforces uniqueness of `ResourceRoot.Run`, generation one, each static member's unique membership, and inclusion of the root run among its members. Thus `RootRun` identifies at most one retained resource root in the supported profile. `Root.Run` and `Root.RootRun` must remain distinct fields: a permission holder may be a member of a root whose Run is another principal. Equality happens in the simplest separate-root fixture but must not be a hidden semantic assumption. Same-root two-principal provision remains expressible, with holder-specific permissions and revocations.

## Exact safety interpretation approved

For a claimed root locator L in Store namespace N, computing `compositionRootID(N,L)` yields a **candidate identity**, not a grant, proof of membership, or evidence that a row exists. Acceptance additionally requires all of the following:

1. L is nonempty, well-formed, exact generation one for this profile, and not a wildcard or implicit current-root selector.
2. The derived ID resolves to a retained root whose descriptor ID, descriptor Run, profile entry and relational row agree; its Run is exactly L.
3. The holder's exact RunID is an actual static member of that root, with the retained membership index and root member list agreeing. Resolving the holder's actual root and silently replacing a wrong claimed locator is forbidden.
4. Holder source identity/revision, original complete Permission, current generation/retirement/revocation validity, issuer consent, recipient binding and original request/service checks all remain those of the frozen contract. Neither the root owner nor any other member's permission may substitute for the actual holder's permission.
5. The recipient locator resolves analogously to the exact recipient's actual static root. A root that exists but does not contain that recipient fails, just as a nonexistent candidate root fails.

Apply the same checks at grant issuance, first effect authority decision, and each fresh physical gate, with historical interpretation at the corresponding retained prefixes for replay/inspection. The gate must recheck the **selected original locator/witness**; it must not reinterpret the selector as whatever root currently happens to contain a similarly named actor. The composition profile is immutable here, so this adds no dynamic root migration protocol.

Revocation keys continue to include exact holder RunID, exact RootRun, and complete Permission. RootRun alone must not revoke all permissions of a resource root or all of its members. An altered RootRun is a wrong scope, not a route to borrowing another principal's tuple. The unaffected-r control remains mandatory.

## Retained derived identities and denial behavior

The proposed retained additions are safe and useful:

- `AuthorityDecision.RootID`: derived from the claim's Root.RootRun and actual Store namespace, even if the claim is denied because the root does not exist or the holder is not a member.
- `AuthorityDecision.RecipientRootID`: derived from Delegate's claimed RecipientRoot; absent/empty only for Revoke, whose claim has no recipient.
- `EffectAuthority.RootID`: derived from Use.Root.RootRun for every nonnil well-formed Use; empty only for nil Use and the explicit missing-selector denial.
- `EffectAuthority.ServingRootID`: the exact root of the actual committed effect's serving RunID. This field must be derived from actual immutable serving membership, not inferred from a possibly false grant's RecipientRoot. In a valid delegation it necessarily equals the accepted grant's derived RecipientRootID. In a denial it remains the actual serving root and can expose the mismatch. The opt-in profile guarantees serving membership; missing/corrupt actual membership is store corruption or invalid provision, not an empty successful witness.

This last interpretation is the recommended unambiguous definition of ServingRootID. It avoids making the effect's real accounting scope depend on a hostile grant claim and preserves the existing serving-root debit contract. It does not require another field or subsystem.

A candidate ID on a denied witness must not be inserted into an outer composition fact's RootIDs as if it were a real participant. The outer gate RootIDs contain actual participant membership roots only. For a gate following an allowed selected attempt, those roots must coincide with the already validated selected issuer and server roots; inconsistency is corruption, never a repair. Rejected issuance/first decisions can retain candidate IDs within their typed denial descriptors while the enclosing fact remains rooted in the actual committed issuer/target participants under the existing lifecycle/reservation rules.

Every persisted derived ID is checked by recomputing it using retained namespace and locator, then separately validating the existence/membership predicate appropriate to the disposition. Do not check a denied witness by requiring the candidate to exist, and do not check an accepted witness by mere hash equality. Optional omitted IDs must follow the exact rules above, not a general `omitempty means unresolved` convention.

## Falsifiers and qualification needed

Before claiming the correction realized, require:

- The same namespace-free literal program/provision can create two distinct Stores successfully; each derives different namespace-bound actual root IDs while preserving the same intended logical locator and member bindings. No artifacts are patched after Create.
- Correct RootRun with different holder Run (same-root two-principal case) works only with the holder's own legitimate tuple; another member's tuple remains unauthorized.
- Wrong root owner RunID, nonexistent root owner, wrong root generation, root that exists but lacks holder membership, and correct holder with the other principal's root each yield the specified refusal.
- Equivalent recipient-root mutations refuse. Same-root issuer and recipient are permitted only when all other one-hop/request/tuple checks hold.
- Modifying derived RootID/RecipientRootID/ServingRootID while preserving locators, or modifying locators while preserving old IDs, fails independent retained validation.
- Copying accepted decisions/attempts from Store X to Store Y cannot pass merely because both have equal principal names and logical locators. Namespace-bound IDs, exact retained originals and canonical descriptors must agree with the owning Store. Reusing the same program as a new legitimate computation is different from transplanting old retained authority evidence.
- Denied nonexistent-root claims retain their deterministic candidate IDs without fabricating root rows, membership, grants or outer RootIDs. Reopen preserves that refusal and never resolves it into an allow.
- Root revocation still leaves another tuple of the same holder and another holder's tuple intact; no revocation key collapses to RootRun alone.

The positive two-Store construction proves only that the literal locator obstruction is removed. It does not establish the authority chain, permission attenuation, revocation/reopen semantics, or progress; the remaining Task2 tests and independent reviews still apply.

## Version and scope judgment

Changing the written contract to version 2 while retaining the unused protocol string `alan/request-authority/v1` is coherent **on the reported preproduction facts**: no authority implementation or persisted version-13/14 artifact has yet instantiated the old wire contract. Preserve this as a rejected design and correction, not a migrated deployed protocol. If any such old-format artifact is discovered, the assumption must be revisited and it must not be silently read under changed field semantics.

No further simplification is needed. Removing explicit root locators and always deriving from holder membership could be safe in a narrower immutable profile, but would eliminate an explicit selected scope and weaken the wrong-root falsifier. Keeping exact namespace-free RootRun/RecipientRoot locators with derived retained IDs is a small, auditable correction that answers the reproduced obstruction without weakening the authority boundary.

## Final draft comparison

I subsequently read the updated contract-v2 Section 2 and correction addendum in `docs/foundational/phase3b3/02-contract.md`, SHA-256 `c96241fa365807c4322756d0513c9fda05edff6650c4bc09f32f8cde1bbcb863`. It explicitly defines ServingRootID from actual serving membership, keeps candidate IDs separate from membership proof, and restricts outer gate RootIDs to actual membership roots. These match the approved interpretation above. I also verified the rejected contract/report archive paths and runtime RED source/log archive paths exist. **Final verdict remains APPROVE for this contract-v2 correction; no outstanding blocking design finding in the reviewed delta.** Implementation and executable qualification remain outstanding.
