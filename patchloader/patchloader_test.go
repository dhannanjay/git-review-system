package patchloader

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"review-diff/gitrunner"
)

func TestLoadPatchBinary(t *testing.T) {
	dir, err := os.MkdirTemp("", "patchloader-binary-*")
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
	// Create a file with NUL byte to make git detect it as binary
	os.WriteFile(dir+"/file.bin", []byte("hello\x00world\n"), 0644)
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	git("checkout", "feature")
	// Modify the binary file
	os.WriteFile(dir+"/file.bin", []byte("modified\x00content\n"), 0644)
	git("add", ".")
	git("commit", "-m", "modify-binary")

	git("checkout", "main")

	runner := gitrunner.New(dir)
	df, err := LoadPatch(context.Background(), runner, "main", "feature", "file.bin")
	if err != nil {
		t.Fatal(err)
	}
	if !df.IsBinary {
		t.Error("expected IsBinary to be true for binary file")
	}
	if df.OldPath != "file.bin" && df.OldPath != "" {
		t.Logf("old path: %q", df.OldPath)
	}
	if df.NewPath != "file.bin" && df.NewPath != "" {
		t.Logf("new path: %q", df.NewPath)
	}
}

func TestLoadPatchBinaryNoNull(t *testing.T) {
	// Even without NUL, if git treats output as binary we should detect it.
	// Use .gitattributes to force binary.
	dir, err := os.MkdirTemp("", "patchloader-binary2-*")
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

	os.WriteFile(dir+"/data.bin", []byte("some content\n"), 0644)
	os.WriteFile(dir+"/.gitattributes", []byte("*.bin binary\n"), 0644)
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	git("checkout", "feature")
	os.WriteFile(dir+"/data.bin", []byte("modified content\n"), 0644)
	git("add", ".")
	git("commit", "-m", "modify-binary")

	git("checkout", "main")

	runner := gitrunner.New(dir)
	df, err := LoadPatch(context.Background(), runner, "main", "feature", "data.bin")
	if err != nil {
		t.Fatal(err)
	}
	if !df.IsBinary {
		t.Error("expected IsBinary to be true for .gitattributes-forced binary file")
	}
}
