# Primary-source check for the new citation

Task 3 independently opened these primary sources with the web tool on
2026-11-14, after reading the controller's source notes:

- https://www.microsoft.com/en-us/research/publication/conjoining-specifications/
- https://lamport.azurewebsites.net/pubs/abadi-conjoining.pdf

The publication record identifies Martin Abadi and Leslie Lamport,
*Conjoining Specifications*, ACM Transactions on Programming Languages and
Systems, volume 17, issue 3, May 1995, pages 507-534. The author PDF front matter
identifies a 3 November 1995 revised version of SRC Research Report 118; its
article footer on printed page 507 says pages 507-533. This is a source-version
and pagination discrepancy, not silently reconciled. The BibTeX entry uses the
published metadata and the consulted author-copy URL. No DOI was verified.

The introduction distinguishes decomposition of a given system from composition
of reusable open components and makes environmental assumptions explicit.
Section 5.2 (printed page 529) supplies a composition theorem with explicit
hypotheses. The citation is background for hiding and assumption discharge;
Phase 3C.1's direct finite-path argument neither instantiates that general
formal theorem nor claims to have mechanized it.
