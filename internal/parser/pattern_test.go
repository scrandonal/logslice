package parser

import (
	"strings"
	"testing"
)

func makePatternLine(raw string) Line {
	return Line{Raw: raw}
}

func TestPatternCounterNoMatch(t *testing.T) {
	pc, err := NewPatternCounter(map[string]string{"error": `ERROR`})
	if err != nil {
		t.Fatal(err)
	}
	pc.Feed(makePatternLine("2024-01-01 INFO all good"))
	if got := pc.Counts()["error"]; got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
	if len(pc.Matches()) != 0 {
		t.Errorf("expected no matches")
	}
}

func TestPatternCounterMatch(t *testing.T) {
	pc, err := NewPatternCounter(map[string]string{
		"error": `ERROR`,
		"warn":  `WARN`,
	})
	if err != nil {
		t.Fatal(err)
	}
	lines := []string{
		"2024-01-01 ERROR something broke",
		"2024-01-01 INFO ok",
		"2024-01-01 WARN slow response",
		"2024-01-01 ERROR another error",
	}
	for _, raw := range lines {
		pc.Feed(makePatternLine(raw))
	}
	counts := pc.Counts()
	if counts["error"] != 2 {
		t.Errorf("expected error=2, got %d", counts["error"])
	}
	if counts["warn"] != 1 {
		t.Errorf("expected warn=1, got %d", counts["warn"])
	}
	if len(pc.Matches()) != 3 {
		t.Errorf("expected 3 matches, got %d", len(pc.Matches()))
	}
}

func TestPatternCounterInvalidRegex(t *testing.T) {
	_, err := NewPatternCounter(map[string]string{"bad": `[invalid`})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestCountPatternsHelper(t *testing.T) {
	lines := []Line{
		makePatternLine("ERROR disk full"),
		makePatternLine("INFO started"),
		makePatternLine("ERROR timeout"),
	}
	matches, counts, err := CountPatterns(lines, map[string]string{"error": `ERROR`})
	if err != nil {
		t.Fatal(err)
	}
	if counts["error"] != 2 {
		t.Errorf("expected 2, got %d", counts["error"])
	}
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d", len(matches))
	}
}

func TestPatternCountsJSON(t *testing.T) {
	counts := map[string]int{"error": 3, "warn": 1}
	out := PatternCountsJSON(counts)
	if !strings.Contains(out, `"error":3`) {
		t.Errorf("missing error count in JSON: %s", out)
	}
	if !strings.Contains(out, `"warn":1`) {
		t.Errorf("missing warn count in JSON: %s", out)
	}
}

func TestPatternMatchesJSONEmpty(t *testing.T) {
	out := PatternMatchesJSON(nil)
	if out != "[]" {
		t.Errorf("expected [], got %s", out)
	}
}
