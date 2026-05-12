package parser

import (
	"testing"
	"time"
)

func makeSeqLine(raw string, ts *time.Time) Line {
	return Line{Raw: raw, Timestamp: ts}
}

func seqTime(s string) *time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05", s)
	return &t
}

func TestSequenceDetectorEmpty(t *testing.T) {
	sd, err := NewSequenceDetector([]string{"start", "end"}, 0)
	if err != nil {
		t.Fatal(err)
	}
	result := sd.Detect(nil)
	if len(result) != 0 {
		t.Errorf("expected 0 matches, got %d", len(result))
	}
}

func TestSequenceDetectorNoMatch(t *testing.T) {
	sd, _ := NewSequenceDetector([]string{"alpha", "beta"}, 0)
	lines := []Line{
		makeSeqLine("alpha found", nil),
		makeSeqLine("gamma found", nil),
	}
	result := sd.Detect(lines)
	if len(result) != 0 {
		t.Errorf("expected 0 matches, got %d", len(result))
	}
}

func TestSequenceDetectorMatch(t *testing.T) {
	sd, _ := NewSequenceDetector([]string{"START", "PROCESS", "END"}, 0)
	lines := []Line{
		makeSeqLine("START job", seqTime("2024-01-01T10:00:00")),
		makeSeqLine("noise", nil),
		makeSeqLine("PROCESS data", seqTime("2024-01-01T10:00:05")),
		makeSeqLine("more noise", nil),
		makeSeqLine("END job", seqTime("2024-01-01T10:00:10")),
	}
	result := sd.Detect(lines)
	if len(result) != 1 {
		t.Fatalf("expected 1 match, got %d", len(result))
	}
	if len(result[0].Steps) != 3 {
		t.Errorf("expected 3 steps, got %d", len(result[0].Steps))
	}
	if result[0].Elapsed == nil {
		t.Error("expected elapsed to be set")
	} else if *result[0].Elapsed != 10*time.Second {
		t.Errorf("expected 10s elapsed, got %s", result[0].Elapsed)
	}
}

func TestSequenceDetectorWindowExcludes(t *testing.T) {
	sd, _ := NewSequenceDetector([]string{"START", "END"}, 5*time.Second)
	lines := []Line{
		makeSeqLine("START", seqTime("2024-01-01T10:00:00")),
		makeSeqLine("END", seqTime("2024-01-01T10:00:10")),
	}
	result := sd.Detect(lines)
	if len(result) != 0 {
		t.Errorf("expected 0 matches (outside window), got %d", len(result))
	}
}

func TestSequenceDetectorInvalidRegex(t *testing.T) {
	_, err := NewSequenceDetector([]string{"[invalid"}, 0)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestDetectSequencesHelper(t *testing.T) {
	lines := []Line{
		makeSeqLine("login user", nil),
		makeSeqLine("logout user", nil),
	}
	result, err := DetectSequences(lines, []string{"login", "logout"}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 match, got %d", len(result))
	}
}
