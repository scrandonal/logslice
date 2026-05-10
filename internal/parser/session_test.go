package parser

import (
	"strings"
	"testing"
	"time"
)

func makeSessionLine(raw string, ts *time.Time) Line {
	return Line{Raw: raw, Timestamp: ts}
}

func sessionTime(s string) *time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05", s)
	return &t
}

func TestSplitSessionsEmpty(t *testing.T) {
	result := SplitSessions(nil, time.Minute)
	if len(result) != 0 {
		t.Fatalf("expected 0 sessions, got %d", len(result))
	}
}

func TestSplitSessionsSingleSession(t *testing.T) {
	lines := []Line{
		makeSessionLine("a", sessionTime("2024-01-01T10:00:00")),
		makeSessionLine("b", sessionTime("2024-01-01T10:00:30")),
		makeSessionLine("c", sessionTime("2024-01-01T10:01:00")),
	}
	sessions := SplitSessions(lines, 2*time.Minute)
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if len(sessions[0].Lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(sessions[0].Lines))
	}
	if sessions[0].Duration != time.Minute {
		t.Errorf("expected 1m duration, got %s", sessions[0].Duration)
	}
}

func TestSplitSessionsMultipleSessions(t *testing.T) {
	lines := []Line{
		makeSessionLine("a", sessionTime("2024-01-01T10:00:00")),
		makeSessionLine("b", sessionTime("2024-01-01T10:00:30")),
		makeSessionLine("c", sessionTime("2024-01-01T10:10:00")), // gap > 5m
		makeSessionLine("d", sessionTime("2024-01-01T10:10:10")),
	}
	sessions := SplitSessions(lines, 5*time.Minute)
	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
	if len(sessions[0].Lines) != 2 {
		t.Errorf("session 0: expected 2 lines, got %d", len(sessions[0].Lines))
	}
	if len(sessions[1].Lines) != 2 {
		t.Errorf("session 1: expected 2 lines, got %d", len(sessions[1].Lines))
	}
}

func TestSplitSessionsNoTimestamp(t *testing.T) {
	lines := []Line{
		makeSessionLine("a", nil),
		makeSessionLine("b", nil),
	}
	sessions := SplitSessions(lines, time.Minute)
	// No timestamps means no gap detection; all in one session.
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Start != nil || sessions[0].End != nil {
		t.Errorf("expected nil start/end for no-timestamp session")
	}
}

func TestSessionsJSONEmpty(t *testing.T) {
	out := SessionsJSON(nil)
	if out != "[]" {
		t.Errorf("expected [], got %s", out)
	}
}

func TestSessionsJSONStructure(t *testing.T) {
	lines := []Line{
		makeSessionLine("x", sessionTime("2024-06-01T08:00:00")),
		makeSessionLine("y", sessionTime("2024-06-01T08:00:05")),
	}
	sessions := SplitSessions(lines, time.Minute)
	out := SessionsJSON(sessions)
	if !strings.Contains(out, `"session":1`) {
		t.Errorf("missing session field: %s", out)
	}
	if !strings.Contains(out, `"lines":2`) {
		t.Errorf("missing lines field: %s", out)
	}
	if !strings.Contains(out, `"duration_ms":5000`) {
		t.Errorf("missing duration_ms: %s", out)
	}
}
