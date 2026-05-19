# Issue 2: E2E two-ref diff — first changed file side-by-side

## Parent

[PRD.md](../PRD.md) — review-diff v1

## What to build

End-to-end path from CLI to UI for the core compare:

- CLI: `review-diff -C <repo> --base <ref> --head <ref>` using **`git diff` three-dot** (`base...head`) by default
- Backend: review **session** + **git runner** subprocess wrapper; resolve changed files enough to pick the **first text file** in the diff
- Load that file’s unified patch, **parse** into hunks with old/new line numbers
- GUI: **side-by-side** panes with add/remove line highlighting (no syntax highlighting)

Remote-tracking ref names (e.g. `origin/main`) work if present locally; no network calls.

Include **diff parser unit tests** with fixture patches.

## Acceptance criteria

- [ ] Against a test repo with a known branch pair, the GUI shows a side-by-side diff for the first changed text file
- [ ] Three-dot semantics verified by an integration test (merge-base behavior differs from two-dot in a crafted repo)
- [ ] Both columns show **line numbers**; added/removed lines are visually distinct
- [ ] Parser unit tests cover multiple hunks, additions, deletions, and line number mapping
- [ ] Invalid repo or unknown ref fails with a clear CLI error (non-zero exit) without hanging GUI

## Blocked by

- #1 Bootstrap Wails v2 app and blocking CLI

## User stories

2, 3, 8, 9, 10 (partial — single file, not full list yet); contributes to 1
