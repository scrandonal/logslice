package parser

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func buildPivotResult() *PivotResult {
	return NewPivot(
		[]Line{
			makePivotLine(`level=info`, ptrPivotTime("2024-03-01T09:00:00Z")),
			makePivotLine(`level=info`, ptrPivotTime("2024-03-01T09:00:45Z")),
			makePivotLine(`level=error`, ptrPivotTime("2024-03-01T09:01:10Z")),
		},
		PivotOptions{Field: "level", BucketSize: time.Minute},
	)
}

func TestPrintPivotEmpty(t *testing.T) {
	var sb strings.Builder
	PrintPivot(&sb, &PivotResult{Field: "level"})
	if !strings.Contains(sb.String(), "(no pivot data)") {
		t.Errorf("expected no-data message, got: %q", sb.String())
	}
}

func TestPrintPivotHeader(t *testing.T) {
	var sb strings.Builder
	PrintPivot(&sb, buildPivotResult())
	out := sb.String()
	if !strings.Contains(out, "level") {
		t.Errorf("expected field name in output, got: %q", out)
	}
	if !strings.Contains(out, "total") {
		t.Errorf("expected 'total' column header, got: %q", out)
	}
	if !strings.Contains(out, "info") {
		t.Errorf("expected 'info' row, got: %q", out)
	}
}

func TestPivotJSONStructure(t *testing.T) {
	result := buildPivotResult()
	s, err := PivotJSON(result)
	if err != nil {
		t.Fatalf("PivotJSON error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if obj["field"] != "level" {
		t.Errorf("expected field=level in JSON, got %v", obj["field"])
	}
	rows, ok := obj["rows"].([]interface{})
	if !ok || len(rows) == 0 {
		t.Errorf("expected non-empty rows array")
	}
}

func TestPivotJSONEmpty(t *testing.T) {
	s, err := PivotJSON(&PivotResult{Field: "svc"})
	if err != nil {
		t.Fatalf("PivotJSON error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	rows := obj["rows"].([]interface{})
	if len(rows) != 0 {
		t.Errorf("expected empty rows, got %d", len(rows))
	}
}
