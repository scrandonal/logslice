package parser

import (
	"testing"
	"time"
)

func makeWindowLine(raw string, ts *time.Time) Line {
	return Line{Raw: raw, Timestamp: ts}
}

func ptr(t time.Time) *time.Time { return &t }

func TestWindowEmpty(t *testing.T) {
	w := NewWindow(5)
	if w.Len() != 0 {
		t.Errorf("expected empty window")
	}
	if len(w.Lines()) != 0 {
		t.Errorf("expected no lines")
	}
}

func TestWindowPushAndOrder(t *testing.T) {
	w := NewWindow(3)
	w.Push(makeWindowLine("a", nil))
	w.Push(makeWindowLine("b", nil))
	w.Push(makeWindowLine("c", nil))

	lines := w.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3, got %d", len(lines))
	}
	if lines[0].Raw != "a" || lines[1].Raw != "b" || lines[2].Raw != "c" {
		t.Errorf("wrong order: %v", lines)
	}
}

func TestWindowEviction(t *testing.T) {
	w := NewWindow(2)
	w.Push(makeWindowLine("a", nil))
	w.Push(makeWindowLine("b", nil))
	w.Push(makeWindowLine("c", nil))

	lines := w.Lines()
	if len(lines) != 2 {
		t.Fatalf("expected 2, got %d", len(lines))
	}
	if lines[0].Raw != "b" || lines[1].Raw != "c" {
		t.Errorf("expected b,c got %v", lines)
	}
}

func TestWindowInRange(t *testing.T) {
	w := NewWindow(10)
	t1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC)
	t3 := time.Date(2024, 1, 1, 10, 10, 0, 0, time.UTC)

	w.Push(makeWindowLine("early", ptr(t1)))
	w.Push(makeWindowLine("mid", ptr(t2)))
	w.Push(makeWindowLine("late", ptr(t3)))
	w.Push(makeWindowLine("no-ts", nil))

	from := time.Date(2024, 1, 1, 10, 3, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 10, 7, 0, 0, time.UTC)

	result := w.InRange(from, to)
	if len(result) != 2 {
		t.Fatalf("expected 2 (mid + no-ts), got %d", len(result))
	}
	if result[0].Raw != "mid" {
		t.Errorf("expected mid, got %s", result[0].Raw)
	}
	if result[1].Raw != "no-ts" {
		t.Errorf("expected no-ts, got %s", result[1].Raw)
	}
}
