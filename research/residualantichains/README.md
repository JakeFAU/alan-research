# Finite residual antichains (Phase 3C.5)

`Export(io.Writer) error` is the sole public boundary. From repository root:

```
go run ./research/residualantichains/cmd/export
```

This decodes authenticated Phase 3C.4 automata, recomputes all language inclusion,
profile and product results, and emits `alan/phase3c5/v1`. Old Go packages are not
imported and old verdicts are not oracles. The new original minimal fixtures
retain strict FULL inequality with J equality. Generic marked-domain counterexamples
stay outside the nine-monitor roster.

The contract, conditional argument and exact schema are in
`docs/foundational/phase3c5/00-contract.md`, `01-argument.md`, `02-evidence-schema.md`.
The tests include complete automata certificates, finite-set algebra, generic
finite trie controls and conditional product closure. Source arguments and bounded
checks do not establish runtime refinement, novelty, speedup or global necessity.

```
GOCACHE=/tmp/alan-phase3c5/go-cache go test -v -timeout 30m ./research/residualantichains/...
GOCACHE=/tmp/alan-phase3c5/go-cache go test -race -v -timeout 30m ./research/residualantichains/...
GOCACHE=/tmp/alan-phase3c5/go-cache STATICCHECK_CACHE=/tmp/alan-phase3c5/staticcheck-cache staticcheck ./research/residualantichains/...
python3 docs/foundational/phase3c5/tools/capture.py /tmp/fresh-phase3c5-capture
```

`ALAN_FALSE_HYPOTHESIS` deliberately enables failing research hypotheses in
`TestFalseHypotheses`; unset it for passing tests. Their nonzero archived results
are falsifications, not successful qualification runs. Capture directories are
immutable; overwrite attempts are rejected before commands or writes.
