package parser

import (
	"testing"
	"time"
)

func TestMaskerRedactMode(t *testing.T) {
	m, err := NewMasker([]string{`\d{4}-\d{4}-\d{4}-\d{4}`}, MaskRedact)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "card=1234-5678-9012-3456 processed"
	want := "card=[REDACTED] processed"
	got := m.Mask(input)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMaskerPartialMode(t *testing.T) {
	m, err := NewMasker([]string{`password=\S+`}, MaskPartial)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "auth password=secret123 ok"
	got := m.Mask(input)
	// original match is "password=secret123" (18 chars), partial keeps first+last
	if got == input {
		t.Error("expected line to be masked")
	}
	if got[len("auth "):len("auth ")+1] != "p" {
		t.Errorf("expected partial mask to start with 'p', got: %q", got)
	}
}

func TestMaskerMultiplePatterns(t *testing.T) {
	m, err := NewMasker([]string{
		`token=[A-Za-z0-9]+`,
		`ip=\d+\.\d+\.\d+\.\d+`,
	}, MaskRedact)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "token=abc123 ip=192.168.1.1 ok"
	want := "[REDACTED] [REDACTED] ok"
	got := m.Mask(input)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMaskerInvalidPattern(t *testing.T) {
	_, err := NewMasker([]string{`[invalid`}, MaskRedact)
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestMaskerNoMatch(t *testing.T) {
	m, err := NewMasker([]string{`secret=\S+`}, MaskRedact)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "nothing sensitive here"
	got := m.Mask(input)
	if got != input {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestMaskLinesHelper(t *testing.T) {
	now := time.Now()
	lines := []Line{
		{Raw: "user=admin pass=hunter2", Timestamp: &now},
		{Raw: "no sensitive data", Timestamp: nil},
	}
	out, err := MaskLines(lines, []string{`pass=\S+`}, MaskRedact)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	if out[0].Raw != "user=admin [REDACTED]" {
		t.Errorf("got %q", out[0].Raw)
	}
	if out[0].Timestamp != lines[0].Timestamp {
		t.Error("timestamp should be preserved")
	}
	if out[1].Raw != "no sensitive data" {
		t.Errorf("second line should be unchanged, got %q", out[1].Raw)
	}
}
