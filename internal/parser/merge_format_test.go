package parser

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPrintMerged(t *testing.T) {
	lines := []MergedLine{
		{Line: makeLine("hello world", makeTime("2024-03-01T08:00:00Z")), Source: 0},
		{Line: makeLine("no timestamp", nil), Source: 1},
	}

	var buf bytes.Buffer
	PrintMerged(&buf, lines)
	out := buf.String()

	if !strings.Contains(out, "[src:0]") {
		t.Error("expected source index 0 in output")
	}
	if !strings.Contains(out, "2024-03-01T08:00:00Z") {
		t.Error("expected timestamp in output")
	}
	if !strings.Contains(out, "[src:1]") {
		t.Error("expected source index 1 in output")
	}
	if !strings.Contains(out, "[-]") {
		t.Error("expected dash for nil timestamp")
	}
}

func TestMergedJSON(t *testing.T) {
	lines := []MergedLine{
		{Line: makeLine("entry one", makeTime("2024-03-01T09:00:00Z")), Source: 0},
		{Line: makeLine("entry two", makeTime("2024-03-01T09:01:00Z")), Source: 1},
	}

	out, err := MergedJSON(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(parsed) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(parsed))
	}
	if parsed[0]["raw"] != "entry one" {
		t.Errorf("unexpected raw value: %v", parsed[0]["raw"])
	}
	if parsed[1]["source"].(float64) != 1 {
		t.Errorf("unexpected source: %v", parsed[1]["source"])
	}
}

func TestMergedJSONEmpty(t *testing.T) {
	out, err := MergedJSON(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[]") {
		t.Errorf("expected empty JSON array, got: %s", out)
	}
}
