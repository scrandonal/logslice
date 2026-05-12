package parser

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func buildSequenceMatches() []SequenceMatch {
	t1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 1, 10, 0, 5, 0, time.UTC)
	elapsed := t2.Sub(t1)
	return []SequenceMatch{
		{
			Steps:   []Line{{Raw: "START"}, {Raw: "END"}},
			Start:   &t1,
			End:     &t2,
			Elapsed: &elapsed,
		},
	}
}

func TestPrintSequencesEmpty(t *testing.T) {
	var buf bytes.Buffer
	printSequencesTo(&buf, nil)
	if !strings.Contains(buf.String(), "no sequences") {
		t.Errorf("expected 'no sequences' message, got: %s", buf.String())
	}
}

func TestPrintSequencesSingle(t *testing.T) {
	var buf bytes.Buffer
	matches := buildSequenceMatches()
	printSequencesTo(&buf, matches)
	out := buf.String()
	if !strings.Contains(out, "sequence 1") {
		t.Error("expected sequence header")
	}
	if !strings.Contains(out, "START") || !strings.Contains(out, "END") {
		t.Error("expected step lines in output")
	}
	if !strings.Contains(out, "5s") {
		t.Error("expected elapsed time in output")
	}
}

func TestSequencesJSONStructure(t *testing.T) {
	matches := buildSequenceMatches()
	out, err := SequencesJSON(matches)
	if err != nil {
		t.Fatal(err)
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &arr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(arr) != 1 {
		t.Fatalf("expected 1 element, got %d", len(arr))
	}
	if _, ok := arr[0]["start"]; !ok {
		t.Error("expected 'start' field")
	}
	if _, ok := arr[0]["elapsed"]; !ok {
		t.Error("expected 'elapsed' field")
	}
	steps, ok := arr[0]["steps"].([]interface{})
	if !ok || len(steps) != 2 {
		t.Errorf("expected 2 steps, got %v", arr[0]["steps"])
	}
}

func TestSequencesJSONEmpty(t *testing.T) {
	out, err := SequencesJSON(nil)
	if err != nil {
		t.Fatal(err)
	}
	if out != "[]" && out != "null" {
		t.Errorf("unexpected output for empty: %s", out)
	}
}
