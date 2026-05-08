package parser

import (
	"bytes"
	"strings"
	"testing"
)

func buildContextResults() []ContextLine {
	return []ContextLine{
		{
			Line:    makeCtxLine("ERROR: disk full"),
			Before:  []Line{makeCtxLine("checking disk")},
			After:   []Line{makeCtxLine("cleanup started")},
			Matched: true,
		},
	}
}

func TestPrintContextEmpty(t *testing.T) {
	var buf bytes.Buffer
	PrintContext(&buf, nil)
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got: %s", buf.String())
	}
}

func TestPrintContextSingle(t *testing.T) {
	var buf bytes.Buffer
	PrintContext(&buf, buildContextResults())
	out := buf.String()
	if !strings.Contains(out, "> ERROR: disk full") {
		t.Errorf("missing matched line marker: %s", out)
	}
	if !strings.Contains(out, "checking disk") {
		t.Errorf("missing before line: %s", out)
	}
	if !strings.Contains(out, "cleanup started") {
		t.Errorf("missing after line: %s", out)
	}
}

func TestPrintContextSeparator(t *testing.T) {
	results := []ContextLine{
		{Line: makeCtxLine("first match"), Matched: true},
		{Line: makeCtxLine("second match"), Matched: true},
	}
	var buf bytes.Buffer
	PrintContext(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "--") {
		t.Errorf("expected separator between results: %s", out)
	}
}

func TestContextJSONStructure(t *testing.T) {
	results := buildContextResults()
	got := ContextJSON(results)
	if !strings.HasPrefix(got, "[") || !strings.HasSuffix(got, "]") {
		t.Errorf("not a JSON array: %s", got)
	}
	if !strings.Contains(got, `"matched"`) {
		t.Errorf("missing matched field: %s", got)
	}
	if !strings.Contains(got, `"before"`) {
		t.Errorf("missing before field: %s", got)
	}
	if !strings.Contains(got, `"after"`) {
		t.Errorf("missing after field: %s", got)
	}
}

func TestContextJSONMultiple(t *testing.T) {
	results := []ContextLine{
		{Line: makeCtxLine("match one"), Matched: true},
		{Line: makeCtxLine("match two"), Matched: true},
	}
	got := ContextJSON(results)
	count := strings.Count(got, `"matched"`)
	if count != 2 {
		t.Errorf("expected 2 matched fields, got %d in: %s", count, got)
	}
}
