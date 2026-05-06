package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func buildChunks(t *testing.T) []*Chunk {
	t.Helper()
	base := time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC)
	cs := NewChunkSplitter(time.Minute)
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			ts := base.Add(time.Duration(i)*2*time.Minute + time.Duration(j)*time.Second)
			cs.Add(fmt.Sprintf("[%s] msg %d-%d", ts.Format("2006-01-02T15:04:05Z"), i, j))
		}
	}
	return cs.Flush()
}

func TestPrintChunkSummaries(t *testing.T) {
	chunks := buildChunks(t)
	var buf bytes.Buffer
	PrintChunkSummaries(&buf, chunks)
	out := buf.String()
	if !strings.Contains(out, "start") {
		t.Error("expected header 'start' in output")
	}
	if !strings.Contains(out, "2024-06-01T08:00:00Z") {
		t.Errorf("expected first chunk start timestamp in output, got:\n%s", out)
	}
}

func TestChunksJSON(t *testing.T) {
	chunks := buildChunks(t)
	raw, err := ChunksJSON(chunks)
	if err != nil {
		t.Fatalf("ChunksJSON error: %v", err)
	}
	var summaries []ChunkSummary
	if err := json.Unmarshal([]byte(raw), &summaries); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}
	if summaries[0].LineCount != 3 {
		t.Errorf("expected 3 lines in first chunk, got %d", summaries[0].LineCount)
	}
}

func TestChunksJSONEmpty(t *testing.T) {
	raw, err := ChunksJSON(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if raw != "[]" {
		t.Errorf("expected '[]', got %q", raw)
	}
}

func TestPrintChunkSummariesEmpty(t *testing.T) {
	var buf bytes.Buffer
	PrintChunkSummaries(&buf, nil)
	out := buf.String()
	if !strings.Contains(out, "start") {
		t.Error("expected header even for empty chunk list")
	}
}
