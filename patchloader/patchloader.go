package patchloader

import (
	"context"
	"fmt"

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
	return diffparser.ParseUnified(out)
}
