package main

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"

	"review-diff/changelist"
	"review-diff/gitrunner"
	"review-diff/patchloader"
	"review-diff/session"
)

func TestUsageText(t *testing.T) {
	text := usageText()
	if !strings.Contains(text, "review-diff") {
		t.Error("usage text should contain binary name")
	}
	if !strings.Contains(text, "--help") {
		t.Error("usage text should mention --help")
	}
	if !strings.Contains(text, "--base") {
		t.Error("usage text should mention --base")
	}
	if !strings.Contains(text, "--head") {
		t.Error("usage text should mention --head")
	}
	if !strings.Contains(text, "-C") {
		t.Error("usage text should mention -C")
	}
	if !strings.Contains(text, "--base-worktree") {
		t.Error("usage text should mention --base-worktree")
	}
	if !strings.Contains(text, "--head-worktree") {
		t.Error("usage text should mention --head-worktree")
	}
	if !strings.Contains(text, "--fetch") {
		t.Error("usage text should mention --fetch")
	}
}

func TestUsageTextIncludesBlockingNote(t *testing.T) {
	text := usageText()
	if !strings.Contains(text, "blocks") {
		t.Error("usage text should mention blocking behavior")
	}
}

func TestUsageTextIncludesKeyboard(t *testing.T) {
	text := usageText()
	if !strings.Contains(text, "Keyboard") {
		t.Error("usage text should include keyboard section")
	}
	if !strings.Contains(text, "next file") {
		t.Error("usage text should document next file keybinding")
	}
	if !strings.Contains(text, "previous file") {
		t.Error("usage text should document previous file keybinding")
	}
	if !strings.Contains(text, "next hunk") {
		t.Error("usage text should document next hunk keybinding")
	}
	if !strings.Contains(text, "previous hunk") {
		t.Error("usage text should document previous hunk keybinding")
	}
	if !strings.Contains(text, "quit") {
		t.Error("usage text should document quit keybinding")
	}
}

func TestLazyPerFileLoading(t *testing.T) {
	dir, err := os.MkdirTemp("", "review-diff-lazy-*")
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

	os.WriteFile(dir+"/alpha.txt", []byte("alpha\n"), 0644)
	os.WriteFile(dir+"/beta.txt", []byte("beta\n"), 0644)
	os.WriteFile(dir+"/gamma.txt", []byte("gamma\n"), 0644)
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	os.WriteFile(dir+"/alpha.txt", []byte("alpha\nmain\n"), 0644)
	git("add", ".")
	git("commit", "-m", "main-change-alpha")

	git("checkout", "feature")
	os.WriteFile(dir+"/alpha.txt", []byte("alpha\nfeature\n"), 0644)
	os.WriteFile(dir+"/beta.txt", []byte("beta\nfeature\n"), 0644)
	os.WriteFile(dir+"/gamma.txt", []byte("gamma\nfeature\n"), 0644)
	git("add", ".")
	git("commit", "-m", "feature-changes")

	git("checkout", "main")

	runner := gitrunner.New(dir)
	ctx := context.Background()
	base := "main"
	head := "feature"

	changes, err := changelist.ListChanges(ctx, runner, base, head)
	if err != nil {
		t.Fatal(err)
	}

	if len(changes) < 3 {
		t.Fatalf("expected at least 3 changed files, got %d", len(changes))
	}

	found := map[string]bool{}
	for _, c := range changes {
		found[c.NewPath] = true
	}
	for _, name := range []string{"alpha.txt", "beta.txt", "gamma.txt"} {
		if !found[name] {
			t.Errorf("expected %s in changed files", name)
		}
	}

	patch, err := patchloader.LoadPatch(ctx, runner, base, head, "alpha.txt")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(patch.OldPath, "alpha.txt") && !strings.Contains(patch.NewPath, "alpha.txt") {
		t.Errorf("patch should reference alpha.txt, got old=%q new=%q", patch.OldPath, patch.NewPath)
	}
	if len(patch.Hunks) == 0 {
		t.Error("expected hunks in alpha.txt patch")
	}
	for _, h := range patch.Hunks {
		for _, line := range h.Lines {
			if strings.Contains(line.Content, "beta") {
				t.Errorf("alpha.txt patch should not contain beta content: %q", line.Content)
			}
			if strings.Contains(line.Content, "gamma") {
				t.Errorf("alpha.txt patch should not contain gamma content: %q", line.Content)
			}
		}
	}

	patch, err = patchloader.LoadPatch(ctx, runner, base, head, "gamma.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(patch.Hunks) == 0 {
		t.Error("expected hunks in gamma.txt patch")
	}
}

func TestEmptyDiff(t *testing.T) {
	dir, err := os.MkdirTemp("", "review-diff-empty-*")
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
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	runner := gitrunner.New(dir)
	ctx := context.Background()

	changes, err := changelist.ListChanges(ctx, runner, "main", "feature")
	if err != nil {
		t.Fatal(err)
	}

	if len(changes) != 0 {
		t.Errorf("expected no changes for identical branches, got %d", len(changes))
	}
}

