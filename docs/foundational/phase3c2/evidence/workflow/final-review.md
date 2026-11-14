# Phase 3C.2 whole-branch adversarial review

Reviewed base `69534188f796efab2e866750f6a76e2c7a3a1706` through head
`a86146fe041ede25753bf470205c1cb5ecf565f2`, including all five commits. HEAD
matched the requested head. The supplied review package's diff body exactly
matches `git diff -U10 69534188 a86146fe`.

## Question and scope

Do the counterexamples, conditional source argument, executable resource-aware
relation, separate checker, retained evidence and manuscript describe the same
fixed-chain result without importing a runtime or unrestricted composition claim?

Yes, within the frozen premises. This review independently examined the whole
logical change, rather than treating the task-review verdicts as proof. It read
the contract, completed plan, progress ledger and review record; all new Go
implementation/test/exporter source; the Python checker, tests and capture
tooling; argument/schema/results/qualification; manuscript changes and relevant
unchanged delegated semantics. Generated archives were checked by exact content
and provenance checks described below rather than by printing every JSON record.
The prior reports were used for attribution and historical concerns, not as an
oracle for the present semantic assessment.

## Strengths and substantive assessment

1. **The adverse witness proves the claimed information loss.** At
   `research/capacitycomposition/evidence.go:89`, both holder prefixes are built
   from genuine empty custody using legal transitions. Their complete base
   records match; the different original-held bits are not manufactured by
   editing an endpoint. With both originals waiting and serial B, AB holding Q
   prevents both applications, whereas BC holding Q has a finite completion
   suffix. Transport and recovery cannot release Q. Ready-now genuinely allows
   the AB acquisition while BC is U, and commitment of BC invalidates readiness.

2. **The positive argument has the necessary premises.** The source proof at
   `docs/foundational/phase3c2/01-argument.md:77` and `:113` separates correspondence
   from stable-ready sufficiency. In `research/capacitycomposition/model.go:145`,
   stable readiness follows from immutable servicing or absorbing outgoing D.
   Application is the only hold-clearing transition. Existing holders can be
   applied, applied originals consumed without Q, then BC and AB finished in
   dependency order. This construction needs the stated atomic root update,
   atomic application/result-retention/Q-release, durable custody, fixed exact
   grants, sufficient E and recoverable-owner premises. It establishes a finite
   possible suffix, not an inevitable schedule or runtime result.

3. **Lifting uses the correct abstraction and actual custody.** At
   `research/capacitycomposition/concrete.go:83`, an action must be enabled in the
   resource-aware relation. Reserve prepares real admission; a held apply is
   already admitted; consume obtains bytes from recorded, volatile, channel or
   retained service custody. Preparation changes neither another stage nor a
   hold. The unchanged contractcomposition/twoagent source supports these steps:
   service work and receipt transport do not require the remote caller to be up,
   durable ACK release leaves a durable caller source, and recovery preserves
   admission/holds. WorldResult is carried as historical data and written by
   actual application, never tested to enable an operation or choose custody.

4. **Stable-ready's observable cost is retained.** The old AB-before-BC-commit
   trace is legal under early reservation. Stable-ready cannot lift it while
   hiding only reservation. The Python reserve-only epsilon-closure check at
   `docs/foundational/phase3c2/tools/audit.py:334` tests that trace language fact;
   visible BC progress is not silently treated as preparation. The argument,
   results, roadmap, Abstract, new section and Conclusion preserve this limit.

5. **Composition preserves identity without inventing a second pool.**
   `research/capacitycomposition/compose.go:73` validates fixed memberships and
   roots, joins equal IDs once, and relies on the base owner/port compatibility
   checks. Close validates the exact chain. Grouping is over static declarations;
   it does not merge independent dynamic histories. Fixed-root validation and
   cloned slices avoid the copied-pool and append-aliasing hazards in this path.
   No general open-component substitution result is implied.

