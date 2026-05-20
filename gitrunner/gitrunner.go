package gitrunner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

var allowedGitSubcommands = map[string]bool{
	"diff":       true,
	"log":        true,
	"status":     true,
	"merge-base": true,
	"rev-parse":  true,
}

type Runner struct {
	repo string
}

func New(repo string) *Runner {
	return &Runner{repo: repo}
}

func (r *Runner) Run(ctx context.Context, args ...string) (string, error) {
	if err := validateGitArgs(args); err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, "git", args...) // #nosec G204 — args validated by validateGitArgs
	cmd.Dir = r.repo
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %v: %s\n%s", args, stderr.String(), err)
	}
	return stdout.String(), nil
}

func validateGitArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("git: no arguments provided")
	}
	if !allowedGitSubcommands[args[0]] {
		return fmt.Errorf("git: disallowed subcommand: %s", args[0])
	}
	for _, arg := range args {
		if arg == "-c" || arg == "--config-env" {
			return fmt.Errorf("git: disallowed option: %s", arg)
		}
	}
	return nil
}