func TestRenameDiff(t *testing.T) {
	dir, err := os.MkdirTemp("", "review-diff-rename-*")
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
	os.WriteFile(dir+"/old-name.txt", []byte("content\n"), 0644)
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	git("checkout", "feature")
	git("mv", "old-name.txt", "new-name.txt")
	git("commit", "-m", "rename-file")

	git("checkout", "main")

	runner := gitrunner.New(dir)
	ctx := context.Background()

	changes, err := changelist.ListChanges(ctx, runner, "main", "feature")
	if err != nil {
		t.Fatal(err)
	}

	var renameFound bool
	for _, c := range changes {
		if c.Status == "R" {
			renameFound = true
			if c.OldPath != "old-name.txt" {
				t.Errorf("rename old path: got %q, want %q", c.OldPath, "old-name.txt")
			}
			if c.NewPath != "new-name.txt" {
				t.Errorf("rename new path: got %q, want %q", c.NewPath, "new-name.txt")
			}
		}
	}
	if !renameFound {
		t.Fatal("expected rename in changes")
	}

	// Load patch via new name
	patch, err := patchloader.LoadPatch(ctx, runner, "main", "feature", "new-name.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(patch.Hunks) == 0 {
		t.Error("expected hunks in rename patch")
	}
}

func TestDeleteDiff(t *testing.T) {
	dir, err := os.MkdirTemp("", "review-diff-delete-*")
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
	os.WriteFile(dir+"/todelete.txt", []byte("line1\nline2\nline3\n"), 0644)
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	git("checkout", "feature")
	git("rm", "todelete.txt")
	git("commit", "-m", "delete-file")

	git("checkout", "main")

	runner := gitrunner.New(dir)
	ctx := context.Background()

	changes, err := changelist.ListChanges(ctx, runner, "main", "feature")
	if err != nil {
		t.Fatal(err)
	}

	var deleteFound bool
	for _, c := range changes {
		if c.Status == "D" && c.NewPath == "todelete.txt" {
			deleteFound = true
		}
	}
	if !deleteFound {
		t.Fatal("expected deleted file in changes")
	}

	// Load patch for deleted file should not crash
	patch, err := patchloader.LoadPatch(ctx, runner, "main", "feature", "todelete.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(patch.Hunks) == 0 {
		t.Error("expected hunks in delete patch")
	}
	// All lines should be removed type
	for _, h := range patch.Hunks {
		for _, line := range h.Lines {
			if line.Type != 2 {
				t.Errorf("expected removed lines in delete patch, got type %d: %q", line.Type, line.Content)
			}
		}
	}
}

func TestBinaryDiff(t *testing.T) {
	dir, err := os.MkdirTemp("", "review-diff-binary-*")
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
	os.WriteFile(dir+"/file.bin", []byte("hello\x00world\n"), 0644)
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	git("checkout", "feature")
	os.WriteFile(dir+"/file.bin", []byte("modified\x00content\n"), 0644)
	git("add", ".")
	git("commit", "-m", "modify-binary")

	git("checkout", "main")

	runner := gitrunner.New(dir)
	ctx := context.Background()

	changes, err := changelist.ListChanges(ctx, runner, "main", "feature")
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 {
		t.Fatal("expected changes in binary test")
	}

	// Load patch should detect binary and not crash
	patch, err := patchloader.LoadPatch(ctx, runner, "main", "feature", "file.bin")
	if err != nil {
		t.Fatal(err)
	}
	if !patch.IsBinary {
		t.Error("expected IsBinary true for binary file")
	}
}

func TestLinkedWorktreeDiff(t *testing.T) {
	dir := t.TempDir()
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

	os.WriteFile(dir+"/shared.txt", []byte("shared\n"), 0644)
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	featureDir := t.TempDir()
	git("worktree", "add", featureDir, "feature")

	git2 := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = featureDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}

	os.WriteFile(featureDir+"/feature-only.txt", []byte("feature change\n"), 0644)
	git2("add", ".")
	git2("commit", "-m", "feature-change")

	ctx := context.Background()
	sess, err := session.New(dir, "", "", dir, featureDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := sess.ResolveWorktreeRefs(ctx); err != nil {
		t.Fatal(err)
	}
	if sess.Base == "" {
		t.Fatal("expected resolved base ref")
	}
	if sess.Head == "" {
		t.Fatal("expected resolved head ref")
	}

	runner := gitrunner.New(dir)
	changes, err := changelist.ListChanges(ctx, runner, sess.Base, sess.Head)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, c := range changes {
		if c.NewPath == "feature-only.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected feature-only.txt in changes, got %v", changes)
	}

	patch, err := patchloader.LoadPatch(ctx, runner, sess.Base, sess.Head, "feature-only.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(patch.Hunks) == 0 {
		t.Error("expected hunks for feature-only.txt")
	}
}

