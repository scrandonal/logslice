package parser

import (
	"strings"
	"testing"
	"time"
)

const sampleLog = `[2024-01-15T10:00:00Z] INFO  service started
[2024-01-15T10:01:00Z] DEBUG processing request id=42
[2024-01-15T10:02:00Z] ERROR failed to connect to db
[2024-01-15T10:03:00Z] INFO  retry attempt 1
[2024-01-15T10:04:00Z] INFO  connection restored
`

func TestReadLinesAllLines(t *testing.T) {
	r := strings.NewReader(sampleLog)
	result, err := ReadLines(r, ReadOptions{Format: "raw"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(result.Lines))
	}
	if result.Collector.Total() != 5 {
		t.Errorf("expected total=5, got %d", result.Collector.Total())
	}
}

func TestReadLinesWithTimeRange(t *testing.T) {
	from := mustParseTime("2024-01-15T10:01:00Z")
	to := mustParseTime("2024-01-15T10:03:00Z")

	r := strings.NewReader(sampleLog)
	result, err := ReadLines(r, ReadOptions{
		From:   &from,
		To:     &to,
		Format: "raw",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(result.Lines))
	}
}

func TestReadLinesJSONFormat(t *testing.T) {
	r := strings.NewReader(sampleLog)
	result, err := ReadLines(r, ReadOptions{Format: "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Lines) == 0 {
		t.Fatal("expected at least one line")
	}
	for _, line := range result.Lines {
		if len(line) == 0 || line[0] != '{' {
			t.Errorf("expected JSON line, got: %s", line)
		}
	}
}

func TestReadLinesEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	result, err := ReadLines(r, ReadOptions{Format: "raw"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Lines) != 0 {
		t.Errorf("expected 0 lines, got %d", len(result.Lines))
	}
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
