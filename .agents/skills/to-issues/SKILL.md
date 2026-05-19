---
name: to-issues
description: Break a plan, spec, or PRD into independently-grabbable issues on the project issue tracker using vertical-slice tracer bullets. Use when converting a PRD to issues, creating implementation tickets, or mentions "prd to issues", "to issues", or "/to-issues".
---

# To Issues (PRD → issues)

Break a plan into independently grabbable issues using **vertical slices** (tracer bullets).

## Issue tracker

Before publishing issues:

1. If `docs/agents/issue-tracker.md` exists in the project repo, follow it.
2. Else if the repo uses GitHub and `gh` is available, use GitHub Issues.
3. Else ask the user (tracker URL, tool, or markdown-only output).

If `docs/agents/triage-labels.md` defines labels for AFK-ready work, apply them when publishing slices marked AFK. Otherwise ask or omit labels.

## Process

### 1. Gather context

Use conversation context. If the user passes an issue reference (number, URL, or path), fetch the full body and comments from the issue tracker.

### 2. Explore the codebase (optional)

If needed, explore the codebase. Use domain vocabulary from project docs (e.g. `CONTEXT.md`). Respect ADRs in the affected area.

### 3. Draft vertical slices

Break the plan into **tracer bullet** issues. Each slice is a thin vertical cut through **all** integration layers end-to-end — not a single horizontal layer.

Slices are **HITL** (needs human input: architecture, design review) or **AFK** (can be implemented and merged without human interaction). Prefer AFK where possible.

<vertical-slice-rules>
- Each slice delivers a narrow but **complete** path through every layer (schema, API, UI, tests, etc.)
- A completed slice is demoable or verifiable on its own
- Prefer many thin slices over few thick ones
</vertical-slice-rules>

### 4. Quiz the user

Present the breakdown as a numbered list. For each slice:

- **Title**: short descriptive name
- **Type**: HITL / AFK
- **Blocked by**: which slices must complete first (if any)
- **User stories covered**: from the source material (if applicable)

Ask:

- Is the granularity right? (too coarse / too fine)
- Are dependencies correct?
- Should any slices be merged or split?
- Are HITL vs AFK labels correct?

Iterate until the user approves.

### 5. Publish issues

For each approved slice, create an issue using the template below. Publish in dependency order (blockers first) so real issue IDs can be referenced in **Blocked by**.

Do **not** close or modify the parent issue.

<issue-template>

## Parent

Reference to the parent issue (omit if there is no parent).

## What to build

Concise description of this vertical slice — end-to-end behavior, not layer-by-layer tasks.

Avoid file paths and code snippets unless a prototype encodes a decision (state machine, reducer, schema, type). Trim to decision-rich parts only.

## Acceptance criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## Blocked by

- Reference to blocking issue(s)

Or "None — can start immediately".

</issue-template>
