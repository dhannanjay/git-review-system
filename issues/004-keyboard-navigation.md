# Issue 4: Keyboard navigation (files, hunks, quit)

## Parent

[PRD.md](../PRD.md) — review-diff v1

## What to build

Keyboard-driven review for the file list and active diff:

- **Files:** `j` / `k` or `]` / `[` — next / previous file in the sidebar (updates selection and lazy-loads patch)
- **Hunks:** `n` / `p` — next / previous hunk in the current file (scroll or focus)
- **Quit:** `q` — close window; CLI unblocks

Bindings work when the main window is focused; avoid stealing keys from text inputs if any exist.

## Acceptance criteria

- [ ] User can move through all changed files without the mouse
- [ ] User can jump hunk-to-hunk in a multi-hunk file without the mouse
- [ ] `q` closes the app and returns exit code 0 on clean quit
- [ ] Documented in `--help` and briefly in README (full overlay comes in #7)

## Blocked by

- #3 Changed-files sidebar and lazy per-file patches

## User stories

11, 12, 13
