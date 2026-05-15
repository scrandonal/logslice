package parser

import (
	"testing"
	"time"
)

func makeRetraceLine(raw string, ts *time.Time) Line {
	return Line{Raw: raw, Timestamp: ts}
}

func TestRetraceEmpty(t *testing.T) {
	result := RetraceLines(nil, RetraceOptions{})
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestRetraceNoFrames(t *testing.T) {
	lines := []Line{
		makeRetraceLine("2024-01-01 INFO starting server", nil),
		makeRetraceLine("2024-01-01 INFO listening on :8080", nil),
	}
	result := RetraceLines(lines, RetraceOptions{})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	for _, e := range result {
		if len(e.Frames) != 0 {
			t.Errorf("expected no frames, got %d", len(e.Frames))
		}
	}
}

func TestRetraceJavaStyle(t *testing.T) {
	lines := []Line{
		makeRetraceLine("ERROR NullPointerException", nil),
		makeRetraceLine("\tat com.example.Foo.bar(Foo.java:42)", nil),
		makeRetraceLine("\tat com.example.Main.main(Main.java:10)", nil),
		makeRetraceLine("INFO recovered", nil),
	}
	result := RetraceLines(lines, RetraceOptions{})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if len(result[0].Frames) != 2 {
		t.Errorf("expected 2 frames, got %d", len(result[0].Frames))
	}
	if result[0].Trigger.Raw != "ERROR NullPointerException" {
		t.Errorf("unexpected trigger: %s", result[0].Trigger.Raw)
	}
}

func TestRetraceMaxFrames(t *testing.T) {
	lines := []Line{
		makeRetraceLine("FATAL panic: runtime error", nil),
		makeRetraceLine("\tgoroutine 1 [running]:", nil),
		makeRetraceLine("\tmain.doWork(main.go:55)", nil),
		makeRetraceLine("\tmain.main(main.go:10)", nil),
	}
	result := RetraceLines(lines, RetraceOptions{MaxFrames: 2})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if len(result[0].Frames) != 2 {
		t.Errorf("expected 2 frames (capped), got %d", len(result[0].Frames))
	}
}

func TestRetraceCustomPattern(t *testing.T) {
	lines := []Line{
		makeRetraceLine("ERROR something failed", nil),
		makeRetraceLine("  >> module/pkg/file.go:99", nil),
		makeRetraceLine("  >> module/pkg/other.go:12", nil),
	}
	result := RetraceLines(lines, RetraceOptions{FramePattern: `^\s+>>`})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if len(result[0].Frames) != 2 {
		t.Errorf("expected 2 frames, got %d", len(result[0].Frames))
	}
}

func TestRetraceInvalidPatternFallback(t *testing.T) {
	lines := []Line{
		makeRetraceLine("WARN something", nil),
		makeRetraceLine("\tat pkg.Foo(foo.go:1)", nil),
	}
	// invalid regex should fall back to default
	result := RetraceLines(lines, RetraceOptions{FramePattern: "[invalid"})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry with fallback, got %d", len(result))
	}
	if len(result[0].Frames) != 1 {
		t.Errorf("expected 1 frame with fallback pattern, got %d", len(result[0].Frames))
	}
}
