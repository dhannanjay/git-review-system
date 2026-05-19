# Issue (v2): Optional `--fetch` before diff

## Parent

[PRD.md](../PRD.md) — review-diff v2 (deferred)

## What to build

Add optional **`--fetch`** to refresh remote-tracking refs before building the change list and diffs. Run `git fetch` (scoped remotes or documented default) inside the git runner; **no UI changes** if session/list API stays the same. Network use is explicit and opt-in.

Design should match the v1 session/git-runner split so fetch is a pre-step before `ListChanges`.

## Acceptance criteria

- [ ] `review-diff ... --fetch` updates refs then shows the same GUI diff as without fetch when network succeeds
- [ ] Without `--fetch`, behavior unchanged from v1 (no network)
- [ ] Fetch failure surfaces clear CLI error; does not hang GUI
- [ ] Integration test with a local bare remote or mock documents expected behavior

## Blocked by

- v1 complete (through #8)

## User stories

22