6. **The independent checker is substantive but accurately bounded.** At
   `docs/foundational/phase3c2/tools/audit.py:145`, all graph arcs and reverse
   completion labels are reconstructed from the specification. At `:218` and
   `:308`, concrete post-states are derived from predecessor events and then
   compared field by field. Exact type checks distinguish Booleans from counts.
   The checker neither imports nor executes Go semantics. It shares the frozen
   specification, not an independently discovered specification. Probe metadata
   remains a structural report; the 144 Python-generated histories share Python
   replay rules and are correctly called checker fixtures.

7. **Evidence categories and manuscript attribution remain aligned.** The
   72/1,632/2,412 abstract roster, eight noncompletable states, 1,854,036 factored
   probe cases, ten captured histories/161 steps and 144 Python fixtures/7,200
   steps are not conflated. `docs/foundational/phase3c2/04-results.md:56` and
   `05-qualification.md:20` correctly distinguish fresh scoped checks from
   historical source-bound package evidence. The new manuscript content is
   aligned across Markdown/TeX. The old Phase 3C.1 disjoint-resource premises
   remain intact; this branch does not retroactively qualify them under shared Q.

## Findings

**Critical:** none.

**Important:** none.

**Minor:** none requiring a change. Task 1 M1 is resolved: the README race command
contains the recorded `-timeout 30m`. Historical task-status wording in the frozen
argument remains appropriately superseded by current results/qualification and
is explicitly explained in the Task 3 report.

## Checks actually executed

These checks were read-only except this report. No Go suite, exporter, Python
semantic suite/audit, build, module operation or PDF rebuild was rerun.

- `cat`, bounded `sed`, `rg`, `git diff --stat`, `git diff --name-status`,
  `git rev-parse HEAD`, and `git status --short` for the assigned scope and
  dependencies. The working tree was clean before this ignored report.
- `git diff --check 69534188 a86146fe`: exit 0.
- `PYTHONDONTWRITEBYTECODE=1 python3 docs/foundational/phase3c2/evidence/task2/capture.py --verify`:
  exit 0; all 201 artifact/specification records and packed/raw archive bindings
  checked; 32 producer hashes match; Task 1 raw export SHA256 is
  `a7cf316219b30cdc83d455b95a16549c022d55aa422b642b8f5976315e0c0fd2`.
- `python3 docs/foundational/phase3c2/evidence/manuscript/verify.py`: exit 0;
  637 historical inputs, 20 supplemental inputs, 12 new Go inputs, 32 producer
  inputs, 201 Task 2 artifacts, 23 qualification records, eight compressed
  manuscript artifacts and 29 local links match. Changed section normalization,
  Abstract/Conclusion equality, paper hashes and race timeout agree. Visual
  attestations from this verifier were not mistaken for fresh visual inspection.
- Independent inline `python3 -B` using `git cat-file --batch`: all 657 inherited
  source/module inputs match current bytes and merged-base blobs; all 20
  supplemental inputs also match the specified historical baseline commit.
  Scope inventory confirms 302 changed paths: five existing documentation/paper
  files modified; every addition belongs to phase3c2 or capacitycomposition.
  No old model, production, dependency or historical evidence file changed.
- Independent archive checks: five preflight and ten workflow records match
  their hashes, including decompressed brief hashes/lengths; Task 1/2 exports
  are byte-identical; all 43 hostile inputs match their recorded raw hashes,
  exit 1, expected diagnostic, exact stderr and empty stdout. The removed
  non-tree-arc control retains 24 states and every one of its 23 BFS parent arcs,
  while its corrected actual count is 33 versus the genuine 34.
- Inspected retained ordinary/race test output, scoped lint/export logs,
  independent audit and unittest output, capture/qualification command metadata
  and prior review dispositions. Passing results remain attributed to those
  actual runs, not this reviewer.
- Independent inline PDF checks using `pdftotext -layout paper/alan-dissertation.pdf -`:
  output exactly equals retained compressed current text. Actual warning lines
  extracted from both TeX logs match each other and the seven disclosed warnings.
  All 24 existing PNG hashes match the render manifest. Personally viewed
  existing pages 44 and 45: the complete new section has readable text, intact
  margins and no observed clipping, overlap or missing glyphs. No new rendering
  or authoring marker was needed for this read-only check.
