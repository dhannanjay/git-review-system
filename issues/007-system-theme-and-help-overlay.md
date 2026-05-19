# Issue 7: System theme and help overlay

## Parent

[PRD.md](../PRD.md) — review-diff v1

## What to build

UI polish aligned with the PRD:

- **Theme:** follow OS **light/dark** via `prefers-color-scheme` in Svelte/CSS
- **Help overlay:** `?` toggles a short panel listing default keybindings and flag reminders
- **CLI help:** complete `--help` for all v1 flags (including worktree flags if #6 merged) and pointer to keyboard defaults

## Acceptance criteria

- [ ] App respects system theme on Linux and macOS (manual spot-check; no automated pixel tests required)
- [ ] `?` opens/closes help overlay with file, hunk, and quit bindings
- [ ] `review-diff --help` documents `-C`, `--base`, `--head`, worktree flags, and default keys

## Blocked by

- #4 Keyboard navigation (files, hunks, quit)

## User stories

17, 21
