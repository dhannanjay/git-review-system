# PRD: review-diff — Local Git Diff Reviewer (v1)

## Problem Statement

Reviewing changes between Git branches is slow and heavy when done through a browser, a full IDE, or raw terminal diffs. Developers need a **fast, lightweight desktop viewer** on Linux and macOS that focuses on **reading** branch-to-branch changes—including **linked git worktrees**—without the bloat of a full Git hosting UI or IDE.

## Solution

**review-diff** is a CLI-first desktop app: run `review-diff` with a repo path and two refs (or two worktree paths); a **native GUI** opens showing a **changed-files list** and **side-by-side** diffs with keyboard navigation. The tool uses local Git only (no network in v1), loads each file’s patch **lazily**, and exits when the user quits—returning control to the shell.

## User Stories

1. As a developer, I want to run `review-diff -C <repo> --base main --head feature` from my terminal, so that a GUI opens immediately with the PR-style diff between those branches.
2. As a developer, I want the default comparison to use the three-dot range (`base...head`), so that I see what my branch would merge into the base (like a pull request).
3. As a developer, I want to pass remote-tracking refs such as `origin/main` as `--base` or `--head`, so that I can review against my last fetch without the tool calling the network.
4. As a developer, I want to compare two linked git worktrees via `--base-worktree` and `--head-worktree`, so that I can review what is checked out in each worktree without manually naming branches.
5. As a developer, I want optional `--base` / `--head` overrides alongside worktree paths, so that I can compare a specific ref in one worktree against another worktree’s `HEAD` or ref.
6. As a developer, I want the tool to reject worktree paths that are not linked to the same repository, so that I do not get confusing or empty diffs from unrelated directories.
7. As a developer, I want a sidebar list of all changed files on open, so that I can see the scope of the branch before reading any hunks.
8. As a developer, I want only the selected file’s patch loaded at a time, so that opening a large branch stays responsive.
9. As a developer, I want a side-by-side view of old and new lines with add/remove highlighting, so that I can scan changes efficiently without syntax coloring noise in v1.
10. As a developer, I want line numbers on both the left (old) and right (new) columns, so that I can correlate lines with blame, issues, or chat.
11. As a developer, I want to jump to the next or previous hunk with `n` and `p`, so that I can move through large files without scrolling blindly.
12. As a developer, I want to jump to the next or previous file with `j`/`k` or `]`/`[`, so that I can review many files in sequence from the keyboard.
13. As a developer, I want to quit with `q`, so that I can return to my terminal workflow quickly.
14. As a developer, I want the CLI process to block until I close the GUI, so that scripts and shell habits get a predictable exit.
15. As a developer, I want renames and deletes shown clearly in the file list, so that I do not miss structural changes.
16. As a developer, I want binary and submodule changes called out as stubs instead of broken panes, so that I know they changed without opening unusable views.
17. As a developer, I want the UI to follow my system light/dark preference, so that the app feels native on macOS and Linux desktops.
18. As a developer on macOS, I want to install from a release `.app` bundle, so that I can run the tool like other desktop apps.
19. As a developer on Linux, I want a release tarball with a single binary, so that I can put it on my PATH without a heavy runtime.
20. As an Arch Linux user, I want an AUR `PKGBUILD` that wraps the official Linux tarball, so that I can install with my usual package workflow without a separate Arch-specific build.
21. As a developer, I want `review-diff --help` to document flags and default keybindings, so that I do not need external docs for daily use.
22. As a developer (v2), I want an optional `--fetch` flag to refresh remotes before diffing, so that `origin/*` refs are up to date when I choose to use the network.

## Implementation Decisions

### Architecture overview

Single Go module, **Wails v2** application with embedded **Svelte** frontend. One binary **`review-diff`** exposes a CLI subcommand that validates arguments, prepares a **review session**, and launches the Wails window. The GUI communicates with Go via Wails bindings/services; the frontend never shells out to Git directly.

### Modules (deep interfaces)

| Module | Responsibility | Public interface (conceptual) |
|--------|----------------|------------------------------|
| **CLI** | Flag parsing, validation, process lifecycle (block until GUI exit), help text | `Run(args) → exit code` |
| **Review session** | Immutable context: repo path, base/head refs, optional worktree paths, diff range mode (three-dot default) | `NewSession(opts) → Session` or error |
| **Git runner** | Subprocess `git` with `-C`, consistent env, timeout policy for hung commands | `Run(ctx, repo, args...) → stdout, err` |
| **Worktree validator** | Confirm paths share `git-common-dir`; resolve `HEAD` for worktree paths | `ValidateLinked(repo, wtA, wtB) → error` |
| **Change list** | Produce ordered file list with status (modified, added, deleted, renamed, binary) | `ListChanges(session) → []FileChange` |
| **Patch loader** | Lazy unified diff for one path | `LoadPatch(session, path) → Patch` |
| **Diff parser** | Turn unified diff into hunks with old/new line ranges and lines | `ParseUnified(patchText) → PatchModel` |
| **Wails bridge** | Bind session + list + load patch to frontend; emit errors | JSON-friendly DTOs for Svelte |
| **Svelte UI** | File list, side-by-side panes, keyboard routing, system theme CSS | Consumes bridge API only |

