# Phase 3B.7 whole-branch adversarial review

Reviewed base `d8645d739445d4b09c0d38dca4454ca6297fdf22` through exact head
`98c9b04b8c5c62a249357fac669b0d7b2e5aefef`, following the supplied final-review
instructions, frozen contract, research working agreements and reviewer template.
The supplied package's diff is byte-identical to `git diff -U10` for that range.
HEAD matched the reviewed head and the tracked checkout was clean.

## Strengths and substantive assessment

The branch supplies the requested concrete matched experiment within its
declared scope. Each of the normal and race captures contains a genuine
admitted-but-unapplied history and an applied-without-local-result history.
The seed copies match by opaque file-set hashes. The complete post-open,
stopped, reopened and final exports and complete final Run reports match byte
for byte within each pair. The model projection is separately checked; it is
not used to hide differences in the richer local observation.

The causal apparatus supports the distinction. The worker's pre-call observer
precedes `client.Do`; both branches synchronize that observer with an
acknowledged scheduler hold. The zero worker is never released. The applied
worker is released only after the scheduler acknowledgment, and the endpoint
post-commit acknowledgment precedes confirmed owner death. Successful endpoint
response release follows that death. The owner excludes test-context
cancellation from its scheduler lifetime, treats unexpected scheduler return
as an abort, and keeps release pipes until reaping. The endpoint callback may
leave with a cancellation error; continuous callback occupancy is correctly
not claimed. I checked these interactions against the pinned owning source,
including the inherited observer, transport wrapper, pre-call emission,
endpoint journal/commit/barrier order, disposable-copy capture and controlled
revocation helpers.

The retained endpoint dossiers positively distinguish zero contacts and
applications from one contact and application. Both endpoints have complete
stopped and final coverage, and the applied receipt is explicitly idempotent.
The raw binding adapter preserves the existing authority and continuity
checks: original request/grant/root/wire identity, historical allowance and
debit, exact two-principal renewal, exact-b revocation, current-epoch fresh
refusal, surviving c, and absent selected local receipt, settlement and
consumption. The adapter does not pretend to be a second implementation of
that prior validator.

The comparison implementation preserves ordered arrays, exact strings and
number tokens, null/empty distinctions and duplicate-key rejection. Required
native endpoint coverage and complete final reports are validated before
support is certified. The retained synthetic controls cover missing evidence,
coherently resealed false inputs, contradictory identities, wrong chronology,
same-applied pairs, and richer state differences hidden by a model projection.

The post-capture correction is justified by the frozen source, rather than by
changing the intended observation to fit the captures. The transport wrapper
sets the selected authority admission ordinal to 1 before emitting the passive
authorization event. The corrected checker requires exactly 1 and its join to
the worker; it does not accept both 0 and 1. The original rejected reports,
original binary and freeze, mistaken prior review interpretation, correction
and re-review remain available. The six-path freeze delta changes no producer,
contract, equality comparator or raw evidence.

The changed Markdown and LaTeX claims accurately distinguish this operational
witness from the finite model and from unproved distributed/recovery claims.
The body and Conclusion identify the idempotent endpoint, common closed seed,
sequential isolated processes and shared Store owner. The Abstract, body and
Conclusion disclose acceptance after a checker correction. The PDF and source
archive bind to the submitted manuscript bytes. No new literature claim or
bibliography change is introduced.

## Issues

- Critical: none found.
- Important: none found.
- Minor: none requiring correction for this branch's stated experiment.

No substantiated prior finding is waived. The four original Task 2 REJECT
findings are resolved in the current code: explicit complete zero-world
coverage, complete native Run schema, full native authorization coordinates,
and effective unique test selection. The subsequent ordinal mistake is
separately resolved by the source-supported round-two correction. The original
reports remain historical evidence, including their erroneous intermediate
interpretation.

## Checks actually performed

Read-only commands used Git inspection, `cat`/`sed`/`rg`, and inline Python
standard-library parsing, gzip/tar decoding and hashing. No new runtime
capture, checker execution, Go test/race suite, build, lint campaign or SQLite
open occurred. This report is the only file written by this reviewer.

- Verified the exact review package against the specified Git range and read
  the producer, checker, adapter, controls, protocol, contract, results,
  qualification, runbook, manuscript changes and archived review dispositions.
- Rehashed all 629 current corrected-freeze inputs and independently computed
  the six-path delta from the original 627-input freeze.
- Verified all 637 pinned runtime Git blobs against the producer manifest,
  including their lengths and SHA-256/Git blob identities, and all four overlay
  length/hash entries against current source. Both captured source manifests
  equal the frozen manifest. No `.envrc` was read.
- Rehashed both still-present original and corrected checker binaries against
  their distinct recorded identities. These checks establish consistency, not
  independent binary-to-source attestation.
