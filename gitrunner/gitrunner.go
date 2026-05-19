package gitrunner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type Runner struct {
	repo string
}

func New(repo string) *Runner {
	return &Runner{repo: repo}
}

func (r *Runner) Run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = r.repo
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %v: %s\n%s", args, stderr.String(), err)
	}
	return stdout.String(), nil
}
