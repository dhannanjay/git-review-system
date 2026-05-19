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
		if c.Path == "a.txt" {
			found = true
		}
		if c.Path == "b.txt" {
			t.Error("b.txt should not be in three-dot diff (unchanged in feature)")
		}
		if c.Path == "c.txt" {
			t.Error("c.txt should not be in three-dot diff (only in main)")
		}
	}
	if !found {
		t.Error("expected a.txt in changes")
	}
}
