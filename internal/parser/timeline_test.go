package parser

import (
	"strings"
	"testing"
	"time"
)

func makeTimelineLine(ts *time.Time, raw string) Line {
	return Line{Timestamp: ts, Raw: raw}
}

func tlTime(s string) *time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return &t
}

func TestBuildTimelineEmpty(t *testing.T) {
	res := BuildTimeline(nil, time.Minute)
	if len(res.Buckets) != 0 {
		t.Fatalf("expected 0 buckets, got %d", len(res.Buckets))
	}
	if res.Total != 0 {
		t.Fatalf("expected total 0, got %d", res.Total)
	}
}

func TestBuildTimelineSingleBucket(t *testing.T) {
	lines := []Line{
		makeTimelineLine(tlTime("2024-01-01 10:00:30"), "a"),
		makeTimelineLine(tlTime("2024-01-01 10:00:45"), "b"),
		makeTimelineLine(tlTime("2024-01-01 10:00:59"), "c"),
	}
	res := BuildTimeline(lines, time.Minute)
	if len(res.Buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(res.Buckets))
	}
	if res.Buckets[0].Count != 3 {
		t.Errorf("expected count 3, got %d", res.Buckets[0].Count)
	}
	if res.Total != 3 {
		t.Errorf("expected total 3, got %d", res.Total)
	}
}

func TestBuildTimelineMultipleBuckets(t *testing.T) {
	lines := []Line{
		makeTimelineLine(tlTime("2024-01-01 10:00:10"), "a"),
		makeTimelineLine(tlTime("2024-01-01 10:01:10"), "b"),
		makeTimelineLine(tlTime("2024-01-01 10:01:50"), "c"),
		makeTimelineLine(nil, "no-ts"),
	}
	res := BuildTimeline(lines, time.Minute)
	if len(res.Buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(res.Buckets))
	}
	if res.Buckets[0].Count != 1 {
		t.Errorf("bucket[0] count: want 1, got %d", res.Buckets[0].Count)
	}
	if res.Buckets[1].Count != 2 {
		t.Errorf("bucket[1] count: want 2, got %d", res.Buckets[1].Count)
	}
	if res.Total != 4 {
		t.Errorf("total: want 4, got %d", res.Total)
	}
}

func TestPrintTimelineEmpty(t *testing.T) {
	var sb strings.Builder
	PrintTimeline(&sb, TimelineResult{})
	if !strings.Contains(sb.String(), "no timestamped") {
		t.Errorf("expected no-timestamp message, got: %s", sb.String())
	}
}

func TestPrintTimelineOutput(t *testing.T) {
	lines := []Line{
		makeTimelineLine(tlTime("2024-01-01 10:00:10"), "a"),
		makeTimelineLine(tlTime("2024-01-01 10:01:10"), "b"),
	}
	res := BuildTimeline(lines, time.Minute)
	var sb strings.Builder
	PrintTimeline(&sb, res)
	out := sb.String()
	if !strings.Contains(out, "10:00:00") {
		t.Errorf("expected bucket label in output, got: %s", out)
	}
	if !strings.Contains(out, "total:") {
		t.Errorf("expected total line in output, got: %s", out)
	}
}

func TestTimelineJSON(t *testing.T) {
	lines := []Line{
		makeTimelineLine(tlTime("2024-01-01 10:00:10"), "a"),
	}
	res := BuildTimeline(lines, time.Minute)
	j := TimelineJSON(res)
	if !strings.Contains(j, `"total":1`) {
		t.Errorf("expected total in JSON, got: %s", j)
	}
	if !strings.Contains(j, `"count":1`) {
		t.Errorf("expected count in JSON, got: %s", j)
	}
}

func TestTimelineJSONEmpty(t *testing.T) {
	j := TimelineJSON(TimelineResult{Interval: time.Minute})
	if !strings.Contains(j, `"buckets":[]`) {
		t.Errorf("expected empty buckets array, got: %s", j)
	}
}
