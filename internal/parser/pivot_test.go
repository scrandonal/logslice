package parser

import (
	"testing"
	"time"
)

func makePivotLine(raw string, ts *time.Time) Line {
	return Line{Raw: raw, Timestamp: ts}
}

func ptrPivotTime(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return &t
}

func TestPivotEmpty(t *testing.T) {
	result := NewPivot(nil, PivotOptions{Field: "level", BucketSize: time.Minute})
	if len(result.Rows) != 0 {
		t.Fatalf("expected 0 rows, got %d", len(result.Rows))
	}
	if len(result.Buckets) != 0 {
		t.Fatalf("expected 0 buckets, got %d", len(result.Buckets))
	}
}

func TestPivotSingleField(t *testing.T) {
	lines := []Line{
		makePivotLine(`level=info msg="ok"`, ptrPivotTime("2024-01-01T10:00:00Z")),
		makePivotLine(`level=error msg="fail"`, ptrPivotTime("2024-01-01T10:00:30Z")),
		makePivotLine(`level=info msg="ok2"`, ptrPivotTime("2024-01-01T10:01:00Z")),
	}
	result := NewPivot(lines, PivotOptions{Field: "level", BucketSize: time.Minute})

	if result.Field != "level" {
		t.Errorf("expected field=level, got %q", result.Field)
	}
	if len(result.Buckets) != 2 {
		t.Errorf("expected 2 buckets, got %d", len(result.Buckets))
	}
	if len(result.Rows) != 2 {
		t.Errorf("expected 2 rows (info, error), got %d", len(result.Rows))
	}
	// info should be first (higher total)
	if result.Rows[0].Value != "info" {
		t.Errorf("expected first row to be info, got %q", result.Rows[0].Value)
	}
	if result.Rows[0].Total != 2 {
		t.Errorf("expected info total=2, got %d", result.Rows[0].Total)
	}
}

func TestPivotMaxValues(t *testing.T) {
	lines := []Line{
		makePivotLine(`level=info`, ptrPivotTime("2024-01-01T10:00:00Z")),
		makePivotLine(`level=warn`, ptrPivotTime("2024-01-01T10:00:00Z")),
		makePivotLine(`level=error`, ptrPivotTime("2024-01-01T10:00:00Z")),
	}
	result := NewPivot(lines, PivotOptions{Field: "level", BucketSize: time.Minute, MaxValues: 2})
	if len(result.Rows) != 2 {
		t.Errorf("expected 2 rows due to MaxValues, got %d", len(result.Rows))
	}
}

func TestPivotNoTimestamp(t *testing.T) {
	lines := []Line{
		{Raw: `level=info`, Timestamp: nil},
		{Raw: `level=info`, Timestamp: nil},
	}
	result := NewPivot(lines, PivotOptions{Field: "level", BucketSize: time.Minute})
	if len(result.Buckets) != 1 || result.Buckets[0] != "(no time)" {
		t.Errorf("expected single '(no time)' bucket, got %v", result.Buckets)
	}
	if result.Rows[0].Total != 2 {
		t.Errorf("expected total=2, got %d", result.Rows[0].Total)
	}
}
