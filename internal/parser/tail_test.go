package parser

import (
	"strings"
	"testing"
	"time"
)

func TestTailReaderAllLines(t *testing.T) {
	input := "[2024-01-01 10:00:00] alpha\n[2024-01-01 10:00:01] beta\n[2024-01-01 10:00:02] gamma\n"
	r := strings.NewReader(input)
	out := make(chan TailLine, 10)

	go TailReader(r, TailOptions{}, out)

	var lines []TailLine
	for tl := range out {
		lines = append(lines, tl)
	}

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0].Line.Raw != "[2024-01-01 10:00:00] alpha" {
		t.Errorf("unexpected first line: %s", lines[0].Line.Raw)
	}
	for _, tl := range lines {
		if tl.Err != nil {
			t.Errorf("unexpected error: %v", tl.Err)
		}
	}
}

func TestTailReaderMaxLines(t *testing.T) {
	input := "[2024-01-01 10:00:00] a\n[2024-01-01 10:00:01] b\n[2024-01-01 10:00:02] c\n[2024-01-01 10:00:03] d\n"
	r := strings.NewReader(input)
	out := make(chan TailLine, 10)

	go TailReader(r, TailOptions{MaxLines: 2}, out)

	var lines []TailLine
	for tl := range out {
		lines = append(lines, tl)
	}

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestTailReaderTimestampParsed(t *testing.T) {
	input := "[2024-06-15 08:30:00] startup complete\n"
	r := strings.NewReader(input)
	out := make(chan TailLine, 5)

	go TailReader(r, TailOptions{}, out)

	tl := <-out
	if tl.Line.Timestamp == nil {
		t.Fatal("expected timestamp to be parsed")
	}
	expected := time.Date(2024, 6, 15, 8, 30, 0, 0, time.UTC)
	if !tl.Line.Timestamp.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, *tl.Line.Timestamp)
	}
}

func TestTailReaderEmptyInput(t *testing.T) {
	r := strings.NewReader("")
	out := make(chan TailLine, 5)

	go TailReader(r, TailOptions{}, out)

	var lines []TailLine
	for tl := range out {
		lines = append(lines, tl)
	}
	if len(lines) != 0 {
		t.Errorf("expected no lines, got %d", len(lines))
	}
}
