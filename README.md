# review-diff

A native desktop diff viewer for local Git branches.

## Development

```bash
# Run in live development mode (hot reloads frontend on changes)
wails dev

# Build a production binary
wails build

# Run tests
go test ./...
```

## Usage

```bash
# Launch the GUI (blocks until window closes)
review-diff

# Show help
review-diff --help
```

Branch-diff flags (`--base`, `--head`, etc.) are coming in a later release.

## Requirements

- Go 1.23+
- Node.js 20+
- Wails v2 CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
