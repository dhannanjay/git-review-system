# Issue 3: Changed-files sidebar and lazy per-file patches

## Parent

[PRD.md](../PRD.md) — review-diff v1

## What to build

Expand the E2E flow to full branch scope without loading every patch up front:

- On session open, run **`git diff --name-status`** (with rename detection) and show an ordered **file list** in the sidebar
- **Lazy load**: when the user selects a file, fetch that path’s patch only and render side-by-side (reuse parser/UI from #2)
- Default selection: first **text** file in the list

Wails bridge exposes list + load-patch to Svelte; frontend does not call `git` directly.

## Acceptance criteria

- [ ] Opening a compare shows all changed paths in the sidebar before any per-file patch is loaded
- [ ] Selecting a file loads and displays its diff; switching files does not reload the full branch patch at once
- [ ] Integration test: temp repo with ≥3 changed files — only the selected file’s content is requested from git layer (mock or spy acceptable)
- [ ] Empty diff (no changes) shows an empty list and a clear empty state in the main pane

## Blocked by

- #2 E2E two-ref diff: first changed file side-by-side

## User stories

7, 8; completes 1 (open branch compare in GUI)
