package parser

import (
	"testing"
)

func makeFreqLine(raw string) Line {
	return Line{Raw: raw}
}

func TestCalcFreqEmpty(t *testing.T) {
	result := CalcFreq(nil, FreqOptions{})
	if len(result.Entries) != 0 {
		t.Fatalf("expected no entries, got %d", len(result.Entries))
	}
	if result.Total != 0 {
		t.Fatalf("expected total 0, got %d", result.Total)
	}
}

func TestCalcFreqRawLines(t *testing.T) {
	lines := []Line{
		makeFreqLine("foo"),
		makeFreqLine("bar"),
		makeFreqLine("foo"),
		makeFreqLine("foo"),
		makeFreqLine("bar"),
	}
	result := CalcFreq(lines, FreqOptions{})
	if result.Total != 5 {
		t.Fatalf("expected total 5, got %d", result.Total)
	}
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
	if result.Entries[0].Token != "foo" || result.Entries[0].Count != 3 {
		t.Errorf("expected foo=3, got %+v", result.Entries[0])
	}
	if result.Entries[1].Token != "bar" || result.Entries[1].Count != 2 {
		t.Errorf("expected bar=2, got %+v", result.Entries[1])
	}
}

func TestCalcFreqTopN(t *testing.T) {
	lines := []Line{
		makeFreqLine("a"),
		makeFreqLine("b"),
		makeFreqLine("c"),
		makeFreqLine("a"),
		makeFreqLine("a"),
		makeFreqLine("b"),
	}
	result := CalcFreq(lines, FreqOptions{TopN: 2})
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries with TopN=2, got %d", len(result.Entries))
	}
}

func TestCalcFreqCaseInsensitive(t *testing.T) {
	lines := []Line{
		makeFreqLine("ERROR"),
		makeFreqLine("error"),
		makeFreqLine("Error"),
	}
	result := CalcFreq(lines, FreqOptions{CaseInsensitive: true})
	if len(result.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result.Entries))
	}
	if result.Entries[0].Token != "error" || result.Entries[0].Count != 3 {
		t.Errorf("unexpected entry: %+v", result.Entries[0])
	}
}

func TestCalcFreqFieldExtraction(t *testing.T) {
	lines := []Line{
		makeFreqLine(`level=info msg="started"`),
		makeFreqLine(`level=error msg="failed"`),
		makeFreqLine(`level=info msg="done"`),
	}
	result := CalcFreq(lines, FreqOptions{Field: "level"})
	if result.Total != 3 {
		t.Fatalf("expected total 3, got %d", result.Total)
	}
	if result.Entries[0].Token != "info" || result.Entries[0].Count != 2 {
		t.Errorf("expected info=2, got %+v", result.Entries[0])
	}
}

func TestCalcFreqFieldMissing(t *testing.T) {
	lines := []Line{
		makeFreqLine("no fields here"),
		makeFreqLine("also nothing"),
	}
	result := CalcFreq(lines, FreqOptions{Field: "level"})
	if result.Total != 0 {
		t.Fatalf("expected total 0 when field missing, got %d", result.Total)
	}
}
