package parser

import (
	"strings"
	"testing"
)

func buildPatternCounts() map[string]int {
	return map[string]int{"error": 5, "timeout": 2}
}

func TestPrintPatternCountsEmpty(t *testing.T) {
	var sb strings.Builder
	PrintPatternCounts(&sb, map[string]int{})
	if !strings.Contains(sb.String(), "no patterns matched") {
		t.Errorf("expected empty message, got: %s", sb.String())
	}
}

func TestPrintPatternCounts(t *testing.T) {
	var sb strings.Builder
	PrintPatternCounts(&sb, buildPatternCounts())
	out := sb.String()
	if !strings.Contains(out, "error") {
		t.Errorf("expected 'error' in output: %s", out)
	}
	if !strings.Contains(out, "5") {
		t.Errorf("expected count 5 in output: %s", out)
	}
	if !strings.Contains(out, "timeout") {
		t.Errorf("expected 'timeout' in output: %s", out)
	}
}

func TestPrintPatternMatches(t *testing.T) {
	matches := []PatternMatch{
		{Line: Line{Raw: "ERROR disk full"}, Pattern: "error", Count: 1},
		{Line: Line{Raw: "ERROR timeout"}, Pattern: "error", Count: 2},
	}
	var sb strings.Builder
	PrintPatternMatches(&sb, matches)
	out := sb.String()
	if !strings.Contains(out, "[error]") {
		t.Errorf("expected pattern prefix, got: %s", out)
	}
	if !strings.Contains(out, "disk full") {
		t.Errorf("expected raw line content, got: %s", out)
	}
}

func TestPatternCountsJSONEmpty(t *testing.T) {
	out := PatternCountsJSON(map[string]int{})
	if out != "{}" {
		t.Errorf("expected {}, got %s", out)
	}
}

func TestPatternMatchesJSONStructure(t *testing.T) {
	matches := []PatternMatch{
		{Line: Line{Raw: "ERROR something"}, Pattern: "error", Count: 1},
	}
	out := PatternMatchesJSON(matches)
	if !strings.HasPrefix(out, "[") || !strings.HasSuffix(out, "]") {
		t.Errorf("expected JSON array, got: %s", out)
	}
	if !strings.Contains(out, `"pattern"`) {
		t.Errorf("expected pattern key in JSON: %s", out)
	}
	if !strings.Contains(out, `"raw"`) {
		t.Errorf("expected raw key in JSON: %s", out)
	}
	if !strings.Contains(out, `"count"`) {
		t.Errorf("expected count key in JSON: %s", out)
	}
}
