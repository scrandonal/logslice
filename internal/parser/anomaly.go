package parser

import (
	"math"
	"time"
)

// AnomalyResult holds a log line flagged as a rate anomaly.
type AnomalyResult struct {
	Line      Line
	Bucket    time.Time
	Count     int
	Mean      float64
	StdDev    float64
	ZScore    float64
}

// AnomalyDetector finds time buckets whose line count deviates significantly
// from the mean count across all buckets (z-score threshold).
type AnomalyDetector struct {
	bucketSize time.Duration
	threshold  float64 // z-score cutoff
}

// NewAnomalyDetector creates an AnomalyDetector with the given bucket size and
// z-score threshold. A threshold of 2.0 is a reasonable default.
func NewAnomalyDetector(bucketSize time.Duration, threshold float64) *AnomalyDetector {
	if bucketSize <= 0 {
		bucketSize = time.Minute
	}
	if threshold <= 0 {
		threshold = 2.0
	}
	return &AnomalyDetector{bucketSize: bucketSize, threshold: threshold}
}

// Detect scans lines, groups them into time buckets, computes mean and standard
// deviation of per-bucket counts, and returns lines from buckets whose count
// exceeds mean + threshold*stddev.
func (a *AnomalyDetector) Detect(lines []Line) []AnomalyResult {
	if len(lines) == 0 {
		return nil
	}

	type bucket struct {
		key   time.Time
		lines []Line
	}

	order := []time.Time{}
	buckets := map[time.Time]*bucket{}

	for _, l := range lines {
		if l.Timestamp == nil {
			continue
		}
		key := l.Timestamp.Truncate(a.bucketSize)
		if _, ok := buckets[key]; !ok {
			buckets[key] = &bucket{key: key}
			order = append(order, key)
		}
		buckets[key].lines = append(buckets[key].lines, l)
	}

	if len(order) == 0 {
		return nil
	}

	counts := make([]float64, len(order))
	for i, k := range order {
		counts[i] = float64(len(buckets[k].lines))
	}

	mean, stddev := meanStddev(counts)

	var results []AnomalyResult
	for i, k := range order {
		z := 0.0
		if stddev > 0 {
			z = (counts[i] - mean) / stddev
		}
		if z >= a.threshold {
			for _, l := range buckets[k].lines {
				results = append(results, AnomalyResult{
					Line:   l,
					Bucket: k,
					Count:  len(buckets[k].lines),
					Mean:   mean,
					StdDev: stddev,
					ZScore: z,
				})
			}
		}
	}
	return results
}

func meanStddev(vals []float64) (mean, stddev float64) {
	if len(vals) == 0 {
		return 0, 0
	}
	for _, v := range vals {
		mean += v
	}
	mean /= float64(len(vals))
	var variance float64
	for _, v := range vals {
		d := v - mean
		variance += d * d
	}
	variance /= float64(len(vals))
	stddev = math.Sqrt(variance)
	return
}
