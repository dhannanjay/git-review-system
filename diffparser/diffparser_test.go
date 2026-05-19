package diffparser

import (
	"testing"
)

const singleHunkAdd = `--- a/hello.txt
+++ b/hello.txt
@@ -1,3 +1,4 @@
 line one
 line two
-line three
+line three updated
+line four
`

const multiHunkDiff = `--- a/numbers.txt
+++ b/numbers.txt
@@ -1,5 +1,6 @@
 1
 2
+two-and-half
 3
 4
 5
@@ -10,4 +11,5 @@
 ten
 eleven
 twelve
+twelve-point-five
 thirteen
`

const mixedAddDelete = `--- a/colors.txt
+++ b/colors.txt
@@ -1,6 +1,5 @@
 red
 green
-blue
 yellow
+orange
 purple
-indigo
`

func TestParseUnified_singleHunkAdd(t *testing.T) {
	df, err := ParseUnified(singleHunkAdd)
	if err != nil {
		t.Fatal(err)
	}
	if df.OldPath != "a/hello.txt" {
		t.Errorf("old path: got %q, want %q", df.OldPath, "a/hello.txt")
	}
	if df.NewPath != "b/hello.txt" {
		t.Errorf("new path: got %q, want %q", df.NewPath, "b/hello.txt")
	}
	if len(df.Hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(df.Hunks))
	}
	h := df.Hunks[0]
	if h.OldStart != 1 || h.OldCount != 3 {
		t.Errorf("old range: got %d,%d want 1,3", h.OldStart, h.OldCount)
	}
	if h.NewStart != 1 || h.NewCount != 4 {
		t.Errorf("new range: got %d,%d want 1,4", h.NewStart, h.NewCount)
	}
	if len(h.Lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(h.Lines))
	}
	// Context, context, removed, added, added
	checkLine(t, h.Lines[0], LineContext, "line one", 1, 1)
	checkLine(t, h.Lines[1], LineContext, "line two", 2, 2)
	checkLine(t, h.Lines[2], LineRemoved, "line three", 3, 0)
	checkLine(t, h.Lines[3], LineAdded, "line three updated", 0, 3)
	checkLine(t, h.Lines[4], LineAdded, "line four", 0, 4)
}

func TestParseUnified_multiHunk(t *testing.T) {
	df, err := ParseUnified(multiHunkDiff)
	if err != nil {
		t.Fatal(err)
	}
	if len(df.Hunks) != 2 {
		t.Fatalf("expected 2 hunks, got %d", len(df.Hunks))
	}

	// First hunk
	h0 := df.Hunks[0]
	if h0.OldStart != 1 || h0.OldCount != 5 {
		t.Errorf("h0 old range: got %d,%d want 1,5", h0.OldStart, h0.OldCount)
	}
	if h0.NewStart != 1 || h0.NewCount != 6 {
		t.Errorf("h0 new range: got %d,%d want 1,6", h0.NewStart, h0.NewCount)
	}
	if len(h0.Lines) != 6 {
		t.Fatalf("h0 expected 6 lines, got %d", len(h0.Lines))
	}
	checkLine(t, h0.Lines[0], LineContext, "1", 1, 1)
	checkLine(t, h0.Lines[1], LineContext, "2", 2, 2)
	checkLine(t, h0.Lines[2], LineAdded, "two-and-half", 0, 3)
	checkLine(t, h0.Lines[3], LineContext, "3", 3, 4)
	checkLine(t, h0.Lines[4], LineContext, "4", 4, 5)
	checkLine(t, h0.Lines[5], LineContext, "5", 5, 6)

	// Second hunk
	h1 := df.Hunks[1]
	if h1.OldStart != 10 || h1.OldCount != 4 {
		t.Errorf("h1 old range: got %d,%d want 10,4", h1.OldStart, h1.OldCount)
	}
	if h1.NewStart != 11 || h1.NewCount != 5 {
		t.Errorf("h1 new range: got %d,%d want 11,5", h1.NewStart, h1.NewCount)
	}
	if len(h1.Lines) != 5 {
		t.Fatalf("h1 expected 5 lines, got %d", len(h1.Lines))
	}
	checkLine(t, h1.Lines[0], LineContext, "ten", 10, 11)
	checkLine(t, h1.Lines[1], LineContext, "eleven", 11, 12)
	checkLine(t, h1.Lines[2], LineContext, "twelve", 12, 13)
	checkLine(t, h1.Lines[3], LineAdded, "twelve-point-five", 0, 14)
	checkLine(t, h1.Lines[4], LineContext, "thirteen", 13, 15)
}

func TestParseUnified_mixedAddDelete(t *testing.T) {
	df, err := ParseUnified(mixedAddDelete)
	if err != nil {
		t.Fatal(err)
	}
	if len(df.Hunks) != 1 {
		t.Fatalf("expected 1 hunk, got %d", len(df.Hunks))
	}
	h := df.Hunks[0]
	if h.OldStart != 1 || h.OldCount != 6 {
		t.Errorf("old range: got %d,%d want 1,6", h.OldStart, h.OldCount)
	}
	if h.NewStart != 1 || h.NewCount != 5 {
		t.Errorf("new range: got %d,%d want 1,5", h.NewStart, h.NewCount)
	}
	if len(h.Lines) != 7 {
		t.Fatalf("expected 7 lines, got %d", len(h.Lines))
	}
	checkLine(t, h.Lines[0], LineContext, "red", 1, 1)
	checkLine(t, h.Lines[1], LineContext, "green", 2, 2)
	checkLine(t, h.Lines[2], LineRemoved, "blue", 3, 0)
	checkLine(t, h.Lines[3], LineContext, "yellow", 4, 3)
	checkLine(t, h.Lines[4], LineAdded, "orange", 0, 4)
	checkLine(t, h.Lines[5], LineContext, "purple", 5, 5)
	checkLine(t, h.Lines[6], LineRemoved, "indigo", 6, 0)
}

func TestParseUnified_empty(t *testing.T) {
	df, err := ParseUnified("")
	if err != nil {
		t.Fatal(err)
	}
	if df == nil {
		t.Fatal("expected non-nil DiffFile")
	}
	if len(df.Hunks) != 0 {
		t.Errorf("expected 0 hunks, got %d", len(df.Hunks))
	}
}

func checkLine(t *testing.T, l Line, typ LineType, content string, oldNum, newNum int) {
	t.Helper()
	if l.Type != typ {
		t.Errorf("line type: got %d want %d for content %q", l.Type, typ, l.Content)
	}
	if l.Content != content {
		t.Errorf("line content: got %q want %q", l.Content, content)
	}
	if l.OldNum != oldNum {
		t.Errorf("line oldNum: got %d want %d for content %q", l.OldNum, oldNum, l.Content)
	}
	if l.NewNum != newNum {
		t.Errorf("line newNum: got %d want %d for content %q", l.NewNum, newNum, l.Content)
	}
}