func TestLinkedWorktreeDiffWithRefOverride(t *testing.T) {
	dir := t.TempDir()
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

	os.WriteFile(dir+"/shared.txt", []byte("shared\n"), 0644)
	git("add", ".")
	git("commit", "-m", "init")
	git("branch", "feature")

	featureDir := t.TempDir()
	git("worktree", "add", featureDir, "feature")

	git2 := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = featureDir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}

	os.WriteFile(featureDir+"/branch-file.txt", []byte("feature change\n"), 0644)
	git2("add", ".")
	git2("commit", "-m", "feature-change")

	ctx := context.Background()
	sess, err := session.New(dir, "main", "feature", dir, featureDir)
	if err != nil {
		t.Fatal(err)
	}
	if err := sess.ResolveWorktreeRefs(ctx); err != nil {
		t.Fatal(err)
	}
	if sess.Base != "main" {
		t.Errorf("base should remain 'main', got %q", sess.Base)
	}
	if sess.Head != "feature" {
		t.Errorf("head should remain 'feature', got %q", sess.Head)
	}

	runner := gitrunner.New(dir)
	changes, err := changelist.ListChanges(ctx, runner, sess.Base, sess.Head)
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, c := range changes {
		if c.NewPath == "branch-file.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected branch-file.txt in changes, got %v", changes)
	}
}

func TestUnlinkedWorktreePaths(t *testing.T) {
	repo1 := t.TempDir()
	repo2 := t.TempDir()

	for _, d := range []string{repo1, repo2} {
		cmd := exec.Command("git", "init", "-b", "main")
		cmd.Dir = d
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git init in %s: %s", d, out)
		}
		cmd = exec.Command("git", "config", "user.email", "t@t.com")
		cmd.Dir = d
		cmd.Run()
		cmd = exec.Command("git", "config", "user.name", "T")
		cmd.Dir = d
		cmd.Run()
		os.WriteFile(d+"/init.txt", []byte("x\n"), 0644)
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = d
		cmd.Run()
		cmd = exec.Command("git", "commit", "-m", "init")
		cmd.Dir = d
		cmd.Run()
	}

	ctx := context.Background()
	sess, err := session.New(repo1, "", "", repo1, repo2)
	if err != nil {
		t.Fatal(err)
	}
	err = sess.ResolveWorktreeRefs(ctx)
	if err == nil {
		t.Fatal("expected error for unlinked worktree paths")
	}
	if !strings.Contains(err.Error(), "not linked") {
		t.Errorf("expected error about not linked, got: %v", err)
	}
}

func TestFetchIntegration(t *testing.T) {
	// Create a local bare remote
	bareDir := t.TempDir()
	cmd := exec.Command("git", "init", "--bare", bareDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init --bare: %s", out)
	}

	// Working copy A: push main + feature branches to bare remote
	dirA := t.TempDir()
	gitA := func(args ...string) {
		t.Helper()
		c := exec.Command("git", args...)
		c.Dir = dirA
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}
	gitA("init", "-b", "main")
	gitA("config", "user.email", "t@t.com")
	gitA("config", "user.name", "T")
	gitA("remote", "add", "origin", bareDir)

	os.WriteFile(dirA+"/init.txt", []byte("init\n"), 0644)
	gitA("add", ".")
	gitA("commit", "-m", "init")
	gitA("push", "-u", "origin", "main")

	gitA("checkout", "-b", "feature")
	os.WriteFile(dirA+"/feature.txt", []byte("feature change\n"), 0644)
	gitA("add", ".")
	gitA("commit", "-m", "feature-commit")
	gitA("push", "-u", "origin", "feature")

	gitA("checkout", "main")

	// Working copy B: clone but only fetch main ref, so origin/feature is missing
	dirB := t.TempDir()
	gitB := func(args ...string) {
		t.Helper()
		c := exec.Command("git", args...)
		c.Dir = dirB
		if out, err := c.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %s", args, out)
		}
	}
	gitB("init", "-b", "main")
	gitB("config", "user.email", "t@t.com")
	gitB("config", "user.name", "T")
	gitB("remote", "add", "origin", bareDir)
	gitB("fetch", "origin", "main")

	// Verify feature branch ref does not exist in B yet
	checkRef := exec.Command("git", "show-ref", "--verify", "refs/remotes/origin/feature")
	checkRef.Dir = dirB
	if err := checkRef.Run(); err == nil {
		t.Fatal("expected origin/feature to not exist before fetch")
	}

	ctx := context.Background()
	runner := gitrunner.New(dirB)

	// ListChanges against the missing ref should fail
	_, err := changelist.ListChanges(ctx, runner, "origin/main", "origin/feature")
	if err == nil {
		t.Fatal("expected error listing changes with non-existent ref origin/feature")
	}

	// Run fetch --all to bring in the missing ref
	if _, err := runner.Run(ctx, "fetch", "--all"); err != nil {
		t.Fatalf("fetch failed: %v", err)
	}

	// Now origin/feature should exist and ListChanges should succeed
	changes, err := changelist.ListChanges(ctx, runner, "origin/main", "origin/feature")
	if err != nil {
		t.Fatalf("ListChanges after fetch: %v", err)
	}
	var found bool
	for _, c := range changes {
		if c.NewPath == "feature.txt" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected feature.txt in changes after fetch, got %v", changes)
	}
}
