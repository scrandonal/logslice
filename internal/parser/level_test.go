package parser

import (
	"testing"
)

func TestParseLevelKnown(t *testing.T) {
	cases := []struct {
		input    string
		expected LogLevel
	}{
		{"[INFO] server started", LevelInfo},
		{"[DEBUG] connecting to db", LevelDebug},
		{"[WARN] high memory usage", LevelWarn},
		{"[ERROR] connection refused", LevelError},
		{"[FATAL] out of memory", LevelFatal},
		{"2024-01-01 INFO request ok", LevelInfo},
		{"no level here", LevelUnknown},
	}
	for _, tc := range cases {
		got := ParseLevel(tc.input)
		if got != tc.expected {
			t.Errorf("ParseLevel(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestLevelString(t *testing.T) {
	if LevelInfo.String() != "INFO" {
		t.Errorf("expected INFO, got %s", LevelInfo.String())
	}
	if LevelUnknown.String() != "UNKNOWN" {
		t.Errorf("expected UNKNOWN, got %s", LevelUnknown.String())
	}
}

func TestLevelFilterMatch(t *testing.T) {
	f := NewLevelFilter(LevelWarn)

	if f.Match("[DEBUG] verbose output") {
		t.Error("DEBUG should not match WARN filter")
	}
	if f.Match("[INFO] all good") {
		t.Error("INFO should not match WARN filter")
	}
	if !f.Match("[WARN] disk space low") {
		t.Error("WARN should match WARN filter")
	}
	if !f.Match("[ERROR] failed to open file") {
		t.Error("ERROR should match WARN filter")
	}
	if !f.Match("[FATAL] crash") {
		t.Error("FATAL should match WARN filter")
	}
}

func TestLevelFilterPassesUnknown(t *testing.T) {
	f := NewLevelFilter(LevelError)
	if !f.Match("plain text with no level token") {
		t.Error("unknown level lines should pass through")
	}
}

func TestFilterByLevel(t *testing.T) {
	lines := []Line{
		{Raw: "[DEBUG] step 1"},
		{Raw: "[INFO] step 2"},
		{Raw: "[WARN] step 3"},
		{Raw: "[ERROR] step 4"},
	}
	result := FilterByLevel(lines, LevelWarn)
	if len(result) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(result))
	}
	if ParseLevel(result[0].Raw) != LevelWarn {
		t.Errorf("expected first result to be WARN, got %v", ParseLevel(result[0].Raw))
	}
	if ParseLevel(result[1].Raw) != LevelError {
		t.Errorf("expected second result to be ERROR, got %v", ParseLevel(result[1].Raw))
	}
}

func TestFilterByLevelEmpty(t *testing.T) {
	result := FilterByLevel(nil, LevelInfo)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d lines", len(result))
	}
}
