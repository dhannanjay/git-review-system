#!/usr/bin/env bash
set -euo pipefail

VERSION="${VERSION:-0.1.0}"
BUILD_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$BUILD_DIR/.." && pwd)"
OUTDIR="$BUILD_DIR/release"

usage() {
  cat <<EOF
Usage: VERSION=x.y.z ./build/release.sh

Builds a distributable release artifact for the current platform.

Environment:
  VERSION   Version string (default: 0.1.0)

Output:
  macOS:   $OUTDIR/review-diff-\$VERSION-darwin-amd64.zip  (.app bundle inside)
  Linux:   $OUTDIR/review-diff-\$VERSION-linux-amd64.tar.gz (single binary)
EOF
  exit 0
}

[[ "${1:-}" == "--help" || "${1:-}" == "-h" ]] && usage

echo "==> review-diff release v$VERSION ($(uname -s))"

cd "$REPO_ROOT"

# Build frontend
echo "==> Building frontend…"
npm --prefix frontend run build

# Build Go binary via Wails
echo "==> Running wails build…"
wails build -clean -tags webkit2_41

case "$(uname -s)" in
  Darwin)
    mkdir -p "$OUTDIR"
    APP_SRC="$BUILD_DIR/bin/review-diff.app"
    if [ ! -d "$APP_SRC" ]; then
      echo "!! Expected .app at $APP_SRC — check wails build output"
      exit 1
    fi
    ARCHIVE="$OUTDIR/review-diff-$VERSION-darwin-amd64.zip"
    echo "==> Packaging macOS .app → $ARCHIVE"
    cd "$BUILD_DIR/bin"
    zip -r "$ARCHIVE" "review-diff.app"
    cd "$REPO_ROOT"
    echo "==> Done: $ARCHIVE"
    ;;

  Linux)
    mkdir -p "$OUTDIR"
    BIN_SRC="$BUILD_DIR/bin/review-diff"
    if [ ! -f "$BIN_SRC" ]; then
      echo "!! Expected binary at $BIN_SRC — check wails build output"
      exit 1
    fi
    ARCHIVE="$OUTDIR/review-diff-$VERSION-linux-amd64.tar.gz"
    echo "==> Packaging Linux binary → $ARCHIVE"
    tar -C "$BUILD_DIR/bin" -czf "$ARCHIVE" "review-diff"
    echo "==> Done: $ARCHIVE"
    ;;

  *)
    echo "Unsupported platform: $(uname -s)"
    exit 1
    ;;
esac
