package parser

import (
	"fmt"
	"testing"
	"time"
)

func makeChunkLine(ts time.Time, msg string) string {
	return fmt.Sprintf("[%s] %s", ts.Format("2006-01-02T15:04:05Z"), msg)
}

func TestChunkSplitterEmpty(t *testing.T) {
	cs := NewChunkSplitter(time.Minute)
	chunks := cs.Flush()
	if len(chunks) != 0 {
		t.Fatalf("expected 0 chunks, got %d", len(chunks))
	}
}

func TestChunkSplitterSingleBucket(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	cs := NewChunkSplitter(time.Minute)
	for i := 0; i < 3; i++ {
		cs.Add(makeChunkLine(base.Add(time.Duration(i)*time.Second), fmt.Sprintf("msg%d", i)))
	}
	chunks := cs.Flush()
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	if len(chunks[0].Lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(chunks[0].Lines))
	}
}

func TestChunkSplitterMultipleBuckets(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	cs := NewChunkSplitter(time.Minute)
	cs.Add(makeChunkLine(base, "first"))
	cs.Add(makeChunkLine(base.Add(90*time.Second), "second bucket"))
	cs.Add(makeChunkLine(base.Add(200*time.Second), "third bucket"))
	chunks := cs.Flush()
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
}

func TestChunkSplitterNoTimestampLine(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	cs := NewChunkSplitter(time.Minute)
	cs.Add(makeChunkLine(base, "anchor"))
	cs.Add("no timestamp here")
	chunks := cs.Flush()
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	if len(chunks[0].Lines) != 2 {
		t.Errorf("expected 2 lines in chunk, got %d", len(chunks[0].Lines))
	}
}

func TestChunkSplitterNoTimestampBeforeAnchor(t *testing.T) {
	cs := NewChunkSplitter(time.Minute)
	cs.Add("orphan line") // no current chunk yet, should be dropped
	chunks := cs.Flush()
	if len(chunks) != 0 {
		t.Fatalf("expected 0 chunks, got %d", len(chunks))
	}
}

func TestChunkStartEnd(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	end := base.Add(30 * time.Second)
	cs := NewChunkSplitter(time.Minute)
	cs.Add(makeChunkLine(base, "start"))
	cs.Add(makeChunkLine(end, "end"))
	chunks := cs.Flush()
	if !chunks[0].Start.Equal(base) {
		t.Errorf("start mismatch: got %v want %v", chunks[0].Start, base)
	}
	if !chunks[0].End.Equal(end) {
		t.Errorf("end mismatch: got %v want %v", chunks[0].End, end)
	}
}
