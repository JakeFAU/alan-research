# Phase 3B.2 reproduction

Run from the repository root with Go 1.27.3 and the existing cached dependencies. There are no new dependencies, live model calls or paid endpoints. Integration tests need local loopback listeners and subprocesses. Read `evidence/source.json` for the final immutable executable revision and `evidence/commands.json` for actual statuses and tool versions; development logs alone are not the final qualification.

```sh
export GOTOOLCHAIN=local
export GOCACHE=/tmp/alan-phase3b2-reproduction-cache
export STATICCHECK_CACHE=/tmp/alan-phase3b2-reproduction-staticcheck
export GOPROXY=off
export GOSUMDB=off

go test -count=1 -timeout=60m ./...
go test -race -count=1 -v ./research/twoagent/... ./research/mutualobligations/... ./research/authorityservice/...
go test -count=1 ./examples/communication/authority
go test -tags=integration -count=1 -v ./examples/communication/mutual
go test -race -tags=integration -count=1 -v ./examples/communication/mutual
ALAN_PHASE2_EVIDENCE=/tmp/alan-phase3b2-runtime-normal \
  go test -tags=integration -count=1 -v ./internal/durabletest -run '^TestPhase3B2'
ALAN_PHASE2_EVIDENCE=/tmp/alan-phase3b2-runtime-race \
  go test -race -tags=integration -count=1 -v ./internal/durabletest -run '^TestPhase3B2'
staticcheck -tags=integration ./research/twoagent/... ./research/mutualobligations/... ./research/authorityservice/... ./examples/communication/authority ./examples/communication/mutual ./internal/durabletest
```

The integration-tagged Staticcheck command is expected to report the six unchanged older-test diagnostics listed in the results; the final run checked model/fixtures and durable tests separately and matched durable-test diagnostics byte-for-byte against the merged base. It is not a green lint gate.

The runtime child helper returns immediately when not invoked as a child; it is exercised by the real process matrix. A successful test run must include the individual crash cuts, not merely the helper name. The ordinary root command excludes integration-tagged tests and separate consumer modules. Focused integration/race commands qualify the new witness and unchanged reciprocal fixture; this is not a claim to have run every integration test in the repository.

## Deterministic operational export

```sh
go build -buildvcs=false -o /tmp/alan-phase3b2-export ./research/authorityservice/cmd/export
/tmp/alan-phase3b2-export > /tmp/alan-phase3b2-model.json
gzip -dc docs/foundational/phase3b2/evidence/model.json.gz > /tmp/alan-phase3b2-archived-model.json
cmp /tmp/alan-phase3b2-model.json /tmp/alan-phase3b2-archived-model.json
```

The export retains all sixteen graph counts, shortest witnesses, temporal cuts and cyclic witnesses without serializing every state. BFS has no depth cutoff. Reverse reachability distinguishes raw and authority-safe completion. Repeated output must be identical for the frozen source. State/arc counts do not count the infinitely many cyclic histories.

Expected-failure controls:

```sh
/tmp/alan-phase3b2-export -config S3/Missing -assert-authority-safe
/tmp/alan-phase3b2-export -config Cartesian/ScopeMix -assert-authority-safe
/tmp/alan-phase3b2-export -config S2/Missing -assert-universal-completion
/tmp/alan-phase3b2-export -config S2/LateGrant -assert-universal-completion
/tmp/alan-phase3b2-export -config S2/Revocable -assert-universal-completion
```

Each must exit **1**, preserving JSON on stdout and explaining the falsified candidate on stderr. Exit **2** is an argument/output/apparatus error, not the intended counterexample. A positive S2/Explicit run with both assertion flags must exit 0. The shortest unsafe authorization/admission/application prefixes are distinct four/five/six-event witnesses.

## Independent retained ledger

