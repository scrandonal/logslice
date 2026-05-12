package parser

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func buildLatencyStats() LatencyStats {
	return CalcLatency([]LatencyLine{
		makeLatencyLine(10 * time.Millisecond),
		makeLatencyLine(50 * time.Millisecond),
		makeLatencyLine(100 * time.Millisecond),
		makeLatencyLine(200 * time.Millisecond),
		makeLatencyLine(500 * time.Millisecond),
	})
}

func TestPrintLatencyEmpty(t *testing.T) {
	var buf bytes.Buffer
	printLatencyTo(&buf, LatencyStats{})
	if !strings.Contains(buf.String(), "no latency data") {
		t.Errorf("expected 'no latency data', got: %s", buf.String())
	}
}

func TestPrintLatencyFields(t *testing.T) {
	var buf bytes.Buffer
	printLatencyTo(&buf, buildLatencyStats())
	out := buf.String()
	for _, label := range []string{"count", "min", "max", "mean", "p50", "p90", "p99", "stddev"} {
		if !strings.Contains(out, label) {
			t.Errorf("output missing label %q", label)
		}
	}
}

func TestLatencyJSONEmpty(t *testing.T) {
	out, err := LatencyJSON(LatencyStats{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"count\":0") {
		t.Errorf("expected count:0 in JSON, got: %s", out)
	}
}

func TestLatencyJSONValues(t *testing.T) {
	stats := buildLatencyStats()
	out, err := LatencyJSON(stats)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"count\":5") {
		t.Errorf("expected count:5, got: %s", out)
	}
	if !strings.Contains(out, "\"min_ms\":10") {
		t.Errorf("expected min_ms:10, got: %s", out)
	}
	if !strings.Contains(out, "\"max_ms\":500") {
		t.Errorf("expected max_ms:500, got: %s", out)
	}
}
