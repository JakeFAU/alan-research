# Reviewed Phase 3B.4 PDF build

The installed `paper/alan-dissertation.pdf` is the 41-page A4 artifact built from substantively approved research commit `5db0ff147ea4f16b037976a1a6ff968202b2d0fb`. The later locator-disclosure commit changes no manuscript source. PDF SHA256: `c79be8583f18baf3832c7059b4035fd509f9e4d26a8860f1b8eb544c1a0fe4ae`.

`build-record.json` records the exact cached Tectonic command, exit0, before/after source hashes and prior/new PDF hashes. `tools.json` records Tectonic0.17.0 and Poppler26.10.0. `build-and-renders.tar.gz` retains the exact source inputs, build helper and logs, extraction, metadata, render commands and PNGs. `manifest.json` verifies every archive member against its source bytes and the installed PDF identity.

The controller visually inspected all ten retained pages:1,2,9,28,29,37–41. The new section, updated Abstract/Conclusion and continuations are legible without clipping, overlap or missing glyphs. The build retains five unique box warnings in unchanged earlier passages; exact warnings and source comparison are retained and their pages inspected. This is not a warning-free build.

The new section is Markdown6.12 / LaTeX-PDF6.7, contents page2 and body37–39. Existing numbering and the prior Markdown-only Formal Falsification section remain different; only the affected mathematical text is synchronized. See `visual-review.md` and the current results/review index. The initial substantive review predates this build; final Task6 artifact closure and broad whole-change review are separate recorded gates.
