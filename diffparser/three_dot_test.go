package diffparser

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestThreeDotSemantics(t *testing.T) {
	dir, err := os.MkdirTemp("", "review-diff-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %s\n%s", args, out, err)
		}
	}

	git("init", "-b", "main")
	git("config", "user.email", "test@test.com")
	git("config", "user.name", "Test")

	// Commit A: initial on main
	writeFile(t, dir, "file.txt", "line1\nline2\nline3\n")
	git("add", ".")
	git("commit", "-m", "initial")

	// Create feature branch at commit A
	git("branch", "feature", "main")

	// Commit B: main-only change
	writeFile(t, dir, "file.txt", "line1\nline2\nline3\nmain-only\n")
	git("add", ".")
	git("commit", "-m", "main-only")

	// Commit C: feature-only change
	git("checkout", "feature")
	writeFile(t, dir, "file.txt", "line1\nline2\nline3\nfeature-only\n")
	git("add", ".")
	git("commit", "-m", "feature-only")

	// Verify merge-base is the initial commit
	git("checkout", "main")

	// three-dot: main...feature = diff between merge-base (A) and feature (C)
	threeDotOut, err := exec.Command("git", "-C", dir, "diff", "main...feature").CombinedOutput()
	if err != nil {
		t.Fatalf("three-dot diff failed: %s", threeDotOut)
	}

	// two-dot: main..feature = diff between main (B) and feature (C)
	twoDotOut, err := exec.Command("git", "-C", dir, "diff", "main..feature").CombinedOutput()
	if err != nil {
		t.Fatalf("two-dot diff failed: %s", twoDotOut)
	}

	// Three-dot should show only feature-only change (one addition)
	threeDotDF, err := ParseUnified(string(threeDotOut))
	if err != nil {
		t.Fatal(err)
	}
	if len(threeDotDF.Hunks) == 0 {
		t.Fatal("three-dot: expected at least one hunk")
	}
	hasFeatureOnly := false
	hasMainOnly := false
	for _, h := range threeDotDF.Hunks {
		for _, l := range h.Lines {
			if l.Type == LineAdded && strings.Contains(l.Content, "feature-only") {
				hasFeatureOnly = true
			}
			if l.Type == LineAdded && strings.Contains(l.Content, "main-only") {
				hasMainOnly = true
			}
		}
	}
	if !hasFeatureOnly {
		t.Error("three-dot: expected to see 'feature-only' addition (diff between merge-base and feature)")
	}
	if hasMainOnly {
		t.Error("three-dot: should NOT see 'main-only' (that is between merge-base and main, not in three-dot range)")
	}

	// Two-dot should show both sides: removal of main-only and addition of feature-only
	twoDotDF, err := ParseUnified(string(twoDotOut))
	if err != nil {
		t.Fatal(err)
	}
	// In two-dot, main..feature means diff from main to feature
	// So main-only disappears and feature-only appears
	hasFeatureAdded := false
	hasMainRemoved := false
	for _, h := range twoDotDF.Hunks {
		for _, l := range h.Lines {
			if l.Type == LineAdded && strings.Contains(l.Content, "feature-only") {
				hasFeatureAdded = true
			}
			if l.Type == LineRemoved && strings.Contains(l.Content, "main-only") {
				hasMainRemoved = true
			}
		}
	}
	if !hasFeatureAdded {
		t.Error("two-dot: expected to see 'feature-only' added")
	}
	if !hasMainRemoved {
		t.Error("two-dot: expected to see 'main-only' removed")
	}

	t.Logf("Three-dot output:\n%s", threeDotOut)
	t.Logf("Two-dot output:\n%s", twoDotOut)
}

func writeFile(t *testing.T, dir, path, content string) {
	t.Helper()
	if err := os.WriteFile(dir+"/"+path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