The independent checker reconstructs roots, grant ancestry, exact service identities, decisions, admission opportunities and custody from events. Its graph check validates one BFS-prefix ledger per state and every outgoing edge, with merged-target ledger comparison. This checks the finite abstraction's edges; it does not enumerate every cyclic history. Retained malformed histories preserve the visible application/result tail while changing decisive authority edges. Source-defined controls and final exported audit coverage are recorded in the evidence directory.

The standalone independent export is reproducible as follows:

```sh
go build -buildvcs=false -o /tmp/alan-phase3b2-audit ./research/authorityservice/cmd/audit
/tmp/alan-phase3b2-audit > /tmp/alan-phase3b2-audit.json
gzip -dc docs/foundational/phase3b2/evidence/audit.json.gz > /tmp/alan-phase3b2-archived-audit.json
cmp /tmp/alan-phase3b2-audit.json /tmp/alan-phase3b2-archived-audit.json
```

## Runtime evidence and boundary labels

Archive the full verbose (`-v`) command output as well as `ALAN_PHASE2_EVIDENCE`. The existing checked kill helper writes temporary `evidence-*.json` files and logs their complete signal/disposition records as `DURABLE` JSON; it does not implement an `ALAN_EVIDENCE` environment switch. The final archive extracts those records from the verbose logs. Each exported runtime case includes ordinary provisioned programs, permissions, original messages, exact inspections, selected history, coherent compressed SQLite backups and independent endpoint ingress/application evidence. Crash cases additionally retain the observed barrier, actual SIGKILL/fresh PID, cut and recovered stores, and paid/unpaid counts. Reopening may append facts but must preserve original decisions, results, identities and permissions. Endpoint addresses, store namespaces and process identities vary; byte-identical runtime output is not claimed.

`intent.authorized` records source disposition. `dispatch.consider/pre-admission` follows an allowed adapter reservation but precedes physical admission. Despite its name, `adapter.authorized/pre-transport` follows physical registration and precedes the adapter call. `adapter.result-recorded/pre-settlement` and `adapter.settled/pre-consumption` separate receipt retention, settlement and authored consumption. Generic allowed reservations without a result remain unknown after reopening. The test's knowledge that a loopback endpoint saw no contact does not become a runtime receipt or grant.

The final normal/race raw directories are preserved as `runtime-normal.tar.gz` and `runtime-race.tar.gz`, with per-member digests in `runtime-index.json`. Each archive preserves its original verbose log and command provenance. `kill-records-normal.json` and `kill-records-race.json` contain the actual extracted checked-signal records. The post-run archive utility is identified by its own delivery artifact hash, separately from the frozen Go execution source. Recheck the archives with:

```sh
python3 docs/foundational/phase3b2/evidence/runtime-audit.py \
  docs/foundational/phase3b2/evidence /tmp/alan-phase3b2-runtime-audit.json
cmp docs/foundational/phase3b2/evidence/runtime-audit.json /tmp/alan-phase3b2-runtime-audit.json
```

No general S2 grant exists in the inspected public API, so the grant-retained runtime cut is explicitly inapplicable. Static permission persistence and receipt-adoption grants are not substitutes. The model's S2 and revocation results are qualified independently; there is no claimed full model/runtime refinement.

## Manuscript and artifact identity

From `paper/`, run `tectonic --keep-logs alan-dissertation.tex`, then render the affected PDF pages with Poppler and inspect them. Markdown and LaTeX must agree on the new propositions, counts, assumptions and limits. The new BibTeX attribution is Experimental RFC 2693, verified from the primary source, not a new ALAN attenuation principle. The edited Abstract and Conclusion are also format-normalized against their counterparts. Historical manuscript sections outside this slice are not asserted to have been normalized. The PDF has 37 pages; pages 1–3 and 34–37 were visually inspected. Existing warnings at TeX lines 215 and 559 remain, including the pre-existing 5.2184pt overflow at 559; no new warning or visible defect was found in the edited authority section.

Final evidence records exact source hashes, artifact byte lengths/digests, commands/statuses, skipped cases and independent review dispositions. No production source or Phase 3A/3B.1 model/fixture/evidence file is changed by this slice.
