package gitrunner

import (
	"context"
	"os"
	"os/exec"
	"testing"
)

func TestRun_success(t *testing.T) {
	dir, err := os.MkdirTemp("", "gitrunner-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	exec.Command("git", "-C", dir, "init", "-b", "main").Run()
	exec.Command("git", "-C", dir, "config", "user.email", "t@t.com").Run()
	exec.Command("git", "-C", dir, "config", "user.name", "T").Run()
	os.WriteFile(dir+"/a.txt", []byte("hello"), 0644)
	exec.Command("git", "-C", dir, "add", ".").Run()
	exec.Command("git", "-C", dir, "commit", "-m", "init").Run()

	r := New(dir)
	out, err := r.Run(context.Background(), "log", "--oneline")
	if err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Error("expected git log output")
	}
}

func TestRun_error(t *testing.T) {
	dir, err := os.MkdirTemp("", "gitrunner-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	r := New(dir)
	_, err = r.Run(context.Background(), "status")
	if err == nil {
		t.Error("expected error for non-git directory")
	}
}
