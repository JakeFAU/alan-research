# Phase 3C.4 Task 1 archive

[`qualify.py`](qualify.py) is the driver retained from the original capture.
**Do not rerun it in place.** It writes command records and captures beside
itself, so rerunning it can overwrite frozen evidence even if it exits early.
The archived driver has no overwrite guard.

See the [qualification reproduction instructions](../../05-qualification.md#checking-the-archive-and-recording-a-new-attempt)
for read-only checks, the guarded command for a fresh checker, export, and CLI
attempt, and the separate Go ordinary, race, and Staticcheck commands.
