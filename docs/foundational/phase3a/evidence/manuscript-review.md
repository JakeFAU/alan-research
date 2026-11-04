# Manuscript build and visual inspection

Command, from `paper/`:

```sh
tectonic --keep-logs alan-dissertation.tex
```

The command completed successfully in this task. The resulting PDF has 34 pages.
The new Phase 3A text has matching assumptions, counterexample, two-world
argument and conditional custody induction in Markdown and LaTeX. It does not
add a numbered theorem or alter existing theorem counters. Their pre-existing
section layouts differ; this check concerns the added mathematical content,
not a claim of full historical manuscript equivalence.

Saltzer/Reed/Clark (1984) was added to both the Markdown reference list and
`references.bib`. The primary MIT PDF supplies author names, journal, volume,
issue, year, page range and the delivery-guarantees discussion; the new text
attributes the established principle without claiming it is novel.

Poppler rendered pages 2, 31, 32, 33 and 34 at scale-to 1600. They were visually
inspected: the new table-of-contents entry, section boundary, formula, paragraphs,
following conclusion and new bibliography entry are legible without clipping
or overlap. The new section begins on page 31 and continues on page 32.

`tectonic.log` retains the engine output. The build reports underfull boxes at
unchanged lines 215 and 559, and a 5.2184pt overfull box at unchanged line 559.
No warning points to the new Phase 3A section. This is not a warning-free build;
unrelated existing formatting was left outside this research slice.
