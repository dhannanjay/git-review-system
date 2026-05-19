package main

import (
	"strings"
	"testing"
)

func TestUsageText(t *testing.T) {
	text := usageText()
	if !strings.Contains(text, "review-diff") {
		t.Error("usage text should contain binary name")
	}
	if !strings.Contains(text, "--help") {
		t.Error("usage text should mention --help")
	}
	if !strings.Contains(text, "blocks") {
		t.Error("usage text should mention blocking behavior")
	}
}

func TestUsageTextIncludesFlagsNote(t *testing.T) {
	text := usageText()
	if !strings.Contains(text, "--base") || !strings.Contains(text, "--head") {
		t.Error("usage text should mention forthcoming branch-diff flags")
	}
}
