package changelist

import (
	"context"
	"fmt"
	"strings"

	"review-diff/gitrunner"
)

type FileChange struct {
	Path   string
	Status string // M, A, D, R, etc.
}

func ListChanges(ctx context.Context, runner *gitrunner.Runner, base, head string) ([]FileChange, error) {
	out, err := runner.Run(ctx, "diff", "--name-status", base+"..."+head)
	if err != nil {
		return nil, fmt.Errorf("list changes: %w", err)
	}
	out = strings.TrimSpace(out)
	if out == "" {
		return nil, nil
	}
	lines := strings.Split(out, "\n")
	var changes []FileChange
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) < 2 {
			continue
		}
		changes = append(changes, FileChange{Status: parts[0], Path: parts[1]})
	}
	return changes, nil
}
