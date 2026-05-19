# Issue 5: Renames, deletes, and binary/submodule stubs

## Parent

[PRD.md](../PRD.md) — review-diff v1

## What to build

Handle non-plain-modify entries in the change list and main pane:

- **Renames:** show old → new in the sidebar (from name-status)
- **Deletes:** selecting shows an appropriate stub or empty-old vs message (no broken side-by-side)
- **Binary / submodule:** list entry + main-pane stub (e.g. “Binary file changed”) instead of parsing as text

Integration tests in temporary repos covering rename, delete, and binary file in one compare.

## Acceptance criteria

- [ ] Rename appears in file list with both paths identifiable
- [ ] Deleted file selection does not crash; user sees a clear stub or one-sided view
- [ ] Binary (and submodule if detected) shows stub, not garbled text or parser errors
- [ ] Integration tests pass for rename, delete, and binary scenarios

## Blocked by

- #3 Changed-files sidebar and lazy per-file patches

## User stories

15, 16
