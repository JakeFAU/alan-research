# Residual interfaces

Finite epsilon-NFA residual-language profiles over authenticated Phase 3C.3
fixtures. The package exports only `Export(io.Writer) error`:

```
go run ./research/residualinterfaces/cmd/export
```

Run from the repository root. JSON contains complete language/profile certificates,
all 2,772 rebuilt products, contextual nonblocking rows, negative and private-loop
witnesses, a separate generic strictness control and complete source bindings.
No predecessor outcomes or historical captures are used as verdict oracles.

`automaton.go`, `language.go` and `product.go` retain the predecessor's small
finite-relation representation/traversal conventions, adapted locally without
importing predecessor packages. All product outcomes are recomputed from pinned
automata. This algorithmic similarity is not independent verification; the
separate Task 2 Python reconstruction provides that challenge.

See `docs/foundational/phase3c4/00-contract.md`, `01-argument.md` and
`02-evidence-schema.md` for assumptions and exact schema. Nonblocking is an
existential suffix property at every reachable state, not inevitable completion.
Theorem 4.8 and unrestricted ALAN substitution remain open.

The escapable tau cycle is an abstract monitor fixture. Its realization by the
runtime's strongly terminating local instruction grammar is unestablished; the
lasso does not show private divergence of validated ALAN local programs.
