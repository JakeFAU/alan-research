# Task 1 implementation report

Status: DONE, awaiting the controller's independent task review. No reviewer
approval is claimed. Task base: `4b8a0b8b135e314660a2146401f494aad4525cdb`.
Implementation commit: `5f907cda569d97ec5c08923b35e77de5b3df06ce`
(`Model named progress contracts under fixed composition`).
The resulting implementation commit is returned directly to the controller.
This report remains ignored workflow scratch for the controller to archive
with the task/review evidence; it is not part of the implementation commit.

## Question and result

Can the fixed chain/cycle preserve all-original completion and existential
completion under grouping while retaining named service dependencies and the
original retry-enabled durable-ACK transport semantics?

Implemented the four-stage progress contract and compatible leaf/join/close,
the concrete adapter and finite lifting, complete macro BFS/reverse reachability,
dependency certificates, seven genuine concrete histories, and deterministic
versioned export. The argument is a direct source-level conditional proposition,
not mechanized proof or complete concrete-product exploration. R0 and R1 remain
falsified candidate summaries. Abstract service eligibility is explicitly
separate from concrete immediate readiness.

Executed macro results: every chain mask has 16 states, 22 arcs when B is serial
or 24 when B services, and every state can complete. The all-serial cycle has
63 states, 117 arcs and exactly one noncompleting state (WWW). Other cycle masks
have 64 states, 128–144 arcs and every state can complete. Counts are generated
from actual graphs, not embedded expectations. The static cyclic certificate
is not treated as necessary for failure from a particular state.

## Implementation and assigned files

* `research/contractcomposition/model.go`, `compose.go`, `concrete.go`,
  `explore.go` and their tests: exact named bindings; validated config, state,
  action and component records; U/W/P/D progress; physical-owner transport
  adapter; old independent causal ledger checks; finite lifting from retained
  data; full small graph and deterministic evidence construction.
* `research/contractcomposition/cmd/export/main.go`, `main_test.go`, `README.md`:
  deterministic JSON CLI and package/API scope.
* `docs/foundational/phase3c1/01-argument.md`: every concrete transition class,
  invariant induction, universal representative lifting, both finite path
  directions, preserved Complete, allowed grouping grammar, sufficient
  dependency ordering and explicit limits.
* `docs/foundational/phase3c1/02-evidence-schema.md`: exact v1 fields, ordering,
  identities, graphs, open intermediate components, witness/projection records,
  retry certificate and claim categories for Task 2.
* `docs/foundational/phase3c1/evidence/task1/`: actual RED/GREEN and final logs,
  full `export.json`, source hashes/tool versions/argv/exits in
  `source-manifest.json`.
* This task report is ignored workflow scratch, to be archived by controller.

No old models, production/runtime source, dependency files, manuscript files,
contract/plan, retained SQLite or `.envrc` were modified or opened. No subagent
or reviewer was dispatched. Controller confirmed hostile JSON mutation records
and independent checker belong to Task 2; Task 1 supplies genuine records and
hostile Go tests. No extra helper source files beyond the assigned layout were
needed. Scratch runner/caches remain under `/tmp/alan-phase3c1/task1/`.

## TDD evidence and actual corrections

All log paths below are relative to `docs/foundational/phase3c1/evidence/task1/`.
The manifest contains argv arrays and exits; RED exits are 1, GREEN exits 0.

| Phase | Command (from module root) | RED finding / GREEN evidence |
| --- | --- | --- |
| Progress | `go test -count=1 ./research/contractcomposition/...` | `red-progress.log`: hidden debt/readiness lost, bad contracts accepted, missing graph/grouping; `green-progress.log` after implementation. |
| Concrete | `go test -count=1 -run 'TestConcrete\|TestPreRecord\|TestAbstractEligibility\|TestWitnessHostile\|TestEveryOld' ./research/contractcomposition` | `red-concrete.log`: missing transitions/lifts/witnesses; `green-concrete.log` uses the same filter plus `-v` after implementation. |
| Export | `go test -count=1 -run 'TestEvidence\|TestExport' ./research/contractcomposition/...` | `red-export.log`: missing domain/candidates and empty JSON; `green-export.log` after implementation. |
| Retained premises | `go test -count=1 -run TestLiftRejects ./research/contractcomposition` | `red-premises.log`: forged concrete binding and invalid activation were accepted by Lift; `green-premises.log` includes `TestLiftRejects\|TestCrashDisables\|TestOldLedger` after source validation. The RED command was repeated after changing the assertion to report both violations. |
| Oracle noninterference | `go test -count=1 -run TestLiftingDoesNotConsultWorldOracle ./research/contractcomposition` | `red-lifting-oracle.log`: validation's oracle check blocked lifting when the inaccessible oracle was erased. Fixed validation to inspect only actual retained records; `green-lifting-oracle.log` passes. |
| Constructors | `go test -count=1 -run TestStateAndActionConstructorsRejectUndeclaredRecords ./research/contractcomposition` | `red-constructors.log`: accepted absent CA/invalid stage/unsupported action and omitted binding; `green-constructors.log` after validated constructors. These commands used the task-local GOCACHE. |

The oracle failure was substantive: merely avoiding oracle reads in the raw
successor relation was insufficient when the lifting helper's validator read
that oracle. Both actor stepping and lifting selection now ignore WorldResult;
receipt authenticity comes from the actual trace ledger/custody premises.

