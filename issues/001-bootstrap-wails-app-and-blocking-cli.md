# Issue 1: Bootstrap Wails v2 app and blocking CLI

## Parent

[PRD.md](../PRD.md) — review-diff v1

## What to build

Establish the project skeleton: Go module, **Wails v2** + **Svelte** frontend, single binary **`review-diff`**. Running the command opens a native window (initially empty or placeholder content) and the **CLI blocks until the window closes**, then exits with a predictable code. Ship minimal **`--help`** describing the tool and forthcoming flags.

This slice is the foundation: conventions for Wails bindings, frontend build, and CLI entrypoint are decided here.

## Acceptance criteria

- [ ] `review-diff` (no args or with placeholder flags) builds and runs on **Linux and macOS** dev machines
- [ ] A Wails window opens and can be closed; the shell command does not return until the GUI exits
- [ ] `review-diff --help` prints usage (binary name, brief description, note that branch diff flags come in later slices)
- [ ] README or CONTRIBUTING notes how to run dev (`wails dev`) and build (`wails build`)
- [ ] Repository layout is suitable for subsequent slices (Go backend package, Svelte app, one main)

## Blocked by

None — can start immediately.
