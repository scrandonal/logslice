package parser

import (
	"testing"
	"time"
)

func makeTime(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return &t
}

func makeLine(raw string, ts *time.Time) Line {
	return Line{Raw: raw, Timestamp: ts}
}

func TestMergeLinesOrdering(t *testing.T) {
	src0 := []Line{
		makeLine("a", makeTime("2024-01-01T10:00:00Z")),
		makeLine("c", makeTime("2024-01-01T10:02:00Z")),
	}
	src1 := []Line{
		makeLine("b", makeTime("2024-01-01T10:01:00Z")),
		makeLine("d", makeTime("2024-01-01T10:03:00Z")),
	}

	merged := MergeLines([][]Line{src0, src1})
	if len(merged) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(merged))
	}
	expected := []string{"a", "b", "c", "d"}
	for i, ml := range merged {
		if ml.Raw != expected[i] {
			t.Errorf("pos %d: want %q, got %q", i, expected[i], ml.Raw)
		}
	}
}

func TestMergeLinesNoTimestamp(t *testing.T) {
	src := []Line{
		makeLine("no-ts", nil),
		makeLine("ts", makeTime("2024-01-01T10:00:00Z")),
	}
	merged := MergeLines([][]Line{src})
	if len(merged) != 2 {
		t.Fatalf("expected 2, got %d", len(merged))
	}
	if merged[0].Raw != "ts" {
		t.Errorf("timestamped line should come first")
	}
	if merged[1].Raw != "no-ts" {
		t.Errorf("no-timestamp line should come last")
	}
}

func TestMergeTimeRange(t *testing.T) {
	src := []Line{
		makeLine("early", makeTime("2024-01-01T09:00:00Z")),
		makeLine("in", makeTime("2024-01-01T10:00:00Z")),
		makeLine("late", makeTime("2024-01-01T11:00:00Z")),
	}
	from := *makeTime("2024-01-01T09:30:00Z")
	to := *makeTime("2024-01-01T10:30:00Z")

	result := MergeTimeRange([][]Line{src}, from, to)
	if len(result) != 1 {
		t.Fatalf("expected 1 line, got %d", len(result))
	}
	if result[0].Raw != "in" {
		t.Errorf("expected 'in', got %q", result[0].Raw)
	}
}

func TestMergeTimeRangeEmpty(t *testing.T) {
	result := MergeTimeRange([][]Line{}, time.Now(), time.Now())
	if len(result) != 0 {
		t.Errorf("expected empty result for empty sources")
	}
}
