package parser

import (
	"testing"
	"time"
)

func makeLatencyLine(d time.Duration) LatencyLine {
	return LatencyLine{Line: Line{Raw: "entry"}, Duration: d}
}

func TestCalcLatencyEmpty(t *testing.T) {
	stats := CalcLatency(nil)
	if stats.Count != 0 {
		t.Fatalf("expected count 0, got %d", stats.Count)
	}
}

func TestCalcLatencySingle(t *testing.T) {
	stats := CalcLatency([]LatencyLine{makeLatencyLine(100 * time.Millisecond)})
	if stats.Count != 1 {
		t.Fatalf("expected count 1, got %d", stats.Count)
	}
	if stats.Min != 100*time.Millisecond {
		t.Errorf("min: got %s, want 100ms", stats.Min)
	}
	if stats.Max != 100*time.Millisecond {
		t.Errorf("max: got %s, want 100ms", stats.Max)
	}
	if stats.Stddev != 0 {
		t.Errorf("stddev should be 0 for single value, got %s", stats.Stddev)
	}
}

func TestCalcLatencyPercentiles(t *testing.T) {
	lines := []LatencyLine{
		makeLatencyLine(10 * time.Millisecond),
		makeLatencyLine(20 * time.Millisecond),
		makeLatencyLine(30 * time.Millisecond),
		makeLatencyLine(40 * time.Millisecond),
		makeLatencyLine(200 * time.Millisecond),
	}
	stats := CalcLatency(lines)
	if stats.Count != 5 {
		t.Fatalf("expected count 5, got %d", stats.Count)
	}
	if stats.Min != 10*time.Millisecond {
		t.Errorf("min: got %s", stats.Min)
	}
	if stats.Max != 200*time.Millisecond {
		t.Errorf("max: got %s", stats.Max)
	}
	if stats.P50 < 25*time.Millisecond || stats.P50 > 35*time.Millisecond {
		t.Errorf("p50 out of expected range: %s", stats.P50)
	}
	if stats.P99 < 150*time.Millisecond {
		t.Errorf("p99 should be near max: %s", stats.P99)
	}
}

func TestCalcLatencyMean(t *testing.T) {
	lines := []LatencyLine{
		makeLatencyLine(10 * time.Millisecond),
		makeLatencyLine(20 * time.Millisecond),
		makeLatencyLine(30 * time.Millisecond),
	}
	stats := CalcLatency(lines)
	if stats.Mean != 20*time.Millisecond {
		t.Errorf("mean: got %s, want 20ms", stats.Mean)
	}
}

func TestLatencyJSON(t *testing.T) {
	lines := []LatencyLine{
		makeLatencyLine(50 * time.Millisecond),
		makeLatencyLine(100 * time.Millisecond),
	}
	stats := CalcLatency(lines)
	out, err := LatencyJSON(stats)
	if err != nil {
		t.Fatalf("LatencyJSON error: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected non-empty JSON")
	}
	for _, key := range []string{"count", "min_ms", "max_ms", "mean_ms", "p50_ms", "p90_ms", "p99_ms"} {
		if !containsStr(out, key) {
			t.Errorf("JSON missing key %q", key)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstr(s, sub))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
