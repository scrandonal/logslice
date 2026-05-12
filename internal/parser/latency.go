package parser

import (
	"math"
	"sort"
	"time"
)

// LatencyStats holds percentile and summary statistics for a set of durations.
type LatencyStats struct {
	Count  int
	Min    time.Duration
	Max    time.Duration
	Mean   time.Duration
	P50    time.Duration
	P90    time.Duration
	P99    time.Duration
	Stddev time.Duration
}

// LatencyLine is a log line paired with an extracted duration value.
type LatencyLine struct {
	Line     Line
	Duration time.Duration
}

// CalcLatency extracts durations from lines using the given field name
// and returns a LatencyStats summary. Duration values must be parseable
// by time.ParseDuration (e.g. "12ms", "1.5s").
func CalcLatency(lines []LatencyLine) LatencyStats {
	if len(lines) == 0 {
		return LatencyStats{}
	}

	durations := make([]float64, len(lines))
	for i, l := range lines {
		durations[i] = float64(l.Duration)
	}
	sort.Float64s(durations)

	n := len(durations)
	var sum float64
	for _, d := range durations {
		sum += d
	}
	mean := sum / float64(n)

	var variance float64
	for _, d := range durations {
		diff := d - mean
		variance += diff * diff
	}
	variance /= float64(n)

	return LatencyStats{
		Count:  n,
		Min:    time.Duration(durations[0]),
		Max:    time.Duration(durations[n-1]),
		Mean:   time.Duration(mean),
		P50:    time.Duration(percentile(durations, 50)),
		P90:    time.Duration(percentile(durations, 90)),
		P99:    time.Duration(percentile(durations, 99)),
		Stddev: time.Duration(math.Sqrt(variance)),
	}
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := p / 100.0 * float64(len(sorted)-1)
	lo := int(math.Floor(idx))
	hi := int(math.Ceil(idx))
	if lo == hi {
		return sorted[lo]
	}
	frac := idx - float64(lo)
	return sorted[lo]*(1-frac) + sorted[hi]*frac
}