- Compared the supplied diff body with `git diff -U10 69534188 a86146fe`: exact
  equality. Initial comparison with default three-line context failed because
  the package uses ten lines of context; this was a reader assumption, not a
  missing change. Another exploratory integrity snippet assumed all preflight
  manifest records had byte lengths; it failed with KeyError and was corrected
  to check lengths where declared. Both corrected checks passed without any
  repository change.

## Declined judgments and exact remaining obligations

These items were sent to the controller before the verdict for explicit
disposition. They do not disappear because this review approves the bounded
packet.

1. **Fresh repeat execution of unchanged suites/exports/builds.** Declined because
   passing captures are source-bound and the brief prohibits repetition merely
   for confirmation. No present acceptance work remains on these commands;
   changed relevant source/dependencies or a new unexplained failure would
   require appropriate requalification. Historical runs remain historical.
2. **Independent enumeration of the million-case probes or entire concrete
   product; mechanized universal proof.** Not performed. A stronger claim needs
   the corresponding independent enumeration or formal proof, including all
   custody/availability combinations. The current universal statement remains
   a conditional source argument supplemented by nonexhaustive factored probes.
3. **Runtime/distributed and excluded semantic behavior.** No judgment of runtime
   implementation, distributed allocation, unknown split effects, changing
   grants/membership/policies/generations, cancellation/revocation, permanent
   owner loss, cycles/arbitrary topology or contexts, open-component
   substitution, fairness/divergence transfer, exact-cost equivalence, reduced
   verification cost, stable-ready necessity or optimality. Each requires its
   own contract and proof/experiment before any such claim; none is delivered
   here. E is assumed sufficient rather than modeled as a balance.
4. **Global inherited manuscript and historical-theorem audit.** This review
   checks changed claims and relevant delegated contracts. It does not repair
   Markdown-only Section 7, inherited numbering/bibliography differences or the
   seven existing TeX warnings. Global consistency remains a separate obligation
   if requested; the branch expressly disclaims global equivalence.
5. **Fresh PDF build or all-page visual review.** No new build/render was needed;
   actual text/log/hash integrity was checked and pages 44–45 viewed. The other
   22 rendered pages retain Task 3 worker/reviewer visual attribution. A claim
   that this reviewer independently inspected them would require that work.
6. **New external literature/DOI verification.** Metadata is checked against the
   explicitly attributed controller inspection and aligned bibliography. No
   new browse or DOI verification occurred. The retained 403 limitation stays;
   a stronger verified-DOI claim would need a successful primary-source check.
7. **Independent RED chronology or historical environment reconstruction.**
   The preserved failures and run records were inspected, not re-created from
   preimplementation snapshots or old host environments. Exact current/base
   source identity supports the stated reuse, not a claim of independently
   reproduced historical execution. Such a chronology/environment claim would
   need corresponding snapshots and environment evidence.

## Verdict and readiness

The controller explicitly dispositioned all seven declined items before this
final verdict: retained same-source runs discharge current scoped verification;
probe/product/mechanization limits remain explicit; all expanded-domain/runtime
claims remain unestablished; changed-content manuscript alignment is the
accepted scope; visual and literature attribution stays exact; and RED/historical
execution remains attributed evidence rather than reconstructed chronology or
environment reproducibility. These are scope dispositions, not waivers of a
found defect. This reviewer agrees they leave no missing acceptance gate within
the frozen contract.

**APPROVE** for draft-PR delivery of this research-only logical change at the
reviewed head, preserving the stated assumptions, limitations and attributions.
No actionable finding remains from this whole-branch review. The controller
must record its disposition of the declined items and update the pending
whole-branch review/delivery status without rewriting the historical reviews.
This is not proof by reviewer agreement, runtime qualification, general Phase
3C closure, approval of later semantic edits, or merge authorization.
