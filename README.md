# review-diff

A native desktop diff viewer for local Git branches.

## Requirements

- **Git ≥ 2.30** (required at runtime; must be on `PATH`)
- **Go 1.23+** (to build from source)
- **Node.js 20+** (to build from source)
- **Wails v2 CLI** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`) (to build from source)

## Install

### macOS — download the release `.app`

1. Download `review-diff-<version>-darwin-amd64.zip` from the
   [releases page](https://github.com/dhannanjay/git-review-system/releases).
2. Unzip and move `review-diff.app` to `/Applications`.
3. Run from Spotlight or from the terminal:

   ```bash
   open /Applications/review-diff.app
   ```

### Linux — release tarball

1. Download `review-diff-<version>-linux-amd64.tar.gz` from the
   [releases page](https://github.com/dhannanjay/git-review-system/releases).
2. Extract and place the binary on your `PATH`:

   ```bash
   tar xzf review-diff-<version>-linux-amd64.tar.gz
   sudo mv review-diff /usr/local/bin/
   ```

### Arch Linux — AUR / PKGBUILD

The repository includes a [`PKGBUILD`](PKGBUILD) that wraps the official
Linux tarball.  To install:

```bash
git clone https://github.com/dhannanjay/git-review-system
cd git-review-system
makepkg -fi
```

Or publish to the AUR following the usual workflow.  See `PKGBUILD` for
the version bump process.

## Build from source

```bash
# Install the Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Run in live development mode (hot reloads frontend on changes)
wails dev

# Build a production binary
wails build

# Run tests
go test ./...
```

The compiled binary and/or `.app` bundle are written to `build/bin/`.
To produce a distributable archive, run:

```bash
VERSION=x.y.z ./build/release.sh
```

## Usage

```bash
# Launch the GUI (blocks until window closes)
review-diff

# Show help
review-diff --help
```

## Keyboard

| Key | Action |
|-----|--------|
| `j` / `]` | Next file |
| `k` / `[` | Previous file |
| `n` | Next hunk |
| `p` | Previous hunk |
| `q` | Quit |
