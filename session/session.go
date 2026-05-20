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

	mainCommonAbs, err := s.resolveMainRepoCommonDir(ctx)
	if err != nil {
		return err
	}

	if s.BaseWorktree != "" {
		if err := s.checkLinked(ctx, mainCommonAbs, s.BaseWorktree); err != nil {
			return err
		}
		if s.Base == "" {
			commit, err := s.resolveHEAD(ctx, s.BaseWorktree)
			if err != nil {
				return err
			}
			s.Base = commit
		}
	}

	if s.HeadWorktree != "" {
		if err := s.checkLinked(ctx, mainCommonAbs, s.HeadWorktree); err != nil {
			return err
		}
		if s.Head == "" {
			commit, err := s.resolveHEAD(ctx, s.HeadWorktree)
			if err != nil {
				return err
			}
			s.Head = commit
		}
	}

	return nil
}

func (s *Session) resolveMainRepoCommonDir(ctx context.Context) (string, error) {
	runner := gitrunner.New(s.Repo)
	out, err := runner.Run(ctx, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("resolve common dir for %s: %w", s.Repo, err)
	}
	return resolveCommonDir(s.Repo, strings.TrimSpace(out))
}

func (s *Session) checkLinked(ctx context.Context, mainCommonAbs, wtPath string) error {
	runner := gitrunner.New(wtPath)
	out, err := runner.Run(ctx, "rev-parse", "--git-common-dir")
	if err != nil {
		return fmt.Errorf("worktree %q is not a git repository", wtPath)
	}
	wtCommonAbs, err := resolveCommonDir(wtPath, strings.TrimSpace(out))
	if err != nil {
		return err
	}
	if mainCommonAbs != wtCommonAbs {
		return fmt.Errorf("worktree %q is not linked to repository %s", wtPath, s.Repo)
	}
	return nil
}

func (s *Session) resolveHEAD(ctx context.Context, wtPath string) (string, error) {
	runner := gitrunner.New(wtPath)
	out, err := runner.Run(ctx, "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("resolve HEAD in %q: %w", wtPath, err)
	}
	return strings.TrimSpace(out), nil
}
