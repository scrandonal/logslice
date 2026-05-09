package parser

import (
	"strings"
	"testing"
)

func TestExtractSingleField(t *testing.T) {
	fe := NewFieldExtractor([]string{"level"})
	got := fe.Extract(`[2024-01-01T00:00:00Z] level=info msg="hello world"`)
	if got["level"] != "info" {
		t.Errorf("expected info, got %q", got["level"])
	}
}

func TestExtractQuotedField(t *testing.T) {
	fe := NewFieldExtractor([]string{"msg"})
	got := fe.Extract(`level=warn msg="something went wrong" code=42`)
	if got["msg"] != "something went wrong" {
		t.Errorf("expected 'something went wrong', got %q", got["msg"])
	}
}

func TestExtractMissingField(t *testing.T) {
	fe := NewFieldExtractor([]string{"missing"})
	got := fe.Extract(`level=info msg=ok`)
	if _, ok := got["missing"]; ok {
		t.Error("expected missing field to be absent")
	}
}

func TestExtractAll(t *testing.T) {
	got := ExtractAll(`level=info msg=ok code=200`)
	if got["level"] != "info" || got["msg"] != "ok" || got["code"] != "200" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestExtractAllQuoted(t *testing.T) {
	got := ExtractAll(`level=error msg="disk full" host=srv1`)
	if got["msg"] != "disk full" {
		t.Errorf("expected 'disk full', got %q", got["msg"])
	}
	if got["host"] != "srv1" {
		t.Errorf("expected srv1, got %q", got["host"])
	}
}

func TestExtractMultipleFields(t *testing.T) {
	fe := NewFieldExtractor([]string{"level", "code"})
	got := fe.Extract(`level=warn msg="bad request" code=400`)
	if got["level"] != "warn" {
		t.Errorf("expected warn, got %q", got["level"])
	}
	if got["code"] != "400" {
		t.Errorf("expected 400, got %q", got["code"])
	}
	if _, ok := got["msg"]; ok {
		t.Error("expected msg to be absent when not in extractor fields")
	}
}

func TestFieldsJSON(t *testing.T) {
	lines := []Line{
		{Text: `level=info msg=started`},
		{Text: `level=error msg="oh no"`},
		{Text: `no fields here`},
	}
	out := FieldsJSON(lines, []string{"level", "msg"})
	if !strings.Contains(out, `"level":"info"`) {
		t.Errorf("missing level=info in JSON: %s", out)
	}
	if !strings.Contains(out, `"msg":"oh no"`) {
		t.Errorf("missing msg=oh no in JSON: %s", out)
	}
	if strings.HasPrefix(out, "[") == false {
		t.Errorf("expected JSON array, got: %s", out)
	}
}

func TestFieldsJSONEmpty(t *testing.T) {
	out := FieldsJSON([]Line{}, []string{"level"})
	if out != "[]" {
		t.Errorf("expected [], got %s", out)
	}
}

func TestPrintFields(t *testing.T) {
	lines := []Line{
		{Text: `level=info code=200`},
		{Text: `no match`},
	}
	var sb strings.Builder
	PrintFields(&sb, lines, []string{"level", "code"})
	out := sb.String()
	if !strings.Contains(out, "level=info") {
		t.Errorf("expected level=info in output: %q", out)
	}
	if !strings.Contains(out, "code=200") {
		t.Errorf("expected code=200 in output: %q", out)
	}
}
