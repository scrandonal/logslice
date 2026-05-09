package parser

import (
	"testing"
	"time"
)

func makeAnomalyLine(ts time.Time, text string) Line {
	return Line{Timestamp: &ts, Raw: text}
}

func anomalyTime(minuteOffset int) time.Time {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return base.Add(time.Duration(minuteOffset) * time.Minute)
}

func TestAnomalyDetectorEmpty(t *testing.T) {
	det := NewAnomalyDetector(time.Minute, 2.0)
	results := det.Detect(nil)
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestAnomalyDetectorNoTimestamp(t *testing.T) {
	det := NewAnomalyDetector(time.Minute, 2.0)
	lines := []Line{
		{Raw: "no timestamp"},
		{Raw: "also no timestamp"},
	}
	results := det.Detect(lines)
	if len(results) != 0 {
		t.Fatalf("expected no results for lines without timestamps, got %d", len(results))
	}
}

func TestAnomalyDetectorBelowThreshold(t *testing.T) {
	det := NewAnomalyDetector(time.Minute, 2.0)
	// Uniform distribution — no bucket should be anomalous.
	var lines []Line
	for m := 0; m < 5; m++ {
		for i := 0; i < 3; i++ {
			lines = append(lines, makeAnomalyLine(anomalyTime(m), "msg"))
		}
	}
	results := det.Detect(lines)
	if len(results) != 0 {
		t.Fatalf("expected no anomalies for uniform distribution, got %d", len(results))
	}
}

func TestAnomalyDetectorDetectsSpike(t *testing.T) {
	det := NewAnomalyDetector(time.Minute, 2.0)
	// Minutes 0-4: 2 lines each. Minute 5: 20 lines (spike).
	var lines []Line
	for m := 0; m < 5; m++ {
		for i := 0; i < 2; i++ {
			lines = append(lines, makeAnomalyLine(anomalyTime(m), "normal"))
		}
	}
	for i := 0; i < 20; i++ {
		lines = append(lines, makeAnomalyLine(anomalyTime(5), "spike"))
	}
	results := det.Detect(lines)
	if len(results) == 0 {
		t.Fatal("expected anomaly results for spike bucket")
	}
	for _, r := range results {
		if r.ZScore < 2.0 {
			t.Errorf("expected z-score >= 2.0, got %.2f", r.ZScore)
		}
		if r.Count != 20 {
			t.Errorf("expected spike bucket count 20, got %d", r.Count)
		}
	}
}

func TestAnomalyDetectorZScoreFields(t *testing.T) {
	det := NewAnomalyDetector(time.Minute, 1.5)
	var lines []Line
	for i := 0; i < 1; i++ {
		lines = append(lines, makeAnomalyLine(anomalyTime(0), "low"))
	}
	for i := 0; i < 50; i++ {
		lines = append(lines, makeAnomalyLine(anomalyTime(1), "high"))
	}
	results := det.Detect(lines)
	if len(results) == 0 {
		t.Fatal("expected at least one anomaly result")
	}
	r := results[0]
	if r.Mean <= 0 {
		t.Errorf("mean should be positive, got %.2f", r.Mean)
	}
	if r.StdDev <= 0 {
		t.Errorf("stddev should be positive, got %.2f", r.StdDev)
	}
	if r.ZScore <= 0 {
		t.Errorf("z-score should be positive, got %.2f", r.ZScore)
	}
}
