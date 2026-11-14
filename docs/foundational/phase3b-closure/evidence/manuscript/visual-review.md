# Manuscript build and visual inspection

The PDF skill edit marker ran exactly once before manuscript edits, with exit 0. Dependency discovery located the bundled Node/Poppler tools. Tectonic ran as:

```
tectonic --only-cached --keep-logs /tmp/alan-phase3b-closure/manuscript/alan-dissertation.tex
```

Exit 0. The scratch directory held copies of the updated Markdown/TeX and unchanged bibliography. The resulting 45-page A4 PDF was installed in `paper/alan-dissertation.pdf`. The prior PDF had 44 pages and SHA-256 `7fc751e059f61423dc307e346826dd17099ac3d99c8a8787dd5d48fed17e8f76`.

`pdftotext -layout` comparisons identified changed pages 42, 43, 44 and 45. Bundled `pdftoppm -f 41 -l 45 -r 100 -png` rendered those pages plus their preceding continuation. The author opened all five PNGs: the Phase 3B.7 text and post-capture correction continue from 41 to 42; the new closure text and Conclusion heading fit on 42; the Conclusion continues on 43; references continue on 44–45. Text, math, path wrapping and page numbers were legible, without observed clipping, overlap or missing glyphs. Page 45's shorter reference continuation is intentional pagination.

The build retains five unique inherited warnings: underfull badness 2744 at lines 214–215; underfull badness 5161, 1087 and 3229 plus overfull 5.2184pt at lines 558–559. They recur on TeX passes and match the prior packet. No changed source occurs at those lines.

The inherited Markdown-only Section 7 and all earlier section-number differences are preserved. Added prose is exactly equal after removing the Markdown link versus TeX path syntax. The bibliography is unchanged.

[Build sources, logs, extracted text and five inspected PNGs](build-and-renders.tar.gz) · [Archive/member/output hashes](manifest.json)
