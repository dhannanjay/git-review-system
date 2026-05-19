---
name: to-prd
description: Turn the current conversation context into a PRD and publish it to the project issue tracker. Use when the user wants to create a PRD, write a product requirements doc, or mentions "to prd" or "/to-prd".
---

# To PRD

Synthesize the current conversation and codebase understanding into a PRD. Do **not** interview the user — use what you already know.

## Issue tracker

Before publishing, determine where the PRD should live:

1. If `docs/agents/issue-tracker.md` exists in the project repo, follow it.
2. Else if the repo uses GitHub and `gh` is available, create a GitHub issue.
3. Else ask the user (issue tracker URL, tool, or paste-only output).

If `docs/agents/triage-labels.md` defines a label for work ready for autonomous implementation, apply it when publishing. Otherwise ask or omit labels.

## Process

1. Explore the repo to understand the current codebase, if you have not already. Use the project's domain vocabulary (e.g. `CONTEXT.md`, glossaries). Respect ADRs in the area you are touching.

2. Sketch the major modules to build or modify. Look for **deep modules** — lots of functionality behind a simple, stable, testable interface.

   Confirm with the user that these modules match their expectations and which modules should have tests.

3. Write the PRD using the template below, then publish to the issue tracker (or deliver the markdown to the user if no tracker is configured).

<prd-template>

## Problem Statement

The problem the user faces, from the user's perspective.

## Solution

The solution, from the user's perspective.

## User Stories

A long, numbered list. Each story:

1. As an <actor>, I want a <feature>, so that <benefit>

Cover all aspects of the feature extensively.

<user-story-example>
1. As a mobile bank customer, I want to see balance on my accounts, so that I can make better informed decisions about my spending
</user-story-example>

## Implementation Decisions

Decisions made, including:

- Modules to build or modify
- Interfaces that will change
- Technical clarifications
- Architecture, schema, API contracts
- Specific interactions

Do **not** include file paths or code snippets unless a prototype encodes a decision more precisely than prose (state machine, reducer, schema, type shape). In that case, inline only the decision-rich parts and note it came from a prototype.

## Testing Decisions

- What makes a good test here (prefer external behavior over implementation details)
- Which modules will be tested
- Prior art for similar tests in the codebase

## Out of Scope

What this PRD explicitly does not cover.

## Further Notes

Any other notes.

</prd-template>
