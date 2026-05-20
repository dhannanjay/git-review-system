package changelist

import (
	"context"
	"fmt"
	"strings"

	"review-diff/gitrunner"
)

type FileChange struct {
	OldPath string
	NewPath string
	Status  string // M, A, D, R (rename score stripped)
}

func ListChanges(ctx context.Context, runner *gitrunner.Runner, base, head string) ([]FileChange, error) {
	out, err := runner.Run(ctx, "diff", "--name-status", "--find-renames", base+"..."+head)
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
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 2 {
			continue
		}
		status := parts[0]
		if len(status) > 0 && status[0] == 'R' {
			status = "R"
		}
		fc := FileChange{Status: status}
		if status == "R" && len(parts) >= 3 {
			fc.OldPath = parts[1]
			fc.NewPath = parts[2]
		} else {
			fc.NewPath = parts[1]
		}
		changes = append(changes, fc)
	}
	return changes, nil
}
