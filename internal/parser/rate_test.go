package parser

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeRateLine(ts *time.Time, raw string) Line {
	return Line{Timestamp: ts, Raw: raw}
}

func rateTime(s string) *time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05", s)
	return &t
}

func TestCalcRateEmpty(t *testing.T) {
	pts := CalcRate(nil, time.Minute)
	if len(pts) != 0 {
		t.Fatalf("expected 0 points, got %d", len(pts))
	}
}

func TestCalcRateSingleBucket(t *testing.T) {
	lines := []Line{
		makeRateLine(rateTime("2024-01-01T10:00:10"), "a"),
		makeRateLine(rateTime("2024-01-01T10:00:45"), "b"),
		makeRateLine(rateTime("2024-01-01T10:00:59"), "c"),
	}
	pts := CalcRate(lines, time.Minute)
	if len(pts) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(pts))
	}
	if pts[0].Count != 3 {
		t.Errorf("expected count 3, got %d", pts[0].Count)
	}
}

func TestCalcRateMultipleBuckets(t *testing.T) {
	lines := []Line{
		makeRateLine(rateTime("2024-01-01T10:00:05"), "a"),
		makeRateLine(rateTime("2024-01-01T10:01:10"), "b"),
		makeRateLine(rateTime("2024-01-01T10:01:50"), "c"),
		makeRateLine(rateTime("2024-01-01T10:02:01"), "d"),
	}
	pts := CalcRate(lines, time.Minute)
	if len(pts) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(pts))
	}
	if pts[1].Count != 2 {
		t.Errorf("expected bucket[1] count 2, got %d", pts[1].Count)
	}
}

func TestCalcRateNoTimestamp(t *testing.T) {
	lines := []Line{
		makeRateLine(nil, "no ts 1"),
		makeRateLine(nil, "no ts 2"),
	}
	pts := CalcRate(lines, time.Minute)
	if len(pts) != 1 {
		t.Fatalf("expected 1 bucket for nil timestamps, got %d", len(pts))
	}
	if pts[0].Count != 2 {
		t.Errorf("expected count 2, got %d", pts[0].Count)
	}
}

func TestPrintRateEmpty(t *testing.T) {
	var buf bytes.Buffer
	PrintRate(nil, time.Minute, &buf)
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected 'no data', got %q", buf.String())
	}
}

func TestRateJSON(t *testing.T) {
	lines := []Line{
		makeRateLine(rateTime("2024-03-15T09:01:00"), "x"),
		makeRateLine(rateTime("2024-03-15T09:02:00"), "y"),
	}
	pts := CalcRate(lines, time.Minute)
	out := RateJSON(pts, time.Minute)
	if !strings.HasPrefix(out, "[") {
		t.Errorf("expected JSON array, got %q", out)
	}
	if !strings.Contains(out, "count") {
		t.Errorf("expected 'count' key in JSON, got %q", out)
	}
}
