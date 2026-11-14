# Task 1 evidence index

Use `capture-3/export.json.gz` as the final fixture and `capture-4/` as the
independent fresh deterministic repetition. Both have clean dependency enumeration
and exporter command exits. Their semantic bytes equal captures 1 and 2 too.
`final-capture-equality.json` records all byte identities;
`final-dependency-validation.json` checks clean stderr, no package errors and the
module-relative compiled roster. `final-overwrite-refusal.json` confirms no
existing capture file changed in content, size or modification time.

Captures 1 and 2 are preserved verbatim: go-list exited zero but reported a
nonfatal permission warning writing the default module stat cache. The helper
was corrected to use `/tmp/alan-phase3c5/mod-cache`; its original source and the
resolution are in `capture-diagnostics/`. Go code, argument and schema did not
change. The helper is not a producer input and therefore does not enter semantic
export hashes; final helper bytes are separately bound by capture environment
records and `source-freeze-final.json`. `source-freeze.json` remains the initial
freeze, including the original helper hash. No accepted semantic bytes were edited.

`qualification/` retains final fresh ordinary/race/Staticcheck command records.
`red/` retains the initial semantic failures against the insufficient equality/FULL
baseline, with exact source snapshots. `false-hypotheses/` retains three separately
labeled intentional assertion failures against the final model, distinct from
passing counterexample tests. `measured-summary.json` adds derived coverage and
pair details, without replacing the complete graphs. The artifact manifest binds
all files below this directory except itself.

These are finite model and source-argument artifacts. Independent Python audit,
root qualification, manuscript/PDF reconciliation and delivery are later tasks.
