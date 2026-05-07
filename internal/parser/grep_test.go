package parser

import (
	"testing"
)

func makeGrepLine(raw string) Line {
	return Line{Raw: raw}
}

func TestGrepLiteralMatch(t *testing.T) {
	g, err := NewGrepper("ERROR", GrepLiteral)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !g.Match("[2024-01-01] ERROR something broke") {
		t.Error("expected match for line containing ERROR")
	}
	if g.Match("[2024-01-01] INFO all good") {
		t.Error("expected no match for line without ERROR")
	}
}

func TestGrepInvertMatch(t *testing.T) {
	g, err := NewGrepper("DEBUG", GrepInvert)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !g.Match("[2024-01-01] INFO hello") {
		t.Error("expected match: line does not contain DEBUG")
	}
	if g.Match("[2024-01-01] DEBUG verbose output") {
		t.Error("expected no match: line contains DEBUG")
	}
}

func TestGrepRegexMatch(t *testing.T) {
	g, err := NewGrepper(`user_id=\d+`, GrepRegex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !g.Match("request user_id=42 processed") {
		t.Error("expected regex match")
	}
	if g.Match("request user_id=abc processed") {
		t.Error("expected no regex match")
	}
}

func TestGrepRegexInvalid(t *testing.T) {
	_, err := NewGrepper(`[invalid`, GrepRegex)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestGrepLines(t *testing.T) {
	lines := []Line{
		makeGrepLine("INFO startup complete"),
		makeGrepLine("ERROR disk full"),
		makeGrepLine("WARN low memory"),
		makeGrepLine("ERROR connection refused"),
	}
	g, _ := NewGrepper("ERROR", GrepLiteral)
	got := GrepLines(lines, g)
	if len(got) != 2 {
		t.Fatalf("expected 2 matched lines, got %d", len(got))
	}
	for _, l := range got {
		if !g.Match(l.Raw) {
			t.Errorf("line %q should have matched", l.Raw)
		}
	}
}

func TestGrepLinesEmpty(t *testing.T) {
	g, _ := NewGrepper("ERROR", GrepLiteral)
	got := GrepLines([]Line{}, g)
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d lines", len(got))
	}
}
