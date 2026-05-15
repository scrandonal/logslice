package parser

import (
	"testing"
	"time"
)

func makeNormLine(raw string) Line {
	return Line{Raw: raw}
}

func TestNormalizerInvalidPattern(t *testing.T) {
	_, err := NewNormalizer([][2]string{{`[invalid`, "x"}})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNormalizerApplyUUID(t *testing.T) {
	n, err := NewNormalizer([][2]string{
		{`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`, "<UUID>"},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := n.Apply("request id=550e8400-e29b-41d4-a716-446655440000 done")
	want := "request id=<UUID> done"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestNormalizerApplyIP(t *testing.T) {
	n, err := NewNormalizer([][2]string{
		{`\b(?:\d{1,3}\.){3}\d{1,3}\b`, "<IP>"},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := n.Apply("connection from 192.168.1.42 accepted")
	want := "connection from <IP> accepted"
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestNormalizerMultipleRules(t *testing.T) {
	n := DefaultNormalizer()
	got := n.Apply("user 550e8400-e29b-41d4-a716-446655440000 from 10.0.0.1 took 12345 ms")
	if got != "user <UUID> from <IP> took <NUM> ms" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestNormalizeLinePreservesTimestamp(t *testing.T) {
	n := DefaultNormalizer()
	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	l := Line{Raw: "error port 54321 unreachable", Timestamp: &ts}
	out := n.NormalizeLine(l)
	if out.Timestamp == nil || !out.Timestamp.Equal(ts) {
		t.Error("timestamp should be preserved")
	}
	if out.Raw != "error port <NUM> unreachable" {
		t.Errorf("unexpected raw: %q", out.Raw)
	}
}

func TestNormalizeLinesHelper(t *testing.T) {
	lines := []Line{
		makeNormLine("id=1001 status=ok"),
		makeNormLine("id=2002 status=fail"),
	}
	out, err := NormalizeLines(lines, [][2]string{{`\b\d+\b`, "<NUM>"}})
	if err != nil {
		t.Fatal(err)
	}
	if out[0].Raw != "id=<NUM> status=ok" {
		t.Errorf("line 0: %q", out[0].Raw)
	}
	if out[1].Raw != "id=<NUM> status=fail" {
		t.Errorf("line 1: %q", out[1].Raw)
	}
}

func TestNormalizeLinesInvalidPattern(t *testing.T) {
	_, err := NormalizeLines(nil, [][2]string{{`[bad`, "x"}})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNormalizedKey(t *testing.T) {
	got := NormalizedKey("  Hello   World  ")
	if got != "hello world" {
		t.Errorf("got %q", got)
	}
}
