package parser

import (
	"strings"
	"testing"
)

func TestTruncateNoop(t *testing.T) {
	tr := NewTruncator(0, TruncateRight)
	s := "this string should not be changed"
	if got := tr.Truncate(s); got != s {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestTruncateRight(t *testing.T) {
	tr := NewTruncator(10, TruncateRight)
	input := "hello world, this is a long line"
	got := tr.Truncate(input)
	if len(got) != 10 {
		t.Fatalf("expected length 10, got %d: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected ellipsis suffix, got %q", got)
	}
}

func TestTruncateRightShortInput(t *testing.T) {
	tr := NewTruncator(100, TruncateRight)
	input := "short"
	if got := tr.Truncate(input); got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestTruncateMiddle(t *testing.T) {
	tr := NewTruncator(11, TruncateMiddle)
	input := "abcdefghijklmnopqrstuvwxyz"
	got := tr.Truncate(input)
	if len(got) != 11 {
		t.Fatalf("expected length 11, got %d: %q", len(got), got)
	}
	if !strings.Contains(got, "...") {
		t.Fatalf("expected ellipsis in middle, got %q", got)
	}
	// prefix and suffix should come from original
	if !strings.HasPrefix(got, "abc") {
		t.Fatalf("expected prefix 'abc', got %q", got)
	}
	if !strings.HasSuffix(got, "xyz") {
		t.Fatalf("expected suffix 'xyz', got %q", got)
	}
}

func TestTruncateMiddleExact(t *testing.T) {
	tr := NewTruncator(5, TruncateMiddle)
	input := "12345"
	if got := tr.Truncate(input); got != input {
		t.Fatalf("expected unchanged, got %q", got)
	}
}

func TestTruncateLinesHelper(t *testing.T) {
	lines := []Line{
		{Raw: "short line"},
		{Raw: "this is a very long line that exceeds the limit we set"},
	}
	out := TruncateLines(lines, 20, TruncateRight)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	if len(out[0].Raw) > 20 {
		t.Fatalf("line 0 should be unchanged: %q", out[0].Raw)
	}
	if len(out[1].Raw) != 20 {
		t.Fatalf("line 1 should be truncated to 20, got %d: %q", len(out[1].Raw), out[1].Raw)
	}
}

func TestTruncateLinesDisabled(t *testing.T) {
	lines := []Line{
		{Raw: strings.Repeat("x", 200)},
	}
	out := TruncateLines(lines, 0, TruncateRight)
	if len(out[0].Raw) != 200 {
		t.Fatalf("expected 200, got %d", len(out[0].Raw))
	}
}
