# TASK

Fix issue {{TASK_ID}}: {{ISSUE_TITLE}}

Pull in the issue using `gh issue view {{TASK_ID}}`. Read [PRD.md](PRD.md) for product context.

Only work on the issue specified.

Work on branch {{BRANCH}}. **You must `git add` and `git commit`** — Sandcastle only sees commits, not unstaged files.

# STACK (review-diff)

- **Go** backend + **Wails v2** desktop shell + **Svelte** frontend
- Binary name: `review-diff`
- Git subprocesses for diffing (no GitHub API in v1)
- Do **not** use Electron

If there is no `go.mod` yet, scaffold with Wails v2 (Svelte template) in the repo root.

# CONTEXT

Here are the last 10 commits:

<recent-commits>

!`git log -n 10 --format="%H%n%ad%n%B---" --date=short`

</recent-commits>

# EXPLORATION

Explore the repo. The worktree may start as docs-only; your job is to add the application code required by the issue.

# EXECUTION

1. Implement the smallest vertical slice that satisfies the issue acceptance criteria.
2. Prefer `go test ./...` for tests; add packages under test as you add code.
3. Verify the project builds:
   - `go build -o /dev/null ./...` (or `go build ./...`)
   - `wails build` once a Wails app exists (skip only if Wails is not yet scaffolded)

# FEEDBACK LOOPS

Before each commit, run the checks above. Fix failures before committing.

Do **not** rely on `npm run test` unless `package.json` defines those scripts.

# COMMIT

**Every iteration must end with at least one commit** if you changed files.

Commit message rules:

1. Start with `RALPH:` prefix
2. Reference issue {{TASK_ID}} and PRD
3. Brief summary of what was done

Example: `RALPH: #1 bootstrap Wails v2 app and blocking CLI`

# THE ISSUE

If the task is not fully complete, leave a comment on the issue with progress.

Do not close the issue — the merge agent does that.

When the issue acceptance criteria are met, output:

<promise>COMPLETE</promise>

# FINAL RULES

ONLY WORK ON A SINGLE TASK.
