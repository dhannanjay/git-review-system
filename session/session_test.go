package session

import "testing"

func TestNew_valid(t *testing.T) {
	s, err := New("/repo", "main", "feature", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if s.Repo != "/repo" {
		t.Errorf("repo: got %q", s.Repo)
	}
	if s.Base != "main" {
		t.Errorf("base: got %q", s.Base)
	}
	if s.Head != "feature" {
		t.Errorf("head: got %q", s.Head)
	}
}

func TestNew_missingRepo(t *testing.T) {
	_, err := New("", "main", "feature", "", "")
	if err == nil {
		t.Fatal("expected error for missing repo")
	}
}

func TestNew_missingBase(t *testing.T) {
	_, err := New("/repo", "", "feature", "", "")
	if err == nil {
		t.Fatal("expected error for missing base")
	}
}

func TestNew_missingHead(t *testing.T) {
	_, err := New("/repo", "main", "", "", "")
	if err == nil {
		t.Fatal("expected error for missing head")
	}
}

func TestNew_worktreeValid(t *testing.T) {
	s, err := New("/repo", "", "", "/wt/base", "/wt/head")
	if err != nil {
		t.Fatal(err)
	}
	if s.BaseWorktree != "/wt/base" {
		t.Errorf("baseWorktree: got %q", s.BaseWorktree)
	}
	if s.HeadWorktree != "/wt/head" {
		t.Errorf("headWorktree: got %q", s.HeadWorktree)
	}
}

func TestNew_worktreeWithRef(t *testing.T) {
	s, err := New("/repo", "main", "feature", "/wt/base", "/wt/head")
	if err != nil {
		t.Fatal(err)
	}
	if s.Base != "main" {
		t.Errorf("base: got %q", s.Base)
	}
	if s.Head != "feature" {
		t.Errorf("head: got %q", s.Head)
	}
	if s.BaseWorktree != "/wt/base" {
		t.Errorf("baseWorktree: got %q", s.BaseWorktree)
	}
	if s.HeadWorktree != "/wt/head" {
		t.Errorf("headWorktree: got %q", s.HeadWorktree)
	}
}

func TestNew_missingBaseAndWorktree(t *testing.T) {
	_, err := New("/repo", "", "", "", "")
	if err == nil {
		t.Fatal("expected error for missing both --base and --base-worktree")
	}
}

func TestNew_missingHeadAndWorktree(t *testing.T) {
	_, err := New("/repo", "main", "", "", "")
	if err == nil {
		t.Fatal("expected error for missing both --head and --head-worktree")
	}
}
