package parser

import (
	"testing"
	"time"
)

func TestCollectorInitialState(t *testing.T) {
	c := NewCollector()
	s := c.Finalize()

	if s.LinesScanned != 0 {
		t.Errorf("expected 0 scanned, got %d", s.LinesScanned)
	}
	if s.LinesMatched != 0 {
		t.Errorf("expected 0 matched, got %d", s.LinesMatched)
	}
	if s.FirstMatch != nil || s.LastMatch != nil {
		t.Error("expected nil first/last match")
	}
}

func TestCollectorCounting(t *testing.T) {
	c := NewCollector()

	for i := 0; i < 10; i++ {
		c.RecordScanned()
	}
	for i := 0; i < 3; i++ {
		c.RecordSkipped()
	}
	c.RecordParseError()

	s := c.Finalize()
	if s.LinesScanned != 10 {
		t.Errorf("expected 10 scanned, got %d", s.LinesScanned)
	}
	if s.LinesSkipped != 3 {
		t.Errorf("expected 3 skipped, got %d", s.LinesSkipped)
	}
	if s.ParseErrors != 1 {
		t.Errorf("expected 1 parse error, got %d", s.ParseErrors)
	}
}

func TestCollectorMatchTimestamps(t *testing.T) {
	c := NewCollector()

	t1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	c.RecordMatch(&t1)
	c.RecordMatch(&t2)
	c.RecordMatch(&t3)

	s := c.Finalize()
	if s.LinesMatched != 3 {
		t.Errorf("expected 3 matched, got %d", s.LinesMatched)
	}
	if s.FirstMatch == nil || !s.FirstMatch.Equal(t1) {
		t.Errorf("expected FirstMatch=%v, got %v", t1, s.FirstMatch)
	}
	if s.LastMatch == nil || !s.LastMatch.Equal(t3) {
		t.Errorf("expected LastMatch=%v, got %v", t3, s.LastMatch)
	}
}

func TestCollectorElapsed(t *testing.T) {
	c := NewCollector()
	time.Sleep(5 * time.Millisecond)
	s := c.Finalize()
	if s.Elapsed < 5*time.Millisecond {
		t.Errorf("expected elapsed >= 5ms, got %v", s.Elapsed)
	}
}

func TestCollectorMatchNilTimestamp(t *testing.T) {
	c := NewCollector()
	c.RecordMatch(nil)
	s := c.Finalize()
	if s.LinesMatched != 1 {
		t.Errorf("expected 1 matched, got %d", s.LinesMatched)
	}
	if s.FirstMatch != nil || s.LastMatch != nil {
		t.Error("expected nil first/last match when ts is nil")
	}
}
