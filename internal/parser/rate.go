package parser

import (
	"time"
)

// RatePoint holds a timestamp bucket and the count of lines in that bucket.
type RatePoint struct {
	Bucket    time.Time
	Count     int
	FirstLine *Line
}

// RateCalculator buckets log lines by a fixed time interval and counts them.
type RateCalculator struct {
	interval time.Duration
	buckets  map[int64]*RatePoint
	order    []int64
}

// NewRateCalculator creates a RateCalculator with the given bucket interval.
func NewRateCalculator(interval time.Duration) *RateCalculator {
	if interval <= 0 {
		interval = time.Minute
	}
	return &RateCalculator{
		interval: interval,
		buckets:  make(map[int64]*RatePoint),
	}
}

// Add inserts a line into the appropriate time bucket.
func (r *RateCalculator) Add(l Line) {
	var key int64
	if l.Timestamp != nil {
		key = l.Timestamp.Truncate(r.interval).UnixNano()
	}
	if _, ok := r.buckets[key]; !ok {
		var t time.Time
		if l.Timestamp != nil {
			t = l.Timestamp.Truncate(r.interval)
		}
		r.buckets[key] = &RatePoint{Bucket: t, FirstLine: &l}
		r.order = append(r.order, key)
	}
	r.buckets[key].Count++
}

// Points returns the rate points in insertion order.
func (r *RateCalculator) Points() []RatePoint {
	out := make([]RatePoint, 0, len(r.order))
	for _, k := range r.order {
		out = append(out, *r.buckets[k])
	}
	return out
}

// CalcRate buckets a slice of lines and returns the resulting RatePoints.
func CalcRate(lines []Line, interval time.Duration) []RatePoint {
	r := NewRateCalculator(interval)
	for _, l := range lines {
		r.Add(l)
	}
	return r.Points()
}
