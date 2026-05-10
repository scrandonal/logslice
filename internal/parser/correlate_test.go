package parser

import (
	"strings"
	"testing"
	"time"
)

func makeCorrelateLine(raw string, ts *time.Time) Line {
	return Line{Raw: raw, Timestamp: ts}
}

func corrTime(s string) *time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05Z", s)
	return &t
}

func TestCorrelatorNoMatch(t *testing.T) {
	c, err := NewCorrelator(`req_id=(\w+)`)
	if err != nil {
		t.Fatal(err)
	}
	c.Add(makeCorrelateLine("no id here", nil))
	res := c.Results()
	if len(res) != 1 || res[0].Key != "(unmatched)" {
		t.Fatalf("expected unmatched group, got %+v", res)
	}
}

func TestCorrelatorGrouping(t *testing.T) {
	c, err := NewCorrelator(`req=(\w+)`)
	if err != nil {
		t.Fatal(err)
	}
	c.Add(makeCorrelateLine("req=abc started", corrTime("2024-01-01T10:00:00Z")))
	c.Add(makeCorrelateLine("req=xyz started", corrTime("2024-01-01T10:01:00Z")))
	c.Add(makeCorrelateLine("req=abc finished", corrTime("2024-01-01T10:02:00Z")))
	res := c.Results()
	if len(res) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(res))
	}
	if res[0].Key != "abc" || res[0].Count != 2 {
		t.Errorf("unexpected group: %+v", res[0])
	}
	if res[1].Key != "xyz" || res[1].Count != 1 {
		t.Errorf("unexpected group: %+v", res[1])
	}
}

func TestCorrelatorTimestampRange(t *testing.T) {
	c, _ := NewCorrelator(`id=(\d+)`)
	c.Add(makeCorrelateLine("id=1 early", corrTime("2024-01-01T08:00:00Z")))
	c.Add(makeCorrelateLine("id=1 late", corrTime("2024-01-01T09:00:00Z")))
	res := c.Results()
	if res[0].First == nil || res[0].Last == nil {
		t.Fatal("expected non-nil timestamps")
	}
	if !res[0].First.Before(*res[0].Last) {
		t.Error("first should be before last")
	}
}

func TestCorrelateLinesHelper(t *testing.T) {
	lines := []Line{
		makeCorrelateLine("user=alice login", nil),
		makeCorrelateLine("user=bob login", nil),
		makeCorrelateLine("user=alice logout", nil),
	}
	res, err := CorrelateLines(lines, `user=(\w+)`)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(res))
	}
}

func TestCorrelateInvalidPattern(t *testing.T) {
	_, err := NewCorrelator(`[invalid`)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestCorrelateJSON(t *testing.T) {
	lines := []Line{
		makeCorrelateLine("req=abc ok", corrTime("2024-01-01T10:00:00Z")),
	}
	res, _ := CorrelateLines(lines, `req=(\w+)`)
	out := CorrelateJSON(res)
	if !strings.Contains(out, `"abc"`) {
		t.Errorf("expected key in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"count":1`) {
		t.Errorf("expected count in JSON, got: %s", out)
	}
}

func TestCorrelateJSONEmpty(t *testing.T) {
	out := CorrelateJSON(nil)
	if out != "[]" {
		t.Errorf("expected [], got %s", out)
	}
}
