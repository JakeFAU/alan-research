# Phase 3B.7 manuscript build and visual review

Cached Tectonic 0.17.0 built the synchronized manuscript with
`tectonic --only-cached --keep-logs alan-dissertation.tex`. The command
exited 0 in 1.007 s. The source files and bibliography hashes remained
stable during the build. The installed A4 PDF has 44 pages and SHA-256
`7fc751e059f61423dc307e346826dd17099ac3d99c8a8787dd5d48fed17e8f76`.

Poppler's layout-text comparison against the prior 43-page PDF identified
changes on pages 1, 2 and 41–44. I visually inspected those PNG renders at
110 DPI and page 40, which starts the Phase 3B.6 section that continues
onto page 41. The Abstract and table of contents fit on pages 1–2. The
Phase 3B.7 subsection begins on page 41, and its post-capture correction
and causal limits continue on page 42. The Conclusion spans pages 42–43
and the references continue onto page 44. Headings, math, links as
typeset, wrapping, margins and page numbers were legible, with no
clipping, overlap or missing glyphs observed in those renders.

The Abstract, Phase 3B.6 closing paragraph, new Phase 3B.7 subsection
and Conclusion match after ordinary formatting and link/path
normalization. Phase 3B.7 is Markdown Section 6.15 and LaTeX/PDF Section
6.10. The Conclusion is Markdown Section 8 and LaTeX/PDF Section 7.
The inherited Markdown-only Section 7 remains byte-identical and absent
from LaTeX/PDF. The bibliography is unchanged. The same five unique
TeX box warnings on unchanged lines 214–215 and 558–559 are present in
the Phase 3B.6 build evidence.
[Source, build and archive cross-checks](checks.json)

The build archive contains the copied sources, PDF, build command and
record, TeX and BibTeX logs, tool versions, extracted text, PDF metadata,
render arguments and inspected PNGs. The manifest binds the archive and
each member. The PDF edit operation marker ran exactly once before the
first manuscript source edit, with exit 0 and empty output.
[Build and render archive](build-and-renders.tar.gz) ·
[Manifest](manifest.json)
