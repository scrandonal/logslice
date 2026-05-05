package parser

import (
	"strings"
	"testing"
)

func mustFilter(t *testing.T, from, to string) Filter {
	t.Helper()
	f, err := NewFilter(from, to)
	if err != nil {
		t.Fatalf("NewFilter(%q, %q): %v", from, to, err)
	}
	return f
}

func TestFilterMatch(t *testing.T) {
	f := mustFilter(t, "2024-01-01T10:00:00Z", "2024-01-01T12:00:00Z")

	cases := []struct {
		ts   string
		want bool
	}{
		{"2024-01-01T09:59:59Z", false},
		{"2024-01-01T10:00:00Z", true},
		{"2024-01-01T11:00:00Z", true},
		{"2024-01-01T12:00:00Z", true},
		{"2024-01-01T12:00:01Z", false},
	}
	for _, c := range cases {
		ts, err := ParseTimestamp(c.ts)
		if err != nil {
			t.Fatalf("ParseTimestamp(%q): %v", c.ts, err)
		}
		if got := f.Match(ts); got != c.want {
			t.Errorf("Match(%s) = %v, want %v", c.ts, got, c.want)
		}
	}
}

func TestFilterSlice(t *testing.T) {
	input := strings.Join([]string{
		"[2024-01-01T09:00:00Z] before range",
		"[2024-01-01T10:00:00Z] first match",
		"[2024-01-01T11:00:00Z] second match",
		"[2024-01-01T12:00:00Z] third match",
		"[2024-01-01T13:00:00Z] after range",
	}, "\n")

	f := mustFilter(t, "2024-01-01T10:00:00Z", "2024-01-01T12:00:00Z")
	s := NewScanner(strings.NewReader(input))
	lines := f.Slice(s)

	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3: %v", len(lines), lines)
	}
	if !strings.Contains(lines[0], "first match") {
		t.Errorf("unexpected first line: %q", lines[0])
	}
	if !strings.Contains(lines[2], "third match") {
		t.Errorf("unexpected last line: %q", lines[2])
	}
}

func TestFilterSliceEmpty(t *testing.T) {
	f := mustFilter(t, "2024-01-01T10:00:00Z", "2024-01-01T12:00:00Z")
	s := NewScanner(strings.NewReader(""))
	lines := f.Slice(s)
	if len(lines) != 0 {
		t.Errorf("expected empty result, got %v", lines)
	}
}
