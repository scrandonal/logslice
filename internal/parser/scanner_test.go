package parser

import (
	"strings"
	"testing"
	"time"
)

const sampleLog = `[2024-01-15T10:00:00Z] INFO  server started on :8080
[2024-01-15T10:01:23Z] DEBUG received request GET /health
not a valid log line
[2024-01-15T10:02:45Z] ERROR connection refused: dial tcp 127.0.0.1:5432
[2024-01-15T10:03:10Z] INFO  request completed in 12ms
`

func TestScannerBasic(t *testing.T) {
	r := strings.NewReader(sampleLog)
	s := NewScanner(r)

	var entries []Entry
	for s.Scan() {
		entries = append(entries, s.Entry())
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected scanner error: %v", err)
	}

	// Should skip the invalid line and return 4 valid entries.
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}
}

func TestScannerTimestamps(t *testing.T) {
	r := strings.NewReader(sampleLog)
	s := NewScanner(r)

	expected := []string{
		"2024-01-15T10:00:00Z",
		"2024-01-15T10:01:23Z",
		"2024-01-15T10:02:45Z",
		"2024-01-15T10:03:10Z",
	}

	i := 0
	for s.Scan() {
		entry := s.Entry()
		want, err := time.Parse(time.RFC3339, expected[i])
		if err != nil {
			t.Fatalf("failed to parse expected timestamp: %v", err)
		}
		if !entry.Timestamp.Equal(want) {
			t.Errorf("entry %d: expected timestamp %v, got %v", i, want, entry.Timestamp)
		}
		i++
	}
}

func TestScannerEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	s := NewScanner(r)

	if s.Scan() {
		t.Fatal("expected no entries for empty input")
	}
	if s.Err() != nil {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestScannerNoValidLines(t *testing.T) {
	input := "no timestamps here\njust plain text\n"
	r := strings.NewReader(input)
	s := NewScanner(r)

	if s.Scan() {
		t.Fatal("expected no entries when no valid timestamps present")
	}
}
