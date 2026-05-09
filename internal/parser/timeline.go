package parser

import (
	"time"
)

// TimelineBucket holds a time bucket label and the count of log lines in it.
type TimelineBucket struct {
	Label string
	Start time.Time
	Count int
}

// TimelineResult holds all buckets produced by BuildTimeline.
type TimelineResult struct {
	Buckets  []TimelineBucket
	Interval time.Duration
	Total    int
}

// BuildTimeline groups lines into fixed-duration buckets and returns a
// TimelineResult suitable for ASCII or JSON rendering.
func BuildTimeline(lines []Line, interval time.Duration) TimelineResult {
	if interval <= 0 {
		interval = time.Minute
	}

	result := TimelineResult{Interval: interval}
	if len(lines) == 0 {
		return result
	}

	buckets := make(map[int64]*TimelineBucket)
	var keys []int64

	for _, l := range lines {
		result.Total++
		if l.Timestamp == nil {
			continue
		}
		ts := l.Timestamp.Truncate(interval)
		key := ts.UnixNano()
		if _, ok := buckets[key]; !ok {
			buckets[key] = &TimelineBucket{
				Label: ts.Format("2006-01-02 15:04:05"),
				Start: ts,
			}
			keys = append(keys, key)
		}
		buckets[key].Count++
	}

	// Sort keys ascending.
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}

	for _, k := range keys {
		result.Buckets = append(result.Buckets, *buckets[k])
	}
	return result
}
