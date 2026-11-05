## Completion follow-up: genuine M03 target-generation-two parent

The independent review `alan-phase3b4-campaign-review-r1.md` **APPROVED the bounded Task4 campaign** at runtime `2b8f6e19a013d2b918146a7193451843004d2f86` / research `5d7473735d550ac1550228e9991863c3601afde7`. It identified a missing M03 producer dependency and a P3 overbroad supplied-parent statement; it did not reject the campaign. Its inspection found 942 message appearances in the selected normal/race captures, all targeting generation 1. False claim-generation controls did not supply the required genuine target 2 message. The controller authorized this narrow completion before Task5, with scoped addition review to follow.

Runtime completion commit `fc60391f48b46a8c8635afc46882db8338b6c93b` adds only `internal/durabletest/phase3b4_generation_parent_test.go`. `TestPhase3B4CampaignAuxiliaryTargetGeneration` genuinely authors q2 `["b",1,6,0]` with target `a`, generation 2 before Create. The existing P10 insertion keeps planned service/reply locators coherent: principal BA serviceA5/replyA8 and AB serviceB10/replyB13. Public Commit emits q2 at B6/PC5. Public AuthorizeIntent records `Allowed=true`, reason `authorized`, source B1. Public DeliverMessage retains `Accepted=false`, reason `target_generation_changed`, `Bound=A1` and an empty observation ID. Bound identifies the current run considered by the refusal; it is not an accepted A1 delivery. The actual q2 message has the same obligation/payload as q but a distinct full ID and requested generation.

Both final captures contain exactly A1/B1, unchanged actual root records, unchanged subject generations, no q2 observation or accepted admission, no principal service emission and no endpoint contact/application. Both runs remain open and unpaid. This deliberately incomplete public prefix is outside C and demonstrates no live generation transition, principal completion, or mutant kill. No A2 run/root, fabricated target 2 intent, retained repair, direct SQL, private consumption or production policy change is introduced.

Separate evidence is `docs/foundational/phase3b4/evidence/task4-completion-m03/`. Its `case-index.json` names one new case per mode and ordered `emitted-before-authorization` → `final` captures. Actual program/intent/admission-descriptor/result/source bytes, raw and typed inspection, Replay, paged history, endpoints and sidecars are archived. `runtime-source.json` pins all 328 Go/module files; every original 327 source-file hash still matches the approved 2b source. `approved-package-before-completion.json` verifies every original Task4 evidence file against research 5d before this narrative/manifest correction. Original 80-case indices, raw archives, costs, correspondence and source archives remain byte-identical. Only the original report, its manifest hash and README handoff are updated. The exact APPROVE review and its check.py/check.json/extra.json/m03.json artifacts are copied without alteration.

| Mode | Actual case path | Case SHA256 |
|---|---|---|
| Normal | `normal/M03-auxiliary-target-generation-two-1165515867` | `62e3e9f602551984e3d7e630b1a851399a0c1d3a457f4ffd4d78a0b2c0b4e1ad` |
| Race | `race/M03-auxiliary-target-generation-two-341226844` | `0dca65ce497732d5eb35f7af206c1ed1ebecd93037b7983000e376714a5281e4` |

Focused commands ran from `/private/tmp/alan-phase3b3-runtime` using the previously documented offline environment. The wrapper `/tmp/alan-task4-run.py` retains exact argv/environment, full source hashes, exit and logs in the completion package's `commands/`:

```sh
ALAN_PHASE3B4_EVIDENCE=/tmp/alan-task4-m03-final-normal python3 /tmp/alan-task4-run.py m03-final-normal go test -count=1 -timeout=3m -tags integration ./internal/durabletest -run '^TestPhase3B4CampaignAuxiliaryTargetGeneration$'
ALAN_PHASE3B4_EVIDENCE=/tmp/alan-task4-m03-final-race python3 /tmp/alan-task4-run.py m03-final-race go test -count=1 -timeout=3m -race -tags integration ./internal/durabletest -run '^TestPhase3B4CampaignAuxiliaryTargetGeneration$'
python3 /tmp/alan-task4-run.py m03-staticcheck staticcheck -tags integration ./internal/durabletest/...
gofmt -l internal/durabletest/phase3b4_generation_parent_test.go
git diff --cached --check
```

Normal test output: `ok github.com/JakeFAU/alan/internal/durabletest 2.097s`; race: `ok github.com/JakeFAU/alan/internal/durabletest 3.800s`. Staticcheck, format and staged whitespace checks succeeded with empty output. All three final command source maps match the completion commit exactly. The unchanged full campaign was not repeated; the new name matches its future `^TestPhase3B4Campaign` selector. Task5 still owes final frozen-pair qualification.

The first normal probe obtained the genuine target-generation refusal but then failed an apparatus assertion comparing normalized retained roots (assigned IDs/version and canonical ordering) to unnormalized provision roots. It did not expose a runtime failure. `apparatus-red.tar.gz`, `commands/m03-normal.*` and `alan-task4-m03-root-order-red.go` preserve the actual failure/raw prefix/exact source. The corrected assertion compares exact retained root records before and after delivery and separately checks the original subjects and only A1/B1. Both focused modes passed afterward. The failed run's source is reconstructible from approved 2b plus the preserved new file and checked against its command source map.

The genuine M03 producer parent is now supplied. Task5 must still construct the coherent contradictory accepted-A1 admission/observation/index mirrors, demonstrate sound-checker rejection and generation-join-mutant acceptance, and retain its own outcomes. No Task5 attack or mutant result is claimed here. The prior N10e/f construction boundaries, ancillary internal-writer-only promoted14/16 compatibility scope, and Task6 manuscript obligations remain unchanged.
