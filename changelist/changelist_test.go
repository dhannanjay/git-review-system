package changelist

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"review-diff/gitrunner"
)

func TestListChanges(t *testing.T) {
	dir, err := os.MkdirTemp("", "changelist-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}

	git("init", "-b", "main")
	git("config", "user.email", "t@t.com")
	git("config", "user.name", "T")
	os.WriteFile(dir+"/a.txt", []byte("a\n"), 0644)
	os.WriteFile(dir+"/b.txt", []byte("b\n"), 0644)
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	os.WriteFile(dir+"/a.txt", []byte("a\nmodified\n"), 0644)
	os.WriteFile(dir+"/c.txt", []byte("new\n"), 0644)
	git("add", ".")
	git("commit", "-m", "main-change")

	git("checkout", "feature")
	os.WriteFile(dir+"/a.txt", []byte("a\nfeature\n"), 0644)
	git("add", ".")
	git("commit", "-m", "feature-change")

	git("checkout", "main")

	runner := gitrunner.New(dir)
	changes, err := ListChanges(context.Background(), runner, "main", "feature")
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 {
		t.Fatal("expected changes")
	}
	// three-dot between main and feature (merge-base at init)
	// feature only changed a.txt, so only a.txt should appear
	found := false
	for _, c := range changes {
		if c.NewPath == "a.txt" {
			found = true
		}
		if c.NewPath == "b.txt" {
			t.Error("b.txt should not be in three-dot diff (unchanged in feature)")
		}
		if c.NewPath == "c.txt" {
			t.Error("c.txt should not be in three-dot diff (only in main)")
		}
	}
	if !found {
		t.Error("expected a.txt in changes")
	}
}

func TestListChangesRename(t *testing.T) {
	dir, err := os.MkdirTemp("", "changelist-rename-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}

	git("init", "-b", "main")
	git("config", "user.email", "t@t.com")
	git("config", "user.name", "T")
	os.WriteFile(dir+"/old.txt", []byte("content\n"), 0644)
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	git("checkout", "feature")
	git("mv", "old.txt", "new.txt")
	git("commit", "-m", "rename-file")

	git("checkout", "main")

	runner := gitrunner.New(dir)
	changes, err := ListChanges(context.Background(), runner, "main", "feature")
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 {
		t.Fatal("expected changes")
	}
	var found bool
	for _, c := range changes {
		if c.Status == "R" {
			found = true
			if c.OldPath != "old.txt" {
				t.Errorf("old path: got %q, want %q", c.OldPath, "old.txt")
			}
			if c.NewPath != "new.txt" {
				t.Errorf("new path: got %q, want %q", c.NewPath, "new.txt")
			}
		}
	}
	if !found {
		t.Error("expected rename in changes")
	}
}

func TestListChangesDelete(t *testing.T) {
	dir, err := os.MkdirTemp("", "changelist-delete-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}

	git("init", "-b", "main")
	git("config", "user.email", "t@t.com")
	git("config", "user.name", "T")
	os.WriteFile(dir+"/keep.txt", []byte("keep\n"), 0644)
	os.WriteFile(dir+"/delete-me.txt", []byte("bye\n"), 0644)
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	git("checkout", "feature")
	git("rm", "delete-me.txt")
	git("commit", "-m", "delete-file")

	git("checkout", "main")

	runner := gitrunner.New(dir)
	changes, err := ListChanges(context.Background(), runner, "main", "feature")
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 {
		t.Fatal("expected changes")
	}
	var found bool
	for _, c := range changes {
		if c.Status == "D" && c.NewPath == "delete-me.txt" {
			found = true
		}
	}
	if !found {
		t.Error("expected delete-me.txt with D status in changes")
	}
}
