package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	fs := flag.NewFlagSet("review-diff", flag.ContinueOnError)
	help := fs.Bool("help", false, "show usage")
	fs.Usage = func() { fmt.Fprint(os.Stderr, usageText()) }

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}

	if *help {
		fmt.Print(usageText())
		os.Exit(0)
	}

	if fs.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "review-diff: unrecognized argument: %s\n", fs.Arg(0))
		fs.Usage()
		os.Exit(2)
	}

	app := NewApp()

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
	return `Usage: review-diff [--help]

A native desktop diff viewer for local Git branches.

Opens a GUI window that blocks until closed. Branch-diff flags
(--base, --head, etc.) will be added in a later release.

Flags:
  --help  show this message and exit
`
}
