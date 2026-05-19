# Issue 6: Linked worktree compare

## Parent

[PRD.md](../PRD.md) — review-diff v1

## What to build

CLI and session support for comparing across **linked git worktrees**:

```
review-diff -C <repo> --base-worktree <path> --head-worktree <path> [--base <ref>] [--head <ref>]
```

- Default each side to that worktree’s **`HEAD`** unless `--base` / `--head` overrides that side
- **Validate** all worktree paths share the same `git-common-dir` as `-C`; fail fast with a clear error if not
- Reuse existing GUI and diff pipeline (#2+); no UI redesign required beyond passing resolved refs

Integration tests: two worktrees on one repo, compare; negative test with unrelated directories.

## Acceptance criteria

- [ ] Compare two worktrees checked out on different branches shows expected diff in GUI
- [ ] `--base` / `--head` override ref for the corresponding worktree side when provided
- [ ] Unlinked worktree paths error at CLI with actionable message
- [ ] Integration tests cover valid linked pair and invalid unrelated paths

## Blocked by

- #2 E2E two-ref diff: first changed file side-by-side

## User stories

4, 5, 6
