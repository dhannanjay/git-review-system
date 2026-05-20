package session

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"review-diff/gitrunner"
)

// resolveCommonDir converts the output of git rev-parse --git-common-dir to an
// absolute path. The output may be relative (e.g. ".git") or absolute (for
// linked worktrees, the path points into the main repo's .git/worktrees/ dir).
func resolveCommonDir(repoPath, commonDir string) (string, error) {
	if filepath.IsAbs(commonDir) {
		return filepath.Clean(commonDir), nil
	}
	return filepath.Abs(filepath.Join(repoPath, commonDir))
}

type Session struct {
	Repo         string
	Base         string
	Head         string
	BaseWorktree string
	HeadWorktree string
}

func New(repo, base, head, baseWorktree, headWorktree string) (*Session, error) {
	if repo == "" {
		return nil, fmt.Errorf("repo path (-C) is required")
	}
	if base == "" && baseWorktree == "" {
		return nil, fmt.Errorf("either --base or --base-worktree is required")
	}
	if head == "" && headWorktree == "" {
		return nil, fmt.Errorf("either --head or --head-worktree is required")
	}
	return &Session{
		Repo:         repo,
		Base:         base,
		Head:         head,
		BaseWorktree: baseWorktree,
		HeadWorktree: headWorktree,
	}, nil
}

func (s *Session) ResolveWorktreeRefs(ctx context.Context) error {
	if s.BaseWorktree == "" && s.HeadWorktree == "" {
		return nil
	}

	runner := gitrunner.New(s.Repo)
	mainCommonOut, err := runner.Run(ctx, "rev-parse", "--git-common-dir")
	if err != nil {
		return fmt.Errorf("resolve common dir for %s: %w", s.Repo, err)
	}
	mainCommon := strings.TrimSpace(mainCommonOut)
	mainCommonAbs, err := resolveCommonDir(s.Repo, mainCommon)
	if err != nil {
		return err
	}

	checkLinked := func(wtPath string) error {
		wtRunner := gitrunner.New(wtPath)
		wtCommonOut, err := wtRunner.Run(ctx, "rev-parse", "--git-common-dir")
		if err != nil {
			return fmt.Errorf("worktree %q is not a git repository", wtPath)
		}
		wtCommon := strings.TrimSpace(wtCommonOut)
		wtCommonAbs, err := resolveCommonDir(wtPath, wtCommon)
		if err != nil {
			return err
		}
		if mainCommonAbs != wtCommonAbs {
			return fmt.Errorf("worktree %q is not linked to repository %s", wtPath, s.Repo)
		}
		return nil
	}

	resolveHEAD := func(wtPath string) (string, error) {
		wtRunner := gitrunner.New(wtPath)
		out, err := wtRunner.Run(ctx, "rev-parse", "HEAD")
		if err != nil {
			return "", fmt.Errorf("resolve HEAD in %q: %w", wtPath, err)
		}
		return strings.TrimSpace(out), nil
	}

	if s.BaseWorktree != "" {
		if err := checkLinked(s.BaseWorktree); err != nil {
			return err
		}
		if s.Base == "" {
			headRef, err := resolveHEAD(s.BaseWorktree)
			if err != nil {
				return err
			}
			s.Base = headRef
		}
	}

	if s.HeadWorktree != "" {
		if err := checkLinked(s.HeadWorktree); err != nil {
			return err
		}
		if s.Head == "" {
			headRef, err := resolveHEAD(s.HeadWorktree)
			if err != nil {
				return err
			}
			s.Head = headRef
		}
	}

	return nil
}