**Three-dot default:** For `--base` and `--head` ref names, use `git diff base...head` (merge-base semantics). Document and optionally support two-dot range in a later iteration; not required for v1 unless trivial to add as `--two-dot`.

**Worktree semantics (v1):**

- `--repo` / `-C`: primary repository context (any worktree path in the linked set).
- `--base-worktree` / `--head-worktree`: paths to linked worktrees; default ref for each side is that worktree’s `HEAD` unless `--base` / `--head` override that side.
- When only `--base` / `--head` are set, behave as a normal in-repo ref diff at `--repo`.
- Validate both worktree flags reference the same `git-common-dir` as `--repo` when any worktree flag is used.

**Lazy loading:** On session open, call `git diff --name-status` (with rename detection) only. On file selection, `git diff` for that path only. Optionally prefetch adjacent files in a later release; not v1.

**Non-text files:** Include in the change list with status; in the main pane show a short stub (e.g. “Binary file changed”) instead of a side-by-side text view.

**Keyboard defaults:** `j`/`k` or `]`/`[` next/prev file; `n`/`p` next/prev hunk; `q` quit. Optional `?` overlay in GUI mirroring help.

**CLI contract (v1):**

```
review-diff -C <repo> --base <ref> --head <ref>
review-diff -C <repo> --base-worktree <path> --head-worktree <path> [--base <ref>] [--head <ref>]
```

**Distribution:** Wails `build` → macOS `.app`, Linux tarball containing single binary. Repository includes `PKGBUILD` (Arch) that downloads the release tarball and verifies checksums—no separate Arch binary artifact.

**v2 hook:** Design session and Git runner so `--fetch` can run `git fetch` (scoped remotes or `--all`) before `ListChanges` without UI changes.

### Technical stack

- **Wails v2**, **Go** (current stable), **Svelte** + TypeScript for frontend.
- **Git** assumed on `PATH`; minimum version TBD in implementation (document in README; target 2.30+).
- No Electron; no syntax highlighter in v1.

### Prototype note

No prototype exists yet; module boundaries above are the initial architecture contract.

## Testing Decisions

Prefer **behavioral tests** at module boundaries, not Wails/webview integration, for v1.

| Module | Test approach |
|--------|----------------|
| **Diff parser** | Unit tests with fixture unified-diff strings; assert hunk count, line numbers, added/removed lines |
| **Git runner / change list / patch loader** | Integration tests in temporary repos created in test setup: branches, renames, binary file, two linked worktrees |
| **Worktree validator** | Integration: valid pair passes; unrelated paths fail with clear error |
| **Review session / CLI validation** | Table-driven unit tests for flag combinations and error messages (no GUI) |
| **Svelte UI** | Manual QA for v1; defer component/e2e tests unless a lightweight harness is already standard |

**Prior art:** Greenfield repo; no existing test patterns.

**Good tests:** Assert file list order and statuses, three-dot vs accidental two-dot behavior via known merge-base fixtures, lazy load returns patch only for requested path, worktree `HEAD` resolution matches `git rev-parse` in each worktree.

**Not worth testing in v1:** Pixel layout, system theme detection, Wails window lifecycle.

## Out of Scope (v1)

- Pull request URLs and GitHub/GitLab (or other) APIs
- Posting review comments or approvals
- `git fetch` / network (deferred to v2 `--fetch`)
- Syntax highlighting (e.g. Shiki)
- Terminal UI (TUI)
- Review queue or multi-PR dashboard
- Cross-repository diff (unrelated repos)
- Binary/image preview, submodule patch rendering
- Configurable keymaps or persisted UI preferences beyond system theme
- Windows support
- `git review-diff`-style `git-*` alias packaging (binary name is `review-diff` only unless user aliases manually)

## Further Notes

- **Origin requirement:** [requirement.md](requirement.md) called out PR efficiency, remote branches, and workspaces; v1 delivers **local ref + worktree** review; **remote-tracking ref names** without fetch; **PR-style three-dot** diff. PR URL workflow remains a later epic.
- **v2 priority (agreed):** Optional `--fetch` before diff; then consider PR URL open and syntax highlighting.
- **Issue tracker:** [GitHub Issues](https://github.com/dhannanjay/git-review-system/issues) (#1–#9); local copies in [issues/](issues/README.md).
- **Modules / tests confirmation:** Suggested modules and test focus are listed above; adjust before implementation if you want fewer integration tests or an earlier `--two-dot` flag.
