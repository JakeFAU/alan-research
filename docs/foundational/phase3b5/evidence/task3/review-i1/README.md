# Superseded opportunity evidence

This directory preserves the evidence rejected by review finding I1 at commit
`331186c7b44cdc18b9f7a0970992fe00ce89714b`. The historical CLI serialized one
BA opportunity as `"AQ=="` at each N07 intermediate cut. The historical
comparison script interpreted those four characters as four separate BA
entries. See [review.md](review.md) and [provenance.json](provenance.json)
for the review, original output observations, and archived hashes.

The corrected comparison artifact is [../verdicts.json](../verdicts.json).
It records one BA entry in each intermediate N07 cut and an empty list in
each final cut. The historical report hashes include the recorded paths.
