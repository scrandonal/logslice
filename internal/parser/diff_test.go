package parser

import (
	"strings"
	"testing"
	"time"
)

func makeDiffLine(raw string, ts *time.Time) Line {
	return Line{Raw: raw, Timestamp: ts}
}

func diffTime(s string) *time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return &t
}

func TestDiffNoLines(t *testing.T) {
	res := DiffLogs(nil, nil, DiffOptions{})
	if len(res.OnlyLeft) != 0 || len(res.OnlyRight) != 0 || len(res.Common) != 0 {
		t.Fatal("expected empty result")
	}
}

func TestDiffIdentical(t *testing.T) {
	lines := []Line{
		makeDiffLine("[2024-01-01T00:00:00Z] foo", diffTime("2024-01-01T00:00:00Z")),
		makeDiffLine("[2024-01-01T00:00:01Z] bar", diffTime("2024-01-01T00:00:01Z")),
	}
	res := DiffLogs(lines, lines, DiffOptions{})
	if len(res.Common) != 2 {
		t.Fatalf("expected 2 common, got %d", len(res.Common))
	}
	if len(res.OnlyLeft) != 0 || len(res.OnlyRight) != 0 {
		t.Fatal("expected no unique lines")
	}
}

func TestDiffDisjoint(t *testing.T) {
	left := []Line{makeDiffLine("alpha", nil)}
	right := []Line{makeDiffLine("beta", nil)}
	res := DiffLogs(left, right, DiffOptions{})
	if len(res.OnlyLeft) != 1 || res.OnlyLeft[0].Raw != "alpha" {
		t.Fatal("expected alpha in only-left")
	}
	if len(res.OnlyRight) != 1 || res.OnlyRight[0].Raw != "beta" {
		t.Fatal("expected beta in only-right")
	}
	if len(res.Common) != 0 {
		t.Fatal("expected no common lines")
	}
}

func TestDiffIgnoreTimestamp(t *testing.T) {
	t1 := diffTime("2024-01-01T10:00:00Z")
	t2 := diffTime("2024-01-01T10:00:01Z")
	left := []Line{makeDiffLine("[2024-01-01T10:00:00Z] hello world", t1)}
	right := []Line{makeDiffLine("[2024-01-01T10:00:01Z] hello world", t2)}

	// Without ignore: different raw text → not matched
	res := DiffLogs(left, right, DiffOptions{})
	if len(res.Common) != 0 {
		t.Fatal("expected no common without IgnoreTimestamp")
	}

	// With ignore: same body → matched
	res = DiffLogs(left, right, DiffOptions{IgnoreTimestamp: true})
	if len(res.Common) != 1 {
		t.Fatalf("expected 1 common with IgnoreTimestamp, got %d", len(res.Common))
	}
}

func TestDiffWindowExcludes(t *testing.T) {
	t1 := diffTime("2024-01-01T10:00:00Z")
	t2 := diffTime("2024-01-01T10:30:00Z") // 30 min apart
	left := []Line{makeDiffLine("same line", t1)}
	right := []Line{makeDiffLine("same line", t2)}

	res := DiffLogs(left, right, DiffOptions{Window: 5 * time.Minute})
	if len(res.Common) != 0 {
		t.Fatal("expected no common: timestamps too far apart")
	}
}

func TestPrintDiffNoDifferences(t *testing.T) {
	var sb strings.Builder
	printDiffTo(&sb, DiffResult{})
	if !strings.Contains(sb.String(), "no differences") {
		t.Fatalf("unexpected output: %q", sb.String())
	}
}

func TestDiffJSON(t *testing.T) {
	res := DiffResult{
		OnlyLeft:  []Line{{Raw: "left-only"}},
		OnlyRight: []Line{{Raw: "right-only"}},
		Common:    []Line{{Raw: "shared"}},
	}
	out := DiffJSON(res)
	if !strings.Contains(out, "left-only") {
		t.Error("missing left-only in JSON")
	}
	if !strings.Contains(out, "right-only") {
		t.Error("missing right-only in JSON")
	}
	if !strings.Contains(out, `"common":1`) {
		t.Error("missing common count in JSON")
	}
}
