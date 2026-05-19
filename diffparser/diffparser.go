package diffparser

import (
	"strconv"
	"strings"
)

type DiffFile struct {
	OldPath string
	NewPath string
	Hunks   []Hunk
}

type Hunk struct {
	OldStart int
	OldCount int
	NewStart int
	NewCount int
	Header   string
	Lines    []Line
}

type LineType int

const (
	LineContext LineType = iota
	LineAdded
	LineRemoved
)

type Line struct {
	Type    LineType
	Content string
	OldNum  int
	NewNum  int
}

func ParseUnified(input string) (*DiffFile, error) {
	lines := strings.Split(input, "\n")

	// Remove trailing empty line from split
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	df := &DiffFile{}

	i := 0
	// Skip header lines (anything before the first @@)
	for i < len(lines) {
		if strings.HasPrefix(lines[i], "--- ") {
			df.OldPath = strings.TrimPrefix(lines[i], "--- ")
			i++
		} else if strings.HasPrefix(lines[i], "+++ ") {
			df.NewPath = strings.TrimPrefix(lines[i], "+++ ")
			i++
		} else if strings.HasPrefix(lines[i], "@@") {
			break
		} else {
			i++
		}
	}

	for i < len(lines) {
		if !strings.HasPrefix(lines[i], "@@") {
			i++
			continue
		}

		hunk, nextIdx, err := parseHunk(lines, i)
		if err != nil {
			return nil, err
		}
		df.Hunks = append(df.Hunks, hunk)
		i = nextIdx
	}

	return df, nil
}

func parseHunk(lines []string, start int) (Hunk, int, error) {
	hunkLine := lines[start]
	// @@ -oldStart,oldCount +newStart,newCount @@ optional header
	parts := strings.SplitN(hunkLine, "@@", 3)
	if len(parts) < 3 {
		// still try to parse
	}
	rangePart := strings.TrimSpace(parts[1])

	header := ""
	if len(parts) == 3 {
		header = strings.TrimSpace(parts[2])
	}

	// Parse -oldStart,oldCount +newStart,newCount
	ranges := strings.SplitN(rangePart, " ", 3)
	var oldRange, newRange string
	for _, r := range ranges {
		if strings.HasPrefix(r, "-") {
			oldRange = strings.TrimPrefix(r, "-")
		} else if strings.HasPrefix(r, "+") {
			newRange = strings.TrimPrefix(r, "+")
		}
	}

	oldStart, oldCount := parseRange(oldRange)
	newStart, newCount := parseRange(newRange)

	h := Hunk{
		OldStart: oldStart,
		OldCount: oldCount,
		NewStart: newStart,
		NewCount: newCount,
		Header:   header,
	}

	i := start + 1
	oldNum := oldStart
	newNum := newStart

	for i < len(lines) {
		line := lines[i]
		if strings.HasPrefix(line, "@@") {
			break
		}
		var l Line
		if strings.HasPrefix(line, "+") {
			l = Line{Type: LineAdded, Content: line[1:], OldNum: 0, NewNum: newNum}
			newNum++
		} else if strings.HasPrefix(line, "-") {
			l = Line{Type: LineRemoved, Content: line[1:], OldNum: oldNum, NewNum: 0}
			oldNum++
		} else if strings.HasPrefix(line, " ") {
			l = Line{Type: LineContext, Content: line[1:], OldNum: oldNum, NewNum: newNum}
			oldNum++
			newNum++
		} else if strings.HasPrefix(line, "\\") {
			// No newline at end of file marker - skip or include as context
			l = Line{Type: LineContext, Content: line, OldNum: 0, NewNum: 0}
		}
		h.Lines = append(h.Lines, l)
		i++
	}

	return h, i, nil
}

func parseRange(r string) (int, int) {
	parts := strings.SplitN(r, ",", 2)
	start, _ := strconv.Atoi(parts[0])
	count := 1
	if len(parts) == 2 {
		c, err := strconv.Atoi(parts[1])
		if err == nil {
			count = c
		}
	}
	return start, count
}
