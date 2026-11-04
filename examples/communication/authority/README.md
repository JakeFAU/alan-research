# Phase 3B.2 ordinary same-Store authority witness

`New` builds exactly two ordinary retained roots, A and B, without invocation or
special runtime agent roles. Both commit reciprocal original requests. In the
local case A serves before awaiting its own result and B awaits before serving.
Only actual consumption of both authentic service results completes the originals;
transport acknowledgments, refusal observations and runner quiescence do not.

The trusted deployment registry fixes two distinct endpoint configurations:
`phase3b2.counter-A.red` means increment/counter-A/red, and
`phase3b2.counter-B.blue` means increment/counter-B/blue. Their exact configuration
is retained in `Prepared.Endpoint` and `Prepared.Model`. These names encode whole
finite tuples for this fixture. The runtime does not independently interpret
resource/scope labels from arbitrary payloads or enforce a general tuple algebra.
A retains only the first static adapter permission and B only the second. When A
attempts B's operation, construction removes the compiler's automatically supplied
permission before public `Create`, while preserving A's original local permission.
No positive SQLite records are rewritten.

The missing case retains an actual `grant_denied` adapter attempt and unsuccessful
result. A grant-shaped claim in a request remains data. The result laundering case
changes service placement so A first consumes an authentic result from B, then
includes that result's mutation value in its forbidden operation. Possessing that
result confers no adapter permission. Neither refused original is paid by its
refusal. In the result laundering case one original is genuinely paid and the
other remains unpaid.

The public runtime has no retained, operation-bound, rooted delegator/attenuation
chain primitive for the model's S2 contract. Static permissions and the specialized
receipt-query grants do not realize it. S2 positive construction and delegation
retention interruption cuts are therefore inapplicable here. This is an explicit
implementation obstruction, not a test of a substitute grant protocol.

`internal/durabletest/phase3b2_runtime_test.go` uses public execution and existing
observation/process helpers. It kills real executors with SIGKILL, then starts fresh
owner processes. Its nine cuts distinguish before source authorization, after
source authorization, allowed adapter reservation before physical admission,
after physical admission before transport, result before settlement, settlement
before consumption, genuine grant refusal, mutual wait, and one original paid
with one pending. `adapter.authorized/pre-transport` occurs after physical
admission; `dispatch.consider/pre-admission` occurs after durable reservation.

An allowed generic adapter reservation without a retained result never retransmits
after reopen, even when killed before any endpoint contact. The original remains
unknown and unpaid. This differs from the model's atomic application-plus-receipt
premise and is deliberately preserved. Known results settle/consume without
physical reissue. Same-Store ownership and fixed local endpoints establish no
independent-owner, distributed, arbitrary-context or general delegation theorem.

Set `ALAN_PHASE2_EVIDENCE` to a fresh temporary directory to export the reused
process helpers' JSON records, coherent SQLite backups, endpoint ingress journals,
raw program JSON, and explicit Phase 3B.2 endpoint/config environment records.
The latter supplement the historical helper's Phase 2 environment whitelist.

Archive the verbose `-v` run log too: the existing kill helper writes its checked
signal/PID record as `DURABLE` JSON in that log. It does not implement an
`ALAN_EVIDENCE` switch. The Phase 3B.2 runbook preserves those records alongside
the raw phase artifacts.
