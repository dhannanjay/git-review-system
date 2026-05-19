# Issue 8: Release builds and Arch PKGBUILD

## Parent

[PRD.md](../PRD.md) — review-diff v1

## What to build

Shippable artifacts and install docs:

- **macOS:** Wails build produces a distributable **`.app`**
- **Linux:** release **tarball** containing a single `review-diff` binary
- **Arch:** **`PKGBUILD`** in repo that downloads the official Linux tarball and verifies checksums (no separate Arch binary build)
- README sections: install from release, build from source, Git ≥2.30 on PATH

CI/release automation is optional in this slice; documented manual release steps are acceptable for v1.

## Acceptance criteria

- [ ] Documented steps produce working `.app` on macOS and binary tarball on Linux
- [ ] Tarball binary runs on a clean machine with git installed
- [ ] `PKGBUILD` builds/installs against published tarball (document version bump process)
- [ ] README covers Linux, macOS, and Arch (AUR) install paths

## Blocked by

- #2 E2E two-ref diff: first changed file side-by-side

## User stories

18, 19, 20
