package parser

import (
	"strings"
	"testing"
	"time"
)

func makeClassifyLine(raw string) Line {
	now := time.Now()
	return Line{Raw: raw, Timestamp: &now}
}

func TestClassifierNoRules(t *testing.T) {
	c, err := NewClassifier(nil, "other")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	l := makeClassifyLine("some log line")
	r := c.Classify(l)
	if r.Category != "other" {
		t.Errorf("expected 'other', got %q", r.Category)
	}
}

func TestClassifierMatch(t *testing.T) {
	rules := map[string]string{
		"error":   `\berror\b`,
		"warning": `\bwarn(ing)?\b`,
	}
	c, err := NewClassifier(rules, "info")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tests := []struct {
		raw  string
		want string
	}{
		{"2024/01/01 ERROR something failed", "error"},
		{"2024/01/01 Warning: disk low", "warning"},
		{"2024/01/01 INFO startup complete", "info"},
	}
	for _, tc := range tests {
		r := c.Classify(makeClassifyLine(tc.raw))
		if r.Category != tc.want {
			t.Errorf("raw=%q: expected category %q, got %q", tc.raw, tc.want, r.Category)
		}
	}
}

func TestClassifierInvalidPattern(t *testing.T) {
	_, err := NewClassifier(map[string]string{"bad": "[invalid"}, "")
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestClassifyLinesHelper(t *testing.T) {
	lines := []Line{
		makeClassifyLine("ERROR disk full"),
		makeClassifyLine("INFO all good"),
	}
	rules := map[string]string{"error": `error`}
	results, err := ClassifyLines(lines, rules, "other")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Category != "error" {
		t.Errorf("expected 'error', got %q", results[0].Category)
	}
	if results[1].Category != "other" {
		t.Errorf("expected 'other', got %q", results[1].Category)
	}
}

func TestCategorySummary(t *testing.T) {
	results := []ClassifyResult{
		{Line: makeClassifyLine("a"), Category: "error"},
		{Line: makeClassifyLine("b"), Category: "error"},
		{Line: makeClassifyLine("c"), Category: "info"},
	}
	counts := CategorySummary(results)
	if counts["error"] != 2 {
		t.Errorf("expected error count 2, got %d", counts["error"])
	}
	if counts["info"] != 1 {
		t.Errorf("expected info count 1, got %d", counts["info"])
	}
}

func TestClassifySummaryJSON(t *testing.T) {
	results := []ClassifyResult{
		{Line: makeClassifyLine("x"), Category: "error"},
	}
	out := ClassifySummaryJSON(results)
	if !strings.Contains(out, `"error"`) {
		t.Errorf("expected 'error' in JSON, got %s", out)
	}
}

func TestClassifyJSON(t *testing.T) {
	results := []ClassifyResult{
		{Line: makeClassifyLine("msg"), Category: "warn"},
	}
	out := ClassifyJSON(results)
	if !strings.Contains(out, `"warn"`) {
		t.Errorf("expected category in JSON output, got %s", out)
	}
	if !strings.Contains(out, `"msg"`) {
		t.Errorf("expected raw in JSON output, got %s", out)
	}
}
