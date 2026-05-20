package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"review-diff/session"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	fs := flag.NewFlagSet("review-diff", flag.ContinueOnError)
	help := fs.Bool("help", false, "show usage")
	repo := fs.String("C", "", "path to git repository")
	base := fs.String("base", "", "base ref (e.g. main, origin/main)")
	head := fs.String("head", "", "head ref (e.g. feature-branch)")
	baseWorktree := fs.String("base-worktree", "", "base worktree path (resolves to its HEAD unless --base overrides)")
	headWorktree := fs.String("head-worktree", "", "head worktree path (resolves to its HEAD unless --head overrides)")
	fs.Usage = func() { fmt.Fprint(os.Stderr, usageText()) }

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}

	if *help {
		fmt.Print(usageText())
		os.Exit(0)
	}

	// If no flags given (no -C, no --base, no --head, no worktree flags), show usage
	// This also handles the case where the binary is run for wails module generation
	// which provides no arguments - in that case we start the app without a session
	// so bindings can be registered.
	hasFlags := *repo != "" || *base != "" || *head != "" || *baseWorktree != "" || *headWorktree != ""
	if hasFlags {
		sess, err := session.New(*repo, *base, *head, *baseWorktree, *headWorktree)
		if err != nil {
			fmt.Fprintf(os.Stderr, "review-diff: %v\n", err)
			fs.Usage()
			os.Exit(2)
		}
		if err := sess.ResolveWorktreeRefs(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "review-diff: %v\n", err)
			os.Exit(2)
		}
		if err := validateRefs(context.Background(), sess); err != nil {
			fmt.Fprintf(os.Stderr, "review-diff: %v\n", err)
			os.Exit(2)
		}
		app := NewApp(sess)
		run(app)
	} else {
		// Start GUI without flags (for wails generate module or direct launch)
		app := NewApp(nil)
		run(app)
	}
}

func run(app *App) {
	if err := wails.Run(&options.App{
		Title:  "review-diff",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func usageText() string {
	return `Usage: review-diff -C <repo> --base <ref> --head <ref>
       review-diff -C <repo> --base-worktree <path> --head-worktree <path> [--base <ref>] [--head <ref>]

A native desktop diff viewer for local Git branches.
Opens a GUI window that blocks until closed.

Flags:
  -C <repo>             path to git repository
  --base <ref>          base ref (e.g. main, origin/main)
  --head <ref>          head ref (e.g. feature-branch)
  --base-worktree <path> base worktree path (resolves to its HEAD unless --base overrides)
  --head-worktree <path> head worktree path (resolves to its HEAD unless --head overrides)
  --help                show this message and exit

Keyboard:
  j / ]      next file
  k / [      previous file
  n          next hunk
  p          previous hunk
  q          quit
`
}