The original factored run checked AB in a chain: 7,464 contexts, 49,992 arcs,
4,320 lifts, retained in `green-concrete.log`. After controller feedback, coverage
expanded to all AB/BC/CA cycle renamings and terminal chain BC; that newer run is
retained separately in `green-factored-expanded.log`. Earlier results were not
silently relabeled as final coverage.

## Final checks actually executed

Final verification used `GOCACHE=/tmp/alan-phase3c1/task1/go-cache` and
`STATICCHECK_CACHE=/tmp/alan-phase3c1/task1/staticcheck-cache` where applicable.
Tools: `go version go1.27.4 darwin/arm64`; `staticcheck 2026.2.2 (0.8.2)`.

* `go test -count=1 -v ./research/contractcomposition/...`: exit 0;
  `normal.log`; 19 top-level tests plus 7 named config subtests pass.
* `go test -race -count=1 ./research/contractcomposition/...`: exit 0;
  `race.log`; both packages pass.
* `staticcheck ./research/contractcomposition/...`: exit 0; empty clean
  `staticcheck.log`.
* `gofmt -l` over all ten new Go files: exit 0; empty `formatting.log`.
* `go run ./research/contractcomposition/cmd/export`: exit 0 under task-local
  GOCACHE. A second final execution was byte-equal to the earlier retained
  1,584,501-byte export. Export SHA-256:
  `48b4cfb4cc10b7dde4979874f669014f4863e1e05e48737ef4561f119dfc3602`.
* All seven genuine histories pass legal concrete adjacency and each original's
  separately implemented old `twoagent.CheckTrace` ledger. Hostile tests reject
  forged event/state bindings, unsupported records, premature release and
  duplicated application. Physical B crash disables both roles while retaining
  durable custody. The exact retry loop is retained, not pruned from semantics.
* Full source-input manifest covers 17 inputs: ten new Go files, package README,
  contract/argument/schema, and unchanged twoagent model/explore/oracle files.
  The three old source hashes were compared to `git show <task-base>:<path>` and
  match. The entire exported object was traversed for unexpected null records;
  only declared initial-node Via nulls occur.

Actual factored coverage: all 311 old durable-ACK reachable edge states, crossed
for each named cycle edge with four callee outgoing stages, two policies and
three physical phases; plus terminal chain BC with no outgoing C original,
both C policies and three phases. Result: **24,258 representatives, 184,256
micro arcs projected and 14,058 enabled named-edge lifts**. The surrounding
edges use canonical durable representatives. These numbers include overlapping
contexts and do not enumerate reachable full-product microstates.

Normal/race passed before the final constructor addition as well; those logs
remain `normal-before-constructors.log` / `race-before-constructors.log`. The
final required checks were rerun after that source change. No duplicate root
suite was run. Controller reports the frozen 70-package baseline completed
exit 0 separately; that claim is attributed to controller evidence, not this
worker's execution. Root build/module/global Staticcheck accounting remains
controller/Task 3 work.

## Tool failures and self-review

The first `go run` hit a sandbox denial in the default Go cache and produced an
empty redirected sample; a following parser failed on it. No successful export
was claimed. Switching GOCACHE to the authorized task scratch directory produced
the retained genuine sample. Staticcheck initially failed exit 1 trying to write
its default cache; the actual diagnostic remains
`staticcheck-cache-failure.log`. It passed with the task-local STATICCHECK_CACHE.
No sandbox escalation or rule change was needed. Read-only discovery of an
as-yet absent Task 2 brief and absent Makefile/CI directory also returned missing
path diagnostics; no build/check was inferred from those paths.

Self-review covered the complete authored source/test/document diff and the
full staged path list. It rechecked the serial-only Apply filter, all physical
roles, exact bindings through old-edge renaming, caller-independent callee
steps, original retry/duplicate paths, durable ACK release, no reapplication,
projection labels, all-original debt, and the universal lifting cases. The
large generated JSON was reviewed through deterministic regeneration, full
record traversal, graph count/arc checks and genuine-witness replay; the
independent semantic audit remains Task 2. Self-review also clarified shared incoming/outgoing roles (without implying
shared component ownership) and described flat/left/right runs as separately
explored, avoiding an independent-oracle claim. Source-manifest input hashes and
retained output names were checked. Initial staging was blocked by the
read-only worktree Git index; the narrowly scoped Git escalation succeeded.
The report was excluded from the index per controller instruction. Only
assigned new source/docs/evidence paths are included;
logs are explicitly retained despite the repository's general `*.log` ignore.

## Remaining obligations and boundaries

The conditional source argument must receive independent review. Task 2 must
check exported graphs and histories without relying on producer successors,
and retain its hostile JSON controls/rejection reasons. Task 3 owns manuscript
synchronization and overall qualification/delivery. These are pending tasks,
not completed evidence here.

No whole concrete-product enumeration, mechanized universal theorem, arbitrary
context/topology result, inevitable progress/fairness/divergence equivalence,
payload equivalence, runtime/distributed refinement, shared-authority fan-out,
changed identity/grant support or reduced verification-cost result is claimed.
R1 only refutes the declared readiness summary. Permanent owner loss, split
application/retention, erased custody, spoofed binding and resource depletion
violate the fixed proposition's premises. Phase 3C as a whole remains open.

Final commit check: `git diff --cached --check` exited 0; all 17 staged source
hashes and all recorded command output paths verified; 35 assigned paths only.
Post-commit `git status --short` is clean (ignored workflow report excluded).
