package parser

import (
	"testing"
	"time"
)

func makeBurstLine(ts time.Time, raw string) Line {
	return Line{Raw: raw, Timestamp: &ts}
}

func burstTime(h, m, s int) time.Time {
	return time.Date(2024, 1, 1, h, m, s, 0, time.UTC)
}

func TestBurstDetectorEmpty(t *testing.T) {
	bd := NewBurstDetector(time.Minute, 3)
	result := bd.Detect(nil)
	if len(result) != 0 {
		t.Fatalf("expected no bursts, got %d", len(result))
	}
}

func TestBurstDetectorNoTimestamp(t *testing.T) {
	lines := []Line{
		{Raw: "no timestamp"},
		{Raw: "also none"},
	}
	bd := NewBurstDetector(time.Minute, 2)
	result := bd.Detect(lines)
	if len(result) != 0 {
		t.Fatalf("expected no bursts for lines without timestamps, got %d", len(result))
	}
}

func TestBurstDetectorBelowThreshold(t *testing.T) {
	lines := []Line{
		makeBurstLine(burstTime(10, 0, 0), "line1"),
		makeBurstLine(burstTime(10, 0, 30), "line2"),
	}
	bd := NewBurstDetector(time.Minute, 5)
	result := bd.Detect(lines)
	if len(result) != 0 {
		t.Fatalf("expected no bursts below threshold, got %d", len(result))
	}
}

func TestBurstDetectorSingleBurst(t *testing.T) {
	lines := []Line{
		makeBurstLine(burstTime(10, 0, 0), "a"),
		makeBurstLine(burstTime(10, 0, 10), "b"),
		makeBurstLine(burstTime(10, 0, 20), "c"),
		makeBurstLine(burstTime(10, 5, 0), "d"),
	}
	bd := NewBurstDetector(time.Minute, 3)
	result := bd.Detect(lines)
	if len(result) == 0 {
		t.Fatal("expected at least one burst")
	}
	got := result[0]
	if got.Count < 3 {
		t.Errorf("expected burst count >= 3, got %d", got.Count)
	}
	if !got.Start.Equal(burstTime(10, 0, 0)) {
		t.Errorf("unexpected burst start: %v", got.Start)
	}
}

func TestBurstDetectorLinesAttached(t *testing.T) {
	lines := []Line{
		makeBurstLine(burstTime(9, 0, 0), "x"),
		makeBurstLine(burstTime(9, 0, 5), "y"),
		makeBurstLine(burstTime(9, 0, 10), "z"),
	}
	bd := NewBurstDetector(30*time.Second, 3)
	result := bd.Detect(lines)
	if len(result) == 0 {
		t.Fatal("expected a burst")
	}
	if len(result[0].Lines) != 3 {
		t.Errorf("expected 3 lines in burst, got %d", len(result[0].Lines))
	}
}

func TestBurstDetectorDefaultParams(t *testing.T) {
	// zero values should not panic
	bd := NewBurstDetector(0, 0)
	if bd.windowSize != time.Minute {
		t.Errorf("expected default window 1m, got %v", bd.windowSize)
	}
	if bd.threshold != 1 {
		t.Errorf("expected default threshold 1, got %d", bd.threshold)
	}
}
