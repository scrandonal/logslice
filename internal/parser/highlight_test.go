package parser

import (
	"testing"
)

func TestHighlightNoneMode(t *testing.T) {
	h := NewHighlighter([]string{"error"}, HighlightNone)
	input := "an error occurred"
	if got := h.Highlight(input); got != input {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestHighlightNoKeywords(t *testing.T) {
	h := NewHighlighter(nil, HighlightBracket)
	input := "nothing to highlight"
	if got := h.Highlight(input); got != input {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestHighlightBracketSingle(t *testing.T) {
	h := NewHighlighter([]string{"error"}, HighlightBracket)
	got := h.Highlight("an error occurred")
	want := "an [[error]] occurred"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHighlightBracketMultipleOccurrences(t *testing.T) {
	h := NewHighlighter([]string{"warn"}, HighlightBracket)
	got := h.Highlight("warn: disk warn threshold")
	want := "[[warn]]: disk [[warn]] threshold"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHighlightCaseInsensitive(t *testing.T) {
	h := NewHighlighter([]string{"error"}, HighlightBracket)
	got := h.Highlight("An ERROR occurred")
	want := "An [[ERROR]] occurred"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHighlightMultipleKeywords(t *testing.T) {
	h := NewHighlighter([]string{"error", "warn"}, HighlightBracket)
	got := h.Highlight("error and warn both present")
	want := "[[error]] and [[warn]] both present"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHighlightANSI(t *testing.T) {
	h := NewHighlighter([]string{"ok"}, HighlightANSI)
	got := h.Highlight("status ok")
	want := "status \x1b[1;33mok\x1b[0m"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHighlightNoMatch(t *testing.T) {
	h := NewHighlighter([]string{"fatal"}, HighlightBracket)
	input := "everything is fine"
	if got := h.Highlight(input); got != input {
		t.Errorf("expected unchanged line, got %q", got)
	}
}
