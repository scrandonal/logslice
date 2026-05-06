package parser

import (
	"testing"
	"time"
)

func makeSampleLine(ts *time.Time, raw string) Line {
	return Line{Timestamp: ts, Raw: raw}
}

func ptrTime(t time.Time) *time.Time { return &t }

func TestSamplerKeepAll(t *testing.T) {
	s := NewSampler(SampleConfig{Rate: 1})
	lines := []Line{
		makeSampleLine(nil, "a"),
		makeSampleLine(nil, "b"),
		makeSampleLine(nil, "c"),
	}
	out := s.Sample(lines)
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
}

func TestSamplerRateTwo(t *testing.T) {
	s := NewSampler(SampleConfig{Rate: 2})
	lines := []Line{
		makeSampleLine(nil, "1"),
		makeSampleLine(nil, "2"),
		makeSampleLine(nil, "3"),
		makeSampleLine(nil, "4"),
	}
	out := s.Sample(lines)
	// counter hits 2,4 → lines[1] and lines[3]
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	if out[0].Raw != "2" || out[1].Raw != "4" {
		t.Errorf("unexpected lines: %v", out)
	}
}

func TestSamplerBucket(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	s := NewSampler(SampleConfig{Rate: 1, Bucket: time.Minute})

	lines := []Line{
		makeSampleLine(ptrTime(base), "first in bucket"),
		makeSampleLine(ptrTime(base.Add(30*time.Second)), "second in same bucket"),
		makeSampleLine(ptrTime(base.Add(90*time.Second)), "first in next bucket"),
	}
	out := s.Sample(lines)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	if out[0].Raw != "first in bucket" {
		t.Errorf("unexpected first line: %s", out[0].Raw)
	}
	if out[1].Raw != "first in next bucket" {
		t.Errorf("unexpected second line: %s", out[1].Raw)
	}
}

func TestSamplerReset(t *testing.T) {
	s := NewSampler(SampleConfig{Rate: 2})
	lines := []Line{
		makeSampleLine(nil, "a"),
		makeSampleLine(nil, "b"),
	}
	s.Sample(lines) // counter now at 2
	s.Reset()
	out := s.Sample(lines) // counter restarts: hits 2 → "b"
	if len(out) != 1 || out[0].Raw != "b" {
		t.Errorf("expected ['b'] after reset, got %v", out)
	}
}

func TestSamplerZeroRateTreatedAsOne(t *testing.T) {
	s := NewSampler(SampleConfig{Rate: 0})
	lines := []Line{
		makeSampleLine(nil, "x"),
		makeSampleLine(nil, "y"),
	}
	out := s.Sample(lines)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines with rate=0, got %d", len(out))
	}
}
