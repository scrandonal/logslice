package parser

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestFormatterRaw(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatRaw)

	line := []byte("[2024-01-15T10:00:00Z] INFO hello world")
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	if err := f.WriteLine(line, ts); err != nil {
		t.Fatalf("WriteLine error: %v", err)
	}

	got := buf.String()
	want := string(line) + "\n"
	if got != want {
		t.Errorf("raw: got %q, want %q", got, want)
	}
}

func TestFormatterJSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON)

	line := []byte("[2024-01-15T10:00:00Z] INFO hello world")
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	if err := f.WriteLine(line, ts); err != nil {
		t.Fatalf("WriteLine error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"timestamp":"2024-01-15T10:00:00Z"`) {
		t.Errorf("JSON missing timestamp field: %s", got)
	}
	if !strings.Contains(got, `"line":`) {
		t.Errorf("JSON missing line field: %s", got)
	}
	if !strings.HasSuffix(strings.TrimSpace(got), "}") {
		t.Errorf("JSON output not closed properly: %s", got)
	}
}

func TestFormatterJSONNullTimestamp(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON)

	if err := f.WriteLine([]byte("no timestamp here"), time.Time{}); err != nil {
		t.Fatalf("WriteLine error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"timestamp":null`) {
		t.Errorf("expected null timestamp, got: %s", got)
	}
}

func TestFormatterJSONEscaping(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON)

	line := []byte(`tab:	here "quoted" back\slash`)
	if err := f.WriteLine(line, time.Time{}); err != nil {
		t.Fatalf("WriteLine error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `\t`) {
		t.Errorf("expected escaped tab in output: %s", got)
	}
	if !strings.Contains(got, `\"`) {
		t.Errorf("expected escaped quote in output: %s", got)
	}
	if !strings.Contains(got, `\\`) {
		t.Errorf("expected escaped backslash in output: %s", got)
	}
}

func TestFormatterMultipleLines(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatRaw)

	lines := []string{"line one", "line two", "line three"}
	for _, l := range lines {
		if err := f.WriteLine([]byte(l), time.Time{}); err != nil {
			t.Fatalf("WriteLine error: %v", err)
		}
	}

	got := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(got) != len(lines) {
		t.Fatalf("expected %d lines, got %d", len(lines), len(got))
	}
	for i, want := range lines {
		if got[i] != want {
			t.Errorf("line %d: got %q, want %q", i, got[i], want)
		}
	}
}
