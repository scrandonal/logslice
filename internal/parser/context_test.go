package parser

import (
	"testing"
)

func makeCtxLine(raw string) Line {
	return Line{Raw: raw}
}

func TestContextExtractorNoMatch(t *testing.T) {
	ex := NewContextExtractor(1, 1)
	lines := []Line{makeCtxLine("aaa"), makeCtxLine("bbb")}
	var out []ContextLine
	for _, l := range lines {
		out = append(out, ex.Push(l, false)...)
	}
	out = append(out, ex.Flush()...)
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}

func TestContextExtractorMatchNoBefore(t *testing.T) {
	ex := NewContextExtractor(0, 0)
	out := ex.Push(makeCtxLine("match"), true)
	if len(out) != 1 || out[0].Line.Raw != "match" {
		t.Fatalf("unexpected result: %+v", out)
	}
}

func TestContextExtractorBeforeAndAfter(t *testing.T) {
	ex := NewContextExtractor(1, 1)
	lines := []string{"before", "match", "after", "extra"}
	var out []ContextLine
	for i, raw := range lines {
		out = append(out, ex.Push(makeCtxLine(raw), i == 1)...)
	}
	out = append(out, ex.Flush()...)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	r := out[0]
	if len(r.Before) != 1 || r.Before[0].Raw != "before" {
		t.Errorf("wrong before: %+v", r.Before)
	}
	if len(r.After) != 1 || r.After[0].Raw != "after" {
		t.Errorf("wrong after: %+v", r.After)
	}
}

func TestExtractContextHelper(t *testing.T) {
	lines := []Line{
		makeCtxLine("line one"),
		makeCtxLine("ERROR occurred"),
		makeCtxLine("line three"),
		makeCtxLine("all good"),
		makeCtxLine("ERROR again"),
	}
	out := ExtractContext(lines, "ERROR", 1, 1, false)
	if len(out) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(out))
	}
	if out[0].Line.Raw != "ERROR occurred" {
		t.Errorf("wrong first match: %s", out[0].Line.Raw)
	}
	if out[1].Line.Raw != "ERROR again" {
		t.Errorf("wrong second match: %s", out[1].Line.Raw)
	}
}

func TestExtractContextCaseInsensitive(t *testing.T) {
	lines := []Line{
		makeCtxLine("warn: disk low"),
		makeCtxLine("WARN: memory low"),
	}
	out := ExtractContext(lines, "warn", 0, 0, true)
	if len(out) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(out))
	}
}

func TestContextJSONEmpty(t *testing.T) {
	got := ContextJSON(nil)
	if got != "[]" {
		t.Errorf("expected [], got %s", got)
	}
}

func TestContextJSON(t *testing.T) {
	results := []ContextLine{
		{
			Line:    makeCtxLine("matched line"),
			Before:  []Line{makeCtxLine("before line")},
			After:   []Line{makeCtxLine("after line")},
			Matched: true,
		},
	}
	got := ContextJSON(results)
	if got == "" || got == "[]" {
		t.Errorf("unexpected empty JSON: %s", got)
	}
	if len(got) < 10 {
		t.Errorf("JSON too short: %s", got)
	}
}
