package parser

import (
	"time"
)

// BurstWindow holds a detected burst of log activity.
type BurstWindow struct {
	Start    time.Time
	End      time.Time
	Count    int
	Lines    []Line
}

// BurstDetector finds time windows where log volume exceeds a threshold.
type BurstDetector struct {
	windowSize time.Duration
	threshold  int
}

// NewBurstDetector creates a BurstDetector with the given window size and
// minimum line count threshold to qualify as a burst.
func NewBurstDetector(windowSize time.Duration, threshold int) *BurstDetector {
	if windowSize <= 0 {
		windowSize = time.Minute
	}
	if threshold <= 0 {
		threshold = 1
	}
	return &BurstDetector{windowSize: windowSize, threshold: threshold}
}

// Detect scans the provided lines and returns all burst windows where the
// number of log entries within windowSize meets or exceeds the threshold.
// Lines without a timestamp are skipped for windowing purposes.
func (b *BurstDetector) Detect(lines []Line) []BurstWindow {
	// collect only timestamped lines
	tsLines := make([]Line, 0, len(lines))
	for _, l := range lines {
		if l.Timestamp != nil {
			tsLines = append(tsLines, l)
		}
	}
	if len(tsLines) == 0 {
		return nil
	}

	var bursts []BurstWindow
	i := 0
	for i < len(tsLines) {
		start := *tsLines[i].Timestamp
		end := start.Add(b.windowSize)
		j := i
		for j < len(tsLines) && !(*tsLines[j].Timestamp).After(end) {
			j++
		}
		count := j - i
		if count >= b.threshold {
			bursts = append(bursts, BurstWindow{
				Start: start,
				End:   *tsLines[j-1].Timestamp,
				Count: count,
				Lines: tsLines[i:j],
			})
		}
		i++
	}
	return bursts
}
