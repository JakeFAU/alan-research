# Task 2 implementation status: frozen literal-root obstruction

Status: **not implemented**. A concrete representability defect was found in the
frozen contract before production edits. No authority validation was weakened.

## Question and result

Can precompiled literal `AuthorityUse` and `DelegationClaim` values identify the
exact existing composition resource roots, without receiving a Store namespace
before ordinary `Create`?

The frozen construction does not currently provide a route. `AuthorityRoot.RootID`
and `DelegationClaim.RecipientRootID` name actual existing resource-root IDs.
`compositionRootID(namespace, run)` hashes the randomly generated namespace and
root run (`internal/executor/composition_types.go:137`). `Create` generates that
namespace internally after the caller has constructed the provision and artifacts
(`internal/executor/store.go:328`). Supplying a predeclared `ResourceRoot.ID` does
not select a stable root name: retention rejects it unless it matches the newly
generated namespace-bound ID (`internal/executor/composition_store.go`,
`retainCompositionProfile`).

The contract deliberately avoids a namespace field in program claims because
literal programs cannot know it. The required root IDs retain that dependency
indirectly. Dynamic payload construction alone would not solve the defect:
`Effect.Authority` is frozen as a literal pointer copied unchanged by emit, not an
operand or a runtime lookup.

This is a construction obstruction in this contract, not a proof that the
request-bound primitive needs a larger engine or is impossible.

## Executed RED evidence

Added only
`/private/tmp/alan-phase3b3-runtime/internal/executor/request_authority_test.go`.
`TestRequestAuthorityLiteralRootLocator` first creates a real ordinary public
composition Store, obtains its canonical root ID, then attempts a second ordinary
Create with the same provision and that exact predeclared ID. The intended
literal-root construction fails with `root identity mismatch`.

Executed from `/private/tmp/alan-phase3b3-runtime`:

```sh
GOTOOLCHAIN=local GOCACHE=/tmp/alan-phase3a-go-build GOPROXY=off GOSUMDB=off go test ./internal/executor -run '^TestRequestAuthorityLiteralRootLocator$' -count=1 -v > /tmp/alan-phase3b3-literal-root-red.log 2>&1
```

Exit status: **1**. One test executed, one failed at the intended assertion.
The first Create succeeded. Failure was a semantic root-identity mismatch, not a
compiler failure or missing environment dependency. The exact log is retained at
`/tmp/alan-phase3b3-literal-root-red.log`.

No green implementation, race, Staticcheck, legacy-format qualification, positive
two-agent Run, replay/inspection authority validation, or physical-gate
qualification is claimed. The three existing baseline files are unchanged.
No production source, dependencies, archives, branches, or commits were changed.

## Required contract decision

Parent was notified with the failing command and source evidence before any
representation change. A narrowly scoped option is a canonical logical root
locator in program claims/selectors, resolved to the namespace-bound existing
root at decision time and retained explicitly in decision/gate evidence. This
requires a documented contract correction and validation of wrong-root claims;
it must not silently interpret empty or arbitrary root strings as authorization.
Other options such as caller-controlled namespaces or post-Create artifact
rewriting have different provenance consequences and are not implemented.

The new test deliberately remains RED pending the parent's contract decision;
retain its original source and log as a rejected construction witness before
changing the assertion or fixture. All remaining Task 2 implementation and
qualification work remains outstanding.
