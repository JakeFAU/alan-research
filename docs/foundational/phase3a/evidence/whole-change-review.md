# Phase 3A independent whole-change review

Reviewed the frozen initial Phase 3A slice against `AGENTS.md`, the charter,
contract, plan ledger and `whole-change-diff.patch`, with direct inspection of
the model, oracle, tests, results/runbook, manuscript additions, evidence records
and relevant Store admission/commit code. No production or research source was
modified. This review does not close the wider Phase 3 program.

## Verdict

**Research/source verdict: READY. No Critical or Important implementation,
argument, manuscript-consistency or evidence-claim defect found.**

**Final verification/finishing verdict: NOT YET READY at this review snapshot.**
`evidence/commands.json` correctly records the original repository suite as
running. Its final result and package-coverage reconciliation are required before
claiming completed repository verification. This is a pending verification gate,
not a discovered implementation defect. Preserve the running suite; there is no
reason to start another full run blindly.

Concrete completion: retain the terminal output and exit status, verify that
its package outputs plus the independently tested new model packages cover the
final root-module package set, update the command record and applicable plan
checkboxes, and refresh affected source hashes. If a test fails, investigate and
report it; do not convert pending status to a passing claim.

## Skeptical findings

1. **The negative claim is supported at its stated boundary.** The antecedent
   is weak local committed-state preservation plus original-request at-most-once
   application. The witness does not silently redefine that antecedent as
   end-to-end correctness. The eight-event custody violation and ten-event
   recovered information loss are correctly distinguished. The latter has no
   result in either durable store, A's volatile state or the result channel.
   Original-request retransmission cannot reconstruct or reapply the result.

2. **The paired-world argument covers the complete declared recovery state.**
   `Recoverable` removes only `WorldResult`; identity, grant, activation,
   committed facts, all channels, and all volatile/durable receipt fields remain.
   I independently compared the exported JSON endpoints after removing that
   single field: they agree, and their removed actual values differ. Transition
   inspection confirms that the oracle field supplies neither guards nor bytes;
   future recovery cannot regain access through the model. Logs and researcher
   traces are explicitly outside the recovery interface. This establishes the
   scoped information-loss argument, conditional on opaque results and the
   absence of an additional observation source.

3. **The custody induction is noncircular.** A enables durable ACK using only its
   own recorded receipt. B enables release using only the authentic ACK channel;
   it does not inspect A's private state or the custody predicate. Authentic
   receipt provenance and monotone A retention carry the recorded fact through
   ACK delay, duplication and A's interruption. Application establishes the
   invariant and release is the only modeled destructive receipt step. The
   theorem is sufficient under these premises, not a necessity result.

4. **The model/oracle separation is useful but has an explicit limit.** The
   causal checker does not invoke the transition generator, explorer or state
   checker. Every reachable edge appended to a shortest source prefix is checked,
   and malformed traces challenge its acceptance boundary. Both implementations
   nevertheless implement the same declared specification; agreement cannot
   prove that specification refines ALAN or validate omitted behaviors. The
   artifacts say this explicitly. Inspection found no missing transition within
   the declared one-delegation/one-crash/one-slot domain that invalidates the
   result. The separate source-independent case argument supports the custody
   claim without treating shared-code agreement as proof.

5. **Reachability is not presented as liveness.** BFS retains revisits and cycles,
   and reverse reachability is seeded by A's consumption marker. The counts mean
   that some suffix consumes, not that every schedule terminates, B releases,
   channels drain or the pair becomes quiescent. The unfair retry cycle and
   explicit service premises preserve that distinction. One-slot channels,
   atomic B application/retention and fixed grants remain visible assumptions.

6. **Current Store implications are appropriately narrow.** Accepted admission
   calls `insertObservation` in the same transaction that stores the admission;
   `Commit` records consumption together with history, intents and successor
   state. These inspected code paths support the claimed contrast with volatile
   authored receipt in the abstract weak profile. Neither code inspection nor
   the finite model is represented as distributed-owner/process correspondence.

7. **The new mathematical text is synchronized.** Markdown and LaTeX contain
   matching premises, witness, indistinguishability argument, custody formula,
   proof cases and open obligations. The new Saltzer/Reed/Clark reference appears
   in both the Markdown list and BibTeX, and LaTeX cites its key. Existing section
   layouts differ and are not newly claimed identical. New repository links
   resolve. The parent's archived PDF visual review/build record supports the
   rendering statement; I did not independently render the PDF. Existing build
   warnings remain disclosed. Primary-literature inspection is attributed to
   the archived argument/manuscript review, not newly performed here.

## Evidence checked during this review

- Recomputed all 19 SHA-256 entries in `evidence/source.json`: no mismatch.
- Executed the retained exporter: exit 0, byte-identical `model.json` output.
- Executed its rejected-hypothesis mode: expected exit 1, `FALSIFIED` diagnostic,
  byte-identical `naive-red.json` output. This probe uses the retained binary;
  its build provenance and prior final-source tests are recorded separately.
- Confirmed the report's 407/1,848 and 311/1,372 state/edge counts, zero local
  state/transition violations in both profiles, delivery custody violations 64,
  and completion counts 391/16 versus 311/0.
- Read the archived normal/race PASS logs, clean Staticcheck/build records,
  intentional development failures, manuscript build output and prior reviews.
  These are task evidence, not test runs newly claimed by this reviewer.
- Checked new relative links and the current admission/commit source paths.
  A simple link scan also emitted two unchanged README image-title false
  positives; those are not broken links introduced by this change.

No full repository suite was rerun. Dynamic authority/cancellation, B or
simultaneous crashes, non-atomic B effects, mutual waits, executable runtime
correspondence and arbitrary-agent closure remain open as documented.
