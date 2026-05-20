package patchloader

import (
	"context"
	"fmt"
	"strings"

	"review-diff/diffparser"
	"review-diff/gitrunner"
)

func LoadPatch(ctx context.Context, runner *gitrunner.Runner, base, head, filePath string) (*diffparser.DiffFile, error) {
	out, err := runner.Run(ctx, "diff", base+"..."+head, "--", filePath)
	if err != nil {
		return nil, fmt.Errorf("load patch: %w", err)
	}
	if out == "" {
		return nil, fmt.Errorf("no diff for %s", filePath)
	}

	if strings.Contains(out, "Binary files") {
		df := &diffparser.DiffFile{IsBinary: true}
		for _, line := range strings.Split(out, "\n") {
			if strings.HasPrefix(line, "diff --git ") {
				fields := strings.Fields(line)
				if len(fields) >= 4 {
					df.OldPath = strings.TrimPrefix(fields[2], "a/")
					df.NewPath = strings.TrimPrefix(fields[3], "b/")
				}
				break
			}
		}
		return df, nil
	}

	df, err := diffparser.ParseUnified(out)
	if err != nil {
		return nil, err
	}

	if strings.Contains(out, "Subproject commit") {
		df.IsSubmodule = true
	}

	return df, nil
}
