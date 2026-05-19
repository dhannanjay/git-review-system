# Implementation issues (v1)

Vertical slices for [PRD.md](../PRD.md). Work in dependency order; **Blocked by** references use issue numbers below.

| # | Title | Type | Blocked by |
|---|--------|------|------------|
| 1 | [Bootstrap Wails v2 app and blocking CLI](001-bootstrap-wails-app-and-blocking-cli.md) | HITL | — |
| 2 | [E2E two-ref diff: first changed file side-by-side](002-e2e-two-ref-diff-first-file-side-by-side.md) | AFK | 1 |
| 3 | [Changed-files sidebar and lazy per-file patches](003-changed-files-sidebar-lazy-patches.md) | AFK | 2 |
| 4 | [Keyboard navigation (files, hunks, quit)](004-keyboard-navigation.md) | AFK | 3 |
| 5 | [Renames, deletes, and binary/submodule stubs](005-renames-deletes-binary-stubs.md) | AFK | 3 |
| 6 | [Linked worktree compare](006-linked-worktree-compare.md) | AFK | 2 |
| 7 | [System theme and help overlay](007-system-theme-and-help-overlay.md) | AFK | 4 |
| 8 | [Release builds and Arch PKGBUILD](008-release-builds-and-arch-pkgbuild.md) | AFK | 2 |

**v2 (deferred):** [Optional `--fetch` before diff](v2-optional-fetch-before-diff.md)