- Decoded both raw archives in memory and checked their archive lengths and
  digests, all 122 regular member names/modes/lengths/digests, and their complete
  121-member inner seals. Strict JSON parsing rejected duplicate keys in the
  verification path; journal records had complete line endings. SQLite
  members were treated only as opaque bytes.
- Directly checked seed file-set equality, all four paired export byte
  equalities, final report equality, stopped/reopened/final run coordinates,
  and absent external receipts. Independently checked consecutive transcript
  ordinals, acknowledgment/release/death order, SIGKILL/reap results, absent
  abort artifacts, endpoint journal cardinality, stable stopped/final endpoint
  snapshots, and idempotent zero/one application facts.
- Checked all eight retained checker reports against the complete exact
  input-member digest sets. Initial reports preserve Protocol=false and
  Supported=false; corrected reports preserve support. Seal/Source/Raw/World
  pass in both revisions. Repeated reports are byte-identical per mode and
  revision; recorded direct exits are respectively 1 and 0.
- Verified all 23 indexed command records and their indexed file hashes,
  decompressed output lengths/digests and direct exits. Independently counted
  422 test/subtest pass events in each retained affected ordinary/race corpus
  log, with no failures or named skips. Each capture log records exactly one
  prescribed zero/applied pair. Retained correction GREEN logs contain no
  failure or race markers. These are inspected prior executions, not reruns.
- Verified the archived workflow manifest and the manuscript archive's
  complete member inventory. Archived Markdown, LaTeX, BibTeX and PDF bytes
  equal the current files. Read the retained visual-review/build account and
  compared the changed manuscript claims. PDF visual inspection is attributed
  to the worker and Task 3 reviewer; I did not rerender or independently view
  pages in this final review.

One initial review-only assertion compared the supplied diff against Git's
default three-line context and failed. Inspection showed the supplied package
uses ten-line context; comparison with `git diff -U10` passed exactly. This
was a comparison-command assumption, not a package discrepancy.

## Declined to judge

- Artifact authenticity or independent binary-to-source attestation: the
  experiment explicitly assumes an honest pinned producer and endpoint;
  unsigned metadata and hashes establish internal consistency only.
- Independent SQLite-to-JSON content equivalence: opening captured databases
  is prohibited, and the frozen disposable-copy/export procedure is the
  declared trust boundary. I inspected that procedure and byte custody.
- Independently provisioned or concurrent cloned deployments, distributed
  ownership, arbitrary revocation, non-idempotent endpoint behavior, general
  recovery equivalence, liveness and N-agent closure: these are explicit
  exclusions, and the submitted claims do not present them as established.
- Exhaustive adversarial input security or a complete proof of the legacy
  authority/model validators: this task adds a bounded offline protocol and
  narrow adapter; relevant integration paths and negative controls were
  reviewed without representing them as universal verification.
- Full forensic-dossier, physical SQLite-layout or timing equality: these are
  excluded from the predeclared observation; process and endpoint evidence
  must distinguish the branches.
- Repair of the inherited Markdown-only Section 7/TeX numbering divergence
  and global Staticcheck debt: these are disclosed pre-existing boundaries,
  unchanged by this research slice. No full manuscript equivalence or globally
  clean lint claim is made.
- New statistical confidence or repeatability rates: two prescribed pairs
  establish bounded examples, not a statistical campaign.
- Fresh broad/race qualification, new captures and another PDF rendering:
  existing scoped evidence and independent Task 3 review cover those checks;
  the dispatch explicitly avoids duplicating campaigns. No concrete unresolved
  doubt required a new probe.
- Remote draft-PR delivery or merge readiness of a later head: these are
  controller delivery steps beyond this exact-head review. No merge is
  authorized by this report.

## Recommendations and remaining delivery checks

Archive this report unchanged, record the controller's disposition of every
declined-to-judge item, update the pending final-review entry in
`docs/foundational/phase3b7/07-reviews.md`, and complete the planned draft-PR
delivery. Inspect the final status/stat/diff after those evidence-only changes
and confirm that the 629 frozen inputs, raw archives and manuscript remain
unchanged. The final report/PR should identify the exact reviewed head and
subsequent delivery-only commit, retain the post-capture correction caveat and
source-conditioned limits, and leave merging to the user. These are delivery
obligations, not unresolved implementation defects.

## Assessment

**Specification: APPROVE.**

**Quality: APPROVE.**

**Whole-branch verdict: APPROVE for the bounded research result and draft-PR
delivery.** The evidence supports the matched zero/one witness under its
explicit source and observation assumptions, with no actionable defect found.
This is a bounded adversarial audit result, not a mathematical proof or merge
authorization.
